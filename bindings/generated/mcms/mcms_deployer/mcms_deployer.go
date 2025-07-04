// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package module_mcms_deployer

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

type IMcmsDeployer interface {
	AuthorizeUpgrade(param module_common.OwnerCap, state bind.Object, policy byte, digest []byte, packageAddress string) bind.IMethod
	// Connect adds/changes the client used in the contract
	Connect(client suiclient.ClientImpl)
}

type McmsDeployerContract struct {
	packageID *sui.Address
	client    suiclient.ClientImpl
}

var _ IMcmsDeployer = (*McmsDeployerContract)(nil)

func NewMcmsDeployer(packageID string, client suiclient.ClientImpl) (*McmsDeployerContract, error) {
	pkgObjectId, err := bind.ToSuiAddress(packageID)
	if err != nil {
		return nil, fmt.Errorf("package ID is not a Sui address: %w", err)
	}

	return &McmsDeployerContract{
		packageID: pkgObjectId,
		client:    client,
	}, nil
}

func (c *McmsDeployerContract) Connect(client suiclient.ClientImpl) {
	c.client = client
}

// Structs

type DeployerState struct {
	Id string `move:"sui::object::UID"`
}

type UpgradeCapRegistered struct {
	PrevOwner      string `move:"address"`
	PackageAddress string `move:"address"`
	Version        uint64 `move:"u64"`
	Policy         byte   `move:"u8"`
}

type UpgradeTicketAuthorized struct {
	PackageAddress string `move:"address"`
	Policy         byte   `move:"u8"`
	Digest         []byte `move:"vector<u8>"`
}

type UpgradeReceiptCommitted struct {
	PackageAddress string `move:"address"`
	OldVersion     uint64 `move:"u64"`
	NewVersion     uint64 `move:"u64"`
}

type MCMS_DEPLOYER struct {
}

// Functions

func (c *McmsDeployerContract) AuthorizeUpgrade(param module_common.OwnerCap, state bind.Object, policy byte, digest []byte, packageAddress string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "mcms_deployer", "authorize_upgrade", false, "", "", param, state, policy, digest, packageAddress)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "mcms_deployer", "authorize_upgrade", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}
