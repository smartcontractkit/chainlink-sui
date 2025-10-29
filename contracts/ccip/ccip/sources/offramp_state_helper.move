module ccip::offramp_state_helper;

use ccip::client::{Self, Any2SuiMessage, Any2SuiTokenAmount};
use ccip::ownable::OwnerCap;
use ccip::receiver_registry;
use ccip::state_object::{Self, CCIPObjectRef};
use ccip::token_admin_registry as registry;
use std::ascii;
use std::type_name;
use sui::address;

const ENoMessageToExtract: u64 = 1;
const ETypeProofMismatch: u64 = 2;
const ECCIPReceiveFailed: u64 = 3;
const EWrongReceiptAndTokenTransfer: u64 = 4;
const ETokenTransferMismatch: u64 = 5;
const ETokenTransferAlreadyExists: u64 = 6;
const ETokenTransferDoesNotExist: u64 = 7;
const ETokenTransferAlreadyCompleted: u64 = 8;
const EMessageAlreadyExists: u64 = 9;
const EInvalidOwnerCap: u64 = 10;

public struct OFFRAMP_STATE_HELPER has drop {}

public struct ReceiverParams {
    // if this CCIP message contains token transfers, this vector will be non-empty.
    token_transfer: Option<DestTokenTransfer>,
    // if this CCIP message needs to call a function on the receiver, this will be populated.
    message: Option<Any2SuiMessage>,
    source_chain_selector: u64,
    receipt: Option<CompletedDestTokenTransfer>,
}

/// the cap to be stored in the offramp state to control the updates to ReceiverParams
public struct DestTransferCap has key, store {
    id: UID,
}

public struct CompletedDestTokenTransfer {
    token_receiver: address,
    dest_token_address: address,
}

public struct DestTokenTransfer has copy, drop {
    token_receiver: address,
    remote_chain_selector: u64,
    // the amount of token to transfer, denoted from the source chain
    source_amount: u256,
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

// create a new dest transfer cap if we need to create a new offramp state object
public fun new_dest_transfer_cap(
    ref: &CCIPObjectRef,
    owner_cap: &OwnerCap,
    ctx: &mut TxContext,
): DestTransferCap {
    assert!(object::id(owner_cap) == state_object::owner_cap_id(ref), EInvalidOwnerCap);

    DestTransferCap {
        id: object::new(ctx),
    }
}

public fun create_receiver_params(_: &DestTransferCap, source_chain_selector: u64): ReceiverParams {
    ReceiverParams {
        token_transfer: option::none(),
        message: option::none(),
        source_chain_selector,
        receipt: option::none(),
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
    token_receiver: address,
    remote_chain_selector: u64,
    source_amount: u256,
    dest_token_address: address,
    dest_token_pool_package_id: address,
    source_pool_address: vector<u8>,
    source_pool_data: vector<u8>,
    offchain_data: vector<u8>,
) {
    assert!(receiver_params.token_transfer.is_none(), ETokenTransferAlreadyExists);

    receiver_params
        .token_transfer
        .fill(DestTokenTransfer {
            token_receiver,
            remote_chain_selector,
            source_amount,
            dest_token_address,
            dest_token_pool_package_id,
            source_pool_address,
            source_pool_data,
            offchain_token_data: offchain_data,
        });
}

/// if this CCIP message requires calling a function on a receiver in SUI, this function
/// should be called to populate the message field in the ReceiverParams object.
/// this is permissioned by the DestTransferCap, which is stored in the offramp state.
public fun populate_message(
    _: &DestTransferCap,
    receiver_params: &mut ReceiverParams,
    any2sui_message: Any2SuiMessage,
) {
    assert!(receiver_params.message.is_none(), EMessageAlreadyExists);
    receiver_params.message.fill(any2sui_message);
}

public fun get_dest_token_transfer_data(
    receiver_params: &ReceiverParams,
): (address, u64, u256, address, address, vector<u8>, vector<u8>, vector<u8>) {
    assert!(receiver_params.token_transfer.is_some(), ETokenTransferDoesNotExist);

    let token_transfer = receiver_params.token_transfer.borrow();
    (
        token_transfer.token_receiver,
        token_transfer.remote_chain_selector,
        token_transfer.source_amount,
        token_transfer.dest_token_address,
        token_transfer.dest_token_pool_package_id,
        token_transfer.source_pool_address,
        token_transfer.source_pool_data,
        token_transfer.offchain_token_data,
    )
}

public fun get_token_param_data(
    receiver_params: &ReceiverParams,
): (address, u256, address, vector<u8>, vector<u8>, vector<u8>) {
    assert!(receiver_params.token_transfer.is_some(), ETokenTransferDoesNotExist);
    let token_param = receiver_params.token_transfer.borrow();

    (
        token_param.token_receiver,
        token_param.source_amount,
        token_param.dest_token_address,
        token_param.source_pool_address,
        token_param.source_pool_data, // this is the encoded decimals
        token_param.offchain_token_data,
    )
}

/// only the token pool with a proper type proof can call this function to
/// add a receipt to the receiver params.
public fun complete_token_transfer<TypeProof: drop>(
    ref: &CCIPObjectRef,
    receiver_params: &mut ReceiverParams,
    _: TypeProof,
) {
    let dest_token_transfer = receiver_params.token_transfer.borrow();
    let token_receiver = dest_token_transfer.token_receiver;
    let dest_token_address = dest_token_transfer.dest_token_address;
    let (_, _, _, _, _, type_proof, _, _) = registry::get_token_config_data(
        ref,
        dest_token_address,
    );

    let proof_tn = type_name::with_defining_ids<TypeProof>();
    let proof_tn_str = type_name::into_string(proof_tn);
    assert!(type_proof == proof_tn_str, ETypeProofMismatch);

    let receipt = CompletedDestTokenTransfer {
        token_receiver: token_receiver,
        dest_token_address: dest_token_address,
    };

    assert!(receiver_params.receipt.is_none(), ETokenTransferAlreadyCompleted);

    receiver_params.receipt.fill(receipt);
}

public fun extract_any2sui_message(receiver_params: &mut ReceiverParams): Any2SuiMessage {
    assert!(receiver_params.message.is_some(), ENoMessageToExtract);

    receiver_params.message.extract()
}

public fun new_any2sui_message(
    _: &DestTransferCap,
    message_id: vector<u8>,
    source_chain_selector: u64,
    sender: vector<u8>,
    data: vector<u8>,
    message_receiver: address,
    token_receiver: address,
    token_addresses: vector<address>,
    token_amounts: vector<u256>,
): Any2SuiMessage {
    client::new_any2sui_message(
        message_id,
        source_chain_selector,
        sender,
        data,
        message_receiver,
        token_receiver,
        client::new_dest_token_amounts(token_addresses, token_amounts),
    )
}

public fun consume_any2sui_message<TypeProof: drop>(
    ref: &CCIPObjectRef,
    message: Any2SuiMessage,
    _: TypeProof,
): (vector<u8>, u64, vector<u8>, vector<u8>, address, address, vector<Any2SuiTokenAmount>) {
    let proof_tn = type_name::with_defining_ids<TypeProof>();
    let address_str = type_name::address_string(&proof_tn);
    let receiver_package_id = address::from_ascii_bytes(&ascii::into_bytes(address_str));

    let receiver_config = receiver_registry::get_receiver_config(ref, receiver_package_id);
    let (_, proof_typename) = receiver_registry::get_receiver_config_fields(receiver_config);
    assert!(proof_typename == proof_tn.into_string(), ETypeProofMismatch);

    client::consume_any2sui_message(message, receiver_package_id)
}

/// this function is called by ccip offramp directly, permissioned by the dest transfer cap.
/// it compares token transfers vectors from both hot potatoes and ensures that the message
/// in receiver params is empty.
public fun deconstruct_receiver_params(_: &DestTransferCap, receiver_params: ReceiverParams) {
    let ReceiverParams {
        token_transfer: mut token_transfer_op,
        message: message_op,
        source_chain_selector: _,
        receipt: mut receipt_op,
    } = receiver_params;

    // make sure all token transfers are completed
    assert!(
        token_transfer_op.is_none() && receipt_op.is_none() || (token_transfer_op.is_some() && receipt_op.is_some()),
        EWrongReceiptAndTokenTransfer,
    );
    if (token_transfer_op.is_some()) {
        let token_transfer = token_transfer_op.extract();
        let receipt = receipt_op.extract();
        let CompletedDestTokenTransfer {
            token_receiver,
            dest_token_address,
        } = receipt;

        assert!(
            token_receiver == token_transfer.token_receiver &&
            dest_token_address == token_transfer.dest_token_address,
            ETokenTransferMismatch,
        );
    };

    token_transfer_op.destroy_none();
    receipt_op.destroy_none();

    assert!(message_op.is_none(), ECCIPReceiveFailed);
    message_op.destroy_none();
}

// =========================== Test Functions =========================== //

#[test_only]
public fun test_init(ctx: &mut TxContext) {
    init(OFFRAMP_STATE_HELPER {}, ctx);
}

#[test_only]
public fun deconstruct_receiver_params_with_message_for_test(
    _: &DestTransferCap,
    receiver_package_id: address,
    receiver_params: ReceiverParams,
) {
    let ReceiverParams {
        token_transfer: _,
        message: mut message_op,
        source_chain_selector: _,
        receipt: mut r,
    } = receiver_params;

    if (r.is_some()) {
        let completed_transfer = r.extract();
        let CompletedDestTokenTransfer {
            token_receiver: _,
            dest_token_address: _,
        } = completed_transfer;
    };

    if (message_op.is_some()) {
        let message = message_op.extract();
        client::consume_any2sui_message(message, receiver_package_id);
    };
    message_op.destroy_none();
    r.destroy_none();
}
