package mcmsops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/mcms"

	suisdk "github.com/smartcontractkit/mcms/sdk/sui"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/contracts"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	"github.com/smartcontractkit/chainlink-sui/deployment/utils"
)

type UpgradeCCIPInput struct {
	PackageName     contracts.Package `json:"packageName" validate:"required"`
	TargetPackageId string            `json:"targetPackageId" validate:"required"`
	NamedAddresses  map[string]string `json:"namedAddresses"`

	// Chain related
	ChainSelector uint64 `json:"chainSelector"`

	// MCMS related
	MmcsPackageID      string `json:"mcmsPackageID"`
	McmsStateObjID     string `json:"mcmsStateObjID"`
	RegistryObjID      string `json:"registryObjID"`
	TimelockObjID      string `json:"timelockObjID"`
	AccountObjID       string `json:"accountObjID"`
	DeployerStateObjID string `json:"deployerStateObjID"`
	OwnerCapObjID      string `json:"ownerCapObjID"`

	// Timelock related
	TimelockConfig utils.TimelockConfig `json:"timelockConfig"`
}

var upgradeHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input UpgradeCCIPInput) (output mcms.TimelockProposal, err error) {
	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer

	// Get the actual signer address for compilation
	signerAddress, err := deps.Signer.GetAddress()
	if err != nil {
		return mcms.TimelockProposal{}, err
	}

	namedAddresses := make(map[string]string)
	for k, v := range input.NamedAddresses {
		namedAddresses[k] = v
	}
	namedAddresses["signer"] = signerAddress

	artifact, err := bind.CompilePackage(input.PackageName, namedAddresses, true, deps.SuiRPC)
	if err != nil {
		return mcms.TimelockProposal{}, err
	}

	mcmsTx, err := suisdk.CreateUpgradeTransaction(artifact, input.MmcsPackageID, input.DeployerStateObjID, input.RegistryObjID, input.OwnerCapObjID, input.TargetPackageId)
	if err != nil {
		return mcms.TimelockProposal{}, err
	}

	op := mcmstypes.BatchOperation{
		ChainSelector: mcmstypes.ChainSelector(input.ChainSelector),
		Transactions:  []mcmstypes.Transaction{mcmsTx},
	}

	proposalInput := utils.GenerateProposalInput{
		Client:             deps.Client,
		MCMSPackageID:      input.MmcsPackageID,
		MCMSStateObjID:     input.McmsStateObjID,
		TimelockObjID:      input.TimelockObjID,
		AccountObjID:       input.AccountObjID,
		RegistryObjID:      input.RegistryObjID,
		DeployerStateObjID: input.DeployerStateObjID,
		ChainSelector:      input.ChainSelector,
		TimelockConfig:     input.TimelockConfig,
		Description:        fmt.Sprintf("Upgrade the %s package to the latest version", input.PackageName),
		BatchOp:            op,
	}
	timelockProposal, err := utils.GenerateProposal(b.GetContext(), proposalInput)
	if err != nil {
		return mcms.TimelockProposal{}, fmt.Errorf("failed to build proposal: %w", err)
	}

	return *timelockProposal, nil
}

var UpgradeCCIPOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("mcms", "package", "upgrade"),
	semver.MustParse("0.1.0"),
	"Returns the MCMS proposal that upgrades a CCIP package",
	upgradeHandler,
)
