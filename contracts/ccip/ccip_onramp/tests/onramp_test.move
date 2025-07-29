#[test_only]
module ccip_onramp::onramp_test {
    use ccip_onramp::onramp::{Self, OnRampState};
    use ccip_onramp::ownable::{OwnerCap};
    use ccip::state_object::{Self, CCIPObjectRef};
    use sui::test_scenario::{Self as ts, Scenario};
    use ccip::dynamic_dispatcher as dd;
    use ccip::nonce_manager::{Self, NonceManagerCap};
    use ccip::token_admin_registry;
    use ccip::rmn_remote;
    use ccip::fee_quoter;
    use sui::coin::{Self, CoinMetadata};
    use std::string;

    use mcms::mcms_registry::{Self, Registry};
    use mcms::mcms_account;
    use mcms::mcms_deployer;

    const DEST_CHAIN_SELECTOR_1: u64 = 1;
    const DEST_CHAIN_SELECTOR_2: u64 = 2;
    const ALLOWED_SENDER_1: address = @0x11;
    const ALLOWED_SENDER_2: address = @0x22;
    const ALLOWED_SENDER_3: address = @0x33;
    const OWNER: address = @0x123;
    const FEE_AGGREGATOR: address = @0x456;
    const ALLOWLIST_ADMIN: address = @0x789;

    // Chain family selectors
    const CHAIN_FAMILY_SELECTOR_EVM: vector<u8> = x"2812d52c";
    const CHAIN_FAMILY_SELECTOR_SVM: vector<u8> = x"1e10bdc4";

    // Test addresses
    const EVM_RECEIVER_ADDRESS: vector<u8> = x"000000000000000000000000f4030086522a5beea4988f8ca5b36dbc97bee88c";
    const SVM_RECEIVER_ADDRESS: vector<u8> = x"1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef";
    const SVM_TOKEN_RECEIVER: vector<u8> = x"1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef";

    // Fee quoter configuration constants
    const DEFAULT_MAX_FEE_JUELS: u256 = 1000000;
    const DEFAULT_TOKEN_PRICE_STALENESS_THRESHOLD: u64 = 3600;
    const DEFAULT_TOKEN_PRICE: u256 = 150_000_000_000_000_000_000; // $150 * 1e18
    const DEFAULT_GAS_PRICE: u256 = 1_000_000_000_000; // 1e12
    const DEFAULT_PREMIUM_MULTIPLIER: u64 = 1_000_000_000_000_000_000; // 1e18 = 100%

    // Test error constants
    const EFeeCalculationFailed: u64 = 1;

    // Test token for fee testing - using proper one-time witness pattern
    public struct ONRAMP_TEST has drop {}

    // Helper function to create test token with proper one-time witness
    #[test_only]
    public fun create_test_token(ctx: &mut TxContext): (coin::TreasuryCap<ONRAMP_TEST>, CoinMetadata<ONRAMP_TEST>) {
        coin::create_currency(
            ONRAMP_TEST {},
            8, // decimals
            b"TEST",
            b"Test Token",
            b"Test token for testing",
            option::none(),
            ctx
        )
    }

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
        token_admin_registry::initialize(&mut ref, &ccip_owner_cap, scenario.ctx());
        rmn_remote::initialize(&mut ref, &ccip_owner_cap, 1000, scenario.ctx());
        fee_quoter::initialize(
            &mut ref, 
            &ccip_owner_cap, 
            1000000, // max_fee_juels_per_msg
            @0x1, // link_token address
            3600, // token_price_staleness_threshold (1 hour)
            vector[], // fee_tokens
            scenario.ctx()
        );

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
            FEE_AGGREGATOR,
            ALLOWLIST_ADMIN,
            vector[DEST_CHAIN_SELECTOR_1, DEST_CHAIN_SELECTOR_2], // dest_chain_selectors
            vector[true, false], // dest_chain_enabled
            vector[true, false], // dest_chain_allowlist_enabled
            ctx
        );
    }

    // === Fee Quoter Setup Helpers ===
    
    /// Configure fee quoter for a destination chain with the specified chain family
    fun configure_fee_quoter_dest_chain(
        ref: &mut CCIPObjectRef,
        ccip_owner_cap: &state_object::OwnerCap,
        dest_chain_selector: u64,
        chain_family_selector: vector<u8>
    ) {
        fee_quoter::apply_dest_chain_config_updates(
            ref,
            ccip_owner_cap,
            dest_chain_selector,
            true, // is_enabled
            1000, // max_number_of_tokens_per_msg
            30_000, // max_data_bytes
            30_000_000, // max_per_msg_gas_limit
            250_000, // dest_gas_overhead
            16, // dest_gas_per_payload_byte_base
            200, // dest_gas_per_payload_byte_high
            300, // dest_gas_per_payload_byte_threshold
            400000, // dest_data_availability_overhead_gas
            500, // dest_gas_per_data_availability_byte
            600, // dest_data_availability_multiplier_bps
            chain_family_selector,
            false, // enforce_out_of_order
            50, // default_token_fee_usd_cents
            90_000, // default_token_dest_gas_overhead
            200_000, // default_tx_gas_limit
            DEFAULT_PREMIUM_MULTIPLIER, // gas_multiplier_wei_per_eth
            1_000_000, // gas_price_staleness_threshold
            50 // network_fee_usd_cents
        );
    }

    /// Setup token and price configuration for fee calculation
    fun setup_fee_token_and_prices(
        ref: &mut CCIPObjectRef,
        ccip_owner_cap: &state_object::OwnerCap,
        clock: &sui::clock::Clock,
        coin_metadata_addr: address,
        dest_chain_selector: u64,
        ctx: &mut TxContext
    ): fee_quoter::FeeQuoterCap {
        // Add the test token as a fee token
        fee_quoter::apply_fee_token_updates(
            ref,
            ccip_owner_cap,
            vector[], // fee_tokens_to_remove
            vector[coin_metadata_addr], // fee_tokens_to_add
        );

        // Set up premium multipliers
        fee_quoter::apply_premium_multiplier_wei_per_eth_updates(
            ref,
            ccip_owner_cap,
            vector[coin_metadata_addr],
            vector[DEFAULT_PREMIUM_MULTIPLIER],
        );

        // Set up price updates
        let fee_quoter_cap = fee_quoter::create_fee_quoter_cap(ctx);
        fee_quoter::update_prices(
            ref,
            &fee_quoter_cap,
            clock,
            vector[coin_metadata_addr],
            vector[DEFAULT_TOKEN_PRICE],
            vector[dest_chain_selector],
            vector[DEFAULT_GAS_PRICE]
        );

        fee_quoter_cap
    }

    /// Create a complete standalone test environment for fee testing
    fun setup_standalone_fee_test_env(): (
        ts::Scenario,
        CCIPObjectRef,
        OnRampState,
        Registry,
        sui::clock::Clock,
        state_object::OwnerCap,
        OwnerCap,
        coin::TreasuryCap<ONRAMP_TEST>,
        CoinMetadata<ONRAMP_TEST>,
        address
    ) {
        let mut scenario = ts::begin(OWNER);
        let ctx = scenario.ctx();
        let mut clock = sui::clock::create_for_testing(ctx);
        clock.set_for_testing(1_000_000_000);

        // Initialize all required modules
        mcms_account::test_init(ctx);
        mcms_registry::test_init(ctx);
        mcms_deployer::test_init(ctx);
        state_object::test_init(ctx);
        dd::test_init(ctx);
        onramp::test_init(ctx);

        scenario.next_tx(OWNER);

        let mut ref = ts::take_shared<CCIPObjectRef>(&scenario);
        let registry = ts::take_shared<Registry>(&scenario);
        let mut state = ts::take_shared<OnRampState>(&scenario);
        
        let ccip_owner_cap = ts::take_from_sender<state_object::OwnerCap>(&scenario);
        let onramp_owner_cap = ts::take_from_sender<OwnerCap>(&scenario);
        
        // Initialize CCIP modules
        nonce_manager::initialize(&mut ref, &ccip_owner_cap, scenario.ctx());
        token_admin_registry::initialize(&mut ref, &ccip_owner_cap, scenario.ctx());
        rmn_remote::initialize(&mut ref, &ccip_owner_cap, 1000, scenario.ctx());
        fee_quoter::initialize(
            &mut ref, 
            &ccip_owner_cap, 
            DEFAULT_MAX_FEE_JUELS,
            @0x1, // link_token address
            DEFAULT_TOKEN_PRICE_STALENESS_THRESHOLD,
            vector[], // fee_tokens
            scenario.ctx()
        );

        scenario.next_tx(OWNER);
        
        let source_transfer_cap = ts::take_from_sender<dd::SourceTransferCap>(&scenario);
        let nonce_manager_cap = ts::take_from_sender<NonceManagerCap>(&scenario);

        // Initialize onramp
        onramp::initialize(
            &mut state,
            &onramp_owner_cap,
            nonce_manager_cap,
            source_transfer_cap,
            123, // chain_selector
            FEE_AGGREGATOR,
            ALLOWLIST_ADMIN,
            vector[DEST_CHAIN_SELECTOR_1],
            vector[true], // dest_chain_enabled
            vector[false], // dest_chain_allowlist_enabled
            scenario.ctx()
        );

        // Create test token
        scenario.next_tx(OWNER);
        let ctx = scenario.ctx();
        let (treasury_cap, coin_metadata) = create_test_token(ctx);
        let coin_metadata_addr = object::id_to_address(object::borrow_id(&coin_metadata));

        (scenario, ref, state, registry, clock, ccip_owner_cap, onramp_owner_cap, treasury_cap, coin_metadata, coin_metadata_addr)
    }

    /// Cleanup standalone test environment
    fun cleanup_standalone_fee_test_env(
        scenario: ts::Scenario,
        ref: CCIPObjectRef,
        state: OnRampState,
        registry: Registry,
        clock: sui::clock::Clock,
        ccip_owner_cap: state_object::OwnerCap,
        onramp_owner_cap: OwnerCap,
        treasury_cap: coin::TreasuryCap<ONRAMP_TEST>,
        coin_metadata: CoinMetadata<ONRAMP_TEST>,
        fee_quoter_cap: fee_quoter::FeeQuoterCap
    ) {
        fee_quoter::destroy_fee_quoter_cap(fee_quoter_cap);
        transfer::public_transfer(treasury_cap, OWNER);
        transfer::public_freeze_object(coin_metadata);
        
        ts::return_to_address(OWNER, ccip_owner_cap);
        ts::return_to_address(OWNER, onramp_owner_cap);
        ts::return_shared(state);
        ts::return_shared(registry);
        ts::return_shared(ref);
        clock.destroy_for_testing();
        ts::end(scenario);
    }

    // === Basic Tests ===

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

    #[test]
    public fun test_set_dynamic_config() {
        let (mut env, nonce_manager_cap, source_transfer_cap) = setup();
        let owner_cap = ts::take_from_sender<OwnerCap>(&env.scenario);
        initialize_onramp(&mut env, &owner_cap, nonce_manager_cap, source_transfer_cap);

        // Test initial config
        let dynamic_config = onramp::get_dynamic_config(&env.state);
        let (fee_aggregator, allowlist_admin) = onramp::get_dynamic_config_fields(dynamic_config);
        assert!(fee_aggregator == FEE_AGGREGATOR);
        assert!(allowlist_admin == ALLOWLIST_ADMIN);

        // Update config
        let new_fee_aggregator = @0x999;
        let new_allowlist_admin = @0x888;
        onramp::set_dynamic_config(&mut env.state, &owner_cap, new_fee_aggregator, new_allowlist_admin);

        // Verify config was updated
        let updated_config = onramp::get_dynamic_config(&env.state);
        let (updated_fee_aggregator, updated_allowlist_admin) = onramp::get_dynamic_config_fields(updated_config);
        assert!(updated_fee_aggregator == new_fee_aggregator);
        assert!(updated_allowlist_admin == new_allowlist_admin);

        env.tear_down();
        ts::return_to_address(OWNER, owner_cap);
    }

    #[test]
    public fun test_apply_dest_chain_config_updates() {
        let (mut env, nonce_manager_cap, source_transfer_cap) = setup();
        let owner_cap = ts::take_from_sender<OwnerCap>(&env.scenario);
        initialize_onramp(&mut env, &owner_cap, nonce_manager_cap, source_transfer_cap);

        // Add a new destination chain
        let new_chain_selector = 999;
        onramp::apply_dest_chain_config_updates(
            &mut env.state,
            &owner_cap,
            vector[new_chain_selector],
            vector[true], // enabled
            vector[false] // allowlist disabled
        );

        // Verify new chain was added
        assert!(onramp::is_chain_supported(&env.state, new_chain_selector));
        let (enabled, seq, allowlist_enabled, allowed_senders) = onramp::get_dest_chain_config(&env.state, new_chain_selector);
        assert!(enabled == true);
        assert!(seq == 0);
        assert!(allowlist_enabled == false);
        assert!(allowed_senders == vector[]);

        // Update existing chain config
        onramp::apply_dest_chain_config_updates(
            &mut env.state,
            &owner_cap,
            vector[DEST_CHAIN_SELECTOR_2],
            vector[true], // enable previously disabled chain
            vector[true] // enable allowlist
        );

        // Verify chain was updated
        let (enabled, seq, allowlist_enabled, allowed_senders) = onramp::get_dest_chain_config(&env.state, DEST_CHAIN_SELECTOR_2);
        assert!(enabled == true);
        assert!(seq == 0);
        assert!(allowlist_enabled == true);
        assert!(allowed_senders == vector[]);

        env.tear_down();
        ts::return_to_address(OWNER, owner_cap);
    }

    #[test]
    public fun test_get_static_config() {
        let (mut env, nonce_manager_cap, source_transfer_cap) = setup();
        let owner_cap = ts::take_from_sender<OwnerCap>(&env.scenario);
        initialize_onramp(&mut env, &owner_cap, nonce_manager_cap, source_transfer_cap);

        // Test static config getters
        let static_config = onramp::get_static_config(&env.state);
        let chain_selector = onramp::get_static_config_fields(static_config);
        assert!(chain_selector == 123); // From initialize call

        env.tear_down();
        ts::return_to_address(OWNER, owner_cap);
    }

    #[test]
    public fun test_get_outbound_nonce() {
        let (mut env, nonce_manager_cap, source_transfer_cap) = setup();
        let owner_cap = ts::take_from_sender<OwnerCap>(&env.scenario);
        initialize_onramp(&mut env, &owner_cap, nonce_manager_cap, source_transfer_cap);

        // Test get_outbound_nonce
        let nonce = onramp::get_outbound_nonce(&env.ref, DEST_CHAIN_SELECTOR_1, OWNER);
        assert!(nonce == 0); // Initial nonce should be 0

        env.tear_down();
        ts::return_to_address(OWNER, owner_cap);
    }

    #[test]
    public fun test_ownership_functions() {
        let (mut env, nonce_manager_cap, source_transfer_cap) = setup();
        let owner_cap = ts::take_from_sender<OwnerCap>(&env.scenario);
        initialize_onramp(&mut env, &owner_cap, nonce_manager_cap, source_transfer_cap);

        // Test ownership getters
        assert!(onramp::owner(&env.state) == OWNER);
        assert!(!onramp::has_pending_transfer(&env.state));
        assert!(onramp::pending_transfer_from(&env.state).is_none());
        assert!(onramp::pending_transfer_to(&env.state).is_none());
        assert!(onramp::pending_transfer_accepted(&env.state).is_none());

        // Test transfer ownership
        let new_owner = @0x999;
        onramp::transfer_ownership(&mut env.state, &owner_cap, new_owner, env.scenario.ctx());

        // Verify pending transfer
        assert!(onramp::has_pending_transfer(&env.state));
        assert!(onramp::pending_transfer_from(&env.state).extract() == OWNER);
        assert!(onramp::pending_transfer_to(&env.state).extract() == new_owner);
        assert!(onramp::pending_transfer_accepted(&env.state).extract() == false);

        // Test accept ownership
        env.scenario.next_tx(new_owner);
        onramp::accept_ownership(&mut env.state, env.scenario.ctx());

        // Verify acceptance
        assert!(onramp::pending_transfer_accepted(&env.state).extract() == true);

        env.tear_down();
        ts::return_to_address(OWNER, owner_cap);
    }

    #[test]
    public fun test_apply_allowlist_updates_by_admin() {
        let (mut env, nonce_manager_cap, source_transfer_cap) = setup();
        let owner_cap = ts::take_from_sender<OwnerCap>(&env.scenario);
        initialize_onramp(&mut env, &owner_cap, nonce_manager_cap, source_transfer_cap);

        // Test allowlist updates by admin (not owner)
        env.scenario.next_tx(ALLOWLIST_ADMIN);
        onramp::apply_allowlist_updates_by_admin(
            &mut env.state,
            vector[DEST_CHAIN_SELECTOR_1],
            vector[true],
            vector[vector[ALLOWED_SENDER_1]],
            vector[vector[]],
            env.scenario.ctx()
        );

        // Verify update was applied
        let (allowlist_enabled, allowed_senders) = onramp::get_allowed_senders_list(&env.state, DEST_CHAIN_SELECTOR_1);
        assert!(allowlist_enabled == true);
        assert!(allowed_senders == vector[ALLOWED_SENDER_1]);

        env.tear_down();
        ts::return_to_address(OWNER, owner_cap);
    }

    #[test]
    #[expected_failure(abort_code = onramp::EOnlyCallableByAllowlistAdmin)]
    public fun test_apply_allowlist_updates_by_admin_unauthorized() {
        let (mut env, nonce_manager_cap, source_transfer_cap) = setup();
        let owner_cap = ts::take_from_sender<OwnerCap>(&env.scenario);
        initialize_onramp(&mut env, &owner_cap, nonce_manager_cap, source_transfer_cap);

        // Test allowlist updates by unauthorized user
        env.scenario.next_tx(@0x999); // Not the allowlist admin
        onramp::apply_allowlist_updates_by_admin(
            &mut env.state,
            vector[DEST_CHAIN_SELECTOR_1],
            vector[true],
            vector[vector[ALLOWED_SENDER_1]],
            vector[vector[]],
            env.scenario.ctx()
        );

        env.tear_down();
        ts::return_to_address(OWNER, owner_cap);
    }

    // === Function-specific Tests ===

    #[test]
    public fun test_type_and_version() {
        let version = onramp::type_and_version();
        assert!(version == string::utf8(b"OnRamp 1.6.0"));
    }

    #[test]
    public fun test_get_ccip_package_id() {
        let package_id = onramp::get_ccip_package_id();
        assert!(package_id == @ccip);
    }

    // === Fee Tests ===

    #[test]
    #[expected_failure(abort_code = ccip::fee_quoter::EUnknownDestChainSelector)]
    public fun test_get_fee_unconfigured_destination() {
        let (mut env, nonce_manager_cap, source_transfer_cap) = setup();
        let owner_cap = ts::take_from_sender<OwnerCap>(&env.scenario);
        initialize_onramp(&mut env, &owner_cap, nonce_manager_cap, source_transfer_cap);

        // Create a test token for fee calculation
        env.scenario.next_tx(OWNER);
        let ctx = env.scenario.ctx();
        let (treasury_cap, coin_metadata) = create_test_token(ctx);

        // Test get_fee function with basic parameters
        // This should fail because destination chain is not configured in fee quoter
        let _fee = onramp::get_fee(
            &env.ref,
            &env.clock,
            DEST_CHAIN_SELECTOR_1,
            EVM_RECEIVER_ADDRESS,
            b"test_data", // data
            vector[], // token_addresses (empty for no tokens)
            vector[], // token_amounts (empty for no tokens)
            &coin_metadata,
            b"extra_args"
        );

        // Clean up test objects (won't be reached due to expected failure)
        transfer::public_transfer(treasury_cap, OWNER);
        transfer::public_freeze_object(coin_metadata);
        env.tear_down();
        ts::return_to_address(OWNER, owner_cap);
    }

    #[test]
    public fun test_get_fee_success() {
        let (
            mut scenario, mut ref, state, registry, clock, ccip_owner_cap, onramp_owner_cap, 
            treasury_cap, coin_metadata, coin_metadata_addr
        ) = setup_standalone_fee_test_env();

        // Configure fee quoter for EVM
        configure_fee_quoter_dest_chain(&mut ref, &ccip_owner_cap, DEST_CHAIN_SELECTOR_1, CHAIN_FAMILY_SELECTOR_EVM);
        
        // Setup token and prices
        let fee_quoter_cap = setup_fee_token_and_prices(
            &mut ref, &ccip_owner_cap, &clock, coin_metadata_addr, DEST_CHAIN_SELECTOR_1, scenario.ctx()
        );

        // Test get_fee function - this should succeed
        let fee = onramp::get_fee(
            &ref,
            &clock,
            DEST_CHAIN_SELECTOR_1,
            EVM_RECEIVER_ADDRESS,
            b"test_data",
            vector[], // token_addresses (empty for no tokens)
            vector[], // token_amounts (empty for no tokens)
            &coin_metadata,
            b"" // empty extra_args will use defaults
        );

        // Verify that we got a non-zero fee
        assert!(fee > 0, EFeeCalculationFailed);

        cleanup_standalone_fee_test_env(
            scenario, ref, state, registry, clock, ccip_owner_cap, onramp_owner_cap, 
            treasury_cap, coin_metadata, fee_quoter_cap
        );
    }

    #[test]
    public fun test_get_fee_success_svm() {
        let (
            mut scenario, mut ref, state, registry, clock, ccip_owner_cap, onramp_owner_cap, 
            treasury_cap, coin_metadata, coin_metadata_addr
        ) = setup_standalone_fee_test_env();

        // Configure fee quoter for SVM
        configure_fee_quoter_dest_chain(&mut ref, &ccip_owner_cap, DEST_CHAIN_SELECTOR_1, CHAIN_FAMILY_SELECTOR_SVM);
        
        // Setup token and prices
        let fee_quoter_cap = setup_fee_token_and_prices(
            &mut ref, &ccip_owner_cap, &clock, coin_metadata_addr, DEST_CHAIN_SELECTOR_1, scenario.ctx()
        );

        // Create SVM-specific extra args
        let svm_extra_args = {
            use ccip::client;
            client::encode_svm_extra_args_v1(
                100000, // compute_units
                0, // account_is_writable_bitmap
                false, // allow_out_of_order_execution
                SVM_TOKEN_RECEIVER,
                vector[] // accounts
            )
        };

        // Test get_fee function with SVM configuration - this should succeed
        let fee = onramp::get_fee(
            &ref,
            &clock,
            DEST_CHAIN_SELECTOR_1,
            SVM_RECEIVER_ADDRESS,
            b"test_data",
            vector[], // token_addresses (empty for no tokens)
            vector[], // token_amounts (empty for no tokens)
            &coin_metadata,
            svm_extra_args
        );

        // Verify that we got a non-zero fee
        assert!(fee > 0, EFeeCalculationFailed);

        cleanup_standalone_fee_test_env(
            scenario, ref, state, registry, clock, ccip_owner_cap, onramp_owner_cap, 
            treasury_cap, coin_metadata, fee_quoter_cap
        );
    }

    // === Error Code Tests ===

    #[test]
    #[expected_failure(abort_code = onramp::EDestChainArgumentMismatch)]
    public fun test_error_dest_chain_argument_mismatch() {
        let (mut env, nonce_manager_cap, source_transfer_cap) = setup();
        let owner_cap = ts::take_from_sender<OwnerCap>(&env.scenario);

        onramp::initialize(
            &mut env.state,
            &owner_cap,
            nonce_manager_cap,
            source_transfer_cap,
            123, // chain_selector
            FEE_AGGREGATOR,
            ALLOWLIST_ADMIN,
            vector[DEST_CHAIN_SELECTOR_1, DEST_CHAIN_SELECTOR_2], // 2 elements
            vector[true], // 1 element - mismatch!
            vector[true, false], // 2 elements
            env.scenario.ctx()
        );

        env.tear_down();
        ts::return_to_address(OWNER, owner_cap);
    }

    #[test]
    #[expected_failure(abort_code = onramp::EInvalidDestChainSelector)]
    public fun test_error_invalid_dest_chain_selector() {
        let (mut env, nonce_manager_cap, source_transfer_cap) = setup();
        let owner_cap = ts::take_from_sender<OwnerCap>(&env.scenario);

        onramp::initialize(
            &mut env.state,
            &owner_cap,
            nonce_manager_cap,
            source_transfer_cap,
            123, // chain_selector
            FEE_AGGREGATOR,
            ALLOWLIST_ADMIN,
            vector[0], // Invalid zero chain selector
            vector[true],
            vector[false],
            env.scenario.ctx()
        );

        env.tear_down();
        ts::return_to_address(OWNER, owner_cap);
    }

    #[test]
    #[expected_failure(abort_code = onramp::EUnknownDestChainSelector)]
    public fun test_error_unknown_dest_chain_selector() {
        let (mut env, nonce_manager_cap, source_transfer_cap) = setup();
        let owner_cap = ts::take_from_sender<OwnerCap>(&env.scenario);
        initialize_onramp(&mut env, &owner_cap, nonce_manager_cap, source_transfer_cap);

        // Test EUnknownDestChainSelector (3) - querying unknown chain
        let (_enabled, _seq, _allowlist_enabled, _allowed_senders): (bool, u64, bool, vector<address>) = 
            onramp::get_dest_chain_config(&env.state, 999); // Unknown chain

        env.tear_down();
        ts::return_to_address(OWNER, owner_cap);
    }

    #[test]
    #[expected_failure(abort_code = ccip::fee_quoter::EUnknownDestChainSelector)]
    public fun test_error_sender_not_allowed() {
        let (mut env, nonce_manager_cap, source_transfer_cap) = setup();
        let owner_cap = ts::take_from_sender<OwnerCap>(&env.scenario);
        initialize_onramp(&mut env, &owner_cap, nonce_manager_cap, source_transfer_cap);

        // Create a test token for ccip_send
        env.scenario.next_tx(OWNER);
        let ctx = env.scenario.ctx();
        let (mut treasury_cap, coin_metadata) = create_test_token(ctx);
        let mut test_coin = coin::mint(&mut treasury_cap, 1000, ctx);

        // Switch to unauthorized sender
        env.scenario.next_tx(@0x999); // Not in allowlist for DEST_CHAIN_SELECTOR_1
        
        // Create empty token params for testing
        let token_params = dd::create_token_params(DEST_CHAIN_SELECTOR_1, EVM_RECEIVER_ADDRESS);

        // Test ccip_send function - this will fail on fee quoter validation first
        // before reaching the sender allowlist check
        let Env { scenario: _, state, registry: _, ref, clock } = &mut env;
        
        // This should fail due to unconfigured destination chain in fee quoter
        onramp::ccip_send(
            ref,
            state,
            clock,
            b"data",
            token_params,
            &coin_metadata,
            &mut test_coin,
            b"extra_args",
            env.scenario.ctx()
        );

        // Clean up test objects (won't be reached due to expected failure)
        transfer::public_transfer(treasury_cap, OWNER);
        transfer::public_freeze_object(coin_metadata);
        transfer::public_transfer(test_coin, OWNER);
        env.tear_down();
        ts::return_to_address(OWNER, owner_cap);
    }

    #[test]
    #[expected_failure(abort_code = onramp::EInvalidAllowlistAddress)]
    public fun test_error_invalid_allowlist_address() {
        let (mut env, nonce_manager_cap, source_transfer_cap) = setup();
        let owner_cap = ts::take_from_sender<OwnerCap>(&env.scenario);
        initialize_onramp(&mut env, &owner_cap, nonce_manager_cap, source_transfer_cap);

        // Test EInvalidAllowlistAddress (8) - zero address in allowlist
        onramp::apply_allowlist_updates(
            &mut env.state,
            &owner_cap,
            vector[DEST_CHAIN_SELECTOR_1],
            vector[true], // enable allowlist
            vector[vector[@0x0]], // Invalid zero address
            vector[vector[]],
            env.scenario.ctx()
        );

        env.tear_down();
        ts::return_to_address(OWNER, owner_cap);
    }

    #[test]
    #[expected_failure(abort_code = sui::dynamic_field::EFieldDoesNotExist)]
    public fun test_error_unexpected_withdraw_amount() {
        let (mut env, nonce_manager_cap, source_transfer_cap) = setup();
        let owner_cap = ts::take_from_sender<OwnerCap>(&env.scenario);
        initialize_onramp(&mut env, &owner_cap, nonce_manager_cap, source_transfer_cap);

        // Create a test token
        env.scenario.next_tx(OWNER);
        let ctx = env.scenario.ctx();
        let (treasury_cap, coin_metadata) = create_test_token(ctx);

        // Test withdraw_fee_tokens when no fees exist
        // This should fail because there are no fee tokens to withdraw
        onramp::withdraw_fee_tokens(&mut env.state, &owner_cap, &coin_metadata);

        // Clean up test objects (won't be reached due to expected failure)
        transfer::public_transfer(treasury_cap, OWNER);
        transfer::public_freeze_object(coin_metadata);
        env.tear_down();
        ts::return_to_address(OWNER, owner_cap);
    }
}
