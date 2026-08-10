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

	"github.com/smartcontractkit/chainlink-sui/integration-tests/load/config"
	"github.com/smartcontractkit/chainlink-sui/integration-tests/load/sui"
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

	// Prepare results
	results := &config.RunResults{
		RunName:             cfg.RunName,
		EnvName:             cfg.EnvName,
		SourceChainSelector: cfg.SourceChainSelector,
		DestChainSelector:   cfg.DestChainSelector,
		TotalMessages:       cfg.MessageCount,
		RunStarted:          time.Now().Format(time.RFC3339),
		Messages:            make([]config.SentMessage, 0, cfg.MessageCount),
	}

	// Save results on normal exit, panic, or interrupt (Ctrl+C).
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

	// Send messages sequentially
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

		var tokenCoinID string
		if isTokenMode {
			tokenCoinID, err = tokenCoinPool.Pop(ctx)
			if err != nil {
				t.Fatalf("Failed to get token coin from pool for message %d: %v", i+1, err)
			}
			msg.TokenAmount = fmt.Sprintf("%d", cfg.TokenConfig.Amount)
			msg.TokenIdentifier = cfg.TokenConfig.TokenIdentifier
		}

		var tokenIdentifier string
		if isTokenMode {
			tokenIdentifier = cfg.TokenConfig.TokenIdentifier
		}

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
			tokenIdentifier,
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

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}
