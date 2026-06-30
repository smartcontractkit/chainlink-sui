package suigrpcconn

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewConnectionTransportModes(t *testing.T) {
	t.Parallel()

	tlsClient := NewConnectionWithAuth("fullnode.testnet.sui.io:443", "tok", true)
	require.NotNil(t, tlsClient)
	require.NotEmpty(t, tlsClient.dialOptions())

	plainClient := NewConnectionWithAuth("127.0.0.1:9000", "tok", false)
	require.NotNil(t, plainClient)
	require.NotEmpty(t, plainClient.dialOptions())
}
