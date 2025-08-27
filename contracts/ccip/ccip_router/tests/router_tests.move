#[test_only]
module ccip_router::router_tests;

use ccip_router::ownable::OwnerCap;
use ccip_router::router::{Self, RouterState};
use sui::test_scenario::{Self as ts, Scenario};

const SENDER_1: address = @0x1;
const SENDER_2: address = @0x2;

const ETH_CHAIN_SELECTOR: u64 = 5009297550715157269;
const AVAX_CHAIN_SELECTOR: u64 = 6433500567565415381;
const BSC_CHAIN_SELECTOR: u64 = 4380317901350075273;
const ARBITRARY_CHAIN_SELECTOR: u64 = 123456789;
const ETH_ON_RAMP_ADDRESS: address = @0x111;
const AVAX_ON_RAMP_ADDRESS: address = @0x222;

const VERSION_1_6_0: vector<u8> = vector[1, 6, 0];
const INVALID_VERSION: vector<u8> = vector[1, 2]; // Invalid because it has 2 elements, not 3

fun create_test_scenario(): Scenario {
    ts::begin(SENDER_1)
}

#[test]
fun test_initialization() {
    let mut scenario = create_test_scenario();

    {
        let ctx = scenario.ctx();
        router::test_init(ctx);
    };

    {
        scenario.next_tx(@0xB);
        let router = scenario.take_shared<RouterState>();
        assert!(router::type_and_version() == std::string::utf8(b"Router 1.6.0"));
        ts::return_shared(router);
    };

    ts::end(scenario);
}

#[test]
fun test_set_and_get_on_ramp_infos() {
    let mut scenario = create_test_scenario();

    {
        let ctx = scenario.ctx();
        router::test_init(ctx);
    };

    scenario.next_tx(SENDER_1);
    {
        let mut router = scenario.take_shared<RouterState>();
        let owner_cap = scenario.take_from_sender<OwnerCap>();

        let dest_chain_selectors = vector[ETH_CHAIN_SELECTOR, AVAX_CHAIN_SELECTOR];
        let on_ramp_addresses = vector[ETH_ON_RAMP_ADDRESS, AVAX_ON_RAMP_ADDRESS];
        let on_ramp_versions = vector[VERSION_1_6_0, VERSION_1_6_0];

        router::set_on_ramp_infos(
            &owner_cap,
            &mut router,
            dest_chain_selectors,
            on_ramp_addresses,
            on_ramp_versions,
        );

        let infos = router::get_on_ramp_infos(&router, dest_chain_selectors);
        assert!(infos.length() == 2);
        assert!(router::get_on_ramp_address(infos[0]) == ETH_ON_RAMP_ADDRESS);
        assert!(router::get_on_ramp_version(infos[0]) == VERSION_1_6_0);
        assert!(router::get_on_ramp_address(infos[1]) == AVAX_ON_RAMP_ADDRESS);
        assert!(router::get_on_ramp_version(infos[1]) == VERSION_1_6_0);

        scenario.return_to_sender(owner_cap);
        ts::return_shared(router);
    };

    // a random user should not be able to get on ramp info
    scenario.next_tx(SENDER_2);
    {
        let router = scenario.take_shared<RouterState>();

        let non_existent_chain = vector[BSC_CHAIN_SELECTOR];
        let non_existent_infos = router::get_on_ramp_infos(&router, non_existent_chain);
        assert!(non_existent_infos.length() == 1);
        assert!(router::get_on_ramp_address(non_existent_infos[0]) == @0x0);
        assert!(router::get_on_ramp_version(non_existent_infos[0]).is_empty());

        // Test get_on_ramp_versions with mixed existing and non-existing chains
        let mixed_chains = vector[ETH_CHAIN_SELECTOR, BSC_CHAIN_SELECTOR, AVAX_CHAIN_SELECTOR];
        let mixed_infos = router::get_on_ramp_infos(&router, mixed_chains);
        assert!(mixed_infos.length() == 3);
        assert!(router::get_on_ramp_address(mixed_infos[0]) == ETH_ON_RAMP_ADDRESS);
        assert!(router::get_on_ramp_version(mixed_infos[0]) == VERSION_1_6_0);
        assert!(router::get_on_ramp_address(mixed_infos[1]) == @0x0);
        assert!(router::get_on_ramp_version(mixed_infos[1]).is_empty());
        assert!(router::get_on_ramp_address(mixed_infos[2]) == AVAX_ON_RAMP_ADDRESS);
        assert!(router::get_on_ramp_version(mixed_infos[2]) == VERSION_1_6_0);

        ts::return_shared(router);
    };

    ts::end(scenario);
}

#[test]
fun test_remove_on_ramp_info() {
    let mut scenario = create_test_scenario();

    {
        let ctx = scenario.ctx();
        router::test_init(ctx);
    };

    scenario.next_tx(SENDER_1);
    {
        let mut router = scenario.take_shared<RouterState>();
        let owner_cap = scenario.take_from_sender<OwnerCap>();

        let dest_chain_selectors = vector[ETH_CHAIN_SELECTOR, AVAX_CHAIN_SELECTOR];
        let on_ramp_addresses = vector[ETH_ON_RAMP_ADDRESS, AVAX_ON_RAMP_ADDRESS];
        let on_ramp_versions = vector[VERSION_1_6_0, VERSION_1_6_0];
        router::set_on_ramp_infos(
            &owner_cap,
            &mut router,
            dest_chain_selectors,
            on_ramp_addresses,
            on_ramp_versions,
        );

        // Verify they were added
        assert!(router::is_chain_supported(&router, ETH_CHAIN_SELECTOR));
        assert!(router::is_chain_supported(&router, AVAX_CHAIN_SELECTOR));
        scenario.return_to_sender(owner_cap);
        ts::return_shared(router);
    };

    scenario.next_tx(SENDER_1);
    {
        let mut router = scenario.take_shared<RouterState>();
        let owner_cap = scenario.take_from_sender<OwnerCap>();

        let dest_chain_selectors = vector[ETH_CHAIN_SELECTOR, AVAX_CHAIN_SELECTOR];
        let on_ramp_addresses = vector[ETH_ON_RAMP_ADDRESS, AVAX_ON_RAMP_ADDRESS];
        let on_ramp_versions = vector[VERSION_1_6_0, VERSION_1_6_0];
        router::set_on_ramp_infos(
            &owner_cap,
            &mut router,
            dest_chain_selectors,
            on_ramp_addresses,
            on_ramp_versions,
        );

        // Verify they were added
        assert!(router::is_chain_supported(&router, ETH_CHAIN_SELECTOR));
        assert!(router::is_chain_supported(&router, AVAX_CHAIN_SELECTOR));

        // Now remove one of them by setting an empty version
        let remove_selectors = vector[ETH_CHAIN_SELECTOR];
        let remove_versions = vector[vector[]]; // Empty version removes the chain
        router::set_on_ramp_infos(
            &owner_cap,
            &mut router,
            remove_selectors,
            vector[@0x0],
            remove_versions,
        );

        // Verify it was removed
        assert!(!router::is_chain_supported(&router, ETH_CHAIN_SELECTOR));
        assert!(router::is_chain_supported(&router, AVAX_CHAIN_SELECTOR)); // This one should still exist

        // Check with get_on_ramp_versions
        let check_selectors = vector[ETH_CHAIN_SELECTOR, AVAX_CHAIN_SELECTOR];
        let infos = router::get_on_ramp_infos(&router, check_selectors);
        assert!(infos.length() == 2);
        assert!(router::get_on_ramp_address(infos[0]) == @0x0); // ETH should be removed
        assert!(router::get_on_ramp_version(infos[0]).is_empty());
        assert!(router::get_on_ramp_address(infos[1]) == AVAX_ON_RAMP_ADDRESS);
        assert!(router::get_on_ramp_version(infos[1]) == VERSION_1_6_0);

        scenario.return_to_sender(owner_cap);
        ts::return_shared(router);
    };

    ts::end(scenario);
}

#[test]
#[expected_failure(abort_code = router::EInvalidOnrampVersion)]
fun test_set_invalid_on_ramp_version() {
    let mut scenario = create_test_scenario();

    {
        let ctx = scenario.ctx();
        router::test_init(ctx);
    };

    scenario.next_tx(SENDER_1);
    {
        let mut router = scenario.take_shared<RouterState>();
        let owner_cap = scenario.take_from_sender<OwnerCap>();

        router::set_on_ramp_infos(
            &owner_cap,
            &mut router,
            vector[ETH_CHAIN_SELECTOR],
            vector[ETH_ON_RAMP_ADDRESS],
            vector[INVALID_VERSION],
        );

        scenario.return_to_sender(owner_cap);
        ts::return_shared(router);
    };

    ts::end(scenario);
}

#[test]
#[expected_failure(abort_code = router::EOnrampInfoNotFound)]
fun test_get_on_ramp_unsupported_chain() {
    let mut scenario = create_test_scenario();

    {
        let ctx = scenario.ctx();
        router::test_init(ctx);
    };

    scenario.next_tx(@0xB);
    {
        let router = scenario.take_shared<RouterState>();
        // This should fail because the chain is not supported
        router::get_on_ramp_info(&router, ARBITRARY_CHAIN_SELECTOR);
        ts::return_shared(router);
    };

    ts::end(scenario);
}
