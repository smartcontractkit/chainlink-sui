package indexer

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	suirpcv2 "github.com/block-vision/sui-go-sdk/pb/sui/rpc/v2"
	"github.com/mr-tron/base58"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"

	chainreader_util "github.com/smartcontractkit/chainlink-sui/relayer/chainreader/chainreader_util"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/config"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/database"
	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
	"github.com/smartcontractkit/chainlink-sui/relayer/common"
)

// TransactionsIndexer consumes checkpoint transaction batches from a channel,
// identifies failed transactions from known transmitters, and creates synthetic
// ExecutionStateChanged events for those failures.
type TransactionsIndexer struct {
	db     *database.DBStore
	logger logger.Logger
	wg     sync.WaitGroup

	// Event configs for reading source chain config
	eventConfigs map[string]*config.ChainReaderEvent

	// Offramp package tracking
	offrampPackageID       string
	latestOfframpPackageID string
	mu                     sync.RWMutex
	offrampPackageIDReady  chan struct{}
	offrampPackageOnce     sync.Once

	// Constants for event/module identification
	executionEventModuleKey string
	executionEventKey       string
	configEventModuleKey    string
	configEventKey          string
	executeFunction         string

	starter services.StateMachine
}

// TransactionsIndexerApi defines the interface for the transactions indexer.
type TransactionsIndexerApi interface {
	// Start begins consuming from the transactions channel.
	// The channel is closed when the poller stops.
	Start(ctx context.Context, transactionsCh <-chan CheckpointTransactionsBatch) error
	// ProcessCheckpointTransactions processes a batch of transactions from a single checkpoint.
	ProcessCheckpointTransactions(ctx context.Context, batch CheckpointTransactionsBatch) error
	// SetOffRampPackage sets the offramp package IDs.
	SetOffRampPackage(pkg string, latestPkg string)
	Ready() error
	Close() error
}

// NewTransactionsIndexer creates a new TransactionsIndexer instance.
func NewTransactionsIndexer(
	db sqlutil.DataSource,
	lggr logger.Logger,
	eventConfigs map[string]*config.ChainReaderEvent,
) TransactionsIndexerApi {
	logInstance := logger.Named(lggr, "SuiTransactionsIndexer")
	dataStore := database.NewDBStore(db, logInstance)

	return &TransactionsIndexer{
		db:                      dataStore,
		logger:                  logInstance,
		eventConfigs:            eventConfigs,
		executionEventModuleKey: "offramp",
		executionEventKey:       "ExecutionStateChanged",
		configEventModuleKey:    "ocr3_base",
		configEventKey:          "ConfigSet",
		executeFunction:         "init_execute",
		offrampPackageIDReady:   make(chan struct{}),
	}
}

// Start begins consuming transactions from the provided channel.
// The consumer drains the channel immediately so the ChainPoller is not blocked
// while waiting for the initial ConfigSet event to identify transmitters.
func (tIndexer *TransactionsIndexer) Start(ctx context.Context, transactionsCh <-chan CheckpointTransactionsBatch) error {
	return tIndexer.starter.StartOnce("TransactionsIndexer", func() error {
		tIndexer.wg.Add(1)
		go func() {
			defer tIndexer.wg.Done()
			tIndexer.run(ctx, transactionsCh)
		}()

		return nil
	})
}

// run is the main consumer loop.
func (tIndexer *TransactionsIndexer) run(ctx context.Context, transactionsCh <-chan CheckpointTransactionsBatch) {
	tIndexer.logger.Info("TransactionsIndexer starting")

	ready := make(chan error, 1)
	go func() {
		ready <- tIndexer.waitForInitialEvent(ctx)
	}()

	processingEnabled := false

	for {
		select {
		case <-ctx.Done():
			tIndexer.logger.Info("TransactionsIndexer stopping")
			return
		case err := <-ready:
			if err != nil {
				tIndexer.logger.Infow("Transaction processing disabled", "error", err)
				continue
			}
			tIndexer.logger.Info("TransactionsIndexer ready to process checkpoint transactions")
			processingEnabled = true
		case batch, ok := <-transactionsCh:
			if !ok {
				tIndexer.logger.Info("TransactionsIndexer channel closed, exiting")
				return
			}
			if !processingEnabled {
				continue
			}
			if err := tIndexer.ProcessCheckpointTransactions(ctx, batch); err != nil {
				tIndexer.logger.Errorw("Failed to process checkpoint transactions",
					"sequence", batch.Checkpoint.SequenceNumber,
					"error", err)
			}
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
	t.offrampPackageID = pkg
	t.latestOfframpPackageID = latestPkg
	t.mu.Unlock()

	t.logger.Infow("OffRamp package set", "offrampPackageId", pkg, "latestOfframpPackageId", latestPkg)

	t.offrampPackageOnce.Do(func() { close(t.offrampPackageIDReady) })
}

// waitForOffRampPackage blocks until the OffRamp package ID is available or the
// provided context is canceled.
func (t *TransactionsIndexer) waitForOffRampPackage(ctx context.Context) (string, error) {
	t.mu.RLock()
	pkg := t.offrampPackageID
	ch := t.offrampPackageIDReady
	t.mu.RUnlock()
	if pkg != "" {
		return pkg, nil
	}
	t.logger.Info("Waiting for OffRamp package...")
	select {
	case <-ch:
		t.mu.RLock()
		pkg = t.offrampPackageID
		t.mu.RUnlock()
		if pkg == "" {
			return "", fmt.Errorf("package ready signaled but empty")
		}
		return pkg, nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

// waitForInitialEvent waits for the initial ConfigSet event to be indexed.
func (tIndexer *TransactionsIndexer) waitForInitialEvent(ctx context.Context) error {
	tIndexer.logger.Infow("waitForInitialEvent start", "idx_ptr", fmt.Sprintf("%p", tIndexer))

	moduleKey := tIndexer.configEventModuleKey
	eventKey := tIndexer.configEventKey

	tIndexer.logger.Infof("Waiting for initial %s::%s event before starting transaction processing...", moduleKey, eventKey)

	// Wait until Bind provides the OffRamp package
	pkg, err := tIndexer.waitForOffRampPackage(ctx)
	if err != nil {
		tIndexer.logger.Infow("Transaction processing stopped during initial wait (no OffRamp pkg).")
		return err
	}
	tIndexer.logger.Infow("OffRamp package ready", "package", pkg)

	// Poll the DB for the first ConfigSet event
	ticker := time.NewTicker(100 * time.Millisecond)
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
			tIndexer.logger.Infow(fmt.Sprintf("Found initial %s::%s event, starting tx processing.", moduleKey, eventKey), "count", len(events))
			return nil
		}

		select {
		case <-ticker.C:
			// tIndexer.logger.Infow(fmt.Sprintf("No %s::%s events found yet, waiting...", moduleKey, eventKey))
			continue
		case <-ctx.Done():
			tIndexer.logger.Infow(fmt.Sprintf("Transaction processing stopped during initial wait for %s::%s event.", moduleKey, eventKey))
			return ctx.Err()
		}
	}
}

// ProcessCheckpointTransactions processes a batch of transactions from a single checkpoint.
func (tIndexer *TransactionsIndexer) ProcessCheckpointTransactions(ctx context.Context, batch CheckpointTransactionsBatch) error {
	transmitters, err := tIndexer.getTransmitters(ctx)
	if err != nil {
		return fmt.Errorf("failed to get transmitters: %w", err)
	}

	if len(transmitters) == 0 {
		return nil
	}

	// Build transmitter set for O(1) lookup
	transmitterSet := make(map[string]struct{})
	for _, t := range transmitters {
		transmitterSet[t] = struct{}{}
	}

	tIndexer.logger.Debugw("Processing checkpoint transactions",
		"sequence", batch.Checkpoint.SequenceNumber,
		"transactions", len(batch.Transactions),
		"transmitters", len(transmitters))

	eventAccountAddress, latestOfframpPackageId, err := tIndexer.getEventPackageIdFromConfig()
	if err != nil {
		return fmt.Errorf("failed to get event package: %w", err)
	}
	eventHandle := fmt.Sprintf("%s::%s::%s", eventAccountAddress, tIndexer.executionEventModuleKey, tIndexer.executionEventKey)

	var records []database.EventRecord

	for _, tx := range batch.Transactions {
		if !tIndexer.shouldProcessTransaction(tx, transmitterSet) {
			continue
		}

		record, err := tIndexer.processFailedTransaction(ctx, tx, eventHandle, eventAccountAddress, latestOfframpPackageId, batch.Checkpoint)
		if err != nil {
			tIndexer.logger.Errorw("Failed to process failed transaction",
				"txDigest", tx.GetDigest(),
				"error", err)
			continue
		}
		if record != nil {
			records = append(records, *record)
		}
	}

	if len(records) == 0 {
		return nil
	}

	// Batch insert with fallback
	if err := tIndexer.db.InsertEvents(ctx, records); err != nil {
		tIndexer.logger.Errorw("Batch insert failed, falling back to per-event insert", "error", err)

		for _, record := range records {
			if err := tIndexer.db.InsertEvents(ctx, []database.EventRecord{record}); err != nil {
				tIndexer.logger.Errorw("Failed to insert single synthetic event, skipping",
					"error", err,
					"txDigest", record.TxDigest)
			}
		}
	}

	tIndexer.logger.Debugw("Inserted synthetic ExecutionStateChanged events",
		"count", len(records))

	return nil
}

// shouldProcessTransaction determines if a transaction should be processed.
func (tIndexer *TransactionsIndexer) shouldProcessTransaction(tx *suirpcv2.ExecutedTransaction, transmitterSet map[string]struct{}) bool {
	// Must be from a known transmitter
	sender := tx.GetTransaction().GetSender()
	if sender == "" {
		return false
	}
	if _, ok := transmitterSet[sender]; !ok {
		return false
	}

	// Must be a failed transaction
	if tx.GetEffects().GetStatus().GetSuccess() {
		return false
	}

	// Must be a programmable transaction
	kind := tx.GetTransaction().GetKind()
	if kind == nil || kind.GetProgrammableTransaction() == nil {
		return false
	}

	return true
}

// processFailedTransaction processes a single failed transaction and creates a synthetic event record.
func (tIndexer *TransactionsIndexer) processFailedTransaction(
	ctx context.Context,
	tx *suirpcv2.ExecutedTransaction,
	eventHandle string,
	eventAccountAddress string,
	latestOfframpPackageID string,
	checkpoint CheckpointMeta,
) (*database.EventRecord, error) {
	tIndexer.logger.Debugw("Processing failed transaction",
		"txDigest", tx.GetDigest(),
		"eventHandle", eventHandle,
		"eventAccountAddress", eventAccountAddress,
		"latestOfframpPackageId", latestOfframpPackageID,
		"checkpoint", checkpoint)

	execErr := tx.GetEffects().GetStatus().GetError()
	moveAbort, err := tIndexer.parseMoveAbortFromExecutionError(execErr)
	if err != nil {
		tIndexer.logger.Debugw("Failed to parse move abort (not necessarily an error)", "error", err, "moveErrString", execErr)
		// Continue processing - some failures may not have abort info
		return nil, nil
	}

	// Find init_execute command and extract call args
	kind := tx.GetTransaction().GetKind()
	if kind == nil {
		return nil, errors.New("transaction has no kind")
	}
	programmableTx := kind.GetProgrammableTransaction()
	if programmableTx == nil {
		return nil, errors.New("not a programmable transaction")
	}

	var executionMethodIndex int
	includesValidPTBCommand := false

	for i, cmd := range programmableTx.Commands {
		moveCall := cmd.GetMoveCall()
		if moveCall == nil {
			continue
		}

		pkg := moveCall.GetPackage()
		module := moveCall.GetModule()
		function := moveCall.GetFunction()

		if (pkg == eventAccountAddress || pkg == latestOfframpPackageID) &&
			module == tIndexer.executionEventModuleKey &&
			function == tIndexer.executeFunction {
			executionMethodIndex = i
			includesValidPTBCommand = true
			break
		}
	}

	if !includesValidPTBCommand {
		tIndexer.logger.Warnw("Expected PTB command not found",
			"txDigest", tx.GetDigest())
		return nil, nil
	}

	// The failure should NOT take place at init_execute
	if moveAbort.Location.FunctionName != nil &&
		*moveAbort.Location.FunctionName == tIndexer.executeFunction {
		tIndexer.logger.Debugw("Skipping - failure at init_execute",
			"txDigest", tx.GetDigest())
		return nil, nil
	}

	// Extract call args for the report from the ProgrammableTransaction
	if executionMethodIndex >= len(programmableTx.Commands) {
		return nil, fmt.Errorf("command index %d out of range", executionMethodIndex)
	}

	moveCall := programmableTx.Commands[executionMethodIndex].GetMoveCall()
	if moveCall == nil {
		return nil, errors.New("command is not a MoveCall")
	}

	// Extract arguments from the move call
	inputs := programmableTx.Inputs
	var callArgs []*suirpcv2.Input
	for _, arg := range moveCall.Arguments {
		// Only process INPUT type arguments
		if arg.GetKind() != suirpcv2.Argument_INPUT {
			continue
		}
		inputIdx := arg.GetInput()
		if int(inputIdx) >= len(inputs) {
			return nil, fmt.Errorf("input index %d out of range", inputIdx)
		}
		callArgs = append(callArgs, inputs[inputIdx])
	}

	if len(callArgs) < 5 {
		tIndexer.logger.Errorw("Expected report in arg position 4", "callArgs", len(callArgs))
		return nil, nil
	}

	// Extract report bytes from input - Pure input type contains raw bytes
	reportInput := callArgs[4]
	if reportInput == nil {
		return nil, errors.New("report input is nil")
	}
	reportBytes := reportInput.GetPure()
	if reportBytes == nil {
		return nil, errors.New("report input is not Pure type")
	}

	// Deserialize execution report
	execReport, err := codec.DeserializeExecutionReportFromPure(reportBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize execution report: %w", err)
	}

	sourceChainSelector := execReport.Message.Header.SourceChainSelector
	sourceChainConfig, err := tIndexer.getSourceChainConfig(ctx, sourceChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to get source chain config: %w", err)
	}
	if sourceChainConfig == nil {
		tIndexer.logger.Debugw("No source chain config found",
			"sourceChainSelector", sourceChainSelector)
		return nil, nil
	}

	// Calculate message hash
	hasher := chainreader_util.NewMessageHasherV1(tIndexer.logger)
	messageHash, err := hasher.Hash(ctx, execReport, sourceChainConfig.OnRamp)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate message hash: %w", err)
	}

	// Create synthetic ExecutionStateChanged event
	executionStateChanged := map[string]any{
		"source_chain_selector": strconv.FormatUint(sourceChainSelector, 10),
		"sequence_number":       strconv.FormatUint(execReport.Message.Header.SequenceNumber, 10),
		"message_id":            codec.BytesToAnySlice(execReport.Message.Header.MessageID),
		"message_hash":          codec.BytesToAnySlice(messageHash[:]),
		"state":                 uint8(3), // 3 = FAILURE
	}

	// Normalize keys to camelCase
	if normalized := common.ConvertMapKeysToCamelCase(executionStateChanged); normalized != nil {
		if mapData, ok := normalized.(map[string]any); ok {
			executionStateChanged = mapData
		}
	}

	// Convert tx digest to hex format
	txDigest := tx.GetDigest()
	txDigestHex := txDigest
	if !strings.HasPrefix(txDigest, "0x") {
		if base58Bytes, err := base58.Decode(txDigest); err == nil {
			txDigestHex = "0x" + hex.EncodeToString(base58Bytes)
		}
	}

	// Convert block hash to bytes
	blockHashBytes := []byte(checkpoint.Digest)
	if decoded, err := base58.Decode(checkpoint.Digest); err == nil {
		blockHashBytes = decoded
	}

	// Build record
	record := database.EventRecord{
		EventAccountAddress: eventAccountAddress,
		EventHandle:         eventHandle,
		EventOffset:         0, // Synthetic events have offset 0
		TxDigest:            txDigestHex,
		BlockVersion:        0,
		BlockHeight:         strconv.FormatUint(checkpoint.SequenceNumber, 10),
		BlockHash:           blockHashBytes,
		BlockTimestamp:      checkpoint.TimestampMs / 1000, // Convert ms to seconds
		Data:                executionStateChanged,
		IsSynthetic:         true,
	}

	return &record, nil
}

// getTransmitters retrieves the transmitters from the ConfigSet event.
func (tIndexer *TransactionsIndexer) getTransmitters(ctx context.Context) ([]string, error) {
	moduleKey := tIndexer.configEventModuleKey
	eventKey := tIndexer.configEventKey

	eventAccountAddress, _, err := tIndexer.getEventPackageIdFromConfig()
	if err != nil {
		tIndexer.logger.Errorw("Failed to get ConfigSet event config", "error", err)
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
		tIndexer.logger.Errorw("Failed to query ConfigSet events", "error", err)
		return nil, err
	}

	if len(events) == 0 {
		tIndexer.logger.Warnw("No ConfigSet events found")
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
		tIndexer.logger.Warnw("No transmitters found in ConfigSet event")
		return nil, nil
	}

	tIndexer.logger.Infow("Found transmitters in ConfigSet event", "count", len(transmitters))

	return transmitters, nil
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
		return nil, nil
	}

	var configEvent codec.SourceChainConfigSet
	if err := codec.DecodeSuiJsonValue(events[0].Data, &configEvent); err != nil {
		return nil, fmt.Errorf("failed to decode SourceChainConfigSet event: %w", err)
	}

	return &configEvent.SourceChainConfig, nil
}

// getEventPackageIdFromConfig returns the cached OffRamp package.
func (t *TransactionsIndexer) getEventPackageIdFromConfig() (string, string, error) {
	t.mu.RLock()
	pkg := t.offrampPackageID
	latestPkg := t.latestOfframpPackageID
	t.mu.RUnlock()

	if pkg != "" {
		return pkg, latestPkg, nil
	}
	return "", "", fmt.Errorf("offramp package not set yet")
}

// ModuleId represents Move's ModuleId { address, name }
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

// Regex to capture MoveAbort details
var abortRe = regexp.MustCompile(
	`^MoveAbort\(` +
		`MoveLocation \{ module: ModuleId \{ address: ([0-9a-f]+), name: Identifier\("([^"]+)"\) \}, ` +
		`function: (\d+), instruction: (\d+), function_name: (Some\("([^"]+)"\)|None) \}, ` +
		`(\d+)\) in command (\d+)$`,
)

// parseMoveAbortFromExecutionError extracts abort metadata from gRPC v2 ExecutionError.
func (tIndexer *TransactionsIndexer) parseMoveAbortFromExecutionError(execErr *suirpcv2.ExecutionError) (*MoveAbort, error) {
	if execErr == nil {
		return nil, errors.New("execution error is nil")
	}

	abort := execErr.GetAbort()
	if abort == nil {
		return nil, errors.New("execution error has no abort details")
	}

	location := abort.GetLocation()
	if location == nil {
		return nil, errors.New("abort has no location")
	}

	var functionName *string
	if name := location.GetFunctionName(); name != "" {
		functionName = &name
	}

	return &MoveAbort{
		Location: MoveLocation{
			Module: ModuleId{
				Address: location.GetPackage(),
				Name:    location.GetModule(),
			},
			Function:     uint64(location.GetFunction()),
			Instruction:  uint64(location.GetInstruction()),
			FunctionName: functionName,
		},
		AbortCode:    abort.GetAbortCode(),
		CommandIndex: execErr.GetCommand(),
	}, nil
}

// parseMoveAbort parses the legacy Rust debug error string into a MoveAbort struct.
func (tIndexer *TransactionsIndexer) parseMoveAbort(s string) (*MoveAbort, error) {
	m := abortRe.FindStringSubmatch(s)
	if m == nil {
		return nil, fmt.Errorf("input does not match MoveAbort pattern")
	}

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

// Ready returns nil if the indexer has started successfully.
func (tIndexer *TransactionsIndexer) Ready() error {
	return tIndexer.starter.Ready()
}

// Close stops the transactions indexer.
func (tIndexer *TransactionsIndexer) Close() error {
	return tIndexer.starter.StopOnce("TransactionsIndexer", func() error {
		tIndexer.wg.Wait()
		return nil
	})
}
