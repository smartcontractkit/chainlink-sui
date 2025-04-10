//go:build integration

package client

import (
	"context"
	"testing"
	"time"

	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/test-go/testify/require"

	"github.com/smartcontractkit/chainlink-sui/relayer/keystore"
	"github.com/smartcontractkit/chainlink-sui/relayer/testutils"
)

//nolint:paralleltest
func TestClient(t *testing.T) {
	log := logger.Test(t)

	cmd, err := testutils.StartSuiNode(testutils.CLI)
	require.NoError(t, err)
	t.Cleanup(func() {
		if cmd.Process != nil {
			perr := cmd.Process.Kill()
			if perr != nil {
				t.Logf("Failed to kill process: %v", perr)
			}
		}
	})

	accountAddress := testutils.GetAccountAndKeyFromSui(t, log)
	keystoreInstance, err := keystore.NewSuiKeystore(log, "", keystore.PrivateKeySigner)
	require.NoError(t, err)
	signer, err := keystoreInstance.GetSignerFromAddress(accountAddress)
	require.NoError(t, err)
	relayerClient, err := NewClient(log, testutils.LocalUrl, nil, 10*time.Second, &signer)
	require.NoError(t, err)

	err = testutils.FundWithFaucet(log, constant.SuiLocalnet, accountAddress)
	require.NoError(t, err)

	contractPath := testutils.BuildSetup(t, "contracts/test")
	testutils.BuildContract(t, contractPath)

	packageId, publishOutput, err := testutils.PublishContract(t, "TestContract", contractPath, accountAddress, nil)
	require.NoError(t, err)

	log.Debugw("Published Contract", "packageId", packageId)

	counterObjectId, err := testutils.QueryCreatedObjectID(publishOutput.ObjectChanges, packageId, "counter", "Counter")
	require.NoError(t, err)

	// Test GetLatestValue for different data types
	//nolint:paralleltest
	t.Run("FunctionRead", func(t *testing.T) {
		args := []any{counterObjectId}
		argTypes := []string{"address"}

		response, err := relayerClient.ReadFunction(
			context.Background(),
			packageId,
			"counter",
			"get_count",
			args,
			argTypes,
		)
		require.NoError(t, err)
		require.NotNil(t, response)
	})
}
