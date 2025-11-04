#[test_only]
module mcms::bcs_stream_test;

use mcms::bcs_stream as bs;
use std::string;
use sui::address;

const MOCK_ADDRESS_1: address = @0x1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b;

#[test]
public fun test_assert_is_consumed() {
    let s = bs::new(vector[]);

    bs::assert_is_consumed(&s);
}

#[test]
#[expected_failure(abort_code = bs::ENotConsumed)]
public fun test_assert_is_consumed_not_consumed() {
    let s = bs::new(vector[1, 2, 3]);

    bs::assert_is_consumed(&s);
}

#[test]
public fun test_deserialize_bool() {
    let mut s = bs::new(vector[1, 2, 3]);
    let b = bs::deserialize_bool(&mut s);

    assert!(b);
    assert!(bs::get_cur(&s) == 1);
}

#[test]
#[expected_failure(abort_code = bs::EMalformedData)]
public fun test_deserialize_bool_malformed() {
    let mut s = bs::new(vector[2, 2, 3]);
    bs::deserialize_bool(&mut s);
}

#[test]
public fun test_deserialize_address() {
    let bytes = address::to_bytes(MOCK_ADDRESS_1);
    let mut s = bs::new(bytes);

    let addr = bs::deserialize_address(&mut s);
    assert!(addr == MOCK_ADDRESS_1);
    assert!(bs::get_cur(&s) == 32);
}

#[test]
#[expected_failure(abort_code = bs::EOutOfBytes)]
public fun test_deserialize_address_out_of_bytes() {
    let mut s = bs::new(vector[1, 2, 3]);
    bs::deserialize_address(&mut s);
}

#[test]
public fun test_deserialize_u8() {
    let mut s = bs::new(vector[3, 2, 1]);
    let u8_val = bs::deserialize_u8(&mut s);

    assert!(u8_val == 3);
    assert!(bs::get_cur(&s) == 1);
}

#[test]
public fun test_deserialize_u16() {
    let mut s = bs::new(vector[3, 2, 1]);
    let u16_val = bs::deserialize_u16(&mut s);

    assert!(u16_val == 515); // 3 + 2 << 8
    assert!(bs::get_cur(&s) == 2);
}

#[test]
public fun test_deserialize_u32() {
    let mut s = bs::new(vector[3, 2, 1, 1]);
    let u32_val = bs::deserialize_u32(&mut s);

    assert!(u32_val == 16843267); // 3 + 2 << 8 + 1 << 16 + 1 << 24 = 3 + 512 + 65536 + 16777216 = 16843267
    assert!(bs::get_cur(&s) == 4);
}

#[test]
public fun test_deserialize_u64() {
    let mut s = bs::new(vector[3, 0, 0, 0, 0, 0, 0, 1]);
    let u64_val = bs::deserialize_u64(&mut s);

    assert!(u64_val == 72057594037927939); // 3 + 1 << 56 = 3 + 72057594037927936 = 72057594037927939
    assert!(bs::get_cur(&s) == 8);
}

#[test]
public fun test_deserialize_u128() {
    let mut s = bs::new(vector[3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1]);
    let u128_val = bs::deserialize_u128(&mut s);

    assert!(u128_val == 1329227995784915872903807060280344579); // 3 + 1 << 120
    assert!(bs::get_cur(&s) == 16);
}

#[test]
public fun test_deserialize_u256() {
    let mut s = bs::new(vector[
        3,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        1,
    ]);
    let u256_val = bs::deserialize_u256(&mut s);

    assert!(
        u256_val == 452312848583266388373324160190187140051835877600158453279131187530910662659,
    ); // 3 + 1 << 248
    assert!(bs::get_cur(&s) == 32);
}

#[test]
public fun test_deserialize_uleb128() {
    let mut s = bs::new(vector[0x81, 0x3, 0, 0]);

    let u64_val = bs::deserialize_uleb128(&mut s);

    // 0x81 & 0x7F = 0x01 = 1
    // 0x03 << 7 = 384
    // 1 + 384 = 385
    assert!(u64_val == 385);
    assert!(bs::get_cur(&s) == 2);
}

#[test]
public fun test_deserialize_fixed_vector_u8() {
    let mut s = bs::new(x"01020304050607");
    let vec = bs::deserialize_fixed_vector_u8(&mut s, 3);

    assert!(vec == vector[0x01, 0x02, 0x03]);
    assert!(bs::get_cur(&s) == 0);

    let u32_val = bs::deserialize_u32(&mut s);
    assert!(u32_val == 117835012); // 4 + 5 << 8 + 6 << 16 + 7 << 24 = 4 + 1280 + 3932160 + 117440512
    assert!(bs::get_cur(&s) == 4);
}

#[test]
public fun test_deserialize_string() {
    let mut s = bs::new(x"0A3333333333333333333403020101");
    let str_val = bs::deserialize_string(&mut s);

    assert!(str_val == string::utf8(b"3333333334"));
    assert!(bs::get_cur(&s) == 0);

    let u32_val = bs::deserialize_u32(&mut s);
    assert!(u32_val == 16843267); // 3 + 2 << 8 + 1 << 16 + 1 << 24 = 3 + 512 + 65536 + 16777216 = 16843267
    assert!(bs::get_cur(&s) == 4);
}

#[test]
public fun test_deserialize_vector_u8() {
    let mut s = bs::new(x"0A31323333333333333334");
    let vec = bs::deserialize_vector_u8(&mut s);

    assert!(vec == vector[0x31, 0x32, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33, 0x34]);
    assert!(bs::get_cur(&s) == 0);
}

#[test]
public fun test_deserialize_option_u8() {
    let mut s = bs::new(x"0108");
    let opt = bs::deserialize_option!(&mut s, |s| bs::deserialize_u8(s));

    assert!(opt == option::some(8));
    assert!(bs::get_cur(&s) == 2); // 1 byte for presence + 1 byte for u8 value
}

#[test]
public fun test_deserialize_option_none() {
    let mut s = bs::new(x"00");
    let opt = bs::deserialize_option!(&mut s, |s| bs::deserialize_u8(s));

    assert!(opt == option::none());
    assert!(bs::get_cur(&s) == 1); // 1 byte for presence
}

#[test]
public fun test_deserialize_option_string() {
    let mut s = bs::new(x"010A31323333333333333334");
    let opt = bs::deserialize_option!(&mut s, |s| bs::deserialize_string(s));

    assert!(opt == option::some(string::utf8(b"1233333334")));
    assert!(bs::get_cur(&s) == 0); // after string deserialization, the cursor is reset to 0 bc the data is updated
}
