// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package module_burn_mint_token_pool

import (
	"context"
	"fmt"
	"math/big"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/mystenbcs"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/block-vision/sui-go-sdk/transaction"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
)

var (
	_ = big.NewInt
)

type IBurnMintTokenPool interface {
	TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	Initialize(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, coinMetadata bind.Object, treasuryCap bind.Object, burnMintTokenPoolPackageId string, tokenPoolAdministrator string, lockOrBurnParams []string, releaseOrMintParams []string) (*models.SuiTransactionBlockResponse, error)
	InitializeByCcipAdmin(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, ownerCap bind.Object, coinMetadata bind.Object, treasuryCap bind.Object, burnMintTokenPoolPackageId string, tokenPoolAdministrator string, lockOrBurnParams []string, releaseOrMintParams []string) (*models.SuiTransactionBlockResponse, error)
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
	LockOrBurn(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, c_ bind.Object, tokenParams bind.Object, state bind.Object, clock bind.Object) (*models.SuiTransactionBlockResponse, error)
	ReleaseOrMint(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, receiverParams bind.Object, index uint64, pool bind.Object, clock bind.Object) (*models.SuiTransactionBlockResponse, error)
	SetChainRateLimiterConfigs(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object, clock bind.Object, remoteChainSelectors []uint64, outboundIsEnableds []bool, outboundCapacities []uint64, outboundRates []uint64, inboundIsEnableds []bool, inboundCapacities []uint64, inboundRates []uint64) (*models.SuiTransactionBlockResponse, error)
	SetChainRateLimiterConfig(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object, clock bind.Object, remoteChainSelector uint64, outboundIsEnabled bool, outboundCapacity uint64, outboundRate uint64, inboundIsEnabled bool, inboundCapacity uint64, inboundRate uint64) (*models.SuiTransactionBlockResponse, error)
	DestroyTokenPool(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object) (*models.SuiTransactionBlockResponse, error)
	Owner(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	HasPendingTransfer(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	PendingTransferFrom(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	PendingTransferTo(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	PendingTransferAccepted(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	TransferOwnership(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object, newOwner string) (*models.SuiTransactionBlockResponse, error)
	AcceptOwnership(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	AcceptOwnershipFromObject(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, from string) (*models.SuiTransactionBlockResponse, error)
	ExecuteOwnershipTransfer(ctx context.Context, opts *bind.CallOpts, ownerCap bind.Object, ownableState bind.Object, to string) (*models.SuiTransactionBlockResponse, error)
	McmsRegisterEntrypoint(ctx context.Context, opts *bind.CallOpts, typeArgs []string, registry bind.Object, state bind.Object, ownerCap bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsRegisterUpgradeCap(ctx context.Context, opts *bind.CallOpts, upgradeCap bind.Object, registry bind.Object, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsEntrypoint(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	DevInspect() IBurnMintTokenPoolDevInspect
	Encoder() BurnMintTokenPoolEncoder
}

type IBurnMintTokenPoolDevInspect interface {
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
	ReleaseOrMint(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, receiverParams bind.Object, index uint64, pool bind.Object, clock bind.Object) (bind.Object, error)
	DestroyTokenPool(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object) (any, error)
	Owner(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (string, error)
	HasPendingTransfer(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (bool, error)
	PendingTransferFrom(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*string, error)
	PendingTransferTo(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*string, error)
	PendingTransferAccepted(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*bool, error)
}

type BurnMintTokenPoolEncoder interface {
	TypeAndVersion() (*bind.EncodedCall, error)
	TypeAndVersionWithArgs(args ...any) (*bind.EncodedCall, error)
	Initialize(typeArgs []string, ref bind.Object, coinMetadata bind.Object, treasuryCap bind.Object, burnMintTokenPoolPackageId string, tokenPoolAdministrator string, lockOrBurnParams []string, releaseOrMintParams []string) (*bind.EncodedCall, error)
	InitializeWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	InitializeByCcipAdmin(typeArgs []string, ref bind.Object, ownerCap bind.Object, coinMetadata bind.Object, treasuryCap bind.Object, burnMintTokenPoolPackageId string, tokenPoolAdministrator string, lockOrBurnParams []string, releaseOrMintParams []string) (*bind.EncodedCall, error)
	InitializeByCcipAdminWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
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
	LockOrBurn(typeArgs []string, ref bind.Object, c_ bind.Object, tokenParams bind.Object, state bind.Object, clock bind.Object) (*bind.EncodedCall, error)
	LockOrBurnWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	ReleaseOrMint(typeArgs []string, ref bind.Object, receiverParams bind.Object, index uint64, pool bind.Object, clock bind.Object) (*bind.EncodedCall, error)
	ReleaseOrMintWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	SetChainRateLimiterConfigs(typeArgs []string, state bind.Object, ownerCap bind.Object, clock bind.Object, remoteChainSelectors []uint64, outboundIsEnableds []bool, outboundCapacities []uint64, outboundRates []uint64, inboundIsEnableds []bool, inboundCapacities []uint64, inboundRates []uint64) (*bind.EncodedCall, error)
	SetChainRateLimiterConfigsWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	SetChainRateLimiterConfig(typeArgs []string, state bind.Object, ownerCap bind.Object, clock bind.Object, remoteChainSelector uint64, outboundIsEnabled bool, outboundCapacity uint64, outboundRate uint64, inboundIsEnabled bool, inboundCapacity uint64, inboundRate uint64) (*bind.EncodedCall, error)
	SetChainRateLimiterConfigWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	DestroyTokenPool(typeArgs []string, state bind.Object, ownerCap bind.Object) (*bind.EncodedCall, error)
	DestroyTokenPoolWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
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
	ExecuteOwnershipTransfer(ownerCap bind.Object, ownableState bind.Object, to string) (*bind.EncodedCall, error)
	ExecuteOwnershipTransferWithArgs(args ...any) (*bind.EncodedCall, error)
	McmsRegisterEntrypoint(typeArgs []string, registry bind.Object, state bind.Object, ownerCap bind.Object) (*bind.EncodedCall, error)
	McmsRegisterEntrypointWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	McmsRegisterUpgradeCap(upgradeCap bind.Object, registry bind.Object, state bind.Object) (*bind.EncodedCall, error)
	McmsRegisterUpgradeCapWithArgs(args ...any) (*bind.EncodedCall, error)
	McmsEntrypoint(typeArgs []string, state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsEntrypointWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
}

type BurnMintTokenPoolContract struct {
	*bind.BoundContract
	burnMintTokenPoolEncoder
	devInspect *BurnMintTokenPoolDevInspect
}

type BurnMintTokenPoolDevInspect struct {
	contract *BurnMintTokenPoolContract
}

var _ IBurnMintTokenPool = (*BurnMintTokenPoolContract)(nil)
var _ IBurnMintTokenPoolDevInspect = (*BurnMintTokenPoolDevInspect)(nil)

func NewBurnMintTokenPool(packageID string, client sui.ISuiAPI) (*BurnMintTokenPoolContract, error) {
	contract, err := bind.NewBoundContract(packageID, "burn_mint_token_pool", "burn_mint_token_pool", client)
	if err != nil {
		return nil, err
	}

	c := &BurnMintTokenPoolContract{
		BoundContract:            contract,
		burnMintTokenPoolEncoder: burnMintTokenPoolEncoder{BoundContract: contract},
	}
	c.devInspect = &BurnMintTokenPoolDevInspect{contract: c}
	return c, nil
}

func (c *BurnMintTokenPoolContract) Encoder() BurnMintTokenPoolEncoder {
	return c.burnMintTokenPoolEncoder
}

func (c *BurnMintTokenPoolContract) DevInspect() IBurnMintTokenPoolDevInspect {
	return c.devInspect
}

func (c *BurnMintTokenPoolContract) BuildPTB(ctx context.Context, ptb *transaction.Transaction, encoded *bind.EncodedCall) (*transaction.Argument, error) {
	var callArgManager *bind.CallArgManager
	if ptb.Data.V1 != nil && ptb.Data.V1.Kind.ProgrammableTransaction != nil &&
		ptb.Data.V1.Kind.ProgrammableTransaction.Inputs != nil {
		callArgManager = bind.NewCallArgManagerWithExisting(ptb.Data.V1.Kind.ProgrammableTransaction.Inputs)
	} else {
		callArgManager = bind.NewCallArgManager()
	}

	arguments, err := callArgManager.ConvertEncodedCallArgsToArguments(encoded.CallArgs)
	if err != nil {
		return nil, fmt.Errorf("failed to convert EncodedCallArguments to Arguments: %w", err)
	}

	ptb.Data.V1.Kind.ProgrammableTransaction.Inputs = callArgManager.GetInputs()

	typeTagValues := make([]transaction.TypeTag, len(encoded.TypeArgs))
	for i, tag := range encoded.TypeArgs {
		if tag != nil {
			typeTagValues[i] = *tag
		}
	}

	argumentValues := make([]transaction.Argument, len(arguments))
	for i, arg := range arguments {
		if arg != nil {
			argumentValues[i] = *arg
		}
	}

	result := ptb.MoveCall(
		models.SuiAddress(encoded.Module.PackageID),
		encoded.Module.ModuleName,
		encoded.Function,
		typeTagValues,
		argumentValues,
	)

	return &result, nil
}

type BurnMintTokenPoolState struct {
	Id             string      `move:"sui::object::UID"`
	TokenPoolState bind.Object `move:"TokenPoolState"`
	TreasuryCap    bind.Object `move:"TreasuryCap<T>"`
	OwnableState   bind.Object `move:"OwnableState"`
}

type TypeProof struct {
}

type McmsCallback struct {
}

func init() {
	bind.RegisterStructDecoder("burn_mint_token_pool::burn_mint_token_pool::BurnMintTokenPoolState", func(data []byte) (interface{}, error) {
		var result BurnMintTokenPoolState
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("burn_mint_token_pool::burn_mint_token_pool::TypeProof", func(data []byte) (interface{}, error) {
		var result TypeProof
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("burn_mint_token_pool::burn_mint_token_pool::McmsCallback", func(data []byte) (interface{}, error) {
		var result McmsCallback
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
}

// TypeAndVersion executes the type_and_version Move function.
func (c *BurnMintTokenPoolContract) TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.burnMintTokenPoolEncoder.TypeAndVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// Initialize executes the initialize Move function.
func (c *BurnMintTokenPoolContract) Initialize(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, coinMetadata bind.Object, treasuryCap bind.Object, burnMintTokenPoolPackageId string, tokenPoolAdministrator string, lockOrBurnParams []string, releaseOrMintParams []string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.burnMintTokenPoolEncoder.Initialize(typeArgs, ref, coinMetadata, treasuryCap, burnMintTokenPoolPackageId, tokenPoolAdministrator, lockOrBurnParams, releaseOrMintParams)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// InitializeByCcipAdmin executes the initialize_by_ccip_admin Move function.
func (c *BurnMintTokenPoolContract) InitializeByCcipAdmin(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, ownerCap bind.Object, coinMetadata bind.Object, treasuryCap bind.Object, burnMintTokenPoolPackageId string, tokenPoolAdministrator string, lockOrBurnParams []string, releaseOrMintParams []string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.burnMintTokenPoolEncoder.InitializeByCcipAdmin(typeArgs, ref, ownerCap, coinMetadata, treasuryCap, burnMintTokenPoolPackageId, tokenPoolAdministrator, lockOrBurnParams, releaseOrMintParams)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetToken executes the get_token Move function.
func (c *BurnMintTokenPoolContract) GetToken(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.burnMintTokenPoolEncoder.GetToken(typeArgs, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetTokenDecimals executes the get_token_decimals Move function.
func (c *BurnMintTokenPoolContract) GetTokenDecimals(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.burnMintTokenPoolEncoder.GetTokenDecimals(typeArgs, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetRemotePools executes the get_remote_pools Move function.
func (c *BurnMintTokenPoolContract) GetRemotePools(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.burnMintTokenPoolEncoder.GetRemotePools(typeArgs, state, remoteChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// IsRemotePool executes the is_remote_pool Move function.
func (c *BurnMintTokenPoolContract) IsRemotePool(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, remoteChainSelector uint64, remotePoolAddress []byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.burnMintTokenPoolEncoder.IsRemotePool(typeArgs, state, remoteChainSelector, remotePoolAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetRemoteToken executes the get_remote_token Move function.
func (c *BurnMintTokenPoolContract) GetRemoteToken(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.burnMintTokenPoolEncoder.GetRemoteToken(typeArgs, state, remoteChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// AddRemotePool executes the add_remote_pool Move function.
func (c *BurnMintTokenPoolContract) AddRemotePool(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object, remoteChainSelector uint64, remotePoolAddress []byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.burnMintTokenPoolEncoder.AddRemotePool(typeArgs, state, ownerCap, remoteChainSelector, remotePoolAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// RemoveRemotePool executes the remove_remote_pool Move function.
func (c *BurnMintTokenPoolContract) RemoveRemotePool(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object, remoteChainSelector uint64, remotePoolAddress []byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.burnMintTokenPoolEncoder.RemoveRemotePool(typeArgs, state, ownerCap, remoteChainSelector, remotePoolAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// IsSupportedChain executes the is_supported_chain Move function.
func (c *BurnMintTokenPoolContract) IsSupportedChain(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.burnMintTokenPoolEncoder.IsSupportedChain(typeArgs, state, remoteChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetSupportedChains executes the get_supported_chains Move function.
func (c *BurnMintTokenPoolContract) GetSupportedChains(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.burnMintTokenPoolEncoder.GetSupportedChains(typeArgs, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ApplyChainUpdates executes the apply_chain_updates Move function.
func (c *BurnMintTokenPoolContract) ApplyChainUpdates(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object, remoteChainSelectorsToRemove []uint64, remoteChainSelectorsToAdd []uint64, remotePoolAddressesToAdd [][][]byte, remoteTokenAddressesToAdd [][]byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.burnMintTokenPoolEncoder.ApplyChainUpdates(typeArgs, state, ownerCap, remoteChainSelectorsToRemove, remoteChainSelectorsToAdd, remotePoolAddressesToAdd, remoteTokenAddressesToAdd)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetAllowlistEnabled executes the get_allowlist_enabled Move function.
func (c *BurnMintTokenPoolContract) GetAllowlistEnabled(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.burnMintTokenPoolEncoder.GetAllowlistEnabled(typeArgs, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetAllowlist executes the get_allowlist Move function.
func (c *BurnMintTokenPoolContract) GetAllowlist(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.burnMintTokenPoolEncoder.GetAllowlist(typeArgs, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// SetAllowlistEnabled executes the set_allowlist_enabled Move function.
func (c *BurnMintTokenPoolContract) SetAllowlistEnabled(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object, enabled bool) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.burnMintTokenPoolEncoder.SetAllowlistEnabled(typeArgs, state, ownerCap, enabled)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ApplyAllowlistUpdates executes the apply_allowlist_updates Move function.
func (c *BurnMintTokenPoolContract) ApplyAllowlistUpdates(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object, removes []string, adds []string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.burnMintTokenPoolEncoder.ApplyAllowlistUpdates(typeArgs, state, ownerCap, removes, adds)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// LockOrBurn executes the lock_or_burn Move function.
func (c *BurnMintTokenPoolContract) LockOrBurn(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, c_ bind.Object, tokenParams bind.Object, state bind.Object, clock bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.burnMintTokenPoolEncoder.LockOrBurn(typeArgs, ref, c_, tokenParams, state, clock)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ReleaseOrMint executes the release_or_mint Move function.
func (c *BurnMintTokenPoolContract) ReleaseOrMint(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, receiverParams bind.Object, index uint64, pool bind.Object, clock bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.burnMintTokenPoolEncoder.ReleaseOrMint(typeArgs, ref, receiverParams, index, pool, clock)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// SetChainRateLimiterConfigs executes the set_chain_rate_limiter_configs Move function.
func (c *BurnMintTokenPoolContract) SetChainRateLimiterConfigs(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object, clock bind.Object, remoteChainSelectors []uint64, outboundIsEnableds []bool, outboundCapacities []uint64, outboundRates []uint64, inboundIsEnableds []bool, inboundCapacities []uint64, inboundRates []uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.burnMintTokenPoolEncoder.SetChainRateLimiterConfigs(typeArgs, state, ownerCap, clock, remoteChainSelectors, outboundIsEnableds, outboundCapacities, outboundRates, inboundIsEnableds, inboundCapacities, inboundRates)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// SetChainRateLimiterConfig executes the set_chain_rate_limiter_config Move function.
func (c *BurnMintTokenPoolContract) SetChainRateLimiterConfig(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object, clock bind.Object, remoteChainSelector uint64, outboundIsEnabled bool, outboundCapacity uint64, outboundRate uint64, inboundIsEnabled bool, inboundCapacity uint64, inboundRate uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.burnMintTokenPoolEncoder.SetChainRateLimiterConfig(typeArgs, state, ownerCap, clock, remoteChainSelector, outboundIsEnabled, outboundCapacity, outboundRate, inboundIsEnabled, inboundCapacity, inboundRate)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// DestroyTokenPool executes the destroy_token_pool Move function.
func (c *BurnMintTokenPoolContract) DestroyTokenPool(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.burnMintTokenPoolEncoder.DestroyTokenPool(typeArgs, state, ownerCap)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// Owner executes the owner Move function.
func (c *BurnMintTokenPoolContract) Owner(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.burnMintTokenPoolEncoder.Owner(typeArgs, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// HasPendingTransfer executes the has_pending_transfer Move function.
func (c *BurnMintTokenPoolContract) HasPendingTransfer(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.burnMintTokenPoolEncoder.HasPendingTransfer(typeArgs, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// PendingTransferFrom executes the pending_transfer_from Move function.
func (c *BurnMintTokenPoolContract) PendingTransferFrom(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.burnMintTokenPoolEncoder.PendingTransferFrom(typeArgs, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// PendingTransferTo executes the pending_transfer_to Move function.
func (c *BurnMintTokenPoolContract) PendingTransferTo(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.burnMintTokenPoolEncoder.PendingTransferTo(typeArgs, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// PendingTransferAccepted executes the pending_transfer_accepted Move function.
func (c *BurnMintTokenPoolContract) PendingTransferAccepted(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.burnMintTokenPoolEncoder.PendingTransferAccepted(typeArgs, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// TransferOwnership executes the transfer_ownership Move function.
func (c *BurnMintTokenPoolContract) TransferOwnership(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object, newOwner string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.burnMintTokenPoolEncoder.TransferOwnership(typeArgs, state, ownerCap, newOwner)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// AcceptOwnership executes the accept_ownership Move function.
func (c *BurnMintTokenPoolContract) AcceptOwnership(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.burnMintTokenPoolEncoder.AcceptOwnership(typeArgs, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// AcceptOwnershipFromObject executes the accept_ownership_from_object Move function.
func (c *BurnMintTokenPoolContract) AcceptOwnershipFromObject(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, from string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.burnMintTokenPoolEncoder.AcceptOwnershipFromObject(typeArgs, state, from)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ExecuteOwnershipTransfer executes the execute_ownership_transfer Move function.
func (c *BurnMintTokenPoolContract) ExecuteOwnershipTransfer(ctx context.Context, opts *bind.CallOpts, ownerCap bind.Object, ownableState bind.Object, to string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.burnMintTokenPoolEncoder.ExecuteOwnershipTransfer(ownerCap, ownableState, to)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsRegisterEntrypoint executes the mcms_register_entrypoint Move function.
func (c *BurnMintTokenPoolContract) McmsRegisterEntrypoint(ctx context.Context, opts *bind.CallOpts, typeArgs []string, registry bind.Object, state bind.Object, ownerCap bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.burnMintTokenPoolEncoder.McmsRegisterEntrypoint(typeArgs, registry, state, ownerCap)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsRegisterUpgradeCap executes the mcms_register_upgrade_cap Move function.
func (c *BurnMintTokenPoolContract) McmsRegisterUpgradeCap(ctx context.Context, opts *bind.CallOpts, upgradeCap bind.Object, registry bind.Object, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.burnMintTokenPoolEncoder.McmsRegisterUpgradeCap(upgradeCap, registry, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsEntrypoint executes the mcms_entrypoint Move function.
func (c *BurnMintTokenPoolContract) McmsEntrypoint(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.burnMintTokenPoolEncoder.McmsEntrypoint(typeArgs, state, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// TypeAndVersion executes the type_and_version Move function using DevInspect to get return values.
//
// Returns: 0x1::string::String
func (d *BurnMintTokenPoolDevInspect) TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (string, error) {
	encoded, err := d.contract.burnMintTokenPoolEncoder.TypeAndVersion()
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
func (d *BurnMintTokenPoolDevInspect) GetToken(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (string, error) {
	encoded, err := d.contract.burnMintTokenPoolEncoder.GetToken(typeArgs, state)
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
func (d *BurnMintTokenPoolDevInspect) GetTokenDecimals(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (byte, error) {
	encoded, err := d.contract.burnMintTokenPoolEncoder.GetTokenDecimals(typeArgs, state)
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
func (d *BurnMintTokenPoolDevInspect) GetRemotePools(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, remoteChainSelector uint64) ([][]byte, error) {
	encoded, err := d.contract.burnMintTokenPoolEncoder.GetRemotePools(typeArgs, state, remoteChainSelector)
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
func (d *BurnMintTokenPoolDevInspect) IsRemotePool(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, remoteChainSelector uint64, remotePoolAddress []byte) (bool, error) {
	encoded, err := d.contract.burnMintTokenPoolEncoder.IsRemotePool(typeArgs, state, remoteChainSelector, remotePoolAddress)
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
func (d *BurnMintTokenPoolDevInspect) GetRemoteToken(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, remoteChainSelector uint64) ([]byte, error) {
	encoded, err := d.contract.burnMintTokenPoolEncoder.GetRemoteToken(typeArgs, state, remoteChainSelector)
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
func (d *BurnMintTokenPoolDevInspect) IsSupportedChain(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, remoteChainSelector uint64) (bool, error) {
	encoded, err := d.contract.burnMintTokenPoolEncoder.IsSupportedChain(typeArgs, state, remoteChainSelector)
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
func (d *BurnMintTokenPoolDevInspect) GetSupportedChains(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) ([]uint64, error) {
	encoded, err := d.contract.burnMintTokenPoolEncoder.GetSupportedChains(typeArgs, state)
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
func (d *BurnMintTokenPoolDevInspect) GetAllowlistEnabled(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (bool, error) {
	encoded, err := d.contract.burnMintTokenPoolEncoder.GetAllowlistEnabled(typeArgs, state)
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
func (d *BurnMintTokenPoolDevInspect) GetAllowlist(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) ([]string, error) {
	encoded, err := d.contract.burnMintTokenPoolEncoder.GetAllowlist(typeArgs, state)
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

// ReleaseOrMint executes the release_or_mint Move function using DevInspect to get return values.
//
// Returns: osh::ReceiverParams
func (d *BurnMintTokenPoolDevInspect) ReleaseOrMint(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, receiverParams bind.Object, index uint64, pool bind.Object, clock bind.Object) (bind.Object, error) {
	encoded, err := d.contract.burnMintTokenPoolEncoder.ReleaseOrMint(typeArgs, ref, receiverParams, index, pool, clock)
	if err != nil {
		return bind.Object{}, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return bind.Object{}, err
	}
	if len(results) == 0 {
		return bind.Object{}, fmt.Errorf("no return value")
	}
	result, ok := results[0].(bind.Object)
	if !ok {
		return bind.Object{}, fmt.Errorf("unexpected return type: expected bind.Object, got %T", results[0])
	}
	return result, nil
}

// DestroyTokenPool executes the destroy_token_pool Move function using DevInspect to get return values.
//
// Returns: TreasuryCap<T>
func (d *BurnMintTokenPoolDevInspect) DestroyTokenPool(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object) (any, error) {
	encoded, err := d.contract.burnMintTokenPoolEncoder.DestroyTokenPool(typeArgs, state, ownerCap)
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

// Owner executes the owner Move function using DevInspect to get return values.
//
// Returns: address
func (d *BurnMintTokenPoolDevInspect) Owner(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (string, error) {
	encoded, err := d.contract.burnMintTokenPoolEncoder.Owner(typeArgs, state)
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
func (d *BurnMintTokenPoolDevInspect) HasPendingTransfer(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (bool, error) {
	encoded, err := d.contract.burnMintTokenPoolEncoder.HasPendingTransfer(typeArgs, state)
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
func (d *BurnMintTokenPoolDevInspect) PendingTransferFrom(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*string, error) {
	encoded, err := d.contract.burnMintTokenPoolEncoder.PendingTransferFrom(typeArgs, state)
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
func (d *BurnMintTokenPoolDevInspect) PendingTransferTo(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*string, error) {
	encoded, err := d.contract.burnMintTokenPoolEncoder.PendingTransferTo(typeArgs, state)
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
func (d *BurnMintTokenPoolDevInspect) PendingTransferAccepted(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*bool, error) {
	encoded, err := d.contract.burnMintTokenPoolEncoder.PendingTransferAccepted(typeArgs, state)
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

type burnMintTokenPoolEncoder struct {
	*bind.BoundContract
}

// TypeAndVersion encodes a call to the type_and_version Move function.
func (c burnMintTokenPoolEncoder) TypeAndVersion() (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("type_and_version", typeArgsList, typeParamsList, []string{}, []any{}, []string{
		"0x1::string::String",
	})
}

// TypeAndVersionWithArgs encodes a call to the type_and_version Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c burnMintTokenPoolEncoder) TypeAndVersionWithArgs(args ...any) (*bind.EncodedCall, error) {
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
func (c burnMintTokenPoolEncoder) Initialize(typeArgs []string, ref bind.Object, coinMetadata bind.Object, treasuryCap bind.Object, burnMintTokenPoolPackageId string, tokenPoolAdministrator string, lockOrBurnParams []string, releaseOrMintParams []string) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("initialize", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"&CoinMetadata<T>",
		"TreasuryCap<T>",
		"address",
		"address",
		"vector<address>",
		"vector<address>",
	}, []any{
		ref,
		coinMetadata,
		treasuryCap,
		burnMintTokenPoolPackageId,
		tokenPoolAdministrator,
		lockOrBurnParams,
		releaseOrMintParams,
	}, nil)
}

// InitializeWithArgs encodes a call to the initialize Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c burnMintTokenPoolEncoder) InitializeWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"&CoinMetadata<T>",
		"TreasuryCap<T>",
		"address",
		"address",
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
	return c.EncodeCallArgsWithGenerics("initialize", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// InitializeByCcipAdmin encodes a call to the initialize_by_ccip_admin Move function.
func (c burnMintTokenPoolEncoder) InitializeByCcipAdmin(typeArgs []string, ref bind.Object, ownerCap bind.Object, coinMetadata bind.Object, treasuryCap bind.Object, burnMintTokenPoolPackageId string, tokenPoolAdministrator string, lockOrBurnParams []string, releaseOrMintParams []string) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("initialize_by_ccip_admin", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"&state_object::OwnerCap",
		"&CoinMetadata<T>",
		"TreasuryCap<T>",
		"address",
		"address",
		"vector<address>",
		"vector<address>",
	}, []any{
		ref,
		ownerCap,
		coinMetadata,
		treasuryCap,
		burnMintTokenPoolPackageId,
		tokenPoolAdministrator,
		lockOrBurnParams,
		releaseOrMintParams,
	}, nil)
}

// InitializeByCcipAdminWithArgs encodes a call to the initialize_by_ccip_admin Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c burnMintTokenPoolEncoder) InitializeByCcipAdminWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"&state_object::OwnerCap",
		"&CoinMetadata<T>",
		"TreasuryCap<T>",
		"address",
		"address",
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
	return c.EncodeCallArgsWithGenerics("initialize_by_ccip_admin", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// GetToken encodes a call to the get_token Move function.
func (c burnMintTokenPoolEncoder) GetToken(typeArgs []string, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("get_token", typeArgsList, typeParamsList, []string{
		"&BurnMintTokenPoolState<T>",
	}, []any{
		state,
	}, []string{
		"address",
	})
}

// GetTokenWithArgs encodes a call to the get_token Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c burnMintTokenPoolEncoder) GetTokenWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&BurnMintTokenPoolState<T>",
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
func (c burnMintTokenPoolEncoder) GetTokenDecimals(typeArgs []string, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("get_token_decimals", typeArgsList, typeParamsList, []string{
		"&BurnMintTokenPoolState<T>",
	}, []any{
		state,
	}, []string{
		"u8",
	})
}

// GetTokenDecimalsWithArgs encodes a call to the get_token_decimals Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c burnMintTokenPoolEncoder) GetTokenDecimalsWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&BurnMintTokenPoolState<T>",
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
func (c burnMintTokenPoolEncoder) GetRemotePools(typeArgs []string, state bind.Object, remoteChainSelector uint64) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("get_remote_pools", typeArgsList, typeParamsList, []string{
		"&BurnMintTokenPoolState<T>",
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
func (c burnMintTokenPoolEncoder) GetRemotePoolsWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&BurnMintTokenPoolState<T>",
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
func (c burnMintTokenPoolEncoder) IsRemotePool(typeArgs []string, state bind.Object, remoteChainSelector uint64, remotePoolAddress []byte) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("is_remote_pool", typeArgsList, typeParamsList, []string{
		"&BurnMintTokenPoolState<T>",
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
func (c burnMintTokenPoolEncoder) IsRemotePoolWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&BurnMintTokenPoolState<T>",
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
func (c burnMintTokenPoolEncoder) GetRemoteToken(typeArgs []string, state bind.Object, remoteChainSelector uint64) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("get_remote_token", typeArgsList, typeParamsList, []string{
		"&BurnMintTokenPoolState<T>",
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
func (c burnMintTokenPoolEncoder) GetRemoteTokenWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&BurnMintTokenPoolState<T>",
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
func (c burnMintTokenPoolEncoder) AddRemotePool(typeArgs []string, state bind.Object, ownerCap bind.Object, remoteChainSelector uint64, remotePoolAddress []byte) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("add_remote_pool", typeArgsList, typeParamsList, []string{
		"&mut BurnMintTokenPoolState<T>",
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
func (c burnMintTokenPoolEncoder) AddRemotePoolWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut BurnMintTokenPoolState<T>",
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
func (c burnMintTokenPoolEncoder) RemoveRemotePool(typeArgs []string, state bind.Object, ownerCap bind.Object, remoteChainSelector uint64, remotePoolAddress []byte) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("remove_remote_pool", typeArgsList, typeParamsList, []string{
		"&mut BurnMintTokenPoolState<T>",
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
func (c burnMintTokenPoolEncoder) RemoveRemotePoolWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut BurnMintTokenPoolState<T>",
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
func (c burnMintTokenPoolEncoder) IsSupportedChain(typeArgs []string, state bind.Object, remoteChainSelector uint64) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("is_supported_chain", typeArgsList, typeParamsList, []string{
		"&BurnMintTokenPoolState<T>",
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
func (c burnMintTokenPoolEncoder) IsSupportedChainWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&BurnMintTokenPoolState<T>",
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
func (c burnMintTokenPoolEncoder) GetSupportedChains(typeArgs []string, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("get_supported_chains", typeArgsList, typeParamsList, []string{
		"&BurnMintTokenPoolState<T>",
	}, []any{
		state,
	}, []string{
		"vector<u64>",
	})
}

// GetSupportedChainsWithArgs encodes a call to the get_supported_chains Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c burnMintTokenPoolEncoder) GetSupportedChainsWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&BurnMintTokenPoolState<T>",
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
func (c burnMintTokenPoolEncoder) ApplyChainUpdates(typeArgs []string, state bind.Object, ownerCap bind.Object, remoteChainSelectorsToRemove []uint64, remoteChainSelectorsToAdd []uint64, remotePoolAddressesToAdd [][][]byte, remoteTokenAddressesToAdd [][]byte) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("apply_chain_updates", typeArgsList, typeParamsList, []string{
		"&mut BurnMintTokenPoolState<T>",
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
func (c burnMintTokenPoolEncoder) ApplyChainUpdatesWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut BurnMintTokenPoolState<T>",
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
func (c burnMintTokenPoolEncoder) GetAllowlistEnabled(typeArgs []string, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("get_allowlist_enabled", typeArgsList, typeParamsList, []string{
		"&BurnMintTokenPoolState<T>",
	}, []any{
		state,
	}, []string{
		"bool",
	})
}

// GetAllowlistEnabledWithArgs encodes a call to the get_allowlist_enabled Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c burnMintTokenPoolEncoder) GetAllowlistEnabledWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&BurnMintTokenPoolState<T>",
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
func (c burnMintTokenPoolEncoder) GetAllowlist(typeArgs []string, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("get_allowlist", typeArgsList, typeParamsList, []string{
		"&BurnMintTokenPoolState<T>",
	}, []any{
		state,
	}, []string{
		"vector<address>",
	})
}

// GetAllowlistWithArgs encodes a call to the get_allowlist Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c burnMintTokenPoolEncoder) GetAllowlistWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&BurnMintTokenPoolState<T>",
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
func (c burnMintTokenPoolEncoder) SetAllowlistEnabled(typeArgs []string, state bind.Object, ownerCap bind.Object, enabled bool) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("set_allowlist_enabled", typeArgsList, typeParamsList, []string{
		"&mut BurnMintTokenPoolState<T>",
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
func (c burnMintTokenPoolEncoder) SetAllowlistEnabledWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut BurnMintTokenPoolState<T>",
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
func (c burnMintTokenPoolEncoder) ApplyAllowlistUpdates(typeArgs []string, state bind.Object, ownerCap bind.Object, removes []string, adds []string) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("apply_allowlist_updates", typeArgsList, typeParamsList, []string{
		"&mut BurnMintTokenPoolState<T>",
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
func (c burnMintTokenPoolEncoder) ApplyAllowlistUpdatesWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut BurnMintTokenPoolState<T>",
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

// LockOrBurn encodes a call to the lock_or_burn Move function.
func (c burnMintTokenPoolEncoder) LockOrBurn(typeArgs []string, ref bind.Object, c_ bind.Object, tokenParams bind.Object, state bind.Object, clock bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("lock_or_burn", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"Coin<T>",
		"&mut dd::TokenParams",
		"&mut BurnMintTokenPoolState<T>",
		"&Clock",
	}, []any{
		ref,
		c_,
		tokenParams,
		state,
		clock,
	}, nil)
}

// LockOrBurnWithArgs encodes a call to the lock_or_burn Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c burnMintTokenPoolEncoder) LockOrBurnWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"Coin<T>",
		"&mut dd::TokenParams",
		"&mut BurnMintTokenPoolState<T>",
		"&Clock",
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
func (c burnMintTokenPoolEncoder) ReleaseOrMint(typeArgs []string, ref bind.Object, receiverParams bind.Object, index uint64, pool bind.Object, clock bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("release_or_mint", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"osh::ReceiverParams",
		"u64",
		"&mut BurnMintTokenPoolState<T>",
		"&Clock",
	}, []any{
		ref,
		receiverParams,
		index,
		pool,
		clock,
	}, []string{
		"osh::ReceiverParams",
	})
}

// ReleaseOrMintWithArgs encodes a call to the release_or_mint Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c burnMintTokenPoolEncoder) ReleaseOrMintWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"osh::ReceiverParams",
		"u64",
		"&mut BurnMintTokenPoolState<T>",
		"&Clock",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("release_or_mint", typeArgsList, typeParamsList, expectedParams, args, []string{
		"osh::ReceiverParams",
	})
}

// SetChainRateLimiterConfigs encodes a call to the set_chain_rate_limiter_configs Move function.
func (c burnMintTokenPoolEncoder) SetChainRateLimiterConfigs(typeArgs []string, state bind.Object, ownerCap bind.Object, clock bind.Object, remoteChainSelectors []uint64, outboundIsEnableds []bool, outboundCapacities []uint64, outboundRates []uint64, inboundIsEnableds []bool, inboundCapacities []uint64, inboundRates []uint64) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("set_chain_rate_limiter_configs", typeArgsList, typeParamsList, []string{
		"&mut BurnMintTokenPoolState<T>",
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
func (c burnMintTokenPoolEncoder) SetChainRateLimiterConfigsWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut BurnMintTokenPoolState<T>",
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
func (c burnMintTokenPoolEncoder) SetChainRateLimiterConfig(typeArgs []string, state bind.Object, ownerCap bind.Object, clock bind.Object, remoteChainSelector uint64, outboundIsEnabled bool, outboundCapacity uint64, outboundRate uint64, inboundIsEnabled bool, inboundCapacity uint64, inboundRate uint64) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("set_chain_rate_limiter_config", typeArgsList, typeParamsList, []string{
		"&mut BurnMintTokenPoolState<T>",
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
func (c burnMintTokenPoolEncoder) SetChainRateLimiterConfigWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut BurnMintTokenPoolState<T>",
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

// DestroyTokenPool encodes a call to the destroy_token_pool Move function.
func (c burnMintTokenPoolEncoder) DestroyTokenPool(typeArgs []string, state bind.Object, ownerCap bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("destroy_token_pool", typeArgsList, typeParamsList, []string{
		"BurnMintTokenPoolState<T>",
		"OwnerCap",
	}, []any{
		state,
		ownerCap,
	}, []string{
		"TreasuryCap<T>",
	})
}

// DestroyTokenPoolWithArgs encodes a call to the destroy_token_pool Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c burnMintTokenPoolEncoder) DestroyTokenPoolWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"BurnMintTokenPoolState<T>",
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
		"TreasuryCap<T>",
	})
}

// Owner encodes a call to the owner Move function.
func (c burnMintTokenPoolEncoder) Owner(typeArgs []string, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("owner", typeArgsList, typeParamsList, []string{
		"&BurnMintTokenPoolState<T>",
	}, []any{
		state,
	}, []string{
		"address",
	})
}

// OwnerWithArgs encodes a call to the owner Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c burnMintTokenPoolEncoder) OwnerWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&BurnMintTokenPoolState<T>",
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
func (c burnMintTokenPoolEncoder) HasPendingTransfer(typeArgs []string, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("has_pending_transfer", typeArgsList, typeParamsList, []string{
		"&BurnMintTokenPoolState<T>",
	}, []any{
		state,
	}, []string{
		"bool",
	})
}

// HasPendingTransferWithArgs encodes a call to the has_pending_transfer Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c burnMintTokenPoolEncoder) HasPendingTransferWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&BurnMintTokenPoolState<T>",
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
func (c burnMintTokenPoolEncoder) PendingTransferFrom(typeArgs []string, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("pending_transfer_from", typeArgsList, typeParamsList, []string{
		"&BurnMintTokenPoolState<T>",
	}, []any{
		state,
	}, []string{
		"0x1::option::Option<address>",
	})
}

// PendingTransferFromWithArgs encodes a call to the pending_transfer_from Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c burnMintTokenPoolEncoder) PendingTransferFromWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&BurnMintTokenPoolState<T>",
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
func (c burnMintTokenPoolEncoder) PendingTransferTo(typeArgs []string, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("pending_transfer_to", typeArgsList, typeParamsList, []string{
		"&BurnMintTokenPoolState<T>",
	}, []any{
		state,
	}, []string{
		"0x1::option::Option<address>",
	})
}

// PendingTransferToWithArgs encodes a call to the pending_transfer_to Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c burnMintTokenPoolEncoder) PendingTransferToWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&BurnMintTokenPoolState<T>",
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
func (c burnMintTokenPoolEncoder) PendingTransferAccepted(typeArgs []string, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("pending_transfer_accepted", typeArgsList, typeParamsList, []string{
		"&BurnMintTokenPoolState<T>",
	}, []any{
		state,
	}, []string{
		"0x1::option::Option<bool>",
	})
}

// PendingTransferAcceptedWithArgs encodes a call to the pending_transfer_accepted Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c burnMintTokenPoolEncoder) PendingTransferAcceptedWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&BurnMintTokenPoolState<T>",
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
func (c burnMintTokenPoolEncoder) TransferOwnership(typeArgs []string, state bind.Object, ownerCap bind.Object, newOwner string) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("transfer_ownership", typeArgsList, typeParamsList, []string{
		"&mut BurnMintTokenPoolState<T>",
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
func (c burnMintTokenPoolEncoder) TransferOwnershipWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut BurnMintTokenPoolState<T>",
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
func (c burnMintTokenPoolEncoder) AcceptOwnership(typeArgs []string, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("accept_ownership", typeArgsList, typeParamsList, []string{
		"&mut BurnMintTokenPoolState<T>",
	}, []any{
		state,
	}, nil)
}

// AcceptOwnershipWithArgs encodes a call to the accept_ownership Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c burnMintTokenPoolEncoder) AcceptOwnershipWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut BurnMintTokenPoolState<T>",
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
func (c burnMintTokenPoolEncoder) AcceptOwnershipFromObject(typeArgs []string, state bind.Object, from string) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("accept_ownership_from_object", typeArgsList, typeParamsList, []string{
		"&mut BurnMintTokenPoolState<T>",
		"&mut UID",
	}, []any{
		state,
		from,
	}, nil)
}

// AcceptOwnershipFromObjectWithArgs encodes a call to the accept_ownership_from_object Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c burnMintTokenPoolEncoder) AcceptOwnershipFromObjectWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut BurnMintTokenPoolState<T>",
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

// ExecuteOwnershipTransfer encodes a call to the execute_ownership_transfer Move function.
func (c burnMintTokenPoolEncoder) ExecuteOwnershipTransfer(ownerCap bind.Object, ownableState bind.Object, to string) (*bind.EncodedCall, error) {
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
func (c burnMintTokenPoolEncoder) ExecuteOwnershipTransferWithArgs(args ...any) (*bind.EncodedCall, error) {
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

// McmsRegisterEntrypoint encodes a call to the mcms_register_entrypoint Move function.
func (c burnMintTokenPoolEncoder) McmsRegisterEntrypoint(typeArgs []string, registry bind.Object, state bind.Object, ownerCap bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("mcms_register_entrypoint", typeArgsList, typeParamsList, []string{
		"&mut Registry",
		"&mut BurnMintTokenPoolState<T>",
		"OwnerCap",
	}, []any{
		registry,
		state,
		ownerCap,
	}, nil)
}

// McmsRegisterEntrypointWithArgs encodes a call to the mcms_register_entrypoint Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c burnMintTokenPoolEncoder) McmsRegisterEntrypointWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut Registry",
		"&mut BurnMintTokenPoolState<T>",
		"OwnerCap",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("mcms_register_entrypoint", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// McmsRegisterUpgradeCap encodes a call to the mcms_register_upgrade_cap Move function.
func (c burnMintTokenPoolEncoder) McmsRegisterUpgradeCap(upgradeCap bind.Object, registry bind.Object, state bind.Object) (*bind.EncodedCall, error) {
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
func (c burnMintTokenPoolEncoder) McmsRegisterUpgradeCapWithArgs(args ...any) (*bind.EncodedCall, error) {
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
func (c burnMintTokenPoolEncoder) McmsEntrypoint(typeArgs []string, state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("mcms_entrypoint", typeArgsList, typeParamsList, []string{
		"&mut BurnMintTokenPoolState<T>",
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
func (c burnMintTokenPoolEncoder) McmsEntrypointWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut BurnMintTokenPoolState<T>",
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
	return c.EncodeCallArgsWithGenerics("mcms_entrypoint", typeArgsList, typeParamsList, expectedParams, args, nil)
}
