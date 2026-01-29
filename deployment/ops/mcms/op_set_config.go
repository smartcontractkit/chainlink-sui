package mcmsops

import (
	"fmt"
	"math/big"

	"github.com/Masterminds/semver/v3"

	cselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/mcms/sdk"
	suisdk "github.com/smartcontractkit/mcms/sdk/sui"
	"github.com/smartcontractkit/mcms/types"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	modulemcms "github.com/smartcontractkit/chainlink-sui/bindings/generated/mcms/mcms"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
)

type MCMSSetConfigInput struct {
	ChainSelector uint64 `yaml:"chainSelector"`
	// MCMS related
	McmsPackageID string `yaml:"mcmsPackageID"`
	OwnerCap      string `yaml:"ownerCap"`
	McmsObjectID  string `yaml:"mcmsObjectID"`
	// Timelock related
	Role suisdk.TimelockRole `yaml:"role"`
	// Config related
	Config    types.Config `yaml:"config"`
	ClearRoot bool         `yaml:"clearRoot"`
}

var SetConfigMCMSOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("mcms", "mcms", "set_config"),
	semver.MustParse("0.1.0"),
	"Set config in the MCMS contract",
	setConfigMcmsHandler,
)

var setConfigMcmsHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input MCMSSetConfigInput) (output sui_ops.OpTxResult[cld_ops.EmptyInput], err error) {
	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	mcms, err := modulemcms.NewMcms(input.McmsPackageID, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[cld_ops.EmptyInput]{}, err
	}

	chainID, err := cselectors.SuiChainIdFromSelector(input.ChainSelector)
	if err != nil {
		return sui_ops.OpTxResult[cld_ops.EmptyInput]{}, err
	}
	chainIDBig := new(big.Int).SetUint64(chainID)
	groupQuorum, groupParents, signerAddresses, signerGroups, err := sdk.ExtractSetConfigInputs(&input.Config)

	if err != nil {
		return sui_ops.OpTxResult[cld_ops.EmptyInput]{}, fmt.Errorf("unable to extract set config inputs: %w", err)
	}
	signers := make([][]byte, len(signerAddresses))
	for i, addr := range signerAddresses {
		signers[i] = addr.Bytes()
	}

	encodedCall, err := mcms.Encoder().SetConfig(
		bind.Object{Id: input.OwnerCap},
		bind.Object{Id: input.McmsObjectID},
		input.Role.Byte(),
		chainIDBig,
		signers,
		signerGroups,
		groupQuorum[:],
		groupParents[:],
		input.ClearRoot,
	)
	if err != nil {
		return sui_ops.OpTxResult[cld_ops.EmptyInput]{}, fmt.Errorf("failed to encode SetConfig call: %w", err)
	}
	call, err := sui_ops.ToTransactionCall(encodedCall, input.McmsObjectID)
	if err != nil {
		return sui_ops.OpTxResult[cld_ops.EmptyInput]{}, fmt.Errorf("failed to convert encoded call to TransactionCall: %w", err)
	}

	if deps.Signer == nil {
		b.Logger.Infow("Skipping execution of AcceptOwnership on StateObject as per no Signer provided")
		return sui_ops.OpTxResult[cld_ops.EmptyInput]{
			Digest:    "",
			PackageId: input.McmsPackageID,
			Call:      call,
		}, nil
	}

	tx, err := mcms.Bound().ExecuteTransaction(
		b.GetContext(),
		opts,
		encodedCall,
	)

	return sui_ops.OpTxResult[cld_ops.EmptyInput]{
		Call:      call,
		Digest:    tx.Digest,
		PackageId: input.McmsPackageID,
	}, err
}
