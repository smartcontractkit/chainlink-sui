module test::counter {
    use sui::object::{Self, UID, ID};
    use sui::derived_object;
    use sui::transfer;
    use sui::tx_context::{Self, TxContext};
    use sui::event;
    use sui::address;
    use sui::coin::{Self as coin, Coin};
    use sui::balance::{Self, Balance};
    use std::vector;
    use std::ascii;
    use std::type_name;
    use std::string;
    use test_secondary::state_object::CCIPObjectRef;

    const EInvalidCounterValue: u64 = 1;
    const EInvalidBytesLength: u64 = 2;

    public struct COUNTER has drop {}

    // Event emitted when counter is incremented
    public struct CounterIncremented has copy, drop {
        counter_id: ID,
        new_value: u64
    }

    public struct CounterDecremented has copy, drop {
        event_type: std::string::String,
        counter_id: ID,
        new_value: u64
    }

    public struct CounterBytes has copy, drop {
        bytes: vector<u8>,
        nested: NestedCounterBytes,
        values: vector<u64>,
    }

    public struct NestedCounterBytes has copy, drop {
        value: u64,
        bytes: vector<u8>,
    }

    public struct DoubleCheckCCIPPointer has copy, drop {
        addresses: vector<address>
    }

    public struct AdminCap has key, store {
        id: UID
    }

    public struct Counter has key, store {
        id: UID,
        value: u64
    }

    public struct CounterObject has key {
        id: UID,
    }

    public struct SomeObject has key, store {
        id: UID,
    }

    // Pointer to reference both Counter and AdminCap objects
    public struct CounterPointer has key, store {
        id: UID,
        counter_object_id: address,
    }

    // Struct that contains a list of addresses
    public struct AddressList has copy, drop {
        addresses: vector<address>,
        count: u64,
    }

    // Simple struct for testing BCS decoding
    public struct SimpleResult has copy, drop {
        value: u64,
    }

    public struct ComplexResult has copy, drop {
        count: u64,
        addr: address,
        is_complex: bool,
        bytes: vector<u8>,
    }

    public struct NestedStruct has copy, drop {
        is_nested: bool,
        double_count: u64,
        nested_struct: ComplexResult,
        nested_simple_struct: SimpleResult,
    }

    public struct MultiNestedStruct has copy, drop {
        is_multi_nested: bool,
        double_count: u64,
        nested_struct: NestedStruct,
        nested_simple_struct: SimpleResult,
    }

    public struct ConfigInfo has store, drop, copy {
        config_digest: vector<u8>,
        big_f: u8,
        n: u8,
        is_signature_verification_enabled: bool
    }

    public struct OCRConfig has store, drop, copy {
        config_info: ConfigInfo,
        signers: vector<vector<u8>>,
        transmitters: vector<address>
    }

    fun init(_witness: COUNTER, ctx: &mut TxContext) {
        let mut counter_object = CounterObject {
            id: object::new(ctx),
        };

        let counter = Counter {
            id: derived_object::claim(&mut counter_object.id, b"Counter"),
            value: 0
        };

        let admin_cap = AdminCap {
            id: object::new(ctx)
        };

        // Create the pointer that references both objects
        let pointer = CounterPointer {
            id: object::new(ctx),
            counter_object_id: object::id_to_address(object::borrow_id(&counter_object)),
        };

        let pointer2 = CounterPointer {
            id: object::new(ctx),
            counter_object_id: object::id_to_address(object::borrow_id(&counter_object)),
        };

        let tn = type_name::get_with_original_ids<COUNTER>();
        let package_bytes = ascii::into_bytes(tn.get_address());
        let package_id = address::from_ascii_bytes(&package_bytes);

        transfer::share_object(counter);
        transfer::share_object(counter_object);

        transfer::transfer(admin_cap, tx_context::sender(ctx));
        transfer::transfer(pointer, package_id);
        transfer::transfer(pointer2, tx_context::sender(ctx));
    }

    /// Create and share a Counter object
    public fun initialize(ctx: &mut TxContext) {
        let counter = Counter {
            id: object::new(ctx),
            value: 0
        };
        transfer::share_object(counter);
    }

    public fun type_and_version(): std::string::String {
        std::string::utf8(b"Counter 1.6.0")
    }

    /// Increment counter by 1
    public fun increment(counter: &mut Counter) {
        counter.value = counter.value + 1;

        // Emit an event
        event::emit(CounterIncremented {
            counter_id: object::id(counter),
            new_value: counter.value
        });
    }

    /// Increment counter with pointer dependency
    public fun get_value_with_pointer_dependency(counter: &mut Counter, pointer: &CCIPObjectRef): u64 {
        event::emit(DoubleCheckCCIPPointer {
            addresses: test_secondary::state_object::get_package_ids(pointer),
        });

        5
    }

    /// Decrement counter by 1
    public fun decrement(counter: &mut Counter) {
        counter.value = counter.value - 1;

        // Emit an event
        event::emit(CounterDecremented {
            event_type: string::utf8(b"CounterDecremented"),
            counter_id: object::id(counter),
            new_value: counter.value
        });
    }

    /// Create a "Counter" object and return it (give up ownership without sharing)
    public fun create(ctx: &mut TxContext): Counter {
        Counter {
            id: object::new(ctx),
            value: 0
        }
    }

    public fun increment_by_one(counter: &mut Counter, _ctx: &mut TxContext): u64 {
        counter.value = counter.value + 1;

        // Emit an event
        event::emit(CounterIncremented {
            counter_id: object::id(counter),
            new_value: counter.value
        });

        counter.value
    }

    public fun increment_by_one_no_context(counter: &mut Counter): u64 {
        counter.value = counter.value + 1;

        // Emit an event
        event::emit(CounterIncremented {
            counter_id: object::id(counter),
            new_value: counter.value
        });

        counter.value
    }

    public fun increment_by_two(_admin: &AdminCap, counter: &mut Counter, _ctx: &mut TxContext) {
        counter.value = counter.value + 2;

        // Emit an event
        event::emit(CounterIncremented {
            counter_id: object::id(counter),
            new_value: counter.value
        });
    }

    public fun increment_by_two_no_context(_admin: &AdminCap, counter: &mut Counter) {
        counter.value = counter.value + 2;

        // Emit an event
        event::emit(CounterIncremented {
            counter_id: object::id(counter),
            new_value: counter.value
        });
    }

    public fun increment_by(counter: &mut Counter, by: u64) {
        // This is a test to ensure that the failed transaction is findable by the indexer
        assert!(by <= 999, EInvalidCounterValue);
        
        counter.value = counter.value + by;

        // Emit an event
        event::emit(CounterIncremented {
            counter_id: object::id(counter),
            new_value: counter.value
        });
    }

    /// Increment counter by a*b
    public fun increment_mult(
        counter: &mut Counter,
        a: u64,
        b: u64,
        _ctx: &mut TxContext
    ) {
        counter.value = counter.value + (a * b);

        // Emit an event
        event::emit(CounterIncremented {
            counter_id: object::id(counter),
            new_value: counter.value
        });
    }

    /// receive byte array
    public fun increment_by_bytes_length(
        counter: &mut Counter,
        bytes: vector<u8>,
    ) {
        assert!(bytes.length() <= 4, EInvalidBytesLength);
        counter.value = counter.value + (bytes.length());

        // Emit an event
        event::emit(CounterIncremented {
            counter_id: object::id(counter),
            new_value: counter.value
        });
    }

    /// Emit counter bytes
    public fun emit_counter_bytes(_ctx: &mut TxContext) {
        // Emit an event
        event::emit(CounterBytes {
            bytes: b"test",
            values: vector[1, 2, 3, 4],
            nested: NestedCounterBytes {
                value: 42,
                bytes: b"test",
            },
        });
    }

    /// Get the value of the count
    public fun get_count(counter: &Counter): u64 {
        counter.value
    }

    // Used to test the pointer functionality in CR
    public fun get_count_using_pointer(counter: &Counter): u64 {
        counter.value
    }

    public fun get_count_no_entry(counter: &Counter): u64 {
        counter.value
    }
    /// Get the value of a coin of generic type T
    public fun get_coin_value<T>(coin: &Coin<T>): u64 {
        coin::value(coin)
    }

    /// Returns a struct containing a list of addresses
    public fun get_address_list(): AddressList {
        let mut addresses = vector::empty<address>();

        // Add some sample addresses
        vector::push_back(&mut addresses, @0x1);
        vector::push_back(&mut addresses, @0x2);
        vector::push_back(&mut addresses, @0x3);
        vector::push_back(&mut addresses, @0x4);

        AddressList {
            addresses,
            count: vector::length(&addresses),
        }
    }

    public fun get_simple_result(): SimpleResult {
        SimpleResult {
            value: 42,
        }
    }

    /// Returns a simple struct with a single value for BCS testing
    public fun get_result_struct(): ComplexResult {
        let mut random_bytes = vector::empty<u8>();
        vector::push_back(&mut random_bytes, 1);
        vector::push_back(&mut random_bytes, 2);
        vector::push_back(&mut random_bytes, 3);
        vector::push_back(&mut random_bytes, 4);

        ComplexResult {
            count: 42,
            addr: @0x1,
            is_complex: true,
            bytes: random_bytes,
        }
    }

    /// Returns a nested struct with a complex struct and a simple struct
    public fun get_nested_result_struct(): NestedStruct {
        let mut random_bytes = vector::empty<u8>();
        vector::push_back(&mut random_bytes, 1);
        vector::push_back(&mut random_bytes, 2);
        vector::push_back(&mut random_bytes, 3);
        vector::push_back(&mut random_bytes, 4);

        NestedStruct {
            is_nested: true,
            double_count: 42,
            nested_struct: ComplexResult {
                count: 42,
                addr: @0x1,
                is_complex: true,
                bytes: random_bytes,
            },
            nested_simple_struct: SimpleResult {
                value: 42,
            },
        }
    }

    /// Returns a multi nested struct with a nested struct and a simple struct
    public fun get_multi_nested_result_struct(): MultiNestedStruct {
        MultiNestedStruct {
            is_multi_nested: true,
            double_count: 42,
            nested_struct: get_nested_result_struct(),
            nested_simple_struct: SimpleResult {
                value: 42,
            },
        }
    }

    public fun get_tuple_struct(): (u64, address, bool, MultiNestedStruct) {
        let nested_result = MultiNestedStruct {
            is_multi_nested: true,
            double_count: 42,
            nested_struct: get_nested_result_struct(),
            nested_simple_struct: SimpleResult {
                value: 42,
            },
        };

        (42, @0x1, true, nested_result)
    }

    public fun get_ocr_config(): OCRConfig {
        let mut config_digest = vector::empty<u8>();
        vector::push_back(&mut config_digest, 2);
        vector::push_back(&mut config_digest, 3);
        vector::push_back(&mut config_digest, 4);
        vector::push_back(&mut config_digest, 5);

        let mut signers = vector::empty<vector<u8>>();
        vector::push_back(&mut signers, @0x5.to_bytes());
        let mut signer2 = vector::empty<u8>();
        vector::push_back(&mut signer2, 0x6);
        vector::push_back(&mut signer2, 0x6);
        vector::push_back(&mut signers, signer2);

        let mut transmitters = vector::empty<address>();
        vector::push_back(&mut transmitters, @0x7);

        OCRConfig {
            config_info: ConfigInfo {
                config_digest,
                big_f: 1,
                n: 2,
                is_signature_verification_enabled: true,
            },
            signers,
            transmitters,
        }
    }

    public fun get_vector_of_u8(): vector<u8> {
        let mut bytes = vector::empty<u8>();
        vector::push_back(&mut bytes, 1);
        vector::push_back(&mut bytes, 2);
        vector::push_back(&mut bytes, 3);
        vector::push_back(&mut bytes, 4);
        bytes
    }

    public fun get_vector_of_addresses(): vector<address> {
        let mut addresses = vector::empty<address>();
        vector::push_back(&mut addresses, @0x1);
        vector::push_back(&mut addresses, @0x2);
        vector::push_back(&mut addresses, @0x3);
        vector::push_back(&mut addresses, @0x4);
        addresses
    }

    public fun get_vector_of_vectors_of_u8(): vector<vector<u8>> {
        let mut vectors = vector::empty<vector<u8>>();
        vector::push_back(&mut vectors, @0x1.to_bytes());
        vector::push_back(&mut vectors, @0x2.to_bytes());
        vector::push_back(&mut vectors, @0x3.to_bytes());
        vector::push_back(&mut vectors, @0x4.to_bytes());
        vectors
    }

    #[allow(lint(self_transfer))]
    public fun create_many_objects(count: u64, ctx: &mut TxContext) {
        let mut i = 0;
        
        while (i < count) {
            let obj = SomeObject {
                id: object::new(ctx),
            };
            transfer::transfer(obj, tx_context::sender(ctx));
            i = i + 1;
        }
    }
}
