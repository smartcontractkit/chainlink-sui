package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/block-vision/sui-go-sdk/models"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
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

	return nil
}

type EventRecord struct {
	EventAccountAddress string
	EventHandle         string
	EventOffset         uint64
	TxDigest            string
	BlockVersion        uint64
	BlockHeight         string
	BlockHash           []byte
	BlockTimestamp      uint64
	Data                map[string]any
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

	if len(limitAndSort.SortBy) > 0 {
		direction := "ASC"
		if sortDir, ok := limitAndSort.SortBy[0].(query.SortBySequence); ok && sortDir.GetDirection() == query.Desc {
			direction = "DESC"
		}
		baseSQL += " ORDER BY event_offset " + direction
	} else {
		// default to descending order if no sort is provided
		baseSQL += " ORDER BY event_offset ASC"
	}

	if limitAndSort.Limit.Count > 0 {
		baseSQL += fmt.Sprintf(" LIMIT %d", limitAndSort.Limit.Count)
	}

	store.lgr.Debugw("querying events", "sql", baseSQL, "args", args)

	rows, err := store.ds.QueryContext(ctx, baseSQL, args...)
	if err != nil {
		return nil, fmt.Errorf("query events failed: %w", err)
	}
	defer rows.Close()

	var records []EventRecord
	for rows.Next() {
		var record EventRecord
		var dataBytes []byte
		err := rows.Scan(&record.EventAccountAddress, &record.EventHandle, &record.EventOffset, &record.BlockVersion, &record.BlockHeight, &record.BlockHash, &record.BlockTimestamp, &record.TxDigest, &dataBytes)
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

	store.lgr.Debugw("fetched DB events", "records", records)

	return records, nil
}

// GetLatestOffset returns a cursor (of type EventId) based on the latest event recorded in the DB for a given type
func (store *DBStore) GetLatestOffset(ctx context.Context, eventAccountAddress, eventHandle string) (*models.EventId, uint64, error) {
	var offset uint64
	var txDigest string
	var totalCount uint64
	err := store.ds.QueryRowxContext(ctx, QueryEventsOffset, eventAccountAddress, eventHandle).Scan(&offset, &txDigest, &totalCount)
	if err != nil {
		store.lgr.Errorw("failed to get latest offset", "error", err, "eventAccountAddress", eventAccountAddress, "eventHandle", eventHandle)
		// no rows found in DB, return a nil index
		//nolint:nilnil
		if errors.Is(err, sql.ErrNoRows) {
			// this is not an error, just nothing to return
			return nil, 0, nil
		}

		return nil, 0, fmt.Errorf("failed to get latest offset: %w", err)
	}

	store.lgr.Debugw("latest offset", "offset", offset, "txDigest", txDigest)

	return &models.EventId{
		TxDigest: txDigest,
		// EventSeq is scoped per transaction, the first (and only) event in a tx always has eventSeq = "0".
		// We use (txDigest, eventSeq) as the pagination cursor to resume fetching events reliably.
		EventSeq: "0",
	}, totalCount, nil
}

func (store *DBStore) GetTxDigestByEventId(ctx context.Context, eventID uint64) (string, error) {
	var txDigest string
	err := store.ds.QueryRowxContext(ctx, GetTxDigestById, eventID).Scan(&txDigest)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("no transaction found for event ID %d: %w", eventID, err)
		}
		return "", fmt.Errorf("failed to get transaction digest by event ID %d: %w", eventID, err)
	}
	return txDigest, nil
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
