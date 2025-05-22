module test::counter {
    use sui::object::{Self, UID, ID};
    use sui::transfer;
    use sui::tx_context::{Self, TxContext};
    use sui::event;
    use sui::address;
    use std::vector;

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
    
    // Pointer to reference both Counter and AdminCap objects
    public struct CounterPointer has key, store {
        id: UID,
        counter_id: address,
        admin_cap_id: address,
    }

    fun init(ctx: &mut TxContext) {
        let counter = Counter { 
            id: object::new(ctx), 
            value: 0 
        };

        let admin_cap = AdminCap {
            id: object::new(ctx)
        };
        
        // Create the pointer that references both objects
        let pointer = CounterPointer {
            id: object::new(ctx),
            counter_id: object::id_to_address(object::borrow_id(&counter)),
            admin_cap_id: object::id_to_address(object::borrow_id(&admin_cap)),
        };

        transfer::share_object(counter);
        transfer::transfer(admin_cap, tx_context::sender(ctx));
        transfer::transfer(pointer, ctx.sender());
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

    public entry fun increment_by_two_no_context(_admin: &AdminCap, counter: &mut Counter) {
        counter.value = counter.value + 2;
        
        // Emit an event
        event::emit(CounterIncremented {
            counter_id: object::id(counter),
            new_value: counter.value
        });
    }

    public entry fun increment_by(counter: &mut Counter, by: u64) {
        counter.value = counter.value + by;

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


    public fun array_size<T: drop>(arr: vector<T>): u64 {
        vector::length(&arr)
    }

}
