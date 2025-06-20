#[test_only]
module ccip::client_test;

use sui::bcs;

use ccip::client;
use ccip::eth_abi;

use mcms::bcs_stream;

#[test]
fun test_encode_decode_vector_u8() {
    let input = vector[1, 2, 3, 4, 5];
    let encoded = bcs::to_bytes(&input);

    let mut decode_stream = bcs_stream::new(encoded);
    let decoded = bcs_stream::deserialize_vector_u8(&mut decode_stream);
    assert!(input == decoded, 0);
}

#[test]
fun test_generic_extra_args_v2_encoding() {
    // Test basic encoding
    let gas_limit = 500000u256;
    let allow_ooo = true;
    let encoded = client::encode_generic_extra_args_v2(gas_limit, allow_ooo);

    // Verify structure: tag (4 bytes) + u256 (32 bytes) + bool (1 byte) = 37 bytes
    assert!(encoded.length() == 37, 0);

    // Verify tag
    let tag = eth_abi::slice(&encoded, 0, 4);
    assert!(tag == client::generic_extra_args_v2_tag(), 1);

    // Test with different values
    let gas_limit2 = 0u256;
    let allow_ooo2 = false;
    let encoded2 = client::encode_generic_extra_args_v2(gas_limit2, allow_ooo2);
    assert!(encoded2.length() == 37, 2);

    // Verify they're different (except for tag)
    let data1 = eth_abi::slice(&encoded, 4, encoded.length() - 4);
    let data2 = eth_abi::slice(&encoded2, 4, encoded2.length() - 4);
    assert!(data1 != data2, 3);
}

#[test]
fun test_svm_extra_args_v1_encoding() {
    let compute_units = 100000u32;
    let bitmap = 255u64;
    let allow_ooo = true;
    let token_receiver =
        x"1234567890123456789012345678901234567890123456789012345678901234";
    let accounts = vector[
        x"8f2a9c4b7d6e1f3a5c8b9e2d4f7a1c6b8e5d2f9a4c7b1e6d3f8a5c2b9e4d7f1a",
        x"8f2a9c4b7d6e1f3a5c8b9e2d4f7a1c6b8e5d2f9a4c7b1e6d3f8a5c2b9e4d7f1b"
    ];

    let encoded =
        client::encode_svm_extra_args_v1(
            compute_units,
            bitmap,
            allow_ooo,
            token_receiver,
            accounts
        );

    // Verify tag
    let tag = eth_abi::slice(&encoded, 0, 4);
    assert!(tag == client::svm_extra_args_v1_tag());

    // Verify minimum size (tag + u32 + u64 + bool + token_receiver + accounts)
    assert!(encoded.length() >= 4 + 4 + 8 + 1 + 32);
}

#[test]
fun test_bcs_u256_consistency() {
    // Test that large u256 values encode/decode correctly
    let large_values = vector[
        100000u256,
        1000000u256,
        18446744073709551615u256, // Max u64
        115792089237316195423570985008687907853269984665640564039457584007913129639935u256 // Max u256
    ];

    large_values.do_ref!(
        |value| {
            let encoded = client::encode_generic_extra_args_v2(*value, true);
            assert!(encoded.length() == 37);

            // Extract the encoded u256 bytes (skip tag, take 32 bytes)
            let u256_bytes = eth_abi::slice(&encoded, 4, 32);
            assert!(u256_bytes.length() == 32);
        }
    );
}

#[test]
fun test_bcs_boolean_consistency() {
    // Test that boolean values encode consistently
    let encoded_true = client::encode_generic_extra_args_v2(100u256, true);
    let encoded_false = client::encode_generic_extra_args_v2(100u256, false);

    // Should be same length
    assert!(encoded_true.length() == encoded_false.length());

    // Should differ only in the last byte (the boolean)
    let true_bool_byte = encoded_true[encoded_true.length() - 1];
    let false_bool_byte = encoded_false[encoded_false.length() - 1];

    assert!(true_bool_byte != false_bool_byte);
    assert!(true_bool_byte == 1); // BCS encodes true as 0x01
    assert!(false_bool_byte == 0); // BCS encodes false as 0x00
}

#[test]
fun test_empty_accounts_svm_args() {
    let compute_units = 50000u32;
    let bitmap = 0u64;
    let allow_ooo = false;
    let token_receiver =
        x"0000000000000000000000000000000000000000000000000000000000000000";
    let empty_accounts = vector[];

    let encoded =
        client::encode_svm_extra_args_v1(
            compute_units,
            bitmap,
            allow_ooo,
            token_receiver,
            empty_accounts
        );

    // Should encode successfully
    assert!(encoded.length() >= 4 + 4 + 8 + 1 + 32);

    // Test with single account
    let single_account = vector[x"8f2a9c4b7d6e1f3a5c8b9e2d4f7a1c6b8e5d2f9a4c7b1e6d3f8a5c2b9e4d7f1a"];
    let encoded_with_account =
        client::encode_svm_extra_args_v1(
            compute_units,
            bitmap,
            allow_ooo,
            token_receiver,
            single_account
        );

    // Should be larger than empty accounts version
    assert!(encoded_with_account.length() > encoded.length());
}

#[test]
#[expected_failure(abort_code = client::EInvalidSVMTokenReceiverLength)]
fun test_svm_args_rejects_long_token_receiver() {
    // Test that token receivers longer than 32 bytes are rejected
    let long_receiver =
        x"00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"; // 50 bytes
    // EInvalidSVMTokenReceiverLength
    client::encode_svm_extra_args_v1(100, 0, false, long_receiver, vector[]);
}

#[test]
fun test_encode_svm_extra_args_v1_basic() {
    let token_receiver = x"8f2a9c4b7d6e1f3a5c8b9e2d4f7a1c6b8e5d2f9a4c7b1e6d3f8a5c2b9e4d7f1c";
    let accounts = vector[
        x"8f2a9c4b7d6e1f3a5c8b9e2d4f7a1c6b8e5d2f9a4c7b1e6d3f8a5c2b9e4d7f1a",
        x"8f2a9c4b7d6e1f3a5c8b9e2d4f7a1c6b8e5d2f9a4c7b1e6d3f8a5c2b9e4d7f1b"
    ];

    let result =
        client::encode_svm_extra_args_v1(
            1000u32, 0u64, true, token_receiver, accounts
        );

    // Verify the result starts with the correct tag
    let tag_len = client::svm_extra_args_v1_tag().length();
    let mut i = 0;
    while (i < tag_len) {
        assert!(
            *result.borrow(i) == *client::svm_extra_args_v1_tag().borrow(i),
            0
        );
        i = i + 1;
    };

    // Result should be non-empty and contain the tag
    assert!(result.length() > tag_len, 1);
}

#[test]
fun test_encode_svm_extra_args_v1_empty_accounts() {
    let token_receiver = x"8f2a9c4b7d6e1f3a5c8b9e2d4f7a1c6b8e5d2f9a4c7b1e6d3f8a5c2b9e4d7f1d";
    let accounts = vector[];

    let result =
        client::encode_svm_extra_args_v1(
            500u32, 0u64, false, token_receiver, accounts
        );

    // Should not fail and should contain the tag
    let tag_len = client::svm_extra_args_v1_tag().length();
    assert!(result.length() > tag_len, 0);
}

#[test]
fun test_encode_svm_extra_args_v1_32_byte_addresses() {
    let mut token_receiver = vector[];
    let mut account1 = vector[];
    let mut account2 = vector[];

    // Create exactly 32-byte addresses
    let mut i = 0;
    while (i < 32) {
        token_receiver.push_back((i as u8));
        account1.push_back(((i + 100) as u8));
        account2.push_back(((i + 200) as u8));
        i = i + 1;
    };

    let accounts = vector[account1, account2];

    let result =
        client::encode_svm_extra_args_v1(
            2000u32,
            0xFFFFFFFFFFFFFFFFu64,
            true,
            token_receiver,
            accounts
        );

    let tag_len = client::svm_extra_args_v1_tag().length();
    assert!(result.length() > tag_len, 0);
}

#[test]
#[expected_failure(abort_code = client::EInvalidSVMTokenReceiverLength)]
fun test_encode_svm_extra_args_v1_invalid_token_receiver_length() {
    // This test should fail because we're creating a token_receiver that's longer than 32 bytes
    let mut token_receiver = vector[];
    let mut i = 0;
    while (i < 33) { // 33 bytes - too long
        token_receiver.push_back((i as u8));
        i = i + 1;
    };
    let accounts = vector[];
    client::encode_svm_extra_args_v1(1000u32, 0u64, true, token_receiver, accounts);
}

#[test]
#[expected_failure(abort_code = client::EInvalidSVMAccountLength)]
fun test_encode_svm_extra_args_v1_invalid_account_length() {
    // This test should fail because we're creating an account that's longer than 32 bytes
    let token_receiver = x"8f2a9c4b7d6e1f3a5c8b9e2d4f7a1c6b8e5d2f9a4c7b1e6d3f8a5c2b9e4d7f1e";
    let accounts = vector[x"aaaaaaaa8f2a9c4b7d6e1f3a5c8b9e2d4f7a1c6b8e5d2f9a4c7b1e6d3f8a5c2b9e4d7f1f"];

    client::encode_svm_extra_args_v1(1000u32, 0u64, true, token_receiver, accounts);
}