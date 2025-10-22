module counter_b::counter {
    use sui::object::{Self, UID, ID};
    use sui::tx_context::TxContext;
    use sui::transfer;
    use sui::event;

    use counter_a::counter::{Self as counter_a, Counter as CounterA};

    /// This must match the module name in uppercase (COUNTER)
    public struct COUNTER has drop {}

    /// Event emitted when we interact with counter_a
    public struct CounterIntegratedEvent has copy, drop {
        counter_a_id: ID,
        new_value: u64,
    }

    /// Local object in counter_b that references a counter_a object
    public struct Counter has key, store {
        id: UID,
        counter_a_id: ID,
    }

    /// Package-level init function (called once at publish time)
    fun init(_witness: COUNTER, ctx: &mut TxContext) {
        let counter_a_instance = counter_a::new(ctx);

        let counter_b_instance = Counter {
            id: object::new(ctx),
            counter_a_id: object::id(&counter_a_instance),
        };

        // Share counter_b object
        transfer::share_object(counter_b_instance);

        // Publicly share counter_a instance
        transfer::public_share_object(counter_a_instance);
    }

    /// Increment the linked counter_a
    public fun increment_with_dependency(counter_a_ref: &mut CounterA) {
        counter_a::increment(counter_a_ref);
        event::emit(CounterIntegratedEvent {
            counter_a_id: object::id(counter_a_ref),
            new_value: counter_a::val(counter_a_ref),
        });
    }

    /// Read from counter_a dependency
    public fun get_value_from_dependency(counter_a_ref: &CounterA): u64 {
        counter_a::val(counter_a_ref)
    }
}
