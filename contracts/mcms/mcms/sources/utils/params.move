module mcms::params;

use sui::bcs;

const E_CMP_VECTORS_DIFF_LEN: u64 = 0;

public fun encode_uint<T: drop>(input: T, num_bytes: u64): vector<u8> {
    let mut bcs_bytes = bcs::to_bytes(&input);

    let len = bcs_bytes.length();
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
    assert!(len == b.length(), E_CMP_VECTORS_DIFF_LEN);

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
