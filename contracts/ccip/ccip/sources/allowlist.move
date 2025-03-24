module ccip::allowlist {
    use sui::event;

    public struct AllowlistState has key, store {
        id: UID,
        allowlist_enabled: bool,
        // it's possible to use vec_set for allowlist
        allowlist: vector<address>,
    }

    public struct AllowlistRemove has copy, drop {
        sender: address,
    }

    public struct AllowlistAdd has copy, drop {
        sender: address,
    }

    const E_ALLOWLIST_NOT_ENABLED: u64 = 1;

    public fun new(allowlist: vector<address>, ctx: &mut TxContext): AllowlistState {
        AllowlistState {
            id: object::new(ctx),
            allowlist_enabled: !vector::is_empty(&allowlist),
            allowlist,
        }
    }

    public fun get_allowlist_enabled(state: &AllowlistState): bool {
        state.allowlist_enabled
    }

    public fun set_allowlist_enabled(state: &mut AllowlistState, enabled: bool) {
        state.allowlist_enabled = enabled
    }

    public fun get_allowlist(state: &AllowlistState): vector<address> {
        state.allowlist
    }

    public fun is_allowed(state: &AllowlistState, sender: address): bool {
        if (!state.allowlist_enabled) {
            return true
        };
        vector::contains(&state.allowlist, &sender)
    }

    public fun apply_allowlist_updates(
        state: &mut AllowlistState,
        removes: vector<address>,
        adds: vector<address>,
    ) {
        let mut i = 0;
        let mut len = vector::length(&removes);
        while (i < len) {
            // get the address to remove
            let remove_address = vector::borrow(&removes, i);
            // find the index of the address in the allowlist
            let (found, j) = vector::index_of(&state.allowlist, remove_address);
            if (found) {
                vector::swap_remove(&mut state.allowlist, j);
                event::emit(AllowlistRemove {
                    sender: *remove_address,
                });
            };
            i = i + 1;
        };

        if (!vector::is_empty(&adds)) {
            assert!(state.allowlist_enabled, E_ALLOWLIST_NOT_ENABLED);

            i = 0;
            len = vector::length(&adds);
            while (i < len) {
                let add_address = vector::borrow(&adds, i);
                let (found, _) = vector::index_of(&state.allowlist, add_address);
                if (add_address != @0x0 && !found) {
                    vector::push_back(&mut state.allowlist, *add_address);
                    event::emit(AllowlistAdd {
                        sender: *add_address,
                    });
                };
                i = i + 1;
            };
        }
    }

    public fun destroy_allowlist(state: AllowlistState) {
        let AllowlistState {
            id,
            allowlist_enabled: _,
            allowlist: _,
        } = state;
        object::delete(id);
    }
}

#[test_only]
module ccip::allowlist_test {
    use ccip::allowlist;
    use sui::test_scenario;

    fun set_up_test(allowlist: vector<address>, ctx: &mut TxContext): allowlist::AllowlistState {
        allowlist::new(allowlist, ctx)
    }

    #[test]
    public fun init_empty_is_empty_and_disabled() {
        let mut scenario = test_scenario::begin(@0x1);

        let state = set_up_test(vector::empty(), scenario.ctx());

        assert!(!allowlist::get_allowlist_enabled(&state), 1);
        assert!(vector::is_empty(&allowlist::get_allowlist(&state)), 1);

        // Any address is allowed when the allowlist is disabled
        assert!(allowlist::is_allowed(&state, @0x1111111111111), 1);

        allowlist::destroy_allowlist(state);

        test_scenario::end(scenario);
    }

    #[test]
    public fun init_non_empty_is_non_empty_and_enabled() {
        let mut scenario = test_scenario::begin(@0x1);

        let init_allowlist = vector[@0x1, @0x2];

        let state = set_up_test(init_allowlist, scenario.ctx());

        assert!(allowlist::get_allowlist_enabled(&state), 1);
        assert!(vector::length(&allowlist::get_allowlist(&state)) == 2, 1);

        // The given addresses are allowed
        assert!(allowlist::is_allowed(&state, *vector::borrow(&init_allowlist, 0)), 1);
        assert!(allowlist::is_allowed(&state, *vector::borrow(&init_allowlist, 1)), 1);

        assert!(!allowlist::is_allowed(&state, @0x3), 1);

        allowlist::destroy_allowlist(state);

        test_scenario::end(scenario);
    }

    #[test]
    #[expected_failure(abort_code = allowlist::E_ALLOWLIST_NOT_ENABLED, location = allowlist)]
    public fun cannot_add_to_disabled_allowlist() {
        let mut scenario = test_scenario::begin(@0x1);

        let mut state = set_up_test(vector::empty(), scenario.ctx());

        let adds = vector[@0x1];

        allowlist::apply_allowlist_updates(&mut state, vector::empty(), adds);

        allowlist::destroy_allowlist(state);

        test_scenario::end(scenario);
    }

    #[test]
    public fun apply_allowlist_updates_mutates_state() {
        let mut scenario = test_scenario::begin(@0x1);

        let mut state = set_up_test(vector::empty(), scenario.ctx());

        allowlist::set_allowlist_enabled(&mut state, true);

        assert!(vector::is_empty(&allowlist::get_allowlist(&state)), 1);

        allowlist::apply_allowlist_updates(&mut state, vector::empty(), vector::empty());

        assert!(vector::is_empty(&allowlist::get_allowlist(&state)), 1);

        let adds = vector[@0x1, @0x2];

        allowlist::apply_allowlist_updates(&mut state, vector::empty(), adds);

        let removes = vector[@0x1];

        allowlist::apply_allowlist_updates(&mut state, removes, vector::empty());

        assert!(vector::length(&allowlist::get_allowlist(&state)) == 1, 1);
        assert!(allowlist::is_allowed(&state, @0x2), 1);
        assert!(!allowlist::is_allowed(&state, @0x1), 1);

        allowlist::destroy_allowlist(state);

        test_scenario::end(scenario);
    }

    #[test]
    public fun apply_allowlist_updates_removes_before_adds() {
        let mut scenario = test_scenario::begin(@0x1);

        let mut state = set_up_test(vector::empty(), scenario.ctx());
        let account_to_allow = @0x1;

        allowlist::set_allowlist_enabled(&mut state, true);

        let adds_and_removes = vector[account_to_allow];

        allowlist::apply_allowlist_updates(&mut state, vector::empty(), adds_and_removes);

        assert!(vector::length(&allowlist::get_allowlist(&state)) == 1, 1);
        assert!(allowlist::is_allowed(&state, account_to_allow), 1);

        allowlist::apply_allowlist_updates(&mut state, adds_and_removes, adds_and_removes);

        // Since removes happen before adds, the account should still be allowed
        assert!(allowlist::is_allowed(&state, account_to_allow), 1);

        allowlist::destroy_allowlist(state);

        test_scenario::end(scenario);
    }
}
