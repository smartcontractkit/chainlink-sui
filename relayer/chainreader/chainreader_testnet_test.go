//go:build integration && testnet

package chainreader

// TODO: FIXME
// func TestChainReaderDevnet(t *testing.T) {
// 	runTestnetTest(t, testutils.DevnetUrl)
// }

// func TestChainReaderTestnet(t *testing.T) {
// 	runTestnetTest(t, testutils.TestnetUrl)
// }

// func runTestnetTest(t *testing.T, rpcUrl string) {
// 	logger := logger.Test(t)

// 	privateKey, publicKey, accountAddress := testutils.LoadAccountFromEnv(t, logger)
// 	if privateKey == nil {
// 		t.Skip("PRIVATE_KEY or ADDRESS environment variable is not set, skipping testnet test")
// 	}

// 	runChainReaderTest(t, logger, rpcUrl, accountAddress, publicKey, privateKey)
// }
