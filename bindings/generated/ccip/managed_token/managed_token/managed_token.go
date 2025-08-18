// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package module_managed_token

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

type IManagedToken interface {
	TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	Initialize(ctx context.Context, opts *bind.CallOpts, typeArgs []string, treasuryCap bind.Object) (*models.SuiTransactionBlockResponse, error)
	InitializeWithDenyCap(ctx context.Context, opts *bind.CallOpts, typeArgs []string, treasuryCap bind.Object, denyCap bind.Object) (*models.SuiTransactionBlockResponse, error)
	MintAllowance(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, mintCap bind.Object) (*models.SuiTransactionBlockResponse, error)
	TotalSupply(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	IsAuthorizedMintCap(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, id bind.Object) (*models.SuiTransactionBlockResponse, error)
	ConfigureNewMinter(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object, minter string, allowance uint64, isUnlimited bool) (*models.SuiTransactionBlockResponse, error)
	IncrementMintAllowance(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object, mintCapId bind.Object, denyList bind.Object, allowanceIncrement uint64) (*models.SuiTransactionBlockResponse, error)
	SetUnlimitedMintAllowances(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object, mintCapId bind.Object, denyList bind.Object, isUnlimited bool) (*models.SuiTransactionBlockResponse, error)
	GetAllMintCaps(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	MintAndTransfer(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, mintCap bind.Object, denyList bind.Object, amount uint64, recipient string) (*models.SuiTransactionBlockResponse, error)
	Mint(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, mintCap bind.Object, denyList bind.Object, amount uint64, recipient string) (*models.SuiTransactionBlockResponse, error)
	Burn(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, mintCap bind.Object, denyList bind.Object, coin bind.Object, from string) (*models.SuiTransactionBlockResponse, error)
	Blocklist(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object, denyList bind.Object, addr string) (*models.SuiTransactionBlockResponse, error)
	Unblocklist(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object, denyList bind.Object, addr string) (*models.SuiTransactionBlockResponse, error)
	Pause(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object, denyList bind.Object) (*models.SuiTransactionBlockResponse, error)
	Unpause(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object, denyList bind.Object) (*models.SuiTransactionBlockResponse, error)
	DestroyManagedToken(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ownerCap bind.Object, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	BorrowTreasuryCap(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object) (*models.SuiTransactionBlockResponse, error)
	Owner(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	HasPendingTransfer(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	PendingTransferFrom(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	PendingTransferTo(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	PendingTransferAccepted(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	TransferOwnership(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object, newOwner string) (*models.SuiTransactionBlockResponse, error)
	AcceptOwnership(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	AcceptOwnershipFromObject(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, from string) (*models.SuiTransactionBlockResponse, error)
	AcceptOwnershipAsMcms(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	ExecuteOwnershipTransfer(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ownerCap bind.Object, state bind.Object, to string) (*models.SuiTransactionBlockResponse, error)
	ExecuteOwnershipTransferToMcms(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ownerCap bind.Object, state bind.Object, registry bind.Object, to string) (*models.SuiTransactionBlockResponse, error)
	McmsRegisterUpgradeCap(ctx context.Context, opts *bind.CallOpts, upgradeCap bind.Object, registry bind.Object, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsEntrypoint(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, registry bind.Object, denyList bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	DevInspect() IManagedTokenDevInspect
	Encoder() ManagedTokenEncoder
}

type IManagedTokenDevInspect interface {
	TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (string, error)
	MintAllowance(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, mintCap bind.Object) ([]any, error)
	TotalSupply(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (uint64, error)
	IsAuthorizedMintCap(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, id bind.Object) (bool, error)
	GetAllMintCaps(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) ([]bind.Object, error)
	Mint(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, mintCap bind.Object, denyList bind.Object, amount uint64, recipient string) (any, error)
	DestroyManagedToken(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ownerCap bind.Object, state bind.Object) ([]any, error)
	BorrowTreasuryCap(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object) (any, error)
	Owner(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (string, error)
	HasPendingTransfer(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (bool, error)
	PendingTransferFrom(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*string, error)
	PendingTransferTo(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*string, error)
	PendingTransferAccepted(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*bool, error)
}

type ManagedTokenEncoder interface {
	TypeAndVersion() (*bind.EncodedCall, error)
	TypeAndVersionWithArgs(args ...any) (*bind.EncodedCall, error)
	Initialize(typeArgs []string, treasuryCap bind.Object) (*bind.EncodedCall, error)
	InitializeWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	InitializeWithDenyCap(typeArgs []string, treasuryCap bind.Object, denyCap bind.Object) (*bind.EncodedCall, error)
	InitializeWithDenyCapWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	MintAllowance(typeArgs []string, state bind.Object, mintCap bind.Object) (*bind.EncodedCall, error)
	MintAllowanceWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	TotalSupply(typeArgs []string, state bind.Object) (*bind.EncodedCall, error)
	TotalSupplyWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	IsAuthorizedMintCap(typeArgs []string, state bind.Object, id bind.Object) (*bind.EncodedCall, error)
	IsAuthorizedMintCapWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	ConfigureNewMinter(typeArgs []string, state bind.Object, ownerCap bind.Object, minter string, allowance uint64, isUnlimited bool) (*bind.EncodedCall, error)
	ConfigureNewMinterWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	IncrementMintAllowance(typeArgs []string, state bind.Object, ownerCap bind.Object, mintCapId bind.Object, denyList bind.Object, allowanceIncrement uint64) (*bind.EncodedCall, error)
	IncrementMintAllowanceWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	SetUnlimitedMintAllowances(typeArgs []string, state bind.Object, ownerCap bind.Object, mintCapId bind.Object, denyList bind.Object, isUnlimited bool) (*bind.EncodedCall, error)
	SetUnlimitedMintAllowancesWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	GetAllMintCaps(typeArgs []string, state bind.Object) (*bind.EncodedCall, error)
	GetAllMintCapsWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	MintAndTransfer(typeArgs []string, state bind.Object, mintCap bind.Object, denyList bind.Object, amount uint64, recipient string) (*bind.EncodedCall, error)
	MintAndTransferWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	Mint(typeArgs []string, state bind.Object, mintCap bind.Object, denyList bind.Object, amount uint64, recipient string) (*bind.EncodedCall, error)
	MintWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	Burn(typeArgs []string, state bind.Object, mintCap bind.Object, denyList bind.Object, coin bind.Object, from string) (*bind.EncodedCall, error)
	BurnWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	Blocklist(typeArgs []string, state bind.Object, ownerCap bind.Object, denyList bind.Object, addr string) (*bind.EncodedCall, error)
	BlocklistWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	Unblocklist(typeArgs []string, state bind.Object, ownerCap bind.Object, denyList bind.Object, addr string) (*bind.EncodedCall, error)
	UnblocklistWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	Pause(typeArgs []string, state bind.Object, ownerCap bind.Object, denyList bind.Object) (*bind.EncodedCall, error)
	PauseWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	Unpause(typeArgs []string, state bind.Object, ownerCap bind.Object, denyList bind.Object) (*bind.EncodedCall, error)
	UnpauseWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	DestroyManagedToken(typeArgs []string, ownerCap bind.Object, state bind.Object) (*bind.EncodedCall, error)
	DestroyManagedTokenWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	BorrowTreasuryCap(typeArgs []string, state bind.Object, ownerCap bind.Object) (*bind.EncodedCall, error)
	BorrowTreasuryCapWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
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
	AcceptOwnershipAsMcms(typeArgs []string, state bind.Object, params bind.Object) (*bind.EncodedCall, error)
	AcceptOwnershipAsMcmsWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	ExecuteOwnershipTransfer(typeArgs []string, ownerCap bind.Object, state bind.Object, to string) (*bind.EncodedCall, error)
	ExecuteOwnershipTransferWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	ExecuteOwnershipTransferToMcms(typeArgs []string, ownerCap bind.Object, state bind.Object, registry bind.Object, to string) (*bind.EncodedCall, error)
	ExecuteOwnershipTransferToMcmsWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	McmsRegisterUpgradeCap(upgradeCap bind.Object, registry bind.Object, state bind.Object) (*bind.EncodedCall, error)
	McmsRegisterUpgradeCapWithArgs(args ...any) (*bind.EncodedCall, error)
	McmsEntrypoint(typeArgs []string, state bind.Object, registry bind.Object, denyList bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsEntrypointWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
}

type ManagedTokenContract struct {
	*bind.BoundContract
	managedTokenEncoder
	devInspect *ManagedTokenDevInspect
}

type ManagedTokenDevInspect struct {
	contract *ManagedTokenContract
}

var _ IManagedToken = (*ManagedTokenContract)(nil)
var _ IManagedTokenDevInspect = (*ManagedTokenDevInspect)(nil)

func NewManagedToken(packageID string, client sui.ISuiAPI) (*ManagedTokenContract, error) {
	contract, err := bind.NewBoundContract(packageID, "managed_token", "managed_token", client)
	if err != nil {
		return nil, err
	}

	c := &ManagedTokenContract{
		BoundContract:       contract,
		managedTokenEncoder: managedTokenEncoder{BoundContract: contract},
	}
	c.devInspect = &ManagedTokenDevInspect{contract: c}
	return c, nil
}

func (c *ManagedTokenContract) Encoder() ManagedTokenEncoder {
	return c.managedTokenEncoder
}

func (c *ManagedTokenContract) DevInspect() IManagedTokenDevInspect {
	return c.devInspect
}

type TokenState struct {
	Id                string       `move:"sui::object::UID"`
	TreasuryCap       bind.Object  `move:"TreasuryCap<T>"`
	DenyCap           *bind.Object `move:"0x1::option::Option<DenyCapV2<T>>"`
	MintAllowancesMap bind.Object  `move:"VecMap<ID, MintAllowance<T>>"`
	OwnableState      bind.Object  `move:"OwnableState<T>"`
}

type MintCap struct {
	Id string `move:"sui::object::UID"`
}

type MintCapCreated struct {
	MintCap bind.Object `move:"ID"`
}

type MinterConfigured struct {
	MintCapOwner string      `move:"address"`
	MintCap      bind.Object `move:"ID"`
	Allowance    uint64      `move:"u64"`
	IsUnlimited  bool        `move:"bool"`
}

type Minted struct {
	MintCap bind.Object `move:"ID"`
	Minter  string      `move:"address"`
	To      string      `move:"address"`
	Amount  uint64      `move:"u64"`
}

type Burnt struct {
	MintCap bind.Object `move:"ID"`
	Burner  string      `move:"address"`
	From    string      `move:"address"`
	Amount  uint64      `move:"u64"`
}

type Blocklisted struct {
	Address string `move:"address"`
}

type Unblocklisted struct {
	Address string `move:"address"`
}

type Paused struct {
}

type Unpaused struct {
}

type MinterAllowanceIncremented struct {
	MintCap            bind.Object `move:"ID"`
	AllowanceIncrement uint64      `move:"u64"`
	NewAllowance       uint64      `move:"u64"`
}

type MinterUnlimitedAllowanceSet struct {
	MintCap bind.Object `move:"ID"`
}

type McmsCallback struct {
}

type bcsMinterConfigured struct {
	MintCapOwner [32]byte
	MintCap      bind.Object
	Allowance    uint64
	IsUnlimited  bool
}

func convertMinterConfiguredFromBCS(bcs bcsMinterConfigured) (MinterConfigured, error) {

	return MinterConfigured{
		MintCapOwner: fmt.Sprintf("0x%x", bcs.MintCapOwner),
		MintCap:      bcs.MintCap,
		Allowance:    bcs.Allowance,
		IsUnlimited:  bcs.IsUnlimited,
	}, nil
}

type bcsMinted struct {
	MintCap bind.Object
	Minter  [32]byte
	To      [32]byte
	Amount  uint64
}

func convertMintedFromBCS(bcs bcsMinted) (Minted, error) {

	return Minted{
		MintCap: bcs.MintCap,
		Minter:  fmt.Sprintf("0x%x", bcs.Minter),
		To:      fmt.Sprintf("0x%x", bcs.To),
		Amount:  bcs.Amount,
	}, nil
}

type bcsBurnt struct {
	MintCap bind.Object
	Burner  [32]byte
	From    [32]byte
	Amount  uint64
}

func convertBurntFromBCS(bcs bcsBurnt) (Burnt, error) {

	return Burnt{
		MintCap: bcs.MintCap,
		Burner:  fmt.Sprintf("0x%x", bcs.Burner),
		From:    fmt.Sprintf("0x%x", bcs.From),
		Amount:  bcs.Amount,
	}, nil
}

type bcsBlocklisted struct {
	Address [32]byte
}

func convertBlocklistedFromBCS(bcs bcsBlocklisted) (Blocklisted, error) {

	return Blocklisted{
		Address: fmt.Sprintf("0x%x", bcs.Address),
	}, nil
}

type bcsUnblocklisted struct {
	Address [32]byte
}

func convertUnblocklistedFromBCS(bcs bcsUnblocklisted) (Unblocklisted, error) {

	return Unblocklisted{
		Address: fmt.Sprintf("0x%x", bcs.Address),
	}, nil
}

func init() {
	bind.RegisterStructDecoder("managed_token::managed_token::TokenState", func(data []byte) (interface{}, error) {
		var result TokenState
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("managed_token::managed_token::MintCap", func(data []byte) (interface{}, error) {
		var result MintCap
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("managed_token::managed_token::MintCapCreated", func(data []byte) (interface{}, error) {
		var result MintCapCreated
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("managed_token::managed_token::MinterConfigured", func(data []byte) (interface{}, error) {
		var temp bcsMinterConfigured
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertMinterConfiguredFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("managed_token::managed_token::Minted", func(data []byte) (interface{}, error) {
		var temp bcsMinted
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertMintedFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("managed_token::managed_token::Burnt", func(data []byte) (interface{}, error) {
		var temp bcsBurnt
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertBurntFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("managed_token::managed_token::Blocklisted", func(data []byte) (interface{}, error) {
		var temp bcsBlocklisted
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertBlocklistedFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("managed_token::managed_token::Unblocklisted", func(data []byte) (interface{}, error) {
		var temp bcsUnblocklisted
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertUnblocklistedFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("managed_token::managed_token::Paused", func(data []byte) (interface{}, error) {
		var result Paused
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("managed_token::managed_token::Unpaused", func(data []byte) (interface{}, error) {
		var result Unpaused
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("managed_token::managed_token::MinterAllowanceIncremented", func(data []byte) (interface{}, error) {
		var result MinterAllowanceIncremented
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("managed_token::managed_token::MinterUnlimitedAllowanceSet", func(data []byte) (interface{}, error) {
		var result MinterUnlimitedAllowanceSet
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("managed_token::managed_token::McmsCallback", func(data []byte) (interface{}, error) {
		var result McmsCallback
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
}

// TypeAndVersion executes the type_and_version Move function.
func (c *ManagedTokenContract) TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenEncoder.TypeAndVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// Initialize executes the initialize Move function.
func (c *ManagedTokenContract) Initialize(ctx context.Context, opts *bind.CallOpts, typeArgs []string, treasuryCap bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenEncoder.Initialize(typeArgs, treasuryCap)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// InitializeWithDenyCap executes the initialize_with_deny_cap Move function.
func (c *ManagedTokenContract) InitializeWithDenyCap(ctx context.Context, opts *bind.CallOpts, typeArgs []string, treasuryCap bind.Object, denyCap bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenEncoder.InitializeWithDenyCap(typeArgs, treasuryCap, denyCap)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// MintAllowance executes the mint_allowance Move function.
func (c *ManagedTokenContract) MintAllowance(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, mintCap bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenEncoder.MintAllowance(typeArgs, state, mintCap)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// TotalSupply executes the total_supply Move function.
func (c *ManagedTokenContract) TotalSupply(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenEncoder.TotalSupply(typeArgs, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// IsAuthorizedMintCap executes the is_authorized_mint_cap Move function.
func (c *ManagedTokenContract) IsAuthorizedMintCap(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, id bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenEncoder.IsAuthorizedMintCap(typeArgs, state, id)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ConfigureNewMinter executes the configure_new_minter Move function.
func (c *ManagedTokenContract) ConfigureNewMinter(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object, minter string, allowance uint64, isUnlimited bool) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenEncoder.ConfigureNewMinter(typeArgs, state, ownerCap, minter, allowance, isUnlimited)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// IncrementMintAllowance executes the increment_mint_allowance Move function.
func (c *ManagedTokenContract) IncrementMintAllowance(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object, mintCapId bind.Object, denyList bind.Object, allowanceIncrement uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenEncoder.IncrementMintAllowance(typeArgs, state, ownerCap, mintCapId, denyList, allowanceIncrement)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// SetUnlimitedMintAllowances executes the set_unlimited_mint_allowances Move function.
func (c *ManagedTokenContract) SetUnlimitedMintAllowances(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object, mintCapId bind.Object, denyList bind.Object, isUnlimited bool) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenEncoder.SetUnlimitedMintAllowances(typeArgs, state, ownerCap, mintCapId, denyList, isUnlimited)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetAllMintCaps executes the get_all_mint_caps Move function.
func (c *ManagedTokenContract) GetAllMintCaps(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenEncoder.GetAllMintCaps(typeArgs, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// MintAndTransfer executes the mint_and_transfer Move function.
func (c *ManagedTokenContract) MintAndTransfer(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, mintCap bind.Object, denyList bind.Object, amount uint64, recipient string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenEncoder.MintAndTransfer(typeArgs, state, mintCap, denyList, amount, recipient)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// Mint executes the mint Move function.
func (c *ManagedTokenContract) Mint(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, mintCap bind.Object, denyList bind.Object, amount uint64, recipient string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenEncoder.Mint(typeArgs, state, mintCap, denyList, amount, recipient)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// Burn executes the burn Move function.
func (c *ManagedTokenContract) Burn(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, mintCap bind.Object, denyList bind.Object, coin bind.Object, from string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenEncoder.Burn(typeArgs, state, mintCap, denyList, coin, from)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// Blocklist executes the blocklist Move function.
func (c *ManagedTokenContract) Blocklist(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object, denyList bind.Object, addr string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenEncoder.Blocklist(typeArgs, state, ownerCap, denyList, addr)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// Unblocklist executes the unblocklist Move function.
func (c *ManagedTokenContract) Unblocklist(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object, denyList bind.Object, addr string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenEncoder.Unblocklist(typeArgs, state, ownerCap, denyList, addr)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// Pause executes the pause Move function.
func (c *ManagedTokenContract) Pause(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object, denyList bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenEncoder.Pause(typeArgs, state, ownerCap, denyList)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// Unpause executes the unpause Move function.
func (c *ManagedTokenContract) Unpause(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object, denyList bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenEncoder.Unpause(typeArgs, state, ownerCap, denyList)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// DestroyManagedToken executes the destroy_managed_token Move function.
func (c *ManagedTokenContract) DestroyManagedToken(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ownerCap bind.Object, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenEncoder.DestroyManagedToken(typeArgs, ownerCap, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// BorrowTreasuryCap executes the borrow_treasury_cap Move function.
func (c *ManagedTokenContract) BorrowTreasuryCap(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenEncoder.BorrowTreasuryCap(typeArgs, state, ownerCap)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// Owner executes the owner Move function.
func (c *ManagedTokenContract) Owner(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenEncoder.Owner(typeArgs, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// HasPendingTransfer executes the has_pending_transfer Move function.
func (c *ManagedTokenContract) HasPendingTransfer(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenEncoder.HasPendingTransfer(typeArgs, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// PendingTransferFrom executes the pending_transfer_from Move function.
func (c *ManagedTokenContract) PendingTransferFrom(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenEncoder.PendingTransferFrom(typeArgs, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// PendingTransferTo executes the pending_transfer_to Move function.
func (c *ManagedTokenContract) PendingTransferTo(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenEncoder.PendingTransferTo(typeArgs, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// PendingTransferAccepted executes the pending_transfer_accepted Move function.
func (c *ManagedTokenContract) PendingTransferAccepted(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenEncoder.PendingTransferAccepted(typeArgs, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// TransferOwnership executes the transfer_ownership Move function.
func (c *ManagedTokenContract) TransferOwnership(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object, newOwner string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenEncoder.TransferOwnership(typeArgs, state, ownerCap, newOwner)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// AcceptOwnership executes the accept_ownership Move function.
func (c *ManagedTokenContract) AcceptOwnership(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenEncoder.AcceptOwnership(typeArgs, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// AcceptOwnershipFromObject executes the accept_ownership_from_object Move function.
func (c *ManagedTokenContract) AcceptOwnershipFromObject(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, from string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenEncoder.AcceptOwnershipFromObject(typeArgs, state, from)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// AcceptOwnershipAsMcms executes the accept_ownership_as_mcms Move function.
func (c *ManagedTokenContract) AcceptOwnershipAsMcms(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenEncoder.AcceptOwnershipAsMcms(typeArgs, state, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ExecuteOwnershipTransfer executes the execute_ownership_transfer Move function.
func (c *ManagedTokenContract) ExecuteOwnershipTransfer(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ownerCap bind.Object, state bind.Object, to string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenEncoder.ExecuteOwnershipTransfer(typeArgs, ownerCap, state, to)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ExecuteOwnershipTransferToMcms executes the execute_ownership_transfer_to_mcms Move function.
func (c *ManagedTokenContract) ExecuteOwnershipTransferToMcms(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ownerCap bind.Object, state bind.Object, registry bind.Object, to string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenEncoder.ExecuteOwnershipTransferToMcms(typeArgs, ownerCap, state, registry, to)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsRegisterUpgradeCap executes the mcms_register_upgrade_cap Move function.
func (c *ManagedTokenContract) McmsRegisterUpgradeCap(ctx context.Context, opts *bind.CallOpts, upgradeCap bind.Object, registry bind.Object, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenEncoder.McmsRegisterUpgradeCap(upgradeCap, registry, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsEntrypoint executes the mcms_entrypoint Move function.
func (c *ManagedTokenContract) McmsEntrypoint(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, registry bind.Object, denyList bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenEncoder.McmsEntrypoint(typeArgs, state, registry, denyList, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// TypeAndVersion executes the type_and_version Move function using DevInspect to get return values.
//
// Returns: 0x1::string::String
func (d *ManagedTokenDevInspect) TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (string, error) {
	encoded, err := d.contract.managedTokenEncoder.TypeAndVersion()
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

// MintAllowance executes the mint_allowance Move function using DevInspect to get return values.
//
// Returns:
//
//	[0]: u64
//	[1]: bool
func (d *ManagedTokenDevInspect) MintAllowance(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, mintCap bind.Object) ([]any, error) {
	encoded, err := d.contract.managedTokenEncoder.MintAllowance(typeArgs, state, mintCap)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	return d.contract.Call(ctx, opts, encoded)
}

// TotalSupply executes the total_supply Move function using DevInspect to get return values.
//
// Returns: u64
func (d *ManagedTokenDevInspect) TotalSupply(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (uint64, error) {
	encoded, err := d.contract.managedTokenEncoder.TotalSupply(typeArgs, state)
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

// IsAuthorizedMintCap executes the is_authorized_mint_cap Move function using DevInspect to get return values.
//
// Returns: bool
func (d *ManagedTokenDevInspect) IsAuthorizedMintCap(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, id bind.Object) (bool, error) {
	encoded, err := d.contract.managedTokenEncoder.IsAuthorizedMintCap(typeArgs, state, id)
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

// GetAllMintCaps executes the get_all_mint_caps Move function using DevInspect to get return values.
//
// Returns: vector<ID>
func (d *ManagedTokenDevInspect) GetAllMintCaps(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) ([]bind.Object, error) {
	encoded, err := d.contract.managedTokenEncoder.GetAllMintCaps(typeArgs, state)
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
	result, ok := results[0].([]bind.Object)
	if !ok {
		return nil, fmt.Errorf("unexpected return type: expected []bind.Object, got %T", results[0])
	}
	return result, nil
}

// Mint executes the mint Move function using DevInspect to get return values.
//
// Returns: Coin<T>
func (d *ManagedTokenDevInspect) Mint(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, mintCap bind.Object, denyList bind.Object, amount uint64, recipient string) (any, error) {
	encoded, err := d.contract.managedTokenEncoder.Mint(typeArgs, state, mintCap, denyList, amount, recipient)
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

// DestroyManagedToken executes the destroy_managed_token Move function using DevInspect to get return values.
//
// Returns:
//
//	[0]: TreasuryCap<T>
//	[1]: 0x1::option::Option<DenyCapV2<T>>
func (d *ManagedTokenDevInspect) DestroyManagedToken(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ownerCap bind.Object, state bind.Object) ([]any, error) {
	encoded, err := d.contract.managedTokenEncoder.DestroyManagedToken(typeArgs, ownerCap, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	return d.contract.Call(ctx, opts, encoded)
}

// BorrowTreasuryCap executes the borrow_treasury_cap Move function using DevInspect to get return values.
//
// Returns: &TreasuryCap<T>
func (d *ManagedTokenDevInspect) BorrowTreasuryCap(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, ownerCap bind.Object) (any, error) {
	encoded, err := d.contract.managedTokenEncoder.BorrowTreasuryCap(typeArgs, state, ownerCap)
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
func (d *ManagedTokenDevInspect) Owner(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (string, error) {
	encoded, err := d.contract.managedTokenEncoder.Owner(typeArgs, state)
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
func (d *ManagedTokenDevInspect) HasPendingTransfer(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (bool, error) {
	encoded, err := d.contract.managedTokenEncoder.HasPendingTransfer(typeArgs, state)
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
func (d *ManagedTokenDevInspect) PendingTransferFrom(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*string, error) {
	encoded, err := d.contract.managedTokenEncoder.PendingTransferFrom(typeArgs, state)
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
func (d *ManagedTokenDevInspect) PendingTransferTo(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*string, error) {
	encoded, err := d.contract.managedTokenEncoder.PendingTransferTo(typeArgs, state)
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
func (d *ManagedTokenDevInspect) PendingTransferAccepted(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (*bool, error) {
	encoded, err := d.contract.managedTokenEncoder.PendingTransferAccepted(typeArgs, state)
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

type managedTokenEncoder struct {
	*bind.BoundContract
}

// TypeAndVersion encodes a call to the type_and_version Move function.
func (c managedTokenEncoder) TypeAndVersion() (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("type_and_version", typeArgsList, typeParamsList, []string{}, []any{}, []string{
		"0x1::string::String",
	})
}

// TypeAndVersionWithArgs encodes a call to the type_and_version Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenEncoder) TypeAndVersionWithArgs(args ...any) (*bind.EncodedCall, error) {
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
func (c managedTokenEncoder) Initialize(typeArgs []string, treasuryCap bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("initialize", typeArgsList, typeParamsList, []string{
		"TreasuryCap<T>",
	}, []any{
		treasuryCap,
	}, nil)
}

// InitializeWithArgs encodes a call to the initialize Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenEncoder) InitializeWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"TreasuryCap<T>",
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

// InitializeWithDenyCap encodes a call to the initialize_with_deny_cap Move function.
func (c managedTokenEncoder) InitializeWithDenyCap(typeArgs []string, treasuryCap bind.Object, denyCap bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("initialize_with_deny_cap", typeArgsList, typeParamsList, []string{
		"TreasuryCap<T>",
		"DenyCapV2<T>",
	}, []any{
		treasuryCap,
		denyCap,
	}, nil)
}

// InitializeWithDenyCapWithArgs encodes a call to the initialize_with_deny_cap Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenEncoder) InitializeWithDenyCapWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"TreasuryCap<T>",
		"DenyCapV2<T>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("initialize_with_deny_cap", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// MintAllowance encodes a call to the mint_allowance Move function.
func (c managedTokenEncoder) MintAllowance(typeArgs []string, state bind.Object, mintCap bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("mint_allowance", typeArgsList, typeParamsList, []string{
		"&TokenState<T>",
		"ID",
	}, []any{
		state,
		mintCap,
	}, []string{
		"u64",
		"bool",
	})
}

// MintAllowanceWithArgs encodes a call to the mint_allowance Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenEncoder) MintAllowanceWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&TokenState<T>",
		"ID",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("mint_allowance", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u64",
		"bool",
	})
}

// TotalSupply encodes a call to the total_supply Move function.
func (c managedTokenEncoder) TotalSupply(typeArgs []string, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("total_supply", typeArgsList, typeParamsList, []string{
		"&TokenState<T>",
	}, []any{
		state,
	}, []string{
		"u64",
	})
}

// TotalSupplyWithArgs encodes a call to the total_supply Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenEncoder) TotalSupplyWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&TokenState<T>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("total_supply", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u64",
	})
}

// IsAuthorizedMintCap encodes a call to the is_authorized_mint_cap Move function.
func (c managedTokenEncoder) IsAuthorizedMintCap(typeArgs []string, state bind.Object, id bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("is_authorized_mint_cap", typeArgsList, typeParamsList, []string{
		"&TokenState<T>",
		"ID",
	}, []any{
		state,
		id,
	}, []string{
		"bool",
	})
}

// IsAuthorizedMintCapWithArgs encodes a call to the is_authorized_mint_cap Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenEncoder) IsAuthorizedMintCapWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&TokenState<T>",
		"ID",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("is_authorized_mint_cap", typeArgsList, typeParamsList, expectedParams, args, []string{
		"bool",
	})
}

// ConfigureNewMinter encodes a call to the configure_new_minter Move function.
func (c managedTokenEncoder) ConfigureNewMinter(typeArgs []string, state bind.Object, ownerCap bind.Object, minter string, allowance uint64, isUnlimited bool) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("configure_new_minter", typeArgsList, typeParamsList, []string{
		"&mut TokenState<T>",
		"&OwnerCap<T>",
		"address",
		"u64",
		"bool",
	}, []any{
		state,
		ownerCap,
		minter,
		allowance,
		isUnlimited,
	}, nil)
}

// ConfigureNewMinterWithArgs encodes a call to the configure_new_minter Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenEncoder) ConfigureNewMinterWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut TokenState<T>",
		"&OwnerCap<T>",
		"address",
		"u64",
		"bool",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("configure_new_minter", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// IncrementMintAllowance encodes a call to the increment_mint_allowance Move function.
func (c managedTokenEncoder) IncrementMintAllowance(typeArgs []string, state bind.Object, ownerCap bind.Object, mintCapId bind.Object, denyList bind.Object, allowanceIncrement uint64) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("increment_mint_allowance", typeArgsList, typeParamsList, []string{
		"&mut TokenState<T>",
		"&OwnerCap<T>",
		"ID",
		"&DenyList",
		"u64",
	}, []any{
		state,
		ownerCap,
		mintCapId,
		denyList,
		allowanceIncrement,
	}, nil)
}

// IncrementMintAllowanceWithArgs encodes a call to the increment_mint_allowance Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenEncoder) IncrementMintAllowanceWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut TokenState<T>",
		"&OwnerCap<T>",
		"ID",
		"&DenyList",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("increment_mint_allowance", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// SetUnlimitedMintAllowances encodes a call to the set_unlimited_mint_allowances Move function.
func (c managedTokenEncoder) SetUnlimitedMintAllowances(typeArgs []string, state bind.Object, ownerCap bind.Object, mintCapId bind.Object, denyList bind.Object, isUnlimited bool) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("set_unlimited_mint_allowances", typeArgsList, typeParamsList, []string{
		"&mut TokenState<T>",
		"&OwnerCap<T>",
		"ID",
		"&DenyList",
		"bool",
	}, []any{
		state,
		ownerCap,
		mintCapId,
		denyList,
		isUnlimited,
	}, nil)
}

// SetUnlimitedMintAllowancesWithArgs encodes a call to the set_unlimited_mint_allowances Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenEncoder) SetUnlimitedMintAllowancesWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut TokenState<T>",
		"&OwnerCap<T>",
		"ID",
		"&DenyList",
		"bool",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("set_unlimited_mint_allowances", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// GetAllMintCaps encodes a call to the get_all_mint_caps Move function.
func (c managedTokenEncoder) GetAllMintCaps(typeArgs []string, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("get_all_mint_caps", typeArgsList, typeParamsList, []string{
		"&TokenState<T>",
	}, []any{
		state,
	}, []string{
		"vector<ID>",
	})
}

// GetAllMintCapsWithArgs encodes a call to the get_all_mint_caps Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenEncoder) GetAllMintCapsWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&TokenState<T>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("get_all_mint_caps", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<ID>",
	})
}

// MintAndTransfer encodes a call to the mint_and_transfer Move function.
func (c managedTokenEncoder) MintAndTransfer(typeArgs []string, state bind.Object, mintCap bind.Object, denyList bind.Object, amount uint64, recipient string) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("mint_and_transfer", typeArgsList, typeParamsList, []string{
		"&mut TokenState<T>",
		"&MintCap<T>",
		"&DenyList",
		"u64",
		"address",
	}, []any{
		state,
		mintCap,
		denyList,
		amount,
		recipient,
	}, nil)
}

// MintAndTransferWithArgs encodes a call to the mint_and_transfer Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenEncoder) MintAndTransferWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut TokenState<T>",
		"&MintCap<T>",
		"&DenyList",
		"u64",
		"address",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("mint_and_transfer", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// Mint encodes a call to the mint Move function.
func (c managedTokenEncoder) Mint(typeArgs []string, state bind.Object, mintCap bind.Object, denyList bind.Object, amount uint64, recipient string) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("mint", typeArgsList, typeParamsList, []string{
		"&mut TokenState<T>",
		"&MintCap<T>",
		"&DenyList",
		"u64",
		"address",
	}, []any{
		state,
		mintCap,
		denyList,
		amount,
		recipient,
	}, []string{
		"Coin<T>",
	})
}

// MintWithArgs encodes a call to the mint Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenEncoder) MintWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut TokenState<T>",
		"&MintCap<T>",
		"&DenyList",
		"u64",
		"address",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("mint", typeArgsList, typeParamsList, expectedParams, args, []string{
		"Coin<T>",
	})
}

// Burn encodes a call to the burn Move function.
func (c managedTokenEncoder) Burn(typeArgs []string, state bind.Object, mintCap bind.Object, denyList bind.Object, coin bind.Object, from string) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("burn", typeArgsList, typeParamsList, []string{
		"&mut TokenState<T>",
		"&MintCap<T>",
		"&DenyList",
		"Coin<T>",
		"address",
	}, []any{
		state,
		mintCap,
		denyList,
		coin,
		from,
	}, nil)
}

// BurnWithArgs encodes a call to the burn Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenEncoder) BurnWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut TokenState<T>",
		"&MintCap<T>",
		"&DenyList",
		"Coin<T>",
		"address",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("burn", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// Blocklist encodes a call to the blocklist Move function.
func (c managedTokenEncoder) Blocklist(typeArgs []string, state bind.Object, ownerCap bind.Object, denyList bind.Object, addr string) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("blocklist", typeArgsList, typeParamsList, []string{
		"&mut TokenState<T>",
		"&OwnerCap<T>",
		"&mut DenyList",
		"address",
	}, []any{
		state,
		ownerCap,
		denyList,
		addr,
	}, nil)
}

// BlocklistWithArgs encodes a call to the blocklist Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenEncoder) BlocklistWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut TokenState<T>",
		"&OwnerCap<T>",
		"&mut DenyList",
		"address",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("blocklist", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// Unblocklist encodes a call to the unblocklist Move function.
func (c managedTokenEncoder) Unblocklist(typeArgs []string, state bind.Object, ownerCap bind.Object, denyList bind.Object, addr string) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("unblocklist", typeArgsList, typeParamsList, []string{
		"&mut TokenState<T>",
		"&OwnerCap<T>",
		"&mut DenyList",
		"address",
	}, []any{
		state,
		ownerCap,
		denyList,
		addr,
	}, nil)
}

// UnblocklistWithArgs encodes a call to the unblocklist Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenEncoder) UnblocklistWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut TokenState<T>",
		"&OwnerCap<T>",
		"&mut DenyList",
		"address",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("unblocklist", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// Pause encodes a call to the pause Move function.
func (c managedTokenEncoder) Pause(typeArgs []string, state bind.Object, ownerCap bind.Object, denyList bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("pause", typeArgsList, typeParamsList, []string{
		"&mut TokenState<T>",
		"&OwnerCap<T>",
		"&mut DenyList",
	}, []any{
		state,
		ownerCap,
		denyList,
	}, nil)
}

// PauseWithArgs encodes a call to the pause Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenEncoder) PauseWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut TokenState<T>",
		"&OwnerCap<T>",
		"&mut DenyList",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("pause", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// Unpause encodes a call to the unpause Move function.
func (c managedTokenEncoder) Unpause(typeArgs []string, state bind.Object, ownerCap bind.Object, denyList bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("unpause", typeArgsList, typeParamsList, []string{
		"&mut TokenState<T>",
		"&OwnerCap<T>",
		"&mut DenyList",
	}, []any{
		state,
		ownerCap,
		denyList,
	}, nil)
}

// UnpauseWithArgs encodes a call to the unpause Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenEncoder) UnpauseWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut TokenState<T>",
		"&OwnerCap<T>",
		"&mut DenyList",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("unpause", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// DestroyManagedToken encodes a call to the destroy_managed_token Move function.
func (c managedTokenEncoder) DestroyManagedToken(typeArgs []string, ownerCap bind.Object, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("destroy_managed_token", typeArgsList, typeParamsList, []string{
		"OwnerCap<T>",
		"TokenState<T>",
	}, []any{
		ownerCap,
		state,
	}, []string{
		"TreasuryCap<T>",
		"0x1::option::Option<DenyCapV2<T>>",
	})
}

// DestroyManagedTokenWithArgs encodes a call to the destroy_managed_token Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenEncoder) DestroyManagedTokenWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"OwnerCap<T>",
		"TokenState<T>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("destroy_managed_token", typeArgsList, typeParamsList, expectedParams, args, []string{
		"TreasuryCap<T>",
		"0x1::option::Option<DenyCapV2<T>>",
	})
}

// BorrowTreasuryCap encodes a call to the borrow_treasury_cap Move function.
func (c managedTokenEncoder) BorrowTreasuryCap(typeArgs []string, state bind.Object, ownerCap bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("borrow_treasury_cap", typeArgsList, typeParamsList, []string{
		"&TokenState<T>",
		"&OwnerCap<T>",
	}, []any{
		state,
		ownerCap,
	}, []string{
		"&TreasuryCap<T>",
	})
}

// BorrowTreasuryCapWithArgs encodes a call to the borrow_treasury_cap Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenEncoder) BorrowTreasuryCapWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&TokenState<T>",
		"&OwnerCap<T>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("borrow_treasury_cap", typeArgsList, typeParamsList, expectedParams, args, []string{
		"&TreasuryCap<T>",
	})
}

// Owner encodes a call to the owner Move function.
func (c managedTokenEncoder) Owner(typeArgs []string, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("owner", typeArgsList, typeParamsList, []string{
		"&TokenState<T>",
	}, []any{
		state,
	}, []string{
		"address",
	})
}

// OwnerWithArgs encodes a call to the owner Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenEncoder) OwnerWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&TokenState<T>",
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
func (c managedTokenEncoder) HasPendingTransfer(typeArgs []string, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("has_pending_transfer", typeArgsList, typeParamsList, []string{
		"&TokenState<T>",
	}, []any{
		state,
	}, []string{
		"bool",
	})
}

// HasPendingTransferWithArgs encodes a call to the has_pending_transfer Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenEncoder) HasPendingTransferWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&TokenState<T>",
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
func (c managedTokenEncoder) PendingTransferFrom(typeArgs []string, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("pending_transfer_from", typeArgsList, typeParamsList, []string{
		"&TokenState<T>",
	}, []any{
		state,
	}, []string{
		"0x1::option::Option<address>",
	})
}

// PendingTransferFromWithArgs encodes a call to the pending_transfer_from Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenEncoder) PendingTransferFromWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&TokenState<T>",
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
func (c managedTokenEncoder) PendingTransferTo(typeArgs []string, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("pending_transfer_to", typeArgsList, typeParamsList, []string{
		"&TokenState<T>",
	}, []any{
		state,
	}, []string{
		"0x1::option::Option<address>",
	})
}

// PendingTransferToWithArgs encodes a call to the pending_transfer_to Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenEncoder) PendingTransferToWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&TokenState<T>",
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
func (c managedTokenEncoder) PendingTransferAccepted(typeArgs []string, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("pending_transfer_accepted", typeArgsList, typeParamsList, []string{
		"&TokenState<T>",
	}, []any{
		state,
	}, []string{
		"0x1::option::Option<bool>",
	})
}

// PendingTransferAcceptedWithArgs encodes a call to the pending_transfer_accepted Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenEncoder) PendingTransferAcceptedWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&TokenState<T>",
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
func (c managedTokenEncoder) TransferOwnership(typeArgs []string, state bind.Object, ownerCap bind.Object, newOwner string) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("transfer_ownership", typeArgsList, typeParamsList, []string{
		"&mut TokenState<T>",
		"&OwnerCap<T>",
		"address",
	}, []any{
		state,
		ownerCap,
		newOwner,
	}, nil)
}

// TransferOwnershipWithArgs encodes a call to the transfer_ownership Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenEncoder) TransferOwnershipWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut TokenState<T>",
		"&OwnerCap<T>",
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
func (c managedTokenEncoder) AcceptOwnership(typeArgs []string, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("accept_ownership", typeArgsList, typeParamsList, []string{
		"&mut TokenState<T>",
	}, []any{
		state,
	}, nil)
}

// AcceptOwnershipWithArgs encodes a call to the accept_ownership Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenEncoder) AcceptOwnershipWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut TokenState<T>",
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
func (c managedTokenEncoder) AcceptOwnershipFromObject(typeArgs []string, state bind.Object, from string) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("accept_ownership_from_object", typeArgsList, typeParamsList, []string{
		"&mut TokenState<T>",
		"&mut UID",
	}, []any{
		state,
		from,
	}, nil)
}

// AcceptOwnershipFromObjectWithArgs encodes a call to the accept_ownership_from_object Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenEncoder) AcceptOwnershipFromObjectWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut TokenState<T>",
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

// AcceptOwnershipAsMcms encodes a call to the accept_ownership_as_mcms Move function.
func (c managedTokenEncoder) AcceptOwnershipAsMcms(typeArgs []string, state bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("accept_ownership_as_mcms", typeArgsList, typeParamsList, []string{
		"&mut TokenState<T>",
		"ExecutingCallbackParams",
	}, []any{
		state,
		params,
	}, nil)
}

// AcceptOwnershipAsMcmsWithArgs encodes a call to the accept_ownership_as_mcms Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenEncoder) AcceptOwnershipAsMcmsWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut TokenState<T>",
		"ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("accept_ownership_as_mcms", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// ExecuteOwnershipTransfer encodes a call to the execute_ownership_transfer Move function.
func (c managedTokenEncoder) ExecuteOwnershipTransfer(typeArgs []string, ownerCap bind.Object, state bind.Object, to string) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("execute_ownership_transfer", typeArgsList, typeParamsList, []string{
		"OwnerCap<T>",
		"&mut TokenState<T>",
		"address",
	}, []any{
		ownerCap,
		state,
		to,
	}, nil)
}

// ExecuteOwnershipTransferWithArgs encodes a call to the execute_ownership_transfer Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenEncoder) ExecuteOwnershipTransferWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"OwnerCap<T>",
		"&mut TokenState<T>",
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
func (c managedTokenEncoder) ExecuteOwnershipTransferToMcms(typeArgs []string, ownerCap bind.Object, state bind.Object, registry bind.Object, to string) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("execute_ownership_transfer_to_mcms", typeArgsList, typeParamsList, []string{
		"OwnerCap<T>",
		"&mut TokenState<T>",
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
func (c managedTokenEncoder) ExecuteOwnershipTransferToMcmsWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"OwnerCap<T>",
		"&mut TokenState<T>",
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
func (c managedTokenEncoder) McmsRegisterUpgradeCap(upgradeCap bind.Object, registry bind.Object, state bind.Object) (*bind.EncodedCall, error) {
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
func (c managedTokenEncoder) McmsRegisterUpgradeCapWithArgs(args ...any) (*bind.EncodedCall, error) {
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
func (c managedTokenEncoder) McmsEntrypoint(typeArgs []string, state bind.Object, registry bind.Object, denyList bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("mcms_entrypoint", typeArgsList, typeParamsList, []string{
		"&mut TokenState<T>",
		"&mut Registry",
		"&mut DenyList",
		"ExecutingCallbackParams",
	}, []any{
		state,
		registry,
		denyList,
		params,
	}, nil)
}

// McmsEntrypointWithArgs encodes a call to the mcms_entrypoint Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenEncoder) McmsEntrypointWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut TokenState<T>",
		"&mut Registry",
		"&mut DenyList",
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
