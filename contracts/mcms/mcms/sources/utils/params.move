module mcms::params;

use std::string::{Self, String};
use std::type_name::TypeName;
use sui::address;
use sui::bcs;
use sui::hex;

const ECmpVectorsDiffLen: u64 = 1;
const EOutOfBytes: u64 = 2;
const ENonModuleType: u64 = 3;
const EInputTooLargeForNumBytes: u64 = 4;

const ASCII_COLON: u8 = 58;
const ASCII_LESS_THAN: u8 = 60; // '<' character for generics

public fun encode_uint<T: drop>(input: T, num_bytes: u64): vector<u8> {
    let mut bcs_bytes = bcs::to_bytes(&input);

    let len = bcs_bytes.length();
    assert!(len <= num_bytes, EInputTooLargeForNumBytes);

    if (len < num_bytes) {
        let bytes_to_pad = num_bytes - len;
        let mut i = 0;
        while (i < bytes_to_pad) {
            bcs_bytes.push_back(0);
            i = i + 1;
        };
    };

    // little endian to big endian
    bcs_bytes.reverse();

    bcs_bytes
}

public fun right_pad_vec(v: &mut vector<u8>, num_bytes: u64) {
    let len = v.length();
    if (len < num_bytes) {
        let bytes_to_pad = num_bytes - len;
        let mut i = 0;
        while (i < bytes_to_pad) {
            v.push_back(0);
            i = i + 1;
        };
    };
}

/// compares two vectors of equal length, returns true if a > b, false otherwise.
public fun vector_u8_gt(a: &vector<u8>, b: &vector<u8>): bool {
    let len = a.length();
    assert!(len == b.length(), ECmpVectorsDiffLen);

    if (len == 0) {
        return false
    };

    // compare each byte until not equal
    let mut i = 0;
    while (i < len) {
        let byte_a = a[i];
        let byte_b = b[i];
        if (byte_a > byte_b) {
            return true
        } else if (byte_a < byte_b) {
            return false
        };
        i = i + 1;
    };

    // vectors are equal, a == b
    false
}

public fun get_account_address_and_module_name(proof_type: TypeName): (address, String) {
    let account_address_bytes = hex::decode(proof_type.address_string().into_bytes());
    let account_address = address::from_bytes(account_address_bytes);
    let module_name = string::from_ascii(proof_type.module_string());
    (account_address, module_name)
}

/// Get the struct name from a TypeName.
/// e.g. "0x1::option::Option<u64>" -> "Option"
public fun get_struct_name(type_name: &TypeName): vector<u8> {
    assert!(!type_name.is_primitive(), ENonModuleType);

    let str_bytes = type_name.as_string().as_bytes();

    // Skip address (2 chars per byte) and "::"
    let mut i = address::length() * 2 + 2;

    // Skip module name - find the next "::"
    let colon = ASCII_COLON;
    while (i < str_bytes.length()) {
        if (str_bytes[i] == colon && i + 1 < str_bytes.length() && str_bytes[i + 1] == colon) {
            i = i + 2; // skip "::"
            break
        };
        i = i + 1;
    };

    // Now collect the struct name until we hit '<' (generics) or end
    let mut struct_name = vector[];
    let less_than = ASCII_LESS_THAN;
    while (i < str_bytes.length()) {
        let char = str_bytes[i];
        if (char == less_than) {
            break
        };
        struct_name.push_back(char);
        i = i + 1;
    };

    struct_name
}

public fun slice<T: copy>(v: &vector<T>, start: u64, len: u64): vector<T> {
    let v_len = v.length();
    assert!(start + len <= v_len, EOutOfBytes);

    let mut result = vector[];
    let mut i = start;
    while (i < start + len) {
        result.push_back(v[i]);
        i = i + 1;
    };
    result
}
