#[test_only]
module burn_mint_token_pool::burn_mint_token_pool_ownable_test {
    use sui::test_scenario::{Self as ts, Scenario};
    use sui::coin;
    use ccip::state_object::{Self, CCIPObjectRef};
    use ccip::token_admin_registry;
    use ccip::rmn_remote;

    use burn_mint_token_pool::burn_mint_token_pool::{Self, BurnMintTokenPoolState};
    use ccip_token_pool::ownable::{Self, OwnerCap};

    public struct BURN_MINT_TOKEN_POOL_OWNABLE_TEST has drop {}

    const OWNER: address = @0x123;
    const NEW_OWNER: address = @0x456;
    const OTHER_USER: address = @0x789;
    const Decimals: u8 = 8;

    public struct TestEnv {
        scenario: Scenario,
        state: BurnMintTokenPoolState<BURN_MINT_TOKEN_POOL_OWNABLE_TEST>,
        ccip_ref: CCIPObjectRef,
    }

    fun setup(): (TestEnv, OwnerCap) {
        let mut scenario = ts::begin(OWNER);
        let ctx = scenario.ctx();

        // Setup CCIP environment
        state_object::test_init(ctx);
        
        scenario.next_tx(OWNER);
        let ccip_owner_cap = scenario.take_from_sender<state_object::OwnerCap>();
        let mut ccip_ref = scenario.take_shared<CCIPObjectRef>();
        
        // Initialize required CCIP modules
        token_admin_registry::initialize(&mut ccip_ref, &ccip_owner_cap, scenario.ctx());
        rmn_remote::initialize(&mut ccip_ref, &ccip_owner_cap, 1000, scenario.ctx());
        
        scenario.next_tx(OWNER);
        
        // Create token and initialize pool
        let (treasury_cap, coin_metadata) = coin::create_currency(
            BURN_MINT_TOKEN_POOL_OWNABLE_TEST {},
            Decimals,
            b"BMTP",
            b"BurnMintTestToken",
            b"burn_mint_test_token",
            option::none(),
            scenario.ctx()
        );
        
        burn_mint_token_pool::initialize_by_ccip_admin(
            &mut ccip_ref,
            &coin_metadata,
            treasury_cap,
            @burn_mint_token_pool,
            @0x123,
            vector[],
            vector[],
            scenario.ctx()
        );
        
        transfer::public_freeze_object(coin_metadata);
        transfer::public_transfer(ccip_owner_cap, @0x0);

        scenario.next_tx(OWNER);
        let state = scenario.take_shared<BurnMintTokenPoolState<BURN_MINT_TOKEN_POOL_OWNABLE_TEST>>();
        let owner_cap = scenario.take_from_sender<OwnerCap>();

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
        assert!(burn_mint_token_pool::owner(&env.state) == OWNER);
        assert!(!burn_mint_token_pool::has_pending_transfer(&env.state));
        assert!(burn_mint_token_pool::pending_transfer_from(&env.state).is_none());
        assert!(burn_mint_token_pool::pending_transfer_to(&env.state).is_none());
        assert!(burn_mint_token_pool::pending_transfer_accepted(&env.state).is_none());

        tear_down(env);
        ts::return_to_address(OWNER, owner_cap);
    }

    #[test]
    public fun test_ownership_transfer_flow() {
        let (mut env, owner_cap) = setup();

        // Transfer ownership
        burn_mint_token_pool::transfer_ownership(&mut env.state, &owner_cap, NEW_OWNER, env.scenario.ctx());

        // Verify pending transfer
        assert!(burn_mint_token_pool::has_pending_transfer(&env.state));
        assert!(burn_mint_token_pool::pending_transfer_from(&env.state).extract() == OWNER);
        assert!(burn_mint_token_pool::pending_transfer_to(&env.state).extract() == NEW_OWNER);
        assert!(burn_mint_token_pool::pending_transfer_accepted(&env.state).extract() == false);

        // Accept ownership as new owner
        env.scenario.next_tx(NEW_OWNER);
        burn_mint_token_pool::accept_ownership(&mut env.state, env.scenario.ctx());

        // Verify acceptance
        assert!(burn_mint_token_pool::pending_transfer_accepted(&env.state).extract() == true);

        tear_down(env);
        ts::return_to_address(OWNER, owner_cap);
    }

    #[test]
    public fun test_accept_ownership_from_object() {
        let (mut env, owner_cap) = setup();

        // Create a UID object for testing
        let mut uid_object = object::new(env.scenario.ctx());
        let uid_address = uid_object.to_address();

        // Transfer ownership to the UID object's address
        burn_mint_token_pool::transfer_ownership(&mut env.state, &owner_cap, uid_address, env.scenario.ctx());

        // Accept ownership from the UID object
        burn_mint_token_pool::accept_ownership_from_object(&mut env.state, &mut uid_object, env.scenario.ctx());

        // Verify acceptance
        assert!(burn_mint_token_pool::pending_transfer_accepted(&env.state).extract() == true);

        object::delete(uid_object);
        tear_down(env);
        ts::return_to_address(OWNER, owner_cap);
    }

    #[test]
    #[expected_failure(abort_code = ownable::EUnauthorizedOwnershipTransfer)]
    public fun test_transfer_ownership_unauthorized() {
        let (mut env, owner_cap) = setup();

        // Try to transfer ownership from unauthorized user
        env.scenario.next_tx(OTHER_USER);
        burn_mint_token_pool::transfer_ownership(&mut env.state, &owner_cap, NEW_OWNER, env.scenario.ctx());

        tear_down(env);
        ts::return_to_address(OWNER, owner_cap);
    }

    #[test]
    #[expected_failure(abort_code = ownable::EUnauthorizedAcceptance)]
    public fun test_accept_ownership_unauthorized() {
        let (mut env, owner_cap) = setup();

        // Transfer ownership
        burn_mint_token_pool::transfer_ownership(&mut env.state, &owner_cap, NEW_OWNER, env.scenario.ctx());

        // Try to accept ownership from unauthorized user
        env.scenario.next_tx(OTHER_USER);
        burn_mint_token_pool::accept_ownership(&mut env.state, env.scenario.ctx());

        tear_down(env);
        ts::return_to_address(OWNER, owner_cap);
    }

    #[test]
    #[expected_failure(abort_code = ownable::ENoPendingTransfer)]
    public fun test_accept_ownership_no_pending_transfer() {
        let (mut env, owner_cap) = setup();

        // Try to accept ownership without pending transfer
        env.scenario.next_tx(NEW_OWNER);
        burn_mint_token_pool::accept_ownership(&mut env.state, env.scenario.ctx());

        tear_down(env);
        ts::return_to_address(OWNER, owner_cap);
    }

    #[test]
    #[expected_failure(abort_code = ownable::ETransferAlreadyAccepted)]
    public fun test_accept_ownership_already_accepted() {
        let (mut env, owner_cap) = setup();

        // Transfer ownership
        burn_mint_token_pool::transfer_ownership(&mut env.state, &owner_cap, NEW_OWNER, env.scenario.ctx());

        // Accept ownership
        env.scenario.next_tx(NEW_OWNER);
        burn_mint_token_pool::accept_ownership(&mut env.state, env.scenario.ctx());

        // Try to accept again
        burn_mint_token_pool::accept_ownership(&mut env.state, env.scenario.ctx());

        tear_down(env);
        ts::return_to_address(OWNER, owner_cap);
    }

    #[test]
    public fun test_ownable_functions_with_owner_cap_validation() {
        let (mut env, owner_cap) = setup();

        // Test that ownable functions work with proper owner cap validation
        // These should succeed because we have the correct owner cap

        // First, configure a chain so we can add remote pools
        burn_mint_token_pool::apply_chain_updates(
            &mut env.state,
            &owner_cap,
            vector[], // remove nothing
            vector[1000], // add chain 1000
            vector[vector[b"test_remote_pool"]], // pool addresses
            vector[b"test_remote_token"] // token addresses
        );

        // Test add_remote_pool (now that chain is configured)
        burn_mint_token_pool::add_remote_pool(
            &mut env.state,
            &owner_cap,
            1000,
            b"additional_remote_pool"
        );

        // Test set_allowlist_enabled
        burn_mint_token_pool::set_allowlist_enabled(&mut env.state, &owner_cap, true);

        // Test apply_allowlist_updates
        burn_mint_token_pool::apply_allowlist_updates(
            &mut env.state,
            &owner_cap,
            vector[], // removes
            vector[@0x123] // adds
        );

        // Verify the changes were applied
        assert!(burn_mint_token_pool::get_allowlist_enabled(&env.state));
        let allowlist = burn_mint_token_pool::get_allowlist(&env.state);
        assert!(allowlist.length() == 1);
        assert!(allowlist[0] == @0x123);

        // Verify remote pool was added
        assert!(burn_mint_token_pool::is_remote_pool(&env.state, 1000, b"additional_remote_pool"));

        tear_down(env);
        ts::return_to_address(OWNER, owner_cap);
    }
} 