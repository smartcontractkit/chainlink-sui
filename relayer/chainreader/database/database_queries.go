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
    ALTER TABLE sui.events ADD COLUMN IF NOT EXISTS is_synthetic BOOLEAN DEFAULT FALSE;
    ALTER TABLE sui.events ADD COLUMN IF NOT EXISTS tx_idx BIGINT NOT NULL DEFAULT 0;
    `

	CreateCheckpointCursorsTable = `
	CREATE TABLE IF NOT EXISTS sui.checkpoint_cursors (
		id TEXT PRIMARY KEY,
		last_sequence BIGINT NOT NULL,
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);
	`

	CreateIndices = `
	CREATE INDEX IF NOT EXISTS idx_events_account_handle_timestamp ON sui.events(event_account_address, event_handle, block_timestamp DESC);
	CREATE INDEX IF NOT EXISTS idx_events_offset ON sui.events(event_account_address, event_handle, event_offset);
	CREATE INDEX IF NOT EXISTS idx_events_data_gin ON sui.events USING gin(data);
	CREATE INDEX IF NOT EXISTS idx_events_account_handle_id ON sui.events(event_account_address, event_handle, id DESC);
	CREATE INDEX IF NOT EXISTS idx_events_order ON sui.events(event_account_address, event_handle, (block_height::BIGINT), tx_idx, event_offset);
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
		data,
		is_synthetic,
		tx_idx
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	ON CONFLICT DO NOTHING;
    `

	QueryEventsBase = `
	SELECT event_account_address, event_handle, tx_idx, event_offset, block_version, block_height, block_hash, block_timestamp, tx_digest, data
	FROM sui.events
	WHERE event_account_address = $1 AND event_handle = $2
    `

	GetCheckpointCursorQuery = `
	SELECT last_sequence
	FROM sui.checkpoint_cursors
	WHERE id = $1
	`

	// UpsertCheckpointCursorQuery is monotonic: the stored sequence never regresses even if
	// callers race or submit stale values.
	UpsertCheckpointCursorQuery = `
	INSERT INTO sui.checkpoint_cursors (id, last_sequence) VALUES ($1, $2)
	ON CONFLICT (id) DO UPDATE
	SET last_sequence = GREATEST(sui.checkpoint_cursors.last_sequence, EXCLUDED.last_sequence),
	    updated_at = NOW();
	`
)
