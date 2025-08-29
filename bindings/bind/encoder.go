package bind

import (
	"context"
	"errors"
	"fmt"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/transaction"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	bindutils "github.com/smartcontractkit/chainlink-sui/bindings/utils"
)

// EncodedCall represents an encoded Move function call
type EncodedCall struct {
	Module      ModuleInformation
	Function    string
	TypeArgs    []*transaction.TypeTag
	CallArgs    []*EncodedCallArgument
	ReturnTypes []string
	TypeParams  []string
}

func (e *EncodedCall) String() string {
	// TODO: also add args
	return fmt.Sprintf("%s.%s", e.Module, e.Function)
}

// ValidateCallArgs ensures all CallArgs are valid (Pure or Object) and returns an error if an
// unresolved SDK variant is detected
func (e *EncodedCall) ValidateCallArgs() error {
	for i, encArg := range e.CallArgs {
		if encArg == nil {
			return fmt.Errorf("nil EncodedCallArgument at index %d", i)
		}

		if err := encArg.Validate(); err != nil {
			return fmt.Errorf("invalid EncodedCallArgument at index %d: %w", i, err)
		}

		if encArg.CallArg != nil {
			arg := encArg.CallArg

			if arg.UnresolvedObject != nil {
				return fmt.Errorf("UnresolvedObject detected at parameter %d: this is an SDK-only type not recognized by Sui blockchain. "+
					"Use bind.GetObject() to get fully resolved objects before encoding", i)
			}

			if arg.UnresolvedPure != nil {
				return fmt.Errorf("UnresolvedPure detected at parameter %d: this is an SDK-only type not recognized by Sui blockchain. "+
					"Pure values must be BCS-encoded as Pure([]byte)", i)
			}

			if arg.Pure == nil && arg.Object == nil {
				return fmt.Errorf("invalid CallArg at parameter %d: must be either Pure or Object (blockchain types), "+
					"not UnresolvedPure or UnresolvedObject (SDK-only types)", i)
			}
		}
	}

	return nil
}

// EncodeCallArgs encodes function parameters to CallArgs for use in PTB
func (c *BoundContract) EncodeCallArgs(function string, typeArgs []string, paramTypes []string, paramValues []any) (*EncodedCall, error) {
	return c.EncodeCallArgsWithReturnTypes(function, typeArgs, paramTypes, paramValues, nil)
}

// EncodeCallArgsWithReturnTypes encodes function parameters to CallArgs for use in PTB with return type information
func (c *BoundContract) EncodeCallArgsWithReturnTypes(function string, typeArgs []string, paramTypes []string, paramValues []any, returnTypes []string) (*EncodedCall, error) {
	return c.EncodeCallArgsWithGenerics(function, typeArgs, nil, paramTypes, paramValues, returnTypes)
}

// EncodeCallArgsWithGenerics encodes function parameters with full generic type information
func (c *BoundContract) EncodeCallArgsWithGenerics(function string, typeArgs []string, typeParams []string, paramTypes []string, paramValues []any, returnTypes []string) (*EncodedCall, error) {
	if len(paramTypes) != len(paramValues) {
		return nil, fmt.Errorf("paramTypes and paramValues must have the same length")
	}

	typeTags := make([]*transaction.TypeTag, len(typeArgs))
	for i, typeArg := range typeArgs {
		typeTag, err := bindutils.ConvertTypeStringToTypeTag(typeArg)
		if err != nil {
			return nil, fmt.Errorf("failed to parse type argument %q: %w", typeArg, err)
		}
		typeTags[i] = typeTag
	}

	// create generic type resolver if we have type parameters
	var resolver *GenericTypeResolver
	var err error
	if len(typeParams) > 0 && len(typeArgs) > 0 {
		resolver, err = NewGenericTypeResolver(typeParams, typeArgs)
		if err != nil {
			return nil, fmt.Errorf("failed to create generic type resolver: %w", err)
		}
	}

	encodedArgs := make([]*EncodedCallArgument, len(paramTypes))
	for i := range paramTypes {
		typeName := paramTypes[i]
		typeValue := paramValues[i]

		if resolver != nil {
			resolvedType := resolver.ResolveType(typeName)
			if resolvedType != typeName {
				// type was resolved, use the concrete type
				typeName = resolvedType
			}
		}

		// Check if the value is already a transaction.Argument (for PTB chaining)
		switch v := typeValue.(type) {
		case transaction.Argument:
			encodedArgs[i] = NewEncodedCallArgFromArgumentWithType(&v, typeName)
		case *transaction.Argument:
			encodedArgs[i] = NewEncodedCallArgFromArgumentWithType(v, typeName)
		default:
			if typeName == "vector<u8>" {
				switch v := typeValue.(type) {
				case [32]byte:
					typeValue = v[:]
				case *[32]byte:
					typeValue = v[:]
				}
			}

			callArg, err := ConvertToCallArg(typeName, typeValue)
			if err != nil {
				return nil, fmt.Errorf("failed to convert parameter %d (%s): %w", i, typeName, err)
			}
			encodedArgs[i] = NewEncodedCallArgFromCallArgWithType(callArg, typeName)
		}
	}

	return &EncodedCall{
		Module: ModuleInformation{
			PackageID:   c.packageID,
			PackageName: c.packageName,
			ModuleName:  c.moduleName,
		},
		Function:    function,
		TypeArgs:    typeTags,
		CallArgs:    encodedArgs,
		ReturnTypes: returnTypes,
		TypeParams:  typeParams,
	}, nil
}

// AppendPTB adds an EncodedCall to an existing PTB and returns the result argument
func (c *BoundContract) AppendPTB(ctx context.Context, opts *CallOpts, ptb *transaction.Transaction, encoded *EncodedCall) (*transaction.Argument, error) {
	lggr, err := logger.New()
	if err != nil {
		return nil, err
	}

	lggr.Infow(">>> ENTERED AppendPTB <<<",
		"numCallArgs", len(encoded.CallArgs),
		"hasPTBData", ptb.Data,
	)

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
