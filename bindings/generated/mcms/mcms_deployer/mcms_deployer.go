// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package module_mcms_deployer

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

type IMcmsDeployer interface {
	RegisterUpgradeCap(ctx context.Context, opts *bind.CallOpts, state bind.Object, registry bind.Object, upgradeCap bind.Object) (*models.SuiTransactionBlockResponse, error)
	AuthorizeUpgrade(ctx context.Context, opts *bind.CallOpts, param bind.Object, state bind.Object, policy byte, digest []byte, packageAddress string) (*models.SuiTransactionBlockResponse, error)
	AuthorizeUpgradeBypassCap(ctx context.Context, opts *bind.CallOpts, state bind.Object, policy byte, digest []byte, packageAddress string) (*models.SuiTransactionBlockResponse, error)
	CommitUpgrade(ctx context.Context, opts *bind.CallOpts, state bind.Object, receipt bind.Object) (*models.SuiTransactionBlockResponse, error)
	DevInspect() IMcmsDeployerDevInspect
	Encoder() McmsDeployerEncoder
	Bound() bind.IBoundContract
}

type IMcmsDeployerDevInspect interface {
	AuthorizeUpgrade(ctx context.Context, opts *bind.CallOpts, param bind.Object, state bind.Object, policy byte, digest []byte, packageAddress string) (bind.Object, error)
	AuthorizeUpgradeBypassCap(ctx context.Context, opts *bind.CallOpts, state bind.Object, policy byte, digest []byte, packageAddress string) (bind.Object, error)
}

type McmsDeployerEncoder interface {
	RegisterUpgradeCap(state bind.Object, registry bind.Object, upgradeCap bind.Object) (*bind.EncodedCall, error)
	RegisterUpgradeCapWithArgs(args ...any) (*bind.EncodedCall, error)
	AuthorizeUpgrade(param bind.Object, state bind.Object, policy byte, digest []byte, packageAddress string) (*bind.EncodedCall, error)
	AuthorizeUpgradeWithArgs(args ...any) (*bind.EncodedCall, error)
	AuthorizeUpgradeBypassCap(state bind.Object, policy byte, digest []byte, packageAddress string) (*bind.EncodedCall, error)
	AuthorizeUpgradeBypassCapWithArgs(args ...any) (*bind.EncodedCall, error)
	CommitUpgrade(state bind.Object, receipt bind.Object) (*bind.EncodedCall, error)
	CommitUpgradeWithArgs(args ...any) (*bind.EncodedCall, error)
}

type McmsDeployerContract struct {
	*bind.BoundContract
	mcmsDeployerEncoder
	devInspect *McmsDeployerDevInspect
}

type McmsDeployerDevInspect struct {
	contract *McmsDeployerContract
}

var _ IMcmsDeployer = (*McmsDeployerContract)(nil)
var _ IMcmsDeployerDevInspect = (*McmsDeployerDevInspect)(nil)

func NewMcmsDeployer(packageID string, client sui.ISuiAPI) (IMcmsDeployer, error) {
	contract, err := bind.NewBoundContract(packageID, "mcms", "mcms_deployer", client)
	if err != nil {
		return nil, err
	}

	c := &McmsDeployerContract{
		BoundContract:       contract,
		mcmsDeployerEncoder: mcmsDeployerEncoder{BoundContract: contract},
	}
	c.devInspect = &McmsDeployerDevInspect{contract: c}
	return c, nil
}

func (c *McmsDeployerContract) Bound() bind.IBoundContract {
	return c.BoundContract
}

func (c *McmsDeployerContract) Encoder() McmsDeployerEncoder {
	return c.mcmsDeployerEncoder
}

func (c *McmsDeployerContract) DevInspect() IMcmsDeployerDevInspect {
	return c.devInspect
}

type DeployerState struct {
	Id           string      `move:"sui::object::UID"`
	UpgradeCaps  bind.Object `move:"Table<address, UpgradeCap>"`
	CapToPackage bind.Object `move:"Table<ID, address>"`
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
	OldPackageAddress string `move:"address"`
	NewPackageAddress string `move:"address"`
	OldVersion        uint64 `move:"u64"`
	NewVersion        uint64 `move:"u64"`
}

type MCMS_DEPLOYER struct {
}

type bcsUpgradeCapRegistered struct {
	PrevOwner      [32]byte
	PackageAddress [32]byte
	Version        uint64
	Policy         byte
}

func convertUpgradeCapRegisteredFromBCS(bcs bcsUpgradeCapRegistered) (UpgradeCapRegistered, error) {

	return UpgradeCapRegistered{
		PrevOwner:      fmt.Sprintf("0x%x", bcs.PrevOwner),
		PackageAddress: fmt.Sprintf("0x%x", bcs.PackageAddress),
		Version:        bcs.Version,
		Policy:         bcs.Policy,
	}, nil
}

type bcsUpgradeTicketAuthorized struct {
	PackageAddress [32]byte
	Policy         byte
	Digest         []byte
}

func convertUpgradeTicketAuthorizedFromBCS(bcs bcsUpgradeTicketAuthorized) (UpgradeTicketAuthorized, error) {

	return UpgradeTicketAuthorized{
		PackageAddress: fmt.Sprintf("0x%x", bcs.PackageAddress),
		Policy:         bcs.Policy,
		Digest:         bcs.Digest,
	}, nil
}

type bcsUpgradeReceiptCommitted struct {
	OldPackageAddress [32]byte
	NewPackageAddress [32]byte
	OldVersion        uint64
	NewVersion        uint64
}

func convertUpgradeReceiptCommittedFromBCS(bcs bcsUpgradeReceiptCommitted) (UpgradeReceiptCommitted, error) {

	return UpgradeReceiptCommitted{
		OldPackageAddress: fmt.Sprintf("0x%x", bcs.OldPackageAddress),
		NewPackageAddress: fmt.Sprintf("0x%x", bcs.NewPackageAddress),
		OldVersion:        bcs.OldVersion,
		NewVersion:        bcs.NewVersion,
	}, nil
}

func init() {
	bind.RegisterStructDecoder("mcms::mcms_deployer::DeployerState", func(data []byte) (interface{}, error) {
		var result DeployerState
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("mcms::mcms_deployer::UpgradeCapRegistered", func(data []byte) (interface{}, error) {
		var temp bcsUpgradeCapRegistered
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertUpgradeCapRegisteredFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("mcms::mcms_deployer::UpgradeTicketAuthorized", func(data []byte) (interface{}, error) {
		var temp bcsUpgradeTicketAuthorized
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertUpgradeTicketAuthorizedFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("mcms::mcms_deployer::UpgradeReceiptCommitted", func(data []byte) (interface{}, error) {
		var temp bcsUpgradeReceiptCommitted
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertUpgradeReceiptCommittedFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("mcms::mcms_deployer::MCMS_DEPLOYER", func(data []byte) (interface{}, error) {
		var result MCMS_DEPLOYER
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
}

// RegisterUpgradeCap executes the register_upgrade_cap Move function.
func (c *McmsDeployerContract) RegisterUpgradeCap(ctx context.Context, opts *bind.CallOpts, state bind.Object, registry bind.Object, upgradeCap bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsDeployerEncoder.RegisterUpgradeCap(state, registry, upgradeCap)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// AuthorizeUpgrade executes the authorize_upgrade Move function.
func (c *McmsDeployerContract) AuthorizeUpgrade(ctx context.Context, opts *bind.CallOpts, param bind.Object, state bind.Object, policy byte, digest []byte, packageAddress string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsDeployerEncoder.AuthorizeUpgrade(param, state, policy, digest, packageAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// AuthorizeUpgradeBypassCap executes the authorize_upgrade_bypass_cap Move function.
func (c *McmsDeployerContract) AuthorizeUpgradeBypassCap(ctx context.Context, opts *bind.CallOpts, state bind.Object, policy byte, digest []byte, packageAddress string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsDeployerEncoder.AuthorizeUpgradeBypassCap(state, policy, digest, packageAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// CommitUpgrade executes the commit_upgrade Move function.
func (c *McmsDeployerContract) CommitUpgrade(ctx context.Context, opts *bind.CallOpts, state bind.Object, receipt bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsDeployerEncoder.CommitUpgrade(state, receipt)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// AuthorizeUpgrade executes the authorize_upgrade Move function using DevInspect to get return values.
//
// Returns: UpgradeTicket
func (d *McmsDeployerDevInspect) AuthorizeUpgrade(ctx context.Context, opts *bind.CallOpts, param bind.Object, state bind.Object, policy byte, digest []byte, packageAddress string) (bind.Object, error) {
	encoded, err := d.contract.mcmsDeployerEncoder.AuthorizeUpgrade(param, state, policy, digest, packageAddress)
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

// AuthorizeUpgradeBypassCap executes the authorize_upgrade_bypass_cap Move function using DevInspect to get return values.
//
// Returns: UpgradeTicket
func (d *McmsDeployerDevInspect) AuthorizeUpgradeBypassCap(ctx context.Context, opts *bind.CallOpts, state bind.Object, policy byte, digest []byte, packageAddress string) (bind.Object, error) {
	encoded, err := d.contract.mcmsDeployerEncoder.AuthorizeUpgradeBypassCap(state, policy, digest, packageAddress)
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

type mcmsDeployerEncoder struct {
	*bind.BoundContract
}

// RegisterUpgradeCap encodes a call to the register_upgrade_cap Move function.
func (c mcmsDeployerEncoder) RegisterUpgradeCap(state bind.Object, registry bind.Object, upgradeCap bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("register_upgrade_cap", typeArgsList, typeParamsList, []string{
		"&mut DeployerState",
		"&Registry",
		"UpgradeCap",
	}, []any{
		state,
		registry,
		upgradeCap,
	}, nil)
}

// RegisterUpgradeCapWithArgs encodes a call to the register_upgrade_cap Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsDeployerEncoder) RegisterUpgradeCapWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut DeployerState",
		"&Registry",
		"UpgradeCap",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("register_upgrade_cap", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// AuthorizeUpgrade encodes a call to the authorize_upgrade Move function.
func (c mcmsDeployerEncoder) AuthorizeUpgrade(param bind.Object, state bind.Object, policy byte, digest []byte, packageAddress string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("authorize_upgrade", typeArgsList, typeParamsList, []string{
		"&OwnerCap",
		"&mut DeployerState",
		"u8",
		"vector<u8>",
		"address",
	}, []any{
		param,
		state,
		policy,
		digest,
		packageAddress,
	}, []string{
		"UpgradeTicket",
	})
}

// AuthorizeUpgradeWithArgs encodes a call to the authorize_upgrade Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsDeployerEncoder) AuthorizeUpgradeWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&OwnerCap",
		"&mut DeployerState",
		"u8",
		"vector<u8>",
		"address",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("authorize_upgrade", typeArgsList, typeParamsList, expectedParams, args, []string{
		"UpgradeTicket",
	})
}

// AuthorizeUpgradeBypassCap encodes a call to the authorize_upgrade_bypass_cap Move function.
func (c mcmsDeployerEncoder) AuthorizeUpgradeBypassCap(state bind.Object, policy byte, digest []byte, packageAddress string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("authorize_upgrade_bypass_cap", typeArgsList, typeParamsList, []string{
		"&mut DeployerState",
		"u8",
		"vector<u8>",
		"address",
	}, []any{
		state,
		policy,
		digest,
		packageAddress,
	}, []string{
		"UpgradeTicket",
	})
}

// AuthorizeUpgradeBypassCapWithArgs encodes a call to the authorize_upgrade_bypass_cap Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsDeployerEncoder) AuthorizeUpgradeBypassCapWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut DeployerState",
		"u8",
		"vector<u8>",
		"address",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("authorize_upgrade_bypass_cap", typeArgsList, typeParamsList, expectedParams, args, []string{
		"UpgradeTicket",
	})
}

// CommitUpgrade encodes a call to the commit_upgrade Move function.
func (c mcmsDeployerEncoder) CommitUpgrade(state bind.Object, receipt bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("commit_upgrade", typeArgsList, typeParamsList, []string{
		"&mut DeployerState",
		"UpgradeReceipt",
	}, []any{
		state,
		receipt,
	}, nil)
}

// CommitUpgradeWithArgs encodes a call to the commit_upgrade Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsDeployerEncoder) CommitUpgradeWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut DeployerState",
		"UpgradeReceipt",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("commit_upgrade", typeArgsList, typeParamsList, expectedParams, args, nil)
}
