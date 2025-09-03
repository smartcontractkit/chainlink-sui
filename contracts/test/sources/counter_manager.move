module test::counter_manager {
    use test::counter::Counter;
    use sui::tx_context::TxContext;
    use sui::object::{Self, UID};
    use sui::borrow;

    public struct CounterManager has key {
        id: UID,
        counter: borrow::Referent<Counter>
    }

    /// Create a new CounterManager for a shared Counter object
    public fun create(counter: Counter, ctx: &mut TxContext) {
        let manager = CounterManager {
            id: object::new(ctx),
            counter: borrow::new(counter, ctx)
        };
        transfer::share_object(manager);
    }

    /// Methods to borrow and return the counter
    public fun borrow_counter(manager: &mut CounterManager): (Counter, borrow::Borrow) {
        manager.counter.borrow()
    }

    public fun return_counter(manager: &mut CounterManager, counter: Counter, b: borrow::Borrow) {
        manager.counter.put_back(counter, b)
    }
}