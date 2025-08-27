module ccip::bcs_helper;

use mcms::bcs_stream::{Self, BCSStream};

const EInvalidObjectAddress: u64 = 1;

public fun validate_obj_addr(addr: address, stream: &mut BCSStream) {
    let deserialized_address = bcs_stream::deserialize_address(stream);
    assert!(deserialized_address == addr, EInvalidObjectAddress);
}

public fun validate_obj_addrs(addrs: vector<address>, stream: &mut BCSStream) {
    let mut i = 0;
    while (i < addrs.length()) {
        validate_obj_addr(addrs[i], stream);
        i = i + 1;
    }
}
