package changesets

import (
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/deployment"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	ccipops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip"
)

var _ cldf.ChangeSetV2[UnregisterTokenPoolConfig] = UnregisterTokenPool{}

// UnregisterTokenPool removes a token's registration from the CCIP Token Admin
// Registry using the deployer key. It calls the public admin-gated
// token_admin_registry::unregister_pool, which asserts the caller is the
// token's configured administrator — not the pool OwnerCap holder. This is used
// to clear a stale on-chain registration so a pool can be redeployed and
// re-registered, since pool initialization registers unconditionally and aborts
// with ETokenAlreadyRegistered when the coin is already present.
type UnregisterTokenPool struct{}

type UnregisterTokenPoolConfig struct {
	ChainSelector uint64 `json:"chainSelector" yaml:"chainSelector"`

	// CoinMetadataAddress is the CoinMetadata object id of the token to
	// unregister. The Token Admin Registry keys registrations on the CoinMetadata
	// object id, not the coin-type string.
	CoinMetadataAddress string `json:"coinMetadataAddress" yaml:"coinMetadataAddress"`

	// OwnerCapObjectId is only used to build MCMS callback data, which the EOA
	// path does not send. Leave empty for deployer-key execution.
	OwnerCapObjectId string `json:"ownerCapObjectId,omitempty" yaml:"ownerCapObjectId,omitempty"`
}

// Apply implements deployment.ChangeSetV2.
func (d UnregisterTokenPool) Apply(e cldf.Environment, config UnregisterTokenPoolConfig) (cldf.ChangesetOutput, error) {
	suiState, err := deployment.LoadOnchainStatesui(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	state, ok := suiState[config.ChainSelector]
	if !ok {
		return cldf.ChangesetOutput{}, fmt.Errorf("no Sui chain state for chain selector %d", config.ChainSelector)
	}

	ccipPackageId := state.EffectiveCCIPPackageID()
	if ccipPackageId == "" {
		return cldf.ChangesetOutput{}, fmt.Errorf("ccip package id not found for chain selector %d", config.ChainSelector)
	}
	if state.CCIPObjectRef == "" {
		return cldf.ChangesetOutput{}, fmt.Errorf("ccip object ref not found for chain selector %d", config.ChainSelector)
	}

	suiChain := e.BlockChains.SuiChains()[config.ChainSelector]
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

	report, err := cld_ops.ExecuteOperation(e.OperationsBundle, ccipops.TokenAdminRegistryUnregisterPoolOp, deps, ccipops.UnregisterPoolInput{
		CCIPPackageId:       ccipPackageId,
		CCIPObjectRef:       state.CCIPObjectRef,
		OwnerCapObjectId:    config.OwnerCapObjectId,
		CoinMetadataAddress: config.CoinMetadataAddress,
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to unregister token pool: %w", err)
	}

	return cldf.ChangesetOutput{
		Reports: []cld_ops.Report[any, any]{report.ToGenericReport()},
	}, nil
}

// VerifyPreconditions implements deployment.ChangeSetV2.
func (d UnregisterTokenPool) VerifyPreconditions(e cldf.Environment, config UnregisterTokenPoolConfig) error {
	if config.ChainSelector == 0 {
		return fmt.Errorf("chainSelector is required")
	}
	if config.CoinMetadataAddress == "" {
		return fmt.Errorf("coinMetadataAddress is required")
	}
	return nil
}
