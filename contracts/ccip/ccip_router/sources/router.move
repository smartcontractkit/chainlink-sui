module ccip_router::router {
    use std::string::{Self, String};

    use sui::clock::Clock;

    use ccip::onramp;
    use ccip::state_object::CCIPObjectRef;

    public fun type_and_version(): String {
        string::utf8(b"Router 1.6.0")
    }

    public fun is_chain_supported(ref: &CCIPObjectRef, dest_chain_selector: u64): bool {
        onramp::is_chain_supported(ref, dest_chain_selector)
    }

    public fun get_fee(
        ref: &CCIPObjectRef,
        clock: &Clock,
        dest_chain_selector: u64,
        receiver: vector<u8>,
        data: vector<u8>,
        token_addresses: vector<address>,
        token_amounts: vector<u64>,
        token_store_addresses: vector<address>,
        fee_token: address,
        fee_token_store: address,
        extra_args: vector<u8>
    ): u64 {
        onramp::get_fee(
            ref,
            clock,
            dest_chain_selector,
            receiver,
            data,
            token_addresses,
            token_amounts,
            token_store_addresses,
            fee_token,
            fee_token_store,
            extra_args
        )
    }

    public fun ccip_send(
        ref: &mut CCIPObjectRef,
        clock: &Clock,
        dest_chain_selector: u64,
        receiver: vector<u8>,
        data: vector<u8>,
        token_addresses: vector<address>,
        token_amounts: vector<u64>,
        token_store_addresses: vector<address>,
        fee_token: address,
        fee_token_store: address,
        extra_args: vector<u8>,
        ctx: &mut TxContext
    ) {
        onramp::ccip_send(
            ref,
            clock,
            dest_chain_selector,
            receiver,
            data,
            token_addresses,
            token_amounts,
            token_store_addresses,
            fee_token,
            fee_token_store,
            extra_args,
            ctx
        );
    }

    public fun ccip_send_with_message_id(
        ref: &mut CCIPObjectRef,
        clock: &Clock,
        dest_chain_selector: u64,
        receiver: vector<u8>,
        data: vector<u8>,
        token_addresses: vector<address>,
        token_amounts: vector<u64>,
        token_store_addresses: vector<address>,
        fee_token: address,
        fee_token_store: address,
        extra_args: vector<u8>,
        ctx: &mut TxContext
    ): vector<u8> {
        onramp::ccip_send(
            ref,
            clock,
            dest_chain_selector,
            receiver,
            data,
            token_addresses,
            token_amounts,
            token_store_addresses,
            fee_token,
            fee_token_store,
            extra_args,
            ctx
        )
    }
}
