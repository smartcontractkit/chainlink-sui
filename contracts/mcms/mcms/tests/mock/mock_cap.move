module 0x987::mock_cap;

public struct MockCap has key, store {
    id: UID,
}

public struct MockWitness has drop {}

public fun new(ctx: &mut TxContext): MockCap {
    MockCap { id: object::new(ctx) }
}

#[test_only]
public fun test_witness(): MockWitness {
    MockWitness {}
}
