package main

import (
	"fmt"
	"reflect"
	"strings"

	opregistry "github.com/smartcontractkit/chainlink-sui/deployment/ops/registry"
)

func main() {
	fmt.Println("All Operations with Input Templates:")
	fmt.Println("====================================")
	fmt.Println()

	// Use AllOperationsTyped from registry to stay in sync
	allOps := opregistry.AllOperationsTyped

	for i, op := range allOps {
		// op is of type any but contains an Operation[IN, OUT, DEP] value (not pointer)
		// Use reflection to get the input type from the operation
		opVal := reflect.ValueOf(op)

		opType := opVal.Type()

		if opType.Kind() == reflect.Struct && opType.NumField() >= 2 {
			def := opregistry.AllOperations[i].Def()

			// handler field is at index 1, even though it's unexported
			handlerField := opType.Field(1)
			handlerType := handlerField.Type

			if handlerType.Kind() == reflect.Func && handlerType.NumIn() >= 3 {
				// Function signature is: func(Bundle, DEP, IN) (OUT, error)
				// IN is the 3rd parameter (index 2)
				inputType := handlerType.In(2)

				// Access def fields directly
				id := def.ID
				versionField := def.Version.String()
				description := def.Description

				fmt.Printf("%d. ID: %s\n", i+1, id)
				fmt.Printf("   Version: %s\n", versionField) // Most operations use 0.1.0
				fmt.Printf("   Description: %s\n", description)
				fmt.Printf("   Input Type: %s\n", inputType.String())
				fmt.Println("   YAML Template:")

				yaml := generateYAMLTemplate(inputType, "   ", 0, make(map[reflect.Type]bool), 5)
				fmt.Print(yaml)
				fmt.Println()
			}
		}
	}

	fmt.Printf("Total operations: %d\n", len(allOps))
}

func generateYAMLTemplate(t reflect.Type, indent string, depth int, visited map[reflect.Type]bool, maxDepth int) string {
	// Check depth limit
	if depth > maxDepth {
		return " # ... (max depth reached)\n"
	}

	// Handle pointers
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// Check for cycles
	if visited[t] {
		return fmt.Sprintf(" # ... (circular reference to %s)\n", t.String())
	}

	switch t.Kind() {
	case reflect.Struct:
		visited[t] = true
		defer func() { delete(visited, t) }()

		var result strings.Builder
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)

			// Skip unexported fields
			if !field.IsExported() {
				continue
			}

			// Skip fields with json:"-" tag
			jsonTag := field.Tag.Get("json")
			if jsonTag == "-" {
				continue
			}

			// Get field name from json tag or use field name
			fieldName := getFieldName(field)

			result.WriteString(fmt.Sprintf("%s  %s:", indent, fieldName))

			fieldValue := generateFieldValue(field.Type, indent+"    ", depth+1, visited, maxDepth)
			result.WriteString(fieldValue)

			if !strings.HasSuffix(fieldValue, "\n") {
				result.WriteString("\n")
			}
		}
		return result.String()
	default:
		return " # " + t.String() + "\n"
	}
}

func generateFieldValue(t reflect.Type, indent string, depth int, visited map[reflect.Type]bool, maxDepth int) string {
	// Check depth limit
	if depth > maxDepth {
		return " ... (max depth reached)"
	}

	// Handle pointers
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	switch t.Kind() {
	case reflect.String:
		return " # string"
	case reflect.Bool:
		return " # bool"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return " # " + t.String()
	case reflect.Slice, reflect.Array:
		// Special case for []byte
		if t.Elem().Kind() == reflect.Uint8 {
			return " # " + t.String()
		}
		elemValue := generateFieldValue(t.Elem(), indent+"  ", depth+1, visited, maxDepth)
		return fmt.Sprintf("\n%s- %s", indent, strings.TrimSpace(elemValue))
	case reflect.Struct:
		structYAML := generateYAMLTemplate(t, indent, depth+1, visited, maxDepth)
		return "\n" + structYAML
	case reflect.Map:
		keyType := t.Key()
		valueType := t.Elem()
		valueStr := generateFieldValue(valueType, indent+"  ", depth+1, visited, maxDepth)

		keyExample := "example_key"
		if keyType.Kind() >= reflect.Int && keyType.Kind() <= reflect.Uint64 {
			keyExample = "123"
		}

		return fmt.Sprintf("\n%s%s: %s", indent, keyExample, strings.TrimSpace(valueStr))
	case reflect.Interface:
		return ` "interface{} - provide appropriate value"`
	default:
		return fmt.Sprintf(` "unknown_type_%s"`, t.Kind().String())
	}
}

func getFieldName(field reflect.StructField) string {
	// Try json tag first
	if jsonTag := field.Tag.Get("json"); jsonTag != "" {
		if parts := strings.Split(jsonTag, ","); len(parts) > 0 && parts[0] != "" {
			return parts[0]
		}
	}

	// Try yaml tag
	if yamlTag := field.Tag.Get("yaml"); yamlTag != "" {
		if parts := strings.Split(yamlTag, ","); len(parts) > 0 && parts[0] != "" {
			return parts[0]
		}
	}

	// Fall back to field name in lowercase
	return strings.ToLower(field.Name)
}
