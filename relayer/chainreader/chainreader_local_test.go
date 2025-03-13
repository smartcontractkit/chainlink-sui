//go:build integration

package chainreader

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"math/big"
	"testing"

	sui "github.com/block-vision/sui-go-sdk"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	// "github.com/smartcontractkit/chainlink-internal-integrations/sui/relayer/testutils"
)

func TestChainReaderLocal(t *testing.T) {
	logger := logger.Test(t)

	privateKey, publicKey, accountAddress := testutils.LoadAccountFromEnv(t, logger)
	if privateKey == nil {
		newPublicKey, newPrivateKey, err := ed25519.GenerateKey(rand.Reader)
		require.NoError(t, err)
		privateKey = newPrivateKey
		publicKey = newPublicKey

		// Generate Sui address from public key
		accountAddress = testutils.DeriveAddressFromPublicKey(publicKey)

		logger.Debugw("Created account", "publicKey", hex.EncodeToString([]byte(publicKey)), "accountAddress", accountAddress)
	}

	err := testutils.StartSuiNode()
	require.NoError(t, err)
	logger.Debugw("Started Sui node")

	rpcUrl := "http://localhost:9000"
	client := sui.NewClient(rpcUrl)

	err = testutils.FundWithFaucet(logger, client, accountAddress)
	require.NoError(t, err)

	runChainReaderTest(t, logger, rpcUrl, accountAddress, publicKey, privateKey)
}

func runChainReaderTest(t *testing.T, logger logger.Logger, rpcUrl string, accountAddress string, publicKey ed25519.PublicKey, privateKey ed25519.PrivateKey) {
	// keystore := testutils.NewTestKeystore(t)
	// keystore.AddKey(privateKey)

	// client := sui.NewClient(rpcUrl)
	// getClient := func() (*sui.SuiClient, error) { return client, nil }

	// txmConfig := txm.DefaultConfigSet
	// txmgr, err := txm.New(logger, keystore, txmConfig, getClient)
	// require.NoError(t, err)

	// err = txmgr.Start(context.Background())
	// require.NoError(t, err)

	// publicKeyHex := hex.EncodeToString([]byte(publicKey))

	// // Compile and publish test module
	// packageObjectId := testutils.CompileAndPublishTestModule(t, client, privateKey, accountAddress)

	// txId := uuid.New().String()
	// // Sui transaction enqueue is different from Aptos
	// err = txmgr.Enqueue(
	// 	txId,
	// 	getSampleTxMetadata(),
	// 	accountAddress,
	// 	publicKeyHex,
	// 	"move_call", // Sui transaction type
	// 	[]string{},  // type args
	// 	[]string{packageObjectId, "echo", "init"}, // target, module, function
	// 	[]any{}, // args
	// 	true,    // simulate
	// )
	// require.NoError(t, err)

	// confirmed := false
	// for i := 0; i < 10; i++ {
	// 	time.Sleep(time.Second * 1)
	// 	status, err := txmgr.GetStatus(txId)
	// 	require.NoError(t, err)
	// 	if status != commontypes.Unconfirmed {
	// 		confirmed = true
	// 		break
	// 	}
	// }
	// require.True(t, confirmed)

	// config := ChainReaderConfig{
	// 	Modules: map[string]*Module{
	// 		"testContract": {
	// 			Name: "echo",
	// 			Functions: map[string]*Function{
	// 				"replacementNameEchoU64": {
	// 					Name: "echo_u64",
	// 					Params: []Parameter{
	// 						{
	// 							Name:     "Value1",
	// 							Type:     "u64",
	// 							Required: true,
	// 						},
	// 					},
	// 				},
	// 				"echo_u32_u64_tuple": {
	// 					Params: []Parameter{
	// 						{
	// 							Name:     "Value1",
	// 							Type:     "u32",
	// 							Required: true,
	// 						},
	// 						{
	// 							Name:     "Value2",
	// 							Type:     "u64",
	// 							Required: true,
	// 						},
	// 					},
	// 				},
	// 				"echo_string": {
	// 					Params: []Parameter{
	// 						{
	// 							Name:     "Value1",
	// 							Type:     "string",
	// 							Required: true,
	// 						},
	// 					},
	// 				},
	// 				"echo_byte_vector": {
	// 					Params: []Parameter{
	// 						{
	// 							Name:     "Value1",
	// 							Type:     "vector<u8>",
	// 							Required: true,
	// 						},
	// 					},
	// 				},
	// 				"echo_byte_vector_vector": {
	// 					Params: []Parameter{
	// 						{
	// 							Name:     "Value1",
	// 							Type:     "vector<vector<u8>>",
	// 							Required: true,
	// 						},
	// 					},
	// 				},
	// 				"echo_u256": {
	// 					Params: []Parameter{
	// 						{
	// 							Name:     "Value1",
	// 							Type:     "u256",
	// 							Required: true,
	// 						},
	// 					},
	// 				},
	// 			},
	// 		},
	// 	},
	// }

	// binding := commontypes.BoundContract{
	// 	Name:    "testContract",
	// 	Address: packageObjectId,
	// }

	// chainReader := NewChainReader(logger, client, config)
	// err = chainReader.Bind(context.Background(), []commontypes.BoundContract{binding})
	// require.NoError(t, err)

	// confidenceLevel := primitives.Finalized

	// // Test cases similar to Aptos but adapted for Sui
	// expectedUint64 := uint64(42)
	// var retUint64 uint64
	// err = chainReader.GetLatestValue(context.Background(), binding.ReadIdentifier("replacementNameEchoU64"), confidenceLevel, struct {
	// 	Value1 uint64
	// }{Value1: expectedUint64}, &retUint64)
	// require.NoError(t, err)
	// require.Equal(t, expectedUint64, retUint64)

	// // Continue with other test cases...
	// // Note: Some of these might need to be adjusted based on Sui's specific type system
}

func getSampleTxMetadata() *commontypes.TxMeta {
	workflowID := "sample-workflow-id"
	return &commontypes.TxMeta{
		WorkflowExecutionID: &workflowID,
		GasLimit:            big.NewInt(21000),
	}
}
