# Mock Contracts CLI Tool

A Go CLI tool for managing Sui Move contracts from the `contracts/test/sources` directory. This tool provides environment setup, contract deployment, and event emission capabilities.

## Prerequisites

1. **Sui CLI**: Must be installed and available in PATH
   ```bash
   # Install Sui CLI
   cargo install --locked --git https://github.com/MystenLabs/sui.git --branch devnet sui
   ```

2. **Sui Account**: For deployment, you need an active Sui account with SUI tokens
   ```bash
   # Create new account (if needed)
   sui client new-address ed25519
   
   # Check current account
   sui client active-address
   
   # Check balance
   sui client gas
   ```

3. **Network Access**: 
   - For local deployment: Local Sui node running
   - For devnet/testnet: Internet connection

## Usage

The tool provides multiple commands for different phases of contract management. All commands use consistent flag-based argument parsing with built-in help documentation.

### CLI Commands

#### 1. Setup Environment
```bash
# Setup local Sui node and fund account
go run mockcontracts/main.go setup

# Get help for setup command
go run mockcontracts/main.go setup -h
```

#### 2. Parse Deployment Output
```bash
# Parse deployment JSON output and extract package ID (flag-based)
go run mockcontracts/main.go post-publish -file deployment_output.json

# Get help for post-publish command
go run mockcontracts/main.go post-publish -h
```

#### 3. Emit Events
```bash
# Emit events from all contracts
go run mockcontracts/main.go emit-events -package-id-file package_id.txt

# Use default package ID file location
go run mockcontracts/main.go emit-events

# Emit a single specific event
go run mockcontracts/main.go emit-single-event \
    -package-id-file package_id.txt \
    -function-name emit_static_config_set_event \
    -contract-name offramp

# Get help for any command
go run mockcontracts/main.go emit-events -h
go run mockcontracts/main.go emit-single-event -h
```

#### 4. General Help
```bash
# Show all available commands
go run mockcontracts/main.go -h

# Get help for any specific command
go run mockcontracts/main.go <command> -h
```

### Task Runner Integration (Recommended)

The easiest way to use the tool is through the integrated Task runner:

```bash
# Complete workflow: setup + publish + parse
task mockcontracts:setup        # Setup environment
task mockcontracts:publish      # Deploy contracts (depends on setup)
task mockcontracts:post-publish # Parse output (depends on publish)

# Emit events after deployment
task mockcontracts:emit-events

# Emit single event with task runner
task mockcontracts:emit-single-event -- -function-name emit_static_config_set_event -contract-name offramp

# Get help for CLI commands
task mockcontracts:help
```

### Direct Sui CLI Usage

For advanced users who want direct control:

```bash
# Publish contracts directly (after setup)
sui client publish --gas-budget 2000000000 --json --silence-warnings --dev \
    --with-unpublished-dependencies $PWD/contracts/test/sources/ > deployment_output.json
```

### Command Reference

| Command | Description | Example |
|---------|-------------|---------|
| `setup` | Setup local Sui node and fund account | `go run mockcontracts/main.go setup` |
| `post-publish` | Parse deployment output file | `go run mockcontracts/main.go post-publish -file deployment_output.json` |
| `emit-events` | Emit events from all contracts | `go run mockcontracts/main.go emit-events -package-id-file package_id.txt` |
| `emit-single-event` | Emit a specific event | `go run mockcontracts/main.go emit-single-event -function-name <EVENT> -contract-name <CONTRACT>` |

### Command-Specific Options

#### Setup Command
```bash
go run mockcontracts/main.go setup -h  # No additional flags currently
```

#### Post-Publish Command
| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `-file` | string | - | Deployment output file (required) |

```bash
# Flag-based (recommended)
go run mockcontracts/main.go post-publish -file deployment_output.json

# Positional argument (backwards compatible)
go run mockcontracts/main.go post-publish deployment_output.json
```

#### Emit Events Command
| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `-package-id-file` | string | `package_id.txt` | File containing the package ID |

#### Emit Single Event Command
| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `-package-id-file` | string | `package_id.txt` | File containing the package ID |
| `-function-name` | string | - | Name of the function to emit (required) |
| `-contract-name` | string | - | Name of the contract (required) |

### Supported Contracts and Events

The CLI now supports multiple contract emitters:

#### `offramp` Contract Events:
- `emit_static_config_set_event`
- `emit_dynamic_config_set_event`
- `emit_source_chain_config_set_event`
- `emit_skipped_already_executed_event`
- `emit_execution_state_changed_event`
- `emit_commit_report_accepted_event`
- `emit_skipped_report_execution_event`
- `emit_ocr_config_event`

#### Other Available Contracts:
- `token_admin_registry` - Token administration events
- `fee_quoter` - Fee calculation events
- `rmn_remote` - RMN (Risk Management Network) events
- `router` - Routing configuration events
- `onramp` - Message sending events
- `managed_token_pool` - Managed token pool events
- `burn_mint_token_pool` - Burn/mint token pool events
- `lock_release_token_pool` - Lock/release token pool events

> **Note**: Use `go run mockcontracts/main.go emit-single-event -h` to see the exact function names for each contract, or check the respective emitter files in `mockcontracts/events/`.


## Complete Workflow Example

### Step 1: Environment Setup
```bash
$ task mockcontracts:setup
# OR
$ go run mockcontracts/main.go setup

# Output:
# {"level":"info","msg":"Starting setup command"}
# {"level":"info","msg":"Sui node PID","pid":12345}
# {"level":"info","msg":"Active address","activeAddress":"0x..."}
# {"level":"info","msg":"Setup completed successfully"}
```

### Step 2: Deploy Contracts
```bash
$ task mockcontracts:publish
# OR
$ sui client publish --gas-budget 2000000000 --json --silence-warnings --dev \
    --with-unpublished-dependencies $PWD/contracts/test/sources/ > deployment_output.json

# Creates: deployment_output.json with transaction details
```

### Step 3: Parse Deployment Output
```bash
$ task mockcontracts:post-publish
# OR (new flag-based approach)
$ go run mockcontracts/main.go post-publish -file deployment_output.json
# OR (backwards compatible)
$ go run mockcontracts/main.go post-publish deployment_output.json

# Output:
# {"level":"info","msg":"Found published package","packageId":"0x1b1b08..."}
# {"level":"info","msg":"Package ID saved to package_id.txt"}

# Creates: package_id.txt with the package ID
```

### Step 4: Emit Events
```bash
# Emit all events in batch (new flag-based approach)
$ go run mockcontracts/main.go emit-events -package-id-file package_id.txt
# OR use default package ID file
$ go run mockcontracts/main.go emit-events

# OR emit a single event
$ go run mockcontracts/main.go emit-single-event \
    -package-id-file package_id.txt \
    -function-name emit_static_config_set_event \
    -contract-name offramp

# OR using task runner
$ task mockcontracts:emit-single-event -- -function-name emit_static_config_set_event -contract-name offramp

# Output:
# {"level":"info","msg":"Emitting static config set event"}
# {"level":"info","msg":"Executed PTB","tx":"..."}
# {"level":"info","msg":"Event emission completed"}
```

### Generated Files

After running the workflow, you'll have:
- `sui.pid` - Process ID of the local Sui node
- `deployment_output.json` - Full deployment transaction details
- `package_id.txt` - Extracted package ID for easy reference

## Network Configuration

### Local Network
```bash
# Default - assumes local node at localhost:9000
go run mockcontracts/deploy.go
```


## Task Runner Integration

The tool is fully integrated with the project's Task runner for streamlined workflows:

```bash
# Individual tasks
task mockcontracts:setup        # Setup local environment
task mockcontracts:publish      # Deploy contracts (auto-runs setup first)
task mockcontracts:post-publish # Parse deployment output (auto-runs publish first)
task mockcontracts:emit-events  # Emit events from all contracts

# Single event emission (pass arguments after --)
task mockcontracts:emit-single-event -- -function-name emit_static_config_set_event -contract-name offramp

# Get CLI help
task mockcontracts:help

# Full workflow (recommended)
task mockcontracts:post-publish # Runs setup → publish → post-publish
```

### Task Dependencies

The tasks have built-in dependencies to ensure proper execution order:
- `publish` depends on `setup` (environment must be ready)
- `post-publish` depends on `publish` (contracts must be deployed first)

### New Task Features

- **`emit-single-event`**: Now supports passing CLI arguments via `--`
- **`help`**: New task to show CLI help and usage examples
- **Flag-based commands**: All tasks now use the improved flag-based CLI

## Troubleshooting

### Common Issues

1. **"sui: command not found"**
   ```bash
   # Install Sui CLI
   cargo install --locked --git https://github.com/MystenLabs/sui.git --branch devnet sui
   ```

2. **"InsufficientGas" errors**
   ```bash
   # Check account balance
   sui client gas
   
   # Fund account manually
   sui client faucet
   
   # Or run setup again
   go run mockcontracts/main.go setup
   ```

3. **"Package ID not found"**
   - Ensure `deployment_output.json` exists and contains valid JSON
   - Check that the deployment was successful
   - Verify the file wasn't corrupted

4. **"Invalid function name" or "Unknown contract" errors**
   - Use `go run mockcontracts/main.go emit-single-event -h` for help
   - Check supported contracts and events list above
   - Ensure function name matches exactly (case-sensitive)
   - Verify the contract name is correct
   - Use `-function-name` instead of `-event-name` (updated flag name)

5. **Local node connection issues**
   ```bash
   # Check if node is running
   ps aux | grep sui
   
   # Restart if needed
   pkill sui
   go run mockcontracts/main.go setup
   ```

### File Locations

The tool creates several files in the current directory:
- `sui.pid` - Local Sui node process ID
- `deployment_output.json` - Full deployment transaction details  
- `package_id.txt` - Extracted package ID for easy reference

