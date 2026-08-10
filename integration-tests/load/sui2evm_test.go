package load

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-testing-framework/wasp"

	bindutils "github.com/smartcontractkit/chainlink-sui/bindings/utils"
	"github.com/smartcontractkit/chainlink-sui/integration-tests/load/config"
	"github.com/smartcontractkit/chainlink-sui/integration-tests/load/sui"
	"github.com/smartcontractkit/chainlink-sui/integration-tests/load/wallet"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
)

var runName = flag.String("run-name", "", "Name of the run config file (without .toml) in runs/ directory")

func TestSui2EVM(t *testing.T) {
	if *runName == "" {
		t.Fatal("--run-name flag is required (e.g., --run-name my-first-sui-to-evm-run)")
	}
	if strings.Contains(*runName, "evm-to-sui") {
		t.Fatalf("--run-name %q is an EVM→Sui config; run TestEVM2Sui instead", *runName)
	}

	ctx := context.Background()

	// Load full config (all 4 layers)
	cfg, err := config.LoadFullConfig(*runName)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Find source chain RPC URL from network config
	srcNetwork, err := config.FindNetworkBySelector(cfg.Networks, cfg.SourceChainSelector)
	if err != nil {
		t.Fatalf("Source chain not found in network config: %v", err)
	}
	if len(srcNetwork.RPCs) == 0 {
		t.Fatal("Source chain has no RPC endpoints")
	}
	suiRPCURL := srcNetwork.RPCs[0].HTTPURL

	// Create Sui client and signer
	ptbClient, err := sui.NewSuiClient(t, suiRPCURL)
	if err != nil {
		t.Fatalf("Failed to create Sui client: %v", err)
	}

	suiSigner, senderAddress, err := sui.NewSuiSigner(cfg.SuiPrivateKey)
	if err != nil {
		t.Fatalf("Failed to create Sui signer: %v", err)
	}

	slog.Info("Sui signer ready", "address", senderAddress)

	// Resolve Sui addresses from address book
	suiAddrs, err := config.ResolveAddressesForChain(cfg.AddressBook, cfg.SourceChainSelector)
	if err != nil {
		t.Fatalf("Failed to resolve Sui addresses: %v", err)
	}

	ccipPkgID, ok := suiAddrs["SuiCCIP"]
	if !ok {
		t.Fatal("SuiCCIP address not found in address book")
	}
	onRampPkgID, ok := suiAddrs["SuiOnRamp"]
	if !ok {
		t.Fatal("SuiOnRamp address not found in address book")
	}
	ccipObjectRefID, ok := suiAddrs["SuiCCIPObjectRef"]
	if !ok {
		t.Fatal("SuiCCIPObjectRef not found in address book")
	}
	onRampStateID, ok := suiAddrs["SuiOnRampStateObjectID"]
	if !ok {
		t.Fatal("SuiOnRampStateObjectID not found in address book")
	}
	linkTokenMetadataID, ok := suiAddrs["SuiLinkTokenObjectMetadataID"]
	if !ok {
		t.Fatal("SuiLinkTokenObjectMetadataID not found in address book")
	}

	// Fee token: use native SUI for simplicity.
	// CoinMetadata<0x2::sui::SUI> is a well-known system object.
	const feeTokenType = "0x2::sui::SUI"
	const suiCoinMetadataID = "0x587c29de216efd4219573e08a1f6964d4fa7cb714518c2c8a0f29abfa264327d"

	receiverBytes, err := sui.ParseEVMReceiver32(cfg.ReceiverAddress)
	if err != nil {
		t.Fatalf("Invalid EVM receiver address %q: %v", cfg.ReceiverAddress, err)
	}

	// Determine mode: message-only vs token transfer
	isTokenMode := cfg.TokenConfig != nil
	var tokenCoinPool *sui.TokenCoinPool
	var tokenPoolConfig *sui.TokenPoolConfig
	var tokenCoinType string
	var tokenPoolPkgID string
	var tokenPoolStateID string
	var denyListObjectID string
	var tokenStateObjectID string

	if isTokenMode {
		tokenPoolConfig, err = sui.ResolveTokenPoolConfig(
			ctx,
			ptbClient,
			suiSigner,
			ccipPkgID,
			ccipObjectRefID,
			cfg.TokenConfig.TokenIdentifier,
		)
		if err != nil {
			t.Fatalf("Failed to resolve token pool config: %v", err)
		}

		tokenCoinType = tokenPoolConfig.CoinType
		tokenPoolPkgID = tokenPoolConfig.LatestPoolPackageID
		tokenPoolStateID = tokenPoolConfig.PoolStateObjectID
		// DenyList and TokenState object IDs are derived from lock_or_burn_params for managed pools.
		// managed_token_pool lock_or_burn params: [clock, denyList, managedTokenState, poolState]
		// The registry returns the shared object IDs in that order.
		if tokenPoolConfig.PoolKind == "managed" && len(tokenPoolConfig.LockOrBurnParams) > 3 {
			denyListObjectID = tokenPoolConfig.LockOrBurnParams[1]
			tokenStateObjectID = tokenPoolConfig.LockOrBurnParams[2]
		}

		slog.Info("Resolved token pool config",
			"poolPackage", tokenPoolPkgID,
			"poolModule", tokenPoolConfig.PoolModule,
			"poolKind", tokenPoolConfig.PoolKind,
			"poolState", tokenPoolStateID,
			"coinType", tokenCoinType,
		)

		tokenCoinPool, err = sui.PrepareTokenCoinPool(
			ctx,
			ptbClient,
			suiSigner,
			senderAddress,
			tokenCoinType,
			cfg.MessageCount,
			cfg.TokenConfig.Amount,
		)
		if err != nil {
			t.Fatalf("Failed to prepare token coin pool: %v", err)
		}
	}

	estimatedFee, err := sui.EstimateSuiToEVMFee(
		ctx,
		ptbClient,
		suiSigner,
		ccipPkgID,
		onRampPkgID,
		ccipObjectRefID,
		onRampStateID,
		feeTokenType,
		suiCoinMetadataID,
		cfg.DestChainSelector,
		receiverBytes,
		cfg.MessageData,
		cfg.EvmCallbackGasLimit,
	)
	if err != nil {
		t.Fatalf("Failed to estimate Sui fee: %v", err)
	}
	splitAmountPerCoin := sui.RecommendedSplitAmountPerCoin(estimatedFee)
	slog.Info("Using dynamic SUI split amount",
		"estimatedFee", estimatedFee,
		"splitAmountPerCoin", splitAmountPerCoin,
	)

	coinPool, err := sui.PrepareSuiCoinPool(
		ctx,
		ptbClient,
		suiSigner,
		senderAddress,
		cfg.MessageCount,
		splitAmountPerCoin,
	)
	if err != nil {
		t.Fatalf("Failed to prepare SUI coin pool: %v", err)
	}

	slog.Info("Resolved Sui addresses",
		"ccip", ccipPkgID,
		"onramp", onRampPkgID,
		"ccipObjectRef", ccipObjectRefID,
		"onRampState", onRampStateID,
		"linkTokenMetadata", linkTokenMetadataID,
		"coinPoolSize", coinPool.Size(),
		"tokenCoinPoolSize", func() int {
			if tokenCoinPool == nil {
				return 0
			}
			return tokenCoinPool.Size()
		}(),
	)

	// Token transfers stay sequential; message-only uses WASP wallet pool.
	if isTokenMode {
		runSui2EVMSequential(t, cfg, ptbClient, suiSigner, coinPool, tokenCoinPool,
			ccipPkgID, onRampPkgID, ccipObjectRefID, onRampStateID,
			feeTokenType, suiCoinMetadataID,
			tokenCoinType, tokenPoolPkgID, tokenPoolStateID, denyListObjectID, tokenStateObjectID,
			receiverBytes)
		return
	}

	runSui2EVMWASP(t, cfg, ptbClient, suiSigner, senderAddress,
		ccipPkgID, onRampPkgID, ccipObjectRefID, onRampStateID,
		feeTokenType, suiCoinMetadataID,
		receiverBytes)
}

// runSui2EVMSequential sends token-transfer messages sequentially using the main signer.
func runSui2EVMSequential(
	t *testing.T,
	cfg *config.LoadTestConfig,
	ptbClient *client.PTBClient,
	suiSigner bindutils.SuiSigner,
	coinPool *sui.SuiCoinPool,
	tokenCoinPool *sui.TokenCoinPool,
	ccipPkgID, onRampPkgID, ccipObjectRefID, onRampStateID string,
	feeTokenType, suiCoinMetadataID string,
	tokenCoinType, tokenPoolPkgID, tokenPoolStateID, denyListObjectID, tokenStateObjectID string,
	receiverBytes []byte,
) {
	ctx := context.Background()

	results := &config.RunResults{
		RunName:             cfg.RunName,
		EnvName:             cfg.EnvName,
		SourceChainSelector: cfg.SourceChainSelector,
		DestChainSelector:   cfg.DestChainSelector,
		TotalMessages:       cfg.MessageCount,
		RunStarted:          time.Now().Format(time.RFC3339),
		Messages:            make([]config.SentMessage, 0, cfg.MessageCount),
	}

	saveResults := func() {
		results.RunEnded = time.Now().Format(time.RFC3339)
		if err := config.SaveResults(results); err != nil {
			slog.Error("Failed to save results", "err", err)
		}
	}
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalCh
		slog.Warn("Interrupt received, saving results...")
		saveResults()
		os.Exit(1)
	}()
	defer saveResults()

	for i := 0; i < cfg.MessageCount; i++ {
		slog.Info("Sending message", "progress", fmt.Sprintf("%d/%d", i+1, cfg.MessageCount))

		gasCoinID, err := coinPool.Pop(ctx)
		if err != nil {
			t.Fatalf("Failed to get gas coin from pool for message %d: %v", i+1, err)
		}
		feeCoinID, err := coinPool.Pop(ctx)
		if err != nil {
			t.Fatalf("Failed to get fee coin from pool for message %d: %v", i+1, err)
		}
		if gasCoinID == feeCoinID {
			t.Fatalf("Invalid coin pool state: gas and fee coin IDs are equal for message %d (%s)", i+1, gasCoinID)
		}

		msg := config.SentMessage{
			SourceChainSelector: cfg.SourceChainSelector,
			DestChainSelector:   cfg.DestChainSelector,
			Timestamp:           time.Now().Format(time.RFC3339),
		}

		tokenCoinID, err := tokenCoinPool.Pop(ctx)
		if err != nil {
			t.Fatalf("Failed to get token coin from pool for message %d: %v", i+1, err)
		}
		msg.TokenAmount = fmt.Sprintf("%d", cfg.TokenConfig.Amount)
		msg.TokenIdentifier = cfg.TokenConfig.TokenIdentifier

		messageID, txDigest, seqNum, sendErr := sui.SendTokenMessage(
			ctx,
			ptbClient,
			suiSigner,
			ccipPkgID,
			onRampPkgID,
			ccipObjectRefID,
			onRampStateID,
			gasCoinID,
			feeTokenType,
			suiCoinMetadataID,
			feeCoinID,
			tokenCoinID,
			cfg.TokenConfig.TokenIdentifier,
			tokenCoinType,
			tokenPoolPkgID,
			tokenPoolStateID,
			denyListObjectID,
			tokenStateObjectID,
			cfg.DestChainSelector,
			receiverBytes,
			cfg.MessageData,
			cfg.SuiGasBudget,
			cfg.EvmCallbackGasLimit,
		)

		if sendErr != nil {
			slog.Error("Failed to send message",
				"progress", fmt.Sprintf("%d/%d", i+1, cfg.MessageCount),
				"error", sendErr,
			)
			msg.Success = false
			msg.Error = sendErr.Error()
			results.FailedMessages++
		} else {
			slog.Info("Message sent",
				"progress", fmt.Sprintf("%d/%d", i+1, cfg.MessageCount),
				"messageID", messageID,
				"txDigest", txDigest,
				"seqNum", seqNum,
				"tokenAmount", msg.TokenAmount,
				"tokenIdentifier", msg.TokenIdentifier,
			)
			msg.Success = true
			msg.MessageID = messageID
			msg.TransactionHash = txDigest
			msg.SequenceNumber = fmt.Sprintf("%d", seqNum)
			results.SuccessfulMessages++
		}

		results.Messages = append(results.Messages, msg)
	}

	slog.Info("Run complete",
		"successful", results.SuccessfulMessages,
		"failed", results.FailedMessages,
		"total", results.TotalMessages,
	)

	if results.FailedMessages > 0 {
		t.Fatalf("Run completed with failed messages: failed=%d total=%d", results.FailedMessages, results.TotalMessages)
	}
}

// runSui2EVMWASP sends message-only CCIP messages using WASP with N parallel wallets.
func runSui2EVMWASP(
	t *testing.T,
	cfg *config.LoadTestConfig,
	ptbClient *client.PTBClient,
	mainSigner bindutils.SuiSigner,
	mainAddress string,
	ccipPkgID, onRampPkgID, ccipObjectRefID, onRampStateID string,
	feeTokenType, suiCoinMetadataID string,
	receiverBytes []byte,
) {
	ctx := context.Background()
	N := cfg.LoadWallets

	// Step 1: Generate N wallets.
	wallets, err := wallet.GenerateSuiWallets(N)
	if err != nil {
		t.Fatalf("Failed to generate Sui wallets: %v", err)
	}
	t.Cleanup(func() {
		if err := wallet.SweepSuiWallets(ctx, ptbClient, wallets, mainAddress); err != nil {
			slog.Error("Failed to sweep Sui wallets during cleanup", "error", err)
		}
	})

	// Step 2: Estimate fee and compute split amount before funding.
	estimatedFee, err := sui.EstimateSuiToEVMFee(
		ctx, ptbClient, mainSigner,
		ccipPkgID, onRampPkgID, ccipObjectRefID, onRampStateID,
		feeTokenType, suiCoinMetadataID,
		cfg.DestChainSelector, receiverBytes, cfg.MessageData, cfg.EvmCallbackGasLimit,
	)
	if err != nil {
		t.Fatalf("Failed to estimate Sui fee: %v", err)
	}
	splitAmount := sui.RecommendedSplitAmountPerCoin(estimatedFee)

	// Step 3: Fund each wallet using the actual split amount requirement.
	fundingPerWallet := estimateSuiFunding(cfg, N, splitAmount)
	if err := wallet.FundSuiWallets(ctx, ptbClient, mainSigner, mainAddress, wallets, fundingPerWallet); err != nil {
		t.Fatalf("Failed to fund Sui wallets: %v", err)
	}

	// Step 4: Prepare all gas coin pools BEFORE creating any generators.
	// CRITICAL: wasp.NewGenerator starts a context.WithTimeout at construction time,
	// not at Run() time. Any slow work (like Sui RPC calls for coin splitting) done
	// between NewGenerator and p.Run() eats into the generator's lifetime.
	// With 5 wallets × ~3s per PrepareSuiCoinPool = ~15s of sequential prep,
	// a 2s duration would expire before p.Run() is even called for early wallets.
	msgPerWallet := cfg.MessageCount / N
	if msgPerWallet < 1 {
		msgPerWallet = 1
	}
	type walletPool struct {
		w       *wallet.Wallet
		gasPool *sui.SuiCoinPool
	}
	pools := make([]walletPool, N)
	for i := 0; i < N; i++ {
		w := wallets[i]
		gasPool, err := sui.PrepareSuiCoinPool(ctx, ptbClient, w.SuiSigner, w.Address, msgPerWallet, splitAmount)
		if err != nil {
			t.Fatalf("Failed to prepare SUI coin pool for wallet %d: %v", i, err)
		}
		pools[i] = walletPool{w: w, gasPool: gasPool}
	}

	// Step 5: Create all generators and run immediately.
	// Generators are created in a tight loop right before p.Run() to minimize
	// the gap between NewGenerator (context start) and Run (actual execution).
	resultsCh := make(chan config.SentMessage, cfg.MessageCount)
	p := wasp.NewProfile()
	rpsPerWallet := int64(cfg.LoadRPS / N)
	if rpsPerWallet < 1 {
		rpsPerWallet = 1
	}
	// Each wallet sends msgPerWallet messages at rpsPerWallet RPS.
	// duration = msgPerWallet / rpsPerWallet seconds.
	duration := time.Duration(msgPerWallet) * time.Second / time.Duration(rpsPerWallet)

	for i := 0; i < N; i++ {
		pool := pools[i]
		gun := &Sui2EVMMsgGun{
			wallet:            pool.w,
			gasPool:           pool.gasPool,
			ptbClient:         ptbClient,
			ccipPkgID:         ccipPkgID,
			onRampPkgID:       onRampPkgID,
			ccipObjectRefID:   ccipObjectRefID,
			onRampStateID:     onRampStateID,
			feeTokenType:      feeTokenType,
			suiCoinMetaID:     suiCoinMetadataID,
			destChainSelector: cfg.DestChainSelector,
			receiver:          receiverBytes,
			data:              cfg.MessageData,
			gasBudget:         cfg.SuiGasBudget,
			evmCallbackGas:    cfg.EvmCallbackGasLimit,
			resultsCh:         resultsCh,
		}

		gen, err := wasp.NewGenerator(&wasp.Config{
			GenName:  fmt.Sprintf("sui2evm-w%d", i),
			LoadType: wasp.RPS,
			Schedule: wasp.Plain(rpsPerWallet, duration),
			Gun:      gun,
		})
		p.Add(gen, err)
	}

	// Step 6: Run all generators in parallel.
	if _, err := p.Run(true); err != nil {
		t.Fatalf("WASP profile run failed: %v", err)
	}
	close(resultsCh)

	// Step 7: Collect and save results.
	results := collectResults(cfg, resultsCh)
	if err := config.SaveResults(results); err != nil {
		t.Fatalf("Failed to save results: %v", err)
	}

	if results.FailedMessages > 0 {
		t.Fatalf("Run completed with failed messages: failed=%d total=%d", results.FailedMessages, results.TotalMessages)
	}
}

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}
