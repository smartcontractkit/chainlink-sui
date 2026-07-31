package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
)

const (
	MaxEventsQueryLimit = 2000
)

type DBStore struct {
	ds  sqlutil.DataSource
	lgr logger.Logger
}

func NewDBStore(ds sqlutil.DataSource, lgr logger.Logger) *DBStore {
	return &DBStore{
		ds:  ds,
		lgr: logger.Named(lgr, "SuiDBStore"),
	}
}

func (store *DBStore) EnsureSchema(ctx context.Context) error {
	_, err := store.ds.ExecContext(ctx, CreateSchema)
	if err != nil {
		return fmt.Errorf("failed to create sui schema: %w", err)
	}

	_, err = store.ds.ExecContext(ctx, CreateEventsTable)
	if err != nil {
		return fmt.Errorf("failed to create sui.events table: %w", err)
	}

	_, err = store.ds.ExecContext(ctx, CreateIndices)
	if err != nil {
		return fmt.Errorf("failed to create sui indexes: %w", err)
	}

	_, err = store.ds.ExecContext(ctx, CreateCheckpointCursorsTable)
	if err != nil {
		return fmt.Errorf("failed to create sui.checkpoint_cursors table: %w", err)
	}

	return nil
}

type EventRecord struct {
	EventAccountAddress string
	EventHandle         string
	// TxIndex is the position of the transaction within its checkpoint, as returned by the node.
	// Together with (block_height, event_offset) it provides a total order across all events.
	TxIndex        uint64
	EventOffset    uint64
	TxDigest       string
	BlockVersion   uint64
	BlockHeight    string
	BlockHash      []byte
	BlockTimestamp uint64
	Data           map[string]any
	IsSynthetic    bool
}

func (store *DBStore) InsertEvents(ctx context.Context, records []EventRecord) error {
	if len(records) == 0 {
		return nil
	}

	for _, record := range records {
		data, err := json.Marshal(record.Data)
		if err != nil {
			return fmt.Errorf("failed to marshal event data for handle %s at offset %d: %w", record.EventHandle, record.EventOffset, err)
		}

		_, err = store.ds.ExecContext(ctx, InsertEvent,
			record.EventAccountAddress,
			record.EventHandle,
			record.EventOffset,
			record.TxDigest,
			record.BlockVersion,
			record.BlockHeight,
			record.BlockHash,
			record.BlockTimestamp,
			data,
			record.IsSynthetic,
			record.TxIndex,
		)
		if err != nil {
			return fmt.Errorf("failed to insert event (handle: %s, offset: %d): %w", record.EventHandle, record.EventOffset, err)
		}
	}

	return nil
}

func (store *DBStore) QueryEvents(ctx context.Context, eventAccountAddress, eventHandle string, expressions []query.Expression, limitAndSort query.LimitAndSort) ([]EventRecord, error) {
	baseSQL := QueryEventsBase

	args := []any{eventAccountAddress, eventHandle}
	argCount := 3

	if len(expressions) > 0 {
		var conditions []string
		for _, expr := range expressions {
			sqlCondition, err := BuildSQLCondition(expr, &args, &argCount)
			if err != nil {
				return nil, fmt.Errorf("failed to build SQL condition: %w", err)
			}
			conditions = append(conditions, sqlCondition)
		}

		if len(conditions) > 0 {
			baseSQL += " AND " + strings.Join(conditions, " AND ")
		}
	}

	direction := "DESC" // default to descending order if no sort is provided
	if len(limitAndSort.SortBy) > 0 {
		direction = "ASC"
		if sortDir, ok := limitAndSort.SortBy[0].(query.SortBySequence); ok && sortDir.GetDirection() == query.Desc {
			direction = "DESC"
		}
	}
	// (block_height, tx_idx, event_offset) totally orders events as returned by the node:
	// checkpoint, then transaction position within the checkpoint, then event position within the tx.
	baseSQL += fmt.Sprintf(" ORDER BY block_height::BIGINT %s, tx_idx %s, event_offset %s", direction, direction, direction)

	limit := limitAndSort.Limit.Count
	if limit > MaxEventsQueryLimit {
		store.lgr.Warnw("query limit is greater than max limit, using max limit", "limit", limit, "maxLimit", MaxEventsQueryLimit)
		limit = MaxEventsQueryLimit
	} else if limit == 0 {
		// use max limit if no limit is provided
		limit = MaxEventsQueryLimit
	}
	baseSQL += fmt.Sprintf(" LIMIT %d", limit)

	rows, err := store.ds.QueryContext(ctx, baseSQL, args...)
	if err != nil {
		return nil, fmt.Errorf("query events failed: %w", err)
	}
	defer rows.Close()

	var records []EventRecord
	for rows.Next() {
		var record EventRecord
		var dataBytes []byte
		err := rows.Scan(&record.EventAccountAddress, &record.EventHandle, &record.TxIndex, &record.EventOffset, &record.BlockVersion, &record.BlockHeight, &record.BlockHash, &record.BlockTimestamp, &record.TxDigest, &dataBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event record: %w", err)
		}

		var data map[string]any
		if err := json.Unmarshal(dataBytes, &data); err != nil {
			return nil, fmt.Errorf("failed to unmarshal event data: %w", err)
		}
		record.Data = data
		records = append(records, record)
	}

	return records, nil
}

// GetCheckpointCursor returns the persisted last-processed checkpoint sequence for the given
// cursor id. found is false when no cursor has been persisted yet.
func (store *DBStore) GetCheckpointCursor(ctx context.Context, id string) (seq uint64, found bool, err error) {
	err = store.ds.QueryRowxContext(ctx, GetCheckpointCursorQuery, id).Scan(&seq)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, false, nil
		}
		return 0, false, fmt.Errorf("failed to get checkpoint cursor %q: %w", id, err)
	}

	return seq, true, nil
}

// UpsertCheckpointCursor persists the last-processed checkpoint sequence for the given cursor id.
// The stored value is monotonic: it never regresses to a lower sequence.
func (store *DBStore) UpsertCheckpointCursor(ctx context.Context, id string, seq uint64) error {
	_, err := store.ds.ExecContext(ctx, UpsertCheckpointCursorQuery, id, seq)
	if err != nil {
		return fmt.Errorf("failed to upsert checkpoint cursor %q: %w", id, err)
	}

	return nil
}

func operatorSQL(op primitives.ComparisonOperator) string {
	switch op {
	case primitives.Eq:
		return "="
	case primitives.Neq:
		return "!="
	case primitives.Gt:
		return ">"
	case primitives.Gte:
		return ">="
	case primitives.Lt:
		return "<"
	case primitives.Lte:
		return "<="
	default:
		// Default to equality if unknown
		return "="
	}
}
