package sui

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/transaction"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	bindutils "github.com/smartcontractkit/chainlink-sui/bindings/utils"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
)

const suiCoinObjectType = "0x2::coin::Coin<0x2::sui::SUI>"

func normalizeCoinObjectType(typeTag string) string {
	trimmed := strings.TrimSpace(typeTag)
	if trimmed == "" {
		return trimmed
	}
	if strings.HasPrefix(trimmed, "0x2::coin::Coin<") {
		return trimmed
	}
	return "0x2::coin::Coin<" + trimmed + ">"
}

// SuiCoinPool provides pre-split SUI coin object IDs.
// Each message consumes two coin objects: one for gas, one for CCIP fee token.
type SuiCoinPool struct {
	coins chan string
}

func newSuiCoinPool(ids []string) *SuiCoinPool {
	ch := make(chan string, len(ids))
	for _, id := range ids {
		ch <- id
	}
	return &SuiCoinPool{coins: ch}
}

func (p *SuiCoinPool) Pop(ctx context.Context) (string, error) {
	select {
	case id := <-p.coins:
		return id, nil
	case <-ctx.Done():
		return "", fmt.Errorf("context canceled while waiting for SUI coin: %w", ctx.Err())
	}
}

func (p *SuiCoinPool) Size() int {
	return len(p.coins)
}

// TokenCoinPool provides pre-split token coin object IDs for Sui→EVM transfers.
// Each message consumes exactly one token coin object via lock_or_burn.
type TokenCoinPool struct {
	coins         chan string
	coinType      string
	amountPerCoin uint64
}

func newTokenCoinPool(ids []string, coinType string, amountPerCoin uint64) *TokenCoinPool {
	ch := make(chan string, len(ids))
	for _, id := range ids {
		ch <- id
	}
	return &TokenCoinPool{
		coins:         ch,
		coinType:      coinType,
		amountPerCoin: amountPerCoin,
	}
}

func (p *TokenCoinPool) Pop(ctx context.Context) (string, error) {
	select {
	case id := <-p.coins:
		return id, nil
	case <-ctx.Done():
		return "", fmt.Errorf("context canceled while waiting for token coin: %w", ctx.Err())
	}
}

func (p *TokenCoinPool) Size() int {
	return len(p.coins)
}

func (p *TokenCoinPool) CoinType() string {
	return p.coinType
}

func (p *TokenCoinPool) AmountPerCoin() uint64 {
	return p.amountPerCoin
}

// CalculateRequiredSuiCoins returns how many SUI objects to pre-split for a run.
// We need 2 per message (gas + fee), plus a 20% buffer for retries/headroom.
func CalculateRequiredSuiCoins(messageCount int) int {
	if messageCount < 1 {
		return 0
	}
	base := messageCount * 2
	withBuffer := (base*12 + 9) / 10 // round up base*1.2
	if withBuffer < 4 {
		withBuffer = 4
	}
	return withBuffer
}

// PrepareSuiCoinPool pre-splits SUI into multiple coin objects that can be consumed by load messages.
func PrepareSuiCoinPool(
	ctx context.Context,
	ptbClient *client.PTBClient,
	signer bindutils.SuiSigner,
	signerAddress string,
	messageCount int,
	amountPerCoin uint64,
) (*SuiCoinPool, error) {
	requiredCoins := CalculateRequiredSuiCoins(messageCount)
	if requiredCoins == 0 {
		return nil, fmt.Errorf("invalid message count: %d", messageCount)
	}

	sourceCoinID, sourceBalance, err := fetchLargestSuiCoin(ctx, ptbClient, signerAddress)
	if err != nil {
		return nil, err
	}

	usableCoinIDs, err := fetchUsableSuiCoinIDs(ctx, ptbClient, signerAddress, sourceCoinID, amountPerCoin)
	if err != nil {
		return nil, err
	}

	if len(usableCoinIDs) >= requiredCoins {
		slog.Info("Using existing SUI coins for load test pool",
			"requiredCoins", requiredCoins,
			"availableCoins", len(usableCoinIDs),
			"minBalancePerCoin", amountPerCoin,
		)
		return newSuiCoinPool(usableCoinIDs[:requiredCoins]), nil
	}

	missingCoins := requiredCoins - len(usableCoinIDs)
	totalSplitNeeded := uint64(missingCoins) * amountPerCoin //nolint:gosec // bounded by run config input
	if sourceBalance <= totalSplitNeeded {
		return nil, fmt.Errorf(
			"insufficient SUI balance for split: source coin %s has %d, need more than %d to create %d additional coins",
			sourceCoinID,
			sourceBalance,
			totalSplitNeeded,
			missingCoins,
		)
	}

	slog.Info("Preparing pre-split SUI coin pool",
		"requiredCoins", requiredCoins,
		"existingUsableCoins", len(usableCoinIDs),
		"missingCoins", missingCoins,
		"amountPerCoin", amountPerCoin,
		"sourceCoinID", sourceCoinID,
		"sourceBalance", sourceBalance,
	)

	err = splitSuiCoinObjects(ctx, ptbClient, signer, signerAddress, sourceCoinID, missingCoins, amountPerCoin)
	if err != nil {
		return nil, err
	}

	usableCoinIDs, err = fetchUsableSuiCoinIDs(ctx, ptbClient, signerAddress, sourceCoinID, amountPerCoin)
	if err != nil {
		return nil, err
	}
	if len(usableCoinIDs) < requiredCoins {
		return nil, fmt.Errorf(
			"insufficient usable SUI coins after split: have %d need %d (minBalancePerCoin=%d)",
			len(usableCoinIDs),
			requiredCoins,
			amountPerCoin,
		)
	}

	slog.Info("Prepared pre-split SUI coin pool", "poolCoins", requiredCoins)
	return newSuiCoinPool(usableCoinIDs[:requiredCoins]), nil
}

func splitSuiCoinObjects(
	ctx context.Context,
	ptbClient *client.PTBClient,
	signer bindutils.SuiSigner,
	signerAddress string,
	sourceCoinID string,
	numCoins int,
	amountPerCoin uint64,
) error {
	if numCoins < 1 {
		return nil
	}

	batchSize := 20
	splitGasBudget := uint64(100_000_000)

	for batchStart := 0; batchStart < numCoins; batchStart += batchSize {
		remaining := numCoins - batchStart
		if remaining > batchSize {
			remaining = batchSize
		}

		ptb := transaction.NewTransaction()
		gasArg := ptb.Gas()
		splitArgs := make([]transaction.Argument, 0, remaining)

		for i := 0; i < remaining; i++ {
			splitCoin := ptb.SplitCoins(gasArg, []transaction.Argument{ptb.Pure(amountPerCoin)})
			splitArgs = append(splitArgs, splitCoin)
		}

		ptb.TransferObjects(splitArgs, ptb.Pure(signerAddress))

		callOpts := &bind.CallOpts{
			Signer:           signer,
			GasObject:        sourceCoinID,
			GasBudget:        &splitGasBudget,
			WaitForExecution: true,
		}

		resp, err := bind.ExecutePTB(ctx, callOpts, ptbClient, ptb)
		if err != nil {
			return fmt.Errorf("failed to execute split PTB at batch %d: %w", batchStart, err)
		}

		slog.Info("Completed SUI split batch",
			"batchStart", batchStart,
			"batchSize", remaining,
			"txDigest", resp.Digest,
		)
	}

	return nil
}

// PrepareTokenCoinPool pre-splits the sender's token coins into exact per-message amounts.
//
// coinType is the full Move coin type (e.g. "0x...::ccip_burn_mint_token::CCIP_BURN_MINT_TOKEN").
// amountPerCoin is the exact transfer amount per message in base units.
// requiredCoins equals messageCount (one coin consumed per message).
func PrepareTokenCoinPool(
	ctx context.Context,
	ptbClient *client.PTBClient,
	signer bindutils.SuiSigner,
	signerAddress string,
	coinType string,
	messageCount int,
	amountPerCoin uint64,
) (*TokenCoinPool, error) {
	if messageCount < 1 {
		return nil, fmt.Errorf("invalid message count: %d", messageCount)
	}
	if amountPerCoin == 0 {
		return nil, fmt.Errorf("amount per coin must be > 0")
	}

	requiredCoins := messageCount
	totalNeeded := uint64(requiredCoins) * amountPerCoin //nolint:gosec // bounded by run config input

	coinObjectType := normalizeCoinObjectType(coinType)

	sourceCoinID, sourceBalance, err := fetchLargestCoinOfType(ctx, ptbClient, signerAddress, coinObjectType)
	if err != nil {
		return nil, err
	}
	if sourceBalance < totalNeeded {
		return nil, fmt.Errorf(
			"insufficient token balance for split: coin type %s, balance %d, need %d for %d messages",
			coinType,
			sourceBalance,
			totalNeeded,
			requiredCoins,
		)
	}

	exactCoinIDs, err := fetchExactCoinIDsOfType(ctx, ptbClient, signerAddress, coinObjectType, amountPerCoin)
	if err != nil {
		return nil, err
	}

	if len(exactCoinIDs) >= requiredCoins {
		slog.Info("Using existing token coins for load test pool",
			"coinType", coinType,
			"coinObjectType", coinObjectType,
			"requiredCoins", requiredCoins,
			"availableCoins", len(exactCoinIDs),
			"amountPerCoin", amountPerCoin,
		)
		return newTokenCoinPool(exactCoinIDs[:requiredCoins], coinType, amountPerCoin), nil
	}

	missingCoins := requiredCoins - len(exactCoinIDs)
	slog.Info("Preparing pre-split token coin pool",
		"coinType", coinType,
		"coinObjectType", coinObjectType,
		"requiredCoins", requiredCoins,
		"existingUsableCoins", len(exactCoinIDs),
		"missingCoins", missingCoins,
		"amountPerCoin", amountPerCoin,
		"sourceCoinID", sourceCoinID,
		"sourceBalance", sourceBalance,
	)

	err = splitTokenCoinObjects(ctx, ptbClient, signer, signerAddress, sourceCoinID, coinType, missingCoins, amountPerCoin)
	if err != nil {
		return nil, err
	}

	exactCoinIDs, err = fetchExactCoinIDsOfType(ctx, ptbClient, signerAddress, coinObjectType, amountPerCoin)
	if err != nil {
		return nil, err
	}
	if len(exactCoinIDs) < requiredCoins {
		return nil, fmt.Errorf(
			"insufficient usable token coins after split: have %d need %d (coinType=%s amountPerCoin=%d)",
			len(exactCoinIDs),
			requiredCoins,
			coinType,
			amountPerCoin,
		)
	}

	slog.Info("Prepared pre-split token coin pool", "coinType", coinType, "coinObjectType", coinObjectType, "poolCoins", requiredCoins)
	return newTokenCoinPool(exactCoinIDs[:requiredCoins], coinType, amountPerCoin), nil
}

func splitTokenCoinObjects(
	ctx context.Context,
	ptbClient *client.PTBClient,
	signer bindutils.SuiSigner,
	signerAddress string,
	sourceCoinID string,
	coinType string,
	numCoins int,
	amountPerCoin uint64,
) error {
	if numCoins < 1 {
		return nil
	}

	batchSize := 20
	splitGasBudget := uint64(100_000_000)

	for batchStart := 0; batchStart < numCoins; batchStart += batchSize {
		remaining := numCoins - batchStart
		if remaining > batchSize {
			remaining = batchSize
		}

		ptb := transaction.NewTransaction()
		objIdBytes, convErr := transaction.ConvertSuiAddressStringToBytes(models.SuiAddress(sourceCoinID))
		if convErr != nil {
			return fmt.Errorf("failed to convert source coin ID to bytes: %w", convErr)
		}
		sourceArg := ptb.Object(transaction.CallArg{Object: &transaction.ObjectArg{ImmOrOwnedObject: &transaction.SuiObjectRef{ObjectId: *objIdBytes}}})
		splitArgs := make([]transaction.Argument, 0, remaining)

		for i := 0; i < remaining; i++ {
			splitCoin := ptb.SplitCoins(sourceArg, []transaction.Argument{ptb.Pure(amountPerCoin)})
			splitArgs = append(splitArgs, splitCoin)
		}

		ptb.TransferObjects(splitArgs, ptb.Pure(signerAddress))

		callOpts := &bind.CallOpts{
			Signer:           signer,
			GasObject:        sourceCoinID,
			GasBudget:        &splitGasBudget,
			WaitForExecution: true,
		}

		resp, err := bind.ExecutePTB(ctx, callOpts, ptbClient, ptb)
		if err != nil {
			return fmt.Errorf("failed to execute token split PTB at batch %d: %w", batchStart, err)
		}

		slog.Info("Completed token split batch",
			"coinType", coinType,
			"batchStart", batchStart,
			"batchSize", remaining,
			"txDigest", resp.Digest,
		)
	}

	return nil
}

func fetchUsableCoinIDsOfType(
	ctx context.Context,
	ptbClient *client.PTBClient,
	signerAddress string,
	coinType string,
	excludeCoinID string,
	minBalance uint64,
) ([]string, error) {
	coins, err := ptbClient.QueryCoinsByAddress(ctx, signerAddress, coinType)
	if err != nil {
		return nil, fmt.Errorf("failed to query coins of type %s for address %s: %w", coinType, signerAddress, err)
	}

	ids := make([]string, 0, len(coins))
	for _, coin := range coins {
		if coin.GetObjectId() == excludeCoinID {
			continue
		}
		if coin.GetBalance() < minBalance {
			continue
		}
		ids = append(ids, coin.GetObjectId())
	}

	return ids, nil
}

func fetchExactCoinIDsOfType(
	ctx context.Context,
	ptbClient *client.PTBClient,
	signerAddress string,
	coinType string,
	minBalance uint64,
) ([]string, error) {
	coins, err := ptbClient.QueryCoinsByAddress(ctx, signerAddress, coinType)
	if err != nil {
		return nil, fmt.Errorf("failed to query coins of type %s for exact balance lookup: %w", coinType, err)
	}

	ids := make([]string, 0, len(coins))
	for _, coin := range coins {
		if coin.GetBalance() != minBalance {
			continue
		}
		ids = append(ids, coin.GetObjectId())
	}

	return ids, nil
}

func fetchLargestCoinOfType(ctx context.Context, ptbClient *client.PTBClient, signerAddress string, coinType string) (coinID string, balance uint64, err error) {
	coins, err := ptbClient.QueryCoinsByAddress(ctx, signerAddress, coinType)
	if err != nil {
		return "", 0, fmt.Errorf("failed to query coins of type %s: %w", coinType, err)
	}
	if len(coins) == 0 {
		return "", 0, fmt.Errorf("no coins of type %s found for address %s", coinType, signerAddress)
	}

	for _, c := range coins {
		b := c.GetBalance()
		if b >= balance {
			balance = b
			coinID = c.GetObjectId()
		}
	}

	if coinID == "" {
		return "", 0, fmt.Errorf("could not determine a source coin of type %s", coinType)
	}

	return coinID, balance, nil
}

func fetchUsableSuiCoinIDs(
	ctx context.Context,
	ptbClient *client.PTBClient,
	signerAddress string,
	excludeCoinID string,
	minBalance uint64,
) ([]string, error) {
	coins, err := ptbClient.QueryCoinsByAddress(ctx, signerAddress, suiCoinObjectType)
	if err != nil {
		return nil, fmt.Errorf("failed to query SUI coins for address %s: %w", signerAddress, err)
	}

	ids := make([]string, 0, len(coins))
	for _, coin := range coins {
		if coin.GetObjectId() == excludeCoinID {
			continue
		}
		if coin.GetBalance() < minBalance {
			continue
		}
		ids = append(ids, coin.GetObjectId())
	}

	return ids, nil
}

func fetchLargestSuiCoin(ctx context.Context, ptbClient *client.PTBClient, signerAddress string) (coinID string, balance uint64, err error) {
	coins, err := ptbClient.QueryCoinsByAddress(ctx, signerAddress, suiCoinObjectType)
	if err != nil {
		return "", 0, fmt.Errorf("failed to query SUI coins: %w", err)
	}
	if len(coins) == 0 {
		return "", 0, fmt.Errorf("no SUI coins found for address %s", signerAddress)
	}

	for _, c := range coins {
		b := c.GetBalance()
		if b >= balance {
			balance = b
			coinID = c.GetObjectId()
		}
	}

	if coinID == "" {
		return "", 0, fmt.Errorf("could not determine a source SUI coin")
	}

	return coinID, balance, nil
}
