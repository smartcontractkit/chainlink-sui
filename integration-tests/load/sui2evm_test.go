package load

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"math"
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
	suiRPC := srcNetwork.RPCs[0]
	grpcToken := cfg.SuiGrpcToken
	if grpcToken == "" {
		grpcToken = suiRPC.GrpcToken
	}

	// Create Sui client and signer
	ptbClient, err := sui.NewSuiClient(t, suiRPC.HTTPURL, suiRPC.GrpcTarget, grpcToken)
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

	// Token transfers stay sequential; message-only uses WASP wallet pool.
	if isTokenMode {
		coinPool, prepErr := sui.PrepareSuiCoinPool(
			ctx,
			ptbClient,
			suiSigner,
			senderAddress,
			cfg.MessageCount,
			splitAmountPerCoin,
		)
		if prepErr != nil {
			t.Fatalf("Failed to prepare SUI coin pool: %v", prepErr)
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

		runSui2EVMSequential(t, cfg, ptbClient, suiSigner, coinPool, tokenCoinPool,
			ccipPkgID, onRampPkgID, ccipObjectRefID, onRampStateID,
			feeTokenType, suiCoinMetadataID,
			estimatedFee, splitAmountPerCoin,
			tokenCoinType, tokenPoolPkgID, tokenPoolStateID, denyListObjectID, tokenStateObjectID,
			receiverBytes)
		return
	}

	slog.Info("Resolved Sui addresses",
		"ccip", ccipPkgID,
		"onramp", onRampPkgID,
		"ccipObjectRef", ccipObjectRefID,
		"onRampState", onRampStateID,
		"linkTokenMetadata", linkTokenMetadataID,
	)

	runSui2EVMWASP(t, cfg, ptbClient, suiSigner, senderAddress,
		ccipPkgID, onRampPkgID, ccipObjectRefID, onRampStateID,
		feeTokenType, suiCoinMetadataID,
		estimatedFee, splitAmountPerCoin,
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
	estimatedFee, splitAmountPerCoin uint64,
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
			splitAmountPerCoin,
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

	_ = estimatedFee // reserved for run-level diagnostics parity with WASP mode

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
	estimatedFee, splitAmount uint64,
	receiverBytes []byte,
) {
	ctx := context.Background()
	N := cfg.LoadWallets

	// Step 1: Generate N wallets (deterministic when WALLET_SEED is set).
	seed, err := wallet.ParseSeed(cfg.WalletSeed)
	if err != nil {
		t.Fatalf("Failed to parse wallet seed: %v", err)
	}
	wallets, err := wallet.GenerateSuiWallets(N, seed)
	if err != nil {
		t.Fatalf("Failed to generate Sui wallets: %v", err)
	}

	// Step 2: Fund each wallet using the precomputed split amount requirement.
	fundingPerWallet := estimateSuiFunding(cfg, N, splitAmount)
	if err := wallet.FundSuiWallets(ctx, ptbClient, mainSigner, mainAddress, wallets, fundingPerWallet); err != nil {
		t.Fatalf("Failed to fund Sui wallets: %v", err)
	}
	t.Cleanup(func() {
		if err := wallet.SweepSuiWallets(ctx, ptbClient, wallets, mainAddress); err != nil {
			slog.Error("Failed to sweep Sui wallets during cleanup", "error", err)
		}
	})

	// Step 3: Prepare one gas coin + one reusable fee coin per wallet BEFORE creating generators.
	// CRITICAL: wasp.NewGenerator starts a context.WithTimeout at construction time,
	// not at Run() time. Any slow work (like Sui RPC calls for coin splitting) done
	// between NewGenerator and p.Run() eats into the generator's lifetime.
	// With 5 wallets × ~3s of sequential prep,
	// a 2s duration would expire before p.Run() is even called for early wallets.
	msgPerWalletBase := cfg.MessageCount / N
	msgRemainder := cfg.MessageCount % N
	msgCounts := make([]int, N)
	for i := 0; i < N; i++ {
		msgCounts[i] = msgPerWalletBase
		if i < msgRemainder {
			msgCounts[i]++
		}
		if msgCounts[i] < 1 {
			msgCounts[i] = 1
		}
	}

	targetPerWalletRPS := float64(cfg.LoadRPS) / float64(N)
	if targetPerWalletRPS <= 0 {
		t.Fatalf("invalid target per-wallet RPS %.4f from load.rps=%d wallets=%d", targetPerWalletRPS, cfg.LoadRPS, N)
	}
	minInterval := time.Duration(float64(time.Second) / targetPerWalletRPS)

	// WASP schedule requires integer RPS. We over-drive the scheduler and enforce
	// the exact fractional per-wallet rate in the gun via minCallInterval.
	scheduleRPS := int64(math.Ceil(targetPerWalletRPS))
	if scheduleRPS < 1 {
		scheduleRPS = 1
	}

	slog.Info("WASP rate configuration",
		"globalTargetRPS", cfg.LoadRPS,
		"wallets", N,
		"targetPerWalletRPS", targetPerWalletRPS,
		"scheduleRPSPerWallet", scheduleRPS,
		"minCallIntervalMs", minInterval.Milliseconds(),
	)
	type walletCoins struct {
		w         *wallet.Wallet
		gasCoinID string
		feeCoinID string
	}
	coins := make([]walletCoins, N)
	for i := 0; i < N; i++ {
		w := wallets[i]
		gasCoinID, feeCoinID, err := sui.SetupWalletCoins(ctx, ptbClient, w.SuiSigner, w.Address, msgCounts[i], cfg.SuiGasBudget, splitAmount)
		if err != nil {
			t.Fatalf("Failed to setup wallet coins for wallet %d: %v", i, err)
		}
		coins[i] = walletCoins{w: w, gasCoinID: gasCoinID, feeCoinID: feeCoinID}
	}

	// Step 4: Create all generators and run immediately.
	// Generators are created in a tight loop right before p.Run() to minimize
	// the gap between NewGenerator (context start) and Run (actual execution).
	resultsCh := make(chan config.SentMessage, cfg.MessageCount)
	p := wasp.NewProfile()

	for i := 0; i < N; i++ {
		walletCoins := coins[i]
		msgCount := msgCounts[i]
		durationSec := int64(math.Ceil(float64(msgCount) / targetPerWalletRPS))
		if durationSec < 1 {
			durationSec = 1
		}
		duration := time.Duration(durationSec) * time.Second

		gun := &Sui2EVMMsgGun{
			wallet:            walletCoins.w,
			gasCoinID:         walletCoins.gasCoinID,
			feeCoinID:         walletCoins.feeCoinID,
			feeAmount:         splitAmount,
			maxCalls:          msgCount,
			minCallInterval:   minInterval,
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
			GenName:     fmt.Sprintf("sui2evm-w%d", i),
			LoadType:    wasp.RPS,
			Schedule:    wasp.Plain(scheduleRPS, duration),
			CallTimeout: 30 * time.Second,
			Gun:         gun,
		})
		p.Add(gen, err)
	}

	// Step 5: Run all generators in parallel.
	if _, err := p.Run(true); err != nil {
		t.Fatalf("WASP profile run failed: %v", err)
	}

	_ = estimatedFee // retained for parity/debug parity in caller-level logging

	// Step 6: Collect and save results.
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
