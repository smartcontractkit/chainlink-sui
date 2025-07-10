module ccip::dynamic_dispatcher;

use std::type_name;

use ccip::state_object::CCIPObjectRef;
use ccip::token_admin_registry as registry;

const ETypeProofMismatch: u64 = 1;
const EInvalidDestinationChainSelector: u64 = 2;
const EInvalidReceiver: u64 = 3;

public struct DYNAMIC_DISPATCHER has drop {}

/// the cap to be stored in the onramp state to control the source token transfer
public struct SourceTransferCap has key, store {
    id: UID,
}

public struct TokenParams {
    destination_chain_selector: u64,
    receiver: vector<u8>,
    params: vector<SourceTokenTransfer>
}

public struct SourceTokenTransfer has copy, drop {
    // the source token pool package id in SUI
    source_pool: address,
    // the amount of token to transfer
    amount: u64,
    // the source token's coin metadata object id
    source_token_address: address,
    // the destination chain's token address
    dest_token_address: vector<u8>,
    extra_data: vector<u8>,
}

fun init(_witness: DYNAMIC_DISPATCHER, ctx: &mut TxContext) {
    let source_cap = SourceTransferCap {
        id: object::new(ctx),
    };

    transfer::transfer(source_cap, ctx.sender());
}

/// create a new TokenParams object. this should generally be the first step in the PTB flow
/// from onramp side.
public fun create_token_params(destination_chain_selector: u64, receiver: vector<u8>): TokenParams {
    assert!(destination_chain_selector != 0, EInvalidDestinationChainSelector);
    assert!(receiver.length() == 32, EInvalidReceiver);
    TokenParams {
        destination_chain_selector,
        receiver,
        params: vector[]
    }
}

public fun get_destination_chain_selector(token_params: &TokenParams): u64 {
    token_params.destination_chain_selector
}

public fun get_receiver(token_params: &TokenParams): vector<u8> {
    token_params.receiver
}

/// only the token pool with a proper type proof can add the corresponding token transfer.
/// this is not permissioned by a cap because this function is used by token pools in txs
/// signed by a CCIP user.
public fun add_source_token_transfer<TypeProof: drop>(
    ref: &CCIPObjectRef,
    mut token_params: TokenParams,
    amount: u64,
    source_token_address: address,
    dest_token_address: vector<u8>,
    extra_data: vector<u8>,
    _: TypeProof,
): TokenParams {
    let (token_pool_package_id, _, _, _, _, _, type_proof, _, _) = registry::get_token_config(ref, source_token_address);
    let proof_tn = type_name::get<TypeProof>();
    let proof_tn_str = type_name::into_string(proof_tn);
    assert!(type_proof == proof_tn_str, ETypeProofMismatch);
    token_params.params.push_back(
        SourceTokenTransfer {
            source_pool: token_pool_package_id,
            amount,
            source_token_address,
            dest_token_address,
            extra_data,
        }
    );
    token_params
}

/// deconstruct the TokenParams object to get the destination chain selector, receiver, and the list of source token transfers.
/// this is permissioned by the source transfer cap, which is stored in the onramp state.
public fun deconstruct_token_params(_: &SourceTransferCap, token_params: TokenParams): (u64, vector<u8>, vector<SourceTokenTransfer>) {
    let TokenParams {
        destination_chain_selector,
        receiver,
        params
    } = token_params;
    (destination_chain_selector, receiver, params)
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

// =========================== Test Functions =========================== //

#[test_only]
public fun test_init(ctx: &mut TxContext) {
    init(DYNAMIC_DISPATCHER {}, ctx);
}
