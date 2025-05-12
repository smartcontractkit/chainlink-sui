module ccip_router::router;

use std::string::{Self, String};

use sui::clock::Clock;
use sui::coin::{Coin, CoinMetadata};

use ccip::state_object::CCIPObjectRef;

use ccip_onramp::onramp::{Self, OnRampState};

use dynamic_dispatcher::dynamic_dispatcher as dd;

public fun type_and_version(): String {
    string::utf8(b"Router 1.6.0")
}

public fun is_chain_supported(state: &OnRampState, dest_chain_selector: u64): bool {
    onramp::is_chain_supported(state, dest_chain_selector)
}

public fun get_fee<T>(
    ref: &CCIPObjectRef,
    clock: &Clock,
    dest_chain_selector: u64,
    receiver: vector<u8>,
    data: vector<u8>,
    token_addresses: vector<address>,
    token_amounts: vector<u64>,
    fee_token: &CoinMetadata<T>,
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
        fee_token,
        extra_args
    )
}

// ccip_send does not have a return value. EOA calls cannot receive a return value.
public fun ccip_send<T>(
    ref: &mut CCIPObjectRef,
    state: &mut OnRampState,
    clock: &Clock,
    dest_chain_selector: u64,
    receiver: vector<u8>,
    data: vector<u8>,
    token_params: dd::TokenParams,
    fee_token_metadata: &CoinMetadata<T>,
    fee_token: Coin<T>,
    extra_args: vector<u8>,
    ctx: &mut TxContext
) {
    onramp::ccip_send(
        ref,
        state,
        clock,
        dest_chain_selector,
        receiver,
        data,
        token_params,
        fee_token_metadata,
        fee_token,
        extra_args,
        ctx,
    );
}

// ccip_send_with_message_id has a return value. Contract calls can receive a return value.
public fun ccip_send_with_message_id<T>(
    ref: &mut CCIPObjectRef,
    state: &mut OnRampState,
    clock: &Clock,
    dest_chain_selector: u64,
    receiver: vector<u8>,
    data: vector<u8>,
    token_params: dd::TokenParams,
    fee_token_metadata: &CoinMetadata<T>,
    fee_token: Coin<T>,
    extra_args: vector<u8>,
    ctx: &mut TxContext
): vector<u8> {
    onramp::ccip_send(
        ref,
        state,
        clock,
        dest_chain_selector,
        receiver,
        data,
        token_params,
        fee_token_metadata,
        fee_token,
        extra_args,
        ctx,
    )
}
