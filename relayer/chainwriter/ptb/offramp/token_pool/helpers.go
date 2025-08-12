package token_pool

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/transaction"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cwConfig "github.com/smartcontractkit/chainlink-sui/relayer/chainwriter/config"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainwriter/ptb/generics"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
)

func encodeHexByteArray(bytes []byte) string {
	hexString := hex.EncodeToString(bytes)
	if !strings.HasPrefix(hexString, "0x") {
		hexString = "0x" + hexString
	}
	return hexString
}

func formatSuiObjectString(objectString string) string {
	if !strings.HasPrefix(objectString, "0x") {
		objectString = "0x" + objectString
	}
	return objectString
}

func DecodeBase64ParamsArray(base64Params []string) ([]string, error) {
	decodedParams := []string{}
	for _, param := range base64Params {
		decodedParam, err := base64.StdEncoding.DecodeString(param)
		if err != nil {
			return nil, fmt.Errorf("failed to decode base64 param: %w", err)
		}
		decodedParams = append(decodedParams, encodeHexByteArray(decodedParam))
	}
	return decodedParams, nil
}

func ProcessMoveCall(
	ctx context.Context,
	lggr logger.Logger,
	builder *transaction.Transaction,
	cmd cwConfig.ChainWriterPTBCommand,
	arguments cwConfig.Arguments,
	cachedArgs *map[string]transaction.Argument,
	ptbClient *client.PTBClient,
) (*transaction.Argument, error) {
	lggr.Debugw("Processing move call", "Command", cmd, "Args", arguments)

	// All three fields below are required for a successful move call
	if cmd.PackageId == nil {
		return nil, fmt.Errorf("missing required parameter 'PackageId' for move call PTB command")
	}
	if cmd.ModuleId == nil {
		return nil, fmt.Errorf("missing required parameter 'ModuleId' for move call PTB command")
	}
	if cmd.Function == nil {
		return nil, fmt.Errorf("missing required parameter 'Function' for move call PTB command")
	}

	// Convert package ID to Address
	packageId := models.SuiAddress(*cmd.PackageId)

	// Process arguments
	processedArgs, err := ProcessArgsForCommand(ctx, lggr, builder, cmd.Params, arguments, cachedArgs, ptbClient)
	if err != nil {
		return nil, err
	}

	processedArgTypes, err := generics.ResolveGenericTypeTags(lggr, cmd.Params, arguments)
	if err != nil {
		return nil, err
	}

	lggr.Debugw("Processed Type Tags", "Type Tags", processedArgTypes)
	lggr.Debugw("Processed args", "Args", processedArgs)
	// Add the move call to the builder
	ptbArgument := builder.MoveCall(packageId, *cmd.ModuleId, *cmd.Function, processedArgTypes, processedArgs)

	return &ptbArgument, nil
}

// ProcessArgsForCommand converts parametedsr specifications into concrete arguments
func ProcessArgsForCommand(
	ctx context.Context,
	lggr logger.Logger,
	builder *transaction.Transaction,
	params []codec.SuiFunctionParam,
	arguments cwConfig.Arguments,
	cachedArgs *map[string]transaction.Argument,
	ptbClient *client.PTBClient,
) ([]transaction.Argument, error) {
	processedArgs := make([]transaction.Argument, 0, len(params))
	lggr.Debugw("Processing args", "Args", arguments)
	for _, param := range params {
		lggr.Debugw("Processing PTB parameter", "Param", param)

		// specify if the value is Mutable, this is used specifically for object PTB args
		isMutable := true
		if param.IsMutable != nil {
			isMutable = *param.IsMutable
		}

		// check if this is a PTB result dependency
		if param.PTBDependency != nil {
			// if the config does not specify a ResultIndex, then the dependency is
			// on the entire result of the dependee command
			if param.PTBDependency.ResultIndex == nil {
				processedArgs = append(processedArgs, transaction.Argument{
					Result: &param.PTBDependency.CommandIndex,
				})

				continue
			}

			// otherwise, we need a specific result from the dependee command
			processedArgs = append(processedArgs, transaction.Argument{
				NestedResult: &transaction.NestedResult{
					Index:       param.PTBDependency.CommandIndex,
					ResultIndex: *param.PTBDependency.ResultIndex,
				},
			})

			continue
		}

		// otherwise, check if the parameter is in the provided args
		if argRawValue, exists := arguments.Args[param.Name]; exists {
			// check if the param has already been converted and cached
			if cachedArg, exists := (*cachedArgs)[param.Name]; exists {
				processedArgs = append(processedArgs, cachedArg)
				continue
			}

			if param.Type == "object_id" {
				id, ok := argRawValue.(string)
				if !ok {
					return nil, fmt.Errorf("expected string for object id for param %s, got %T", param.Name, argRawValue)
				}
				argRawValue = id
			}

			// append to the array of args
			processedArgValue, err := ptbClient.TransformTransactionArg(ctx, builder, argRawValue, param.Type, isMutable)
			if err != nil {
				return nil, fmt.Errorf("failed to build argument for %s: %w, %s", param.Name, err, argRawValue)
			}
			processedArgs = append(processedArgs, *processedArgValue)
			// add the processed arg to the cache
			(*cachedArgs)[param.Name] = *processedArgValue

			continue
		}

		// fallback to the default value. Some assumptions:
		// - arguments with default values do NOT have type arguments
		// - objects do not have default values
		if param.DefaultValue != nil {
			value := param.DefaultValue
			if param.Type == "object_id" {
				id, ok := param.DefaultValue.(string)
				if !ok {
					return nil, fmt.Errorf("expected string for object id for param %s, got %T", param.Name, param.DefaultValue)
				}
				value = id
			}
			ptbArg, err := ptbClient.TransformTransactionArg(ctx, builder, value, param.Type, isMutable)
			if err != nil {
				return nil, fmt.Errorf("failed to build default value for %s: %w", param.Name, err)
			}
			// append to the array of args
			processedArgs = append(processedArgs, *ptbArg)

			continue
		}

		// Value not found for required param
		if param.Required {
			return nil, fmt.Errorf("required parameter %s has no value", param.Name)
		}

		// append an empty argument since it is not required and no value found
		processedArgs = append(processedArgs, transaction.Argument{})
	}

	return processedArgs, nil
}
