# Lock Release Token Pool

A CCIP token pool implementation that locks tokens on the source chain and releases them on the destination chain. This pool maintains a reserve of tokens that can be locked for outbound transfers and released for inbound transfers.

## Overview

The Lock Release Token Pool is designed for tokens that exist on multiple chains and need to be transferred cross-chain while maintaining a 1:1 ratio. When tokens are sent cross-chain, they are locked in the pool's reserve. When tokens are received from other chains, they are released from the reserve to the recipient.

**⚠️ Important**: This token pool requires active liquidity management to function properly. A rebalancer must monitor and manage the pool's liquidity to ensure there are sufficient tokens in the reserve for outbound transfers and to handle inbound transfers appropriately.

## Key Features

- ✅ **Token Locking**: Locks tokens in reserve for outbound cross-chain transfers
- ✅ **Token Release**: Releases tokens from reserve for inbound cross-chain transfers
- ✅ **Liquidity Management**: Rebalancer can add/withdraw liquidity to maintain pool balance (required for proper operation)
- ✅ **Multi-Chain Support**: Supports multiple remote chains with configurable token addresses
- ✅ **Rate Limiting**: Per-chain rate limiting to prevent abuse
- ✅ **Access Control**: Owner-only operations with MCMS integration
- ✅ **Allowlist Support**: Optional sender allowlist for restricted access
- ✅ **RMN Integration**: Respects Risk Management Network curse status

## Architecture

### State Structure

```move
public struct LockReleaseTokenPoolState<phantom T> has key {
    id: UID,
    token_pool_state: TokenPoolState,  // Base token pool functionality
    reserve: Coin<T>,                  // Token reserve for locking/releasing
    rebalancer: address,               // Address that can manage liquidity
    ownable_state: OwnableState,       // Ownership management
}
```

### Core Operations

1. **Lock Operation**: When tokens are sent cross-chain, they are locked in the reserve
2. **Release Operation**: When tokens are received from other chains, they are released from the reserve
3. **Liquidity Management**: Rebalancer can add/withdraw liquidity to maintain pool balance

## Usage

### Initialization

```move
use lock_release_token_pool::lock_release_token_pool;

// Initialize the lock release token pool
lock_release_token_pool::initialize<MyToken>(
    &mut owner_cap,
    &mut ccip_ref,
    &coin_metadata,
    &treasury_cap,
    token_pool_administrator,
    rebalancer_address,
    ctx
);
```

### Cross-Chain Operations

```move
// Lock tokens for outbound transfer
lock_release_token_pool::lock_or_burn<MyToken>(
    &mut ccip_ref,
    &clock,
    &mut pool_state,
    sender,
    destination_chain,
    amount,
    ctx
);

// Release tokens for inbound transfer
lock_release_token_pool::release_or_mint<MyToken>(
    &mut ccip_ref,
    &clock,
    &mut pool_state,
    recipient,
    source_chain,
    amount,
    ctx
);
```

### Liquidity Management

#### Direct Operations (Rebalancer)

```move
// Add liquidity to the pool (rebalancer cap required)
lock_release_token_pool::provide_liquidity<MyToken>(
    &mut pool_state,
    &rebalancer_cap,
    liquidity_coin,
    ctx
);

// Withdraw liquidity from the pool (rebalancer cap required)
let withdrawn_coin = lock_release_token_pool::withdraw_liquidity<MyToken>(
    &mut pool_state,
    &rebalancer_cap,
    amount,
    ctx
);
```

#### MCMS Operations (Multi-Sig Governance)

For enhanced security, liquidity operations can be controlled via MCMS:

```move
// Add liquidity via MCMS (requires RebalancerCap in McmsCap)
lock_release_token_pool::mcms_provide_liquidity<MyToken>(
    &mut pool_state,
    &mut registry,
    liquidity_coin,
    params,
    ctx
);

// Withdraw liquidity via MCMS (requires RebalancerCap in McmsCap)
lock_release_token_pool::mcms_withdraw_liquidity<MyToken>(
    &mut pool_state,
    &mut registry,
    params,
    ctx
);
```

**MCMS Liquidity Benefits:**

- **Multi-Sig Security**: Requires multiple signers to approve liquidity changes
- **Time-Delayed Execution**: Operations go through timelock for additional safety

## Core Functions Reference

### Initialization Functions

| Function          | Description                       | Parameters                                                               |
| ----------------- | --------------------------------- | ------------------------------------------------------------------------ |
| `initialize<T>()` | Initialize pool with treasury cap | owner_cap, ccip_ref, coin_metadata, treasury_cap, admin, rebalancer, ctx |

### Cross-Chain Operations

| Function               | Description                         | Parameters                                            |
| ---------------------- | ----------------------------------- | ----------------------------------------------------- |
| `lock_or_burn<T>()`    | Lock tokens for outbound transfer   | ccip_ref, clock, state, sender, chain, amount, ctx    |
| `release_or_mint<T>()` | Release tokens for inbound transfer | ccip_ref, clock, state, recipient, chain, amount, ctx |

### Liquidity Management

| Function                       | Description                           | Parameters                                   |
| ------------------------------ | ------------------------------------- | -------------------------------------------- |
| `provide_liquidity<T>()`       | Add liquidity to pool (direct)        | state, rebalancer_cap, liquidity_coin, ctx   |
| `withdraw_liquidity<T>()`      | Withdraw liquidity from pool (direct) | state, rebalancer_cap, amount, ctx           |
| `mcms_provide_liquidity<T>()`  | Add liquidity via MCMS                | state, registry, liquidity_coin, params, ctx |
| `mcms_withdraw_liquidity<T>()` | Withdraw liquidity via MCMS           | state, registry, params, ctx                 |
| `set_rebalancer<T>()`          | Set new rebalancer address            | state, owner_cap, new_rebalancer, ctx        |
| `get_rebalancer<T>()`          | Get current rebalancer                | state                                        |
| `get_balance<T>()`             | Get pool reserve balance              | state                                        |

### Chain Management

| Function                    | Description                   | Parameters                                                              |
| --------------------------- | ----------------------------- | ----------------------------------------------------------------------- |
| `apply_chain_updates<T>()`  | Add/remove supported chains   | state, chains_to_remove, chains_to_add, pool_addresses, token_addresses |
| `add_remote_pool<T>()`      | Add remote pool address       | state, owner_cap, chain_selector, pool_address, ctx                     |
| `remove_remote_pool<T>()`   | Remove remote pool address    | state, owner_cap, chain_selector, pool_address, ctx                     |
| `get_supported_chains<T>()` | Get supported chain selectors | state                                                                   |
| `is_supported_chain<T>()`   | Check if chain is supported   | state, chain_selector                                                   |

### Query Functions

| Function                  | Description                     | Returns            |
| ------------------------- | ------------------------------- | ------------------ |
| `get_token<T>()`          | Get token metadata address      | address            |
| `get_token_decimals<T>()` | Get token decimals              | u8                 |
| `get_remote_pools<T>()`   | Get remote pool addresses       | vector<vector<u8>> |
| `is_remote_pool<T>()`     | Check if address is remote pool | bool               |
| `get_remote_token<T>()`   | Get remote token address        | vector<u8>         |

### Allowlist Management

| Function                       | Description                   | Parameters                                                   |
| ------------------------------ | ----------------------------- | ------------------------------------------------------------ |
| `get_allowlist_enabled<T>()`   | Check if allowlist is enabled | state                                                        |
| `set_allowlist_enabled<T>()`   | Enable/disable allowlist      | state, owner_cap, enabled, ctx                               |
| `get_allowlist<T>()`           | Get allowlist addresses       | state                                                        |
| `apply_allowlist_updates<T>()` | Update allowlist              | state, owner_cap, addresses_to_remove, addresses_to_add, ctx |

### Rate Limiting

| Function                              | Description                         | Parameters                                           |
| ------------------------------------- | ----------------------------------- | ---------------------------------------------------- |
| `set_chain_rate_limiter_configs<T>()` | Set rate limits for multiple chains | clock, state, owner_cap, configs, ctx                |
| `set_chain_rate_limiter_config<T>()`  | Set rate limit for single chain     | clock, state, owner_cap, chain_selector, config, ctx |

### Ownership Management

| Function                         | Description                | Parameters                       |
| -------------------------------- | -------------------------- | -------------------------------- |
| `owner<T>()`                     | Get current owner          | state                            |
| `transfer_ownership<T>()`        | Transfer ownership         | state, owner_cap, new_owner, ctx |
| `accept_ownership<T>()`          | Accept ownership transfer  | state, ctx                       |
| `has_pending_transfer<T>()`      | Check for pending transfer | state                            |
| `pending_transfer_from<T>()`     | Get pending transfer from  | state                            |
| `pending_transfer_to<T>()`       | Get pending transfer to    | state                            |
| `pending_transfer_accepted<T>()` | Check if transfer accepted | state                            |

### MCMS Integration Functions

| Function                                   | Description                 | Purpose                         |
| ------------------------------------------ | --------------------------- | ------------------------------- |
| `mcms_accept_ownership<T>()`               | Accept ownership via MCMS   | Multi-sig ownership transfer    |
| `mcms_register_upgrade_cap()`              | Register upgrade capability | Multi-sig upgrade management    |
| `mcms_set_rebalancer<T>()`                 | Set rebalancer via MCMS     | Multi-sig rebalancer management |
| `mcms_provide_liquidity<T>()`              | Add liquidity via MCMS      | Multi-sig liquidity management  |
| `mcms_withdraw_liquidity<T>()`             | Withdraw liquidity via MCMS | Multi-sig liquidity management  |
| `mcms_set_allowlist_enabled<T>()`          | Set allowlist via MCMS      | Multi-sig access control        |
| `mcms_apply_allowlist_updates<T>()`        | Update allowlist via MCMS   | Multi-sig allowlist management  |
| `mcms_apply_chain_updates<T>()`            | Update chains via MCMS      | Multi-sig chain management      |
| `mcms_add_remote_pool<T>()`                | Add remote pool via MCMS    | Multi-sig pool management       |
| `mcms_remove_remote_pool<T>()`             | Remove remote pool via MCMS | Multi-sig pool management       |
| `mcms_set_chain_rate_limiter_configs<T>()` | Set rate limits via MCMS    | Multi-sig rate limiting         |
| `mcms_set_chain_rate_limiter_config<T>()`  | Set rate limit via MCMS     | Multi-sig rate limiting         |
| `mcms_transfer_ownership<T>()`             | Transfer ownership via MCMS | Multi-sig ownership transfer    |
| `mcms_execute_ownership_transfer<T>()`     | Execute transfer via MCMS   | Multi-sig ownership transfer    |

### Utility Functions

| Function                  | Description        | Parameters            |
| ------------------------- | ------------------ | --------------------- |
| `destroy_token_pool<T>()` | Destroy pool state | state, owner_cap, ctx |

## Integration

This module integrates with:

- **CCIPTokenPool**: Base token pool functionality
- **ChainlinkCCIP**: Core CCIP protocol functionality
- **MCMS**: Multi-Chain Multi-Sig governance
- **RMN Remote**: Risk Management Network for security
- **Sui Framework**: Clock, events, and standard library functions

## Security Considerations

- All cross-chain operations validate RMN curse status
- Rate limiting prevents rapid token drainage
- Ownership changes require two-step confirmation
- Remote pool addresses are validated against configured lists
- Rebalancer can only manage liquidity, not cross-chain operations
- **Liquidity Management Required**: Pool requires active rebalancer for proper operation
- Decimal overflow protection prevents arithmetic errors

## Package Information

- **Package Name**: `LockReleaseTokenPool`
- **Version**: 1.6.0
- **Edition**: 2024
- **Dependencies**:
  - `CCIPTokenPool` (Base token pool functionality)
  - `ChainlinkManyChainMultisig` (MCMS integration)

### Module Structure

- `lock_release_token_pool.move` - Core lock/release functionality

## Testing

Run tests with:

```bash
sui move test
```

The test suite covers:

- Initialization and configuration
- Lock and release operations
- Liquidity management
- Chain configuration
- Rate limiting
- Access control
- Ownership management
- MCMS integration
- Error conditions and edge cases
