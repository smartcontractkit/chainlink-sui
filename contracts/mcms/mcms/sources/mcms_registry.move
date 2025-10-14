module mcms::mcms_registry;

use mcms::params;
use std::string::String;
use std::type_name::{Self, TypeName};
use sui::address;
use sui::bag::{Self, Bag};
use sui::event;
use sui::table::{Self, Table};

public struct Registry has key {
    id: UID,
    /// Maps account address -> package cap
    /// Only one cap per account address/package
    package_caps: Bag,
    /// Maps package_address -> proof_type
    registered_proof_types: Table<address, TypeName>,
    /// Maps package_address -> allowed module names (as bytes)
    allowed_modules: Table<address, vector<vector<u8>>>,
    /// Tracks batch execution state to enforce callback ordering
    batch_execution: Table<vector<u8>, BatchExecutionState>,
    /// Tracks completed batches for predecessor validation
    completed_batches: Table<vector<u8>, bool>,
}

/// Tracks execution progress of a batch to enforce MCMS operation ordering
public struct BatchExecutionState has store {
    total_callbacks: u64,
    next_expected_sequence: u64,
}

/// `ExecutingCallbackParams` is created when an operation is ready to be executed from MCMS
public struct ExecutingCallbackParams {
    target: address,
    module_name: String,
    function_name: String,
    data: vector<u8>,
    batch_id: vector<u8>,
    sequence_number: u64,
    total_in_batch: u64,
    expected_proof_type: TypeName,
}

public struct EntrypointRegistered has copy, drop {
    registry_id: ID,
    account_address: address,
    allowed_modules: vector<vector<u8>>,
}

const EPackageCapAlreadyRegistered: u64 = 1;
const EPackageCapNotRegistered: u64 = 2;
const EPackageIdMismatch: u64 = 3;
const EOutOfOrderExecution: u64 = 5;
const EWrongProofType: u64 = 6;
const EPackageNotRegistered: u64 = 7;
const EModuleNotRegistered: u64 = 8;
const EModuleNotAllowed: u64 = 9;

public struct MCMS_REGISTRY has drop {}

fun init(_witness: MCMS_REGISTRY, ctx: &mut TxContext) {
    let registry = Registry {
        id: object::new(ctx),
        package_caps: bag::new(ctx),
        registered_proof_types: table::new(ctx),
        allowed_modules: table::new(ctx),
        batch_execution: table::new(ctx),
        completed_batches: table::new(ctx),
    };

    transfer::share_object(registry);
}

fun enforce_execution_order(
    registry: &mut Registry,
    batch_id: vector<u8>,
    sequence_number: u64,
    total_in_batch: u64,
) {
    if (!registry.batch_execution.contains(batch_id)) {
        registry
            .batch_execution
            .add(
                batch_id,
                BatchExecutionState {
                    total_callbacks: total_in_batch,
                    next_expected_sequence: 0,
                },
            );
    };

    let state = registry.batch_execution.borrow_mut(batch_id);
    assert!(sequence_number == state.next_expected_sequence, EOutOfOrderExecution);

    state.next_expected_sequence = state.next_expected_sequence + 1;

    // When batch completes, mark as completed and clean up execution state
    if (state.next_expected_sequence == state.total_callbacks) {
        registry.completed_batches.add(batch_id, true);
        let BatchExecutionState { total_callbacks: _, next_expected_sequence: _ } = registry
            .batch_execution
            .remove(batch_id);
    }
}

public fun register_entrypoint<T: drop, C: key + store>(
    registry: &mut Registry,
    _proof: T,
    package_cap: C,
    allowed_modules: vector<vector<u8>>,
    _ctx: &mut TxContext,
) {
    let proof_type = type_name::with_original_ids<T>();
    let (proof_account_address, _proof_module_name) = params::get_account_address_and_module_name(
        proof_type,
    );

    assert!(!registry.package_caps.contains(proof_account_address), EPackageCapAlreadyRegistered);

    // Register package cap for package address
    registry.package_caps.add(proof_account_address, package_cap);

    // Register proof type for package address
    registry.registered_proof_types.add(proof_account_address, proof_type);

    // Register allowed modules for package address
    registry.allowed_modules.add(proof_account_address, allowed_modules);

    event::emit(EntrypointRegistered {
        registry_id: object::id(registry),
        account_address: proof_account_address,
        allowed_modules,
    });
}

public fun get_callback_params_with_caps<T: drop, C: key + store>(
    registry: &mut Registry,
    _proof: T,
    params: ExecutingCallbackParams,
): (&C, String, vector<u8>) {
    let ExecutingCallbackParams {
        target,
        module_name,
        function_name,
        data,
        batch_id,
        sequence_number,
        total_in_batch,
        expected_proof_type,
    } = params;

    enforce_execution_order(registry, batch_id, sequence_number, total_in_batch);

    let proof_type = type_name::with_original_ids<T>();
    let (proof_account_address, _proof_module_name) = params::get_account_address_and_module_name(
        proof_type,
    );

    assert!(proof_type == expected_proof_type, EWrongProofType);
    assert!(target == proof_account_address, EPackageIdMismatch);

    // Validate the proof comes from same package ID
    assert!(registry.package_caps.contains(proof_account_address), EPackageCapNotRegistered);
    assert!(registry.allowed_modules.contains(proof_account_address), EPackageNotRegistered);

    // Validate that the `module_name` is in the allowed modules list
    let allowed = registry.allowed_modules.borrow(proof_account_address);
    assert!(allowed.contains(module_name.as_bytes()), EModuleNotAllowed);

    let package_cap = registry.package_caps.borrow(proof_account_address);
    (package_cap, function_name, data)
}

public fun release_cap<T: drop, C: key + store>(registry: &mut Registry, _witness: T): C {
    let proof_type = type_name::with_original_ids<T>();
    let (proof_account_address, _) = params::get_account_address_and_module_name(
        proof_type,
    );

    assert!(registry.package_caps.contains(proof_account_address), EPackageCapNotRegistered);
    assert!(registry.registered_proof_types.contains(proof_account_address), EModuleNotRegistered);

    let expected_type = registry.registered_proof_types.borrow(proof_account_address);
    assert!(proof_type == *expected_type, EWrongProofType);

    registry.package_caps.remove(proof_account_address)
}

public(package) fun borrow_owner_cap<C: key + store>(registry: &Registry): &C {
    registry.package_caps.borrow(get_multisig_address())
}

public fun get_callback_params<T: drop>(
    registry: &mut Registry,
    params: ExecutingCallbackParams,
    _proof: T,
): (address, String, String, vector<u8>) {
    let ExecutingCallbackParams {
        target,
        module_name,
        function_name,
        data,
        batch_id,
        sequence_number,
        total_in_batch,
        expected_proof_type,
    } = params;

    enforce_execution_order(registry, batch_id, sequence_number, total_in_batch);

    let proof_type = type_name::with_original_ids<T>();
    let (proof_account_address, _) = params::get_account_address_and_module_name(
        proof_type,
    );

    assert!(target == proof_account_address, EPackageIdMismatch);

    // Validate the proof type matches the expected proof type
    assert!(proof_type == expected_proof_type, EWrongProofType);

    (target, module_name, function_name, data)
}

public(package) fun get_callback_params_from_mcms(
    registry: &mut Registry,
    params: ExecutingCallbackParams,
): (address, String, String, vector<u8>, TypeName) {
    let ExecutingCallbackParams {
        target,
        module_name,
        function_name,
        data,
        batch_id,
        sequence_number,
        total_in_batch,
        expected_proof_type,
    } = params;

    enforce_execution_order(registry, batch_id, sequence_number, total_in_batch);

    (target, module_name, function_name, data, expected_proof_type)
}

public(package) fun create_executing_callback_params(
    target: address,
    module_name: String,
    function_name: String,
    data: vector<u8>,
    batch_id: vector<u8>,
    sequence_number: u64,
    total_in_batch: u64,
    expected_proof_type: TypeName,
): ExecutingCallbackParams {
    ExecutingCallbackParams {
        target,
        module_name,
        function_name,
        data,
        batch_id,
        sequence_number,
        total_in_batch,
        expected_proof_type,
    }
}

public fun is_package_registered(registry: &Registry, package_address: address): bool {
    registry.package_caps.contains(package_address)
}

public(package) fun get_registered_proof_type(
    registry: &Registry,
    package_address: address,
): TypeName {
    assert!(registry.registered_proof_types.contains(package_address), EPackageNotRegistered);
    *registry.registered_proof_types.borrow(package_address)
}

public fun target(params: &ExecutingCallbackParams): address {
    params.target
}

public fun module_name(params: &ExecutingCallbackParams): String {
    params.module_name
}

public fun function_name(params: &ExecutingCallbackParams): String {
    params.function_name
}

public fun data(params: &ExecutingCallbackParams): vector<u8> {
    params.data
}

/// Check if a batch has been completed (all callbacks executed in order)
public fun is_batch_completed(registry: &Registry, batch_id: vector<u8>): bool {
    registry.completed_batches.contains(batch_id)
}

/// Get the next expected sequence number for a batch
public fun get_next_expected_sequence(registry: &Registry, batch_id: vector<u8>): u64 {
    if (!registry.batch_execution.contains(batch_id)) {
        return 0
    };
    registry.batch_execution.borrow(batch_id).next_expected_sequence
}

public fun get_multisig_address(): address {
    address::from_ascii_bytes(
        &type_name::with_defining_ids<McmsProof>().address_string().into_bytes(),
    )
}

public struct McmsProof has drop {}

public(package) fun create_mcms_proof(): McmsProof {
    McmsProof {}
}

// ===================== TESTS =====================

#[test_only]
/// Initialize the registry for testing
public fun test_init(ctx: &mut TxContext) {
    init(MCMS_REGISTRY {}, ctx)
}

#[test_only]
/// Create executing callback params for testing
public fun test_create_executing_callback_params(
    target: address,
    module_name: String,
    function_name: String,
    data: vector<u8>,
    batch_id: vector<u8>,
    sequence_number: u64,
    total_in_batch: u64,
    expected_proof_type: TypeName,
): ExecutingCallbackParams {
    create_executing_callback_params(
        target,
        module_name,
        function_name,
        data,
        batch_id,
        sequence_number,
        total_in_batch,
        expected_proof_type,
    )
}
