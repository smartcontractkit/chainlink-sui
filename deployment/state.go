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
	LockReleaseStateID         string
	FeeQuoterCapID             string
	OnRampAddress              string
	OnRampStateObjectID        string
	OffRampAddress             string
	OffRampOwnerCapID          string
	OffRampStateObjectID       string
	LinkTokenAddress           string
	LinkTokenCoinMetadataID    string
	LinkTokenTreasuryCapID     string
	CCIPBurnMintTokenPool      string
	CCIPBurnMintTokenPoolState string
}

// LoadOnchainStatesui loads chain state for sui chains from env
func LoadOnchainStatesui(env cldf.Environment) (map[uint64]CCIPChainState, error) {
	rawChains := env.BlockChains.SuiChains()
	suiChains := make(map[uint64]CCIPChainState)

	for chainSelector := range rawChains {
		addresses, err := env.ExistingAddresses.AddressesForChain(chainSelector) //nolint:staticcheck // we need to migrate to datastore
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
			chainState.LockReleaseStateID = addr

		case SuiMCMSType:
			chainState.MCMsAddress = addr

		case SuiTokenPoolType:
			chainState.TokenPoolAddress = addr

		case SuiCCIPObjectRefType:
			chainState.CCIPObjectRef = addr

		case SuiFeeQuoterCapType:
			chainState.FeeQuoterCapID = addr

		case SuiOnRampType:
			chainState.OnRampAddress = addr

		case SuiOnRampStateObjectIDType:
			chainState.OnRampStateObjectID = addr

		case SuiOffRampType:
			chainState.OffRampAddress = addr

		case SuiOffRampStateObjectIDType:
			chainState.OffRampStateObjectID = addr

		case SuiOffRampOwnerCapObjectIDType:
			chainState.OffRampOwnerCapID = addr

		case SuiLinkTokenType:
			chainState.LinkTokenAddress = addr

		case SuiLinkTokenObjectMetadataID:
			chainState.LinkTokenCoinMetadataID = addr

		case SuiLinkTokenTreasuryCapID:
			chainState.LinkTokenTreasuryCapID = addr

		case SuiBnMTokenPoolType:
			chainState.CCIPBurnMintTokenPool = addr

		case SuiBnMTokenPoolStateType:
			chainState.CCIPBurnMintTokenPoolState = addr
		}
		// Set address based on type
	}
	return chainState, nil
}
