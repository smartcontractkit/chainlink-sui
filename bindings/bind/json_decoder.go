package bind

import (
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"strings"

	"github.com/mitchellh/mapstructure"
)

// DecodeJSONReturn decodes gRPC/JSON Move return values into the provided target.
func DecodeJSONReturn(data any, target any) error {
	if target == nil {
		return fmt.Errorf("target cannot be nil")
	}

	if raw, ok := data.(json.RawMessage); ok {
		var intermediate any
		if err := json.Unmarshal(raw, &intermediate); err != nil {
			return fmt.Errorf("json unmarshal failed: %w", err)
		}

		return DecodeJSONReturn(intermediate, target)
	}

	if reflect.TypeOf(data) == reflect.TypeOf(target).Elem() {
		reflect.ValueOf(target).Elem().Set(reflect.ValueOf(data))
		return nil
	}

	targetType := reflect.TypeOf(target).Elem()

	bigPtrT := reflect.TypeOf((*big.Int)(nil))
	bigValT := bigPtrT.Elem()
	if targetType == bigValT || targetType == bigPtrT {
		return decodeBigInt(data, target)
	}

	return decodeWithMapstructure(data, target)
}

func decodeBigInt(data any, target any) error {
	str, ok := data.(string)
	if !ok {
		return fmt.Errorf("big.Int decode: expected string, got %T", data)
	}

	bi, success := new(big.Int).SetString(str, 10)
	if !success {
		return fmt.Errorf("big.Int decode: invalid number %q", str)
	}

	targetValue := reflect.ValueOf(target).Elem()
	targetType := targetValue.Type()
	bigPtrT := reflect.TypeOf((*big.Int)(nil))
	bigValT := bigPtrT.Elem()

	if targetType == bigValT {
		targetValue.Set(reflect.ValueOf(*bi))
	} else {
		targetValue.Set(reflect.ValueOf(bi))
	}

	return nil
}

func decodeWithMapstructure(data any, target any) error {
	config := &mapstructure.DecoderConfig{
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			UnifiedTypeConverterHook,
			mapstructure.StringToTimeDurationHookFunc(),
		),
		Result:           target,
		WeaklyTypedInput: true,
		TagName:          "json",
		MatchName:        fuzzyFieldMatcher,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return fmt.Errorf("failed to create decoder: %w", err)
	}

	return decoder.Decode(data)
}

func fuzzyFieldMatcher(mapKey, fieldName string) bool {
	mk := strings.ReplaceAll(mapKey, "_", "")
	fn := strings.ReplaceAll(fieldName, "_", "")
	return strings.EqualFold(mk, fn)
}

func handleSingleFieldStruct(t reflect.Type, data any, decodeFn func(any, any) error) (any, error) {
	field := t.Field(0)
	newStructVal := reflect.New(t).Elem()
	fieldPtr := newStructVal.Field(0).Addr().Interface()

	if err := decodeFn(data, fieldPtr); err != nil {
		return nil, fmt.Errorf("failed decoding for single-field struct %v field %s (%v): %w",
			t, field.Name, field.Type, err)
	}

	return newStructVal.Interface(), nil
}
