package txm

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"sort"

	"github.com/google/uuid"
	"github.com/pattonkan/sui-go/sui"
	"github.com/pattonkan/sui-go/sui/suiptb"
	"github.com/pattonkan/sui-go/suiclient"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/fardream/go-bcs/bcs"

	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/client/suierrors"
	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
	"github.com/smartcontractkit/chainlink-sui/relayer/signer"
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
	Payload       []byte
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
//   - lggr: Logger for recording operation details and errors
//   - signerService: Service for signing the generated transaction bytes
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
	lggr logger.Logger,
	signerService signer.SuiSigner,
	suiClient client.SuiPTBClient,
	requestType string,
	transactionID string, txMetadata *commontypes.TxMeta,
	signerAddress string, function *SuiFunction,
	typeArgs []string, paramTypes []string, paramValues []any,
) (*SuiTx, error) {
	packageObjectId := function.PackageId
	moduleName := function.Module
	functionName := function.Name

	if len(paramTypes) != len(paramValues) {
		msg := fmt.Sprintf("unexpected number of parameters, expected %d, got %d", len(paramTypes), len(paramValues))
		lggr.Error(msg)

		return nil, errors.New(msg)
	}

	functionValues := make([]any, len(paramValues))
	for i, v := range paramValues {
		value, err := codec.EncodeToSuiValue(paramTypes[i], v)
		if err != nil {
			lggr.Errorf("failed to encode value: %v", err)
			return nil, err
		}

		functionValues[i] = value
	}

	rsp, err := suiClient.MoveCall(ctx, client.MoveCallRequest{
		Signer:          signerAddress,
		PackageObjectId: packageObjectId,
		Module:          moduleName,
		Function:        functionName,
		// We will only need to pass the type arguments if the function is generic
		TypeArguments: []any{},
		Arguments:     functionValues,
		GasBudget:     txMetadata.GasLimit.String(),
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

	signatures, err := signerService.Sign(txBytes)
	if err != nil {
		lggr.Errorf("Error signing transaction: %v", err)
	}

	return &SuiTx{
		TransactionID: transactionID,
		Sender:        signerAddress,
		Metadata:      &commontypes.TxMeta{},
		Timestamp:     GetCurrentUnixTimestamp(),
		Payload:       txBytes,
		Functions:     []*SuiFunction{function},
		Signatures:    signatures,
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
//   - lggr: Logger for recording operation details and errors.
//   - signerService: Service for signing the generated transaction bytes.
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
	lggr logger.Logger,
	signerService signer.SuiSigner,
	suiClient client.SuiPTBClient,
	requestType string,
	transactionID string,
	txMetadata *commontypes.TxMeta,
	signerAddress string,
	ptb *suiptb.ProgrammableTransaction,
	simulateTx bool,
) (*SuiTx, error) {
	address, err := sui.AddressFromHex(signerAddress)
	if err != nil {
		lggr.Errorf("failed to get address from hex: %v", err)
		return nil, err
	}

	// Define gasBudget outside the if statement
	var gasBudget uint64
	if txMetadata.GasLimit != nil {
		gasBudget = txMetadata.GasLimit.Uint64()
	} else {
		gasBudget = uint64(defaultGasBudget)
	}

	// Get all available coins for this address
	coinData, err := suiClient.GetCoinsByAddress(ctx, signerAddress)
	if err != nil {
		lggr.Errorf("failed to get coins by address: %v", err)
		return nil, err
	}

	// Use the new function to select coins for gas budget
	gasBudgetCoins, err := SelectCoinsForGasBudget(gasBudget, coinData)
	if err != nil {
		lggr.Errorf("failed to select coins for gas budget: %v", err)
		return nil, err
	}

	lggr.Debugw("Gas budget coins selected",
		"gasBudget", gasBudget,
		"numCoins", len(gasBudgetCoins),
		"gasBudgetCoins", gasBudgetCoins)

	// We can get the reference from the current epoch -> https://docs.sui.io/sui-api-ref#suix_getreferencegasprice
	// TODO: decide if we need this or not
	gasPrice := suiclient.DefaultGasPrice

	// Create transaction data
	txData := suiptb.NewTransactionData(
		address,
		(*ptb),
		gasBudgetCoins,
		gasBudget,
		gasPrice,
	)

	lggr.Debugw("Transaction data", "txData", txData)

	// Marshal the transaction data using BCS
	ptbBytes, err := bcs.Marshal(txData)
	if err != nil {
		lggr.Errorf("failed to marshal transaction data: %v", err)
		return nil, err
	}

	signatures, err := signerService.Sign(ptbBytes)
	if err != nil {
		lggr.Errorf("Error signing transaction: %v", err)
	}

	functions := []*SuiFunction{}
	for _, command := range ptb.Commands {
		functions = append(functions, &SuiFunction{
			PackageId: command.MoveCall.Package.String(),
			Module:    command.MoveCall.Module,
			Name:      command.MoveCall.Function,
		})
	}

	return &SuiTx{
		TransactionID: transactionID,
		Sender:        signerAddress,
		Metadata:      txMetadata,
		Timestamp:     GetCurrentUnixTimestamp(),
		Payload:       ptbBytes,
		Functions:     functions,
		Signatures:    signatures,
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
func SelectCoinsForGasBudget(gasBudget uint64, availableCoins []client.CoinData) ([]*sui.ObjectRef, error) {
	if len(availableCoins) == 0 {
		return nil, fmt.Errorf("no coins available for gas budget")
	}

	// Convert CoinData to a more workable format with parsed balances
	type coinInfo struct {
		coinData  client.CoinData
		balance   uint64
		objectRef *sui.ObjectRef
	}

	coins := make([]coinInfo, 0, len(availableCoins))
	for _, coin := range availableCoins {
		// Parse balance to uint64
		var balance uint64
		_, err := fmt.Sscanf(coin.Balance, "%d", &balance)
		if err != nil {
			return nil, fmt.Errorf("failed to parse coin balance: %w", err)
		}

		coinObjectId, err := sui.ObjectIdFromHex(coin.CoinObjectId)
		if err != nil {
			return nil, fmt.Errorf("failed to parse coin object id: %w", err)
		}

		digest, err := sui.NewDigest(coin.Digest)
		if err != nil {
			return nil, fmt.Errorf("failed to parse coin digest: %w", err)
		}

		// Convert version string to uint64
		var version uint64
		_, err = fmt.Sscanf(coin.Version, "%d", &version)
		if err != nil {
			return nil, fmt.Errorf("failed to parse coin version: %w", err)
		}

		ref := &sui.ObjectRef{
			ObjectId: coinObjectId,
			Version:  version,
			Digest:   digest,
		}

		coins = append(coins, coinInfo{
			coinData:  coin,
			balance:   balance,
			objectRef: ref,
		})
	}

	// Sort coins by balance in descending order
	sort.Slice(coins, func(i, j int) bool {
		return coins[i].balance > coins[j].balance
	})

	// First, check if there's a single coin that's close to the gas budget
	for _, coin := range coins {
		if coin.balance >= gasBudget {
			return []*sui.ObjectRef{coin.objectRef}, nil
		}
	}

	// If no single coin is sufficient, try to find a combination
	// This is a simplified greedy approach (not guaranteed optimal for all cases)
	selected := make([]*sui.ObjectRef, 0)
	var totalBalance uint64

	for _, coin := range coins {
		selected = append(selected, coin.objectRef)
		totalBalance += coin.balance

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
