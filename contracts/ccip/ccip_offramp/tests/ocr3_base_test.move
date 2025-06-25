#[test_only]
module ccip_offramp::ocr3_base_test {
    use ccip_offramp::ocr3_base::{Self, OCR3BaseState};
    use sui::test_scenario::{Self as ts, Scenario};

    // === Constants ===
    const ADMIN: address = @0x1;
    const TRANSMITTER_1: address = @0x100;
    const TRANSMITTER_2: address = @0x200;
    const TRANSMITTER_3: address = @0x300;
    const TRANSMITTER_4: address = @0x400;
    const TRANSMITTER_5: address = @0x500;
    
    // Valid ED25519 public keys (32 bytes each)
    const SIGNER_1: vector<u8> = x"1111111111111111111111111111111111111111111111111111111111111111";
    const SIGNER_2: vector<u8> = x"2222222222222222222222222222222222222222222222222222222222222222";
    const SIGNER_3: vector<u8> = x"3333333333333333333333333333333333333333333333333333333333333333";
    const SIGNER_4: vector<u8> = x"4444444444444444444444444444444444444444444444444444444444444444";
    const SIGNER_5: vector<u8> = x"5555555555555555555555555555555555555555555555555555555555555555";

    // Test constants
    const CONFIG_DIGEST_32_BYTES: vector<u8> = x"abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890";
    const BIG_F: u8 = 1;

    // === Helper Functions ===

    fun create_test_scenario(): Scenario {
        ts::begin(ADMIN)
    }

    fun setup_basic_ocr_state(scenario: &mut Scenario): OCR3BaseState {
        let ctx = scenario.ctx();
        ocr3_base::new(ctx)
    }

    fun setup_valid_signers(): vector<vector<u8>> {
        vector[SIGNER_1, SIGNER_2, SIGNER_3, SIGNER_4, SIGNER_5]
    }

    fun setup_valid_transmitters(): vector<address> {
        vector[TRANSMITTER_1, TRANSMITTER_2, TRANSMITTER_3, TRANSMITTER_4, TRANSMITTER_5]
    }

    fun cleanup_scenario(scenario: Scenario, state: OCR3BaseState) {
        // Transfer the state to a dummy address since we can't drop it
        transfer::public_transfer(state, @0x0);
        ts::end(scenario);
    }

    // === Basic Functionality Tests ===

    #[test]
    public fun test_new_ocr3_base_state() {
        let mut scenario = create_test_scenario();
        let state = setup_basic_ocr_state(&mut scenario);
        
        // State should be created successfully
        // We can't inspect internal state directly, but creation success is verified
        
        cleanup_scenario(scenario, state);
    }

    #[test]
    public fun test_ocr_plugin_types() {
        assert!(ocr3_base::ocr_plugin_type_commit() == 0);
        assert!(ocr3_base::ocr_plugin_type_execution() == 1);
    }

    #[test]
    public fun test_set_ocr3_config_commit() {
        let mut scenario = create_test_scenario();
        let mut state = setup_basic_ocr_state(&mut scenario);
        
        let signers = setup_valid_signers();
        let transmitters = setup_valid_transmitters();
        
        ocr3_base::set_ocr3_config(
            &mut state,
            CONFIG_DIGEST_32_BYTES,
            ocr3_base::ocr_plugin_type_commit(),
            BIG_F,
            true, // signature verification enabled for commit
            signers,
            transmitters
        );
        
        // Verify config was set by getting it back
        let config = ocr3_base::latest_config_details(&state, ocr3_base::ocr_plugin_type_commit());
        let (config_digest, big_f, n, is_sig_enabled, config_signers, config_transmitters) = 
            ocr3_base::latest_config_details_fields(config);
            
        assert!(config_digest == CONFIG_DIGEST_32_BYTES);
        assert!(big_f == BIG_F);
        assert!(n == 5); // 5 signers
        assert!(is_sig_enabled == true);
        assert!(config_signers.length() == 5);
        assert!(config_transmitters.length() == 5);
        
        cleanup_scenario(scenario, state);
    }

    #[test]
    public fun test_set_ocr3_config_execution() {
        let mut scenario = create_test_scenario();
        let mut state = setup_basic_ocr_state(&mut scenario);
        
        let transmitters = setup_valid_transmitters();
        
        ocr3_base::set_ocr3_config(
            &mut state,
            CONFIG_DIGEST_32_BYTES,
            ocr3_base::ocr_plugin_type_execution(),
            BIG_F,
            false, // signature verification disabled for execution
            vector[], // no signers needed for execution
            transmitters
        );
        
        // Verify config was set
        let config = ocr3_base::latest_config_details(&state, ocr3_base::ocr_plugin_type_execution());
        let (config_digest, big_f, n, is_sig_enabled, config_signers, config_transmitters) = 
            ocr3_base::latest_config_details_fields(config);
            
        assert!(config_digest == CONFIG_DIGEST_32_BYTES);
        assert!(big_f == BIG_F);
        assert!(n == 0); // no signers for execution
        assert!(is_sig_enabled == false);
        assert!(config_signers.length() == 0);
        assert!(config_transmitters.length() == 5);
        
        cleanup_scenario(scenario, state);
    }

    #[test]
    public fun test_config_signers_and_transmitters() {
        let mut scenario = create_test_scenario();
        let mut state = setup_basic_ocr_state(&mut scenario);
        
        let signers = setup_valid_signers();
        let transmitters = setup_valid_transmitters();
        
        ocr3_base::set_ocr3_config(
            &mut state,
            CONFIG_DIGEST_32_BYTES,
            ocr3_base::ocr_plugin_type_commit(),
            BIG_F,
            true,
            signers,
            transmitters
        );
        
        let config = ocr3_base::latest_config_details(&state, ocr3_base::ocr_plugin_type_commit());
        
        // Test config_signers function
        let retrieved_signers = ocr3_base::config_signers(&config);
        assert!(retrieved_signers.length() == 5);
        assert!(retrieved_signers[0] == SIGNER_1);
        assert!(retrieved_signers[4] == SIGNER_5);
        
        // Test config_transmitters function
        let retrieved_transmitters = ocr3_base::config_transmitters(&config);
        assert!(retrieved_transmitters.length() == 5);
        assert!(retrieved_transmitters[0] == TRANSMITTER_1);
        assert!(retrieved_transmitters[4] == TRANSMITTER_5);
        
        cleanup_scenario(scenario, state);
    }

    #[test]
    public fun test_deserialize_sequence_bytes() {
        // Test with sequence number 9
        let sequence_bytes_9 = x"0000000000000000000000000000000000000000000000000000000000000009";
        let result = ocr3_base::deserialize_sequence_bytes(sequence_bytes_9);
        assert!(result == 9);
        
        // Test with sequence number 255
        let sequence_bytes_255 = x"00000000000000000000000000000000000000000000000000000000000000ff";
        let result = ocr3_base::deserialize_sequence_bytes(sequence_bytes_255);
        assert!(result == 255);
        
        // Test with sequence number 65536 (0x10000)
        let sequence_bytes_65536 = x"0000000000000000000000000000000000000000000000000000000000010000";
        let result = ocr3_base::deserialize_sequence_bytes(sequence_bytes_65536);
        assert!(result == 65536);
    }

    #[test]
    public fun test_transmit_execution_no_signatures() {
        let mut scenario = create_test_scenario();
        let mut state = setup_basic_ocr_state(&mut scenario);
        
        let transmitters = setup_valid_transmitters();
        
        // Set up execution config (no signature verification)
        ocr3_base::set_ocr3_config(
            &mut state,
            CONFIG_DIGEST_32_BYTES,
            ocr3_base::ocr_plugin_type_execution(),
            BIG_F,
            false, // no signature verification
            vector[], // no signers
            transmitters
        );
        
        // Test transmit function
        let report_context = vector[
            CONFIG_DIGEST_32_BYTES,
            x"0000000000000000000000000000000000000000000000000000000000000001" // sequence 1
        ];
        let report = b"test_report_data";
        let signatures = vector[]; // no signatures needed for execution
        
        scenario.next_tx(TRANSMITTER_1);
        let ctx = scenario.ctx();
        
        ocr3_base::transmit(
            &state,
            TRANSMITTER_1,
            ocr3_base::ocr_plugin_type_execution(),
            report_context,
            report,
            signatures,
            ctx
        );
        
        cleanup_scenario(scenario, state);
    }

    // === Error Condition Tests ===

    #[test]
    #[expected_failure(abort_code = ocr3_base::EBigFMustBePositive)]
    public fun test_set_config_zero_big_f() {
        let mut scenario = create_test_scenario();
        let mut state = setup_basic_ocr_state(&mut scenario);
        
        ocr3_base::set_ocr3_config(
            &mut state,
            CONFIG_DIGEST_32_BYTES,
            ocr3_base::ocr_plugin_type_commit(),
            0, // big_f = 0 should fail
            true,
            setup_valid_signers(),
            setup_valid_transmitters()
        );
        
        cleanup_scenario(scenario, state);
    }

    #[test]
    #[expected_failure(abort_code = ocr3_base::EInvalidConfigDigestLength)]
    public fun test_set_config_invalid_digest_length() {
        let mut scenario = create_test_scenario();
        let mut state = setup_basic_ocr_state(&mut scenario);
        
        ocr3_base::set_ocr3_config(
            &mut state,
            b"short_digest", // not 32 bytes
            ocr3_base::ocr_plugin_type_commit(),
            BIG_F,
            true,
            setup_valid_signers(),
            setup_valid_transmitters()
        );
        
        cleanup_scenario(scenario, state);
    }

    #[test]
    #[expected_failure(abort_code = ocr3_base::ENoTransmitters)]
    public fun test_set_config_no_transmitters() {
        let mut scenario = create_test_scenario();
        let mut state = setup_basic_ocr_state(&mut scenario);
        
        ocr3_base::set_ocr3_config(
            &mut state,
            CONFIG_DIGEST_32_BYTES,
            ocr3_base::ocr_plugin_type_commit(),
            BIG_F,
            true,
            setup_valid_signers(),
            vector[] // no transmitters should fail
        );
        
        cleanup_scenario(scenario, state);
    }

    #[test]
    #[expected_failure(abort_code = ocr3_base::EBigFTooHigh)]
    public fun test_set_config_big_f_too_high() {
        let mut scenario = create_test_scenario();
        let mut state = setup_basic_ocr_state(&mut scenario);
        
        // With 5 signers, big_f must be <= 1 (since we need > 3*big_f signers)
        // 5 > 3*1 = 3 ✓, but 5 > 3*2 = 6 ✗
        ocr3_base::set_ocr3_config(
            &mut state,
            CONFIG_DIGEST_32_BYTES,
            ocr3_base::ocr_plugin_type_commit(),
            2, // big_f = 2 is too high for 5 signers
            true,
            setup_valid_signers(),
            setup_valid_transmitters()
        );
        
        cleanup_scenario(scenario, state);
    }

    #[test]
    #[expected_failure(abort_code = ocr3_base::ERepeatedSigners)]
    public fun test_set_config_repeated_signers() {
        let mut scenario = create_test_scenario();
        let mut state = setup_basic_ocr_state(&mut scenario);
        
        let duplicate_signers = vector[
            SIGNER_1, SIGNER_2, SIGNER_1, SIGNER_3, SIGNER_4 // SIGNER_1 repeated
        ];
        
        ocr3_base::set_ocr3_config(
            &mut state,
            CONFIG_DIGEST_32_BYTES,
            ocr3_base::ocr_plugin_type_commit(),
            BIG_F,
            true,
            duplicate_signers,
            setup_valid_transmitters()
        );
        
        cleanup_scenario(scenario, state);
    }

    #[test]
    #[expected_failure(abort_code = ocr3_base::ERepeatedTransmitters)]
    public fun test_set_config_repeated_transmitters() {
        let mut scenario = create_test_scenario();
        let mut state = setup_basic_ocr_state(&mut scenario);
        
        let duplicate_transmitters = vector[
            TRANSMITTER_1, TRANSMITTER_2, TRANSMITTER_1, TRANSMITTER_3, TRANSMITTER_4 // TRANSMITTER_1 repeated
        ];
        
        ocr3_base::set_ocr3_config(
            &mut state,
            CONFIG_DIGEST_32_BYTES,
            ocr3_base::ocr_plugin_type_commit(),
            BIG_F,
            true,
            setup_valid_signers(),
            duplicate_transmitters
        );
        
        cleanup_scenario(scenario, state);
    }

    #[test]
    #[expected_failure(abort_code = ccip::address::EZeroAddressNotAllowed)]
    public fun test_set_config_zero_transmitter() {
        let mut scenario = create_test_scenario();
        let mut state = setup_basic_ocr_state(&mut scenario);
        
        let transmitters_with_zero = vector[
            TRANSMITTER_1, @0x0, TRANSMITTER_3 // zero address should fail
        ];
        
        ocr3_base::set_ocr3_config(
            &mut state,
            CONFIG_DIGEST_32_BYTES,
            ocr3_base::ocr_plugin_type_execution(),
            BIG_F,
            false,
            vector[], // no signers for execution
            transmitters_with_zero
        );
        
        cleanup_scenario(scenario, state);
    }

    #[test]
    #[expected_failure(abort_code = ccip::address::EZeroAddressNotAllowed)]
    public fun test_set_config_zero_signer() {
        let mut scenario = create_test_scenario();
        let mut state = setup_basic_ocr_state(&mut scenario);
        
        let signers_with_zero = vector[
            SIGNER_1, 
            SIGNER_2,
            SIGNER_3,
            SIGNER_4,
            vector[], // empty/zero signer should fail
        ];
        
        ocr3_base::set_ocr3_config(
            &mut state,
            CONFIG_DIGEST_32_BYTES,
            ocr3_base::ocr_plugin_type_commit(),
            BIG_F, // With 5 signers, big_f=1 satisfies 5 > 3*1
            true,
            signers_with_zero,
            setup_valid_transmitters()
        );
        
        cleanup_scenario(scenario, state);
    }

    #[test]
    #[expected_failure(abort_code = ocr3_base::EStaticConfigCannotBeChanged)]
    public fun test_set_config_change_signature_verification() {
        let mut scenario = create_test_scenario();
        let mut state = setup_basic_ocr_state(&mut scenario);
        
        // First set config with signature verification enabled
        ocr3_base::set_ocr3_config(
            &mut state,
            CONFIG_DIGEST_32_BYTES,
            ocr3_base::ocr_plugin_type_commit(),
            BIG_F,
            true, // signature verification enabled
            setup_valid_signers(),
            setup_valid_transmitters()
        );
        
        // Try to change signature verification setting - should fail
        ocr3_base::set_ocr3_config(
            &mut state,
            x"1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef", // different digest
            ocr3_base::ocr_plugin_type_commit(),
            BIG_F,
            false, // try to disable signature verification - should fail
            vector[], // no signers
            setup_valid_transmitters()
        );
        
        cleanup_scenario(scenario, state);
    }

    #[test]
    #[expected_failure(abort_code = ocr3_base::EConfigNotSet)]
    public fun test_latest_config_details_not_set() {
        let mut scenario = create_test_scenario();
        let state = setup_basic_ocr_state(&mut scenario);
        
        // Try to get config that was never set
        let _config = ocr3_base::latest_config_details(&state, ocr3_base::ocr_plugin_type_commit());
        
        cleanup_scenario(scenario, state);
    }

    #[test]
    #[expected_failure(abort_code = ocr3_base::EUnauthorizedTransmitter)]
    public fun test_transmit_unauthorized_transmitter() {
        let mut scenario = create_test_scenario();
        let mut state = setup_basic_ocr_state(&mut scenario);
        
        let transmitters = vector[TRANSMITTER_1, TRANSMITTER_2]; // only these are authorized
        
        ocr3_base::set_ocr3_config(
            &mut state,
            CONFIG_DIGEST_32_BYTES,
            ocr3_base::ocr_plugin_type_execution(),
            BIG_F,
            false,
            vector[],
            transmitters
        );
        
        let report_context = vector[
            CONFIG_DIGEST_32_BYTES,
            x"0000000000000000000000000000000000000000000000000000000000000001"
        ];
        
        scenario.next_tx(TRANSMITTER_3); // unauthorized transmitter
        let ctx = scenario.ctx();
        
        ocr3_base::transmit(
            &state,
            TRANSMITTER_3, // not in authorized list
            ocr3_base::ocr_plugin_type_execution(),
            report_context,
            b"test_report",
            vector[],
            ctx
        );
        
        cleanup_scenario(scenario, state);
    }

    #[test]
    #[expected_failure(abort_code = ocr3_base::EConfigDigestMismatch)]
    public fun test_transmit_wrong_config_digest() {
        let mut scenario = create_test_scenario();
        let mut state = setup_basic_ocr_state(&mut scenario);
        
        ocr3_base::set_ocr3_config(
            &mut state,
            CONFIG_DIGEST_32_BYTES,
            ocr3_base::ocr_plugin_type_execution(),
            BIG_F,
            false,
            vector[],
            setup_valid_transmitters()
        );
        
        let wrong_digest = x"deadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef";
        let report_context = vector[
            wrong_digest, // wrong config digest
            x"0000000000000000000000000000000000000000000000000000000000000001"
        ];
        
        scenario.next_tx(TRANSMITTER_1);
        let ctx = scenario.ctx();
        
        ocr3_base::transmit(
            &state,
            TRANSMITTER_1,
            ocr3_base::ocr_plugin_type_execution(),
            report_context,
            b"test_report",
            vector[],
            ctx
        );
        
        cleanup_scenario(scenario, state);
    }

    #[test]
    #[expected_failure(abort_code = ocr3_base::EInvalidReportContextLength)]
    public fun test_transmit_invalid_report_context_length() {
        let mut scenario = create_test_scenario();
        let mut state = setup_basic_ocr_state(&mut scenario);
        
        ocr3_base::set_ocr3_config(
            &mut state,
            CONFIG_DIGEST_32_BYTES,
            ocr3_base::ocr_plugin_type_execution(),
            BIG_F,
            false,
            vector[],
            setup_valid_transmitters()
        );
        
        let invalid_report_context = vector[
            CONFIG_DIGEST_32_BYTES // only one element, should be two
        ];
        
        scenario.next_tx(TRANSMITTER_1);
        let ctx = scenario.ctx();
        
        ocr3_base::transmit(
            &state,
            TRANSMITTER_1,
            ocr3_base::ocr_plugin_type_execution(),
            invalid_report_context,
            b"test_report",
            vector[],
            ctx
        );
        
        cleanup_scenario(scenario, state);
    }
} 