// ================================================================
//          THIS IS A TEST CONTRACT FOR THE OFFRAMP
// ================================================================

module test::offramp {
    use std::ascii;
    use std::string::{Self, String};
    use std::type_name;
    use std::u256;

    use sui::address;
    use sui::clock;
    use sui::event;
    use sui::hash;
    use sui::package::UpgradeCap;
    use sui::table::{Self, Table};
    use sui::vec_map::{Self, VecMap};
    use sui::object::{Self, UID};
    use sui::derived_object;

    use test::ocr3_base::{Self, OCR3BaseState, OCRConfig};

    public struct CCIPObjectRef has key, store {
        id: UID,
    }

    public struct OffRampState has key, store {
        id: UID,
        package_ids: vector<address>,
    }

    public struct OffRampObject has key {
        id: UID,
    }

    public struct OffRampStatePointer has key, store {
        id: UID,
        off_ramp_state_id: address,
        owner_cap_id: address,
        off_ramp_object_id: address,
    }

    public struct SourceChainConfig has store, drop, copy {
        router: address,
        is_enabled: bool,
        min_seq_nr: u64,
        is_rmn_verification_disabled: bool,
        on_ramp: vector<u8>
    }

    // report public structs
    public struct RampMessageHeader has drop {
        message_id: vector<u8>,
        source_chain_selector: u64,
        dest_chain_selector: u64,
        sequence_number: u64,
        nonce: u64
    }

    public struct Any2SuiRampMessage has drop {
        header: RampMessageHeader,
        sender: vector<u8>,
        data: vector<u8>,
        receiver: address,
        gas_limit: u256,
        token_amounts: vector<Any2SuiTokenTransfer>
    }

    public struct Any2SuiTokenTransfer has drop {
        source_pool_address: vector<u8>,
        // the token's coin metadata object id on SUI
        dest_token_address: address,
        dest_gas_amount: u32,
        extra_data: vector<u8>,

        amount: u256 // This is the amount to transfer, as set on the source chain.
    }

    public struct ExecutionReport has drop {
        source_chain_selector: u64,
        message: Any2SuiRampMessage,
        offchain_token_data: vector<vector<u8>>,
        proofs: vector<vector<u8>>, // Proofs used to construct the merkle root
    }

    // Matches the EVM public struct
    public struct CommitReport has store, drop, copy {
        price_updates: PriceUpdates, // Price updates for the fee_quoter
        blessed_merkle_roots: vector<MerkleRoot>, // Merkle roots that have been blessed by RMN
        unblessed_merkle_roots: vector<MerkleRoot>, // Merkle roots that don't require RMN blessing
        rmn_signatures: vector<vector<u8>> // The signatures for the blessed merkle roots
    }

    public struct PriceUpdates has store, drop, copy {
        token_price_updates: vector<TokenPriceUpdate>,
        gas_price_updates: vector<GasPriceUpdate>
    }

    public struct TokenPriceUpdate has store, drop, copy {
        source_token: address,
        usd_per_token: u256
    }

    public struct GasPriceUpdate has store, drop, copy {
        dest_chain_selector: u64,
        usd_per_unit_gas: u256
    }

    public struct MerkleRoot has store, drop, copy {
        source_chain_selector: u64,
        on_ramp_address: vector<u8>,
        min_seq_nr: u64,
        max_seq_nr: u64,
        merkle_root: vector<u8>
    }

    public struct StaticConfig has store, drop, copy {
        chain_selector: u64,
        rmn_remote: address,
        token_admin_registry: address,
        nonce_manager: address
    }

    // On EVM, the feeQuoter is a dynamic address but due to the Sui implementation using a static
    // upgradable FeeQuoter stored within the state ref, this value is actually static and cannot be
    // accessed by its object id/address directly by users.
    // For compatibility reasons, we keep it as a dynamic config.
    public struct DynamicConfig has store, drop, copy {
        fee_quoter: address,
        permissionless_execution_threshold_seconds: u32 // The delay before manual exec is enabled
    }

    public struct StaticConfigSet has copy, drop {
        chain_selector: u64
    }

    public struct DynamicConfigSet has copy, drop {
        dynamic_config: DynamicConfig
    }

    public struct SourceChainConfigSet has copy, drop {
        source_chain_selector: u64,
        source_chain_config: SourceChainConfig
    }

    public struct SkippedAlreadyExecuted has copy, drop {
        source_chain_selector: u64,
        sequence_number: u64
    }

    public struct ExecutionStateChanged has copy, drop {
        source_chain_selector: u64,
        sequence_number: u64,
        message_id: vector<u8>,
        message_hash: vector<u8>,
        state: u8
    }

    public struct CommitReportAccepted has copy, drop {
        blessed_merkle_roots: vector<MerkleRoot>,
        unblessed_merkle_roots: vector<MerkleRoot>,
        price_updates: PriceUpdates
    }

    public struct SkippedReportExecution has copy, drop {
        source_chain_selector: u64
    }

    // Simple methods to emit events for testing purposes

    /// Emit a StaticConfigSet event
    public fun emit_static_config_set_event(chain_selector: u64) {
        event::emit(StaticConfigSet { chain_selector });
    }

    /// Emit a DynamicConfigSet event
    public fun emit_dynamic_config_set_event(dynamic_config: DynamicConfig) {
        event::emit(DynamicConfigSet { dynamic_config });
    }

    /// Emit a SourceChainConfigSet event
    public fun emit_source_chain_config_set_event(
        source_chain_selector: u64,
        source_chain_config: SourceChainConfig
    ) {
        event::emit(SourceChainConfigSet {
            source_chain_selector,
            source_chain_config
        });
    }

    /// Emit a SkippedAlreadyExecuted event
    public fun emit_skipped_already_executed_event(
        source_chain_selector: u64,
        sequence_number: u64
    ) {
        event::emit(SkippedAlreadyExecuted {
            source_chain_selector,
            sequence_number
        });
    }

    /// Emit an ExecutionStateChanged event
    public fun emit_execution_state_changed_event(
        source_chain_selector: u64,
        sequence_number: u64,
        message_id: vector<u8>,
        message_hash: vector<u8>,
        state: u8
    ) {
        event::emit(ExecutionStateChanged {
            source_chain_selector,
            sequence_number,
            message_id,
            message_hash,
            state
        });
    }

    public fun emit_sample_execution_state_changed_event() {
        let source_chain_selector = 24;
        let sequence_number = 12;
        let message_id = x"1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef";
        let message_hash = x"1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef";
        let state = 1;
        emit_execution_state_changed_event(source_chain_selector, sequence_number, message_id, message_hash, state);
    }

    /// Emit a CommitReportAccepted event
    public fun emit_commit_report_accepted_event(
        blessed_merkle_roots: vector<MerkleRoot>,
        unblessed_merkle_roots: vector<MerkleRoot>,
        price_updates: PriceUpdates
    ) {
        event::emit(CommitReportAccepted {
            blessed_merkle_roots,
            unblessed_merkle_roots,
            price_updates
        });
    }

    /// Emit a SkippedReportExecution event
    public fun emit_skipped_report_execution_event(source_chain_selector: u64) {
        event::emit(SkippedReportExecution { source_chain_selector });
    }

    // Helper functions to create test structures

    /// Create a test SourceChainConfig
    public fun create_test_source_chain_config(
        router: address,
        is_enabled: bool,
        min_seq_nr: u64,
        is_rmn_verification_disabled: bool,
        on_ramp: vector<u8>
    ): SourceChainConfig {
        SourceChainConfig {
            router,
            is_enabled,
            min_seq_nr,
            is_rmn_verification_disabled,
            on_ramp
        }
    }

    /// Create a default test SourceChainConfig with reasonable values
    public fun create_default_test_source_chain_config(): SourceChainConfig {
        create_test_source_chain_config(
            @0x1,
            true,
            1,
            false,
            x"1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
        )
    }

    /// Create a test DynamicConfig
    public fun create_test_dynamic_config(
        fee_quoter: address,
        permissionless_execution_threshold_seconds: u32
    ): DynamicConfig {
        DynamicConfig {
            fee_quoter,
            permissionless_execution_threshold_seconds
        }
    }

    /// Create a default test DynamicConfig with reasonable values
    public fun create_default_test_dynamic_config(): DynamicConfig {
        create_test_dynamic_config(
            @0x2,
            3600
        )
    }

    /// Create a test RampMessageHeader
    public fun create_test_ramp_message_header(
        message_id: vector<u8>,
        source_chain_selector: u64,
        dest_chain_selector: u64,
        sequence_number: u64,
        nonce: u64
    ): RampMessageHeader {
        RampMessageHeader {
            message_id,
            source_chain_selector,
            dest_chain_selector,
            sequence_number,
            nonce
        }
    }

    /// Create a test Any2SuiTokenTransfer
    public fun create_test_any2sui_token_transfer(
        source_pool_address: vector<u8>,
        dest_token_address: address,
        dest_gas_amount: u32,
        extra_data: vector<u8>,
        amount: u256
    ): Any2SuiTokenTransfer {
        Any2SuiTokenTransfer {
            source_pool_address,
            dest_token_address,
            dest_gas_amount,
            extra_data,
            amount
        }
    }

    /// Create a test Any2SuiRampMessage
    public fun create_test_any2sui_ramp_message(
        header: RampMessageHeader,
        sender: vector<u8>,
        data: vector<u8>,
        receiver: address,
        gas_limit: u256,
        token_amounts: vector<Any2SuiTokenTransfer>
    ): Any2SuiRampMessage {
        Any2SuiRampMessage {
            header,
            sender,
            data,
            receiver,
            gas_limit,
            token_amounts
        }
    }

    /// Create a test MerkleRoot
    public fun create_test_merkle_root(
        source_chain_selector: u64,
        on_ramp_address: vector<u8>,
        min_seq_nr: u64,
        max_seq_nr: u64,
        merkle_root: vector<u8>
    ): MerkleRoot {
        MerkleRoot {
            source_chain_selector,
            on_ramp_address,
            min_seq_nr,
            max_seq_nr,
            merkle_root
        }
    }

    /// Create a test PriceUpdates
    public fun create_test_price_updates(
        token_price_updates: vector<TokenPriceUpdate>,
        gas_price_updates: vector<GasPriceUpdate>
    ): PriceUpdates {
        PriceUpdates {
            token_price_updates,
            gas_price_updates
        }
    }

    /// Create a test TokenPriceUpdate
    public fun create_test_token_price_update(
        source_token: address,
        usd_per_token: u256
    ): TokenPriceUpdate {
        TokenPriceUpdate {
            source_token,
            usd_per_token
        }
    }

    /// Create a test GasPriceUpdate
    public fun create_test_gas_price_update(
        dest_chain_selector: u64,
        usd_per_unit_gas: u256
    ): GasPriceUpdate {
        GasPriceUpdate {
            dest_chain_selector,
            usd_per_unit_gas
        }
    }

    /// Create test chain selectors
    public fun create_test_chain_selectors(): vector<u64> {
        vector[1, 137, 42161, 10]  // Ethereum, Polygon, Arbitrum, Optimism
    }

    /// Create test message IDs
    public fun create_test_message_ids(): vector<vector<u8>> {
        vector[
            x"1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
            x"2234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdea"
        ]
    }

    /// Create test on-ramp addresses
    public fun create_test_on_ramp_addresses(): vector<vector<u8>> {
        vector[
            x"3234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdeb",
            x"4234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdec"
        ]
    }

    public fun type_and_version(): String {
        string::utf8(b"OffRamp 1.6.0")
    }

    public struct OFFRAMP has drop {}

    fun init(_witness: OFFRAMP, ctx: &mut TxContext) {
        let tn = type_name::get_with_original_ids<OFFRAMP>();
        let package_bytes = ascii::into_bytes(tn.get_address());
        let package_id = address::from_ascii_bytes(&package_bytes);

        let state = OffRampState {
            id: object::new(ctx),
            package_ids: vector[
                package_id
            ]
        };
        let ref = CCIPObjectRef {
            id: object::new(ctx)
        };
        let config = SourceChainConfig {
            router: ctx.sender(),
            is_enabled: true,
            min_seq_nr: 0,
            is_rmn_verification_disabled: true,
            on_ramp: vector[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1,1,1,1]
        };
        let config_set = SourceChainConfigSet {
            source_chain_selector: 16015286601757825753,
            source_chain_config: config
        };
        event::emit(config_set);

        let config_set2 = SourceChainConfigSet {
            source_chain_selector: 14767482510784806043,
            source_chain_config: config
        };
        event::emit(config_set2);

        let mut off_ramp_object = OffRampObject {
            id: object::new(ctx)
        };

        let pointer = OffRampStatePointer {
            id: object::new(ctx),
            off_ramp_state_id: object::uid_to_address(&state.id),
            owner_cap_id: @0x0,
            off_ramp_object_id: object::id_address(&off_ramp_object),
        };

        let another_state_object = OffRampState {
            id: derived_object::claim(&mut off_ramp_object.id, b"OffRampState"),
            package_ids: vector[
                package_id
            ],
        };
        
        transfer::share_object(state);
        transfer::share_object(another_state_object);
        transfer::share_object(off_ramp_object);
        transfer::public_share_object(ref);
        transfer::transfer(pointer, package_id);
    }

    public fun add_package_id(state: &mut OffRampState, package_id: address) {
        state.package_ids.push_back(package_id);
    }

    public fun remove_package_id(state: &mut OffRampState, package_id: address) {
        state.package_ids.remove(0);
    }

    public fun get_all_source_chain_configs(): (vector<u64>, vector<SourceChainConfig>) {
        let mut chain_selectors = vector[];
        let mut chain_configs = vector[];

        chain_selectors.push_back(16015286);
        chain_configs.push_back(create_default_test_source_chain_config());

        (chain_selectors, chain_configs)
    }

    // ================================================================
    // |                          Execution                           |
    // ================================================================

    const ENotImplemented: u64 = 1;

    public fun init_execute(
        ref: &CCIPObjectRef,
        state: &mut OffRampState,
        clock: &clock::Clock,
        report_context: vector<vector<u8>>,
        report: vector<u8>,
        ctx: &mut TxContext
    ) {}

    public fun finish_execute() {
        assert!(false, ENotImplemented);
    }
}