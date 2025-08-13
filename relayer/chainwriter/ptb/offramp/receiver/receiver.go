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
	signerAddress string,
	messages []ccipocr3.Message,
	previousCommandIndex uint16,
	ccipObjectRef string,
	ccipPackageId string,
	client ptbClient.SuiPTBClient,
) ([]transaction.Argument, error) {
	suiClient := client.GetClient()
	devInspectSigner := signer.NewDevInspectSigner(signerAddress)
	registeredReceivers, err := FilterRegisteredReceivers(ctx, lggr, messages, signerAddress, client, ccipObjectRef, ccipPackageId)
	if err != nil {
		return nil, err
	}

	lggr.Info("registered receivers", "count", len(registeredReceivers))

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

		signature, err := GetCCIPReceiverSignature(ctx, lggr, signerAddress, receiverPackageId, moduleName, client)
		if err != nil {
			return nil, err
		}

		lggr.Infow("signature", "signature", signature)
	}

	return []transaction.Argument{}, nil
}

func GetCCIPReceiverSignature(
	ctx context.Context,
	lggr logger.Logger,
	signerAddress string,
	receiverPackageId string,
	receiverModule string,
	client ptbClient.SuiPTBClient,
) (string, error) {
	lggr.Infow("getting ccip receiver signature", "receiverPackageId", receiverPackageId, "receiverModule", receiverModule)
	normalizedModule, err := client.GetNormalizedModule(ctx, receiverPackageId, receiverModule)
	if err != nil {
		lggr.Error("failed to get normalized module", "error", err)
		return "", err
	}

	functions := normalizedModule.ExposedFunctions
	if functions[CCIP_RECEIVER_FUNCTION] == nil {
		lggr.Error("ccip_receive function not found", "receiverPackageId", receiverPackageId)
		return "", fmt.Errorf("ccip_receive function not found: %s", receiverPackageId)
	}

	function := functions[CCIP_RECEIVER_FUNCTION].(map[string]any)
	lggr.Infow("function", "function", function)

	return "", nil
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
