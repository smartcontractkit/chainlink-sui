package txm

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"strconv"

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
	TransactionID string
	Sender        string
	Metadata      *commontypes.TxMeta
	Timestamp     uint64
	Payload       string // BCS base64 encoded transaction bytes
	Functions     []*SuiFunction
	Signatures    []string
	RequestType   string
	Attempt       int
	State         TransactionState
	Digest        string
	LastUpdatedAt uint64
	TxError       *suierrors.SuiError
}

func (tx *SuiTx) IncrementAttempts() {
	tx.Attempt += 1
}

func TransactionIDGenerator() string {
	return fmt.Sprintf("0x%s", uuid.New().String())
}

// GenerateTransaction creates a new SuiTx transaction.
// RELEVANT NOTES:
//  1. This function currently uses the unsafe MoveCall API to generate the transaction bytes.
//     This is a temporary solution until we migrate to a way to generate the BCS bytes locally.
//  2. Only supports MoveCall transactions. We will extend this to support PTB transactions in the future.
//
// END OF NOTES SECTION -----
// This function constructs a transaction for calling a Move function on the Sui blockchain by:
// 1. Parsing and validating the provided function string (format: "packageId::module::function")
// 2. Encoding parameter values to their Sui representation based on the provided types
// 3. Making a MoveCall request to the Sui client to generate transaction bytes
// 4. Signing the transaction bytes using the provided signer service
// 5. Creating and returning a complete SuiTx object in "Pending" state
//
// Parameters:
//   - ctx: Context for the operation, used for cancellation and timeouts
//   - pubKey: Public key of the account that will sign and submit the transaction
//   - lggr: Logger for recording operation details and errors
//   - keystoreService: Service for signing the generated transaction bytes
//   - suiClient: Client for interacting with the Sui blockchain
//   - transactionID: Unique identifier for the transaction
//   - txMetadata: Transaction metadata including gas configuration
//   - signerAddress: Address of the account that will sign and submit the transaction
//   - function: A SuiFunction struct containing the package ID, module, and function name
//   - typeArgs: Type arguments for generic functions (corresponds to the <T> parameters in Move)
//   - paramTypes: Array of parameter types as strings (must match the function signature)
//   - paramValues: Array of parameter values (must match the types in paramTypes)
//
// Returns:
//   - *SuiTx: A complete transaction object ready for submission
//   - error: An error if any step of transaction generation fails
//
// The returned SuiTx will have:
//   - State: "Pending"
//   - Attempt: 1
//   - RequestType: "WaitForEffectsCert" or "WaitForLocalExecution"
//   - Timestamp: Current UTC timestamp
//   - Payload: The BCS-serialized transaction bytes
//   - Signatures: Array of signatures produced by the signer service
func GenerateTransaction(
	ctx context.Context,
	pubKey []byte,
	lggr logger.Logger,
	keystoreService loop.Keystore,
	suiClient client.SuiPTBClient,
	requestType string,
	transactionID string, txMetadata *commontypes.TxMeta,
	function *SuiFunction,
	typeArgs []string,
	paramTypes []string,
	paramValues []any,
) (*SuiTx, error) {
	packageObjectId := function.PackageId
	moduleName := function.Module
	functionName := function.Name

	if len(paramTypes) != len(paramValues) {
		msg := fmt.Sprintf("unexpected number of parameters, expected %d, got %d", len(paramTypes), len(paramValues))
		lggr.Error(msg)

		return nil, errors.New(msg)
	}

	signerAddress, err := client.GetAddressFromPublicKey(pubKey)
	if err != nil {
		lggr.Errorf("failed to get address from public key: %v", err)
		return nil, err
	}

	// TODO: we will need to replace this by a BSC serialization
	rsp, err := suiClient.MoveCall(ctx, client.MoveCallRequest{
		Signer:          signerAddress,
		PackageObjectId: packageObjectId,
		Module:          moduleName,
		Function:        functionName,
		// We will only need to pass the type arguments if the function is generic
		TypeArguments: []any{},
		Arguments:     paramValues,
		GasBudget:     txMetadata.GasLimit.Uint64(),
	})

	if err != nil {
		msg := fmt.Sprintf("failed to move call: %v", err)
		lggr.Error(msg)

		return nil, errors.New(msg)
	}

	txBytes, err := base64.StdEncoding.DecodeString(rsp.TxBytes)
	if err != nil {
		msg := fmt.Sprintf("failed to decode tx bytes: %v", err)
		lggr.Error(msg)

		return nil, errors.New(msg)
	}

	signatures, err := keystoreService.Sign(ctx, signerAddress, txBytes)
	if err != nil {
		lggr.Errorf("Error signing transaction: %v", err)
		return nil, err
	}

	// Convert []byte signature to []string
	signatureStrings := []string{client.SerializeSuiSignature(signatures, pubKey)}

	return &SuiTx{
		TransactionID: transactionID,
		Sender:        signerAddress,
		Metadata:      &commontypes.TxMeta{},
		Timestamp:     GetCurrentUnixTimestamp(),
		Payload:       rsp.TxBytes,
		Functions:     []*SuiFunction{function},
		Signatures:    signatureStrings,
		RequestType:   requestType,
		Attempt:       0,
		State:         StatePending,
		Digest:        "",
		LastUpdatedAt: GetCurrentUnixTimestamp(),
		TxError:       nil,
	}, nil
}

// GeneratePTBTransaction creates a new SuiTx transaction for a Programmable Transaction Block (PTB).
// This function constructs a PTB transaction by:
// 1. Determining the gas budget from metadata or using a default value.
// 2. Fetching available SUI coins for the specified signer address.
// 3. Selecting the optimal set of coins to cover the gas budget using SelectCoinsForGasBudget.
// 4. Creating the transaction data (sender, PTB, gas coins, budget, price).
// 5. Marshaling the transaction data into BCS bytes.
// 6. Signing the BCS bytes using the provided signer service.
// 7. Creating and returning a SuiTx object in "Pending" state, ready for submission.
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
//   - signerAddress: Address of the account that will sign and submit the transaction.
//   - ptb: The ProgrammableTransaction block containing the commands to be executed.
//   - simulateTx: Boolean flag indicating whether to simulate the transaction (currently unused).
//
// Returns:
//   - *SuiTx: A complete transaction object ready for submission.
//   - error: An error if any step of transaction generation fails (e.g., fetching coins, selecting gas coins, marshaling, signing).
func GeneratePTBTransaction(
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
) (*SuiTx, error) {
	signerAddress, err := client.GetAddressFromPublicKey(pubKey)
	if err != nil {
		lggr.Errorf("failed to get address from public key: %v", err)
		return nil, err
	}

	// Define gasBudget
	var gasBudget uint64
	if txMetadata.GasLimit != nil {
		gasBudget = txMetadata.GasLimit.Uint64()
	} else {
		gasBudget = uint64(defaultGasBudget)
	}

	// Get available coins
	coinData, err := suiClient.GetCoinsByAddress(ctx, signerAddress)
	if err != nil {
		lggr.Errorf("failed to get coins by address: %v", err)
		return nil, err
	}

	// Select coins for gas budget
	gasBudgetCoins, err := SelectCoinsForGasBudget(gasBudget, coinData)
	if err != nil {
		lggr.Errorf("failed to select coins for gas budget: %v", err)
		return nil, err
	}

	lggr.Debugw("Gas budget coins selected",
		"gasBudget", gasBudget,
		"numCoins", len(gasBudgetCoins),
		"gasBudgetCoins", gasBudgetCoins)

	// Create payment coins using block-vision SDK format
	paymentCoins := make([]transaction.SuiObjectRef, 0, len(gasBudgetCoins))
	for _, coin := range gasBudgetCoins {
		coinObjectIdBytes, coinErr := transaction.ConvertSuiAddressStringToBytes(models.SuiAddress(coin.CoinObjectId))
		if coinErr != nil {
			return nil, coinErr
		}
		versionUint, coinErr := strconv.ParseUint(coin.Version, 10, 64)
		if coinErr != nil {
			return nil, fmt.Errorf("failed to parse version: %w", err)
		}
		digestBytes, coinErr := transaction.ConvertObjectDigestStringToBytes(models.ObjectDigest(coin.Digest))
		if coinErr != nil {
			return nil, fmt.Errorf("failed to convert object digest for payment coin: %w", err)
		}
		paymentCoins = append(paymentCoins, transaction.SuiObjectRef{
			ObjectId: *coinObjectIdBytes,
			Version:  versionUint,
			Digest:   *digestBytes,
		})
	}

	ptb.SetGasBudget(gasBudget)
	ptb.SetSender(models.SuiAddress(signerAddress))
	ptb.SetGasOwner(models.SuiAddress(signerAddress))
	ptb.SetGasPayment(paymentCoins)

	// Use the toBCSBase64 to get transaction bytes for signing (similar to MoveCall)
	txBytes, err := toBCSBase64(ctx, ptb, signerAddress, lggr)
	if err != nil {
		lggr.Errorf("failed to get bcs bytes: %v", err)
		return nil, err
	}

	bytesTx, err := base64.StdEncoding.DecodeString(txBytes)
	if err != nil {
		lggr.Errorf("failed to decode tx bytes: %v", err)
		return nil, err
	}

	// Sign using keystore (similar to working examples)
	signature, err := keystoreService.Sign(ctx, signerAddress, bytesTx)
	if err != nil {
		lggr.Errorf("Error signing transaction: %v", err)
		return nil, err
	}

	// Serialize signature (same as working code)
	signatureStrings := []string{client.SerializeSuiSignature(signature, pubKey)}

	// Extract functions from PTB commands
	functions := []*SuiFunction{}
	for _, command := range ptb.Data.V1.Kind.ProgrammableTransaction.Commands {
		packageIDstr := "0x" + hex.EncodeToString(command.MoveCall.Package[:])
		functions = append(functions, &SuiFunction{
			PackageId: packageIDstr,
			Module:    command.MoveCall.Module,
			Name:      command.MoveCall.Function,
		})
	}

	return &SuiTx{
		TransactionID: transactionID,
		Sender:        signerAddress,
		Metadata:      txMetadata,
		Timestamp:     GetCurrentUnixTimestamp(),
		Payload:       txBytes, // Use base64 encoded bytes
		Functions:     functions,
		Signatures:    signatureStrings,
		RequestType:   requestType,
		Attempt:       0,
		State:         StatePending,
		Digest:        "",
		LastUpdatedAt: GetCurrentUnixTimestamp(),
		TxError:       nil,
	}, nil
}

// SelectCoinsForGasBudget selects the optimal set of coins that match the required gas budget.
// It tries to find coins whose total balance is equal to or greater than the gas budget.
// If exact match isn't possible, it returns coins with the smallest excess over the budget.
func SelectCoinsForGasBudget(gasBudget uint64, availableCoins []models.CoinData) ([]models.CoinData, error) {
	if len(availableCoins) == 0 {
		return nil, fmt.Errorf("no coins available for gas budget")
	}

	// parse all balances once
	balances := make([]uint64, len(availableCoins))
	for i, coin := range availableCoins {
		var balance uint64
		_, err := fmt.Sscanf(coin.Balance, "%d", &balance)
		if err != nil {
			return nil, fmt.Errorf("failed to parse coin balance: %w", err)
		}
		balances[i] = balance
	}

	// create index slice and sort by balance (descending)
	indices := make([]int, len(availableCoins))
	for i := range indices {
		indices[i] = i
	}
	sort.Slice(indices, func(i, j int) bool {
		return balances[indices[i]] > balances[indices[j]]
	})

	// check if there's a single coin that covers the gas budget
	for _, idx := range indices {
		if balances[idx] >= gasBudget {
			return []models.CoinData{availableCoins[idx]}, nil
		}
	}

	// if no single coin is sufficient, find the minimal combination
	selected := make([]models.CoinData, 0)
	var totalBalance uint64

	for _, idx := range indices {
		selected = append(selected, availableCoins[idx])
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

// toBCSBase64 converts a transaction to a BCS base64 string.
// This is taken from the block-vision SDK to gain more control over the signing process.
func toBCSBase64(ctx context.Context, tx *transaction.Transaction, signerAddress string, lggr logger.Logger) (string, error) {
	if tx.Data.V1.GasData.Price == nil {
		if tx.SuiClient != nil {
			rsp, err := tx.SuiClient.SuiXGetReferenceGasPrice(ctx)
			if err != nil {
				return "", err
			}
			tx.SetGasPrice(rsp)
		}
	}
	tx.SetGasBudgetIfNotSet(defaultGasBudget)
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

	lggr.Infow("Transaction Data", "Transaction Data", tx.Data)

	bcsEncodedMsg, err := tx.Data.Marshal()
	if err != nil {
		return "", err
	}
	bcsBase64 := mystenbcs.ToBase64(bcsEncodedMsg)

	return bcsBase64, nil
}
