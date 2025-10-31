#[test_only]
module managed_token_faucet::faucet_tests;

use managed_token_faucet::faucet;
use std::string;

#[test]
public fun test_type_and_version() {
    let version = faucet::type_and_version();
    assert!(string::as_bytes(&version) == b"ManagedTokenFaucet 1.6.0", 0);
}
