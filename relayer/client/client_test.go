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

func TestClient(t *testing.T) {
	t.Parallel()

	log := logger.Test(t)

	_, err := testutils.StartSuiNode(testutils.CLI)
	require.NoError(t, err)

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

	packageId, _, err := testutils.PublishContract(t, "TestContract", contractPath, accountAddress, nil)
	require.NoError(t, err)

	log.Debugw("Published Contract", "packageId", packageId)

	initializeOutput := testutils.CallContractFromCLI(t, packageId, accountAddress, "counter", "initialize", nil)
	require.NoError(t, err)

	counterObjectId, err := testutils.QueryCreatedObjectID(initializeOutput.ObjectChanges, packageId, "counter", "Counter")
	require.NoError(t, err)

	// Test GetLatestValue for different data types
	t.Run("FunctionRead", func(t *testing.T) {
		t.Parallel()

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
