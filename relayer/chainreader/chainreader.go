package chainreader

import (
	"context"
	"fmt"
	"github.com/smartcontractkit/chainlink-internal-integrations/sui/relayer/codec"
	"strings"

	"github.com/smartcontractkit/chainlink-common/pkg/services"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/mitchellh/mapstructure"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	pkgtypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
	// TODO: enable after codec and txm is implemented
	// "github.com/smartcontractkit/chainlink-internal-integrations/sui/relayer/txm"
)

type suiChainReader struct {
	pkgtypes.UnimplementedContractReader

	logger           logger.Logger
	config           ChainReaderConfig
	starter          services.StateMachine
	packageAddresses map[string]string
	client           sui.ISuiAPI
}

func NewChainReader(lgr logger.Logger, client sui.ISuiAPI, config ChainReaderConfig) pkgtypes.ContractReader {
	return &suiChainReader{
		logger:           logger.Named(lgr, "SuiChainReader"),
		client:           client,
		config:           config,
		packageAddresses: map[string]string{},
	}
}

func (s *suiChainReader) Name() string {
	return s.logger.Name()
}

func (s *suiChainReader) Ready() error {
	return s.starter.Ready()
}

func (s *suiChainReader) HealthReport() map[string]error {
	return map[string]error{s.Name(): s.starter.Healthy()}
}

func (s *suiChainReader) Start(ctx context.Context) error {
	return s.starter.StartOnce(s.Name(), func() error {
		return nil
	})
}

func (s *suiChainReader) Close() error {
	return s.starter.StopOnce(s.Name(), func() error {
		return nil
	})
}

func (s *suiChainReader) Bind(ctx context.Context, bindings []types.BoundContract) error {
	newBindings := map[string]string{}
	for _, binding := range bindings {
		// In Sui, we don't need to parse addresses, they're already in the correct format
		if !strings.HasPrefix(binding.Address, "0x") {
			return fmt.Errorf("invalid Sui package address format: %s", binding.Address)
		}

		newBindings[binding.Name] = binding.Address
	}

	for name, address := range newBindings {
		s.packageAddresses[name] = address
	}

	return nil
}

func (s *suiChainReader) Unbind(ctx context.Context, bindings []types.BoundContract) error {
	for _, binding := range bindings {
		key := binding.Name
		if _, ok := s.packageAddresses[key]; ok {
			delete(s.packageAddresses, key)
		} else {
			return fmt.Errorf("no such binding: %s", key)
		}
	}
	return nil
}

func (s *suiChainReader) GetLatestValue(ctx context.Context, readIdentifier string, confidenceLevel primitives.ConfidenceLevel, params, returnVal any) error {
	// Decode the readIdentifier - a combination of address, contract, and readName as a concatenated string
	readComponents := strings.Split(readIdentifier, "-")
	if len(readComponents) != 3 {
		return fmt.Errorf("invalid read identifier: %s", readIdentifier)
	}
	_address, contractName, method := readComponents[0], readComponents[1], readComponents[2]

	// Source the read configuration, by contract name
	address, ok := s.packageAddresses[contractName]
	if !ok {
		return fmt.Errorf("no bound address for package %s", contractName)
	}

	// Notice: the address in the readIdentifier should match the bound address, by contract name
	if address != _address {
		return fmt.Errorf("bound address %s for package %s does not match read address %s", address, contractName, _address)
	}

	moduleConfig, ok := s.config.Modules[contractName]
	if !ok {
		return fmt.Errorf("no such contract: %s", contractName)
	}

	functionConfig, ok := moduleConfig.Functions[method]
	if !ok {
		return fmt.Errorf("no such method: %s", method)
	}

	// Extract parameters from the params object
	argMap := make(map[string]interface{})
	if err := mapstructure.Decode(params, &argMap); err != nil {
		return fmt.Errorf("failed to parse arguments: %w", err)
	}

	// Prepare arguments for the function call
	args := []interface{}{}

	if functionConfig.Params != nil {
		for _, paramConfig := range functionConfig.Params {
			argValue, ok := argMap[paramConfig.Name]
			if !ok {
				if paramConfig.Required {
					return fmt.Errorf("missing argument: %s", paramConfig.Name)
				}
				argValue = paramConfig.DefaultValue
			}

			// No need for BCS serialization in Sui calls via JSON-RPC
			args = append(args, argValue)
		}
	}

	var moduleName string
	if moduleConfig.Name != "" {
		moduleName = moduleConfig.Name
	} else {
		moduleName = contractName
	}

	var functionName string
	if functionConfig.Name != "" {
		functionName = functionConfig.Name
	} else {
		functionName = method
	}

	// Call the move function
	resp, err := s.client.MoveCall(ctx, models.MoveCallRequest{
		PackageObjectId: address,
		Module:          moduleName,
		Function:        functionName,
		TypeArguments:   []interface{}{},
		Arguments:       args,
	})

	if err != nil {
		return fmt.Errorf("failed to call view function: %w", err)
	}

	fmt.Print(resp)

	// TODO: objectId is currently assumed to me in the `method` section of the `readIdentifier`
	object, err := s.client.SuiGetObject(ctx, models.SuiGetObjectRequest{
		ObjectId: method,
		Options: models.SuiObjectDataOptions{
			ShowContent: true,
		},
	})

	if err != nil {
		return fmt.Errorf("failed to get object: %w", err)
	}

	// Decode the return value into the provided structure
	return codec.DecodeSuiJsonValue(object.Data.Content, returnVal)
}

func (s *suiChainReader) BatchGetLatestValues(ctx context.Context, request types.BatchGetLatestValuesRequest) (types.BatchGetLatestValuesResult, error) {
	// not implemented
	return types.BatchGetLatestValuesResult{}, nil
}

func (s *suiChainReader) QueryKey(ctx context.Context, contract types.BoundContract, filter query.KeyFilter, limitAndSort query.LimitAndSort, sequenceDataType any) ([]types.Sequence, error) {
	// not implemented
	return nil, nil
}
