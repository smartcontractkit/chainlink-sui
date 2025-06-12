package chainreader

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/pattonkan/sui-go/suiclient"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"

	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/database"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/indexer"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/codec"

	"maps"

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

	logger                    logger.Logger
	config                    ChainReaderConfig
	starter                   services.StateMachine
	packageAddresses          map[string]string
	client                    *client.PTBClient
	dbStore                   *database.DBStore
	eventsIndexer             indexer.EventsIndexerApi
	eventsIndexerCancel       *context.CancelFunc
	transactionsIndexer       indexer.TransactionsIndexerApi
	transactionsIndexerCancel *context.CancelFunc
}

var _ pkgtypes.ContractTypeProvider = &suiChainReader{}

type ExtendedContractReader interface {
	pkgtypes.ContractReader
	QueryKeyWithMetadata(ctx context.Context, contract pkgtypes.BoundContract, filter query.KeyFilter, limitAndSort query.LimitAndSort, sequenceDataType any) ([]SequenceWithMetadata, error)
}

// readIdentifier represents the parsed components of a read identifier
type readIdentifier struct {
	address      string
	contractName string
	readName     string
}

func NewChainReader(ctx context.Context, lgr logger.Logger, abstractClient *client.PTBClient, config ChainReaderConfig, db sqlutil.DataSource) (pkgtypes.ContractReader, error) {
	dbStore := database.NewDBStore(db, lgr)

	err := dbStore.EnsureSchema(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to ensure database schema: %w", err)
	}

	// Create a list of all event selectors to pass to indexer
	eventConfigurations := make([]*client.EventSelector, 0)
	for _, moduleConfig := range config.Modules {
		if moduleConfig.Events != nil {
			for _, eventConfig := range moduleConfig.Events {
				eventConfigurations = append(eventConfigurations, &eventConfig.EventSelector)
			}
		}
	}

	eventsIndexer := indexer.NewEventIndexer(
		dbStore,
		lgr,
		abstractClient,
		eventConfigurations,
		config.EventsIndexer.PollingInterval,
		config.EventsIndexer.SyncTimeout,
	)

	transactionsIndexer := indexer.NewTransactionsIndexer(
		dbStore,
		config.TransactionsIndexer.PollingInterval,
		config.TransactionsIndexer.SyncTimeout,
	)

	return &suiChainReader{
		logger:           logger.Named(lgr, "SuiChainReader"),
		client:           abstractClient,
		config:           config,
		dbStore:          dbStore,
		packageAddresses: map[string]string{},
		// indexers
		eventsIndexer:             eventsIndexer,
		transactionsIndexer:       transactionsIndexer,
		eventsIndexerCancel:       nil,
		transactionsIndexerCancel: nil,
	}, nil
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
		// start events indexer
		eventsIndexerCtx, cancelEventsIndexerCtx := context.WithCancel(ctx)
		go func() {
			err := s.eventsIndexer.Start(eventsIndexerCtx)
			if err != nil {
				s.logger.Error("Indexer failed to start", "error", err)
				if s.eventsIndexerCancel != nil {
					(*s.eventsIndexerCancel)()
				}
			}
			s.logger.Info("Events indexer started")
			// set the cancel function
			s.eventsIndexerCancel = &cancelEventsIndexerCtx
		}()

		// start transactions indexer
		transactionsIndexerCtx, cancelTransactionsIndexerCtx := context.WithCancel(ctx)
		go func() {
			err := s.transactionsIndexer.Start(transactionsIndexerCtx)
			if err != nil {
				s.logger.Error("Indexer failed to start", "error", err)
				if s.transactionsIndexerCancel != nil {
					(*s.transactionsIndexerCancel)()
				}
			}
			s.logger.Info("Transactions indexer started")
			// set the cancel function
			s.transactionsIndexerCancel = &cancelTransactionsIndexerCtx
		}()

		return nil
	})
}

func (s *suiChainReader) Close() error {
	return s.starter.StopOnce(s.Name(), func() error {
		// stop events indexer
		if s.eventsIndexerCancel != nil {
			(*s.eventsIndexerCancel)()
		}
		s.logger.Info("Events indexer stopped")

		// stop transactions indexer
		if s.transactionsIndexerCancel != nil {
			(*s.transactionsIndexerCancel)()
		}
		s.logger.Info("Transactions indexer stopped")

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

	maps.Copy(s.packageAddresses, newBindings)

	// Update the indexer's package addresses and event configurations
	// This ensures the indexer knows about the newly bound contracts
	s.updateIndexerConfiguration()

	return nil
}

// updateIndexerConfiguration updates the indexer with current bindings and configurations
func (s *suiChainReader) updateIndexerConfiguration() {
	// Create event configurations for the indexer based on the chainreader config
	// TODO: Update the indexer's configuration dynamically
	// For now, the indexer will be created with empty configurations
	// and will need to be recreated when bindings change
	s.logger.Warnw("Updated indexer configuration")
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

// QueryKey queries events from the indexer database for events that were populated from the RPC node
func (s *suiChainReader) QueryKey(ctx context.Context, contract pkgtypes.BoundContract, filter query.KeyFilter, limitAndSort query.LimitAndSort, sequenceDataType any) ([]pkgtypes.Sequence, error) {
	// Validate contract binding
	if err := s.validateContractBinding(contract); err != nil {
		return nil, err
	}

	// Get module and event configuration
	moduleConfig := s.config.Modules[contract.Name]
	eventConfig, err := s.getEventConfig(moduleConfig, filter.Key)
	// No event config found, construct a config
	if err != nil {
		// construct a new config ad-hoc
		eventConfig = &ChainReaderEvent{
			Name:      filter.Key,
			EventType: filter.Key,
			EventSelector: client.EventSelector{
				Package: contract.Address,
				Module:  contract.Name,
				Event:   filter.Key,
			},
		}
	}

	// Sync the event in case it's not already in the database
	err = s.eventsIndexer.SyncEvent(ctx, &eventConfig.EventSelector)
	if err != nil {
		return nil, err
	}

	// Query events from database
	eventRecords, err := s.queryEvents(ctx, contract, eventConfig, limitAndSort)
	if err != nil {
		return nil, err
	}

	// Transform events to sequences
	return s.transformEventsToSequences(eventRecords, sequenceDataType)
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

func (s *suiChainReader) CreateContractType(readName string, forEncoding bool) (any, error) {
	// only called when LOOP plugin
	// TODO: should something be added to the LOOP plugin?
	return &[]byte{}, nil
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

// isObjectRead determines if the read operation is for an object (starts with 0x) or a function call
func isObjectRead(readName string) bool {
	return strings.HasPrefix(readName, objectIdPrefix)
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

// queryEvents queries events from the database instead of the Sui blockchain
func (s *suiChainReader) queryEvents(ctx context.Context, contract pkgtypes.BoundContract, eventConfig *ChainReaderEvent, limitAndSort query.LimitAndSort) ([]database.EventRecord, error) {
	// Create the event handle for database lookup
	eventHandle := fmt.Sprintf("%s::%s::%s", contract.Address, contract.Name, eventConfig.EventType)

	s.logger.Debugw("Querying events from database",
		"address", contract.Address,
		"module", contract.Name,
		"eventType", eventConfig.EventType,
		"eventHandle", eventHandle,
		"limit", limitAndSort.Limit.Count,
	)

	// Query events from database
	records, err := s.dbStore.QueryEvents(ctx, contract.Address, eventHandle, []query.Expression{}, limitAndSort)
	if err != nil {
		s.logger.Errorw("Failed to query events from database",
			"error", err,
			"address", contract.Address,
			"module", contract.Name,
			"eventType", eventConfig.EventType,
			"eventHandle", eventHandle,
		)

		return nil, fmt.Errorf("failed to query events from database: %w", err)
	}

	s.logger.Debugw("Successfully queried events from database",
		"eventCount", len(records),
		"eventHandle", eventHandle,
	)

	return records, nil
}

// transformEventsToSequences converts database event records to sequence format
func (s *suiChainReader) transformEventsToSequences(eventRecords []database.EventRecord, sequenceDataType any) ([]pkgtypes.Sequence, error) {
	sequences := make([]pkgtypes.Sequence, 0, len(eventRecords))

	for _, record := range eventRecords {
		eventData := reflect.New(reflect.TypeOf(sequenceDataType).Elem()).Interface()

		s.logger.Debugw("Processing database event record", "data", record.Data, "offset", record.EventOffset)

		// if we are running in loop plugin mode, we will want to decode into JSON and then into JSON bytes always
		if s.config.IsLoopPlugin {
			// decode into JSON and then into JSON bytes
			jsonData, err := json.Marshal(record.Data)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal data for LOOP: %w", err)
			}
			eventData = &jsonData
		} else if err := codec.DecodeSuiJsonValue(record.Data, eventData); err != nil {
			return nil, fmt.Errorf("failed to decode event data: %w", err)
		}

		// Create cursor from the event offset - this is simpler than the blockchain event ID
		// TODO: change this to match what's expected in DB lookups
		cursor := fmt.Sprintf(`{"event_offset": %d}`, record.EventOffset)

		sequence := pkgtypes.Sequence{
			Cursor: cursor,
			Data:   eventData,
			Head: pkgtypes.Head{
				Timestamp: record.BlockTimestamp,
				Hash:      record.BlockHash,
				Height:    record.BlockHeight,
			},
		}

		sequences = append(sequences, sequence)
	}

	s.logger.Debugw("Successfully transformed events to sequences", "sequenceCount", len(sequences))

	return sequences, nil
}
