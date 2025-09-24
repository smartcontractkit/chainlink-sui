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

type TokenAdminRegistryEmitter interface {
	EmitEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool, eventName string) (*models.SuiTransactionBlockResponse, error)
	BatchEmitEvents(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction) (*models.SuiTransactionBlockResponse, error)
	EmitPoolRegisteredEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool) (*models.SuiTransactionBlockResponse, error)
	EmitAdministratorTransferredEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool) (*models.SuiTransactionBlockResponse, error)
	EmitPoolSetEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool) (*models.SuiTransactionBlockResponse, error)
	EmitPoolUnregisteredEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool) (*models.SuiTransactionBlockResponse, error)
	EmitAdministratorTransferRequestedEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool) (*models.SuiTransactionBlockResponse, error)
}

type TokenAdminRegistryEmitterImpl struct {
	lggr           logger.Logger
	packageId      string
	boundContract  *bind.BoundContract
	callOpts       bind.CallOpts
	client         sui.ISuiAPI
	accountAddress string
}

func NewTokenAdminRegistryEmitter(
	lggr logger.Logger,
	packageId string,
	signer rel.SuiSigner,
	callOpts bind.CallOpts,
	client sui.ISuiAPI,
	accountAddress string,
) *TokenAdminRegistryEmitterImpl {
	boundContract, err := bind.NewBoundContract(
		packageId,
		"test",
		"token_admin_registry",
		client,
	)
	if err != nil {
		lggr.Errorw("Failed to create token_admin_registry bound contract", "error", err)
		return nil
	}

	return &TokenAdminRegistryEmitterImpl{
		lggr:           lggr,
		packageId:      packageId,
		boundContract:  boundContract,
		callOpts:       callOpts,
		client:         client,
		accountAddress: accountAddress,
	}
}

// executePTBCall handles the common PTB execution logic for event emission
func (e *TokenAdminRegistryEmitterImpl) executePTBCall(ctx context.Context, ptb *transaction.Transaction, addToPTB bool, encodedCall *bind.EncodedCall, operationName string) (*models.SuiTransactionBlockResponse, error) {
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

func (e *TokenAdminRegistryEmitterImpl) EmitEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool, eventName string) (*models.SuiTransactionBlockResponse, error) {
	e.lggr.Infow("Emitting event", "eventName", eventName)

	if ptb == nil {
		ptb = transaction.NewTransaction()
		addToPTB = false
	}

	switch eventName {
	case "emit_pool_registered_event":
		return e.EmitPoolRegisteredEvent(ctx, gasBudget, ptb, addToPTB)
	case "emit_administrator_transferred_event":
		return e.EmitAdministratorTransferredEvent(ctx, gasBudget, ptb, addToPTB)
	case "emit_pool_set_event":
		return e.EmitPoolSetEvent(ctx, gasBudget, ptb, addToPTB)
	case "emit_pool_unregistered_event":
		return e.EmitPoolUnregisteredEvent(ctx, gasBudget, ptb, addToPTB)
	case "emit_administrator_transfer_requested_event":
		return e.EmitAdministratorTransferRequestedEvent(ctx, gasBudget, ptb, addToPTB)
	default:
		e.lggr.Errorw("Invalid event name", "eventName", eventName)
		return nil, fmt.Errorf("invalid event name: %s", eventName)
	}
}

func (e *TokenAdminRegistryEmitterImpl) BatchEmitEvents(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction) (*models.SuiTransactionBlockResponse, error) {
	if ptb == nil {
		ptb = transaction.NewTransaction()
	}

	e.lggr.Infow("Emitting batch of events for token admin registry")

	e.EmitPoolRegisteredEvent(ctx, gasBudget, ptb, true)
	e.EmitAdministratorTransferredEvent(ctx, gasBudget, ptb, true)
	e.EmitPoolSetEvent(ctx, gasBudget, ptb, true)
	e.EmitPoolUnregisteredEvent(ctx, gasBudget, ptb, true)
	e.EmitAdministratorTransferRequestedEvent(ctx, gasBudget, ptb, true)

	tx, err := bind.ExecutePTB(ctx, &e.callOpts, e.client, ptb)
	if err != nil {
		e.lggr.Errorw("Failed to execute PTB", "error", err)
		return nil, err
	}
	return tx, nil
}

func (e *TokenAdminRegistryEmitterImpl) EmitPoolRegisteredEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool) (*models.SuiTransactionBlockResponse, error) {
	e.lggr.Infow("Emitting pool registered event")
	emitPoolRegisteredEvent, err := e.boundContract.EncodeCallArgs(
		"emit_pool_registered_event",
		[]string{},
		[]string{},
		nil,
	)
	if err != nil {
		e.lggr.Errorw("Failed to encode call", "error", err)
		return nil, err
	}

	return e.executePTBCall(ctx, ptb, addToPTB, emitPoolRegisteredEvent, "emit_pool_registered_event")
}

func (e *TokenAdminRegistryEmitterImpl) EmitAdministratorTransferredEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool) (*models.SuiTransactionBlockResponse, error) {
	e.lggr.Infow("Emitting administrator transferred event")
	emitAdministratorTransferredEvent, err := e.boundContract.EncodeCallArgs(
		"emit_administrator_transferred_event",
		[]string{},
		[]string{},
		nil,
	)
	if err != nil {
		e.lggr.Errorw("Failed to encode call", "error", err)
		return nil, err
	}

	return e.executePTBCall(ctx, ptb, addToPTB, emitAdministratorTransferredEvent, "emit_administrator_transferred_event")
}

func (e *TokenAdminRegistryEmitterImpl) EmitPoolSetEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool) (*models.SuiTransactionBlockResponse, error) {
	e.lggr.Infow("Emitting pool set event")
	emitPoolSetEvent, err := e.boundContract.EncodeCallArgs(
		"emit_pool_set_event",
		[]string{},
		[]string{},
		nil,
	)
	if err != nil {
		e.lggr.Errorw("Failed to encode call", "error", err)
		return nil, err
	}

	return e.executePTBCall(ctx, ptb, addToPTB, emitPoolSetEvent, "emit_pool_set_event")
}

func (e *TokenAdminRegistryEmitterImpl) EmitPoolUnregisteredEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool) (*models.SuiTransactionBlockResponse, error) {
	e.lggr.Infow("Emitting pool unregistered event")
	emitPoolUnregisteredEvent, err := e.boundContract.EncodeCallArgs(
		"emit_pool_unregistered_event",
		[]string{},
		[]string{},
		nil,
	)
	if err != nil {
		e.lggr.Errorw("Failed to encode call", "error", err)
		return nil, err
	}

	return e.executePTBCall(ctx, ptb, addToPTB, emitPoolUnregisteredEvent, "emit_pool_unregistered_event")
}

func (e *TokenAdminRegistryEmitterImpl) EmitAdministratorTransferRequestedEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool) (*models.SuiTransactionBlockResponse, error) {
	e.lggr.Infow("Emitting administrator transfer requested event")
	emitAdministratorTransferRequestedEvent, err := e.boundContract.EncodeCallArgs(
		"emit_administrator_transfer_requested_event",
		[]string{},
		[]string{},
		nil,
	)
	if err != nil {
		e.lggr.Errorw("Failed to encode call", "error", err)
		return nil, err
	}

	return e.executePTBCall(ctx, ptb, addToPTB, emitAdministratorTransferRequestedEvent, "emit_administrator_transfer_requested_event")
}
