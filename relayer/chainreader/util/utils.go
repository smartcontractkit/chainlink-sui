package util

import (
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
)

func ExtractTimestampFilter(expressions []query.Expression) (uint64, bool) {
	for _, expr := range expressions {
		if expr.IsPrimitive() {
			if tsExpr, ok := expr.Primitive.(*primitives.Timestamp); ok {
				if tsExpr.Operator == primitives.Gte {
					return tsExpr.Timestamp, true
				}
			}
		}
	}

	return 0, false
}

func IsNumeric(value any) bool {
	_, ok := value.(uint64)
	return ok
}
