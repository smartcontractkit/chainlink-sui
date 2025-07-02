# Managed Token Pool - CCIP Integration

This package provides CCIP (Cross-Chain Interoperability Protocol) integration for managed tokens through the `managed_token_pool` module.

## Overview

The managed token pool enables cross-chain transfers for managed tokens by:
- **Registering with Token Admin Registry**: Integrates with CCIP's token registration system
- **Handling Cross-Chain Operations**: Manages lock/burn on source chain and release/mint on destination
- **Rate Limiting**: Configurable rate limits for different chains
- **Access Control**: Pool-specific ownership and administration

## Key Design Principle

**Complete Independence**: The managed token pool works with existing managed tokens without requiring any modifications to the managed token contract. The managed token can exist and operate completely independently.

## Prerequisites

Before deploying a managed token pool, you need:
1. **Deployed Managed Token**: A functional managed token contract (see `../managed_token/README.md`)
2. **Mint Cap**: A mint cap created by the managed token owner for the pool
3. **CCIP Environment**: Deployed CCIP infrastructure including token admin registry

## Deployment Process

### Step 1: Prepare Managed Token

If you don't have a managed token, deploy one first:

```move
// Deploy managed token independently (see managed_token documentation)
managed_token::initialize<MyToken>(treasury_cap, ctx);
```

### Step 2: Create Mint Cap for Token Pool

The managed token owner creates a mint cap for the token pool:

```move
// Create mint cap for the token pool
managed_token::configure_new_minter<MyToken>(
    &mut managed_token_state,
    &owner_cap,
    token_pool_administrator, // Who will control the pool
    1000000,                 // Initial allowance
    false,                   // is_unlimited (set to true for unlimited minting)
    ctx,
);

// This transfers the MintCap to the token_pool_administrator
```

### Step 3: Deploy Token Pool

Deploy the token pool using the recommended function:

```move
// Deploy token pool for existing managed token
managed_token_pool::initialize_with_managed_token<MyToken>(
    ref,                     // CCIP object ref
    &managed_token_state,    // Managed token state (immutable reference)
    &owner_cap,             // Managed token owner cap (for treasury cap access)
    coin_metadata,          // Token metadata
    mint_cap,               // Mint cap from step 2
    token_pool_package_id,  // Token pool package ID
    token_pool_administrator, // Pool administrator
    ctx,
);
```

This function:
- Takes an immutable reference to managed token state
- Accesses the treasury cap internally for registration purposes
- Registers the pool with the token admin registry
- Creates the pool state with proper ownership

## Alternative Initialization Methods

### Standard Initialization (Requires Treasury Cap Reference)

```move
managed_token_pool::initialize<MyToken>(
    ref,
    treasury_cap_ref,
    coin_metadata,
    mint_cap,
    token_pool_package_id,
    token_pool_administrator,
    ctx,
);
```

### CCIP Admin Initialization

```move
managed_token_pool::initialize_by_ccip_admin<MyToken>(
    ref,
    coin_metadata,
    mint_cap,
    token_pool_package_id,
    token_pool_administrator,
    ctx,
);
```

## Pool Management Functions

### Chain Configuration

```move
// Add/remove supported chains
managed_token_pool::apply_chain_updates(
    &mut pool_state,
    &owner_cap,
    chains_to_remove,
    chains_to_add,
    remote_pool_addresses,
    remote_token_addresses
);

// Add individual remote pool
managed_token_pool::add_remote_pool(
    &mut pool_state,
    &owner_cap,
    remote_chain_selector,
    remote_pool_address
);
```

### Rate Limiting

```move
// Configure rate limits for a chain
managed_token_pool::set_chain_rate_limiter_config(
    &mut pool_state,
    &owner_cap,
    &clock,
    remote_chain_selector,
    outbound_enabled, outbound_capacity, outbound_rate,
    inbound_enabled, inbound_capacity, inbound_rate
);
```

### Allowlist Management

```move
// Enable/disable allowlist
managed_token_pool::set_allowlist_enabled(&mut pool_state, &owner_cap, enabled);

// Update allowlist
managed_token_pool::apply_allowlist_updates(
    &mut pool_state,
    &owner_cap,
    addresses_to_remove,
    addresses_to_add
);
```

## CCIP Operations

### Cross-Chain Transfer Functions

```move
// Handle outbound transfers (lock or burn)
managed_token_pool::lock_or_burn<MyToken>(
    ref,
    &clock,
    &mut pool_state,
    &deny_list,
    &mut token_state,
    coin,
    token_params,
    ctx
);

// Handle inbound transfers (release or mint)
managed_token_pool::release_or_mint<MyToken>(
    ref,
    &clock,
    &mut pool_state,
    &mut token_state,
    &deny_list,
    receiver_params,
    index,
    ctx
);
```

## Pool Information Functions

```move
// Get basic pool information
let token_address = managed_token_pool::get_token(&pool_state);
let decimals = managed_token_pool::get_token_decimals(&pool_state);
let supported_chains = managed_token_pool::get_supported_chains(&pool_state);

// Get remote chain information
let remote_pools = managed_token_pool::get_remote_pools(&pool_state, chain_selector);
let remote_token = managed_token_pool::get_remote_token(&pool_state, chain_selector);
let is_supported = managed_token_pool::is_supported_chain(&pool_state, chain_selector);

// Get allowlist information
let allowlist_enabled = managed_token_pool::get_allowlist_enabled(&pool_state);
let allowlist = managed_token_pool::get_allowlist(&pool_state);
```

## Pool Lifecycle Management

### Ownership Transfer

```move
// Transfer pool ownership
managed_token_pool::transfer_ownership(&mut pool_state, &owner_cap, new_owner, ctx);
managed_token_pool::accept_ownership(&mut pool_state, ctx);
```

### Pool Destruction

```move
// Destroy pool and recover mint cap (must unregister from token admin registry first)
let recovered_mint_cap = managed_token_pool::destroy_token_pool(pool_state, owner_cap, ctx);
```

## Testing

Run the comprehensive test suite:

```bash
# Test managed token pool package
sui move test -d

# Run specific test for new initialization function
sui move test -d test_initialize_with_managed_token_function
```

### Test Coverage

- `test_initialize_with_managed_token_function()`: Tests the recommended initialization pattern
- `test_lock_or_burn_functionality()`: Cross-chain transfer operations
- `test_release_or_mint_functionality()`: Inbound transfer handling
- `test_rate_limiter_configuration()`: Rate limiting setup
- `test_chain_configuration_management()`: Remote chain management
- And 8 additional comprehensive tests covering edge cases and error scenarios

## Best Practices

### Security
1. **Rate Limiting**: Always configure appropriate rate limits for production
2. **Access Control**: Carefully manage pool owner cap
3. **Allowlists**: Use allowlists for additional transfer restrictions when needed
4. **Monitoring**: Monitor cross-chain operations and rate limit usage

### Operational
1. **Independent Deployment**: Deploy managed token first, then token pool
2. **Mint Cap Management**: Monitor mint cap allowances and adjust as needed
3. **Chain Management**: Add chains incrementally and test thoroughly
4. **Upgrade Planning**: Pools can be upgraded independently of managed tokens

### Integration
1. **CCIP Dependencies**: Ensure all CCIP infrastructure is properly deployed
2. **Token Admin Registry**: Verify registration before enabling transfers
3. **Remote Chains**: Coordinate with remote chain token pool deployments
4. **Testing**: Test on testnets before mainnet deployment

## Error Handling

Common pool-specific errors:

- `EInvalidArguments`: Check array lengths in bulk operations
- `EInvalidOwnerCap`: Ensure correct pool owner cap usage
- Rate limit errors: Check rate limit configurations
- CCIP integration errors: Verify CCIP infrastructure setup

## Architecture Benefits

- **Decoupled Design**: Pool operates independently of managed token core logic
- **Flexible Integration**: Same managed token can work with multiple pool versions
- **Upgrade Path**: Pools can be upgraded without affecting managed token
- **Separation of Concerns**: CCIP logic contained within pool module

## Related Documentation

- **Managed Token**: `../managed_token/README.md` - Core token functionality
- **CCIP Documentation**: Chainlink CCIP documentation for cross-chain concepts
- **Token Admin Registry**: Documentation on token registration and management

---

**Note**: This pool requires an existing managed token. The managed token operates completely independently and can be used without any token pool. 