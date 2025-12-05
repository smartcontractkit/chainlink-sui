module mcms::mcms_registry;

use mcms::params;
use std::ascii;
use std::string::String;
use std::type_name::{Self, TypeName};
use sui::address;
use sui::bag::{Self, Bag};
use sui::event;
use sui::package::Publisher;
use sui::table::{Self, Table};

public struct Registry has key {
    id: UID,
    /// Maps account address -> package cap
    /// Only one cap per account address/package
    package_caps: Bag,
    /// Maps package_address -> proof_type
    registered_proof_types: Table<ascii::String, TypeName>,
    /// Reverse lookup of proof type -> package address
    proof_type_to_package: Table<TypeName, ascii::String>,
    /// Maps package_address -> allowed module names (as bytes)
    allowed_modules: Table<ascii::String, vector<vector<u8>>>,
    /// Tracks batch execution state to enforce callback ordering
    batch_execution: Table<vector<u8>, BatchExecutionState>,
    /// Tracks completed batches for predecessor validation
    completed_batches: Table<vector<u8>, bool>,
}

/// Wrapper for Publisher validation that encodes proof type in type parameter
/// This allows compile-time enforcement that Publisher matches the proof type
public struct PublisherWrapper<phantom T: drop> {
    package_address: ascii::String,
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
}

public struct EntrypointRegistered has copy, drop {
    registry_id: ID,
    account_address: ascii::String,
    allowed_modules: vector<vector<u8>>,
    proof_type: TypeName,
}

public struct ModulesAdded has copy, drop {
    registry_id: ID,
    package_address: ascii::String,
    module_names: vector<vector<u8>>,
}

public struct ModulesRemoved has copy, drop {
    registry_id: ID,
    package_address: ascii::String,
    module_names: vector<vector<u8>>,
}

const EPackageCapAlreadyRegistered: u64 = 1;
const EPackageCapNotRegistered: u64 = 2;
const EPackageIdMismatch: u64 = 3;
const EOutOfOrderExecution: u64 = 5;
const EWrongProofType: u64 = 6;
const EPackageNotRegistered: u64 = 7;
const EModuleNotAllowed: u64 = 8;
const EModuleAlreadyAllowed: u64 = 9;
const EModuleNotInAllowlist: u64 = 10;
const EOnlyAcceptOwnershipAllowed: u64 = 11;
const EOnlyMcmsAcceptOwnershipProofAllowed: u64 = 12;
const EInvalidModuleName: u64 = 13;
const EProofNotAtPublisherAddressAndModule: u64 = 14;
const EProofTypeNotRegistered: u64 = 15;
const ECapAddressMismatch: u64 = 16;

public struct MCMS_REGISTRY has drop {}

fun init(_witness: MCMS_REGISTRY, ctx: &mut TxContext) {
    let registry = Registry {
        id: object::new(ctx),
        package_caps: bag::new(ctx),
        registered_proof_types: table::new(ctx),
        proof_type_to_package: table::new(ctx),
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

/// Having reference to Publisher means you have access to `Publisher` object.
/// This is only sent to the package deployer, therefore we know only the owner can call this.
public fun create_publisher_wrapper<T: drop>(
    publisher: &Publisher,
    _proof: T,
): PublisherWrapper<T> {
    assert!(publisher.from_module<T>(), EProofNotAtPublisherAddressAndModule);
    PublisherWrapper<T> { package_address: *publisher.package() }
}

/// Registers a package with MCMS.
/// `PublisherWrapper` asserts that the proof is at the publisher address and module.
/// `C` must be the same package as the PublisherWrapper.
public fun register_entrypoint<T: drop, C: key + store>(
    registry: &mut Registry,
    publisher_wrapper: PublisherWrapper<T>,
    _proof: T,
    package_cap: C,
    allowed_modules: vector<vector<u8>>,
    _ctx: &mut TxContext,
) {
    let PublisherWrapper { package_address } = publisher_wrapper;
    let proof_type = type_name::with_original_ids<T>();

    let cap_address = type_name::with_original_ids<C>().address_string();
    assert!(cap_address == package_address, ECapAddressMismatch);

    // Assert publisher is not already registered
    assert!(!registry.package_caps.contains(package_address), EPackageCapAlreadyRegistered);

    // Register package cap directly (Publisher stays in OwnerCap)
    registry.package_caps.add(package_address, package_cap);

    // Register proof type for package address
    registry.registered_proof_types.add(package_address, proof_type);

    // Register proof type to package address (reverse lookup)
    registry.proof_type_to_package.add(proof_type, package_address);

    // Register allowed modules for package address
    registry.allowed_modules.add(package_address, allowed_modules);

    event::emit(EntrypointRegistered {
        registry_id: object::id(registry),
        account_address: package_address,
        allowed_modules,
        proof_type,
    });
}

/// Add a new module to the allowed modules list for a registered package.
public fun add_allowed_modules<T: drop>(
    registry: &mut Registry,
    _proof: T,
    new_allowed_modules: vector<vector<u8>>,
    _ctx: &mut TxContext,
) {
    let proof_type = type_name::with_original_ids<T>();
    let proof_account_address = proof_type.address_string();

    // Validate the package is registered
    assert!(registry.allowed_modules.contains(proof_account_address), EPackageNotRegistered);

    // Validate proof type matches the expected proof type
    let expected_proof_type = *registry.registered_proof_types.borrow(proof_account_address);
    assert!(proof_type == expected_proof_type, EWrongProofType);

    let allowed_modules = registry.allowed_modules.borrow_mut(proof_account_address);
    let mut i = 0;
    while (i < new_allowed_modules.length()) {
        let new_module = new_allowed_modules[i];
        assert!(!allowed_modules.contains(&new_module), EModuleAlreadyAllowed);
        allowed_modules.push_back(new_module);
        i = i + 1;
    };

    event::emit(ModulesAdded {
        registry_id: object::id(registry),
        package_address: proof_account_address,
        module_names: new_allowed_modules,
    });
}

/// Remove modules from the allowed modules list for a registered package.
public fun remove_allowed_modules<T: drop>(
    registry: &mut Registry,
    _proof: T,
    modules_to_remove: vector<vector<u8>>,
    _ctx: &mut TxContext,
) {
    let proof_type = type_name::with_original_ids<T>();
    let proof_account_address = proof_type.address_string();

    // Validate the package is registered
    assert!(registry.allowed_modules.contains(proof_account_address), EPackageNotRegistered);

    // Validate proof type matches the expected proof type
    let expected_proof_type = *registry.registered_proof_types.borrow(proof_account_address);
    assert!(proof_type == expected_proof_type, EWrongProofType);

    let allowed_modules = registry.allowed_modules.borrow_mut(proof_account_address);

    let mut i = 0;
    while (i < modules_to_remove.length()) {
        let (found, index) = allowed_modules.index_of(&modules_to_remove[i]);
        assert!(found, EModuleNotInAllowlist);

        allowed_modules.remove(index);
        i = i + 1;
    };

    event::emit(ModulesRemoved {
        registry_id: object::id(registry),
        package_address: proof_account_address,
        module_names: modules_to_remove,
    });
}

public fun get_callback_params_with_caps<T: drop, C: key + store>(
    registry: &mut Registry,
    _proof: T,
    params: ExecutingCallbackParams,
): (&mut C, String, vector<u8>) {
    let ExecutingCallbackParams {
        target,
        module_name,
        function_name,
        data,
        batch_id,
        sequence_number,
        total_in_batch,
    } = params;

    enforce_execution_order(registry, batch_id, sequence_number, total_in_batch);

    let proof_type = type_name::with_original_ids<T>();
    let proof_account_address = proof_type.address_string();

    let expected_proof_type = get_registered_proof_type(registry, proof_account_address);
    assert!(proof_type == expected_proof_type, EWrongProofType);
    assert!(registry.proof_type_to_package.contains(proof_type), EProofTypeNotRegistered);
    assert!(target.to_ascii_string() == proof_account_address, EPackageIdMismatch);

    // Validate the proof comes from same package ID
    assert!(registry.package_caps.contains(proof_account_address), EPackageCapNotRegistered);
    assert!(registry.allowed_modules.contains(proof_account_address), EPackageNotRegistered);

    // Validate that the `module_name` is in the allowed modules list
    let allowed = registry.allowed_modules.borrow(proof_account_address);
    assert!(allowed.contains(module_name.as_bytes()), EModuleNotAllowed);

    let package_cap = registry.package_caps.borrow_mut(proof_account_address);
    (package_cap, function_name, data)
}

/// Release the package cap for a registered package.
/// All four states from `Registry` are removed:
/// - package_caps
/// - registered_proof_types
/// - proof_type_to_package
/// - allowed_modules
public fun release_cap<T: drop, C: key + store>(registry: &mut Registry, _proof: T): C {
    let proof_type = type_name::with_original_ids<T>();
    let proof_account_address = proof_type.address_string();

    // Assert the package is registered
    assert!(registry.package_caps.contains(proof_account_address), EPackageCapNotRegistered);
    assert!(registry.proof_type_to_package.contains(proof_type), EProofTypeNotRegistered);
    assert!(registry.registered_proof_types.contains(proof_account_address), EPackageNotRegistered);

    let expected_type = registry.registered_proof_types.remove(proof_account_address);
    assert!(proof_type == expected_type, EWrongProofType);

    let cap = registry.package_caps.remove(proof_account_address);
    registry.proof_type_to_package.remove(proof_type);
    registry.allowed_modules.remove(proof_account_address);

    cap
}

public(package) fun borrow_owner_cap<C: key + store>(registry: &Registry): &C {
    registry.package_caps.borrow(get_multisig_address_ascii())
}

/// Only proof with struct name "McmsAcceptOwnershipProof" are allowed to be used with this function.
/// Validate the target, module name and function name are as expected.
/// This is only ever called by `mcms_accept_ownership`, precursor to the package being registered with MCMS
/// ExecutingCallbackParams can only be created by MCMS
public fun get_accept_ownership_data<T: drop>(
    registry: &mut Registry,
    params: ExecutingCallbackParams,
    _proof: T,
): vector<u8> {
    let ExecutingCallbackParams {
        target,
        module_name,
        function_name,
        data,
        batch_id,
        sequence_number,
        total_in_batch,
    } = params;

    enforce_execution_order(registry, batch_id, sequence_number, total_in_batch);

    let proof_type = type_name::with_original_ids<T>();
    let proof_account_address = proof_type.address_string();
    let proof_module_name = proof_type.module_string();
    let struct_name = params::get_struct_name(&proof_type);

    assert!(target.to_ascii_string() == proof_account_address, EPackageIdMismatch);
    assert!(proof_module_name.to_string() == module_name, EInvalidModuleName);
    assert!(function_name.as_bytes() == b"accept_ownership", EOnlyAcceptOwnershipAllowed);
    assert!(struct_name == b"McmsAcceptOwnershipProof", EOnlyMcmsAcceptOwnershipProofAllowed);

    data
}

public(package) fun get_callback_params_from_mcms(
    registry: &mut Registry,
    params: ExecutingCallbackParams,
): (address, String, String, vector<u8>) {
    let ExecutingCallbackParams {
        target,
        module_name,
        function_name,
        data,
        batch_id,
        sequence_number,
        total_in_batch,
    } = params;

    enforce_execution_order(registry, batch_id, sequence_number, total_in_batch);

    (target, module_name, function_name, data)
}

public(package) fun create_executing_callback_params(
    target: address,
    module_name: String,
    function_name: String,
    data: vector<u8>,
    batch_id: vector<u8>,
    sequence_number: u64,
    total_in_batch: u64,
): ExecutingCallbackParams {
    ExecutingCallbackParams {
        target,
        module_name,
        function_name,
        data,
        batch_id,
        sequence_number,
        total_in_batch,
    }
}

public fun is_package_registered(registry: &Registry, package_address: ascii::String): bool {
    registry.package_caps.contains(package_address)
}

public(package) fun get_registered_proof_type(
    registry: &Registry,
    package_address: ascii::String,
): TypeName {
    assert!(registry.registered_proof_types.contains(package_address), EPackageNotRegistered);
    *registry.registered_proof_types.borrow(package_address)
}

/// Returns the list of allowed module names for a registered package
public fun get_allowed_modules(
    registry: &Registry,
    package_address: ascii::String,
): vector<vector<u8>> {
    assert!(registry.allowed_modules.contains(package_address), EPackageNotRegistered);
    *registry.allowed_modules.borrow(package_address)
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
    address::from_ascii_bytes(&get_multisig_address_ascii().into_bytes())
}

public fun get_multisig_address_ascii(): ascii::String {
    type_name::with_defining_ids<MCMS_REGISTRY>().address_string()
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
): ExecutingCallbackParams {
    create_executing_callback_params(
        target,
        module_name,
        function_name,
        data,
        batch_id,
        sequence_number,
        total_in_batch,
    )
}

#[test_only]
public fun test_get_cap_address<C: key + store>(
    registry: &Registry,
    package_address: ascii::String,
): address {
    assert!(registry.package_caps.contains(package_address), EPackageCapNotRegistered);
    let cap: &C = registry.package_caps.borrow(package_address);
    object::id_address(cap)
}

#[test_only]
public fun get_cap<C: key + store>(registry: &Registry, package_address: ascii::String): &C {
    registry.package_caps.borrow(package_address)
}
