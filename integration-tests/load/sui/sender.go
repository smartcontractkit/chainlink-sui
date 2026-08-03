package sui

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"
	"strconv"
	"strings"

	"github.com/block-vision/sui-go-sdk/transaction"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	onramp "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_onramp/onramp"
	bindutils "github.com/smartcontractkit/chainlink-sui/bindings/utils"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
)

// FetchFeeCoin finds a coin of the given type owned by the signer address.
// Returns the coin's object ID, or an error if no suitable coin is found.
func FetchFeeCoin(ctx context.Context, ptbClient *client.PTBClient, signerAddress string, coinType string) (string, error) {
	coins, err := ptbClient.QueryCoinsByAddress(ctx, signerAddress, coinType)
	if err != nil {
		return "", fmt.Errorf("failed to query coins for type %s: %w", coinType, err)
	}
	if len(coins) == 0 {
		return "", fmt.Errorf("no coins of type %s found for address %s", coinType, signerAddress)
	}

	// Use the coin with the largest balance
	var bestCoin string
	var bestBalance uint64
	for _, coin := range coins {
		bal := coin.GetBalance()
		if bal >= bestBalance {
			bestBalance = bal
			bestCoin = coin.GetObjectId()
		}
	}

	slog.Info("Selected fee coin",
		"coinId", bestCoin,
		"balance", bestBalance,
		"coinType", coinType,
	)

	return bestCoin, nil
}

// SendMessage builds and executes a CCIP message from Sui to EVM.
//
// ccipObjectRefID and onRampStateID are shared object IDs resolved from the address book.
// feeTokenType is the full Move coin type (e.g. "0x2::sui::SUI" for native SUI).
// feeTokenMetadataID is the CoinMetadata object ID for the fee token.
// feeTokenCoinID is a coin object owned by the signer, used to pay CCIP fees.
//
// Returns the message ID, transaction digest, and sequence number.
func SendMessage(
	ctx context.Context,
	ptbClient *client.PTBClient,
	signer bindutils.SuiSigner,
	ccipPkgID string,
	onRampPkgID string,
	ccipObjectRefID string,
	onRampStateID string,
	feeTokenType string,
	feeTokenMetadataID string,
	feeTokenCoinID string,
	destChainSelector uint64,
	receiver []byte,
	data []byte,
	gasBudget uint64,
	evmCallbackGasLimit uint64,
) (messageID string, txDigest string, seqNum uint64, err error) {
	// Resolve latest package IDs (handle contract upgrades)
	latestCcipPkg, err := ptbClient.GetLatestPackageId(ctx, ccipPkgID, "ccip")
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to resolve latest CCIP package: %w", err)
	}
	latestOnRampPkg, err := ptbClient.GetLatestPackageId(ctx, onRampPkgID, "ccip_onramp")
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to resolve latest OnRamp package: %w", err)
	}

	slog.Info("Resolved latest package IDs",
		"ccip", latestCcipPkg,
		"onramp", latestOnRampPkg,
	)

	// Build PTB
	ptb := transaction.NewTransaction()

	callOpts := &bind.CallOpts{
		Signer:    signer,
		GasBudget: &gasBudget,
	}

	// Step 1: create_token_transfer_params (onramp_state_helper)
	// Even for message-only transfers, this must be called with an empty 32-byte receiver.
	helperContract, err := bind.NewBoundContract(latestCcipPkg, "ccip", "onramp_state_helper", ptbClient)
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to bind onramp_state_helper: %w", err)
	}

	encodedHelper, err := helperContract.EncodeCallArgsWithGenerics(
		"create_token_transfer_params",
		[]string{},             // typeArgs
		[]string{},             // typeParams
		[]string{"vector<u8>"}, // paramTypes
		[]any{make([]byte, 32)}, // zeroed 32-byte receiver
		nil, // returnTypes
	)
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to encode create_token_transfer_params: %w", err)
	}

	tokenParamsResult, err := helperContract.AppendPTB(ctx, callOpts, ptb, encodedHelper)
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to append create_token_transfer_params: %w", err)
	}

	// Step 2: ccip_send (onramp)
	// Use the generated OnrampContract for proper encoding
	onRampContract, err := onramp.NewOnramp(latestOnRampPkg, ptbClient)
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to create Onramp contract: %w", err)
	}

	// Build extra args for Sui→EVM
	// The gasLimit in GenericExtraArgsV2 is the EVM callback gas on the destination chain,
	// NOT the Sui PTB gas budget.
	extraArgs := MakeBCSEVMExtraArgsV2(big.NewInt(int64(evmCallbackGasLimit)), false)

	// Use the encoder's CcipSendWithArgs to pass the PTB argument from step 1.
	// All object IDs are resolved from the address book — no placeholders.
	encodedSend, err := onRampContract.Encoder().CcipSendWithArgs(
		[]string{feeTokenType},           // fee token type (e.g. "0x2::sui::SUI")
		bind.Object{Id: ccipObjectRefID}, // CCIPObjectRef (shared, from address book)
		bind.Object{Id: onRampStateID},   // OnRampState (shared, from address book)
		bind.Object{Id: "0x6"},           // Clock (always 0x6)
		destChainSelector,
		receiver,
		data,
		tokenParamsResult,                    // TokenTransferParams from step 1
		bind.Object{Id: feeTokenMetadataID},  // CoinMetadata<T> (from address book)
		bind.Object{Id: feeTokenCoinID},      // Coin<T> (real coin from signer's wallet)
		extraArgs,
	)
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to encode ccip_send: %w", err)
	}

	_, err = onRampContract.Bound().AppendPTB(ctx, callOpts, ptb, encodedSend)
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to append ccip_send: %w", err)
	}

	// Execute PTB
	resp, err := bind.ExecutePTB(ctx, callOpts, ptbClient, ptb)
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to execute PTB: %w", err)
	}

	txDigest = resp.Digest

	// Extract CCIPMessageSent event
	for _, event := range resp.Events {
		if event.PackageId == latestOnRampPkg && strings.HasSuffix(event.Type, "CCIPMessageSent") {
			parsed := event.ParsedJson
			if seqStr, ok := parsed["sequence_number"]; ok {
				seqStrVal, _ := seqStr.(string)
				seqNum, _ = strconv.ParseUint(seqStrVal, 10, 64)
			}
			// Extract message ID from the message header
			if msg, ok := parsed["message"]; ok {
				msgMap, _ := msg.(map[string]any)
				if header, ok := msgMap["header"]; ok {
					headerMap, _ := header.(map[string]any)
					if mid, ok := headerMap["message_id"]; ok {
						messageID = fmt.Sprintf("%v", mid)
					}
				}
			}
			break
		}
	}

	if messageID == "" {
		// Fallback: use tx digest as message ID if event extraction fails
		messageID = txDigest
		slog.Warn("Could not extract CCIPMessageSent event, using tx digest as message ID",
			"txDigest", txDigest,
		)
	}

	return messageID, txDigest, seqNum, nil
}
