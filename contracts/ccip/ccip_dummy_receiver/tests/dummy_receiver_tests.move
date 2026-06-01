#[test_only]
module ccip_dummy_receiver::dummy_receiver_tests;

use ccip::client;
use ccip::ownable::OwnerCap;
use ccip::receiver_registry;
use ccip::state_object::{Self, CCIPObjectRef};
use ccip::upgrade_registry;
use ccip_dummy_receiver::dummy_receiver::{Self, CCIPReceiverState, DummyReceiverProof, OwnerCap as DummyOwnerCap};
use std::ascii;
use std::string;
use std::type_name;
use sui::address;
use sui::clock;
use sui::test_scenario::{Self as ts};

const OWNER: address = @0x1000;
const SOURCE_CHAIN: u64 = 2000;

fun receiver_package_id(): address {
    let proof_tn = type_name::with_defining_ids<DummyReceiverProof>();
    let address_str = type_name::address_string(&proof_tn);
    address::from_ascii_bytes(&ascii::into_bytes(address_str))
}

#[test]
public fun test_type_and_version() {
    let version = dummy_receiver::type_and_version();
    assert!(string::as_bytes(&version) == b"DummyReceiver 1.6.0", 0);
}

#[test]
public fun test_ccip_receive_v2_happy_path() {
    let mut sc = ts::begin(OWNER);
    state_object::test_init(sc.ctx());
    sc.next_tx(OWNER);

    let mut ref = ts::take_shared<CCIPObjectRef>(&sc);
    let ccip_owner_cap = ts::take_from_sender<OwnerCap>(&sc);
    upgrade_registry::initialize(&mut ref, &ccip_owner_cap, sc.ctx());
    receiver_registry::initialize(&mut ref, &ccip_owner_cap, sc.ctx());

    sc.next_tx(OWNER);
    dummy_receiver::test_setup(sc.ctx());

    sc.next_tx(OWNER);
    let mut state = ts::take_shared<CCIPReceiverState>(&sc);
    let owner_cap = ts::take_from_sender<DummyOwnerCap>(&sc);
    dummy_receiver::register_receiver(&owner_cap, &mut ref);

    let mut clock = clock::create_for_testing(sc.ctx());
    let message_id = b"message_id_32_bytes_padding_ok!!";
    let state_addr = object::id_address(&state);
    let message = client::new_any2sui_message_v2_for_test(
        message_id,
        SOURCE_CHAIN,
        b"sender",
        b"payload",
        receiver_package_id(),
        @0x0,
        vector[state_addr],
        vector[],
    );

    dummy_receiver::ccip_receive(message_id, &ref, message, &clock, &mut state);

    assert!(dummy_receiver::get_counter(&state) == 1);

    ts::return_shared(state);
    ts::return_to_sender(&sc, owner_cap);
    ts::return_to_sender(&sc, ccip_owner_cap);
    ts::return_shared(ref);
    clock.destroy_for_testing();
    ts::end(sc);
}
