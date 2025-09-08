package changesets

import (
	"math/big"

	ccip_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip"
)

type SuiChainDefinition struct {
	// Selector is the chain selector of this chain.
	Selector uint64 `json:"selector"`
	// GasPrice defines the USD price (18 decimals) per unit gas for this chain as a destination.
	GasPrice *big.Int `json:"gasPrice"`
}

type DeploySuiChainConfig struct {
	ContractParamsPerChain map[uint64]ChainContractParams
}

type MintSuiTokenConfig struct {
	ChainSelector  uint64
	TokenPackageId string
	TreasuryCapId  string
	Amount         uint64
}

type DeployDummyRecieverConfig struct {
	ChainSelector uint64
	CCIPPackageId string
	McmsPackageId string
	McmsOwner     string
}

type UpdateSuiPriceConfig struct {
	ChainSelector       uint64
	DestChainSelector   []uint64
	CCIPPackageId       string
	CCIPObjectRef       string
	SourceTokenMetadata []string
	SourceUsdPerToken   []*big.Int
	GasUsdPerUnitGas    []*big.Int
}

type ApplyFeeTokenUpdateConfig struct {
	ChainSelector     uint64
	CCIPPackageId     string
	StateObjectId     string
	OwnerCapObjectId  string
	FeeTokensToRemove []string
	FeeTokensToAdd    []string
}

type ApplyPremiumMultiplierWeiPerEthConfig struct {
	ChainSelector                 uint64
	CCIPPackageId                 string
	CCIPObjectRef                 string
	StateObjectId                 string
	OwnerCapObjectId              string
	Tokens                        []string
	PremiumMultiplierWeiPerEth    []uint64
	TokenTransferFeeConfigUpdates []TokenTransferFeeConfig
}

type TokenTransferFeeConfig struct {
	TokenMetadata                 string
	MinFeeUsdCents                uint32
	MaxFeeUsdCents                uint32
	DeciBps                       uint16
	DestGasOverhead               uint32
	DestBytesOverhead             uint32
	IsEnabled                     bool
	AggregateRateLimitCapacity    uint64
	AggregateRateLimitRate        uint64
	AggregateRateLimitIsEnabled   bool
	TokenBucketRateLimitCapacity  uint64
	TokenBucketRateLimitRate      uint64
	TokenBucketRateLimitIsEnabled bool
}

// ChainContractParams stores configuration to call initialize in CCIP contracts
type ChainContractParams struct {
	DestChainSelector uint64
	FeeQuoterParams   ccip_ops.InitFeeQuoterInput
	OffRampParams     OffRampParams
	OnRampParams      OnRampParams
}

type OffRampParams struct {
	ChainSelector                    uint64
	PermissionlessExecutionThreshold uint32
	IsRMNVerificationDisabled        []bool
	SourceChainSelectors             []uint64
	SourceChainIsEnabled             []bool
	SourceChainsOnRamp               [][]byte
}

type OnRampParams struct {
	ChainSelector  uint64
	AllowlistAdmin string
	FeeAggregator  string
}

type UpdateSuiLaneConfig struct {
	ChainSelector           uint64
	RemoteChainSelector     uint64
	CCIPPackageId           string
	OnRampStateObjectId     string
	OffRampStateObjectId    string
	OffRampOwnerCap         string
	RemoteOnRampBytes       []byte
	MaxNumberOfTokensPerMsg uint16
	MaxDataBytes            uint32
	MaxPerMsgGasLimit       uint64
}

type SetOCR3OffRampConfig struct {
	ChainSelector         uint64
	CCIPPackageId         string
	OffRampStateObjectId  string
	OffRampOwnerCap       string
	OffchainConfigVersion uint64
	Signers               []string
	Transmitters          []string
	F                     byte
	OnchainConfig         []byte
	OffchainConfigBytes   []byte
}
