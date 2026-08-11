package lanes

import (
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestRemoteAddressBytesToHex(t *testing.T) {
	t.Parallel()

	evmAddr := common.HexToAddress("0xd3E190f381f06DC0d289590fd452C42Fa2DAC586")

	t.Run("20 byte EVM address is left padded", func(t *testing.T) {
		got, err := remoteAddressBytesToHex(evmAddr.Bytes())
		require.NoError(t, err)
		require.Equal(t, "0x000000000000000000000000d3e190f381f06dc0d289590fd452c42fa2dac586", got)
	})

	t.Run("32 byte address is preserved", func(t *testing.T) {
		raw := make([]byte, 32)
		for i := range raw {
			raw[i] = byte(i + 1)
		}
		got, err := remoteAddressBytesToHex(raw)
		require.NoError(t, err)
		require.Equal(t, "0x"+hex.EncodeToString(raw), got)
	})

	t.Run("32 byte left padded EVM address is preserved", func(t *testing.T) {
		padded := make([]byte, 32)
		copy(padded[12:], evmAddr.Bytes())
		got, err := remoteAddressBytesToHex(padded)
		require.NoError(t, err)
		require.Equal(t, "0x"+hex.EncodeToString(padded), got)
	})

	t.Run("empty bytes rejected", func(t *testing.T) {
		_, err := remoteAddressBytesToHex(nil)
		require.Error(t, err)
		require.Contains(t, err.Error(), "empty address bytes")
	})

	t.Run("longer than 32 bytes rejected", func(t *testing.T) {
		_, err := remoteAddressBytesToHex(make([]byte, 33))
		require.Error(t, err)
		require.Contains(t, err.Error(), "address longer than 32 bytes")
	})
}
