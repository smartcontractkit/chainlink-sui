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

    // currently we only create 1 capability to manage all the CCIP components like fee_quoter, rmn_remote
    // TODO: figure out if we need to create multiple capabilities for each CCIP component
    public struct OwnerCap has key, store {
        id: UID
    }

    public struct UserCap has key, store {
        id: UID,
    }

    public struct AcceptCap has key, store {
        id: UID
    }

    public struct CCIPObjectRef has key, store {
        id: UID,
        // this is the owner_cap's owner
        // this object is a shared object and cannot be transferred
        // the ownership of the entire CCIP object ref is equivalent to the ownership of its owner cap object
        current_owner_cap_owner: address,
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
            current_owner_cap_owner: ctx.sender(),
            pending_transfer: option::none()
        };

        let owner_cap = OwnerCap {
            id: object::new(ctx)
        };

        let user_cap = UserCap {
            id: object::new(ctx)
        };

        // TODO: another approach is to transfer to owner.
        // however, in that case, any other addresses cannot access this object onchain
        // this includes the new proposed owner
        transfer::share_object(ref);
        transfer::transfer(owner_cap, ctx.sender());
        transfer::freeze_object(user_cap);
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

    public(package) fun borrow_mut_from_user<T: key + store>(_: &UserCap, ref: &mut CCIPObjectRef, name: vector<u8>): &mut T {
        dof::borrow_mut(&mut ref.id, name)
    }

    public fun transfer_ownership(
        _: &OwnerCap, ref: &mut CCIPObjectRef, to: address, ctx: &mut TxContext
    ) {
        let caller = ctx.sender();
        assert!(caller != to, E_CANNOT_TRANSFER_TO_SELF);

        ref.pending_transfer = option::some(
            PendingTransfer { from: caller, to, accepted: false }
        );

        // create and transfer the accept cap to the proposed owner
        let accept_cap = AcceptCap {
            id: object::new(ctx)
        };
        transfer::transfer(accept_cap, to);

        event::emit(OwnershipTransferRequested { from: caller, to });
    }

    // if CCIPObjectRef is not a shared object, the proposed new owner cannot accept this object onchain
    public fun accept_ownership(accept_cap: AcceptCap, ref: &mut CCIPObjectRef, ctx: &mut TxContext) {
        assert!(
            ref.pending_transfer.is_some(),
            E_NO_PENDING_TRANSFER
        );

        let caller = ctx.sender();
        let pending_transfer = ref.pending_transfer.borrow_mut();

        assert!(
            pending_transfer.from == ref.current_owner_cap_owner,
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

        // destroy the accept cap once the proposed owner has accepted the ownership transfer
        let AcceptCap { id } = accept_cap;
        object::delete(id);

        event::emit(OwnershipTransferAccepted { from: pending_transfer.from, to: caller });
    }

    // the only thing needs to be transfered is the owner cap
    public fun execute_ownership_transfer(
        owner_cap: OwnerCap, ref: &mut CCIPObjectRef, to: address, ctx: &mut TxContext
    ) {
        let caller_address = ctx.sender();

        let pending_transfer = ref.pending_transfer.extract();

        // since ref is a shared object now, it's impossible to transfer its ownership
        assert!(
            pending_transfer.from == ref.current_owner_cap_owner,
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
        ref.current_owner_cap_owner = pending_transfer.to;
        transfer::public_transfer(owner_cap, pending_transfer.to);
        // the extract will remove the object within option wrapper
        // state.pending_transfer = option::none();

        event::emit(OwnershipTransferred { from: caller_address, to })
    }

    #[test_only]
    public(package) fun create(ctx: &mut TxContext): (OwnerCap, UserCap, CCIPObjectRef) {
        let ref = CCIPObjectRef {
            id: object::new(ctx),
            current_owner_cap_owner: ctx.sender(),
            pending_transfer: option::none()
        };

        let owner_cap = OwnerCap {
            id: object::new(ctx)
        };

        let user_cap = UserCap {
            id: object::new(ctx)
        };

        (owner_cap, user_cap, ref)
    }

    #[test_only]
    public(package) fun get_current_owner_cap_owner(ref: &CCIPObjectRef): address {
        ref.current_owner_cap_owner
    }

    #[test_only]
    public(package) fun destroy_owner_cap(cap: OwnerCap) {
        let OwnerCap { id } = cap;
        object::delete(id);
    }

    #[test_only]
    public(package) fun destroy_user_cap(cap: UserCap) {
        let UserCap { id } = cap;
        object::delete(id);
    }

    #[test_only]
    public(package) fun destroy_state_object(ref: CCIPObjectRef) {
        let CCIPObjectRef { id, current_owner_cap_owner: _owner, pending_transfer: _pending_transfer } = ref;
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
    use ccip::state_object::{Self, OwnerCap, UserCap, CCIPObjectRef};
    use sui::test_scenario::{Self, Scenario};

    const SENDER_1: address = @0x1;
    const SENDER_2: address = @0x2;

    fun set_up_test(): (Scenario, OwnerCap, UserCap, CCIPObjectRef, TestObject) {
        let mut scenario = test_scenario::begin(@0x1);
        let ctx = scenario.ctx();

        let (owner_cap, user_cap, ref) = state_object::create(ctx);
        let obj = TestObject {
            id: object::new(ctx)
        };
        (scenario, owner_cap, user_cap, ref, obj)
    }

    fun tear_down_test(scenario: Scenario, owner_cap: OwnerCap, user_cap: UserCap, ref: CCIPObjectRef) {
        state_object::destroy_owner_cap(owner_cap);
        state_object::destroy_user_cap(user_cap);
        state_object::destroy_state_object(ref);
        test_scenario::end(scenario);
    }

    public struct TestObject has key, store {
        id: UID
    }

    #[test]
    public fun test_add() {
        let (scenario, owner_cap, user_cap, mut ref, obj) = set_up_test();

        state_object::add(&owner_cap, &mut ref, b"test", obj);
        assert!(state_object::contains(&ref, b"test"));

        tear_down_test(scenario, owner_cap, user_cap, ref);
    }

    #[test]
    public fun test_remove() {
        let (scenario, owner_cap, user_cap, mut ref, obj) = set_up_test();

        state_object::add(&owner_cap, &mut ref, b"test", obj);
        assert!(state_object::contains(&ref, b"test"));

        let obj2: TestObject = state_object::remove(&owner_cap, &mut ref, b"test");
        assert!(!state_object::contains(&ref, b"test"));

        let TestObject { id } = obj2;
        object::delete(id);

        tear_down_test(scenario, owner_cap, user_cap, ref);
    }

    #[test]
    public fun test_borrow() {
        let (scenario, owner_cap, user_cap, mut ref, obj) = set_up_test();

        state_object::add(&owner_cap, &mut ref, b"test", obj);
        assert!(state_object::contains(&ref, b"test"));

        let _obj2: &TestObject = state_object::borrow(&ref, b"test");
        assert!(state_object::contains(&ref, b"test"));

        tear_down_test(scenario, owner_cap, user_cap, ref);
    }

    #[test]
    public fun test_borrow_mut() {
        let (scenario, owner_cap, user_cap, mut ref, obj) = set_up_test();

        state_object::add(&owner_cap, &mut ref, b"test", obj);
        assert!(state_object::contains(&ref, b"test"));

        let _obj2: &mut TestObject = state_object::borrow_mut(&owner_cap, &mut ref, b"test");
        assert!(state_object::contains(&ref, b"test"));

        tear_down_test(scenario, owner_cap, user_cap, ref);
    }

    #[test]
    public fun test_transfer_ownership() {
        let (mut scenario, owner_cap, user_cap, mut ref, obj) = set_up_test();

        state_object::add(&owner_cap, &mut ref, b"test", obj);

        let ctx = scenario.ctx();
        let new_owner = SENDER_2;
        state_object::transfer_ownership(&owner_cap, &mut ref, new_owner, ctx);

        let (from, to, accepted) = state_object::pending_transfer(&ref);
        assert!(from == SENDER_1);
        assert!(to == new_owner);
        assert!(!accepted);

        // after transfer, the owner is still the original owner
        let owner = state_object::get_current_owner_cap_owner(&ref);
        assert!(owner == SENDER_1);

        tear_down_test(scenario, owner_cap, user_cap, ref);
    }

    #[test]
    public fun test_accept_and_execute_ownership() {
        let (mut scenario_1, owner_cap, user_cap, mut ref, obj) = set_up_test();
        state_object::add(&owner_cap, &mut ref, b"test", obj);

        // tx 1: SENDER_1 transfer ownership to SENDER_2
        let ctx_1 = scenario_1.ctx();
        let new_owner = SENDER_2;
        state_object::transfer_ownership(&owner_cap, &mut ref, new_owner, ctx_1);
        let (from, to, accepted) = state_object::pending_transfer(&ref);
        assert!(from == SENDER_1);
        assert!(to == new_owner);
        assert!(!accepted);

        let effects_1 = test_scenario::end(scenario_1);
        let transferred = test_scenario::transferred_to_account(&effects_1);

        assert!(transferred.keys().length() == 1); // the accept_cap should be the only object transferred
        let id = transferred.keys()[0];
        let addr = transferred[&id];
        assert!(addr == new_owner); // the accept_cap should be transferred to the new owner

        // tx 2: SENDER_2 accepts the ownership transfer
        let mut scenario_2 = test_scenario::begin(new_owner);
        let accept_cap = test_scenario::take_from_address<state_object::AcceptCap>(&scenario_2, new_owner);
        let ctx_2 = scenario_2.ctx();

        state_object::accept_ownership(
            accept_cap,
            &mut ref,
            ctx_2
        );
        let (from, to, accepted) = state_object::pending_transfer(&ref);
        assert!(from == SENDER_1);
        assert!(to == new_owner);
        assert!(accepted);
        // after accept, the owner is still the original owner
        let owner_1 = state_object::get_current_owner_cap_owner(&ref);
        assert!(owner_1 == SENDER_1);

        let effects_2 = test_scenario::end(scenario_2);
        let deleted_2 = test_scenario::deleted(&effects_2);
        assert!(deleted_2.length() == 1); // the accept_cap should be deleted

        // tx 3: SENDER_1 executes the ownership transfer
        let mut scenario_3 = test_scenario::begin(SENDER_1);
        let ctx_3 = scenario_3.ctx();
        state_object::execute_ownership_transfer(owner_cap, &mut ref, new_owner, ctx_3);

        let effects_3 = test_scenario::end(scenario_3);
        let (from, to, accepted) = state_object::pending_transfer(&ref);
        assert!(from == @0x0);
        assert!(to == @0x0);
        assert!(!accepted);
        // after execute, the owner is the new owner
        let owner_2 = state_object::get_current_owner_cap_owner(&ref);
        assert!(owner_2 == SENDER_2);

        let transferred_3 = test_scenario::transferred_to_account(&effects_3);
        assert!(transferred_3.keys().length() == 1); // the owner_cap should be the only object transferred
        let owner_cap_id = transferred_3.keys()[0];
        let owner_cap_owner = transferred_3[&owner_cap_id];
        assert!(owner_cap_owner == new_owner); // the accept_cap should be transferred to the new owner

        // tx 4: SENDER_2 can now update the state object
        let scenario_4 = test_scenario::begin(SENDER_2);
        let owner_cap = test_scenario::take_from_address<state_object::OwnerCap>(&scenario_4, SENDER_2);

        // the new owner use the owner_cap to access the state object
        let obj2: TestObject = state_object::remove(&owner_cap, &mut ref, b"test");
        assert!(!state_object::contains(&ref, b"test"));
        let TestObject { id } = obj2;
        object::delete(id);

        tear_down_test(scenario_4, owner_cap, user_cap, ref);
    }
}