package adapters

import (
	"fmt"
	"math/big"
	"strings"

	semver "github.com/Masterminds/semver/v3"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldfsui "github.com/smartcontractkit/chainlink-deployments-framework/chain/sui"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cldf_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	tokensapi "github.com/smartcontractkit/chainlink-ccip/deployment/tokens"
	"github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	suideploy "github.com/smartcontractkit/chainlink-sui/deployment"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	burnminttokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_burn_mint_token_pool"
	managedtokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_managed_token_pool"
	coin_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops/coin"
	suideployutils "github.com/smartcontractkit/chainlink-sui/deployment/utils"
)

var (
	_ tokensapi.TokenAdapter    = &SuiTokenAdapter{}
	_ tokensapi.TokenRefResolver = &SuiTokenAdapter{}
)

// SuiTokenAdapter implements tokensapi.TokenAdapter and tokensapi.TokenRefResolver for
// Sui CCIP token pools.
//
// Addressing model: on Sui a token is identified by its coin type string
// (0x<package>::<module>::<STRUCT>) and a token pool by its Move package ID (a hex object
// id). The generic TokenAdapter abstraction is address-string based, so this adapter maps:
//   - token AddressRef.Address  -> coin type string
//   - pool  AddressRef.Address  -> pool package id; Qualifier -> token symbol
//
// Sui pool state and owner-cap objects are stored in the datastore keyed by a symbol label
// (not a qualifier), so pool-object resolution filters by contract type then matches the
// symbol label. The framework has no label filter, so a local helper is used.
//
// Sequence methods run in MCMS proposal mode: the Sui ops are invoked with Signer=nil so
// they return a TransactionCall, which is bridged into OnChainOutput.BatchOps for the
// generic changesets to assemble MCMS proposals. Methods that publish packages or execute
// reads are the exception.
type SuiTokenAdapter struct{}

// ================================================================
// === Trivial methods                                           ===
// ================================================================

// AddressRefToBytes serializes a Sui AddressRef to bytes. Pool package ids and object ids
// are hex; coin type strings contain "::" and are returned as their raw UTF-8 bytes so the
// coin type round-trips through DeriveTokenDecimals.
func (a *SuiTokenAdapter) AddressRefToBytes(ref datastore.AddressRef) ([]byte, error) {
	if ref.Address == "" {
		return nil, fmt.Errorf("address ref has empty address")
	}
	if strings.Contains(ref.Address, "::") {
		return []byte(ref.Address), nil
	}
	return suideploy.StrToBytes(ref.Address)
}

// DeriveTokenPoolCounterpart returns the token pool bytes unchanged. Sui token pools are
// objects addressed directly; there is no PDA-style derivation from the token like Solana.
func (a *SuiTokenAdapter) DeriveTokenPoolCounterpart(_ deployment.Environment, _ uint64, tokenPool []byte, _ []byte) ([]byte, error) {
	return tokenPool, nil
}

// MigrateLockReleasePoolLiquiditySequence is not supported on Sui. The lockbox-based v2.0
// liquidity migration is an EVM-only flow; nil signals no support.
func (a *SuiTokenAdapter) MigrateLockReleasePoolLiquiditySequence() *cldf_ops.Sequence[tokensapi.MigrateLockReleasePoolLiquidityInput, sequences.OnChainOutput, cldf_chain.BlockChains] {
	return nil
}

// DeployTokenVerify validates a Sui token deployment input. Sui has no fixed decimal cap
// equivalent to EVM's 18 and the deployment is verified on-chain by the publish transaction,
// so this is a no-op for now, mirroring the Solana adapter.
func (a *SuiTokenAdapter) DeployTokenVerify(_ deployment.Environment, _ tokensapi.DeployTokenInput) error {
	return nil
}

// ================================================================
// === Medium methods                                            ===
// ================================================================

// DeriveTokenDecimals reads the token decimals from on-chain coin metadata. The token
// bytes carry the Sui coin type string (see AddressRefToBytes).
func (a *SuiTokenAdapter) DeriveTokenDecimals(e deployment.Environment, chainSelector uint64, _ datastore.AddressRef, token []byte) (uint8, error) {
	chain, ok := e.BlockChains.SuiChains()[chainSelector]
	if !ok {
		return 0, fmt.Errorf("sui chain with selector %d not found", chainSelector)
	}
	coinType := string(token)
	if coinType == "" {
		return 0, fmt.Errorf("token bytes are empty; expected the sui coin type string")
	}
	report, err := cldf_ops.ExecuteOperation(e.OperationsBundle, coin_ops.GetCoinSymbolOp, suiDeps(chain), coinType)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch coin metadata for type %s: %w", coinType, err)
	}
	if report.Output.Decimals < 0 || report.Output.Decimals > 255 {
		return 0, fmt.Errorf("invalid decimals %d for coin type %s", report.Output.Decimals, coinType)
	}
	return uint8(report.Output.Decimals), nil
}

// SetTokenPoolRateLimits sets the default-lane rate limits on a Sui token pool as an MCMS
// proposal. Sui pools have a single bucket per remote lane, so fastFinality buckets are
// not supported and only the default bucket is applied.
func (a *SuiTokenAdapter) SetTokenPoolRateLimits() *cldf_ops.Sequence[tokensapi.TPRLRemotes, sequences.OnChainOutput, cldf_chain.BlockChains] {
	return cldf_ops.NewSequence(
		"sui-adapter:set-token-pool-rate-limits",
		semver.MustParse("1.6.0"),
		"Set rate limits on a Sui token pool as an MCMS proposal",
		func(b cldf_ops.Bundle, chains cldf_chain.BlockChains, input tokensapi.TPRLRemotes) (sequences.OnChainOutput, error) {
			rl, ok := input.GetBucketForFinality(false)
			if !ok {
				b.Logger.Warnf("sui SetTokenPoolRateLimits: no default rate-limit bucket for pool %s on chain %d; skipping", input.TokenPoolRef.Address, input.ChainSelector)
				return sequences.OnChainOutput{}, nil
			}

			chain, ok := chains.SuiChains()[input.ChainSelector]
			if !ok {
				return sequences.OnChainOutput{}, fmt.Errorf("sui chain with selector %d not found", input.ChainSelector)
			}
			coinType := input.TokenRef.Address
			if coinType == "" {
				return sequences.OnChainOutput{}, fmt.Errorf("token ref has no coin type address on chain %d", input.ChainSelector)
			}
			stateObjID, ownerCapID, err := resolveSuiPoolObjects(input.ExistingDataStore, input.ChainSelector, input.TokenPoolRef.Type, input.TokenPoolRef.Qualifier)
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("failed to resolve sui pool objects: %w", err)
			}

			obCap, err := bigToU64(rl.OutboundRateLimiterConfig.Capacity)
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("outbound capacity: %w", err)
			}
			obRate, err := bigToU64(rl.OutboundRateLimiterConfig.Rate)
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("outbound rate: %w", err)
			}
			ibCap, err := bigToU64(rl.InboundRateLimiterConfig.Capacity)
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("inbound capacity: %w", err)
			}
			ibRate, err := bigToU64(rl.InboundRateLimiterConfig.Rate)
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("inbound rate: %w", err)
			}

			deps := suiDeps(chain)
			remote := input.RemoteChainSelector
			var call sui_ops.TransactionCall
			switch input.TokenPoolRef.Type {
			case datastore.ContractType(suideploy.SuiBnMTokenPoolType):
				r, err := cldf_ops.ExecuteOperation(b, burnminttokenpoolops.BurnMintTokenPoolSetChainRateLimiterOp, deps, burnminttokenpoolops.BurnMintTokenPoolSetChainRateLimiterInput{
					BurnMintPackageId:    input.TokenPoolRef.Address,
					CoinObjectTypeArg:    coinType,
					StateObjectId:        stateObjID,
					OwnerCap:             ownerCapID,
					RemoteChainSelectors: []uint64{remote},
					OutboundIsEnableds:   []bool{rl.OutboundRateLimiterConfig.IsEnabled},
					OutboundCapacities:   []uint64{obCap},
					OutboundRates:        []uint64{obRate},
					InboundIsEnableds:    []bool{rl.InboundRateLimiterConfig.IsEnabled},
					InboundCapacities:    []uint64{ibCap},
					InboundRates:         []uint64{ibRate},
				})
				if err != nil {
					return sequences.OnChainOutput{}, fmt.Errorf("failed to set rate limits on burn-mint pool: %w", err)
				}
				call = r.Output.Call
			case datastore.ContractType(suideploy.SuiManagedTokenPoolType):
				r, err := cldf_ops.ExecuteOperation(b, managedtokenpoolops.ManagedTokenPoolSetChainRateLimiterOp, deps, managedtokenpoolops.ManagedTokenPoolSetChainRateLimiterInput{
					ManagedTokenPoolPackageId: input.TokenPoolRef.Address,
					CoinObjectTypeArg:         coinType,
					StateObjectId:             stateObjID,
					OwnerCap:                  ownerCapID,
					RemoteChainSelectors:      []uint64{remote},
					OutboundIsEnableds:        []bool{rl.OutboundRateLimiterConfig.IsEnabled},
					OutboundCapacities:        []uint64{obCap},
					OutboundRates:             []uint64{obRate},
					InboundIsEnableds:         []bool{rl.InboundRateLimiterConfig.IsEnabled},
					InboundCapacities:         []uint64{ibCap},
					InboundRates:              []uint64{ibRate},
				})
				if err != nil {
					return sequences.OnChainOutput{}, fmt.Errorf("failed to set rate limits on managed pool: %w", err)
				}
				call = r.Output.Call
			default:
				return sequences.OnChainOutput{}, fmt.Errorf("unsupported sui token pool type %s for SetTokenPoolRateLimits", input.TokenPoolRef.Type)
			}
			return batchOpFromCall(input.ChainSelector, call)
		},
	)
}

// UpdateAuthorities transfers a Sui token pool's ownership to the MCMS timelock as an MCMS
// proposal. On Sui this is a single execute_ownership_transfer_to_mcms call (not the EVM
// transfer+accept pair); the target is the normal MCMS package id with the MCMS registry.
func (a *SuiTokenAdapter) UpdateAuthorities() *cldf_ops.Sequence[tokensapi.UpdateAuthoritiesInput, sequences.OnChainOutput, *deployment.Environment] {
	return cldf_ops.NewSequence(
		"sui-adapter:update-authorities",
		semver.MustParse("1.6.0"),
		"Transfer Sui token pool ownership to the MCMS timelock as an MCMS proposal",
		func(b cldf_ops.Bundle, e *deployment.Environment, input tokensapi.UpdateAuthoritiesInput) (sequences.OnChainOutput, error) {
			chain, ok := e.BlockChains.SuiChains()[input.ChainSelector]
			if !ok {
				return sequences.OnChainOutput{}, fmt.Errorf("sui chain with selector %d not found", input.ChainSelector)
			}
			coinType := input.TokenRef.Address
			if coinType == "" {
				return sequences.OnChainOutput{}, fmt.Errorf("token ref has no coin type address on chain %d", input.ChainSelector)
			}
			stateObjID, ownerCapID, err := resolveSuiPoolObjects(e.DataStore, input.ChainSelector, input.TokenPoolRef.Type, input.TokenPoolRef.Qualifier)
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("failed to resolve sui pool objects: %w", err)
			}

			mcmsPkgRef, ok := findRefExcludingLabel(findRefsByType(e.DataStore, input.ChainSelector, datastore.ContractType(suideploy.SuiMcmsPackageIDType)), suideploy.MCMSFastCurseLabel)
			if !ok {
				return sequences.OnChainOutput{}, fmt.Errorf("normal (non-fastcurse) MCMS package not found on chain %d", input.ChainSelector)
			}
			mcmsRegRef, ok := findRefExcludingLabel(findRefsByType(e.DataStore, input.ChainSelector, datastore.ContractType(suideploy.SuiMcmsRegistryObjectIDType)), suideploy.MCMSFastCurseLabel)
			if !ok {
				return sequences.OnChainOutput{}, fmt.Errorf("normal (non-fastcurse) MCMS registry not found on chain %d", input.ChainSelector)
			}

			deps := suiDeps(chain)
			typeArgs := []string{coinType}
			var call sui_ops.TransactionCall
			switch input.TokenPoolRef.Type {
			case datastore.ContractType(suideploy.SuiBnMTokenPoolType):
				r, err := cldf_ops.ExecuteOperation(b, burnminttokenpoolops.ExecuteOwnershipTransferToMcmsBurnMintTokenPoolOp, deps, burnminttokenpoolops.ExecuteOwnershipTransferToMcmsBurnMintTokenPoolInput{
					BurnMintTokenPoolPackageId: input.TokenPoolRef.Address,
					TypeArgs:                   typeArgs,
					OwnerCapObjectId:           ownerCapID,
					StateObjectId:              stateObjID,
					RegistryObjectId:           mcmsRegRef.Address,
					To:                         mcmsPkgRef.Address,
				})
				if err != nil {
					return sequences.OnChainOutput{}, fmt.Errorf("failed to transfer burn-mint pool ownership to MCMS: %w", err)
				}
				call = r.Output.Call
			case datastore.ContractType(suideploy.SuiManagedTokenPoolType):
				r, err := cldf_ops.ExecuteOperation(b, managedtokenpoolops.ExecuteOwnershipTransferToMcmsManagedTokenPoolOp, deps, managedtokenpoolops.ExecuteOwnershipTransferToMcmsManagedTokenPoolInput{
					ManagedTokenPoolPackageId: input.TokenPoolRef.Address,
					TypeArgs:                  typeArgs,
					OwnerCapObjectId:          ownerCapID,
					StateObjectId:             stateObjID,
					RegistryObjectId:          mcmsRegRef.Address,
					To:                        mcmsPkgRef.Address,
				})
				if err != nil {
					return sequences.OnChainOutput{}, fmt.Errorf("failed to transfer managed pool ownership to MCMS: %w", err)
				}
				call = r.Output.Call
			default:
				return sequences.OnChainOutput{}, fmt.Errorf("unsupported sui token pool type %s for UpdateAuthorities", input.TokenPoolRef.Type)
			}
			return batchOpFromCall(input.ChainSelector, call)
		},
	)
}

// ================================================================
// === TokenRefResolver                                          ===
// ================================================================

// ResolveTokenPoolRef resolves a Sui pool package id to its datastore AddressRef. The
// package id is looked up in the datastore to recover the pool contract type and the
// token symbol label. The returned ref carries the symbol as its Qualifier so downstream
// adapter methods can resolve the pool's state and owner-cap objects by symbol.
func (a *SuiTokenAdapter) ResolveTokenPoolRef(_ cldf_ops.Bundle, _ cldf_chain.BlockChains, ds datastore.DataStore, chainSelector uint64, address string) (datastore.AddressRef, error) {
	refs := findRefsByAddress(ds, chainSelector, address)
	if len(refs) == 0 {
		return datastore.AddressRef{}, fmt.Errorf("no sui address ref found for %s on chain %d", address, chainSelector)
	}
	for _, r := range refs {
		if !isSuiPoolType(r.Type) {
			continue
		}
		return datastore.AddressRef{
			ChainSelector: chainSelector,
			Type:          r.Type,
			Address:       r.Address,
			Qualifier:     symbolFromLabels(r),
			Version:       r.Version,
			Labels:        r.Labels,
		}, nil
	}
	return datastore.AddressRef{}, fmt.Errorf("address %s on chain %d is not a recognized sui token pool type", address, chainSelector)
}

// ResolveTokenRef resolves a Sui coin type string to an AddressRef. The coin type is the
// canonical token identifier on Sui. The symbol is fetched from on-chain coin metadata on
// a best-effort basis (falling back to the coin type) so it can be used as the qualifier.
func (a *SuiTokenAdapter) ResolveTokenRef(b cldf_ops.Bundle, chains cldf_chain.BlockChains, _ datastore.DataStore, chainSelector uint64, address string) (datastore.AddressRef, error) {
	coinType := normalizeCoinType(address)
	if !strings.Contains(coinType, "::") {
		return datastore.AddressRef{}, fmt.Errorf("sui token address %q is not a valid coin type (expected 0x<package>::<module>::<STRUCT>)", address)
	}
	chain, ok := chains.SuiChains()[chainSelector]
	if !ok {
		return datastore.AddressRef{}, fmt.Errorf("sui chain with selector %d not found", chainSelector)
	}
	symbol := coinType
	if report, err := cldf_ops.ExecuteOperation(b, coin_ops.GetCoinSymbolOp, suiDeps(chain), coinType); err != nil {
		b.Logger.Warnf("failed to fetch coin metadata for %s, using coin type as qualifier: %v", coinType, err)
	} else if report.Output.Symbol != "" {
		symbol = report.Output.Symbol
	}
	return datastore.AddressRef{
		ChainSelector: chainSelector,
		Type:          suiTokenType(coinType),
		Address:       coinType,
		Qualifier:     symbol,
		Version:       semver.MustParse("1.0.0"),
	}, nil
}

// ================================================================
// === Stubs: not yet implemented                                ===
// ================================================================

func (a *SuiTokenAdapter) ConfigureTokenForTransfersSequence() *cldf_ops.Sequence[tokensapi.ConfigureTokenForTransfersInput, sequences.OnChainOutput, cldf_chain.BlockChains] {
	return nil
}

// DeriveTokenAddress is not implemented yet. Recovering the coin type from a pool requires
// a DevInspect GetToken call whose typeArgs themselves need the coin type, so the coin type
// must come from the token ref / datastore rather than the pool. This is wired as part of
// the Tier 3 ConfigureTokenForTransfers work.
func (a *SuiTokenAdapter) DeriveTokenAddress(_ deployment.Environment, _ uint64, _ datastore.AddressRef) (string, error) {
	return "", fmt.Errorf("DeriveTokenAddress is not implemented on SuiTokenAdapter yet")
}

func (a *SuiTokenAdapter) ManualRegistration() *cldf_ops.Sequence[tokensapi.ManualRegistrationSequenceInput, sequences.OnChainOutput, cldf_chain.BlockChains] {
	return nil
}

func (a *SuiTokenAdapter) DeployToken() *cldf_ops.Sequence[tokensapi.DeployTokenInput, sequences.OnChainOutput, cldf_chain.BlockChains] {
	return nil
}

// DeployTokenPoolForToken is not implemented yet. The Sui DeployAndInitAllTokenPoolsSequence
// publishes a Move package and requires several fields the generic DeployTokenPoolInput does
// not carry (CCIP/MCMS package ids, fastcurse package, CCIPObjectRef, token administrator,
// and the coin type). Package publish also executes directly rather than via MCMS. The
// mapping must be validated against the Sui deploy sequence before wiring.
func (a *SuiTokenAdapter) DeployTokenPoolForToken() *cldf_ops.Sequence[tokensapi.DeployTokenPoolInput, sequences.OnChainOutput, cldf_chain.BlockChains] {
	return nil
}

// ================================================================
// === Helpers                                                   ===
// ================================================================

// suiDeps builds OpTxDeps in proposal mode: Signer is nil so ops return a TransactionCall
// instead of executing. Reads (coin metadata) also work with a nil signer.
func suiDeps(chain cldfsui.Chain) sui_ops.OpTxDeps {
	gas := uint64(400_000_000)
	return sui_ops.OpTxDeps{
		Client: chain.Client,
		Signer: nil,
		GetCallOpts: func() *bind.CallOpts {
			return &bind.CallOpts{WaitForExecution: true, GasBudget: &gas}
		},
		SuiRPC: chain.URL,
	}
}

// batchOpFromCall bridges a Sui TransactionCall into an OnChainOutput carrying a single
// MCMS BatchOperation, matching the pattern used by the Sui RMN curse/uncurse sequences.
func batchOpFromCall(chainSelector uint64, call sui_ops.TransactionCall) (sequences.OnChainOutput, error) {
	tx, err := suideployutils.TransactionCallToMCMSTransaction(call)
	if err != nil {
		return sequences.OnChainOutput{}, fmt.Errorf("failed to convert sui call to mcms transaction: %w", err)
	}
	return sequences.OnChainOutput{
		BatchOps: []mcmstypes.BatchOperation{{
			ChainSelector: mcmstypes.ChainSelector(chainSelector),
			Transactions:  []mcmstypes.Transaction{tx},
		}},
	}, nil
}

func findRefsByType(ds datastore.DataStore, selector uint64, contractType datastore.ContractType) []datastore.AddressRef {
	if ds == nil {
		return nil
	}
	return ds.Addresses().Filter(
		datastore.AddressRefByChainSelector(selector),
		datastore.AddressRefByType(contractType),
	)
}

func findRefsByAddress(ds datastore.DataStore, selector uint64, address string) []datastore.AddressRef {
	if ds == nil {
		return nil
	}
	return ds.Addresses().Filter(
		datastore.AddressRefByChainSelector(selector),
		datastore.AddressRefByAddress(address),
	)
}

func findRefByLabel(refs []datastore.AddressRef, label string) (datastore.AddressRef, bool) {
	for _, r := range refs {
		if r.Labels.Contains(label) {
			return r, true
		}
	}
	return datastore.AddressRef{}, false
}

func findRefExcludingLabel(refs []datastore.AddressRef, label string) (datastore.AddressRef, bool) {
	for _, r := range refs {
		if !r.Labels.Contains(label) {
			return r, true
		}
	}
	return datastore.AddressRef{}, false
}

// resolveSuiPoolObjects returns the pool's state object id and owner-cap object id by
// matching the token symbol label against the pool's state and owner-cap contract types.
func resolveSuiPoolObjects(ds datastore.DataStore, selector uint64, poolType datastore.ContractType, symbol string) (stateObjID, ownerCapID string, err error) {
	if symbol == "" {
		return "", "", fmt.Errorf("pool ref has no symbol qualifier; cannot resolve sui pool state and owner-cap")
	}
	stateType, ownerType, err := poolObjectTypes(poolType)
	if err != nil {
		return "", "", err
	}
	stateRef, ok := findRefByLabel(findRefsByType(ds, selector, stateType), symbol)
	if !ok {
		return "", "", fmt.Errorf("no sui pool state ref for symbol %s (type %s) on chain %d", symbol, stateType, selector)
	}
	ownerRef, ok := findRefByLabel(findRefsByType(ds, selector, ownerType), symbol)
	if !ok {
		return "", "", fmt.Errorf("no sui pool owner-cap ref for symbol %s (type %s) on chain %d", symbol, ownerType, selector)
	}
	return stateRef.Address, ownerRef.Address, nil
}

func poolObjectTypes(poolType datastore.ContractType) (stateType, ownerType datastore.ContractType, err error) {
	switch poolType {
	case datastore.ContractType(suideploy.SuiBnMTokenPoolType):
		return datastore.ContractType(suideploy.SuiBnMTokenPoolStateType), datastore.ContractType(suideploy.SuiBnMTokenPoolOwnerIDType), nil
	case datastore.ContractType(suideploy.SuiManagedTokenPoolType):
		return datastore.ContractType(suideploy.SuiManagedTokenPoolStateType), datastore.ContractType(suideploy.SuiManagedTokenPoolOwnerIDType), nil
	case datastore.ContractType(suideploy.SuiLnRTokenPoolType):
		return datastore.ContractType(suideploy.SuiLnRTokenPoolStateType), datastore.ContractType(suideploy.SuiLnRTokenPoolOwnerIDType), nil
	default:
		return "", "", fmt.Errorf("unsupported sui token pool type %s", poolType)
	}
}

func isSuiPoolType(t datastore.ContractType) bool {
	switch t {
	case datastore.ContractType(suideploy.SuiBnMTokenPoolType),
		datastore.ContractType(suideploy.SuiManagedTokenPoolType),
		datastore.ContractType(suideploy.SuiLnRTokenPoolType):
		return true
	}
	return false
}

// symbolFromLabels returns the first label, which is the token symbol for Sui pool refs.
func symbolFromLabels(r datastore.AddressRef) string {
	labels := r.Labels.List()
	if len(labels) > 0 {
		return labels[0]
	}
	return ""
}

// suiTokenType maps a coin type to a datastore contract type on a best-effort basis. CCIP
// burn-mint test tokens on Sui are managed coins, so non-LINK coins default to the managed
// token type.
func suiTokenType(coinType string) datastore.ContractType {
	switch {
	case strings.Contains(coinType, "::link::"):
		return datastore.ContractType(suideploy.SuiLinkTokenType)
	case strings.Contains(coinType, "managed_token"):
		return datastore.ContractType(suideploy.SuiManagedTokenType)
	default:
		return datastore.ContractType(suideploy.SuiManagedTokenType)
	}
}

func normalizeCoinType(s string) string {
	if strings.HasPrefix(s, "0x") {
		return s
	}
	return "0x" + s
}

func bigToU64(x *big.Int) (uint64, error) {
	if x == nil {
		return 0, nil
	}
	if !x.IsUint64() {
		return 0, fmt.Errorf("value %s does not fit in uint64", x.String())
	}
	return x.Uint64(), nil
}
