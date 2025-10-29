package view

import (
	"fmt"
)

type CCIPView struct {
	ContractMetaData

	FeeQuoter          FeeQuoterView          `json:"feeQuoter"`
	RMNRemote          RMNRemoteView          `json:"rmnRemote"`
	TokenAdminRegistry TokenAdminRegistryView `json:"tokenAdminRegistry"`
	NonceManager       NonceManagerView       `json:"nonceManager"`
	ReceiverRegistry   ReceiverRegistryView   `json:"receiverRegistry"`
}

type FeeQuoterView struct {
	ContractMetaData

	FeeTokens               []string                            `json:"feeTokens"`
	StaticConfig            FeeQuoterStaticConfig               `json:"staticConfig"`
	DestinationChainConfigs map[uint64]FeeQuoterDestChainConfig `json:"destinationChainConfigs"`
}

type FeeQuoterStaticConfig struct {
	MaxFeeJuelsPerMsg            string `json:"maxFeeJuelsPerMsg"`
	LinkToken                    string `json:"linkToken"`
	TokenPriceStalenessThreshold uint64 `json:"tokenPriceStalenessThreshold"`
}

type FeeQuoterDestChainConfig struct {
	IsEnabled                         bool   `json:"isEnabled"`
	MaxNumberOfTokensPerMsg           uint16 `json:"maxNumberOfTokensPerMsg"`
	MaxDataBytes                      uint32 `json:"maxDataBytes"`
	MaxPerMsgGasLimit                 uint32 `json:"maxPerMsgGasLimit"`
	DestGasOverhead                   uint32 `json:"destGasOverhead"`
	DestGasPerPayloadByteBase         uint8  `json:"destGasPerPayloadByteBase"`
	DestGasPerPayloadByteHigh         uint8  `json:"destGasPerPayloadByteHigh"`
	DestGasPerPayloadByteThreshold    uint16 `json:"destGasPerPayloadByteThreshold"`
	DestDataAvailabilityOverheadGas   uint32 `json:"destDataAvailabilityOverheadGas"`
	DestGasPerDataAvailabilityByte    uint16 `json:"destGasPerDataAvailabilityByte"`
	DestDataAvailabilityMultiplierBps uint16 `json:"destDataAvailabilityMultiplierBps"`
	ChainFamilySelector               string `json:"chainFamilySelector"`
	EnforceOutOfOrder                 bool   `json:"enforceOutOfOrder"`
	DefaultTokenFeeUsdCents           uint16 `json:"defaultTokenFeeUsdCents"`
	DefaultTokenDestGasOverhead       uint32 `json:"defaultTokenDestGasOverhead"`
	DefaultTxGasLimit                 uint32 `json:"defaultTxGasLimit"`
	GasMultiplierWeiPerEth            uint64 `json:"gasMultiplierWeiPerEth"`
	GasPriceStalenessThreshold        uint32 `json:"gasPriceStalenessThreshold"`
	NetworkFeeUsdCents                uint32 `json:"networkFeeUsdCents"`
}

type RMNRemoteView struct {
	ContractMetaData
	IsCursed             bool                     `json:"isCursed"`
	Config               RMNRemoteVersionedConfig `json:"config"`
	CursedSubjectEntries []RMNRemoteCurseEntry    `json:"cursedSubjectEntries"`
}

type RMNRemoteVersionedConfig struct {
	Version uint32            `json:"version"`
	Signers []RMNRemoteSigner `json:"signers"`
	Fsign   uint64            `json:"fSign"`
}

type RMNRemoteSigner struct {
	OnchainPublicKey string `json:"onchain_public_key"` // Follow EVM snake_case
	NodeIndex        uint64 `json:"node_index"`
}

type RMNRemoteCurseEntry struct {
	Subject  string `json:"subject"`
	Selector uint64 `json:"selector"`
}

type TokenAdminRegistryView struct {
	ContractMetaData
}

type NonceManagerView struct {
	ContractMetaData
}

type ReceiverRegistryView struct {
	ContractMetaData
}

// GenerateCCIPView generates a mocked CCIP view for SUI
func GenerateCCIPView(ccipAddress string) (CCIPView, error) {
	if ccipAddress == "" {
		return CCIPView{}, fmt.Errorf("ccipAddress cannot be empty")
	}

	// Mocked implementation - replace with actual on-chain queries later
	return CCIPView{
		ContractMetaData: ContractMetaData{
			Address:        ccipAddress,
			Owner:          "0xmocked_owner_address",
			TypeAndVersion: "CCIP 1.0.0",
		},
		FeeQuoter: FeeQuoterView{
			ContractMetaData: ContractMetaData{
				Address:        ccipAddress,
				TypeAndVersion: "FeeQuoter 1.0.0",
			},
			FeeTokens: []string{
				"0xmocked_link_token",
				"0xmocked_sui_token",
			},
			StaticConfig: FeeQuoterStaticConfig{
				MaxFeeJuelsPerMsg:            "1000000000000000000",
				LinkToken:                    "0xmocked_link_token",
				TokenPriceStalenessThreshold: 3600,
			},
			DestinationChainConfigs: map[uint64]FeeQuoterDestChainConfig{
				1: {
					IsEnabled:                         true,
					MaxNumberOfTokensPerMsg:           10,
					MaxDataBytes:                      10000,
					MaxPerMsgGasLimit:                 4000000,
					DestGasOverhead:                   100000,
					DestGasPerPayloadByteBase:         16,
					DestGasPerPayloadByteHigh:         32,
					DestGasPerPayloadByteThreshold:    1024,
					DestDataAvailabilityOverheadGas:   50000,
					DestGasPerDataAvailabilityByte:    100,
					DestDataAvailabilityMultiplierBps: 10000,
					ChainFamilySelector:               "evm",
					EnforceOutOfOrder:                 false,
					DefaultTokenFeeUsdCents:           50,
					DefaultTokenDestGasOverhead:       34000,
					DefaultTxGasLimit:                 200000,
					GasMultiplierWeiPerEth:            1000000000000000000,
					GasPriceStalenessThreshold:        3600,
					NetworkFeeUsdCents:                100,
				},
			},
		},
		RMNRemote: RMNRemoteView{
			ContractMetaData: ContractMetaData{
				Address:        ccipAddress,
				TypeAndVersion: "RMNRemote 1.0.0",
			},
			IsCursed: false,
			Config: RMNRemoteVersionedConfig{
				Version: 1,
				Signers: []RMNRemoteSigner{
					{
						OnchainPublicKey: "0xmocked_signer_1",
						NodeIndex:        0,
					},
					{
						OnchainPublicKey: "0xmocked_signer_2",
						NodeIndex:        1,
					},
				},
				Fsign: 2,
			},
			CursedSubjectEntries: []RMNRemoteCurseEntry{},
		},
		TokenAdminRegistry: TokenAdminRegistryView{
			ContractMetaData: ContractMetaData{
				Address:        ccipAddress,
				TypeAndVersion: "TokenAdminRegistry 1.0.0",
			},
		},
		NonceManager: NonceManagerView{
			ContractMetaData: ContractMetaData{
				Address:        ccipAddress,
				TypeAndVersion: "NonceManager 1.0.0",
			},
		},
		ReceiverRegistry: ReceiverRegistryView{
			ContractMetaData: ContractMetaData{
				Address:        ccipAddress,
				TypeAndVersion: "ReceiverRegistry 1.0.0",
			},
		},
	}, nil
}
