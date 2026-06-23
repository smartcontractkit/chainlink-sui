// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package module_mcms_deployer

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

const FunctionInfo = `[{"package":"mcms","module":"mcms_deployer","name":"authorize_upgrade","parameters":[{"name":"_","type":"OwnerCap"},{"name":"state","type":"DeployerState"},{"name":"policy","type":"u8"},{"name":"digest","type":"vector<u8>"},{"name":"package_address","type":"address"}]},{"package":"mcms","module":"mcms_deployer","name":"commit_upgrade","parameters":[{"name":"state","type":"DeployerState"},{"name":"receipt","type":"UpgradeReceipt"}]},{"package":"mcms","module":"mcms_deployer","name":"has_upgrade_cap","parameters":[{"name":"state","type":"DeployerState"},{"name":"package_address","type":"address"}]},{"package":"mcms","module":"mcms_deployer","name":"register_upgrade_cap","parameters":[{"name":"state","type":"DeployerState"},{"name":"registry","type":"Registry"},{"name":"upgrade_cap","type":"UpgradeCap"}]},{"package":"mcms","module":"mcms_deployer","name":"release_upgrade_cap","parameters":[{"name":"state","type":"DeployerState"},{"name":"registry","type":"Registry"},{"name":"_proof","type":"T"}]}]`

type IMcmsDeployer interface {
	RegisterUpgradeCap(ctx context.Context, opts *bind.CallOpts, state bind.Object, registry bind.Object, upgradeCap bind.Object) (*models.SuiTransactionBlockResponse, error)
	AuthorizeUpgrade(ctx context.Context, opts *bind.CallOpts, param bind.Object, state bind.Object, policy byte, digest []byte, packageAddress string) (*models.SuiTransactionBlockResponse, error)
	CommitUpgrade(ctx context.Context, opts *bind.CallOpts, state bind.Object, receipt bind.Object) (*models.SuiTransactionBlockResponse, error)
	ReleaseUpgradeCap(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, registry bind.Object, proof bind.Object) (*models.SuiTransactionBlockResponse, error)
	HasUpgradeCap(ctx context.Context, opts *bind.CallOpts, state bind.Object, packageAddress string) (*models.SuiTransactionBlockResponse, error)
	DevInspect() IMcmsDeployerDevInspect
	Encoder() McmsDeployerEncoder
	Bound() bind.IBoundContract
}

type IMcmsDeployerDevInspect interface {
	AuthorizeUpgrade(ctx context.Context, opts *bind.CallOpts, param bind.Object, state bind.Object, policy byte, digest []byte, packageAddress string) (bind.Object, error)
	ReleaseUpgradeCap(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, registry bind.Object, proof bind.Object) (bind.Object, error)
	HasUpgradeCap(ctx context.Context, opts *bind.CallOpts, state bind.Object, packageAddress string) (bool, error)
}

type McmsDeployerEncoder interface {
	RegisterUpgradeCap(state bind.Object, registry bind.Object, upgradeCap bind.Object) (*bind.EncodedCall, error)
	RegisterUpgradeCapWithArgs(args ...any) (*bind.EncodedCall, error)
	AuthorizeUpgrade(param bind.Object, state bind.Object, policy byte, digest []byte, packageAddress string) (*bind.EncodedCall, error)
	AuthorizeUpgradeWithArgs(args ...any) (*bind.EncodedCall, error)
	CommitUpgrade(state bind.Object, receipt bind.Object) (*bind.EncodedCall, error)
	CommitUpgradeWithArgs(args ...any) (*bind.EncodedCall, error)
	ReleaseUpgradeCap(typeArgs []string, state bind.Object, registry bind.Object, proof bind.Object) (*bind.EncodedCall, error)
	ReleaseUpgradeCapWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	HasUpgradeCap(state bind.Object, packageAddress string) (*bind.EncodedCall, error)
	HasUpgradeCapWithArgs(args ...any) (*bind.EncodedCall, error)
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

func NewMcmsDeployer(packageID string, chainClient client.BindingsClient) (IMcmsDeployer, error) {
	contract, err := bind.NewBoundContract(packageID, "mcms", "mcms_deployer", chainClient)
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

// CommitUpgrade executes the commit_upgrade Move function.
func (c *McmsDeployerContract) CommitUpgrade(ctx context.Context, opts *bind.CallOpts, state bind.Object, receipt bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsDeployerEncoder.CommitUpgrade(state, receipt)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ReleaseUpgradeCap executes the release_upgrade_cap Move function.
func (c *McmsDeployerContract) ReleaseUpgradeCap(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, registry bind.Object, proof bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsDeployerEncoder.ReleaseUpgradeCap(typeArgs, state, registry, proof)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// HasUpgradeCap executes the has_upgrade_cap Move function.
func (c *McmsDeployerContract) HasUpgradeCap(ctx context.Context, opts *bind.CallOpts, state bind.Object, packageAddress string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsDeployerEncoder.HasUpgradeCap(state, packageAddress)
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
	var result bind.Object
	if err := bind.DecodeJSONReturn(results[0], &result); err != nil {
		return bind.Object{}, fmt.Errorf("failed to decode return value: %w", err)
	}
	return result, nil
}

// ReleaseUpgradeCap executes the release_upgrade_cap Move function using DevInspect to get return values.
//
// Returns: UpgradeCap
func (d *McmsDeployerDevInspect) ReleaseUpgradeCap(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, registry bind.Object, proof bind.Object) (bind.Object, error) {
	encoded, err := d.contract.mcmsDeployerEncoder.ReleaseUpgradeCap(typeArgs, state, registry, proof)
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
	var result bind.Object
	if err := bind.DecodeJSONReturn(results[0], &result); err != nil {
		return bind.Object{}, fmt.Errorf("failed to decode return value: %w", err)
	}
	return result, nil
}

// HasUpgradeCap executes the has_upgrade_cap Move function using DevInspect to get return values.
//
// Returns: bool
func (d *McmsDeployerDevInspect) HasUpgradeCap(ctx context.Context, opts *bind.CallOpts, state bind.Object, packageAddress string) (bool, error) {
	encoded, err := d.contract.mcmsDeployerEncoder.HasUpgradeCap(state, packageAddress)
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
	var result bool
	if err := bind.DecodeJSONReturn(results[0], &result); err != nil {
		return false, fmt.Errorf("failed to decode return value: %w", err)
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

// ReleaseUpgradeCap encodes a call to the release_upgrade_cap Move function.
func (c mcmsDeployerEncoder) ReleaseUpgradeCap(typeArgs []string, state bind.Object, registry bind.Object, proof bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("release_upgrade_cap", typeArgsList, typeParamsList, []string{
		"&mut DeployerState",
		"&Registry",
		"T",
	}, []any{
		state,
		registry,
		proof,
	}, []string{
		"UpgradeCap",
	})
}

// ReleaseUpgradeCapWithArgs encodes a call to the release_upgrade_cap Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsDeployerEncoder) ReleaseUpgradeCapWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut DeployerState",
		"&Registry",
		"T",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("release_upgrade_cap", typeArgsList, typeParamsList, expectedParams, args, []string{
		"UpgradeCap",
	})
}

// HasUpgradeCap encodes a call to the has_upgrade_cap Move function.
func (c mcmsDeployerEncoder) HasUpgradeCap(state bind.Object, packageAddress string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("has_upgrade_cap", typeArgsList, typeParamsList, []string{
		"&DeployerState",
		"address",
	}, []any{
		state,
		packageAddress,
	}, []string{
		"bool",
	})
}

// HasUpgradeCapWithArgs encodes a call to the has_upgrade_cap Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsDeployerEncoder) HasUpgradeCapWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&DeployerState",
		"address",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("has_upgrade_cap", typeArgsList, typeParamsList, expectedParams, args, []string{
		"bool",
	})
}
