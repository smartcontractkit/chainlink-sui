# CCIP Sui Receiver Integration Guide

This guide explains how to build a CCIP message receiver on Sui that remains safe under permissionless manual execution. It applies to the V1 destination execution format shipped today.

Reference implementation: [`ccip_dummy_receiver`](../ccip_dummy_receiver/sources/ccip_dummy_receiver.move).

CCIP receivers are third-party applications. You are responsible for your object model, security properties, and whether your design is safe under manual execution.

### Summary

Destination V1 intentionally omits `receiver_object_ids` from the OCR execution report and merkle leaf. Source senders still encode those IDs in `extra_args`, and the relayer uses them when building the execution PTB, but the committed report on Sui does not bind them. That is deliberate: if a sender specified wrong object IDs and execution fails or the message sits unexecuted, an honest party can retry manual execution or a corrected relayer run with the right IDs instead of leaving the message permanently stuck. The tradeoff is that tail object choice is not authenticated on destination — your receiver design must prevent an executor from substituting a different valid object of the same type.

---

## Threat model

On the **source** chain, senders pass Sui-specific extra args via `client::encode_sui_extra_args_v1`. The source `message_id` commits to those bytes, including `receiver_object_ids`.

On the **destination**, the V1 offramp OCR execution report and leaf hash authenticate:

- receiver package address
- message body fields in `Any2SuiMessage`
- token destination via `token_receiver`

They do **not** authenticate which object references are passed as PTB tail arguments to `ccip_receive`.

After `permissionless_execution_threshold_seconds`, any address may call `offramp::manually_init_execute` and assemble the execution PTB. If your receiver accepts multiple valid shared objects of the same tail type, an executor can deliver a committed message to the wrong object while the protocol marks execution successful. The message cannot be replayed against the intended object.

Token delivery is unaffected: `token_receiver` is in the destination leaf hash.

---

## `receiver_object_ids`

### What it is

`receiver_object_ids` is a vector of 32-byte Sui object IDs, one per object-reference parameter in your `ccip_receive` function after the three protocol-fixed arguments:

| Index | Parameter | Set by |
|------:|-----------|--------|
| 0 | `expected_message_id: vector<u8>` | Relayer |
| 1 | `ref: &CCIPObjectRef` | Relayer |
| 2 | `message: Any2SuiMessage` | Relayer via `extract_any2sui_message` |
| 3+ | Your `&T` / `&mut T` tail args | Relayer, one ID per slot in ABI order |

Senders encode IDs on the source chain:

```move
let extra_args = client::encode_sui_extra_args_v1(
    gas_limit,
    allow_out_of_order_execution,
    token_receiver_bytes,           // 32 bytes
    vector[
        bcs::to_bytes(&clock_id),
        bcs::to_bytes(&state_id),
    ],
);
```

Count invariant:

```text
len(receiver_object_ids) == number of object-reference parameters after message
```

The relayer reads `receiverObjectIds` from source extra args metadata when building the execution PTB and appends matching object inputs to the `ccip_receive` call.

### Why IDs are not in the OCR execution report

Destination V1 deliberately does **not** include `receiver_object_ids` in the merkle leaf or deserialized execution report. That is an intentional tradeoff:

| Benefit | Cost |
|---------|------|
| Wrong or stale IDs can be corrected at manual execution time | Tail object choice is not cryptographically bound on destination |
| Receiver upgrades or singleton rotation do not permanently stick messages | Receiver design must prevent object substitution |

If a sender supplies a bad object ID that does not exist or has the wrong type, the PTB fails before submission or reverts. A manual executor or updated relayer run can retry with the correct singleton ID without the message being permanently stuck.

Binding IDs into the OCR report would make typos and operational changes permanent. The middle-path mitigation is **receiver design**, not protocol-level ID binding.

---

## Execution PTB

CCIP message execution is always a single atomic PTB. Typical command order:

1. `offramp::init_execute` or `offramp::manually_init_execute`
2. Token pool `release_or_mint` when the report includes tokens
3. `offramp_state_helper::extract_any2sui_message`
4. `your_package::ccip_receive` with tail object inputs
5. `offramp::finish_execute`

All steps succeed or the entire transaction reverts. Shared objects used in the PTB must be passed with the correct mutability: read-only refs for `&T`, mutable refs for `&mut T`.

---

## Safe receiver shapes

A receiver is safe from manual-execution object substitution when an attacker-controlled PTB cannot cause a committed message to mutate unintended state.

### Tier 1 — structurally safe

These patterns need no per-call object ID checks.

#### Stateless — no tail objects

```move
public fun ccip_receive(
    expected_message_id: vector<u8>,
    ref: &CCIPObjectRef,
    message: client::Any2SuiMessage,
) {
    let (msg_id, _, _, data, _, _, _) =
        osh::consume_any2sui_message(ref, message, MyProof {});
    assert!(msg_id == expected_message_id, EMessageIdMismatch);
    // handle data; no mutable tail state
}
```

`receiver_object_ids = []`

#### Singleton state — one shared object per tail type

```move
public struct ReceiverState has key {
    id: UID,
    accounts: Table<address, AccountData>,
}

fun init(otw: MY_PKG, ctx: &mut TxContext) {
    let state = ReceiverState { id: object::new(ctx), accounts: table::new(ctx) };
    transfer::share_object(state);
    // no other function creates ReceiverState
}

public fun ccip_receive(
    expected_message_id: vector<u8>,
    ref: &CCIPObjectRef,
    message: client::Any2SuiMessage,
    state: &mut ReceiverState,
) {
    let (msg_id, _, sender, data, _, _, _) =
        osh::consume_any2sui_message(ref, message, MyProof {});
    assert!(msg_id == expected_message_id, EMessageIdMismatch);
    let account = state.accounts.borrow_mut(parse_recipient(&data));
    apply(account, &data);
}
```

`receiver_object_ids = [state_id]`

Requirements:

- Exactly one valid instance of each tail type for the deployment lifetime
- No public or entry function that creates a second instance
- Prefer `has key` without `store` on delivery state so only your module controls share and transfer

#### Multi-singleton — several tail types, each a singleton

The dummy receiver uses `&Clock` at `@0x6` plus one `CCIPReceiverState`:

```move
public fun ccip_receive(
    expected_message_id: vector<u8>,
    ref: &CCIPObjectRef,
    message: client::Any2SuiMessage,
    clock: &Clock,
    state: &mut CCIPReceiverState,
) { /* ... */ }
```

`receiver_object_ids = [clock_object_id, state_object_id]`

System singletons like `Clock` satisfy the invariant by design.

#### Two-step pull — singleton inbox, authorized claim

`ccip_receive` only writes to a shared inbox. State owners pull later with an explicit capability:

```move
public fun ccip_receive(
    expected_message_id: vector<u8>,
    ref: &CCIPObjectRef,
    message: client::Any2SuiMessage,
    inbox: &mut Inbox,
) {
    let (msg_id, _, _, data, _, _, _) =
        osh::consume_any2sui_message(ref, message, MyProof {});
    assert!(msg_id == expected_message_id, EMessageIdMismatch);
    inbox.pending.add(msg_id, Pending { recipient: decode_recipient(&data), data });
}

public fun claim(
    inbox: &mut Inbox,
    account: &mut Account,
    cap: &AccountOwnerCap,
    expected_message_id: vector<u8>,
) {
    assert!(cap.account_id == object::id(account), EWrongCap);
    let pending = inbox.pending.remove(expected_message_id);
    assert!(pending.recipient == cap.owner, EWrongRecipient);
    apply(account, pending.data);
}
```

Binding moves to `claim`, where the owner signs the PTB.

---

## Non-singleton external objects

When logical recipients are separate shared or owned objects, do **not** pass them as mutable tail arguments unless you add explicit binding checks. Prefer one of these patterns.

### Hub routing — recommended default

One singleton hub; route using fields from the authenticated message:

```move
public fun ccip_receive(
    expected_message_id: vector<u8>,
    ref: &CCIPObjectRef,
    message: client::Any2SuiMessage,
    hub: &mut Hub,
) {
    let (msg_id, _, sender, data, _, _, _) =
        osh::consume_any2sui_message(ref, message, MyProof {});
    assert!(msg_id == expected_message_id, EMessageIdMismatch);

    let target = decode_target_address(&data);
    assert!(is_allowed_recipient(sender, target), EUnauthorizedTarget);

    let entry = hub.accounts.borrow_mut(target);
    apply(entry, &data);
}
```

The executor can only pass the one `Hub`. Routing uses `data` and `sender`, which are authenticated in `Any2SuiMessage`.

Define a versioned payload schema in `data` so senders and receivers agree on encoding.

### ID binding from message `data` — advanced

If a non-singleton object must appear as a tail arg, bind it to an ID encoded in authenticated `data` **before** any mutation:

```move
public fun ccip_receive(
    expected_message_id: vector<u8>,
    ref: &CCIPObjectRef,
    message: client::Any2SuiMessage,
    vault: &mut Vault,
) {
    let (msg_id, _, _, data, _, _, _) =
        osh::consume_any2sui_message(ref, message, MyProof {});
    assert!(msg_id == expected_message_id, EMessageIdMismatch);

    let expected_vault_id = decode_vault_id(&data);
    assert!(object::id_address(vault) == expected_vault_id, EWrongVault);

    apply(vault, &data);
}
```

This works but places security on per-function discipline. Verify the assert on every mutation path in your own code review and tests. Prefer hub routing or inbox pull when possible.

**Do not** read expected IDs from offchain relayer metadata inside the receiver. `Any2SuiMessage` does not carry source `extra_args`; only message fields are available onchain.

### Sender-pinned routing

When each source sender may only affect their own row:

```move
let entry = hub.by_sender.borrow_mut(sender_as_address);
```

Suitable for per-sender accounts without embedding a target ID in `data`. Not suitable when one sender must deliver to arbitrary recipients unless `data` carries an authorized target and you validate it.

---

## Bad patterns

### Multiple shared vaults as tail args without binding

```move
// UNSAFE: attacker passes vault_B while message was intended for vault_A
public fun ccip_receive(
    expected_message_id: vector<u8>,
    ref: &CCIPObjectRef,
    message: client::Any2SuiMessage,
    vault: &mut Vault,
) {
    let (_, _, _, data, _, _, _) =
        osh::consume_any2sui_message(ref, message, MyProof {});
    vault.balance = vault.balance + decode_amount(&data);
}
```

Any shared `Vault` of this type is a valid PTB input after the manual execution threshold.

### Public factory for tail state

```move
// UNSAFE: enables a second CCIPReceiverState for substitution
public fun create_receiver_state(ctx: &mut TxContext): CCIPReceiverState {
    CCIPReceiverState { id: object::new(ctx), /* ... */ }
}
```

Singleton receivers must create delivery state only in `init`, not via public factories.

### Trusting the relayer without onchain checks

```move
// UNSAFE: assumes the PTB caller always passes the intended object
public fun ccip_receive(
    expected_message_id: vector<u8>,
    ref: &CCIPObjectRef,
    message: client::Any2SuiMessage,
    state: &mut UserState,
) {
    let (msg_id, _, _, data, _, _, _) =
        osh::consume_any2sui_message(ref, message, MyProof {});
    assert!(msg_id == expected_message_id, EMessageIdMismatch);
    state.apply(&data);
}
```

Manual executors are not your relayer. Design for attacker-controlled PTBs.

### Per-user owned state as the only tail arg without inbox pull

Owned objects are harder to substitute in permissionless PTBs because the owner must sign. They also complicate relayer-driven execution. Use inbox + `claim` instead of relying on ownership alone.

---

## Security checklist

Before deploying your receiver, verify:

- [ ] Every type in `ccip_receive` tail args is listed
- [ ] Each tail type is a singleton, or binding is enforced before mutation, or only an inbox is mutated
- [ ] No public or entry path creates a second valid tail object
- [ ] `consume_any2sui_message` and `expected_message_id` check run before state changes
- [ ] Substitution or wrong-target tests exist in the receiver test suite
- [ ] `receiver_object_ids` length matches the tail object parameter count

---

## Quick reference

| Pattern | Tail args | `receiver_object_ids` | Manual exec safe? |
|---------|-----------|----------------------|-------------------|
| Stateless | none | `[]` | Yes |
| Singleton state | `&mut State` | `[state_id]` | Yes |
| Multi-singleton | `&Clock`, `&mut State` | `[clock_id, state_id]` | Yes |
| Hub routing | `&mut Hub` | `[hub_id]` | Yes |
| Inbox + claim | `&mut Inbox` | `[inbox_id]` | Yes |
| Multi vault tail, no assert | `&mut Vault` | `[vault_id]` | **No** |
| Multi vault + ID in `data` | `&mut Vault` | `[vault_id]` | Yes if assert enforced |

---

## Related reading

- [`ccip_dummy_receiver`](../ccip_dummy_receiver/sources/ccip_dummy_receiver.move) — canonical singleton example
- [`ccip_offramp` README](../ccip_offramp/README.md) — execution PTB flow
- [`client::encode_sui_extra_args_v1`](../ccip/sources/client.move) — source extra args encoding
