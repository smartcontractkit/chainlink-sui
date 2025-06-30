#[test_only]
module ccip_router::router_ownable_test {
    use ccip_router::router::{Self, RouterState};
    use ccip_router::ownable::OwnerCap;
    use sui::test_scenario::{Self as ts, Scenario};

    const OWNER: address = @0x123;
    const NEW_OWNER: address = @0x456;
    const OTHER_USER: address = @0x789;

    public struct TestEnv {
        scenario: Scenario,
        state: RouterState,
    }

    fun setup(): (TestEnv, OwnerCap) {
        let mut scenario = ts::begin(OWNER);
        let ctx = scenario.ctx();

        // Initialize router
        router::test_init(ctx);

        scenario.next_tx(OWNER);

        let state = ts::take_shared<RouterState>(&scenario);
        let owner_cap = ts::take_from_sender<OwnerCap>(&scenario);

        let env = TestEnv {
            scenario,
            state,
        };

        (env, owner_cap)
    }

    fun tear_down(env: TestEnv) {
        let TestEnv { scenario, state } = env;

        ts::return_shared(state);
        ts::end(scenario);
    }

    #[test]
    public fun test_basic_ownable_functionality() {
        let (env, owner_cap) = setup();

        // Test owner function
        let current_owner = router::owner(&env.state);
        assert!(current_owner == OWNER);

        // Test ownership query functions
        assert!(!router::has_pending_transfer(&env.state));
        assert!(router::pending_transfer_from(&env.state).is_none());
        assert!(router::pending_transfer_to(&env.state).is_none());
        assert!(router::pending_transfer_accepted(&env.state).is_none());

        tear_down(env);
        ts::return_to_address(OWNER, owner_cap);
    }

    #[test]
    public fun test_ownership_transfer_flow() {
        let (mut env, owner_cap) = setup();

        // Test ownership transfer initiation
        router::transfer_ownership(&mut env.state, &owner_cap, NEW_OWNER, env.scenario.ctx());

        // Verify pending transfer state
        assert!(router::has_pending_transfer(&env.state));
        assert!(router::pending_transfer_from(&env.state) == option::some(OWNER));
        assert!(router::pending_transfer_to(&env.state) == option::some(NEW_OWNER));
        assert!(router::pending_transfer_accepted(&env.state) == option::some(false));

        // Owner should still be the same until transfer is accepted
        assert!(router::owner(&env.state) == OWNER);

        // Accept ownership transfer
        env.scenario.next_tx(NEW_OWNER);
        router::accept_ownership(&mut env.state, env.scenario.ctx());

        // Verify the transfer is accepted but not yet completed
        assert!(router::pending_transfer_accepted(&env.state) == option::some(true));

        tear_down(env);
        ts::return_to_address(OWNER, owner_cap);
    }

    #[test]
    public fun test_owner_authorization_for_set_on_ramp_infos() {
        let (mut env, owner_cap) = setup();

        // Test that owner can set on ramp infos
        let dest_chain_selector = 1000;
        let on_ramp_address = @0x111;
        let on_ramp_version = vector[1, 6, 0];

        router::set_on_ramp_infos(
            &owner_cap,
            &mut env.state,
            vector[dest_chain_selector],
            vector[on_ramp_address],
            vector[on_ramp_version]
        );

        // Verify the on ramp info was set
        assert!(router::is_chain_supported(&env.state, dest_chain_selector));
        let (address, version) = router::get_on_ramp_info(&env.state, dest_chain_selector);
        assert!(address == on_ramp_address);
        assert!(version == on_ramp_version);

        tear_down(env);
        ts::return_to_address(OWNER, owner_cap);
    }

    #[test]
    #[expected_failure(abort_code = router::EInvalidOwnerCap)]
    public fun test_unauthorized_set_on_ramp_infos() {
        let (mut env, _owner_cap) = setup();

        // Create a fake owner cap from a different router
        env.scenario.next_tx(OTHER_USER);
        let ctx = env.scenario.ctx();
        router::test_init(ctx);

        env.scenario.next_tx(OTHER_USER);
        let fake_owner_cap = ts::take_from_sender<OwnerCap>(&env.scenario);
        let fake_state = ts::take_shared<RouterState>(&env.scenario);

        // Try to use fake owner cap on original state - should fail
        router::set_on_ramp_infos(
            &fake_owner_cap,
            &mut env.state,
            vector[1000],
            vector[@0x111],
            vector[vector[1, 6, 0]]
        );

        // Clean up (should not be reached due to expected failure)
        ts::return_to_address(OTHER_USER, fake_owner_cap);
        ts::return_shared(fake_state);
        tear_down(env);
        ts::return_to_address(OWNER, _owner_cap);
    }

    #[test]
    #[expected_failure(abort_code = ccip_router::ownable::ECannotTransferToSelf)]
    public fun test_cannot_transfer_to_self() {
        let (mut env, owner_cap) = setup();

        // Try to transfer ownership to self - should fail
        router::transfer_ownership(&mut env.state, &owner_cap, OWNER, env.scenario.ctx());

        tear_down(env);
        ts::return_to_address(OWNER, owner_cap);
    }

    #[test]
    #[expected_failure(abort_code = ccip_router::ownable::EMustBeProposedOwner)]
    public fun test_unauthorized_accept_ownership() {
        let (mut env, owner_cap) = setup();

        // Transfer ownership to NEW_OWNER
        router::transfer_ownership(&mut env.state, &owner_cap, NEW_OWNER, env.scenario.ctx());

        // Try to accept ownership from wrong address - should fail
        env.scenario.next_tx(OTHER_USER); // Not NEW_OWNER
        router::accept_ownership(&mut env.state, env.scenario.ctx());

        tear_down(env);
        ts::return_to_address(OWNER, owner_cap);
    }
} 