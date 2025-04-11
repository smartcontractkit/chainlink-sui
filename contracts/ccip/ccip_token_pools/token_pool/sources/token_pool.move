module ccip_token_pool::token_pool {
    use sui::clock::Clock;
    use sui::coin::CoinMetadata;
    use sui::event;
    use sui::vec_map::{Self, VecMap};

    use ccip::eth_abi;
    use ccip::state_object;
    use ccip::token_admin_registry;
    use ccip::rmn_remote;
    use ccip::allowlist;

    use ccip_token_pool::token_pool_rate_limiter;

    // TODO: verify if this needs to be an object
    public struct TokenPoolState has store {
        allowlist_state: allowlist::AllowlistState,
        // TODO: check if we need to store decimals here
        coin_metadata: address,
        remote_chain_configs: VecMap<u64, RemoteChainConfig>,
        rate_limiter_config: token_pool_rate_limiter::RateLimitState
    }

    public struct RemoteChainConfig has store, drop, copy {
        remote_token_address: vector<u8>,
        remote_pools: vector<vector<u8>>
    }

    public struct Locked has copy, drop {
        local_token: address,
        amount: u64
    }

    public struct Released has copy, drop {
        local_token: address,
        recipient: address,
        amount: u64
    }

    public struct AllowlistRemove has copy, drop {
        sender: address
    }

    public struct AllowlistAdd has copy, drop {
        sender: address
    }

    public struct RemotePoolAdded has copy, drop {
        remote_chain_selector: u64,
        remote_pool_address: vector<u8>
    }

    public struct RemotePoolRemoved has copy, drop {
        remote_chain_selector: u64,
        remote_pool_address: vector<u8>
    }

    public struct ChainAdded has copy, drop {
        remote_chain_selector: u64,
        remote_token_address: vector<u8>
    }

    const E_NOT_PUBLISHER: u64 = 1;
    const E_UNKNOWN_FUNGIBLE_ASSET: u64 = 2;
    const E_UNKNOWN_REMOTE_CHAIN_SELECTOR: u64 = 3;
    const E_ZERO_ADDRESS_NOT_ALLOWED: u64 = 4;
    const E_REMOTE_POOL_ALREADY_ADDED: u64 = 5;
    const E_UNKNOWN_REMOTE_POOL: u64 = 6;
    const E_REMOTE_CHAIN_TO_ADD_MISMATCH: u64 = 7;
    const E_REMOTE_CHAIN_ALREADY_EXISTS: u64 = 8;
    const E_INVALID_REMOTE_CHAIN_DECIMALS: u64 = 9;
    const E_INVALID_ENCODED_AMOUNT: u64 = 10;

    // ================================================================
    // |                    Initialize and state                      |
    // ================================================================

    public fun initialize<T>(
        coin_metadata: &CoinMetadata<T>, allowlist: vector<address>, ctx: &mut TxContext
    ): TokenPoolState {
        assert_can_initialize(ctx.sender());

        TokenPoolState {
            allowlist_state: allowlist::new(allowlist, ctx),
            coin_metadata: object::id_to_address(&object::id(coin_metadata)),
            remote_chain_configs: vec_map::empty<u64, RemoteChainConfig>(),
            rate_limiter_config: token_pool_rate_limiter::new(ctx)
        }
    }

    // TODO: figure out the assertion rules
    fun assert_can_initialize(caller_address: address) {
        if (caller_address == @ccip_token_pool) { return };

        // if (object::is_object(@ccip_token_pool)) {
        //     let token_pool_object =
        //         object::address_to_object<ObjectCore>(@ccip_token_pool);
        //     if (caller_address == object::owner(token_pool_object)
        //         || caller_address == object::root_owner(token_pool_object)) { return };
        // };

        abort E_NOT_PUBLISHER
    }

    public fun get_router(): address {
        @ccip
    }

    public fun get_token(state: &TokenPoolState): address {
        state.coin_metadata
    }

    public fun get_token_decimals<T>(coin_metadata: &CoinMetadata<T>): u8 {
        coin_metadata.get_decimals()
    }

    // ================================================================
    // |                        Remote Chains                         |
    // ================================================================

    public fun get_supported_chains(state: &TokenPoolState): vector<u64> {
        state.remote_chain_configs.keys()
    }

    public fun is_supported_chain(
        state: &TokenPoolState, remote_chain_selector: u64
    ): bool {
        state.remote_chain_configs.contains(&remote_chain_selector)
    }

    public fun apply_chain_updates(
        state: &mut TokenPoolState,
        remote_chain_selectors_to_remove: vector<u64>,
        remote_chain_selectors_to_add: vector<u64>,
        remote_pool_addresses_to_add: vector<vector<vector<u8>>>,
        remote_token_addresses_to_add: vector<vector<u8>>
    ) {
        remote_chain_selectors_to_remove.do_ref!(
            |remote_chain_selector| {
                assert!(
                    state.remote_chain_configs.contains(remote_chain_selector),
                    E_UNKNOWN_REMOTE_CHAIN_SELECTOR
                );
                state.remote_chain_configs.remove(remote_chain_selector);
            }
        );

        let add_len = remote_chain_selectors_to_add.length();
        assert!(
            add_len == remote_pool_addresses_to_add.length(),
            E_REMOTE_CHAIN_TO_ADD_MISMATCH
        );
        assert!(
            add_len == remote_token_addresses_to_add.length(),
            E_REMOTE_CHAIN_TO_ADD_MISMATCH
        );

        let mut i = 0;
        while (i < add_len) {
            let remote_chain_selector = remote_chain_selectors_to_add[i];
            assert!(
                !state.remote_chain_configs.contains(&remote_chain_selector),
                E_REMOTE_CHAIN_ALREADY_EXISTS
            );
            let remote_pool_addresses = remote_pool_addresses_to_add[i];
            let remote_token_address = remote_token_addresses_to_add[i];
            assert!(
                !remote_token_address.is_empty(),
                E_ZERO_ADDRESS_NOT_ALLOWED
            );

            let mut remote_chain_config = RemoteChainConfig {
                remote_token_address,
                remote_pools: vector[]
            };

            remote_pool_addresses.do_ref!(
                |remote_pool_address| {
                    let remote_pool_address: vector<u8> = *remote_pool_address;
                    let (found, _) =
                        remote_chain_config.remote_pools.index_of(&remote_pool_address);
                    assert!(!found, E_REMOTE_POOL_ALREADY_ADDED);

                    remote_chain_config.remote_pools.push_back(remote_pool_address);

                    event::emit(
                        RemotePoolAdded { remote_chain_selector, remote_pool_address }
                    );
                }
            );

            state.remote_chain_configs.insert(remote_chain_selector, remote_chain_config);

            event::emit(ChainAdded { remote_chain_selector, remote_token_address });
            i = i + 1;
        };
    }

    // ================================================================
    // |                        Remote Pools                          |
    // ================================================================

    public fun get_remote_pools(
        state: &TokenPoolState, remote_chain_selector: u64
    ): vector<vector<u8>> {
        assert!(
            state.remote_chain_configs.contains(&remote_chain_selector),
            E_UNKNOWN_REMOTE_CHAIN_SELECTOR
        );
        let remote_chain_config =
            state.remote_chain_configs.get(&remote_chain_selector);
        remote_chain_config.remote_pools
    }

    public fun is_remote_pool(
        state: &TokenPoolState, remote_chain_selector: u64, remote_pool_address: vector<u8>
    ): bool {
        let remote_pools = get_remote_pools(state, remote_chain_selector);
        let (found, _) = remote_pools.index_of(&remote_pool_address);
        found
    }

    public fun get_remote_token(
        state: &TokenPoolState, remote_chain_selector: u64
    ): vector<u8> {
        assert!(
            state.remote_chain_configs.contains(&remote_chain_selector),
            E_UNKNOWN_REMOTE_CHAIN_SELECTOR
        );
        let remote_chain_config =
            state.remote_chain_configs.get(&remote_chain_selector);
        remote_chain_config.remote_token_address
    }

    public fun add_remote_pool(
        state: &mut TokenPoolState,
        remote_chain_selector: u64,
        remote_pool_address: vector<u8>
    ) {
        assert!(
            !remote_pool_address.is_empty(),
            E_ZERO_ADDRESS_NOT_ALLOWED
        );

        assert!(
            state.remote_chain_configs.contains(&remote_chain_selector),
            E_UNKNOWN_REMOTE_CHAIN_SELECTOR
        );
        let remote_chain_config =
            state.remote_chain_configs.get_mut(&remote_chain_selector);

        let (found, _) = remote_chain_config.remote_pools.index_of(&remote_pool_address);
        assert!(!found, E_REMOTE_POOL_ALREADY_ADDED);

        remote_chain_config.remote_pools.push_back(remote_pool_address);

        event::emit(RemotePoolAdded { remote_chain_selector, remote_pool_address });
    }

    public fun remove_remote_pool(
        state: &mut TokenPoolState,
        remote_chain_selector: u64,
        remote_pool_address: vector<u8>
    ) {
        assert!(
            state.remote_chain_configs.contains(&remote_chain_selector),
            E_UNKNOWN_REMOTE_CHAIN_SELECTOR
        );
        let remote_chain_config =
            state.remote_chain_configs.get_mut(&remote_chain_selector);

        let (found, i) = remote_chain_config.remote_pools.index_of(&remote_pool_address);
        assert!(found, E_UNKNOWN_REMOTE_POOL);

        // remove instead of swap_remove for readability, so the newest added pool is always at the end.
        remote_chain_config.remote_pools.remove(i);

        event::emit(RemotePoolRemoved { remote_chain_selector, remote_pool_address });
    }

    // ================================================================
    // |                         Validation                           |
    // ================================================================

    // Returns the remote token as bytes
    public fun validate_lock_or_burn<T>(
        ref: &state_object::CCIPObjectRef,
        clock: &Clock,
        state: &mut TokenPoolState,
        coin_metadata: &CoinMetadata<T>, // we can pass only the ID but it still does not remove the type param
        input: &token_admin_registry::LockOrBurnInputV1,
        local_amount: u64
    ): vector<u8> {
        // Validate the fungible asset
        let configured_token = get_token(state);

        // make sure the caller is requesting this pool's fungible asset.
        assert!(
            configured_token == object::id_to_address(object::borrow_id(coin_metadata)),
            E_UNKNOWN_FUNGIBLE_ASSET
        );

        // Check RMN curse status
        let remote_chain_selector =
            token_admin_registry::get_lock_or_burn_remote_chain_selector(input);
        assert!(!rmn_remote::is_cursed_u128(ref, (remote_chain_selector as u128)));

        // Allowlist check
        let _sender = token_admin_registry::get_lock_or_burn_sender(input);
        if (allowlist::get_allowlist_enabled(&state.allowlist_state)) {
            assert!(
                allowlist::is_allowed(&state.allowlist_state, _sender),
                E_NOT_PUBLISHER
            );
        };

        if (!is_supported_chain(state, remote_chain_selector)) {
            abort E_UNKNOWN_REMOTE_CHAIN_SELECTOR
        };

        token_pool_rate_limiter::consume_outbound(
            clock, &mut state.rate_limiter_config, remote_chain_selector, local_amount
        );

        get_remote_token(state, remote_chain_selector)
    }

    public fun validate_release_or_mint(
        ref: &state_object::CCIPObjectRef,
        clock: &Clock,
        state: &mut TokenPoolState,
        input: &token_admin_registry::ReleaseOrMintInputV1,
        local_amount: u64
    ) {
        // Validate the fungible asset
        let local_token = token_admin_registry::get_release_or_mint_local_token(input);
        let configured_token = get_token(state);

        // make sure the caller is requesting this pool's fungible asset.
        assert!(
            configured_token == local_token,
            E_UNKNOWN_FUNGIBLE_ASSET
        );

        // Check RMN curse status
        let remote_chain_selector =
            token_admin_registry::get_release_or_mint_remote_chain_selector(input);
        assert!(!rmn_remote::is_cursed_u128(ref, (remote_chain_selector as u128)));

        let source_pool_address =
            token_admin_registry::get_release_or_mint_source_pool_address(input);

        // This checks if the remote chain selector and the source pool are valid.
        assert!(
            is_remote_pool(state, remote_chain_selector, source_pool_address),
            E_UNKNOWN_REMOTE_POOL
        );

        token_pool_rate_limiter::consume_inbound(
            clock, &mut state.rate_limiter_config, remote_chain_selector, local_amount
        );
    }

    // ================================================================
    // |                           Events                             |
    // ================================================================

    public fun emit_released_or_minted(
        state: &mut TokenPoolState, recipient: address, amount: u64
    ) {
        event::emit(Released { local_token: state.coin_metadata, recipient, amount });
    }

    public fun emit_locked_or_burned(
        state: &mut TokenPoolState, amount: u64
    ) {
        event::emit(Locked { local_token: state.coin_metadata, amount });
    }

    // ================================================================
    // |                          Decimals                            |
    // ================================================================

    // for a token, CoinMetadata is supposed to be shared
    public fun encode_local_decimals<T>(coin_metadata: &CoinMetadata<T>): vector<u8> {
        let decimals = coin_metadata.get_decimals();
        let mut ret = vector[];
        eth_abi::encode_u8(&mut ret, decimals);
        ret
    }

    public fun parse_remote_decimals(
        source_pool_data: vector<u8>, local_decimals: u8
    ): u8 {
        let data_len = source_pool_data.length();
        if (data_len == 0) {
            // Fallback to the local value.
            return local_decimals
        };

        assert!(data_len == 32, E_INVALID_REMOTE_CHAIN_DECIMALS);

        let remote_decimals = eth_abi::decode_u256_value(source_pool_data);
        assert!(
            remote_decimals <= 255,
            E_INVALID_REMOTE_CHAIN_DECIMALS
        );

        remote_decimals as u8
    }

    public fun calculate_local_amount(
        remote_amount: u256, remote_decimals: u8, local_decimals: u8
    ): u64 {
        let local_amount =
            calculate_local_amount_internal(
                remote_amount, remote_decimals, local_decimals
            );
        // check that the calculated amount fits in a u64
        assert!(
            local_amount <= 18446744073709551615,
            E_INVALID_ENCODED_AMOUNT
        );
        local_amount as u64
    }

    fun calculate_local_amount_internal(
        remote_amount: u256, remote_decimals: u8, local_decimals: u8
    ): u256 {
        // TODO: check for overflows
        if (remote_decimals == local_decimals) {
            return remote_amount
        } else if (remote_decimals > local_decimals) {
            let decimals_diff = remote_decimals - local_decimals;
            let mut current_amount = remote_amount;
            let mut i = 0;
            while (i < decimals_diff) {
                current_amount = current_amount / 10;
                i = i + 1;
            };
            return current_amount
        } else {
            let decimals_diff = local_decimals - remote_decimals;
            let mut current_amount = remote_amount;
            let mut i = 0;
            while (i < decimals_diff) {
                current_amount = current_amount * 10;
                i = i + 1;
            };
            return current_amount
        }
    }

    public fun calculate_release_or_mint_amount<T>(
        coin_metadata: &CoinMetadata<T>, state: &TokenPoolState, input: &token_admin_registry::ReleaseOrMintInputV1
    ): u64 {
        // make sure the caller is requesting this pool's fungible asset.
        assert!(
            get_token(state) == object::id_to_address(object::borrow_id(coin_metadata)),
            E_UNKNOWN_FUNGIBLE_ASSET
        );

        let local_decimals = get_token_decimals(coin_metadata);
        let source_amount =
            token_admin_registry::get_release_or_mint_source_amount(input);
        let source_pool_data =
            token_admin_registry::get_release_or_mint_source_pool_data(input);
        let remote_decimals = parse_remote_decimals(source_pool_data, local_decimals);
        let local_amount =
            calculate_local_amount(source_amount, remote_decimals, local_decimals);
        local_amount
    }

    // ================================================================
    // |                    Rate limit config                         |
    // ================================================================

    public fun set_chain_rate_limiter_config(
        clock: &Clock,
        state: &mut TokenPoolState,
        remote_chain_selector: u64,
        outbound_is_enabled: bool,
        outbound_capacity: u64,
        outbound_rate: u64,
        inbound_is_enabled: bool,
        inbound_capacity: u64,
        inbound_rate: u64
    ) {
        token_pool_rate_limiter::set_chain_rate_limiter_config(
            clock,
            &mut state.rate_limiter_config,
            remote_chain_selector,
            outbound_is_enabled,
            outbound_capacity,
            outbound_rate,
            inbound_is_enabled,
            inbound_capacity,
            inbound_rate
        );
    }

    // ================================================================
    // |                          Allowlist                           |
    // ================================================================

    public fun get_allowlist_enabled(state: &TokenPoolState): bool {
        allowlist::get_allowlist_enabled(&state.allowlist_state)
    }

    public fun get_allowlist(state: &TokenPoolState): vector<address> {
        allowlist::get_allowlist(&state.allowlist_state)
    }

    public fun apply_allowlist_updates(
        state: &mut TokenPoolState, removes: vector<address>, adds: vector<address>
    ) {
        allowlist::apply_allowlist_updates(&mut state.allowlist_state, removes, adds);
    }

    // ================================================================
    // |                          Destroy                             |
    // ================================================================

    public fun destroy_token_pool(state: TokenPoolState) {
        let TokenPoolState {
            allowlist_state,
            coin_metadata: _coin_metadata,
            remote_chain_configs: _remote_chain_configs,
            rate_limiter_config
        } = state;

        allowlist::destroy_allowlist(allowlist_state);

        token_pool_rate_limiter::destroy_rate_limiter(rate_limiter_config);
    }
}

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

        let mut state = token_pool::initialize(&coin_metadata, vector[], ctx);

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
}
