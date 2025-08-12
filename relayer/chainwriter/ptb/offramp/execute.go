// / A package to build all the bespoke code (PTB) along with its commands for the OffRampExecute operation.
// / There will be no dependency on the PTBConstructor interface here due to writing entirely custom code that is not meant to be re-usable.
// / This package does not generate CW configs but rather generates the actual PTB along with its commands directly.
package offramp

import (
	"context"
	"fmt"
	"strings"

	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/block-vision/sui-go-sdk/transaction"
	"github.com/mitchellh/mapstructure"
	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_offramp "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_offramp/offramp"
	module_complex "github.com/smartcontractkit/chainlink-sui/bindings/generated/test/complex"
	module_counter "github.com/smartcontractkit/chainlink-sui/bindings/generated/test/counter"
	"github.com/smartcontractkit/chainlink-sui/bindings/packages/ccip"
	"github.com/smartcontractkit/chainlink-sui/bindings/packages/offramp"

	"github.com/smartcontractkit/chainlink-sui/relayer/chainwriter/config"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
)

var (
	DEFAULT_NR_OFFRAMP_PTB_COMMANDS  = 2
	OFFRAMP_TOKEN_POOL_FUNCTION_NAME = "release_or_mint"
	SUI_PATH_COMPONENTS_COUNT        = 3
)

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
	client sui.ISuiAPI,
	ptb *transaction.Transaction,
	args config.Arguments,
	ptbConfigs *config.ChainWriterFunction, // TODO: needed?
	signerPublicKey []byte,
	addressMappings OffRampAddressMappings,
) (err error) {
	offrampArgs := &SuiOffRampExecCallArgs{}
	err = mapstructure.Decode(args.Args, &offrampArgs)
	if err != nil {
		return fmt.Errorf("failed to decode args for offramp execute PTB: %w", err)
	}

	tokenAmounts := make([]ccipocr3.RampTokenAmount, 0)
	messages := make([]ccipocr3.Message, 0)

	// save all messages in a single slice
	for _, report := range offrampArgs.Info.AbstractReports {
		for _, message := range report.Messages {
			tokenAmounts = append(tokenAmounts, message.TokenAmounts...)
			messages = append(messages, message)
		}
	}

	tokenPoolStateAddresses, err := GetTokenPoolByTokenAddress(ctx, lggr, tokenAmounts, signerPublicKey)
	if err != nil {
		return fmt.Errorf("failed to get token pool by token address offramp execute PTB: %w", err)
	}

	// Set the offramp package interface from bindings
	offrampPkg, err := offramp.NewOfframp(addressMappings.OffRampPackageId, client)
	if err != nil {
		return err
	}
	offrampContract := offrampPkg.Offramp().(*module_offramp.OfframpContract)
	offrampEncoder := offrampContract.Encoder()

	// Create an encoder for the `init_execute` offramp method to be attached to the PTB.
	// This is being done using the bindings to re-use code but can otherwise be done using the SDK directly.
	encodedInitExecute, err := offrampEncoder.InitExecute(
		bind.Object{Id: addressMappings.CcipObjectRef}, // TODO: double check this
		bind.Object{Id: addressMappings.OffRampState},
		bind.Object{Id: addressMappings.ClockObject},
		[][]byte{
			offrampArgs.ReportContext[0][:],
			offrampArgs.ReportContext[1][:],
		},
		offrampArgs.Report,
	)
	if err != nil {
		return fmt.Errorf("failed to encode move call (init_execute) using bindings: %w", err)
	}

	initExecuteResult, err := offrampContract.BuildPTB(ctx, ptb, encodedInitExecute)
	if err != nil {
		return fmt.Errorf("failed to build PTB (init_execute) using bindings: %w", err)
	}

	// Generate N token pool commands and attach them to the PTB, each command must return a result
	// that will subsequently be used to make a vector of hot potatoes before finishing execution.
	generatedTokenPoolCommands, err := GeneratePTBCommandsForTokenPools(ptb, lggr, tokenPoolStateAddresses)
	if err != nil {
		return err
	}

	// TODO: filter out messages that have a receiver that is not registered

	// TODO: move into its own file related to receives
	// Generate receiver call commands
	//nolint:gosec // G115:
	receiverCommands, err := GenerateReceiverCallCommands(lggr, messages, uint16(len(generatedTokenPoolCommands)))
	if err != nil {
		return err
	}

	// TODO: add make move vec command here (does this need to be done before or after the receiver calls?)
	hotPotatoVecResult := ptb.MakeMoveVec(strPtr("osh::...."), generatedTokenPoolCommands)

	// add the final PTB command (finish_execute) to the PTB using the interface from bindings
	// "&mut OffRampState",
	// "osh::ReceiverParams",
	// "vector<osh::CompletedDestTokenTransfer>",
	// TODO: how do we pass the vec of hot potatoes directly into this method since they are PTB dependencies?
	_, err = offrampEncoder.FinishExecuteWithArgs(bind.Object{Id: addressMappings.OffRampState}, initExecuteResult, hotPotatoVecResult)
	if err != nil {
		return fmt.Errorf("failed to encode move call (finish_execute) using bindings: %w", err)
	}

	return nil
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

func FilterRegisteredReceivers(
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
	ptb *transaction.Transaction,
	lggr logger.Logger,
	tokenPools []TokenPool,
) ([]transaction.Argument, error) {
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
