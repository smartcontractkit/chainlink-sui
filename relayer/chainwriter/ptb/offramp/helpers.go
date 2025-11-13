package offramp

import (
	"context"
	"fmt"
	"strings"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_token_admin_registry "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/token_admin_registry"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
)

func AnyPointer[T any](v T) *T {
	return &v
}

type OffRampAddressMappings struct {
	CcipPackageId    string `json:"ccipPackageId"`
	CcipObjectRef    string `json:"ccipObjectRef"`
	CcipOwnerCap     string `json:"ccipOwnerCap"`
	ClockObject      string `json:"clockObject"`
	OffRampPackageId string `json:"offRampPackageId"`
	OffRampState     string `json:"offRampState"`
}

// GetOfframpAddressMappings initializes and populates all required address mappings for PTB expansion operations.
//
// This function performs discovery and resolution of critical CCIP infrastructure addresses by:
// 1. Using the provided OffRamp package ID to query and discover the CCIP package ID
// 2. Reading owned objects to locate the OffRamp state pointer and extract the state address
// 3. Reading CCIP package objects to find the CCIP object reference and owner capability addresses
// 4. Assembling a complete address mapping required for subsequent PTB operations
//
// Parameters:
//   - ctx: Context for the operation, used for request lifecycle management
//   - lggr: Logger instance for debugging and operational visibility
//   - ptbClient: Sui PTB client for reading blockchain state and objects
//   - offRampPackageId: The OffRamp package identifier to start discovery from
//   - publicKey: Public key bytes for generating signer address for read operations
//
// Returns:
//   - OffRampAddressMappings: A struct containing all resolved addresses
//   - error: Error if any discovery step fails, objects are missing, or network issues occur
func GetOfframpAddressMappings(
	ctx context.Context,
	lggr logger.Logger,
	ptbClient client.SuiPTBClient,
	offRampPackageId string,
	publicKey []byte,
) (OffRampAddressMappings, error) {
	// address mappings for the expander
	addressMappings := OffRampAddressMappings{
		CcipPackageId:    "",
		CcipObjectRef:    "",
		CcipOwnerCap:     "",
		ClockObject:      "0x6",
		OffRampPackageId: offRampPackageId,
		OffRampState:     "",
	}

	// Use the `toAddress` (offramp package ID) from the config overrides to get the offramp pointer object
	signerAddress, err := client.GetAddressFromPublicKey(publicKey)
	if err != nil {
		lggr.Errorw("Error getting signer address", "error", err)
		return OffRampAddressMappings{}, err
	}

	ccipPkgID, err := ptbClient.GetCCIPPackageID(ctx, addressMappings.OffRampPackageId, signerAddress)
	if err != nil {
		return OffRampAddressMappings{}, err
	}

	addressMappings.CcipPackageId = ccipPkgID

	lggr.Debugw("ccipPackageId", "ccipPackageId", addressMappings.CcipPackageId)
	lggr.Debugw("offRampPackageId", "offrampPackageId", addressMappings.OffRampPackageId)

	// Get the offramp parent object ID and derive the OffRampState object ID
	offRampObjectIDKey := "off_ramp_object_id"
	var offRampObjectID string
	if cached, ok := ptbClient.GetCachedValue(offRampObjectIDKey); ok {
		offRampObjectID = cached.(string)
	} else {
		offRampObjectID, err = ptbClient.GetParentObjectID(ctx, addressMappings.OffRampPackageId, "offramp", "OffRampStatePointer")
		if err != nil {
			lggr.Errorw("Error getting offramp parent object ID", "error", err)
			return OffRampAddressMappings{}, err
		}
		ptbClient.SetCachedValue(offRampObjectIDKey, offRampObjectID)
	}

	// Derive OffRampState from parent object ID
	offRampStateID, err := bind.DeriveObjectIDWithVectorU8Key(offRampObjectID, []byte("OffRampState"))
	if err != nil {
		lggr.Errorw("Error deriving offramp state object ID", "error", err)
		return OffRampAddressMappings{}, err
	}
	addressMappings.OffRampState = offRampStateID

	// Get the CCIP parent object ID and derive CCIPObjectRef and OwnerCap
	ccipObjectIDKey := "ccip_object_id"
	var ccipObjectID string
	if cached, ok := ptbClient.GetCachedValue(ccipObjectIDKey); ok {
		ccipObjectID = cached.(string)
	} else {
		ccipObjectID, err = ptbClient.GetParentObjectID(ctx, ccipPkgID, "state_object", "CCIPObjectRefPointer")
		if err != nil {
			lggr.Errorw("Error getting ccip parent object ID", "error", err)
			return OffRampAddressMappings{}, err
		}
		ptbClient.SetCachedValue(ccipObjectIDKey, ccipObjectID)
	}

	// Derive CCIPObjectRef from parent object ID
	ccipObjectRefID, err := bind.DeriveObjectIDWithVectorU8Key(ccipObjectID, []byte("CCIPObjectRef"))
	if err != nil {
		lggr.Errorw("Error deriving ccip object ref ID", "error", err)
		return OffRampAddressMappings{}, err
	}
	addressMappings.CcipObjectRef = ccipObjectRefID

	// Derive CCIP OwnerCap from parent object ID
	ccipOwnerCapID, err := bind.DeriveObjectIDWithVectorU8Key(ccipObjectID, []byte("CCIP_OWNABLE"))
	if err != nil {
		lggr.Errorw("Error deriving ccip owner cap ID", "error", err)
		return OffRampAddressMappings{}, err
	}
	addressMappings.CcipOwnerCap = ccipOwnerCapID

	return addressMappings, nil
}

type TypeParameter struct {
	TypeParameter float64 `json:"TypeParameter"`
}

type SuiArgumentMetadata struct {
	Address       string          `json:"address"`
	Module        string          `json:"module"`
	Name          string          `json:"name"`
	TypeArguments []TypeParameter `json:"typeArguments"`
	Reference     string          `json:"reference"`
	Type          string          `json:"type"`
}

func decodeParam(lggr logger.Logger, param any, reference string) SuiArgumentMetadata {
	// Handle primitive types (strings like "U64", "Bool", etc.)
	if str, ok := param.(string); ok {
		return SuiArgumentMetadata{
			Address:       "",
			Module:        "",
			Name:          str,
			Reference:     reference,
			TypeArguments: []TypeParameter{},
			Type:          ParseParamType(lggr, str),
		}
	}

	// Handle complex types (maps)
	m := param.(map[string]any)
	for k, v := range m {
		switch k {
		case "Struct":
			// Direct struct
			s := v.(map[string]any)
			typeArguments := []TypeParameter{}
			for _, ta := range s["typeArguments"].([]any) {
				typeArgument := ta.(map[string]any)
				typeArguments = append(typeArguments, TypeParameter{TypeParameter: typeArgument["TypeParameter"].(float64)})
			}
			return SuiArgumentMetadata{
				Address:       s["address"].(string),
				Module:        s["module"].(string),
				Name:          s["name"].(string),
				Reference:     reference,
				TypeArguments: typeArguments,
				Type:          ParseParamType(lggr, v),
			}
		case "Reference", "MutableReference", "Vector":
			// Reference and MutableReference are the same thing
			// We need to unwrap the struct
			return decodeParam(lggr, v, k)
		default:
			inner := v.(map[string]any)["Struct"].(map[string]any)
			typeArguments := []TypeParameter{}
			for _, ta := range inner["typeArguments"].([]any) {
				typeArgument := ta.(map[string]any)
				typeArguments = append(typeArguments, TypeParameter{TypeParameter: typeArgument["TypeParameter"].(float64)})
			}
			return SuiArgumentMetadata{
				Address:       inner["address"].(string),
				Module:        inner["module"].(string),
				Name:          inner["name"].(string),
				Reference:     k,
				TypeArguments: typeArguments,
				Type:          ParseParamType(lggr, v),
			}
		}
	}
	return SuiArgumentMetadata{}
}

func ParseParamType(lggr logger.Logger, param interface{}) string {
	// Case 1: string primitive
	if str, ok := param.(string); ok {
		switch str {
		case "U8":
			return "u8"
		case "U16":
			return "u16"
		case "U32":
			return "u32"
		case "U64":
			return "u64"
		case "U128":
			return "u128"
		case "U256":
			return "u256"
		case "Bool":
			return "bool"
		case "Address":
			return "object_id"
		default:
			return "unknown"
		}
	}

	// Case 2: map structure (e.g., Vector, Reference, Struct)
	if m, ok := param.(map[string]interface{}); ok {
		if vectorVal, ok := m["Vector"]; ok {
			return "vector<" + ParseParamType(lggr, vectorVal) + ">"
		}
		if refVal, ok := m["Reference"]; ok {
			return ParseParamType(lggr, refVal)
		}
		if mutRefVal, ok := m["MutableReference"]; ok {
			return ParseParamType(lggr, mutRefVal)
		}
		if _, ok := m["Struct"]; ok {
			// Special case for strings
			if m["address"] == "String" {
				return "string"
			}
			return "object_id"
		}
		// Handle direct struct content (when called from decodeParam with unwrapped struct)
		if address, ok := m["address"]; ok {
			if address == "String" {
				return "string"
			}
			return "object_id"
		}
	}

	// Fallback
	return "unknown"
}

func DecodeParameters(lggr logger.Logger, function map[string]any, key string) ([]string, error) {
	parametersRaw, exists := function[key]
	if !exists || parametersRaw == nil {
		lggr.Errorw("key field is missing or nil", "function", function, "key", key)
		return nil, fmt.Errorf("key field is missing or nil")
	}

	parameters, ok := parametersRaw.([]any)
	if !ok {
		lggr.Errorw("key field is not an array", "parametersRaw", parametersRaw, "key", key)
		return nil, fmt.Errorf("key field is not an array")
	}

	lggr.Debugw("Raw parameters", "parameters", parameters, "key", key)

	defaultReference := "Reference"
	decodedParameters := make([]SuiArgumentMetadata, len(parameters))
	for i, parameter := range parameters {
		decodedParameters[i] = decodeParam(lggr, parameter, defaultReference)
	}

	lggr.Debugw("decoded parameters", "decodedParameters", decodedParameters)

	paramTypes := make([]string, 0, len(decodedParameters))
	for _, param := range decodedParameters {
		if param.Name == "TxContext" {
			continue
		}

		if param.Reference == "Reference" {
			if strings.HasPrefix(param.Type, "u") || param.Type == "bool" {
				// It's a primitive, not an object reference
				paramTypes = append(paramTypes, param.Type)
			} else {
				paramTypes = append(paramTypes, "&object")
			}
			continue
		}

		if param.Reference == "MutableReference" {
			paramTypes = append(paramTypes, "&mut object")
			continue
		}

		if param.Reference == "Vector" {
			paramTypes = append(paramTypes, "vector<"+param.Type+">")
			continue
		}

		paramTypes = append(paramTypes, strings.ToLower(param.Type))
	}

	return paramTypes, nil
}

// IsValidTokenPoolConfig does a basic check to ensure that the token pool config contains enough
// data to be used for offramp execution
func IsValidTokenPoolConfig(tokenConfig *module_token_admin_registry.TokenConfig) bool {
	return tokenConfig.TokenPoolPackageId != "" &&
		tokenConfig.TokenPoolModule != "" &&
		tokenConfig.TokenType != ""
}
