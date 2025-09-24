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

type OnRampEmitter interface {
	EmitEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool, eventName string) (*models.SuiTransactionBlockResponse, error)
	BatchEmitEvents(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction) (*models.SuiTransactionBlockResponse, error)
	EmitConfigSetEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool) (*models.SuiTransactionBlockResponse, error)
	EmitDestChainConfigSetEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool) (*models.SuiTransactionBlockResponse, error)
	EmitCCIPMessageSentEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool, destChainSelector uint64, sequenceNumber uint64) (*models.SuiTransactionBlockResponse, error)
	EmitAllowlistSendersAddedEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool, destChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	EmitAllowlistSendersRemovedEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool, destChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	EmitFeeTokenWithdrawnEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool) (*models.SuiTransactionBlockResponse, error)
}

type OnRampEmitterImpl struct {
	lggr           logger.Logger
	packageId      string
	boundContract  *bind.BoundContract
	callOpts       bind.CallOpts
	client         sui.ISuiAPI
	accountAddress string
}

func NewOnRampEmitter(
	lggr logger.Logger,
	packageId string,
	signer rel.SuiSigner,
	callOpts bind.CallOpts,
	client sui.ISuiAPI,
	accountAddress string,
) *OnRampEmitterImpl {
	boundContract, err := bind.NewBoundContract(
		packageId,
		"test",
		"onramp",
		client,
	)
	if err != nil {
		lggr.Errorw("Failed to create onramp bound contract", "error", err)
		return nil
	}

	return &OnRampEmitterImpl{
		lggr:           lggr,
		packageId:      packageId,
		boundContract:  boundContract,
		callOpts:       callOpts,
		client:         client,
		accountAddress: accountAddress,
	}
}

// executePTBCall handles the common PTB execution logic for event emission
func (e *OnRampEmitterImpl) executePTBCall(ctx context.Context, ptb *transaction.Transaction, addToPTB bool, encodedCall *bind.EncodedCall, operationName string) (*models.SuiTransactionBlockResponse, error) {
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

func (e *OnRampEmitterImpl) EmitEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool, eventName string) (*models.SuiTransactionBlockResponse, error) {
	e.lggr.Infow("Emitting event", "eventName", eventName)

	if ptb == nil {
		ptb = transaction.NewTransaction()
		addToPTB = false
	}

	switch eventName {
	case "emit_config_set_event":
		return e.EmitConfigSetEvent(ctx, gasBudget, ptb, addToPTB)
	case "emit_dest_chain_config_set_event":
		return e.EmitDestChainConfigSetEvent(ctx, gasBudget, ptb, addToPTB)
	case "emit_ccip_message_sent_event":
		return e.EmitCCIPMessageSentEvent(ctx, gasBudget, ptb, addToPTB, 1, 1)
	case "emit_allowlist_senders_added_event":
		return e.EmitAllowlistSendersAddedEvent(ctx, gasBudget, ptb, addToPTB, 1)
	case "emit_allowlist_senders_removed_event":
		return e.EmitAllowlistSendersRemovedEvent(ctx, gasBudget, ptb, addToPTB, 1)
	case "emit_fee_token_withdrawn_event":
		return e.EmitFeeTokenWithdrawnEvent(ctx, gasBudget, ptb, addToPTB)
	default:
		e.lggr.Errorw("Invalid event name", "eventName", eventName)
		return nil, fmt.Errorf("invalid event name: %s", eventName)
	}
}

func (e *OnRampEmitterImpl) BatchEmitEvents(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction) (*models.SuiTransactionBlockResponse, error) {
	if ptb == nil {
		ptb = transaction.NewTransaction()
	}

	e.lggr.Infow("Emitting batch of events for onramp")

	e.EmitConfigSetEvent(ctx, gasBudget, ptb, true)
	e.EmitDestChainConfigSetEvent(ctx, gasBudget, ptb, true)
	e.EmitCCIPMessageSentEvent(ctx, gasBudget, ptb, true, 1, 1)
	e.EmitAllowlistSendersAddedEvent(ctx, gasBudget, ptb, true, 1)
	e.EmitAllowlistSendersRemovedEvent(ctx, gasBudget, ptb, true, 1)
	e.EmitFeeTokenWithdrawnEvent(ctx, gasBudget, ptb, true)

	tx, err := bind.ExecutePTB(ctx, &e.callOpts, e.client, ptb)
	if err != nil {
		e.lggr.Errorw("Failed to execute PTB", "error", err)
		return nil, err
	}
	return tx, nil
}

func (e *OnRampEmitterImpl) EmitConfigSetEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool) (*models.SuiTransactionBlockResponse, error) {
	e.lggr.Infow("Emitting config set event")
	emitConfigSetEvent, err := e.boundContract.EncodeCallArgs(
		"emit_config_set_event",
		[]string{},
		[]string{},
		nil,
	)
	if err != nil {
		e.lggr.Errorw("Failed to encode call", "error", err)
		return nil, err
	}

	return e.executePTBCall(ctx, ptb, addToPTB, emitConfigSetEvent, "emit_config_set_event")
}

func (e *OnRampEmitterImpl) EmitDestChainConfigSetEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool) (*models.SuiTransactionBlockResponse, error) {
	e.lggr.Infow("Emitting dest chain config set event")
	emitDestChainConfigSetEvent, err := e.boundContract.EncodeCallArgs(
		"emit_dest_chain_config_set_event",
		[]string{},
		[]string{},
		nil,
	)
	if err != nil {
		e.lggr.Errorw("Failed to encode call", "error", err)
		return nil, err
	}

	return e.executePTBCall(ctx, ptb, addToPTB, emitDestChainConfigSetEvent, "emit_dest_chain_config_set_event")
}

func (e *OnRampEmitterImpl) EmitCCIPMessageSentEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool, destChainSelector uint64, sequenceNumber uint64) (*models.SuiTransactionBlockResponse, error) {
	e.lggr.Infow("Emitting CCIP message sent event", "destChainSelector", destChainSelector, "sequenceNumber", sequenceNumber)
	emitCCIPMessageSentEvent, err := e.boundContract.EncodeCallArgs(
		"emit_ccip_message_sent_event",
		[]string{},
		[]string{"u64", "u64"},
		[]any{destChainSelector, sequenceNumber},
	)
	if err != nil {
		e.lggr.Errorw("Failed to encode call", "error", err)
		return nil, err
	}

	return e.executePTBCall(ctx, ptb, addToPTB, emitCCIPMessageSentEvent, "emit_ccip_message_sent_event")
}

func (e *OnRampEmitterImpl) EmitAllowlistSendersAddedEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool, destChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	e.lggr.Infow("Emitting allowlist senders added event", "destChainSelector", destChainSelector)
	emitAllowlistSendersAddedEvent, err := e.boundContract.EncodeCallArgs(
		"emit_allowlist_senders_added_event",
		[]string{},
		[]string{"u64"},
		[]any{destChainSelector},
	)
	if err != nil {
		e.lggr.Errorw("Failed to encode call", "error", err)
		return nil, err
	}

	return e.executePTBCall(ctx, ptb, addToPTB, emitAllowlistSendersAddedEvent, "emit_allowlist_senders_added_event")
}

func (e *OnRampEmitterImpl) EmitAllowlistSendersRemovedEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool, destChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	e.lggr.Infow("Emitting allowlist senders removed event", "destChainSelector", destChainSelector)
	emitAllowlistSendersRemovedEvent, err := e.boundContract.EncodeCallArgs(
		"emit_allowlist_senders_removed_event",
		[]string{},
		[]string{"u64"},
		[]any{destChainSelector},
	)
	if err != nil {
		e.lggr.Errorw("Failed to encode call", "error", err)
		return nil, err
	}

	return e.executePTBCall(ctx, ptb, addToPTB, emitAllowlistSendersRemovedEvent, "emit_allowlist_senders_removed_event")
}

func (e *OnRampEmitterImpl) EmitFeeTokenWithdrawnEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool) (*models.SuiTransactionBlockResponse, error) {
	e.lggr.Infow("Emitting fee token withdrawn event")
	emitFeeTokenWithdrawnEvent, err := e.boundContract.EncodeCallArgs(
		"emit_fee_token_withdrawn_event",
		[]string{},
		[]string{},
		nil,
	)
	if err != nil {
		e.lggr.Errorw("Failed to encode call", "error", err)
		return nil, err
	}

	return e.executePTBCall(ctx, ptb, addToPTB, emitFeeTokenWithdrawnEvent, "emit_fee_token_withdrawn_event")
}
