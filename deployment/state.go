package deployment

import (
	"errors"
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

type CCIPChainState struct {
	CCIPRouterAddress            string
	CCIPAddress                  string
	CCIPObjectRef                string
	CCIPOwnerCapObjectId         string
	CCIPUpgradeCapObjectId       string
	MCMsAddress                  string
	TokenPoolAddress             string
	LockReleaseAddress           string
	LockReleaseStateId           string
	FeeQuoterCapId               string
	OnRampAddress                string
	OnRampStateObjectId          string
	OnRampOwnerCapObjectId       string
	OnRampUpgradeCapId           string
	OffRampAddress               string
	OffRampOwnerCapId            string
	OffRampUpgradeCapId          string
	OffRampStateObjectId         string
	LinkTokenAddress             string
	LinkTokenCoinMetadataId      string
	LinkTokenTreasuryCapId       string
	CCIPBurnMintTokenPool        string
	CCIPBurnMintTokenPoolState   string
	CCIPBurnMintTokenPoolOwnerId string
}

// LoadOnchainStatesui loads chain state for sui chains from env
func LoadOnchainStatesui(env cldf.Environment) (map[uint64]CCIPChainState, error) {
	rawChains := env.BlockChains.SuiChains()
	suiChains := make(map[uint64]CCIPChainState)

	for chainSelector := range rawChains {
		addresses, err := env.ExistingAddresses.AddressesForChain(chainSelector)
		if err != nil {
			// Chain not found in address book, initialize empty state
			if !errors.Is(err, cldf.ErrChainNotFound) {
				return nil, fmt.Errorf("failed to get addresses for chain %d: %w", chainSelector, err)
			}
			addresses = make(map[string]cldf.TypeAndVersion)
		}

		chainState, err := loadsuiChainStateFromAddresses(addresses)
		if err != nil {
			return nil, fmt.Errorf("failed to load chain state for chain %d: %w", chainSelector, err)
		}

		suiChains[chainSelector] = chainState
	}

	return suiChains, nil
}

func loadsuiChainStateFromAddresses(addresses map[string]cldf.TypeAndVersion) (CCIPChainState, error) {
	chainState := CCIPChainState{}
	for addr, typeAndVersion := range addresses {
		// Parse addresss

		switch typeAndVersion.Type {
		case SuiCCIPRouterType:
			chainState.CCIPRouterAddress = addr

		case SuiCCIPType:
			chainState.CCIPAddress = addr

		case SuiLockReleaseTPType:
			chainState.LockReleaseAddress = addr

		case SuiLockReleaseTPStateType:
			chainState.LockReleaseStateId = addr

		case SuiMcmsPackageIDType:
			chainState.MCMsAddress = addr

		case SuiCCIPObjectRefType:
			chainState.CCIPObjectRef = addr

		case SuiCCIPOwnerCapObjectIDType:
			chainState.CCIPOwnerCapObjectId = addr

		case SuiCCIPUpgradeCapObjectIDType:
			chainState.CCIPUpgradeCapObjectId = addr

		case SuiFeeQuoterCapType:
			chainState.FeeQuoterCapId = addr

		case SuiOnRampType:
			chainState.OnRampAddress = addr

		case SuiOnRampStateObjectIDType:
			chainState.OnRampStateObjectId = addr

		case SuiOnRampOwnerCapObjectIDType:
			chainState.OnRampOwnerCapObjectId = addr

		case SuiOnRampUpgradeCapObjectIDType:
			chainState.OnRampUpgradeCapId = addr

		case SuiOffRampType:
			chainState.OffRampAddress = addr

		case SuiOffRampStateObjectIDType:
			chainState.OffRampStateObjectId = addr

		case SuiOffRampOwnerCapObjectIDType:
			chainState.OffRampOwnerCapId = addr

		case SuiOffRampUpgradeCapObjectIDType:
			chainState.OffRampUpgradeCapId = addr

		case SuiLinkTokenType:
			chainState.LinkTokenAddress = addr

		case SuiLinkTokenObjectMetadataID:
			chainState.LinkTokenCoinMetadataId = addr

		case SuiLinkTokenTreasuryCapID:
			chainState.LinkTokenTreasuryCapId = addr

		case SuiBnMTokenPoolType:
			chainState.CCIPBurnMintTokenPool = addr

		case SuiBnMTokenPoolStateType:
			chainState.CCIPBurnMintTokenPoolState = addr

		case SuiBnMTokenPoolOwnerIDType:
			chainState.CCIPBurnMintTokenPoolOwnerId = addr
		}
		// Set address based on type
	}
	return chainState, nil
}
