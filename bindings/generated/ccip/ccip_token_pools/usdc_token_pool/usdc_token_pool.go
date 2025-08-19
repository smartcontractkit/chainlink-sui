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

type IUsdcTokenPool interface {
	TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	Initialize(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, ccipAdminProof bind.Object, coinMetadata bind.Object, localDomainIdentifier uint32, tokenPoolPackageId string, tokenPoolAdministrator string) (*models.SuiTransactionBlockResponse, error)
	GetToken(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	GetTokenDecimals(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	GetRemotePools(ctx context.Context, opts *bind.CallOpts, state bind.Object, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	IsRemotePool(ctx context.Context, opts *bind.CallOpts, state bind.Object, remoteChainSelector uint64, remotePoolAddress []byte) (*models.SuiTransactionBlockResponse, error)
	GetRemoteToken(ctx context.Context, opts *bind.CallOpts, state bind.Object, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	AddRemotePool(ctx context.Context, opts *bind.CallOpts, state bind.Object, ownerCap bind.Object, remoteChainSelector uint64, remotePoolAddress []byte) (*models.SuiTransactionBlockResponse, error)
	RemoveRemotePool(ctx context.Context, opts *bind.CallOpts, state bind.Object, ownerCap bind.Object, remoteChainSelector uint64, remotePoolAddress []byte) (*models.SuiTransactionBlockResponse, error)
	IsSupportedChain(ctx context.Context, opts *bind.CallOpts, state bind.Object, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	GetSupportedChains(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	ApplyChainUpdates(ctx context.Context, opts *bind.CallOpts, state bind.Object, ownerCap bind.Object, remoteChainSelectorsToRemove []uint64, remoteChainSelectorsToAdd []uint64, remotePoolAddressesToAdd [][][]byte, remoteTokenAddressesToAdd [][]byte) (*models.SuiTransactionBlockResponse, error)
	GetAllowlistEnabled(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	GetAllowlist(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	SetAllowlistEnabled(ctx context.Context, opts *bind.CallOpts, state bind.Object, ownerCap bind.Object, enabled bool) (*models.SuiTransactionBlockResponse, error)
	ApplyAllowlistUpdates(ctx context.Context, opts *bind.CallOpts, state bind.Object, ownerCap bind.Object, removes []string, adds []string) (*models.SuiTransactionBlockResponse, error)
	GetPackageAuthCaller(ctx context.Context, opts *bind.CallOpts, typeArgs []string) (*models.SuiTransactionBlockResponse, error)
	LockOrBurn(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, tokenTransferParams bind.Object, c_ bind.Object, remoteChainSelector uint64, receiver string, clock bind.Object, denyList bind.Object, pool bind.Object, state bind.Object, messageTransmitterState bind.Object, treasury bind.Object) (*models.SuiTransactionBlockResponse, error)
	ReleaseOrMint(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, receiverParams bind.Object, tokenTransfer bind.Object, clock bind.Object, denyList bind.Object, pool bind.Object, state bind.Object, messageTransmitterState bind.Object, treasury bind.Object) (*models.SuiTransactionBlockResponse, error)
	GetDomain(ctx context.Context, opts *bind.CallOpts, pool bind.Object, chainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	SetDomains(ctx context.Context, opts *bind.CallOpts, pool bind.Object, ownerCap bind.Object, remoteChainSelectors []uint64, remoteDomainIdentifiers []uint32, allowedRemoteCallers [][]byte, enableds []bool) (*models.SuiTransactionBlockResponse, error)
	SetChainRateLimiterConfigs(ctx context.Context, opts *bind.CallOpts, state bind.Object, ownerCap bind.Object, clock bind.Object, remoteChainSelectors []uint64, outboundIsEnableds []bool, outboundCapacities []uint64, outboundRates []uint64, inboundIsEnableds []bool, inboundCapacities []uint64, inboundRates []uint64) (*models.SuiTransactionBlockResponse, error)
	SetChainRateLimiterConfig(ctx context.Context, opts *bind.CallOpts, state bind.Object, ownerCap bind.Object, clock bind.Object, remoteChainSelector uint64, outboundIsEnabled bool, outboundCapacity uint64, outboundRate uint64, inboundIsEnabled bool, inboundCapacity uint64, inboundRate uint64) (*models.SuiTransactionBlockResponse, error)
	Owner(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	HasPendingTransfer(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	PendingTransferFrom(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	PendingTransferTo(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	PendingTransferAccepted(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	TransferOwnership(ctx context.Context, opts *bind.CallOpts, state bind.Object, ownerCap bind.Object, newOwner string) (*models.SuiTransactionBlockResponse, error)
	AcceptOwnership(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	AcceptOwnershipFromObject(ctx context.Context, opts *bind.CallOpts, state bind.Object, from string) (*models.SuiTransactionBlockResponse, error)
	ExecuteOwnershipTransfer(ctx context.Context, opts *bind.CallOpts, ownerCap bind.Object, ownableState bind.Object, to string) (*models.SuiTransactionBlockResponse, error)
	ExecuteOwnershipTransferToMcms(ctx context.Context, opts *bind.CallOpts, ownerCap bind.Object, state bind.Object, registry bind.Object, to string) (*models.SuiTransactionBlockResponse, error)
	McmsRegisterUpgradeCap(ctx context.Context, opts *bind.CallOpts, upgradeCap bind.Object, registry bind.Object, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsEntrypoint(ctx context.Context, opts *bind.CallOpts, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	DestroyTokenPool(ctx context.Context, opts *bind.CallOpts, state bind.Object, ownerCap bind.Object) (*models.SuiTransactionBlockResponse, error)
	DevInspect() IUsdcTokenPoolDevInspect
	Encoder() UsdcTokenPoolEncoder
}

type IUsdcTokenPoolDevInspect interface {
	TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (string, error)
	GetToken(ctx context.Context, opts *bind.CallOpts, state bind.Object) (string, error)
	GetTokenDecimals(ctx context.Context, opts *bind.CallOpts, state bind.Object) (byte, error)
	GetRemotePools(ctx context.Context, opts *bind.CallOpts, state bind.Object, remoteChainSelector uint64) ([][]byte, error)
	IsRemotePool(ctx context.Context, opts *bind.CallOpts, state bind.Object, remoteChainSelector uint64, remotePoolAddress []byte) (bool, error)
	GetRemoteToken(ctx context.Context, opts *bind.CallOpts, state bind.Object, remoteChainSelector uint64) ([]byte, error)
	IsSupportedChain(ctx context.Context, opts *bind.CallOpts, state bind.Object, remoteChainSelector uint64) (bool, error)
	GetSupportedChains(ctx context.Context, opts *bind.CallOpts, state bind.Object) ([]uint64, error)
	GetAllowlistEnabled(ctx context.Context, opts *bind.CallOpts, state bind.Object) (bool, error)
	GetAllowlist(ctx context.Context, opts *bind.CallOpts, state bind.Object) ([]string, error)
	GetPackageAuthCaller(ctx context.Context, opts *bind.CallOpts, typeArgs []string) (string, error)
	GetDomain(ctx context.Context, opts *bind.CallOpts, pool bind.Object, chainSelector uint64) (Domain, error)
	Owner(ctx context.Context, opts *bind.CallOpts, state bind.Object) (string, error)
	HasPendingTransfer(ctx context.Context, opts *bind.CallOpts, state bind.Object) (bool, error)
	PendingTransferFrom(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*string, error)
	PendingTransferTo(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*string, error)
	PendingTransferAccepted(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*bool, error)
}

type UsdcTokenPoolEncoder interface {
	TypeAndVersion() (*bind.EncodedCall, error)
	TypeAndVersionWithArgs(args ...any) (*bind.EncodedCall, error)
	Initialize(typeArgs []string, ref bind.Object, ccipAdminProof bind.Object, coinMetadata bind.Object, localDomainIdentifier uint32, tokenPoolPackageId string, tokenPoolAdministrator string) (*bind.EncodedCall, error)
	InitializeWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	GetToken(state bind.Object) (*bind.EncodedCall, error)
	GetTokenWithArgs(args ...any) (*bind.EncodedCall, error)
	GetTokenDecimals(state bind.Object) (*bind.EncodedCall, error)
	GetTokenDecimalsWithArgs(args ...any) (*bind.EncodedCall, error)
	GetRemotePools(state bind.Object, remoteChainSelector uint64) (*bind.EncodedCall, error)
	GetRemotePoolsWithArgs(args ...any) (*bind.EncodedCall, error)
	IsRemotePool(state bind.Object, remoteChainSelector uint64, remotePoolAddress []byte) (*bind.EncodedCall, error)
	IsRemotePoolWithArgs(args ...any) (*bind.EncodedCall, error)
	GetRemoteToken(state bind.Object, remoteChainSelector uint64) (*bind.EncodedCall, error)
	GetRemoteTokenWithArgs(args ...any) (*bind.EncodedCall, error)
	AddRemotePool(state bind.Object, ownerCap bind.Object, remoteChainSelector uint64, remotePoolAddress []byte) (*bind.EncodedCall, error)
	AddRemotePoolWithArgs(args ...any) (*bind.EncodedCall, error)
	RemoveRemotePool(state bind.Object, ownerCap bind.Object, remoteChainSelector uint64, remotePoolAddress []byte) (*bind.EncodedCall, error)
	RemoveRemotePoolWithArgs(args ...any) (*bind.EncodedCall, error)
	IsSupportedChain(state bind.Object, remoteChainSelector uint64) (*bind.EncodedCall, error)
	IsSupportedChainWithArgs(args ...any) (*bind.EncodedCall, error)
	GetSupportedChains(state bind.Object) (*bind.EncodedCall, error)
	GetSupportedChainsWithArgs(args ...any) (*bind.EncodedCall, error)
	ApplyChainUpdates(state bind.Object, ownerCap bind.Object, remoteChainSelectorsToRemove []uint64, remoteChainSelectorsToAdd []uint64, remotePoolAddressesToAdd [][][]byte, remoteTokenAddressesToAdd [][]byte) (*bind.EncodedCall, error)
	ApplyChainUpdatesWithArgs(args ...any) (*bind.EncodedCall, error)
	GetAllowlistEnabled(state bind.Object) (*bind.EncodedCall, error)
	GetAllowlistEnabledWithArgs(args ...any) (*bind.EncodedCall, error)
	GetAllowlist(state bind.Object) (*bind.EncodedCall, error)
	GetAllowlistWithArgs(args ...any) (*bind.EncodedCall, error)
	SetAllowlistEnabled(state bind.Object, ownerCap bind.Object, enabled bool) (*bind.EncodedCall, error)
	SetAllowlistEnabledWithArgs(args ...any) (*bind.EncodedCall, error)
	ApplyAllowlistUpdates(state bind.Object, ownerCap bind.Object, removes []string, adds []string) (*bind.EncodedCall, error)
	ApplyAllowlistUpdatesWithArgs(args ...any) (*bind.EncodedCall, error)
	GetPackageAuthCaller(typeArgs []string) (*bind.EncodedCall, error)
	GetPackageAuthCallerWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	LockOrBurn(typeArgs []string, ref bind.Object, tokenTransferParams bind.Object, c_ bind.Object, remoteChainSelector uint64, receiver string, clock bind.Object, denyList bind.Object, pool bind.Object, state bind.Object, messageTransmitterState bind.Object, treasury bind.Object) (*bind.EncodedCall, error)
	LockOrBurnWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	ReleaseOrMint(typeArgs []string, ref bind.Object, receiverParams bind.Object, tokenTransfer bind.Object, clock bind.Object, denyList bind.Object, pool bind.Object, state bind.Object, messageTransmitterState bind.Object, treasury bind.Object) (*bind.EncodedCall, error)
	ReleaseOrMintWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	GetDomain(pool bind.Object, chainSelector uint64) (*bind.EncodedCall, error)
	GetDomainWithArgs(args ...any) (*bind.EncodedCall, error)
	SetDomains(pool bind.Object, ownerCap bind.Object, remoteChainSelectors []uint64, remoteDomainIdentifiers []uint32, allowedRemoteCallers [][]byte, enableds []bool) (*bind.EncodedCall, error)
	SetDomainsWithArgs(args ...any) (*bind.EncodedCall, error)
	SetChainRateLimiterConfigs(state bind.Object, ownerCap bind.Object, clock bind.Object, remoteChainSelectors []uint64, outboundIsEnableds []bool, outboundCapacities []uint64, outboundRates []uint64, inboundIsEnableds []bool, inboundCapacities []uint64, inboundRates []uint64) (*bind.EncodedCall, error)
	SetChainRateLimiterConfigsWithArgs(args ...any) (*bind.EncodedCall, error)
	SetChainRateLimiterConfig(state bind.Object, ownerCap bind.Object, clock bind.Object, remoteChainSelector uint64, outboundIsEnabled bool, outboundCapacity uint64, outboundRate uint64, inboundIsEnabled bool, inboundCapacity uint64, inboundRate uint64) (*bind.EncodedCall, error)
	SetChainRateLimiterConfigWithArgs(args ...any) (*bind.EncodedCall, error)
	Owner(state bind.Object) (*bind.EncodedCall, error)
	OwnerWithArgs(args ...any) (*bind.EncodedCall, error)
	HasPendingTransfer(state bind.Object) (*bind.EncodedCall, error)
	HasPendingTransferWithArgs(args ...any) (*bind.EncodedCall, error)
	PendingTransferFrom(state bind.Object) (*bind.EncodedCall, error)
	PendingTransferFromWithArgs(args ...any) (*bind.EncodedCall, error)
	PendingTransferTo(state bind.Object) (*bind.EncodedCall, error)
	PendingTransferToWithArgs(args ...any) (*bind.EncodedCall, error)
	PendingTransferAccepted(state bind.Object) (*bind.EncodedCall, error)
	PendingTransferAcceptedWithArgs(args ...any) (*bind.EncodedCall, error)
	TransferOwnership(state bind.Object, ownerCap bind.Object, newOwner string) (*bind.EncodedCall, error)
	TransferOwnershipWithArgs(args ...any) (*bind.EncodedCall, error)
	AcceptOwnership(state bind.Object) (*bind.EncodedCall, error)
	AcceptOwnershipWithArgs(args ...any) (*bind.EncodedCall, error)
	AcceptOwnershipFromObject(state bind.Object, from string) (*bind.EncodedCall, error)
	AcceptOwnershipFromObjectWithArgs(args ...any) (*bind.EncodedCall, error)
	ExecuteOwnershipTransfer(ownerCap bind.Object, ownableState bind.Object, to string) (*bind.EncodedCall, error)
	ExecuteOwnershipTransferWithArgs(args ...any) (*bind.EncodedCall, error)
	ExecuteOwnershipTransferToMcms(ownerCap bind.Object, state bind.Object, registry bind.Object, to string) (*bind.EncodedCall, error)
	ExecuteOwnershipTransferToMcmsWithArgs(args ...any) (*bind.EncodedCall, error)
	McmsRegisterUpgradeCap(upgradeCap bind.Object, registry bind.Object, state bind.Object) (*bind.EncodedCall, error)
	McmsRegisterUpgradeCapWithArgs(args ...any) (*bind.EncodedCall, error)
	McmsEntrypoint(state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsEntrypointWithArgs(args ...any) (*bind.EncodedCall, error)
	DestroyTokenPool(state bind.Object, ownerCap bind.Object) (*bind.EncodedCall, error)
	DestroyTokenPoolWithArgs(args ...any) (*bind.EncodedCall, error)
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

func NewUsdcTokenPool(packageID string, client sui.ISuiAPI) (*UsdcTokenPoolContract, error) {
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

func (c *UsdcTokenPoolContract) Encoder() UsdcTokenPoolEncoder {
	return c.usdcTokenPoolEncoder
}

func (c *UsdcTokenPoolContract) DevInspect() IUsdcTokenPoolDevInspect {
	return c.devInspect
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

type TypeProof struct {
}

type McmsCallback struct {
}

func init() {
	bind.RegisterStructDecoder("usdc_token_pool::usdc_token_pool::Domain", func(data []byte) (interface{}, error) {
		var result Domain
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("usdc_token_pool::usdc_token_pool::DomainsSet", func(data []byte) (interface{}, error) {
		var result DomainsSet
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("usdc_token_pool::usdc_token_pool::USDCTokenPoolState", func(data []byte) (interface{}, error) {
		var result USDCTokenPoolState
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("usdc_token_pool::usdc_token_pool::TypeProof", func(data []byte) (interface{}, error) {
		var result TypeProof
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("usdc_token_pool::usdc_token_pool::McmsCallback", func(data []byte) (interface{}, error) {
		var result McmsCallback
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
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
func (c *UsdcTokenPoolContract) Initialize(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, ccipAdminProof bind.Object, coinMetadata bind.Object, localDomainIdentifier uint32, tokenPoolPackageId string, tokenPoolAdministrator string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.Initialize(typeArgs, ref, ccipAdminProof, coinMetadata, localDomainIdentifier, tokenPoolPackageId, tokenPoolAdministrator)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetToken executes the get_token Move function.
func (c *UsdcTokenPoolContract) GetToken(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.GetToken(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetTokenDecimals executes the get_token_decimals Move function.
func (c *UsdcTokenPoolContract) GetTokenDecimals(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.GetTokenDecimals(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetRemotePools executes the get_remote_pools Move function.
func (c *UsdcTokenPoolContract) GetRemotePools(ctx context.Context, opts *bind.CallOpts, state bind.Object, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.GetRemotePools(state, remoteChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// IsRemotePool executes the is_remote_pool Move function.
func (c *UsdcTokenPoolContract) IsRemotePool(ctx context.Context, opts *bind.CallOpts, state bind.Object, remoteChainSelector uint64, remotePoolAddress []byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.IsRemotePool(state, remoteChainSelector, remotePoolAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetRemoteToken executes the get_remote_token Move function.
func (c *UsdcTokenPoolContract) GetRemoteToken(ctx context.Context, opts *bind.CallOpts, state bind.Object, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.GetRemoteToken(state, remoteChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// AddRemotePool executes the add_remote_pool Move function.
func (c *UsdcTokenPoolContract) AddRemotePool(ctx context.Context, opts *bind.CallOpts, state bind.Object, ownerCap bind.Object, remoteChainSelector uint64, remotePoolAddress []byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.AddRemotePool(state, ownerCap, remoteChainSelector, remotePoolAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// RemoveRemotePool executes the remove_remote_pool Move function.
func (c *UsdcTokenPoolContract) RemoveRemotePool(ctx context.Context, opts *bind.CallOpts, state bind.Object, ownerCap bind.Object, remoteChainSelector uint64, remotePoolAddress []byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.RemoveRemotePool(state, ownerCap, remoteChainSelector, remotePoolAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// IsSupportedChain executes the is_supported_chain Move function.
func (c *UsdcTokenPoolContract) IsSupportedChain(ctx context.Context, opts *bind.CallOpts, state bind.Object, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.IsSupportedChain(state, remoteChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetSupportedChains executes the get_supported_chains Move function.
func (c *UsdcTokenPoolContract) GetSupportedChains(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.GetSupportedChains(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ApplyChainUpdates executes the apply_chain_updates Move function.
func (c *UsdcTokenPoolContract) ApplyChainUpdates(ctx context.Context, opts *bind.CallOpts, state bind.Object, ownerCap bind.Object, remoteChainSelectorsToRemove []uint64, remoteChainSelectorsToAdd []uint64, remotePoolAddressesToAdd [][][]byte, remoteTokenAddressesToAdd [][]byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.ApplyChainUpdates(state, ownerCap, remoteChainSelectorsToRemove, remoteChainSelectorsToAdd, remotePoolAddressesToAdd, remoteTokenAddressesToAdd)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetAllowlistEnabled executes the get_allowlist_enabled Move function.
func (c *UsdcTokenPoolContract) GetAllowlistEnabled(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.GetAllowlistEnabled(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetAllowlist executes the get_allowlist Move function.
func (c *UsdcTokenPoolContract) GetAllowlist(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.GetAllowlist(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// SetAllowlistEnabled executes the set_allowlist_enabled Move function.
func (c *UsdcTokenPoolContract) SetAllowlistEnabled(ctx context.Context, opts *bind.CallOpts, state bind.Object, ownerCap bind.Object, enabled bool) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.SetAllowlistEnabled(state, ownerCap, enabled)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ApplyAllowlistUpdates executes the apply_allowlist_updates Move function.
func (c *UsdcTokenPoolContract) ApplyAllowlistUpdates(ctx context.Context, opts *bind.CallOpts, state bind.Object, ownerCap bind.Object, removes []string, adds []string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.ApplyAllowlistUpdates(state, ownerCap, removes, adds)
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
func (c *UsdcTokenPoolContract) LockOrBurn(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, tokenTransferParams bind.Object, c_ bind.Object, remoteChainSelector uint64, receiver string, clock bind.Object, denyList bind.Object, pool bind.Object, state bind.Object, messageTransmitterState bind.Object, treasury bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.LockOrBurn(typeArgs, ref, tokenTransferParams, c_, remoteChainSelector, receiver, clock, denyList, pool, state, messageTransmitterState, treasury)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ReleaseOrMint executes the release_or_mint Move function.
func (c *UsdcTokenPoolContract) ReleaseOrMint(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, receiverParams bind.Object, tokenTransfer bind.Object, clock bind.Object, denyList bind.Object, pool bind.Object, state bind.Object, messageTransmitterState bind.Object, treasury bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.ReleaseOrMint(typeArgs, ref, receiverParams, tokenTransfer, clock, denyList, pool, state, messageTransmitterState, treasury)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetDomain executes the get_domain Move function.
func (c *UsdcTokenPoolContract) GetDomain(ctx context.Context, opts *bind.CallOpts, pool bind.Object, chainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.GetDomain(pool, chainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// SetDomains executes the set_domains Move function.
func (c *UsdcTokenPoolContract) SetDomains(ctx context.Context, opts *bind.CallOpts, pool bind.Object, ownerCap bind.Object, remoteChainSelectors []uint64, remoteDomainIdentifiers []uint32, allowedRemoteCallers [][]byte, enableds []bool) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.SetDomains(pool, ownerCap, remoteChainSelectors, remoteDomainIdentifiers, allowedRemoteCallers, enableds)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// SetChainRateLimiterConfigs executes the set_chain_rate_limiter_configs Move function.
func (c *UsdcTokenPoolContract) SetChainRateLimiterConfigs(ctx context.Context, opts *bind.CallOpts, state bind.Object, ownerCap bind.Object, clock bind.Object, remoteChainSelectors []uint64, outboundIsEnableds []bool, outboundCapacities []uint64, outboundRates []uint64, inboundIsEnableds []bool, inboundCapacities []uint64, inboundRates []uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.SetChainRateLimiterConfigs(state, ownerCap, clock, remoteChainSelectors, outboundIsEnableds, outboundCapacities, outboundRates, inboundIsEnableds, inboundCapacities, inboundRates)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// SetChainRateLimiterConfig executes the set_chain_rate_limiter_config Move function.
func (c *UsdcTokenPoolContract) SetChainRateLimiterConfig(ctx context.Context, opts *bind.CallOpts, state bind.Object, ownerCap bind.Object, clock bind.Object, remoteChainSelector uint64, outboundIsEnabled bool, outboundCapacity uint64, outboundRate uint64, inboundIsEnabled bool, inboundCapacity uint64, inboundRate uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.SetChainRateLimiterConfig(state, ownerCap, clock, remoteChainSelector, outboundIsEnabled, outboundCapacity, outboundRate, inboundIsEnabled, inboundCapacity, inboundRate)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// Owner executes the owner Move function.
func (c *UsdcTokenPoolContract) Owner(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.Owner(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// HasPendingTransfer executes the has_pending_transfer Move function.
func (c *UsdcTokenPoolContract) HasPendingTransfer(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.HasPendingTransfer(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// PendingTransferFrom executes the pending_transfer_from Move function.
func (c *UsdcTokenPoolContract) PendingTransferFrom(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.PendingTransferFrom(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// PendingTransferTo executes the pending_transfer_to Move function.
func (c *UsdcTokenPoolContract) PendingTransferTo(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.PendingTransferTo(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// PendingTransferAccepted executes the pending_transfer_accepted Move function.
func (c *UsdcTokenPoolContract) PendingTransferAccepted(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.PendingTransferAccepted(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// TransferOwnership executes the transfer_ownership Move function.
func (c *UsdcTokenPoolContract) TransferOwnership(ctx context.Context, opts *bind.CallOpts, state bind.Object, ownerCap bind.Object, newOwner string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.TransferOwnership(state, ownerCap, newOwner)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// AcceptOwnership executes the accept_ownership Move function.
func (c *UsdcTokenPoolContract) AcceptOwnership(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.AcceptOwnership(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// AcceptOwnershipFromObject executes the accept_ownership_from_object Move function.
func (c *UsdcTokenPoolContract) AcceptOwnershipFromObject(ctx context.Context, opts *bind.CallOpts, state bind.Object, from string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.AcceptOwnershipFromObject(state, from)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ExecuteOwnershipTransfer executes the execute_ownership_transfer Move function.
func (c *UsdcTokenPoolContract) ExecuteOwnershipTransfer(ctx context.Context, opts *bind.CallOpts, ownerCap bind.Object, ownableState bind.Object, to string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.ExecuteOwnershipTransfer(ownerCap, ownableState, to)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ExecuteOwnershipTransferToMcms executes the execute_ownership_transfer_to_mcms Move function.
func (c *UsdcTokenPoolContract) ExecuteOwnershipTransferToMcms(ctx context.Context, opts *bind.CallOpts, ownerCap bind.Object, state bind.Object, registry bind.Object, to string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.ExecuteOwnershipTransferToMcms(ownerCap, state, registry, to)
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

// McmsEntrypoint executes the mcms_entrypoint Move function.
func (c *UsdcTokenPoolContract) McmsEntrypoint(ctx context.Context, opts *bind.CallOpts, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.McmsEntrypoint(state, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// DestroyTokenPool executes the destroy_token_pool Move function.
func (c *UsdcTokenPoolContract) DestroyTokenPool(ctx context.Context, opts *bind.CallOpts, state bind.Object, ownerCap bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.usdcTokenPoolEncoder.DestroyTokenPool(state, ownerCap)
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
func (d *UsdcTokenPoolDevInspect) GetToken(ctx context.Context, opts *bind.CallOpts, state bind.Object) (string, error) {
	encoded, err := d.contract.usdcTokenPoolEncoder.GetToken(state)
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
func (d *UsdcTokenPoolDevInspect) GetTokenDecimals(ctx context.Context, opts *bind.CallOpts, state bind.Object) (byte, error) {
	encoded, err := d.contract.usdcTokenPoolEncoder.GetTokenDecimals(state)
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
func (d *UsdcTokenPoolDevInspect) GetRemotePools(ctx context.Context, opts *bind.CallOpts, state bind.Object, remoteChainSelector uint64) ([][]byte, error) {
	encoded, err := d.contract.usdcTokenPoolEncoder.GetRemotePools(state, remoteChainSelector)
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
func (d *UsdcTokenPoolDevInspect) IsRemotePool(ctx context.Context, opts *bind.CallOpts, state bind.Object, remoteChainSelector uint64, remotePoolAddress []byte) (bool, error) {
	encoded, err := d.contract.usdcTokenPoolEncoder.IsRemotePool(state, remoteChainSelector, remotePoolAddress)
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
func (d *UsdcTokenPoolDevInspect) GetRemoteToken(ctx context.Context, opts *bind.CallOpts, state bind.Object, remoteChainSelector uint64) ([]byte, error) {
	encoded, err := d.contract.usdcTokenPoolEncoder.GetRemoteToken(state, remoteChainSelector)
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
func (d *UsdcTokenPoolDevInspect) IsSupportedChain(ctx context.Context, opts *bind.CallOpts, state bind.Object, remoteChainSelector uint64) (bool, error) {
	encoded, err := d.contract.usdcTokenPoolEncoder.IsSupportedChain(state, remoteChainSelector)
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
func (d *UsdcTokenPoolDevInspect) GetSupportedChains(ctx context.Context, opts *bind.CallOpts, state bind.Object) ([]uint64, error) {
	encoded, err := d.contract.usdcTokenPoolEncoder.GetSupportedChains(state)
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
func (d *UsdcTokenPoolDevInspect) GetAllowlistEnabled(ctx context.Context, opts *bind.CallOpts, state bind.Object) (bool, error) {
	encoded, err := d.contract.usdcTokenPoolEncoder.GetAllowlistEnabled(state)
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
func (d *UsdcTokenPoolDevInspect) GetAllowlist(ctx context.Context, opts *bind.CallOpts, state bind.Object) ([]string, error) {
	encoded, err := d.contract.usdcTokenPoolEncoder.GetAllowlist(state)
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
func (d *UsdcTokenPoolDevInspect) GetDomain(ctx context.Context, opts *bind.CallOpts, pool bind.Object, chainSelector uint64) (Domain, error) {
	encoded, err := d.contract.usdcTokenPoolEncoder.GetDomain(pool, chainSelector)
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

// Owner executes the owner Move function using DevInspect to get return values.
//
// Returns: address
func (d *UsdcTokenPoolDevInspect) Owner(ctx context.Context, opts *bind.CallOpts, state bind.Object) (string, error) {
	encoded, err := d.contract.usdcTokenPoolEncoder.Owner(state)
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
func (d *UsdcTokenPoolDevInspect) HasPendingTransfer(ctx context.Context, opts *bind.CallOpts, state bind.Object) (bool, error) {
	encoded, err := d.contract.usdcTokenPoolEncoder.HasPendingTransfer(state)
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
func (d *UsdcTokenPoolDevInspect) PendingTransferFrom(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*string, error) {
	encoded, err := d.contract.usdcTokenPoolEncoder.PendingTransferFrom(state)
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
func (d *UsdcTokenPoolDevInspect) PendingTransferTo(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*string, error) {
	encoded, err := d.contract.usdcTokenPoolEncoder.PendingTransferTo(state)
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
func (d *UsdcTokenPoolDevInspect) PendingTransferAccepted(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*bool, error) {
	encoded, err := d.contract.usdcTokenPoolEncoder.PendingTransferAccepted(state)
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
func (c usdcTokenPoolEncoder) Initialize(typeArgs []string, ref bind.Object, ccipAdminProof bind.Object, coinMetadata bind.Object, localDomainIdentifier uint32, tokenPoolPackageId string, tokenPoolAdministrator string) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("initialize", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"state_object::CCIPAdminProof",
		"&CoinMetadata<T>",
		"u32",
		"address",
		"address",
	}, []any{
		ref,
		ccipAdminProof,
		coinMetadata,
		localDomainIdentifier,
		tokenPoolPackageId,
		tokenPoolAdministrator,
	}, nil)
}

// InitializeWithArgs encodes a call to the initialize Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c usdcTokenPoolEncoder) InitializeWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"state_object::CCIPAdminProof",
		"&CoinMetadata<T>",
		"u32",
		"address",
		"address",
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
func (c usdcTokenPoolEncoder) GetToken(state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_token", typeArgsList, typeParamsList, []string{
		"&USDCTokenPoolState",
	}, []any{
		state,
	}, []string{
		"address",
	})
}

// GetTokenWithArgs encodes a call to the get_token Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c usdcTokenPoolEncoder) GetTokenWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&USDCTokenPoolState",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_token", typeArgsList, typeParamsList, expectedParams, args, []string{
		"address",
	})
}

// GetTokenDecimals encodes a call to the get_token_decimals Move function.
func (c usdcTokenPoolEncoder) GetTokenDecimals(state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_token_decimals", typeArgsList, typeParamsList, []string{
		"&USDCTokenPoolState",
	}, []any{
		state,
	}, []string{
		"u8",
	})
}

// GetTokenDecimalsWithArgs encodes a call to the get_token_decimals Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c usdcTokenPoolEncoder) GetTokenDecimalsWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&USDCTokenPoolState",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_token_decimals", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u8",
	})
}

// GetRemotePools encodes a call to the get_remote_pools Move function.
func (c usdcTokenPoolEncoder) GetRemotePools(state bind.Object, remoteChainSelector uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_remote_pools", typeArgsList, typeParamsList, []string{
		"&USDCTokenPoolState",
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
func (c usdcTokenPoolEncoder) GetRemotePoolsWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&USDCTokenPoolState",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_remote_pools", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<vector<u8>>",
	})
}

// IsRemotePool encodes a call to the is_remote_pool Move function.
func (c usdcTokenPoolEncoder) IsRemotePool(state bind.Object, remoteChainSelector uint64, remotePoolAddress []byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("is_remote_pool", typeArgsList, typeParamsList, []string{
		"&USDCTokenPoolState",
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
func (c usdcTokenPoolEncoder) IsRemotePoolWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&USDCTokenPoolState",
		"u64",
		"vector<u8>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("is_remote_pool", typeArgsList, typeParamsList, expectedParams, args, []string{
		"bool",
	})
}

// GetRemoteToken encodes a call to the get_remote_token Move function.
func (c usdcTokenPoolEncoder) GetRemoteToken(state bind.Object, remoteChainSelector uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_remote_token", typeArgsList, typeParamsList, []string{
		"&USDCTokenPoolState",
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
func (c usdcTokenPoolEncoder) GetRemoteTokenWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&USDCTokenPoolState",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_remote_token", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<u8>",
	})
}

// AddRemotePool encodes a call to the add_remote_pool Move function.
func (c usdcTokenPoolEncoder) AddRemotePool(state bind.Object, ownerCap bind.Object, remoteChainSelector uint64, remotePoolAddress []byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("add_remote_pool", typeArgsList, typeParamsList, []string{
		"&mut USDCTokenPoolState",
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
func (c usdcTokenPoolEncoder) AddRemotePoolWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut USDCTokenPoolState",
		"&OwnerCap",
		"u64",
		"vector<u8>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("add_remote_pool", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// RemoveRemotePool encodes a call to the remove_remote_pool Move function.
func (c usdcTokenPoolEncoder) RemoveRemotePool(state bind.Object, ownerCap bind.Object, remoteChainSelector uint64, remotePoolAddress []byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("remove_remote_pool", typeArgsList, typeParamsList, []string{
		"&mut USDCTokenPoolState",
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
func (c usdcTokenPoolEncoder) RemoveRemotePoolWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut USDCTokenPoolState",
		"&OwnerCap",
		"u64",
		"vector<u8>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("remove_remote_pool", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// IsSupportedChain encodes a call to the is_supported_chain Move function.
func (c usdcTokenPoolEncoder) IsSupportedChain(state bind.Object, remoteChainSelector uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("is_supported_chain", typeArgsList, typeParamsList, []string{
		"&USDCTokenPoolState",
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
func (c usdcTokenPoolEncoder) IsSupportedChainWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&USDCTokenPoolState",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("is_supported_chain", typeArgsList, typeParamsList, expectedParams, args, []string{
		"bool",
	})
}

// GetSupportedChains encodes a call to the get_supported_chains Move function.
func (c usdcTokenPoolEncoder) GetSupportedChains(state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_supported_chains", typeArgsList, typeParamsList, []string{
		"&USDCTokenPoolState",
	}, []any{
		state,
	}, []string{
		"vector<u64>",
	})
}

// GetSupportedChainsWithArgs encodes a call to the get_supported_chains Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c usdcTokenPoolEncoder) GetSupportedChainsWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&USDCTokenPoolState",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_supported_chains", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<u64>",
	})
}

// ApplyChainUpdates encodes a call to the apply_chain_updates Move function.
func (c usdcTokenPoolEncoder) ApplyChainUpdates(state bind.Object, ownerCap bind.Object, remoteChainSelectorsToRemove []uint64, remoteChainSelectorsToAdd []uint64, remotePoolAddressesToAdd [][][]byte, remoteTokenAddressesToAdd [][]byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("apply_chain_updates", typeArgsList, typeParamsList, []string{
		"&mut USDCTokenPoolState",
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
func (c usdcTokenPoolEncoder) ApplyChainUpdatesWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut USDCTokenPoolState",
		"&OwnerCap",
		"vector<u64>",
		"vector<u64>",
		"vector<vector<vector<u8>>>",
		"vector<vector<u8>>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("apply_chain_updates", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// GetAllowlistEnabled encodes a call to the get_allowlist_enabled Move function.
func (c usdcTokenPoolEncoder) GetAllowlistEnabled(state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_allowlist_enabled", typeArgsList, typeParamsList, []string{
		"&USDCTokenPoolState",
	}, []any{
		state,
	}, []string{
		"bool",
	})
}

// GetAllowlistEnabledWithArgs encodes a call to the get_allowlist_enabled Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c usdcTokenPoolEncoder) GetAllowlistEnabledWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&USDCTokenPoolState",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_allowlist_enabled", typeArgsList, typeParamsList, expectedParams, args, []string{
		"bool",
	})
}

// GetAllowlist encodes a call to the get_allowlist Move function.
func (c usdcTokenPoolEncoder) GetAllowlist(state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_allowlist", typeArgsList, typeParamsList, []string{
		"&USDCTokenPoolState",
	}, []any{
		state,
	}, []string{
		"vector<address>",
	})
}

// GetAllowlistWithArgs encodes a call to the get_allowlist Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c usdcTokenPoolEncoder) GetAllowlistWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&USDCTokenPoolState",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_allowlist", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<address>",
	})
}

// SetAllowlistEnabled encodes a call to the set_allowlist_enabled Move function.
func (c usdcTokenPoolEncoder) SetAllowlistEnabled(state bind.Object, ownerCap bind.Object, enabled bool) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("set_allowlist_enabled", typeArgsList, typeParamsList, []string{
		"&mut USDCTokenPoolState",
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
func (c usdcTokenPoolEncoder) SetAllowlistEnabledWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut USDCTokenPoolState",
		"&OwnerCap",
		"bool",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("set_allowlist_enabled", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// ApplyAllowlistUpdates encodes a call to the apply_allowlist_updates Move function.
func (c usdcTokenPoolEncoder) ApplyAllowlistUpdates(state bind.Object, ownerCap bind.Object, removes []string, adds []string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("apply_allowlist_updates", typeArgsList, typeParamsList, []string{
		"&mut USDCTokenPoolState",
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
func (c usdcTokenPoolEncoder) ApplyAllowlistUpdatesWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut USDCTokenPoolState",
		"&OwnerCap",
		"vector<address>",
		"vector<address>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
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
func (c usdcTokenPoolEncoder) LockOrBurn(typeArgs []string, ref bind.Object, tokenTransferParams bind.Object, c_ bind.Object, remoteChainSelector uint64, receiver string, clock bind.Object, denyList bind.Object, pool bind.Object, state bind.Object, messageTransmitterState bind.Object, treasury bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("lock_or_burn", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"&mut onramp_sh::TokenTransferParams",
		"Coin<T>",
		"u64",
		"address",
		"&Clock",
		"&DenyList",
		"&mut USDCTokenPoolState",
		"&MinterState",
		"&mut MessageTransmitterState",
		"&mut Treasury<T>",
	}, []any{
		ref,
		tokenTransferParams,
		c_,
		remoteChainSelector,
		receiver,
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
		"address",
		"&Clock",
		"&DenyList",
		"&mut USDCTokenPoolState",
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
func (c usdcTokenPoolEncoder) ReleaseOrMint(typeArgs []string, ref bind.Object, receiverParams bind.Object, tokenTransfer bind.Object, clock bind.Object, denyList bind.Object, pool bind.Object, state bind.Object, messageTransmitterState bind.Object, treasury bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("release_or_mint", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"&mut offramp_sh::ReceiverParams",
		"offramp_sh::DestTokenTransfer",
		"&Clock",
		"&DenyList",
		"&mut USDCTokenPoolState",
		"&mut MinterState",
		"&mut MessageTransmitterState",
		"&mut Treasury<T>",
	}, []any{
		ref,
		receiverParams,
		tokenTransfer,
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
		"offramp_sh::DestTokenTransfer",
		"&Clock",
		"&DenyList",
		"&mut USDCTokenPoolState",
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
func (c usdcTokenPoolEncoder) GetDomain(pool bind.Object, chainSelector uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_domain", typeArgsList, typeParamsList, []string{
		"&USDCTokenPoolState",
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
func (c usdcTokenPoolEncoder) GetDomainWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&USDCTokenPoolState",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_domain", typeArgsList, typeParamsList, expectedParams, args, []string{
		"usdc_token_pool::usdc_token_pool::Domain",
	})
}

// SetDomains encodes a call to the set_domains Move function.
func (c usdcTokenPoolEncoder) SetDomains(pool bind.Object, ownerCap bind.Object, remoteChainSelectors []uint64, remoteDomainIdentifiers []uint32, allowedRemoteCallers [][]byte, enableds []bool) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("set_domains", typeArgsList, typeParamsList, []string{
		"&mut USDCTokenPoolState",
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
func (c usdcTokenPoolEncoder) SetDomainsWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut USDCTokenPoolState",
		"&OwnerCap",
		"vector<u64>",
		"vector<u32>",
		"vector<vector<u8>>",
		"vector<bool>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("set_domains", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// SetChainRateLimiterConfigs encodes a call to the set_chain_rate_limiter_configs Move function.
func (c usdcTokenPoolEncoder) SetChainRateLimiterConfigs(state bind.Object, ownerCap bind.Object, clock bind.Object, remoteChainSelectors []uint64, outboundIsEnableds []bool, outboundCapacities []uint64, outboundRates []uint64, inboundIsEnableds []bool, inboundCapacities []uint64, inboundRates []uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("set_chain_rate_limiter_configs", typeArgsList, typeParamsList, []string{
		"&mut USDCTokenPoolState",
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
func (c usdcTokenPoolEncoder) SetChainRateLimiterConfigsWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut USDCTokenPoolState",
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
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("set_chain_rate_limiter_configs", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// SetChainRateLimiterConfig encodes a call to the set_chain_rate_limiter_config Move function.
func (c usdcTokenPoolEncoder) SetChainRateLimiterConfig(state bind.Object, ownerCap bind.Object, clock bind.Object, remoteChainSelector uint64, outboundIsEnabled bool, outboundCapacity uint64, outboundRate uint64, inboundIsEnabled bool, inboundCapacity uint64, inboundRate uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("set_chain_rate_limiter_config", typeArgsList, typeParamsList, []string{
		"&mut USDCTokenPoolState",
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
func (c usdcTokenPoolEncoder) SetChainRateLimiterConfigWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut USDCTokenPoolState",
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
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("set_chain_rate_limiter_config", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// Owner encodes a call to the owner Move function.
func (c usdcTokenPoolEncoder) Owner(state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("owner", typeArgsList, typeParamsList, []string{
		"&USDCTokenPoolState",
	}, []any{
		state,
	}, []string{
		"address",
	})
}

// OwnerWithArgs encodes a call to the owner Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c usdcTokenPoolEncoder) OwnerWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&USDCTokenPoolState",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("owner", typeArgsList, typeParamsList, expectedParams, args, []string{
		"address",
	})
}

// HasPendingTransfer encodes a call to the has_pending_transfer Move function.
func (c usdcTokenPoolEncoder) HasPendingTransfer(state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("has_pending_transfer", typeArgsList, typeParamsList, []string{
		"&USDCTokenPoolState",
	}, []any{
		state,
	}, []string{
		"bool",
	})
}

// HasPendingTransferWithArgs encodes a call to the has_pending_transfer Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c usdcTokenPoolEncoder) HasPendingTransferWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&USDCTokenPoolState",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("has_pending_transfer", typeArgsList, typeParamsList, expectedParams, args, []string{
		"bool",
	})
}

// PendingTransferFrom encodes a call to the pending_transfer_from Move function.
func (c usdcTokenPoolEncoder) PendingTransferFrom(state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("pending_transfer_from", typeArgsList, typeParamsList, []string{
		"&USDCTokenPoolState",
	}, []any{
		state,
	}, []string{
		"0x1::option::Option<address>",
	})
}

// PendingTransferFromWithArgs encodes a call to the pending_transfer_from Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c usdcTokenPoolEncoder) PendingTransferFromWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&USDCTokenPoolState",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("pending_transfer_from", typeArgsList, typeParamsList, expectedParams, args, []string{
		"0x1::option::Option<address>",
	})
}

// PendingTransferTo encodes a call to the pending_transfer_to Move function.
func (c usdcTokenPoolEncoder) PendingTransferTo(state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("pending_transfer_to", typeArgsList, typeParamsList, []string{
		"&USDCTokenPoolState",
	}, []any{
		state,
	}, []string{
		"0x1::option::Option<address>",
	})
}

// PendingTransferToWithArgs encodes a call to the pending_transfer_to Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c usdcTokenPoolEncoder) PendingTransferToWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&USDCTokenPoolState",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("pending_transfer_to", typeArgsList, typeParamsList, expectedParams, args, []string{
		"0x1::option::Option<address>",
	})
}

// PendingTransferAccepted encodes a call to the pending_transfer_accepted Move function.
func (c usdcTokenPoolEncoder) PendingTransferAccepted(state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("pending_transfer_accepted", typeArgsList, typeParamsList, []string{
		"&USDCTokenPoolState",
	}, []any{
		state,
	}, []string{
		"0x1::option::Option<bool>",
	})
}

// PendingTransferAcceptedWithArgs encodes a call to the pending_transfer_accepted Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c usdcTokenPoolEncoder) PendingTransferAcceptedWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&USDCTokenPoolState",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("pending_transfer_accepted", typeArgsList, typeParamsList, expectedParams, args, []string{
		"0x1::option::Option<bool>",
	})
}

// TransferOwnership encodes a call to the transfer_ownership Move function.
func (c usdcTokenPoolEncoder) TransferOwnership(state bind.Object, ownerCap bind.Object, newOwner string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("transfer_ownership", typeArgsList, typeParamsList, []string{
		"&mut USDCTokenPoolState",
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
func (c usdcTokenPoolEncoder) TransferOwnershipWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut USDCTokenPoolState",
		"&OwnerCap",
		"address",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("transfer_ownership", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// AcceptOwnership encodes a call to the accept_ownership Move function.
func (c usdcTokenPoolEncoder) AcceptOwnership(state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("accept_ownership", typeArgsList, typeParamsList, []string{
		"&mut USDCTokenPoolState",
	}, []any{
		state,
	}, nil)
}

// AcceptOwnershipWithArgs encodes a call to the accept_ownership Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c usdcTokenPoolEncoder) AcceptOwnershipWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut USDCTokenPoolState",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("accept_ownership", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// AcceptOwnershipFromObject encodes a call to the accept_ownership_from_object Move function.
func (c usdcTokenPoolEncoder) AcceptOwnershipFromObject(state bind.Object, from string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("accept_ownership_from_object", typeArgsList, typeParamsList, []string{
		"&mut USDCTokenPoolState",
		"&mut UID",
	}, []any{
		state,
		from,
	}, nil)
}

// AcceptOwnershipFromObjectWithArgs encodes a call to the accept_ownership_from_object Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c usdcTokenPoolEncoder) AcceptOwnershipFromObjectWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut USDCTokenPoolState",
		"&mut UID",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("accept_ownership_from_object", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// ExecuteOwnershipTransfer encodes a call to the execute_ownership_transfer Move function.
func (c usdcTokenPoolEncoder) ExecuteOwnershipTransfer(ownerCap bind.Object, ownableState bind.Object, to string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("execute_ownership_transfer", typeArgsList, typeParamsList, []string{
		"OwnerCap",
		"&mut OwnableState",
		"address",
	}, []any{
		ownerCap,
		ownableState,
		to,
	}, nil)
}

// ExecuteOwnershipTransferWithArgs encodes a call to the execute_ownership_transfer Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c usdcTokenPoolEncoder) ExecuteOwnershipTransferWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"OwnerCap",
		"&mut OwnableState",
		"address",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("execute_ownership_transfer", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// ExecuteOwnershipTransferToMcms encodes a call to the execute_ownership_transfer_to_mcms Move function.
func (c usdcTokenPoolEncoder) ExecuteOwnershipTransferToMcms(ownerCap bind.Object, state bind.Object, registry bind.Object, to string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("execute_ownership_transfer_to_mcms", typeArgsList, typeParamsList, []string{
		"OwnerCap",
		"&mut USDCTokenPoolState",
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
func (c usdcTokenPoolEncoder) ExecuteOwnershipTransferToMcmsWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"OwnerCap",
		"&mut USDCTokenPoolState",
		"&mut Registry",
		"address",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
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

// McmsEntrypoint encodes a call to the mcms_entrypoint Move function.
func (c usdcTokenPoolEncoder) McmsEntrypoint(state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_entrypoint", typeArgsList, typeParamsList, []string{
		"&mut USDCTokenPoolState",
		"&mut Registry",
		"ExecutingCallbackParams",
	}, []any{
		state,
		registry,
		params,
	}, nil)
}

// McmsEntrypointWithArgs encodes a call to the mcms_entrypoint Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c usdcTokenPoolEncoder) McmsEntrypointWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut USDCTokenPoolState",
		"&mut Registry",
		"ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_entrypoint", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// DestroyTokenPool encodes a call to the destroy_token_pool Move function.
func (c usdcTokenPoolEncoder) DestroyTokenPool(state bind.Object, ownerCap bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("destroy_token_pool", typeArgsList, typeParamsList, []string{
		"usdc_token_pool::usdc_token_pool::USDCTokenPoolState",
		"OwnerCap",
	}, []any{
		state,
		ownerCap,
	}, nil)
}

// DestroyTokenPoolWithArgs encodes a call to the destroy_token_pool Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c usdcTokenPoolEncoder) DestroyTokenPoolWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"usdc_token_pool::usdc_token_pool::USDCTokenPoolState",
		"OwnerCap",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("destroy_token_pool", typeArgsList, typeParamsList, expectedParams, args, nil)
}
