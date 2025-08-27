# ChainReader

The ChainReader is a core component of the Chainlink SUI Relayer that provides comprehensive blockchain data access capabilities. It enables efficient reading of Sui blockchain state, object data, and events with support for real-time monitoring and historical queries. It consists of 3 main components:

**ContractReader Implementation**: Exposes methods for reading values from contracts and querying events. These methods include:

- `GetLatestValue` - gets the latest value of some state in a bound contract
- `QueryKey` - queries a certain event that was emitted on-chain (and potentially filters by the values of those events)
- Other utility methods

**Events Indexer**: Continuously polls for a set of events defined in the ChainReader's configurations. The events indexer also inserts the found events into the database to ensure that querying them can be done easily with the full extent of SQL querying capabilities, unlike the possibly limited RPC querying features.

> **Note**: It is important to recognize that each database instance is isolated to each individual relayer instance, therefore we must avoid relying completely on the database to answer for those events and should always query the RPC to ensure that the database is caught up before responding to events queries.

**Transactions Indexer**: Finds the transmitters (accounts making on-chain calls to contracts) and watches for failed transactions originating from those accounts. This is useful because in Sui, unlike EVM, events from failed transactions are not indexed and are not findable by querying the RPC's events. Instead, we must generate synthetic events in cases like ExecutionStateChanged in the case of failures.



## Events Indexer Overview

During the initialization of the ChainReader abstraction, the events that we are interested in querying are received as part of the ChainReader's configuration. The ChainReader also receives polling frequency configs (interval and timeout) that will be used as polling constraints in the events indexer.

ChainReader then initializes the events indexer and stores it in the state (self struct). This takes places in the NewChainReader method in /relayer/chainreader/reader/chainreader.go. 

Below is a code snippet that shows how this is done (code is omitted for brevity):

```go
// File: /relayer/chainreader/reader/chainreader.go

func NewChainReader(..., configs config.ChainReaderConfig, ...) (...) {
	//... omitted

    // Create a list of all event selectors to pass to indexers
	eventConfigurations := make([]*client.EventSelector, 0)
	eventConfigurationsMap := make(map[string]*config.ChainReaderEvent)
	for _, moduleConfig := range configs.Modules {
		if moduleConfig.Events != nil {
			for _, eventConfig := range moduleConfig.Events {
				eventConfigurations = append(eventConfigurations, &eventConfig.EventSelector)
				eventConfigurationsMap[fmt.Sprintf("%s::%s", eventConfig.Name, eventConfig.EventType)] = eventConfig
			}
		}
	}

	eventsIndexer := indexer.NewEventIndexer(
		dbStore, // Abstraction over the database connection with helper methods
		lgr,
		abstractClient, // PTB Client (abstraction over the RPC SDK)
		eventConfigurations,
		configs.EventsIndexer.PollingInterval,
		configs.EventsIndexer.SyncTimeout,
	)

	//... omitted

	return &suiChainReader{
		//... omitted
		eventsIndexer:             eventsIndexer,
		eventsIndexerCancel:       nil,
	}, nil
}
```


NOTE: eventConfigurations is a slice of the EventSelector type which in Sui refers to 3 values (package ID, module and type) delimited by ::.  

For example `0x...::offramp::StaticConfigSet.`

```go
// File: /relayer/client/models.go

type EventFilterByMoveEventModule struct {
	Package string `json:"package"`
	Module  string `json:"module"`
	Event   string `json:"event"`
}

// EventSelector is an alias for EventFilterByMoveEventModule
type EventSelector = EventFilterByMoveEventModule
​
```

Once the ChainReader is initialized, a subsequent call to the Start method of ChainReader will also start the events indexer. This can also be found in same file referenced above.

NOTE: we keep track of the cancel method of the context that the events indexer was started with to ensure that a clean stop is achievable later if ChainReader’s Stop method is called.


```go
// File: /relayer/chainreader/reader/chainreader.go

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

		// ... omitted

		return nil
	})
}
```


Once the events indexer starts, it will continuously polls the RPC endpoint for events that are listed within the specified selectors.

```go
// File: /relayer/chainreader/indexer/events_indexer.go

func (eIndexer *EventsIndexer) Start(ctx context.Context) error {
	ticker := time.NewTicker(eIndexer.pollingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			syncCtx, cancel := context.WithTimeout(ctx, eIndexer.syncTimeout)
			start := time.Now()

			err := eIndexer.SyncAllEvents(syncCtx)
			elapsed := time.Since(start)

			if err != nil && !errors.Is(err, context.DeadlineExceeded) {
				eIndexer.logger.Warnw("EventSync completed with errors", "error", err, "duration", elapsed)
			} else if err != nil {
				eIndexer.logger.Warnw("EventSync timed out", "duration", elapsed)
			} else {
				eIndexer.logger.Debugw("Event sync completed successfully", "duration", elapsed)
			}

			cancel()
		case <-ctx.Done():
			eIndexer.logger.Infow("Event polling stopped")
			return nil
		}
	}
}
```

The eIndexer.SyncAllEvents call does the following:

- For each event listed in eventConfigurations
- Fetch the latest cursor (from the DB) for it
- Query the Sui RPC endpoint (using the PTB client) to get the latest events (of that type)
- Insert each into the database

The database events table is created with the following schema:

| Column Name | Data Type | Constraints | Description |
|-------------|-----------|-------------|-------------|
| `id` | `BIGSERIAL` | `PRIMARY KEY` | Auto-incrementing unique identifier |
| `event_account_address` | `TEXT` | `NOT NULL` | Address of the account that emitted the event |
| `event_handle` | `TEXT` | `NOT NULL` | Fully qualified Sui event selector |
| `event_offset` | `BIGINT` | `NOT NULL` | Offset position of the event within the transaction |
| `tx_digest` | `TEXT` | `NOT NULL` | Unique transaction digest/hash |
| `block_version` | `BIGINT` | `NOT NULL` | Version number of the block |
| `block_height` | `TEXT` | `NOT NULL` | Height of the block in the chain |
| `block_hash` | `BYTEA` | `NOT NULL` | Hash of the block (binary data) |
| `block_timestamp` | `BIGINT` | `NOT NULL` | Unix timestamp when the block was created |
| `data` | `JSONB` | `NOT NULL` | Event data as a JSON blob for efficient querying |

**Unique Constraint**: `UNIQUE (event_account_address, event_handle, tx_digest, event_offset)`

> **Note**: The `event_handle` field stores the fully qualified Sui event selector in the format `package::module::event_type`.
​
The event_handle being a string field in the fully qualified Sui event selector format discussed above. Also note that data is simply a JSON blob due to the ability of Postgres to query JSON fields efficiently.

To view the exact fields of each event, you can refer to the corresponding contract or the /relayer/codec/types.go file. All event types will match exactly what is available in Aptos and other implementations as they are cast into strong types in Chainlink Core.

NOTE: the types.go file in the codec module may not include all the types we are indexing since we don’t need to deserialize them in the Relayer but can be easily added to serve as a reference.



In the case of queries for events, QueryKey in ChainReader, we must make update the configuration and sync that specific event before making a call to the database. This helps ensure that responses to queries are always upto date without over reliance on the database records which maybe stale.


One of the challenges faced here is that the ChainReader does not know all the package IDs at compile time (at the time the configuration is set and ChainReader is started). QueryKey method must therefore update its configuration to set the package ID to ensure it’s available for querying.

```go
// File: /relayer/chainreader/reader/chainreader.go

// QueryKey queries events from the indexer database for events that were populated from the RPC node
func (s *suiChainReader) QueryKey(..., contract pkgtypes.BoundContract, filter query.KeyFilter, ...) (...) {
	// ...omitted

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
	}

	// only write contract address, rest will be handled during chainreader config
	eventConfig.EventSelector.Package = contract.Address

	// ...omitted
}
```

## Transactions Indexer Overview

The Transactions Indexer addresses a unique challenge in Sui blockchain: unlike EVM chains, events from failed transactions are not indexed by the RPC and cannot be queried directly. To solve this, the Transactions Indexer monitors transmitter accounts for failed transactions and generates synthetic events that would have been emitted if the transactions had succeeded.

During ChainReader initialization, the Transactions Indexer is created alongside the Events Indexer. The indexer is configured to monitor specific accounts (transmitters) that are responsible for executing cross-chain transactions, particularly looking for failed `ExecutionStateChanged` events.

```go
// File: /relayer/chainreader/reader/chainreader.go

func NewChainReader(..., configs config.ChainReaderConfig, ...) (...) {
    //... omitted

    // Create transactions indexer for synthetic event generation
    transactionIndexer := indexer.NewTransactionsIndexer(
        dataStore.DB(),
        lgr,
        abstractClient,
        configs.TransactionsIndexer.PollingInterval,
        configs.TransactionsIndexer.SyncTimeout,
        eventConfigurationsMap,
    )

    //... omitted

    return &suiChainReader{
        //... omitted
        transactionIndexer:        transactionIndexer,
        transactionIndexerCancel:  nil,
    }, nil
}
```

The Transactions Indexer is initialized with several key parameters:

- **Polling Configuration**: Controls how frequently to check for new transactions
- **Execute Functions**: Monitors specific contract functions (`finish_execute`) that can fail
- **Event Configurations**: Maps of events that need synthetic generation when transactions fail
- **Transmitter Tracking**: Maintains cursors for each monitored transmitter account

### Synthetic Event Generation Process

Once started, the Transactions Indexer follows a systematic process to identify failed transactions and generate synthetic events:

```go
// File: /relayer/chainreader/indexer/transactions_indexer.go

func (tIndexer *TransactionsIndexer) Start(ctx context.Context) error {
    // Wait for initial ExecutionStateChanged event before starting
    if err := tIndexer.waitForInitialEvent(ctx); err != nil {
        return err
    }

    ticker := time.NewTicker(tIndexer.pollingInterval)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            syncCtx, cancel := context.WithTimeout(ctx, tIndexer.syncTimeout)
            err := tIndexer.SyncAllTransmittersTransactions(syncCtx)
            
            // Handle sync results and continue polling
            
        case <-ctx.Done():
            return nil
        }
    }
}
```

The indexer waits for an initial `ExecutionStateChanged` event to ensure the system is properly bootstrapped before monitoring for failures. This prevents the indexer from starting before the contracts are deployed and configured.

### Failed Transaction Detection and Processing

For each polling cycle, the Transactions Indexer performs the following steps:

1. **Query Transmitter Transactions**: Retrieves recent transactions from each known transmitter account
2. **Filter Failed Transactions**: Identifies transactions with `status != "success"`
3. **Validate Transaction Type**: Ensures the failed transaction is a programmable transaction
4. **Parse Error Details**: Extracts Move abort information from the transaction error
5. **Validate Execution Context**: Confirms the failure occurred in the expected module and function

```go
// File: /relayer/chainreader/indexer/transactions_indexer.go

// Process each failed transaction
for _, transactionRecord := range queryResponse.Data {
    if transactionRecord.Effects.Status.Status == "success" {
        continue // Skip successful transactions
    }

    // Parse the Move abort error to understand failure context
    errMessage := transactionRecord.Effects.Status.Error
    moveAbort, err := tIndexer.parseMoveAbort(errMessage)
    if err != nil {
        continue
    }

    // Validate the failure occurred in the expected module/function
    if moveAbort.Location.Module.Name != moduleKey || 
       !slices.Contains(tIndexer.executeFunctions, *moveAbort.Location.FunctionName) {
        continue
    }

    // Extract execution report from transaction arguments
    // Generate synthetic ExecutionStateChanged event
}
```

### Synthetic ExecutionStateChanged Event Creation

When a valid failed execution is detected, the indexer creates a synthetic `ExecutionStateChanged` event that mirrors what would have been emitted if the transaction succeeded but with a failure state:

```go
// Create synthetic ExecutionStateChanged event
executionStateChanged := map[string]any{
    "source_chain_selector": fmt.Sprintf("%d", sourceChainSelector),
    "sequence_number":       fmt.Sprintf("%d", execReport.Message.Header.SequenceNumber),
    "message_id":            "0x" + hex.EncodeToString(execReport.Message.Header.MessageID),
    "message_hash":          "0x" + hex.EncodeToString(messageHash[:]),
    "state":                 uint8(3), // 3 = FAILURE
}
```

The synthetic event contains the same data fields as a real `ExecutionStateChanged` event, but with the execution state set to `FAILURE` (value 3). This enables downstream systems to properly track failed cross-chain message executions.

### Database Integration

Synthetic events are inserted into the same `sui.events` table used by the Events Indexer, ensuring consistent querying capabilities:

```go
record := database.EventRecord{
    EventAccountAddress: eventAccountAddress,
    EventHandle:         eventHandle,        // Same format as real events
    EventOffset:         0,
    TxDigest:            transactionRecord.Digest,
    BlockHeight:         checkpointResponse.SequenceNumber,
    BlockHash:           []byte(checkpointResponse.Digest),
    BlockTimestamp:      blockTimestamp,
    Data:                executionStateChanged, // Synthetic event data
}
```

The indexer employs a batch insertion strategy with individual fallback to ensure maximum reliability when persisting synthetic events.

### Transmitter Discovery and Management  

The Transactions Indexer automatically discovers transmitter accounts by monitoring `ConfigSet` events from the OCR3 base contract. This ensures it stays synchronized with the current set of authorized transmitters without manual configuration:

```go
// File: /relayer/chainreader/indexer/transactions_indexer.go

func (tIndexer *TransactionsIndexer) getTransmitters(ctx context.Context) ([]models.SuiAddress, error) {
    // Query ConfigSet events to find current transmitters
    // Extract transmitter addresses from OCR configuration
    // Update internal transmitter tracking
}
```

This dynamic discovery mechanism ensures the indexer automatically adapts to OCR configuration changes without requiring restarts or manual intervention.

**NOTE**: The Transactions Indexer maintains separate cursors for each transmitter account, enabling efficient incremental processing and avoiding duplicate synthetic event generation.

