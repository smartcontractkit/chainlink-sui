package indexer

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/mr-tron/base58"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"

	crUtil "github.com/smartcontractkit/chainlink-sui/relayer/chainreader/chainreader_util"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/config"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/database"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
	"github.com/smartcontractkit/chainlink-sui/relayer/common"
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
	offrampPackageId        string
	latestOfframpPackageId  string
	executionEventModuleKey string
	executionEventKey       string
	configEventModuleKey    string
	configEventKey          string
	executeFunction         string

	// configs
	eventConfigs map[string]*config.ChainReaderEvent

	mu                    sync.RWMutex
	offrampPackageIdReady chan struct{}
	offrampPackageOnce    sync.Once
}

type TransactionsIndexerApi interface {
	Start(ctx context.Context) error
	SetOffRampPackage(pkg string, latestPkg string)
	Ready() error
	Close() error
}

func NewTransactionsIndexer(
	db sqlutil.DataSource,
	lggr logger.Logger,
	sdkClient client.SuiPTBClient,
	pollingInterval time.Duration,
	syncTimeout time.Duration,
	eventConfigs map[string]*config.ChainReaderEvent,
) TransactionsIndexerApi {
	logInstance := logger.Named(lggr, "SuiTransactionsIndexer")
	dataStore := database.NewDBStore(db, logInstance)

	return &TransactionsIndexer{
		db:                      dataStore,
		client:                  sdkClient,
		logger:                  logInstance,
		pollingInterval:         pollingInterval,
		syncTimeout:             syncTimeout,
		transmitters:            make(map[models.SuiAddress]string),
		executionEventModuleKey: "offramp",
		executionEventKey:       "ExecutionStateChanged",
		configEventModuleKey:    "ocr3_base",
		configEventKey:          "ConfigSet",
		executeFunction:         "init_execute",
		eventConfigs:            eventConfigs,
		offrampPackageIdReady:   make(chan struct{}),
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

// SetOffRampPackage sets offramp called by chainreader Bind.
func (t *TransactionsIndexer) SetOffRampPackage(pkg string, latestPkg string) {
	if pkg == "" {
		t.logger.Warn("SetOffRampPackage called with empty package id")
		return
	}
	t.mu.Lock()
	t.offrampPackageId = pkg
	t.latestOfframpPackageId = latestPkg
	t.mu.Unlock()

	t.logger.Infow("OffRamp package set", "offrampPackageId", pkg, "latestOfframpPackageId", latestPkg)

	t.offrampPackageOnce.Do(func() { close(t.offrampPackageIdReady) })
}

// waitForOffRampPackage blocks until the OffRamp package ID is available or the
// provided context is canceled.
func (t *TransactionsIndexer) waitForOffRampPackage(ctx context.Context) (string, error) {
	t.mu.RLock()
	pkg := t.offrampPackageId
	ch := t.offrampPackageIdReady
	t.mu.RUnlock()
	if pkg != "" {
		return pkg, nil
	}
	t.logger.Info("Waiting for OffRamp package...")
	select {
	case <-ch:
		t.mu.RLock()
		pkg = t.offrampPackageId
		t.mu.RUnlock()
		if pkg == "" {
			return "", fmt.Errorf("package ready signaled but empty")
		}
		return pkg, nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

// waitForInitialEvent method waits for the initial ExecutionStateChanged event to be indexed
// in the database before starting the transaction polling loop.
func (tIndexer *TransactionsIndexer) waitForInitialEvent(ctx context.Context) error {
	tIndexer.logger.Infow("waitForInitialEvent start", "idx_ptr", fmt.Sprintf("%p", tIndexer))

	moduleKey := tIndexer.configEventModuleKey // "ocr3_base"
	eventKey := tIndexer.configEventKey        // "ConfigSet"

	tIndexer.logger.Infof("Waiting for initial %s::%s event before starting transaction polling...", moduleKey, eventKey)

	// 1) Wait until Bind provides the OffRamp package (or ctx is cancelled)
	pkg, err := tIndexer.waitForOffRampPackage(ctx)
	if err != nil {
		tIndexer.logger.Infow("Transaction polling stopped during initial wait (no OffRamp pkg).")
		return err
	}
	tIndexer.logger.Infow("OffRamp package ready", "package", pkg)

	// 2) Poll the DB for the first ConfigSet event
	ticker := time.NewTicker(tIndexer.pollingInterval)
	defer ticker.Stop()

	for {
		eventAccountAddress := pkg
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

	tIndexer.logger.Debugw("syncTransmittersTransactions start", "transmitters", transmitters)

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

	tIndexer.logger.Debugw("syncTransmitterTransactions start", "transmitter", transmitter, "cursor", cursor)

	eventAccountAddress, latestOfframpPackageId, err := tIndexer.getEventPackageIdFromConfig()
	if err != nil {
		return 0, fmt.Errorf("failed to get ExecutionStateChanged event config: %w", err)
	}
	eventHandle := fmt.Sprintf("%s::%s::%s", eventAccountAddress, moduleKey, eventKey)

	select {
	case <-ctx.Done():
		return totalProcessed, ctx.Err()
	default:
		if cursor == "" {
			// Get the cursor from the DB store
			transmitterCursorFromDB, err := tIndexer.db.GetTransmitterCursor(ctx, transmitter)
			if err != nil {
				tIndexer.logger.Warnw("Failed to get transmitter cursor from DB store", "error", err)
			}
			// Attempt to check if a cursor exists in the DB store
			if transmitterCursorFromDB != "" {
				tIndexer.logger.Debugw("Found transmitter cursor in DB store", "transmitter", transmitter, "cursor", transmitterCursorFromDB)
				cursor = transmitterCursorFromDB
			} else {
				tIndexer.logger.Debugw("No transmitter cursor found in DB store, starting fresh sync", "transmitter", transmitter)
			}
		}

		queryResponse, err := tIndexer.client.QueryTransactions(ctx, string(transmitter), &cursor, &batchSize)
		if err != nil {
			return totalProcessed, fmt.Errorf("failed to fetch transactions for transmitter %s: %w", transmitter, err)
		}

		if len(queryResponse.Data) == 0 {
			return totalProcessed, nil
		}

		lastDigest := queryResponse.Data[len(queryResponse.Data)-1].Digest
		defer func() {
			// Update the cursor to the last transaction digest regardless of the code path below
			tIndexer.transmitters[transmitter] = lastDigest

			// Update the cursor in the DB store
			err := tIndexer.db.UpdateTransmitterCursor(ctx, transmitter, lastDigest)
			if err != nil {
				tIndexer.logger.Errorw("Failed to update transmitter cursor in DB store", "error", err)
			}
		}()

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

			tIndexer.logger.Debugw("Extracted move abort from failed transaction", "moveAbort", moveAbort, "digest", transactionRecord.Digest)

			executionMethodIndex := 0
			includesValidPTBCommand := false

			// Check if any of the transaction's commands match with the expected (offramp) package and module
			for i, raw := range transactionRecord.Transaction.Data.Transaction.Transactions {
				if moveCall := models.MoveCall(raw); moveCall != nil {
					packageID := moveCall.Package
					moduleName := moveCall.Module
					functionName := moveCall.Function

					if (packageID == eventAccountAddress || packageID == latestOfframpPackageId) &&
						moduleName == tIndexer.executionEventModuleKey &&
						functionName == tIndexer.executeFunction {
						executionMethodIndex = i
						includesValidPTBCommand = true
						break
					}
				}
			}

			// NOTE: The check below does not guarantee that a malicious (known) transmitter is not sending a failed PTB
			// with the expected package and module. However, it is considered as the worst case scenario simply involves
			// creating an event record with a failure state against an digest that is not checked.
			if !includesValidPTBCommand {
				tIndexer.logger.Warnw(
					"Expected PTB command (_::offramp::init_execute) not found in commands of failed PTB originating from known transmitter",
					"transmitter", transmitter,
					"digest", transactionRecord.Digest,
					"transactionRecord", transactionRecord,
				)
				continue
			}

			// The failure should NOT take place at `init_execute`. This command must be valid to ensure that the report can be extracted.
			if moveAbort.Location.FunctionName == nil || *moveAbort.Location.FunctionName == tIndexer.executeFunction {
				tIndexer.logger.Debugw("Skipping transaction for failed function against init_execute function",
					"transmitter", transmitter,
					"location", moveAbort.Location,
					"functionName", *moveAbort.Location.FunctionName,
					"digest", transactionRecord.Digest,
				)

				continue
			}

			// The command from which to extract the report (init_execute) should always be the first command in the PTB, however,
			// we use the index here as a simple safety check to ensure that the command index is valid in case that a function
			// is added before init_execute in the future.
			commandIndex := uint64(executionMethodIndex)
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
			reportValue, ok := reportArg["value"].([]any)
			if !ok {
				tIndexer.logger.Errorw("Expected report value to be a []any",
					"transmitter", transmitter,
					"txDigest", transactionRecord.Digest,
					"reportArg", reportArg,
					"valueType", fmt.Sprintf("%T", reportArg["value"]))
				continue
			}

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

			hasher := crUtil.NewMessageHasherV1(tIndexer.logger)
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
				// The conversion to []any is needed to avoid the default Go DB SDK behaviour of converting the byte slice to encoded base64 string.
				"message_id":   codec.BytesToAnySlice(execReport.Message.Header.MessageID),
				"message_hash": codec.BytesToAnySlice(messageHash[:]),
				"state":        uint8(3), // 3 = FAILURE
			}

			tIndexer.logger.Debugw("About to insert synthetic ExecutionStateChanged event", "executionStateChanged", executionStateChanged)

			// normalize keys
			executionStateChanged = common.ConvertMapKeysToCamelCase(executionStateChanged).(map[string]any)

			blockTimestamp, err := strconv.ParseUint(checkpointResponse.TimestampMs, 10, 64)
			if err != nil {
				tIndexer.logger.Errorw("Failed to parse block timestamp", "error", err)
				continue
			}

			// Convert the txDigest to hex
			txDigestHex := transactionRecord.Digest
			if base64Bytes, err := base58.Decode(txDigestHex); err == nil {
				hexTxId := hex.EncodeToString(base64Bytes)
				txDigestHex = "0x" + hexTxId
			}

			blockHashBytes, err := base58.Decode(checkpointResponse.Digest)
			if err != nil {
				tIndexer.logger.Errorw("Failed to decode block hash", "error", err)
				// fallback
				blockHashBytes = []byte(checkpointResponse.Digest)
			}

			record := database.EventRecord{
				EventAccountAddress: eventAccountAddress,
				EventHandle:         eventHandle,
				EventOffset:         0,
				TxDigest:            txDigestHex,
				BlockHeight:         checkpointResponse.SequenceNumber,
				BlockHash:           blockHashBytes,
				// Convert to seconds for consistency with events indexer.
				BlockTimestamp:      blockTimestamp / 1000,
				Data:                executionStateChanged,
				IsSynthetic:         true,
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

	eventAccountAddress, _, err := tIndexer.getEventPackageIdFromConfig()
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
		selector  = "sourceChainSelector"
	)

	eventAccountAddress, _, err := tIndexer.getEventPackageIdFromConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get SourceChainConfigSet event config: %w", err)
	}
	eventHandle := fmt.Sprintf("%s::%s::%s", eventAccountAddress, moduleKey, eventKey)

	filter := []query.Expression{
		query.Comparator(selector,
			primitives.ValueComparator{Value: strconv.FormatUint(sourceChainSelector, 10), Operator: primitives.Eq},
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

// Prefer the cached OffRamp package
func (t *TransactionsIndexer) getEventPackageIdFromConfig() (string, string, error) {
	t.mu.RLock()
	pkg := t.offrampPackageId
	latestPkg := t.latestOfframpPackageId
	t.mu.RUnlock()

	if pkg != "" {
		return pkg, latestPkg, nil
	}
	return "", "", fmt.Errorf("offramp package not set yet")
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
		} else if argIndex >= float64(len(inputCallArgs)) {
			return nil, fmt.Errorf("arg index out of range for failed transaction, argIndex: %d, inputCallArgs length: %d", int(argIndex), len(inputCallArgs))
		} else if inputCallArgs[uint64(argIndex)] == nil {
			return nil, fmt.Errorf("arg value is nil for failed transaction, argIndex: %d", int(argIndex))
		}

		commandArgs = append(commandArgs, inputCallArgs[uint64(argIndex)])
	}

	return commandArgs, nil
}

func (tIndexer *TransactionsIndexer) Ready() error {
	// TODO: implement
	return nil
}

func (tIndexer *TransactionsIndexer) Close() error {
	// TODO: implement
	return nil
}
