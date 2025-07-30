module ccip::onramp_state_helper;

use std::type_name;

use ccip::state_object::CCIPObjectRef;
use ccip::token_admin_registry as registry;

const ETypeProofMismatch: u64 = 1;

public struct ONRAMP_STATE_HELPER has drop {}

/// the cap to be stored in the onramp state to control the source token transfer
public struct SourceTransferCap has key, store {
    id: UID,
}

public struct TokenTransferParams {
    remote_chain_selector: u64,
    source_pool_package_id: address,
    amount: u64,
    source_token_coin_metadata_address: address,
    dest_token_address: vector<u8>,
    extra_data: vector<u8>,
}

fun init(_witness: ONRAMP_STATE_HELPER, ctx: &mut TxContext) {
    let source_cap = SourceTransferCap {
        id: object::new(ctx),
    };

    transfer::transfer(source_cap, ctx.sender());
}

public fun create_token_transfer_params<TypeProof: drop>(
    ref: &CCIPObjectRef,
    remote_chain_selector: u64,
    amount: u64,
    source_token_coin_metadata_address: address,
    dest_token_address: vector<u8>,
    extra_data: vector<u8>,
    _: TypeProof,
 ): TokenTransferParams {
    let (token_pool_package_id, _, _, _, _, type_proof, _, _) = registry::get_token_config(ref, source_token_coin_metadata_address);
    let proof_tn = type_name::get<TypeProof>();
    let proof_tn_str = type_name::into_string(proof_tn);
    assert!(type_proof == proof_tn_str, ETypeProofMismatch);
    TokenTransferParams {
        remote_chain_selector,
        source_pool_package_id: token_pool_package_id,
        amount,
        source_token_coin_metadata_address,
        dest_token_address,
        extra_data,
    }
}

public fun deconstruct_token_params(_: &SourceTransferCap, mut token_transfer_params: vector<TokenTransferParams>) {
    while (token_transfer_params.length() > 0) {
        let TokenTransferParams {
            remote_chain_selector: _,
            source_pool_package_id: _,
            amount: _,
            source_token_coin_metadata_address: _,
            dest_token_address: _,
            extra_data: _,
        } = token_transfer_params.pop_back();
    };
    token_transfer_params.destroy_empty()
}

public fun get_source_token_transfer_data(token_transfer_params: &vector<TokenTransferParams>, index: u64): (u64, address, u64, address, vector<u8>, vector<u8>) {
    (
        token_transfer_params[index].remote_chain_selector,
        token_transfer_params[index].source_pool_package_id,
        token_transfer_params[index].amount,
        token_transfer_params[index].source_token_coin_metadata_address,
        token_transfer_params[index].dest_token_address,
        token_transfer_params[index].extra_data,
    )
}

// =========================== Test Functions =========================== //

#[test_only]
public fun test_init(ctx: &mut TxContext) {
    init(ONRAMP_STATE_HELPER {}, ctx);
}
