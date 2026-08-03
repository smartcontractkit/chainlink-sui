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

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_3/message_hasher"

	"github.com/smartcontractkit/chainlink-sui/integration-tests/load/config"
	"github.com/smartcontractkit/chainlink-sui/integration-tests/load/evm"
)

func TestEVM2Sui(t *testing.T) {
	if *runName == "" {
		t.Fatal("--run-name flag is required (e.g., --run-name my-first-evm-to-sui-run)")
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

	extraArgs, err := evm.SerializeClientSUIExtraArgsV1(suiExtraArgs)
	if err != nil {
		t.Fatalf("Failed to serialize SuiExtraArgsV1: %v", err)
	}

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

	// Ensure results are saved even on panic
	defer func() {
		results.RunEnded = time.Now().Format(time.RFC3339)
		if err := config.SaveResults(results); err != nil {
			t.Errorf("Failed to save results: %v", err)
		}
	}()

	// Send messages sequentially
	for i := 0; i < cfg.MessageCount; i++ {
		slog.Info("Sending message", "progress", fmt.Sprintf("%d/%d", i+1, cfg.MessageCount))

		msg := config.SentMessage{
			SourceChainSelector: cfg.SourceChainSelector,
			DestChainSelector:   cfg.DestChainSelector,
			Timestamp:           time.Now().Format(time.RFC3339),
		}

		messageID, txHash, sendErr := evm.SendMessage(
			ctx,
			ethClient,
			auth,
			routerAddress,
			cfg.DestChainSelector,
			receiver[:],
			cfg.MessageData,
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
