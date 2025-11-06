module mcms::bcs_stream;

use mcms::params;
use std::string::{Self, String};
use sui::bcs;

const EMalformedData: u64 = 1;
const EOutOfBytes: u64 = 2;
const ENotConsumed: u64 = 3;
const EInvalidObjectAddress: u64 = 4;

public struct BCSStream has drop {
    /// Byte buffer containing the serialized data.
    data: vector<u8>,
    /// Cursor indicating the current position in the byte buffer.
    cur: u64,
}

public fun assert_is_consumed(stream: &BCSStream) {
    assert!(stream.cur == stream.data.length(), ENotConsumed);
}

public fun deserialize_uleb128(stream: &mut BCSStream): u64 {
    let mut res = 0;
    let mut shift = 0;

    while (stream.cur < stream.data.length()) {
        let byte = stream.data[stream.cur];
        stream.cur = stream.cur + 1;

        let val = ((byte & 0x7f) as u64);
        if (((val << shift) >> shift) != val) {
            abort EMalformedData
        };
        res = res | (val << shift);

        if ((byte & 0x80) == 0) {
            if (shift > 0 && val == 0) {
                abort EMalformedData
            };
            return res
        };

        shift = shift + 7;
        if (shift > 64) {
            abort EMalformedData
        };
    };

    abort EOutOfBytes
}

public fun deserialize_bool(stream: &mut BCSStream): bool {
    assert!(stream.cur < stream.data.length(), EOutOfBytes);

    let byte = stream.data[stream.cur];
    stream.cur = stream.cur + 1;
    if (byte == 0) {
        return false
    } else if (byte == 1) {
        return true
    };

    abort EMalformedData
}

public fun deserialize_address(stream: &mut BCSStream): address {
    let data = &stream.data;
    let cur = stream.cur;

    assert!(cur + 32 <= data.length(), EOutOfBytes);

    let address_bytes = params::slice(data, cur, 32);
    let mut bcs_instance = bcs::new(address_bytes);
    stream.cur = stream.cur + 32;
    bcs::peel_address(&mut bcs_instance)
}

public fun deserialize_u8(stream: &mut BCSStream): u8 {
    let data = &stream.data;
    let cur = stream.cur;

    assert!(cur < data.length(), EOutOfBytes);

    let res = data[cur];

    stream.cur = cur + 1;
    res
}

public fun deserialize_u16(stream: &mut BCSStream): u16 {
    let data = &stream.data;
    let cur = stream.cur;

    assert!(cur + 2 <= data.length(), EOutOfBytes);
    let res = (data[cur] as u16) | ((data[cur + 1] as u16) << 8);

    stream.cur = stream.cur + 2;
    res
}

public fun deserialize_u32(stream: &mut BCSStream): u32 {
    let data = &stream.data;
    let cur = stream.cur;

    assert!(cur + 4 <= data.length(), EOutOfBytes);
    let res =
        (data[cur] as u32)
            | ((data[cur + 1]  as u32) << 8)
            | ((data[cur + 2] as u32) << 16)
            | ((data[cur + 3]  as u32) << 24);

    stream.cur = stream.cur + 4;
    res
}

public fun deserialize_u64(stream: &mut BCSStream): u64 {
    let data = &stream.data;
    let cur = stream.cur;

    assert!(cur + 8 <= data.length(), EOutOfBytes);
    let res =
        (data[cur] as u64)
            | ((data[cur + 1] as u64) << 8)
            | ((data[cur + 2] as u64) << 16)
            | ((data[cur + 3] as u64) << 24)
            | ((data[cur + 4] as u64) << 32)
            | ((data[cur + 5] as u64) << 40)
            | ((data[cur + 6] as u64) << 48)
            | ((data[cur + 7] as u64) << 56);

    stream.cur = stream.cur + 8;
    res
}

public fun deserialize_u128(stream: &mut BCSStream): u128 {
    let data = &stream.data;
    let cur = stream.cur;

    assert!(cur + 16 <= data.length(), EOutOfBytes);
    let res =
        (data[cur]  as u128)
            | ((data[cur + 1] as u128) << 8)
            | ((data[cur + 2] as u128) << 16)
            | ((data[cur + 3] as u128) << 24)
            | ((data[cur + 4] as u128) << 32)
            | ((data[cur + 5] as u128) << 40)
            | ((data[cur + 6] as u128) << 48)
            | ((data[cur + 7] as u128) << 56)
            | ((data[cur + 8] as u128) << 64)
            | ((data[cur + 9] as u128) << 72)
            | ((data[cur + 10] as u128) << 80)
            | ((data[cur + 11] as u128) << 88)
            | ((data[cur + 12] as u128) << 96)
            | ((data[cur + 13] as u128) << 104)
            | ((data[cur + 14] as u128) << 112)
            | ((data[cur + 15] as u128) << 120);

    stream.cur = stream.cur + 16;
    res
}

public fun deserialize_u256(stream: &mut BCSStream): u256 {
    let data = &stream.data;
    let cur = stream.cur;

    assert!(cur + 32 <= data.length(), EOutOfBytes);
    let res =
        (data[cur] as u256)
            | ((data[cur + 1] as u256) << 8)
            | ((data[cur + 2] as u256) << 16)
            | ((data[cur + 3] as u256) << 24)
            | ((data[cur + 4] as u256) << 32)
            | ((data[cur + 5] as u256) << 40)
            | ((data[cur + 6] as u256) << 48)
            | ((data[cur + 7] as u256) << 56)
            | ((data[cur + 8] as u256) << 64)
            | ((data[cur + 9] as u256) << 72)
            | ((data[cur + 10] as u256) << 80)
            | ((data[cur + 11] as u256) << 88)
            | ((data[cur + 12] as u256) << 96)
            | ((data[cur + 13] as u256) << 104)
            | ((data[cur + 14] as u256) << 112)
            | ((data[cur + 15] as u256) << 120)
            | ((data[cur + 16] as u256) << 128)
            | ((data[cur + 17] as u256) << 136)
            | ((data[cur + 18] as u256) << 144)
            | ((data[cur + 19] as u256) << 152)
            | ((data[cur + 20] as u256) << 160)
            | ((data[cur + 21] as u256) << 168)
            | ((data[cur + 22] as u256) << 176)
            | ((data[cur + 23] as u256) << 184)
            | ((data[cur + 24] as u256) << 192)
            | ((data[cur + 25] as u256) << 200)
            | ((data[cur + 26] as u256) << 208)
            | ((data[cur + 27] as u256) << 216)
            | ((data[cur + 28] as u256) << 224)
            | ((data[cur + 29] as u256) << 232)
            | ((data[cur + 30] as u256) << 240)
            | ((data[cur + 31] as u256) << 248);

    stream.cur = stream.cur + 32;
    res
}

/// Deserializes a `u256` value from the stream.
public fun deserialize_u256_entry(data: vector<u8>, cursor: u64) {
    let mut stream = BCSStream { data: data, cur: cursor };
    deserialize_u256(&mut stream);
}

public fun new(data: vector<u8>): BCSStream {
    BCSStream { data, cur: 0 }
}

public fun deserialize_fixed_vector_u8(stream: &mut BCSStream, len: u64): vector<u8> {
    let data = &mut stream.data;
    let cur = stream.cur;

    assert!(cur + len <= data.length(), EOutOfBytes);

    let mut res = trim(data, cur);
    stream.data = trim(&mut res, len);
    stream.cur = 0;

    res
}

public fun deserialize_string(stream: &mut BCSStream): String {
    let len = deserialize_uleb128(stream);
    let data = &mut stream.data;
    let cur = stream.cur;

    assert!(cur + len <= data.length(), EOutOfBytes);

    let mut res = trim(data, cur);
    stream.data = trim(&mut res, len);
    stream.cur = 0;

    string::utf8(res)
}

/// First, reads the length of the vector, which is in uleb128 format.
/// After determining the length, it then reads the contents of the vector.
/// The `elem_deserializer` lambda expression is used sequentially to deserialize each element of the vector.
public macro fun deserialize_vector<$E>(
    $stream: &mut BCSStream,
    $elem_deserializer: |&mut BCSStream| -> $E,
): vector<$E> {
    let len = deserialize_uleb128($stream);
    let mut v = vector::empty();

    let mut i = 0;
    while (i < len) {
        v.push_back($elem_deserializer($stream));
        i = i + 1;
    };

    v
}

public fun deserialize_vector_u8(stream: &mut BCSStream): vector<u8> {
    let len = deserialize_uleb128(stream);
    let data = &mut stream.data;
    let cur = stream.cur;

    assert!(cur + len <= data.length(), EOutOfBytes);

    let mut res = trim(data, cur);
    stream.data = trim(&mut res, len);
    stream.cur = 0;

    res
}

public fun validate_obj_addr(addr: address, stream: &mut BCSStream) {
    let deserialized_address = deserialize_address(stream);
    assert!(deserialized_address == addr, EInvalidObjectAddress);
}

public fun validate_obj_addrs(addrs: vector<address>, stream: &mut BCSStream) {
    let mut i = 0;
    while (i < addrs.length()) {
        validate_obj_addr(addrs[i], stream);
        i = i + 1;
    }
}

/// Deserializes `Option` from the stream.
/// First, reads a single byte representing the presence (0x01) or absence (0x00) of data.
/// After determining the presence of data, it then reads the actual data if present.
/// The `f` lambda expression is used to deserialize the element contained within the `Option`.
public macro fun deserialize_option<$E>(
    $stream: &mut BCSStream,
    $f: |&mut BCSStream| -> $E,
): Option<$E> {
    let is_data = deserialize_bool($stream);
    if (is_data) {
        option::some($f($stream))
    } else {
        option::none()
    }
}

// this is the equivalent of vector::trim in Aptos Move
fun trim<T: copy>(vec: &mut vector<T>, new_len: u64): vector<T> {
    let mut removed = vector::empty<T>();
    let orig_len = vec.length();

    // If new_len is greater than or equal to the current length, nothing to remove.
    if (new_len >= orig_len) {
        return removed
    };

    // Remove elements from the back until the vector's length equals new_len.
    while (vec.length() > new_len) {
        let elem = vec.pop_back();
        removed.push_back(elem);
    };

    // The elements in `removed` are in reverse order relative to their original order.
    // Reverse the vector to restore the original order.
    vector::reverse(&mut removed);
    removed
}

#[test_only]
public fun get_cur(stream: &BCSStream): u64 {
    stream.cur
}
