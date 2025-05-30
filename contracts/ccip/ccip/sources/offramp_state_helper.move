module ccip::offramp_state_helper;

use std::type_name;
use ccip::receiver_registry;

use ccip::client;
use ccip::state_object::CCIPObjectRef;
use ccip::token_admin_registry as registry;

const EWrongIndexInReceiverParams: u64 = 1;
const ETokenTransferAlreadyCompleted: u64 = 2;
const ETokenPoolAddressMismatch: u64 = 3;
const ETypeProofMismatch: u64 = 4;

public struct OFFRAMP_STATE_HELPER has drop {}

public struct ReceiverParams {
    params: vector<DestTokenTransfer>,
    message: Option<client::Any2SuiMessage>,
}

public struct DestTransferCap has key, store {
    id: UID,
}

public struct DestTokenTransfer has copy, drop {
    // sender: vector<u8>,
    receiver: address,
    source_amount: u64,
    local_amount: u64,
    // source_chain_selector: u64,
    dest_token_address: address,
    token_pool_address: address,
    source_pool_address: vector<u8>,
    source_pool_data: vector<u8>,
    offchain_token_data: vector<u8>,
    completed: bool
}

fun init(_witness: OFFRAMP_STATE_HELPER, ctx: &mut TxContext) {
    let dest_cap = DestTransferCap {
        id: object::new(ctx),
    };

    transfer::transfer(dest_cap, ctx.sender());
}

public fun get_completed(transfer: DestTokenTransfer): bool {
    transfer.completed
}

public fun create_receiver_params(_: &DestTransferCap): ReceiverParams {
    ReceiverParams {
        params: vector[],
        message: option::none(),
    }
}

public fun add_dest_token_transfer(
    _: &DestTransferCap,
    receiver_params: &mut ReceiverParams,
    receiver: address,
    source_amount: u64,
    dest_token_address: address,
    token_pool_address: address,
    source_pool_address: vector<u8>,
    source_pool_data: vector<u8>,
    offchain_data: vector<u8>,
) {
    receiver_params.params.push_back(
        DestTokenTransfer {
            // sender: message.sender,
            receiver,
            source_amount,
            local_amount: 0, // to be calculated by the destination token pool
            // source_chain_selector: message.header.source_chain_selector,
            dest_token_address,
            token_pool_address,
            source_pool_address,
            source_pool_data,
            offchain_token_data: offchain_data,
            completed: false,
        }
    );
}

public fun populate_message(
    _: &DestTransferCap,
    receiver_params: &mut ReceiverParams,
    any2sui_message: client::Any2SuiMessage,
) {
    receiver_params.message.fill(any2sui_message);
}

public fun get_token_param_data(
    receiver_params: &ReceiverParams, index: u64
): (address, u64, address, vector<u8>, vector<u8>) {
    assert!(
        index < receiver_params.params.length(),
        EWrongIndexInReceiverParams
    );
    let token_param = receiver_params.params[index];

    (
        // token_param.sender,
        token_param.receiver,
        token_param.source_amount,
        token_param.dest_token_address,
        token_param.source_pool_address,
        token_param.source_pool_data, // this is the encoded decimals
    )
}

// only the token pool with a proper type proof can mark the corresponding token transfer as completed
public fun complete_token_transfer<TypeProof: drop>(
    ref: &CCIPObjectRef,
    mut receiver_params: ReceiverParams,
    index: u64,
    local_amount: u64,
    _: TypeProof,
): ReceiverParams {
    assert!(
        index < receiver_params.params.length(),
        EWrongIndexInReceiverParams,
    );

    let param = receiver_params.params[index];
    assert!(!param.completed, ETokenTransferAlreadyCompleted);
    let (token_pool_package_id, _, _, _, _, _, type_proof) = registry::get_token_config(ref, param.dest_token_address);
    assert!(
        param.token_pool_address == token_pool_package_id,
        ETokenPoolAddressMismatch,
    );
    let proof_tn = type_name::get<TypeProof>();
    let proof_tn_str = type_name::into_string(proof_tn);
    assert!(type_proof == proof_tn_str, ETypeProofMismatch);

    receiver_params.params[index].completed = true;
    receiver_params.params[index].local_amount = local_amount;

    receiver_params
}

// called by ccip receiver directly, or by PTB to extract the message and send to the receiver
public fun extract_any2sui_message<TypeProof: drop>(
    ref: &CCIPObjectRef,
    mut receiver_params: ReceiverParams,
    package_id: address,
    _: TypeProof,
): (Option<client::Any2SuiMessage>, ReceiverParams) {
    let receiver_config = receiver_registry::get_receiver_config(ref, package_id);
    let proof_tn = type_name::get<TypeProof>();
    let (_, _, _, proof_typename) = receiver_registry::get_receiver_config_fields(receiver_config);
    assert!(
        proof_typename == proof_tn,
        ETypeProofMismatch
    );

    let message = receiver_params.message;
    receiver_params.message = option::none();

    (message, receiver_params)
}

public fun deconstruct_receiver_params(receiver_params: ReceiverParams): (vector<DestTokenTransfer>, Option<client::Any2SuiMessage>) {
    let ReceiverParams {
        params,
        message,
    } = receiver_params;
    (params, message)
}