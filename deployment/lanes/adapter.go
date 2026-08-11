package lanes

import (
	"encoding/binary"
	"math/big"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"

	laneapi "github.com/smartcontractkit/chainlink-ccip/deployment/lanes"
)

// suiFamilySelector is bytes4(keccak256("CCIP ChainFamilySelector Sui")) = 0xc4e05953.
var suiFamilySelector = [4]byte{0xc4, 0xe0, 0x59, 0x53}

// SuiAdapter implements laneapi.LaneAdapter for Sui CCIP 1.6.0 lanes.
// Package IDs are resolved via LoadOnchainStatesui (datastore address refs first,
// legacy addresses.json fallback); the ds parameter is unused.
// ConnectChains must run inside WithConnectChainsEnvironment.
type SuiAdapter struct{}

func (a *SuiAdapter) GetOnRampAddress(_ datastore.DataStore, chainSelector uint64) ([]byte, error) {
	e, err := connectChainsEnvironment()
	if err != nil {
		return nil, err
	}
	return resolveOnRampPackageID(e, chainSelector)
}

func (a *SuiAdapter) GetOffRampAddress(_ datastore.DataStore, chainSelector uint64) ([]byte, error) {
	e, err := connectChainsEnvironment()
	if err != nil {
		return nil, err
	}
	return resolveOffRampPackageID(e, chainSelector)
}

func (a *SuiAdapter) GetFQAddress(_ datastore.DataStore, chainSelector uint64) ([]byte, error) {
	e, err := connectChainsEnvironment()
	if err != nil {
		return nil, err
	}
	return resolveCCIPPackageID(e, chainSelector)
}

func (a *SuiAdapter) GetRouterAddress(_ datastore.DataStore, chainSelector uint64) ([]byte, error) {
	e, err := connectChainsEnvironment()
	if err != nil {
		return nil, err
	}
	return resolveRouterPackageID(e, chainSelector)
}

// GetFeeQuoterDestChainConfig returns defaults applied on remote FeeQuoters when Sui is
// the destination chain. Values align with wired Sui testnet lanes (evm_feequoter_dest_configure)
// and Sui execution limits (MaxDataBytes 16_000).
func (a *SuiAdapter) GetFeeQuoterDestChainConfig() laneapi.FeeQuoterDestChainConfig {
	return laneapi.FeeQuoterDestChainConfig{
		IsEnabled:                   true,
		MaxDataBytes:                16_000,
		MaxPerMsgGasLimit:           3_000_000,
		DestGasOverhead:             300_000,
		DestGasPerPayloadByteBase:   16,
		ChainFamilySelector:         binary.BigEndian.Uint32(suiFamilySelector[:]),
		DefaultTokenFeeUSDCents:     25,
		DefaultTokenDestGasOverhead: 90_000,
		DefaultTxGasLimit:           200_000,
		NetworkFeeUSDCents:          10,
		V1Params: &laneapi.FeeQuoterV1Params{
			MaxNumberOfTokensPerMsg:           10,
			DestGasPerPayloadByteHigh:         40,
			DestGasPerPayloadByteThreshold:    3_000,
			DestDataAvailabilityOverheadGas:   100,
			DestGasPerDataAvailabilityByte:    16,
			DestDataAvailabilityMultiplierBps: 1,
			EnforceOutOfOrder:                 true,
			GasMultiplierWeiPerEth:            11e17,
		},
	}
}

func (a *SuiAdapter) GetDefaultGasPrice() *big.Int {
	return big.NewInt(15e11)
}

func (a *SuiAdapter) GetChainFamilySelector() [4]byte {
	return suiFamilySelector
}
