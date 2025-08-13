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

	// Generate N token pool commands and attach them to the PTB, each command must return a result
	// that will subsequently be used to make a vector of hot potatoes before finishing execution.
	tokenConfigs, err := tokenAdminRegistryDevInspect.GetTokenConfigs(ctx, callOpts, bind.Object{Id: addressMappings.CcipObjectRef}, coinMetadataAddresses)
	tokenPoolCommandsResults := make([]transaction.Argument, 0)
	for _, tokenPoolConfigs := range tokenConfigs {
		tokenPoolNormalizedModule, err := ptbClient.GetNormalizedModule(ctx, tokenPoolConfigs.TokenPoolPackageId, tokenPoolConfigs.TokenPoolModule)
		if err != nil {
			return fmt.Errorf("failed to get normalized module for token pool: %w", err)
		}

		tokenPoolCommandResult, err := AppendPTBCommandForTokenPool(ctx, lggr, sdkClient, ptb, callOpts, &addressMappings, &tokenPoolConfigs, &tokenPoolNormalizedModule)
		if err != nil {
			return fmt.Errorf("failed to append token pool command to PTB: %w", err)
		}

		tokenPoolCommandsResults = append(tokenPoolCommandsResults, *tokenPoolCommandResult)
	}

	// Create a receiver binding interface to filter out non-registered receivers
	receiverRegistryPkg, err := receiver_registry.NewReceiverRegistry(addressMappings.CcipPackageId, sdkClient)
	if err != nil {
		return err
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
			return fmt.Errorf("invalid receiver format, expected packageID:moduleID:functionName, got %s", message.Receiver)
		}

		receiverPackageId, receiverModule, receiverFunction := receiverParts[0], receiverParts[1], receiverParts[2]
		isRegistered, err := receiverRegistryDevInspect.IsRegisteredReceiver(ctx, callOpts, bind.Object{Id: addressMappings.CcipObjectRef}, receiverPackageId)
		if err != nil {
			return fmt.Errorf("failed to check if receiver is registered in offramp execution: %w", err)
		}
		if !isRegistered {
			// TODO: should this fail the whole execution?
			lggr.Debugw("receiver is not registered in offramp execution, skipping message...", "receiver", message.Receiver)
			continue
		}

		receiverCommandResult, err := AppendPTBCommandForReceiver(
			ctx,
			lggr,
			sdkClient,
			ptb,
			callOpts,
			initExecuteResult,
			receiverPackageId,
			receiverModule,
			receiverFunction,
			&addressMappings,
		)
		if err != nil {
			return err
		}
		receiverCommandsResults = append(receiverCommandsResults, *receiverCommandResult)
	}

	var ccipReceiveCommandResult *transaction.Argument
	if len(receiverCommandsResults) > 0 {
		// TODO: handle CCIP receive
	}

	// Make a vector of hot potatoes from all the token pool commands' results.
	// This will be passed into the final `finish_execute` call.
	// TODO: check if passing nil as a type is allowed for make_move_vec
	hotPotatoVecResult := ptb.MakeMoveVec(nil, tokenPoolCommandsResults)

	// TODO: check if the hot potato from the init or the ccip_receive is passed in here
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
	addressMappings *OffRampAddressMappings,
	tokenPoolConfigs *module_token_admin_registry.TokenConfig,
	normalizedModule *models.GetNormalizedMoveModuleResponse,
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

func AppendPTBCommandForReceiver(
	ctx context.Context,
	lggr logger.Logger,
	sdkClient sui.ISuiAPI,
	ptb *transaction.Transaction,
	callOpts *bind.CallOpts,
	initHotPotatoResult *transaction.Argument,
	packageId string,
	moduleId string,
	functionName string,
	addressMappings *OffRampAddressMappings,
) (*transaction.Argument, error) {
	boundReceiverContract, err := bind.NewBoundContract(packageId, packageId, moduleId, sdkClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create receiver bound contract when appending PTB command: %w", err)
	}

	typeArgsList := []string{}
	typeParamsList := []string{}
	encodedReceiverCall, err := boundReceiverContract.EncodeCallArgsWithGenerics(functionName, typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"address",
		//"_" // TODO: figure out the type for this
	}, []any{
		bind.Object{
			Id: addressMappings.CcipObjectRef,
		},
		//ownerCap,
		//maxFeeJuelsPerMsg,
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to encode receiver call: %w", err)
	}

	receiverCommandResult, err := boundReceiverContract.AppendPTB(ctx, callOpts, ptb, encodedReceiverCall)
	if err != nil {
		return nil, fmt.Errorf("failed to build PTB (receiver call) using bindings: %w", err)
	}

	return receiverCommandResult, nil
}
