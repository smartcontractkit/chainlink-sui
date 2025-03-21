//go:build integration

package chainreader

import (
	"context"
	"strings"
	"testing"

	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"

	"github.com/smartcontractkit/chainlink-sui/relayer/testutils"
)

func TestChainReaderLocal(t *testing.T) {
	t.Parallel()

	log := logger.Test(t)

	var err error
	privateKey, _, accountAddress := testutils.LoadAccountFromEnv(t, log)
	// if the env does not contain a private key to be loaded, create one
	if privateKey == nil {
		_, _, accountAddress, err = testutils.GenerateAccountKeyPair(t, log)
	}
	require.NoError(t, err)

	err = testutils.StartSuiNode(testutils.CLI)
	require.NoError(t, err)
	log.Debugw("Started Sui node")

	err = testutils.FundWithFaucet(log, constant.SuiLocalnet, accountAddress)
	require.NoError(t, err)

	runChainReaderCounterTest(t, log, testutils.LocalUrl)
}

func runChainReaderCounterTest(t *testing.T, log logger.Logger, rpcUrl string) {
	t.Helper()

	client := sui.NewSuiClient(rpcUrl)

	// start by deploying the counter contract to local net
	packageId, counterObjectId, err := testutils.DeployCounterContract(t)
	require.NoError(t, err)

	// Set up the ChainReader
	chainReaderConfig := ChainReaderConfig{
		Modules: map[string]*ChainReaderModule{
			"counter": {
				Name: "counter",
			},
		},
	}

	counterBinding := types.BoundContract{
		Name:    "counter",
		Address: packageId, // Package ID of the deployed counter contract
	}

	chainReader := NewChainReader(log, client, chainReaderConfig)
	err = chainReader.Bind(context.Background(), []types.BoundContract{counterBinding})
	require.NoError(t, err)

	log.Debugw("ChainReader setup complete")

	// Test GetLatestValue for different data types
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
}
