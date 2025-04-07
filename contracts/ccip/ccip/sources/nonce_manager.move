module ccip::nonce_manager {
    use std::string::{Self, String};
    use sui::table::{Self, Table};

    use ccip::state_object::{Self, OwnerCap, UserCap, CCIPObjectRef};

    public struct NonceManagerState has key, store {
        id: UID,
        // dest chain selector -> sender -> nonce
        outbound_nonces: Table<u64, Table<address, u64>>
    }

    const NONCE_MANAGER_STATE_NAME: vector<u8> = b"NonceManagerState";
    const E_ALREADY_INITIALIZED: u64 = 1;

    public fun type_and_version(): String {
        string::utf8(b"NonceManager 1.6.0")
    }

    public fun initialize(owner_cap: &OwnerCap, ref: &mut CCIPObjectRef, ctx: &mut TxContext) {
        assert!(
            !state_object::contains(ref, NONCE_MANAGER_STATE_NAME),
            E_ALREADY_INITIALIZED
        );

        let state = NonceManagerState {
            id: object::new(ctx),
            outbound_nonces: table::new(ctx)
        };
        state_object::add(owner_cap, ref, NONCE_MANAGER_STATE_NAME, state);
    }

    public fun get_outbound_nonce(
        ref: &CCIPObjectRef,
        dest_chain_selector: u64,
        sender: address
    ): u64 {
        let state = state_object::borrow<NonceManagerState>(ref, NONCE_MANAGER_STATE_NAME);

        if (!table::contains(&state.outbound_nonces, dest_chain_selector)) {
            return 0
        };

        let dest_chain_nonces = &state.outbound_nonces[dest_chain_selector];
        if (!table::contains(dest_chain_nonces, sender)) {
            return 0
        };
        dest_chain_nonces[sender]
    }

    public(package) fun get_incremented_outbound_nonce(
        user_cap: &UserCap,
        ref: &mut CCIPObjectRef,
        dest_chain_selector: u64,
        sender: address,
        ctx: &mut TxContext
    ): u64 {
        let state = state_object::borrow_mut_from_user<NonceManagerState>(user_cap, ref, NONCE_MANAGER_STATE_NAME);

        if (!table::contains(&state.outbound_nonces, dest_chain_selector)) {
            table::add(
                &mut state.outbound_nonces, dest_chain_selector, table::new(ctx)
            );
        };
        let dest_chain_nonces =
            table::borrow_mut(&mut state.outbound_nonces, dest_chain_selector);
        if (!table::contains(dest_chain_nonces, sender)) {
            table::add(dest_chain_nonces, sender, 0);
        };

        let nonce_ref = table::borrow_mut(dest_chain_nonces, sender);
        let incremented_nonce = *nonce_ref + 1;
        *nonce_ref = incremented_nonce;
        incremented_nonce
    }
}

#[test_only]
module ccip::nonce_manager_test {
    use ccip::nonce_manager;
    use ccip::state_object::{Self, OwnerCap, UserCap, CCIPObjectRef};
    use sui::test_scenario::{Self, Scenario};

    const NONCE_MANAGER_STATE_NAME: vector<u8> = b"NonceManagerState";

    fun set_up_test(): (Scenario, OwnerCap, UserCap, CCIPObjectRef) {
        let mut scenario = test_scenario::begin(@0x1);
        let ctx = scenario.ctx();

        let (owner_cap, user_cap, ref) = state_object::create(ctx);

        (scenario, owner_cap, user_cap, ref)
    }

    fun initialize(owner_cap: &OwnerCap, ref: &mut CCIPObjectRef, ctx: &mut TxContext) {
        nonce_manager::initialize(
            owner_cap,
            ref,
            ctx
        );
    }

    fun tear_down_test(scenario: Scenario, owner_cap: OwnerCap, user_cap: UserCap, ref: CCIPObjectRef) {
        state_object::destroy_owner_cap(owner_cap);
        state_object::destroy_user_cap(user_cap);
        state_object::destroy_state_object(ref);
        test_scenario::end(scenario);
    }

    #[test]
    public fun test_initialize() {
        let (mut scenario, owner_cap, user_cap, mut ref) = set_up_test();
        let ctx = scenario.ctx();
        initialize(&owner_cap, &mut ref, ctx);

        let _state = state_object::borrow<nonce_manager::NonceManagerState>(&ref, NONCE_MANAGER_STATE_NAME);

        assert!(
            state_object::contains(
                &ref,
                NONCE_MANAGER_STATE_NAME
            )
        );

        tear_down_test(scenario, owner_cap, user_cap, ref);
    }

    #[test]
    public fun test_get_incremented_outbound_nonce() {
        let (mut scenario, owner_cap, user_cap, mut ref) = set_up_test();
        let ctx = scenario.ctx();
        initialize(&owner_cap, &mut ref, ctx);

        let mut nonce = nonce_manager::get_outbound_nonce(&ref, 1, @0x1);
        assert!(nonce == 0);

        let mut incremented_nonce = nonce_manager::get_incremented_outbound_nonce(
            &user_cap,
            &mut ref,
            1,
            @0x1,
            ctx
        );
        assert!(incremented_nonce == 1);

        nonce = nonce_manager::get_outbound_nonce(&ref, 1, @0x1);
        assert!(nonce == 1);

        incremented_nonce = nonce_manager::get_incremented_outbound_nonce(
            &user_cap,
            &mut ref,
            1,
            @0x1,
            ctx
        );
        assert!(incremented_nonce == 2);

        nonce = nonce_manager::get_outbound_nonce(&ref, 1, @0x1);
        assert!(nonce == 2);

        tear_down_test(scenario, owner_cap, user_cap, ref);
    }
}