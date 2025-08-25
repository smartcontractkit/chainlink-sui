package bind

import (
	"fmt"

	"github.com/block-vision/sui-go-sdk/transaction"

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
