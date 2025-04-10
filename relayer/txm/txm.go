package txm

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/keystore"
	"github.com/smartcontractkit/chainlink-sui/relayer/signer"
)

const expectedFunctionTokens = 3
const numberGoroutines = 2

type TxManager interface {
	services.Service
	Enqueue(ctx context.Context, transactionID string, txMetadata *commontypes.TxMeta, signerAddress, function string, typeArgs []string, paramTypes []string, paramValues []any, simulateTx bool) (*SuiTx, error)
	GetTransactionStatus(ctx context.Context, transactionID string) (commontypes.TransactionStatus, error)
}

type SuiTxm struct {
	lggr                  logger.Logger
	suiGateway            client.SuiClient
	keyStoreRepository    keystore.Keystore
	transactionRepository TxmStore
	configuration         Config
	signer                signer.SuiSigner
	done                  sync.WaitGroup
	broadcastChannel      chan string
	stopChannel           chan struct{}
}

var _ TxManager = (*SuiTxm)(nil)

func NewSuiTxm(lggr logger.Logger, gateway client.SuiClient, k keystore.Keystore, conf Config, sig signer.SuiSigner, transactionsRepository TxmStore) (*SuiTxm, error) {
	return &SuiTxm{
		lggr:                  logger.Named(lggr, "SuiTxm"),
		suiGateway:            gateway,
		keyStoreRepository:    k,
		transactionRepository: transactionsRepository,
		configuration:         conf,
		signer:                sig,
		broadcastChannel:      make(chan string, conf.BroadcastChanSize),
		stopChannel:           make(chan struct{}),
	}, nil
}

func (txm *SuiTxm) Enqueue(ctx context.Context, transactionID string, txMetadata *commontypes.TxMeta, signerAddress, function string, typeArgs []string, paramTypes []string, paramValues []any, simulateTx bool) (*SuiTx, error) {
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
		ctx, txm.lggr, txm.signer, txm.suiGateway,
		txm.configuration.RequestType, transactionID, txMetadata, signerAddress,
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
	txm.lggr.Infow("Closing SuiTxm")
	txm.stopChannel <- struct{}{}
	txm.lggr.Infow("SuiTxm closed")

	return nil
}

// TODO: implement
func (txm *SuiTxm) HealthReport() map[string]error {
	panic("unimplemented")
}

// TODO: implement
func (txm *SuiTxm) Name() string {
	panic("unimplemented")
}

// TODO: implement
func (txm *SuiTxm) Ready() error {
	panic("unimplemented")
}

func (txm *SuiTxm) Start(ctx context.Context) error {
	txm.lggr.Infow("Starting SuiTxm")
	txm.done.Add(numberGoroutines) // waitgroup: broadcaster, confirmer
	go txm.broadcastLoop(ctx)
	go txm.confimerLoop(ctx)

	return nil
}
