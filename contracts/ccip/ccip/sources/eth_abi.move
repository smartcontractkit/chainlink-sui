// module to do the equivalent packing as ethereum's abi.encode and abi.encodePacked
module ccip::eth_abi {
    use sui::bcs;
    use sui::address;

    const E_OUT_OF_BYTES: u64 = 1;
    // E_INVALID_BYTES32 is not used. keep it for now to match Aptos error code
    // const E_INVALID_BYTES32: u64 = 2;
    const E_INVALID_ADDRESS: u64 = 3;
    const E_INVALID_BOOL: u64 = 4;
    const E_INVALID_SELECTOR: u64 = 5;
    const E_INVALID_U256_LENGTH: u64 = 6;
    const E_INVALID_LENGTH: u64 = 7;
    const ENCODED_BOOL_FALSE: vector<u8> = vector[
        0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
        0, 0, 0, 0
    ];
    const ENCODED_BOOL_TRUE: vector<u8> = vector[
        0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
        0, 0, 0, 1
    ];

    public fun encode_object_id(out: &mut vector<u8>, value: object::ID) {
        encode_address(out, object::id_to_address(&value))
    }

    public fun encode_address(out: &mut vector<u8>, value: address) {
        vector::append(out, bcs::to_bytes(&value))
    }

    public fun encode_u8(out: &mut vector<u8>, value: u8) {
        encode_u256(out, value as u256);
    }

    public fun encode_u32(out: &mut vector<u8>, value: u32) {
        encode_u256(out, value as u256)
    }

    public fun encode_u64(out: &mut vector<u8>, value: u64) {
        encode_u256(out, value as u256)
    }

    public fun encode_u256(out: &mut vector<u8>, value: u256) {
        let mut value_bytes = bcs::to_bytes(&value);
        // little endian to big endian
        vector::reverse(&mut value_bytes);
        vector::append(out, value_bytes)
    }

    public fun encode_bool(out: &mut vector<u8>, value: bool) {
        vector::append(out, if (value) ENCODED_BOOL_TRUE
        else ENCODED_BOOL_FALSE)
    }

    public fun encode_bytes32(
        out: &mut vector<u8>, value: vector<u8>
    ) {
        assert!(value.length() <= 32, E_INVALID_LENGTH);
        let padding_len = 32 - value.length();
        let mut i = 0;
        while (i < padding_len) {
            out.push_back(0);
            i = i + 1;
        };
        vector::append(out, value)
    }

    public fun encode_bytes(out: &mut vector<u8>, value: vector<u8>) {
        encode_u256(out, (value.length() as u256));

        vector::append(out, value);
        let padding_len = 32 - (value.length() % 32);
        let mut i = 0;
        while (i < padding_len) {
            out.push_back(0);
            i = i + 1;
        }
    }

    public fun encode_selector(out: &mut vector<u8>, value: vector<u8>) {
        assert!(value.length() == 4, E_INVALID_SELECTOR);
        vector::append(out, value);
    }

    // TODO: not used onchain. verify if used offchain
    public fun encode_packed_address(
        out: &mut vector<u8>, value: address
    ) {
        vector::append(out, bcs::to_bytes(&value))
    }

    public fun encode_packed_bytes(
        out: &mut vector<u8>, value: vector<u8>
    ) {
        vector::append(out, value)
    }

    public fun encode_packed_bytes32(
        out: &mut vector<u8>, value: vector<u8>
    ) {
        assert!(value.length() <= 32, E_INVALID_LENGTH);
        vector::append(out, value)
    }

    public fun encode_packed_u8(out: &mut vector<u8>, value: u8) {
        out.push_back(value);
    }

    public fun encode_packed_u32(out: &mut vector<u8>, value: u32) {
        let mut value_bytes = bcs::to_bytes(&value);
        // little endian to big endian
        vector::reverse(&mut value_bytes);
        vector::append(out, value_bytes)
    }

    public fun encode_packed_u64(out: &mut vector<u8>, value: u64) {
        let mut value_bytes = bcs::to_bytes(&value);
        // little endian to big endian
        vector::reverse(&mut value_bytes);
        vector::append(out, value_bytes)
    }

    public fun encode_packed_u256(out: &mut vector<u8>, value: u256) {
        let mut value_bytes = bcs::to_bytes(&value);
        // little endian to big endian
        vector::reverse(&mut value_bytes);
        vector::append(out, value_bytes)
    }

    // ABIStream won't be published. no need to add a key
    public struct ABIStream has drop {
        data: vector<u8>,
        cur: u64
    }

    #[test_only]
    public fun get_cur(stream: &ABIStream): u64 {
        stream.cur
    }

    public fun new_stream(data: vector<u8>): ABIStream {
        ABIStream { data, cur: 0 }
    }

    public fun decode_address(stream: &mut ABIStream): address {
        let data = &stream.data;
        let cur = stream.cur;

        assert!(
            cur + 32 <= data.length(),
            E_OUT_OF_BYTES
        );

        // Verify first 12 bytes are zero
        // This is to decode Ethereum address not Sui address
        let mut i = 0;
        let mut value_bytes = vector[];
        while (i < 12) {
            assert!(
                data[cur + i] == 0,
                E_INVALID_ADDRESS
            );
            value_bytes.push_back(0);
            i = i + 1;
        };

        // add the remaining 20 bytes
        while (i < 32) {
            value_bytes.push_back(data[cur + i]);
            i = i + 1;
        };
        stream.cur = cur + 32;

        address::from_bytes(value_bytes)
    }

    public fun decode_u256_value(mut value_bytes: vector<u8>): u256 {
        assert!(
            value_bytes.length() == 32,
            E_INVALID_U256_LENGTH
        );
        vector::reverse(&mut value_bytes);

        // Deserialize to u256
        bcs::peel_u256(&mut bcs::new(value_bytes))
    }

    public fun decode_u256(stream: &mut ABIStream): u256 {
        let data = &stream.data;
        let cur = stream.cur;

        assert!(
            cur + 32 <= data.length(),
            E_OUT_OF_BYTES
        );

        let mut value_bytes = slice(data, cur, 32);
        // Convert from big endian to little endian
        vector::reverse(&mut value_bytes);

        stream.cur = cur + 32;
        bcs::peel_u256(&mut bcs::new(value_bytes))
    }

    /// Returns a new vector containing `len` elements from `vec`
    /// starting at index `start`. Panics if `start + len` exceeds the vector length.
    fun slice<T: copy>(vec: &vector<T>, start: u64, len: u64): vector<T> {
        let vec_len = vec.length();
        // Ensure we have enough elements for the slice.
        assert!(start + len <= vec_len, E_OUT_OF_BYTES);
        let mut new_vec = vector::empty<T>();
        let mut i = start;
        while (i < start + len) {
            // Copy each element from the original vector into the new vector.
            new_vec.push_back(vec[i]);
            i = i + 1;
        };
        new_vec
    }

    public fun decode_u8(stream: &mut ABIStream): u8 {
        (decode_u256(stream) as u8)
    }

    public fun decode_u32(stream: &mut ABIStream): u32 {
        (decode_u256(stream) as u32)
    }

    public fun decode_u64(stream: &mut ABIStream): u64 {
        (decode_u256(stream) as u64)
    }

    public fun decode_bool(stream: &mut ABIStream): bool {
        let data = &stream.data;
        let cur = stream.cur;

        assert!(
            cur + 32 <= data.length(),
            E_OUT_OF_BYTES
        );

        let value = slice(data, cur, 32);
        stream.cur = cur + 32;

        if (value == ENCODED_BOOL_FALSE) { false }
        else if (value == ENCODED_BOOL_TRUE) { true }
        else {
            abort E_INVALID_BOOL
        }
    }

    public macro fun decode_vector<$E>(
        $stream: &mut ABIStream, $f: |&mut ABIStream| -> $E
    ): vector<$E> {
        let len = decode_u256($stream);
        let mut v = vector::empty();
        let mut i = 0;

        while (i < len) {
            v.push_back($f($stream));
            i = i + 1;
        };

        v
    }

    public fun decode_bytes(stream: &mut ABIStream): vector<u8> {
        // First read length as u256
        let length = (decode_u256(stream) as u64);

        let padding_len = if (length % 32 == 0) { 0 }
        else {
            32 - (length % 32)
        };

        let data = &stream.data;
        let cur = stream.cur;

        assert!(
            cur + length + padding_len <= data.length(),
            E_OUT_OF_BYTES
        );

        let bytes = slice(data, cur, length);

        stream.cur = cur + length + padding_len;

        bytes
    }

    public fun decode_bytes32(stream: &mut ABIStream): vector<u8> {
        let data = &stream.data;
        let cur = stream.cur;

        assert!(
            cur + 32 <= data.length(),
            E_OUT_OF_BYTES
        );

        let bytes = slice(data, cur, 32);
        stream.cur = cur + 32;
        bytes
    }
}

#[test_only]
module ccip::eth_abi_test {
    use ccip::eth_abi;

    #[test]
    public fun encode_eth_bool() {
        let mut v = vector[];
        let encoded_bool_true: vector<u8> = vector[
            0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1
        ];

        eth_abi::encode_bool(&mut v, true);

        assert!(v == encoded_bool_true);
    }

    #[test]
    public fun encode_eth_u256() {
        let mut v = vector[];
        // 256 is 0x100 in hex
        let value: u256 = 256;

        eth_abi::encode_u256(&mut v, value);

        let encoded_uint256: vector<u8> = vector[
            0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0
        ];
        assert!(v == encoded_uint256);
    }

    #[test]
    public fun encode_eth_u64() {
        let mut v = vector[];
        // 99 is 0x63 in hex
        let value: u64 = 99;

        eth_abi::encode_u64(&mut v, value);

        let encoded_u64: vector<u8> = vector[
            0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x63
        ];
        assert!(v == encoded_u64);
    }

    #[test]
    public fun encode_eth_u32() {
        let mut v = vector[];
        let value: u32 = 4294967295; // 0xFFFFFFFF max u32

        eth_abi::encode_u32(&mut v, value);

        let encoded_u32: vector<u8> = vector[
            0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xFF, 0xFF, 0xFF, 0xFF
        ];
        assert!(v == encoded_u32);
    }

    #[test]
    public fun encode_eth_u8() {
        let mut v = vector[];
        let value: u8 = 0;

        eth_abi::encode_u8(&mut v, value);

        let encode_u8: vector<u8> = vector[
            0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0
        ];
        assert!(v == encode_u8);
    }

    #[test]
    public fun encode_eth_address() {
        let mut v = vector[];
        let value: address = @0xADa80b6ae7F00960C3020b5E97AAACCc3a4674f9;

        eth_abi::encode_address(&mut v, value);

        let encode_address: vector<u8> =  vector[
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x00, 0x00, 0xAD, 0xa8, 0x0b, 0x6a,
            0xe7, 0xF0, 0x09, 0x60, 0xC3, 0x02, 0x0b, 0x5E,
            0x97, 0xAA, 0xAC, 0xCc, 0x3a, 0x46, 0x74, 0xf9
        ];
        assert!(v == encode_address);
    }

    #[test]
    public fun encode_bytes32() {
        let mut v = vector[];
        let value: vector<u8> = vector[0x01, 0x02, 0x03, 0x04];

        eth_abi::encode_bytes32(&mut v, value);
        let encode_bytes32: vector<u8> = vector[
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x00, 0x00, 0x01, 0x02, 0x03, 0x04
        ];
        assert!(v == encode_bytes32);
    }

    #[test]
    #[expected_failure(abort_code = eth_abi::E_INVALID_LENGTH)]
    public fun encode_bytes32_failed() {
        let mut v = vector[];
        let value: vector<u8> = vector[
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x00, 0x00, 0x01, 0x02, 0x03, 0x04
        ];

        eth_abi::encode_bytes32(&mut v, value);
    }

    #[test]
    public fun encode_bytes() {
        let mut v = vector[];
        let value: vector<u8> = vector[0x01, 0x02, 0x03, 0x04, 0x05];
        eth_abi::encode_bytes(&mut v, value);

        let encoded_bytes: vector<u8> = vector[
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05,
            0x01, 0x02, 0x03, 0x04, 0x05, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00
        ];

        assert!(v == encoded_bytes);
    }

    #[test]
    public fun encode_selector() {
        let mut v = vector[];
        let value: vector<u8> = vector[0x01, 0x02, 0x03, 0x04];
        eth_abi::encode_selector(&mut v, value);

        let encoded_selector: vector<u8> = vector[
            0x01, 0x02, 0x03, 0x04
        ];

        assert!(v == encoded_selector);
    }

    #[test]
    #[expected_failure(abort_code = eth_abi::E_INVALID_SELECTOR)]
    public fun encode_selector_failed() {
        let mut v = vector[];
        let value: vector<u8> = vector[0x01, 0x02, 0x03, 0x04, 0x05];
        eth_abi::encode_selector(&mut v, value);
    }

    #[test]
    public fun encode_packed_bytes32() {
        let mut v = vector[];
        let value: vector<u8> = vector[0x01, 0x02, 0x03, 0x04, 0x05];

        eth_abi::encode_packed_bytes32(&mut v, value);
        let encoded_packed_bytes32: vector<u8> = vector[
            0x01, 0x02, 0x03, 0x04, 0x05
        ];
        assert!(v == encoded_packed_bytes32);
    }

    #[test]
    #[expected_failure(abort_code = eth_abi::E_INVALID_LENGTH)]
    public fun encode_packed_bytes32_failed() {
        let mut v = vector[];
        let value: vector<u8> = vector[
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05,
            0x01, 0x02, 0x03, 0x04, 0x05, 0x00, 0x00, 0x00
        ];

        eth_abi::encode_packed_bytes32(&mut v, value);
    }

    #[test]
    public fun encode_packed_bytes() {
        let mut v = vector[];
        let value: vector<u8> = vector[0x01, 0x02, 0x03, 0x04, 0x05];
        eth_abi::encode_packed_bytes(&mut v, value);

        let encoded_bytes: vector<u8> = vector[0x01, 0x02, 0x03, 0x04, 0x05];

        assert!(v == encoded_bytes);
    }

    #[test]
    public fun encode_packed_u32() {
        let mut v = vector[];
        let value: u32 = 4294967295; // 0xFFFFFFFF max u32

        eth_abi::encode_packed_u32(&mut v, value);

        let encoded_packed_u32: vector<u8> = vector[
            0xFF, 0xFF, 0xFF, 0xFF
        ];
        assert!(v == encoded_packed_u32);
    }

    #[test]
    public fun encode_packed_u64() {
        let mut v = vector[];
        let value: u64 = 99;

        eth_abi::encode_packed_u64(&mut v, value);

        let encoded_packed_u64: vector<u8> = vector[
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x63
        ];
        assert!(v == encoded_packed_u64);
    }

    #[test]
    public fun encode_packed_u256() {
        let mut v = vector[];
        let value: u256 = 256;

        eth_abi::encode_packed_u256(&mut v, value);

        let encoded_packed_u256: vector<u8> = vector[
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00
        ];
        assert!(v == encoded_packed_u256);
    }

    #[test]
    public fun decode_address() {
        let addr: vector<u8> =  vector[
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x00, 0x00, 0xAD, 0xa8, 0x0b, 0x6a,
            0xe7, 0xF0, 0x09, 0x60, 0xC3, 0x02, 0x0b, 0x5E,
            0x97, 0xAA, 0xAC, 0xCc, 0x3a, 0x46, 0x74, 0xf9
        ];

        let mut stream: eth_abi::ABIStream = eth_abi::new_stream(addr);

        let decoded_addr: address = eth_abi::decode_address(&mut stream);

        assert!(decoded_addr == @0xADa80b6ae7F00960C3020b5E97AAACCc3a4674f9);
        assert!(eth_abi::get_cur(&stream) == 32);
    }

    #[test]
    #[expected_failure(abort_code = eth_abi::E_OUT_OF_BYTES)]
    public fun decode_address_failed() {
        let addr: vector<u8> =  vector[
            0x00, 0x00, 0x00, 0x00, 0xAD, 0xa8, 0x0b, 0x6a,
            0xe7, 0xF0, 0x09, 0x60, 0xC3, 0x02, 0x0b, 0x5E,
            0x97, 0xAA, 0xAC, 0xCc, 0x3a, 0x46, 0x74, 0xf9
        ];

        let mut stream: eth_abi::ABIStream = eth_abi::new_stream(addr);

        eth_abi::decode_address(&mut stream);
    }


    #[test]
    public fun decode_u256_value() {
        let encoded_uint256: vector<u8> = vector[
            0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0
        ];

        let value: u256 = eth_abi::decode_u256_value(encoded_uint256);

        assert!(value == 256);
    }

    #[test]
    #[expected_failure(abort_code = eth_abi::E_INVALID_U256_LENGTH)]
    public fun decode_u256_value_failed() {
        let encoded_uint256: vector<u8> = vector[
            0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0
        ];

        eth_abi::decode_u256_value(encoded_uint256);
    }

    #[test]
    public fun decode_u256() {
        let data: vector<u8> =  vector[
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00
        ];

        let mut stream: eth_abi::ABIStream = eth_abi::new_stream(data);

        let value: u256 = eth_abi::decode_u256(&mut stream);
        assert!(value == 256);
        assert!(eth_abi::get_cur(&stream) == 32);
    }

    #[test]
    public fun decode_u32() {
        let data: vector<u8> =  vector[
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00
        ];

        let mut stream: eth_abi::ABIStream = eth_abi::new_stream(data);

        let value: u64 = eth_abi::decode_u64(&mut stream);
        assert!(value == 256);
        assert!(eth_abi::get_cur(&stream) == 32);
    }

    #[test]
    public fun decode_bool() {
        let data: vector<u8> =  vector[
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01
        ];

        let mut stream: eth_abi::ABIStream = eth_abi::new_stream(data);

        let value: bool = eth_abi::decode_bool(&mut stream);
        assert!(value == true);
        assert!(eth_abi::get_cur(&stream) == 32);
    }

    #[test]
    #[expected_failure(abort_code = eth_abi::E_INVALID_BOOL)]
    public fun decode_bool_failed_invalid_bool() {
        let data: vector<u8> =  vector[
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02
        ];

        let mut stream: eth_abi::ABIStream = eth_abi::new_stream(data);

        eth_abi::decode_bool(&mut stream);
    }

    #[test]
    #[expected_failure(abort_code = eth_abi::E_OUT_OF_BYTES)]
    public fun decode_bool_failed_out_of_bytes() {
        let data: vector<u8> =  vector[
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02
        ];

        let mut stream: eth_abi::ABIStream = eth_abi::new_stream(data);

        eth_abi::decode_bool(&mut stream);
    }

    #[test]
    public fun decode_bytes32() {
        let data: vector<u8> =  vector[
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x02
        ];

        let mut stream: eth_abi::ABIStream = eth_abi::new_stream(data);

        let b32: vector<u8> = eth_abi::decode_bytes32(&mut stream);

        assert!(data == b32);
        assert!(eth_abi::get_cur(&stream) == 32);
    }

    #[test]
    #[expected_failure(abort_code = eth_abi::E_OUT_OF_BYTES)]
    public fun decode_bytes32_failed() {
        let data: vector<u8> =  vector[
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x02
        ];

        let mut stream: eth_abi::ABIStream = eth_abi::new_stream(data);

        eth_abi::decode_bytes32(&mut stream);
    }

    #[test]
    public fun decode_bytes() {
        let data: vector<u8> = vector[
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02,
            0xFF, 0xAA, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00
        ];

        let mut stream: eth_abi::ABIStream = eth_abi::new_stream(data);

        let bytes: vector<u8> = eth_abi::decode_bytes(&mut stream);

        assert!(bytes == vector[0xFF, 0xAA]);
        assert!(eth_abi::get_cur(&stream) == 64);
    }

    #[test]
    public fun decode_bytes_long() {
        let data: vector<u8> = vector[
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x21,
            0xFF, 0xAA, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
            0xFF, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00
        ];

        let mut stream: eth_abi::ABIStream = eth_abi::new_stream(data);

        let bytes: vector<u8> = eth_abi::decode_bytes(&mut stream);

        assert!(bytes == vector[
            0xFF, 0xAA, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
            0xFF
        ]);
        assert!(eth_abi::get_cur(&stream) == 96);
    }

    #[test]
    #[expected_failure(abort_code = eth_abi::E_OUT_OF_BYTES)]
    public fun decode_bytes_failed() {
        let data: vector<u8> = vector[
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04,
            0xFF, 0xAA
        ];

        let mut stream: eth_abi::ABIStream = eth_abi::new_stream(data);

        eth_abi::decode_bytes(&mut stream);
    }

    #[test]
    #[expected_failure(abort_code = eth_abi::E_OUT_OF_BYTES)]
    public fun decode_bytes_failed_padding() {
        let data: vector<u8> = vector[
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04,
            0xFF, 0xAA, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
        ];

        let mut stream: eth_abi::ABIStream = eth_abi::new_stream(data);

        eth_abi::decode_bytes(&mut stream);
    }
}
