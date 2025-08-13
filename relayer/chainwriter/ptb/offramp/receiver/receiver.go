package receiver

import (
	"context"
	"fmt"
	"strings"

	"github.com/block-vision/sui-go-sdk/transaction"
	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	receiver_registry "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/receiver_registry"
	ptbClient "github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/signer"
)

var (
	SUI_PATH_COMPONENTS_COUNT = 3
	CCIP_RECEIVER_FUNCTION    = "ccip_receive"
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

func DecodeParameters(lggr logger.Logger, function map[string]any, key string) ([]SuiArgumentMetadata, error) {
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

	lggr.Infow("decoded parameters", "decodedParameters", decodedParameters)

	return decodedParameters, nil
}

func FilterRegisteredReceivers(
	ctx context.Context,
	lggr logger.Logger,
	messages []ccipocr3.Message,
	signerAddress string,
	client ptbClient.SuiPTBClient,
	ccipObjectRef string,
	ccipPackageId string,
) ([]ccipocr3.Message, error) {
	registeredReceivers := make([]ccipocr3.Message, 0)

	suiClient := client.GetClient()

	for _, message := range messages {
		if len(message.Receiver) > 0 && len(message.Data) > 0 {
			receiverParts := strings.Split(string(message.Receiver), "::")
			if len(receiverParts) != SUI_PATH_COMPONENTS_COUNT {
				return nil, fmt.Errorf("invalid receiver format, expected packageID:moduleID:functionName, got %s", message.Receiver)
			}

			receiverFactory, err := receiver_registry.NewReceiverRegistry(ccipPackageId, suiClient)
			if err != nil {
				return nil, err
			}

			receiverService := receiverFactory.DevInspect()

			devInspectSigner := signer.NewDevInspectSigner(signerAddress)

			opts := &bind.CallOpts{
				Signer:           devInspectSigner,
				WaitForExecution: true,
			}

			ref := bind.Object{
				Id: ccipObjectRef,
			}

			receiverPackageId := receiverParts[0]
			isRegistered, err := receiverService.IsRegisteredReceiver(ctx, opts, ref, receiverPackageId)
			if err != nil {
				lggr.Error("failed to check if receiver is registered", "error", err)
				return nil, err
			}

			if isRegistered {
				lggr.Info("receiver is registered ", "receiver ", message.Receiver)
				registeredReceivers = append(registeredReceivers, message)
			}
		}
	}

	return registeredReceivers, nil
}

func AddReceiverCallCommands(
	ctx context.Context,
	lggr logger.Logger,
	ptb *transaction.Transaction,
	signerAddress string,
	messages []ccipocr3.Message,
	previousCommandIndex uint16,
	ccipObjectRef string,
	ccipPackageId string,
	client ptbClient.SuiPTBClient,
) ([]*transaction.Argument, error) {
	suiClient := client.GetClient()
	devInspectSigner := signer.NewDevInspectSigner(signerAddress)
	registeredReceivers, err := FilterRegisteredReceivers(ctx, lggr, messages, signerAddress, client, ccipObjectRef, ccipPackageId)
	if err != nil {
		return nil, err
	}

	lggr.Info("registered receivers", "count", len(registeredReceivers))

	finalCommands := []*transaction.Argument{}

	for _, message := range registeredReceivers {
		receiverParts := strings.Split(string(message.Receiver), "::")
		receiverPackageId := receiverParts[0]

		receiverFactory, err := receiver_registry.NewReceiverRegistry(ccipPackageId, suiClient)
		if err != nil {
			return nil, err
		}

		moduleName, stateParams, err := getReceiverInfo(ctx, &devInspectSigner, receiverFactory, receiverPackageId, ccipObjectRef)
		if err != nil {
			return nil, err
		}

		lggr.Infow("receiver info", "receiver", receiverPackageId, "module", moduleName, "stateParams", stateParams)

		// TODO: remove, this is a debug function
		receiverFunction := "echo"

		signature, err := GetCCIPReceiverSignature(ctx, lggr, signerAddress, receiverPackageId, moduleName, client, receiverFunction)
		if err != nil {
			lggr.Error("failed to get receiver info", "error", err)
			return nil, err
		}

		decodedParameters, err := DecodeParameters(lggr, signature, "parameters")
		if err != nil {
			lggr.Error("failed to decode parameters", "error", err)
			return nil, err
		}

		receiverBoundContract, err := bind.NewBoundContract(
			receiverPackageId,
			receiverPackageId,
			moduleName,
			suiClient,
		)
		if err != nil {
			lggr.Error("failed to get receiver bound contract", "error", err)
			return nil, err
		}

		// Use the decoded parameters to populate the paramTypes for the bound contract
		paramTypes := []string{}
		for _, param := range decodedParameters {
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

		decodedReturnTypes, err := DecodeParameters(lggr, signature, "return")
		if err != nil {
			lggr.Error("failed to decode return types", "error", err)
			return nil, err
		}

		returnTypes := []string{}
		for _, param := range decodedReturnTypes {
			returnTypes = append(returnTypes, strings.ToLower(param.Type))
		}

		lggr.Infow("return types", "returnTypes", returnTypes)

		paramValues := []any{ccipObjectRef, []byte("Hello World")}

		encodedCall, err := receiverBoundContract.EncodeCallArgsWithGenerics(
			receiverFunction,
			[]string{},
			[]string{},
			paramTypes,
			paramValues,
			returnTypes,
		)
		if err != nil {
			lggr.Error("failed to encode call", "error", err)
			return nil, err
		}

		opts := &bind.CallOpts{
			Signer:           devInspectSigner,
			WaitForExecution: true,
		}

		arg, err := receiverBoundContract.AppendPTB(ctx, opts, ptb, encodedCall)
		if err != nil {
			return nil, err
		}
		finalCommands = append(finalCommands, arg)

		lggr.Infow("signature", "signature", signature)
	}

	return finalCommands, nil
}

func GetCCIPReceiverSignature(
	ctx context.Context,
	lggr logger.Logger,
	signerAddress string,
	receiverPackageId string,
	receiverModule string,
	client ptbClient.SuiPTBClient,
	receiverFunction string,
) (map[string]any, error) {
	lggr.Infow("getting ccip receiver signature", "receiverPackageId", receiverPackageId, "receiverModule", receiverModule)
	normalizedModule, err := client.GetNormalizedModule(ctx, receiverPackageId, receiverModule)
	if err != nil {
		lggr.Error("failed to get normalized module", "error", err)
		return nil, err
	}

	lggr.Infow("normalized module", "normalizedModule", normalizedModule)

	functions := normalizedModule.ExposedFunctions
	if functions[receiverFunction] == nil {
		lggr.Error("ccip_receive function not found", "receiverPackageId", receiverPackageId)
		return nil, fmt.Errorf("ccip_receive function not found: %s", receiverPackageId)
	}

	function := functions[receiverFunction].(map[string]any)

	return function, nil
}

func getReceiverInfo(
	ctx context.Context,
	devInspectSigner *signer.DevInspectSuiSigner,
	receiverFactory *receiver_registry.ReceiverRegistryContract,
	receiverPackageId string,
	ccipObjectRef string,
) (string, []string, error) {
	receiverService := receiverFactory.DevInspect()
	opts := &bind.CallOpts{
		Signer:           *devInspectSigner,
		WaitForExecution: true,
	}

	ref := bind.Object{
		Id: ccipObjectRef,
	}

	receiverInfo, err := receiverService.GetReceiverInfo(ctx, opts, ref, receiverPackageId)
	if err != nil {
		return "", nil, err
	}

	if len(receiverInfo) < 2 {
		return "", nil, fmt.Errorf("invalid receiver info response: expected 2 fields, got %d", len(receiverInfo))
	}

	moduleName, ok := receiverInfo[0].(string)
	if !ok {
		return "", nil, fmt.Errorf("invalid module name type: expected string, got %T", receiverInfo[0])
	}

	stateParams, ok := receiverInfo[1].([]string)
	if !ok {
		return "", nil, fmt.Errorf("invalid state params type: expected []string, got %T", receiverInfo[1])
	}

	return moduleName, stateParams, nil
}
