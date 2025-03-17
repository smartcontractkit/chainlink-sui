# Deployment Scripts

This directory contains scripts for deploying and interacting with the Sui contracts.

## Available Scripts

### `deploy_contracts.sh`

This script deploys the sample contracts to a local Sui network.

#### Prerequisites

- A local Sui network must be running
- The `sui` CLI must be installed and configured
- `jq` must be installed for JSON parsing

#### Usage

```bash
./deploy_contracts.sh
```

The script will:
1. Build the contracts from the `contracts/test` directory
2. Publish the package to the local Sui network
3. Create a Counter object by calling the `initialize` function
4. Display the package ID and example commands for interacting with the contracts

#### Example Output

```
Starting deployment of Sui sample contracts...
Building contracts...
Package path: /path/to/package.json
Deploying contracts...
Contracts deployed successfully!
Package ID: 0x1234567890abcdef1234567890abcdef
Creating a counter object...
Deployment complete!
You can now interact with the deployed contracts.
Example commands:
  sui client call --package 0x1234567890abcdef1234567890abcdef --module counter --function increment --args $COUNTER_ID --gas-budget 10000000
  sui client call --package 0x1234567890abcdef1234567890abcdef --module counter --function increment_mult --args $COUNTER_ID 5 10 --gas-budget 10000000
```

## Adding New Scripts

When adding new scripts to this directory:
1. Make sure they are executable (`chmod +x script_name.sh`)
2. Update this README with information about the script
3. Follow the same formatting and error handling patterns as existing scripts 