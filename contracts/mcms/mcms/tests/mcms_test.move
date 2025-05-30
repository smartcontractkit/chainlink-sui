#[test_only]
module mcms::mcms_test;

use mcms::mcms::{
    Self,
    MultisigState,
    Timelock,
};
use mcms::mcms_account::{Self, AccountState, OwnerCap};
use mcms::mcms_deployer::{Self, DeployerState};
use mcms::mcms_registry::{Self, Registry};
use mcms::params;
use std::string;
use std::vector;
use sui::ecdsa_k1;
use sui::hash::keccak256;
use sui::hex;
use sui::test_scenario::{Self as ts, Scenario};

const OWNER: address = @0x123;
const SENDER: address = @0x456;
const MCMS_TEST_MODULE_NAME: vector<u8> = b"mcms_test";

// keccak256("MANY_CHAIN_MULTI_SIG_DOMAIN_SEPARATOR_OP_APTOS")
const MANY_CHAIN_MULTI_SIG_DOMAIN_SEPARATOR_OP: vector<u8> =
    x"e5a6d1256b00d7ec22512b6b60a3f4d75c559745d2dbf309f77b8b756caabe14";

const CHAIN_ID: u256 = 4;
const TIMESTAMP: u64 = 1744315405;

const MIN_DELAY: u64 = 3600; // 1 hour delay
const TEST_TARGET_ADDRESS: address = @0xabc;
const TEST_SALT: vector<u8> = x"1234567890abcdef";
const TEST_PREDECESSOR: vector<u8> = x"";

// Proposer signers from the logs (already in ascending order)
const PROPOSER_ADDR1: vector<u8> = x"5916431f0ea809587757df994233861e1271be55";
const PROPOSER_ADDR2: vector<u8> = x"8803c3ed076e57d51e28301933418094bd961cc5";
const PROPOSER_ADDR3: vector<u8> = x"8950e6c6832c9b0591801418684d27b2853b2c74";

// test config: 2-of-3 multisig
const SIGNER_GROUPS: vector<u8> = vector[0, 0, 0];

const GROUP_QUORUMS: vector<u8> = vector[
    2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
];

const GROUP_PARENTS: vector<u8> = vector[
    0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
];

public struct Env {
    scenario: ts::Scenario,
    state: MultisigState,
    timelock: Timelock,
    registry: Registry,
    account_state: AccountState,
    deployer_state: DeployerState,
    clock: sui::clock::Clock,
}

public struct SetRootArgs has drop {
    role: u8,
    root: vector<u8>,
    valid_until: u64,
    chain_id: u256,
    multisig: vector<u8>,
    pre_op_count: u64,
    post_op_count: u64,
    override_previous_root: bool,
    metadata_proof: vector<vector<u8>>,
    signatures: vector<vector<u8>>,
}

const ROOT: vector<u8> = x"207e468568a65fc396a84da1291685cc64756b7ade81eababdcdc32ba6ba26da";
const VALID_UNTIL: u64 = 1747669800045;
const PRE_OP_COUNT: u64 = 0;
const POST_OP_COUNT: u64 = 4;

const METADATA_PROOF: vector<vector<u8>> = vector[
    x"66cf50cb9a50c740313fd0f889b676af60d35ef700711d94df6eeff3f1ba66c2",
    x"951f1094081a858642cc6635f0885317828a7fddddd00668391c50f1e9e1bb66",
    x"597116801e22b18150f2abc4ca2ecd63e147bb67e24e4b5f900d49b909e1919f",
];

const OP1_PROOF: vector<vector<u8>> = vector[
    x"a619565e90c1c564293b59b344ed0e12ed06eafb3c45b70baf6fdf299a046297", // metadata hash
    x"951f1094081a858642cc6635f0885317828a7fddddd00668391c50f1e9e1bb66", // sibling at level 1
    x"597116801e22b18150f2abc4ca2ecd63e147bb67e24e4b5f900d49b909e1919f", // sibling at level 2
];

// The OPs contained are
// {
// 			Target:      mcmsAccount,
// 			ModuleName:  "mcms_account",
// 			Function:    "accept_ownership",
// 			Data:        []byte{},
// 			Delay:       1,
// 			Predecessor: []byte{},
// 			Salt:        []byte{},
// 		},
// 		{
// 			Target:      mcmsAccount,
// 			ModuleName:  "mcms_deployer",
// 			Function:    "stage_code_chunk_and_publish_to_object",
// 			Data:        stageCodeChunkAndPublishToAccountBytes,
// 			Delay:       1,
// 			Predecessor: []byte{},
// 			Salt:        []byte{},
// 		},
// 		{
// 			Target:      userModuleAccount,
// 			ModuleName:  "mcms_user",
// 			Function:    "function_one",
// 			Data:        functionOneParamBytes,
// 			Delay:       1,
// 			Predecessor: []byte{},
// 			Salt:        []byte{},
// 		},
// 		{
// 			Target:      userModuleAccount,
// 			ModuleName:  "mcms_user",
// 			Function:    "function_two",
// 			Data:        functionTwoParamBytes,
// 			Delay:       1,
// 			Predecessor: []byte{},
// 			Salt:        []byte{},
// 		},
const LEAVES: vector<vector<u8>> = vector[
    x"a619565e90c1c564293b59b344ed0e12ed06eafb3c45b70baf6fdf299a046297", // metadata hash
    x"66cf50cb9a50c740313fd0f889b676af60d35ef700711d94df6eeff3f1ba66c2", // op1 hash
    x"2feec0e3a232c5c847874246203e62c43db473fe85245095122e166be9114e13", // op2 hash
    x"411a4726f8a920fc0a814bd9897a06f3dd0f1c799a047deaa6469f105f5a6705", // op3 hash
    x"cb4dffef33843b197cd33346d3339d8432b14789504167c63fb9f74a73baaea5", // op4 hash
];

const SIGNATURES: vector<vector<u8>> = vector[
    x"72398e2f325e707217fa8108a08c126f49f4144c30c7e93896d139c9f1d9468c30424b060c19fa5c7820a17b57badb19375207c787878533834618688a4780581c",
    x"9bb8ba839f9152cdc61556fcc70b0ebcb4d442654263a3d1c323e1eed85ebc6016e87c8e59f12d850b9d6b789ccafd93f19dcf65eb6fc75fd4351d5970214c1d1c",
    x"225739c80de11d50f3dca8fbb8288881abad17439690abd8eee32d48ff2f6dd204c48aff2fac4dfe8ce7176fd12d2633d7892bfb3d3f3cfeb00352773fe55c8c1b",
];

const OP1_NONCE: u64 = 0;
const OP1_DATA: vector<u8> =
    x"01a969156fce9a4f08bcdc07b90f338efc630bff8dfa8340500cb6414aca762a4e010c6d636d735f6163636f756e7401106163636570745f6f776e6572736869700100200000000000000000000000000000000000000000000000000000000000000000000000000000000000";

fun default_set_root_args(override_previous_root: bool): SetRootArgs {
    SetRootArgs {
        role: mcms::proposer_role(),
        root: ROOT,
        valid_until: VALID_UNTIL,
        chain_id: CHAIN_ID,
        multisig: @mcms.to_bytes(),
        pre_op_count: PRE_OP_COUNT,
        post_op_count: POST_OP_COUNT,
        override_previous_root,
        metadata_proof: METADATA_PROOF,
        signatures: SIGNATURES,
    }
}

fun call_set_root(env: &mut Env, args: SetRootArgs) {
    mcms::set_root(
        &mut env.state,
        &env.clock,
        args.role,
        args.root,
        args.valid_until,
        args.chain_id,
        args.multisig,
        args.pre_op_count,
        args.post_op_count,
        args.override_previous_root,
        args.metadata_proof,
        args.signatures,
        env.scenario.ctx(),
    );
}

fun setup(): Env {
    let mut scenario = ts::begin(OWNER);
    let ctx = scenario.ctx();

    let mut clock = sui::clock::create_for_testing(ctx);
    clock.set_for_testing(1_000_000_000);

    mcms_account::test_init(ctx);
    mcms_registry::test_init(ctx);
    mcms_deployer::test_init(ctx);
    mcms::test_init(ctx);

    scenario.next_tx(OWNER);

    let state = ts::take_shared<MultisigState>(&scenario);
    let timelock = ts::take_shared<Timelock>(&scenario);
    let registry = ts::take_shared<Registry>(&scenario);
    let account_state = ts::take_shared<AccountState>(&scenario);
    let deployer_state = ts::take_shared<DeployerState>(&scenario);

    Env {
        scenario,
        state,
        timelock,
        registry,
        account_state,
        deployer_state,
        clock,
    }
}

#[test]
#[expected_failure(abort_code = mcms::EAlreadySeenHash, location = mcms)]
fun test_set_root__already_seen_hash() {
    let mut env = setup();
    let owner_cap = ts::take_from_sender<OwnerCap>(&env.scenario);

    let role = mcms::proposer_role();
    mcms::set_config(
        &owner_cap,
        &mut env.state,
        role,
        CHAIN_ID,
        vector[PROPOSER_ADDR1, PROPOSER_ADDR2, PROPOSER_ADDR3],
        SIGNER_GROUPS,
        GROUP_QUORUMS,
        GROUP_PARENTS,
        false,
        env.scenario.ctx(),
    );

    let signed_hash = mcms::compute_eth_message_hash(ROOT, VALID_UNTIL);
    mcms::test_set_hash_seen(
        &mut env.state,
        role,
        signed_hash,
        true,
    );

    call_set_root(&mut env, default_set_root_args(false));

    ts::return_to_sender(&env.scenario, owner_cap);
    env.destroy();
}

// Need valid proofs to test this
// #[test]
// #[expected_failure(abort_code = mcms::EValidUntilExpired, location = mcms)]
// public fun test_set_root__valid_until_expired() {
//     let mut env = setup();
//     let mut set_root_args = default_set_root_args(false);
//     set_root_args.valid_until = TIMESTAMP - 1; // set valid_until to a time in the past
//     call_set_root(&mut env, set_root_args);

//     env.destroy()
// }

#[test]
#[expected_failure(abort_code = mcms::EInvalidRootLen, location = mcms)]
fun test_set_root__invalid_root_len() {
    let mut env = setup();
    let invalid_root = x"8ad6edb34398f637ca17e46b0b51ce50e18f56287aa0bf728ae3b5c4119c16";
    let mut set_root_args = default_set_root_args(false);
    set_root_args.root = invalid_root;
    call_set_root(&mut env, set_root_args);

    env.destroy()
}

#[test]
#[expected_failure(abort_code = mcms::EWrongMultisig, location = mcms)]
fun test_set_root__invalid_multisig_addr() {
    let mut env = setup();
    let mut set_root_args = default_set_root_args(false);
    set_root_args.multisig = @0x999.to_bytes();
    call_set_root(&mut env, set_root_args);

    env.destroy()
}

#[test]
#[expected_failure(abort_code = mcms::EPendingOps, location = mcms)]
public fun test_set_root__pending_ops() {
    let mut env = setup();
    let role = mcms::proposer_role();
    mcms::test_set_expiring_root_and_op_count(&mut env.state, role, ROOT, VALID_UNTIL, 1);
    mcms::test_set_root_metadata(
        &mut env.state,
        role,
        CHAIN_ID,
        @mcms.to_bytes(),
        0,
        2, // 1 more than the current op_count
        false,
    );

    call_set_root(&mut env, default_set_root_args(false));

    env.destroy()
}

#[test]
#[expected_failure(abort_code = mcms::EProofCannotBeVerified, location = mcms)]
public fun test_set_root__override_previous_root() {
    let mut env = setup();
    let owner_cap = ts::take_from_sender<OwnerCap>(&env.scenario);
    mcms::set_config(
        &owner_cap,
        &mut env.state,
        mcms::proposer_role(),
        CHAIN_ID,
        vector[PROPOSER_ADDR1, PROPOSER_ADDR2, PROPOSER_ADDR3],
        SIGNER_GROUPS,
        GROUP_QUORUMS,
        GROUP_PARENTS,
        false,
        env.scenario.ctx(),
    );
    let mut set_root_args = default_set_root_args(false);
    set_root_args.post_op_count = 20;
    // Change the post_op_count to a value that is not equal to the proof's post_op_count
    call_set_root(&mut env, set_root_args);

    ts::return_to_sender(&env.scenario, owner_cap);
    env.destroy();
}

#[test]
#[expected_failure(abort_code = mcms::EWrongPreOpCount, location = mcms)]
public fun test_set_root__wrong_pre_op_count() {
    let mut env = setup();
    let mut set_root_args = default_set_root_args(false);
    set_root_args.pre_op_count = 1; // wrong pre op count, should equal op count (0)
    call_set_root(&mut env, set_root_args);
    env.destroy();
}

#[test]
#[expected_failure(abort_code = mcms::EWrongPostOpCount, location = mcms)]
public fun test_set_root__wrong_post_op_count() {
    let mut env = setup();
    let role = mcms::proposer_role();
    mcms::test_set_expiring_root_and_op_count(&mut env.state, role, ROOT, VALID_UNTIL, 1);
    mcms::test_set_root_metadata(
        &mut env.state,
        role,
        CHAIN_ID,
        @mcms.to_bytes(),
        0,
        1,
        false,
    );

    let mut set_root_args = default_set_root_args(false);
    set_root_args.pre_op_count = PRE_OP_COUNT + 1; // correct pre op count after state updates
    set_root_args.post_op_count = PRE_OP_COUNT; // post op count should be >= pre op count
    call_set_root(&mut env, set_root_args);

    env.destroy();
}

#[test]
#[expected_failure(abort_code = mcms::EProofCannotBeVerified, location = mcms)]
public fun test_set_root__empty_metadata_proof() {
    let mut env = setup();
    let mut set_root_args = default_set_root_args(false);
    set_root_args.metadata_proof = vector[]; // empty proof
    call_set_root(&mut env, set_root_args);
    env.destroy();
}

#[test]
#[expected_failure(abort_code = mcms::EProofCannotBeVerified, location = mcms)]
public fun test_set_root__metadata_not_consistent_with_proof() {
    let mut env = setup();
    let mut set_root_args = default_set_root_args(false);
    set_root_args.post_op_count = POST_OP_COUNT + 1; // post op count modified
    call_set_root(&mut env, set_root_args);
    env.destroy();
}

// ============== Need valid proofs to test these =================

// #[test]
// #[expected_failure(abort_code = mcms::EMissingConfig, location = mcms)]
// fun test_set_root__config_not_set() {
//     let mut env = setup();
//     let mut set_root_args = default_set_root_args(false);
//     set_root_args.signatures = vector[]; // no signatures
//     call_set_root(&mut env, set_root_args);
//     env.destroy();
// }

// #[test]
// #[expected_failure(abort_code = mcms::ESignerAddrMustBeIncreasing, location = mcms)]
// fun test_set_root__out_of_order_signatures() {
//     let mut env = setup();
//     let owner_cap = ts::take_from_sender<OwnerCap>(&env.scenario);
//     let role = mcms::proposer_role();
//         mcms::set_config(
//         &owner_cap,
//         &mut env.state,
//         role,
//         CHAIN_ID,
//         vector[PROPOSER_ADDR1, PROPOSER_ADDR2, PROPOSER_ADDR3],
//         SIGNER_GROUPS,
//         GROUP_QUORUMS,
//         GROUP_PARENTS,
//         false,
//         env.scenario.ctx(),
//     );
//     let mut set_root_args = default_set_root_args(false);
//     let sig0 = set_root_args.signatures[0];
//     let sig1 = set_root_args.signatures[1];
//     let sig2 = set_root_args.signatures[2];
//     set_root_args.signatures = vector[sig0, sig2, sig1]; // shuffle signature order
//     call_set_root(&mut env, set_root_args);

//     ts::return_to_sender(&env.scenario, owner_cap);
//     env.destroy();
// }

// #[test]
// #[expected_failure(abort_code = mcms::EInvalidSigner, location = mcms)]
// fun test_set_root__signature_from_invalid_signer() {
//     let mut env = setup();
//     let owner_cap = ts::take_from_sender<OwnerCap>(&env.scenario);
//     let role = mcms::proposer_role();
//     mcms::set_config(
//         &owner_cap,
//         &mut env.state,
//         role,
//         CHAIN_ID,
//         vector[PROPOSER_ADDR1, PROPOSER_ADDR2, PROPOSER_ADDR3],
//         SIGNER_GROUPS,
//         GROUP_QUORUMS,
//         GROUP_PARENTS,
//         false,
//         env.scenario.ctx(),
//     );
//     let mut set_root_args = default_set_root_args(false);
//     let invalid_signer_sig =
//         x"bb7f7e44b8d9c8f978c255c7efd6abb64e8fa9a33dcb6db2e2203d8aacd51dd471113ca6c8d1ed56bb0395f0bef0daf2fae6ef2cb5c86c57d148c7de473383461B";
//     set_root_args.signatures = vector[invalid_signer_sig]; // add signature from invalid signer
//     call_set_root(&mut env, set_root_args);

//     ts::return_to_sender(&env.scenario, owner_cap);
//     env.destroy();
// }

// #[test]
// #[expected_failure(abort_code = mcms::EInsufficientSigners, location = mcms)]
// fun test_set_root__signer_quorum_not_met() {
//     let mut env = setup();
//     let owner_cap = ts::take_from_sender<OwnerCap>(&env.scenario);
//     let role = mcms::proposer_role();
//     mcms::set_config(
//         &owner_cap,
//         &mut env.state,
//         role,
//         CHAIN_ID,
//         vector[PROPOSER_ADDR1, PROPOSER_ADDR2, PROPOSER_ADDR3],
//         SIGNER_GROUPS,
//         GROUP_QUORUMS,
//         GROUP_PARENTS,
//         false,
//         env.scenario.ctx(),
//     );
//     let mut set_root_args = default_set_root_args(false);
//     let signer1 = set_root_args.signatures[0];
//     set_root_args.signatures = vector[signer1]; // only 1 signature, quorum is 2
//     call_set_root(&mut env, set_root_args);

//     ts::return_to_sender(&env.scenario, owner_cap);
//     env.destroy();
// }

// #[test]
// fun test_set_root__success() {
//     let mut env = setup();
//     let owner_cap = ts::take_from_sender<OwnerCap>(&env.scenario);
//     let expected_role = mcms::proposer_role();
//     mcms::set_config(
//         &owner_cap,
//         &mut env.state,
//         expected_role,
//         CHAIN_ID,
//         vector[PROPOSER_ADDR1, PROPOSER_ADDR2, PROPOSER_ADDR3],
//         SIGNER_GROUPS,
//         GROUP_QUORUMS,
//         GROUP_PARENTS,
//         false,
//         env.scenario.ctx(),
//     );
//     let set_root_args = default_set_root_args(false);
//     call_set_root(&mut env, set_root_args);

//     let (root, valid_until, op_count) = mcms::expiring_root_and_op_count(&env.state, expected_role);
//     assert!(root == ROOT);
//     assert!(valid_until == VALID_UNTIL);
//     assert!(op_count == PRE_OP_COUNT);

//     let root_metadata = &mcms::get_root_metadata(&env.state, expected_role);
//     assert!(mcms::role(root_metadata) == expected_role);
//     assert!(mcms::chain_id(root_metadata) == CHAIN_ID);
//     assert!(mcms::root_metadata_multisig(root_metadata) == @mcms.to_bytes());
//     assert!(mcms::pre_op_count(root_metadata) == PRE_OP_COUNT);
//     assert!(mcms::post_op_count(root_metadata) == POST_OP_COUNT);
//     assert!(mcms::override_previous_root(root_metadata) == false);

//     ts::return_to_sender(&env.scenario, owner_cap);
//     env.destroy();
// }

// ============== set_config tests ================= //

#[test]
#[expected_failure(abort_code = mcms::EInvalidNumSigners, location = mcms)]
fun test_set_config__invalid_number_of_signers() {
    let mut env = setup();
    let owner_cap = ts::take_from_sender<OwnerCap>(&env.scenario);
    // empty signer addresses and groups
    mcms::set_config(
        &owner_cap,
        &mut env.state,
        mcms::proposer_role(),
        CHAIN_ID,
        vector[], // signer_addresses
        vector[], // signer_groups
        GROUP_QUORUMS,
        GROUP_PARENTS,
        false,
        env.scenario.ctx(),
    );
    ts::return_to_sender(&env.scenario, owner_cap);
    env.destroy();
}

#[test]
#[expected_failure(abort_code = mcms::ESignerAddrMustBeIncreasing, location = mcms)]
fun test_set_config__signers_must_be_distinct() {
    let mut env = setup();
    let owner_cap = ts::take_from_sender<OwnerCap>(&env.scenario);
    // signer addresses out of order
    let signer_addresses = vector[PROPOSER_ADDR1, PROPOSER_ADDR3, PROPOSER_ADDR2];
    mcms::set_config(
        &owner_cap,
        &mut env.state,
        mcms::proposer_role(),
        CHAIN_ID,
        signer_addresses, // signer_addresses
        SIGNER_GROUPS,
        GROUP_QUORUMS,
        GROUP_PARENTS,
        false,
        env.scenario.ctx(),
    );
    ts::return_to_sender(&env.scenario, owner_cap);
    env.destroy();
}

#[test]
#[expected_failure(abort_code = mcms::ESignerAddrMustBeIncreasing, location = mcms)]
fun test_set_config__signers_must_be_increasing() {
    let mut env = setup();
    let owner_cap = ts::take_from_sender<OwnerCap>(&env.scenario);
    let signer_addresses = vector[PROPOSER_ADDR1, PROPOSER_ADDR2, PROPOSER_ADDR2];
    mcms::set_config(
        &owner_cap,
        &mut env.state,
        mcms::proposer_role(),
        CHAIN_ID,
        signer_addresses, // signer_addresses
        SIGNER_GROUPS,
        GROUP_QUORUMS,
        GROUP_PARENTS,
        false,
        env.scenario.ctx(),
    );
    ts::return_to_sender(&env.scenario, owner_cap);
    env.destroy();
}

#[test]
#[expected_failure(abort_code = mcms::EInvalidSignerAddrLen, location = mcms)]
fun test_set_config__invalid_signer_address() {
    let mut env = setup();
    let owner_cap = ts::take_from_sender<OwnerCap>(&env.scenario);
    let invalid_signer_addr = x"E37ca797F7fCCFbd9bb3bf8f812F19C3184df1";
    let signer_addresses = vector[PROPOSER_ADDR1, PROPOSER_ADDR2, invalid_signer_addr];

    mcms::set_config(
        &owner_cap,
        &mut env.state,
        mcms::proposer_role(),
        CHAIN_ID,
        signer_addresses,
        SIGNER_GROUPS,
        GROUP_QUORUMS,
        GROUP_PARENTS,
        false,
        env.scenario.ctx(),
    );
    ts::return_to_sender(&env.scenario, owner_cap);
    env.destroy();
}

#[test]
#[expected_failure(abort_code = mcms::EOutOfBoundsGroup, location = mcms)]
fun test_set_config__out_of_bounds_signer_group() {
    let mut env = setup();
    let owner_cap = ts::take_from_sender<OwnerCap>(&env.scenario);

    let signer_addresses = vector[PROPOSER_ADDR1, PROPOSER_ADDR2, PROPOSER_ADDR3];
    // signer group out of bounds
    let signer_groups: vector<u8> = vector[1, 2, mcms::num_groups() as u8];

    mcms::set_config(
        &owner_cap,
        &mut env.state,
        mcms::proposer_role(),
        CHAIN_ID,
        signer_addresses,
        signer_groups,
        GROUP_QUORUMS,
        GROUP_PARENTS,
        false,
        env.scenario.ctx(),
    );
    ts::return_to_sender(&env.scenario, owner_cap);
    env.destroy();
}

#[test]
#[expected_failure(abort_code = mcms::EOutOfBoundsGroupQuorum, location = mcms)]
fun test_set_config__out_of_bounds_group_quorum() {
    let mut env = setup();
    let owner_cap = ts::take_from_sender<OwnerCap>(&env.scenario);

    let signer_addresses = vector[PROPOSER_ADDR1, PROPOSER_ADDR2, PROPOSER_ADDR3];
    // group quorum out of bounds (greater than num signers)
    let mut group_quorums = vector[2, 1, 1, (mcms::max_num_signers() as u8) + 1];
    params::right_pad_vec(&mut group_quorums, mcms::num_groups());

    mcms::set_config(
        &owner_cap,
        &mut env.state,
        mcms::proposer_role(),
        CHAIN_ID,
        signer_addresses,
        SIGNER_GROUPS,
        group_quorums,
        GROUP_PARENTS,
        false,
        env.scenario.ctx(),
    );
    ts::return_to_sender(&env.scenario, owner_cap);
    env.destroy();
}

#[test]
#[expected_failure(abort_code = mcms::EGroupTreeNotWellFormed, location = mcms)]
fun test_set_config__root_is_not_its_own_parent() {
    let mut env = setup();
    let owner_cap = ts::take_from_sender<OwnerCap>(&env.scenario);

    // group parent of root is group 1 (should be itself = group 0)
    let mut group_parents = vector[1];
    params::right_pad_vec(&mut group_parents, mcms::num_groups());

    mcms::set_config(
        &owner_cap,
        &mut env.state,
        mcms::proposer_role(),
        CHAIN_ID,
        vector[PROPOSER_ADDR1, PROPOSER_ADDR2, PROPOSER_ADDR3],
        SIGNER_GROUPS,
        GROUP_QUORUMS,
        group_parents,
        false,
        env.scenario.ctx(),
    );
    ts::return_to_sender(&env.scenario, owner_cap);
    env.destroy();
}

#[test]
#[expected_failure(abort_code = mcms::EGroupTreeNotWellFormed, location = mcms)]
fun test_set_config__non_root_is_its_own_parent() {
    let mut env = setup();
    let owner_cap = ts::take_from_sender<OwnerCap>(&env.scenario);

    // group parent of group 1 is itself (should be lower index group)
    let mut group_parents = vector[0, 1];
    params::right_pad_vec(&mut group_parents, mcms::num_groups());

    mcms::set_config(
        &owner_cap,
        &mut env.state,
        mcms::proposer_role(),
        CHAIN_ID,
        vector[PROPOSER_ADDR1, PROPOSER_ADDR2, PROPOSER_ADDR3],
        SIGNER_GROUPS,
        GROUP_QUORUMS,
        group_parents,
        false,
        env.scenario.ctx(),
    );
    ts::return_to_sender(&env.scenario, owner_cap);
    env.destroy();
}

#[test]
#[expected_failure(abort_code = mcms::EGroupTreeNotWellFormed, location = mcms)]
fun test_set_config__group_parent_higher_index() {
    let mut env = setup();
    let owner_cap = ts::take_from_sender<OwnerCap>(&env.scenario);

    // group parent of group 1 is group 2 (should be lower index group)
    let mut group_parents = vector[0, 2];
    params::right_pad_vec(&mut group_parents, mcms::num_groups());

    mcms::set_config(
        &owner_cap,
        &mut env.state,
        mcms::proposer_role(),
        CHAIN_ID,
        vector[PROPOSER_ADDR1, PROPOSER_ADDR2, PROPOSER_ADDR3],
        SIGNER_GROUPS,
        GROUP_QUORUMS,
        group_parents,
        false,
        env.scenario.ctx(),
    );
    ts::return_to_sender(&env.scenario, owner_cap);
    env.destroy();
}

#[test]
#[expected_failure(abort_code = mcms::EOutOfBoundsGroupQuorum, location = mcms)]
fun test_set_config__quorum_cannot_be_met() {
    let mut env = setup();
    let owner_cap = ts::take_from_sender<OwnerCap>(&env.scenario);

    // group quorum of group 0 (root) is 4, which can never be met because there are only three child groups
    let mut group_quorum = vector[4, 1, 1, 1];
    params::right_pad_vec(&mut group_quorum, mcms::num_groups());

    mcms::set_config(
        &owner_cap,
        &mut env.state,
        mcms::proposer_role(),
        CHAIN_ID,
        vector[PROPOSER_ADDR1, PROPOSER_ADDR2, PROPOSER_ADDR3],
        SIGNER_GROUPS,
        group_quorum,
        GROUP_PARENTS,
        false,
        env.scenario.ctx(),
    );
    ts::return_to_sender(&env.scenario, owner_cap);
    env.destroy();
}

#[test]
#[expected_failure(abort_code = mcms::ESignerInDisabledGroup, location = mcms)]
fun test_set_config__signer_in_disabled_group() {
    let mut env = setup();
    let owner_cap = ts::take_from_sender<OwnerCap>(&env.scenario);

    // group 31 is disabled (quorum = 0) but signer 3 is in group 31
    let signer_groups = vector[1, 2, 31];

    mcms::set_config(
        &owner_cap,
        &mut env.state,
        mcms::proposer_role(),
        CHAIN_ID,
        vector[PROPOSER_ADDR1, PROPOSER_ADDR2, PROPOSER_ADDR3],
        signer_groups,
        GROUP_QUORUMS,
        GROUP_PARENTS,
        false,
        env.scenario.ctx(),
    );
    ts::return_to_sender(&env.scenario, owner_cap);
    env.destroy();
}

#[test]
#[expected_failure(abort_code = mcms::ESignerGroupsLenMismatch, location = mcms)]
fun test_set_config__signer_group_len_mismatch() {
    let mut env = setup();
    let owner_cap = ts::take_from_sender<OwnerCap>(&env.scenario);

    // len of signer groups does not match len of signers
    let signer_groups = vector[1, 2, 3, 3];

    mcms::set_config(
        &owner_cap,
        &mut env.state,
        mcms::proposer_role(),
        CHAIN_ID,
        vector[PROPOSER_ADDR1, PROPOSER_ADDR2, PROPOSER_ADDR3],
        signer_groups,
        GROUP_QUORUMS,
        GROUP_PARENTS,
        false,
        env.scenario.ctx(),
    );
    ts::return_to_sender(&env.scenario, owner_cap);
    env.destroy();
}

#[test]
fun test_set_config__success() {
    let mut env = setup();
    let owner_cap = ts::take_from_sender<OwnerCap>(&env.scenario);

    // manually modify root state to check for modifications
    let role = mcms::proposer_role();
    let new_op_count = 5;
    mcms::test_set_expiring_root_and_op_count(
        &mut env.state, role, ROOT, VALID_UNTIL, new_op_count
    );
    mcms::test_set_root_metadata(
        &mut env.state,
        role,
        CHAIN_ID,
        @mcms.to_bytes(),
        new_op_count,
        new_op_count,
        false
    );

    mcms::set_config(
        &owner_cap,
        &mut env.state,
        mcms::proposer_role(),
        CHAIN_ID,
        vector[PROPOSER_ADDR1, PROPOSER_ADDR2, PROPOSER_ADDR3],
        SIGNER_GROUPS,
        GROUP_QUORUMS,
        GROUP_PARENTS,
        false,
        env.scenario.ctx(),
    );

    // let config = mcms::get_config(role);
    //     assert!(vector::length(&mcms::config_signers(&config)) == 3, 1);
    //     let (addr1, index1, group1) =
    //         mcms::signer_view(&mcms::config_signers(&config)[0]);
    //     assert!(addr1 == PROPOSER_ADDR1, 2);
    //     assert!(index1 == 0, 3);
    //     assert!(group1 == 0, 4);
    //     let (addr2, index2, group2) =
    //         mcms::signer_view(&mcms::config_signers(&config)[1]);
    //     assert!(addr2 == PROPOSER_ADDR2, 5);
    //     assert!(index2 == 1, 6);
    //     assert!(group2 == 0, 7);
    //     let (addr3, index3, group3) =
    //         mcms::signer_view(&mcms::config_signers(&config)[2]);
    //     assert!(addr3 == PROPOSER_ADDR3, 8);
    //     assert!(index3 == 2, 9);
    //     assert!(group3 == 0, 10);
    //     assert!(mcms::config_group_quorums(&config) == GROUP_QUORUMS, 11);
    //     assert!(mcms::config_group_parents(&config) == GROUP_PARENTS, 12);

    //     let (root, valid_until, op_count) = mcms::expiring_root_and_op_count(multisig);
    //     assert!(root == ROOT, 7);
    //     assert!(valid_until == VALID_UNTIL, 8);
    //     assert!(op_count == new_op_count, 9);

    //     let root_metadata = mcms::root_metadata(multisig);
    //     assert!(mcms::role(root_metadata) == role, 10);
    //     assert!(mcms::chain_id(root_metadata) == 1, 11);
    //     assert!(mcms::root_metadata_multisig(root_metadata) == @0xabc, 12);
    //     assert!(mcms::pre_op_count(root_metadata) == 5, 13);
    //     assert!(mcms::post_op_count(root_metadata) == 5, 14);
    //     assert!(!mcms::override_previous_root(root_metadata), 15);

    //     // test set config with clear_root=true, change to 1-of-2 multisig with a nested 2-of-2 multisig
    //     let signer_addr = vector[PROPOSER_ADDR1, PROPOSER_ADDR2, PROPOSER_ADDR3];
    //     let signer_groups = vector[1, 3, 4];
    //     let group_quorums = vector[1, 1, 2, 1, 1];
    //     params::right_pad_vec(&mut group_quorums, mcms::num_groups());
    //     let group_parents = vector[0, 0, 0, 2, 2];
    //     params::right_pad_vec(&mut group_parents, mcms::num_groups());
    //     mcms::set_config(
    //         owner,
    //         mcms::proposer_role(),
    //         signer_addr,
    //         signer_groups,
    //         group_quorums,
    //         group_parents,
    //         true
    //     );
    //     let config = mcms::get_config(role);
    //     assert!(vector::length(&mcms::config_signers(&config)) == 3, 14);
    //     let (addr1, index1, group1) =
    //         mcms::signer_view(&mcms::config_signers(&config)[0]);
    //     let (addr2, index2, group2) =
    //         mcms::signer_view(&mcms::config_signers(&config)[1]);
    //     let (addr3, index3, group3) =
    //         mcms::signer_view(&mcms::config_signers(&config)[2]);
    //     assert!(addr1 == PROPOSER_ADDR1, 15);
    //     assert!(index1 == 0, 16);
    //     assert!(group1 == 1, 17);
    //     assert!(addr2 == PROPOSER_ADDR2, 18);
    //     assert!(index2 == 1, 19);
    //     assert!(group2 == 3, 20);
    //     assert!(addr3 == PROPOSER_ADDR3, 21);
    //     assert!(index3 == 2, 22);
    //     assert!(group3 == 4, 23);
    //     assert!(group_quorums == group_quorums, 24);
    //     assert!(group_parents == group_parents, 25);

    //     let (root, valid_until, op_count) = mcms::expiring_root_and_op_count(multisig);
    //     assert!(root == vector[], 20);
    //     assert!(valid_until == 0, 21);
    //     assert!(op_count == new_op_count, 22);

    //     let root_metadata = &mcms::get_root_metadata(&env.state, expected_role);
    //     assert!(mcms::role(root_metadata) == role, 23);
    //     assert!(mcms::chain_id(root_metadata) == CHAIN_ID, 24);
    //     assert!(mcms::root_metadata_multisig(root_metadata) == @mcms, 25);
    //     assert!(mcms::pre_op_count(root_metadata) == new_op_count, 26);
    //     assert!(mcms::post_op_count(root_metadata) == new_op_count, 27);
    //     assert!(mcms::override_previous_root(root_metadata), 28);

    ts::return_to_sender(&env.scenario, owner_cap);
    env.destroy();
}

fun destroy(env: Env) {
    let Env {
        scenario,
        state,
        timelock,
        registry,
        account_state,
        deployer_state,
        clock,
    } = env;

    ts::return_shared(registry);
    ts::return_shared(timelock);
    ts::return_shared(state);
    ts::return_shared(account_state);
    ts::return_shared(deployer_state);
    clock.destroy_for_testing();

    scenario.end();
}
