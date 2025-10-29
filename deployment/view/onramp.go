package view

type OnRampView struct {
	ContractMetaData

	StaticConfig          OnRampStaticConfig               `json:"staticConfig"`
	DynamicConfig         OnRampDynamicConfig              `json:"dynamicConfig"`
	SourceTokenToPool     map[string]string                `json:"sourceTokenToPool"`
	DestChainSpecificData map[uint64]DestChainSpecificData `json:"destChainSpecificData"`
}

type OnRampStaticConfig struct {
	ChainSelector uint64 `json:"chainSelector"`
}

type OnRampDynamicConfig struct {
	FeeAggregator  string `json:"feeAggregator"`
	AllowlistAdmin string `json:"allowlistAdmin"`
}

type DestChainSpecificData struct {
	AllowedSendersList []string              `json:"allowedSendersList"`
	DestChainConfig    OnRampDestChainConfig `json:"destChainConfig"`
	ExpectedNextSeqNum uint64                `json:"expectedNextSeqNum"`
}

type OnRampDestChainConfig struct {
	SequenceNumber   uint64 `json:"sequenceNumber"`
	AllowlistEnabled bool   `json:"allowlistEnabled"`
	Router           string `json:"router"`
}

// GenerateOnRampView generates an onramp view for a given onramp address
// This is a mocked implementation for now
func GenerateOnRampView(onRampAddress string, routerAddress string, ccipAddress string) (OnRampView, error) {
	// Mock data for now
	sourceTokenToPool := map[string]string{
		"0xtoken1": "0xpool1",
		"0xtoken2": "0xpool2",
	}

	destChainSpecificData := map[uint64]DestChainSpecificData{
		123456789: {
			AllowedSendersList: []string{"0xsender1", "0xsender2"},
			DestChainConfig: OnRampDestChainConfig{
				SequenceNumber:   1,
				AllowlistEnabled: true,
				Router:           routerAddress,
			},
			ExpectedNextSeqNum: 100,
		},
		987654321: {
			AllowedSendersList: []string{"0xsender3"},
			DestChainConfig: OnRampDestChainConfig{
				SequenceNumber:   5,
				AllowlistEnabled: false,
				Router:           routerAddress,
			},
			ExpectedNextSeqNum: 250,
		},
	}

	return OnRampView{
		ContractMetaData: ContractMetaData{
			Address:        onRampAddress,
			Owner:          "0xmock_onramp_owner",
			TypeAndVersion: "OnRamp 1.0.0",
		},
		StaticConfig: OnRampStaticConfig{
			ChainSelector: 123456789,
		},
		DynamicConfig: OnRampDynamicConfig{
			FeeAggregator:  "0xmock_fee_aggregator",
			AllowlistAdmin: "0xmock_allowlist_admin",
		},
		SourceTokenToPool:     sourceTokenToPool,
		DestChainSpecificData: destChainSpecificData,
	}, nil
}
