#[test_only]
module cw_tests::cw_tests_tests;

use cw_tests::cw_tests;

//const ENotImplemented: u64 = 0;


#[test]
fun test_hello_world() {
    assert!(cw_tests::hello_world() == b"Hello, World!".to_string(), 0);
}
