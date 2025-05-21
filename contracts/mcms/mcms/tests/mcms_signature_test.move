#[test_only]
module mcms::mcms_signature_test;

use mcms::mcms::{Self};

#[test]
fun test_ecdsa_recover_evm_addr() {
    let root = vector[32, 126, 70, 133, 104, 166, 95, 195, 150, 168, 77, 161, 41, 22, 133, 204, 100, 117, 107, 122, 222, 129, 234, 186, 189, 205, 195, 43, 166, 186, 38, 218];
    let valid_until = 1747415002647;
    let message_hash = mcms::test_compute_eth_message_hash(root, valid_until);

    let mut sig = vector[106, 24, 116, 54, 214, 108, 191, 140, 11, 163, 94, 71, 128, 23, 54, 81, 94, 0, 164, 32, 93, 58, 132, 100, 159, 78, 102, 244, 248, 18, 155, 94, 82, 190, 98, 222, 171, 108, 75, 54, 119, 212, 73, 145, 251, 25, 187, 189, 55, 43, 107, 198, 98, 189, 1, 254, 138, 36, 15, 230, 176, 138, 143, 72, 27];
    let addr_recovered = mcms::test_ecdsa_recover_evm_addr(message_hash, sig);
    
    std::debug::print(&std::string::utf8(b"Recovered address:"));
    std::debug::print(&addr_recovered);

    // 6a07dba2435035822e36091f1db4855d0b28de94 - Go EVM recovered address
    let go_recovered_addr = vector[106, 7, 219, 162, 67, 80, 53, 130, 46, 54, 9, 31, 29, 180, 133, 93, 11, 40, 222, 148];

    // Below fails as Sui ecrecover is not the same as Go ecrecover
    assert!(addr_recovered == go_recovered_addr);
}

fun to_hex_string(bytes: &vector<u8>): std::string::String {
    let hex_chars = vector[
        // ASCII values for '0' through '9'
        48, 49, 50, 51, 52, 53, 54, 55, 56, 57,
        // ASCII values for 'a' through 'f'
        97, 98, 99, 100, 101, 102
    ];
    
    let mut result = vector::empty<u8>();
    let mut i = 0;
    while (i < bytes.length()) {
        let byte = bytes[i];
        vector::push_back(&mut result, hex_chars[(byte >> 4) as u64]);
        vector::push_back(&mut result, hex_chars[(byte & 0xf) as u64]);
        i = i + 1;
    };
    
    std::string::utf8(result)
}
