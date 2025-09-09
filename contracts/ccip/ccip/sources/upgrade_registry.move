module ccip::upgrade_registry;

use ccip::ownable::OwnerCap;
use ccip::state_object::{Self, CCIPObjectRef};
use std::string::String;
use sui::event;
use sui::table::{Self, Table};

public struct FunctionRestrictionsUpdated has copy, drop {
    module_name: String,
    function_name: String,
    blocked_versions: vector<u64>,
}

public struct ModuleRestrictionsUpdated has copy, drop {
    module_name: String,
    blocked_versions: vector<u64>,
}

public struct FunctionKey has copy, drop, store {
    module_name: String,
    function_name: String,
}

const EFunctionNotAllowed: u64 = 1;

public struct UpgradeRegistry has key, store {
    id: UID,
    // (module_name, function_name) -> blocked function versions
    function_restrictions: Table<FunctionKey, vector<u64>>,
    //  module_name -> blocked versions
    module_restrictions: Table<String, vector<u64>>,
}

public fun initialize(ref: &mut CCIPObjectRef, owner_cap: &OwnerCap, ctx: &mut TxContext) {
    let registry = UpgradeRegistry {
        id: object::new(ctx),
        function_restrictions: table::new(ctx),
        module_restrictions: table::new(ctx),
    };

    state_object::add(ref, owner_cap, registry, ctx);
}

// =================== Function Restrictions =================== //

public fun update_function_restrictions(
    ref: &mut CCIPObjectRef,
    _: &OwnerCap,
    module_name: String,
    function_name: String,
    blocked_versions: vector<u64>,
    _: &mut TxContext,
) {
    let registry = state_object::borrow_mut<UpgradeRegistry>(ref);

    let key = FunctionKey { module_name, function_name };
    if (registry.function_restrictions.contains(key)) {
        registry.function_restrictions.remove(key);
    };
    registry.function_restrictions.add(key, blocked_versions);
    event::emit(FunctionRestrictionsUpdated {
        module_name,
        function_name,
        blocked_versions,
    });
}

public fun get_function_restrictions(
    ref: &CCIPObjectRef,
    module_name: String,
    function_name: String,
): vector<u64> {
    let registry = state_object::borrow<UpgradeRegistry>(ref);

    let key = FunctionKey { module_name, function_name };
    if (!registry.function_restrictions.contains(key)) {
        vector::empty()
    } else {
        *registry.function_restrictions.borrow(key)
    }
}

// if this entire module is allowed, and this function is allowed, return true
public fun is_function_allowed(
    ref: &CCIPObjectRef,
    module_name: String,
    function_name: String,
    contract_version: u64,
): bool {
    let registry = state_object::borrow<UpgradeRegistry>(ref);

    let key = FunctionKey { module_name, function_name };

    // if the module version is not allowed, then the function is not allowed
    if (!is_module_allowed(ref, module_name, contract_version)) {
        return false
    };

    // if the function version is not allowed, then the function is not allowed
    if (!registry.function_restrictions.contains(key)) {
        true
    } else {
        let blocked_versions = registry.function_restrictions.borrow(key);
        !blocked_versions.contains(&contract_version)
    }
}

public fun verify_function_allowed(
    ref: &CCIPObjectRef,
    module_name: String,
    function_name: String,
    contract_version: u64,
) {
    assert!(
        is_function_allowed(
            ref,
            module_name,
            function_name,
            contract_version,
        ),
        EFunctionNotAllowed,
    );
}

// =================== Module Restrictions =================== //

public fun update_module_restrictions(
    ref: &mut CCIPObjectRef,
    _: &OwnerCap,
    module_name: String,
    blocked_versions: vector<u64>,
    _: &mut TxContext,
) {
    let registry = state_object::borrow_mut<UpgradeRegistry>(ref);

    if (registry.module_restrictions.contains(module_name)) {
        registry.module_restrictions.remove(module_name);
    };
    registry.module_restrictions.add(module_name, blocked_versions);
    event::emit(ModuleRestrictionsUpdated {
        module_name,
        blocked_versions,
    });
}

public fun get_module_restrictions(ref: &CCIPObjectRef, module_name: String): vector<u64> {
    let registry = state_object::borrow<UpgradeRegistry>(ref);

    if (!registry.module_restrictions.contains(module_name)) {
        vector::empty()
    } else {
        *registry.module_restrictions.borrow(module_name)
    }
}

public fun is_module_allowed(
    ref: &CCIPObjectRef,
    module_name: String,
    contract_version: u64,
): bool {
    let registry = state_object::borrow<UpgradeRegistry>(ref);

    if (!registry.module_restrictions.contains(module_name)) {
        true
    } else {
        let blocked_versions = registry.module_restrictions.borrow(module_name);
        !blocked_versions.contains(&contract_version)
    }
}
