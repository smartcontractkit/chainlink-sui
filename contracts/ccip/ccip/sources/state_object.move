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
}