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

type BurnMintTokenPoolEmitter interface {
	EmitEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool, eventName string) (*models.SuiTransactionBlockResponse, error)
	BatchEmitEvents(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction) (*models.SuiTransactionBlockResponse, error)
	EmitLockedOrBurnedEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool, amount uint64, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	EmitReleasedOrMintedEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool, amount uint64, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	EmitRemotePoolAddedEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	EmitRemotePoolRemovedEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	EmitChainAddedEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	EmitChainRemovedEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	EmitLiquidityAddedEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool, amount uint64) (*models.SuiTransactionBlockResponse, error)
	EmitLiquidityRemovedEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool, amount uint64) (*models.SuiTransactionBlockResponse, error)
	EmitRebalancerSetEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool) (*models.SuiTransactionBlockResponse, error)
}

type BurnMintTokenPoolEmitterImpl struct {
	lggr           logger.Logger
	packageId      string
	boundContract  *bind.BoundContract
	callOpts       bind.CallOpts
	client         sui.ISuiAPI
	accountAddress string
}

func NewBurnMintTokenPoolEmitter(
	lggr logger.Logger,
	packageId string,
	signer rel.SuiSigner,
	callOpts bind.CallOpts,
	client sui.ISuiAPI,
	accountAddress string,
) *BurnMintTokenPoolEmitterImpl {
	boundContract, err := bind.NewBoundContract(
		packageId,
		"test",
		"burn_mint_token_pool",
		client,
	)
	if err != nil {
		lggr.Errorw("Failed to create burn_mint_token_pool bound contract", "error", err)
		return nil
	}

	return &BurnMintTokenPoolEmitterImpl{
		lggr:           lggr,
		packageId:      packageId,
		boundContract:  boundContract,
		callOpts:       callOpts,
		client:         client,
		accountAddress: accountAddress,
	}
}

// executePTBCall handles the common PTB execution logic for event emission
func (e *BurnMintTokenPoolEmitterImpl) executePTBCall(ctx context.Context, ptb *transaction.Transaction, addToPTB bool, encodedCall *bind.EncodedCall, operationName string) (*models.SuiTransactionBlockResponse, error) {
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

func (e *BurnMintTokenPoolEmitterImpl) EmitEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool, eventName string) (*models.SuiTransactionBlockResponse, error) {
	e.lggr.Infow("Emitting event", "eventName", eventName)

	if ptb == nil {
		ptb = transaction.NewTransaction()
		addToPTB = false
	}

	switch eventName {
	case "emit_locked_or_burned_event":
		return e.EmitLockedOrBurnedEvent(ctx, gasBudget, ptb, addToPTB, 1000, 1)
	case "emit_released_or_minted_event":
		return e.EmitReleasedOrMintedEvent(ctx, gasBudget, ptb, addToPTB, 1000, 1)
	case "emit_remote_pool_added_event":
		return e.EmitRemotePoolAddedEvent(ctx, gasBudget, ptb, addToPTB, 1)
	case "emit_remote_pool_removed_event":
		return e.EmitRemotePoolRemovedEvent(ctx, gasBudget, ptb, addToPTB, 1)
	case "emit_chain_added_event":
		return e.EmitChainAddedEvent(ctx, gasBudget, ptb, addToPTB, 1)
	case "emit_chain_removed_event":
		return e.EmitChainRemovedEvent(ctx, gasBudget, ptb, addToPTB, 1)
	case "emit_liquidity_added_event":
		return e.EmitLiquidityAddedEvent(ctx, gasBudget, ptb, addToPTB, 1000)
	case "emit_liquidity_removed_event":
		return e.EmitLiquidityRemovedEvent(ctx, gasBudget, ptb, addToPTB, 1000)
	case "emit_rebalancer_set_event":
		return e.EmitRebalancerSetEvent(ctx, gasBudget, ptb, addToPTB)
	default:
		e.lggr.Errorw("Invalid event name", "eventName", eventName)
		return nil, fmt.Errorf("invalid event name: %s", eventName)
	}
}

func (e *BurnMintTokenPoolEmitterImpl) BatchEmitEvents(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction) (*models.SuiTransactionBlockResponse, error) {
	if ptb == nil {
		ptb = transaction.NewTransaction()
	}

	e.lggr.Infow("Emitting batch of events for burn mint token pool")

	e.EmitLockedOrBurnedEvent(ctx, gasBudget, ptb, true, 1000, 1)
	e.EmitReleasedOrMintedEvent(ctx, gasBudget, ptb, true, 1000, 1)
	e.EmitRemotePoolAddedEvent(ctx, gasBudget, ptb, true, 1)
	e.EmitRemotePoolRemovedEvent(ctx, gasBudget, ptb, true, 1)
	e.EmitChainAddedEvent(ctx, gasBudget, ptb, true, 1)
	e.EmitChainRemovedEvent(ctx, gasBudget, ptb, true, 1)
	e.EmitLiquidityAddedEvent(ctx, gasBudget, ptb, true, 1000)
	e.EmitLiquidityRemovedEvent(ctx, gasBudget, ptb, true, 1000)
	e.EmitRebalancerSetEvent(ctx, gasBudget, ptb, true)

	tx, err := bind.ExecutePTB(ctx, &e.callOpts, e.client, ptb)
	if err != nil {
		e.lggr.Errorw("Failed to execute PTB", "error", err)
		return nil, err
	}
	return tx, nil
}

func (e *BurnMintTokenPoolEmitterImpl) EmitLockedOrBurnedEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool, amount uint64, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	e.lggr.Infow("Emitting locked or burned event", "amount", amount, "remoteChainSelector", remoteChainSelector)
	emitLockedOrBurnedEvent, err := e.boundContract.EncodeCallArgs(
		"emit_locked_or_burned_event",
		[]string{},
		[]string{"u64", "u64"},
		[]any{amount, remoteChainSelector},
	)
	if err != nil {
		e.lggr.Errorw("Failed to encode call", "error", err)
		return nil, err
	}

	return e.executePTBCall(ctx, ptb, addToPTB, emitLockedOrBurnedEvent, "emit_locked_or_burned_event")
}

func (e *BurnMintTokenPoolEmitterImpl) EmitReleasedOrMintedEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool, amount uint64, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	e.lggr.Infow("Emitting released or minted event", "amount", amount, "remoteChainSelector", remoteChainSelector)
	emitReleasedOrMintedEvent, err := e.boundContract.EncodeCallArgs(
		"emit_released_or_minted_event",
		[]string{},
		[]string{"u64", "u64"},
		[]any{amount, remoteChainSelector},
	)
	if err != nil {
		e.lggr.Errorw("Failed to encode call", "error", err)
		return nil, err
	}

	return e.executePTBCall(ctx, ptb, addToPTB, emitReleasedOrMintedEvent, "emit_released_or_minted_event")
}

func (e *BurnMintTokenPoolEmitterImpl) EmitRemotePoolAddedEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
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

func (e *BurnMintTokenPoolEmitterImpl) EmitRemotePoolRemovedEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
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

func (e *BurnMintTokenPoolEmitterImpl) EmitChainAddedEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
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

func (e *BurnMintTokenPoolEmitterImpl) EmitChainRemovedEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
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

func (e *BurnMintTokenPoolEmitterImpl) EmitLiquidityAddedEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool, amount uint64) (*models.SuiTransactionBlockResponse, error) {
	e.lggr.Infow("Emitting liquidity added event", "amount", amount)
	emitLiquidityAddedEvent, err := e.boundContract.EncodeCallArgs(
		"emit_liquidity_added_event",
		[]string{},
		[]string{"u64"},
		[]any{amount},
	)
	if err != nil {
		e.lggr.Errorw("Failed to encode call", "error", err)
		return nil, err
	}

	return e.executePTBCall(ctx, ptb, addToPTB, emitLiquidityAddedEvent, "emit_liquidity_added_event")
}

func (e *BurnMintTokenPoolEmitterImpl) EmitLiquidityRemovedEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool, amount uint64) (*models.SuiTransactionBlockResponse, error) {
	e.lggr.Infow("Emitting liquidity removed event", "amount", amount)
	emitLiquidityRemovedEvent, err := e.boundContract.EncodeCallArgs(
		"emit_liquidity_removed_event",
		[]string{},
		[]string{"u64"},
		[]any{amount},
	)
	if err != nil {
		e.lggr.Errorw("Failed to encode call", "error", err)
		return nil, err
	}

	return e.executePTBCall(ctx, ptb, addToPTB, emitLiquidityRemovedEvent, "emit_liquidity_removed_event")
}

func (e *BurnMintTokenPoolEmitterImpl) EmitRebalancerSetEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool) (*models.SuiTransactionBlockResponse, error) {
	e.lggr.Infow("Emitting rebalancer set event")
	emitRebalancerSetEvent, err := e.boundContract.EncodeCallArgs(
		"emit_rebalancer_set_event",
		[]string{},
		[]string{},
		nil,
	)
	if err != nil {
		e.lggr.Errorw("Failed to encode call", "error", err)
		return nil, err
	}

	return e.executePTBCall(ctx, ptb, addToPTB, emitRebalancerSetEvent, "emit_rebalancer_set_event")
}
