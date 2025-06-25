#[test_only]
module managed_token::managed_token_test;

use std::string;
use sui::coin::{Self, CoinMetadata};
use sui::deny_list;
use sui::test_scenario::{Self, Scenario};
use sui::test_utils;

use managed_token::managed_token::{Self, TokenState, MintCap};
use managed_token::ownable::OwnerCap;

// Test token types - each test gets its own unique witness type
public struct MANAGED_TOKEN_TEST has drop {}

const OWNER: address = @0x1;
const MINTER: address = @0x2;
const RECIPIENT: address = @0x3;
const OTHER_USER: address = @0x4;

// === Helper Functions ===

// Generic setup function that takes a witness type
fun setup_managed_token_test(): (Scenario, TokenState<MANAGED_TOKEN_TEST>, OwnerCap<MANAGED_TOKEN_TEST>, CoinMetadata<MANAGED_TOKEN_TEST>) {
    let mut scenario = test_scenario::begin(OWNER);
    let ctx = scenario.ctx();

    // Create the test token
    let (treasury_cap, coin_metadata) = coin::create_currency(
        MANAGED_TOKEN_TEST {},
        8,
        b"TEST",
        b"Test Token",
        b"A test token for managed token tests",
        option::none(),
        ctx
    );

    // Initialize managed token
    managed_token::initialize(treasury_cap, ctx);

    // Get the shared objects from scenario
    scenario.next_tx(OWNER);
    let state = scenario.take_shared<TokenState<MANAGED_TOKEN_TEST>>();
    let owner_cap = scenario.take_from_sender<OwnerCap<MANAGED_TOKEN_TEST>>();

    (scenario, state, owner_cap, coin_metadata)
}

fun setup_regulated_token_test(): (Scenario, TokenState<MANAGED_TOKEN_TEST>, OwnerCap<MANAGED_TOKEN_TEST>, CoinMetadata<MANAGED_TOKEN_TEST>) {
    let mut scenario = test_scenario::begin(OWNER);
    let ctx = scenario.ctx();

    // Create the test token with deny cap (required for pause/blocklist functionality)
    let (treasury_cap, deny_cap, coin_metadata) = coin::create_regulated_currency_v2(
        MANAGED_TOKEN_TEST {},
        8,
        b"TEST",
        b"Test Token",
        b"A test token for managed token tests",
        option::none(),
        true,
        ctx
    );

    // Initialize with deny cap
    managed_token::initialize_with_deny_cap(treasury_cap, deny_cap, ctx);

    scenario.next_tx(OWNER);
    let state = scenario.take_shared<TokenState<MANAGED_TOKEN_TEST>>();
    let owner_cap = scenario.take_from_sender<OwnerCap<MANAGED_TOKEN_TEST>>();

    (scenario, state, owner_cap, coin_metadata)
}

fun cleanup_managed_token_test<T>(
    scenario: Scenario,
    state: TokenState<T>,
    owner_cap: OwnerCap<T>,
    coin_metadata: CoinMetadata<T>
) {
    transfer::public_transfer(owner_cap, OWNER);
    transfer::public_freeze_object(coin_metadata);
    test_scenario::return_shared(state);
    scenario.end();
}

fun setup_minter_with_allowance(
    state: &mut TokenState<MANAGED_TOKEN_TEST>,
    owner_cap: &OwnerCap<MANAGED_TOKEN_TEST>,
    minter: address,
    allowance: u64,
    is_unlimited: bool,
    ctx: &mut TxContext
) {
    managed_token::configure_new_minter(
        state,
        owner_cap,
        minter,
        allowance,
        is_unlimited,
        ctx
    );
}

// === Basic Functionality Tests ===

#[test]
fun test_basic_initialization_and_configuration() {
    let (mut scenario, mut state, owner_cap, coin_metadata) = setup_managed_token_test();

    // Verify initial state
    assert!(managed_token::owner(&state) == OWNER);
    assert!(managed_token::total_supply(&state) == 0);

    // Test type and version
    let version = managed_token::type_and_version();
    assert!(version == string::utf8(b"ManagedToken 1.0.0"));

    // Configure multiple minters
    setup_minter_with_allowance(&mut state, &owner_cap, MINTER, 1000, false, scenario.ctx());
    setup_minter_with_allowance(&mut state, &owner_cap, OTHER_USER, 2000, true, scenario.ctx());

    // Verify mint caps were created
    let mint_caps = managed_token::get_all_mint_caps(&state);
    assert!(mint_caps.length() == 2);

    // Clean up mint caps
    scenario.next_tx(MINTER);
    let mint_cap1 = scenario.take_from_sender<MintCap<MANAGED_TOKEN_TEST>>();
    transfer::public_transfer(mint_cap1, MINTER);

    scenario.next_tx(OTHER_USER);
    let mint_cap2 = scenario.take_from_sender<MintCap<MANAGED_TOKEN_TEST>>();
    transfer::public_transfer(mint_cap2, OTHER_USER);

    cleanup_managed_token_test(scenario, state, owner_cap, coin_metadata);
}

#[test]
fun test_treasury_cap_access_and_total_supply() {
    let (mut scenario, mut state, owner_cap, coin_metadata) = setup_managed_token_test();

    // Test borrow_treasury_cap function
    let treasury_cap_ref = managed_token::borrow_treasury_cap(&owner_cap, &state);
    assert!(treasury_cap_ref.total_supply() == 0);
    assert!(managed_token::total_supply(&state) == 0);

    // Configure a minter and mint tokens to test supply tracking
    setup_minter_with_allowance(&mut state, &owner_cap, MINTER, 500, false, scenario.ctx());

    scenario.next_tx(MINTER);
    let mint_cap = scenario.take_from_sender<MintCap<MANAGED_TOKEN_TEST>>();
    let deny_list = deny_list::new_for_testing(scenario.ctx());
    
    // Mint tokens and verify supply tracking
    let amount = 300;
    let coin = managed_token::mint(&mut state, &mint_cap, &deny_list, amount, RECIPIENT, scenario.ctx());
    assert!(managed_token::total_supply(&state) == amount);
    
    let treasury_cap_ref = managed_token::borrow_treasury_cap(&owner_cap, &state);
    assert!(treasury_cap_ref.total_supply() == amount);

    // Clean up
    transfer::public_transfer(coin, RECIPIENT);
    transfer::public_transfer(mint_cap, MINTER);
    test_utils::destroy(deny_list);
    cleanup_managed_token_test(scenario, state, owner_cap, coin_metadata);
}

// === Minting Tests ===

#[test]
fun test_mint_operations_with_sufficient_allowance() {
    let (mut scenario, mut state, owner_cap, coin_metadata) = setup_managed_token_test();

    let allowance = 1000;
    setup_minter_with_allowance(&mut state, &owner_cap, MINTER, allowance, false, scenario.ctx());

    scenario.next_tx(MINTER);
    let mint_cap = scenario.take_from_sender<MintCap<MANAGED_TOKEN_TEST>>();
    let deny_list = deny_list::new_for_testing(scenario.ctx());
    
    // Test mint function
    let mint_amount = 500;
    let coin = managed_token::mint(&mut state, &mint_cap, &deny_list, mint_amount, RECIPIENT, scenario.ctx());
    assert!(coin.value() == mint_amount);
    assert!(managed_token::total_supply(&state) == mint_amount);

    // Verify remaining allowance
    let (remaining_allowance, is_unlimited) = managed_token::mint_allowance(&state, object::id(&mint_cap));
    assert!(remaining_allowance == allowance - mint_amount);
    assert!(!is_unlimited);

    // Test mint_and_transfer function
    let transfer_amount = 200;
    managed_token::mint_and_transfer(&mut state, &mint_cap, &deny_list, transfer_amount, RECIPIENT, scenario.ctx());
    assert!(managed_token::total_supply(&state) == mint_amount + transfer_amount);

    // Verify final allowance
    let (final_allowance, _) = managed_token::mint_allowance(&state, object::id(&mint_cap));
    assert!(final_allowance == allowance - mint_amount - transfer_amount);

    // Clean up
    transfer::public_transfer(coin, RECIPIENT);
    transfer::public_transfer(mint_cap, MINTER);
    test_utils::destroy(deny_list);
    cleanup_managed_token_test(scenario, state, owner_cap, coin_metadata);
}

#[test]
fun test_unlimited_allowance_operations() {
    let (mut scenario, mut state, owner_cap, coin_metadata) = setup_managed_token_test();

    // Configure minter with unlimited allowance
    setup_minter_with_allowance(&mut state, &owner_cap, MINTER, 0, true, scenario.ctx());

    scenario.next_tx(MINTER);
    let mint_cap = scenario.take_from_sender<MintCap<MANAGED_TOKEN_TEST>>();
    let deny_list = deny_list::new_for_testing(scenario.ctx());

    // Verify unlimited allowance
    let (_allowance, is_unlimited) = managed_token::mint_allowance(&state, object::id(&mint_cap));
    assert!(is_unlimited);

    // Test minting large amounts
    let amounts = vector[1000, 5000, 10000];
    let mut total_minted = 0;

    let mut i = 0;
    while (i < amounts.length()) {
        let amount = amounts[i];
        let coin = managed_token::mint(&mut state, &mint_cap, &deny_list, amount, RECIPIENT, scenario.ctx());
        assert!(coin.value() == amount);
        total_minted = total_minted + amount;
        
        // Test burning
        managed_token::burn(&mut state, &mint_cap, &deny_list, coin, RECIPIENT, scenario.ctx());
        i = i + 1;
    };

    // Verify final supply is zero after burns
    assert!(managed_token::total_supply(&state) == 0);

    // Clean up
    transfer::public_transfer(mint_cap, MINTER);
    test_utils::destroy(deny_list);
    cleanup_managed_token_test(scenario, state, owner_cap, coin_metadata);
}

// === Allowance Management Tests ===

#[test]
fun test_allowance_management_comprehensive() {
    let (mut scenario, mut state, owner_cap, coin_metadata) = setup_managed_token_test();

    let initial_allowance = 1000;
    setup_minter_with_allowance(&mut state, &owner_cap, MINTER, initial_allowance, false, scenario.ctx());

    scenario.next_tx(MINTER);
    let mint_cap = scenario.take_from_sender<MintCap<MANAGED_TOKEN_TEST>>();
    let mint_cap_id = object::id(&mint_cap);
    let deny_list = deny_list::new_for_testing(scenario.ctx());

    // Test increment_mint_allowance
    scenario.next_tx(OWNER);
    let increment_amount = 500;
    managed_token::increment_mint_allowance(&mut state, &owner_cap, mint_cap_id, &deny_list, increment_amount, scenario.ctx());

    let (new_allowance, is_unlimited) = managed_token::mint_allowance(&state, mint_cap_id);
    assert!(new_allowance == initial_allowance + increment_amount);
    assert!(!is_unlimited);

    // Test set_unlimited_mint_allowances to unlimited
    managed_token::set_unlimited_mint_allowances(&mut state, &owner_cap, mint_cap_id, &deny_list, true, scenario.ctx());
    let (_allowance, is_unlimited) = managed_token::mint_allowance(&state, mint_cap_id);
    assert!(is_unlimited);

    // Test set back to limited
    managed_token::set_unlimited_mint_allowances(&mut state, &owner_cap, mint_cap_id, &deny_list, false, scenario.ctx());
    let (limited_allowance, is_unlimited) = managed_token::mint_allowance(&state, mint_cap_id);
    assert!(limited_allowance == 0);
    assert!(!is_unlimited);

    // Clean up
    transfer::public_transfer(mint_cap, MINTER);
    test_utils::destroy(deny_list);
    cleanup_managed_token_test(scenario, state, owner_cap, coin_metadata);
}

#[test]
fun test_multiple_allowance_increments() {
    let (mut scenario, mut state, owner_cap, coin_metadata) = setup_managed_token_test();

    let initial_allowance = 1000;
    setup_minter_with_allowance(&mut state, &owner_cap, MINTER, initial_allowance, false, scenario.ctx());

    scenario.next_tx(MINTER);
    let mint_cap = scenario.take_from_sender<MintCap<MANAGED_TOKEN_TEST>>();
    let mint_cap_id = object::id(&mint_cap);
    let deny_list = deny_list::new_for_testing(scenario.ctx());

    // Multiple increments
    let increments = vector[500, 300, 200];
    let mut expected_allowance = initial_allowance;

    scenario.next_tx(OWNER);
    let mut i = 0;
    while (i < increments.length()) {
        let increment = increments[i];
        managed_token::increment_mint_allowance(&mut state, &owner_cap, mint_cap_id, &deny_list, increment, scenario.ctx());
        expected_allowance = expected_allowance + increment;
        
        let (current_allowance, is_unlimited) = managed_token::mint_allowance(&state, mint_cap_id);
        assert!(current_allowance == expected_allowance);
        assert!(!is_unlimited);
        i = i + 1;
    };

    // Test minting with final allowance
    scenario.next_tx(MINTER);
    let mint_amount = 1500; // More than initial but less than final allowance
    let coin = managed_token::mint(&mut state, &mint_cap, &deny_list, mint_amount, RECIPIENT, scenario.ctx());
    assert!(coin.value() == mint_amount);

    // Clean up
    transfer::public_transfer(coin, RECIPIENT);
    transfer::public_transfer(mint_cap, MINTER);
    test_utils::destroy(deny_list);
    cleanup_managed_token_test(scenario, state, owner_cap, coin_metadata);
}

// === Ownership Tests ===

#[test]
fun test_ownership_transfer_flow() {
    let (mut scenario, mut state, owner_cap, coin_metadata) = setup_managed_token_test();

    // Verify initial owner
    assert!(managed_token::owner(&state) == OWNER);

    // Initiate ownership transfer
    managed_token::transfer_ownership(&mut state, &owner_cap, OTHER_USER, scenario.ctx());

    // Verify pending transfer state
    assert!(managed_token::has_pending_transfer(&state));
    assert!(managed_token::pending_transfer_from(&state) == option::some(OWNER));
    assert!(managed_token::pending_transfer_to(&state) == option::some(OTHER_USER));

    // Accept ownership transfer
    scenario.next_tx(OTHER_USER);
    managed_token::accept_ownership(&mut state, scenario.ctx());

    cleanup_managed_token_test(scenario, state, owner_cap, coin_metadata);
}

// === Regulated Token Tests (Pause/Blocklist) ===

#[test]
fun test_pause_and_unpause_operations() {
    let (mut scenario, mut state, owner_cap, coin_metadata) = setup_regulated_token_test();
    let mut deny_list = deny_list::new_for_testing(scenario.ctx());

    // Test pause function
    managed_token::pause(&mut state, &owner_cap, &mut deny_list, scenario.ctx());

    // Test unpause function
    managed_token::unpause(&mut state, &owner_cap, &mut deny_list, scenario.ctx());

    // Clean up
    test_utils::destroy(deny_list);
    cleanup_managed_token_test(scenario, state, owner_cap, coin_metadata);
}

#[test]
fun test_blocklist_and_unblocklist_operations() {
    let (mut scenario, mut state, owner_cap, coin_metadata) = setup_regulated_token_test();
    let mut deny_list = deny_list::new_for_testing(scenario.ctx());

    let address_to_block = @0x123;

    // Test blocklist function
    managed_token::blocklist(&mut state, &owner_cap, &mut deny_list, address_to_block, scenario.ctx());

    // Test unblocklist function
    managed_token::unblocklist(&mut state, &owner_cap, &mut deny_list, address_to_block, scenario.ctx());

    // Clean up
    test_utils::destroy(deny_list);
    cleanup_managed_token_test(scenario, state, owner_cap, coin_metadata);
}

// === Authorization and Access Control Tests ===

#[test]
fun test_unauthorized_access_scenarios() {
    let (mut scenario, mut state, owner_cap, coin_metadata) = setup_managed_token_test();

    // Configure a minter
    setup_minter_with_allowance(&mut state, &owner_cap, MINTER, 1000, false, scenario.ctx());

    scenario.next_tx(MINTER);
    let mint_cap = scenario.take_from_sender<MintCap<MANAGED_TOKEN_TEST>>();
    let mint_cap_id = object::id(&mint_cap);

    // Test unauthorized mint cap queries
    let fake_mint_cap_id = object::id_from_address(@0x999);
    let (allowance, is_unlimited) = managed_token::mint_allowance(&state, fake_mint_cap_id);
    assert!(allowance == 0);
    assert!(!is_unlimited);
    assert!(!managed_token::is_authorized_mint_cap(&state, fake_mint_cap_id));

    // Verify authorized mint cap works
    assert!(managed_token::is_authorized_mint_cap(&state, mint_cap_id));
    let (real_allowance, real_unlimited) = managed_token::mint_allowance(&state, mint_cap_id);
    assert!(real_allowance == 1000);
    assert!(!real_unlimited);

    // Clean up
    transfer::public_transfer(mint_cap, MINTER);
    cleanup_managed_token_test(scenario, state, owner_cap, coin_metadata);
}

// === Burn Operations Tests ===

#[test]
fun test_burn_operations_comprehensive() {
    let (mut scenario, mut state, owner_cap, coin_metadata) = setup_managed_token_test();

    // Configure minters
    setup_minter_with_allowance(&mut state, &owner_cap, MINTER, 1000, false, scenario.ctx());
    setup_minter_with_allowance(&mut state, &owner_cap, OTHER_USER, 1000, false, scenario.ctx());

    // Get mint caps
    scenario.next_tx(MINTER);
    let mint_cap1 = scenario.take_from_sender<MintCap<MANAGED_TOKEN_TEST>>();

    scenario.next_tx(OTHER_USER);
    let mint_cap2 = scenario.take_from_sender<MintCap<MANAGED_TOKEN_TEST>>();

    let deny_list = deny_list::new_for_testing(scenario.ctx());

    // Both minters mint tokens
    scenario.next_tx(MINTER);
    let coin1 = managed_token::mint(&mut state, &mint_cap1, &deny_list, 300, MINTER, scenario.ctx());

    scenario.next_tx(OTHER_USER);
    let coin2 = managed_token::mint(&mut state, &mint_cap2, &deny_list, 400, OTHER_USER, scenario.ctx());

    // Verify total supply
    assert!(managed_token::total_supply(&state) == 700);

    // Sequential burns
    scenario.next_tx(MINTER);
    managed_token::burn(&mut state, &mint_cap1, &deny_list, coin1, MINTER, scenario.ctx());
    assert!(managed_token::total_supply(&state) == 400);

    scenario.next_tx(OTHER_USER);
    managed_token::burn(&mut state, &mint_cap2, &deny_list, coin2, OTHER_USER, scenario.ctx());
    assert!(managed_token::total_supply(&state) == 0);

    // Clean up
    transfer::public_transfer(mint_cap1, MINTER);
    transfer::public_transfer(mint_cap2, OTHER_USER);
    test_utils::destroy(deny_list);
    cleanup_managed_token_test(scenario, state, owner_cap, coin_metadata);
}

// === Destruction Tests ===

#[test]
fun test_managed_token_destruction_scenarios() {
    let (mut scenario, mut state, owner_cap, coin_metadata) = setup_managed_token_test();

    // Configure multiple minters
    setup_minter_with_allowance(&mut state, &owner_cap, MINTER, 1000, false, scenario.ctx());
    setup_minter_with_allowance(&mut state, &owner_cap, OTHER_USER, 0, true, scenario.ctx());

    // Get mint caps
    scenario.next_tx(MINTER);
    let mint_cap1 = scenario.take_from_sender<MintCap<MANAGED_TOKEN_TEST>>();

    scenario.next_tx(OTHER_USER);
    let mint_cap2 = scenario.take_from_sender<MintCap<MANAGED_TOKEN_TEST>>();

    // Verify mint caps exist
    let mint_caps = managed_token::get_all_mint_caps(&state);
    assert!(mint_caps.length() == 2);

    // Destroy the managed token
    scenario.next_tx(OWNER);
    let (treasury_cap, deny_cap_option) = managed_token::destroy_managed_token(owner_cap, state, scenario.ctx());

    // Verify results
    assert!(deny_cap_option.is_none());
    assert!(treasury_cap.total_supply() == 0);

    // Clean up
    transfer::public_transfer(mint_cap1, MINTER);
    transfer::public_transfer(mint_cap2, OTHER_USER);
    transfer::public_transfer(treasury_cap, OWNER);
    transfer::public_freeze_object(coin_metadata);
    deny_cap_option.destroy_none();
    scenario.end();
}

#[test]
fun test_destruction_with_regulated_token() {
    let (mut scenario, mut state, owner_cap, coin_metadata) = setup_regulated_token_test();

    // Configure a minter and mint some tokens
    setup_minter_with_allowance(&mut state, &owner_cap, MINTER, 1500, false, scenario.ctx());

    scenario.next_tx(MINTER);
    let mint_cap = scenario.take_from_sender<MintCap<MANAGED_TOKEN_TEST>>();
    let deny_list = deny_list::new_for_testing(scenario.ctx());

    // Mint tokens to increase supply
    let mint_amount = 500;
    let coin = managed_token::mint(&mut state, &mint_cap, &deny_list, mint_amount, RECIPIENT, scenario.ctx());
    assert!(managed_token::total_supply(&state) == mint_amount);

    // Destroy the managed token
    scenario.next_tx(OWNER);
    let (treasury_cap, deny_cap_option) = managed_token::destroy_managed_token(owner_cap, state, scenario.ctx());

    // Verify results
    assert!(deny_cap_option.is_some());
    assert!(treasury_cap.total_supply() == mint_amount);
    assert!(coin.value() == mint_amount);

    // Clean up
    let mut deny_cap_option = deny_cap_option;
    let deny_cap = deny_cap_option.extract();
    transfer::public_transfer(coin, RECIPIENT);
    transfer::public_transfer(mint_cap, MINTER);
    transfer::public_transfer(treasury_cap, OWNER);
    transfer::public_transfer(deny_cap, OWNER);
    transfer::public_freeze_object(coin_metadata);
    test_utils::destroy(deny_list);
    deny_cap_option.destroy_none();
    scenario.end();
}

// === Error Condition Tests ===

#[test]
#[expected_failure(abort_code = managed_token::EZeroAmount)]
fun test_zero_amount_mint_failure() {
    let (mut scenario, mut state, owner_cap, coin_metadata) = setup_managed_token_test();
    setup_minter_with_allowance(&mut state, &owner_cap, MINTER, 1000, false, scenario.ctx());

    scenario.next_tx(MINTER);
    let mint_cap = scenario.take_from_sender<MintCap<MANAGED_TOKEN_TEST>>();
    let deny_list = deny_list::new_for_testing(scenario.ctx());
    
    // Try to mint zero amount - should fail
    let coin = managed_token::mint(&mut state, &mint_cap, &deny_list, 0, RECIPIENT, scenario.ctx());
    transfer::public_transfer(coin, RECIPIENT);

    // Cleanup (should not be reached)
    transfer::public_transfer(mint_cap, MINTER);
    test_utils::destroy(deny_list);
    cleanup_managed_token_test(scenario, state, owner_cap, coin_metadata);
}

#[test]
#[expected_failure(abort_code = managed_token::EInsufficientAllowance)]
fun test_insufficient_allowance_failure() {
    let (mut scenario, mut state, owner_cap, coin_metadata) = setup_managed_token_test();
    
    let allowance = 500;
    setup_minter_with_allowance(&mut state, &owner_cap, MINTER, allowance, false, scenario.ctx());

    scenario.next_tx(MINTER);
    let mint_cap = scenario.take_from_sender<MintCap<MANAGED_TOKEN_TEST>>();
    let deny_list = deny_list::new_for_testing(scenario.ctx());

    // Try to mint more than allowance - should fail
    managed_token::mint_and_transfer(&mut state, &mint_cap, &deny_list, allowance + 1, RECIPIENT, scenario.ctx());

    // Cleanup (should not be reached)
    transfer::public_transfer(mint_cap, MINTER);
    test_utils::destroy(deny_list);
    cleanup_managed_token_test(scenario, state, owner_cap, coin_metadata);
}

#[test]
#[expected_failure(abort_code = managed_token::EZeroAmount)]
fun test_zero_amount_burn_failure() {
    let (mut scenario, mut state, owner_cap, coin_metadata) = setup_managed_token_test();
    setup_minter_with_allowance(&mut state, &owner_cap, MINTER, 1000, false, scenario.ctx());

    scenario.next_tx(MINTER);
    let mint_cap = scenario.take_from_sender<MintCap<MANAGED_TOKEN_TEST>>();
    let deny_list = deny_list::new_for_testing(scenario.ctx());

    // Create zero-value coin and try to burn - should fail
    let zero_coin = coin::zero<MANAGED_TOKEN_TEST>(scenario.ctx());
    managed_token::burn(&mut state, &mint_cap, &deny_list, zero_coin, MINTER, scenario.ctx());

    // Cleanup (should not be reached)
    transfer::public_transfer(mint_cap, MINTER);
    test_utils::destroy(deny_list);
    cleanup_managed_token_test(scenario, state, owner_cap, coin_metadata);
}

#[test]
#[expected_failure(abort_code = managed_token::EZeroAmount)]
fun test_increment_allowance_zero_amount_failure() {
    let (mut scenario, mut state, owner_cap, coin_metadata) = setup_managed_token_test();
    setup_minter_with_allowance(&mut state, &owner_cap, MINTER, 1000, false, scenario.ctx());

    scenario.next_tx(MINTER);
    let mint_cap = scenario.take_from_sender<MintCap<MANAGED_TOKEN_TEST>>();
    let mint_cap_id = object::id(&mint_cap);
    let deny_list = deny_list::new_for_testing(scenario.ctx());

    // Try to increment with zero amount - should fail
    scenario.next_tx(OWNER);
    managed_token::increment_mint_allowance(&mut state, &owner_cap, mint_cap_id, &deny_list, 0, scenario.ctx());

    // Cleanup (should not be reached)
    transfer::public_transfer(mint_cap, MINTER);
    test_utils::destroy(deny_list);
    cleanup_managed_token_test(scenario, state, owner_cap, coin_metadata);
}

#[test]
#[expected_failure(abort_code = managed_token::ECannotIncreaseUnlimitedAllowance)]
fun test_increment_unlimited_allowance_failure() {
    let (mut scenario, mut state, owner_cap, coin_metadata) = setup_managed_token_test();
    setup_minter_with_allowance(&mut state, &owner_cap, MINTER, 0, true, scenario.ctx());

    scenario.next_tx(MINTER);
    let mint_cap = scenario.take_from_sender<MintCap<MANAGED_TOKEN_TEST>>();
    let mint_cap_id = object::id(&mint_cap);
    let deny_list = deny_list::new_for_testing(scenario.ctx());

    // Try to increment unlimited allowance - should fail
    scenario.next_tx(OWNER);
    managed_token::increment_mint_allowance(&mut state, &owner_cap, mint_cap_id, &deny_list, 500, scenario.ctx());

    // Cleanup (should not be reached)
    transfer::public_transfer(mint_cap, MINTER);
    test_utils::destroy(deny_list);
    cleanup_managed_token_test(scenario, state, owner_cap, coin_metadata);
}

#[test]
#[expected_failure(abort_code = managed_token::EUnauthorizedMintCap)]
fun test_unauthorized_mint_cap_operations() {
    let (mut scenario, mut state, owner_cap, coin_metadata) = setup_managed_token_test();
    let deny_list = deny_list::new_for_testing(scenario.ctx());

    // Use fake mint cap ID
    let fake_mint_cap_id = object::id_from_address(@0x999);

    // Try various operations with unauthorized mint cap - should fail
    scenario.next_tx(OWNER);
    managed_token::increment_mint_allowance(&mut state, &owner_cap, fake_mint_cap_id, &deny_list, 500, scenario.ctx());

    // Cleanup (should not be reached)
    test_utils::destroy(deny_list);
    cleanup_managed_token_test(scenario, state, owner_cap, coin_metadata);
}

#[test]
#[expected_failure(abort_code = managed_token::EPaused)]
fun test_operations_while_paused() {
    let (mut scenario, mut state, owner_cap, coin_metadata) = setup_regulated_token_test();
    let mut deny_list = deny_list::new_for_testing(scenario.ctx());

    setup_minter_with_allowance(&mut state, &owner_cap, MINTER, 1000, false, scenario.ctx());

    scenario.next_tx(MINTER);
    let mint_cap = scenario.take_from_sender<MintCap<MANAGED_TOKEN_TEST>>();
    let mint_cap_id = object::id(&mint_cap);

    // Pause the token
    scenario.next_tx(OWNER);
    managed_token::pause(&mut state, &owner_cap, &mut deny_list, scenario.ctx());

    // Try to increment allowance while paused - should fail
    managed_token::increment_mint_allowance(&mut state, &owner_cap, mint_cap_id, &deny_list, 500, scenario.ctx());

    // Cleanup (should not be reached)
    transfer::public_transfer(mint_cap, MINTER);
    test_utils::destroy(deny_list);
    cleanup_managed_token_test(scenario, state, owner_cap, coin_metadata);
}

#[test]
#[expected_failure(abort_code = managed_token::EDeniedAddress)]
fun test_operations_with_blocklisted_addresses() {
    let (mut scenario, mut state, owner_cap, coin_metadata) = setup_regulated_token_test();
    let mut deny_list = deny_list::new_for_testing(scenario.ctx());

    setup_minter_with_allowance(&mut state, &owner_cap, MINTER, 1000, false, scenario.ctx());

    scenario.next_tx(MINTER);
    let mint_cap = scenario.take_from_sender<MintCap<MANAGED_TOKEN_TEST>>();

    // Blocklist the minter
    scenario.next_tx(OWNER);
    managed_token::blocklist(&mut state, &owner_cap, &mut deny_list, MINTER, scenario.ctx());

    // Try to mint while minter is blocklisted - should fail
    scenario.next_tx(MINTER);
    managed_token::mint_and_transfer(&mut state, &mint_cap, &deny_list, 500, RECIPIENT, scenario.ctx());

    // Cleanup (should not be reached)
    transfer::public_transfer(mint_cap, MINTER);
    test_utils::destroy(deny_list);
    cleanup_managed_token_test(scenario, state, owner_cap, coin_metadata);
}

#[test]
#[expected_failure(abort_code = managed_token::EInvalidOwnerCap)]
fun test_invalid_owner_cap_operations() {
    let (mut scenario, mut state, owner_cap, coin_metadata) = setup_managed_token_test();

    // Create another managed token to get different owner cap
    scenario.next_tx(OTHER_USER);
    let ctx = scenario.ctx();
    let (fake_treasury_cap, fake_coin_metadata) = coin::create_currency(
        MANAGED_TOKEN_TEST {},
        8,
        b"FAKE",
        b"Fake Token",
        b"A fake token",
        option::none(),
        ctx
    );
    managed_token::initialize(fake_treasury_cap, ctx);

    scenario.next_tx(OTHER_USER);
    let fake_state = scenario.take_shared<TokenState<MANAGED_TOKEN_TEST>>();
    let fake_owner_cap = scenario.take_from_sender<OwnerCap<MANAGED_TOKEN_TEST>>();

         // Try to destroy with wrong owner cap - should fail
     let (treasury_cap_result, mut deny_cap_result) = managed_token::destroy_managed_token(fake_owner_cap, state, scenario.ctx());
     
     // This code should not be reached due to expected failure
     transfer::public_transfer(treasury_cap_result, OTHER_USER);
     if (deny_cap_result.is_some()) {
         transfer::public_transfer(deny_cap_result.extract(), OTHER_USER);
     };
     deny_cap_result.destroy_none();

    // Cleanup (should not be reached)
    transfer::public_transfer(owner_cap, OWNER);
    transfer::public_freeze_object(coin_metadata);
    transfer::public_freeze_object(fake_coin_metadata);
    test_scenario::return_shared(fake_state);
    scenario.end();
}
