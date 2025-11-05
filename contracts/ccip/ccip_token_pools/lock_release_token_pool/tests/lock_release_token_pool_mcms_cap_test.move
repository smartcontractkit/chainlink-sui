#[test_only]
module lock_release_token_pool::lock_release_token_pool_mcms_cap_test;

use ccip::offramp_state_helper;
use ccip::onramp_state_helper;
use ccip::ownable::OwnerCap as CCIPOwnerCap;
use ccip::rmn_remote;
use ccip::state_object::{Self, CCIPObjectRef};
use ccip::token_admin_registry;
use ccip::upgrade_registry;
use lock_release_token_pool::lock_release_token_pool::{
    Self,
    LockReleaseTokenPoolState,
    RebalancerCap,
    McmsCap
};
use lock_release_token_pool::ownable::{Self, OwnerCap};
use mcms::mcms_registry::{Self, Registry, ExecutingCallbackParams};
use std::bcs;
use std::string;
use sui::coin::{Self, Coin};
use sui::test_scenario::{Self as ts, Scenario};

const OWNER: address = @0x123;
const NEW_OWNER: address = @0x456;
const EOA_REBALANCER: address = @0x789;
const EOA_REBALANCER_2: address = @0xabc;
const OTHER_USER: address = @0xdef;
const CCIP_ADMIN: address = @0x400;
const TOKEN_ADMIN: address = @0x200;

public struct LOCK_RELEASE_TOKEN_POOL_MCMS_CAP_TEST has drop {}

public struct TestEnv {
    scenario: Scenario,
    state: LockReleaseTokenPoolState<LOCK_RELEASE_TOKEN_POOL_MCMS_CAP_TEST>,
    ccip_ref: CCIPObjectRef,
    mcms_registry: Registry,
}

// ================================================================
// |                      Setup Helpers                           |
// ================================================================

fun setup(): (TestEnv, OwnerCap, RebalancerCap<LOCK_RELEASE_TOKEN_POOL_MCMS_CAP_TEST>) {
    let mut scenario = ts::begin(OWNER);

    // Setup CCIP environment
    scenario.next_tx(CCIP_ADMIN);
    let ctx = scenario.ctx();
    state_object::test_init(ctx);

    scenario.next_tx(CCIP_ADMIN);
    let ccip_owner_cap = scenario.take_from_sender<CCIPOwnerCap>();
    let mut ccip_ref = scenario.take_shared<CCIPObjectRef>();

    // Initialize required CCIP modules
    upgrade_registry::initialize(&mut ccip_ref, &ccip_owner_cap, scenario.ctx());
    rmn_remote::initialize(&mut ccip_ref, &ccip_owner_cap, 1000, scenario.ctx());
    token_admin_registry::initialize(&mut ccip_ref, &ccip_owner_cap, scenario.ctx());
    onramp_state_helper::test_init(scenario.ctx());
    offramp_state_helper::test_init(scenario.ctx());

    // Initialize MCMS registry
    mcms_registry::test_init(scenario.ctx());

    // Initialize token pool
    scenario.next_tx(OWNER);
    let ctx = scenario.ctx();
    let (treasury_cap, coin_metadata) = coin::create_currency(
        LOCK_RELEASE_TOKEN_POOL_MCMS_CAP_TEST {},
        8, // decimals
        b"TEST",
        b"TestToken",
        b"test_token",
        option::none(),
        ctx,
    );

    // Call test_init to create owner_cap
    lock_release_token_pool::test_init(ctx);

    transfer::public_freeze_object(coin_metadata);
    transfer::public_transfer(treasury_cap, OWNER);
    transfer::public_transfer(ccip_owner_cap, @0x0);
    ts::return_shared(ccip_ref);

    // Now take the owner_cap that was created by test_init and initialize the pool
    scenario.next_tx(OWNER);
    let mut owner_cap_for_init = ts::take_from_sender<OwnerCap>(&scenario);
    let mut ccip_ref = ts::take_shared<CCIPObjectRef>(&scenario);
    let coin_metadata = ts::take_immutable<
        coin::CoinMetadata<LOCK_RELEASE_TOKEN_POOL_MCMS_CAP_TEST>,
    >(&scenario);
    let treasury_cap = ts::take_from_sender<
        coin::TreasuryCap<LOCK_RELEASE_TOKEN_POOL_MCMS_CAP_TEST>,
    >(&scenario);

    lock_release_token_pool::initialize(
        &mut owner_cap_for_init,
        &mut ccip_ref,
        &coin_metadata,
        &treasury_cap,
        TOKEN_ADMIN,
        EOA_REBALANCER,
        scenario.ctx(),
    );

    transfer::public_transfer(owner_cap_for_init, OWNER);
    transfer::public_transfer(treasury_cap, OWNER);
    ts::return_immutable(coin_metadata);
    ts::return_shared(ccip_ref);

    scenario.next_tx(OWNER);
    let state = ts::take_shared<LockReleaseTokenPoolState<LOCK_RELEASE_TOKEN_POOL_MCMS_CAP_TEST>>(
        &scenario,
    );
    let owner_cap = ts::take_from_sender<OwnerCap>(
        &scenario,
    );
    let ccip_ref = ts::take_shared<CCIPObjectRef>(&scenario);

    scenario.next_tx(EOA_REBALANCER);
    let rebalancer_cap = ts::take_from_sender<RebalancerCap<LOCK_RELEASE_TOKEN_POOL_MCMS_CAP_TEST>>(
        &scenario,
    );

    scenario.next_tx(OWNER);
    let mcms_registry = ts::take_shared<Registry>(&scenario);

    let env = TestEnv {
        scenario,
        state,
        ccip_ref,
        mcms_registry,
    };

    (env, owner_cap, rebalancer_cap)
}

fun tear_down(env: TestEnv) {
    let TestEnv { scenario, state, ccip_ref, mcms_registry } = env;

    ts::return_shared(state);
    ts::return_shared(ccip_ref);
    ts::return_shared(mcms_registry);
    ts::end(scenario);
}

fun create_mcms_callback_params(
    target: address,
    function_name: vector<u8>,
    data: vector<u8>,
    batch_id: vector<u8>,
    sequence_number: u64,
): ExecutingCallbackParams {
    mcms_registry::test_create_executing_callback_params(
        target,
        string::utf8(b"lock_release_token_pool"),
        string::utf8(function_name),
        data,
        batch_id,
        sequence_number,
        1, // total_in_batch
    )
}

fun encode_set_rebalancer_data(
    state_id: address,
    owner_cap_id: address,
    rebalancer: address,
): vector<u8> {
    let mut data = vector::empty();
    vector::append(&mut data, bcs::to_bytes(&state_id));
    vector::append(&mut data, bcs::to_bytes(&owner_cap_id));
    vector::append(&mut data, bcs::to_bytes(&rebalancer));
    data
}

fun encode_provide_liquidity_data(
    state_id: address,
    rebalancer_cap_id: address,
    coin_id: address,
): vector<u8> {
    let mut data = vector::empty();
    vector::append(&mut data, bcs::to_bytes(&state_id));
    vector::append(&mut data, bcs::to_bytes(&rebalancer_cap_id));
    vector::append(&mut data, bcs::to_bytes(&coin_id));
    data
}

fun encode_withdraw_liquidity_data(
    state_id: address,
    rebalancer_cap_id: address,
    amount: u64,
    to: address,
): vector<u8> {
    let mut data = vector::empty();
    vector::append(&mut data, bcs::to_bytes(&state_id));
    vector::append(&mut data, bcs::to_bytes(&rebalancer_cap_id));
    vector::append(&mut data, bcs::to_bytes(&amount));
    vector::append(&mut data, bcs::to_bytes(&to));
    data
}

fun setup_mcms_ownership(env: &mut TestEnv, owner_cap: OwnerCap) {
    // Transfer ownership to MCMS
    env.scenario.next_tx(OWNER);
    lock_release_token_pool::transfer_ownership(
        &mut env.state,
        &owner_cap,
        mcms_registry::get_multisig_address(),
        env.scenario.ctx(),
    );

    // Accept ownership as MCMS
    env.scenario.next_tx(mcms_registry::get_multisig_address());
    let state_addr = object::id_address(&env.state);
    let params = create_mcms_callback_params(
        @lock_release_token_pool,
        b"accept_ownership",
        bcs::to_bytes(&state_addr),
        x"0000000000000000000000000000000000000000000000000000000000000001",
        0,
    );
    lock_release_token_pool::mcms_accept_ownership(
        &mut env.state,
        &mut env.mcms_registry,
        params,
        env.scenario.ctx(),
    );

    // Execute ownership transfer
    env.scenario.next_tx(OWNER);
    lock_release_token_pool::execute_ownership_transfer_to_mcms(
        owner_cap,
        &mut env.state,
        &mut env.mcms_registry,
        mcms_registry::get_multisig_address(),
        env.scenario.ctx(),
    );
}

// Setup MCMS ownership AND take rebalancer control
fun setup_mcms_with_rebalancer(
    env: &mut TestEnv,
    owner_cap: OwnerCap,
    owner_cap_id: address,
): address {
    setup_mcms_ownership(env, owner_cap);

    // MCMS takes rebalancer control
    env.scenario.next_tx(mcms_registry::get_multisig_address());
    let mcms_address = mcms_registry::get_multisig_address();
    let data = encode_set_rebalancer_data(
        object::id_address(&env.state),
        owner_cap_id,
        mcms_address,
    );
    let params = create_mcms_callback_params(
        @lock_release_token_pool,
        b"set_rebalancer",
        data,
        x"0000000000000000000000000000000000000000000000000000000000000002",
        0,
    );
    lock_release_token_pool::mcms_set_rebalancer(
        &mut env.state,
        &mut env.mcms_registry,
        params,
        env.scenario.ctx(),
    );

    // Return rebalancer cap ID
    let mcms_cap = mcms_registry::get_cap<McmsCap<LOCK_RELEASE_TOKEN_POOL_MCMS_CAP_TEST>>(
        &env.mcms_registry,
        @lock_release_token_pool.to_ascii_string(),
    );
    lock_release_token_pool::mcms_rebalancer_cap_address(mcms_cap)
}

#[test]
public fun test_set_rebalancer_to_eoa() {
    let (mut env, owner_cap, rebalancer_cap) = setup();

    // Get initial rebalancer cap ID
    let initial_rebalancer_id = object::id(&rebalancer_cap);
    assert!(
        lock_release_token_pool::get_rebalancer(&env.state) == object::id_to_address(&initial_rebalancer_id),
    );

    // Set rebalancer to new EOA
    env.scenario.next_tx(OWNER);
    lock_release_token_pool::set_rebalancer(
        &mut env.state,
        &owner_cap,
        EOA_REBALANCER_2,
        env.scenario.ctx(),
    );

    // Verify new rebalancer cap was created and transferred
    env.scenario.next_tx(EOA_REBALANCER_2);
    let new_rebalancer_cap = ts::take_from_sender<
        RebalancerCap<LOCK_RELEASE_TOKEN_POOL_MCMS_CAP_TEST>,
    >(&env.scenario);
    let new_rebalancer_id = object::id(&new_rebalancer_cap);

    // Verify state was updated
    assert!(
        lock_release_token_pool::get_rebalancer(&env.state) == object::id_to_address(&new_rebalancer_id),
    );
    assert!(new_rebalancer_id != initial_rebalancer_id);

    ts::return_to_address(EOA_REBALANCER_2, new_rebalancer_cap);
    ts::return_to_address(EOA_REBALANCER, rebalancer_cap);
    ts::return_to_address(OWNER, owner_cap);
    tear_down(env);
}

#[test]
public fun test_set_rebalancer_rotation() {
    let (mut env, owner_cap, rebalancer_cap) = setup();

    // Rotate to first new rebalancer
    env.scenario.next_tx(OWNER);
    lock_release_token_pool::set_rebalancer(
        &mut env.state,
        &owner_cap,
        EOA_REBALANCER_2,
        env.scenario.ctx(),
    );

    env.scenario.next_tx(EOA_REBALANCER_2);
    let new_rebalancer_cap = ts::take_from_sender<
        RebalancerCap<LOCK_RELEASE_TOKEN_POOL_MCMS_CAP_TEST>,
    >(&env.scenario);

    // Rotate again to another rebalancer
    env.scenario.next_tx(OWNER);
    lock_release_token_pool::set_rebalancer(
        &mut env.state,
        &owner_cap,
        NEW_OWNER,
        env.scenario.ctx(),
    );

    env.scenario.next_tx(NEW_OWNER);
    let newest_rebalancer_cap = ts::take_from_sender<
        RebalancerCap<LOCK_RELEASE_TOKEN_POOL_MCMS_CAP_TEST>,
    >(&env.scenario);

    // Now destroy the old caps (from first rebalancer and second rebalancer)
    env.scenario.next_tx(EOA_REBALANCER);
    lock_release_token_pool::destroy_rebalancer_cap(
        &mut env.state,
        rebalancer_cap,
        env.scenario.ctx(),
    );

    env.scenario.next_tx(EOA_REBALANCER_2);
    lock_release_token_pool::destroy_rebalancer_cap(
        &mut env.state,
        new_rebalancer_cap,
        env.scenario.ctx(),
    );

    // Verify the newest cap is still valid
    assert!(
        lock_release_token_pool::get_rebalancer(&env.state) == object::id_to_address(&object::id(&newest_rebalancer_cap)),
    );

    ts::return_to_address(NEW_OWNER, newest_rebalancer_cap);
    ts::return_to_address(OWNER, owner_cap);
    tear_down(env);
}

#[test]
public fun test_execute_ownership_transfer_to_mcms_with_rebalancer() {
    let (mut env, owner_cap, rebalancer_cap) = setup();
    let owner_cap_id = object::id_address(&owner_cap);

    // Setup MCMS with rebalancer control
    setup_mcms_with_rebalancer(&mut env, owner_cap, owner_cap_id);

    // Verify ownership transferred
    assert!(lock_release_token_pool::owner(&env.state) == mcms_registry::get_multisig_address());

    // Assert package is registered
    assert!(
        mcms_registry::is_package_registered(
            &env.mcms_registry,
            @lock_release_token_pool.to_ascii_string(),
        ),
    );

    // Assert old rebalancer cap is not registered anymore
    // New rebalancer cap was created and stored in MCMS
    assert!(
        lock_release_token_pool::get_rebalancer(&env.state) != object::id_address(&rebalancer_cap),
    );

    ts::return_to_address(EOA_REBALANCER, rebalancer_cap);
    tear_down(env);
}

#[test]
public fun test_execute_ownership_transfer_to_mcms_without_rebalancer() {
    let (mut env, owner_cap, rebalancer_cap) = setup();

    // Setup MCMS ownership only (no rebalancer)
    setup_mcms_ownership(&mut env, owner_cap);

    // Verify ownership transferred
    assert!(lock_release_token_pool::owner(&env.state) == mcms_registry::get_multisig_address());

    // Verify LockReleaseTokenPoolState rebalancer cap id not changed
    assert!(
        lock_release_token_pool::get_rebalancer(&env.state) == object::id_address(&rebalancer_cap),
    );

    ts::return_to_address(EOA_REBALANCER, rebalancer_cap);
    tear_down(env);
}

#[test]
public fun test_mcms_set_rebalancer_to_mcms_address() {
    let (mut env, owner_cap, rebalancer_cap) = setup();
    let owner_cap_id = object::id_address(&owner_cap);

    // Setup MCMS with rebalancer control
    setup_mcms_with_rebalancer(&mut env, owner_cap, owner_cap_id);

    // Verify LockReleaseTokenPoolState rebalancer cap id was updated to MCMS rebalancer cap address
    let mcms_cap = mcms_registry::get_cap<McmsCap<LOCK_RELEASE_TOKEN_POOL_MCMS_CAP_TEST>>(
        &env.mcms_registry,
        @lock_release_token_pool.to_ascii_string(),
    );
    let mcms_rebalancer_cap_addr = lock_release_token_pool::mcms_rebalancer_cap_address<
        LOCK_RELEASE_TOKEN_POOL_MCMS_CAP_TEST,
    >(mcms_cap);
    assert!(mcms_rebalancer_cap_addr == lock_release_token_pool::get_rebalancer(&env.state));

    ts::return_to_address(EOA_REBALANCER, rebalancer_cap);
    tear_down(env);
}

#[test]
public fun test_mcms_retake_rebalancer_after_delegation() {
    let (mut env, owner_cap, rebalancer_cap) = setup();
    let owner_cap_id = object::id_address(&owner_cap);

    // Setup MCMS with rebalancer control
    setup_mcms_with_rebalancer(&mut env, owner_cap, owner_cap_id);

    // Verify MCMS has rebalancer control
    let mcms_cap = mcms_registry::get_cap<McmsCap<LOCK_RELEASE_TOKEN_POOL_MCMS_CAP_TEST>>(
        &env.mcms_registry,
        @lock_release_token_pool.to_ascii_string(),
    );
    let mcms_rebalancer_cap_addr_1 = lock_release_token_pool::mcms_rebalancer_cap_address(mcms_cap);
    assert!(mcms_rebalancer_cap_addr_1 == lock_release_token_pool::get_rebalancer(&env.state));
    let mcms_address = mcms_registry::get_multisig_address();

    // Step 3: MCMS delegates rebalancer to EOA
    env.scenario.next_tx(mcms_registry::get_multisig_address());
    let data = encode_set_rebalancer_data(
        object::id_address(&env.state),
        owner_cap_id,
        EOA_REBALANCER_2,
    );
    let params = create_mcms_callback_params(
        @lock_release_token_pool,
        b"set_rebalancer",
        data,
        x"0000000000000000000000000000000000000000000000000000000000000003",
        0,
    );
    lock_release_token_pool::mcms_set_rebalancer(
        &mut env.state,
        &mut env.mcms_registry,
        params,
        env.scenario.ctx(),
    );

    // Verify EOA has rebalancer control
    env.scenario.next_tx(EOA_REBALANCER_2);
    let eoa_rebalancer_cap = ts::take_from_sender<
        RebalancerCap<LOCK_RELEASE_TOKEN_POOL_MCMS_CAP_TEST>,
    >(&env.scenario);
    assert!(
        object::id_address(&eoa_rebalancer_cap) == lock_release_token_pool::get_rebalancer(&env.state),
    );

    // Step 4: MCMS re-takes rebalancer control
    env.scenario.next_tx(mcms_registry::get_multisig_address());
    let data = encode_set_rebalancer_data(
        object::id_address(&env.state),
        owner_cap_id,
        mcms_address,
    );
    let params = create_mcms_callback_params(
        @lock_release_token_pool,
        b"set_rebalancer",
        data,
        x"0000000000000000000000000000000000000000000000000000000000000004",
        0,
    );
    lock_release_token_pool::mcms_set_rebalancer(
        &mut env.state,
        &mut env.mcms_registry,
        params,
        env.scenario.ctx(),
    );

    // Verify MCMS has rebalancer control again
    let mcms_cap = mcms_registry::get_cap<McmsCap<LOCK_RELEASE_TOKEN_POOL_MCMS_CAP_TEST>>(
        &env.mcms_registry,
        @lock_release_token_pool.to_ascii_string(),
    );
    let mcms_rebalancer_cap_addr_2 = lock_release_token_pool::mcms_rebalancer_cap_address(mcms_cap);
    assert!(mcms_rebalancer_cap_addr_2 == lock_release_token_pool::get_rebalancer(&env.state));

    // New MCMS cap should be different from the first one
    assert!(mcms_rebalancer_cap_addr_1 != mcms_rebalancer_cap_addr_2);

    // Step 5: Verify old EOA cap can be destroyed
    env.scenario.next_tx(EOA_REBALANCER_2);
    lock_release_token_pool::destroy_rebalancer_cap(
        &mut env.state,
        eoa_rebalancer_cap,
        env.scenario.ctx(),
    );

    ts::return_to_address(EOA_REBALANCER, rebalancer_cap);
    tear_down(env);
}

#[test]
public fun test_mcms_set_rebalancer_idempotent() {
    let (mut env, owner_cap, rebalancer_cap) = setup();
    let owner_cap_id = object::id_address(&owner_cap);

    // Setup MCMS with rebalancer control
    setup_mcms_with_rebalancer(&mut env, owner_cap, owner_cap_id);

    let mcms_address = mcms_registry::get_multisig_address();

    // Capture the rebalancer cap ID
    let mcms_cap = mcms_registry::get_cap<McmsCap<LOCK_RELEASE_TOKEN_POOL_MCMS_CAP_TEST>>(
        &env.mcms_registry,
        @lock_release_token_pool.to_ascii_string(),
    );
    let initial_rebalancer_cap_addr = lock_release_token_pool::mcms_rebalancer_cap_address(
        mcms_cap,
    );
    assert!(initial_rebalancer_cap_addr == lock_release_token_pool::get_rebalancer(&env.state));

    // Call mcms_set_rebalancer again with MCMS address (should be idempotent no-op)
    env.scenario.next_tx(mcms_registry::get_multisig_address());
    let data = encode_set_rebalancer_data(
        object::id_address(&env.state),
        owner_cap_id,
        mcms_address,
    );
    let params = create_mcms_callback_params(
        @lock_release_token_pool,
        b"set_rebalancer",
        data,
        x"0000000000000000000000000000000000000000000000000000000000000003",
        0,
    );
    lock_release_token_pool::mcms_set_rebalancer(
        &mut env.state,
        &mut env.mcms_registry,
        params,
        env.scenario.ctx(),
    );

    // Verify rebalancer cap ID unchanged (idempotent)
    let mcms_cap = mcms_registry::get_cap<McmsCap<LOCK_RELEASE_TOKEN_POOL_MCMS_CAP_TEST>>(
        &env.mcms_registry,
        @lock_release_token_pool.to_ascii_string(),
    );
    let final_rebalancer_cap_addr = lock_release_token_pool::mcms_rebalancer_cap_address(mcms_cap);
    assert!(final_rebalancer_cap_addr == initial_rebalancer_cap_addr);
    assert!(final_rebalancer_cap_addr == lock_release_token_pool::get_rebalancer(&env.state));

    ts::return_to_address(EOA_REBALANCER, rebalancer_cap);
    tear_down(env);
}

#[test]
public fun test_mcms_multiple_delegation_cycles() {
    let (mut env, owner_cap, rebalancer_cap) = setup();
    let owner_cap_id = object::id_address(&owner_cap);

    // Setup MCMS with rebalancer control
    setup_mcms_with_rebalancer(&mut env, owner_cap, owner_cap_id);

    let mcms_address = mcms_registry::get_multisig_address();

    // Verify Cycle 1: MCMS has rebalancer control
    let mcms_cap = mcms_registry::get_cap<McmsCap<LOCK_RELEASE_TOKEN_POOL_MCMS_CAP_TEST>>(
        &env.mcms_registry,
        @lock_release_token_pool.to_ascii_string(),
    );
    let mcms_rebalancer_cap_addr_1 = lock_release_token_pool::mcms_rebalancer_cap_address(mcms_cap);
    assert!(mcms_rebalancer_cap_addr_1 == lock_release_token_pool::get_rebalancer(&env.state));

    // Cycle 1 → 2: MCMS delegates to EOA_REBALANCER_2
    env.scenario.next_tx(mcms_registry::get_multisig_address());
    let data = encode_set_rebalancer_data(
        object::id_address(&env.state),
        owner_cap_id,
        EOA_REBALANCER_2,
    );
    let params = create_mcms_callback_params(
        @lock_release_token_pool,
        b"set_rebalancer",
        data,
        x"0000000000000000000000000000000000000000000000000000000000000003",
        0,
    );
    lock_release_token_pool::mcms_set_rebalancer(
        &mut env.state,
        &mut env.mcms_registry,
        params,
        env.scenario.ctx(),
    );

    env.scenario.next_tx(EOA_REBALANCER_2);
    let eoa_rebalancer_cap_2 = ts::take_from_sender<
        RebalancerCap<LOCK_RELEASE_TOKEN_POOL_MCMS_CAP_TEST>,
    >(&env.scenario);
    assert!(
        object::id_address(&eoa_rebalancer_cap_2) == lock_release_token_pool::get_rebalancer(&env.state),
    );

    // Cycle 2 → 3: MCMS re-takes control
    env.scenario.next_tx(mcms_registry::get_multisig_address());
    let data = encode_set_rebalancer_data(
        object::id_address(&env.state),
        owner_cap_id,
        mcms_address,
    );
    let params = create_mcms_callback_params(
        @lock_release_token_pool,
        b"set_rebalancer",
        data,
        x"0000000000000000000000000000000000000000000000000000000000000004",
        0,
    );
    lock_release_token_pool::mcms_set_rebalancer(
        &mut env.state,
        &mut env.mcms_registry,
        params,
        env.scenario.ctx(),
    );

    let mcms_cap = mcms_registry::get_cap<McmsCap<LOCK_RELEASE_TOKEN_POOL_MCMS_CAP_TEST>>(
        &env.mcms_registry,
        @lock_release_token_pool.to_ascii_string(),
    );
    let mcms_rebalancer_cap_addr_2 = lock_release_token_pool::mcms_rebalancer_cap_address(mcms_cap);
    assert!(mcms_rebalancer_cap_addr_2 == lock_release_token_pool::get_rebalancer(&env.state));
    assert!(mcms_rebalancer_cap_addr_1 != mcms_rebalancer_cap_addr_2);

    // Cycle 3 → 4: MCMS delegates to NEW_OWNER
    env.scenario.next_tx(mcms_registry::get_multisig_address());
    let data = encode_set_rebalancer_data(
        object::id_address(&env.state),
        owner_cap_id,
        NEW_OWNER,
    );
    let params = create_mcms_callback_params(
        @lock_release_token_pool,
        b"set_rebalancer",
        data,
        x"0000000000000000000000000000000000000000000000000000000000000005",
        0,
    );
    lock_release_token_pool::mcms_set_rebalancer(
        &mut env.state,
        &mut env.mcms_registry,
        params,
        env.scenario.ctx(),
    );

    env.scenario.next_tx(NEW_OWNER);
    let eoa_rebalancer_cap_new = ts::take_from_sender<
        RebalancerCap<LOCK_RELEASE_TOKEN_POOL_MCMS_CAP_TEST>,
    >(&env.scenario);
    assert!(
        object::id_address(&eoa_rebalancer_cap_new) == lock_release_token_pool::get_rebalancer(&env.state),
    );

    // Cycle 4 → Final: MCMS takes final control
    env.scenario.next_tx(mcms_registry::get_multisig_address());
    let data = encode_set_rebalancer_data(
        object::id_address(&env.state),
        owner_cap_id,
        mcms_address,
    );
    let params = create_mcms_callback_params(
        @lock_release_token_pool,
        b"set_rebalancer",
        data,
        x"0000000000000000000000000000000000000000000000000000000000000006",
        0,
    );
    lock_release_token_pool::mcms_set_rebalancer(
        &mut env.state,
        &mut env.mcms_registry,
        params,
        env.scenario.ctx(),
    );

    // Verify final MCMS control
    let mcms_cap = mcms_registry::get_cap<McmsCap<LOCK_RELEASE_TOKEN_POOL_MCMS_CAP_TEST>>(
        &env.mcms_registry,
        @lock_release_token_pool.to_ascii_string(),
    );
    let mcms_rebalancer_cap_addr_final = lock_release_token_pool::mcms_rebalancer_cap_address(
        mcms_cap,
    );
    assert!(mcms_rebalancer_cap_addr_final == lock_release_token_pool::get_rebalancer(&env.state));
    assert!(mcms_rebalancer_cap_addr_final != mcms_rebalancer_cap_addr_2);

    // Clean up old caps
    env.scenario.next_tx(EOA_REBALANCER_2);
    lock_release_token_pool::destroy_rebalancer_cap(
        &mut env.state,
        eoa_rebalancer_cap_2,
        env.scenario.ctx(),
    );

    env.scenario.next_tx(NEW_OWNER);
    lock_release_token_pool::destroy_rebalancer_cap(
        &mut env.state,
        eoa_rebalancer_cap_new,
        env.scenario.ctx(),
    );

    ts::return_to_address(EOA_REBALANCER, rebalancer_cap);
    tear_down(env);
}

// ================================================================
// |                    Liquidity Management Tests                |
// ================================================================

#[test]
public fun test_mcms_provide_liquidity_success() {
    let (mut env, owner_cap, rebalancer_cap) = setup();
    let owner_cap_id = object::id_address(&owner_cap);

    // Setup MCMS with rebalancer control
    let rebalancer_cap_id = setup_mcms_with_rebalancer(&mut env, owner_cap, owner_cap_id);

    // Get initial balance
    let initial_balance = lock_release_token_pool::get_balance(&env.state);

    // Mint test coins
    env.scenario.next_tx(OWNER);
    let mut treasury_cap = ts::take_from_sender<
        coin::TreasuryCap<LOCK_RELEASE_TOKEN_POOL_MCMS_CAP_TEST>,
    >(&env.scenario);
    let test_coin = coin::mint(&mut treasury_cap, 1000000, env.scenario.ctx());
    let coin_id = object::id_address(&test_coin);

    // Provide liquidity via MCMS
    env.scenario.next_tx(mcms_registry::get_multisig_address());
    let data = encode_provide_liquidity_data(
        object::id_address(&env.state),
        rebalancer_cap_id,
        coin_id,
    );
    let params = create_mcms_callback_params(
        @lock_release_token_pool,
        b"provide_liquidity",
        data,
        x"0000000000000000000000000000000000000000000000000000000000000003",
        0,
    );
    lock_release_token_pool::mcms_provide_liquidity(
        &mut env.state,
        &mut env.mcms_registry,
        test_coin,
        params,
        env.scenario.ctx(),
    );

    // Verify balance increased
    let final_balance = lock_release_token_pool::get_balance(&env.state);
    assert!(final_balance == initial_balance + 1000000);

    ts::return_to_address(OWNER, treasury_cap);
    ts::return_to_address(EOA_REBALANCER, rebalancer_cap);
    tear_down(env);
}

#[test]
public fun test_mcms_withdraw_liquidity_success() {
    let (mut env, owner_cap, rebalancer_cap) = setup();
    let owner_cap_id = object::id_address(&owner_cap);

    // Setup MCMS with rebalancer control
    let rebalancer_cap_id = setup_mcms_with_rebalancer(&mut env, owner_cap, owner_cap_id);

    // Provide initial liquidity
    env.scenario.next_tx(OWNER);
    let mut treasury_cap = ts::take_from_sender<
        coin::TreasuryCap<LOCK_RELEASE_TOKEN_POOL_MCMS_CAP_TEST>,
    >(&env.scenario);
    let test_coin = coin::mint(&mut treasury_cap, 1000000, env.scenario.ctx());
    let coin_id = object::id_address(&test_coin);

    env.scenario.next_tx(mcms_registry::get_multisig_address());
    let data = encode_provide_liquidity_data(
        object::id_address(&env.state),
        rebalancer_cap_id,
        coin_id,
    );
    let params = create_mcms_callback_params(
        @lock_release_token_pool,
        b"provide_liquidity",
        data,
        x"0000000000000000000000000000000000000000000000000000000000000003",
        0,
    );
    lock_release_token_pool::mcms_provide_liquidity(
        &mut env.state,
        &mut env.mcms_registry,
        test_coin,
        params,
        env.scenario.ctx(),
    );

    // Get balance before withdrawal
    let balance_before = lock_release_token_pool::get_balance(&env.state);
    assert!(balance_before == 1000000);

    // Withdraw liquidity via MCMS
    env.scenario.next_tx(mcms_registry::get_multisig_address());
    let data = encode_withdraw_liquidity_data(
        object::id_address(&env.state),
        rebalancer_cap_id,
        500000,
        NEW_OWNER, // recipient
    );
    let params = create_mcms_callback_params(
        @lock_release_token_pool,
        b"withdraw_liquidity",
        data,
        x"0000000000000000000000000000000000000000000000000000000000000004",
        0,
    );
    lock_release_token_pool::mcms_withdraw_liquidity(
        &mut env.state,
        &mut env.mcms_registry,
        params,
        env.scenario.ctx(),
    );

    // Verify balance decreased
    let balance_after = lock_release_token_pool::get_balance(&env.state);
    assert!(balance_after == 500000);

    // Verify recipient received the coin
    env.scenario.next_tx(NEW_OWNER);
    let received_coin = ts::take_from_sender<Coin<LOCK_RELEASE_TOKEN_POOL_MCMS_CAP_TEST>>(
        &env.scenario,
    );
    assert!(received_coin.value() == 500000);

    ts::return_to_address(NEW_OWNER, received_coin);
    ts::return_to_address(OWNER, treasury_cap);
    ts::return_to_address(EOA_REBALANCER, rebalancer_cap);
    tear_down(env);
}

#[test]
public fun test_mcms_provide_and_withdraw_liquidity_full_cycle() {
    let (mut env, owner_cap, rebalancer_cap) = setup();
    let owner_cap_id = object::id_address(&owner_cap);

    // Setup MCMS with rebalancer control
    let rebalancer_cap_id = setup_mcms_with_rebalancer(&mut env, owner_cap, owner_cap_id);

    // Step 1: Provide initial liquidity (1,000,000)
    env.scenario.next_tx(OWNER);
    let mut treasury_cap = ts::take_from_sender<
        coin::TreasuryCap<LOCK_RELEASE_TOKEN_POOL_MCMS_CAP_TEST>,
    >(&env.scenario);
    let test_coin_1 = coin::mint(&mut treasury_cap, 1000000, env.scenario.ctx());
    let coin_id_1 = object::id_address(&test_coin_1);

    env.scenario.next_tx(mcms_registry::get_multisig_address());
    let data = encode_provide_liquidity_data(
        object::id_address(&env.state),
        rebalancer_cap_id,
        coin_id_1,
    );
    let params = create_mcms_callback_params(
        @lock_release_token_pool,
        b"provide_liquidity",
        data,
        x"0000000000000000000000000000000000000000000000000000000000000003",
        0,
    );
    lock_release_token_pool::mcms_provide_liquidity(
        &mut env.state,
        &mut env.mcms_registry,
        test_coin_1,
        params,
        env.scenario.ctx(),
    );
    assert!(lock_release_token_pool::get_balance(&env.state) == 1000000);

    // Step 2: Withdraw some liquidity (300,000)
    env.scenario.next_tx(mcms_registry::get_multisig_address());
    let data = encode_withdraw_liquidity_data(
        object::id_address(&env.state),
        rebalancer_cap_id,
        300000,
        NEW_OWNER,
    );
    let params = create_mcms_callback_params(
        @lock_release_token_pool,
        b"withdraw_liquidity",
        data,
        x"0000000000000000000000000000000000000000000000000000000000000004",
        0,
    );
    lock_release_token_pool::mcms_withdraw_liquidity(
        &mut env.state,
        &mut env.mcms_registry,
        params,
        env.scenario.ctx(),
    );
    assert!(lock_release_token_pool::get_balance(&env.state) == 700000);

    // Step 3: Provide more liquidity (500,000)
    env.scenario.next_tx(OWNER);
    let test_coin_2 = coin::mint(&mut treasury_cap, 500000, env.scenario.ctx());
    let coin_id_2 = object::id_address(&test_coin_2);

    env.scenario.next_tx(mcms_registry::get_multisig_address());
    let data = encode_provide_liquidity_data(
        object::id_address(&env.state),
        rebalancer_cap_id,
        coin_id_2,
    );
    let params = create_mcms_callback_params(
        @lock_release_token_pool,
        b"provide_liquidity",
        data,
        x"0000000000000000000000000000000000000000000000000000000000000005",
        0,
    );
    lock_release_token_pool::mcms_provide_liquidity(
        &mut env.state,
        &mut env.mcms_registry,
        test_coin_2,
        params,
        env.scenario.ctx(),
    );
    assert!(lock_release_token_pool::get_balance(&env.state) == 1200000);

    // Step 4: Withdraw to different address
    env.scenario.next_tx(mcms_registry::get_multisig_address());
    let data = encode_withdraw_liquidity_data(
        object::id_address(&env.state),
        rebalancer_cap_id,
        200000,
        OTHER_USER,
    );
    let params = create_mcms_callback_params(
        @lock_release_token_pool,
        b"withdraw_liquidity",
        data,
        x"0000000000000000000000000000000000000000000000000000000000000006",
        0,
    );
    lock_release_token_pool::mcms_withdraw_liquidity(
        &mut env.state,
        &mut env.mcms_registry,
        params,
        env.scenario.ctx(),
    );
    assert!(lock_release_token_pool::get_balance(&env.state) == 1000000);

    // Verify recipients received coins
    env.scenario.next_tx(NEW_OWNER);
    let coin_new_owner = ts::take_from_sender<Coin<LOCK_RELEASE_TOKEN_POOL_MCMS_CAP_TEST>>(
        &env.scenario,
    );
    assert!(coin_new_owner.value() == 300000);

    env.scenario.next_tx(OTHER_USER);
    let coin_other_user = ts::take_from_sender<Coin<LOCK_RELEASE_TOKEN_POOL_MCMS_CAP_TEST>>(
        &env.scenario,
    );
    assert!(coin_other_user.value() == 200000);

    ts::return_to_address(NEW_OWNER, coin_new_owner);
    ts::return_to_address(OTHER_USER, coin_other_user);
    ts::return_to_address(OWNER, treasury_cap);
    ts::return_to_address(EOA_REBALANCER, rebalancer_cap);
    tear_down(env);
}

// ================================================================
// |                      Negative Tests                          |
// ================================================================

#[test]
#[expected_failure(abort_code = lock_release_token_pool::EInvalidRebalancer)]
public fun test_set_rebalancer_fails_if_target_is_mcms() {
    let (mut env, owner_cap, rebalancer_cap) = setup();

    // Try to set rebalancer to MCMS address using set_rebalancer (should fail)
    env.scenario.next_tx(OWNER);
    lock_release_token_pool::set_rebalancer(
        &mut env.state,
        &owner_cap,
        mcms_registry::get_multisig_address(),
        env.scenario.ctx(),
    );

    ts::return_to_address(EOA_REBALANCER, rebalancer_cap);
    ts::return_to_address(OWNER, owner_cap);
    tear_down(env);
}

#[test]
#[expected_failure(abort_code = lock_release_token_pool::ERebalancerCapIsInUse)]
public fun test_destroy_rebalancer_cap_fails_if_active() {
    let (mut env, owner_cap, rebalancer_cap) = setup();

    // Try to destroy the current active rebalancer cap (should fail)
    env.scenario.next_tx(EOA_REBALANCER);
    lock_release_token_pool::destroy_rebalancer_cap(
        &mut env.state,
        rebalancer_cap,
        env.scenario.ctx(),
    );

    ts::return_to_address(OWNER, owner_cap);
    tear_down(env);
}

#[test]
#[expected_failure(abort_code = lock_release_token_pool::EInvalidOwnerCap)]
public fun test_set_rebalancer_unauthorized() {
    let (mut env, owner_cap, rebalancer_cap) = setup();

    // Create a fake owner cap
    env.scenario.next_tx(OTHER_USER);
    let fake_owner_cap = ownable::create_test_owner_cap(env.scenario.ctx());

    // Try to set rebalancer with fake cap (should fail)
    lock_release_token_pool::set_rebalancer(
        &mut env.state,
        &fake_owner_cap,
        NEW_OWNER,
        env.scenario.ctx(),
    );

    ownable::test_destroy_owner_cap(fake_owner_cap);
    ts::return_to_address(EOA_REBALANCER, rebalancer_cap);
    ts::return_to_address(OWNER, owner_cap);
    tear_down(env);
}

#[test]
#[expected_failure(abort_code = lock_release_token_pool::ERebalancerCapDoesNotExist)]
public fun test_mcms_provide_liquidity_fails_without_rebalancer_cap() {
    let (mut env, owner_cap, rebalancer_cap) = setup();

    // Transfer ownership to MCMS but DON'T take rebalancer control
    setup_mcms_ownership(&mut env, owner_cap);

    // Try to provide liquidity without having rebalancer cap (should fail)
    env.scenario.next_tx(OWNER);
    let mut treasury_cap = ts::take_from_sender<
        coin::TreasuryCap<LOCK_RELEASE_TOKEN_POOL_MCMS_CAP_TEST>,
    >(&env.scenario);
    let test_coin = coin::mint(&mut treasury_cap, 1000000, env.scenario.ctx());
    let coin_id = object::id_address(&test_coin);

    env.scenario.next_tx(mcms_registry::get_multisig_address());
    let data = encode_provide_liquidity_data(
        object::id_address(&env.state),
        @0x0, // dummy rebalancer cap id
        coin_id,
    );
    let params = create_mcms_callback_params(
        @lock_release_token_pool,
        b"provide_liquidity",
        data,
        x"0000000000000000000000000000000000000000000000000000000000000002",
        0,
    );
    lock_release_token_pool::mcms_provide_liquidity(
        &mut env.state,
        &mut env.mcms_registry,
        test_coin,
        params,
        env.scenario.ctx(),
    );

    ts::return_to_address(OWNER, treasury_cap);
    ts::return_to_address(EOA_REBALANCER, rebalancer_cap);
    tear_down(env);
}

#[test]
#[expected_failure(abort_code = lock_release_token_pool::ERebalancerCapDoesNotExist)]
public fun test_mcms_withdraw_liquidity_fails_without_rebalancer_cap() {
    let (mut env, owner_cap, rebalancer_cap) = setup();

    // Transfer ownership to MCMS but DON'T take rebalancer control
    setup_mcms_ownership(&mut env, owner_cap);

    // Try to withdraw liquidity without having rebalancer cap (should fail)
    env.scenario.next_tx(mcms_registry::get_multisig_address());
    let data = encode_withdraw_liquidity_data(
        object::id_address(&env.state),
        @0x0, // dummy rebalancer cap id
        100000,
        NEW_OWNER,
    );
    let params = create_mcms_callback_params(
        @lock_release_token_pool,
        b"withdraw_liquidity",
        data,
        x"0000000000000000000000000000000000000000000000000000000000000002",
        0,
    );
    lock_release_token_pool::mcms_withdraw_liquidity(
        &mut env.state,
        &mut env.mcms_registry,
        params,
        env.scenario.ctx(),
    );

    ts::return_to_address(EOA_REBALANCER, rebalancer_cap);
    tear_down(env);
}

#[test]
#[expected_failure(abort_code = lock_release_token_pool::ETokenPoolBalanceTooLow)]
public fun test_mcms_withdraw_liquidity_fails_insufficient_balance() {
    let (mut env, owner_cap, rebalancer_cap) = setup();
    let owner_cap_id = object::id_address(&owner_cap);

    // Setup MCMS with rebalancer control
    let rebalancer_cap_id = setup_mcms_with_rebalancer(&mut env, owner_cap, owner_cap_id);

    // Provide small amount of liquidity
    env.scenario.next_tx(OWNER);
    let mut treasury_cap = ts::take_from_sender<
        coin::TreasuryCap<LOCK_RELEASE_TOKEN_POOL_MCMS_CAP_TEST>,
    >(&env.scenario);
    let test_coin = coin::mint(&mut treasury_cap, 100000, env.scenario.ctx());
    let coin_id = object::id_address(&test_coin);

    env.scenario.next_tx(mcms_registry::get_multisig_address());
    let data = encode_provide_liquidity_data(
        object::id_address(&env.state),
        rebalancer_cap_id,
        coin_id,
    );
    let params = create_mcms_callback_params(
        @lock_release_token_pool,
        b"provide_liquidity",
        data,
        x"0000000000000000000000000000000000000000000000000000000000000003",
        0,
    );
    lock_release_token_pool::mcms_provide_liquidity(
        &mut env.state,
        &mut env.mcms_registry,
        test_coin,
        params,
        env.scenario.ctx(),
    );

    // Try to withdraw more than available (should fail)
    env.scenario.next_tx(mcms_registry::get_multisig_address());
    let data = encode_withdraw_liquidity_data(
        object::id_address(&env.state),
        rebalancer_cap_id,
        500000, // More than the 100000 provided
        NEW_OWNER,
    );
    let params = create_mcms_callback_params(
        @lock_release_token_pool,
        b"withdraw_liquidity",
        data,
        x"0000000000000000000000000000000000000000000000000000000000000004",
        0,
    );
    lock_release_token_pool::mcms_withdraw_liquidity(
        &mut env.state,
        &mut env.mcms_registry,
        params,
        env.scenario.ctx(),
    );

    ts::return_to_address(OWNER, treasury_cap);
    ts::return_to_address(EOA_REBALANCER, rebalancer_cap);
    tear_down(env);
}

#[test]
#[expected_failure(abort_code = mcms::bcs_stream::EInvalidObjectAddress)]
public fun test_mcms_provide_liquidity_fails_wrong_coin_id() {
    let (mut env, owner_cap, rebalancer_cap) = setup();
    let owner_cap_id = object::id_address(&owner_cap);

    // Setup MCMS with rebalancer control
    let rebalancer_cap_id = setup_mcms_with_rebalancer(&mut env, owner_cap, owner_cap_id);

    // Mint test coin but encode wrong coin ID
    env.scenario.next_tx(OWNER);
    let mut treasury_cap = ts::take_from_sender<
        coin::TreasuryCap<LOCK_RELEASE_TOKEN_POOL_MCMS_CAP_TEST>,
    >(&env.scenario);
    let test_coin = coin::mint(&mut treasury_cap, 1000000, env.scenario.ctx());

    // Provide liquidity with WRONG coin ID in BCS data (should fail during validation)
    env.scenario.next_tx(mcms_registry::get_multisig_address());
    let data = encode_provide_liquidity_data(
        object::id_address(&env.state),
        rebalancer_cap_id,
        @0x999, // Wrong coin ID
    );
    let params = create_mcms_callback_params(
        @lock_release_token_pool,
        b"provide_liquidity",
        data,
        x"0000000000000000000000000000000000000000000000000000000000000003",
        0,
    );
    lock_release_token_pool::mcms_provide_liquidity(
        &mut env.state,
        &mut env.mcms_registry,
        test_coin,
        params,
        env.scenario.ctx(),
    );

    ts::return_to_address(OWNER, treasury_cap);
    ts::return_to_address(EOA_REBALANCER, rebalancer_cap);
    tear_down(env);
}
