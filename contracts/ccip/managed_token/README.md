# Managed Token - Standalone Token Management System

This package provides a comprehensive, standalone managed token system with advanced administrative capabilities. **The managed token operates completely independently** and does not require any additional components, token pools, or CCIP infrastructure.

## Key Design Principle

**Complete Independence**: The managed token is a fully self-contained system that can be deployed, managed, and used without requiring any external dependencies beyond the standard Sui framework. It provides all essential token management capabilities as a standalone solution.

## Core Capabilities

The managed token provides comprehensive token management features:

### 🔐 **Treasury Management**
- **Treasury Cap Storage**: Securely stores the treasury cap within the token state
- **Access Control**: Provides controlled access to treasury cap for authorized operations
- **Mint Cap Issuance**: Issues mint capabilities to authorized parties with configurable allowances

### 💰 **Minting & Burning**
- **Controlled Minting**: Mint tokens through authorized mint caps with allowance limits
- **Flexible Burning**: Burn tokens to reduce total supply
- **Allowance Management**: Configure unlimited or limited minting allowances
- **Multiple Minters**: Support for multiple authorized minting parties

### 🚫 **Access Control & Security**
- **Address Blocklisting**: Block specific addresses from all token operations
- **Global Pause**: Pause all token transfers system-wide
- **Deny Cap Integration**: Optional integration with Sui's deny list system
- **Owner-Only Functions**: Critical functions restricted to token owner

### 👑 **Ownership Management**
- **Secure Ownership Transfer**: Two-step ownership transfer process
- **Pending Transfer Tracking**: Track and manage ownership transfer states
- **MCMS Integration**: Optional multi-signature governance support

## Use Cases

The managed token is perfect for various scenarios **without requiring any additional components**:

### Traditional Token Management
- Corporate treasury tokens
- Stablecoin implementations  
- Reward token systems
- Governance tokens with advanced controls

### DeFi Applications
- Lending protocol tokens
- Liquidity mining rewards
- Yield farming tokens
- Protocol governance tokens

### Enterprise Solutions
- Supply chain tokens
- Loyalty point systems
- Internal company tokens
- Partnership reward systems

### Future CCIP Integration (Optional)
- While completely standalone, the managed token can optionally be integrated with CCIP token pools for cross-chain functionality
- See `../ccip_token_pools/managed_token_pool/README.md` for cross-chain integration details

## Quick Start

### Basic Deployment

```move
// Simple deployment without deny capabilities
managed_token::initialize<MyToken>(treasury_cap, ctx);

// This creates:
// - TokenState<MyToken> (shared object containing treasury cap)
// - OwnerCap<MyToken> (transferred to deployer for admin functions)
```

### Full-Featured Deployment

```move
// Deployment with blocklist and pause capabilities
managed_token::initialize_with_deny_cap<MyToken>(
    treasury_cap,
    deny_cap,    // Enables blocklist and pause functionality
    ctx,
);
```

### Setting Up Authorized Minters

```move
// Create mint capabilities for other parties
managed_token::configure_new_minter<MyToken>(
    &mut token_state,
    &owner_cap,
    minter_address,
    1000000,  // Initial allowance (0 if unlimited)
    false,    // is_unlimited (true for unlimited minting)
    ctx,
);

// The MintCap is automatically transferred to minter_address
```

### Basic Token Operations

```move
// Mint tokens (by authorized minter)
let coin = managed_token::mint<MyToken>(
    &mut token_state,
    &mint_cap,
    &deny_list,
    1000,           // Amount to mint
    recipient,      // Who receives the tokens
    ctx,
);

// Burn tokens (by authorized minter)
managed_token::burn<MyToken>(
    &mut token_state,
    &mint_cap,
    &deny_list,
    coin,           // Coin object to burn
    from_address,   // Original owner (for events)
    ctx,
);
```

## Core Functions Reference

### Initialization Functions

| Function | Description | Use Case |
|----------|-------------|----------|
| `initialize<T>()` | Basic token initialization | Simple token deployments |
| `initialize_with_deny_cap<T>()` | Initialize with blocklist/pause support | Enterprise/regulated tokens |

### Minter Management

| Function | Description | Parameters |
|----------|-------------|------------|
| `configure_new_minter<T>()` | Create new mint capability | minter, allowance, is_unlimited |
| `increment_mint_allowance<T>()` | Increase minter's allowance | mint_cap_id, increment_amount |
| `set_unlimited_mint_allowances<T>()` | Toggle unlimited minting | mint_cap_id, is_unlimited |

### Token Operations

| Function | Description | Returns |
|----------|-------------|---------|
| `mint<T>()` | Mint tokens and return coin | `Coin<T>` |
| `mint_and_transfer<T>()` | Mint and transfer directly | void |
| `burn<T>()` | Burn token coin | void |

### Administrative Functions

| Function | Description | Requirements |
|----------|-------------|--------------|
| `blocklist<T>()` | Block an address | Owner cap, deny cap |
| `unblocklist<T>()` | Unblock an address | Owner cap, deny cap |
| `pause<T>()` | Pause all transfers | Owner cap, deny cap |
| `unpause<T>()` | Resume all transfers | Owner cap, deny cap |

### Query Functions

| Function | Description | Returns |
|----------|-------------|---------|
| `total_supply<T>()` | Get total token supply | u64 |
| `mint_allowance<T>()` | Check mint cap allowance | (u64, bool) |
| `is_authorized_mint_cap<T>()` | Verify mint cap authorization | bool |
| `get_all_mint_caps<T>()` | List all mint cap IDs | vector<ID> |

### Access Functions

| Function | Description | Purpose |
|----------|-------------|---------|
| `borrow_treasury_cap<T>()` | Get treasury cap reference | External integrations |

## Advanced Features

### Ownership Management

```move
// Transfer ownership (two-step process)
managed_token::transfer_ownership<MyToken>(&mut state, &owner_cap, new_owner, ctx);

// New owner accepts ownership
managed_token::accept_ownership<MyToken>(&mut state, ctx);

// Check ownership state
let current_owner = managed_token::owner<MyToken>(&state);
let has_pending = managed_token::has_pending_transfer<MyToken>(&state);
```

### Security Features

```move
// Emergency pause (stops all transfers)
managed_token::pause<MyToken>(&mut state, &owner_cap, &mut deny_list, ctx);

// Block problematic addresses
managed_token::blocklist<MyToken>(&mut state, &owner_cap, &mut deny_list, bad_address, ctx);

// Incremental allowance increases (no unlimited conversion)
managed_token::increment_mint_allowance<MyToken>(
    &mut state, &owner_cap, mint_cap_id, &deny_list, 500000, ctx
);
```

### Multi-Signature Governance (Optional)

```move
// Register with MCMS for multi-sig governance
managed_token::mcms_register_entrypoint<MyToken>(
    &mut registry,
    &mut state, 
    owner_cap,  // Transferred to MCMS
    ctx,
);
```

## Error Handling

Common error codes and their meanings:

| Error Code | Constant | Description |
|------------|----------|-------------|
| 1 | `EDeniedAddress` | Address is blocklisted |
| 2 | `EDenyCapNotFound` | Deny cap required but not found |
| 3 | `EInsufficientAllowance` | Mint cap has insufficient allowance |
| 4 | `EInvalidOwnerCap` | Invalid or unauthorized owner cap |
| 5 | `EPaused` | Token is globally paused |
| 6 | `EUnauthorizedMintCap` | Mint cap is not authorized |
| 7 | `EZeroAmount` | Amount must be greater than zero |

## Testing

Run the comprehensive test suite:

```bash
# Navigate to managed token directory
cd contracts/ccip/managed_token

# Run all tests
sui move test -d

# Run specific test categories
sui move test -d -f mint_functionality
sui move test -d -f ownership_management
sui move test -d -f security_features
```

### Test Coverage (22 Tests)

The package includes comprehensive tests covering:
- **Initialization**: Basic and deny cap initialization
- **Minting Operations**: All minting patterns and error cases
- **Burning Operations**: Token burning and validation
- **Access Control**: Blocklist and pause functionality  
- **Ownership**: Transfer and acceptance workflows
- **Allowance Management**: Increment and unlimited settings
- **Security**: Error conditions and unauthorized access
- **Edge Cases**: Boundary conditions and state validation

## Best Practices

### Security
1. **Owner Cap Security**: Store owner cap securely - it controls all administrative functions
2. **Mint Cap Distribution**: Only distribute mint caps to trusted parties
3. **Conservative Allowances**: Start with conservative mint allowances, increase as needed
4. **Regular Monitoring**: Monitor mint cap usage and token supply

### Operational  
1. **Phased Deployment**: Test thoroughly on devnet/testnet before mainnet
2. **Allowance Management**: Regularly review and adjust mint allowances
3. **Event Monitoring**: Set up monitoring for all token events
4. **Backup Procedures**: Have procedures for emergency pause/blocklist

### Integration
1. **Treasury Cap Access**: Use `borrow_treasury_cap()` for secure external integrations
2. **Coin vs Transfer**: Use `mint()` for programmatic use, `mint_and_transfer()` for direct transfers
3. **State Management**: Token state is a shared object - plan concurrent access carefully

## Architecture Benefits

- **🔧 Complete Independence**: No external dependencies beyond Sui framework
- **🎯 Single Responsibility**: Focused solely on token management
- **🔄 Flexible Integration**: Can integrate with various systems without modification
- **⚡ High Performance**: Optimized for efficient token operations
- **🛡️ Enterprise Ready**: Comprehensive security and administrative features
- **📈 Scalable**: Supports unlimited minters and complex allowance structures

## Example Implementation

Here's a complete example of deploying and using a managed token:

```move
module my_company::company_token {
    use sui::coin::{Self, TreasuryCap, DenyCapV2};
    use managed_token::managed_token;
    
    // One-time witness for token creation
    public struct COMPANY_TOKEN has drop {}
    
    // Initialize company token with full features
    fun init(witness: COMPANY_TOKEN, ctx: &mut TxContext) {
        // Create coin and get treasury + deny caps
        let (treasury_cap, deny_cap, metadata) = coin::create_regulated_currency_v2(
            witness,
            9,                              // decimals
            b"COMP",                       // symbol
            b"Company Token",              // name  
            b"Internal company token",     // description
            option::none(),               // icon url
            false,                        // deny by default
            ctx
        );
        
        // Deploy managed token with full capabilities
        managed_token::initialize_with_deny_cap<COMPANY_TOKEN>(
            treasury_cap,
            deny_cap,
            ctx
        );
        
        // Transfer metadata to sender for reference
        transfer::public_transfer(metadata, ctx.sender());
    }
    
    // Company can create authorized minters as needed
    public entry fun create_minter(
        state: &mut managed_token::TokenState<COMPANY_TOKEN>,
        owner_cap: &managed_token::OwnerCap<COMPANY_TOKEN>,
        minter_address: address,
        initial_allowance: u64,
        ctx: &mut TxContext
    ) {
        managed_token::configure_new_minter<COMPANY_TOKEN>(
            state,
            owner_cap,
            minter_address,
            initial_allowance,
            false, // not unlimited
            ctx
        );
    }
}
```

---

**Note**: This managed token system is completely standalone. For cross-chain functionality, see the optional `managed_token_pool` integration guide at `../ccip_token_pools/managed_token_pool/README.md`. 