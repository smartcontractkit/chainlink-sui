// / A package to build all the bespoke code (PTB) along with its commands for the OffRampExecute operation.
// / There will be no dependency on the PTBConstructor interface here due to writing entirely custom code that is not meant to be re-usable.
// / This package does not generate CW configs but rather generates the actual PTB along with its commands directly.
package offramp

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
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
	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
	"github.com/smartcontractkit/chainlink-sui/relayer/signer"

	"github.com/smartcontractkit/chainlink-sui/relayer/chainwriter/config"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
)

const OfframpTokenPoolFunctionName = "release_or_mint"

type SuiOffRampExecCallArgs struct {
	ReportContext [2][32]byte                `mapstructure:"ReportContext"`
	Report        []byte                     `mapstructure:"Report"`
	Info          ccipocr3.ExecuteReportInfo `mapstructure:"Info"`
	ExtraData     ExtraDataDecoded           `mapstructure:"ExtraData"`
}

type ExtraDataDecoded struct {
	// ExtraArgsDecoded contain message specific extra args.
	ExtraArgsDecoded map[string]any
	// DestExecDataDecoded contain token transfer specific extra args.
	DestExecDataDecoded []map[string]any
}

func DecodeOffRampExecCallArgs(args map[string]any) (*SuiOffRampExecCallArgs, error) {
	offrampArgs := &SuiOffRampExecCallArgs{}
	err := mapstructure.Decode(args, &offrampArgs)
	if err != nil {
		return nil, fmt.Errorf("failed to decode args for offramp execute PTB: %w", err)
	}

	return offrampArgs, nil
}

// BuildOffRampExecutePTB builds the PTB for the OffRampExecute operation (V1).
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
	offrampArgs, err := DecodeOffRampExecCallArgs(args.Args)
	if err != nil {
		return fmt.Errorf("failed to decode args for offramp execute PTB: %w", err)
	}

	coinMetadataAddresses := make([]string, 0)
	messages := make([]ccipocr3.Message, 0)

	// Prepare some data from the execution report for easy access when constructing the PTB commands.
	for _, report := range offrampArgs.Info.AbstractReports {
		for _, message := range report.Messages {
			messages = append(messages, message)
			for _, tokenAmount := range message.TokenAmounts {
				destTokenAddress := "0x" + hex.EncodeToString(tokenAmount.DestTokenAddress)

				lggr.Debugw("found token metadata address", "address", destTokenAddress)

				coinMetadataAddresses = append(coinMetadataAddresses, destTokenAddress)
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

	// Get the latest package ID from the offramp module
	latestOfframpPackageId, err := ptbClient.GetLatestPackageId(ctx, addressMappings.OffRampPackageId, "offramp")
	if err != nil {
		return err
	}

	latestCcipPackageId, err := ptbClient.GetLatestPackageId(ctx, addressMappings.CcipPackageId, "state_object")
	if err != nil {
		return err
	}

	// Update both the offramp and ccip package IDs to the latest versions
	addressMappings.OffRampPackageId = latestOfframpPackageId
	addressMappings.CcipPackageId = latestCcipPackageId

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
	if err != nil {
		return err
	}

	lggr.Debugw("finished processing token pool calls", "tokenPoolCalls", tokenPoolCommandsResults)

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
		offrampArgs.ExtraData.ExtraArgsDecoded,
	)
	if err != nil {
		return err
	}

	// add the final PTB command (finish_execute) to the PTB using the interface from bindings
	encodedFinishExecute, err := offrampEncoder.FinishExecuteWithArgs(bind.Object{Id: addressMappings.CcipObjectRef}, bind.Object{Id: addressMappings.OffRampState}, initExecuteResult)
	if err != nil {
		return fmt.Errorf("failed to encode move call (finish_execute) using bindings: %w", err)
	}

	_, err = offrampContract.AppendPTB(ctx, callOpts, ptb, encodedFinishExecute)
	if err != nil {
		return fmt.Errorf("failed to build PTB (finish_execute) using bindings: %w", err)
	}

	return nil
}

// BuildOffRampExecutePTBV2 builds the PTB for the OffRampExecute operation using V2 functions
// that enforce receiver object binding at the protocol level. The V2 leaf hash includes
// receiver_object_ids, and extract_any2sui_message_v2 verifies them before message delivery.
func BuildOffRampExecutePTBV2(
	ctx context.Context,
	lggr logger.Logger,
	ptbClient client.SuiPTBClient,
	ptb *transaction.Transaction,
	args config.Arguments,
	signerAddress string,
	addressMappings OffRampAddressMappings,
) (err error) {
	sdkClient := ptbClient.GetClient()
	offrampArgs, err := DecodeOffRampExecCallArgs(args.Args)
	if err != nil {
		return fmt.Errorf("failed to decode args for offramp execute PTB v2: %w", err)
	}

	coinMetadataAddresses := make([]string, 0)
	messages := make([]ccipocr3.Message, 0)

	for _, report := range offrampArgs.Info.AbstractReports {
		for _, message := range report.Messages {
			messages = append(messages, message)
			for _, tokenAmount := range message.TokenAmounts {
				destTokenAddress := "0x" + hex.EncodeToString(tokenAmount.DestTokenAddress)
				lggr.Debugw("found token metadata address", "address", destTokenAddress)
				coinMetadataAddresses = append(coinMetadataAddresses, destTokenAddress)
			}
		}
	}

	devInspectSigner := signer.NewDevInspectSigner(signerAddress)
	callOpts := &bind.CallOpts{
		Signer:           devInspectSigner,
		WaitForExecution: true,
	}

	latestOfframpPackageId, err := ptbClient.GetLatestPackageId(ctx, addressMappings.OffRampPackageId, "offramp")
	if err != nil {
		return err
	}

	latestCcipPackageId, err := ptbClient.GetLatestPackageId(ctx, addressMappings.CcipPackageId, "state_object")
	if err != nil {
		return err
	}

	addressMappings.OffRampPackageId = latestOfframpPackageId
	addressMappings.CcipPackageId = latestCcipPackageId

	offrampPkg, err := offramp.NewOfframp(addressMappings.OffRampPackageId, sdkClient)
	if err != nil {
		return err
	}
	offrampContract := offrampPkg.Offramp().(*module_offramp.OfframpContract)
	offrampEncoder := offrampContract.Encoder()

	encodedInitExecute, err := offrampEncoder.InitExecuteV2(
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
		return fmt.Errorf("failed to encode move call (init_execute_v2) using bindings: %w", err)
	}

	initExecuteResult, err := offrampContract.AppendPTB(ctx, callOpts, ptb, encodedInitExecute)
	if err != nil {
		return fmt.Errorf("failed to build PTB (init_execute_v2) using bindings: %w", err)
	}

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
	if err != nil {
		return err
	}

	lggr.Debugw("finished processing token pool calls", "tokenPoolCalls", tokenPoolCommandsResults)

	_, err = ProcessReceiversV2(
		ctx,
		lggr,
		ptbClient,
		ptb,
		messages,
		&addressMappings,
		callOpts,
		initExecuteResult,
		offrampArgs.ExtraData.ExtraArgsDecoded,
	)
	if err != nil {
		return err
	}

	encodedFinishExecute, err := offrampEncoder.FinishExecuteV2WithArgs(bind.Object{Id: addressMappings.CcipObjectRef}, bind.Object{Id: addressMappings.OffRampState}, initExecuteResult)
	if err != nil {
		return fmt.Errorf("failed to encode move call (finish_execute_v2) using bindings: %w", err)
	}

	_, err = offrampContract.AppendPTB(ctx, callOpts, ptb, encodedFinishExecute)
	if err != nil {
		return fmt.Errorf("failed to build PTB (finish_execute_v2) using bindings: %w", err)
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

	lggr.Debugw("processing token pools for offramp execution...", "coinMetadataAddresses", coinMetadataAddresses)

	// Set the ccip package interface from bindings
	ccipPkg, err := ccip.NewCCIP(addressMappings.CcipPackageId, sdkClient)
	if err != nil {
		return nil, err
	}
	tokenAdminRegistryContract := ccipPkg.TokenAdminRegistry().(*module_token_admin_registry.TokenAdminRegistryContract)
	tokenAdminRegistryDevInspect := tokenAdminRegistryContract.DevInspect()

	// Generate N token pool commands and attach them to the PTB, each command must return a result
	// that will subsequently be used to make a vector of hot potatoes before finishing execution.
	tokenPoolCommandsResults := make([]transaction.Argument, 0)

	// NOTE: there will only ever be one token pool per offramp execution, but we loop over the addresses
	// as they are provided in the form of an array from the core node. We can alternatively simply read
	// the first index but we do this to allow for simplified future updates.
	for _, coinMetadataAddress := range coinMetadataAddresses {
		tokenConfig, err := tokenAdminRegistryDevInspect.GetTokenConfigStruct(ctx, callOpts, bind.Object{Id: addressMappings.CcipObjectRef}, coinMetadataAddress)
		if err != nil {
			return nil, fmt.Errorf("failed to get token configs for offramp execution: %w", err)
		} else if !IsValidTokenPoolConfig(&tokenConfig) {
			return nil, fmt.Errorf("invalid token pool config for token metadata address: %s: %+v", coinMetadataAddress, tokenConfig)
		}

		lggr.Debugw("fetched token configs via dev inspect call", "tokenConfig", tokenConfig)

		// Get the move normalized module to dynamically construct the parameters for the token pool call
		tokenPoolNormalizedModule, err := ptbClient.GetNormalizedModule(ctx, tokenConfig.TokenPoolPackageId, tokenConfig.TokenPoolModule)
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
			&tokenConfig,
			&tokenPoolNormalizedModule,
			receiverParams,
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

	tokenType := tokenPoolConfigs.TokenType
	if !strings.HasPrefix(tokenType, "0x") {
		tokenType = "0x" + tokenType
	}

	typeArgsList := []string{tokenType}
	typeParamsList := []string{}
	// The fixed arguments that must be present for every token pool call.
	paramValues := []any{
		bind.Object{Id: addressMappings.CcipObjectRef},
		receiverParams,
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
	paramTypes, err := DecodeParameters(lggr, functionSignature.(map[string]any), "parameters")
	if err != nil {
		return nil, fmt.Errorf("failed to decode parameters for token pool function: %w", err)
	}

	lggr.Debugw("calling token pool", "paramTypes", paramTypes, "paramValues", paramValues)

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
	extraArgs map[string]any,
) ([]transaction.Argument, error) {
	sdkClient := ptbClient.GetClient()

	receiverRegistryPkg, err := receiver_registry.NewReceiverRegistry(addressMappings.CcipPackageId, sdkClient)
	if err != nil {
		return nil, err
	}
	receiverRegistryDevInspect := receiverRegistryPkg.DevInspect()

	receiverCommandsResults := make([]transaction.Argument, 0)
	for _, message := range messages {
		if len(message.Receiver) == 0 || message.Receiver == nil {
			lggr.Errorw("unexpected nil or zero length receiver, skipping message in offramp execution...", "message", message)
			continue
		}

		if bytes.Equal(message.Receiver, codec.AccountZero) {
			lggr.Debugw("receiver is zero address, skipping message in offramp execution...", "message", message)
			continue
		}

		// Mirror on-chain gating: skip receiver call when on-chain would not populate the message.
		// On-chain: has_valid_message_receiver = (!data.is_empty() || gas_limit != 0) && is_registered_receiver
		if !needsAppDelivery(message, extraArgs) {
			lggr.Debugw("message has no data and zero gas limit, skipping receiver call",
				"receiver", hex.EncodeToString(message.Receiver))
			continue
		}

		receiverPackageId := "0x" + hex.EncodeToString(message.Receiver)

		isRegistered, err := receiverRegistryDevInspect.IsRegisteredReceiver(ctx, callOpts, bind.Object{Id: addressMappings.CcipObjectRef}, receiverPackageId)
		if err != nil {
			return nil, fmt.Errorf("failed to check if receiver is registered in offramp execution: %w", err)
		}
		if !isRegistered {
			lggr.Warnw("receiver not registered, skipping receiver call (on-chain will not populate message)",
				"receiver", receiverPackageId)
			continue
		}

		receiverConfig, err := receiverRegistryDevInspect.GetReceiverConfig(ctx, callOpts, bind.Object{Id: addressMappings.CcipObjectRef}, receiverPackageId)
		if err != nil {
			// RPC/network error — propagate so the caller can retry later.
			return nil, fmt.Errorf("failed to get receiver config in offramp execution: %w", err)
		}

		receiverNormalizedModule, err := ptbClient.GetNormalizedModule(ctx, receiverPackageId, receiverConfig.ModuleName)
		if err != nil {
			// RPC/network error — propagate so the caller can retry later.
			return nil, fmt.Errorf("failed to get normalized module for receiver: %w", err)
		}

		receiverCommandResult, err := AppendPTBCommandForReceiver(
			ctx,
			lggr,
			sdkClient,
			ptb,
			callOpts,
			receiverPackageId,
			receiverConfig.ModuleName,
			"ccip_receive",
			addressMappings,
			message.Header.MessageID,
			&receiverNormalizedModule,
			receiverParams,
			extraArgs,
		)
		skip, retErr := classifyReceiverBuildError(receiverPackageId, err)
		if skip {
			lggr.Errorw("skipping receiver command due to unsupported ABI; PTB will fail on-chain",
				"receiver", receiverPackageId,
				"error", err)
			continue
		}
		if retErr != nil {
			return nil, retErr
		}
		receiverCommandsResults = append(receiverCommandsResults, *receiverCommandResult)
	}

	return receiverCommandsResults, nil
}

func needsAppDelivery(message ccipocr3.Message, extraArgs map[string]any) bool {
	if len(message.Data) > 0 {
		return true
	}
	if val, ok := extraArgs["gasLimit"]; ok {
		switch gl := val.(type) {
		case *big.Int:
			return gl != nil && gl.Sign() > 0
		case uint64:
			return gl > 0
		}
	}
	return false
}

// extractReceiverObjectIdStrings parses receiver_object_ids from the extraArgs map
// into hex-prefixed address strings. Returns an empty slice if the key is missing or
// contains no valid entries.
func extractReceiverObjectIdStrings(extraArgs map[string]any) []string {
	raw, ok := extraArgs["receiverObjectIds"]
	if !ok {
		return []string{}
	}

	var result []string
	switch vals := raw.(type) {
	case [][]byte:
		for _, v := range vals {
			result = append(result, "0x"+hex.EncodeToString(v))
		}
	case []any:
		for _, v := range vals {
			if b, ok := v.([]byte); ok {
				result = append(result, "0x"+hex.EncodeToString(b))
			}
		}
	}
	if result == nil {
		return []string{}
	}
	return result
}

// classifyReceiverBuildError distinguishes permanent unsupported-ABI failures from
// transient build errors. Unsupported ABI errors are skipped so the PTB can be
// submitted and fail on-chain with ECCIPReceiveFailed.
func classifyReceiverBuildError(receiverPackageId string, err error) (skip bool, retErr error) {
	if err == nil {
		return false, nil
	}
	if errors.Is(err, ErrUnsupportedReceiverABI) {
		return true, nil
	}
	return false, fmt.Errorf("failed to build receiver command for %s: %w", receiverPackageId, err)
}

func appendCcipReceiveCommand(
	ctx context.Context,
	lggr logger.Logger,
	ptb *transaction.Transaction,
	callOpts *bind.CallOpts,
	boundReceiverContract *bind.BoundContract,
	functionName string,
	addressMappings *OffRampAddressMappings,
	messageID [32]byte,
	normalizedModule *models.GetNormalizedMoveModuleResponse,
	extractedAny2SuiMessageResult *transaction.Argument,
	receiverObjectIdStrings []string,
) (*transaction.Argument, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	paramValues := []any{
		messageID,
		bind.Object{Id: addressMappings.CcipObjectRef},
		extractedAny2SuiMessageResult,
	}

	functionSignature, ok := normalizedModule.ExposedFunctions[functionName]
	if !ok {
		return nil, fmt.Errorf("%w: function %q not found in module", ErrUnsupportedReceiverABI, functionName)
	}

	paramTypes, err := DecodeParameters(lggr, functionSignature.(map[string]any), "parameters")
	if err != nil {
		return nil, fmt.Errorf("failed to decode receiver parameters: %w", err)
	}

	lggr.Debugw("calling receiver", "paramTypes", paramTypes, "paramValues", paramValues)

	for _, objectId := range receiverObjectIdStrings {
		paramValues = append(paramValues, bind.Object{Id: objectId})
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
	messageID [32]byte,
	normalizedModule *models.GetNormalizedMoveModuleResponse,
	receiverParams *transaction.Argument,
	extraArgs map[string]any,
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
	paramTypes := []string{
		"&mut object",
	}
	paramValues := []any{
		receiverParams,
	}

	lggr.Debugw("calling offramp state helper to extract any2sui message", "paramTypes", paramTypes, "paramValues", paramValues)

	encodedAny2SuiExtractCall, err := offrampStateHelperContract.EncodeCallArgsWithGenerics(
		"extract_any2sui_message",
		typeArgsList,
		typeParamsList,
		paramTypes,
		paramValues,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to encode extract_any2sui_message call: %w", err)
	}

	extractedAny2SuiMessageResult, err := offrampStateHelperContract.AppendPTB(ctx, callOpts, ptb, encodedAny2SuiExtractCall)
	if err != nil {
		return nil, fmt.Errorf("failed to build PTB (extract_any2sui_message) using bindings: %w", err)
	}

	return appendCcipReceiveCommand(
		ctx,
		lggr,
		ptb,
		callOpts,
		boundReceiverContract,
		functionName,
		addressMappings,
		messageID,
		normalizedModule,
		extractedAny2SuiMessageResult,
		extractReceiverObjectIdStrings(extraArgs),
	)
}

// ProcessReceiversV2 processes receiver calls using V2 protocol with object binding enforcement.
// It calls extract_any2sui_message_v2 (which verifies receiver_object_ids match the committed values)
// before invoking the receiver's ccip_receive function.
func ProcessReceiversV2(
	ctx context.Context,
	lggr logger.Logger,
	ptbClient client.SuiPTBClient,
	ptb *transaction.Transaction,
	messages []ccipocr3.Message,
	addressMappings *OffRampAddressMappings,
	callOpts *bind.CallOpts,
	receiverParams *transaction.Argument,
	extraArgs map[string]any,
) ([]transaction.Argument, error) {
	sdkClient := ptbClient.GetClient()

	receiverRegistryPkg, err := receiver_registry.NewReceiverRegistry(addressMappings.CcipPackageId, sdkClient)
	if err != nil {
		return nil, err
	}
	receiverRegistryDevInspect := receiverRegistryPkg.DevInspect()

	receiverCommandsResults := make([]transaction.Argument, 0)
	for _, message := range messages {
		if len(message.Receiver) == 0 || message.Receiver == nil {
			lggr.Errorw("unexpected nil or zero length receiver, skipping message in offramp execution...", "message", message)
			continue
		}

		if bytes.Equal(message.Receiver, codec.AccountZero) {
			lggr.Debugw("receiver is zero address, skipping message in offramp execution...", "message", message)
			continue
		}

		if !needsAppDelivery(message, extraArgs) {
			lggr.Debugw("message has no data and zero gas limit, skipping receiver call",
				"receiver", hex.EncodeToString(message.Receiver))
			continue
		}

		receiverPackageId := "0x" + hex.EncodeToString(message.Receiver)

		isRegistered, err := receiverRegistryDevInspect.IsRegisteredReceiver(ctx, callOpts, bind.Object{Id: addressMappings.CcipObjectRef}, receiverPackageId)
		if err != nil {
			return nil, fmt.Errorf("failed to check if receiver is registered in offramp execution: %w", err)
		}
		if !isRegistered {
			lggr.Warnw("receiver not registered, skipping receiver call (on-chain will not populate message)",
				"receiver", receiverPackageId)
			continue
		}

		receiverConfig, err := receiverRegistryDevInspect.GetReceiverConfig(ctx, callOpts, bind.Object{Id: addressMappings.CcipObjectRef}, receiverPackageId)
		if err != nil {
			return nil, fmt.Errorf("failed to get receiver config in offramp execution: %w", err)
		}

		receiverNormalizedModule, err := ptbClient.GetNormalizedModule(ctx, receiverPackageId, receiverConfig.ModuleName)
		if err != nil {
			return nil, fmt.Errorf("failed to get normalized module for receiver: %w", err)
		}

		receiverCommandResult, err := AppendPTBCommandForReceiverV2(
			ctx,
			lggr,
			sdkClient,
			ptb,
			callOpts,
			receiverPackageId,
			receiverConfig.ModuleName,
			"ccip_receive",
			addressMappings,
			message.Header.MessageID,
			&receiverNormalizedModule,
			receiverParams,
			extraArgs,
		)
		skip, retErr := classifyReceiverBuildError(receiverPackageId, err)
		if skip {
			lggr.Errorw("skipping receiver command due to unsupported ABI; PTB will fail on-chain",
				"receiver", receiverPackageId,
				"error", err)
			continue
		}
		if retErr != nil {
			return nil, retErr
		}
		receiverCommandsResults = append(receiverCommandsResults, *receiverCommandResult)
	}

	return receiverCommandsResults, nil
}

// AppendPTBCommandForReceiverV2 builds the receiver call using V2 extract which enforces
// receiver_object_ids binding at the protocol level.
func AppendPTBCommandForReceiverV2(
	ctx context.Context,
	lggr logger.Logger,
	sdkClient sui.ISuiAPI,
	ptb *transaction.Transaction,
	callOpts *bind.CallOpts,
	packageId string,
	moduleId string,
	functionName string,
	addressMappings *OffRampAddressMappings,
	messageID [32]byte,
	normalizedModule *models.GetNormalizedMoveModuleResponse,
	receiverParams *transaction.Argument,
	extraArgs map[string]any,
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

	receiverObjectIdStrings := extractReceiverObjectIdStrings(extraArgs)

	// Call extract_any2sui_message_v2 with receiver_object_ids for protocol-level binding enforcement
	typeArgsList := []string{}
	typeParamsList := []string{}
	paramTypes := []string{
		"&mut object",
		"vector<object_id>",
	}
	paramValues := []any{
		receiverParams,
		receiverObjectIdStrings,
	}

	lggr.Debugw("calling extract_any2sui_message_v2",
		"receiverObjectIds", receiverObjectIdStrings)

	encodedAny2SuiExtractCall, err := offrampStateHelperContract.EncodeCallArgsWithGenerics(
		"extract_any2sui_message_v2",
		typeArgsList,
		typeParamsList,
		paramTypes,
		paramValues,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to encode extract_any2sui_message_v2 call: %w", err)
	}

	extractedAny2SuiMessageResult, err := offrampStateHelperContract.AppendPTB(ctx, callOpts, ptb, encodedAny2SuiExtractCall)
	if err != nil {
		return nil, fmt.Errorf("failed to build PTB (extract_any2sui_message_v2) using bindings: %w", err)
	}

	return appendCcipReceiveCommand(
		ctx,
		lggr,
		ptb,
		callOpts,
		boundReceiverContract,
		functionName,
		addressMappings,
		messageID,
		normalizedModule,
		extractedAny2SuiMessageResult,
		receiverObjectIdStrings,
	)
}
