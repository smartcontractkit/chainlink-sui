package rmn

import (
	"testing"

	"github.com/smartcontractkit/chainlink-ccip/deployment/fastcurse"
	"github.com/stretchr/testify/require"
)

func TestSubjectsToBytes_distinctBackingArrays(t *testing.T) {
	t.Parallel()

	a := fastcurse.Subject{1}
	b := fastcurse.Subject{2}
	got := subjectsToBytes([]fastcurse.Subject{a, b})

	require.Len(t, got, 2)
	require.NotEqual(t, got[0], got[1])
	require.Equal(t, []byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, got[0])
	require.Equal(t, []byte{2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, got[1])

	// Mutating one slice must not affect the other.
	got[0][0] = 99
	require.Equal(t, byte(1), got[1][0])
}
