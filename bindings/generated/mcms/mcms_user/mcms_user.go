// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package module_mcms_user

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

type IMcmsUser interface {
	FunctionOne(ctx context.Context, opts *bind.CallOpts, userData bind.Object, ownerCap bind.Object, arg1 string, arg2 []byte) (*models.SuiTransactionBlockResponse, error)
	FunctionTwo(ctx context.Context, opts *bind.CallOpts, userData bind.Object, ownerCap bind.Object, arg1 string, arg2 *big.Int) (*models.SuiTransactionBlockResponse, error)
	RegisterMcmsEntrypoint(ctx context.Context, opts *bind.CallOpts, ownerCap bind.Object, registry bind.Object, userData bind.Object) (*models.SuiTransactionBlockResponse, error)
	RegisterUpgradeCap(ctx context.Context, opts *bind.CallOpts, state bind.Object, upgradeCap bind.Object, registry bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsEntrypoint(ctx context.Context, opts *bind.CallOpts, userData bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	GetOwnerCap(ctx context.Context, opts *bind.CallOpts, userData bind.Object) (*models.SuiTransactionBlockResponse, error)
	GetInvocations(ctx context.Context, opts *bind.CallOpts, userData bind.Object) (*models.SuiTransactionBlockResponse, error)
	GetFieldA(ctx context.Context, opts *bind.CallOpts, userData bind.Object) (*models.SuiTransactionBlockResponse, error)
	GetFieldB(ctx context.Context, opts *bind.CallOpts, userData bind.Object) (*models.SuiTransactionBlockResponse, error)
	GetFieldC(ctx context.Context, opts *bind.CallOpts, userData bind.Object) (*models.SuiTransactionBlockResponse, error)
	GetFieldD(ctx context.Context, opts *bind.CallOpts, userData bind.Object) (*models.SuiTransactionBlockResponse, error)
	DevInspect() IMcmsUserDevInspect
	Encoder() McmsUserEncoder
}

type IMcmsUserDevInspect interface {
	GetOwnerCap(ctx context.Context, opts *bind.CallOpts, userData bind.Object) (bind.Object, error)
	GetInvocations(ctx context.Context, opts *bind.CallOpts, userData bind.Object) (byte, error)
	GetFieldA(ctx context.Context, opts *bind.CallOpts, userData bind.Object) (string, error)
	GetFieldB(ctx context.Context, opts *bind.CallOpts, userData bind.Object) ([]byte, error)
	GetFieldC(ctx context.Context, opts *bind.CallOpts, userData bind.Object) (string, error)
	GetFieldD(ctx context.Context, opts *bind.CallOpts, userData bind.Object) (*big.Int, error)
}

type McmsUserEncoder interface {
	FunctionOne(userData bind.Object, ownerCap bind.Object, arg1 string, arg2 []byte) (*bind.EncodedCall, error)
	FunctionOneWithArgs(args ...any) (*bind.EncodedCall, error)
	FunctionTwo(userData bind.Object, ownerCap bind.Object, arg1 string, arg2 *big.Int) (*bind.EncodedCall, error)
	FunctionTwoWithArgs(args ...any) (*bind.EncodedCall, error)
	RegisterMcmsEntrypoint(ownerCap bind.Object, registry bind.Object, userData bind.Object) (*bind.EncodedCall, error)
	RegisterMcmsEntrypointWithArgs(args ...any) (*bind.EncodedCall, error)
	RegisterUpgradeCap(state bind.Object, upgradeCap bind.Object, registry bind.Object) (*bind.EncodedCall, error)
	RegisterUpgradeCapWithArgs(args ...any) (*bind.EncodedCall, error)
	McmsEntrypoint(userData bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsEntrypointWithArgs(args ...any) (*bind.EncodedCall, error)
	GetOwnerCap(userData bind.Object) (*bind.EncodedCall, error)
	GetOwnerCapWithArgs(args ...any) (*bind.EncodedCall, error)
	GetInvocations(userData bind.Object) (*bind.EncodedCall, error)
	GetInvocationsWithArgs(args ...any) (*bind.EncodedCall, error)
	GetFieldA(userData bind.Object) (*bind.EncodedCall, error)
	GetFieldAWithArgs(args ...any) (*bind.EncodedCall, error)
	GetFieldB(userData bind.Object) (*bind.EncodedCall, error)
	GetFieldBWithArgs(args ...any) (*bind.EncodedCall, error)
	GetFieldC(userData bind.Object) (*bind.EncodedCall, error)
	GetFieldCWithArgs(args ...any) (*bind.EncodedCall, error)
	GetFieldD(userData bind.Object) (*bind.EncodedCall, error)
	GetFieldDWithArgs(args ...any) (*bind.EncodedCall, error)
}

type McmsUserContract struct {
	*bind.BoundContract
	mcmsUserEncoder
	devInspect *McmsUserDevInspect
}

type McmsUserDevInspect struct {
	contract *McmsUserContract
}

var _ IMcmsUser = (*McmsUserContract)(nil)
var _ IMcmsUserDevInspect = (*McmsUserDevInspect)(nil)

func NewMcmsUser(packageID string, client sui.ISuiAPI) (*McmsUserContract, error) {
	contract, err := bind.NewBoundContract(packageID, "mcms_test", "mcms_user", client)
	if err != nil {
		return nil, err
	}

	c := &McmsUserContract{
		BoundContract:   contract,
		mcmsUserEncoder: mcmsUserEncoder{BoundContract: contract},
	}
	c.devInspect = &McmsUserDevInspect{contract: c}
	return c, nil
}

func (c *McmsUserContract) Encoder() McmsUserEncoder {
	return c.mcmsUserEncoder
}

func (c *McmsUserContract) DevInspect() IMcmsUserDevInspect {
	return c.devInspect
}

func (c *McmsUserContract) BuildPTB(ctx context.Context, ptb *transaction.Transaction, encoded *bind.EncodedCall) (*transaction.Argument, error) {
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

type UserData struct {
	Id          string      `move:"sui::object::UID"`
	Invocations byte        `move:"u8"`
	A           string      `move:"0x1::string::String"`
	B           []byte      `move:"vector<u8>"`
	C           string      `move:"address"`
	D           *big.Int    `move:"u128"`
	OwnerCap    bind.Object `move:"ID"`
}

type OwnerCap struct {
	Id string `move:"sui::object::UID"`
}

type MCMS_USER struct {
}

type SampleMcmsCallback struct {
}

type bcsUserData struct {
	Id          string
	Invocations byte
	A           string
	B           []byte
	C           [32]byte
	D           *big.Int
	OwnerCap    bind.Object
}

func convertUserDataFromBCS(bcs bcsUserData) UserData {
	return UserData{
		Id:          bcs.Id,
		Invocations: bcs.Invocations,
		A:           bcs.A,
		B:           bcs.B,
		C:           fmt.Sprintf("0x%x", bcs.C),
		D:           bcs.D,
		OwnerCap:    bcs.OwnerCap,
	}
}

func init() {
	bind.RegisterStructDecoder("mcms_test::mcms_user::UserData", func(data []byte) (interface{}, error) {
		var temp bcsUserData
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result := convertUserDataFromBCS(temp)
		return result, nil
	})
	bind.RegisterStructDecoder("mcms_test::mcms_user::OwnerCap", func(data []byte) (interface{}, error) {
		var result OwnerCap
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("mcms_test::mcms_user::MCMS_USER", func(data []byte) (interface{}, error) {
		var result MCMS_USER
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("mcms_test::mcms_user::SampleMcmsCallback", func(data []byte) (interface{}, error) {
		var result SampleMcmsCallback
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
}

// FunctionOne executes the function_one Move function.
func (c *McmsUserContract) FunctionOne(ctx context.Context, opts *bind.CallOpts, userData bind.Object, ownerCap bind.Object, arg1 string, arg2 []byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsUserEncoder.FunctionOne(userData, ownerCap, arg1, arg2)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// FunctionTwo executes the function_two Move function.
func (c *McmsUserContract) FunctionTwo(ctx context.Context, opts *bind.CallOpts, userData bind.Object, ownerCap bind.Object, arg1 string, arg2 *big.Int) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsUserEncoder.FunctionTwo(userData, ownerCap, arg1, arg2)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// RegisterMcmsEntrypoint executes the register_mcms_entrypoint Move function.
func (c *McmsUserContract) RegisterMcmsEntrypoint(ctx context.Context, opts *bind.CallOpts, ownerCap bind.Object, registry bind.Object, userData bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsUserEncoder.RegisterMcmsEntrypoint(ownerCap, registry, userData)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// RegisterUpgradeCap executes the register_upgrade_cap Move function.
func (c *McmsUserContract) RegisterUpgradeCap(ctx context.Context, opts *bind.CallOpts, state bind.Object, upgradeCap bind.Object, registry bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsUserEncoder.RegisterUpgradeCap(state, upgradeCap, registry)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsEntrypoint executes the mcms_entrypoint Move function.
func (c *McmsUserContract) McmsEntrypoint(ctx context.Context, opts *bind.CallOpts, userData bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsUserEncoder.McmsEntrypoint(userData, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetOwnerCap executes the get_owner_cap Move function.
func (c *McmsUserContract) GetOwnerCap(ctx context.Context, opts *bind.CallOpts, userData bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsUserEncoder.GetOwnerCap(userData)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetInvocations executes the get_invocations Move function.
func (c *McmsUserContract) GetInvocations(ctx context.Context, opts *bind.CallOpts, userData bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsUserEncoder.GetInvocations(userData)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetFieldA executes the get_field_a Move function.
func (c *McmsUserContract) GetFieldA(ctx context.Context, opts *bind.CallOpts, userData bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsUserEncoder.GetFieldA(userData)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetFieldB executes the get_field_b Move function.
func (c *McmsUserContract) GetFieldB(ctx context.Context, opts *bind.CallOpts, userData bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsUserEncoder.GetFieldB(userData)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetFieldC executes the get_field_c Move function.
func (c *McmsUserContract) GetFieldC(ctx context.Context, opts *bind.CallOpts, userData bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsUserEncoder.GetFieldC(userData)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetFieldD executes the get_field_d Move function.
func (c *McmsUserContract) GetFieldD(ctx context.Context, opts *bind.CallOpts, userData bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsUserEncoder.GetFieldD(userData)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetOwnerCap executes the get_owner_cap Move function using DevInspect to get return values.
//
// Returns: ID
func (d *McmsUserDevInspect) GetOwnerCap(ctx context.Context, opts *bind.CallOpts, userData bind.Object) (bind.Object, error) {
	encoded, err := d.contract.mcmsUserEncoder.GetOwnerCap(userData)
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

// GetInvocations executes the get_invocations Move function using DevInspect to get return values.
//
// Returns: u8
func (d *McmsUserDevInspect) GetInvocations(ctx context.Context, opts *bind.CallOpts, userData bind.Object) (byte, error) {
	encoded, err := d.contract.mcmsUserEncoder.GetInvocations(userData)
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

// GetFieldA executes the get_field_a Move function using DevInspect to get return values.
//
// Returns: 0x1::string::String
func (d *McmsUserDevInspect) GetFieldA(ctx context.Context, opts *bind.CallOpts, userData bind.Object) (string, error) {
	encoded, err := d.contract.mcmsUserEncoder.GetFieldA(userData)
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

// GetFieldB executes the get_field_b Move function using DevInspect to get return values.
//
// Returns: vector<u8>
func (d *McmsUserDevInspect) GetFieldB(ctx context.Context, opts *bind.CallOpts, userData bind.Object) ([]byte, error) {
	encoded, err := d.contract.mcmsUserEncoder.GetFieldB(userData)
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

// GetFieldC executes the get_field_c Move function using DevInspect to get return values.
//
// Returns: address
func (d *McmsUserDevInspect) GetFieldC(ctx context.Context, opts *bind.CallOpts, userData bind.Object) (string, error) {
	encoded, err := d.contract.mcmsUserEncoder.GetFieldC(userData)
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

// GetFieldD executes the get_field_d Move function using DevInspect to get return values.
//
// Returns: u128
func (d *McmsUserDevInspect) GetFieldD(ctx context.Context, opts *bind.CallOpts, userData bind.Object) (*big.Int, error) {
	encoded, err := d.contract.mcmsUserEncoder.GetFieldD(userData)
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
	result, ok := results[0].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("unexpected return type: expected *big.Int, got %T", results[0])
	}
	return result, nil
}

type mcmsUserEncoder struct {
	*bind.BoundContract
}

// FunctionOne encodes a call to the function_one Move function.
func (c mcmsUserEncoder) FunctionOne(userData bind.Object, ownerCap bind.Object, arg1 string, arg2 []byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("function_one", typeArgsList, typeParamsList, []string{
		"&mut UserData",
		"&OwnerCap",
		"String",
		"vector<u8>",
	}, []any{
		userData,
		ownerCap,
		arg1,
		arg2,
	}, nil)
}

// FunctionOneWithArgs encodes a call to the function_one Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsUserEncoder) FunctionOneWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut UserData",
		"&OwnerCap",
		"String",
		"vector<u8>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("function_one", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// FunctionTwo encodes a call to the function_two Move function.
func (c mcmsUserEncoder) FunctionTwo(userData bind.Object, ownerCap bind.Object, arg1 string, arg2 *big.Int) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("function_two", typeArgsList, typeParamsList, []string{
		"&mut UserData",
		"&OwnerCap",
		"address",
		"u128",
	}, []any{
		userData,
		ownerCap,
		arg1,
		arg2,
	}, nil)
}

// FunctionTwoWithArgs encodes a call to the function_two Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsUserEncoder) FunctionTwoWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut UserData",
		"&OwnerCap",
		"address",
		"u128",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("function_two", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// RegisterMcmsEntrypoint encodes a call to the register_mcms_entrypoint Move function.
func (c mcmsUserEncoder) RegisterMcmsEntrypoint(ownerCap bind.Object, registry bind.Object, userData bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("register_mcms_entrypoint", typeArgsList, typeParamsList, []string{
		"OwnerCap",
		"&mut Registry",
		"&UserData",
	}, []any{
		ownerCap,
		registry,
		userData,
	}, nil)
}

// RegisterMcmsEntrypointWithArgs encodes a call to the register_mcms_entrypoint Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsUserEncoder) RegisterMcmsEntrypointWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"OwnerCap",
		"&mut Registry",
		"&UserData",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("register_mcms_entrypoint", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// RegisterUpgradeCap encodes a call to the register_upgrade_cap Move function.
func (c mcmsUserEncoder) RegisterUpgradeCap(state bind.Object, upgradeCap bind.Object, registry bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("register_upgrade_cap", typeArgsList, typeParamsList, []string{
		"&mut DeployerState",
		"UpgradeCap",
		"&mut Registry",
	}, []any{
		state,
		upgradeCap,
		registry,
	}, nil)
}

// RegisterUpgradeCapWithArgs encodes a call to the register_upgrade_cap Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsUserEncoder) RegisterUpgradeCapWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut DeployerState",
		"UpgradeCap",
		"&mut Registry",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("register_upgrade_cap", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// McmsEntrypoint encodes a call to the mcms_entrypoint Move function.
func (c mcmsUserEncoder) McmsEntrypoint(userData bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_entrypoint", typeArgsList, typeParamsList, []string{
		"&mut UserData",
		"&mut Registry",
		"ExecutingCallbackParams",
	}, []any{
		userData,
		registry,
		params,
	}, nil)
}

// McmsEntrypointWithArgs encodes a call to the mcms_entrypoint Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsUserEncoder) McmsEntrypointWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut UserData",
		"&mut Registry",
		"ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_entrypoint", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// GetOwnerCap encodes a call to the get_owner_cap Move function.
func (c mcmsUserEncoder) GetOwnerCap(userData bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_owner_cap", typeArgsList, typeParamsList, []string{
		"&UserData",
	}, []any{
		userData,
	}, []string{
		"ID",
	})
}

// GetOwnerCapWithArgs encodes a call to the get_owner_cap Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsUserEncoder) GetOwnerCapWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&UserData",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_owner_cap", typeArgsList, typeParamsList, expectedParams, args, []string{
		"ID",
	})
}

// GetInvocations encodes a call to the get_invocations Move function.
func (c mcmsUserEncoder) GetInvocations(userData bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_invocations", typeArgsList, typeParamsList, []string{
		"&UserData",
	}, []any{
		userData,
	}, []string{
		"u8",
	})
}

// GetInvocationsWithArgs encodes a call to the get_invocations Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsUserEncoder) GetInvocationsWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&UserData",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_invocations", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u8",
	})
}

// GetFieldA encodes a call to the get_field_a Move function.
func (c mcmsUserEncoder) GetFieldA(userData bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_field_a", typeArgsList, typeParamsList, []string{
		"&UserData",
	}, []any{
		userData,
	}, []string{
		"0x1::string::String",
	})
}

// GetFieldAWithArgs encodes a call to the get_field_a Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsUserEncoder) GetFieldAWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&UserData",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_field_a", typeArgsList, typeParamsList, expectedParams, args, []string{
		"0x1::string::String",
	})
}

// GetFieldB encodes a call to the get_field_b Move function.
func (c mcmsUserEncoder) GetFieldB(userData bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_field_b", typeArgsList, typeParamsList, []string{
		"&UserData",
	}, []any{
		userData,
	}, []string{
		"vector<u8>",
	})
}

// GetFieldBWithArgs encodes a call to the get_field_b Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsUserEncoder) GetFieldBWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&UserData",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_field_b", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<u8>",
	})
}

// GetFieldC encodes a call to the get_field_c Move function.
func (c mcmsUserEncoder) GetFieldC(userData bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_field_c", typeArgsList, typeParamsList, []string{
		"&UserData",
	}, []any{
		userData,
	}, []string{
		"address",
	})
}

// GetFieldCWithArgs encodes a call to the get_field_c Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsUserEncoder) GetFieldCWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&UserData",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_field_c", typeArgsList, typeParamsList, expectedParams, args, []string{
		"address",
	})
}

// GetFieldD encodes a call to the get_field_d Move function.
func (c mcmsUserEncoder) GetFieldD(userData bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_field_d", typeArgsList, typeParamsList, []string{
		"&UserData",
	}, []any{
		userData,
	}, []string{
		"u128",
	})
}

// GetFieldDWithArgs encodes a call to the get_field_d Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsUserEncoder) GetFieldDWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&UserData",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_field_d", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u128",
	})
}
