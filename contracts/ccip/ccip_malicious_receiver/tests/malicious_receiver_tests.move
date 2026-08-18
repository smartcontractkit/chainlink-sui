#[test_only]
module ccip_malicious_receiver::malicious_receiver_tests;

use ccip_malicious_receiver::malicious_receiver;
use std::string;

// === Basic Tests ===

#[test]
public fun test_type_and_version() {
    let version = malicious_receiver::type_and_version();
    assert!(string::as_bytes(&version) == b"MaliciousReceiver 1.0.0", 0);
}
