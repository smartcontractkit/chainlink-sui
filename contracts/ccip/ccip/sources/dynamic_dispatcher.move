module ccip::dynamic_dispatcher;

use ccip::client;

const E_WRONG_INDEX_IN_RECEIVER_PARAMS: u64 = 1;
const E_TOKEN_TRANSFER_ALREADY_COMPLETED: u64 = 2;
const E_TOKEN_POOL_ADDRESS_MISMATCH: u64 = 3;

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

public fun create_token_params(): TokenParams {
    TokenParams {
        params: vector[]
    }
}

public fun add_token_param(
    mut token_params: TokenParams,
    source_pool: address,
    amount: u64,
    source_token_address: address,
    dest_token_address: vector<u8>,
    extra_data: vector<u8>,
): TokenParams {
    token_params.params.push_back(
        SourceTokenTransfer {
            source_pool,
            amount,
            source_token_address,
            dest_token_address,
            extra_data, // encoded decimals
        }
    );
    token_params
}

public fun deconstruct_token_params(token_params: TokenParams): vector<SourceTokenTransfer> {
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

public struct ReceiverParams {
    params: vector<DestTokenTransfer>,
    message: Option<client::Any2SuiMessage>,
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
    offchain_token_data: vector<u8>, // not used?
    completed: bool
}

public fun get_completed(transfer: DestTokenTransfer): bool {
    transfer.completed
}

public fun create_receiver_params(): ReceiverParams {
    ReceiverParams {
        params: vector[],
        message: option::none(),
    }
}

public fun add_dest_token_transfer(
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
        E_WRONG_INDEX_IN_RECEIVER_PARAMS
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

// called by token pool to mark token transfers as completed
public fun complete_token_transfer(
    mut receiver_params: ReceiverParams,
    index: u64,
    local_amount: u64,
    token_pool_address: address
): ReceiverParams {
    assert!(
        index < receiver_params.params.length(),
        E_WRONG_INDEX_IN_RECEIVER_PARAMS
    );
    assert!(
        !receiver_params.params[index].completed,
        E_TOKEN_TRANSFER_ALREADY_COMPLETED
    );
    assert!(
        receiver_params.params[index].token_pool_address == token_pool_address,
        E_TOKEN_POOL_ADDRESS_MISMATCH
    );
    receiver_params.params[index].completed = true;
    receiver_params.params[index].local_amount = local_amount;

    receiver_params
}

// called by ccip receiver directly, or by PTB to extract the message and send to the receiver
public fun extract_any2sui_message(
    mut receiver_params: ReceiverParams
): (Option<client::Any2SuiMessage>, ReceiverParams) {
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
