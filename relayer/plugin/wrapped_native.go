package plugin

import "math/big"

// SUI's FeeQuoter key is the legacy coin::CoinMetadata<SUI> object ID, which differs per network,
// so CCIP config refresh needs the key resolved from the node's chain identifier.
var wrappedNativeTokenByChain = map[string]string{
	"0x35834a8a": "0x9258181f5ceac8dbffb7030890243caed69a9599d2886d957a9cb7656af3bdb3",
	"0x4c78adac": "0x587c29de216efd4219573e08a1f6964d4fa7cb714518c2c8a0f29abfa264327d",
}

func wrappedNativeTokenForChain(chainID *big.Int) string {
	if chainID == nil {
		return ""
	}

	return wrappedNativeTokenByChain["0x"+chainID.Text(16)]
}
