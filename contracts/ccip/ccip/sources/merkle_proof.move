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
