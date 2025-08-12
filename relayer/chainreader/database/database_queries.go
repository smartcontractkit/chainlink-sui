package database

const (
	CreateSchema = `
        CREATE SCHEMA IF NOT EXISTS sui;
    `

	CreateEventsTable = `
	CREATE TABLE IF NOT EXISTS sui.events (
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
    `

	InsertEvent = `
	INSERT INTO sui.events (
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
    `

	QueryEventsBase = `
	SELECT event_account_address, event_handle, event_offset, block_version, block_height, block_hash, block_timestamp, tx_digest, data
	FROM sui.events
	WHERE event_account_address = $1 AND event_handle = $2
    `

	QueryEventsOffset = `
	SELECT COALESCE(event_offset, 0) as event_offset, tx_digest, COUNT(*) OVER() as total_count
	FROM sui.events 
	WHERE event_account_address = $1 AND event_handle = $2 
	ORDER BY id DESC 
	LIMIT 1
	`

	CountEvents = `
	SELECT COUNT(*) as total_count
	FROM sui.events 
	WHERE event_account_address = $1 AND event_handle = $2 
	`

	GetTxDigestById = `
	SELECT tx_digest
	FROM sui.events
	WHERE id = $1
	`
)
