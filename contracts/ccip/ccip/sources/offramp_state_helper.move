module ccip::offramp_state_helper;

use std::ascii;
use std::type_name;

use sui::address;

use ccip::client::{Self, Any2SuiMessage, Any2SuiTokenAmount};
use ccip::receiver_registry;
use ccip::state_object::CCIPObjectRef;
use ccip::token_admin_registry as registry;

const EWrongIndexInReceiverParams: u64 = 1;
const ENoMessageToExtract: u64 = 2;
const ETypeProofMismatch: u64 = 3;
const ECCIPReceiveFailed: u64 = 4;
const EWrongNumberOfCompletedTransfers: u64 = 5;
const ETokenTransferMismatch: u64 = 6;

public struct OFFRAMP_STATE_HELPER has drop {}

public struct ReceiverParams {
    // if this CCIP message contains token transfers, this vector will be non-empty.
    params: vector<DestTokenTransfer>,
    // if this CCIP message needs to call a function on the receiver, this will be populated.
    message: Option<Any2SuiMessage>,
    source_chain_selector: u64,
    receipts: vector<CompletedDestTokenTransfer>,
}

/// the cap to be stored in the offramp state to control the updates to ReceiverParams
public struct DestTransferCap has key, store {
    id: UID,
}

public struct CompletedDestTokenTransfer {
    token_pool_id: ID,
    receiver: address,
    dest_token_address: address,
}

public struct DestTokenTransfer has copy, drop {
    receiver: address,
    remote_chain_selector: u64,
    // the amount of token to transfer, denoted from the source chain
    source_amount: u64,
    // the token's coin metadata object id on SUI
    dest_token_address: address,
    // the destination token pool package id on SUI
    dest_token_pool_package_id: address,
    // the source pool address on the source chain
    source_pool_address: vector<u8>,
    source_pool_data: vector<u8>,
    offchain_token_data: vector<u8>,
}

fun init(_witness: OFFRAMP_STATE_HELPER, ctx: &mut TxContext) {
    let dest_cap = DestTransferCap {
        id: object::new(ctx),
    };

    transfer::transfer(dest_cap, ctx.sender());
}

public fun create_receiver_params(_: &DestTransferCap, source_chain_selector: u64): ReceiverParams {
    ReceiverParams {
        params: vector[],
        message: option::none(),
        source_chain_selector,
        receipts: vector[],
    }
}

public fun get_source_chain_selector(receiver_params: &ReceiverParams): u64 {
    receiver_params.source_chain_selector
}

/// add a new token transfer to the ReceiverParams object, which is done within offramp.
/// this is permissioned by the DestTransferCap, which is stored in the offramp state.
public fun add_dest_token_transfer(
    _: &DestTransferCap,
    receiver_params: &mut ReceiverParams,
    receiver: address,
    remote_chain_selector: u64,
    source_amount: u64,
    dest_token_address: address,
    dest_token_pool_package_id: address,
    source_pool_address: vector<u8>,
    source_pool_data: vector<u8>,
    offchain_data: vector<u8>,
) {
    receiver_params.params.push_back(
        DestTokenTransfer {
            receiver,
            remote_chain_selector,
            source_amount,
            dest_token_address,
            dest_token_pool_package_id,
            source_pool_address,
            source_pool_data,
            offchain_token_data: offchain_data,
        }
    );
}

/// if this CCIP message requires calling a function on a receiver in SUI, this function
/// should be called to populate the message field in the ReceiverParams object.
/// this is permissioned by the DestTransferCap, which is stored in the offramp state.
public fun populate_message(
    _: &DestTransferCap,
    receiver_params: &mut ReceiverParams,
    any2sui_message: Any2SuiMessage,
) {
    receiver_params.message.fill(any2sui_message);
}

public fun create_dest_token_transfer(
    receiver: address,
    remote_chain_selector: u64,
    source_amount: u64,
    dest_token_address: address,
    dest_token_pool_package_id: address,
    source_pool_address: vector<u8>,
    source_pool_data: vector<u8>,
    offchain_token_data: vector<u8>,
): DestTokenTransfer {
    DestTokenTransfer {
        receiver,
        remote_chain_selector,
        source_amount,
        dest_token_address,
        dest_token_pool_package_id,
        source_pool_address,
        source_pool_data,
        offchain_token_data,
    }
}

public fun get_dest_token_transfer_data(
    d: DestTokenTransfer,
): (
    address,
    u64,
    u64,
    address,
    address,
    vector<u8>,
    vector<u8>,
    vector<u8>,
) {
    (
        d.receiver,
        d.remote_chain_selector,
        d.source_amount,
        d.dest_token_address,
        d.dest_token_pool_package_id,
        d.source_pool_address,
        d.source_pool_data,
        d.offchain_token_data,
    )
}

public fun get_dest_token_transfer(
    receiver_params: &ReceiverParams,
    index: u64,
): DestTokenTransfer {
    assert!(
        index < receiver_params.params.length(),
        EWrongIndexInReceiverParams
    );

    receiver_params.params[index]
}

public fun get_token_param_data(
    receiver_params: &ReceiverParams,
    index: u64,
): (
    address,
    u64,
    address,
    vector<u8>,
    vector<u8>,
    vector<u8>,
) {
    assert!(
        index < receiver_params.params.length(),
        EWrongIndexInReceiverParams
    );
    let token_param = receiver_params.params[index];

    (
        token_param.receiver,
        token_param.source_amount,
        token_param.dest_token_address,
        token_param.source_pool_address,
        token_param.source_pool_data, // this is the encoded decimals
        token_param.offchain_token_data,
    )
}

/// only the token pool with a proper type proof can call this function to
/// return a CompletedDestTokenTransfer object.
public fun complete_token_transfer<TypeProof: drop>(
    ref: &CCIPObjectRef,
    receiver_params: &mut ReceiverParams,
    receiver: address,
    dest_token_address: address,
    token_pool_id: ID,
    _: TypeProof,
) {
    let token_config = registry::get_token_config(ref, dest_token_address);
    let (_, _, _, _, _, type_proof, _, _) = registry::get_token_config_data(token_config);

    let proof_tn = type_name::get<TypeProof>();
    let proof_tn_str = type_name::into_string(proof_tn);
    assert!(type_proof == proof_tn_str, ETypeProofMismatch);

    let receipt = CompletedDestTokenTransfer {
        token_pool_id: token_pool_id,
        receiver: receiver,
        dest_token_address: dest_token_address,
    };

    receiver_params.receipts.push_back(receipt);
}

public fun extract_any2sui_message(
    receiver_params: &mut ReceiverParams,
): Any2SuiMessage {
    assert!(
        receiver_params.message.is_some(),
        ENoMessageToExtract
    );

    receiver_params.message.extract()
}

public fun consume_any2sui_message<TypeProof: drop>(
    ref: &CCIPObjectRef,
    message: Any2SuiMessage,
    _: TypeProof,
): (
    vector<u8>,
    u64,
    vector<u8>,
    vector<u8>,
    vector<Any2SuiTokenAmount>,
) {
    let proof_tn = type_name::get<TypeProof>();
    let address_str = type_name::get_address(&proof_tn);
    let receiver_package_id = address::from_ascii_bytes(&ascii::into_bytes(address_str));

    let receiver_config = receiver_registry::get_receiver_config(ref, receiver_package_id);
    let (_, proof_typename) = receiver_registry::get_receiver_config_fields(receiver_config);
    assert!(
        proof_typename == proof_tn.into_string(),
        ETypeProofMismatch
    );

    client::consume_any2sui_message(message)
}

/// this function is called by ccip offramp directly, permissioned by the dest transfer cap.
/// it compares token transfers vectors from both hot potatoes and ensures that the message
/// in receiver params is empty.
public fun deconstruct_receiver_params(
    _: &DestTransferCap,
    receiver_params: ReceiverParams,
) {
    let ReceiverParams {
        mut params,
        message: message_op,
        source_chain_selector: _,
        receipts: mut completed_transfers,
    } = receiver_params;

    // make sure all token transfers are completed
    assert!(completed_transfers.length() == params.length(), EWrongNumberOfCompletedTransfers);
    while (!completed_transfers.is_empty()) {
        let completed_transfer = completed_transfers.pop_back();
        let CompletedDestTokenTransfer {
            token_pool_id: _,
            receiver,
            dest_token_address,
        } = completed_transfer;
        let token_transfer = params.pop_back();
        assert!(
            receiver == token_transfer.receiver &&
            dest_token_address == token_transfer.dest_token_address,
            ETokenTransferMismatch
        );
    };
    completed_transfers.destroy_empty();

    // make sure the any2sui message is extracted
    assert!(message_op.is_none(), ECCIPReceiveFailed);
    message_op.destroy_none();
}

// =========================== Test Functions =========================== //

#[test_only]
public fun test_init(ctx: &mut TxContext) {
    init(OFFRAMP_STATE_HELPER {}, ctx);
}

#[test_only]
public fun deconstruct_receiver_params_for_test(
    _: &DestTransferCap,
    receiver_params: ReceiverParams,
) {
    let ReceiverParams {
        params: _,
        message: message_op,
        source_chain_selector: _,
        receipts: mut r,
    } = receiver_params;

    while (r.is_empty() == false) {
        let completed_transfer = r.pop_back();
        let CompletedDestTokenTransfer {
            token_pool_id: _,
            receiver: _,
            dest_token_address: _,
        } = completed_transfer;
    };

    assert!(message_op.is_none(), ECCIPReceiveFailed);
    message_op.destroy_none();
    r.destroy_empty();
}

#[test_only]
public fun deconstruct_receiver_params_with_message_for_test(
    _: &DestTransferCap,
    receiver_params: ReceiverParams,
) {
    let ReceiverParams {
        params: _,
        message: mut message_op,
        source_chain_selector: _,
        receipts: mut r,
    } = receiver_params;

    while (r.is_empty() == false) {
        let completed_transfer = r.pop_back();
        let CompletedDestTokenTransfer {
            token_pool_id: _,
            receiver: _,
            dest_token_address: _,
        } = completed_transfer;
    };

    if (message_op.is_some()) {
        let message = message_op.extract();
        client::consume_any2sui_message(message);
    };
    message_op.destroy_none();
    r.destroy_empty();
}
