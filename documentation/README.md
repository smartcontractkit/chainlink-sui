# Chainlink SUI - CCIP Integration

![Chainlink SUI](https://img.shields.io/badge/Chainlink-SUI-blue)
![CCIP](https://img.shields.io/badge/CCIP-Cross%20Chain-green)
![Sui](https://img.shields.io/badge/Sui-Blockchain-orange)
![Go](https://img.shields.io/badge/Go-1.23+-blue)

Welcome to the comprehensive documentation for **Chainlink SUI**, a robust implementation of Chainlink's Cross-Chain Interoperability Protocol (CCIP) for the Sui blockchain. This project enables secure, reliable, and efficient cross-chain communication and token transfers between Sui and other supported blockchain networks.

## 🚀 What is Chainlink SUI?

Chainlink SUI is a complete integration solution that brings Chainlink's battle-tested CCIP infrastructure to the Sui blockchain ecosystem. It provides:

- **CCIP Contracts**: SUI Smart contracts for cross-chain functionality
- **Relayer Service**: blockchain communication and event indexing service
- **Binding generation**: Go bindings for all contracts
- **Operations**: Deployment and management of all contracts

## 📋 Quick Navigation

### 🏃‍♂️ **New to Chainlink SUI?**
- [Quick Start Guide](getting-started/quick-start.md) - Get up and running in minutes
- [Installation](getting-started/installation.md) - Step-by-step setup instructions
- [Development Environment](getting-started/development-environment.md) - Set up your dev environment

### 🔧 **Core Components**
- [Relayer](relayer/introduction.md) - Blockchain communication and event indexing service
- [CCIP Contracts](contracts/overview.md) - Smart contracts for cross-chain functionality
- [Bindings](bindings/overview.md) - Go bindings for all contracts
- [Operations](ops/deployment.md) - Deployment and management of all contracts

### 🚀 **Integration & Deployment**
- [CCIP Integration](integration/ccip.md) - How to integrate CCIP into your application
- [Deployment Guide](ops/deployment.md) - Production deployment instructions
- [Go Bindings](bindings/overview.md) - Using generated Go bindings
- [Operations](ops/deployment.md) - Deployment and management of all contracts



## 🛠️ Technology Stack

| Component | Technology | Purpose |
|-----------|------------|---------|
| **Blockchain** | [Sui](https://sui.io/) | Fast, secure, object-centric blockchain |
| **Smart Contracts** | [Move](https://move-language.github.io/move/) | Safe, resource-oriented programming |
| **Backend Services** | [Go](https://golang.org/) | Relayer, transaction manager, bindings |
| **Cross-Chain Protocol** | [CCIP](https://chain.link/ccip) | Secure cross-chain communication |
| **Database** | PostgreSQL | Event indexing and state management |

## 🏃‍♂️ Quick Start

Get started with Chainlink SUI in just a few steps:

<!-- tabs:start -->

#### **Using Nix (Recommended)**

```bash
# Clone the repository
git clone https://github.com/smartcontractkit/chainlink-sui.git
cd chainlink-sui

# Enter the development environment
nix develop

# Start local Sui network
docker compose up -d sui

# Deploy sample contracts
./scripts/deploy_sample_contracts.sh
```

#### **Manual Setup**

```bash
# Install dependencies
go mod download

# Install Sui CLI
cargo install --git https://github.com/MystenLabs/sui.git --branch main sui

# Set up local environment
sui client new-env --alias local --rpc http://127.0.0.1:9000
sui client switch --env local
```

#### **Using Docker**

```bash
# Quick start with Docker Compose
docker compose up -d

# Deploy contracts
docker compose exec sui bash -c "./scripts/deploy_sample_contracts.sh"
```

<!-- tabs:end -->

## 🔧 System Requirements

### Minimum Requirements
- **Go**: 1.23 or later
- **Sui CLI**: Latest stable version
- **Node.js**: 18+ (for some tooling)
- **Docker**: For local development
- **PostgreSQL**: 12+ (for event indexing)
