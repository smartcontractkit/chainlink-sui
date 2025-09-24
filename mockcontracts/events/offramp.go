package events

import (
	"context"
	"errors"
	"fmt"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/block-vision/sui-go-sdk/transaction"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	rel "github.com/smartcontractkit/chainlink-sui/relayer/signer"
)

type OffRampEmitter interface {
	EmitEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool, eventName string) (*models.SuiTransactionBlockResponse, error)
	BatchEmitEvents(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction) error
	EmitStaticConfigSetEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool) (*models.SuiTransactionBlockResponse, error)
	EmitDynamicConfigSetEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool) (*models.SuiTransactionBlockResponse, error)
	EmitSourceChainConfigSetEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool) (*models.SuiTransactionBlockResponse, error)
	EmitSkippedAlreadyExecutedEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool) (*models.SuiTransactionBlockResponse, error)
	EmitExecutionStateChangedEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool) (*models.SuiTransactionBlockResponse, error)
	EmitCommitReportAcceptedEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool) (*models.SuiTransactionBlockResponse, error)
	EmitSkippedReportExecutionEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool) (*models.SuiTransactionBlockResponse, error)
	EmitOcrConfigEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool) (*models.SuiTransactionBlockResponse, error)
}

type OffRampEmitterImpl struct {
	lggr           logger.Logger
	packageId      string
	boundContract  *bind.BoundContract
	callOpts       bind.CallOpts
	client         sui.ISuiAPI
	accountAddress string
}

func NewOffRampEmitter(
	lggr logger.Logger,
	packageId string,
	signer rel.SuiSigner,
	callOpts bind.CallOpts,
	client sui.ISuiAPI,
	accountAddress string,
) *OffRampEmitterImpl {
	boundContract, err := bind.NewBoundContract(
		packageId,
		"test",
		"offramp",
		client,
	)
	if err != nil {
		lggr.Errorw("Failed to create offramp bound contract", "error", err)
		return nil
	}

	return &OffRampEmitterImpl{
		lggr:           lggr,
		packageId:      packageId,
		boundContract:  boundContract,
		callOpts:       callOpts,
		client:         client,
		accountAddress: accountAddress,
	}
}

// executePTBCall handles the common PTB execution logic for event emission
func (e *OffRampEmitterImpl) executePTBCall(ctx context.Context, ptb *transaction.Transaction, addToPTB bool, encodedCall *bind.EncodedCall, operationName string) (*models.SuiTransactionBlockResponse, error) {
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

func (e *OffRampEmitterImpl) EmitEvent(ctx context.Context, gasBudget uint64, functionName string) (*models.SuiTransactionBlockResponse, error) {
	e.lggr.Infow("Emitting event", "functionName", functionName)

	if functionName == "" {
		e.lggr.Errorw("Function name is required")
		return nil, errors.New("function name is required")
	}

	ptb := transaction.NewTransaction()
	addToPTB := false

	switch functionName {
	case "emit_static_config_set_event":
		return e.EmitStaticConfigSetEvent(ctx, gasBudget, ptb, addToPTB)
	case "emit_dynamic_config_set_event":
		return e.EmitDynamicConfigSetEvent(ctx, gasBudget, ptb, addToPTB)
	case "emit_source_chain_config_set_event":
		return e.EmitSourceChainConfigSetEvent(ctx, gasBudget, ptb, addToPTB)
	case "emit_skipped_already_executed_event":
		return e.EmitSkippedAlreadyExecutedEvent(ctx, gasBudget, ptb, addToPTB)
	case "emit_execution_state_changed_event":
		return e.EmitExecutionStateChangedEvent(ctx, gasBudget, ptb, addToPTB)
	case "emit_commit_report_accepted_event":
		return e.EmitCommitReportAcceptedEvent(ctx, gasBudget, ptb, addToPTB)
	case "emit_skipped_report_execution_event":
		return e.EmitSkippedReportExecutionEvent(ctx, gasBudget, ptb, addToPTB)
	case "emit_ocr_config_event":
		return e.EmitOcrConfigEvent(ctx, gasBudget, ptb, addToPTB)
	default:
		e.lggr.Errorw("Invalid event name", "functionName", functionName)
		return nil, fmt.Errorf("invalid function name: %s", functionName)
	}
}

func (e *OffRampEmitterImpl) BatchEmitEvents(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction) (*models.SuiTransactionBlockResponse, error) {
	if ptb == nil {
		ptb = transaction.NewTransaction()
	}

	e.EmitStaticConfigSetEvent(ctx, gasBudget, ptb, true)
	e.EmitDynamicConfigSetEvent(ctx, gasBudget, ptb, true)
	e.EmitSourceChainConfigSetEvent(ctx, gasBudget, ptb, true)
	e.EmitCommitReportAcceptedEvent(ctx, gasBudget, ptb, true)
	e.EmitSkippedAlreadyExecutedEvent(ctx, gasBudget, ptb, true)
	// e.EmitExecutionStateChangedEvent(ctx, gasBudget, ptb, true)
	e.EmitSkippedReportExecutionEvent(ctx, gasBudget, ptb, true)
	e.EmitOcrConfigEvent(ctx, gasBudget, ptb, true)
	tx, err := bind.ExecutePTB(ctx, &e.callOpts, e.client, ptb)
	if err != nil {
		e.lggr.Errorw("Failed to execute PTB", "error", err)
		return nil, err
	}
	return tx, nil
}

func (e *OffRampEmitterImpl) EmitStaticConfigSetEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool) (*models.SuiTransactionBlockResponse, error) {
	// emit_static_config_set_event
	e.lggr.Infow("Emitting static config set event")
	emitStaticConfigSetEvent, err := e.boundContract.EncodeCallArgsWithGenerics(
		"emit_static_config_set_event",
		[]string{},
		[]string{},
		[]string{"u64"},
		[]any{1},
		nil,
	)
	if err != nil {
		e.lggr.Errorw("Failed to encode call", "error", err)
		return nil, err
	}

	return e.executePTBCall(ctx, ptb, addToPTB, emitStaticConfigSetEvent, "emit_static_config_set_event")
}

func (e *OffRampEmitterImpl) EmitDynamicConfigSetEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool) (*models.SuiTransactionBlockResponse, error) {
	e.lggr.Infow("Emitting dynamic config set event")
	emitDynamicConfigSetEvent, err := e.boundContract.EncodeCallArgs(
		"emit_dynamic_config_set_event",
		[]string{},
		[]string{},
		nil,
	)
	if err != nil {
		e.lggr.Errorw("Failed to encode call", "error", err)
		return nil, err
	}

	return e.executePTBCall(ctx, ptb, addToPTB, emitDynamicConfigSetEvent, "emit_dynamic_config_set_event")
}

func (e *OffRampEmitterImpl) EmitSourceChainConfigSetEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool) (*models.SuiTransactionBlockResponse, error) {
	e.lggr.Infow("Emitting source chain config set event")
	emitSourceChainConfigSetEvent, err := e.boundContract.EncodeCallArgs(
		"emit_source_chain_config_set_event",
		[]string{},
		[]string{"u64"},
		[]any{1},
	)
	if err != nil {
		e.lggr.Errorw("Failed to encode call", "error", err)
		return nil, err
	}

	return e.executePTBCall(ctx, ptb, addToPTB, emitSourceChainConfigSetEvent, "emit_source_chain_config_set_event")
}

func (e *OffRampEmitterImpl) EmitSkippedAlreadyExecutedEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool) (*models.SuiTransactionBlockResponse, error) {
	e.lggr.Infow("Emitting skipped already executed event")
	emitSkippedAlreadyExecutedEvent, err := e.boundContract.EncodeCallArgs(
		"emit_skipped_already_executed_event",
		[]string{},
		[]string{"u64", "u64"},
		[]any{1, 1},
	)
	if err != nil {
		e.lggr.Errorw("Failed to encode call", "error", err)
		return nil, err
	}

	return e.executePTBCall(ctx, ptb, addToPTB, emitSkippedAlreadyExecutedEvent, "emit_skipped_already_executed_event")
}

func (e *OffRampEmitterImpl) EmitExecutionStateChangedEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool) (*models.SuiTransactionBlockResponse, error) {
	e.lggr.Infow("Emitting execution state changed event")
	emitExecutionStateChangedEvent, err := e.boundContract.EncodeCallArgs(
		"emit_execution_state_changed_event",
		[]string{},
		[]string{"u64"},
		[]any{1},
	)
	if err != nil {
		e.lggr.Errorw("Failed to encode call", "error", err)
		return nil, err
	}

	return e.executePTBCall(ctx, ptb, addToPTB, emitExecutionStateChangedEvent, "emit_execution_state_changed_event")
}

func (e *OffRampEmitterImpl) EmitCommitReportAcceptedEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool) (*models.SuiTransactionBlockResponse, error) {
	e.lggr.Infow("Emitting commit report accepted event")
	emitCommitReportAcceptedEvent, err := e.boundContract.EncodeCallArgs(
		"emit_commit_report_accepted_event",
		[]string{},
		[]string{},
		nil,
	)
	if err != nil {
		e.lggr.Errorw("Failed to encode call", "error", err)
		return nil, err
	}

	return e.executePTBCall(ctx, ptb, addToPTB, emitCommitReportAcceptedEvent, "emit_commit_report_accepted_event")
}

func (e *OffRampEmitterImpl) EmitSkippedReportExecutionEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool) (*models.SuiTransactionBlockResponse, error) {
	e.lggr.Infow("Emitting skipped report execution event")
	emitSkippedReportExecutionEvent, err := e.boundContract.EncodeCallArgs(
		"emit_skipped_report_execution_event",
		[]string{},
		[]string{"u64"},
		[]any{1},
	)
	if err != nil {
		e.lggr.Errorw("Failed to encode call", "error", err)
		return nil, err
	}

	return e.executePTBCall(ctx, ptb, addToPTB, emitSkippedReportExecutionEvent, "emit_skipped_report_execution_event")
}

func (e *OffRampEmitterImpl) EmitOcrConfigEvent(ctx context.Context, gasBudget uint64, ptb *transaction.Transaction, addToPTB bool) (*models.SuiTransactionBlockResponse, error) {
	// emit_ocr_config_event
	e.lggr.Infow("Emitting ocr config event")
	emitOcrConfigEvent, err := e.boundContract.EncodeCallArgs(
		"emit_ocr_config_event",
		[]string{},
		[]string{},
		nil,
	)
	if err != nil {
		e.lggr.Errorw("Failed to encode call", "error", err)
		return nil, err
	}

	return e.executePTBCall(ctx, ptb, addToPTB, emitOcrConfigEvent, "emit_ocr_config_event")
}
