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

public struct PackageHistory has copy, drop, store {
    package_id: address,
    version: u64,
    timestamp: u64,
}

public struct FunctionKey has copy, drop, store {
    module_name: String,
    function_name: String,
}

public struct UpgradeRegistry has key, store {
    id: UID,
    // (module_name, function_name) -> blocked function versions
    function_restrictions: Table<FunctionKey, vector<u64>>,
    //  module_name -> blocked versions
    module_restrictions: Table<String, vector<u64>>,
    package_history: Table<String, vector<PackageHistory>>,
}

public fun initialize(ref: &mut CCIPObjectRef, owner_cap: &OwnerCap, ctx: &mut TxContext) {
    let registry = UpgradeRegistry {
        id: object::new(ctx),
        function_restrictions: table::new(ctx),
        module_restrictions: table::new(ctx),
        package_history: table::new(ctx),
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

// =================== Package History =================== //

// To be used by off-chain systems to get the package history of a given module
public fun get_package_history(
    ref: &CCIPObjectRef,
    package_name: String,
): (vector<address>, vector<u64>, vector<u64>) {
    let registry = state_object::borrow<UpgradeRegistry>(ref);

    if (!registry.package_history.contains(package_name)) {
        (vector::empty(), vector::empty(), vector::empty())
    } else {
        let package_history = registry.package_history.borrow(package_name);
        let package_ids = package_history.map_ref!(|ph| ph.package_id);
        let versions = package_history.map_ref!(|ph| ph.version);
        let timestamps = package_history.map_ref!(|ph| ph.timestamp);
        (package_ids, versions, timestamps)
    }
}
