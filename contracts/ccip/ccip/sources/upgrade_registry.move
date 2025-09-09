module ccip::upgrade_registry;

use ccip::ownable::OwnerCap;
use ccip::state_object::{Self, CCIPObjectRef};
use std::string::String;
use sui::event;
use sui::table::{Self, Table};

public struct VersionBlocked has copy, drop {
    module_name: String,
    version: u8,
}

public struct FunctionBlocked has copy, drop {
    module_name: String,
    function_name: String,
    version: u8,
}

const EFunctionNotAllowed: u64 = 1;

public struct UpgradeRegistry has key, store {
    id: UID,
    // module_name -> vector[vector<u8>, vector<u8>, ...]
    // the module_name represents the module under which the function is located, e.g. "fee_quoter", "offramp", etc.
    // the outer vector includes all the blocked function versions for this given module_name
    // the inner vectors can be:
    //  1. vector with a single u8, representing an entire version is blocked
    //  2. vector with multiple u8s, representing a version following by the function name, e.g. [1, b"get_fee"]
    //     this means v1 of "get_fee" is blocked
    function_restrictions: Table<String, vector<vector<u8>>>,
}

public fun initialize(ref: &mut CCIPObjectRef, owner_cap: &OwnerCap, ctx: &mut TxContext) {
    let registry = UpgradeRegistry {
        id: object::new(ctx),
        function_restrictions: table::new(ctx),
    };

    state_object::add(ref, owner_cap, registry, ctx);
}

// =================== Function Restrictions =================== //

public fun block_version(
    ref: &mut CCIPObjectRef,
    _: &OwnerCap,
    module_name: String,
    version: u8,
    _: &mut TxContext,
) {
    let registry = state_object::borrow_mut<UpgradeRegistry>(ref);
    if (!registry.function_restrictions.contains(module_name)) {
        registry.function_restrictions.add(module_name, vector[]);
    };
    registry.function_restrictions.borrow_mut(module_name).push_back(vector[version]);
    event::emit(VersionBlocked {
        module_name,
        version,
    });
}

public fun block_function(
    ref: &mut CCIPObjectRef,
    _: &OwnerCap,
    module_name: String,
    function_name: String,
    version: u8,
    _: &mut TxContext,
) {
    let registry = state_object::borrow_mut<UpgradeRegistry>(ref);
    if (!registry.function_restrictions.contains(module_name)) {
        registry.function_restrictions.add(module_name, vector[]);
    };
    let mut blocked_function = vector[version];
    blocked_function.append(function_name.into_bytes());
    registry.function_restrictions.borrow_mut(module_name).push_back(blocked_function);
    event::emit(FunctionBlocked {
        module_name,
        function_name,
        version,
    });
}

public fun get_module_restrictions(ref: &CCIPObjectRef, module_name: String): vector<vector<u8>> {
    let registry = state_object::borrow<UpgradeRegistry>(ref);

    if (!registry.function_restrictions.contains(module_name)) {
        vector::empty()
    } else {
        *registry.function_restrictions.borrow(module_name)
    }
}

// if this entire module is allowed, and this function is allowed, return true
public fun is_function_allowed(
    ref: &CCIPObjectRef,
    module_name: String,
    function_name: String,
    version: u8,
): bool {
    let registry = state_object::borrow<UpgradeRegistry>(ref);

    if (!registry.function_restrictions.contains(module_name)) {
        return true
    };

    let blocked_functions = registry.function_restrictions.borrow(module_name);
    let v = vector[version];
    let mut function_name_bytes = vector[];
    function_name_bytes.push_back(version);
    function_name_bytes.append(function_name.into_bytes());

    !blocked_functions.contains(&function_name_bytes) && !blocked_functions.contains(&v)
}

public fun verify_function_allowed(
    ref: &CCIPObjectRef,
    module_name: String,
    function_name: String,
    version: u8,
) {
    assert!(
        is_function_allowed(
            ref,
            module_name,
            function_name,
            version,
        ),
        EFunctionNotAllowed,
    );
}
