package plugin

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWrappedNativeTokenForChain(t *testing.T) {
	t.Parallel()

	const (
		mainnetToken = "0x9258181f5ceac8dbffb7030890243caed69a9599d2886d957a9cb7656af3bdb3"
		testnetToken = "0x587c29de216efd4219573e08a1f6964d4fa7cb714518c2c8a0f29abfa264327d"
	)

	tests := []struct {
		name    string
		chainID *big.Int
		want    string
	}{
		{name: "mainnet hex identifier", chainID: big.NewInt(0x35834a8a), want: mainnetToken},
		{name: "testnet hex identifier", chainID: big.NewInt(0x4c78adac), want: testnetToken},
		{name: "mainnet decimal identifier", chainID: big.NewInt(897796746), want: mainnetToken},
		{name: "testnet decimal identifier", chainID: big.NewInt(1282977196), want: testnetToken},
		{name: "nil chain id", chainID: nil, want: ""},
		{name: "unknown chain id", chainID: big.NewInt(1), want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tt.want, wrappedNativeTokenForChain(tt.chainID))
		})
	}
}
