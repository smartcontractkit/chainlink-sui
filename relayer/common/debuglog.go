package common

// TEMPORARY debug instrumentation for session 39639b (Sui CCIP merkle-root read-path investigation).
// Writes NDJSON lines to a fixed log file. Remove this file once debugging is complete.

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

const debugLogPath = "/Users/faisal/Desktop/chainlink-sui/.cursor/debug-39639b.log"

const debugSessionID = "39639b"

var debugLogMu sync.Mutex

// DebugLog appends a single NDJSON line describing a runtime event to the debug log file.
// It is intentionally best-effort: any error is silently ignored so instrumentation never
// affects program behavior.
func DebugLog(location, message string, data map[string]any) {
	entry := map[string]any{
		"sessionId": debugSessionID,
		"timestamp": time.Now().UnixMilli(),
		"location":  location,
		"message":   message,
		"data":      data,
	}

	line, err := json.Marshal(entry)
	if err != nil {
		return
	}

	debugLogMu.Lock()
	defer debugLogMu.Unlock()

	f, err := os.OpenFile(debugLogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return
	}
	defer f.Close()

	_, _ = f.Write(append(line, '\n'))
}
