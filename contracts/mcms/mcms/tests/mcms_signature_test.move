#[test_only]
module mcms::mcms_signature_test;

use mcms::mcms::{Self};
use sui::hash::keccak256;

#[test]
fun test_ecdsa_recover_evm_addr() {
    let go_signed_hash = x"82beeb06ccbbd8c12deea82a4df9b87c6f034c1e8c290391d8c3f9b9e2ee739f";
    let root = x"207e468568a65fc396a84da1291685cc64756b7ade81eababdcdc32ba6ba26da";
    let valid_until = 1747669800045;
    let message_hash = mcms::compute_eth_message_hash(root, valid_until);

    // Sui needs the pre-image, (prior to keccak256(msg))
    // So we verify the keccak256(msg) is the same as the go implementation
    let move_signed_hash = keccak256(&message_hash);
    assert!(move_signed_hash == go_signed_hash);

    // vector[191,131,46,209,180,199,136,235,95,80,234,47,111,183,245,233,219,11,130,237,242,152,45,154,107,188,164,60,228,105,166,62,60,185,39,100,63,128,40,63,206,104,210,94,14,80,192,59,246,55,253,49,118,139,7,98,94,21,39,20,68,11,18,248,28];
    let mut sig = x"bf832ed1b4c788eb5f50ea2f6fb7f5e9db0b82edf2982d9a6bbca43ce469a63e3cb927643f80283fce68d25e0e50c03bf637fd31768b07625e152714440b12f81c";
    let addr_recovered = mcms::test_ecdsa_recover_evm_addr(message_hash, sig);

    let go_recovered_addr = x"35bc1834a0d8f8e116e4f2243c44088654139e9b";
    assert!(addr_recovered == go_recovered_addr);
}
