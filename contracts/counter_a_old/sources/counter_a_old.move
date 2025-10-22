module 0x4f9a563b120c09717b76662b4256b10f470abdaf8dc10f0239e6ab5d744ce45e::counter {
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

