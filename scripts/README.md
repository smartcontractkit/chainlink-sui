# Deployment Scripts

This directory contains scripts for deploying and interacting with the Sui contracts.

## Available Scripts

### `deploy_sample_contracts.sh`

This script deploys the sample contracts to a local Sui network.

#### Prerequisites

- A local Sui network must be running
- The `sui` CLI must be installed and configured
  - Once installed, use `sui client active-address` and proceed with the CLI instructions
  - Then run `sui start` in a separate terminal window to start a local Sui network before running the deploy scripts

#### Usage

```bash
./deploy_sample_contracts.sh
```

The script will:
1. Build the contracts from the `contracts/test` directory
2. Publish the package to the local Sui network
3. Create a Counter object by calling the `initialize` function
4. Display the package ID and example commands for interacting with the contracts

#### Example Output

```
Result: Transaction Digest: 5e2eBKVzFBWeWknCQ91WxB23Fdfg7GrhKYZJjhkHeQ13
╭──────────────────────────────────────────────────────────────────────────────────────────────────────────────╮
│ Transaction Data                                                                                             │
├──────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
│ Sender: 0xfb40f8c84ef92c377df954cb99cd603b51decbc980dfcb4f2bc99410d66cc22b                                   │
│ Gas Owner: 0xfb40f8c84ef92c377df954cb99cd603b51decbc980dfcb4f2bc99410d66cc22b                                │
│ Gas Budget: 20000000 MIST                                                                                    │
│ Gas Price: 1000 MIST                                                                                         │
│ Gas Payment:                                                                                                 │
│  ┌──                                                                                                         │
│  │ ID: 0x340c175c5ee5132bd25990df9c987e2296ffbe52514f3a380481ed49b93786a1                                    │
│  │ Version: 10                                                                                               │
│  │ Digest: J342gqVVGV2MRALUi8BX9PjcwodUt8NbuRhy9qGgANGk                                                      │
│  └──                                                                                                         │
│                                                                                                              │
│ Transaction Kind: Programmable                                                                               │
│ ╭──────────────────────────────────────────────────────────────────────────────────────────────────────────╮ │
│ │ Input Objects                                                                                            │ │
│ ├──────────────────────────────────────────────────────────────────────────────────────────────────────────┤ │
│ │ 0   Pure Arg: Type: address, Value: "0xfb40f8c84ef92c377df954cb99cd603b51decbc980dfcb4f2bc99410d66cc22b" │ │
│ ╰──────────────────────────────────────────────────────────────────────────────────────────────────────────╯ │
│ ╭─────────────────────────────────────────────────────────────────────────╮                                  │
│ │ Commands                                                                │                                  │
│ ├─────────────────────────────────────────────────────────────────────────┤                                  │
│ │ 0  Publish:                                                             │                                  │
│ │  ┌                                                                      │                                  │
│ │  │ Dependencies:                                                        │                                  │
│ │  │   0x0000000000000000000000000000000000000000000000000000000000000001 │                                  │
│ │  │   0x0000000000000000000000000000000000000000000000000000000000000002 │                                  │
│ │  └                                                                      │                                  │
│ │                                                                         │                                  │
│ │ 1  TransferObjects:                                                     │                                  │
│ │  ┌                                                                      │                                  │
│ │  │ Arguments:                                                           │                                  │
│ │  │   Result 0                                                           │                                  │
│ │  │ Address: Input  0                                                    │                                  │
│ │  └                                                                      │                                  │
│ ╰─────────────────────────────────────────────────────────────────────────╯                                  │
│                                                                                                              │
│ Signatures:                                                                                                  │
│    VYTGY1Jf4NnzEsDS7uIXB4Iqd3srMXmU6rGiy9/Bai9U3F+BoQruZ1RF+FdpGnK2wcX8cNr8nTXM7gplimKJCw==                  │
│                                                                                                              │
╰──────────────────────────────────────────────────────────────────────────────────────────────────────────────╯
╭───────────────────────────────────────────────────────────────────────────────────────────────────╮
│ Transaction Effects                                                                               │
├───────────────────────────────────────────────────────────────────────────────────────────────────┤
│ Digest: 5e2eBKVzFBWeWknCQ91WxB23Fdfg7GrhKYZJjhkHeQ13                                              │
│ Status: Success                                                                                   │
│ Executed Epoch: 0                                                                                 │
│                                                                                                   │
│ Created Objects:                                                                                  │
│  ┌──                                                                                              │
│  │ ID: 0x6c362339352f55944a135bc43265fd02d9a12aff35567c97e35153e88358a015                         │
│  │ Owner: Shared( 11 )                                                                            │
│  │ Version: 11                                                                                    │
│  │ Digest: 8Hit9rKmNDS6nD9ZBMqKi4CUU6M5TNTSQz41UKMKZqbf                                           │
│  └──                                                                                              │
│  ┌──                                                                                              │
│  │ ID: 0xabf6d4f0df648f9305bb32f8de07486b438d597779d9611e64c3add4044e71fc                         │
│  │ Owner: Immutable                                                                               │
│  │ Version: 1                                                                                     │
│  │ Digest: BAeswdrWgEx3CFHtbDhr2dY8vzHPq3Utz7tTQUqQf3sY                                           │
│  └──                                                                                              │
│  ┌──                                                                                              │
│  │ ID: 0xf122d4b44169c16c83747b9b5c244447c5e573bfd9a5c155165c469cecc97a2e                         │
│  │ Owner: Account Address ( 0xfb40f8c84ef92c377df954cb99cd603b51decbc980dfcb4f2bc99410d66cc22b )  │
│  │ Version: 11                                                                                    │
│  │ Digest: HqsnG9bNT7G7XhPCEtsA86coNb7mbqsU5nockKibngxD                                           │
│  └──                                                                                              │
│ Mutated Objects:                                                                                  │
│  ┌──                                                                                              │
│  │ ID: 0x340c175c5ee5132bd25990df9c987e2296ffbe52514f3a380481ed49b93786a1                         │
│  │ Owner: Account Address ( 0xfb40f8c84ef92c377df954cb99cd603b51decbc980dfcb4f2bc99410d66cc22b )  │
│  │ Version: 11                                                                                    │
│  │ Digest: JAUFqmgDdWTvvCSbsFtXBLTZD59s7BNgkW69YqveB2yb                                           │
│  └──                                                                                              │
│ Gas Object:                                                                                       │
│  ┌──                                                                                              │
│  │ ID: 0x340c175c5ee5132bd25990df9c987e2296ffbe52514f3a380481ed49b93786a1                         │
│  │ Owner: Account Address ( 0xfb40f8c84ef92c377df954cb99cd603b51decbc980dfcb4f2bc99410d66cc22b )  │
│  │ Version: 11                                                                                    │
│  │ Digest: JAUFqmgDdWTvvCSbsFtXBLTZD59s7BNgkW69YqveB2yb                                           │
│  └──                                                                                              │
│ Gas Cost Summary:                                                                                 │
│    Storage Cost: 16142400 MIST                                                                    │
│    Computation Cost: 1000000 MIST                                                                 │
│    Storage Rebate: 978120 MIST                                                                    │
│    Non-refundable Storage Fee: 9880 MIST                                                          │
│                                                                                                   │
│ Transaction Dependencies:                                                                         │
│    FiLTSJFq3vdvDHb71g8wB2s28pHkXJGwFQ8rQ9eqxz8                                                    │
│    HjRN4q5iN2TTNC8FPX9Vi6SbANBTCNAm6JFGxu4Kb8vK                                                   │
╰───────────────────────────────────────────────────────────────────────────────────────────────────╯
╭─────────────────────────────╮
│ No transaction block events │
╰─────────────────────────────╯

╭──────────────────────────────────────────────────────────────────────────────────────────────────────╮
│ Object Changes                                                                                       │
├──────────────────────────────────────────────────────────────────────────────────────────────────────┤
│ Created Objects:                                                                                     │
│  ┌──                                                                                                 │
│  │ ObjectID: 0x6c362339352f55944a135bc43265fd02d9a12aff35567c97e35153e88358a015                      │
│  │ Sender: 0xfb40f8c84ef92c377df954cb99cd603b51decbc980dfcb4f2bc99410d66cc22b                        │
│  │ Owner: Shared( 11 )                                                                               │
│  │ ObjectType: 0xabf6d4f0df648f9305bb32f8de07486b438d597779d9611e64c3add4044e71fc::echo::EventStore  │
│  │ Version: 11                                                                                       │
│  │ Digest: 8Hit9rKmNDS6nD9ZBMqKi4CUU6M5TNTSQz41UKMKZqbf                                              │
│  └──                                                                                                 │
│  ┌──                                                                                                 │
│  │ ObjectID: 0xf122d4b44169c16c83747b9b5c244447c5e573bfd9a5c155165c469cecc97a2e                      │
│  │ Sender: 0xfb40f8c84ef92c377df954cb99cd603b51decbc980dfcb4f2bc99410d66cc22b                        │
│  │ Owner: Account Address ( 0xfb40f8c84ef92c377df954cb99cd603b51decbc980dfcb4f2bc99410d66cc22b )     │
│  │ ObjectType: 0x2::package::UpgradeCap                                                              │
│  │ Version: 11                                                                                       │
│  │ Digest: HqsnG9bNT7G7XhPCEtsA86coNb7mbqsU5nockKibngxD                                              │
│  └──                                                                                                 │
│ Mutated Objects:                                                                                     │
│  ┌──                                                                                                 │
│  │ ObjectID: 0x340c175c5ee5132bd25990df9c987e2296ffbe52514f3a380481ed49b93786a1                      │
│  │ Sender: 0xfb40f8c84ef92c377df954cb99cd603b51decbc980dfcb4f2bc99410d66cc22b                        │
│  │ Owner: Account Address ( 0xfb40f8c84ef92c377df954cb99cd603b51decbc980dfcb4f2bc99410d66cc22b )     │
│  │ ObjectType: 0x2::coin::Coin<0x2::sui::SUI>                                                        │
│  │ Version: 11                                                                                       │
│  │ Digest: JAUFqmgDdWTvvCSbsFtXBLTZD59s7BNgkW69YqveB2yb                                              │
│  └──                                                                                                 │
│ Published Objects:                                                                                   │
│  ┌──                                                                                                 │
│  │ PackageID: 0xabf6d4f0df648f9305bb32f8de07486b438d597779d9611e64c3add4044e71fc                     │
│  │ Version: 1                                                                                        │
│  │ Digest: BAeswdrWgEx3CFHtbDhr2dY8vzHPq3Utz7tTQUqQf3sY                                              │
│  │ Modules: counter, echo                                                                            │
│  └──                                                                                                 │
╰──────────────────────────────────────────────────────────────────────────────────────────────────────╯
╭───────────────────────────────────────────────────────────────────────────────────────────────────╮
│ Balance Changes                                                                                   │
├───────────────────────────────────────────────────────────────────────────────────────────────────┤
│  ┌──                                                                                              │
│  │ Owner: Account Address ( 0xfb40f8c84ef92c377df954cb99cd603b51decbc980dfcb4f2bc99410d66cc22b )  │
│  │ CoinType: 0x2::sui::SUI                                                                        │
│  │ Amount: -16164280                                                                              │
│  └──                                                                                              │
╰───────────────────────────────────────────────────────────────────────────────────────────────────╯
Contracts deployed successfully!
Package ID: 0xabf6d4f0df648f9305bb32f8de07486b438d597779d9611e64c3add4044e71fc
Creating a counter object...
Transaction Digest: DNkL3w6WmJqDQAcBeVyKbiMQdkBqMm2Uyo8JMN9EkSKC
╭─────────────────────────────────────────────────────────────────────────────────────────────╮
│ Transaction Data                                                                            │
├─────────────────────────────────────────────────────────────────────────────────────────────┤
│ Sender: 0xfb40f8c84ef92c377df954cb99cd603b51decbc980dfcb4f2bc99410d66cc22b                  │
│ Gas Owner: 0xfb40f8c84ef92c377df954cb99cd603b51decbc980dfcb4f2bc99410d66cc22b               │
│ Gas Budget: 20000000 MIST                                                                   │
│ Gas Price: 1000 MIST                                                                        │
│ Gas Payment:                                                                                │
│  ┌──                                                                                        │
│  │ ID: 0x340c175c5ee5132bd25990df9c987e2296ffbe52514f3a380481ed49b93786a1                   │
│  │ Version: 11                                                                              │
│  │ Digest: JAUFqmgDdWTvvCSbsFtXBLTZD59s7BNgkW69YqveB2yb                                     │
│  └──                                                                                        │
│                                                                                             │
│ Transaction Kind: Programmable                                                              │
│   No input objects for this transaction                                                     │
│ ╭──────────────────────────────────────────────────────────────────────────────────╮        │
│ │ Commands                                                                         │        │
│ ├──────────────────────────────────────────────────────────────────────────────────┤        │
│ │ 0  MoveCall:                                                                     │        │
│ │  ┌                                                                               │        │
│ │  │ Function:  initialize                                                         │        │
│ │  │ Module:    counter                                                            │        │
│ │  │ Package:   0xabf6d4f0df648f9305bb32f8de07486b438d597779d9611e64c3add4044e71fc │        │
│ │  └                                                                               │        │
│ ╰──────────────────────────────────────────────────────────────────────────────────╯        │
│                                                                                             │
│ Signatures:                                                                                 │
│    wHTzWcIsV2T+0ZVMiPoN/Q0Ystc4NS2S/LhdJDfumYY+MPT0Xec0uhaQwukIo/peLQfeh7E8CQtEncIFuYACAw== │
│                                                                                             │
╰─────────────────────────────────────────────────────────────────────────────────────────────╯
╭───────────────────────────────────────────────────────────────────────────────────────────────────╮
│ Transaction Effects                                                                               │
├───────────────────────────────────────────────────────────────────────────────────────────────────┤
│ Digest: DNkL3w6WmJqDQAcBeVyKbiMQdkBqMm2Uyo8JMN9EkSKC                                              │
│ Status: Success                                                                                   │
│ Executed Epoch: 0                                                                                 │
│                                                                                                   │
│ Created Objects:                                                                                  │
│  ┌──                                                                                              │
│  │ ID: 0xd2aeeaf835710b7c69fd74c19195bebce7554f068c4a6c505205f898e4070926                         │
│  │ Owner: Shared( 12 )                                                                            │
│  │ Version: 12                                                                                    │
│  │ Digest: 6tuLer5aTA4pUV6TZB7LvjnEhenvnq5tyQvfmSxRXudN                                           │
│  └──                                                                                              │
│ Mutated Objects:                                                                                  │
│  ┌──                                                                                              │
│  │ ID: 0x340c175c5ee5132bd25990df9c987e2296ffbe52514f3a380481ed49b93786a1                         │
│  │ Owner: Account Address ( 0xfb40f8c84ef92c377df954cb99cd603b51decbc980dfcb4f2bc99410d66cc22b )  │
│  │ Version: 12                                                                                    │
│  │ Digest: 7E2eTbZKUu11DsBdPnybKJ9pnAa134XiXAciSddLYsFs                                           │
│  └──                                                                                              │
│ Gas Object:                                                                                       │
│  ┌──                                                                                              │
│  │ ID: 0x340c175c5ee5132bd25990df9c987e2296ffbe52514f3a380481ed49b93786a1                         │
│  │ Owner: Account Address ( 0xfb40f8c84ef92c377df954cb99cd603b51decbc980dfcb4f2bc99410d66cc22b )  │
│  │ Version: 12                                                                                    │
│  │ Digest: 7E2eTbZKUu11DsBdPnybKJ9pnAa134XiXAciSddLYsFs                                           │
│  └──                                                                                              │
│ Gas Cost Summary:                                                                                 │
│    Storage Cost: 2348400 MIST                                                                     │
│    Computation Cost: 1000000 MIST                                                                 │
│    Storage Rebate: 978120 MIST                                                                    │
│    Non-refundable Storage Fee: 9880 MIST                                                          │
│                                                                                                   │
│ Transaction Dependencies:                                                                         │
│    5e2eBKVzFBWeWknCQ91WxB23Fdfg7GrhKYZJjhkHeQ13                                                   │
╰───────────────────────────────────────────────────────────────────────────────────────────────────╯
╭─────────────────────────────╮
│ No transaction block events │
╰─────────────────────────────╯

╭──────────────────────────────────────────────────────────────────────────────────────────────────────╮
│ Object Changes                                                                                       │
├──────────────────────────────────────────────────────────────────────────────────────────────────────┤
│ Created Objects:                                                                                     │
│  ┌──                                                                                                 │
│  │ ObjectID: 0xd2aeeaf835710b7c69fd74c19195bebce7554f068c4a6c505205f898e4070926                      │
│  │ Sender: 0xfb40f8c84ef92c377df954cb99cd603b51decbc980dfcb4f2bc99410d66cc22b                        │
│  │ Owner: Shared( 12 )                                                                               │
│  │ ObjectType: 0xabf6d4f0df648f9305bb32f8de07486b438d597779d9611e64c3add4044e71fc::counter::Counter  │
│  │ Version: 12                                                                                       │
│  │ Digest: 6tuLer5aTA4pUV6TZB7LvjnEhenvnq5tyQvfmSxRXudN                                              │
│  └──                                                                                                 │
│ Mutated Objects:                                                                                     │
│  ┌──                                                                                                 │
│  │ ObjectID: 0x340c175c5ee5132bd25990df9c987e2296ffbe52514f3a380481ed49b93786a1                      │
│  │ Sender: 0xfb40f8c84ef92c377df954cb99cd603b51decbc980dfcb4f2bc99410d66cc22b                        │
│  │ Owner: Account Address ( 0xfb40f8c84ef92c377df954cb99cd603b51decbc980dfcb4f2bc99410d66cc22b )     │
│  │ ObjectType: 0x2::coin::Coin<0x2::sui::SUI>                                                        │
│  │ Version: 12                                                                                       │
│  │ Digest: 7E2eTbZKUu11DsBdPnybKJ9pnAa134XiXAciSddLYsFs                                              │
│  └──                                                                                                 │
╰──────────────────────────────────────────────────────────────────────────────────────────────────────╯
╭───────────────────────────────────────────────────────────────────────────────────────────────────╮
│ Balance Changes                                                                                   │
├───────────────────────────────────────────────────────────────────────────────────────────────────┤
│  ┌──                                                                                              │
│  │ Owner: Account Address ( 0xfb40f8c84ef92c377df954cb99cd603b51decbc980dfcb4f2bc99410d66cc22b )  │
│  │ CoinType: 0x2::sui::SUI                                                                        │
│  │ Amount: -2370280                                                                               │
│  └──                                                                                              │
╰───────────────────────────────────────────────────────────────────────────────────────────────────╯
Deployment complete!
You can now interact with the deployed contracts.
Example commands:
  sui client --url http://localhost:9000 call --package 0xabf6d4f0df648f9305bb32f8de07486b438d597779d9611e64c3add4044e71fc --module counter --function increment --args $COUNTER_ID --gas-budget 10000000
  sui client --url http://localhost:9000 call --package 0xabf6d4f0df648f9305bb32f8de07486b438d597779d9611e64c3add4044e71fc --module counter --function increment_mult --args $COUNTER_ID 5 10 --gas-budget 10000000
```
