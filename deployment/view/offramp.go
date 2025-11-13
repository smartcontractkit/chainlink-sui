package view

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain/sui"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_offramp "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_offramp/offramp"
)

type OffRampView struct {
	ContractMetaData

	StaticConfig       OffRampStaticConfig                 `json:"staticConfig"`
	DynamicConfig      OffRampDynamicConfig                `json:"dynamicConfig"`
	SourceChainConfigs map[uint64]OffRampSourceChainConfig `json:"sourceChainConfigs"`
}

type OffRampStaticConfig struct {
	ChainSelector      uint64 `json:"chainSelector"`
	RMNRemote          string `json:"rmnRemote"`
	TokenAdminRegistry string `json:"tokenAdminRegistry"`
	NonceManager       string `json:"nonceManager"`
}

type OffRampDynamicConfig struct {
	FeeQuoter                               string `json:"feeQuoter"`
	PermissionlessExecutionThresholdSeconds uint32 `json:"permissionlessExecutionThresholdSeconds"`
}

type OffRampSourceChainConfig struct {
	Router                    string `json:"router"`
	IsEnabled                 bool   `json:"isEnabled"`
	MinSeqNr                  uint64 `json:"minSeqNr"`
	IsRMNVerificationDisabled bool   `json:"isRMNVerificationDisabled"`
	OnRamp                    string `json:"onRamp"`
}

// GenerateOffRampView generates an offramp view for a given offramp by querying the on-chain state
func GenerateOffRampView(
	ctx context.Context,
	chain sui.Chain,
	offRampPackageID string,
	offRampStateObjectID string,
	ccipObjectRef string,
) (OffRampView, error) {
	if offRampPackageID == "" || offRampStateObjectID == "" {
		return OffRampView{}, fmt.Errorf("offRampPackageID and offRampStateObjectID cannot be empty")
	}

	offRampContract, err := module_offramp.NewOfframp(offRampPackageID, chain.Client)
	if err != nil {
		return OffRampView{}, fmt.Errorf("failed to create offramp contract binding: %w", err)
	}
	offRampStateObj := bind.Object{Id: offRampStateObjectID}
	ccipRefObj := bind.Object{Id: ccipObjectRef}
	callOpts := &bind.CallOpts{Signer: chain.Signer}

	// Get owner
	owner, err := offRampContract.DevInspect().Owner(ctx, callOpts, offRampStateObj)
	if err != nil {
		return OffRampView{}, fmt.Errorf("failed to get owner: %w", err)
	}

	// Get type and version
	typeAndVersion, err := offRampContract.DevInspect().TypeAndVersion(ctx, callOpts)
	if err != nil {
		return OffRampView{}, fmt.Errorf("failed to get type and version: %w", err)
	}

	// Get static config
	staticConfig, err := offRampContract.DevInspect().GetStaticConfig(ctx, callOpts, ccipRefObj, offRampStateObj)
	if err != nil {
		return OffRampView{}, fmt.Errorf("failed to get static config: %w", err)
	}

	// Get dynamic config
	dynamicConfig, err := offRampContract.DevInspect().GetDynamicConfig(ctx, callOpts, ccipRefObj, offRampStateObj)
	if err != nil {
		return OffRampView{}, fmt.Errorf("failed to get dynamic config: %w", err)
	}

	// Get all source chain configs
	sourceChainConfigsRaw, err := offRampContract.DevInspect().GetAllSourceChainConfigs(ctx, callOpts, ccipRefObj, offRampStateObj)
	if err != nil {
		return OffRampView{}, fmt.Errorf("failed to get source chain configs: %w", err)
	}

	// Parse source chain configs
	// GetAllSourceChainConfigs returns [0]: vector<u64>, [1]: vector<SourceChainConfig>
	sourceChainConfigs := make(map[uint64]OffRampSourceChainConfig)
	if len(sourceChainConfigsRaw) >= 2 {
		selectors, ok := sourceChainConfigsRaw[0].([]uint64)
		if !ok {
			return OffRampView{}, fmt.Errorf("unexpected type for source chain selectors: got %T", sourceChainConfigsRaw[0])
		}
		configs, ok := sourceChainConfigsRaw[1].([]module_offramp.SourceChainConfig)
		if !ok {
			return OffRampView{}, fmt.Errorf("unexpected type for source chain configs: got %T", sourceChainConfigsRaw[1])
		}

		for i, selector := range selectors {
			if i < len(configs) {
				sourceChainConfigs[selector] = OffRampSourceChainConfig{
					Router:                    configs[i].Router,
					IsEnabled:                 configs[i].IsEnabled,
					MinSeqNr:                  configs[i].MinSeqNr,
					IsRMNVerificationDisabled: configs[i].IsRmnVerificationDisabled,
					OnRamp:                    fmt.Sprintf("0x%x", configs[i].OnRamp),
				}
			}
		}
	}

	return OffRampView{
		ContractMetaData: ContractMetaData{
			Address:        offRampPackageID,
			Owner:          owner,
			TypeAndVersion: typeAndVersion,
			StateObjectID:  offRampStateObjectID,
		},
		StaticConfig: OffRampStaticConfig{
			ChainSelector:      staticConfig.ChainSelector,
			RMNRemote:          staticConfig.RmnRemote,
			TokenAdminRegistry: staticConfig.TokenAdminRegistry,
			NonceManager:       staticConfig.NonceManager,
		},
		DynamicConfig: OffRampDynamicConfig{
			FeeQuoter:                               dynamicConfig.FeeQuoter,
			PermissionlessExecutionThresholdSeconds: dynamicConfig.PermissionlessExecutionThresholdSeconds,
		},
		SourceChainConfigs: sourceChainConfigs,
	}, nil

}
