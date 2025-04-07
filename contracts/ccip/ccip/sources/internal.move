module ccip::internal {

    const E_TOKEN_ARGUMENTS_MISMATCH: u64 = 1;

    public struct Sui2AnyMessage has drop {
        receiver: vector<u8>,
        data: vector<u8>,
        token_amounts: vector<Sui2AnyTokenAmount>,
        fee_token: address,
        fee_token_store: address,
        extra_args: vector<u8>
    }

    public struct Sui2AnyTokenAmount has drop {
        token: address,
        amount: u64,
        token_store: address
    }

    public fun new_sui2any_message(
        receiver: vector<u8>,
        data: vector<u8>,
        token_addresses: vector<address>,
        token_amounts: vector<u64>,
        token_store_addresses: vector<address>,
        fee_token: address,
        fee_token_store: address,
        extra_args: vector<u8>
    ): Sui2AnyMessage {
        let tokens_len = token_addresses.length();
        assert!(
            tokens_len == token_amounts.length(),
            E_TOKEN_ARGUMENTS_MISMATCH
        );
        assert!(
            tokens_len == token_store_addresses.length(),
            E_TOKEN_ARGUMENTS_MISMATCH
        );
        let mut converted_token_amounts = vector[];
        let mut i = 0;
        while (i < tokens_len) {
            let token = token_addresses[i];
            let amount = token_amounts[i];
            let token_store = token_store_addresses[i];
            converted_token_amounts.push_back(
                Sui2AnyTokenAmount {
                    token,
                    amount,
                    token_store
                }
            );
            i = i + 1;
        };
        Sui2AnyMessage {
            receiver,
            data,
            token_amounts: converted_token_amounts,
            fee_token,
            fee_token_store,
            extra_args
        }
    }

    /// Returns all fields except for FungibleAsset
    public fun get_sui2any_fields(
        message: &Sui2AnyMessage
    ): (vector<u8>, vector<u8>, address, address, vector<u8>) {
        (
            message.receiver,
            message.data,
            message.fee_token,
            message.fee_token_store,
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
