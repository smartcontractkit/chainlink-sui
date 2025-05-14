module ccip::dynamic_dispatcher;

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
