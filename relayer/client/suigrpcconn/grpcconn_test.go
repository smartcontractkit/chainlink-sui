package suigrpcconn

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewGrpcConnTransportModes(t *testing.T) {
	t.Parallel()

	tlsConn := NewGrpcConn("fullnode.testnet.sui.io:443", WithTransportSecurity(true))
	require.NotNil(t, tlsConn)
	require.NotEmpty(t, tlsConn.dialOpts)

	plainConn := NewGrpcConn("127.0.0.1:9000", WithTransportSecurity(false))
	require.NotNil(t, plainConn)
	require.NotEmpty(t, plainConn.dialOpts)
}
