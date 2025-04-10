module ccip::onramp {
    use sui::clock::Clock;
    use sui::event;
    use sui::hash;
    use std::string::{Self, String};
    use sui::table::{Self, Table};

    use ccip::eth_abi;
    use ccip::fee_quoter;
    use ccip::internal;
    use ccip::merkle_proof;
    use ccip::nonce_manager;
    use ccip::rmn_remote;
    use ccip::state_object::{Self, CCIPObjectRef};

    public struct OnRampState has key, store {
        id: UID,
        chain_selector: u64,
        allowlist_admin: address,

        // dest chain selector -> config
        dest_chain_configs: Table<u64, DestChainConfig>
    }

    public struct DestChainConfig has store, drop {
        // on EVM, transfers can be stopped by zeroing the router address,
        // since we don't have a router address here, we add an is_enabled flag.
        // ref: https://github.com/smartcontractkit/chainlink/blob/62a9b78e1c32174ccec11f1ed487edf3b0b4e8fd/contracts/src/v0.8/ccip/onRamp/OnRamp.sol#L181
        is_enabled: bool,
        sequence_number: u64,
        allowlist_enabled: bool,
        // TODO: should we use a Table/SmartTable here?
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
        dest_token_address: vector<u8>,
        extra_data: vector<u8>,
        amount: u64,
        dest_exec_data: vector<u8>
    }

    public struct StaticConfig has store, drop {
        chain_selector: u64
    }

    public struct DynamicConfig has store, drop {
        allowlist_admin: address
    }

    public struct ConfigSet has copy, drop {
        chain_selector: u64,
        allowlist_admin: address
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

    const ON_RAMP_STATE_NAME: vector<u8> = b"OnRampState";
    const E_ALREADY_INITIALIZED: u64 = 1;
    const E_DEST_CHAIN_ARGUMENT_MISMATCH: u64 = 2;
    const E_INVALID_DEST_CHAIN_SELECTOR: u64 = 3;
    const E_UNKNOWN_DEST_CHAIN_SELECTOR: u64 = 4;
    const E_DEST_CHAIN_NOT_ENABLED: u64 = 5;
    const E_SENDER_NOT_ALLOWED: u64 = 6;
    // const E_ONLY_CALLABLE_BY_OWNER_OR_ALLOWLIST_ADMIN: u64 = 7;
    const E_INVALID_ALLOWLIST_REQUEST: u64 = 8;
    const E_INVALID_ALLOWLIST_ADDRESS: u64 = 9;
    // const E_UNSUPPORTED_TOKEN: u64 = 10;
    // const E_INVALID_FEE_TOKEN: u64 = 11;
    const E_CURSED_BY_RMN: u64 = 12;
    const E_BAD_RMN_SIGNAL: u64 = 13;
    // const E_INVALID_TOKEN: u64 = 14;
    // const E_INVALID_TOKEN_STORE: u64 = 15;
    // const E_UNEXPECTED_WITHDRAW_AMOUNT: u64 = 16;
    // const E_UNEXPECTED_FUNGIBLE_ASSET: u64 = 17;
    // const E_UNKNOWN_FUNCTION: u64 = 18;
    const E_ONLY_CALLABLE_BY_OWNER: u64 = 19;

    public fun type_and_version(): String {
        string::utf8(b"OnRamp 1.6.0")
    }

    // fun init_module(publisher: &signer) {
    //     if (@mcms_register_entrypoints != @0x0) {
    //         mcms_registry::register_entrypoint(
    //             publisher, string::utf8(b"onramp"), McmsCallback {}
    //         );
    //     };
    // }

    public fun initialize(
        ref: &mut CCIPObjectRef,
        chain_selector: u64,
        allowlist_admin: address,
        dest_chain_selectors: vector<u64>,
        dest_chain_enabled: vector<bool>,
        dest_chain_allowlist_enabled: vector<bool>,
        ctx: &mut TxContext
    ) {
        assert!(
            ctx.sender() == state_object::get_current_owner(ref),
            E_ONLY_CALLABLE_BY_OWNER
        );
        assert!(
            !state_object::contains(ref, ON_RAMP_STATE_NAME),
            E_ALREADY_INITIALIZED
        );

        let mut state = OnRampState {
            id: object::new(ctx),
            chain_selector,
            allowlist_admin: @0x0,
            dest_chain_configs: table::new<u64, DestChainConfig>(ctx)
        };

        set_dynamic_config_internal(&mut state, allowlist_admin);

        apply_dest_chain_config_updates_internal(
            &mut state,
            dest_chain_selectors,
            dest_chain_enabled,
            dest_chain_allowlist_enabled
        );

        state_object::add(ref, ON_RAMP_STATE_NAME, state, ctx);
    }

    public fun is_chain_supported(ref: &CCIPObjectRef, dest_chain_selector: u64): bool {
        let state = state_object::borrow<OnRampState>(ref, ON_RAMP_STATE_NAME);

        state.dest_chain_configs.contains(dest_chain_selector)
    }

    public fun get_expected_next_sequence_number(ref: &CCIPObjectRef, dest_chain_selector: u64): u64 {
        let state = state_object::borrow<OnRampState>(ref, ON_RAMP_STATE_NAME);
        assert!(
            state.dest_chain_configs.contains(dest_chain_selector),
            E_UNKNOWN_DEST_CHAIN_SELECTOR
        );
        let dest_chain_config = &state.dest_chain_configs[dest_chain_selector];
        dest_chain_config.sequence_number + 1
    }

    fun set_dynamic_config_internal(
        state: &mut OnRampState, allowlist_admin: address
    ) {
        state.allowlist_admin = allowlist_admin;

        event::emit(ConfigSet { chain_selector: state.chain_selector, allowlist_admin });
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

            if (!table::contains(&state.dest_chain_configs, dest_chain_selector)) {
                table::add(
                    &mut state.dest_chain_configs,
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

    public fun get_fee(
        ref: &CCIPObjectRef,
        clock: &Clock,
        dest_chain_selector: u64,
        receiver: vector<u8>,
        data: vector<u8>,
        token_addresses: vector<address>,
        token_amounts: vector<u64>,
        token_store_addresses: vector<address>,
        fee_token: address,
        fee_token_store: address,
        extra_args: vector<u8>
    ): u64 {
        let message =
            internal::new_sui2any_message(
                receiver,
                data,
                token_addresses,
                token_amounts,
                token_store_addresses,
                fee_token,
                fee_token_store,
                extra_args
            );
        get_fee_internal(ref, clock, dest_chain_selector, &message)
    }

    fun get_fee_internal(
        ref: &CCIPObjectRef,
        clock: &Clock,
        dest_chain_selector: u64,
        message: &internal::Sui2AnyMessage
    ): u64 {
        assert!(
            !rmn_remote::is_cursed_u128(ref, dest_chain_selector as u128),
            E_CURSED_BY_RMN
        );
        fee_quoter::get_validated_fee(ref, clock, dest_chain_selector, message)
    }

    public fun set_dynamic_config(
        ref: &mut CCIPObjectRef,
        allowlist_admin: address,
        ctx: &mut TxContext
    ) {
        assert!(
            ctx.sender() == state_object::get_current_owner(ref),
            E_ONLY_CALLABLE_BY_OWNER
        );
        let state = state_object::borrow_mut_with_ctx<OnRampState>(ref, ON_RAMP_STATE_NAME, ctx);

        set_dynamic_config_internal(state, allowlist_admin);
    }

    public fun apply_dest_chain_config_updates(
        ref: &mut CCIPObjectRef,
        dest_chain_selectors: vector<u64>,
        dest_chain_enabled: vector<bool>,
        dest_chain_allowlist_enabled: vector<bool>,
        ctx: &mut TxContext
    ) {
        assert!(
            ctx.sender() == state_object::get_current_owner(ref),
            E_ONLY_CALLABLE_BY_OWNER
        );
        let state = state_object::borrow_mut_with_ctx<OnRampState>(ref, ON_RAMP_STATE_NAME, ctx);

        apply_dest_chain_config_updates_internal(
            state,
            dest_chain_selectors,
            dest_chain_enabled,
            dest_chain_allowlist_enabled
        )
    }

    public fun get_dest_chain_config(ref: &CCIPObjectRef, dest_chain_selector: u64): (bool, u64, bool) {
        let state = state_object::borrow<OnRampState>(ref, ON_RAMP_STATE_NAME);

        assert!(
            state.dest_chain_configs.contains(dest_chain_selector),
            E_UNKNOWN_DEST_CHAIN_SELECTOR
        );

        let dest_chain_config = &state.dest_chain_configs[dest_chain_selector];

        (
            dest_chain_config.is_enabled,
            dest_chain_config.sequence_number,
            dest_chain_config.allowlist_enabled
        )
    }

    public fun get_allowed_senders_list(ref: &CCIPObjectRef, dest_chain_selector: u64): (bool, vector<address>) {
        let state = state_object::borrow<OnRampState>(ref, ON_RAMP_STATE_NAME);

        assert!(
            state.dest_chain_configs.contains(dest_chain_selector),
            E_UNKNOWN_DEST_CHAIN_SELECTOR
        );

        let dest_chain_config = &state.dest_chain_configs[dest_chain_selector];

        (dest_chain_config.allowlist_enabled, dest_chain_config.allowed_senders)
    }

    // TODO: verify this:
    // in aptos, this function can be called by either the owner or the allowlist admin
    // but in current implementation, only the owner can call this function
    public fun apply_allowlist_updates(
        ref: &mut CCIPObjectRef,
        dest_chain_selectors: vector<u64>,
        dest_chain_allowlist_enabled: vector<bool>,
        dest_chain_add_allowed_senders: vector<vector<address>>,
        dest_chain_remove_allowed_senders: vector<vector<address>>,
        ctx: &mut TxContext
    ) {
        let state = state_object::borrow_mut_with_ctx<OnRampState>(ref, ON_RAMP_STATE_NAME, ctx);

        // assert!(
        //     signer::address_of(caller) == auth::owner()
        //         || signer::address_of(caller) == state.allowlist_admin,
        //     E_ONLY_CALLABLE_BY_OWNER_OR_ALLOWLIST_ADMIN
        // );

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

    public fun get_static_config(ref: &CCIPObjectRef): StaticConfig {
        let state = state_object::borrow<OnRampState>(ref, ON_RAMP_STATE_NAME);
        StaticConfig { chain_selector: state.chain_selector }
    }

    public fun get_dynamic_config(ref: &CCIPObjectRef): DynamicConfig {
        let state = state_object::borrow<OnRampState>(ref, ON_RAMP_STATE_NAME);
        DynamicConfig { allowlist_admin: state.allowlist_admin }
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

    // TODO: this function needs to be permissionless but it also updates some states
    public fun ccip_send(
        ref: &mut CCIPObjectRef,
        clock: &Clock,
        dest_chain_selector: u64,
        receiver: vector<u8>,
        data: vector<u8>,
        token_addresses: vector<address>,
        token_amounts: vector<u64>,
        token_store_addresses: vector<address>,
        fee_token: address,
        fee_token_store: address,
        extra_args: vector<u8>,
        ctx: &mut TxContext
    ): vector<u8> {
        assert!(!rmn_remote::is_cursed_global(ref), E_BAD_RMN_SIGNAL);

        let message =
            internal::new_sui2any_message(
                receiver,
                data,
                token_addresses,
                token_amounts,
                token_store_addresses,
                fee_token,
                fee_token_store,
                extra_args
            );

        let fee_token_amount = get_fee_internal(ref, clock, dest_chain_selector, &message);

        if (fee_token_amount != 0) {
            // // deposit the fee in the state object's primary fungible store.
            // let fa_metadata = resolve_fungible_asset(fee_token);
            // let resolved_store =
            //     resolve_fungible_store(
            //         signer::address_of(caller), fa_metadata, fee_token_store
            //     );
            //
            // let fa =
            //     dispatchable_fungible_asset::withdraw(
            //         caller, resolved_store, fee_token_amount
            //     );
            // // validate the withdrawn asset since we're potentially calling dispatchable fungible asset functions.
            // assert!(
            //     fa_metadata == fungible_asset::metadata_from_asset(&fa),
            //     E_UNEXPECTED_FUNGIBLE_ASSET
            // );
            // assert!(
            //     fee_token_amount == fungible_asset::amount(&fa),
            //     E_UNEXPECTED_WITHDRAW_AMOUNT
            // );
            //
            // primary_fungible_store::deposit(state_object::object_address(), fa);
        };

        let sender = ctx.sender();
        verify_sender(ref, dest_chain_selector, sender);

        let mut dest_token_addresses = vector[];
        let mut dest_pool_datas = vector[];

        let tokens_len = token_addresses.length();
        let mut token_transfers = vector[];
        let mut i = 0;
        while (i < tokens_len) {
            let _token = token_addresses[i];
            let amount = token_amounts[i];
            let _token_store = token_store_addresses[i];

            // let fa_metadata = resolve_fungible_asset(token);
            // let resolved_store = resolve_fungible_store(sender, fa_metadata, token_store);
            //
            // let fa = dispatchable_fungible_asset::withdraw(
            //     caller, resolved_store, amount
            // );

            // // validate the withdrawn asset since we're potentially calling dispatchable fungible asset functions.
            // assert!(
            //     fa_metadata == fungible_asset::metadata_from_asset(&fa),
            //     E_UNEXPECTED_FUNGIBLE_ASSET
            // );
            // assert!(
            //     amount == fungible_asset::amount(&fa),
            //     E_UNEXPECTED_WITHDRAW_AMOUNT
            // );

            // let token_pool_address = token_admin_registry::get_pool(token);
            // assert!(token_pool_address != @0x0, E_UNSUPPORTED_TOKEN);
            //
            // let (dest_token_address, dest_pool_data) =
            //     token_admin_dispatcher::dispatch_lock_or_burn(
            //         token_pool_address,
            //         fa,
            //         sender,
            //         dest_chain_selector,
            //         receiver
            //     );

            let token_pool_address = @0x0;
            let dest_token_address = vector[];
            let dest_pool_data = vector[];

            dest_token_addresses.push_back(dest_token_address);
            dest_pool_datas.push_back(dest_pool_data);

            token_transfers.push_back(
                Sui2AnyTokenTransfer {
                    source_pool_address: token_pool_address,
                    dest_token_address,
                    extra_data: dest_pool_data,
                    amount,
                    dest_exec_data: vector[]
                }
            );
            
            i = i + 1;
        };

        let sequence_number = get_incremented_sequence_number(ref, dest_chain_selector);

        let (
            fee_value_juels,
            is_out_of_order_execution,
            converted_extra_args,
            mut dest_exec_data_per_token
        ) =
            fee_quoter::process_message_args(
                ref,
                dest_chain_selector,
                fee_token,
                fee_token_amount,
                extra_args,
                dest_token_addresses,
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
            dest_chain_selector,
            is_out_of_order_execution,
            sender,
            sequence_number,
            data,
            receiver,
            converted_extra_args,
            fee_token,
            fee_token_amount,
            fee_value_juels,
            token_transfers,
            ctx
        );

        event::emit(CCIPMessageSent { dest_chain_selector, sequence_number, message });

        message.header.message_id
    }

    fun verify_sender(ref: &CCIPObjectRef, dest_chain_selector: u64, sender: address) {
        let state = state_object::borrow<OnRampState>(ref, ON_RAMP_STATE_NAME);

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
        ref: &mut CCIPObjectRef, dest_chain_selector: u64
    ): u64 {
        let state = state_object::borrow_mut_from_user<OnRampState>(ref, ON_RAMP_STATE_NAME);
        let dest_chain_config = state.dest_chain_configs.borrow_mut(dest_chain_selector);
        dest_chain_config.sequence_number = dest_chain_config.sequence_number + 1;

        dest_chain_config.sequence_number
    }

    fun construct_message(
        ref: &mut CCIPObjectRef,
        dest_chain_selector: u64,
        is_out_of_order_execution: bool,
        sender: address,
        sequence_number: u64,
        data: vector<u8>,
        receiver: vector<u8>,
        converted_extra_args: vector<u8>,
        fee_token: address,
        fee_token_amount: u64,
        fee_value_juels: u64,
        token_transfers: vector<Sui2AnyTokenTransfer>,
        ctx: &mut TxContext
    ): Sui2AnyRampMessage {
        // calculate nonce
        let mut nonce = 0;
        if (!is_out_of_order_execution) {
            nonce = nonce_manager::get_incremented_outbound_nonce(
                ref, dest_chain_selector, sender, ctx
            );
        };

        let state = state_object::borrow<OnRampState>(ref, ON_RAMP_STATE_NAME);
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
            fee_token,
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

    // //
    // // MCMS entrypoint
    // //
    //
    // struct McmsCallback has drop {}
    //
    // public fun mcms_entrypoint<T: key>(
    //     _metadata: Object<T>
    // ): option::Option<u128> acquires OnRampState {
    //     let (caller, function, data) =
    //         mcms_registry::get_callback_params(@ccip, McmsCallback {});
    //
    //     let function_bytes = *string::bytes(&function);
    //     let stream = bcs_stream::new(data);
    //
    //     if (function_bytes == b"initialize") {
    //         let chain_selector = bcs_stream::deserialize_u64(&mut stream);
    //         let allowlist_admin = bcs_stream::deserialize_address(&mut stream);
    //         let dest_chain_selectors =
    //             bcs_stream::deserialize_vector(
    //                 &mut stream,
    //                 |stream| bcs_stream::deserialize_u64(stream)
    //             );
    //         let dest_chain_enabled =
    //             bcs_stream::deserialize_vector(
    //                 &mut stream,
    //                 |stream| bcs_stream::deserialize_bool(stream)
    //             );
    //         let dest_chain_allowlist_enabled =
    //             bcs_stream::deserialize_vector(
    //                 &mut stream,
    //                 |stream| bcs_stream::deserialize_bool(stream)
    //             );
    //         bcs_stream::assert_is_consumed(&stream);
    //         initialize(
    //             &caller,
    //             chain_selector,
    //             allowlist_admin,
    //             dest_chain_selectors,
    //             dest_chain_enabled,
    //             dest_chain_allowlist_enabled
    //         );
    //     } else if (function_bytes == b"set_dynamic_config") {
    //         let allowlist_admin = bcs_stream::deserialize_address(&mut stream);
    //         bcs_stream::assert_is_consumed(&stream);
    //         set_dynamic_config(&caller, allowlist_admin);
    //     } else if (function_bytes == b"apply_dest_chain_config_updates") {
    //         let dest_chain_selectors =
    //             bcs_stream::deserialize_vector(
    //                 &mut stream,
    //                 |stream| bcs_stream::deserialize_u64(stream)
    //             );
    //         let dest_chain_enabled =
    //             bcs_stream::deserialize_vector(
    //                 &mut stream,
    //                 |stream| bcs_stream::deserialize_bool(stream)
    //             );
    //         let dest_chain_allowlist_enabled =
    //             bcs_stream::deserialize_vector(
    //                 &mut stream,
    //                 |stream| bcs_stream::deserialize_bool(stream)
    //             );
    //         bcs_stream::assert_is_consumed(&stream);
    //         apply_dest_chain_config_updates(
    //             &caller,
    //             dest_chain_selectors,
    //             dest_chain_enabled,
    //             dest_chain_allowlist_enabled
    //         );
    //     } else if (function_bytes == b"apply_allowlist_updates") {
    //         let dest_chain_selectors =
    //             bcs_stream::deserialize_vector(
    //                 &mut stream,
    //                 |stream| bcs_stream::deserialize_u64(stream)
    //             );
    //         let dest_chain_allowlist_enabled =
    //             bcs_stream::deserialize_vector(
    //                 &mut stream,
    //                 |stream| bcs_stream::deserialize_bool(stream)
    //             );
    //         let dest_chain_add_allowed_senders =
    //             bcs_stream::deserialize_vector(
    //                 &mut stream,
    //                 |stream| bcs_stream::deserialize_vector(
    //                     stream,
    //                     |stream| bcs_stream::deserialize_address(stream)
    //                 )
    //             );
    //         let dest_chain_remove_allowed_senders =
    //             bcs_stream::deserialize_vector(
    //                 &mut stream,
    //                 |stream| bcs_stream::deserialize_vector(
    //                     stream,
    //                     |stream| bcs_stream::deserialize_address(stream)
    //                 )
    //             );
    //         bcs_stream::assert_is_consumed(&stream);
    //         apply_allowlist_updates(
    //             &caller,
    //             dest_chain_selectors,
    //             dest_chain_allowlist_enabled,
    //             dest_chain_add_allowed_senders,
    //             dest_chain_remove_allowed_senders
    //         );
    //     } else {
    //         abort E_UNKNOWN_FUNCTION)
    //     };
    //
    //     option::none()
    // }
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