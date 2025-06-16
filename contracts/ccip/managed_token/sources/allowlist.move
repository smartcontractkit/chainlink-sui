module managed_token::allowlist {
    use sui::event;
    use std::string::{Self, String};

    public struct AllowlistState has store {
        allowlist_name: String,
        allowlist_enabled: bool,
        allowlist: vector<address>,
    }

    public struct AllowlistRemoved has copy, drop {
        allowlist_name: String,
        sender: address,
    }

    public struct AllowlistAdded has copy, drop {
        allowlist_name: String,
        sender: address,
    }

    const E_ALLOWLIST_NOT_ENABLED: u64 = 1;

    public fun new(allowlist: vector<address>): AllowlistState {
        new_with_name(allowlist, string::utf8(b"default"))
    }

    public fun new_with_name(allowlist: vector<address>, allowlist_name: String
    ): AllowlistState {
        AllowlistState {
            allowlist_name,
            allowlist_enabled: !allowlist.is_empty(),
            allowlist,
        }
    }

    public fun get_allowlist_enabled(state: &AllowlistState): bool {
        state.allowlist_enabled
    }

    public fun set_allowlist_enabled(
        state: &mut AllowlistState, enabled: bool
    ) {
        state.allowlist_enabled = enabled;
    }

    public fun get_allowlist(state: &AllowlistState): vector<address> {
        state.allowlist
    }

    public fun is_allowed(state: &AllowlistState, sender: address): bool {
        if (!state.allowlist_enabled) {
            return true
        };

        state.allowlist.contains(&sender)
    }

    public fun apply_allowlist_updates(
        state: &mut AllowlistState, removes: vector<address>, adds: vector<address>
    ) {
        removes.do_ref!(
            |remove_address| {
                let (found, i) = state.allowlist.index_of(remove_address);
                if (found) {
                    state.allowlist.swap_remove(i);
                    event::emit(
                        AllowlistRemoved {
                            allowlist_name: state.allowlist_name,
                            sender: *remove_address
                        }
                    );
                }
            }
        );

        if (!adds.is_empty()) {
            assert!(
                state.allowlist_enabled,
                E_ALLOWLIST_NOT_ENABLED
            );

            adds.do_ref!(
                |add_address| {
                    let add_address: address = *add_address;
                    let (found, _) = state.allowlist.index_of(&add_address);
                    if (add_address != @0x0 && !found) {
                        state.allowlist.push_back(add_address);
                        event::emit(
                            AllowlistAdded {
                                allowlist_name: state.allowlist_name,
                                sender: add_address
                            }
                        );
                    }
                }
            );
        }
    }

    public fun destroy_allowlist(state: AllowlistState) {
        let AllowlistState {
            allowlist_name: _,
            allowlist_enabled: _,
            allowlist: _,
        } = state;
    }
}
