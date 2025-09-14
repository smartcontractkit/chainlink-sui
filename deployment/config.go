package deployment

// this is config.go file
import (
	ccipops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip"
	offrampops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_offramp"
)

// These are the static/default FeeQuoter + chain params
var DefaultCCIPSeqConfig = ccipops.DeployAndInitCCIPSeqInput{
	// Initialize FeeQuoter
	MaxFeeJuelsPerMsg:            "200000000000000000000",
	TokenPriceStalenessThreshold: 1000000,

	// Fee Quoter configuration
	AddMinFeeUsdCents:    []uint32{3000},
	AddMaxFeeUsdCents:    []uint32{30000},
	AddDeciBps:           []uint16{1000},
	AddDestGasOverhead:   []uint32{1000000},
	AddDestBytesOverhead: []uint32{1000},
	AddIsEnabled:         []bool{true},
	RemoveTokens:         []string{},

	// Fee Quoter destination chain configuration
	IsEnabled:                         true,
	MaxNumberOfTokensPerMsg:           1,
	MaxDataBytes:                      30_000,
	MaxPerMsgGasLimit:                 3_000_000,
	DestGasOverhead:                   300_000,
	DestGasPerPayloadByteBase:         byte(16),
	DestGasPerPayloadByteHigh:         byte(40),
	DestGasPerPayloadByteThreshold:    uint16(3000),
	DestDataAvailabilityOverheadGas:   100,
	DestGasPerDataAvailabilityByte:    16,
	DestDataAvailabilityMultiplierBps: 1,
	ChainFamilySelector:               []byte{40, 18, 213, 44},
	EnforceOutOfOrder:                 false,
	DefaultTokenFeeUsdCents:           25,
	DefaultTokenDestGasOverhead:       90_000,
	DefaultTxGasLimit:                 200_000,
	GasMultiplierWeiPerEth:            1_000_000_000_000_000_000,
	GasPriceStalenessThreshold:        1_000_000,
	NetworkFeeUsdCents:                10,

	// apply_premium_multiplier_wei_per_eth_updates
	PremiumMultiplierWeiPerEth: []uint64{900_000_000_000_000_000},
}

var DefaultOffRampSeqConfig = offrampops.DeployAndInitCCIPOffRampSeqInput{
	InitializeOffRampInput: offrampops.InitializeOffRampInput{
		PremissionExecThresholdSeconds:        uint32(60 * 60 * 8), // 8 hours
		SourceChainsIsEnabled:                 []bool{true},
		SourceChainsIsRMNVerificationDisabled: []bool{true},
	},
}
