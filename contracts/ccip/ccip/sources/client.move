module ccip::client {

    public struct Any2SuiTokenAmount has store, drop, copy {
        token: address,
        amount: u64
    }

    public struct Any2SuiMessage has store, drop, copy {
        message_id: vector<u8>,
        source_chain_selector: u64,
        sender: vector<u8>,
        data: vector<u8>,
        dest_token_amounts: vector<Any2SuiTokenAmount>
    }

    public fun get_source_chain_selector(input: &Any2SuiMessage): u64 {
        input.source_chain_selector
    }

    public fun get_message_id(input: &Any2SuiMessage): vector<u8> {
        input.message_id
    }

    public fun get_sender(input: &Any2SuiMessage): vector<u8> {
        input.sender
    }

    public fun get_data(input: &Any2SuiMessage): vector<u8> {
        input.data
    }

    public fun get_token(input: &Any2SuiTokenAmount): address {
        input.token
    }

    public fun get_amount(input: &Any2SuiTokenAmount): u64 {
        input.amount
    }

    public fun get_dest_token_amounts(input: &Any2SuiMessage): vector<Any2SuiTokenAmount> {
        input.dest_token_amounts
    }

    public fun new_any2sui_message(
        message_id: vector<u8>,
        source_chain_selector: u64,
        sender: vector<u8>,
        data: vector<u8>,
        dest_token_amounts: vector<Any2SuiTokenAmount>
    ): Any2SuiMessage {
        Any2SuiMessage {
            message_id,
            source_chain_selector,
            sender,
            data,
            dest_token_amounts
        }
    }

    public fun new_dest_token_amounts(
        token_addresses: vector<address>, token_amounts: vector<u64>
    ): vector<Any2SuiTokenAmount> {
        vector::zip_map_ref!(
            &token_addresses,
            &token_amounts,
            |token_address, token_amount| {
                Any2SuiTokenAmount { token: *token_address, amount: *token_amount }
            }
        )
    }
}
