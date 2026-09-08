package deployment

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// paddedSuiAddr pads a short hex suffix to the canonical 64-hex-character form.
func paddedSuiAddr(shortHex string) string {
	return "0x" + strings.Repeat("0", 64-len(shortHex)) + shortHex
}

func TestSuiAddressesEqual(t *testing.T) {
	t.Parallel()

	// Hex case and zero padding are spellings, not identity; surrounding whitespace is
	// tolerated for addresses.
	require.True(t, SuiAddressesEqual("0xab", "0x00000aB"))
	require.True(t, SuiAddressesEqual("0xAB", paddedSuiAddr("ab")))
	require.True(t, SuiAddressesEqual(paddedSuiAddr("AB"), "0xab"))
	require.True(t, SuiAddressesEqual(" 0xaB ", paddedSuiAddr("ab")))

	require.False(t, SuiAddressesEqual("0xab", paddedSuiAddr("ac")))
	require.False(t, SuiAddressesEqual("", paddedSuiAddr("ab")))
	require.True(t, SuiAddressesEqual("", ""))

	// Non-address strings compare byte-identical only — whitespace included.
	require.True(t, SuiAddressesEqual("0xccip-pkg", "0xccip-pkg"))
	require.False(t, SuiAddressesEqual("0xccip-pkg", "0xCCIP-pkg"))
	require.False(t, SuiAddressesEqual(" 0xccip-pkg ", "0xccip-pkg"))
}

func TestQualifierContainsSuiAddress(t *testing.T) {
	t.Parallel()

	const short, qualType = "0xab", "SuiCCIP"
	padded := paddedSuiAddr("ab")

	// The shim pattern in every spelling of the embedded address.
	for _, qualifier := range []string{
		short + "-" + qualType,
		padded + "-" + qualType,
		"0xAB-" + qualType,
		paddedSuiAddr("AB") + "-" + qualType,
		strings.ToUpper(padded),
	} {
		require.True(t, QualifierContainsSuiAddress(qualifier, short), "qualifier %q", qualifier)
		require.True(t, QualifierContainsSuiAddress(qualifier, padded), "qualifier %q", qualifier)
	}

	for _, qualifier := range []string{ChainSingletonQualifier, "CCIP-BnM", "CLLCCIP", "RMNMCMS"} {
		require.False(t, QualifierContainsSuiAddress(qualifier, short), "qualifier %q", qualifier)
	}

	// Non-hex addresses keep the plain case-insensitive substring behaviour.
	coinType := "0x<pkg>::module::STRUCT"
	require.True(t, QualifierContainsSuiAddress(coinType+"-x", coinType))
	require.False(t, QualifierContainsSuiAddress("CCIP-BnM", coinType))

	// The all-zero address skips the unpadded needle (bare "0x" matches too much).
	require.True(t, QualifierContainsSuiAddress("0x"+strings.Repeat("0", 64), "0x0"))
	require.False(t, QualifierContainsSuiAddress("0x", "0x0"))
}
