#[test_only]
module ccip_onramp::onramp_test {
    use ccip_onramp::onramp::{Self, OnRampState};
    use ccip_onramp::ownable::{OwnerCap};
    use ccip::state_object::{Self, CCIPObjectRef};
    use sui::test_scenario::{Self as ts, Scenario};
    use ccip::dynamic_dispatcher as dd;
    use ccip::nonce_manager::{Self, NonceManagerCap};

    use mcms::mcms_registry::{Self, Registry};
    use mcms::mcms_account::{Self};
    use mcms::mcms_deployer::{Self};

    const DEST_CHAIN_SELECTOR_1: u64 = 1;
    const DEST_CHAIN_SELECTOR_2: u64 = 2;
    const ALLOWED_SENDER_1: address = @0x11;
    const ALLOWED_SENDER_2: address = @0x22;
    const ALLOWED_SENDER_3: address = @0x33;
    const OWNER: address = @0x123;

    public struct Env {
        scenario: Scenario,
        state: OnRampState,
        registry: Registry,
        ref: CCIPObjectRef,
        clock: sui::clock::Clock,
    }

    fun setup(): (Env, NonceManagerCap, dd::SourceTransferCap) {
        let mut scenario = ts::begin(OWNER);
        let ctx = scenario.ctx();
        let mut clock = sui::clock::create_for_testing(ctx);
        clock.set_for_testing(1_000_000_000);

        mcms_account::test_init(ctx);
        mcms_registry::test_init(ctx);
        mcms_deployer::test_init(ctx);

        state_object::test_init(ctx);
        dd::test_init(ctx);

        onramp::test_init(ctx);

        scenario.next_tx(OWNER);

        let registry = ts::take_shared<Registry>(&scenario);
        let mut ref = ts::take_shared<CCIPObjectRef>(&scenario);
        let state = ts::take_shared<OnRampState>(&scenario);

        let ccip_owner_cap = ts::take_from_sender<state_object::OwnerCap>(&scenario);
        nonce_manager::initialize(&mut ref, &ccip_owner_cap, scenario.ctx());

        scenario.next_tx(OWNER);
        
        let source_transfer_cap = ts::take_from_sender<dd::SourceTransferCap>(&scenario);
        let nonce_manager_cap = ts::take_from_sender<NonceManagerCap>(&scenario);

        ts::return_to_address(OWNER, ccip_owner_cap);

        let env = Env {
            scenario,
            state,
            registry,
            ref,
            clock,
        };

        (env, nonce_manager_cap, source_transfer_cap)
    }

    fun tear_down(env: Env) {
        let Env { scenario, state, registry, ref, clock } = env;

        ts::return_shared(state);
        ts::return_shared(registry);
        ts::return_shared(ref);

        clock.destroy_for_testing();

        ts::end(scenario);
    }

    fun initialize_onramp(
        env: &mut Env,
        onramp_owner_cap: &OwnerCap,
        nonce_manager_cap: NonceManagerCap,
        source_transfer_cap: dd::SourceTransferCap,
    ) {
        let ctx = env.scenario.ctx();
        onramp::initialize(
            &mut env.state,
            onramp_owner_cap,
            nonce_manager_cap,
            source_transfer_cap,
            123, // chain_selector
            ctx.sender(),
            ctx.sender(),
            vector[DEST_CHAIN_SELECTOR_1, DEST_CHAIN_SELECTOR_2], // dest_chain_selectors
            vector[true, false], // dest_chain_enabled
            vector[true, false], // dest_chain_allowlist_enabled
            ctx
        );
    }

    #[test]
    public fun test_initialize() {
        let (mut env, nonce_manager_cap, source_transfer_cap) = setup();
        let owner_cap = ts::take_from_sender<OwnerCap>(&env.scenario);
        initialize_onramp(&mut env, &owner_cap, nonce_manager_cap, source_transfer_cap);

        assert!(onramp::is_chain_supported(&env.state, DEST_CHAIN_SELECTOR_1));
        assert!(onramp::is_chain_supported(&env.state, DEST_CHAIN_SELECTOR_2));

        assert!(onramp::get_expected_next_sequence_number(&env.state, DEST_CHAIN_SELECTOR_1) == 1);
        assert!(onramp::get_expected_next_sequence_number(&env.state, DEST_CHAIN_SELECTOR_2) == 1);

        let (enabled, seq, allowlist_enabled, allowed_senders) = onramp::get_dest_chain_config(&env.state, DEST_CHAIN_SELECTOR_1);
        assert!(enabled == true);
        assert!(seq == 0);
        assert!(allowlist_enabled == true);
        assert!(allowed_senders == vector[]);

        let (enabled, seq, allowlist_enabled, allowed_senders) = onramp::get_dest_chain_config(&env.state, DEST_CHAIN_SELECTOR_2);
        assert!(enabled == false);
        assert!(seq == 0);
        assert!(allowlist_enabled == false);
        assert!(allowed_senders == vector[]);

        env.tear_down();

        ts::return_to_address(OWNER, owner_cap);
    }

    #[test]
    public fun test_apply_allowlist_updates() {
        let (mut env, nonce_manager_cap, source_transfer_cap) = setup();
        let owner_cap = ts::take_from_sender<OwnerCap>(&env.scenario);
        initialize_onramp(&mut env, &owner_cap, nonce_manager_cap, source_transfer_cap);

        onramp::apply_allowlist_updates(
            &mut env.state,
            &owner_cap,
            vector[DEST_CHAIN_SELECTOR_1, DEST_CHAIN_SELECTOR_2], // dest_chain_selectors
            vector[true, true], // dest_chain_allowlist_enabled
            vector[
                vector[ALLOWED_SENDER_1, ALLOWED_SENDER_2],
                vector[ALLOWED_SENDER_3]
            ], // dest_chain_add_allowed_senders
            vector[
                vector[],
                vector[]
            ], // dest_chain_remove_allowed_senders
            env.scenario.ctx()
        );

        let (allowlist_enabled, allowed_senders) = onramp::get_allowed_senders_list(&env.state, DEST_CHAIN_SELECTOR_1);
        assert!(allowlist_enabled == true);
        assert!(allowed_senders == vector[ALLOWED_SENDER_1, ALLOWED_SENDER_2]);
        let (allowlist_enabled, allowed_senders) = onramp::get_allowed_senders_list(&env.state, DEST_CHAIN_SELECTOR_2);
        assert!(allowlist_enabled == true);
        assert!(allowed_senders == vector[ALLOWED_SENDER_3]);

        onramp::apply_allowlist_updates(
            &mut env.state,
            &owner_cap,
            vector[DEST_CHAIN_SELECTOR_1, DEST_CHAIN_SELECTOR_2], // dest_chain_selectors
            vector[true, false], // dest_chain_allowlist_enabled
            vector[
                vector[],
                vector[]
            ], // dest_chain_add_allowed_senders
            vector[
                vector[ALLOWED_SENDER_2],
                vector[]
            ], // dest_chain_remove_allowed_senders
            env.scenario.ctx()
        );

        let (allowlist_enabled, allowed_senders) = onramp::get_allowed_senders_list(&env.state, DEST_CHAIN_SELECTOR_1);
        assert!(allowlist_enabled == true);
        assert!(allowed_senders == vector[ALLOWED_SENDER_1]);
        let (allowlist_enabled, allowed_senders) = onramp::get_allowed_senders_list(&env.state, DEST_CHAIN_SELECTOR_2);
        assert!(allowlist_enabled == false);
        assert!(allowed_senders == vector[ALLOWED_SENDER_3]);

        env.tear_down();

        ts::return_to_address(OWNER, owner_cap);
    }
}
