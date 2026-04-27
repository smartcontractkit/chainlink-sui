package bind

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/block-vision/sui-go-sdk/transaction"

	bindutils "github.com/smartcontractkit/chainlink-sui/bindings/utils"
)

// Polling parameters for waitForTransactionIndexed. Exposed as package-level vars so
// callers / tests can tune them without changing the SignAndSendTx signature.
var (
	// WaitForTxIndexedTimeout bounds the total time SignAndSendTx will poll for a
	// successfully-submitted transaction to become visible on the fullnode's JSON-RPC
	// view (i.e. queryable via sui_getTransactionBlock). This is an additional wall
	// clock cost on top of the transaction's own execution time, paid only when the
	// caller sets CallOpts.WaitForExecution = true.
	WaitForTxIndexedTimeout = 30 * time.Second

	// WaitForTxIndexedInitialBackoff is the first delay between poll attempts; it
	// doubles on each failure up to WaitForTxIndexedMaxBackoff.
	WaitForTxIndexedInitialBackoff = 100 * time.Millisecond

	// WaitForTxIndexedMaxBackoff caps the poll interval so steady-state polling does
	// not add unbounded latency for slow-indexing fullnodes.
	WaitForTxIndexedMaxBackoff = 1 * time.Second
)

// ErrTxIndexingTimeout is returned by SignAndSendTx when the transaction was accepted
// by validators (effects certified) but the fullnode JSON-RPC did not surface the tx
// within WaitForTxIndexedTimeout. The caller can distinguish this from a real tx
// failure (the tx is fine; just not yet visible for reads).
var ErrTxIndexingTimeout = errors.New("tx submitted but fullnode indexing wait timed out")

// SuiTransactionBlockGetter is the minimal JSON-RPC surface required to poll for a
// transaction digest via sui_getTransactionBlock. Implementations such as
// github.com/block-vision/sui-go-sdk/sui.ISuiAPI satisfy it.
type SuiTransactionBlockGetter interface {
	SuiGetTransactionBlock(ctx context.Context, req models.SuiGetTransactionBlockRequest) (models.SuiTransactionBlockResponse, error)
}

type Object struct {
	Id                   string
	InitialSharedVersion *uint64
}

// EmptyMoveStructWitness is a placeholder for a zero-sized Move struct (typical `has drop` proof
// witnesses). It encodes as an empty Pure argument (empty BCS for no fields).
type EmptyMoveStructWitness struct{}

type CallOpts struct {
	Signer           bindutils.SuiSigner
	GasObject        string
	GasBudget        *uint64
	GasPrice         *uint64
	WaitForExecution bool

	ObjectResolver *ObjectResolver
}

func SignAndSendTx(ctx context.Context, signer bindutils.SuiSigner, client sui.ISuiAPI, txBytes []byte, waitForExecution bool) (*models.SuiTransactionBlockResponse, error) {
	signatures, err := signer.Sign(txBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to sign tx: %w", err)
	}

	b64bytes := bindutils.EncodeBase64(txBytes)

	// Convert signatures to string array
	signatureStrings := make([]string, 0, len(signatures))
	signatureStrings = append(signatureStrings, signatures...)

	// We always use WaitForEffectsCert on the server side. WaitForLocalExecution has
	// been silently ignored at the JSON-RPC layer since Sui v1.33
	// (https://forums.sui.io/t/deprecating-waitforlocalexecution/45988), so setting it
	// provides no additional consistency guarantee. When the caller actually needs
	// read-after-write consistency (waitForExecution = true), we emulate it below by
	// polling sui_getTransactionBlock until the fullnode surfaces the tx — the same
	// approach the Typescript SDK takes in waitForTransaction:
	// https://github.com/MystenLabs/ts-sdks/blob/502ad7f2803bf6443f7cb000c802d78110585b6f/packages/typescript/src/experimental/core.ts#L114
	blockReq := models.SuiExecuteTransactionBlockRequest{
		TxBytes:   b64bytes,
		Signature: signatureStrings,
		Options: models.SuiTransactionBlockOptions{
			ShowInput:          true,
			ShowRawInput:       true,
			ShowEffects:        true,
			ShowObjectChanges:  true,
			ShowBalanceChanges: true,
			ShowEvents:         true,
		},
		RequestType: "WaitForEffectsCert",
	}

	tx, err := client.SuiExecuteTransactionBlock(ctx, blockReq)
	if err != nil {
		msg := fmt.Errorf("tx failed calling move method: %w", err)
		return nil, msg
	}

	if err := GetFailedTxError(&tx); err != nil {
		return &tx, err
	}

	if waitForExecution && tx.Digest != "" {
		if waitErr := WaitForTransactionIndexed(ctx, client, tx.Digest); waitErr != nil {
			// The tx itself succeeded on-chain; surface the indexing wait failure so
			// callers can decide whether to retry their follow-up RPC calls. We keep
			// the tx response so the caller still has access to effects / digest.
			return &tx, waitErr
		}
	}

	return &tx, nil
}

// WaitForTransactionIndexed polls sui_getTransactionBlock until the fullnode surfaces
// the given digest, providing the read-after-write consistency that
// WaitForLocalExecution used to give before being silently disabled in JSON-RPC
// (Sui >= v1.33). Without this, a tight "tx A -> tx B referencing objects mutated by
// tx A" sequence can race with the fullnode's indexer and have tx B rejected pre-
// consensus with "Object ... Version ... is not available for consumption" because the
// local view of owned-object versions (notably the default gas coin) is stale.
//
// The poll uses exponential backoff between WaitForTxIndexedInitialBackoff and
// WaitForTxIndexedMaxBackoff, bounded by WaitForTxIndexedTimeout.
//
// Relayer and other packages that execute with WaitForEffectsCert but still need the
// same consistency should call this after a successful execute when they would have
// previously used WaitForLocalExecution.
func WaitForTransactionIndexed(ctx context.Context, client SuiTransactionBlockGetter, digest string) error {
	pollCtx, cancel := context.WithTimeout(ctx, WaitForTxIndexedTimeout)
	defer cancel()

	req := models.SuiGetTransactionBlockRequest{
		Digest: digest,
		Options: models.SuiTransactionBlockOptions{
			ShowEffects: true,
		},
	}

	backoff := WaitForTxIndexedInitialBackoff
	var lastErr error
	for {
		resp, err := client.SuiGetTransactionBlock(pollCtx, req)
		if err == nil && resp.Digest == digest {
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

func DevInspectTx(ctx context.Context, signerAddress string, client sui.ISuiAPI, txBytes []byte) (*models.SuiTransactionBlockResponse, error) {
	b64bytes := bindutils.EncodeBase64(txBytes)

	devInspectReq := models.SuiDevInspectTransactionBlockRequest{
		Sender:  signerAddress,
		TxBytes: b64bytes,
	}

	tx, err := client.SuiDevInspectTransactionBlock(ctx, devInspectReq)
	if err != nil {
		msg := fmt.Errorf("tx failed calling dev inspect method: %w", err)
		return nil, msg
	}

	return &tx, nil
}

// DevInspectPTB executes a PTB using DevInspect
func DevInspectPTB(ctx context.Context, signerAddress string, client sui.ISuiAPI, ptb *transaction.Transaction) (*models.SuiTransactionBlockResponse, error) {
	// ensure the PTB has the required data
	if ptb.Data.V1 == nil || ptb.Data.V1.Kind == nil {
		return nil, fmt.Errorf("PTB is not properly initialized")
	}

	// at this stage, we do not have any type information, and all unresolved variants should be resolved.
	if ptb.Data.V1.Kind.ProgrammableTransaction != nil && len(ptb.Data.V1.Kind.ProgrammableTransaction.Inputs) > 0 {
		for _, input := range ptb.Data.V1.Kind.ProgrammableTransaction.Inputs {
			if input.UnresolvedPure != nil {
				return nil, fmt.Errorf("UnresolvedPure found in PTB inputs")
			}
			if input.UnresolvedObject != nil {
				return nil, fmt.Errorf("UnresolvedObject found in PTB inputs")
			}
		}
	}

	txBytes, err := ptb.Data.V1.Kind.Marshal()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transaction kind: %w", err)
	}

	b64TxBytes := bindutils.EncodeBase64(txBytes)

	devInspectReq := models.SuiDevInspectTransactionBlockRequest{
		Sender:  signerAddress,
		TxBytes: b64TxBytes,
	}

	tx, err := client.SuiDevInspectTransactionBlock(ctx, devInspectReq)
	if err != nil {
		msg := fmt.Errorf("tx failed calling dev inspect method: %w", err)
		return nil, msg
	}

	return &tx, nil
}
