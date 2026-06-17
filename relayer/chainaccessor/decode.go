package chainaccessor

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
)

// The values returned by client.ReadFunction (devInspect) are JSON-decoded Move
// return values. Sui encodes integers wider than 32 bits as strings, so these
// helpers tolerate string, json.Number, float64 and native integer forms.

func asUint64(v any) (uint64, error) {
	switch t := v.(type) {
	case uint64:
		return t, nil
	case int64:
		return uint64(t), nil
	case int:
		return uint64(t), nil
	case float64:
		return uint64(t), nil
	case json.Number:
		n, err := t.Int64()
		if err != nil {
			return 0, fmt.Errorf("invalid numeric value %q: %w", t.String(), err)
		}
		return uint64(n), nil
	case string:
		n, ok := new(big.Int).SetString(t, 10)
		if !ok {
			return 0, fmt.Errorf("invalid uint64 string %q", t)
		}
		return n.Uint64(), nil
	default:
		return 0, fmt.Errorf("cannot convert %T to uint64", v)
	}
}

func asBigInt(v any) (*big.Int, error) {
	switch t := v.(type) {
	case string:
		n, ok := new(big.Int).SetString(t, 10)
		if !ok {
			return nil, fmt.Errorf("invalid big.Int string %q", t)
		}
		return n, nil
	case float64:
		return big.NewInt(int64(t)), nil
	case json.Number:
		n, ok := new(big.Int).SetString(t.String(), 10)
		if !ok {
			return nil, fmt.Errorf("invalid big.Int value %q", t.String())
		}
		return n, nil
	case uint64:
		return new(big.Int).SetUint64(t), nil
	case int64:
		return big.NewInt(t), nil
	default:
		return nil, fmt.Errorf("cannot convert %T to big.Int", v)
	}
}

func asBool(v any) (bool, error) {
	b, ok := v.(bool)
	if !ok {
		return false, fmt.Errorf("cannot convert %T to bool", v)
	}
	return b, nil
}

func asMap(v any) (map[string]any, error) {
	m, ok := v.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("expected struct/object, got %T", v)
	}
	return m, nil
}

// field looks up a key in a Move struct result map, accepting either the raw
// field name or a JSON-tag style variant.
func field(m map[string]any, key string) (any, error) {
	if v, ok := m[key]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("missing field %q in result", key)
}

func uint64Field(m map[string]any, key string) (uint64, error) {
	v, err := field(m, key)
	if err != nil {
		return 0, err
	}
	return asUint64(v)
}

// asBytes tolerantly decodes a Move vector<u8> as represented in event JSON:
// a 0x-prefixed hex string, a base64 string, or a slice of byte-valued numbers.
func asBytes(v any) ([]byte, error) {
	switch t := v.(type) {
	case nil:
		return nil, nil
	case []byte:
		return t, nil
	case string:
		if t == "" {
			return nil, nil
		}
		if strings.HasPrefix(t, "0x") || strings.HasPrefix(t, "0X") {
			return hex.DecodeString(t[2:])
		}
		if b, err := base64.StdEncoding.DecodeString(t); err == nil {
			return b, nil
		}
		return hex.DecodeString(t)
	case []any:
		b := make([]byte, len(t))
		for i, e := range t {
			n, err := asUint64(e)
			if err != nil {
				return nil, fmt.Errorf("byte at index %d: %w", i, err)
			}
			b[i] = byte(n)
		}
		return b, nil
	default:
		return nil, fmt.Errorf("cannot convert %T to bytes", v)
	}
}

func bytesField(m map[string]any, key string) ([]byte, error) {
	v, err := field(m, key)
	if err != nil {
		return nil, err
	}
	return asBytes(v)
}

// toBytes32 left-pads (or right-truncates) a byte slice into a [32]byte.
func toBytes32(b []byte) [32]byte {
	var out [32]byte
	if len(b) >= 32 {
		copy(out[:], b[len(b)-32:])
		return out
	}
	copy(out[32-len(b):], b)
	return out
}

func mapField(m map[string]any, key string) (map[string]any, error) {
	v, err := field(m, key)
	if err != nil {
		return nil, err
	}
	return asMap(v)
}

// firstReturn returns the single expected return value from a view call.
func firstReturn(values []any) (any, error) {
	if len(values) == 0 {
		return nil, fmt.Errorf("view call returned no values")
	}
	return values[0], nil
}
