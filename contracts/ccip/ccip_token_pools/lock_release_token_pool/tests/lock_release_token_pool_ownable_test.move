#[test_only]
module lock_release_token_pool::lock_release_token_pool_ownable_test {
    use sui::coin;
    use sui::test_scenario::{Self as ts, Scenario};

    use ccip::state_object::{Self, OwnerCap as CCIPOwnerCap, CCIPObjectRef};
    use ccip::dynamic_dispatcher;
    use ccip::offramp_state_helper;
    use ccip::rmn_remote;
    use ccip::token_admin_registry;

    use lock_release_token_pool::lock_release_token_pool::{Self, LockReleaseTokenPoolState};
    use ccip_token_pool::ownable::{Self, OwnerCap};

    const OWNER: address = @0x123;
    const NEW_OWNER: address = @0x456;
    const OTHER_USER: address = @0x789;
    const CCIP_ADMIN: address = @0x400;
    const TOKEN_ADMIN: address = @0x200;
    const REBALANCER: address = @0x100;

    public struct LOCK_RELEASE_TOKEN_POOL_OWNABLE_TEST has drop {}

    public struct TestEnv {
        scenario: Scenario,
        state: LockReleaseTokenPoolState<LOCK_RELEASE_TOKEN_POOL_OWNABLE_TEST>,
        ccip_ref: CCIPObjectRef,
    }

    fun setup(): (TestEnv, OwnerCap) {
        let mut scenario = ts::begin(OWNER);
        
        // Setup CCIP environment
        scenario.next_tx(CCIP_ADMIN);
        let ctx = scenario.ctx();
        state_object::test_init(ctx);

        scenario.next_tx(CCIP_ADMIN);
        let ccip_owner_cap = scenario.take_from_sender<CCIPOwnerCap>();
        let mut ccip_ref = scenario.take_shared<CCIPObjectRef>();

        // Initialize required CCIP modules
        rmn_remote::initialize(&mut ccip_ref, &ccip_owner_cap, 1000, scenario.ctx());
        token_admin_registry::initialize(&mut ccip_ref, &ccip_owner_cap, scenario.ctx());
        dynamic_dispatcher::test_init(scenario.ctx());
        offramp_state_helper::test_init(scenario.ctx());

        // Initialize token pool using the existing test token from the main test module
        scenario.next_tx(OWNER);
        let ctx = scenario.ctx();
        let (treasury_cap, coin_metadata) = coin::create_currency(
            LOCK_RELEASE_TOKEN_POOL_OWNABLE_TEST {},
            8, // decimals
            b"TEST",
            b"TestToken",
            b"test_token",
            option::none(),
            ctx
        );

        lock_release_token_pool::initialize(
            &mut ccip_ref,
            &coin_metadata,
            &treasury_cap,
            TOKEN_ADMIN,
            REBALANCER,
            ctx
        );

        transfer::public_freeze_object(coin_metadata);
        transfer::public_transfer(treasury_cap, ctx.sender());
        transfer::public_transfer(ccip_owner_cap, @0x0);

        scenario.next_tx(OWNER);
        let state = ts::take_shared<LockReleaseTokenPoolState<LOCK_RELEASE_TOKEN_POOL_OWNABLE_TEST>>(&scenario);
        let owner_cap = ts::take_from_sender<OwnerCap>(&scenario);

        let env = TestEnv {
            scenario,
            state,
            ccip_ref,
        };

        (env, owner_cap)
    }

    fun tear_down(env: TestEnv) {
        let TestEnv { scenario, state, ccip_ref } = env;

        ts::return_shared(state);
        ts::return_shared(ccip_ref);
        ts::end(scenario);
    }

    #[test]
    public fun test_basic_ownable_functionality() {
        let (env, owner_cap) = setup();

        // Test basic ownership queries
        assert!(lock_release_token_pool::owner(&env.state) == OWNER);
        assert!(!lock_release_token_pool::has_pending_transfer(&env.state));
        assert!(lock_release_token_pool::pending_transfer_from(&env.state).is_none());
        assert!(lock_release_token_pool::pending_transfer_to(&env.state).is_none());
        assert!(lock_release_token_pool::pending_transfer_accepted(&env.state).is_none());

        tear_down(env);
        ts::return_to_address(OWNER, owner_cap);
    }

    #[test]
    public fun test_ownership_transfer_flow() {
        let (mut env, owner_cap) = setup();

        // Transfer ownership
        lock_release_token_pool::transfer_ownership(&mut env.state, &owner_cap, NEW_OWNER, env.scenario.ctx());

        // Verify pending transfer
        assert!(lock_release_token_pool::has_pending_transfer(&env.state));
        assert!(lock_release_token_pool::pending_transfer_from(&env.state).extract() == OWNER);
        assert!(lock_release_token_pool::pending_transfer_to(&env.state).extract() == NEW_OWNER);
        assert!(lock_release_token_pool::pending_transfer_accepted(&env.state).extract() == false);

        // Accept ownership as new owner
        env.scenario.next_tx(NEW_OWNER);
        lock_release_token_pool::accept_ownership(&mut env.state, env.scenario.ctx());

        // Verify the transfer is accepted
        assert!(lock_release_token_pool::pending_transfer_accepted(&env.state).extract() == true);

        tear_down(env);
        ts::return_to_address(OWNER, owner_cap);
    }

    #[test]
    public fun test_accept_ownership_from_object() {
        let (mut env, owner_cap) = setup();

        // Create a test object that represents NEW_OWNER
        let mut test_obj = object::new(env.scenario.ctx());
        let test_obj_address = test_obj.to_address();

        // Transfer ownership to the object's address instead of NEW_OWNER
        lock_release_token_pool::transfer_ownership(&mut env.state, &owner_cap, test_obj_address, env.scenario.ctx());

        // Accept ownership from the object
        env.scenario.next_tx(NEW_OWNER); // Still use NEW_OWNER as the transaction sender
        lock_release_token_pool::accept_ownership_from_object(&mut env.state, &mut test_obj, env.scenario.ctx());

        // Verify the transfer is accepted
        assert!(lock_release_token_pool::pending_transfer_accepted(&env.state).extract() == true);

        // Clean up the test object
        test_obj.delete();

        tear_down(env);
        ts::return_to_address(OWNER, owner_cap);
    }

    #[test]
    #[expected_failure(abort_code = ownable::EUnauthorizedOwnershipTransfer)]
    public fun test_transfer_ownership_unauthorized() {
        let (mut env, owner_cap) = setup();

        // Try to transfer ownership from unauthorized user
        env.scenario.next_tx(OTHER_USER);
        lock_release_token_pool::transfer_ownership(&mut env.state, &owner_cap, NEW_OWNER, env.scenario.ctx());

        tear_down(env);
        ts::return_to_address(OWNER, owner_cap);
    }

    #[test]
    #[expected_failure(abort_code = ownable::EUnauthorizedAcceptance)]
    public fun test_accept_ownership_unauthorized() {
        let (mut env, owner_cap) = setup();

        // Transfer ownership
        lock_release_token_pool::transfer_ownership(&mut env.state, &owner_cap, NEW_OWNER, env.scenario.ctx());

        // Try to accept ownership from unauthorized user
        env.scenario.next_tx(OTHER_USER);
        lock_release_token_pool::accept_ownership(&mut env.state, env.scenario.ctx());

        tear_down(env);
        ts::return_to_address(OWNER, owner_cap);
    }

    #[test]
    #[expected_failure(abort_code = ownable::ENoPendingTransfer)]
    public fun test_accept_ownership_no_pending_transfer() {
        let (mut env, owner_cap) = setup();

        // Try to accept ownership without pending transfer
        env.scenario.next_tx(NEW_OWNER);
        lock_release_token_pool::accept_ownership(&mut env.state, env.scenario.ctx());

        tear_down(env);
        ts::return_to_address(OWNER, owner_cap);
    }

    #[test]
    #[expected_failure(abort_code = ownable::ETransferAlreadyAccepted)]
    public fun test_accept_ownership_already_accepted() {
        let (mut env, owner_cap) = setup();

        // Transfer ownership
        lock_release_token_pool::transfer_ownership(&mut env.state, &owner_cap, NEW_OWNER, env.scenario.ctx());

        // Accept ownership
        env.scenario.next_tx(NEW_OWNER);
        lock_release_token_pool::accept_ownership(&mut env.state, env.scenario.ctx());

        // Try to accept again (should fail)
        lock_release_token_pool::accept_ownership(&mut env.state, env.scenario.ctx());

        tear_down(env);
        ts::return_to_address(OWNER, owner_cap);
    }

    #[test]
    public fun test_ownable_functions_with_owner_cap_validation() {
        let (mut env, owner_cap) = setup();

        // Test that owner cap validation works in other functions
        // These should succeed because we have the correct owner cap
        lock_release_token_pool::set_rebalancer(&owner_cap, &mut env.state, @0x999);
        assert!(lock_release_token_pool::get_rebalancer(&env.state) == @0x999);

        lock_release_token_pool::set_allowlist_enabled(&mut env.state, &owner_cap, true);
        assert!(lock_release_token_pool::get_allowlist_enabled(&env.state) == true);

        // Test balance check to ensure the ownable integration doesn't break token functionality
        assert!(lock_release_token_pool::get_balance<LOCK_RELEASE_TOKEN_POOL_OWNABLE_TEST>(&env.state) == 0);

        tear_down(env);
        ts::return_to_address(OWNER, owner_cap);
    }
} 