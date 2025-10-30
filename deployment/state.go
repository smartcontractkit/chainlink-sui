package deployment

import (
	"errors"
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

type CCIPChainState struct {
	// MCMS related
	MCMSPackageID               string
	MCMSStateObjectID           string
	MCMSRegistryObjectID        string
	MCMSAccountStateObjectID    string
	MCMSAccountOwnerCapObjectID string
	MCMSTimelockObjectID        string

	// CCIP related
	CCIPAddress            string
	CCIPObjectRef          string
	CCIPOwnerCapObjectId   string
	CCIPUpgradeCapObjectId string
	FeeQuoterCapId         string

	// CCIP Router related
	CCIPRouterAddress       string
	CCIPRouterStateObjectID string

	// OnRamp related
	OnRampAddress          string
	OnRampStateObjectId    string
	OnRampOwnerCapObjectId string
	OnRampUpgradeCapId     string

	// OffRamp related
	OffRampAddress       string
	OffRampStateObjectId string
	OffRampOwnerCapId    string
	OffRampUpgradeCapId  string

	// LINK token related
	LinkTokenAddress        string
	LinkTokenCoinMetadataId string
	LinkTokenTreasuryCapId  string

	// Token related
	ManagedToken                 string
	ManagedTokenOwnerCapObjectID string
	ManagedTokenStateObjectID    string
	ManagedTokenMinterCapID      string

	// Token pools related
	LockReleaseAddress           string
	LockReleaseStateId           string
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
		// Parse addresss based on type and label
		switch typeAndVersion.Type {

		// MCMS related
		case SuiMcmsPackageIDType:
			chainState.MCMSPackageID = addr
		case SuiMcmsRegistryObjectIDType:
			chainState.MCMSRegistryObjectID = addr
		case SuiMcmsObjectIDType:
			chainState.MCMSStateObjectID = addr
		case SuiMcmsAccountStateObjectIDType:
			chainState.MCMSAccountStateObjectID = addr
		case SuiMcmsAccountOwnerCapObjectIDType:
			chainState.MCMSAccountOwnerCapObjectID = addr
		case SuiMcmsTimelockObjectIDType:
			chainState.MCMSTimelockObjectID = addr

		// CCIP Router related
		case SuiCCIPRouterType:
			chainState.CCIPRouterAddress = addr
		case SuiCCIPRouterStateObjectType:
			chainState.CCIPRouterStateObjectID = addr

		// CCIP related
		case SuiCCIPType:
			chainState.CCIPAddress = addr
		case SuiCCIPObjectRefType:
			chainState.CCIPObjectRef = addr
		case SuiCCIPOwnerCapObjectIDType:
			chainState.CCIPOwnerCapObjectId = addr
		case SuiCCIPUpgradeCapObjectIDType:
			chainState.CCIPUpgradeCapObjectId = addr
		case SuiFeeQuoterCapType:
			chainState.FeeQuoterCapId = addr

		// OnRamp related
		case SuiOnRampType:
			chainState.OnRampAddress = addr
		case SuiOnRampStateObjectIDType:
			chainState.OnRampStateObjectId = addr
		case SuiOnRampOwnerCapObjectIDType:
			chainState.OnRampOwnerCapObjectId = addr
		case SuiOnRampUpgradeCapObjectIDType:
			chainState.OnRampUpgradeCapId = addr

		// OffRamp related
		case SuiOffRampType:
			chainState.OffRampAddress = addr
		case SuiOffRampStateObjectIDType:
			chainState.OffRampStateObjectId = addr
		case SuiOffRampOwnerCapObjectIDType:
			chainState.OffRampOwnerCapId = addr
		case SuiOffRampUpgradeCapObjectIDType:
			chainState.OffRampUpgradeCapId = addr

		// LINK Token related
		case SuiLinkTokenType:
			chainState.LinkTokenAddress = addr
		case SuiLinkTokenObjectMetadataID:
			chainState.LinkTokenCoinMetadataId = addr
		case SuiLinkTokenTreasuryCapID:
			chainState.LinkTokenTreasuryCapId = addr

		// Token related
		case SuiManagedTokenType:
			chainState.ManagedToken = addr
		case SuiManagedTokenOwnerCapObjectID:
			chainState.ManagedTokenOwnerCapObjectID = addr
		case SuiManagedTokenStateObjectID:
			chainState.ManagedTokenStateObjectID = addr
		case SuiManagedTokenMinterCapID:
			chainState.ManagedTokenMinterCapID = addr

		// Token pools related
		case SuiLockReleaseTPType:
			chainState.LockReleaseAddress = addr
		case SuiLockReleaseTPStateType:
			chainState.LockReleaseStateId = addr
		case SuiBnMTokenPoolType:
			chainState.CCIPBurnMintTokenPool = addr
		case SuiBnMTokenPoolStateType:
			chainState.CCIPBurnMintTokenPoolState = addr
		case SuiBnMTokenPoolOwnerIDType:
			chainState.CCIPBurnMintTokenPoolOwnerId = addr
		}
	}
	return chainState, nil
}
