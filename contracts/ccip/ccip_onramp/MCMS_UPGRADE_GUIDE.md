# CCIP OnRamp Upgrade Guide with MCMS

This guide documents how to upgrade CCIP OnRamp packages using Multi-Chain Multi-Sig (MCMS) governance integrated with native Sui package upgrade functionality.

## Overview

### Sui Package Upgrades

Sui package upgrades work through a three-step process:

1. **Authorize**: Create `UpgradeTicket` via `package::authorize_upgrade(cap, policy, digest)`
2. **Upgrade**: Execute upgrade transaction within PTB using `tx.upgrade()`
3. **Commit**: Complete upgrade via `package::commit_upgrade(cap, receipt)`

### MCMS Integration

MCMS provides governance-controlled upgrades through:

- Multi-signature approval requirements
- Timelock delays for security
- Hot potato pattern ensuring atomic execution
- Version compatibility checks

## Prerequisites

Before upgrading, ensure:

1. **Package Deployed with Upgrade Support**

   - `UpgradeCap` registered with MCMS deployer
   - `OwnerCap` registered with MCMS registry
   - Version management implemented in contract

2. **MCMS Configuration**

   - Signers configured with appropriate quorum
   - Timelock delays set as required
   - Registry and deployer state objects deployed

3. **New Package Prepared**
   - Code changes maintain compatibility (public function signatures, struct layouts)
   - Version constants updated if needed
   - Migration functions implemented if required

## Complete Upgrade Flow

### Phase 1: MCMS Proposal & Scheduling

#### Step 1: Create Upgrade Proposal

```bash
# Off-chain: Generate upgrade proposal with merkle tree
# Include operation data for authorize_upgrade call
```

#### Step 2: Sign and Submit Root

```move
// MCMS signers call set_root with signed merkle root
mcms::set_root(
    state: &mut MultisigState,
    clock: &Clock,
    role: u8, // PROPOSER_ROLE
    root: vector<u8>, // Merkle root containing upgrade operation
    valid_until: u64,
    // ... other parameters
    metadata_proof: vector<vector<u8>>,
    signatures: vector<vector<u8>>,
)
```

#### Step 3: Schedule Upgrade Operation

```move
// Execute and schedule via MCMS timelock
let timelock_params = mcms::execute(
    state: &mut MultisigState,
    // ... merkle proof for authorize_upgrade operation
);

mcms::dispatch_timelock_schedule_batch(
    timelock: &mut Timelock,
    clock: &Clock,
    timelock_params, // hot potato from execute
    ctx: &mut TxContext,
);
```

### Phase 2: Execute Upgrade (After Timelock Delay)

#### Step 4: Execute MCMS Authorization

```move
// After timelock delay expires, execute the scheduled operation
let timelock_params = mcms::execute(
    state: &mut MultisigState,
    // ... merkle proof for authorize_upgrade operation
);

let executing_params = mcms::dispatch_timelock_execute_batch(
    timelock: &mut Timelock,
    clock: &Clock,
    timelock_params, // hot potato
    ctx: &mut TxContext,
); // Returns vector<ExecutingCallbackParams>
```

#### Step 5: Complete MCMS-Controlled Upgrade in Single PTB

The upgrade must happen atomically in a single PTB with three commands:

```typescript
// PTB construction (TypeScript SDK example)
const tx = new TransactionBlock();

// 1. MCMS authorizes upgrade and returns UpgradeTicket
const [upgradeTicket] = tx.moveCall({
  target: `${mcmsPackageId}::mcms::execute_dispatch_to_deployer`,
  arguments: [
    tx.object(registryId),
    tx.object(deployerStateId),
    tx.object(executingCallbackParamsId), // from MCMS execution
  ],
});

// 2. Execute the package upgrade using the ticket
const [upgradeReceipt] = tx.upgrade({
  modules: newPackageModules,
  dependencies: packageDependencies,
  packageId: currentPackageId,
  ticket: upgradeTicket, // Hot potato from step 1
});

// 3. Commit upgrade to complete the process
tx.moveCall({
  target: `${mcmsPackageId}::mcms_deployer::commit_upgrade`,
  arguments: [
    tx.object(deployerStateId),
    upgradeReceipt, // Hot potato from step 2
  ],
});

// Execute the PTB
const result = await suiClient.signAndExecuteTransactionBlock({
  transactionBlock: tx,
  signer: keypair,
  options: { showEffects: true, showObjectChanges: true },
});
```

### Phase 3: Post-Upgrade Migration (If Needed)

#### Step 6: Execute State Migration via MCMS

After the package upgrade is complete, shared object state migration must be done through MCMS entrypoint since MCMS owns the `OwnerCap`:

**IMPORTANT:** The migration function must first be added to the `mcms_entrypoint` in the NEW package version:

```move
// In the NEW package's onramp.move mcms_entrypoint function:
} else if (function_bytes == b"migrate_to_v2") {
    bcs_stream::assert_is_consumed(&stream); // No additional parameters
    migrate_to_v2(state, owner_cap);
} else if (function_bytes == b"migrate_to_v3") {
    bcs_stream::assert_is_consumed(&stream); // No additional parameters
    migrate_to_v3(state, owner_cap);
```

#### Migration via MCMS (Production Process)

1. **Create MCMS Proposal for Migration:**

```typescript
// Create MCMS operation for migration
const migrationOp = {
  target: NEW_PACKAGE_ID, // The upgraded package ID
  module_name: "onramp",
  function_name: "migrate_to_v2",
  data: encodeParameters([], []), // No additional data needed
};

// Submit for MCMS approval and timelock
```

2. **Execute Migration through MCMS:**

```bash
# After MCMS approval and timelock delay
sui client call \
  --package NEW_PACKAGE_ID \
  --module onramp \
  --function mcms_entrypoint \
  --args ONRAMP_STATE_ID MCMS_REGISTRY_ID EXECUTING_CALLBACK_PARAMS_ID \
  --gas-budget 10000000
```

**Important Notes:**

- Migration functions belong in the **new package version**, not the old one
- Migration functions must be **added to mcms_entrypoint** to be callable by MCMS
- Each migration function can only be called **once** (version protection)
- The function checks current version: `assert!(state.version == 1, EIncompatibleVersion)`
- Migration is **separate** from the package upgrade itself
- MCMS owns the `OwnerCap`, so direct calls to migration functions will fail

## Technical Implementation Details

### PTB Structure for Upgrades

The complete upgrade must happen atomically in a single PTB with exactly 3 commands:

```move
// PTB structure for MCMS-controlled upgrade:
programmable_transaction_block {
    // 1. MCMS authorizes upgrade and returns UpgradeTicket hot potato
    command_1: mcms::execute_dispatch_to_deployer(
        registry, deployer_state, executing_callback_params
    ) -> UpgradeTicket

    // 2. Execute package upgrade using the UpgradeTicket
    command_2: tx.upgrade(
        modules, dependencies, packageId, ticket
    ) -> UpgradeReceipt

    // 3. Commit upgrade consuming the UpgradeReceipt hot potato
    command_3: mcms_deployer::commit_upgrade(
        deployer_state, receipt
    )

    // All commands must succeed or entire PTB fails atomically
}
```

### Version Compatibility

OnRamp contract includes version management:

```move
const VERSION: u64 = 1;

public struct OnRampState has key, store {
    id: UID,
    version: u64,  // Tracks contract version
    // ... other fields
}

// All public functions check version compatibility
fun assert_compatible_version(state: &OnRampState) {
    assert!(state.version == VERSION, EIncompatibleVersion);
}
```

### Migration Functions

Migration functions are placed in the **NEW package version**, not the old one:

```move
// This function belongs in VERSION 2 package, not VERSION 1
public fun migrate_to_v2(state: &mut OnRampState, _: &OwnerCap) {
    assert!(state.version == 1, EIncompatibleVersion); // Ensures one-time execution

    // Perform migration logic here:
    // - Initialize new fields added in v2
    // - Transform existing data if needed
    // - Update data structures

    state.version = 2; // Update version last
}

// Version 3 would have migrate_to_v3 function
public fun migrate_to_v3(state: &mut OnRampState, _: &OwnerCap) {
    assert!(state.version == 2, EIncompatibleVersion); // From v2 to v3
    // Migration logic for v2 → v3
    state.version = 3;
}
```

**Key Principles:**

- **Each migration function can only be called once** (version assertion protection)
- **Migration functions belong in target version**, not source version
- **Migration functions must be added to mcms_entrypoint** for MCMS to call them
- **Error code 19** indicates version assertion failure (trying to migrate incompatible version)
- **Sequential migrations required**: v1→v2→v3, cannot skip versions

## Example: Complete Upgrade Workflow

### 1. Prepare Upgrade

```bash
# Build new package version
sui move build

# Verify compatibility
sui move test

# Generate package digest
DIGEST=$(sui client publish --dry-run --json | jq -r '.digest')
```

### 2. Create MCMS Proposal

```typescript
// Generate merkle tree with authorize_upgrade operation
const upgradeOperation = {
  target: MCMS_DEPLOYER_ADDRESS,
  module_name: "mcms_deployer",
  function_name: "authorize_upgrade",
  data: encodeParameters(
    ["u8", "vector<u8>", "address"],
    [POLICY_COMPATIBLE, packageDigest, PACKAGE_ID]
  ),
};

// Submit to MCMS signers for approval and scheduling
```

### 3. Execute Upgrade PTB (After Timelock)

```typescript
const tx = new TransactionBlock();

// 1. Get UpgradeTicket from MCMS
const [upgradeTicket] = tx.moveCall({
  target: `${MCMS_PACKAGE}::mcms::execute_dispatch_to_deployer`,
  arguments: [registryId, deployerStateId, executingCallbackParamsId],
});

// 2. Execute upgrade
const [upgradeReceipt] = tx.upgrade({
  modules: compiledModules,
  dependencies: [CCIP_PACKAGE, MCMS_PACKAGE, SUI_FRAMEWORK],
  packageId: CURRENT_PACKAGE_ID,
  ticket: upgradeTicket,
});

// 3. Commit upgrade
tx.moveCall({
  target: `${MCMS_PACKAGE}::mcms_deployer::commit_upgrade`,
  arguments: [deployerStateId, upgradeReceipt],
});

// Execute atomically
const result = await suiClient.signAndExecuteTransactionBlock({
  transactionBlock: tx,
  signer: keypair,
  options: { showEffects: true, showObjectChanges: true },
});
```

### 4. Execute Migration via MCMS (Separate Transaction)

```bash
# Migration must go through MCMS entrypoint since MCMS owns the OwnerCap
sui client call \
  --package NEW_PACKAGE_ID \
  --module onramp \
  --function mcms_entrypoint \
  --args ONRAMP_STATE_ID MCMS_REGISTRY_ID EXECUTING_CALLBACK_PARAMS_ID \
  --gas-budget 10000000
```

**Prerequisites:**

- Migration function added to `mcms_entrypoint` in the new package
- MCMS proposal created and approved for the migration operation
- Timelock delay satisfied (if required)

### 5. Verify Upgrade

```bash
# Check package version
sui client object $PACKAGE_ID

# Verify contract version
sui client call --function get_version --module onramp --package $PACKAGE_ID
```

## Troubleshooting

### Debug Commands

```bash
# Check MCMS state
sui client call --function get_root --module mcms --package $MCMS_PACKAGE

# Verify upgrade cap registration
sui client call --function is_package_registered --module mcms_registry --package $MCMS_PACKAGE

# Check timelock status
sui client call --function timelock_is_operation_ready --module mcms --package $MCMS_PACKAGE
```
