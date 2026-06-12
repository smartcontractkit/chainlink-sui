// Regression test for audit finding F-01:
//   "USDC + programmatic message inbound transfer permanently aborts (stranded burned funds)"
//
// Before the fix, offramp_state_helper::complete_token_transfer re-derived the destination
// amount by parsing the token's source_pool_data as an ABI-encoded decimals word
// (assert!(data_len == 32, EInvalidRemoteChainDecimals)). The USDC/CCTP pool encodes its
// source_pool_data as encode_u64(nonce) ++ encode_u32(domain) = 64 bytes, so every USDC
// transfer that also carried a programmatic message aborted on that assert — stranding the
// already-CCTP-burned source funds.
//
// The fix makes the pool the single source of truth: it passes the local_amount it actually
// minted/released into complete_token_transfer, which surfaces that value to the receiver
// verbatim and never re-parses source_pool_data. This test proves the 64-byte USDC + message
// path now COMPLETES and surfaces the pool's amount.
//
// It uses only the audited public API and the genuine owner-minted DestTransferCap (the same
// object the offramp consumes at initialize); nothing is forged. The 64-byte blob is exactly
// what usdc_token_pool::encode_source_pool_data produces.

#[test_only]
module ccip::f01_usdc_message_regression_test;

use ccip::client;
use ccip::eth_abi;
use ccip::offramp_state_helper::{Self, DestTransferCap};
use ccip::ownable::OwnerCap;
use ccip::state_object::{Self, CCIPObjectRef};
use ccip::token_admin_registry as registry;
use ccip::upgrade_registry;
use std::ascii;
use std::string;
use std::type_name;
use sui::test_scenario::{Self as ts, Scenario};

public struct TestTypeProof has drop {}

const OWNER: address = @0x1000;
const RECEIVER_ADDRESS: address = @0x2000;
const TOKEN_ADDRESS: address =
    @0x1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b;
const TOKEN_POOL_ADDRESS: address =
    @0xdeeb7a4662eec9f2f3def03fb937a663dddaa2e215b8078a284d026b7946c270;
const SOURCE_CHAIN_SELECTOR: u64 = 1000;
const SOURCE_AMOUNT: u256 = 1_000_000; // 1 USDC at 6 decimals (source-chain denominated)
// The amount the USDC pool actually mints on Sui (CCTP-attested), passed into
// complete_token_transfer as the single source of truth.
const LOCAL_MINTED_AMOUNT: u64 = 1_000_000;
const LOCAL_DECIMALS: u8 = 6; // USDC is 6 decimals on Sui
const USDC_DOMAIN: u32 = 0; // Ethereum CCTP domain
const USDC_NONCE: u64 = 12345;

fun setup_test(): (Scenario, OwnerCap, CCIPObjectRef, DestTransferCap) {
    let mut scenario = ts::begin(OWNER);
    let ctx = scenario.ctx();

    state_object::test_init(ctx);
    scenario.next_tx(OWNER);

    let owner_cap = scenario.take_from_sender<OwnerCap>();
    let mut ref = scenario.take_shared<CCIPObjectRef>();

    upgrade_registry::initialize(&mut ref, &owner_cap, scenario.ctx());
    registry::initialize(&mut ref, &owner_cap, scenario.ctx());
    registry::initialize_local_decimals(&mut ref, &owner_cap, scenario.ctx());

    offramp_state_helper::test_init(scenario.ctx());
    scenario.next_tx(OWNER);
    let dest_cap = scenario.take_from_sender<DestTransferCap>();

    (scenario, owner_cap, ref, dest_cap)
}

fun cleanup_test(
    scenario: Scenario,
    owner_cap: OwnerCap,
    ref: CCIPObjectRef,
    dest_cap: DestTransferCap,
) {
    ts::return_to_sender(&scenario, owner_cap);
    ts::return_shared(ref);
    transfer::public_transfer(dest_cap, @0x0);
    ts::end(scenario);
}

fun register_usdc_pool(owner_cap: &OwnerCap, ref: &mut CCIPObjectRef, scenario: &mut Scenario) {
    registry::register_pool_as_owner_v2(
        owner_cap,
        ref,
        TOKEN_ADDRESS,
        TOKEN_POOL_ADDRESS,
        string::utf8(b"usdc_token_pool"),
        ascii::string(b"TestType"),
        OWNER,
        type_name::into_string(type_name::with_defining_ids<TestTypeProof>()),
        vector<address>[],
        vector<address>[],
        LOCAL_DECIMALS,
        scenario.ctx(),
    );
}

/// Byte-exact mirror of usdc_token_pool::encode_source_pool_data:
///   nonce (u64) then domain (u32), each ABI-encoded into a 32-byte big-endian word => 64 bytes.
fun usdc_source_pool_data(domain: u32, nonce: u64): vector<u8> {
    let mut spd = vector[];
    eth_abi::encode_u64(&mut spd, nonce); // 32 bytes
    eth_abi::encode_u32(&mut spd, domain); // 32 bytes
    spd
}

/// Precondition: the bytes the USDC pool emits are 64 long (not 0 and not 32). This is the
/// blob that the pre-fix `assert!(data_len == 32)` rejected.
#[test]
fun test_usdc_source_pool_data_is_64_bytes() {
    let spd = usdc_source_pool_data(USDC_DOMAIN, USDC_NONCE);
    assert!(spd.length() == 64, 9999);
}

/// THE F-01 REGRESSION: a USDC transfer carrying a programmatic message and a 64-byte
/// source_pool_data now COMPLETES (no abort) and surfaces the pool's local minted amount to
/// the receiver. This is the 64-byte case the prior test suite never covered.
#[test]
fun test_usdc_plus_message_completes_and_surfaces_local_amount() {
    let (mut scenario, owner_cap, mut ref, dest_cap) = setup_test();
    register_usdc_pool(&owner_cap, &mut ref, &mut scenario);

    let mut receiver_params = offramp_state_helper::create_receiver_params(
        &dest_cap,
        SOURCE_CHAIN_SELECTOR,
    );

    // The 64-byte USDC blob (domain ++ nonce) — exactly what the offramp copies into
    // receiver_params from the source pool. With the fix it is no longer parsed as decimals.
    offramp_state_helper::add_dest_token_transfer(
        &dest_cap,
        &mut receiver_params,
        RECEIVER_ADDRESS,
        SOURCE_CHAIN_SELECTOR,
        SOURCE_AMOUNT,
        TOKEN_ADDRESS,
        TOKEN_POOL_ADDRESS,
        b"source_pool_address",
        usdc_source_pool_data(USDC_DOMAIN, USDC_NONCE),
        b"offchain_data",
    );

    // The programmatic-message leg: a populated message makes message.is_some(), the branch
    // that pre-fix triggered the decimals re-parse and aborted.
    let dest_token_amounts = client::new_dest_token_amounts(
        vector[TOKEN_ADDRESS],
        vector[SOURCE_AMOUNT],
    );
    let test_message = client::new_any2sui_message(
        b"message_id_32_bytes_long_test_msg",
        SOURCE_CHAIN_SELECTOR,
        b"sender_address",
        b"call_my_receiver",
        @0x5432,
        @0x12345,
        dest_token_amounts,
    );
    offramp_state_helper::populate_message(&dest_cap, &mut receiver_params, test_message);

    // Pre-fix this aborted with EInvalidRemoteChainDecimals; post-fix it completes and writes
    // the pool's already-minted local amount into the message.
    offramp_state_helper::complete_token_transfer(
        &ref,
        &mut receiver_params,
        LOCAL_MINTED_AMOUNT,
        TestTypeProof {},
    );

    let message = offramp_state_helper::extract_any2sui_message(&mut receiver_params);
    let (_, _, _, _, _, _, extracted) = client::consume_any2sui_message(message, @0x5432);
    // The receiver sees exactly the amount the pool minted on Sui — the single source of truth.
    assert!(client::get_amount(&extracted[0]) == (LOCAL_MINTED_AMOUNT as u256));

    offramp_state_helper::deconstruct_receiver_params(&dest_cap, receiver_params);
    cleanup_test(scenario, owner_cap, ref, dest_cap);
}

/// Control: the same 64-byte USDC transfer with NO message also completes (it always did —
/// the pure token path never re-parsed source_pool_data).
#[test]
fun test_usdc_token_only_succeeds() {
    let (mut scenario, owner_cap, mut ref, dest_cap) = setup_test();
    register_usdc_pool(&owner_cap, &mut ref, &mut scenario);

    let mut receiver_params = offramp_state_helper::create_receiver_params(
        &dest_cap,
        SOURCE_CHAIN_SELECTOR,
    );

    offramp_state_helper::add_dest_token_transfer(
        &dest_cap,
        &mut receiver_params,
        RECEIVER_ADDRESS,
        SOURCE_CHAIN_SELECTOR,
        SOURCE_AMOUNT,
        TOKEN_ADDRESS,
        TOKEN_POOL_ADDRESS,
        b"source_pool_address",
        usdc_source_pool_data(USDC_DOMAIN, USDC_NONCE),
        b"offchain_data",
    );

    offramp_state_helper::complete_token_transfer(
        &ref,
        &mut receiver_params,
        LOCAL_MINTED_AMOUNT,
        TestTypeProof {},
    );

    offramp_state_helper::deconstruct_receiver_params(&dest_cap, receiver_params);
    cleanup_test(scenario, owner_cap, ref, dest_cap);
}
