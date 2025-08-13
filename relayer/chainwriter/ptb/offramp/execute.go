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
	module_token_admin_registry "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/token_admin_registry"
	module_offramp "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_offramp/offramp"
	"github.com/smartcontractkit/chainlink-sui/bindings/packages/ccip"
	"github.com/smartcontractkit/chainlink-sui/bindings/packages/offramp"
	"github.com/smartcontractkit/chainlink-sui/relayer/signer"

	"github.com/smartcontractkit/chainlink-sui/relayer/chainwriter/config"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
)

const OfframpTokenPoolFunctionName = "release_or_mint"

type SuiOffRampExecCallArgs struct {
	ReportContext [2][32]byte                `mapstructure:"ReportContext"`
	Report        []byte                     `mapstructure:"Report"`
	Info          ccipocr3.ExecuteReportInfo `mapstructure:"Info"`
}

// BuildOffRampExecutePTB builds the PTB for the OffRampExecute operation
func BuildOffRampExecutePTB(
	ctx context.Context,
	lggr logger.Logger,
	ptbClient client.SuiPTBClient,
	ptb *transaction.Transaction,
	args config.Arguments,
	signerAddress string,
	addressMappings OffRampAddressMappings,
) (err error) {
	sdkClient := ptbClient.GetClient()
	offrampArgs := &SuiOffRampExecCallArgs{}
	err = mapstructure.Decode(args.Args, &offrampArgs)
	if err != nil {
		return fmt.Errorf("failed to decode args for offramp execute PTB: %w", err)
	}

	coinMetadataAddresses := make([]string, 0)
	tokenAmounts := make([]ccipocr3.RampTokenAmount, 0)
	messages := make([]ccipocr3.Message, 0)

	// save all messages in a single slice
	for _, report := range offrampArgs.Info.AbstractReports {
		for _, message := range report.Messages {
			tokenAmounts = append(tokenAmounts, message.TokenAmounts...)
			messages = append(messages, message)
			for _, tokenAmount := range message.TokenAmounts {
				coinMetadataAddresses = append(coinMetadataAddresses, tokenAmount.DestTokenAddress.String())
			}
		}
	}

	devInspectSigner := signer.NewDevInspectSigner(signerAddress)

	// Call options for bindings DevInspect calls
	callOpts := &bind.CallOpts{
		Signer:           devInspectSigner,
		WaitForExecution: true,
	}

	// Set the ccip package interface from bindings
	ccipPkg, err := ccip.NewCCIP(addressMappings.CcipPackageId, sdkClient)
	if err != nil {
		return err
	}
	tokenAdminRegistryContract := ccipPkg.TokenAdminRegistry().(*module_token_admin_registry.TokenAdminRegistryContract)
	tokenAdminRegistryDevInspect := tokenAdminRegistryContract.DevInspect()

	// TODO: remove this, it's not needed since we make a `GetTokenConfigs` call which includes the pools
	_, err = tokenAdminRegistryDevInspect.GetPools(ctx, callOpts, bind.Object{Id: addressMappings.CcipObjectRef}, coinMetadataAddresses)
	if err != nil {
		return fmt.Errorf("failed to get pools from token admin registry: %w", err)
	}

	// Set the offramp package interface from bindings
	offrampPkg, err := offramp.NewOfframp(addressMappings.OffRampPackageId, sdkClient)
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

	initExecuteResult, err := offrampContract.AppendPTB(ctx, callOpts, ptb, encodedInitExecute)
	if err != nil {
		return fmt.Errorf("failed to build PTB (init_execute) using bindings: %w", err)
	}

	// Generate N token pool commands and attach them to the PTB, each command must return a result
	// that will subsequently be used to make a vector of hot potatoes before finishing execution.
	tokenConfigs, err := tokenAdminRegistryDevInspect.GetTokenConfigs(ctx, callOpts, bind.Object{Id: addressMappings.CcipObjectRef}, coinMetadataAddresses)
	tokenPoolCommandsResults := make([]transaction.Argument, 0)
	for _, tokenPoolConfigs := range tokenConfigs {
		tokenPoolCommandResult, err := AppendPTBCommandForTokenPool(ctx, lggr, sdkClient, ptb, callOpts, tokenPoolConfigs)
		if err != nil {
			return fmt.Errorf("failed to append token pool command to PTB: %w", err)
		}
		tokenPoolCommandsResults = append(tokenPoolCommandsResults, *tokenPoolCommandResult)
	}

	// TODO: filter out messages that have a receiver that is not registered
	// Generate receiver call commands
	//nolint:gosec // G115:
	receiverCommands, err := GenerateReceiverCallCommands(lggr, messages, uint16(len(tokenPoolCommandsResults)))
	if err != nil {
		return err
	}

	// Make a vector of hot potatoes from all the token pool commands' results.
	// This will be passed into the final `finish_execute` call.
	// TODO: check if passing nil as a type is allowed for make_move_vec
	hotPotatoVecResult := ptb.MakeMoveVec(nil, tokenPoolCommandsResults)

	// add the final PTB command (finish_execute) to the PTB using the interface from bindings
	encodedFinishExecute, err := offrampEncoder.FinishExecuteWithArgs(bind.Object{Id: addressMappings.OffRampState}, initExecuteResult, hotPotatoVecResult)
	if err != nil {
		return fmt.Errorf("failed to encode move call (finish_execute) using bindings: %w", err)
	}

	_, err = offrampContract.AppendPTB(ctx, callOpts, ptb, encodedFinishExecute)
	if err != nil {
		return fmt.Errorf("failed to build PTB (finish_execute) using bindings: %w", err)
	}

	return nil
}

func AppendPTBCommandForTokenPool(
	ctx context.Context,
	lggr logger.Logger,
	sdkClient sui.ISuiAPI,
	ptb *transaction.Transaction,
	callOpts *bind.CallOpts,
	tokenPoolConfigs module_token_admin_registry.TokenConfig,
) (*transaction.Argument, error) {
	poolBoundContract, err := bind.NewBoundContract(
		tokenPoolConfigs.TokenPoolPackageId,
		tokenPoolConfigs.TokenPoolPackageId,
		tokenPoolConfigs.TokenPoolModule,
		sdkClient,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create token pool bound contract when appending PTB command: %w", err)
	}

	typeArgsList := []string{}
	typeParamsList := []string{}
	encodedTokenPoolCall, err := poolBoundContract.EncodeCallArgsWithGenerics(OfframpTokenPoolFunctionName, typeArgsList, typeParamsList, []string{
		//"&mut CCIPObjectRef",
		//"&OwnerCap",
		//"u256",
		//"address",
		//"u64",
		//"vector<address>",
	}, []any{
		//ref,
		//ownerCap,
		//maxFeeJuelsPerMsg,
		//linkToken,
		//tokenPriceStalenessThreshold,
		//feeTokens,
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to encode token pool call: %w", err)
	}

	tokenPoolCommandResult, err := poolBoundContract.AppendPTB(ctx, callOpts, ptb, encodedTokenPoolCall)
	if err != nil {
		return nil, fmt.Errorf("failed to build PTB (token pool call) using bindings: %w", err)
	}

	return tokenPoolCommandResult, nil
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
