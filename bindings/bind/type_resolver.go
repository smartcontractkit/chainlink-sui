package bind

import (
	"fmt"
	"strings"

	"github.com/smartcontractkit/chainlink-sui/bindings/utils"
)

type TypeResolver struct {
	typeParams map[string]string
}

func NewTypeResolver(genericParams []string, concreteTypes []string) (*TypeResolver, error) {
	if len(genericParams) != len(concreteTypes) {
		return nil, fmt.Errorf("mismatch between generic parameters (%d) and concrete types (%d)",
			len(genericParams), len(concreteTypes))
	}

	typeParams := make(map[string]string)
	for i, param := range genericParams {
		typeParams[param] = concreteTypes[i]
	}

	return &TypeResolver{
		typeParams: typeParams,
	}, nil
}

// ResolveType resolves a potentially generic type to a concrete type
func (r *TypeResolver) ResolveType(moveType string) string {
	// direct type parameter (T, U, etc.)
	if concrete, ok := r.typeParams[moveType]; ok {
		return concrete
	}

	// handle generic structs like Box<T> or Pair<T,U>
	if idx := strings.Index(moveType, "<"); idx > 0 {
		baseName := moveType[:idx]
		typeParamsStr := moveType[idx+1 : len(moveType)-1]

		params := utils.SplitTypeParams(typeParamsStr)
		resolvedParams := make([]string, len(params))

		for i, param := range params {
			param = strings.TrimSpace(param)
			if concrete, ok := r.typeParams[param]; ok {
				resolvedParams[i] = concrete
			} else {
				// not a type param, return as-is
				resolvedParams[i] = param
			}
		}

		return baseName + "<" + strings.Join(resolvedParams, ",") + ">"
	}

	return moveType
}

// IsGenericType checks if a type contains generic parameters
func IsGenericType(moveType string) bool {
	// Check for single letter type parameters
	if len(moveType) == 1 && moveType[0] >= 'A' && moveType[0] <= 'Z' {
		return true
	}

	// Check for generic in vector
	if strings.HasPrefix(moveType, "vector<") && strings.HasSuffix(moveType, ">") {
		inner := moveType[7 : len(moveType)-1]
		return IsGenericType(inner)
	}

	// Check for generic in struct types
	if idx := strings.Index(moveType, "<"); idx > 0 {
		typeParamsStr := moveType[idx+1 : len(moveType)-1]
		params := utils.SplitTypeParams(typeParamsStr)
		for _, param := range params {
			if IsGenericType(strings.TrimSpace(param)) {
				return true
			}
		}
	}

	return false
}
