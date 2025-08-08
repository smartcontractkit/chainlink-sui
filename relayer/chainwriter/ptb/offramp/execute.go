// / A package to build all the bespoke code (PTB) along with its commands for the OffRampExecute operation.
// / There will be no dependency on the PTBConstructor interface here due to writing entirely custom code that is not meant to be re-usable.
// / This package does not generate CW configs but rather generates the actual PTB along with its commands directly.
package offramp

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/block-vision/sui-go-sdk/transaction"
	"github.com/mitchellh/mapstructure"
	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-sui/relayer/chainwriter/config"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
)

var (
	DEFAULT_NR_OFFRAMP_PTB_COMMANDS  = 2
	OFFRAMP_TOKEN_POOL_FUNCTION_NAME = "release_or_mint"
	SUI_PATH_COMPONENTS_COUNT        = 3
)

type TokenPool struct {
	CoinMetadata          string
	TokenType             string // e.g. "sui:0x66::link_module::LINK"
	PackageId             string
	ModuleId              string
	Function              string
	TokenPoolStateAddress string
	Index                 int
}

type SuiOffRampExecCallArgs struct {
	ReportContext [2][32]byte                `mapstructure:"ReportContext"`
	Report        []byte                     `mapstructure:"Report"`
	Info          ccipocr3.ExecuteReportInfo `mapstructure:"Info"`
}

// OffRampPTBArgs represents arguments for OffRamp PTB expansion operations
type OffRampPTBArgs struct {
	ExecArgs   SuiOffRampExecCallArgs
	PTBConfigs *config.ChainWriterFunction
}

// OffRampPTBResult represents the result of OffRamp PTB expansion operations
type OffRampPTBResult struct {
	PTBCommands []config.ChainWriterPTBCommand
	UpdatedArgs map[string]any
	TypeArgs    map[string]string
}

type GetPoolInfosResult struct {
	TokenPoolPackageIds     []SuiAddress `json:"token_pool_package_ids"`
	TokenPoolStateAddresses []SuiAddress `json:"token_pool_state_addresses"`
	TokenPoolModules        []string     `json:"token_pool_modules"`
	TokenTypes              []string     `json:"token_types"`
}

// BuildOffRampExecutePTB builds the PTB for the OffRampExecute operation
func BuildOffRampExecutePTB(
	ctx context.Context,
	lggr logger.Logger,
	args config.Arguments,
	ptbConfigs *config.ChainWriterFunction, // TODO: needed?
	signerPublicKey []byte,
	addressMappings OffRampAddressMappings,
) (ptb *transaction.Transaction, err error) {
	offrampArgs := &SuiOffRampExecCallArgs{}
	err = mapstructure.Decode(args.Args, &offrampArgs)
	if err != nil {
		return nil, fmt.Errorf("failed to decode args: %w", err)
	}

	tokenAmounts := make([]ccipocr3.RampTokenAmount, 0)
	messages := make([]ccipocr3.Message, 0)

	// save all messages in a single slice
	for _, report := range args.Info.AbstractReports {
		for _, message := range report.Messages {
			tokenAmounts = append(tokenAmounts, message.TokenAmounts...)
			messages = append(messages, message)
		}
	}

	tokenPoolStateAddresses, err := s.getTokenPoolByTokenAddress(ctx, lggr, tokenAmounts, signerPublicKey)
	if err != nil {
		return nil, nil, nil, err
	}

	generatedTokenPoolCommands, err := GeneratePTBCommandsForTokenPools(lggr, tokenPoolStateAddresses)
	if err != nil {
		return nil, nil, nil, err
	}

	// TODO: filter  out messages that have a receiver that is not registered

	// Generate receiver call commands
	//nolint:gosec // G115:
	receiverCommands, err := GenerateReceiverCallCommands(lggr, messages, uint16(len(generatedTokenPoolCommands)))
	if err != nil {
		return nil, nil, nil, err
	}

	// Construct the final PTB commands by inserting generated commands between config commands
	//nolint:gosec // G115:
	finalPTBCommands := make([]config.ChainWriterPTBCommand, 0, len(ptbConfigs.PTBCommands)+len(generatedTokenPoolCommands)+len(receiverCommands))

	// Add the first command from config (init_execute)
	finalPTBCommands = append(finalPTBCommands, ptbConfigs.PTBCommands[0])

	// Insert all generated token pool commands
	finalPTBCommands = append(finalPTBCommands, generatedTokenPoolCommands...)

	// Insert all generated receiver commands
	finalPTBCommands = append(finalPTBCommands, receiverCommands...)

	// Add the remaining commands from config (finish_execute)
	endCommand := ptbConfigs.PTBCommands[len(ptbConfigs.PTBCommands)-1]

	// Find and update the PTB dependency in the existing parameters
	for i := range endCommand.Params {
		if endCommand.Params[i].PTBDependency != nil {
			//nolint:gosec // G115: PTB commands are typically small in number, overflow extremely unlikely
			endCommand.Params[i].PTBDependency.CommandIndex = uint16(len(finalPTBCommands) - 1)
		}
	}

	finalPTBCommands = append(finalPTBCommands, endCommand)

	// Generate token pool arguments
	tokenPoolArgs, typeArgs, err := GenerateArgumentsForTokenPools(s.AddressMappings["ccipObjectRef"], s.AddressMappings["clockObject"], lggr, tokenPoolStateAddresses)
	if err != nil {
		return nil, nil, nil, err
	}

	filteredMessages, err := s.FilterRegisteredReceivers(ctx, lggr, messages, signerPublicKey)
	if err != nil {
		return nil, nil, nil, err
	}

	// Generate receiver call arguments
	//nolint:gosec // G115:
	receiverArgs, err := GenerateReceiverCallArguments(lggr, filteredMessages, uint16(len(generatedTokenPoolCommands)), s.AddressMappings["ccipObjectRef"])
	if err != nil {
		return nil, nil, nil, err
	}

	// Merge token pool and receiver arguments
	ptbArguments := make(map[string]any)
	for k, v := range tokenPoolArgs {
		ptbArguments[k] = v
	}
	for k, v := range receiverArgs {
		ptbArguments[k] = v
	}

	return finalPTBCommands, ptbArguments, typeArgs, nil
}

// getTokenPoolByTokenAddress gets token pool addresses for given token addresses (internal method)
func (s *SuiPTBExpander) getTokenPoolByTokenAddress(
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

func GenerateReceiverCallCommands(
	lggr logger.Logger,
	messages []ccipocr3.Message,
	previousCommandIndex uint16,
) ([]config.ChainWriterPTBCommand, error) {
	var receiverCommands []config.ChainWriterPTBCommand
	receiverIndex := previousCommandIndex + 1
	for _, message := range messages {
		if len(message.Receiver) > 0 && len(message.Data) > 0 {
			// Parse the receiver string into packageID:moduleID:functionName format
			receiverParts := strings.Split(string(message.Receiver), "::")
			if len(receiverParts) != SUI_PATH_COMPONENTS_COUNT {
				return nil, fmt.Errorf("invalid receiver format, expected packageID:moduleID:functionName, got %s", message.Receiver)
			}

			receiverCommands = append(receiverCommands, config.ChainWriterPTBCommand{
				Type:      codec.SuiPTBCommandMoveCall,
				PackageId: AnyPointer(receiverParts[0]),
				ModuleId:  AnyPointer(receiverParts[1]),
				Function:  AnyPointer(receiverParts[2]),
				Params: []codec.SuiFunctionParam{
					{
						Name:     "ccip_object_ref",
						Type:     "object_id",
						Required: true,
					},
					{
						Name:     fmt.Sprintf("package_id_%d", receiverIndex),
						Type:     "address",
						Required: true,
					},
					{
						Name:     fmt.Sprintf("receiver_params_%d", receiverIndex),
						Type:     "ptb_dependency",
						Required: true,
						PTBDependency: &codec.PTBCommandDependency{
							// PTB commands are typically small in number, overflow extremely unlikely
							//nolint:gosec
							CommandIndex: receiverIndex - 1,
						},
					},
				},
			})
			receiverIndex++
		}
	}

	return receiverCommands, nil
}

func (s *SuiPTBExpander) FilterRegisteredReceivers(
	ctx context.Context,
	lggr logger.Logger,
	messages []ccipocr3.Message,
	signerPublicKey []byte,
) ([]ccipocr3.Message, error) {
	registeredReceivers := make([]ccipocr3.Message, 0)
	for _, message := range messages {
		if len(message.Receiver) > 0 && len(message.Data) > 0 {
			receiverParts := strings.Split(string(message.Receiver), "::")
			if len(receiverParts) != SUI_PATH_COMPONENTS_COUNT {
				return nil, fmt.Errorf("invalid receiver format, expected packageID:moduleID:functionName, got %s", message.Receiver)
			}

			receiverPackageId := receiverParts[0]

			signerAddress, err := client.GetAddressFromPublicKey(signerPublicKey)
			if err != nil {
				return nil, err
			}

			lggr.Debugw("Getting receiver config", "receiverPackageId", receiverPackageId, "receiverFunctionName", receiverParts[2])

			result, err := s.ptbClient.ReadFunction(
				ctx,
				signerAddress,
				s.AddressMappings["ccipPackageId"],
				"receiver_registry",
				"is_registered_receiver",
				[]any{
					s.AddressMappings["ccipObjectRef"],
					receiverPackageId,
				},
				[]string{
					"object_id",
					"address",
				},
			)
			if err != nil {
				lggr.Errorw("Error getting pool infos", "error", err)
				return nil, err
			}

			var isRegistered bool
			lggr.Debugw("isRegistered", "isRegistered", result[0])
			err = codec.DecodeSuiJsonValue(result[0], &isRegistered)
			if err != nil {
				return nil, err
			}

			if isRegistered {
				registeredReceivers = append(registeredReceivers, message)
			}
		}
	}

	lggr.Debugw("registeredReceivers", "registeredReceivers", registeredReceivers)

	return registeredReceivers, nil
}

func GenerateReceiverCallArguments(
	lggr logger.Logger,
	messages []ccipocr3.Message,
	previousCommandIndex uint16,
	ccipObjectRef string,
) (map[string]any, error) {
	arguments := make(map[string]any)

	arguments["ccip_object_ref"] = ccipObjectRef

	commandIndex := previousCommandIndex + 1

	for _, message := range messages {
		if len(message.Receiver) > 0 && len(message.Data) > 0 {
			lggr.Debugw("receiverParts", "receiverParts", message.Receiver)
			receiverParts := strings.Split(string(message.Receiver), "::")
			if len(receiverParts) != SUI_PATH_COMPONENTS_COUNT {
				return nil, fmt.Errorf("invalid receiver format, expected packageID:moduleID:functionName, got %s", message.Receiver)
			}
			arguments[fmt.Sprintf("package_id_%d", commandIndex)] = receiverParts[0]
			commandIndex++
		}
	}

	return arguments, nil
}

// Auxiliary functions

// GeneratePTBCommands generates PTB commands for token addresses
func GeneratePTBCommandsForTokenPools(
	lggr logger.Logger,
	tokenPools []TokenPool,
) ([]config.ChainWriterPTBCommand, error) {
	ptbCommands := make([]config.ChainWriterPTBCommand, len(tokenPools))
	for i, tokenPool := range tokenPools {
		cmdIndex := i + 1
		previousCommandIndex := i
		// We need to increment the index by 1 because the first command (`init_execute`) is already added
		lggr.Infow("Generating PTB command from token pool", "tokenPool", tokenPool, "with index", i)
		packageID := tokenPool.PackageId
		if !strings.HasPrefix(packageID, "0x") {
			packageID = "0x" + packageID
		}
		ptbCommands[i] = config.ChainWriterPTBCommand{
			Type:      codec.SuiPTBCommandMoveCall,
			PackageId: AnyPointer(packageID),
			ModuleId:  AnyPointer(tokenPool.ModuleId),
			Function:  AnyPointer(tokenPool.Function),
			Params: []codec.SuiFunctionParam{
				{
					Name:      "ccip_object_ref",
					Type:      "object_id",
					Required:  true,
					IsMutable: AnyPointer(false),
				},
				{
					Name:     "receiver_params",
					Type:     "ptb_dependency",
					Required: true,
					PTBDependency: &codec.PTBCommandDependency{
						//nolint:gosec // G115: PTB commands are typically small in number, overflow extremely unlikely
						CommandIndex: uint16(previousCommandIndex),
					},
				},
				{
					Name:     fmt.Sprintf("index_%d", cmdIndex),
					Type:     "u64",
					Required: true,
				},
				{
					Name:      fmt.Sprintf("pool_%d", cmdIndex),
					Type:      "object_id",
					Required:  true,
					IsMutable: AnyPointer(true),
					IsGeneric: true,
				},
				{
					Name:      "clock",
					Type:      "object_id",
					Required:  true,
					IsMutable: AnyPointer(false),
				},
			},
		}
	}

	return ptbCommands, nil
}

// GenerateArgumentsFromTokenAmounts generates PTB arguments for token addresses
//
//nolint:gosec // G115:
func GenerateArgumentsForTokenPools(
	ccipObjectRef string,
	clockRef string,
	lggr logger.Logger,
	tokenPools []TokenPool,
) (map[string]any, map[string]string, error) {
	arguments := make(map[string]any)
	typeArgs := make(map[string]string)

	arguments["ccip_object_ref"] = ccipObjectRef
	arguments["clock"] = clockRef

	for i, tokenPool := range tokenPools {
		cmdIndex := i + 1
		arguments[fmt.Sprintf("pool_%d", cmdIndex)] = tokenPool.TokenPoolStateAddress
		arguments[fmt.Sprintf("index_%d", cmdIndex)] = uint64(tokenPool.Index)
		typeArgs[fmt.Sprintf("pool_%d", cmdIndex)] = tokenPool.TokenType
	}

	return arguments, typeArgs, nil
}

func AnyPointer[T any](v T) *T {
	return &v
}
