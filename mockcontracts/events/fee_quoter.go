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

type FeeQuoterEmitter interface {
	EmitEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool, eventName string) (*models.SuiTransactionBlockResponse, error)
	BatchEmitEvents(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction) (*models.SuiTransactionBlockResponse, error)
	EmitFeeTokenAddedEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool) (*models.SuiTransactionBlockResponse, error)
	EmitFeeTokenRemovedEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool) (*models.SuiTransactionBlockResponse, error)
	EmitTokenTransferFeeConfigAddedEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool, destChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	EmitTokenTransferFeeConfigRemovedEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool, destChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	EmitUsdPerTokenUpdatedEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool) (*models.SuiTransactionBlockResponse, error)
	EmitUsdPerUnitGasUpdatedEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool) (*models.SuiTransactionBlockResponse, error)
	EmitDestChainAddedEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool, destChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	EmitDestChainConfigUpdatedEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool, destChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	EmitPremiumMultiplierWeiPerEthUpdatedEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool) (*models.SuiTransactionBlockResponse, error)
}

type FeeQuoterEmitterImpl struct {
	lggr           logger.Logger
	packageId      string
	boundContract  *bind.BoundContract
	callOpts       bind.CallOpts
	client         sui.ISuiAPI
	accountAddress string
	clockObjectId  string
}

func NewFeeQuoterEmitter(
	lggr logger.Logger,
	packageId string,
	signer rel.SuiSigner,
	callOpts bind.CallOpts,
	client sui.ISuiAPI,
	accountAddress string,
) *FeeQuoterEmitterImpl {
	boundContract, err := bind.NewBoundContract(
		packageId,
		"test",
		"fee_quoter",
		client,
	)
	if err != nil {
		lggr.Errorw("Failed to create fee_quoter bound contract", "error", err)
		return nil
	}

	clockObjectID := "0x6"

	return &FeeQuoterEmitterImpl{
		lggr:           lggr,
		packageId:      packageId,
		boundContract:  boundContract,
		callOpts:       callOpts,
		client:         client,
		accountAddress: accountAddress,
		clockObjectId:  clockObjectID,
	}
}

// executePTBCall handles the common PTB execution logic for event emission
func (e *FeeQuoterEmitterImpl) executePTBCall(ctx context.Context, ptb *transaction.Transaction, addToPTB bool, encodedCall *bind.EncodedCall, operationName string) (*models.SuiTransactionBlockResponse, error) {
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

func (e *FeeQuoterEmitterImpl) EmitEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool, eventName string) (*models.SuiTransactionBlockResponse, error) {
	e.lggr.Infow("Emitting event", "eventName", eventName)

	if ptb == nil {
		ptb = transaction.NewTransaction()
		addToPTB = false
	}

	switch eventName {
	case "emit_fee_token_added_event":
		return e.EmitFeeTokenAddedEvent(ctx, gasBudget, ptb, addToPTB)
	case "emit_fee_token_removed_event":
		return e.EmitFeeTokenRemovedEvent(ctx, gasBudget, ptb, addToPTB)
	case "emit_token_transfer_fee_config_added_event":
		return e.EmitTokenTransferFeeConfigAddedEvent(ctx, gasBudget, ptb, addToPTB, 1)
	case "emit_token_transfer_fee_config_removed_event":
		return e.EmitTokenTransferFeeConfigRemovedEvent(ctx, gasBudget, ptb, addToPTB, 1)
	case "emit_usd_per_token_updated_event":
		return e.EmitUsdPerTokenUpdatedEvent(ctx, gasBudget, ptb, addToPTB)
	case "emit_usd_per_unit_gas_updated_event":
		return e.EmitUsdPerUnitGasUpdatedEvent(ctx, gasBudget, ptb, addToPTB)
	case "emit_dest_chain_added_event":
		return e.EmitDestChainAddedEvent(ctx, gasBudget, ptb, addToPTB, 1)
	case "emit_dest_chain_config_updated_event":
		return e.EmitDestChainConfigUpdatedEvent(ctx, gasBudget, ptb, addToPTB, 1)
	case "emit_premium_multiplier_wei_per_eth_updated_event":
		return e.EmitPremiumMultiplierWeiPerEthUpdatedEvent(ctx, gasBudget, ptb, addToPTB)
	default:
		e.lggr.Errorw("Invalid event name", "eventName", eventName)
		return nil, fmt.Errorf("invalid event name: %s", eventName)
	}
}

func (e *FeeQuoterEmitterImpl) BatchEmitEvents(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction) (*models.SuiTransactionBlockResponse, error) {
	if ptb == nil {
		ptb = transaction.NewTransaction()
	}

	e.lggr.Infow("Emitting batch of events for fee quoter")

	e.EmitFeeTokenAddedEvent(ctx, gasBudget, ptb, true)
	e.EmitFeeTokenRemovedEvent(ctx, gasBudget, ptb, true)
	e.EmitTokenTransferFeeConfigAddedEvent(ctx, gasBudget, ptb, true, 1)
	e.EmitTokenTransferFeeConfigRemovedEvent(ctx, gasBudget, ptb, true, 1)
	e.EmitUsdPerTokenUpdatedEvent(ctx, gasBudget, ptb, true)
	e.EmitUsdPerUnitGasUpdatedEvent(ctx, gasBudget, ptb, true)
	e.EmitDestChainAddedEvent(ctx, gasBudget, ptb, true, 1)
	e.EmitDestChainConfigUpdatedEvent(ctx, gasBudget, ptb, true, 1)
	e.EmitPremiumMultiplierWeiPerEthUpdatedEvent(ctx, gasBudget, ptb, true)

	tx, err := bind.ExecutePTB(ctx, &e.callOpts, e.client, ptb)
	if err != nil {
		e.lggr.Errorw("Failed to execute PTB", "error", err)
		return nil, err
	}
	return tx, nil
}

func (e *FeeQuoterEmitterImpl) EmitFeeTokenAddedEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool) (*models.SuiTransactionBlockResponse, error) {
	e.lggr.Infow("Emitting fee token added event")
	emitFeeTokenAddedEvent, err := e.boundContract.EncodeCallArgs(
		"emit_fee_token_added_event",
		[]string{},
		[]string{},
		nil,
	)
	if err != nil {
		e.lggr.Errorw("Failed to encode call", "error", err)
		return nil, err
	}

	return e.executePTBCall(ctx, ptb, addToPTB, emitFeeTokenAddedEvent, "emit_fee_token_added_event")
}

func (e *FeeQuoterEmitterImpl) EmitFeeTokenRemovedEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool) (*models.SuiTransactionBlockResponse, error) {
	e.lggr.Infow("Emitting fee token removed event")
	emitFeeTokenRemovedEvent, err := e.boundContract.EncodeCallArgs(
		"emit_fee_token_removed_event",
		[]string{},
		[]string{},
		nil,
	)
	if err != nil {
		e.lggr.Errorw("Failed to encode call", "error", err)
		return nil, err
	}

	return e.executePTBCall(ctx, ptb, addToPTB, emitFeeTokenRemovedEvent, "emit_fee_token_removed_event")
}

func (e *FeeQuoterEmitterImpl) EmitTokenTransferFeeConfigAddedEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool, destChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	e.lggr.Infow("Emitting token transfer fee config added event", "destChainSelector", destChainSelector)
	emitTokenTransferFeeConfigAddedEvent, err := e.boundContract.EncodeCallArgs(
		"emit_token_transfer_fee_config_added_event",
		[]string{},
		[]string{"u64"},
		[]any{destChainSelector},
	)
	if err != nil {
		e.lggr.Errorw("Failed to encode call", "error", err)
		return nil, err
	}

	return e.executePTBCall(ctx, ptb, addToPTB, emitTokenTransferFeeConfigAddedEvent, "emit_token_transfer_fee_config_added_event")
}

func (e *FeeQuoterEmitterImpl) EmitTokenTransferFeeConfigRemovedEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool, destChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	e.lggr.Infow("Emitting token transfer fee config removed event", "destChainSelector", destChainSelector)
	emitTokenTransferFeeConfigRemovedEvent, err := e.boundContract.EncodeCallArgs(
		"emit_token_transfer_fee_config_removed_event",
		[]string{},
		[]string{"u64"},
		[]any{destChainSelector},
	)
	if err != nil {
		e.lggr.Errorw("Failed to encode call", "error", err)
		return nil, err
	}

	return e.executePTBCall(ctx, ptb, addToPTB, emitTokenTransferFeeConfigRemovedEvent, "emit_token_transfer_fee_config_removed_event")
}

func (e *FeeQuoterEmitterImpl) EmitUsdPerTokenUpdatedEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool) (*models.SuiTransactionBlockResponse, error) {
	e.lggr.Infow("Emitting usd per token updated event")
	emitUsdPerTokenUpdatedEvent, err := e.boundContract.EncodeCallArgs(
		"emit_usd_per_token_updated_event",
		[]string{},
		[]string{"&clock::Clock"},
		[]any{e.clockObjectId},
	)
	if err != nil {
		e.lggr.Errorw("Failed to encode call", "error", err)
		return nil, err
	}

	return e.executePTBCall(ctx, ptb, addToPTB, emitUsdPerTokenUpdatedEvent, "emit_usd_per_token_updated_event")
}

func (e *FeeQuoterEmitterImpl) EmitUsdPerUnitGasUpdatedEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool) (*models.SuiTransactionBlockResponse, error) {
	e.lggr.Infow("Emitting usd per unit gas updated event")
	emitUsdPerUnitGasUpdatedEvent, err := e.boundContract.EncodeCallArgs(
		"emit_usd_per_unit_gas_updated_event",
		[]string{},
		[]string{"&clock::Clock"},
		[]any{e.clockObjectId},
	)
	if err != nil {
		e.lggr.Errorw("Failed to encode call", "error", err)
		return nil, err
	}

	return e.executePTBCall(ctx, ptb, addToPTB, emitUsdPerUnitGasUpdatedEvent, "emit_usd_per_unit_gas_updated_event")
}

func (e *FeeQuoterEmitterImpl) EmitDestChainAddedEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool, destChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	e.lggr.Infow("Emitting dest chain added event", "destChainSelector", destChainSelector)
	emitDestChainAddedEvent, err := e.boundContract.EncodeCallArgs(
		"emit_dest_chain_added_event",
		[]string{},
		[]string{"u64"},
		[]any{destChainSelector},
	)
	if err != nil {
		e.lggr.Errorw("Failed to encode call", "error", err)
		return nil, err
	}

	return e.executePTBCall(ctx, ptb, addToPTB, emitDestChainAddedEvent, "emit_dest_chain_added_event")
}

func (e *FeeQuoterEmitterImpl) EmitDestChainConfigUpdatedEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool, destChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	e.lggr.Infow("Emitting dest chain config updated event", "destChainSelector", destChainSelector)
	emitDestChainConfigUpdatedEvent, err := e.boundContract.EncodeCallArgs(
		"emit_dest_chain_config_updated_event",
		[]string{},
		[]string{"u64"},
		[]any{destChainSelector},
	)
	if err != nil {
		e.lggr.Errorw("Failed to encode call", "error", err)
		return nil, err
	}

	return e.executePTBCall(ctx, ptb, addToPTB, emitDestChainConfigUpdatedEvent, "emit_dest_chain_config_updated_event")
}

func (e *FeeQuoterEmitterImpl) EmitPremiumMultiplierWeiPerEthUpdatedEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool) (*models.SuiTransactionBlockResponse, error) {
	e.lggr.Infow("Emitting premium multiplier wei per eth updated event")
	emitPremiumMultiplierWeiPerEthUpdatedEvent, err := e.boundContract.EncodeCallArgs(
		"emit_premium_multiplier_wei_per_eth_updated_event",
		[]string{},
		[]string{},
		nil,
	)
	if err != nil {
		e.lggr.Errorw("Failed to encode call", "error", err)
		return nil, err
	}

	return e.executePTBCall(ctx, ptb, addToPTB, emitPremiumMultiplierWeiPerEthUpdatedEvent, "emit_premium_multiplier_wei_per_eth_updated_event")
}
