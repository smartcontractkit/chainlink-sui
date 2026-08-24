// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package module_ccip_lock_release_token

import (
	"context"
	"fmt"
	"math/big"

	"github.com/block-vision/sui-go-sdk/models"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
)

var (
	_ = big.NewInt
)

const FunctionInfo = `[{"package":"ccip_lock_release_token","module":"ccip_lock_release_token","name":"mint","parameters":[{"name":"treasury_cap","type":"TreasuryCap<CCIP_LOCK_RELEASE_TOKEN>"},{"name":"amount","type":"u64"}]},{"package":"ccip_lock_release_token","module":"ccip_lock_release_token","name":"mint_and_transfer","parameters":[{"name":"treasury_cap","type":"TreasuryCap<CCIP_LOCK_RELEASE_TOKEN>"},{"name":"amount","type":"u64"},{"name":"recipient","type":"address"}]}]`

type ICcipLockReleaseToken interface {
	MintAndTransfer(ctx context.Context, opts *bind.CallOpts, treasuryCap bind.Object, amount uint64, recipient string) (*models.SuiTransactionBlockResponse, error)
	Mint(ctx context.Context, opts *bind.CallOpts, treasuryCap bind.Object, amount uint64) (*models.SuiTransactionBlockResponse, error)
	DevInspect() ICcipLockReleaseTokenDevInspect
	Encoder() CcipLockReleaseTokenEncoder
	Bound() bind.IBoundContract
}

type ICcipLockReleaseTokenDevInspect interface {
}

type CcipLockReleaseTokenEncoder interface {
	MintAndTransfer(treasuryCap bind.Object, amount uint64, recipient string) (*bind.EncodedCall, error)
	MintAndTransferWithArgs(args ...any) (*bind.EncodedCall, error)
	Mint(treasuryCap bind.Object, amount uint64) (*bind.EncodedCall, error)
	MintWithArgs(args ...any) (*bind.EncodedCall, error)
}

type CcipLockReleaseTokenContract struct {
	*bind.BoundContract
	ccipLockReleaseTokenEncoder
	devInspect *CcipLockReleaseTokenDevInspect
}

type CcipLockReleaseTokenDevInspect struct {
	contract *CcipLockReleaseTokenContract
}

var _ ICcipLockReleaseToken = (*CcipLockReleaseTokenContract)(nil)
var _ ICcipLockReleaseTokenDevInspect = (*CcipLockReleaseTokenDevInspect)(nil)

func NewCcipLockReleaseToken(packageID string, chainClient client.BindingsClient) (ICcipLockReleaseToken, error) {
	contract, err := bind.NewBoundContract(packageID, "ccip_lock_release_token", "ccip_lock_release_token", chainClient)
	if err != nil {
		return nil, err
	}

	c := &CcipLockReleaseTokenContract{
		BoundContract:               contract,
		ccipLockReleaseTokenEncoder: ccipLockReleaseTokenEncoder{BoundContract: contract},
	}
	c.devInspect = &CcipLockReleaseTokenDevInspect{contract: c}
	return c, nil
}

func (c *CcipLockReleaseTokenContract) Bound() bind.IBoundContract {
	return c.BoundContract
}

func (c *CcipLockReleaseTokenContract) Encoder() CcipLockReleaseTokenEncoder {
	return c.ccipLockReleaseTokenEncoder
}

func (c *CcipLockReleaseTokenContract) DevInspect() ICcipLockReleaseTokenDevInspect {
	return c.devInspect
}

type CCIP_LOCK_RELEASE_TOKEN struct {
}

// MintAndTransfer executes the mint_and_transfer Move function.
func (c *CcipLockReleaseTokenContract) MintAndTransfer(ctx context.Context, opts *bind.CallOpts, treasuryCap bind.Object, amount uint64, recipient string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.ccipLockReleaseTokenEncoder.MintAndTransfer(treasuryCap, amount, recipient)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// Mint executes the mint Move function.
func (c *CcipLockReleaseTokenContract) Mint(ctx context.Context, opts *bind.CallOpts, treasuryCap bind.Object, amount uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.ccipLockReleaseTokenEncoder.Mint(treasuryCap, amount)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

type ccipLockReleaseTokenEncoder struct {
	*bind.BoundContract
}

// MintAndTransfer encodes a call to the mint_and_transfer Move function.
func (c ccipLockReleaseTokenEncoder) MintAndTransfer(treasuryCap bind.Object, amount uint64, recipient string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mint_and_transfer", typeArgsList, typeParamsList, []string{
		"&mut TreasuryCap<CCIP_LOCK_RELEASE_TOKEN>",
		"u64",
		"address",
	}, []any{
		treasuryCap,
		amount,
		recipient,
	}, nil)
}

// MintAndTransferWithArgs encodes a call to the mint_and_transfer Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c ccipLockReleaseTokenEncoder) MintAndTransferWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut TreasuryCap<CCIP_LOCK_RELEASE_TOKEN>",
		"u64",
		"address",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mint_and_transfer", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// Mint encodes a call to the mint Move function.
func (c ccipLockReleaseTokenEncoder) Mint(treasuryCap bind.Object, amount uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mint", typeArgsList, typeParamsList, []string{
		"&mut TreasuryCap<CCIP_LOCK_RELEASE_TOKEN>",
		"u64",
	}, []any{
		treasuryCap,
		amount,
	}, nil)
}

// MintWithArgs encodes a call to the mint Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c ccipLockReleaseTokenEncoder) MintWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut TreasuryCap<CCIP_LOCK_RELEASE_TOKEN>",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mint", typeArgsList, typeParamsList, expectedParams, args, nil)
}
