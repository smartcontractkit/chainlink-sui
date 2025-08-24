package offramp

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
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
	getCCIPPackageIdResponse, err := ptbClient.ReadFunction(ctx, signerAddress, addressMappings.OffRampPackageId, "offramp", "get_ccip_package_id", []any{}, []string{})
	if err != nil {
		lggr.Errorw("Error reading ccip package id", "error", err)
		return OffRampAddressMappings{}, err
	}
	lggr.Debugw("getCCIPPackageIdResponse", "getCCIPPackageIdResponse", getCCIPPackageIdResponse)
	// Parse the response to get the returned address as a hex string
	var addressBytes []byte

	// Handle both byte slice and base64 string responses
	switch v := getCCIPPackageIdResponse[0].(type) {
	case []byte:
		// Response is already raw bytes ([]byte and []uint8 are the same type)
		addressBytes = v
	case string:
		// Response is base64-encoded string, decode it
		var decodeErr error
		addressBytes, decodeErr = base64.StdEncoding.DecodeString(v)
		if decodeErr != nil {
			lggr.Errorw("Error decoding base64 ccip package id", "error", decodeErr)
			return OffRampAddressMappings{}, decodeErr
		}
	default:
		lggr.Errorw("Unexpected type for ccip package id response", "type", fmt.Sprintf("%T", getCCIPPackageIdResponse[0]))
		return OffRampAddressMappings{}, fmt.Errorf("unexpected type for ccip package id response, got %T", getCCIPPackageIdResponse[0])
	}
	// Convert bytes to hex string with "0x" prefix
	ccipPackageId := "0x" + hex.EncodeToString(addressBytes)
	addressMappings.CcipPackageId = ccipPackageId

	lggr.Debugw("ccipPackageId", "ccipPackageId", addressMappings.CcipPackageId)
	lggr.Debugw("offRampPackageId", "offrampPackageId", addressMappings.OffRampPackageId)

	// get the offramp state object
	offrampOwnedObjects, err := ptbClient.ReadOwnedObjects(ctx, addressMappings.OffRampPackageId, nil)
	if err != nil {
		lggr.Errorw("Error reading offramp state object", "error", err)
		return OffRampAddressMappings{}, err
	}
	for _, ccipOwnedObject := range offrampOwnedObjects {
		if ccipOwnedObject.Data.Type != "" && strings.Contains(ccipOwnedObject.Data.Type, "offramp::OffRampStatePointer") {
			lggr.Debugw("Found offramp state object pointer", "fields", ccipOwnedObject.Data.Content.Fields)
			// parse the object into a map
			parsedObject := ccipOwnedObject.Data.Content.Fields
			lggr.Debugw("offRampStatePointer Parsed", "offRampStatePointer", parsedObject)
			addressMappings.OffRampState = parsedObject["off_ramp_state_id"].(string)

			break
		}
	}
	if addressMappings.OffRampState == "" {
		lggr.Errorw("Address mappings are not populated", "addressMappings", addressMappings)
		return OffRampAddressMappings{}, fmt.Errorf("address mappings are missing required fields for expander (offRampState)")
	}

	// Get the object pointer present in the CCIP package ID
	ccipOwnedObjects, err := ptbClient.ReadOwnedObjects(ctx, addressMappings.CcipPackageId, nil)
	if err != nil {
		lggr.Errorw("Error reading ccip object ref", "error", err)
		return OffRampAddressMappings{}, err
	}
	for _, ccipOwnedObject := range ccipOwnedObjects {
		if ccipOwnedObject.Data.Type != "" && strings.Contains(ccipOwnedObject.Data.Type, "state_object::CCIPObjectRefPointer") {
			// parse the object into a map
			parsedObject := ccipOwnedObject.Data.Content.Fields
			if err != nil {
				lggr.Errorw("Error parsing ccip object ref", "error", err)
				return OffRampAddressMappings{}, err
			}
			lggr.Debugw("ccipObjectRefPointer", "ccipObjectRefPointer", parsedObject)
			addressMappings.CcipObjectRef = parsedObject["object_ref_id"].(string)
			addressMappings.CcipOwnerCap = parsedObject["owner_cap_id"].(string)

			break
		}
	}
	// check that address mappings are populated
	if addressMappings.CcipObjectRef == "" || addressMappings.CcipOwnerCap == "" {
		lggr.Errorw("Address mappings are not populated", "addressMappings", addressMappings)
		return OffRampAddressMappings{}, fmt.Errorf("address mappings are missing required fields for expander (ccipObjectRef, ccipOwnerCap)")
	}

	lggr.Debugw("Address mappings for expander", "addressMappings", addressMappings)

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
			paramTypes = append(paramTypes, "&object")
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

func ConvertFunctionParams(argMap map[string]interface{}, params []codec.SuiFunctionParam) ([]string, []any, error) {
	var types []string
	var values []any

	for _, paramConfig := range params {
		argValue, ok := argMap[paramConfig.Name]
		if !ok {
			// If it's required and has no default, it's an error
			if paramConfig.Required {
				return nil, nil, fmt.Errorf("missing argument: %s", paramConfig.Name)
			}
			// If default is set, use it
			if paramConfig.DefaultValue != nil {
				argValue = paramConfig.DefaultValue
			} else {
				// Otherwise, skip this param — assume it will be appended later
				continue
			}
		}

		types = append(types, paramConfig.Type)
		values = append(values, argValue)
	}

	return types, values, nil
}
