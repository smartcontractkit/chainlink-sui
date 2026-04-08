// THIS CONTRACT IS ONLY FOR TESTING PURPOSES. IT IS NOT INTENDED FOR PRODUCTION USE.
module ccip_dummy_receiver::dummy_receiver;

use ccip::client;
use ccip::offramp_state_helper as osh;
use ccip::publisher_wrapper;
use ccip::receiver_registry;
use ccip::state_object::CCIPObjectRef;
use std::string::{Self, String};
use sui::clock::Clock;
use sui::coin::Coin;
use sui::dynamic_field as df;
use sui::event;
use sui::package::{Self, Publisher};
use sui::transfer::Receiving;

const EMessageIdMismatch: u64 = 0;

public struct DUMMY_RECEIVER has drop {}

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
    message_receiver: address,
    token_receiver: address,
    dest_token_transfer_length: u64,
    dest_token_amounts: vector<TokenAmount>,
}

public struct DummyReceiverProof has drop {}

public struct PublisherKey has copy, drop, store {}

public struct TokenAmount has copy, drop, store {
    token: address,
    amount: u256,
}

public fun type_and_version(): String {
    string::utf8(b"DummyReceiver 1.6.0")
}

fun init(otw: DUMMY_RECEIVER, ctx: &mut TxContext) {
    let state = CCIPReceiverState {
        id: object::new(ctx),
        counter: 0,
        message_id: vector[],
        source_chain_selector: 0,
        sender: vector[],
        data: vector[],
        message_receiver: @0x0,
        token_receiver: @0x0,
        dest_token_transfer_length: 0,
        dest_token_amounts: vector[],
    };

    let mut owner_cap = OwnerCap {
        id: object::new(ctx),
        receiver_address: object::id_to_address(object::borrow_id(&state)),
    };

    let publisher = package::claim(otw, ctx);
    df::add(&mut owner_cap.id, PublisherKey {}, publisher);

    transfer::share_object(state);
    transfer::transfer(owner_cap, ctx.sender());
}

public fun register_receiver(owner_cap: &OwnerCap, ref: &mut CCIPObjectRef) {
    let publisher: &Publisher = df::borrow(&owner_cap.id, PublisherKey {});
    let publisher_wrapper = publisher_wrapper::create(publisher, DummyReceiverProof {});
    receiver_registry::register_receiver(ref, publisher_wrapper, DummyReceiverProof {});
}

public fun get_counter(state: &CCIPReceiverState): u64 {
    state.counter
}

public fun get_dest_token_amounts(state: &CCIPReceiverState): vector<TokenAmount> {
    state.dest_token_amounts
}

public fun get_token_receiver(state: &CCIPReceiverState): address {
    state.token_receiver
}

public fun get_token_amount_token(token_amount: &TokenAmount): address {
    token_amount.token
}

public fun get_token_amount_amount(token_amount: &TokenAmount): u256 {
    token_amount.amount
}

// if coin objects are sent to an object (in this case, the receiver state object), this function must be implemented
// in order to "receive" those coin objects. otherwise, the coin objects will be locked in the object until the package
// is upgraded with such a function to receive coin objects from this object.
public fun receive_and_send_coin<T>(
    state: &mut CCIPReceiverState,
    _: &OwnerCap,
    coin_receiving: Receiving<Coin<T>>,
    recipient: address,
) {
    let c = transfer::public_receive<Coin<T>>(&mut state.id, coin_receiving);
    transfer::public_transfer(c, recipient);
}

public fun receive_coin<T>(
    state: &mut CCIPReceiverState,
    _: &OwnerCap,
    coin_receiving: Receiving<Coin<T>>,
): Coin<T> {
    transfer::public_receive<Coin<T>>(&mut state.id, coin_receiving)
}

// DO NOT USE THIS FUNCTION IN PRODUCTION. IT IS ONLY FOR TESTING PURPOSES.
public fun receive_and_send_coin_no_owner_cap<T>(
    state: &mut CCIPReceiverState,
    coin_receiving: Receiving<Coin<T>>,
    recipient: address,
) {
    let c = transfer::public_receive<Coin<T>>(&mut state.id, coin_receiving);
    transfer::public_transfer(c, recipient);
}

// DO NOT USE THIS FUNCTION IN PRODUCTION. IT IS ONLY FOR TESTING PURPOSES.
public fun receive_coin_no_owner_cap<T>(
    state: &mut CCIPReceiverState,
    coin_receiving: Receiving<Coin<T>>,
): Coin<T> {
    transfer::public_receive<Coin<T>>(&mut state.id, coin_receiving)
}

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
        message_receiver,
        token_receiver,
        dest_token_amounts,
    ) = osh::consume_any2sui_message(ref, message, DummyReceiverProof {});

    assert!(message_id == expected_message_id, EMessageIdMismatch);

    state.counter = state.counter + 1;
    state.message_id = message_id;
    state.source_chain_selector = source_chain_selector;
    state.sender = sender;
    state.data = data;
    state.message_receiver = message_receiver;
    state.token_receiver = token_receiver;
    state.dest_token_transfer_length = dest_token_amounts.length() as u64;
    state.dest_token_amounts = vector[];

    let mut i = 0;
    while (i < state.dest_token_transfer_length) {
        let (token, amount) = client::get_token_and_amount(&dest_token_amounts[i]);
        state.dest_token_amounts.push_back(TokenAmount { token, amount });
        i = i + 1;
    };

    event::emit(ReceivedMessage {
        message_id,
        source_chain_selector,
        sender,
        data,
        dest_token_transfer_length: state.dest_token_transfer_length,
        dest_token_amounts: state.dest_token_amounts,
    });
}
