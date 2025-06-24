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

#[test]
fun test_mint_with_sufficient_allowance() {
    let (mut scenario, mut state, owner_cap, coin_metadata) = setup_managed_token_test();

    // Configure a new minter with allowance
    let allowance = 1000;
    managed_token::configure_new_minter(
        &mut state,
        &owner_cap,
        MINTER,
        allowance,
        false, // not unlimited
        scenario.ctx()
    );

    // Switch to minter
    scenario.next_tx(MINTER);
    let mint_cap = scenario.take_from_sender<MintCap<MANAGED_TOKEN_TEST>>();

    // For testing without deny list, we'll create a dummy deny list
    let deny_list = deny_list::new_for_testing(scenario.ctx());
    
    // Mint tokens
    let amount = 500;
    let coin = managed_token::mint(
        &mut state,
        &mint_cap,
        &deny_list,
        amount,
        RECIPIENT,
        scenario.ctx()
    );

    // Verify the coin amount
    assert!(coin.value() == amount);

    // Verify total supply increased
    assert!(managed_token::total_supply(&state) == amount);

    // Verify remaining allowance
    let (remaining_allowance, is_unlimited) = managed_token::mint_allowance(&state, object::id(&mint_cap));
    assert!(remaining_allowance == allowance - amount);
    assert!(!is_unlimited);

    // Clean up
    transfer::public_transfer(coin, RECIPIENT);
    transfer::public_transfer(mint_cap, MINTER);
    test_utils::destroy(deny_list);
    cleanup_managed_token_test(scenario, state, owner_cap, coin_metadata);
}

#[test]
fun test_mint_with_insufficient_allowance() {
    let (mut scenario, mut state, owner_cap, coin_metadata) = setup_managed_token_test();

    // Configure a new minter with small allowance
    let allowance = 100;
    managed_token::configure_new_minter(
        &mut state,
        &owner_cap,
        MINTER,
        allowance,
        false,
        scenario.ctx()
    );

    // Switch to minter
    scenario.next_tx(MINTER);
    let mint_cap = scenario.take_from_sender<MintCap<MANAGED_TOKEN_TEST>>();

    // Verify allowance is set correctly
    let (current_allowance, is_unlimited) = managed_token::mint_allowance(&state, object::id(&mint_cap));
    assert!(current_allowance == allowance);
    assert!(!is_unlimited);

    // Clean up
    transfer::public_transfer(mint_cap, MINTER);
    cleanup_managed_token_test(scenario, state, owner_cap, coin_metadata);
}

#[test]
fun test_basic_configuration() {
    let (mut scenario, mut state, owner_cap, coin_metadata) = setup_managed_token_test();

    // Verify initial owner
    assert!(managed_token::owner(&state) == OWNER);

    // Configure multiple minters
    managed_token::configure_new_minter(
        &mut state,
        &owner_cap,
        MINTER,
        1000,
        false,
        scenario.ctx()
    );

    managed_token::configure_new_minter(
        &mut state,
        &owner_cap,
        OTHER_USER,
        2000,
        true, // unlimited
        scenario.ctx()
    );

    // Get all mint caps
    let mint_caps = managed_token::get_all_mint_caps(&state);
    assert!(mint_caps.length() == 2);

    // Verify type and version
    let version = managed_token::type_and_version();
    assert!(version == string::utf8(b"ManagedToken 1.0.0"));

    // Clean up - collect mint caps
    scenario.next_tx(MINTER);
    let mint_cap1 = scenario.take_from_sender<MintCap<MANAGED_TOKEN_TEST>>();
    transfer::public_transfer(mint_cap1, MINTER);

    scenario.next_tx(OTHER_USER);
    let mint_cap2 = scenario.take_from_sender<MintCap<MANAGED_TOKEN_TEST>>();
    transfer::public_transfer(mint_cap2, OTHER_USER);

    cleanup_managed_token_test(scenario, state, owner_cap, coin_metadata);
}

#[test]
fun test_unlimited_allowance_and_burn() {
    let (mut scenario, mut state, owner_cap, coin_metadata) = setup_managed_token_test();

    // Configure a new minter with unlimited allowance
    managed_token::configure_new_minter(
        &mut state,
        &owner_cap,
        MINTER,
        0, // amount doesn't matter for unlimited
        true, // unlimited
        scenario.ctx()
    );

    // Switch to minter
    scenario.next_tx(MINTER);
    let mint_cap = scenario.take_from_sender<MintCap<MANAGED_TOKEN_TEST>>();

    // Verify unlimited allowance
    let (_allowance, is_unlimited) = managed_token::mint_allowance(&state, object::id(&mint_cap));
    assert!(is_unlimited);

    // Create a dummy deny list for testing
    let deny_list = deny_list::new_for_testing(scenario.ctx());

    // Mint tokens
    let amount = 1000;
    let coin = managed_token::mint(
        &mut state,
        &mint_cap,
        &deny_list,
        amount,
        RECIPIENT,
        scenario.ctx()
    );

    assert!(coin.value() == amount);
    assert!(managed_token::total_supply(&state) == amount);

    // Test burning
    managed_token::burn(
        &mut state,
        &mint_cap,
        &deny_list,
        coin,
        RECIPIENT,
        scenario.ctx()
    );

    // Verify total supply decreased
    assert!(managed_token::total_supply(&state) == 0);

    // Verify type and version
    let version = managed_token::type_and_version();
    assert!(version == string::utf8(b"ManagedToken 1.0.0"));

    // Clean up
    transfer::public_transfer(mint_cap, MINTER);
    test_utils::destroy(deny_list);
    cleanup_managed_token_test(scenario, state, owner_cap, coin_metadata);
}

#[test]
fun test_ownership_transfer() {
    let (mut scenario, mut state, owner_cap, coin_metadata) = setup_managed_token_test();

    // Verify initial owner
    assert!(managed_token::owner(&state) == OWNER);

    // Test ownership transfer initiation
    managed_token::transfer_ownership(
        &mut state,
        &owner_cap,
        OTHER_USER,
        scenario.ctx()
    );

    // Verify pending transfer
    assert!(managed_token::has_pending_transfer(&state));
    assert!(managed_token::pending_transfer_from(&state) == option::some(OWNER));
    assert!(managed_token::pending_transfer_to(&state) == option::some(OTHER_USER));

    // Accept ownership transfer
    scenario.next_tx(OTHER_USER);
    managed_token::accept_ownership(&mut state, scenario.ctx());

    // Clean up
    cleanup_managed_token_test(scenario, state, owner_cap, coin_metadata);
}

#[test]
fun test_mint_and_transfer_function() {
    let (mut scenario, mut state, owner_cap, coin_metadata) = setup_managed_token_test();

    // Configure a new minter with allowance
    let allowance = 2000;
    managed_token::configure_new_minter(
        &mut state,
        &owner_cap,
        MINTER,
        allowance,
        false,
        scenario.ctx()
    );

    // Switch to minter
    scenario.next_tx(MINTER);
    let mint_cap = scenario.take_from_sender<MintCap<MANAGED_TOKEN_TEST>>();

    // Create deny list for testing
    let deny_list = deny_list::new_for_testing(scenario.ctx());
    
    // Test mint_and_transfer function (directly transfers to recipient)
    let amount = 750;
    managed_token::mint_and_transfer(
        &mut state,
        &mint_cap,
        &deny_list,
        amount,
        RECIPIENT,
        scenario.ctx()
    );

    // Verify total supply increased
    assert!(managed_token::total_supply(&state) == amount);

    // Verify remaining allowance decreased
    let (remaining_allowance, is_unlimited) = managed_token::mint_allowance(&state, object::id(&mint_cap));
    assert!(remaining_allowance == allowance - amount);
    assert!(!is_unlimited);

    // Verify recipient received the tokens
    scenario.next_tx(RECIPIENT);
    let coin = scenario.take_from_sender<sui::coin::Coin<MANAGED_TOKEN_TEST>>();
    assert!(coin.value() == amount);

    // Clean up
    transfer::public_transfer(coin, RECIPIENT);
    transfer::public_transfer(mint_cap, MINTER);
    test_utils::destroy(deny_list);
    cleanup_managed_token_test(scenario, state, owner_cap, coin_metadata);
}

#[test]
fun test_allowance_management_functions() {
    let (mut scenario, mut state, owner_cap, coin_metadata) = setup_managed_token_test();

    // Configure a new minter with limited allowance
    let initial_allowance = 1000;
    managed_token::configure_new_minter(
        &mut state,
        &owner_cap,
        MINTER,
        initial_allowance,
        false, // not unlimited initially
        scenario.ctx()
    );

    scenario.next_tx(MINTER);
    let mint_cap = scenario.take_from_sender<MintCap<MANAGED_TOKEN_TEST>>();
    let mint_cap_id = object::id(&mint_cap);

    // Create deny list for testing
    let deny_list = deny_list::new_for_testing(scenario.ctx());

    // Test set_unlimited_mint_allowances function
    scenario.next_tx(OWNER);
    managed_token::set_unlimited_mint_allowances(
        &mut state,
        &owner_cap,
        mint_cap_id,
        &deny_list,
        true, // set to unlimited
        scenario.ctx()
    );

    // Verify allowance is now unlimited
    let (_allowance, is_unlimited) = managed_token::mint_allowance(&state, mint_cap_id);
    assert!(is_unlimited);

    // Test that we can mint more than the original allowance
    scenario.next_tx(MINTER);
    let large_amount = 5000; // More than initial allowance
    let coin = managed_token::mint(
        &mut state,
        &mint_cap,
        &deny_list,
        large_amount,
        RECIPIENT,
        scenario.ctx()
    );

    assert!(coin.value() == large_amount);
    assert!(managed_token::total_supply(&state) == large_amount);

    // Clean up
    transfer::public_transfer(coin, RECIPIENT);
    transfer::public_transfer(mint_cap, MINTER);
    test_utils::destroy(deny_list);
    cleanup_managed_token_test(scenario, state, owner_cap, coin_metadata);
}

#[test]
fun test_unauthorized_access_errors() {
    let (mut scenario, mut state, owner_cap, coin_metadata) = setup_managed_token_test();

    // Configure a minter
    managed_token::configure_new_minter(
        &mut state,
        &owner_cap,
        MINTER,
        1000,
        false,
        scenario.ctx()
    );

    scenario.next_tx(MINTER);
    let mint_cap = scenario.take_from_sender<MintCap<MANAGED_TOKEN_TEST>>();
    let mint_cap_id = object::id(&mint_cap);

    // Create deny list
    let deny_list = deny_list::new_for_testing(scenario.ctx());

    // Test unauthorized mint cap usage by using a random ID
    let fake_mint_cap_id = object::id_from_address(@0x999);
    
    // Verify mint_allowance returns (0, false) for unauthorized mint cap
    let (allowance, is_unlimited) = managed_token::mint_allowance(&state, fake_mint_cap_id);
    assert!(allowance == 0);
    assert!(!is_unlimited);

    // Verify is_authorized_mint_cap returns false for unauthorized mint cap
    assert!(!managed_token::is_authorized_mint_cap(&state, fake_mint_cap_id));

    // Test get_all_mint_caps function
    let all_mint_caps = managed_token::get_all_mint_caps(&state);
    assert!(all_mint_caps.length() == 1);
    assert!(all_mint_caps[0] == mint_cap_id);

    // Clean up
    transfer::public_transfer(mint_cap, MINTER);
    test_utils::destroy(deny_list);
    cleanup_managed_token_test(scenario, state, owner_cap, coin_metadata);
}

#[test]
fun test_treasury_cap_access_and_validation() {
    let (mut scenario, mut state, owner_cap, coin_metadata) = setup_managed_token_test();

    // Test borrow_treasury_cap function
    let treasury_cap_ref = managed_token::borrow_treasury_cap(&owner_cap, &state);
    
    // Verify we can read treasury cap properties
    let initial_supply = treasury_cap_ref.total_supply();
    assert!(initial_supply == 0); // Should start with 0 supply

    // Also test the total_supply function
    assert!(managed_token::total_supply(&state) == initial_supply);

    // Test owner function
    assert!(managed_token::owner(&state) == OWNER);

    // Configure a minter and mint some tokens to test supply tracking
    managed_token::configure_new_minter(
        &mut state,
        &owner_cap,
        MINTER,
        500,
        false,
        scenario.ctx()
    );

    scenario.next_tx(MINTER);
    let mint_cap = scenario.take_from_sender<MintCap<MANAGED_TOKEN_TEST>>();

    let deny_list = deny_list::new_for_testing(scenario.ctx());
    
    // Mint some tokens
    let amount = 300;
    let coin = managed_token::mint(
        &mut state,
        &mint_cap,
        &deny_list,
        amount,
        RECIPIENT,
        scenario.ctx()
    );

    // Verify supply increased
    assert!(managed_token::total_supply(&state) == amount);
    
    // Verify treasury cap also reflects the change
    let treasury_cap_ref = managed_token::borrow_treasury_cap(&owner_cap, &state);
    assert!(treasury_cap_ref.total_supply() == amount);

    // Clean up
    transfer::public_transfer(coin, RECIPIENT);
    transfer::public_transfer(mint_cap, MINTER);
    test_utils::destroy(deny_list);
    cleanup_managed_token_test(scenario, state, owner_cap, coin_metadata);
}

#[test]
#[expected_failure(abort_code = managed_token::EZeroAmount)]
fun test_zero_amount_mint_failure() {
    let (mut scenario, mut state, owner_cap, coin_metadata) = setup_managed_token_test();

    // Configure a minter
    managed_token::configure_new_minter(
        &mut state,
        &owner_cap,
        MINTER,
        1000,
        false,
        scenario.ctx()
    );

    scenario.next_tx(MINTER);
    let mint_cap = scenario.take_from_sender<MintCap<MANAGED_TOKEN_TEST>>();

    let deny_list = deny_list::new_for_testing(scenario.ctx());
    
    // Try to mint zero amount - should fail
    let coin = managed_token::mint(
        &mut state,
        &mint_cap,
        &deny_list,
        0, // Zero amount should cause failure
        RECIPIENT,
        scenario.ctx()
    );
    transfer::public_transfer(coin, RECIPIENT);

    // This should not be reached due to expected failure
    transfer::public_transfer(mint_cap, MINTER);
    test_utils::destroy(deny_list);
    cleanup_managed_token_test(scenario, state, owner_cap, coin_metadata);
}

#[test]
fun test_increment_mint_allowance() {
    let (mut scenario, mut state, owner_cap, coin_metadata) = setup_managed_token_test();

    // Configure a new minter with limited allowance
    let initial_allowance = 1000;
    managed_token::configure_new_minter(
        &mut state,
        &owner_cap,
        MINTER,
        initial_allowance,
        false, // not unlimited
        scenario.ctx()
    );

    scenario.next_tx(MINTER);
    let mint_cap = scenario.take_from_sender<MintCap<MANAGED_TOKEN_TEST>>();
    let mint_cap_id = object::id(&mint_cap);

    // Create deny list for testing
    let deny_list = deny_list::new_for_testing(scenario.ctx());

    // Test increment_mint_allowance function
    scenario.next_tx(OWNER);
    let increment_amount = 500;
    managed_token::increment_mint_allowance(
        &mut state,
        &owner_cap,
        mint_cap_id,
        &deny_list,
        increment_amount,
        scenario.ctx()
    );

    // Verify allowance was incremented
    let (new_allowance, is_unlimited) = managed_token::mint_allowance(&state, mint_cap_id);
    assert!(new_allowance == initial_allowance + increment_amount);
    assert!(!is_unlimited);

    // Clean up
    transfer::public_transfer(mint_cap, MINTER);
    test_utils::destroy(deny_list);
    cleanup_managed_token_test(scenario, state, owner_cap, coin_metadata);
}

#[test]
fun test_initialize_with_deny_cap() {
    let mut scenario = test_scenario::begin(OWNER);
    let ctx = scenario.ctx();

    // Create the test token with deny cap
    let (treasury_cap, coin_metadata) = coin::create_currency(
        MANAGED_TOKEN_TEST {},
        8,
        b"TEST",
        b"Test Token",
        b"A test token for managed token tests",
        option::none(),
        ctx
    );

    // Initialize managed token with deny cap
    managed_token::initialize(treasury_cap, ctx);

    // Get the shared objects from scenario
    scenario.next_tx(OWNER);
    let state = scenario.take_shared<TokenState<MANAGED_TOKEN_TEST>>();
    let owner_cap = scenario.take_from_sender<OwnerCap<MANAGED_TOKEN_TEST>>();

    // Verify initialization worked
    assert!(managed_token::owner(&state) == OWNER);
    assert!(managed_token::total_supply(&state) == 0);

    // Clean up
    transfer::public_transfer(owner_cap, OWNER);
    transfer::public_freeze_object(coin_metadata);
    test_scenario::return_shared(state);
    scenario.end();
}

#[test]
fun test_set_unlimited_mint_allowances_detailed() {
    let (mut scenario, mut state, owner_cap, coin_metadata) = setup_managed_token_test();

    // Configure a new minter with limited allowance
    managed_token::configure_new_minter(
        &mut state,
        &owner_cap,
        MINTER,
        500,
        false, // not unlimited initially
        scenario.ctx()
    );

    scenario.next_tx(MINTER);
    let mint_cap = scenario.take_from_sender<MintCap<MANAGED_TOKEN_TEST>>();
    let mint_cap_id = object::id(&mint_cap);

    // Create deny list for testing
    let deny_list = deny_list::new_for_testing(scenario.ctx());

    // Verify initial limited allowance
    let (initial_allowance, is_unlimited) = managed_token::mint_allowance(&state, mint_cap_id);
    assert!(initial_allowance == 500);
    assert!(!is_unlimited);

    // Set to unlimited
    scenario.next_tx(OWNER);
    managed_token::set_unlimited_mint_allowances(
        &mut state,
        &owner_cap,
        mint_cap_id,
        &deny_list,
        true,
        scenario.ctx()
    );

    // Verify it's now unlimited
    let (_allowance, is_unlimited) = managed_token::mint_allowance(&state, mint_cap_id);
    assert!(is_unlimited);

    // Test minting a large amount (more than original allowance)
    scenario.next_tx(MINTER);
    let large_amount = 2000;
    let coin = managed_token::mint(
        &mut state,
        &mint_cap,
        &deny_list,
        large_amount,
        RECIPIENT,
        scenario.ctx()
    );

    assert!(coin.value() == large_amount);

    // Set back to limited
    scenario.next_tx(OWNER);
    managed_token::set_unlimited_mint_allowances(
        &mut state,
        &owner_cap,
        mint_cap_id,
        &deny_list,
        false,
        scenario.ctx()
    );

    // Verify it's now limited again
    let (_allowance, is_unlimited) = managed_token::mint_allowance(&state, mint_cap_id);
    assert!(!is_unlimited);

    // Clean up
    transfer::public_transfer(coin, RECIPIENT);
    transfer::public_transfer(mint_cap, MINTER);
    test_utils::destroy(deny_list);
    cleanup_managed_token_test(scenario, state, owner_cap, coin_metadata);
}

#[test]
fun test_pause_and_unpause() {
    let mut scenario = test_scenario::begin(OWNER);
    let ctx = scenario.ctx();

    // Create the test token with deny cap (required for pause/unpause)
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
    let mut state = scenario.take_shared<TokenState<MANAGED_TOKEN_TEST>>();
    let owner_cap = scenario.take_from_sender<OwnerCap<MANAGED_TOKEN_TEST>>();
    // Create deny list for testing
    let mut deny_list = deny_list::new_for_testing(scenario.ctx());

    // Test pause function
    managed_token::pause(
        &mut state,
        &owner_cap,
        &mut deny_list,
        scenario.ctx()
    );

    // Test unpause function
    managed_token::unpause(
        &mut state,
        &owner_cap,
        &mut deny_list,
        scenario.ctx()
    );

    // Clean up
    transfer::public_transfer(owner_cap, OWNER);
    transfer::public_freeze_object(coin_metadata);
    test_utils::destroy(deny_list);
    test_scenario::return_shared(state);
    scenario.end();
}

#[test]
fun test_blocklist_and_unblocklist() {
    let mut scenario = test_scenario::begin(OWNER);
    let ctx = scenario.ctx();

    let (treasury_cap, deny_cap, coin_metadata) = coin::create_regulated_currency_v2(
        MANAGED_TOKEN_TEST {},
        8,
        b"TEST_REG",
        b"Test Regulated Token",
        b"A regulated test token",
        option::none(),
        false,
        ctx
    );

    // Initialize with deny cap
    managed_token::initialize_with_deny_cap(treasury_cap, deny_cap, ctx);

    scenario.next_tx(OWNER);
    let mut state = scenario.take_shared<TokenState<MANAGED_TOKEN_TEST>>();
    let owner_cap = scenario.take_from_sender<OwnerCap<MANAGED_TOKEN_TEST>>();
    // Create deny list for testing
    let mut deny_list = deny_list::new_for_testing(scenario.ctx());

    // Test blocklist function
    let address_to_block = @0x123;
    managed_token::blocklist(
        &mut state,
        &owner_cap,
        &mut deny_list,
        address_to_block,
        scenario.ctx()
    );

    // Test unblocklist function
    managed_token::unblocklist(
        &mut state,
        &owner_cap,
        &mut deny_list,
        address_to_block,
        scenario.ctx()
    );

    // Clean up
    transfer::public_transfer(owner_cap, OWNER);
    transfer::public_freeze_object(coin_metadata);
    test_utils::destroy(deny_list);
    test_scenario::return_shared(state);
    scenario.end();
}

#[test]
#[expected_failure(abort_code = managed_token::EZeroAmount)]
fun test_increment_mint_allowance_zero_amount() {
    let (mut scenario, mut state, owner_cap, coin_metadata) = setup_managed_token_test();

    // Configure a new minter with limited allowance
    managed_token::configure_new_minter(
        &mut state,
        &owner_cap,
        MINTER,
        1000,
        false, // not unlimited
        scenario.ctx()
    );

    scenario.next_tx(MINTER);
    let mint_cap = scenario.take_from_sender<MintCap<MANAGED_TOKEN_TEST>>();
    let mint_cap_id = object::id(&mint_cap);

    // Create deny list for testing
    let deny_list = deny_list::new_for_testing(scenario.ctx());

    // Try to increment with zero amount - should fail
    scenario.next_tx(OWNER);
    managed_token::increment_mint_allowance(
        &mut state,
        &owner_cap,
        mint_cap_id,
        &deny_list,
        0, // Zero increment should cause failure
        scenario.ctx()
    );

    // Clean up (should not be reached)
    transfer::public_transfer(mint_cap, MINTER);
    test_utils::destroy(deny_list);
    cleanup_managed_token_test(scenario, state, owner_cap, coin_metadata);
}

#[test]
#[expected_failure(abort_code = managed_token::EUnauthorizedMintCap)]
fun test_increment_mint_allowance_unauthorized_mint_cap() {
    let (mut scenario, mut state, owner_cap, coin_metadata) = setup_managed_token_test();

    // Create deny list for testing
    let deny_list = deny_list::new_for_testing(scenario.ctx());

    // Use a fake mint cap ID that doesn't exist
    let fake_mint_cap_id = object::id_from_address(@0x999);

    // Try to increment allowance for unauthorized mint cap - should fail
    scenario.next_tx(OWNER);
    managed_token::increment_mint_allowance(
        &mut state,
        &owner_cap,
        fake_mint_cap_id,
        &deny_list,
        500,
        scenario.ctx()
    );

    // Clean up (should not be reached)
    test_utils::destroy(deny_list);
    cleanup_managed_token_test(scenario, state, owner_cap, coin_metadata);
}

#[test]
#[expected_failure(abort_code = managed_token::ECannotIncreaseUnlimitedAllowance)]
fun test_increment_mint_allowance_unlimited_allowance() {
    let (mut scenario, mut state, owner_cap, coin_metadata) = setup_managed_token_test();

    // Configure a new minter with unlimited allowance
    managed_token::configure_new_minter(
        &mut state,
        &owner_cap,
        MINTER,
        0, // amount doesn't matter for unlimited
        true, // unlimited allowance
        scenario.ctx()
    );

    scenario.next_tx(MINTER);
    let mint_cap = scenario.take_from_sender<MintCap<MANAGED_TOKEN_TEST>>();
    let mint_cap_id = object::id(&mint_cap);

    // Create deny list for testing
    let deny_list = deny_list::new_for_testing(scenario.ctx());

    // Try to increment unlimited allowance - should fail
    scenario.next_tx(OWNER);
    managed_token::increment_mint_allowance(
        &mut state,
        &owner_cap,
        mint_cap_id,
        &deny_list,
        500,
        scenario.ctx()
    );

    // Clean up (should not be reached)
    transfer::public_transfer(mint_cap, MINTER);
    test_utils::destroy(deny_list);
    cleanup_managed_token_test(scenario, state, owner_cap, coin_metadata);
}

#[test]
#[expected_failure(abort_code = managed_token::EPaused)]
fun test_increment_mint_allowance_paused() {
    let mut scenario = test_scenario::begin(OWNER);
    let ctx = scenario.ctx();

    // Create the test token with deny cap (required for pause functionality)
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
    let mut state = scenario.take_shared<TokenState<MANAGED_TOKEN_TEST>>();
    let owner_cap = scenario.take_from_sender<OwnerCap<MANAGED_TOKEN_TEST>>();
    let mut deny_list = deny_list::new_for_testing(scenario.ctx());

    // Configure a new minter with limited allowance
    managed_token::configure_new_minter(
        &mut state,
        &owner_cap,
        MINTER,
        1000,
        false, // not unlimited
        scenario.ctx()
    );

    scenario.next_tx(MINTER);
    let mint_cap = scenario.take_from_sender<MintCap<MANAGED_TOKEN_TEST>>();
    let mint_cap_id = object::id(&mint_cap);

    // Pause the token
    scenario.next_tx(OWNER);
    managed_token::pause(
        &mut state,
        &owner_cap,
        &mut deny_list,
        scenario.ctx()
    );

    // Try to increment allowance while paused - should fail
    managed_token::increment_mint_allowance(
        &mut state,
        &owner_cap,
        mint_cap_id,
        &deny_list,
        500,
        scenario.ctx()
    );

    // Clean up (should not be reached)
    transfer::public_transfer(mint_cap, MINTER);
    transfer::public_transfer(owner_cap, OWNER);
    transfer::public_freeze_object(coin_metadata);
    test_utils::destroy(deny_list);
    test_scenario::return_shared(state);
    scenario.end();
}

#[test]
fun test_increment_mint_allowance_success_multiple_increments() {
    let (mut scenario, mut state, owner_cap, coin_metadata) = setup_managed_token_test();

    // Configure a new minter with limited allowance
    let initial_allowance = 1000;
    managed_token::configure_new_minter(
        &mut state,
        &owner_cap,
        MINTER,
        initial_allowance,
        false, // not unlimited
        scenario.ctx()
    );

    scenario.next_tx(MINTER);
    let mint_cap = scenario.take_from_sender<MintCap<MANAGED_TOKEN_TEST>>();
    let mint_cap_id = object::id(&mint_cap);

    // Create deny list for testing
    let deny_list = deny_list::new_for_testing(scenario.ctx());

    // First increment
    scenario.next_tx(OWNER);
    let first_increment = 500;
    managed_token::increment_mint_allowance(
        &mut state,
        &owner_cap,
        mint_cap_id,
        &deny_list,
        first_increment,
        scenario.ctx()
    );

    // Verify first increment
    let (allowance_after_first, is_unlimited) = managed_token::mint_allowance(&state, mint_cap_id);
    assert!(allowance_after_first == initial_allowance + first_increment);
    assert!(!is_unlimited);

    // Second increment
    let second_increment = 300;
    managed_token::increment_mint_allowance(
        &mut state,
        &owner_cap,
        mint_cap_id,
        &deny_list,
        second_increment,
        scenario.ctx()
    );

    // Verify final allowance
    let (final_allowance, is_unlimited) = managed_token::mint_allowance(&state, mint_cap_id);
    assert!(final_allowance == initial_allowance + first_increment + second_increment);
    assert!(!is_unlimited);

    // Test that minting still works with the increased allowance
    scenario.next_tx(MINTER);
    let mint_amount = 1500; // More than initial allowance but less than incremented allowance
    let coin = managed_token::mint(
        &mut state,
        &mint_cap,
        &deny_list,
        mint_amount,
        RECIPIENT,
        scenario.ctx()
    );

    assert!(coin.value() == mint_amount);

    // Verify remaining allowance
    let (remaining_allowance, _) = managed_token::mint_allowance(&state, mint_cap_id);
    assert!(remaining_allowance == final_allowance - mint_amount);

    // Clean up
    transfer::public_transfer(coin, RECIPIENT);
    transfer::public_transfer(mint_cap, MINTER);
    test_utils::destroy(deny_list);
    cleanup_managed_token_test(scenario, state, owner_cap, coin_metadata);
}

#[test]
#[expected_failure(abort_code = managed_token::EInvalidOwnerCap)]
fun test_set_unlimited_mint_allowances_invalid_owner_cap() {
    let (mut scenario, mut state, owner_cap, coin_metadata) = setup_managed_token_test();

    // Configure a new minter
    managed_token::configure_new_minter(
        &mut state,
        &owner_cap,
        MINTER,
        1000,
        false,
        scenario.ctx()
    );

    scenario.next_tx(MINTER);
    let mint_cap = scenario.take_from_sender<MintCap<MANAGED_TOKEN_TEST>>();
    let mint_cap_id = object::id(&mint_cap);

    // Create deny list for testing
    let deny_list = deny_list::new_for_testing(scenario.ctx());

    // Create a fake owner cap by creating another managed token
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

    // Try to use fake owner cap on original state - should fail
    scenario.next_tx(OTHER_USER);
    managed_token::set_unlimited_mint_allowances(
        &mut state,
        &fake_owner_cap, // Wrong owner cap
        mint_cap_id,
        &deny_list,
        true,
        scenario.ctx()
    );

    // Clean up (should not be reached)
    transfer::public_transfer(mint_cap, MINTER);
    transfer::public_transfer(fake_owner_cap, OTHER_USER);
    transfer::public_freeze_object(fake_coin_metadata);
    test_scenario::return_shared(fake_state);
    test_utils::destroy(deny_list);
    cleanup_managed_token_test(scenario, state, owner_cap, coin_metadata);
}

#[test]
#[expected_failure(abort_code = managed_token::EPaused)]
fun test_set_unlimited_mint_allowances_paused() {
    let mut scenario = test_scenario::begin(OWNER);
    let ctx = scenario.ctx();

    // Create the test token with deny cap (required for pause functionality)
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
    let mut state = scenario.take_shared<TokenState<MANAGED_TOKEN_TEST>>();
    let owner_cap = scenario.take_from_sender<OwnerCap<MANAGED_TOKEN_TEST>>();
    let mut deny_list = deny_list::new_for_testing(scenario.ctx());

    // Configure a new minter
    managed_token::configure_new_minter(
        &mut state,
        &owner_cap,
        MINTER,
        1000,
        false,
        scenario.ctx()
    );

    scenario.next_tx(MINTER);
    let mint_cap = scenario.take_from_sender<MintCap<MANAGED_TOKEN_TEST>>();
    let mint_cap_id = object::id(&mint_cap);

    // Pause the token
    scenario.next_tx(OWNER);
    managed_token::pause(
        &mut state,
        &owner_cap,
        &mut deny_list,
        scenario.ctx()
    );

    // Try to set unlimited allowance while paused - should fail
    managed_token::set_unlimited_mint_allowances(
        &mut state,
        &owner_cap,
        mint_cap_id,
        &deny_list,
        true,
        scenario.ctx()
    );

    // Clean up (should not be reached)
    transfer::public_transfer(mint_cap, MINTER);
    transfer::public_transfer(owner_cap, OWNER);
    transfer::public_freeze_object(coin_metadata);
    test_utils::destroy(deny_list);
    test_scenario::return_shared(state);
    scenario.end();
}

#[test]
#[expected_failure(abort_code = managed_token::EUnauthorizedMintCap)]
fun test_set_unlimited_mint_allowances_unauthorized_mint_cap() {
    let (mut scenario, mut state, owner_cap, coin_metadata) = setup_managed_token_test();

    // Create deny list for testing
    let deny_list = deny_list::new_for_testing(scenario.ctx());

    // Use a fake mint cap ID that doesn't exist
    let fake_mint_cap_id = object::id_from_address(@0x999);

    // Try to set unlimited allowance for unauthorized mint cap - should fail
    scenario.next_tx(OWNER);
    managed_token::set_unlimited_mint_allowances(
        &mut state,
        &owner_cap,
        fake_mint_cap_id,
        &deny_list,
        true,
        scenario.ctx()
    );

    // Clean up (should not be reached)
    test_utils::destroy(deny_list);
    cleanup_managed_token_test(scenario, state, owner_cap, coin_metadata);
}

#[test]
fun test_set_unlimited_mint_allowances_success_scenarios() {
    let (mut scenario, mut state, owner_cap, coin_metadata) = setup_managed_token_test();

    // Configure a new minter with limited allowance
    let initial_allowance = 1000;
    managed_token::configure_new_minter(
        &mut state,
        &owner_cap,
        MINTER,
        initial_allowance,
        false, // not unlimited initially
        scenario.ctx()
    );

    scenario.next_tx(MINTER);
    let mint_cap = scenario.take_from_sender<MintCap<MANAGED_TOKEN_TEST>>();
    let mint_cap_id = object::id(&mint_cap);

    // Create deny list for testing
    let deny_list = deny_list::new_for_testing(scenario.ctx());

    // Verify initial limited allowance
    let (allowance, is_unlimited) = managed_token::mint_allowance(&state, mint_cap_id);
    assert!(allowance == initial_allowance);
    assert!(!is_unlimited);

    // Set to unlimited
    scenario.next_tx(OWNER);
    managed_token::set_unlimited_mint_allowances(
        &mut state,
        &owner_cap,
        mint_cap_id,
        &deny_list,
        true,
        scenario.ctx()
    );

    // Verify it's now unlimited
    let (_allowance, is_unlimited) = managed_token::mint_allowance(&state, mint_cap_id);
    assert!(is_unlimited);

    // Test minting a large amount (more than original allowance)
    scenario.next_tx(MINTER);
    let large_amount = 5000;
    let coin1 = managed_token::mint(
        &mut state,
        &mint_cap,
        &deny_list,
        large_amount,
        RECIPIENT,
        scenario.ctx()
    );
    assert!(coin1.value() == large_amount);

    // Test another large mint to verify it's truly unlimited
    let another_large_amount = 10000;
    let coin2 = managed_token::mint(
        &mut state,
        &mint_cap,
        &deny_list,
        another_large_amount,
        RECIPIENT,
        scenario.ctx()
    );
    assert!(coin2.value() == another_large_amount);

    // Set back to limited (with allowance = 0)
    managed_token::set_unlimited_mint_allowances(
        &mut state,
        &owner_cap,
        mint_cap_id,
        &deny_list,
        false,
        scenario.ctx()
    );

    // Verify it's now limited with 0 allowance
    let (limited_allowance, is_unlimited) = managed_token::mint_allowance(&state, mint_cap_id);
    assert!(limited_allowance == 0);
    assert!(!is_unlimited);

    // Clean up
    transfer::public_transfer(coin1, RECIPIENT);
    transfer::public_transfer(coin2, RECIPIENT);
    transfer::public_transfer(mint_cap, MINTER);
    test_utils::destroy(deny_list);
    cleanup_managed_token_test(scenario, state, owner_cap, coin_metadata);
}

#[test]
fun test_set_unlimited_mint_allowances_toggle_behavior() {
    let (mut scenario, mut state, owner_cap, coin_metadata) = setup_managed_token_test();

    // Configure a new minter with unlimited allowance initially
    managed_token::configure_new_minter(
        &mut state,
        &owner_cap,
        MINTER,
        0, // amount doesn't matter for unlimited
        true, // unlimited initially
        scenario.ctx()
    );

    scenario.next_tx(MINTER);
    let mint_cap = scenario.take_from_sender<MintCap<MANAGED_TOKEN_TEST>>();
    let mint_cap_id = object::id(&mint_cap);

    // Create deny list for testing
    let deny_list = deny_list::new_for_testing(scenario.ctx());

    // Verify initial unlimited allowance
    let (_allowance, is_unlimited) = managed_token::mint_allowance(&state, mint_cap_id);
    assert!(is_unlimited);

    // Set to limited (should set allowance to 0)
    scenario.next_tx(OWNER);
    managed_token::set_unlimited_mint_allowances(
        &mut state,
        &owner_cap,
        mint_cap_id,
        &deny_list,
        false,
        scenario.ctx()
    );

    // Verify it's now limited with 0 allowance
    let (limited_allowance, is_unlimited) = managed_token::mint_allowance(&state, mint_cap_id);
    assert!(limited_allowance == 0);
    assert!(!is_unlimited);

    // Set back to unlimited
    managed_token::set_unlimited_mint_allowances(
        &mut state,
        &owner_cap,
        mint_cap_id,
        &deny_list,
        true,
        scenario.ctx()
    );

    // Verify it's unlimited again
    let (_allowance, is_unlimited) = managed_token::mint_allowance(&state, mint_cap_id);
    assert!(is_unlimited);

    // Test that unlimited minting works
    scenario.next_tx(MINTER);
    let test_amount = 1000000; // Very large amount
    let coin = managed_token::mint(
        &mut state,
        &mint_cap,
        &deny_list,
        test_amount,
        RECIPIENT,
        scenario.ctx()
    );
    assert!(coin.value() == test_amount);

    // Clean up
    transfer::public_transfer(coin, RECIPIENT);
    transfer::public_transfer(mint_cap, MINTER);
    test_utils::destroy(deny_list);
    cleanup_managed_token_test(scenario, state, owner_cap, coin_metadata);
}

#[test]
#[expected_failure(abort_code = managed_token::EPaused)]
fun test_burn_paused() {
    let mut scenario = test_scenario::begin(OWNER);
    let ctx = scenario.ctx();

    // Create the test token with deny cap (required for pause functionality)
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
    let mut state = scenario.take_shared<TokenState<MANAGED_TOKEN_TEST>>();
    let owner_cap = scenario.take_from_sender<OwnerCap<MANAGED_TOKEN_TEST>>();
    let mut deny_list = deny_list::new_for_testing(scenario.ctx());

    // Configure a new minter
    managed_token::configure_new_minter(
        &mut state,
        &owner_cap,
        MINTER,
        1000,
        false,
        scenario.ctx()
    );

    scenario.next_tx(MINTER);
    let mint_cap = scenario.take_from_sender<MintCap<MANAGED_TOKEN_TEST>>();

    // Mint some tokens first
    let coin = managed_token::mint(
        &mut state,
        &mint_cap,
        &deny_list,
        500,
        MINTER,
        scenario.ctx()
    );

    // Pause the token
    scenario.next_tx(OWNER);
    managed_token::pause(
        &mut state,
        &owner_cap,
        &mut deny_list,
        scenario.ctx()
    );

    // Try to burn while paused - should fail
    scenario.next_tx(MINTER);
    managed_token::burn(
        &mut state,
        &mint_cap,
        &deny_list,
        coin,
        MINTER,
        scenario.ctx()
    );

    // Clean up (should not be reached)
    transfer::public_transfer(mint_cap, MINTER);
    transfer::public_transfer(owner_cap, OWNER);
    transfer::public_freeze_object(coin_metadata);
    test_utils::destroy(deny_list);
    test_scenario::return_shared(state);
    scenario.end();
}

#[test]
#[expected_failure(abort_code = managed_token::EDeniedAddress)]
fun test_burn_blocklisted_address() {
    let mut scenario = test_scenario::begin(OWNER);
    let ctx = scenario.ctx();

    // Create the test token with deny cap (required for blocklist functionality)
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
    let mut state = scenario.take_shared<TokenState<MANAGED_TOKEN_TEST>>();
    let owner_cap = scenario.take_from_sender<OwnerCap<MANAGED_TOKEN_TEST>>();
    let mut deny_list = deny_list::new_for_testing(scenario.ctx());

    // Configure a new minter
    managed_token::configure_new_minter(
        &mut state,
        &owner_cap,
        MINTER,
        1000,
        false,
        scenario.ctx()
    );

    scenario.next_tx(MINTER);
    let mint_cap = scenario.take_from_sender<MintCap<MANAGED_TOKEN_TEST>>();

    // Mint some tokens first
    let coin = managed_token::mint(
        &mut state,
        &mint_cap,
        &deny_list,
        500,
        MINTER,
        scenario.ctx()
    );

    // Blocklist the minter
    scenario.next_tx(OWNER);
    managed_token::blocklist(
        &mut state,
        &owner_cap,
        &mut deny_list,
        MINTER,
        scenario.ctx()
    );

    // Try to burn while blocklisted - should fail
    scenario.next_tx(MINTER);
    managed_token::burn(
        &mut state,
        &mint_cap,
        &deny_list,
        coin,
        MINTER,
        scenario.ctx()
    );

    // Clean up (should not be reached)
    transfer::public_transfer(mint_cap, MINTER);
    transfer::public_transfer(owner_cap, OWNER);
    transfer::public_freeze_object(coin_metadata);
    test_utils::destroy(deny_list);
    test_scenario::return_shared(state);
    scenario.end();
}

#[test]
#[expected_failure(abort_code = managed_token::EUnauthorizedMintCap)]
fun test_burn_unauthorized_mint_cap() {
    let (mut scenario, mut state, owner_cap, coin_metadata) = setup_managed_token_test();

    // Configure a new minter
    managed_token::configure_new_minter(
        &mut state,
        &owner_cap,
        MINTER,
        1000,
        false,
        scenario.ctx()
    );

    scenario.next_tx(MINTER);
    let mint_cap = scenario.take_from_sender<MintCap<MANAGED_TOKEN_TEST>>();

    // Create deny list for testing
    let deny_list = deny_list::new_for_testing(scenario.ctx());

    // Mint some tokens first
    let coin = managed_token::mint(
        &mut state,
        &mint_cap,
        &deny_list,
        500,
        MINTER,
        scenario.ctx()
    );

    // Create another managed token to get a different mint cap
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
    let mut fake_state = scenario.take_shared<TokenState<MANAGED_TOKEN_TEST>>();
    let fake_owner_cap = scenario.take_from_sender<OwnerCap<MANAGED_TOKEN_TEST>>();

    // Configure a mint cap on the fake state
    managed_token::configure_new_minter(
        &mut fake_state,
        &fake_owner_cap,
        OTHER_USER,
        1000,
        false,
        scenario.ctx()
    );

    scenario.next_tx(OTHER_USER);
    let fake_mint_cap = scenario.take_from_sender<MintCap<MANAGED_TOKEN_TEST>>();

    // Try to burn with mint cap from different state - should fail
    scenario.next_tx(MINTER);
    managed_token::burn(
        &mut state,
        &fake_mint_cap, // Unauthorized mint cap
        &deny_list,
        coin,
        MINTER,
        scenario.ctx()
    );

    // Clean up (should not be reached)
    transfer::public_transfer(mint_cap, MINTER);
    transfer::public_transfer(fake_mint_cap, OTHER_USER);
    transfer::public_transfer(fake_owner_cap, OTHER_USER);
    transfer::public_freeze_object(fake_coin_metadata);
    test_scenario::return_shared(fake_state);
    test_utils::destroy(deny_list);
    cleanup_managed_token_test(scenario, state, owner_cap, coin_metadata);
}

#[test]
#[expected_failure(abort_code = managed_token::EZeroAmount)]
fun test_burn_zero_amount() {
    let (mut scenario, mut state, owner_cap, coin_metadata) = setup_managed_token_test();

    // Configure a new minter
    managed_token::configure_new_minter(
        &mut state,
        &owner_cap,
        MINTER,
        1000,
        false,
        scenario.ctx()
    );

    scenario.next_tx(MINTER);
    let mint_cap = scenario.take_from_sender<MintCap<MANAGED_TOKEN_TEST>>();

    // Create deny list for testing
    let deny_list = deny_list::new_for_testing(scenario.ctx());

    // Create a zero-value coin
    let zero_coin = coin::zero<MANAGED_TOKEN_TEST>(scenario.ctx());

    // Try to burn zero amount coin - should fail
    scenario.next_tx(MINTER);
    managed_token::burn(
        &mut state,
        &mint_cap,
        &deny_list,
        zero_coin,
        MINTER,
        scenario.ctx()
    );

    // Clean up (should not be reached)
    transfer::public_transfer(mint_cap, MINTER);
    test_utils::destroy(deny_list);
    cleanup_managed_token_test(scenario, state, owner_cap, coin_metadata);
}

#[test]
fun test_burn_success() {
    let (mut scenario, mut state, owner_cap, coin_metadata) = setup_managed_token_test();

    // Configure a new minter
    managed_token::configure_new_minter(
        &mut state,
        &owner_cap,
        MINTER,
        1000,
        false,
        scenario.ctx()
    );

    scenario.next_tx(MINTER);
    let mint_cap = scenario.take_from_sender<MintCap<MANAGED_TOKEN_TEST>>();

    // Create deny list for testing
    let deny_list = deny_list::new_for_testing(scenario.ctx());

    // Mint some tokens first
    let mint_amount = 500;
    let coin = managed_token::mint(
        &mut state,
        &mint_cap,
        &deny_list,
        mint_amount,
        MINTER,
        scenario.ctx()
    );

    // Verify initial total supply
    assert!(managed_token::total_supply(&state) == mint_amount);

    // Burn the tokens
    scenario.next_tx(MINTER);
    managed_token::burn(
        &mut state,
        &mint_cap,
        &deny_list,
        coin,
        MINTER,
        scenario.ctx()
    );

    // Verify total supply decreased
    assert!(managed_token::total_supply(&state) == 0);

    // Clean up
    transfer::public_transfer(mint_cap, MINTER);
    test_utils::destroy(deny_list);
    cleanup_managed_token_test(scenario, state, owner_cap, coin_metadata);
}

#[test]
fun test_burn_partial_amount() {
    let (mut scenario, mut state, owner_cap, coin_metadata) = setup_managed_token_test();

    // Configure a new minter with unlimited allowance
    managed_token::configure_new_minter(
        &mut state,
        &owner_cap,
        MINTER,
        0,
        true, // unlimited
        scenario.ctx()
    );

    scenario.next_tx(MINTER);
    let mint_cap = scenario.take_from_sender<MintCap<MANAGED_TOKEN_TEST>>();

    // Create deny list for testing
    let deny_list = deny_list::new_for_testing(scenario.ctx());

    // Mint some tokens
    let total_mint_amount = 1000;
    let coin1 = managed_token::mint(
        &mut state,
        &mint_cap,
        &deny_list,
        total_mint_amount,
        MINTER,
        scenario.ctx()
    );

    let coin2 = managed_token::mint(
        &mut state,
        &mint_cap,
        &deny_list,
        total_mint_amount,
        MINTER,
        scenario.ctx()
    );

    // Verify initial total supply
    assert!(managed_token::total_supply(&state) == total_mint_amount * 2);

    // Burn only one coin
    scenario.next_tx(MINTER);
    managed_token::burn(
        &mut state,
        &mint_cap,
        &deny_list,
        coin1,
        MINTER,
        scenario.ctx()
    );

    // Verify total supply decreased by the burned amount only
    assert!(managed_token::total_supply(&state) == total_mint_amount);

    // Burn the second coin
    managed_token::burn(
        &mut state,
        &mint_cap,
        &deny_list,
        coin2,
        MINTER,
        scenario.ctx()
    );

    // Verify total supply is now zero
    assert!(managed_token::total_supply(&state) == 0);

    // Clean up
    transfer::public_transfer(mint_cap, MINTER);
    test_utils::destroy(deny_list);
    cleanup_managed_token_test(scenario, state, owner_cap, coin_metadata);
}

#[test]
fun test_burn_multiple_minters() {
    let (mut scenario, mut state, owner_cap, coin_metadata) = setup_managed_token_test();

    // Configure two minters
    managed_token::configure_new_minter(
        &mut state,
        &owner_cap,
        MINTER,
        1000,
        false,
        scenario.ctx()
    );

    managed_token::configure_new_minter(
        &mut state,
        &owner_cap,
        OTHER_USER,
        1000,
        false,
        scenario.ctx()
    );

    // Get mint caps
    scenario.next_tx(MINTER);
    let mint_cap1 = scenario.take_from_sender<MintCap<MANAGED_TOKEN_TEST>>();

    scenario.next_tx(OTHER_USER);
    let mint_cap2 = scenario.take_from_sender<MintCap<MANAGED_TOKEN_TEST>>();

    // Create deny list for testing
    let deny_list = deny_list::new_for_testing(scenario.ctx());

    // Both minters mint tokens
    scenario.next_tx(MINTER);
    let coin1 = managed_token::mint(
        &mut state,
        &mint_cap1,
        &deny_list,
        300,
        MINTER,
        scenario.ctx()
    );

    scenario.next_tx(OTHER_USER);
    let coin2 = managed_token::mint(
        &mut state,
        &mint_cap2,
        &deny_list,
        400,
        OTHER_USER,
        scenario.ctx()
    );

    // Verify total supply
    assert!(managed_token::total_supply(&state) == 700);

    // First minter burns their tokens
    scenario.next_tx(MINTER);
    managed_token::burn(
        &mut state,
        &mint_cap1,
        &deny_list,
        coin1,
        MINTER,
        scenario.ctx()
    );

    // Verify partial burn
    assert!(managed_token::total_supply(&state) == 400);

    // Second minter burns their tokens
    scenario.next_tx(OTHER_USER);
    managed_token::burn(
        &mut state,
        &mint_cap2,
        &deny_list,
        coin2,
        OTHER_USER,
        scenario.ctx()
    );

    // Verify all tokens burned
    assert!(managed_token::total_supply(&state) == 0);

    // Clean up
    transfer::public_transfer(mint_cap1, MINTER);
    transfer::public_transfer(mint_cap2, OTHER_USER);
    test_utils::destroy(deny_list);
    cleanup_managed_token_test(scenario, state, owner_cap, coin_metadata);
}

#[test]
#[expected_failure(abort_code = managed_token::EPaused)]
fun test_validate_mint_paused() {
    let mut scenario = test_scenario::begin(OWNER);
    let ctx = scenario.ctx();

    // Create the test token with deny cap (required for pause functionality)
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
    let mut state = scenario.take_shared<TokenState<MANAGED_TOKEN_TEST>>();
    let owner_cap = scenario.take_from_sender<OwnerCap<MANAGED_TOKEN_TEST>>();
    let mut deny_list = deny_list::new_for_testing(scenario.ctx());

    // Configure a new minter
    managed_token::configure_new_minter(
        &mut state,
        &owner_cap,
        MINTER,
        1000,
        false,
        scenario.ctx()
    );

    scenario.next_tx(MINTER);
    let mint_cap = scenario.take_from_sender<MintCap<MANAGED_TOKEN_TEST>>();

    // Pause the token
    scenario.next_tx(OWNER);
    managed_token::pause(
        &mut state,
        &owner_cap,
        &mut deny_list,
        scenario.ctx()
    );

    // Try to mint while paused - should fail
    scenario.next_tx(MINTER);
    managed_token::mint_and_transfer(
        &mut state,
        &mint_cap,
        &deny_list,
        500,
        RECIPIENT,
        scenario.ctx()
    );

    // Clean up (should not be reached)
    transfer::public_transfer(mint_cap, MINTER);
    transfer::public_transfer(owner_cap, OWNER);
    transfer::public_freeze_object(coin_metadata);
    test_utils::destroy(deny_list);
    test_scenario::return_shared(state);
    scenario.end();
}

#[test]
#[expected_failure(abort_code = managed_token::EDeniedAddress)]
fun test_validate_mint_blocklisted_minter() {
    let mut scenario = test_scenario::begin(OWNER);
    let ctx = scenario.ctx();

    // Create the test token with deny cap (required for blocklist functionality)
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
    let mut state = scenario.take_shared<TokenState<MANAGED_TOKEN_TEST>>();
    let owner_cap = scenario.take_from_sender<OwnerCap<MANAGED_TOKEN_TEST>>();
    let mut deny_list = deny_list::new_for_testing(scenario.ctx());

    // Configure a new minter
    managed_token::configure_new_minter(
        &mut state,
        &owner_cap,
        MINTER,
        1000,
        false,
        scenario.ctx()
    );

    scenario.next_tx(MINTER);
    let mint_cap = scenario.take_from_sender<MintCap<MANAGED_TOKEN_TEST>>();

    // Blocklist the minter
    scenario.next_tx(OWNER);
    managed_token::blocklist(
        &mut state,
        &owner_cap,
        &mut deny_list,
        MINTER,
        scenario.ctx()
    );

    // Try to mint while minter is blocklisted - should fail
    scenario.next_tx(MINTER);
    managed_token::mint_and_transfer(
        &mut state,
        &mint_cap,
        &deny_list,
        500,
        RECIPIENT,
        scenario.ctx()
    );

    // Clean up (should not be reached)
    transfer::public_transfer(mint_cap, MINTER);
    transfer::public_transfer(owner_cap, OWNER);
    transfer::public_freeze_object(coin_metadata);
    test_utils::destroy(deny_list);
    test_scenario::return_shared(state);
    scenario.end();
}

#[test]
#[expected_failure(abort_code = managed_token::EDeniedAddress)]
fun test_validate_mint_blocklisted_recipient() {
    let mut scenario = test_scenario::begin(OWNER);
    let ctx = scenario.ctx();

    // Create the test token with deny cap (required for blocklist functionality)
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
    let mut state = scenario.take_shared<TokenState<MANAGED_TOKEN_TEST>>();
    let owner_cap = scenario.take_from_sender<OwnerCap<MANAGED_TOKEN_TEST>>();
    let mut deny_list = deny_list::new_for_testing(scenario.ctx());

    // Configure a new minter
    managed_token::configure_new_minter(
        &mut state,
        &owner_cap,
        MINTER,
        1000,
        false,
        scenario.ctx()
    );

    scenario.next_tx(MINTER);
    let mint_cap = scenario.take_from_sender<MintCap<MANAGED_TOKEN_TEST>>();

    // Blocklist the recipient
    scenario.next_tx(OWNER);
    managed_token::blocklist(
        &mut state,
        &owner_cap,
        &mut deny_list,
        RECIPIENT,
        scenario.ctx()
    );

    // Try to mint to blocklisted recipient - should fail
    scenario.next_tx(MINTER);
    managed_token::mint_and_transfer(
        &mut state,
        &mint_cap,
        &deny_list,
        500,
        RECIPIENT, // Blocklisted recipient
        scenario.ctx()
    );

    // Clean up (should not be reached)
    transfer::public_transfer(mint_cap, MINTER);
    transfer::public_transfer(owner_cap, OWNER);
    transfer::public_freeze_object(coin_metadata);
    test_utils::destroy(deny_list);
    test_scenario::return_shared(state);
    scenario.end();
}

#[test]
#[expected_failure(abort_code = managed_token::EUnauthorizedMintCap)]
fun test_validate_mint_unauthorized_mint_cap() {
    let (mut scenario, mut state, owner_cap, coin_metadata) = setup_managed_token_test();

    // Configure a new minter
    managed_token::configure_new_minter(
        &mut state,
        &owner_cap,
        MINTER,
        1000,
        false,
        scenario.ctx()
    );

    scenario.next_tx(MINTER);
    let mint_cap = scenario.take_from_sender<MintCap<MANAGED_TOKEN_TEST>>();

    // Create deny list for testing
    let deny_list = deny_list::new_for_testing(scenario.ctx());

    // Create another managed token to get a different mint cap
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
    let mut fake_state = scenario.take_shared<TokenState<MANAGED_TOKEN_TEST>>();
    let fake_owner_cap = scenario.take_from_sender<OwnerCap<MANAGED_TOKEN_TEST>>();

    // Configure a mint cap on the fake state
    managed_token::configure_new_minter(
        &mut fake_state,
        &fake_owner_cap,
        OTHER_USER,
        1000,
        false,
        scenario.ctx()
    );

    scenario.next_tx(OTHER_USER);
    let fake_mint_cap = scenario.take_from_sender<MintCap<MANAGED_TOKEN_TEST>>();

    // Try to mint with mint cap from different state - should fail
    scenario.next_tx(MINTER);
    managed_token::mint_and_transfer(
        &mut state,
        &fake_mint_cap, // Unauthorized mint cap
        &deny_list,
        500,
        RECIPIENT,
        scenario.ctx()
    );

    // Clean up (should not be reached)
    transfer::public_transfer(mint_cap, MINTER);
    transfer::public_transfer(fake_mint_cap, OTHER_USER);
    transfer::public_transfer(fake_owner_cap, OTHER_USER);
    transfer::public_freeze_object(fake_coin_metadata);
    test_scenario::return_shared(fake_state);
    test_utils::destroy(deny_list);
    cleanup_managed_token_test(scenario, state, owner_cap, coin_metadata);
}

#[test]
#[expected_failure(abort_code = managed_token::EZeroAmount)]
fun test_validate_mint_zero_amount() {
    let (mut scenario, mut state, owner_cap, coin_metadata) = setup_managed_token_test();

    // Configure a new minter
    managed_token::configure_new_minter(
        &mut state,
        &owner_cap,
        MINTER,
        1000,
        false,
        scenario.ctx()
    );

    scenario.next_tx(MINTER);
    let mint_cap = scenario.take_from_sender<MintCap<MANAGED_TOKEN_TEST>>();

    // Create deny list for testing
    let deny_list = deny_list::new_for_testing(scenario.ctx());

    // Try to mint zero amount - should fail
    scenario.next_tx(MINTER);
    managed_token::mint_and_transfer(
        &mut state,
        &mint_cap,
        &deny_list,
        0, // Zero amount should cause failure
        RECIPIENT,
        scenario.ctx()
    );

    // Clean up (should not be reached)
    transfer::public_transfer(mint_cap, MINTER);
    test_utils::destroy(deny_list);
    cleanup_managed_token_test(scenario, state, owner_cap, coin_metadata);
}

#[test]
#[expected_failure(abort_code = managed_token::EInsufficientAllowance)]
fun test_validate_mint_insufficient_allowance() {
    let (mut scenario, mut state, owner_cap, coin_metadata) = setup_managed_token_test();

    // Configure a new minter with limited allowance
    let allowance = 500;
    managed_token::configure_new_minter(
        &mut state,
        &owner_cap,
        MINTER,
        allowance,
        false, // not unlimited
        scenario.ctx()
    );

    scenario.next_tx(MINTER);
    let mint_cap = scenario.take_from_sender<MintCap<MANAGED_TOKEN_TEST>>();

    // Create deny list for testing
    let deny_list = deny_list::new_for_testing(scenario.ctx());

    // Try to mint more than allowance - should fail
    scenario.next_tx(MINTER);
    managed_token::mint_and_transfer(
        &mut state,
        &mint_cap,
        &deny_list,
        allowance + 1, // More than allowance
        RECIPIENT,
        scenario.ctx()
    );

    // Clean up (should not be reached)
    transfer::public_transfer(mint_cap, MINTER);
    test_utils::destroy(deny_list);
    cleanup_managed_token_test(scenario, state, owner_cap, coin_metadata);
}

#[test]
fun test_validate_mint_success_scenarios() {
    let (mut scenario, mut state, owner_cap, coin_metadata) = setup_managed_token_test();

    // Configure a new minter with sufficient allowance
    let allowance = 1000;
    managed_token::configure_new_minter(
        &mut state,
        &owner_cap,
        MINTER,
        allowance,
        false,
        scenario.ctx()
    );

    scenario.next_tx(MINTER);
    let mint_cap = scenario.take_from_sender<MintCap<MANAGED_TOKEN_TEST>>();

    // Create deny list for testing
    let deny_list = deny_list::new_for_testing(scenario.ctx());

    // Test successful mint within allowance
    scenario.next_tx(MINTER);
    let mint_amount = 300;
    let coin1 = managed_token::mint(
        &mut state,
        &mint_cap,
        &deny_list,
        mint_amount,
        RECIPIENT,
        scenario.ctx()
    );

    assert!(coin1.value() == mint_amount);
    assert!(managed_token::total_supply(&state) == mint_amount);

    // Test mint_and_transfer function
    let transfer_amount = 200;
    managed_token::mint_and_transfer(
        &mut state,
        &mint_cap,
        &deny_list,
        transfer_amount,
        RECIPIENT,
        scenario.ctx()
    );

    assert!(managed_token::total_supply(&state) == mint_amount + transfer_amount);

    // Verify remaining allowance
    let (remaining_allowance, is_unlimited) = managed_token::mint_allowance(&state, object::id(&mint_cap));
    assert!(remaining_allowance == allowance - mint_amount - transfer_amount);
    assert!(!is_unlimited);

    // Test minting exactly the remaining allowance
    let final_coin = managed_token::mint(
        &mut state,
        &mint_cap,
        &deny_list,
        remaining_allowance,
        RECIPIENT,
        scenario.ctx()
    );

    assert!(final_coin.value() == remaining_allowance);

    // Verify allowance is now exhausted
    let (final_allowance, _) = managed_token::mint_allowance(&state, object::id(&mint_cap));
    assert!(final_allowance == 0);

    // Clean up
    transfer::public_transfer(coin1, RECIPIENT);
    transfer::public_transfer(final_coin, RECIPIENT);
    transfer::public_transfer(mint_cap, MINTER);
    test_utils::destroy(deny_list);
    cleanup_managed_token_test(scenario, state, owner_cap, coin_metadata);
}

#[test]
fun test_validate_mint_unlimited_allowance() {
    let (mut scenario, mut state, owner_cap, coin_metadata) = setup_managed_token_test();

    // Configure a new minter with unlimited allowance
    managed_token::configure_new_minter(
        &mut state,
        &owner_cap,
        MINTER,
        0, // amount doesn't matter for unlimited
        true, // unlimited
        scenario.ctx()
    );

    scenario.next_tx(MINTER);
    let mint_cap = scenario.take_from_sender<MintCap<MANAGED_TOKEN_TEST>>();

    // Create deny list for testing
    let deny_list = deny_list::new_for_testing(scenario.ctx());

    // Test minting very large amounts with unlimited allowance
    scenario.next_tx(MINTER);
    let large_amount1 = 1000000;
    let coin1 = managed_token::mint(
        &mut state,
        &mint_cap,
        &deny_list,
        large_amount1,
        RECIPIENT,
        scenario.ctx()
    );

    assert!(coin1.value() == large_amount1);

    // Test another large mint to verify unlimited behavior
    let large_amount2 = 5000000;
    let coin2 = managed_token::mint(
        &mut state,
        &mint_cap,
        &deny_list,
        large_amount2,
        RECIPIENT,
        scenario.ctx()
    );

    assert!(coin2.value() == large_amount2);

    // Verify allowance is still unlimited
    let (_allowance, is_unlimited) = managed_token::mint_allowance(&state, object::id(&mint_cap));
    assert!(is_unlimited);

    // Verify total supply
    assert!(managed_token::total_supply(&state) == large_amount1 + large_amount2);

    // Clean up
    transfer::public_transfer(coin1, RECIPIENT);
    transfer::public_transfer(coin2, RECIPIENT);
    transfer::public_transfer(mint_cap, MINTER);
    test_utils::destroy(deny_list);
    cleanup_managed_token_test(scenario, state, owner_cap, coin_metadata);
}

#[test]
#[expected_failure(abort_code = managed_token::EInvalidOwnerCap)]
fun test_destroy_managed_token_invalid_owner_cap() {
    let (mut scenario, state, owner_cap, coin_metadata) = setup_managed_token_test();

    // Create another managed token to get a different owner cap
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
    scenario.next_tx(OTHER_USER);
    let (treasury_cap_result, mut deny_cap_result) = managed_token::destroy_managed_token(
        fake_owner_cap, // Wrong owner cap
        state, // Original state
        scenario.ctx()
    );
    
    // This code should not be reached due to expected failure
    transfer::public_transfer(treasury_cap_result, OTHER_USER);
    if (deny_cap_result.is_some()) {
        transfer::public_transfer(deny_cap_result.extract(), OTHER_USER);
    };
    deny_cap_result.destroy_none();

    // Clean up (should not be reached)
    transfer::public_transfer(owner_cap, OWNER);
    transfer::public_freeze_object(coin_metadata);
    transfer::public_freeze_object(fake_coin_metadata);
    test_scenario::return_shared(fake_state);
    scenario.end();
}

#[test]
fun test_destroy_managed_token_without_deny_cap() {
    let (mut scenario, mut state, owner_cap, coin_metadata) = setup_managed_token_test();

    // Configure some minters to test cleanup
    managed_token::configure_new_minter(
        &mut state,
        &owner_cap,
        MINTER,
        1000,
        false,
        scenario.ctx()
    );

    managed_token::configure_new_minter(
        &mut state,
        &owner_cap,
        OTHER_USER,
        2000,
        true, // unlimited
        scenario.ctx()
    );

    // Get the mint caps
    scenario.next_tx(MINTER);
    let mint_cap1 = scenario.take_from_sender<MintCap<MANAGED_TOKEN_TEST>>();

    scenario.next_tx(OTHER_USER);
    let mint_cap2 = scenario.take_from_sender<MintCap<MANAGED_TOKEN_TEST>>();

    // Verify we have mint caps before destruction
    let mint_caps = managed_token::get_all_mint_caps(&state);
    assert!(mint_caps.length() == 2);

    // Destroy the managed token
    scenario.next_tx(OWNER);
    let (treasury_cap, deny_cap_option) = managed_token::destroy_managed_token(
        owner_cap,
        state,
        scenario.ctx()
    );

    // Verify deny cap is None (since we didn't initialize with one)
    assert!(deny_cap_option.is_none());

    // Verify we can still use the treasury cap
    let initial_supply = treasury_cap.total_supply();
    assert!(initial_supply == 0); // Should be 0 since we didn't mint anything

    // Clean up
    transfer::public_transfer(mint_cap1, MINTER);
    transfer::public_transfer(mint_cap2, OTHER_USER);
    transfer::public_transfer(treasury_cap, OWNER);
    transfer::public_freeze_object(coin_metadata);
    deny_cap_option.destroy_none();
    scenario.end();
}

#[test]
fun test_destroy_managed_token_with_deny_cap() {
    let mut scenario = test_scenario::begin(OWNER);
    let ctx = scenario.ctx();

    // Create the test token with deny cap
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
    let mut state = scenario.take_shared<TokenState<MANAGED_TOKEN_TEST>>();
    let owner_cap = scenario.take_from_sender<OwnerCap<MANAGED_TOKEN_TEST>>();

    // Configure some minters to test cleanup
    managed_token::configure_new_minter(
        &mut state,
        &owner_cap,
        MINTER,
        1500,
        false,
        scenario.ctx()
    );

    scenario.next_tx(MINTER);
    let mint_cap = scenario.take_from_sender<MintCap<MANAGED_TOKEN_TEST>>();

    // Verify we have mint caps before destruction
    let mint_caps = managed_token::get_all_mint_caps(&state);
    assert!(mint_caps.length() == 1);

    // Destroy the managed token
    scenario.next_tx(OWNER);
    let (returned_treasury_cap, deny_cap_option) = managed_token::destroy_managed_token(
        owner_cap,
        state,
        scenario.ctx()
    );

    // Verify deny cap is Some (since we initialized with one)
    assert!(deny_cap_option.is_some());
    let mut deny_cap_option = deny_cap_option;
    let returned_deny_cap = deny_cap_option.extract();

    // Verify we can still use the treasury cap
    let initial_supply = returned_treasury_cap.total_supply();
    assert!(initial_supply == 0);

    // Clean up
    transfer::public_transfer(mint_cap, MINTER);
    transfer::public_transfer(returned_treasury_cap, OWNER);
    transfer::public_transfer(returned_deny_cap, OWNER);
    transfer::public_freeze_object(coin_metadata);
    deny_cap_option.destroy_none();
    scenario.end();
}

#[test]
fun test_destroy_managed_token_with_minted_tokens() {
    let (mut scenario, mut state, owner_cap, coin_metadata) = setup_managed_token_test();

    // Configure a minter
    managed_token::configure_new_minter(
        &mut state,
        &owner_cap,
        MINTER,
        1000,
        false,
        scenario.ctx()
    );

    scenario.next_tx(MINTER);
    let mint_cap = scenario.take_from_sender<MintCap<MANAGED_TOKEN_TEST>>();

    // Create deny list for testing
    let deny_list = deny_list::new_for_testing(scenario.ctx());

    // Mint some tokens to increase total supply
    let mint_amount = 500;
    let coin = managed_token::mint(
        &mut state,
        &mint_cap,
        &deny_list,
        mint_amount,
        RECIPIENT,
        scenario.ctx()
    );

    // Verify total supply before destruction
    assert!(managed_token::total_supply(&state) == mint_amount);

    // Destroy the managed token
    scenario.next_tx(OWNER);
    let (treasury_cap, deny_cap_option) = managed_token::destroy_managed_token(
        owner_cap,
        state,
        scenario.ctx()
    );

    // Verify deny cap is None
    assert!(deny_cap_option.is_none());

    // Verify treasury cap still tracks the minted supply
    assert!(treasury_cap.total_supply() == mint_amount);

    // The minted coin should still be valid
    assert!(coin.value() == mint_amount);

    // Clean up
    transfer::public_transfer(coin, RECIPIENT);
    transfer::public_transfer(mint_cap, MINTER);
    transfer::public_transfer(treasury_cap, OWNER);
    transfer::public_freeze_object(coin_metadata);
    test_utils::destroy(deny_list);
    deny_cap_option.destroy_none();
    scenario.end();
}

#[test]
fun test_destroy_managed_token_multiple_mint_caps() {
    let (mut scenario, mut state, owner_cap, coin_metadata) = setup_managed_token_test();

    // Configure multiple minters with different allowances
    managed_token::configure_new_minter(
        &mut state,
        &owner_cap,
        MINTER,
        1000,
        false,
        scenario.ctx()
    );

    managed_token::configure_new_minter(
        &mut state,
        &owner_cap,
        OTHER_USER,
        0,
        true, // unlimited
        scenario.ctx()
    );

    managed_token::configure_new_minter(
        &mut state,
        &owner_cap,
        RECIPIENT,
        500,
        false,
        scenario.ctx()
    );

    // Get all mint caps
    scenario.next_tx(MINTER);
    let mint_cap1 = scenario.take_from_sender<MintCap<MANAGED_TOKEN_TEST>>();

    scenario.next_tx(OTHER_USER);
    let mint_cap2 = scenario.take_from_sender<MintCap<MANAGED_TOKEN_TEST>>();

    scenario.next_tx(RECIPIENT);
    let mint_cap3 = scenario.take_from_sender<MintCap<MANAGED_TOKEN_TEST>>();

    // Verify we have all mint caps before destruction
    let mint_caps = managed_token::get_all_mint_caps(&state);
    assert!(mint_caps.length() == 3);

    // Test allowances before destruction
    let (allowance1, unlimited1) = managed_token::mint_allowance(&state, object::id(&mint_cap1));
    let (_allowance2, unlimited2) = managed_token::mint_allowance(&state, object::id(&mint_cap2));
    let (allowance3, unlimited3) = managed_token::mint_allowance(&state, object::id(&mint_cap3));

    assert!(allowance1 == 1000 && !unlimited1);
    assert!(unlimited2);
    assert!(allowance3 == 500 && !unlimited3);

    // Destroy the managed token (should clean up all mint allowances)
    scenario.next_tx(OWNER);
    let (treasury_cap, deny_cap_option) = managed_token::destroy_managed_token(
        owner_cap,
        state,
        scenario.ctx()
    );

    // Verify deny cap is None
    assert!(deny_cap_option.is_none());

    // Verify treasury cap is returned correctly
    assert!(treasury_cap.total_supply() == 0);

    // Clean up
    transfer::public_transfer(mint_cap1, MINTER);
    transfer::public_transfer(mint_cap2, OTHER_USER);
    transfer::public_transfer(mint_cap3, RECIPIENT);
    transfer::public_transfer(treasury_cap, OWNER);
    transfer::public_freeze_object(coin_metadata);
    deny_cap_option.destroy_none();
    scenario.end();
}

#[test]
fun test_destroy_managed_token_after_ownership_transfer() {
    let (mut scenario, mut state, owner_cap, coin_metadata) = setup_managed_token_test();

    // Transfer ownership to OTHER_USER
    managed_token::transfer_ownership(
        &mut state,
        &owner_cap,
        OTHER_USER,
        scenario.ctx()
    );

    // Accept ownership
    scenario.next_tx(OTHER_USER);
    managed_token::accept_ownership(&mut state, scenario.ctx());

    // Execute ownership transfer
    scenario.next_tx(managed_token::owner(&state));
    managed_token::execute_ownership_transfer(
        owner_cap,
        managed_token::get_ownable_state(&mut state),
        OTHER_USER,
        scenario.ctx()
    );

    // Verify ownership changed
    assert!(managed_token::owner(&state) == OTHER_USER);

    scenario.next_tx(OTHER_USER);
    let owner_cap_other_user = scenario.take_from_sender<OwnerCap<MANAGED_TOKEN_TEST>>();
    let (treasury_cap, deny_cap_option) = managed_token::destroy_managed_token(
        owner_cap_other_user, // This should work because ownership was transferred
        state,
        scenario.ctx()
    );

    // Verify deny cap is None
    assert!(deny_cap_option.is_none());

    // Clean up
    transfer::public_transfer(treasury_cap, OTHER_USER);
    transfer::public_freeze_object(coin_metadata);
    deny_cap_option.destroy_none();
    scenario.end();
}
