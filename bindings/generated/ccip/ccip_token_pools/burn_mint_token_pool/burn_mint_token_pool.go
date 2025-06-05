// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package module_burn_mint_token_pool

import (
	"context"
	"fmt"
	"math/big"

	"github.com/holiman/uint256"
	"github.com/pattonkan/sui-go/sui"
	"github.com/pattonkan/sui-go/sui/suiptb"
	"github.com/pattonkan/sui-go/suiclient"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_common "github.com/smartcontractkit/chainlink-sui/bindings/common"
)

// Unused vars used for unused imports
var (
	_ = big.NewInt
	_ = uint256.NewInt
)

type IBurnMintTokenPool interface {
	TypeAndVersion() bind.IMethod
	Initialize(ref module_common.CCIPObjectRef, coinMetadata bind.Object, treasuryCap bind.Object, tokenPoolPackageId string) bind.IMethod
	GetRouter() bind.IMethod
	// Connect adds/changes the client used in the contract
	Connect(client suiclient.ClientImpl)
}

type BurnMintTokenPoolContract struct {
	packageID *sui.Address
	client    suiclient.ClientImpl
}

var _ IBurnMintTokenPool = (*BurnMintTokenPoolContract)(nil)

func NewBurnMintTokenPool(packageID string, client suiclient.ClientImpl) (*BurnMintTokenPoolContract, error) {
	pkgObjectId, err := bind.ToSuiAddress(packageID)
	if err != nil {
		return nil, fmt.Errorf("package ID is not a Sui address: %w", err)
	}

	return &BurnMintTokenPoolContract{
		packageID: pkgObjectId,
		client:    client,
	}, nil
}

func (c *BurnMintTokenPoolContract) Connect(client suiclient.ClientImpl) {
	c.client = client
}

// Structs

type OwnerCap struct {
	Id string `move:"sui::object::UID"`
}

type BurnMintTokenPoolState struct {
	Id          string      `move:"sui::object::UID"`
	TreasuryCap bind.Object `move:"TreasuryCap<T>"`
}

type TypeProof struct {
}

// Functions

func (c *BurnMintTokenPoolContract) TypeAndVersion() bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "burn_mint_token_pool", "type_and_version", false, "")
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "burn_mint_token_pool", "type_and_version", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *BurnMintTokenPoolContract) Initialize(ref module_common.CCIPObjectRef, coinMetadata bind.Object, treasuryCap bind.Object, tokenPoolPackageId string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "burn_mint_token_pool", "initialize", false, "", ref, coinMetadata, treasuryCap, tokenPoolPackageId)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "burn_mint_token_pool", "initialize", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *BurnMintTokenPoolContract) GetRouter() bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "burn_mint_token_pool", "get_router", false, "")
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "burn_mint_token_pool", "get_router", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}
