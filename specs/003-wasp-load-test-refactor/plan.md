# Plan: Wallet-Pool WASP Load Tests for CCIP Message Sends

## Scope
- **In**: Message-only arbitrary message sends (Sui→EVM and EVM→Sui) using WASP with N parallel wallets
- **Out**: Token transfer tests — keep as sequential sends (existing `TestSui2EVM` / `TestEVM2Sui` token paths unchanged)

---

## Architecture

```
┌──────────────────────────────────────────────────────────────┐
│  Setup                                                       │
│  1. Generate N wallets from main key                         │
│  2. Fund each wallet from main (native token)                │
│  3. Each wallet prepares its own gas coin pool (Sui only)    │
└──────────────────────────┬───────────────────────────────────┘
                           │
┌──────────────────────────▼───────────────────────────────────┐
│  WASP Profile (N generators, 1 per wallet)                   │
│                                                              │
│  Generator 0 ──goroutine──> Gun0.Call() → Gun0.Call() → ... │
│    wallet: W0, gasPool: GP0                                  │
│                                                              │
│  Generator 1 ──goroutine──> Gun1.Call() → Gun1.Call() → ... │
│    wallet: W1, gasPool: GP1                                  │
│                                                              │
│  ...                                                         │
│                                                              │
│  Generator N-1 ──goroutine──> GunN.Call() → GunN.Call()     │
│    wallet: WN-1, gasPool: GPN-1                              │
│                                                              │
│  Each generator: 1 RPS per wallet                            │
│  Total throughput: N RPS                                     │
└──────────────────────────┬───────────────────────────────────┘
                           │
┌──────────────────────────▼───────────────────────────────────┐
│  Teardown                                                    │
│  Sweep remaining funds from all wallets back to main         │
└──────────────────────────────────────────────────────────────┘
```

**Why this works**: WASP generators are single-threaded (one goroutine per generator). Each generator owns its wallet exclusively → no nonce contention (EVM), no auth mutation (EVM), no coin object contention (Sui).

---

## Steps

### Phase 1: Config Changes

**1. Add `[load]` section to TOML run config**

File: `integration-tests/load/config/types.go`

```go
// In RunConfig:
type RunConfig struct {
    Run         RunParams         `toml:"run"`
    Receiver    ReceiverParams    `toml:"receiver"`
    Gas         GasParams         `toml:"gas"`
    Load        *LoadParams       `toml:"load,omitempty"`       // NEW
    Token       *TokenParams      `toml:"token,omitempty"`
    SuiReceiver *SuiReceiverParams `toml:"sui_receiver,omitempty"`
}

type LoadParams struct {
    RPS     int `toml:"rps"`               // target total RPS across all wallets
    Wallets int `toml:"wallets"`            // number of parallel wallets (default: 1)
}
```

File: `integration-tests/load/config/runconfig.go`
- Add validation: `load.rps >= 1`, `load.wallets >= 1`
- Defaults: if `[load]` section absent, default to `rps=1, wallets=1` (backward compatible)

File: `integration-tests/load/config/types.go` — `LoadTestConfig`
```go
type LoadTestConfig struct {
    // ... existing fields ...
    LoadRPS     int  // from LoadParams.RPS
    LoadWallets int  // from LoadParams.Wallets
}
```

**Example TOML:**
```toml
[run]
env = "staging"
source_chain_selector = "14480252757250338401"
dest_chain_selector = "16015286601757825753"
message_count = 100
message_data = "hello ccip"

[receiver]
address = "0x..."

[gas]
sui_gas_budget = 100000000
evm_gas_limit = 200000
evm_callback_gas_limit = 200000

[load]
rps = 5        # 5 messages/second total
wallets = 5    # 5 parallel wallets → 1 RPS per wallet
```

### Phase 2: Wallet Generation & Funding

**2. Create `integration-tests/load/wallet/wallet.go`**

```go
package wallet

// Wallet holds a funded account ready for load testing.
type Wallet struct {
    Address string          // Sui bech32 address or EVM hex address
    // For Sui:
    SuiSigner bindutils.SuiSigner  // nil for EVM wallets
    // For EVM:
    EVMTransactOpts *bind.TransactOpts  // nil for Sui wallets
    EVMPrivKey     *ecdsa.PrivateKey    // nil for Sui wallets
}

// GenerateSuiWallets creates N independent Ed25519 keypairs.
func GenerateSuiWallets(n int) ([]*Wallet, error)

// GenerateEVMWallets creates N independent secp256k1 keypairs with TransactOpts.
func GenerateEVMWallets(n int, chainID *big.Int) ([]*Wallet, error)
```

**3. Create `integration-tests/load/wallet/funding.go`**

```go
// FundSuiWallets sends `amount` MIST of native SUI from mainSigner to each wallet.
// Uses parallel transfers (one tx per wallet).
func FundSuiWallets(
    ctx context.Context,
    ptbClient *client.PTBClient,
    mainSigner bindutils.SuiSigner,
    mainAddress string,
    wallets []*Wallet,
    amountPerWallet uint64,  // in MIST
) error

// FundEVMWallets sends `amount` wei of native ETH from mainAuth to each wallet.
// Uses parallel transfers.
func FundEVMWallets(
    ctx context.Context,
    ethClient *ethclient.Client,
    mainAuth *bind.TransactOpts,
    wallets []*Wallet,
    amountPerWallet *big.Int,
) error

// SweepSuiWallets sends remaining SUI balance from each wallet back to mainAddress.
func SweepSuiWallets(
    ctx context.Context,
    ptbClient *client.PTBClient,
    wallets []*Wallet,
    mainAddress string,
) error

// SweepEVMWallets sends remaining ETH balance from each wallet back to mainAddress.
func SweepEVMWallets(
    ctx context.Context,
    ethClient *ethclient.Client,
    wallets []*Wallet,
    mainAddress common.Address,
) error
```

### Phase 3: Gun Implementations

**4. Create `integration-tests/load/sui2evm_gun.go`**

```go
package load

// Sui2EVMMsgGun sends message-only CCIP messages from Sui to EVM.
// Each instance owns one wallet and its gas coin pool exclusively.
type Sui2EVMMsgGun struct {
    wallet          *wallet.Wallet
    gasPool         *sui.SuiCoinPool
    ptbClient       *client.PTBClient
    // CCIP addresses (shared, read-only)
    ccipPkgID       string
    onRampPkgID     string
    ccipObjectRefID string
    onRampStateID   string
    feeTokenType    string
    suiCoinMetaID   string
    // Message params
    destChainSelector uint64
    receiver          []byte
    data              []byte
    gasBudget         uint64
    evmCallbackGas    uint64
    // Results
    resultsCh         chan<- config.SentMessage
}

func (g *Sui2EVMMsgGun) Call(_ *wasp.Generator) *wasp.Response {
    ctx := context.Background()

    // Pop 2 coins: gas + fee
    gasCoin, err := g.gasPool.Pop(ctx)
    if err != nil { return failResponse(err) }
    feeCoin, err := g.gasPool.Pop(ctx)
    if err != nil { return failResponse(err) }

    msgID, txDigest, seqNum, err := sui.SendMessage(
        ctx, g.ptbClient, g.wallet.SuiSigner,
        g.ccipPkgID, g.onRampPkgID, g.ccipObjectRefID, g.onRampStateID,
        gasCoin, g.feeTokenType, g.suiCoinMetaID, feeCoin,
        g.destChainSelector, g.receiver, g.data,
        g.gasBudget, g.evmCallbackGas,
    )

    // Push result to channel
    g.resultsCh <- config.SentMessage{...}

    if err != nil {
        return &wasp.Response{Failed: true, Error: err.Error(), Group: "sui->evm"}
    }
    return &wasp.Response{Failed: false, Group: "sui->evm",
        Data: map[string]string{"messageID": msgID, "txDigest": txDigest}}
}
```

**5. Create `integration-tests/load/evm2sui_gun.go`**

```go
package load

// EVM2SuiMsgGun sends message-only CCIP messages from EVM to Sui.
// Each instance owns one wallet (TransactOpts) exclusively.
type EVM2SuiMsgGun struct {
    wallet          *wallet.Wallet  // owns TransactOpts exclusively
    ethClient       *ethclient.Client
    routerAddress   common.Address
    destChainSelector uint64
    receiver        [32]byte
    data            []byte
    extraArgs       []byte
    resultsCh       chan<- config.SentMessage
}

func (g *EVM2SuiMsgGun) Call(_ *wasp.Generator) *wasp.Response {
    ctx := context.Background()

    // No nonce management needed — single goroutine owns this auth
    msgID, txHash, err := evm.SendMessage(
        ctx, g.ethClient, g.wallet.EVMTransactOpts,
        g.routerAddress, g.destChainSelector,
        g.receiver[:], g.data, g.extraArgs,
    )

    g.resultsCh <- config.SentMessage{...}

    if err != nil {
        return &wasp.Response{Failed: true, Error: err.Error(), Group: "evm->sui"}
    }
    return &wasp.Response{Failed: false, Group: "evm->sui",
        Data: map[string]string{"messageID": msgID, "txHash": txHash}}
}
```

### Phase 4: Refactor Test Functions

**6. Refactor `TestSui2EVM` — message-only path**

```go
func TestSui2EVM(t *testing.T) {
    // ... existing setup (config, addresses, receiver) ...

    // Token mode → keep sequential (existing code path)
    if cfg.TokenConfig != nil {
        runSui2EVMSequential(t, cfg, ...)  // extract existing code
        return
    }

    // Message-only → WASP with wallet pool
    runSui2EVMWASP(t, cfg, ...)
}

func runSui2EVMWASP(t *testing.T, cfg *config.LoadTestConfig, ...) {
    ctx := context.Background()
    N := cfg.LoadWallets  // e.g., 5

    // Step 1: Generate N wallets
    wallets, err := wallet.GenerateSuiWallets(N)  // returns []*Wallet with len=N

    // Step 2: Fund each wallet
    // Estimate total cost: N wallets × (msgPerWallet × 2 × splitAmount + buffer)
    fundingPerWallet := estimateSuiFunding(cfg, N)
    err = wallet.FundSuiWallets(ctx, ptbClient, mainSigner, mainAddr, wallets, fundingPerWallet)

    // Step 3: Create N generators (one per wallet)
    // Each generator gets its own gun with its own wallet and gas pool
    msgPerWallet := cfg.MessageCount / N
    splitAmount := sui.RecommendedSplitAmountPerCoin(estimatedFee)
    resultsCh := make(chan config.SentMessage, cfg.MessageCount)
    
    p := wasp.NewProfile()
    rpsPerWallet := float64(cfg.LoadRPS) / float64(N)
    duration := time.Duration(cfg.MessageCount/N) * time.Second

    // THIS LOOP CREATES N GENERATORS — one per wallet
    for i := 0; i < N; i++ {
        w := wallets[i]
        
        // Each wallet prepares its own gas coin pool
        gasPool, _ := sui.PrepareSuiCoinPool(ctx, ptbClient, w.SuiSigner, w.Address,
            msgPerWallet, splitAmount)

        // Create a gun for this wallet
        gun := &Sui2EVMMsgGun{
            wallet:    w,           // wallet #i
            gasPool:   gasPool,     // gas pool #i
            ptbClient: ptbClient,
            // ... fill in CCIP addresses, message params ...
            resultsCh: resultsCh,
        }

        // Create ONE generator for this gun
        gen := wasp.NewGenerator(&wasp.Config{
            GenName:  fmt.Sprintf("sui2evm-w%d", i),
            LoadType: wasp.RPS,
            Schedule: wasp.Plain(rpsPerWallet, duration),
            Gun:      gun,  // this generator uses this gun (and its wallet)
        })

        // Add generator #i to profile
        p.Add(gen)
    }

    // Step 4: Run all N generators in parallel
    // WASP starts N goroutines, each calling its gun's Call() method
    p.Run(true)
    close(resultsCh)

    // Step 5: Collect results
    results := collectResults(cfg, resultsCh)
    config.SaveResults(results)

    // Step 6: Sweep funds back
    wallet.SweepSuiWallets(ctx, ptbClient, wallets, mainAddr)
}
```

**7. Refactor `TestEVM2Sui` — message-only path**

Same pattern:
```go
func TestEVM2Sui(t *testing.T) {
    // ... existing setup ...

    if cfg.TokenConfig != nil {
        runEVM2SuiSequential(t, cfg, ...)  // extract existing code
        return
    }

    runEVM2SuiWASP(t, cfg, ...)
}

func runEVM2SuiWASP(t *testing.T, cfg *config.LoadTestConfig, ...) {
    N := cfg.LoadWallets  // e.g., 5

    // Step 1: Generate N EVM wallets
    wallets, _ := wallet.GenerateEVMWallets(N, chainID)  // returns []*Wallet with len=N

    // Step 2: Fund each wallet with ETH
    fundingPerWallet := estimateEvmFunding(cfg, N)
    wallet.FundEVMWallets(ctx, ethClient, mainAuth, wallets, fundingPerWallet)

    // Step 3: Create N generators (one per wallet)
    // Each generator gets its own gun with its own wallet
    resultsCh := make(chan config.SentMessage, cfg.MessageCount)
    
    p := wasp.NewProfile()
    rpsPerWallet := float64(cfg.LoadRPS) / float64(N)
    duration := time.Duration(cfg.MessageCount/N) * time.Second

    // THIS LOOP CREATES N GENERATORS — one per wallet
    for i := 0; i < N; i++ {
        w := wallets[i]

        // Create a gun for this wallet
        gun := &EVM2SuiMsgGun{
            wallet:        w,  // wallet #i
            ethClient:     ethClient,
            routerAddress: routerAddr,
            // ... fill in params ...
            resultsCh: resultsCh,
        }

        // Create ONE generator for this gun
        gen := wasp.NewGenerator(&wasp.Config{
            GenName:  fmt.Sprintf("evm2sui-w%d", i),
            LoadType: wasp.RPS,
            Schedule: wasp.Plain(rpsPerWallet, duration),
            Gun:      gun,  // this generator uses this gun (and its wallet)
        })

        // Add generator #i to profile
        p.Add(gen)
    }

    // Step 4: Run all N generators in parallel
    // WASP starts N goroutines, each calling its gun's Call() method
    p.Run(true)
    close(resultsCh)

    // Step 5: Collect results
    results := collectResults(cfg, resultsCh)
    config.SaveResults(results)

    // Step 6: Sweep funds back
    wallet.SweepEVMWallets(ctx, ethClient, wallets, mainAddress)
}
```

### Phase 5: Helpers

**8. Create `integration-tests/load/results_collector.go`**

```go
// collectResults drains the results channel and builds RunResults.
func collectResults(cfg *config.LoadTestConfig, ch <-chan config.SentMessage) *config.RunResults {
    results := &config.RunResults{
        RunName:             cfg.RunName,
        EnvName:             cfg.EnvName,
        SourceChainSelector: cfg.SourceChainSelector,
        DestChainSelector:   cfg.DestChainSelector,
        TotalMessages:       cfg.MessageCount,
        RunStarted:          ...,
        Messages:            make([]config.SentMessage, 0),
    }
    for msg := range ch {
        results.Messages = append(results.Messages, msg)
        if msg.Success { results.SuccessfulMessages++ } else { results.FailedMessages++ }
    }
    results.RunEnded = time.Now().Format(time.RFC3339)
    return results
}
```

**9. Create `integration-tests/load/funding_estimator.go`**

```go
// estimateSuiFunding calculates how much SUI (in MIST) each wallet needs.
// = (msgPerWallet × 2 coins × splitAmount) + gasBudgetBuffer
func estimateSuiFunding(cfg *config.LoadTestConfig, numWallets int) uint64

// estimateEvmFunding calculates how much ETH (in wei) each wallet needs.
// = msgPerWallet × (estimatedFee + gasCost) × 1.5 buffer
func estimateEvmFunding(cfg *config.LoadTestConfig, numWallets int) *big.Int
```

---

## Files Summary

### New files
| File | Purpose |
|------|---------|
| `load/wallet/wallet.go` | Wallet generation (Sui + EVM) |
| `load/wallet/funding.go` | Fund wallets, sweep funds back |
| `load/sui2evm_gun.go` | `Sui2EVMMsgGun` implementing `wasp.Gun` |
| `load/evm2sui_gun.go` | `EVM2SuiMsgGun` implementing `wasp.Gun` |
| `load/results_collector.go` | Channel-based results aggregation |
| `load/funding_estimator.go` | Calculate funding amounts per wallet |

### Modified files
| File | Change |
|------|--------|
| `load/config/types.go` | Add `LoadParams` struct, `LoadRPS`/`LoadWallets` to `LoadTestConfig` |
| `load/config/runconfig.go` | Add validation for `[load]` section |
| `load/sui2evm_test.go` | Split into sequential (token) + WASP (message-only) paths |
| `load/evm2sui_test.go` | Split into sequential (token) + WASP (message-only) paths |
| `go.mod` | Add `chainlink-testing-framework/wasp` dependency |

### Unchanged files
| File | Reason |
|------|--------|
| `load/sui/sender.go` | `SendMessage` already concurrent-safe |
| `load/sui/coin_pool.go` | Channel-based, concurrent-safe |
| `load/sui/client.go` | Connection-pooled, concurrent-safe |
| `load/evm/sender.go` | `SendMessage` works with dedicated auth |
| `load/evm/client.go` | No changes needed |
| `load/config/config.go` | `LoadFullConfig` passes through new fields |

---

## Problems & Solutions

| # | Problem | Solution |
|---|---------|----------|
| 1 | **EVM nonce contention** | Not a problem — 1 generator = 1 wallet = 1 goroutine. go-ethereum auto-manages nonce. |
| 2 | **EVM auth mutation** | Not a problem — each generator exclusively owns its `*bind.TransactOpts`. |
| 3 | **Sui coin object contention** | Not a problem — each wallet has its own gas coins from its own pool. |
| 4 | **Funding calculation** | Estimate per-wallet cost: `(msgCount/wallets) × 2 × splitAmount + gasBuffer`. Fund with 1.5× safety margin. |
| 5 | **Sweep failures** | If a wallet's remaining balance < gas cost of sweep tx, skip it. Log warning. |
| 6 | **Wallet generation determinism** | Use `crypto.GenerateKey()` (random). If crash recovery needed, can switch to BIP39 mnemonic later. |
| 7 | **RPS vs message count** | `message_count` determines total messages. `rps` determines rate. `duration = message_count / rps`. WASP schedule: `Plain(rpsPerWallet, duration)`. |
| 8 | **Uneven wallet usage** | Not an issue — each wallet sends exactly `messageCount / N` messages. |

---

## Verification

1. **Single wallet (N=1)**: Should behave identically to current sequential test
2. **Multiple wallets (N=5)**: Verify ~5× throughput improvement
3. **Funding**: Verify wallets receive correct amounts, sweep returns funds
4. **Results parity**: JSON output matches existing format
5. **Error handling**: Kill RPC mid-test → WASP records failures, sweep still runs
6. **Token mode unchanged**: Verify token transfer tests still work sequentially

---

## Decisions

- **Message-only only**: Token transfers stay sequential. Separate concern, separate PR.
- **1 generator = 1 wallet**: Simpler than semaphore pool. No contention by design.
- **Random keypairs**: Simpler than BIP39 derivation. Crash recovery is a future concern.
- **Channel-based results**: Avoids mutex on shared slice. Gun pushes to channel, collector drains.
- **Sweep on teardown**: Always attempt sweep, even on partial failure. Best-effort.
- **`[load]` section optional**: Defaults to `rps=1, wallets=1` for backward compatibility.
