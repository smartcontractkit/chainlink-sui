module test::counter {
    use sui::object::{Self, UID};
    use sui::transfer;
    use sui::tx_context::{Self, TxContext};

    struct Counter has key, store {
        id: UID,
        value: u64
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
    public entry fun increment(counter: &mut Counter, _ctx: &mut TxContext) {
        counter.value = counter.value + 1;
    }

    /// Increment counter by a*b
    public entry fun increment_mult(
        counter: &mut Counter,
        a: u64,
        b: u64,
        _ctx: &mut TxContext
    ) {
        counter.value = counter.value + (a * b);
    }

    /// Get the value of the count
    public entry fun get_count(counter: &Counter): u64 {
        counter.value
    }
}
