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

    use test::ocr3_base::{Self, OCR3BaseState, OCRConfig};

    public struct CCIPObjectRef has key, store {
        id: UID,
    }

    public struct OffRampState has key, store {
        id: UID,
    }

    public struct OffRampStatePointer has key, store {
        id: UID,
        off_ramp_state_id: address,
        owner_cap_id: address,
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

    public fun type_and_version(): String {
        string::utf8(b"OffRamp 1.6.0")
    }

    public struct OFFRAMP has drop {}

    fun init(_witness: OFFRAMP, ctx: &mut TxContext) {
        let state = OffRampState {
            id: object::new(ctx)
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
        
        transfer::share_object(state);
        transfer::public_share_object(ref);
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