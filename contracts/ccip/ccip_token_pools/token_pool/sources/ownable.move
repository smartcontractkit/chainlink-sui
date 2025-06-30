/// Ownable functionality for the Burn Mint Token Pool module
/// Provides ownership management with two-step ownership transfer process
module ccip_token_pool::ownable {
    use sui::event;

    const EInvalidOwnerCap: u64 = 0;
    const EUnauthorizedOwnershipTransfer: u64 = 1;
    const ENoPendingTransfer: u64 = 2;
    const EUnauthorizedAcceptance: u64 = 3;
    const ETransferAlreadyAccepted: u64 = 4;
    const EInvalidRecipient: u64 = 5;
    const ETransferNotAccepted: u64 = 6;

    public struct OwnerCap has key, store {
        id: UID,
    }

    public struct OwnableState has key, store {
        id: UID,
        owner: address,
        pending_transfer: Option<PendingTransfer>,
        owner_cap_id: ID,
    }

    public struct PendingTransfer has drop, store {
        from: address,
        to: address,
        accepted: bool,
    }

    // =================== Events =================== //

    public struct NewOwnableStateEvent has copy, drop, store {
        ownable_state_id: ID,
        owner_cap_id: ID,
        owner: address,
    }

    public struct OwnershipTransferRequested has copy, drop, store {
        from: address,
        to: address,
    }

    public struct OwnershipTransferAccepted has copy, drop, store {
        from: address,
        to: address,
    }

    public struct OwnershipTransferred has copy, drop, store {
        from: address,
        to: address,
    }

    // =================== Functions =================== //

    public fun new(ctx: &mut TxContext): (OwnableState, OwnerCap) {
        let owner = ctx.sender();
        let owner_cap = OwnerCap {
            id: object::new(ctx),
        };

        let ownable_state = OwnableState {
            id: object::new(ctx),
            owner,
            pending_transfer: option::none(),
            owner_cap_id: object::id(&owner_cap),
        };

        event::emit(NewOwnableStateEvent {
            ownable_state_id: object::id(&ownable_state),
            owner_cap_id: object::id(&owner_cap),
            owner,
        });

        (ownable_state, owner_cap)
    }

    public fun owner(ownable_state: &OwnableState): address {
        ownable_state.owner
    }

    public fun owner_cap_id(ownable_state: &OwnableState): ID {
        ownable_state.owner_cap_id
    }

    public fun has_pending_transfer(ownable_state: &OwnableState): bool {
        ownable_state.pending_transfer.is_some()
    }

    public fun pending_transfer_from(ownable_state: &OwnableState): Option<address> {
        if (ownable_state.pending_transfer.is_some()) {
            option::some(ownable_state.pending_transfer.borrow().from)
        } else {
            option::none()
        }
    }

    public fun pending_transfer_to(ownable_state: &OwnableState): Option<address> {
        if (ownable_state.pending_transfer.is_some()) {
            option::some(ownable_state.pending_transfer.borrow().to)
        } else {
            option::none()
        }
    }

    public fun pending_transfer_accepted(ownable_state: &OwnableState): Option<bool> {
        if (ownable_state.pending_transfer.is_some()) {
            option::some(ownable_state.pending_transfer.borrow().accepted)
        } else {
            option::none()
        }
    }

    public fun transfer_ownership(
        owner_cap: &OwnerCap,
        ownable_state: &mut OwnableState,
        new_owner: address,
        ctx: &mut TxContext,
    ) {
        assert!(object::id(owner_cap) == ownable_state.owner_cap_id, EInvalidOwnerCap);
        assert!(ctx.sender() == ownable_state.owner, EUnauthorizedOwnershipTransfer);

        ownable_state.pending_transfer = option::some(PendingTransfer {
            from: ownable_state.owner,
            to: new_owner,
            accepted: false,
        });

        event::emit(OwnershipTransferRequested {
            from: ownable_state.owner,
            to: new_owner,
        });
    }

    public fun accept_ownership(
        ownable_state: &mut OwnableState,
        ctx: &mut TxContext,
    ) {
        assert!(ownable_state.pending_transfer.is_some(), ENoPendingTransfer);
        let pending = ownable_state.pending_transfer.borrow_mut();
        assert!(ctx.sender() == pending.to, EUnauthorizedAcceptance);
        assert!(!pending.accepted, ETransferAlreadyAccepted);

        pending.accepted = true;

        event::emit(OwnershipTransferAccepted {
            from: pending.from,
            to: pending.to,
        });
    }

    public fun accept_ownership_from_object(
        ownable_state: &mut OwnableState,
        from: &mut UID,
        _ctx: &mut TxContext,
    ) {
        assert!(ownable_state.pending_transfer.is_some(), ENoPendingTransfer);
        let pending = ownable_state.pending_transfer.borrow_mut();
        assert!(from.to_address() == pending.to, EUnauthorizedAcceptance);
        assert!(!pending.accepted, ETransferAlreadyAccepted);

        pending.accepted = true;

        event::emit(OwnershipTransferAccepted {
            from: pending.from,
            to: pending.to,
        });
    }

    public fun accept_ownership_as_mcms(
        ownable_state: &mut OwnableState,
        mcms: address,
        _ctx: &mut TxContext,
    ) {
        assert!(ownable_state.pending_transfer.is_some(), ENoPendingTransfer);
        let pending = ownable_state.pending_transfer.borrow_mut();
        assert!(mcms == pending.to, EUnauthorizedAcceptance);
        assert!(!pending.accepted, ETransferAlreadyAccepted);

        pending.accepted = true;

        event::emit(OwnershipTransferAccepted {
            from: pending.from,
            to: pending.to,
        });
    }

    public fun execute_ownership_transfer(
        owner_cap: OwnerCap,
        ownable_state: &mut OwnableState,
        to: address,
        ctx: &mut TxContext,
    ) {
        assert!(object::id(&owner_cap) == ownable_state.owner_cap_id, EInvalidOwnerCap);
        assert!(ownable_state.pending_transfer.is_some(), ENoPendingTransfer);
        let pending = ownable_state.pending_transfer.borrow();
        assert!(pending.to == to, EInvalidRecipient);
        assert!(pending.accepted, ETransferNotAccepted);

        let old_owner = ownable_state.owner;
        ownable_state.owner = to;
        ownable_state.pending_transfer = option::none();

        // Create new owner cap for the new owner
        let new_owner_cap = OwnerCap {
            id: object::new(ctx),
        };
        ownable_state.owner_cap_id = object::id(&new_owner_cap);

        // Destroy the old owner cap
        let OwnerCap { id } = owner_cap;
        object::delete(id);

        event::emit(OwnershipTransferred {
            from: old_owner,
            to,
        });

        // Transfer the new owner cap to the new owner
        transfer::public_transfer(new_owner_cap, to);
    }

    public fun set_owner(
        owner_cap: &OwnerCap,
        ownable_state: &mut OwnableState,
        new_owner: address,
        _ctx: &mut TxContext,
    ) {
        assert!(object::id(owner_cap) == ownable_state.owner_cap_id, EInvalidOwnerCap);
        
        let old_owner = ownable_state.owner;
        ownable_state.owner = new_owner;

        event::emit(OwnershipTransferred {
            from: old_owner,
            to: new_owner,
        });
    }

    public fun destroy_ownable_state(ownable_state: OwnableState) {
        let OwnableState {
            id: ownable_id,
            owner: _,
            pending_transfer: _,
            owner_cap_id: _,
        } = ownable_state;
        object::delete(ownable_id);
    }

    public fun destroy_owner_cap(owner_cap: OwnerCap) {
        let OwnerCap { id: owner_cap_id } = owner_cap;
        object::delete(owner_cap_id);
    }
} 