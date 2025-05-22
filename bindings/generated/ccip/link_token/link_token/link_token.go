// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package module_link_token

import (
	"context"
	"fmt"
	"math/big"

	"github.com/pattonkan/sui-go/sui"
	"github.com/pattonkan/sui-go/sui/suiptb"
	"github.com/pattonkan/sui-go/suiclient"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
)

// Unused vars used for unused imports
var (
	_ = big.NewInt
)

type ILinkToken interface {
	MintAndTransfer(treasuryCap string, amount uint64, recipient string) bind.IMethod
	Mint(treasuryCap string, amount uint64) bind.IMethod
	// Connect adds/changes the client used in the contract
	Connect(client suiclient.ClientImpl)
}

type LinkTokenContract struct {
	packageID *sui.Address
	client    suiclient.ClientImpl
}

var _ ILinkToken = (*LinkTokenContract)(nil)

func NewLinkToken(packageID string, client suiclient.ClientImpl) (*LinkTokenContract, error) {
	pkgObjectId, err := bind.ToSuiAddress(packageID)
	if err != nil {
		return nil, fmt.Errorf("package ID is not a Sui address: %w", err)
	}

	return &LinkTokenContract{
		packageID: pkgObjectId,
		client:    client,
	}, nil
}

func (c *LinkTokenContract) Connect(client suiclient.ClientImpl) {
	c.client = client
}

// Structs

type LINK_TOKEN struct {
}

// Functions

func (c *LinkTokenContract) MintAndTransfer(treasuryCap string, amount uint64, recipient string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "link_token", "mint_and_transfer", false, "", treasuryCap, amount, recipient)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "link_token", "mint_and_transfer", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *LinkTokenContract) Mint(treasuryCap string, amount uint64) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "link_token", "mint", false, "", treasuryCap, amount)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "link_token", "mint", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}
