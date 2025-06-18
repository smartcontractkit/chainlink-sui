module ccip_onramp::ownable {
    use mcms::mcms_registry::{Self, Registry};
    use sui::event;

    public struct OwnerCap has key, store {
        id: UID,
        ownable_state_id: ID,
    }

    public struct OwnableState has key, store {
        id: UID,
        owner: address,
        pending_transfer: Option<PendingTransfer>,
    }

    public struct PendingTransfer has drop, store {
        from: address,
        to: address,
        accepted: bool,
    }

    public struct McmsProof has drop {}

    // =================== Events =================== //

    public struct NewOwnableStateEvent has copy, store, drop {
        ownable_state_id: ID,
        owner_cap_id: ID,
        owner: address,
    }

    public struct OwnershipTransferRequested has copy, drop {
        from: address,
        to: address,
    }

    public struct OwnershipTransferAccepted has copy, drop {
        from: address,
        to: address,
    }

    public struct OwnershipTransferred has copy, drop {
        from: address,
        to: address,
    }

    const EInvalidOwnerCap: u64 = 1;
    const ECannotTransferToSelf: u64 = 2;
    const EMustBeProposedOwner: u64 = 3;
    const ENoPendingTransfer: u64 = 4;
    const ETransferAlreadyAccepted: u64 = 5;
    const EOwnerChanged: u64 = 6;
    const EProposedOwnerMismatch: u64 = 7;
    const ETransferNotAccepted: u64 = 8;

    public fun new(ctx: &mut TxContext): (OwnableState, OwnerCap) {
        let owner = ctx.sender();

        let state = OwnableState {
            id: object::new(ctx),
            owner,
            pending_transfer: option::none(),
        };

        let owner_cap = OwnerCap {
            id: object::new(ctx),
            ownable_state_id: object::id(&state),
        };

        event::emit(NewOwnableStateEvent {
            ownable_state_id: object::id(&state),
            owner_cap_id: object::id(&owner_cap),
            owner,
        });

        (state, owner_cap)
    }

    public fun set_owner(
        owner_cap: &OwnerCap,
        state: &mut OwnableState,
        owner: address,
        _ctx: &mut TxContext,
    ) {
        assert!(owner_cap.ownable_state_id == object::id(state), EInvalidOwnerCap);
        state.owner = owner;
    }

    public fun transfer_ownership(
        owner_cap: &OwnerCap,
        state: &mut OwnableState,
        to: address,
        _ctx: &mut TxContext,
    ) {
        assert!(owner_cap.ownable_state_id == object::id(state), EInvalidOwnerCap);
        assert!(state.owner != to, ECannotTransferToSelf);

        state.pending_transfer =
            option::some(PendingTransfer {
                from: state.owner,
                to,
                accepted: false,
            });

        event::emit(OwnershipTransferRequested { from: state.owner, to });
    }

    public fun accept_ownership(state: &mut OwnableState, ctx: &mut TxContext) {
        accept_ownership_internal(state, ctx.sender());
    }

    /// UID is a privileged type that is only accessible by the object owner.
    public fun accept_ownership_from_object(state: &mut OwnableState, from: &mut UID) {
        accept_ownership_internal(state, from.to_address());
    }

    public(package) fun accept_ownership_as_mcms(state: &mut OwnableState, _ctx: &mut TxContext) {
        accept_ownership_internal(state, @mcms);
    }

    fun accept_ownership_internal(state: &mut OwnableState, caller: address) {
        assert!(state.pending_transfer.is_some(), ENoPendingTransfer);

        let pending_transfer = state.pending_transfer.borrow_mut();
        let current_owner = state.owner;

        // check that the owner has not changed from a direct call to 0x1::transfer::public_transfer,
        // in which case the transfer flow should be restarted.
        assert!(current_owner == pending_transfer.from, EOwnerChanged);
        assert!(caller == pending_transfer.to, EMustBeProposedOwner);
        assert!(!pending_transfer.accepted, ETransferAlreadyAccepted);

        pending_transfer.accepted = true;

        event::emit(OwnershipTransferAccepted { from: pending_transfer.from, to: caller });
    }

    #[allow(lint(custom_state_change))]
    public fun execute_ownership_transfer(
        owner_cap: OwnerCap,
        state: &mut OwnableState,
        registry: &mut Registry,
        to: address,
        ctx: &mut TxContext,
    ) {
        assert!(owner_cap.ownable_state_id == object::id(state), EInvalidOwnerCap);
        assert!(state.pending_transfer.is_some(), ENoPendingTransfer);

        let pending_transfer = state.pending_transfer.extract();
        let current_owner = state.owner;
        let new_owner = pending_transfer.to;

        // check that the owner has not changed from a direct call to 0x1::transfer::public_transfer,
        // in which case the transfer flow should be restarted.
        assert!(pending_transfer.from == current_owner, EOwnerChanged);
        assert!(new_owner == to, EProposedOwnerMismatch);
        assert!(pending_transfer.accepted, ETransferNotAccepted);

        // if the new owner is mcms, we need to add the `OwnerCap` to the registry.
        if (new_owner == @mcms) {
            mcms_registry::register_entrypoint(
                registry,
                McmsProof {},
                option::some(owner_cap),
                ctx,
            );
        } else {
            transfer::transfer(owner_cap, new_owner);
        };

        state.owner = new_owner;
        state.pending_transfer = option::none();

        event::emit(OwnershipTransferred { from: current_owner, to: new_owner });
    }
}
