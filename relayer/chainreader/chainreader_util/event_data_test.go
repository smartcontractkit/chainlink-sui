package chainreaderutil

import (
	"encoding/binary"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
)

func TestOverlayLeadingU64FromBCS_CCIPMessageSentSelector(t *testing.T) {
	t.Parallel()

	const destChainSelector = uint64(909606746561742123)
	const sequenceNumber = uint64(1)

	bcs := make([]byte, 16)
	binary.LittleEndian.PutUint64(bcs[0:8], destChainSelector)
	binary.LittleEndian.PutUint64(bcs[8:16], sequenceNumber)

	data := map[string]any{
		"destChainSelector": float64(909606746561742100),
		"sequenceNumber":    float64(1),
	}

	OverlayLeadingU64FromBCS(data, bcs, "destChainSelector", "sequenceNumber")

	require.Equal(t, "909606746561742123", data["destChainSelector"])
	require.Equal(t, "1", data["sequenceNumber"])
}

func TestNormalizeLargeIntComparatorsToString(t *testing.T) {
	t.Parallel()

	const destChainSelector = uint64(909606746561742123)

	expressions := []query.Expression{
		{
			Primitive: &primitives.Comparator{
				Name: "destChainSelector",
				ValueComparators: []primitives.ValueComparator{
					{Value: destChainSelector, Operator: primitives.Eq},
				},
			},
		},
	}

	normalized := NormalizeLargeIntComparatorsToString(expressions)
	comp := normalized[0].Primitive.(*primitives.Comparator)
	require.Equal(t, "909606746561742123", comp.ValueComparators[0].Value)
}

func TestParseEventDataFromJSON_PreservesLargeIntAsString(t *testing.T) {
	t.Parallel()

	jsonVal, err := structpb.NewValue(map[string]any{
		"dest_chain_selector": "909606746561742123",
		"sequence_number":     json.Number("1"),
	})
	require.NoError(t, err)

	data := ParseEventDataFromJSON(jsonVal)
	require.Equal(t, "909606746561742123", data["dest_chain_selector"])
	require.Equal(t, "1", data["sequence_number"])
}
