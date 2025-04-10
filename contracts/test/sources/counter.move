module test::counter {
    use sui::object::{Self, UID, ID};
    use sui::transfer;
    use sui::tx_context::{Self, TxContext};
    use sui::event;

    // Event emitted when counter is incremented
    public struct CounterIncremented has copy, drop {
        counter_id: ID,
        new_value: u64
    }

    public struct AdminCap has key, store {
        id: UID
    }

    public struct Counter has key, store {
        id: UID,
        value: u64
    }

    fun init(ctx: &mut TxContext) {
        let counter = Counter { 
            id: object::new(ctx), 
            value: 0 
        };

        let admin_cap = AdminCap {
            id: object::new(ctx)
        };

        transfer::share_object(counter);
        transfer::transfer(admin_cap, ctx.sender());
    }

    /// Create and share a Counter object
    public entry fun initialize(ctx: &mut TxContext) {
        let counter = Counter { 
            id: object::new(ctx), 
            value: 0 
        };
        transfer::share_object(counter);
    }

    /// Increment counter by 1
    public entry fun increment(counter: &mut Counter) {
        counter.value = counter.value + 1;
        
        // Emit an event
        event::emit(CounterIncremented {
            counter_id: object::id(counter),
            new_value: counter.value
        });
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

    public entry fun increment_by_two_no_context(_admin: &AdminCap, counter: &mut Counter) {
        counter.value = counter.value + 2;
        
        // Emit an event
        event::emit(CounterIncremented {
            counter_id: object::id(counter),
            new_value: counter.value
        });
    }

    /// Increment counter by a*b
    public entry fun increment_mult(
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

    /// Get the value of the count
    public entry fun get_count(counter: &Counter): u64 {
        counter.value
    }

    public fun get_count_no_entry(counter: &Counter): u64 {
        counter.value
    }
}
