// / A package to build all the bespoke code (PTB) along with its commands for the OffRampExecute operation.
// / There will be no dependency on the PTBConstructor interface here due to writing entirely custom code that is not meant to be re-usable.
// / This package does not generate CW configs but rather generates the actual PTB along with its commands directly.
package offramp

import (
	"context"
	"fmt"
	"strings"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/block-vision/sui-go-sdk/transaction"
	"github.com/mitchellh/mapstructure"
	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	receiver_registry "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/receiver_registry"
	module_token_admin_registry "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/token_admin_registry"
	module_offramp "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_offramp/offramp"
	"github.com/smartcontractkit/chainlink-sui/bindings/packages/ccip"
	"github.com/smartcontractkit/chainlink-sui/bindings/packages/offramp"
	"github.com/smartcontractkit/chainlink-sui/relayer/signer"

	"github.com/smartcontractkit/chainlink-sui/relayer/chainwriter/config"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
)

const OfframpTokenPoolFunctionName = "release_or_mint"

type SuiOffRampExecCallArgs struct {
	ReportContext [2][32]byte                `mapstructure:"ReportContext"`
	Report        []byte                     `mapstructure:"Report"`
	Info          ccipocr3.ExecuteReportInfo `mapstructure:"Info"`
}

// TODO: remove once hot potato approach validated
//type ReceiverParams struct {
//	RemoteChainSelector uint64
//	Receiver            [32]byte
//	SourceAmount        uint64
//	DestTokenAddress    [32]byte
//	SourcePoolAddress   []byte
//	SourcePoolData      []byte
//	OffchainTokenData   []byte
//}

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
	messages := make([]ccipocr3.Message, 0)

	// TODO: remove once hot potato approach validated
	//receiverParamsData := map[string]ReceiverParams{}

	// Prepare some data from the execution report for easy access when constructing the PTB commands.
	for _, report := range offrampArgs.Info.AbstractReports {
		for _, message := range report.Messages {
			messages = append(messages, message)
			for _, tokenAmount := range message.TokenAmounts {
				destTokenAddress := tokenAmount.DestTokenAddress.String()
				coinMetadataAddresses = append(coinMetadataAddresses, destTokenAddress)

				// TODO: remove once hot potato approach validated
				//receiverParamsData[destTokenAddress] = ReceiverParams{
				//	RemoteChainSelector: uint64(report.SourceChainSelector),
				//	Receiver:            [32]byte(message.Receiver),
				//	SourceAmount:        tokenAmount.Amount.Uint64(),
				//	DestTokenAddress:    [32]byte(tokenAmount.DestTokenAddress),
				//	SourcePoolAddress:   tokenAmount.SourcePoolAddress,
				//	// TODO: double check if the following fields are correct
				//	SourcePoolData:    tokenAmount.ExtraData,
				//	OffchainTokenData: tokenAmount.DestExecData,
				//}
			}
		}
	}

	// An interface used to make dev inspect calls in bindings, actual signing does not happen here.
	devInspectSigner := signer.NewDevInspectSigner(signerAddress)

	// Call options for bindings DevInspect calls
	callOpts := &bind.CallOpts{
		Signer:           devInspectSigner,
		WaitForExecution: true,
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
		bind.Object{Id: addressMappings.CcipObjectRef},
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

	// Process each token pool from this offramp execution after getting their configs
	// from the registry. Attach the commands to the PTB and return their argument results.
	tokenPoolCommandsResults, err := ProcessTokenPools(
		ctx,
		lggr,
		ptbClient,
		ptb,
		&addressMappings,
		callOpts,
		coinMetadataAddresses,
		initExecuteResult,
	)

	// Process each message and create PTB commands for each (valid) receiver.
	_, err = ProcessReceivers(
		ctx,
		lggr,
		ptbClient,
		ptb,
		messages,
		&addressMappings,
		callOpts,
		initExecuteResult,
	)
	if err != nil {
		return err
	}

	// Make a vector of hot potatoes from all the token pool commands' results.
	// This will be passed into the final `finish_execute` call.
	hotPotatoVecResult := ptb.MakeMoveVec(AnyPointer("_"), tokenPoolCommandsResults)

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

func ProcessTokenPools(
	ctx context.Context,
	lggr logger.Logger,
	ptbClient client.SuiPTBClient,
	ptb *transaction.Transaction,
	addressMappings *OffRampAddressMappings,
	callOpts *bind.CallOpts,
	coinMetadataAddresses []string,
	receiverParams *transaction.Argument,
) ([]transaction.Argument, error) {
	sdkClient := ptbClient.GetClient()

	// Set the ccip package interface from bindings
	ccipPkg, err := ccip.NewCCIP(addressMappings.CcipPackageId, sdkClient)
	if err != nil {
		return nil, err
	}
	tokenAdminRegistryContract := ccipPkg.TokenAdminRegistry().(*module_token_admin_registry.TokenAdminRegistryContract)
	tokenAdminRegistryDevInspect := tokenAdminRegistryContract.DevInspect()

	// Generate N token pool commands and attach them to the PTB, each command must return a result
	// that will subsequently be used to make a vector of hot potatoes before finishing execution.
	tokenConfigs, err := tokenAdminRegistryDevInspect.GetTokenConfigs(ctx, callOpts, bind.Object{Id: addressMappings.CcipObjectRef}, coinMetadataAddresses)
	if err != nil {
		return nil, fmt.Errorf("failed to get token configs for offramp execution: %w", err)
	}

	tokenPoolCommandsResults := make([]transaction.Argument, 0)
	for idx, tokenPoolConfigs := range tokenConfigs {
		// TODO: remove once hot potato approach validated
		//// Get the relevant receiver params data for this token pool
		//tokenPoolEncodedData := receiverParamsData[coinMetadataAddresses[idx]]

		// Get the move normalized module to dynamically construct the parameters for the token pool call
		tokenPoolNormalizedModule, err := ptbClient.GetNormalizedModule(ctx, tokenPoolConfigs.TokenPoolPackageId, tokenPoolConfigs.TokenPoolModule)
		if err != nil {
			return nil, fmt.Errorf("failed to get normalized module for token pool: %w", err)
		}

		tokenPoolCommandResult, err := AppendPTBCommandForTokenPool(
			ctx,
			lggr,
			sdkClient,
			ptb,
			callOpts,
			addressMappings,
			&tokenPoolConfigs,
			&tokenPoolNormalizedModule,
			receiverParams,
			idx,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to append token pool command to PTB: %w", err)
		}

		tokenPoolCommandsResults = append(tokenPoolCommandsResults, *tokenPoolCommandResult)
	}

	return tokenPoolCommandsResults, nil
}

func AppendPTBCommandForTokenPool(
	ctx context.Context,
	lggr logger.Logger,
	sdkClient sui.ISuiAPI,
	ptb *transaction.Transaction,
	callOpts *bind.CallOpts,
	addressMappings *OffRampAddressMappings,
	tokenPoolConfigs *module_token_admin_registry.TokenConfig,
	normalizedModule *models.GetNormalizedMoveModuleResponse,
	receiverParams *transaction.Argument,
	index int,
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

	// TODO: replace this call with generated bindings call once ready
	offrampStateHelperContract, err := bind.NewBoundContract(
		addressMappings.CcipPackageId,
		addressMappings.CcipPackageId,
		"offramp_state_helper",
		sdkClient,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create offramp state helper bound contract when appending PTB command: %w", err)
	}

	typeArgsList := []string{}
	typeParamsList := []string{}
	paramTypes := []string{
		"&mut ReceiverParams",
		"u64",
	}
	paramValues := []any{
		receiverParams,
		index,
	}

	encodedGetTokenParamDataCall, err := offrampStateHelperContract.EncodeCallArgsWithGenerics(
		"get_dest_token_transfer",
		typeArgsList,
		typeParamsList,
		paramTypes,
		paramValues,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to encode get_token_param_data call: %w", err)
	}

	getTokenParamDataCommandResult, err := offrampStateHelperContract.AppendPTB(ctx, callOpts, ptb, encodedGetTokenParamDataCall)
	if err != nil {
		return nil, fmt.Errorf("failed to build PTB (get_token_param_data) using bindings: %w", err)
	}

	// TODO: remove once hot potato approach validated
	//// Encode the data needed by the token pool call
	//bcsEncoder := &aptosBCS.Serializer{}
	//bcsEncoder.U64(data.RemoteChainSelector)
	//bcsEncoder.FixedBytes(data.Receiver[:])
	//bcsEncoder.U64(data.SourceAmount)
	//bcsEncoder.FixedBytes(data.DestTokenAddress[:])
	//bcsEncoder.WriteBytes(data.SourcePoolAddress)
	//bcsEncoder.WriteBytes(data.SourcePoolData)
	//bcsEncoder.WriteBytes(data.OffchainTokenData)
	//encodedData := bcsEncoder.ToBytes()

	typeArgsList = []string{}
	typeParamsList = []string{}
	paramTypes = []string{}
	// The fixed arguments that must be present for every token pool call.
	paramValues = []any{
		bind.Object{Id: addressMappings.CcipObjectRef},
		getTokenParamDataCommandResult,
	}

	// Append dynamic values (addresses) to the paramValues for the token pool call.
	// This allows an unknown set of addresses to be passed in.
	for _, value := range tokenPoolConfigs.ReleaseOrMintParams {
		paramValues = append(paramValues, value)
	}

	// Use the normalized module to populate the paramTypes and paramValues for the bound contract
	functionSignature, ok := normalizedModule.ExposedFunctions[OfframpTokenPoolFunctionName]
	if !ok {
		return nil, fmt.Errorf("missing function signature for token pool function not found in module (%s)", OfframpTokenPoolFunctionName)
	}

	// Figure out the parameter types from the normalized module of the token pool
	paramTypes, err = DecodeParameters(lggr, functionSignature.(map[string]any), "parameters")
	if err != nil {
		return nil, fmt.Errorf("failed to decode parameters for token pool function: %w", err)
	}

	encodedTokenPoolCall, err := poolBoundContract.EncodeCallArgsWithGenerics(
		OfframpTokenPoolFunctionName,
		typeArgsList,
		typeParamsList,
		paramTypes,
		paramValues,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to encode token pool call: %w", err)
	}

	tokenPoolCommandResult, err := poolBoundContract.AppendPTB(ctx, callOpts, ptb, encodedTokenPoolCall)
	if err != nil {
		return nil, fmt.Errorf("failed to build PTB (token pool call) using bindings: %w", err)
	}

	return tokenPoolCommandResult, nil
}

func ProcessReceivers(
	ctx context.Context,
	lggr logger.Logger,
	ptbClient client.SuiPTBClient,
	ptb *transaction.Transaction,
	messages []ccipocr3.Message,
	addressMappings *OffRampAddressMappings,
	callOpts *bind.CallOpts,
	receiverParams *transaction.Argument,
) ([]transaction.Argument, error) {
	sdkClient := ptbClient.GetClient()

	// Create a receiver binding interface to filter out non-registered receivers
	receiverRegistryPkg, err := receiver_registry.NewReceiverRegistry(addressMappings.CcipPackageId, sdkClient)
	if err != nil {
		return nil, err
	}
	receiverRegistryDevInspect := receiverRegistryPkg.DevInspect()

	receiverCommandsResults := make([]transaction.Argument, 0)
	// Generate receiver call commands
	for _, message := range messages {
		// If there is no receiver, skip this message
		if len(message.Receiver) == 0 || message.Receiver == nil {
			lggr.Debugw("no receiver specified, skipping message in offramp execution...", "message", message)
			continue
		}
		// If there is no data, skip this message
		if len(message.Data) == 0 {
			lggr.Debugw("no data specified, skipping message in offramp execution...", "message", message)
			continue
		}

		// Parse the receiver string into `packageID::moduleID::functionName` format
		receiverParts := strings.Split(string(message.Receiver), "::")
		if len(receiverParts) != 3 {
			return nil, fmt.Errorf("invalid receiver format, expected packageID:moduleID:functionName, got %s", message.Receiver)
		}

		receiverPackageId, receiverModule, receiverFunction := receiverParts[0], receiverParts[1], receiverParts[2]
		isRegistered, err := receiverRegistryDevInspect.IsRegisteredReceiver(ctx, callOpts, bind.Object{Id: addressMappings.CcipObjectRef}, receiverPackageId)
		if err != nil {
			return nil, fmt.Errorf("failed to check if receiver is registered in offramp execution: %w", err)
		}
		// If the receiver is not registered, fail the entire execution
		if !isRegistered {
			return nil, fmt.Errorf("receiver is not registered in offramp execution. error: %s", message.Receiver)
		}

		// Get the receiver config via the receiver registry binding
		receiverConfig, err := receiverRegistryDevInspect.GetReceiverConfig(
			ctx,
			callOpts,
			bind.Object{Id: addressMappings.CcipObjectRef},
			receiverPackageId,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to get receiver config for offramp execution: %w", err)
		}

		receiverNormalizedModule, err := ptbClient.GetNormalizedModule(ctx, receiverPackageId, receiverModule)
		if err != nil {
			return nil, fmt.Errorf("failed to get normalized module for token pool: %w", err)
		}

		receiverCommandResult, err := AppendPTBCommandForReceiver(
			ctx,
			lggr,
			sdkClient,
			ptb,
			callOpts,
			receiverPackageId,
			receiverModule,
			receiverFunction,
			addressMappings,
			&receiverConfig,
			&receiverNormalizedModule,
			receiverParams,
		)
		if err != nil {
			return nil, err
		}
		receiverCommandsResults = append(receiverCommandsResults, *receiverCommandResult)
	}

	return receiverCommandsResults, nil
}

func AppendPTBCommandForReceiver(
	ctx context.Context,
	lggr logger.Logger,
	sdkClient sui.ISuiAPI,
	ptb *transaction.Transaction,
	callOpts *bind.CallOpts,
	packageId string,
	moduleId string,
	functionName string,
	addressMappings *OffRampAddressMappings,
	receiverConfig *receiver_registry.ReceiverConfig,
	normalizedModule *models.GetNormalizedMoveModuleResponse,
	receiverParams *transaction.Argument,
) (*transaction.Argument, error) {
	boundReceiverContract, err := bind.NewBoundContract(packageId, packageId, moduleId, sdkClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create receiver bound contract when appending PTB command: %w", err)
	}

	offrampStateHelperContract, err := bind.NewBoundContract(
		addressMappings.CcipPackageId,
		addressMappings.CcipPackageId,
		"offramp_state_helper",
		sdkClient,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create offramp state helper bound contract when appending PTB command: %w", err)
	}

	typeArgsList := []string{}
	typeParamsList := []string{}
	paramTypes := []string{}
	paramValues := []any{
		bind.Object{Id: addressMappings.CcipObjectRef},
		receiverParams,
		// TODO: figure out what else is needed
	}

	encodedAny2SuiExtractCall, err := offrampStateHelperContract.EncodeCallArgsWithGenerics(
		"extract_any2sui_message",
		typeArgsList,
		typeParamsList,
		paramTypes,
		paramValues,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to encode get_token_param_data call: %w", err)
	}

	extractedAny2SuiMessageResult, err := offrampStateHelperContract.AppendPTB(ctx, callOpts, ptb, encodedAny2SuiExtractCall)
	if err != nil {
		return nil, fmt.Errorf("failed to build PTB (get_token_param_data) using bindings: %w", err)
	}

	typeArgsList = []string{}
	typeParamsList = []string{}
	paramTypes = []string{}
	paramValues = []any{
		bind.Object{Id: addressMappings.CcipObjectRef},
		extractedAny2SuiMessageResult,
	}

	// Use the normalized module to populate the paramTypes and paramValues for the bound contract
	functionSignature, ok := normalizedModule.ExposedFunctions[OfframpTokenPoolFunctionName]
	if !ok {
		return nil, fmt.Errorf("missing function signature for token pool function not found in module (%s)", OfframpTokenPoolFunctionName)
	}

	// Figure out the parameter types from the normalized module of the token pool
	paramTypes, err = DecodeParameters(lggr, functionSignature.(map[string]any), "parameters")
	if err != nil {
		return nil, fmt.Errorf("failed to decode parameters for token pool function: %w", err)
	}

	// Append dynamic values (addresses) to the paramValues for the receiver call.
	// This is used for state references for the receiver (similar to the token pool call).
	for _, value := range receiverConfig.ReceiverStateParams {
		paramValues = append(paramValues, value)
	}

	encodedReceiverCall, err := boundReceiverContract.EncodeCallArgsWithGenerics(
		functionName,
		typeArgsList,
		typeParamsList,
		paramTypes,
		paramValues,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to encode receiver call: %w", err)
	}

	receiverCommandResult, err := boundReceiverContract.AppendPTB(ctx, callOpts, ptb, encodedReceiverCall)
	if err != nil {
		return nil, fmt.Errorf("failed to build PTB (receiver call) using bindings: %w", err)
	}

	return receiverCommandResult, nil
}
