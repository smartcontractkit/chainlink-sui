module ccip::token_admin_registry {

    use std::string::{Self, String};
    use std::type_name;

    use sui::address;
    use sui::coin::TreasuryCap;
    use sui::event;
    use sui::table::{Self, Table};

    use ccip::state_object::{Self, CCIPObjectRef, OwnerCap};

    const TOKEN_ADMIN_REGISTRY_STATE_NAME: vector<u8> = b"TokenAdminRegistry";

    public struct FunctionInfo has store, copy, drop {
        module_address: address,
        module_name: String,
        function_name: String
    }

    public struct TokenAdminRegistryState has key, store {
        id: UID,
        token_configs: Table<address, TokenConfig>,
    }

    public struct TokenConfig has store, drop, copy {
        token_pool_address: address,
        administrator: address,
        pending_administrator: address
    }

    public struct PoolSet has copy, drop {
        coin_metadata_address: address,
        previous_pool_address: address,
        new_pool_address: address
    }

    public struct AdministratorTransferRequested has copy, drop {
        coin_metadata_address: address,
        current_admin: address,
        new_admin: address
    }

    public struct AdministratorTransferred has copy, drop {
        coin_metadata_address: address,
        new_admin: address
    }

    const E_PROOF_NOT_IN_TOKEN_POOL_MODULE: u64 = 1;
    const E_PROOF_NOT_AT_TOKEN_POOL_ADDRESS: u64 = 2;
    const E_ALREADY_INITIALIZED: u64 = 3;
    const E_FUNGIBLE_ASSET_ALREADY_REGISTERED: u64 = 4;
    const E_FUNGIBLE_ASSET_NOT_REGISTERED: u64 = 5;
    const E_NOT_ADMINISTRATOR: u64 = 6;
    const E_NOT_PENDING_ADMINISTRATOR: u64 = 7;

    public fun type_and_version(): String {
        string::utf8(b"TokenAdminRegistry 1.6.0")
    }

    public fun initialize(
        ref: &mut CCIPObjectRef,
        _: &OwnerCap,
        ctx: &mut TxContext
    ) {
        assert!(
            !state_object::contains(ref, TOKEN_ADMIN_REGISTRY_STATE_NAME),
            E_ALREADY_INITIALIZED
        );
        let state = TokenAdminRegistryState {
            id: object::new(ctx),
            token_configs: table::new(ctx)
        };

        state_object::add(ref, TOKEN_ADMIN_REGISTRY_STATE_NAME, state, ctx);
    }

    public fun get_pools(
        ref: &CCIPObjectRef,
        coin_metadata_addresses: vector<address>
    ): vector<address>{
        let state = state_object::borrow<TokenAdminRegistryState>(ref, TOKEN_ADMIN_REGISTRY_STATE_NAME);

        let mut token_pool_addresses: vector<address> = vector::empty();
        coin_metadata_addresses.do_ref!(
            |metadata_address| {
                let metadata_address: address = *metadata_address;
                if (state.token_configs.contains(metadata_address)) {
                    let token_config = state.token_configs.borrow(metadata_address);
                    token_pool_addresses.push_back(token_config.token_pool_address);
                } else {
                    // returns @0x0 for assets without token pools.
                    token_pool_addresses.push_back(@0x0);
                }
            }
        );

        token_pool_addresses
    }

    // this function can also take a coin metadata or a coin::zero
    // but that requires adding a type parameter to the function
    public fun get_pool(ref: &CCIPObjectRef, coin_metadata_address: address): address {
        let state = state_object::borrow<TokenAdminRegistryState>(ref, TOKEN_ADMIN_REGISTRY_STATE_NAME);

        if (state.token_configs.contains(coin_metadata_address)) {
            let token_config = state.token_configs.borrow(coin_metadata_address);
            token_config.token_pool_address
        } else {
            // returns @0x0 for assets without token pools.
            @0x0
        }
    }

    public fun get_token_config(
        ref: &CCIPObjectRef, coin_metadata_address: address
    ): (address, address, address) {
        let state = state_object::borrow<TokenAdminRegistryState>(ref, TOKEN_ADMIN_REGISTRY_STATE_NAME);

        if (state.token_configs.contains(coin_metadata_address)) {
            let token_config = state.token_configs.borrow(coin_metadata_address);
            (
                token_config.token_pool_address,
                token_config.administrator,
                token_config.pending_administrator
            )
        } else {
            (@0x0, @0x0, @0x0)
        }
    }

    // TODO: this cannot be supported with basic Sui libraries unless we implement it
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

    // only the token owner with the treasury cap can call this function.
    #[allow(lint(self_transfer))]
    public fun register_pool<T, ProofType: drop>(
        ref: &mut CCIPObjectRef,
        _: &TreasuryCap<T>, // passing in the treasury cap to demonstrate ownership over the token
        coin_metadata_address: address,
        token_pool_address: address,
        token_pool_module_name: vector<u8>,
        _proof: ProofType,
        ctx: &mut TxContext
    ) {
        register_pool_internal(
            ref,
            coin_metadata_address,
            token_pool_address,
            token_pool_module_name,
            _proof,
            ctx
        );
    }

    // only the CCIP owner can call this function.
    public fun register_pool_by_admin<ProofType: drop>(
        ref: &mut CCIPObjectRef,
        coin_metadata_address: address,
        token_pool_address: address,
        token_pool_module_name: vector<u8>,
        _proof: ProofType,
        ctx: &mut TxContext
    ) {
        assert!(
            ctx.sender() == state_object::get_current_owner(ref),
            E_NOT_ADMINISTRATOR
        );

        register_pool_internal(
            ref,
            coin_metadata_address,
            token_pool_address,
            token_pool_module_name,
            _proof,
            ctx
        );
    }

    fun register_pool_internal<ProofType: drop>(
        ref: &mut CCIPObjectRef,
        coin_metadata_address: address,
        token_pool_address: address,
        token_pool_module_name: vector<u8>,
        _proof: ProofType, // use this proof type to validate the token pool address & token pool module name
        ctx: &TxContext
    ) {
        let state = state_object::borrow_mut<TokenAdminRegistryState>(ref, TOKEN_ADMIN_REGISTRY_STATE_NAME);

        let proof_tn = type_name::get<ProofType>();
        let pool_bytes = proof_tn.get_address().into_bytes();
        assert!(
            token_pool_address == address::from_ascii_bytes(&pool_bytes),
            E_PROOF_NOT_AT_TOKEN_POOL_ADDRESS
        );
        assert!(
            string::utf8(token_pool_module_name).to_ascii() == proof_tn.get_module(),
            E_PROOF_NOT_IN_TOKEN_POOL_MODULE
        );
        assert!(
            !state.token_configs.contains(coin_metadata_address),
            E_FUNGIBLE_ASSET_ALREADY_REGISTERED
        );

        // the initial administrator will always be the either the token pool owner or CCIP admin
        let token_config = TokenConfig {
            token_pool_address,
            administrator: ctx.sender(),
            pending_administrator: @0x0,
        };

        state.token_configs.add(coin_metadata_address, token_config);
    }

    public fun set_pool<ProofType: drop>(
        ref: &mut CCIPObjectRef,
        coin_metadata_address: address,
        _proof: ProofType,
        ctx: &mut TxContext
    ) {
        let state = state_object::borrow_mut<TokenAdminRegistryState>(ref, TOKEN_ADMIN_REGISTRY_STATE_NAME);

        // let token_pool_address = object::id_to_address(&object::id(reg));
        // let coin_tn = type_name::get<T>();

        let proof_tn = type_name::get<ProofType>();
        let pool_bytes = proof_tn.get_address().into_bytes();
        let token_pool_address = address::from_ascii_bytes(&pool_bytes);

        assert!(
            state.token_configs.contains(coin_metadata_address),
            E_FUNGIBLE_ASSET_NOT_REGISTERED
        );

        let token_config = state.token_configs.borrow_mut(coin_metadata_address);

        // the tx signer must be the administrator of the token pool.
        assert!(
            token_config.administrator == ctx.sender(),
            E_NOT_ADMINISTRATOR
        );

        let previous_pool_address = token_config.token_pool_address;
        if (previous_pool_address != token_pool_address) {
            token_config.token_pool_address = token_pool_address;

            event::emit(
                PoolSet {
                    coin_metadata_address,
                    previous_pool_address,
                    new_pool_address: token_pool_address
                }
            );
        }
    }

    public fun transfer_admin_role(
        ref: &mut CCIPObjectRef,
        coin_metadata_address: address,
        new_admin: address,
        ctx: &mut TxContext
    ) {
        let state = state_object::borrow_mut<TokenAdminRegistryState>(ref, TOKEN_ADMIN_REGISTRY_STATE_NAME);

        assert!(
            state.token_configs.contains(coin_metadata_address),
            E_FUNGIBLE_ASSET_NOT_REGISTERED
        );

        let token_config = state.token_configs.borrow_mut(coin_metadata_address);

        assert!(
            token_config.administrator == ctx.sender(),
            E_NOT_ADMINISTRATOR
        );

        // can be @0x0 to cancel a pending transfer.
        token_config.pending_administrator = new_admin;

        event::emit(
            AdministratorTransferRequested {
                coin_metadata_address,
                current_admin: token_config.administrator,
                new_admin
            }
        );
    }

    public fun accept_admin_role(
        ref: &mut CCIPObjectRef,
        coin_metadata_address: address,
        ctx: &mut TxContext
    ) {
        let state = state_object::borrow_mut<TokenAdminRegistryState>(ref, TOKEN_ADMIN_REGISTRY_STATE_NAME);

        assert!(
            state.token_configs.contains(coin_metadata_address),
            E_FUNGIBLE_ASSET_NOT_REGISTERED
        );

        let token_config = state.token_configs.borrow_mut(coin_metadata_address);

        assert!(
            token_config.pending_administrator == ctx.sender(),
            E_NOT_PENDING_ADMINISTRATOR
        );

        token_config.administrator = token_config.pending_administrator;
        token_config.pending_administrator = @0x0;

        event::emit(
            AdministratorTransferred { coin_metadata_address, new_admin: token_config.administrator }
        );
    }

    public fun is_administrator(
        ref: &CCIPObjectRef, coin_metadata_address: address, administrator: address
    ): bool {
        let state = state_object::borrow<TokenAdminRegistryState>(ref, TOKEN_ADMIN_REGISTRY_STATE_NAME);

        assert!(
            state.token_configs.contains(coin_metadata_address),
            E_FUNGIBLE_ASSET_NOT_REGISTERED
        );

        let token_config = state.token_configs.borrow(coin_metadata_address);
        token_config.administrator == administrator
    }

    // #[test_only]
    // public(package) fun get_registration(
    //     reg: &TokenPoolRegistration
    // ): (FunctionInfo, FunctionInfo, TypeName, u8) {
    //     (
    //         reg.lock_or_burn_function,
    //         reg.release_or_mint_function,
    //         reg.proof_typename,
    //         reg.execution_state
    //     )
    // }

    // #[test_only]
    // public(package) fun get_function_info(
    //     reg: &FunctionInfo
    // ): (address, String, String) {
    //     (
    //         reg.module_address,
    //         reg.module_name,
    //         reg.function_name
    //     )
    // }
    //
    // #[test_only]
    // public(package) fun destroy_registration(
    //     reg: TokenPoolRegistration
    // ) {
    //     let TokenPoolRegistration {
    //         id,
    //         lock_or_burn_function: _lbf,
    //         release_or_mint_function: _rmf,
    //         proof_typename: _pt,
    //         execution_state: _es,
    //         executing_lock_or_burn_input_v1: _lobi,
    //         executing_release_or_mint_input_v1: _rmi,
    //         executing_lock_or_burn_output_v1: _lobo,
    //         executing_release_or_mint_output_v1: _rmo
    //     } = reg;
    //     object::delete(id);
    // }
}

// the name here is not conventional, but it is used to avoid check failures in register_pool
// #[test_only]
// module ccip::token_pool {
//     use std::string;
//
//     use sui::coin;
//     use sui::test_scenario::{Self, Scenario};
//
//     use ccip::token_admin_registry as registry;
//     use ccip::token_admin_registry::TokenPoolRegistration;
//     use ccip::state_object::{Self, CCIPObjectRef};
//
//     public struct TOKEN_POOL has drop {}
//     const Decimals: u8 = 8;
//     const EXECUTION_STATE_IDLE: u8 = 1;
//
//     fun set_up_test(): (Scenario, CCIPObjectRef) {
//         let mut scenario = test_scenario::begin(@0x1000);
//         let ctx = scenario.ctx();
//
//         let ref = state_object::create(ctx);
//
//         (scenario, ref)
//     }
//
//     fun initialize(ref: &mut CCIPObjectRef, ctx: &mut TxContext) {
//         registry::initialize(ref, ctx);
//     }
//
//     fun tear_down_test(scenario: Scenario, ref: CCIPObjectRef) {
//         state_object::destroy_state_object(ref);
//         test_scenario::end(scenario);
//     }
//
//     #[test]
//     public fun test_initialize() {
//         let (mut scenario, mut ref) = set_up_test();
//         let ctx = scenario.ctx();
//         initialize(&mut ref, ctx);
//
//         let (token_pool_address, administrator, pending_administrator) = registry::get_token_config(&ref, @0x2);
//         assert!(token_pool_address == @0x0);
//         assert!(administrator == @0x0);
//         assert!(pending_administrator == @0x0);
//
//         tear_down_test(scenario, ref);
//     }
//
//     #[test]
//     #[expected_failure(abort_code = registry::E_FUNGIBLE_ASSET_NOT_REGISTERED)]
//     public fun test_transfer_admin_role_not_registered() {
//         let (mut scenario, mut ref) = set_up_test();
//         let ctx = scenario.ctx();
//         initialize(&mut ref, ctx);
//
//         registry::transfer_admin_role(&mut ref, @0x2, @0x3, ctx);
//
//         tear_down_test(scenario, ref);
//     }
//
//     #[test]
//     public fun test_register_and_set_pool() {
//         let (mut scenario, mut ref) = set_up_test();
//         let ctx = scenario.ctx();
//         initialize(&mut ref, ctx);
//
//         let (treasury_cap, coin_metadata) = coin::create_currency(
//             TOKEN_POOL {},
//             Decimals,
//             b"TEST",
//             b"TestToken",
//             b"test_token",
//             option::none(),
//             ctx
//         );
//
//         registry::register_pool(
//             &mut ref,
//             b"token_pool",
//             &treasury_cap,
//             &coin_metadata,
//             ctx
//         );
//
//         let local_token = object::id_to_address(&object::id(&coin_metadata));
//         let pool_addresses = registry::get_pools(&ref, vector[local_token]);
//         assert!(pool_addresses.length() == 1);
//         assert!(pool_addresses[0] == ctx.sender());
//
//         assert!(registry::is_administrator(&ref, local_token, ctx.sender()));
//
//         let (token_pool_address, administrator, pending_administrator) = registry::get_token_config(&ref, local_token);
//         assert!(token_pool_address == ctx.sender());
//         assert!(administrator == ctx.sender());
//         assert!(pending_administrator == @0x0);
//
//         transfer::public_freeze_object(coin_metadata);
//         transfer::public_transfer(treasury_cap, ctx.sender());
//
//         let effects = test_scenario::next_epoch(&mut scenario, @0x1000);
//         let created = test_scenario::created(&effects);
//         assert!(created.length() == 3);
//         test_scenario::end(scenario);
//
//         let mut scenario_2 = test_scenario::begin(@0x1000);
//         // the registration is the first object
//         let registration = test_scenario::take_from_sender_by_id<TokenPoolRegistration>(&scenario_2, created[0]);
//         let ctx_2 = scenario_2.ctx();
//
//         let (lock_and_burn, release_and_mint, type_name, execution_state) = registry::get_registration(&registration);
//         assert!(execution_state == EXECUTION_STATE_IDLE);
//         let (lock_and_burn_address, lock_and_burn_module_name, lock_and_burn_function_name) = registry::get_function_info(&lock_and_burn);
//         assert!(lock_and_burn_address == ctx_2.sender());
//         assert!(lock_and_burn_module_name == string::utf8(b"token_pool"));
//         assert!(lock_and_burn_function_name == string::utf8(b"lock_or_burn"));
//         let (release_and_mint_address, release_and_mint_module_name, release_and_mint_function_name) = registry::get_function_info(&release_and_mint);
//         assert!(release_and_mint_address == ctx_2.sender());
//         assert!(release_and_mint_module_name == string::utf8(b"token_pool"));
//         assert!(release_and_mint_function_name == string::utf8(b"release_or_mint"));
//         assert!(type_name.get_address() == ctx_2.sender().to_ascii_string());
//         assert!(type_name.get_module() == string::utf8(b"token_pool").to_ascii());
//
//         registry::set_pool(&mut ref, local_token, &registration, ctx_2);
//
//         registry::transfer_admin_role(&mut ref, local_token, @0x3000, ctx_2);
//         scenario_2.end();
//
//         let mut scenario_3 = test_scenario::begin(@0x3000);
//         let ctx_3 = scenario_3.ctx();
//         registry::accept_admin_role(&mut ref, local_token, ctx_3);
//         assert!(registry::is_administrator(&ref, local_token, @0x3000));
//
//         scenario_3.end();
//         state_object::destroy_state_object(ref);
//         registry::destroy_registration(registration);
//     }
// }