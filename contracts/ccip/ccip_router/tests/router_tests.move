#[test_only]
module ccip_router::router_tests;

use ccip_router::ownable::{Self, OwnerCap};
use ccip_router::router::{Self, RouterState, RouterObject};
use sui::derived_object;
use sui::test_scenario::{Self as ts, Scenario};

const SENDER_1: address = @0x1;

const ETH_CHAIN_SELECTOR: u64 = 5009297550715157269;
const AVAX_CHAIN_SELECTOR: u64 = 6433500567565415381;
const BSC_CHAIN_SELECTOR: u64 = 4380317901350075273;
const ARBITRARY_CHAIN_SELECTOR: u64 = 123456789;
const ETH_ON_RAMP_ADDRESS: address = @0x111;
const AVAX_ON_RAMP_ADDRESS: address = @0x222;

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
fun test_set_and_get_on_ramps() {
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

        router::set_on_ramps(
            &owner_cap,
            &mut router,
            dest_chain_selectors,
            on_ramp_addresses,
        );

        // Test individual chain support
        assert!(router::is_chain_supported(&router, ETH_CHAIN_SELECTOR));
        assert!(router::is_chain_supported(&router, AVAX_CHAIN_SELECTOR));
        assert!(!router::is_chain_supported(&router, BSC_CHAIN_SELECTOR));

        // Test getting on ramp addresses
        let eth_on_ramp = router::get_on_ramp(&router, ETH_CHAIN_SELECTOR);
        let avax_on_ramp = router::get_on_ramp(&router, AVAX_CHAIN_SELECTOR);
        assert!(eth_on_ramp == ETH_ON_RAMP_ADDRESS);
        assert!(avax_on_ramp == AVAX_ON_RAMP_ADDRESS);

        scenario.return_to_sender(owner_cap);
        ts::return_shared(router);
    };

    ts::end(scenario);
}

#[test]
fun test_update_on_ramp() {
    let mut scenario = create_test_scenario();

    {
        let ctx = scenario.ctx();
        router::test_init(ctx);
    };

    scenario.next_tx(SENDER_1);
    {
        let mut router = scenario.take_shared<RouterState>();
        let owner_cap = scenario.take_from_sender<OwnerCap>();

        // First add ETH chain
        router::set_on_ramps(
            &owner_cap,
            &mut router,
            vector[ETH_CHAIN_SELECTOR],
            vector[ETH_ON_RAMP_ADDRESS],
        );

        // Verify it was added
        assert!(router::is_chain_supported(&router, ETH_CHAIN_SELECTOR));
        let eth_on_ramp = router::get_on_ramp(&router, ETH_CHAIN_SELECTOR);
        assert!(eth_on_ramp == ETH_ON_RAMP_ADDRESS);

        // Now update ETH chain to use AVAX address
        router::set_on_ramps(
            &owner_cap,
            &mut router,
            vector[ETH_CHAIN_SELECTOR],
            vector[AVAX_ON_RAMP_ADDRESS],
        );

        // Verify it was updated
        assert!(router::is_chain_supported(&router, ETH_CHAIN_SELECTOR));
        let updated_eth_on_ramp = router::get_on_ramp(&router, ETH_CHAIN_SELECTOR);
        assert!(updated_eth_on_ramp == AVAX_ON_RAMP_ADDRESS);

        scenario.return_to_sender(owner_cap);
        ts::return_shared(router);
    };

    ts::end(scenario);
}

#[test]
#[expected_failure(abort_code = router::EInvalidOnrampAddress)]
fun test_set_zero_address_fails() {
    let mut scenario = create_test_scenario();

    {
        let ctx = scenario.ctx();
        router::test_init(ctx);
    };

    scenario.next_tx(SENDER_1);
    {
        let mut router = scenario.take_shared<RouterState>();
        let owner_cap = scenario.take_from_sender<OwnerCap>();

        // Try to set zero address - should fail
        router::set_on_ramps(
            &owner_cap,
            &mut router,
            vector[ETH_CHAIN_SELECTOR],
            vector[@0x0],
        );

        scenario.return_to_sender(owner_cap);
        ts::return_shared(router);
    };

    ts::end(scenario);
}

#[test]
#[expected_failure(abort_code = router::EParamsLengthMismatch)]
fun test_set_on_ramps_length_mismatch() {
    let mut scenario = create_test_scenario();

    {
        let ctx = scenario.ctx();
        router::test_init(ctx);
    };

    scenario.next_tx(SENDER_1);
    {
        let mut router = scenario.take_shared<RouterState>();
        let owner_cap = scenario.take_from_sender<OwnerCap>();

        // Mismatched lengths should fail
        router::set_on_ramps(
            &owner_cap,
            &mut router,
            vector[ETH_CHAIN_SELECTOR, AVAX_CHAIN_SELECTOR], // 2 selectors
            vector[ETH_ON_RAMP_ADDRESS], // 1 address
        );

        scenario.return_to_sender(owner_cap);
        ts::return_shared(router);
    };

    ts::end(scenario);
}

#[test]
#[expected_failure(abort_code = router::EOnrampNotFound)]
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
        router::get_on_ramp(&router, ARBITRARY_CHAIN_SELECTOR);
        ts::return_shared(router);
    };

    ts::end(scenario);
}

#[test]
fun test_derive_address() {
    let mut scenario = create_test_scenario();
    let ctx = scenario.ctx();
    router::test_init(ctx);

    scenario.next_tx(SENDER_1);
    let router_object = scenario.take_shared<RouterObject>();

    // Test OwnerCap derivation
    let derived_owner_cap_addr = derived_object::derive_address(
        object::id(&router_object),
        ownable::default_key(),
    );
    let owner_cap = scenario.take_from_sender<OwnerCap>();
    let owner_cap_id = object::id(&owner_cap).to_address();

    assert!(derived_owner_cap_addr == owner_cap_id);
    assert!(
        derived_owner_cap_addr == @0x8e574462de77f45ea5bf3e8c1da19bcf081d25796376367e311926a8f993177e,
    );

    // Test RouterState derivation
    let derived_router_state_addr = derived_object::derive_address(
        object::id(&router_object),
        b"RouterState",
    );
    let router_state = scenario.take_shared<RouterState>();
    let router_state_id = object::id(&router_state).to_address();
    assert!(derived_router_state_addr == router_state_id);
    assert!(
        derived_router_state_addr == @0xc2ab753588210ab5de22dca1caf6e6d18a0b514c28c1975655c5769117d6f9ef,
    );

    ts::return_to_address(SENDER_1, owner_cap);
    ts::return_shared(router_state);
    ts::return_shared(router_object);
    ts::end(scenario);
}
