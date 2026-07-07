package bind

import (
	"math/big"
	"testing"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/stretchr/testify/require"
)

// TestDeserializeBCS_EmptyVector_Address exercises the regression case where an
// empty vector<address> panicked in reflect.SliceOf because the element type was
// never populated. See bcs_deserialize.go bcsDeserializeSlice.
func TestDeserializeBCS_EmptyVector_Address(t *testing.T) {
	t.Parallel()

	// BCS-encoded empty vector: single ULEB128 zero byte.
	data := []byte{0x00}

	got, err := DeserializeBCS(data, []string{"vector<address>"})
	require.NoError(t, err)
	require.Len(t, got, 1)

	slice, ok := got[0].([]models.SuiAddress)
	require.True(t, ok, "expected []models.SuiAddress, got %T", got[0])
	require.Empty(t, slice)
}

// TestDeserializeBCS_EmptyVector_ByInnerType covers the same regression across
// every inner type reachable from bcsDeserializeType.
func TestDeserializeBCS_EmptyVector_ByInnerType(t *testing.T) {
	t.Parallel()

	// BCS-encoded empty vector: single ULEB128 zero byte.
	empty := []byte{0x00}

	cases := []struct {
		moveType string
		expected any
	}{
		{"vector<bool>", []bool{}},
		{"vector<u16>", []uint16{}},
		{"vector<u32>", []uint32{}},
		{"vector<u64>", []uint64{}},
		{"vector<u128>", []*big.Int{}},
		{"vector<u256>", []*big.Int{}},
		{"vector<address>", []models.SuiAddress{}},
		{"vector<0x1::string::String>", []string{}},
		{"vector<vector<u64>>", [][]uint64{}},
		{"vector<vector<address>>", [][]models.SuiAddress{}},
		// Custom Move structs fall through to bcsDeserializeAddress at decode time.
		{"vector<0x2::coin::Coin<0x2::sui::SUI>>", []models.SuiAddress{}},
	}

	for _, tc := range cases {
		t.Run(tc.moveType, func(t *testing.T) {
			t.Parallel()

			got, err := DeserializeBCS(empty, []string{tc.moveType})
			require.NoError(t, err)
			require.Len(t, got, 1)
			require.IsType(t, tc.expected, got[0])
			require.Empty(t, got[0])
		})
	}
}

// TestDeserializeBCS_NonEmptyVector_Address is a smoke test to make sure the
// fix preserves the happy path for non-empty vectors.
func TestDeserializeBCS_NonEmptyVector_Address(t *testing.T) {
	t.Parallel()

	// One address = 32 bytes, all 0xAB.
	var addr [32]byte
	for i := range addr {
		addr[i] = 0xAB
	}
	data := make([]byte, 0, 33)
	data = append(data, 0x01) // ULEB128 length = 1
	data = append(data, addr[:]...)

	got, err := DeserializeBCS(data, []string{"vector<address>"})
	require.NoError(t, err)
	require.Len(t, got, 1)
	slice, ok := got[0].([]models.SuiAddress)
	require.True(t, ok, "expected []models.SuiAddress, got %T", got[0])
	require.Len(t, slice, 1)
	require.Equal(t, models.SuiAddress("0xabababababababababababababababababababababababababababababababab"), slice[0])
}
