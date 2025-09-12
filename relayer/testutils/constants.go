package testutils

// Network URLs
const (
	DevnetUrl  = "https://fullnode.devnet.sui.io:443"
	TestnetUrl = "https://fullnode.testnet.sui.io:443"
)

// Network environments
const (
	SuiMainnet  = "mainnet"
	SuiTestnet  = "testnet"
	SuiDevnet   = "devnet"
	SuiLocalnet = "localnet"
)

const DefaultByteSize = 32
const SignatureComponents = 2 // R and S components in ECDSA signature
