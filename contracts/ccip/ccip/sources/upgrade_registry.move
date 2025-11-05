module ccip::upgrade_registry;

use ccip::ownable::OwnerCap;
use ccip::state_object::{Self, CCIPObjectRef};
use mcms::bcs_stream;
use mcms::mcms_registry::{Self, Registry, ExecutingCallbackParams};
use std::string::{Self, String};
use sui::event;
use sui::table::{Self, Table};

public struct VersionBlocked has copy, drop {
    module_name: String,
    version: u8,
}

public struct VersionUnblocked has copy, drop {
    module_name: String,
    version: u8,
}

public struct FunctionBlocked has copy, drop {
    module_name: String,
    function_name: String,
    version: u8,
}

public struct FunctionUnblocked has copy, drop {
    module_name: String,
    function_name: String,
    version: u8,
}

const EFunctionNotAllowed: u64 = 1;
const EInvalidOwnerCap: u64 = 2;
const EAlreadyInitialized: u64 = 3;
const EInvalidFunction: u64 = 4;

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
    assert!(object::id(owner_cap) == state_object::owner_cap_id(ref), EInvalidOwnerCap);
    assert!(!state_object::contains<UpgradeRegistry>(ref), EAlreadyInitialized);
    let registry = UpgradeRegistry {
        id: object::new(ctx),
        function_restrictions: table::new(ctx),
    };

    state_object::add(ref, owner_cap, registry, ctx);
}

// =================== Function Restrictions =================== //

public fun block_version(
    ref: &mut CCIPObjectRef,
    owner_cap: &OwnerCap,
    module_name: String,
    version: u8,
    _: &mut TxContext,
) {
    assert!(object::id(owner_cap) == state_object::owner_cap_id(ref), EInvalidOwnerCap);

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

public fun unblock_version(
    ref: &mut CCIPObjectRef,
    owner_cap: &OwnerCap,
    module_name: String,
    version: u8,
    _: &mut TxContext,
) {
    assert!(object::id(owner_cap) == state_object::owner_cap_id(ref), EInvalidOwnerCap);

    let registry = state_object::borrow_mut<UpgradeRegistry>(ref);
    if (!registry.function_restrictions.contains(module_name)) {
        return
    };
    let blocked_versions = registry.function_restrictions.borrow_mut(module_name);
    let mut i = 0;
    while (i < blocked_versions.length()) {
        let blocked_version = &blocked_versions[i];
        if (blocked_version[0] == version) {
            blocked_versions.swap_remove(i);
            event::emit(VersionUnblocked {
                module_name,
                version,
            });
            return
        };
        i = i + 1;
    };
}

public fun block_function(
    ref: &mut CCIPObjectRef,
    owner_cap: &OwnerCap,
    module_name: String,
    function_name: String,
    version: u8,
    _: &mut TxContext,
) {
    assert!(object::id(owner_cap) == state_object::owner_cap_id(ref), EInvalidOwnerCap);

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

public fun unblock_function(
    ref: &mut CCIPObjectRef,
    owner_cap: &OwnerCap,
    module_name: String,
    function_name: String,
    version: u8,
    _: &mut TxContext,
) {
    assert!(object::id(owner_cap) == state_object::owner_cap_id(ref), EInvalidOwnerCap);

    let registry = state_object::borrow_mut<UpgradeRegistry>(ref);
    if (!registry.function_restrictions.contains(module_name)) {
        return
    };
    let blocked_functions = registry.function_restrictions.borrow_mut(module_name);
    let mut unblock_function = vector[version];
    unblock_function.append(function_name.into_bytes());
    let mut i = 0;
    while (i < blocked_functions.length()) {
        let blocked_function = &blocked_functions[i];
        if (blocked_function == unblock_function) {
            blocked_functions.swap_remove(i);
            event::emit(FunctionUnblocked {
                module_name,
                function_name,
                version,
            });
            return
        };
        i = i + 1;
    };
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

// =================== MCMS Functions =================== //

public struct McmsCallback has drop {}

public fun mcms_block_version(
    ref: &mut CCIPObjectRef,
    registry: &mut Registry,
    params: ExecutingCallbackParams,
    ctx: &mut TxContext,
) {
    let (owner_cap, function, data) = mcms_registry::get_callback_params_with_caps<
        McmsCallback,
        OwnerCap,
    >(
        registry,
        McmsCallback {},
        params,
    );
    assert!(function == string::utf8(b"block_version"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    bcs_stream::validate_obj_addrs(
        vector[object::id_address(ref), object::id_address(owner_cap)],
        &mut stream,
    );

    let module_name = bcs_stream::deserialize_string(&mut stream);
    let version = bcs_stream::deserialize_u8(&mut stream);
    bcs_stream::assert_is_consumed(&stream);

    block_version(ref, owner_cap, module_name, version, ctx);
}

public fun mcms_unblock_version(
    ref: &mut CCIPObjectRef,
    registry: &mut Registry,
    params: ExecutingCallbackParams,
    ctx: &mut TxContext,
) {
    let (owner_cap, function, data) = mcms_registry::get_callback_params_with_caps<
        McmsCallback,
        OwnerCap,
    >(
        registry,
        McmsCallback {},
        params,
    );
    assert!(function == string::utf8(b"unblock_version"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    bcs_stream::validate_obj_addrs(
        vector[object::id_address(ref), object::id_address(owner_cap)],
        &mut stream,
    );

    let module_name = bcs_stream::deserialize_string(&mut stream);
    let version = bcs_stream::deserialize_u8(&mut stream);
    bcs_stream::assert_is_consumed(&stream);

    unblock_version(ref, owner_cap, module_name, version, ctx);
}

public fun mcms_block_function(
    ref: &mut CCIPObjectRef,
    registry: &mut Registry,
    params: ExecutingCallbackParams,
    ctx: &mut TxContext,
) {
    let (owner_cap, function, data) = mcms_registry::get_callback_params_with_caps<
        McmsCallback,
        OwnerCap,
    >(
        registry,
        McmsCallback {},
        params,
    );
    assert!(function == string::utf8(b"block_function"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    bcs_stream::validate_obj_addrs(
        vector[object::id_address(ref), object::id_address(owner_cap)],
        &mut stream,
    );

    let module_name = bcs_stream::deserialize_string(&mut stream);
    let function_name = bcs_stream::deserialize_string(&mut stream);
    let version = bcs_stream::deserialize_u8(&mut stream);
    bcs_stream::assert_is_consumed(&stream);

    block_function(ref, owner_cap, module_name, function_name, version, ctx);
}

public fun mcms_unblock_function(
    ref: &mut CCIPObjectRef,
    registry: &mut Registry,
    params: ExecutingCallbackParams,
    ctx: &mut TxContext,
) {
    let (owner_cap, function, data) = mcms_registry::get_callback_params_with_caps<
        McmsCallback,
        OwnerCap,
    >(
        registry,
        McmsCallback {},
        params,
    );
    assert!(function == string::utf8(b"unblock_function"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    bcs_stream::validate_obj_addrs(
        vector[object::id_address(ref), object::id_address(owner_cap)],
        &mut stream,
    );

    let module_name = bcs_stream::deserialize_string(&mut stream);
    let function_name = bcs_stream::deserialize_string(&mut stream);
    let version = bcs_stream::deserialize_u8(&mut stream);
    bcs_stream::assert_is_consumed(&stream);

    unblock_function(ref, owner_cap, module_name, function_name, version, ctx);
}
