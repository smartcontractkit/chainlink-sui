module ccip::internal {

    const E_TOKEN_ARGUMENTS_MISMATCH: u64 = 1;

    public struct Sui2AnyMessage has drop {
        receiver: vector<u8>,
        data: vector<u8>,
        token_amounts: vector<Sui2AnyTokenAmount>,
        fee_token_metadata: address,
        fee_token_balance: u64,
        extra_args: vector<u8>
    }

    public struct Sui2AnyTokenAmount has drop {
        token: address,
        amount: u64,
    }

    public fun new_sui2any_message(
        receiver: vector<u8>,
        data: vector<u8>,
        token_addresses: vector<address>,
        token_amounts: vector<u64>,
        fee_token_metadata: address,
        fee_token_balance: u64,
        extra_args: vector<u8>
    ): Sui2AnyMessage {
        assert!(
            token_addresses.length() == token_amounts.length(),
            E_TOKEN_ARGUMENTS_MISMATCH
        );
        let mut converted_token_amounts = vector[];
        let tokens_len = token_addresses.length();
        let mut i = 0;
        while (i < tokens_len) {
            let token = token_addresses[i];
            let amount = token_amounts[i];
            converted_token_amounts.push_back(
                Sui2AnyTokenAmount {
                    token,
                    amount
                }
            );
            i = i + 1;
        };
        Sui2AnyMessage {
            receiver,
            data,
            token_amounts: converted_token_amounts,
            fee_token_metadata,
            fee_token_balance,
            extra_args
        }
    }

    public fun get_sui2any_fields(
        message: &Sui2AnyMessage
    ): (vector<u8>, vector<u8>, address, u64, vector<u8>) {
        (
            message.receiver,
            message.data,
            message.fee_token_metadata,
            message.fee_token_balance,
            message.extra_args
        )
    }

    public fun get_sui2any_token_transfers(
        message: &Sui2AnyMessage
    ): (vector<address>, vector<u64>) {
        let mut token_addresses = vector[];
        let mut token_amounts = vector[];
        let len = message.token_amounts.length();
        let mut i = 0;
        while (i < len) {
            let token_amount = &message.token_amounts[i];
            token_addresses.push_back(token_amount.token);
            token_amounts.push_back(token_amount.amount);
            i = i + 1;
        };
        (token_addresses, token_amounts)
    }
}
