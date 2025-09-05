# Multi-Chain Multi-Sig (MCMS) System

MCMS provides secure multi-signature governance for Sui Move contracts with time-delayed execution. It allows teams to safely manage critical operations like configuration changes, ownership transfers, and package upgrades across multiple contracts.

## How MCMS Works

MCMS uses a **propose → sign → execute** workflow:

1. **Propose**: Create operations and generate a Merkle tree root
2. **Sign**: Signers approve the root with their private keys  
3. **Execute**: After timelock delay, execute the approved operations

### Core Concepts

- **Operations**: Individual function calls to be executed
- **Merkle Tree**: Batches operations together with cryptographic proof
- **Timelock**: Configurable delay before operations can be executed
- **Hot Potato**: Ensures all operations complete successfully or fail atomically

## Quick Start Guide

### 1. Setup MCMS

First, deploy and configure MCMS with your signers:

```bash
# Deploy MCMS contracts
sui client publish

# Configure signers and quorum (example: 2-of-3 multisig)
sui client call --function set_config \
  --args $MCMS_STATE \
  --args 2 \  # PROPOSER role
  --args "[$SIGNER1,$SIGNER2,$SIGNER3]" \
  --args "[0,0,0]" \  # All in group 0
  --args "[2,0,0,...]" \  # Group 0 needs 2 signatures
  --args "[0,0,1,...]"    # Group hierarchy
```

### 2. Prepare Your Contract

Add MCMS integration to contracts you want to control:

```move
// Add to your contract
public fun mcms_entrypoint(
    state: &mut YourState,
    registry: &mut Registry, 
    params: ExecutingCallbackParams,
    ctx: &mut TxContext,
) {
    let (owner_cap, function_name, data) = 
        mcms_registry::get_callback_params<YourProof, OwnerCap>(
            registry, YourProof {}, params
        );
    
    if (function_name == string::utf8(b"update_config")) {
        // Deserialize data and call your function
        let new_value = deserialize_u64(data);
        your_update_function(state, owner_cap, new_value);
    }
    // Handle other functions...
}

// Register with MCMS during deployment
mcms_registry::register_entrypoint(registry, YourProof {}, owner_cap, ctx);
```

### 3. Create and Execute Operations

#### Step 1: Create Operations Off-chain

```typescript
// Define what you want to do
const operations = [
    {
        target: YOUR_CONTRACT_ADDRESS,
        module_name: "your_module", 
        function_name: "update_config",
        data: encodeU64(42), // New config value
    }
];

// Generate merkle tree
const { merkleTree, root } = createMerkleTree(operations);
```

#### Step 2: Get Signatures

```typescript
// Each signer signs the root hash
const message = ethers.utils.solidityKeccak256(
    ["bytes32", "uint256"], 
    [root, validUntil]
);

const signatures = await Promise.all(
    signers.map(signer => signer.signMessage(message))
);
```

#### Step 3: Submit Root

```move
// Submit signed root to MCMS
mcms::set_root(
    state,
    clock,
    2, // PROPOSER role
    root,
    valid_until,
    chain_id,
    multisig_address,
    current_op_count,
    new_op_count, 
    false, // don't override
    metadata_proof,
    signatures,
    ctx
);
```

#### Step 4: Execute Operations

```move
// Execute each operation from the approved tree
let hot_potato = mcms::execute(
    state,
    clock, 
    2, // PROPOSER role
    chain_id,
    multisig_address,
    nonce,
    YOUR_CONTRACT_ADDRESS,
    string::utf8(b"mcms"),
    string::utf8(b"timelock_schedule_batch"),
    operation_data,
    merkle_proof,
);

// Schedule with timelock
mcms::dispatch_timelock_schedule_batch(timelock, clock, hot_potato, ctx);

// Later, after delay...
let hot_potato = mcms::execute(/* same params but timelock_execute_batch */);
let callbacks = mcms::dispatch_timelock_execute_batch(timelock, clock, hot_potato, ctx);

// Execute in your contract
your_contract::mcms_entrypoint(your_state, registry, callbacks[0], ctx);
```

## Common Use Cases

### 1. Configuration Updates

Update contract parameters safely:

```typescript
// Create operation to update a configuration value
const updateConfigOp = {
    target: CONTRACT_ADDRESS,
    module_name: "your_contract",
    function_name: "update_config", 
    data: encodeParameters(["uint64"], [newValue])
};
```

### 2. Ownership Transfer

Transfer contract ownership through MCMS:

```typescript
const transferOwnershipOp = {
    target: CONTRACT_ADDRESS,
    module_name: "your_contract", 
    function_name: "transfer_ownership",
    data: encodeParameters(["address"], [newOwner])
};
```

### 3. Package Upgrades

Safely upgrade contract packages:

```typescript
const upgradeOp = {
    target: MCMS_DEPLOYER_ADDRESS,
    module_name: "mcms_deployer",
    function_name: "authorize_upgrade", 
    data: encodeParameters(["u8", "vector<u8>", "address"], [policy, digest, packageId])
};
```

### 4. Emergency Actions

Bypass timelock in emergencies (requires BYPASSER role):

```move
// Skip timelock delay for urgent fixes
let hot_potato = mcms::execute(state, clock, BYPASSER_ROLE, /* ... */);
let callbacks = mcms::dispatch_timelock_bypasser_execute_batch(hot_potato, ctx);
your_contract::mcms_entrypoint(state, registry, callbacks[0], ctx);
```

## Roles and Permissions

| Role | Value | Purpose | Can Do |
|------|-------|---------|---------|
| **PROPOSER** | 2 | Normal operations | Schedule operations with timelock delay |
| **BYPASSER** | 0 | Emergency actions | Execute operations immediately |
| **CANCELLER** | 1 | Operation control | Cancel pending operations |

## Complete Example Workflow

Here's a complete example of updating a contract configuration:

### 1. Off-chain Preparation

```typescript
// 1. Define the operation
const operations = [{
    target: "0x123...", // Your contract address  
    module_name: "config",
    function_name: "set_fee_rate",
    data: encodeU64(250), // New fee rate: 2.5%
}];

// 2. Create merkle tree
const tree = new MerkleTree(operations);
const root = tree.getRoot();

// 3. Get signatures (each signer signs root + expiry)
const message = keccak256(encodePackedParams(["bytes32", "uint256"], [root, validUntil]));
const signatures = await getMultisigSignatures(signers, message);
```

### 2. On-chain Execution

```bash
# 1. Submit the signed root
sui client call --function set_root \
  --args $MCMS_STATE $CLOCK 2 $ROOT $VALID_UNTIL ... $SIGNATURES

# 2. Schedule the operation (with timelock delay)
sui client call --function execute \
  --args $MCMS_STATE $CLOCK 2 ... $MERKLE_PROOF | \
sui client call --function dispatch_timelock_schedule_batch

# 3. Wait for timelock delay...

# 4. Execute the operation
sui client call --function execute \
  --args $MCMS_STATE $CLOCK 2 ... $MERKLE_PROOF | \
sui client call --function dispatch_timelock_execute_batch | \
sui client call --function mcms_entrypoint --args $YOUR_CONTRACT
```

The operation is now complete! Your contract's fee rate has been updated to 2.5%.

## Advanced Features

### Hierarchical Signer Groups

MCMS supports complex approval hierarchies:

```bash
# Example: 2-of-3 setup with nested groups
# Group 0 (root): Needs 2 of groups 1,2
# Group 1: Needs 1 of 2 signers  
# Group 2: Needs 1 of 1 signer

sui client call --function set_config \
  --args 2 \  # PROPOSER role
  --args "[$ADDR1,$ADDR2,$ADDR3]" \  # Signer addresses
  --args "[1,1,2]" \  # Group assignments
  --args "[2,1,1,0,...]" \  # Group quorums  
  --args "[0,0,0,1,...]"     # Group parents
```

### Batch Operations

Execute multiple operations atomically:

```typescript
const batchOps = [
    { target: CONTRACT_A, function_name: "update_fee", data: encodeU64(100) },
    { target: CONTRACT_B, function_name: "set_admin", data: encodeAddress(newAdmin) }, 
    { target: CONTRACT_C, function_name: "pause", data: encodeBool(true) }
];

const tree = new MerkleTree(batchOps);
// All operations execute together or all fail
```

### Operation Dependencies

Ensure operations execute in order:

```move
// Use predecessor hash to create dependencies
let operation1_hash = hash_operation_batch(calls1, ZERO_HASH, salt1);
let operation2_hash = hash_operation_batch(calls2, operation1_hash, salt2);
// Operation 2 can only execute after operation 1 completes
```

## Troubleshooting

### Common Issues

**"EInsufficientSigners" (Error 12)**
- Not enough signers approved the root
- Check that quorum requirements are met
- Verify all signatures are valid and sorted

**"EValidUntilExpired" (Error 8)**  
- The root has expired
- Submit a new root with later expiry time
- Check clock synchronization

**"EProofCannotBeVerified" (Error 6)**
- Merkle proof doesn't match the root
- Verify operation data matches exactly
- Check proof generation logic

**"EOperationNotReady" (Error 30)**
- Timelock delay hasn't passed yet
- Wait for the full delay period
- Check operation was properly scheduled

### Debug Commands

```bash
# Check current root and operation count
sui client call --function get_root --args $MCMS_STATE 2

# Verify operation is scheduled
sui client call --function timelock_is_operation_ready --args $TIMELOCK $OPERATION_ID

# Check signer configuration
sui client call --function get_config --args $MCMS_STATE 2
```

## Best Practices

### Security
- **Set appropriate timelock delays** (e.g., 24 hours for config changes)
- **Use different roles appropriately** (BYPASSER only for emergencies)
- **Test all operations on devnet first** 
- **Monitor for unexpected proposals**
- **Maintain offline backup keys**

### Operations  
- **Coordinate signer availability** before creating proposals
- **Document all operations** with clear descriptions
- **Batch related changes** to reduce coordination overhead
- **Plan rollback procedures** for critical changes
- **Set reasonable root expiry times** (e.g., 7 days)

### Development
- **Keep `mcms_entrypoint` functions simple** and focused
- **Validate all input data** before executing operations
- **Use consistent function naming** across contracts
- **Handle partial failures gracefully**
- **Add operation logging** for auditability

## How Quorum and Signers Work

MCMS uses a sophisticated hierarchical group system for signature validation that supports complex governance structures.

### Signer Groups and Hierarchy

MCMS organizes signers into **up to 32 groups** (indexed 0-31) arranged in a **tree structure**:

- **Group 0** is always the **root group**
- Each group has a **parent group** (except Group 0, which is its own parent)
- Each group has a **quorum requirement** (number of approvals needed)
- Groups form a tree where children must satisfy their quorum before contributing to parent groups

### How Quorum Validation Works

When signatures are submitted via `set_root()`, MCMS validates them through this process:

1. **Signature Recovery**: Each signature is validated and the signer address recovered
2. **Address Sorting**: Signatures must be sorted by signer address (ascending order)
3. **Group Counting**: For each valid signature:
   - Find the signer's group
   - Increment vote count for that group
   - If group reaches its quorum, increment vote count for parent group
   - Continue up the tree until root group or quorum not met
4. **Root Validation**: The root group (Group 0) must meet its quorum for approval

### Configuration Examples

#### Simple 2-of-3 Multisig
```bash
# 3 signers, all in group 0, need 2 signatures
signer_addresses = [addr1, addr2, addr3]
signer_groups = [0, 0, 0]           # All in group 0
group_quorums = [2, 0, 0, ...]      # Group 0 needs 2, others disabled
group_parents = [0, 0, 1, ...]      # Group 0 is root
```

#### Hierarchical 3-of-5 with Departments
```bash
# Engineering (2 people, need 1) + Security (2 people, need 1) + CEO (1 person)
# Root needs 2 of the 3 groups
signer_addresses = [eng1, eng2, sec1, sec2, ceo]
signer_groups = [1, 1, 2, 2, 3]     # Groups 1=eng, 2=sec, 3=ceo  
group_quorums = [2, 1, 1, 1, 0...]  # Root needs 2, each dept needs 1
group_parents = [0, 0, 0, 0, 1...]  # Groups 1,2,3 report to root (0)
```

#### Complex Nested Structure
```bash
# Multi-tier approval: Team leads → Department heads → C-level
signer_addresses = [team1, team2, dept1, dept2, cto, ceo]
signer_groups = [3, 4, 2, 2, 1, 1]     # Nested groups
group_quorums = [2, 1, 1, 1, 1, 0...]  # Various quorum requirements
group_parents = [0, 0, 1, 3, 2, 2...]  # Tree structure
```

### Signature Process Detail

#### 1. Message Construction
```move
// Signers sign: keccak256(abi.encode(root, valid_until))
let message_hash = compute_eth_message_hash(root, valid_until);
```

#### 2. Off-chain Signing
```typescript
// Each signer creates an Ethereum-compatible signature
const message = ethers.utils.arrayify(messageHash);
const signature = await signer.signMessage(message);
```

#### 3. On-chain Validation
```move
// MCMS validates each signature in order
let signer_addr = ecdsa_recover_evm_addr(signed_hash, signature);
assert!(signer_addr > prev_signer_addr, ESignerAddrMustBeIncreasing);
assert!(multisig.signers.contains(&signer_addr), EInvalidSigner);
```

### Group Vote Counting Algorithm

Here's how MCMS counts votes through the hierarchy:

```move
// For each valid signature:
let mut group = signer.group;
while (true) {
    // Increment vote count for current group
    group_vote_counts[group] += 1;
    
    let quorum = group_quorums[group];
    if (group_vote_counts[group] != quorum) {
        break; // Haven't reached quorum yet
    }
    
    if (group == 0) {
        break; // Reached root group
    }
    
    // Move to parent group
    group = group_parents[group];
}

// Final check: root group must have met its quorum
assert!(group_vote_counts[0] >= group_quorums[0]);
```

### Practical Examples

#### Example 1: Simple Team (3-of-5)
```bash
# Setup
group_quorums = [3, 0, 0, ...]  # Need 3 signatures in group 0
signer_groups = [0, 0, 0, 0, 0] # All 5 signers in group 0

# Signing scenario: Alice, Bob, Charlie sign
# Result: group_vote_counts[0] = 3 >= 3 ✅ APPROVED
```

#### Example 2: Department Structure (2-of-3 departments)
```bash
# Setup  
group_quorums = [2, 1, 1, 1, 0...] # Root needs 2, each dept needs 1
signer_groups = [1, 1, 2, 2, 3]    # 2 eng, 2 sec, 1 exec
group_parents = [0, 0, 0, 0, 1...] # All report to root

# Signing scenario: eng1 + sec1 + exec sign
# eng1 signs → group_vote_counts[1] = 1 (meets quorum) → group_vote_counts[0] = 1
# sec1 signs → group_vote_counts[2] = 1 (meets quorum) → group_vote_counts[0] = 2  
# exec signs → group_vote_counts[3] = 1 (meets quorum) → group_vote_counts[0] = 3
# Result: group_vote_counts[0] = 3 >= 2 ✅ APPROVED
```

#### Example 3: Insufficient Signatures
```bash
# Same setup as Example 2, but only engineering signs
# eng1 + eng2 sign
# eng1 signs → group_vote_counts[1] = 1 (meets quorum) → group_vote_counts[0] = 1
# eng2 signs → group_vote_counts[1] = 2 (already met, no effect on parent)
# Result: group_vote_counts[0] = 1 < 2 ❌ REJECTED
```

### Security Features

#### Replay Protection
```move
// Each signed hash can only be used once
assert!(!multisig.seen_signed_hashes.contains(&signed_hash), EAlreadySeenHash);
multisig.seen_signed_hashes.insert(signed_hash, true);
```

#### Address Ordering
```move
// Signatures must be in ascending order by signer address
assert!(params::vector_u8_gt(&signer_addr, &prev_address), ESignerAddrMustBeIncreasing);
```

#### Group Validation
```move
// Groups must form a valid tree structure
assert!(i == 0 || group_parent < i, EGroupTreeNotWellFormed);  // Parent has lower index
assert!(i != 0 || group_parent == 0, EGroupTreeNotWellFormed); // Root is self-parent
```

### Configuration Constraints

- **Maximum 200 signers** across all groups
- **32 groups maximum** (indices 0-31)  
- **Group 0 must be the root** and its own parent
- **Parent groups must have lower indices** than child groups
- **Quorum = 0 disables a group** (cannot have children)
- **Signer addresses must be unique** and 20 bytes (Ethereum format)

This hierarchical system enables sophisticated governance models while maintaining cryptographic security and preventing various attack vectors through careful validation of the group tree structure and signature ordering requirements.

## Summary

MCMS provides secure, auditable governance for Sui Move contracts through:

- **Multi-signature approval** with flexible quorum requirements
- **Time-delayed execution** for additional security
- **Atomic operations** ensuring all-or-nothing execution  
- **Emergency bypass capabilities** for critical situations
- **Merkle tree batching** for efficient bulk operations

The system is designed to be both secure and practical for managing production smart contracts across multiple protocols and packages.