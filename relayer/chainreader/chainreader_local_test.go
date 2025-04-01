//go:build integration

package chainreader

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
	"github.com/smartcontractkit/chainlink-sui/relayer/keystore"

	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"

	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/testutils"
)

//nolint:paralleltest
func TestChainReaderLocal(t *testing.T) {
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

	err = testutils.FundWithFaucet(log, constant.SuiLocalnet, accountAddress)
	require.NoError(t, err)

	runChainReaderCounterTest(t, log, testutils.LocalUrl)
}

//nolint:paralleltest
func runChainReaderCounterTest(t *testing.T, log logger.Logger, rpcUrl string) {
	t.Helper()

	accountAddress := testutils.GetAccountAndKeyFromSui(t, log)
	keystoreInstance, err := keystore.NewSuiKeystore(log, "", keystore.PrivateKeySigner)
	signer, err := keystoreInstance.GetSignerFromAddress(accountAddress)
	relayerClient, err := client.NewClient(log, rpcUrl, nil, 10*time.Second, &signer)
	require.NoError(t, err)

	err = testutils.FundWithFaucet(log, constant.SuiLocalnet, accountAddress)
	require.NoError(t, err)

	contractPath := testutils.BuildSetup(t, "contracts/test")
	testutils.BuildContract(t, contractPath)

	packageId, _, err := testutils.PublishContract(t, "TestContract", contractPath, accountAddress, nil)
	require.NoError(t, err)

	log.Debugw("Published Contract", "packageId", packageId)

	initializeOutput := testutils.CallContractFromCLI(t, packageId, accountAddress, "counter", "initialize", nil)
	require.NoError(t, err)

	counterObjectId, err := testutils.QueryCreatedObjectID(initializeOutput.ObjectChanges, packageId, "counter", "Counter")
	require.NoError(t, err)

	// Set up the ChainReader
	chainReaderConfig := ChainReaderConfig{
		Modules: map[string]*ChainReaderModule{
			"counter": {
				Name: "counter",
				Functions: map[string]*ChainReaderFunction{
					"get_count": {
						Name: "get_count",
						Params: []codec.SuiFunctionParam{
							{
								Type:         "address",
								Name:         "counter_id",
								DefaultValue: counterObjectId,
								Required:     true,
							},
						},
					},
				},
			},
		},
	}

	counterBinding := types.BoundContract{
		Name:    "counter",
		Address: packageId, // Package ID of the deployed counter contract
	}

	chainReader := NewChainReader(log, *relayerClient, chainReaderConfig)
	err = chainReader.Bind(context.Background(), []types.BoundContract{counterBinding})
	require.NoError(t, err)

	log.Debugw("ChainReader setup complete")

	// Test GetLatestValue for different data types
	//nolint:paralleltest
	t.Run("GetLatestValue_Uint64", func(t *testing.T) {
		t.Parallel()
		expectedUint64 := uint64(0)
		var retUint64 uint64
		err = chainReader.GetLatestValue(
			context.Background(),
			strings.Join([]string{packageId, counterBinding.Name, counterObjectId}, "-"),
			primitives.Finalized,
			struct {
				Value uint64
			}{Value: expectedUint64},
			&retUint64,
		)
		require.NoError(t, err)
		require.Equal(t, expectedUint64, retUint64)
	})

	//nolint:paralleltest
	t.Run("GetLatestValue_FunctionRead", func(t *testing.T) {
		t.Parallel()
		expectedUint64 := uint64(0)
		var retUint64 uint64

		log.Debugw("Testing get_count",
			"counterObjectId", counterObjectId,
			"packageId", packageId,
		)

		err = chainReader.GetLatestValue(
			context.Background(),
			strings.Join([]string{packageId, counterBinding.Name, "get_count"}, "-"),
			primitives.Finalized,
			map[string]interface{}{
				"counter_id": counterObjectId,
			},
			&retUint64,
		)
		require.NoError(t, err)
		require.Equal(t, expectedUint64, retUint64)
	})
}
