package deployment

import (
	"errors"
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

type CCIPChainState struct {
	CCIPRouterAddress          string
	CCIPAddress                string
	CCIPObjectRef              string
	MCMsAddress                string
	TokenPoolAddress           string
	LockReleaseAddress         string
	LockReleaseStateId         string
	FeeQuoterCapId             string
	OnRampAddress              string
	OnRampStateObjectId        string
	OffRampAddress             string
	OffRampOwnerCapId          string
	OffRampStateObjectId       string
	LinkTokenAddress           string
	LinkTokenCoinMetadataId    string
	LinkTokenTreasuryCapId     string
	CCIPBurnMintTokenPool      string
	CCIPBurnMintTokenPoolState string
}

// LoadOnchainStatesui loads chain state for sui chains from env
func LoadOnchainStatesui(env cldf.Environment) (map[uint64]CCIPChainState, error) {
	rawChains := env.BlockChains.SuiChains()
	suiChains := make(map[uint64]CCIPChainState)

	for chainSelector := range rawChains {
		addresses, err := env.ExistingAddresses.AddressesForChain(chainSelector) //nolint
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

		case SuiMCMSType:
			chainState.MCMsAddress = addr

		case SuiTokenPoolType:
			chainState.TokenPoolAddress = addr

		case SuiCCIPObjectRefType:
			chainState.CCIPObjectRef = addr

		case SuiFeeQuoterCapType:
			chainState.FeeQuoterCapId = addr

		case SuiOnRampType:
			chainState.OnRampAddress = addr

		case SuiOnRampStateObjectIdType:
			chainState.OnRampStateObjectId = addr

		case SuiOffRampType:
			chainState.OffRampAddress = addr

		case SuiOffRampStateObjectIdType:
			chainState.OffRampStateObjectId = addr

		case SuiOffRampOwnerCapObjectIdType:
			chainState.OffRampOwnerCapId = addr

		case SuiLinkTokenType:
			chainState.LinkTokenAddress = addr

		case SuiLinkTokenObjectMetadataId:
			chainState.LinkTokenCoinMetadataId = addr

		case SuiLinkTokenTreasuryCapId:
			chainState.LinkTokenTreasuryCapId = addr

		case SuiBnMTokenPoolType:
			chainState.CCIPBurnMintTokenPool = addr

		case SuiBnMTokenPoolStateType:
			chainState.CCIPBurnMintTokenPoolState = addr
		}
		// Set address based on type

	}
	return chainState, nil
}
