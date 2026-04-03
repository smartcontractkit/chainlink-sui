package changesets

import (
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	ccipops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip"
	offrampops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_offramp"
	onrampops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_onramp"
)

type AddPackageIdTarget string

const (
	AddPackageIdTargetCCIP    AddPackageIdTarget = "ccip"
	AddPackageIdTargetOnRamp  AddPackageIdTarget = "onramp"
	AddPackageIdTargetOffRamp AddPackageIdTarget = "offramp"
)

type AddCCIPPackageIdConfig struct {
	SuiChainSelector uint64             `yaml:"suiChainSelector"`
	Target           AddPackageIdTarget `yaml:"target"`
	ModulePackageId  string             `yaml:"modulePackageId"`
	StateObjectId    string             `yaml:"stateObjectId"`
	OwnerCapObjectId string             `yaml:"ownerCapObjectId"`
	PackageId        string             `yaml:"packageId"`
}

var _ cldf.ChangeSetV2[AddCCIPPackageIdConfig] = AddCCIPPackageId{}

type AddCCIPPackageId struct{}

func (d AddCCIPPackageId) Apply(e cldf.Environment, config AddCCIPPackageIdConfig) (cldf.ChangesetOutput, error) {
	ab := cldf.NewMemoryAddressBook()

	suiChain := e.BlockChains.SuiChains()[config.SuiChainSelector]

	deps := sui_ops.OpTxDeps{
		Client: suiChain.Client,
		Signer: suiChain.Signer,
		GetCallOpts: func() *bind.CallOpts {
			b := uint64(400_000_000)
			return &bind.CallOpts{
				WaitForExecution: true,
				GasBudget:        &b,
			}
		},
		SuiRPC: suiChain.URL,
	}

	report, err := executeAddPackageId(e, config, deps)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	return cldf.ChangesetOutput{
		AddressBook: ab,
		Reports:     []operations.Report[any, any]{report},
	}, nil
}

func executeAddPackageId(e cldf.Environment, config AddCCIPPackageIdConfig, deps sui_ops.OpTxDeps) (operations.Report[any, any], error) {
	target := config.Target
	if target == "" {
		target = AddPackageIdTargetCCIP
	}

	switch target {
	case AddPackageIdTargetCCIP:
		r, err := operations.ExecuteOperation(e.OperationsBundle, ccipops.AddPackageIdStateObjectOp, deps, ccipops.AddPackageIdStateObjectInput{
			CCIPPackageId:         config.ModulePackageId,
			CCIPObjectRefObjectId: config.StateObjectId,
			OwnerCapObjectId:      config.OwnerCapObjectId,
			PackageId:             config.PackageId,
		})
		if err != nil {
			return operations.Report[any, any]{}, fmt.Errorf("failed to add package ID to CCIP state object for Sui chain %d: %w", config.SuiChainSelector, err)
		}
		return r.ToGenericReport(), nil

	case AddPackageIdTargetOnRamp:
		r, err := operations.ExecuteOperation(e.OperationsBundle, onrampops.AddPackageIdOp, deps, onrampops.AddPackageIdInput{
			OnRampPackageId:  config.ModulePackageId,
			StateObjectId:    config.StateObjectId,
			OwnerCapObjectId: config.OwnerCapObjectId,
			PackageId:        config.PackageId,
		})
		if err != nil {
			return operations.Report[any, any]{}, fmt.Errorf("failed to add package ID to OnRamp for Sui chain %d: %w", config.SuiChainSelector, err)
		}
		return r.ToGenericReport(), nil

	case AddPackageIdTargetOffRamp:
		r, err := operations.ExecuteOperation(e.OperationsBundle, offrampops.AddPackageIdOffRampOp, deps, offrampops.AddPackageIdOffRampInput{
			OffRampPackageId: config.ModulePackageId,
			StateObjectId:    config.StateObjectId,
			OwnerCapObjectId: config.OwnerCapObjectId,
			PackageId:        config.PackageId,
		})
		if err != nil {
			return operations.Report[any, any]{}, fmt.Errorf("failed to add package ID to OffRamp for Sui chain %d: %w", config.SuiChainSelector, err)
		}
		return r.ToGenericReport(), nil

	default:
		return operations.Report[any, any]{}, fmt.Errorf("unknown target %q: must be one of %q, %q, %q",
			config.Target, AddPackageIdTargetCCIP, AddPackageIdTargetOnRamp, AddPackageIdTargetOffRamp)
	}
}

func (d AddCCIPPackageId) VerifyPreconditions(e cldf.Environment, config AddCCIPPackageIdConfig) error {
	if config.Target != "" && config.Target != AddPackageIdTargetCCIP &&
		config.Target != AddPackageIdTargetOnRamp && config.Target != AddPackageIdTargetOffRamp {
		return fmt.Errorf("invalid target %q: must be one of %q, %q, %q",
			config.Target, AddPackageIdTargetCCIP, AddPackageIdTargetOnRamp, AddPackageIdTargetOffRamp)
	}
	return nil
}
