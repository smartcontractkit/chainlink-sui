module ccip::state_object {
    use sui::event;
    use sui::dynamic_object_field as dof;

    const E_MODULE_ALREADY_EXISTS: u64 = 1;
    const E_MODULE_DOES_NOT_EXISTS: u64 = 2;
    const E_CANNOT_TRANSFER_TO_SELF: u64 = 3;
    const E_OWNER_CHANGED: u64 = 4;
    const E_NO_PENDING_TRANSFER: u64 = 5;
    const E_TRANSFER_NOT_ACCEPTED: u64 = 6;
    const E_TRANSFER_ALREADY_ACCEPTED: u64 = 7;
    const E_MUST_BE_PROPOSED_OWNER: u64 = 8;
    const E_PROPOSED_OWNER_MISMATCH: u64 = 9;
    const E_UNAUTHORIZED: u64 = 10;

    public struct OwnerCap has key, store {
        id: UID,
    }

    public struct OwnershipTransferRequested has copy, drop {
        from: address,
        to: address
    }

    public struct OwnershipTransferAccepted has copy, drop {
        from: address,
        to: address
    }

    public struct OwnershipTransferred has copy, drop {
        from: address,
        to: address
    }

    public struct CCIPObjectRef has key, store {
        id: UID,
        // this is not the owner of the CCIP object ref in SUI's concept
        // this object is a shared object and cannot be transferred and has no owner according to SUI
        // the owner here refers to the address which has the power to manage this ref
        current_owner: address,
        pending_transfer: Option<PendingTransfer>
    }

    public struct PendingTransfer has store, copy, drop {
        from: address,
        to: address,
        accepted: bool
    }

    fun init(ctx: &mut TxContext) {
        let ref = CCIPObjectRef {
            id: object::new(ctx),
            current_owner: ctx.sender(),
            pending_transfer: option::none()
        };
        let owner_cap = OwnerCap {
            id: object::new(ctx),
        };

        transfer::share_object(ref);
        transfer::transfer(owner_cap, ctx.sender());
    }

    // TODO: we may need to include link token here
    public(package) fun add<T: key + store>(ref: &mut CCIPObjectRef, name: vector<u8>, obj: T, ctx: &TxContext) {
        // TODO: or remove an existing object with this name?
        assert!(ctx.sender() == ref.current_owner, E_UNAUTHORIZED);
        assert!(!dof::exists_(&ref.id, name), E_MODULE_ALREADY_EXISTS);
        dof::add(&mut ref.id, name, obj);
    }

    public(package) fun contains(ref: &CCIPObjectRef, name: vector<u8>): bool {
        dof::exists_(&ref.id, name)
    }

    public(package) fun remove<T: key + store>(ref: &mut CCIPObjectRef, name: vector<u8>, ctx: &TxContext): T {
        assert!(ctx.sender() == ref.current_owner, E_UNAUTHORIZED);
        assert!(dof::exists_(&ref.id, name), E_MODULE_DOES_NOT_EXISTS);
        dof::remove(&mut ref.id, name)
    }

    public(package) fun borrow<T: key + store>(ref: &CCIPObjectRef, name: vector<u8>): &T {
        dof::borrow(&ref.id, name)
    }

    public(package) fun borrow_mut_with_ctx<T: key + store>(ref: &mut CCIPObjectRef, name: vector<u8>, ctx: &TxContext): &mut T {
        assert!(ctx.sender() == ref.current_owner, E_UNAUTHORIZED);

        dof::borrow_mut(&mut ref.id, name)
    }

    public fun transfer_ownership(
        ref: &mut CCIPObjectRef, to: address, ctx: &mut TxContext
    ) {
        let caller = ctx.sender();
        assert!(caller != to, E_CANNOT_TRANSFER_TO_SELF);
        assert!(ref.current_owner == caller, E_UNAUTHORIZED);

        ref.pending_transfer = option::some(
            PendingTransfer { from: caller, to, accepted: false }
        );

        event::emit(OwnershipTransferRequested { from: caller, to });
    }

    public fun accept_ownership(ref: &mut CCIPObjectRef, ctx: &mut TxContext) {
        assert!(
            ref.pending_transfer.is_some(),
            E_NO_PENDING_TRANSFER
        );

        let caller = ctx.sender();
        let pending_transfer = ref.pending_transfer.borrow_mut();

        assert!(
            pending_transfer.from == ref.current_owner,
            E_OWNER_CHANGED
        );
        assert!(
            pending_transfer.to == caller,
            E_MUST_BE_PROPOSED_OWNER
        );
        assert!(
            !pending_transfer.accepted,
            E_TRANSFER_ALREADY_ACCEPTED
        );

        pending_transfer.accepted = true;

        event::emit(OwnershipTransferAccepted { from: pending_transfer.from, to: caller });
    }

    public fun execute_ownership_transfer(
        ref: &mut CCIPObjectRef, to: address, ctx: &mut TxContext
    ) {
        let caller = ctx.sender();
        assert!(caller == ref.current_owner, E_UNAUTHORIZED);

        let pending_transfer = ref.pending_transfer.extract();

        // since ref is a shared object now, it's impossible to transfer its ownership
        assert!(
            pending_transfer.from == ref.current_owner,
            E_OWNER_CHANGED
        );
        assert!(
            pending_transfer.to == to, E_PROPOSED_OWNER_MISMATCH
        );
        assert!(
            pending_transfer.accepted, E_TRANSFER_NOT_ACCEPTED
        );

        // transfer the owner cap to the new owner
        // cannot transfer the shared object anymore
        ref.current_owner = pending_transfer.to;
        // the extract will remove the object within option wrapper
        // state.pending_transfer = option::none();

        event::emit(OwnershipTransferred { from: caller, to })
    }

    public(package) fun get_current_owner(ref: &CCIPObjectRef): address {
        ref.current_owner
    }

    #[test_only]
    public(package) fun create(ctx: &mut TxContext): CCIPObjectRef {
        CCIPObjectRef {
            id: object::new(ctx),
            current_owner: ctx.sender(),
            pending_transfer: option::none()
        }
    }

    #[test_only]
    public(package) fun destroy_state_object(ref: CCIPObjectRef) {
        let CCIPObjectRef { id, current_owner: _owner, pending_transfer: _pending_transfer } = ref;
        object::delete(id);
    }

    #[test_only]
    public(package) fun pending_transfer(ref: &CCIPObjectRef): (address, address, bool) {
        let pt = ref.pending_transfer;
        if (pt.is_none()) {
            return (@0x0, @0x0, false)
        };
        let pt = option::borrow(&ref.pending_transfer);

        (pt.from, pt.to, pt.accepted)
    }
}

#[test_only]
module ccip::state_object_test {
    use sui::test_scenario::{Self, Scenario};

    use ccip::state_object::{Self, CCIPObjectRef};

    const SENDER_1: address = @0x1;
    const SENDER_2: address = @0x2;

    fun set_up_test(): (Scenario, CCIPObjectRef, TestObject) {
        let mut scenario = test_scenario::begin(SENDER_1);
        let ctx = scenario.ctx();

        let ref = state_object::create(ctx);
        let obj = TestObject {
            id: object::new(ctx)
        };
        (scenario, ref, obj)
    }

    fun tear_down_test(scenario: Scenario, ref: CCIPObjectRef) {
        state_object::destroy_state_object(ref);
        test_scenario::end(scenario);
    }

    public struct TestObject has key, store {
        id: UID
    }

    #[test]
    public fun test_add() {
        let (mut scenario, mut ref, obj) = set_up_test();
        let ctx = scenario.ctx();

        state_object::add(&mut ref, b"test", obj, ctx);
        assert!(state_object::contains(&ref, b"test"));

        tear_down_test(scenario, ref);
    }

    #[test]
    public fun test_remove() {
        let (mut scenario, mut ref, obj) = set_up_test();
        let ctx = scenario.ctx();

        state_object::add(&mut ref, b"test", obj, ctx);
        assert!(state_object::contains(&ref, b"test"));

        let obj2: TestObject = state_object::remove(&mut ref, b"test", ctx);
        assert!(!state_object::contains(&ref, b"test"));

        let TestObject { id } = obj2;
        object::delete(id);

        tear_down_test(scenario, ref);
    }

    #[test]
    public fun test_borrow() {
        let (mut scenario, mut ref, obj) = set_up_test();
        let ctx = scenario.ctx();

        state_object::add(&mut ref, b"test", obj, ctx);
        assert!(state_object::contains(&ref, b"test"));

        let _obj2: &TestObject = state_object::borrow(&ref, b"test");
        assert!(state_object::contains(&ref, b"test"));

        tear_down_test(scenario, ref);
    }

    #[test]
    public fun test_borrow_mut() {
        let (mut scenario, mut ref, obj) = set_up_test();
        let ctx = scenario.ctx();

        state_object::add(&mut ref, b"test", obj, ctx);
        assert!(state_object::contains(&ref, b"test"));

        let _obj2: &mut TestObject = state_object::borrow_mut_with_ctx(&mut ref, b"test", ctx);
        assert!(state_object::contains(&ref, b"test"));

        tear_down_test(scenario, ref);
    }

    #[test]
    public fun test_transfer_ownership() {
        let (mut scenario, mut ref, obj) = set_up_test();
        let ctx = scenario.ctx();

        state_object::add(&mut ref, b"test", obj, ctx);

        let ctx = scenario.ctx();
        let new_owner = SENDER_2;
        state_object::transfer_ownership(&mut ref, new_owner, ctx);

        let (from, to, accepted) = state_object::pending_transfer(&ref);
        assert!(from == SENDER_1);
        assert!(to == new_owner);
        assert!(!accepted);

        // after transfer, the owner is still the original owner
        let owner = state_object::get_current_owner(&ref);
        assert!(owner == SENDER_1);

        tear_down_test(scenario, ref);
    }

    #[test]
    public fun test_accept_and_execute_ownership() {
        let (mut scenario_1, mut ref, obj) = set_up_test();
        let ctx_1 = scenario_1.ctx();
        state_object::add(&mut ref, b"test", obj, ctx_1);

        // tx 1: SENDER_1 transfer ownership to SENDER_2
        // let ctx_1 = scenario_1.ctx();
        let new_owner = SENDER_2;
        state_object::transfer_ownership(&mut ref, new_owner, ctx_1);
        let (from, to, accepted) = state_object::pending_transfer(&ref);
        assert!(from == SENDER_1);
        assert!(to == new_owner);
        assert!(!accepted);

        test_scenario::end(scenario_1);

        // tx 2: SENDER_2 accepts the ownership transfer
        let mut scenario_2 = test_scenario::begin(new_owner);
        // let accept_cap = test_scenario::take_from_address<state_object::AcceptCap>(&scenario_2, new_owner);
        let ctx_2 = scenario_2.ctx();

        state_object::accept_ownership(&mut ref, ctx_2);
        let (from, to, accepted) = state_object::pending_transfer(&ref);
        assert!(from == SENDER_1);
        assert!(to == new_owner);
        assert!(accepted);
        // after accept, the owner is still the original owner
        let owner_1 = state_object::get_current_owner(&ref);
        assert!(owner_1 == SENDER_1);

        test_scenario::end(scenario_2);

        // tx 3: SENDER_1 executes the ownership transfer
        let mut scenario_3 = test_scenario::begin(SENDER_1);
        let ctx_3 = scenario_3.ctx();
        state_object::execute_ownership_transfer(&mut ref, new_owner, ctx_3);
        test_scenario::end(scenario_3);

        let (from, to, accepted) = state_object::pending_transfer(&ref);
        assert!(from == @0x0);
        assert!(to == @0x0);
        assert!(!accepted);
        // after execute, the owner is the new owner
        let owner_2 = state_object::get_current_owner(&ref);
        assert!(owner_2 == SENDER_2);

        // tx 4: SENDER_2 can now update the state object
        let mut scenario_4 = test_scenario::begin(SENDER_2);

        let obj2: TestObject = state_object::remove(&mut ref, b"test", scenario_4.ctx());
        assert!(!state_object::contains(&ref, b"test"));
        let TestObject { id } = obj2;
        object::delete(id);

        tear_down_test(scenario_4, ref);
    }
}