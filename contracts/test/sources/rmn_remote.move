module test::rmn_remote;

use std::bcs;
use std::string::{Self, String};
use std::type_name;

use sui::address;
use sui::ecdsa_k1;
use sui::event;
use sui::hash;
use sui::vec_map::{Self, VecMap};

const SIGNATURE_NUM_BYTES: u64 = 64;
const GLOBAL_CURSE_SUBJECT: vector<u8> = x"01000000000000000000000000000001";

public struct RMNRemoteState has key, store {
    id: UID,
    local_chain_selector: u64,
    config: Config,
    config_count: u32,
    // most operations are O(n) with vec map, but it's easy to retrieve all the keys
    signers: VecMap<vector<u8>, bool>,
    cursed_subjects: VecMap<vector<u8>, bool>
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

// TODO: figure out what to do with chain_id. Cannot get it from Sui library.
public struct Report has drop {
    dest_chain_selector: u64,
    rmn_remote_contract_address: address,
    off_ramp_address: address,
    rmn_home_contract_config_digest: vector<u8>,
    merkle_roots: vector<MerkleRoot>
}

public struct MerkleRoot has drop {
    source_chain_selector: u64,
    on_ramp_address: vector<u8>,
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

// Simple methods to emit events for testing purposes

/// Emit a ConfigSet event
public fun emit_config_set_event(version: u32, config: Config) {
    event::emit(ConfigSet { version, config });
}

/// Emit a Cursed event for a single subject
public fun emit_cursed_event(subject: vector<u8>) {
    event::emit(Cursed { subjects: vector[subject] });
}

/// Emit a Cursed event for multiple subjects
public fun emit_cursed_multiple_event(subjects: vector<vector<u8>>) {
    event::emit(Cursed { subjects });
}

/// Emit an Uncursed event for a single subject
public fun emit_uncursed_event(subject: vector<u8>) {
    event::emit(Uncursed { subjects: vector[subject] });
}

/// Emit an Uncursed event for multiple subjects
public fun emit_uncursed_multiple_event(subjects: vector<vector<u8>>) {
    event::emit(Uncursed { subjects });
}
