package testutils

// Network URLs
const (
	DevnetURL      = "https://fullnode.devnet.sui.io:443"
	TestnetURL     = "https://fullnode.testnet.sui.io:443"
	LocalURL       = "http://127.0.0.1:9000"
	LocalFaucetURL = "http://127.0.0.1:9123/gas"
	LocalGrpcURL   = "127.0.0.1:9000"
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
