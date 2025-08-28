module ccip::onramp_state_helper;

use ccip::state_object::CCIPObjectRef;
use ccip::token_admin_registry as registry;
use std::type_name;

const ETypeProofMismatch: u64 = 1;
const ETokenTransferAlreadyExists: u64 = 2;
const ETokenTransferDoesNotExist: u64 = 3;

public struct ONRAMP_STATE_HELPER has drop {}

/// the cap to be stored in the onramp state to control the source token transfer
public struct SourceTransferCap has key, store {
    id: UID,
}

public struct TokenTransferParams {
    token_transfer: Option<TokenTransferMetadata>,
    token_receiver: address,
}

public fun get_token_receiver(params: &TokenTransferParams): address {
    params.token_receiver
}

public struct TokenTransferMetadata {
    remote_chain_selector: u64,
    token_pool_package_id: address,
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

public fun create_token_transfer_params(token_receiver: address): TokenTransferParams {
    TokenTransferParams {
        token_transfer: option::none(),
        token_receiver,
    }
}

/// add a new token transfer to the TokenTransferParams object, which is done within onramp.
/// this is permissioned by the SourceTransferCap, which is stored in the onramp state.
public fun add_token_transfer_param<TypeProof: drop>(
    ref: &CCIPObjectRef,
    token_transfer_params: &mut TokenTransferParams,
    remote_chain_selector: u64,
    amount: u64,
    source_token_coin_metadata_address: address,
    dest_token_address: vector<u8>,
    extra_data: vector<u8>,
    _: TypeProof,
) {
    let token_config = registry::get_token_config(ref, source_token_coin_metadata_address);
    let (token_pool_package_id, _, _, _, _, type_proof, _, _) = registry::get_token_config_data(
        token_config,
    );

    let proof_tn = type_name::get<TypeProof>();
    let proof_tn_str = type_name::into_string(proof_tn);
    assert!(type_proof == proof_tn_str, ETypeProofMismatch);

    assert!(token_transfer_params.token_transfer.is_none(), ETokenTransferAlreadyExists);

    token_transfer_params
        .token_transfer
        .fill(TokenTransferMetadata {
            remote_chain_selector,
            token_pool_package_id,
            amount,
            source_token_coin_metadata_address,
            dest_token_address,
            extra_data,
        })
}

public fun has_token_transfer(token_transfer_params: &TokenTransferParams): bool {
    token_transfer_params.token_transfer.is_some()
}

public fun deconstruct_token_params(
    _: &SourceTransferCap,
    token_transfer_params: TokenTransferParams,
) {
    let TokenTransferParams { token_transfer: mut token_transfer, token_receiver: _ } =
        token_transfer_params;
    if (option::is_some(&token_transfer)) {
        let TokenTransferMetadata {
            remote_chain_selector: _,
            token_pool_package_id: _,
            amount: _,
            source_token_coin_metadata_address: _,
            dest_token_address: _,
            extra_data: _,
        } = token_transfer.extract();
    };
    token_transfer.destroy_none();
}

public fun get_source_token_transfer_data(
    token_transfer_params: &TokenTransferParams,
): (u64, address, u64, address, vector<u8>, vector<u8>) {
    assert!(token_transfer_params.token_transfer.is_some(), ETokenTransferDoesNotExist);
    let token_transfer = token_transfer_params.token_transfer.borrow();
    (
        token_transfer.remote_chain_selector,
        token_transfer.token_pool_package_id,
        token_transfer.amount,
        token_transfer.source_token_coin_metadata_address,
        token_transfer.dest_token_address,
        token_transfer.extra_data,
    )
}

// =========================== Test Functions =========================== //

#[test_only]
public fun test_init(ctx: &mut TxContext) {
    init(ONRAMP_STATE_HELPER {}, ctx);
}
