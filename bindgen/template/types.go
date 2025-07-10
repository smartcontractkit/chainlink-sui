package template

import (
	"fmt"
	"strings"

	"github.com/smartcontractkit/chainlink-sui/bindgen/parse"
)

func createGoTypeFromMove(s string, localStructs map[string]parse.Struct) (tmplType, error) {
	switch s {
	case "u8":
		return tmplType{
			GoType:   "byte",
			MoveType: s,
		}, nil
	case "u16":
		return tmplType{
			GoType:   "uint16",
			MoveType: s,
		}, nil
	case "u32":
		return tmplType{
			GoType:   "uint32",
			MoveType: s,
		}, nil
	case "u64":
		return tmplType{
			GoType:   "uint64",
			MoveType: s,
		}, nil
	case "u128":
		return tmplType{
			GoType:   "*big.Int",
			MoveType: s,
		}, nil
	case "u256":
		return tmplType{
			GoType:   "*big.Int",
			MoveType: s,
		}, nil
	case "bool":
		return tmplType{
			GoType:   "bool",
			MoveType: s,
		}, nil
	case "address":
		return tmplType{
			GoType:   "string",
			MoveType: s,
		}, nil
	case "String", "string::String", "std::string::String", "0x1::string::String":
		return tmplType{
			GoType:   "string",
			MoveType: "0x1::string::String",
		}, nil
	case "UID", "object::UID", "sui::object::UID":
		return tmplType{
			GoType:   "string",
			MoveType: "sui::object::UID",
		}, nil
	default:
		if strings.HasPrefix(s, "vector<") && strings.HasSuffix(s, ">") {
			innerTypeName := strings.TrimSuffix(strings.TrimPrefix(s, "vector<"), ">")
			innerType, err := createGoTypeFromMove(innerTypeName, localStructs)
			if err != nil {
				return tmplType{}, err
			}

			return tmplType{
				GoType:   "[]" + innerType.GoType,
				MoveType: s,
			}, nil
		}

		optionPrefixes := []string{"Option<", "option::Option<", "std::option::Option<", "0x1::option::Option<"}
		for _, prefix := range optionPrefixes {
			if strings.HasPrefix(s, prefix) && strings.HasSuffix(s, ">") {
				innerTypeName := strings.TrimSuffix(strings.TrimPrefix(s, prefix), ">")
				innerType, err := createGoTypeFromMove(innerTypeName, localStructs)
				if err != nil {
					return tmplType{}, err
				}

				return tmplType{
					GoType:   "*" + innerType.GoType,
					MoveType: fmt.Sprintf("0x1::option::Option<%s>", innerType.MoveType),
					Option: &tmplOption{
						UnderlyingGoType: innerType.GoType,
					},
				}, nil
			}
		}

		baseType := stripGenericType(s)
		if _, ok := localStructs[baseType]; ok {
			if isSuiObjectStruct(localStructs[baseType]) {
				return tmplType{
					GoType:   "bind.Object",
					MoveType: s,
				}, nil
			}

			// if generic struct, use the base type name without generic parameters
			return tmplType{
				GoType:   baseType,
				MoveType: s,
			}, nil
		}

		return tmplType{
			GoType:   "bind.Object",
			MoveType: s,
		}, nil
	}
}

func isSuiObjectStruct(s parse.Struct) bool {
	if s.IsEvent {
		return false
	}

	for _, field := range s.Fields {
		if field.Name == "id" && field.Type == "UID" {
			return true
		}
	}

	return false
}

func stripGenericType(s string) string {
	if i := strings.Index(s, "<"); i != -1 {
		return s[:i]
	}

	return s
}

func containsGenericTypeParam(moveType string, typeParams []string) bool {
	// check if the type itself is a type parameter
	for _, param := range typeParams {
		if moveType == param {
			return true
		}
	}

	// check for generic in vectors and structs
	if idx := strings.Index(moveType, "<"); idx > 0 {
		typeParamsStr := moveType[idx+1 : len(moveType)-1]
		params := splitTypeParams(typeParamsStr)
		for _, param := range params {
			if containsGenericTypeParam(strings.TrimSpace(param), typeParams) {
				return true
			}
		}
	}

	return false
}

func splitTypeParams(params string) []string {
	var result []string
	var current strings.Builder
	depth := 0

	for _, ch := range params {
		if ch == '<' {
			depth++
		} else if ch == '>' {
			depth--
		} else if ch == ',' && depth == 0 {
			result = append(result, strings.TrimSpace(current.String()))
			current.Reset()

			continue
		}
		current.WriteRune(ch)
	}

	if current.Len() > 0 {
		result = append(result, strings.TrimSpace(current.String()))
	}

	return result
}
