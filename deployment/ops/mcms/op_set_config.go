package mcmsops

import (
	"github.com/Masterminds/semver/v3"
	"github.com/block-vision/sui-go-sdk/models"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	suisdk "github.com/smartcontractkit/mcms/sdk/sui"
	"github.com/smartcontractkit/mcms/types"
)

type MCMSSetConfigInput struct {
	ChainSelector uint64 `json:"chainSelector"`
	// MCMS related
	McmsPackageID string `json:"mcmsPackageID"`
	OwnerCap      string `json:"ownerCap"`
	McmsObjectID  string `json:"mcmsObjectID"`
	// Timelock related
	Role suisdk.TimelockRole `json:"role"`
	// Config related
	Config types.Config `json:"config"`
}

var setConfigMcmsHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input MCMSSetConfigInput) (output sui_ops.OpTxResult[cld_ops.EmptyInput], err error) {
	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	mcmsConfigurer, err := suisdk.NewConfigurer(deps.Client, deps.Signer, input.Role, input.McmsPackageID, input.OwnerCap, input.ChainSelector)
	if err != nil {
		return sui_ops.OpTxResult[cld_ops.EmptyInput]{}, err
	}

	mcmsTx, err := mcmsConfigurer.SetConfig(b.GetContext(), input.McmsObjectID, &input.Config, true)
	if err != nil {
		return sui_ops.OpTxResult[cld_ops.EmptyInput]{}, err
	}

	suiTx := mcmsTx.RawData.(*models.SuiTransactionBlockResponse)

	return sui_ops.OpTxResult[cld_ops.EmptyInput]{
		Digest:    suiTx.Digest,
		PackageId: input.McmsPackageID,
	}, err
}

var SetConfigMCMSOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("mcms", "mcms", "set_config"),
	semver.MustParse("0.1.0"),
	"Set config in the MCMS contract",
	setConfigMcmsHandler,
)
