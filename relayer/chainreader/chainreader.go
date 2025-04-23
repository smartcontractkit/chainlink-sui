package chainreader

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/mitchellh/mapstructure"

	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/codec"

	"github.com/smartcontractkit/chainlink-common/pkg/services"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	pkgtypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
)

const defaultQueryLimit = 25

type suiChainReader struct {
	pkgtypes.UnimplementedContractReader

	logger           logger.Logger
	config           ChainReaderConfig
	starter          services.StateMachine
	packageAddresses map[string]string
	client           client.PTBClient
}

func NewChainReader(lgr logger.Logger, abstractClient client.PTBClient, config ChainReaderConfig) pkgtypes.ContractReader {
	return &suiChainReader{
		logger:           logger.Named(lgr, "SuiChainReader"),
		client:           abstractClient,
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

		if _, ok := s.packageAddresses[key]; !ok {
			return fmt.Errorf("no such binding: %s", key)
		}
		delete(s.packageAddresses, key)
	}

	return nil
}

// GetLatestValue A method to get the latest value of an object managed by one of the contracts in the Sui network integration.
// Note that the `readIdentifier` here is split into 3 parts in a `-` delimited string. The third part being the Object ID.
func (s *suiChainReader) GetLatestValue(ctx context.Context, readIdentifier string, confidenceLevel primitives.ConfidenceLevel, params, returnVal any) error {
	// Decode the readIdentifier - a combination of address, contract, and readName as a concatenated string
	readComponents := strings.Split(readIdentifier, "-")
	expectedComponents := 3
	if len(readComponents) != expectedComponents {
		return fmt.Errorf("invalid read identifier: %s", readIdentifier)
	}
	_address, contractName, objectIdOrFunction := readComponents[0], readComponents[1], readComponents[2]

	// Source the read configuration, by contract name
	address, ok := s.packageAddresses[contractName]
	if !ok {
		return fmt.Errorf("no bound address for package %s", contractName)
	}

	// The address in the readIdentifier should match the bound address, by contract name
	if address != _address {
		return fmt.Errorf("bound address %s for package %s does not match read address %s", address, contractName, _address)
	}

	_, ok = s.config.Modules[contractName]
	if !ok {
		return fmt.Errorf("no such contract: %s", contractName)
	}

	var valueField any
	// Since the last part of the readIdentifier can be either a function or an object ID, we need to check to determine
	// how to proceed to get the value.
	if strings.HasPrefix(objectIdOrFunction, "0x") {
		objectId := objectIdOrFunction

		object, err := s.client.ReadObjectId(ctx, objectId)
		if err != nil {
			return fmt.Errorf("failed to get object: %w", err)
		}

		// Extract the value field from the object
		valueField, ok = object["value"]
		if !ok {
			return fmt.Errorf("object does not contain a 'value' field")
		}
	} else {
		method := objectIdOrFunction
		// We need to call the function from the contract
		moduleConfig, ok := s.config.Modules[contractName]
		if !ok {
			return fmt.Errorf("no such contract: %s", contractName)
		}

		functionConfig, ok := moduleConfig.Functions[method]
		if !ok {
			return fmt.Errorf("no such method: %s", method)
		}

		// Extract parameters from the params object
		argMap := make(map[string]any)
		if err := mapstructure.Decode(params, &argMap); err != nil {
			return fmt.Errorf("failed to parse arguments: %w", err)
		}

		// Prepare arguments for the function call
		args := []any{}
		argTypes := []string{}

		if functionConfig.Params != nil {
			for _, paramConfig := range functionConfig.Params {
				argValue, ok := argMap[paramConfig.Name]
				if !ok {
					if paramConfig.Required {
						return fmt.Errorf("missing argument: %s", paramConfig.Name)
					}
					argValue = paramConfig.DefaultValue
				}

				args = append(args, argValue)
				argTypes = append(argTypes, paramConfig.Type)
			}
		}

		s.logger.Debugw("Calling ReadFunction",
			"address", address,
			"module", moduleConfig.Name,
			"method", method,
			"encodedArgs", args,
			"argTypes", argTypes,
		)

		response, err := s.client.ReadFunction(ctx, address, moduleConfig.Name, method, args, argTypes)
		if err != nil {
			s.logger.Errorw("ReadFunction failed",
				"error", err,
				"address", address,
				"module", moduleConfig.Name,
				"method", method,
				"args", args,
				"argTypes", argTypes,
			)

			return fmt.Errorf("failed to call function: %w", err)
		}

		s.logger.Debugw("Sui ReadFunction", "response", response.ReturnValues[0])

		// Extract the array from the response
		rawArray := response.ReturnValues[0].([]any)
		s.logger.Debugw("Raw array value", "array", rawArray)

		// TODO: move this into a helper when merging code with bindings
		// Convert the array of interface{} to []byte
		byteArray := make([]byte, len(rawArray))
		for i, v := range rawArray {
			// Convert each element to byte
			if num, ok := v.(float64); ok {
				byteArray[i] = byte(num)
			}
		}

		valueField = byteArray
	}

	// Decode the return value into the provided structure
	return codec.DecodeSuiJsonValue(valueField, returnVal)
}

func (s *suiChainReader) BatchGetLatestValues(ctx context.Context, request types.BatchGetLatestValuesRequest) (types.BatchGetLatestValuesResult, error) {
	result := make(types.BatchGetLatestValuesResult)

	for contract, batch := range request {
		batchResults := make(types.ContractBatchResults, len(batch))
		resultChan := make(chan struct {
			index  int
			result types.BatchReadResult
		}, len(batch))

		for i, read := range batch {
			go func(index int, read types.BatchRead) {
				readResult := types.BatchReadResult{ReadName: read.ReadName}

				err := s.GetLatestValue(ctx, contract.ReadIdentifier(read.ReadName), primitives.Finalized, read.Params, read.ReturnVal)
				readResult.SetResult(read.ReturnVal, err)

				select {
				case resultChan <- struct {
					index  int
					result types.BatchReadResult
				}{index, readResult}:
				case <-ctx.Done():
					return
				}
			}(i, read)
		}

		for range batch {
			select {
			case res := <-resultChan:
				batchResults[res.index] = res.result
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}

		result[contract] = batchResults
	}

	return result, nil
}

func (s *suiChainReader) QueryKey(ctx context.Context, contract types.BoundContract, filter query.KeyFilter, limitAndSort query.LimitAndSort, sequenceDataType any) ([]types.Sequence, error) {
	// Validate contract has a bound address
	address, ok := s.packageAddresses[contract.Name]
	if !ok {
		return nil, fmt.Errorf("no bound address for package %s", contract.Name)
	}

	// Check that the address in the contract matches the bound address
	if address != contract.Address {
		return nil, fmt.Errorf("bound address %s for package %s does not match provided address %s", address, contract.Name, contract.Address)
	}

	// Check for module configuration
	moduleConfig, ok := s.config.Modules[contract.Name]
	if !ok {
		return nil, fmt.Errorf("no such contract: %s", contract.Name)
	}

	// Extract event field name from the filter
	eventFieldName := filter.Key

	// Check event configuration
	eventConfig, ok := moduleConfig.Events[eventFieldName]
	if !ok {
		return nil, fmt.Errorf("no such event: %s", eventFieldName)
	}

	// Extract limit from limitAndSort
	limit := uint(defaultQueryLimit)
	if count := limitAndSort.Limit.Count; count > 0 {
		limit = uint(count)
	}

	// Determine sorting direction
	descending := true
	for _, sortBy := range limitAndSort.SortBy {
		if seqSort, ok := sortBy.(query.SortBySequence); ok {
			if seqSort.GetDirection() == query.Asc {
				descending = false
				break
			}
		}
	}

	var cursor *client.EventId
	if limitAndSort.Limit.Cursor != "" {
		// unmarshal the cursor
		err := json.Unmarshal([]byte(limitAndSort.Limit.Cursor), &cursor)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal cursor: %w", err)
		}
	} else {
		cursor = nil
	}

	// Query events from Sui blockchain
	eventsResponse, err := s.client.QueryEvents(
		ctx,
		client.EventFilterByMoveEventModule{
			Package: address,
			Module:  contract.Name,
			Event:   eventConfig.EventType,
		},
		&limit,
		cursor,
		&client.QuerySortOptions{
			Descending: descending,
		},
	)
	if err != nil {
		s.logger.Errorw("Failed to query events",
			"error", err,
			"address", address,
			"module", moduleConfig.Name,
			"eventType", eventConfig.EventType,
			"limit", limit,
		)

		return nil, fmt.Errorf("failed to query events: %+w", err)
	}

	// Transform events into the expected Sequence format
	sequences := make([]types.Sequence, 0, len(eventsResponse.Data))
	for _, event := range eventsResponse.Data {
		// Create new instance of eventData for each event
		eventData := reflect.New(reflect.TypeOf(sequenceDataType).Elem()).Interface()

		s.logger.Debugw("event", "ParsedJson", event.ParsedJson)

		// Decode the event data
		err := codec.DecodeSuiJsonValue(event.ParsedJson, eventData)
		if err != nil {
			return nil, fmt.Errorf("failed to decode event data: %+w", err)
		}

		marshalledCursor, err := json.Marshal(client.EventId{
			TxDigest: event.Id.TxDigest.String(),
			EventSeq: event.Id.EventSeq,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to marshal event cursor: %+w", err)
		}

		block, err := s.client.BlockByDigest(ctx, event.Id.TxDigest.String())
		if err != nil {
			return nil, fmt.Errorf("failed to get block by digest: %+w", err)
		}

		sequence := types.Sequence{
			Cursor: string(marshalledCursor),
			Data:   eventData,
			Head: pkgtypes.Head{
				Timestamp: block.Timestamp,
				Hash:      []byte(block.TxDigest),
				Height:    strconv.FormatUint(block.Height, 10),
			},
		}
		sequences = append(sequences, sequence)
	}

	return sequences, nil
}
