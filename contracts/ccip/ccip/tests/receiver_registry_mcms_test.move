#[test_only]
#[allow(implicit_const_copy)]
module ccip::receiver_registry_mcms_test;

use ccip::ownable::OwnerCap;
use ccip::publisher_wrapper;
use ccip::receiver_registry;
use ccip::state_object::{Self, CCIPObjectRef};
use ccip::upgrade_registry;
use mcms::mcms_account;
use mcms::mcms_deployer;
use mcms::mcms_registry::{Self, Registry};
use std::string;
use std::type_name;
use sui::address;
use sui::bcs;
use sui::package;
use sui::test_scenario::{Self as ts, Scenario};

const OWNER: address = @0x123;
const RECEIVER_REGISTRY_MODULE_NAME: vector<u8> = b"receiver_registry";

public struct RECEIVER_REGISTRY_MCMS_TEST has drop {}
public struct TestReceiverProof has drop, copy {}

public struct Env {
    scenario: Scenario,
    ref: CCIPObjectRef,
    registry: Registry,
}

fun get_package_id_from_proof<ProofType>(): address {
    let proof_tn = type_name::with_defining_ids<ProofType>();
    let address_str = type_name::address_string(&proof_tn);
    address::from_ascii_bytes(&std::ascii::into_bytes(address_str))
}

fun setup(): Env {
    let mut scenario = ts::begin(OWNER);

    {
        let ctx = scenario.ctx();
        mcms_account::test_init(ctx);
        mcms_registry::test_init(ctx);
        mcms_deployer::test_init(ctx);
        state_object::test_init(ctx);
    };

    scenario.next_tx(OWNER);

    let registry = ts::take_shared<Registry>(&scenario);
    let mut ref = ts::take_shared<CCIPObjectRef>(&scenario);
    let owner_cap = ts::take_from_sender<OwnerCap>(&scenario);

    upgrade_registry::initialize(&mut ref, &owner_cap, scenario.ctx());
    receiver_registry::initialize(&mut ref, &owner_cap, scenario.ctx());
    ts::return_to_address(OWNER, owner_cap);

    scenario.next_tx(OWNER);

    Env {
        scenario,
        ref,
        registry,
    }
}

fun tear_down(env: Env) {
    let Env { scenario, ref, registry } = env;
    ts::return_shared(ref);
    ts::return_shared(registry);
    ts::end(scenario);
}

fun register_test_receiver(ref: &mut CCIPObjectRef, ctx: &mut TxContext) {
    let proof = TestReceiverProof {};
    let publisher = package::test_claim(RECEIVER_REGISTRY_MCMS_TEST {}, ctx);
    let publisher_wrapper = publisher_wrapper::create(&publisher, proof);
    receiver_registry::register_receiver(ref, publisher_wrapper, proof);
    package::burn_publisher(publisher);
}

fun transfer_ownership_to_mcms(env: &mut Env, owner_cap: OwnerCap) {
    state_object::transfer_ownership(
        &mut env.ref,
        &owner_cap,
        mcms_registry::get_multisig_address(),
        env.scenario.ctx(),
    );

    env.scenario.next_tx(mcms_registry::get_multisig_address());
    state_object::accept_ownership(&mut env.ref, env.scenario.ctx());

    state_object::execute_ownership_transfer_to_mcms(
        &mut env.ref,
        owner_cap,
        &mut env.registry,
        @mcms,
        env.scenario.ctx(),
    );

    env.scenario.next_tx(OWNER);
}

#[test]
public fun test_mcms_unregister_receiver() {
    let mut env = setup();
    let package_id = get_package_id_from_proof<TestReceiverProof>();

    register_test_receiver(&mut env.ref, env.scenario.ctx());
    assert!(receiver_registry::is_registered_receiver(&env.ref, package_id));

    let owner_cap = ts::take_from_sender<OwnerCap>(&env.scenario);
    transfer_ownership_to_mcms(&mut env, owner_cap);

    let owner_cap_address = mcms_registry::test_get_cap_address<OwnerCap>(
        &env.registry,
        @ccip.to_ascii_string(),
    );

    let mut data = vector[];
    data.append(bcs::to_bytes(&object::id_address(&env.ref)));
    data.append(bcs::to_bytes(&owner_cap_address));
    data.append(bcs::to_bytes(&package_id));

    let params = mcms_registry::test_create_executing_callback_params(
        @ccip,
        string::utf8(RECEIVER_REGISTRY_MODULE_NAME),
        string::utf8(b"unregister_receiver"),
        data,
        x"0000000000000000000000000000000000000000000000000000000000000001",
        0,
        1,
    );

    receiver_registry::mcms_unregister_receiver(
        &mut env.ref,
        &mut env.registry,
        params,
        env.scenario.ctx(),
    );

    assert!(!receiver_registry::is_registered_receiver(&env.ref, package_id));

    env.tear_down();
}

#[test]
#[expected_failure(abort_code = receiver_registry::EInvalidFunction)]
public fun test_mcms_unregister_receiver_wrong_function_name() {
    let mut env = setup();
    let package_id = get_package_id_from_proof<TestReceiverProof>();

    register_test_receiver(&mut env.ref, env.scenario.ctx());

    let owner_cap = ts::take_from_sender<OwnerCap>(&env.scenario);
    transfer_ownership_to_mcms(&mut env, owner_cap);

    let owner_cap_address = mcms_registry::test_get_cap_address<OwnerCap>(
        &env.registry,
        @ccip.to_ascii_string(),
    );

    let mut data = vector[];
    data.append(bcs::to_bytes(&object::id_address(&env.ref)));
    data.append(bcs::to_bytes(&owner_cap_address));
    data.append(bcs::to_bytes(&package_id));

    let params = mcms_registry::test_create_executing_callback_params(
        @ccip,
        string::utf8(RECEIVER_REGISTRY_MODULE_NAME),
        string::utf8(b"register_receiver"),
        data,
        x"0000000000000000000000000000000000000000000000000000000000000001",
        0,
        1,
    );

    receiver_registry::mcms_unregister_receiver(
        &mut env.ref,
        &mut env.registry,
        params,
        env.scenario.ctx(),
    );

    env.tear_down();
}
