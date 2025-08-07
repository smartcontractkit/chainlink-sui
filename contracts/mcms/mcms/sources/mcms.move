module mcms::mcms;

use mcms::bcs_stream;
use mcms::mcms_account::{Self, OwnerCap, AccountState};
use mcms::mcms_deployer::{Self, DeployerState};
use mcms::mcms_registry::{Self, ExecutingCallbackParams, Registry};
use mcms::params;
use std::string::String;
use sui::bcs;
use sui::clock::Clock;
use sui::ecdsa_k1;
use sui::event;
use sui::hash::keccak256;
use sui::package::UpgradeTicket;
use sui::table::{Self, Table};
use sui::vec_map::{Self, VecMap};
use sui::vec_set::{Self, VecSet};

const BYPASSER_ROLE: u8 = 0;
const CANCELLER_ROLE: u8 = 1;
const PROPOSER_ROLE: u8 = 2;
const TIMELOCK_ROLE: u8 = 3;
const MAX_ROLE: u8 = 4;

const NUM_GROUPS: u64 = 32;
const MAX_NUM_SIGNERS: u64 = 200;

// equivalent to initializing empty uint8[NUM_GROUPS] in Solidity
const VEC_NUM_GROUPS: vector<u8> = vector[
    0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
];

// TODO: Change hash to SUI (MANY_CHAIN_MULTI_SIG_DOMAIN_SEPARATOR_METADATA_SUI)
// keccak256("MANY_CHAIN_MULTI_SIG_DOMAIN_SEPARATOR_METADATA_APTOS")
const MANY_CHAIN_MULTI_SIG_DOMAIN_SEPARATOR_METADATA: vector<u8> =
    x"a71d47b6c00b64ee21af96a1d424cb2dcbbed12becdcd3b4e6c7fc4c2f80a697";

// TODO: Change hash to SUI (MANY_CHAIN_MULTI_SIG_DOMAIN_SEPARATOR_OP_SUI)
// keccak256("MANY_CHAIN_MULTI_SIG_DOMAIN_SEPARATOR_OP_APTOS")
const MANY_CHAIN_MULTI_SIG_DOMAIN_SEPARATOR_OP: vector<u8> =
    x"e5a6d1256b00d7ec22512b6b60a3f4d75c559745d2dbf309f77b8b756caabe14";

/// Special timestamp value indicating an operation is done
const DONE_TIMESTAMP: u64 = 1;

const ZERO_HASH: vector<u8> = vector[
    0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
];

public struct MultisigState has key, store {
    id: UID,
    bypasser: Multisig,
    canceller: Multisig,
    proposer: Multisig,
}

public struct Multisig has store {
    role: u8,
    /// signers is used to easily validate the existence of the signer by its address. We still
    /// have signers stored in config in order to easily deactivate them when a new config is set.
    signers: VecMap<vector<u8>, Signer>,
    config: Config,
    /// Remember signed hashes that this contract has seen. Each signed hash can only be set once.
    seen_signed_hashes: VecMap<vector<u8>, bool>,
    expiring_root_and_op_count: ExpiringRootAndOpCount,
    root_metadata: RootMetadata,
}

public struct Signer has copy, drop, store {
    addr: vector<u8>,
    index: u8,
    group: u8,
}

public struct Config has copy, drop, store {
    signers: vector<Signer>,
    /// group_quorums[i] stores the quorum for the i-th signer group. Any group with
    /// group_quorums[i] = 0 is considered disabled. The i-th group is successful if
    /// it is enabled and at least group_quorums[i] of its children are successful.
    group_quorums: vector<u8>,
    /// group_parents[i] stores the parent group of the i-th signer group. We ensure that the
    /// groups form a tree structure (where the root/0-th signer group points to itself as
    /// parent) by enforcing
    /// - (i != 0) implies (group_parents[i] < i)
    /// - group_parents[0] == 0
    group_parents: vector<u8>,
}

public struct ExpiringRootAndOpCount has copy, drop, store {
    root: vector<u8>,
    valid_until: u64,
    op_count: u64,
}

public struct Op has copy, drop {
    role: u8,
    chain_id: u256,
    multisig: address,
    nonce: u64,
    to: address,
    module_name: String,
    function_name: String,
    data: vector<u8>,
}

public struct RootMetadata has copy, drop, store {
    role: u8,
    chain_id: u256,
    multisig: address,
    pre_op_count: u64,
    post_op_count: u64,
    override_previous_root: bool,
}

/// Hot potato only issued when merkle proof is verified against the multisig root
/// This is issued when successful validation from `mcms::execute`
public struct TimelockCallbackParams {
    module_name: String,
    function_name: String,
    data: vector<u8>,
    role: u8,
}

public struct MultisigStateInitialized has copy, drop {
    bypasser: u8,
    canceller: u8,
    proposer: u8,
}

public struct ConfigSet has copy, drop {
    role: u8,
    config: Config,
    is_root_cleared: bool,
}

public struct NewRoot has copy, drop {
    role: u8,
    root: vector<u8>,
    valid_until: u64,
    metadata: RootMetadata,
}

public struct OpExecuted has copy, drop {
    role: u8,
    chain_id: u256,
    multisig: address,
    nonce: u64,
    to: address,
    module_name: String,
    function_name: String,
    data: vector<u8>,
}

const EInvalidRole: u64 = 1;
const EInvalidRootLen: u64 = 2;
const EMissingConfig: u64 = 3;
const EWrongPreOpCount: u64 = 4;
const EWrongPostOpCount: u64 = 5;
const EProofCannotBeVerified: u64 = 6;
const EAlreadySeenHash: u64 = 7;
const EValidUntilExpired: u64 = 8;
const EWrongMultisig: u64 = 9;
const EInvalidSigner: u64 = 10;
const ESignerInDisabledGroup: u64 = 11;
const EInsufficientSigners: u64 = 12;
const EInvalidGroupQuorumLen: u64 = 13;
const EInvalidGroupParentsLen: u64 = 14;
const EOutOfBoundsGroup: u64 = 15;
const EOutOfBoundsGroupQuorum: u64 = 16;
const ESignerAddrMustBeIncreasing: u64 = 17;
const EPendingOps: u64 = 18;
const EInvalidNumSigners: u64 = 19;
const ESignerGroupsLenMismatch: u64 = 20;
const EGroupTreeNotWellFormed: u64 = 21;
const EInvalidSignerAddrLen: u64 = 22;
const EPostOpCountReached: u64 = 23;
const EWrongNonce: u64 = 24;
const EInvalidModuleName: u64 = 25;
const EInvalidFunctionName: u64 = 26;
const ENotAuthorizedRole: u64 = 27;
const EInsufficientDelay: u64 = 28;
const EOperationAlreadyScheduled: u64 = 29;
const EOperationNotReady: u64 = 30;
const EMissingDependency: u64 = 31;
const ENotTimeLockRole: u64 = 32;
const EInvalidIndex: u64 = 33;
const EFunctionBlocked: u64 = 34;
const EInvalidParameters: u64 = 35;
const EOperationCannotBeCancelled: u64 = 36;
const EUnknownMCMSAccountModuleFunction: u64 = 37;
const EUnknownMCMSModule: u64 = 38;
const EUnknownMCMSDeployerModuleFunction: u64 = 39;
const EInvalidMCMS: u64 = 40;

public struct MCMS has drop {}

/// This acts as a proof type for the MCMS module.
public struct McmsCallback has drop {}

fun init(_witness: MCMS, ctx: &mut TxContext) {
    let bypasser = create_multisig(BYPASSER_ROLE);
    let canceller = create_multisig(CANCELLER_ROLE);
    let proposer = create_multisig(PROPOSER_ROLE);

    let multisig_state = MultisigState {
        id: object::new(ctx),
        bypasser,
        canceller,
        proposer,
    };

    let timelock = Timelock {
        id: object::new(ctx),
        min_delay: 0,
        timestamps: table::new(ctx),
        blocked_functions: vec_set::empty(),
    };

    event::emit(MultisigStateInitialized {
        bypasser: multisig_state.bypasser.role,
        canceller: multisig_state.canceller.role,
        proposer: multisig_state.proposer.role,
    });

    event::emit(TimelockInitialized {
        min_delay: timelock.min_delay,
    });

    transfer::share_object(multisig_state);
    transfer::share_object(timelock);
}

fun create_multisig(role: u8): Multisig {
    Multisig {
        role,
        signers: vec_map::empty(),
        config: Config {
            signers: vector[],
            group_quorums: VEC_NUM_GROUPS,
            group_parents: VEC_NUM_GROUPS,
        },
        seen_signed_hashes: vec_map::empty(),
        expiring_root_and_op_count: ExpiringRootAndOpCount {
            root: vector[],
            valid_until: 0,
            op_count: 0,
        },
        root_metadata: RootMetadata {
            role,
            chain_id: 0,
            multisig: mcms_registry::get_multisig_address(),
            pre_op_count: 0,
            post_op_count: 0,
            override_previous_root: false,
        },
    }
}

/// @notice set_root Sets a new expiring root.
///
/// @param root is the new expiring root.
/// @param valid_until is the time by which root is valid
/// @param chain_id is the chain id of the chain on which the root is valid
/// @param multisig is the address of the multisig to set the root for
/// @param pre_op_count is the number of operations that have been executed before this root was set
/// @param post_op_count is the number of operations that have been executed after this root was set
/// @param override_previous_root is a boolean that indicates whether to override the previous root
/// @param metadata_proof is the MerkleProof of inclusion of the metadata in the Merkle tree.
/// @param signatures the ECDSA signatures on (root, valid_until).
///
/// @dev the message (root, valid_until) should be signed by a sufficient set of signers.
/// This signature authenticates also the metadata.
///
/// @dev this method can be executed by anyone who has the root and valid signatures.
/// as we validate the correctness of signatures, this imposes no risk.
public entry fun set_root(
    state: &mut MultisigState,
    clock: &Clock,
    role: u8,
    root: vector<u8>,
    valid_until: u64,
    chain_id: u256,
    multisig_addr: address,
    pre_op_count: u64,
    post_op_count: u64,
    override_previous_root: bool,
    metadata_proof: vector<vector<u8>>,
    signatures: vector<vector<u8>>,
    _ctx: &mut TxContext,
) {
    assert!(is_valid_role(role), EInvalidRole);

    let metadata = RootMetadata {
        role,
        chain_id,
        multisig: multisig_addr,
        pre_op_count,
        post_op_count,
        override_previous_root,
    };

    let signed_hash = compute_eth_message_hash(root, valid_until);

    // Validate that `multisig` is a registered multisig for `role`.
    let multisig = borrow_multisig_mut(state, role);

    assert!(!multisig.seen_signed_hashes.contains(&signed_hash), EAlreadySeenHash);
    assert!(get_timestamp_seconds(clock) <= valid_until, EValidUntilExpired);

    // TODO: No support for chain_ids yet
    // assert!(metadata.chain_id == (chain_ids::get() as u256), EWrongChainId);

    assert!(metadata.multisig == mcms_registry::get_multisig_address(), EWrongMultisig);

    let op_count = multisig.expiring_root_and_op_count.op_count;
    assert!(
        override_previous_root || op_count == multisig.root_metadata.post_op_count,
        EPendingOps,
    );

    assert!(op_count == metadata.pre_op_count, EWrongPreOpCount);
    assert!(metadata.pre_op_count <= metadata.post_op_count, EWrongPostOpCount);

    let metadata_leaf_hash = hash_metadata_leaf(metadata);
    assert!(verify_merkle_proof(metadata_proof, root, metadata_leaf_hash), EProofCannotBeVerified);

    let mut prev_address = vector[];
    let mut group_vote_counts: vector<u8> = vector[];
    params::right_pad_vec(&mut group_vote_counts, NUM_GROUPS);

    let signatures_len = signatures.length();
    let mut i = 0;
    while (i < signatures_len) {
        let signature = signatures[i];
        let signer_addr = ecdsa_recover_evm_addr(signed_hash, signature);
        // the off-chain system is required to sort the signatures by the
        // signer address in an increasing order
        if (i > 0) {
            assert!(params::vector_u8_gt(&signer_addr, &prev_address), ESignerAddrMustBeIncreasing);
        };
        assert!(multisig.signers.contains(&signer_addr), EInvalidSigner);
        prev_address = signer_addr;

        let signer = *multisig.signers.get(&signer_addr);

        // check group quorums
        let mut group: u8 = signer.group;
        while (true) {
            let group_vote_count = group_vote_counts.borrow_mut((group as u64));
            *group_vote_count = *group_vote_count + 1;

            let quorum = multisig.config.group_quorums.borrow((group as u64));
            if (*group_vote_count != *quorum) {
                // bail out unless we just hit the quorum. we only hit each quorum once,
                // so we never move on to the parent of a group more than once.
                break
            };

            if (group == 0) {
                // root group reached
                break
            };

            // group quorum reached, restart loop and check parent group
            group = multisig.config.group_parents[(group as u64)];
        };
        i = i + 1;
    };

    // the group at the root of the tree (with index 0) determines whether the vote passed,
    // we cannot proceed if it isn't configured with a valid (non-zero) quorum
    let root_group_quorum = multisig.config.group_quorums[0];
    assert!(root_group_quorum != 0, EMissingConfig);

    // check root group reached quorum
    let root_group_vote_count = group_vote_counts[0];
    assert!(root_group_vote_count >= root_group_quorum, EInsufficientSigners);

    multisig.seen_signed_hashes.insert(signed_hash, true);
    multisig.expiring_root_and_op_count =
        ExpiringRootAndOpCount {
            root,
            valid_until,
            op_count: metadata.pre_op_count,
        };
    multisig.root_metadata = metadata;

    event::emit(NewRoot {
        role,
        root,
        valid_until,
        metadata: RootMetadata {
            role,
            chain_id,
            multisig: multisig_addr,
            pre_op_count: metadata.pre_op_count,
            post_op_count: metadata.post_op_count,
            override_previous_root: metadata.override_previous_root,
        },
    });
}

// https://github.com/MystenLabs/sui/blob/main/examples/move/crypto/ecdsa_k1/sources/example.move#L62
fun ecdsa_recover_evm_addr(msg: vector<u8>, mut sig: vector<u8>): vector<u8> {
    // Normalize the last byte of the signature to be 0 or 1.
    let v = &mut sig[64];
    if (*v == 27) {
        *v = 0;
    } else if (*v == 28) {
        *v = 1;
    } else if (*v > 35) {
        *v = (*v - 1) % 2;
    };

    // Ethereum signature is produced with Keccak256 hash of the message, so the last param is 0.
    let pubkey = ecdsa_k1::secp256k1_ecrecover(&sig, &msg, 0);
    let uncompressed = ecdsa_k1::decompress_pubkey(&pubkey);

    // Take the last 64 bytes of the uncompressed pubkey.
    let mut uncompressed_64 = vector[];
    let mut i = 1;
    while (i < 65) {
        uncompressed_64.push_back(uncompressed[i]);
        i = i + 1;
    };

    // Take the last 20 bytes of the hash of the 64-bytes uncompressed pubkey.
    let hashed = keccak256(&uncompressed_64);
    let mut addr = vector[];
    let mut i = 12;
    while (i < 32) {
        addr.push_back(hashed[i]);
        i = i + 1;
    };

    addr
}

public fun execute(
    state: &mut MultisigState,
    clock: &Clock,
    role: u8,
    chain_id: u256,
    multisig_addr: address,
    nonce: u64,
    to: address,
    module_name: String,
    function_name: String,
    data: vector<u8>,
    proof: vector<vector<u8>>,
): TimelockCallbackParams {
    assert!(is_valid_role(role), EInvalidRole);

    let op = Op {
        role,
        chain_id,
        multisig: multisig_addr,
        nonce,
        to,
        module_name,
        function_name,
        data,
    };
    let multisig = borrow_multisig_mut(state, role);

    assert!(
        multisig.root_metadata.post_op_count
                > multisig.expiring_root_and_op_count.op_count,
        EPostOpCountReached,
    );

    // TODO: No support for chain_ids yet
    // assert!(chain_id == (chain_ids::get() as u256), EWrongChainId);
    assert!(
        get_timestamp_seconds(clock) <= multisig.expiring_root_and_op_count.valid_until,
        EValidUntilExpired,
    );
    assert!(
        multisig.root_metadata.multisig == mcms_registry::get_multisig_address(),
        EWrongMultisig,
    );
    assert!(op.multisig == mcms_registry::get_multisig_address(), EWrongMultisig);
    assert!(nonce == multisig.expiring_root_and_op_count.op_count, EWrongNonce);

    // computes keccak256(abi.encode(MANY_CHAIN_MULTI_SIG_DOMAIN_SEPARATOR_OP, op))
    let hashed_leaf = hash_op_leaf(MANY_CHAIN_MULTI_SIG_DOMAIN_SEPARATOR_OP, op);
    assert!(
        verify_merkle_proof(proof, multisig.expiring_root_and_op_count.root, hashed_leaf),
        EProofCannotBeVerified,
    );

    multisig.expiring_root_and_op_count.op_count = multisig.expiring_root_and_op_count.op_count + 1;

    // Only allow dispatching to timelock functions
    assert!(mcms_registry::get_multisig_address() == op.to, EInvalidMCMS);
    assert!(*op.module_name.as_bytes() == b"mcms", EInvalidModuleName);

    event::emit(OpExecuted {
        role,
        chain_id,
        multisig: multisig_addr,
        nonce,
        to,
        module_name,
        function_name,
        data,
    });

    // Create TimelockCallbackParams hot potato
    // Caller must call into Timelock functions with the hot potato
    TimelockCallbackParams {
        module_name,
        function_name,
        data,
        role,
    }
}

// ================================ Timelock Callback Functions ================================ //

// These functions are called with `TimelockCallbackParams` from `execute`
// `TimelockCallbackParams` is only issued when merkle proof is verified against the multisig root
// These functions are ready to be executed

/// `dispatch_timelock_` functions should only be called after calling mcms::execute_timelock_schedule_batch
/// This can be public as `TimelockCallbackParams` is only issued when merkle proof is verified against the multisig root
public fun dispatch_timelock_schedule_batch(
    timelock: &mut Timelock,
    clock: &Clock,
    timelock_callback_params: TimelockCallbackParams, // hot potato
    ctx: &mut TxContext,
) {
    let TimelockCallbackParams {
        module_name,
        function_name,
        data,
        role,
    } = timelock_callback_params;

    assert!(*module_name.as_bytes() == b"mcms", EInvalidModuleName);
    assert!(*function_name.as_bytes() == b"timelock_schedule_batch", EInvalidFunctionName);

    let (
        targets,
        module_names,
        function_names,
        datas,
        predecessor,
        salt,
        delay,
    ) = deserialize_timelock_schedule_batch(data);

    timelock_schedule_batch(
        timelock,
        clock,
        role,
        targets,
        module_names,
        function_names,
        datas,
        predecessor,
        salt,
        delay,
        ctx,
    )
}

public fun dispatch_timelock_execute_batch(
    timelock: &mut Timelock,
    clock: &Clock,
    timelock_callback_params: TimelockCallbackParams,
    ctx: &mut TxContext,
): vector<ExecutingCallbackParams> {
    let TimelockCallbackParams {
        module_name,
        function_name,
        data,
        role,
    } = timelock_callback_params;

    assert!(is_valid_role(role), EInvalidRole);
    assert!(*module_name.as_bytes() == b"mcms", EInvalidModuleName);
    assert!(*function_name.as_bytes() == b"timelock_execute_batch", EInvalidFunctionName);

    let (
        targets,
        module_names,
        function_names,
        datas,
        predecessor,
        salt,
    ) = deserialize_timelock_execute_batch(data);

    timelock_execute_batch(
        timelock,
        clock,
        targets,
        module_names,
        function_names,
        datas,
        predecessor,
        salt,
        ctx,
    )
}

public fun dispatch_timelock_bypasser_execute_batch(
    timelock_callback_params: TimelockCallbackParams,
    ctx: &mut TxContext,
): vector<ExecutingCallbackParams> {
    let TimelockCallbackParams {
        module_name,
        function_name,
        data,
        role,
    } = timelock_callback_params;

    assert!(*module_name.as_bytes() == b"mcms", EInvalidModuleName);
    assert!(*function_name.as_bytes() == b"timelock_bypasser_execute_batch", EInvalidFunctionName);

    let (
        targets,
        module_names,
        function_names,
        datas,
    ) = deserialize_timelock_bypasser_execute_batch(data);

    timelock_bypasser_execute_batch(
        role,
        targets,
        module_names,
        function_names,
        datas,
        ctx,
    )
}

public fun dispatch_timelock_cancel(
    timelock: &mut Timelock,
    timelock_callback_params: TimelockCallbackParams,
    ctx: &mut TxContext,
) {
    let TimelockCallbackParams { module_name, function_name, data, role } =
        timelock_callback_params;

    assert!(*module_name.as_bytes() == b"mcms", EInvalidModuleName);
    assert!(*function_name.as_bytes() == b"timelock_cancel", EInvalidFunctionName);

    let id = deserialize_timelock_cancel(data);
    timelock_cancel(timelock, role, id, ctx)
}

public fun dispatch_timelock_update_min_delay(
    timelock: &mut Timelock,
    timelock_callback_params: TimelockCallbackParams,
    ctx: &mut TxContext,
) {
    let TimelockCallbackParams { module_name, function_name, data, role } =
        timelock_callback_params;

    assert!(*module_name.as_bytes() == b"mcms", EInvalidModuleName);
    assert!(*function_name.as_bytes() == b"timelock_update_min_delay", EInvalidFunctionName);

    let new_min_delay = deserialize_timelock_update_min_delay(data);
    timelock_update_min_delay(timelock, role, new_min_delay, ctx)
}

public fun dispatch_timelock_block_function(
    timelock: &mut Timelock,
    timelock_callback_params: TimelockCallbackParams,
    ctx: &mut TxContext,
) {
    let TimelockCallbackParams { module_name, function_name, data, role } =
        timelock_callback_params;

    assert!(*module_name.as_bytes() == b"mcms", EInvalidModuleName);
    assert!(*function_name.as_bytes() == b"timelock_block_function", EInvalidFunctionName);

    let (target, module_name, function_name) = deserialize_timelock_function_action(data);
    timelock_block_function(timelock, role, target, module_name, function_name, ctx)
}

public fun dispatch_timelock_unblock_function(
    timelock: &mut Timelock,
    timelock_callback_params: TimelockCallbackParams,
    ctx: &mut TxContext,
) {
    let TimelockCallbackParams { module_name, function_name, data, role } =
        timelock_callback_params;

    assert!(*module_name.as_bytes() == b"mcms", EInvalidModuleName);
    assert!(*function_name.as_bytes() == b"timelock_unblock_function", EInvalidFunctionName);

    let (target, module_name, function_name) = deserialize_timelock_function_action(data);
    timelock_unblock_function(timelock, role, target, module_name, function_name, ctx)
}

/*
    ================================ Execute Timelock Functions ================================

    These functions are called with `ExecutingCallbackParams`
    `ExecutingCallbackParams` is only issued from `timelock_execute_batch` or `timelock_bypasser_execute_batch`
    These functions are ready to be executed therefore no validation is needed
    */

public fun execute_dispatch_to_account(
    registry: &mut Registry,
    account_state: &mut AccountState,
    executing_callback_params: ExecutingCallbackParams,
    ctx: &mut TxContext,
) {
    let (target, module_name, function_name, data) = mcms_registry::get_callback_params_from_mcms(
        executing_callback_params,
    );
    assert!(target == mcms_registry::get_multisig_address(), EWrongMultisig);
    assert!(*module_name.as_bytes() == b"mcms_account", EUnknownMCMSAccountModuleFunction);

    let function_name_bytes = *function_name.as_bytes();
    let mut stream = bcs_stream::new(data);

    if (function_name_bytes == b"transfer_ownership") {
        let target = bcs_stream::deserialize_address(&mut stream);
        bcs_stream::assert_is_consumed(&stream);
        let cap = mcms_registry::borrow_owner_cap(registry);
        mcms_account::transfer_ownership(cap, account_state, target, ctx);
    } else if (function_name_bytes == b"accept_ownership_as_timelock") {
        mcms_account::accept_ownership_as_timelock(
            account_state,
            ctx,
        );
    } else if (function_name_bytes == b"execute_ownership_transfer") {
        let target = bcs_stream::deserialize_address(&mut stream);
        bcs_stream::assert_is_consumed(&stream);
        let owner_cap = mcms_registry::release_cap(registry, mcms_registry::create_mcms_proof());
        mcms_account::execute_ownership_transfer(owner_cap, account_state, registry, target, ctx);
    } else {
        abort EUnknownMCMSAccountModuleFunction
    }
}

public fun execute_dispatch_to_deployer(
    registry: &mut Registry,
    deployer_state: &mut DeployerState,
    executing_callback_params: ExecutingCallbackParams,
    ctx: &mut TxContext,
): UpgradeTicket {
    let (target, module_name, function_name, data) = mcms_registry::get_callback_params_from_mcms(
        executing_callback_params,
    );

    assert!(target == mcms_registry::get_multisig_address(), EUnknownMCMSModule);
    assert!(*module_name.as_bytes() == b"mcms_deployer", EUnknownMCMSDeployerModuleFunction);

    let function_name_bytes = *function_name.as_bytes();
    let mut stream = bcs_stream::new(data);

    if (function_name_bytes == b"authorize_upgrade") {
        let policy = bcs_stream::deserialize_u8(&mut stream);
        bcs_stream::assert_is_consumed(&stream);
        let digest = bcs_stream::deserialize_vector_u8(&mut stream);
        bcs_stream::assert_is_consumed(&stream);
        let code_address = bcs_stream::deserialize_address(&mut stream);
        bcs_stream::assert_is_consumed(&stream);

        let owner_cap = mcms_registry::borrow_owner_cap(registry);
        mcms_deployer::authorize_upgrade(
            owner_cap,
            deployer_state,
            policy,
            digest,
            code_address,
            ctx,
        )
    } else {
        abort EUnknownMCMSDeployerModuleFunction
    }
}

public fun execute_timelock_schedule_batch(
    timelock: &mut Timelock,
    clock: &Clock,
    executing_callback_params: ExecutingCallbackParams,
    ctx: &mut TxContext,
) {
    let (_, module_name, function_name, data) = mcms_registry::get_callback_params_from_mcms(
        executing_callback_params,
    );
    assert!(*module_name.as_bytes() == b"mcms", EInvalidModuleName);
    assert!(*function_name.as_bytes() == b"timelock_schedule_batch", EInvalidFunctionName);

    let (
        targets,
        module_names,
        function_names,
        datas,
        predecessor,
        salt,
        delay,
    ) = deserialize_timelock_schedule_batch(data);

    timelock_schedule_batch(
        timelock,
        clock,
        TIMELOCK_ROLE,
        targets,
        module_names,
        function_names,
        datas,
        predecessor,
        salt,
        delay,
        ctx,
    )
}

public fun execute_timelock_execute_batch(
    timelock: &mut Timelock,
    clock: &Clock,
    executing_callback_params: ExecutingCallbackParams,
    ctx: &mut TxContext,
): vector<ExecutingCallbackParams> {
    let (_, module_name, function_name, data) = mcms_registry::get_callback_params_from_mcms(
        executing_callback_params,
    );
    assert!(*module_name.as_bytes() == b"mcms", EInvalidModuleName);
    assert!(*function_name.as_bytes() == b"timelock_execute_batch", EInvalidFunctionName);

    let (
        targets,
        module_names,
        function_names,
        datas,
        predecessor,
        salt,
    ) = deserialize_timelock_execute_batch(data);

    timelock_execute_batch(
        timelock,
        clock,
        targets,
        module_names,
        function_names,
        datas,
        predecessor,
        salt,
        ctx,
    )
}

public fun execute_timelock_bypasser_execute_batch(
    executing_callback_params: ExecutingCallbackParams,
    ctx: &mut TxContext,
): vector<ExecutingCallbackParams> {
    let (_, module_name, function_name, data) = mcms_registry::get_callback_params_from_mcms(
        executing_callback_params,
    );
    assert!(*module_name.as_bytes() == b"mcms", EInvalidModuleName);
    assert!(*function_name.as_bytes() == b"timelock_bypasser_execute_batch", EInvalidFunctionName);

    let (
        targets,
        module_names,
        function_names,
        datas,
    ) = deserialize_timelock_bypasser_execute_batch(data);

    timelock_bypasser_execute_batch(
        TIMELOCK_ROLE,
        targets,
        module_names,
        function_names,
        datas,
        ctx,
    )
}

public fun execute_timelock_cancel(
    timelock: &mut Timelock,
    executing_callback_params: ExecutingCallbackParams,
    ctx: &mut TxContext,
) {
    let (_, module_name, function_name, data) = mcms_registry::get_callback_params_from_mcms(
        executing_callback_params,
    );
    assert!(*module_name.as_bytes() == b"mcms", EInvalidModuleName);
    assert!(*function_name.as_bytes() == b"timelock_cancel", EInvalidFunctionName);

    let id = deserialize_timelock_cancel(data);
    timelock_cancel(timelock, TIMELOCK_ROLE, id, ctx)
}

public fun execute_timelock_update_min_delay(
    timelock: &mut Timelock,
    executing_callback_params: ExecutingCallbackParams,
    ctx: &mut TxContext,
) {
    let (_, module_name, function_name, data) = mcms_registry::get_callback_params_from_mcms(
        executing_callback_params,
    );
    assert!(*module_name.as_bytes() == b"mcms", EInvalidModuleName);
    assert!(*function_name.as_bytes() == b"timelock_update_min_delay", EInvalidFunctionName);

    let new_min_delay = deserialize_timelock_update_min_delay(data);
    timelock_update_min_delay(timelock, TIMELOCK_ROLE, new_min_delay, ctx)
}

public fun execute_timelock_block_function(
    timelock: &mut Timelock,
    executing_callback_params: ExecutingCallbackParams,
    ctx: &mut TxContext,
) {
    let (_, module_name, function_name, data) = mcms_registry::get_callback_params_from_mcms(
        executing_callback_params,
    );
    assert!(*module_name.as_bytes() == b"mcms", EInvalidModuleName);
    assert!(*function_name.as_bytes() == b"timelock_block_function", EInvalidFunctionName);

    let (target, module_name, function_name) = deserialize_timelock_function_action(data);
    timelock_block_function(timelock, TIMELOCK_ROLE, target, module_name, function_name, ctx)
}

public fun execute_timelock_unblock_function(
    timelock: &mut Timelock,
    executing_callback_params: ExecutingCallbackParams,
    ctx: &mut TxContext,
) {
    let (_, module_name, function_name, data) = mcms_registry::get_callback_params_from_mcms(
        executing_callback_params,
    );
    assert!(*module_name.as_bytes() == b"mcms", EInvalidModuleName);
    assert!(*function_name.as_bytes() == b"timelock_unblock_function", EInvalidFunctionName);

    let (target, module_name, function_name) = deserialize_timelock_function_action(data);
    timelock_unblock_function(timelock, TIMELOCK_ROLE, target, module_name, function_name, ctx)
}

public fun execute_set_config(
    registry: &mut Registry,
    state: &mut MultisigState,
    executing_callback_params: ExecutingCallbackParams,
    ctx: &mut TxContext,
) {
    let (_, module_name, function_name, data) = mcms_registry::get_callback_params_from_mcms(
        executing_callback_params,
    );
    assert!(*module_name.as_bytes() == b"mcms", EInvalidModuleName);
    assert!(*function_name.as_bytes() == b"set_config", EInvalidFunctionName);

    let stream = &mut bcs_stream::new(data);
    let role_param = bcs_stream::deserialize_u8(stream);
    let signer_addresses = bcs_stream::deserialize_vector!(
        stream,
        |stream| { bcs_stream::deserialize_vector_u8(stream) },
    );
    let signer_groups = bcs_stream::deserialize_vector_u8(stream);
    let group_quorums = bcs_stream::deserialize_vector_u8(stream);
    let group_parents = bcs_stream::deserialize_vector_u8(stream);
    let clear_root = bcs_stream::deserialize_bool(stream);
    bcs_stream::assert_is_consumed(stream);

    let owner_cap = mcms_registry::borrow_owner_cap(registry);
    set_config(
        owner_cap,
        state,
        role_param,
        12,
        signer_addresses,
        signer_groups,
        group_quorums,
        group_parents,
        clear_root,
        ctx,
    );
}

/// Updates the multisig configuration, including signer addresses and group settings.
public entry fun set_config(
    _: &OwnerCap,
    state: &mut MultisigState,
    role: u8,
    chain_id: u256,
    signer_addresses: vector<vector<u8>>,
    signer_groups: vector<u8>,
    group_quorums: vector<u8>,
    group_parents: vector<u8>,
    clear_root: bool,
    _ctx: &TxContext,
) {
    assert!(
        signer_addresses.length() != 0
                    && signer_addresses.length() <= MAX_NUM_SIGNERS,
        EInvalidNumSigners,
    );
    assert!(signer_addresses.length() == signer_groups.length(), ESignerGroupsLenMismatch);
    assert!(group_quorums.length() == NUM_GROUPS, EInvalidGroupQuorumLen);
    assert!(group_parents.length() == NUM_GROUPS, EInvalidGroupParentsLen);

    // validate group structure
    // counts number of children of each group
    let mut group_children_counts = vector[];
    params::right_pad_vec(&mut group_children_counts, NUM_GROUPS);
    // first, we count the signers as children
    signer_groups.do_ref!(|group| {
        let group: u64 = *group as u64;
        assert!(group < NUM_GROUPS, EOutOfBoundsGroup);
        let count = group_children_counts.borrow_mut(group);
        *count = *count + 1;
    });

    // second, we iterate backwards so as to check each group and propagate counts from
    // child group to parent groups up the tree to the root
    let mut j = 0;
    while (j < NUM_GROUPS) {
        let i = NUM_GROUPS - j - 1;
        // ensure we have a well-formed group tree:
        // - the root should have itself as parent
        // - all other groups should have a parent group with a lower index
        let group_parent = group_parents[i] as u64;
        assert!(i == 0 || group_parent < i, EGroupTreeNotWellFormed);
        assert!(i != 0 || group_parent == 0, EGroupTreeNotWellFormed);

        let group_quorum = group_quorums[i];
        let disabled = group_quorum == 0;
        let group_children_count = group_children_counts[i];
        if (disabled) {
            // if group is disabled, ensure it has no children
            assert!(group_children_count == 0, ESignerInDisabledGroup);
        } else {
            // if group is enabled, ensure group quorum can be met
            assert!(group_children_count >= group_quorum, EOutOfBoundsGroupQuorum);

            // propagate children counts to parent group
            let count = group_children_counts.borrow_mut(group_parent);
            *count = *count + 1;
        };
        j = j + 1;
    };

    let multisig = borrow_multisig_mut(state, role);

    // remove old signer addresses
    multisig.signers = vec_map::empty();
    multisig.config.signers = vector[];

    // save group quorums and parents to timelock
    multisig.config.group_quorums = group_quorums;
    multisig.config.group_parents = group_parents;

    // check signer addresses are in increasing order and save signers to timelock
    // evm zero address (20 bytes of 0) is the smallest address possible
    let mut prev_signer_addr = vector[];
    let mut i = 0;
    while (i < signer_addresses.length()) {
        let signer_addr = signer_addresses[i];
        assert!(signer_addr.length() == 20, EInvalidSignerAddrLen);

        if (i > 0) {
            assert!(
                params::vector_u8_gt(&signer_addr, &prev_signer_addr),
                ESignerAddrMustBeIncreasing,
            );
        };

        let signer = Signer {
            addr: signer_addr,
            index: (i as u8),
            group: signer_groups[i],
        };
        multisig.signers.insert(signer_addr, signer);
        multisig.config.signers.push_back(signer);
        prev_signer_addr = signer_addr;

        i = i + 1;
    };

    if (clear_root) {
        // clearRoot is equivalent to overriding with a completely empty root
        let op_count = multisig.expiring_root_and_op_count.op_count;
        multisig.expiring_root_and_op_count =
            ExpiringRootAndOpCount {
                root: vector[],
                valid_until: 0,
                op_count,
            };
        multisig.root_metadata =
            RootMetadata {
                role,
                chain_id,
                multisig: mcms_registry::get_multisig_address(),
                pre_op_count: op_count,
                post_op_count: op_count,
                override_previous_root: true,
            };
    };

    event::emit(ConfigSet { role, config: multisig.config, is_root_cleared: clear_root });
}

public fun verify_merkle_proof(
    proof: vector<vector<u8>>,
    root: vector<u8>,
    leaf: vector<u8>,
): bool {
    let mut computed_hash = leaf;
    proof.do_ref!(|proof_element| {
        let (left, right) = if (params::vector_u8_gt(&computed_hash, proof_element)) {
            (*proof_element, computed_hash)
        } else {
            (computed_hash, *proof_element)
        };
        let mut hash_input: vector<u8> = left;
        hash_input.append(right);
        computed_hash = keccak256(&hash_input);
    });
    computed_hash == root
}

public fun compute_eth_message_hash(root: vector<u8>, valid_until: u64): vector<u8> {
    // abi.encode(root (bytes32), valid_until)
    let valid_until_bytes = params::encode_uint(valid_until, 32);
    assert!(root.length() == 32, EInvalidRootLen); // root should be 32 bytes
    let mut abi_encoded_params = root;
    abi_encoded_params.append(valid_until_bytes);

    // keccak256(abi_encoded_params)
    let hashed_encoded_params = keccak256(&abi_encoded_params);

    // ECDSA.toEthSignedMessageHash()
    let mut eth_msg_prefix = b"\x19Ethereum Signed Message:\n32";
    eth_msg_prefix.append(hashed_encoded_params);
    eth_msg_prefix
}

public fun hash_op_leaf(domain_separator: vector<u8>, op: Op): vector<u8> {
    let mut packed = vector[];
    packed.append(domain_separator);
    packed.append(bcs::to_bytes(&op.role));
    packed.append(bcs::to_bytes(&op.chain_id));
    packed.append(bcs::to_bytes(&op.multisig));
    packed.append(bcs::to_bytes(&op.nonce));
    packed.append(bcs::to_bytes(&op.to));
    packed.append(bcs::to_bytes(&op.module_name));
    packed.append(bcs::to_bytes(&op.function_name));
    packed.append(bcs::to_bytes(&op.data));
    keccak256(&packed)
}

fun hash_metadata_leaf(metadata: RootMetadata): vector<u8> {
    let mut packed = vector[];
    packed.append(MANY_CHAIN_MULTI_SIG_DOMAIN_SEPARATOR_METADATA);
    packed.append(bcs::to_bytes(&metadata.role));
    packed.append(bcs::to_bytes(&metadata.chain_id));
    packed.append(bcs::to_bytes(&metadata.multisig));
    packed.append(bcs::to_bytes(&metadata.pre_op_count));
    packed.append(bcs::to_bytes(&metadata.post_op_count));
    packed.append(bcs::to_bytes(&metadata.override_previous_root));
    keccak256(&packed)
}

fun deserialize_timelock_schedule_batch(
    data: vector<u8>,
): (
    vector<address>,
    vector<String>,
    vector<String>,
    vector<vector<u8>>,
    vector<u8>,
    vector<u8>,
    u64,
) {
    let stream = &mut bcs_stream::new(data);
    let targets = bcs_stream::deserialize_vector!(
        stream,
        |stream| bcs_stream::deserialize_address(stream),
    );
    let module_names = bcs_stream::deserialize_vector!(
        stream,
        |stream| bcs_stream::deserialize_string(stream),
    );
    let function_names = bcs_stream::deserialize_vector!(
        stream,
        |stream| bcs_stream::deserialize_string(stream),
    );
    let datas = bcs_stream::deserialize_vector!(
        stream,
        |stream| bcs_stream::deserialize_vector_u8(stream),
    );
    let predecessor = bcs_stream::deserialize_vector_u8(stream);
    let salt = bcs_stream::deserialize_vector_u8(stream);
    let delay = bcs_stream::deserialize_u64(stream);
    bcs_stream::assert_is_consumed(stream);

    (targets, module_names, function_names, datas, predecessor, salt, delay)
}

fun deserialize_timelock_execute_batch(
    data: vector<u8>,
): (vector<address>, vector<String>, vector<String>, vector<vector<u8>>, vector<u8>, vector<u8>) {
    let stream = &mut bcs_stream::new(data);
    let targets = bcs_stream::deserialize_vector!(
        stream,
        |stream| bcs_stream::deserialize_address(stream),
    );
    let module_names = bcs_stream::deserialize_vector!(
        stream,
        |stream| bcs_stream::deserialize_string(stream),
    );
    let function_names = bcs_stream::deserialize_vector!(
        stream,
        |stream| bcs_stream::deserialize_string(stream),
    );
    let datas = bcs_stream::deserialize_vector!(
        stream,
        |stream| bcs_stream::deserialize_vector_u8(stream),
    );
    let predecessor = bcs_stream::deserialize_vector_u8(stream);
    let salt = bcs_stream::deserialize_vector_u8(stream);
    bcs_stream::assert_is_consumed(stream);

    (targets, module_names, function_names, datas, predecessor, salt)
}

fun deserialize_timelock_bypasser_execute_batch(
    data: vector<u8>,
): (vector<address>, vector<String>, vector<String>, vector<vector<u8>>) {
    let stream = &mut bcs_stream::new(data);
    let targets = bcs_stream::deserialize_vector!(
        stream,
        |stream| bcs_stream::deserialize_address(stream),
    );
    let module_names = bcs_stream::deserialize_vector!(
        stream,
        |stream| bcs_stream::deserialize_string(stream),
    );
    let function_names = bcs_stream::deserialize_vector!(
        stream,
        |stream| bcs_stream::deserialize_string(stream),
    );
    let datas = bcs_stream::deserialize_vector!(
        stream,
        |stream| bcs_stream::deserialize_vector_u8(stream),
    );
    bcs_stream::assert_is_consumed(stream);

    (targets, module_names, function_names, datas)
}

fun deserialize_timelock_cancel(data: vector<u8>): vector<u8> {
    let stream = &mut bcs_stream::new(data);
    let id = bcs_stream::deserialize_vector_u8(stream);
    bcs_stream::assert_is_consumed(stream);
    id
}

fun deserialize_timelock_update_min_delay(data: vector<u8>): u64 {
    let stream = &mut bcs_stream::new(data);
    let new_min_delay = bcs_stream::deserialize_u64(stream);
    bcs_stream::assert_is_consumed(stream);
    new_min_delay
}

fun deserialize_timelock_function_action(data: vector<u8>): (address, String, String) {
    let stream = &mut bcs_stream::new(data);
    let target = bcs_stream::deserialize_address(stream);
    let module_name = bcs_stream::deserialize_string(stream);
    let function_name = bcs_stream::deserialize_string(stream);
    bcs_stream::assert_is_consumed(stream);

    (target, module_name, function_name)
}

public fun seen_signed_hashes(state: &MultisigState, role: u8): VecMap<vector<u8>, bool> {
    borrow_multisig(state, role).seen_signed_hashes
}

public fun expiring_root_and_op_count(state: &MultisigState, role: u8): (vector<u8>, u64, u64) {
    let multisig = borrow_multisig(state, role);
    (
        multisig.expiring_root_and_op_count.root,
        multisig.expiring_root_and_op_count.valid_until,
        multisig.expiring_root_and_op_count.op_count,
    )
}

public fun root_metadata(multisig: &Multisig): RootMetadata {
    multisig.root_metadata
}

public fun get_root_metadata(state: &MultisigState, role: u8): RootMetadata {
    let multisig = borrow_multisig(state, role);
    multisig.root_metadata
}

public fun get_op_count(state: &MultisigState, role: u8): u64 {
    let multisig = borrow_multisig(state, role);
    multisig.expiring_root_and_op_count.op_count
}

public fun get_root(state: &MultisigState, role: u8): (vector<u8>, u64) {
    let multisig = borrow_multisig(state, role);
    (multisig.expiring_root_and_op_count.root, multisig.expiring_root_and_op_count.valid_until)
}

public fun get_config(state: &MultisigState, role: u8): Config {
    borrow_multisig(state, role).config
}

public fun num_groups(): u64 {
    NUM_GROUPS
}

public fun max_num_signers(): u64 {
    MAX_NUM_SIGNERS
}

public fun bypasser_role(): u8 {
    BYPASSER_ROLE
}

public fun canceller_role(): u8 {
    CANCELLER_ROLE
}

public fun proposer_role(): u8 {
    PROPOSER_ROLE
}

public fun timelock_role(): u8 {
    TIMELOCK_ROLE
}

public fun is_valid_role(role: u8): bool {
    role < MAX_ROLE
}

public fun zero_hash(): vector<u8> {
    ZERO_HASH
}

fun borrow_multisig(state: &MultisigState, role: u8): &Multisig {
    if (role == BYPASSER_ROLE) {
        return &state.bypasser
    } else if (role == CANCELLER_ROLE) {
        return &state.canceller
    } else if (role == PROPOSER_ROLE) {
        return &state.proposer
    } else {
        abort EInvalidRole
    }
}

fun borrow_multisig_mut(state: &mut MultisigState, role: u8): &mut Multisig {
    if (role == BYPASSER_ROLE) {
        return &mut state.bypasser
    } else if (role == CANCELLER_ROLE) {
        return &mut state.canceller
    } else if (role == PROPOSER_ROLE) {
        return &mut state.proposer
    } else {
        abort EInvalidRole
    }
}

public fun role(root_metadata: &RootMetadata): u8 {
    root_metadata.role
}

public fun chain_id(root_metadata: &RootMetadata): u256 {
    root_metadata.chain_id
}

public fun root_metadata_multisig(root_metadata: &RootMetadata): address {
    root_metadata.multisig
}

public fun pre_op_count(root_metadata: &RootMetadata): u64 {
    root_metadata.pre_op_count
}

public fun post_op_count(root_metadata: &RootMetadata): u64 {
    root_metadata.post_op_count
}

public fun override_previous_root(root_metadata: &RootMetadata): bool {
    root_metadata.override_previous_root
}

public fun config_signers(config: &Config): vector<Signer> {
    config.signers
}

public fun config_group_quorums(config: &Config): vector<u8> {
    config.group_quorums
}

public fun config_group_parents(config: &Config): vector<u8> {
    config.group_parents
}

// =======================================================================================
// |                                 Timelock Implementation                              |
// =======================================================================================

public struct Timelock has key, store {
    id: UID,
    min_delay: u64,
    /// hashed batch of hashed calls -> timestamp
    timestamps: Table<vector<u8>, u64>,
    /// blocked functions
    blocked_functions: VecSet<Function>,
}

public struct Call has copy, drop, store {
    function: Function,
    data: vector<u8>,
}

public struct Function has copy, drop, store {
    target: address,
    module_name: String,
    function_name: String,
}

public struct TimelockInitialized has copy, drop {
    min_delay: u64,
}

public struct BypasserCallInitiated has copy, drop {
    index: u64,
    target: address,
    module_name: String,
    function_name: String,
    data: vector<u8>,
}

public struct Cancelled has copy, drop {
    id: vector<u8>,
}

public struct CallScheduled has copy, drop {
    id: vector<u8>,
    index: u64,
    target: address,
    module_name: String,
    function_name: String,
    data: vector<u8>,
    predecessor: vector<u8>,
    salt: vector<u8>,
    delay: u64,
}

public struct CallInitiated has copy, drop {
    id: vector<u8>,
    index: u64,
    target: address,
    module_name: String,
    function_name: String,
    data: vector<u8>,
}

public struct UpdateMinDelay has copy, drop {
    old_min_delay: u64,
    new_min_delay: u64,
}

public struct FunctionBlocked has copy, drop {
    target: address,
    module_name: String,
    function_name: String,
}

public struct FunctionUnblocked has copy, drop {
    target: address,
    module_name: String,
    function_name: String,
}

/// Schedule a batch of calls to be executed after a delay.
/// This function can only be called by PROPOSER or TIMELOCK role.
fun timelock_schedule_batch(
    timelock: &mut Timelock,
    clock: &Clock,
    role: u8,
    targets: vector<address>,
    module_names: vector<String>,
    function_names: vector<String>,
    datas: vector<vector<u8>>,
    predecessor: vector<u8>,
    salt: vector<u8>,
    delay: u64,
    _ctx: &mut TxContext,
) {
    assert!(role == PROPOSER_ROLE || role == TIMELOCK_ROLE, ENotAuthorizedRole);

    let calls = create_calls(targets, module_names, function_names, datas);
    let id = hash_operation_batch(calls, predecessor, salt);

    timelock_schedule(timelock, clock, id, delay);

    let mut i = 0;
    while (i < calls.length()) {
        assert_not_blocked(timelock, &calls[i].function);
        event::emit(CallScheduled {
            id,
            index: i,
            target: calls[i].function.target,
            module_name: calls[i].function.module_name,
            function_name: calls[i].function.function_name,
            data: calls[i].data,
            predecessor,
            salt,
            delay,
        });
        i = i + 1;
    };
}

fun timelock_schedule(timelock: &mut Timelock, clock: &Clock, id: vector<u8>, delay: u64) {
    assert!(!timelock_is_operation_internal(timelock, id), EOperationAlreadyScheduled);
    assert!(delay >= timelock.min_delay, EInsufficientDelay);

    let timestamp = clock.timestamp_ms() + timelock.min_delay + delay;
    timelock.timestamps.add(id, timestamp);
}

fun timelock_before_call(
    timelock: &Timelock,
    clock: &Clock,
    id: vector<u8>,
    predecessor: vector<u8>,
) {
    assert!(timelock_is_operation_ready(timelock, clock, id), EOperationNotReady);
    assert!(
        predecessor == ZERO_HASH || timelock_is_operation_done(timelock, predecessor),
        EMissingDependency,
    );
}

fun timelock_after_call(
    timelock: &mut Timelock,
    clock: &Clock,
    id: vector<u8>,
    _ctx: &mut TxContext,
) {
    assert!(timelock_is_operation_ready(timelock, clock, id), EOperationNotReady);
    *timelock.timestamps.borrow_mut(id) = DONE_TIMESTAMP;
}

/// Anyone can call this as it checks if the operation was scheduled by a bypasser or proposer.
public fun timelock_execute_batch(
    timelock: &mut Timelock,
    clock: &Clock,
    targets: vector<address>,
    module_names: vector<String>,
    function_names: vector<String>,
    datas: vector<vector<u8>>,
    predecessor: vector<u8>,
    salt: vector<u8>,
    _ctx: &mut TxContext,
): vector<ExecutingCallbackParams> {
    let calls = create_calls(targets, module_names, function_names, datas);
    let id = hash_operation_batch(calls, predecessor, salt);

    timelock_before_call(timelock, clock, id, predecessor);

    let mut calls_to_execute = vector[];
    let mut i = 0;

    while (i < calls.length()) {
        let function = calls[i].function;
        let target = function.target;
        let module_name = function.module_name;
        let function_name = function.function_name;
        let data = calls[i].data;

        let params = mcms_registry::create_executing_callback_params(
            target,
            module_name,
            function_name,
            data,
        );
        calls_to_execute.push_back(params);

        event::emit(CallInitiated { id, index: i, target, module_name, function_name, data });

        i = i + 1;
    };

    // Timestamps can be safely updated as `vector<ExecutingCallbackParams>` are returned to the caller
    // These must be consumed by relevant modules
    timelock_after_call(timelock, clock, id, _ctx);

    // Must call module's mcms_entrypoint with `calls_to_execute`
    calls_to_execute
}

fun timelock_bypasser_execute_batch(
    role: u8,
    targets: vector<address>,
    module_names: vector<String>,
    function_names: vector<String>,
    datas: vector<vector<u8>>,
    _ctx: &mut TxContext,
): vector<ExecutingCallbackParams> {
    assert!(role == BYPASSER_ROLE || role == TIMELOCK_ROLE, ENotAuthorizedRole);

    let mut calls_to_execute = vector[];
    let mut i = 0;

    while (i < targets.length()) {
        let target = targets[i];
        let module_name = module_names[i];
        let function_name = function_names[i];
        let data = datas[i];

        let params = mcms_registry::create_executing_callback_params(
            target,
            module_name,
            function_name,
            data,
        );
        calls_to_execute.push_back(params);

        event::emit(BypasserCallInitiated { index: i, target, module_name, function_name, data });

        i = i + 1;
    };

    // Must call module's mcms_entrypoint with each `ExecutingCallbackParams` in `calls_to_execute`
    calls_to_execute
}

// =======================================================================================
// |                                 Code Ownership Transfer                              |
// =======================================================================================

// Each package defines an `OwnerCap` object
// Ownership of `OwnerCap` dictates the owner of the code object
// Transferring ownership of `OwnerCap` must be done by calling `mcms_entrypoint` in CCIP `ownable.move`
// This replaces `timelock_dispatch_to_registry`

fun timelock_cancel(timelock: &mut Timelock, role: u8, id: vector<u8>, _ctx: &mut TxContext) {
    assert!(role == CANCELLER_ROLE || role == TIMELOCK_ROLE, ENotAuthorizedRole);
    assert!(timelock_is_operation_pending(timelock, id), EOperationCannotBeCancelled);

    timelock.timestamps.remove(id);
    event::emit(Cancelled { id });
}

fun timelock_update_min_delay(
    timelock: &mut Timelock,
    role: u8,
    new_min_delay: u64,
    _ctx: &mut TxContext,
) {
    assert!(role == TIMELOCK_ROLE, EInvalidRole);

    let old_min_delay = timelock.min_delay;
    timelock.min_delay = new_min_delay;

    event::emit(UpdateMinDelay { old_min_delay, new_min_delay });
}

fun timelock_block_function(
    timelock: &mut Timelock,
    role: u8,
    target: address,
    module_name: String,
    function_name: String,
    _ctx: &mut TxContext,
) {
    assert!(role == TIMELOCK_ROLE, ENotTimeLockRole);

    let mut already_blocked = false;
    let new_function = Function { target, module_name, function_name };
    let mut i = 0;

    while (i < timelock.blocked_functions.size()) {
        let blocked_function = timelock.blocked_functions.keys()[i];
        if (equals(&new_function, &blocked_function)) {
            already_blocked = true;
            break
        };
        i = i + 1;
    };

    if (!already_blocked) {
        timelock.blocked_functions.insert(new_function);
        event::emit(FunctionBlocked { target, module_name, function_name });
    };
}

fun timelock_unblock_function(
    timelock: &mut Timelock,
    role: u8,
    target: address,
    module_name: String,
    function_name: String,
    _ctx: &mut TxContext,
) {
    assert!(role == TIMELOCK_ROLE, ENotTimeLockRole);

    let function_to_unblock = Function { target, module_name, function_name };
    let mut i = 0;

    while (i < timelock.blocked_functions.size()) {
        let blocked_function = timelock.blocked_functions.keys()[i];
        if (equals(&function_to_unblock, &blocked_function)) {
            timelock.blocked_functions.remove(&blocked_function);
            event::emit(FunctionUnblocked { target, module_name, function_name });
            break
        };
        i = i + 1;
    };
}

fun assert_not_blocked(timelock: &Timelock, function: &Function) {
    let mut i = 0;
    while (i < timelock.blocked_functions.size()) {
        let blocked_function = timelock.blocked_functions.keys()[i];
        if (equals(function, &blocked_function)) {
            abort EFunctionBlocked
        };
        i = i + 1;
    };
}

public fun timelock_get_blocked_function(timelock: &Timelock, index: u64): Function {
    assert!(index < timelock.blocked_functions.size(), EInvalidIndex);
    timelock.blocked_functions.keys()[index]
}

public fun timelock_is_operation(timelock: &Timelock, id: vector<u8>): bool {
    timelock_is_operation_internal(timelock, id)
}

fun timelock_is_operation_internal(timelock: &Timelock, id: vector<u8>): bool {
    timelock.timestamps.contains(id) && *timelock.timestamps.borrow(id) > 0
}

public fun timelock_is_operation_pending(timelock: &Timelock, id: vector<u8>): bool {
    timelock.timestamps.contains(id)
                && *timelock.timestamps.borrow(id) > DONE_TIMESTAMP
}

public fun timelock_is_operation_ready(timelock: &Timelock, clock: &Clock, id: vector<u8>): bool {
    if (!timelock.timestamps.contains(id)) {
        return false
    };

    let timestamp_value = *timelock.timestamps.borrow(id);
    timestamp_value > DONE_TIMESTAMP && timestamp_value <= clock.timestamp_ms()
}

public fun timelock_is_operation_done(timelock: &Timelock, id: vector<u8>): bool {
    timelock.timestamps.contains(id)
                && *timelock.timestamps.borrow(id) == DONE_TIMESTAMP
}

public fun timelock_get_timestamp(timelock: &Timelock, id: vector<u8>): u64 {
    if (timelock.timestamps.contains(id)) {
        *timelock.timestamps.borrow(id)
    } else { 0 }
}

public fun timelock_min_delay(timelock: &Timelock): u64 {
    timelock.min_delay
}

public fun timelock_get_blocked_functions(timelock: &Timelock): vector<Function> {
    let mut blocked_functions = vector[];
    let contents = timelock.blocked_functions.keys();
    let mut i = 0;
    while (i < contents.length()) {
        blocked_functions.push_back(contents[i]);
        i = i + 1;
    };
    blocked_functions
}

public fun timelock_get_blocked_functions_count(timelock: &Timelock): u64 {
    timelock.blocked_functions.size()
}

public fun create_calls(
    targets: vector<address>,
    module_names: vector<String>,
    function_names: vector<String>,
    datas: vector<vector<u8>>,
): vector<Call> {
    let len = targets.length();
    assert!(
        len == module_names.length()
                    && len == function_names.length()
                    && len == datas.length(),
        EInvalidParameters,
    );

    let mut calls = vector[];
    let mut i = 0;
    while (i < len) {
        let target = targets[i];
        let module_name = module_names[i];
        let function_name = function_names[i];
        let data = datas[i];
        let function = Function { target, module_name, function_name };
        let call = Call { function, data };
        calls.push_back(call);
        i = i + 1;
    };

    calls
}

public fun hash_operation_batch(
    calls: vector<Call>,
    predecessor: vector<u8>,
    salt: vector<u8>,
): vector<u8> {
    let mut packed = vector[];
    packed.append(bcs::to_bytes(&calls));
    packed.append(predecessor);
    packed.append(salt);
    keccak256(&packed)
}

fun equals(fn1: &Function, fn2: &Function): bool {
    fn1.target == fn2.target
                && fn1.module_name.as_bytes() == fn2.module_name.as_bytes()
                && fn1.function_name.as_bytes() == fn2.function_name.as_bytes()
}

public fun signer_view(signer_: &Signer): (vector<u8>, u8, u8) {
    (signer_.addr, signer_.index, signer_.group)
}

public fun function_name(function: Function): String {
    function.function_name
}

public fun module_name(function: Function): String {
    function.module_name
}

public fun target(function: Function): address {
    function.target
}

public fun data(call: Call): vector<u8> {
    call.data
}

fun get_timestamp_seconds(clock: &Clock): u64 {
    clock.timestamp_ms() / 1000
}

// ===================== TESTS =====================

#[test_only]
public fun test_init(ctx: &mut TxContext) {
    init(MCMS {}, ctx)
}

#[test_only]
public fun test_timelock_schedule_batch(
    timelock: &mut Timelock,
    clock: &Clock,
    role: u8,
    targets: vector<address>,
    module_names: vector<String>,
    function_names: vector<String>,
    datas: vector<vector<u8>>,
    predecessor: vector<u8>,
    salt: vector<u8>,
    delay: u64,
    ctx: &mut TxContext,
) {
    timelock_schedule_batch(
        timelock,
        clock,
        role,
        targets,
        module_names,
        function_names,
        datas,
        predecessor,
        salt,
        delay,
        ctx,
    )
}

#[test_only]
public fun test_ecdsa_recover_evm_addr(msg: vector<u8>, signature: vector<u8>): vector<u8> {
    ecdsa_recover_evm_addr(msg, signature)
}

#[test_only]
public fun test_compute_eth_message_hash(root: vector<u8>, valid_until: u64): vector<u8> {
    compute_eth_message_hash(root, valid_until)
}

#[test_only]
public fun test_set_hash_seen(state: &mut MultisigState, role: u8, hash: vector<u8>, seen: bool) {
    let multisig = borrow_multisig_mut(state, role);
    multisig.seen_signed_hashes.insert(hash, seen);
}

#[test_only]
public fun test_set_expiring_root_and_op_count(
    state: &mut MultisigState,
    role: u8,
    root: vector<u8>,
    valid_until: u64,
    op_count: u64,
) {
    let multisig = borrow_multisig_mut(state, role);
    multisig.expiring_root_and_op_count.root = root;
    multisig.expiring_root_and_op_count.valid_until = valid_until;
    multisig.expiring_root_and_op_count.op_count = op_count;
}

#[test_only]
public fun test_set_root_metadata(
    state: &mut MultisigState,
    role: u8,
    chain_id: u256,
    multisig_addr: address,
    pre_op_count: u64,
    post_op_count: u64,
    override_previous_root: bool,
) {
    let multisig = borrow_multisig_mut(state, role);
    multisig.root_metadata.role = role;
    multisig.root_metadata.chain_id = chain_id;
    multisig.root_metadata.multisig = multisig_addr;
    multisig.root_metadata.pre_op_count = pre_op_count;
    multisig.root_metadata.post_op_count = post_op_count;
    multisig.root_metadata.override_previous_root = override_previous_root;
}
