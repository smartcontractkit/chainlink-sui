// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package module_lock_release_token_pool

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

type ILockReleaseTokenPool interface {
	TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	Initialize(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, coinMetadata bind.Object, treasuryCap bind.Object, lockReleaseTokenPoolPackageId string, tokenPoolAdministrator string, rebalancer string) (*models.SuiTransactionBlockResponse, error)
	InitializeByCcipAdmin(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, ownerCap bind.Object, coinMetadata bind.Object, lockReleaseTokenPoolPackageId string, tokenPoolAdministrator string, rebalancer string) (*models.SuiTransactionBlockResponse, error)
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
	LockOrBurn(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, c_ bind.Object, tokenParams bind.Object, clock bind.Object, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	ReleaseOrMint(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, receiverParams bind.Object, index uint64, clock bind.Object, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	SetChainRateLimiterConfigs(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object, clock bind.Object, remoteChainSelectors []uint64, outboundIsEnableds []bool, outboundCapacities []uint64, outboundRates []uint64, inboundIsEnableds []bool, inboundCapacities []uint64, inboundRates []uint64) (*models.SuiTransactionBlockResponse, error)
	SetChainRateLimiterConfig(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object, clock bind.Object, remoteChainSelector uint64, outboundIsEnabled bool, outboundCapacity uint64, outboundRate uint64, inboundIsEnabled bool, inboundCapacity uint64, inboundRate uint64) (*models.SuiTransactionBlockResponse, error)
	ProvideLiquidity(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, c_ bind.Object) (*models.SuiTransactionBlockResponse, error)
	WithdrawLiquidity(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, amount uint64) (*models.SuiTransactionBlockResponse, error)
	SetRebalancer(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ownerCap bind.Object, state bind.Object, rebalancer string) (*models.SuiTransactionBlockResponse, error)
	GetRebalancer(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	GetBalance(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error)
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
	DestroyTokenPool(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object) (*models.SuiTransactionBlockResponse, error)
	DevInspect() ILockReleaseTokenPoolDevInspect
	Encoder() LockReleaseTokenPoolEncoder
}

type ILockReleaseTokenPoolDevInspect interface {
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
	ReleaseOrMint(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, receiverParams bind.Object, index uint64, clock bind.Object, state bind.Object) (bind.Object, error)
	WithdrawLiquidity(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, amount uint64) (any, error)
	GetRebalancer(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (string, error)
	GetBalance(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (uint64, error)
	Owner(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (string, error)
	HasPendingTransfer(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (bool, error)
	PendingTransferFrom(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*string, error)
	PendingTransferTo(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*string, error)
	PendingTransferAccepted(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*bool, error)
	DestroyTokenPool(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object) (any, error)
}

type LockReleaseTokenPoolEncoder interface {
	TypeAndVersion() (*bind.EncodedCall, error)
	TypeAndVersionWithArgs(args ...any) (*bind.EncodedCall, error)
	Initialize(typeArgs []string, ref bind.Object, coinMetadata bind.Object, treasuryCap bind.Object, lockReleaseTokenPoolPackageId string, tokenPoolAdministrator string, rebalancer string) (*bind.EncodedCall, error)
	InitializeWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	InitializeByCcipAdmin(typeArgs []string, ref bind.Object, ownerCap bind.Object, coinMetadata bind.Object, lockReleaseTokenPoolPackageId string, tokenPoolAdministrator string, rebalancer string) (*bind.EncodedCall, error)
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
	LockOrBurn(typeArgs []string, ref bind.Object, c_ bind.Object, tokenParams bind.Object, clock bind.Object, state bind.Object) (*bind.EncodedCall, error)
	LockOrBurnWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	ReleaseOrMint(typeArgs []string, ref bind.Object, receiverParams bind.Object, index uint64, clock bind.Object, state bind.Object) (*bind.EncodedCall, error)
	ReleaseOrMintWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	SetChainRateLimiterConfigs(typeArgs []string, state bind.Object, ownerCap bind.Object, clock bind.Object, remoteChainSelectors []uint64, outboundIsEnableds []bool, outboundCapacities []uint64, outboundRates []uint64, inboundIsEnableds []bool, inboundCapacities []uint64, inboundRates []uint64) (*bind.EncodedCall, error)
	SetChainRateLimiterConfigsWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	SetChainRateLimiterConfig(typeArgs []string, state bind.Object, ownerCap bind.Object, clock bind.Object, remoteChainSelector uint64, outboundIsEnabled bool, outboundCapacity uint64, outboundRate uint64, inboundIsEnabled bool, inboundCapacity uint64, inboundRate uint64) (*bind.EncodedCall, error)
	SetChainRateLimiterConfigWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	ProvideLiquidity(typeArgs []string, state bind.Object, c_ bind.Object) (*bind.EncodedCall, error)
	ProvideLiquidityWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	WithdrawLiquidity(typeArgs []string, state bind.Object, amount uint64) (*bind.EncodedCall, error)
	WithdrawLiquidityWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	SetRebalancer(typeArgs []string, ownerCap bind.Object, state bind.Object, rebalancer string) (*bind.EncodedCall, error)
	SetRebalancerWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	GetRebalancer(typeArgs []string, state bind.Object) (*bind.EncodedCall, error)
	GetRebalancerWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	GetBalance(typeArgs []string, state bind.Object) (*bind.EncodedCall, error)
	GetBalanceWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
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
	DestroyTokenPool(typeArgs []string, state bind.Object, ownerCap bind.Object) (*bind.EncodedCall, error)
	DestroyTokenPoolWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
}

type LockReleaseTokenPoolContract struct {
	*bind.BoundContract
	lockReleaseTokenPoolEncoder
	devInspect *LockReleaseTokenPoolDevInspect
}

type LockReleaseTokenPoolDevInspect struct {
	contract *LockReleaseTokenPoolContract
}

var _ ILockReleaseTokenPool = (*LockReleaseTokenPoolContract)(nil)
var _ ILockReleaseTokenPoolDevInspect = (*LockReleaseTokenPoolDevInspect)(nil)

func NewLockReleaseTokenPool(packageID string, client sui.ISuiAPI) (*LockReleaseTokenPoolContract, error) {
	contract, err := bind.NewBoundContract(packageID, "lock_release_token_pool", "lock_release_token_pool", client)
	if err != nil {
		return nil, err
	}

	c := &LockReleaseTokenPoolContract{
		BoundContract:               contract,
		lockReleaseTokenPoolEncoder: lockReleaseTokenPoolEncoder{BoundContract: contract},
	}
	c.devInspect = &LockReleaseTokenPoolDevInspect{contract: c}
	return c, nil
}

func (c *LockReleaseTokenPoolContract) Encoder() LockReleaseTokenPoolEncoder {
	return c.lockReleaseTokenPoolEncoder
}

func (c *LockReleaseTokenPoolContract) DevInspect() ILockReleaseTokenPoolDevInspect {
	return c.devInspect
}

func (c *LockReleaseTokenPoolContract) BuildPTB(ctx context.Context, ptb *transaction.Transaction, encoded *bind.EncodedCall) (*transaction.Argument, error) {
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

type LockReleaseTokenPoolState struct {
	Id             string      `move:"sui::object::UID"`
	TokenPoolState bind.Object `move:"TokenPoolState"`
	Reserve        bind.Object `move:"Coin<T>"`
	Rebalancer     string      `move:"address"`
	OwnableState   bind.Object `move:"OwnableState"`
}

type TypeProof struct {
}

type McmsCallback struct {
}

type bcsLockReleaseTokenPoolState struct {
	Id             string
	TokenPoolState bind.Object
	Reserve        bind.Object
	Rebalancer     [32]byte
	OwnableState   bind.Object
}

func convertLockReleaseTokenPoolStateFromBCS(bcs bcsLockReleaseTokenPoolState) LockReleaseTokenPoolState {
	return LockReleaseTokenPoolState{
		Id:             bcs.Id,
		TokenPoolState: bcs.TokenPoolState,
		Reserve:        bcs.Reserve,
		Rebalancer:     fmt.Sprintf("0x%x", bcs.Rebalancer),
		OwnableState:   bcs.OwnableState,
	}
}

func init() {
	bind.RegisterStructDecoder("lock_release_token_pool::lock_release_token_pool::LockReleaseTokenPoolState", func(data []byte) (interface{}, error) {
		var temp bcsLockReleaseTokenPoolState
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result := convertLockReleaseTokenPoolStateFromBCS(temp)
		return result, nil
	})
	bind.RegisterStructDecoder("lock_release_token_pool::lock_release_token_pool::TypeProof", func(data []byte) (interface{}, error) {
		var result TypeProof
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("lock_release_token_pool::lock_release_token_pool::McmsCallback", func(data []byte) (interface{}, error) {
		var result McmsCallback
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
}

// TypeAndVersion executes the type_and_version Move function.
func (c *LockReleaseTokenPoolContract) TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.lockReleaseTokenPoolEncoder.TypeAndVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// Initialize executes the initialize Move function.
func (c *LockReleaseTokenPoolContract) Initialize(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, coinMetadata bind.Object, treasuryCap bind.Object, lockReleaseTokenPoolPackageId string, tokenPoolAdministrator string, rebalancer string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.lockReleaseTokenPoolEncoder.Initialize(typeArgs, ref, coinMetadata, treasuryCap, lockReleaseTokenPoolPackageId, tokenPoolAdministrator, rebalancer)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// InitializeByCcipAdmin executes the initialize_by_ccip_admin Move function.
func (c *LockReleaseTokenPoolContract) InitializeByCcipAdmin(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, ownerCap bind.Object, coinMetadata bind.Object, lockReleaseTokenPoolPackageId string, tokenPoolAdministrator string, rebalancer string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.lockReleaseTokenPoolEncoder.InitializeByCcipAdmin(typeArgs, ref, ownerCap, coinMetadata, lockReleaseTokenPoolPackageId, tokenPoolAdministrator, rebalancer)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetToken executes the get_token Move function.
func (c *LockReleaseTokenPoolContract) GetToken(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.lockReleaseTokenPoolEncoder.GetToken(typeArgs, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetTokenDecimals executes the get_token_decimals Move function.
func (c *LockReleaseTokenPoolContract) GetTokenDecimals(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.lockReleaseTokenPoolEncoder.GetTokenDecimals(typeArgs, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetRemotePools executes the get_remote_pools Move function.
func (c *LockReleaseTokenPoolContract) GetRemotePools(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.lockReleaseTokenPoolEncoder.GetRemotePools(typeArgs, state, remoteChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// IsRemotePool executes the is_remote_pool Move function.
func (c *LockReleaseTokenPoolContract) IsRemotePool(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, remoteChainSelector uint64, remotePoolAddress []byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.lockReleaseTokenPoolEncoder.IsRemotePool(typeArgs, state, remoteChainSelector, remotePoolAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetRemoteToken executes the get_remote_token Move function.
func (c *LockReleaseTokenPoolContract) GetRemoteToken(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.lockReleaseTokenPoolEncoder.GetRemoteToken(typeArgs, state, remoteChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// AddRemotePool executes the add_remote_pool Move function.
func (c *LockReleaseTokenPoolContract) AddRemotePool(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object, remoteChainSelector uint64, remotePoolAddress []byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.lockReleaseTokenPoolEncoder.AddRemotePool(typeArgs, state, ownerCap, remoteChainSelector, remotePoolAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// RemoveRemotePool executes the remove_remote_pool Move function.
func (c *LockReleaseTokenPoolContract) RemoveRemotePool(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object, remoteChainSelector uint64, remotePoolAddress []byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.lockReleaseTokenPoolEncoder.RemoveRemotePool(typeArgs, state, ownerCap, remoteChainSelector, remotePoolAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// IsSupportedChain executes the is_supported_chain Move function.
func (c *LockReleaseTokenPoolContract) IsSupportedChain(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.lockReleaseTokenPoolEncoder.IsSupportedChain(typeArgs, state, remoteChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetSupportedChains executes the get_supported_chains Move function.
func (c *LockReleaseTokenPoolContract) GetSupportedChains(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.lockReleaseTokenPoolEncoder.GetSupportedChains(typeArgs, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ApplyChainUpdates executes the apply_chain_updates Move function.
func (c *LockReleaseTokenPoolContract) ApplyChainUpdates(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object, remoteChainSelectorsToRemove []uint64, remoteChainSelectorsToAdd []uint64, remotePoolAddressesToAdd [][][]byte, remoteTokenAddressesToAdd [][]byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.lockReleaseTokenPoolEncoder.ApplyChainUpdates(typeArgs, state, ownerCap, remoteChainSelectorsToRemove, remoteChainSelectorsToAdd, remotePoolAddressesToAdd, remoteTokenAddressesToAdd)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetAllowlistEnabled executes the get_allowlist_enabled Move function.
func (c *LockReleaseTokenPoolContract) GetAllowlistEnabled(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.lockReleaseTokenPoolEncoder.GetAllowlistEnabled(typeArgs, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetAllowlist executes the get_allowlist Move function.
func (c *LockReleaseTokenPoolContract) GetAllowlist(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.lockReleaseTokenPoolEncoder.GetAllowlist(typeArgs, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// SetAllowlistEnabled executes the set_allowlist_enabled Move function.
func (c *LockReleaseTokenPoolContract) SetAllowlistEnabled(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object, enabled bool) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.lockReleaseTokenPoolEncoder.SetAllowlistEnabled(typeArgs, state, ownerCap, enabled)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ApplyAllowlistUpdates executes the apply_allowlist_updates Move function.
func (c *LockReleaseTokenPoolContract) ApplyAllowlistUpdates(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object, removes []string, adds []string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.lockReleaseTokenPoolEncoder.ApplyAllowlistUpdates(typeArgs, state, ownerCap, removes, adds)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// LockOrBurn executes the lock_or_burn Move function.
func (c *LockReleaseTokenPoolContract) LockOrBurn(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, c_ bind.Object, tokenParams bind.Object, clock bind.Object, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.lockReleaseTokenPoolEncoder.LockOrBurn(typeArgs, ref, c_, tokenParams, clock, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ReleaseOrMint executes the release_or_mint Move function.
func (c *LockReleaseTokenPoolContract) ReleaseOrMint(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, receiverParams bind.Object, index uint64, clock bind.Object, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.lockReleaseTokenPoolEncoder.ReleaseOrMint(typeArgs, ref, receiverParams, index, clock, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// SetChainRateLimiterConfigs executes the set_chain_rate_limiter_configs Move function.
func (c *LockReleaseTokenPoolContract) SetChainRateLimiterConfigs(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object, clock bind.Object, remoteChainSelectors []uint64, outboundIsEnableds []bool, outboundCapacities []uint64, outboundRates []uint64, inboundIsEnableds []bool, inboundCapacities []uint64, inboundRates []uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.lockReleaseTokenPoolEncoder.SetChainRateLimiterConfigs(typeArgs, state, ownerCap, clock, remoteChainSelectors, outboundIsEnableds, outboundCapacities, outboundRates, inboundIsEnableds, inboundCapacities, inboundRates)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// SetChainRateLimiterConfig executes the set_chain_rate_limiter_config Move function.
func (c *LockReleaseTokenPoolContract) SetChainRateLimiterConfig(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object, clock bind.Object, remoteChainSelector uint64, outboundIsEnabled bool, outboundCapacity uint64, outboundRate uint64, inboundIsEnabled bool, inboundCapacity uint64, inboundRate uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.lockReleaseTokenPoolEncoder.SetChainRateLimiterConfig(typeArgs, state, ownerCap, clock, remoteChainSelector, outboundIsEnabled, outboundCapacity, outboundRate, inboundIsEnabled, inboundCapacity, inboundRate)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ProvideLiquidity executes the provide_liquidity Move function.
func (c *LockReleaseTokenPoolContract) ProvideLiquidity(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, c_ bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.lockReleaseTokenPoolEncoder.ProvideLiquidity(typeArgs, state, c_)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// WithdrawLiquidity executes the withdraw_liquidity Move function.
func (c *LockReleaseTokenPoolContract) WithdrawLiquidity(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, amount uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.lockReleaseTokenPoolEncoder.WithdrawLiquidity(typeArgs, state, amount)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// SetRebalancer executes the set_rebalancer Move function.
func (c *LockReleaseTokenPoolContract) SetRebalancer(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ownerCap bind.Object, state bind.Object, rebalancer string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.lockReleaseTokenPoolEncoder.SetRebalancer(typeArgs, ownerCap, state, rebalancer)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetRebalancer executes the get_rebalancer Move function.
func (c *LockReleaseTokenPoolContract) GetRebalancer(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.lockReleaseTokenPoolEncoder.GetRebalancer(typeArgs, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetBalance executes the get_balance Move function.
func (c *LockReleaseTokenPoolContract) GetBalance(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.lockReleaseTokenPoolEncoder.GetBalance(typeArgs, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// Owner executes the owner Move function.
func (c *LockReleaseTokenPoolContract) Owner(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.lockReleaseTokenPoolEncoder.Owner(typeArgs, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// HasPendingTransfer executes the has_pending_transfer Move function.
func (c *LockReleaseTokenPoolContract) HasPendingTransfer(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.lockReleaseTokenPoolEncoder.HasPendingTransfer(typeArgs, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// PendingTransferFrom executes the pending_transfer_from Move function.
func (c *LockReleaseTokenPoolContract) PendingTransferFrom(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.lockReleaseTokenPoolEncoder.PendingTransferFrom(typeArgs, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// PendingTransferTo executes the pending_transfer_to Move function.
func (c *LockReleaseTokenPoolContract) PendingTransferTo(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.lockReleaseTokenPoolEncoder.PendingTransferTo(typeArgs, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// PendingTransferAccepted executes the pending_transfer_accepted Move function.
func (c *LockReleaseTokenPoolContract) PendingTransferAccepted(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.lockReleaseTokenPoolEncoder.PendingTransferAccepted(typeArgs, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// TransferOwnership executes the transfer_ownership Move function.
func (c *LockReleaseTokenPoolContract) TransferOwnership(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object, newOwner string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.lockReleaseTokenPoolEncoder.TransferOwnership(typeArgs, state, ownerCap, newOwner)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// AcceptOwnership executes the accept_ownership Move function.
func (c *LockReleaseTokenPoolContract) AcceptOwnership(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.lockReleaseTokenPoolEncoder.AcceptOwnership(typeArgs, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// AcceptOwnershipFromObject executes the accept_ownership_from_object Move function.
func (c *LockReleaseTokenPoolContract) AcceptOwnershipFromObject(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, from string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.lockReleaseTokenPoolEncoder.AcceptOwnershipFromObject(typeArgs, state, from)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ExecuteOwnershipTransfer executes the execute_ownership_transfer Move function.
func (c *LockReleaseTokenPoolContract) ExecuteOwnershipTransfer(ctx context.Context, opts *bind.CallOpts, ownerCap bind.Object, ownableState bind.Object, to string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.lockReleaseTokenPoolEncoder.ExecuteOwnershipTransfer(ownerCap, ownableState, to)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsRegisterEntrypoint executes the mcms_register_entrypoint Move function.
func (c *LockReleaseTokenPoolContract) McmsRegisterEntrypoint(ctx context.Context, opts *bind.CallOpts, typeArgs []string, registry bind.Object, state bind.Object, ownerCap bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.lockReleaseTokenPoolEncoder.McmsRegisterEntrypoint(typeArgs, registry, state, ownerCap)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsRegisterUpgradeCap executes the mcms_register_upgrade_cap Move function.
func (c *LockReleaseTokenPoolContract) McmsRegisterUpgradeCap(ctx context.Context, opts *bind.CallOpts, upgradeCap bind.Object, registry bind.Object, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.lockReleaseTokenPoolEncoder.McmsRegisterUpgradeCap(upgradeCap, registry, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsEntrypoint executes the mcms_entrypoint Move function.
func (c *LockReleaseTokenPoolContract) McmsEntrypoint(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.lockReleaseTokenPoolEncoder.McmsEntrypoint(typeArgs, state, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// DestroyTokenPool executes the destroy_token_pool Move function.
func (c *LockReleaseTokenPoolContract) DestroyTokenPool(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.lockReleaseTokenPoolEncoder.DestroyTokenPool(typeArgs, state, ownerCap)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// TypeAndVersion executes the type_and_version Move function using DevInspect to get return values.
//
// Returns: 0x1::string::String
func (d *LockReleaseTokenPoolDevInspect) TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (string, error) {
	encoded, err := d.contract.lockReleaseTokenPoolEncoder.TypeAndVersion()
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
func (d *LockReleaseTokenPoolDevInspect) GetToken(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (string, error) {
	encoded, err := d.contract.lockReleaseTokenPoolEncoder.GetToken(typeArgs, state)
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
func (d *LockReleaseTokenPoolDevInspect) GetTokenDecimals(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (byte, error) {
	encoded, err := d.contract.lockReleaseTokenPoolEncoder.GetTokenDecimals(typeArgs, state)
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
func (d *LockReleaseTokenPoolDevInspect) GetRemotePools(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, remoteChainSelector uint64) ([][]byte, error) {
	encoded, err := d.contract.lockReleaseTokenPoolEncoder.GetRemotePools(typeArgs, state, remoteChainSelector)
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
func (d *LockReleaseTokenPoolDevInspect) IsRemotePool(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, remoteChainSelector uint64, remotePoolAddress []byte) (bool, error) {
	encoded, err := d.contract.lockReleaseTokenPoolEncoder.IsRemotePool(typeArgs, state, remoteChainSelector, remotePoolAddress)
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
func (d *LockReleaseTokenPoolDevInspect) GetRemoteToken(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, remoteChainSelector uint64) ([]byte, error) {
	encoded, err := d.contract.lockReleaseTokenPoolEncoder.GetRemoteToken(typeArgs, state, remoteChainSelector)
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
func (d *LockReleaseTokenPoolDevInspect) IsSupportedChain(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, remoteChainSelector uint64) (bool, error) {
	encoded, err := d.contract.lockReleaseTokenPoolEncoder.IsSupportedChain(typeArgs, state, remoteChainSelector)
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
func (d *LockReleaseTokenPoolDevInspect) GetSupportedChains(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) ([]uint64, error) {
	encoded, err := d.contract.lockReleaseTokenPoolEncoder.GetSupportedChains(typeArgs, state)
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
func (d *LockReleaseTokenPoolDevInspect) GetAllowlistEnabled(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (bool, error) {
	encoded, err := d.contract.lockReleaseTokenPoolEncoder.GetAllowlistEnabled(typeArgs, state)
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
func (d *LockReleaseTokenPoolDevInspect) GetAllowlist(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) ([]string, error) {
	encoded, err := d.contract.lockReleaseTokenPoolEncoder.GetAllowlist(typeArgs, state)
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
func (d *LockReleaseTokenPoolDevInspect) ReleaseOrMint(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, receiverParams bind.Object, index uint64, clock bind.Object, state bind.Object) (bind.Object, error) {
	encoded, err := d.contract.lockReleaseTokenPoolEncoder.ReleaseOrMint(typeArgs, ref, receiverParams, index, clock, state)
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

// WithdrawLiquidity executes the withdraw_liquidity Move function using DevInspect to get return values.
//
// Returns: Coin<T>
func (d *LockReleaseTokenPoolDevInspect) WithdrawLiquidity(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, amount uint64) (any, error) {
	encoded, err := d.contract.lockReleaseTokenPoolEncoder.WithdrawLiquidity(typeArgs, state, amount)
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

// GetRebalancer executes the get_rebalancer Move function using DevInspect to get return values.
//
// Returns: address
func (d *LockReleaseTokenPoolDevInspect) GetRebalancer(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (string, error) {
	encoded, err := d.contract.lockReleaseTokenPoolEncoder.GetRebalancer(typeArgs, state)
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

// GetBalance executes the get_balance Move function using DevInspect to get return values.
//
// Returns: u64
func (d *LockReleaseTokenPoolDevInspect) GetBalance(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (uint64, error) {
	encoded, err := d.contract.lockReleaseTokenPoolEncoder.GetBalance(typeArgs, state)
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
	result, ok := results[0].(uint64)
	if !ok {
		return 0, fmt.Errorf("unexpected return type: expected uint64, got %T", results[0])
	}
	return result, nil
}

// Owner executes the owner Move function using DevInspect to get return values.
//
// Returns: address
func (d *LockReleaseTokenPoolDevInspect) Owner(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (string, error) {
	encoded, err := d.contract.lockReleaseTokenPoolEncoder.Owner(typeArgs, state)
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
func (d *LockReleaseTokenPoolDevInspect) HasPendingTransfer(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (bool, error) {
	encoded, err := d.contract.lockReleaseTokenPoolEncoder.HasPendingTransfer(typeArgs, state)
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
func (d *LockReleaseTokenPoolDevInspect) PendingTransferFrom(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*string, error) {
	encoded, err := d.contract.lockReleaseTokenPoolEncoder.PendingTransferFrom(typeArgs, state)
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
func (d *LockReleaseTokenPoolDevInspect) PendingTransferTo(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*string, error) {
	encoded, err := d.contract.lockReleaseTokenPoolEncoder.PendingTransferTo(typeArgs, state)
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
func (d *LockReleaseTokenPoolDevInspect) PendingTransferAccepted(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*bool, error) {
	encoded, err := d.contract.lockReleaseTokenPoolEncoder.PendingTransferAccepted(typeArgs, state)
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
// Returns: Coin<T>
func (d *LockReleaseTokenPoolDevInspect) DestroyTokenPool(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object) (any, error) {
	encoded, err := d.contract.lockReleaseTokenPoolEncoder.DestroyTokenPool(typeArgs, state, ownerCap)
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

type lockReleaseTokenPoolEncoder struct {
	*bind.BoundContract
}

// TypeAndVersion encodes a call to the type_and_version Move function.
func (c lockReleaseTokenPoolEncoder) TypeAndVersion() (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("type_and_version", typeArgsList, typeParamsList, []string{}, []any{}, []string{
		"0x1::string::String",
	})
}

// TypeAndVersionWithArgs encodes a call to the type_and_version Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c lockReleaseTokenPoolEncoder) TypeAndVersionWithArgs(args ...any) (*bind.EncodedCall, error) {
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
func (c lockReleaseTokenPoolEncoder) Initialize(typeArgs []string, ref bind.Object, coinMetadata bind.Object, treasuryCap bind.Object, lockReleaseTokenPoolPackageId string, tokenPoolAdministrator string, rebalancer string) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("initialize", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"&CoinMetadata<T>",
		"&TreasuryCap<T>",
		"address",
		"address",
		"address",
	}, []any{
		ref,
		coinMetadata,
		treasuryCap,
		lockReleaseTokenPoolPackageId,
		tokenPoolAdministrator,
		rebalancer,
	}, nil)
}

// InitializeWithArgs encodes a call to the initialize Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c lockReleaseTokenPoolEncoder) InitializeWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"&CoinMetadata<T>",
		"&TreasuryCap<T>",
		"address",
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

// InitializeByCcipAdmin encodes a call to the initialize_by_ccip_admin Move function.
func (c lockReleaseTokenPoolEncoder) InitializeByCcipAdmin(typeArgs []string, ref bind.Object, ownerCap bind.Object, coinMetadata bind.Object, lockReleaseTokenPoolPackageId string, tokenPoolAdministrator string, rebalancer string) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("initialize_by_ccip_admin", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"&state_object::OwnerCap",
		"&CoinMetadata<T>",
		"address",
		"address",
		"address",
	}, []any{
		ref,
		ownerCap,
		coinMetadata,
		lockReleaseTokenPoolPackageId,
		tokenPoolAdministrator,
		rebalancer,
	}, nil)
}

// InitializeByCcipAdminWithArgs encodes a call to the initialize_by_ccip_admin Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c lockReleaseTokenPoolEncoder) InitializeByCcipAdminWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"&state_object::OwnerCap",
		"&CoinMetadata<T>",
		"address",
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
	return c.EncodeCallArgsWithGenerics("initialize_by_ccip_admin", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// GetToken encodes a call to the get_token Move function.
func (c lockReleaseTokenPoolEncoder) GetToken(typeArgs []string, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("get_token", typeArgsList, typeParamsList, []string{
		"&LockReleaseTokenPoolState<T>",
	}, []any{
		state,
	}, []string{
		"address",
	})
}

// GetTokenWithArgs encodes a call to the get_token Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c lockReleaseTokenPoolEncoder) GetTokenWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&LockReleaseTokenPoolState<T>",
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
func (c lockReleaseTokenPoolEncoder) GetTokenDecimals(typeArgs []string, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("get_token_decimals", typeArgsList, typeParamsList, []string{
		"&LockReleaseTokenPoolState<T>",
	}, []any{
		state,
	}, []string{
		"u8",
	})
}

// GetTokenDecimalsWithArgs encodes a call to the get_token_decimals Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c lockReleaseTokenPoolEncoder) GetTokenDecimalsWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&LockReleaseTokenPoolState<T>",
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
func (c lockReleaseTokenPoolEncoder) GetRemotePools(typeArgs []string, state bind.Object, remoteChainSelector uint64) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("get_remote_pools", typeArgsList, typeParamsList, []string{
		"&LockReleaseTokenPoolState<T>",
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
func (c lockReleaseTokenPoolEncoder) GetRemotePoolsWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&LockReleaseTokenPoolState<T>",
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
func (c lockReleaseTokenPoolEncoder) IsRemotePool(typeArgs []string, state bind.Object, remoteChainSelector uint64, remotePoolAddress []byte) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("is_remote_pool", typeArgsList, typeParamsList, []string{
		"&LockReleaseTokenPoolState<T>",
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
func (c lockReleaseTokenPoolEncoder) IsRemotePoolWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&LockReleaseTokenPoolState<T>",
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
func (c lockReleaseTokenPoolEncoder) GetRemoteToken(typeArgs []string, state bind.Object, remoteChainSelector uint64) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("get_remote_token", typeArgsList, typeParamsList, []string{
		"&LockReleaseTokenPoolState<T>",
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
func (c lockReleaseTokenPoolEncoder) GetRemoteTokenWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&LockReleaseTokenPoolState<T>",
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
func (c lockReleaseTokenPoolEncoder) AddRemotePool(typeArgs []string, state bind.Object, ownerCap bind.Object, remoteChainSelector uint64, remotePoolAddress []byte) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("add_remote_pool", typeArgsList, typeParamsList, []string{
		"&mut LockReleaseTokenPoolState<T>",
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
func (c lockReleaseTokenPoolEncoder) AddRemotePoolWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut LockReleaseTokenPoolState<T>",
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
func (c lockReleaseTokenPoolEncoder) RemoveRemotePool(typeArgs []string, state bind.Object, ownerCap bind.Object, remoteChainSelector uint64, remotePoolAddress []byte) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("remove_remote_pool", typeArgsList, typeParamsList, []string{
		"&mut LockReleaseTokenPoolState<T>",
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
func (c lockReleaseTokenPoolEncoder) RemoveRemotePoolWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut LockReleaseTokenPoolState<T>",
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
func (c lockReleaseTokenPoolEncoder) IsSupportedChain(typeArgs []string, state bind.Object, remoteChainSelector uint64) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("is_supported_chain", typeArgsList, typeParamsList, []string{
		"&LockReleaseTokenPoolState<T>",
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
func (c lockReleaseTokenPoolEncoder) IsSupportedChainWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&LockReleaseTokenPoolState<T>",
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
func (c lockReleaseTokenPoolEncoder) GetSupportedChains(typeArgs []string, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("get_supported_chains", typeArgsList, typeParamsList, []string{
		"&LockReleaseTokenPoolState<T>",
	}, []any{
		state,
	}, []string{
		"vector<u64>",
	})
}

// GetSupportedChainsWithArgs encodes a call to the get_supported_chains Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c lockReleaseTokenPoolEncoder) GetSupportedChainsWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&LockReleaseTokenPoolState<T>",
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
func (c lockReleaseTokenPoolEncoder) ApplyChainUpdates(typeArgs []string, state bind.Object, ownerCap bind.Object, remoteChainSelectorsToRemove []uint64, remoteChainSelectorsToAdd []uint64, remotePoolAddressesToAdd [][][]byte, remoteTokenAddressesToAdd [][]byte) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("apply_chain_updates", typeArgsList, typeParamsList, []string{
		"&mut LockReleaseTokenPoolState<T>",
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
func (c lockReleaseTokenPoolEncoder) ApplyChainUpdatesWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut LockReleaseTokenPoolState<T>",
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
func (c lockReleaseTokenPoolEncoder) GetAllowlistEnabled(typeArgs []string, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("get_allowlist_enabled", typeArgsList, typeParamsList, []string{
		"&LockReleaseTokenPoolState<T>",
	}, []any{
		state,
	}, []string{
		"bool",
	})
}

// GetAllowlistEnabledWithArgs encodes a call to the get_allowlist_enabled Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c lockReleaseTokenPoolEncoder) GetAllowlistEnabledWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&LockReleaseTokenPoolState<T>",
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
func (c lockReleaseTokenPoolEncoder) GetAllowlist(typeArgs []string, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("get_allowlist", typeArgsList, typeParamsList, []string{
		"&LockReleaseTokenPoolState<T>",
	}, []any{
		state,
	}, []string{
		"vector<address>",
	})
}

// GetAllowlistWithArgs encodes a call to the get_allowlist Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c lockReleaseTokenPoolEncoder) GetAllowlistWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&LockReleaseTokenPoolState<T>",
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
func (c lockReleaseTokenPoolEncoder) SetAllowlistEnabled(typeArgs []string, state bind.Object, ownerCap bind.Object, enabled bool) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("set_allowlist_enabled", typeArgsList, typeParamsList, []string{
		"&mut LockReleaseTokenPoolState<T>",
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
func (c lockReleaseTokenPoolEncoder) SetAllowlistEnabledWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut LockReleaseTokenPoolState<T>",
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
func (c lockReleaseTokenPoolEncoder) ApplyAllowlistUpdates(typeArgs []string, state bind.Object, ownerCap bind.Object, removes []string, adds []string) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("apply_allowlist_updates", typeArgsList, typeParamsList, []string{
		"&mut LockReleaseTokenPoolState<T>",
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
func (c lockReleaseTokenPoolEncoder) ApplyAllowlistUpdatesWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut LockReleaseTokenPoolState<T>",
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
func (c lockReleaseTokenPoolEncoder) LockOrBurn(typeArgs []string, ref bind.Object, c_ bind.Object, tokenParams bind.Object, clock bind.Object, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("lock_or_burn", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"Coin<T>",
		"&mut dd::TokenParams",
		"&Clock",
		"&mut LockReleaseTokenPoolState<T>",
	}, []any{
		ref,
		c_,
		tokenParams,
		clock,
		state,
	}, nil)
}

// LockOrBurnWithArgs encodes a call to the lock_or_burn Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c lockReleaseTokenPoolEncoder) LockOrBurnWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"Coin<T>",
		"&mut dd::TokenParams",
		"&Clock",
		"&mut LockReleaseTokenPoolState<T>",
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
func (c lockReleaseTokenPoolEncoder) ReleaseOrMint(typeArgs []string, ref bind.Object, receiverParams bind.Object, index uint64, clock bind.Object, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("release_or_mint", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"osh::ReceiverParams",
		"u64",
		"&Clock",
		"&mut LockReleaseTokenPoolState<T>",
	}, []any{
		ref,
		receiverParams,
		index,
		clock,
		state,
	}, []string{
		"osh::ReceiverParams",
	})
}

// ReleaseOrMintWithArgs encodes a call to the release_or_mint Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c lockReleaseTokenPoolEncoder) ReleaseOrMintWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"osh::ReceiverParams",
		"u64",
		"&Clock",
		"&mut LockReleaseTokenPoolState<T>",
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
func (c lockReleaseTokenPoolEncoder) SetChainRateLimiterConfigs(typeArgs []string, state bind.Object, ownerCap bind.Object, clock bind.Object, remoteChainSelectors []uint64, outboundIsEnableds []bool, outboundCapacities []uint64, outboundRates []uint64, inboundIsEnableds []bool, inboundCapacities []uint64, inboundRates []uint64) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("set_chain_rate_limiter_configs", typeArgsList, typeParamsList, []string{
		"&mut LockReleaseTokenPoolState<T>",
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
func (c lockReleaseTokenPoolEncoder) SetChainRateLimiterConfigsWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut LockReleaseTokenPoolState<T>",
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
func (c lockReleaseTokenPoolEncoder) SetChainRateLimiterConfig(typeArgs []string, state bind.Object, ownerCap bind.Object, clock bind.Object, remoteChainSelector uint64, outboundIsEnabled bool, outboundCapacity uint64, outboundRate uint64, inboundIsEnabled bool, inboundCapacity uint64, inboundRate uint64) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("set_chain_rate_limiter_config", typeArgsList, typeParamsList, []string{
		"&mut LockReleaseTokenPoolState<T>",
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
func (c lockReleaseTokenPoolEncoder) SetChainRateLimiterConfigWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut LockReleaseTokenPoolState<T>",
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

// ProvideLiquidity encodes a call to the provide_liquidity Move function.
func (c lockReleaseTokenPoolEncoder) ProvideLiquidity(typeArgs []string, state bind.Object, c_ bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("provide_liquidity", typeArgsList, typeParamsList, []string{
		"&mut LockReleaseTokenPoolState<T>",
		"Coin<T>",
	}, []any{
		state,
		c_,
	}, nil)
}

// ProvideLiquidityWithArgs encodes a call to the provide_liquidity Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c lockReleaseTokenPoolEncoder) ProvideLiquidityWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut LockReleaseTokenPoolState<T>",
		"Coin<T>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("provide_liquidity", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// WithdrawLiquidity encodes a call to the withdraw_liquidity Move function.
func (c lockReleaseTokenPoolEncoder) WithdrawLiquidity(typeArgs []string, state bind.Object, amount uint64) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("withdraw_liquidity", typeArgsList, typeParamsList, []string{
		"&mut LockReleaseTokenPoolState<T>",
		"u64",
	}, []any{
		state,
		amount,
	}, []string{
		"Coin<T>",
	})
}

// WithdrawLiquidityWithArgs encodes a call to the withdraw_liquidity Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c lockReleaseTokenPoolEncoder) WithdrawLiquidityWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut LockReleaseTokenPoolState<T>",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("withdraw_liquidity", typeArgsList, typeParamsList, expectedParams, args, []string{
		"Coin<T>",
	})
}

// SetRebalancer encodes a call to the set_rebalancer Move function.
func (c lockReleaseTokenPoolEncoder) SetRebalancer(typeArgs []string, ownerCap bind.Object, state bind.Object, rebalancer string) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("set_rebalancer", typeArgsList, typeParamsList, []string{
		"&OwnerCap",
		"&mut LockReleaseTokenPoolState<T>",
		"address",
	}, []any{
		ownerCap,
		state,
		rebalancer,
	}, nil)
}

// SetRebalancerWithArgs encodes a call to the set_rebalancer Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c lockReleaseTokenPoolEncoder) SetRebalancerWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&OwnerCap",
		"&mut LockReleaseTokenPoolState<T>",
		"address",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("set_rebalancer", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// GetRebalancer encodes a call to the get_rebalancer Move function.
func (c lockReleaseTokenPoolEncoder) GetRebalancer(typeArgs []string, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("get_rebalancer", typeArgsList, typeParamsList, []string{
		"&LockReleaseTokenPoolState<T>",
	}, []any{
		state,
	}, []string{
		"address",
	})
}

// GetRebalancerWithArgs encodes a call to the get_rebalancer Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c lockReleaseTokenPoolEncoder) GetRebalancerWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&LockReleaseTokenPoolState<T>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("get_rebalancer", typeArgsList, typeParamsList, expectedParams, args, []string{
		"address",
	})
}

// GetBalance encodes a call to the get_balance Move function.
func (c lockReleaseTokenPoolEncoder) GetBalance(typeArgs []string, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("get_balance", typeArgsList, typeParamsList, []string{
		"&LockReleaseTokenPoolState<T>",
	}, []any{
		state,
	}, []string{
		"u64",
	})
}

// GetBalanceWithArgs encodes a call to the get_balance Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c lockReleaseTokenPoolEncoder) GetBalanceWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&LockReleaseTokenPoolState<T>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("get_balance", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u64",
	})
}

// Owner encodes a call to the owner Move function.
func (c lockReleaseTokenPoolEncoder) Owner(typeArgs []string, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("owner", typeArgsList, typeParamsList, []string{
		"&LockReleaseTokenPoolState<T>",
	}, []any{
		state,
	}, []string{
		"address",
	})
}

// OwnerWithArgs encodes a call to the owner Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c lockReleaseTokenPoolEncoder) OwnerWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&LockReleaseTokenPoolState<T>",
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
func (c lockReleaseTokenPoolEncoder) HasPendingTransfer(typeArgs []string, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("has_pending_transfer", typeArgsList, typeParamsList, []string{
		"&LockReleaseTokenPoolState<T>",
	}, []any{
		state,
	}, []string{
		"bool",
	})
}

// HasPendingTransferWithArgs encodes a call to the has_pending_transfer Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c lockReleaseTokenPoolEncoder) HasPendingTransferWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&LockReleaseTokenPoolState<T>",
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
func (c lockReleaseTokenPoolEncoder) PendingTransferFrom(typeArgs []string, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("pending_transfer_from", typeArgsList, typeParamsList, []string{
		"&LockReleaseTokenPoolState<T>",
	}, []any{
		state,
	}, []string{
		"0x1::option::Option<address>",
	})
}

// PendingTransferFromWithArgs encodes a call to the pending_transfer_from Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c lockReleaseTokenPoolEncoder) PendingTransferFromWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&LockReleaseTokenPoolState<T>",
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
func (c lockReleaseTokenPoolEncoder) PendingTransferTo(typeArgs []string, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("pending_transfer_to", typeArgsList, typeParamsList, []string{
		"&LockReleaseTokenPoolState<T>",
	}, []any{
		state,
	}, []string{
		"0x1::option::Option<address>",
	})
}

// PendingTransferToWithArgs encodes a call to the pending_transfer_to Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c lockReleaseTokenPoolEncoder) PendingTransferToWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&LockReleaseTokenPoolState<T>",
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
func (c lockReleaseTokenPoolEncoder) PendingTransferAccepted(typeArgs []string, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("pending_transfer_accepted", typeArgsList, typeParamsList, []string{
		"&LockReleaseTokenPoolState<T>",
	}, []any{
		state,
	}, []string{
		"0x1::option::Option<bool>",
	})
}

// PendingTransferAcceptedWithArgs encodes a call to the pending_transfer_accepted Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c lockReleaseTokenPoolEncoder) PendingTransferAcceptedWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&LockReleaseTokenPoolState<T>",
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
func (c lockReleaseTokenPoolEncoder) TransferOwnership(typeArgs []string, state bind.Object, ownerCap bind.Object, newOwner string) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("transfer_ownership", typeArgsList, typeParamsList, []string{
		"&mut LockReleaseTokenPoolState<T>",
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
func (c lockReleaseTokenPoolEncoder) TransferOwnershipWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut LockReleaseTokenPoolState<T>",
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
func (c lockReleaseTokenPoolEncoder) AcceptOwnership(typeArgs []string, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("accept_ownership", typeArgsList, typeParamsList, []string{
		"&mut LockReleaseTokenPoolState<T>",
	}, []any{
		state,
	}, nil)
}

// AcceptOwnershipWithArgs encodes a call to the accept_ownership Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c lockReleaseTokenPoolEncoder) AcceptOwnershipWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut LockReleaseTokenPoolState<T>",
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
func (c lockReleaseTokenPoolEncoder) AcceptOwnershipFromObject(typeArgs []string, state bind.Object, from string) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("accept_ownership_from_object", typeArgsList, typeParamsList, []string{
		"&mut LockReleaseTokenPoolState<T>",
		"&mut UID",
	}, []any{
		state,
		from,
	}, nil)
}

// AcceptOwnershipFromObjectWithArgs encodes a call to the accept_ownership_from_object Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c lockReleaseTokenPoolEncoder) AcceptOwnershipFromObjectWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut LockReleaseTokenPoolState<T>",
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
func (c lockReleaseTokenPoolEncoder) ExecuteOwnershipTransfer(ownerCap bind.Object, ownableState bind.Object, to string) (*bind.EncodedCall, error) {
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
func (c lockReleaseTokenPoolEncoder) ExecuteOwnershipTransferWithArgs(args ...any) (*bind.EncodedCall, error) {
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
func (c lockReleaseTokenPoolEncoder) McmsRegisterEntrypoint(typeArgs []string, registry bind.Object, state bind.Object, ownerCap bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("mcms_register_entrypoint", typeArgsList, typeParamsList, []string{
		"&mut Registry",
		"&mut LockReleaseTokenPoolState<T>",
		"OwnerCap",
	}, []any{
		registry,
		state,
		ownerCap,
	}, nil)
}

// McmsRegisterEntrypointWithArgs encodes a call to the mcms_register_entrypoint Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c lockReleaseTokenPoolEncoder) McmsRegisterEntrypointWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut Registry",
		"&mut LockReleaseTokenPoolState<T>",
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
func (c lockReleaseTokenPoolEncoder) McmsRegisterUpgradeCap(upgradeCap bind.Object, registry bind.Object, state bind.Object) (*bind.EncodedCall, error) {
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
func (c lockReleaseTokenPoolEncoder) McmsRegisterUpgradeCapWithArgs(args ...any) (*bind.EncodedCall, error) {
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
func (c lockReleaseTokenPoolEncoder) McmsEntrypoint(typeArgs []string, state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("mcms_entrypoint", typeArgsList, typeParamsList, []string{
		"&mut LockReleaseTokenPoolState<T>",
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
func (c lockReleaseTokenPoolEncoder) McmsEntrypointWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut LockReleaseTokenPoolState<T>",
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

// DestroyTokenPool encodes a call to the destroy_token_pool Move function.
func (c lockReleaseTokenPoolEncoder) DestroyTokenPool(typeArgs []string, state bind.Object, ownerCap bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("destroy_token_pool", typeArgsList, typeParamsList, []string{
		"LockReleaseTokenPoolState<T>",
		"OwnerCap",
	}, []any{
		state,
		ownerCap,
	}, []string{
		"Coin<T>",
	})
}

// DestroyTokenPoolWithArgs encodes a call to the destroy_token_pool Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c lockReleaseTokenPoolEncoder) DestroyTokenPoolWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"LockReleaseTokenPoolState<T>",
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
		"Coin<T>",
	})
}
