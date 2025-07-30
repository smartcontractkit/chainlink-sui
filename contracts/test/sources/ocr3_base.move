// ================================================================
//          THIS IS A TEST CONTRACT FOR THE OCR3 BASE
// ================================================================

module test::ocr3_base;

use std::bit_vector;

use sui::ed25519;
use sui::event;
use sui::hash;
use sui::table::{Self, Table};

use sui::address;

const MAX_NUM_ORACLES: u64 = 256;
const OCR_PLUGIN_TYPE_COMMIT: u8 = 0;
const OCR_PLUGIN_TYPE_EXECUTION: u8 = 1;
const PUBLIC_KEY_NUM_BYTES: u64 = 32;

public struct UnvalidatedPublicKey has copy, drop, store {
    bytes: vector<u8>
}

public struct ConfigInfo has store, drop, copy {
    config_digest: vector<u8>,
    big_f: u8,
    n: u8,
    is_signature_verification_enabled: bool
}

public struct OCRConfig has store, drop, copy {
    config_info: ConfigInfo,
    signers: vector<vector<u8>>,
    transmitters: vector<address>
}

// this struct is stored in offramp state object
public struct OCR3BaseState has key, store {
    id: UID,
    // ocr plugin type -> ocr config
    ocr3_configs: Table<u8, OCRConfig>,
    // ocr plugin type -> signers
    signer_oracles: Table<u8, vector<UnvalidatedPublicKey>>,
    // ocr plugin type -> transmitters
    transmitter_oracles: Table<u8, vector<address>>
}

public struct ConfigSet has copy, drop {
    ocr_plugin_type: u8,
    config_digest: vector<u8>,
    signers: vector<vector<u8>>,
    transmitters: vector<address>,
    big_f: u8
}

public struct Transmitted has copy, drop {
    ocr_plugin_type: u8,
    config_digest: vector<u8>,
    sequence_number: u64
}

const EWrongPubkeySize: u64 = 1;
const EBigFMustBePositive: u64 = 2;
const EStaticConfigCannotBeChanged: u64 = 3;
const ETooManySigners: u64 = 4;
const EBigFTooHigh: u64 = 5;
const ETooManyTransmitters: u64 = 6;
const ENoTransmitters: u64 = 7;
const ERepeatedSigners: u64 = 8;
const ERepeatedTransmitters: u64 = 9;
const EConfigDigestMismatch: u64 = 10;
const EUnauthorizedTransmitter: u64 = 11;
const EWrongNumberOfSignatures: u64 = 12;
const ECouldNotValidateSignerKey: u64 = 13;
const EInvalidReportContextLength: u64 = 14;
const EInvalidConfigDigestLength: u64 = 15;
const EInvalidSequenceLength: u64 = 16;
const EUnauthorizedSigner: u64 = 17;
const ENonUniqueSignatures: u64 = 18;
const EInvalidSignature: u64 = 19;
const EOutOfBytes: u64 = 20;
const EConfigNotSet: u64 = 21;
const EInvalidSignatureLength: u64 = 22;

// there is no init or initialize functions in ocr3 base
// ocr3 base state is only created and stored in offramp state
public fun new(ctx: &mut TxContext): OCR3BaseState {
    OCR3BaseState {
        id: object::new(ctx),
        ocr3_configs: table::new(ctx),
        signer_oracles: table::new(ctx),
        transmitter_oracles: table::new(ctx)
    }
}

public fun ocr_plugin_type_commit(): u8 {
    OCR_PLUGIN_TYPE_COMMIT
}

public fun ocr_plugin_type_execution(): u8 {
    OCR_PLUGIN_TYPE_EXECUTION
}

public fun set_ocr3_config(
    config_digest: vector<u8>,
    ocr_plugin_type: u8,
    big_f: u8,
    signers: vector<vector<u8>>,
    transmitters: vector<address>
) {
    event::emit(
        ConfigSet { ocr_plugin_type, config_digest, signers, transmitters, big_f }
    );
}
