module ccip::dynamic_dispatcher;

use std::type_name;

use ccip::state_object::CCIPObjectRef;
use ccip::token_admin_registry as registry;

const ETypeProofMismatch: u64 = 1;

public struct DYNAMIC_DISPATCHER has drop {}

public struct SourceTransferCap has key, store {
    id: UID,
}

public struct TokenParams {
    params: vector<SourceTokenTransfer>
}

public struct SourceTokenTransfer has copy, drop {
    source_pool: address,
    amount: u64,
    source_token_address: address,
    dest_token_address: vector<u8>,
    extra_data: vector<u8>,
}

fun init(_witness: DYNAMIC_DISPATCHER, ctx: &mut TxContext) {
    let source_cap = SourceTransferCap {
        id: object::new(ctx),
    };

    transfer::transfer(source_cap, ctx.sender());
}

public fun create_token_params(): TokenParams {
    TokenParams {
        params: vector[]
    }
}

// only the token pool with a proper type proof can add the corresponding token transfer
public fun add_source_token_transfer<TypeProof: drop>(
    ref: &CCIPObjectRef,
    mut token_params: TokenParams,
    amount: u64,
    source_token_address: address,
    dest_token_address: vector<u8>,
    extra_data: vector<u8>,
    _: TypeProof,
): TokenParams {
    let (source_pool_address, _, _, type_proof) = registry::get_token_config(ref, source_token_address);
    let proof_tn = type_name::get<TypeProof>();
    let proof_tn_str = type_name::into_string(proof_tn);
    assert!(type_proof == proof_tn_str, ETypeProofMismatch);
    token_params.params.push_back(
        SourceTokenTransfer {
            source_pool: source_pool_address,
            amount,
            source_token_address,
            dest_token_address,
            extra_data,
        }
    );
    token_params
}

public fun deconstruct_token_params(_: &SourceTransferCap, token_params: TokenParams): vector<SourceTokenTransfer> {
    let TokenParams {
        params
    } = token_params;
    params
}

public fun get_source_token_transfer_data(token_transfer: SourceTokenTransfer): (address, u64, address, vector<u8>, vector<u8>) {
    (
        token_transfer.source_pool,
        token_transfer.amount,
        token_transfer.source_token_address,
        token_transfer.dest_token_address,
        token_transfer.extra_data,
    )
}
