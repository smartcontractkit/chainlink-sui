module test::echo {
    use std::string::String;
    use sui::event;

    public struct SingleValueEvent has copy, drop {
        value: u64
    }

    public struct NoConfigSingleValueEvent has copy, drop {
            value: u64
        }

    public struct DoubleValueEvent has copy, drop {
        number: u64,
        text: String
    }

    public struct TripleValueEvent has copy, drop {
        values: vector<vector<u8>>
    }

    public struct EventStore has key {
        id: UID
    }

    // Initialization function - automatically called once when module is published
    fun init(ctx: &mut TxContext) {
        transfer::share_object(
            EventStore {
                id: object::new(ctx)
            }
        );
    }

    public entry fun echo_with_events(
        _event_store: &EventStore,
        number: u64,
        text: String,
        bytes: vector<u8>,
        _ctx: &mut TxContext
    ) {
        // Emit events directly using Sui's event system
        event::emit(SingleValueEvent { value: number });
        event::emit(DoubleValueEvent { number, text });

        let mut values = vector::empty<vector<u8>>();
        vector::push_back(&mut values, bytes);
        event::emit(TripleValueEvent { values });
    }

    public fun simple_event_echo(
        val: u64,
        _ctx: &mut TxContext
    ): u64 {
        event::emit(SingleValueEvent { value: val });
        val
    }

        public fun no_config_event_echo(
            val: u64,
            _ctx: &mut TxContext
        ): u64 {
            event::emit(NoConfigSingleValueEvent { value: val });
            val
        }

    public fun echo_u64(val: u64): u64 {
        val
    }

    public fun echo_u256(val: u256): u256 {
        val
    }

    public fun echo_u32_u64_tuple(val1: u32, val2: u64): (u32, u64) {
        (val1, val2)
    }

    public fun echo_string(val: String): String {
        val
    }

    public fun echo_byte_vector(val: vector<u8>): vector<u8> {
        val
    }

    public fun echo_byte_vector_vector(val: vector<vector<u8>>): vector<vector<u8>> {
        val
    }
}
