module mcms_test::mcms_user;

use mcms::bcs_stream;
use mcms::mcms_deployer::{Self, DeployerState};
use mcms::mcms_registry::{Self, ExecutingCallbackParams, Registry};
use std::string::{Self, String};
use sui::package::UpgradeCap;

const EInvalidAdminCap: u64 = 1;
const EUnknownFunction: u64 = 2;

public struct UserData has key, store {
    id: UID,
    invocations: u8,
    a: String,
    b: vector<u8>,
    c: address,
    d: u128,
    owner_cap: ID,
}

public struct OwnerCap has key, store {
    id: UID,
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

public fun function_two(user_data: &mut UserData, owner_cap: &OwnerCap, arg1: address, arg2: u128) {
    assert_valid_owner_cap(user_data, owner_cap);

    user_data.invocations = user_data.invocations + 1;
    user_data.c = arg1;
    user_data.d = arg2;
}

public struct MCMS_USER has drop {}

fun init(_witness: MCMS_USER, ctx: &mut TxContext) {
    let owner_cap = OwnerCap {
        id: object::new(ctx),
    };

    transfer::share_object(UserData {
        id: object::new(ctx),
        invocations: 0,
        a: string::utf8(b""),
        b: vector[],
        c: @0x0,
        d: 0,
        owner_cap: object::id(&owner_cap),
    });

    transfer::transfer(owner_cap, ctx.sender());
}

public fun initialize(
    owner_cap: OwnerCap,
    upgrade_cap: UpgradeCap,
    user_data: &UserData,
    registry: &mut Registry,
    state: &mut DeployerState,
    ctx: &mut TxContext,
) {
    assert_valid_owner_cap(user_data, &owner_cap);

    // Transfer owner_cap to MCMS
    mcms_registry::register_entrypoint(
        registry,
        SampleMcmsCallback {},
        option::some(owner_cap),
        ctx,
    );

    // Transfer upgrade permissions to MCMS
    mcms_deployer::register_upgrade_cap(
        state,
        registry,
        upgrade_cap,
        ctx,
    );
}

fun assert_valid_owner_cap(user_data: &UserData, owner_cap: &OwnerCap) {
    assert!(user_data.owner_cap == object::id(owner_cap), EInvalidAdminCap);
}

public struct SampleMcmsCallback has drop {}

public fun mcms_entrypoint(
    registry: &mut Registry,
    user_data: &mut UserData,
    params: ExecutingCallbackParams, // hot potato
    _ctx: &mut TxContext,
) {
    let (owner_cap, function, data) = mcms_registry::get_callback_params<
        SampleMcmsCallback,
        OwnerCap,
    >(
        registry,
        SampleMcmsCallback {},
        params,
    );

    let function_bytes = *function.as_bytes();
    let mut stream = bcs_stream::new(data);

    if (function_bytes == b"function_one") {
        let arg1 = bcs_stream::deserialize_string(&mut stream);
        let arg2 = bcs_stream::deserialize_vector_u8(&mut stream);
        bcs_stream::assert_is_consumed(&stream);
        function_one(user_data, owner_cap, arg1, arg2);
    } else if (function_bytes == b"function_two") {
        let arg1 = bcs_stream::deserialize_address(&mut stream);
        let arg2 = bcs_stream::deserialize_u128(&mut stream);
        bcs_stream::assert_is_consumed(&stream);
        function_two(user_data, owner_cap, arg1, arg2);
    } else {
        abort EUnknownFunction
    };
}

public fun get_owner_cap(user_data: &UserData): ID {
    user_data.owner_cap
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
public fun test_create_user_data(ctx: &mut TxContext, owner_cap: ID): UserData {
    UserData {
        id: object::new(ctx),
        invocations: 0,
        a: string::utf8(b""),
        b: vector[],
        c: @0x0,
        d: 0,
        owner_cap,
    }
}
