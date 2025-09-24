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

type ManagedTokenPoolEmitter interface {
	EmitEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool, eventName string) (*models.SuiTransactionBlockResponse, error)
	BatchEmitEvents(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction) (*models.SuiTransactionBlockResponse, error)
	EmitTokenLockedOrBurnedEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool, amount uint64, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	EmitChainAddedEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	EmitChainRemovedEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	EmitRemotePoolAddedEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	EmitRemotePoolRemovedEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	EmitTokenReleasedOrMintedEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool, amount uint64, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
}

type ManagedTokenPoolEmitterImpl struct {
	lggr           logger.Logger
	packageId      string
	boundContract  *bind.BoundContract
	callOpts       bind.CallOpts
	client         sui.ISuiAPI
	accountAddress string
}

func NewManagedTokenPoolEmitter(
	lggr logger.Logger,
	packageId string,
	signer rel.SuiSigner,
	callOpts bind.CallOpts,
	client sui.ISuiAPI,
	accountAddress string,
) *ManagedTokenPoolEmitterImpl {
	boundContract, err := bind.NewBoundContract(
		packageId,
		"test",
		"managed_token_pool",
		client,
	)
	if err != nil {
		lggr.Errorw("Failed to create managed_token_pool bound contract", "error", err)
		return nil
	}

	return &ManagedTokenPoolEmitterImpl{
		lggr:           lggr,
		packageId:      packageId,
		boundContract:  boundContract,
		callOpts:       callOpts,
		client:         client,
		accountAddress: accountAddress,
	}
}

// executePTBCall handles the common PTB execution logic for event emission
func (e *ManagedTokenPoolEmitterImpl) executePTBCall(ctx context.Context, ptb *transaction.Transaction, addToPTB bool, encodedCall *bind.EncodedCall, operationName string) (*models.SuiTransactionBlockResponse, error) {
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

func (e *ManagedTokenPoolEmitterImpl) EmitEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool, eventName string) (*models.SuiTransactionBlockResponse, error) {
	e.lggr.Infow("Emitting event", "eventName", eventName)

	if ptb == nil {
		ptb = transaction.NewTransaction()
		addToPTB = false
	}

	switch eventName {
	case "emit_token_locked_or_burned_event":
		return e.EmitTokenLockedOrBurnedEvent(ctx, gasBudget, ptb, addToPTB, 1000, 1)
	case "emit_chain_added_event":
		return e.EmitChainAddedEvent(ctx, gasBudget, ptb, addToPTB, 1)
	case "emit_chain_removed_event":
		return e.EmitChainRemovedEvent(ctx, gasBudget, ptb, addToPTB, 1)
	case "emit_remote_pool_added_event":
		return e.EmitRemotePoolAddedEvent(ctx, gasBudget, ptb, addToPTB, 1)
	case "emit_remote_pool_removed_event":
		return e.EmitRemotePoolRemovedEvent(ctx, gasBudget, ptb, addToPTB, 1)
	case "emit_token_released_or_minted_event":
		return e.EmitTokenReleasedOrMintedEvent(ctx, gasBudget, ptb, addToPTB, 1000, 1)
	default:
		e.lggr.Errorw("Invalid event name", "eventName", eventName)
		return nil, fmt.Errorf("invalid event name: %s", eventName)
	}
}

func (e *ManagedTokenPoolEmitterImpl) BatchEmitEvents(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction) (*models.SuiTransactionBlockResponse, error) {
	if ptb == nil {
		ptb = transaction.NewTransaction()
	}

	e.lggr.Infow("Emitting batch of events for managed token pool")

	e.EmitTokenLockedOrBurnedEvent(ctx, gasBudget, ptb, true, 1000, 1)
	e.EmitChainAddedEvent(ctx, gasBudget, ptb, true, 1)
	e.EmitChainRemovedEvent(ctx, gasBudget, ptb, true, 1)
	e.EmitRemotePoolAddedEvent(ctx, gasBudget, ptb, true, 1)
	e.EmitRemotePoolRemovedEvent(ctx, gasBudget, ptb, true, 1)
	e.EmitTokenReleasedOrMintedEvent(ctx, gasBudget, ptb, true, 1000, 1)

	tx, err := bind.ExecutePTB(ctx, &e.callOpts, e.client, ptb)
	if err != nil {
		e.lggr.Errorw("Failed to execute PTB", "error", err)
		return nil, err
	}
	return tx, nil
}

func (e *ManagedTokenPoolEmitterImpl) EmitTokenLockedOrBurnedEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool, amount uint64, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	e.lggr.Infow("Emitting token locked or burned event", "amount", amount, "remoteChainSelector", remoteChainSelector)
	emitTokenLockedOrBurnedEvent, err := e.boundContract.EncodeCallArgs(
		"emit_token_locked_or_burned_event",
		[]string{},
		[]string{"u64", "u64"},
		[]any{amount, remoteChainSelector},
	)
	if err != nil {
		e.lggr.Errorw("Failed to encode call", "error", err)
		return nil, err
	}

	return e.executePTBCall(ctx, ptb, addToPTB, emitTokenLockedOrBurnedEvent, "emit_token_locked_or_burned_event")
}

func (e *ManagedTokenPoolEmitterImpl) EmitChainAddedEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	e.lggr.Infow("Emitting chain added event", "remoteChainSelector", remoteChainSelector)
	emitChainAddedEvent, err := e.boundContract.EncodeCallArgs(
		"emit_chain_added_event",
		[]string{},
		[]string{"u64"},
		[]any{remoteChainSelector},
	)
	if err != nil {
		e.lggr.Errorw("Failed to encode call", "error", err)
		return nil, err
	}

	return e.executePTBCall(ctx, ptb, addToPTB, emitChainAddedEvent, "emit_chain_added_event")
}

func (e *ManagedTokenPoolEmitterImpl) EmitChainRemovedEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	e.lggr.Infow("Emitting chain removed event", "remoteChainSelector", remoteChainSelector)
	emitChainRemovedEvent, err := e.boundContract.EncodeCallArgs(
		"emit_chain_removed_event",
		[]string{},
		[]string{"u64"},
		[]any{remoteChainSelector},
	)
	if err != nil {
		e.lggr.Errorw("Failed to encode call", "error", err)
		return nil, err
	}

	return e.executePTBCall(ctx, ptb, addToPTB, emitChainRemovedEvent, "emit_chain_removed_event")
}

func (e *ManagedTokenPoolEmitterImpl) EmitRemotePoolAddedEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	e.lggr.Infow("Emitting remote pool added event", "remoteChainSelector", remoteChainSelector)
	emitRemotePoolAddedEvent, err := e.boundContract.EncodeCallArgs(
		"emit_remote_pool_added_event",
		[]string{},
		[]string{"u64"},
		[]any{remoteChainSelector},
	)
	if err != nil {
		e.lggr.Errorw("Failed to encode call", "error", err)
		return nil, err
	}

	return e.executePTBCall(ctx, ptb, addToPTB, emitRemotePoolAddedEvent, "emit_remote_pool_added_event")
}

func (e *ManagedTokenPoolEmitterImpl) EmitRemotePoolRemovedEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	e.lggr.Infow("Emitting remote pool removed event", "remoteChainSelector", remoteChainSelector)
	emitRemotePoolRemovedEvent, err := e.boundContract.EncodeCallArgs(
		"emit_remote_pool_removed_event",
		[]string{},
		[]string{"u64"},
		[]any{remoteChainSelector},
	)
	if err != nil {
		e.lggr.Errorw("Failed to encode call", "error", err)
		return nil, err
	}

	return e.executePTBCall(ctx, ptb, addToPTB, emitRemotePoolRemovedEvent, "emit_remote_pool_removed_event")
}

func (e *ManagedTokenPoolEmitterImpl) EmitTokenReleasedOrMintedEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool, amount uint64, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	e.lggr.Infow("Emitting token released or minted event", "amount", amount, "remoteChainSelector", remoteChainSelector)
	emitTokenReleasedOrMintedEvent, err := e.boundContract.EncodeCallArgs(
		"emit_token_released_or_minted_event",
		[]string{},
		[]string{"u64", "u64"},
		[]any{amount, remoteChainSelector},
	)
	if err != nil {
		e.lggr.Errorw("Failed to encode call", "error", err)
		return nil, err
	}

	return e.executePTBCall(ctx, ptb, addToPTB, emitTokenReleasedOrMintedEvent, "emit_token_released_or_minted_event")
}
