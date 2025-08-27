#[test_only]
module mcms::mcms_account_test;

use mcms::mcms_account::{Self, AccountState, OwnerCap};
use mcms::mcms_registry::{Self, Registry};
use sui::test_scenario::{Self as ts};

const OWNER: address = @0x123;

public struct Env {
    scenario: ts::Scenario,
    state: AccountState,
    registry: Registry,
}

public fun setup(): Env {
    let mut scenario = ts::begin(OWNER);
    let ctx = scenario.ctx();

    mcms_account::test_init(ctx);
    mcms_registry::test_init(ctx);

    scenario.next_tx(OWNER);

    let registry = ts::take_shared<Registry>(&scenario);
    let state = ts::take_shared<AccountState>(&scenario);

    Env { scenario, state, registry }
}

#[test]
fun test_transfer_ownership_to_self_flow() {
    let mut env = setup();
    let owner_cap = ts::take_from_sender<OwnerCap>(&env.scenario);

    mcms_account::transfer_ownership_to_self(
        &owner_cap,
        &mut env.state,
        env.scenario.ctx(),
    );
    assert!(mcms_account::pending_transfer_from(&env.state) == option::some(OWNER));
    assert!(mcms_account::pending_transfer_to(&env.state) == option::some(mcms_registry::get_multisig_address()));
    assert!(mcms_account::pending_transfer_accepted(&env.state) == option::some(false));

    mcms_account::test_accept_ownership_as_timelock(
        &mut env.state,
        env.scenario.ctx(),
    );
    assert!(mcms_account::pending_transfer_accepted(&env.state) == option::some(true));

    mcms_account::execute_ownership_transfer(
        owner_cap,
        &mut env.state,
        &mut env.registry,
        mcms_registry::get_multisig_address(),
        env.scenario.ctx(),
    );
    assert!(mcms_account::pending_transfer_from(&env.state) == option::none());
    assert!(mcms_account::pending_transfer_to(&env.state) == option::none());
    assert!(mcms_account::pending_transfer_accepted(&env.state) == option::none());

    assert!(mcms_registry::is_package_registered(&env.registry, mcms_registry::get_multisig_address()));

    env.destroy();
}

public fun destroy(env: Env) {
    let Env {
        scenario,
        state,
        registry,
    } = env;

    ts::return_shared(registry);
    ts::return_shared(state);

    scenario.end();
}
