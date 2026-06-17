package chainreaderutil

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"strconv"

	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
	"google.golang.org/protobuf/types/known/structpb"
)

var largeIntJSONFieldNames = map[string]struct{}{
	"destChainSelector":     {},
	"sourceChainSelector":   {},
	"sequenceNumber":        {},
	"dest_chain_selector":   {},
	"source_chain_selector": {},
	"sequence_number":       {},
}

// ParseEventDataFromJSON converts structpb event JSON into a map while preserving u64 precision
// where possible. structpb stores JSON numbers as float64, so values above 2^53 lose precision
// when read via AsInterface. json.Number avoids further loss during decode; exact u64 fields
// should still be overlaid from BCS when available.
func ParseEventDataFromJSON(jsonVal *structpb.Value) map[string]any {
	if jsonVal == nil {
		return map[string]any{}
	}

	raw, err := jsonVal.MarshalJSON()
	if err != nil {
		return map[string]any{}
	}

	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()

	var data map[string]any
	if err := dec.Decode(&data); err != nil || data == nil {
		return map[string]any{}
	}

	normalizeLargeIntFieldsInMap(data)
	return data
}

// NormalizeLargeIntEventFields converts known u64 event fields to decimal strings in-place.
func NormalizeLargeIntEventFields(data map[string]any) {
	normalizeLargeIntFieldsInMap(data)
}

func normalizeLargeIntFieldsInMap(data map[string]any) {
	for key, value := range data {
		if _, isLargeIntField := largeIntJSONFieldNames[key]; isLargeIntField {
			if s, ok := coerceToDecimalString(value); ok {
				data[key] = s
			}
			continue
		}

		switch nested := value.(type) {
		case map[string]any:
			normalizeLargeIntFieldsInMap(nested)
		case []any:
			for i, item := range nested {
				if nestedMap, ok := item.(map[string]any); ok {
					normalizeLargeIntFieldsInMap(nestedMap)
					nested[i] = nestedMap
				}
			}
		}
	}
}

func coerceToDecimalString(value any) (string, bool) {
	switch v := value.(type) {
	case json.Number:
		return v.String(), true
	case string:
		return v, true
	case uint64:
		return strconv.FormatUint(v, 10), true
	case uint32:
		return strconv.FormatUint(uint64(v), 10), true
	case int64:
		if v < 0 {
			return "", false
		}
		return strconv.FormatUint(uint64(v), 10), true
	case int:
		if v < 0 {
			return "", false
		}
		return strconv.FormatUint(uint64(v), 10), true
	case float64:
		return strconv.FormatUint(uint64(v), 10), true
	default:
		return "", false
	}
}

// OverlayLeadingU64FromBCS writes leading little-endian u64 fields from event BCS bytes into data.
// Used when structpb JSON rendering has already lost precision for large chain selectors.
func OverlayLeadingU64FromBCS(data map[string]any, bcs []byte, fieldNames ...string) {
	requiredLen := len(fieldNames) * 8
	if len(bcs) < requiredLen {
		return
	}

	for i, fieldName := range fieldNames {
		offset := i * 8
		value := binary.LittleEndian.Uint64(bcs[offset : offset+8])
		data[fieldName] = strconv.FormatUint(value, 10)
	}
}

// NormalizeLargeIntComparatorsToString converts uint64 filter values for known large-int event
// fields into decimal strings so SQL uses text comparison against string-stored JSONB values.
func NormalizeLargeIntComparatorsToString(expressions []query.Expression) []query.Expression {
	normalized := make([]query.Expression, len(expressions))
	for i, expr := range expressions {
		normalized[i] = normalizeExpressionLargeIntComparators(expr)
	}
	return normalized
}

func normalizeExpressionLargeIntComparators(expr query.Expression) query.Expression {
	if expr.Primitive != nil {
		if comp, ok := expr.Primitive.(*primitives.Comparator); ok {
			if _, isLargeIntField := largeIntJSONFieldNames[comp.Name]; isLargeIntField {
				newComp := *comp
				newComp.ValueComparators = make([]primitives.ValueComparator, len(comp.ValueComparators))
				for i, vc := range comp.ValueComparators {
					newComp.ValueComparators[i] = vc
					if s, ok := comparatorValueToDecimalString(vc.Value); ok {
						newComp.ValueComparators[i].Value = s
					}
				}
				return query.Expression{Primitive: &newComp}
			}
		}
	}

	if len(expr.BoolExpression.Expressions) == 0 {
		return expr
	}

	newExpr := expr
	newExpr.BoolExpression.Expressions = make([]query.Expression, len(expr.BoolExpression.Expressions))
	for i, subExpr := range expr.BoolExpression.Expressions {
		newExpr.BoolExpression.Expressions[i] = normalizeExpressionLargeIntComparators(subExpr)
	}
	return newExpr
}

func comparatorValueToDecimalString(value any) (string, bool) {
	return coerceToDecimalString(value)
}
