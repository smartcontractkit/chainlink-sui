package deployment

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	module_burn_mint_token_pool "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_token_pools/burn_mint_token_pool"
	module_lock_release_token_pool "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_token_pools/lock_release_token_pool"
	module_managed_token_pool "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_token_pools/managed_token_pool"
	"github.com/smartcontractkit/chainlink-sui/deployment/view"
)

type TokenPoolType string

const (
	TokenPoolTypeBurnMint    TokenPoolType = "bnm"
	TokenPoolTypeLockRelease TokenPoolType = "lnr"
	TokenPoolTypeManaged     TokenPoolType = "managed"
)

type SuiChainView struct {
	ChainSelector uint64 `json:"chainSelector,omitempty"`
	ChainID       string `json:"chainID,omitempty"`

	MCMSWithTimelock          view.MCMSWithTimelockView `json:"mcmsWithTimelock"`
	FastCurseMCMSWithTimelock view.MCMSWithTimelockView `json:"fastcurseMcmsWithTimelock,omitempty"`

	CCIP    view.CCIPView               `json:"ccip,omitempty"`
	OnRamp  map[string]view.OnRampView  `json:"onRamp,omitempty"`
	OffRamp map[string]view.OffRampView `json:"offRamp,omitempty"`
	Router  view.RouterView             `json:"router,omitempty"`

	TokenPools map[string]map[string]view.TokenPoolView `json:"tokenPools,omitempty"` // TokenSymbol => TokenPool Address => PoolView
}

type CCIPPoolState struct {
	PackageID        string
	StateObjectId    string
	OwnerCapObjectId string
	RebalancerCapIds []string // only applicable for LR TP
}

type ManagedTokenState struct {
	TokenPackageID      string
	TokenCoinMetadataID string
	TokenTreasuryCapID  string
	TokenUpgradeCapID   string
	PackageID           string
	StateObjectId       string
	OwnerCapObjectId    string
	MinterCapObjectIds  []string
	PublisherObjectId   string
}

type ManagedTokenFaucetState struct {
	PackageID          string
	StateObjectId      string
	UpgradeCapObjectId string
}

// MCMSLabel is the address book label applied to the fastcurse MCMS instance.
const MCMSFastCurseLabel = "fastcurse"

// MCMSStateFields holds the seven object IDs that describe one MCMS deployment.
type MCMSStateFields struct {
	PackageID               string
	StateObjectID           string
	RegistryObjectID        string
	DeployerStateObjectID   string
	AccountStateObjectID    string
	AccountOwnerCapObjectID string
	TimelockObjectID        string
}

type CCIPChainState struct {
	// MCMS related (normal governance instance)
	MCMSPackageID               string
	MCMSStateObjectID           string
	MCMSRegistryObjectID        string
	MCMSDeployerStateObjectID   string
	MCMSAccountStateObjectID    string
	MCMSAccountOwnerCapObjectID string
	MCMSTimelockObjectID        string

	// FastCurse MCMS related (fastcurse governance instance, stored with label "fastcurse")
	FastCurseMCMSPackageID               string
	FastCurseMCMSStateObjectID           string
	FastCurseMCMSRegistryObjectID        string
	FastCurseMCMSDeployerStateObjectID   string
	FastCurseMCMSAccountStateObjectID    string
	FastCurseMCMSAccountOwnerCapObjectID string
	FastCurseMCMSTimelockObjectID        string

	// CCIP related
	CCIPAddress            string
	CCIPObjectRef          string
	CCIPOwnerCapObjectId   string
	CCIPUpgradeCapObjectId string
	FeeQuoterCapId         string

	// CCIP Router related
	CCIPRouterAddress          string
	CCIPRouterStateObjectID    string
	CCIPRouterOwnerCapObjectId string
	CCIPRouterUpgradeCapId     string

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
	LinkTokenUpgradeCapId   string

	// Managed Token related
	ManagedTokens       map[string]ManagedTokenState
	ManagedTokenFaucets map[string]ManagedTokenFaucetState

	// Token pools related
	LnRTokenPools     map[string]CCIPPoolState
	BnMTokenPools     map[string]CCIPPoolState
	ManagedTokenPools map[string]CCIPPoolState

	// mock upgrade related
	OnRampMockV2PackageId  string
	OffRampMockV2PackageId string
	CCIPMockV2PackageId    string
}

// MCMSState returns the MCMS object IDs for the requested instance.
// When isFastCurse is true the fastcurse instance fields are returned;
// otherwise the normal governance instance fields are returned.
func (s CCIPChainState) MCMSState(isFastCurse bool) MCMSStateFields {
	if isFastCurse {
		return MCMSStateFields{
			PackageID:               s.FastCurseMCMSPackageID,
			StateObjectID:           s.FastCurseMCMSStateObjectID,
			RegistryObjectID:        s.FastCurseMCMSRegistryObjectID,
			DeployerStateObjectID:   s.FastCurseMCMSDeployerStateObjectID,
			AccountStateObjectID:    s.FastCurseMCMSAccountStateObjectID,
			AccountOwnerCapObjectID: s.FastCurseMCMSAccountOwnerCapObjectID,
			TimelockObjectID:        s.FastCurseMCMSTimelockObjectID,
		}
	}
	return MCMSStateFields{
		PackageID:               s.MCMSPackageID,
		StateObjectID:           s.MCMSStateObjectID,
		RegistryObjectID:        s.MCMSRegistryObjectID,
		DeployerStateObjectID:   s.MCMSDeployerStateObjectID,
		AccountStateObjectID:    s.MCMSAccountStateObjectID,
		AccountOwnerCapObjectID: s.MCMSAccountOwnerCapObjectID,
		TimelockObjectID:        s.MCMSTimelockObjectID,
	}
}

func (s CCIPChainState) GenerateView(e *cldf.Environment, selector uint64, chainName string) (SuiChainView, error) {
	lggr := e.Logger
	chainView := SuiChainView{
		ChainSelector: selector,
		TokenPools:    make(map[string]map[string]view.TokenPoolView),
		OnRamp:        make(map[string]view.OnRampView),
		OffRamp:       make(map[string]view.OffRampView),
	}

	lggr.Infow("generating Sui chain view", "chain", chainName, "selector", selector)

	suiChain := e.BlockChains.SuiChains()[selector]
	ctx := context.Background()

	var mu sync.Mutex
	g, ctxG1 := errgroup.WithContext(ctx)

	// Normal MCMS
	if s.MCMSStateObjectID != "" {
		g.Go(func() error {
			mcmsView, err := view.GenerateMCMSWithTimelockView(ctxG1, suiChain, s.MCMSPackageID, s.MCMSStateObjectID, s.MCMSTimelockObjectID, s.MCMSAccountStateObjectID)
			if err != nil {
				return fmt.Errorf("failed to generate mcms view for mcms %s: %w", s.MCMSStateObjectID, err)
			}
			mu.Lock()
			chainView.MCMSWithTimelock = mcmsView
			mu.Unlock()
			lggr.Infow("generated MCMS view", "mcmsStateObjectID", s.MCMSStateObjectID, "chain", chainName)
			return nil
		})
	}

	// FastCurse MCMS
	if s.FastCurseMCMSStateObjectID != "" {
		g.Go(func() error {
			mcmsView, err := view.GenerateMCMSWithTimelockView(ctxG1, suiChain, s.FastCurseMCMSPackageID, s.FastCurseMCMSStateObjectID, s.FastCurseMCMSTimelockObjectID, s.FastCurseMCMSAccountStateObjectID)
			if err != nil {
				return fmt.Errorf("failed to generate fastcurse mcms view for mcms %s: %w", s.FastCurseMCMSStateObjectID, err)
			}
			mu.Lock()
			chainView.FastCurseMCMSWithTimelock = mcmsView
			mu.Unlock()
			lggr.Infow("generated FastCurse MCMS view", "fastCurseMCMSStateObjectID", s.FastCurseMCMSStateObjectID, "chain", chainName)
			return nil
		})
	}

	// CCIP
	if s.CCIPAddress != "" {
		g.Go(func() error {
			ccipView, err := view.GenerateCCIPView(ctxG1, suiChain, s.CCIPAddress, s.CCIPObjectRef, s.CCIPRouterAddress, s.CCIPRouterStateObjectID)
			if err != nil {
				return fmt.Errorf("failed to generate ccip view for ccip %s: %w", s.CCIPAddress, err)
			}
			mu.Lock()
			chainView.CCIP = ccipView
			mu.Unlock()
			lggr.Infow("generated CCIP view", "ccipAddress", s.CCIPAddress, "chain", chainName)
			return nil
		})
	}

	// Router
	if s.CCIPRouterAddress != "" && s.CCIPRouterStateObjectID != "" {
		g.Go(func() error {
			routerView, err := view.GenerateRouterView(ctxG1, suiChain, s.CCIPRouterAddress, s.CCIPRouterStateObjectID)
			if err != nil {
				return fmt.Errorf("failed to generate router view for router %s: %w", s.CCIPRouterAddress, err)
			}
			mu.Lock()
			chainView.Router = routerView
			mu.Unlock()
			lggr.Infow("generated router view", "routerAddress", s.CCIPRouterAddress, "chain", chainName)
			return nil
		})
	}

	// OnRamp
	if s.OnRampAddress != "" {
		g.Go(func() error {
			onRampView, err := view.GenerateOnRampView(ctxG1, suiChain, s.OnRampAddress, s.OnRampStateObjectId, s.CCIPRouterAddress, s.CCIPRouterStateObjectID)
			if err != nil {
				return fmt.Errorf("failed to generate onramp view for onramp %s: %w", s.OnRampAddress, err)
			}
			mu.Lock()
			chainView.OnRamp[s.OnRampAddress] = onRampView
			mu.Unlock()
			lggr.Infow("generated onRamp view", "onRampAddress", s.OnRampAddress, "chain", chainName)
			return nil
		})
	}

	// OffRamp
	if s.OffRampAddress != "" {
		g.Go(func() error {
			offRampView, err := view.GenerateOffRampView(ctxG1, suiChain, s.OffRampAddress, s.OffRampStateObjectId, s.CCIPObjectRef)
			if err != nil {
				return fmt.Errorf("failed to generate offramp view for offramp %s: %w", s.OffRampAddress, err)
			}
			mu.Lock()
			chainView.OffRamp[s.OffRampAddress] = offRampView
			mu.Unlock()
			lggr.Infow("generated offRamp view", "offRampAddress", s.OffRampAddress, "chain", chainName)
			return nil
		})
	}

	// Wait here because pools depend on tokenAdminRegistry from CCIP view
	if err := g.Wait(); err != nil {
		return SuiChainView{}, err
	}

	// Token pools
	tokenConfigs := chainView.CCIP.TokenAdminRegistry.TokenConfigs
	g, ctxG2 := errgroup.WithContext(ctx)

	// BurnMint Token Pools
	for symbol, pool := range s.BnMTokenPools {
		if pool.PackageID == "" || pool.StateObjectId == "" {
			lggr.Warnw("Skipping BnM token pool with missing data", "symbol", symbol, "chain", chainName)
			continue
		}

		g.Go(func() error {
			contract, err := module_burn_mint_token_pool.NewBurnMintTokenPool(pool.PackageID, suiChain.Client)
			if err != nil {
				return fmt.Errorf("failed to create BnM token pool contract for symbol %s: %w", symbol, err)
			}

			poolView, err := view.GenerateTokenPoolView(ctxG2, suiChain, pool.PackageID, pool.StateObjectId, tokenConfigs, contract.DevInspect(), lggr)
			if err != nil {
				return fmt.Errorf("failed to generate BnM token pool view for symbol %s: %w", symbol, err)
			}

			if len(pool.RebalancerCapIds) > 0 {
				return fmt.Errorf("BnM token pool %s has rebalancer cap ids, but it is not applicable", symbol)
			}

			mu.Lock()
			if chainView.TokenPools[symbol] == nil {
				chainView.TokenPools[symbol] = make(map[string]view.TokenPoolView)
			}
			chainView.TokenPools[symbol][poolView.Address] = poolView
			mu.Unlock()

			lggr.Infow("generated BnM token pool view", "symbol", symbol, "poolAddress", pool.PackageID, "chain", chainName)
			return nil
		})
	}

	// LockRelease Token Pools
	for symbol, pool := range s.LnRTokenPools {
		if pool.PackageID == "" || pool.StateObjectId == "" {
			lggr.Warnw("Skipping LnR token pool with missing data", "symbol", symbol, "chain", chainName)
			continue
		}

		g.Go(func() error {
			contract, err := module_lock_release_token_pool.NewLockReleaseTokenPool(pool.PackageID, suiChain.Client)
			if err != nil {
				return fmt.Errorf("failed to create LnR token pool contract for symbol %s: %w", symbol, err)
			}

			poolView, err := view.GenerateTokenPoolView(ctxG2, suiChain, pool.PackageID, pool.StateObjectId, tokenConfigs, contract.DevInspect(), lggr)
			if err != nil {
				return fmt.Errorf("failed to generate LnR token pool view for symbol %s: %w", symbol, err)
			}

			poolView.RebalancerCapIds = pool.RebalancerCapIds

			mu.Lock()
			if chainView.TokenPools[symbol] == nil {
				chainView.TokenPools[symbol] = make(map[string]view.TokenPoolView)
			}
			chainView.TokenPools[symbol][poolView.Address] = poolView
			mu.Unlock()

			lggr.Infow("generated LnR token pool view", "symbol", symbol, "poolAddress", pool.PackageID, "chain", chainName)
			return nil
		})
	}

	// Managed Token Pools
	for symbol, pool := range s.ManagedTokenPools {
		if pool.PackageID == "" || pool.StateObjectId == "" {
			lggr.Warnw("Skipping managed token pool with missing data", "symbol", symbol, "chain", chainName)
			continue
		}

		g.Go(func() error {
			contract, err := module_managed_token_pool.NewManagedTokenPool(pool.PackageID, suiChain.Client)
			if err != nil {
				return fmt.Errorf("failed to create managed token pool contract for symbol %s: %w", symbol, err)
			}

			poolView, err := view.GenerateTokenPoolView(ctxG2, suiChain, pool.PackageID, pool.StateObjectId, tokenConfigs, contract.DevInspect(), lggr)
			if err != nil {
				return fmt.Errorf("failed to generate managed token pool view for symbol %s: %w", symbol, err)
			}

			if len(pool.RebalancerCapIds) > 0 {
				return fmt.Errorf("managed token pool %s has rebalancer cap ids, but it is not applicable", symbol)
			}

			mu.Lock()
			if chainView.TokenPools[symbol] == nil {
				chainView.TokenPools[symbol] = make(map[string]view.TokenPoolView)
			}
			chainView.TokenPools[symbol][poolView.Address] = poolView
			mu.Unlock()

			lggr.Infow("generated managed token pool view", "symbol", symbol, "poolAddress", pool.PackageID, "chain", chainName)
			return nil
		})
	}

	return chainView, g.Wait()
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
		ManagedTokens:       make(map[string]ManagedTokenState),
		ManagedTokenFaucets: make(map[string]ManagedTokenFaucetState),
		BnMTokenPools:       make(map[string]CCIPPoolState),
		LnRTokenPools:       make(map[string]CCIPPoolState),
		ManagedTokenPools:   make(map[string]CCIPPoolState),
	}
	for addr, typeAndVersion := range addresses {
		// Determine whether this address belongs to the fastcurse MCMS instance.
		isFastCurse := typeAndVersion.Labels.Contains(MCMSFastCurseLabel)

		switch typeAndVersion.Type {

		// MCMS related — route to normal or fastcurse fields based on the label.
		case SuiMcmsPackageIDType:
			if isFastCurse {
				chainState.FastCurseMCMSPackageID = addr
			} else {
				chainState.MCMSPackageID = addr
			}
		case SuiMcmsRegistryObjectIDType:
			if isFastCurse {
				chainState.FastCurseMCMSRegistryObjectID = addr
			} else {
				chainState.MCMSRegistryObjectID = addr
			}
		case SuiMcmsObjectIDType:
			if isFastCurse {
				chainState.FastCurseMCMSStateObjectID = addr
			} else {
				chainState.MCMSStateObjectID = addr
			}
		case SuiMcmsAccountStateObjectIDType:
			if isFastCurse {
				chainState.FastCurseMCMSAccountStateObjectID = addr
			} else {
				chainState.MCMSAccountStateObjectID = addr
			}
		case SuiMcmsAccountOwnerCapObjectIDType:
			if isFastCurse {
				chainState.FastCurseMCMSAccountOwnerCapObjectID = addr
			} else {
				chainState.MCMSAccountOwnerCapObjectID = addr
			}
		case SuiMcmsTimelockObjectIDType:
			if isFastCurse {
				chainState.FastCurseMCMSTimelockObjectID = addr
			} else {
				chainState.MCMSTimelockObjectID = addr
			}
		case SuiMcmsDeployerObjectIDType:
			if isFastCurse {
				chainState.FastCurseMCMSDeployerStateObjectID = addr
			} else {
				chainState.MCMSDeployerStateObjectID = addr
			}

		// CCIP Router related
		case SuiCCIPRouterType:
			chainState.CCIPRouterAddress = addr
		case SuiCCIPRouterStateObjectType:
			chainState.CCIPRouterStateObjectID = addr
		case SuiCCIPRouterOwnerCapObjectIDType:
			chainState.CCIPRouterOwnerCapObjectId = addr
		case SuiRouterUpgradeCapObjectIDType:
			chainState.CCIPRouterUpgradeCapId = addr

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
		case SuiLinkTokenUpgradeCapID:
			chainState.LinkTokenUpgradeCapId = addr

		// Managed Token related
		case SuiManagedTokenPackageIDType:
			symbol, err := getTokenSymbol(typeAndVersion)
			if err != nil {
				return CCIPChainState{}, fmt.Errorf("failed to get token symbol for Managed token: %w", err)
			}
			managed_token := chainState.ManagedTokens[symbol]
			managed_token.PackageID = addr
			chainState.ManagedTokens[symbol] = managed_token
		case SuiManagedTokenCoinMetadataIDType:
			symbol, err := getTokenSymbol(typeAndVersion)
			if err != nil {
				return CCIPChainState{}, fmt.Errorf("failed to get token symbol for Managed token: %w", err)
			}
			managed_token := chainState.ManagedTokens[symbol]
			managed_token.TokenCoinMetadataID = addr
			chainState.ManagedTokens[symbol] = managed_token
		case SuiManagedTokenUpgradeCapIDType:
			symbol, err := getTokenSymbol(typeAndVersion)
			if err != nil {
				return CCIPChainState{}, fmt.Errorf("failed to get token symbol for Managed token: %w", err)
			}
			managed_token := chainState.ManagedTokens[symbol]
			managed_token.TokenUpgradeCapID = addr
			chainState.ManagedTokens[symbol] = managed_token
		case SuiManagedTokenTreasuryCapIDType:
			symbol, err := getTokenSymbol(typeAndVersion)
			if err != nil {
				return CCIPChainState{}, fmt.Errorf("failed to get token symbol for Managed token: %w", err)
			}
			managed_token := chainState.ManagedTokens[symbol]
			managed_token.TokenTreasuryCapID = addr
			chainState.ManagedTokens[symbol] = managed_token
		case SuiManagedTokenType:
			symbol, err := getTokenSymbol(typeAndVersion)
			if err != nil {
				return CCIPChainState{}, fmt.Errorf("failed to get token symbol for Managed token: %w", err)
			}
			managed_token := chainState.ManagedTokens[symbol]
			managed_token.TokenPackageID = addr
			chainState.ManagedTokens[symbol] = managed_token
		case SuiManagedTokenOwnerCapObjectID:
			symbol, err := getTokenSymbol(typeAndVersion)
			if err != nil {
				return CCIPChainState{}, fmt.Errorf("failed to get token symbol for Managed token: %w", err)
			}
			managed_token := chainState.ManagedTokens[symbol]
			managed_token.OwnerCapObjectId = addr
			chainState.ManagedTokens[symbol] = managed_token
		case SuiManagedTokenStateObjectID:
			symbol, err := getTokenSymbol(typeAndVersion)
			if err != nil {
				return CCIPChainState{}, fmt.Errorf("failed to get token symbol for Managed token: %w", err)
			}
			managed_token := chainState.ManagedTokens[symbol]
			managed_token.StateObjectId = addr
			chainState.ManagedTokens[symbol] = managed_token
		case SuiManagedTokenMinterCapID:
			symbol, err := getTokenSymbol(typeAndVersion)
			if err != nil {
				return CCIPChainState{}, fmt.Errorf("failed to get token symbol for Managed token: %w", err)
			}
			managed_token := chainState.ManagedTokens[symbol]
			managed_token.MinterCapObjectIds = append(managed_token.MinterCapObjectIds, addr)
			chainState.ManagedTokens[symbol] = managed_token
		case SuiManagedTokenPublisherObjectId:
			symbol, err := getTokenSymbol(typeAndVersion)
			if err != nil {
				return CCIPChainState{}, fmt.Errorf("failed to get token symbol for Managed token: %w", err)
			}
			managed_token := chainState.ManagedTokens[symbol]
			managed_token.PublisherObjectId = addr
			chainState.ManagedTokens[symbol] = managed_token
		case SuiManagedTokenFaucetPackageIDType:
			symbol, err := getTokenSymbol(typeAndVersion)
			if err != nil {
				return CCIPChainState{}, fmt.Errorf("failed to get token symbol for Managed token faucet: %w", err)
			}
			faucet := chainState.ManagedTokenFaucets[symbol]
			faucet.PackageID = addr
			chainState.ManagedTokenFaucets[symbol] = faucet
		case SuiManagedTokenFaucetStateObjectIDType:
			symbol, err := getTokenSymbol(typeAndVersion)
			if err != nil {
				return CCIPChainState{}, fmt.Errorf("failed to get token symbol for Managed token faucet: %w", err)
			}
			faucet := chainState.ManagedTokenFaucets[symbol]
			faucet.StateObjectId = addr
			chainState.ManagedTokenFaucets[symbol] = faucet
		case SuiManagedTokenFaucetUpgradeCapObjectIDType:
			symbol, err := getTokenSymbol(typeAndVersion)
			if err != nil {
				return CCIPChainState{}, fmt.Errorf("failed to get token symbol for Managed token faucet: %w", err)
			}
			faucet := chainState.ManagedTokenFaucets[symbol]
			faucet.UpgradeCapObjectId = addr
			chainState.ManagedTokenFaucets[symbol] = faucet

		// mock upgrade related
		case SuiOnRampMockV2:
			chainState.OnRampMockV2PackageId = addr
		case SuiOffRampMockV2:
			chainState.OffRampMockV2PackageId = addr
		case SuiCCIPMockV2:
			chainState.CCIPMockV2PackageId = addr

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
		case SuiLnRTokenPoolRebalancerCapIDType:
			symbol, err := getTokenSymbol(typeAndVersion)
			if err != nil {
				return CCIPChainState{}, fmt.Errorf("failed to get token symbol for LnR token pool: %w", err)
			}
			pool := chainState.LnRTokenPools[symbol]
			pool.RebalancerCapIds = append(pool.RebalancerCapIds, addr)
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
