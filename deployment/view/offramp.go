package view

type OffRampView struct {
	ContractMetaData

	LatestPriceSequenceNumber uint64                              `json:"latestPriceSequenceNumber"`
	StaticConfig              OffRampStaticConfig                 `json:"staticConfig"`
	DynamicConfig             OffRampDynamicConfig                `json:"dynamicConfig"`
	SourceChainConfigs        map[uint64]OffRampSourceChainConfig `json:"sourceChainConfigs"`
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

// GenerateOffRampView generates an offramp view for a given offramp address
// This is a mocked implementation for now
func GenerateOffRampView(offRampAddress string, routerAddress string) (OffRampView, error) {
	// Mock data for now
	sourceChainConfigs := map[uint64]OffRampSourceChainConfig{
		123456789: {
			Router:                    routerAddress,
			IsEnabled:                 true,
			MinSeqNr:                  1,
			IsRMNVerificationDisabled: false,
			OnRamp:                    "0xmock_onramp_1",
		},
		987654321: {
			Router:                    routerAddress,
			IsEnabled:                 true,
			MinSeqNr:                  5,
			IsRMNVerificationDisabled: false,
			OnRamp:                    "0xmock_onramp_2",
		},
	}

	return OffRampView{
		ContractMetaData: ContractMetaData{
			Address:        offRampAddress,
			Owner:          "0xmock_offramp_owner",
			TypeAndVersion: "OffRamp 1.0.0",
		},
		LatestPriceSequenceNumber: 12345,
		StaticConfig: OffRampStaticConfig{
			ChainSelector:      123456789,
			RMNRemote:          "0xmock_rmn_remote",
			TokenAdminRegistry: "0xmock_token_admin_registry",
			NonceManager:       "0xmock_nonce_manager",
		},
		DynamicConfig: OffRampDynamicConfig{
			FeeQuoter:                               "0xmock_fee_quoter",
			PermissionlessExecutionThresholdSeconds: 3600,
		},
		SourceChainConfigs: sourceChainConfigs,
	}, nil
}
