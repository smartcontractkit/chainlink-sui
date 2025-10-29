module usdc_token_pool::usdc_token_pool;

use ccip::eth_abi;
use ccip::offramp_state_helper as offramp_sh;
use ccip::onramp_state_helper as onramp_sh;
use ccip::state_object::{Self, CCIPObjectRef};
use ccip::token_admin_registry;
use mcms::bcs_stream;
use mcms::mcms_deployer::{Self, DeployerState};
use mcms::mcms_registry::{Self, Registry, ExecutingCallbackParams};
use message_transmitter::auth::auth_caller_identifier;
use message_transmitter::message;
use message_transmitter::receive_message::{Self, Receipt, ReceiveMessageTicket};
use message_transmitter::state::State as MessageTransmitterState;
use stablecoin::treasury::Treasury;
use std::string::{Self, String};
use std::type_name;
use sui::address;
use sui::clock::Clock;
use sui::coin::{Coin, CoinMetadata};
use sui::deny_list::DenyList;
use sui::event;
use sui::package::{Self, Publisher, UpgradeCap};
use sui::table::{Self, Table};
use token_messenger_minter::burn_message;
use token_messenger_minter::deposit_for_burn::{Self, DepositForBurnWithCallerTicket};
use token_messenger_minter::handle_receive_message;
use token_messenger_minter::state::State as MinterState;
use usdc_token_pool::ownable::{Self, OwnerCap, OwnableState};
use usdc_token_pool::token_pool::{Self, TokenPoolState};

public struct USDC_TOKEN_POOL has drop {}

fun init(otw: USDC_TOKEN_POOL, ctx: &mut TxContext) {
    package::claim_and_keep(otw, ctx);
}

// We restrict to the first version. New pool may be required for subsequent versions.
const SUPPORTED_USDC_VERSION_U64: u64 = 0;

const CLOCK_ADDRESS: address = @0x6;
const DENY_LIST_ADDRESS: address = @0x403;

/// A domain is a USDC representation of a destination chain.
/// @dev Zero is a valid domain identifier.
/// @dev The address to mint on the destination chain is the corresponding USDC pool.
/// @dev The allowedCaller represents the contract authorized to call receiveMessage on the destination CCTP message transmitter.
/// For EVM dest pool version 1.6.1, this is the MessageTransmitterProxy of the destination chain.
/// For EVM dest pool version 1.5.1, this is the destination chain's token pool.
public struct Domain has copy, drop, store {
    allowed_caller: vector<u8>, //  Address allowed to mint on the domain
    domain_identifier: u32, // Unique domain ID
    enabled: bool,
}

public struct DomainsSet has copy, drop {
    allowed_caller: vector<u8>,
    domain_identifier: u32,
    remote_chain_selector: u64,
    enabled: bool,
}

public struct USDCTokenPoolState<phantom T> has key {
    id: UID,
    token_pool_state: TokenPoolState,
    chain_to_domain: Table<u64, Domain>,
    local_domain_identifier: u32,
    ownable_state: OwnableState<T>,
}

const EInvalidCoinMetadata: u64 = 1;
const EInvalidArguments: u64 = 2;
const EInvalidOwnerCap: u64 = 4;
const EZeroChainSelector: u64 = 5;
const EEmptyAllowedCaller: u64 = 6;
const EInvalidMessageVersion: u64 = 7;
const EDomainMismatch: u64 = 8;
const ENonceMismatch: u64 = 9;
const EDomainNotFound: u64 = 10;
const EDomainDisabled: u64 = 11;
const ETokenAmountOverflow: u64 = 12;
const EInvalidFunction: u64 = 13;
const EPoolStillRegistered: u64 = 14;
const EInvalidProof: u64 = 15;
const EInvalidPackageId: u64 = 16;
const EInvalidModuleName: u64 = 17;
const EInvalidFunctionName: u64 = 18;

// ================================================================
// |                             Init                             |
// ================================================================

public fun type_and_version(): String {
    string::utf8(b"USDCTokenPool 1.6.0")
}

#[allow(lint(self_transfer))]
public fun initialize_by_ccip_admin<T: drop>(
    ref: &mut CCIPObjectRef,
    mut ccip_admin_proof: state_object::CCIPAdminProof,
    coin_metadata: &CoinMetadata<T>, // this can be provided as an address or in Move.toml
    publisher: Publisher,
    ctx: &mut TxContext,
) {
    assert!(!state_object::get_ccip_admin_proof_validated(&ccip_admin_proof), EInvalidProof);

    let data = state_object::get_ccip_admin_proof_data(&ccip_admin_proof);
    let mut stream = bcs_stream::new(data);

    let target_package_id = bcs_stream::deserialize_address(&mut stream);
    let target_module_name = bcs_stream::deserialize_string(&mut stream);
    let target_function_name = bcs_stream::deserialize_string(&mut stream);
    let local_domain_identifier = bcs_stream::deserialize_u32(&mut stream);
    let token_pool_package_id = bcs_stream::deserialize_address(&mut stream);
    let token_pool_administrator = bcs_stream::deserialize_address(&mut stream);
    bcs_stream::assert_is_consumed(&stream);

    assert!(target_package_id == @usdc_token_pool, EInvalidPackageId);
    assert!(target_module_name == string::utf8(b"usdc_token_pool"), EInvalidModuleName);
    assert!(
        target_function_name == string::utf8(b"initialize_by_ccip_admin"),
        EInvalidFunctionName,
    );

    state_object::set_ccip_admin_proof_validated(&mut ccip_admin_proof, true);

    let coin_metadata_address: address = object::id_to_address(&object::id(coin_metadata));
    assert!(coin_metadata_address == @usdc_coin_metadata_object_id, EInvalidCoinMetadata);

    let (ownable_state, mut token_pool_owner_cap) = ownable::new(ctx);
    ownable::attach_publisher(&mut token_pool_owner_cap, publisher);

    let usdc_token_pool = USDCTokenPoolState<T> {
        id: object::new(ctx),
        token_pool_state: token_pool::initialize(
            coin_metadata_address,
            coin_metadata.get_decimals(),
            vector[],
            ctx,
        ),
        chain_to_domain: table::new(ctx),
        local_domain_identifier,
        ownable_state,
    };

    let token_type = type_name::with_defining_ids<T>();
    let proof_type = type_name::with_defining_ids<TypeProof>();
    let token_pool_state_address = object::id_to_address(&object::id(&usdc_token_pool));

    token_admin_registry::register_pool_by_admin(
        ref,
        ccip_admin_proof,
        coin_metadata_address,
        token_pool_package_id,
        string::utf8(b"usdc_token_pool"),
        token_type.into_string(),
        token_pool_administrator,
        proof_type.into_string(),
        // these addresses match the lock_or_burn and release_or_mint functions' last 6 arguments, excluding the ctx
        vector[
            CLOCK_ADDRESS,
            DENY_LIST_ADDRESS,
            token_pool_state_address,
            @token_messenger_minter_state,
            @message_transmitter_state,
            @treasury,
        ],
        vector[
            CLOCK_ADDRESS,
            DENY_LIST_ADDRESS,
            token_pool_state_address,
            @token_messenger_minter_state,
            @message_transmitter_state,
            @treasury,
        ],
        ctx,
    );

    transfer::share_object(usdc_token_pool);
    transfer::public_transfer(token_pool_owner_cap, ctx.sender());
}

public fun set_pool<T>(
    ref: &mut CCIPObjectRef,
    state: &mut USDCTokenPoolState<T>,
    owner_cap: &OwnerCap<T>,
    coin_metadata_address: address,
    ctx: &mut TxContext,
) {
    set_pool_internal(ref, state, owner_cap, coin_metadata_address, ctx.sender());
}

fun set_pool_internal<T>(
    ref: &mut CCIPObjectRef,
    state: &USDCTokenPoolState<T>,
    owner_cap: &OwnerCap<T>,
    coin_metadata_address: address,
    caller: address,
) {
    assert!(object::id(owner_cap) == ownable::owner_cap_id(&state.ownable_state), EInvalidOwnerCap);

    let token_pool_state_address = object::uid_to_address(&state.id);
    token_admin_registry::set_pool(
        ref,
        coin_metadata_address,
        vector[
            CLOCK_ADDRESS,
            DENY_LIST_ADDRESS,
            token_pool_state_address,
            @token_messenger_minter_state,
            @message_transmitter_state,
            @treasury,
        ],
        vector[
            CLOCK_ADDRESS,
            DENY_LIST_ADDRESS,
            token_pool_state_address,
            @token_messenger_minter_state,
            @message_transmitter_state,
            @treasury,
        ],
        TypeProof {},
        caller,
    );
}

// ================================================================
// |                 Exposing token_pool functions                |
// ================================================================

// this now returns the address of coin metadata
public fun get_token<T>(state: &USDCTokenPoolState<T>): address {
    token_pool::get_token(&state.token_pool_state)
}

public fun get_token_decimals<T>(state: &USDCTokenPoolState<T>): u8 {
    state.token_pool_state.get_local_decimals()
}

public fun get_remote_pools<T>(
    state: &USDCTokenPoolState<T>,
    remote_chain_selector: u64,
): vector<vector<u8>> {
    token_pool::get_remote_pools(&state.token_pool_state, remote_chain_selector)
}

public fun is_remote_pool<T>(
    state: &USDCTokenPoolState<T>,
    remote_chain_selector: u64,
    remote_pool_address: vector<u8>,
): bool {
    token_pool::is_remote_pool(
        &state.token_pool_state,
        remote_chain_selector,
        remote_pool_address,
    )
}

public fun get_remote_token<T>(
    state: &USDCTokenPoolState<T>,
    remote_chain_selector: u64,
): vector<u8> {
    token_pool::get_remote_token(&state.token_pool_state, remote_chain_selector)
}

public fun add_remote_pool<T>(
    state: &mut USDCTokenPoolState<T>,
    owner_cap: &OwnerCap<T>,
    remote_chain_selector: u64,
    remote_pool_address: vector<u8>,
) {
    assert!(object::id(owner_cap) == ownable::owner_cap_id(&state.ownable_state), EInvalidOwnerCap);
    token_pool::add_remote_pool(
        &mut state.token_pool_state,
        remote_chain_selector,
        remote_pool_address,
    );
}

public fun remove_remote_pool<T>(
    state: &mut USDCTokenPoolState<T>,
    owner_cap: &OwnerCap<T>,
    remote_chain_selector: u64,
    remote_pool_address: vector<u8>,
) {
    assert!(object::id(owner_cap) == ownable::owner_cap_id(&state.ownable_state), EInvalidOwnerCap);
    token_pool::remove_remote_pool(
        &mut state.token_pool_state,
        remote_chain_selector,
        remote_pool_address,
    );
}

public fun is_supported_chain<T>(state: &USDCTokenPoolState<T>, remote_chain_selector: u64): bool {
    token_pool::is_supported_chain(&state.token_pool_state, remote_chain_selector)
}

public fun get_supported_chains<T>(state: &USDCTokenPoolState<T>): vector<u64> {
    token_pool::get_supported_chains(&state.token_pool_state)
}

public fun apply_chain_updates<T>(
    state: &mut USDCTokenPoolState<T>,
    owner_cap: &OwnerCap<T>,
    remote_chain_selectors_to_remove: vector<u64>,
    remote_chain_selectors_to_add: vector<u64>,
    remote_pool_addresses_to_add: vector<vector<vector<u8>>>,
    remote_token_addresses_to_add: vector<vector<u8>>,
) {
    assert!(object::id(owner_cap) == ownable::owner_cap_id(&state.ownable_state), EInvalidOwnerCap);
    token_pool::apply_chain_updates(
        &mut state.token_pool_state,
        remote_chain_selectors_to_remove,
        remote_chain_selectors_to_add,
        remote_pool_addresses_to_add,
        remote_token_addresses_to_add,
    );
}

public fun get_allowlist_enabled<T>(state: &USDCTokenPoolState<T>): bool {
    token_pool::get_allowlist_enabled(&state.token_pool_state)
}

public fun get_allowlist<T>(state: &USDCTokenPoolState<T>): vector<address> {
    token_pool::get_allowlist(&state.token_pool_state)
}

public fun set_allowlist_enabled<T>(
    state: &mut USDCTokenPoolState<T>,
    owner_cap: &OwnerCap<T>,
    enabled: bool,
) {
    assert!(object::id(owner_cap) == ownable::owner_cap_id(&state.ownable_state), EInvalidOwnerCap);
    token_pool::set_allowlist_enabled(&mut state.token_pool_state, enabled);
}

public fun apply_allowlist_updates<T>(
    state: &mut USDCTokenPoolState<T>,
    owner_cap: &OwnerCap<T>,
    removes: vector<address>,
    adds: vector<address>,
) {
    assert!(object::id(owner_cap) == ownable::owner_cap_id(&state.ownable_state), EInvalidOwnerCap);
    token_pool::apply_allowlist_updates(&mut state.token_pool_state, removes, adds);
}

// ================================================================
// |                         Burn/Mint                            |
// ================================================================

public struct TypeProof has drop {}

// This function calculates the package auth caller based on the TypeProof defined in the pool.
// When sending USDC to Sui chain, the destination caller needs to be set to the package auth caller.
// CCTP will validate that the destination caller set by the source chain matches the package auth caller.
// See https://github.com/circlefin/sui-cctp/blob/70290f70d7c3d6caf23a91b379cac08a20f0d762/packages/message_transmitter/sources/receive_message.move#L118-L150
// and https://developers.circle.com/cctp/sui-packages#destination-callers-for-sui-as-destination-chain for more details.
// This is critical for Sui because we cannot set destination caller to a single CL node.
public fun get_package_auth_caller<TypeProof: drop>(): address {
    auth_caller_identifier<TypeProof>()
}

public fun lock_or_burn<T: drop>(
    ref: &CCIPObjectRef,
    token_transfer_params: &mut onramp_sh::TokenTransferParams,
    c: Coin<T>,
    remote_chain_selector: u64,
    clock: &Clock,
    deny_list: &DenyList,
    pool: &mut USDCTokenPoolState<T>,
    state: &MinterState,
    message_transmitter_state: &mut MessageTransmitterState,
    treasury: &mut Treasury<T>,
    ctx: &mut TxContext,
) {
    let amount = c.value();
    let sender = ctx.sender();
    let mint_recipient = address::from_bytes(onramp_sh::get_token_receiver(token_transfer_params));

    assert!(pool.chain_to_domain.contains(remote_chain_selector), EDomainNotFound);
    let remote_domain_info = pool.chain_to_domain.borrow(remote_chain_selector);
    assert!(remote_domain_info.enabled, EDomainDisabled);

    // This metod validates various aspects of the lock or burn operation. If any of the
    // validations fail, the transaction will abort.
    let dest_token_address = token_pool::validate_lock_or_burn(
        ref,
        clock,
        &mut pool.token_pool_state,
        sender,
        remote_chain_selector,
        amount,
    );

    let ticket: DepositForBurnWithCallerTicket<
        T,
        TypeProof,
    > = deposit_for_burn::create_deposit_for_burn_with_caller_ticket(
        TypeProof {},
        c,
        remote_domain_info.domain_identifier,
        mint_recipient,
        address::from_bytes(remote_domain_info.allowed_caller),
    );

    let (_, msg) = deposit_for_burn::deposit_for_burn_with_caller_with_package_auth(
        ticket,
        state,
        message_transmitter_state,
        deny_list,
        treasury,
        ctx,
    );

    let nonce = message::nonce(&msg);
    let source_pool_data = encode_source_pool_data(pool.local_domain_identifier, nonce);

    token_pool::emit_locked_or_burned(&pool.token_pool_state, amount, remote_chain_selector);

    onramp_sh::add_token_transfer_param(
        ref,
        token_transfer_params,
        remote_chain_selector,
        amount,
        get_token(pool),
        dest_token_address,
        source_pool_data,
        TypeProof {},
    )
}

public fun release_or_mint<T: drop>(
    ref: &CCIPObjectRef,
    receiver_params: &mut offramp_sh::ReceiverParams,
    clock: &Clock,
    deny_list: &DenyList,
    pool: &mut USDCTokenPoolState<T>,
    state: &mut MinterState,
    message_transmitter_state: &mut MessageTransmitterState,
    treasury: &mut Treasury<T>,
    ctx: &mut TxContext,
) {
    let (
        receiver,
        remote_chain_selector,
        _,
        dest_token_address,
        _,
        source_pool_address,
        source_pool_data,
        offchain_token_data,
    ) = offramp_sh::get_dest_token_transfer_data(receiver_params);
    let (message_bytes, attestation) = parse_message_and_attestation(offchain_token_data);

    // Prepare the ReceiveMessageTicket by calling create_receive_message_ticket() from within your package.
    let ticket: ReceiveMessageTicket<TypeProof> = receive_message::create_receive_message_ticket(
        TypeProof {},
        message_bytes,
        attestation,
    );

    // Receive the message on MessageTransmitter.
    let receipt: Receipt = receive_message::receive_message_with_package_auth(
        ticket,
        message_transmitter_state,
    );
    let (source_domain_identifier, nonce) = decode_source_pool_data(source_pool_data);
    // local domain identifier is checked in receive_message_with_package_auth
    validate_receipt(&receipt, source_domain_identifier, nonce);

    // Pass the Receipt into TokenMessengerMinter to mint the USDC.
    let ticket_with_burn_message = handle_receive_message::handle_receive_message(
        receipt,
        state,
        deny_list,
        treasury,
        ctx,
    );

    let (
        stamp_receipt_ticket,
        burn_message,
    ) = handle_receive_message::deconstruct_stamp_receipt_ticket_with_burn_message(
        ticket_with_burn_message,
    );

    // Stamp the receipt
    let stamped_receipt = receive_message::stamp_receipt(
        stamp_receipt_ticket,
        message_transmitter_state,
    );

    // Complete the message and destroy the StampedReceipt
    receive_message::complete_receive_message(stamped_receipt, message_transmitter_state);

    let local_amount = burn_message::amount(&burn_message);
    let mut amount_op = local_amount.try_as_u64();
    assert!(amount_op.is_some(), ETokenAmountOverflow);
    let amount = amount_op.extract();

    token_pool::validate_release_or_mint(
        ref,
        clock,
        &mut pool.token_pool_state,
        remote_chain_selector,
        dest_token_address,
        source_pool_address,
        amount,
    );

    token_pool::emit_released_or_minted(
        &pool.token_pool_state,
        receiver,
        amount,
        remote_chain_selector,
    );

    offramp_sh::complete_token_transfer(
        ref,
        receiver_params,
        TypeProof {},
    );
}

fun parse_message_and_attestation(payload: vector<u8>): (vector<u8>, vector<u8>) {
    let mut stream = eth_abi::new_stream(payload);

    let message = eth_abi::decode_bytes(&mut stream);
    let attestation = eth_abi::decode_bytes(&mut stream);

    (message, attestation)
}

fun encode_source_pool_data(local_domain_identifier: u32, nonce: u64): vector<u8> {
    let mut source_pool_data = vector[];
    eth_abi::encode_u64(&mut source_pool_data, nonce);
    eth_abi::encode_u32(&mut source_pool_data, local_domain_identifier);
    source_pool_data
}

fun decode_source_pool_data(source_pool_data: vector<u8>): (u32, u64) {
    let mut stream = eth_abi::new_stream(source_pool_data);
    let nonce = eth_abi::decode_u64(&mut stream);
    let local_domain_identifier = eth_abi::decode_u32(&mut stream);

    (local_domain_identifier, nonce)
}

fun validate_receipt(receipt: &Receipt, expected_source_domain: u32, expected_nonce: u64) {
    let version = receive_message::current_version(receipt);
    assert!(version == SUPPORTED_USDC_VERSION_U64, EInvalidMessageVersion);

    let source_domain = receive_message::source_domain(receipt);
    let nonce = receive_message::nonce(receipt);

    assert!(source_domain == expected_source_domain, EDomainMismatch);

    assert!(nonce == expected_nonce, ENonceMismatch);
}

// ================================================================
// |                      USDC Domains                            |
// ================================================================

public fun get_domain<T>(pool: &USDCTokenPoolState<T>, chain_selector: u64): Domain {
    assert!(pool.chain_to_domain.contains(chain_selector), EDomainNotFound);
    *pool.chain_to_domain.borrow(chain_selector)
}

public fun set_domains<T>(
    pool: &mut USDCTokenPoolState<T>,
    owner_cap: &OwnerCap<T>,
    remote_chain_selectors: vector<u64>,
    remote_domain_identifiers: vector<u32>,
    allowed_remote_callers: vector<vector<u8>>,
    enableds: vector<bool>,
) {
    assert!(object::id(owner_cap) == ownable::owner_cap_id(&pool.ownable_state), EInvalidOwnerCap);

    let number_of_chains = remote_chain_selectors.length();

    assert!(
        number_of_chains == remote_domain_identifiers.length()
            && number_of_chains == allowed_remote_callers.length()
            && number_of_chains == enableds.length(),
        EInvalidArguments,
    );

    let mut i = 0;
    while (i < number_of_chains) {
        let allowed_caller = allowed_remote_callers[i];
        let domain_identifier = remote_domain_identifiers[i];
        let remote_chain_selector = remote_chain_selectors[i];
        let enabled = enableds[i];

        assert!(remote_chain_selector != 0, EZeroChainSelector);

        assert!(allowed_caller.length() != 0, EEmptyAllowedCaller);

        if (pool.chain_to_domain.contains(remote_chain_selector)) {
            pool.chain_to_domain.remove(remote_chain_selector);
        };
        pool
            .chain_to_domain
            .add(
                remote_chain_selector,
                Domain { allowed_caller, domain_identifier, enabled },
            );

        event::emit(DomainsSet {
            allowed_caller,
            domain_identifier,
            remote_chain_selector,
            enabled,
        });
        i = i + 1;
    };
}

// ================================================================
// |                    Rate limit config                         |
// ================================================================

public fun set_chain_rate_limiter_configs<T>(
    state: &mut USDCTokenPoolState<T>,
    owner_cap: &OwnerCap<T>,
    clock: &Clock,
    remote_chain_selectors: vector<u64>,
    outbound_is_enableds: vector<bool>,
    outbound_capacities: vector<u64>,
    outbound_rates: vector<u64>,
    inbound_is_enableds: vector<bool>,
    inbound_capacities: vector<u64>,
    inbound_rates: vector<u64>,
) {
    assert!(object::id(owner_cap) == ownable::owner_cap_id(&state.ownable_state), EInvalidOwnerCap);
    let number_of_chains = remote_chain_selectors.length();

    assert!(
        number_of_chains == outbound_is_enableds.length()
            && number_of_chains == outbound_capacities.length()
            && number_of_chains == outbound_rates.length()
            && number_of_chains == inbound_is_enableds.length()
            && number_of_chains == inbound_capacities.length()
            && number_of_chains == inbound_rates.length(),
        EInvalidArguments,
    );

    let mut i = 0;
    while (i < number_of_chains) {
        token_pool::set_chain_rate_limiter_config(
            clock,
            &mut state.token_pool_state,
            remote_chain_selectors[i],
            outbound_is_enableds[i],
            outbound_capacities[i],
            outbound_rates[i],
            inbound_is_enableds[i],
            inbound_capacities[i],
            inbound_rates[i],
        );
        i = i + 1;
    };
}

public fun set_chain_rate_limiter_config<T>(
    state: &mut USDCTokenPoolState<T>,
    owner_cap: &OwnerCap<T>,
    clock: &Clock,
    remote_chain_selector: u64,
    outbound_is_enabled: bool,
    outbound_capacity: u64,
    outbound_rate: u64,
    inbound_is_enabled: bool,
    inbound_capacity: u64,
    inbound_rate: u64,
) {
    assert!(object::id(owner_cap) == ownable::owner_cap_id(&state.ownable_state), EInvalidOwnerCap);
    token_pool::set_chain_rate_limiter_config(
        clock,
        &mut state.token_pool_state,
        remote_chain_selector,
        outbound_is_enabled,
        outbound_capacity,
        outbound_rate,
        inbound_is_enabled,
        inbound_capacity,
        inbound_rate,
    );
}

// ================================================================
// |                      Ownable Functions                       |
// ================================================================

public fun owner<T>(state: &USDCTokenPoolState<T>): address {
    ownable::owner(&state.ownable_state)
}

public fun has_pending_transfer<T>(state: &USDCTokenPoolState<T>): bool {
    ownable::has_pending_transfer(&state.ownable_state)
}

public fun pending_transfer_from<T>(state: &USDCTokenPoolState<T>): Option<address> {
    ownable::pending_transfer_from(&state.ownable_state)
}

public fun pending_transfer_to<T>(state: &USDCTokenPoolState<T>): Option<address> {
    ownable::pending_transfer_to(&state.ownable_state)
}

public fun pending_transfer_accepted<T>(state: &USDCTokenPoolState<T>): Option<bool> {
    ownable::pending_transfer_accepted(&state.ownable_state)
}

public fun transfer_ownership<T>(
    state: &mut USDCTokenPoolState<T>,
    owner_cap: &OwnerCap<T>,
    new_owner: address,
    ctx: &mut TxContext,
) {
    ownable::transfer_ownership(owner_cap, &mut state.ownable_state, new_owner, ctx);
}

public fun accept_ownership<T>(state: &mut USDCTokenPoolState<T>, ctx: &mut TxContext) {
    ownable::accept_ownership(&mut state.ownable_state, ctx);
}

public fun accept_ownership_from_object<T>(
    state: &mut USDCTokenPoolState<T>,
    from: &mut UID,
    ctx: &mut TxContext,
) {
    ownable::accept_ownership_from_object(&mut state.ownable_state, from, ctx);
}

public fun mcms_accept_ownership<T>(
    state: &mut USDCTokenPoolState<T>,
    registry: &mut Registry,
    params: ExecutingCallbackParams,
    ctx: &mut TxContext,
) {
    let data = mcms_registry::get_accept_ownership_data(
        registry,
        params,
        McmsAcceptOwnershipProof<T> {},
    );

    let mut stream = bcs_stream::new(data);
    bcs_stream::validate_obj_addr(object::id_address(state), &mut stream);
    bcs_stream::assert_is_consumed(&stream);

    let mcms = mcms_registry::get_multisig_address();
    ownable::mcms_accept_ownership(&mut state.ownable_state, mcms, ctx);
}

public fun execute_ownership_transfer<T>(
    owner_cap: OwnerCap<T>,
    state: &mut USDCTokenPoolState<T>,
    to: address,
    ctx: &mut TxContext,
) {
    ownable::execute_ownership_transfer(owner_cap, &mut state.ownable_state, to, ctx);
}

public fun execute_ownership_transfer_to_mcms<T>(
    owner_cap: OwnerCap<T>,
    state: &mut USDCTokenPoolState<T>,
    registry: &mut Registry,
    to: address,
    ctx: &mut TxContext,
) {
    let publisher_wrapper = mcms_registry::create_publisher_wrapper(
        ownable::borrow_publisher(&owner_cap),
        McmsCallback<T> {},
    );

    ownable::execute_ownership_transfer_to_mcms(
        owner_cap,
        &mut state.ownable_state,
        registry,
        to,
        publisher_wrapper,
        McmsCallback<T> {},
        vector[b"usdc_token_pool"],
        ctx,
    );
}

public fun mcms_register_upgrade_cap(
    upgrade_cap: UpgradeCap,
    registry: &mut Registry,
    state: &mut DeployerState,
    ctx: &mut TxContext,
) {
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

public struct McmsCallback<phantom T> has drop {}

/// Proof for MCMS Accept Ownership
public struct McmsAcceptOwnershipProof<phantom T> has drop {}

public fun mcms_set_allowlist_enabled<T>(
    state: &mut USDCTokenPoolState<T>,
    registry: &mut Registry,
    params: ExecutingCallbackParams,
) {
    let (owner_cap, function, data) = mcms_registry::get_callback_params_with_caps<
        McmsCallback<T>,
        OwnerCap<T>,
    >(
        registry,
        McmsCallback<T> {},
        params,
    );
    assert!(function == string::utf8(b"set_allowlist_enabled"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    bcs_stream::validate_obj_addrs(
        vector[object::id_address(state), object::id_address(owner_cap)],
        &mut stream,
    );

    let enabled = bcs_stream::deserialize_bool(&mut stream);
    bcs_stream::assert_is_consumed(&stream);

    set_allowlist_enabled(state, owner_cap, enabled);
}

public fun mcms_apply_allowlist_updates<T>(
    state: &mut USDCTokenPoolState<T>,
    registry: &mut Registry,
    params: ExecutingCallbackParams,
) {
    let (owner_cap, function, data) = mcms_registry::get_callback_params_with_caps<
        McmsCallback<T>,
        OwnerCap<T>,
    >(
        registry,
        McmsCallback<T> {},
        params,
    );
    assert!(function == string::utf8(b"apply_allowlist_updates"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    bcs_stream::validate_obj_addrs(
        vector[object::id_address(state), object::id_address(owner_cap)],
        &mut stream,
    );

    let removes = bcs_stream::deserialize_vector!(
        &mut stream,
        |stream| bcs_stream::deserialize_address(stream),
    );
    let adds = bcs_stream::deserialize_vector!(
        &mut stream,
        |stream| bcs_stream::deserialize_address(stream),
    );
    bcs_stream::assert_is_consumed(&stream);

    apply_allowlist_updates(state, owner_cap, removes, adds);
}

public fun mcms_apply_chain_updates<T>(
    state: &mut USDCTokenPoolState<T>,
    registry: &mut Registry,
    params: ExecutingCallbackParams,
) {
    let (owner_cap, function, data) = mcms_registry::get_callback_params_with_caps<
        McmsCallback<T>,
        OwnerCap<T>,
    >(
        registry,
        McmsCallback<T> {},
        params,
    );
    assert!(function == string::utf8(b"apply_chain_updates"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    bcs_stream::validate_obj_addrs(
        vector[object::id_address(state), object::id_address(owner_cap)],
        &mut stream,
    );

    let remote_chain_selectors_to_remove = bcs_stream::deserialize_vector!(
        &mut stream,
        |stream| bcs_stream::deserialize_u64(stream),
    );
    let remote_chain_selectors_to_add = bcs_stream::deserialize_vector!(
        &mut stream,
        |stream| bcs_stream::deserialize_u64(stream),
    );
    let remote_pool_addresses_to_add = bcs_stream::deserialize_vector!(
        &mut stream,
        |stream| bcs_stream::deserialize_vector!(
            stream,
            |stream| bcs_stream::deserialize_vector_u8(stream),
        ),
    );
    let remote_token_addresses_to_add = bcs_stream::deserialize_vector!(
        &mut stream,
        |stream| bcs_stream::deserialize_vector_u8(stream),
    );
    bcs_stream::assert_is_consumed(&stream);

    apply_chain_updates(
        state,
        owner_cap,
        remote_chain_selectors_to_remove,
        remote_chain_selectors_to_add,
        remote_pool_addresses_to_add,
        remote_token_addresses_to_add,
    );
}

public fun mcms_add_remote_pool<T>(
    state: &mut USDCTokenPoolState<T>,
    registry: &mut Registry,
    params: ExecutingCallbackParams,
) {
    let (owner_cap, function, data) = mcms_registry::get_callback_params_with_caps<
        McmsCallback<T>,
        OwnerCap<T>,
    >(
        registry,
        McmsCallback<T> {},
        params,
    );
    assert!(function == string::utf8(b"add_remote_pool"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    bcs_stream::validate_obj_addrs(
        vector[object::id_address(state), object::id_address(owner_cap)],
        &mut stream,
    );

    let remote_chain_selector = bcs_stream::deserialize_u64(&mut stream);
    let remote_pool_address = bcs_stream::deserialize_vector_u8(&mut stream);
    bcs_stream::assert_is_consumed(&stream);

    add_remote_pool(state, owner_cap, remote_chain_selector, remote_pool_address);
}

public fun mcms_remove_remote_pool<T>(
    state: &mut USDCTokenPoolState<T>,
    registry: &mut Registry,
    params: ExecutingCallbackParams,
) {
    let (owner_cap, function, data) = mcms_registry::get_callback_params_with_caps<
        McmsCallback<T>,
        OwnerCap<T>,
    >(
        registry,
        McmsCallback<T> {},
        params,
    );
    assert!(function == string::utf8(b"remove_remote_pool"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    bcs_stream::validate_obj_addrs(
        vector[object::id_address(state), object::id_address(owner_cap)],
        &mut stream,
    );
    let remote_chain_selector = bcs_stream::deserialize_u64(&mut stream);
    let remote_pool_address = bcs_stream::deserialize_vector_u8(&mut stream);
    bcs_stream::assert_is_consumed(&stream);

    remove_remote_pool(state, owner_cap, remote_chain_selector, remote_pool_address);
}

public fun mcms_set_chain_rate_limiter_configs<T>(
    state: &mut USDCTokenPoolState<T>,
    registry: &mut Registry,
    params: ExecutingCallbackParams,
    clock: &Clock,
) {
    let (owner_cap, function, data) = mcms_registry::get_callback_params_with_caps<
        McmsCallback<T>,
        OwnerCap<T>,
    >(
        registry,
        McmsCallback<T> {},
        params,
    );
    assert!(function == string::utf8(b"set_chain_rate_limiter_configs"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    bcs_stream::validate_obj_addrs(
        vector[object::id_address(state), object::id_address(owner_cap), object::id_address(clock)],
        &mut stream,
    );

    let remote_chain_selectors = bcs_stream::deserialize_vector!(
        &mut stream,
        |stream| bcs_stream::deserialize_u64(stream),
    );
    let outbound_is_enableds = bcs_stream::deserialize_vector!(
        &mut stream,
        |stream| bcs_stream::deserialize_bool(stream),
    );
    let outbound_capacities = bcs_stream::deserialize_vector!(
        &mut stream,
        |stream| bcs_stream::deserialize_u64(stream),
    );
    let outbound_rates = bcs_stream::deserialize_vector!(
        &mut stream,
        |stream| bcs_stream::deserialize_u64(stream),
    );
    let inbound_is_enableds = bcs_stream::deserialize_vector!(
        &mut stream,
        |stream| bcs_stream::deserialize_bool(stream),
    );
    let inbound_capacities = bcs_stream::deserialize_vector!(
        &mut stream,
        |stream| bcs_stream::deserialize_u64(stream),
    );
    let inbound_rates = bcs_stream::deserialize_vector!(
        &mut stream,
        |stream| bcs_stream::deserialize_u64(stream),
    );
    bcs_stream::assert_is_consumed(&stream);

    set_chain_rate_limiter_configs(
        state,
        owner_cap,
        clock,
        remote_chain_selectors,
        outbound_is_enableds,
        outbound_capacities,
        outbound_rates,
        inbound_is_enableds,
        inbound_capacities,
        inbound_rates,
    );
}

public fun mcms_set_chain_rate_limiter_config<T>(
    state: &mut USDCTokenPoolState<T>,
    registry: &mut Registry,
    params: ExecutingCallbackParams,
    clock: &Clock,
) {
    let (owner_cap, function, data) = mcms_registry::get_callback_params_with_caps<
        McmsCallback<T>,
        OwnerCap<T>,
    >(
        registry,
        McmsCallback<T> {},
        params,
    );
    assert!(function == string::utf8(b"set_chain_rate_limiter_config"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    bcs_stream::validate_obj_addrs(
        vector[object::id_address(state), object::id_address(owner_cap), object::id_address(clock)],
        &mut stream,
    );

    let remote_chain_selector = bcs_stream::deserialize_u64(&mut stream);
    let outbound_is_enabled = bcs_stream::deserialize_bool(&mut stream);
    let outbound_capacity = bcs_stream::deserialize_u64(&mut stream);
    let outbound_rate = bcs_stream::deserialize_u64(&mut stream);
    let inbound_is_enabled = bcs_stream::deserialize_bool(&mut stream);
    let inbound_capacity = bcs_stream::deserialize_u64(&mut stream);
    let inbound_rate = bcs_stream::deserialize_u64(&mut stream);
    bcs_stream::assert_is_consumed(&stream);

    set_chain_rate_limiter_config(
        state,
        owner_cap,
        clock,
        remote_chain_selector,
        outbound_is_enabled,
        outbound_capacity,
        outbound_rate,
        inbound_is_enabled,
        inbound_capacity,
        inbound_rate,
    );
}

/// destroy the USDC token pool state and the owner cap
/// this should only be called after unregistering the pool from the token admin registry
public fun destroy_token_pool<T>(
    ref: &mut CCIPObjectRef,
    state: USDCTokenPoolState<T>,
    owner_cap: OwnerCap<T>,
    ctx: &mut TxContext,
) {
    assert!(
        object::id(&owner_cap) == ownable::owner_cap_id(&state.ownable_state),
        EInvalidOwnerCap,
    );
    assert!(
        !token_admin_registry::is_pool_registered(ref, get_token(&state)),
        EPoolStillRegistered,
    );

    let USDCTokenPoolState<T> {
        id: state_id,
        token_pool_state,
        chain_to_domain,
        local_domain_identifier: _,
        ownable_state,
    } = state;
    token_pool::destroy_token_pool(token_pool_state);
    table::drop(chain_to_domain);
    object::delete(state_id);

    // Destroy ownable state and owner cap using helper functions
    ownable::destroy(ownable_state, owner_cap, ctx);
}

public fun mcms_transfer_ownership<T>(
    state: &mut USDCTokenPoolState<T>,
    registry: &mut Registry,
    params: ExecutingCallbackParams,
    ctx: &mut TxContext,
) {
    let (owner_cap, function, data) = mcms_registry::get_callback_params_with_caps<
        McmsCallback<T>,
        OwnerCap<T>,
    >(
        registry,
        McmsCallback<T> {},
        params,
    );
    assert!(function == string::utf8(b"transfer_ownership"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    bcs_stream::validate_obj_addrs(
        vector[object::id_address(state), object::id_address(owner_cap)],
        &mut stream,
    );

    let to = bcs_stream::deserialize_address(&mut stream);
    bcs_stream::assert_is_consumed(&stream);

    transfer_ownership(state, owner_cap, to, ctx);
}

public fun mcms_execute_ownership_transfer<T>(
    state: &mut USDCTokenPoolState<T>,
    registry: &mut Registry,
    params: ExecutingCallbackParams,
    ctx: &mut TxContext,
) {
    let (_owner_cap, function, data) = mcms_registry::get_callback_params_with_caps<
        McmsCallback<T>,
        OwnerCap<T>,
    >(
        registry,
        McmsCallback<T> {},
        params,
    );
    assert!(function == string::utf8(b"execute_ownership_transfer"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    bcs_stream::validate_obj_addrs(
        vector[object::id_address(_owner_cap), object::id_address(state)],
        &mut stream,
    );

    let to = bcs_stream::deserialize_address(&mut stream);
    bcs_stream::assert_is_consumed(&stream);

    let owner_cap = mcms_registry::release_cap<McmsCallback<T>, OwnerCap<T>>(
        registry,
        McmsCallback<T> {},
    );
    execute_ownership_transfer(owner_cap, state, to, ctx);
}

public fun mcms_set_pool<T>(
    ref: &mut CCIPObjectRef,
    state: &USDCTokenPoolState<T>,
    registry: &mut Registry,
    params: ExecutingCallbackParams,
    _: &mut TxContext,
) {
    let (owner_cap, function, data) = mcms_registry::get_callback_params_with_caps<
        McmsCallback<T>,
        OwnerCap<T>,
    >(
        registry,
        McmsCallback<T> {},
        params,
    );
    assert!(function == string::utf8(b"set_pool"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    bcs_stream::validate_obj_addrs(
        vector[object::id_address(ref), object::id_address(state), object::id_address(owner_cap)],
        &mut stream,
    );
    let coin_metadata_address = bcs_stream::deserialize_address(&mut stream);
    bcs_stream::assert_is_consumed(&stream);

    set_pool_internal(
        ref,
        state,
        owner_cap,
        coin_metadata_address,
        mcms_registry::get_multisig_address(),
    );
}

public fun mcms_add_allowed_modules<T>(
    registry: &mut Registry,
    params: ExecutingCallbackParams,
    ctx: &mut TxContext,
) {
    let (_owner_cap, function, data) = mcms_registry::get_callback_params_with_caps<
        McmsCallback<T>,
        OwnerCap<T>,
    >(
        registry,
        McmsCallback<T> {},
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

    mcms_registry::add_allowed_modules(registry, McmsCallback<T> {}, new_module_names, ctx);
}

public fun mcms_remove_allowed_modules<T>(
    registry: &mut Registry,
    params: ExecutingCallbackParams,
    ctx: &mut TxContext,
) {
    let (_owner_cap, function, data) = mcms_registry::get_callback_params_with_caps<
        McmsCallback<T>,
        OwnerCap<T>,
    >(
        registry,
        McmsCallback<T> {},
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

    mcms_registry::remove_allowed_modules(registry, McmsCallback<T> {}, module_names, ctx);
}
