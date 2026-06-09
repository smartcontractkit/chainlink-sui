#[allow(implicit_const_copy)]
#[test_only]
module ccip::sui_fee_billing_tests;

use ccip::client;
use ccip::fee_quoter;
use ccip::ownable::OwnerCap;
use ccip::state_object::{Self, CCIPObjectRef};
use ccip::upgrade_registry;
use sui::clock;
use sui::test_scenario;

const CHAIN_FAMILY_SELECTOR_SUI: vector<u8> = x"c4e05953";
const ONE_E_18: u256 = 1_000_000_000_000_000_000;
const FEE_TOKEN: address =
    @0x1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b;
const TOKEN_2: address = @0x000000000000000000000000F4030086522a5bEEa4988F8cA5B36dbC97BeE88c;
const SUI_RECEIVER: vector<u8> =
    x"1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef";
const TOKEN_RECEIVER: vector<u8> =
    x"aabbccdd11223344aabbccdd11223344aabbccdd11223344aabbccdd11223344";
const OBJ1: vector<u8> =
    x"0101010101010101010101010101010101010101010101010101010101010101";
const OBJ2: vector<u8> =
    x"0202020202020202020202020202020202020202020202020202020202020202";
const OBJ3: vector<u8> =
    x"0303030303030303030303030303030303030303030303030303030303030303";
const OBJ4: vector<u8> =
    x"0404040404040404040404040404040404040404040404040404040404040404";
const OBJ5: vector<u8> =
    x"0505050505050505050505050505050505050505050505050505050505050505";
const OBJ6: vector<u8> =
    x"0606060606060606060606060606060606060606060606060606060606060606";
const OBJ7: vector<u8> =
    x"0707070707070707070707070707070707070707070707070707070707070707";
const OBJ8: vector<u8> =
    x"0808080808080808080808080808080808080808080808080808080808080808";
const OBJ9: vector<u8> =
    x"0909090909090909090909090909090909090909090909090909090909090909";
const OBJ10: vector<u8> =
    x"0a0a0a0a0a0a0a0a0a0a0a0a0a0a0a0a0a0a0a0a0a0a0a0a0a0a0a0a0a0a0a0a";
const ADMIN: address = @0x1;

fun setup(): (test_scenario::Scenario, OwnerCap, CCIPObjectRef) {
    let mut scenario = test_scenario::begin(ADMIN);
    let ctx = scenario.ctx();
    state_object::test_init(ctx);
    scenario.next_tx(ADMIN);
    let owner_cap = scenario.take_from_sender<OwnerCap>();
    let ref = scenario.take_shared<CCIPObjectRef>();
    (scenario, owner_cap, ref)
}

fun init_fee_quoter(ref: &mut CCIPObjectRef, cap: &OwnerCap, ctx: &mut TxContext) {
    upgrade_registry::initialize(ref, cap, ctx);
    fee_quoter::initialize(
        ref,
        cap,
        200 * ONE_E_18,
        FEE_TOKEN,
        1000,
        vector[FEE_TOKEN, TOKEN_2],
        ctx,
    );
}

fun setup_sui_dest(ref: &mut CCIPObjectRef, cap: &OwnerCap, ctx: &mut TxContext) {
    fee_quoter::apply_dest_chain_config_updates(
        ref,
        cap,
        100,
        true,
        10,
        30_000,
        3_000_000,
        300_000,
        16,
        40,
        300,
        100,
        16,
        600,
        CHAIN_FAMILY_SELECTOR_SUI,
        false,
        50,
        90_000,
        200_000,
        ONE_E_18 as u64,
        1_000_000,
        50,
        ctx,
    );
}

fun setup_prices(
    ref: &mut CCIPObjectRef,
    fq_cap: &fee_quoter::FeeQuoterCap,
    clock: &clock::Clock,
    ctx: &mut TxContext,
) {
    fee_quoter::update_prices(
        ref,
        fq_cap,
        clock,
        vector[FEE_TOKEN, TOKEN_2],
        vector[150_000_000_000 * ONE_E_18, 150_000_000_000 * ONE_E_18],
        vector[100],
        vector[7_500_000_000_000],
        ctx,
    );
}

fun cleanup(scenario: test_scenario::Scenario, cap: OwnerCap, ref: CCIPObjectRef) {
    test_scenario::return_to_sender(&scenario, cap);
    test_scenario::return_shared(ref);
    test_scenario::end(scenario);
}

fun many_object_ids(count: u64): vector<vector<u8>> {
    let mut ids = vector[];
    let mut i = 0;
    while (i < count) {
        ids.push_back(x"2234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdea");
        i = i + 1;
    };
    ids
}

#[test]
public fun test_sui_fee_increases_with_receiver_object_id_count() {
    let (mut scenario, cap, mut ref) = setup();
    let ctx = scenario.ctx();
    init_fee_quoter(&mut ref, &cap, ctx);
    let fq_cap = fee_quoter::new_fee_quoter_cap(&ref, &cap, ctx);
    let mut clock = clock::create_for_testing(ctx);
    clock::increment_for_testing(&mut clock, 20000);
    setup_prices(&mut ref, &fq_cap, &clock, ctx);
    setup_sui_dest(&mut ref, &cap, ctx);
    fee_quoter::apply_premium_multiplier_wei_per_eth_updates(
        &mut ref,
        &cap,
        vector[FEE_TOKEN],
        vector[ONE_E_18 as u64],
        ctx,
    );

    let fee_0_ids = fee_quoter::get_validated_fee(
        &ref,
        &clock,
        100,
        SUI_RECEIVER,
        b"test_payload",
        vector[],
        vector[],
        FEE_TOKEN,
        client::encode_sui_extra_args_v1(200_000, false, TOKEN_RECEIVER, vector[]),
    );
    let fee_10_ids = fee_quoter::get_validated_fee(
        &ref,
        &clock,
        100,
        SUI_RECEIVER,
        b"test_payload",
        vector[],
        vector[],
        FEE_TOKEN,
        client::encode_sui_extra_args_v1(
            200_000,
            false,
            TOKEN_RECEIVER,
            vector[OBJ1, OBJ2, OBJ3, OBJ4, OBJ5, OBJ6, OBJ7, OBJ8, OBJ9, OBJ10],
        ),
    );

    assert!(fee_0_ids > 0, 1);
    assert!(fee_10_ids > fee_0_ids, 2);

    fee_quoter::destroy_fee_quoter_cap(&ref, &cap, fq_cap);
    clock::destroy_for_testing(clock);
    cleanup(scenario, cap, ref);
}

#[test]
#[expected_failure(abort_code = fee_quoter::ETooManySuiExtraArgsReceiverObjectIds)]
public fun test_sui_fee_reverts_when_too_many_receiver_object_ids() {
    let (mut scenario, cap, mut ref) = setup();
    let ctx = scenario.ctx();
    init_fee_quoter(&mut ref, &cap, ctx);
    let fq_cap = fee_quoter::new_fee_quoter_cap(&ref, &cap, ctx);
    let mut clock = clock::create_for_testing(ctx);
    clock::increment_for_testing(&mut clock, 20000);
    setup_prices(&mut ref, &fq_cap, &clock, ctx);
    setup_sui_dest(&mut ref, &cap, ctx);
    fee_quoter::apply_premium_multiplier_wei_per_eth_updates(
        &mut ref,
        &cap,
        vector[FEE_TOKEN],
        vector[ONE_E_18 as u64],
        ctx,
    );

    fee_quoter::get_validated_fee(
        &ref,
        &clock,
        100,
        SUI_RECEIVER,
        b"test_payload",
        vector[],
        vector[],
        FEE_TOKEN,
        client::encode_sui_extra_args_v1(200_000, false, TOKEN_RECEIVER, many_object_ids(65)),
    );

    fee_quoter::destroy_fee_quoter_cap(&ref, &cap, fq_cap);
    clock::destroy_for_testing(clock);
    cleanup(scenario, cap, ref);
}

#[test]
#[expected_failure(abort_code = fee_quoter::ETooManySuiExtraArgsReceiverObjectIds)]
public fun test_sui_fee_reverts_zero_receiver_with_object_ids() {
    let (mut scenario, cap, mut ref) = setup();
    let ctx = scenario.ctx();
    init_fee_quoter(&mut ref, &cap, ctx);
    let fq_cap = fee_quoter::new_fee_quoter_cap(&ref, &cap, ctx);
    let mut clock = clock::create_for_testing(ctx);
    clock::increment_for_testing(&mut clock, 20000);
    setup_prices(&mut ref, &fq_cap, &clock, ctx);
    setup_sui_dest(&mut ref, &cap, ctx);
    fee_quoter::apply_premium_multiplier_wei_per_eth_updates(
        &mut ref,
        &cap,
        vector[FEE_TOKEN],
        vector[ONE_E_18 as u64],
        ctx,
    );

    let zero_receiver = x"0000000000000000000000000000000000000000000000000000000000000000";
    fee_quoter::get_validated_fee(
        &ref,
        &clock,
        100,
        zero_receiver,
        b"test_payload",
        vector[],
        vector[],
        FEE_TOKEN,
        client::encode_sui_extra_args_v1(200_000, false, TOKEN_RECEIVER, vector[OBJ1]),
    );

    fee_quoter::destroy_fee_quoter_cap(&ref, &cap, fq_cap);
    clock::destroy_for_testing(clock);
    cleanup(scenario, cap, ref);
}

#[test]
#[expected_failure(abort_code = fee_quoter::EMessageTooLarge)]
public fun test_sui_expanded_payload_respects_max_data_bytes() {
    let (mut scenario, cap, mut ref) = setup();
    let ctx = scenario.ctx();
    init_fee_quoter(&mut ref, &cap, ctx);
    let fq_cap = fee_quoter::new_fee_quoter_cap(&ref, &cap, ctx);
    let mut clock = clock::create_for_testing(ctx);
    clock::increment_for_testing(&mut clock, 20000);
    setup_prices(&mut ref, &fq_cap, &clock, ctx);

    // Configure with a low max_data_bytes (512) so object IDs push it over the limit
    fee_quoter::apply_dest_chain_config_updates(
        &mut ref,
        &cap,
        100,
        true,
        10,
        512,
        3_000_000,
        300_000,
        16,
        40,
        300,
        100,
        16,
        600,
        CHAIN_FAMILY_SELECTOR_SUI,
        false,
        50,
        90_000,
        200_000,
        ONE_E_18 as u64,
        1_000_000,
        50,
        ctx,
    );

    fee_quoter::apply_premium_multiplier_wei_per_eth_updates(
        &mut ref,
        &cap,
        vector[FEE_TOKEN],
        vector[ONE_E_18 as u64],
        ctx,
    );

    // 10 object IDs = (10 + 1) * 32 = 352 bytes overhead + 12 bytes data = 364. Should pass.
    // But 15 IDs = (15 + 1) * 32 = 512 + 12 = 524 > 512. Should fail.
    fee_quoter::get_validated_fee(
        &ref,
        &clock,
        100,
        SUI_RECEIVER,
        b"test_payload",
        vector[],
        vector[],
        FEE_TOKEN,
        client::encode_sui_extra_args_v1(200_000, false, TOKEN_RECEIVER, many_object_ids(15)),
    );

    fee_quoter::destroy_fee_quoter_cap(&ref, &cap, fq_cap);
    clock::destroy_for_testing(clock);
    cleanup(scenario, cap, ref);
}
