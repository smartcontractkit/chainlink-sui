#[test_only]
module ccip::mcms_ownership_transfer_out_test;

use ccip::ownable::{Self, OwnerCap};
use ccip::state_object::{Self, CCIPObjectRef};
use mcms::mcms_account;
use mcms::mcms_deployer::{Self, DeployerState};
use mcms::mcms_registry::{Self, Registry};
use std::string;
use sui::bcs;
use sui::package::{Self, UpgradeCap};
use sui::test_scenario::{Self as ts, Scenario};

const ADMIN: address = @0xA;
const PACKAGE_OWNER: address = @0xB;
const TRANSFER_TARGET: address = @0xC;

fun init_mcms_and_ccip(ctx: &mut TxContext) {
    mcms_account::test_init(ctx);
    mcms_registry::test_init(ctx);
    mcms_deployer::test_init(ctx);
    state_object::test_init(ctx);
}

fun transfer_ownership_to_mcms(scenario: &mut Scenario) {
    ts::next_tx(scenario, ADMIN);
    let mut registry = ts::take_shared<Registry>(scenario);
    let mut ref = ts::take_shared<CCIPObjectRef>(scenario);
    let owner_cap = ts::take_from_sender<ownable::OwnerCap>(scenario);

    state_object::transfer_ownership(
        &mut ref,
        &owner_cap,
        mcms_registry::get_multisig_address(),
        ts::ctx(scenario),
    );

    ts::next_tx(scenario, mcms_registry::get_multisig_address());
    state_object::accept_ownership(&mut ref, ts::ctx(scenario));

    ts::next_tx(scenario, ADMIN);

    state_object::execute_ownership_transfer_to_mcms(
        &mut ref,
        owner_cap,
        &mut registry,
        mcms_registry::get_multisig_address(),
        ts::ctx(scenario),
    );

    ts::return_shared(registry);
    ts::return_shared(ref);
}

fun register_upgrade_cap(scenario: &mut Scenario): address {
    ts::next_tx(scenario, PACKAGE_OWNER);
    let mut deployer_state = ts::take_shared<DeployerState>(scenario);
    let registry = ts::take_shared<Registry>(scenario);
    let ctx = ts::ctx(scenario);

    let upgrade_cap = package::test_publish(@ccip.to_id(), ctx);
    let package_address = upgrade_cap.package().to_address();

    mcms_deployer::register_upgrade_cap(
        &mut deployer_state,
        &registry,
        upgrade_cap,
        ctx,
    );

    ts::return_shared(deployer_state);
    ts::return_shared(registry);

    package_address
}

fun initiate_and_accept_transfer_from_mcms(
    scenario: &mut Scenario,
    ref: &mut CCIPObjectRef,
    registry: &mut Registry,
) {
    ts::next_tx(scenario, ADMIN);
    let owner_cap_address = mcms_registry::test_get_cap_address<OwnerCap>(
        registry,
        @ccip.to_ascii_string(),
    );

    let mut data = vector::empty<u8>();
    data.append(bcs::to_bytes(&object::id_address(ref)));
    data.append(bcs::to_bytes(&owner_cap_address));
    data.append(bcs::to_bytes(&TRANSFER_TARGET));

    let params = mcms_registry::test_create_executing_callback_params(
        @ccip,
        string::utf8(b"state_object"),
        string::utf8(b"transfer_ownership"),
        data,
        x"0000000000000000000000000000000000000000000000000000000000000001",
        0,
        1,
    );

    state_object::mcms_transfer_ownership(
        ref,
        registry,
        params,
        ts::ctx(scenario),
    );

    ts::next_tx(scenario, TRANSFER_TARGET);
    state_object::accept_ownership(ref, ts::ctx(scenario));
}

fun execute_mcms_ownership_transfer(
    scenario: &mut Scenario,
    ref: &mut CCIPObjectRef,
    registry: &mut Registry,
    deployer_state: &mut DeployerState,
    package_address: address,
    batch_id: vector<u8>,
) {
    let owner_cap_address = mcms_registry::test_get_cap_address<OwnerCap>(
        registry,
        @ccip.to_ascii_string(),
    );

    let mut data = vector::empty<u8>();
    data.append(bcs::to_bytes(&object::id_address(ref)));
    data.append(bcs::to_bytes(&owner_cap_address));
    data.append(bcs::to_bytes(&TRANSFER_TARGET));
    data.append(bcs::to_bytes(&package_address));

    let params = mcms_registry::test_create_executing_callback_params(
        @ccip,
        string::utf8(b"state_object"),
        string::utf8(b"execute_ownership_transfer"),
        data,
        batch_id,
        0,
        1,
    );

    state_object::mcms_execute_ownership_transfer(
        ref,
        registry,
        deployer_state,
        params,
        ts::ctx(scenario),
    );
}

fun assert_recipient_received_caps(
    scenario: &mut Scenario,
    recipient: address,
    expected_package_address: address,
) {
    ts::next_tx(scenario, recipient);
    let owner_cap = ts::take_from_sender<OwnerCap>(scenario);
    let upgrade_cap = ts::take_from_sender<UpgradeCap>(scenario);
    assert!(upgrade_cap.package().to_address() == expected_package_address);
    ts::return_to_address(recipient, owner_cap);
    ts::return_to_address(recipient, upgrade_cap);
}

#[test]
fun test_mcms_execute_ownership_transfer_with_upgrade_cap() {
    let mut scenario = ts::begin(ADMIN);

    init_mcms_and_ccip(ts::ctx(&mut scenario));
    transfer_ownership_to_mcms(&mut scenario);
    let package_address = register_upgrade_cap(&mut scenario);

    ts::next_tx(&mut scenario, ADMIN);
    {
        let registry = ts::take_shared<Registry>(&scenario);
        let deployer_state = ts::take_shared<DeployerState>(&scenario);

        assert!(mcms_registry::is_package_registered(&registry, @ccip.to_ascii_string()));
        assert!(mcms_deployer::has_upgrade_cap(&deployer_state, package_address));

        ts::return_shared(registry);
        ts::return_shared(deployer_state);
    };

    ts::next_tx(&mut scenario, ADMIN);
    let mut ref = ts::take_shared<CCIPObjectRef>(&scenario);
    let mut registry = ts::take_shared<Registry>(&scenario);
    initiate_and_accept_transfer_from_mcms(&mut scenario, &mut ref, &mut registry);
    ts::return_shared(registry);
    ts::return_shared(ref);

    ts::next_tx(&mut scenario, ADMIN);
    let mut ref = ts::take_shared<CCIPObjectRef>(&scenario);
    let mut registry = ts::take_shared<Registry>(&scenario);
    let mut deployer_state = ts::take_shared<DeployerState>(&scenario);

    execute_mcms_ownership_transfer(
        &mut scenario,
        &mut ref,
        &mut registry,
        &mut deployer_state,
        package_address,
        x"0000000000000000000000000000000000000000000000000000000000000002",
    );

    assert!(state_object::owner(&ref) == TRANSFER_TARGET);
    assert!(!mcms_deployer::has_upgrade_cap(&deployer_state, package_address));
    assert!(!mcms_registry::is_package_registered(&registry, @ccip.to_ascii_string()));

    assert_recipient_received_caps(&mut scenario, TRANSFER_TARGET, package_address);

    ts::return_shared(deployer_state);
    ts::return_shared(registry);
    ts::return_shared(ref);
    ts::end(scenario);
}

#[test]
fun test_mcms_execute_ownership_transfer_after_upgrade() {
    let mut scenario = ts::begin(ADMIN);

    init_mcms_and_ccip(ts::ctx(&mut scenario));
    transfer_ownership_to_mcms(&mut scenario);
    let old_package_address = register_upgrade_cap(&mut scenario);

    ts::next_tx(&mut scenario, ADMIN);
    let mut deployer_state = ts::take_shared<DeployerState>(&scenario);
    let owner_cap = ts::take_from_sender<mcms_account::OwnerCap>(&scenario);
    let ctx = ts::ctx(&mut scenario);

    let ticket = mcms_deployer::authorize_upgrade(
        &owner_cap,
        &mut deployer_state,
        0,
        vector[],
        old_package_address,
        ctx,
    );
    let receipt = package::test_upgrade(ticket);
    let new_package_address = receipt.package().to_address();

    mcms_deployer::commit_upgrade(&mut deployer_state, receipt, ctx);

    assert!(!mcms_deployer::has_upgrade_cap(&deployer_state, old_package_address));
    assert!(mcms_deployer::has_upgrade_cap(&deployer_state, new_package_address));

    ts::return_to_sender(&scenario, owner_cap);
    ts::return_shared(deployer_state);

    ts::next_tx(&mut scenario, ADMIN);
    let mut ref = ts::take_shared<CCIPObjectRef>(&scenario);
    let mut registry = ts::take_shared<Registry>(&scenario);
    initiate_and_accept_transfer_from_mcms(&mut scenario, &mut ref, &mut registry);
    ts::return_shared(registry);
    ts::return_shared(ref);

    ts::next_tx(&mut scenario, ADMIN);
    let mut ref = ts::take_shared<CCIPObjectRef>(&scenario);
    let mut registry = ts::take_shared<Registry>(&scenario);
    let mut deployer_state = ts::take_shared<DeployerState>(&scenario);

    execute_mcms_ownership_transfer(
        &mut scenario,
        &mut ref,
        &mut registry,
        &mut deployer_state,
        new_package_address,
        x"0000000000000000000000000000000000000000000000000000000000000003",
    );

    assert!(state_object::owner(&ref) == TRANSFER_TARGET);
    assert!(!mcms_deployer::has_upgrade_cap(&deployer_state, new_package_address));

    assert_recipient_received_caps(&mut scenario, TRANSFER_TARGET, new_package_address);

    ts::return_shared(deployer_state);
    ts::return_shared(registry);
    ts::return_shared(ref);
    ts::end(scenario);
}
