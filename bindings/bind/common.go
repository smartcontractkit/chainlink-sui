package bind

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/block-vision/sui-go-sdk/models"
	suirpcv2 "github.com/block-vision/sui-go-sdk/pb/sui/rpc/v2"
	"github.com/block-vision/sui-go-sdk/signer"
	"github.com/block-vision/sui-go-sdk/transaction"

	bindutils "github.com/smartcontractkit/chainlink-sui/bindings/utils"
	"github.com/smartcontractkit/chainlink-sui/codec"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
)

// Polling parameters for waitForTransactionIndexed.
var (
	WaitForTxIndexedTimeout        = 30 * time.Second
	WaitForTxIndexedInitialBackoff = 100 * time.Millisecond
	WaitForTxIndexedMaxBackoff     = 1 * time.Second
)

var ErrTxIndexingTimeout = errors.New("tx submitted but fullnode indexing wait timed out")

type transactionStatusClient interface {
	GetTransactionStatus(ctx context.Context, digest string) (client.TransactionResult, error)
}

type transactionObjectChangesClient interface {
	GetTransactionChangedObjects(ctx context.Context, digest string) ([]*suirpcv2.ChangedObject, error)
}

// Object is an alias for codec.Object for backward compatibility.
//
// Deprecated: use codec.Object directly.
type Object = codec.Object

type EmptyMoveStructWitness struct{}

type CallOpts struct {
	Signer           bindutils.SuiSigner
	GasObject        string
	GasBudget        *uint64
	GasPrice         *uint64
	WaitForExecution bool

	ObjectResolver *ObjectResolver
}

func SignAndSendTx(ctx context.Context, signer bindutils.SuiSigner, chainClient client.BindingsClient, txBytes []byte, waitForExecution bool) (*models.SuiTransactionBlockResponse, error) {
	signatures, err := signer.Sign(txBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to sign tx: %w", err)
	}

	execReq, err := buildExecuteRequest(txBytes, signatures)
	if err != nil {
		return nil, fmt.Errorf("failed to build execute request: %w", err)
	}

	resp, err := chainClient.SendTransaction(ctx, execReq)
	if err != nil {
		return nil, fmt.Errorf("tx failed calling move method: %w", err)
	}

	tx, err := mapExecuteResponseToModels(resp)
	if err != nil {
		return nil, err
	}

	if tx.Digest != "" {
		if fetcher, ok := chainClient.(transactionObjectChangesClient); ok {
			if changed, fetchErr := fetcher.GetTransactionChangedObjects(ctx, tx.Digest); fetchErr == nil && len(changed) > 0 {
				tx.ObjectChanges = mapChangedObjectsToModels(changed)
			}
		}
	}

	if err := GetFailedTxError(tx); err != nil {
		return tx, err
	}

	if waitForExecution && tx.Digest != "" {
		if waitErr := WaitForTransactionIndexed(ctx, chainClient, tx.Digest); waitErr != nil {
			return tx, waitErr
		}
	}

	return tx, nil
}

func WaitForTransactionIndexed(ctx context.Context, chainClient transactionStatusClient, digest string) error {
	pollCtx, cancel := context.WithTimeout(ctx, WaitForTxIndexedTimeout)
	defer cancel()

	backoff := WaitForTxIndexedInitialBackoff
	var lastErr error
	for {
		result, err := chainClient.GetTransactionStatus(pollCtx, digest)
		if err == nil && result.Status == "success" {
			return nil
		}
		lastErr = err

		select {
		case <-pollCtx.Done():
			if lastErr != nil {
				return fmt.Errorf("%w (digest=%s): %w", ErrTxIndexingTimeout, digest, lastErr)
			}
			return fmt.Errorf("%w (digest=%s)", ErrTxIndexingTimeout, digest)
		case <-time.After(backoff):
		}

		if backoff < WaitForTxIndexedMaxBackoff {
			backoff *= 2
			if backoff > WaitForTxIndexedMaxBackoff {
				backoff = WaitForTxIndexedMaxBackoff
			}
		}
	}
}

func buildSimulateBCS(ctx context.Context, chainClient client.BindingsClient, ptb *transaction.Transaction, opts *CallOpts) ([]byte, error) {
	if ptb.Data.V1.GasData.Budget == nil {
		budget := DefaultGasBudget
		if opts != nil && opts.GasBudget != nil {
			budget = *opts.GasBudget
		}
		ptb.SetGasBudget(budget)
	}

	if ptb.Data.V1.GasData.Price == nil {
		var gasPrice uint64
		if opts != nil && opts.GasPrice != nil {
			gasPrice = *opts.GasPrice
		} else {
			price, err := chainClient.GetReferenceGasPrice(ctx)
			if err != nil {
				gasPrice = defaultGasPrice
			} else {
				gasPrice = price.Uint64()
			}
		}
		ptb.SetGasPrice(gasPrice)
	}

	if ptb.Data.V1.Sender != nil {
		addr := transaction.ConvertSuiAddressBytesToString(*ptb.Data.V1.Sender)
		ptb.SetSigner(&signer.Signer{Address: string(addr)})
	}

	return ptb.BuildBCSBytes(ctx)
}
