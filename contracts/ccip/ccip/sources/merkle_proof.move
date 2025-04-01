module ccip::merkle_proof {
    use sui::hash;

    const E_VECTOR_LENGTH_MISMATCH: u64 = 1;

    const LEAF_DOMAIN_SEPARATOR: vector<u8> = x"0000000000000000000000000000000000000000000000000000000000000000";
    const INTERNAL_DOMAIN_SEPARATOR: vector<u8> = x"0000000000000000000000000000000000000000000000000000000000000001";

    public fun leaf_domain_separator(): vector<u8> {
        LEAF_DOMAIN_SEPARATOR
    }

    public fun vector_u8_gt(a: &vector<u8>, b: &vector<u8>): bool {
        let len = vector::length(a);
        assert!(
            len == vector::length(b), E_VECTOR_LENGTH_MISMATCH
        );

        let mut i = 0;
        // compare each byte until not equal
        while (i < len) {
            let byte_a = *vector::borrow(a, i);
            let byte_b = *vector::borrow(b, i);
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

    public fun merkle_root_simple(leaf: vector<u8>, proofs: vector<vector<u8>>): vector<u8> {
        vector::fold!(proofs, leaf, |acc, proof| hash_pair(acc, proof))
    }

    // preserve this function in case we need to bring it back for multi-proof verification
    // public fun merkle_root(
    //     leaves: &vector<vector<u8>>,
    //     proofs: &vector<vector<u8>>,
    //     proof_flag_bits: u256
    // ): vector<u8> {
    //     let leaves_len = vector::length(leaves);
    //     let proofs_len = vector::length(proofs);
    //
    //     assert!(leaves_len > 0, E_LEAVES_CANNOT_BE_EMPTY);
    //     assert!(
    //         leaves_len <= (MAX_NUM_HASHES + 1) && proofs_len <= (MAX_NUM_HASHES + 1),
    //         E_INVALID_PROOF
    //     );
    //
    //     let total_hashes = leaves_len + proofs_len - 1;
    //     assert!(total_hashes <= MAX_NUM_HASHES, E_INVALID_PROOF);
    //     assert!(total_hashes > 0, E_SINGLE_LEAF);
    //
    //     let mut hashes = vector[];
    //     let mut leaf_pos = 0u64;
    //     let mut hash_pos = 0u64;
    //     let mut proof_pos = 0u64;
    //     let mut i = 0u64;
    //
    //     while (i < total_hashes) {
    //         let mut a;
    //         // total_hashes <= MAX_NUM_HASHES so i < MAX_NUM_HASHES and fit inside a u8.
    //         let current_bit = 1 << (i as u8);
    //         if ((proof_flag_bits & current_bit) == current_bit) {
    //             if (leaf_pos < leaves_len) {
    //                 a = *vector::borrow(leaves, leaf_pos);
    //                 leaf_pos = leaf_pos + 1;
    //             } else {
    //                 assert!(
    //                     hash_pos < vector::length(&hashes),
    //                     E_INVALID_PROOF
    //                 );
    //                 a = *vector::borrow(&hashes, hash_pos);
    //                 hash_pos = hash_pos + 1;
    //             }
    //         } else {
    //             assert!(proof_pos < proofs_len, E_INVALID_PROOF);
    //             a = *vector::borrow(proofs, proof_pos);
    //             proof_pos = proof_pos + 1;
    //         };
    //
    //         let mut b;
    //         if (leaf_pos < leaves_len) {
    //             b = *vector::borrow(leaves, leaf_pos);
    //             leaf_pos = leaf_pos + 1;
    //         } else {
    //             assert!(
    //                 hash_pos < vector::length(&hashes),
    //                 E_INVALID_PROOF
    //             );
    //             b = *vector::borrow(&hashes, hash_pos);
    //             hash_pos = hash_pos + 1;
    //         };
    //
    //         assert!(hash_pos <= i, E_INVALID_PROOF);
    //
    //         let hash = hash_pair(a, b);
    //         vector::push_back(&mut hashes, hash);
    //         i = i + 1;
    //     };
    //
    //     assert!(
    //         hash_pos == (total_hashes - 1)
    //             && leaf_pos == leaves_len
    //             && proof_pos == proofs_len,
    //         E_INVALID_PROOF
    //     );
    //
    //     *vector::borrow(&hashes, total_hashes - 1)
    // }

    fun hash_internal_node(left: vector<u8>, right: vector<u8>): vector<u8> {
        let mut data = INTERNAL_DOMAIN_SEPARATOR;
        vector::append(&mut data, left);
        vector::append(&mut data, right);
        hash::keccak256(&data)
    }

    /// Hashes a pair of byte vectors, ordering them lexographically
    fun hash_pair(a: vector<u8>, b: vector<u8>): vector<u8> {
        if (!vector_u8_gt(&a, &b)) {
            hash_internal_node(a, b)
        } else {
            hash_internal_node(b, a)
        }
    }
}

#[test_only]
module ccip::merkle_proof_test {
    use ccip::merkle_proof;

    #[test]
    #[expected_failure(abort_code = merkle_proof::E_VECTOR_LENGTH_MISMATCH)]
    public fun vector_u8_gt_failed() {
        let a = vector[1, 2, 3];
        let b = vector[1, 2];
        merkle_proof::vector_u8_gt(&a, &b);
    }

    #[test]
    public fun vector_u8_gt() {
        let a = vector[1, 2, 5];
        let b = vector[1, 2, 4];
        assert!(merkle_proof::vector_u8_gt(&a, &b));
    }

    #[test]
    public fun merkle_root_simple_1() {
        let leaf = x"a20c0244af79697a4ef4e2378c9d5d14cbd49ddab3427b12594c7cfa67a7f240";
        let proofs = vector[
            x"7b43f2a9158ed1c62f904d7e8b195a3cb467e5821d9f46a0c873b25df831e994",
            x"2c8fd561a437b9e04a1f85c26d930b7ef148ac259b3e70d4168ac75fb309e24d",
            x"531ab8c46f920d75e34c9a2785f63b14d08e59a27c1fb543ec962d805bf43a17",
            x"9c25e74fb138d60a72c95e831bf047ae952c6d138f54ba317ec509d24a86f35b",
            x"41f82ea763b905dc8a347fc15b920ed648b31ae57c29846fd03ba5129e57c83f",
            x"842bf60d953cae711fc758b249e0268da35f9138d40c7ae26bf5138c47b925d0",
            x"16c94fb32a857d31e85c940b73df42a61e873cf569b02d954a8e61c703d85fa2",
            x"6d359a0cf47b238e51b72fd945ac6813e0974c821bd53af68f27b950e47c159d",
            x"a8541cf73e926b0dc5782ab38f41e65917d0843bac259e70f34c8d16b52e9751",
            x"3fd28a51c70e942b76f81da54c8319e0b26f37c9058e40d79a13f52c7b846e0a",
            x"9147c38e25ba0f6d349c52e817d57a43b61f852ca960f43b7e950dc8512ab34f",
            x"5e03d9724ab1861fc52d970e63fa38854cb27d19e650a8239f14d76b823cf50a",
            x"c78e24b951f63a820d751ce49b37a05fd26813f94c852bae067d319a4ec58f27",
            x"2ad57f138c46b059e29a31c7046df8258e3ba4721fd950973e860cb54a1de763",
            x"951ce8427bd309f65a832db46f17c08e359c4fa126d8703b851ef7629d04b853",
            x"4ea731c85f920b74d61e8a43b52cf960873de51972ab249c530fd148b67a35ec",
            x"b3298f56e047c13a820d751ce49b37a05fd26813f94c852bae067d319a4ec58f",
            x"7af42d9138c50e67b249831cd56a04ef952c7b138f54ba317ec509d24a86f35b",
            x"51e8723f9c04ba256df1388e47a013c95e821bd47f369a0ce528b34f97610d8a",
            x"c43b950d72f81e864ab327d95c01e47b932fa8560cd148bf359e67128d4ca53e",
            x"08d46ba23f951ce8427bd309f65a832db46f17c08e359c4fa126d8703b851ef7",
            x"832af54c97610db8359e67128d4ca53ef07921da568f37c20b741de6924fa813",
            x"1fc758b249e0268da35f9138d40c7ae26bf5138c47b925d0842bf60d953cae71",
            x"943d860f52ca711be459a227d08b35fc479e631ab52c78d10e964fa35b8219e7",
            x"2eb741d80f9653ac257e148bf239c06da51c843be65f9207d34ab1682f950c73",
            x"a95017e46b238cf53ad178429b0e855cb32af7610d943bc8561fd27945ac832e",
            x"751ce49b37a05fd26813f94c852bae067d319a4ec58f27b3298f56e047c13a82",
            x"4d853cf61ab27e09d45b832c970e65ac31f84ab56d238940e71cd692580fa47b",
            x"b62f985124c40a73e91d863fb54c27905ed8127fa349e0268db7540c953e71f8",
            x"39e257ac841bd56f039a42c8157eb029f48d36a15c02eb741d863fb54c27905e",
            x"0c953e71f82ab62f985124c40a73e91d863fb54c27905ed8127fa349e0268db7",
            x"65ac31f84ab56d238940e71cd692580fa47b4d853cf61ab27e09d45b832c970e",
            x"e259ac841bd56f039a42c8157eb029f48d36a15c02eb741d863fb54c27905e39"
        ];
        let res: vector<u8> = merkle_proof::merkle_root_simple(leaf, proofs);

        let hash_bytes = x"96a04377b175ac7aadd2e6babf30ef75b6df6d671f4709f765e65bbbd7b71339";
        assert!(res == hash_bytes);
    }

    #[test]
    public fun merkle_root_simple_2() {
        let leaf = x"8beab00297b94bf079fcd5893b0a33ebf6b0ce862cd06be07c87d3c63e1c4acf";

        let proofs = vector[
            x"1fa3b2c9458ed1762f904dae8b195a3cb467e5821d9f46a0c873b25df831e994",
            x"2c8fd561a437b9e04a1f85c26d930b7ef148ac259b3e70d4168ac75fb309e24d",
            x"531ab8c46f920d75e34c9a2785f63b14d08e59a27c1fb543ec962d805bf43a17",
            x"9c25e74fb138d60a72c95e831bf047ae952c6d138f54ba317ec509d24a86f35b",
            x"41f82ea763b905dc8a347fc15b920ed648b31ae57c29846fd03ba5129e57c83f",
            x"842bf60d953cae711fc758b249e0268da35f9138d40c7ae26bf5138c47b925d0",
            x"16c94fb32a857d31e85c940b73df42a61e873cf569b02d954a8e61c703d85fa2",
            x"6d359a0cf47b238e51b72fd945ac6813e0974c821bd53af68f27b950e47c159d",
            x"a8541cf73e926b0dc5782ab38f41e65917d0843bac259e70f34c8d16b52e9751",
            x"3fd28a51c70e942b76f81da54c8319e0b26f37c9058e40d79a13f52c7b846e0a",
            x"9147c38e25ba0f6d349c52e817d57a43b61f852ca960f43b7e950dc8512ab34f",
            x"5e03d9724ab1861fc52d970e63fa38854cb27d19e650a8239f14d76b823cf50a",
            x"c78e24b951f63a820d751ce49b37a05fd26813f94c852bae067d319a4ec58f27",
            x"2ad57f138c46b059e29a31c7046df8258e3ba4721fd950973e860cb54a1de763",
            x"951ce8427bd309f65a832db46f17c08e359c4fa126d8703b851ef7629d04b853",
            x"4ea731c85f920b74d61e8a43b52cf960873de51972ab249c530fd148b67a35ec",
            x"b3298f56e047c13a820d751ce49b37a05fd26813f94c852bae067d319a4ec58f",
            x"7af42d9138c50e67b249831cd56a04ef952c7b138f54ba317ec509d24a86f35b",
            x"51e8723f9c04ba256df1388e47a013c95e821bd47f369a0ce528b34f97610d8a",
            x"c43b950d72f81e864ab327d95c01e47b932fa8560cd148bf359e67128d4ca53e",
            x"08d46ba23f951ce8427bd309f65a832db46f17c08e359c4fa126d8703b851ef7",
            x"832af54c97610db8359e67128d4ca53ef07921da568f37c20b741de6924fa813",
            x"1fc758b249e0268da35f9138d40c7ae26bf5138c47b925d0842bf60d953cae71",
            x"943d860f52ca711be459a227d08b35fc479e631ab52c78d10e964fa35b8219e7",
            x"2eb741d80f9653ac257e148bf239c06da51c843be65f9207d34ab1682f950c73",
            x"a95017e46b238cf53ad178429b0e855cb32af7610d943bc8561fd27945ac832e",
            x"751ce49b37a05fd26813f94c852bae067d319a4ec58f27b3298f56e047c13a82",
            x"4d853cf61ab27e09d45b832c970e65ac31f84ab56d238940e71cd692580fa47b",
            x"b62f985124c40a73e91d863fb54c27905ed8127fa349e0268db7540c953e71f8",
            x"39e257ac841bd56f039a42c8157eb029f48d36a15c02eb741d863fb54c27905e",
            x"0c953e71f82ab62f985124c40a73e91d863fb54c27905ed8127fa349e0268db7",
            x"65ac31f84ab56d238940e71cd692580fa47b4d853cf61ab27e09d45b832c970e",
            x"e259ac841bd56f039a42c8157eb029f48d36a15c02eb741d863fb54c27905e39"
        ];

        let res: vector<u8> = merkle_proof::merkle_root_simple(leaf, proofs);

        let hash_bytes = x"397b8848147b5cfb6ea44a8c79515dac64bbcb8013925cce11df06432565e354";

        assert!(res == hash_bytes);
    }
}
