package indexer

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"time"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"

	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/config"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/database"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/util"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
)

type TransactionsIndexer struct {
	db              *database.DBStore
	client          client.SuiPTBClient
	logger          logger.Logger
	pollingInterval time.Duration
	syncTimeout     time.Duration
	// map of transmitter address to cursor (the last processed transaction digest)
	transmitters map[models.SuiAddress]string
	// event selectors
	executionEventModuleKey string
	executionEventKey       string
	configEventModuleKey    string
	configEventKey          string
	executeFunctions        []string
	// configs
	eventConfigs map[string]*config.ChainReaderEvent
}

type TransactionsIndexerApi interface {
	Start(ctx context.Context) error
	UpdateEventConfig(eventConfig *config.ChainReaderEvent)
}

func NewTransactionsIndexer(
	db *database.DBStore,
	lggr logger.Logger,
	sdkClient client.SuiPTBClient,
	pollingInterval time.Duration,
	syncTimeout time.Duration,
	eventConfigs map[string]*config.ChainReaderEvent,
) TransactionsIndexerApi {
	return &TransactionsIndexer{
		db:                      db,
		client:                  sdkClient,
		logger:                  lggr,
		pollingInterval:         pollingInterval,
		syncTimeout:             syncTimeout,
		transmitters:            make(map[models.SuiAddress]string),
		executionEventModuleKey: "offramp",
		executionEventKey:       "ExecutionStateChanged",
		configEventModuleKey:    "ocr3_base",
		configEventKey:          "ConfigSet",
		executeFunctions:        []string{"finish_execute"},
		eventConfigs:            eventConfigs,
	}
}

// Start method initiates the polling loop for the transactions indexer to enable
// indexing synthetic events for failed transactions.
func (tIndexer *TransactionsIndexer) Start(ctx context.Context) error {
	if err := tIndexer.waitForInitialEvent(ctx); err != nil {
		return err
	}

	tIndexer.logger.Infow("Transaction polling goroutine started")
	defer tIndexer.logger.Infow("Transaction polling goroutine exited")

	ticker := time.NewTicker(tIndexer.pollingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			syncCtx, cancel := context.WithTimeout(ctx, tIndexer.syncTimeout)
			start := time.Now()

			err := tIndexer.SyncAllTransmittersTransactions(syncCtx)
			elapsed := time.Since(start)

			if err != nil && !errors.Is(err, context.DeadlineExceeded) {
				tIndexer.logger.Warnw("TxSync completed with errors", "error", err, "duration", elapsed)
			} else if err != nil {
				tIndexer.logger.Warnw("Transaction sync timed out", "duration", elapsed)
			} else {
				tIndexer.logger.Debugw("Transaction sync completed successfully", "duration", elapsed)
			}

			cancel()
		case <-ctx.Done():
			tIndexer.logger.Infow("Transaction polling stopped")
			return nil
		}
	}
}

// UpdateEventConfig method either edits or inserts the event config into the map of configs
func (tIndexer *TransactionsIndexer) UpdateEventConfig(eventConfig *config.ChainReaderEvent) {
	key := fmt.Sprintf("%s::%s", eventConfig.Module, eventConfig.Name)
	tIndexer.eventConfigs[key] = eventConfig
}

// waitForInitialEvent method waits for the initial ExecutionStateChanged event to be indexed
// in the database before starting the transaction polling loop.
func (tIndexer *TransactionsIndexer) waitForInitialEvent(ctx context.Context) error {
	var (
		moduleKey = tIndexer.configEventModuleKey
		eventKey  = tIndexer.configEventKey
	)

	tIndexer.logger.Infof("Waiting for initial %s::%s event before starting transaction polling...", moduleKey, eventKey)

	ticker := time.NewTicker(tIndexer.pollingInterval)
	defer ticker.Stop()

	for {
		eventAccountAddress, err := tIndexer.getEventPackageIdFromConfig(moduleKey, eventKey)
		if err != nil {
			tIndexer.logger.Warnw(fmt.Sprintf("Failed to get %s::%s event config, retrying...", moduleKey, eventKey), "error", err)
		} else {
			eventHandle := fmt.Sprintf("%s::%s::%s", eventAccountAddress, moduleKey, eventKey)
			events, err := tIndexer.db.QueryEvents(
				ctx,
				eventAccountAddress,
				eventHandle,
				[]query.Expression{},
				query.LimitAndSort{
					Limit: query.CountLimit(1),
					SortBy: []query.SortBy{
						query.NewSortBySequence(query.Desc),
					},
				},
			)
			if err != nil {
				tIndexer.logger.Warnw(fmt.Sprintf("Failed to query for %s::%s events, retrying...", moduleKey, eventKey), "error", err)
			} else if len(events) > 0 {
				tIndexer.logger.Infow(fmt.Sprintf("Found initial %s::%s event, starting tx poller.", moduleKey, eventKey), "count", len(events))
				return nil
			}
		}

		select {
		case <-ticker.C:
			tIndexer.logger.Infow(fmt.Sprintf("No %s::%s events found yet, waiting...", moduleKey, eventKey))
			continue
		case <-ctx.Done():
			tIndexer.logger.Infow(fmt.Sprintf("Transaction polling stopped during initial wait for %s::%s event.", moduleKey, eventKey))
			return ctx.Err()
		}
	}
}

// SyncTransmittersTransactions method syncs the transactions for each known transmitter.
func (tIndexer *TransactionsIndexer) SyncAllTransmittersTransactions(ctx context.Context) error {
	transmitters, err := tIndexer.getTransmitters(ctx)
	if err != nil {
		return fmt.Errorf("failed to get transmitters: %w", err)
	}

	if len(transmitters) == 0 {
		return nil
	}

	var batchSize uint64 = 50
	var totalProcessed int

	for _, transmitter := range transmitters {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if _, exists := tIndexer.transmitters[transmitter]; !exists {
				tIndexer.logger.Debugw("Initializing cursor for transmitter", "transmitter", transmitter)
				tIndexer.transmitters[transmitter] = ""
			}

			processed, err := tIndexer.syncTransmitterTransactions(ctx, transmitter, batchSize)
			if err != nil {
				tIndexer.logger.Errorw("Failed to sync transmitter transactions", "transmitter", transmitter, "error", err)

				continue
			}
			totalProcessed += processed
		}
	}

	if totalProcessed > 0 {
		tIndexer.logger.Debugw("All transmitters' failed transactions processed", "totalProcessed", totalProcessed)
	}

	return nil
}

func (tIndexer *TransactionsIndexer) syncTransmitterTransactions(ctx context.Context, transmitter models.SuiAddress, batchSize uint64) (int, error) {
	var (
		moduleKey = tIndexer.executionEventModuleKey
		eventKey  = tIndexer.executionEventKey
	)

	cursor := tIndexer.transmitters[transmitter]
	totalProcessed := 0

	eventAccountAddress, err := tIndexer.getEventPackageIdFromConfig(moduleKey, eventKey)
	if err != nil {
		return 0, fmt.Errorf("failed to get ExecutionStateChanged event config: %w", err)
	}
	eventHandle := fmt.Sprintf("%s::%s::%s", eventAccountAddress, moduleKey, eventKey)

	select {
	case <-ctx.Done():
		return totalProcessed, ctx.Err()
	default:
		queryResponse, err := tIndexer.client.QueryTransactions(ctx, string(transmitter), &cursor, &batchSize)
		if err != nil {
			return totalProcessed, fmt.Errorf("failed to fetch transactions for transmitter %s: %w", transmitter, err)
		}

		if len(queryResponse.Data) == 0 {
			return totalProcessed, nil
		}

		var records []database.EventRecord
		for _, transactionRecord := range queryResponse.Data {
			if transactionRecord.Effects.Status.Status == "success" {
				tIndexer.logger.Debugw("Skipping successful transaction",
					"transmitter", transmitter, "digest", transactionRecord.Digest)

				continue
			}

			tIndexer.logger.Infow("Found failed transaction",
				"transmitter", transmitter, "digest", transactionRecord.Digest)

			if transactionRecord.Transaction.Data.Transaction.Kind != "ProgrammableTransaction" {
				tIndexer.logger.Debugw("Skipping non-programmable transaction",
					"transmitter", transmitter, "digest", transactionRecord.Digest)

				continue
			}

			// get the checkpoint / block details
			checkpointResponse, err := tIndexer.client.GetBlockById(ctx, transactionRecord.Checkpoint)
			if err != nil {
				tIndexer.logger.Errorw("Failed to get checkpoint", "error", err)
				continue
			}

			// parse the transaction error
			errMessage := transactionRecord.Effects.Status.Error
			moveAbort, err := tIndexer.parseMoveAbort(errMessage)
			if err != nil {
				tIndexer.logger.Errorw("Failed to parse move abort", "error", err)
				continue
			}

			if moveAbort.Location.Module.Name != moduleKey {
				tIndexer.logger.Debugw("Skipping transaction with different module",
					"transmitter", transmitter, "module", moveAbort.Location.Module.Name)

				continue
			}

			if moveAbort.Location.FunctionName == nil || !slices.Contains(tIndexer.executeFunctions, *moveAbort.Location.FunctionName) {
				tIndexer.logger.Debugw("Skipping transaction for non-execute function",
					"transmitter", transmitter, "function", *moveAbort.Location.FunctionName)

				continue
			}

			// we always get the report from the init_execute function call (index 0), the "finish_execute" function call
			// does not contain an argument which contains the report
			// NOTE: we assume that init_execute (which contains the report) is always the first command in the PTB
			commandIndex := uint64(0)
			callArgs, err := tIndexer.extractCommandCallArgs(&transactionRecord, commandIndex)
			if err != nil {
				tIndexer.logger.Errorw("Failed to extract command call args", "error", err)
				continue
			}

			tIndexer.logger.Debugw("Extracted command call args in transactions indexer", "transmitter", transmitter, "txDigest", transactionRecord.Digest, "args", callArgs)

			if len(callArgs) < 5 {
				tIndexer.logger.Errorw("Expected report to be a hex string", "transmitter", transmitter, "txDigest", transactionRecord.Digest, "callArgs", callArgs)
				continue
			}

			reportArg := callArgs[4]
			tIndexer.logger.Debugw("Report arg", "reportArg", reportArg)

			// Handle the conversion from []interface{} to []byte
			reportValue := reportArg["value"].([]any)
			reportBytes := make([]byte, len(reportValue))
			for i, val := range reportValue {
				num, ok := val.(float64)
				if !ok {
					tIndexer.logger.Errorw("Expected numeric value in byte array",
						"transmitter", transmitter, "txDigest", transactionRecord.Digest, "value", val, "type", fmt.Sprintf("%T", val))

					continue
				}
				reportBytes[i] = byte(num)
			}

			tIndexer.logger.Infow("Report bytes", "reportBytes", reportBytes)

			execReport, err := codec.DeserializeExecutionReport(reportBytes)
			if err != nil {
				tIndexer.logger.Errorw("Failed to deserialize execution report",
					"transmitter", transmitter, "txDigest", transactionRecord.Digest, "error", err)

				continue
			}

			tIndexer.logger.Debugw("Deserialized execution report", "execReport", execReport)

			sourceChainSelector := execReport.Message.Header.SourceChainSelector
			sourceChainConfig, err := tIndexer.getSourceChainConfig(ctx, sourceChainSelector)
			if err != nil {
				tIndexer.logger.Errorw("Failed to get source chain config",
					"transmitter", transmitter, "sourceChainSelector", sourceChainSelector, "error", err)

				continue
			}

			if sourceChainConfig == nil {
				tIndexer.logger.Debugw("No source chain config found for selector",
					"transmitter", transmitter, "sourceChainSelector", sourceChainSelector)

				continue
			}

			tIndexer.logger.Debugw("Source chain config", "sourceChainConfig", sourceChainConfig)
			tIndexer.logger.Debugw("Execution report", "execReport", execReport)

			hasher := util.NewMessageHasherV1(tIndexer.logger)
			messageHash, err := hasher.Hash(ctx, execReport, sourceChainConfig.OnRamp)
			if err != nil {
				tIndexer.logger.Errorw("Failed to calculate message hash",
					"transmitter", transmitter, "txDigest", transactionRecord.Digest, "error", err)

				continue
			}

			// Create synthetic ExecutionStateChanged event
			// The fields map one-to-one the onchain event
			executionStateChanged := map[string]any{
				"source_chain_selector": fmt.Sprintf("%d", sourceChainSelector),
				"sequence_number":       fmt.Sprintf("%d", execReport.Message.Header.SequenceNumber),
				"message_id":            "0x" + hex.EncodeToString(execReport.Message.Header.MessageID),
				"message_hash":          "0x" + hex.EncodeToString(messageHash[:]),
				"state":                 uint8(3), // 3 = FAILURE
			}

			blockTimestamp, err := strconv.ParseUint(checkpointResponse.TimestampMs, 10, 64)
			if err != nil {
				tIndexer.logger.Errorw("Failed to parse block timestamp", "error", err)
				continue
			}

			record := database.EventRecord{
				EventAccountAddress: eventAccountAddress,
				EventHandle:         eventHandle,
				EventOffset:         0,
				TxDigest:            transactionRecord.Digest,
				BlockHeight:         checkpointResponse.SequenceNumber,
				BlockHash:           []byte(checkpointResponse.Digest),
				BlockTimestamp:      blockTimestamp,
				Data:                executionStateChanged,
			}

			records = append(records, record)
			totalProcessed++
		}

		if len(records) > 0 {
			// Try batch insert first
			if err := tIndexer.db.InsertEvents(ctx, records); err != nil {
				tIndexer.logger.Errorw("Batch insert failed, falling back to per-event insert", "error", err)
				// Fallback: insert each record individually, skip bad ones
				totalProcessedFallback := 0
				for _, record := range records {
					if err := tIndexer.db.InsertEvents(ctx, []database.EventRecord{record}); err != nil {
						tIndexer.logger.Errorw("Failed to insert single synthetic event, skipping",
							"error", err,
							"transmitter", transmitter,
							"txDigest", record.TxDigest)

						continue
					}

					totalProcessedFallback++
				}
				tIndexer.logger.Debugw("Inserted synthetic ExecutionStateChanged events", "count", totalProcessed, "transmitter", transmitter)

				return totalProcessedFallback, nil
			}

			// update the cursor to the last transaction digest
			tIndexer.transmitters[transmitter] = queryResponse.Data[len(queryResponse.Data)-1].Digest
			tIndexer.logger.Debugw("Inserted synthetic ExecutionStateChanged events",
				"count", len(records), "transmitter", transmitter)
		}

		tIndexer.logger.Debugw("Inserted synthetic ExecutionStateChanged events", "records", records)

		return totalProcessed, nil
	}
}

// getTransmitters method retrieves the transmitters from the OCRConfigSet event in the 'ocr3_base.move' contract.
func (tIndexer *TransactionsIndexer) getTransmitters(ctx context.Context) ([]models.SuiAddress, error) {
	var (
		moduleKey = tIndexer.configEventModuleKey
		eventKey  = tIndexer.configEventKey
	)

	eventAccountAddress, err := tIndexer.getEventPackageIdFromConfig(moduleKey, eventKey)
	if err != nil {
		tIndexer.logger.Errorw("Failed to get OCRConfigSet event config", "error", err)
		return nil, err
	}
	eventHandle := fmt.Sprintf("%s::%s::%s", eventAccountAddress, moduleKey, eventKey)

	events, err := tIndexer.db.QueryEvents(
		ctx,
		eventAccountAddress,
		eventHandle,
		[]query.Expression{},
		query.LimitAndSort{
			Limit: query.CountLimit(1),
			SortBy: []query.SortBy{
				query.NewSortBySequence(query.Desc),
			},
		},
	)

	if err != nil {
		tIndexer.logger.Errorw("Failed to query OCRConfigSet events", "error", err)
		return nil, err
	}

	if len(events) == 0 {
		tIndexer.logger.Warnw("No OCRConfigSet events found")
		return nil, nil
	}

	var configSet codec.ConfigSet
	if err := codec.DecodeSuiJsonValue(events[0].Data, &configSet); err != nil {
		tIndexer.logger.Errorw("Failed to decode ConfigSet event", "error", err)
		return nil, fmt.Errorf("failed to decode ConfigSet event: %w", err)
	}

	tIndexer.logger.Infow("Found ConfigSet event", "data", events[0].Data)

	transmitters := configSet.Transmitters
	if len(transmitters) == 0 {
		tIndexer.logger.Warnw("`No transmitters` found in OCRConfigSet event")
		return nil, nil
	}

	suiAddresses := make([]models.SuiAddress, 0, len(transmitters))
	for _, transmitter := range transmitters {
		suiAddresses = append(suiAddresses, models.SuiAddress(transmitter))
	}

	tIndexer.logger.Infow("Found transmitters in OCRConfigSet event", "count", len(suiAddresses))

	return suiAddresses, nil
}

func (tIndexer *TransactionsIndexer) getSourceChainConfig(ctx context.Context, sourceChainSelector uint64) (*codec.SourceChainConfig, error) {
	const (
		moduleKey = "offramp"
		eventKey  = "SourceChainConfigSet"
		selector  = "SourceChainSelector"
	)

	eventAccountAddress, err := tIndexer.getEventPackageIdFromConfig(moduleKey, eventKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get SourceChainConfigSet event config: %w", err)
	}
	eventHandle := fmt.Sprintf("%s::%s::%s", eventAccountAddress, moduleKey, eventKey)

	filter := []query.Expression{
		query.Comparator(selector,
			primitives.ValueComparator{Value: sourceChainSelector, Operator: primitives.Eq},
		),
	}

	events, err := tIndexer.db.QueryEvents(
		ctx,
		eventAccountAddress,
		eventHandle,
		filter,
		query.LimitAndSort{
			Limit: query.CountLimit(1),
			SortBy: []query.SortBy{
				query.NewSortBySequence(query.Desc),
			},
		},
	)

	if err != nil {
		return nil, fmt.Errorf("failed to query SourceChainConfigSet event: %w", err)
	}

	if len(events) == 0 {
		tIndexer.logger.Debugw("No SourceChainConfigSet event found", "sourceChainSelector", sourceChainSelector)
		//nolint:nilnil
		return nil, nil
	}

	var configEvent codec.SourceChainConfigSet
	if err := codec.DecodeSuiJsonValue(events[0].Data, &configEvent); err != nil {
		return nil, fmt.Errorf("failed to decode SourceChainConfigSet event: %w", err)
	}

	return &configEvent.SourceChainConfig, nil
}

// getEventConfig method retrieves the event config from the map of configs
// returns the package, module, and event account address for the event config
func (tIndexer *TransactionsIndexer) getEventPackageIdFromConfig(moduleKey, eventKey string) (string, error) {
	key := fmt.Sprintf("%s::%s", moduleKey, eventKey)

	if eventConfig, ok := tIndexer.eventConfigs[key]; ok {
		if eventConfig.Package == "" {
			return "", fmt.Errorf("event package ID not found for %s", key)
		}

		return eventConfig.Package, nil
	}

	return "", fmt.Errorf("event config not found for %s", key)
}

// ModuleId represents Move’s ModuleId { address, name }
type ModuleId struct {
	Address string
	Name    string
}

// MoveLocation corresponds to MoveLocation { module, function, instruction, function_name }
type MoveLocation struct {
	Module       ModuleId
	Function     uint64
	Instruction  uint64
	FunctionName *string // nil if None
}

// MoveAbort wraps a MoveLocation plus abort code and PTB command index
type MoveAbort struct {
	Location     MoveLocation
	AbortCode    uint64
	CommandIndex uint64
}

// regex to capture:
//
//	1: address (hex)
//	2: module name
//	3: function (decimal)
//	4: instruction (decimal)
//	5: either Some("X") or None
//	6: inner X from Some("X") (empty if None)
//	7: abort code
//	8: command index
var abortRe = regexp.MustCompile(
	`^MoveAbort\(` +
		`MoveLocation \{ module: ModuleId \{ address: ([0-9a-f]+), name: Identifier\("([^"]+)"\) \}, ` +
		`function: (\d+), instruction: (\d+), function_name: (Some\("([^"]+)"\)|None) \}, ` +
		`(\d+)\) in command (\d+)$`,
)

// ParseMoveAbort parses the error string into a MoveAbort struct.
func (tIndexer *TransactionsIndexer) parseMoveAbort(s string) (*MoveAbort, error) {
	m := abortRe.FindStringSubmatch(s)
	if m == nil {
		return nil, fmt.Errorf("input does not match MoveAbort pattern")
	}
	// m[1]=address, m[2]=modName, m[3]=func, m[4]=instr,
	// m[5]=full (Some("…")|None), m[6]=inner name or "",
	// m[7]=abortCode, m[8]=cmdIndex

	// parse integers
	fn, err := strconv.ParseUint(m[3], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("bad function number: %w", err)
	}
	instr, err := strconv.ParseUint(m[4], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("bad instruction number: %w", err)
	}
	abortCode, err := strconv.ParseUint(m[7], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("bad abort code: %w", err)
	}
	cmdIdx, err := strconv.ParseUint(m[8], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("bad command index: %w", err)
	}

	// optional function name
	var fname *string
	if m[5] != "None" {
		fname = new(string)
		*fname = m[6]
	}

	loc := MoveLocation{
		Module: ModuleId{
			Address: m[1],
			Name:    m[2],
		},
		Function:     fn,
		Instruction:  instr,
		FunctionName: fname,
	}

	return &MoveAbort{
		Location:     loc,
		AbortCode:    abortCode,
		CommandIndex: cmdIdx,
	}, nil
}

// extractCommandCallArgs zips the input indices with the input call args to output a slice of call arg details
func (tIndexer *TransactionsIndexer) extractCommandCallArgs(transactionRecord *models.SuiTransactionBlockResponse, commandIndex uint64) ([]models.SuiCallArg, error) {
	// this refers to the indexed inputs of the command call which failed
	commandDetails, ok := transactionRecord.Transaction.Data.Transaction.Transactions[commandIndex].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("failed to read command details for failed transaction")
	}
	// this refers to the indexed inputs of the entire PTB transaction
	inputCallArgs := transactionRecord.Transaction.Data.Transaction.Inputs

	moveCall, ok := commandDetails["MoveCall"].(map[string]any)
	if !ok {
		tIndexer.logger.Debugw("Failed to read MoveCall details for failed transaction", "commandDetails", commandDetails)
		return nil, fmt.Errorf("failed to read MoveCall details for failed transaction")
	}

	moveCallArguments, ok := moveCall["arguments"].([]any)
	if !ok {
		tIndexer.logger.Debugw("Failed to read MoveCall arguments for failed transaction", "moveCall", moveCall)
		return nil, fmt.Errorf("failed to read MoveCall arguments for failed transaction")
	}

	// construct a slice of call arg details based on the command call arguments
	commandArgs := make([]models.SuiCallArg, 0)
	for _, arg := range moveCallArguments {
		argEntry, ok := arg.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("failed to read arg entry for failed transaction")
		}
		argIndex, ok := argEntry["Input"].(float64)
		if !ok {
			return nil, fmt.Errorf("failed to read arg index for failed transaction")
		}
		commandArgs = append(commandArgs, inputCallArgs[uint64(argIndex)])
	}

	return commandArgs, nil
}
