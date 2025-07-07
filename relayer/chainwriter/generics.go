package chainwriter

import (
	"fmt"
	"strings"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/transaction"

	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
)

// TypeTagBuilder provides methods to create TypeTag instances
type TypeTagBuilder struct{}

// NewTypeTagBuilder creates a new TypeTagBuilder instance
func NewTypeTagBuilder() *TypeTagBuilder {
	return &TypeTagBuilder{}
}

func (p *PTBConstructor) ResolveGenericTypeTags(
	params []codec.SuiFunctionParam,
	arguments Arguments,
) ([]transaction.TypeTag, error) {
	if len(params) == 0 {
		return []transaction.TypeTag{}, nil
	}

	// Filter generic parameters
	genericParams := make([]codec.SuiFunctionParam, 0)
	for _, param := range params {
		if param.IsGeneric {
			genericParams = append(genericParams, param)
		}
	}

	if len(genericParams) == 0 {
		return []transaction.TypeTag{}, nil
	}

	// Use a map to track unique type tags by their string representation
	uniqueTags := make(map[string]transaction.TypeTag)

	builder := NewTypeTagBuilder()

	p.log.Debugw("Resolving generic type tags",
		"genericParamCount", len(genericParams),
		"arguments", arguments)

	for _, param := range genericParams {
		if param.Name == "" {
			return nil, fmt.Errorf("generic parameter missing name: %+v", param)
		}

		// Get the actual generic type from ArgTypes
		genericType, exists := arguments.ArgTypes[param.Name]
		if !exists {
			return nil, fmt.Errorf("generic parameter %q not found in ArgTypes", param.Name)
		}

		// Build the appropriate TypeTag
		var typeTag transaction.TypeTag
		var err error

		typeTag, err = builder.createTypeTag(genericType)
		if err != nil {
			return nil, fmt.Errorf("failed to create type tag for param %q with type %q: %w",
				param.Name, genericType, err)
		}

		// Use the genericType as the deduplication key
		uniqueTags[genericType] = typeTag

		p.log.Debugw("Created type tag",
			"param", param.Name,
			"genericType", genericType,
			"isVector", isVectorType(param.Type))
	}

	// Convert map to slice maintaining order by re-iterating over genericParams
	result := make([]transaction.TypeTag, 0, len(uniqueTags))
	seen := make(map[string]bool)

	for _, param := range genericParams {
		if genericType, exists := arguments.ArgTypes[param.Name]; exists {
			if !seen[genericType] {
				if typeTag, exists := uniqueTags[genericType]; exists {
					result = append(result, typeTag)
					seen[genericType] = true
				}
			}
		}
	}

	p.log.Debugw("Resolved type tags", "count", len(result))

	return result, nil
}

// Helper function to check if a type is a vector type
func isVectorType(paramType string) bool {
	return strings.HasPrefix(paramType, "vector<") && strings.HasSuffix(paramType, ">")
}

// createTypeTag creates a TypeTag for the given type string
func (b *TypeTagBuilder) createTypeTag(typeStr string) (transaction.TypeTag, error) {
	if typeStr == "" {
		return transaction.TypeTag{}, fmt.Errorf("type string cannot be empty")
	}

	// Handle vector types
	if strings.HasPrefix(typeStr, "vector<") && strings.HasSuffix(typeStr, ">") {
		// in the case of vector<T>, we need to create a TypeTag for T only
		// actual vector<G> for the generic T is not supported and causes BCS errors when marshalling
		innerType := extractVectorInnerType(typeStr)
		if innerType == "" {
			return transaction.TypeTag{}, fmt.Errorf("failed to extract vector inner type from %s", typeStr)
		}

		typeTag, err := b.createTypeTag(innerType)
		if err != nil {
			return transaction.TypeTag{}, fmt.Errorf("failed to create type tag for vector inner type %s: %w", innerType, err)
		}

		return typeTag, nil
	}

	// Handle struct types (package::module::name)
	if strings.Contains(typeStr, "::") {
		return b.createStructTypeTag(typeStr)
	}

	// If is not a vector or struct, it is a primitive type
	// TODO: Block vision SDK does not support primitive types, so this call will return an error
	return b.createPrimitiveTypeTag(typeStr)
}

// createPrimitiveTypeTag creates TypeTag for primitive types
func (b *TypeTagBuilder) createPrimitiveTypeTag(typeStr string) (transaction.TypeTag, error) {
	baseTag := transaction.TypeTag{}
	return baseTag, fmt.Errorf("block vision SDK does not support primitive types: %s", typeStr)
}

// createStructTypeTag creates TypeTag for struct types
func (b *TypeTagBuilder) createStructTypeTag(typeStr string) (transaction.TypeTag, error) {
	parts := strings.Split(typeStr, "::")
	tagLen := 3
	if len(parts) != tagLen {
		return transaction.TypeTag{}, fmt.Errorf("invalid struct type format %q, expected package::module::name", typeStr)
	}

	packageID, module, name := parts[0], parts[1], parts[2]

	// Convert package ID to address bytes
	packageAddr := models.SuiAddress(packageID)
	addressBytes, err := transaction.ConvertSuiAddressStringToBytes(packageAddr)
	if err != nil {
		return transaction.TypeTag{}, fmt.Errorf("failed to convert package address %q: %w", packageID, err)
	}

	// Initialize all boolean fields to false for struct types
	return transaction.TypeTag{
		Struct: &transaction.StructTag{
			Address:    *addressBytes,
			Module:     module,
			Name:       name,
			TypeParams: []*transaction.TypeTag{},
		},
	}, nil
}

// extractVectorInnerType extracts the inner type from a vector type string
func extractVectorInnerType(vectorType string) string {
	if !strings.HasPrefix(vectorType, "vector<") || !strings.HasSuffix(vectorType, ">") {
		return ""
	}

	return strings.TrimSuffix(strings.TrimPrefix(vectorType, "vector<"), ">")
}
