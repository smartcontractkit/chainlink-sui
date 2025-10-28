#[test_only]
module mcms::mcms_order_test;

use mcms::mcms_registry::{Self, Registry};
use std::string;
use sui::test_scenario::{Self as ts, Scenario};

const BATCH_ID_1: vector<u8> = x"0000000000000000000000000000000000000000000000000000000000000001";
const BATCH_ID_2: vector<u8> = x"0000000000000000000000000000000000000000000000000000000000000002";

public struct TestMcmsCallback has drop {}

fun create_test_scenario(): Scenario {
    ts::begin(@0xA)
}

#[test]
fun test_sequential_execution_success() {
    let mut scenario = create_test_scenario();

    {
        let ctx = scenario.ctx();
        mcms_registry::test_init(ctx);
    };

    scenario.next_tx(@0xB);
    {
        let mut registry = scenario.take_shared<Registry>();

        // Create 3 callbacks with same batch_id, sequences 0→1→2

        let param0 = mcms_registry::test_create_executing_callback_params(
            @0x123,
            string::utf8(b"test_module"),
            string::utf8(b"test_function"),
            vector::empty(),
            BATCH_ID_1,
            0, // sequence_number
            3, // total_in_batch
        );

        let param1 = mcms_registry::test_create_executing_callback_params(
            @0x123,
            string::utf8(b"test_module"),
            string::utf8(b"test_function"),
            vector::empty(),
            BATCH_ID_1,
            1, // sequence_number
            3, // total_in_batch
        );

        let param2 = mcms_registry::test_create_executing_callback_params(
            @0x123,
            string::utf8(b"test_module"),
            string::utf8(b"test_function"),
            vector::empty(),
            BATCH_ID_1,
            2, // sequence_number
            3, // total_in_batch
        );

        // Execute in correct order
        assert!(mcms_registry::get_next_expected_sequence(&registry, BATCH_ID_1) == 0);
        let (_, _, _, _) = mcms_registry::get_callback_params_from_mcms(&mut registry, param0);

        assert!(mcms_registry::get_next_expected_sequence(&registry, BATCH_ID_1) == 1);
        let (_, _, _, _) = mcms_registry::get_callback_params_from_mcms(&mut registry, param1);

        assert!(mcms_registry::get_next_expected_sequence(&registry, BATCH_ID_1) == 2);
        assert!(!mcms_registry::is_batch_completed(&registry, BATCH_ID_1));

        let (_, _, _, _) = mcms_registry::get_callback_params_from_mcms(&mut registry, param2);

        // Verify batch is completed and cleaned up
        assert!(mcms_registry::is_batch_completed(&registry, BATCH_ID_1));
        assert!(mcms_registry::get_next_expected_sequence(&registry, BATCH_ID_1) == 0); // Batch removed

        ts::return_shared(registry);
    };

    ts::end(scenario);
}

#[test]
#[expected_failure(abort_code = mcms_registry::EOutOfOrderExecution)]
fun test_out_of_order_execution_fails() {
    let mut scenario = create_test_scenario();

    // Initialize registry
    {
        let ctx = scenario.ctx();
        mcms_registry::test_init(ctx);
    };

    scenario.next_tx(@0xB);
    {
        let mut registry = scenario.take_shared<Registry>();

        // Create 3 callbacks with same batch_id
        let param0 = mcms_registry::test_create_executing_callback_params(
            @0x123,
            string::utf8(b"test_module"),
            string::utf8(b"test_function"),
            vector::empty(),
            BATCH_ID_1, // batch_id
            0,
            3,
        );

        let param2 = mcms_registry::test_create_executing_callback_params(
            @0x123,
            string::utf8(b"test_module"),
            string::utf8(b"test_function"),
            vector::empty(),
            BATCH_ID_1,
            2, // Skip sequence 1!
            3,
        );

        // Execute callback 0 successfully
        let (_, _, _, _) = mcms_registry::get_callback_params_from_mcms(&mut registry, param0);

        // Try to execute callback 2 (skip callback 1) - should fail!
        let (_, _, _, _) = mcms_registry::get_callback_params_from_mcms(&mut registry, param2);

        ts::return_shared(registry);
    };

    ts::end(scenario);
}

#[test]
#[expected_failure(abort_code = mcms_registry::EOutOfOrderExecution)]
fun test_duplicate_sequence_execution_fails() {
    let mut scenario = create_test_scenario();

    {
        let ctx = scenario.ctx();
        mcms_registry::test_init(ctx);
    };

    scenario.next_tx(@0xB);
    {
        let mut registry = scenario.take_shared<Registry>();

        // Create 2 callbacks with same batch_id
        let param0_first = mcms_registry::test_create_executing_callback_params(
            @0x123,
            string::utf8(b"test_module"),
            string::utf8(b"test_function"),
            vector::empty(),
            BATCH_ID_1,
            0,
            2,
        );

        let param0_duplicate = mcms_registry::test_create_executing_callback_params(
            @0x123,
            string::utf8(b"test_module"),
            string::utf8(b"test_function"),
            vector::empty(),
            BATCH_ID_1,
            0, // Same sequence number!
            2,
        );

        // Execute callback 0 successfully
        let (_, _, _, _) = mcms_registry::get_callback_params_from_mcms(
            &mut registry,
            param0_first,
        );

        // Try to execute callback 0 again - should fail!
        let (_, _, _, _) = mcms_registry::get_callback_params_from_mcms(
            &mut registry,
            param0_duplicate,
        );

        ts::return_shared(registry);
    };

    ts::end(scenario);
}

#[test]
fun test_batch_completion_tracking() {
    let mut scenario = create_test_scenario();

    {
        let ctx = scenario.ctx();
        mcms_registry::test_init(ctx);
    };

    scenario.next_tx(@0xB);
    {
        let mut registry = scenario.take_shared<Registry>();

        // Create batch with 3 callbacks
        let param0 = mcms_registry::test_create_executing_callback_params(
            @0x123,
            string::utf8(b"test_module"),
            string::utf8(b"test_function"),
            vector::empty(),
            BATCH_ID_1,
            0,
            3,
        );

        let param1 = mcms_registry::test_create_executing_callback_params(
            @0x123,
            string::utf8(b"test_module"),
            string::utf8(b"test_function"),
            vector::empty(),
            BATCH_ID_1,
            1,
            3,
        );

        let param2 = mcms_registry::test_create_executing_callback_params(
            @0x123,
            string::utf8(b"test_module"),
            string::utf8(b"test_function"),
            vector::empty(),
            BATCH_ID_1,
            2,
            3,
        );

        // Execute callbacks 0, 1 (batch not complete yet)
        let (_, _, _, _) = mcms_registry::get_callback_params_from_mcms(&mut registry, param0);
        assert!(!mcms_registry::is_batch_completed(&registry, BATCH_ID_1), 0);

        let (_, _, _, _) = mcms_registry::get_callback_params_from_mcms(&mut registry, param1);
        assert!(!mcms_registry::is_batch_completed(&registry, BATCH_ID_1), 1);

        // Execute callback 2 (final callback)
        let (_, _, _, _) = mcms_registry::get_callback_params_from_mcms(&mut registry, param2);

        // Verify batch is completed
        assert!(mcms_registry::is_batch_completed(&registry, BATCH_ID_1));

        // Verify state cleaned up from batch_execution table
        assert!(mcms_registry::get_next_expected_sequence(&registry, BATCH_ID_1) == 0);

        ts::return_shared(registry);
    };

    ts::end(scenario);
}

#[test]
fun test_multiple_independent_batches() {
    let mut scenario = create_test_scenario();

    {
        let ctx = scenario.ctx();
        mcms_registry::test_init(ctx);
    };

    scenario.next_tx(@0xB);
    {
        let mut registry = scenario.take_shared<Registry>();

        // Create 2 batches with different batch_ids
        // Batch A has 3 callbacks
        let batch_a_param0 = mcms_registry::test_create_executing_callback_params(
            @0x123,
            string::utf8(b"module_a"),
            string::utf8(b"function_a"),
            vector::empty(),
            BATCH_ID_1,
            0,
            3,
        );

        let batch_a_param1 = mcms_registry::test_create_executing_callback_params(
            @0x123,
            string::utf8(b"module_a"),
            string::utf8(b"function_a"),
            vector::empty(),
            BATCH_ID_1,
            1,
            3,
        );

        let batch_a_param2 = mcms_registry::test_create_executing_callback_params(
            @0x123,
            string::utf8(b"module_a"),
            string::utf8(b"function_a"),
            vector::empty(),
            BATCH_ID_1,
            2,
            3,
        );

        // Batch B has 2 callbacks
        let batch_b_param0 = mcms_registry::test_create_executing_callback_params(
            @0x456,
            string::utf8(b"module_b"),
            string::utf8(b"function_b"),
            vector::empty(),
            BATCH_ID_2,
            0,
            2,
        );

        let batch_b_param1 = mcms_registry::test_create_executing_callback_params(
            @0x456,
            string::utf8(b"module_b"),
            string::utf8(b"function_b"),
            vector::empty(),
            BATCH_ID_2,
            1,
            2,
        );

        // Execute interleaved: A0, A1, B0, A2, B1
        let (_, _, _, _) = mcms_registry::get_callback_params_from_mcms(
            &mut registry,
            batch_a_param0,
        );
        assert!(mcms_registry::get_next_expected_sequence(&registry, BATCH_ID_1) == 1);
        assert!(!mcms_registry::is_batch_completed(&registry, BATCH_ID_1));

        let (_, _, _, _) = mcms_registry::get_callback_params_from_mcms(
            &mut registry,
            batch_a_param1,
        );
        assert!(mcms_registry::get_next_expected_sequence(&registry, BATCH_ID_1) == 2);

        let (_, _, _, _) = mcms_registry::get_callback_params_from_mcms(
            &mut registry,
            batch_b_param0,
        );
        assert!(mcms_registry::get_next_expected_sequence(&registry, BATCH_ID_2) == 1);
        assert!(!mcms_registry::is_batch_completed(&registry, BATCH_ID_2));

        let (_, _, _, _) = mcms_registry::get_callback_params_from_mcms(
            &mut registry,
            batch_a_param2,
        );
        assert!(mcms_registry::is_batch_completed(&registry, BATCH_ID_1));

        let (_, _, _, _) = mcms_registry::get_callback_params_from_mcms(
            &mut registry,
            batch_b_param1,
        );
        assert!(mcms_registry::is_batch_completed(&registry, BATCH_ID_2));

        // Both batches completed independently
        assert!(mcms_registry::get_next_expected_sequence(&registry, BATCH_ID_1) == 0);
        assert!(mcms_registry::get_next_expected_sequence(&registry, BATCH_ID_2) == 0);

        ts::return_shared(registry);
    };

    ts::end(scenario);
}

#[test]
fun test_get_next_expected_sequence() {
    let mut scenario = create_test_scenario();

    {
        let ctx = scenario.ctx();
        mcms_registry::test_init(ctx);
    };

    scenario.next_tx(@0xB);
    {
        let mut registry = scenario.take_shared<Registry>();

        // Create batch with 3 callbacks
        let param0 = mcms_registry::test_create_executing_callback_params(
            @0x123,
            string::utf8(b"test_module"),
            string::utf8(b"test_function"),
            vector::empty(),
            BATCH_ID_1,
            0,
            3,
        );

        let param1 = mcms_registry::test_create_executing_callback_params(
            @0x123,
            string::utf8(b"test_module"),
            string::utf8(b"test_function"),
            vector::empty(),
            BATCH_ID_1,
            1,
            3,
        );

        let param2 = mcms_registry::test_create_executing_callback_params(
            @0x123,
            string::utf8(b"test_module"),
            string::utf8(b"test_function"),
            vector::empty(),
            BATCH_ID_1,
            2,
            3,
        );

        // Check get_next_expected_sequence returns 0 initially (batch not started)
        assert!(mcms_registry::get_next_expected_sequence(&registry, BATCH_ID_1) == 0);

        // Execute callback 0, verify it returns 1
        let (_, _, _, _) = mcms_registry::get_callback_params_from_mcms(&mut registry, param0);
        assert!(mcms_registry::get_next_expected_sequence(&registry, BATCH_ID_1) == 1);

        // Execute callback 1, verify it returns 2
        let (_, _, _, _) = mcms_registry::get_callback_params_from_mcms(&mut registry, param1);
        assert!(mcms_registry::get_next_expected_sequence(&registry, BATCH_ID_1) == 2);

        // Execute callback 2, verify it returns 0 (batch removed)
        let (_, _, _, _) = mcms_registry::get_callback_params_from_mcms(&mut registry, param2);
        assert!(mcms_registry::get_next_expected_sequence(&registry, BATCH_ID_1) == 0);

        ts::return_shared(registry);
    };

    ts::end(scenario);
}
