#[test_only]
module ccip_dummy_receiver::dummy_receiver_tests;

use ccip_dummy_receiver::dummy_receiver;
use std::string;

// === Basic Tests ===

#[test]
public fun test_type_and_version() {
    let version = dummy_receiver::type_and_version();
    assert!(string::as_bytes(&version) == b"DummyReceiver 1.6.0", 0);
}
