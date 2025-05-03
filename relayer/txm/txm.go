package txm

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/pattonkan/sui-go/sui/suiptb"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	commonutils "github.com/smartcontractkit/chainlink-common/pkg/utils"

	"github.com/smartcontractkit/chainlink-sui/relayer/client"
)

const expectedFunctionTokens = 3
const numberGoroutines = 2

type TxManager interface {
	services.Service
	Enqueue(ctx context.Context, transactionID string, txMetadata *commontypes.TxMeta, signerPublicKey []byte, function string, typeArgs []string, paramTypes []string, paramValues []any, simulateTx bool) (*SuiTx, error)
	EnqueuePTB(ctx context.Context, transactionID string, txMetadata *commontypes.TxMeta, signerPublicKey []byte, ptb *suiptb.ProgrammableTransaction, simulateTx bool) (*SuiTx, error)
	GetTransactionStatus(ctx context.Context, transactionID string) (commontypes.TransactionStatus, error)
	GetClient() client.SuiPTBClient
}

type SuiTxm struct {
	lggr                  logger.Logger
	suiGateway            client.SuiPTBClient
	keystoreService       loop.Keystore
	transactionRepository TxmStore
	retryManager          RetryManager
	gasManager            GasManager
	configuration         Config
	Starter               commonutils.StartStopOnce
	done                  sync.WaitGroup
	broadcastChannel      chan string
	stopChannel           chan struct{}
}

func NewSuiTxm(
	lggr logger.Logger, gateway client.SuiPTBClient, k loop.Keystore,
	conf Config, transactionsRepository TxmStore,
	retryManager RetryManager, gasManager GasManager,
) (*SuiTxm, error) {
	return &SuiTxm{
		lggr:                  logger.Named(lggr, "SuiTxm"),
		suiGateway:            gateway,
		keystoreService:       k,
		transactionRepository: transactionsRepository,
		retryManager:          retryManager,
		gasManager:            gasManager,
		configuration:         conf,
		broadcastChannel:      make(chan string, conf.BroadcastChanSize),
		stopChannel:           make(chan struct{}),
	}, nil
}

// Enqueue generates a standard Move call transaction, adds it to the transaction store,
// and queues it for broadcasting. It handles transaction ID generation (if needed),
// function signature parsing, parameter encoding, transaction generation, signing, and storage.
//
// Parameters:
//   - ctx: Context for the operation.
//   - transactionID: A specific ID for the transaction. If empty, a new ID is generated.
//   - txMetadata: Transaction metadata, potentially including gas limits.
//   - signerPublicKey: The public key of the account signing and sending the transaction.
//   - function: The full function signature (e.g., "packageId::module::functionName").
//   - typeArgs: Type arguments for generic Move functions.
//   - paramTypes: The types of the parameters being passed to the Move function.
//   - paramValues: The actual values of the parameters.
//   - simulateTx: Boolean flag indicating whether to simulate the transaction (currently unused in GenerateTransaction).
//
// Returns:
//   - *SuiTx: The generated and stored transaction object.
//   - error: An error if the transaction ID exists, function parsing fails, transaction generation fails, or storage fails.
func (txm *SuiTxm) Enqueue(ctx context.Context, transactionID string, txMetadata *commontypes.TxMeta, signerPublicKey []byte, function string, typeArgs []string, paramTypes []string, paramValues []any, simulateTx bool) (*SuiTx, error) {
	if transactionID == "" {
		transactionID = TransactionIDGenerator()
	} else {
		// Check if the transaction ID already exists in the transactions map
		// If the transaction ID already exists, return an error
		_, err := txm.transactionRepository.GetTransaction(transactionID)
		if err == nil {
			return nil, errors.New("transaction already exists")
		}
	}

	functionTokens := strings.Split(function, "::")
	if len(functionTokens) != expectedFunctionTokens {
		msg := fmt.Sprintf("unexpected function name, expected 3 tokens, got %d", len(functionTokens))
		txm.lggr.Error(msg)

		return nil, errors.New(msg)
	}

	suiFunction := &SuiFunction{
		PackageId: functionTokens[0],
		Module:    functionTokens[1],
		Name:      functionTokens[2],
	}

	transaction, err := GenerateTransaction(
		ctx, signerPublicKey, txm.lggr, txm.keystoreService, txm.suiGateway,
		txm.configuration.RequestType, transactionID, txMetadata,
		suiFunction, typeArgs, paramTypes, paramValues,
	)
	if err != nil {
		txm.lggr.Errorw("Failed to generate transaction", "error", err)
		return nil, err
	}

	err = txm.transactionRepository.AddTransaction(*transaction)
	if err != nil {
		txm.lggr.Errorw("Failed to add transaction to repository", "error", err)
		return nil, err
	}

	txm.broadcastChannel <- transactionID
	txm.lggr.Infow("Transaction added to broadcast channel", "transactionID", transactionID)
	txm.lggr.Infow("Transaction enqueued", "transactionID", transactionID)

	return transaction, nil
}

// EnqueuePTB generates a transaction based on a pre-constructed Programmable Transaction Block (PTB),
// adds it to the transaction store, and queues it for broadcasting.
// It determines gas limits, selects gas coins, signs the transaction, and stores it.
// It's part of the TxManager interface implementation.
//
// Parameters:
//   - ctx: Context for the operation.
//   - transactionID: Unique identifier for the transaction.
//   - txMetadata: Transaction metadata, potentially including gas limits.
//   - signerPublicKey: The public key of the account signing and sending the transaction.
//   - ptb: The ProgrammableTransaction block containing the sequence of commands.
//   - simulateTx: Boolean flag indicating whether to simulate the transaction (currently unused).
//
// Returns:
//   - *SuiTx: The generated and stored transaction object.
//   - error: An error if transaction generation or storage fails.
func (txm *SuiTxm) EnqueuePTB(ctx context.Context, transactionID string, txMetadata *commontypes.TxMeta, signerPublicKey []byte, ptb *suiptb.ProgrammableTransaction, simulateTx bool) (*SuiTx, error) {
	txm.lggr.Infow("Enqueuing PTB", "transactionID", transactionID, "ptb", ptb)

	transaction, err := GeneratePTBTransaction(
		ctx, signerPublicKey, txm.lggr, txm.keystoreService, txm.suiGateway,
		txm.configuration.RequestType, transactionID, txMetadata,
		ptb, simulateTx,
	)
	if err != nil {
		txm.lggr.Errorw("Failed to generate PTB transaction", "error", err)
		return nil, err
	}

	txm.lggr.Infow("PTB transaction generated", "transactionID", transactionID, "ptb", transaction)

	err = txm.transactionRepository.AddTransaction(*transaction)
	if err != nil {
		txm.lggr.Errorw("Failed to add transaction to repository", "error", err)
		return nil, err
	}

	txm.broadcastChannel <- transactionID
	txm.lggr.Infow("PTB Transaction added to broadcast channel", "transactionID", transactionID)
	txm.lggr.Infow("PTB Transaction enqueued", "transactionID", transactionID)

	return transaction, nil
}

// GetTransactionStatus implements TxManager.
func (txm *SuiTxm) GetTransactionStatus(ctx context.Context, transactionID string) (commontypes.TransactionStatus, error) {
	tx, err := txm.transactionRepository.GetTransaction(transactionID)
	if err != nil {
		txm.lggr.Errorw("Failed to get transaction", "transactionID", transactionID, "error", err)
		return commontypes.Unknown, err
	}

	switch tx.State {
	case StatePending:
		txm.lggr.Infow("Transaction is pending", "transactionID", transactionID)
		return commontypes.Pending, nil
	case StateSubmitted:
		txm.lggr.Infow("Transaction is submitted", "transactionID", transactionID)
		return commontypes.Unconfirmed, nil
	case StateFinalized:
		txm.lggr.Infow("Transaction is finalized", "transactionID", transactionID)
		return commontypes.Finalized, nil
	case StateRetriable:
		txm.lggr.Infow("Transaction is retriable", "transactionID", transactionID)
		return commontypes.Failed, nil
	case StateFailed:
		txm.lggr.Infow("Transaction has failed", "transactionID", transactionID)
		return commontypes.Fatal, nil
	default:
		txm.lggr.Errorw("Unknown transaction state", "transactionID", transactionID, "state", tx.State)
		return commontypes.Unknown, errors.New("unknown transaction state")
	}
}

func (txm *SuiTxm) Close() error {
	return txm.Starter.StopOnce("SuiTxm", func() error {
		txm.lggr.Infow("Closing SuiTxm")
		close(txm.stopChannel)
		txm.done.Wait()
		txm.lggr.Infow("SuiTxm closed")

		return nil
	})
}

// TODO: implement
func (txm *SuiTxm) HealthReport() map[string]error {
	return map[string]error{txm.Name(): txm.Starter.Healthy()}
}

// TODO: implement
func (txm *SuiTxm) Name() string {
	return txm.lggr.Name()
}

// TODO: implement
func (txm *SuiTxm) Ready() error {
	return txm.Starter.Ready()
}

func (txm *SuiTxm) Start(ctx context.Context) error {
	return txm.Starter.StartOnce("SuiTxm", func() error {
		txm.lggr.Infow("Starting SuiTxm")
		txm.done.Add(numberGoroutines) // waitgroup: broadcaster, confirmer
		go txm.broadcastLoop(ctx)
		go txm.confirmerLoop(ctx)

		return nil
	})
}

// GetClient returns the Sui client instance used by the transaction manager.
func (txm *SuiTxm) GetClient() client.SuiPTBClient {
	return txm.suiGateway
}

var _ TxManager = (*SuiTxm)(nil)
