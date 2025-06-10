//go:build integration

package loop

import (
	"context"
	"crypto/ed25519"
	"fmt"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil/sqltest"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
	"github.com/smartcontractkit/chainlink-sui/relayer/keystore"
	"github.com/smartcontractkit/chainlink-sui/relayer/testutils"
)

//nolint:paralleltest
func TestLoopChainReaderLocal(t *testing.T) {
	log := logger.Test(t)

	var err error
	accountAddress := testutils.GetAccountAndKeyFromSui(t, log)
	cmd, err := testutils.StartSuiNode(testutils.CLI)
	require.NoError(t, err)

	// Ensure the process is killed when the test completes.
	t.Cleanup(func() {
		if cmd.Process != nil {
			perr := cmd.Process.Kill()
			if perr != nil {
				t.Logf("Failed to kill process: %v", perr)
			}
		}
	})

	log.Debugw("Started Sui node")

	err = testutils.FundWithFaucet(log, testutils.SuiLocalnet, accountAddress)
	require.NoError(t, err)

	runLoopChainReaderEchoTest(t, log, testutils.LocalUrl)
}

func runLoopChainReaderEchoTest(t *testing.T, log logger.Logger, rpcUrl string) {
	t.Helper()
	ctx := context.Background()

	accountAddress := testutils.GetAccountAndKeyFromSui(t, log)
	keystoreInstance, keystoreErr := keystore.NewSuiKeystore(log, "")
	require.NoError(t, keystoreErr)

	privateKey, err := keystoreInstance.GetPrivateKeyByAddress(accountAddress)
	require.NoError(t, err)
	publicKey := privateKey.Public().(ed25519.PublicKey)
	_ = []byte(publicKey) // Not used directly in this test

	relayerClient, clientErr := client.NewPTBClient(log, rpcUrl, nil, 10*time.Second, keystoreInstance, 5, "WaitForLocalExecution")
	require.NoError(t, clientErr)

	faucetFundErr := testutils.FundWithFaucet(log, testutils.SuiLocalnet, accountAddress)
	require.NoError(t, faucetFundErr)

	contractPath := testutils.BuildSetup(t, "contracts/test")
	testutils.BuildContract(t, contractPath)

	packageId, publishOutput, err := testutils.PublishContract(t, "TestContract", contractPath, accountAddress, nil)
	require.NoError(t, err)

	log.Debugw("Published Contract", "packageId", packageId)
	log.Debugw("Publish output", "output", publishOutput)

	// Set up the base ChainReader with echo function configurations
	chainReaderConfig := chainreader.ChainReaderConfig{
		IsLoopPlugin: true,
		Modules: map[string]*chainreader.ChainReaderModule{
			"echo": {
				Name: "echo",
				Functions: map[string]*chainreader.ChainReaderFunction{
					"echo_u64": {
						Name:          "echo_u64",
						SignerAddress: accountAddress,
						Params: []codec.SuiFunctionParam{
							{
								Type:     "u64",
								Name:     "val",
								Required: true,
							},
						},
					},
					"echo_u256": {
						Name:          "echo_u256",
						SignerAddress: accountAddress,
						Params: []codec.SuiFunctionParam{
							{
								Type:     "u256",
								Name:     "val",
								Required: true,
							},
						},
					},
					"echo_u32_u64_tuple": {
						Name:          "echo_u32_u64_tuple",
						SignerAddress: accountAddress,
						Params: []codec.SuiFunctionParam{
							{
								Type:     "u32",
								Name:     "val1",
								Required: true,
							},
							{
								Type:     "u64",
								Name:     "val2",
								Required: true,
							},
						},
					},
					"echo_string": {
						Name:          "echo_string",
						SignerAddress: accountAddress,
						Params: []codec.SuiFunctionParam{
							{
								Type:     "0x1::string::String",
								Name:     "val",
								Required: true,
							},
						},
					},
					"echo_byte_vector": {
						Name:          "echo_byte_vector",
						SignerAddress: accountAddress,
						Params: []codec.SuiFunctionParam{
							{
								Type:     "vector<u8>",
								Name:     "val",
								Required: true,
							},
						},
					},
					"echo_byte_vector_vector": {
						Name:          "echo_byte_vector_vector",
						SignerAddress: accountAddress,
						Params: []codec.SuiFunctionParam{
							{
								Type:     "vector<vector<u8>>",
								Name:     "val",
								Required: true,
							},
						},
					},
				},
				Events: map[string]*chainreader.ChainReaderEvent{
					"single_value_event": {
						Name:      "single_value_event",
						EventType: "SingleValueEvent",
					},
					"double_value_event": {
						Name:      "double_value_event",
						EventType: "DoubleValueEvent",
					},
					"triple_value_event": {
						Name:      "triple_value_event",
						EventType: "TripleValueEvent",
					},
				},
			},
		},
	}

	echoBinding := types.BoundContract{
		Name:    "echo",
		Address: packageId, // Package ID of the deployed echo contract
	}

	// Set up DB
	datastoreUrl := os.Getenv("TEST_DB_URL")
	if datastoreUrl == "" {
		t.Skip("Skipping persistent tests as TEST_DB_URL is not set in CI")
	}
	db := sqltest.NewDB(t, datastoreUrl)

	// Create the base chain reader
	chainReader, err := chainreader.NewChainReader(ctx, log, relayerClient, chainReaderConfig, db)
	require.NoError(t, err)

	// Wrap the base chain reader with loop chain reader
	loopReader := NewLoopChainReader(log, chainReader)

	// Bind the contracts to the loop reader
	err = loopReader.Bind(context.Background(), []types.BoundContract{echoBinding})
	require.NoError(t, err)

	log.Debugw("LoopChainReader setup complete")

	// Test 1: echo_u64 function call
	t.Run("LoopReader_GetLatestValue_EchoU64", func(t *testing.T) {
		testValue := uint64(42)
		var retUint64 uint64

		err = loopReader.GetLatestValue(
			context.Background(),
			strings.Join([]string{packageId, echoBinding.Name, "echo_u64"}, "-"),
			primitives.Finalized,
			map[string]any{
				"val": testValue,
			},
			&retUint64,
		)
		require.NoError(t, err)
		require.Equal(t, testValue, retUint64)
	})

	// Test 2: echo_u64 with different values
	t.Run("LoopReader_GetLatestValue_EchoU64_VariousValues", func(t *testing.T) {
		testCases := []uint64{0, 1, 100, 1000, 1000000000}

		for _, testValue := range testCases {
			t.Run(fmt.Sprintf("Value_%d", testValue), func(t *testing.T) {
				var retUint64 uint64
				err = loopReader.GetLatestValue(
					context.Background(),
					strings.Join([]string{packageId, echoBinding.Name, "echo_u64"}, "-"),
					primitives.Finalized,
					map[string]any{
						"val": testValue,
					},
					&retUint64,
				)
				require.NoError(t, err)
				require.Equal(t, testValue, retUint64)
			})
		}
	})

	// Test 3: echo_u256 function call
	t.Run("LoopReader_GetLatestValue_EchoU256", func(t *testing.T) {
		t.Skip("Skipping u256 test")

		testValue := big.NewInt(123456789)
		var retBigInt *big.Int
		err = loopReader.GetLatestValue(
			context.Background(),
			strings.Join([]string{packageId, echoBinding.Name, "echo_u256"}, "-"),
			primitives.Finalized,
			map[string]any{
				"val": testValue,
			},
			&retBigInt,
		)
		require.NoError(t, err)
		require.Equal(t, testValue, retBigInt)
	})

	// Test 4: echo_u256 with large values
	t.Run("LoopReader_GetLatestValue_EchoU256_LargeValue", func(t *testing.T) {
		t.Skip("Skipping large value test")

		// Test with a very large number
		testValue := new(big.Int)
		testValue.SetString("123456789012345678901234567890", 10)
		var retBigInt *big.Int
		err = loopReader.GetLatestValue(
			context.Background(),
			strings.Join([]string{packageId, echoBinding.Name, "echo_u256"}, "-"),
			primitives.Finalized,
			map[string]any{
				"val": testValue,
			},
			&retBigInt,
		)
		require.NoError(t, err)
		require.Equal(t, testValue, retBigInt)
	})

	// Test 5: echo_u32_u64_tuple function call
	t.Run("LoopReader_GetLatestValue_EchoTuple", func(t *testing.T) {
		t.Skip("Skipping tuple test")

		testVal1 := uint32(100)
		testVal2 := uint64(200)

		type TupleResult struct {
			First  uint32 `json:"0"`
			Second uint64 `json:"1"`
		}

		var retTuple TupleResult
		err = loopReader.GetLatestValue(
			context.Background(),
			strings.Join([]string{packageId, echoBinding.Name, "echo_u32_u64_tuple"}, "-"),
			primitives.Finalized,
			map[string]any{
				"val1": testVal1,
				"val2": testVal2,
			},
			&retTuple,
		)
		require.NoError(t, err)
		require.Equal(t, testVal1, retTuple.First)
		require.Equal(t, testVal2, retTuple.Second)
	})

	// Test 6: echo_string function call
	t.Run("LoopReader_GetLatestValue_EchoString", func(t *testing.T) {
		t.Skip("Skipping string test")

		testString := "Hello, Sui!"
		var retString string
		err = loopReader.GetLatestValue(
			context.Background(),
			strings.Join([]string{packageId, echoBinding.Name, "echo_string"}, "-"),
			primitives.Finalized,
			map[string]any{
				"val": testString,
			},
			&retString,
		)
		require.NoError(t, err)
		require.Equal(t, testString, retString)
	})
}
