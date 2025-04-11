/// This module is responsible for storage and retrieval of fee token and token transfer
/// information and pricing.
module ccip::fee_quoter {
    use std::string::{Self, String};
    use sui::clock;
    use sui::event;
    use sui::table;

    use ccip::eth_abi;
    use ccip::internal;
    use ccip::state_object::{Self, CCIPObjectRef};

    const FEE_QUOTER_STATE_NAME: vector<u8> = b"FeeQuoterState";
    const CHAIN_FAMILY_SELECTOR_EVM: vector<u8> = x"2812d52c";
    const CHAIN_FAMILY_SELECTOR_SVM: vector<u8> = x"1e10bdc4";
    const CHAIN_FAMILY_SELECTOR_APTOS: vector<u8> = x"ac77ffec";
    const EVM_PRECOMPILE_SPACE: u256 = 1024;
    const SVM_EXTRA_ARGS_V1_TAG: vector<u8> = x"1f3b3aba";
    const GENERIC_EXTRA_ARGS_V2_TAG: vector<u8> = x"181dcf10";

    const GAS_PRICE_BITS: u8 = 112;

    const MESSAGE_FIXED_BYTES: u64 = 32 * 15;
    const MESSAGE_FIXED_BYTES_PER_TOKEN: u64 = 32 * (4 + (3 + 2));

    const CCIP_LOCK_OR_BURN_V1_RET_BYTES: u32 = 32;

    const MAX_U64: u256 = 18446744073709551615;
    const MAX_U160: u256 = 1461501637330902918203684832716283019655932542975;
    const MAX_U256: u256 =
        115792089237316195423570985008687907853269984665640564039457584007913129639935;
    const VAL_1E5: u256 = 100_000;
    const VAL_1E14: u256 = 100_000_000_000_000;
    const VAL_1E16: u256 = 10_000_000_000_000_000;
    const VAL_1E18: u256 = 1_000_000_000_000_000_000;

    public struct FeeQuoterState has key, store {
        id: UID,
        max_fee_juels_per_msg: u64,
        // TODO: figure out if we should use CoinMetadata for link token or object ID for Treasury Cap object
        // TODO: we will need a link token contract
        link_token: address,
        token_price_staleness_threshold: u64,
        fee_tokens: vector<address>,
        usd_per_unit_gas_by_dest_chain: table::Table<u64, TimestampedPrice>,
        // TODO: we need to know the token address for common tokens like USDC, SUI, etc
        usd_per_token: table::Table<address, TimestampedPrice>,
        dest_chain_configs: table::Table<u64, DestChainConfig>,
        // dest chain selector -> local token -> TokenTransferFeeConfig
        token_transfer_fee_configs: table::Table<u64, table::Table<address, TokenTransferFeeConfig>>,
        // TODO: update calculations - should this be octa per apt?
        premium_multiplier_wei_per_eth: table::Table<address, u64>,
    }

    public struct StaticConfig has drop {
        max_fee_juels_per_msg: u64,
        link_token: address,
        token_price_staleness_threshold: u64
    }

    public struct DestChainConfig has store, drop, copy {
        is_enabled: bool,
        max_number_of_tokens_per_msg: u16,
        max_data_bytes: u32,
        max_per_msg_gas_limit: u32,
        dest_gas_overhead: u32,
        dest_gas_per_payload_byte_base: u8,
        dest_gas_per_payload_byte_high: u8,
        dest_gas_per_payload_byte_threshold: u16,
        dest_data_availability_overhead_gas: u32,
        dest_gas_per_data_availability_byte: u16,
        dest_data_availability_multiplier_bps: u16,
        chain_family_selector: vector<u8>,
        enforce_out_of_order: bool,
        default_token_fee_usd_cents: u16,
        default_token_dest_gas_overhead: u32,
        default_tx_gas_limit: u32,
        // TODO: should this be octa per apt?
        gas_multiplier_wei_per_eth: u64,
        gas_price_staleness_threshold: u32,
        network_fee_usd_cents: u32
    }

    public struct TokenTransferFeeConfig has store, drop, copy {
        min_fee_usd_cents: u32,
        max_fee_usd_cents: u32,
        deci_bps: u16,
        dest_gas_overhead: u32,
        dest_bytes_overhead: u32,
        is_enabled: bool
    }

    public struct TimestampedPrice has store, drop, copy {
        value: u256,
        timestamp: u64
    }

    public struct FeeTokenAdded has copy, drop {
        fee_token: address
    }

    public struct FeeTokenRemoved has copy, drop {
        fee_token: address
    }

    public struct TokenTransferFeeConfigAdded has copy, drop {
        dest_chain_selector: u64,
        token: address,
        token_transfer_fee_config: TokenTransferFeeConfig
    }

    public struct TokenTransferFeeConfigRemoved has copy, drop {
        dest_chain_selector: u64,
        token: address
    }

    public struct UsdPerTokenUpdated has copy, drop {
        token: address,
        usd_per_token: u256,
        timestamp: u64
    }

    public struct UsdPerUnitGasUpdated has copy, drop {
        dest_chain_selector: u64,
        usd_per_unit_gas: u256,
        timestamp: u64
    }

    public struct DestChainAdded has copy, drop {
        dest_chain_selector: u64,
        dest_chain_config: DestChainConfig
    }

    public struct DestChainConfigUpdated has copy, drop {
        dest_chain_selector: u64,
        dest_chain_config: DestChainConfig
    }

    public struct PremiumMultiplierWeiPerEthUpdated has copy, drop {
        token: address,
        premium_multiplier_wei_per_eth: u64
    }

    const E_ALREADY_INITIALIZED: u64 = 1;
    // const E_INVALID_LINK_TOKEN: u64 = 2;
    const E_UNKNOWN_DEST_CHAIN_SELECTOR: u64 = 3;
    const E_UNKNOWN_TOKEN: u64 = 4;
    const E_DEST_CHAIN_NOT_ENABLED: u64 = 5;
    const E_TOKEN_UPDATE_MISMATCH: u64 = 6;
    const E_GAS_UPDATE_MISMATCH: u64 = 7;
    const E_TOKEN_TRANSFER_FEE_CONFIG_MISMATCH: u64 = 8;
    const E_FEE_TOKEN_NOT_SUPPORTED: u64 = 9;
    const E_TOKEN_NOT_SUPPORTED: u64 = 10;
    const E_UNKNOWN_CHAIN_FAMILY_SELECTOR: u64 = 11;
    const E_STALE_GAS_PRICE: u64 = 12;
    const E_MESSAGE_TOO_LARGE: u64 = 13;
    const E_UNSUPPORTED_NUMBER_OF_TOKENS: u64 = 14;
    const E_INVALID_EVM_ADDRESS: u64 = 15;
    const E_INVALID_SVM_ADDRESS: u64 = 16;
    const E_FEE_TOKEN_COST_TOO_HIGH: u64 = 17;
    const E_MESSAGE_GAS_LIMIT_TOO_HIGH: u64 = 18;
    const E_EXTRA_ARG_OUT_OF_ORDER_EXECUTION_MUST_BE_TRUE: u64 = 19;
    const E_INVALID_EXTRA_ARGS_TAG: u64 = 20;
    const E_INVALID_EXTRA_ARGS_DATA: u64 = 21;
    const E_INVALID_TOKEN_RECEIVER: u64 = 22;
    const E_MESSAGE_COMPUTE_UNIT_LIMIT_TOO_HIGH: u64 = 23;
    const E_MESSAGE_FEE_TOO_HIGH: u64 = 24;
    const E_SOURCE_TOKEN_DATA_TOO_LARGE: u64 = 25;
    const E_INVALID_DEST_CHAIN_SELECTOR: u64 = 26;
    const E_INVALID_GAS_LIMIT: u64 = 27;
    const E_INVALID_CHAIN_FAMILY_SELECTOR: u64 = 28;
    const E_TO_TOKEN_AMOUNT_TOO_LARGE: u64 = 29;
    // const E_UNKNOWN_FUNCTION: u64 = 30;
    const E_OUT_OF_BOUND: u64 = 31;
    const E_ONLY_CALLABLE_BY_OWNER: u64 = 32;

    public fun type_and_version(): String {
        string::utf8(b"FeeQuoter 1.6.0")
    }

    // fun init_module(publisher: &signer) {
    //     if (@mcms_register_entrypoints != @0x0) {
    //         mcms_registry::register_entrypoint(
    //             publisher, string::utf8(b"fee_quoter"), McmsCallback {}
    //         );
    //     };
    // }

    public fun initialize(
        ref: &mut CCIPObjectRef,
        max_fee_juels_per_msg: u64,
        link_token: address,
        token_price_staleness_threshold: u64,
        fee_tokens: vector<address>,
        ctx: &mut TxContext
    ) {
        assert!(
            ctx.sender() == state_object::get_current_owner(ref),
            E_ONLY_CALLABLE_BY_OWNER
        );
        assert!(
            !state_object::contains(ref, FEE_QUOTER_STATE_NAME),
            E_ALREADY_INITIALIZED
        );

        // TODO: how to perform some checks for link token. this technically just needs
        // to be a unique identifier of link token. we don't request any data from it.
        // Could be the object ID of treasury cap or coin metadata
        // assert!(
        //     object::object_exists<CoinMetadata>(link_token),
        //     E_INVALID_LINK_TOKEN
        // );

        let state = FeeQuoterState {
            id: object::new(ctx),
            max_fee_juels_per_msg,
            link_token,
            token_price_staleness_threshold,
            fee_tokens,
            usd_per_unit_gas_by_dest_chain: table::new(ctx),
            usd_per_token: table::new(ctx),
            dest_chain_configs: table::new(ctx),
            token_transfer_fee_configs: table::new(ctx),
            premium_multiplier_wei_per_eth: table::new(ctx),
        };
        state_object::add(ref, FEE_QUOTER_STATE_NAME, state, ctx);
    }

    public fun apply_fee_token_updates(
        ref: &mut CCIPObjectRef,
        fee_tokens_to_remove: vector<address>,
        fee_tokens_to_add: vector<address>,
        ctx: &mut TxContext
    ) {
        assert!(
            ctx.sender() == state_object::get_current_owner(ref),
            E_ONLY_CALLABLE_BY_OWNER
        );
        let state = state_object::borrow_mut_with_ctx<FeeQuoterState>(ref, FEE_QUOTER_STATE_NAME, ctx);

        // Remove tokens
        vector::do_ref!(
            &fee_tokens_to_remove,
            |fee_token| {
                let fee_token = *fee_token;
                let (found, index) = vector::index_of(&state.fee_tokens, &fee_token);
                if (found) {
                    vector::remove(&mut state.fee_tokens, index);
                    event::emit(FeeTokenRemoved { fee_token });
                };
            }
        );

        // Add new tokens
        vector::do_ref!(
            &fee_tokens_to_add,
            |fee_token| {
                let fee_token = *fee_token;
                let (found, _) = vector::index_of(&state.fee_tokens, &fee_token);
                if (!found) {
                    state.fee_tokens.push_back(fee_token);
                    event::emit(FeeTokenAdded { fee_token });
                };
            }
        );
    }

    // Note that unlike EVM, this only allows changes for a single dest chain selector at a time.
    public fun apply_token_transfer_fee_config_updates(
        ref: &mut CCIPObjectRef,
        dest_chain_selector: u64,
        add_tokens: vector<address>,
        add_min_fee_usd_cents: vector<u32>,
        add_max_fee_usd_cents: vector<u32>,
        add_deci_bps: vector<u16>,
        add_dest_gas_overhead: vector<u32>,
        add_dest_bytes_overhead: vector<u32>,
        add_is_enabled: vector<bool>,
        remove_tokens: vector<address>,
        ctx: &mut TxContext
    ) {
        assert!(
            ctx.sender() == state_object::get_current_owner(ref),
            E_ONLY_CALLABLE_BY_OWNER
        );
        let state = state_object::borrow_mut_with_ctx<FeeQuoterState>(ref, FEE_QUOTER_STATE_NAME, ctx);

        if (!table::contains(
            &state.token_transfer_fee_configs, dest_chain_selector
        )) {
            table::add(
                &mut state.token_transfer_fee_configs,
                dest_chain_selector,
                table::new(ctx)
            );
        };
        let token_transfer_fee_configs =
            table::borrow_mut(
                &mut state.token_transfer_fee_configs, dest_chain_selector
            );

        let add_tokens_len = add_tokens.length();
        assert!(
            add_tokens_len == add_min_fee_usd_cents.length(),
            E_TOKEN_TRANSFER_FEE_CONFIG_MISMATCH
        );
        assert!(
            add_tokens_len == add_max_fee_usd_cents.length(),
            E_TOKEN_TRANSFER_FEE_CONFIG_MISMATCH
        );
        assert!(
            add_tokens_len == add_deci_bps.length(),
            E_TOKEN_TRANSFER_FEE_CONFIG_MISMATCH
        );
        assert!(
            add_tokens_len == add_dest_gas_overhead.length(),
            E_TOKEN_TRANSFER_FEE_CONFIG_MISMATCH
        );
        assert!(
            add_tokens_len == add_dest_bytes_overhead.length(),
            E_TOKEN_TRANSFER_FEE_CONFIG_MISMATCH
        );
        assert!(
            add_tokens_len == add_is_enabled.length(),
            E_TOKEN_TRANSFER_FEE_CONFIG_MISMATCH
        );

        let mut i = 0;
        while (i < add_tokens_len) {
            let token = add_tokens[i];
            let min_fee_usd_cents = add_min_fee_usd_cents[i];
            let max_fee_usd_cents = add_max_fee_usd_cents[i];
            let deci_bps = add_deci_bps[i];
            let dest_gas_overhead = add_dest_gas_overhead[i];
            let dest_bytes_overhead = add_dest_bytes_overhead[i];
            let is_enabled = add_is_enabled[i];

            let token_transfer_fee_config = TokenTransferFeeConfig {
                min_fee_usd_cents,
                max_fee_usd_cents,
                deci_bps,
                dest_gas_overhead,
                dest_bytes_overhead,
                is_enabled
            };

            table::add(
                token_transfer_fee_configs, token, token_transfer_fee_config
            );

            event::emit(
                TokenTransferFeeConfigAdded {
                    dest_chain_selector,
                    token,
                    token_transfer_fee_config
                }
            );

            i = i + 1;
        };

        vector::do_ref!(
            &remove_tokens,
            |token| {
                let token = *token;
                if (table::contains(token_transfer_fee_configs, token)) {
                    table::remove(token_transfer_fee_configs, token);

                    event::emit(
                        TokenTransferFeeConfigRemoved { dest_chain_selector, token }
                    );
                }
            }
        );
    }

    public fun apply_dest_chain_config_updates(
        ref: &mut CCIPObjectRef,
        dest_chain_selector: u64,
        is_enabled: bool,
        max_number_of_tokens_per_msg: u16,
        max_data_bytes: u32,
        max_per_msg_gas_limit: u32,
        dest_gas_overhead: u32,
        dest_gas_per_payload_byte_base: u8,
        dest_gas_per_payload_byte_high: u8,
        dest_gas_per_payload_byte_threshold: u16,
        dest_data_availability_overhead_gas: u32,
        dest_gas_per_data_availability_byte: u16,
        dest_data_availability_multiplier_bps: u16,
        chain_family_selector: vector<u8>,
        enforce_out_of_order: bool,
        default_token_fee_usd_cents: u16,
        default_token_dest_gas_overhead: u32,
        default_tx_gas_limit: u32,
        gas_multiplier_wei_per_eth: u64,
        gas_price_staleness_threshold: u32,
        network_fee_usd_cents: u32,
        ctx: &mut TxContext
    ) {
        assert!(
            ctx.sender() == state_object::get_current_owner(ref),
            E_ONLY_CALLABLE_BY_OWNER
        );
        let state = state_object::borrow_mut_with_ctx<FeeQuoterState>(ref, FEE_QUOTER_STATE_NAME, ctx);

        assert!(
            dest_chain_selector != 0,
            E_INVALID_DEST_CHAIN_SELECTOR
        );
        assert!(
            default_tx_gas_limit != 0 && default_tx_gas_limit <= max_per_msg_gas_limit,
            E_INVALID_GAS_LIMIT
        );

        assert!(
            chain_family_selector == CHAIN_FAMILY_SELECTOR_EVM
                || chain_family_selector == CHAIN_FAMILY_SELECTOR_SVM
                || chain_family_selector == CHAIN_FAMILY_SELECTOR_APTOS,
            E_INVALID_CHAIN_FAMILY_SELECTOR
        );

        let dest_chain_config = DestChainConfig {
            is_enabled,
            max_number_of_tokens_per_msg,
            max_data_bytes,
            max_per_msg_gas_limit,
            dest_gas_overhead,
            dest_gas_per_payload_byte_base,
            dest_gas_per_payload_byte_high,
            dest_gas_per_payload_byte_threshold,
            dest_data_availability_overhead_gas,
            dest_gas_per_data_availability_byte,
            dest_data_availability_multiplier_bps,
            chain_family_selector,
            enforce_out_of_order,
            default_token_fee_usd_cents,
            default_token_dest_gas_overhead,
            default_tx_gas_limit,
            gas_multiplier_wei_per_eth,
            gas_price_staleness_threshold,
            network_fee_usd_cents
        };

        if (table::contains(&state.dest_chain_configs, dest_chain_selector)) {
            let dest_chain_config_ref =
                table::borrow_mut(
                    &mut state.dest_chain_configs, dest_chain_selector
                );
            *dest_chain_config_ref = dest_chain_config;
            event::emit(DestChainConfigUpdated { dest_chain_selector, dest_chain_config });
        } else {
            table::add(
                &mut state.dest_chain_configs, dest_chain_selector, dest_chain_config
            );
            event::emit(DestChainAdded { dest_chain_selector, dest_chain_config });
        }
    }

    public fun apply_premium_multiplier_wei_per_eth_updates(
        ref: &mut CCIPObjectRef,
        tokens: vector<address>,
        premium_multiplier_wei_per_eth: vector<u64>,
        ctx: &mut TxContext
    ) {
        assert!(
            ctx.sender() == state_object::get_current_owner(ref),
            E_ONLY_CALLABLE_BY_OWNER
        );
        let state = state_object::borrow_mut_with_ctx<FeeQuoterState>(ref, FEE_QUOTER_STATE_NAME, ctx);

        vector::zip_do_ref!(
            &tokens,
            &premium_multiplier_wei_per_eth,
            |token, premium_multiplier_wei_per_eth| {
                let token: address = *token;
                let premium_multiplier_wei_per_eth: u64 = *premium_multiplier_wei_per_eth;
        
                if (table::contains(&state.premium_multiplier_wei_per_eth, token)) {
                    let _old_value = table::remove(&mut state.premium_multiplier_wei_per_eth, token);
                };
                table::add(&mut state.premium_multiplier_wei_per_eth, token, premium_multiplier_wei_per_eth);
        
                event::emit(
                    PremiumMultiplierWeiPerEthUpdated {
                        token,
                        premium_multiplier_wei_per_eth
                    }
                );
            }
        );
    }

    public fun get_static_config(ref: &CCIPObjectRef): StaticConfig {
        let state = state_object::borrow<FeeQuoterState>(ref, FEE_QUOTER_STATE_NAME);
        StaticConfig {
            max_fee_juels_per_msg: state.max_fee_juels_per_msg,
            link_token: state.link_token,
            token_price_staleness_threshold: state.token_price_staleness_threshold
        }
    }

    public fun get_token_transfer_fee_config(
        ref: &CCIPObjectRef,
        dest_chain_selector: u64,
        token: address
    ): TokenTransferFeeConfig {
        let state = state_object::borrow<FeeQuoterState>(ref, FEE_QUOTER_STATE_NAME);
        *get_token_transfer_fee_config_internal(
            state, dest_chain_selector, token
        )
    }

    public fun get_token_price(ref: &CCIPObjectRef, token: address): TimestampedPrice {
        let state = state_object::borrow<FeeQuoterState>(ref, FEE_QUOTER_STATE_NAME);
        get_token_price_internal(state, token)
    }

    public fun get_token_prices(
        ref: &CCIPObjectRef,
        tokens: vector<address>
    ): (vector<TimestampedPrice>) {
        let state = state_object::borrow<FeeQuoterState>(ref, FEE_QUOTER_STATE_NAME);
        tokens.map_ref!(|token| get_token_price_internal(state, *token))
    }

    public fun get_dest_chain_gas_price(
        ref: &CCIPObjectRef,
        dest_chain_selector: u64
    ): TimestampedPrice {
        let state = state_object::borrow<FeeQuoterState>(ref, FEE_QUOTER_STATE_NAME);
        get_dest_chain_gas_price_internal(state, dest_chain_selector)
    }

    public fun get_token_and_gas_prices(
        ref: &CCIPObjectRef, clock: &clock::Clock, token: address, dest_chain_selector: u64
    ): (u256, u256) {
        let state = state_object::borrow<FeeQuoterState>(ref, FEE_QUOTER_STATE_NAME);
        let dest_chain_config = get_dest_chain_config_internal(
            state, dest_chain_selector
        );
        assert!(
            dest_chain_config.is_enabled,
            E_DEST_CHAIN_NOT_ENABLED
        );
        let token_price = get_token_price_internal(state, token);
        let gas_price_value =
            get_validated_gas_price_internal(
                state, clock, dest_chain_config, dest_chain_selector
            );
        (token_price.value, gas_price_value)
    }

    public fun get_dest_chain_config(
        ref: &CCIPObjectRef, dest_chain_selector: u64
    ): DestChainConfig {
        let state = state_object::borrow<FeeQuoterState>(ref, FEE_QUOTER_STATE_NAME);
        *get_dest_chain_config_internal(state, dest_chain_selector)
    }

    fun get_dest_chain_config_internal(
        state: &FeeQuoterState, dest_chain_selector: u64
    ): &DestChainConfig {
        assert!(
            table::contains(&state.dest_chain_configs, dest_chain_selector),
            E_UNKNOWN_DEST_CHAIN_SELECTOR
        );
        table::borrow(&state.dest_chain_configs, dest_chain_selector)
    }

    fun get_dest_chain_gas_price_internal(
        state: &FeeQuoterState, dest_chain_selector: u64
    ): TimestampedPrice {
        assert!(
            table::contains(
                &state.usd_per_unit_gas_by_dest_chain, dest_chain_selector
            ),
            E_UNKNOWN_DEST_CHAIN_SELECTOR
        );
        *table::borrow(
            &state.usd_per_unit_gas_by_dest_chain, dest_chain_selector
        )
    }

    fun get_validated_gas_price_internal(
        state: &FeeQuoterState,
        clock: &clock::Clock,
        dest_chain_config: &DestChainConfig,
        dest_chain_selector: u64
    ): u256 {
        let gas_price = get_dest_chain_gas_price_internal(state, dest_chain_selector);
        if (dest_chain_config.gas_price_staleness_threshold > 0) {
            let time_passed_secs = clock::timestamp_ms(clock) / 1000 - gas_price.timestamp;
            assert!(
                time_passed_secs <= (dest_chain_config.gas_price_staleness_threshold as u64),
                E_STALE_GAS_PRICE
            );
        };
        gas_price.value
    }

    fun get_token_price_internal(
        state: &FeeQuoterState, token: address
    ): TimestampedPrice {
        assert!(
            table::contains(&state.usd_per_token, token),
            E_UNKNOWN_TOKEN
        );
        *table::borrow(&state.usd_per_token, token)
    }

    fun convert_token_amount_internal(
        state: &FeeQuoterState,
        from_token: address,
        from_token_amount: u64,
        to_token: address
    ): u64 {
        let from_token_price = get_token_price_internal(state, from_token);
        let to_token_price = get_token_price_internal(state, to_token);

        let to_token_amount =
            ((from_token_amount as u256) * from_token_price.value) / to_token_price.value;
        assert!(
            to_token_amount <= MAX_U64,
            E_TO_TOKEN_AMOUNT_TOO_LARGE
        );
        to_token_amount as u64
    }

    // TODO: revisit the permission control
    public(package) fun update_prices(
        ref: &mut CCIPObjectRef,
        clock: &clock::Clock,
        source_tokens: vector<address>,
        source_usd_per_token: vector<u256>,
        gas_dest_chain_selectors: vector<u64>,
        gas_usd_per_unit_gas: vector<u256>,
        ctx: &TxContext
    ) {
        assert!(
            source_tokens.length() == source_usd_per_token.length(),
            E_TOKEN_UPDATE_MISMATCH
        );
        assert!(
            gas_dest_chain_selectors.length() == gas_usd_per_unit_gas.length(),
            E_GAS_UPDATE_MISMATCH
        );

        let state = state_object::borrow_mut_with_ctx<FeeQuoterState>(ref, FEE_QUOTER_STATE_NAME, ctx);
        let timestamp = clock.timestamp_ms() / 1000;

        vector::zip_do_ref!(
            &source_tokens,
            &source_usd_per_token,
            |token, usd_per_token| {
                let timestamped_price = TimestampedPrice {
                    value: *usd_per_token,
                    timestamp
                };

                if (table::contains(&state.usd_per_token, *token)) {
                    let _old_value = table::remove(&mut state.usd_per_token, *token);
                };
                table::add(&mut state.usd_per_token, *token, timestamped_price);

                event::emit(
                    UsdPerTokenUpdated {
                        token: *token,
                        usd_per_token: *usd_per_token,
                        timestamp
                    }
                );
            }
        );

        vector::zip_do_ref!(
            &gas_dest_chain_selectors,
            &gas_usd_per_unit_gas,
            |dest_chain_selector, usd_per_unit_gas| {
                let timestamped_price = TimestampedPrice {
                    value: *usd_per_unit_gas,
                    timestamp
                };

                if (table::contains(&state.usd_per_unit_gas_by_dest_chain, *dest_chain_selector)) {
                    let _old_value = table::remove(&mut state.usd_per_unit_gas_by_dest_chain, *dest_chain_selector);
                };
                table::add(&mut state.usd_per_unit_gas_by_dest_chain, *dest_chain_selector, timestamped_price);

                event::emit(
                    UsdPerUnitGasUpdated {
                        dest_chain_selector: *dest_chain_selector,
                        usd_per_unit_gas: *usd_per_unit_gas,
                        timestamp
                    }
                );
            }
        );
    }

    public(package) fun get_validated_fee(
        ref: &CCIPObjectRef,
        clock: &clock::Clock,
        dest_chain_selector: u64,
        message: &internal::Sui2AnyMessage
    ): u64 {
        let (receiver, data, fee_token, _fee_token_store, extra_args) =
            internal::get_sui2any_fields(message);

        let (local_token_addresses, local_token_amounts) =
            internal::get_sui2any_token_transfers(message);

        let state = state_object::borrow<FeeQuoterState>(ref, FEE_QUOTER_STATE_NAME);

        let dest_chain_config = get_dest_chain_config_internal(
            state, dest_chain_selector
        );
        assert!(
            dest_chain_config.is_enabled,
            E_DEST_CHAIN_NOT_ENABLED
        );

        assert!(
            vector::contains(&state.fee_tokens, &fee_token),
            E_FEE_TOKEN_NOT_SUPPORTED
        );

        let chain_family_selector = dest_chain_config.chain_family_selector;

        let data_len = data.length();
        let tokens_len = local_token_addresses.length();
        validate_message(dest_chain_config, data_len, tokens_len);

        let gas_limit =
            if (chain_family_selector == CHAIN_FAMILY_SELECTOR_EVM) {
                validate_evm_address(receiver);
                resolve_generic_gas_limit(dest_chain_config, extra_args)
            } else if (chain_family_selector == CHAIN_FAMILY_SELECTOR_APTOS) {
                validate_32byte_address(receiver, true);
                resolve_generic_gas_limit(dest_chain_config, extra_args)
            } else if (chain_family_selector == CHAIN_FAMILY_SELECTOR_SVM) {
                let require_valid_token_receiver = tokens_len > 0;
                let svm_gas_limit =
                    resolve_svm_gas_limit(
                        dest_chain_config, extra_args, require_valid_token_receiver
                    );
                let must_be_non_zero = svm_gas_limit > 0;
                validate_32byte_address(receiver, must_be_non_zero);
                svm_gas_limit
            } else {
                abort E_UNKNOWN_CHAIN_FAMILY_SELECTOR
            };

        let fee_token_price = get_token_price_internal(state, fee_token);
        let packed_gas_price =
            get_validated_gas_price_internal(
                state, clock, dest_chain_config, dest_chain_selector
            );

        // TODO: this should probably be premium_fee_usd_octa for sui?
        let (mut premium_fee_usd_wei, token_transfer_gas, token_transfer_bytes_overhead) =
        if (tokens_len > 0) {
            get_token_transfer_cost(
                state,
                dest_chain_config,
                dest_chain_selector,
                fee_token,
                fee_token_price,
                local_token_addresses,
                local_token_amounts
            )
        } else {
            ((dest_chain_config.network_fee_usd_cents as u256) * VAL_1E16, 0, 0)
        };
        let premium_multiplier =
            get_premium_multiplier_wei_per_eth_internal(state, fee_token);
        premium_fee_usd_wei = premium_fee_usd_wei * (premium_multiplier as u256); // Apply premium multiplier in wei/eth units

        let data_availability_cost_usd_36_decimals =
            if (dest_chain_config.dest_data_availability_multiplier_bps > 0) {
                // TODO: on EVM, the gas price is uint224 and the top 112 bits are used. here we're using a u256
                // and expecting that the extra top 22 bits are zeroes. update this and `gas_cost` below
                // if needed.
                let data_availability_gas_price = packed_gas_price >> GAS_PRICE_BITS;
                get_data_availability_cost(
                    dest_chain_config,
                    data_availability_gas_price,
                    data_len,
                    tokens_len,
                    token_transfer_bytes_overhead
                )
            } else { 0 };

        let call_data_length: u256 =
            (data_len as u256) * (token_transfer_bytes_overhead as u256);
        let mut dest_call_data_cost =
            call_data_length
                * (dest_chain_config.dest_gas_per_payload_byte_base as u256);
        if (call_data_length
            > (dest_chain_config.dest_gas_per_payload_byte_threshold as u256)) {
            dest_call_data_cost = (
                dest_chain_config.dest_gas_per_payload_byte_base as u256
            ) * (dest_chain_config.dest_gas_per_payload_byte_threshold as u256)
                + (
                call_data_length
                    - (dest_chain_config.dest_gas_per_payload_byte_threshold as u256)
            ) * (dest_chain_config.dest_gas_per_payload_byte_high as u256);
        };

        let total_dest_chain_gas =
            (dest_chain_config.dest_gas_overhead as u256) + (token_transfer_gas as u256)
                + dest_call_data_cost + gas_limit;

        let gas_cost = packed_gas_price & (MAX_U256 >> (255 - GAS_PRICE_BITS + 1));

        let total_cost_usd =
            (
                total_dest_chain_gas * gas_cost
                    * (dest_chain_config.gas_multiplier_wei_per_eth as u256)
            ) + premium_fee_usd_wei + data_availability_cost_usd_36_decimals;

        let fee_token_cost = total_cost_usd / fee_token_price.value;

        // we need to convert back to a u64 which is what the fungible asset module uses for amounts.
        assert!(
            fee_token_cost <= MAX_U64, E_FEE_TOKEN_COST_TOO_HIGH
        );
        fee_token_cost as u64
    }

    fun validate_message(
        dest_chain_config: &DestChainConfig, data_len: u64, tokens_len: u64
    ) {
        assert!(
            data_len <= (dest_chain_config.max_data_bytes as u64),
            E_MESSAGE_TOO_LARGE
        );
        assert!(
            tokens_len <= (dest_chain_config.max_number_of_tokens_per_msg as u64),
            E_UNSUPPORTED_NUMBER_OF_TOKENS
        );
    }

    fun validate_evm_address(encoded_address: vector<u8>) {
        let encoded_address_len = encoded_address.length();
        // TODO: why this is checking against 32 bytes?
        assert!(
            encoded_address_len == 32, E_INVALID_EVM_ADDRESS
        );

        let encoded_address_uint = eth_abi::decode_u256_value(encoded_address);

        assert!(
            encoded_address_uint >= EVM_PRECOMPILE_SPACE,
            E_INVALID_EVM_ADDRESS
        );
        assert!(
            encoded_address_uint <= MAX_U160,
            E_INVALID_EVM_ADDRESS
        );
    }

    fun validate_32byte_address(
        encoded_address: vector<u8>, must_be_non_zero: bool
    ) {
        let encoded_address_len = encoded_address.length();
        assert!(
            encoded_address_len == 32, E_INVALID_SVM_ADDRESS
        );

        if (must_be_non_zero) {
            assert!(
                encoded_address.length()== 32,
                E_INVALID_SVM_ADDRESS
            );
            let encoded_address_uint = eth_abi::decode_u256_value(encoded_address);
            assert!(
                encoded_address_uint > 0,
                E_INVALID_SVM_ADDRESS
            );
        };
    }

    fun resolve_generic_gas_limit(
        dest_chain_config: &DestChainConfig, extra_args: vector<u8>
    ): u256 {
        let extra_args_len = extra_args.length();
        if (extra_args_len == 0) {
            dest_chain_config.default_tx_gas_limit as u256
        } else {
            let (gas_limit, allow_out_of_order_execution) =
                decode_generic_extra_args(extra_args);
            assert!(
                gas_limit <= (dest_chain_config.max_per_msg_gas_limit as u256),
                E_MESSAGE_GAS_LIMIT_TOO_HIGH
            );
            assert!(
                !dest_chain_config.enforce_out_of_order || allow_out_of_order_execution,
                E_EXTRA_ARG_OUT_OF_ORDER_EXECUTION_MUST_BE_TRUE
            );
            gas_limit
        }
    }

    fun resolve_svm_gas_limit(
        dest_chain_config: &DestChainConfig,
        extra_args: vector<u8>,
        require_valid_token_receiver: bool
    ): u256 {
        let extra_args_len = extra_args.length();
        assert!(extra_args_len > 0, E_INVALID_EXTRA_ARGS_DATA);
        let (
            compute_units,
            _account_is_writable_bitmap,
            allow_out_of_order_execution,
            token_receiver,
            _accounts
        ) = decode_svm_extra_args(extra_args);
        assert!(
            !dest_chain_config.enforce_out_of_order || allow_out_of_order_execution,
            E_EXTRA_ARG_OUT_OF_ORDER_EXECUTION_MUST_BE_TRUE
        );
        assert!(
            compute_units <= dest_chain_config.max_per_msg_gas_limit,
            E_MESSAGE_COMPUTE_UNIT_LIMIT_TOO_HIGH
        );
        if (require_valid_token_receiver) {
            assert!(
                token_receiver.length() == 32,
                E_INVALID_TOKEN_RECEIVER
            );
            let token_receiver_uint = eth_abi::decode_u256_value(token_receiver);
            assert!(
                token_receiver_uint > 0,
                E_INVALID_TOKEN_RECEIVER
            );
        };
        compute_units as u256
    }

    fun get_token_transfer_cost(
        state: &FeeQuoterState,
        dest_chain_config: &DestChainConfig,
        dest_chain_selector: u64,
        fee_token: address,
        fee_token_price: TimestampedPrice,
        local_token_addresses: vector<address>,
        local_token_amounts: vector<u64>
    ): (u256, u32, u32) {
        let mut token_transfer_fee_wei: u256 = 0;
        let mut token_transfer_gas: u32 = 0;
        let mut token_transfer_bytes_overhead: u32 = 0;

        vector::zip_do_ref!(
            &local_token_addresses,
            &local_token_amounts,
            |local_token_address, local_token_amount| {
                let local_token_address: address = *local_token_address;
                let local_token_amount: u64 = *local_token_amount;

                let transfer_fee_config =
                    get_token_transfer_fee_config_internal(
                        state, dest_chain_selector, local_token_address
                    );

                if (!transfer_fee_config.is_enabled) {
                    token_transfer_fee_wei = token_transfer_fee_wei
                    + ((dest_chain_config.default_token_fee_usd_cents as u256)
                    * VAL_1E16);
                    token_transfer_gas = token_transfer_gas
                    + dest_chain_config.default_token_dest_gas_overhead;
                    token_transfer_bytes_overhead = token_transfer_bytes_overhead
                    + CCIP_LOCK_OR_BURN_V1_RET_BYTES;
                } else {
                    let mut bps_fee_usd_wei = 0;
                    if (transfer_fee_config.deci_bps > 0) {
                        let token_price =
                        if (local_token_address == fee_token) {
                            fee_token_price
                        } else {
                            get_token_price_internal(state, local_token_address)
                        };
                        let token_usd_value =
                            calc_usd_value_from_token_amount(
                                local_token_amount, token_price.value
                            );
                        bps_fee_usd_wei = (
                            token_usd_value * (transfer_fee_config.deci_bps as u256)
                        ) / VAL_1E5;
                    };

                    token_transfer_gas = token_transfer_gas + transfer_fee_config.dest_gas_overhead;
                    token_transfer_bytes_overhead = token_transfer_bytes_overhead + transfer_fee_config.dest_bytes_overhead;

                    let min_fee_usd_wei =
                        (transfer_fee_config.min_fee_usd_cents as u256) * VAL_1E16;
                    let max_fee_usd_wei =
                        (transfer_fee_config.max_fee_usd_cents as u256) * VAL_1E16;
                    let selected_fee_usd_wei =
                        if (bps_fee_usd_wei < min_fee_usd_wei) {
                            min_fee_usd_wei
                        } else if (bps_fee_usd_wei > max_fee_usd_wei) {
                            max_fee_usd_wei
                        } else {
                            bps_fee_usd_wei
                        };
                    token_transfer_fee_wei = token_transfer_fee_wei+ selected_fee_usd_wei;
                }
            }
        );

        (token_transfer_fee_wei, token_transfer_gas, token_transfer_bytes_overhead)
    }

    fun calc_usd_value_from_token_amount(
        token_amount: u64, token_price: u256
    ): u256 {
        (token_amount as u256) * (token_price as u256) / VAL_1E18
    }

    fun get_token_transfer_fee_config_internal(
        state: &FeeQuoterState, dest_chain_selector: u64, token: address
    ): &TokenTransferFeeConfig {
        assert!(
            table::contains(
                &state.token_transfer_fee_configs, dest_chain_selector
            ),
            E_UNKNOWN_DEST_CHAIN_SELECTOR
        );
        let dest_chain_fee_configs =
            table::borrow(&state.token_transfer_fee_configs, dest_chain_selector);
        assert!(
            table::contains(dest_chain_fee_configs, token),
            E_TOKEN_NOT_SUPPORTED
        );
        table::borrow(dest_chain_fee_configs, token)
    }

    public fun get_premium_multiplier_wei_per_eth(
        ref: &CCIPObjectRef,
        token: address
    ): u64 {
        let state = state_object::borrow<FeeQuoterState>(ref, FEE_QUOTER_STATE_NAME);
        get_premium_multiplier_wei_per_eth_internal(state, token)
    }

    fun get_premium_multiplier_wei_per_eth_internal(
        state: &FeeQuoterState, token: address
    ): u64 {
        assert!(
            table::contains(&state.premium_multiplier_wei_per_eth, token),
            E_UNKNOWN_TOKEN
        );
        *table::borrow(&state.premium_multiplier_wei_per_eth, token)
    }

    fun decode_generic_extra_args(extra_args: vector<u8>): (u256, bool) {
        let extra_args_len = extra_args.length();
        let args_tag = slice(&extra_args, 0, 4);
        let args_data = slice(&extra_args, 4, extra_args_len - 4);

        if (args_tag == GENERIC_EXTRA_ARGS_V2_TAG) {
            decode_generic_extra_args_v2(args_data)
        } else {
            abort E_INVALID_EXTRA_ARGS_TAG
        }
    }

    fun decode_generic_extra_args_v2(extra_args: vector<u8>): (u256, bool) {
        let mut stream = eth_abi::new_stream(extra_args);
        let gas_limit = eth_abi::decode_u256(&mut stream);
        let allow_out_of_order_execution = eth_abi::decode_bool(&mut stream);
        (gas_limit, allow_out_of_order_execution)
    }

    fun encode_generic_extra_args_v2(
        gas_limit: u256, allow_out_of_order_execution: bool
    ): vector<u8> {
        let mut extra_args = vector[];
        eth_abi::encode_selector(&mut extra_args, GENERIC_EXTRA_ARGS_V2_TAG);
        eth_abi::encode_u256(&mut extra_args, gas_limit);
        eth_abi::encode_bool(&mut extra_args, allow_out_of_order_execution);
        extra_args
    }

    fun decode_svm_extra_args(
        extra_args: vector<u8>
    ): (u32, u64, bool, vector<u8>, vector<vector<u8>>) {
        // TODO: we need extra validation here. if extra_args length is less than tag length + data length,
        let extra_args_len = extra_args.length();
        let args_tag = slice(&extra_args, 0, 4);
        assert!(
            args_tag == SVM_EXTRA_ARGS_V1_TAG,
            E_INVALID_EXTRA_ARGS_TAG
        );
        let args_data = slice(&extra_args, 4, extra_args_len - 4);
        decode_svm_extra_args_v1(args_data)
    }

    fun decode_svm_extra_args_v1(
        extra_args: vector<u8>
    ): (u32, u64, bool, vector<u8>, vector<vector<u8>>) {
        let mut stream = eth_abi::new_stream(extra_args);
        let compute_units = eth_abi::decode_u32(&mut stream);
        let account_is_writable_bitmap = eth_abi::decode_u64(&mut stream);
        let allow_out_of_order_execution = eth_abi::decode_bool(&mut stream);
        let token_receiver = eth_abi::decode_bytes32(&mut stream);
        let accounts =
            eth_abi::decode_vector!(
                &mut stream,
                |stream| { eth_abi::decode_bytes32(stream) }
            );
        (
            compute_units,
            account_is_writable_bitmap,
            allow_out_of_order_execution,
            token_receiver,
            accounts
        )
    }

    /// Returns a new vector containing `len` elements from `vec`
    /// starting at index `start`. Panics if `start + len` exceeds the vector length.
    fun slice<T: copy>(vec: &vector<T>, start: u64, len: u64): vector<T> {
        let vec_len = vec.length();
        // Ensure we have enough elements for the slice.
        assert!(start + len <= vec_len, E_OUT_OF_BOUND);
        let mut new_vec = vector::empty<T>();
        let mut i = start;
        while (i < start + len) {
            // Copy each element from the original vector into the new vector.
            new_vec.push_back(vec[i]);
            i = i + 1;
        };
        new_vec
    }

    fun get_data_availability_cost(
        dest_chain_config: &DestChainConfig,
        data_availability_gas_price: u256,
        data_len: u64,
        tokens_len: u64,
        total_transfer_bytes_overhead: u32
    ): u256 {
        let data_availability_length_bytes =
            MESSAGE_FIXED_BYTES + data_len + (tokens_len
                * MESSAGE_FIXED_BYTES_PER_TOKEN)
                + (total_transfer_bytes_overhead as u64);

        let data_availability_gas =
            ((data_availability_length_bytes as u256)
                * (dest_chain_config.dest_gas_per_data_availability_byte as u256)) + (
                dest_chain_config.dest_data_availability_overhead_gas as u256
            );

        data_availability_gas * data_availability_gas_price
            * (dest_chain_config.dest_data_availability_multiplier_bps as u256)
            * VAL_1E14
    }

    fun process_chain_family_selector(
        dest_chain_config: &DestChainConfig,
        is_message_with_token_transfers: bool,
        extra_args: vector<u8>
    ): (vector<u8>, bool) {
        let chain_family_selector = dest_chain_config.chain_family_selector;
        if (chain_family_selector == CHAIN_FAMILY_SELECTOR_EVM
            || chain_family_selector == CHAIN_FAMILY_SELECTOR_APTOS) {
            let (gas_limit, allow_out_of_order_execution) =
                decode_generic_extra_args(extra_args);
            let extra_args_v2 =
                encode_generic_extra_args_v2(gas_limit, allow_out_of_order_execution);
            (extra_args_v2, allow_out_of_order_execution)
        } else if (chain_family_selector == CHAIN_FAMILY_SELECTOR_SVM) {
            let (
                _compute_units,
                _account_is_writable_bitmap,
                allow_out_of_order_execution,
                token_receiver,
                _accounts
            ) = decode_svm_extra_args(extra_args);
            if (is_message_with_token_transfers) {
                assert!(
                    token_receiver.length() == 32,
                    E_INVALID_TOKEN_RECEIVER
                );
                let token_receiver_uint = eth_abi::decode_u256_value(token_receiver);
                assert!(
                    token_receiver_uint > 0,
                    E_INVALID_TOKEN_RECEIVER
                );
            };
            (extra_args, allow_out_of_order_execution)
        } else {
            // TODO: add support for Aptos?
            abort E_UNKNOWN_CHAIN_FAMILY_SELECTOR
        }
    }

    fun process_pool_return_data(
        dest_chain_config: &DestChainConfig,
        token_transfer_fee_config: &TokenTransferFeeConfig,
        dest_token_addresses: vector<vector<u8>>,
        dest_pool_datas: vector<vector<u8>>
    ): vector<vector<u8>> {
        let chain_family_selector = dest_chain_config.chain_family_selector;

        let tokens_len = dest_token_addresses.length();

        let mut dest_exec_data_per_token = vector[];
        let mut i = 0;
        while (i < tokens_len) {
            let dest_token_address = dest_token_addresses[i];
            let dest_pool_data = &dest_pool_datas[i];
            let dest_pool_data_len = dest_pool_data.length();
            if (dest_pool_data_len > (CCIP_LOCK_OR_BURN_V1_RET_BYTES as u64)) {
                assert!(
                    dest_pool_data_len <= (token_transfer_fee_config.dest_bytes_overhead as u64),
                    E_SOURCE_TOKEN_DATA_TOO_LARGE
                );
            };

            if (chain_family_selector == CHAIN_FAMILY_SELECTOR_EVM) {
                validate_evm_address(dest_token_address);
            } else if (chain_family_selector == CHAIN_FAMILY_SELECTOR_SVM
                || chain_family_selector == CHAIN_FAMILY_SELECTOR_APTOS) {
                validate_32byte_address(dest_token_address, /* must_be_non_zero= */ true);
            };

            let dest_gas_amount =
                if (token_transfer_fee_config.is_enabled) {
                    token_transfer_fee_config.dest_gas_overhead
                } else {
                    dest_chain_config.default_token_dest_gas_overhead
                };

            let mut dest_exec_data = vector[];
            eth_abi::encode_u32(&mut dest_exec_data, dest_gas_amount);
            dest_exec_data_per_token.push_back(dest_exec_data);

            i = i + 1;
        };

        dest_exec_data_per_token
    }

    public(package) fun process_message_args(
        ref: &CCIPObjectRef,
        dest_chain_selector: u64,
        fee_token: address,
        fee_token_amount: u64,
        extra_args: vector<u8>,
        dest_token_addresses: vector<vector<u8>>,
        dest_pool_datas: vector<vector<u8>>
    ): (u64, bool, vector<u8>, vector<vector<u8>>) {
        let state = state_object::borrow<FeeQuoterState>(ref, FEE_QUOTER_STATE_NAME);
        let msg_fee_juels =
            if (fee_token == state.link_token) {
                fee_token_amount
            } else {
                convert_token_amount_internal(
                    state,
                    fee_token,
                    fee_token_amount,
                    state.link_token
                )
            };

        assert!(
            msg_fee_juels <= state.max_fee_juels_per_msg,
            E_MESSAGE_FEE_TOO_HIGH
        );

        let dest_chain_config = get_dest_chain_config_internal(
            state, dest_chain_selector
        );
        let token_transfer_fee_config =
            get_token_transfer_fee_config_internal(state, dest_chain_selector, fee_token);

        let (converted_extra_args, is_out_of_order_execution) =
            process_chain_family_selector(
                dest_chain_config,
                !vector::is_empty(&dest_token_addresses),
                extra_args
            );

        let dest_exec_data_per_token =
            process_pool_return_data(
                dest_chain_config,
                token_transfer_fee_config,
                dest_token_addresses,
                dest_pool_datas
            );

        (
            msg_fee_juels,
            is_out_of_order_execution,
            converted_extra_args,
            dest_exec_data_per_token
        )
    }

    #[test_only]
    public fun get_fee_tokens(ref: &CCIPObjectRef ): vector<address> {
        let state = state_object::borrow<FeeQuoterState>(ref, FEE_QUOTER_STATE_NAME);
        state.fee_tokens
    }
}

#[test_only]
module ccip::fee_quoter_test {
    use std::bcs;
    use sui::test_scenario::{Self, Scenario};
    use sui::clock;

    use ccip::internal;
    use ccip::fee_quoter;
    use ccip::state_object::{Self, CCIPObjectRef};

    const CHAIN_FAMILY_SELECTOR_EVM: vector<u8> = x"2812d52c";
    const CHAIN_FAMILY_SELECTOR_SVM: vector<u8> = x"1e10bdc4";
    const FEE_QUOTER_STATE_NAME: vector<u8> = b"FeeQuoterState";
    const MOCK_ADDRESS_1: address = @0x1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b;
    const MOCK_ADDRESS_2: address = @0x000000000000000000000000F4030086522a5bEEa4988F8cA5B36dbC97BeE88c; // EVM token address
    const MOCK_ADDRESS_3: address = @0x8a7b6c5d4e3f2a1b0c9d8e7f6a5b4c3d2e1f0a9b8c7d6e5f4a3b2c1d0e9f8a7;
    const MOCK_ADDRESS_4: address = @0x3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d;
    const MOCK_ADDRESS_5: address = @0xd1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d5e6f7a8b9c0d1e2;

    fun set_up_test(): (Scenario, CCIPObjectRef) {
        let mut scenario = test_scenario::begin(@0x1);
        let ctx = scenario.ctx();

        let ref = state_object::create(ctx);

        (scenario, ref)
    }

    fun initialize(ref: &mut CCIPObjectRef, ctx: &mut TxContext) {
        fee_quoter::initialize(
            ref,
            2000,
            MOCK_ADDRESS_1,
            1000,
            vector[
                MOCK_ADDRESS_1,
                MOCK_ADDRESS_2,
                MOCK_ADDRESS_3
            ],
            ctx
        );
    }

    fun tear_down_test(scenario: Scenario, ref: CCIPObjectRef) {
        state_object::destroy_state_object(ref);
        test_scenario::end(scenario);
    }

    #[test]
    public fun test_initialize() {
        let (mut scenario, mut ref) = set_up_test();
        let ctx = scenario.ctx();
        initialize(&mut ref, ctx);

        let _state = state_object::borrow<fee_quoter::FeeQuoterState>(&ref, FEE_QUOTER_STATE_NAME);

        let fee_tokens = fee_quoter::get_fee_tokens(&ref);
        assert!(fee_tokens == vector[
            MOCK_ADDRESS_1,
            MOCK_ADDRESS_2,
            MOCK_ADDRESS_3
        ]);

        tear_down_test(scenario, ref);
    }

    #[test]
    public fun test_apply_fee_token_updates() {
        let (mut scenario, mut ref) = set_up_test();
        let ctx = scenario.ctx();
        initialize(&mut ref, ctx);

        fee_quoter::apply_fee_token_updates(
            &mut ref,
            vector[
                MOCK_ADDRESS_1,
                MOCK_ADDRESS_2
            ],
            vector[
                MOCK_ADDRESS_4,
                MOCK_ADDRESS_5
            ],
            ctx
        );

        let fee_tokens = fee_quoter::get_fee_tokens(&ref);
        assert!(fee_tokens == vector[
            MOCK_ADDRESS_3,
            MOCK_ADDRESS_4,
            MOCK_ADDRESS_5
        ]);

        tear_down_test(scenario, ref);
    }

    #[test]
    public fun test_apply_token_transfer_fee_config_updates() {
        let (mut scenario, mut ref) = set_up_test();
        let ctx = scenario.ctx();
        initialize(&mut ref, ctx);

        fee_quoter::apply_token_transfer_fee_config_updates(
            &mut ref,
            10,
            vector[MOCK_ADDRESS_1, MOCK_ADDRESS_2],
            vector[100, 200],
            vector[3000, 4000],
            vector[500, 600],
            vector[700, 800],
            vector[900, 1000],
            vector[true, false],
            vector[],
            ctx
        );

        // a successful get means the config is created.
        // we can verify the content of config but that requires an additional
        // function to expose fields within the config due to the fact that this
        // test is outside the module.
        let _config1 = fee_quoter::get_token_transfer_fee_config(&ref, 10, MOCK_ADDRESS_1);
        let _config2 = fee_quoter::get_token_transfer_fee_config(&ref, 10, MOCK_ADDRESS_2);

        tear_down_test(scenario, ref);
    }

    #[test]
    #[expected_failure(abort_code = fee_quoter::E_TOKEN_TRANSFER_FEE_CONFIG_MISMATCH)]
    public fun test_apply_token_transfer_fee_config_updates_config_mismatch() {
        let (mut scenario, mut ref) = set_up_test();
        let ctx = scenario.ctx();
        initialize(&mut ref, ctx);

        fee_quoter::apply_token_transfer_fee_config_updates(
            &mut ref,
            10,
            vector[MOCK_ADDRESS_1, MOCK_ADDRESS_2],
            vector[100, 200],
            vector[3000], // only one value
            vector[500, 600],
            vector[700, 800],
            vector[900, 1000],
            vector[true, false],
            vector[],
            ctx
        );

        tear_down_test(scenario, ref);
    }

    #[test]
    #[expected_failure(abort_code = fee_quoter::E_TOKEN_NOT_SUPPORTED)]
    public fun test_apply_token_transfer_fee_config_updates_remove_token() {
        let (mut scenario, mut ref) = set_up_test();
        let ctx = scenario.ctx();
        initialize(&mut ref, ctx);

        fee_quoter::apply_token_transfer_fee_config_updates(
            &mut ref,
            10,
            vector[MOCK_ADDRESS_1, MOCK_ADDRESS_2],
            vector[100, 200],
            vector[3000, 4000],
            vector[500, 600],
            vector[700, 800],
            vector[900, 1000],
            vector[true, false],
            vector[],
            ctx
        );

        fee_quoter::apply_token_transfer_fee_config_updates(
            &mut ref,
            10, // dest_chain_selector
            vector[], // source_tokens
            vector[], // min_fee_usd_cents
            vector[], // max_fee_usd_cents
            vector[], // dest_gas_overhead
            vector[], // dest_bytes_overhead
            vector[], // deci_bps
            vector[], // is_enabled
            vector[MOCK_ADDRESS_1], // remove MOCK_ADDRESS_1
            ctx
        );

        fee_quoter::get_token_transfer_fee_config(&ref, 10, MOCK_ADDRESS_1);

        tear_down_test(scenario, ref);
    }

    #[test]
    public fun test_apply_premium_multiplier_wei_per_eth_updates() {
        let (mut scenario, mut ref) = set_up_test();
        let ctx = scenario.ctx();
        initialize(&mut ref, ctx);

        fee_quoter::apply_premium_multiplier_wei_per_eth_updates(
            &mut ref,
            vector[MOCK_ADDRESS_1, MOCK_ADDRESS_2], // source_tokens
            vector[1000, 2000], // premium_multiplier_wei_per_eth
            ctx
        );

        assert!(fee_quoter::get_premium_multiplier_wei_per_eth(&ref, MOCK_ADDRESS_1) == 1000);
        assert!(fee_quoter::get_premium_multiplier_wei_per_eth(&ref, MOCK_ADDRESS_2) == 2000);

        tear_down_test(scenario, ref);
    }

    #[test]
    public fun test_update_prices() {
        let (mut scenario, mut ref) = set_up_test();
        let ctx = scenario.ctx();
        initialize(&mut ref, ctx);

        let mut clock = clock::create_for_testing(ctx);
        clock::increment_for_testing(&mut clock, 20000);
        fee_quoter::update_prices(
            &mut ref,
            &clock,
            vector[MOCK_ADDRESS_1, MOCK_ADDRESS_2], // source_tokens
            vector[1000, 2000], // source_usd_per_token
            vector[100, 1000], // gas_dest_chain_selectors
            vector[3000, 4000], // gas_usd_per_unit_gas
            ctx
        );

        // prices are successfully updated if we can find the config for the dest chain selector / token address
        let _timestamp_price = fee_quoter::get_dest_chain_gas_price(&ref, 100);
        let _token_price = fee_quoter::get_token_price(&ref, MOCK_ADDRESS_1);

        clock::destroy_for_testing(clock);
        tear_down_test(scenario, ref);
    }

    #[test]
    public fun test_apply_dest_chain_config_updates() {
        let (mut scenario, mut ref) = set_up_test();
        let ctx = scenario.ctx();
        initialize(&mut ref, ctx);

        fee_quoter::apply_dest_chain_config_updates(
            &mut ref,
            100, // dest_chain_selector
            true, // is_enabled
            1000, // max_number_of_tokens_per_msg
            20000, // max_data_bytes
            5000000, // max_per_msg_gas_limit
            100000, // dest_gas_overhead
            100, // dest_gas_per_payload_byte_base
            200, // dest_gas_per_payload_byte_high
            300, // dest_gas_per_payload_byte_threshold
            400000, // dest_data_availability_overhead_gas
            500, // dest_gas_per_data_availability_byte
            600, // dest_data_availability_multiplier_bps
            CHAIN_FAMILY_SELECTOR_EVM, // chain_family_selector
            true, // enforce_out_of_order
            1000, // default_token_fee_usd_cents
            2000, // default_token_dest_gas_overhead
            3000000, // default_tx_gas_limit
            4000000, // gas_multiplier_wei_per_eth
            5000000, // gas_price_staleness_threshold
            6000000, // network_fee_usd_cents
            ctx
        );

        let _config = fee_quoter::get_dest_chain_config(&ref, 100);

        tear_down_test(scenario, ref);
    }

    #[allow(implicit_const_copy)]
    #[test]
    public fun test_process_message_args_evm() {
        let (mut scenario, mut ref) = set_up_test();
        let ctx = scenario.ctx();
        initialize(&mut ref, ctx);

        fee_quoter::apply_dest_chain_config_updates(
            &mut ref,
            100, // dest_chain_selector
            true, // is_enabled
            1000, // max_number_of_tokens_per_msg
            20000, // max_data_bytes
            5000000, // max_per_msg_gas_limit
            100000, // dest_gas_overhead
            100, // dest_gas_per_payload_byte_base
            200, // dest_gas_per_payload_byte_high
            300, // dest_gas_per_payload_byte_threshold
            400000, // dest_data_availability_overhead_gas
            500, // dest_gas_per_data_availability_byte
            600, // dest_data_availability_multiplier_bps
            CHAIN_FAMILY_SELECTOR_EVM, // chain_family_selector
            true, // enforce_out_of_order
            1000, // default_token_fee_usd_cents
            2000, // default_token_dest_gas_overhead
            3000000, // default_tx_gas_limit
            4000000, // gas_multiplier_wei_per_eth
            5000000, // gas_price_staleness_threshold
            6000000, // network_fee_usd_cents
            ctx
        );

        fee_quoter::apply_token_transfer_fee_config_updates(
            &mut ref,
            100,
            vector[MOCK_ADDRESS_1, MOCK_ADDRESS_2], // source_tokens
            vector[100, 200], // source_usd_per_token
            vector[3000, 4000], // gas_dest_chain_selectors
            vector[500, 600], // gas_usd_per_unit_gas
            vector[700, 800], // dest_gas_overhead
            vector[900, 1000], // dest_bytes_overhead
            vector[true, false], // is_enabled
            vector[], // dest_chain_selectors
            ctx
        );

        let evm_extra_args = x"181dcf10181dcf10181dcf10181dcf10181dcf10181dcf10181dcf10181dcf10181dcf100000000000000000000000000000000000000000000000000000000000000001";

        let (
            msg_fee_juels,
            is_out_of_order_execution,
            converted_extra_args,
            dest_exec_data_per_token
        ) = fee_quoter::process_message_args(
            &ref,
            100, // dest_chain_selector
            MOCK_ADDRESS_1, // fee_token
            1000, // fee_token_amount
            evm_extra_args, // extra_args
            vector[
                bcs::to_bytes(&MOCK_ADDRESS_2)
            ], // dest_token_addresses
            vector[
                bcs::to_bytes(&MOCK_ADDRESS_3)
            ] // dest_pool_datas
        );

        assert!(msg_fee_juels == 1000);
        assert!(is_out_of_order_execution == true);
        assert!(converted_extra_args == evm_extra_args);
        assert!(dest_exec_data_per_token == vector[x"00000000000000000000000000000000000000000000000000000000000002bc"]);

        tear_down_test(scenario, ref);
    }

    #[allow(implicit_const_copy)]
    #[test]
    public fun test_process_message_args_svm() {
        let (mut scenario, mut ref) = set_up_test();
        let ctx = scenario.ctx();
        initialize(&mut ref, ctx);

        fee_quoter::apply_dest_chain_config_updates(
            &mut ref,
            100, // dest_chain_selector
            true, // is_enabled
            1000, // max_number_of_tokens_per_msg
            20000, // max_data_bytes
            5000000, // max_per_msg_gas_limit
            100000, // dest_gas_overhead
            100, // dest_gas_per_payload_byte_base
            200, // dest_gas_per_payload_byte_high
            300, // dest_gas_per_payload_byte_threshold
            400000, // dest_data_availability_overhead_gas
            500, // dest_gas_per_data_availability_byte
            600, // dest_data_availability_multiplier_bps
            CHAIN_FAMILY_SELECTOR_SVM, // chain_family_selector
            true, // enforce_out_of_order
            1000, // default_token_fee_usd_cents
            2000, // default_token_dest_gas_overhead
            3000000, // default_tx_gas_limit
            4000000, // gas_multiplier_wei_per_eth
            5000000, // gas_price_staleness_threshold
            6000000, // network_fee_usd_cents
            ctx
        );

        fee_quoter::apply_token_transfer_fee_config_updates(
            &mut ref,
            100, // dest_chain_selector
            vector[MOCK_ADDRESS_1, MOCK_ADDRESS_4], // source_tokens
            vector[100, 200], // source_usd_per_token
            vector[3000, 4000], // gas_dest_chain_selectors
            vector[500, 600], // gas_usd_per_unit_gas
            vector[700, 800], // dest_gas_overhead
            vector[900, 1000], // dest_bytes_overhead
            vector[true, false], // is_enabled
            vector[], // dest_chain_selectors
            ctx
        );

        let svm_extra_args = x"1f3b3aba00000000000000000000000000000000000000000000000000000000000dcf00000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000abc00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000111";

        let (
            msg_fee_juels,
            is_out_of_order_execution,
            converted_extra_args,
            dest_exec_data_per_token
        ) = fee_quoter::process_message_args(
            &ref,
            100, // dest_chain_selector
            MOCK_ADDRESS_1, // fee_token
            1000, // fee_token_amount
            svm_extra_args, // extra_args
            vector[
                bcs::to_bytes(&MOCK_ADDRESS_4)
            ], // dest_token_addresses
            vector[
                bcs::to_bytes(&MOCK_ADDRESS_3)
            ] // dest_pool_datas
        );

        assert!(msg_fee_juels == 1000);
        assert!(is_out_of_order_execution == true);
        assert!(converted_extra_args == svm_extra_args);
        assert!(dest_exec_data_per_token == vector[x"00000000000000000000000000000000000000000000000000000000000002bc"]);

        tear_down_test(scenario, ref);
    }

    #[test]
    public fun test_get_validated_fee() {
        let (mut scenario, mut ref) = set_up_test();
        let ctx = scenario.ctx();
        initialize(&mut ref, ctx);

        let mut clock = clock::create_for_testing(ctx);
        clock::increment_for_testing(&mut clock, 20000);
        fee_quoter::update_prices(
            &mut ref,
            &clock,
            vector[MOCK_ADDRESS_1, MOCK_ADDRESS_2], // source_tokens
            vector[1000, 2000], // source_usd_per_token
            vector[100, 1000], // gas_dest_chain_selectors
            vector[3000, 4000], // gas_usd_per_unit_gas
            ctx
        );

        fee_quoter::apply_dest_chain_config_updates(
            &mut ref,
            100, // dest_chain_selector
            true, // is_enabled
            1000, // max_number_of_tokens_per_msg
            20000, // max_data_bytes
            5000000, // max_per_msg_gas_limit
            100000, // dest_gas_overhead
            1, // dest_gas_per_payload_byte_base
            4, // dest_gas_per_payload_byte_high
            5, // dest_gas_per_payload_byte_threshold
            40000, // dest_data_availability_overhead_gas
            5, // dest_gas_per_data_availability_byte
            6, // dest_data_availability_multiplier_bps
            CHAIN_FAMILY_SELECTOR_EVM, // chain_family_selector
            true, // enforce_out_of_order
            1000, // default_token_fee_usd_cents
            2000, // default_token_dest_gas_overhead
            3000000, // default_tx_gas_limit
            40, // gas_multiplier_wei_per_eth
            50000, // gas_price_staleness_threshold
            600, // network_fee_usd_cents
            ctx
        );

        fee_quoter::apply_token_transfer_fee_config_updates(
            &mut ref,
            100, // dest_chain_selector
            vector[MOCK_ADDRESS_1, MOCK_ADDRESS_2], // add_tokens
            vector[100, 200], // add_min_fee_usd_cents
            vector[300, 400], // add_max_fee_usd_cents
            vector[500, 600], // add_deci_bps
            vector[700, 800], // add_dest_gas_overhead
            vector[900, 1000], // add_dest_bytes_overhead
            vector[true, false], // add_is_enabled
            vector[], // remove_tokens
            ctx
        );

        fee_quoter::apply_premium_multiplier_wei_per_eth_updates(
            &mut ref,
            vector[MOCK_ADDRESS_1, MOCK_ADDRESS_2], // source_tokens
            vector[10000, 200000], // premium_multiplier_wei_per_eth
            ctx
        );

        let evm_extra_args = x"181dcf1000000000000000000000000000000000000000000000000000000000001dcf100000000000000000000000000000000000000000000000000000000000000001";

        let message =
            internal::new_sui2any_message(
                x"000000000000000000000000f4030086522a5beea4988f8ca5b36dbc97bee88c", // receiver
                b"456abc", // data
                vector[MOCK_ADDRESS_1], // token_addresses
                vector[100], // token_amounts
                vector[MOCK_ADDRESS_2], // token_store_addresses
                MOCK_ADDRESS_1, // fee_token
                MOCK_ADDRESS_2, // fee_token_store
                evm_extra_args // extra_args
            );

        let val = fee_quoter::get_validated_fee(&ref, &clock, 100, &message);
        assert!(val == 10000000000249100440);

        clock::destroy_for_testing(clock);
        tear_down_test(scenario, ref);
    }
}