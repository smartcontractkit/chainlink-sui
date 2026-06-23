#[test_only]
module ccip::decimal_parity_test;

use ccip::offramp_state_helper;

#[test]
public fun test_parity_18_to_9() {
    let result = offramp_state_helper::test_calculate_local_amount(
        1_000_000_000_000_000_000, // 1e18
        18,
        9,
    );
    assert!(result == 1_000_000_000); // 1e9
}

#[test]
public fun test_parity_18_to_6() {
    let result = offramp_state_helper::test_calculate_local_amount(
        1_000_000_000_000_000_000, // 1e18
        18,
        6,
    );
    assert!(result == 1_000_000); // 1e6
}

#[test]
public fun test_parity_6_to_9() {
    let result = offramp_state_helper::test_calculate_local_amount(
        1_000_000, // 1e6
        6,
        9,
    );
    assert!(result == 1_000_000_000); // 1e9
}

#[test]
public fun test_parity_same_decimals() {
    let result = offramp_state_helper::test_calculate_local_amount(
        123_456_789,
        8,
        8,
    );
    assert!(result == 123_456_789);
}

#[test]
public fun test_parity_18_to_8() {
    let result = offramp_state_helper::test_calculate_local_amount(
        5_000_000_000_000_000_000, // 5e18
        18,
        8,
    );
    assert!(result == 500_000_000); // 5e8
}

#[test]
public fun test_parity_8_to_18() {
    let result = offramp_state_helper::test_calculate_local_amount(
        500_000_000, // 5e8
        8,
        18,
    );
    assert!(result == 5_000_000_000_000_000_000); // 5e18
}

#[test]
public fun test_parity_truncation_rounds_down() {
    // 1_999_999_999 with 18->9 should truncate to 1 (not round up)
    let result = offramp_state_helper::test_calculate_local_amount(
        1_999_999_999, // less than 1e10
        18,
        9,
    );
    assert!(result == 1); // floor(1_999_999_999 / 1e9) = 1
}

#[test]
public fun test_parity_zero_amount() {
    let result = offramp_state_helper::test_calculate_local_amount(0, 18, 9);
    assert!(result == 0);
}

#[test]
public fun test_parse_remote_decimals_valid() {
    // ABI-encoded u256(18) = 0x12
    let data = x"0000000000000000000000000000000000000000000000000000000000000012";
    let result = offramp_state_helper::test_parse_remote_decimals(data, 9);
    assert!(result == 18);
}

#[test]
public fun test_parse_remote_decimals_empty_falls_back_to_local() {
    let result = offramp_state_helper::test_parse_remote_decimals(vector[], 9);
    assert!(result == 9);
}

#[test]
#[expected_failure]
public fun test_parse_remote_decimals_invalid_length() {
    let data = x"001122"; // 3 bytes, not 0 or 32
    offramp_state_helper::test_parse_remote_decimals(data, 9);
}

#[test]
public fun test_parity_matrix() {
    // Comprehensive matrix: (source_amount, remote_decimals, local_decimals, expected)
    // These values match what burn_mint_token_pool::token_pool::calculate_local_amount produces
    let source_amounts: vector<u256> = vector[
        1_000_000_000_000_000_000, // 1e18
        500_000_000_000_000_000,   // 0.5e18
        1,                          // minimum
        999_999_999_999_999_999,   // just under 1e18
    ];
    let remote_decimals: vector<u8> = vector[18, 18, 18, 18];
    let local_decimals: vector<u8> = vector[9, 9, 9, 9];
    let expected: vector<u64> = vector[
        1_000_000_000,  // 1e18 / 1e9
        500_000_000,    // 0.5e18 / 1e9
        0,              // 1 / 1e9 rounds to 0
        999_999_999,    // floor(999_999_999_999_999_999 / 1e9)
    ];

    let mut i: u64 = 0;
    while (i < source_amounts.length()) {
        let result = offramp_state_helper::test_calculate_local_amount(
            source_amounts[i],
            remote_decimals[i],
            local_decimals[i],
        );
        assert!(result == expected[i]);
        i = i + 1;
    };
}
