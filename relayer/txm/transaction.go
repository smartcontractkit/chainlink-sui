package txm

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/google/uuid"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
	"github.com/smartcontractkit/chainlink-sui/relayer/signer"
)

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
	Signatures    []string
	RequestType   string
	Attempt       int
	State         TransactionState
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
//   - RequestType: "Call"
//   - Timestamp: Current UTC timestamp
//   - Payload: The BCS-serialized transaction bytes
//   - Signatures: Array of signatures produced by the signer service
func GenerateTransaction(
	ctx context.Context,
	lggr logger.Logger,
	signerService signer.SuiSigner,
	suiClient client.SuiClient,
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

	rsp, err := suiClient.MoveCall(ctx, models.MoveCallRequest{
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
		Signatures:    signatures,
		RequestType:   "Call",
		Attempt:       0,
		State:         StatePending,
	}, nil
}
