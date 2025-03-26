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
        let tokens_len = vector::length(&token_addresses);
        assert!(
            tokens_len == vector::length(&token_amounts),
            E_TOKEN_ARGUMENTS_MISMATCH
        );
        assert!(
            tokens_len == vector::length(&token_store_addresses),
            E_TOKEN_ARGUMENTS_MISMATCH
        );
        let mut converted_token_amounts = vector[];
        let mut i = 0;
        while (i < tokens_len) {
            let token = *vector::borrow(&token_addresses, i);
            let amount = *vector::borrow(&token_amounts, i);
            let token_store = *vector::borrow(&token_store_addresses, i);
            vector::push_back(
                &mut converted_token_amounts,
                Sui2AnyTokenAmount { token, amount, token_store }
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
        let len = vector::length(&message.token_amounts);
        let mut i = 0;
        while (i < len) {
            let token_amount = vector::borrow(&message.token_amounts, i);
            vector::push_back(&mut token_addresses, token_amount.token);
            vector::push_back(&mut token_amounts, token_amount.amount);
            i = i + 1;
        };
        (token_addresses, token_amounts)
    }
}
