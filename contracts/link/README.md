# LINK Token - ChainLink Token Implementation for Sui

This package provides the ChainLink Token (LINK) implementation for the Sui blockchain, designed to maintain feature parity with the LINK token on Ethereum while leveraging Sui's unique capabilities.

## Overview

The LINK token is a fundamental component of the ChainLink ecosystem, serving as the native utility token for ChainLink's decentralized oracle network. This implementation brings LINK to the Sui blockchain with careful attention to maintaining the same properties and behaviors as the Ethereum version.

## Key Design Principles

### 🔗 **Ethereum Feature Parity**
- **No Deny List**: LINK tokens maintain the same permissionless nature as on Ethereum
- **No Blocklist Capability**: Users cannot be blocked from holding or transferring LINK tokens
- **Open Access**: Preserves the decentralized and permissionless nature of the original LINK token

### 🏗️ **Multi-Minter Architecture**
- **Managed Token Integration**: Designed to work in conjunction with the `managed_token` package
- **Treasury Cap Storage**: The treasury cap is stored in the managed_token for secure multi-minter functionality
- **Authorized Minters**: Owner can issue minter caps to eligible parties (typically LINK token pools)

## Architecture

### Token Structure
```move
public struct LINK has drop {}
```

The LINK token follows Sui's standard coin implementation with these characteristics:
- **Symbol**: LINK
- **Name**: ChainLink Token  
- **Decimals**: 9
- **Icon**: ChainLink official token icon

### Integration with Managed Token

This LINK token implementation is **designed to be used in conjunction with the `managed_token` package** to enable advanced treasury management capabilities:

```
┌─────────────────┐    stores treasury cap    ┌─────────────────┐
│   LINK Token    │ ─────────────────────────▶ │ Managed Token   │
│   Package       │                           │ Package         │
└─────────────────┘                           └─────────────────┘
                                                       │
                                                       │ issues minter caps
                                                       ▼
                                               ┌─────────────────┐
                                               │ LINK Token      │
                                               │ Pools & Other   │
                                               │ Authorized      │
                                               │ Minters         │
                                               └─────────────────┘
```

## Usage Patterns

### Initial Deployment
1. Deploy the LINK token contract
2. Initialize the managed_token with the LINK treasury cap
3. The managed_token stores the treasury cap and enables multi-minter functionality

### Minter Management
The token owner (via the managed_token package) can:
- Issue `MintCap` objects to authorized minters
- Configure allowances for each minter (limited or unlimited)
- Revoke minter capabilities as needed

### Typical Authorized Minters
- **CCIP Token Pools**: For cross-chain LINK transfers
- **Bridge Contracts**: For moving LINK between different networks  
- **Protocol Contracts**: For minting LINK as rewards or operational needs

## Core Functions

### `init(witness: LINK, ctx: &mut TxContext)`
Initializes the LINK token with metadata and creates the treasury cap. The treasury cap should be transferred to a managed_token instance for multi-minter functionality.

### `mint_and_transfer(treasury_cap, amount, recipient, ctx)`
Mints LINK tokens and transfers them directly to a recipient. This function would typically be called by authorized minters through the managed_token system.

### `mint(treasury_cap, amount, ctx): Coin<LINK>`
Mints LINK tokens and returns the coin object. This provides flexibility for minters to handle the tokens as needed.

## Key Differences from Ethereum LINK

### Similarities (Feature Parity)
- ✅ **No Blocklist**: Cannot block users from holding or transferring tokens
- ✅ **Permissionless**: Anyone can hold and transfer LINK tokens
- ✅ **Standard Token**: Follows standard token interface patterns

### Sui-Specific Enhancements  
- ✅ **Object-Based**: Leverages Sui's object model for better composability
- ✅ **Multi-Minter Support**: Built-in support for multiple authorized minters
- ✅ **Secure Treasury Management**: Treasury cap securely managed by managed_token package

## Important Design Decision: No Deny Cap

**This LINK token implementation deliberately does NOT use Sui's deny list functionality.** This design choice ensures:

1. **Ethereum Parity**: Maintains the same permissionless properties as LINK on Ethereum
2. **Decentralization**: Preserves the decentralized nature of the ChainLink ecosystem  
3. **User Freedom**: Users cannot be arbitrarily blocked from using their tokens
4. **Protocol Integrity**: Maintains the trustless nature expected of LINK tokens

## Deployment and Integration

### For Token Pools
When integrating LINK with token pools:
1. Obtain a `MintCap` from the managed_token owner
2. Use the mint capabilities for cross-chain operations
3. Respect allowance limits (if applicable)

### For DeFi Protocols
When integrating LINK into DeFi applications:
1. Use standard Sui coin interfaces
2. No special handling needed for blocklists (since none exist)
3. Treat as a standard, permissionless token

## Security Considerations

- **Treasury Cap Security**: The treasury cap is stored in the managed_token, providing secure multi-minter functionality
- **Minter Authorization**: Only authorized parties can receive mint capabilities
- **No Arbitrary Blocking**: Cannot block legitimate users, maintaining decentralization
- **Standard Coin Security**: Inherits security properties of Sui's coin framework

## Compatibility

This LINK token is designed to be fully compatible with:
- Standard Sui coin interfaces
- CCIP token pools and cross-chain infrastructure
- DeFi protocols expecting standard coin behavior
- Wallet applications and user interfaces

The absence of deny list functionality ensures maximum compatibility and maintains the expected behavior of LINK tokens across all applications. 