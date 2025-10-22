module counter_a::counter {
    const ENegativeDecrement: u64 = 0;

    public struct Counter has key,store{
        id: UID,
        val: u64
    }

    public fun new(ctx: &mut TxContext): Counter {
        Counter {
            id: object::new(ctx),
            val:0
        }
    }

    public fun increment(c: &mut Counter) {
        c.val = c.val + 1
    }

    public fun decrement(c: &mut Counter) {
        assert!(c.val > 0, ENegativeDecrement);
        c.val = c.val - 1
    }

    public fun val(c: &Counter) :u64 {
        c.val
    }

}