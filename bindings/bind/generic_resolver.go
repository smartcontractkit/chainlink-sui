package bind

import (
	"fmt"
	"strings"

	"github.com/smartcontractkit/chainlink-sui/bindings/utils"
)

// GenericTypeResolver resolves generic type parameters to concrete types
type GenericTypeResolver struct {
	typeParams map[string]string
}

// NewGenericTypeResolver creates a new resolver with the given type parameter mappings
func NewGenericTypeResolver(typeParamNames []string, typeArgs []string) (*GenericTypeResolver, error) {
	if len(typeParamNames) != len(typeArgs) {
		return nil, fmt.Errorf("type parameter count mismatch: %d params but %d args", len(typeParamNames), len(typeArgs))
	}

	typeParams := make(map[string]string)
	for i, name := range typeParamNames {
		typeParams[name] = typeArgs[i]
	}

	return &GenericTypeResolver{
		typeParams: typeParams,
	}, nil
}

// ResolveType resolves a type that may contain generic parameters
func (r *GenericTypeResolver) ResolveType(typeName string) string {
	if concrete, ok := r.typeParams[typeName]; ok {
		return concrete
	}

	// handle references
	if strings.HasPrefix(typeName, "&mut ") {
		inner := strings.TrimPrefix(typeName, "&mut ")
		return "&mut " + r.ResolveType(inner)
	}
	if strings.HasPrefix(typeName, "&") {
		inner := strings.TrimPrefix(typeName, "&")
		return "&" + r.ResolveType(inner)
	}

	// handle vectors
	if strings.HasPrefix(typeName, "vector<") && strings.HasSuffix(typeName, ">") {
		inner := typeName[7 : len(typeName)-1]
		return "vector<" + r.ResolveType(inner) + ">"
	}

	// handle generic structs
	if idx := strings.Index(typeName, "<"); idx > 0 && strings.HasSuffix(typeName, ">") {
		baseName := typeName[:idx]
		typeParamsStr := typeName[idx+1 : len(typeName)-1]
		params := utils.SplitTypeParams(typeParamsStr)

		resolvedParams := make([]string, len(params))
		for i, param := range params {
			resolvedParams[i] = r.ResolveType(strings.TrimSpace(param))
		}

		return baseName + "<" + strings.Join(resolvedParams, ", ") + ">"
	}

	return typeName
}
