package load

import (
	"context"
	"encoding/hex"
	"fmt"
	"log/slog"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_3/message_hasher"
	"github.com/smartcontractkit/chainlink-testing-framework/wasp"

	"github.com/smartcontractkit/chainlink-sui/integration-tests/load/config"
	"github.com/smartcontractkit/chainlink-sui/integration-tests/load/evm"
	"github.com/smartcontractkit/chainlink-sui/integration-tests/load/sui"
	"github.com/smartcontractkit/chainlink-sui/integration-tests/load/wallet"
)

func TestEVM2Sui(t *testing.T) {
	if *runName == "" {
		t.Fatal("--run-name flag is required (e.g., --run-name my-first-evm-to-sui-run)")
	}
	if strings.Contains(*runName, "sui-to-evm") {
		t.Fatalf("--run-name %q is a Sui→EVM config; run TestSui2EVM instead", *runName)
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
	evmRPCURL := srcNetwork.RPCs[0].HTTPURL

	// Create EVM client and signer
	ethClient, err := evm.NewEVMClient(evmRPCURL)
	if err != nil {
		t.Fatalf("Failed to create EVM client: %v", err)
	}
	defer ethClient.Close()

	chainID, err := ethClient.ChainID(ctx)
	if err != nil {
		t.Fatalf("Failed to get chain ID: %v", err)
	}

	auth, err := evm.NewEVMSigner(cfg.EVMPrivateKey, chainID)
	if err != nil {
		t.Fatalf("Failed to create EVM signer: %v", err)
	}

	slog.Info("EVM signer ready", "address", auth.From.Hex(), "chainID", chainID)

	// Resolve EVM router address from address book
	evmAddrs, err := config.ResolveAddressesForChain(cfg.AddressBook, cfg.SourceChainSelector)
	if err != nil {
		t.Fatalf("Failed to resolve EVM addresses: %v", err)
	}

	routerAddrStr, ok := evmAddrs["Router"]
	if !ok {
		t.Fatal("Router address not found in address book for source chain")
	}
	routerAddress := common.HexToAddress(routerAddrStr)

	slog.Info("Resolved EVM addresses", "router", routerAddress.Hex())

	// Parse receiver address (Sui package ID, 32 bytes)
	receiverBytes, err := hex.DecodeString(strings.TrimPrefix(cfg.ReceiverAddress, "0x"))
	if err != nil {
		t.Fatalf("Failed to parse receiver address: %v", err)
	}
	var receiver [32]byte
	copy(receiver[:], receiverBytes)

	// Determine mode: message-only vs token transfer
	isTokenMode := cfg.TokenConfig != nil
	var tokenAddress common.Address
	var tokenAmount *big.Int
	var totalTokenAmount *big.Int
	var suiReceiverStateID [32]byte

	if isTokenMode {
		tokenAddress = common.HexToAddress(cfg.TokenConfig.TokenIdentifier)
		tokenAmount = new(big.Int).SetUint64(cfg.TokenConfig.Amount)
		totalTokenAmount = new(big.Int).Mul(tokenAmount, big.NewInt(int64(cfg.MessageCount)))

		// For token-only mode, tokenReceiver is the Sui EOA (signer address converted to bytes32).
		// For token+message mode, tokenReceiver is the receiver state object ID.
		if cfg.TokenConfig.Mode == "token_and_message" {
			if cfg.SuiReceiverConfig == nil || cfg.SuiReceiverConfig.PackageID == "" {
				t.Fatal("sui_receiver.package_id is required for EVM→Sui token_and_message mode")
			}
			// Resolve Sui destination RPC from network config
			destNetwork, err := config.FindNetworkBySelector(cfg.Networks, cfg.DestChainSelector)
			if err != nil {
				t.Fatalf("Destination chain not found in network config: %v", err)
			}
			if len(destNetwork.RPCs) == 0 {
				t.Fatal("Destination chain has no RPC endpoints")
			}
			suiDestRPC := destNetwork.RPCs[0]
			grpcToken := cfg.SuiGrpcToken
			if grpcToken == "" {
				grpcToken = suiDestRPC.GrpcToken
			}
			suiDestClient, err := sui.NewSuiClient(t, suiDestRPC.HTTPURL, suiDestRPC.GrpcTarget, grpcToken)
			if err != nil {
				t.Fatalf("Failed to create Sui destination client: %v", err)
			}
			defer suiDestClient.Close()

			receiverStateStr, err := sui.ResolveReceiverState(ctx, suiDestClient, cfg.SuiReceiverConfig.PackageID)
			if err != nil {
				t.Fatalf("Failed to resolve Sui receiver state: %v", err)
			}
			suiReceiverStateID, err = sui.SuiObjectIdToBytes32(receiverStateStr)
			if err != nil {
				t.Fatalf("Failed to convert receiver state ID to bytes32: %v", err)
			}
		}

		// Approve Router for total token amount once at start
		err = evm.ApproveRouterForTokens(
			ctx,
			ethClient,
			auth,
			tokenAddress,
			routerAddress,
			totalTokenAmount,
		)
		if err != nil {
			t.Fatalf("Failed to approve Router for tokens: %v", err)
		}
	}

	// Construct SuiExtraArgsV1
	// Clock object ID is always 0x6 on Sui
	clockObjID := [32]byte{}
	copy(clockObjID[:], common.FromHex("0x0000000000000000000000000000000000000000000000000000000000000006"))

	suiExtraArgs := message_hasher.ClientSuiExtraArgsV1{
		GasLimit:                 big.NewInt(int64(cfg.EvmGasLimit)),
		AllowOutOfOrderExecution: true,
		TokenReceiver:            [32]byte{}, // zero for message-only
		ReceiverObjectIds:        [][32]byte{clockObjID},
	}

	if isTokenMode {
		if cfg.TokenConfig.Mode == "token_only" {
			// Token-only: tokenReceiver is the Sui EOA (signer address as bytes32)
			signerAddressBytes, err := sui.SuiObjectIdToBytes32(auth.From.Hex())
			if err != nil {
				t.Fatalf("Failed to convert EVM signer address to Sui token receiver bytes: %v", err)
			}
			suiExtraArgs.TokenReceiver = signerAddressBytes
		} else {
			// token+message: tokenReceiver is the receiver state object ID
			suiExtraArgs.TokenReceiver = suiReceiverStateID
			suiExtraArgs.ReceiverObjectIds = [][32]byte{clockObjID, suiReceiverStateID}
		}
	}

	extraArgs, err := evm.SerializeClientSUIExtraArgsV1(suiExtraArgs)
	if err != nil {
		t.Fatalf("Failed to serialize SuiExtraArgsV1: %v", err)
	}

	// Token transfers stay sequential; message-only uses WASP wallet pool.
	if isTokenMode {
		runEVM2SuiSequential(t, cfg, ethClient, auth, routerAddress,
			receiver, tokenAddress, tokenAmount, extraArgs)
		return
	}

	runEVM2SuiWASP(t, cfg, ethClient, auth, routerAddress, receiver, extraArgs)
}

// runEVM2SuiSequential sends token-transfer messages sequentially using the main signer.
func runEVM2SuiSequential(
	t *testing.T,
	cfg *config.LoadTestConfig,
	ethClient *ethclient.Client,
	auth *bind.TransactOpts,
	routerAddress common.Address,
	receiver [32]byte,
	tokenAddress common.Address,
	tokenAmount *big.Int,
	extraArgs []byte,
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

	defer func() {
		results.RunEnded = time.Now().Format(time.RFC3339)
		if err := config.SaveResults(results); err != nil {
			t.Errorf("Failed to save results: %v", err)
		}
	}()

	for i := 0; i < cfg.MessageCount; i++ {
		slog.Info("Sending message", "progress", fmt.Sprintf("%d/%d", i+1, cfg.MessageCount))

		msg := config.SentMessage{
			SourceChainSelector: cfg.SourceChainSelector,
			DestChainSelector:   cfg.DestChainSelector,
			Timestamp:           time.Now().Format(time.RFC3339),
		}

		msg.TokenAmount = fmt.Sprintf("%d", cfg.TokenConfig.Amount)
		msg.TokenIdentifier = cfg.TokenConfig.TokenIdentifier
		messageID, txHash, sendErr := evm.SendTokenMessage(
			ctx,
			ethClient,
			auth,
			routerAddress,
			cfg.DestChainSelector,
			receiver[:],
			cfg.MessageData,
			tokenAddress,
			tokenAmount,
			extraArgs,
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
				"txHash", txHash,
				"tokenAmount", msg.TokenAmount,
				"tokenIdentifier", msg.TokenIdentifier,
			)
			msg.Success = true
			msg.MessageID = messageID
			msg.TransactionHash = txHash
			results.SuccessfulMessages++
		}

		results.Messages = append(results.Messages, msg)
	}

	slog.Info("Run complete",
		"successful", results.SuccessfulMessages,
		"failed", results.FailedMessages,
		"total", results.TotalMessages,
	)
}

// runEVM2SuiWASP sends message-only CCIP messages using WASP with N parallel wallets.
func runEVM2SuiWASP(
	t *testing.T,
	cfg *config.LoadTestConfig,
	ethClient *ethclient.Client,
	mainAuth *bind.TransactOpts,
	routerAddress common.Address,
	receiver [32]byte,
	extraArgs []byte,
) {
	ctx := context.Background()
	N := cfg.LoadWallets

	// Step 1: Generate N wallets.
	chainID, err := ethClient.ChainID(ctx)
	if err != nil {
		t.Fatalf("Failed to get chain ID: %v", err)
	}
	seed, err := wallet.ParseSeed(cfg.WalletSeed)
	if err != nil {
		t.Fatalf("Failed to parse wallet seed: %v", err)
	}
	wallets, err := wallet.GenerateEVMWallets(N, chainID, seed)
	if err != nil {
		t.Fatalf("Failed to generate EVM wallets: %v", err)
	}
	t.Cleanup(func() {
		if err := wallet.SweepEVMWallets(ctx, ethClient, wallets, mainAuth.From); err != nil {
			slog.Error("Failed to sweep EVM wallets during cleanup", "error", err)
		}
	})

	// Step 2: Fund each wallet.
	fundingPerWallet := estimateEvmFunding(cfg, N)
	if err := wallet.FundEVMWallets(ctx, ethClient, mainAuth, wallets, fundingPerWallet); err != nil {
		t.Fatalf("Failed to fund EVM wallets: %v", err)
	}

	// Step 3: Create N generators (one per wallet).
	resultsCh := make(chan config.SentMessage, cfg.MessageCount)
	p := wasp.NewProfile()
	rpsPerWallet := int64(cfg.LoadRPS / N)
	if rpsPerWallet < 1 {
		rpsPerWallet = 1
	}
	msgPerWallet := cfg.MessageCount / N
	if msgPerWallet < 1 {
		msgPerWallet = 1
	}
	// Each wallet sends msgPerWallet messages at rpsPerWallet RPS.
	// duration = msgPerWallet / rpsPerWallet seconds.
	duration := time.Duration(msgPerWallet) * time.Second / time.Duration(rpsPerWallet)

	for i := 0; i < N; i++ {
		w := wallets[i]

		gun := &EVM2SuiMsgGun{
			wallet:            w,
			ethClient:         ethClient,
			routerAddress:     routerAddress,
			destChainSelector: cfg.DestChainSelector,
			receiver:          receiver,
			data:              cfg.MessageData,
			extraArgs:         extraArgs,
			resultsCh:         resultsCh,
		}

		gen, err := wasp.NewGenerator(&wasp.Config{
			GenName:  fmt.Sprintf("evm2sui-w%d", i),
			LoadType: wasp.RPS,
			Schedule: wasp.Plain(rpsPerWallet, duration),
			Gun:      gun,
		})
		p.Add(gen, err)
	}

	// Step 4: Run all generators in parallel.
	if _, err := p.Run(true); err != nil {
		t.Fatalf("WASP profile run failed: %v", err)
	}
	close(resultsCh)

	// Step 5: Collect and save results.
	results := collectResults(cfg, resultsCh)
	if err := config.SaveResults(results); err != nil {
		t.Fatalf("Failed to save results: %v", err)
	}

	if results.FailedMessages > 0 {
		t.Fatalf("Run completed with failed messages: failed=%d total=%d", results.FailedMessages, results.TotalMessages)
	}
}
