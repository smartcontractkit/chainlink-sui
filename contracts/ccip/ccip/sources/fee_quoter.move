/// This module is responsible for storage and retrieval of fee token and token transfer
/// information and pricing.
module ccip::fee_quoter;

use ccip::client;
use ccip::eth_abi;
use ccip::ownable::OwnerCap;
use ccip::state_object::{Self, CCIPObjectRef};
use mcms::bcs_stream;
use mcms::mcms_registry::{Self, Registry, ExecutingCallbackParams};
use std::bcs;
use std::string::{Self, String};
use sui::clock;
use sui::event;
use sui::table;

const CHAIN_FAMILY_SELECTOR_EVM: vector<u8> = x"2812d52c";
const CHAIN_FAMILY_SELECTOR_SVM: vector<u8> = x"1e10bdc4";
const CHAIN_FAMILY_SELECTOR_APTOS: vector<u8> = x"ac77ffec";
const CHAIN_FAMILY_SELECTOR_SUI: vector<u8> = x"c4e05953";

/// @dev We disallow the first 1024 addresses to avoid calling into a range known for hosting precompiles. Calling
/// into precompiles probably won't cause any issues, but to be safe we can disallow this range. It is extremely
/// unlikely that anyone would ever be able to generate an address in this range. There is no official range of
/// precompiles, but EIP-7587 proposes to reserve the range 0x100 to 0x1ff. Our range is more conservative, even
/// though it might not be exhaustive for all chains, which is OK. We also disallow the zero address, which is a
/// common practice.
const EVM_PRECOMPILE_SPACE: u256 = 1024;

/// @dev According to the Aptos docs, the first 0xa addresses are reserved for precompiles.
/// https://github.com/aptos-labs/aptos-core/blob/main/aptos-move/framework/aptos-framework/doc/account.md#function-create_framework_reserved_account-1
/// We use the same range for SUI, even though there is one documented reserved address outside of this range.
/// Since sending a message to this address would not cause any negative side effects, as it would never register
/// a callback with CCIP, there is no negative impact.
/// https://move-book.com/appendix/reserved-addresses.html
const MOVE_PRECOMPILE_SPACE: u256 = 0x0b;

const GAS_PRICE_BITS: u8 = 112;
const GAS_PRICE_MASK_112_BITS: u256 = 0xffffffffffffffffffffffffffff; // 28 f's

const MESSAGE_FIXED_BYTES: u64 = 32 * 15;
const MESSAGE_FIXED_BYTES_PER_TOKEN: u64 = 32 * (4 + (3 + 2));

const CCIP_LOCK_OR_BURN_V1_RET_BYTES: u32 = 32;

/// The maximum number of accounts that can be passed in SVMExtraArgs.
const SVM_EXTRA_ARGS_MAX_ACCOUNTS: u64 = 64;

/// Number of overhead accounts needed for message execution on SVM.
/// These are message.receiver, and the OffRamp Signer PDA specific to the receiver.
const SVM_MESSAGING_ACCOUNTS_OVERHEAD: u64 = 2;

/// The size of each SVM account (in bytes).
const SVM_ACCOUNT_BYTE_SIZE: u64 = 32;

/// The expected static payload size of a token transfer when Borsh encoded and submitted to SVM.
/// TokenPool extra data and offchain data sizes are dynamic, and should be accounted for separately.
const SVM_TOKEN_TRANSFER_DATA_OVERHEAD: u64 =
    (4 + 32) // source_pool
    + 32 // token_address
    + 4 // gas_amount
    + 4 // extra_data overhead
    + 32 // amount
    + 32 // size of the token lookup table account
    + 32 // token-related accounts in the lookup table, over-estimated to 32, typically between 11 - 13
    + 32 // token account belonging to the token receiver, e.g ATA, not included in the token lookup table
    + 32 // per-chain token pool config, not included in the token lookup table
    + 32 // per-chain token billing config, not always included in the token lookup table
    + 32; // OffRamp pool signer PDA, not included in the token lookup table;

const MAX_U64: u256 = 18446744073709551615;
const MAX_U160: u256 = 1461501637330902918203684832716283019655932542975;
const VAL_1E5: u256 = 100_000;
const VAL_1E14: u256 = 100_000_000_000_000;
const VAL_1E16: u256 = 10_000_000_000_000_000;
const VAL_1E18: u256 = 1_000_000_000_000_000_000;

// Link has 8 decimals on Sui and 18 decimals on it's native chain, Ethereum. We want to emit
// the fee in juels (1e18) denomination for consistency across chains. This means we multiply
// the fee by 1e10 on Sui before we emit it in the event.
const LOCAL_8_TO_18_DECIMALS_LINK_MULTIPLIER: u256 = 10_000_000_000;

public struct FeeQuoterState has key, store {
    id: UID,
    max_fee_juels_per_msg: u256,
    link_token: address,
    token_price_staleness_threshold: u64,
    fee_tokens: vector<address>,
    /// @dev The gas price per unit of gas for a given destination chain, in USD with 18 decimals. Multiple gas prices can
    /// be encoded into the same value. Each price takes {Internal.GAS_PRICE_BITS} bits. For example, if Optimism is the
    /// destination chain, gas price can include L1 base fee and L2 gas price. Logic to parse the price components is
    ///  chain-specific, and should live in OnRamp.
    /// @dev Price of 1e18 is 1 USD. Examples:
    ///     Very Expensive:   1 unit of gas costs 1 USD                  -> 1e18.
    ///     Expensive:        1 unit of gas costs 0.1 USD                -> 1e17.
    ///     Cheap:            1 unit of gas costs 0.000001 USD           -> 1e12.
    usd_per_unit_gas_by_dest_chain: table::Table<u64, TimestampedPrice>,
    /// @dev The price, in USD with 18 decimals, per 1e18 of the smallest token denomination.
    /// @dev Price of 1e18 represents 1 USD per 1e18 token amount.
    ///     1 USDC = 1.00 USD per full token, each full token is 1e6 units -> 1 * 1e18 * 1e18 / 1e6 = 1e30.
    ///     1 ETH = 2,000 USD per full token, each full token is 1e18 units -> 2000 * 1e18 * 1e18 / 1e18 = 2_000e18.
    ///     1 LINK = 5.00 USD per full token, each full token is 1e18 units -> 5 * 1e18 * 1e18 / 1e18 = 5e18.
    usd_per_token: table::Table<address, TimestampedPrice>,
    dest_chain_configs: table::Table<u64, DestChainConfig>,
    // dest chain selector -> local token -> TokenTransferFeeConfig
    token_transfer_fee_configs: table::Table<u64, table::Table<address, TokenTransferFeeConfig>>,
    premium_multiplier_wei_per_eth: table::Table<address, u64>,
}

public struct FeeQuoterCap has key, store {
    id: UID,
}

public struct StaticConfig has drop {
    max_fee_juels_per_msg: u256,
    link_token: address,
    token_price_staleness_threshold: u64,
}

public struct DestChainConfig has copy, drop, store {
    is_enabled: bool,
    max_number_of_tokens_per_msg: u16, // Maximum number of distinct tokens transferred per message.
    max_data_bytes: u32, // Maximum data payload size in bytes.
    max_per_msg_gas_limit: u32,
    dest_gas_overhead: u32, // Gas charged on top of the gasLimit to cover destination chain costs.
    dest_gas_per_payload_byte_base: u8, // Default dest-chain gas charged each byte of `data` payload.
    dest_gas_per_payload_byte_high: u8, // High dest-chain gas charged each byte of `data` payload, used to account for eip-7623.
    dest_gas_per_payload_byte_threshold: u16, // The value at which the billing switches from destGasPerPayloadByteBase to destGasPerPayloadByteHigh.
    dest_data_availability_overhead_gas: u32, // Data availability gas charged for overhead costs e.g. for OCR.
    dest_gas_per_data_availability_byte: u16, // Gas units charged per byte of message data that needs availability.
    dest_data_availability_multiplier_bps: u16, // Multiplier for data availability gas, multiples of bps, or 0.0001.
    chain_family_selector: vector<u8>, // Selector that identifies the destination chain's family. Used to determine the correct validations to perform for the dest chain.
    enforce_out_of_order: bool, // Whether to enforce the allowOutOfOrderExecution extraArg value to be true.
    // The following three properties are defaults, they can be overridden by setting the TokenTransferFeeConfig for a token.
    default_token_fee_usd_cents: u16, // Default token fee charged per token transfer.
    default_token_dest_gas_overhead: u32, // Default gas charged to execute a token transfer on the destination chain.
    default_tx_gas_limit: u32, // Default gas limit for a tx.
    gas_multiplier_wei_per_eth: u64, // Multiplier for gas costs, 1e18 based so 11e17 = 10% extra cost.
    gas_price_staleness_threshold: u32, // The amount of time a gas price can be stale before it is considered invalid (0 means disabled).
    network_fee_usd_cents: u32, // Flat network fee to charge for messages, multiples of 0.01 USD.
}

public struct TokenTransferFeeConfig has copy, drop, store {
    min_fee_usd_cents: u32, // Minimum fee to charge per token transfer, multiples of 0.01 USD.
    max_fee_usd_cents: u32, // Maximum fee to charge per token transfer, multiples of 0.01 USD.
    deci_bps: u16, // Basis points charged on token transfers, multiples of 0.1bps, or 1e-5.
    dest_gas_overhead: u32, // Gas charged to execute the token transfer on the destination chain.
    dest_bytes_overhead: u32, // Data availability bytes that are returned from the source pool and sent to the dest pool. Must be >= Pool.CCIP_LOCK_OR_BURN_V1_RET_BYTES. Set as multiple of 32 bytes.
    is_enabled: bool, // Whether this token has custom transfer fees.
}

public struct TimestampedPrice has copy, drop, store {
    value: u256,
    timestamp: u64,
}

public struct FeeTokenAdded has copy, drop {
    fee_token: address,
}

public struct FeeTokenRemoved has copy, drop {
    fee_token: address,
}

public struct TokenTransferFeeConfigAdded has copy, drop {
    dest_chain_selector: u64,
    token: address,
    token_transfer_fee_config: TokenTransferFeeConfig,
}

public struct TokenTransferFeeConfigRemoved has copy, drop {
    dest_chain_selector: u64,
    token: address,
}

public struct UsdPerTokenUpdated has copy, drop {
    token: address,
    usd_per_token: u256,
    timestamp: u64,
}

public struct UsdPerUnitGasUpdated has copy, drop {
    dest_chain_selector: u64,
    usd_per_unit_gas: u256,
    timestamp: u64,
}

public struct DestChainAdded has copy, drop {
    dest_chain_selector: u64,
    dest_chain_config: DestChainConfig,
}

public struct DestChainConfigUpdated has copy, drop {
    dest_chain_selector: u64,
    dest_chain_config: DestChainConfig,
}

public struct PremiumMultiplierWeiPerEthUpdated has copy, drop {
    token: address,
    premium_multiplier_wei_per_eth: u64,
}

const EAlreadyInitialized: u64 = 1;
const EOutOfBound: u64 = 2;
const EUnknownDestChainSelector: u64 = 3;
const EUnknownToken: u64 = 4;
const EDestChainNotEnabled: u64 = 5;
const ETokenUpdateMismatch: u64 = 6;
const EGasUpdateMismatch: u64 = 7;
const ETokenTransferFeeConfigMismatch: u64 = 8;
const EFeeTokenNotSupported: u64 = 9;
const EZeroTokenPrice: u64 = 10;
const EUnknownChainFamilySelector: u64 = 11;
const EStaleGasPrice: u64 = 12;
const EMessageTooLarge: u64 = 13;
const EUnsupportedNumberOfTokens: u64 = 14;
const EInvalidEvmAddress: u64 = 15;
const EInvalid32BytesAddress: u64 = 16;
const EFeeTokenCostTooHigh: u64 = 17;
const EMessageGasLimitTooHigh: u64 = 18;
const EExtraArgOutOfOrderExecutionMustBeTrue: u64 = 19;
const EInvalidExtraArgsTag: u64 = 20;
const EInvalidExtraArgsData: u64 = 21;
const EInvalidTokenReceiver: u64 = 22;
const EMessageComputeUnitLimitTooHigh: u64 = 23;
const EMessageFeeTooHigh: u64 = 24;
const ESourceTokenDataTooLarge: u64 = 25;
const EInvalidDestChainSelector: u64 = 26;
const EInvalidGasLimit: u64 = 27;
const EInvalidChainFamilySelector: u64 = 28;
const EToTokenAmountTooLarge: u64 = 29;
const ETooManySvmExtraArgsAccounts: u64 = 30;
const EInvalidSvmExtraArgsWritableBitmap: u64 = 31;
const EInvalidFeeRange: u64 = 32;
const EInvalidDestBytesOverhead: u64 = 33;
const EInvalidSvmReceiverLength: u64 = 34;
const EInvalidSvmAccountLength: u64 = 35;
const ETokenAmountMismatch: u64 = 36;
const EInvalidOwnerCap: u64 = 37;
const EInvalidFunction: u64 = 38;

public fun type_and_version(): String {
    string::utf8(b"FeeQuoter 1.6.0")
}

#[allow(lint(self_transfer))]
public fun initialize(
    ref: &mut CCIPObjectRef,
    owner_cap: &OwnerCap,
    max_fee_juels_per_msg: u256,
    link_token: address, // can pass in the LINK metadata object to make sure this is a valid token address
    token_price_staleness_threshold: u64,
    fee_tokens: vector<address>,
    ctx: &mut TxContext,
) {
    assert!(!state_object::contains<FeeQuoterState>(ref), EAlreadyInitialized);

    let state = FeeQuoterState {
        id: object::new(ctx),
        max_fee_juels_per_msg,
        link_token,
        token_price_staleness_threshold,
        fee_tokens,
        usd_per_unit_gas_by_dest_chain: table::new(ctx),
        usd_per_token: table::new(ctx),
        dest_chain_configs: table::new(ctx),
        token_transfer_fee_configs: table::new(ctx),
        premium_multiplier_wei_per_eth: table::new(ctx),
    };
    state_object::add(ref, owner_cap, state, ctx);
}

#[allow(lint(self_transfer))]
public fun issue_fee_quoter_cap(_: &OwnerCap, ctx: &mut TxContext) {
    let fee_quoter_cap = FeeQuoterCap {
        id: object::new(ctx),
    };
    transfer::transfer(
        fee_quoter_cap,
        ctx.sender(),
    );
}

public fun get_token_price(ref: &CCIPObjectRef, token: address): TimestampedPrice {
    let state = state_object::borrow<FeeQuoterState>(ref);
    get_token_price_internal(state, token)
}

public fun get_timestamped_price_fields(tp: TimestampedPrice): (u256, u64) {
    (tp.value, tp.timestamp)
}

public fun get_token_prices(
    ref: &CCIPObjectRef,
    tokens: vector<address>,
): (vector<TimestampedPrice>) {
    let state = state_object::borrow<FeeQuoterState>(ref);
    tokens.map_ref!(|token| get_token_price_internal(state, *token))
}

public fun get_dest_chain_gas_price(
    ref: &CCIPObjectRef,
    dest_chain_selector: u64,
): TimestampedPrice {
    let state = state_object::borrow<FeeQuoterState>(ref);
    get_dest_chain_gas_price_internal(state, dest_chain_selector)
}

public fun get_token_and_gas_prices(
    ref: &CCIPObjectRef,
    clock: &clock::Clock,
    token: address,
    dest_chain_selector: u64,
): (u256, u256) {
    let state = state_object::borrow<FeeQuoterState>(ref);
    let dest_chain_config = get_dest_chain_config_internal(
        state,
        dest_chain_selector,
    );
    assert!(dest_chain_config.is_enabled, EDestChainNotEnabled);
    let token_price = get_token_price_internal(state, token);
    let gas_price_value = get_validated_gas_price_internal(
        state,
        clock,
        dest_chain_config,
        dest_chain_selector,
    );
    (token_price.value, gas_price_value)
}

public fun convert_token_amount(
    ref: &CCIPObjectRef,
    from_token: address,
    from_token_amount: u64,
    to_token: address,
): u64 {
    let state = state_object::borrow<FeeQuoterState>(ref);
    convert_token_amount_internal(state, from_token, from_token_amount, to_token)
}

public fun get_fee_tokens(ref: &CCIPObjectRef): vector<address> {
    let state = state_object::borrow<FeeQuoterState>(ref);
    state.fee_tokens
}

public fun apply_fee_token_updates(
    ref: &mut CCIPObjectRef,
    owner_cap: &OwnerCap,
    fee_tokens_to_remove: vector<address>,
    fee_tokens_to_add: vector<address>,
    _ctx: &mut TxContext,
) {
    assert!(object::id(owner_cap) == ref.owner_cap_id(), EInvalidOwnerCap);

    let state = state_object::borrow_mut<FeeQuoterState>(ref);

    // Remove tokens
    fee_tokens_to_remove.do_ref!(|fee_token| {
        let fee_token = *fee_token;
        let (found, index) = state.fee_tokens.index_of(&fee_token);
        if (found) {
            state.fee_tokens.remove(index);
            event::emit(FeeTokenRemoved { fee_token });
        };
    });

    // Add new tokens
    fee_tokens_to_add.do_ref!(|fee_token| {
        let fee_token = *fee_token;
        let (found, _) = state.fee_tokens.index_of(&fee_token);
        if (!found) {
            state.fee_tokens.push_back(fee_token);
            event::emit(FeeTokenAdded { fee_token });
        };
    });
}

public fun get_token_transfer_fee_config(
    ref: &CCIPObjectRef,
    dest_chain_selector: u64,
    token: address,
): TokenTransferFeeConfig {
    let state = state_object::borrow<FeeQuoterState>(ref);
    get_token_transfer_fee_config_internal(
        state,
        dest_chain_selector,
        token,
    )
}

fun get_token_transfer_fee_config_internal(
    state: &FeeQuoterState,
    dest_chain_selector: u64,
    token: address,
): TokenTransferFeeConfig {
    let empty_fee_config = TokenTransferFeeConfig {
        min_fee_usd_cents: 0,
        max_fee_usd_cents: 0,
        deci_bps: 0,
        dest_gas_overhead: 0,
        dest_bytes_overhead: 0,
        is_enabled: false,
    };
    if (!state.token_transfer_fee_configs.contains(dest_chain_selector)) {
        empty_fee_config
    } else {
        let dest_chain_fee_configs = state.token_transfer_fee_configs.borrow(dest_chain_selector);
        if (!dest_chain_fee_configs.contains(token)) {
            empty_fee_config
        } else {
            *dest_chain_fee_configs.borrow(token)
        }
    }
}

// Note that unlike EVM, this only allows changes for a single dest chain selector at a time.
public fun apply_token_transfer_fee_config_updates(
    ref: &mut CCIPObjectRef,
    owner_cap: &OwnerCap,
    dest_chain_selector: u64,
    add_tokens: vector<address>,
    add_min_fee_usd_cents: vector<u32>,
    add_max_fee_usd_cents: vector<u32>,
    add_deci_bps: vector<u16>,
    add_dest_gas_overhead: vector<u32>,
    add_dest_bytes_overhead: vector<u32>,
    add_is_enabled: vector<bool>,
    remove_tokens: vector<address>,
    ctx: &mut TxContext,
) {
    assert!(object::id(owner_cap) == ref.owner_cap_id(), EInvalidOwnerCap);

    let state = state_object::borrow_mut<FeeQuoterState>(ref);

    if (!state.token_transfer_fee_configs.contains(dest_chain_selector)) {
        state
            .token_transfer_fee_configs
            .add(
                dest_chain_selector,
                table::new(ctx),
            );
    };
    let token_transfer_fee_configs = state
        .token_transfer_fee_configs
        .borrow_mut(dest_chain_selector);

    let add_tokens_len = add_tokens.length();
    assert!(add_tokens_len == add_min_fee_usd_cents.length(), ETokenTransferFeeConfigMismatch);
    assert!(add_tokens_len == add_max_fee_usd_cents.length(), ETokenTransferFeeConfigMismatch);
    assert!(add_tokens_len == add_deci_bps.length(), ETokenTransferFeeConfigMismatch);
    assert!(add_tokens_len == add_dest_gas_overhead.length(), ETokenTransferFeeConfigMismatch);
    assert!(add_tokens_len == add_dest_bytes_overhead.length(), ETokenTransferFeeConfigMismatch);
    assert!(add_tokens_len == add_is_enabled.length(), ETokenTransferFeeConfigMismatch);

    let mut i = 0;
    while (i < add_tokens_len) {
        let token = add_tokens[i];
        let min_fee_usd_cents = add_min_fee_usd_cents[i];
        let max_fee_usd_cents = add_max_fee_usd_cents[i];
        let deci_bps = add_deci_bps[i];
        let dest_gas_overhead = add_dest_gas_overhead[i];
        let dest_bytes_overhead = add_dest_bytes_overhead[i];
        let is_enabled = add_is_enabled[i];

        let token_transfer_fee_config = TokenTransferFeeConfig {
            min_fee_usd_cents,
            max_fee_usd_cents,
            deci_bps,
            dest_gas_overhead,
            dest_bytes_overhead,
            is_enabled,
        };

        assert!(
            token_transfer_fee_config.min_fee_usd_cents < token_transfer_fee_config.max_fee_usd_cents,
            EInvalidFeeRange,
        );
        assert!(
            token_transfer_fee_config.dest_bytes_overhead >= CCIP_LOCK_OR_BURN_V1_RET_BYTES,
            EInvalidDestBytesOverhead,
        );

        token_transfer_fee_configs.add(token, token_transfer_fee_config);

        event::emit(TokenTransferFeeConfigAdded {
            dest_chain_selector,
            token,
            token_transfer_fee_config,
        });

        i = i + 1;
    };

    remove_tokens.do_ref!(|token| {
        let token = *token;
        if (token_transfer_fee_configs.contains(token)) {
            token_transfer_fee_configs.remove(token);

            event::emit(TokenTransferFeeConfigRemoved { dest_chain_selector, token });
        }
    });
}

// this should only be called from offramp, hence gated by a fee quoter cap stored in offramp
public fun update_prices(
    ref: &mut CCIPObjectRef,
    _: &FeeQuoterCap,
    clock: &clock::Clock,
    source_tokens: vector<address>,
    source_usd_per_token: vector<u256>,
    gas_dest_chain_selectors: vector<u64>,
    gas_usd_per_unit_gas: vector<u256>,
    _ctx: &mut TxContext,
) {
    assert!(source_tokens.length() == source_usd_per_token.length(), ETokenUpdateMismatch);
    assert!(gas_dest_chain_selectors.length() == gas_usd_per_unit_gas.length(), EGasUpdateMismatch);

    let state = state_object::borrow_mut<FeeQuoterState>(ref);
    let timestamp = clock.timestamp_ms() / 1000;

    source_tokens.zip_do_ref!(&source_usd_per_token, |token, usd_per_token| {
        let timestamped_price = TimestampedPrice {
            value: *usd_per_token,
            timestamp,
        };

        if (state.usd_per_token.contains(*token)) {
            let _old_value = state.usd_per_token.remove(*token);
        };
        state.usd_per_token.add(*token, timestamped_price);

        event::emit(UsdPerTokenUpdated {
            token: *token,
            usd_per_token: *usd_per_token,
            timestamp,
        });
    });

    gas_dest_chain_selectors.zip_do_ref!(
        &gas_usd_per_unit_gas,
        |dest_chain_selector, usd_per_unit_gas| {
            let timestamped_price = TimestampedPrice {
                value: *usd_per_unit_gas,
                timestamp,
            };

            if (state.usd_per_unit_gas_by_dest_chain.contains(*dest_chain_selector)) {
                let _old_value = state.usd_per_unit_gas_by_dest_chain.remove(*dest_chain_selector);
            };
            state.usd_per_unit_gas_by_dest_chain.add(*dest_chain_selector, timestamped_price);

            event::emit(UsdPerUnitGasUpdated {
                dest_chain_selector: *dest_chain_selector,
                usd_per_unit_gas: *usd_per_unit_gas,
                timestamp,
            });
        },
    );
}

public fun get_validated_fee(
    ref: &CCIPObjectRef,
    clock: &clock::Clock,
    dest_chain_selector: u64,
    receiver: vector<u8>,
    data: vector<u8>,
    local_token_addresses: vector<address>, // the token's coin metadata object ids
    local_token_amounts: vector<u64>,
    fee_token: address, // the fee token's coin metadata object id
    extra_args: vector<u8>,
): u64 {
    let state = state_object::borrow<FeeQuoterState>(ref);

    let dest_chain_config = get_dest_chain_config_internal(
        state,
        dest_chain_selector,
    );
    assert!(dest_chain_config.is_enabled, EDestChainNotEnabled);

    assert!(state.fee_tokens.contains(&fee_token), EFeeTokenNotSupported);

    let chain_family_selector = dest_chain_config.chain_family_selector;

    let data_len = data.length();
    let tokens_len = local_token_addresses.length();
    validate_message(dest_chain_config, data_len, tokens_len);

    let gas_limit = if (
        chain_family_selector == CHAIN_FAMILY_SELECTOR_EVM
            || chain_family_selector == CHAIN_FAMILY_SELECTOR_APTOS
            || chain_family_selector == CHAIN_FAMILY_SELECTOR_SUI
    ) {
        resolve_generic_gas_limit(dest_chain_config, extra_args)
    } else if (chain_family_selector == CHAIN_FAMILY_SELECTOR_SVM) {
        resolve_svm_gas_limit(
            dest_chain_config,
            state,
            dest_chain_selector,
            extra_args,
            receiver,
            data_len,
            tokens_len,
            local_token_addresses,
        )
    } else {
        abort EUnknownChainFamilySelector
    };

    validate_dest_family_address(chain_family_selector, receiver, gas_limit);

    let fee_token_price = get_token_price_internal(state, fee_token);
    assert!(fee_token_price.value > 0, EZeroTokenPrice);
    let packed_gas_price = get_validated_gas_price_internal(
        state,
        clock,
        dest_chain_config,
        dest_chain_selector,
    );

    let (mut premium_fee_usd_mysten, token_transfer_gas, token_transfer_bytes_overhead) = if (
        tokens_len > 0
    ) {
        get_token_transfer_cost(
            state,
            dest_chain_config,
            dest_chain_selector,
            fee_token,
            fee_token_price,
            local_token_addresses,
            local_token_amounts,
        )
    } else {
        ((dest_chain_config.network_fee_usd_cents as u256) * VAL_1E16, 0, 0)
    };
    let premium_multiplier = get_premium_multiplier_wei_per_eth_internal(state, fee_token);
    premium_fee_usd_mysten = premium_fee_usd_mysten * (premium_multiplier as u256); // Apply premium multiplier in mysten/sui units

    let data_availability_cost_usd_36_decimals = if (
        dest_chain_config.dest_data_availability_multiplier_bps > 0
    ) {
        // Extract data availability gas price (upper 112 bits) - matches EVM uint112 behavior
        let data_availability_gas_price =
            (packed_gas_price >> GAS_PRICE_BITS) & GAS_PRICE_MASK_112_BITS;
        get_data_availability_cost(
            dest_chain_config,
            data_availability_gas_price,
            data_len,
            tokens_len,
            token_transfer_bytes_overhead,
        )
    } else { 0 };

    let call_data_length: u256 = (data_len as u256) + (token_transfer_bytes_overhead as u256);
    let mut dest_call_data_cost =
        call_data_length * (dest_chain_config.dest_gas_per_payload_byte_base as u256);
    if (call_data_length > (dest_chain_config.dest_gas_per_payload_byte_threshold as u256)) {
        dest_call_data_cost =
            (dest_chain_config.dest_gas_per_payload_byte_base as u256) *
                (dest_chain_config.dest_gas_per_payload_byte_threshold as u256)
           + (call_data_length - (dest_chain_config.dest_gas_per_payload_byte_threshold as u256)) *
                (dest_chain_config.dest_gas_per_payload_byte_high as u256);
    };

    let total_dest_chain_gas =
        (dest_chain_config.dest_gas_overhead as u256) + (token_transfer_gas as u256)
            + dest_call_data_cost + gas_limit;

    let gas_cost = packed_gas_price & GAS_PRICE_MASK_112_BITS;

    let total_cost_usd =
        (total_dest_chain_gas * gas_cost *
            (dest_chain_config.gas_multiplier_wei_per_eth as u256)) +
        premium_fee_usd_mysten + data_availability_cost_usd_36_decimals;

    let fee_token_cost = total_cost_usd / fee_token_price.value;

    // we need to convert back to a u64 which is what the fungible asset module uses for amounts.
    assert!(fee_token_cost <= MAX_U64, EFeeTokenCostTooHigh);
    fee_token_cost as u64
}

public fun apply_premium_multiplier_wei_per_eth_updates(
    ref: &mut CCIPObjectRef,
    owner_cap: &OwnerCap,
    tokens: vector<address>,
    premium_multiplier_wei_per_eth: vector<u64>,
    _ctx: &mut TxContext,
) {
    assert!(object::id(owner_cap) == ref.owner_cap_id(), EInvalidOwnerCap);

    let state = state_object::borrow_mut<FeeQuoterState>(ref);

    tokens.zip_do_ref!(&premium_multiplier_wei_per_eth, |token, premium_multiplier_wei_per_eth| {
        let token: address = *token;
        let premium_multiplier_wei_per_eth: u64 = *premium_multiplier_wei_per_eth;

        if (state.premium_multiplier_wei_per_eth.contains(token)) {
            state.premium_multiplier_wei_per_eth.remove(token);
        };
        state.premium_multiplier_wei_per_eth.add(token, premium_multiplier_wei_per_eth);

        event::emit(PremiumMultiplierWeiPerEthUpdated {
            token,
            premium_multiplier_wei_per_eth,
        });
    });
}

public fun get_premium_multiplier_wei_per_eth(ref: &CCIPObjectRef, token: address): u64 {
    let state = state_object::borrow<FeeQuoterState>(ref);
    get_premium_multiplier_wei_per_eth_internal(state, token)
}

fun get_premium_multiplier_wei_per_eth_internal(state: &FeeQuoterState, token: address): u64 {
    assert!(state.premium_multiplier_wei_per_eth.contains(token), EUnknownToken);
    *state.premium_multiplier_wei_per_eth.borrow(token)
}

fun resolve_generic_gas_limit(dest_chain_config: &DestChainConfig, extra_args: vector<u8>): u256 {
    let (gas_limit, allow_out_of_order_execution) = decode_generic_extra_args(
        dest_chain_config,
        extra_args,
    );
    assert!(
        gas_limit <= (dest_chain_config.max_per_msg_gas_limit as u256),
        EMessageGasLimitTooHigh,
    );
    assert!(
        !dest_chain_config.enforce_out_of_order || allow_out_of_order_execution,
        EExtraArgOutOfOrderExecutionMustBeTrue,
    );
    gas_limit
}

fun resolve_svm_gas_limit(
    dest_chain_config: &DestChainConfig,
    state: &FeeQuoterState,
    dest_chain_selector: u64,
    extra_args: vector<u8>,
    receiver: vector<u8>,
    data_len: u64,
    tokens_len: u64,
    local_token_addresses: vector<address>,
): u256 {
    let extra_args_len = extra_args.length();
    assert!(extra_args_len > 0, EInvalidExtraArgsData);

    let (
        compute_units,
        account_is_writable_bitmap,
        allow_out_of_order_execution,
        token_receiver,
        accounts,
    ) = decode_svm_extra_args(extra_args);

    let gas_limit = compute_units;

    assert!(
        !dest_chain_config.enforce_out_of_order || allow_out_of_order_execution,
        EExtraArgOutOfOrderExecutionMustBeTrue,
    );
    assert!(gas_limit <= dest_chain_config.max_per_msg_gas_limit, EMessageComputeUnitLimitTooHigh);

    let accounts_length = accounts.length();
    // The max payload size for SVM is heavily dependent on the accounts passed into extra args and the number of
    // tokens. Below, token and account overhead will count towards maxDataBytes.
    let mut svm_expanded_data_length = data_len;

    // The receiver length has not yet been validated before this point.
    assert!(receiver.length() == 32, EInvalidSvmReceiverLength);
    let receiver_uint = eth_abi::decode_u256_value(receiver);
    if (receiver_uint == 0) {
        // When message receiver is zero, CCIP receiver is not invoked on SVM.
        // There should not be additional accounts specified for the receiver.
        assert!(accounts_length == 0, ETooManySvmExtraArgsAccounts);
    } else {
        // The messaging accounts needed for CCIP receiver on SVM are:
        // message receiver, offramp PDA signer,
        // plus remaining accounts specified in SVM extraArgs. Each account is 32 bytes.
        svm_expanded_data_length =
            svm_expanded_data_length + (
            accounts_length + SVM_MESSAGING_ACCOUNTS_OVERHEAD
        ) * SVM_ACCOUNT_BYTE_SIZE;
    };

    let mut i = 0;
    while (i < accounts_length) {
        assert!(accounts[i].length() == 32, EInvalidSvmAccountLength);
        i = i + 1;
    };

    if (tokens_len > 0) {
        assert!(
            token_receiver.length() == 32
                    && eth_abi::decode_u256_value(token_receiver) != 0,
            EInvalidTokenReceiver,
        );
    };

    assert!(accounts_length <= SVM_EXTRA_ARGS_MAX_ACCOUNTS, ETooManySvmExtraArgsAccounts);
    assert!(
        (account_is_writable_bitmap >> (accounts_length as u8)) == 0,
        EInvalidSvmExtraArgsWritableBitmap,
    );

    svm_expanded_data_length =
        svm_expanded_data_length + tokens_len * SVM_TOKEN_TRANSFER_DATA_OVERHEAD;

    // The token destBytesOverhead can be very different per token so we have to take it into account as well.
    let mut i = 0;
    while (i < tokens_len) {
        let local_token_address = local_token_addresses[i];
        let destBytesOverhead = get_token_transfer_fee_config_internal(
            state,
            dest_chain_selector,
            local_token_address,
        ).dest_bytes_overhead;

        // Pools get CCIP_LOCK_OR_BURN_V1_RET_BYTES by default, but if an override is set we use that instead.
        if (destBytesOverhead > 0) {
            svm_expanded_data_length = svm_expanded_data_length + (destBytesOverhead as u64);
        } else {
            svm_expanded_data_length =
                svm_expanded_data_length + (CCIP_LOCK_OR_BURN_V1_RET_BYTES as u64);
        };

        i = i + 1;
    };

    assert!(
        svm_expanded_data_length <= (dest_chain_config.max_data_bytes as u64),
        EMessageTooLarge,
    );

    gas_limit as u256
}

fun decode_generic_extra_args(
    dest_chain_config: &DestChainConfig,
    extra_args: vector<u8>,
): (u256, bool) {
    let extra_args_len = extra_args.length();
    if (extra_args_len == 0) {
        // If extra args are empty, generate default values.
        (dest_chain_config.default_tx_gas_limit as u256, false)
    } else {
        assert!(extra_args_len >= 4, EInvalidExtraArgsData);

        let args_tag = slice(&extra_args, 0, 4);
        assert!(args_tag == client::generic_extra_args_v2_tag(), EInvalidExtraArgsTag);

        let args_data = slice(&extra_args, 4, extra_args_len - 4);
        decode_generic_extra_args_v2(args_data)
    }
}

fun decode_generic_extra_args_v2(extra_args: vector<u8>): (u256, bool) {
    let mut stream = bcs_stream::new(extra_args);
    let gas_limit = bcs_stream::deserialize_u256(&mut stream);
    let allow_out_of_order_execution = bcs_stream::deserialize_bool(&mut stream);
    bcs_stream::assert_is_consumed(&stream);
    (gas_limit, allow_out_of_order_execution)
}

fun decode_svm_extra_args(
    extra_args: vector<u8>,
): (u32, u64, bool, vector<u8>, vector<vector<u8>>) {
    let extra_args_len = extra_args.length();
    let args_tag = slice(&extra_args, 0, 4);
    assert!(args_tag == client::svm_extra_args_v1_tag(), EInvalidExtraArgsTag);
    assert!(extra_args_len >= 4, EInvalidExtraArgsData);
    let args_data = slice(&extra_args, 4, extra_args_len - 4);
    decode_svm_extra_args_v1(args_data)
}

fun decode_svm_extra_args_v1(
    extra_args: vector<u8>,
): (u32, u64, bool, vector<u8>, vector<vector<u8>>) {
    let mut stream = bcs_stream::new(extra_args);
    let compute_units = bcs_stream::deserialize_u32(&mut stream);
    let account_is_writable_bitmap = bcs_stream::deserialize_u64(&mut stream);
    let allow_out_of_order_execution = bcs_stream::deserialize_bool(&mut stream);
    let token_receiver = bcs_stream::deserialize_vector_u8(&mut stream);
    let accounts = bcs_stream::deserialize_vector!(
        &mut stream,
        |stream| bcs_stream::deserialize_vector_u8(stream),
    );
    bcs_stream::assert_is_consumed(&stream);
    (
        compute_units,
        account_is_writable_bitmap,
        allow_out_of_order_execution,
        token_receiver,
        accounts,
    )
}

fun get_data_availability_cost(
    dest_chain_config: &DestChainConfig,
    data_availability_gas_price: u256,
    data_len: u64,
    tokens_len: u64,
    total_transfer_bytes_overhead: u32,
): u256 {
    let data_availability_length_bytes =
        MESSAGE_FIXED_BYTES + data_len + (tokens_len
            * MESSAGE_FIXED_BYTES_PER_TOKEN)
            + (total_transfer_bytes_overhead as u64);

    let data_availability_gas =
        ((data_availability_length_bytes as u256)
            * (dest_chain_config.dest_gas_per_data_availability_byte as u256)) + (
            dest_chain_config.dest_data_availability_overhead_gas as u256
        );

    data_availability_gas * data_availability_gas_price
        * (dest_chain_config.dest_data_availability_multiplier_bps as u256)
        * VAL_1E14
}

fun get_token_transfer_cost(
    state: &FeeQuoterState,
    dest_chain_config: &DestChainConfig,
    dest_chain_selector: u64,
    fee_token: address,
    fee_token_price: TimestampedPrice,
    local_token_addresses: vector<address>,
    local_token_amounts: vector<u64>,
): (u256, u32, u32) {
    let mut token_transfer_fee_wei: u256 = 0;
    let mut token_transfer_gas: u32 = 0;
    let mut token_transfer_bytes_overhead: u32 = 0;

    local_token_addresses.zip_do_ref!(
        &local_token_amounts,
        |local_token_address, local_token_amount| {
            let local_token_address: address = *local_token_address;
            let local_token_amount: u64 = *local_token_amount;

            let transfer_fee_config = get_token_transfer_fee_config_internal(
                state,
                dest_chain_selector,
                local_token_address,
            );

            if (!transfer_fee_config.is_enabled) {
                token_transfer_fee_wei =
                    token_transfer_fee_wei + (dest_chain_config.default_token_fee_usd_cents as u256) * VAL_1E16;
                token_transfer_gas =
                    token_transfer_gas + dest_chain_config.default_token_dest_gas_overhead;
                token_transfer_bytes_overhead =
                    token_transfer_bytes_overhead + CCIP_LOCK_OR_BURN_V1_RET_BYTES;
            } else {
                let mut bps_fee_usd_wei = 0;
                if (transfer_fee_config.deci_bps > 0) {
                    let token_price = if (local_token_address == fee_token) {
                        fee_token_price
                    } else {
                        get_token_price_internal(state, local_token_address)
                    };
                    let token_usd_value = calc_usd_value_from_token_amount(
                        local_token_amount,
                        token_price.value,
                    );
                    bps_fee_usd_wei =
                        token_usd_value * (transfer_fee_config.deci_bps as u256) / VAL_1E5;
                };

                token_transfer_gas = token_transfer_gas + transfer_fee_config.dest_gas_overhead;
                token_transfer_bytes_overhead =
                    token_transfer_bytes_overhead + transfer_fee_config.dest_bytes_overhead;

                let min_fee_usd_wei = (transfer_fee_config.min_fee_usd_cents as u256) * VAL_1E16;
                let max_fee_usd_wei = (transfer_fee_config.max_fee_usd_cents as u256) * VAL_1E16;
                let selected_fee_usd_wei = if (bps_fee_usd_wei < min_fee_usd_wei) {
                    min_fee_usd_wei
                } else if (bps_fee_usd_wei > max_fee_usd_wei) {
                    max_fee_usd_wei
                } else {
                    bps_fee_usd_wei
                };
                token_transfer_fee_wei = token_transfer_fee_wei + selected_fee_usd_wei;
            }
        },
    );

    (token_transfer_fee_wei, token_transfer_gas, token_transfer_bytes_overhead)
}

fun calc_usd_value_from_token_amount(token_amount: u64, token_price: u256): u256 {
    (token_amount as u256) * token_price / VAL_1E18
}

public fun get_token_receiver(
    ref: &CCIPObjectRef,
    dest_chain_selector: u64,
    extra_args: vector<u8>,
    message_receiver: vector<u8>,
): vector<u8> {
    let state = state_object::borrow<FeeQuoterState>(ref);

    let chain_family_selector = get_dest_chain_config_internal(
        state,
        dest_chain_selector,
    ).chain_family_selector;
    if (
        chain_family_selector == CHAIN_FAMILY_SELECTOR_EVM
        || chain_family_selector == CHAIN_FAMILY_SELECTOR_APTOS
        || chain_family_selector == CHAIN_FAMILY_SELECTOR_SUI
    ) {
        message_receiver
    } else if (chain_family_selector == CHAIN_FAMILY_SELECTOR_SVM) {
        let (
            _compute_units,
            _account_is_writable_bitmap,
            _allow_out_of_order_execution,
            token_receiver,
            _accounts,
        ) = decode_svm_extra_args(extra_args);
        token_receiver
    } else {
        abort EUnknownChainFamilySelector
    }
}

/// @returns (msg_fee_juels, is_out_of_order_execution, converted_extra_args, dest_exec_data_per_token)
public fun process_message_args(
    ref: &CCIPObjectRef,
    dest_chain_selector: u64,
    fee_token: address,
    fee_token_amount: u64,
    extra_args: vector<u8>,
    local_token_addresses: vector<address>,
    dest_token_addresses: vector<vector<u8>>,
    dest_pool_datas: vector<vector<u8>>,
): (u256, bool, vector<u8>, vector<vector<u8>>) {
    let state = state_object::borrow<FeeQuoterState>(ref);
    // This is the fee in Sui denomination. We convert it to juels (1e18 based) below.
    let msg_fee_link_local_denomination = if (fee_token == state.link_token) {
        fee_token_amount
    } else {
        convert_token_amount_internal(
            state,
            fee_token,
            fee_token_amount,
            state.link_token,
        )
    };

    // We convert the local denomination to juels here. This means that the offchain monitoring will always
    // get a consistent juels amount regardless of the token denomination on the chain.
    let msg_fee_juels =
        (msg_fee_link_local_denomination as u256)
            * LOCAL_8_TO_18_DECIMALS_LINK_MULTIPLIER;

    // max_fee_juels_per_msg is in juels denomination for consistency across chains.
    assert!(msg_fee_juels <= state.max_fee_juels_per_msg, EMessageFeeTooHigh);

    let dest_chain_config = get_dest_chain_config_internal(
        state,
        dest_chain_selector,
    );

    let (converted_extra_args, is_out_of_order_execution) = process_chain_family_selector(
        dest_chain_config,
        !dest_token_addresses.is_empty(),
        extra_args,
    );

    let dest_exec_data_per_token = process_pool_return_data(
        state,
        dest_chain_config,
        dest_chain_selector,
        local_token_addresses,
        dest_token_addresses,
        dest_pool_datas,
    );

    (msg_fee_juels, is_out_of_order_execution, converted_extra_args, dest_exec_data_per_token)
}

fun process_chain_family_selector(
    dest_chain_config: &DestChainConfig,
    is_message_with_token_transfers: bool,
    extra_args: vector<u8>,
): (vector<u8>, bool) {
    let chain_family_selector = dest_chain_config.chain_family_selector;
    if (
        chain_family_selector == CHAIN_FAMILY_SELECTOR_EVM
        || chain_family_selector == CHAIN_FAMILY_SELECTOR_APTOS
        || chain_family_selector == CHAIN_FAMILY_SELECTOR_SUI
    ) {
        let (gas_limit, allow_out_of_order_execution) = decode_generic_extra_args(
            dest_chain_config,
            extra_args,
        );
        let extra_args_v2 = client::encode_generic_extra_args_v2(
            gas_limit,
            allow_out_of_order_execution,
        );
        (extra_args_v2, allow_out_of_order_execution)
    } else if (chain_family_selector == CHAIN_FAMILY_SELECTOR_SVM) {
        let (
            compute_units,
            _account_is_writable_bitmap,
            allow_out_of_order_execution,
            token_receiver,
            _accounts,
        ) = decode_svm_extra_args(extra_args);
        if (is_message_with_token_transfers) {
            assert!(token_receiver.length() == 32, EInvalidTokenReceiver);
            let token_receiver_uint = eth_abi::decode_u256_value(token_receiver);
            assert!(token_receiver_uint > 0, EInvalidTokenReceiver);
        };

        assert!(
            !dest_chain_config.enforce_out_of_order || allow_out_of_order_execution,
            EExtraArgOutOfOrderExecutionMustBeTrue,
        );
        assert!(
            compute_units <= dest_chain_config.max_per_msg_gas_limit,
            EMessageComputeUnitLimitTooHigh,
        );

        (extra_args, allow_out_of_order_execution)
    } else {
        abort EUnknownChainFamilySelector
    }
}

fun process_pool_return_data(
    state: &FeeQuoterState,
    dest_chain_config: &DestChainConfig,
    dest_chain_selector: u64,
    local_token_addresses: vector<address>,
    dest_token_addresses: vector<vector<u8>>,
    dest_pool_datas: vector<vector<u8>>,
): vector<vector<u8>> {
    let chain_family_selector = dest_chain_config.chain_family_selector;

    let tokens_len = dest_token_addresses.length();
    assert!(tokens_len == dest_pool_datas.length(), ETokenAmountMismatch);

    let mut dest_exec_data_per_token = vector[];
    let mut i = 0;
    while (i < tokens_len) {
        let local_token_address = local_token_addresses[i];
        let dest_token_address = dest_token_addresses[i];
        let dest_pool_data_len = dest_pool_datas[i].length();

        let token_transfer_fee_config = get_token_transfer_fee_config_internal(
            state,
            dest_chain_selector,
            local_token_address,
        );

        if (dest_pool_data_len > (CCIP_LOCK_OR_BURN_V1_RET_BYTES as u64)) {
            assert!(
                dest_pool_data_len <= (token_transfer_fee_config.dest_bytes_overhead as u64),
                ESourceTokenDataTooLarge,
            );
        };

        // We pass in 1 as gas_limit as this only matters for SVM address validation. This ensures the address
        // may not be 0x0.
        validate_dest_family_address(chain_family_selector, dest_token_address, 1);

        let dest_gas_amount = if (token_transfer_fee_config.is_enabled) {
            token_transfer_fee_config.dest_gas_overhead
        } else {
            dest_chain_config.default_token_dest_gas_overhead
        };

        let dest_exec_data = bcs::to_bytes(&dest_gas_amount);
        dest_exec_data_per_token.push_back(dest_exec_data);

        i = i + 1;
    };

    dest_exec_data_per_token
}

public fun get_dest_chain_config(ref: &CCIPObjectRef, dest_chain_selector: u64): DestChainConfig {
    let state = state_object::borrow<FeeQuoterState>(ref);
    *get_dest_chain_config_internal(state, dest_chain_selector)
}

fun get_dest_chain_config_internal(
    state: &FeeQuoterState,
    dest_chain_selector: u64,
): &DestChainConfig {
    assert!(state.dest_chain_configs.contains(dest_chain_selector), EUnknownDestChainSelector);
    state.dest_chain_configs.borrow(dest_chain_selector)
}

public fun get_dest_chain_config_fields(
    dest_chain_config: DestChainConfig,
): (
    bool,
    u16,
    u32,
    u32,
    u32,
    u8,
    u8,
    u16,
    u32,
    u16,
    u16,
    vector<u8>,
    bool,
    u16,
    u32,
    u32,
    u64,
    u32,
    u32,
) {
    (
        dest_chain_config.is_enabled,
        dest_chain_config.max_number_of_tokens_per_msg,
        dest_chain_config.max_data_bytes,
        dest_chain_config.max_per_msg_gas_limit,
        dest_chain_config.dest_gas_overhead,
        dest_chain_config.dest_gas_per_payload_byte_base,
        dest_chain_config.dest_gas_per_payload_byte_high,
        dest_chain_config.dest_gas_per_payload_byte_threshold,
        dest_chain_config.dest_data_availability_overhead_gas,
        dest_chain_config.dest_gas_per_data_availability_byte,
        dest_chain_config.dest_data_availability_multiplier_bps,
        dest_chain_config.chain_family_selector,
        dest_chain_config.enforce_out_of_order,
        dest_chain_config.default_token_fee_usd_cents,
        dest_chain_config.default_token_dest_gas_overhead,
        dest_chain_config.default_tx_gas_limit,
        dest_chain_config.gas_multiplier_wei_per_eth,
        dest_chain_config.gas_price_staleness_threshold,
        dest_chain_config.network_fee_usd_cents,
    )
}

public fun apply_dest_chain_config_updates(
    ref: &mut CCIPObjectRef,
    owner_cap: &OwnerCap,
    dest_chain_selector: u64,
    is_enabled: bool,
    max_number_of_tokens_per_msg: u16,
    max_data_bytes: u32,
    max_per_msg_gas_limit: u32,
    dest_gas_overhead: u32,
    dest_gas_per_payload_byte_base: u8,
    dest_gas_per_payload_byte_high: u8,
    dest_gas_per_payload_byte_threshold: u16,
    dest_data_availability_overhead_gas: u32,
    dest_gas_per_data_availability_byte: u16,
    dest_data_availability_multiplier_bps: u16,
    chain_family_selector: vector<u8>,
    enforce_out_of_order: bool,
    default_token_fee_usd_cents: u16,
    default_token_dest_gas_overhead: u32,
    default_tx_gas_limit: u32,
    gas_multiplier_wei_per_eth: u64,
    gas_price_staleness_threshold: u32,
    network_fee_usd_cents: u32,
    _ctx: &mut TxContext,
) {
    assert!(object::id(owner_cap) == ref.owner_cap_id(), EInvalidOwnerCap);

    let state = state_object::borrow_mut<FeeQuoterState>(ref);

    assert!(dest_chain_selector != 0, EInvalidDestChainSelector);
    assert!(
        default_tx_gas_limit != 0 && default_tx_gas_limit <= max_per_msg_gas_limit,
        EInvalidGasLimit,
    );

    assert!(
        chain_family_selector == CHAIN_FAMILY_SELECTOR_EVM
            || chain_family_selector == CHAIN_FAMILY_SELECTOR_SVM
            || chain_family_selector == CHAIN_FAMILY_SELECTOR_APTOS
            || chain_family_selector == CHAIN_FAMILY_SELECTOR_SUI,
        EInvalidChainFamilySelector,
    );

    let dest_chain_config = DestChainConfig {
        is_enabled,
        max_number_of_tokens_per_msg,
        max_data_bytes,
        max_per_msg_gas_limit,
        dest_gas_overhead,
        dest_gas_per_payload_byte_base,
        dest_gas_per_payload_byte_high,
        dest_gas_per_payload_byte_threshold,
        dest_data_availability_overhead_gas,
        dest_gas_per_data_availability_byte,
        dest_data_availability_multiplier_bps,
        chain_family_selector,
        enforce_out_of_order,
        default_token_fee_usd_cents,
        default_token_dest_gas_overhead,
        default_tx_gas_limit,
        gas_multiplier_wei_per_eth,
        gas_price_staleness_threshold,
        network_fee_usd_cents,
    };

    if (state.dest_chain_configs.contains(dest_chain_selector)) {
        let dest_chain_config_ref = state.dest_chain_configs.borrow_mut(dest_chain_selector);
        *dest_chain_config_ref = dest_chain_config;
        event::emit(DestChainConfigUpdated { dest_chain_selector, dest_chain_config });
    } else {
        state.dest_chain_configs.add(dest_chain_selector, dest_chain_config);
        event::emit(DestChainAdded { dest_chain_selector, dest_chain_config });
    }
}

public fun get_static_config(ref: &CCIPObjectRef): StaticConfig {
    let state = state_object::borrow<FeeQuoterState>(ref);
    StaticConfig {
        max_fee_juels_per_msg: state.max_fee_juels_per_msg,
        link_token: state.link_token,
        token_price_staleness_threshold: state.token_price_staleness_threshold,
    }
}

public fun get_static_config_fields(cfg: StaticConfig): (u256, address, u64) {
    (cfg.max_fee_juels_per_msg, cfg.link_token, cfg.token_price_staleness_threshold)
}

fun get_validated_token_price(state: &FeeQuoterState, token: address): TimestampedPrice {
    let token_price = get_token_price_internal(state, token);
    assert!(token_price.value > 0 && token_price.timestamp > 0, EUnknownToken);
    token_price
}

// Token prices can be stale. On EVM we have additional fallbacks to a price feed, if configured.
// Since these fallbacks don't exist on Sui, we simply return the price as is.
fun get_token_price_internal(state: &FeeQuoterState, token: address): TimestampedPrice {
    assert!(state.usd_per_token.contains(token), EUnknownToken);
    *state.usd_per_token.borrow(token)
}

fun get_dest_chain_gas_price_internal(
    state: &FeeQuoterState,
    dest_chain_selector: u64,
): TimestampedPrice {
    assert!(
        state.usd_per_unit_gas_by_dest_chain.contains(dest_chain_selector),
        EUnknownDestChainSelector,
    );
    *state.usd_per_unit_gas_by_dest_chain.borrow(dest_chain_selector)
}

fun get_validated_gas_price_internal(
    state: &FeeQuoterState,
    clock: &clock::Clock,
    dest_chain_config: &DestChainConfig,
    dest_chain_selector: u64,
): u256 {
    let gas_price = get_dest_chain_gas_price_internal(state, dest_chain_selector);
    if (dest_chain_config.gas_price_staleness_threshold > 0) {
        let time_passed_secs = clock::timestamp_ms(clock) / 1000 - gas_price.timestamp;
        assert!(
            time_passed_secs <= (dest_chain_config.gas_price_staleness_threshold as u64),
            EStaleGasPrice,
        );
    };
    gas_price.value
}

fun convert_token_amount_internal(
    state: &FeeQuoterState,
    from_token: address,
    from_token_amount: u64,
    to_token: address,
): u64 {
    let from_token_price = get_validated_token_price(state, from_token);
    let to_token_price = get_validated_token_price(state, to_token);

    let to_token_amount =
        (from_token_amount as u256) * from_token_price.value / to_token_price.value;
    assert!(to_token_amount <= MAX_U64, EToTokenAmountTooLarge);
    to_token_amount as u64
}

fun validate_message(dest_chain_config: &DestChainConfig, data_len: u64, tokens_len: u64) {
    assert!(data_len <= (dest_chain_config.max_data_bytes as u64), EMessageTooLarge);
    assert!(
        tokens_len <= (dest_chain_config.max_number_of_tokens_per_msg as u64),
        EUnsupportedNumberOfTokens,
    );
}

fun validate_dest_family_address(
    chain_family_selector: vector<u8>,
    encoded_address: vector<u8>,
    gas_limit: u256,
) {
    if (chain_family_selector == CHAIN_FAMILY_SELECTOR_EVM) {
        validate_evm_address(encoded_address);
    } else if (chain_family_selector == CHAIN_FAMILY_SELECTOR_SVM) {
        // SVM addresses don't have a precompile space at the first X addresses, instead we validate that if the gasLimit
        // is non-zero, the address must not be 0x0.
        let mut min_address = 0;
        if (gas_limit > 0) {
            min_address = 1;
        };
        validate_32byte_address(encoded_address, min_address);
    } else if (chain_family_selector == CHAIN_FAMILY_SELECTOR_APTOS
        || chain_family_selector == CHAIN_FAMILY_SELECTOR_SUI) {
        validate_32byte_address(encoded_address, MOVE_PRECOMPILE_SPACE);
    };
}

fun validate_evm_address(encoded_address: vector<u8>) {
    let encoded_address_len = encoded_address.length();
    assert!(encoded_address_len == 32, EInvalid32BytesAddress);

    let encoded_address_uint = eth_abi::decode_u256_value(encoded_address);

    assert!(encoded_address_uint >= EVM_PRECOMPILE_SPACE, EInvalidEvmAddress);
    assert!(encoded_address_uint <= MAX_U160, EInvalidEvmAddress);
}

fun validate_32byte_address(encoded_address: vector<u8>, min_value: u256) {
    assert!(encoded_address.length() == 32, EInvalid32BytesAddress);

    let encoded_address_uint = eth_abi::decode_u256_value(encoded_address);
    assert!(encoded_address_uint >= min_value, EInvalid32BytesAddress);
}

public fun get_token_transfer_fee_config_fields(
    cfg: TokenTransferFeeConfig,
): (u32, u32, u16, u32, u32, bool) {
    (
        cfg.min_fee_usd_cents,
        cfg.max_fee_usd_cents,
        cfg.deci_bps,
        cfg.dest_gas_overhead,
        cfg.dest_bytes_overhead,
        cfg.is_enabled,
    )
}

/// Returns a new vector containing `len` elements from `vec`
/// starting at index `start`. Panics if `start + len` exceeds the vector length.
fun slice<T: copy>(vec: &vector<T>, start: u64, len: u64): vector<T> {
    let vec_len = vec.length();
    // Ensure we have enough elements for the slice.
    assert!(start + len <= vec_len, EOutOfBound);
    let mut new_vec = vector::empty<T>();
    let mut i = start;
    while (i < start + len) {
        // Copy each element from the original vector into the new vector.
        new_vec.push_back(vec[i]);
        i = i + 1;
    };
    new_vec
}

// ================================================================
// |                      MCMS Entrypoint                         |
// ================================================================

/// Proof for CCIP admin
public struct CCIPAdminProof has drop {}

public struct McmsCallback has drop {}

public fun mcms_entrypoint(
    ref: &mut CCIPObjectRef,
    registry: &mut Registry,
    params: ExecutingCallbackParams,
    ctx: &mut TxContext,
) {
    let (owner_cap, function, data) = mcms_registry::get_callback_params<McmsCallback, OwnerCap>(
        registry,
        McmsCallback {},
        params,
    );

    let function_bytes = *function.as_bytes();
    let mut stream = bcs_stream::new(data);

    if (function_bytes == b"apply_fee_token_updates") {
        let fee_tokens_to_remove = bcs_stream::deserialize_vector!(
            &mut stream,
            |stream| { bcs_stream::deserialize_address(stream) },
        );
        let fee_tokens_to_add = bcs_stream::deserialize_vector!(
            &mut stream,
            |stream| { bcs_stream::deserialize_address(stream) },
        );
        bcs_stream::assert_is_consumed(&stream);

        apply_fee_token_updates(ref, owner_cap, fee_tokens_to_remove, fee_tokens_to_add, ctx);
    } else if (function_bytes == b"apply_dest_chain_config_updates") {
        let dest_chain_selector = bcs_stream::deserialize_u64(&mut stream);
        let is_enabled = bcs_stream::deserialize_bool(&mut stream);
        let max_number_of_tokens_per_msg = bcs_stream::deserialize_u16(&mut stream);
        let max_data_bytes = bcs_stream::deserialize_u32(&mut stream);
        let max_per_msg_gas_limit = bcs_stream::deserialize_u32(&mut stream);
        let dest_gas_overhead = bcs_stream::deserialize_u32(&mut stream);
        let dest_gas_per_payload_byte_base = bcs_stream::deserialize_u8(&mut stream);
        let dest_gas_per_payload_byte_high = bcs_stream::deserialize_u8(&mut stream);
        let dest_gas_per_payload_byte_threshold = bcs_stream::deserialize_u16(&mut stream);
        let dest_data_availability_overhead_gas = bcs_stream::deserialize_u32(&mut stream);
        let dest_gas_per_data_availability_byte = bcs_stream::deserialize_u16(&mut stream);
        let dest_data_availability_multiplier_bps = bcs_stream::deserialize_u16(&mut stream);
        let chain_family_selector = bcs_stream::deserialize_vector_u8(&mut stream);
        let enforce_out_of_order = bcs_stream::deserialize_bool(&mut stream);
        let default_token_fee_usd_cents = bcs_stream::deserialize_u16(&mut stream);
        let default_token_dest_gas_overhead = bcs_stream::deserialize_u32(&mut stream);
        let default_tx_gas_limit = bcs_stream::deserialize_u32(&mut stream);
        let gas_multiplier_wei_per_eth = bcs_stream::deserialize_u64(&mut stream);
        let gas_price_staleness_threshold = bcs_stream::deserialize_u32(&mut stream);
        let network_fee_usd_cents = bcs_stream::deserialize_u32(&mut stream);
        bcs_stream::assert_is_consumed(&stream);

        apply_dest_chain_config_updates(
            ref,
            owner_cap,
            dest_chain_selector,
            is_enabled,
            max_number_of_tokens_per_msg,
            max_data_bytes,
            max_per_msg_gas_limit,
            dest_gas_overhead,
            dest_gas_per_payload_byte_base,
            dest_gas_per_payload_byte_high,
            dest_gas_per_payload_byte_threshold,
            dest_data_availability_overhead_gas,
            dest_gas_per_data_availability_byte,
            dest_data_availability_multiplier_bps,
            chain_family_selector,
            enforce_out_of_order,
            default_token_fee_usd_cents,
            default_token_dest_gas_overhead,
            default_tx_gas_limit,
            gas_multiplier_wei_per_eth,
            gas_price_staleness_threshold,
            network_fee_usd_cents,
            ctx,
        );
    } else if (function_bytes == b"apply_token_transfer_fee_config_updates") {
        let dest_chain_selector = bcs_stream::deserialize_u64(&mut stream);
        let add_tokens = bcs_stream::deserialize_vector!(
            &mut stream,
            |stream| { bcs_stream::deserialize_address(stream) },
        );
        let add_min_fee_usd_cents = bcs_stream::deserialize_vector!(
            &mut stream,
            |stream| { bcs_stream::deserialize_u32(stream) },
        );
        let add_max_fee_usd_cents = bcs_stream::deserialize_vector!(
            &mut stream,
            |stream| { bcs_stream::deserialize_u32(stream) },
        );
        let add_deci_bps = bcs_stream::deserialize_vector!(
            &mut stream,
            |stream| { bcs_stream::deserialize_u16(stream) },
        );
        let add_dest_gas_overhead = bcs_stream::deserialize_vector!(
            &mut stream,
            |stream| { bcs_stream::deserialize_u32(stream) },
        );
        let add_dest_bytes_overhead = bcs_stream::deserialize_vector!(
            &mut stream,
            |stream| { bcs_stream::deserialize_u32(stream) },
        );
        let add_is_enabled = bcs_stream::deserialize_vector!(
            &mut stream,
            |stream| { bcs_stream::deserialize_bool(stream) },
        );
        let remove_tokens = bcs_stream::deserialize_vector!(
            &mut stream,
            |stream| { bcs_stream::deserialize_address(stream) },
        );
        bcs_stream::assert_is_consumed(&stream);

        apply_token_transfer_fee_config_updates(
            ref,
            owner_cap,
            dest_chain_selector,
            add_tokens,
            add_min_fee_usd_cents,
            add_max_fee_usd_cents,
            add_deci_bps,
            add_dest_gas_overhead,
            add_dest_bytes_overhead,
            add_is_enabled,
            remove_tokens,
            ctx,
        );
    } else if (function_bytes == b"apply_premium_multiplier_wei_per_eth_updates") {
        let tokens = bcs_stream::deserialize_vector!(
            &mut stream,
            |stream| { bcs_stream::deserialize_address(stream) },
        );
        let premium_multiplier_wei_per_eth = bcs_stream::deserialize_vector!(
            &mut stream,
            |stream| { bcs_stream::deserialize_u64(stream) },
        );
        bcs_stream::assert_is_consumed(&stream);

        apply_premium_multiplier_wei_per_eth_updates(
            ref,
            owner_cap,
            tokens,
            premium_multiplier_wei_per_eth,
            ctx,
        );
    } else if (function_bytes == b"issue_fee_quoter_cap") {
        issue_fee_quoter_cap(owner_cap, ctx);
    } else {
        abort EInvalidFunction
    }
}

#[test_only]
public fun create_fee_quoter_cap(ctx: &mut TxContext): FeeQuoterCap {
    FeeQuoterCap {
        id: object::new(ctx),
    }
}

#[test_only]
public fun destroy_fee_quoter_cap(cap: FeeQuoterCap) {
    let FeeQuoterCap { id } = cap;
    object::delete(id);
}

#[test]
fun test_decode_generic_extra_args_v2() {
    let dest_chain_config = DestChainConfig {
        is_enabled: true,
        max_number_of_tokens_per_msg: 1000,
        max_data_bytes: 1000,
        max_per_msg_gas_limit: 1000,
        dest_gas_overhead: 1000,
        dest_gas_per_payload_byte_base: 10,
        dest_gas_per_payload_byte_high: 10,
        dest_gas_per_payload_byte_threshold: 1000,
        dest_data_availability_overhead_gas: 1000,
        dest_gas_per_data_availability_byte: 1000,
        dest_data_availability_multiplier_bps: 1000,
        chain_family_selector: b"test",
        enforce_out_of_order: true,
        default_token_fee_usd_cents: 1000,
        default_token_dest_gas_overhead: 1000,
        default_tx_gas_limit: 1000,
        gas_multiplier_wei_per_eth: 1000,
        gas_price_staleness_threshold: 1000,
        network_fee_usd_cents: 1000,
    };

    let expected_gas_limit = 101;
    let expected_allow_out_of_order_execution = true;

    let extra_args = client::encode_generic_extra_args_v2(
        expected_gas_limit,
        expected_allow_out_of_order_execution,
    );

    let (gas_limit, allow_out_of_order_execution) = decode_generic_extra_args(
        &dest_chain_config,
        extra_args,
    );

    assert!(gas_limit == expected_gas_limit, 0);
    assert!(allow_out_of_order_execution == expected_allow_out_of_order_execution, 0);
}

#[test]
fun test_decode_svm_extra_args_v1() {
    let expected_compute_units = 101;
    let expected_account_is_writable_bitmap = 102;
    let expected_allow_out_of_order_execution = true;
    let expected_token_receiver =
        x"1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef";
    let expected_accounts = vector[
        x"2234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdea",
        x"3234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdeb",
    ];

    let extra_args = client::encode_svm_extra_args_v1(
        expected_compute_units,
        expected_account_is_writable_bitmap,
        expected_allow_out_of_order_execution,
        expected_token_receiver,
        expected_accounts,
    );

    let (
        compute_units,
        account_is_writable_bitmap,
        allow_out_of_order_execution,
        token_receiver,
        accounts,
    ) = decode_svm_extra_args(extra_args);

    assert!(compute_units == expected_compute_units, 0);
    assert!(account_is_writable_bitmap == expected_account_is_writable_bitmap, 0);
    assert!(allow_out_of_order_execution == expected_allow_out_of_order_execution, 0);
    assert!(token_receiver == expected_token_receiver, 0);
    assert!(accounts == expected_accounts, 0);
}
