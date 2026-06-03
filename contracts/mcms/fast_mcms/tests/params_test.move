#[test_only]
module mcms::params_test;

use mcms::params;
use std::type_name;

public struct McmsAcceptOwnershipProof has drop {}

public struct McmsAcceptOwnershipProofWithGeneric<phantom T> has drop {}

#[test]
fun test_get_struct_name_mcms_accept_ownership_proof() {
    let type_name = type_name::with_defining_ids<McmsAcceptOwnershipProof>();
    let struct_name = params::get_struct_name(&type_name);
    assert!(struct_name == b"McmsAcceptOwnershipProof");
}

#[test]
fun test_get_struct_name_mcms_accept_ownership_proof_generic() {
    let type_name = type_name::with_defining_ids<McmsAcceptOwnershipProofWithGeneric<u64>>();
    let struct_name = params::get_struct_name(&type_name);
    assert!(struct_name == b"McmsAcceptOwnershipProofWithGeneric");
}

#[test]
fun test_get_struct_name_option_u64() {
    let type_name = type_name::with_defining_ids<std::option::Option<u64>>();
    let struct_name = params::get_struct_name(&type_name);
    assert!(struct_name == b"Option");
}
