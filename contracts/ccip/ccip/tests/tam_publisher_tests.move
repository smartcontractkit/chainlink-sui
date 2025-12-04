#[test_only]
module ccip::tam_publisher_tests;

use ccip::ownable::OwnerCap;
use ccip::publisher_wrapper;
use ccip::state_object::{Self, CCIPObjectRef};
use ccip::token_admin_registry as registry;
use ccip::upgrade_registry;
use std::string;
use std::type_name;
use sui::address;
use sui::coin;
use sui::coin_registry::{Self, Currency, CoinRegistry};
use sui::package;
use sui::test_scenario as ts;

/// OTW for test module
public struct TAM_PUBLISHER_TESTS has drop {}

/// Type identifier for the coin (not a one-time witness)
/// This struct has key ability for use with coin_registry::new_currency
public struct MyCoin has key { id: UID }

public struct TestTypeProof has drop {}
public struct TestTypeProof2 has drop {}

const ADMIN: address = @0xADDD;
const DECIMALS: u8 = 18;

fun setup_ccip_environment(scenario: &mut ts::Scenario): (OwnerCap, CCIPObjectRef) {
    scenario.next_tx(ADMIN);
    let ctx = scenario.ctx();
    state_object::test_init(ctx);

    scenario.next_tx(ADMIN);
    let owner_cap = scenario.take_from_sender<OwnerCap>();
    let mut ccip_ref = scenario.take_shared<CCIPObjectRef>();

    upgrade_registry::initialize(&mut ccip_ref, &owner_cap, scenario.ctx());
    registry::initialize(&mut ccip_ref, &owner_cap, scenario.ctx());

    (owner_cap, ccip_ref)
}

// ================================================================
// |        Token Admin Registry Integration Tests               |
// ================================================================

#[test]
/// Test that the registered token config contains the correct package address
public fun test_registered_type_proof_package_matches() {
    let mut scenario = ts::begin(ADMIN);
    let (owner_cap, mut ccip_ref) = setup_ccip_environment(&mut scenario);

    // Create a test token
    scenario.next_tx(ADMIN);
    let (treasury_cap, coin_metadata) = coin::create_currency(
        TAM_PUBLISHER_TESTS {},
        DECIMALS,
        b"TEST",
        b"TestToken",
        b"test_token",
        option::none(),
        scenario.ctx(),
    );
    let coin_metadata_address = object::id_to_address(&object::id(&coin_metadata));

    scenario.next_tx(ADMIN);
    {
        let publisher = package::test_claim(TAM_PUBLISHER_TESTS {}, scenario.ctx());
        let expected_package_address = address::from_ascii_bytes(publisher.package().as_bytes());
        let publisher_wrapper = publisher_wrapper::create(&publisher, TestTypeProof {});

        registry::register_pool(
            &mut ccip_ref,
            &treasury_cap,
            &coin_metadata,
            ADMIN,
            vector<address>[],
            vector<address>[],
            publisher_wrapper,
            TestTypeProof {},
        );

        let (
            registered_package_id,
            _module,
            _token_type,
            _admin,
            _pending_admin,
            _type_proof,
            _lock_params,
            _release_params,
        ) = registry::get_token_config_data(&ccip_ref, coin_metadata_address);

        assert!(registered_package_id == expected_package_address);

        package::burn_publisher(publisher);
    };

    transfer::public_transfer(treasury_cap, ADMIN);
    transfer::public_freeze_object(coin_metadata);
    transfer::public_transfer(owner_cap, ADMIN);
    ts::return_shared(ccip_ref);
    scenario.end();
}

// #[test]
// /// Test that the registered token config contains the correct package address
// /// using the new coin_registry::new_currency API instead of coin::create_currency
// public fun test_registered_type_proof_package_matches_with_currency() {
//     let mut scenario = ts::begin(ADMIN);
//     let (owner_cap, mut ccip_ref) = setup_ccip_environment(&mut scenario);

//     // Transaction 0: Create the CoinRegistry for testing (as system address @0x0)
//     scenario.next_tx(@0x0);
//     {
//         let registry = coin_registry::create_coin_data_registry_for_testing(scenario.ctx());
//         coin_registry::share_for_testing(registry);
//     };

//     // Transaction 1: Create currency using new_currency and finalize (shares the Currency)
//     scenario.next_tx(ADMIN);
//     let (treasury_cap, metadata_cap) = {
//         let mut registry = scenario.take_shared<CoinRegistry>();
//         let (initializer, treasury_cap) = coin_registry::new_currency<MyCoin>(
//             &mut registry,
//             DECIMALS,
//             string::utf8(b"TEST"),
//             string::utf8(b"TestToken"),
//             string::utf8(b"test_token"),
//             string::utf8(b""),
//             scenario.ctx(),
//         );
//         // finalize shares the Currency internally and returns MetadataCap
//         let metadata_cap = coin_registry::finalize(initializer, scenario.ctx());
//         ts::return_shared(registry);
//         (treasury_cap, metadata_cap)
//     };

//     // Transaction 2: Delete the metadata cap and get currency address
//     scenario.next_tx(ADMIN);
//     let currency_address = {
//         let mut currency = scenario.take_shared<Currency<MyCoin>>();
//         let addr = object::id_to_address(&object::id(&currency));
//         coin_registry::delete_metadata_cap(&mut currency, metadata_cap);
//         ts::return_shared(currency);
//         addr
//     };

//     // Transaction 3: Register pool with currency
//     scenario.next_tx(ADMIN);
//     {
//         let currency = scenario.take_shared<Currency<MyCoin>>();
//         let publisher = package::test_claim(TAM_PUBLISHER_TESTS {}, scenario.ctx());
//         let expected_package_address = address::from_ascii_bytes(publisher.package().as_bytes());
//         let publisher_wrapper = publisher_wrapper::create(&publisher, TestTypeProof {});

//         registry::register_pool_with_currency(
//             &mut ccip_ref,
//             &treasury_cap,
//             &currency,
//             ADMIN,
//             vector<address>[],
//             vector<address>[],
//             publisher_wrapper,
//             TestTypeProof {},
//         );

//         let (
//             registered_package_id,
//             _module,
//             _token_type,
//             _admin,
//             _pending_admin,
//             _type_proof,
//             _lock_params,
//             _release_params,
//         ) = registry::get_token_config_data(&ccip_ref, currency_address);

//         assert!(registered_package_id == expected_package_address);

//         package::burn_publisher(publisher);
//         ts::return_shared(currency);
//     };

//     transfer::public_transfer(treasury_cap, ADMIN);
//     transfer::public_transfer(owner_cap, ADMIN);
//     ts::return_shared(ccip_ref);
//     scenario.end();
// }

#[test]
/// Test full registration flow with publisher wrapper validation
public fun test_register_pool_with_valid_publisher_wrapper() {
    let mut scenario = ts::begin(ADMIN);
    let (owner_cap, mut ccip_ref) = setup_ccip_environment(&mut scenario);

    scenario.next_tx(ADMIN);
    let (treasury_cap, coin_metadata) = coin::create_currency(
        TAM_PUBLISHER_TESTS {},
        DECIMALS,
        b"TEST",
        b"TestToken",
        b"test_token",
        option::none(),
        scenario.ctx(),
    );
    let coin_metadata_address = object::id_to_address(&object::id(&coin_metadata));

    scenario.next_tx(ADMIN);
    {
        let publisher = package::test_claim(TAM_PUBLISHER_TESTS {}, scenario.ctx());
        let expected_package_address = address::from_ascii_bytes(publisher.package().as_bytes());
        let publisher_wrapper = publisher_wrapper::create(&publisher, TestTypeProof {});

        registry::register_pool(
            &mut ccip_ref,
            &treasury_cap,
            &coin_metadata,
            ADMIN,
            vector<address>[@0x1, @0x2],
            vector<address>[@0x3, @0x4],
            publisher_wrapper,
            TestTypeProof {},
        );

        let pool_address = registry::get_pool(&ccip_ref, coin_metadata_address);
        assert!(pool_address == expected_package_address);

        let (
            package_id,
            module_name,
            _token_type,
            administrator,
            pending_admin,
            type_proof,
            lock_params,
            release_params,
        ) = registry::get_token_config_data(&ccip_ref, coin_metadata_address);

        assert!(package_id == expected_package_address);
        assert!(module_name == string::utf8(b"tam_publisher_tests"));
        assert!(administrator == ADMIN);
        assert!(pending_admin == @0x0);
        assert!(
            type_proof == type_name::into_string(type_name::with_defining_ids<TestTypeProof>()),
        );
        assert!(lock_params.length() == 2);
        assert!(release_params.length() == 2);

        package::burn_publisher(publisher);
    };

    transfer::public_transfer(treasury_cap, ADMIN);
    transfer::public_freeze_object(coin_metadata);
    transfer::public_transfer(owner_cap, ADMIN);
    ts::return_shared(ccip_ref);
    scenario.end();
}

#[test]
/// Test that publisher wrapper can be created and used for registration
/// This verifies that the package address is correctly extracted and stored
public fun test_publisher_wrapper_stores_correct_package_address() {
    let mut scenario = ts::begin(ADMIN);
    let (owner_cap, mut ccip_ref) = setup_ccip_environment(&mut scenario);

    scenario.next_tx(ADMIN);
    let (treasury_cap, coin_metadata) = coin::create_currency(
        TAM_PUBLISHER_TESTS {},
        DECIMALS,
        b"TEST",
        b"TestToken",
        b"test_token",
        option::none(),
        scenario.ctx(),
    );
    let coin_metadata_address = object::id_to_address(&object::id(&coin_metadata));

    scenario.next_tx(ADMIN);
    {
        let publisher = package::test_claim(TAM_PUBLISHER_TESTS {}, scenario.ctx());
        let expected_package_address = address::from_ascii_bytes(publisher.package().as_bytes());

        let publisher_wrapper = publisher_wrapper::create(&publisher, TestTypeProof {});

        registry::register_pool(
            &mut ccip_ref,
            &treasury_cap,
            &coin_metadata,
            ADMIN,
            vector<address>[],
            vector<address>[],
            publisher_wrapper,
            TestTypeProof {},
        );

        let (stored_address, _, _, _, _, _, _, _) = registry::get_token_config_data(
            &ccip_ref,
            coin_metadata_address,
        );
        assert!(stored_address == expected_package_address);

        package::burn_publisher(publisher);
    };

    transfer::public_transfer(treasury_cap, ADMIN);
    transfer::public_freeze_object(coin_metadata);
    transfer::public_transfer(owner_cap, ADMIN);
    ts::return_shared(ccip_ref);
    scenario.end();
}
