package counter

import (
	"context"
	"testing"

	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/block-vision/sui-go-sdk/signer"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/relayer/testutils"

	"github.com/stretchr/testify/require"
)

func TestCounter(t *testing.T) {
	t.Parallel()
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

	pk, _, _, err := testutils.GenerateAccountKeyPair(t, log)
	require.NoError(t, err)

	signer := signer.NewSigner(pk.Seed())
	client := sui.NewSuiClient("http://localhost:9000")

	err = testutils.FundWithFaucet(log, constant.SuiLocalnet, signer.Address)
	require.NoError(t, err)

	counter, tx, err := PublishCounter(context.Background(), bind.TxOpts{}, *signer, client)
	require.NoError(t, err)

	require.NotNil(t, counter)
	require.NotNil(t, tx)

	counterObjectId, err := bind.FindObjectIdFromPublishTx(tx, "counter")
	require.NoError(t, err)

	// two increments
	_, err = counter.Increment(counterObjectId).Execute(context.Background(), bind.TxOpts{}, *signer, client)
	require.NoError(t, err)
	_, err = counter.Increment(counterObjectId).Execute(context.Background(), bind.TxOpts{}, *signer, client)
	require.NoError(t, err)

	// TODO: Check counter value
}
