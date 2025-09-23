package events

import (
	"context"

	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/block-vision/sui-go-sdk/transaction"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	rel "github.com/smartcontractkit/chainlink-sui/relayer/signer"
)

func EmitOffRampEvents(
	ctx context.Context,
	lggr logger.Logger,
	packageId string,
	signer rel.SuiSigner,
	callOpts bind.CallOpts,
	client sui.ISuiAPI,
	accountAddress string,
	gasBudget uint64,
) error {
	lggr.Infow("Starting emit-static-config-set-event command")

	ptb := transaction.NewTransaction()
	sourceChainSelector := uint64(1)

	offrampBoundContract, err := bind.NewBoundContract(
		packageId,
		"test",
		"offramp",
		client,
	)
	if err != nil {
		lggr.Errorw("Failed to create offramp bound contract", "error", err)
		return err
	}

	// emit_static_config_set_event
	emitStaticConfigSetEvent, err := offrampBoundContract.EncodeCallArgsWithGenerics(
		"emit_static_config_set_event",
		[]string{},
		[]string{},
		[]string{"u64"},
		[]any{1},
		nil,
	)
	if err != nil {
		lggr.Errorw("Failed to encode call", "error", err)
		return err
	}

	_, err = offrampBoundContract.AppendPTB(ctx, &callOpts, ptb, emitStaticConfigSetEvent)
	if err != nil {
		lggr.Errorw("Failed to append PTB", "error", err)
		return err
	}

	// emit_dynamic_config_set_event
	// emitDynamicConfigSetEvent, err := offrampBoundContract.EncodeCallArgs(
	// 	"emit_dynamic_config_set_event",
	// 	[]string{},
	// 	[]string{},
	// 	nil,
	// )
	// if err != nil {
	// 	lggr.Errorw("Failed to encode call", "error", err)
	// 	return err
	// }

	// _, err = offrampBoundContract.AppendPTB(ctx, &callOpts, ptb, emitDynamicConfigSetEvent)
	// if err != nil {
	// 	lggr.Errorw("Failed to append PTB", "error", err)
	// 	return err
	// }

	// emit_source_chain_config_set_event
	emitSourceChainConfigSetEvent, err := offrampBoundContract.EncodeCallArgs(
		"emit_source_chain_config_set_event",
		[]string{},
		[]string{"u64"},
		[]any{sourceChainSelector},
	)
	if err != nil {
		lggr.Errorw("Failed to encode call", "error", err)
		return err
	}

	_, err = offrampBoundContract.AppendPTB(ctx, &callOpts, ptb, emitSourceChainConfigSetEvent)
	if err != nil {
		lggr.Errorw("Failed to append PTB", "error", err)
		return err
	}

	// emit_commit_report_accepted_event
	emitCommitReportAcceptedEvent, err := offrampBoundContract.EncodeCallArgs(
		"emit_commit_report_accepted_event",
		[]string{},
		[]string{},
		nil,
	)
	if err != nil {
		lggr.Errorw("Failed to encode call", "error", err)
		return err
	}

	_, err = offrampBoundContract.AppendPTB(ctx, &callOpts, ptb, emitCommitReportAcceptedEvent)
	if err != nil {
		lggr.Errorw("Failed to append PTB", "error", err)
		return err
	}

	tx, err := bind.ExecutePTB(ctx, &callOpts, client, ptb)
	if err != nil {
		lggr.Errorw("Failed to execute transaction", "error", err)
		return err
	}

	lggr.Infow("Event emission completed (scaffolding)", "tx", tx)

	return nil
}
