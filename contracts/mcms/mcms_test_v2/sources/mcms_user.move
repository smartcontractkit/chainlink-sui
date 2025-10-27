module mcms_test::mcms_user;

use mcms::bcs_stream;
use mcms::mcms_deployer::{Self, DeployerState};
use mcms::mcms_registry::{Self, ExecutingCallbackParams, Registry};
use std::string::{Self, String};
use sui::package::{Self, UpgradeCap};
use mcms_test::ownable::{Self, OwnableState, OwnerCap};

const EInvalidAdminCap: u64 = 1;
const EInvalidFunction: u64 = 2;

public struct UserData has key, store {
    id: UID,
    invocations: u8,
    a: String,
    b: vector<u8>,
    c: address,
    d: u128,
    ownable_state: OwnableState,
}

public fun type_and_version(): String {
    string::utf8(b"MCMSUser 2.0.0")
}

public fun function_one(
    user_data: &mut UserData,
    owner_cap: &OwnerCap,
    arg1: String,
    arg2: vector<u8>,
) {
    assert_valid_owner_cap(user_data, owner_cap);

    user_data.invocations = user_data.invocations + 1;
    user_data.a = arg1;
    user_data.b = arg2;
}

public fun function_two(
    user_data: &mut UserData,
    owner_cap: &OwnerCap,
    arg1: address,
    arg2: u128,
) {
    assert_valid_owner_cap(user_data, owner_cap);

    user_data.invocations = user_data.invocations + 1;
    user_data.c = arg1;
    user_data.d = arg2;
}

public struct MCMS_USER has drop {}

fun init(otw: MCMS_USER, ctx: &mut TxContext) {
    let (ownable_state, mut owner_cap) = ownable::new(ctx);

    let publisher = package::claim(otw, ctx);
    ownable::attach_publisher(&mut owner_cap, publisher);

    transfer::share_object(UserData {
        id: object::new(ctx),
        invocations: 0,
        a: string::utf8(b""),
        b: vector[],
        c: @0x0,
        d: 0,
        ownable_state,
    });

    transfer::public_transfer(owner_cap, ctx.sender());
}

public fun register_mcms_entrypoint(
    owner_cap: OwnerCap,
    registry: &mut Registry,
    user_data: &UserData,
    ctx: &mut TxContext,
) {
    assert_valid_owner_cap(user_data, &owner_cap);

    // Create publisher wrapper
    let publisher_wrapper = mcms_registry::create_publisher_wrapper(
        ownable::borrow_publisher(&owner_cap),
        SampleMcmsCallback {},
    );

    // Transfer owner_cap to MCMS
    mcms_registry::register_entrypoint(
        registry,
        publisher_wrapper,
        SampleMcmsCallback {},
        owner_cap,
        vector[b"mcms_user"],
        ctx,
    );
}

public fun register_upgrade_cap(
    state: &mut DeployerState,
    upgrade_cap: UpgradeCap,
    registry: &mut Registry,
    ctx: &mut TxContext,
) {
    // Transfer upgrade permissions to MCMS
    mcms_deployer::register_upgrade_cap(
        state,
        registry,
        upgrade_cap,
        ctx,
    );
}

fun assert_valid_owner_cap(user_data: &UserData, owner_cap: &OwnerCap) {
    assert!(ownable::owner_cap_id(&user_data.ownable_state) == object::id(owner_cap), EInvalidAdminCap);
}

public struct SampleMcmsCallback has drop {}

public fun mcms_function_one(
    user_data: &mut UserData,
    registry: &mut Registry,
    params: ExecutingCallbackParams, // hot potato
    _ctx: &mut TxContext,
) {
    let (owner_cap, function, data) = mcms_registry::get_callback_params_with_caps<
        SampleMcmsCallback,
        OwnerCap,
    >(
        registry,
        SampleMcmsCallback {},
        params,
    );

    assert!(function == string::utf8(b"function_one"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    bcs_stream::validate_obj_addrs(
        vector[object::id_address(user_data), object::id_address(owner_cap)],
        &mut stream,
    );

    let arg1 = bcs_stream::deserialize_string(&mut stream);
    let arg2 = bcs_stream::deserialize_vector_u8(&mut stream);
    bcs_stream::assert_is_consumed(&stream);
    function_one(user_data, owner_cap, arg1, arg2);
}

public fun mcms_function_two(
    user_data: &mut UserData,
    registry: &mut Registry,
    params: ExecutingCallbackParams, // hot potato
    _ctx: &mut TxContext,
) {
    let (owner_cap, function, data) = mcms_registry::get_callback_params_with_caps<
        SampleMcmsCallback,
        OwnerCap,
    >(
        registry,
        SampleMcmsCallback {},
        params,
    );

    assert!(function == string::utf8(b"function_two"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    bcs_stream::validate_obj_addrs(
        vector[object::id_address(user_data), object::id_address(owner_cap)],
        &mut stream,
    );

    let arg1 = bcs_stream::deserialize_address(&mut stream);
    let arg2 = bcs_stream::deserialize_u128(&mut stream);
    bcs_stream::assert_is_consumed(&stream);
    function_two(user_data, owner_cap, arg1, arg2);
}

public fun get_owner_cap_id(user_data: &UserData): ID {
    ownable::owner_cap_id(&user_data.ownable_state)
}

public fun get_invocations(user_data: &UserData): u8 {
    user_data.invocations
}

public fun get_field_a(user_data: &UserData): String {
    user_data.a
}

public fun get_field_b(user_data: &UserData): vector<u8> {
    user_data.b
}

public fun get_field_c(user_data: &UserData): address {
    user_data.c
}

public fun get_field_d(user_data: &UserData): u128 {
    user_data.d
}

// ===================== Test Functions =====================

#[test_only]
public fun test_init(ctx: &mut TxContext) {
    init(MCMS_USER {}, ctx)
}

#[test_only]
public fun test_create_user_data(
    ctx: &mut TxContext,
): (UserData, OwnerCap) {
    let (ownable_state, owner_cap) = ownable::new(ctx);

    (UserData {
        id: object::new(ctx),
        invocations: 0,
        a: string::utf8(b""),
        b: vector[],
        c: @0x0,
        d: 0,
        ownable_state,
    }, owner_cap)
}
