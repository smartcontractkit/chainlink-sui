package txm

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
	"sort"
	"strconv"
	"strings"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/mystenbcs"
	"github.com/block-vision/sui-go-sdk/transaction"

	"github.com/google/uuid"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/client/suierrors"
)

const defaultGasBudget = 200000000

type SuiFunction struct {
	PackageId string
	Module    string
	Name      string
}

// TransactionState represents the lifecycle state of a transaction in the system.
// It's implemented as an integer-based enum for efficient comparisons.
type TransactionState int

const (
	// StatePending represents a transaction that has been created but not yet submitted.
	StatePending TransactionState = iota
	// StateSubmitted represents a transaction that has been submitted to the network.
	StateSubmitted
	// StateFinalized represents a transaction that has been successfully executed and finalized.
	StateFinalized
	// StateRetriable represents a transaction that encountered an issue but can be retried.
	StateRetriable
	// StateFailed represents a transaction that has failed permanently.
	StateFailed
)

type SuiTx struct {
	TransactionID         string
	Sender                string
	PublicKey             []byte
	Metadata              *commontypes.TxMeta
	Timestamp             uint64
	Payload               string // BCS base64 encoded transaction bytes
	Functions             []*SuiFunction
	Signatures            []string
	RequestType           string
	Attempt               int
	State                 TransactionState
	Digest                string
	LastUpdatedAt         uint64
	TxError               *suierrors.SuiError
	BroadcastError        string
	GasBudget             uint64
	Ptb                   *transaction.Transaction
	PaymentCoinsObjectRef []transaction.SuiObjectRef
	CoinManager           GasCoinManager
}

// UpdateBSCPayload regenerates the BCS payload and signatures for the SuiTx.
// This method is typically used after modifying the transaction's PTB or gas budget,
// ensuring that the transaction bytes and signatures are up-to-date before broadcasting.
//
// Parameters:
//   - ctx: Context for cancellation and deadlines.
//   - lggr: Logger for error and debug output.
//   - keystoreService: Service used to sign the transaction bytes.
//   - suiClient: Client for Sui blockchain operations.
//
// Returns:
//   - error: If any step fails (address derivation, transaction preparation, signing, or encoding).
func (tx *SuiTx) UpdateBSCPayload(
	ctx context.Context,
	lggr logger.Logger,
	keystoreService loop.Keystore,
	suiClient client.SuiPTBClient,
) error {
	signerAddress, err := client.GetAddressFromPublicKey(tx.PublicKey)
	if err != nil {
		return fmt.Errorf("failed to get address from public key: %w", err)
	}

	txBytes, paymentCoins, err := preparePTBTransaction(ctx, signerAddress, suiClient, tx.Ptb, tx.GasBudget, lggr, tx.CoinManager)
	if err != nil {
		return fmt.Errorf("failed to prepare PTB transaction: %w", err)
	}

	tx.Payload = txBytes

	// Get the signer ID (in keystore) of the public key
	signerId := fmt.Sprintf("%064x", tx.PublicKey)

	// Sign using keystore
	bytesTx, err := base64.StdEncoding.DecodeString(txBytes)
	if err != nil {
		lggr.Errorf("failed to decode tx bytes: %v", err)
		return fmt.Errorf("failed to decode tx bytes: %w", err)
	}

	signature, err := keystoreService.Sign(ctx, signerId, suiClient.HashTxBytes(bytesTx))
	if err != nil {
		lggr.Errorf("Error signing transaction: %v", err)
		return fmt.Errorf("failed to sign transaction: %w", err)
	}

	// Serialize signatures for new bcs payload
	signatureStrings := []string{client.SerializeSuiSignature(signature, tx.PublicKey)}
	tx.Signatures = signatureStrings
	tx.PaymentCoinsObjectRef = paymentCoins

	return nil
}

func (tx *SuiTx) IncrementAttempts() {
	tx.Attempt += 1
}

func TransactionIDGenerator() string {
	return fmt.Sprintf("0x%s", uuid.New().String())
}

// GeneratePTBTransactionWithGasEstimation creates a new SuiTx transaction for a PTB with accurate gas estimation.
// This function improves upon GeneratePTBTransaction by using the gas manager to estimate gas requirements
// more accurately before finalizing the transaction.
//
// The process follows these steps:
// 1. Build a preliminary transaction with a temporary gas budget to get transaction bytes.
// 2. Use the gas manager to estimate the actual gas requirements.
// 3. Update the gas budget if estimation yields a lower value (keeping the same selected coins).
// 4. Sign the transaction with the final gas budget.
//
// Note: Coins selected for the preliminary (higher) budget will always cover the final (lower or equal) budget,
// so there's no need to re-select coins or re-fetch from the blockchain.
//
// Parameters:
//   - ctx: Context for the operation, used for cancellation and timeouts.
//   - pubKey: Public key of the account that will sign and submit the transaction.
//   - lggr: Logger for recording operation details and errors.
//   - keystoreService: Service for signing the generated transaction bytes.
//   - suiClient: Client for interacting with the Sui blockchain.
//   - requestType: The type of request for transaction execution (e.g., "WaitForEffectsCert").
//   - transactionID: Unique identifier for the transaction.
//   - txMetadata: Transaction metadata including gas configuration (GasLimit).
//   - ptb: The ProgrammableTransaction block containing the commands to be executed.
//   - simulateTx: Boolean flag indicating whether to simulate the transaction (currently unused).
//   - gasManager: Gas manager for estimating gas requirements.
//   - coinManager: Coin manager for managing gas coins.
//
// Returns:
//   - *SuiTx: A complete transaction object ready for submission with accurate gas estimation.
//   - error: An error if any step of transaction generation fails.
func GeneratePTBTransactionWithGasEstimation(
	ctx context.Context,
	pubKey []byte,
	lggr logger.Logger,
	keystoreService loop.Keystore,
	suiClient client.SuiPTBClient,
	requestType string,
	transactionID string,
	txMetadata *commontypes.TxMeta,
	ptb *transaction.Transaction,
	simulateTx bool,
	gasManager GasManager,
	coinManager GasCoinManager,
) (*SuiTx, error) {
	signerAddress, err := client.GetAddressFromPublicKey(pubKey)
	if err != nil {
		lggr.Errorf("failed to get address from public key: %v", err)
		return nil, err
	}

	var finalGasBudget uint64

	// Step 1: Determine initial gas budget for preliminary transaction
	var preliminaryGasBudget uint64
	if txMetadata != nil && txMetadata.GasLimit != nil {
		preliminaryGasBudget = txMetadata.GasLimit.Uint64()
	} else {
		preliminaryGasBudget = uint64(defaultGasBudget)
	}

	// Step 2: Build preliminary transaction to get transaction bytes for gas estimation
	lggr.Debugw("Building preliminary transaction for gas estimation",
		"preliminaryGasBudget", preliminaryGasBudget)

	preliminaryTx, err := buildPreliminaryTransaction(
		ctx, signerAddress, suiClient, ptb, 
		preliminaryGasBudget, lggr, coinManager,
	)
	if err != nil {
		lggr.Errorf("failed to build preliminary transaction: %v", err)
		return nil, err
	}

	// Step 3: Estimate gas using the gas manager
	estimatedGas, err := gasManager.EstimateGasBudget(ctx, preliminaryTx)
	if err != nil {
		lggr.Warnw("Gas estimation failed, falling back to metadata/default", "error", err, "fallbackBudget", preliminaryGasBudget)
		finalGasBudget = preliminaryGasBudget
	} else {
		// If the estimate is bigger than the provided budget, we need to abort
		if estimatedGas > preliminaryGasBudget {
			return nil, fmt.Errorf("estimated gas is greater than preliminary gas budget: %d > %d", estimatedGas, preliminaryGasBudget)
		}
		finalGasBudget = estimatedGas
	}

	lggr.Debugw("Gas estimation complete",
		"finalGasBudget", finalGasBudget,
		"preliminaryGasBudget", preliminaryGasBudget,
		"estimatedGas", estimatedGas,
	)

	// Step 4: If gas budget changed and is lower, update the PTB and regenerate transaction bytes
	// The coins selected for the higher preliminary budget will definitely cover the lower final budget
	if finalGasBudget != preliminaryGasBudget {
		lggr.Debugw("Updating gas budget in PTB",
			"from", preliminaryGasBudget,
			"to", finalGasBudget)

		// Update the PTB's gas budget to the lower estimated value
		ptb.SetGasBudget(finalGasBudget)

		// Regenerate transaction bytes with new gas budget
		txBytes, err := toBCSBase64(ctx, ptb, signerAddress, lggr, finalGasBudget)
		if err != nil {
			return nil, fmt.Errorf("failed to regenerate tx bytes: %w", err)
		}
		preliminaryTx.Payload = txBytes
		preliminaryTx.GasBudget = finalGasBudget
		preliminaryTx.Metadata = &commontypes.TxMeta{GasLimit: big.NewInt(int64(finalGasBudget))}
	}

	// Step 5: Sign the transaction
	bytesTx, err := base64.StdEncoding.DecodeString(preliminaryTx.Payload)
	if err != nil {
		lggr.Errorf("failed to decode tx bytes: %v", err)
		return nil, err
	}

	signerId := fmt.Sprintf("%064x", pubKey)
	signature, err := keystoreService.Sign(ctx, signerId, suiClient.HashTxBytes(bytesTx))
	if err != nil {
		lggr.Errorf("Error signing transaction: %v", err)
		return nil, err
	}

	signatureStrings := []string{client.SerializeSuiSignature(signature, pubKey)}

	// Extract functions from PTB commands (currently disabled for debugging)
	functions := []*SuiFunction{}

	// Step 6: Return the final signed transaction
	return &SuiTx{
		TransactionID:         transactionID,
		PublicKey:             pubKey,
		Sender:                signerAddress,
		Metadata:              &commontypes.TxMeta{GasLimit: big.NewInt(int64(finalGasBudget))},
		Timestamp:             GetCurrentUnixTimestamp(),
		Payload:               preliminaryTx.Payload,
		Functions:             functions,
		Signatures:            signatureStrings,
		RequestType:           requestType,
		Attempt:               0,
		State:                 StatePending,
		Digest:                "",
		LastUpdatedAt:         GetCurrentUnixTimestamp(),
		TxError:               nil,
		GasBudget:             finalGasBudget,
		CoinManager:           coinManager,
		Ptb:                   ptb,
		PaymentCoinsObjectRef: preliminaryTx.PaymentCoinsObjectRef, // Use the coins selected in preliminary tx
	}, nil
}

// SelectCoinsForGasBudget selects the optimal set of coins that match the required gas budget.
// It tries to find coins whose total balance is equal to or greater than the gas budget.
// If exact match isn't possible, it returns coins with the smallest excess over the budget.
func SelectCoinsForGasBudget(gasBudget uint64, availableCoins []models.CoinData) ([]models.CoinData, error) {
	if len(availableCoins) == 0 {
		return nil, fmt.Errorf("no coins available for gas budget")
	}

	// Filter only SUI coins for gas use
	var suiCoins []models.CoinData
	for _, coin := range availableCoins {
		if strings.HasPrefix(coin.CoinType, "0x2::sui::SUI") {
			suiCoins = append(suiCoins, coin)
		}
	}

	if len(suiCoins) == 0 {
		return nil, fmt.Errorf("no SUI coins available for gas budget")
	}

	// parse all balances once
	balances := make([]uint64, len(suiCoins))
	for i, coin := range suiCoins {
		var balance uint64
		_, err := fmt.Sscanf(coin.Balance, "%d", &balance)
		if err != nil {
			return nil, fmt.Errorf("failed to parse coin balance: %w", err)
		}
		balances[i] = balance
	}

	// create index slice and sort by balance (descending)
	indices := make([]int, len(suiCoins))
	for i := range indices {
		indices[i] = i
	}
	sort.Slice(indices, func(i, j int) bool {
		return balances[indices[i]] > balances[indices[j]]
	})

	// check if there's a single coin that covers the gas budget
	for _, idx := range indices {
		if balances[idx] >= gasBudget {
			return []models.CoinData{suiCoins[idx]}, nil
		}
	}

	// if no single coin is sufficient, find the minimal combination
	selected := make([]models.CoinData, 0)
	var totalBalance uint64

	for _, idx := range indices {
		selected = append(selected, suiCoins[idx])
		totalBalance += balances[idx]

		if totalBalance >= gasBudget {
			break
		}
	}

	if totalBalance < gasBudget {
		return nil, fmt.Errorf("insufficient funds for gas budget: required %d, available %d",
			gasBudget, totalBalance)
	}

	return selected, nil
}

// preparePTBTransaction handles the common logic for setting up a PTB transaction.
// This includes fetching coins, selecting gas coins, setting PTB parameters, and converting to BCS bytes.
// It returns the transaction bytes and payment coins for further processing.
func preparePTBTransaction(
	ctx context.Context,
	signerAddress string,
	suiClient client.SuiPTBClient,
	ptb *transaction.Transaction,
	gasBudget uint64,
	lggr logger.Logger,
	coinManager GasCoinManager,
) (txBytes string, paymentCoins []transaction.SuiObjectRef, err error) {
	// Get available coins for gas
	coinData, err := suiClient.GetCoinsByAddress(ctx, signerAddress)
	if err != nil {
		return "", nil, fmt.Errorf("failed to get coins by address: %w", err)
	}

	// Filter coins that are not locked
	filteredCoinData := make([]models.CoinData, 0)
	for _, coin := range coinData {
		coinObjectIdBytes, coinErr := transaction.ConvertSuiAddressStringToBytes(models.SuiAddress(coin.CoinObjectId))
		if coinErr != nil {
			return "", nil, fmt.Errorf("failed to convert coin object ID to bytes: %w", coinErr)
		}

		if coinManager.IsCoinReserved(*coinObjectIdBytes) {
			continue
		}

		filteredCoinData = append(filteredCoinData, coin)
	}

	// Select coins for gas budget
	gasBudgetCoins, err := SelectCoinsForGasBudget(gasBudget, filteredCoinData)
	if err != nil {
		return "", nil, fmt.Errorf("failed to select coins for gas budget: %w", err)
	}

	lggr.Debugw("Gas budget coins selected",
		"gasBudget", gasBudget,
		"numCoins", len(gasBudgetCoins),
		"gasBudgetCoins", gasBudgetCoins)

	// Create payment coins using block-vision SDK format
	paymentCoins = make([]transaction.SuiObjectRef, 0, len(gasBudgetCoins))
	for _, coin := range gasBudgetCoins {
		coinObjectIdBytes, coinErr := transaction.ConvertSuiAddressStringToBytes(models.SuiAddress(coin.CoinObjectId))
		if coinErr != nil {
			return "", nil, coinErr
		}
		versionUint, coinErr := strconv.ParseUint(coin.Version, 10, 64)
		if coinErr != nil {
			return "", nil, fmt.Errorf("failed to parse version: %w", coinErr)
		}
		digestBytes, coinErr := transaction.ConvertObjectDigestStringToBytes(models.ObjectDigest(coin.Digest))
		if coinErr != nil {
			return "", nil, fmt.Errorf("failed to convert object digest for payment coin: %w", coinErr)
		}
		paymentCoins = append(paymentCoins, transaction.SuiObjectRef{
			ObjectId: *coinObjectIdBytes,
			Version:  versionUint,
			Digest:   *digestBytes,
		})
	}

	// Set transaction parameters
	ptb.SetGasBudget(gasBudget)
	ptb.SetSender(models.SuiAddress(signerAddress))
	ptb.SetGasOwner(models.SuiAddress(signerAddress))
	ptb.SetGasPayment(paymentCoins)

	// Get transaction bytes
	txBytes, err = toBCSBase64(ctx, ptb, signerAddress, lggr, gasBudget)
	if err != nil {
		return "", nil, fmt.Errorf("failed to get bcs bytes: %w", err)
	}

	return txBytes, paymentCoins, nil
}

// buildPreliminaryTransaction creates a minimal SuiTx for gas estimation purposes.
// This transaction is not signed and is only used to get the transaction bytes for estimation.
func buildPreliminaryTransaction(
	ctx context.Context,
	signerAddress string,
	suiClient client.SuiPTBClient,
	ptb *transaction.Transaction,
	gasBudget uint64,
	lggr logger.Logger,
	coinManager GasCoinManager,
) (*SuiTx, error) {
	// Use common preparation logic
	txBytes, paymentCoins, err := preparePTBTransaction(ctx, signerAddress, suiClient, ptb, gasBudget, lggr, coinManager)
	if err != nil {
		return nil, err
	}

	// Create a minimal SuiTx for gas estimation (no signatures needed)
	return &SuiTx{
		Payload:               txBytes,
		Metadata:              &commontypes.TxMeta{GasLimit: big.NewInt(int64(gasBudget))},
		Sender:                signerAddress,
		PaymentCoinsObjectRef: paymentCoins,
	}, nil
}

// toBCSBase64 converts a transaction to a BCS base64 string.
// This is taken from the block-vision SDK to gain more control over the signing process.
func toBCSBase64(
	ctx context.Context,
	tx *transaction.Transaction,
	signerAddress string,
	lggr logger.Logger,
	gasBudget uint64,
) (string, error) {
	if tx.Data.V1.GasData.Price == nil {
		if tx.SuiClient != nil {
			rsp, err := tx.SuiClient.SuiXGetReferenceGasPrice(ctx)
			if err != nil {
				return "", err
			}
			tx.SetGasPrice(rsp)
		}
	}
	tx.SetGasBudgetIfNotSet(gasBudget)
	tx.SetSenderIfNotSet(models.SuiAddress(signerAddress))

	if tx.Data.V1.Sender == nil {
		return "", errors.New("sender not set")
	}
	if tx.Data.V1.GasData.Owner == nil {
		tx.SetGasOwner(models.SuiAddress(signerAddress))
	}
	if !tx.Data.V1.GasData.IsAllSet() {
		return "", errors.New("gas data not all set")
	}

	lggr.Debugw("Transaction Data", "Transaction Data", tx.Data)

	bcsEncodedMsg, err := tx.Data.Marshal()
	if err != nil {
		return "", err
	}
	bcsBase64 := mystenbcs.ToBase64(bcsEncodedMsg)

	return bcsBase64, nil
}
