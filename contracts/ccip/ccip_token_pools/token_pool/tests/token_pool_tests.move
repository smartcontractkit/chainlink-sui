#[test_only]
module ccip_token_pool::token_pool_test {
    use sui::coin;
    use sui::test_scenario::{Self, Scenario};

    use ccip_token_pool::token_pool::{Self, TokenPoolState};

    public struct TOKEN_POOL_TEST has drop {}

    const Decimals: u8 = 8;

    const DefaultRemoteChain: u64 = 2000;
    const DefaultRemoteToken: vector<u8> = b"default_remote_token";
    const DefaultRemotePool: vector<u8> = b"default_remote_pool";

    fun set_up_test(): (Scenario, TokenPoolState) {
        let mut scenario = test_scenario::begin(@ccip_token_pool);
        let ctx = scenario.ctx();

        let (treasury_cap, coin_metadata) = coin::create_currency(
            TOKEN_POOL_TEST {},
            Decimals,
            b"TEST",
            b"TestToken",
            b"test_token",
            option::none(),
            ctx
        );

        let mut state = token_pool::initialize(object::id_to_address(&object::id(&coin_metadata)), Decimals, vector[], ctx);

        // Set state in the pool
        set_up_default_remote_chain(&mut state);

        transfer::public_freeze_object(coin_metadata);
        transfer::public_transfer(treasury_cap, ctx.sender());

        (scenario, state)
    }

    fun set_up_default_remote_chain(state: &mut TokenPoolState) {
        token_pool::apply_chain_updates(
            state,
            vector[],
            vector[DefaultRemoteChain],
            vector[vector[DefaultRemotePool]],
            vector[DefaultRemoteToken]
        )
    }

    #[test]
    public fun initialize_correctly_sets_state() {
        let (scenario, state) = set_up_test();

        assert!(token_pool::is_supported_chain(&state, DefaultRemoteChain), 1);

        token_pool::destroy_token_pool(state);
        scenario.end();
    }

    #[test]
    fun add_remote_pool_existing_chain() {
        let (scenario, mut state) = set_up_test();
        let new_remote_pool = b"new_pool";

        assert!(
            !token_pool::is_remote_pool(&state, DefaultRemoteChain, new_remote_pool),
            1
        );
        assert!(
            token_pool::get_remote_pools(&state, DefaultRemoteChain).length() == 1,
            1
        );

        token_pool::add_remote_pool(&mut state, DefaultRemoteChain, new_remote_pool);

        assert!(
            token_pool::is_remote_pool(&state, DefaultRemoteChain, new_remote_pool),
            1
        );
        assert!(
            token_pool::get_remote_pools(&state, DefaultRemoteChain).length() == 2,
            1
        );
        assert!(token_pool::is_supported_chain(&state, DefaultRemoteChain), 1);

        token_pool::destroy_token_pool(state);
        scenario.end();
    }

    #[test]
    fun apply_chain_updates() {
        let (scenario, mut state) = set_up_test();
        let new_remote_chain = 3000;
        let new_remote_token = b"new_remote_token";
        let new_remote_pool = b"new_remote_pool";
        let new_remote_pool_2 = b"new_remote_pool_2";

        assert!(!token_pool::is_supported_chain(&state, new_remote_chain));

        token_pool::apply_chain_updates(
            &mut state,
            vector[],
            vector[new_remote_chain],
            vector[vector[new_remote_pool, new_remote_pool_2]],
            vector[new_remote_token]
        );
        assert!(token_pool::is_supported_chain(&state, new_remote_chain));
        assert!(token_pool::get_remote_pools(&state, new_remote_chain).length() == 2);
        assert!(token_pool::get_remote_token(&state, new_remote_chain) == new_remote_token);

        token_pool::apply_chain_updates(
            &mut state,
            vector[new_remote_chain],
            vector[],
            vector[],
            vector[]
        );
        assert!(!token_pool::is_supported_chain(&state, new_remote_chain));

        token_pool::destroy_token_pool(state);
        scenario.end();
    }

    #[test]
    fun test_calculate_local_amount_same_decimals() {
        // When remote and local decimals are the same, amount should not change
        let remote_amount: u256 = 1000000;
        let remote_decimals: u8 = 8;
        let local_decimals: u8 = 8;

        let local_amount =
            token_pool::calculate_local_amount(
                remote_amount, remote_decimals, local_decimals
            );
        assert!(local_amount == 1000000, 0);
    }

    #[test]
    fun test_calculate_local_amount_more_decimals() {
        // When local has more decimals, amount should increase
        let remote_amount: u256 = 1000000;
        let remote_decimals: u8 = 6; // 6 decimals
        let local_decimals: u8 = 8; // 8 decimals (2 more)

        let local_amount =
            token_pool::calculate_local_amount(
                remote_amount, remote_decimals, local_decimals
            );
        assert!(local_amount == 100000000, 0); // 1000000 * 10^2
    }

    #[test]
    fun test_calculate_local_amount_fewer_decimals() {
        // When local has fewer decimals, amount should decrease
        let remote_amount: u256 = 1000000;
        let remote_decimals: u8 = 8; // 8 decimals
        let local_decimals: u8 = 6; // 6 decimals (2 fewer)

        let local_amount =
            token_pool::calculate_local_amount(
                remote_amount, remote_decimals, local_decimals
            );
        assert!(local_amount == 10000, 0); // 1000000 / 10^2
    }

    #[test]
    #[expected_failure(abort_code = token_pool::E_DECIMAL_OVERFLOW)]
    fun test_decimal_overflow_protection() {
        // Test for overflow protection - when decimal difference exceeds MAX_SAFE_DECIMAL_DIFF
        let remote_amount: u256 = 1000000;
        let remote_decimals: u8 = 1; // 1 decimal
        let local_decimals: u8 = 100; // 100 decimals (99 more - exceeds the limit of 77)

        // E_DECIMAL_OVERFLOW error
        let _local_amount =
            token_pool::calculate_local_amount(
                remote_amount, remote_decimals, local_decimals
            );
    }

    #[test]
    #[expected_failure(abort_code = token_pool::E_INVALID_ENCODED_AMOUNT)]
    fun test_local_amount_u64_overflow() {
        let remote_amount: u256 = 0xffffffffffffffffffffffffffffffff;
        let remote_decimals: u8 = 0;
        let local_decimals: u8 = 18;

        // E_INVALID_ENCODED_AMOUNT error
        let _local_amount =
            token_pool::calculate_local_amount(
                remote_amount, remote_decimals, local_decimals
            );
    }
}