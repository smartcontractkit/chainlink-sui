module ccip::rmn_remote;

use ccip::eth_abi;
use ccip::ownable::{Self, OwnerCap};
use ccip::state_object::{Self, CCIPObjectRef};
use ccip::upgrade_registry::verify_function_allowed;
use mcms::bcs_stream;
use mcms::mcms_registry::{Self, Registry, ExecutingCallbackParams};
use fast_mcms::mcms_registry as fast_mcms_registry;
use fast_mcms::mcms_registry::Registry as FastRegistry;
use fast_mcms::mcms_registry::ExecutingCallbackParams as FastExecutingCallbackParams;
use std::bcs;
use std::string::{Self, String};
use sui::event;
use sui::hash;
use sui::vec_map::{Self, VecMap};

const GLOBAL_CURSE_SUBJECT: vector<u8> = x"01000000000000000000000000000001";

public struct RMNRemoteState has key, store {
    id: UID,
    local_chain_selector: u64,
    config: Config,
    config_count: u32,
    // most operations are O(n) with vec map, but it's easy to retrieve all the keys
    signers: VecMap<vector<u8>, bool>,
    cursed_subjects: VecMap<vector<u8>, bool>,
}

public struct Config has copy, drop, store {
    rmn_home_contract_config_digest: vector<u8>,
    signers: vector<Signer>,
    f_sign: u64,
}

public struct Signer has copy, drop, store {
    onchain_public_key: vector<u8>,
    node_index: u64,
}

public struct ConfigSet has copy, drop {
    version: u32,
    config: Config,
}

public struct Cursed has copy, drop {
    subjects: vector<vector<u8>>,
}

public struct Uncursed has copy, drop {
    subjects: vector<vector<u8>>,
}

/// Capability granting curse-only authority. Minted by the CCIP owner via
/// `create_curser_cap` and intended to be registered as the package cap in a
/// secondary MCMS instance whose only allowed module on `@ccip` is `rmn_remote`.
/// Possession of `&CurserCap` is the authority for `curse_with_curser_cap` and
/// `curse_multiple_with_curser_cap`; uncurse is not reachable through this cap.
public struct CurserCap has key, store {
    id: UID,
}

const EAlreadyInitialized: u64 = 1;
const EAlreadyCursed: u64 = 2;
const EDuplicateSigner: u64 = 3;
const EInvalidSignerOrder: u64 = 4;
const ENotEnoughSigners: u64 = 5;
const ENotCursed: u64 = 6;
const EZeroValueNotAllowed: u64 = 7;
const EInvalidDigestLength: u64 = 8;
const ESignersMismatch: u64 = 9;
const EInvalidSubjectLength: u64 = 10;
const EInvalidPublicKeyLength: u64 = 11;
const EInvalidFunction: u64 = 12;
const EInvalidOwnerCap: u64 = 13;

const VERSION: u8 = 1;

public fun type_and_version(): String {
    string::utf8(b"RMNRemote 1.6.0")
}

public fun initialize(
    ref: &mut CCIPObjectRef,
    owner_cap: &OwnerCap,
    local_chain_selector: u64,
    ctx: &mut TxContext,
) {
    assert!(object::id(owner_cap) == state_object::owner_cap_id(ref), EInvalidOwnerCap);
    assert!(!state_object::contains<RMNRemoteState>(ref), EAlreadyInitialized);
    assert!(local_chain_selector != 0, EZeroValueNotAllowed);

    let state = RMNRemoteState {
        id: object::new(ctx),
        local_chain_selector,
        config: Config {
            rmn_home_contract_config_digest: vector[],
            signers: vector[],
            f_sign: 0,
        },
        config_count: 0,
        signers: vec_map::empty<vector<u8>, bool>(),
        cursed_subjects: vec_map::empty<vector<u8>, bool>(),
    };

    state_object::add(ref, owner_cap, state, ctx);
}

public fun set_config(
    ref: &mut CCIPObjectRef,
    owner_cap: &OwnerCap,
    rmn_home_contract_config_digest: vector<u8>,
    signer_onchain_public_keys: vector<vector<u8>>,
    node_indexes: vector<u64>,
    f_sign: u64,
) {
    verify_function_allowed(
        ref,
        string::utf8(b"rmn_remote"),
        string::utf8(b"set_config"),
        VERSION,
    );

    assert!(object::id(owner_cap) == state_object::owner_cap_id(ref), EInvalidOwnerCap);

    let state = state_object::borrow_mut<RMNRemoteState>(ref);

    assert!(rmn_home_contract_config_digest.length() == 32, EInvalidDigestLength);

    assert!(eth_abi::decode_u256_value(rmn_home_contract_config_digest) != 0, EZeroValueNotAllowed);

    let signers_len = signer_onchain_public_keys.length();
    assert!(signers_len == node_indexes.length(), ESignersMismatch);

    let mut i = 1;
    while (i < signers_len) {
        let previous_node_index = node_indexes[i - 1];
        let current_node_index = node_indexes[i];
        assert!(previous_node_index < current_node_index, EInvalidSignerOrder);
        i = i + 1;
    };

    assert!(signers_len >= (2 * f_sign + 1), ENotEnoughSigners);

    let keys = state.signers.keys();
    let mut i = 0;
    let keys_len = keys.length();
    while (i < keys_len) {
        let key = keys[i];
        state.signers.remove(&key);
        i = i + 1;
    };

    let signers = signer_onchain_public_keys.zip_map_ref!(
        &node_indexes,
        |signer_public_key_bytes, node_indexes| {
            let signer_public_key_bytes: vector<u8> = *signer_public_key_bytes;
            let node_index: u64 = *node_indexes;
            // expect an ethereum address of 20 bytes.
            assert!(signer_public_key_bytes.length() == 20, EInvalidPublicKeyLength);
            assert!(!state.signers.contains(&signer_public_key_bytes), EDuplicateSigner);
            state.signers.insert(signer_public_key_bytes, true);
            Signer {
                onchain_public_key: signer_public_key_bytes,
                node_index,
            }
        },
    );

    let new_config = Config {
        rmn_home_contract_config_digest,
        signers,
        f_sign,
    };
    state.config = new_config;

    let new_config_count = state.config_count + 1;
    state.config_count = new_config_count;

    event::emit(ConfigSet { version: new_config_count, config: new_config });
}

public fun get_versioned_config(ref: &CCIPObjectRef): (u32, Config) {
    verify_function_allowed(
        ref,
        string::utf8(b"rmn_remote"),
        string::utf8(b"get_versioned_config"),
        VERSION,
    );
    let state = state_object::borrow<RMNRemoteState>(ref);

    (state.config_count, state.config)
}

public fun get_local_chain_selector(ref: &CCIPObjectRef): u64 {
    verify_function_allowed(
        ref,
        string::utf8(b"rmn_remote"),
        string::utf8(b"get_local_chain_selector"),
        VERSION,
    );
    let state = state_object::borrow<RMNRemoteState>(ref);

    state.local_chain_selector
}

public fun get_report_digest_header(): vector<u8> {
    hash::keccak256(&b"RMN_V1_6_ANY2SUI_REPORT")
}

public fun curse(ref: &mut CCIPObjectRef, owner_cap: &OwnerCap, subject: vector<u8>) {
    verify_function_allowed(
        ref,
        string::utf8(b"rmn_remote"),
        string::utf8(b"curse"),
        VERSION,
    );
    assert!(object::id(owner_cap) == state_object::owner_cap_id(ref), EInvalidOwnerCap);

    curse_multiple(ref, owner_cap, vector[subject]);
}

public fun curse_multiple(
    ref: &mut CCIPObjectRef,
    owner_cap: &OwnerCap,
    subjects: vector<vector<u8>>,
) {
    verify_function_allowed(
        ref,
        string::utf8(b"rmn_remote"),
        string::utf8(b"curse_multiple"),
        VERSION,
    );
    assert!(object::id(owner_cap) == state_object::owner_cap_id(ref), EInvalidOwnerCap);

    let state = state_object::borrow_mut<RMNRemoteState>(ref);
    insert_cursed_subjects(state, subjects);
}

/// Shared state mutation for every curse path. Keeps OwnerCap-based and
/// CurserCap-based callers on identical semantics so the two auth paths cannot
/// drift apart on validation, insertion, or event emission.
fun insert_cursed_subjects(state: &mut RMNRemoteState, subjects: vector<vector<u8>>) {
    subjects.do_ref!(|subject| {
        let subject: vector<u8> = *subject;
        assert!(subject.length() == 16, EInvalidSubjectLength);
        assert!(!state.cursed_subjects.contains(&subject), EAlreadyCursed);
        state.cursed_subjects.insert(subject, true);
    });
    event::emit(Cursed { subjects });
}

/// Curse a single subject via possession of a `CurserCap`. Intended for use by
/// a secondary MCMS instance whose Registry holds a `CurserCap` for `@ccip`.
public fun curse_with_curser_cap(
    ref: &mut CCIPObjectRef,
    _curser_cap: &CurserCap,
    subject: vector<u8>,
) {
    verify_function_allowed(
        ref,
        string::utf8(b"rmn_remote"),
        string::utf8(b"curse_with_curser_cap"),
        VERSION,
    );
    let state = state_object::borrow_mut<RMNRemoteState>(ref);
    insert_cursed_subjects(state, vector[subject]);
}

/// Curse multiple subjects via possession of a `CurserCap`.
public fun curse_multiple_with_curser_cap(
    ref: &mut CCIPObjectRef,
    _curser_cap: &CurserCap,
    subjects: vector<vector<u8>>,
) {
    verify_function_allowed(
        ref,
        string::utf8(b"rmn_remote"),
        string::utf8(b"curse_multiple_with_curser_cap"),
        VERSION,
    );

    let state = state_object::borrow_mut<RMNRemoteState>(ref);
    insert_cursed_subjects(state, subjects);
}

/// Mint a new `CurserCap`. Owner-only. The fresh cap can be registered as the
/// package cap on a secondary MCMS Registry to enable curse-only governance
/// without granting that instance any other CCIP authority.
public fun create_curser_cap(
    ref: &mut CCIPObjectRef,
    owner_cap: &OwnerCap,
    ctx: &mut TxContext,
): CurserCap {
    verify_function_allowed(
        ref,
        string::utf8(b"rmn_remote"),
        string::utf8(b"create_curser_cap"),
        VERSION,
    );
    assert!(object::id(owner_cap) == state_object::owner_cap_id(ref), EInvalidOwnerCap);

    CurserCap { id: object::new(ctx) }
}

public fun uncurse(ref: &mut CCIPObjectRef, owner_cap: &OwnerCap, subject: vector<u8>) {
    verify_function_allowed(
        ref,
        string::utf8(b"rmn_remote"),
        string::utf8(b"uncurse"),
        VERSION,
    );
    assert!(object::id(owner_cap) == state_object::owner_cap_id(ref), EInvalidOwnerCap);

    uncurse_multiple(ref, owner_cap, vector[subject]);
}

public fun uncurse_multiple(
    ref: &mut CCIPObjectRef,
    owner_cap: &OwnerCap,
    subjects: vector<vector<u8>>,
) {
    verify_function_allowed(
        ref,
        string::utf8(b"rmn_remote"),
        string::utf8(b"uncurse_multiple"),
        VERSION,
    );
    assert!(object::id(owner_cap) == state_object::owner_cap_id(ref), EInvalidOwnerCap);

    let state = state_object::borrow_mut<RMNRemoteState>(ref);

    subjects.do_ref!(|subject| {
        let subject: vector<u8> = *subject;
        assert!(state.cursed_subjects.contains(&subject), ENotCursed);
        state.cursed_subjects.remove(&subject);
    });
    event::emit(Uncursed { subjects });
}

public fun get_cursed_subjects(ref: &CCIPObjectRef): vector<vector<u8>> {
    verify_function_allowed(
        ref,
        string::utf8(b"rmn_remote"),
        string::utf8(b"get_cursed_subjects"),
        VERSION,
    );
    let state = state_object::borrow<RMNRemoteState>(ref);

    state.cursed_subjects.keys()
}

#[allow(implicit_const_copy)]
public fun is_cursed_global(ref: &CCIPObjectRef): bool {
    verify_function_allowed(
        ref,
        string::utf8(b"rmn_remote"),
        string::utf8(b"is_cursed_global"),
        VERSION,
    );
    let state = state_object::borrow<RMNRemoteState>(ref);

    state.cursed_subjects.contains(&GLOBAL_CURSE_SUBJECT)
}

public fun is_cursed(ref: &CCIPObjectRef, subject: vector<u8>): bool {
    verify_function_allowed(
        ref,
        string::utf8(b"rmn_remote"),
        string::utf8(b"is_cursed"),
        VERSION,
    );
    let state = state_object::borrow<RMNRemoteState>(ref);

    state.cursed_subjects.contains(&subject) || is_cursed_global(ref)
}

public fun is_cursed_u128(ref: &CCIPObjectRef, subject_value: u128): bool {
    verify_function_allowed(
        ref,
        string::utf8(b"rmn_remote"),
        string::utf8(b"is_cursed_u128"),
        VERSION,
    );
    let mut subject = bcs::to_bytes(&subject_value);
    subject.reverse();
    is_cursed(ref, subject)
}

// ================================================================
// |                      MCMS Entrypoints                       |
// ================================================================

public fun mcms_set_config(
    ref: &mut CCIPObjectRef,
    registry: &mut Registry,
    params: ExecutingCallbackParams,
) {
    let (owner_cap, function, data) = mcms_registry::get_callback_params_with_caps<
        state_object::McmsCallback,
        OwnerCap,
    >(
        registry,
        state_object::mcms_callback(),
        params,
    );
    assert!(function == string::utf8(b"set_config"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    bcs_stream::validate_obj_addrs(
        vector[object::id_address(ref), object::id_address(owner_cap)],
        &mut stream,
    );

    let rmn_home_contract_config_digest = bcs_stream::deserialize_vector_u8(&mut stream);
    let signer_onchain_public_keys = bcs_stream::deserialize_vector!(
        &mut stream,
        |stream| { bcs_stream::deserialize_vector_u8(stream) },
    );
    let node_indexes = bcs_stream::deserialize_vector!(
        &mut stream,
        |stream| { bcs_stream::deserialize_u64(stream) },
    );
    let f_sign = bcs_stream::deserialize_u64(&mut stream);
    bcs_stream::assert_is_consumed(&stream);

    set_config(
        ref,
        owner_cap,
        rmn_home_contract_config_digest,
        signer_onchain_public_keys,
        node_indexes,
        f_sign,
    );
}

public fun mcms_curse(
    ref: &mut CCIPObjectRef,
    registry: &mut Registry,
    params: ExecutingCallbackParams,
) {
    let (owner_cap, function, data) = mcms_registry::get_callback_params_with_caps<
        state_object::McmsCallback,
        OwnerCap,
    >(
        registry,
        state_object::mcms_callback(),
        params,
    );
    assert!(function == string::utf8(b"curse"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    bcs_stream::validate_obj_addrs(
        vector[object::id_address(ref), object::id_address(owner_cap)],
        &mut stream,
    );

    let subject = bcs_stream::deserialize_vector_u8(&mut stream);
    bcs_stream::assert_is_consumed(&stream);

    curse(ref, owner_cap, subject);
}

public fun mcms_curse_multiple(
    ref: &mut CCIPObjectRef,
    registry: &mut Registry,
    params: ExecutingCallbackParams,
) {
    let (owner_cap, function, data) = mcms_registry::get_callback_params_with_caps<
        state_object::McmsCallback,
        OwnerCap,
    >(
        registry,
        state_object::mcms_callback(),
        params,
    );
    assert!(function == string::utf8(b"curse_multiple"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    bcs_stream::validate_obj_addrs(
        vector[object::id_address(ref), object::id_address(owner_cap)],
        &mut stream,
    );

    let subjects = bcs_stream::deserialize_vector!(
        &mut stream,
        |stream| { bcs_stream::deserialize_vector_u8(stream) },
    );
    bcs_stream::assert_is_consumed(&stream);

    curse_multiple(ref, owner_cap, subjects);
}

public fun mcms_uncurse(
    ref: &mut CCIPObjectRef,
    registry: &mut Registry,
    params: ExecutingCallbackParams,
) {
    let (owner_cap, function, data) = mcms_registry::get_callback_params_with_caps<
        state_object::McmsCallback,
        OwnerCap,
    >(
        registry,
        state_object::mcms_callback(),
        params,
    );
    assert!(function == string::utf8(b"uncurse"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    bcs_stream::validate_obj_addrs(
        vector[object::id_address(ref), object::id_address(owner_cap)],
        &mut stream,
    );

    let subject = bcs_stream::deserialize_vector_u8(&mut stream);
    bcs_stream::assert_is_consumed(&stream);

    uncurse(ref, owner_cap, subject);
}

public fun mcms_uncurse_multiple(
    ref: &mut CCIPObjectRef,
    registry: &mut Registry,
    params: ExecutingCallbackParams,
) {
    let (owner_cap, function, data) = mcms_registry::get_callback_params_with_caps<
        state_object::McmsCallback,
        OwnerCap,
    >(
        registry,
        state_object::mcms_callback(),
        params,
    );
    assert!(function == string::utf8(b"uncurse_multiple"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    bcs_stream::validate_obj_addrs(
        vector[object::id_address(ref), object::id_address(owner_cap)],
        &mut stream,
    );

    let subjects = bcs_stream::deserialize_vector!(
        &mut stream,
        |stream| { bcs_stream::deserialize_vector_u8(stream) },
    );
    bcs_stream::assert_is_consumed(&stream);

    uncurse_multiple(ref, owner_cap, subjects);
}

// ================================================================
// |               Fast Curse MCMS Entrypoints                    |
// ================================================================
//
// These callbacks are intended to be invoked from a SECONDARY MCMS instance
// whose `Registry` holds a `CurserCap` (not `OwnerCap`) for `@ccip` with
// `allowed_modules = [b"rmn_remote"]`. They use `state_object::McmsCallback`
// as the proof type because the CCIP `Publisher` was claimed by the
// `STATE_OBJECT` OTW; `sui::package::Publisher::from_module<T>` requires the
// proof type's defining module to match the Publisher's `module_name` field.
// The two MCMS instances are separate shared `Registry` objects, so reusing
// the proof type does not create a collision: the slow instance holds
// `(state_object::McmsCallback, OwnerCap)` for `@ccip` and the fast instance
// holds `(state_object::McmsCallback, CurserCap)` for `@ccip`.

/// Fast-path callback. Releases the `CurserCap` from the fast Registry and
/// curses a single subject. The Registry's `allowed_modules` allowlist plus
/// the absence of any `uncurse`-bearing callback in this entrypoint family
/// constrains the fast MCMS to curse operations only.
public fun mcms_curse_with_curser_cap(
    ref: &mut CCIPObjectRef,
    registry: &mut FastRegistry,
    params: FastExecutingCallbackParams,
) {
    let (curser_cap, function, data) = fast_mcms_registry::get_callback_params_with_caps<
        state_object::McmsCallback,
        CurserCap,
    >(
        registry,
        state_object::mcms_callback(),
        params,
    );
    assert!(function == string::utf8(b"curse_with_curser_cap"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    bcs_stream::validate_obj_addrs(
        vector[object::id_address(ref), object::id_address(curser_cap)],
        &mut stream,
    );

    let subject = bcs_stream::deserialize_vector_u8(&mut stream);
    bcs_stream::assert_is_consumed(&stream);

    curse_with_curser_cap(ref, curser_cap, subject);
}

public fun mcms_curse_multiple_with_curser_cap(
    ref: &mut CCIPObjectRef,
    registry: &mut FastRegistry,
    params: FastExecutingCallbackParams,
) {
    let (curser_cap, function, data) = fast_mcms_registry::get_callback_params_with_caps<
        state_object::McmsCallback,
        CurserCap,
    >(
        registry,
        state_object::mcms_callback(),
        params,
    );
    assert!(function == string::utf8(b"curse_multiple_with_curser_cap"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    bcs_stream::validate_obj_addrs(
        vector[object::id_address(ref), object::id_address(curser_cap)],
        &mut stream,
    );

    let subjects = bcs_stream::deserialize_vector!(
        &mut stream,
        |stream| { bcs_stream::deserialize_vector_u8(stream) },
    );
    bcs_stream::assert_is_consumed(&stream);

    curse_multiple_with_curser_cap(ref, curser_cap, subjects);
}

// ================================================================
// |          Slow-MCMS Bootstrap Wrappers for Fast Curse          |
// ================================================================
//
// These callbacks are invoked by the slow MCMS (which holds CCIP `OwnerCap`)
// to mint a `CurserCap` and register it in the fast MCMS Registry. The slow
// proposal can batch a mint op then a register op (`mcms_create_curser_cap`
// → `mcms_register_curser_cap`), or use the single combined op
// (`mcms_mint_and_register_curser_cap`).
//
// TODO: consolidate the two-op flow into the single combined flow once the
// PTB-chain ergonomics and external proposal tooling are agreed; the combined
// path is strictly safer because no transient `CurserCap` exists between txs.

/// Slow-MCMS wrapper that mints a `CurserCap`. Returns the cap by value so a
/// PTB can chain it into `mcms_register_curser_cap` in the same batch.
public fun mcms_create_curser_cap(
    ref: &mut CCIPObjectRef,
    registry: &mut Registry,
    params: ExecutingCallbackParams,
    ctx: &mut TxContext,
): CurserCap {
    let (owner_cap, function, data) = mcms_registry::get_callback_params_with_caps<
        state_object::McmsCallback,
        OwnerCap,
    >(
        registry,
        state_object::mcms_callback(),
        params,
    );
    assert!(function == string::utf8(b"create_curser_cap"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    bcs_stream::validate_obj_addrs(
        vector[object::id_address(ref), object::id_address(owner_cap)],
        &mut stream,
    );
    bcs_stream::assert_is_consumed(&stream);

    create_curser_cap(ref, owner_cap, ctx)
}

/// Slow-MCMS wrapper that registers a `CurserCap` as the package cap on a
/// fast MCMS Registry with `allowed_modules = [b"rmn_remote"]`. Built from
/// the OwnerCap's stored `Publisher` so the proof type binding is rooted in
/// the same module that claimed the CCIP Publisher (`state_object`).
public fun mcms_register_curser_cap(
    ref: &mut CCIPObjectRef,
    slow_registry: &mut Registry,
    fast_registry: &mut FastRegistry,
    params: ExecutingCallbackParams,
    curser_cap: CurserCap,
    ctx: &mut TxContext,
) {
    let (owner_cap, function, data) = mcms_registry::get_callback_params_with_caps<
        state_object::McmsCallback,
        OwnerCap,
    >(
        slow_registry,
        state_object::mcms_callback(),
        params,
    );
    assert!(function == string::utf8(b"register_curser_cap"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    bcs_stream::validate_obj_addrs(
        vector[
            object::id_address(ref),
            object::id_address(owner_cap),
            object::id_address(fast_registry),
        ],
        &mut stream,
    );
    bcs_stream::assert_is_consumed(&stream);

    verify_function_allowed(
        ref,
        string::utf8(b"rmn_remote"),
        string::utf8(b"register_curser_cap"),
        VERSION,
    );

    let publisher_wrapper = fast_mcms_registry::create_publisher_wrapper(
        ownable::borrow_publisher(owner_cap),
        state_object::mcms_callback(),
    );

    fast_mcms_registry::register_entrypoint(
        fast_registry,
        publisher_wrapper,
        state_object::mcms_callback(),
        curser_cap,
        vector[b"rmn_remote"],
        ctx,
    );
}

/// Slow-MCMS wrapper that mints a `CurserCap` and atomically registers it in
/// the fast MCMS Registry. Preferred over the two-op variant: no PTB caller
/// can substitute a different `CurserCap` between mint and register.
public fun mcms_mint_and_register_curser_cap(
    ref: &mut CCIPObjectRef,
    slow_registry: &mut Registry,
    fast_registry: &mut FastRegistry,
    params: ExecutingCallbackParams,
    ctx: &mut TxContext,
) {
    let (owner_cap, function, data) = mcms_registry::get_callback_params_with_caps<
        state_object::McmsCallback,
        OwnerCap,
    >(
        slow_registry,
        state_object::mcms_callback(),
        params,
    );
    assert!(function == string::utf8(b"mint_and_register_curser_cap"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    bcs_stream::validate_obj_addrs(
        vector[
            object::id_address(ref),
            object::id_address(owner_cap),
            object::id_address(fast_registry),
        ],
        &mut stream,
    );
    bcs_stream::assert_is_consumed(&stream);

    verify_function_allowed(
        ref,
        string::utf8(b"rmn_remote"),
        string::utf8(b"mint_and_register_curser_cap"),
        VERSION,
    );

    let publisher_wrapper = fast_mcms_registry::create_publisher_wrapper(
        ownable::borrow_publisher(owner_cap),
        state_object::mcms_callback(),
    );

    assert!(object::id(owner_cap) == state_object::owner_cap_id(ref), EInvalidOwnerCap);
    let curser_cap = CurserCap { id: object::new(ctx) };

    fast_mcms_registry::register_entrypoint(
        fast_registry,
        publisher_wrapper,
        state_object::mcms_callback(),
        curser_cap,
        vector[b"rmn_remote"],
        ctx,
    );
}

#[test_only]
public fun get_config(config: &Config): (vector<u8>, vector<Signer>, u64) {
    (config.rmn_home_contract_config_digest, config.signers, config.f_sign)
}
