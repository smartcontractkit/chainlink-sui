# SUI CCIP CLI Example

To install dependencies:

```bash
bun install
```

## Configuration

### Network-Specific Environment Files

The CLI supports network-specific environment files to easily switch between different networks (localnet, testnet, mainnet, devnet).

1. Create your network-specific env files:
   ```bash
   cp .env.example .env.localnet
   cp .env.example .env.testnet
   ```

2. Fill in the values for each network in the respective files:
   - `.env.localnet` - for local development
   - `.env.testnet` - for testnet deployment

3. The CLI will automatically load the correct env file based on the `--network` flag:
   ```bash
   # Uses .env.localnet
   bun ccip_send --dest-chain-selector 2 --receiver 0x... --pool-kind lock_release --network localnet
   
   # Uses .env.testnet
   bun ccip_send --dest-chain-selector 16015286601757825753 --receiver 0x... --pool-kind burn_mint --network testnet
   ```

### Default Environment File

You can also use a default `.env` file that will be used when no `--network` flag is specified or when using a custom fullnode URL.

## Usage

To see all available options:

```bash
bun ccip_send --help
```

The command above will display all the arguments needed to send a CCIP Send message via the onramp package.
