package expander

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainwriter/config"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
)

type TypeParameter struct {
	TypeParameter int `json:"TypeParameter"`
}

type SuiArgumentMetadata struct {
	Address       string          `json:"address"`
	Module        string          `json:"module"`
	Name          string          `json:"name"`
	TypeArguments []TypeParameter `json:"typeArguments"`
	Reference     string          `json:"reference"`
}

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
		function := functions[tokenPoolInfo.Function]
		isValid, decodedParameters := isFunctionValid(lggr, function.(map[string]any))
		if !isValid {
			// So the decoded parameters are not available for use in the PTB command. They are not needed.
			// Just log and return the error.
			lggr.Errorw("function is not valid", "function", function)
			return nil, fmt.Errorf("function is not valid: %s", tokenPoolInfo.Function)
		}

		ptbParams := []codec.SuiFunctionParam{}
		for _, param := range decodedParameters {
			if param.Name == "TokenParams" && param.Module == "dynamic_dispatcher" {
				lggr.Debugw("Skipping out hot potato TokenParams", "param", param)
				continue
			}
			isMutable := param.Reference == "MutableReference"
			ptbParams = append(ptbParams, codec.SuiFunctionParam{
				Name:      buildParameterName(param, tokenPoolInfo.Index),
				Type:      param.Module + "::" + param.Name,
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

func buildParameterName(param SuiArgumentMetadata, tokenPoolIndex int) string {
	suffix := fmt.Sprintf("%s_%s", param.Module, param.Name)
	return fmt.Sprintf("token_pool_%d_%s", tokenPoolIndex, suffix)
}

func decodeParam(param any) SuiArgumentMetadata {
	m := param.(map[string]any)
	for k, v := range m {
		if k == "Struct" {
			// Direct struct
			s := v.(map[string]any)
			return SuiArgumentMetadata{
				Address:   s["address"].(string),
				Module:    s["module"].(string),
				Name:      s["name"].(string),
				Reference: "",
			}
		} else {
			// Note, we do not support Vector, so this logic holds
			// Adding support for Vectors
			// Wrapped (Reference/MutableReference/etc) - unwrap once
			inner := v.(map[string]any)["Struct"].(map[string]any)
			return SuiArgumentMetadata{
				Address:   inner["address"].(string),
				Module:    inner["module"].(string),
				Name:      inner["name"].(string),
				Reference: k,
			}
		}
	}
	return SuiArgumentMetadata{}
}

func isFunctionValid(lggr logger.Logger, function map[string]any) (bool, []SuiArgumentMetadata) {
	parameters := function["parameters"].([]any)
	lggr.Debugw("parameters", "parameters", parameters)

	decodedParameters := make([]SuiArgumentMetadata, len(parameters))
	for i, parameter := range parameters {
		decodedParameters[i] = decodeParam(parameter)
	}

	if len(parameters) < 3 {
		lggr.Errorw("Not enough parameters", "parameters", parameters)
		return false, nil
	}

	// Decode and validate parameters
	param0 := decodedParameters[0]
	param1 := decodedParameters[1]
	param2 := decodedParameters[2]

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

	return true, decodedParameters
}
