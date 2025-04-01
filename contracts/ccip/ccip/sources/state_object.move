module ccip::state_object {
    use sui::dynamic_object_field as dof;

    const E_MODULE_ALREADY_EXISTS: u64 = 1;
    const E_MODULE_DOES_NOT_EXISTS: u64 = 2;

    // currently we only create 1 capability to manage all the CCIP components like fee_quoter, rmn_remote
    // TODO: figure out if we need to create multiple capabilities for each CCIP component
    public struct OwnerCap has key, store {
        id: UID
    }

    public struct CCIPObjectRef has key {
        id: UID
    }

    fun init(ctx: &mut TxContext) {
        let ref = CCIPObjectRef {
            id: object::new(ctx)
        };

        let owner_cap = OwnerCap {
            id: object::new(ctx)
        };

        // TODO: another approach is to transfer to owner.
        transfer::share_object(ref);
        transfer::transfer(owner_cap, ctx.sender());
    }

    // TODO: we may need to include link token here
    public fun add<T: key + store>(_: &OwnerCap, ref: &mut CCIPObjectRef, name: vector<u8>, obj: T) {
        // TODO: or remove an existing object with this name?
        assert!(!dof::exists_(&ref.id, name), E_MODULE_ALREADY_EXISTS);
        dof::add(&mut ref.id, name, obj);
    }

    public fun contains(ref: &CCIPObjectRef, name: vector<u8>): bool {
        dof::exists_(&ref.id, name)
    }

    public fun remove<T: key + store>(_: &OwnerCap, ref: &mut CCIPObjectRef, name: vector<u8>): T {
        assert!(dof::exists_(&ref.id, name), E_MODULE_DOES_NOT_EXISTS);
        dof::remove(&mut ref.id, name)
    }

    public(package) fun borrow<T: key + store>(ref: &CCIPObjectRef, name: vector<u8>): &T {
        dof::borrow(&ref.id, name)
    }

    public(package) fun borrow_mut<T: key + store>(_: &OwnerCap, ref: &mut CCIPObjectRef, name: vector<u8>): &mut T {
        dof::borrow_mut(&mut ref.id, name)
    }

    #[test_only]
    public(package) fun create(ctx: &mut TxContext): (OwnerCap, CCIPObjectRef) {
        let ref = CCIPObjectRef {
            id: object::new(ctx)
        };

        let owner_cap = OwnerCap {
            id: object::new(ctx)
        };

        (owner_cap, ref)
    }

    #[test_only]
    public(package) fun destroy_owner_cap(cap: OwnerCap) {
        let OwnerCap { id } = cap;
        object::delete(id);
    }

    #[test_only]
    public(package) fun destroy_state_object(ref: CCIPObjectRef) {
        let CCIPObjectRef { id } = ref;
        object::delete(id);
    }
}

#[test_only]
module ccip::state_object_test {
    use ccip::state_object::{Self, OwnerCap, CCIPObjectRef};
    use sui::test_scenario::{Self, Scenario};

    fun set_up_test(): (Scenario, OwnerCap, CCIPObjectRef, TestObject) {
        let mut scenario = test_scenario::begin(@0x1);
        let ctx = scenario.ctx();

        let (owner_cap, ref) = state_object::create(ctx);
        let obj = TestObject {
            id: object::new(ctx)
        };
        (scenario, owner_cap, ref, obj)
    }

    fun tear_down_test(scenario: Scenario, owner_cap: OwnerCap, ref: CCIPObjectRef) {
        state_object::destroy_owner_cap(owner_cap);
        state_object::destroy_state_object(ref);
        test_scenario::end(scenario);
    }

    public struct TestObject has key, store {
        id: UID
    }

    #[test]
    public fun test_add() {
        let (scenario, owner_cap, mut ref, obj) = set_up_test();

        state_object::add(&owner_cap, &mut ref, b"test", obj);
        assert!(state_object::contains(&ref, b"test"));

        tear_down_test(scenario, owner_cap, ref);
    }

    #[test]
    public fun test_remove() {
        let (scenario, owner_cap, mut ref, obj) = set_up_test();

        state_object::add(&owner_cap, &mut ref, b"test", obj);
        assert!(state_object::contains(&ref, b"test"));

        let obj2: TestObject = state_object::remove(&owner_cap, &mut ref, b"test");
        assert!(!state_object::contains(&ref, b"test"));

        let TestObject { id } = obj2;
        object::delete(id);

        tear_down_test(scenario, owner_cap, ref);
    }

    #[test]
    public fun test_borrow() {
        let (scenario, owner_cap, mut ref, obj) = set_up_test();

        state_object::add(&owner_cap, &mut ref, b"test", obj);
        assert!(state_object::contains(&ref, b"test"));

        let _obj2: &TestObject = state_object::borrow(&ref, b"test");
        assert!(state_object::contains(&ref, b"test"));

        tear_down_test(scenario, owner_cap, ref);
    }

    #[test]
    public fun test_borrow_mut() {
        let (scenario, owner_cap, mut ref, obj) = set_up_test();

        state_object::add(&owner_cap, &mut ref, b"test", obj);
        assert!(state_object::contains(&ref, b"test"));

        let _obj2: &mut TestObject = state_object::borrow_mut(&owner_cap, &mut ref, b"test");
        assert!(state_object::contains(&ref, b"test"));

        tear_down_test(scenario, owner_cap, ref);
    }
}