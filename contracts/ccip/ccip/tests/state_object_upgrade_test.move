#[test_only]
module ccip::state_object_upgrade_test;

use ccip::ownable;
use ccip::state_object;
use mcms::mcms_account;
use mcms::mcms_deployer::{Self, DeployerState};
use mcms::mcms_registry::{Self, Registry};
use sui::package;
use sui::test_scenario::{Self as ts, Scenario};

const ADMIN: address = @0xA;
const PACKAGE_OWNER: address = @0xB;

fun init_mcms_and_ccip(ctx: &mut TxContext) {
    mcms_account::test_init(ctx);
    mcms_registry::test_init(ctx);
    mcms_deployer::test_init(ctx);
    state_object::test_init(ctx);
}

// Transfer CCIP ownership to MCMS (which registers @ccip with the registry)
fun transfer_ownership_to_mcms(scenario: &mut Scenario) {
    ts::next_tx(scenario, ADMIN);
    let mut registry = ts::take_shared<Registry>(scenario);
    let mut ref = ts::take_shared<state_object::CCIPObjectRef>(scenario);
    let owner_cap = ts::take_from_sender<ownable::OwnerCap>(scenario);

    // Step 1: transfer_ownership
    state_object::transfer_ownership(
        &mut ref,
        &owner_cap,
        mcms_registry::get_multisig_address(),
        ts::ctx(scenario),
    );

    // Step 2: accept_ownership
    ts::next_tx(scenario, mcms_registry::get_multisig_address());
    state_object::accept_ownership(&mut ref, ts::ctx(scenario));

    ts::next_tx(scenario, ADMIN);

    // Step 3: execute_ownership_transfer_to_mcms (registers with MCMS)
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

#[test]
fun test_upgrade_flow_updates_state_correctly() {
    let mut scenario = ts::begin(ADMIN);

    init_mcms_and_ccip(ts::ctx(&mut scenario));
    transfer_ownership_to_mcms(&mut scenario);

    ts::next_tx(&mut scenario, PACKAGE_OWNER);
    let mut deployer_state = ts::take_shared<DeployerState>(&scenario);
    let registry = ts::take_shared<Registry>(&scenario);
    let ctx = ts::ctx(&mut scenario);

    let upgrade_cap = package::test_publish(@ccip.to_id(), ctx);
    let old_package_address = upgrade_cap.package().to_address();
    let cap_id = object::id(&upgrade_cap);

    mcms_deployer::register_upgrade_cap(
        &mut deployer_state,
        &registry,
        upgrade_cap,
        ctx,
    );

    // Verify initial state
    assert!(mcms_deployer::has_upgrade_cap(&deployer_state, old_package_address), 0);

    ts::return_shared(deployer_state);
    ts::return_shared(registry);

    // Perform upgrade: authorize -> upgrade -> commit
    ts::next_tx(&mut scenario, ADMIN);
    let mut deployer_state = ts::take_shared<DeployerState>(&scenario);
    let owner_cap = ts::take_from_sender<mcms_account::OwnerCap>(&scenario);
    let ctx = ts::ctx(&mut scenario);

    // Step 1: Authorize upgrade
    let ticket = mcms_deployer::authorize_upgrade(
        &owner_cap,
        &mut deployer_state,
        0, // policy
        vector[], // digest
        old_package_address,
        ctx,
    );

    // Step 2: Perform upgrade (simulated with test_upgrade)
    let receipt = package::test_upgrade(ticket);
    let new_package_address = receipt.package().to_address();

    // Verify receipt points to new package but same cap
    assert!(receipt.cap() == cap_id, 1);
    assert!(new_package_address != old_package_address, 2);

    // Step 3: Commit upgrade
    mcms_deployer::commit_upgrade(
        &mut deployer_state,
        receipt,
        ctx,
    );

    // 1. Old package address should NO LONGER be accessible
    assert!(!mcms_deployer::has_upgrade_cap(&deployer_state, old_package_address), 3);

    // 2. New package address SHOULD be accessible
    assert!(mcms_deployer::has_upgrade_cap(&deployer_state, new_package_address), 4);

    ts::return_to_sender(&scenario, owner_cap);
    ts::return_shared(deployer_state);

    // Verify we can authorize another upgrade using NEW package address
    ts::next_tx(&mut scenario, ADMIN);
    let mut deployer_state = ts::take_shared<DeployerState>(&scenario);
    let owner_cap = ts::take_from_sender<mcms_account::OwnerCap>(&scenario);
    let ctx = ts::ctx(&mut scenario);

    // This should work with NEW package address
    let ticket = mcms_deployer::authorize_upgrade(
        &owner_cap,
        &mut deployer_state,
        0,
        vector[],
        new_package_address,
        ctx,
    );

    let receipt = package::test_upgrade(ticket);
    mcms_deployer::commit_upgrade(
        &mut deployer_state,
        receipt,
        ctx,
    );

    ts::return_to_sender(&scenario, owner_cap);
    ts::return_shared(deployer_state);

    ts::end(scenario);
}

#[test]
#[expected_failure(abort_code = mcms::mcms_deployer::EPackageAddressNotRegistered)]
fun test_cannot_authorize_upgrade_with_old_package_address_after_upgrade() {
    let mut scenario = ts::begin(ADMIN);

    init_mcms_and_ccip(ts::ctx(&mut scenario));
    transfer_ownership_to_mcms(&mut scenario);

    let old_package_address: address;

    // Register and perform upgrade
    {
        ts::next_tx(&mut scenario, PACKAGE_OWNER);
        let mut deployer_state = ts::take_shared<DeployerState>(&scenario);
        let registry = ts::take_shared<Registry>(&scenario);
        let ctx = ts::ctx(&mut scenario);

        let upgrade_cap = package::test_publish(@ccip.to_id(), ctx);
        old_package_address = upgrade_cap.package().to_address();

        mcms_deployer::register_upgrade_cap(
            &mut deployer_state,
            &registry,
            upgrade_cap,
            ctx,
        );

        ts::return_shared(deployer_state);
        ts::return_shared(registry);
    };

    // Perform one upgrade
    {
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
        mcms_deployer::commit_upgrade(&mut deployer_state, receipt, ctx);

        ts::return_to_sender(&scenario, owner_cap);
        ts::return_shared(deployer_state);
    };

    // Try to authorize upgrade with OLD package address - should FAIL
    {
        ts::next_tx(&mut scenario, ADMIN);
        let mut deployer_state = ts::take_shared<DeployerState>(&scenario);
        let owner_cap = ts::take_from_sender<mcms_account::OwnerCap>(&scenario);
        let ctx = ts::ctx(&mut scenario);

        // This MUST fail because old_package_address is no longer in upgrade_caps
        let ticket = mcms_deployer::authorize_upgrade(
            &owner_cap,
            &mut deployer_state,
            0,
            vector[],
            old_package_address, // Using OLD address - should fail!
            ctx,
        );

        // This line should never execute due to expected failure above
        // Need to consume hot potato ticket and receipt
        let receipt = package::test_upgrade(ticket);
        mcms_deployer::commit_upgrade(&mut deployer_state, receipt, ctx);

        ts::return_to_sender(&scenario, owner_cap);
        ts::return_shared(deployer_state);
    };

    ts::end(scenario);
}

#[test]
fun test_multiple_upgrades_chain_correctly() {
    let mut scenario = ts::begin(ADMIN);

    init_mcms_and_ccip(ts::ctx(&mut scenario));
    transfer_ownership_to_mcms(&mut scenario);

    let mut current_package_address: address;

    // Register initial UpgradeCap
    {
        ts::next_tx(&mut scenario, PACKAGE_OWNER);
        let mut deployer_state = ts::take_shared<DeployerState>(&scenario);
        let registry = ts::take_shared<Registry>(&scenario);
        let ctx = ts::ctx(&mut scenario);

        let upgrade_cap = package::test_publish(@ccip.to_id(), ctx);
        current_package_address = upgrade_cap.package().to_address();

        mcms_deployer::register_upgrade_cap(
            &mut deployer_state,
            &registry,
            upgrade_cap,
            ctx,
        );

        ts::return_shared(deployer_state);
        ts::return_shared(registry);
    };

    // Perform 3 consecutive upgrades
    let mut i = 0;
    while (i < 3) {
        ts::next_tx(&mut scenario, ADMIN);
        let mut deployer_state = ts::take_shared<DeployerState>(&scenario);
        let owner_cap = ts::take_from_sender<mcms_account::OwnerCap>(&scenario);
        let ctx = ts::ctx(&mut scenario);

        let old_addr = current_package_address;

        let ticket = mcms_deployer::authorize_upgrade(
            &owner_cap,
            &mut deployer_state,
            0,
            vector[],
            current_package_address,
            ctx,
        );

        let receipt = package::test_upgrade(ticket);
        current_package_address = receipt.package().to_address();

        mcms_deployer::commit_upgrade(&mut deployer_state, receipt, ctx);

        // Verify state updated correctly
        assert!(!mcms_deployer::has_upgrade_cap(&deployer_state, old_addr));
        assert!(mcms_deployer::has_upgrade_cap(&deployer_state, current_package_address));

        ts::return_to_sender(&scenario, owner_cap);
        ts::return_shared(deployer_state);

        i = i + 1;
    };

    ts::end(scenario);
}
