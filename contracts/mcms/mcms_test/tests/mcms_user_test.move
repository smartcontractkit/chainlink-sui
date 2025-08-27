#[test_only]
module mcms_test::mcms_user_test;

use mcms::mcms_deployer::{Self, DeployerState};
use mcms::mcms_registry::{Self, Registry};
use mcms_test::mcms_user::{Self, UserData, OwnerCap};
use std::string;
use sui::bcs;
use sui::package;
use sui::test_scenario::{Self as ts, Scenario};
use sui::test_utils;

const SENDER: address = @0xA;
const MODULE_NAME: vector<u8> = b"mcms_user";
const FUNCTION_ONE: vector<u8> = b"function_one";
const FUNCTION_TWO: vector<u8> = b"function_two";

const TEST_ARG_STRING: vector<u8> = b"test_string";
const TEST_ARG_BYTES: vector<u8> = b"test_bytes";
const TEST_ARG_ADDRESS: address = @0xCAFE;
const TEST_ARG_U128: u128 = 12345;

fun create_test_scenario(): Scenario {
    ts::begin(@0xA)
}

/// Set up test environment with mcms_registry and user_data
/// Sends owner_cap to the registry
fun setup_mcms_registry_and_user_data(scenario: &mut Scenario): address {
    // Transaction 1: Initialize registry
    {
        ts::next_tx(scenario, SENDER);
        let ctx = ts::ctx(scenario);
        mcms_registry::test_init(ctx);
    };

    // Transaction 2: Initialize user data
    {
        ts::next_tx(scenario, SENDER);
        let ctx = ts::ctx(scenario);
        mcms_user::test_init(ctx);
    };

    // Transaction 3: Initialize deployer state
    {
        ts::next_tx(scenario, SENDER);
        let ctx = ts::ctx(scenario);
        mcms_deployer::test_init(ctx);
    };

    // Transaction 4: Register entrypoint
    {
        ts::next_tx(scenario, SENDER);

        let mut registry = ts::take_shared<Registry>(scenario);
        let user_data = ts::take_shared<UserData>(scenario);
        let owner_cap = ts::take_from_sender<OwnerCap>(scenario);
        let ctx = ts::ctx(scenario);


        // Initialize the user data with mcms_registry
        // This creates a owner_cap and mcms_registry owns this cap
        mcms_user::register_mcms_entrypoint(
            owner_cap,
            &mut registry,
            &user_data,
            ctx,
        );

        ts::return_shared(user_data);
        ts::return_shared(registry);
    };

    {
        ts::next_tx(scenario, SENDER);

        let mut deployer_state = ts::take_shared<DeployerState>(scenario);
        let mut registry = ts::take_shared<Registry>(scenario);
        let user_data = ts::take_shared<UserData>(scenario);
        let ctx = ts::ctx(scenario);

        let upgrade_cap = package::test_publish(@mcms_test.to_id(), ctx);

        mcms_user::register_upgrade_cap(
            &mut deployer_state,
            upgrade_cap,
            &mut registry,
            ctx,
        );

        ts::return_shared(user_data);
        ts::return_shared(registry);
        ts::return_shared(deployer_state);
    };

    SENDER
}

#[test]
fun test_register_and_initialize() {
    let mut scenario = create_test_scenario();

    let sender = setup_mcms_registry_and_user_data(&mut scenario);

    // Verify user data is initialized correctly
    {
        ts::next_tx(&mut scenario, sender);

        let user_data = ts::take_shared<UserData>(&scenario);
        assert!(mcms_user::get_invocations(&user_data) == 0);

        ts::return_shared(user_data);
    };

    ts::end(scenario);
}

#[test]
fun test_mcms_entrypoint_function_one() {
    let mut scenario = create_test_scenario();

    let sender = setup_mcms_registry_and_user_data(&mut scenario);

    // Transaction 4: Execute mcms_entrypoint with function_one
    {
        ts::next_tx(&mut scenario, sender);

        let mut registry = ts::take_shared<Registry>(&scenario);
        let mut user_data = ts::take_shared<UserData>(&scenario);

        let arg1 = string::utf8(TEST_ARG_STRING);
        let arg2 = TEST_ARG_BYTES;

        // Serialize arguments for BCS
        let mut data = vector::empty<u8>();
        vector::append(&mut data, bcs::to_bytes(&arg1));
        vector::append(&mut data, bcs::to_bytes(&arg2));

        let ctx = ts::ctx(&mut scenario);

        // PTB Construction
        // Command 1: Create callback params (same return result as mcms::execute)
        let params = mcms_registry::test_create_executing_callback_params(
            @mcms_test,
            string::utf8(MODULE_NAME),
            string::utf8(FUNCTION_ONE),
            data,
        );

        // Command 2: Call mcms_entrypoint with params hot potato
        mcms_user::mcms_entrypoint(
            &mut user_data,
            &mut registry,
            params,
            ctx,
        );

        // Verify user_data has been updated using accessor functions
        assert!(mcms_user::get_invocations(&user_data) == 1);
        assert!(*mcms_user::get_field_a(&user_data).as_bytes() == TEST_ARG_STRING);
        assert!(mcms_user::get_field_b(&user_data) == TEST_ARG_BYTES);

        ts::return_shared(user_data);
        ts::return_shared(registry);
    };

    ts::end(scenario);
}

#[test]
fun test_mcms_entrypoint_function_two() {
    let mut scenario = create_test_scenario();

    let sender = setup_mcms_registry_and_user_data(&mut scenario);

    // Transaction 4: Execute mcms_entrypoint with function_two
    {
        ts::next_tx(&mut scenario, sender);

        let mut registry = ts::take_shared<Registry>(&scenario);
        let mut user_data = ts::take_shared<UserData>(&scenario);

        let arg1 = TEST_ARG_ADDRESS;
        let arg2 = TEST_ARG_U128;

        let mut data = vector::empty<u8>();
        vector::append(&mut data, bcs::to_bytes(&arg1));
        vector::append(&mut data, bcs::to_bytes(&arg2));

        let ctx = ts::ctx(&mut scenario);

        // PTB Construction
        // Command 1: Create callback params (same return result as mcms::execute)
        let params = mcms_registry::test_create_executing_callback_params(
            @mcms_test,
            string::utf8(MODULE_NAME),
            string::utf8(FUNCTION_TWO),
            data,
        );

        // Command 2: Call mcms_entrypoint
        mcms_user::mcms_entrypoint(
            &mut user_data,
            &mut registry,
            params,
            ctx,
        );

        // Verify user_data has been updated
        assert!(mcms_user::get_invocations(&user_data) == 1);
        assert!(mcms_user::get_field_c(&user_data) == TEST_ARG_ADDRESS);
        assert!(mcms_user::get_field_d(&user_data) == TEST_ARG_U128);

        ts::return_shared(user_data);
        ts::return_shared(registry);
    };

    ts::end(scenario);
}

#[test]
#[expected_failure(abort_code = mcms_user::EUnknownFunction)]
fun test_mcms_entrypoint_unknown_function() {
    let mut scenario = create_test_scenario();

    let sender = setup_mcms_registry_and_user_data(&mut scenario);

    // Transaction 4: Execute mcms_entrypoint with unknown function
    {
        ts::next_tx(&mut scenario, sender);

        let mut registry = ts::take_shared<Registry>(&scenario);
        let mut user_data = ts::take_shared<UserData>(&scenario);

        // PTB Construction
        // Command 1: Create callback params with unknown function name
        let params = mcms_registry::test_create_executing_callback_params(
            @mcms_test,
            string::utf8(MODULE_NAME),
            string::utf8(b"unknown_function"),
            vector::empty(),
        );

        let ctx = ts::ctx(&mut scenario);
        // Command 2: This should fail because the function name is unknown
        mcms_user::mcms_entrypoint(
            &mut user_data,
            &mut registry,
            params,
            ctx,
        );

        ts::return_shared(user_data);
        ts::return_shared(registry);
    };

    ts::end(scenario);
}

#[test]
#[expected_failure(abort_code = mcms_registry::EModuleNameMismatch)]
fun test_mcms_entrypoint_wrong_module_name() {
    let mut scenario = create_test_scenario();

    let sender = setup_mcms_registry_and_user_data(&mut scenario);

    // Transaction 4: Execute mcms_entrypoint with wrong package name
    {
        ts::next_tx(&mut scenario, sender);

        let mut registry = ts::take_shared<Registry>(&scenario);
        let mut user_data = ts::take_shared<UserData>(&scenario);

        // PTB Construction
        // Command 1: Create callback params with wrong module name
        let params = mcms_registry::test_create_executing_callback_params(
            @mcms_test,
            string::utf8(b"wrong_module_name"),
            string::utf8(FUNCTION_ONE),
            vector::empty(),
        );

        let ctx = ts::ctx(&mut scenario);
        // Command 2: This should fail because the module name doesn't match
        mcms_user::mcms_entrypoint(
            &mut user_data,
            &mut registry,
            params,
            ctx,
        );

        ts::return_shared(user_data);
        ts::return_shared(registry);
    };

    ts::end(scenario);
}

#[test]
fun test_sequential_function_calls() {
    let mut scenario = create_test_scenario();

    let sender = setup_mcms_registry_and_user_data(&mut scenario);

    // Transaction 4: Call function_one
    {
        ts::next_tx(&mut scenario, sender);

        let mut registry = ts::take_shared<Registry>(&scenario);
        let mut user_data = ts::take_shared<UserData>(&scenario);

        let arg1 = string::utf8(TEST_ARG_STRING);
        let arg2 = TEST_ARG_BYTES;

        // Serialize arguments for BCS
        let mut data = vector::empty<u8>();
        vector::append(&mut data, bcs::to_bytes(&arg1));
        vector::append(&mut data, bcs::to_bytes(&arg2));

        let ctx = ts::ctx(&mut scenario);

        // PTB Construction
        // Command 1: Create callback params
        let params = mcms_registry::test_create_executing_callback_params(
            @mcms_test,
            string::utf8(MODULE_NAME),
            string::utf8(FUNCTION_ONE),
            data,
        );

        // Command 2: Call mcms_entrypoint
        mcms_user::mcms_entrypoint(
            &mut user_data,
            &mut registry,
            params,
            ctx,
        );

        // Verify invocation count
        assert!(mcms_user::get_invocations(&user_data) == 1);

        ts::return_shared(user_data);
        ts::return_shared(registry);
    };

    // Transaction 5: Call function_two
    {
        ts::next_tx(&mut scenario, sender);

        let mut registry = ts::take_shared<Registry>(&scenario);
        let mut user_data = ts::take_shared<UserData>(&scenario);

        let arg1 = TEST_ARG_ADDRESS;
        let arg2 = TEST_ARG_U128;

        // Serialize arguments for BCS
        let mut data = vector::empty<u8>();
        vector::append(&mut data, bcs::to_bytes(&arg1));
        vector::append(&mut data, bcs::to_bytes(&arg2));

        // PTB Construction
        // Command 1: Create callback params
        let params = mcms_registry::test_create_executing_callback_params(
            @mcms_test,
            string::utf8(MODULE_NAME),
            string::utf8(FUNCTION_TWO),
            data,
        );

        let ctx = ts::ctx(&mut scenario);

        // Command 2: Call mcms_entrypoint
        mcms_user::mcms_entrypoint(
            &mut user_data,
            &mut registry,
            params,
            ctx,
        );

        // Verify invocation count increased
        assert!(mcms_user::get_invocations(&user_data) == 2);

        // Verify the data from both calls is present
        assert!(*mcms_user::get_field_a(&user_data).as_bytes() == TEST_ARG_STRING);
        assert!(mcms_user::get_field_b(&user_data) == TEST_ARG_BYTES);
        assert!(mcms_user::get_field_c(&user_data) == TEST_ARG_ADDRESS);
        assert!(mcms_user::get_field_d(&user_data) == TEST_ARG_U128);

        ts::return_shared(user_data);
        ts::return_shared(registry);
    };

    ts::end(scenario);
}

#[test]
#[expected_failure(abort_code = mcms_user::EInvalidAdminCap)]
fun test_call_function_with_invalid_user_data() {
    let mut scenario = create_test_scenario();

    let sender = setup_mcms_registry_and_user_data(&mut scenario);

    // Test fake user_data creation
    ts::next_tx(&mut scenario, sender);
    {
        let mut registry = ts::take_shared<Registry>(&scenario);

        let ctx = ts::ctx(&mut scenario);

        let arg1 = string::utf8(TEST_ARG_STRING);
        let arg2 = TEST_ARG_BYTES;

        let mut data = vector::empty<u8>();
        vector::append(&mut data, bcs::to_bytes(&arg1));
        vector::append(&mut data, bcs::to_bytes(&arg2));

        let params = mcms_registry::test_create_executing_callback_params(
            @mcms_test,
            string::utf8(MODULE_NAME),
            string::utf8(FUNCTION_ONE),
            data,
        );

        let fake_owner_cap = object::new(ctx);
        // Create a fake user_data
        let mut fake_user_data = mcms_user::test_create_user_data(ctx, fake_owner_cap.to_inner());

        // Command 2: This should fail because we provide an unregistered user_data
        // The cap does not exist for this user_data
        mcms_user::mcms_entrypoint(
            &mut fake_user_data,
            &mut registry,
            params,
            ctx,
        );

        ts::return_shared(fake_user_data);
        ts::return_shared(registry);
        test_utils::destroy(fake_owner_cap);
    };

    ts::end(scenario);
}
