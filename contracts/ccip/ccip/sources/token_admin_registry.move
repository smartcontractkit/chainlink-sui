module ccip::token_admin_registry {

    use std::string::{Self, String};
    use std::type_name::{Self, TypeName};

    use sui::coin::{CoinMetadata, TreasuryCap};
    use sui::event;
    use sui::table::{Self, Table};

    use ccip::state_object::{Self, CCIPObjectRef};

    const TOKEN_ADMIN_REGISTRY_STATE_NAME: vector<u8> = b"TokenAdminRegistry";
    const EXECUTION_STATE_IDLE: u8 = 1;
    const EXECUTION_STATE_LOCK_OR_BURN: u8 = 2;
    const EXECUTION_STATE_RELEASE_OR_MINT: u8 = 3;

    public struct FunctionInfo has store, copy, drop {
        module_address: address,
        module_name: String,
        function_name: String
    }

    public struct TokenAdminRegistryState has key, store {
        id: UID,

        // fungible asset metadata address -> TokenConfig
        // TODO: there were previously raised concerns during an audit that a user could maliciously calculate the bucket for a key and
        // cause repeated splitting, but we need to retrieve all the keys, which isn't available in Table.
        // consider other solutions.
        token_configs: Table<address, TokenConfig>,
    }

    public struct TokenConfig has store, drop, copy {
        token_pool_address: address,
        administrator: address,
        pending_administrator: address
    }

    public struct TokenPoolRegistration has key, store {
        id: UID,
        lock_or_burn_function: FunctionInfo,
        release_or_mint_function: FunctionInfo,
        proof_typename: TypeName,
        execution_state: u8,
        executing_lock_or_burn_input_v1: Option<LockOrBurnInputV1>,
        executing_release_or_mint_input_v1: Option<ReleaseOrMintInputV1>,
        executing_lock_or_burn_output_v1: Option<LockOrBurnOutputV1>,
        executing_release_or_mint_output_v1: Option<ReleaseOrMintOutputV1>
    }

    public struct LockOrBurnInputV1 has store, drop {
        sender: address,
        remote_chain_selector: u64,
        receiver: vector<u8>
    }

    public struct LockOrBurnOutputV1 has store, drop {
        dest_token_address: vector<u8>,
        dest_pool_data: vector<u8>
    }

    public struct ReleaseOrMintInputV1 has store, drop {
        sender: vector<u8>,
        receiver: address,
        source_amount: u256,
        local_token: address,
        remote_chain_selector: u64,
        source_pool_address: vector<u8>,
        source_pool_data: vector<u8>,
        offchain_token_data: vector<u8>
    }

    // TODO: consider removing ReleaseOrMintOutput, it exists only for a consistent UX across lock and release.
    public struct ReleaseOrMintOutputV1 has store, drop {
        destination_amount: u64
    }

    public struct PoolSet has copy, drop {
        local_token: address,
        previous_pool_address: address,
        new_pool_address: address
    }

    public struct AdministratorTransferRequested has copy, drop {
        local_token: address,
        current_admin: address,
        new_admin: address
    }

    public struct AdministratorTransferred has copy, drop {
        local_token: address,
        new_admin: address
    }

    const E_PROOF_NOT_IN_TOKEN_POOL_MODULE: u64 = 1;
    const E_PROOF_NOT_AT_TOKEN_POOL_ADDRESS: u64 = 2;
    const E_UNKNOWN_PROOF_TYPE: u64 = 3;
    const E_NOT_IN_IDLE_STATE: u64 = 4;
    const E_NOT_IN_LOCK_OR_BURN_STATE: u64 = 5;
    const E_NOT_IN_RELEASE_OR_MINT_STATE: u64 = 6;
    const E_NON_EMPTY_LOCK_OR_BURN_INPUT: u64 = 7;
    const E_NON_EMPTY_LOCK_OR_BURN_OUTPUT: u64 = 8;
    const E_NON_EMPTY_RELEASE_OR_MINT_INPUT: u64 = 9;
    const E_NON_EMPTY_RELEASE_OR_MINT_OUTPUT: u64 = 10;
    const E_MISSING_LOCK_OR_BURN_INPUT: u64 = 11;
    const E_MISSING_LOCK_OR_BURN_OUTPUT: u64 = 12;
    const E_MISSING_RELEASE_OR_MINT_INPUT: u64 = 13;
    const E_MISSING_RELEASE_OR_MINT_OUTPUT: u64 = 14;
    const E_FUNGIBLE_ASSET_ALREADY_REGISTERED: u64 = 15;
    const E_FUNGIBLE_ASSET_NOT_REGISTERED: u64 = 16;
    const E_NOT_ADMINISTRATOR: u64 = 17;
    const E_NOT_PENDING_ADMINISTRATOR: u64 = 18;

    public fun type_and_version(): String {
        string::utf8(b"TokenAdminRegistry 1.6.0")
    }

    // TODO: add MCMS support
    public fun initialize(ref: &mut CCIPObjectRef, ctx: &mut TxContext) {
        // if (@mcms_register_entrypoints != @0x0) {
        //     mcms_registry::register_entrypoint(
        //         publisher, string::utf8(b"token_admin_registry"), McmsCallback {}
        //     );
        // };

        let state = TokenAdminRegistryState {
            id: object::new(ctx),
            token_configs: table::new(ctx)
        };

        state_object::add(ref, TOKEN_ADMIN_REGISTRY_STATE_NAME, state, ctx);
    }

    public fun get_pools(
        ref: &CCIPObjectRef,
        local_tokens: vector<address>
    ): vector<address> {

        let state = state_object::borrow<TokenAdminRegistryState>(ref, TOKEN_ADMIN_REGISTRY_STATE_NAME);

        local_tokens.map_ref!(
            |local_token| {
                let local_token: address = *local_token;
                if (state.token_configs.contains(local_token)) {
                    let token_config = state.token_configs.borrow(local_token);
                    token_config.token_pool_address
                } else {
                    // returns @0x0 for assets without token pools.
                    @0x0
                }
            }
        )
    }

    // returns the token pool address for the given local token, or @0x0 if the token is not registered.
    public fun get_pool(ref: &CCIPObjectRef, local_token: address): address {
        let state = state_object::borrow<TokenAdminRegistryState>(ref, TOKEN_ADMIN_REGISTRY_STATE_NAME);

        if (state.token_configs.contains(local_token)) {
            let token_config = state.token_configs.borrow(local_token);
            token_config.token_pool_address
        } else {
            // returns @0x0 for assets without token pools.
            @0x0
        }
    }

    // returns (token_pool_address, administrator, pending_administrator)
    public fun get_token_config(
        ref: &CCIPObjectRef, local_token: address
    ): (address, address, address) {
        let state = state_object::borrow<TokenAdminRegistryState>(ref, TOKEN_ADMIN_REGISTRY_STATE_NAME);

        if (state.token_configs.contains(local_token)) {
            let token_config = state.token_configs.borrow(local_token);
            (
                token_config.token_pool_address,
                token_config.administrator,
                token_config.pending_administrator
            )
        } else {
            (@0x0, @0x0, @0x0)
        }
    }

    // public fun get_all_configured_tokens(
    //     starting_bucket_index: u64, starting_vector_index: u64, max_count: u64
    // ): (vector<address>, Option<u64>, Option<u64>) acquires TokenAdminRegistryState {
    //     // see the SmartTable documentation for descriptions of the function paramters and return values.
    //     // ref: https://github.com/aptos-labs/aptos-core/blob/6593fb81261f25490ffddc2252a861c994234c2a/aptos-move/framework/aptos-stdlib/sources/data_structures/smart_table.move#L212
    //
    //     let state = borrow_state();
    //     state.token_configs.keys_paginated(
    //         starting_bucket_index, starting_vector_index, max_count
    //     )
    // }

    // ================================================================
    // |                       Register Pool                          |
    // ================================================================

    #[allow(lint(self_transfer))]
    public fun register_pool<T: drop>(
        ref: &mut CCIPObjectRef,
        token_pool_module_name: vector<u8>,
        _treasury_cap: &TreasuryCap<T>, // pass in the treasury cap to ensure the caller owns it?
        // local_token: address,
        coin_metadata: &CoinMetadata<T>,
        // _proof: ProofType,
        ctx: &mut TxContext
    ) {
        let state = state_object::borrow_mut_from_user<TokenAdminRegistryState>(ref, TOKEN_ADMIN_REGISTRY_STATE_NAME);

        let local_token = object::id_to_address(&object::id(coin_metadata));
        // no need to validate it's a valid asset, bc coin_metadata will do that for us.

        let token_pool_address = ctx.sender();
        // assert!(
        //     !exists<TokenPoolRegistration>(token_pool_address),
        //      E_ALREADY_REGISTERED
        // );

        // TODO: figure out the permissioning model for the token pool registration

        assert!(
            !state.token_configs.contains(local_token),
            E_FUNGIBLE_ASSET_ALREADY_REGISTERED
        );

        // the initial administrator will always be the token pool account.
        // callers can immediately propose a new administrator afterwards if
        // needed.
        let token_config = TokenConfig {
            token_pool_address,
            administrator: token_pool_address,
            pending_administrator: @0x0
        };

        state.token_configs.add(local_token, token_config);

        let lock_or_burn_function =
            FunctionInfo {
                module_address: token_pool_address,
                module_name: string::utf8(token_pool_module_name),
                function_name: string::utf8(b"lock_or_burn")
            };
        let proof_typename = type_name::get<T>();
        assert!(
            proof_typename.get_address() == token_pool_address.to_ascii_string(),
            E_PROOF_NOT_AT_TOKEN_POOL_ADDRESS
        );
        assert!(
            proof_typename.get_module() == string::utf8(token_pool_module_name).to_ascii(),
            E_PROOF_NOT_IN_TOKEN_POOL_MODULE
        );

        let release_or_mint_function =
            FunctionInfo {
                module_address: token_pool_address,
                module_name: string::utf8(token_pool_module_name),
                function_name: string::utf8(b"release_or_mint")
            };

        transfer::transfer(
            TokenPoolRegistration {
                id: object::new(ctx),
                lock_or_burn_function,
                release_or_mint_function,
                proof_typename,
                execution_state: EXECUTION_STATE_IDLE,
                executing_lock_or_burn_input_v1: option::none(),
                executing_release_or_mint_input_v1: option::none(),
                executing_lock_or_burn_output_v1: option::none(),
                executing_release_or_mint_output_v1: option::none()
            },
            token_pool_address
        );
    }

    // depending on future development, we may need to save this registration object
    public fun set_pool(
        ref: &mut CCIPObjectRef, local_token: address, reg: &TokenPoolRegistration, ctx: &mut TxContext
    ) {
        let state = state_object::borrow_mut_from_user<TokenAdminRegistryState>(ref, TOKEN_ADMIN_REGISTRY_STATE_NAME);

        let token_pool_address = object::id_to_address(&object::id(reg));

        assert!(
            state.token_configs.contains(local_token),
            E_FUNGIBLE_ASSET_NOT_REGISTERED
        );

        let token_config = state.token_configs.borrow_mut(local_token);

        assert!(
            token_config.administrator == ctx.sender(),
            E_NOT_ADMINISTRATOR
        );

        let previous_pool_address = token_config.token_pool_address;
        if (previous_pool_address != token_pool_address) {
            token_config.token_pool_address = token_pool_address;

            event::emit(
                PoolSet {
                    local_token,
                    previous_pool_address,
                    new_pool_address: token_pool_address
                }
            );
        }
    }

    public fun transfer_admin_role(
        ref: &mut CCIPObjectRef, local_token: address, new_admin: address, ctx: &mut TxContext
    ) {
        let state = state_object::borrow_mut_from_user<TokenAdminRegistryState>(ref, TOKEN_ADMIN_REGISTRY_STATE_NAME);

        assert!(
            state.token_configs.contains(local_token),
            E_FUNGIBLE_ASSET_NOT_REGISTERED
        );

        let token_config = state.token_configs.borrow_mut(local_token);

        assert!(
            token_config.administrator == ctx.sender(),
            E_NOT_ADMINISTRATOR
        );

        // can be @0x0 to cancel a pending transfer.
        token_config.pending_administrator = new_admin;

        event::emit(
            AdministratorTransferRequested {
                local_token,
                current_admin: token_config.administrator,
                new_admin
            }
        );
    }

    public fun accept_admin_role(
        ref: &mut CCIPObjectRef, local_token: address, ctx: &mut TxContext
    ) {
        let state = state_object::borrow_mut_from_user<TokenAdminRegistryState>(ref, TOKEN_ADMIN_REGISTRY_STATE_NAME);

        assert!(
            state.token_configs.contains(local_token),
            E_FUNGIBLE_ASSET_NOT_REGISTERED
        );

        let token_config = state.token_configs.borrow_mut(local_token);

        assert!(
            token_config.pending_administrator == ctx.sender(),
            E_NOT_PENDING_ADMINISTRATOR
        );

        token_config.administrator = token_config.pending_administrator;
        token_config.pending_administrator = @0x0;

        event::emit(
            AdministratorTransferred { local_token, new_admin: token_config.administrator }
        );
    }

    public fun is_administrator(
        ref: &CCIPObjectRef, local_token: address, administrator: address
    ): bool {
        let state = state_object::borrow<TokenAdminRegistryState>(ref, TOKEN_ADMIN_REGISTRY_STATE_NAME);
        assert!(
            state.token_configs.contains(local_token),
            E_FUNGIBLE_ASSET_NOT_REGISTERED
        );

        let token_config = state.token_configs.borrow(local_token);
        token_config.administrator == administrator
    }

    // ================================================================
    // |                         Pool I/O V1                          |
    // ================================================================

    public fun get_lock_or_burn_input_v1<ProofType: drop>(
        registration: &mut TokenPoolRegistration,
        _proof: ProofType
    ): LockOrBurnInputV1 {
        assert!(
            type_name::get<ProofType>() == registration.proof_typename,
            E_UNKNOWN_PROOF_TYPE
        );

        assert!(
            registration.execution_state == EXECUTION_STATE_LOCK_OR_BURN,
            E_NOT_IN_LOCK_OR_BURN_STATE
        );
        assert!(
            registration.executing_lock_or_burn_input_v1.is_some(),
            E_MISSING_LOCK_OR_BURN_INPUT
        );
        assert!(
            registration.executing_lock_or_burn_output_v1.is_none(),
            E_NON_EMPTY_LOCK_OR_BURN_OUTPUT
        );
        assert!(
            registration.executing_release_or_mint_input_v1.is_none(),
            E_NON_EMPTY_RELEASE_OR_MINT_INPUT
        );
        assert!(
            registration.executing_release_or_mint_output_v1.is_none(),
            E_NON_EMPTY_RELEASE_OR_MINT_OUTPUT
        );

        registration.executing_lock_or_burn_input_v1.extract()
    }

    public fun set_lock_or_burn_output_v1<ProofType: drop>(
        registration: &mut TokenPoolRegistration,
        _proof: ProofType,
        dest_token_address: vector<u8>,
        dest_pool_data: vector<u8>
    ) {
        assert!(
            type_name::get<ProofType>() == registration.proof_typename,
            E_UNKNOWN_PROOF_TYPE
        );

        assert!(
            registration.execution_state == EXECUTION_STATE_LOCK_OR_BURN,
            E_NOT_IN_LOCK_OR_BURN_STATE
        );
        assert!(
            registration.executing_lock_or_burn_input_v1.is_none(),
            E_NON_EMPTY_LOCK_OR_BURN_INPUT
        );
        assert!(
            registration.executing_lock_or_burn_output_v1.is_none(),
            E_NON_EMPTY_LOCK_OR_BURN_OUTPUT
        );
        assert!(
            registration.executing_release_or_mint_input_v1.is_none(),
            E_NON_EMPTY_RELEASE_OR_MINT_INPUT
        );
        assert!(
            registration.executing_release_or_mint_output_v1.is_none(),
            E_NON_EMPTY_RELEASE_OR_MINT_OUTPUT
        );

        registration.executing_lock_or_burn_output_v1.fill(
            LockOrBurnOutputV1 { dest_token_address, dest_pool_data }
        )
    }

    public fun get_release_or_mint_input_v1<ProofType: drop>(
        registration: &mut TokenPoolRegistration,
        _proof: ProofType
    ): ReleaseOrMintInputV1 {

        assert!(
            type_name::get<ProofType>() == registration.proof_typename,
            E_UNKNOWN_PROOF_TYPE
        );

        assert!(
            registration.execution_state == EXECUTION_STATE_RELEASE_OR_MINT,
            E_NOT_IN_RELEASE_OR_MINT_STATE
        );
        assert!(
            registration.executing_release_or_mint_input_v1.is_some(),
            E_MISSING_RELEASE_OR_MINT_INPUT
        );
        assert!(
            registration.executing_release_or_mint_output_v1.is_none(),
            E_NON_EMPTY_RELEASE_OR_MINT_OUTPUT
        );
        assert!(
            registration.executing_lock_or_burn_input_v1.is_none(),
            E_NON_EMPTY_LOCK_OR_BURN_INPUT
        );
        assert!(
            registration.executing_lock_or_burn_output_v1.is_none(),
            E_NON_EMPTY_LOCK_OR_BURN_OUTPUT
        );

        registration.executing_release_or_mint_input_v1.extract()
    }

    public fun set_release_or_mint_output_v1<ProofType: drop>(
        registration: &mut TokenPoolRegistration,
        _proof: ProofType,
        destination_amount: u64
    ) {

        assert!(
            type_name::get<ProofType>() == registration.proof_typename,
            E_UNKNOWN_PROOF_TYPE
        );

        assert!(
            registration.execution_state == EXECUTION_STATE_RELEASE_OR_MINT,
            E_NOT_IN_RELEASE_OR_MINT_STATE
        );
        assert!(
            registration.executing_release_or_mint_input_v1.is_none(),
            E_NON_EMPTY_RELEASE_OR_MINT_INPUT
        );
        assert!(
            registration.executing_release_or_mint_output_v1.is_none(),
            E_NON_EMPTY_RELEASE_OR_MINT_OUTPUT
        );
        assert!(
            registration.executing_lock_or_burn_input_v1.is_none(),
            E_NON_EMPTY_LOCK_OR_BURN_INPUT
        );
        assert!(
            registration.executing_lock_or_burn_output_v1.is_none(),
            E_NON_EMPTY_LOCK_OR_BURN_OUTPUT
        );

        registration.executing_release_or_mint_output_v1.fill(
            ReleaseOrMintOutputV1 { destination_amount }
        )
    }

    // LockOrBurnInput accessors
    public fun get_lock_or_burn_sender(input: &LockOrBurnInputV1): address {
        input.sender
    }

    public fun get_lock_or_burn_remote_chain_selector(
        input: &LockOrBurnInputV1
    ): u64 {
        input.remote_chain_selector
    }

    public fun get_lock_or_burn_receiver(input: &LockOrBurnInputV1): vector<u8> {
        input.receiver
    }

    // ReleaseOrMintInput accessors
    public fun get_release_or_mint_sender(input: &ReleaseOrMintInputV1): vector<u8> {
        input.sender
    }

    public fun get_release_or_mint_receiver(input: &ReleaseOrMintInputV1): address {
        input.receiver
    }

    public fun get_release_or_mint_source_amount(
        input: &ReleaseOrMintInputV1
    ): u256 {
        input.source_amount
    }

    public fun get_release_or_mint_local_token(
        input: &ReleaseOrMintInputV1
    ): address {
        input.local_token
    }

    public fun get_release_or_mint_remote_chain_selector(
        input: &ReleaseOrMintInputV1
    ): u64 {
        input.remote_chain_selector
    }

    public fun get_release_or_mint_source_pool_address(
        input: &ReleaseOrMintInputV1
    ): vector<u8> {
        input.source_pool_address
    }

    public fun get_release_or_mint_source_pool_data(
        input: &ReleaseOrMintInputV1
    ): vector<u8> {
        input.source_pool_data
    }

    public fun get_release_or_mint_offchain_token_data(
        input: &ReleaseOrMintInputV1
    ): vector<u8> {
        input.offchain_token_data
    }

    // ================================================================
    // |                        Lock or Burn                          |
    // ================================================================

    // TODO: revisit the return value after dynamic dispatch work-around is decided.
    public(package) fun start_lock_or_burn(
        registration: &mut TokenPoolRegistration,
        sender: address,
        remote_chain_selector: u64,
        receiver: vector<u8>
    ) {
        assert!(
            registration.execution_state == EXECUTION_STATE_IDLE,
            E_NOT_IN_IDLE_STATE
        );
        assert!(
            registration.executing_lock_or_burn_input_v1.is_none(),
            E_NON_EMPTY_LOCK_OR_BURN_INPUT
        );
        assert!(
            registration.executing_lock_or_burn_output_v1.is_none(),
            E_NON_EMPTY_LOCK_OR_BURN_OUTPUT
        );
        assert!(
            registration.executing_release_or_mint_input_v1.is_none(),
            E_NON_EMPTY_RELEASE_OR_MINT_INPUT
        );
        assert!(
            registration.executing_release_or_mint_output_v1.is_none(),
            E_NON_EMPTY_RELEASE_OR_MINT_OUTPUT
        );

        registration.execution_state = EXECUTION_STATE_LOCK_OR_BURN;
        registration.executing_lock_or_burn_input_v1.fill(
            LockOrBurnInputV1 { sender, remote_chain_selector, receiver }
        );
    }

    public(package) fun finish_lock_or_burn(
        registration: &mut TokenPoolRegistration
    ): (vector<u8>, vector<u8>) {
        assert!(
            registration.execution_state == EXECUTION_STATE_LOCK_OR_BURN,
            E_NOT_IN_LOCK_OR_BURN_STATE
        );
        assert!(
            registration.executing_lock_or_burn_input_v1.is_none(),
            E_NON_EMPTY_LOCK_OR_BURN_INPUT
        );
        assert!(
            registration.executing_lock_or_burn_output_v1.is_some(),
            E_MISSING_LOCK_OR_BURN_OUTPUT
        );
        assert!(
            registration.executing_release_or_mint_input_v1.is_none(),
            E_NON_EMPTY_RELEASE_OR_MINT_INPUT
        );
        assert!(
            registration.executing_release_or_mint_output_v1.is_none(),
            E_NON_EMPTY_RELEASE_OR_MINT_OUTPUT
        );

        registration.execution_state = EXECUTION_STATE_IDLE;

        // // the dispatch callback is passed a fungible_asset::TransferRef reference which could allow the store to be frozen,
        // // causing future deposit/withdraw callbacks to fail. note that this fungible store is only used as part of the dispatch
        // // mechanism.
        // // ref: https://github.com/aptos-labs/aptos-core/blob/7fc73792e9db11462c9a42038c4a9eb41cc00192/aptos-move/framework/aptos-framework/sources/fungible_asset.move#L923
        // if (fungible_asset::is_frozen(registration.dispatch_deposit_fungible_store)) {
        //     fungible_asset::set_frozen_flag(
        //         &registration.dispatch_fa_transfer_ref,
        //         registration.dispatch_deposit_fungible_store,
        //         false
        //     );
        // };

        let output = registration.executing_lock_or_burn_output_v1.extract();
        (output.dest_token_address, output.dest_pool_data)
    }

    // ================================================================
    // |                       Release or Mint                        |
    // ================================================================

    // TODO: revisit the return value after dynamic dispatch work-around is decided.
    public(package) fun start_release_or_mint(
        registration: &mut TokenPoolRegistration,
        sender: vector<u8>,
        receiver: address,
        source_amount: u256,
        local_token: address,
        remote_chain_selector: u64,
        source_pool_address: vector<u8>,
        source_pool_data: vector<u8>,
        offchain_token_data: vector<u8>
    ) {
        assert!(
            registration.execution_state == EXECUTION_STATE_IDLE,
            E_NOT_IN_IDLE_STATE
        );
        assert!(
            registration.executing_release_or_mint_input_v1.is_none(),
            E_NON_EMPTY_RELEASE_OR_MINT_INPUT
        );
        assert!(
            registration.executing_release_or_mint_output_v1.is_none(),
            E_NON_EMPTY_RELEASE_OR_MINT_OUTPUT
        );
        assert!(
            registration.executing_lock_or_burn_input_v1.is_none(),
            E_NON_EMPTY_LOCK_OR_BURN_INPUT
        );
        assert!(
            registration.executing_lock_or_burn_output_v1.is_none(),
            E_NON_EMPTY_LOCK_OR_BURN_OUTPUT
        );

        registration.execution_state = EXECUTION_STATE_RELEASE_OR_MINT;
        registration.executing_release_or_mint_input_v1.fill(
            ReleaseOrMintInputV1 {
                sender,
                receiver,
                source_amount,
                local_token,
                remote_chain_selector,
                source_pool_address,
                source_pool_data,
                offchain_token_data
            }
        );

        // (
        //     object::generate_signer_for_extending(&registration.dispatch_extend_ref),
        //     registration.dispatch_deposit_fungible_store
        // )
    }

    public(package) fun finish_release_or_mint(
        registration: &mut TokenPoolRegistration
    ): u64 {
        assert!(
            registration.execution_state == EXECUTION_STATE_RELEASE_OR_MINT,
            E_NOT_IN_RELEASE_OR_MINT_STATE
        );
        assert!(
            registration.executing_release_or_mint_input_v1.is_none(),
            E_NON_EMPTY_RELEASE_OR_MINT_INPUT
        );
        assert!(
            registration.executing_release_or_mint_output_v1.is_some(),
            E_MISSING_RELEASE_OR_MINT_OUTPUT
        );
        assert!(
            registration.executing_lock_or_burn_input_v1.is_none(),
            E_NON_EMPTY_LOCK_OR_BURN_INPUT
        );
        assert!(
            registration.executing_lock_or_burn_output_v1.is_none(),
            E_NON_EMPTY_LOCK_OR_BURN_OUTPUT
        );

        registration.execution_state = EXECUTION_STATE_IDLE;

        // // the dispatch callback is passed a fungible_asset::TransferRef reference which could allow the store to be frozen,
        // // causing future deposit/withdraw callbacks to fail. note that this fungible store is only used as part of the dispatch
        // // mechanism.
        // // ref: https://github.com/aptos-labs/aptos-core/blob/7fc73792e9db11462c9a42038c4a9eb41cc00192/aptos-move/framework/aptos-framework/sources/fungible_asset.move#L936
        // if (fungible_asset::is_frozen(registration.dispatch_deposit_fungible_store)) {
        //     fungible_asset::set_frozen_flag(
        //         &registration.dispatch_fa_transfer_ref,
        //         registration.dispatch_deposit_fungible_store,
        //         false
        //     );
        // };

        let output = registration.executing_release_or_mint_output_v1.extract();

        output.destination_amount
    }

    // // ================================================================
    // // |                      MCMS Entrypoint                         |
    // // ================================================================
    //
    // struct McmsCallback has drop {}
    //
    // public fun mcms_entrypoint<T: key>(
    //     _metadata: Object<T>
    // ): option::Option<u128> acquires TokenAdminRegistryState {
    //     let (caller, function, data) =
    //         mcms_registry::get_callback_params(@ccip, McmsCallback {});
    //
    //     let function_bytes = *string::bytes(&function);
    //     let stream = bcs_stream::new(data);
    //
    //     if (function_bytes == b"set_pool") {
    //         let local_token = bcs_stream::deserialize_address(&mut stream);
    //         let token_pool_address = bcs_stream::deserialize_address(&mut stream);
    //         bcs_stream::assert_is_consumed(&stream);
    //         set_pool(&caller, local_token, token_pool_address)
    //     } else if (function_bytes == b"transfer_admin_role") {
    //         let local_token = bcs_stream::deserialize_address(&mut stream);
    //         let new_admin = bcs_stream::deserialize_address(&mut stream);
    //         bcs_stream::assert_is_consumed(&stream);
    //         transfer_admin_role(&caller, local_token, new_admin)
    //     } else if (function_bytes == b"accept_admin_role") {
    //         let local_token = bcs_stream::deserialize_address(&mut stream);
    //         bcs_stream::assert_is_consumed(&stream);
    //         accept_admin_role(&caller, local_token)
    //     } else {
    //         abort  E_UNKNOWN_FUNCTION)
    //     };
    //
    //     option::none()
    // }

    #[test_only]
    public(package) fun get_registration(
        reg: &TokenPoolRegistration
    ): (FunctionInfo, FunctionInfo, TypeName, u8) {
        (
            reg.lock_or_burn_function,
            reg.release_or_mint_function,
            reg.proof_typename,
            reg.execution_state
        )
    }

    #[test_only]
    public(package) fun get_function_info(
        reg: &FunctionInfo
    ): (address, String, String) {
        (
            reg.module_address,
            reg.module_name,
            reg.function_name
        )
    }

    #[test_only]
    public(package) fun destroy_registration(
        reg: TokenPoolRegistration
    ) {
        let TokenPoolRegistration {
            id,
            lock_or_burn_function: _lbf,
            release_or_mint_function: _rmf,
            proof_typename: _pt,
            execution_state: _es,
            executing_lock_or_burn_input_v1: _lobi,
            executing_release_or_mint_input_v1: _rmi,
            executing_lock_or_burn_output_v1: _lobo,
            executing_release_or_mint_output_v1: _rmo
        } = reg;
        object::delete(id);
    }
}

// the name here is not conventional, but it is used to avoid check failures in register_pool
#[test_only]
module ccip::token_pool {
    use std::string;

    use sui::coin;
    use sui::test_scenario::{Self, Scenario};

    use ccip::token_admin_registry as registry;
    use ccip::token_admin_registry::TokenPoolRegistration;
    use ccip::state_object::{Self, CCIPObjectRef};

    public struct TOKEN_POOL has drop {}
    const Decimals: u8 = 8;
    const EXECUTION_STATE_IDLE: u8 = 1;

    fun set_up_test(): (Scenario, CCIPObjectRef) {
        let mut scenario = test_scenario::begin(@0x1000);
        let ctx = scenario.ctx();

        let ref = state_object::create(ctx);

        (scenario, ref)
    }

    fun initialize(ref: &mut CCIPObjectRef, ctx: &mut TxContext) {
        registry::initialize(ref, ctx);
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

        let (token_pool_address, administrator, pending_administrator) = registry::get_token_config(&ref, @0x2);
        assert!(token_pool_address == @0x0);
        assert!(administrator == @0x0);
        assert!(pending_administrator == @0x0);

        tear_down_test(scenario, ref);
    }

    #[test]
    #[expected_failure(abort_code = registry::E_FUNGIBLE_ASSET_NOT_REGISTERED)]
    public fun test_transfer_admin_role_not_registered() {
        let (mut scenario, mut ref) = set_up_test();
        let ctx = scenario.ctx();
        initialize(&mut ref, ctx);

        registry::transfer_admin_role(&mut ref, @0x2, @0x3, ctx);

        tear_down_test(scenario, ref);
    }

    #[test]
    public fun test_register_and_set_pool() {
        let (mut scenario, mut ref) = set_up_test();
        let ctx = scenario.ctx();
        initialize(&mut ref, ctx);

        let (treasury_cap, coin_metadata) = coin::create_currency(
            TOKEN_POOL {},
            Decimals,
            b"TEST",
            b"TestToken",
            b"test_token",
            option::none(),
            ctx
        );

        registry::register_pool(
            &mut ref,
            b"token_pool",
            &treasury_cap,
            &coin_metadata,
            ctx
        );

        let local_token = object::id_to_address(&object::id(&coin_metadata));
        let pool_addresses = registry::get_pools(&ref, vector[local_token]);
        assert!(pool_addresses.length() == 1);
        assert!(pool_addresses[0] == ctx.sender());

        assert!(registry::is_administrator(&ref, local_token, ctx.sender()));

        let (token_pool_address, administrator, pending_administrator) = registry::get_token_config(&ref, local_token);
        assert!(token_pool_address == ctx.sender());
        assert!(administrator == ctx.sender());
        assert!(pending_administrator == @0x0);

        transfer::public_freeze_object(coin_metadata);
        transfer::public_transfer(treasury_cap, ctx.sender());

        let effects = test_scenario::next_epoch(&mut scenario, @0x1000);
        let created = test_scenario::created(&effects);
        assert!(created.length() == 3);
        test_scenario::end(scenario);

        let mut scenario_2 = test_scenario::begin(@0x1000);
        // the registration is the first object
        let registration = test_scenario::take_from_sender_by_id<TokenPoolRegistration>(&scenario_2, created[0]);
        let ctx_2 = scenario_2.ctx();

        let (lock_and_burn, release_and_mint, type_name, execution_state) = registry::get_registration(&registration);
        assert!(execution_state == EXECUTION_STATE_IDLE);
        let (lock_and_burn_address, lock_and_burn_module_name, lock_and_burn_function_name) = registry::get_function_info(&lock_and_burn);
        assert!(lock_and_burn_address == ctx_2.sender());
        assert!(lock_and_burn_module_name == string::utf8(b"token_pool"));
        assert!(lock_and_burn_function_name == string::utf8(b"lock_or_burn"));
        let (release_and_mint_address, release_and_mint_module_name, release_and_mint_function_name) = registry::get_function_info(&release_and_mint);
        assert!(release_and_mint_address == ctx_2.sender());
        assert!(release_and_mint_module_name == string::utf8(b"token_pool"));
        assert!(release_and_mint_function_name == string::utf8(b"release_or_mint"));
        assert!(type_name.get_address() == ctx_2.sender().to_ascii_string());
        assert!(type_name.get_module() == string::utf8(b"token_pool").to_ascii());

        registry::set_pool(&mut ref, local_token, &registration, ctx_2);

        registry::transfer_admin_role(&mut ref, local_token, @0x3000, ctx_2);
        scenario_2.end();

        let mut scenario_3 = test_scenario::begin(@0x3000);
        let ctx_3 = scenario_3.ctx();
        registry::accept_admin_role(&mut ref, local_token, ctx_3);
        assert!(registry::is_administrator(&ref, local_token, @0x3000));

        scenario_3.end();
        state_object::destroy_state_object(ref);
        registry::destroy_registration(registration);
    }
}