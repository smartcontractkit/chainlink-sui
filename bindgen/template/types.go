package template

import (
	"fmt"
	"strings"

	"github.com/smartcontractkit/chainlink-sui/bindgen/parse"
)

func createGoTypeFromMove(s string, localStructs map[string]parse.Struct, externalStructs []parse.ExternalStruct) (tmplType, error) {
	aliasMap := map[string]string{
		"dd::SourceTransferCap": "ccip::common::SourceTransferCap",
		"dd::TokenParams":       "ccip::common::TokenParams",
		"osh::DestTransferCap":  "ccip::common::DestTransferCap",
		"osh::ReceiverParams":   "ccip::common::ReceiverParams",
	}
	if realParam, ok := aliasMap[s]; ok {
		s = realParam
	}

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
			GoType:   "uint256.Int",
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
	case "String", "string::String", "std::string::String":
		return tmplType{
			GoType:   "string",
			MoveType: "0x1::string::String",
		}, nil
	case "UID":
		return tmplType{
			GoType:   "string",
			MoveType: "sui::object::UID",
		}, nil
	default:
		if isSuiObject(s) {
			return tmplType{
				GoType:   "bind.Object",
				MoveType: s,
			}, nil
		}
		if strings.HasPrefix(s, "vector<") && strings.HasSuffix(s, ">") {
			innerTypeName := strings.TrimSuffix(strings.TrimPrefix(s, "vector<"), ">")
			innerType, err := createGoTypeFromMove(innerTypeName, localStructs, externalStructs)
			if err != nil {
				return tmplType{}, err
			}

			return tmplType{
				GoType:   "[]" + innerType.GoType,
				MoveType: s,
			}, nil
		}
		if strings.HasPrefix(s, "Option<") && strings.HasSuffix(s, ">") {
			innerTypeName := strings.TrimSuffix(strings.TrimPrefix(s, "Option<"), ">")
			innerType, err := createGoTypeFromMove(innerTypeName, localStructs, externalStructs)
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
		if strings.HasPrefix(s, "option::Option<") && strings.HasSuffix(s, ">") {
			innerTypeName := strings.TrimSuffix(strings.TrimPrefix(s, "option::Option<"), ">")
			innerType, err := createGoTypeFromMove(innerTypeName, localStructs, externalStructs)
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
		if strings.HasPrefix(s, "std::option::Option<") && strings.HasSuffix(s, ">") {
			innerTypeName := strings.TrimSuffix(strings.TrimPrefix(s, "std::option::Option<"), ">")
			innerType, err := createGoTypeFromMove(innerTypeName, localStructs, externalStructs)
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
		// Check if local struct
		baseType := stripGenericType(s)
		if _, ok := localStructs[baseType]; ok {
			// If it's an object, we only want the object ID (string)
			if isSuiObjectStruct(localStructs[baseType]) {
				return tmplType{
					GoType:   "bind.Object",
					MoveType: s,
				}, nil
			}

			return tmplType{
				GoType:   s,
				MoveType: s,
			}, nil
		}
		// Check if external struct
		for _, externalStruct := range externalStructs {
			// Type could be used as package::module::Struct, module::Struct or Struct directly, depending on the import
			if s == fmt.Sprintf("%s::%s::%s", externalStruct.Package, externalStruct.Module, externalStruct.Name) ||
				s == fmt.Sprintf("%s::%s", externalStruct.Module, externalStruct.Name) ||
				s == externalStruct.Name {
				return tmplType{
					GoType:   fmt.Sprintf("module_%s.%s", externalStruct.Module, ToUpperCamelCase(externalStruct.Name)),
					MoveType: s,
					Import: &tmplImport{
						Path:        externalStruct.ImportPath,
						PackageName: fmt.Sprintf("module_%s", externalStruct.Module),
					},
				}, nil
			}
		}
	}

	return tmplType{}, fmt.Errorf("unknown move type: %s", s)
}

// Add as needed
var hardCodedObjectTypes = []string{
	"TreasuryCap",
	"CoinMetadata",
	"Clock",
	"Coin",
	"MintCap",
	"DenyCapV2",
	"DenyList",
	"TokenState",
	"ID",
}

func isSuiObject(s string) bool {
	for _, hardCodedType := range hardCodedObjectTypes {
		if strings.Contains(s, hardCodedType) {
			return true
		}
	}

	return false
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
