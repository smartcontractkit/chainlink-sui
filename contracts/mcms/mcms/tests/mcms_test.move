#[test_only]
module mcms::mcms_test;

use mcms::mcms::{
    Self,
    MultisigState,
    Timelock,
    TimelockCallbackParams,
};
use mcms::mcms_account::{Self, AccountState, OwnerCap};
use mcms::mcms_deployer::{Self, DeployerState};
use mcms::mcms_registry::{Self, Registry};
use mcms::params;
use std::string::{Self, String};
use std::{bcs};
use sui::test_scenario::{Self as ts};

const OWNER: address = @0x123;

// keccak256("MANY_CHAIN_MULTI_SIG_DOMAIN_SEPARATOR_OP_SUI")
const MANY_CHAIN_MULTI_SIG_DOMAIN_SEPARATOR_OP: vector<u8> =
    x"542b28b7edb99385286abe2b9c308f91a385cbcb48fc98127cfd13deb28a50b8";

const CHAIN_ID: u256 = 2;
const TIMESTAMP: u64 = 1757637134;

const MIN_DELAY: u64 = 3600; // 1 hour delay

// Proposer signers from the logs (in ascending order)
const PROPOSER_ADDR1: vector<u8> = x"a7cec01fef1816a5e8602e5ef996190f49befaba";
const PROPOSER_ADDR2: vector<u8> = x"acff986809fceb011e8c742ea9dbb21b0576c927";
const PROPOSER_ADDR3: vector<u8> = x"ee32db57d81ca94d4622d0870a6a2a9f60fd63e0";

// test config: 2-of-3 multisig
const SIGNER_GROUPS: vector<u8> = vector[0, 0, 0];

const GROUP_QUORUMS: vector<u8> = vector[
    2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
];

const GROUP_PARENTS: vector<u8> = vector[
    0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
];

const ROOT: vector<u8> = x"89e61ce983bda948c3c485d5554730750ef512de48d4c5082ee38dc8609bdd5c";
const VALID_UNTIL: u64 = 1757723533;
const SIGNATURES: vector<vector<u8>> = vector[
    x"b66bd32cb747d33fb65a860214404cbe14a79cd448ecc6804882da15ab53a2c441b7c082638201932ad277b4f4f3370aa099255575f6d4b3e308a77916281b391c",
    x"9686c757c3f961940620bddf267acab420e1d62f3807853f3a81a3a50903d4931f12109bde9467d77a592db23f2342ca608c907390609d06f157a06206fd68301b"
];

const PRE_OP_COUNT: u64 = 0;
const POST_OP_COUNT: u64 = 1;

const METADATA_PROOF: vector<vector<u8>> = vector[
    x"5aaedd07c0bff02c6fde122e5380101680b2645c561c8df04aabdc1b925a3e95", // op hash (sibling proof for metadata)
];

const OP1_PROOF: vector<vector<u8>> = vector[
    x"ebfd00c9a23cdd55c747ed3f4522d7672a5859dab0933a1c19fa7db4b417904f", // metadata hash (sibling proof for op)
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
const LEAVES: vector<vector<u8>> = vector[
    x"ebfd00c9a23cdd55c747ed3f4522d7672a5859dab0933a1c19fa7db4b417904f", // metadata hash
    x"5aaedd07c0bff02c6fde122e5380101680b2645c561c8df04aabdc1b925a3e95", // op1 hash
];

const OP1_NONCE: u64 = 0;
const OP1_DATA: vector<u8> = x"01269d332e49310ea7bda3428dc6df523d9d647022221bbf173bf903f8f9656ff5010c6d636d735f6163636f756e74011c6163636570745f6f776e6572736869705f61735f74696d656c6f636b01002000000000000000000000000000000000000000000000000000000000000000002068c4bb8d000000000000000000000000000000000000000000000000000000000100000000000000";

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
        proof: vector<vector<u8>>
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
            proof: OP1_PROOF
        }
    }

    fun call_execute(env: &mut Env, args: ExecuteArgs): TimelockCallbackParams {
        mcms::execute(
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
        )
    }

    #[test]
    public fun test_e2e() {
        let mut env = setup();

        let role = mcms::proposer_role();
        let owner_cap = ts::take_from_sender<OwnerCap>(&env.scenario);
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

        let signers = mcms::signers(&env.state, role);
        assert!(signers.size() == 3);

        let set_root_args = default_set_root_args(false);
        call_set_root(&mut env, set_root_args);

        let (root, valid_until, op_count) = mcms::expiring_root_and_op_count(&env.state, role);
        assert!(root == ROOT);
        assert!(valid_until == VALID_UNTIL);
        assert!(op_count == 0);

        // First we must transfer ownership to `@mcms` (the multisig/self)
        mcms_account::transfer_ownership_to_self(&owner_cap, &mut env.account_state, env.scenario.ctx());

        // FIRST EXECUTE: Schedule the timelock operation via mcms::execute  
        let schedule_args = default_execute_args(); // Uses timelock_schedule_batch
        let callback_params = call_execute(&mut env, schedule_args);
        dispatch_timelock_schedule_batch_helper(&mut env, callback_params);

        // check op count incremented after scheduling
        let (_post_execute_root, _post_execute_valid_until, post_execute_op_count) =
            mcms::expiring_root_and_op_count(&env.state, role);
        assert!(post_execute_op_count == 1);

        // Wait for delay (10 second)
        env.clock.set_for_testing((TIMESTAMP * 1000) + 10000);
        
        // SECOND EXECUTE: Execute the scheduled timelock operation directly (not via mcms::execute)
        // This should call timelock_execute_batch directly like the Go test's timelockExecutable.Execute()
        let target_from_data = @0x269d332e49310ea7bda3428dc6df523d9d647022221bbf173bf903f8f9656ff5;
        timelock_execute_batch_helper(
            &mut env,
            vector[target_from_data], // targets - use address from OP1_DATA
            vector[string::utf8(b"mcms_account")], // module_names  
            vector[string::utf8(b"accept_ownership_as_timelock")], // function_names
            vector[x""], // datas (empty bytes)
            x"0000000000000000000000000000000000000000000000000000000000000000", // predecessor - 32-byte zero hash (same as scheduling)
            x"68c4bb8d00000000000000000000000000000000000000000000000000000000", // salt - same as scheduling
        );

        let ctx = env.scenario.ctx();
        mcms_account::execute_ownership_transfer(
            owner_cap,
            &mut env.account_state,
            &mut env.registry,
            mcms_registry::get_multisig_address(),
            ctx,
        );

        // Verify new owner is now `@mcms`
        let new_mcms_owner = mcms_account::owner(&env.account_state);
        assert!(new_mcms_owner == mcms_registry::get_multisig_address());

        env.destroy();
    }

fun setup(): Env {
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

#[test]
#[expected_failure(abort_code = mcms::EMissingConfig, location = mcms)]
fun test_set_root__config_not_set() {
    let mut env = setup();
    let mut set_root_args = default_set_root_args(false);
    set_root_args.signatures = vector[]; // no signatures
    call_set_root(&mut env, set_root_args);
    env.destroy();
}

#[test]
#[expected_failure(abort_code = mcms::ESignerAddrMustBeIncreasing, location = mcms)]
fun test_set_root__out_of_order_signatures() {
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
    let sig0 = set_root_args.signatures[0];
    let sig1 = set_root_args.signatures[1];
    // Reverse the order of the 2 signatures (out of order)
    set_root_args.signatures = vector[sig1, sig0]; // shuffle signature order
    call_set_root(&mut env, set_root_args);

    ts::return_to_sender(&env.scenario, owner_cap);
    env.destroy();
}

#[test]
#[expected_failure(abort_code = mcms::EInvalidSigner, location = mcms)]
fun test_set_root__signature_from_invalid_signer() {
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
    let invalid_signer_sig =
        x"bb7f7e44b8d9c8f978c255c7efd6abb64e8fa9a33dcb6db2e2203d8aacd51dd471113ca6c8d1ed56bb0395f0bef0daf2fae6ef2cb5c86c57d148c7de473383461B";
    set_root_args.signatures = vector[invalid_signer_sig]; // add signature from invalid signer
    call_set_root(&mut env, set_root_args);

    ts::return_to_sender(&env.scenario, owner_cap);
    env.destroy();
}

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
        &mut env.state, role, ROOT, VALID_UNTIL, new_op_count
    );
    mcms::test_set_root_metadata(
        &mut env.state,
        role,
        CHAIN_ID,
        mcms_registry::get_multisig_address(),
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

    // Verify the first configuration (clear_root=false)
    let signers = mcms::signers(&env.state, role);
    assert!(signers.size() == 3);
    
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
    let (root_after_clear, valid_until_after_clear, op_count_after_clear) = 
        mcms::expiring_root_and_op_count(&env.state, role);
    assert!(root_after_clear == vector[]);
    assert!(valid_until_after_clear == 0);
    assert!(op_count_after_clear == new_op_count);

    let root_metadata_after_clear = mcms::get_root_metadata(&env.state, role);
    assert!(mcms::role(&root_metadata_after_clear) == role);
    assert!(mcms::chain_id(&root_metadata_after_clear) == CHAIN_ID);
    assert!(mcms::root_metadata_multisig(&root_metadata_after_clear) == mcms_registry::get_multisig_address());
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
        &mut env.state, role, ROOT, VALID_UNTIL, post_op_count
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
        &mut env.state, role, ROOT, TIMESTAMP - 1, 0 // expired valid_until
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
    
    // Schedule a batch operation
    let targets = vector[@0x1];
    let module_names = vector[string::utf8(b"test_module")];
    let function_names = vector[string::utf8(b"test_function")]; 
    let datas = vector[b"test_data"];
    let predecessor = vector[];
    let salt = vector[];
    let delay = MIN_DELAY;
    
    // Need to borrow clock first to avoid borrowing conflict
    let clock = &env.clock;
    mcms::test_timelock_schedule_batch(
        &mut env.timelock,
        clock,
        mcms::proposer_role(),
        targets,
        module_names,
        function_names,
        datas,
        predecessor,
        salt,
        delay,
        env.scenario.ctx(),
    );
    
    env.destroy();
}

#[test]
#[allow(implicit_const_copy)]
fun test_cancel_operation() {
    let mut env = setup();
    let clock = &env.clock;
    let targets = vector[@mcms];
    let module_names = vector[string::utf8(b"mcms")];
    let function_names = vector[string::utf8(b"timelock_update_min_delay")];
    let data = bcs::to_bytes(&MIN_DELAY);
    let datas = vector[data];
    let predecessor = vector[];
    let salt = vector[1u8];
    let delay = 0u64;
    
    // Schedule the batch first
    mcms::test_timelock_schedule_batch(
        &mut env.timelock,
        clock,
        mcms::proposer_role(),
        targets,
        module_names,
        function_names,
        datas,
        predecessor,
        salt,
        delay,
        env.scenario.ctx(),
    );
    
    // Calculate the operation ID
    let calls = mcms::create_calls(targets, module_names, function_names, datas);
    let id = mcms::hash_operation_batch(calls, predecessor, salt);
    
    // Verify operation is pending
    assert!(mcms::timelock_is_operation_pending(&env.timelock, id), 0);
    
    // Cancel the operation
    mcms::test_timelock_cancel(
        &mut env.timelock,
        mcms::canceller_role(),
        id,
        env.scenario.ctx(),
    );
    
    // Verify operation is no longer pending
    assert!(!mcms::timelock_is_operation_pending(&env.timelock, id), 1);
    
    env.destroy();
}

#[test]
#[expected_failure] // EOperationNotExists - exact error code varies
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

// TODO: Implement bypasser execute batch test - requires complex callback setup
// #[test]
// fun test_bypasser_execute_batch() { ... }

// ============== Advanced Timelock Tests ================= //

#[test]
#[expected_failure] // Parameter length mismatch
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
        vector[], // predecessor
        vector[1u8], // salt
        0, // delay
        env.scenario.ctx(),
    );
    
    env.destroy();
}

#[test]
#[expected_failure] // Insufficient delay
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
        vector[], // predecessor
        vector[1u8], // salt
        MIN_DELAY - 1, // delay lower than minimum
        env.scenario.ctx(),
    );
    
    env.destroy();
}

#[test]
#[expected_failure] // Operation already scheduled
fun test_schedule_already_scheduled() {
    let mut env = setup();
    let clock = &env.clock;
    let targets = vector[@mcms];
    let module_names = vector[string::utf8(b"test_module")];
    let function_names = vector[string::utf8(b"test_function")];
    let datas = vector[vector[0u8]];
    let predecessor = vector[];
    let salt = vector[1u8];
    let delay = 0u64;
    
    // Schedule the batch first time
    mcms::test_timelock_schedule_batch(
        &mut env.timelock,
        clock,
        mcms::proposer_role(),
        targets,
        module_names,
        function_names,
        datas,
        predecessor,
        salt,
        delay,
        env.scenario.ctx(),
    );
    
    // Try to schedule the same batch again (should fail)
    mcms::test_timelock_schedule_batch(
        &mut env.timelock,
        clock,
        mcms::proposer_role(),
        targets,
        module_names,
        function_names,
        datas,
        predecessor,
        salt,
        delay,
        env.scenario.ctx(),
    );
    
    env.destroy();
}

#[test] 
#[expected_failure] // Blocked function scheduling
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
        vector[], // predecessor
        vector[1u8], // salt
        0, // delay
        env.scenario.ctx(),
    );
    
    env.destroy();
}

// TODO: Implement execute batch test - requires understanding Sui execution flow
// #[test] 
// #[expected_failure] // Operation not ready yet
// fun test_execute_batch_not_ready() { ... }

#[test]
#[expected_failure] // Invalid index for blocked function
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
    assert!(blocked_count == 1, 0);
    
    let function = mcms::timelock_get_blocked_function(&env.timelock, 0);
    assert!(mcms::function_name(function) == test_function, 1);
    assert!(mcms::module_name(function) == test_module, 2);
    assert!(mcms::target(function) == test_target, 3);
    
    // Unblock the function
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
    assert!(blocked_count_after == 0, 4);
    
    env.destroy();
}

// ============== View/Getter Function Tests ================= //

#[test]
fun test_view_getter_functions() {
    let env = setup();
    
    // Test role constants
    assert!(mcms::bypasser_role() == 0, 0);
    assert!(mcms::canceller_role() == 1, 1);
    assert!(mcms::proposer_role() == 2, 2);
    assert!(mcms::timelock_role() == 3, 3);
    
    // Test role validation
    assert!(mcms::is_valid_role(0), 4);
    assert!(mcms::is_valid_role(1), 5);
    assert!(mcms::is_valid_role(2), 6);
    assert!(mcms::is_valid_role(3), 7);
    assert!(!mcms::is_valid_role(4), 8); // Invalid role
    
    // Test initial timelock state
    assert!(mcms::timelock_min_delay(&env.timelock) == 0, 9);
    assert!(mcms::timelock_get_blocked_functions_count(&env.timelock) == 0, 10);
    
    env.destroy();
}

#[test] 
fun test_timelock_view_functions() {
    let mut env = setup();
    
    // Initial state
    assert!(mcms::timelock_min_delay(&env.timelock) == 0, 0);
    assert!(mcms::timelock_get_blocked_functions_count(&env.timelock) == 0, 1);
    
    // Test after updating min delay
    mcms::test_timelock_update_min_delay(
        &mut env.timelock,
        mcms::timelock_role(),
        MIN_DELAY,
        env.scenario.ctx(),
    );
    assert!(mcms::timelock_min_delay(&env.timelock) == MIN_DELAY, 2);
    
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
    
    assert!(mcms::timelock_get_blocked_functions_count(&env.timelock) == 1, 3);
    
    // Verify the blocked function details
    let blocked_function = mcms::timelock_get_blocked_function(&env.timelock, 0);
    assert!(mcms::function_name(blocked_function) == test_function, 4);
    assert!(mcms::module_name(blocked_function) == test_module, 5);
    assert!(mcms::target(blocked_function) == test_target, 6);
    
    env.destroy();
}

#[test]
#[allow(implicit_const_copy)]
fun test_operation_status_functions() {
    let mut env = setup();
    let clock = &env.clock;
    let targets = vector[@mcms];
    let module_names = vector[string::utf8(b"mcms")];
    let function_names = vector[string::utf8(b"timelock_update_min_delay")];
    let datas = vector[bcs::to_bytes(&MIN_DELAY)];
    let predecessor = vector[];
    let salt = vector[1u8];
    let delay = 0u64;
    
    // Calculate operation ID
    let calls = mcms::create_calls(targets, module_names, function_names, datas);
    let id = mcms::hash_operation_batch(calls, predecessor, salt);
    
    // Initially operation should not exist
    assert!(!mcms::timelock_is_operation_pending(&env.timelock, id), 0);
    assert!(!mcms::timelock_is_operation(&env.timelock, id), 1);
    
    // Schedule the operation
    mcms::test_timelock_schedule_batch(
        &mut env.timelock,
        clock,
        mcms::proposer_role(),
        targets,
        module_names,
        function_names,
        datas,
        predecessor,
        salt,
        delay,
        env.scenario.ctx(),
    );
    
    // Now operation should be pending
    assert!(mcms::timelock_is_operation_pending(&env.timelock, id), 2);
    assert!(mcms::timelock_is_operation(&env.timelock, id), 3);
    
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
        false // override_previous_root
    );
    
    // Assert exact metadata leaf hash matches the expected first leaf
    let expected_metadata_hash = LEAVES[0];
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
    
    // Assert exact OP leaf hash matches the expected second leaf
    let expected_op_hash = LEAVES[1];
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
    let expected_leaf_hash = LEAVES[1];
    let expected_root = ROOT;
    assert!(computed_leaf_hash == expected_leaf_hash);
    assert!(mcms::verify_merkle_proof(OP1_PROOF, expected_root, computed_leaf_hash));
}

// Helper function to dispatch timelock schedule batch without borrowing conflicts
fun dispatch_timelock_schedule_batch_helper(
    env: &mut Env,
    callback_params: TimelockCallbackParams,
) {
    let clock = &env.clock;
    let ctx = env.scenario.ctx();
    mcms::dispatch_timelock_schedule_batch(
        &mut env.timelock,
        clock,
        callback_params,
        ctx,
    );
}

// Helper function to execute timelock batch without borrowing conflicts  
fun timelock_execute_batch_helper(
    env: &mut Env,
    targets: vector<address>,
    module_names: vector<String>,
    function_names: vector<String>,
    datas: vector<vector<u8>>,
    predecessor: vector<u8>,
    salt: vector<u8>,
) {
    let clock = &env.clock;
    let ctx = env.scenario.ctx();
    let executing_callback_params = mcms::timelock_execute_batch(
        &mut env.timelock,
        clock,
        targets,
        module_names,
        function_names,
        datas,
        predecessor,
        salt,
        ctx,
    );
    
    // Process the callback params
    let mut params = executing_callback_params;
    while (!params.is_empty()) {
        let param = params.pop_back();
        let ctx = env.scenario.ctx();
        mcms::execute_dispatch_to_account(
            &mut env.registry,
            &mut env.account_state,
            param,
            ctx,
        );
    };
    params.destroy_empty();
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
