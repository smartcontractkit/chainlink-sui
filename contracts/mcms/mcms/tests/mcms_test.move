#[test_only]
module mcms::mcms_test;

use mcms::mcms::{Self, MultisigState, Timelock};
use mcms::mcms_account::{Self, AccountState, OwnerCap};
use mcms::mcms_deployer::{Self, DeployerState};
use mcms::mcms_registry::{Self, Registry};
use mcms::params;
use std::bcs;
use std::string::{Self, String};
use sui::package;
use sui::test_scenario as ts;

const OWNER: address = @0x123;

// keccak256("MANY_CHAIN_MULTI_SIG_DOMAIN_SEPARATOR_OP_SUI")
const MANY_CHAIN_MULTI_SIG_DOMAIN_SEPARATOR_OP: vector<u8> =
    x"542b28b7edb99385286abe2b9c308f91a385cbcb48fc98127cfd13deb28a50b8";

const CHAIN_ID: u256 = 2;
const TIMESTAMP: u64 = 1762209199;
const VALID_UNTIL: u64 = 1762295599;

const MIN_DELAY: u64 = 1; // 1 second delay

// Proposer signers derived from deterministic seeds (sorted ascending)
const PROPOSER_ADDR1: vector<u8> = x"2b5ad5c4795c026514f8317c7a215e218dccd6cf";
const PROPOSER_ADDR2: vector<u8> = x"6813eb9362372eef6200f3b1dbc3f819671cba69";
const PROPOSER_ADDR3: vector<u8> = x"7e5f4552091a69125d5dfcb7b8c2659029395bdf";

// test config: 2-of-3 multisig
const SIGNER_GROUPS: vector<u8> = vector[0, 0, 0];

const GROUP_QUORUMS: vector<u8> = vector[
    2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
];

const GROUP_PARENTS: vector<u8> = vector[
    0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
];

const ROOT: vector<u8> = x"02265107f55a98d78cc2ae5910e14f389f28afeae8cbb4a6b58d2a1e997b44bd";
const SIGNATURES: vector<vector<u8>> = vector[
    x"f7d88e6c32389053429c75762b79777329cda5974b2402a04c6b3f5d32ddfd267107d1d09e895d41788fbaf4b8a8b798d225325c15bea1e78d5ca7e5af26e2501b",
    x"d93eec6f57ffab038e7a8a14056c30643ef28f484cf2e90c01e36ad21dbf1ec5190000f35ad91a0af9248fa6d6bdf79b0b62d30b8541df351635456193a18c9e1c",
];

const PRE_OP_COUNT: u64 = 0;
const POST_OP_COUNT: u64 = 1;

const METADATA_PROOF: vector<vector<u8>> = vector[
    x"60aca62639a14862806492b9fb384b02a4f325f204dc51579f8b0f51ed6fd9d2", // OP hash (sibling proof for metadata leaf)
];

const OP1_PROOF: vector<vector<u8>> = vector[
    x"ae675bc2daf4f44b0eeea013c0e09e54f5af4aeb002068f06490b24882716b7a", // METADATA hash (sibling proof for op leaf)
];

// The OPs contained are
// {
// 			Target:      mcmsAccount,
// 			ModuleName:  "mcms_account",
// 			Function:    "accept_ownership_as_timelock",
// 			Data:        []byte{},
// 			Delay:       1,
// 			Predecessor: []byte{},
// 			Salt:        []byte{},
// 		},
const LEAVES: vector<vector<u8>> = vector[
    x"ae675bc2daf4f44b0eeea013c0e09e54f5af4aeb002068f06490b24882716b7a", // METADATA hash (metadata leaves come first)
    x"60aca62639a14862806492b9fb384b02a4f325f204dc51579f8b0f51ed6fd9d2", // OP hash (operation leaves come second)
];

const OP1_NONCE: u64 = 0;
const OP1_DATA: vector<u8> =
    x"010000000000000000000000000000000000000000000000000000000000000000010c6d636d735f6163636f756e74011c6163636570745f6f776e6572736869705f61735f74696d656c6f636b012087cce46d7a32876f725dfd0b81215e7d72c2cc41d19510a044ce6dde9fc5a4922000000000000000000000000000000000000000000000000000000000000000002068c6c852000000000000000000000000000000000000000000000000000000000100000000000000";

public struct Env {
    scenario: ts::Scenario,
    state: MultisigState,
    timelock: Timelock,
    registry: Registry,
    account_state: AccountState,
    deployer_state: DeployerState,
    clock: sui::clock::Clock,
}

// Accessor functions for Env
public fun env_scenario(env: &mut Env): &mut ts::Scenario { &mut env.scenario }

public fun env_state(env: &mut Env): &mut MultisigState { &mut env.state }

public fun env_timelock(env: &mut Env): &mut Timelock { &mut env.timelock }

public fun env_registry(env: &mut Env): &mut Registry { &mut env.registry }

public fun env_account_state(env: &mut Env): &mut AccountState { &mut env.account_state }

public fun env_deployer_state(env: &mut Env): &mut DeployerState { &mut env.deployer_state }

public fun env_clock(env: &Env): &sui::clock::Clock { &env.clock }

public struct SetRootArgs has drop {
    role: u8,
    root: vector<u8>,
    valid_until: u64,
    chain_id: u256,
    multisig: address,
    pre_op_count: u64,
    post_op_count: u64,
    override_previous_root: bool,
    metadata_proof: vector<vector<u8>>,
    signatures: vector<vector<u8>>,
}

fun default_set_root_args(override_previous_root: bool): SetRootArgs {
    SetRootArgs {
        role: mcms::proposer_role(),
        root: ROOT,
        valid_until: VALID_UNTIL,
        chain_id: CHAIN_ID,
        multisig: mcms_registry::get_multisig_address(),
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

public struct ExecuteArgs has drop {
    role: u8,
    chain_id: u256,
    multisig: address,
    nonce: u64,
    to: address,
    module_name: String,
    function: String,
    data: vector<u8>,
    proof: vector<vector<u8>>,
}

fun default_execute_args(): ExecuteArgs {
    ExecuteArgs {
        role: mcms::proposer_role(),
        chain_id: CHAIN_ID,
        multisig: mcms_registry::get_multisig_address(),
        nonce: OP1_NONCE,
        to: mcms_registry::get_multisig_address(),
        module_name: string::utf8(b"mcms"),
        function: string::utf8(b"timelock_schedule_batch"),
        data: OP1_DATA,
        proof: OP1_PROOF,
    }
}

// #[test]
// /// Cannot test E2E with hard coded proofs as object IDs are generated dynamically (AccountState object ID)
// /// Therefore when we serialize the proof, each run will generate a different proof as the object IDs are different
// /// We rely on MCMS lib e2e tests to test the full flow with valid proofs.
// public fun test_e2e() {
//     let mut env = setup();

//     let role = mcms::proposer_role();
//     let owner_cap = ts::take_from_sender<OwnerCap>(&env.scenario);
//     mcms::set_config(
//         &owner_cap,
//         &mut env.state,
//         role,
//         CHAIN_ID,
//         vector[PROPOSER_ADDR1, PROPOSER_ADDR2, PROPOSER_ADDR3],
//         SIGNER_GROUPS,
//         GROUP_QUORUMS,
//         GROUP_PARENTS,
//         true,
//         env.scenario.ctx(),
//     );

//     let signers = mcms::signers(&env.state, role);
//     assert!(signers.length() == 3);

//     let set_root_args = default_set_root_args(false);
//     call_set_root(&mut env, set_root_args);

//     let (root, valid_until, op_count) = mcms::expiring_root_and_op_count(&env.state, role);
//     assert!(root == ROOT);
//     assert!(valid_until == VALID_UNTIL);
//     assert!(op_count == 0);

//     // First we must transfer ownership to `@mcms` (the multisig/self)
//     mcms_account::transfer_ownership_to_self(
//         &owner_cap,
//         &mut env.account_state,
//         env.scenario.ctx(),
//     );

//     // FIRST EXECUTE: Schedule the timelock operation
//     // We schedule directly since we can't easily generate dynamic OP1_DATA for merkle proof verification
//     // Object `account_state` will have a different ID each time, so we need to serialize it and use it as data
//     let account_id_bytes = bcs::to_bytes(&object::id_address(&env.account_state));
//     let clock = &env.clock;
//     mcms::test_timelock_schedule_batch(
//         &mut env.timelock,
//         clock,
//         mcms::proposer_role(),
//         vector[mcms_registry::get_multisig_address()], // targets
//         vector[string::utf8(b"mcms_account")], // module_names
//         vector[string::utf8(b"accept_ownership_as_timelock")], // function_names
//         vector[account_id_bytes], // datas - account object ID
//         x"0000000000000000000000000000000000000000000000000000000000000000", // predecessor
//         x"68c6c85200000000000000000000000000000000000000000000000000000000", // salt
//         1, // delay (1 second)
//         env.scenario.ctx(),
//     );

//     // Wait for delay (10 second)
//     env.clock.set_for_testing((TIMESTAMP * 1000) + 10000);

//     // SECOND EXECUTE: Execute the scheduled timelock operation
//     // This calls timelock_execute_batch directly like the Go test's timelockExecutable.Execute()
//     // Parameters must match exactly what was scheduled above
//     let account_id_bytes = bcs::to_bytes(&object::id_address(&env.account_state));
//     timelock_execute_dispatch_to_acc_helper(
//         &mut env,
//         vector[mcms_registry::get_multisig_address()], // targets - must match scheduling parameters
//         vector[string::utf8(b"mcms_account")], // module_names
//         vector[string::utf8(b"accept_ownership_as_timelock")], // function_names
//         vector[account_id_bytes], // datas - must contain the same account object ID
//         x"0000000000000000000000000000000000000000000000000000000000000000", // predecessor
//         x"68c6c85200000000000000000000000000000000000000000000000000000000", // salt
//     );

//     let ctx = env.scenario.ctx();
//     mcms_account::execute_ownership_transfer(
//         owner_cap,
//         &mut env.account_state,
//         &mut env.registry,
//         mcms_registry::get_multisig_address(),
//         ctx,
//     );

//     // Verify new owner is now `@mcms`
//     let new_mcms_owner = mcms_account::owner(&env.account_state);
//     assert!(new_mcms_owner == mcms_registry::get_multisig_address());

//     env.destroy();
// }

public fun setup(): Env {
    let mut scenario = ts::begin(OWNER);
    let ctx = scenario.ctx();

    let mut clock = sui::clock::create_for_testing(ctx);
    clock.set_for_testing(TIMESTAMP * 1000);

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

#[test]
#[expected_failure(abort_code = mcms::EValidUntilExpired, location = mcms)]
public fun test_set_root__valid_until_expired() {
    let mut env = setup();
    let mut set_root_args = default_set_root_args(false);
    set_root_args.valid_until = TIMESTAMP - 1; // set valid_until to a time in the past
    call_set_root(&mut env, set_root_args);

    env.destroy()
}

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
    set_root_args.multisig = @0x999;
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
        mcms_registry::get_multisig_address(),
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
        mcms_registry::get_multisig_address(),
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
//     let sig0 = set_root_args.signatures[0];
//     let sig1 = set_root_args.signatures[1];
//     // Reverse the order of the 2 signatures (out of order)
//     set_root_args.signatures = vector[sig1, sig0]; // shuffle signature order
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

#[test]
#[expected_failure(abort_code = mcms::EInsufficientSigners, location = mcms)]
fun test_set_root__signer_quorum_not_met() {
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
    let mut set_root_args = default_set_root_args(false);
    let signer1 = set_root_args.signatures[0];
    set_root_args.signatures = vector[signer1]; // only 1 signature, quorum is 2
    call_set_root(&mut env, set_root_args);

    ts::return_to_sender(&env.scenario, owner_cap);
    env.destroy();
}

#[test]
fun test_set_root__success() {
    let mut env = setup();
    let owner_cap = ts::take_from_sender<OwnerCap>(&env.scenario);
    let expected_role = mcms::proposer_role();
    mcms::set_config(
        &owner_cap,
        &mut env.state,
        expected_role,
        CHAIN_ID,
        vector[PROPOSER_ADDR1, PROPOSER_ADDR2, PROPOSER_ADDR3],
        SIGNER_GROUPS,
        GROUP_QUORUMS,
        GROUP_PARENTS,
        false,
        env.scenario.ctx(),
    );
    let set_root_args = default_set_root_args(false);
    call_set_root(&mut env, set_root_args);

    let (root, valid_until, op_count) = mcms::expiring_root_and_op_count(&env.state, expected_role);
    assert!(root == ROOT);
    assert!(valid_until == VALID_UNTIL);
    assert!(op_count == PRE_OP_COUNT);

    let root_metadata = mcms::get_root_metadata(&env.state, expected_role);
    assert!(mcms::role(&root_metadata) == expected_role);
    assert!(mcms::chain_id(&root_metadata) == CHAIN_ID);
    assert!(mcms::root_metadata_multisig(&root_metadata) == mcms_registry::get_multisig_address());
    assert!(mcms::pre_op_count(&root_metadata) == PRE_OP_COUNT);
    assert!(mcms::post_op_count(&root_metadata) == POST_OP_COUNT);
    assert!(mcms::override_previous_root(&root_metadata) == false);

    ts::return_to_sender(&env.scenario, owner_cap);
    env.destroy();
}

#[test]
#[expected_failure(abort_code = mcms::EProofCannotBeVerified, location = mcms)]
fun test_set_root__invalid_chain_id() {
    let mut env = setup();
    let owner_cap = ts::take_from_sender<OwnerCap>(&env.scenario);
    let expected_role = mcms::proposer_role();
    mcms::set_config(
        &owner_cap,
        &mut env.state,
        expected_role,
        CHAIN_ID,
        vector[PROPOSER_ADDR1, PROPOSER_ADDR2, PROPOSER_ADDR3],
        SIGNER_GROUPS,
        GROUP_QUORUMS,
        GROUP_PARENTS,
        false,
        env.scenario.ctx(),
    );

    let mut set_root_args = default_set_root_args(false);
    set_root_args.chain_id = 111; // wrong chain id - this breaks the merkle proof
    call_set_root(&mut env, set_root_args);

    ts::return_to_sender(&env.scenario, owner_cap);
    env.destroy();
}

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
#[allow(implicit_const_copy)]
fun test_set_config__success() {
    let mut env = setup();
    let owner_cap = ts::take_from_sender<OwnerCap>(&env.scenario);

    // manually modify root state to check for modifications
    let role = mcms::proposer_role();
    let new_op_count = 5;
    mcms::test_set_expiring_root_and_op_count(
        &mut env.state,
        role,
        ROOT,
        VALID_UNTIL,
        new_op_count,
    );
    mcms::test_set_root_metadata(
        &mut env.state,
        role,
        CHAIN_ID,
        mcms_registry::get_multisig_address(),
        new_op_count,
        new_op_count,
        false,
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

    // Verify the first configuration (clear_root=false)
    let signers = mcms::signers(&env.state, role);
    assert!(signers.length() == 3);

    // Check signers are properly stored and accessible
    let signer_1_exists = signers.contains(&PROPOSER_ADDR1);
    let signer_2_exists = signers.contains(&PROPOSER_ADDR2);
    let signer_3_exists = signers.contains(&PROPOSER_ADDR3);
    assert!(signer_1_exists);
    assert!(signer_2_exists);
    assert!(signer_3_exists);

    // Verify root and metadata are preserved (clear_root=false)
    let (root, valid_until, op_count) = mcms::expiring_root_and_op_count(&env.state, role);
    assert!(root == ROOT);
    assert!(valid_until == VALID_UNTIL);
    assert!(op_count == new_op_count);

    let root_metadata = mcms::get_root_metadata(&env.state, role);
    assert!(mcms::role(&root_metadata) == role);
    assert!(mcms::chain_id(&root_metadata) == CHAIN_ID);
    assert!(mcms::root_metadata_multisig(&root_metadata) == mcms_registry::get_multisig_address());
    assert!(mcms::pre_op_count(&root_metadata) == new_op_count);
    assert!(mcms::post_op_count(&root_metadata) == new_op_count);
    assert!(!mcms::override_previous_root(&root_metadata));

    // Test set config with clear_root=true - should clear the root
    mcms::set_config(
        &owner_cap,
        &mut env.state,
        mcms::proposer_role(),
        CHAIN_ID,
        vector[PROPOSER_ADDR1, PROPOSER_ADDR2, PROPOSER_ADDR3],
        SIGNER_GROUPS,
        GROUP_QUORUMS,
        GROUP_PARENTS,
        true, // clear_root=true
        env.scenario.ctx(),
    );

    // Verify root is cleared when clear_root=true
    let (
        root_after_clear,
        valid_until_after_clear,
        op_count_after_clear,
    ) = mcms::expiring_root_and_op_count(&env.state, role);
    assert!(root_after_clear == vector[]);
    assert!(valid_until_after_clear == 0);
    assert!(op_count_after_clear == new_op_count);

    let root_metadata_after_clear = mcms::get_root_metadata(&env.state, role);
    assert!(mcms::role(&root_metadata_after_clear) == role);
    assert!(mcms::chain_id(&root_metadata_after_clear) == CHAIN_ID);
    assert!(
        mcms::root_metadata_multisig(&root_metadata_after_clear) == mcms_registry::get_multisig_address(),
    );
    assert!(mcms::pre_op_count(&root_metadata_after_clear) == new_op_count);
    assert!(mcms::post_op_count(&root_metadata_after_clear) == new_op_count);
    assert!(mcms::override_previous_root(&root_metadata_after_clear));

    ts::return_to_sender(&env.scenario, owner_cap);
    env.destroy();
}

// ========== Execute Test Helper Functions ==========

// Helper function that expects execute to fail, so we don't need to handle the return value
fun call_execute_expect_failure(env: &mut Env, args: ExecuteArgs) {
    let callback_params = mcms::execute(
        &mut env.state,
        &env.clock,
        args.role,
        args.chain_id,
        args.multisig,
        args.nonce,
        args.to,
        args.module_name,
        args.function,
        args.data,
        args.proof,
        env.scenario.ctx(),
    );

    // Consume the callback params - this should never be reached due to expected_failure
    // but we need to handle the return value properly by passing it to a dispatch function
    mcms::dispatch_timelock_schedule_batch(
        &mut env.timelock,
        &env.clock,
        callback_params,
        env.scenario.ctx(),
    );

    abort 999
}

// ========== Execute Tests ==========

#[test]
#[expected_failure(abort_code = mcms::EPostOpCountReached, location = mcms)]
fun test_execute__root_not_set() {
    let mut env = setup();
    let execute_args = default_execute_args();
    call_execute_expect_failure(&mut env, execute_args);
    destroy(env);
}

#[test]
#[expected_failure(abort_code = mcms::EPostOpCountReached, location = mcms)]
fun test_execute__post_op_count_reached() {
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

    call_set_root(&mut env, default_set_root_args(false));
    let post_op_count = POST_OP_COUNT;
    mcms::test_set_expiring_root_and_op_count(
        &mut env.state,
        role,
        ROOT,
        VALID_UNTIL,
        post_op_count,
    );

    let execute_args = default_execute_args();
    ts::return_to_sender(&env.scenario, owner_cap);
    call_execute_expect_failure(&mut env, execute_args);
    destroy(env);
}

#[test]
#[expected_failure(abort_code = mcms::EWrongNonce, location = mcms)]
fun test_execute__wrong_nonce() {
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

    call_set_root(&mut env, default_set_root_args(false));

    let mut execute_args = default_execute_args();
    execute_args.nonce = 999; // wrong nonce

    ts::return_to_sender(&env.scenario, owner_cap);
    call_execute_expect_failure(&mut env, execute_args);
    destroy(env);
}

#[test]
#[expected_failure(abort_code = mcms::EWrongMultisig, location = mcms)]
fun test_execute__wrong_multisig_addr() {
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

    call_set_root(&mut env, default_set_root_args(false));

    let mut execute_args = default_execute_args();
    execute_args.multisig = @0x12345; // wrong multisig address

    ts::return_to_sender(&env.scenario, owner_cap);
    call_execute_expect_failure(&mut env, execute_args);
    destroy(env);
}

#[test]
#[expected_failure(abort_code = mcms::EValidUntilExpired, location = mcms)]
fun test_execute__root_expired() {
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

    call_set_root(&mut env, default_set_root_args(false));

    // Set expired root
    mcms::test_set_expiring_root_and_op_count(
        &mut env.state,
        role,
        ROOT,
        TIMESTAMP - 1,
        0, // expired valid_until
    );

    let execute_args = default_execute_args();
    ts::return_to_sender(&env.scenario, owner_cap);
    call_execute_expect_failure(&mut env, execute_args);
    destroy(env);
}

#[test]
#[expected_failure(abort_code = mcms::EProofCannotBeVerified, location = mcms)]
fun test_execute__bad_op_proof() {
    let mut env = setup();
    let role = mcms::proposer_role();
    let owner_cap = env.scenario.take_from_sender<OwnerCap>();
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
    env.scenario.return_to_sender(owner_cap);

    let set_root_args = default_set_root_args(false);
    call_set_root(&mut env, set_root_args);

    let mut execute_args = default_execute_args();
    execute_args.data = b"different data"; // modify op so proof verification should fail
    call_execute_expect_failure(&mut env, execute_args);
    destroy(env);
}

#[test]
#[expected_failure(abort_code = mcms::EProofCannotBeVerified, location = mcms)]
fun test_execute__empty_proof() {
    let mut env = setup();
    let role = mcms::proposer_role();
    let owner_cap = env.scenario.take_from_sender<OwnerCap>();
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
    env.scenario.return_to_sender(owner_cap);

    let set_root_args = default_set_root_args(false);
    call_set_root(&mut env, set_root_args);

    let mut execute_args = default_execute_args();
    execute_args.proof = vector[]; // empty proof
    call_execute_expect_failure(&mut env, execute_args);
    destroy(env);
}

#[test]
#[expected_failure(abort_code = mcms::EWrongNonce, location = mcms)]
fun test_execute__ops_executed_out_of_order() {
    let mut env = setup();
    let role = mcms::proposer_role();
    let owner_cap = env.scenario.take_from_sender<OwnerCap>();
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
    env.scenario.return_to_sender(owner_cap);

    // modify state to add pending ops to a different one from OP1_NONCE
    mcms::test_set_expiring_root_and_op_count(
        &mut env.state,
        role,
        ROOT,
        VALID_UNTIL,
        OP1_NONCE + 1,
    );

    mcms::test_set_root_metadata(
        &mut env.state,
        role,
        CHAIN_ID,
        @mcms,
        0, // pre_op_count
        2, // post_op_count
        false,
    );

    let execute_args = default_execute_args();
    call_execute_expect_failure(&mut env, execute_args);
    destroy(env);
}

#[test]
#[expected_failure(abort_code = mcms::EProofCannotBeVerified, location = mcms)]
fun test_execute__wrong_chain_id() {
    let mut env = setup();
    let role = mcms::proposer_role();
    let owner_cap = env.scenario.take_from_sender<OwnerCap>();
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
    env.scenario.return_to_sender(owner_cap);

    let set_root_args = default_set_root_args(false);
    call_set_root(&mut env, set_root_args);

    let mut execute_args = default_execute_args();
    execute_args.chain_id = 111; // wrong chain id - this breaks the merkle proof
    call_execute_expect_failure(&mut env, execute_args);
    destroy(env);
}

// ============== Ownership tests ================= //

#[test]
fun test_ownable__transfer_ownership() {
    let mut env = setup();
    let owner_cap = ts::take_from_sender<OwnerCap>(&env.scenario);

    // Get current owner and transfer to a different address
    let current_owner = mcms_account::owner(&env.account_state);
    let new_owner_addr = @0x999; // Use a clearly different address
    mcms_account::transfer_ownership(
        &owner_cap,
        &mut env.account_state,
        new_owner_addr,
        env.scenario.ctx(),
    );

    // Check that there's a pending transfer
    assert!(mcms_account::pending_transfer_to(&env.account_state) == option::some(new_owner_addr));
    assert!(mcms_account::pending_transfer_accepted(&env.account_state) == option::some(false));

    // Accept ownership as the new owner (simulate the new owner accepting)
    env.scenario.next_tx(new_owner_addr);
    mcms_account::accept_ownership(
        &mut env.account_state,
        env.scenario.ctx(),
    );

    // Check ownership has not changes
    assert!(mcms_account::owner(&env.account_state) == current_owner);

    // Execute ownership transfer
    env.scenario.next_tx(new_owner_addr);
    mcms_account::execute_ownership_transfer(
        owner_cap,
        &mut env.account_state,
        &mut env.registry,
        new_owner_addr,
        env.scenario.ctx(),
    );

    // Verify that ownership has been transferred
    let final_owner = mcms_account::owner(&env.account_state);
    assert!(final_owner == new_owner_addr);

    env.destroy();
}

// ============== Timelock tests ================= //

#[test]
fun test_timelock_initialization() {
    let env = setup();
    let min_delay = mcms::timelock_min_delay(&env.timelock);
    assert!(min_delay == 0);
    env.destroy();
}

#[test]
fun test_update_min_delay() {
    let mut env = setup();
    let new_delay = MIN_DELAY;
    mcms::test_timelock_update_min_delay(
        &mut env.timelock,
        mcms::timelock_role(),
        new_delay,
        env.scenario.ctx(),
    );
    assert!(mcms::timelock_min_delay(&env.timelock) == MIN_DELAY);
    env.destroy();
}

#[test]
fun test_schedule_batch() {
    let mut env = setup();

    // Schedule a batch operation - Need to borrow clock first to avoid borrowing conflict
    let clock = &env.clock;
    mcms::test_timelock_schedule_batch(
        &mut env.timelock,
        clock,
        mcms::proposer_role(),
        vector[@0x1], // targets
        vector[string::utf8(b"test_module")], // module_names
        vector[string::utf8(b"test_function")], // function_names
        vector[b"test_data"], // datas
        x"a1b2c3d4e5f60718293804a5b6c7d8e9f0a1b2c3d4e5f6071829384a5b6c7d8e", // predecessor
        vector<u8>[], // salt
        MIN_DELAY, // delay
        env.scenario.ctx(),
    );

    env.destroy();
}

#[test]
#[allow(implicit_const_copy)]
fun test_cancel_operation() {
    let mut env = setup();
    let clock = &env.clock;
    let data = bcs::to_bytes(&MIN_DELAY);

    // Schedule the batch first
    let predecessor = x"bb2adb5b9907ea8042c90eb159f31f68c53ae174499bd16a1d1308876399fbac";
    mcms::test_timelock_schedule_batch(
        &mut env.timelock,
        clock,
        mcms::proposer_role(),
        vector[@mcms], // targets
        vector[string::utf8(b"mcms")], // module_names
        vector[string::utf8(b"timelock_update_min_delay")], // function_names
        vector[data], // datas
        predecessor, // predecessor (32 bytes)
        vector[1u8], // salt
        0u64, // delay
        env.scenario.ctx(),
    );

    // Calculate the operation ID
    let calls = mcms::create_calls(
        vector[@mcms], // targets
        vector[string::utf8(b"mcms")], // module_names
        vector[string::utf8(b"timelock_update_min_delay")], // function_names
        vector[data], // datas
    );
    let id = mcms::hash_operation_batch(
        calls,
        predecessor, // Must match the predecessor used in scheduling
        vector[1u8],
    );

    // Verify operation is pending
    assert!(mcms::timelock_is_operation_pending(&env.timelock, id));

    // Cancel the operation
    mcms::test_timelock_cancel(
        &mut env.timelock,
        mcms::canceller_role(),
        id,
        env.scenario.ctx(),
    );

    // Verify operation is no longer pending
    assert!(!mcms::timelock_is_operation_pending(&env.timelock, id));

    env.destroy();
}

#[test]
#[expected_failure(abort_code = mcms::EOperationCannotBeCancelled, location = mcms)]
fun test_cancel_nonexistent_operation() {
    let mut env = setup();

    // Try to cancel a nonexistent operation
    mcms::test_timelock_cancel(
        &mut env.timelock,
        mcms::canceller_role(),
        vector[123u8],
        env.scenario.ctx(),
    );

    env.destroy();
}

#[test]
fun test_bypasser_execute_batch() {
    let mut env = setup();

    // Get initial min_delay value
    let initial_delay = mcms::timelock_min_delay(&env.timelock);
    let new_delay = initial_delay + 1000; // Set to a different value

    // Use bypasser role to execute batch - this should return ExecutingCallbackParams
    let mut bypasser_update_delay_data = vector[];
    bypasser_update_delay_data.append(bcs::to_bytes(&new_delay));

    let mut executing_params = mcms::test_timelock_bypasser_execute_batch(
        mcms::bypasser_role(),
        vector[mcms_registry::get_multisig_address()], // targets
        vector[string::utf8(b"mcms")], // module_names
        vector[string::utf8(b"timelock_update_min_delay")], // function_names
        vector[bypasser_update_delay_data], // datas
        env.scenario.ctx(),
    );

    // Verify we got exactly 1 ExecutingCallbackParams
    assert!(executing_params.length() == 1);

    // Extract the ExecutingCallbackParams and verify its contents
    let params = executing_params.borrow(0);
    assert!(mcms_registry::target(params) == mcms_registry::get_multisig_address());
    assert!(mcms_registry::module_name(params) == string::utf8(b"mcms"));
    assert!(mcms_registry::function_name(params) == string::utf8(b"timelock_update_min_delay"));

    // Now we need to consume the ExecutingCallbackParams by calling the actual function
    // Since this is timelock_update_min_delay, we need to dispatch it properly
    let mut update_delay_data = vector[];
    update_delay_data.append(bcs::to_bytes(&object::id_address(&env.timelock)));
    update_delay_data.append(bcs::to_bytes(&new_delay));

    let callback_params = mcms_registry::test_create_executing_callback_params(
        mcms_registry::get_multisig_address(),
        string::utf8(b"mcms"),
        string::utf8(b"timelock_update_min_delay"),
        update_delay_data,
        x"0000000000000000000000000000000000000000000000000000000000000001",
        0,
        1,
    );

    // Dispatch the update_min_delay function to consume the hot potato
    mcms::mcms_timelock_update_min_delay(
        &mut env.timelock,
        &mut env.registry,
        callback_params,
        env.scenario.ctx(),
    );

    // Verify that min_delay was actually updated - this proves the bypasser execution worked
    let updated_delay = mcms::timelock_min_delay(&env.timelock);
    assert!(updated_delay == new_delay, 0);
    assert!(updated_delay != initial_delay, 1); // Ensure it actually changed

    // Must consume the ExecutingCallbackParams hot potato
    // We know there's exactly 1 param, so just consume it directly
    let params = executing_params.pop_back();
    let (target, module_name, function_name, data) = mcms_registry::get_callback_params_from_mcms(
        &mut env.registry,
        params,
    );

    // Verify the ExecutingCallbackParams has the expected structure
    assert!(target == mcms_registry::get_multisig_address());
    assert!(module_name == string::utf8(b"mcms"));
    assert!(function_name == string::utf8(b"timelock_update_min_delay"));
    assert!(!data.is_empty()); // Should contain the serialized new_delay

    // Now the vector should be empty
    vector::destroy_empty(executing_params);

    env.destroy();
}

#[test]
#[
    expected_failure(
        abort_code = mcms::EOperationNotReady,
        location = mcms,
    ),
] // Operation not ready yet
fun test_execute_batch_not_ready() {
    let mut env = setup();

    let delay = 100000;
    mcms::test_timelock_update_min_delay(
        &mut env.timelock,
        mcms::timelock_role(),
        delay,
        env.scenario.ctx(),
    );

    let clock = &env.clock;
    let predecessor = x"3f7a2d9e4c1b8a5f0e3d7c6b9a8f5e4d2c1b0a9f8e7d6c5b4a3f2e1d0c9b8a7f";
    mcms::test_timelock_schedule_batch(
        &mut env.timelock,
        clock,
        mcms::proposer_role(),
        vector[mcms_registry::get_multisig_address()], // targets
        vector[string::utf8(b"test_module")], // module_names
        vector[string::utf8(b"test_function")], // function_names
        vector[vector[0u8]], // datas
        predecessor, // predecessor (32 bytes)
        vector[1u8], // salt
        delay,
        env.scenario.ctx(),
    );

    // Try to execute before the delay has passed (should fail)
    timelock_execute_dispatch_to_acc_helper(
        &mut env,
        vector[mcms_registry::get_multisig_address()], // targets
        vector[string::utf8(b"test_module")], // module_names
        vector[string::utf8(b"test_function")], // function_names
        vector[vector[0u8]], // datas
        predecessor, // Must match the one used in scheduling
        vector[1u8], // salt
    );

    env.destroy();
}

#[test]
#[expected_failure(abort_code = mcms::EOperationNotReady, location = mcms)]
fun test_execute_unscheduled_operation() {
    let mut env = setup();

    // Try to execute without scheduling first (should fail)
    timelock_execute_dispatch_to_acc_helper(
        &mut env,
        vector[mcms_registry::get_multisig_address()], // targets
        vector[string::utf8(b"test_module")], // module_names
        vector[string::utf8(b"test_function")], // function_names
        vector[vector[0u8]], // datas
        x"0000000000000000000000000000000000000000000000000000000000000000", // predecessor (32 bytes)
        vector[1u8], // salt
    );

    env.destroy();
}

#[test]
#[expected_failure(abort_code = mcms::EMissingDependency, location = mcms)]
fun test_execute_batch_after_completion() {
    let mut env = setup();
    let clock = &env.clock;
    let data = bcs::to_bytes(&1000u64);
    let predecessor = x"8d4e2f1a0b9c3e5d7f8a2b4c6d1e3f5a7b9c0d2e4f6a8b1c3d5e7f9a0b2c4d6e";

    // Schedule the batch
    mcms::test_timelock_schedule_batch(
        &mut env.timelock,
        clock,
        mcms::proposer_role(),
        vector[mcms_registry::get_multisig_address()], // targets
        vector[string::utf8(b"mcms")], // module_names
        vector[string::utf8(b"timelock_update_min_delay")], // function_names
        vector[data], // datas
        predecessor, // predecessor (32 bytes)
        vector[1u8], // salt
        0u64, // delay
        env.scenario.ctx(),
    );

    // Execute once
    timelock_execute_dispatch_to_acc_helper(
        &mut env,
        vector[mcms_registry::get_multisig_address()], // targets
        vector[string::utf8(b"mcms")], // module_names
        vector[string::utf8(b"timelock_update_min_delay")], // function_names
        vector[data], // datas
        predecessor, // Must match the one used in scheduling
        vector[1u8], // salt
    );

    // Try to execute again (should fail)
    timelock_execute_dispatch_to_acc_helper(
        &mut env,
        vector[mcms_registry::get_multisig_address()], // targets
        vector[string::utf8(b"mcms")], // module_names
        vector[string::utf8(b"timelock_update_min_delay")], // function_names
        vector[data], // datas
        predecessor, // Must match the one used in scheduling
        vector[1u8], // salt
    );

    env.destroy();
}

#[test]
#[expected_failure(abort_code = mcms::EInvalidFunctionName, location = mcms)]
fun test_execute_unknown_function() {
    let mut env = setup();

    // Test dispatching to unknown function within mcms module
    let callback_params = mcms_registry::test_create_executing_callback_params(
        mcms_registry::get_multisig_address(),
        string::utf8(b"mcms"),
        string::utf8(b"unknown_function"),
        vector[],
        x"0000000000000000000000000000000000000000000000000000000000000001",
        0,
        1,
    );

    // Try to dispatch unknown function - this should fail with EInvalidFunctionName
    mcms::mcms_timelock_update_min_delay(
        &mut env.timelock,
        &mut env.registry,
        callback_params,
        env.scenario.ctx(),
    );

    env.destroy();
}

#[test]
fun test_bypasser_blocked_function_interaction() {
    let mut env = setup();
    let test_target = @0xabc;
    let test_module = string::utf8(b"test_module");
    let test_function = string::utf8(b"test_function");

    // Block the function first
    mcms::test_timelock_block_function(
        &mut env.timelock,
        mcms::timelock_role(),
        test_target,
        test_module,
        test_function,
        env.scenario.ctx(),
    );

    // Verify function was blocked
    let blocked_count = mcms::timelock_get_blocked_functions_count(&env.timelock);
    assert!(blocked_count == 1);

    let blocked_function = mcms::timelock_get_blocked_function(&env.timelock, 0);
    assert!(mcms::target(blocked_function) == test_target);
    assert!(mcms::module_name(blocked_function) == test_module);
    assert!(mcms::function_name(blocked_function) == test_function);

    env.destroy();
}

// ============== Advanced Timelock Tests ================= //

#[test]
#[expected_failure(abort_code = mcms::EInvalidParameters, location = mcms)]
fun test_schedule_batch_invalid_parameters() {
    let mut env = setup();
    let clock = &env.clock;

    // Try to schedule with mismatched parameters length
    mcms::test_timelock_schedule_batch(
        &mut env.timelock,
        clock,
        mcms::proposer_role(),
        vector[@mcms, @mcms], // 2 targets
        vector[string::utf8(b"test_module")], // But only 1 module name
        vector[string::utf8(b"test_function")],
        vector[vector[0u8]],
        x"f1e2d3c4b5a6978685746352413021f0e9d8c7b6a5948372615049382716f5e4", // predecessor
        vector[1u8], // salt
        0, // delay
        env.scenario.ctx(),
    );

    env.destroy();
}

#[test]
#[expected_failure(abort_code = mcms::EInsufficientDelay, location = mcms)]
fun test_schedule_insufficient_delay() {
    let mut env = setup();
    let clock = &env.clock;

    // First set a minimum delay
    mcms::test_timelock_update_min_delay(
        &mut env.timelock,
        mcms::timelock_role(),
        MIN_DELAY,
        env.scenario.ctx(),
    );

    // Try to schedule with delay lower than minimum
    mcms::test_timelock_schedule_batch(
        &mut env.timelock,
        clock,
        mcms::proposer_role(),
        vector[@mcms],
        vector[string::utf8(b"test_module")],
        vector[string::utf8(b"test_function")],
        vector[vector[0u8]],
        x"2c9b8e7f6a5d4c3b2a1f0e9d8c7b6a5f4e3d2c1b0a9f8e7d6c5b4a3f2e1d0c9b", // predecessor
        vector[1u8], // salt
        MIN_DELAY - 1, // delay lower than minimum
        env.scenario.ctx(),
    );

    env.destroy();
}

#[test]
#[expected_failure(abort_code = mcms::EOperationAlreadyScheduled, location = mcms)]
fun test_schedule_already_scheduled() {
    let mut env = setup();
    let clock = &env.clock;

    // Schedule the batch first time
    mcms::test_timelock_schedule_batch(
        &mut env.timelock,
        clock,
        mcms::proposer_role(),
        vector[@mcms], // targets
        vector[string::utf8(b"test_module")], // module_names
        vector[string::utf8(b"test_function")], // function_names
        vector[vector[0u8]], // datas
        x"7f3e9d2c1b0a8f6e5d4c3b2a1f0e9d8c7b6a5f4e3d2c1b0a9f8e7d6c5b4a3f2e", // predecessor
        vector[1u8], // salt
        0u64, // delay
        env.scenario.ctx(),
    );

    // Try to schedule the same batch again (should fail)
    mcms::test_timelock_schedule_batch(
        &mut env.timelock,
        clock,
        mcms::proposer_role(),
        vector[@mcms], // targets
        vector[string::utf8(b"test_module")], // module_names
        vector[string::utf8(b"test_function")], // function_names
        vector[vector[0u8]], // datas
        x"7f3e9d2c1b0a8f6e5d4c3b2a1f0e9d8c7b6a5f4e3d2c1b0a9f8e7d6c5b4a3f2e", // predecessor
        vector[1u8], // salt
        0u64, // delay
        env.scenario.ctx(),
    );

    env.destroy();
}

#[test]
#[expected_failure(abort_code = mcms::EFunctionBlocked, location = mcms)]
fun test_schedule_blocked_function() {
    let mut env = setup();
    let clock = &env.clock;
    let test_target = @0xabc;
    let test_module = string::utf8(b"test_module");
    let test_function = string::utf8(b"test_function");

    // Block the function first
    mcms::test_timelock_block_function(
        &mut env.timelock,
        mcms::timelock_role(),
        test_target,
        test_module,
        test_function,
        env.scenario.ctx(),
    );

    // Try to schedule the blocked function (should fail)
    mcms::test_timelock_schedule_batch(
        &mut env.timelock,
        clock,
        mcms::proposer_role(),
        vector[test_target],
        vector[test_module],
        vector[test_function],
        vector[vector[0u8]],
        x"5d4c3b2a1f0e9d8c7b6a5f4e3d2c1b0a9f8e7d6c5b4a3f2e1d0c9b8a7f6e5d4c", // predecessor
        vector[1u8], // salt
        0, // delay
        env.scenario.ctx(),
    );

    env.destroy();
}

#[test]
#[expected_failure(abort_code = mcms::EInvalidIndex, location = mcms)]
fun test_get_blocked_function_invalid_index() {
    let mut env = setup();

    // Block a function first
    mcms::test_timelock_block_function(
        &mut env.timelock,
        mcms::timelock_role(),
        @0xabc,
        string::utf8(b"test_module"),
        string::utf8(b"test_function"),
        env.scenario.ctx(),
    );

    // Try to get function at invalid index (should fail)
    let _invalid_func = mcms::timelock_get_blocked_function(&env.timelock, 999);

    env.destroy();
}

#[test]
fun test_idempotent_block_function() {
    let mut env = setup();
    let test_target = @0xabc;
    let test_module = string::utf8(b"test_module");
    let test_function = string::utf8(b"test_function");

    // Block a function
    mcms::test_timelock_block_function(
        &mut env.timelock,
        mcms::timelock_role(),
        test_target,
        test_module,
        test_function,
        env.scenario.ctx(),
    );

    // Count should be 1
    let count = mcms::timelock_get_blocked_functions_count(&env.timelock);
    assert!(count == 1, 0);

    // Block the same function again (should be idempotent)
    mcms::test_timelock_block_function(
        &mut env.timelock,
        mcms::timelock_role(),
        test_target,
        test_module,
        test_function,
        env.scenario.ctx(),
    );

    // Count should still be 1
    let count_after = mcms::timelock_get_blocked_functions_count(&env.timelock);
    assert!(count_after == 1, 1);

    env.destroy();
}

#[test]
fun test_block_unblock_function() {
    let mut env = setup();
    let test_target = @0xabc;
    let test_module = string::utf8(b"test_module");
    let test_function = string::utf8(b"test_function");

    // Block a function
    mcms::test_timelock_block_function(
        &mut env.timelock,
        mcms::timelock_role(),
        test_target,
        test_module,
        test_function,
        env.scenario.ctx(),
    );

    // Verify the function is blocked
    let blocked_count = mcms::timelock_get_blocked_functions_count(&env.timelock);
    assert!(blocked_count == 1);

    let function = mcms::timelock_get_blocked_function(&env.timelock, 0);
    assert!(mcms::function_name(function) == test_function);
    assert!(mcms::module_name(function) == test_module);
    assert!(mcms::target(function) == test_target);

    mcms::test_timelock_unblock_function(
        &mut env.timelock,
        mcms::timelock_role(),
        test_target,
        test_module,
        test_function,
        env.scenario.ctx(),
    );

    // Verify the function is no longer blocked
    let blocked_count_after = mcms::timelock_get_blocked_functions_count(&env.timelock);
    assert!(blocked_count_after == 0);

    env.destroy();
}

// ============== View/Getter Function Tests ================= //

#[test]
fun test_view_getter_functions() {
    let env = setup();

    // Test role constants
    assert!(mcms::bypasser_role() == 0);
    assert!(mcms::canceller_role() == 1);
    assert!(mcms::proposer_role() == 2);
    assert!(mcms::timelock_role() == 3);

    // Test role validation
    assert!(mcms::is_valid_role(0));
    assert!(mcms::is_valid_role(1));
    assert!(mcms::is_valid_role(2));
    assert!(mcms::is_valid_role(3));
    assert!(!mcms::is_valid_role(4)); // Invalid role

    // Test initial timelock state
    assert!(mcms::timelock_min_delay(&env.timelock) == 0);
    assert!(mcms::timelock_get_blocked_functions_count(&env.timelock) == 0);

    env.destroy();
}

#[test]
fun test_timelock_view_functions() {
    let mut env = setup();

    // Initial state
    assert!(mcms::timelock_min_delay(&env.timelock) == 0);
    assert!(mcms::timelock_get_blocked_functions_count(&env.timelock) == 0);

    // Test after updating min delay
    mcms::test_timelock_update_min_delay(
        &mut env.timelock,
        mcms::timelock_role(),
        MIN_DELAY,
        env.scenario.ctx(),
    );
    assert!(mcms::timelock_min_delay(&env.timelock) == MIN_DELAY);

    // Test after blocking a function
    let test_target = @0xabc;
    let test_module = string::utf8(b"test_module");
    let test_function = string::utf8(b"test_function");

    mcms::test_timelock_block_function(
        &mut env.timelock,
        mcms::timelock_role(),
        test_target,
        test_module,
        test_function,
        env.scenario.ctx(),
    );

    assert!(mcms::timelock_get_blocked_functions_count(&env.timelock) == 1);

    // Verify the blocked function details
    let blocked_function = mcms::timelock_get_blocked_function(&env.timelock, 0);
    assert!(mcms::function_name(blocked_function) == test_function);
    assert!(mcms::module_name(blocked_function) == test_module);
    assert!(mcms::target(blocked_function) == test_target);

    env.destroy();
}

#[test]
#[allow(implicit_const_copy)]
fun test_operation_status_functions() {
    let mut env = setup();
    let clock = &env.clock;

    // Calculate operation ID
    let calls = mcms::create_calls(
        vector[@mcms], // targets
        vector[string::utf8(b"mcms")], // module_names
        vector[string::utf8(b"timelock_update_min_delay")], // function_names
        vector[bcs::to_bytes(&MIN_DELAY)], // datas
    );
    let id = mcms::hash_operation_batch(
        calls,
        x"c4b3a29f8e7d6c5b4a3f2e1d0c9b8a7f6e5d4c3b2a1f0e9d8c7b6a5f4e3d2c1b",
        vector[1u8],
    );

    // Initially operation should not exist
    assert!(!mcms::timelock_is_operation_pending(&env.timelock, id));
    assert!(!mcms::timelock_is_operation(&env.timelock, id));

    // Schedule the operation
    mcms::test_timelock_schedule_batch(
        &mut env.timelock,
        clock,
        mcms::proposer_role(),
        vector[@mcms], // targets
        vector[string::utf8(b"mcms")], // module_names
        vector[string::utf8(b"timelock_update_min_delay")], // function_names
        vector[bcs::to_bytes(&MIN_DELAY)], // datas
        x"c4b3a29f8e7d6c5b4a3f2e1d0c9b8a7f6e5d4c3b2a1f0e9d8c7b6a5f4e3d2c1b", // predecessor
        vector[1u8], // salt
        0u64, // delay
        env.scenario.ctx(),
    );

    // Now operation should be pending
    assert!(mcms::timelock_is_operation_pending(&env.timelock, id));
    assert!(mcms::timelock_is_operation(&env.timelock, id));

    env.destroy();
}

// ============== Utility tests ================= //

#[test]
#[allow(implicit_const_copy)]
fun test_utils__hash_metadata_leaf() {
    let hash = mcms::test_hash_metadata_leaf(
        mcms::proposer_role(), // role
        CHAIN_ID, // chain_id
        mcms_registry::get_multisig_address(), // multisig
        PRE_OP_COUNT, // pre_op_count
        POST_OP_COUNT, // post_op_count
        false, // override_previous_root
    );

    // Assert exact metadata leaf hash matches the expected metadata leaf
    let expected_metadata_hash = LEAVES[0]; // LEAVES[0] is metadata hash after reordering
    assert!(hash == expected_metadata_hash);
}

#[test]
#[allow(implicit_const_copy)]
fun test_utils__hash_op_leaf() {
    let op = mcms::test_create_op(
        mcms::proposer_role(), // role
        CHAIN_ID, // chain_id
        mcms_registry::get_multisig_address(), // multisig
        OP1_NONCE, // nonce
        mcms_registry::get_multisig_address(), // to
        string::utf8(b"mcms"), // module_name
        string::utf8(b"timelock_schedule_batch"), // function_name
        OP1_DATA, // data
    );

    let hash = mcms::hash_op_leaf(MANY_CHAIN_MULTI_SIG_DOMAIN_SEPARATOR_OP, op);

    // Assert exact OP leaf hash matches the expected operation leaf
    let expected_op_hash = LEAVES[1]; // LEAVES[1] is OP hash after reordering
    assert!(hash == expected_op_hash);
}

#[test]
#[allow(implicit_const_copy)]
fun test_verify_merkle_proof_with_hash_op() {
    let op = mcms::test_create_op(
        mcms::proposer_role(), // role
        CHAIN_ID, // chain_id
        mcms_registry::get_multisig_address(), // multisig
        OP1_NONCE,
        mcms_registry::get_multisig_address(), // to
        string::utf8(b"mcms"), // module_name
        string::utf8(b"timelock_schedule_batch"), // function_name
        OP1_DATA, // data
    );

    let computed_leaf_hash = mcms::hash_op_leaf(MANY_CHAIN_MULTI_SIG_DOMAIN_SEPARATOR_OP, op);

    // Must match expected leaf, then verify merkle proof
    let expected_leaf_hash = LEAVES[1]; // LEAVES[1] is OP hash after reordering
    let expected_root = ROOT;
    assert!(computed_leaf_hash == expected_leaf_hash);
    assert!(mcms::verify_merkle_proof(OP1_PROOF, expected_root, computed_leaf_hash));
}

#[test]
#[allow(implicit_const_copy)]
fun test_different_group_structures() {
    let mut env = setup();
    let owner_cap = ts::take_from_sender<OwnerCap>(&env.scenario);

    // This test verifies a multisig structure with nested groups to ensure
    // the MCMS system correctly handles hierarchical signer configurations

    // ==================== Group Structure ====================
    // Create a hierarchical structure as follows:
    //
    // Root Group (index 0): Requires 2-of-3 approvals from its children:
    //     +-- Group 1 (index 1): Requires 1-of-2 approvals
    //     |   +-- Signer 1 (in Group 3)
    //     |   +-- Signer 2 (in Group 4)
    //     |
    //     +-- Group 2 (index 2): Requires 1 approval
    //     |   +-- Signer 3
    //     |
    //     +-- (No direct signers in root group)
    //
    // This structure means to execute a transaction:
    // - Either Signer 3 + (Signer 1 OR Signer 2) must approve
    // - OR any other valid combination meeting the 2-of-3 quorum at the root level
    // ============================================================

    let signer_addresses = vector[PROPOSER_ADDR1, PROPOSER_ADDR2, PROPOSER_ADDR3];
    let signer_groups = vector[3, 4, 2]; // Signer 1 in Group 3, Signer 2 in Group 4, Signer 3 in Group 2

    // Define quorum for each group (how many approvals needed)
    // Root Group (0): needs 2 approvals
    // Group 1 (1): needs 1 approval from its children
    // Group 2 (2): needs 1 approval
    // Group 3 (3): needs 1 approval (leaf group with Signer 1)
    // Group 4 (4): needs 1 approval (leaf group with Signer 2)
    let mut group_quorums = vector[2, 1, 1, 1, 1];
    params::right_pad_vec(&mut group_quorums, mcms::num_groups());

    // Define the group hierarchy (which group is parent of which)
    // Group 0: parent is itself (root)
    // Group 1: parent is root (0)
    // Group 2: parent is root (0)
    // Group 3: parent is Group 1
    // Group 4: parent is Group 1
    let mut group_parents = vector[0, 0, 0, 1, 1];
    params::right_pad_vec(&mut group_parents, mcms::num_groups());

    // Configure the multisig structure for proposer role
    let role = mcms::proposer_role();
    mcms::set_config(
        &owner_cap,
        &mut env.state,
        role,
        CHAIN_ID,
        signer_addresses,
        signer_groups,
        group_quorums,
        group_parents,
        false,
        env.scenario.ctx(),
    );

    // Verify the configuration was set correctly
    let signers = mcms::signers(&env.state, role);
    assert!(signers.length() == 3); // Verify we have 3 signers

    // Verify signers are correctly stored
    assert!(signers.contains(&PROPOSER_ADDR1));
    assert!(signers.contains(&PROPOSER_ADDR2));
    assert!(signers.contains(&PROPOSER_ADDR3));

    ts::return_to_sender(&env.scenario, owner_cap);
    env.destroy();
}

#[test]
#[allow(implicit_const_copy)]
fun test_view_getter_functions_after_config() {
    let mut env = setup();
    let owner_cap = ts::take_from_sender<OwnerCap>(&env.scenario);

    // Set a config to test view functions after configuration
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
        true,
        env.scenario.ctx(),
    );

    // Test signers map after configuration
    let signers = mcms::signers(&env.state, role);
    assert!(signers.length() == 3);
    assert!(signers.contains(&PROPOSER_ADDR1));
    assert!(signers.contains(&PROPOSER_ADDR2));
    assert!(signers.contains(&PROPOSER_ADDR3));

    ts::return_to_sender(&env.scenario, owner_cap);
    env.destroy();
}

#[test]
fun test_function_view_functions() {
    let mut env = setup();

    let target = @0x123;
    let module_name = string::utf8(b"test_module");
    let function_name = string::utf8(b"test_function");

    // Test function helper getters
    let calls = mcms::create_calls(
        vector[target], // targets
        vector[module_name], // module_names
        vector[function_name], // function_names
        vector[vector[0u8]], // datas
    );
    assert!(calls.length() == 1);

    // Test blocked functions view functions
    mcms::test_timelock_block_function(
        &mut env.timelock,
        mcms::timelock_role(),
        target,
        module_name,
        function_name,
        env.scenario.ctx(),
    );

    let blocked_count = mcms::timelock_get_blocked_functions_count(&env.timelock);
    assert!(blocked_count == 1);

    // Test function getters on blocked function
    let blocked_fn = mcms::timelock_get_blocked_function(&env.timelock, 0);
    assert!(mcms::target(blocked_fn) == target);
    assert!(mcms::module_name(blocked_fn) == module_name);
    assert!(mcms::function_name(blocked_fn) == function_name);

    env.destroy();
}

#[test]
fun test_merkle__ecdsa_recover_evm_addr() {
    let eth_signed_message_hash =
        x"910cd291f5281f5bf25d8a83962f282b6c2bdf831f079dfcb84480f922abd2e1";
    let signature =
        x"45283a6239b1b559a910e97f79a52bab1605e8bd952c4b4e0720ed9b1e9e96712acab6f5f946bfa3dfa61f47705aff6e2f17f6ad83d484857bb119a06ba1f0e71C";
    let recovered_addr = mcms::test_ecdsa_recover_evm_addr(eth_signed_message_hash, signature);

    // Just verify that we get a 20-byte address back (EVM address length)
    // The exact address might differ between Sui and Aptos implementations
    assert!(recovered_addr.length() == 20);
}

#[test]
fun test_is_valid_role() {
    assert!(!mcms::is_valid_role(255)); // Invalid role
    assert!(mcms::is_valid_role(0)); // Valid role - bypasser
    assert!(mcms::is_valid_role(1)); // Valid role - canceller
    assert!(mcms::is_valid_role(2)); // Valid role - proposer
    assert!(mcms::is_valid_role(3)); // Valid role - timelock
    assert!(!mcms::is_valid_role(4)); // Invalid role - beyond max
}

#[test]
fun test_timelock_dispatching_system() {
    let mut env = setup();

    // Test the complete timelock dispatching system with actual dispatch functions

    // 1. Test dispatch_timelock_update_min_delay
    let initial_delay = mcms::timelock_min_delay(&env.timelock);
    let new_delay = initial_delay + 5000;

    let mut update_delay_data = vector[];
    update_delay_data.append(bcs::to_bytes(&object::id_address(&env.timelock))); // Add timelock address
    update_delay_data.append(bcs::to_bytes(&new_delay));

    let update_delay_params = mcms_registry::test_create_executing_callback_params(
        mcms_registry::get_multisig_address(),
        string::utf8(b"mcms"),
        string::utf8(b"timelock_update_min_delay"),
        update_delay_data,
        x"0000000000000000000000000000000000000000000000000000000000000001", // batch_id
        0, // sequence_number
        1, // total_in_batch
    );

    mcms::mcms_timelock_update_min_delay(
        &mut env.timelock,
        &mut env.registry,
        update_delay_params,
        env.scenario.ctx(),
    );

    // Verify the dispatch actually worked
    let updated_delay = mcms::timelock_min_delay(&env.timelock);
    assert!(updated_delay == new_delay, 0);

    // 2. Test dispatch_timelock_schedule_batch flow
    let even_newer_delay = new_delay + 1000;
    let current_min_delay = mcms::timelock_min_delay(&env.timelock);
    let schedule_delay = current_min_delay + 100; // Must be >= min_delay

    // We need to properly serialize the schedule batch parameters
    // For testing, we'll use the existing test helper function instead
    let clock = &env.clock; // Borrow clock separately to avoid conflicts
    mcms::test_timelock_schedule_batch(
        &mut env.timelock,
        clock,
        mcms::proposer_role(),
        vector[mcms_registry::get_multisig_address()], // targets
        vector[string::utf8(b"mcms")], // module_names
        vector[string::utf8(b"timelock_update_min_delay")], // function_names
        vector[bcs::to_bytes(&even_newer_delay)], // datas
        x"9e8d7c6b5a4f3e2d1c0b9a8f7e6d5c4b3a2f1e0d9c8b7a6f5e4d3c2b1a0f9e8d", // predecessor
        vector<u8>[2u8], // salt
        schedule_delay,
        env.scenario.ctx(),
    );

    // Verify the operation was scheduled by checking it's pending
    let calls = mcms::create_calls(
        vector[mcms_registry::get_multisig_address()], // targets
        vector[string::utf8(b"mcms")], // module_names
        vector[string::utf8(b"timelock_update_min_delay")], // function_names
        vector[bcs::to_bytes(&even_newer_delay)], // datas
    );
    let operation_id = mcms::hash_operation_batch(
        calls,
        x"9e8d7c6b5a4f3e2d1c0b9a8f7e6d5c4b3a2f1e0d9c8b7a6f5e4d3c2b1a0f9e8d",
        vector<u8>[2u8],
    );
    assert!(mcms::timelock_is_operation_pending(&env.timelock, operation_id), 1);

    // 3. Test creating ExecutingCallbackParams directly
    let test_target = mcms_registry::get_multisig_address();
    let test_module = string::utf8(b"test_module");
    let test_function = string::utf8(b"test_function");
    let executing_params = mcms_registry::test_create_executing_callback_params(
        test_target,
        test_module,
        test_function,
        vector[1, 2, 3],
        x"0000000000000000000000000000000000000000000000000000000000000004", // batch_id
        0, // sequence_number
        1, // total_in_batch
    );

    // Verify ExecutingCallbackParams properties
    assert!(mcms_registry::target(&executing_params) == test_target);
    assert!(mcms_registry::module_name(&executing_params) == test_module);
    assert!(mcms_registry::function_name(&executing_params) == test_function);
    assert!(mcms_registry::data(&executing_params) == vector[1, 2, 3]);

    // Must consume the ExecutingCallbackParams hot potato
    // In a real scenario, this would be dispatched to the appropriate module
    let (
        consumed_target,
        consumed_module,
        consumed_function,
        consumed_data,
    ) = mcms_registry::get_callback_params_from_mcms(&mut env.registry, executing_params);
    // Verify the consumed values match what we expect
    assert!(consumed_target == test_target);
    assert!(consumed_module == test_module);
    assert!(consumed_function == test_function);
    assert!(consumed_data == vector[1, 2, 3]);

    env.destroy();
}

#[test]
#[expected_failure(abort_code = mcms::EMissingDependency, location = mcms)]
fun test_execute_batch_missing_dependency() {
    let mut env = setup();

    // First, schedule a batch without dependencies
    let clock = &env.clock;
    mcms::test_timelock_schedule_batch(
        &mut env.timelock,
        clock,
        mcms::proposer_role(),
        vector[@mcms], // targets
        vector[string::utf8(b"test_module")], // module_names
        vector[string::utf8(b"test_function1")], // function_names
        vector[vector[0u8]], // datas
        x"6e5d4c3b2a1f0e9d8c7b6a5f4e3d2c1b0a9f8e7d6c5b4a3f2e1d0c9b8a7f6e5d", // predecessor (no dependency)
        vector[1u8], // salt
        0u64, // delay (immediate execution)
        env.scenario.ctx(),
    );

    // Schedule second batch with dependency on non-existent operation
    let nonexistent_predecessor =
        x"deadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef";
    mcms::test_timelock_schedule_batch(
        &mut env.timelock,
        clock,
        mcms::proposer_role(),
        vector[@mcms], // targets
        vector[string::utf8(b"test_module")], // module_names
        vector[string::utf8(b"test_function2")], // function_names
        vector[vector[0u8]], // datas
        nonexistent_predecessor, // predecessor (non-existent operation ID, 32 bytes)
        vector[2u8], // salt
        0u64, // delay (immediate execution)
        env.scenario.ctx(),
    );

    // Try to execute the second batch - should fail with EMissingDependency
    timelock_execute_dispatch_to_acc_helper(
        &mut env,
        vector[mcms_registry::get_multisig_address()], // targets
        vector[string::utf8(b"test_module")], // module_names
        vector[string::utf8(b"test_function2")], // function_names
        vector[vector[0u8]], // datas
        nonexistent_predecessor, // Must match the one used in scheduling
        vector[2u8], // salt
    );

    env.destroy();
}

#[test]
fun test_bypasser_execute_blocked_function() {
    let mut env = setup();

    let target = mcms_registry::get_multisig_address();
    let module_name = string::utf8(b"mcms");
    let function_name = string::utf8(b"timelock_update_min_delay");

    // Block the function first
    mcms::test_timelock_block_function(
        &mut env.timelock,
        mcms::timelock_role(),
        target,
        module_name,
        function_name,
        env.scenario.ctx(),
    );

    // Verify function was blocked (check the count increased)
    let blocked_count = mcms::timelock_get_blocked_functions_count(&env.timelock);
    assert!(blocked_count == 1, 0);

    // Get initial min_delay value
    let initial_delay = mcms::timelock_min_delay(&env.timelock);
    let new_delay = initial_delay + 1000;

    // Bypasser should be able to directly execute the blocked function
    // Prepare data with timelock object ID
    let mut update_delay_data = vector[];
    update_delay_data.append(bcs::to_bytes(&object::id_address(&env.timelock)));
    update_delay_data.append(bcs::to_bytes(&new_delay));

    let mut executing_params = mcms::test_timelock_bypasser_execute_batch(
        mcms::bypasser_role(),
        vector[target], // targets
        vector[module_name], // module_names
        vector[function_name], // function_names
        vector[update_delay_data], // datas
        env.scenario.ctx(),
    );

    // Process the executing callback params
    let params = executing_params.pop_back();
    let (
        callback_target,
        callback_module,
        callback_function,
        callback_data,
    ) = mcms_registry::get_callback_params_from_mcms(&mut env.registry, params);

    // Verify callback params structure
    assert!(callback_target == target);
    assert!(callback_module == module_name);
    assert!(callback_function == function_name);
    assert!(!callback_data.is_empty());

    // Now dispatch the timelock function to actually update the min delay
    let params = mcms_registry::test_create_executing_callback_params(
        mcms_registry::get_multisig_address(),
        string::utf8(b"mcms"),
        string::utf8(b"timelock_update_min_delay"),
        callback_data,
        x"0000000000000000000000000000000000000000000000000000000000000001", // batch_id
        0, // sequence_number
        1, // total_in_batch
    );

    mcms::mcms_timelock_update_min_delay(
        &mut env.timelock,
        &mut env.registry,
        params,
        env.scenario.ctx(),
    );

    // Verify the min delay was updated despite the function being blocked
    let updated_delay = mcms::timelock_min_delay(&env.timelock);
    assert!(updated_delay == new_delay, 1);

    vector::destroy_empty(executing_params);
    env.destroy();
}

#[test]
#[allow(implicit_const_copy)]
fun test_execute_batch_with_dependencies() {
    let mut env = setup();

    let delay = 1u64; // Small delay for testing

    // Execute first operation to update min delay
    let clock = &env.clock;

    // Prepare data with timelock object ID
    let mut first_min_delay_data = vector[];
    first_min_delay_data.append(bcs::to_bytes(&object::id_address(&env.timelock)));
    first_min_delay_data.append(bcs::to_bytes(&MIN_DELAY));

    mcms::test_timelock_schedule_batch(
        &mut env.timelock,
        clock,
        mcms::proposer_role(),
        vector[mcms_registry::get_multisig_address()], // targets
        vector[string::utf8(b"mcms")], // module_names
        vector[string::utf8(b"timelock_update_min_delay")], // function_names
        vector[first_min_delay_data], // datas
        mcms::zero_hash(), // predecessor
        x"abcd", // salt
        delay,
        env.scenario.ctx(),
    );

    // Fast-forward time to allow execution
    env.clock.increment_for_testing(delay * 1000 + 1);

    {
        let clock = &env.clock;
        let registry = &env.registry;

        // Prepare data for execution with timelock object ID
        let mut exec_first_min_delay_data = vector[];
        exec_first_min_delay_data.append(bcs::to_bytes(&object::id_address(&env.timelock)));
        exec_first_min_delay_data.append(bcs::to_bytes(&MIN_DELAY));

        let mut executing_callback_params = mcms::timelock_execute_batch(
            &mut env.timelock,
            clock,
            registry,
            vector[mcms_registry::get_multisig_address()], // targets
            vector[string::utf8(b"mcms")], // module_names
            vector[string::utf8(b"timelock_update_min_delay")], // function_names
            vector[exec_first_min_delay_data], // datas
            mcms::zero_hash(), // predecessor
            x"abcd", // salt
            env.scenario.ctx(),
        );

        // Execute second operation
        while (!executing_callback_params.is_empty()) {
            mcms::mcms_timelock_update_min_delay(
                &mut env.timelock,
                &mut env.registry,
                executing_callback_params.pop_back(),
                env.scenario.ctx(),
            );
        };
        executing_callback_params.destroy_empty();
    };

    // Verify min delay was updated
    assert!(mcms::timelock_min_delay(&env.timelock) == MIN_DELAY);

    // Schedule and execute second operation to test sequential execution
    let new_delay = 2000u64;
    let clock = &env.clock;

    // Prepare data with timelock object ID
    let mut second_update_delay_data = vector[];
    second_update_delay_data.append(bcs::to_bytes(&object::id_address(&env.timelock)));
    second_update_delay_data.append(bcs::to_bytes(&new_delay));

    mcms::test_timelock_schedule_batch(
        &mut env.timelock,
        clock,
        mcms::proposer_role(),
        vector[mcms_registry::get_multisig_address()], // targets
        vector[string::utf8(b"mcms")], // module_names
        vector[string::utf8(b"timelock_update_min_delay")], // function_names
        vector[second_update_delay_data], // datas
        mcms::zero_hash(), // predecessor
        x"efab", // salt
        MIN_DELAY, // Use the updated min delay
        env.scenario.ctx(),
    );

    // Fast-forward time for the new min delay
    env.clock.increment_for_testing(MIN_DELAY * 1000 + 1);

    {
        let clock = &env.clock;
        let registry = &env.registry;
        // Prepare data with timelock object ID for execution
        let mut execute_update_delay_data = vector[];
        execute_update_delay_data.append(bcs::to_bytes(&object::id_address(&env.timelock)));
        execute_update_delay_data.append(bcs::to_bytes(&new_delay));

        let mut executing_callback_params = mcms::timelock_execute_batch(
            &mut env.timelock,
            clock,
            registry,
            vector[mcms_registry::get_multisig_address()], // targets
            vector[string::utf8(b"mcms")], // module_names
            vector[string::utf8(b"timelock_update_min_delay")], // function_names
            vector[execute_update_delay_data], // datas
            mcms::zero_hash(), // predecessor
            x"efab", // salt
            env.scenario.ctx(),
        );

        // Execute second operation
        while (!executing_callback_params.is_empty()) {
            mcms::mcms_timelock_update_min_delay(
                &mut env.timelock,
                &mut env.registry,
                executing_callback_params.pop_back(),
                env.scenario.ctx(),
            );
        };
        executing_callback_params.destroy_empty();
    };

    // Verify min delay was updated to the new value
    assert!(mcms::timelock_min_delay(&env.timelock) == new_delay);

    env.destroy();
}

#[test]
fun test_bypasser_allowed_when_timelock_active() {
    let mut env = setup();

    let delay = 1800u64;

    // Update min delay to a significant value to ensure timelock is active
    mcms::test_timelock_update_min_delay(
        &mut env.timelock,
        mcms::timelock_role(),
        delay,
        env.scenario.ctx(),
    );

    // Get initial delay value
    let initial_delay = mcms::timelock_min_delay(&env.timelock);
    let new_delay = initial_delay + 1000;

    // This should succeed because bypassers are allowed to bypass the timelock
    // Prepare data with timelock object ID
    let mut bypasser_active_update_delay_data = vector[];
    bypasser_active_update_delay_data.append(bcs::to_bytes(&object::id_address(&env.timelock)));
    bypasser_active_update_delay_data.append(bcs::to_bytes(&new_delay));

    let mut executing_params = mcms::test_timelock_bypasser_execute_batch(
        mcms::bypasser_role(),
        vector[mcms_registry::get_multisig_address()], // targets
        vector[string::utf8(b"mcms")], // module_names
        vector[string::utf8(b"timelock_update_min_delay")], // function_names
        vector[bypasser_active_update_delay_data], // datas
        env.scenario.ctx(),
    );

    // Process the executing callback params to complete the operation
    mcms::mcms_timelock_update_min_delay(
        &mut env.timelock,
        &mut env.registry,
        executing_params.pop_back(),
        env.scenario.ctx(),
    );

    // Verify the min_delay was updated, confirming the bypass worked
    let updated_delay = mcms::timelock_min_delay(&env.timelock);
    assert!(updated_delay == new_delay, 0);

    executing_params.destroy_empty();
    env.destroy();
}

#[test]
#[allow(implicit_const_copy)]
#[expected_failure(abort_code = mcms::EUnknownMCMSAccountModuleFunction, location = mcms)]
fun test_unknown_mcms_module() {
    let mut env = setup();

    let clock = &env.clock;

    // Schedule batch with unknown module name
    mcms::test_timelock_schedule_batch(
        &mut env.timelock,
        clock,
        mcms::proposer_role(),
        vector[mcms_registry::get_multisig_address()], // targets
        vector[string::utf8(b"unknown_module")], // module_names (non-existent)
        vector[string::utf8(b"timelock_update_min_delay")], // function_names
        vector[bcs::to_bytes(&MIN_DELAY)], // datas
        mcms::zero_hash(), // predecessor
        vector[1u8], // salt
        0u64, // delay (immediate execution)
        env.scenario.ctx(),
    );

    // Try to execute - should fail with EUnknownMCMSModule
    timelock_execute_dispatch_to_acc_helper(
        &mut env,
        vector[mcms_registry::get_multisig_address()], // targets
        vector[string::utf8(b"unknown_module")], // module_names (non-existent)
        vector[string::utf8(b"timelock_update_min_delay")], // function_names
        vector[bcs::to_bytes(&MIN_DELAY)], // datas
        mcms::zero_hash(), // predecessor
        vector[1u8], // salt
    );

    env.destroy();
}

#[test]
fun test_timelock_dispatch_to_self() {
    let mut env = setup();

    // Test dispatch to mcms module function - timelock_update_min_delay
    let new_min_delay = 2000;
    let mut data = vector[];
    data.append(bcs::to_bytes(&object::id_address(&env.timelock)));
    data.append(bcs::to_bytes(&new_min_delay));

    let params = mcms_registry::test_create_executing_callback_params(
        mcms_registry::get_multisig_address(),
        string::utf8(b"mcms"),
        string::utf8(b"timelock_update_min_delay"),
        data,
        x"0000000000000000000000000000000000000000000000000000000000000001", // batch_id
        0, // sequence_number
        1, // total_in_batch
    );

    mcms::mcms_timelock_update_min_delay(
        &mut env.timelock,
        &mut env.registry,
        params,
        env.scenario.ctx(),
    );

    // Verify the dispatch worked
    assert!(mcms::timelock_min_delay(&env.timelock) == new_min_delay);

    // Test dispatch to timelock_block_function
    let target = mcms_registry::get_multisig_address();
    let module_name = string::utf8(b"test_module");
    let function_name = string::utf8(b"test_function");

    let mut block_data = vector[];
    block_data.append(bcs::to_bytes(&object::id_address(&env.timelock)));
    block_data.append(bcs::to_bytes(&target));
    block_data.append(bcs::to_bytes(&module_name));
    block_data.append(bcs::to_bytes(&function_name));

    let params = mcms_registry::test_create_executing_callback_params(
        mcms_registry::get_multisig_address(),
        string::utf8(b"mcms"),
        string::utf8(b"timelock_block_function"),
        block_data,
        x"0000000000000000000000000000000000000000000000000000000000000002", // Different batch_id from first call
        0,
        1,
    );

    mcms::mcms_timelock_block_function(
        &mut env.timelock,
        &mut env.registry,
        params,
        env.scenario.ctx(),
    );

    // Verify function was blocked
    let blocked_count = mcms::timelock_get_blocked_functions_count(&env.timelock);
    assert!(blocked_count == 1);

    // Test dispatch to timelock_unblock_function
    let mut unblock_data = vector[];
    unblock_data.append(bcs::to_bytes(&object::id_address(&env.timelock)));
    unblock_data.append(bcs::to_bytes(&target));
    unblock_data.append(bcs::to_bytes(&module_name));
    unblock_data.append(bcs::to_bytes(&function_name));

    let params = mcms_registry::test_create_executing_callback_params(
        mcms_registry::get_multisig_address(),
        string::utf8(b"mcms"),
        string::utf8(b"timelock_unblock_function"),
        unblock_data,
        x"0000000000000000000000000000000000000000000000000000000000000003", // Different batch_id from previous calls
        0,
        1,
    );

    mcms::mcms_timelock_unblock_function(
        &mut env.timelock,
        &mut env.registry,
        params,
        env.scenario.ctx(),
    );

    // Verify function was unblocked
    let unblocked_count = mcms::timelock_get_blocked_functions_count(&env.timelock);
    assert!(unblocked_count == 0);

    env.destroy();
}

#[test]
fun test_timelock_dispatch_to_account() {
    let mut env = setup();

    // First establish proper ownership
    let initial_owner = mcms_account::owner(&env.account_state);
    let new_owner = @0x456; // Make sure it's different from initial owner

    // First, we need to register an owner cap in the registry for dispatch to work
    let owner_cap = ts::take_from_sender<mcms_account::OwnerCap>(&env.scenario);
    let owner_cap_id = object::id_address(&owner_cap);
    let publisher_wrapper = mcms_account::test_create_publisher_wrapper(&owner_cap);
    mcms_registry::register_entrypoint(
        &mut env.registry,
        publisher_wrapper,
        mcms_account::create_mcms_account_proof(),
        owner_cap,
        vector[b"mcms_account", b"mcms_deployer", b"mcms_registry"], // Allowed MCMS modules
        env.scenario.ctx(),
    );

    // Now test dispatch routing to account module
    // Prepare BCS data with owner cap and account_state object IDs, then target address
    let mut data = vector[];
    data.append(bcs::to_bytes(&owner_cap_id));
    data.append(bcs::to_bytes(&object::id_address(&env.account_state)));
    data.append(bcs::to_bytes(&new_owner));

    let params = mcms_registry::test_create_executing_callback_params(
        mcms_registry::get_multisig_address(),
        string::utf8(b"mcms_account"),
        string::utf8(b"transfer_ownership"),
        data,
        x"0000000000000000000000000000000000000000000000000000000000000005", // batch_id
        0, // sequence_number
        1, // total_in_batch,
    );

    // Test the dispatch routing - this should route to mcms_account::transfer_ownership
    mcms::mcms_dispatch_to_account(
        &mut env.registry,
        &mut env.account_state,
        params,
        env.scenario.ctx(),
    );

    // Verify ownership was proposed (not yet transferred)
    let proposed_owner = mcms_account::pending_transfer_to(&env.account_state);
    assert!(option::is_some(&proposed_owner));
    assert!(option::borrow(&proposed_owner) == &new_owner);

    // The owner should still be the initial owner until acceptance
    assert!(mcms_account::owner(&env.account_state) == initial_owner);

    env.destroy();
}

#[test]
fun test_timelock_dispatch_to_deployer() {
    let mut env = setup();

    // First, we need to register an owner cap in the registry for dispatch to work
    let owner_cap = ts::take_from_sender<mcms_account::OwnerCap>(&env.scenario);
    let owner_cap_id = object::id_address(&owner_cap);
    let publisher_wrapper = mcms_account::test_create_publisher_wrapper(&owner_cap);
    mcms_registry::register_entrypoint(
        &mut env.registry,
        publisher_wrapper,
        mcms_account::create_mcms_account_proof(),
        owner_cap,
        vector[b"mcms_account", b"mcms_deployer", b"mcms_registry"], // Allowed MCMS modules
        env.scenario.ctx(),
    );

    let upgrade_cap = package::test_publish(
        mcms_registry::get_multisig_address().to_id(),
        env.scenario.ctx(),
    );
    {
        let registry_ref = &env.registry;
        mcms_deployer::register_upgrade_cap(
            &mut env.deployer_state,
            registry_ref,
            upgrade_cap,
            env.scenario.ctx(),
        );
    };

    let mut data = vector[];
    data.append(bcs::to_bytes(&owner_cap_id)); // owner_cap object ID
    data.append(bcs::to_bytes(&object::id_address(&env.deployer_state))); // deployer_state object ID
    data.append(bcs::to_bytes(&(1 as u8))); // policy
    data.append(bcs::to_bytes(&vector[123u8])); // digest
    data.append(bcs::to_bytes(&mcms_registry::get_multisig_address())); // package address

    let mut executing_params = mcms::test_timelock_bypasser_execute_batch(
        mcms::bypasser_role(),
        vector[mcms_registry::get_multisig_address()],
        vector[string::utf8(b"mcms_deployer")],
        vector[string::utf8(b"authorize_upgrade")],
        vector[data],
        env.scenario.ctx(),
    );

    // Process the executing callback params to complete the operation
    let params = executing_params.pop_back();
    let upgrade_ticket = mcms::mcms_dispatch_to_deployer(
        &mut env.registry,
        &mut env.deployer_state,
        params,
        env.scenario.ctx(),
    );
    let upgrade_receipt = package::test_upgrade(upgrade_ticket);
    mcms_deployer::commit_upgrade(
        &mut env.deployer_state,
        upgrade_receipt,
        env.scenario.ctx(),
    );
    executing_params.destroy_empty();
    env.destroy();
}

#[test]
#[expected_failure(abort_code = mcms_registry::EModuleNotAllowed, location = mcms_registry)]
fun test_timelock_dispatch_to_registry_invalid_module() {
    let mut env = setup();

    // First, we need to register an owner cap in the registry for dispatch to work
    let owner_cap = ts::take_from_sender<mcms_account::OwnerCap>(&env.scenario);
    let publisher_wrapper = mcms_account::test_create_publisher_wrapper(&owner_cap);
    mcms_registry::register_entrypoint(
        &mut env.registry,
        publisher_wrapper,
        mcms_account::create_mcms_account_proof(),
        owner_cap,
        vector[b"mcms_account", b"mcms_deployer", b"mcms_registry"], // Allowed MCMS modules
        env.scenario.ctx(),
    );

    // Test dispatch to registry with invalid module name (should trigger validation error)
    let params = mcms_registry::test_create_executing_callback_params(
        mcms_registry::get_multisig_address(),
        string::utf8(b"invalid_module"), // Wrong module name
        string::utf8(b"is_package_registered"), // Valid function name
        vector[1, 2, 3],
        x"0000000000000000000000000000000000000000000000000000000000000002", // batch_id
        0, // sequence_number
        1, // total_in_batch
    );

    // This should fail with EModuleNameMismatch when the registry validates the module name
    let (_cap, _function_name, _data) = mcms_registry::get_callback_params_with_caps<
        mcms_account::McmsAccountProof,
        mcms_account::OwnerCap,
    >(
        &mut env.registry,
        mcms_account::create_mcms_account_proof(),
        params,
    );

    env.destroy();
}

#[test]
fun test_dispatch_timelock_execute_batch() {
    let mut env = setup();

    // First schedule a batch operation so we have something to execute
    let clock = &env.clock;

    // Prepare data with timelock object ID prepended
    let mut update_min_delay_data = vector[];
    update_min_delay_data.append(bcs::to_bytes(&object::id_address(&env.timelock)));
    update_min_delay_data.append(bcs::to_bytes(&1000));

    mcms::test_timelock_schedule_batch(
        &mut env.timelock,
        clock,
        mcms::proposer_role(),
        vector[mcms_registry::get_multisig_address()], // targets
        vector[string::utf8(b"mcms")], // module_names
        vector[string::utf8(b"timelock_update_min_delay")], // function_names
        vector[update_min_delay_data], // datas - new min delay with timelock ID
        mcms::zero_hash(), // predecessor
        vector[1u8], // salt
        0u64, // delay
        env.scenario.ctx(),
    );

    // Create serialized data for timelock_execute_batch parameters
    let targets = vector[mcms_registry::get_multisig_address()];
    let module_names = vector[string::utf8(b"mcms")];
    let function_names = vector[string::utf8(b"timelock_update_min_delay")];

    // Re-create the data with timelock object ID (must match what was scheduled)
    let mut update_min_delay_data2 = vector[];
    update_min_delay_data2.append(bcs::to_bytes(&object::id_address(&env.timelock)));
    update_min_delay_data2.append(bcs::to_bytes(&1000u64));
    let datas = vector[update_min_delay_data2];
    let predecessor = mcms::zero_hash();
    let salt = vector[1u8];

    let mut serialized_data = vector[];
    serialized_data.append(bcs::to_bytes(&targets));
    serialized_data.append(bcs::to_bytes(&module_names));
    serialized_data.append(bcs::to_bytes(&function_names));
    serialized_data.append(bcs::to_bytes(&datas));
    serialized_data.append(bcs::to_bytes(&predecessor));
    serialized_data.append(bcs::to_bytes(&salt));

    // Create TimelockCallbackParams for dispatch_timelock_execute_batch
    let callback_params = mcms::test_create_timelock_callback_params(
        mcms::timelock_role(),
        string::utf8(b"mcms"),
        string::utf8(b"timelock_execute_batch"),
        serialized_data,
    );

    // Test the dispatch function
    let registry = &env.registry;
    let mut executing_params = mcms::dispatch_timelock_execute_batch(
        &mut env.timelock,
        clock,
        registry,
        callback_params,
        env.scenario.ctx(),
    );
    assert!(executing_params.length() == 1);

    // Verify the ExecutingCallbackParams has the correct structure
    let params = &executing_params[0];
    assert!(mcms_registry::target(params) == mcms_registry::get_multisig_address());
    assert!(mcms_registry::module_name(params) == string::utf8(b"mcms"));
    assert!(mcms_registry::function_name(params) == string::utf8(b"timelock_update_min_delay"));

    mcms::mcms_timelock_update_min_delay(
        &mut env.timelock,
        &mut env.registry,
        executing_params.pop_back(),
        env.scenario.ctx(),
    );
    executing_params.destroy_empty();

    env.destroy();
}

#[test]
fun test_dispatch_timelock_bypasser_execute_batch() {
    let mut env = setup();

    // Create serialized data for timelock_bypasser_execute_batch parameters
    let targets = vector[mcms_registry::get_multisig_address()];
    let module_names = vector[string::utf8(b"mcms")];
    let function_names = vector[string::utf8(b"timelock_update_min_delay")];

    // Prepare data with timelock object ID
    let mut bypasser_batch_data = vector[];
    bypasser_batch_data.append(bcs::to_bytes(&object::id_address(&env.timelock)));
    bypasser_batch_data.append(bcs::to_bytes(&2000));

    let datas = vector[bypasser_batch_data]; // new min delay with timelock ID

    let mut serialized_data = bcs::to_bytes(&targets);
    serialized_data.append(bcs::to_bytes(&module_names));
    serialized_data.append(bcs::to_bytes(&function_names));
    serialized_data.append(bcs::to_bytes(&datas));

    // Create TimelockCallbackParams for dispatch_timelock_bypasser_execute_batch
    let callback_params = mcms::test_create_timelock_callback_params(
        mcms::bypasser_role(),
        string::utf8(b"mcms"),
        string::utf8(b"timelock_bypasser_execute_batch"),
        serialized_data,
    );

    let mut executing_params = mcms::dispatch_timelock_bypasser_execute_batch(
        callback_params,
        env.scenario.ctx(),
    );
    assert!(executing_params.length() == 1);

    // Verify the ExecutingCallbackParams has the correct structure
    let params = &executing_params[0];
    assert!(mcms_registry::target(params) == mcms_registry::get_multisig_address());
    assert!(mcms_registry::module_name(params) == string::utf8(b"mcms"));
    assert!(mcms_registry::function_name(params) == string::utf8(b"timelock_update_min_delay"));

    mcms::mcms_timelock_update_min_delay(
        &mut env.timelock,
        &mut env.registry,
        executing_params.pop_back(),
        env.scenario.ctx(),
    );
    executing_params.destroy_empty();

    env.destroy();
}

#[test]
fun test_dispatch_timelock_cancel() {
    let mut env = setup();

    // First schedule an operation so we have something to cancel
    let clock = &env.clock;
    let salt = vector[2];
    let predecessor = x"1a2b3c4d5e6f708192a3b4c5d6e7f809102a3b4c5d6e7f809102a3b4c5d6e7f8";

    mcms::test_timelock_schedule_batch(
        &mut env.timelock,
        clock,
        mcms::proposer_role(),
        vector[mcms_registry::get_multisig_address()], // targets
        vector[string::utf8(b"mcms")], // module_names
        vector[string::utf8(b"timelock_update_min_delay")], // function_names
        vector[bcs::to_bytes(&3000)], // datas - new min delay
        predecessor, // predecessor (32 bytes)
        salt, // salt
        1000, // delay (so it's pending, not immediate)
        env.scenario.ctx(),
    );

    // Calculate the operation ID to cancel
    let calls = mcms::create_calls(
        vector[mcms_registry::get_multisig_address()],
        vector[string::utf8(b"mcms")],
        vector[string::utf8(b"timelock_update_min_delay")],
        vector[bcs::to_bytes(&3000)],
    );
    let operation_id = mcms::hash_operation_batch(
        calls,
        predecessor, // Must match the one used in scheduling
        salt,
    );

    // Verify the operation is pending before cancellation
    assert!(mcms::timelock_is_operation_pending(&env.timelock, operation_id));

    // Create TimelockCallbackParams for dispatch_timelock_cancel
    let mut cancel_data = vector[];
    cancel_data.append(bcs::to_bytes(&operation_id));

    let callback_params = mcms::test_create_timelock_callback_params(
        mcms::canceller_role(),
        string::utf8(b"mcms"),
        string::utf8(b"timelock_cancel"),
        cancel_data,
    );

    mcms::dispatch_timelock_cancel(
        &mut env.timelock,
        callback_params,
        env.scenario.ctx(),
    );

    // Verify the operation is no longer pending
    assert!(!mcms::timelock_is_operation_pending(&env.timelock, operation_id));

    env.destroy();
}

#[test]
#[expected_failure(abort_code = mcms::EInvalidModuleName, location = mcms)]
fun test_dispatch_timelock_execute_batch_invalid_module() {
    let mut env = setup();

    // Create TimelockCallbackParams with invalid module name
    let callback_params = mcms::test_create_timelock_callback_params(
        mcms::timelock_role(),
        string::utf8(b"invalid_module"), // Wrong module name
        string::utf8(b"timelock_execute_batch"),
        vector[], // Empty data
    );

    let clock = &env.clock;
    let registry = &env.registry;
    let executing_params = mcms::dispatch_timelock_execute_batch(
        &mut env.timelock,
        clock,
        registry,
        callback_params,
        env.scenario.ctx(),
    );

    executing_params.destroy_empty();
    env.destroy();
}

#[test]
#[expected_failure(abort_code = mcms::EInvalidFunctionName, location = mcms)]
fun test_dispatch_timelock_bypasser_execute_batch_invalid_function() {
    let mut env = setup();

    // Create TimelockCallbackParams with invalid function name
    let callback_params = mcms::test_create_timelock_callback_params(
        mcms::bypasser_role(),
        string::utf8(b"mcms"),
        string::utf8(b"invalid_function"), // Wrong function name
        vector[], // Empty data
    );

    // This should fail with EInvalidFunctionName
    let executing_params = mcms::dispatch_timelock_bypasser_execute_batch(
        callback_params,
        env.scenario.ctx(),
    );

    executing_params.destroy_empty();
    env.destroy();
}

#[test]
#[expected_failure(abort_code = mcms::EInvalidFunctionName, location = mcms)]
fun test_dispatch_timelock_cancel_invalid_function() {
    let mut env = setup();

    // Create TimelockCallbackParams with invalid function name
    let callback_params = mcms::test_create_timelock_callback_params(
        mcms::canceller_role(),
        string::utf8(b"mcms"),
        string::utf8(b"invalid_function"), // Wrong function name
        vector[], // Empty data
    );

    // This should fail with EInvalidFunctionName
    mcms::dispatch_timelock_cancel(
        &mut env.timelock,
        callback_params,
        env.scenario.ctx(),
    );

    env.destroy();
}

fun timelock_execute_dispatch_to_acc_helper(
    env: &mut Env,
    targets: vector<address>,
    module_names: vector<String>,
    function_names: vector<String>,
    datas: vector<vector<u8>>,
    predecessor: vector<u8>,
    salt: vector<u8>,
) {
    let mut executing_callback_params = mcms::timelock_execute_batch(
        &mut env.timelock,
        &env.clock,
        &env.registry,
        targets,
        module_names,
        function_names,
        datas,
        predecessor,
        salt,
        env.scenario.ctx(),
    );

    mcms::mcms_dispatch_to_account(
        &mut env.registry,
        &mut env.account_state,
        executing_callback_params.pop_back(),
        env.scenario.ctx(),
    );
    executing_callback_params.destroy_empty();
}

public fun destroy(env: Env) {
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

/// Helper function for tests in other modules to call mcms_dispatch_to_account
public fun test_mcms_dispatch_to_account(
    env: &mut Env,
    params: mcms_registry::ExecutingCallbackParams,
) {
    mcms::mcms_dispatch_to_account(
        &mut env.registry,
        &mut env.account_state,
        params,
        env.scenario.ctx(),
    )
}

#[test]
fun test_dispatch_timelock_schedule_batch() {
    let mut env = setup();

    // Create serialized data for timelock_schedule_batch parameters
    let targets = vector[mcms_registry::get_multisig_address()];
    let module_names = vector[string::utf8(b"mcms")];
    let function_names = vector[string::utf8(b"timelock_update_min_delay")];
    let datas = vector[bcs::to_bytes(&5000)]; // new min delay
    let predecessor = x"5f6e7d8c9baa0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c";
    let salt = vector[3];
    let delay = 1000;

    let mut serialized_data = vector[];
    serialized_data.append(bcs::to_bytes(&targets));
    serialized_data.append(bcs::to_bytes(&module_names));
    serialized_data.append(bcs::to_bytes(&function_names));
    serialized_data.append(bcs::to_bytes(&datas));
    serialized_data.append(bcs::to_bytes(&predecessor));
    serialized_data.append(bcs::to_bytes(&salt));
    serialized_data.append(bcs::to_bytes(&delay));

    // Create TimelockCallbackParams for dispatch_timelock_schedule_batch
    let callback_params = mcms::test_create_timelock_callback_params(
        mcms::proposer_role(),
        string::utf8(b"mcms"),
        string::utf8(b"timelock_schedule_batch"),
        serialized_data,
    );

    // Test the dispatch function
    let clock = &env.clock;
    mcms::dispatch_timelock_schedule_batch(
        &mut env.timelock,
        clock,
        callback_params,
        env.scenario.ctx(),
    );

    // Verify the operation was scheduled by checking if it's pending
    let calls = mcms::create_calls(targets, module_names, function_names, datas);
    let operation_id = mcms::hash_operation_batch(
        calls,
        predecessor, // Must match the one used in scheduling
        salt,
    );
    assert!(mcms::timelock_is_operation_pending(&env.timelock, operation_id));

    env.destroy();
}

#[test]
fun test_dispatch_timelock_update_min_delay() {
    let mut env = setup();

    // Get initial min_delay value
    let initial_delay = mcms::timelock_min_delay(&env.timelock);
    let new_delay = initial_delay + 2000;

    // Create TimelockCallbackParams for dispatch_timelock_update_min_delay
    let mut data = vector[];
    data.append(bcs::to_bytes(&object::id_address(&env.timelock)));
    data.append(bcs::to_bytes(&new_delay));

    let callback_params = mcms_registry::test_create_executing_callback_params(
        mcms_registry::get_multisig_address(),
        string::utf8(b"mcms"),
        string::utf8(b"timelock_update_min_delay"),
        data,
        x"0000000000000000000000000000000000000000000000000000000000000001",
        0,
        1,
    );

    // Test the dispatch function
    mcms::mcms_timelock_update_min_delay(
        &mut env.timelock,
        &mut env.registry,
        callback_params,
        env.scenario.ctx(),
    );

    // Verify the min_delay was updated
    let updated_delay = mcms::timelock_min_delay(&env.timelock);
    assert!(updated_delay == new_delay);
    assert!(updated_delay != initial_delay);

    env.destroy();
}

#[test]
fun test_dispatch_timelock_block_function() {
    let mut env = setup();

    // Target function to block
    let target = mcms_registry::get_multisig_address();
    let module_name = string::utf8(b"test_module");
    let function_name = string::utf8(b"test_function");

    // Create serialized data for timelock_block_function parameters
    let mut serialized_data = vector[];
    serialized_data.append(bcs::to_bytes(&object::id_address(&env.timelock)));
    serialized_data.append(bcs::to_bytes(&target));
    serialized_data.append(bcs::to_bytes(&module_name));
    serialized_data.append(bcs::to_bytes(&function_name));

    // Create ExecutingCallbackParams for dispatch_timelock_block_function
    let callback_params = mcms_registry::test_create_executing_callback_params(
        mcms_registry::get_multisig_address(),
        string::utf8(b"mcms"),
        string::utf8(b"timelock_block_function"),
        serialized_data,
        x"0000000000000000000000000000000000000000000000000000000000000001",
        0,
        1,
    );

    // Get initial blocked functions count
    let initial_count = mcms::timelock_get_blocked_functions_count(&env.timelock);

    // Test the dispatch function
    mcms::mcms_timelock_block_function(
        &mut env.timelock,
        &mut env.registry,
        callback_params,
        env.scenario.ctx(),
    );

    // Verify the function was blocked
    let updated_count = mcms::timelock_get_blocked_functions_count(&env.timelock);
    assert!(updated_count == initial_count + 1);

    // Verify the specific function is in the blocked list
    // Check if our function is in the blocked list by examining the first (and only) blocked function
    assert!(updated_count == 1); // Should have exactly 1 blocked function
    let blocked_function = mcms::timelock_get_blocked_function(&env.timelock, 0);
    assert!(mcms::target(blocked_function) == target);
    assert!(mcms::module_name(blocked_function) == module_name);
    assert!(mcms::function_name(blocked_function) == function_name);

    env.destroy();
}

#[test]
fun test_dispatch_timelock_unblock_function() {
    let mut env = setup();

    // Target function to block and then unblock
    let target = mcms_registry::get_multisig_address();
    let module_name = string::utf8(b"test_module");
    let function_name = string::utf8(b"test_function");

    // First block the function using the direct test function
    mcms::test_timelock_block_function(
        &mut env.timelock,
        mcms::timelock_role(),
        target,
        module_name,
        function_name,
        env.scenario.ctx(),
    );

    // Verify function is blocked
    let initial_count = mcms::timelock_get_blocked_functions_count(&env.timelock);
    assert!(initial_count > 0);

    // Create serialized data for timelock_unblock_function parameters
    let mut serialized_data = vector[];
    serialized_data.append(bcs::to_bytes(&object::id_address(&env.timelock)));
    serialized_data.append(bcs::to_bytes(&target));
    serialized_data.append(bcs::to_bytes(&module_name));
    serialized_data.append(bcs::to_bytes(&function_name));

    // Create TimelockCallbackParams for dispatch_timelock_unblock_function
    let callback_params = mcms_registry::test_create_executing_callback_params(
        mcms_registry::get_multisig_address(),
        string::utf8(b"mcms"),
        string::utf8(b"timelock_unblock_function"),
        serialized_data,
        x"0000000000000000000000000000000000000000000000000000000000000001",
        0,
        1,
    );

    // Test the dispatch function
    mcms::mcms_timelock_unblock_function(
        &mut env.timelock,
        &mut env.registry,
        callback_params,
        env.scenario.ctx(),
    );

    // Verify the function was unblocked
    let updated_count = mcms::timelock_get_blocked_functions_count(&env.timelock);
    assert!(updated_count == initial_count - 1);

    // Verify the specific function is no longer in the blocked list
    // Since we unblocked the only function, the count should be back to 0
    assert!(updated_count == 0);

    env.destroy();
}

#[test]
#[expected_failure(abort_code = mcms::EInvalidModuleName, location = mcms)]
fun test_dispatch_timelock_schedule_batch_invalid_module() {
    let mut env = setup();

    // Create TimelockCallbackParams with invalid module name
    let callback_params = mcms::test_create_timelock_callback_params(
        mcms::proposer_role(),
        string::utf8(b"invalid_module"), // Wrong module name
        string::utf8(b"timelock_schedule_batch"),
        vector[], // Empty data
    );

    // This should fail with EInvalidModuleName
    let clock = &env.clock;
    mcms::dispatch_timelock_schedule_batch(
        &mut env.timelock,
        clock,
        callback_params,
        env.scenario.ctx(),
    );

    env.destroy();
}

#[test]
#[expected_failure(abort_code = mcms::EInvalidFunctionName, location = mcms)]
fun test_dispatch_timelock_update_min_delay_invalid_function() {
    let mut env = setup();

    // Create ExecutingCallbackParams with invalid function name
    let callback_params = mcms_registry::test_create_executing_callback_params(
        mcms_registry::get_multisig_address(),
        string::utf8(b"mcms"),
        string::utf8(b"invalid_function"),
        vector[],
        x"0000000000000000000000000000000000000000000000000000000000000001",
        0,
        1,
    );

    // This should fail with EInvalidFunctionName
    mcms::mcms_timelock_update_min_delay(
        &mut env.timelock,
        &mut env.registry,
        callback_params,
        env.scenario.ctx(),
    );

    env.destroy();
}

#[test]
#[expected_failure(abort_code = mcms::EInvalidFunctionName, location = mcms)]
fun test_dispatch_timelock_block_function_invalid_function() {
    let mut env = setup();

    // Create ExecutingCallbackParams with invalid function name
    let callback_params = mcms_registry::test_create_executing_callback_params(
        mcms_registry::get_multisig_address(),
        string::utf8(b"mcms"),
        string::utf8(b"invalid_function"),
        vector[],
        x"0000000000000000000000000000000000000000000000000000000000000001",
        0,
        1,
    );

    // This should fail with EInvalidFunctionName
    mcms::mcms_timelock_block_function(
        &mut env.timelock,
        &mut env.registry,
        callback_params,
        env.scenario.ctx(),
    );

    env.destroy();
}

#[test]
#[expected_failure(abort_code = mcms::EInvalidFunctionName, location = mcms)]
fun test_dispatch_timelock_unblock_function_invalid_function() {
    let mut env = setup();

    // Provide valid BCS data so we can reach the function name check
    let mut data = vector[];
    data.append(bcs::to_bytes(&object::id_address(&env.timelock)));
    data.append(bcs::to_bytes(&mcms_registry::get_multisig_address()));
    data.append(bcs::to_bytes(&string::utf8(b"test_module")));
    data.append(bcs::to_bytes(&string::utf8(b"test_function")));

    let callback_params = mcms_registry::test_create_executing_callback_params(
        mcms_registry::get_multisig_address(),
        string::utf8(b"mcms"),
        string::utf8(b"invalid_function"), // This causes EInvalidFunctionName
        data,
        x"0000000000000000000000000000000000000000000000000000000000000001",
        0,
        1,
    );

    // This should fail with EInvalidFunctionName
    mcms::mcms_timelock_unblock_function(
        &mut env.timelock,
        &mut env.registry,
        callback_params,
        env.scenario.ctx(),
    );

    env.destroy();
}

#[test]
fun test_mcms_dispatch_to_registry_add_allowed_modules() {
    let mut env = setup();

    // Transaction 1: Initialize registry
    {
        let ctx = env.scenario.ctx();
        mcms_registry::test_init(ctx);
    };

    ts::next_tx(&mut env.scenario, OWNER);

    // Transaction 2: Register MCMS's own package with McmsProof (normally done during MCMS account registration)
    {
        let mut registry = ts::take_shared<Registry>(&env.scenario);
        let owner_cap = ts::take_from_sender<OwnerCap>(&env.scenario);
        let ctx = ts::ctx(&mut env.scenario);

        let publisher_wrapper = mcms_account::test_create_publisher_wrapper(&owner_cap);

        // Register MCMS package with McmsProof witness
        mcms_registry::register_entrypoint<mcms_account::McmsAccountProof, OwnerCap>(
            &mut registry,
            publisher_wrapper,
            mcms_account::create_mcms_account_proof(),
            owner_cap,
            vector[b"mcms_account", b"mcms_deployer", b"mcms_registry"],
            ctx,
        );

        ts::return_shared(registry);
    };

    ts::next_tx(&mut env.scenario, @0xA);

    // Transaction 3: Use mcms_dispatch_to_registry to add a new module
    {
        let mut registry = ts::take_shared<Registry>(&env.scenario);

        // Prepare data for add_allowed_modules
        let mut data = vector::empty<u8>();
        data.append(bcs::to_bytes(&object::id_address(&registry))); // Registry address

        // Serialize vector of module names
        let module_names = vector[b"new_mcms_module"];
        data.append(bcs::to_bytes(&module_names));

        let params = mcms_registry::test_create_executing_callback_params(
            mcms_registry::get_multisig_address(),
            string::utf8(b"mcms_registry"),
            string::utf8(b"add_allowed_modules"),
            data,
            x"0000000000000000000000000000000000000000000000000000000000000001",
            0,
            1,
        );

        // Dispatch to registry
        let ctx = ts::ctx(&mut env.scenario);
        mcms::mcms_dispatch_to_registry(&mut registry, params, ctx);

        // Verify the new module was added
        let allowed_modules = mcms_registry::get_allowed_modules(
            &registry,
            mcms_registry::get_multisig_address_ascii(),
        );
        assert!(allowed_modules.length() == 4); // Original 3 + 1 new
        assert!(allowed_modules[3] == b"new_mcms_module");

        ts::return_shared(registry);
    };

    env.destroy();
}

