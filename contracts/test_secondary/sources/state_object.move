module test_secondary::state_object;

use std::ascii;
use std::string;
use std::type_name;
use sui::address;
use sui::derived_object;

const EModuleAlreadyExists: u64 = 1;
const EModuleDoesNotExist: u64 = 2;
const EInvalidFunction: u64 = 3;
const EInvalidOwnerCap: u64 = 4;
const EPackageIdNotFound: u64 = 5;

public struct CCIPObject has key {
    id: UID,
}

public struct CCIPObjectRef has key, store {
    id: UID,
    package_ids: vector<address>,
}

public struct CCIPObjectRefPointer has key, store {
    id: UID,
    ccip_object_id: address,
}

public struct STATE_OBJECT has drop {}

fun init(_witness: STATE_OBJECT, ctx: &mut TxContext) {
    let mut ccip_object = CCIPObject { id: object::new(ctx) };

    let mut ref = CCIPObjectRef {
        id: derived_object::claim(&mut ccip_object.id, b"CCIPObjectRef"),
        package_ids: vector[],
    };

    let pointer = CCIPObjectRefPointer {
        id: object::new(ctx),
        ccip_object_id: object::id_address(&ccip_object),
    };

    let tn = type_name::get_with_original_ids<STATE_OBJECT>();
    let package_bytes = ascii::into_bytes(tn.get_address());
    let package_id = address::from_ascii_bytes(&package_bytes);
    ref.package_ids.push_back(package_id);

    transfer::share_object(ref);
    transfer::share_object(ccip_object);

    transfer::transfer(pointer, package_id);
}

public fun get_package_ids(pointer: &CCIPObjectRef): vector<address> {
    pointer.package_ids
}
