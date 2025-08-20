# Database Integration

The Database Integration component provides the persistence layer for the Chainlink SUI Relayer, offering robust PostgreSQL-backed storage for events, transactions, and system state. It includes optimized schemas, indexing strategies, and query patterns designed for high-performance blockchain data access.

## Table of Contents

1. [Overview](#overview)
2. [Database Architecture](#database-architecture)
3. [Schema Design](#schema-design)
4. [Connection Management](#connection-management)
5. [Event Storage](#event-storage)
6. [Query Optimization](#query-optimization)
7. [Indexing Strategies](#indexing-strategies)
8. [Performance Tuning](#performance-tuning)
9. [Backup and Recovery](#backup-and-recovery)
10. [Monitoring](#monitoring)
11. [Troubleshooting](#troubleshooting)

## Overview

The database integration layer provides:

- **Event Storage**: Structured storage for blockchain events with JSONB data fields
- **Transaction Tracking**: Persistent transaction state and metadata
- **Cursor Management**: Position tracking for event indexing resumption
- **Query Optimization**: Efficient indexes and query patterns for fast data access
- **Connection Pooling**: Optimized database connection management
- **Schema Evolution**: Support for database migrations and schema updates

```mermaid
graph TB
    subgraph "Application Layer"
        CR[ChainReader]
        EI[Event Indexer]
        TXM[Transaction Manager]
    end
    
    subgraph "Database Layer"
        DBS[DB Store]
        CP[Connection Pool]
        QC[Query Cache]
    end
    
    subgraph "PostgreSQL"
        ET[Events Table]
        TT[Transactions Table]
        IDX[Indexes]
        JSONB[JSONB Fields]
    end
    
    CR --> DBS
    EI --> DBS
    TXM --> DBS
    
    DBS --> CP
    DBS --> QC
    
    CP --> ET
    CP --> TT
    CP --> IDX
    ET --> JSONB
```

## Database Architecture

### Core Components

The database integration consists of several key components:

#### 1. DBStore
The main database interface providing CRUD operations:

```go
type DBStore struct {
    conn   sqlutil.DataSource
    logger logger.Logger
}

func NewDBStore(conn sqlutil.DataSource, logger logger.Logger) *DBStore {
    return &DBStore{
        conn:   conn,
        logger: logger,
    }
}
```

#### 2. Connection Pool Management
Optimized connection pooling for high concurrency:

```go
// Connection pool configuration
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(10)
db.SetConnMaxLifetime(5 * time.Minute)
db.SetConnMaxIdleTime(30 * time.Second)
```

#### 3. Schema Management
Automated schema creation and migration:

```go
func (db *DBStore) EnsureSchema(ctx context.Context) error {
    if err := db.createEventsTable(ctx); err != nil {
        return fmt.Errorf("failed to create events table: %w", err)
    }
    
    if err := db.createIndexes(ctx); err != nil {
        return fmt.Errorf("failed to create indexes: %w", err)
    }
    
    return nil
}
```

### Data Flow

1. **Event Ingestion**: Events from blockchain are batch-inserted into events table
2. **Index Updates**: Automatic index maintenance for optimal query performance
3. **Query Processing**: Optimized queries using strategic indexes
4. **Connection Pooling**: Efficient connection reuse across operations

## Schema Design

### Events Table

The primary table for storing blockchain events:

```sql
CREATE TABLE events (
    id BIGSERIAL PRIMARY KEY,
    event_account_address TEXT NOT NULL,
    event_handle TEXT NOT NULL,
    event_offset BIGINT NOT NULL,
    tx_digest TEXT NOT NULL,
    block_version BIGINT NOT NULL,
    block_height TEXT NOT NULL,
    block_hash BYTEA NOT NULL,
    block_timestamp BIGINT NOT NULL,
    data JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
```

#### Field Descriptions

| Field | Type | Purpose |
|-------|------|---------|
| `id` | BIGSERIAL | Auto-incrementing primary key |
| `event_account_address` | TEXT | Contract package address |
| `event_handle` | TEXT | Full event identifier (package::module::event) |
| `event_offset` | BIGINT | Event sequence number for ordering |
| `tx_digest` | TEXT | Transaction hash containing the event |
| `block_version` | BIGINT | Block version number |
| `block_height` | TEXT | Block height as string |
| `block_hash` | BYTEA | Block hash as bytes |
| `block_timestamp` | BIGINT | Unix timestamp of block |
| `data` | JSONB | Event data as structured JSON |
| `created_at` | TIMESTAMP | Record creation timestamp |

### Event Data Structure

Events store structured data in JSONB format:

```json
{
  "counter_id": "0x123...",
  "old_value": "42",
  "new_value": "43",
  "user": "0xabc...",
  "amount": "1000000",
  "timestamp": 1640995200
}
```

### Transactions Table (Optional)

For transaction state tracking:

```sql
CREATE TABLE transactions (
    id TEXT PRIMARY KEY,
    sender TEXT NOT NULL,
    payload TEXT NOT NULL,
    functions TEXT[] NOT NULL,
    signatures TEXT[] NOT NULL,
    request_type TEXT NOT NULL,
    attempt INTEGER NOT NULL DEFAULT 1,
    state INTEGER NOT NULL DEFAULT 0,
    digest TEXT,
    last_updated_at BIGINT NOT NULL,
    tx_error TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);
```

## Connection Management

### Database Configuration

Configure PostgreSQL connection:

```bash
# Environment variable
export DATABASE_URL="postgresql://username:password@localhost:5432/chainlink"

# Or in Go
config := &DatabaseConfig{
    URL:                "postgresql://username:password@localhost:5432/chainlink",
    MaxOpenConns:       25,
    MaxIdleConns:       10,
    ConnMaxLifetime:    5 * time.Minute,
    ConnMaxIdleTime:    30 * time.Second,
}
```

### Connection Pool Optimization

```go
func ConfigureConnectionPool(db *sql.DB) {
    // Maximum number of open connections
    db.SetMaxOpenConns(25)
    
    // Maximum number of idle connections
    db.SetMaxIdleConns(10)
    
    // Maximum lifetime of connections
    db.SetConnMaxLifetime(5 * time.Minute)
    
    // Maximum idle time for connections
    db.SetConnMaxIdleTime(30 * time.Second)
}
```

### Connection Health Monitoring

```go
func (db *DBStore) HealthCheck(ctx context.Context) error {
    // Test basic connectivity
    if err := db.conn.PingContext(ctx); err != nil {
        return fmt.Errorf("database ping failed: %w", err)
    }
    
    // Test with simple query
    var result int
    err := db.conn.QueryRowContext(ctx, "SELECT 1").Scan(&result)
    if err != nil {
        return fmt.Errorf("database query test failed: %w", err)
    }
    
    return nil
}
```

## Event Storage

### Batch Insert Operations

Optimized batch insertion for high-throughput event ingestion:

```go
func (db *DBStore) InsertEvents(ctx context.Context, records []EventRecord) error {
    if len(records) == 0 {
        return nil
    }

    // Use transaction for batch operation
    tx, err := db.conn.BeginTx(ctx, nil)
    if err != nil {
        return fmt.Errorf("failed to begin transaction: %w", err)
    }
    defer tx.Rollback()

    // Prepare statement
    stmt, err := tx.PrepareContext(ctx, `
        INSERT INTO events (
            event_account_address, event_handle, event_offset,
            tx_digest, block_version, block_height, block_hash,
            block_timestamp, data
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `)
    if err != nil {
        return fmt.Errorf("failed to prepare statement: %w", err)
    }
    defer stmt.Close()

    // Execute batch inserts
    for _, record := range records {
        _, err := stmt.ExecContext(ctx,
            record.EventAccountAddress,
            record.EventHandle,
            record.EventOffset,
            record.TxDigest,
            record.BlockVersion,
            record.BlockHeight,
            record.BlockHash,
            record.BlockTimestamp,
            record.Data,
        )
        if err != nil {
            return fmt.Errorf("failed to insert event: %w", err)
        }
    }

    // Commit transaction
    if err := tx.Commit(); err != nil {
        return fmt.Errorf("failed to commit transaction: %w", err)
    }

    db.logger.Debugw("Inserted batch of events", "count", len(records))
    return nil
}
```

### COPY-based Bulk Loading

For maximum performance with large datasets:

```go
func (db *DBStore) BulkInsertEvents(ctx context.Context, records []EventRecord) error {
    // Use COPY for bulk loading
    conn, err := db.conn.Conn(ctx)
    if err != nil {
        return err
    }
    defer conn.Close()

    // Execute COPY command
    _, err = conn.ExecContext(ctx, "COPY events (event_account_address, event_handle, event_offset, tx_digest, block_version, block_height, block_hash, block_timestamp, data) FROM STDIN")
    if err != nil {
        return err
    }

    // Stream data
    for _, record := range records {
        line := fmt.Sprintf("%s\t%s\t%d\t%s\t%d\t%s\t%x\t%d\t%s\n",
            record.EventAccountAddress,
            record.EventHandle,
            record.EventOffset,
            record.TxDigest, 
            record.BlockVersion,
            record.BlockHeight,
            record.BlockHash,
            record.BlockTimestamp,
            string(record.Data),
        )
        
        if _, err := conn.Write([]byte(line)); err != nil {
            return err
        }
    }
    
    return nil
}
```

## Query Optimization

### Event Queries

Optimized queries for common access patterns:

```go
func (db *DBStore) QueryEvents(ctx context.Context, packageAddress, eventHandle string, cursor *EventCursor, limitAndSort query.LimitAndSort) ([]EventRecord, error) {
    var conditions []string
    var args []interface{}
    argIndex := 1

    // Base conditions
    conditions = append(conditions, fmt.Sprintf("event_account_address = $%d", argIndex))
    args = append(args, packageAddress)
    argIndex++

    conditions = append(conditions, fmt.Sprintf("event_handle = $%d", argIndex))
    args = append(args, eventHandle)
    argIndex++

    // Cursor-based pagination
    if cursor != nil {
        conditions = append(conditions, fmt.Sprintf("block_timestamp > $%d", argIndex))
        args = append(args, cursor.Timestamp)
        argIndex++
    }

    // Build query
    whereClause := strings.Join(conditions, " AND ")
    orderBy := "block_timestamp DESC"
    if limitAndSort.Sort.SortType == query.Asc {
        orderBy = "block_timestamp ASC"
    }

    query := fmt.Sprintf(`
        SELECT event_account_address, event_handle, event_offset, tx_digest,
               block_version, block_height, block_hash, block_timestamp, data
        FROM events
        WHERE %s
        ORDER BY %s
        LIMIT $%d
    `, whereClause, orderBy, argIndex)

    args = append(args, limitAndSort.Limit.Count)

    // Execute query
    rows, err := db.conn.QueryContext(ctx, query, args...)
    if err != nil {
        return nil, fmt.Errorf("failed to query events: %w", err)
    }
    defer rows.Close()

    // Parse results
    var events []EventRecord
    for rows.Next() {
        var event EventRecord
        err := rows.Scan(
            &event.EventAccountAddress,
            &event.EventHandle,
            &event.EventOffset,
            &event.TxDigest,
            &event.BlockVersion,
            &event.BlockHeight,
            &event.BlockHash,
            &event.BlockTimestamp,
            &event.Data,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to scan event: %w", err)
        }
        events = append(events, event)
    }

    return events, nil
}
```

### JSONB Queries

Query event data using JSONB operators:

```go
func (db *DBStore) QueryEventsByValue(ctx context.Context, eventHandle, jsonPath string, value interface{}) ([]EventRecord, error) {
    query := `
        SELECT event_account_address, event_handle, event_offset, tx_digest,
               block_version, block_height, block_hash, block_timestamp, data
        FROM events
        WHERE event_handle = $1 AND data->>$2 = $3
        ORDER BY block_timestamp DESC
        LIMIT 100
    `

    rows, err := db.conn.QueryContext(ctx, query, eventHandle, jsonPath, value)
    if err != nil {
        return nil, fmt.Errorf("failed to query events by value: %w", err)
    }
    defer rows.Close()

    return scanEventRows(rows)
}
```

### Complex JSONB Queries

Advanced querying with JSONB:

```sql
-- Query events with numeric comparisons
SELECT * FROM events 
WHERE event_handle = 'package::module::Transfer'
AND (data->>'amount')::bigint > 1000000;

-- Query events with array containment
SELECT * FROM events
WHERE data->'recipients' ? '0x123...';

-- Query events with nested object access
SELECT * FROM events
WHERE data->'metadata'->>'type' = 'ccip_message';

-- Query events with existence checks
SELECT * FROM events
WHERE data ? 'error_code'
AND event_handle LIKE '%Error%';
```

## Indexing Strategies

### Primary Indexes

Essential indexes for query performance:

```sql
-- Primary key index (automatic)
CREATE UNIQUE INDEX events_pkey ON events(id);

-- Handle + timestamp for time-based queries
CREATE INDEX idx_events_handle_timestamp ON events(event_handle, block_timestamp DESC);

-- Offset tracking for cursor management
CREATE INDEX idx_events_offset ON events(event_account_address, event_handle, event_offset);

-- Transaction lookup
CREATE INDEX idx_events_tx_digest ON events(tx_digest);
```

### JSONB Indexes

Indexes for efficient JSONB queries:

```sql
-- GIN index for general JSONB queries
CREATE INDEX idx_events_data_gin ON events USING gin(data);

-- Specific path indexes for common queries
CREATE INDEX idx_events_amount ON events((data->>'amount'));
CREATE INDEX idx_events_user ON events((data->>'user'));
CREATE INDEX idx_events_counter_id ON events((data->>'counter_id'));

-- Partial indexes for specific event types
CREATE INDEX idx_transfer_events ON events(block_timestamp DESC)
WHERE event_handle LIKE '%Transfer%';
```

### Composite Indexes

Multi-column indexes for complex queries:

```sql
-- Common query pattern: handle + timestamp + account
CREATE INDEX idx_events_composite ON events(
    event_handle, 
    block_timestamp DESC, 
    event_account_address
);

-- Time range queries with specific data
CREATE INDEX idx_events_time_data ON events(
    block_timestamp DESC,
    (data->>'type')
) WHERE data ? 'type';
```

### Index Monitoring

Monitor index usage and performance:

```sql
-- Check index usage statistics
SELECT 
    schemaname,
    tablename,
    indexname,
    idx_scan,
    idx_tup_read,
    idx_tup_fetch
FROM pg_stat_user_indexes
WHERE tablename = 'events'
ORDER BY idx_scan DESC;

-- Identify unused indexes
SELECT 
    schemaname,
    tablename,
    indexname,
    idx_scan
FROM pg_stat_user_indexes
WHERE idx_scan = 0 
AND tablename = 'events';

-- Check index sizes
SELECT 
    indexname,
    pg_size_pretty(pg_relation_size(indexname::regclass)) as size
FROM pg_indexes
WHERE tablename = 'events';
```

## Performance Tuning

### PostgreSQL Configuration

Optimize PostgreSQL settings for blockchain workloads:

```ini
# postgresql.conf

# Memory settings
shared_buffers = 256MB                  # 25% of system RAM
effective_cache_size = 1GB              # 75% of system RAM
maintenance_work_mem = 64MB
work_mem = 4MB

# Checkpoint settings
checkpoint_completion_target = 0.9
wal_buffers = 16MB
max_wal_size = 1GB
min_wal_size = 80MB

# Connection settings
max_connections = 100
shared_preload_libraries = 'pg_stat_statements'

# Logging
log_min_duration_statement = 1000       # Log slow queries
log_line_prefix = '%t [%p]: [%l-1] user=%u,db=%d,app=%a,client=%h '
log_checkpoints = on
log_connections = on
log_disconnections = on
log_lock_waits = on
```

### Query Performance

Optimize common query patterns:

```go
// Use prepared statements for repeated queries
type PreparedQueries struct {
    getEventsByHandle  *sql.Stmt
    getEventsByTxHash  *sql.Stmt
    getLatestEvents    *sql.Stmt
}

func (db *DBStore) PrepareQueries(ctx context.Context) (*PreparedQueries, error) {
    pq := &PreparedQueries{}
    
    var err error
    pq.getEventsByHandle, err = db.conn.PrepareContext(ctx, `
        SELECT * FROM events 
        WHERE event_handle = $1 
        ORDER BY block_timestamp DESC 
        LIMIT $2
    `)
    if err != nil {
        return nil, err
    }
    
    pq.getEventsByTxHash, err = db.conn.PrepareContext(ctx, `
        SELECT * FROM events 
        WHERE tx_digest = $1
    `)
    if err != nil {
        return nil, err
    }
    
    return pq, nil
}
```

### Batch Operations

Optimize batch operations:

```go
func (db *DBStore) BatchInsertWithCopy(ctx context.Context, records []EventRecord) error {
    // Use COPY for large batches (>1000 records)
    if len(records) > 1000 {
        return db.copyInsert(ctx, records)
    }
    
    // Use prepared statements for smaller batches
    return db.batchInsert(ctx, records)
}

func (db *DBStore) copyInsert(ctx context.Context, records []EventRecord) error {
    // Implementation using COPY command
    conn, err := db.conn.Conn(ctx)
    if err != nil {
        return err
    }
    defer conn.Close()
    
    // Stream data using COPY
    // ... implementation details
    
    return nil
}
```

### Connection Pool Tuning

Optimize connection pool settings:

```go
func TuneConnectionPool(db *sql.DB, maxConcurrency int) {
    // Rule of thumb: MaxOpenConns = (core_count * 2) - 1
    maxOpen := maxConcurrency * 2
    if maxOpen > 50 {
        maxOpen = 50 // Cap at reasonable limit
    }
    
    db.SetMaxOpenConns(maxOpen)
    db.SetMaxIdleConns(maxOpen / 2)
    db.SetConnMaxLifetime(5 * time.Minute)
    db.SetConnMaxIdleTime(30 * time.Second)
}
```

## Backup and Recovery

### Automated Backups

Set up automated backup strategy:

```bash
#!/bin/bash
# backup.sh

BACKUP_DIR="/var/backups/postgresql"
DB_NAME="chainlink"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)

# Create full backup
pg_dump -h localhost -U postgres -d $DB_NAME \
    --format=custom \
    --compress=9 \
    --file="$BACKUP_DIR/chainlink_$TIMESTAMP.dump"

# Create schema-only backup
pg_dump -h localhost -U postgres -d $DB_NAME \
    --schema-only \
    --file="$BACKUP_DIR/chainlink_schema_$TIMESTAMP.sql"

# Cleanup old backups (keep 7 days)
find $BACKUP_DIR -name "chainlink_*.dump" -mtime +7 -delete
```

### Point-in-Time Recovery

Configure WAL archiving for point-in-time recovery:

```ini
# postgresql.conf
archive_mode = on
archive_command = 'cp %p /archive/%f'
max_wal_senders = 3
wal_level = replica
```

### Recovery Procedures

```bash
# Restore from backup
pg_restore -h localhost -U postgres -d chainlink_restored \
    --clean --create --verbose \
    /var/backups/postgresql/chainlink_20231215_120000.dump

# Point-in-time recovery
pg_basebackup -h localhost -U postgres -D /tmp/recovery \
    --checkpoint=fast --write-recovery-conf
```

## Monitoring

### Database Metrics

Monitor key database performance metrics:

```go
type DatabaseMetrics struct {
    ConnectionsActive     prometheus.Gauge
    ConnectionsIdle       prometheus.Gauge
    QueriesPerSecond     prometheus.Counter
    QueryDuration        prometheus.Histogram
    TableSize            prometheus.Gauge
    IndexUsage           prometheus.GaugeVec
}

func (db *DBStore) CollectMetrics() *DatabaseMetrics {
    metrics := &DatabaseMetrics{
        ConnectionsActive: prometheus.NewGauge(prometheus.GaugeOpts{
            Name: "db_connections_active",
            Help: "Number of active database connections",
        }),
        QueriesPerSecond: prometheus.NewCounter(prometheus.CounterOpts{
            Name: "db_queries_total",
            Help: "Total number of database queries",
        }),
        QueryDuration: prometheus.NewHistogram(prometheus.HistogramOpts{
            Name: "db_query_duration_seconds",
            Help: "Database query duration",
        }),
    }
    
    return metrics
}
```

### Performance Monitoring Queries

```sql
-- Active connections
SELECT count(*) as active_connections
FROM pg_stat_activity
WHERE state = 'active';

-- Slow queries
SELECT query, mean_time, calls, total_time
FROM pg_stat_statements
ORDER BY mean_time DESC
LIMIT 10;

-- Table sizes
SELECT 
    tablename,
    pg_size_pretty(pg_total_relation_size(tablename::regclass)) as size
FROM pg_tables
WHERE schemaname = 'public';

-- Index usage
SELECT 
    indexname,
    idx_scan,
    idx_tup_read,
    idx_tup_fetch
FROM pg_stat_user_indexes
WHERE tablename = 'events'
ORDER BY idx_scan DESC;
```

### Alerting Rules

Set up alerts for critical database issues:

```go
func (db *DBStore) CheckHealth() []Alert {
    var alerts []Alert
    
    // Check connection count
    if activeConns := db.getActiveConnections(); activeConns > 80 {
        alerts = append(alerts, Alert{
            Type:    "HIGH_CONNECTION_COUNT",
            Message: fmt.Sprintf("Active connections: %d", activeConns),
        })
    }
    
    // Check table size growth
    if tableSize := db.getTableSize("events"); tableSize > 100*1024*1024*1024 { // 100GB
        alerts = append(alerts, Alert{
            Type:    "LARGE_TABLE_SIZE",
            Message: fmt.Sprintf("Events table size: %d bytes", tableSize),
        })
    }
    
    // Check slow queries
    if slowQueries := db.getSlowQueryCount(); slowQueries > 10 {
        alerts = append(alerts, Alert{
            Type:    "SLOW_QUERIES",
            Message: fmt.Sprintf("Slow queries in last hour: %d", slowQueries),
        })
    }
    
    return alerts
}
```

## Troubleshooting

### Common Issues

| Issue | Symptoms | Solution |
|-------|----------|----------|
| **Connection Pool Exhaustion** | "too many connections" errors | Increase pool size, check for connection leaks |
| **Slow Queries** | High response times | Add indexes, optimize queries, check EXPLAIN plans |
| **Disk Space** | Insert failures, maintenance errors | Archive old data, increase disk space |
| **Lock Contention** | Query timeouts, blocking | Optimize transaction scope, add indexes |
| **Index Bloat** | Degraded performance | REINDEX or VACUUM operations |

### Diagnostic Queries

```sql
-- Find blocking queries
SELECT 
    blocked_locks.pid AS blocked_pid,
    blocked_activity.usename AS blocked_user,
    blocking_locks.pid AS blocking_pid,
    blocking_activity.usename AS blocking_user,
    blocked_activity.query AS blocked_statement,
    blocking_activity.query AS current_statement_in_blocking_process
FROM pg_catalog.pg_locks blocked_locks
JOIN pg_catalog.pg_activity blocked_activity ON blocked_activity.pid = blocked_locks.pid
JOIN pg_catalog.pg_locks blocking_locks ON blocking_locks.locktype = blocked_locks.locktype
    AND blocking_locks.database IS NOT DISTINCT FROM blocked_locks.database
    AND blocking_locks.relation IS NOT DISTINCT FROM blocked_locks.relation
    AND blocking_locks.page IS NOT DISTINCT FROM blocked_locks.page
    AND blocking_locks.tuple IS NOT DISTINCT FROM blocked_locks.tuple
    AND blocking_locks.virtualxid IS NOT DISTINCT FROM blocked_locks.virtualxid
    AND blocking_locks.transactionid IS NOT DISTINCT FROM blocked_locks.transactionid
    AND blocking_locks.classid IS NOT DISTINCT FROM blocked_locks.classid
    AND blocking_locks.objid IS NOT DISTINCT FROM blocked_locks.objid
    AND blocking_locks.objsubid IS NOT DISTINCT FROM blocked_locks.objsubid
    AND blocking_locks.pid != blocked_locks.pid
JOIN pg_catalog.pg_activity blocking_activity ON blocking_activity.pid = blocking_locks.pid
WHERE NOT blocked_locks.granted;

-- Check table and index sizes
SELECT 
    tablename,
    pg_size_pretty(pg_relation_size(tablename::regclass)) as table_size,
    pg_size_pretty(pg_total_relation_size(tablename::regclass)) as total_size
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(tablename::regclass) DESC;

-- Analyze query performance
EXPLAIN (ANALYZE, BUFFERS) SELECT * FROM events 
WHERE event_handle = 'package::module::Event' 
AND block_timestamp > extract(epoch from now() - interval '1 hour')
ORDER BY block_timestamp DESC;
```

### Maintenance Operations

```sql
-- Regular maintenance tasks
VACUUM ANALYZE events;
REINDEX TABLE events;

-- Check for index bloat
SELECT 
    schemaname,
    tablename,
    indexname,
    pg_size_pretty(pg_relation_size(indexname::regclass)) as index_size,
    idx_scan,
    idx_tup_read
FROM pg_stat_user_indexes
JOIN pg_indexes USING (schemaname, tablename, indexname)
WHERE tablename = 'events'
ORDER BY pg_relation_size(indexname::regclass) DESC;

-- Update statistics
ANALYZE events;

-- Check table bloat
SELECT 
    schemaname,
    tablename,
    n_tup_ins,
    n_tup_upd,
    n_tup_del,
    n_dead_tup
FROM pg_stat_user_tables
WHERE tablename = 'events';
```

The Database Integration layer provides a robust, scalable foundation for the Chainlink SUI Relayer's data persistence needs. By following the schema design, indexing strategies, and performance optimization techniques outlined in this guide, you can ensure efficient and reliable blockchain data storage and retrieval. 