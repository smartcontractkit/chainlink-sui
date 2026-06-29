package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGrpcTargetUsesTLS(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		target string
		want   bool
	}{
		{name: "public testnet full node", target: "fullnode.testnet.sui.io:443", want: true},
		{name: "public hostname without port", target: "fullnode.testnet.sui.io", want: true},
		{name: "local node", target: "127.0.0.1:9000", want: false},
		{name: "local host without port", target: "127.0.0.1", want: false},
		{name: "ipv6 local node", target: "[::1]:9000", want: false},
		{name: "empty target", target: "", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, grpcTargetUsesTLS(tt.target))
		})
	}
}
