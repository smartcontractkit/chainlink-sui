// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package module_managed_token_pool

import (
	"context"
	"fmt"
	"math/big"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/mystenbcs"
	"github.com/block-vision/sui-go-sdk/sui"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
)

var (
	_ = big.NewInt
)

const FunctionInfo = `[{"package":"managed_token_pool","module":"managed_token_pool","name":"accept_ownership","parameters":[{"name":"state","type":"ManagedTokenPoolState<T>"}]},{"package":"managed_token_pool","module":"managed_token_pool","name":"accept_ownership_from_object","parameters":[{"name":"state","type":"ManagedTokenPoolState<T>"},{"name":"from","type":"sui::object::UID"}]},{"package":"managed_token_pool","module":"managed_token_pool","name":"add_remote_pool","parameters":[{"name":"state","type":"ManagedTokenPoolState<T>"},{"name":"owner_cap","type":"OwnerCap"},{"name":"remote_chain_selector","type":"u64"},{"name":"remote_pool_address","type":"vector<u8>"}]},{"package":"managed_token_pool","module":"managed_token_pool","name":"apply_allowlist_updates","parameters":[{"name":"state","type":"ManagedTokenPoolState<T>"},{"name":"owner_cap","type":"OwnerCap"},{"name":"removes","type":"vector<address>"},{"name":"adds","type":"vector<address>"}]},{"package":"managed_token_pool","module":"managed_token_pool","name":"apply_chain_updates","parameters":[{"name":"state","type":"ManagedTokenPoolState<T>"},{"name":"owner_cap","type":"OwnerCap"},{"name":"remote_chain_selectors_to_remove","type":"vector<u64>"},{"name":"remote_chain_selectors_to_add","type":"vector<u64>"},{"name":"remote_pool_addresses_to_add","type":"vector<vector<vector<u8>>>"},{"name":"remote_token_addresses_to_add","type":"vector<vector<u8>>"}]},{"package":"managed_token_pool","module":"managed_token_pool","name":"destroy_token_pool","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"state","type":"ManagedTokenPoolState<T>"},{"name":"owner_cap","type":"OwnerCap"}]},{"package":"managed_token_pool","module":"managed_token_pool","name":"execute_ownership_transfer","parameters":[{"name":"owner_cap","type":"OwnerCap"},{"name":"state","type":"ManagedTokenPoolState<T>"},{"name":"to","type":"address"}]},{"package":"managed_token_pool","module":"managed_token_pool","name":"execute_ownership_transfer_to_mcms","parameters":[{"name":"owner_cap","type":"OwnerCap"},{"name":"state","type":"ManagedTokenPoolState<T>"},{"name":"registry","type":"Registry"},{"name":"to","type":"address"}]},{"package":"managed_token_pool","module":"managed_token_pool","name":"get_allowlist","parameters":[{"name":"state","type":"ManagedTokenPoolState<T>"}]},{"package":"managed_token_pool","module":"managed_token_pool","name":"get_allowlist_enabled","parameters":[{"name":"state","type":"ManagedTokenPoolState<T>"}]},{"package":"managed_token_pool","module":"managed_token_pool","name":"get_current_inbound_rate_limiter_state","parameters":[{"name":"clock","type":"Clock"},{"name":"state","type":"ManagedTokenPoolState<T>"},{"name":"remote_chain_selector","type":"u64"}]},{"package":"managed_token_pool","module":"managed_token_pool","name":"get_current_outbound_rate_limiter_state","parameters":[{"name":"clock","type":"Clock"},{"name":"state","type":"ManagedTokenPoolState<T>"},{"name":"remote_chain_selector","type":"u64"}]},{"package":"managed_token_pool","module":"managed_token_pool","name":"get_remote_pools","parameters":[{"name":"state","type":"ManagedTokenPoolState<T>"},{"name":"remote_chain_selector","type":"u64"}]},{"package":"managed_token_pool","module":"managed_token_pool","name":"get_remote_token","parameters":[{"name":"state","type":"ManagedTokenPoolState<T>"},{"name":"remote_chain_selector","type":"u64"}]},{"package":"managed_token_pool","module":"managed_token_pool","name":"get_supported_chains","parameters":[{"name":"state","type":"ManagedTokenPoolState<T>"}]},{"package":"managed_token_pool","module":"managed_token_pool","name":"get_token","parameters":[{"name":"state","type":"ManagedTokenPoolState<T>"}]},{"package":"managed_token_pool","module":"managed_token_pool","name":"get_token_decimals","parameters":[{"name":"state","type":"ManagedTokenPoolState<T>"}]},{"package":"managed_token_pool","module":"managed_token_pool","name":"get_token_symbol","parameters":[{"name":"state","type":"ManagedTokenPoolState<T>"}]},{"package":"managed_token_pool","module":"managed_token_pool","name":"has_pending_transfer","parameters":[{"name":"state","type":"ManagedTokenPoolState<T>"}]},{"package":"managed_token_pool","module":"managed_token_pool","name":"initialize_with_managed_token","parameters":[{"name":"owner_cap","type":"OwnerCap"},{"name":"ref","type":"CCIPObjectRef"},{"name":"managed_token_state","type":"TokenState<T>"},{"name":"managed_token_owner_cap","type":"ManagedTokenOwnerCap<T>"},{"name":"coin_metadata","type":"CoinMetadata<T>"},{"name":"mint_cap","type":"MintCap<T>"},{"name":"token_pool_administrator","type":"address"}]},{"package":"managed_token_pool","module":"managed_token_pool","name":"is_remote_pool","parameters":[{"name":"state","type":"ManagedTokenPoolState<T>"},{"name":"remote_chain_selector","type":"u64"},{"name":"remote_pool_address","type":"vector<u8>"}]},{"package":"managed_token_pool","module":"managed_token_pool","name":"is_supported_chain","parameters":[{"name":"state","type":"ManagedTokenPoolState<T>"},{"name":"remote_chain_selector","type":"u64"}]},{"package":"managed_token_pool","module":"managed_token_pool","name":"lock_or_burn","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"token_transfer_params","type":"onramp_sh::TokenTransferParams"},{"name":"c","type":"Coin<T>"},{"name":"remote_chain_selector","type":"u64"},{"name":"clock","type":"Clock"},{"name":"deny_list","type":"DenyList"},{"name":"token_state","type":"TokenState<T>"},{"name":"state","type":"ManagedTokenPoolState<T>"}]},{"package":"managed_token_pool","module":"managed_token_pool","name":"owner","parameters":[{"name":"state","type":"ManagedTokenPoolState<T>"}]},{"package":"managed_token_pool","module":"managed_token_pool","name":"pending_transfer_accepted","parameters":[{"name":"state","type":"ManagedTokenPoolState<T>"}]},{"package":"managed_token_pool","module":"managed_token_pool","name":"pending_transfer_from","parameters":[{"name":"state","type":"ManagedTokenPoolState<T>"}]},{"package":"managed_token_pool","module":"managed_token_pool","name":"pending_transfer_to","parameters":[{"name":"state","type":"ManagedTokenPoolState<T>"}]},{"package":"managed_token_pool","module":"managed_token_pool","name":"release_or_mint","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"receiver_params","type":"offramp_sh::ReceiverParams"},{"name":"clock","type":"Clock"},{"name":"deny_list","type":"DenyList"},{"name":"token_state","type":"TokenState<T>"},{"name":"state","type":"ManagedTokenPoolState<T>"}]},{"package":"managed_token_pool","module":"managed_token_pool","name":"remove_remote_pool","parameters":[{"name":"state","type":"ManagedTokenPoolState<T>"},{"name":"owner_cap","type":"OwnerCap"},{"name":"remote_chain_selector","type":"u64"},{"name":"remote_pool_address","type":"vector<u8>"}]},{"package":"managed_token_pool","module":"managed_token_pool","name":"set_allowlist_enabled","parameters":[{"name":"state","type":"ManagedTokenPoolState<T>"},{"name":"owner_cap","type":"OwnerCap"},{"name":"enabled","type":"bool"}]},{"package":"managed_token_pool","module":"managed_token_pool","name":"set_chain_rate_limiter_config","parameters":[{"name":"state","type":"ManagedTokenPoolState<T>"},{"name":"owner_cap","type":"OwnerCap"},{"name":"clock","type":"Clock"},{"name":"remote_chain_selector","type":"u64"},{"name":"outbound_is_enabled","type":"bool"},{"name":"outbound_capacity","type":"u64"},{"name":"outbound_rate","type":"u64"},{"name":"inbound_is_enabled","type":"bool"},{"name":"inbound_capacity","type":"u64"},{"name":"inbound_rate","type":"u64"}]},{"package":"managed_token_pool","module":"managed_token_pool","name":"set_chain_rate_limiter_configs","parameters":[{"name":"state","type":"ManagedTokenPoolState<T>"},{"name":"owner_cap","type":"OwnerCap"},{"name":"clock","type":"Clock"},{"name":"remote_chain_selectors","type":"vector<u64>"},{"name":"outbound_is_enableds","type":"vector<bool>"},{"name":"outbound_capacities","type":"vector<u64>"},{"name":"outbound_rates","type":"vector<u64>"},{"name":"inbound_is_enableds","type":"vector<bool>"},{"name":"inbound_capacities","type":"vector<u64>"},{"name":"inbound_rates","type":"vector<u64>"}]},{"package":"managed_token_pool","module":"managed_token_pool","name":"transfer_ownership","parameters":[{"name":"state","type":"ManagedTokenPoolState<T>"},{"name":"owner_cap","type":"OwnerCap"},{"name":"new_owner","type":"address"}]},{"package":"managed_token_pool","module":"managed_token_pool","name":"type_and_version","parameters":null}]`

type IManagedTokenPool interface {
	TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	InitializeWithManagedToken(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ownerCap bind.Object, ref bind.Object, managedTokenState bind.Object, managedTokenOwnerCap bind.Object, coinMetadata bind.Object, mintCap bind.Object, tokenPoolAdministrator string) (*models.SuiTransactionBlockResponse, error)
	AddRemotePool(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object, remoteChainSelector uint64, remotePoolAddress []byte) (*models.SuiTransactionBlockResponse, error)
	RemoveRemotePool(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object, remoteChainSelector uint64, remotePoolAddress []byte) (*models.SuiTransactionBlockResponse, error)
	IsSupportedChain(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	GetSupportedChains(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	ApplyChainUpdates(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object, remoteChainSelectorsToRemove []uint64, remoteChainSelectorsToAdd []uint64, remotePoolAddressesToAdd [][][]byte, remoteTokenAddressesToAdd [][]byte) (*models.SuiTransactionBlockResponse, error)
	GetAllowlistEnabled(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	GetAllowlist(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	SetAllowlistEnabled(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object, enabled bool) (*models.SuiTransactionBlockResponse, error)
	ApplyAllowlistUpdates(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object, removes []string, adds []string) (*models.SuiTransactionBlockResponse, error)
	GetToken(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	GetTokenDecimals(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	GetTokenSymbol(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	GetRemotePools(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	IsRemotePool(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, remoteChainSelector uint64, remotePoolAddress []byte) (*models.SuiTransactionBlockResponse, error)
	GetRemoteToken(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	LockOrBurn(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, tokenTransferParams bind.Object, c_ bind.Object, remoteChainSelector uint64, clock bind.Object, denyList bind.Object, tokenState bind.Object, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	ReleaseOrMint(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, receiverParams bind.Object, clock bind.Object, denyList bind.Object, tokenState bind.Object, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	SetChainRateLimiterConfigs(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object, clock bind.Object, remoteChainSelectors []uint64, outboundIsEnableds []bool, outboundCapacities []uint64, outboundRates []uint64, inboundIsEnableds []bool, inboundCapacities []uint64, inboundRates []uint64) (*models.SuiTransactionBlockResponse, error)
	SetChainRateLimiterConfig(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object, clock bind.Object, remoteChainSelector uint64, outboundIsEnabled bool, outboundCapacity uint64, outboundRate uint64, inboundIsEnabled bool, inboundCapacity uint64, inboundRate uint64) (*models.SuiTransactionBlockResponse, error)
	GetCurrentInboundRateLimiterState(ctx context.Context, opts *bind.CallOpts, typeArgs []string, clock bind.Object, state bind.Object, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	GetCurrentOutboundRateLimiterState(ctx context.Context, opts *bind.CallOpts, typeArgs []string, clock bind.Object, state bind.Object, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	Owner(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	HasPendingTransfer(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	PendingTransferFrom(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	PendingTransferTo(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	PendingTransferAccepted(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	TransferOwnership(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object, newOwner string) (*models.SuiTransactionBlockResponse, error)
	AcceptOwnership(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	AcceptOwnershipFromObject(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, from string) (*models.SuiTransactionBlockResponse, error)
	McmsAcceptOwnership(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	ExecuteOwnershipTransfer(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ownerCap bind.Object, state bind.Object, to string) (*models.SuiTransactionBlockResponse, error)
	ExecuteOwnershipTransferToMcms(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ownerCap bind.Object, state bind.Object, registry bind.Object, to string) (*models.SuiTransactionBlockResponse, error)
	McmsRegisterUpgradeCap(ctx context.Context, opts *bind.CallOpts, upgradeCap bind.Object, registry bind.Object, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsSetAllowlistEnabled(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsApplyAllowlistUpdates(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsApplyChainUpdates(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsAddRemotePool(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsRemoveRemotePool(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsSetChainRateLimiterConfigs(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, registry bind.Object, params bind.Object, clock bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsSetChainRateLimiterConfig(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, registry bind.Object, params bind.Object, clock bind.Object) (*models.SuiTransactionBlockResponse, error)
	DestroyTokenPool(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, state bind.Object, ownerCap bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsDestroyTokenPool(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsTransferOwnership(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsExecuteOwnershipTransfer(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, registry bind.Object, deployerState bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsAddAllowedModules(ctx context.Context, opts *bind.CallOpts, typeArgs []string, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsRemoveAllowedModules(ctx context.Context, opts *bind.CallOpts, typeArgs []string, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	DevInspect() IManagedTokenPoolDevInspect
	Encoder() ManagedTokenPoolEncoder
	Bound() bind.IBoundContract
}

type IManagedTokenPoolDevInspect interface {
	TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (string, error)
	IsSupportedChain(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, remoteChainSelector uint64) (bool, error)
	GetSupportedChains(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) ([]uint64, error)
	GetAllowlistEnabled(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (bool, error)
	GetAllowlist(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) ([]string, error)
	GetToken(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (string, error)
	GetTokenDecimals(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (byte, error)
	GetTokenSymbol(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (string, error)
	GetRemotePools(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, remoteChainSelector uint64) ([][]byte, error)
	IsRemotePool(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, remoteChainSelector uint64, remotePoolAddress []byte) (bool, error)
	GetRemoteToken(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, remoteChainSelector uint64) ([]byte, error)
	GetCurrentInboundRateLimiterState(ctx context.Context, opts *bind.CallOpts, typeArgs []string, clock bind.Object, state bind.Object, remoteChainSelector uint64) (TokenBucketWrapper, error)
	GetCurrentOutboundRateLimiterState(ctx context.Context, opts *bind.CallOpts, typeArgs []string, clock bind.Object, state bind.Object, remoteChainSelector uint64) (TokenBucketWrapper, error)
	Owner(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (string, error)
	HasPendingTransfer(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (bool, error)
	PendingTransferFrom(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*string, error)
	PendingTransferTo(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*string, error)
	PendingTransferAccepted(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*bool, error)
	DestroyTokenPool(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, state bind.Object, ownerCap bind.Object) (any, error)
}

type ManagedTokenPoolEncoder interface {
	TypeAndVersion() (*bind.EncodedCall, error)
	TypeAndVersionWithArgs(args ...any) (*bind.EncodedCall, error)
	InitializeWithManagedToken(typeArgs []string, ownerCap bind.Object, ref bind.Object, managedTokenState bind.Object, managedTokenOwnerCap bind.Object, coinMetadata bind.Object, mintCap bind.Object, tokenPoolAdministrator string) (*bind.EncodedCall, error)
	InitializeWithManagedTokenWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	AddRemotePool(typeArgs []string, state bind.Object, ownerCap bind.Object, remoteChainSelector uint64, remotePoolAddress []byte) (*bind.EncodedCall, error)
	AddRemotePoolWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	RemoveRemotePool(typeArgs []string, state bind.Object, ownerCap bind.Object, remoteChainSelector uint64, remotePoolAddress []byte) (*bind.EncodedCall, error)
	RemoveRemotePoolWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	IsSupportedChain(typeArgs []string, state bind.Object, remoteChainSelector uint64) (*bind.EncodedCall, error)
	IsSupportedChainWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	GetSupportedChains(typeArgs []string, state bind.Object) (*bind.EncodedCall, error)
	GetSupportedChainsWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	ApplyChainUpdates(typeArgs []string, state bind.Object, ownerCap bind.Object, remoteChainSelectorsToRemove []uint64, remoteChainSelectorsToAdd []uint64, remotePoolAddressesToAdd [][][]byte, remoteTokenAddressesToAdd [][]byte) (*bind.EncodedCall, error)
	ApplyChainUpdatesWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	GetAllowlistEnabled(typeArgs []string, state bind.Object) (*bind.EncodedCall, error)
	GetAllowlistEnabledWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	GetAllowlist(typeArgs []string, state bind.Object) (*bind.EncodedCall, error)
	GetAllowlistWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	SetAllowlistEnabled(typeArgs []string, state bind.Object, ownerCap bind.Object, enabled bool) (*bind.EncodedCall, error)
	SetAllowlistEnabledWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	ApplyAllowlistUpdates(typeArgs []string, state bind.Object, ownerCap bind.Object, removes []string, adds []string) (*bind.EncodedCall, error)
	ApplyAllowlistUpdatesWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	GetToken(typeArgs []string, state bind.Object) (*bind.EncodedCall, error)
	GetTokenWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	GetTokenDecimals(typeArgs []string, state bind.Object) (*bind.EncodedCall, error)
	GetTokenDecimalsWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	GetTokenSymbol(typeArgs []string, state bind.Object) (*bind.EncodedCall, error)
	GetTokenSymbolWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	GetRemotePools(typeArgs []string, state bind.Object, remoteChainSelector uint64) (*bind.EncodedCall, error)
	GetRemotePoolsWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	IsRemotePool(typeArgs []string, state bind.Object, remoteChainSelector uint64, remotePoolAddress []byte) (*bind.EncodedCall, error)
	IsRemotePoolWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	GetRemoteToken(typeArgs []string, state bind.Object, remoteChainSelector uint64) (*bind.EncodedCall, error)
	GetRemoteTokenWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	LockOrBurn(typeArgs []string, ref bind.Object, tokenTransferParams bind.Object, c_ bind.Object, remoteChainSelector uint64, clock bind.Object, denyList bind.Object, tokenState bind.Object, state bind.Object) (*bind.EncodedCall, error)
	LockOrBurnWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	ReleaseOrMint(typeArgs []string, ref bind.Object, receiverParams bind.Object, clock bind.Object, denyList bind.Object, tokenState bind.Object, state bind.Object) (*bind.EncodedCall, error)
	ReleaseOrMintWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	SetChainRateLimiterConfigs(typeArgs []string, state bind.Object, ownerCap bind.Object, clock bind.Object, remoteChainSelectors []uint64, outboundIsEnableds []bool, outboundCapacities []uint64, outboundRates []uint64, inboundIsEnableds []bool, inboundCapacities []uint64, inboundRates []uint64) (*bind.EncodedCall, error)
	SetChainRateLimiterConfigsWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	SetChainRateLimiterConfig(typeArgs []string, state bind.Object, ownerCap bind.Object, clock bind.Object, remoteChainSelector uint64, outboundIsEnabled bool, outboundCapacity uint64, outboundRate uint64, inboundIsEnabled bool, inboundCapacity uint64, inboundRate uint64) (*bind.EncodedCall, error)
	SetChainRateLimiterConfigWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	GetCurrentInboundRateLimiterState(typeArgs []string, clock bind.Object, state bind.Object, remoteChainSelector uint64) (*bind.EncodedCall, error)
	GetCurrentInboundRateLimiterStateWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	GetCurrentOutboundRateLimiterState(typeArgs []string, clock bind.Object, state bind.Object, remoteChainSelector uint64) (*bind.EncodedCall, error)
	GetCurrentOutboundRateLimiterStateWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	Owner(typeArgs []string, state bind.Object) (*bind.EncodedCall, error)
	OwnerWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	HasPendingTransfer(typeArgs []string, state bind.Object) (*bind.EncodedCall, error)
	HasPendingTransferWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	PendingTransferFrom(typeArgs []string, state bind.Object) (*bind.EncodedCall, error)
	PendingTransferFromWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	PendingTransferTo(typeArgs []string, state bind.Object) (*bind.EncodedCall, error)
	PendingTransferToWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	PendingTransferAccepted(typeArgs []string, state bind.Object) (*bind.EncodedCall, error)
	PendingTransferAcceptedWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	TransferOwnership(typeArgs []string, state bind.Object, ownerCap bind.Object, newOwner string) (*bind.EncodedCall, error)
	TransferOwnershipWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	AcceptOwnership(typeArgs []string, state bind.Object) (*bind.EncodedCall, error)
	AcceptOwnershipWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	AcceptOwnershipFromObject(typeArgs []string, state bind.Object, from string) (*bind.EncodedCall, error)
	AcceptOwnershipFromObjectWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	McmsAcceptOwnership(typeArgs []string, state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsAcceptOwnershipWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	ExecuteOwnershipTransfer(typeArgs []string, ownerCap bind.Object, state bind.Object, to string) (*bind.EncodedCall, error)
	ExecuteOwnershipTransferWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	ExecuteOwnershipTransferToMcms(typeArgs []string, ownerCap bind.Object, state bind.Object, registry bind.Object, to string) (*bind.EncodedCall, error)
	ExecuteOwnershipTransferToMcmsWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	McmsRegisterUpgradeCap(upgradeCap bind.Object, registry bind.Object, state bind.Object) (*bind.EncodedCall, error)
	McmsRegisterUpgradeCapWithArgs(args ...any) (*bind.EncodedCall, error)
	McmsSetAllowlistEnabled(typeArgs []string, state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsSetAllowlistEnabledWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	McmsApplyAllowlistUpdates(typeArgs []string, state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsApplyAllowlistUpdatesWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	McmsApplyChainUpdates(typeArgs []string, state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsApplyChainUpdatesWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	McmsAddRemotePool(typeArgs []string, state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsAddRemotePoolWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	McmsRemoveRemotePool(typeArgs []string, state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsRemoveRemotePoolWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	McmsSetChainRateLimiterConfigs(typeArgs []string, state bind.Object, registry bind.Object, params bind.Object, clock bind.Object) (*bind.EncodedCall, error)
	McmsSetChainRateLimiterConfigsWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	McmsSetChainRateLimiterConfig(typeArgs []string, state bind.Object, registry bind.Object, params bind.Object, clock bind.Object) (*bind.EncodedCall, error)
	McmsSetChainRateLimiterConfigWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	DestroyTokenPool(typeArgs []string, ref bind.Object, state bind.Object, ownerCap bind.Object) (*bind.EncodedCall, error)
	DestroyTokenPoolWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	McmsDestroyTokenPool(typeArgs []string, ref bind.Object, state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsDestroyTokenPoolWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	McmsTransferOwnership(typeArgs []string, state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsTransferOwnershipWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	McmsExecuteOwnershipTransfer(typeArgs []string, state bind.Object, registry bind.Object, deployerState bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsExecuteOwnershipTransferWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	McmsAddAllowedModules(typeArgs []string, registry bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsAddAllowedModulesWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	McmsRemoveAllowedModules(typeArgs []string, registry bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsRemoveAllowedModulesWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
}

type ManagedTokenPoolContract struct {
	*bind.BoundContract
	managedTokenPoolEncoder
	devInspect *ManagedTokenPoolDevInspect
}

type ManagedTokenPoolDevInspect struct {
	contract *ManagedTokenPoolContract
}

var _ IManagedTokenPool = (*ManagedTokenPoolContract)(nil)
var _ IManagedTokenPoolDevInspect = (*ManagedTokenPoolDevInspect)(nil)

func NewManagedTokenPool(packageID string, client sui.ISuiAPI) (IManagedTokenPool, error) {
	contract, err := bind.NewBoundContract(packageID, "managed_token_pool", "managed_token_pool", client)
	if err != nil {
		return nil, err
	}

	c := &ManagedTokenPoolContract{
		BoundContract:           contract,
		managedTokenPoolEncoder: managedTokenPoolEncoder{BoundContract: contract},
	}
	c.devInspect = &ManagedTokenPoolDevInspect{contract: c}
	return c, nil
}

func (c *ManagedTokenPoolContract) Bound() bind.IBoundContract {
	return c.BoundContract
}

func (c *ManagedTokenPoolContract) Encoder() ManagedTokenPoolEncoder {
	return c.managedTokenPoolEncoder
}

func (c *ManagedTokenPoolContract) DevInspect() IManagedTokenPoolDevInspect {
	return c.devInspect
}

type MANAGED_TOKEN_POOL struct {
}

type ManagedTokenPoolObject struct {
	Id string `move:"sui::object::UID"`
}

type ManagedTokenPoolStatePointer struct {
	Id                       string `move:"sui::object::UID"`
	ManagedTokenPoolObjectId string `move:"address"`
}

type ManagedTokenPoolState struct {
	Id             string      `move:"sui::object::UID"`
	TokenPoolState bind.Object `move:"TokenPoolState"`
	MintCap        bind.Object `move:"MintCap<T>"`
	OwnableState   bind.Object `move:"OwnableState"`
}

type TokenBucketWrapper struct {
	Tokens      uint64 `move:"u64"`
	LastUpdated uint64 `move:"u64"`
	IsEnabled   bool   `move:"bool"`
	Capacity    uint64 `move:"u64"`
	Rate        uint64 `move:"u64"`
}

type TypeProof struct {
}

type McmsCallback struct {
}

type McmsAcceptOwnershipProof struct {
}

type bcsManagedTokenPoolStatePointer struct {
	Id                       string
	ManagedTokenPoolObjectId [32]byte
}

func convertManagedTokenPoolStatePointerFromBCS(bcs bcsManagedTokenPoolStatePointer) (ManagedTokenPoolStatePointer, error) {

	return ManagedTokenPoolStatePointer{
		Id:                       bcs.Id,
		ManagedTokenPoolObjectId: fmt.Sprintf("0x%x", bcs.ManagedTokenPoolObjectId),
	}, nil
}

func init() {
	bind.RegisterStructDecoder("managed_token_pool::managed_token_pool::MANAGED_TOKEN_POOL", func(data []byte) (interface{}, error) {
		var result MANAGED_TOKEN_POOL
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for MANAGED_TOKEN_POOL
	bind.RegisterStructDecoder("vector<managed_token_pool::managed_token_pool::MANAGED_TOKEN_POOL>", func(data []byte) (interface{}, error) {
		var results []MANAGED_TOKEN_POOL
		_, err := mystenbcs.Unmarshal(data, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
	})
	bind.RegisterStructDecoder("managed_token_pool::managed_token_pool::ManagedTokenPoolObject", func(data []byte) (interface{}, error) {
		var result ManagedTokenPoolObject
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for ManagedTokenPoolObject
	bind.RegisterStructDecoder("vector<managed_token_pool::managed_token_pool::ManagedTokenPoolObject>", func(data []byte) (interface{}, error) {
		var results []ManagedTokenPoolObject
		_, err := mystenbcs.Unmarshal(data, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
	})
	bind.RegisterStructDecoder("managed_token_pool::managed_token_pool::ManagedTokenPoolStatePointer", func(data []byte) (interface{}, error) {
		var temp bcsManagedTokenPoolStatePointer
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertManagedTokenPoolStatePointerFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for ManagedTokenPoolStatePointer
	bind.RegisterStructDecoder("vector<managed_token_pool::managed_token_pool::ManagedTokenPoolStatePointer>", func(data []byte) (interface{}, error) {
		var temps []bcsManagedTokenPoolStatePointer
		_, err := mystenbcs.Unmarshal(data, &temps)
		if err != nil {
			return nil, err
		}

		results := make([]ManagedTokenPoolStatePointer, len(temps))
		for i, temp := range temps {
			result, err := convertManagedTokenPoolStatePointerFromBCS(temp)
			if err != nil {
				return nil, fmt.Errorf("failed to convert element %d: %w", i, err)
			}
			results[i] = result
		}
		return results, nil
	})
	bind.RegisterStructDecoder("managed_token_pool::managed_token_pool::ManagedTokenPoolState", func(data []byte) (interface{}, error) {
		var result ManagedTokenPoolState
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for ManagedTokenPoolState
	bind.RegisterStructDecoder("vector<managed_token_pool::managed_token_pool::ManagedTokenPoolState>", func(data []byte) (interface{}, error) {
		var results []ManagedTokenPoolState
		_, err := mystenbcs.Unmarshal(data, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
	})
	bind.RegisterStructDecoder("managed_token_pool::managed_token_pool::TokenBucketWrapper", func(data []byte) (interface{}, error) {
		var result TokenBucketWrapper
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for TokenBucketWrapper
	bind.RegisterStructDecoder("vector<managed_token_pool::managed_token_pool::TokenBucketWrapper>", func(data []byte) (interface{}, error) {
		var results []TokenBucketWrapper
		_, err := mystenbcs.Unmarshal(data, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
	})
	bind.RegisterStructDecoder("managed_token_pool::managed_token_pool::TypeProof", func(data []byte) (interface{}, error) {
		var result TypeProof
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for TypeProof
	bind.RegisterStructDecoder("vector<managed_token_pool::managed_token_pool::TypeProof>", func(data []byte) (interface{}, error) {
		var results []TypeProof
		_, err := mystenbcs.Unmarshal(data, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
	})
	bind.RegisterStructDecoder("managed_token_pool::managed_token_pool::McmsCallback", func(data []byte) (interface{}, error) {
		var result McmsCallback
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for McmsCallback
	bind.RegisterStructDecoder("vector<managed_token_pool::managed_token_pool::McmsCallback>", func(data []byte) (interface{}, error) {
		var results []McmsCallback
		_, err := mystenbcs.Unmarshal(data, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
	})
	bind.RegisterStructDecoder("managed_token_pool::managed_token_pool::McmsAcceptOwnershipProof", func(data []byte) (interface{}, error) {
		var result McmsAcceptOwnershipProof
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for McmsAcceptOwnershipProof
	bind.RegisterStructDecoder("vector<managed_token_pool::managed_token_pool::McmsAcceptOwnershipProof>", func(data []byte) (interface{}, error) {
		var results []McmsAcceptOwnershipProof
		_, err := mystenbcs.Unmarshal(data, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
	})
}

// TypeAndVersion executes the type_and_version Move function.
func (c *ManagedTokenPoolContract) TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenPoolEncoder.TypeAndVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// InitializeWithManagedToken executes the initialize_with_managed_token Move function.
func (c *ManagedTokenPoolContract) InitializeWithManagedToken(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ownerCap bind.Object, ref bind.Object, managedTokenState bind.Object, managedTokenOwnerCap bind.Object, coinMetadata bind.Object, mintCap bind.Object, tokenPoolAdministrator string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenPoolEncoder.InitializeWithManagedToken(typeArgs, ownerCap, ref, managedTokenState, managedTokenOwnerCap, coinMetadata, mintCap, tokenPoolAdministrator)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// AddRemotePool executes the add_remote_pool Move function.
func (c *ManagedTokenPoolContract) AddRemotePool(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object, remoteChainSelector uint64, remotePoolAddress []byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenPoolEncoder.AddRemotePool(typeArgs, state, ownerCap, remoteChainSelector, remotePoolAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// RemoveRemotePool executes the remove_remote_pool Move function.
func (c *ManagedTokenPoolContract) RemoveRemotePool(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object, remoteChainSelector uint64, remotePoolAddress []byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenPoolEncoder.RemoveRemotePool(typeArgs, state, ownerCap, remoteChainSelector, remotePoolAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// IsSupportedChain executes the is_supported_chain Move function.
func (c *ManagedTokenPoolContract) IsSupportedChain(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenPoolEncoder.IsSupportedChain(typeArgs, state, remoteChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetSupportedChains executes the get_supported_chains Move function.
func (c *ManagedTokenPoolContract) GetSupportedChains(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenPoolEncoder.GetSupportedChains(typeArgs, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ApplyChainUpdates executes the apply_chain_updates Move function.
func (c *ManagedTokenPoolContract) ApplyChainUpdates(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object, remoteChainSelectorsToRemove []uint64, remoteChainSelectorsToAdd []uint64, remotePoolAddressesToAdd [][][]byte, remoteTokenAddressesToAdd [][]byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenPoolEncoder.ApplyChainUpdates(typeArgs, state, ownerCap, remoteChainSelectorsToRemove, remoteChainSelectorsToAdd, remotePoolAddressesToAdd, remoteTokenAddressesToAdd)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetAllowlistEnabled executes the get_allowlist_enabled Move function.
func (c *ManagedTokenPoolContract) GetAllowlistEnabled(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenPoolEncoder.GetAllowlistEnabled(typeArgs, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetAllowlist executes the get_allowlist Move function.
func (c *ManagedTokenPoolContract) GetAllowlist(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenPoolEncoder.GetAllowlist(typeArgs, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// SetAllowlistEnabled executes the set_allowlist_enabled Move function.
func (c *ManagedTokenPoolContract) SetAllowlistEnabled(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object, enabled bool) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenPoolEncoder.SetAllowlistEnabled(typeArgs, state, ownerCap, enabled)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ApplyAllowlistUpdates executes the apply_allowlist_updates Move function.
func (c *ManagedTokenPoolContract) ApplyAllowlistUpdates(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object, removes []string, adds []string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenPoolEncoder.ApplyAllowlistUpdates(typeArgs, state, ownerCap, removes, adds)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetToken executes the get_token Move function.
func (c *ManagedTokenPoolContract) GetToken(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenPoolEncoder.GetToken(typeArgs, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetTokenDecimals executes the get_token_decimals Move function.
func (c *ManagedTokenPoolContract) GetTokenDecimals(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenPoolEncoder.GetTokenDecimals(typeArgs, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetTokenSymbol executes the get_token_symbol Move function.
func (c *ManagedTokenPoolContract) GetTokenSymbol(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenPoolEncoder.GetTokenSymbol(typeArgs, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetRemotePools executes the get_remote_pools Move function.
func (c *ManagedTokenPoolContract) GetRemotePools(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenPoolEncoder.GetRemotePools(typeArgs, state, remoteChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// IsRemotePool executes the is_remote_pool Move function.
func (c *ManagedTokenPoolContract) IsRemotePool(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, remoteChainSelector uint64, remotePoolAddress []byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenPoolEncoder.IsRemotePool(typeArgs, state, remoteChainSelector, remotePoolAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetRemoteToken executes the get_remote_token Move function.
func (c *ManagedTokenPoolContract) GetRemoteToken(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenPoolEncoder.GetRemoteToken(typeArgs, state, remoteChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// LockOrBurn executes the lock_or_burn Move function.
func (c *ManagedTokenPoolContract) LockOrBurn(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, tokenTransferParams bind.Object, c_ bind.Object, remoteChainSelector uint64, clock bind.Object, denyList bind.Object, tokenState bind.Object, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenPoolEncoder.LockOrBurn(typeArgs, ref, tokenTransferParams, c_, remoteChainSelector, clock, denyList, tokenState, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ReleaseOrMint executes the release_or_mint Move function.
func (c *ManagedTokenPoolContract) ReleaseOrMint(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, receiverParams bind.Object, clock bind.Object, denyList bind.Object, tokenState bind.Object, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenPoolEncoder.ReleaseOrMint(typeArgs, ref, receiverParams, clock, denyList, tokenState, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// SetChainRateLimiterConfigs executes the set_chain_rate_limiter_configs Move function.
func (c *ManagedTokenPoolContract) SetChainRateLimiterConfigs(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object, clock bind.Object, remoteChainSelectors []uint64, outboundIsEnableds []bool, outboundCapacities []uint64, outboundRates []uint64, inboundIsEnableds []bool, inboundCapacities []uint64, inboundRates []uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenPoolEncoder.SetChainRateLimiterConfigs(typeArgs, state, ownerCap, clock, remoteChainSelectors, outboundIsEnableds, outboundCapacities, outboundRates, inboundIsEnableds, inboundCapacities, inboundRates)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// SetChainRateLimiterConfig executes the set_chain_rate_limiter_config Move function.
func (c *ManagedTokenPoolContract) SetChainRateLimiterConfig(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object, clock bind.Object, remoteChainSelector uint64, outboundIsEnabled bool, outboundCapacity uint64, outboundRate uint64, inboundIsEnabled bool, inboundCapacity uint64, inboundRate uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenPoolEncoder.SetChainRateLimiterConfig(typeArgs, state, ownerCap, clock, remoteChainSelector, outboundIsEnabled, outboundCapacity, outboundRate, inboundIsEnabled, inboundCapacity, inboundRate)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetCurrentInboundRateLimiterState executes the get_current_inbound_rate_limiter_state Move function.
func (c *ManagedTokenPoolContract) GetCurrentInboundRateLimiterState(ctx context.Context, opts *bind.CallOpts, typeArgs []string, clock bind.Object, state bind.Object, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenPoolEncoder.GetCurrentInboundRateLimiterState(typeArgs, clock, state, remoteChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetCurrentOutboundRateLimiterState executes the get_current_outbound_rate_limiter_state Move function.
func (c *ManagedTokenPoolContract) GetCurrentOutboundRateLimiterState(ctx context.Context, opts *bind.CallOpts, typeArgs []string, clock bind.Object, state bind.Object, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenPoolEncoder.GetCurrentOutboundRateLimiterState(typeArgs, clock, state, remoteChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// Owner executes the owner Move function.
func (c *ManagedTokenPoolContract) Owner(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenPoolEncoder.Owner(typeArgs, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// HasPendingTransfer executes the has_pending_transfer Move function.
func (c *ManagedTokenPoolContract) HasPendingTransfer(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenPoolEncoder.HasPendingTransfer(typeArgs, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// PendingTransferFrom executes the pending_transfer_from Move function.
func (c *ManagedTokenPoolContract) PendingTransferFrom(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenPoolEncoder.PendingTransferFrom(typeArgs, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// PendingTransferTo executes the pending_transfer_to Move function.
func (c *ManagedTokenPoolContract) PendingTransferTo(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenPoolEncoder.PendingTransferTo(typeArgs, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// PendingTransferAccepted executes the pending_transfer_accepted Move function.
func (c *ManagedTokenPoolContract) PendingTransferAccepted(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenPoolEncoder.PendingTransferAccepted(typeArgs, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// TransferOwnership executes the transfer_ownership Move function.
func (c *ManagedTokenPoolContract) TransferOwnership(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object, newOwner string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenPoolEncoder.TransferOwnership(typeArgs, state, ownerCap, newOwner)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// AcceptOwnership executes the accept_ownership Move function.
func (c *ManagedTokenPoolContract) AcceptOwnership(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenPoolEncoder.AcceptOwnership(typeArgs, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// AcceptOwnershipFromObject executes the accept_ownership_from_object Move function.
func (c *ManagedTokenPoolContract) AcceptOwnershipFromObject(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, from string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenPoolEncoder.AcceptOwnershipFromObject(typeArgs, state, from)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsAcceptOwnership executes the mcms_accept_ownership Move function.
func (c *ManagedTokenPoolContract) McmsAcceptOwnership(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenPoolEncoder.McmsAcceptOwnership(typeArgs, state, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ExecuteOwnershipTransfer executes the execute_ownership_transfer Move function.
func (c *ManagedTokenPoolContract) ExecuteOwnershipTransfer(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ownerCap bind.Object, state bind.Object, to string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenPoolEncoder.ExecuteOwnershipTransfer(typeArgs, ownerCap, state, to)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ExecuteOwnershipTransferToMcms executes the execute_ownership_transfer_to_mcms Move function.
func (c *ManagedTokenPoolContract) ExecuteOwnershipTransferToMcms(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ownerCap bind.Object, state bind.Object, registry bind.Object, to string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenPoolEncoder.ExecuteOwnershipTransferToMcms(typeArgs, ownerCap, state, registry, to)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsRegisterUpgradeCap executes the mcms_register_upgrade_cap Move function.
func (c *ManagedTokenPoolContract) McmsRegisterUpgradeCap(ctx context.Context, opts *bind.CallOpts, upgradeCap bind.Object, registry bind.Object, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenPoolEncoder.McmsRegisterUpgradeCap(upgradeCap, registry, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsSetAllowlistEnabled executes the mcms_set_allowlist_enabled Move function.
func (c *ManagedTokenPoolContract) McmsSetAllowlistEnabled(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenPoolEncoder.McmsSetAllowlistEnabled(typeArgs, state, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsApplyAllowlistUpdates executes the mcms_apply_allowlist_updates Move function.
func (c *ManagedTokenPoolContract) McmsApplyAllowlistUpdates(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenPoolEncoder.McmsApplyAllowlistUpdates(typeArgs, state, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsApplyChainUpdates executes the mcms_apply_chain_updates Move function.
func (c *ManagedTokenPoolContract) McmsApplyChainUpdates(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenPoolEncoder.McmsApplyChainUpdates(typeArgs, state, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsAddRemotePool executes the mcms_add_remote_pool Move function.
func (c *ManagedTokenPoolContract) McmsAddRemotePool(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenPoolEncoder.McmsAddRemotePool(typeArgs, state, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsRemoveRemotePool executes the mcms_remove_remote_pool Move function.
func (c *ManagedTokenPoolContract) McmsRemoveRemotePool(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenPoolEncoder.McmsRemoveRemotePool(typeArgs, state, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsSetChainRateLimiterConfigs executes the mcms_set_chain_rate_limiter_configs Move function.
func (c *ManagedTokenPoolContract) McmsSetChainRateLimiterConfigs(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, registry bind.Object, params bind.Object, clock bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenPoolEncoder.McmsSetChainRateLimiterConfigs(typeArgs, state, registry, params, clock)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsSetChainRateLimiterConfig executes the mcms_set_chain_rate_limiter_config Move function.
func (c *ManagedTokenPoolContract) McmsSetChainRateLimiterConfig(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, registry bind.Object, params bind.Object, clock bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenPoolEncoder.McmsSetChainRateLimiterConfig(typeArgs, state, registry, params, clock)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// DestroyTokenPool executes the destroy_token_pool Move function.
func (c *ManagedTokenPoolContract) DestroyTokenPool(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, state bind.Object, ownerCap bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenPoolEncoder.DestroyTokenPool(typeArgs, ref, state, ownerCap)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsDestroyTokenPool executes the mcms_destroy_token_pool Move function.
func (c *ManagedTokenPoolContract) McmsDestroyTokenPool(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenPoolEncoder.McmsDestroyTokenPool(typeArgs, ref, state, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsTransferOwnership executes the mcms_transfer_ownership Move function.
func (c *ManagedTokenPoolContract) McmsTransferOwnership(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenPoolEncoder.McmsTransferOwnership(typeArgs, state, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsExecuteOwnershipTransfer executes the mcms_execute_ownership_transfer Move function.
func (c *ManagedTokenPoolContract) McmsExecuteOwnershipTransfer(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, registry bind.Object, deployerState bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenPoolEncoder.McmsExecuteOwnershipTransfer(typeArgs, state, registry, deployerState, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsAddAllowedModules executes the mcms_add_allowed_modules Move function.
func (c *ManagedTokenPoolContract) McmsAddAllowedModules(ctx context.Context, opts *bind.CallOpts, typeArgs []string, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenPoolEncoder.McmsAddAllowedModules(typeArgs, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsRemoveAllowedModules executes the mcms_remove_allowed_modules Move function.
func (c *ManagedTokenPoolContract) McmsRemoveAllowedModules(ctx context.Context, opts *bind.CallOpts, typeArgs []string, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenPoolEncoder.McmsRemoveAllowedModules(typeArgs, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// TypeAndVersion executes the type_and_version Move function using DevInspect to get return values.
//
// Returns: 0x1::string::String
func (d *ManagedTokenPoolDevInspect) TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (string, error) {
	encoded, err := d.contract.managedTokenPoolEncoder.TypeAndVersion()
	if err != nil {
		return "", fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return "", err
	}
	if len(results) == 0 {
		return "", fmt.Errorf("no return value")
	}
	result, ok := results[0].(string)
	if !ok {
		return "", fmt.Errorf("unexpected return type: expected string, got %T", results[0])
	}
	return result, nil
}

// IsSupportedChain executes the is_supported_chain Move function using DevInspect to get return values.
//
// Returns: bool
func (d *ManagedTokenPoolDevInspect) IsSupportedChain(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, remoteChainSelector uint64) (bool, error) {
	encoded, err := d.contract.managedTokenPoolEncoder.IsSupportedChain(typeArgs, state, remoteChainSelector)
	if err != nil {
		return false, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return false, err
	}
	if len(results) == 0 {
		return false, fmt.Errorf("no return value")
	}
	result, ok := results[0].(bool)
	if !ok {
		return false, fmt.Errorf("unexpected return type: expected bool, got %T", results[0])
	}
	return result, nil
}

// GetSupportedChains executes the get_supported_chains Move function using DevInspect to get return values.
//
// Returns: vector<u64>
func (d *ManagedTokenPoolDevInspect) GetSupportedChains(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) ([]uint64, error) {
	encoded, err := d.contract.managedTokenPoolEncoder.GetSupportedChains(typeArgs, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no return value")
	}
	result, ok := results[0].([]uint64)
	if !ok {
		return nil, fmt.Errorf("unexpected return type: expected []uint64, got %T", results[0])
	}
	return result, nil
}

// GetAllowlistEnabled executes the get_allowlist_enabled Move function using DevInspect to get return values.
//
// Returns: bool
func (d *ManagedTokenPoolDevInspect) GetAllowlistEnabled(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (bool, error) {
	encoded, err := d.contract.managedTokenPoolEncoder.GetAllowlistEnabled(typeArgs, state)
	if err != nil {
		return false, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return false, err
	}
	if len(results) == 0 {
		return false, fmt.Errorf("no return value")
	}
	result, ok := results[0].(bool)
	if !ok {
		return false, fmt.Errorf("unexpected return type: expected bool, got %T", results[0])
	}
	return result, nil
}

// GetAllowlist executes the get_allowlist Move function using DevInspect to get return values.
//
// Returns: vector<address>
func (d *ManagedTokenPoolDevInspect) GetAllowlist(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) ([]string, error) {
	encoded, err := d.contract.managedTokenPoolEncoder.GetAllowlist(typeArgs, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no return value")
	}
	result, ok := results[0].([]string)
	if !ok {
		return nil, fmt.Errorf("unexpected return type: expected []string, got %T", results[0])
	}
	return result, nil
}

// GetToken executes the get_token Move function using DevInspect to get return values.
//
// Returns: address
func (d *ManagedTokenPoolDevInspect) GetToken(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (string, error) {
	encoded, err := d.contract.managedTokenPoolEncoder.GetToken(typeArgs, state)
	if err != nil {
		return "", fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return "", err
	}
	if len(results) == 0 {
		return "", fmt.Errorf("no return value")
	}
	result, ok := results[0].(string)
	if !ok {
		return "", fmt.Errorf("unexpected return type: expected string, got %T", results[0])
	}
	return result, nil
}

// GetTokenDecimals executes the get_token_decimals Move function using DevInspect to get return values.
//
// Returns: u8
func (d *ManagedTokenPoolDevInspect) GetTokenDecimals(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (byte, error) {
	encoded, err := d.contract.managedTokenPoolEncoder.GetTokenDecimals(typeArgs, state)
	if err != nil {
		return 0, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return 0, err
	}
	if len(results) == 0 {
		return 0, fmt.Errorf("no return value")
	}
	result, ok := results[0].(byte)
	if !ok {
		return 0, fmt.Errorf("unexpected return type: expected byte, got %T", results[0])
	}
	return result, nil
}

// GetTokenSymbol executes the get_token_symbol Move function using DevInspect to get return values.
//
// Returns: ascii::String
func (d *ManagedTokenPoolDevInspect) GetTokenSymbol(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (string, error) {
	encoded, err := d.contract.managedTokenPoolEncoder.GetTokenSymbol(typeArgs, state)
	if err != nil {
		return "", fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return "", err
	}
	if len(results) == 0 {
		return "", fmt.Errorf("no return value")
	}
	result, ok := results[0].(string)
	if !ok {
		return "", fmt.Errorf("unexpected return type: expected string, got %T", results[0])
	}
	return result, nil
}

// GetRemotePools executes the get_remote_pools Move function using DevInspect to get return values.
//
// Returns: vector<vector<u8>>
func (d *ManagedTokenPoolDevInspect) GetRemotePools(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, remoteChainSelector uint64) ([][]byte, error) {
	encoded, err := d.contract.managedTokenPoolEncoder.GetRemotePools(typeArgs, state, remoteChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no return value")
	}
	result, ok := results[0].([][]byte)
	if !ok {
		return nil, fmt.Errorf("unexpected return type: expected [][]byte, got %T", results[0])
	}
	return result, nil
}

// IsRemotePool executes the is_remote_pool Move function using DevInspect to get return values.
//
// Returns: bool
func (d *ManagedTokenPoolDevInspect) IsRemotePool(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, remoteChainSelector uint64, remotePoolAddress []byte) (bool, error) {
	encoded, err := d.contract.managedTokenPoolEncoder.IsRemotePool(typeArgs, state, remoteChainSelector, remotePoolAddress)
	if err != nil {
		return false, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return false, err
	}
	if len(results) == 0 {
		return false, fmt.Errorf("no return value")
	}
	result, ok := results[0].(bool)
	if !ok {
		return false, fmt.Errorf("unexpected return type: expected bool, got %T", results[0])
	}
	return result, nil
}

// GetRemoteToken executes the get_remote_token Move function using DevInspect to get return values.
//
// Returns: vector<u8>
func (d *ManagedTokenPoolDevInspect) GetRemoteToken(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, remoteChainSelector uint64) ([]byte, error) {
	encoded, err := d.contract.managedTokenPoolEncoder.GetRemoteToken(typeArgs, state, remoteChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no return value")
	}
	result, ok := results[0].([]byte)
	if !ok {
		return nil, fmt.Errorf("unexpected return type: expected []byte, got %T", results[0])
	}
	return result, nil
}

// GetCurrentInboundRateLimiterState executes the get_current_inbound_rate_limiter_state Move function using DevInspect to get return values.
//
// Returns: TokenBucketWrapper
func (d *ManagedTokenPoolDevInspect) GetCurrentInboundRateLimiterState(ctx context.Context, opts *bind.CallOpts, typeArgs []string, clock bind.Object, state bind.Object, remoteChainSelector uint64) (TokenBucketWrapper, error) {
	encoded, err := d.contract.managedTokenPoolEncoder.GetCurrentInboundRateLimiterState(typeArgs, clock, state, remoteChainSelector)
	if err != nil {
		return TokenBucketWrapper{}, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return TokenBucketWrapper{}, err
	}
	if len(results) == 0 {
		return TokenBucketWrapper{}, fmt.Errorf("no return value")
	}
	result, ok := results[0].(TokenBucketWrapper)
	if !ok {
		return TokenBucketWrapper{}, fmt.Errorf("unexpected return type: expected TokenBucketWrapper, got %T", results[0])
	}
	return result, nil
}

// GetCurrentOutboundRateLimiterState executes the get_current_outbound_rate_limiter_state Move function using DevInspect to get return values.
//
// Returns: TokenBucketWrapper
func (d *ManagedTokenPoolDevInspect) GetCurrentOutboundRateLimiterState(ctx context.Context, opts *bind.CallOpts, typeArgs []string, clock bind.Object, state bind.Object, remoteChainSelector uint64) (TokenBucketWrapper, error) {
	encoded, err := d.contract.managedTokenPoolEncoder.GetCurrentOutboundRateLimiterState(typeArgs, clock, state, remoteChainSelector)
	if err != nil {
		return TokenBucketWrapper{}, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return TokenBucketWrapper{}, err
	}
	if len(results) == 0 {
		return TokenBucketWrapper{}, fmt.Errorf("no return value")
	}
	result, ok := results[0].(TokenBucketWrapper)
	if !ok {
		return TokenBucketWrapper{}, fmt.Errorf("unexpected return type: expected TokenBucketWrapper, got %T", results[0])
	}
	return result, nil
}

// Owner executes the owner Move function using DevInspect to get return values.
//
// Returns: address
func (d *ManagedTokenPoolDevInspect) Owner(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (string, error) {
	encoded, err := d.contract.managedTokenPoolEncoder.Owner(typeArgs, state)
	if err != nil {
		return "", fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return "", err
	}
	if len(results) == 0 {
		return "", fmt.Errorf("no return value")
	}
	result, ok := results[0].(string)
	if !ok {
		return "", fmt.Errorf("unexpected return type: expected string, got %T", results[0])
	}
	return result, nil
}

// HasPendingTransfer executes the has_pending_transfer Move function using DevInspect to get return values.
//
// Returns: bool
func (d *ManagedTokenPoolDevInspect) HasPendingTransfer(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (bool, error) {
	encoded, err := d.contract.managedTokenPoolEncoder.HasPendingTransfer(typeArgs, state)
	if err != nil {
		return false, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return false, err
	}
	if len(results) == 0 {
		return false, fmt.Errorf("no return value")
	}
	result, ok := results[0].(bool)
	if !ok {
		return false, fmt.Errorf("unexpected return type: expected bool, got %T", results[0])
	}
	return result, nil
}

// PendingTransferFrom executes the pending_transfer_from Move function using DevInspect to get return values.
//
// Returns: 0x1::option::Option<address>
func (d *ManagedTokenPoolDevInspect) PendingTransferFrom(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*string, error) {
	encoded, err := d.contract.managedTokenPoolEncoder.PendingTransferFrom(typeArgs, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no return value")
	}
	result, ok := results[0].(*string)
	if !ok {
		return nil, fmt.Errorf("unexpected return type: expected *string, got %T", results[0])
	}
	return result, nil
}

// PendingTransferTo executes the pending_transfer_to Move function using DevInspect to get return values.
//
// Returns: 0x1::option::Option<address>
func (d *ManagedTokenPoolDevInspect) PendingTransferTo(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*string, error) {
	encoded, err := d.contract.managedTokenPoolEncoder.PendingTransferTo(typeArgs, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no return value")
	}
	result, ok := results[0].(*string)
	if !ok {
		return nil, fmt.Errorf("unexpected return type: expected *string, got %T", results[0])
	}
	return result, nil
}

// PendingTransferAccepted executes the pending_transfer_accepted Move function using DevInspect to get return values.
//
// Returns: 0x1::option::Option<bool>
func (d *ManagedTokenPoolDevInspect) PendingTransferAccepted(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*bool, error) {
	encoded, err := d.contract.managedTokenPoolEncoder.PendingTransferAccepted(typeArgs, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no return value")
	}
	result, ok := results[0].(*bool)
	if !ok {
		return nil, fmt.Errorf("unexpected return type: expected *bool, got %T", results[0])
	}
	return result, nil
}

// DestroyTokenPool executes the destroy_token_pool Move function using DevInspect to get return values.
//
// Returns: MintCap<T>
func (d *ManagedTokenPoolDevInspect) DestroyTokenPool(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, state bind.Object, ownerCap bind.Object) (any, error) {
	encoded, err := d.contract.managedTokenPoolEncoder.DestroyTokenPool(typeArgs, ref, state, ownerCap)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no return value")
	}
	return results[0], nil
}

type managedTokenPoolEncoder struct {
	*bind.BoundContract
}

// TypeAndVersion encodes a call to the type_and_version Move function.
func (c managedTokenPoolEncoder) TypeAndVersion() (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("type_and_version", typeArgsList, typeParamsList, []string{}, []any{}, []string{
		"0x1::string::String",
	})
}

// TypeAndVersionWithArgs encodes a call to the type_and_version Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenPoolEncoder) TypeAndVersionWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("type_and_version", typeArgsList, typeParamsList, expectedParams, args, []string{
		"0x1::string::String",
	})
}

// InitializeWithManagedToken encodes a call to the initialize_with_managed_token Move function.
func (c managedTokenPoolEncoder) InitializeWithManagedToken(typeArgs []string, ownerCap bind.Object, ref bind.Object, managedTokenState bind.Object, managedTokenOwnerCap bind.Object, coinMetadata bind.Object, mintCap bind.Object, tokenPoolAdministrator string) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("initialize_with_managed_token", typeArgsList, typeParamsList, []string{
		"&mut OwnerCap",
		"&mut CCIPObjectRef",
		"&TokenState<T>",
		"&ManagedTokenOwnerCap<T>",
		"&CoinMetadata<T>",
		"MintCap<T>",
		"address",
	}, []any{
		ownerCap,
		ref,
		managedTokenState,
		managedTokenOwnerCap,
		coinMetadata,
		mintCap,
		tokenPoolAdministrator,
	}, nil)
}

// InitializeWithManagedTokenWithArgs encodes a call to the initialize_with_managed_token Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenPoolEncoder) InitializeWithManagedTokenWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut OwnerCap",
		"&mut CCIPObjectRef",
		"&TokenState<T>",
		"&ManagedTokenOwnerCap<T>",
		"&CoinMetadata<T>",
		"MintCap<T>",
		"address",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("initialize_with_managed_token", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// AddRemotePool encodes a call to the add_remote_pool Move function.
func (c managedTokenPoolEncoder) AddRemotePool(typeArgs []string, state bind.Object, ownerCap bind.Object, remoteChainSelector uint64, remotePoolAddress []byte) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("add_remote_pool", typeArgsList, typeParamsList, []string{
		"&mut ManagedTokenPoolState<T>",
		"&OwnerCap",
		"u64",
		"vector<u8>",
	}, []any{
		state,
		ownerCap,
		remoteChainSelector,
		remotePoolAddress,
	}, nil)
}

// AddRemotePoolWithArgs encodes a call to the add_remote_pool Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenPoolEncoder) AddRemotePoolWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut ManagedTokenPoolState<T>",
		"&OwnerCap",
		"u64",
		"vector<u8>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("add_remote_pool", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// RemoveRemotePool encodes a call to the remove_remote_pool Move function.
func (c managedTokenPoolEncoder) RemoveRemotePool(typeArgs []string, state bind.Object, ownerCap bind.Object, remoteChainSelector uint64, remotePoolAddress []byte) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("remove_remote_pool", typeArgsList, typeParamsList, []string{
		"&mut ManagedTokenPoolState<T>",
		"&OwnerCap",
		"u64",
		"vector<u8>",
	}, []any{
		state,
		ownerCap,
		remoteChainSelector,
		remotePoolAddress,
	}, nil)
}

// RemoveRemotePoolWithArgs encodes a call to the remove_remote_pool Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenPoolEncoder) RemoveRemotePoolWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut ManagedTokenPoolState<T>",
		"&OwnerCap",
		"u64",
		"vector<u8>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("remove_remote_pool", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// IsSupportedChain encodes a call to the is_supported_chain Move function.
func (c managedTokenPoolEncoder) IsSupportedChain(typeArgs []string, state bind.Object, remoteChainSelector uint64) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("is_supported_chain", typeArgsList, typeParamsList, []string{
		"&ManagedTokenPoolState<T>",
		"u64",
	}, []any{
		state,
		remoteChainSelector,
	}, []string{
		"bool",
	})
}

// IsSupportedChainWithArgs encodes a call to the is_supported_chain Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenPoolEncoder) IsSupportedChainWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&ManagedTokenPoolState<T>",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("is_supported_chain", typeArgsList, typeParamsList, expectedParams, args, []string{
		"bool",
	})
}

// GetSupportedChains encodes a call to the get_supported_chains Move function.
func (c managedTokenPoolEncoder) GetSupportedChains(typeArgs []string, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("get_supported_chains", typeArgsList, typeParamsList, []string{
		"&ManagedTokenPoolState<T>",
	}, []any{
		state,
	}, []string{
		"vector<u64>",
	})
}

// GetSupportedChainsWithArgs encodes a call to the get_supported_chains Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenPoolEncoder) GetSupportedChainsWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&ManagedTokenPoolState<T>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("get_supported_chains", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<u64>",
	})
}

// ApplyChainUpdates encodes a call to the apply_chain_updates Move function.
func (c managedTokenPoolEncoder) ApplyChainUpdates(typeArgs []string, state bind.Object, ownerCap bind.Object, remoteChainSelectorsToRemove []uint64, remoteChainSelectorsToAdd []uint64, remotePoolAddressesToAdd [][][]byte, remoteTokenAddressesToAdd [][]byte) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("apply_chain_updates", typeArgsList, typeParamsList, []string{
		"&mut ManagedTokenPoolState<T>",
		"&OwnerCap",
		"vector<u64>",
		"vector<u64>",
		"vector<vector<vector<u8>>>",
		"vector<vector<u8>>",
	}, []any{
		state,
		ownerCap,
		remoteChainSelectorsToRemove,
		remoteChainSelectorsToAdd,
		remotePoolAddressesToAdd,
		remoteTokenAddressesToAdd,
	}, nil)
}

// ApplyChainUpdatesWithArgs encodes a call to the apply_chain_updates Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenPoolEncoder) ApplyChainUpdatesWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut ManagedTokenPoolState<T>",
		"&OwnerCap",
		"vector<u64>",
		"vector<u64>",
		"vector<vector<vector<u8>>>",
		"vector<vector<u8>>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("apply_chain_updates", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// GetAllowlistEnabled encodes a call to the get_allowlist_enabled Move function.
func (c managedTokenPoolEncoder) GetAllowlistEnabled(typeArgs []string, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("get_allowlist_enabled", typeArgsList, typeParamsList, []string{
		"&ManagedTokenPoolState<T>",
	}, []any{
		state,
	}, []string{
		"bool",
	})
}

// GetAllowlistEnabledWithArgs encodes a call to the get_allowlist_enabled Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenPoolEncoder) GetAllowlistEnabledWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&ManagedTokenPoolState<T>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("get_allowlist_enabled", typeArgsList, typeParamsList, expectedParams, args, []string{
		"bool",
	})
}

// GetAllowlist encodes a call to the get_allowlist Move function.
func (c managedTokenPoolEncoder) GetAllowlist(typeArgs []string, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("get_allowlist", typeArgsList, typeParamsList, []string{
		"&ManagedTokenPoolState<T>",
	}, []any{
		state,
	}, []string{
		"vector<address>",
	})
}

// GetAllowlistWithArgs encodes a call to the get_allowlist Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenPoolEncoder) GetAllowlistWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&ManagedTokenPoolState<T>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("get_allowlist", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<address>",
	})
}

// SetAllowlistEnabled encodes a call to the set_allowlist_enabled Move function.
func (c managedTokenPoolEncoder) SetAllowlistEnabled(typeArgs []string, state bind.Object, ownerCap bind.Object, enabled bool) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("set_allowlist_enabled", typeArgsList, typeParamsList, []string{
		"&mut ManagedTokenPoolState<T>",
		"&OwnerCap",
		"bool",
	}, []any{
		state,
		ownerCap,
		enabled,
	}, nil)
}

// SetAllowlistEnabledWithArgs encodes a call to the set_allowlist_enabled Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenPoolEncoder) SetAllowlistEnabledWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut ManagedTokenPoolState<T>",
		"&OwnerCap",
		"bool",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("set_allowlist_enabled", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// ApplyAllowlistUpdates encodes a call to the apply_allowlist_updates Move function.
func (c managedTokenPoolEncoder) ApplyAllowlistUpdates(typeArgs []string, state bind.Object, ownerCap bind.Object, removes []string, adds []string) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("apply_allowlist_updates", typeArgsList, typeParamsList, []string{
		"&mut ManagedTokenPoolState<T>",
		"&OwnerCap",
		"vector<address>",
		"vector<address>",
	}, []any{
		state,
		ownerCap,
		removes,
		adds,
	}, nil)
}

// ApplyAllowlistUpdatesWithArgs encodes a call to the apply_allowlist_updates Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenPoolEncoder) ApplyAllowlistUpdatesWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut ManagedTokenPoolState<T>",
		"&OwnerCap",
		"vector<address>",
		"vector<address>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("apply_allowlist_updates", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// GetToken encodes a call to the get_token Move function.
func (c managedTokenPoolEncoder) GetToken(typeArgs []string, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("get_token", typeArgsList, typeParamsList, []string{
		"&ManagedTokenPoolState<T>",
	}, []any{
		state,
	}, []string{
		"address",
	})
}

// GetTokenWithArgs encodes a call to the get_token Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenPoolEncoder) GetTokenWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&ManagedTokenPoolState<T>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("get_token", typeArgsList, typeParamsList, expectedParams, args, []string{
		"address",
	})
}

// GetTokenDecimals encodes a call to the get_token_decimals Move function.
func (c managedTokenPoolEncoder) GetTokenDecimals(typeArgs []string, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("get_token_decimals", typeArgsList, typeParamsList, []string{
		"&ManagedTokenPoolState<T>",
	}, []any{
		state,
	}, []string{
		"u8",
	})
}

// GetTokenDecimalsWithArgs encodes a call to the get_token_decimals Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenPoolEncoder) GetTokenDecimalsWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&ManagedTokenPoolState<T>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("get_token_decimals", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u8",
	})
}

// GetTokenSymbol encodes a call to the get_token_symbol Move function.
func (c managedTokenPoolEncoder) GetTokenSymbol(typeArgs []string, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("get_token_symbol", typeArgsList, typeParamsList, []string{
		"&ManagedTokenPoolState<T>",
	}, []any{
		state,
	}, []string{
		"ascii::String",
	})
}

// GetTokenSymbolWithArgs encodes a call to the get_token_symbol Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenPoolEncoder) GetTokenSymbolWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&ManagedTokenPoolState<T>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("get_token_symbol", typeArgsList, typeParamsList, expectedParams, args, []string{
		"ascii::String",
	})
}

// GetRemotePools encodes a call to the get_remote_pools Move function.
func (c managedTokenPoolEncoder) GetRemotePools(typeArgs []string, state bind.Object, remoteChainSelector uint64) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("get_remote_pools", typeArgsList, typeParamsList, []string{
		"&ManagedTokenPoolState<T>",
		"u64",
	}, []any{
		state,
		remoteChainSelector,
	}, []string{
		"vector<vector<u8>>",
	})
}

// GetRemotePoolsWithArgs encodes a call to the get_remote_pools Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenPoolEncoder) GetRemotePoolsWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&ManagedTokenPoolState<T>",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("get_remote_pools", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<vector<u8>>",
	})
}

// IsRemotePool encodes a call to the is_remote_pool Move function.
func (c managedTokenPoolEncoder) IsRemotePool(typeArgs []string, state bind.Object, remoteChainSelector uint64, remotePoolAddress []byte) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("is_remote_pool", typeArgsList, typeParamsList, []string{
		"&ManagedTokenPoolState<T>",
		"u64",
		"vector<u8>",
	}, []any{
		state,
		remoteChainSelector,
		remotePoolAddress,
	}, []string{
		"bool",
	})
}

// IsRemotePoolWithArgs encodes a call to the is_remote_pool Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenPoolEncoder) IsRemotePoolWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&ManagedTokenPoolState<T>",
		"u64",
		"vector<u8>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("is_remote_pool", typeArgsList, typeParamsList, expectedParams, args, []string{
		"bool",
	})
}

// GetRemoteToken encodes a call to the get_remote_token Move function.
func (c managedTokenPoolEncoder) GetRemoteToken(typeArgs []string, state bind.Object, remoteChainSelector uint64) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("get_remote_token", typeArgsList, typeParamsList, []string{
		"&ManagedTokenPoolState<T>",
		"u64",
	}, []any{
		state,
		remoteChainSelector,
	}, []string{
		"vector<u8>",
	})
}

// GetRemoteTokenWithArgs encodes a call to the get_remote_token Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenPoolEncoder) GetRemoteTokenWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&ManagedTokenPoolState<T>",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("get_remote_token", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<u8>",
	})
}

// LockOrBurn encodes a call to the lock_or_burn Move function.
func (c managedTokenPoolEncoder) LockOrBurn(typeArgs []string, ref bind.Object, tokenTransferParams bind.Object, c_ bind.Object, remoteChainSelector uint64, clock bind.Object, denyList bind.Object, tokenState bind.Object, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("lock_or_burn", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"&mut onramp_sh::TokenTransferParams",
		"Coin<T>",
		"u64",
		"&Clock",
		"&DenyList",
		"&mut TokenState<T>",
		"&mut ManagedTokenPoolState<T>",
	}, []any{
		ref,
		tokenTransferParams,
		c_,
		remoteChainSelector,
		clock,
		denyList,
		tokenState,
		state,
	}, nil)
}

// LockOrBurnWithArgs encodes a call to the lock_or_burn Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenPoolEncoder) LockOrBurnWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"&mut onramp_sh::TokenTransferParams",
		"Coin<T>",
		"u64",
		"&Clock",
		"&DenyList",
		"&mut TokenState<T>",
		"&mut ManagedTokenPoolState<T>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("lock_or_burn", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// ReleaseOrMint encodes a call to the release_or_mint Move function.
func (c managedTokenPoolEncoder) ReleaseOrMint(typeArgs []string, ref bind.Object, receiverParams bind.Object, clock bind.Object, denyList bind.Object, tokenState bind.Object, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("release_or_mint", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"&mut offramp_sh::ReceiverParams",
		"&Clock",
		"&DenyList",
		"&mut TokenState<T>",
		"&mut ManagedTokenPoolState<T>",
	}, []any{
		ref,
		receiverParams,
		clock,
		denyList,
		tokenState,
		state,
	}, nil)
}

// ReleaseOrMintWithArgs encodes a call to the release_or_mint Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenPoolEncoder) ReleaseOrMintWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"&mut offramp_sh::ReceiverParams",
		"&Clock",
		"&DenyList",
		"&mut TokenState<T>",
		"&mut ManagedTokenPoolState<T>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("release_or_mint", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// SetChainRateLimiterConfigs encodes a call to the set_chain_rate_limiter_configs Move function.
func (c managedTokenPoolEncoder) SetChainRateLimiterConfigs(typeArgs []string, state bind.Object, ownerCap bind.Object, clock bind.Object, remoteChainSelectors []uint64, outboundIsEnableds []bool, outboundCapacities []uint64, outboundRates []uint64, inboundIsEnableds []bool, inboundCapacities []uint64, inboundRates []uint64) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("set_chain_rate_limiter_configs", typeArgsList, typeParamsList, []string{
		"&mut ManagedTokenPoolState<T>",
		"&OwnerCap",
		"&Clock",
		"vector<u64>",
		"vector<bool>",
		"vector<u64>",
		"vector<u64>",
		"vector<bool>",
		"vector<u64>",
		"vector<u64>",
	}, []any{
		state,
		ownerCap,
		clock,
		remoteChainSelectors,
		outboundIsEnableds,
		outboundCapacities,
		outboundRates,
		inboundIsEnableds,
		inboundCapacities,
		inboundRates,
	}, nil)
}

// SetChainRateLimiterConfigsWithArgs encodes a call to the set_chain_rate_limiter_configs Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenPoolEncoder) SetChainRateLimiterConfigsWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut ManagedTokenPoolState<T>",
		"&OwnerCap",
		"&Clock",
		"vector<u64>",
		"vector<bool>",
		"vector<u64>",
		"vector<u64>",
		"vector<bool>",
		"vector<u64>",
		"vector<u64>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("set_chain_rate_limiter_configs", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// SetChainRateLimiterConfig encodes a call to the set_chain_rate_limiter_config Move function.
func (c managedTokenPoolEncoder) SetChainRateLimiterConfig(typeArgs []string, state bind.Object, ownerCap bind.Object, clock bind.Object, remoteChainSelector uint64, outboundIsEnabled bool, outboundCapacity uint64, outboundRate uint64, inboundIsEnabled bool, inboundCapacity uint64, inboundRate uint64) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("set_chain_rate_limiter_config", typeArgsList, typeParamsList, []string{
		"&mut ManagedTokenPoolState<T>",
		"&OwnerCap",
		"&Clock",
		"u64",
		"bool",
		"u64",
		"u64",
		"bool",
		"u64",
		"u64",
	}, []any{
		state,
		ownerCap,
		clock,
		remoteChainSelector,
		outboundIsEnabled,
		outboundCapacity,
		outboundRate,
		inboundIsEnabled,
		inboundCapacity,
		inboundRate,
	}, nil)
}

// SetChainRateLimiterConfigWithArgs encodes a call to the set_chain_rate_limiter_config Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenPoolEncoder) SetChainRateLimiterConfigWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut ManagedTokenPoolState<T>",
		"&OwnerCap",
		"&Clock",
		"u64",
		"bool",
		"u64",
		"u64",
		"bool",
		"u64",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("set_chain_rate_limiter_config", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// GetCurrentInboundRateLimiterState encodes a call to the get_current_inbound_rate_limiter_state Move function.
func (c managedTokenPoolEncoder) GetCurrentInboundRateLimiterState(typeArgs []string, clock bind.Object, state bind.Object, remoteChainSelector uint64) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("get_current_inbound_rate_limiter_state", typeArgsList, typeParamsList, []string{
		"&Clock",
		"&ManagedTokenPoolState<T>",
		"u64",
	}, []any{
		clock,
		state,
		remoteChainSelector,
	}, []string{
		"managed_token_pool::managed_token_pool::TokenBucketWrapper",
	})
}

// GetCurrentInboundRateLimiterStateWithArgs encodes a call to the get_current_inbound_rate_limiter_state Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenPoolEncoder) GetCurrentInboundRateLimiterStateWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&Clock",
		"&ManagedTokenPoolState<T>",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("get_current_inbound_rate_limiter_state", typeArgsList, typeParamsList, expectedParams, args, []string{
		"managed_token_pool::managed_token_pool::TokenBucketWrapper",
	})
}

// GetCurrentOutboundRateLimiterState encodes a call to the get_current_outbound_rate_limiter_state Move function.
func (c managedTokenPoolEncoder) GetCurrentOutboundRateLimiterState(typeArgs []string, clock bind.Object, state bind.Object, remoteChainSelector uint64) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("get_current_outbound_rate_limiter_state", typeArgsList, typeParamsList, []string{
		"&Clock",
		"&ManagedTokenPoolState<T>",
		"u64",
	}, []any{
		clock,
		state,
		remoteChainSelector,
	}, []string{
		"managed_token_pool::managed_token_pool::TokenBucketWrapper",
	})
}

// GetCurrentOutboundRateLimiterStateWithArgs encodes a call to the get_current_outbound_rate_limiter_state Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenPoolEncoder) GetCurrentOutboundRateLimiterStateWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&Clock",
		"&ManagedTokenPoolState<T>",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("get_current_outbound_rate_limiter_state", typeArgsList, typeParamsList, expectedParams, args, []string{
		"managed_token_pool::managed_token_pool::TokenBucketWrapper",
	})
}

// Owner encodes a call to the owner Move function.
func (c managedTokenPoolEncoder) Owner(typeArgs []string, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("owner", typeArgsList, typeParamsList, []string{
		"&ManagedTokenPoolState<T>",
	}, []any{
		state,
	}, []string{
		"address",
	})
}

// OwnerWithArgs encodes a call to the owner Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenPoolEncoder) OwnerWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&ManagedTokenPoolState<T>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("owner", typeArgsList, typeParamsList, expectedParams, args, []string{
		"address",
	})
}

// HasPendingTransfer encodes a call to the has_pending_transfer Move function.
func (c managedTokenPoolEncoder) HasPendingTransfer(typeArgs []string, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("has_pending_transfer", typeArgsList, typeParamsList, []string{
		"&ManagedTokenPoolState<T>",
	}, []any{
		state,
	}, []string{
		"bool",
	})
}

// HasPendingTransferWithArgs encodes a call to the has_pending_transfer Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenPoolEncoder) HasPendingTransferWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&ManagedTokenPoolState<T>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("has_pending_transfer", typeArgsList, typeParamsList, expectedParams, args, []string{
		"bool",
	})
}

// PendingTransferFrom encodes a call to the pending_transfer_from Move function.
func (c managedTokenPoolEncoder) PendingTransferFrom(typeArgs []string, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("pending_transfer_from", typeArgsList, typeParamsList, []string{
		"&ManagedTokenPoolState<T>",
	}, []any{
		state,
	}, []string{
		"0x1::option::Option<address>",
	})
}

// PendingTransferFromWithArgs encodes a call to the pending_transfer_from Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenPoolEncoder) PendingTransferFromWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&ManagedTokenPoolState<T>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("pending_transfer_from", typeArgsList, typeParamsList, expectedParams, args, []string{
		"0x1::option::Option<address>",
	})
}

// PendingTransferTo encodes a call to the pending_transfer_to Move function.
func (c managedTokenPoolEncoder) PendingTransferTo(typeArgs []string, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("pending_transfer_to", typeArgsList, typeParamsList, []string{
		"&ManagedTokenPoolState<T>",
	}, []any{
		state,
	}, []string{
		"0x1::option::Option<address>",
	})
}

// PendingTransferToWithArgs encodes a call to the pending_transfer_to Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenPoolEncoder) PendingTransferToWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&ManagedTokenPoolState<T>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("pending_transfer_to", typeArgsList, typeParamsList, expectedParams, args, []string{
		"0x1::option::Option<address>",
	})
}

// PendingTransferAccepted encodes a call to the pending_transfer_accepted Move function.
func (c managedTokenPoolEncoder) PendingTransferAccepted(typeArgs []string, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("pending_transfer_accepted", typeArgsList, typeParamsList, []string{
		"&ManagedTokenPoolState<T>",
	}, []any{
		state,
	}, []string{
		"0x1::option::Option<bool>",
	})
}

// PendingTransferAcceptedWithArgs encodes a call to the pending_transfer_accepted Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenPoolEncoder) PendingTransferAcceptedWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&ManagedTokenPoolState<T>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("pending_transfer_accepted", typeArgsList, typeParamsList, expectedParams, args, []string{
		"0x1::option::Option<bool>",
	})
}

// TransferOwnership encodes a call to the transfer_ownership Move function.
func (c managedTokenPoolEncoder) TransferOwnership(typeArgs []string, state bind.Object, ownerCap bind.Object, newOwner string) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("transfer_ownership", typeArgsList, typeParamsList, []string{
		"&mut ManagedTokenPoolState<T>",
		"&OwnerCap",
		"address",
	}, []any{
		state,
		ownerCap,
		newOwner,
	}, nil)
}

// TransferOwnershipWithArgs encodes a call to the transfer_ownership Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenPoolEncoder) TransferOwnershipWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut ManagedTokenPoolState<T>",
		"&OwnerCap",
		"address",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("transfer_ownership", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// AcceptOwnership encodes a call to the accept_ownership Move function.
func (c managedTokenPoolEncoder) AcceptOwnership(typeArgs []string, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("accept_ownership", typeArgsList, typeParamsList, []string{
		"&mut ManagedTokenPoolState<T>",
	}, []any{
		state,
	}, nil)
}

// AcceptOwnershipWithArgs encodes a call to the accept_ownership Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenPoolEncoder) AcceptOwnershipWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut ManagedTokenPoolState<T>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("accept_ownership", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// AcceptOwnershipFromObject encodes a call to the accept_ownership_from_object Move function.
func (c managedTokenPoolEncoder) AcceptOwnershipFromObject(typeArgs []string, state bind.Object, from string) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("accept_ownership_from_object", typeArgsList, typeParamsList, []string{
		"&mut ManagedTokenPoolState<T>",
		"&mut UID",
	}, []any{
		state,
		from,
	}, nil)
}

// AcceptOwnershipFromObjectWithArgs encodes a call to the accept_ownership_from_object Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenPoolEncoder) AcceptOwnershipFromObjectWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut ManagedTokenPoolState<T>",
		"&mut UID",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("accept_ownership_from_object", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// McmsAcceptOwnership encodes a call to the mcms_accept_ownership Move function.
func (c managedTokenPoolEncoder) McmsAcceptOwnership(typeArgs []string, state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("mcms_accept_ownership", typeArgsList, typeParamsList, []string{
		"&mut ManagedTokenPoolState<T>",
		"&mut Registry",
		"ExecutingCallbackParams",
	}, []any{
		state,
		registry,
		params,
	}, nil)
}

// McmsAcceptOwnershipWithArgs encodes a call to the mcms_accept_ownership Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenPoolEncoder) McmsAcceptOwnershipWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut ManagedTokenPoolState<T>",
		"&mut Registry",
		"ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("mcms_accept_ownership", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// ExecuteOwnershipTransfer encodes a call to the execute_ownership_transfer Move function.
func (c managedTokenPoolEncoder) ExecuteOwnershipTransfer(typeArgs []string, ownerCap bind.Object, state bind.Object, to string) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("execute_ownership_transfer", typeArgsList, typeParamsList, []string{
		"OwnerCap",
		"&mut ManagedTokenPoolState<T>",
		"address",
	}, []any{
		ownerCap,
		state,
		to,
	}, nil)
}

// ExecuteOwnershipTransferWithArgs encodes a call to the execute_ownership_transfer Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenPoolEncoder) ExecuteOwnershipTransferWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"OwnerCap",
		"&mut ManagedTokenPoolState<T>",
		"address",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("execute_ownership_transfer", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// ExecuteOwnershipTransferToMcms encodes a call to the execute_ownership_transfer_to_mcms Move function.
func (c managedTokenPoolEncoder) ExecuteOwnershipTransferToMcms(typeArgs []string, ownerCap bind.Object, state bind.Object, registry bind.Object, to string) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("execute_ownership_transfer_to_mcms", typeArgsList, typeParamsList, []string{
		"OwnerCap",
		"&mut ManagedTokenPoolState<T>",
		"&mut Registry",
		"address",
	}, []any{
		ownerCap,
		state,
		registry,
		to,
	}, nil)
}

// ExecuteOwnershipTransferToMcmsWithArgs encodes a call to the execute_ownership_transfer_to_mcms Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenPoolEncoder) ExecuteOwnershipTransferToMcmsWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"OwnerCap",
		"&mut ManagedTokenPoolState<T>",
		"&mut Registry",
		"address",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("execute_ownership_transfer_to_mcms", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// McmsRegisterUpgradeCap encodes a call to the mcms_register_upgrade_cap Move function.
func (c managedTokenPoolEncoder) McmsRegisterUpgradeCap(upgradeCap bind.Object, registry bind.Object, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_register_upgrade_cap", typeArgsList, typeParamsList, []string{
		"UpgradeCap",
		"&mut Registry",
		"&mut DeployerState",
	}, []any{
		upgradeCap,
		registry,
		state,
	}, nil)
}

// McmsRegisterUpgradeCapWithArgs encodes a call to the mcms_register_upgrade_cap Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenPoolEncoder) McmsRegisterUpgradeCapWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"UpgradeCap",
		"&mut Registry",
		"&mut DeployerState",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_register_upgrade_cap", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// McmsSetAllowlistEnabled encodes a call to the mcms_set_allowlist_enabled Move function.
func (c managedTokenPoolEncoder) McmsSetAllowlistEnabled(typeArgs []string, state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("mcms_set_allowlist_enabled", typeArgsList, typeParamsList, []string{
		"&mut ManagedTokenPoolState<T>",
		"&mut Registry",
		"ExecutingCallbackParams",
	}, []any{
		state,
		registry,
		params,
	}, nil)
}

// McmsSetAllowlistEnabledWithArgs encodes a call to the mcms_set_allowlist_enabled Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenPoolEncoder) McmsSetAllowlistEnabledWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut ManagedTokenPoolState<T>",
		"&mut Registry",
		"ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("mcms_set_allowlist_enabled", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// McmsApplyAllowlistUpdates encodes a call to the mcms_apply_allowlist_updates Move function.
func (c managedTokenPoolEncoder) McmsApplyAllowlistUpdates(typeArgs []string, state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("mcms_apply_allowlist_updates", typeArgsList, typeParamsList, []string{
		"&mut ManagedTokenPoolState<T>",
		"&mut Registry",
		"ExecutingCallbackParams",
	}, []any{
		state,
		registry,
		params,
	}, nil)
}

// McmsApplyAllowlistUpdatesWithArgs encodes a call to the mcms_apply_allowlist_updates Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenPoolEncoder) McmsApplyAllowlistUpdatesWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut ManagedTokenPoolState<T>",
		"&mut Registry",
		"ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("mcms_apply_allowlist_updates", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// McmsApplyChainUpdates encodes a call to the mcms_apply_chain_updates Move function.
func (c managedTokenPoolEncoder) McmsApplyChainUpdates(typeArgs []string, state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("mcms_apply_chain_updates", typeArgsList, typeParamsList, []string{
		"&mut ManagedTokenPoolState<T>",
		"&mut Registry",
		"ExecutingCallbackParams",
	}, []any{
		state,
		registry,
		params,
	}, nil)
}

// McmsApplyChainUpdatesWithArgs encodes a call to the mcms_apply_chain_updates Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenPoolEncoder) McmsApplyChainUpdatesWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut ManagedTokenPoolState<T>",
		"&mut Registry",
		"ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("mcms_apply_chain_updates", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// McmsAddRemotePool encodes a call to the mcms_add_remote_pool Move function.
func (c managedTokenPoolEncoder) McmsAddRemotePool(typeArgs []string, state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("mcms_add_remote_pool", typeArgsList, typeParamsList, []string{
		"&mut ManagedTokenPoolState<T>",
		"&mut Registry",
		"ExecutingCallbackParams",
	}, []any{
		state,
		registry,
		params,
	}, nil)
}

// McmsAddRemotePoolWithArgs encodes a call to the mcms_add_remote_pool Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenPoolEncoder) McmsAddRemotePoolWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut ManagedTokenPoolState<T>",
		"&mut Registry",
		"ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("mcms_add_remote_pool", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// McmsRemoveRemotePool encodes a call to the mcms_remove_remote_pool Move function.
func (c managedTokenPoolEncoder) McmsRemoveRemotePool(typeArgs []string, state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("mcms_remove_remote_pool", typeArgsList, typeParamsList, []string{
		"&mut ManagedTokenPoolState<T>",
		"&mut Registry",
		"ExecutingCallbackParams",
	}, []any{
		state,
		registry,
		params,
	}, nil)
}

// McmsRemoveRemotePoolWithArgs encodes a call to the mcms_remove_remote_pool Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenPoolEncoder) McmsRemoveRemotePoolWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut ManagedTokenPoolState<T>",
		"&mut Registry",
		"ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("mcms_remove_remote_pool", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// McmsSetChainRateLimiterConfigs encodes a call to the mcms_set_chain_rate_limiter_configs Move function.
func (c managedTokenPoolEncoder) McmsSetChainRateLimiterConfigs(typeArgs []string, state bind.Object, registry bind.Object, params bind.Object, clock bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("mcms_set_chain_rate_limiter_configs", typeArgsList, typeParamsList, []string{
		"&mut ManagedTokenPoolState<T>",
		"&mut Registry",
		"ExecutingCallbackParams",
		"&Clock",
	}, []any{
		state,
		registry,
		params,
		clock,
	}, nil)
}

// McmsSetChainRateLimiterConfigsWithArgs encodes a call to the mcms_set_chain_rate_limiter_configs Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenPoolEncoder) McmsSetChainRateLimiterConfigsWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut ManagedTokenPoolState<T>",
		"&mut Registry",
		"ExecutingCallbackParams",
		"&Clock",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("mcms_set_chain_rate_limiter_configs", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// McmsSetChainRateLimiterConfig encodes a call to the mcms_set_chain_rate_limiter_config Move function.
func (c managedTokenPoolEncoder) McmsSetChainRateLimiterConfig(typeArgs []string, state bind.Object, registry bind.Object, params bind.Object, clock bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("mcms_set_chain_rate_limiter_config", typeArgsList, typeParamsList, []string{
		"&mut ManagedTokenPoolState<T>",
		"&mut Registry",
		"ExecutingCallbackParams",
		"&Clock",
	}, []any{
		state,
		registry,
		params,
		clock,
	}, nil)
}

// McmsSetChainRateLimiterConfigWithArgs encodes a call to the mcms_set_chain_rate_limiter_config Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenPoolEncoder) McmsSetChainRateLimiterConfigWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut ManagedTokenPoolState<T>",
		"&mut Registry",
		"ExecutingCallbackParams",
		"&Clock",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("mcms_set_chain_rate_limiter_config", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// DestroyTokenPool encodes a call to the destroy_token_pool Move function.
func (c managedTokenPoolEncoder) DestroyTokenPool(typeArgs []string, ref bind.Object, state bind.Object, ownerCap bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("destroy_token_pool", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"ManagedTokenPoolState<T>",
		"OwnerCap",
	}, []any{
		ref,
		state,
		ownerCap,
	}, []string{
		"MintCap<T>",
	})
}

// DestroyTokenPoolWithArgs encodes a call to the destroy_token_pool Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenPoolEncoder) DestroyTokenPoolWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"ManagedTokenPoolState<T>",
		"OwnerCap",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("destroy_token_pool", typeArgsList, typeParamsList, expectedParams, args, []string{
		"MintCap<T>",
	})
}

// McmsDestroyTokenPool encodes a call to the mcms_destroy_token_pool Move function.
func (c managedTokenPoolEncoder) McmsDestroyTokenPool(typeArgs []string, ref bind.Object, state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("mcms_destroy_token_pool", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"ManagedTokenPoolState<T>",
		"&mut Registry",
		"ExecutingCallbackParams",
	}, []any{
		ref,
		state,
		registry,
		params,
	}, nil)
}

// McmsDestroyTokenPoolWithArgs encodes a call to the mcms_destroy_token_pool Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenPoolEncoder) McmsDestroyTokenPoolWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"ManagedTokenPoolState<T>",
		"&mut Registry",
		"ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("mcms_destroy_token_pool", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// McmsTransferOwnership encodes a call to the mcms_transfer_ownership Move function.
func (c managedTokenPoolEncoder) McmsTransferOwnership(typeArgs []string, state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("mcms_transfer_ownership", typeArgsList, typeParamsList, []string{
		"&mut ManagedTokenPoolState<T>",
		"&mut Registry",
		"ExecutingCallbackParams",
	}, []any{
		state,
		registry,
		params,
	}, nil)
}

// McmsTransferOwnershipWithArgs encodes a call to the mcms_transfer_ownership Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenPoolEncoder) McmsTransferOwnershipWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut ManagedTokenPoolState<T>",
		"&mut Registry",
		"ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("mcms_transfer_ownership", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// McmsExecuteOwnershipTransfer encodes a call to the mcms_execute_ownership_transfer Move function.
func (c managedTokenPoolEncoder) McmsExecuteOwnershipTransfer(typeArgs []string, state bind.Object, registry bind.Object, deployerState bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("mcms_execute_ownership_transfer", typeArgsList, typeParamsList, []string{
		"&mut ManagedTokenPoolState<T>",
		"&mut Registry",
		"&mut DeployerState",
		"ExecutingCallbackParams",
	}, []any{
		state,
		registry,
		deployerState,
		params,
	}, nil)
}

// McmsExecuteOwnershipTransferWithArgs encodes a call to the mcms_execute_ownership_transfer Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenPoolEncoder) McmsExecuteOwnershipTransferWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut ManagedTokenPoolState<T>",
		"&mut Registry",
		"&mut DeployerState",
		"ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("mcms_execute_ownership_transfer", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// McmsAddAllowedModules encodes a call to the mcms_add_allowed_modules Move function.
func (c managedTokenPoolEncoder) McmsAddAllowedModules(typeArgs []string, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("mcms_add_allowed_modules", typeArgsList, typeParamsList, []string{
		"&mut Registry",
		"ExecutingCallbackParams",
	}, []any{
		registry,
		params,
	}, nil)
}

// McmsAddAllowedModulesWithArgs encodes a call to the mcms_add_allowed_modules Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenPoolEncoder) McmsAddAllowedModulesWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut Registry",
		"ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("mcms_add_allowed_modules", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// McmsRemoveAllowedModules encodes a call to the mcms_remove_allowed_modules Move function.
func (c managedTokenPoolEncoder) McmsRemoveAllowedModules(typeArgs []string, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("mcms_remove_allowed_modules", typeArgsList, typeParamsList, []string{
		"&mut Registry",
		"ExecutingCallbackParams",
	}, []any{
		registry,
		params,
	}, nil)
}

// McmsRemoveAllowedModulesWithArgs encodes a call to the mcms_remove_allowed_modules Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenPoolEncoder) McmsRemoveAllowedModulesWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut Registry",
		"ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("mcms_remove_allowed_modules", typeArgsList, typeParamsList, expectedParams, args, nil)
}
