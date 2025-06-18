#[test_only]
module ccip_onramp::onramp_mcms_test {
    use sui::test_scenario::{Self as ts, Scenario};
    use sui::package;
    use std::string;
    use sui::bcs;

    use ccip_onramp::onramp::{Self, OnRampState};
    use ccip_onramp::ownable::{OwnerCap};
    use ccip::state_object::{Self, CCIPObjectRef};
    use ccip::dynamic_dispatcher as dd;
    use ccip::nonce_manager::{Self, NonceManagerCap};

    use mcms::mcms_registry::{Self, Registry};
    use mcms::mcms_account::{Self};
    use mcms::mcms_deployer::{Self, DeployerState};

    const DEST_CHAIN_SELECTOR_1: u64 = 1;
    const DEST_CHAIN_SELECTOR_2: u64 = 2;
    const ALLOWED_SENDER_1: address = @0x11;
    const ALLOWED_SENDER_2: address = @0x22;
    const ALLOWED_SENDER_3: address = @0x33;
    const OWNER: address = @0x123;
    const SENDER: address = @0x456;

    const MODULE_NAME: vector<u8> = b"onramp";

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

        let state_object_owner_cap = ts::take_from_sender<state_object::OwnerCap>(&scenario);
        nonce_manager::initialize(&mut ref, &state_object_owner_cap, scenario.ctx());
        ts::return_to_address(OWNER, state_object_owner_cap);

        scenario.next_tx(OWNER);
        
        let source_transfer_cap = ts::take_from_sender<dd::SourceTransferCap>(&scenario);
        let nonce_manager_cap = ts::take_from_sender<NonceManagerCap>(&scenario);

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

    /// Set up onramp with MCMS
    /// Sends owner_cap to the registry
    fun mcms_register_entrypoint(env: &mut Env, onramp_owner_cap: OwnerCap) {
        let ctx = env.scenario.ctx();

        // Initialize the user data with mcms_registry
        // This creates a owner_cap and mcms_registry owns this cap
        onramp::mcms_register_entrypoint(
            &mut env.registry,
            &mut env.state,
            onramp_owner_cap,
            @mcms,
            ctx,
        );

        env.scenario.next_tx(OWNER);
    }

    fun mcms_register_upgrade_cap(env: &mut Env) {
        let mut registry = ts::take_shared<Registry>(&env.scenario);
        let mut deployer_state = ts::take_shared<DeployerState>(&env.scenario);

        let upgrade_cap = package::test_publish(@ccip_onramp.to_id(), env.scenario.ctx());

        // Initialize the user data with mcms_registry
        // This creates a owner_cap and mcms_registry owns this cap
        onramp::mcms_register_upgrade_cap(upgrade_cap, &mut registry, &mut deployer_state, env.scenario.ctx());

        ts::return_shared(registry);
        ts::return_shared(deployer_state);

        env.scenario.next_tx(OWNER);
    }

    #[test]
    public fun test_mcms_set_dynamic_config() {
        let (mut env, nonce_manager_cap, source_transfer_cap) = setup();

        let owner_cap = ts::take_from_sender<OwnerCap>(&env.scenario);
        initialize_onramp(&mut env, &owner_cap, nonce_manager_cap, source_transfer_cap);

        // Initialize owner_cap with MCMS
        mcms_register_entrypoint(&mut env, owner_cap);

        let mut data = vector::empty<u8>();
        data.append(bcs::to_bytes(&@0x123));
        data.append(bcs::to_bytes(&@0x456));

        let params = mcms_registry::test_create_executing_callback_params(
            @ccip_onramp,
            string::utf8(MODULE_NAME),
            string::utf8(b"set_dynamic_config"),
            data
        );

        onramp::mcms_entrypoint(
            &mut env.state,
            &mut env.registry,
            params,
            env.scenario.ctx()
        );

        let (fee_aggregator, allowlist_admin) = onramp::get_dynamic_config_fields(onramp::get_dynamic_config(&env.state));
        assert!(fee_aggregator == @0x123);
        assert!(allowlist_admin == @0x456);

        env.tear_down();
    }

    #[test]
    public fun test_mcms_apply_dest_chain_config_updates() {
        let (mut env, nonce_manager_cap, source_transfer_cap) = setup();

        let owner_cap = ts::take_from_sender<OwnerCap>(&env.scenario);
        initialize_onramp(&mut env, &owner_cap, nonce_manager_cap, source_transfer_cap);

        // Initialize owner_cap with MCMS
        mcms_register_entrypoint(&mut env, owner_cap);

        // Prepare data: dest_chain_selectors, dest_chain_enabled, dest_chain_allowlist_enabled
        let mut data = vector::empty<u8>();
        data.append(bcs::to_bytes(&vector[DEST_CHAIN_SELECTOR_1, DEST_CHAIN_SELECTOR_2])); // dest_chain_selectors
        data.append(bcs::to_bytes(&vector[true, false])); // dest_chain_enabled
        data.append(bcs::to_bytes(&vector[false, true])); // dest_chain_allowlist_enabled

        let params = mcms_registry::test_create_executing_callback_params(
            @ccip_onramp,
            string::utf8(MODULE_NAME),
            string::utf8(b"apply_dest_chain_config_updates"),
            data
        );

        onramp::mcms_entrypoint(
            &mut env.state,
            &mut env.registry,
            params,
            env.scenario.ctx()
        );

        let (is_enabled, _sequence_number, allowlist_enabled, _allowed_senders) = onramp::get_dest_chain_config(&env.state, DEST_CHAIN_SELECTOR_1);
        assert!(is_enabled == true);
        assert!(allowlist_enabled == false);

        let (is_enabled, _sequence_number, allowlist_enabled, _allowed_senders) = onramp::get_dest_chain_config(&env.state, DEST_CHAIN_SELECTOR_2);
        assert!(is_enabled == false);
        assert!(allowlist_enabled == true);

        env.tear_down();
    }

    #[test]
    public fun test_mcms_apply_allowlist_updates() {
        let (mut env, nonce_manager_cap, source_transfer_cap) = setup();

        let owner_cap = ts::take_from_sender<OwnerCap>(&env.scenario);
        initialize_onramp(&mut env, &owner_cap, nonce_manager_cap, source_transfer_cap);

        mcms_register_entrypoint(&mut env, owner_cap);

        let mut data = vector::empty<u8>();
        data.append(bcs::to_bytes(&vector[DEST_CHAIN_SELECTOR_1, DEST_CHAIN_SELECTOR_2])); // dest_chain_selectors
        data.append(bcs::to_bytes(&vector[true, true])); // dest_chain_allowlist_enabled
        data.append(bcs::to_bytes(&vector[
            vector[ALLOWED_SENDER_1, ALLOWED_SENDER_2],
            vector[ALLOWED_SENDER_3]
        ])); // dest_chain_add_allowed_senders
        data.append(bcs::to_bytes(&vector[
            vector<address>[],
            vector<address>[]
        ])); // dest_chain_remove_allowed_senders

        let params = mcms_registry::test_create_executing_callback_params(
            @ccip_onramp,
            string::utf8(MODULE_NAME),
            string::utf8(b"apply_allowlist_updates"),
            data
        );

        onramp::mcms_entrypoint(
            &mut env.state,
            &mut env.registry,
            params,
            env.scenario.ctx()
        );

        let (allowlist_enabled, allowed_senders) = onramp::get_allowed_senders_list(&env.state, DEST_CHAIN_SELECTOR_1);
        assert!(allowlist_enabled == true);
        assert!(allowed_senders == vector[ALLOWED_SENDER_1, ALLOWED_SENDER_2]);

        let (allowlist_enabled, allowed_senders) = onramp::get_allowed_senders_list(&env.state, DEST_CHAIN_SELECTOR_2);
        assert!(allowlist_enabled == true);
        assert!(allowed_senders == vector[ALLOWED_SENDER_3]);

        env.tear_down();
    }

    #[test]
    public fun test_mcms_transfer_ownership_e2e() {
        let (mut env, nonce_manager_cap, source_transfer_cap) = setup();

        let owner_cap = ts::take_from_sender<OwnerCap>(&env.scenario);
        initialize_onramp(&mut env, &owner_cap, nonce_manager_cap, source_transfer_cap);

        // Transfer ownership to MCMS first
        mcms_register_entrypoint(&mut env, owner_cap);

        let new_owner = @0x999;
        let mut data = vector::empty<u8>();
        data.append(bcs::to_bytes(&new_owner));

        // Transfer ownership to `new_owner` via MCMS
        let params = mcms_registry::test_create_executing_callback_params(
            @ccip_onramp,
            string::utf8(MODULE_NAME),
            string::utf8(b"transfer_ownership"),
            data
        );

        onramp::mcms_entrypoint(
            &mut env.state,
            &mut env.registry,
            params,
            env.scenario.ctx()
        );

        // Accept ownership as `new_owner`
        env.scenario.next_tx(new_owner);
        onramp::accept_ownership(&mut env.state, env.scenario.ctx());

        // Execute ownership transfer as MCMS to `new_owner`
        let mut data = vector::empty<u8>();
        data.append(bcs::to_bytes(&new_owner));

        let params = mcms_registry::test_create_executing_callback_params(
            @ccip_onramp,
            string::utf8(MODULE_NAME),
            string::utf8(b"execute_ownership_transfer"),
            data
        );

        onramp::mcms_entrypoint(
            &mut env.state,
            &mut env.registry,
            params,
            env.scenario.ctx()
        );

        // Verify ownership transfer was completed
        assert!(onramp::owner(&env.state) == new_owner);

        assert!(!onramp::has_pending_transfer(&env.state));
        assert!(onramp::pending_transfer_from(&env.state) == option::none());
        assert!(onramp::pending_transfer_to(&env.state) == option::none());
        assert!(onramp::pending_transfer_accepted(&env.state) == option::none());

        env.tear_down();
    }
}
