module ccip::rmn_remote {
    use std::bcs;
    use sui::hash;
    use sui::event;
    use sui::vec_map;
    use sui::ecdsa_k1;
    use std::string::{Self, String};

    use ccip::eth_abi;
    use ccip::state_object::{Self, OwnerCap, CCIPObjectRef};
    use ccip::merkle_proof;

    public struct RMNRemoteState has key, store {
        id: UID,
        local_chain_selector: u64,
        config: Config,
        config_count: u32,
        // most operations are O(n) with vec map, but it's easy to retrieve all the keys
        signers: vec_map::VecMap<vector<u8>, bool>,
        cursed_subjects: vec_map::VecMap<vector<u8>, bool>
    }

    public struct Config has copy, drop, store {
        rmn_home_contract_config_digest: vector<u8>,
        signers: vector<Signer>,
        f_sign: u64
    }

    public struct Signer has copy, drop, store {
        onchain_public_key: vector<u8>,
        node_index: u64
    }

    public struct Report has drop {
        // dest_chain_id: u64,
        dest_chain_selector: u64,
        rmn_remote_contract_address: address,
        off_ramp_address: address,
        rmn_home_contract_config_digest: vector<u8>,
        merkle_roots: vector<MerkleRoot>
    }

    public struct MerkleRoot has drop {
        source_chain_selector: u64,
        min_seq_nr: u64,
        max_seq_nr: u64,
        merkle_root: vector<u8>
    }

    public struct ConfigSet has copy, drop {
        version: u32,
        config: Config
    }

    public struct Cursed has copy, drop {
        subjects: vector<vector<u8>>
    }

    public struct Uncursed has copy, drop {
        subjects: vector<vector<u8>>
    }

    const SIGNATURE_NUM_BYTES: u64 = 64;
    const GLOBAL_CURSE_SUBJECT: vector<u8> = x"01000000000000000000000000000001";
    const RMN_REMOTE_STATE_NAME: vector<u8> = b"RMNRemoteState";

    const E_ALREADY_INITIALIZED: u64 = 1;
    const E_ALREADY_CURSED: u64 = 2;
    const E_CONFIG_NOT_SET: u64 = 3;
    const E_DUPLICATE_SIGNER: u64 = 4;
    const E_INVALID_SIGNATURE: u64 = 5;
    const E_INVALID_SIGNER_ORDER: u64 = 6;
    const E_NOT_ENOUGH_SIGNERS: u64 = 7;
    const E_NOT_CURSED: u64 = 8;
    const E_OUT_OF_ORDER_SIGNATURES: u64 = 9;
    const E_THRESHOLD_NOT_MET: u64 = 10;
    const E_UNEXPECTED_SIGNER: u64 = 11;
    const E_ZERO_VALUE_NOT_ALLOWED: u64 = 12;
    const E_MERKLE_ROOT_LENGTH_MISMATCH: u64 = 13;
    const E_INVALID_DIGEST_LENGTH: u64 = 14;
    const E_SIGNERS_MISMATCH: u64 = 15;
    const E_INVALID_SUBJECT_LENGTH: u64 = 16;
    const E_INVALID_PUBLIC_KEY_LENGTH: u64 = 17;
    // const E_UNKNOWN_FUNCTION: u64 = 18;

    public fun type_and_version(): String {
        string::utf8(b"RMNRemote 1.6.0")
    }

    // fun init_module(publisher: &signer) {
    //     if (@mcms_register_entrypoints != @0x0) {
    //         mcms_registry::register_entrypoint(
    //             publisher, string::utf8(b"rmn_remote"), McmsCallback {}
    //         );
    //     };
    // }

    public fun initialize(
        ownerCap: &OwnerCap,
        ref: &mut CCIPObjectRef,
        local_chain_selector: u64,
        ctx: &mut TxContext
    ) {
        assert!(
            !state_object::contains(ref, RMN_REMOTE_STATE_NAME),
            E_ALREADY_INITIALIZED
        );

        assert!(
            local_chain_selector != 0,
            E_ZERO_VALUE_NOT_ALLOWED
        );

        let state = RMNRemoteState {
            id: object::new(ctx),
            local_chain_selector,
            config: Config {
                rmn_home_contract_config_digest: vector[],
                signers: vector[],
                f_sign: 0
            },
            config_count: 0,
            signers: vec_map::empty<vector<u8>, bool>(),
            cursed_subjects: vec_map::empty<vector<u8>, bool>()
        };

        state_object::add(ownerCap, ref, RMN_REMOTE_STATE_NAME, state);
    }

    fun calculate_report(report: &Report): vector<u8> {
        let mut digest = vector[];
        eth_abi::encode_bytes32(&mut digest, get_report_digest_header());
        eth_abi::encode_u64(&mut digest, report.dest_chain_selector);
        eth_abi::encode_address(&mut digest, report.rmn_remote_contract_address);
        eth_abi::encode_address(&mut digest, report.off_ramp_address);
        eth_abi::encode_bytes32(&mut digest, report.rmn_home_contract_config_digest);
        vector::do_ref!(
            &report.merkle_roots,
            |merkle_root| {
                let merkle_root: &MerkleRoot = merkle_root;
                eth_abi::encode_u64(&mut digest, merkle_root.source_chain_selector);
                eth_abi::encode_u64(&mut digest, merkle_root.min_seq_nr);
                eth_abi::encode_u64(&mut digest, merkle_root.max_seq_nr);
                eth_abi::encode_bytes32(&mut digest, merkle_root.merkle_root);
            }
        );
        digest
    }

    public fun verify(
        ref: &CCIPObjectRef,
        merkle_root_source_chain_selectors: vector<u64>,
        merkle_root_min_seq_nrs: vector<u64>,
        merkle_root_max_seq_nrs: vector<u64>,
        merkle_root_values: vector<vector<u8>>,
        signatures: vector<vector<u8>>
    ): bool {
        let state = state_object::borrow<RMNRemoteState>(ref, RMN_REMOTE_STATE_NAME);

        assert!(state.config_count > 0, E_CONFIG_NOT_SET);

        let signatures_len = signatures.length();
        assert!(
            signatures_len >= (state.config.f_sign + 1),
            E_THRESHOLD_NOT_MET
        );

        let merkle_root_len = merkle_root_source_chain_selectors.length();
        assert!(
            merkle_root_len == merkle_root_min_seq_nrs.length(),
            E_MERKLE_ROOT_LENGTH_MISMATCH
        );
        assert!(
            merkle_root_len == merkle_root_max_seq_nrs.length(),
            E_MERKLE_ROOT_LENGTH_MISMATCH
        );
        assert!(
            merkle_root_len == merkle_root_values.length(),
            E_MERKLE_ROOT_LENGTH_MISMATCH
        );

        // Since we cannot pass public structs, we need to reconpublic struct it from the individual components.
        let mut merkle_roots = vector[];
        let mut i = 0;
        while (i < merkle_root_len) {
            let source_chain_selector =
                *vector::borrow(&merkle_root_source_chain_selectors, i);
            let min_seq_nr =
                *vector::borrow(&merkle_root_min_seq_nrs, i);
            let max_seq_nr =
                *vector::borrow(&merkle_root_max_seq_nrs, i);
            let merkle_root = *vector::borrow(&merkle_root_values, i);
            vector::push_back(
                &mut merkle_roots,
                MerkleRoot {
                    source_chain_selector,
                    min_seq_nr,
                    max_seq_nr,
                    merkle_root
                }
            );
            i = i + 1;
        };

        // there is no direct way to get chain id from Sui Move, removing dest_chain_id
        let report = Report {
            dest_chain_selector: state.local_chain_selector,
            rmn_remote_contract_address: @ccip,
            off_ramp_address: @ccip,
            rmn_home_contract_config_digest: state.config.rmn_home_contract_config_digest,
            merkle_roots
        };

        let digest = calculate_report(&report);

        let mut previous_eth_address = vector[];
        let mut i = 0;
        while (i < signatures_len) {
            let signature_bytes = *vector::borrow(&signatures, i);

            assert!(signature_bytes.length() == SIGNATURE_NUM_BYTES, E_INVALID_SIGNATURE);

            // rmn only generates signatures with v = 27, subtract the ethereum recover id offset of 27 to get zero.
            // according to Sui Move document: https://docs.sui.io/references/framework/sui/ecdsa_k1#sui_ecdsa_k1_secp256k1_ecrecover
            // the digest is not hashed. hence the digest in calculate_report is not keccack256 hashed
            let public_key_bytes = ecdsa_k1::secp256k1_ecrecover(&signature_bytes, &digest, 0);

            // trim the first 12 bytes of the hash to recover the ethereum address.
            // there is no vector::trim function available
            let mut eth_address = vector::empty();
            let key_hash = &hash::keccak256(&public_key_bytes);
            let len = 32;
            let mut j: u64 = 12;
            while (j < len) {
                // Copy each element starting at index 12 into the new vector.
                vector::push_back(&mut eth_address, *vector::borrow(key_hash, j));
                j = j + 1;
            };

            assert!(
                vec_map::contains(&state.signers, &eth_address),
                E_UNEXPECTED_SIGNER
            );
            assert!(
                vec_map::contains(&state.signers, &eth_address),
                E_UNEXPECTED_SIGNER
            );
            if (i > 0) {
                assert!(
                    merkle_proof::vector_u8_gt(&eth_address, &previous_eth_address),
                    E_OUT_OF_ORDER_SIGNATURES
                );
            };
            previous_eth_address = eth_address;

            i = i + 1;
        };

        true
    }

    public fun set_config(
        ownerCap: &OwnerCap,
        ref: &mut CCIPObjectRef,
        rmn_home_contract_config_digest: vector<u8>,
        signer_onchain_public_keys: vector<vector<u8>>,
        node_indexes: vector<u64>,
        f_sign: u64,
        _ctx: &mut TxContext
    ) {
        let state = state_object::borrow_mut<RMNRemoteState>(ownerCap, ref, RMN_REMOTE_STATE_NAME);

        assert!(
            rmn_home_contract_config_digest.length() == 32,
            E_INVALID_DIGEST_LENGTH
        );

        assert!(
            eth_abi::decode_u256_value(rmn_home_contract_config_digest) != 0,
            E_ZERO_VALUE_NOT_ALLOWED
        );

        let signers_len = signer_onchain_public_keys.length();
        assert!(
            signers_len == node_indexes.length(),
            E_SIGNERS_MISMATCH
        );

        let mut i = 1;
        while (i < signers_len) {
            let previous_node_index = *vector::borrow(&node_indexes, i - 1);
            let current_node_index = *vector::borrow(&node_indexes, i);
            assert!(
                previous_node_index < current_node_index,
                E_INVALID_SIGNER_ORDER
            );
            i = i + 1;
        };

        assert!(
            signers_len >= (2 * f_sign + 1),
            E_NOT_ENOUGH_SIGNERS
        );

        // smart_table::clear(&mut state.signers);
        let keys = vec_map::keys(&state.signers);
        let mut i = 0;
        let keys_len = keys.length();
        while (i < keys_len) {
            let key = *vector::borrow(&keys, i);
            vec_map::remove(&mut state.signers, &key);
            i = i + 1;
        };

        let signers = vector::zip_map_ref!(
            &signer_onchain_public_keys,
            &node_indexes,
            |signer_public_key_bytes, node_indexes| {
                let signer_public_key_bytes: vector<u8> = *signer_public_key_bytes;
                let node_index: u64 = *node_indexes;
                // expect an ethereum address of 20 bytes.
                assert!(
                    signer_public_key_bytes.length() == 20,
                    E_INVALID_PUBLIC_KEY_LENGTH
                );
                assert!(
                    !vec_map::contains(&state.signers, &signer_public_key_bytes),
                    E_DUPLICATE_SIGNER
                );
                vec_map::insert(&mut state.signers, signer_public_key_bytes, true);
                Signer {
                    onchain_public_key: signer_public_key_bytes,
                    node_index
                }
            }
        );

        let new_config = Config {
            rmn_home_contract_config_digest,
            signers,
            f_sign
        };
        state.config = new_config;

        let new_config_count = state.config_count + 1;
        state.config_count = new_config_count;

        event::emit(ConfigSet { version: new_config_count, config: new_config });
    }

    public fun get_versioned_config(ref: &CCIPObjectRef): (u32, Config) {
        let state = state_object::borrow<RMNRemoteState>(ref, RMN_REMOTE_STATE_NAME);

        (state.config_count, state.config)
    }

    public fun get_local_chain_selector(ref: &CCIPObjectRef): u64 {
        let state = state_object::borrow<RMNRemoteState>(ref, RMN_REMOTE_STATE_NAME);

        state.local_chain_selector
    }

    public fun get_report_digest_header(): vector<u8> {
        hash::keccak256(&b"RMN_V1_6_ANY2SUI_REPORT")
    }

    public fun curse(ownerCap: &OwnerCap, ref: &mut CCIPObjectRef, subject: vector<u8>, ctx: &mut TxContext) {
        curse_multiple(ownerCap, ref, vector[subject], ctx);
    }

    public fun curse_multiple(
        ownerCap: &OwnerCap, ref: &mut CCIPObjectRef, subjects: vector<vector<u8>>, _ctx: &mut TxContext
    ) {
        let state = state_object::borrow_mut<RMNRemoteState>(ownerCap, ref, RMN_REMOTE_STATE_NAME);

        vector::do_ref!(
            &subjects,
            |subject| {
                let subject: vector<u8> = *subject;
                assert!(
                    subject.length() == 16,
                    E_INVALID_SUBJECT_LENGTH
                );
                assert!(
                    !vec_map::contains(&state.cursed_subjects, &subject),
                    E_ALREADY_CURSED
                );
                vec_map::insert(&mut state.cursed_subjects, subject, true);
            }
        );
        event::emit(Cursed { subjects });
    }

    public fun uncurse(
        ownerCap: &OwnerCap,
        ref: &mut CCIPObjectRef,
        subject: vector<u8>,
        ctx: &mut TxContext
    ) {
        uncurse_multiple(ownerCap, ref, vector[subject], ctx);
    }

    public fun uncurse_multiple(
        ownerCap: &OwnerCap,
        ref: &mut CCIPObjectRef,
        subjects: vector<vector<u8>>,
        _ctx: &mut TxContext
    ) {
        let state = state_object::borrow_mut<RMNRemoteState>(ownerCap, ref, RMN_REMOTE_STATE_NAME);

        vector::do_ref!(
            &subjects,
            |subject| {
                let subject: vector<u8> = *subject;
                assert!(
                    vec_map::contains(&state.cursed_subjects, &subject),
                    E_NOT_CURSED
                );
                vec_map::remove(&mut state.cursed_subjects, &subject);
            }
        );
        event::emit(Uncursed { subjects });
    }

    public fun get_cursed_subjects(ref: &CCIPObjectRef): vector<vector<u8>> {
        let state = state_object::borrow<RMNRemoteState>(ref, RMN_REMOTE_STATE_NAME);

        vec_map::keys(&state.cursed_subjects)
    }

    #[allow(implicit_const_copy)]
    public fun is_cursed_global(ref: &CCIPObjectRef): bool {
        let state = state_object::borrow<RMNRemoteState>(ref, RMN_REMOTE_STATE_NAME);

        vec_map::contains(&state.cursed_subjects, &GLOBAL_CURSE_SUBJECT)
    }

    public fun is_cursed(ref: &CCIPObjectRef, subject: vector<u8>): bool {
        let state = state_object::borrow<RMNRemoteState>(ref, RMN_REMOTE_STATE_NAME);

        vec_map::contains(&state.cursed_subjects, &subject) || is_cursed_global(ref)
    }

    public fun is_cursed_u128(ref: &CCIPObjectRef, subject_value: u128): bool {
        let mut subject = bcs::to_bytes(&subject_value);
        vector::reverse(&mut subject);
        is_cursed(ref, subject)
    }

    #[test_only]
    public fun get_config(config: &Config): (vector<u8>, vector<Signer>, u64) {
        (config.rmn_home_contract_config_digest, config.signers, config.f_sign)
    }

    //
    // MCMS entrypoint
    //

    // public struct McmsCallback has drop {}
    //
    // public fun mcms_entrypoint<T: key>(
    //     _metadata: object::Object<T>
    // ): option::Option<u128> acquires RMNRemoteState {
    //     let (caller, function, data) =
    //         mcms_registry::get_callback_params(@ccip, McmsCallback {});
    //
    //     let function_bytes = *string::bytes(&function);
    //     let stream = bcs_stream::new(data);
    //
    //     if (function_bytes == b"initialize") {
    //         let local_chain_selector = bcs_stream::deserialize_u64(&mut stream);
    //         bcs_stream::assert_is_consumed(&stream);
    //         initialize(&caller, local_chain_selector);
    //     } else if (function_bytes == b"set_config") {
    //         let rmn_home_contract_config_digest =
    //             bcs_stream::deserialize_vector_u8(&mut stream);
    //         let signer_onchain_public_keys =
    //             bcs_stream::deserialize_vector(
    //                 &mut stream,
    //                 |stream| bcs_stream::deserialize_vector_u8(stream)
    //             );
    //         let node_indexes =
    //             bcs_stream::deserialize_vector(
    //                 &mut stream,
    //                 |stream| bcs_stream::deserialize_u64(stream)
    //             );
    //         let f_sign = bcs_stream::deserialize_u64(&mut stream);
    //         bcs_stream::assert_is_consumed(&stream);
    //         set_config(
    //             &caller,
    //             rmn_home_contract_config_digest,
    //             signer_onchain_public_keys,
    //             node_indexes,
    //             f_sign
    //         )
    //     } else if (function_bytes == b"curse") {
    //         let subject = bcs_stream::deserialize_vector_u8(&mut stream);
    //         bcs_stream::assert_is_consumed(&stream);
    //         curse(&caller, subject)
    //     } else if (function_bytes == b"curse_multiple") {
    //         let subjects =
    //             bcs_stream::deserialize_vector(
    //                 &mut stream,
    //                 |stream| bcs_stream::deserialize_vector_u8(stream)
    //             );
    //         bcs_stream::assert_is_consumed(&stream);
    //         curse_multiple(&caller, subjects)
    //     } else if (function_bytes == b"uncurse") {
    //         let subject = bcs_stream::deserialize_vector_u8(&mut stream);
    //         bcs_stream::assert_is_consumed(&stream);
    //         uncurse(&caller, subject)
    //     } else if (function_bytes == b"uncurse_multiple") {
    //         let subjects =
    //             bcs_stream::deserialize_vector(
    //                 &mut stream,
    //                 |stream| bcs_stream::deserialize_vector_u8(stream)
    //             );
    //         bcs_stream::assert_is_consumed(&stream);
    //         uncurse_multiple(&caller, subjects)
    //     } else {
    //         abort E_UNKNOWN_FUNCTION)
    //     };
    //
    //     option::none()
    // }
}

#[test_only]
module ccip::rmn_remote_test {
    use ccip::state_object::{Self, OwnerCap, CCIPObjectRef};
    use ccip::rmn_remote;
    use sui::test_scenario::{Self, Scenario};

    const RMN_REMOTE_STATE_NAME: vector<u8> = b"RMNRemoteState";

    fun set_up_test(): (Scenario, OwnerCap, CCIPObjectRef) {
        let mut scenario = test_scenario::begin(@0x1);
        let ctx = scenario.ctx();

        let (owner_cap, ref) = state_object::create(ctx);
        (scenario, owner_cap, ref)
    }

    fun tear_down_test(scenario: Scenario, owner_cap: OwnerCap, ref: CCIPObjectRef) {
        state_object::destroy_owner_cap(owner_cap);
        state_object::destroy_state_object(ref);
        test_scenario::end(scenario);
    }

    #[test]
    public fun test_initialize() {
        let (mut scenario, owner_cap, mut ref) = set_up_test();
        let ctx = scenario.ctx();

        rmn_remote::initialize(&owner_cap, &mut ref, 1, ctx);
        let _state = state_object::borrow<rmn_remote::RMNRemoteState>(&ref, RMN_REMOTE_STATE_NAME);
        assert!(rmn_remote::get_local_chain_selector(&ref) == 1);

        tear_down_test(scenario, owner_cap, ref);
    }

    #[test]
    #[expected_failure(abort_code = rmn_remote::E_ZERO_VALUE_NOT_ALLOWED)]
    public fun test_initialize_zero_chain_selector() {
        let (mut scenario, owner_cap, mut ref) = set_up_test();
        let ctx = scenario.ctx();

        rmn_remote::initialize(&owner_cap, &mut ref, 0, ctx);

        tear_down_test(scenario, owner_cap, ref);
    }

    #[test]
    #[expected_failure(abort_code = rmn_remote::E_ALREADY_INITIALIZED)]
    public fun test_initialize_already_initialized() {
        let (mut scenario, owner_cap, mut ref) = set_up_test();
        let ctx = scenario.ctx();

        rmn_remote::initialize(&owner_cap, &mut ref, 1, ctx);
        rmn_remote::initialize(&owner_cap, &mut ref, 1, ctx);

        tear_down_test(scenario, owner_cap, ref);
    }

    #[test]
    public fun test_set_config() {
        let (mut scenario, owner_cap, mut ref) = set_up_test();
        let ctx = scenario.ctx();

        rmn_remote::initialize(&owner_cap, &mut ref, 1, ctx);
        rmn_remote::set_config(
            &owner_cap,
            &mut ref,
            b"00000000000000000000000000000001",
            vector[
                b"00000000000000000002",
                b"00000000000000000003",
                b"00000000000000000004"
            ],
            vector[0, 1, 2],
            1,
            ctx
        );

        let (version, config) = rmn_remote::get_versioned_config(&ref);

        assert!(version == 1);

        let (digest, signers, f_sign) = rmn_remote::get_config(&config);
        assert!(digest == b"00000000000000000000000000000001");
        assert!(signers.length() == 3);
        assert!(f_sign == 1);

        tear_down_test(scenario, owner_cap, ref);
    }

    #[test]
    #[expected_failure(abort_code = rmn_remote::E_INVALID_DIGEST_LENGTH)]
    public fun test_set_config_invalid_digest_length() {
        let( mut scenario, owner_cap, mut ref) = set_up_test();
        let ctx = scenario.ctx();

        rmn_remote::initialize(&owner_cap, &mut ref, 1, ctx);
        rmn_remote::set_config(
            &owner_cap,
            &mut ref,
            b"000000000000000000000000000000", // invalid digest length
            vector[
                b"00000000000000000002",
                b"00000000000000000003",
                b"00000000000000000004"
            ],
            vector[0, 1, 2],
            1,
            ctx
        );

        tear_down_test(scenario, owner_cap, ref);
    }

    #[test]
    #[expected_failure(abort_code = rmn_remote::E_ZERO_VALUE_NOT_ALLOWED)]
    public fun test_set_config_zero_digest() {
        let( mut scenario, owner_cap, mut ref) = set_up_test();
        let ctx = scenario.ctx();

        rmn_remote::initialize(&owner_cap, &mut ref, 1, ctx);
        rmn_remote::set_config(
            &owner_cap,
            &mut ref,
            x"0000000000000000000000000000000000000000000000000000000000000000", // zero digest
            vector[
                b"00000000000000000002",
                b"00000000000000000003",
                b"00000000000000000004"
            ],
            vector[0, 1, 2],
            1,
            ctx
        );

        tear_down_test(scenario, owner_cap, ref);
    }

    #[test]
    #[expected_failure(abort_code = rmn_remote::E_NOT_ENOUGH_SIGNERS)]
    public fun test_set_config_not_enough_signers() {
        let (mut scenario, owner_cap, mut ref) = set_up_test();
        let ctx = scenario.ctx();

        rmn_remote::initialize(&owner_cap, &mut ref, 1, ctx);
        rmn_remote::set_config(
            &owner_cap,
            &mut ref,
            b"00000000000000000000000000000001",
            vector[
                b"00000000000000000002",
                b"00000000000000000003",
                b"00000000000000000004"
            ],
            vector[0, 1, 2],
            2, // f_sign is 2, but only 3 signers
            ctx
        );

        tear_down_test(scenario, owner_cap, ref);
    }

    #[test]
    #[expected_failure(abort_code = rmn_remote::E_SIGNERS_MISMATCH)]
    public fun test_set_config_signers_mismatch() {
        let (mut scenario, owner_cap, mut ref) = set_up_test();
        let ctx = scenario.ctx();

        rmn_remote::initialize(&owner_cap, &mut ref, 1, ctx);
        rmn_remote::set_config(
            &owner_cap,
            &mut ref,
            b"00000000000000000000000000000001",
            vector[
                b"00000000000000000002",
                b"00000000000000000003"
            ],
            vector[0, 1, 2], // 3 signers, but 2 pub keys
            1,
            ctx
        );

        tear_down_test(scenario, owner_cap, ref);
    }

    #[test]
    #[expected_failure(abort_code = rmn_remote::E_INVALID_SIGNER_ORDER)]
    public fun test_set_config_invalid_signer_order() {
        let (mut scenario, owner_cap, mut ref) = set_up_test();
        let ctx = scenario.ctx();

        rmn_remote::initialize(&owner_cap, &mut ref, 1, ctx);
        rmn_remote::set_config(
            &owner_cap,
            &mut ref,
            b"00000000000000000000000000000001",
            vector[
                b"00000000000000000002",
                b"00000000000000000003",
                b"00000000000000000004"
            ],
            vector[1, 0, 2], // invalid order
            1,
            ctx
        );

        tear_down_test(scenario, owner_cap, ref);
    }

    #[test]
    public fun test_curse() {
        let (mut scenario, owner_cap, mut ref) = set_up_test();
        let ctx = scenario.ctx();

        rmn_remote::initialize(&owner_cap, &mut ref, 1, ctx);
        rmn_remote::curse(&owner_cap, &mut ref, b"0000000000000003", ctx);

        let cursed_subjects = rmn_remote::get_cursed_subjects(&ref);
        assert!(cursed_subjects.length() == 1);

        assert!(rmn_remote::is_cursed(&ref, b"0000000000000003"));

        tear_down_test(scenario, owner_cap, ref);
    }

    #[test]
    #[expected_failure(abort_code = rmn_remote::E_INVALID_SUBJECT_LENGTH)]
    public fun test_curse_invalid_subject_length() {
        let (mut scenario, owner_cap, mut ref) = set_up_test();
        let ctx = scenario.ctx();

        rmn_remote::initialize(&owner_cap, &mut ref, 1, ctx);
        rmn_remote::curse(&owner_cap, &mut ref, b"00003", ctx);

        tear_down_test(scenario, owner_cap, ref);
    }

    #[test]
    #[expected_failure(abort_code = rmn_remote::E_ALREADY_CURSED)]
    public fun test_curse_already_cursed() {
        let (mut scenario, owner_cap, mut ref) = set_up_test();
        let ctx = scenario.ctx();

        rmn_remote::initialize(&owner_cap, &mut ref, 1, ctx);
        rmn_remote::curse(&owner_cap, &mut ref, b"0000000000000003", ctx);
        rmn_remote::curse(&owner_cap, &mut ref, b"0000000000000003", ctx);

        tear_down_test(scenario, owner_cap, ref);
    }

    #[test]
    public fun test_curse_multiple() {
        let (mut scenario, owner_cap, mut ref) = set_up_test();
        let ctx = scenario.ctx();

        rmn_remote::initialize(&owner_cap, &mut ref, 1, ctx);
        rmn_remote::curse_multiple(
            &owner_cap,
            &mut ref,
            vector[
                b"0000000000000003",
                b"0000000000000004",
            ],
            ctx
        );

        let cursed_subjects = rmn_remote::get_cursed_subjects(&ref);
        assert!(vector::length(&cursed_subjects) == 2);

        assert!(rmn_remote::is_cursed(&ref, b"0000000000000003"));
        assert!(rmn_remote::is_cursed(&ref, b"0000000000000004"));

        tear_down_test(scenario, owner_cap, ref);
    }

    #[test]
    public fun test_uncurse() {
        let (mut scenario, owner_cap, mut ref) = set_up_test();
        let ctx = scenario.ctx();

        rmn_remote::initialize(&owner_cap, &mut ref, 1, ctx);
        rmn_remote::curse(&owner_cap, &mut ref, b"0000000000000003", ctx);
        let mut cursed_subjects = rmn_remote::get_cursed_subjects(&ref);
        assert!(vector::length(&cursed_subjects) == 1);
        assert!(rmn_remote::is_cursed(&ref, b"0000000000000003"));

        rmn_remote::uncurse(&owner_cap, &mut ref, b"0000000000000003", ctx);
        cursed_subjects = rmn_remote::get_cursed_subjects(&ref);
        assert!(cursed_subjects.length() == 0);
        assert!(!rmn_remote::is_cursed(&ref, b"0000000000000003"));

        tear_down_test(scenario, owner_cap, ref);
    }

    #[test]
    public fun test_is_cursed_global() {
        let (mut scenario, owner_cap, mut ref) = set_up_test();
        let ctx = scenario.ctx();

        rmn_remote::initialize(&owner_cap, &mut ref, 1, ctx);
        rmn_remote::curse(&owner_cap, &mut ref, x"01000000000000000000000000000001", ctx);

        let cursed_subjects = rmn_remote::get_cursed_subjects(&ref);
        assert!(cursed_subjects.length() == 1);
        assert!(rmn_remote::is_cursed_global(&ref));

        tear_down_test(scenario, owner_cap, ref);
    }

    #[test]
    public fun test_is_cursed_u128() {
        let (mut scenario, owner_cap, mut ref) = set_up_test();
        let ctx = scenario.ctx();

        rmn_remote::initialize(&owner_cap, &mut ref, 1, ctx);
        rmn_remote::curse(&owner_cap, &mut ref, x"00000000000000000000000000000100", ctx); // hex(256)

        assert!(rmn_remote::is_cursed_u128(&ref, 256));
        assert!(!rmn_remote::is_cursed_u128(&ref, 100));

        tear_down_test(scenario, owner_cap, ref);
    }

    // TODO: add tests for verify
}
