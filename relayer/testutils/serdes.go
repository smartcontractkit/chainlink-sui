package testutils

import (
	"encoding/json"
	"testing"
)

func ExtractStruct[T any](t *testing.T, payload any) *T {
	t.Helper()
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Failed to marshal data: %v", err)
	}

	var obj T
	if err := json.Unmarshal(jsonBytes, &obj); err != nil {
		t.Fatalf("Failed to unmarshal payload: %v", err)
	}

	return &obj
}
