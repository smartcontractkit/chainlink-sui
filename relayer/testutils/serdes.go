package testutils

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/santhosh-tekuri/jsonschema/v5"
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

func ValidateJSON(input any, schemaJSON string) error {
	compiler := jsonschema.NewCompiler()
	// Add the schema as a resource instead of treating it as a URL
	schemaURL := "schema.json"
	err := compiler.AddResource(schemaURL, strings.NewReader(schemaJSON))
	if err != nil {
		return err
	}
	schema, err := compiler.Compile(schemaURL)
	if err != nil {
		return err
	}

	err = schema.Validate(input)
	return err
}
