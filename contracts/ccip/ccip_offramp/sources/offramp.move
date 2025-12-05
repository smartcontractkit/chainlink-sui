/// The OffRamp package handles merkle root commitments and message execution.
/// Future versions of this contract will be deployed as a separate package to avoid any unwanted side effects
/// during upgrades.
module ccip_offramp::offramp;

use ccip::eth_abi;
use ccip::fee_quoter::{Self, FeeQuoterCap};
use ccip::merkle_proof;
use ccip::offramp_state_helper as osh;
use ccip::receiver_registry;
use ccip::rmn_remote;
use ccip::state_object::CCIPObjectRef;
use ccip::token_admin_registry;
use ccip::upgrade_registry::verify_function_allowed;
use ccip_offramp::ocr3_base::{Self, OCR3BaseState, OCRConfig};
use ccip_offramp::ownable::{Self, OwnerCap, OwnableState};
use mcms::bcs_stream::{Self, BCSStream};
use mcms::mcms_deployer::{Self, DeployerState};
use mcms::mcms_registry::{Self, Registry, ExecutingCallbackParams};
use std::ascii;
use std::string::{Self, String};
use std::type_name;
use sui::address;
use sui::clock;
use sui::derived_object;
use sui::event;
use sui::hash;
use sui::package::{Self, UpgradeCap};
use sui::table::{Self, Table};
use sui::vec_map::{Self, VecMap};

public struct OffRampState has key, store {
    id: UID,
    package_ids: vector<address>,
    ocr3_base_state: OCR3BaseState,
    // static config
    chain_selector: u64,
    // dynamic config
    permissionless_execution_threshold_seconds: u32,
    // source chain selector -> config
    source_chain_configs: VecMap<u64, SourceChainConfig>,
    // source chain selector -> seq num -> execution state
    execution_states: Table<u64, Table<u64, u8>>,
    // merkle root -> timestamp in secs
    roots: Table<vector<u8>, u64>,
    // This is the OCR sequence number, not to be confused with the CCIP message sequence number.
    latest_price_sequence_number: u64,
    // provided when initializing the fee quoter in CCIP package
    fee_quoter_cap: Option<FeeQuoterCap>,
    dest_transfer_cap: Option<osh::DestTransferCap>,
    ownable_state: OwnableState,
}

public struct OffRampObject has key {
    id: UID,
}

public struct OffRampStatePointer has key, store {
    id: UID,
    off_ramp_object_id: address,
}

public struct SourceChainConfig has copy, drop, store {
    router: address,
    is_enabled: bool,
    min_seq_nr: u64,
    is_rmn_verification_disabled: bool,
    on_ramp: vector<u8>,
}

// report public structs
public struct RampMessageHeader has drop {
    message_id: vector<u8>,
    source_chain_selector: u64,
    dest_chain_selector: u64,
    sequence_number: u64,
    nonce: u64,
}

public struct Any2SuiRampMessage has drop {
    header: RampMessageHeader,
    sender: vector<u8>,
    data: vector<u8>,
    receiver: address, // this is the message receiver
    gas_limit: u256,
    token_receiver: address,
    token_amounts: vector<Any2SuiTokenTransfer>,
}

public struct Any2SuiTokenTransfer has drop {
    source_pool_address: vector<u8>,
    // the token's coin metadata object id on SUI
    dest_token_address: address,
    dest_gas_amount: u32,
    extra_data: vector<u8>,
    amount: u256, // This is the amount to transfer, as set on the source chain.
}

public struct ExecutionReport has drop {
    source_chain_selector: u64,
    message: Any2SuiRampMessage,
    offchain_token_data: vector<vector<u8>>,
    proofs: vector<vector<u8>>, // Proofs used to construct the merkle root
}

// Matches the EVM public struct
public struct CommitReport has copy, drop, store {
    price_updates: PriceUpdates, // Price updates for the fee_quoter
    blessed_merkle_roots: vector<MerkleRoot>, // Merkle roots that have been blessed by RMN
    unblessed_merkle_roots: vector<MerkleRoot>, // Merkle roots that don't require RMN blessing
    rmn_signatures: vector<vector<u8>>, // The signatures for the blessed merkle roots
}

public struct PriceUpdates has copy, drop, store {
    token_price_updates: vector<TokenPriceUpdate>,
    gas_price_updates: vector<GasPriceUpdate>,
}

public struct TokenPriceUpdate has copy, drop, store {
    source_token: address,
    usd_per_token: u256,
}

public struct GasPriceUpdate has copy, drop, store {
    dest_chain_selector: u64,
    usd_per_unit_gas: u256,
}

public struct MerkleRoot has copy, drop, store {
    source_chain_selector: u64,
    on_ramp_address: vector<u8>,
    min_seq_nr: u64,
    max_seq_nr: u64,
    merkle_root: vector<u8>,
}

public struct StaticConfig has copy, drop, store {
    chain_selector: u64,
    rmn_remote: address,
    token_admin_registry: address,
    nonce_manager: address,
}

// On EVM, the feeQuoter is a dynamic address but due to the Sui implementation using a static
// upgradable FeeQuoter stored within the state ref, this value is actually static and cannot be
// accessed by its object id/address directly by users.
// For compatibility reasons, we keep it as a dynamic config.
public struct DynamicConfig has copy, drop, store {
    fee_quoter: address,
    permissionless_execution_threshold_seconds: u32, // The delay before manual exec is enabled
}

public struct StaticConfigSet has copy, drop {
    chain_selector: u64,
}

public struct DynamicConfigSet has copy, drop {
    dynamic_config: DynamicConfig,
}

public struct SourceChainConfigSet has copy, drop {
    source_chain_selector: u64,
    source_chain_config: SourceChainConfig,
}

public struct SkippedAlreadyExecuted has copy, drop {
    source_chain_selector: u64,
    sequence_number: u64,
}

public struct ExecutionStateChanged has copy, drop {
    source_chain_selector: u64,
    sequence_number: u64,
    message_id: vector<u8>,
    message_hash: vector<u8>,
    state: u8,
}

public struct CommitReportAccepted has copy, drop {
    blessed_merkle_roots: vector<MerkleRoot>,
    unblessed_merkle_roots: vector<MerkleRoot>,
    price_updates: PriceUpdates,
}

public struct SkippedReportExecution has copy, drop {
    source_chain_selector: u64,
}

const TOKEN_TRANSFER_LIMIT: u64 = 1;

/// These have to match the EVM states
/// However, execution in SUI is done in a single PTB,
/// so we don't have the IN_PROGRESS or FAILURE states.
const EXECUTION_STATE_UNTOUCHED: u8 = 0;
// const EXECUTION_STATE_IN_PROGRESS: u8 = 1;
const EXECUTION_STATE_SUCCESS: u8 = 2;
// const EXECUTION_STATE_FAILURE: u8 = 3;
const ZERO_MERKLE_ROOT: vector<u8> = vector[
    0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
];

const ESourceChainSelectorsMismatch: u64 = 1;
const EZeroChainSelector: u64 = 2;
const EUnknownSourceChainSelector: u64 = 3;
const EMustBeOutOfOrderExec: u64 = 4;
const ESourceChainSelectorMismatch: u64 = 5;
const EDestChainSelectorMismatch: u64 = 6;
const ETokenDataMismatch: u64 = 7;
const ERootNotCommitted: u64 = 8;
const EManualExecutionNotYetEnabled: u64 = 9;
const ESourceChainNotEnabled: u64 = 10;
const ECommitOnRampMismatch: u64 = 11;
const EInvalidInterval: u64 = 12;
const EInvalidRoot: u64 = 13;
const ERootAlreadyCommitted: u64 = 14;
const EStaleCommitReport: u64 = 15;
const ECursedByRmn: u64 = 16;
const ESignatureVerificationRequiredInCommitPlugin: u64 = 17;
const ESignatureVerificationNotAllowedInExecutionPlugin: u64 = 18;
const EFeeQuoterCapExists: u64 = 19;
const EDestTransferCapExists: u64 = 20;
const ERmnBlessingMismatch: u64 = 21;
const EUnsupportedToken: u64 = 22;
const EInvalidOnRampUpdate: u64 = 23;
const EDestTransferCapNotSet: u64 = 24;
const ECalculateMessageHashInvalidArguments: u64 = 25;
const EInvalidFunction: u64 = 26;
const EInvalidTokenReceiver: u64 = 27;
const ETokenTransferLimitExceeded: u64 = 28;
const EPackageIdNotFound: u64 = 29;
const EInvalidOwnerCap: u64 = 30;
const EUnknownSequenceNumber: u64 = 31;
const EInvalidReportContextLength: u64 = 32;

const VERSION: u8 = 1;

public fun type_and_version(): String {
    string::utf8(b"OffRamp 1.6.0")
}

public struct OFFRAMP has drop {}

fun init(otw: OFFRAMP, ctx: &mut TxContext) {
    let mut off_ramp_object = OffRampObject { id: object::new(ctx) };
    let (ownable_state, mut owner_cap) = ownable::new(&mut off_ramp_object.id, ctx);

    let state = OffRampState {
        id: derived_object::claim(&mut off_ramp_object.id, b"OffRampState"),
        package_ids: vector[],
        ocr3_base_state: ocr3_base::new(ctx),
        chain_selector: 0,
        permissionless_execution_threshold_seconds: 0,
        source_chain_configs: vec_map::empty<u64, SourceChainConfig>(),
        execution_states: table::new(ctx),
        roots: table::new(ctx),
        latest_price_sequence_number: 0,
        fee_quoter_cap: option::none(),
        dest_transfer_cap: option::none(),
        ownable_state,
    };

    let pointer = OffRampStatePointer {
        id: object::new(ctx),
        off_ramp_object_id: object::id_address(&off_ramp_object),
    };

    let tn = type_name::with_original_ids<OFFRAMP>();
    let package_bytes = ascii::into_bytes(tn.address_string());
    let package_id = address::from_ascii_bytes(&package_bytes);

    transfer::share_object(state);
    transfer::share_object(off_ramp_object);

    let publisher = package::claim(otw, ctx);
    ownable::attach_publisher(&mut owner_cap, publisher);

    transfer::public_transfer(owner_cap, ctx.sender());
    transfer::transfer(pointer, package_id);
}

public fun initialize(
    state: &mut OffRampState,
    owner_cap: &OwnerCap,
    fee_quoter_cap: FeeQuoterCap,
    dest_transfer_cap: osh::DestTransferCap,
    chain_selector: u64,
    permissionless_execution_threshold_seconds: u32,
    source_chains_selectors: vector<u64>,
    source_chains_is_enabled: vector<bool>,
    source_chains_is_rmn_verification_disabled: vector<bool>,
    source_chains_on_ramp: vector<vector<u8>>,
    ctx: &mut TxContext,
) {
    assert!(object::id(owner_cap) == ownable::owner_cap_id(&state.ownable_state), EInvalidOwnerCap);
    assert!(chain_selector != 0, EZeroChainSelector);
    state.chain_selector = chain_selector;

    assert!(state.fee_quoter_cap.is_none(), EFeeQuoterCapExists);
    state.fee_quoter_cap.fill(fee_quoter_cap);
    assert!(state.dest_transfer_cap.is_none(), EDestTransferCapExists);
    state.dest_transfer_cap.fill(dest_transfer_cap);

    event::emit(StaticConfigSet { chain_selector });

    set_dynamic_config_internal(
        state,
        permissionless_execution_threshold_seconds,
    );
    apply_source_chain_config_updates_internal(
        state,
        source_chains_selectors,
        source_chains_is_enabled,
        source_chains_is_rmn_verification_disabled,
        source_chains_on_ramp,
        ctx,
    );

    let tn = type_name::with_original_ids<OFFRAMP>();
    let package_bytes = ascii::into_bytes(tn.address_string());
    let package_id = address::from_ascii_bytes(&package_bytes);
    state.package_ids.push_back(package_id);
}

public fun add_package_id(state: &mut OffRampState, owner_cap: &OwnerCap, package_id: address) {
    assert!(object::id(owner_cap) == ownable::owner_cap_id(&state.ownable_state), EInvalidOwnerCap);
    state.package_ids.push_back(package_id);
}

public fun remove_package_id(state: &mut OffRampState, owner_cap: &OwnerCap, package_id: address) {
    assert!(object::id(owner_cap) == ownable::owner_cap_id(&state.ownable_state), EInvalidOwnerCap);
    let (found, idx) = state.package_ids.index_of(&package_id);
    assert!(found, EPackageIdNotFound);
    state.package_ids.remove(idx);
}

public fun get_ocr3_base(state: &OffRampState): &OCR3BaseState {
    &state.ocr3_base_state
}

fun set_dynamic_config_internal(
    state: &mut OffRampState,
    permissionless_execution_threshold_seconds: u32,
) {
    state.permissionless_execution_threshold_seconds = permissionless_execution_threshold_seconds;
    let dynamic_config = create_dynamic_config(permissionless_execution_threshold_seconds);
    event::emit(DynamicConfigSet { dynamic_config });
}

fun create_dynamic_config(permissionless_execution_threshold_seconds: u32): DynamicConfig {
    DynamicConfig { fee_quoter: @ccip, permissionless_execution_threshold_seconds }
}

fun apply_source_chain_config_updates_internal(
    state: &mut OffRampState,
    source_chains_selector: vector<u64>,
    source_chains_is_enabled: vector<bool>,
    source_chains_is_rmn_verification_disabled: vector<bool>,
    source_chains_on_ramp: vector<vector<u8>>,
    ctx: &mut TxContext,
) {
    let source_chains_len = source_chains_selector.length();
    assert!(source_chains_len == source_chains_is_enabled.length(), ESourceChainSelectorsMismatch);
    assert!(
        source_chains_len == source_chains_is_rmn_verification_disabled.length(),
        ESourceChainSelectorsMismatch,
    );
    assert!(source_chains_len == source_chains_on_ramp.length(), ESourceChainSelectorsMismatch);

    let mut i = 0;
    while (i < source_chains_len) {
        let source_chain_selector = source_chains_selector[i];
        let is_enabled = source_chains_is_enabled[i];
        let is_rmn_verification_disabled = source_chains_is_rmn_verification_disabled[i];
        let on_ramp = source_chains_on_ramp[i];

        assert!(source_chain_selector != 0, EZeroChainSelector);
        ccip::address::assert_non_zero_address_vector(&on_ramp);

        if (state.source_chain_configs.contains(&source_chain_selector)) {
            // OnRamp updates should only happen due to a misconfiguration.
            // If an OnRamp is misconfigured, no reports should have been
            // committed and no messages should have been executed.
            let existing_config = state.source_chain_configs.get(&source_chain_selector);
            assert!(
                existing_config.min_seq_nr == 1 || existing_config.on_ramp == on_ramp,
                EInvalidOnRampUpdate,
            );
        } else {
            state
                .source_chain_configs
                .insert(
                    source_chain_selector,
                    SourceChainConfig {
                        router: @ccip,
                        is_enabled: false,
                        min_seq_nr: 1,
                        is_rmn_verification_disabled: false,
                        on_ramp: vector[],
                    },
                );
            state.execution_states.add(source_chain_selector, table::new(ctx));
        };

        let config = state.source_chain_configs.get_mut(&source_chain_selector);
        config.is_enabled = is_enabled;
        config.on_ramp = on_ramp;
        config.is_rmn_verification_disabled = is_rmn_verification_disabled;

        event::emit(SourceChainConfigSet { source_chain_selector, source_chain_config: *config });
        i = i + 1;
    }
}

fun assert_source_chain_enabled(state: &OffRampState, source_chain_selector: u64) {
    // assert that the source chain is enabled.
    assert!(
        state.source_chain_configs.contains(&source_chain_selector),
        EUnknownSourceChainSelector,
    );
    let source_chain_config = state.source_chain_configs.get(&source_chain_selector);
    assert!(source_chain_config.is_enabled, ESourceChainNotEnabled);
}

// ================================================================
// |                          Execution                           |
// ================================================================

public fun init_execute(
    ref: &CCIPObjectRef,
    state: &mut OffRampState,
    clock: &clock::Clock,
    report_context: vector<vector<u8>>,
    report: vector<u8>,
    ctx: &mut TxContext,
): osh::ReceiverParams {
    verify_function_allowed(
        ref,
        string::utf8(b"offramp"),
        string::utf8(b"init_execute"),
        VERSION,
    );
    assert!(report_context.length() == 2, EInvalidReportContextLength);
    let reports = deserialize_execution_report(report);

    ocr3_base::transmit(
        &state.ocr3_base_state,
        ctx.sender(),
        ocr3_base::ocr_plugin_type_execution(),
        report_context,
        report,
        vector[],
        ctx,
    );

    pre_execute_single_report(ref, state, clock, reports, false)
}

public fun finish_execute(
    ref: &CCIPObjectRef,
    state: &mut OffRampState,
    receiver_params: osh::ReceiverParams,
) {
    verify_function_allowed(
        ref,
        string::utf8(b"offramp"),
        string::utf8(b"finish_execute"),
        VERSION,
    );
    assert!(state.dest_transfer_cap.is_some(), EDestTransferCapNotSet);
    osh::deconstruct_receiver_params(state.dest_transfer_cap.borrow(), receiver_params);
}

// this function does not involve ocr3 transmit & it sets manual_execution to true
public fun manually_init_execute(
    ref: &CCIPObjectRef,
    state: &mut OffRampState,
    clock: &clock::Clock,
    report_bytes: vector<u8>,
): osh::ReceiverParams {
    verify_function_allowed(
        ref,
        string::utf8(b"offramp"),
        string::utf8(b"manually_init_execute"),
        VERSION,
    );
    let reports = deserialize_execution_report(report_bytes);

    pre_execute_single_report(ref, state, clock, reports, true)
}

public fun get_execution_state(
    state: &OffRampState,
    source_chain_selector: u64,
    sequence_number: u64,
): u8 {
    assert!(state.execution_states.contains(source_chain_selector), EUnknownSourceChainSelector);
    let source_chain_execution_states = state.execution_states.borrow(source_chain_selector);
    assert!(source_chain_execution_states.contains(sequence_number), EUnknownSequenceNumber);
    let execution_state = source_chain_execution_states.borrow(sequence_number);
    *execution_state
}

fun deserialize_execution_report(report_bytes: vector<u8>): ExecutionReport {
    let mut stream = bcs_stream::new(report_bytes);
    let source_chain_selector = bcs_stream::deserialize_u64(&mut stream);

    let message_id = bcs_stream::deserialize_fixed_vector_u8(&mut stream, 32);
    let header_source_chain_selector = bcs_stream::deserialize_u64(&mut stream);
    let dest_chain_selector = bcs_stream::deserialize_u64(&mut stream);
    let sequence_number = bcs_stream::deserialize_u64(&mut stream);
    let nonce = bcs_stream::deserialize_u64(&mut stream);

    let header = RampMessageHeader {
        message_id,
        source_chain_selector: header_source_chain_selector,
        dest_chain_selector,
        sequence_number,
        nonce,
    };

    assert!(source_chain_selector == header_source_chain_selector, ESourceChainSelectorMismatch);

    let sender = bcs_stream::deserialize_vector_u8(&mut stream);
    let data = bcs_stream::deserialize_vector_u8(&mut stream);
    let receiver = bcs_stream::deserialize_address(&mut stream);
    let gas_limit = bcs_stream::deserialize_u256(&mut stream);
    let token_receiver = bcs_stream::deserialize_address(&mut stream);

    let token_amounts = bcs_stream::deserialize_vector!(&mut stream, |stream| {
        let source_pool_address = bcs_stream::deserialize_vector_u8(stream);
        let dest_token_address = bcs_stream::deserialize_address(stream);
        let dest_gas_amount = bcs_stream::deserialize_u32(stream);
        let extra_data = bcs_stream::deserialize_vector_u8(stream);
        let amount = bcs_stream::deserialize_u256(stream);

        Any2SuiTokenTransfer {
            source_pool_address,
            dest_token_address,
            dest_gas_amount,
            extra_data,
            amount,
        }
    });

    let message = Any2SuiRampMessage {
        header,
        sender,
        data,
        receiver,
        gas_limit,
        token_receiver,
        token_amounts,
    };

    let offchain_token_data = bcs_stream::deserialize_vector!(
        &mut stream,
        |stream| bcs_stream::deserialize_vector_u8(stream),
    );

    let proofs = bcs_stream::deserialize_vector!(
        &mut stream,
        |stream| { bcs_stream::deserialize_fixed_vector_u8(stream, 32) },
    );

    bcs_stream::assert_is_consumed(&stream);

    ExecutionReport { source_chain_selector, message, offchain_token_data, proofs }
}

#[allow(implicit_const_copy)]
fun pre_execute_single_report(
    ref: &CCIPObjectRef,
    state: &mut OffRampState,
    clock: &clock::Clock,
    execution_report: ExecutionReport,
    manual_execution: bool,
): osh::ReceiverParams {
    let source_chain_selector = execution_report.source_chain_selector;
    assert!(state.dest_transfer_cap.is_some(), EDestTransferCapNotSet);

    if (rmn_remote::is_cursed_u128(ref, source_chain_selector as u128)) {
        assert!(!manual_execution, ECursedByRmn);

        event::emit(SkippedReportExecution { source_chain_selector });

        return osh::create_receiver_params(state.dest_transfer_cap.borrow(), source_chain_selector)
    };

    assert_source_chain_enabled(state, source_chain_selector);

    assert!(
        execution_report.message.header.dest_chain_selector == state.chain_selector,
        EDestChainSelectorMismatch,
    );

    let source_chain_config = state.source_chain_configs[&source_chain_selector];
    let metadata_hash = calculate_metadata_hash(
        ref,
        source_chain_selector,
        state.chain_selector,
        source_chain_config.on_ramp,
    );

    let hashed_leaf = calculate_message_hash_internal(
        &execution_report.message,
        metadata_hash,
    );

    let root = merkle_proof::merkle_root(hashed_leaf, execution_report.proofs);

    // Essential security check
    let is_old_commit_report = is_committed_root(state, clock, root);

    if (manual_execution) {
        assert!(is_old_commit_report, EManualExecutionNotYetEnabled);
    };

    let source_chain_execution_states = state.execution_states.borrow_mut(source_chain_selector);

    let message = &execution_report.message;
    let sequence_number = message.header.sequence_number;
    if (!source_chain_execution_states.contains(sequence_number)) {
        source_chain_execution_states.add(sequence_number, EXECUTION_STATE_UNTOUCHED);
    };
    let execution_state_ref = source_chain_execution_states.borrow_mut(sequence_number);

    if (*execution_state_ref != EXECUTION_STATE_UNTOUCHED) {
        event::emit(SkippedAlreadyExecuted { source_chain_selector, sequence_number });

        return osh::create_receiver_params(state.dest_transfer_cap.borrow(), source_chain_selector)
    };

    // A zero nonce indicates out of order execution which is the only allowed case.
    assert!(message.header.nonce == 0, EMustBeOutOfOrderExec);

    let number_of_tokens_in_msg = message.token_amounts.length();
    assert!(number_of_tokens_in_msg <= TOKEN_TRANSFER_LIMIT, ETokenTransferLimitExceeded);
    assert!(
        number_of_tokens_in_msg == execution_report.offchain_token_data.length(),
        ETokenDataMismatch,
    );
    assert!(
        message.token_receiver == @0x0 && number_of_tokens_in_msg == 0 || // if token_receiver is empty, no tokens should be transferred
            (message.token_receiver != @0x0 && number_of_tokens_in_msg > 0), // if token_receiver is not empty, tokens should be transferred
        EInvalidTokenReceiver,
    );

    let mut receiver_params = osh::create_receiver_params(
        state.dest_transfer_cap.borrow(),
        source_chain_selector,
    );

    let mut token_addresses = vector[];
    let mut token_amounts = vector[];

    if (number_of_tokens_in_msg == TOKEN_TRANSFER_LIMIT) {
        let token_pool_address: address = token_admin_registry::get_pool(
            ref,
            message.token_amounts[0].dest_token_address,
        );
        assert!(token_pool_address != @0x0, EUnsupportedToken);

        osh::add_dest_token_transfer(
            state.dest_transfer_cap.borrow(),
            &mut receiver_params,
            message.token_receiver, // when sending tokens, token receiver will be included in the execution report
            source_chain_selector,
            message.token_amounts[0].amount,
            message.token_amounts[0].dest_token_address,
            token_pool_address,
            message.token_amounts[0].source_pool_address,
            message.token_amounts[0].extra_data,
            execution_report.offchain_token_data[0],
        );
        token_addresses.push_back(message.token_amounts[0].dest_token_address);
        token_amounts.push_back(message.token_amounts[0].amount);
    };

    let has_valid_message_receiver =
        (!message.data.is_empty() || message.gas_limit != 0) && receiver_registry::is_registered_receiver(ref, message.receiver);
    // if the message has a valid message receiver and proper data & gas limit
    if (has_valid_message_receiver) {
        let any2sui_message = osh::new_any2sui_message(
            state.dest_transfer_cap.borrow(),
            message.header.message_id,
            message.header.source_chain_selector,
            message.sender,
            message.data,
            message.receiver,
            message.token_receiver,
            token_addresses,
            token_amounts,
        );

        osh::populate_message(
            state.dest_transfer_cap.borrow(),
            &mut receiver_params,
            any2sui_message,
        );
    };

    // the entire PTB either succeeds or fails so we can set the state to success
    *execution_state_ref = EXECUTION_STATE_SUCCESS;

    event::emit(ExecutionStateChanged {
        source_chain_selector,
        sequence_number,
        message_id: message.header.message_id,
        message_hash: hashed_leaf,
        state: EXECUTION_STATE_SUCCESS,
    });

    // return the hot potato to user/execution DON
    receiver_params
}

/// Throws an error if the root is not committed.
/// Returns true if the root is eligable for manual execution.
fun is_committed_root(state: &OffRampState, clock: &clock::Clock, root: vector<u8>): bool {
    assert!(state.roots.contains(root), ERootNotCommitted);
    let timestamp_committed_secs = state.roots[root];

    (clock.timestamp_ms() / 1000 - timestamp_committed_secs)
            > (state.permissionless_execution_threshold_seconds as u64)
}

// ================================================================
// |                        Metadata hash                         |
// ================================================================

public fun calculate_metadata_hash(
    ref: &CCIPObjectRef,
    source_chain_selector: u64,
    dest_chain_selector: u64,
    on_ramp: vector<u8>,
): vector<u8> {
    verify_function_allowed(
        ref,
        string::utf8(b"offramp"),
        string::utf8(b"calculate_metadata_hash"),
        VERSION,
    );
    let mut packed = vector[];
    eth_abi::encode_right_padded_bytes32(
        &mut packed,
        hash::keccak256(&b"Any2SuiMessageHashV1"),
    );
    eth_abi::encode_u64(&mut packed, source_chain_selector);
    eth_abi::encode_u64(&mut packed, dest_chain_selector);
    eth_abi::encode_right_padded_bytes32(&mut packed, hash::keccak256(&on_ramp));
    hash::keccak256(&packed)
}

public fun calculate_message_hash(
    ref: &CCIPObjectRef,
    message_id: vector<u8>,
    source_chain_selector: u64,
    dest_chain_selector: u64,
    sequence_number: u64,
    nonce: u64,
    sender: vector<u8>,
    receiver: address,
    on_ramp: vector<u8>,
    data: vector<u8>,
    gas_limit: u256,
    token_receiver: address,
    source_pool_addresses: vector<vector<u8>>,
    dest_token_addresses: vector<address>,
    dest_gas_amounts: vector<u32>,
    extra_datas: vector<vector<u8>>,
    amounts: vector<u256>,
): vector<u8> {
    verify_function_allowed(
        ref,
        string::utf8(b"offramp"),
        string::utf8(b"calculate_message_hash"),
        VERSION,
    );
    let source_pool_addresses_len = source_pool_addresses.length();
    assert!(
        source_pool_addresses_len == dest_token_addresses.length()
                && source_pool_addresses_len == dest_gas_amounts.length()
                && source_pool_addresses_len == extra_datas.length()
                && source_pool_addresses_len == amounts.length(),
        ECalculateMessageHashInvalidArguments,
    );

    let metadata_hash = calculate_metadata_hash(
        ref,
        source_chain_selector,
        dest_chain_selector,
        on_ramp,
    );

    let mut token_amounts = vector[];
    let mut i = 0;
    while (i < source_pool_addresses_len) {
        token_amounts.push_back(Any2SuiTokenTransfer {
            source_pool_address: source_pool_addresses[i],
            dest_token_address: dest_token_addresses[i],
            dest_gas_amount: dest_gas_amounts[i],
            extra_data: extra_datas[i],
            amount: amounts[i],
        });
        i = i + 1;
    };

    let message = Any2SuiRampMessage {
        header: RampMessageHeader {
            message_id,
            source_chain_selector,
            dest_chain_selector,
            sequence_number,
            nonce,
        },
        sender,
        data,
        receiver,
        gas_limit,
        token_receiver,
        token_amounts,
    };

    calculate_message_hash_internal(&message, metadata_hash)
}

fun calculate_message_hash_internal(
    message: &Any2SuiRampMessage,
    metadata_hash: vector<u8>,
): vector<u8> {
    let mut outer_hash = vector[];
    eth_abi::encode_right_padded_bytes32(&mut outer_hash, merkle_proof::leaf_domain_separator());
    eth_abi::encode_right_padded_bytes32(&mut outer_hash, metadata_hash);

    let mut inner_hash = vector[];
    eth_abi::encode_right_padded_bytes32(&mut inner_hash, message.header.message_id);
    eth_abi::encode_address(&mut inner_hash, message.receiver);
    eth_abi::encode_u64(&mut inner_hash, message.header.sequence_number);
    eth_abi::encode_u256(&mut inner_hash, message.gas_limit);
    eth_abi::encode_address(&mut inner_hash, message.token_receiver);
    eth_abi::encode_u64(&mut inner_hash, message.header.nonce);
    eth_abi::encode_right_padded_bytes32(&mut outer_hash, hash::keccak256(&inner_hash));

    eth_abi::encode_right_padded_bytes32(&mut outer_hash, hash::keccak256(&message.sender));
    eth_abi::encode_right_padded_bytes32(&mut outer_hash, hash::keccak256(&message.data));

    let mut token_hash = vector[];
    eth_abi::encode_u256(
        &mut token_hash,
        message.token_amounts.length() as u256,
    );
    message.token_amounts.do_ref!(|token_transfer| {
        let token_transfer: &Any2SuiTokenTransfer = token_transfer;
        eth_abi::encode_bytes(&mut token_hash, token_transfer.source_pool_address);
        eth_abi::encode_address(&mut token_hash, token_transfer.dest_token_address);
        eth_abi::encode_u32(&mut token_hash, token_transfer.dest_gas_amount);
        eth_abi::encode_bytes(&mut token_hash, token_transfer.extra_data);
        eth_abi::encode_u256(&mut token_hash, token_transfer.amount);
    });
    eth_abi::encode_right_padded_bytes32(&mut outer_hash, hash::keccak256(&token_hash));

    hash::keccak256(&outer_hash)
}

// ================================================================
// |                       Deserialization                        |
// ================================================================

fun deserialize_commit_report(report_bytes: vector<u8>): CommitReport {
    let mut stream = bcs_stream::new(report_bytes);
    let token_price_updates = bcs_stream::deserialize_vector!(&mut stream, |stream| {
        let source_token_address = bcs_stream::deserialize_address(stream);
        TokenPriceUpdate {
            source_token: source_token_address,
            usd_per_token: bcs_stream::deserialize_u256(stream),
        }
    });

    let gas_price_updates = bcs_stream::deserialize_vector!(&mut stream, |stream| {
        GasPriceUpdate {
            dest_chain_selector: bcs_stream::deserialize_u64(stream),
            usd_per_unit_gas: bcs_stream::deserialize_u256(stream),
        }
    });

    let blessed_merkle_roots = parse_merkle_root(&mut stream);
    let unblessed_merkle_roots = parse_merkle_root(&mut stream);

    let rmn_signatures = bcs_stream::deserialize_vector!(
        &mut stream,
        |stream| { bcs_stream::deserialize_fixed_vector_u8(stream, 64) },
    );

    bcs_stream::assert_is_consumed(&stream);

    CommitReport {
        price_updates: PriceUpdates { token_price_updates, gas_price_updates },
        blessed_merkle_roots,
        unblessed_merkle_roots,
        rmn_signatures,
    }
}

fun parse_merkle_root(stream: &mut BCSStream): vector<MerkleRoot> {
    bcs_stream::deserialize_vector!(stream, |stream| {
        MerkleRoot {
            source_chain_selector: bcs_stream::deserialize_u64(stream),
            on_ramp_address: bcs_stream::deserialize_vector_u8(stream),
            min_seq_nr: bcs_stream::deserialize_u64(stream),
            max_seq_nr: bcs_stream::deserialize_u64(stream),
            merkle_root: bcs_stream::deserialize_fixed_vector_u8(stream, 32),
        }
    })
}

// ================================================================
// |                             OCR                              |
// ================================================================

public fun set_ocr3_config(
    ref: &CCIPObjectRef,
    state: &mut OffRampState,
    owner_cap: &OwnerCap,
    config_digest: vector<u8>,
    ocr_plugin_type: u8,
    big_f: u8,
    is_signature_verification_enabled: bool,
    signers: vector<vector<u8>>,
    transmitters: vector<address>,
) {
    verify_function_allowed(
        ref,
        string::utf8(b"offramp"),
        string::utf8(b"set_ocr3_config"),
        VERSION,
    );
    assert!(object::id(owner_cap) == ownable::owner_cap_id(&state.ownable_state), EInvalidOwnerCap);
    ocr3_base::set_ocr3_config(
        ref,
        &mut state.ocr3_base_state,
        config_digest,
        ocr_plugin_type,
        big_f,
        is_signature_verification_enabled,
        signers,
        transmitters,
    );
    after_ocr3_config_set(state, ocr_plugin_type, is_signature_verification_enabled);
}

fun after_ocr3_config_set(
    state: &mut OffRampState,
    ocr_plugin_type: u8,
    is_signature_verification_enabled: bool,
) {
    if (ocr_plugin_type == ocr3_base::ocr_plugin_type_commit()) {
        assert!(is_signature_verification_enabled, ESignatureVerificationRequiredInCommitPlugin);
        state.latest_price_sequence_number = 0;
    } else if (ocr_plugin_type == ocr3_base::ocr_plugin_type_execution()) {
        assert!(
            !is_signature_verification_enabled,
            ESignatureVerificationNotAllowedInExecutionPlugin,
        );
    };
}

public fun latest_config_details(state: &OffRampState, ocr_plugin_type: u8): OCRConfig {
    ocr3_base::latest_config_details(&state.ocr3_base_state, ocr_plugin_type)
}

public fun latest_config_digest_fields(
    cfg: OCRConfig,
): (vector<u8>, u8, u8, bool, vector<vector<u8>>, vector<address>) {
    ocr3_base::latest_config_details_fields(cfg)
}

public fun config_signers(state: &OCRConfig): vector<vector<u8>> {
    ocr3_base::config_signers(state)
}

public fun config_transmitters(state: &OCRConfig): vector<address> {
    ocr3_base::config_transmitters(state)
}

// ================================================================
// |                            Commit                            |
// ================================================================

public fun commit(
    ref: &mut CCIPObjectRef,
    state: &mut OffRampState,
    clock: &clock::Clock,
    report_context: vector<vector<u8>>,
    report: vector<u8>,
    signatures: vector<vector<u8>>,
    ctx: &mut TxContext,
) {
    verify_function_allowed(
        ref,
        string::utf8(b"offramp"),
        string::utf8(b"commit"),
        VERSION,
    );
    assert!(report_context.length() == 2, EInvalidReportContextLength);
    let commit_report = deserialize_commit_report(report);

    if (
        commit_report.price_updates.token_price_updates.length() > 0
            || commit_report.price_updates.gas_price_updates.length() > 0
    ) {
        let ocr_sequence_number = ocr3_base::deserialize_sequence_bytes(report_context[1]);
        if (state.latest_price_sequence_number < ocr_sequence_number) {
            state.latest_price_sequence_number = ocr_sequence_number;

            let mut source_tokens = vector[];
            let mut source_usd_per_token = vector[];

            commit_report.price_updates.token_price_updates.do_ref!(|token_price_update| {
                source_tokens.push_back(token_price_update.source_token);
                source_usd_per_token.push_back(token_price_update.usd_per_token);
            });

            let mut gas_dest_chain_selectors = vector[];
            let mut gas_usd_per_unit_gas = vector[];
            commit_report.price_updates.gas_price_updates.do_ref!(|gas_price_update| {
                gas_dest_chain_selectors.push_back(gas_price_update.dest_chain_selector);
                gas_usd_per_unit_gas.push_back(gas_price_update.usd_per_unit_gas);
            });

            fee_quoter::update_prices(
                ref,
                state.fee_quoter_cap.borrow(),
                clock,
                source_tokens,
                source_usd_per_token,
                gas_dest_chain_selectors,
                gas_usd_per_unit_gas,
                ctx,
            );
        } else {
            // If no non-stale valid price updates are present and the report contains no unblessed merkle roots,
            // report is stale and should be rejected. there will be no blessed merkle roots
            assert!(commit_report.unblessed_merkle_roots.length() > 0, EStaleCommitReport);
        };
    };

    // Commit the roots that do not require RMN blessing validation.
    commit_merkle_roots(ref, state, clock, commit_report.unblessed_merkle_roots, false);

    event::emit(CommitReportAccepted {
        blessed_merkle_roots: commit_report.blessed_merkle_roots,
        unblessed_merkle_roots: commit_report.unblessed_merkle_roots,
        price_updates: commit_report.price_updates,
    });

    ocr3_base::transmit(
        &state.ocr3_base_state,
        ctx.sender(),
        ocr3_base::ocr_plugin_type_commit(),
        report_context,
        report,
        signatures,
        ctx,
    )
}

fun commit_merkle_roots(
    ref: &CCIPObjectRef,
    state: &mut OffRampState,
    clock: &clock::Clock,
    merkle_roots: vector<MerkleRoot>,
    is_blessed: bool,
) {
    merkle_roots.do_ref!(|root| {
        let root: &MerkleRoot = root;
        let source_chain_selector = root.source_chain_selector;

        assert!(!rmn_remote::is_cursed_u128(ref, source_chain_selector as u128), ECursedByRmn);

        assert_source_chain_enabled(state, source_chain_selector);

        let source_chain_config = state.source_chain_configs.get_mut(&source_chain_selector);

        // If the root is blessed but RMN blessing is disabled for the source chain, or if the root is not
        // blessed but RMN blessing is enabled, we revert.
        assert!(
            is_blessed != source_chain_config.is_rmn_verification_disabled,
            ERmnBlessingMismatch,
        );

        assert!(source_chain_config.on_ramp == root.on_ramp_address, ECommitOnRampMismatch);
        assert!(
            source_chain_config.min_seq_nr == root.min_seq_nr
                    && root.min_seq_nr <= root.max_seq_nr,
            EInvalidInterval,
        );

        let merkle_root = root.merkle_root;
        assert!(merkle_root.length() == 32 && merkle_root != ZERO_MERKLE_ROOT, EInvalidRoot);

        assert!(!state.roots.contains(merkle_root), ERootAlreadyCommitted);

        source_chain_config.min_seq_nr = root.max_seq_nr + 1;
        state.roots.add(merkle_root, clock.timestamp_ms() / 1000);
    })
}

public fun get_latest_price_sequence_number(state: &OffRampState): u64 {
    state.latest_price_sequence_number
}

public fun get_merkle_root(state: &OffRampState, root: vector<u8>): u64 {
    assert!(state.roots.contains(root), EInvalidRoot);

    *table::borrow(&state.roots, root)
}

public fun get_source_chain_config(
    ref: &CCIPObjectRef,
    state: &OffRampState,
    source_chain_selector: u64,
): SourceChainConfig {
    verify_function_allowed(
        ref,
        string::utf8(b"offramp"),
        string::utf8(b"get_source_chain_config"),
        VERSION,
    );
    if (state.source_chain_configs.contains(&source_chain_selector)) {
        let source_chain_config = state.source_chain_configs.get(&source_chain_selector);
        *source_chain_config
    } else {
        SourceChainConfig {
            router: @0x0,
            is_enabled: false,
            min_seq_nr: 0,
            is_rmn_verification_disabled: false,
            on_ramp: vector[],
        }
    }
}

public fun get_source_chain_config_fields(
    ref: &CCIPObjectRef,
    source_chain_config: SourceChainConfig,
): (address, bool, u64, bool, vector<u8>) {
    verify_function_allowed(
        ref,
        string::utf8(b"offramp"),
        string::utf8(b"get_source_chain_config_fields"),
        VERSION,
    );
    (
        source_chain_config.router,
        source_chain_config.is_enabled,
        source_chain_config.min_seq_nr,
        source_chain_config.is_rmn_verification_disabled,
        source_chain_config.on_ramp,
    )
}

public fun get_all_source_chain_configs(
    ref: &CCIPObjectRef,
    state: &OffRampState,
): (vector<u64>, vector<SourceChainConfig>) {
    verify_function_allowed(
        ref,
        string::utf8(b"offramp"),
        string::utf8(b"get_all_source_chain_configs"),
        VERSION,
    );
    let mut chain_selectors = vector[];
    let mut chain_configs = vector[];
    let keys = state.source_chain_configs.keys();
    keys.do_ref!(|key| {
        chain_selectors.push_back(*key);
        chain_configs.push_back(*state.source_chain_configs.get(key));
    });
    (chain_selectors, chain_configs)
}

// ================================================================
// |                           Config                             |
// ================================================================

public fun get_static_config(ref: &CCIPObjectRef, state: &OffRampState): StaticConfig {
    verify_function_allowed(
        ref,
        string::utf8(b"offramp"),
        string::utf8(b"get_static_config"),
        VERSION,
    );
    create_static_config(state.chain_selector)
}

public fun get_static_config_fields(
    ref: &CCIPObjectRef,
    cfg: StaticConfig,
): (u64, address, address, address) {
    verify_function_allowed(
        ref,
        string::utf8(b"offramp"),
        string::utf8(b"get_static_config_fields"),
        VERSION,
    );
    (cfg.chain_selector, cfg.rmn_remote, cfg.token_admin_registry, cfg.nonce_manager)
}

public fun get_dynamic_config(ref: &CCIPObjectRef, state: &OffRampState): DynamicConfig {
    verify_function_allowed(
        ref,
        string::utf8(b"offramp"),
        string::utf8(b"get_dynamic_config"),
        VERSION,
    );
    create_dynamic_config(state.permissionless_execution_threshold_seconds)
}

public fun get_dynamic_config_fields(ref: &CCIPObjectRef, cfg: DynamicConfig): (address, u32) {
    verify_function_allowed(
        ref,
        string::utf8(b"offramp"),
        string::utf8(b"get_dynamic_config_fields"),
        VERSION,
    );
    (cfg.fee_quoter, cfg.permissionless_execution_threshold_seconds)
}

public fun set_dynamic_config(
    ref: &CCIPObjectRef,
    state: &mut OffRampState,
    owner_cap: &OwnerCap,
    permissionless_execution_threshold_seconds: u32,
) {
    verify_function_allowed(
        ref,
        string::utf8(b"offramp"),
        string::utf8(b"set_dynamic_config"),
        VERSION,
    );
    assert!(object::id(owner_cap) == ownable::owner_cap_id(&state.ownable_state), EInvalidOwnerCap);
    set_dynamic_config_internal(
        state,
        permissionless_execution_threshold_seconds,
    )
}

fun create_static_config(chain_selector: u64): StaticConfig {
    StaticConfig {
        chain_selector,
        rmn_remote: @ccip,
        token_admin_registry: @ccip,
        nonce_manager: @ccip,
    }
}

public fun apply_source_chain_config_updates(
    ref: &CCIPObjectRef,
    state: &mut OffRampState,
    owner_cap: &OwnerCap,
    source_chains_selector: vector<u64>,
    source_chains_is_enabled: vector<bool>,
    source_chains_is_rmn_verification_disabled: vector<bool>,
    source_chains_on_ramp: vector<vector<u8>>,
    ctx: &mut TxContext,
) {
    verify_function_allowed(
        ref,
        string::utf8(b"offramp"),
        string::utf8(b"apply_source_chain_config_updates"),
        VERSION,
    );
    assert!(object::id(owner_cap) == ownable::owner_cap_id(&state.ownable_state), EInvalidOwnerCap);
    apply_source_chain_config_updates_internal(
        state,
        source_chains_selector,
        source_chains_is_enabled,
        source_chains_is_rmn_verification_disabled,
        source_chains_on_ramp,
        ctx,
    )
}

public fun get_ccip_package_id(): address {
    @ccip
}

// ================================================================
// |                      Ownable Functions                       |
// ================================================================

public fun owner(state: &OffRampState): address {
    ownable::owner(&state.ownable_state)
}

public fun has_pending_transfer(state: &OffRampState): bool {
    ownable::has_pending_transfer(&state.ownable_state)
}

public fun pending_transfer_from(state: &OffRampState): Option<address> {
    ownable::pending_transfer_from(&state.ownable_state)
}

public fun pending_transfer_to(state: &OffRampState): Option<address> {
    ownable::pending_transfer_to(&state.ownable_state)
}

public fun pending_transfer_accepted(state: &OffRampState): Option<bool> {
    ownable::pending_transfer_accepted(&state.ownable_state)
}

public fun transfer_ownership(
    ref: &CCIPObjectRef,
    state: &mut OffRampState,
    owner_cap: &OwnerCap,
    new_owner: address,
    ctx: &mut TxContext,
) {
    verify_function_allowed(
        ref,
        string::utf8(b"offramp"),
        string::utf8(b"transfer_ownership"),
        VERSION,
    );
    ownable::transfer_ownership(owner_cap, &mut state.ownable_state, new_owner, ctx);
}

public fun accept_ownership(ref: &CCIPObjectRef, state: &mut OffRampState, ctx: &mut TxContext) {
    verify_function_allowed(
        ref,
        string::utf8(b"offramp"),
        string::utf8(b"accept_ownership"),
        VERSION,
    );
    ownable::accept_ownership(&mut state.ownable_state, ctx);
}

public fun accept_ownership_from_object(
    ref: &CCIPObjectRef,
    state: &mut OffRampState,
    from: &mut UID,
    ctx: &mut TxContext,
) {
    verify_function_allowed(
        ref,
        string::utf8(b"offramp"),
        string::utf8(b"accept_ownership_from_object"),
        VERSION,
    );
    ownable::accept_ownership_from_object(&mut state.ownable_state, from, ctx);
}

public fun mcms_accept_ownership(
    ref: &CCIPObjectRef,
    state: &mut OffRampState,
    registry: &mut Registry,
    params: ExecutingCallbackParams,
    ctx: &mut TxContext,
) {
    verify_function_allowed(
        ref,
        string::utf8(b"offramp"),
        string::utf8(b"mcms_accept_ownership"),
        VERSION,
    );
    let data = mcms_registry::get_accept_ownership_data(
        registry,
        params,
        McmsAcceptOwnershipProof {},
    );

    let mut stream = bcs_stream::new(data);
    bcs_stream::validate_obj_addr(object::id_address(state), &mut stream);

    bcs_stream::assert_is_consumed(&stream);

    let mcms = mcms_registry::get_multisig_address();
    ownable::mcms_accept_ownership(&mut state.ownable_state, mcms, ctx);
}

public fun execute_ownership_transfer(
    ref: &CCIPObjectRef,
    owner_cap: OwnerCap,
    state: &mut OffRampState,
    to: address,
    ctx: &mut TxContext,
) {
    verify_function_allowed(
        ref,
        string::utf8(b"offramp"),
        string::utf8(b"execute_ownership_transfer"),
        VERSION,
    );
    ownable::execute_ownership_transfer(owner_cap, &mut state.ownable_state, to, ctx);
}

public fun execute_ownership_transfer_to_mcms(
    ref: &CCIPObjectRef,
    owner_cap: OwnerCap,
    state: &mut OffRampState,
    registry: &mut Registry,
    to: address,
    ctx: &mut TxContext,
) {
    verify_function_allowed(
        ref,
        string::utf8(b"offramp"),
        string::utf8(b"execute_ownership_transfer_to_mcms"),
        VERSION,
    );

    let publisher_wrapper = mcms_registry::create_publisher_wrapper(
        ownable::borrow_publisher(&owner_cap),
        McmsCallback {},
    );

    ownable::execute_ownership_transfer_to_mcms(
        owner_cap,
        &mut state.ownable_state,
        registry,
        to,
        publisher_wrapper,
        McmsCallback {},
        vector[b"offramp"],
        ctx,
    );
}

public fun mcms_register_upgrade_cap(
    ref: &CCIPObjectRef,
    upgrade_cap: UpgradeCap,
    registry: &mut Registry,
    state: &mut DeployerState,
    ctx: &mut TxContext,
) {
    verify_function_allowed(
        ref,
        string::utf8(b"offramp"),
        string::utf8(b"mcms_register_upgrade_cap"),
        VERSION,
    );
    mcms_deployer::register_upgrade_cap(
        state,
        registry,
        upgrade_cap,
        ctx,
    );
}

// ================================================================
// |                      MCMS Entrypoint                         |
// ================================================================

public struct McmsCallback has drop {}

/// Proof for MCMS Accept Ownership
public struct McmsAcceptOwnershipProof has drop {}

public fun mcms_add_package_id(
    state: &mut OffRampState,
    registry: &mut Registry,
    params: ExecutingCallbackParams,
) {
    let (owner_cap, function, data) = mcms_registry::get_callback_params_with_caps<
        McmsCallback,
        OwnerCap,
    >(
        registry,
        McmsCallback {},
        params,
    );
    assert!(function == string::utf8(b"add_package_id"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    bcs_stream::validate_obj_addrs(
        vector[object::id_address(state), object::id_address(owner_cap)],
        &mut stream,
    );
    let package_id = bcs_stream::deserialize_address(&mut stream);
    bcs_stream::assert_is_consumed(&stream);

    add_package_id(state, owner_cap, package_id);
}

public fun mcms_remove_package_id(
    state: &mut OffRampState,
    registry: &mut Registry,
    params: ExecutingCallbackParams,
) {
    let (owner_cap, function, data) = mcms_registry::get_callback_params_with_caps<
        McmsCallback,
        OwnerCap,
    >(
        registry,
        McmsCallback {},
        params,
    );
    assert!(function == string::utf8(b"remove_package_id"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    bcs_stream::validate_obj_addrs(
        vector[object::id_address(state), object::id_address(owner_cap)],
        &mut stream,
    );
    let package_id = bcs_stream::deserialize_address(&mut stream);
    bcs_stream::assert_is_consumed(&stream);

    remove_package_id(state, owner_cap, package_id);
}

public fun mcms_set_dynamic_config(
    ref: &CCIPObjectRef,
    state: &mut OffRampState,
    registry: &mut Registry,
    params: ExecutingCallbackParams,
) {
    let (owner_cap, function, data) = mcms_registry::get_callback_params_with_caps<
        McmsCallback,
        OwnerCap,
    >(
        registry,
        McmsCallback {},
        params,
    );
    assert!(function == string::utf8(b"set_dynamic_config"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    bcs_stream::validate_obj_addrs(
        vector[object::id_address(ref), object::id_address(state), object::id_address(owner_cap)],
        &mut stream,
    );

    let permissionless_execution_threshold_seconds = bcs_stream::deserialize_u32(&mut stream);
    bcs_stream::assert_is_consumed(&stream);

    set_dynamic_config(ref, state, owner_cap, permissionless_execution_threshold_seconds);
}

public fun mcms_apply_source_chain_config_updates(
    ref: &CCIPObjectRef,
    state: &mut OffRampState,
    registry: &mut Registry,
    params: ExecutingCallbackParams,
    ctx: &mut TxContext,
) {
    let (owner_cap, function, data) = mcms_registry::get_callback_params_with_caps<
        McmsCallback,
        OwnerCap,
    >(
        registry,
        McmsCallback {},
        params,
    );
    assert!(function == string::utf8(b"apply_source_chain_config_updates"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    bcs_stream::validate_obj_addrs(
        vector[object::id_address(ref), object::id_address(state), object::id_address(owner_cap)],
        &mut stream,
    );

    let source_chains_selector = bcs_stream::deserialize_vector!(
        &mut stream,
        |stream| bcs_stream::deserialize_u64(stream),
    );
    let source_chains_is_enabled = bcs_stream::deserialize_vector!(
        &mut stream,
        |stream| bcs_stream::deserialize_bool(stream),
    );
    let source_chains_is_rmn_verification_disabled = bcs_stream::deserialize_vector!(
        &mut stream,
        |stream| bcs_stream::deserialize_bool(stream),
    );
    let source_chains_on_ramp = bcs_stream::deserialize_vector!(
        &mut stream,
        |stream| bcs_stream::deserialize_vector_u8(stream),
    );
    bcs_stream::assert_is_consumed(&stream);

    apply_source_chain_config_updates(
        ref,
        state,
        owner_cap,
        source_chains_selector,
        source_chains_is_enabled,
        source_chains_is_rmn_verification_disabled,
        source_chains_on_ramp,
        ctx,
    );
}

public fun mcms_set_ocr3_config(
    ref: &CCIPObjectRef,
    state: &mut OffRampState,
    registry: &mut Registry,
    params: ExecutingCallbackParams,
) {
    let (owner_cap, function, data) = mcms_registry::get_callback_params_with_caps<
        McmsCallback,
        OwnerCap,
    >(
        registry,
        McmsCallback {},
        params,
    );
    assert!(function == string::utf8(b"set_ocr3_config"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    bcs_stream::validate_obj_addrs(
        vector[object::id_address(ref), object::id_address(state), object::id_address(owner_cap)],
        &mut stream,
    );

    let config_digest = bcs_stream::deserialize_vector_u8(&mut stream);
    let ocr_plugin_type = bcs_stream::deserialize_u8(&mut stream);
    let big_f = bcs_stream::deserialize_u8(&mut stream);
    let is_signature_verification_enabled = bcs_stream::deserialize_bool(&mut stream);
    let signers = bcs_stream::deserialize_vector!(
        &mut stream,
        |stream| bcs_stream::deserialize_vector_u8(stream),
    );
    let transmitters = bcs_stream::deserialize_vector!(
        &mut stream,
        |stream| bcs_stream::deserialize_address(stream),
    );
    bcs_stream::assert_is_consumed(&stream);

    set_ocr3_config(
        ref,
        state,
        owner_cap,
        config_digest,
        ocr_plugin_type,
        big_f,
        is_signature_verification_enabled,
        signers,
        transmitters,
    );
}

public fun mcms_transfer_ownership(
    ref: &CCIPObjectRef,
    state: &mut OffRampState,
    registry: &mut Registry,
    params: ExecutingCallbackParams,
    ctx: &mut TxContext,
) {
    let (owner_cap, function, data) = mcms_registry::get_callback_params_with_caps<
        McmsCallback,
        OwnerCap,
    >(
        registry,
        McmsCallback {},
        params,
    );
    assert!(function == string::utf8(b"transfer_ownership"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    bcs_stream::validate_obj_addrs(
        vector[object::id_address(ref), object::id_address(state), object::id_address(owner_cap)],
        &mut stream,
    );

    let to = bcs_stream::deserialize_address(&mut stream);
    bcs_stream::assert_is_consumed(&stream);

    transfer_ownership(ref, state, owner_cap, to, ctx);
}

public fun mcms_execute_ownership_transfer(
    ref: &CCIPObjectRef,
    state: &mut OffRampState,
    registry: &mut Registry,
    deployer_state: &mut DeployerState,
    params: ExecutingCallbackParams,
    ctx: &mut TxContext,
) {
    let (_owner_cap, function, data) = mcms_registry::get_callback_params_with_caps<
        McmsCallback,
        OwnerCap,
    >(
        registry,
        McmsCallback {},
        params,
    );
    assert!(function == string::utf8(b"execute_ownership_transfer"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    bcs_stream::validate_obj_addrs(
        vector[object::id_address(ref), object::id_address(_owner_cap), object::id_address(state)],
        &mut stream,
    );

    let to = bcs_stream::deserialize_address(&mut stream);
    let package_address = bcs_stream::deserialize_address(&mut stream);
    bcs_stream::assert_is_consumed(&stream);

    let owner_cap = mcms_registry::release_cap(registry, McmsCallback {});

    if (mcms_deployer::has_upgrade_cap(deployer_state, package_address)) {
        let upgrade_cap = mcms_deployer::release_upgrade_cap(
            deployer_state,
            registry,
            McmsCallback {},
        );
        transfer::public_transfer(upgrade_cap, to);
    };

    execute_ownership_transfer(ref, owner_cap, state, to, ctx);
}

public fun mcms_add_allowed_modules(
    registry: &mut Registry,
    params: ExecutingCallbackParams,
    ctx: &mut TxContext,
) {
    let (_owner_cap, function, data) = mcms_registry::get_callback_params_with_caps<
        McmsCallback,
        OwnerCap,
    >(
        registry,
        McmsCallback {},
        params,
    );
    assert!(function == string::utf8(b"add_allowed_modules"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    bcs_stream::validate_obj_addr(object::id_address(registry), &mut stream);

    let new_module_names = bcs_stream::deserialize_vector!(
        &mut stream,
        |stream| bcs_stream::deserialize_vector_u8(stream),
    );
    bcs_stream::assert_is_consumed(&stream);

    mcms_registry::add_allowed_modules(registry, McmsCallback {}, new_module_names, ctx);
}

public fun mcms_remove_allowed_modules(
    registry: &mut Registry,
    params: ExecutingCallbackParams,
    ctx: &mut TxContext,
) {
    let (_owner_cap, function, data) = mcms_registry::get_callback_params_with_caps<
        McmsCallback,
        OwnerCap,
    >(
        registry,
        McmsCallback {},
        params,
    );
    assert!(function == string::utf8(b"remove_allowed_modules"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    bcs_stream::validate_obj_addr(object::id_address(registry), &mut stream);

    let module_names = bcs_stream::deserialize_vector!(
        &mut stream,
        |stream| bcs_stream::deserialize_vector_u8(stream),
    );
    bcs_stream::assert_is_consumed(&stream);

    mcms_registry::remove_allowed_modules(registry, McmsCallback {}, module_names, ctx);
}

// ============================== Test Functions ============================== //

#[test_only]
public fun test_init(ctx: &mut TxContext) {
    init(OFFRAMP {}, ctx);
}
