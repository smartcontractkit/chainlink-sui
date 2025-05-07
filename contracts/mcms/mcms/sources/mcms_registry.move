module mcms::mcms_registry;

use std::string::{Self, String};
use std::type_name::{Self, TypeName};
use sui::address;
use sui::bag::{Self, Bag};
use sui::event;
use sui::hex;
use sui::table::{Self, Table};

public struct Registry has key {
    id: UID,
    /// Maps module name -> ModuleInfo
    callback_modules: Table<String, RegisteredModule>,
    /// Maps account address -> module cap
    /// Only one cap per account address/package
    module_caps: Bag,
}

/// `ExecutingCallbackParams` is created when an operation is ready to be executed from MCMS
public struct ExecutingCallbackParams {
    target: address,
    module_name: String,
    function_name: String,
    data: vector<u8>,
}

public struct RegisteredModule has store {
    proof_type: TypeName,
}

public struct EntrypointRegistered has copy, drop {
    registry_id: ID,
    account_address: address,
    module_name: String,
}

const EModuleNotRegistered: u64 = 1;
const EModuleAlreadyRegistered: u64 = 2;
const EModuleCapAlreadyRegistered: u64 = 3;
const EWrongProofType: u64 = 4;
const EModuleNameMismatch: u64 = 5;
const EModuleCapNotRegistered: u64 = 6;
const ETargetNotRegistered: u64 = 7;

public struct MCMS_REGISTRY has drop {}

fun init(_witness: MCMS_REGISTRY, ctx: &mut TxContext) {
    let registry = Registry {
        id: object::new(ctx),
        callback_modules: table::new(ctx),
        module_caps: bag::new(ctx),
    };

    transfer::share_object(registry);
}

public fun register_entrypoint<T: drop, C: key + store>(
    registry: &mut Registry,
    _proof: T,
    module_cap: Option<C>,
    _ctx: &TxContext,
) {
    let proof_type = type_name::get<T>();
    let (proof_account_address, proof_module_name) = get_account_address_and_module_name(
        proof_type,
    );

    if (module_cap.is_some()) {
        assert!(!registry.module_caps.contains(proof_account_address), EModuleCapAlreadyRegistered);
    };

    assert!(!registry.callback_modules.contains(proof_module_name), EModuleAlreadyRegistered);

    // Register module cap for package address
    registry.module_caps.add(proof_account_address, module_cap.destroy_some());

    // Register proof type for module
    registry.callback_modules.add(proof_module_name, RegisteredModule { proof_type });

    event::emit(EntrypointRegistered {
        registry_id: object::id(registry),
        account_address: proof_account_address,
        module_name: proof_module_name,
    });
}

public fun get_callback_params<T: drop, C: key + store>(
    registry: &mut Registry,
    _witness: T,
    params: ExecutingCallbackParams,
): (&C, String, vector<u8>) {
    let ExecutingCallbackParams { target, module_name, function_name, data } = params;

    let proof_type = type_name::get<T>();
    let (proof_account_address, proof_module_name) = get_account_address_and_module_name(
        proof_type,
    );

    assert!(target == proof_account_address, ETargetNotRegistered);
    assert!(module_name == proof_module_name, EModuleNameMismatch);
    assert!(table::contains(&registry.callback_modules, module_name), EModuleNotRegistered);

    let expected_proof_type = registry.callback_modules.borrow(proof_module_name).proof_type;
    assert!(expected_proof_type == proof_type, EWrongProofType);

    assert!(registry.module_caps.contains(proof_account_address), EModuleCapNotRegistered);

    let module_cap = registry.module_caps.borrow(proof_account_address);
    (module_cap, function_name, data)
}

public fun release_cap<T: drop, C: key + store>(registry: &mut Registry, _witness: T): C {
    let proof_type = type_name::get<T>();
    let (proof_account_address, proof_module_name) = get_account_address_and_module_name(
        proof_type,
    );

    assert!(table::contains(&registry.callback_modules, proof_module_name), EModuleNotRegistered);

    assert!(registry.module_caps.contains(proof_account_address), EModuleCapNotRegistered);

    let RegisteredModule { proof_type: expected_proof_type } = registry
        .callback_modules
        .remove(proof_module_name);
    assert!(expected_proof_type == proof_type, EWrongProofType);

    registry.module_caps.remove(proof_account_address)
}

public(package) fun get_callback_params_for_mcms(
    params: ExecutingCallbackParams,
): (address, String, String, vector<u8>) {
    let ExecutingCallbackParams { target, module_name, function_name, data } = params;
    (target, module_name, function_name, data)
}

public(package) fun create_executing_callback_params(
    target: address,
    module_name: String,
    function_name: String,
    data: vector<u8>,
): ExecutingCallbackParams {
    ExecutingCallbackParams {
        target,
        module_name,
        function_name,
        data,
    }
}

fun get_account_address_and_module_name(proof_type: TypeName): (address, String) {
    let account_address_bytes = hex::decode(proof_type.get_address().into_bytes());
    let account_address = address::from_bytes(account_address_bytes);
    let module_name = string::from_ascii(proof_type.get_module());
    (account_address, module_name)
}

public fun is_package_registered(registry: &Registry, package_address: address): bool {
    registry.module_caps.contains(package_address)
}

public fun is_module_registered(registry: &Registry, module_name: String): bool {
    table::contains(&registry.callback_modules, module_name)
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
): ExecutingCallbackParams {
    create_executing_callback_params(target, module_name, function_name, data)
}
