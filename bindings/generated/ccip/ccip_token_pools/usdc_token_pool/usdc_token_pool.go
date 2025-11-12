// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package module_usdc_token_pool

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

const FunctionInfo = `[{"package":"usdc_token_pool","module":"usdc_token_pool","name":"accept_ownership","parameters":[{"name":"state","type":"USDCTokenPoolState<T>"}]},{"package":"usdc_token_pool","module":"usdc_token_pool","name":"accept_ownership_from_object","parameters":[{"name":"state","type":"USDCTokenPoolState<T>"},{"name":"from","type":"sui::object::UID"}]},{"package":"usdc_token_pool","module":"usdc_token_pool","name":"add_remote_pool","parameters":[{"name":"state","type":"USDCTokenPoolState<T>"},{"name":"owner_cap","type":"OwnerCap"},{"name":"remote_chain_selector","type":"u64"},{"name":"remote_pool_address","type":"vector<u8>"}]},{"package":"usdc_token_pool","module":"usdc_token_pool","name":"apply_allowlist_updates","parameters":[{"name":"state","type":"USDCTokenPoolState<T>"},{"name":"owner_cap","type":"OwnerCap"},{"name":"removes","type":"vector<address>"},{"name":"adds","type":"vector<address>"}]},{"package":"usdc_token_pool","module":"usdc_token_pool","name":"apply_chain_updates","parameters":[{"name":"state","type":"USDCTokenPoolState<T>"},{"name":"owner_cap","type":"OwnerCap"},{"name":"remote_chain_selectors_to_remove","type":"vector<u64>"},{"name":"remote_chain_selectors_to_add","type":"vector<u64>"},{"name":"remote_pool_addresses_to_add","type":"vector<vector<vector<u8>>>"},{"name":"remote_token_addresses_to_add","type":"vector<vector<u8>>"}]},{"package":"usdc_token_pool","module":"usdc_token_pool","name":"destroy_token_pool","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"state","type":"USDCTokenPoolState<T>"},{"name":"owner_cap","type":"OwnerCap"}]},{"package":"usdc_token_pool","module":"usdc_token_pool","name":"execute_ownership_transfer","parameters":[{"name":"owner_cap","type":"OwnerCap"},{"name":"state","type":"USDCTokenPoolState<T>"},{"name":"to","type":"address"}]},{"package":"usdc_token_pool","module":"usdc_token_pool","name":"execute_ownership_transfer_to_mcms","parameters":[{"name":"owner_cap","type":"OwnerCap"},{"name":"state","type":"USDCTokenPoolState<T>"},{"name":"registry","type":"Registry"},{"name":"to","type":"address"}]},{"package":"usdc_token_pool","module":"usdc_token_pool","name":"get_allowlist","parameters":[{"name":"state","type":"USDCTokenPoolState<T>"}]},{"package":"usdc_token_pool","module":"usdc_token_pool","name":"get_allowlist_enabled","parameters":[{"name":"state","type":"USDCTokenPoolState<T>"}]},{"package":"usdc_token_pool","module":"usdc_token_pool","name":"get_current_inbound_rate_limiter_state","parameters":[{"name":"clock","type":"Clock"},{"name":"state","type":"USDCTokenPoolState<T>"},{"name":"remote_chain_selector","type":"u64"}]},{"package":"usdc_token_pool","module":"usdc_token_pool","name":"get_current_outbound_rate_limiter_state","parameters":[{"name":"clock","type":"Clock"},{"name":"state","type":"USDCTokenPoolState<T>"},{"name":"remote_chain_selector","type":"u64"}]},{"package":"usdc_token_pool","module":"usdc_token_pool","name":"get_domain","parameters":[{"name":"pool","type":"USDCTokenPoolState<T>"},{"name":"chain_selector","type":"u64"}]},{"package":"usdc_token_pool","module":"usdc_token_pool","name":"get_package_auth_caller","parameters":null},{"package":"usdc_token_pool","module":"usdc_token_pool","name":"get_remote_pools","parameters":[{"name":"state","type":"USDCTokenPoolState<T>"},{"name":"remote_chain_selector","type":"u64"}]},{"package":"usdc_token_pool","module":"usdc_token_pool","name":"get_remote_token","parameters":[{"name":"state","type":"USDCTokenPoolState<T>"},{"name":"remote_chain_selector","type":"u64"}]},{"package":"usdc_token_pool","module":"usdc_token_pool","name":"get_supported_chains","parameters":[{"name":"state","type":"USDCTokenPoolState<T>"}]},{"package":"usdc_token_pool","module":"usdc_token_pool","name":"get_token","parameters":[{"name":"state","type":"USDCTokenPoolState<T>"}]},{"package":"usdc_token_pool","module":"usdc_token_pool","name":"get_token_decimals","parameters":[{"name":"state","type":"USDCTokenPoolState<T>"}]},{"package":"usdc_token_pool","module":"usdc_token_pool","name":"has_pending_transfer","parameters":[{"name":"state","type":"USDCTokenPoolState<T>"}]},{"package":"usdc_token_pool","module":"usdc_token_pool","name":"initialize","parameters":[{"name":"owner_cap","type":"OwnerCap"},{"name":"coin_metadata","type":"CoinMetadata<T>"},{"name":"local_domain_identifier","type":"u32"}]},{"package":"usdc_token_pool","module":"usdc_token_pool","name":"is_remote_pool","parameters":[{"name":"state","type":"USDCTokenPoolState<T>"},{"name":"remote_chain_selector","type":"u64"},{"name":"remote_pool_address","type":"vector<u8>"}]},{"package":"usdc_token_pool","module":"usdc_token_pool","name":"is_supported_chain","parameters":[{"name":"state","type":"USDCTokenPoolState<T>"},{"name":"remote_chain_selector","type":"u64"}]},{"package":"usdc_token_pool","module":"usdc_token_pool","name":"lock_or_burn","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"token_transfer_params","type":"onramp_sh::TokenTransferParams"},{"name":"c","type":"Coin<T>"},{"name":"remote_chain_selector","type":"u64"},{"name":"clock","type":"Clock"},{"name":"deny_list","type":"DenyList"},{"name":"pool","type":"USDCTokenPoolState<T>"},{"name":"state","type":"MinterState"},{"name":"message_transmitter_state","type":"MessageTransmitterState"},{"name":"treasury","type":"Treasury<T>"}]},{"package":"usdc_token_pool","module":"usdc_token_pool","name":"owner","parameters":[{"name":"state","type":"USDCTokenPoolState<T>"}]},{"package":"usdc_token_pool","module":"usdc_token_pool","name":"pending_transfer_accepted","parameters":[{"name":"state","type":"USDCTokenPoolState<T>"}]},{"package":"usdc_token_pool","module":"usdc_token_pool","name":"pending_transfer_from","parameters":[{"name":"state","type":"USDCTokenPoolState<T>"}]},{"package":"usdc_token_pool","module":"usdc_token_pool","name":"pending_transfer_to","parameters":[{"name":"state","type":"USDCTokenPoolState<T>"}]},{"package":"usdc_token_pool","module":"usdc_token_pool","name":"release_or_mint","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"receiver_params","type":"offramp_sh::ReceiverParams"},{"name":"clock","type":"Clock"},{"name":"deny_list","type":"DenyList"},{"name":"pool","type":"USDCTokenPoolState<T>"},{"name":"state","type":"MinterState"},{"name":"message_transmitter_state","type":"MessageTransmitterState"},{"name":"treasury","type":"Treasury<T>"}]},{"package":"usdc_token_pool","module":"usdc_token_pool","name":"remove_remote_pool","parameters":[{"name":"state","type":"USDCTokenPoolState<T>"},{"name":"owner_cap","type":"OwnerCap"},{"name":"remote_chain_selector","type":"u64"},{"name":"remote_pool_address","type":"vector<u8>"}]},{"package":"usdc_token_pool","module":"usdc_token_pool","name":"set_allowlist_enabled","parameters":[{"name":"state","type":"USDCTokenPoolState<T>"},{"name":"owner_cap","type":"OwnerCap"},{"name":"enabled","type":"bool"}]},{"package":"usdc_token_pool","module":"usdc_token_pool","name":"set_chain_rate_limiter_config","parameters":[{"name":"state","type":"USDCTokenPoolState<T>"},{"name":"owner_cap","type":"OwnerCap"},{"name":"clock","type":"Clock"},{"name":"remote_chain_selector","type":"u64"},{"name":"outbound_is_enabled","type":"bool"},{"name":"outbound_capacity","type":"u64"},{"name":"outbound_rate","type":"u64"},{"name":"inbound_is_enabled","type":"bool"},{"name":"inbound_capacity","type":"u64"},{"name":"inbound_rate","type":"u64"}]},{"package":"usdc_token_pool","module":"usdc_token_pool","name":"set_chain_rate_limiter_configs","parameters":[{"name":"state","type":"USDCTokenPoolState<T>"},{"name":"owner_cap","type":"OwnerCap"},{"name":"clock","type":"Clock"},{"name":"remote_chain_selectors","type":"vector<u64>"},{"name":"outbound_is_enableds","type":"vector<bool>"},{"name":"outbound_capacities","type":"vector<u64>"},{"name":"outbound_rates","type":"vector<u64>"},{"name":"inbound_is_enableds","type":"vector<bool>"},{"name":"inbound_capacities","type":"vector<u64>"},{"name":"inbound_rates","type":"vector<u64>"}]},{"package":"usdc_token_pool","module":"usdc_token_pool","name":"set_domains","parameters":[{"name":"pool","type":"USDCTokenPoolState<T>"},{"name":"owner_cap","type":"OwnerCap"},{"name":"remote_chain_selectors","type":"vector<u64>"},{"name":"remote_domain_identifiers","type":"vector<u32>"},{"name":"allowed_remote_callers","type":"vector<vector<u8>>"},{"name":"enableds","type":"vector<bool>"}]},{"package":"usdc_token_pool","module":"usdc_token_pool","name":"transfer_ownership","parameters":[{"name":"state","type":"USDCTokenPoolState<T>"},{"name":"owner_cap","type":"OwnerCap"},{"name":"new_owner","type":"address"}]},{"package":"usdc_token_pool","module":"usdc_token_pool","name":"type_and_version","parameters":null}]`

type IUsdcTokenPool interface {
	TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	Initialize(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ownerCap bind.Object, coinMetadata bind.Object, localDomainIdentifier uint32) (*models.SuiTransactionBlockResponse, error)
	GetToken(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	GetTokenDecimals(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	GetRemotePools(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	IsRemotePool(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, remoteChainSelector uint64, remotePoolAddress []byte) (*models.SuiTransactionBlockResponse, error)
	GetRemoteToken(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	AddRemotePool(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object, remoteChainSelector uint64, remotePoolAddress []byte) (*models.SuiTransactionBlockResponse, error)
	RemoveRemotePool(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object, remoteChainSelector uint64, remotePoolAddress []byte) (*models.SuiTransactionBlockResponse, error)
	IsSupportedChain(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	GetSupportedChains(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	ApplyChainUpdates(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object, remoteChainSelectorsToRemove []uint64, remoteChainSelectorsToAdd []uint64, remotePoolAddressesToAdd [][][]byte, remoteTokenAddressesToAdd [][]byte) (*models.SuiTransactionBlockResponse, error)
	GetAllowlistEnabled(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	GetAllowlist(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	SetAllowlistEnabled(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object, enabled bool) (*models.SuiTransactionBlockResponse, error)
	ApplyAllowlistUpdates(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object, removes []string, adds []string) (*models.SuiTransactionBlockResponse, error)
	GetPackageAuthCaller(ctx context.Context, opts *bind.CallOpts, typeArgs []string) (*models.SuiTransactionBlockResponse, error)
	LockOrBurn(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, tokenTransferParams bind.Object, c_ bind.Object, remoteChainSelector uint64, clock bind.Object, denyList bind.Object, pool bind.Object, state bind.Object, messageTransmitterState bind.Object, treasury bind.Object) (*models.SuiTransactionBlockResponse, error)
	ReleaseOrMint(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, receiverParams bind.Object, clock bind.Object, denyList bind.Object, pool bind.Object, state bind.Object, messageTransmitterState bind.Object, treasury bind.Object) (*models.SuiTransactionBlockResponse, error)
	GetDomain(ctx context.Context, opts *bind.CallOpts, typeArgs []string, pool bind.Object, chainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	SetDomains(ctx context.Context, opts *bind.CallOpts, typeArgs []string, pool bind.Object, ownerCap bind.Object, remoteChainSelectors []uint64, remoteDomainIdentifiers []uint32, allowedRemoteCallers [][]byte, enableds []bool) (*models.SuiTransactionBlockResponse, error)
	McmsSetDomains(ctx context.Context, opts *bind.CallOpts, typeArgs []string, pool bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
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
	DevInspect() IUsdcTokenPoolDevInspect
	Encoder() UsdcTokenPoolEncoder
	Bound() bind.IBoundContract
}

type IUsdcTokenPoolDevInspect interface {
	TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (string, error)
	GetToken(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (string, error)
	GetTokenDecimals(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (byte, error)
	GetRemotePools(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, remoteChainSelector uint64) ([][]byte, error)
	IsRemotePool(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, remoteChainSelector uint64, remotePoolAddress []byte) (bool, error)
	GetRemoteToken(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, remoteChainSelector uint64) ([]byte, error)
	IsSupportedChain(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, remoteChainSelector uint64) (bool, error)
	GetSupportedChains(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) ([]uint64, error)
	GetAllowlistEnabled(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (bool, error)
	GetAllowlist(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) ([]string, error)
	GetPackageAuthCaller(ctx context.Context, opts *bind.CallOpts, typeArgs []string) (string, error)
	GetDomain(ctx context.Context, opts *bind.CallOpts, typeArgs []string, pool bind.Object, chainSelector uint64) (Domain, error)
	GetCurrentInboundRateLimiterState(ctx context.Context, opts *bind.CallOpts, typeArgs []string, clock bind.Object, state bind.Object, remoteChainSelector uint64) (TokenBucketWrapper, error)
	GetCurrentOutboundRateLimiterState(ctx context.Context, opts *bind.CallOpts, typeArgs []string, clock bind.Object, state bind.Object, remoteChainSelector uint64) (TokenBucketWrapper, error)
	Owner(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (string, error)
	HasPendingTransfer(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (bool, error)
	PendingTransferFrom(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*string, error)
	PendingTransferTo(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*string, error)
	PendingTransferAccepted(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*bool, error)
}

type UsdcTokenPoolEncoder interface {
	TypeAndVersion() (*bind.EncodedCall, error)
	TypeAndVersionWithArgs(args ...any) (*bind.EncodedCall, error)
	Initialize(typeArgs []string, ownerCap bind.Object, coinMetadata bind.Object, localDomainIdentifier uint32) (*bind.EncodedCall, error)
	InitializeWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	GetToken(typeArgs []string, state bind.Object) (*bind.EncodedCall, error)
	GetTokenWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	GetTokenDecimals(typeArgs []string, state bind.Object) (*bind.EncodedCall, error)
	GetTokenDecimalsWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	GetRemotePools(typeArgs []string, state bind.Object, remoteChainSelector uint64) (*bind.EncodedCall, error)
	GetRemotePoolsWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	IsRemotePool(typeArgs []string, state bind.Object, remoteChainSelector uint64, remotePoolAddress []byte) (*bind.EncodedCall, error)
	IsRemotePoolWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	GetRemoteToken(typeArgs []string, state bind.Object, remoteChainSelector uint64) (*bind.EncodedCall, error)
	GetRemoteTokenWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
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
	GetPackageAuthCaller(typeArgs []string) (*bind.EncodedCall, error)
	GetPackageAuthCallerWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	LockOrBurn(typeArgs []string, ref bind.Object, tokenTransferParams bind.Object, c_ bind.Object, remoteChainSelector uint64, clock bind.Object, denyList bind.Object, pool bind.Object, state bind.Object, messageTransmitterState bind.Object, treasury bind.Object) (*bind.EncodedCall, error)
	LockOrBurnWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	ReleaseOrMint(typeArgs []string, ref bind.Object, receiverParams bind.Object, clock bind.Object, denyList bind.Object, pool bind.Object, state bind.Object, messageTransmitterState bind.Object, treasury bind.Object) (*bind.EncodedCall, error)
	ReleaseOrMintWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	GetDomain(typeArgs []string, pool bind.Object, chainSelector uint64) (*bind.EncodedCall, error)
	GetDomainWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	SetDomains(typeArgs []string, pool bind.Object, ownerCap bind.Object, remoteChainSelectors []uint64, remoteDomainIdentifiers []uint32, allowedRemoteCallers [][]byte, enableds []bool) (*bind.EncodedCall, error)
	SetDomainsWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	McmsSetDomains(typeArgs []string, pool bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsSetDomainsWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
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

type UsdcTokenPoolContract struct {
	*bind.BoundContract
	usdcTokenPoolEncoder
	devInspect *UsdcTokenPoolDevInspect
}

type UsdcTokenPoolDevInspect struct {
	contract *UsdcTokenPoolContract
}

var _ IUsdcTokenPool = (*UsdcTokenPoolContract)(nil)
var _ IUsdcTokenPoolDevInspect = (*UsdcTokenPoolDevInspect)(nil)

func NewUsdcTokenPool(packageID string, client sui.ISuiAPI) (IUsdcTokenPool, error) {
	contract, err := bind.NewBoundContract(packageID, "usdc_token_pool", "usdc_token_pool", client)
	if err != nil {
		return nil, err
	}

	c := &UsdcTokenPoolContract{
		BoundContract:        contract,
		usdcTokenPoolEncoder: usdcTokenPoolEncoder{BoundContract: contract},
	}
	c.devInspect = &UsdcTokenPoolDevInspect{contract: c}
	return c, nil
}

func (c *UsdcTokenPoolContract) Bound() bind.IBoundContract {
	return c.BoundContract
}

func (c *UsdcTokenPoolContract) Encoder() UsdcTokenPoolEncoder {
	return c.usdcTokenPoolEncoder
}

func (c *UsdcTokenPoolContract) DevInspect() IUsdcTokenPoolDevInspect {
	return c.devInspect
}

type USDC_TOKEN_POOL struct {
}

type USDCTokenPoolObject struct {
	Id string `move:"sui::object::UID"`
}

type USDCTokenPoolStatePointer struct {
	Id                    string `move:"sui::object::UID"`
	UsdcTokenPoolObjectId string `move:"address"`
}

type Domain struct {
	AllowedCaller    []byte `move:"vector<u8>"`
	DomainIdentifier uint32 `move:"u32"`
	Enabled          bool   `move:"bool"`
}

type DomainsSet struct {
	AllowedCaller       []byte `move:"vector<u8>"`
	DomainIdentifier    uint32 `move:"u32"`
	RemoteChainSelector uint64 `move:"u64"`
	Enabled             bool   `move:"bool"`
}

type USDCTokenPoolState struct {
	Id                    string      `move:"sui::object::UID"`
	TokenPoolState        bind.Object `move:"TokenPoolState"`
	ChainToDomain         bind.Object `move:"Table<u64, Domain>"`
	LocalDomainIdentifier uint32      `move:"u32"`
	OwnableState          bind.Object `move:"OwnableState"`
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

type bcsUSDCTokenPoolStatePointer struct {
	Id                    string
	UsdcTokenPoolObjectId [32]byte
}

func convertUSDCTokenPoolStatePointerFromBCS(bcs bcsUSDCTokenPoolStatePointer) (USDCTokenPoolStatePointer, error) {

	return USDCTokenPoolStatePointer{
		Id:                    bcs.Id,
		UsdcTokenPoolObjectId: fmt.Sprintf("0x%x", bcs.UsdcTokenPoolObjectId),
	}, nil
}

func init() {
	bind.RegisterStructDecoder("usdc_token_pool::usdc_token_pool::USDC_TOKEN_POOL", func(data []byte) (interface{}, error) {
		var result USDC_TOKEN_POOL
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for USDC_TOKEN_POOL
	bind.RegisterStructDecoder("vector<usdc_token_pool::usdc_token_pool::USDC_TOKEN_POOL>", func(data []byte) (interface{}, error) {
		var results []USDC_TOKEN_POOL
		_, err := mystenbcs.Unmarshal(data, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
	})
	bind.RegisterStructDecoder("usdc_token_pool::usdc_token_pool::USDCTokenPoolObject", func(data []byte) (interface{}, error) {
		var result USDCTokenPoolObject
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for USDCTokenPoolObject
	bind.RegisterStructDecoder("vector<usdc_token_pool::usdc_token_pool::USDCTokenPoolObject>", func(data []byte) (interface{}, error) {
		var results []USDCTokenPoolObject
		_, err := mystenbcs.Unmarshal(data, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
	})
	bind.RegisterStructDecoder("usdc_token_pool::usdc_token_pool::USDCTokenPoolStatePointer", func(data []byte) (interface{}, error) {
		var temp bcsUSDCTokenPoolStatePointer
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertUSDCTokenPoolStatePointerFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for USDCTokenPoolStatePointer
	bind.RegisterStructDecoder("vector<usdc_token_pool::usdc_token_pool::USDCTokenPoolStatePointer>", func(data []byte) (interface{}, error) {
		var temps []bcsUSDCTokenPoolStatePointer
		_, err := mystenbcs.Unmarshal(data, &temps)
		if err != nil {
			return nil, err
		}

		results := make([]USDCTokenPoolStatePointer, len(temps))
		for i, temp := range temps {
			result, err := convertUSDCTokenPoolStatePointerFromBCS(temp)
			if err != nil {
				return nil, fmt.Errorf("failed to convert element %d: %w", i, err)
			}
			results[i] = result
		}
		return results, nil
	})
	bind.RegisterStructDecoder("usdc_token_pool::usdc_token_pool::Domain", func(data []byte) (interface{}, error) {
		var result Domain
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for Domain
	bind.RegisterStructDecoder("vector<usdc_token_pool::usdc_token_pool::Domain>", func(data []byte) (interface{}, error) {
		var results []Domain
		_, err := mystenbcs.Unmarshal(data, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
	})
	bind.RegisterStructDecoder("usdc_token_pool::usdc_token_pool::DomainsSet", func(data []byte) (interface{}, error) {
		var result DomainsSet
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for DomainsSet
	bind.RegisterStructDecoder("vector<usdc_token_pool::usdc_token_pool::DomainsSet>", func(data []byte) (interface{}, error) {
		var results []DomainsSet
		_, err := mystenbcs.Unmarshal(data, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
	})
	bind.RegisterStructDecoder("usdc_token_pool::usdc_token_pool::USDCTokenPoolState", func(data []byte) (interface{}, error) {
		var result USDCTokenPoolState
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for USDCTokenPoolState
	bind.RegisterStructDecoder("vector<usdc_token_pool::usdc_token_pool::USDCTokenPoolState>", func(data []byte) (interface{}, error) {
		var results []USDCTokenPoolState
		_, err := mystenbcs.Unmarshal(data, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
	})
	bind.RegisterStructDecoder("usdc_token_pool::usdc_token_pool::TokenBucketWrapper", func(data []byte) (interface{}, error) {
		var result TokenBucketWrapper
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for TokenBucketWrapper
	bind.RegisterStructDecoder("vector<usdc_token_pool::usdc_token_pool::TokenBucketWrapper>", func(data []byte) (interface{}, error) {
		var results []TokenBucketWrapper
		_, err := mystenbcs.Unmarshal(data, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
	})
	bind.RegisterStructDecoder("usdc_token_pool::usdc_token_pool::TypeProof", func(data []byte) (interface{}, error) {
		var result TypeProof
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for TypeProof
	bind.RegisterStructDecoder("vector<usdc_token_pool::usdc_token_pool::TypeProof>", func(data []byte) (interface{}, error) {
		var results []TypeProof
		_, err := mystenbcs.Unmarshal(data, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
	})
	bind.RegisterStructDecoder("usdc_token_pool::usdc_token_pool::McmsCallback", func(data []byte) (interface{}, error) {
		var result McmsCallback
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for McmsCallback
	bind.RegisterStructDecoder("vector<usdc_token_pool::usdc_token_pool::McmsCallback>", func(data []byte) (interface{}, error) {
		var results []McmsCallback
		_, err := mystenbcs.Unmarshal(data, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
	})
	bind.RegisterStructDecoder("usdc_token_pool::usdc_token_pool::McmsAcceptOwnershipProof", func(data []byte) (interface{}, error) {
		var result McmsAcceptOwnershipProof
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for McmsAcceptOwnershipProof
	bind.RegisterStructDecoder("vector<usdc_token_pool::usdc_token_pool::McmsAcceptOwnershipProof>", func(data []byte) (interface{}, error) {
		var results []McmsAcceptOwnershipProof
		_, err := mystenbcs.Unmarshal(data, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
	})
}

// TypeAndVersion executes the type_and_version Move function.
func (c *UsdcTokenPoolContract) TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.TypeAndVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// Initialize executes the initialize Move function.
func (c *UsdcTokenPoolContract) Initialize(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ownerCap bind.Object, coinMetadata bind.Object, localDomainIdentifier uint32) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.Initialize(typeArgs, ownerCap, coinMetadata, localDomainIdentifier)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetToken executes the get_token Move function.
func (c *UsdcTokenPoolContract) GetToken(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.GetToken(typeArgs, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetTokenDecimals executes the get_token_decimals Move function.
func (c *UsdcTokenPoolContract) GetTokenDecimals(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.GetTokenDecimals(typeArgs, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetRemotePools executes the get_remote_pools Move function.
func (c *UsdcTokenPoolContract) GetRemotePools(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.GetRemotePools(typeArgs, state, remoteChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// IsRemotePool executes the is_remote_pool Move function.
func (c *UsdcTokenPoolContract) IsRemotePool(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, remoteChainSelector uint64, remotePoolAddress []byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.IsRemotePool(typeArgs, state, remoteChainSelector, remotePoolAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetRemoteToken executes the get_remote_token Move function.
func (c *UsdcTokenPoolContract) GetRemoteToken(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.GetRemoteToken(typeArgs, state, remoteChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// AddRemotePool executes the add_remote_pool Move function.
func (c *UsdcTokenPoolContract) AddRemotePool(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object, remoteChainSelector uint64, remotePoolAddress []byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.AddRemotePool(typeArgs, state, ownerCap, remoteChainSelector, remotePoolAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// RemoveRemotePool executes the remove_remote_pool Move function.
func (c *UsdcTokenPoolContract) RemoveRemotePool(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object, remoteChainSelector uint64, remotePoolAddress []byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.RemoveRemotePool(typeArgs, state, ownerCap, remoteChainSelector, remotePoolAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// IsSupportedChain executes the is_supported_chain Move function.
func (c *UsdcTokenPoolContract) IsSupportedChain(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.IsSupportedChain(typeArgs, state, remoteChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetSupportedChains executes the get_supported_chains Move function.
func (c *UsdcTokenPoolContract) GetSupportedChains(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.GetSupportedChains(typeArgs, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ApplyChainUpdates executes the apply_chain_updates Move function.
func (c *UsdcTokenPoolContract) ApplyChainUpdates(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object, remoteChainSelectorsToRemove []uint64, remoteChainSelectorsToAdd []uint64, remotePoolAddressesToAdd [][][]byte, remoteTokenAddressesToAdd [][]byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.ApplyChainUpdates(typeArgs, state, ownerCap, remoteChainSelectorsToRemove, remoteChainSelectorsToAdd, remotePoolAddressesToAdd, remoteTokenAddressesToAdd)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetAllowlistEnabled executes the get_allowlist_enabled Move function.
func (c *UsdcTokenPoolContract) GetAllowlistEnabled(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.GetAllowlistEnabled(typeArgs, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetAllowlist executes the get_allowlist Move function.
func (c *UsdcTokenPoolContract) GetAllowlist(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.GetAllowlist(typeArgs, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// SetAllowlistEnabled executes the set_allowlist_enabled Move function.
func (c *UsdcTokenPoolContract) SetAllowlistEnabled(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object, enabled bool) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.SetAllowlistEnabled(typeArgs, state, ownerCap, enabled)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ApplyAllowlistUpdates executes the apply_allowlist_updates Move function.
func (c *UsdcTokenPoolContract) ApplyAllowlistUpdates(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object, removes []string, adds []string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.ApplyAllowlistUpdates(typeArgs, state, ownerCap, removes, adds)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetPackageAuthCaller executes the get_package_auth_caller Move function.
func (c *UsdcTokenPoolContract) GetPackageAuthCaller(ctx context.Context, opts *bind.CallOpts, typeArgs []string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.GetPackageAuthCaller(typeArgs)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// LockOrBurn executes the lock_or_burn Move function.
func (c *UsdcTokenPoolContract) LockOrBurn(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, tokenTransferParams bind.Object, c_ bind.Object, remoteChainSelector uint64, clock bind.Object, denyList bind.Object, pool bind.Object, state bind.Object, messageTransmitterState bind.Object, treasury bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.LockOrBurn(typeArgs, ref, tokenTransferParams, c_, remoteChainSelector, clock, denyList, pool, state, messageTransmitterState, treasury)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ReleaseOrMint executes the release_or_mint Move function.
func (c *UsdcTokenPoolContract) ReleaseOrMint(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, receiverParams bind.Object, clock bind.Object, denyList bind.Object, pool bind.Object, state bind.Object, messageTransmitterState bind.Object, treasury bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.ReleaseOrMint(typeArgs, ref, receiverParams, clock, denyList, pool, state, messageTransmitterState, treasury)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetDomain executes the get_domain Move function.
func (c *UsdcTokenPoolContract) GetDomain(ctx context.Context, opts *bind.CallOpts, typeArgs []string, pool bind.Object, chainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.GetDomain(typeArgs, pool, chainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// SetDomains executes the set_domains Move function.
func (c *UsdcTokenPoolContract) SetDomains(ctx context.Context, opts *bind.CallOpts, typeArgs []string, pool bind.Object, ownerCap bind.Object, remoteChainSelectors []uint64, remoteDomainIdentifiers []uint32, allowedRemoteCallers [][]byte, enableds []bool) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.SetDomains(typeArgs, pool, ownerCap, remoteChainSelectors, remoteDomainIdentifiers, allowedRemoteCallers, enableds)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsSetDomains executes the mcms_set_domains Move function.
func (c *UsdcTokenPoolContract) McmsSetDomains(ctx context.Context, opts *bind.CallOpts, typeArgs []string, pool bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.McmsSetDomains(typeArgs, pool, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// SetChainRateLimiterConfigs executes the set_chain_rate_limiter_configs Move function.
func (c *UsdcTokenPoolContract) SetChainRateLimiterConfigs(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object, clock bind.Object, remoteChainSelectors []uint64, outboundIsEnableds []bool, outboundCapacities []uint64, outboundRates []uint64, inboundIsEnableds []bool, inboundCapacities []uint64, inboundRates []uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.SetChainRateLimiterConfigs(typeArgs, state, ownerCap, clock, remoteChainSelectors, outboundIsEnableds, outboundCapacities, outboundRates, inboundIsEnableds, inboundCapacities, inboundRates)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// SetChainRateLimiterConfig executes the set_chain_rate_limiter_config Move function.
func (c *UsdcTokenPoolContract) SetChainRateLimiterConfig(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object, clock bind.Object, remoteChainSelector uint64, outboundIsEnabled bool, outboundCapacity uint64, outboundRate uint64, inboundIsEnabled bool, inboundCapacity uint64, inboundRate uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.SetChainRateLimiterConfig(typeArgs, state, ownerCap, clock, remoteChainSelector, outboundIsEnabled, outboundCapacity, outboundRate, inboundIsEnabled, inboundCapacity, inboundRate)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetCurrentInboundRateLimiterState executes the get_current_inbound_rate_limiter_state Move function.
func (c *UsdcTokenPoolContract) GetCurrentInboundRateLimiterState(ctx context.Context, opts *bind.CallOpts, typeArgs []string, clock bind.Object, state bind.Object, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.GetCurrentInboundRateLimiterState(typeArgs, clock, state, remoteChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetCurrentOutboundRateLimiterState executes the get_current_outbound_rate_limiter_state Move function.
func (c *UsdcTokenPoolContract) GetCurrentOutboundRateLimiterState(ctx context.Context, opts *bind.CallOpts, typeArgs []string, clock bind.Object, state bind.Object, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.GetCurrentOutboundRateLimiterState(typeArgs, clock, state, remoteChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// Owner executes the owner Move function.
func (c *UsdcTokenPoolContract) Owner(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.Owner(typeArgs, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// HasPendingTransfer executes the has_pending_transfer Move function.
func (c *UsdcTokenPoolContract) HasPendingTransfer(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.HasPendingTransfer(typeArgs, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// PendingTransferFrom executes the pending_transfer_from Move function.
func (c *UsdcTokenPoolContract) PendingTransferFrom(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.PendingTransferFrom(typeArgs, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// PendingTransferTo executes the pending_transfer_to Move function.
func (c *UsdcTokenPoolContract) PendingTransferTo(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.PendingTransferTo(typeArgs, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// PendingTransferAccepted executes the pending_transfer_accepted Move function.
func (c *UsdcTokenPoolContract) PendingTransferAccepted(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.PendingTransferAccepted(typeArgs, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// TransferOwnership executes the transfer_ownership Move function.
func (c *UsdcTokenPoolContract) TransferOwnership(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object, newOwner string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.TransferOwnership(typeArgs, state, ownerCap, newOwner)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// AcceptOwnership executes the accept_ownership Move function.
func (c *UsdcTokenPoolContract) AcceptOwnership(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.AcceptOwnership(typeArgs, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// AcceptOwnershipFromObject executes the accept_ownership_from_object Move function.
func (c *UsdcTokenPoolContract) AcceptOwnershipFromObject(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, from string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.AcceptOwnershipFromObject(typeArgs, state, from)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsAcceptOwnership executes the mcms_accept_ownership Move function.
func (c *UsdcTokenPoolContract) McmsAcceptOwnership(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.McmsAcceptOwnership(typeArgs, state, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ExecuteOwnershipTransfer executes the execute_ownership_transfer Move function.
func (c *UsdcTokenPoolContract) ExecuteOwnershipTransfer(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ownerCap bind.Object, state bind.Object, to string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.ExecuteOwnershipTransfer(typeArgs, ownerCap, state, to)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ExecuteOwnershipTransferToMcms executes the execute_ownership_transfer_to_mcms Move function.
func (c *UsdcTokenPoolContract) ExecuteOwnershipTransferToMcms(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ownerCap bind.Object, state bind.Object, registry bind.Object, to string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.ExecuteOwnershipTransferToMcms(typeArgs, ownerCap, state, registry, to)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsRegisterUpgradeCap executes the mcms_register_upgrade_cap Move function.
func (c *UsdcTokenPoolContract) McmsRegisterUpgradeCap(ctx context.Context, opts *bind.CallOpts, upgradeCap bind.Object, registry bind.Object, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.McmsRegisterUpgradeCap(upgradeCap, registry, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsSetAllowlistEnabled executes the mcms_set_allowlist_enabled Move function.
func (c *UsdcTokenPoolContract) McmsSetAllowlistEnabled(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.McmsSetAllowlistEnabled(typeArgs, state, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsApplyAllowlistUpdates executes the mcms_apply_allowlist_updates Move function.
func (c *UsdcTokenPoolContract) McmsApplyAllowlistUpdates(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.McmsApplyAllowlistUpdates(typeArgs, state, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsApplyChainUpdates executes the mcms_apply_chain_updates Move function.
func (c *UsdcTokenPoolContract) McmsApplyChainUpdates(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.McmsApplyChainUpdates(typeArgs, state, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsAddRemotePool executes the mcms_add_remote_pool Move function.
func (c *UsdcTokenPoolContract) McmsAddRemotePool(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.McmsAddRemotePool(typeArgs, state, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsRemoveRemotePool executes the mcms_remove_remote_pool Move function.
func (c *UsdcTokenPoolContract) McmsRemoveRemotePool(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.McmsRemoveRemotePool(typeArgs, state, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsSetChainRateLimiterConfigs executes the mcms_set_chain_rate_limiter_configs Move function.
func (c *UsdcTokenPoolContract) McmsSetChainRateLimiterConfigs(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, registry bind.Object, params bind.Object, clock bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.McmsSetChainRateLimiterConfigs(typeArgs, state, registry, params, clock)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsSetChainRateLimiterConfig executes the mcms_set_chain_rate_limiter_config Move function.
func (c *UsdcTokenPoolContract) McmsSetChainRateLimiterConfig(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, registry bind.Object, params bind.Object, clock bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.McmsSetChainRateLimiterConfig(typeArgs, state, registry, params, clock)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// DestroyTokenPool executes the destroy_token_pool Move function.
func (c *UsdcTokenPoolContract) DestroyTokenPool(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, state bind.Object, ownerCap bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.DestroyTokenPool(typeArgs, ref, state, ownerCap)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsDestroyTokenPool executes the mcms_destroy_token_pool Move function.
func (c *UsdcTokenPoolContract) McmsDestroyTokenPool(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.McmsDestroyTokenPool(typeArgs, ref, state, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsTransferOwnership executes the mcms_transfer_ownership Move function.
func (c *UsdcTokenPoolContract) McmsTransferOwnership(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.McmsTransferOwnership(typeArgs, state, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsExecuteOwnershipTransfer executes the mcms_execute_ownership_transfer Move function.
func (c *UsdcTokenPoolContract) McmsExecuteOwnershipTransfer(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, registry bind.Object, deployerState bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.McmsExecuteOwnershipTransfer(typeArgs, state, registry, deployerState, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsAddAllowedModules executes the mcms_add_allowed_modules Move function.
func (c *UsdcTokenPoolContract) McmsAddAllowedModules(ctx context.Context, opts *bind.CallOpts, typeArgs []string, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.McmsAddAllowedModules(typeArgs, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsRemoveAllowedModules executes the mcms_remove_allowed_modules Move function.
func (c *UsdcTokenPoolContract) McmsRemoveAllowedModules(ctx context.Context, opts *bind.CallOpts, typeArgs []string, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.McmsRemoveAllowedModules(typeArgs, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// TypeAndVersion executes the type_and_version Move function using DevInspect to get return values.
//
// Returns: 0x1::string::String
func (d *UsdcTokenPoolDevInspect) TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (string, error) {
	encoded, err := d.contract.usdcTokenPoolEncoder.TypeAndVersion()
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

// GetToken executes the get_token Move function using DevInspect to get return values.
//
// Returns: address
func (d *UsdcTokenPoolDevInspect) GetToken(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (string, error) {
	encoded, err := d.contract.usdcTokenPoolEncoder.GetToken(typeArgs, state)
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
func (d *UsdcTokenPoolDevInspect) GetTokenDecimals(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (byte, error) {
	encoded, err := d.contract.usdcTokenPoolEncoder.GetTokenDecimals(typeArgs, state)
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

// GetRemotePools executes the get_remote_pools Move function using DevInspect to get return values.
//
// Returns: vector<vector<u8>>
func (d *UsdcTokenPoolDevInspect) GetRemotePools(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, remoteChainSelector uint64) ([][]byte, error) {
	encoded, err := d.contract.usdcTokenPoolEncoder.GetRemotePools(typeArgs, state, remoteChainSelector)
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
func (d *UsdcTokenPoolDevInspect) IsRemotePool(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, remoteChainSelector uint64, remotePoolAddress []byte) (bool, error) {
	encoded, err := d.contract.usdcTokenPoolEncoder.IsRemotePool(typeArgs, state, remoteChainSelector, remotePoolAddress)
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
func (d *UsdcTokenPoolDevInspect) GetRemoteToken(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, remoteChainSelector uint64) ([]byte, error) {
	encoded, err := d.contract.usdcTokenPoolEncoder.GetRemoteToken(typeArgs, state, remoteChainSelector)
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

// IsSupportedChain executes the is_supported_chain Move function using DevInspect to get return values.
//
// Returns: bool
func (d *UsdcTokenPoolDevInspect) IsSupportedChain(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, remoteChainSelector uint64) (bool, error) {
	encoded, err := d.contract.usdcTokenPoolEncoder.IsSupportedChain(typeArgs, state, remoteChainSelector)
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
func (d *UsdcTokenPoolDevInspect) GetSupportedChains(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) ([]uint64, error) {
	encoded, err := d.contract.usdcTokenPoolEncoder.GetSupportedChains(typeArgs, state)
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
func (d *UsdcTokenPoolDevInspect) GetAllowlistEnabled(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (bool, error) {
	encoded, err := d.contract.usdcTokenPoolEncoder.GetAllowlistEnabled(typeArgs, state)
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
func (d *UsdcTokenPoolDevInspect) GetAllowlist(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) ([]string, error) {
	encoded, err := d.contract.usdcTokenPoolEncoder.GetAllowlist(typeArgs, state)
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

// GetPackageAuthCaller executes the get_package_auth_caller Move function using DevInspect to get return values.
//
// Returns: address
func (d *UsdcTokenPoolDevInspect) GetPackageAuthCaller(ctx context.Context, opts *bind.CallOpts, typeArgs []string) (string, error) {
	encoded, err := d.contract.usdcTokenPoolEncoder.GetPackageAuthCaller(typeArgs)
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

// GetDomain executes the get_domain Move function using DevInspect to get return values.
//
// Returns: Domain
func (d *UsdcTokenPoolDevInspect) GetDomain(ctx context.Context, opts *bind.CallOpts, typeArgs []string, pool bind.Object, chainSelector uint64) (Domain, error) {
	encoded, err := d.contract.usdcTokenPoolEncoder.GetDomain(typeArgs, pool, chainSelector)
	if err != nil {
		return Domain{}, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return Domain{}, err
	}
	if len(results) == 0 {
		return Domain{}, fmt.Errorf("no return value")
	}
	result, ok := results[0].(Domain)
	if !ok {
		return Domain{}, fmt.Errorf("unexpected return type: expected Domain, got %T", results[0])
	}
	return result, nil
}

// GetCurrentInboundRateLimiterState executes the get_current_inbound_rate_limiter_state Move function using DevInspect to get return values.
//
// Returns: TokenBucketWrapper
func (d *UsdcTokenPoolDevInspect) GetCurrentInboundRateLimiterState(ctx context.Context, opts *bind.CallOpts, typeArgs []string, clock bind.Object, state bind.Object, remoteChainSelector uint64) (TokenBucketWrapper, error) {
	encoded, err := d.contract.usdcTokenPoolEncoder.GetCurrentInboundRateLimiterState(typeArgs, clock, state, remoteChainSelector)
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
func (d *UsdcTokenPoolDevInspect) GetCurrentOutboundRateLimiterState(ctx context.Context, opts *bind.CallOpts, typeArgs []string, clock bind.Object, state bind.Object, remoteChainSelector uint64) (TokenBucketWrapper, error) {
	encoded, err := d.contract.usdcTokenPoolEncoder.GetCurrentOutboundRateLimiterState(typeArgs, clock, state, remoteChainSelector)
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
func (d *UsdcTokenPoolDevInspect) Owner(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (string, error) {
	encoded, err := d.contract.usdcTokenPoolEncoder.Owner(typeArgs, state)
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
func (d *UsdcTokenPoolDevInspect) HasPendingTransfer(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (bool, error) {
	encoded, err := d.contract.usdcTokenPoolEncoder.HasPendingTransfer(typeArgs, state)
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
func (d *UsdcTokenPoolDevInspect) PendingTransferFrom(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*string, error) {
	encoded, err := d.contract.usdcTokenPoolEncoder.PendingTransferFrom(typeArgs, state)
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
func (d *UsdcTokenPoolDevInspect) PendingTransferTo(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*string, error) {
	encoded, err := d.contract.usdcTokenPoolEncoder.PendingTransferTo(typeArgs, state)
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
func (d *UsdcTokenPoolDevInspect) PendingTransferAccepted(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*bool, error) {
	encoded, err := d.contract.usdcTokenPoolEncoder.PendingTransferAccepted(typeArgs, state)
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

type usdcTokenPoolEncoder struct {
	*bind.BoundContract
}

// TypeAndVersion encodes a call to the type_and_version Move function.
func (c usdcTokenPoolEncoder) TypeAndVersion() (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("type_and_version", typeArgsList, typeParamsList, []string{}, []any{}, []string{
		"0x1::string::String",
	})
}

// TypeAndVersionWithArgs encodes a call to the type_and_version Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c usdcTokenPoolEncoder) TypeAndVersionWithArgs(args ...any) (*bind.EncodedCall, error) {
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

// Initialize encodes a call to the initialize Move function.
func (c usdcTokenPoolEncoder) Initialize(typeArgs []string, ownerCap bind.Object, coinMetadata bind.Object, localDomainIdentifier uint32) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("initialize", typeArgsList, typeParamsList, []string{
		"&mut OwnerCap",
		"&CoinMetadata<T>",
		"u32",
	}, []any{
		ownerCap,
		coinMetadata,
		localDomainIdentifier,
	}, nil)
}

// InitializeWithArgs encodes a call to the initialize Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c usdcTokenPoolEncoder) InitializeWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut OwnerCap",
		"&CoinMetadata<T>",
		"u32",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("initialize", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// GetToken encodes a call to the get_token Move function.
func (c usdcTokenPoolEncoder) GetToken(typeArgs []string, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("get_token", typeArgsList, typeParamsList, []string{
		"&USDCTokenPoolState<T>",
	}, []any{
		state,
	}, []string{
		"address",
	})
}

// GetTokenWithArgs encodes a call to the get_token Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c usdcTokenPoolEncoder) GetTokenWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&USDCTokenPoolState<T>",
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
func (c usdcTokenPoolEncoder) GetTokenDecimals(typeArgs []string, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("get_token_decimals", typeArgsList, typeParamsList, []string{
		"&USDCTokenPoolState<T>",
	}, []any{
		state,
	}, []string{
		"u8",
	})
}

// GetTokenDecimalsWithArgs encodes a call to the get_token_decimals Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c usdcTokenPoolEncoder) GetTokenDecimalsWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&USDCTokenPoolState<T>",
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

// GetRemotePools encodes a call to the get_remote_pools Move function.
func (c usdcTokenPoolEncoder) GetRemotePools(typeArgs []string, state bind.Object, remoteChainSelector uint64) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("get_remote_pools", typeArgsList, typeParamsList, []string{
		"&USDCTokenPoolState<T>",
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
func (c usdcTokenPoolEncoder) GetRemotePoolsWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&USDCTokenPoolState<T>",
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
func (c usdcTokenPoolEncoder) IsRemotePool(typeArgs []string, state bind.Object, remoteChainSelector uint64, remotePoolAddress []byte) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("is_remote_pool", typeArgsList, typeParamsList, []string{
		"&USDCTokenPoolState<T>",
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
func (c usdcTokenPoolEncoder) IsRemotePoolWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&USDCTokenPoolState<T>",
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
func (c usdcTokenPoolEncoder) GetRemoteToken(typeArgs []string, state bind.Object, remoteChainSelector uint64) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("get_remote_token", typeArgsList, typeParamsList, []string{
		"&USDCTokenPoolState<T>",
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
func (c usdcTokenPoolEncoder) GetRemoteTokenWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&USDCTokenPoolState<T>",
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

// AddRemotePool encodes a call to the add_remote_pool Move function.
func (c usdcTokenPoolEncoder) AddRemotePool(typeArgs []string, state bind.Object, ownerCap bind.Object, remoteChainSelector uint64, remotePoolAddress []byte) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("add_remote_pool", typeArgsList, typeParamsList, []string{
		"&mut USDCTokenPoolState<T>",
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
func (c usdcTokenPoolEncoder) AddRemotePoolWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut USDCTokenPoolState<T>",
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
func (c usdcTokenPoolEncoder) RemoveRemotePool(typeArgs []string, state bind.Object, ownerCap bind.Object, remoteChainSelector uint64, remotePoolAddress []byte) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("remove_remote_pool", typeArgsList, typeParamsList, []string{
		"&mut USDCTokenPoolState<T>",
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
func (c usdcTokenPoolEncoder) RemoveRemotePoolWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut USDCTokenPoolState<T>",
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
func (c usdcTokenPoolEncoder) IsSupportedChain(typeArgs []string, state bind.Object, remoteChainSelector uint64) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("is_supported_chain", typeArgsList, typeParamsList, []string{
		"&USDCTokenPoolState<T>",
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
func (c usdcTokenPoolEncoder) IsSupportedChainWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&USDCTokenPoolState<T>",
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
func (c usdcTokenPoolEncoder) GetSupportedChains(typeArgs []string, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("get_supported_chains", typeArgsList, typeParamsList, []string{
		"&USDCTokenPoolState<T>",
	}, []any{
		state,
	}, []string{
		"vector<u64>",
	})
}

// GetSupportedChainsWithArgs encodes a call to the get_supported_chains Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c usdcTokenPoolEncoder) GetSupportedChainsWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&USDCTokenPoolState<T>",
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
func (c usdcTokenPoolEncoder) ApplyChainUpdates(typeArgs []string, state bind.Object, ownerCap bind.Object, remoteChainSelectorsToRemove []uint64, remoteChainSelectorsToAdd []uint64, remotePoolAddressesToAdd [][][]byte, remoteTokenAddressesToAdd [][]byte) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("apply_chain_updates", typeArgsList, typeParamsList, []string{
		"&mut USDCTokenPoolState<T>",
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
func (c usdcTokenPoolEncoder) ApplyChainUpdatesWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut USDCTokenPoolState<T>",
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
func (c usdcTokenPoolEncoder) GetAllowlistEnabled(typeArgs []string, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("get_allowlist_enabled", typeArgsList, typeParamsList, []string{
		"&USDCTokenPoolState<T>",
	}, []any{
		state,
	}, []string{
		"bool",
	})
}

// GetAllowlistEnabledWithArgs encodes a call to the get_allowlist_enabled Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c usdcTokenPoolEncoder) GetAllowlistEnabledWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&USDCTokenPoolState<T>",
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
func (c usdcTokenPoolEncoder) GetAllowlist(typeArgs []string, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("get_allowlist", typeArgsList, typeParamsList, []string{
		"&USDCTokenPoolState<T>",
	}, []any{
		state,
	}, []string{
		"vector<address>",
	})
}

// GetAllowlistWithArgs encodes a call to the get_allowlist Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c usdcTokenPoolEncoder) GetAllowlistWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&USDCTokenPoolState<T>",
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
func (c usdcTokenPoolEncoder) SetAllowlistEnabled(typeArgs []string, state bind.Object, ownerCap bind.Object, enabled bool) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("set_allowlist_enabled", typeArgsList, typeParamsList, []string{
		"&mut USDCTokenPoolState<T>",
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
func (c usdcTokenPoolEncoder) SetAllowlistEnabledWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut USDCTokenPoolState<T>",
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
func (c usdcTokenPoolEncoder) ApplyAllowlistUpdates(typeArgs []string, state bind.Object, ownerCap bind.Object, removes []string, adds []string) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("apply_allowlist_updates", typeArgsList, typeParamsList, []string{
		"&mut USDCTokenPoolState<T>",
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
func (c usdcTokenPoolEncoder) ApplyAllowlistUpdatesWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut USDCTokenPoolState<T>",
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

// GetPackageAuthCaller encodes a call to the get_package_auth_caller Move function.
func (c usdcTokenPoolEncoder) GetPackageAuthCaller(typeArgs []string) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"TypeProof",
	}
	return c.EncodeCallArgsWithGenerics("get_package_auth_caller", typeArgsList, typeParamsList, []string{}, []any{}, []string{
		"address",
	})
}

// GetPackageAuthCallerWithArgs encodes a call to the get_package_auth_caller Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c usdcTokenPoolEncoder) GetPackageAuthCallerWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"TypeProof",
	}
	return c.EncodeCallArgsWithGenerics("get_package_auth_caller", typeArgsList, typeParamsList, expectedParams, args, []string{
		"address",
	})
}

// LockOrBurn encodes a call to the lock_or_burn Move function.
func (c usdcTokenPoolEncoder) LockOrBurn(typeArgs []string, ref bind.Object, tokenTransferParams bind.Object, c_ bind.Object, remoteChainSelector uint64, clock bind.Object, denyList bind.Object, pool bind.Object, state bind.Object, messageTransmitterState bind.Object, treasury bind.Object) (*bind.EncodedCall, error) {
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
		"&mut USDCTokenPoolState<T>",
		"&MinterState",
		"&mut MessageTransmitterState",
		"&mut Treasury<T>",
	}, []any{
		ref,
		tokenTransferParams,
		c_,
		remoteChainSelector,
		clock,
		denyList,
		pool,
		state,
		messageTransmitterState,
		treasury,
	}, nil)
}

// LockOrBurnWithArgs encodes a call to the lock_or_burn Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c usdcTokenPoolEncoder) LockOrBurnWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"&mut onramp_sh::TokenTransferParams",
		"Coin<T>",
		"u64",
		"&Clock",
		"&DenyList",
		"&mut USDCTokenPoolState<T>",
		"&MinterState",
		"&mut MessageTransmitterState",
		"&mut Treasury<T>",
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
func (c usdcTokenPoolEncoder) ReleaseOrMint(typeArgs []string, ref bind.Object, receiverParams bind.Object, clock bind.Object, denyList bind.Object, pool bind.Object, state bind.Object, messageTransmitterState bind.Object, treasury bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("release_or_mint", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"&mut offramp_sh::ReceiverParams",
		"&Clock",
		"&DenyList",
		"&mut USDCTokenPoolState<T>",
		"&mut MinterState",
		"&mut MessageTransmitterState",
		"&mut Treasury<T>",
	}, []any{
		ref,
		receiverParams,
		clock,
		denyList,
		pool,
		state,
		messageTransmitterState,
		treasury,
	}, nil)
}

// ReleaseOrMintWithArgs encodes a call to the release_or_mint Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c usdcTokenPoolEncoder) ReleaseOrMintWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"&mut offramp_sh::ReceiverParams",
		"&Clock",
		"&DenyList",
		"&mut USDCTokenPoolState<T>",
		"&mut MinterState",
		"&mut MessageTransmitterState",
		"&mut Treasury<T>",
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

// GetDomain encodes a call to the get_domain Move function.
func (c usdcTokenPoolEncoder) GetDomain(typeArgs []string, pool bind.Object, chainSelector uint64) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("get_domain", typeArgsList, typeParamsList, []string{
		"&USDCTokenPoolState<T>",
		"u64",
	}, []any{
		pool,
		chainSelector,
	}, []string{
		"usdc_token_pool::usdc_token_pool::Domain",
	})
}

// GetDomainWithArgs encodes a call to the get_domain Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c usdcTokenPoolEncoder) GetDomainWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&USDCTokenPoolState<T>",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("get_domain", typeArgsList, typeParamsList, expectedParams, args, []string{
		"usdc_token_pool::usdc_token_pool::Domain",
	})
}

// SetDomains encodes a call to the set_domains Move function.
func (c usdcTokenPoolEncoder) SetDomains(typeArgs []string, pool bind.Object, ownerCap bind.Object, remoteChainSelectors []uint64, remoteDomainIdentifiers []uint32, allowedRemoteCallers [][]byte, enableds []bool) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("set_domains", typeArgsList, typeParamsList, []string{
		"&mut USDCTokenPoolState<T>",
		"&OwnerCap",
		"vector<u64>",
		"vector<u32>",
		"vector<vector<u8>>",
		"vector<bool>",
	}, []any{
		pool,
		ownerCap,
		remoteChainSelectors,
		remoteDomainIdentifiers,
		allowedRemoteCallers,
		enableds,
	}, nil)
}

// SetDomainsWithArgs encodes a call to the set_domains Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c usdcTokenPoolEncoder) SetDomainsWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut USDCTokenPoolState<T>",
		"&OwnerCap",
		"vector<u64>",
		"vector<u32>",
		"vector<vector<u8>>",
		"vector<bool>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("set_domains", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// McmsSetDomains encodes a call to the mcms_set_domains Move function.
func (c usdcTokenPoolEncoder) McmsSetDomains(typeArgs []string, pool bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("mcms_set_domains", typeArgsList, typeParamsList, []string{
		"&mut USDCTokenPoolState<T>",
		"&mut Registry",
		"ExecutingCallbackParams",
	}, []any{
		pool,
		registry,
		params,
	}, nil)
}

// McmsSetDomainsWithArgs encodes a call to the mcms_set_domains Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c usdcTokenPoolEncoder) McmsSetDomainsWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut USDCTokenPoolState<T>",
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
	return c.EncodeCallArgsWithGenerics("mcms_set_domains", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// SetChainRateLimiterConfigs encodes a call to the set_chain_rate_limiter_configs Move function.
func (c usdcTokenPoolEncoder) SetChainRateLimiterConfigs(typeArgs []string, state bind.Object, ownerCap bind.Object, clock bind.Object, remoteChainSelectors []uint64, outboundIsEnableds []bool, outboundCapacities []uint64, outboundRates []uint64, inboundIsEnableds []bool, inboundCapacities []uint64, inboundRates []uint64) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("set_chain_rate_limiter_configs", typeArgsList, typeParamsList, []string{
		"&mut USDCTokenPoolState<T>",
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
func (c usdcTokenPoolEncoder) SetChainRateLimiterConfigsWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut USDCTokenPoolState<T>",
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
func (c usdcTokenPoolEncoder) SetChainRateLimiterConfig(typeArgs []string, state bind.Object, ownerCap bind.Object, clock bind.Object, remoteChainSelector uint64, outboundIsEnabled bool, outboundCapacity uint64, outboundRate uint64, inboundIsEnabled bool, inboundCapacity uint64, inboundRate uint64) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("set_chain_rate_limiter_config", typeArgsList, typeParamsList, []string{
		"&mut USDCTokenPoolState<T>",
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
func (c usdcTokenPoolEncoder) SetChainRateLimiterConfigWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut USDCTokenPoolState<T>",
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
func (c usdcTokenPoolEncoder) GetCurrentInboundRateLimiterState(typeArgs []string, clock bind.Object, state bind.Object, remoteChainSelector uint64) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("get_current_inbound_rate_limiter_state", typeArgsList, typeParamsList, []string{
		"&Clock",
		"&USDCTokenPoolState<T>",
		"u64",
	}, []any{
		clock,
		state,
		remoteChainSelector,
	}, []string{
		"usdc_token_pool::usdc_token_pool::TokenBucketWrapper",
	})
}

// GetCurrentInboundRateLimiterStateWithArgs encodes a call to the get_current_inbound_rate_limiter_state Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c usdcTokenPoolEncoder) GetCurrentInboundRateLimiterStateWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&Clock",
		"&USDCTokenPoolState<T>",
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
		"usdc_token_pool::usdc_token_pool::TokenBucketWrapper",
	})
}

// GetCurrentOutboundRateLimiterState encodes a call to the get_current_outbound_rate_limiter_state Move function.
func (c usdcTokenPoolEncoder) GetCurrentOutboundRateLimiterState(typeArgs []string, clock bind.Object, state bind.Object, remoteChainSelector uint64) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("get_current_outbound_rate_limiter_state", typeArgsList, typeParamsList, []string{
		"&Clock",
		"&USDCTokenPoolState<T>",
		"u64",
	}, []any{
		clock,
		state,
		remoteChainSelector,
	}, []string{
		"usdc_token_pool::usdc_token_pool::TokenBucketWrapper",
	})
}

// GetCurrentOutboundRateLimiterStateWithArgs encodes a call to the get_current_outbound_rate_limiter_state Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c usdcTokenPoolEncoder) GetCurrentOutboundRateLimiterStateWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&Clock",
		"&USDCTokenPoolState<T>",
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
		"usdc_token_pool::usdc_token_pool::TokenBucketWrapper",
	})
}

// Owner encodes a call to the owner Move function.
func (c usdcTokenPoolEncoder) Owner(typeArgs []string, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("owner", typeArgsList, typeParamsList, []string{
		"&USDCTokenPoolState<T>",
	}, []any{
		state,
	}, []string{
		"address",
	})
}

// OwnerWithArgs encodes a call to the owner Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c usdcTokenPoolEncoder) OwnerWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&USDCTokenPoolState<T>",
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
func (c usdcTokenPoolEncoder) HasPendingTransfer(typeArgs []string, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("has_pending_transfer", typeArgsList, typeParamsList, []string{
		"&USDCTokenPoolState<T>",
	}, []any{
		state,
	}, []string{
		"bool",
	})
}

// HasPendingTransferWithArgs encodes a call to the has_pending_transfer Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c usdcTokenPoolEncoder) HasPendingTransferWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&USDCTokenPoolState<T>",
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
func (c usdcTokenPoolEncoder) PendingTransferFrom(typeArgs []string, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("pending_transfer_from", typeArgsList, typeParamsList, []string{
		"&USDCTokenPoolState<T>",
	}, []any{
		state,
	}, []string{
		"0x1::option::Option<address>",
	})
}

// PendingTransferFromWithArgs encodes a call to the pending_transfer_from Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c usdcTokenPoolEncoder) PendingTransferFromWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&USDCTokenPoolState<T>",
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
func (c usdcTokenPoolEncoder) PendingTransferTo(typeArgs []string, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("pending_transfer_to", typeArgsList, typeParamsList, []string{
		"&USDCTokenPoolState<T>",
	}, []any{
		state,
	}, []string{
		"0x1::option::Option<address>",
	})
}

// PendingTransferToWithArgs encodes a call to the pending_transfer_to Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c usdcTokenPoolEncoder) PendingTransferToWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&USDCTokenPoolState<T>",
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
func (c usdcTokenPoolEncoder) PendingTransferAccepted(typeArgs []string, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("pending_transfer_accepted", typeArgsList, typeParamsList, []string{
		"&USDCTokenPoolState<T>",
	}, []any{
		state,
	}, []string{
		"0x1::option::Option<bool>",
	})
}

// PendingTransferAcceptedWithArgs encodes a call to the pending_transfer_accepted Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c usdcTokenPoolEncoder) PendingTransferAcceptedWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&USDCTokenPoolState<T>",
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
func (c usdcTokenPoolEncoder) TransferOwnership(typeArgs []string, state bind.Object, ownerCap bind.Object, newOwner string) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("transfer_ownership", typeArgsList, typeParamsList, []string{
		"&mut USDCTokenPoolState<T>",
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
func (c usdcTokenPoolEncoder) TransferOwnershipWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut USDCTokenPoolState<T>",
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
func (c usdcTokenPoolEncoder) AcceptOwnership(typeArgs []string, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("accept_ownership", typeArgsList, typeParamsList, []string{
		"&mut USDCTokenPoolState<T>",
	}, []any{
		state,
	}, nil)
}

// AcceptOwnershipWithArgs encodes a call to the accept_ownership Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c usdcTokenPoolEncoder) AcceptOwnershipWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut USDCTokenPoolState<T>",
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
func (c usdcTokenPoolEncoder) AcceptOwnershipFromObject(typeArgs []string, state bind.Object, from string) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("accept_ownership_from_object", typeArgsList, typeParamsList, []string{
		"&mut USDCTokenPoolState<T>",
		"&mut UID",
	}, []any{
		state,
		from,
	}, nil)
}

// AcceptOwnershipFromObjectWithArgs encodes a call to the accept_ownership_from_object Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c usdcTokenPoolEncoder) AcceptOwnershipFromObjectWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut USDCTokenPoolState<T>",
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
func (c usdcTokenPoolEncoder) McmsAcceptOwnership(typeArgs []string, state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("mcms_accept_ownership", typeArgsList, typeParamsList, []string{
		"&mut USDCTokenPoolState<T>",
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
func (c usdcTokenPoolEncoder) McmsAcceptOwnershipWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut USDCTokenPoolState<T>",
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
func (c usdcTokenPoolEncoder) ExecuteOwnershipTransfer(typeArgs []string, ownerCap bind.Object, state bind.Object, to string) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("execute_ownership_transfer", typeArgsList, typeParamsList, []string{
		"OwnerCap",
		"&mut USDCTokenPoolState<T>",
		"address",
	}, []any{
		ownerCap,
		state,
		to,
	}, nil)
}

// ExecuteOwnershipTransferWithArgs encodes a call to the execute_ownership_transfer Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c usdcTokenPoolEncoder) ExecuteOwnershipTransferWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"OwnerCap",
		"&mut USDCTokenPoolState<T>",
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
func (c usdcTokenPoolEncoder) ExecuteOwnershipTransferToMcms(typeArgs []string, ownerCap bind.Object, state bind.Object, registry bind.Object, to string) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("execute_ownership_transfer_to_mcms", typeArgsList, typeParamsList, []string{
		"OwnerCap",
		"&mut USDCTokenPoolState<T>",
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
func (c usdcTokenPoolEncoder) ExecuteOwnershipTransferToMcmsWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"OwnerCap",
		"&mut USDCTokenPoolState<T>",
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
func (c usdcTokenPoolEncoder) McmsRegisterUpgradeCap(upgradeCap bind.Object, registry bind.Object, state bind.Object) (*bind.EncodedCall, error) {
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
func (c usdcTokenPoolEncoder) McmsRegisterUpgradeCapWithArgs(args ...any) (*bind.EncodedCall, error) {
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
func (c usdcTokenPoolEncoder) McmsSetAllowlistEnabled(typeArgs []string, state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("mcms_set_allowlist_enabled", typeArgsList, typeParamsList, []string{
		"&mut USDCTokenPoolState<T>",
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
func (c usdcTokenPoolEncoder) McmsSetAllowlistEnabledWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut USDCTokenPoolState<T>",
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
func (c usdcTokenPoolEncoder) McmsApplyAllowlistUpdates(typeArgs []string, state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("mcms_apply_allowlist_updates", typeArgsList, typeParamsList, []string{
		"&mut USDCTokenPoolState<T>",
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
func (c usdcTokenPoolEncoder) McmsApplyAllowlistUpdatesWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut USDCTokenPoolState<T>",
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
func (c usdcTokenPoolEncoder) McmsApplyChainUpdates(typeArgs []string, state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("mcms_apply_chain_updates", typeArgsList, typeParamsList, []string{
		"&mut USDCTokenPoolState<T>",
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
func (c usdcTokenPoolEncoder) McmsApplyChainUpdatesWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut USDCTokenPoolState<T>",
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
func (c usdcTokenPoolEncoder) McmsAddRemotePool(typeArgs []string, state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("mcms_add_remote_pool", typeArgsList, typeParamsList, []string{
		"&mut USDCTokenPoolState<T>",
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
func (c usdcTokenPoolEncoder) McmsAddRemotePoolWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut USDCTokenPoolState<T>",
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
func (c usdcTokenPoolEncoder) McmsRemoveRemotePool(typeArgs []string, state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("mcms_remove_remote_pool", typeArgsList, typeParamsList, []string{
		"&mut USDCTokenPoolState<T>",
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
func (c usdcTokenPoolEncoder) McmsRemoveRemotePoolWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut USDCTokenPoolState<T>",
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
func (c usdcTokenPoolEncoder) McmsSetChainRateLimiterConfigs(typeArgs []string, state bind.Object, registry bind.Object, params bind.Object, clock bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("mcms_set_chain_rate_limiter_configs", typeArgsList, typeParamsList, []string{
		"&mut USDCTokenPoolState<T>",
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
func (c usdcTokenPoolEncoder) McmsSetChainRateLimiterConfigsWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut USDCTokenPoolState<T>",
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
func (c usdcTokenPoolEncoder) McmsSetChainRateLimiterConfig(typeArgs []string, state bind.Object, registry bind.Object, params bind.Object, clock bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("mcms_set_chain_rate_limiter_config", typeArgsList, typeParamsList, []string{
		"&mut USDCTokenPoolState<T>",
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
func (c usdcTokenPoolEncoder) McmsSetChainRateLimiterConfigWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut USDCTokenPoolState<T>",
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
func (c usdcTokenPoolEncoder) DestroyTokenPool(typeArgs []string, ref bind.Object, state bind.Object, ownerCap bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("destroy_token_pool", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"USDCTokenPoolState<T>",
		"OwnerCap",
	}, []any{
		ref,
		state,
		ownerCap,
	}, nil)
}

// DestroyTokenPoolWithArgs encodes a call to the destroy_token_pool Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c usdcTokenPoolEncoder) DestroyTokenPoolWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"USDCTokenPoolState<T>",
		"OwnerCap",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("destroy_token_pool", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// McmsDestroyTokenPool encodes a call to the mcms_destroy_token_pool Move function.
func (c usdcTokenPoolEncoder) McmsDestroyTokenPool(typeArgs []string, ref bind.Object, state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("mcms_destroy_token_pool", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"USDCTokenPoolState<T>",
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
func (c usdcTokenPoolEncoder) McmsDestroyTokenPoolWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"USDCTokenPoolState<T>",
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
func (c usdcTokenPoolEncoder) McmsTransferOwnership(typeArgs []string, state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("mcms_transfer_ownership", typeArgsList, typeParamsList, []string{
		"&mut USDCTokenPoolState<T>",
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
func (c usdcTokenPoolEncoder) McmsTransferOwnershipWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut USDCTokenPoolState<T>",
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
func (c usdcTokenPoolEncoder) McmsExecuteOwnershipTransfer(typeArgs []string, state bind.Object, registry bind.Object, deployerState bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("mcms_execute_ownership_transfer", typeArgsList, typeParamsList, []string{
		"&mut USDCTokenPoolState<T>",
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
func (c usdcTokenPoolEncoder) McmsExecuteOwnershipTransferWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut USDCTokenPoolState<T>",
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
func (c usdcTokenPoolEncoder) McmsAddAllowedModules(typeArgs []string, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
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
func (c usdcTokenPoolEncoder) McmsAddAllowedModulesWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
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
func (c usdcTokenPoolEncoder) McmsRemoveAllowedModules(typeArgs []string, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
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
func (c usdcTokenPoolEncoder) McmsRemoveAllowedModulesWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
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
