package deployment

import (
	"errors"
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink-sui/deployment/view"
)

type SuiChainView struct {
	ChainSelector uint64 `json:"chainSelector,omitempty"`
	ChainID       string `json:"chainID,omitempty"`

	MCMSWithTimelock string // TODO

	CCIP    string          // TODO
	OnRamp  string          // TODO
	OffRamp string          // TODO
	Router  view.RouterView `json:"router,omitempty"`

	Tokens     map[string]string            `json:"tokens,omitempty"`
	TokenPools map[string]map[string]string // TODO // TokenSymbol => TokenPool Address => PoolView
}

type CCIPPoolState struct {
	PackageID        string
	StateObjectId    string
	OwnerCapObjectId string
}

type CCIPChainState struct {
	// MCMS related
	MCMSPackageID               string
	MCMSStateObjectID           string
	MCMSRegistryObjectID        string
	MCMSDeployerStateObjectID   string
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
	LnRTokenPools     map[string]CCIPPoolState
	BnMTokenPools     map[string]CCIPPoolState
	ManagedTokenPools map[string]CCIPPoolState
}

func (s CCIPChainState) GenerateView(e *cldf.Environment, selector uint64, chainName string) (SuiChainView, error) {
	lggr := e.Logger
	chainView := SuiChainView{
		ChainSelector: selector,
	}

	lggr.Infow("generating Sui chain view", "chain", chainName, "selector", selector)

	// Router
	if s.CCIPRouterAddress != "" {
		routerView, err := view.GenerateRouterView(s.CCIPRouterAddress, []string{s.OffRampAddress}, false)
		if err != nil {
			return SuiChainView{}, fmt.Errorf("failed to generate router view for router %s: %w", s.CCIPRouterAddress, err)
		}
		chainView.Router = routerView
		lggr.Infow("generated router view", "routerAddress", s.CCIPRouterAddress, "chain", chainName)
	}

	return chainView, nil
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
	chainState := CCIPChainState{
		BnMTokenPools:     make(map[string]CCIPPoolState),
		LnRTokenPools:     make(map[string]CCIPPoolState),
		ManagedTokenPools: make(map[string]CCIPPoolState),
	}
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
		case SuiMcmsDeployerObjectIDType:
			chainState.MCMSDeployerStateObjectID = addr

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

		// BnM Token pools related
		case SuiBnMTokenPoolType:
			symbol, err := getTokenSymbol(typeAndVersion)
			if err != nil {
				return CCIPChainState{}, fmt.Errorf("failed to get token symbol for BnM token pool: %w", err)
			}
			pool := chainState.BnMTokenPools[symbol]
			pool.PackageID = addr
			chainState.BnMTokenPools[symbol] = pool
		case SuiBnMTokenPoolStateType:
			symbol, err := getTokenSymbol(typeAndVersion)
			if err != nil {
				return CCIPChainState{}, fmt.Errorf("failed to get token symbol for BnM token pool: %w", err)
			}
			pool := chainState.BnMTokenPools[symbol]
			pool.StateObjectId = addr
			chainState.BnMTokenPools[symbol] = pool
		case SuiBnMTokenPoolOwnerIDType:
			symbol, err := getTokenSymbol(typeAndVersion)
			if err != nil {
				return CCIPChainState{}, fmt.Errorf("failed to get token symbol for BnM token pool: %w", err)
			}
			pool := chainState.BnMTokenPools[symbol]
			pool.OwnerCapObjectId = addr
			chainState.BnMTokenPools[symbol] = pool

		//  LnR Token pools related
		case SuiLnRTokenPoolType:
			symbol, err := getTokenSymbol(typeAndVersion)
			if err != nil {
				return CCIPChainState{}, fmt.Errorf("failed to get token symbol for LnR token pool: %w", err)
			}
			pool := chainState.LnRTokenPools[symbol]
			pool.PackageID = addr
			chainState.LnRTokenPools[symbol] = pool
		case SuiLnRTokenPoolStateType:
			symbol, err := getTokenSymbol(typeAndVersion)
			if err != nil {
				return CCIPChainState{}, fmt.Errorf("failed to get token symbol for LnR token pool: %w", err)
			}
			pool := chainState.LnRTokenPools[symbol]
			pool.StateObjectId = addr
			chainState.LnRTokenPools[symbol] = pool
		case SuiLnRTokenPoolOwnerIDType:
			symbol, err := getTokenSymbol(typeAndVersion)
			if err != nil {
				return CCIPChainState{}, fmt.Errorf("failed to get token symbol for LnR token pool: %w", err)
			}
			pool := chainState.LnRTokenPools[symbol]
			pool.OwnerCapObjectId = addr
			chainState.LnRTokenPools[symbol] = pool

		// Managed Token pools related
		case SuiManagedTokenPoolType:
			symbol, err := getTokenSymbol(typeAndVersion)
			if err != nil {
				return CCIPChainState{}, fmt.Errorf("failed to get token symbol for Managed token pool: %w", err)
			}
			pool := chainState.ManagedTokenPools[symbol]
			pool.PackageID = addr
			chainState.ManagedTokenPools[symbol] = pool
		case SuiManagedTokenPoolStateType:
			symbol, err := getTokenSymbol(typeAndVersion)
			if err != nil {
				return CCIPChainState{}, fmt.Errorf("failed to get token symbol for Managed token pool: %w", err)
			}
			pool := chainState.ManagedTokenPools[symbol]
			pool.StateObjectId = addr
			chainState.ManagedTokenPools[symbol] = pool
		case SuiManagedTokenPoolOwnerIDType:
			symbol, err := getTokenSymbol(typeAndVersion)
			if err != nil {
				return CCIPChainState{}, fmt.Errorf("failed to get token symbol for Managed token pool: %w", err)
			}
			pool := chainState.ManagedTokenPools[symbol]
			pool.OwnerCapObjectId = addr
			chainState.ManagedTokenPools[symbol] = pool
		}
	}
	return chainState, nil
}

func getTokenSymbol(typeAndVersion cldf.TypeAndVersion) (string, error) {
	if typeAndVersion.Labels.IsEmpty() {
		return "", fmt.Errorf("no labels found for type %s", typeAndVersion.Type)
	}
	labels := typeAndVersion.Labels.List()
	symbolStr := labels[0]
	if symbolStr == "" {
		return "", fmt.Errorf("empty symbol label for type %s", typeAndVersion.Type)
	}
	return symbolStr, nil
}
