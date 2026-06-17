package changesets

import (
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/deployment"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	ccipops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip"
)

type BackfillLocalDecimalsConfig struct {
	SuiChainSelector uint64
	CCIPPackageId    string
	StateObjectId    string
	OwnerCapObjectId string
	VerifyOnly       bool
}

var _ cldf.ChangeSetV2[BackfillLocalDecimalsConfig] = BackfillLocalDecimals{}

type BackfillLocalDecimals struct{}

func (d BackfillLocalDecimals) Apply(e cldf.Environment, config BackfillLocalDecimalsConfig) (cldf.ChangesetOutput, error) {
	seqReports := make([]operations.Report[any, any], 0)

	suiChain := e.BlockChains.SuiChains()[config.SuiChainSelector]
	deps := sui_ops.OpTxDeps{
		Client: suiChain.Client,
		Signer: suiChain.Signer,
		GetCallOpts: func() *bind.CallOpts {
			gasBudget := uint64(400_000_000)
			return &bind.CallOpts{
				WaitForExecution: true,
				GasBudget:        &gasBudget,
			}
		},
		SuiRPC: suiChain.URL,
	}

	verifyReport, err := operations.ExecuteOperation(
		e.OperationsBundle,
		ccipops.TokenAdminRegistryVerifyLocalDecimalsOp,
		deps,
		ccipops.VerifyLocalDecimalsInput{
			CCIPPackageId: config.CCIPPackageId,
			StateObjectId: config.StateObjectId,
		},
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to verify local decimals for Sui chain %d: %w", config.SuiChainSelector, err)
	}
	seqReports = append(seqReports, verifyReport.ToGenericReport())

	if len(verifyReport.Output.Objects.Mismatches) > 0 && config.VerifyOnly {
		return cldf.ChangesetOutput{Reports: seqReports}, fmt.Errorf(
			"local decimals verification failed for %d token(s) on chain %d",
			len(verifyReport.Output.Objects.Mismatches),
			config.SuiChainSelector,
		)
	}

	if config.VerifyOnly {
		return cldf.ChangesetOutput{Reports: seqReports}, nil
	}

	backfillReport, err := operations.ExecuteOperation(
		e.OperationsBundle,
		ccipops.TokenAdminRegistryBackfillAllLocalDecimalsOp,
		deps,
		ccipops.BackfillAllLocalDecimalsInput{
			CCIPPackageId:    config.CCIPPackageId,
			StateObjectId:    config.StateObjectId,
			OwnerCapObjectId: config.OwnerCapObjectId,
		},
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to backfill local decimals for Sui chain %d: %w", config.SuiChainSelector, err)
	}
	seqReports = append(seqReports, backfillReport.ToGenericReport())

	verifyAfterReport, err := operations.ExecuteOperation(
		e.OperationsBundle,
		ccipops.TokenAdminRegistryVerifyLocalDecimalsOp,
		deps,
		ccipops.VerifyLocalDecimalsInput{
			CCIPPackageId: config.CCIPPackageId,
			StateObjectId: config.StateObjectId,
		},
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to re-verify local decimals for Sui chain %d: %w", config.SuiChainSelector, err)
	}
	seqReports = append(seqReports, verifyAfterReport.ToGenericReport())

	if len(verifyAfterReport.Output.Objects.Mismatches) > 0 {
		return cldf.ChangesetOutput{Reports: seqReports}, fmt.Errorf(
			"local decimals still mismatched for %d token(s) after backfill on chain %d",
			len(verifyAfterReport.Output.Objects.Mismatches),
			config.SuiChainSelector,
		)
	}

	return cldf.ChangesetOutput{Reports: seqReports}, nil
}

func (d BackfillLocalDecimals) VerifyPreconditions(e cldf.Environment, config BackfillLocalDecimalsConfig) error {
	if config.SuiChainSelector == 0 {
		return fmt.Errorf("sui chain selector is required")
	}
	if config.CCIPPackageId == "" || config.StateObjectId == "" {
		return fmt.Errorf("ccip package id and state object id are required")
	}
	if !config.VerifyOnly && config.OwnerCapObjectId == "" {
		return fmt.Errorf("owner cap object id is required when backfilling")
	}

	state, err := deployment.LoadOnchainStatesui(e)
	if err != nil {
		return err
	}
	chainState, ok := state[config.SuiChainSelector]
	if !ok {
		return fmt.Errorf("no onchain state for sui chain %d", config.SuiChainSelector)
	}
	if config.CCIPPackageId != chainState.CCIPAddress {
		return fmt.Errorf("ccip package id does not match loaded chain state")
	}

	return nil
}
