// THIS CONTRACT IS ONLY FOR TESTING PURPOSES.
// It exposes a ccip_receive with a generic type parameter to produce a normalized ABI
// containing {"Vector": {"TypeParameter": 0}}, which triggers the poison ABI path
// in the relayer's DecodeParameters/decodeParam.
module ccip_broken_receiver::broken_receiver;

use ccip::client;
use ccip::offramp_state_helper as osh;
use ccip::publisher_wrapper;
use ccip::receiver_registry;
use ccip::state_object::CCIPObjectRef;
use std::string::{Self, String};
use sui::dynamic_field as df;
use sui::package::{Self, Publisher};

public struct BROKEN_RECEIVER has drop {}

public struct OwnerCap has key, store {
    id: UID,
}

public struct BrokenReceiverProof has drop {}

public struct PublisherKey has copy, drop, store {}

public fun type_and_version(): String {
    string::utf8(b"BrokenReceiver 1.0.0-test")
}

fun init(otw: BROKEN_RECEIVER, ctx: &mut TxContext) {
    let mut owner_cap = OwnerCap {
        id: object::new(ctx),
    };

    let publisher = package::claim(otw, ctx);
    df::add(&mut owner_cap.id, PublisherKey {}, publisher);

    transfer::transfer(owner_cap, ctx.sender());
}

public fun register_receiver(owner_cap: &OwnerCap, ref: &mut CCIPObjectRef) {
    let publisher: &Publisher = df::borrow(&owner_cap.id, PublisherKey {});
    let publisher_wrapper = publisher_wrapper::create(publisher, BrokenReceiverProof {});
    receiver_registry::register_receiver(ref, publisher_wrapper, BrokenReceiverProof {});
}

/// This function has a generic type parameter T, producing a normalized ABI with
/// {"Vector": {"TypeParameter": 0}} that the relayer cannot decode.
public fun ccip_receive<T: drop>(
    _items: vector<T>,
    _expected_message_id: vector<u8>,
    ref: &CCIPObjectRef,
    message: client::Any2SuiMessage,
) {
    osh::consume_any2sui_message(ref, message, BrokenReceiverProof {});
}
