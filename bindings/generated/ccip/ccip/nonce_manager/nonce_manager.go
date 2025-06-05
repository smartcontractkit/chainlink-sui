// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package module_nonce_manager

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

type INonceManager interface {
	TypeAndVersion() bind.IMethod
	Initialize(ref module_common.CCIPObjectRef, param module_common.OwnerCap) bind.IMethod
	GetOutboundNonce(ref module_common.CCIPObjectRef, destChainSelector uint64, sender string) bind.IMethod
	GetIncrementedOutboundNonce(ref module_common.CCIPObjectRef, param bind.Object, destChainSelector uint64, sender string) bind.IMethod
	// Connect adds/changes the client used in the contract
	Connect(client suiclient.ClientImpl)
}

type NonceManagerContract struct {
	packageID *sui.Address
	client    suiclient.ClientImpl
}

var _ INonceManager = (*NonceManagerContract)(nil)

func NewNonceManager(packageID string, client suiclient.ClientImpl) (*NonceManagerContract, error) {
	pkgObjectId, err := bind.ToSuiAddress(packageID)
	if err != nil {
		return nil, fmt.Errorf("package ID is not a Sui address: %w", err)
	}

	return &NonceManagerContract{
		packageID: pkgObjectId,
		client:    client,
	}, nil
}

func (c *NonceManagerContract) Connect(client suiclient.ClientImpl) {
	c.client = client
}

// Structs

type NonceManagerCap struct {
	Id string `move:"sui::object::UID"`
}

type NonceManagerState struct {
	Id string `move:"sui::object::UID"`
}

// Functions

func (c *NonceManagerContract) TypeAndVersion() bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "nonce_manager", "type_and_version", false, "")
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "nonce_manager", "type_and_version", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *NonceManagerContract) Initialize(ref module_common.CCIPObjectRef, param module_common.OwnerCap) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "nonce_manager", "initialize", false, "", ref, param)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "nonce_manager", "initialize", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *NonceManagerContract) GetOutboundNonce(ref module_common.CCIPObjectRef, destChainSelector uint64, sender string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "nonce_manager", "get_outbound_nonce", false, "", ref, destChainSelector, sender)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "nonce_manager", "get_outbound_nonce", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *NonceManagerContract) GetIncrementedOutboundNonce(ref module_common.CCIPObjectRef, param bind.Object, destChainSelector uint64, sender string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "nonce_manager", "get_incremented_outbound_nonce", false, "", ref, param, destChainSelector, sender)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "nonce_manager", "get_incremented_outbound_nonce", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}
