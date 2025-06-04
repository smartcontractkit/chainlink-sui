package chainreader

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/pattonkan/sui-go/suiclient"

	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/codec"

	"github.com/smartcontractkit/chainlink-common/pkg/services"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	pkgtypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
)

const (
	defaultQueryLimit   = 25
	readIdentifierParts = 3
	objectIdPrefix      = "0x"
)

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

func (s *suiChainReader) Bind(ctx context.Context, bindings []pkgtypes.BoundContract) error {
	newBindings := map[string]string{}
	for _, binding := range bindings {
		if !strings.HasPrefix(binding.Address, objectIdPrefix) {
			return fmt.Errorf("invalid Sui package address format: %s", binding.Address)
		}
		newBindings[binding.Name] = binding.Address
	}

	for name, address := range newBindings {
		s.packageAddresses[name] = address
	}

	return nil
}

func (s *suiChainReader) Unbind(ctx context.Context, bindings []pkgtypes.BoundContract) error {
	for _, binding := range bindings {
		if _, ok := s.packageAddresses[binding.Name]; !ok {
			return fmt.Errorf("no such binding: %s", binding.Name)
		}
		delete(s.packageAddresses, binding.Name)
	}

	return nil
}

// readIdentifier represents the parsed components of a read identifier
type readIdentifier struct {
	address      string
	contractName string
	readName     string
}

// parseReadIdentifier parses a read identifier string into its components
func (s *suiChainReader) parseReadIdentifier(identifier string) (*readIdentifier, error) {
	components := strings.Split(identifier, "-")
	if len(components) != readIdentifierParts {
		return nil, fmt.Errorf("invalid read identifier format: %s (expected format: address-contract-readName)", identifier)
	}

	return &readIdentifier{
		address:      components[0],
		contractName: components[1],
		readName:     components[2],
	}, nil
}

// validateBinding validates that the contract is bound and addresses match
func (s *suiChainReader) validateBinding(parsed *readIdentifier) error {
	boundAddress, ok := s.packageAddresses[parsed.contractName]
	if !ok {
		return fmt.Errorf("no bound address for contract: %s", parsed.contractName)
	}

	if boundAddress != parsed.address {
		return fmt.Errorf("bound address %s for contract %s does not match read address %s",
			boundAddress, parsed.contractName, parsed.address)
	}

	if _, ok := s.config.Modules[parsed.contractName]; !ok {
		return fmt.Errorf("no configuration for contract: %s", parsed.contractName)
	}

	return nil
}

// isObjectRead determines if the read operation is for an object (starts with 0x) or a function call
func isObjectRead(readName string) bool {
	return strings.HasPrefix(readName, objectIdPrefix)
}

// GetLatestValue retrieves the latest value from either an object or function call
func (s *suiChainReader) GetLatestValue(ctx context.Context, readIdentifier string, confidenceLevel primitives.ConfidenceLevel, params, returnVal any) error {
	parsed, err := s.parseReadIdentifier(readIdentifier)
	if err != nil {
		return err
	}

	if err = s.validateBinding(parsed); err != nil {
		return err
	}

	var valueField any

	if isObjectRead(parsed.readName) {
		valueField, err = s.readObject(ctx, parsed.readName)
	} else {
		valueField, err = s.callFunction(ctx, parsed, params)
	}

	if err != nil {
		return err
	}

	return s.encodeResult(valueField, returnVal)
}

// readObject reads a value from a Sui object
func (s *suiChainReader) readObject(ctx context.Context, objectId string) (any, error) {
	object, err := s.client.ReadObjectId(ctx, objectId)
	if err != nil {
		return nil, fmt.Errorf("failed to read object %s: %w", objectId, err)
	}

	valueField, ok := object["value"]
	if !ok {
		return nil, fmt.Errorf("object %s does not contain a 'value' field", objectId)
	}

	return valueField, nil
}

// callFunction calls a contract function and returns the result
func (s *suiChainReader) callFunction(ctx context.Context, parsed *readIdentifier, params any) (any, error) {
	moduleConfig := s.config.Modules[parsed.contractName]
	functionConfig, ok := moduleConfig.Functions[parsed.readName]
	if !ok {
		return nil, fmt.Errorf("no function configuration for: %s", parsed.readName)
	}

	argMap, err := s.parseParams(params, functionConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to parse parameters: %w", err)
	}

	args, argTypes, err := s.prepareArguments(argMap, functionConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare arguments: %w", err)
	}

	response, err := s.executeFunction(ctx, parsed, moduleConfig, functionConfig, args, argTypes)
	if err != nil {
		return nil, err
	}

	return s.parseResponse(response.ReturnValues[0])
}

// parseParams parses input parameters based on whether we're running as a LOOP plugin
func (s *suiChainReader) parseParams(params any, functionConfig *ChainReaderFunction) (map[string]any, error) {
	argMap := make(map[string]any)

	if s.config.IsLoopPlugin {
		return s.parseLoopParams(params, functionConfig)
	}

	if err := mapstructure.Decode(params, &argMap); err != nil {
		return nil, fmt.Errorf("failed to decode parameters: %w", err)
	}

	return argMap, nil
}

// parseLoopParams handles parameter parsing for LOOP plugin mode
func (s *suiChainReader) parseLoopParams(params any, functionConfig *ChainReaderFunction) (map[string]any, error) {
	paramBytes, ok := params.(*[]byte)
	if !ok {
		return nil, fmt.Errorf("expected *[]byte for LOOP plugin params, got %T", params)
	}

	decoder := json.NewDecoder(bytes.NewReader(*paramBytes))
	decoder.UseNumber()

	var rawArgMap map[string]any
	if err := decoder.Decode(&rawArgMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON params: %w", err)
	}

	// Convert JSON-unmarshaled values to proper Go types
	argMap := make(map[string]any)
	if functionConfig.Params != nil {
		for _, paramConfig := range functionConfig.Params {
			if jsonValue, exists := rawArgMap[paramConfig.Name]; exists {
				convertedValue, err := codec.EncodeToSuiValue(paramConfig.Type, jsonValue)
				if err != nil {
					return nil, fmt.Errorf("failed to convert parameter %s of type %s: %w",
						paramConfig.Name, paramConfig.Type, err)
				}
				argMap[paramConfig.Name] = convertedValue
			}
		}
	}

	return argMap, nil
}

// prepareArguments prepares function arguments and types for the call
func (s *suiChainReader) prepareArguments(argMap map[string]any, functionConfig *ChainReaderFunction) ([]any, []string, error) {
	if functionConfig.Params == nil {
		return []any{}, []string{}, nil
	}

	args := make([]any, 0, len(functionConfig.Params))
	argTypes := make([]string, 0, len(functionConfig.Params))

	for _, paramConfig := range functionConfig.Params {
		argValue, ok := argMap[paramConfig.Name]
		if !ok {
			if paramConfig.Required {
				return nil, nil, fmt.Errorf("missing required argument: %s", paramConfig.Name)
			}
			argValue = paramConfig.DefaultValue
		}

		args = append(args, argValue)
		argTypes = append(argTypes, paramConfig.Type)
	}

	return args, argTypes, nil
}

// executeFunction executes the actual function call
func (s *suiChainReader) executeFunction(ctx context.Context, parsed *readIdentifier, moduleConfig *ChainReaderModule, functionConfig *ChainReaderFunction, args []any, argTypes []string) (*suiclient.ExecutionResultType, error) {
	s.logger.Debugw("Calling ReadFunction",
		"address", parsed.address,
		"module", moduleConfig.Name,
		"method", parsed.readName,
		"encodedArgs", args,
		"argTypes", argTypes,
	)

	response, err := s.client.ReadFunction(ctx, functionConfig.SignerAddress, parsed.address, moduleConfig.Name, parsed.readName, args, argTypes)
	if err != nil {
		s.logger.Errorw("ReadFunction failed",
			"error", err,
			"address", parsed.address,
			"module", moduleConfig.Name,
			"method", parsed.readName,
			"args", args,
			"argTypes", argTypes,
		)

		return nil, fmt.Errorf("failed to call function %s: %w", parsed.readName, err)
	}

	s.logger.Debugw("Sui ReadFunction response", "response", response.ReturnValues[0])

	return response, nil
}

// parseResponse parses the function response based on plugin mode
func (s *suiChainReader) parseResponse(rawResponse any) (any, error) {
	if s.config.IsLoopPlugin {
		valueField, err := codec.ParseSuiResponseValue(rawResponse)
		if err != nil {
			return nil, fmt.Errorf("failed to parse Sui response: %w", err)
		}
		s.logger.Debugw("Sui ParseSuiResponseValue", "valueField", valueField)

		return valueField, nil
	}

	// For non-LOOP mode, we'll parse the response when encoding the result
	return rawResponse, nil
}

// encodeResult encodes the final result based on plugin mode
func (s *suiChainReader) encodeResult(valueField any, returnVal any) error {
	if s.config.IsLoopPlugin {
		return s.encodeLoopResult(valueField, returnVal)
	}

	// For non-LOOP mode, handle both parsed responses and direct values
	if responseArray, ok := valueField.([]any); ok && len(responseArray) >= 2 {
		// This is a raw function response that needs parsing
		return codec.ParseSuiResponseValueWithTarget(valueField, returnVal)
	}

	// This is already a parsed value (from object read)
	return codec.DecodeSuiJsonValue(valueField, returnVal)
}

// encodeLoopResult encodes results for LOOP plugin mode
func (s *suiChainReader) encodeLoopResult(valueField any, returnVal any) error {
	resultBytes, err := json.Marshal(valueField)
	if err != nil {
		return fmt.Errorf("failed to marshal data for LOOP: %w", err)
	}

	returnValPtr, ok := returnVal.(*[]byte)
	if !ok {
		return fmt.Errorf("return value is not a pointer to []byte as expected when running as a LOOP plugin")
	}

	*returnValPtr = make([]byte, len(resultBytes))
	copy(*returnValPtr, resultBytes)

	return nil
}

func (s *suiChainReader) BatchGetLatestValues(ctx context.Context, request pkgtypes.BatchGetLatestValuesRequest) (pkgtypes.BatchGetLatestValuesResult, error) {
	result := make(pkgtypes.BatchGetLatestValuesResult)

	for contract, batch := range request {
		batchResults := make(pkgtypes.ContractBatchResults, len(batch))
		resultChan := make(chan struct {
			index  int
			result pkgtypes.BatchReadResult
		}, len(batch))

		for i, read := range batch {
			go func(index int, read pkgtypes.BatchRead) {
				readResult := pkgtypes.BatchReadResult{ReadName: read.ReadName}

				err := s.GetLatestValue(ctx, contract.ReadIdentifier(read.ReadName), primitives.Finalized, read.Params, read.ReturnVal)
				readResult.SetResult(read.ReturnVal, err)

				select {
				case resultChan <- struct {
					index  int
					result pkgtypes.BatchReadResult
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

func (s *suiChainReader) QueryKey(ctx context.Context, contract pkgtypes.BoundContract, filter query.KeyFilter, limitAndSort query.LimitAndSort, sequenceDataType any) ([]pkgtypes.Sequence, error) {
	// Validate contract binding
	if err := s.validateContractBinding(contract); err != nil {
		return nil, err
	}

	// Get module and event configuration
	moduleConfig := s.config.Modules[contract.Name]
	eventConfig, err := s.getEventConfig(moduleConfig, filter.Key)
	if err != nil {
		return nil, err
	}

	// Query events from blockchain
	eventsResponse, err := s.queryEvents(ctx, contract, eventConfig, limitAndSort)
	if err != nil {
		return nil, err
	}

	// Transform events to sequences
	return s.transformEventsToSequences(ctx, eventsResponse, sequenceDataType)
}

// validateContractBinding validates the contract binding for QueryKey
func (s *suiChainReader) validateContractBinding(contract pkgtypes.BoundContract) error {
	address, ok := s.packageAddresses[contract.Name]
	if !ok {
		return fmt.Errorf("no bound address for package %s", contract.Name)
	}

	if address != contract.Address {
		return fmt.Errorf("bound address %s for package %s does not match provided address %s",
			address, contract.Name, contract.Address)
	}

	if _, ok := s.config.Modules[contract.Name]; !ok {
		return fmt.Errorf("no configuration for contract: %s", contract.Name)
	}

	return nil
}

// getEventConfig retrieves event configuration for the given key
func (s *suiChainReader) getEventConfig(moduleConfig *ChainReaderModule, eventKey string) (*ChainReaderEvent, error) {
	if moduleConfig.Events == nil {
		return nil, fmt.Errorf("no events configured for contract")
	}

	eventConfig, ok := moduleConfig.Events[eventKey]
	if !ok {
		return nil, fmt.Errorf("no configuration for event: %s", eventKey)
	}

	return eventConfig, nil
}

// queryEvents queries events from the Sui blockchain
func (s *suiChainReader) queryEvents(ctx context.Context, contract pkgtypes.BoundContract, eventConfig *ChainReaderEvent, limitAndSort query.LimitAndSort) (*suiclient.EventPage, error) {
	limit := uint(defaultQueryLimit)
	if limitAndSort.Limit.Count > 0 {
		limit = uint(limitAndSort.Limit.Count)
	}

	descending := true
	for _, sortBy := range limitAndSort.SortBy {
		if seqSort, ok := sortBy.(query.SortBySequence); ok && seqSort.GetDirection() == query.Asc {
			descending = false
			break
		}
	}

	var cursor *client.EventId
	if limitAndSort.Limit.Cursor != "" {
		if err := json.Unmarshal([]byte(limitAndSort.Limit.Cursor), &cursor); err != nil {
			return nil, fmt.Errorf("failed to unmarshal cursor: %w", err)
		}
	}

	eventsResponse, err := s.client.QueryEvents(
		ctx,
		client.EventFilterByMoveEventModule{
			Package: contract.Address,
			Module:  contract.Name,
			Event:   eventConfig.EventType,
		},
		&limit,
		cursor,
		&client.QuerySortOptions{Descending: descending},
	)
	if err != nil {
		s.logger.Errorw("Failed to query events",
			"error", err,
			"address", contract.Address,
			"module", contract.Name,
			"eventType", eventConfig.EventType,
			"limit", limit,
		)

		return nil, fmt.Errorf("failed to query events: %w", err)
	}

	return eventsResponse, nil
}

// transformEventsToSequences converts blockchain events to sequence format
func (s *suiChainReader) transformEventsToSequences(ctx context.Context, eventsResponse *suiclient.EventPage, sequenceDataType any) ([]pkgtypes.Sequence, error) {
	sequences := make([]pkgtypes.Sequence, 0, len(eventsResponse.Data))

	for _, event := range eventsResponse.Data {
		eventData := reflect.New(reflect.TypeOf(sequenceDataType).Elem()).Interface()

		s.logger.Debugw("Processing event", "ParsedJson", event.ParsedJson)

		if err := codec.DecodeSuiJsonValue(event.ParsedJson, eventData); err != nil {
			return nil, fmt.Errorf("failed to decode event data: %w", err)
		}

		marshalledCursor, err := json.Marshal(client.EventId{
			TxDigest: event.Id.TxDigest.String(),
			EventSeq: event.Id.EventSeq,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to marshal event cursor: %w", err)
		}

		block, err := s.client.BlockByDigest(ctx, event.Id.TxDigest.String())
		if err != nil {
			return nil, fmt.Errorf("failed to get block by digest: %w", err)
		}

		sequence := pkgtypes.Sequence{
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
