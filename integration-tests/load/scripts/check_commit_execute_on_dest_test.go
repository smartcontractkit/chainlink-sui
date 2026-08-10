package scripts

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	evm2evmofframp "github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/evm_2_evm_offramp"
	offramp160 "github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/offramp"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink-sui/integration-tests/load/config"
)

// TestCheckCommitAndExecuteOnDest reads a results file from a previous Sui→EVM
// load test run and queries the destination EVM chain's OffRamp for
// ExecutionStateChanged events for each message.
//
// Usage: go test -v -run TestCheckCommitAndExecuteOnDest ./scripts/
func TestCheckCommitAndExecuteOnDest(t *testing.T) {
	// Config functions use relative paths from load/; when running `go test ./scripts/`
	// the CWD is scripts/, so chdir up.
	cwd, _ := os.Getwd()
	if err := os.Chdir(".."); err != nil {
		t.Fatalf("Failed to chdir: %v", err)
	}
	t.Cleanup(func() { os.Chdir(cwd) })

	ctx := context.Background()

	// ── Load results ──
	resultsRaw, err := os.ReadFile("results/real-sui-to-evm-run-2-testnet-20260810T113420.txt")
	if err != nil {
		t.Fatalf("Failed to read results: %v", err)
	}
	var results config.RunResults
	if err := json.Unmarshal(resultsRaw, &results); err != nil {
		t.Fatalf("Failed to parse results: %v", err)
	}
	t.Logf("Loaded %d messages (env=%s, destChain=%d)", len(results.Messages), results.EnvName, results.DestChainSelector)

	// ── Load config ──
	addressBook, err := config.LoadAddressBook(results.EnvName)
	if err != nil {
		t.Fatalf("LoadAddressBook: %v", err)
	}
	networks, err := config.LoadNetworks(results.EnvName)
	if err != nil {
		t.Fatalf("LoadNetworks: %v", err)
	}
	destNetwork, err := config.FindNetworkBySelector(networks, results.DestChainSelector)
	if err != nil {
		t.Fatalf("FindNetworkBySelector: %v", err)
	}
	client, err := ethclient.Dial(destNetwork.RPCs[0].HTTPURL)
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	defer client.Close()

	chainID, err := client.ChainID(ctx)
	if err != nil {
		t.Fatalf("ChainID: %v", err)
	}
	t.Logf("Connected to dest chain (chainID=%d)", chainID)

	// ── Find OffRamp for source chain ──
	destAddrs, err := addressBook.AddressesForChain(results.DestChainSelector)
	if err != nil {
		t.Fatalf("AddressesForChain: %v", err)
	}
	offRampAddr, isV160 := findOffRamp(ctx, t, client, destAddrs, results.SourceChainSelector)
	t.Logf("OffRamp: %s (v1.6.0=%v)", offRampAddr.Hex(), isV160)

	// ── Estimate block range from run timestamps ──
	runStart, err := time.Parse(time.RFC3339, results.RunStarted)
	if err != nil {
		t.Fatalf("invalid run_started %q: %v", results.RunStarted, err)
	}
	runEnd, err := time.Parse(time.RFC3339, results.RunEnded)
	if err != nil {
		t.Fatalf("invalid run_ended %q: %v", results.RunEnded, err)
	}
	fromTime := runStart.Add(-5 * time.Minute)
	toTime := runEnd.Add(30 * time.Minute)

	latest, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		t.Fatalf("latest header fetch failed: %v", err)
	}
	if latest == nil {
		t.Fatal("latest header fetch returned nil")
	}
	latestNumber := latest.Number.Uint64()
	blockTime := 2.0
	if latestNumber > 1 {
		prev, err := client.HeaderByNumber(ctx, new(big.Int).Sub(latest.Number, big.NewInt(1)))
		if err == nil && prev != nil && prev.Time > 0 && latest.Time > prev.Time {
			blockTime = float64(latest.Time-prev.Time) / float64(latestNumber-prev.Number.Uint64())
		}
	}
	latestTime := time.Unix(int64(latest.Time), 0)
	fromBlock := uint64(0)
	if fromTime.Before(latestTime) {
		sec := latestTime.Sub(fromTime).Seconds()
		if blk := uint64(sec / blockTime); blk < latestNumber {
			fromBlock = latestNumber - blk
		}
	}
	toBlock := latestNumber
	if toTime.Before(latestTime) {
		sec := latestTime.Sub(toTime).Seconds()
		if blk := uint64(sec / blockTime); blk < latestNumber {
			toBlock = latestNumber - blk
		}
	}
	// Clamp range to 100k blocks
	if toBlock-fromBlock > 100_000 {
		fromBlock = toBlock - 100_000
	}
	t.Logf("Block range: %d → %d (%d blocks)", fromBlock, toBlock, toBlock-fromBlock)

	// ── Build message ID list ──
	messageIDs := make([][32]byte, len(results.Messages))
	seqNums := make([]uint64, len(results.Messages))
	for i, msg := range results.Messages {
		messageIDs[i] = common.HexToHash(msg.MessageID)
		seqNums[i] = parseUint64(msg.SequenceNumber)
	}

	// ── Query ExecutionStateChanged for all messages ──
	type execResult struct {
		Block     uint64
		Timestamp time.Time
		State     uint8
		TxHash    string
	}
	const (
		execStateUntouched  uint8 = 0
		execStateInProgress uint8 = 1
		execStateSuccess    uint8 = 2
		execStateFailure    uint8 = 3
	)
	execByMsgID := make(map[[32]byte]execResult, len(results.Messages))
	filterOpts := &bind.FilterOpts{Start: fromBlock, End: &toBlock, Context: ctx}

	if isV160 {
		filterer, err := offramp160.NewOffRampFilterer(offRampAddr, client)
		if err != nil {
			t.Fatalf("NewOffRampFilterer: %v", err)
		}
		iter, err := filterer.FilterExecutionStateChanged(
			filterOpts,
			[]uint64{results.SourceChainSelector},
			seqNums,
			messageIDs,
		)
		if err != nil {
			t.Fatalf("FilterExecutionStateChanged: %v", err)
		}
		defer iter.Close()
		for iter.Next() {
			e := iter.Event
			header, _ := client.HeaderByNumber(ctx, big.NewInt(int64(e.Raw.BlockNumber)))
			ts := time.Time{}
			if header != nil {
				ts = time.Unix(int64(header.Time), 0)
			}
			prev, ok := execByMsgID[e.MessageId]
			if !ok || e.Raw.BlockNumber >= prev.Block {
				execByMsgID[e.MessageId] = execResult{Block: e.Raw.BlockNumber, Timestamp: ts, State: e.State, TxHash: e.Raw.TxHash.Hex()}
			}
		}
	} else {
		filterer, err := evm2evmofframp.NewEVM2EVMOffRampFilterer(offRampAddr, client)
		if err != nil {
			t.Fatalf("NewEVM2EVMOffRampFilterer: %v", err)
		}
		iter, err := filterer.FilterExecutionStateChanged(
			filterOpts,
			seqNums,
			messageIDs,
		)
		if err != nil {
			t.Fatalf("FilterExecutionStateChanged: %v", err)
		}
		defer iter.Close()
		for iter.Next() {
			e := iter.Event
			header, _ := client.HeaderByNumber(ctx, big.NewInt(int64(e.Raw.BlockNumber)))
			ts := time.Time{}
			if header != nil {
				ts = time.Unix(int64(header.Time), 0)
			}
			prev, ok := execByMsgID[e.MessageId]
			if !ok || e.Raw.BlockNumber >= prev.Block {
				execByMsgID[e.MessageId] = execResult{Block: e.Raw.BlockNumber, Timestamp: ts, State: e.State, TxHash: e.Raw.TxHash.Hex()}
			}
		}
	}

	// ── Print table ──
	fmt.Println()
	fmt.Println("=== ExecutionStateChanged on Dest Chain ===")
	fmt.Printf("%-22s  %6s  %10s  %10s  %-15s  %s\n", "MessageID", "Seq", "Block", "Time", "Status", "TxHash")
	fmt.Println(strings.Repeat("─", 170))

	executed := 0
	for _, msg := range results.Messages {
		msgID := common.HexToHash(msg.MessageID)
		seq := parseUint64(msg.SequenceNumber)
		short := strings.TrimPrefix(msg.MessageID, "0x")
		if len(short) > 20 {
			short = short[:20]
		}

		r, found := execByMsgID[msgID]
		if found {
			executed++
			status := fmt.Sprintf("🔶 state=%d", r.State)
			switch r.State {
			case execStateSuccess:
				status = "✅ success"
			case execStateFailure:
				status = "⚠️  failed"
			case execStateInProgress:
				status = "⏳ in_progress"
			case execStateUntouched:
				status = "⭕ untouched"
			}
			fmt.Printf("%-22s  %6d  %10d  %10s  %-15s  %s\n",
				short, seq, r.Block, r.Timestamp.Format("15:04:05"), status, r.TxHash)
			if url := txExplorerURL(chainID.Uint64(), r.TxHash); url != "" {
				fmt.Printf("%55sExplorer: %s\n", "", url)
			}
		} else {
			fmt.Printf("%-22s  %6d  %10s  %10s  %-15s  %s\n",
				short, seq, "-", "-", "❌ not found", "-")
		}
	}

	fmt.Println(strings.Repeat("─", 170))
	fmt.Printf("Summary: %d/%d executed\n", executed, len(results.Messages))
}

// findOffRamp finds the OffRamp contract handling the given source chain.
// Returns the address and whether it's v1.6.0 (true) or v1.5.0 (false).
func findOffRamp(
	ctx context.Context,
	t *testing.T,
	client *ethclient.Client,
	chainAddrs map[string]cldf.TypeAndVersion,
	sourceChainSelector uint64,
) (common.Address, bool) {
	// Try v1.6.0 OffRamp first (single contract, handles all source chains)
	for addrStr, tv := range chainAddrs {
		if string(tv.Type) != "OffRamp" {
			continue
		}
		addr := common.HexToAddress(addrStr)
		caller, err := offramp160.NewOffRampCaller(addr, client)
		if err != nil {
			continue
		}
		if _, err := caller.GetSourceChainConfig(&bind.CallOpts{Context: ctx}, sourceChainSelector); err == nil {
			return addr, true
		}
	}

	// Fall back to v1.5.0 EVM2EVMOffRamp (one per source chain)
	for addrStr, tv := range chainAddrs {
		if string(tv.Type) != "EVM2EVMOffRamp" {
			continue
		}
		addr := common.HexToAddress(addrStr)
		caller, err := evm2evmofframp.NewEVM2EVMOffRampCaller(addr, client)
		if err != nil {
			continue
		}
		sc, err := caller.GetStaticConfig(&bind.CallOpts{Context: ctx})
		if err != nil {
			continue
		}
		if sc.SourceChainSelector == sourceChainSelector {
			return addr, false
		}
	}

	t.Fatalf("No OffRamp found for source chain %d", sourceChainSelector)
	return common.Address{}, false
}

func parseUint64(s string) uint64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	var v uint64
	if strings.HasPrefix(s, "0x") || strings.HasPrefix(s, "0X") {
		fmt.Sscanf(s, "%x", &v)
	} else {
		fmt.Sscanf(s, "%d", &v)
	}
	return v
}

func txExplorerURL(chainID uint64, txHash string) string {
	switch chainID {
	case 84532:
		return "https://sepolia.basescan.org/tx/" + txHash
	case 8453:
		return "https://basescan.org/tx/" + txHash
	case 11155111:
		return "https://sepolia.etherscan.io/tx/" + txHash
	case 1:
		return "https://etherscan.io/tx/" + txHash
	default:
		return ""
	}
}
