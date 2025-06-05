module test::complex {
    public struct SampleObject has key, store {
      id: UID,
      some_id: vector<u8>,
      some_number: u64,
      some_address: address,
      some_addresses: vector<address>,
    }

    public struct DroppableObject has drop {
      some_id: vector<u8>,
      some_number: u64,
      some_address: address,
      some_addresses: vector<address>,
    }

    public fun new_object_with_transfer(
        some_id: vector<u8>,
        some_number: u64,
        some_address: address,
        some_addresses: vector<address>,
        ctx: &mut TxContext
    ) {
        let obj = SampleObject {
            id: object::new(ctx),
            some_id,
            some_number,
            some_address,
            some_addresses
        };
        transfer::share_object(obj);
    }

    public fun new_object(
        some_id: vector<u8>,
        some_number: u64,
        some_address: address,
        some_addresses: vector<address>,
    ): DroppableObject {
        DroppableObject {
            some_id,
            some_number,
            some_address,
            some_addresses
        }
    }

    public fun flatten_address(
        some_address: address,
        some_addresses: vector<address>,
        _ctx: &mut TxContext
    ): vector<address> {
        let mut addresses = vector::empty<address>();
        vector::push_back(&mut addresses, some_address);
        let mut i = 0;
        let len = vector::length(&some_addresses);

        while (i < len) {
            let address = vector::borrow(&some_addresses, i);
            vector::push_back(&mut addresses, *address);
            i = i + 1;
        };

        addresses
    }

    public fun flatten_u8(
        input: vector<vector<u8>>,
        _ctx: &mut TxContext
    ): vector<u8> {
        let mut output = vector::empty<u8>();
        let mut i = 0;
        let len = vector::length(&input);

        while (i < len) {
            let inner_vector = vector::borrow(&input, i);
            let inner_len = vector::length(inner_vector);
            let mut j = 0;

            while (j < inner_len) {
                let byte = vector::borrow(inner_vector, j);
                vector::push_back(&mut output, *byte);
                j = j + 1;
            };

            i = i + 1;
        };

        output
    }

    public fun check_u128(
        input: u128,
        _ctx: &mut TxContext
    ): u128 {
        input
    }

    public fun check_u256(
        input: u256,
        _ctx: &mut TxContext
    ): u256 {
        input
    }
}
