package bind

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/block-vision/sui-go-sdk/transaction"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	bindutils "github.com/smartcontractkit/chainlink-sui/bindings/utils"
)

const (
	// DefaultGasBudget is the default gas budget for transactions
	DefaultGasBudget uint64 = 10_000_000
)

type BoundContract struct {
	packageID   string
	packageName string
	moduleName  string
	client      sui.ISuiAPI
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

func NewBoundContract(packageID string, packageName, moduleName string, client sui.ISuiAPI) (*BoundContract, error) {
	normalizedID, err := bindutils.ConvertAddressToString(packageID)
	if err != nil {
		return nil, fmt.Errorf("invalid package ID %s: %w", packageID, err)
	}

	return &BoundContract{
		packageID:   normalizedID,
		packageName: packageName,
		moduleName:  moduleName,
		client:      client,
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

// TODO: dedupe transaction argument generation code from ExecuteTransaction
func (c *BoundContract) Call(ctx context.Context, opts *CallOpts, encoded *EncodedCall) ([]any, error) {
	if opts == nil || opts.Signer == nil {
		return nil, fmt.Errorf("CallOpts with Signer is required")
	}

	signerAddressStr, err := opts.Signer.GetAddress()
	if err != nil {
		return nil, fmt.Errorf("failed to get signer address: %w", err)
	}

	// normalize signer address
	signerAddress, err := bindutils.ConvertAddressToString(signerAddressStr)
	if err != nil {
		return nil, fmt.Errorf("invalid signer address %v: %w", signerAddressStr, err)
	}

	resolver := opts.ObjectResolver
	if resolver == nil {
		resolver = NewObjectResolver(c.client)
	}

	// resolve any UnresolvedObjects in EncodedCallArguments
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

	txData := ptb.Data
	if txData.V1 == nil || txData.V1.Kind == nil {
		return nil, fmt.Errorf("transaction data not properly initialized")
	}

	txBytes, err := txData.V1.Kind.Marshal()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transaction kind: %w", err)
	}

	devInspectResp, err := DevInspectTx(ctx, signerAddress, c.client, txBytes)
	if err != nil {
		return nil, err
	}

	if devInspectResp.Effects.Status.Status != "success" {
		return nil, fmt.Errorf("dev inspect failed for %s with status: %s, error: %v",
			encoded.String(), devInspectResp.Effects.Status.Status, devInspectResp.Effects.Status.Error)
	}

	if len(devInspectResp.Results) == 0 {
		// no return values
		return []any{}, nil
	}

	if string(devInspectResp.Results) == "null" {
		return []any{}, nil
	}

	if len(encoded.ReturnTypes) == 0 {
		return nil, fmt.Errorf("no return type information: %s", encoded.String())
	}

	// create TypeResolver to decode generic parameters
	var typeResolver *TypeResolver
	if encoded.TypeParams != nil && encoded.TypeArgs != nil {
		genericParams := encoded.TypeParams
		concreteTypes := make([]string, len(encoded.TypeArgs))
		for i, tag := range encoded.TypeArgs {
			concreteTypes[i] = typeTagToString(tag)
		}

		if len(genericParams) > 0 && len(concreteTypes) > 0 {
			typeResolver, err = NewTypeResolver(genericParams, concreteTypes)
			if err != nil {
				return nil, fmt.Errorf("failed to create type resolver: %w", err)
			}
		}
	}

	decodedValues, err := DecodeDevInspectResults(devInspectResp.Results, encoded.ReturnTypes, typeResolver)
	if err != nil {
		return nil, fmt.Errorf("failed to decode dev inspect results: %w", err)
	}

	return decodedValues, nil
}

func (c *BoundContract) ExecuteTransaction(ctx context.Context, opts *CallOpts, encoded *EncodedCall) (*models.SuiTransactionBlockResponse, error) {
	if opts == nil || opts.Signer == nil {
		return nil, fmt.Errorf("CallOpts with Signer is required")
	}

	ptb := transaction.NewTransaction()
	// Add the encoded call to the PTB
	_, err := c.AppendPTB(ctx, opts, ptb, encoded)
	if err != nil {
		return nil, fmt.Errorf("failed to add encoded call to PTB: %w", err)
	}

	return ExecutePTB(ctx, opts, c.client, ptb)
}

func GetObjectRef(ctx context.Context, client sui.ISuiAPI, objectID string) (*models.SuiObjectRef, error) {
	obj, err := ReadObject(ctx, objectID, client)
	if err != nil {
		return nil, err
	}
	if obj.Error != nil || obj.Data == nil || obj.Data.Content == nil {
		return nil, fmt.Errorf("failed to read object %s", objectID)
	}

	version, err := parseVersionString(obj.Data.Version)
	if err != nil {
		return nil, fmt.Errorf("failed to parse object version: %w", err)
	}

	return &models.SuiObjectRef{
		ObjectId: obj.Data.ObjectId,
		Version:  version,
		Digest:   obj.Data.Digest,
	}, nil
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

func parseVersionString(version string) (uint64, error) {
	if version == "" {
		return 0, fmt.Errorf("empty version string")
	}
	// version might be a number or a SequenceNumber type, try to parse as uint64 directly
	var v uint64
	_, err := fmt.Sscanf(version, "%d", &v)
	if err != nil {
		return 0, fmt.Errorf("failed to parse version %s: %w", version, err)
	}

	return v, nil
}

// AppendPTB adds an EncodedCall to an existing PTB and returns the result argument
func (c *BoundContract) AppendPTB(ctx context.Context, opts *CallOpts, ptb *transaction.Transaction, encoded *EncodedCall) (*transaction.Argument, error) {
	lggr, _ := logger.New()

	lggr.Info("APPENDING PTB FOR EXECUTE", opts.ObjectResolver)
	if opts.ObjectResolver == nil {
		opts.ObjectResolver = NewObjectResolver(c.client)
	}

	// resolve any UnresolvedObjects in EncodedCallArguments
	resolvedEncodedArgs := make([]*EncodedCallArgument, len(encoded.CallArgs))
	for i, encArg := range encoded.CallArgs {
		if encArg == nil {
			lggr.Info("ENCODED CALL ARGS EMPTY")
			return nil, fmt.Errorf("nil EncodedCallArgument at index %d", i)
		}

		if encArg.IsArgument() {
			resolvedEncodedArgs[i] = encArg
		} else if encArg.IsCallArg() {
			resolved, resolveErr := opts.ObjectResolver.ResolveCallArg(ctx, encArg.CallArg, encArg.TypeName)
			if resolveErr != nil {
				lggr.Info("FAILED TO RESOLVE CALLARG EXECUTE", resolveErr)
				return nil, fmt.Errorf("failed to resolve CallArg at index %d: %w", i, resolveErr)
			}
			resolvedEncodedArg := NewEncodedCallArgFromCallArgWithType(resolved, encArg.TypeName)
			resolvedEncodedArgs[i] = resolvedEncodedArg
		} else {
			return nil, errors.New("empty EncodedCallArgument")
		}
	}

	// Get existing inputs from PTB to enable proper deduplication across all calls
	var existingInputs []*transaction.CallArg
	if ptb.Data.V1 != nil && ptb.Data.V1.Kind != nil && ptb.Data.V1.Kind.ProgrammableTransaction != nil {
		existingInputs = ptb.Data.V1.Kind.ProgrammableTransaction.Inputs
	}

	callArgManager := NewCallArgManagerWithExisting(existingInputs)

	arguments, err := callArgManager.ConvertEncodedCallArgsToArguments(resolvedEncodedArgs)
	if err != nil {
		lggr.Info("FAILED TO CONVERT ENCODEDCALLARGS TO ARGS: ", err)
		return nil, fmt.Errorf("failed to convert EncodedCallArguments to Arguments: %w", err)
	}

	inputs := callArgManager.GetInputs()
	if ptb.Data.V1 == nil || ptb.Data.V1.Kind == nil || ptb.Data.V1.Kind.ProgrammableTransaction == nil {
		lggr.Info("FAILED TO CONVERT ENCODEDCALLARGS TO ARGS: ", err)
		return nil, errors.New("unexpected PTB with missing fields")
	}
	// Always replace inputs with deduplicated inputs (similar to BuildPTB)
	ptb.Data.V1.Kind.ProgrammableTransaction.Inputs = inputs

	// TODO: switch to non-pointer type in EncodedCall?
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

	lggr.Info("RUNNING MOVECALL EXECUTE")
	arg := ptb.MoveCall(
		models.SuiAddress(encoded.Module.PackageID),
		encoded.Module.ModuleName,
		encoded.Function,
		typeTagValues,
		argumentValues,
	)

	return &arg, nil
}

func ExecutePTB(ctx context.Context, opts *CallOpts, client sui.ISuiAPI, ptb *transaction.Transaction) (*models.SuiTransactionBlockResponse, error) {
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
		gasPrice, gasPriceErr := client.SuiXGetReferenceGasPrice(ctx)
		if gasPriceErr != nil {
			return nil, fmt.Errorf("failed to get reference gas price: %w", gasPriceErr)
		}
		ptb.SetGasPrice(gasPrice)
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
			gasRef, err = ToSuiObjectRef(ctx, client, opts.GasObject, signerAddress)
		} else {
			gasRef, err = FetchDefaultGasCoinRef(ctx, client, signerAddress)
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

	txBytes, err := ptb.Data.Marshal()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transaction: %w", err)
	}

	return SignAndSendTx(ctx, opts.Signer, client, txBytes, opts.WaitForExecution)
}
