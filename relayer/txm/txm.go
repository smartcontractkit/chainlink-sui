package txm

import (
	"context"
	"errors"
	"sync"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	commonutils "github.com/smartcontractkit/chainlink-common/pkg/utils"

	"github.com/block-vision/sui-go-sdk/transaction"

	"github.com/smartcontractkit/chainlink-sui/relayer/client"
)

const numberGoroutines = 2

type TxManager interface {
	services.Service
	EnqueuePTB(ctx context.Context, transactionID string, txMetadata *commontypes.TxMeta, signerPublicKey []byte, ptb *transaction.Transaction) (*SuiTx, error)
	GetTransactionStatus(ctx context.Context, transactionID string) (commontypes.TransactionStatus, error)
	GetClient() client.SuiPTBClient
	GetGasManager() GasManager
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
	lggr.Infof("Creating SuiTxm")
	lggr.Infof("SuiTxm configuration: %+v", conf)
	lggr.Infof("Gas manager Max Gas Budget: %+v", gasManager.MaxGasBudget())

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
func (txm *SuiTxm) EnqueuePTB(ctx context.Context, transactionID string, txMetadata *commontypes.TxMeta, signerPublicKey []byte, ptb *transaction.Transaction) (*SuiTx, error) {
	txm.lggr.Infow("Enqueuing PTB", "transactionID", transactionID, "ptb", ptb)

	simulateTx := true

	txn, err := GeneratePTBTransactionWithGasEstimation(
		ctx, signerPublicKey, txm.lggr, txm.keystoreService, txm.suiGateway,
		txm.configuration.RequestType, transactionID, txMetadata,
		ptb, simulateTx, txm.gasManager,
	)
	if err != nil {
		txm.lggr.Errorw("Failed to generate PTB txn", "error", err)
		return nil, err
	}

	txm.lggr.Infow("PTB txn generated", "transactionID", transactionID, "ptb", txn)

	err = txm.transactionRepository.AddTransaction(*txn)
	if err != nil {
		txm.lggr.Errorw("Failed to add txn to repository", "error", err)
		return nil, err
	}

	txm.broadcastChannel <- transactionID
	txm.lggr.Infow("PTB Transaction added to broadcast channel", "transactionID", transactionID)
	txm.lggr.Infow("PTB Transaction enqueued", "transactionID", transactionID)

	return txn, nil
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

func (txm *SuiTxm) HealthReport() map[string]error {
	return map[string]error{txm.Name(): txm.Starter.Healthy()}
}

func (txm *SuiTxm) Name() string {
	return txm.lggr.Name()
}

func (txm *SuiTxm) Ready() error {
	return txm.Starter.Ready()
}

func (txm *SuiTxm) Start(_ context.Context) error {
	//nolint:contextcheck
	return txm.Starter.StartOnce("SuiTxm", func() error {
		txm.lggr.Infow("Starting SuiTxm")
		txm.done.Add(numberGoroutines) // waitgroup: broadcaster, confirmer
		go txm.broadcastLoop()
		go txm.confirmerLoop()

		return nil
	})
}

// GetClient returns the Sui client instance used by the transaction manager.
func (txm *SuiTxm) GetClient() client.SuiPTBClient {
	return txm.suiGateway
}

func (txm *SuiTxm) GetGasManager() GasManager {
	return txm.gasManager
}

var _ TxManager = (*SuiTxm)(nil)
