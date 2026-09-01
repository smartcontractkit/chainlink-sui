package adapters

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	semver "github.com/Masterminds/semver/v3"
	chainsel "github.com/smartcontractkit/chain-selectors"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldfsui "github.com/smartcontractkit/chainlink-deployments-framework/chain/sui"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cldf_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	tokensapi "github.com/smartcontractkit/chainlink-ccip/deployment/tokens"
	cciputils "github.com/smartcontractkit/chainlink-ccip/deployment/utils"
	"github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_burn_mint_token_pool "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_token_pools/burn_mint_token_pool"
	module_lock_release_token_pool "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_token_pools/lock_release_token_pool"
	module_managed_token_pool "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_token_pools/managed_token_pool"
	suideploy "github.com/smartcontractkit/chainlink-sui/deployment"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	burnminttokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_burn_mint_token_pool"
	lockreleasetokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_lock_release_token_pool"
	managedtokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_managed_token_pool"
	coin_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops/coin"
	suideployutils "github.com/smartcontractkit/chainlink-sui/deployment/utils"
)

var (
	_ tokensapi.TokenAdapter     = &SuiTokenAdapter{}
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
// not supported and only the default bucket is applied. A missing default bucket means no
// rate-limit update for this lane and the sequence no-ops.
func (a *SuiTokenAdapter) SetTokenPoolRateLimits() *cldf_ops.Sequence[tokensapi.TPRLRemotes, sequences.OnChainOutput, cldf_chain.BlockChains] {
	return cldf_ops.NewSequence(
		"sui-adapter:set-token-pool-rate-limits",
		&suideploy.Version1_0_0,
		"Set rate limits on a Sui token pool as an MCMS proposal",
		func(b cldf_ops.Bundle, chains cldf_chain.BlockChains, input tokensapi.TPRLRemotes) (sequences.OnChainOutput, error) {
			rl, ok := input.GetBucketForFinality(false)
			if !ok {
				b.Logger.Warnf("skipping rate limiter config for token pool (%s) on chain %d since no default bucket was provided", input.TokenPoolRef.Address, input.ChainSelector)
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
			stateObjID, ownerCapID, err := resolveSuiPoolObjects(input.ExistingDataStore, input.ChainSelector, input.TokenPoolRef.Type, suiPoolSymbol(input.TokenPoolRef))
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
			case datastore.ContractType(suideploy.SuiLnRTokenPoolType):
				r, err := cldf_ops.ExecuteOperation(b, lockreleasetokenpoolops.LockReleaseTokenPoolSetChainRateLimiterOp, deps, lockreleasetokenpoolops.LockReleaseTokenPoolSetChainRateLimiterInput{
					LockReleasePackageId: input.TokenPoolRef.Address,
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
					return sequences.OnChainOutput{}, fmt.Errorf("failed to set rate limits on lock-release pool: %w", err)
				}
				call = r.Output.Call
			default:
				return sequences.OnChainOutput{}, fmt.Errorf("unsupported sui token pool type %s for SetTokenPoolRateLimits", input.TokenPoolRef.Type)
			}
			return batchOpFromCall(input.ChainSelector, call)
		},
	)
}

// UpdateAuthorities returns ownership step 2 of 3 as an MCMS proposal: the MCMS timelock
// calls accept_ownership on the pool, which asserts ctx.sender() == pending_transfer.to == MCMS.
//
// The full Sui pool→MCMS ownership flow is EOA/MCMS/EOA:
//  1. EOA transfer_ownership(To: MCMS) — performed by DeployTokenPoolForToken at deploy time,
//     setting a pending (not-yet-accepted) transfer.
//  2. MCMS accept_ownership — this method, returned as an MCMS proposal (Signer=nil → Call).
//  3. EOA execute_ownership_transfer_to_mcms — consumes the OwnerCap and registers the MCMS
//     entrypoint; run EOA-direct via the Sui cs_mcms_execute_ownership_transfer changeset
//     after this accept proposal lands. ownable::execute_ownership_transfer_to_mcms asserts
//     pending_transfer.is_some() (ENoPendingTransfer) and pending_transfer.accepted
//     (ETransferNotAccepted), so it cannot run before steps 1-2.
//
// OnChainOutput.BatchOps can only carry MCMS proposals, so only step 2 belongs here.
func (a *SuiTokenAdapter) UpdateAuthorities() *cldf_ops.Sequence[tokensapi.UpdateAuthoritiesInput, sequences.OnChainOutput, *deployment.Environment] {
	return cldf_ops.NewSequence(
		"sui-adapter:update-authorities",
		&suideploy.Version1_0_0,
		"Propose MCMS accept_ownership of a Sui token pool (ownership step 2 of 3)",
		func(b cldf_ops.Bundle, e *deployment.Environment, input tokensapi.UpdateAuthoritiesInput) (sequences.OnChainOutput, error) {
			chain, ok := e.BlockChains.SuiChains()[input.ChainSelector]
			if !ok {
				return sequences.OnChainOutput{}, fmt.Errorf("sui chain with selector %d not found", input.ChainSelector)
			}
			coinType := input.TokenRef.Address
			if coinType == "" {
				return sequences.OnChainOutput{}, fmt.Errorf("token ref has no coin type address on chain %d", input.ChainSelector)
			}
			stateObjID, _, err := resolveSuiPoolObjects(e.DataStore, input.ChainSelector, input.TokenPoolRef.Type, suiPoolSymbol(input.TokenPoolRef))
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("failed to resolve sui pool objects: %w", err)
			}

			deps := suiDeps(chain)
			typeArgs := []string{coinType}
			var call sui_ops.TransactionCall
			switch input.TokenPoolRef.Type {
			case datastore.ContractType(suideploy.SuiBnMTokenPoolType):
				r, err := cldf_ops.ExecuteOperation(b, burnminttokenpoolops.AcceptOwnershipBurnMintTokenPoolOp, deps, burnminttokenpoolops.AcceptOwnershipBurnMintTokenPoolInput{
					BurnMintTokenPoolPackageId: input.TokenPoolRef.Address,
					TypeArgs:                   typeArgs,
					StateObjectId:              stateObjID,
				})
				if err != nil {
					return sequences.OnChainOutput{}, fmt.Errorf("failed to propose accept_ownership for burn-mint pool: %w", err)
				}
				call = r.Output.Call
			case datastore.ContractType(suideploy.SuiManagedTokenPoolType):
				r, err := cldf_ops.ExecuteOperation(b, managedtokenpoolops.AcceptOwnershipManagedTokenPoolOp, deps, managedtokenpoolops.AcceptOwnershipManagedTokenPoolInput{
					ManagedTokenPoolPackageId: input.TokenPoolRef.Address,
					TypeArgs:                  typeArgs,
					StateObjectId:             stateObjID,
				})
				if err != nil {
					return sequences.OnChainOutput{}, fmt.Errorf("failed to propose accept_ownership for managed pool: %w", err)
				}
				call = r.Output.Call
			case datastore.ContractType(suideploy.SuiLnRTokenPoolType):
				r, err := cldf_ops.ExecuteOperation(b, lockreleasetokenpoolops.AcceptOwnershipLockReleaseTokenPoolOp, deps, lockreleasetokenpoolops.AcceptOwnershipLockReleaseTokenPoolInput{
					LockReleaseTokenPoolPackageId: input.TokenPoolRef.Address,
					TypeArgs:                      typeArgs,
					StateObjectId:                 stateObjID,
				})
				if err != nil {
					return sequences.OnChainOutput{}, fmt.Errorf("failed to propose accept_ownership for lock-release pool: %w", err)
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

// ConfigureTokenForTransfersSequence configures a Sui token pool for cross-chain transfers
// as an MCMS proposal. Per remote chain it applies the chain config (remote token + remote
// pool) and the default-lane rate limits, mirroring the Sui ConfigureBurnMintTokenPool flow.
// Sui has no router SetPool and the existing Sui configure flow does not register the pool in
// the TokenAdminRegistry, so this sequence only sets per-remote chain config + rate limits.
func (a *SuiTokenAdapter) ConfigureTokenForTransfersSequence() *cldf_ops.Sequence[tokensapi.ConfigureTokenForTransfersInput, sequences.OnChainOutput, cldf_chain.BlockChains] {
	return cldf_ops.NewSequence(
		"sui-adapter:configure-token-for-transfers",
		&suideploy.Version1_0_0,
		"Configure a Sui token pool for cross-chain transfers as an MCMS proposal",
		func(b cldf_ops.Bundle, chains cldf_chain.BlockChains, input tokensapi.ConfigureTokenForTransfersInput) (sequences.OnChainOutput, error) {
			chain, ok := chains.SuiChains()[input.ChainSelector]
			if !ok {
				return sequences.OnChainOutput{}, fmt.Errorf("sui chain with selector %d not found", input.ChainSelector)
			}
			coinType := input.TokenRef.Address
			if coinType == "" {
				return sequences.OnChainOutput{}, fmt.Errorf("token ref has no coin type address on chain %d", input.ChainSelector)
			}
			pkgID := input.TokenPoolAddress
			if pkgID == "" {
				return sequences.OnChainOutput{}, fmt.Errorf("token pool address is empty on chain %d", input.ChainSelector)
			}
			poolType, err := suiPoolTypeFromStr(input.PoolType)
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("unsupported sui pool type %q on chain %d: %w", input.PoolType, input.ChainSelector, err)
			}
			stateObjID, ownerCapID, err := resolveSuiPoolObjects(input.ExistingDataStore, input.ChainSelector, poolType, suiPoolSymbol(input.TokenRef))
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("failed to resolve sui pool objects: %w", err)
			}
			localDecimals, err := suiDecimals(b, chain, coinType)
			if err != nil {
				return sequences.OnChainOutput{}, err
			}

			// Configure-before-own: while the deployer still owns the pool, run the
			// OwnerCap-gated config entrypoints directly. A bundled MCMS config op would
			// abort because token_expansion never registers the package with mcms_registry;
			// only step 3 execute_ownership_transfer does. MCMS-owned pools fall back to collect.
			deployerAddr, err := chain.Signer.GetAddress()
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("failed to get deployer address on chain %d: %w", input.ChainSelector, err)
			}
			owner, err := suiPoolOwner(b, chain, poolType, pkgID, coinType, stateObjID)
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("failed to read pool owner on chain %d: %w", input.ChainSelector, err)
			}
			eoDirect := suiAddrEqual(owner, deployerAddr)
			var deps sui_ops.OpTxDeps
			if eoDirect {
				deps = suiDepsExec(chain)
			} else {
				deps = suiDeps(chain)
			}
			batchOps := make([]mcmstypes.BatchOperation, 0)
			for remoteSelector, rc := range input.RemoteChains {
				if err := rc.Validate(); err != nil {
					return sequences.OnChainOutput{}, fmt.Errorf("remote chain config for selector %d: %w", remoteSelector, err)
				}
				remoteTokenHex := "0x" + hex.EncodeToString(rc.RemoteToken)
				remotePoolHex := "0x" + hex.EncodeToString(rc.RemotePool)

				calls := make([]sui_ops.TransactionCall, 0, 2)
				switch poolType {
				case datastore.ContractType(suideploy.SuiBnMTokenPoolType):
					r, err := cldf_ops.ExecuteOperation(b, burnminttokenpoolops.BurnMintTokenPoolApplyChainUpdatesOp, deps, burnminttokenpoolops.BurnMintTokenPoolApplyChainUpdatesInput{
						BurnMintPackageId:         pkgID,
						CoinObjectTypeArg:         coinType,
						StateObjectId:             stateObjID,
						OwnerCap:                  ownerCapID,
						RemoteChainSelectorsToAdd: []uint64{remoteSelector},
						RemotePoolAddressesToAdd:  [][]string{{remotePoolHex}},
						RemoteTokenAddressesToAdd: []string{remoteTokenHex},
					})
					if err != nil {
						return sequences.OnChainOutput{}, fmt.Errorf("apply chain updates for remote %d: %w", remoteSelector, err)
					}
					calls = append(calls, r.Output.Call)
				case datastore.ContractType(suideploy.SuiManagedTokenPoolType):
					r, err := cldf_ops.ExecuteOperation(b, managedtokenpoolops.ManagedTokenPoolApplyChainUpdatesOp, deps, managedtokenpoolops.ManagedTokenPoolApplyChainUpdatesInput{
						ManagedTokenPoolPackageId: pkgID,
						CoinObjectTypeArg:         coinType,
						StateObjectId:             stateObjID,
						OwnerCap:                  ownerCapID,
						RemoteChainSelectorsToAdd: []uint64{remoteSelector},
						RemotePoolAddressesToAdd:  [][]string{{remotePoolHex}},
						RemoteTokenAddressesToAdd: []string{remoteTokenHex},
					})
					if err != nil {
						return sequences.OnChainOutput{}, fmt.Errorf("apply chain updates for remote %d: %w", remoteSelector, err)
					}
					calls = append(calls, r.Output.Call)
				case datastore.ContractType(suideploy.SuiLnRTokenPoolType):
					r, err := cldf_ops.ExecuteOperation(b, lockreleasetokenpoolops.LockReleaseTokenPoolApplyChainUpdatesOp, deps, lockreleasetokenpoolops.LockReleaseTokenPoolApplyChainUpdatesInput{
						LockReleasePackageId:      pkgID,
						CoinObjectTypeArg:         coinType,
						StateObjectId:             stateObjID,
						OwnerCap:                  ownerCapID,
						RemoteChainSelectorsToAdd: []uint64{remoteSelector},
						RemotePoolAddressesToAdd:  [][]string{{remotePoolHex}},
						RemoteTokenAddressesToAdd: []string{remoteTokenHex},
					})
					if err != nil {
						return sequences.OnChainOutput{}, fmt.Errorf("apply chain updates for remote %d: %w", remoteSelector, err)
					}
					calls = append(calls, r.Output.Call)
				default:
					return sequences.OnChainOutput{}, fmt.Errorf("unsupported sui token pool type %s for ConfigureTokenForTransfers", poolType)
				}

				obBucket, obOk := rc.GetOutboundRateLimitBuckets().DefaultBucket()
				ibBucket, ibOk := rc.GetInboundRateLimitBuckets().DefaultBucket()
				switch {
				case obOk && ibOk:
					obRL, ibRL := tokensapi.GenerateTPRLConfigs(obBucket.RateLimit, ibBucket.RateLimit, localDecimals, rc.RemoteDecimals, chainsel.FamilySui, semver.MustParse("1.6.0"), string(poolType))
					rlCall, err := suiSetChainRateLimitCall(b, deps, poolType, pkgID, coinType, stateObjID, ownerCapID, remoteSelector, obRL, ibRL)
					if err != nil {
						return sequences.OnChainOutput{}, err
					}
					calls = append(calls, rlCall)
				case !obOk && !ibOk:
					// No rate-limit update for this remote; keep the chain-config call already in calls.
				default:
					return sequences.OnChainOutput{}, fmt.Errorf("default outbound and inbound rate limits must both be specified or both omitted for remote %d", remoteSelector)
				}

				// EOA-direct mode already executed the ops on-chain; skip collecting their
				// Calls into batchOps. Tx digests surface via ExecutionReports.
				if !eoDirect {
					out, err := batchOpFromCalls(input.ChainSelector, calls)
					if err != nil {
						return sequences.OnChainOutput{}, err
					}
					batchOps = append(batchOps, out.BatchOps...)
				}
			}
			return sequences.OnChainOutput{BatchOps: batchOps}, nil
		},
	)
}

// suiSetChainRateLimitCall builds the SetChainRateLimiter TransactionCall for one remote
// chain, converting the decimal-scaled configs to the u64 on-chain representation.
func suiSetChainRateLimitCall(
	b cldf_ops.Bundle,
	deps sui_ops.OpTxDeps,
	poolType datastore.ContractType,
	pkgID, coinType, stateObjID, ownerCapID string,
	remoteSelector uint64,
	obRL, ibRL tokensapi.RateLimiterConfig,
) (sui_ops.TransactionCall, error) {
	obCap, err := bigToU64(obRL.Capacity)
	if err != nil {
		return sui_ops.TransactionCall{}, fmt.Errorf("outbound capacity: %w", err)
	}
	obRate, err := bigToU64(obRL.Rate)
	if err != nil {
		return sui_ops.TransactionCall{}, fmt.Errorf("outbound rate: %w", err)
	}
	ibCap, err := bigToU64(ibRL.Capacity)
	if err != nil {
		return sui_ops.TransactionCall{}, fmt.Errorf("inbound capacity: %w", err)
	}
	ibRate, err := bigToU64(ibRL.Rate)
	if err != nil {
		return sui_ops.TransactionCall{}, fmt.Errorf("inbound rate: %w", err)
	}
	remotes := []uint64{remoteSelector}
	switch poolType {
	case datastore.ContractType(suideploy.SuiBnMTokenPoolType):
		r, err := cldf_ops.ExecuteOperation(b, burnminttokenpoolops.BurnMintTokenPoolSetChainRateLimiterOp, deps, burnminttokenpoolops.BurnMintTokenPoolSetChainRateLimiterInput{
			BurnMintPackageId:    pkgID,
			CoinObjectTypeArg:    coinType,
			StateObjectId:        stateObjID,
			OwnerCap:             ownerCapID,
			RemoteChainSelectors: remotes,
			OutboundIsEnableds:   []bool{obRL.IsEnabled},
			OutboundCapacities:   []uint64{obCap},
			OutboundRates:        []uint64{obRate},
			InboundIsEnableds:    []bool{ibRL.IsEnabled},
			InboundCapacities:    []uint64{ibCap},
			InboundRates:         []uint64{ibRate},
		})
		if err != nil {
			return sui_ops.TransactionCall{}, fmt.Errorf("set rate limits: %w", err)
		}
		return r.Output.Call, nil
	case datastore.ContractType(suideploy.SuiManagedTokenPoolType):
		r, err := cldf_ops.ExecuteOperation(b, managedtokenpoolops.ManagedTokenPoolSetChainRateLimiterOp, deps, managedtokenpoolops.ManagedTokenPoolSetChainRateLimiterInput{
			ManagedTokenPoolPackageId: pkgID,
			CoinObjectTypeArg:         coinType,
			StateObjectId:             stateObjID,
			OwnerCap:                  ownerCapID,
			RemoteChainSelectors:      remotes,
			OutboundIsEnableds:        []bool{obRL.IsEnabled},
			OutboundCapacities:        []uint64{obCap},
			OutboundRates:             []uint64{obRate},
			InboundIsEnableds:         []bool{ibRL.IsEnabled},
			InboundCapacities:         []uint64{ibCap},
			InboundRates:              []uint64{ibRate},
		})
		if err != nil {
			return sui_ops.TransactionCall{}, fmt.Errorf("set rate limits: %w", err)
		}
		return r.Output.Call, nil
	case datastore.ContractType(suideploy.SuiLnRTokenPoolType):
		r, err := cldf_ops.ExecuteOperation(b, lockreleasetokenpoolops.LockReleaseTokenPoolSetChainRateLimiterOp, deps, lockreleasetokenpoolops.LockReleaseTokenPoolSetChainRateLimiterInput{
			LockReleasePackageId: pkgID,
			CoinObjectTypeArg:    coinType,
			StateObjectId:        stateObjID,
			OwnerCap:             ownerCapID,
			RemoteChainSelectors: remotes,
			OutboundIsEnableds:   []bool{obRL.IsEnabled},
			OutboundCapacities:   []uint64{obCap},
			OutboundRates:        []uint64{obRate},
			InboundIsEnableds:    []bool{ibRL.IsEnabled},
			InboundCapacities:    []uint64{ibCap},
			InboundRates:         []uint64{ibRate},
		})
		if err != nil {
			return sui_ops.TransactionCall{}, fmt.Errorf("set rate limits: %w", err)
		}
		return r.Output.Call, nil
	default:
		return sui_ops.TransactionCall{}, fmt.Errorf("unsupported sui token pool type %s for rate limits", poolType)
	}
}

// DeriveTokenAddress derives the Sui coin type for a pool from the datastore. It is a
// fallback used when the caller does not provide a token ref; the primary path threads the
// coin type through the token ref. The coin type is built from the coin package ref matched by
// symbol, using the module::STRUCT suffix that ref carries as a coinType= label. This is
// token-agnostic: BnM, LINK, and any newly added coin derive identically, and the result is the
// same whether the coin sits behind a BurnMint or a Managed token pool.
func (a *SuiTokenAdapter) DeriveTokenAddress(e deployment.Environment, chainSelector uint64, poolRef datastore.AddressRef) (string, error) {
	return deriveSuiCoinType(e.DataStore, chainSelector, suiPoolSymbol(poolRef))
}

// ManualRegistration is a no-op on Sui. Unlike EVM/Solana, a Sui token pool self-registers
// in the TokenAdminRegistry during pool initialization: the Move `initialize` (burn-mint) and
// `initialize_with_managed_token` (managed) functions call `token_admin_registry::register_pool`
// themselves. DeployTokenPoolForToken triggers that init, so by the time a pool exists it is
// already registered. Calling register_pool again would abort with ETokenAlreadyRegistered, so
// this sequence logs and returns an empty output rather than re-registering. A nil return would
// crash the generic ManualRegistration changeset, so a no-op sequence is returned instead.
func (a *SuiTokenAdapter) ManualRegistration() *cldf_ops.Sequence[tokensapi.ManualRegistrationSequenceInput, sequences.OnChainOutput, cldf_chain.BlockChains] {
	return cldf_ops.NewSequence(
		"sui-adapter:manual-registration",
		&suideploy.Version1_0_0,
		"No-op: Sui pools self-register in the TokenAdminRegistry during initialization",
		func(b cldf_ops.Bundle, _ cldf_chain.BlockChains, input tokensapi.ManualRegistrationSequenceInput) (sequences.OnChainOutput, error) {
			b.Logger.Infow("Sui ManualRegistration is a no-op: token pools self-register in the TokenAdminRegistry during pool initialization",
				"chainSelector", input.ChainSelector, "tokenPoolRef", input.TokenPoolRef.Address)
			return sequences.OnChainOutput{}, nil
		},
	)
}

// DeployToken is not supported for Sui through the generic token adapter. Sui token
// deployment publishes a fixed-kind Move package (managed token or BnM token) whose deploy op
// takes only MCMS config, not the symbol/decimals/supply/pre-mint/senders that the generic
// DeployTokenInput carries; the token identity is determined by the published package, so the
// generic input does not map cleanly. Sui tokens are deployed via the Sui-specific token deploy
// changesets instead. This returns a sequence that errors so a misconfigured Sui chain fails
// with a clear message rather than nil-panicking the generic flow (which calls it unguarded).
func (a *SuiTokenAdapter) DeployToken() *cldf_ops.Sequence[tokensapi.DeployTokenInput, sequences.OnChainOutput, cldf_chain.BlockChains] {
	return cldf_ops.NewSequence(
		"sui-adapter:deploy-token",
		&suideploy.Version1_0_0,
		"Not supported: Sui tokens are deployed via the Sui-specific token deploy changesets",
		func(_ cldf_ops.Bundle, _ cldf_chain.BlockChains, _ tokensapi.DeployTokenInput) (sequences.OnChainOutput, error) {
			return sequences.OnChainOutput{}, fmt.Errorf("DeployToken is not supported for Sui via the generic token adapter; deploy Sui tokens with the Sui-specific token deploy changesets")
		},
	)
}

// DeployTokenPoolForToken deploys and initializes a Sui token pool for an existing token.
// It publishes a Move package, initializes the pool, and transfers pool ownership to MCMS
// (ownership step 1), so it executes directly with the chain signer (not as an MCMS proposal)
// and returns the deployed pool's AddressRefs.
//
// PoolType accepts the Sui short form ("bnm"/"managed"/"lnr"), the Sui contract-type string
// ("SuiBnMTokenPool"/"SuiManagedTokenPool"/"SuiLnRTokenPool"), or the generic cross-family
// contract-type strings used by EVM/Solana ("BurnMintTokenPool"/"LockReleaseTokenPool"). The
// generic names let a single token_expansion YAML use one poolType across families. The token's
// coin type comes from TokenRef.Address and the symbol from TokenRef.Qualifier; CCIP/MCMS state
// is resolved from the datastore. For managed pools, the first mint-cap ref found for the symbol
// is used.
//
// After initialize this calls transfer_ownership(To: MCMS) EOA-direct, setting a pending
// transfer that UpdateAuthorities' accept_ownership MCMS proposal (step 2) then accepts. The
// final execute_ownership_transfer_to_mcms (step 3) is EOA-direct and handled by the Sui
// cs_mcms_execute_ownership_transfer changeset after the accept proposal lands.
func (a *SuiTokenAdapter) DeployTokenPoolForToken() *cldf_ops.Sequence[tokensapi.DeployTokenPoolInput, sequences.OnChainOutput, cldf_chain.BlockChains] {
	return cldf_ops.NewSequence(
		"sui-adapter:deploy-token-pool-for-token",
		&suideploy.Version1_0_0,
		"Deploy and initialize a Sui token pool for an existing token (direct execution)",
		func(b cldf_ops.Bundle, chains cldf_chain.BlockChains, input tokensapi.DeployTokenPoolInput) (sequences.OnChainOutput, error) {
			chain, ok := chains.SuiChains()[input.ChainSelector]
			if !ok {
				return sequences.OnChainOutput{}, fmt.Errorf("sui chain with selector %d not found", input.ChainSelector)
			}
			if input.TokenRef == nil {
				return sequences.OnChainOutput{}, fmt.Errorf("token ref is required to deploy a sui token pool")
			}
			coinType := input.TokenRef.Address
			symbol := input.TokenRef.Qualifier
			if coinType == "" || symbol == "" {
				return sequences.OnChainOutput{}, fmt.Errorf("token ref must carry the coin type (address) and symbol (qualifier) on chain %d", input.ChainSelector)
			}
			poolType, err := suiPoolTypeFromStr(input.PoolType)
			if err != nil {
				return sequences.OnChainOutput{}, err
			}
			ds := input.ExistingDataStore
			ccipPkg := firstRefAddress(findRefsByType(ds, input.ChainSelector, datastore.ContractType(suideploy.SuiCCIPType)))
			if ccipPkg == "" {
				return sequences.OnChainOutput{}, fmt.Errorf("CCIP package not found on chain %d", input.ChainSelector)
			}
			ccipObjRef := firstRefAddress(findRefsByType(ds, input.ChainSelector, datastore.ContractType(suideploy.SuiCCIPObjectRefType)))
			if ccipObjRef == "" {
				return sequences.OnChainOutput{}, fmt.Errorf("CCIP object ref not found on chain %d", input.ChainSelector)
			}
			mcmsPkg, ok := findRefExcludingLabel(findRefsByType(ds, input.ChainSelector, datastore.ContractType(suideploy.SuiMcmsPackageIDType)), suideploy.MCMSFastCurseLabel)
			if !ok {
				return sequences.OnChainOutput{}, fmt.Errorf("normal (non-fastcurse) MCMS package not found on chain %d", input.ChainSelector)
			}
			fastMcmsPkg, ok := findRefByLabel(findRefsByType(ds, input.ChainSelector, datastore.ContractType(suideploy.SuiMcmsPackageIDType)), suideploy.MCMSFastCurseLabel)
			if !ok {
				return sequences.OnChainOutput{}, fmt.Errorf("fastcurse MCMS package not found on chain %d", input.ChainSelector)
			}
			deployer, err := chain.Signer.GetAddress()
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("failed to get deployer address: %w", err)
			}
			admin := input.RateLimitAdmin
			if admin == "" {
				admin = deployer
			}
			deps := suiDepsExec(chain)
			// Managed-token deploys do not persist CoinMetadata to the datastore (only BnM
			// does), so read the CoinMetadata object id on-chain by coin type as the single
			// reliable source for both pool kinds.
			coinMetaReport, err := cldf_ops.ExecuteOperation(b, coin_ops.GetCoinSymbolOp, deps, coinType)
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("failed to fetch coin metadata for type %s on chain %d: %w", coinType, input.ChainSelector, err)
			}
			coinMeta := coinMetaReport.Output.Id
			if coinMeta == "" {
				return sequences.OnChainOutput{}, fmt.Errorf("coin metadata object id not found for type %s on chain %d", coinType, input.ChainSelector)
			}
			var addresses []datastore.AddressRef
			switch poolType {
			case datastore.ContractType(suideploy.SuiBnMTokenPoolType):
				treasuryCap := refAddress(findRefByLabel(findRefsByType(ds, input.ChainSelector, datastore.ContractType(suideploy.SuiManagedTokenTreasuryCapIDType)), symbol))
				if treasuryCap == "" {
					return sequences.OnChainOutput{}, fmt.Errorf("BnM token treasury cap not found for symbol %s on chain %d", symbol, input.ChainSelector)
				}
				deployReport, err := cldf_ops.ExecuteOperation(b, burnminttokenpoolops.DeployCCIPBurnMintTokenPoolOp, deps, burnminttokenpoolops.BurnMintTokenPoolDeployInput{
					CCIPPackageId:    ccipPkg,
					MCMSAddress:      mcmsPkg.Address,
					FastMcmsAddress:  fastMcmsPkg.Address,
					MCMSOwnerAddress: deployer,
				})
				if err != nil {
					return sequences.OnChainOutput{}, fmt.Errorf("failed to deploy burn-mint token pool: %w", err)
				}
				poolPkg := deployReport.Output.PackageId
				initReport, err := cldf_ops.ExecuteOperation(b, burnminttokenpoolops.BurnMintTokenPoolInitializeOp, deps, burnminttokenpoolops.BurnMintTokenPoolInitializeInput{
					BurnMintPackageId:      poolPkg,
					OwnerCapObjectId:       deployReport.Output.Objects.OwnerCapObjectId,
					CoinObjectTypeArg:      coinType,
					StateObjectId:          ccipObjRef,
					CoinMetadataObjectId:   coinMeta,
					TreasuryCapObjectId:    treasuryCap,
					TokenPoolAdministrator: admin,
				})
				if err != nil {
					return sequences.OnChainOutput{}, fmt.Errorf("failed to initialize burn-mint token pool: %w", err)
				}
				// Ownership step 1 of 3: transfer pool ownership to MCMS (EOA-direct). Sets a pending
				// transfer that UpdateAuthorities' accept_ownership proposal (step 2) then accepts.
				if _, err := cldf_ops.ExecuteOperation(b, burnminttokenpoolops.TransferOwnershipBurnMintTokenPoolOp, deps, burnminttokenpoolops.TransferOwnershipBurnMintTokenPoolInput{
					BurnMintTokenPoolPackageId: poolPkg,
					TypeArgs:                   []string{coinType},
					StateObjectId:              initReport.Output.Objects.StateObjectId,
					OwnerCapObjectId:           deployReport.Output.Objects.OwnerCapObjectId,
					To:                         mcmsPkg.Address,
				}); err != nil {
					return sequences.OnChainOutput{}, fmt.Errorf("failed to transfer burn-mint pool ownership to MCMS: %w", err)
				}
				addresses = appendSuiPoolAddresses(addresses, input.ChainSelector, poolType, symbol, poolPkg, initReport.Output.Objects.StateObjectId, deployReport.Output.Objects.OwnerCapObjectId)
			case datastore.ContractType(suideploy.SuiManagedTokenPoolType):
				managedPkg := refAddress(findRefByLabel(findRefsByType(ds, input.ChainSelector, datastore.ContractType(suideploy.SuiManagedTokenPackageIDType)), symbol))
				mtState := refAddress(findRefByLabel(findRefsByType(ds, input.ChainSelector, datastore.ContractType(suideploy.SuiManagedTokenStateObjectID)), symbol))
				mtOwnerCap := refAddress(findRefByLabel(findRefsByType(ds, input.ChainSelector, datastore.ContractType(suideploy.SuiManagedTokenOwnerCapObjectID)), symbol))
				mintCap := refAddress(findRefByLabel(findRefsByType(ds, input.ChainSelector, datastore.ContractType(suideploy.SuiManagedTokenMinterCapID)), symbol))
				if managedPkg == "" || mtState == "" || mtOwnerCap == "" || mintCap == "" {
					return sequences.OnChainOutput{}, fmt.Errorf("managed token state objects not found for symbol %s on chain %d", symbol, input.ChainSelector)
				}
				deployReport, err := cldf_ops.ExecuteOperation(b, managedtokenpoolops.DeployCCIPManagedTokenPoolOp, deps, managedtokenpoolops.ManagedTokenPoolDeployInput{
					CCIPPackageId:         ccipPkg,
					ManagedTokenPackageId: managedPkg,
					MCMSAddress:           mcmsPkg.Address,
					FastMcmsAddress:       fastMcmsPkg.Address,
					MCMSOwnerAddress:      deployer,
				})
				if err != nil {
					return sequences.OnChainOutput{}, fmt.Errorf("failed to deploy managed token pool: %w", err)
				}
				poolPkg := deployReport.Output.PackageId
				initReport, err := cldf_ops.ExecuteOperation(b, managedtokenpoolops.ManagedTokenPoolInitializeOp, deps, managedtokenpoolops.ManagedTokenPoolInitializeInput{
					ManagedTokenPoolPackageId: poolPkg,
					OwnerCapObjectId:          deployReport.Output.Objects.OwnerCapObjectId,
					CoinObjectTypeArg:         coinType,
					CCIPObjectRefObjectId:     ccipObjRef,
					ManagedTokenStateObjectId: mtState,
					ManagedTokenOwnerCapId:    mtOwnerCap,
					CoinMetadataObjectId:      coinMeta,
					MintCapObjectId:           mintCap,
					TokenPoolAdministrator:    admin,
				})
				if err != nil {
					return sequences.OnChainOutput{}, fmt.Errorf("failed to initialize managed token pool: %w", err)
				}
				// Ownership step 1 of 3: transfer pool ownership to MCMS (EOA-direct). Sets a pending
				// transfer that UpdateAuthorities' accept_ownership proposal (step 2) then accepts.
				if _, err := cldf_ops.ExecuteOperation(b, managedtokenpoolops.TransferOwnershipManagedTokenPoolOp, deps, managedtokenpoolops.TransferOwnershipManagedTokenPoolInput{
					ManagedTokenPoolPackageId: poolPkg,
					TypeArgs:                  []string{coinType},
					StateObjectId:             initReport.Output.Objects.StateObjectId,
					OwnerCapObjectId:          deployReport.Output.Objects.OwnerCapObjectId,
					To:                        mcmsPkg.Address,
				}); err != nil {
					return sequences.OnChainOutput{}, fmt.Errorf("failed to transfer managed pool ownership to MCMS: %w", err)
				}
				addresses = appendSuiPoolAddresses(addresses, input.ChainSelector, poolType, symbol, poolPkg, initReport.Output.Objects.StateObjectId, deployReport.Output.Objects.OwnerCapObjectId)
			case datastore.ContractType(suideploy.SuiLnRTokenPoolType):
				treasuryCap := refAddress(findRefByLabel(findRefsByType(ds, input.ChainSelector, datastore.ContractType(suideploy.SuiLnRTokenTreasuryCapIDType)), symbol))
				if treasuryCap == "" {
					return sequences.OnChainOutput{}, fmt.Errorf("token treasury cap not found for symbol %s on chain %d", symbol, input.ChainSelector)
				}
				deployReport, err := cldf_ops.ExecuteOperation(b, lockreleasetokenpoolops.DeployCCIPLockReleaseTokenPoolOp, deps, lockreleasetokenpoolops.LockReleaseTokenPoolDeployInput{
					CCIPPackageId:    ccipPkg,
					MCMSAddress:      mcmsPkg.Address,
					FastMcmsAddress:  fastMcmsPkg.Address,
					MCMSOwnerAddress: deployer,
				})
				if err != nil {
					return sequences.OnChainOutput{}, fmt.Errorf("failed to deploy lock-release token pool: %w", err)
				}
				poolPkg := deployReport.Output.PackageId
				initReport, err := cldf_ops.ExecuteOperation(b, lockreleasetokenpoolops.LockReleaseTokenPoolInitializeOp, deps, lockreleasetokenpoolops.LockReleaseTokenPoolInitializeInput{
					LockReleasePackageId:   poolPkg,
					OwnerCapObjectId:       deployReport.Output.Objects.OwnerCapObjectId,
					CoinObjectTypeArg:      coinType,
					StateObjectId:          ccipObjRef,
					CoinMetadataObjectId:   coinMeta,
					TreasuryCapObjectId:    treasuryCap,
					TokenPoolAdministrator: admin,
					Rebalancer:             deployer,
				})
				if err != nil {
					return sequences.OnChainOutput{}, fmt.Errorf("failed to initialize lock-release token pool: %w", err)
				}
				// Ownership step 1 of 3: transfer pool ownership to MCMS (EOA-direct). Sets a pending
				// transfer that UpdateAuthorities' accept_ownership proposal (step 2) then accepts.
				if _, err := cldf_ops.ExecuteOperation(b, lockreleasetokenpoolops.TransferOwnershipLockReleaseTokenPoolOp, deps, lockreleasetokenpoolops.TransferOwnershipLockReleaseTokenPoolInput{
					LockReleaseTokenPoolPackageId: poolPkg,
					TypeArgs:                      []string{coinType},
					StateObjectId:                 initReport.Output.Objects.StateObjectId,
					OwnerCapObjectId:              deployReport.Output.Objects.OwnerCapObjectId,
					To:                            mcmsPkg.Address,
				}); err != nil {
					return sequences.OnChainOutput{}, fmt.Errorf("failed to transfer lock-release pool ownership to MCMS: %w", err)
				}
				addresses = appendSuiPoolAddresses(addresses, input.ChainSelector, poolType, symbol, poolPkg, initReport.Output.Objects.StateObjectId, deployReport.Output.Objects.OwnerCapObjectId)
			default:
				return sequences.OnChainOutput{}, fmt.Errorf("unsupported sui pool type %s for DeployTokenPoolForToken", poolType)
			}
			return sequences.OnChainOutput{Addresses: addresses}, nil
		},
	)
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

// suiDepsExec builds OpTxDeps for direct execution with the chain signer. Used by deploy/init
// ops that publish packages and must execute immediately rather than via MCMS.
func suiDepsExec(chain cldfsui.Chain) sui_ops.OpTxDeps {
	gas := uint64(400_000_000)
	return sui_ops.OpTxDeps{
		Client: chain.Client,
		Signer: chain.Signer,
		GetCallOpts: func() *bind.CallOpts {
			return &bind.CallOpts{WaitForExecution: true, GasBudget: &gas}
		},
		SuiRPC: chain.URL,
	}
}

// suiPoolOwner reads the pool's ownable owner on-chain via DevInspect to decide
// configure-before-own: a deployer-owned pool can run OwnerCap-gated config directly.
func suiPoolOwner(b cldf_ops.Bundle, chain cldfsui.Chain, poolType datastore.ContractType,
	pkgID, coinType, stateObjID string,
) (string, error) {
	opts := &bind.CallOpts{Signer: chain.Signer}
	state := bind.Object{Id: stateObjID}
	typeArgs := []string{coinType}
	switch poolType {
	case datastore.ContractType(suideploy.SuiBnMTokenPoolType):
		contract, err := module_burn_mint_token_pool.NewBurnMintTokenPool(pkgID, chain.Client)
		if err != nil {
			return "", fmt.Errorf("failed to create burn-mint token pool contract: %w", err)
		}
		return contract.DevInspect().Owner(b.GetContext(), opts, typeArgs, state)
	case datastore.ContractType(suideploy.SuiManagedTokenPoolType):
		contract, err := module_managed_token_pool.NewManagedTokenPool(pkgID, chain.Client)
		if err != nil {
			return "", fmt.Errorf("failed to create managed token pool contract: %w", err)
		}
		return contract.DevInspect().Owner(b.GetContext(), opts, typeArgs, state)
	case datastore.ContractType(suideploy.SuiLnRTokenPoolType):
		contract, err := module_lock_release_token_pool.NewLockReleaseTokenPool(pkgID, chain.Client)
		if err != nil {
			return "", fmt.Errorf("failed to create lock-release token pool contract: %w", err)
		}
		return contract.DevInspect().Owner(b.GetContext(), opts, typeArgs, state)
	default:
		return "", fmt.Errorf("unsupported sui token pool type %s for owner read", poolType)
	}
}

// suiAddrEqual compares two Sui addresses after normalizing to lowercase with a 0x prefix.
func suiAddrEqual(a, b string) bool {
	return normalizeSuiAddr(a) == normalizeSuiAddr(b)
}

func normalizeSuiAddr(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)
	if s != "" && !strings.HasPrefix(s, "0x") {
		s = "0x" + s
	}
	return s
}

// suiPoolTypeFromStr maps a PoolType string to a Sui pool contract type, accepting both the
// Sui short form ("bnm"/"managed"/"lnr") and the contract-type string.
func suiPoolTypeFromStr(s string) (datastore.ContractType, error) {
	switch s {
	case "bnm", string(suideploy.SuiBnMTokenPoolType), string(cciputils.BurnMintTokenPool):
		return datastore.ContractType(suideploy.SuiBnMTokenPoolType), nil
	case "managed", string(suideploy.SuiManagedTokenPoolType):
		return datastore.ContractType(suideploy.SuiManagedTokenPoolType), nil
	case "lnr", string(suideploy.SuiLnRTokenPoolType), string(cciputils.LockReleaseTokenPool):
		return datastore.ContractType(suideploy.SuiLnRTokenPoolType), nil
	default:
		return "", fmt.Errorf("unsupported sui pool type %q", s)
	}
}

// firstRefAddress returns the address of the first ref in a slice, or empty if there are none.
func firstRefAddress(refs []datastore.AddressRef) string {
	if len(refs) == 0 {
		return ""
	}
	return refs[0].Address
}

// refAddress returns the address of a found ref, or empty if ok is false.
func refAddress(ref datastore.AddressRef, ok bool) string {
	if !ok {
		return ""
	}
	return ref.Address
}

// appendSuiPoolAddresses builds the AddressRefs for a freshly deployed pool (package, state,
// owner cap) keyed by the token symbol label, matching the refs saved by the Sui
// DeployTPAndConfigure changeset.
func appendSuiPoolAddresses(in []datastore.AddressRef, selector uint64, poolType datastore.ContractType, symbol, poolPkg, stateObjID, ownerCapID string) []datastore.AddressRef {
	var stateType, ownerType datastore.ContractType
	switch poolType {
	case datastore.ContractType(suideploy.SuiBnMTokenPoolType):
		stateType = datastore.ContractType(suideploy.SuiBnMTokenPoolStateType)
		ownerType = datastore.ContractType(suideploy.SuiBnMTokenPoolOwnerIDType)
	case datastore.ContractType(suideploy.SuiManagedTokenPoolType):
		stateType = datastore.ContractType(suideploy.SuiManagedTokenPoolStateType)
		ownerType = datastore.ContractType(suideploy.SuiManagedTokenPoolOwnerIDType)
	case datastore.ContractType(suideploy.SuiLnRTokenPoolType):
		stateType = datastore.ContractType(suideploy.SuiLnRTokenPoolStateType)
		ownerType = datastore.ContractType(suideploy.SuiLnRTokenPoolOwnerIDType)
	}
	version := semver.MustParse("1.0.0")
	return append(in,
		datastore.AddressRef{ChainSelector: selector, Type: poolType, Address: poolPkg, Version: version, Qualifier: fmt.Sprintf("%s-%s", poolPkg, poolType), Labels: datastore.NewLabelSet(symbol)},
		datastore.AddressRef{ChainSelector: selector, Type: stateType, Address: stateObjID, Version: version, Qualifier: fmt.Sprintf("%s-%s", stateObjID, stateType), Labels: datastore.NewLabelSet(symbol)},
		datastore.AddressRef{ChainSelector: selector, Type: ownerType, Address: ownerCapID, Version: version, Qualifier: fmt.Sprintf("%s-%s", ownerCapID, ownerType), Labels: datastore.NewLabelSet(symbol)},
	)
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

// batchOpFromCalls bridges one or more Sui TransactionCalls into a single MCMS BatchOperation
// with ordered transactions. Returns an empty output when there are no calls.
func batchOpFromCalls(chainSelector uint64, calls []sui_ops.TransactionCall) (sequences.OnChainOutput, error) {
	if len(calls) == 0 {
		return sequences.OnChainOutput{}, nil
	}
	txs := make([]mcmstypes.Transaction, 0, len(calls))
	for _, c := range calls {
		tx, err := suideployutils.TransactionCallToMCMSTransaction(c)
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("failed to convert sui call to mcms transaction: %w", err)
		}
		txs = append(txs, tx)
	}
	return sequences.OnChainOutput{
		BatchOps: []mcmstypes.BatchOperation{{
			ChainSelector: mcmstypes.ChainSelector(chainSelector),
			Transactions:  txs,
		}},
	}, nil
}

// suiDecimals reads token decimals from on-chain coin metadata for a coin type.
func suiDecimals(b cldf_ops.Bundle, chain cldfsui.Chain, coinType string) (uint8, error) {
	report, err := cldf_ops.ExecuteOperation(b, coin_ops.GetCoinSymbolOp, suiDeps(chain), coinType)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch coin metadata for type %s: %w", coinType, err)
	}
	if report.Output.Decimals < 0 || report.Output.Decimals > 255 {
		return 0, fmt.Errorf("invalid decimals %d for coin type %s", report.Output.Decimals, coinType)
	}
	return uint8(report.Output.Decimals), nil
}

// deriveSuiCoinType builds the coin type for a token symbol by locating the token's coin package
// ref in the datastore. The coin package ref carries the coin's module::STRUCT suffix as a
// coinType= label (set by whichever changeset saved the coin package id), so the coin type is
// normalizeCoinType(ref.Address) + "::" + that suffix. The ref is matched by symbol label plus the
// coinType label, independent of the contract type it is filed under, so every token — BnM, LINK,
// or any newly added coin behind a BurnMint or Managed pool — derives through the same path.
func deriveSuiCoinType(ds datastore.DataStore, selector uint64, symbol string) (string, error) {
	if symbol == "" {
		return "", fmt.Errorf("symbol is required to derive the sui coin type")
	}
	if ds == nil {
		return "", fmt.Errorf("could not derive sui coin type for symbol %s on chain %d: datastore is nil", symbol, selector)
	}
	for _, ref := range ds.Addresses().Filter(datastore.AddressRefByChainSelector(selector)) {
		if !ref.Labels.Contains(symbol) {
			continue
		}
		suffix := coinTypeSuffixFromLabels(ref)
		if suffix == "" {
			continue
		}
		return normalizeCoinType(ref.Address) + "::" + suffix, nil
	}
	return "", fmt.Errorf("could not derive sui coin type for symbol %s on chain %d: no coin package ref carrying a coinType= label found for that symbol", symbol, selector)
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

// suiPoolSymbol returns the token symbol used to match pool state/owner-cap and token-package
// refs by label. The generic datastore-first resolver returns a raw ref whose Qualifier is the
// synthetic "<addr>-<type>" rather than the symbol, while the symbol is reliably carried as the
// first label, so prefer labels and fall back to the qualifier only when no label is set.
func suiPoolSymbol(ref datastore.AddressRef) string {
	if s := symbolFromLabels(ref); s != "" {
		return s
	}
	return ref.Qualifier
}

// suiCoinTypeLabelPrefix marks a label carrying a coin package's module::STRUCT suffix. The coin
// deploy changesets add it so deriveSuiCoinType can build the correct coin type for tokens whose
// suffix is not one of the hardcoded BnM/LINK fallbacks.
const suiCoinTypeLabelPrefix = "coinType="

// coinTypeSuffixFromLabels returns the module::STRUCT suffix stored in a coinType= label on the
// ref, or "" when no such label is set.
func coinTypeSuffixFromLabels(r datastore.AddressRef) string {
	for _, l := range r.Labels.List() {
		if strings.HasPrefix(l, suiCoinTypeLabelPrefix) {
			return strings.TrimPrefix(l, suiCoinTypeLabelPrefix)
		}
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
	case strings.Contains(coinType, "::ccip_lock_release_token::"):
		return datastore.ContractType(suideploy.SuiLnRTokenType)
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
