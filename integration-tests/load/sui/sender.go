package sui

import (
	"context"
	"encoding/base64"
	"encoding/hex"
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

const minSplitAmountPerCoin uint64 = 150_000_000 // 0.15 SUI

// RecommendedSplitAmountPerCoin returns a conservative split amount for SUI fee coins
// based on estimated onramp fee, plus a small buffer for fee fluctuations.
// Tuned from load-test observations to avoid significant over-funding.
func RecommendedSplitAmountPerCoin(estimatedFee uint64) uint64 {
	if estimatedFee == 0 {
		return minSplitAmountPerCoin
	}
	withBuffer := estimatedFee + 20_000_000 // fee + 0.02 SUI buffer
	if withBuffer < minSplitAmountPerCoin {
		return minSplitAmountPerCoin
	}
	return withBuffer
}

// EstimateSuiToEVMFee returns the current get_fee quote for a message-only Sui->EVM send.
func EstimateSuiToEVMFee(
	ctx context.Context,
	ptbClient *client.PTBClient,
	signer bindutils.SuiSigner,
	ccipPkgID string,
	onRampPkgID string,
	ccipObjectRefID string,
	onRampStateID string,
	feeTokenType string,
	feeTokenMetadataID string,
	destChainSelector uint64,
	receiver []byte,
	data []byte,
	evmCallbackGasLimit uint64,
) (uint64, error) {
	latestCcipPkg, latestOnRampPkg := resolveLatestPackages(ctx, ptbClient, ccipPkgID, onRampPkgID, ccipObjectRefID, onRampStateID)

	onRampContract, err := onramp.NewOnramp(latestOnRampPkg, ptbClient)
	if err != nil {
		return 0, fmt.Errorf("failed to create Onramp contract for fee estimate: %w", err)
	}

	extraArgs := MakeBCSEVMExtraArgsV2(big.NewInt(int64(evmCallbackGasLimit)), true)
	devOpts := &bind.CallOpts{Signer: signer}

	fee, err := onRampContract.DevInspect().GetFee(
		ctx,
		devOpts,
		[]string{feeTokenType},
		bind.Object{Id: ccipObjectRefID},
		bind.Object{Id: "0x6"},
		destChainSelector,
		receiver,
		data,
		[]string{},
		[]uint64{},
		bind.Object{Id: feeTokenMetadataID},
		extraArgs,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to estimate get_fee (ccip=%s onramp=%s): %w", latestCcipPkg, latestOnRampPkg, err)
	}

	return fee, nil
}

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
	gasCoinID string,
	feeTokenType string,
	feeTokenMetadataID string,
	feeTokenCoinID string,
	feeAmount uint64,
	destChainSelector uint64,
	receiver []byte,
	data []byte,
	gasBudget uint64,
	evmCallbackGasLimit uint64,
) (messageID string, txDigest string, seqNum uint64, err error) {
	return SendTokenMessage(
		ctx,
		ptbClient,
		signer,
		ccipPkgID,
		onRampPkgID,
		ccipObjectRefID,
		onRampStateID,
		gasCoinID,
		feeTokenType,
		feeTokenMetadataID,
		feeTokenCoinID,
		feeAmount,
		"", // no token coin for message-only
		"", // no token address for message-only
		"", // no token coin type for message-only
		"", // no token pool package for message-only
		"", // no pool state for message-only
		"", // no deny list for message-only
		"", // no token state for message-only
		destChainSelector,
		receiver,
		data,
		gasBudget,
		evmCallbackGasLimit,
	)
}

// SendTokenMessage builds and executes a CCIP token transfer from Sui to EVM.
//
// For message-only mode, pass empty tokenCoinID, tokenCoinType, and token pool IDs.
// For token transfers, the PTB sequence is:
//
//	create_token_transfer_params -> managed_token_pool.lock_or_burn -> ccip_send
//
// Returns the message ID, transaction digest, and sequence number.
func SendTokenMessage(
	ctx context.Context,
	ptbClient *client.PTBClient,
	signer bindutils.SuiSigner,
	ccipPkgID string,
	onRampPkgID string,
	ccipObjectRefID string,
	onRampStateID string,
	gasCoinID string,
	feeTokenType string,
	feeTokenMetadataID string,
	feeTokenCoinID string,
	feeAmount uint64,
	tokenCoinID string,
	_ string,
	tokenCoinType string,
	tokenPoolPkgID string,
	tokenPoolStateID string,
	denyListObjectID string,
	tokenStateObjectID string,
	destChainSelector uint64,
	receiver []byte,
	data []byte,
	gasBudget uint64,
	evmCallbackGasLimit uint64,
) (messageID string, txDigest string, seqNum uint64, err error) {
	latestCcipPkg, latestOnRampPkg := resolveLatestPackages(ctx, ptbClient, ccipPkgID, onRampPkgID, ccipObjectRefID, onRampStateID)

	slog.Info("Resolved latest package IDs",
		"ccip", latestCcipPkg,
		"onramp", latestOnRampPkg,
	)

	signerAddress, err := signer.GetAddress()
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to get signer address: %w", err)
	}
	if feeAmount == 0 {
		return "", "", 0, fmt.Errorf("feeAmount must be > 0")
	}

	// Build PTB
	ptb := transaction.NewTransaction()

	callOpts := &bind.CallOpts{
		Signer:    signer,
		GasObject: gasCoinID,
		GasBudget: &gasBudget,
	}

	slog.Debug("Using explicit SUI gas/fee coins",
		"gasCoinId", gasCoinID,
		"feeCoinId", feeTokenCoinID,
	)

	// Step 1: create_token_transfer_params (onramp_state_helper)
	// For EVM destinations, use the same normalized receiver bytes as ccip_send.
	helperContract, err := bind.NewBoundContract(latestCcipPkg, "ccip", "onramp_state_helper", ptbClient)
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to bind onramp_state_helper: %w", err)
	}

	encodedHelper, err := helperContract.EncodeCallArgsWithGenerics(
		"create_token_transfer_params",
		[]string{},             // typeArgs
		[]string{},             // typeParams
		[]string{"vector<u8>"}, // paramTypes
		[]any{receiver},        // normalized receiver bytes
		nil,                    // returnTypes
	)
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to encode create_token_transfer_params: %w", err)
	}

	tokenParamsResult, err := helperContract.AppendPTB(ctx, callOpts, ptb, encodedHelper)
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to append create_token_transfer_params: %w", err)
	}

	// Step 2 (optional): managed_token_pool.lock_or_burn
	if tokenCoinID != "" && tokenPoolPkgID != "" {
		tokenCoinBalance, err := fetchCoinBalanceByIDAndType(ctx, ptbClient, signerAddress, tokenCoinID, tokenCoinType)
		if err != nil {
			return "", "", 0, fmt.Errorf("failed to read token coin balance: %w", err)
		}
		if tokenCoinBalance == 0 {
			return "", "", 0, fmt.Errorf("token coin %s has zero balance", tokenCoinID)
		}

		err = AppendManagedTokenPoolLockOrBurn(
			ctx,
			ptbClient,
			signer,
			ptb,
			callOpts,
			tokenPoolPkgID,
			tokenCoinType,
			tokenCoinID,
			*tokenParamsResult,
			destChainSelector,
			ccipObjectRefID,
			tokenPoolStateID,
			denyListObjectID,
			tokenStateObjectID,
		)
		if err != nil {
			return "", "", 0, fmt.Errorf("failed to append lock_or_burn: %w", err)
		}
	}

	// Step 3: ccip_send (onramp)
	// Use the generated OnrampContract for proper encoding
	onRampContract, err := onramp.NewOnramp(latestOnRampPkg, ptbClient)
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to create Onramp contract: %w", err)
	}

	// Build extra args for Sui→EVM
	// The gasLimit in GenericExtraArgsV2 is the EVM callback gas on the destination chain,
	// NOT the Sui PTB gas budget.
	extraArgs := MakeBCSEVMExtraArgsV2(big.NewInt(int64(evmCallbackGasLimit)), true)

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
		tokenParamsResult,                   // TokenTransferParams from step 1
		bind.Object{Id: feeTokenMetadataID}, // CoinMetadata<T> (from address book)
		bind.Object{Id: feeTokenCoinID},
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
						normalized, normalizeErr := normalizeMessageIDToHex(mid)
						if normalizeErr != nil {
							slog.Warn("Failed to normalize message_id to hex",
								"raw", fmt.Sprintf("%v", mid),
								"error", normalizeErr,
							)
							messageID = fmt.Sprintf("%v", mid)
						} else {
							messageID = normalized
						}
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

func resolveLatestPackageIDFromStateObject(ctx context.Context, ptbClient *client.PTBClient, stateObjectID string) (string, error) {
	obj, err := ptbClient.ReadObjectId(ctx, stateObjectID)
	if err != nil {
		return "", fmt.Errorf("read object %s: %w", stateObjectID, err)
	}

	if obj.GetJson() == nil || obj.GetJson().GetStructValue() == nil {
		return "", fmt.Errorf("object %s has no struct json", stateObjectID)
	}

	fields := obj.GetJson().GetStructValue().GetFields()
	packageIDsField, exists := fields["package_ids"]
	if !exists || packageIDsField == nil || packageIDsField.GetListValue() == nil {
		return "", fmt.Errorf("object %s does not expose package_ids", stateObjectID)
	}

	values := packageIDsField.GetListValue().Values
	if len(values) == 0 {
		return "", fmt.Errorf("object %s has empty package_ids", stateObjectID)
	}

	latest := values[len(values)-1].GetStringValue()
	if latest == "" {
		return "", fmt.Errorf("object %s returned empty latest package id", stateObjectID)
	}

	return latest, nil
}

func resolveLatestPackages(
	ctx context.Context,
	ptbClient *client.PTBClient,
	ccipPkgID string,
	onRampPkgID string,
	ccipObjectRefID string,
	onRampStateID string,
) (string, string) {
	latestCcipPkg := ccipPkgID
	if pkg, err := resolveLatestPackageIDFromStateObject(ctx, ptbClient, ccipObjectRefID); err == nil && pkg != "" {
		latestCcipPkg = pkg
	} else if err != nil {
		slog.Warn("Failed to resolve CCIP package from state object, falling back",
			"stateObjectId", ccipObjectRefID,
			"fallback", ccipPkgID,
			"error", err,
		)
	}

	latestOnRampPkg := onRampPkgID
	if pkg, err := resolveLatestPackageIDFromStateObject(ctx, ptbClient, onRampStateID); err == nil && pkg != "" {
		latestOnRampPkg = pkg
	} else if err != nil {
		slog.Warn("Failed to resolve OnRamp package from state object, falling back",
			"stateObjectId", onRampStateID,
			"fallback", onRampPkgID,
			"error", err,
		)
	}

	if fallbackCCIP, err := ptbClient.GetLatestPackageId(ctx, latestCcipPkg, "ccip"); err == nil && fallbackCCIP != "" {
		latestCcipPkg = fallbackCCIP
	}
	if fallbackOnRamp, err := ptbClient.GetLatestPackageId(ctx, latestOnRampPkg, "ccip_onramp"); err == nil && fallbackOnRamp != "" {
		latestOnRampPkg = fallbackOnRamp
	}

	return latestCcipPkg, latestOnRampPkg
}

func fetchCoinBalanceByID(ctx context.Context, ptbClient *client.PTBClient, signerAddress string, coinID string) (uint64, error) {
	coins, err := ptbClient.QueryCoinsByAddress(ctx, signerAddress, suiCoinObjectType)
	if err != nil {
		return 0, fmt.Errorf("failed to query SUI coins for balance lookup: %w", err)
	}

	for _, coin := range coins {
		if coin.GetObjectId() == coinID {
			return coin.GetBalance(), nil
		}
	}

	return 0, fmt.Errorf("coin %s not found in signer wallet", coinID)
}

// fetchCoinBalanceByIDAndType returns the balance of a specific coin object by ID and type.
func fetchCoinBalanceByIDAndType(ctx context.Context, ptbClient *client.PTBClient, signerAddress string, coinID string, coinType string) (uint64, error) {
	coins, err := ptbClient.QueryCoinsByAddress(ctx, signerAddress, normalizeCoinObjectType(coinType))
	if err != nil {
		return 0, fmt.Errorf("failed to query coins of type %s for balance lookup: %w", coinType, err)
	}

	for _, coin := range coins {
		if coin.GetObjectId() == coinID {
			return coin.GetBalance(), nil
		}
	}

	return 0, fmt.Errorf("coin %s of type %s not found in signer wallet", coinID, coinType)
}

func normalizeMessageIDToHex(messageID any) (string, error) {
	switch v := messageID.(type) {
	case string:
		s := strings.TrimSpace(v)
		if s == "" {
			return "", fmt.Errorf("empty message_id string")
		}

		if strings.HasPrefix(s, "0x") || strings.HasPrefix(s, "0X") {
			raw := strings.TrimPrefix(strings.TrimPrefix(s, "0x"), "0X")
			if _, err := hex.DecodeString(raw); err != nil {
				return "", fmt.Errorf("invalid hex message_id: %w", err)
			}
			return "0x" + strings.ToLower(raw), nil
		}

		if _, err := hex.DecodeString(s); err == nil {
			return "0x" + strings.ToLower(s), nil
		}

		decoded, err := base64.StdEncoding.DecodeString(s)
		if err != nil {
			return "", fmt.Errorf("message_id is neither hex nor base64: %w", err)
		}
		return "0x" + hex.EncodeToString(decoded), nil

	case []byte:
		return "0x" + hex.EncodeToString(v), nil

	case []any:
		bytes := make([]byte, 0, len(v))
		for i, item := range v {
			n, ok := item.(float64)
			if !ok {
				return "", fmt.Errorf("unsupported array element type at index %d: %T", i, item)
			}
			if n < 0 || n > 255 {
				return "", fmt.Errorf("byte out of range at index %d: %v", i, n)
			}
			bytes = append(bytes, byte(n))
		}
		return "0x" + hex.EncodeToString(bytes), nil

	default:
		return "", fmt.Errorf("unsupported message_id type: %T", messageID)
	}
}
