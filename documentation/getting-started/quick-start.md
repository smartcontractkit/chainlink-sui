# [WIP] Quick Start Guide

Get up and running with Chainlink SUI in just a few minutes! This guide will walk you through setting up your development environment, deploying contracts, and sending your first cross-chain message.

## Prerequisites

Before you begin, make sure you have:

- **Docker** and **Docker Compose** installed
- **Git** for cloning the repository
- **Basic familiarity** with blockchain concepts and command line tools

> ℹ️ **Using Nix?** If you have Nix installed, you can skip most setup steps by using our development shell: `nix develop`

## Step 1: Clone and Setup

```bash
# Clone the repository
git clone https://github.com/smartcontractkit/chainlink-sui.git
cd chainlink-sui

# Start the local Sui network
docker compose up -d sui

# Wait for the network to be ready (about 30 seconds)
docker compose logs -f sui
```

Look for this log message indicating the network is ready:
```
Sui RPC server listening on 0.0.0.0:9000
```

## Step 2: Set Up Sui CLI

<!-- tabs:start -->

#### **Option A: Using Docker (Recommended)**

```bash
# Access the Sui CLI through Docker
docker compose exec sui bash

# Inside the container, check the environment
sui client envs
sui client active-address
```

#### **Option B: Local Installation**

```bash
# Install Sui CLI (requires Rust)
cargo install --git https://github.com/MystenLabs/sui.git --branch main sui

# Configure for local network
sui client new-env --alias local --rpc http://127.0.0.1:9000
sui client switch --env local

# Verify connection
sui client active-address
```

<!-- tabs:end -->

## Step 3: Deploy Sample Contracts

Deploy the test contracts to get familiar with the system:

```bash
# Deploy sample contracts (counter, echo)
./scripts/deploy_sample_contracts.sh

# Expected output:
# Successfully deployed counter contract: 0x123...
# Successfully deployed echo contract: 0x456...
# Counter object created: 0x789...
```

Save these addresses - you'll need them for testing!

## Step 4: Deploy CCIP Infrastructure

Now let's deploy the core CCIP contracts:

```bash
# Deploy CCIP contracts
./scripts/deploy_ccip_contracts.sh

# This will deploy:
# - CCIP Router
# - OnRamp contract
# - OffRamp contract  
# - Token Admin Registry
# - Sample token pools
```

The script will output important contract addresses. Save them for later use:

```bash
# Example output:
CCIP_ROUTER=0xabc123...
ONRAMP=0xdef456...
OFFRAMP=0x789abc...
TOKEN_ADMIN_REGISTRY=0x123def...
```

## Step 5: Configure the Relayer

Create a basic relayer configuration:

```bash
# Create configuration directory
mkdir -p config

# Create relayer configuration
cat > config/relayer.toml << 'EOF'
[[Chains]]
ChainID = '0x1'
Enabled = true

# Basic transaction settings
BroadcastChanSize = 1000
ConfirmPollPeriod = '1s'
MaxConcurrentRequests = 10
TransactionTimeout = '30s'
NumberRetries = 3
GasLimit = 10000000
RequestType = 'WaitForLocalExecution'

# Transaction Manager settings
[TransactionManager]
BroadcastChanSize = 100
ConfirmPollSecs = 2
DefaultMaxGasAmount = 200000
MaxSimulateAttempts = 3
MaxSubmitRetryAttempts = 5
MaxTxRetryAttempts = 3
PruneIntervalSecs = 3600
PruneTxExpirationSecs = 1800
SubmitDelayDuration = 1
TxExpirationSecs = 30

# Sui node configuration
[[Chains.Nodes]]
Name = 'local-sui'
URL = 'http://localhost:9000'
SolidityURL = 'http://localhost:9000'
EOF
```

## Step 6: Test Basic Functionality

Let's test the basic functionality with the deployed contracts:

### Test Counter Contract

```bash
# Get the counter value (should be 0 initially)
sui client call \
  --package $COUNTER_PACKAGE \
  --module counter \
  --function get_count \
  --args $COUNTER_OBJECT \
  --dev-inspect

# Increment the counter
sui client call \
  --package $COUNTER_PACKAGE \
  --module counter \
  --function increment \
  --args $COUNTER_OBJECT \
  --gas-budget 10000000

# Check the new value (should be 1)
sui client call \
  --package $COUNTER_PACKAGE \
  --module counter \
  --function get_count \
  --args $COUNTER_OBJECT \
  --dev-inspect
```

### Test Event Emission

```bash
# Emit a test event
sui client call \
  --package $ECHO_PACKAGE \
  --module echo \
  --function emit_simple_event \
  --args "Hello Chainlink SUI!" \
  --gas-budget 10000000
```

## Step 7: Start the Relayer

Now let's start the Chainlink SUI relayer:

<!-- tabs:start -->

#### **Development Mode**

```bash
# Start PostgreSQL for event indexing
docker compose up -d postgres

# Set environment variables
export CHAINLINK_SUI_CONFIG=config/relayer.toml
export DATABASE_URL=postgresql://postgres:postgres@localhost:5432/postgres

# Build and run the relayer
go build -o bin/chainlink-sui ./relayer/cmd/chainlink-sui
./bin/chainlink-sui
```

#### **Docker Mode**

```bash
# Create Docker configuration
cat > docker-compose.relayer.yml << 'EOF'
services:
  relayer:
    build:
      context: .
      dockerfile: Dockerfile.relayer
    environment:
      - CHAINLINK_SUI_CONFIG=/config/relayer.toml
      - DATABASE_URL=postgresql://postgres:postgres@postgres:5432/postgres
    volumes:
      - ./config:/config
    depends_on:
      - postgres
      - sui
    networks:
      - sui
EOF

# Start the relayer
docker compose -f docker-compose.yml -f docker-compose.relayer.yml up -d relayer
```

<!-- tabs:end -->

## Step 8: Send Your First Cross-Chain Message

Now for the exciting part - let's send a cross-chain message!

### Prepare the Message

```bash
# Create a simple message payload
MESSAGE_DATA="Hello from Sui to Ethereum!"
DESTINATION_CHAIN_SELECTOR="5009297550715157269"  # Ethereum Sepolia
RECEIVER_ADDRESS="0x742d35Cc6634C0532925a3b8D400c45532c83a36"  # Example receiver

# Estimate fees first
sui client call \
  --package $CCIP_ROUTER \
  --module router \
  --function get_fee \
  --args $DESTINATION_CHAIN_SELECTOR $MESSAGE_DATA \
  --dev-inspect
```

### Send the Message

```bash
# Send cross-chain message
sui client call \
  --package $CCIP_ONRAMP \
  --module onramp \
  --function ccip_send \
  --args $DESTINATION_CHAIN_SELECTOR $RECEIVER_ADDRESS $MESSAGE_DATA \
  --gas-budget 20000000

# Expected output will include:
# - Transaction hash
# - Message ID
# - Events emitted
```

### Monitor the Message

```bash
# Check the transaction status
sui client tx-block $TRANSACTION_HASH

# Watch for CCIP events in the logs
docker compose logs -f relayer | grep "CCIPMessageSent"
```

## Step 9: Verify Cross-Chain Delivery

The message will be processed by the off-chain infrastructure and delivered to the destination chain. You can monitor this through:

### Check Message Status

```bash
# Query the message status (this would be on the destination chain)
# For now, check the events emitted on Sui
sui client events --package $CCIP_ONRAMP
```

### View in Explorer

- **Sui Transactions**: Check on [Sui Explorer](https://explorer.sui.io/)
- **Cross-Chain Messages**: Monitor through CCIP explorer tools

## Step 10: Test Token Transfers

Let's also test cross-chain token transfers:

### Deploy a Test Token

```bash
# Deploy a managed token for testing
sui client publish contracts/managed_token --gas-budget 50000000

# Save the package ID
MANAGED_TOKEN_PACKAGE=<package_id_from_output>
```

### Create Token Pool

```bash
# Deploy a burn/mint token pool
sui client publish contracts/ccip_token_pools/burn_mint_token_pool --gas-budget 50000000

# Initialize the token pool
sui client call \
  --package $BURN_MINT_POOL_PACKAGE \
  --module burn_mint_token_pool \
  --function initialize \
  --args $MANAGED_TOKEN_PACKAGE \
  --gas-budget 10000000
```

### Send Token Cross-Chain

```bash
# First, mint some tokens to transfer
AMOUNT=1000000  # 1 token (assuming 6 decimals)

sui client call \
  --package $MANAGED_TOKEN_PACKAGE \
  --module managed_token \
  --function mint \
  --args $AMOUNT \
  --gas-budget 10000000

# Send tokens cross-chain
sui client call \
  --package $CCIP_ONRAMP \
  --module onramp \
  --function ccip_send_with_tokens \
  --args $DESTINATION_CHAIN_SELECTOR $RECEIVER_ADDRESS $MESSAGE_DATA $TOKEN_TRANSFERS \
  --gas-budget 30000000
```

## Next Steps

Congratulations! 🎉 You've successfully:

- ✅ Set up a local Sui development environment
- ✅ Deployed CCIP infrastructure contracts
- ✅ Configured and started the Chainlink SUI relayer
- ✅ Sent your first cross-chain message
- ✅ Tested token transfers

### Where to go from here:

1. **📚 Learn the Architecture**: Read the [System Overview](../architecture/overview.md) to understand how everything works
2. **🔧 Build Custom Solutions**: Check out the [Integration Guides](../integration/ccip.md) to build your own cross-chain applications
3. **🏗️ Deploy to Testnet**: Follow the [Deployment Guide](../ops/deployment.md) to deploy to Sui testnet
4. **🧪 Advanced Features**: Explore [Programmable Transaction Blocks](../ptb/introduction.md) for complex operations

### Useful Resources

- **[Relayer Documentation](../relayer/introduction.md)**: Deep dive into relayer configuration and usage
- **[Contract Documentation](../contracts/overview.md)**: Understand the CCIP contract architecture
- **[Development Guide](../development/contributing.md)**: Contributing to the project
- **[Troubleshooting](../troubleshooting/common-issues.md)**: Solutions to common problems

### Getting Help

- 🐛 **Found a bug?** [Open an issue](https://github.com/smartcontractkit/chainlink-sui/issues)
- 💬 **Questions?** [Join discussions](https://github.com/smartcontractkit/chainlink-sui/discussions)
- 📖 **Need more info?** Check the detailed documentation sections

---

**Happy building with Chainlink SUI!** 🚀 