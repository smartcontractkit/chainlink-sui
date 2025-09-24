package events

import (
	"context"
	"fmt"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/block-vision/sui-go-sdk/transaction"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	rel "github.com/smartcontractkit/chainlink-sui/relayer/signer"
)

type RMNRemoteEmitter interface {
	EmitEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool, eventName string) (*models.SuiTransactionBlockResponse, error)
	BatchEmitEvents(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction) (*models.SuiTransactionBlockResponse, error)
	EmitConfigSetEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool, version uint32) (*models.SuiTransactionBlockResponse, error)
	EmitCursedEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool) (*models.SuiTransactionBlockResponse, error)
	EmitCursedMultipleEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool) (*models.SuiTransactionBlockResponse, error)
	EmitUncursedEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool) (*models.SuiTransactionBlockResponse, error)
	EmitUncursedMultipleEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool) (*models.SuiTransactionBlockResponse, error)
}

type RMNRemoteEmitterImpl struct {
	lggr           logger.Logger
	packageId      string
	boundContract  *bind.BoundContract
	callOpts       bind.CallOpts
	client         sui.ISuiAPI
	accountAddress string
}

func NewRMNRemoteEmitter(
	lggr logger.Logger,
	packageId string,
	signer rel.SuiSigner,
	callOpts bind.CallOpts,
	client sui.ISuiAPI,
	accountAddress string,
) *RMNRemoteEmitterImpl {
	boundContract, err := bind.NewBoundContract(
		packageId,
		"test",
		"rmn_remote",
		client,
	)
	if err != nil {
		lggr.Errorw("Failed to create rmn_remote bound contract", "error", err)
		return nil
	}

	return &RMNRemoteEmitterImpl{
		lggr:           lggr,
		packageId:      packageId,
		boundContract:  boundContract,
		callOpts:       callOpts,
		client:         client,
		accountAddress: accountAddress,
	}
}

// executePTBCall handles the common PTB execution logic for event emission
func (e *RMNRemoteEmitterImpl) executePTBCall(ctx context.Context, ptb *transaction.Transaction, addToPTB bool, encodedCall *bind.EncodedCall, operationName string) (*models.SuiTransactionBlockResponse, error) {
	if addToPTB && ptb != nil {
		_, err := e.boundContract.AppendPTB(ctx, &e.callOpts, ptb, encodedCall)
		if err != nil {
			e.lggr.Errorw("Failed to append PTB", "error", err, "operation", operationName)
			return nil, err
		}
		return nil, nil
	}

	if ptb == nil {
		ptb = transaction.NewTransaction()
	}

	_, err := e.boundContract.AppendPTB(ctx, &e.callOpts, ptb, encodedCall)
	if err != nil {
		e.lggr.Errorw("Failed to append PTB", "error", err, "operation", operationName)
		return nil, err
	}

	tx, err := bind.ExecutePTB(ctx, &e.callOpts, e.client, ptb)
	if err != nil {
		e.lggr.Errorw("Failed to execute PTB", "error", err, "operation", operationName)
		return nil, err
	}

	e.lggr.Infow("Executed PTB", "tx", tx, "operation", operationName)

	return tx, nil
}

func (e *RMNRemoteEmitterImpl) EmitEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool, eventName string) (*models.SuiTransactionBlockResponse, error) {
	e.lggr.Infow("Emitting event", "eventName", eventName)

	if ptb == nil {
		ptb = transaction.NewTransaction()
		addToPTB = false
	}

	switch eventName {
	case "emit_config_set_event":
		return e.EmitConfigSetEvent(ctx, gasBudget, ptb, addToPTB, 1)
	case "emit_cursed_event":
		return e.EmitCursedEvent(ctx, gasBudget, ptb, addToPTB)
	case "emit_cursed_multiple_event":
		return e.EmitCursedMultipleEvent(ctx, gasBudget, ptb, addToPTB)
	case "emit_uncursed_event":
		return e.EmitUncursedEvent(ctx, gasBudget, ptb, addToPTB)
	case "emit_uncursed_multiple_event":
		return e.EmitUncursedMultipleEvent(ctx, gasBudget, ptb, addToPTB)
	default:
		e.lggr.Errorw("Invalid event name", "eventName", eventName)
		return nil, fmt.Errorf("invalid event name: %s", eventName)
	}
}

func (e *RMNRemoteEmitterImpl) BatchEmitEvents(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction) (*models.SuiTransactionBlockResponse, error) {
	if ptb == nil {
		ptb = transaction.NewTransaction()
	}

	e.lggr.Infow("Emitting batch of events for RMN remote")

	e.EmitConfigSetEvent(ctx, gasBudget, ptb, true, 1)
	e.EmitCursedEvent(ctx, gasBudget, ptb, true)
	e.EmitCursedMultipleEvent(ctx, gasBudget, ptb, true)
	e.EmitUncursedEvent(ctx, gasBudget, ptb, true)
	e.EmitUncursedMultipleEvent(ctx, gasBudget, ptb, true)

	tx, err := bind.ExecutePTB(ctx, &e.callOpts, e.client, ptb)
	if err != nil {
		e.lggr.Errorw("Failed to execute PTB", "error", err)
		return nil, err
	}
	return tx, nil
}

func (e *RMNRemoteEmitterImpl) EmitConfigSetEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool, version uint32) (*models.SuiTransactionBlockResponse, error) {
	e.lggr.Infow("Emitting config set event", "version", version)
	emitConfigSetEvent, err := e.boundContract.EncodeCallArgs(
		"emit_config_set_event",
		[]string{},
		[]string{"u32"},
		[]any{version},
	)
	if err != nil {
		e.lggr.Errorw("Failed to encode call", "error", err)
		return nil, err
	}

	return e.executePTBCall(ctx, ptb, addToPTB, emitConfigSetEvent, "emit_config_set_event")
}

func (e *RMNRemoteEmitterImpl) EmitCursedEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool) (*models.SuiTransactionBlockResponse, error) {
	e.lggr.Infow("Emitting cursed event")
	emitCursedEvent, err := e.boundContract.EncodeCallArgs(
		"emit_cursed_event",
		[]string{},
		[]string{},
		nil,
	)
	if err != nil {
		e.lggr.Errorw("Failed to encode call", "error", err)
		return nil, err
	}

	return e.executePTBCall(ctx, ptb, addToPTB, emitCursedEvent, "emit_cursed_event")
}

func (e *RMNRemoteEmitterImpl) EmitCursedMultipleEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool) (*models.SuiTransactionBlockResponse, error) {
	e.lggr.Infow("Emitting cursed multiple event")
	emitCursedMultipleEvent, err := e.boundContract.EncodeCallArgs(
		"emit_cursed_multiple_event",
		[]string{},
		[]string{},
		nil,
	)
	if err != nil {
		e.lggr.Errorw("Failed to encode call", "error", err)
		return nil, err
	}

	return e.executePTBCall(ctx, ptb, addToPTB, emitCursedMultipleEvent, "emit_cursed_multiple_event")
}

func (e *RMNRemoteEmitterImpl) EmitUncursedEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool) (*models.SuiTransactionBlockResponse, error) {
	e.lggr.Infow("Emitting uncursed event")
	emitUncursedEvent, err := e.boundContract.EncodeCallArgs(
		"emit_uncursed_event",
		[]string{},
		[]string{},
		nil,
	)
	if err != nil {
		e.lggr.Errorw("Failed to encode call", "error", err)
		return nil, err
	}

	return e.executePTBCall(ctx, ptb, addToPTB, emitUncursedEvent, "emit_uncursed_event")
}

func (e *RMNRemoteEmitterImpl) EmitUncursedMultipleEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool) (*models.SuiTransactionBlockResponse, error) {
	e.lggr.Infow("Emitting uncursed multiple event")
	emitUncursedMultipleEvent, err := e.boundContract.EncodeCallArgs(
		"emit_uncursed_multiple_event",
		[]string{},
		[]string{},
		nil,
	)
	if err != nil {
		e.lggr.Errorw("Failed to encode call", "error", err)
		return nil, err
	}

	return e.executePTBCall(ctx, ptb, addToPTB, emitUncursedMultipleEvent, "emit_uncursed_multiple_event")
}
