module ccip::ocr3_base {
    use std::bit_vector;
    use sui::ed25519;
    use sui::table;
    use sui::hash;
    use sui::event;

    const MAX_NUM_ORACLES: u64 = 256;
    const OCR_PLUGIN_TYPE_COMMIT: u8 = 1;
    const OCR_PLUGIN_TYPE_EXECUTION: u8 = 2;
    const PUBLIC_KEY_NUM_BYTES: u64 = 32;

    const E_BIG_F_MUST_BE_POSITIVE: u64 = 1;
    const E_STATIC_CONFIG_CANNOT_BE_CHANGED: u64 = 2;
    const E_TOO_MANY_SIGNERS: u64 = 3;
    const E_BIG_F_TOO_HIGH: u64 = 4;
    const E_TOO_MANY_TRANSMITTERS: u64 = 5;
    const E_NO_TRANSMITTERS: u64 = 6;
    const E_REPEATED_SIGNERS: u64 = 7;
    const E_REPEATED_TRANSMITTERS: u64 = 8;
    const E_CONFIG_DIGEST_MISMATCH: u64 = 9;
    const E_UNAUTHORIZED_TRANSMITTER: u64 = 10;
    const E_WRONG_NUMBER_OF_SIGNATURES: u64 = 11;
    const E_COULD_NOT_VALIDATE_SIGNER_KEY: u64 = 12;
    const E_INVALID_REPORT_CONTEXT_LENGTH: u64 = 13;
    const E_INVALID_CONFIG_DIGEST_LENGTH: u64 = 14;
    const E_INVALID_SEQUENCE_LENGTH: u64 = 15;
    const E_UNAUTHORIZED_SIGNER: u64 = 16;
    const E_NON_UNIQUE_SIGNATURES: u64 = 17;
    const E_INVALID_SIGNATURE: u64 = 18;
    const E_OUT_OF_BYTES: u64 = 19;
    const E_WRONG_PUBKEY_SIZE: u64 = 20;

    public struct UnvalidatedPublicKey has copy, drop, store {
        bytes: vector<u8>
    }

    // TODO: not used?
    // public struct Oracle has store, drop {
    //     index: u8,
    //     role: u8
    // }

    public struct ConfigInfo has store, drop {
        config_digest: vector<u8>,
        big_f: u8,
        n: u8,
        is_signature_verification_enabled: bool
    }

    public struct Transmitted has copy, drop {
        ocr_plugin_type: u8,
        config_digest: vector<u8>,
        sequence_number: u64
    }

    public struct ConfigSet has copy, drop {
        ocr_plugin_type: u8,
        config_digest: vector<u8>,
        signers: vector<vector<u8>>,
        transmitters: vector<address>,
        big_f: u8
    }

    public struct OCRConfig has store, drop {
        config_info: ConfigInfo,
        signers: vector<vector<u8>>,
        transmitters: vector<address>
    }

    public struct OCR3BaseState has store {
        chain_id: u8,
        // ocr plugin type -> ocr config
        ocr3_configs: table::Table<u8, OCRConfig>,
        // ocr plugin type -> signers
        signer_oracles: table::Table<u8, vector<UnvalidatedPublicKey>>,
        // ocr plugin type -> transmitters
        transmitter_oracles: table::Table<u8, vector<address>>
    }

    public fun ocr_plugin_type_commit(): u8 {
        OCR_PLUGIN_TYPE_COMMIT
    }

    public fun ocr_plugin_type_execution(): u8 {
        OCR_PLUGIN_TYPE_EXECUTION
    }

    // TODO: verify if we can get chain id from Sui Move
    public fun new(ctx: &mut TxContext): OCR3BaseState {
        // assert!(chain_id != chain_ids::sui_mainnet() && chain_id != chain_ids::sui_testnet() && chain_id != chain_ids::sui_custom(), E_WRONG_CHAIN_ID);
        OCR3BaseState {
            // chain_id,
            chain_id: 0,
            ocr3_configs: table::new<u8, OCRConfig>(ctx),
            signer_oracles: table::new<u8, vector<UnvalidatedPublicKey>>(ctx),
            transmitter_oracles: table::new<u8, vector<address>>(ctx)
        }
    }

    // this function relies on the ability to get chain id from Sui Move
    // public fun assert_chain_not_forked(ocr3_state: &OCR3BaseState) {
    //     assert!(
    //         chain_id::get() == ocr3_state.chain_id,
    //         E_FORKED_CHAIN
    //     );
    // }

    public fun latest_config_details(
        ocr3_state: &OCR3BaseState, ocr_plugin_type: u8
    ): (vector<u8>, u8, u8, bool, vector<vector<u8>>, vector<address>) {
        let ocr_config = table::borrow(&ocr3_state.ocr3_configs, ocr_plugin_type);
        let config_info = &ocr_config.config_info;
        (
            config_info.config_digest,
            config_info.big_f,
            config_info.n,
            config_info.is_signature_verification_enabled,
            ocr_config.signers,
            ocr_config.transmitters
        )
    }

    // equivalent of uint64(uint256(reportContext[1]))
    public fun deserialize_sequence_bytes(
        sequence_bytes: vector<u8>
    ): u64 {
        let len = vector::length(&sequence_bytes);
        let mut result: u64 = 0;
        let mut i = len - 8;
        while (i < len) {
            result = (result << 8) + (*vector::borrow(&sequence_bytes, i) as u64);
            i = i + 1;
        };
        result
    }

    // equivalent of keccak256(abi.encodePacked(keccak256(report), reportContext))
    fun hash_report(
        report: vector<u8>, config_digest: vector<u8>, sequence_bytes: vector<u8>
    ): vector<u8> {
        let mut bytes = vector[];

        vector::append(&mut bytes, hash::keccak256(&report));
        vector::append(&mut bytes, config_digest);
        vector::append(&mut bytes, sequence_bytes);

        hash::keccak256(&bytes)
    }

    fun has_duplicates<T>(a: &vector<T>): bool {
        let len = vector::length(a);
        let mut i = 0;

        while (i < len) {
            let mut j = i + 1;
            while (j < len) {
                if (vector::borrow(a, i) == vector::borrow(a, j)) {
                    return true
                };
                j = j + 1;
            };
            i = i + 1;
        };
        false
    }

    fun assign_transmitter_oracles(
        transmitter_oracles: &mut table::Table<u8, vector<address>>,
        ocr_plugin_type: u8,
        transmitters: &vector<address>
    ) {
        assert!(
            !has_duplicates(transmitters),
            E_REPEATED_TRANSMITTERS
        );

        if (table::contains(transmitter_oracles, ocr_plugin_type)) {
            let _old_value = table::remove(transmitter_oracles, ocr_plugin_type);
        };
        table::add(transmitter_oracles, ocr_plugin_type, *transmitters);
    }

    // TODO: explore more valid public key checks
    fun assign_signer_oracles(
        signer_oracles: &mut table::Table<u8, vector<UnvalidatedPublicKey>>,
        ocr_plugin_type: u8,
        signers: &vector<vector<u8>>
    ) {
        assert!(!has_duplicates(signers), E_REPEATED_SIGNERS);

        let validated_signers = vector::map_ref!(
            signers,
            |signer| {
                assert!(
                    validate_public_key(signer),
                    E_COULD_NOT_VALIDATE_SIGNER_KEY
                );
                UnvalidatedPublicKey { bytes: *signer }
            }
        );

        if (table::contains(signer_oracles, ocr_plugin_type)) {
            let _old_value = table::remove(signer_oracles, ocr_plugin_type);
        };
        table::add(signer_oracles, ocr_plugin_type, validated_signers);
    }

    // Sui Move ed25519 does not
    // TODO: verify if we can provide more validation for public key
    fun validate_public_key(pubkey: &vector<u8>): bool {
        vector::length(pubkey) == 32
    }

    fun verify_signature(
        signers: &vector<UnvalidatedPublicKey>,
        hashed_report: vector<u8>,
        signatures: vector<vector<u8>>
    ) {
        let mut seen = bit_vector::new(vector::length(signers));
        vector::do_ref!(
            &signatures,
            |signature_bytes| {
                let public_key =
                    new_unvalidated_public_key_from_bytes(slice(signature_bytes, 0, 32));
                let (exists, index) = vector::index_of(signers, &public_key);
                assert!(exists, E_UNAUTHORIZED_SIGNER);
                assert!(
                    !bit_vector::is_index_set(&seen, index),
                    E_NON_UNIQUE_SIGNATURES
                );
                bit_vector::set(&mut seen, index);
                let signature = slice(signature_bytes, 32, 64);

                let verified =
                ed25519::ed25519_verify(
                    &signature, &public_key.bytes, &hashed_report
                );
                assert!(verified, E_INVALID_SIGNATURE);
            }
        );
    }

    fun new_unvalidated_public_key_from_bytes(bytes: vector<u8>): UnvalidatedPublicKey {
        assert!(std::vector::length(&bytes) == PUBLIC_KEY_NUM_BYTES, E_WRONG_PUBKEY_SIZE);
        UnvalidatedPublicKey { bytes }
    }

    /// Returns a new vector containing `len` elements from `vec`
    /// starting at index `start`. Panics if `start + len` exceeds the vector length.
    fun slice<T: copy>(vec: &vector<T>, start: u64, len: u64): vector<T> {
        let vec_len = vector::length(vec);
        // Ensure we have enough elements for the slice.
        assert!(start + len <= vec_len, E_OUT_OF_BYTES);
        let mut new_vec = vector::empty<T>();
        let mut i = start;
        while (i < start + len) {
            // Copy each element from the original vector into the new vector.
            vector::push_back(&mut new_vec, *vector::borrow(vec, i));
            i = i + 1;
        };
        new_vec
    }

    public fun transmit(
        ocr3_state: &mut OCR3BaseState,
        transmitter: address,
        ocr_plugin_type: u8,
        report_context: vector<vector<u8>>,
        report: vector<u8>,
        signatures: vector<vector<u8>>
    ) {
        let ocr_config = table::borrow(&ocr3_state.ocr3_configs, ocr_plugin_type);
        let config_info = &ocr_config.config_info;

        assert!(
            vector::length(&report_context) == 2,
            E_INVALID_REPORT_CONTEXT_LENGTH
        );

        let config_digest = *vector::borrow(&report_context, 0);
        assert!(
            vector::length(&config_digest) == 32,
            E_INVALID_CONFIG_DIGEST_LENGTH
        );

        let sequence_bytes = *vector::borrow(&report_context, 1);
        assert!(
            vector::length(&sequence_bytes) == 32,
            E_INVALID_SEQUENCE_LENGTH
        );

        // TODO: EVM checks transaction data length here

        assert!(
            config_digest == config_info.config_digest,
            E_CONFIG_DIGEST_MISMATCH
        );

        // TODO: verify if this is feasible
        // assert_chain_not_forked(ocr3_state);

        let plugin_transmitters =
            table::borrow(&ocr3_state.transmitter_oracles, ocr_plugin_type);
        assert!(
            vector::contains(plugin_transmitters, &transmitter),
            E_UNAUTHORIZED_TRANSMITTER
        );

        if (config_info.is_signature_verification_enabled) {
            assert!(
                vector::length(&signatures) == (config_info.big_f as u64) + 1,
                E_WRONG_NUMBER_OF_SIGNATURES
            );

            let hashed_report = hash_report(report, config_digest, sequence_bytes);
            let plugin_signers =
                table::borrow(&ocr3_state.signer_oracles, ocr_plugin_type);
            verify_signature(plugin_signers, hashed_report, signatures);
        };

        let sequence_number: u64 = deserialize_sequence_bytes(sequence_bytes);
        event::emit(Transmitted { ocr_plugin_type, config_digest, sequence_number });
    }

    // used by offramp, hence not marked as entry
    public fun set_ocr3_config(
        // caller: address,
        ocr3_state: &mut OCR3BaseState,
        config_digest: vector<u8>,
        ocr_plugin_type: u8,
        big_f: u8,
        is_signature_verification_enabled: bool,
        signers: vector<vector<u8>>,
        transmitters: vector<address>,
        _ctx: &mut TxContext
    ) {
        // TODO: implement this check in auth
        // auth::assert_only_owner(tx_context::sender(ctx));
        assert!(big_f != 0, E_BIG_F_MUST_BE_POSITIVE);

        let ocr_config = if (table::contains(&ocr3_state.ocr3_configs, ocr_plugin_type)) {
            table::borrow_mut(&mut ocr3_state.ocr3_configs, ocr_plugin_type)
        } else {
            table::add(
                &mut ocr3_state.ocr3_configs,
                ocr_plugin_type,
                OCRConfig {
                    config_info: ConfigInfo {
                        config_digest: vector[],
                        big_f: 0,
                        n: 0,
                        is_signature_verification_enabled: false
                    },
                    signers: vector[],
                    transmitters: vector[]
                }
            );
            table::borrow_mut(&mut ocr3_state.ocr3_configs, ocr_plugin_type)
        };

        let config_info = &mut ocr_config.config_info;

        // If F is 0, then the config is not yet set.
        if (config_info.big_f == 0) {
            config_info.is_signature_verification_enabled = is_signature_verification_enabled;
        } else {
            assert!(
                config_info.is_signature_verification_enabled == is_signature_verification_enabled,
                E_STATIC_CONFIG_CANNOT_BE_CHANGED
            );
        };

        assert!(
            vector::length(&transmitters) <= MAX_NUM_ORACLES,
            E_TOO_MANY_TRANSMITTERS
        );
        assert!(
            vector::length(&transmitters) > 0,
            E_NO_TRANSMITTERS
        );

        if (is_signature_verification_enabled) {
            assert!(
                vector::length(&signers) <= MAX_NUM_ORACLES,
                E_TOO_MANY_SIGNERS
            );
            assert!(
                vector::length(&signers) > 3 * (big_f as u64),
                E_BIG_F_TOO_HIGH
            );
            // NOTE: Transmitters cannot exceed signers. Transmitters do not have to be >= 3F + 1 because they can
            // match >= 3fChain + 1, where fChain <= F. fChain is not represented in MultiOCR3Base - so we skip this check.
            assert!(
                vector::length(&signers) >= vector::length(&transmitters),
                E_TOO_MANY_TRANSMITTERS
            );

            config_info.n = vector::length(&signers) as u8;

            ocr_config.signers = signers;

            assign_signer_oracles(
                &mut ocr3_state.signer_oracles, ocr_plugin_type, &signers
            );
        };

        ocr_config.transmitters = transmitters;

        assign_transmitter_oracles(
            &mut ocr3_state.transmitter_oracles, ocr_plugin_type, &transmitters
        );

        config_info.big_f = big_f;
        config_info.config_digest = config_digest;

        event::emit(
            ConfigSet { ocr_plugin_type, config_digest, signers, transmitters, big_f }
        );
    }
}
