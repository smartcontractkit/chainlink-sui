package database

import (
	"fmt"
	"strings"
)

const eventsTableName = "events"

func quoteIdentifier(name string) string {
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}

func buildQueries(schema string) *queries {
	schemaIdent := quoteIdentifier(schema)
	eventsTable := schemaIdent + "." + quoteIdentifier(eventsTableName)

	return &queries{
		createSchema: fmt.Sprintf(`CREATE SCHEMA IF NOT EXISTS %s;`, schemaIdent),
		createEventsTable: fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
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
			UNIQUE (event_account_address, event_handle, tx_digest, event_offset)
		);
	    `, eventsTable),
		insertEvent: fmt.Sprintf(`
		INSERT INTO %s (
			event_account_address,
			event_handle,
			event_offset,
		    tx_digest,
			block_version,
			block_height,
			block_hash,
			block_timestamp,
			data
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT DO NOTHING;
	    `, eventsTable),
		queryEventsBase: fmt.Sprintf(`
		SELECT event_account_address, event_handle, event_offset, block_version, block_height, block_hash, block_timestamp, tx_digest, data
		FROM %s
		WHERE event_account_address = $1 AND event_handle = $2
	    `, eventsTable),
		queryEventsOffset: fmt.Sprintf(`
		SELECT COALESCE(event_offset, 0) as event_offset, tx_digest, COUNT(*) OVER() as total_count
		FROM %s
		WHERE event_account_address = $1 AND event_handle = $2
		ORDER BY id DESC
		LIMIT 1
		`, eventsTable),
		countEvents: fmt.Sprintf(`
		SELECT COUNT(*) as total_count
		FROM %s
		WHERE event_account_address = $1 AND event_handle = $2
		`, eventsTable),
		getTxDigestById: fmt.Sprintf(`
		SELECT tx_digest
		FROM %s
		WHERE id = $1
		`, eventsTable),
	}
}

type queries struct {
	createSchema      string
	createEventsTable string
	insertEvent       string
	queryEventsBase   string
	queryEventsOffset string
	countEvents       string
	getTxDigestById   string
}
