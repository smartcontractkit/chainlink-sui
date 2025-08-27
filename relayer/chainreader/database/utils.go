package database

import (
	"fmt"
	"strings"

	aptosCRUtils "github.com/smartcontractkit/chainlink-aptos/relayer/chainreader/utils"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
)

// Taken from the Aptos ChainReader DB methods: https://github.com/smartcontractkit/chainlink-aptos/blob/develop/relayer/chainreader/db/db.go#L257
func BuildSQLCondition(expr query.Expression, args *[]any, argCount *int) (string, error) {
	//nolint:all
	if expr.IsPrimitive() {
		switch v := expr.Primitive.(type) {
		case *primitives.Comparator:
			conditions := []string{}
			for _, valueCmp := range v.ValueComparators {
				jsonPath, err := aptosCRUtils.BuildJsonPathExpr("data", v.Name)
				if err != nil {
					return "", fmt.Errorf("invalid field name %s: %w", v.Name, err)
				}

				var condition string
				if aptosCRUtils.IsNumeric(valueCmp.Value) {
					condition = fmt.Sprintf("CAST(%s AS numeric) %s $%d", jsonPath, operatorSQL(valueCmp.Operator), *argCount)
				} else {
					condition = fmt.Sprintf("%s %s $%d", jsonPath, operatorSQL(valueCmp.Operator), *argCount)
				}

				*args = append(*args, valueCmp.Value)
				*argCount++
				conditions = append(conditions, condition)
			}

			return "(" + strings.Join(conditions, " AND ") + ")", nil

		case *primitives.Timestamp:
			condition := fmt.Sprintf("block_timestamp %s $%d", operatorSQL(v.Operator), *argCount)
			*args = append(*args, v.Timestamp)
			*argCount++

			return condition, nil

		case *primitives.Confidence:
			// Confidence filter isn't applicable in the context of Aptos
			return "TRUE", nil

		default:
			return "", fmt.Errorf("unsupported primitive type: %T", expr.Primitive)
		}
	} else {
		if len(expr.BoolExpression.Expressions) < 2 {
			return "", fmt.Errorf("boolean expression must have at least 2 expressions")
		}

		var subConditions []string
		for _, subExpr := range expr.BoolExpression.Expressions {
			subCond, err := BuildSQLCondition(subExpr, args, argCount)
			if err != nil {
				return "", err
			}
			subConditions = append(subConditions, subCond)
		}

		operator := " AND "
		if expr.BoolExpression.BoolOperator == query.OR {
			operator = " OR "
		}

		return "(" + strings.Join(subConditions, operator) + ")", nil
	}
}
