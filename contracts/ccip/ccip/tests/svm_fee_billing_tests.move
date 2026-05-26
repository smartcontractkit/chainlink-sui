#[allow(implicit_const_copy)]
#[test_only]
module ccip::svm_fee_billing_tests;

use ccip::client;
use ccip::fee_quoter;
use ccip::ownable::OwnerCap;
use ccip::state_object::{Self, CCIPObjectRef};
use ccip::upgrade_registry;
use sui::clock;
use sui::test_scenario;

const CHAIN_FAMILY_SELECTOR_SVM: vector<u8> = x"1e10bdc4";
const ONE_E_18: u256 = 1_000_000_000_000_000_000;
const FEE_TOKEN: address =
    @0x1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b;
const TOKEN_2: address = @0x000000000000000000000000F4030086522a5bEEa4988F8cA5B36dbC97BeE88c;
const TOKEN_3: address = @0x8a7b6c5d4e3f2a1b0c9d8e7f6a5b4c3d2e1f0a9b8c7d6e5f4a3b2c1d0e9f8a7;
const SVM_RECEIVER: vector<u8> =
    x"1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef";
const TOKEN_RECEIVER: vector<u8> =
    x"aabbccdd11223344aabbccdd11223344aabbccdd11223344aabbccdd11223344";
const ACC1: vector<u8> =
    x"0101010101010101010101010101010101010101010101010101010101010101";
const ACC2: vector<u8> =
    x"0202020202020202020202020202020202020202020202020202020202020202";
const ACC3: vector<u8> =
    x"0303030303030303030303030303030303030303030303030303030303030303";
const ACC4: vector<u8> =
    x"0404040404040404040404040404040404040404040404040404040404040404";
const ACC5: vector<u8> =
    x"0505050505050505050505050505050505050505050505050505050505050505";
const ACC6: vector<u8> =
    x"0606060606060606060606060606060606060606060606060606060606060606";
const ACC7: vector<u8> =
    x"0707070707070707070707070707070707070707070707070707070707070707";
const ACC8: vector<u8> =
    x"0808080808080808080808080808080808080808080808080808080808080808";
const ACC9: vector<u8> =
    x"0909090909090909090909090909090909090909090909090909090909090909";
const ACC10: vector<u8> =
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
        vector[FEE_TOKEN, TOKEN_2, TOKEN_3],
        ctx,
    );
}

fun setup_svm_dest(ref: &mut CCIPObjectRef, cap: &OwnerCap, ctx: &mut TxContext) {
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
        CHAIN_FAMILY_SELECTOR_SVM,
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

#[test]
public fun test_svm_fee_increases_with_account_count() {
    let (mut scenario, cap, mut ref) = setup();
    let ctx = scenario.ctx();
    init_fee_quoter(&mut ref, &cap, ctx);
    let fq_cap = fee_quoter::new_fee_quoter_cap(&ref, &cap, ctx);
    let mut clock = clock::create_for_testing(ctx);
    clock::increment_for_testing(&mut clock, 20000);
    setup_prices(&mut ref, &fq_cap, &clock, ctx);
    setup_svm_dest(&mut ref, &cap, ctx);
    fee_quoter::apply_premium_multiplier_wei_per_eth_updates(
        &mut ref,
        &cap,
        vector[FEE_TOKEN],
        vector[ONE_E_18 as u64],
        ctx,
    );

    let fee_0_accounts = fee_quoter::get_validated_fee(
        &ref,
        &clock,
        100,
        SVM_RECEIVER,
        b"test_payload",
        vector[],
        vector[],
        FEE_TOKEN,
        client::encode_svm_extra_args_v1(200_000, 0, false, TOKEN_RECEIVER, vector[]),
    );
    let fee_10_accounts = fee_quoter::get_validated_fee(
        &ref,
        &clock,
        100,
        SVM_RECEIVER,
        b"test_payload",
        vector[],
        vector[],
        FEE_TOKEN,
        client::encode_svm_extra_args_v1(
            200_000,
            0,
            false,
            TOKEN_RECEIVER,
            vector[ACC1, ACC2, ACC3, ACC4, ACC5, ACC6, ACC7, ACC8, ACC9, ACC10],
        ),
    );

    assert!(fee_0_accounts > 0, 1);
    assert!(fee_10_accounts > fee_0_accounts, 2);

    fee_quoter::destroy_fee_quoter_cap(&ref, &cap, fq_cap);
    clock::destroy_for_testing(clock);
    cleanup(scenario, cap, ref);
}

#[test]
public fun test_svm_fee_increases_with_token_transfer_overhead() {
    let (mut scenario, cap, mut ref) = setup();
    let ctx = scenario.ctx();
    init_fee_quoter(&mut ref, &cap, ctx);
    let fq_cap = fee_quoter::new_fee_quoter_cap(&ref, &cap, ctx);
    let mut clock = clock::create_for_testing(ctx);
    clock::increment_for_testing(&mut clock, 20000);
    setup_prices(&mut ref, &fq_cap, &clock, ctx);
    setup_svm_dest(&mut ref, &cap, ctx);
    fee_quoter::apply_premium_multiplier_wei_per_eth_updates(
        &mut ref,
        &cap,
        vector[FEE_TOKEN],
        vector[ONE_E_18 as u64],
        ctx,
    );
    fee_quoter::apply_token_transfer_fee_config_updates(
        &mut ref,
        &cap,
        100,
        vector[FEE_TOKEN],
        vector[50],
        vector[5000],
        vector[0],
        vector[60_000],
        vector[64],
        vector[true],
        vector[],
        ctx,
    );

    let fee_no_tokens = fee_quoter::get_validated_fee(
        &ref,
        &clock,
        100,
        SVM_RECEIVER,
        b"payload",
        vector[],
        vector[],
        FEE_TOKEN,
        client::encode_svm_extra_args_v1(200_000, 0, false, TOKEN_RECEIVER, vector[]),
    );
    let fee_with_token = fee_quoter::get_validated_fee(
        &ref,
        &clock,
        100,
        SVM_RECEIVER,
        b"payload",
        vector[FEE_TOKEN],
        vector[1000],
        FEE_TOKEN,
        client::encode_svm_extra_args_v1(200_000, 0, false, TOKEN_RECEIVER, vector[]),
    );

    assert!(fee_no_tokens > 0, 1);
    assert!(fee_with_token > fee_no_tokens, 2);

    fee_quoter::destroy_fee_quoter_cap(&ref, &cap, fq_cap);
    clock::destroy_for_testing(clock);
    cleanup(scenario, cap, ref);
}

#[test]
public fun test_svm_fee_max_accounts_succeeds() {
    let (mut scenario, cap, mut ref) = setup();
    let ctx = scenario.ctx();
    init_fee_quoter(&mut ref, &cap, ctx);
    let fq_cap = fee_quoter::new_fee_quoter_cap(&ref, &cap, ctx);
    let mut clock = clock::create_for_testing(ctx);
    clock::increment_for_testing(&mut clock, 20000);
    setup_prices(&mut ref, &fq_cap, &clock, ctx);
    setup_svm_dest(&mut ref, &cap, ctx);
    fee_quoter::apply_premium_multiplier_wei_per_eth_updates(
        &mut ref,
        &cap,
        vector[FEE_TOKEN],
        vector[ONE_E_18 as u64],
        ctx,
    );

    let mut accounts = vector[];
    let mut i = 0u64;
    while (i < 64u64) {
        let mut acc = vector[];
        let mut j = 0u64;
        while (j < 32u64) {
            acc.push_back((i as u8));
            j = j + 1u64;
        };
        accounts.push_back(acc);
        i = i + 1u64;
    };

    let fee = fee_quoter::get_validated_fee(
        &ref,
        &clock,
        100,
        SVM_RECEIVER,
        b"test_payload",
        vector[],
        vector[],
        FEE_TOKEN,
        client::encode_svm_extra_args_v1(
            200_000,
            0xffffffffffffffff,
            false,
            TOKEN_RECEIVER,
            accounts,
        ),
    );

    assert!(fee > 0, 1);

    fee_quoter::destroy_fee_quoter_cap(&ref, &cap, fq_cap);
    clock::destroy_for_testing(clock);
    cleanup(scenario, cap, ref);
}
