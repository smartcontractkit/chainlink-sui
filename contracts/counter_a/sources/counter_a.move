module 0x4c43754f0bf40d1db7b82b977f784b2e26c181d46d162f132333a42b0ee00294::counter {
    public struct Counter has store, key {
        id: 0x2::object::UID,
        val: u64,
    }
    
    public fun new(arg0: &mut 0x2::tx_context::TxContext) : Counter {
        Counter{
            id  : 0x2::object::new(arg0), 
            val : 0,
        }
    }
    
    public fun decrement(arg0: &mut Counter) {
        assert!(arg0.val > 0, 0);
        arg0.val = arg0.val - 1;
    }
    
    public fun increment(arg0: &mut Counter) {
        arg0.val = arg0.val + 1;
    }
    
    public fun val(arg0: &Counter) : u64 {
        arg0.val
    }
    
    // decompiled from Move bytecode v6
}

