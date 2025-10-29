module ccip::client;

use std::bcs;

const GENERIC_EXTRA_ARGS_V2_TAG: vector<u8> = x"181dcf10";
const SVM_EXTRA_ARGS_V1_TAG: vector<u8> = x"1f3b3aba";
const SUI_EXTRA_ARGS_V1_TAG: vector<u8> = x"21ea4ca9";

const EInvalidSVMTokenReceiverLength: u64 = 1;
const EInvalidSVMAccountLength: u64 = 2;
const EInvalidSUITokenReceiverLength: u64 = 3;
const EInvalidSUIReceiverObjectIdLength: u64 = 4;
const ETokenAmountMismatch: u64 = 5;
const EMessageReceiverMismatch: u64 = 6;

public fun generic_extra_args_v2_tag(): vector<u8> {
    GENERIC_EXTRA_ARGS_V2_TAG
}

public fun svm_extra_args_v1_tag(): vector<u8> {
    SVM_EXTRA_ARGS_V1_TAG
}

public fun sui_extra_args_v1_tag(): vector<u8> {
    SUI_EXTRA_ARGS_V1_TAG
}

public fun encode_sui_extra_args_v1(
    gas_limit: u64,
    allow_out_of_order_execution: bool,
    token_receiver: vector<u8>,
    receiver_object_ids: vector<vector<u8>>,
): vector<u8> {
    let mut extra_args = vector[];
    extra_args.append(SUI_EXTRA_ARGS_V1_TAG);
    extra_args.append(bcs::to_bytes(&gas_limit));
    extra_args.append(bcs::to_bytes(&allow_out_of_order_execution));

    assert!(token_receiver.length() == 32, EInvalidSUITokenReceiverLength);
    extra_args.append(bcs::to_bytes(&token_receiver));

    receiver_object_ids.do_ref!(|receiver_object_id| {
        assert!(receiver_object_id.length() == 32, EInvalidSUIReceiverObjectIdLength);
    });
    extra_args.append(bcs::to_bytes(&receiver_object_ids));
    extra_args
}

public fun encode_generic_extra_args_v2(
    gas_limit: u256,
    allow_out_of_order_execution: bool,
): vector<u8> {
    let mut extra_args = vector[];
    extra_args.append(GENERIC_EXTRA_ARGS_V2_TAG);
    extra_args.append(bcs::to_bytes(&gas_limit));
    extra_args.append(bcs::to_bytes(&allow_out_of_order_execution));
    extra_args
}

public fun encode_svm_extra_args_v1(
    compute_units: u32,
    account_is_writable_bitmap: u64,
    allow_out_of_order_execution: bool,
    token_receiver: vector<u8>,
    accounts: vector<vector<u8>>,
): vector<u8> {
    let mut extra_args = vector[];
    extra_args.append(SVM_EXTRA_ARGS_V1_TAG);
    extra_args.append(bcs::to_bytes(&compute_units));
    extra_args.append(bcs::to_bytes(&account_is_writable_bitmap));
    extra_args.append(bcs::to_bytes(&allow_out_of_order_execution));

    assert!(token_receiver.length() == 32, EInvalidSVMTokenReceiverLength);
    // Check that all accounts have length 32
    accounts.do_ref!(|account| {
        assert!(account.length() == 32, EInvalidSVMAccountLength);
    });

    extra_args.append(bcs::to_bytes(&token_receiver));
    extra_args.append(bcs::to_bytes(&accounts));
    extra_args
}

public struct Any2SuiMessage {
    message_id: vector<u8>,
    source_chain_selector: u64,
    sender: vector<u8>,
    data: vector<u8>,
    message_receiver: address,
    token_receiver: address,
    dest_token_amounts: vector<Any2SuiTokenAmount>,
}

public struct Any2SuiTokenAmount has copy, drop {
    token: address,
    amount: u256,
}

public(package) fun new_any2sui_message(
    message_id: vector<u8>,
    source_chain_selector: u64,
    sender: vector<u8>,
    data: vector<u8>,
    message_receiver: address,
    token_receiver: address,
    dest_token_amounts: vector<Any2SuiTokenAmount>,
): Any2SuiMessage {
    Any2SuiMessage {
        message_id,
        source_chain_selector,
        sender,
        data,
        message_receiver,
        token_receiver,
        dest_token_amounts,
    }
}

public(package) fun consume_any2sui_message(
    message: Any2SuiMessage,
    receiver_package_id: address,
): (vector<u8>, u64, vector<u8>, vector<u8>, address, address, vector<Any2SuiTokenAmount>) {
    let Any2SuiMessage {
        message_id,
        source_chain_selector,
        sender,
        data,
        message_receiver,
        token_receiver,
        dest_token_amounts,
    } = message;
    assert!(message_receiver == receiver_package_id, EMessageReceiverMismatch);

    (
        message_id,
        source_chain_selector,
        sender,
        data,
        message_receiver,
        token_receiver,
        dest_token_amounts,
    )
}

public(package) fun new_dest_token_amounts(
    token_addresses: vector<address>,
    token_amounts: vector<u256>,
): vector<Any2SuiTokenAmount> {
    assert!(token_addresses.length() == token_amounts.length(), ETokenAmountMismatch);
    token_addresses.zip_map_ref!(&token_amounts, |token_address, token_amount| {
        Any2SuiTokenAmount { token: *token_address, amount: *token_amount }
    })
}

public fun get_token(input: &Any2SuiTokenAmount): address {
    input.token
}

public fun get_amount(input: &Any2SuiTokenAmount): u256 {
    input.amount
}

public fun get_token_and_amount(input: &Any2SuiTokenAmount): (address, u256) {
    (input.token, input.amount)
}
