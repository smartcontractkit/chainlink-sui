package offramp

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainwriter/config"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
)

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

const (
	LockOrBurn    = "lock_or_burn"
	ReleaseOrMint = "release_or_mint"
)

func GetTokenPoolPTBConfig(
	ctx context.Context,
	lggr logger.Logger,
	ptbClient client.SuiPTBClient,
	tokenPoolInfo TokenPool,
) (*config.ChainWriterPTBCommand, error) {
	normalizedModule, err := ptbClient.GetNormalizedModule(ctx, tokenPoolInfo.PackageId, tokenPoolInfo.ModuleId)
	if err != nil {
		lggr.Errorw("Error getting normalized module", "error", err)
		return nil, err
	}

	functions := normalizedModule.ExposedFunctions
	ptbConfig := config.ChainWriterPTBCommand{}

	if functions[tokenPoolInfo.Function] == nil {
		lggr.Errorw("Function not found", "function", tokenPoolInfo.Function)
		return nil, fmt.Errorf("function not found: %s", tokenPoolInfo.Function)
	} else {
		function := functions[tokenPoolInfo.Function].(map[string]any)
		decodedParameters, err := DecodeParameters(lggr, function)
		if err != nil {
			lggr.Errorw("Error decoding parameters", "error", err)
			return nil, err
		}
		isValid, _ := IsFunctionValid(lggr, decodedParameters, tokenPoolInfo.Function)
		if !isValid {
			// So the decoded parameters are not available for use in the PTB command. They are not needed.
			// Just log and return the error.
			lggr.Errorw("function is not valid", "function", function)
			return nil, fmt.Errorf("function is not valid: %s", tokenPoolInfo.Function)
		}

		lggr.Debugw("decodedParameters", "decodedParameters", decodedParameters)

		ptbParams := []codec.SuiFunctionParam{}
		for _, param := range decodedParameters {
			if param.Name == "TokenParams" && param.Module == "dynamic_dispatcher" {
				lggr.Debugw("Skipping out hot potato TokenParams", "param", param)
				continue
			}

			if param.Module == "tx_context" && param.Address == "0x2" {
				lggr.Debugw("Skipping out tx_context", "param", param)
				continue
			}

			isMutable := param.Reference == "MutableReference"
			ptbParams = append(ptbParams, codec.SuiFunctionParam{
				Name:      buildParameterName(param, tokenPoolInfo.Index),
				Type:      param.Type,
				Required:  true,
				IsMutable: &isMutable,
			})
		}

		ptbConfig = config.ChainWriterPTBCommand{
			Type:      codec.SuiPTBCommandMoveCall,
			PackageId: AnyPointer(tokenPoolInfo.PackageId),
			ModuleId:  AnyPointer(tokenPoolInfo.ModuleId),
			Function:  AnyPointer(tokenPoolInfo.Function),
			Params:    ptbParams,
		}
	}

	return &ptbConfig, nil
}

// getTokenPoolByTokenAddress gets token pool addresses for given token addresses (internal method)
func GetTokenPoolByTokenAddress(
	ctx context.Context,
	lggr logger.Logger,
	tokenAmounts []ccipocr3.RampTokenAmount,
	signerPublicKey []byte,
) ([]TokenPool, error) {
	coinMetadataAddresses := make([]string, len(tokenAmounts))
	for i, tokenAmount := range tokenAmounts {
		address := tokenAmount.DestTokenAddress
		coinMetadataAddresses[i] = "0x" + hex.EncodeToString(address)
	}

	lggr.Debugw("getting token pool infos",
		"packageID", s.AddressMappings["ccipPackageId"],
		"ccipObjectRef", s.AddressMappings["ccipObjectRef"],
		"coinMetadataAddresses", coinMetadataAddresses)

	signerAddress, err := client.GetAddressFromPublicKey(signerPublicKey)
	if err != nil {
		return nil, err
	}

	poolInfos, err := s.ptbClient.ReadFunction(
		ctx,
		signerAddress,
		s.AddressMappings["ccipPackageId"],
		"token_admin_registry",
		"get_pool_infos",
		[]any{
			s.AddressMappings["ccipObjectRef"],
			coinMetadataAddresses,
		},
		[]string{"object_id", "vector<address>"},
	)
	if err != nil {
		lggr.Errorw("Error getting pool infos", "error", err)
		return nil, err
	}

	var tokenPoolInfo GetPoolInfosResult
	lggr.Debugw("tokenPoolInfo", "tokenPoolInfo", poolInfos[0])
	jsonBytes, err := json.Marshal(poolInfos[0])
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jsonBytes, &tokenPoolInfo)
	if err != nil {
		return nil, err
	}

	lggr.Debugw("Decoded tokenPoolInfo", "tokenPoolInfo", tokenPoolInfo)

	tokenPools := make([]TokenPool, len(tokenAmounts))
	for i, tokenAmount := range tokenAmounts {
		lggr.Debugw("\n\nGetting pool address for token",
			"tokenAddress", tokenAmount.DestTokenAddress,
			"poolIndex", i)

		packageId := hex.EncodeToString(tokenPoolInfo.TokenPoolPackageIds[i][:])
		if !strings.HasPrefix(packageId, "0x") {
			packageId = "0x" + packageId
		}

		tokenType := tokenPoolInfo.TokenTypes[i]
		if !strings.HasPrefix(tokenType, "0x") {
			tokenType = "0x" + tokenType
		}

		tokenPoolStateAddress := hex.EncodeToString(tokenPoolInfo.TokenPoolStateAddresses[i][:])
		if !strings.HasPrefix(tokenPoolStateAddress, "0x") {
			tokenPoolStateAddress = "0x" + tokenPoolStateAddress
		}

		tokenPools[i] = TokenPool{
			CoinMetadata:          "0x" + hex.EncodeToString(tokenAmount.DestTokenAddress),
			TokenType:             tokenType,
			PackageId:             packageId,
			ModuleId:              tokenPoolInfo.TokenPoolModules[i],
			Function:              OFFRAMP_TOKEN_POOL_FUNCTION_NAME,
			TokenPoolStateAddress: tokenPoolStateAddress,
			Index:                 i,
		}
	}

	lggr.Debugw("tokenPoolInfo Decoded", "tokenPoolInfo", tokenPools)

	return tokenPools, nil
}

func buildParameterName(param SuiArgumentMetadata, tokenPoolIndex int) string {
	suffix := fmt.Sprintf("%s_%s", param.Module, param.Name)
	return fmt.Sprintf("token_pool_%d_%s", tokenPoolIndex, suffix)
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

	lggr.Debugw("Parsing parameter", "param", param)

	if str, ok := param.(string); ok {
		switch str {
		case "U8":
			return "int8"
		case "U16":
			return "int16"
		case "U32":
			return "int32"
		case "U64":
			return "int64"
		case "U128":
			return "int128"
		case "U256":
			return "int256"
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

func DecodeParameters(lggr logger.Logger, function map[string]any) ([]SuiArgumentMetadata, error) {
	parametersRaw, exists := function["parameters"]
	if !exists || parametersRaw == nil {
		lggr.Errorw("Parameters field is missing or nil", "function", function)
		return nil, fmt.Errorf("parameters field is missing or nil")
	}

	parameters, ok := parametersRaw.([]any)
	if !ok {
		lggr.Errorw("Parameters field is not an array", "parametersRaw", parametersRaw)
		return nil, fmt.Errorf("parameters field is not an array")
	}

	defaultReference := "Reference"
	decodedParameters := make([]SuiArgumentMetadata, len(parameters))
	for i, parameter := range parameters {
		decodedParameters[i] = decodeParam(lggr, parameter, defaultReference)
	}
	return decodedParameters, nil
}

func IsFunctionValid(lggr logger.Logger, decodedParameters []SuiArgumentMetadata, name string) (bool, []SuiArgumentMetadata) {

	if len(decodedParameters) < 3 {
		lggr.Errorw("Not enough parameters", "parameters", decodedParameters)
		return false, nil
	}

	// Decode and validate parameters
	param0 := decodedParameters[0]
	param1 := decodedParameters[1]
	param2 := decodedParameters[2]

	switch name {
	case LockOrBurn:
		if param0.Module != "state_object" || param0.Name != "CCIPObjectRef" {
			lggr.Errorw("CCIPObjectRef is not the first parameter", "module", param0.Module, "name", param0.Name)
			return false, nil
		}

		if param1.Module != "coin" || param1.Name != "Coin" {
			lggr.Errorw("Coin is not the second parameter", "module", param1.Module, "name", param1.Name)
			return false, nil
		}

		if param2.Module != "dynamic_dispatcher" || param2.Name != "TokenParams" {
			lggr.Errorw("Hot potato TokenParams is not the third parameter", "module", param2.Module, "name", param2.Name)
			return false, nil
		}
	case ReleaseOrMint:
		// TODO: Implement
		return false, nil
	default:
		lggr.Errorw("Invalid function name", "name", name)
		return false, nil
	}

	return true, decodedParameters
}
