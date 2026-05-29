#[test_only]
module ccip::offramp_state_helper_v2_tests;

use ccip::client;
use ccip::offramp_state_helper::{Self as osh, DestTransferCap};
use ccip::ownable::OwnerCap;
use ccip::receiver_registry;
use ccip::state_object::{Self, CCIPObjectRef};
use ccip::upgrade_registry;
use sui::test_scenario::{Self as ts};

const OWNER: address = @0x1000;
const SOURCE_CHAIN: u64 = 2000;
const RECEIVER_PKG: address = @0xABCD;

#[test]
public fun test_v2_extract_with_correct_object_ids() {
    let mut sc = ts::begin(OWNER);
    state_object::test_init(sc.ctx());
    osh::test_init(sc.ctx());
    sc.next_tx(OWNER);

    let owner_cap = ts::take_from_sender<OwnerCap>(&sc);
    let mut ref_obj = ts::take_shared<CCIPObjectRef>(&sc);
    let dest_cap = ts::take_from_sender<DestTransferCap>(&sc);

    upgrade_registry::initialize(&mut ref_obj, &owner_cap, sc.ctx());
    receiver_registry::initialize(&mut ref_obj, &owner_cap, sc.ctx());

    let object_id_a = @0x1111;
    let object_id_b = @0x2222;

    let mut receiver_params = osh::create_receiver_params_v2(&dest_cap, SOURCE_CHAIN);

    let msg = osh::new_any2sui_message_v2(
        &dest_cap,
        b"message_id_32_bytes_padding_ok!!",
        SOURCE_CHAIN,
        b"sender",
        b"payload",
        RECEIVER_PKG,
        @0x0,
        vector[object_id_a, object_id_b],
        vector[],
        vector[],
    );

    osh::populate_message_v2(&dest_cap, &mut receiver_params, msg);

    // Extract with CORRECT object IDs — should succeed
    let extracted = osh::extract_any2sui_message_v2(
        &mut receiver_params,
        vector[object_id_a, object_id_b],
    );

    // Consume via test helper
    osh::deconstruct_receiver_params_v2_with_message_for_test(&dest_cap, RECEIVER_PKG, receiver_params, extracted);

    ts::return_to_sender(&sc, owner_cap);
    ts::return_to_sender(&sc, dest_cap);
    ts::return_shared(ref_obj);
    ts::end(sc);
}

#[test]
#[expected_failure(abort_code = osh::EReceiverObjectMismatch)]
public fun test_v2_extract_with_wrong_object_ids_aborts() {
    let mut sc = ts::begin(OWNER);
    state_object::test_init(sc.ctx());
    osh::test_init(sc.ctx());
    sc.next_tx(OWNER);

    let owner_cap = ts::take_from_sender<OwnerCap>(&sc);
    let mut ref_obj = ts::take_shared<CCIPObjectRef>(&sc);
    let dest_cap = ts::take_from_sender<DestTransferCap>(&sc);

    upgrade_registry::initialize(&mut ref_obj, &owner_cap, sc.ctx());
    receiver_registry::initialize(&mut ref_obj, &owner_cap, sc.ctx());

    let mut receiver_params = osh::create_receiver_params_v2(&dest_cap, SOURCE_CHAIN);

    let msg = osh::new_any2sui_message_v2(
        &dest_cap,
        b"message_id_32_bytes_padding_ok!!",
        SOURCE_CHAIN,
        b"sender",
        b"payload",
        RECEIVER_PKG,
        @0x0,
        vector[@0x1111], // committed: object A
        vector[],
        vector[],
    );

    osh::populate_message_v2(&dest_cap, &mut receiver_params, msg);

    // Extract with WRONG object IDs — should abort with EReceiverObjectMismatch
    let _extracted = osh::extract_any2sui_message_v2(
        &mut receiver_params,
        vector[@0x9999], // attacker substitutes different object
    );

    // unreachable — cleanup for compiler
    abort 0
}

#[test]
public fun test_v2_extract_with_empty_object_ids() {
    let mut sc = ts::begin(OWNER);
    state_object::test_init(sc.ctx());
    osh::test_init(sc.ctx());
    sc.next_tx(OWNER);

    let owner_cap = ts::take_from_sender<OwnerCap>(&sc);
    let mut ref_obj = ts::take_shared<CCIPObjectRef>(&sc);
    let dest_cap = ts::take_from_sender<DestTransferCap>(&sc);

    upgrade_registry::initialize(&mut ref_obj, &owner_cap, sc.ctx());
    receiver_registry::initialize(&mut ref_obj, &owner_cap, sc.ctx());

    let mut receiver_params = osh::create_receiver_params_v2(&dest_cap, SOURCE_CHAIN);

    let msg = osh::new_any2sui_message_v2(
        &dest_cap,
        b"message_id_32_bytes_padding_ok!!",
        SOURCE_CHAIN,
        b"sender",
        b"payload",
        RECEIVER_PKG,
        @0x0,
        vector[], // stateless receiver: no objects committed
        vector[],
        vector[],
    );

    osh::populate_message_v2(&dest_cap, &mut receiver_params, msg);

    // Extract with empty used_object_ids — should succeed
    let extracted = osh::extract_any2sui_message_v2(
        &mut receiver_params,
        vector[],
    );

    osh::deconstruct_receiver_params_v2_with_message_for_test(&dest_cap, RECEIVER_PKG, receiver_params, extracted);

    ts::return_to_sender(&sc, owner_cap);
    ts::return_to_sender(&sc, dest_cap);
    ts::return_shared(ref_obj);
    ts::end(sc);
}

#[test]
public fun test_v2_no_message_populated_deconstruct_succeeds() {
    let mut sc = ts::begin(OWNER);
    state_object::test_init(sc.ctx());
    osh::test_init(sc.ctx());
    sc.next_tx(OWNER);

    let owner_cap = ts::take_from_sender<OwnerCap>(&sc);
    let mut ref_obj = ts::take_shared<CCIPObjectRef>(&sc);
    let dest_cap = ts::take_from_sender<DestTransferCap>(&sc);

    upgrade_registry::initialize(&mut ref_obj, &owner_cap, sc.ctx());

    // Token-only or unregistered receiver: no message populated
    let receiver_params = osh::create_receiver_params_v2(&dest_cap, SOURCE_CHAIN);
    osh::deconstruct_receiver_params_v2(&dest_cap, receiver_params);

    ts::return_to_sender(&sc, owner_cap);
    ts::return_to_sender(&sc, dest_cap);
    ts::return_shared(ref_obj);
    ts::end(sc);
}

#[test]
#[expected_failure(abort_code = osh::ECCIPReceiveFailed)]
public fun test_v2_message_not_consumed_deconstruct_aborts() {
    let mut sc = ts::begin(OWNER);
    state_object::test_init(sc.ctx());
    osh::test_init(sc.ctx());
    sc.next_tx(OWNER);

    let owner_cap = ts::take_from_sender<OwnerCap>(&sc);
    let mut ref_obj = ts::take_shared<CCIPObjectRef>(&sc);
    let dest_cap = ts::take_from_sender<DestTransferCap>(&sc);

    upgrade_registry::initialize(&mut ref_obj, &owner_cap, sc.ctx());
    receiver_registry::initialize(&mut ref_obj, &owner_cap, sc.ctx());

    let mut receiver_params = osh::create_receiver_params_v2(&dest_cap, SOURCE_CHAIN);

    let msg = osh::new_any2sui_message_v2(
        &dest_cap,
        b"message_id_32_bytes_padding_ok!!",
        SOURCE_CHAIN,
        b"sender",
        b"payload",
        RECEIVER_PKG,
        @0x0,
        vector[@0x1111],
        vector[],
        vector[],
    );

    osh::populate_message_v2(&dest_cap, &mut receiver_params, msg);

    // Do NOT extract the message — simulates skipping ccip_receive
    // deconstruct should abort because message is still populated
    osh::deconstruct_receiver_params_v2(&dest_cap, receiver_params);

    abort 0 // unreachable
}

#[test]
public fun test_v2_receiver_object_ids_stored_correctly() {
    let msg = client::new_any2sui_message_v2(
        b"message_id_32_bytes_padding_ok!!",
        SOURCE_CHAIN,
        b"sender",
        b"data",
        RECEIVER_PKG,
        @0x0,
        vector[@0x1111, @0x2222, @0x3333],
        vector[],
    );

    assert!(*client::get_receiver_object_ids(&msg) == vector[@0x1111, @0x2222, @0x3333]);

    // Consume to satisfy no-drop
    client::consume_any2sui_message_v2(msg, RECEIVER_PKG);
}
