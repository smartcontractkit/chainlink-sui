package ptb

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
) ([]transaction.TypeTag, error) {
	if len(params) == 0 {
		return []transaction.TypeTag{}, nil
	}

	builder := NewTypeTagBuilder()
	// Use a map to track unique type tags by their string representation
	uniqueTags := make(map[string]transaction.TypeTag)
	// Track the order in which unique keys are first encountered
	keyOrder := make([]string, 0)

	// Filter generic parameters
	for _, param := range params {
		if param.GenericType != nil {
			// Build the appropriate TypeTag
			var typeTag transaction.TypeTag
			var err error

			typeTag, err = builder.createTypeTag(*param.GenericType)
			if err != nil {
				return nil, fmt.Errorf("failed to create type tag for param %q with type %q: %w",
					param.Name, *param.GenericType, err)
			}

			// Only add to keyOrder if this is the first time we see this genericType
			if _, exists := uniqueTags[*param.GenericType]; !exists {
				keyOrder = append(keyOrder, *param.GenericType)
			}

			// Use the genericType as the deduplication key
			uniqueTags[*param.GenericType] = typeTag

			p.log.Debugw("Created type tag",
				"param", param.Name,
				"genericType", *param.GenericType,
				"isVector", isVectorType(*param.GenericType))
		}
	}

	if len(uniqueTags) == 0 {
		return []transaction.TypeTag{}, nil
	}

	// Return type tags in the order they first appeared
	result := make([]transaction.TypeTag, 0, len(uniqueTags))
	for _, k := range keyOrder {
		result = append(result, uniqueTags[k])
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
