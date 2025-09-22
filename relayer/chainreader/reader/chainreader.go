package reader

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/mitchellh/mapstructure"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"

	"github.com/smartcontractkit/chainlink-aptos/relayer/chainreader/loop"

	aptosCRConfig "github.com/smartcontractkit/chainlink-aptos/relayer/chainreader/config"
	aptosCRUtils "github.com/smartcontractkit/chainlink-aptos/relayer/chainreader/utils"

	crUtil "github.com/smartcontractkit/chainlink-sui/relayer/chainreader/chainreader_util"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/config"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/database"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/indexer"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/codec"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	pkgtypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
)

const (
	defaultQueryLimit   = 25
	readIdentifierParts = 3
	offrampName         = "OffRamp"
	ccipPointerKey      = "state_object::CCIPObjectRefPointer"
)

type suiChainReader struct {
	pkgtypes.UnimplementedContractReader

	logger          logger.Logger
	config          config.ChainReaderConfig
	starter         services.StateMachine
	packageResolver *crUtil.PackageResolver
	client          *client.PTBClient
	dbStore         *database.DBStore
	indexer         indexer.IndexerApi
}

var _ pkgtypes.ContractTypeProvider = &suiChainReader{}

type ExtendedContractReader interface {
	pkgtypes.ContractReader
	QueryKeyWithMetadata(ctx context.Context, contract pkgtypes.BoundContract, filter query.KeyFilter, limitAndSort query.LimitAndSort, sequenceDataType any) ([]aptosCRConfig.SequenceWithMetadata, error)
}

// readIdentifier represents the parsed components of a read identifier
type readIdentifier struct {
	address      string
	contractName string
	readName     string
}

func NewChainReader(
	ctx context.Context,
	lgr logger.Logger,
	ptbClient *client.PTBClient,
	configs config.ChainReaderConfig,
	db sqlutil.DataSource,
	indexer indexer.IndexerApi,
) (pkgtypes.ContractReader, error) {
	dbStore := database.NewDBStore(db, lgr)

	err := dbStore.EnsureSchema(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to ensure database schema: %w", err)
	}

	return &suiChainReader{
		logger:          logger.Named(lgr, "SuiChainReader"),
		client:          ptbClient,
		config:          configs,
		dbStore:         dbStore,
		packageResolver: crUtil.NewPackageResolver(lgr, ptbClient),
		// indexers
		indexer: indexer,
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
		return nil
	})
}

func (s *suiChainReader) Close() error {
	return s.starter.StopOnce(s.Name(), func() error {
		return nil
	})
}

func (s *suiChainReader) Bind(ctx context.Context, bindings []pkgtypes.BoundContract) error {
	for _, binding := range bindings {
		err := s.packageResolver.BindPackage(binding.Name, binding.Address)
		if err != nil {
			return fmt.Errorf("failed to bind package: %w", err)
		}
	}

	// If the "OffRamp" package/module is now bound, set the offramp package ID for the tx indexer
	if pkg, err := s.packageResolver.ResolvePackageAddress(offrampName); err == nil {
		s.indexer.GetTransactionIndexer().SetOffRampPackage(pkg)
	}
	return nil
}

func (s *suiChainReader) Unbind(ctx context.Context, bindings []pkgtypes.BoundContract) error {
	for _, binding := range bindings {
		if err := s.packageResolver.UnbindPackage(binding.Name); err != nil {
			return fmt.Errorf("failed to unbind package %s: %w", binding.Name, err)
		}
	}

	return nil
}

// GetLatestValue retrieves the latest value from either an object or function call
func (s *suiChainReader) GetLatestValue(ctx context.Context, readIdentifier string, _ primitives.ConfidenceLevel, params, returnVal any) error {
	parsed, err := s.parseReadIdentifier(readIdentifier)
	if err != nil {
		return err
	}
	_, contractName, method := parsed.address, parsed.contractName, parsed.readName

	if err = s.validateContractBindingAndConfig(parsed.contractName, parsed.address); err != nil {
		return fmt.Errorf("error validating contract binding and config: %w", err)
	}

	// this ensures we are using values from chain-reader config set in core
	moduleConfig, ok := s.config.Modules[contractName]
	if !ok {
		return fmt.Errorf("no such contract: %s", contractName)
	}

	if moduleConfig.Functions == nil {
		return fmt.Errorf("no functions for contract: %s", contractName)
	}

	functionConfig, ok := moduleConfig.Functions[method]
	if !ok {
		return fmt.Errorf("no such method: %s", method)
	}

	if moduleConfig.Name != "" {
		parsed.contractName = moduleConfig.Name
	}

	if functionConfig.Name != "" {
		parsed.readName = functionConfig.Name
	}

	s.logger.Debugw("calling function after overwrite",
		"address", parsed.address,
		"contract", parsed.contractName,
		"function", parsed.readName,
	)

	results, err := s.callFunction(ctx, parsed, params, functionConfig)
	if err != nil {
		return err
	}

	if functionConfig.ResultTupleToStruct != nil {
		structResult := make(map[string]any)
		for i, mapKey := range functionConfig.ResultTupleToStruct {
			structResult[mapKey] = results[i]
		}

		// Apply result field renames if configured
		if functionConfig.ResultFieldRenames != nil {
			err = aptosCRUtils.MaybeRenameFields(structResult, functionConfig.ResultFieldRenames)
			if err != nil {
				return fmt.Errorf("failed to rename result fields in GetLatestValue: %w", err)
			}
		}

		// if we are running in loop plugin mode, we will want to encode the result into JSON bytes
		if s.config.IsLoopPlugin {
			return s.encodeLoopResult(structResult, returnVal)
		}

		return codec.DecodeSuiJsonValue(structResult, returnVal)
	}

	// otherwise, no tuple to struct specification, just a slice of values
	if s.config.IsLoopPlugin {
		// Apply renames to the result slice or contained maps before encoding
		var renamed any = results
		if functionConfig.ResultFieldRenames != nil {
			err = aptosCRUtils.MaybeRenameFields(renamed, functionConfig.ResultFieldRenames)
			if err != nil {
				return fmt.Errorf("failed to rename result fields in GetLatestValue: %w", err)
			}
		}
		return s.encodeLoopResult(renamed, returnVal)
	}

	s.logger.Debugw("GLV results before decoding to SUI json", "results", results, "returnVal", returnVal)

	// Apply renames (if any) to the primary result element before decoding
	var primary any = results[0]
	if functionConfig.ResultFieldRenames != nil {
		err = aptosCRUtils.MaybeRenameFields(primary, functionConfig.ResultFieldRenames)
		if err != nil {
			return fmt.Errorf("failed to rename result fields in GetLatestValue: %w", err)
		}
	}

	// TODO: handle multiple results for non-loop plugin mode
	return codec.DecodeSuiJsonValue(primary, returnVal)
}

// QueryKey queries events from the indexer database for events that were populated from the RPC node
func (s *suiChainReader) QueryKey(ctx context.Context, contract pkgtypes.BoundContract, filter query.KeyFilter, limitAndSort query.LimitAndSort, sequenceDataType any) ([]pkgtypes.Sequence, error) {
	eventConfig, err := s.updateEventConfigs(ctx, contract, filter)
	if err != nil {
		return nil, err
	}

	// Query events from database
	eventRecords, err := s.queryEvents(ctx, eventConfig, filter.Expressions, limitAndSort)
	if err != nil {
		return nil, err
	}

	// Transform events to sequences
	sequences, err := s.transformEventsToSequences(eventRecords, sequenceDataType, false)
	if err != nil {
		return nil, err
	}

	transformedSequences := make([]pkgtypes.Sequence, 0)
	for _, seq := range sequences {
		transformedSequences = append(transformedSequences, seq.Sequence)
	}

	return transformedSequences, nil
}

type cursor struct {
	EventOffset int64 `json:"event_offset"`
}

func (s *suiChainReader) QueryKeyWithMetadata(ctx context.Context, contract pkgtypes.BoundContract, filter query.KeyFilter, limitAndSort query.LimitAndSort, sequenceDataType any) ([]aptosCRConfig.SequenceWithMetadata, error) {
	eventConfig, err := s.updateEventConfigs(ctx, contract, filter)
	if err != nil {
		return nil, err
	}

	// Query events from database
	eventRecords, err := s.queryEvents(ctx, eventConfig, filter.Expressions, limitAndSort)
	if err != nil {
		return nil, err
	}

	// Transform events to sequences
	sequences, err := s.transformEventsToSequences(eventRecords, sequenceDataType, true)
	if err != nil {
		return nil, err
	}

	// Transform events to enriched sequences (include metadata)
	transformedSequences := make([]aptosCRConfig.SequenceWithMetadata, 0)
	for _, seq := range sequences {
		var c cursor
		if err := json.Unmarshal([]byte(seq.Sequence.Cursor), &c); err != nil {
			return nil, fmt.Errorf("failed to unmarshal cursor: %w", err)
		}

		seq.Sequence.Cursor = strconv.FormatInt(c.EventOffset, 10)
		transformedSequences = append(transformedSequences, aptosCRConfig.SequenceWithMetadata{
			Sequence:  seq.Sequence,
			TxVersion: 0,
			TxHash:    seq.Record.TxDigest,
		})
	}

	return transformedSequences, nil
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

func (s *suiChainReader) updateEventConfigs(ctx context.Context, contract pkgtypes.BoundContract, filter query.KeyFilter) (*config.ChainReaderEvent, error) {
	// Validate contract binding
	if err := s.validateContractBindingAndConfig(contract.Name, contract.Address); err != nil {
		return nil, err
	}

	// Get module and event configuration
	moduleConfig := s.config.Modules[contract.Name]
	eventConfig, err := s.getEventConfig(moduleConfig, filter.Key)
	// No event config found, construct a config
	if err == nil && eventConfig == nil {
		// construct a new config ad-hoc
		eventConfig = &config.ChainReaderEvent{
			Name:      filter.Key,
			EventType: filter.Key,
			EventSelector: client.EventSelector{
				Package: contract.Address,
				Module:  contract.Name,
				Event:   filter.Key,
			},
		}
	} else if err != nil {
		return nil, err
	}

	if moduleConfig.Name != "" {
		eventConfig.Name = moduleConfig.Name
	} else {
		// If the module config has no name, use the module name from the event config
		moduleConfig.Name = moduleConfig.Events[filter.Key].Module
	}

	// only write contract address, rest will be handled during chainreader config
	eventConfig.Package = contract.Address

	// repeat the sync call for each package ID (upgrades) of the module
	// using the contract's own address as signer address since we are only ready
	packageIds, err := s.client.LoadModulePackageIds(ctx, contract.Address, moduleConfig.Name, contract.Address)
	if err != nil {
		return nil, err
	}

	s.logger.Debugw("Found package IDs", "packageIds", packageIds)

	evIndexer := s.indexer.GetEventIndexer()
	// create a selector for each package ID including the upgrades and the initial package ID
	// the `LoadModulePackageIds` will fallback to a single package ID if the module does not have the `get_package_ids` function
	for _, packageId := range packageIds {
		selector := client.EventSelector{
			Package: packageId,
			Module:  moduleConfig.Name,
			Event:   eventConfig.EventType,
			// override the DB insert using the initial package ID
			InitialPackageId: &contract.Address,
		}

		// sync the event in case it's not already in the database
		err = evIndexer.SyncEvent(ctx, &selector)
		if err != nil {
			return nil, err
		}
	}

	// update the event config in the transactions indexer to ensure that the package ID is known
	s.indexer.GetTransactionIndexer().UpdateEventConfig(eventConfig)

	return eventConfig, nil
}

// validateContractBinding validates the contract binding for QueryKey
func (s *suiChainReader) validateContractBindingAndConfig(name string, address string) error {
	err := s.packageResolver.ValidateBinding(name, address)
	if err != nil {
		return fmt.Errorf("invalid binding for contract: %s. err: %w", name, err)
	}

	if _, ok := s.config.Modules[name]; !ok {
		return fmt.Errorf("no configuration for contract: %s", name)
	}

	return nil
}

// callFunction calls a contract function and returns the result
func (s *suiChainReader) callFunction(ctx context.Context, parsed *readIdentifier, params any, functionConfig *config.ChainReaderFunction) ([]any, error) {
	fmt.Println("CALLING CALL FUNCTION")
	argMap, err := s.parseParams(params, functionConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to parse parameters: %w", err)
	}
	fmt.Println("PASSED PARSE PARAMS")
	args, argTypes, err := s.prepareArguments(ctx, argMap, functionConfig, parsed)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare arguments: %w", err)
	}
	fmt.Println("PASSED PREPARE ARGS")
	responseValues, err := s.executeFunction(ctx, parsed, functionConfig, args, argTypes)
	if err != nil {
		return nil, err
	}

	return responseValues, nil
}

// parseParams parses input parameters based on whether we're running as a LOOP plugin
func (s *suiChainReader) parseParams(params any, functionConfig *config.ChainReaderFunction) (map[string]any, error) {
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
func (s *suiChainReader) parseLoopParams(params any, functionConfig *config.ChainReaderFunction) (map[string]any, error) {
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

type pointerMapEntry struct {
	field     string // the field name from the Sui object
	paramName string // the parameter name from the function config
}

// prepareArguments prepares function arguments and types for the call
func (s *suiChainReader) prepareArguments(ctx context.Context, argMap map[string]any, functionConfig *config.ChainReaderFunction, identifier *readIdentifier) ([]any, []string, error) {
	if functionConfig.Params == nil {
		return []any{}, []string{}, nil
	}

	// referring to the tag parts "_::module::Pointer::field"
	tagLength := 4

	// a map of object selector "module::object" to array of fields
	pointersMap := make(map[string][]pointerMapEntry)
	pointerSelectors := make(map[string]readIdentifier)

	// make a set of object pointers that need to fetched
	// to read more about pointer tags, see the documentation in "/relayer/documentation/relayer/pointer-tags-in-cr.md"
	for _, paramConfig := range functionConfig.Params {
		// the parameter has a pointer tag, add it to the set
		if paramConfig.PointerTag != nil {
			tag := strings.Split(*paramConfig.PointerTag, "::")
			// must be 4 values, for example: "_::moduleName::pointerName::fieldName"
			if len(tag) != tagLength {
				return nil, nil, fmt.Errorf("invalid pointer tag: %s", *paramConfig.PointerTag)
			}

			moduleName, pointerName, fieldName := tag[1], tag[2], tag[3]

			// append only the middle 2 parts of the tag to represent the pointer
			appendTag := strings.Join([]string{moduleName, pointerName}, "::")
			if _, ok := pointersMap[appendTag]; !ok {
				pointersMap[appendTag] = make([]pointerMapEntry, 0)
			}
			// add the pointer selector to the map which will later be used to fetch the values from the package owned object fields
			if _, ok := pointerSelectors[appendTag]; !ok {
				readIdentifierForPointer := readIdentifier{
					address:      identifier.address,
					contractName: moduleName,
					readName:     pointerName,
				}

				// special case for pointers from the CCIP package object pointer
				// this is needed to override the specified address (will be offramp package ID) with the CCIP package ID
				if identifier.contractName == strings.ToLower(offrampName) && appendTag == ccipPointerKey {
					ccipPackageID, err := s.client.GetCCIPPackageID(ctx, identifier.address, functionConfig.SignerAddress)
					if err != nil {
						return nil, nil, fmt.Errorf("failed to get CCIP package ID: %w", err)
					}
					readIdentifierForPointer.address = ccipPackageID
				}

				pointerSelectors[appendTag] = readIdentifierForPointer
			}

			// each entry within the pointersMap contains an entry for the field name and
			// an entry for the (function config) parameter name
			pointersMap[appendTag] = append(pointersMap[appendTag], pointerMapEntry{
				field:     fieldName,
				paramName: paramConfig.Name,
			})
		}
	}

	// fetch pointers
	for pointerTag, pointerVals := range pointersMap {
		fields := make([]string, 0, len(pointerVals))

		// get the fields from the pointer values
		for _, pointerVal := range pointerVals {
			fields = append(fields, pointerVal.field)
		}

		selector := pointerSelectors[pointerTag]
		pointerFieldValues, err := s.client.GetValuesFromPackageOwnedObjectField(
			ctx, selector.address, selector.contractName, selector.readName, fields,
		)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get values from package owned object fields: %w", err)
		}

		// add the values to the arg map
		for _, pointerVal := range pointerVals {
			argMap[pointerVal.paramName] = pointerFieldValues[pointerVal.field]
		}
	}

	args := make([]any, 0, len(functionConfig.Params))
	argTypes := make([]string, 0, len(functionConfig.Params))

	// ensure that all the required arguments are present
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
func (s *suiChainReader) executeFunction(ctx context.Context, parsed *readIdentifier, functionConfig *config.ChainReaderFunction, args []any, argTypes []string) ([]any, error) {
	s.logger.Debugw("Calling ReadFunction",
		"address", parsed.address,
		"module", parsed.contractName,
		"method", parsed.readName,
		"encodedArgs", args,
		"argTypes", argTypes,
	)

	values, err := s.client.ReadFunction(ctx, functionConfig.SignerAddress, parsed.address, parsed.contractName, parsed.readName, args, argTypes)
	if err != nil {
		s.logger.Errorw("ReadFunction failed",
			"error", err,
			"address", parsed.address,
			"module", parsed.contractName,
			"method", parsed.readName,
			"args", args,
			"argTypes", argTypes,
		)

		return nil, fmt.Errorf("failed to call function %s: %w", parsed.readName, err)
	}

	s.logger.Debugw("Sui ReadFunction response", "returnValues", values)

	return values, nil
}

// encodeLoopResult encodes results for LOOP plugin mode
func (s *suiChainReader) encodeLoopResult(valueField any, returnVal any) error {
	var toMarshal any

	// Check if the value is a map
	if resultMap, mapOk := valueField.(map[string]any); mapOk {
		toMarshal = resultMap
	} else if resultSlice, sliceOk := valueField.([]any); sliceOk {
		// For primitive values like uint64, the data might not be in a map structure
		if len(resultSlice) == 1 {
			// if it's a single value, we can just marshal it
			toMarshal = resultSlice[0]
		} else {
			// if it's a slice of values, we need to marshal the whole slice
			toMarshal = resultSlice
		}
	} else {
		return fmt.Errorf("expected valueField to be map[string]any or []any, got %T", valueField)
	}

	resultBytes, err := json.Marshal(toMarshal)
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
func (s *suiChainReader) getEventConfig(moduleConfig *config.ChainReaderModule, eventKey string) (*config.ChainReaderEvent, error) {
	if moduleConfig.Events == nil {
		return nil, fmt.Errorf("no events configured for contract")
	}

	eventConfig, ok := moduleConfig.Events[eventKey]
	if !ok {
		s.logger.Errorw("No configuration for event", "eventKey", eventKey, "moduleConfig", moduleConfig)
		return nil, fmt.Errorf("no configuration for event: %s", eventKey)
	}

	return eventConfig, nil
}

// queryEvents queries events from the database instead of the Sui blockchain
func (s *suiChainReader) queryEvents(ctx context.Context, eventConfig *config.ChainReaderEvent, expressions []query.Expression, limitAndSort query.LimitAndSort) ([]database.EventRecord, error) {
	// Create the event handle for database lookup
	eventHandle := fmt.Sprintf("%s::%s::%s", eventConfig.Package, eventConfig.Name, eventConfig.EventType)

	s.logger.Debugw("Querying events from database",
		"address", eventConfig.Package,
		"module", eventConfig.Name,
		"eventType", eventConfig.EventType,
		"eventHandle", eventHandle,
		"limit", limitAndSort.Limit.Count,
	)

	if s.config.IsLoopPlugin {
		deserializedExpressions, err := loop.DeserializeExpressions(expressions)
		if err != nil {
			return nil, fmt.Errorf("failed to deserialize expressions: %w", err)
		}
		expressions = deserializedExpressions
	}

	if eventConfig.EventFilterRenames != nil {
		expressions = aptosCRUtils.ApplyEventFilterRenames(expressions, eventConfig.EventFilterRenames)
	}

	// Query events from database
	records, err := s.dbStore.QueryEvents(ctx, eventConfig.Package, eventHandle, expressions, limitAndSort)
	if err != nil {
		s.logger.Errorw("Failed to query events from database",
			"error", err,
			"address", eventConfig.Package,
			"module", eventConfig.Name,
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

type SequenceWithRecord struct {
	Sequence pkgtypes.Sequence
	Record   *database.EventRecord
}

// transformEventsToSequences converts database event records to sequence format
func (s *suiChainReader) transformEventsToSequences(eventRecords []database.EventRecord, sequenceDataType any, includeRecord bool) ([]SequenceWithRecord, error) {
	sequences := make([]SequenceWithRecord, 0, len(eventRecords))

	s.logger.Debugw("Transforming events to sequences", "eventRecords", eventRecords, "sequenceDataType", sequenceDataType)

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

		// If we are simply querying the keys without metadata (non enriched), then we don't need the
		// the original DB record
		if !includeRecord {
			sequences = append(sequences, SequenceWithRecord{
				Sequence: sequence,
				Record:   nil,
			})
			continue
		}

		sequences = append(sequences, SequenceWithRecord{
			Sequence: sequence,
			Record:   &record,
		})
	}

	s.logger.Debugw("Successfully transformed events to sequences", "sequenceCount", len(sequences))

	return sequences, nil
}
