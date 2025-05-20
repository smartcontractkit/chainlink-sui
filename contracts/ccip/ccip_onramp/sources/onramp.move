module ccip_onramp::onramp {
    use sui::clock::Clock;
    use sui::balance;
    use sui::coin::{Self, Coin, CoinMetadata};
    use sui::event;
    use sui::hash;
    use std::string::{Self, String};
    use sui::table::{Self, Table};
    use sui::bag::{Self, Bag};

    use ccip::dynamic_dispatcher as dd;
    use ccip::eth_abi;
    use ccip::fee_quoter;
    use ccip::merkle_proof;
    use ccip::nonce_manager::{Self, NonceManagerCap};
    use ccip::rmn_remote;
    use ccip::state_object::CCIPObjectRef;

    public struct OwnerCap has key, store {
        id: UID
    }

    public struct OnRampState has key, store {
        id: UID,
        chain_selector: u64,
        fee_aggregator: address,
        allowlist_admin: address,

        // dest chain selector -> config
        dest_chain_configs: Table<u64, DestChainConfig>,
        // coin metadata address -> Coin
        fee_tokens: Bag,
        nonce_manager_cap: Option<NonceManagerCap>,
    }

    public struct OnRampStatePointer has key, store {
        id: UID,
        on_ramp_state_id: address,
        owner_cap_id: address,
    }

    public struct DestChainConfig has store, drop {
        // on EVM, transfers can be stopped by zeroing the router address,
        // since we don't have a router address here, we add an is_enabled flag.
        // ref: https://github.com/smartcontractkit/chainlink/blob/62a9b78e1c32174ccec11f1ed487edf3b0b4e8fd/contracts/src/v0.8/ccip/onRamp/OnRamp.sol#L181
        is_enabled: bool,
        sequence_number: u64,
        allowlist_enabled: bool,
        allowed_senders: vector<address>
    }

    public struct RampMessageHeader has store, drop, copy {
        message_id: vector<u8>,
        source_chain_selector: u64,
        dest_chain_selector: u64,
        sequence_number: u64,
        nonce: u64
    }

    public struct Sui2AnyRampMessage has store, drop, copy {
        header: RampMessageHeader,
        sender: address,
        data: vector<u8>,
        receiver: vector<u8>,
        extra_args: vector<u8>,
        fee_token: address,
        fee_token_amount: u64,
        fee_value_juels: u64,
        token_amounts: vector<Sui2AnyTokenTransfer>
    }

    public struct Sui2AnyTokenTransfer has store, drop, copy {
        source_pool_address: address,
        dest_token_address: vector<u8>, // this should become destination token coin metadata address? or remove
        extra_data: vector<u8>, // random bytes provided by token pool, e.g. encoded decimals
        amount: u64,
        dest_exec_data: vector<u8> // destination gas amount
    }

    public struct StaticConfig has copy, drop {
        chain_selector: u64
    }

    public struct DynamicConfig has copy, drop {
        fee_aggregator: address,
        allowlist_admin: address
    }

    public struct ConfigSet has copy, drop {
        static_config: StaticConfig,
        dynamic_config: DynamicConfig
    }

    public struct DestChainConfigSet has copy, drop {
        dest_chain_selector: u64,
        is_enabled: bool,
        sequence_number: u64,
        allowlist_enabled: bool
    }

    public struct CCIPMessageSent has copy, drop {
        dest_chain_selector: u64,
        sequence_number: u64,
        message: Sui2AnyRampMessage
    }

    public struct AllowlistSendersAdded has copy, drop {
        dest_chain_selector: u64,
        senders: vector<address>
    }

    public struct AllowlistSendersRemoved has copy, drop {
        dest_chain_selector: u64,
        senders: vector<address>
    }

    public struct FeeTokenWithdrawn has copy, drop {
        fee_aggregator: address,
        fee_token: address,
        amount: u64
    }

    const E_DEST_CHAIN_ARGUMENT_MISMATCH: u64 = 1;
    const E_INVALID_DEST_CHAIN_SELECTOR: u64 = 2;
    const E_UNKNOWN_DEST_CHAIN_SELECTOR: u64 = 3;
    const E_DEST_CHAIN_NOT_ENABLED: u64 = 4;
    const E_SENDER_NOT_ALLOWED: u64 = 5;
    const E_ONLY_CALLABLE_BY_ALLOWLIST_ADMIN: u64 = 6;
    const E_INVALID_ALLOWLIST_REQUEST: u64 = 7;
    const E_INVALID_ALLOWLIST_ADDRESS: u64 = 8;
    const E_CURSED_BY_RMN: u64 = 9;
    const E_BAD_RMN_SIGNAL: u64 = 10;
    const E_UNEXPECTED_WITHDRAW_AMOUNT: u64 = 11;
    const E_FEE_AGGREGATOR_NOT_SET: u64 = 12;
    const E_NONCE_MANAGER_CAP_EXISTS: u64 = 13;

    public fun type_and_version(): String {
        string::utf8(b"OnRamp 1.6.0")
    }

    fun init(ctx: &mut TxContext) {
        let owner_cap = OwnerCap {
            id: object::new(ctx)
        };

        let state = OnRampState {
            id: object::new(ctx),
            chain_selector: 0,
            fee_aggregator: @0x0,
            allowlist_admin: @0x0,
            dest_chain_configs: table::new(ctx),
            fee_tokens: bag::new(ctx),
            nonce_manager_cap: option::none(),
        };

        let pointer = OnRampStatePointer {
            id: object::new(ctx),
            on_ramp_state_id: object::id_to_address(object::borrow_id(&state)),
            owner_cap_id: object::id_to_address(object::borrow_id(&owner_cap)),
        };

        transfer::share_object(state);
        transfer::transfer(owner_cap, ctx.sender());
        transfer::transfer(pointer, ctx.sender());
    }

    public fun initialize(
        state: &mut OnRampState,
        _: &OwnerCap,
        cap: NonceManagerCap,
        chain_selector: u64,
        fee_aggregator: address,
        allowlist_admin: address,
        dest_chain_selectors: vector<u64>,
        dest_chain_enabled: vector<bool>,
        dest_chain_allowlist_enabled: vector<bool>
    ) {
        state.chain_selector = chain_selector;
        assert!(
            state.nonce_manager_cap.is_none(),
            E_NONCE_MANAGER_CAP_EXISTS
        );
        state.nonce_manager_cap.fill(cap);

        set_dynamic_config_internal(state, fee_aggregator, allowlist_admin);

        apply_dest_chain_config_updates_internal(
            state,
            dest_chain_selectors,
            dest_chain_enabled,
            dest_chain_allowlist_enabled
        );
    }

    public fun is_chain_supported(state: &OnRampState, dest_chain_selector: u64): bool {
        state.dest_chain_configs.contains(dest_chain_selector)
    }

    public fun get_expected_next_sequence_number(state: &OnRampState, dest_chain_selector: u64): u64 {
        assert!(
            state.dest_chain_configs.contains(dest_chain_selector),
            E_UNKNOWN_DEST_CHAIN_SELECTOR
        );
        let dest_chain_config = &state.dest_chain_configs[dest_chain_selector];
        dest_chain_config.sequence_number + 1
    }

    // TODO: verify withdraw fee tokens
    public fun withdraw_fee_tokens<T>(
        state: &mut OnRampState,
        _: &OwnerCap,
        fee_token_metadata: &CoinMetadata<T>
    ) {
        assert!(state.fee_aggregator != @0x0, E_FEE_AGGREGATOR_NOT_SET);

        let fee_token_metadata_addr = object::id_to_address(object::borrow_id(fee_token_metadata));

        let coins: Coin<T> = bag::remove(&mut state.fee_tokens, fee_token_metadata_addr);
        let balance = balance::value(coin::balance(&coins));
        transfer::public_transfer(coins, state.fee_aggregator);
        event::emit(
            FeeTokenWithdrawn {
                fee_aggregator: state.fee_aggregator,
                fee_token: fee_token_metadata_addr,
                amount: balance
            }
        );
    }

    fun set_dynamic_config_internal(
        state: &mut OnRampState, fee_aggregator: address, allowlist_admin: address
    ) {
        state.fee_aggregator = fee_aggregator;
        state.allowlist_admin = allowlist_admin;

        let static_config = StaticConfig { chain_selector: state.chain_selector };
        let dynamic_config = DynamicConfig { fee_aggregator, allowlist_admin };

        event::emit(ConfigSet { static_config, dynamic_config });
    }

    fun apply_dest_chain_config_updates_internal(
        state: &mut OnRampState,
        dest_chain_selectors: vector<u64>,
        dest_chain_enabled: vector<bool>,
        dest_chain_allowlist_enabled: vector<bool>
    ) {
        let dest_chains_len = dest_chain_selectors.length();
        assert!(
            dest_chains_len == dest_chain_enabled.length(),
            E_DEST_CHAIN_ARGUMENT_MISMATCH
        );
        assert!(
            dest_chains_len == dest_chain_allowlist_enabled.length(),
            E_DEST_CHAIN_ARGUMENT_MISMATCH
        );

        let mut i = 0;
        while (i < dest_chains_len) {
            let dest_chain_selector = dest_chain_selectors[i];
            assert!(
                dest_chain_selector != 0,
                E_INVALID_DEST_CHAIN_SELECTOR
            );

            let is_enabled = dest_chain_enabled[i];
            let allowlist_enabled = dest_chain_allowlist_enabled[i];

            if (!state.dest_chain_configs.contains(dest_chain_selector)) {
                state.dest_chain_configs.add(
                    dest_chain_selector,
                    DestChainConfig {
                        is_enabled: false,
                        sequence_number: 0,
                        allowlist_enabled: false,
                        allowed_senders: vector[]
                    }
                );
            };

            let dest_chain_config =
                table::borrow_mut(
                    &mut state.dest_chain_configs, dest_chain_selector
                );

            dest_chain_config.is_enabled = is_enabled;
            dest_chain_config.allowlist_enabled = allowlist_enabled;

            event::emit(
                DestChainConfigSet {
                    dest_chain_selector,
                    is_enabled,
                    sequence_number: dest_chain_config.sequence_number,
                    allowlist_enabled: dest_chain_config.allowlist_enabled
                }
            );

            i = i + 1;
        };
    }

    public fun get_fee<T>(
        ref: &CCIPObjectRef,
        clock: &Clock,
        dest_chain_selector: u64,
        receiver: vector<u8>,
        data: vector<u8>,
        token_addresses: vector<address>,
        token_amounts: vector<u64>,
        fee_token: &CoinMetadata<T>,
        extra_args: vector<u8>
    ): u64 {
        get_fee_internal(
            ref,
            clock,
            dest_chain_selector,
            receiver,
            data,
            token_addresses,
            token_amounts,
            object::id_to_address(object::borrow_id(fee_token)),
            extra_args,
        )
    }

    fun get_fee_internal(
        ref: &CCIPObjectRef,
        clock: &Clock,
        dest_chain_selector: u64,
        receiver: vector<u8>,
        data: vector<u8>,
        token_addresses: vector<address>,
        token_amounts: vector<u64>,
        fee_token: address,
        extra_args: vector<u8>
    ): u64 {
        assert!(
            !rmn_remote::is_cursed_u128(ref, dest_chain_selector as u128),
            E_CURSED_BY_RMN
        );
        fee_quoter::get_validated_fee(
            ref,
            clock,
            dest_chain_selector,
            receiver,
            data,
            token_addresses,
            token_amounts,
            fee_token,
            extra_args,
        )
    }

    public fun set_dynamic_config(
        state: &mut OnRampState,
        _: &OwnerCap,
        fee_aggregator: address,
        allowlist_admin: address
    ) {
        set_dynamic_config_internal(state, fee_aggregator, allowlist_admin);
    }

    public fun apply_dest_chain_config_updates(
        state: &mut OnRampState,
        _: &OwnerCap,
        dest_chain_selectors: vector<u64>,
        dest_chain_enabled: vector<bool>,
        dest_chain_allowlist_enabled: vector<bool>
    ) {
        apply_dest_chain_config_updates_internal(
            state,
            dest_chain_selectors,
            dest_chain_enabled,
            dest_chain_allowlist_enabled
        )
    }

    public fun get_dest_chain_config(state: &OnRampState, dest_chain_selector: u64): (bool, u64, bool, vector<address>) {
        assert!(
            state.dest_chain_configs.contains(dest_chain_selector),
            E_UNKNOWN_DEST_CHAIN_SELECTOR
        );

        let dest_chain_config = &state.dest_chain_configs[dest_chain_selector];

        (
            dest_chain_config.is_enabled,
            dest_chain_config.sequence_number,
            dest_chain_config.allowlist_enabled,
            dest_chain_config.allowed_senders,
        )
    }

    public fun get_allowed_senders_list(state: &OnRampState, dest_chain_selector: u64): (bool, vector<address>) {
        assert!(
            state.dest_chain_configs.contains(dest_chain_selector),
            E_UNKNOWN_DEST_CHAIN_SELECTOR
        );

        let dest_chain_config = &state.dest_chain_configs[dest_chain_selector];

        (dest_chain_config.allowlist_enabled, dest_chain_config.allowed_senders)
    }

    public fun apply_allowlist_updates(
        state: &mut OnRampState,
        _: &OwnerCap,
        dest_chain_selectors: vector<u64>,
        dest_chain_allowlist_enabled: vector<bool>,
        dest_chain_add_allowed_senders: vector<vector<address>>,
        dest_chain_remove_allowed_senders: vector<vector<address>>
    ) {
        apply_allowlist_updates_internal(
            state,
            dest_chain_selectors,
            dest_chain_allowlist_enabled,
            dest_chain_add_allowed_senders,
            dest_chain_remove_allowed_senders
        );
    }
    
    public fun apply_allowlist_updates_by_admin(
        state: &mut OnRampState,
        dest_chain_selectors: vector<u64>,
        dest_chain_allowlist_enabled: vector<bool>,
        dest_chain_add_allowed_senders: vector<vector<address>>,
        dest_chain_remove_allowed_senders: vector<vector<address>>,
        ctx: &mut TxContext,
    ) {
        assert!(
            state.allowlist_admin == ctx.sender(),
            E_ONLY_CALLABLE_BY_ALLOWLIST_ADMIN
        );

        apply_allowlist_updates_internal(
            state,
            dest_chain_selectors,
            dest_chain_allowlist_enabled,
            dest_chain_add_allowed_senders,
            dest_chain_remove_allowed_senders
        );
    }
    
    fun apply_allowlist_updates_internal(
        state: &mut OnRampState,
        dest_chain_selectors: vector<u64>,
        dest_chain_allowlist_enabled: vector<bool>,
        dest_chain_add_allowed_senders: vector<vector<address>>,
        dest_chain_remove_allowed_senders: vector<vector<address>>,
    ) {
        let dest_chains_len = dest_chain_selectors.length();
        assert!(
            dest_chains_len == dest_chain_allowlist_enabled.length(),
            E_DEST_CHAIN_ARGUMENT_MISMATCH
        );
        assert!(
            dest_chains_len == dest_chain_add_allowed_senders.length(),
            E_DEST_CHAIN_ARGUMENT_MISMATCH
        );
        assert!(
            dest_chains_len == dest_chain_remove_allowed_senders.length(),
            E_DEST_CHAIN_ARGUMENT_MISMATCH
        );

        let mut i = 0;
        while (i < dest_chains_len) {
            let dest_chain_selector = dest_chain_selectors[i];
            assert!(
                state.dest_chain_configs.contains(dest_chain_selector),
                E_UNKNOWN_DEST_CHAIN_SELECTOR
            );

            let allowlist_enabled = dest_chain_allowlist_enabled[i];
            let add_allowed_senders = dest_chain_add_allowed_senders[i];
            let remove_allowed_senders = dest_chain_remove_allowed_senders[i];

            let dest_chain_config =
                state.dest_chain_configs.borrow_mut(dest_chain_selector);
            dest_chain_config.allowlist_enabled = allowlist_enabled;

            if (add_allowed_senders.length() > 0) {
                assert!(allowlist_enabled, E_INVALID_ALLOWLIST_REQUEST);

                vector::do_ref!(
                    &add_allowed_senders,
                    |sender_address| {
                        let sender_address: address = *sender_address;
                        assert!(sender_address != @0x0, E_INVALID_ALLOWLIST_ADDRESS);

                        let (found, _) = vector::index_of(
                            &dest_chain_config.allowed_senders, &sender_address
                        );
                        if (!found) {
                            dest_chain_config.allowed_senders.push_back(sender_address);
                        };
                    }
                );

                event::emit(
                    AllowlistSendersAdded {
                        dest_chain_selector,
                        senders: add_allowed_senders
                    }
                );
            };

            if (remove_allowed_senders.length() > 0) {
                vector::do_ref!(
                    &remove_allowed_senders,
                    |sender_address| {
                        let (found, i) = vector::index_of(
                            &dest_chain_config.allowed_senders, sender_address
                        );
                        if (found) {
                            vector::swap_remove(
                                &mut dest_chain_config.allowed_senders, i
                            );
                        }
                    }
                );

                event::emit(
                    AllowlistSendersRemoved {
                        dest_chain_selector,
                        senders: remove_allowed_senders
                    }
                );
            };
            i = i + 1;
        };
    }

    public fun get_outbound_nonce(
        ref: &CCIPObjectRef, dest_chain_selector: u64, sender: address
    ): u64 {
        nonce_manager::get_outbound_nonce(ref, dest_chain_selector, sender)
    }

    public fun get_static_config(state: &OnRampState): StaticConfig {
        StaticConfig { chain_selector: state.chain_selector }
    }

    public fun get_static_config_fields(cfg: StaticConfig): u64 {
        cfg.chain_selector
    }

    public fun get_dynamic_config(state: &OnRampState): DynamicConfig {
        DynamicConfig {
            fee_aggregator: state.fee_aggregator,
            allowlist_admin: state.allowlist_admin
        }
    }

    public fun get_dynamic_config_fields(cfg: DynamicConfig): (address, address) {
        (cfg.fee_aggregator, cfg.allowlist_admin)
    }

    fun calculate_metadata_hash(
        source_chain_selector: u64, dest_chain_selector: u64
    ): vector<u8> {
        let mut packed = vector[];
        eth_abi::encode_bytes32(
            &mut packed, hash::keccak256(&b"Sui2AnyMessageHashV1")
        );
        eth_abi::encode_u64(&mut packed, source_chain_selector);
        eth_abi::encode_u64(&mut packed, dest_chain_selector);
        eth_abi::encode_address(&mut packed, @ccip);
        hash::keccak256(&packed)
    }

    fun calculate_message_hash(
        message: &Sui2AnyRampMessage, metadata_hash: vector<u8>
    ): vector<u8> {
        let mut outer_hash = vector[];
        eth_abi::encode_bytes32(&mut outer_hash, merkle_proof::leaf_domain_separator());
        eth_abi::encode_bytes32(&mut outer_hash, metadata_hash);

        let mut inner_hash = vector[];
        eth_abi::encode_address(&mut inner_hash, message.sender);
        eth_abi::encode_u64(&mut inner_hash, message.header.sequence_number);
        eth_abi::encode_u64(&mut inner_hash, message.header.nonce);
        eth_abi::encode_address(&mut inner_hash, message.fee_token);
        eth_abi::encode_u64(&mut inner_hash, message.fee_token_amount);
        eth_abi::encode_bytes32(&mut outer_hash, hash::keccak256(&inner_hash));

        eth_abi::encode_bytes32(
            &mut outer_hash, hash::keccak256(&message.receiver)
        );
        eth_abi::encode_bytes32(&mut outer_hash, hash::keccak256(&message.data));

        let mut token_hash = vector[];
        eth_abi::encode_u256(&mut token_hash, message.token_amounts.length() as u256);
        message.token_amounts.do_ref!(
            |token_transfer| {
                let token_transfer: &Sui2AnyTokenTransfer = token_transfer;
                eth_abi::encode_address(
                    &mut token_hash, token_transfer.source_pool_address
                );
                eth_abi::encode_bytes(&mut token_hash, token_transfer.dest_token_address);
                eth_abi::encode_bytes(&mut token_hash, token_transfer.extra_data);
                eth_abi::encode_u64(&mut token_hash, token_transfer.amount);
                eth_abi::encode_bytes(&mut token_hash, token_transfer.dest_exec_data);
            }
        );
        eth_abi::encode_bytes32(&mut outer_hash, hash::keccak256(&token_hash));

        eth_abi::encode_bytes32(
            &mut outer_hash, hash::keccak256(&message.extra_args)
        );

        hash::keccak256(&outer_hash)
    }

    #[allow(lint(self_transfer))]
    public fun ccip_send<T>(
        ref: &mut CCIPObjectRef,
        state: &mut OnRampState,
        clock: &Clock,
        dest_chain_selector: u64,
        receiver: vector<u8>,
        data: vector<u8>,
        token_params: dd::TokenParams,
        fee_token_metadata: &CoinMetadata<T>,
        mut fee_token: Coin<T>,
        extra_args: vector<u8>,
        ctx: &mut TxContext
    ): vector<u8> {
        assert!(!rmn_remote::is_cursed_global(ref), E_BAD_RMN_SIGNAL);
        assert!(!rmn_remote::is_cursed_u128(ref, (dest_chain_selector as u128)), E_BAD_RMN_SIGNAL);

        let fee_token_metadata_addr = object::id_to_address(object::borrow_id(fee_token_metadata));
        let fee_token_balance = balance::value(coin::balance(&fee_token));

        // the hot potato is returned and consumed
        let params = dd::deconstruct_token_params(token_params);

        let mut token_amounts = vector[];
        let mut source_tokens = vector[];
        let mut dest_tokens = vector[];
        let mut dest_pool_datas = vector[];
        let mut token_transfers = vector[];
        let mut i = 0;
        let tokens_len = params.length();

        while (i < tokens_len) {
            let (source_pool, amount, source_token_address, dest_token_address, extra_data) = dd::get_source_token_transfer_data(params[i]);
            token_transfers.push_back(
                Sui2AnyTokenTransfer {
                    source_pool_address: source_pool,
                    amount,
                    dest_token_address,
                    extra_data: extra_data, // encoded decimals
                    dest_exec_data: vector[] // destination execution gas amount, populated later by fee quoter
                }
            );
            token_amounts.push_back(amount);
            source_tokens.push_back(source_token_address);
            dest_tokens.push_back(dest_token_address);
            dest_pool_datas.push_back(extra_data);

            i = i + 1;
        };

        let fee_token_amount =
            get_fee_internal(
                ref,
                clock,
                dest_chain_selector,
                receiver,
                data,
                source_tokens,
                token_amounts,
                fee_token_metadata_addr,
                extra_args,
            );

        if (fee_token_amount != 0) {
            assert!(
                fee_token_amount <= fee_token_balance,
                E_UNEXPECTED_WITHDRAW_AMOUNT
            );

            let refund = coin::split(&mut fee_token, fee_token_balance - fee_token_amount, ctx);
            if (state.fee_tokens.contains(fee_token_metadata_addr)) {
                let coins: &mut Coin<T> = bag::borrow_mut(&mut state.fee_tokens, fee_token_metadata_addr);
                coins.join(fee_token);
            } else {
                state.fee_tokens.add(fee_token_metadata_addr, fee_token);
            };
            transfer::public_transfer(refund, ctx.sender());
        } else {
            transfer::public_transfer(fee_token, ctx.sender());
        };

        let sender = ctx.sender();
        verify_sender(state, dest_chain_selector, sender);

        let sequence_number = get_incremented_sequence_number(state, dest_chain_selector);

        let (
            fee_value_juels,
            is_out_of_order_execution,
            converted_extra_args,
            mut dest_exec_data_per_token
        ) =
            fee_quoter::process_message_args(
                ref,
                dest_chain_selector,
                fee_token_metadata_addr,
                fee_token_amount,
                extra_args,
                source_tokens,
                dest_tokens,
                dest_pool_datas
            );

        vector::zip_do_mut!(
            &mut token_transfers,
            &mut dest_exec_data_per_token,
            |token_amount, dest_exec_data| {
                let token_amount: &mut Sui2AnyTokenTransfer = token_amount;
                token_amount.dest_exec_data = *dest_exec_data;
            }
        );

        let message = construct_message(
            ref,
            state,
            dest_chain_selector,
            is_out_of_order_execution,
            sender,
            sequence_number,
            data,
            receiver,
            converted_extra_args,
            fee_token_metadata_addr,
            fee_token_amount,
            fee_value_juels,
            token_transfers,
            ctx
        );

        event::emit(CCIPMessageSent { dest_chain_selector, sequence_number, message });

        message.header.message_id
    }

    fun verify_sender(state: &OnRampState, dest_chain_selector: u64, sender: address) {
        assert!(
            state.dest_chain_configs.contains(dest_chain_selector),
            E_UNKNOWN_DEST_CHAIN_SELECTOR
        );

        let dest_chain_config = &state.dest_chain_configs[dest_chain_selector];
        assert!(dest_chain_config.is_enabled, E_DEST_CHAIN_NOT_ENABLED);

        if (dest_chain_config.allowlist_enabled) {
            assert!(
                dest_chain_config.allowed_senders.contains(&sender),
                E_SENDER_NOT_ALLOWED
            );
        };
    }

    fun get_incremented_sequence_number(
        state: &mut OnRampState, dest_chain_selector: u64
    ): u64 {
        let dest_chain_config = state.dest_chain_configs.borrow_mut(dest_chain_selector);
        dest_chain_config.sequence_number = dest_chain_config.sequence_number + 1;

        dest_chain_config.sequence_number
    }

    fun construct_message(
        ref: &mut CCIPObjectRef,
        state: &OnRampState,
        dest_chain_selector: u64,
        is_out_of_order_execution: bool,
        sender: address,
        sequence_number: u64,
        data: vector<u8>,
        receiver: vector<u8>,
        converted_extra_args: vector<u8>,
        fee_token_metadata: address,
        fee_token_amount: u64,
        fee_value_juels: u64,
        token_transfers: vector<Sui2AnyTokenTransfer>,
        ctx: &mut TxContext
    ): Sui2AnyRampMessage {
        // calculate nonce
        let mut nonce = 0;
        if (!is_out_of_order_execution) {
            nonce = nonce_manager::get_incremented_outbound_nonce(
                ref,
                state.nonce_manager_cap.borrow(),
                dest_chain_selector,
                sender,
                ctx
            );
        };

        // create message
        let mut message = Sui2AnyRampMessage {
            header: RampMessageHeader {
                // populated on completion
                message_id: vector[],
                source_chain_selector: state.chain_selector,
                dest_chain_selector,
                sequence_number,
                nonce
            },
            sender,
            data,
            receiver,
            extra_args: converted_extra_args,
            fee_token: fee_token_metadata,
            fee_token_amount,
            fee_value_juels,
            token_amounts: token_transfers
        };

        // attach message id
        let metadata_hash = calculate_metadata_hash(state.chain_selector, dest_chain_selector);
        let message_id = calculate_message_hash(&message, metadata_hash);
        message.header.message_id = message_id;

        message
    }

    public fun get_onramp_pointer(pointer: &OnRampStatePointer): (address, address) {
        (pointer.on_ramp_state_id, pointer.owner_cap_id)
    }
}

#[test_only]
module ccip::onramp_test {
    use ccip::onramp::{Self, OnRampState};
    use ccip::state_object::{Self, CCIPObjectRef};
    use sui::test_scenario::{Self, Scenario};

    const ON_RAMP_STATE_NAME: vector<u8> = b"OnRampState";
    const DEST_CHAIN_SELECTOR_1: u64 = 1;
    const DEST_CHAIN_SELECTOR_2: u64 = 2;
    const ALLOWED_SENDER_1: address = @0x11;
    const ALLOWED_SENDER_2: address = @0x22;
    const ALLOWED_SENDER_3: address = @0x33;

    fun set_up_test(): (Scenario, CCIPObjectRef) {
        let mut scenario = test_scenario::begin(@0x1);
        let ctx = scenario.ctx();

        let ref = state_object::create(ctx);
        (scenario, ref)
    }

    fun tear_down_test(scenario: Scenario, ref: CCIPObjectRef) {
        state_object::destroy_state_object(ref);
        test_scenario::end(scenario);
    }

    fun initialize(ref: &mut CCIPObjectRef, ctx: &mut TxContext) {
        onramp::initialize(
            ref,
            123, // chain_selector
            ctx.sender(),
            vector[DEST_CHAIN_SELECTOR_1, DEST_CHAIN_SELECTOR_2], // dest_chain_selectors
            vector[true, false], // dest_chain_enabled
            vector[true, false], // dest_chain_allowlist_enabled
            ctx
        );
    }

    #[test]
    public fun test_initialize() {
        let (mut scenario, mut ref) = set_up_test();
        let ctx = scenario.ctx();
        initialize(&mut ref, ctx);

        let _state = state_object::borrow<OnRampState>(&ref, ON_RAMP_STATE_NAME);

        assert!(onramp::is_chain_supported(&ref, DEST_CHAIN_SELECTOR_1));
        assert!(onramp::is_chain_supported(&ref, DEST_CHAIN_SELECTOR_2));

        assert!(onramp::get_expected_next_sequence_number(&ref, DEST_CHAIN_SELECTOR_1) == 1);
        assert!(onramp::get_expected_next_sequence_number(&ref, DEST_CHAIN_SELECTOR_2) == 1);

        let (enabled, seq, allowlist_enabled) = onramp::get_dest_chain_config(&ref, DEST_CHAIN_SELECTOR_1);
        assert!(enabled == true);
        assert!(seq == 0);
        assert!(allowlist_enabled == true);

        let (enabled, seq, allowlist_enabled) = onramp::get_dest_chain_config(&ref, DEST_CHAIN_SELECTOR_2);
        assert!(enabled == false);
        assert!(seq == 0);
        assert!(allowlist_enabled == false);

        tear_down_test(scenario, ref);
    }

    #[test]
    public fun test_apply_allowlist_updates() {
        let (mut scenario, mut ref) = set_up_test();
        let ctx = scenario.ctx();
        initialize(&mut ref, ctx);

        onramp::apply_allowlist_updates(
            &mut ref,
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
            ctx
        );

        let (allowlist_enabled, allowed_senders) = onramp::get_allowed_senders_list(&ref, DEST_CHAIN_SELECTOR_1);
        assert!(allowlist_enabled == true);
        assert!(allowed_senders == vector[ALLOWED_SENDER_1, ALLOWED_SENDER_2]);
        let (allowlist_enabled, allowed_senders) = onramp::get_allowed_senders_list(&ref, DEST_CHAIN_SELECTOR_2);
        assert!(allowlist_enabled == true);
        assert!(allowed_senders == vector[ALLOWED_SENDER_3]);

        onramp::apply_allowlist_updates(
            &mut ref,
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
            ctx
        );

        let (allowlist_enabled, allowed_senders) = onramp::get_allowed_senders_list(&ref, DEST_CHAIN_SELECTOR_1);
        assert!(allowlist_enabled == true);
        assert!(allowed_senders == vector[ALLOWED_SENDER_1]);
        let (allowlist_enabled, allowed_senders) = onramp::get_allowed_senders_list(&ref, DEST_CHAIN_SELECTOR_2);
        assert!(allowlist_enabled == false);
        assert!(allowed_senders == vector[ALLOWED_SENDER_3]);

        tear_down_test(scenario, ref);
    }
}