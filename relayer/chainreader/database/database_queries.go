package database

const (
	CreateSchema = `
        CREATE SCHEMA IF NOT EXISTS sui;
    `

	CreateEventsTable = `
	CREATE TABLE IF NOT EXISTS sui.events (
		event_account_address TEXT NOT NULL,
		event_handle TEXT NOT NULL,
		event_offset BIGINT NOT NULL,
		block_version BIGINT NOT NULL,
		block_height TEXT NOT NULL,
		block_hash BYTEA NOT NULL,
		block_timestamp BIGINT NOT NULL,
		data JSONB NOT NULL,
		PRIMARY KEY (event_account_address, event_handle, event_offset)
	);
    `

	InsertEvent = `
	INSERT INTO sui.events (
		event_account_address,
		event_handle,
		event_offset,
		block_version,
		block_height,
		block_hash,
		block_timestamp,
		data
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	ON CONFLICT DO NOTHING;
    `

	QueryEventsBase = `
	SELECT event_account_address, event_handle, event_offset, block_version, block_height, block_hash, block_timestamp, data
	FROM sui.events
	WHERE event_account_address = $1 AND event_handle = $2
    `

	QueryEventsOffset = `
	SELECT COALESCE(MAX(event_offset), 0) FROM sui.events
	WHERE event_account_address = $1 AND event_handle = $2
	`

	QueryTransactionVersionByID = `
	SELECT tx_version FROM sui.events
	WHERE id = $1
	`
)
