module ccip_dummy_receiver::dummy_receiver;

use std::string::{Self, String};

use sui::event;
use sui::clock::Clock;
use sui::coin::Coin;
use sui::transfer::Receiving;

use ccip::client;
use ccip::receiver_registry;
use ccip::state_object::CCIPObjectRef;
use ccip::offramp_state_helper::{Self as osh};

const EMessageIdMismatch: u64 = 0;

public struct OwnerCap has key, store {
    id: UID,
    receiver_address: address,
}

public struct ReceivedMessage has copy, drop {
    message_id: vector<u8>,
    source_chain_selector: u64,
    sender: vector<u8>,
    data: vector<u8>,
    dest_token_transfer_length: u64,
    dest_token_amounts: vector<TokenAmount>,
}

public struct CCIPReceiverState has key {
    id: UID,
    counter: u64,
    message_id: vector<u8>,
    source_chain_selector: u64,
    sender: vector<u8>,
    data: vector<u8>,
    dest_token_transfer_length: u64,
    dest_token_amounts: vector<TokenAmount>,
}

public struct DummyReceiverProof has drop {}

public struct TokenAmount has copy, drop, store {
    token: address,
    amount: u64,
}

public fun type_and_version(): String {
    string::utf8(b"DummyReceiver 1.6.0")
}

fun init(ctx: &mut TxContext) {
    let state = CCIPReceiverState {
        id: object::new(ctx),
        counter: 0,
        message_id: vector[],
        source_chain_selector: 0,
        sender: vector[],
        data: vector[],
        dest_token_transfer_length: 0,
        dest_token_amounts: vector[],
    };

    let owner_cap = OwnerCap {
        id: object::new(ctx),
        receiver_address: object::id_to_address(object::borrow_id(&state)),
    };

    transfer::share_object(state);
    transfer::transfer(owner_cap, ctx.sender());
}

public fun register_receiver(ref: &mut CCIPObjectRef) {
    receiver_registry::register_receiver(ref, DummyReceiverProof {});
}

public fun get_counter(state: &CCIPReceiverState): u64 {
    state.counter
}

public fun get_dest_token_amounts(state: &CCIPReceiverState): vector<TokenAmount> {
    state.dest_token_amounts
}

public fun get_token_amount_token(token_amount: &TokenAmount): address {
    token_amount.token
}

public fun get_token_amount_amount(token_amount: &TokenAmount): u64 {
    token_amount.amount
}

public fun receive_and_send_coin<T>(state: &mut CCIPReceiverState, _: &OwnerCap, coin_receiving: Receiving<Coin<T>>, recipient: address) {
    let c = transfer::public_receive<Coin<T>>(&mut state.id, coin_receiving);
    transfer::public_transfer(c, recipient);
}

public fun receive_coin<T>(state: &mut CCIPReceiverState, _: &OwnerCap, coin_receiving: Receiving<Coin<T>>): Coin<T> {
    transfer::public_receive<Coin<T>>(&mut state.id, coin_receiving)
}

public fun receive_and_send_coin_no_owner_cap<T>(state: &mut CCIPReceiverState, coin_receiving: Receiving<Coin<T>>, recipient: address) {
    let c = transfer::public_receive<Coin<T>>(&mut state.id, coin_receiving);
    transfer::public_transfer(c, recipient);
}

public fun receive_coin_no_owner_cap<T>(state: &mut CCIPReceiverState, coin_receiving: Receiving<Coin<T>>): Coin<T> {
    transfer::public_receive<Coin<T>>(&mut state.id, coin_receiving)
}

// any ccip receiver must implement this function with the same signature
public fun ccip_receive(
    expected_message_id: vector<u8>,
    ref: &CCIPObjectRef,
    message: client::Any2SuiMessage,
    _: &Clock, // this is a precompile, but remain the same across all messages
    state: &mut CCIPReceiverState, // this is a singleton, but remain the same across all messages
) {
    let (
        message_id,
        source_chain_selector,
        sender,
        data,
        dest_token_amounts,
    ) = osh::consume_any2sui_message(ref, message, DummyReceiverProof {});

    assert!(
        message_id == expected_message_id,
        EMessageIdMismatch
    );

    state.counter = state.counter + 1;
    state.message_id = message_id;
    state.source_chain_selector = source_chain_selector;
    state.sender = sender;
    state.data = data;
    state.dest_token_transfer_length = dest_token_amounts.length() as u64;
    state.dest_token_amounts = vector[];

    let mut i = 0;
    while (i < state.dest_token_transfer_length) {
        let (token, amount) = client::get_token_and_amount(&dest_token_amounts[i]);
        state.dest_token_amounts.push_back(TokenAmount { token, amount });
        i = i + 1;
    };

    event::emit(
        ReceivedMessage {
            message_id,
            source_chain_selector,
            sender,
            data,
            dest_token_transfer_length: state.dest_token_transfer_length,
            dest_token_amounts: state.dest_token_amounts,
        }
    );
}