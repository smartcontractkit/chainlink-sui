module mcms::bcs_stream {
    use std::string::{Self, String};
    use sui::bcs;

    use mcms::params;

    const E_MALFORMED_DATA: u64 = 1;
    const E_OUT_OF_BYTES: u64 = 2;
    const E_NOT_CONSUMED: u64 = 3;

    public struct BCSStream has drop {
        /// Byte buffer containing the serialized data.
        data: vector<u8>,
        /// Cursor indicating the current position in the byte buffer.
        cur: u64,
    }

    public fun assert_is_consumed(stream: &BCSStream) {
        assert!(stream.cur == stream.data.length(), E_NOT_CONSUMED);
    }

    public fun deserialize_uleb128(stream: &mut BCSStream): u64 {
        let mut res = 0;
        let mut shift = 0;

        while (stream.cur < stream.data.length()) {
            let byte = stream.data[stream.cur];
            stream.cur = stream.cur + 1;

            let val = ((byte & 0x7f) as u64);
            if (((val << shift) >> shift) != val) {
                abort E_MALFORMED_DATA
            };
            res = res | (val << shift);

            if ((byte & 0x80) == 0) {
                if (shift > 0 && val == 0) {
                    abort E_MALFORMED_DATA
                };
                return res
            };

            shift = shift + 7;
            if (shift > 64) {
                abort E_MALFORMED_DATA
            };
        };

        abort E_OUT_OF_BYTES
    }

    public fun deserialize_bool(stream: &mut BCSStream): bool {
        assert!(stream.cur < stream.data.length(), E_OUT_OF_BYTES);

        let byte = stream.data[stream.cur];
        stream.cur = stream.cur + 1;
        if (byte == 0) {
            return false
        } else if (byte == 1) {
            return true
        } else {
            abort E_MALFORMED_DATA
        }
    }

    public fun deserialize_address(stream: &mut BCSStream): address {
        let data = &stream.data;
        let cur = stream.cur;

        assert!(cur + 32 <= data.length(), E_OUT_OF_BYTES);

        let address_bytes = params::slice(data, cur, 32);
        let mut bcs_instance = bcs::new(address_bytes);
        stream.cur = stream.cur + 32;
        bcs::peel_address(&mut bcs_instance)
    }

    public fun deserialize_u8(stream: &mut BCSStream): u8 {
        let data = &stream.data;
        let cur = stream.cur;

        assert!(cur < data.length(), E_OUT_OF_BYTES);

        let res = data[cur];

        stream.cur = cur + 1;
        res
    }

    public fun deserialize_u16(stream: &mut BCSStream): u16 {
        let data = &stream.data;
        let cur = stream.cur;

        assert!(cur + 2 <= data.length(), E_OUT_OF_BYTES);
        let res = (data[cur] as u16) | ((data[cur + 1] as u16) << 8);

        stream.cur = stream.cur + 2;
        res
    }

    public fun deserialize_u32(stream: &mut BCSStream): u32 {
        let data = &stream.data;
        let cur = stream.cur;

        assert!(cur + 4 <= data.length(), E_OUT_OF_BYTES);
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

        assert!(cur + 8 <= data.length(), E_OUT_OF_BYTES);
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

        assert!(cur + 16 <= data.length(), E_OUT_OF_BYTES);
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

        assert!(cur + 32 <= data.length(), E_OUT_OF_BYTES);
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
    public entry fun deserialize_u256_entry(data: vector<u8>, cursor: u64) {
        let mut stream = BCSStream { data: data, cur: cursor };
        deserialize_u256(&mut stream);
    }

    public fun new(data: vector<u8>): BCSStream {
        BCSStream { data, cur: 0 }
    }

    public fun deserialize_fixed_vector_u8(stream: &mut BCSStream, len: u64): vector<u8> {
        let data = &mut stream.data;
        let cur = stream.cur;

        assert!(cur + len <= data.length(), E_OUT_OF_BYTES);

        let mut res = trim(data, cur);
        stream.data = trim(&mut res, len);
        stream.cur = 0;

        res
    }

    public fun deserialize_string(stream: &mut BCSStream): String {
        let len = deserialize_uleb128(stream);
        let data = &mut stream.data;
        let cur = stream.cur;

        assert!(cur + len <= data.length(), E_OUT_OF_BYTES);

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

        assert!(cur + len <= data.length(), E_OUT_OF_BYTES);

        let mut res = trim(data, cur);
        stream.data = trim(&mut res, len);
        stream.cur = 0;

        res
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

    // TODO: check if there is a more efficient way to implement this function
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
}

#[test_only]
module mcms::bcs_stream_test {
    use mcms::bcs_stream as bs;
    use std::string;
    use sui::address;

    const MOCK_ADDRESS_1: address =
        @0x1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b;

    #[test]
    public fun test_assert_is_consumed() {
        let s = bs::new(vector[]);

        bs::assert_is_consumed(&s);
    }

    #[test]
    #[expected_failure(abort_code = bs::E_NOT_CONSUMED)]
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
    #[expected_failure(abort_code = bs::E_MALFORMED_DATA)]
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
    #[expected_failure(abort_code = bs::E_OUT_OF_BYTES)]
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
}
