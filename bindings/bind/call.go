package bind

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/signer"
	"github.com/block-vision/sui-go-sdk/transaction"

	bindutils "github.com/smartcontractkit/chainlink-sui/bindings/utils"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
)

const (
	DefaultGasBudget uint64 = 10_000_000_000
	AddressType             = "address"
)

type IBoundContract interface {
	GetPackageID() string
	GetPackageName() string
	GetModuleName() string
	AppendPTB(ctx context.Context, opts *CallOpts, ptb *transaction.Transaction, encoded *EncodedCall) (*transaction.Argument, error)
	Call(ctx context.Context, opts *CallOpts, encoded *EncodedCall) ([]any, error)
	ExecuteTransaction(ctx context.Context, opts *CallOpts, encoded *EncodedCall) (*models.SuiTransactionBlockResponse, error)
}

var _ IBoundContract = (*BoundContract)(nil)

type BoundContract struct {
	packageID   string
	packageName string
	moduleName  string
	client      client.BindingsClient
}

func (c *BoundContract) GetPackageID() string {
	return c.packageID
}

func (c *BoundContract) GetPackageName() string {
	return c.packageName
}

func (c *BoundContract) GetModuleName() string {
	return c.moduleName
}

func NewBoundContract(packageID string, packageName, moduleName string, chainClient client.BindingsClient) (*BoundContract, error) {
	normalizedID, err := bindutils.ConvertAddressToString(packageID)
	if err != nil {
		return nil, fmt.Errorf("invalid package ID %s: %w", packageID, err)
	}

	return &BoundContract{
		packageID:   normalizedID,
		packageName: packageName,
		moduleName:  moduleName,
		client:      chainClient,
	}, nil
}

type ModuleInformation struct {
	PackageID   string
	PackageName string
	ModuleName  string
}

func (m *ModuleInformation) String() string {
	return fmt.Sprintf("%s::%s::%s", m.PackageID, m.PackageName, m.ModuleName)
}

func (c *BoundContract) Call(ctx context.Context, opts *CallOpts, encoded *EncodedCall) ([]any, error) {
	if opts == nil || opts.Signer == nil {
		return nil, fmt.Errorf("CallOpts with Signer is required")
	}

	signerAddressStr, err := opts.Signer.GetAddress()
	if err != nil {
		return nil, fmt.Errorf("failed to get signer address: %w", err)
	}

	signerAddress, err := bindutils.ConvertAddressToString(signerAddressStr)
	if err != nil {
		return nil, fmt.Errorf("invalid signer address %v: %w", signerAddressStr, err)
	}

	resolver := opts.ObjectResolver
	if resolver == nil {
		resolver = NewObjectResolver(c.client)
	}

	resolvedEncodedArgs := make([]*EncodedCallArgument, len(encoded.CallArgs))
	for i, encArg := range encoded.CallArgs {
		if encArg == nil {
			return nil, fmt.Errorf("nil EncodedCallArgument at index %d", i)
		}

		if encArg.IsArgument() {
			resolvedEncodedArgs[i] = encArg
		} else if encArg.IsCallArg() {
			resolved, resolveErr := resolver.ResolveCallArg(ctx, encArg.CallArg, encArg.TypeName)
			if resolveErr != nil {
				return nil, fmt.Errorf("failed to resolve CallArg at index %d: %w", i, resolveErr)
			}
			resolvedEncodedArg := NewEncodedCallArgFromCallArgWithType(resolved, encArg.TypeName)
			resolvedEncodedArgs[i] = resolvedEncodedArg
		} else {
			return nil, errors.New("empty EncodedCallArgument")
		}
	}

	callArgManager := NewCallArgManager()

	arguments, err := callArgManager.ConvertEncodedCallArgsToArguments(resolvedEncodedArgs)
	if err != nil {
		return nil, fmt.Errorf("failed to convert EncodedCallArguments to Arguments: %w", err)
	}

	ptb := transaction.NewTransaction()
	ptb.SetSender(models.SuiAddress(signerAddress))

	inputs := callArgManager.GetInputs()
	if len(inputs) > 0 {
		if ptb.Data.V1 == nil || ptb.Data.V1.Kind == nil || ptb.Data.V1.Kind.ProgrammableTransaction == nil {
			return nil, errors.New("unexpected PTB with missing fields")
		}
		ptb.Data.V1.Kind.ProgrammableTransaction.Inputs = inputs
	}

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

	_ = ptb.MoveCall(
		models.SuiAddress(encoded.Module.PackageID),
		encoded.Module.ModuleName,
		encoded.Function,
		typeTagValues,
		argumentValues,
	)

	bcsBytes, err := buildSimulateBCS(ctx, c.client, ptb, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to build simulate transaction bytes: %w", err)
	}

	results, err := c.client.SimulatePTB(ctx, bcsBytes)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (c *BoundContract) ExecuteTransaction(ctx context.Context, opts *CallOpts, encoded *EncodedCall) (*models.SuiTransactionBlockResponse, error) {
	if opts == nil || opts.Signer == nil {
		return nil, fmt.Errorf("CallOpts with Signer is required")
	}

	ptb := transaction.NewTransaction()
	_, err := c.AppendPTB(ctx, opts, ptb, encoded)
	if err != nil {
		return nil, fmt.Errorf("failed to add encoded call to PTB: %w", err)
	}

	return ExecutePTB(ctx, opts, c.client, ptb)
}

func GetObjectRef(ctx context.Context, chainClient client.BindingsClient, objectID string) (*models.SuiObjectRef, error) {
	normalizedID, err := bindutils.ConvertAddressToString(objectID)
	if err != nil {
		return nil, fmt.Errorf("invalid object ID %v: %w", objectID, err)
	}

	obj, err := chainClient.ReadObjectId(ctx, normalizedID)
	if err != nil {
		return nil, err
	}

	return mapGrpcCoinToObjectRef(obj), nil
}

func typeTagToString(tag *transaction.TypeTag) string {
	if tag == nil {
		return ""
	}

	if tag.Bool != nil {
		return "bool"
	}
	if tag.U8 != nil {
		return "u8"
	}
	if tag.U16 != nil {
		return "u16"
	}
	if tag.U32 != nil {
		return "u32"
	}
	if tag.U64 != nil {
		return "u64"
	}
	if tag.U128 != nil {
		return "u128"
	}
	if tag.U256 != nil {
		return "u256"
	}
	if tag.Address != nil {
		return AddressType
	}

	if tag.Vector != nil {
		innerType := typeTagToString(tag.Vector)
		return fmt.Sprintf("vector<%s>", innerType)
	}

	if tag.Struct != nil {
		addr, err := bindutils.ConvertBytesToAddress(tag.Struct.Address[:])
		if err != nil {
			return fmt.Sprintf("invalid_address::%s::%s", tag.Struct.Module, tag.Struct.Name)
		}
		baseType := fmt.Sprintf("%s::%s::%s", addr, tag.Struct.Module, tag.Struct.Name)

		if len(tag.Struct.TypeParams) > 0 {
			typeParams := make([]string, len(tag.Struct.TypeParams))
			for i, param := range tag.Struct.TypeParams {
				typeParams[i] = typeTagToString(param)
			}

			return fmt.Sprintf("%s<%s>", baseType, strings.Join(typeParams, ","))
		}

		return baseType
	}

	return ""
}

func (c *BoundContract) AppendPTB(ctx context.Context, opts *CallOpts, ptb *transaction.Transaction, encoded *EncodedCall) (*transaction.Argument, error) {
	if opts.ObjectResolver == nil {
		opts.ObjectResolver = NewObjectResolver(c.client)
	}

	resolvedEncodedArgs := make([]*EncodedCallArgument, len(encoded.CallArgs))
	for i, encArg := range encoded.CallArgs {
		if encArg == nil {
			return nil, fmt.Errorf("nil EncodedCallArgument at index %d", i)
		}

		if encArg.IsArgument() {
			resolvedEncodedArgs[i] = encArg
		} else if encArg.IsCallArg() {
			resolved, resolveErr := opts.ObjectResolver.ResolveCallArg(ctx, encArg.CallArg, encArg.TypeName)
			if resolveErr != nil {
				return nil, fmt.Errorf("failed to resolve CallArg at index %d: %w", i, resolveErr)
			}
			resolvedEncodedArg := NewEncodedCallArgFromCallArgWithType(resolved, encArg.TypeName)
			resolvedEncodedArgs[i] = resolvedEncodedArg
		} else {
			return nil, errors.New("empty EncodedCallArgument")
		}
	}

	var existingInputs []*transaction.CallArg
	if ptb.Data.V1 != nil && ptb.Data.V1.Kind != nil && ptb.Data.V1.Kind.ProgrammableTransaction != nil {
		existingInputs = ptb.Data.V1.Kind.ProgrammableTransaction.Inputs
	}

	callArgManager := NewCallArgManagerWithExisting(existingInputs)

	arguments, err := callArgManager.ConvertEncodedCallArgsToArguments(resolvedEncodedArgs)
	if err != nil {
		return nil, fmt.Errorf("failed to convert EncodedCallArguments to Arguments: %w", err)
	}

	inputs := callArgManager.GetInputs()
	if ptb.Data.V1 == nil || ptb.Data.V1.Kind == nil || ptb.Data.V1.Kind.ProgrammableTransaction == nil {
		return nil, errors.New("unexpected PTB with missing fields")
	}
	ptb.Data.V1.Kind.ProgrammableTransaction.Inputs = inputs

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

	arg := ptb.MoveCall(
		models.SuiAddress(encoded.Module.PackageID),
		encoded.Module.ModuleName,
		encoded.Function,
		typeTagValues,
		argumentValues,
	)

	return &arg, nil
}

func ExecutePTB(ctx context.Context, opts *CallOpts, chainClient client.BindingsClient, ptb *transaction.Transaction) (*models.SuiTransactionBlockResponse, error) {
	if opts == nil || opts.Signer == nil {
		return nil, fmt.Errorf("CallOpts with Signer is required")
	}

	signerAddressStr, err := opts.Signer.GetAddress()
	if err != nil {
		return nil, fmt.Errorf("failed to get signer address: %w", err)
	}

	signerAddress, err := bindutils.ConvertAddressToString(signerAddressStr)
	if err != nil {
		return nil, fmt.Errorf("invalid signer address %v: %w", signerAddressStr, err)
	}

	if ptb.Data.V1.Sender == nil {
		ptb.SetSender(models.SuiAddress(signerAddress))
	}

	if ptb.Data.V1.GasData.Budget == nil {
		budget := DefaultGasBudget
		if opts.GasBudget != nil {
			budget = *opts.GasBudget
		}
		ptb.SetGasBudget(budget)
	}

	if ptb.Data.V1.GasData.Price == nil {
		gasPrice, gasPriceErr := chainClient.GetReferenceGasPrice(ctx)
		if gasPriceErr != nil {
			return nil, fmt.Errorf("failed to get reference gas price: %w", gasPriceErr)
		}
		ptb.SetGasPrice(gasPrice.Uint64())
	}

	if ptb.Data.V1.GasData.Owner == nil {
		normalizedSigner, normalizationErr := bindutils.ConvertAddressToString(signerAddressStr)
		if normalizationErr != nil {
			return nil, fmt.Errorf("invalid signer address for gas owner %v: %w", signerAddressStr, normalizationErr)
		}
		ptb.SetGasOwner(models.SuiAddress(normalizedSigner))
	}

	if ptb.Data.V1.GasData.Payment == nil {
		var gasRef *models.SuiObjectRef
		if opts.GasObject != "" {
			gasRef, err = ToSuiObjectRef(ctx, chainClient, opts.GasObject, signerAddress)
		} else {
			gasRef, err = FetchDefaultGasCoinRef(ctx, chainClient, signerAddress)
		}
		if err != nil {
			return nil, fmt.Errorf("failed to get gas object: %w", err)
		}

		if gasRef != nil {
			objIdBytes, objIdErr := bindutils.ConvertStringToAddressBytes(gasRef.ObjectId)
			if objIdErr != nil {
				return nil, fmt.Errorf("failed to convert gas object ID: %w", objIdErr)
			}
			digestBytes, digestErr := bindutils.ConvertStringToDigestBytes(gasRef.Digest)
			if digestErr != nil {
				return nil, fmt.Errorf("failed to convert gas object digest: %w", digestErr)
			}

			payment := []transaction.SuiObjectRef{{
				ObjectId: *objIdBytes,
				Version:  gasRef.Version,
				Digest:   *digestBytes,
			}}
			ptb.SetGasPayment(payment)
		}
	}

	ptb.SetSigner(&signer.Signer{Address: signerAddress})

	txBytes, err := ptb.BuildBCSBytes(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to build transaction bytes: %w", err)
	}

	return SignAndSendTx(ctx, opts.Signer, chainClient, txBytes, opts.WaitForExecution)
}
