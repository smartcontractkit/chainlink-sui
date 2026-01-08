package reader

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/mitchellh/mapstructure"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"

	"github.com/smartcontractkit/chainlink-aptos/relayer/chainreader/loop"

	aptosCRConfig "github.com/smartcontractkit/chainlink-aptos/relayer/chainreader/config"
	aptosCRUtils "github.com/smartcontractkit/chainlink-aptos/relayer/chainreader/utils"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	crUtil "github.com/smartcontractkit/chainlink-sui/relayer/chainreader/chainreader_util"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/config"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/database"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/indexer"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
	"github.com/smartcontractkit/chainlink-sui/relayer/common"

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

	// Cache of parent object IDs for pointer objects
	// Key format: "{packageID}::{module}::{pointerName}"
	// Value: parent object ID
	parentObjectIDs      map[string]string
	parentObjectIDsMutex sync.RWMutex
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
		indexer:         indexer,
		parentObjectIDs: make(map[string]string),
	}, nil
}

func (s *suiChainReader) Name() string {
	return s.logger.Name()
}

func (s *suiChainReader) Ready() error {
	if err := s.starter.Ready(); err != nil {
		return err
	}

	return nil
}

func (s *suiChainReader) HealthReport() map[string]error {
	report := map[string]error{s.Name(): s.starter.Healthy()}

	// Include indexer health status
	if s.indexer != nil {
		for k, v := range s.indexer.HealthReport() {
			report[k] = v
		}
	}

	return report
}

func (s *suiChainReader) Start(ctx context.Context) error {
	return s.starter.StartOnce(s.Name(), func() error {
		// set the event offset overrides for the event indexer if any
		offsetOverrides := make(map[string]client.EventId)

		for _, moduleConfig := range s.config.Modules {
			for _, eventConfig := range moduleConfig.Events {
				if eventConfig.EventSelectorDefaultOffset != nil {
					key := fmt.Sprintf("%s::%s", eventConfig.EventSelector.Module, eventConfig.EventSelector.Event)
					offsetOverrides[key] = *eventConfig.EventSelectorDefaultOffset
				}
			}
		}

		if len(offsetOverrides) > 0 {
			// ignore this error to avoid blocking the start of the chain reader
			_ = s.indexer.GetEventIndexer().SetEventOffsetOverrides(ctx, offsetOverrides)
		}

		return nil
	})
}

func (s *suiChainReader) Close() error {
	return s.starter.StopOnce(s.Name(), func() error {
		return nil
	})
}

func (s *suiChainReader) Bind(ctx context.Context, bindings []pkgtypes.BoundContract) error {
	offrampPackageAddress := ""

	for _, binding := range bindings {
		err := s.packageResolver.BindPackage(binding.Name, binding.Address)
		if err != nil {
			return fmt.Errorf("failed to bind package: %w", err)
		}

		// Pre-load parent object IDs for known pointer types
		if err := s.preloadParentObjectIDs(ctx, binding); err != nil {
			s.logger.Warnw("Failed to pre-load parent object IDs", "contract", binding.Name, "error", err)
		}

		if common.NormalizeName(binding.Name) == common.NormalizeName(offrampName) {
			offrampPackageAddress = binding.Address
		}
	}

	// If the "OffRamp" package/module is now bound, set the offramp package ID for the tx indexer
	if offrampPackageAddress != "" {
		// Get the latest package ID for the offramp module. The transaction indexer watches the latest package ID for the offramp module
		// to capture failed transactions.
		latestPackageID, err := s.client.GetLatestPackageId(ctx, offrampPackageAddress, "offramp")
		if err != nil {
			s.logger.Errorw("Failed to get latest package ID for OffRamp", "error", err)
			// Use the currently available package address for the offramp module as a fallback
			s.indexer.GetTransactionIndexer().SetOffRampPackage(offrampPackageAddress, offrampPackageAddress)
		} else {
			s.indexer.GetTransactionIndexer().SetOffRampPackage(offrampPackageAddress, latestPackageID)
		}

		// Ensure that the event indexer has the SourceChainConfigSet event selector, this is needed to make sure
		// transaction indexer can find the relevant data when inserting synthetic events.
		err = s.indexer.GetEventIndexer().AddEventSelector(ctx, &client.EventSelector{
			Package: offrampPackageAddress,
			Module:  "offramp",
			Event:   "SourceChainConfigSet",
		})
		if err != nil {
			s.logger.Errorw("Failed to update offramp::SourceChainConfigSet event selector when binding OffRamp", "error", err)
		}
		err = s.indexer.GetEventIndexer().AddEventSelector(ctx, &client.EventSelector{
			Package: offrampPackageAddress,
			Module:  "ocr3_base",
			Event:   "ConfigSet",
		})
		if err != nil {
			s.logger.Errorw("Failed to update ocr3_base::ConfigSet event selector when binding OffRamp", "error", err)
		}
	}

	return nil
}

func (s *suiChainReader) Unbind(ctx context.Context, bindings []pkgtypes.BoundContract) error {
	for _, binding := range bindings {
		if err := s.packageResolver.UnbindPackage(binding.Name); err != nil {
			return fmt.Errorf("failed to unbind package %s: %w", binding.Name, err)
		}

		modulePrefix := fmt.Sprintf("%s::%s::", binding.Address, binding.Name)

		// Clear cached parent object IDs for this unbound contract
		s.parentObjectIDsMutex.Lock()
		for key := range s.parentObjectIDs {
			if strings.HasPrefix(key, modulePrefix) {
				delete(s.parentObjectIDs, key)
			}
		}
		s.parentObjectIDsMutex.Unlock()
	}

	return nil
}

// preloadParentObjectIDs pre-loads parent object IDs for known pointer types at binding time.
// This reduces RPC calls during GetLatestValue by caching parent IDs upfront.
func (s *suiChainReader) preloadParentObjectIDs(ctx context.Context, binding pkgtypes.BoundContract) error {
	pointers := common.GetPointerConfigsByContract(binding.Name)
	if len(pointers) == 0 {
		return nil
	}

	// For OffRamp, we also need to load CCIP pointers
	// If offramp is being bound
	// Fetch the offramp's parent object ID
	// From Offramp, we fetch the CCIP package ID
	// Fetch CCIP's parent object ID
	//
	// We need Offramp state pointer too for offramp, not just CCIP
	var ccipPackageID string
	if strings.EqualFold(binding.Name, offrampName) {
		var err error
		ccipPackageID, err = s.client.GetCCIPPackageID(ctx, binding.Address, binding.Address)
		if err != nil {
			s.logger.Warnw("Failed to get CCIP package ID for OffRamp", "error", err)
		} else {
			// Add CCIP pointers with CCIP package ID
			ccipPointers := common.GetPointerConfigsByContract("ccip")
			pointers = append(pointers, ccipPointers...)
		}
	}

	// Load each pointer's parent object ID
	for _, ptr := range pointers {
		packageID := binding.Address

		// Use CCIP package ID for CCIP pointers when binding OffRamp
		if ptr.Module == "state_object" && ccipPackageID != "" {
			packageID = ccipPackageID
		}

		cacheKey := fmt.Sprintf("%s::%s::%s", packageID, ptr.Module, ptr.Pointer)

		s.parentObjectIDsMutex.RLock()
		_, exists := s.parentObjectIDs[cacheKey]
		s.parentObjectIDsMutex.RUnlock()

		if exists {
			continue
		}

		parentObjectID, err := s.client.GetParentObjectID(ctx, packageID, ptr.Module, ptr.Pointer)
		if err != nil {
			s.logger.Debugw("Could not pre-load parent object ID",
				"packageID", packageID,
				"module", ptr.Module,
				"pointer", ptr.Pointer,
				"error", err)
			continue // Skip this pointer, will load on-demand if needed
		}

		s.parentObjectIDsMutex.Lock()
		s.parentObjectIDs[cacheKey] = parentObjectID
		s.parentObjectIDsMutex.Unlock()

		s.logger.Debugw("Pre-loaded parent object ID",
			"cacheKey", cacheKey,
			"parentObjectID", parentObjectID)
	}

	return nil
}

// GetLatestValue retrieves the latest value from either an object or function call
func (s *suiChainReader) GetLatestValue(ctx context.Context, readIdentifier string, confidenceLevel primitives.ConfidenceLevel, params, returnVal any) error {
	parsed, err := s.parseReadIdentifier(readIdentifier)
	if err != nil {
		return err
	}
	_, contractName, method := parsed.address, parsed.contractName, parsed.readName

	if err = s.validateContractBindingAndConfig(parsed.contractName, parsed.address); err != nil {
		return err
	}

	if params == nil || reflect.ValueOf(params).IsNil() {
		params = make(map[string]any)
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
		// Check the length of results to avoid panics
		if len(results) < len(functionConfig.ResultTupleToStruct) {
			return fmt.Errorf("expected %d results, got %d", len(functionConfig.ResultTupleToStruct), len(results))
		}

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
	responseValues := make([]any, len(results))
	for i, result := range results {
		current := result
		if functionConfig.ResultFieldRenames != nil {
			err = aptosCRUtils.MaybeRenameFields(current, functionConfig.ResultFieldRenames)
			if err != nil {
				return fmt.Errorf("failed to rename result fields in GetLatestValue: %w", err)
			}
		}
		responseValues[i] = current
	}

	if len(results) > 1 {
		return codec.DecodeSuiJsonValue(responseValues, returnVal)
	}

	return codec.DecodeSuiJsonValue(results[0], returnVal)
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
	sequences, err := s.transformEventsToSequences(eventRecords, eventConfig.ExpectedEventType, sequenceDataType, false)
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
	sequences, err := s.transformEventsToSequences(eventRecords, eventConfig.ExpectedEventType, sequenceDataType, true)
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

		var waitgroup sync.WaitGroup
		waitgroup.Add(len(batch))

		for i, read := range batch {
			go func(index int, read pkgtypes.BatchRead, contract pkgtypes.BoundContract) {
				defer waitgroup.Done()
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
			}(i, read, contract)
		}

		// wait for all the results to be processed then close the channel
		go func() {
			waitgroup.Wait()
			close(resultChan)
		}()

		resultsReceived := 0
		for res := range resultChan {
			batchResults[res.index] = res.result
			resultsReceived++
		}

		// check if all the results were received
		if resultsReceived != len(batch) {
			if err := ctx.Err(); err != nil {
				return nil, err
			}
			return nil, fmt.Errorf("batch processing failed: expected %d results, received %d", len(batch), resultsReceived)
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

	if moduleConfig.Name != "" && eventConfig.Name == "" {
		eventConfig.Name = moduleConfig.Name
	} else {
		// If the module config has no name, use the module name from the event config
		moduleConfig.Name = moduleConfig.Events[filter.Key].Module
	}

	if eventConfig.EventSelector.Module == "" {
		eventConfig.EventSelector.Module = moduleConfig.Name
	}

	// only write contract address, rest will be handled during chainreader config
	eventConfig.Package = contract.Address

	evIndexer := s.indexer.GetEventIndexer()
	// create a selector for the initial package ID
	selector := client.EventSelector{
		Package: contract.Address,
		Module:  eventConfig.EventSelector.Module,
		Event:   eventConfig.EventType,
	}

	// ensure that the event selector is included in the indexer's set for upcoming polling loop syncs
	err = evIndexer.AddEventSelector(ctx, &selector)
	if err != nil {
		return nil, fmt.Errorf("failed to add event selector: %w", err)
	}

	return eventConfig, nil
}

// validateContractBinding validates the contract binding for QueryKey
func (s *suiChainReader) validateContractBindingAndConfig(name string, address string) error {
	err := s.packageResolver.ValidateBinding(name, address)
	if err != nil {
		return fmt.Errorf("invalid binding for contract: %s", name)
	}

	if _, ok := s.config.Modules[name]; !ok {
		return fmt.Errorf("no configuration for contract: %s", name)
	}

	return nil
}

// callFunction calls a contract function and returns the result
func (s *suiChainReader) callFunction(ctx context.Context, parsed *readIdentifier, params any, functionConfig *config.ChainReaderFunction) ([]any, error) {
	argMap, err := s.parseParams(params, functionConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to parse parameters: %w", err)
	}

	args, argTypes, err := s.prepareArguments(ctx, argMap, functionConfig, parsed)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare arguments: %w", err)
	}

	// Extract generic type tags from function params
	typeArgs, err := s.extractGenericTypeTags(ctx, parsed, functionConfig, args)
	if err != nil {
		return nil, fmt.Errorf("failed to extract generic type tags: %w", err)
	}

	responseValues, err := s.executeFunction(ctx, parsed, functionConfig, args, argTypes, typeArgs)
	if err != nil {
		return nil, err
	}

	return responseValues, nil
}

// Helper function to extract generic type tags
func (s *suiChainReader) extractGenericTypeTags(ctx context.Context, parsed *readIdentifier, functionConfig *config.ChainReaderFunction, args []any) ([]string, error) {
	if functionConfig.Params == nil {
		return []string{}, nil
	}

	// Use a map to track unique type tags and preserve order
	uniqueTags := make(map[string]struct{})
	keyOrder := make([]string, 0)

	for paramIndex, param := range functionConfig.Params {
		if param.GenericType != nil && *param.GenericType != "" {
			genericType := *param.GenericType
			// Only add if not already present
			if _, exists := uniqueTags[genericType]; !exists {
				keyOrder = append(keyOrder, genericType)
				uniqueTags[genericType] = struct{}{}
			}
		} else if param.GenericDependency != nil && *param.GenericDependency != "" && paramIndex < len(args) {
			genericType, err := s.fetchGenericDependency(ctx, &param, args[paramIndex])
			if err != nil {
				return nil, fmt.Errorf("failed to fetch generic dependency: %w", err)
			}
			if _, exists := uniqueTags[genericType]; !exists {
				keyOrder = append(keyOrder, genericType)
				uniqueTags[genericType] = struct{}{}
			}
		}
	}

	return keyOrder, nil
}

func (s *suiChainReader) fetchGenericDependency(
	ctx context.Context,
	param *codec.SuiFunctionParam,
	paramValue any,
) (string, error) {
	if param == nil || param.GenericDependency == nil || *param.GenericDependency == "" {
		return "", fmt.Errorf("generic dependency is not set")
	}

	switch *param.GenericDependency {
	case "get_token_pool_state_type":
		if paramValue == nil || paramValue.(string) == "" {
			return "", fmt.Errorf("param value is nil or empty string")
		}

		// Use the state object ID to deduce the type
		stateObject, err := s.client.ReadObjectId(ctx, paramValue.(string))
		if err != nil {
			return "", fmt.Errorf("failed to read state object: %w", err)
		}

		s.logger.Debugw("stateObjectType", "stateObjectType", stateObject.Type)

		genericType, err := parseGenericTypeFromObjectType(stateObject.Type)
		if err != nil {
			return "", err
		}

		return genericType, nil
	default:
		return "", fmt.Errorf("unknown generic dependency: %s", *param.GenericDependency)
	}
}

// parseParams parses input parameters based on whether we're running as a LOOP plugin
func (s *suiChainReader) parseParams(params any, functionConfig *config.ChainReaderFunction) (map[string]any, error) {
	argMap := make(map[string]any)

	if params == nil {
		return argMap, nil
	}

	if s.config.IsLoopPlugin {
		return s.parseLoopParams(params, functionConfig)
	}

	if err := mapstructure.Decode(params, &argMap); err != nil {
		return nil, fmt.Errorf("failed to decode parameters: %w", err)
	}

	// Ensure that the argMap is not nil
	if len(argMap) == 0 {
		argMap = make(map[string]any)
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
	derivationKey string // the key used to derive the child object address
	paramName     string // the parameter name from the function config
}

// prepareArguments prepares function arguments and types for the call.
// For pointer tags, it looks up cached parent object IDs (pre-loaded during Bind) and derives
// child object IDs using the derivation keys specified in the pointer tags.
func (s *suiChainReader) prepareArguments(ctx context.Context, argMap map[string]any, functionConfig *config.ChainReaderFunction, identifier *readIdentifier) ([]any, []string, error) {
	if functionConfig.Params == nil {
		return []any{}, []string{}, nil
	}

	// a map of object selector "module::object" to array of fields
	pointersMap := make(map[string][]pointerMapEntry)
	pointerSelectors := make(map[string]readIdentifier)

	// make a set of object pointers that need to fetched
	// to read more about pointer tags, see the documentation in "/relayer/documentation/relayer/pointer-tags-in-cr.md"
	for _, paramConfig := range functionConfig.Params {
		// Skip if no pointer tag configured
		if paramConfig.PointerTag == nil {
			continue
		}

		if err := paramConfig.PointerTag.Validate(); err != nil {
			return nil, nil, fmt.Errorf("invalid pointer tag for parameter %s: %w", paramConfig.Name, err)
		}

		pointerTag := paramConfig.PointerTag

		// append only the middle 2 parts of the tag to represent the pointer
		appendTag := strings.Join([]string{pointerTag.Module, pointerTag.PointerName}, "::")
		if _, ok := pointersMap[appendTag]; !ok {
			pointersMap[appendTag] = make([]pointerMapEntry, 0)
		}
		// add the pointer selector to the map which will later be used to fetch the values from the package owned object fields
		if _, ok := pointerSelectors[appendTag]; !ok {
			readIdentifierForPointer := readIdentifier{
				address:      identifier.address,
				contractName: pointerTag.Module,
				readName:     pointerTag.PointerName,
			}

			// If the pointer tag specifies a PackageID, use it (for cross-package dependencies)
			// This is needed when the pointer object is owned by a different package than the calling contract.
			// e.g. When offramp calls a function that needs CCIPObjectRef from CCIP package,
			// We must search for the pointer in CCIP's owned objects, not offramp's owned objects.
			if pointerTag.PackageID != "" {
				readIdentifierForPointer.address = pointerTag.PackageID
			} else if identifier.contractName == strings.ToLower(offrampName) && appendTag == ccipPointerKey {
				// Special case for OffRamp->CCIP pointer (legacy behavior)
				// This is needed to override the specified address (will be offramp package ID) with the CCIP package ID
				// Only handle offRamp case because other modules are within ccip package
				ccipPackageID, err := s.client.GetCCIPPackageID(ctx, identifier.address, functionConfig.SignerAddress)
				if err != nil {
					return nil, nil, fmt.Errorf("failed to get CCIP package ID: %w", err)
				}
				readIdentifierForPointer.address = ccipPackageID
			}

			pointerSelectors[appendTag] = readIdentifierForPointer
		}

		// each entry within the pointersMap contains the derivation key and the (function config) parameter name
		// the parent field name is looked up from common.PointerConfigs when fetching the parent object ID
		pointersMap[appendTag] = append(pointersMap[appendTag], pointerMapEntry{
			derivationKey: pointerTag.DerivationKey,
			paramName:     paramConfig.Name,
		})
	}

	// fetch pointers
	for pointerTag, pointerVals := range pointersMap {
		selector := pointerSelectors[pointerTag]

		// Try to get parent object ID from cache first
		cacheKey := fmt.Sprintf("%s::%s::%s", selector.address, selector.contractName, selector.readName)

		s.parentObjectIDsMutex.RLock()
		parentObjectID, cached := s.parentObjectIDs[cacheKey]
		s.parentObjectIDsMutex.RUnlock()

		if !cached {
			// Not in cache, fetch from RPC (fallback for on-demand loading)
			var err error
			parentObjectID, err = s.client.GetParentObjectID(
				ctx, selector.address, selector.contractName, selector.readName,
			)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to get parent object ID: %w", err)
			}

			// Cache it for next time
			s.parentObjectIDsMutex.Lock()
			s.parentObjectIDs[cacheKey] = parentObjectID
			s.parentObjectIDsMutex.Unlock()

			s.logger.Debugw("Loaded parent object ID on-demand", "cacheKey", cacheKey, "parentObjectId", parentObjectID)
		}

		// Derive each field's object ID from parent using derivation key and add to arg map
		for _, pointerVal := range pointerVals {
			derivedID, err := bind.DeriveObjectIDWithVectorU8Key(parentObjectID, []byte(pointerVal.derivationKey))
			if err != nil {
				return nil, nil, fmt.Errorf("failed to derive object ID for %s using key %s: %w", pointerVal.paramName, pointerVal.derivationKey, err)
			}
			argMap[pointerVal.paramName] = derivedID
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
func (s *suiChainReader) executeFunction(ctx context.Context, parsed *readIdentifier, functionConfig *config.ChainReaderFunction, args []any, argTypes []string, typeArgs []string) ([]any, error) {
	s.logger.Debugw("Calling ReadFunction",
		"address", parsed.address,
		"module", parsed.contractName,
		"method", parsed.readName,
		"encodedArgs", args,
		"argTypes", argTypes,
		"typeArgs", typeArgs,
	)

	// Override the package ID with the latest package ID of the module being called.
	// This ensure we are always using the latestPkgID in case of upgrades.
	latestPackageId, err := s.client.GetLatestPackageId(ctx, parsed.address, common.GetModuleForContract(parsed.contractName))
	if err != nil {
		return []any{}, err
	}

	// this is the upgraded pkgID
	parsed.address = latestPackageId

	if len(functionConfig.StaticResponse) > 0 {
		return functionConfig.StaticResponse, nil
	} else if len(functionConfig.ResponseFromInputs) > 0 {
		for _, pluckFromInput := range functionConfig.ResponseFromInputs {
			switch pluckFromInput {
			case "package_id":
				return []any{latestPackageId}, nil
			default:
				return nil, fmt.Errorf("unknown response from inputs selection: %s", pluckFromInput)
			}
		}
	}

	values, err := s.client.ReadFunction(ctx, functionConfig.SignerAddress, parsed.address, parsed.contractName, parsed.readName, args, argTypes, typeArgs)
	if err != nil {
		s.logger.Errorw("ReadFunction failed",
			"error", err,
			"address", parsed.address,
			"module", parsed.contractName,
			"method", parsed.readName,
			"args", args,
			"argTypes", argTypes,
			"typeArgs", typeArgs,
		)

		return nil, fmt.Errorf("failed to call function %s: %w", parsed.readName, err)
	}

	s.logger.Debugw("Sui ReadFunction response", "returnValues", values)

	// TODO: Remove this once bindings are used in CR, this is a temporary fix for data ingestion
	hexified := common.ConvertBytesToHex(values).([]any)

	return hexified, nil
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
	eventHandle := fmt.Sprintf("%s::%s::%s", eventConfig.Package, eventConfig.EventSelector.Module, eventConfig.EventType)

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

	s.logger.Debugw("QueryKey received request",
		"contract", eventConfig.Package,
		"eventHandle", eventHandle,
		"expressions", expressions,
		"limitAndSort", limitAndSort)

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

	// Apply the event field renames to the returned records if specified in the config
	if len(eventConfig.EventFieldRenames) > 0 {
		for _, rec := range records {
			mappedData := rec.Data
			renameErr := aptosCRUtils.MaybeRenameFields(mappedData, eventConfig.EventFieldRenames)
			if renameErr != nil {
				s.logger.Errorw("Failed to rename event data fields", "error", renameErr)
				continue
			}

			rec.Data = mappedData
		}
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
func (s *suiChainReader) transformEventsToSequences(eventRecords []database.EventRecord, eventDataType any, sequenceDataType any, includeRecord bool) ([]SequenceWithRecord, error) {
	sequences := make([]SequenceWithRecord, 0, len(eventRecords))

	s.logger.Debugw("Transforming events to sequences", "eventRecords", eventRecords, "sequenceDataType", sequenceDataType)

	expectedEventType := sequenceDataType
	if eventDataType != nil {
		expectedEventType = eventDataType
	}

	for _, record := range eventRecords {
		t := reflect.TypeOf(expectedEventType)
		if t == nil || t.Kind() != reflect.Ptr {
			return nil, fmt.Errorf("sequenceDataType must be a non-nil pointer type")
		}
		eventData := reflect.New(t.Elem()).Interface()

		s.logger.Debugw("Processing database event record", "data", record.Data, "offset", record.EventOffset, "eventDataType", reflect.TypeOf(eventData).Elem())

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

		// Transform the data into the original required type
		if reflect.TypeOf(expectedEventType).Elem() != reflect.TypeOf(sequenceDataType).Elem() {
			transformedData := reflect.New(reflect.TypeOf(sequenceDataType).Elem()).Interface()
			if err := codec.DecodeSuiJsonValue(eventData, transformedData); err != nil {
				return nil, fmt.Errorf("failed to decode event data: %w", err)
			}
			eventData = &transformedData
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

		// create a copy of the record to ensure correct memory location
		toSave := record
		sequences = append(sequences, SequenceWithRecord{
			Sequence: sequence,
			Record:   &toSave,
		})
	}

	s.logger.Debugw("Successfully transformed events to sequences", "sequenceCount", len(sequences), "sequences", sequences)

	return sequences, nil
}

func parseGenericTypeFromObjectType(objectType string) (string, error) {
	startIdx := strings.Index(objectType, "<")
	endIdx := strings.LastIndex(objectType, ">")

	if startIdx == -1 || endIdx == -1 || startIdx >= endIdx {
		return "", fmt.Errorf("invalid object type format, expected generic type parameter: %s", objectType)
	}

	genericType := objectType[startIdx+1 : endIdx]
	return genericType, nil
}
