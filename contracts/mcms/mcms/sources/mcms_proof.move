/// Separate module to create MCMS proof, this is needed to avoid circular dependencies.
module mcms::mcms_proof;

use mcms::params;
use std::string::String;
use std::type_name;

public struct McmsProof has drop {}

const ENotMcmsTimelock: u64 = 1;

public(package) fun create_mcms_proof(): McmsProof {
    McmsProof {}
}

public(package) fun assert_is_mcms_timelock<T: drop>(_witness: T): (address, String) {
    let proof_type = type_name::get<T>();
    let (proof_account_address, proof_module_name) = params::get_account_address_and_module_name(
        proof_type,
    );
    assert!(
        proof_account_address == @mcms && *proof_module_name.as_bytes() == b"mcms_proof",
        ENotMcmsTimelock,
    );
    (proof_account_address, proof_module_name)
}
