//go:build integration

package counter

import (
	"context"
	"testing"

	"github.com/pattonkan/sui-go/suiclient"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	rel "github.com/smartcontractkit/chainlink-sui/relayer/signer"
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

	signer := rel.NewPrivateKeySigner(pk)
	client := suiclient.NewClient("http://localhost:9000")

	signerAddress, err := signer.GetAddress()
	require.NoError(t, err)

	err = testutils.FundWithFaucet(log, "localnet", signerAddress)
	require.NoError(t, err)

	ctx := context.Background()

	counter, tx, err := PublishCounter(ctx, bind.TxOpts{}, signer, *client)
	require.NoError(t, err)

	require.NotNil(t, counter)
	require.NotNil(t, tx)

	counterObjectId, err := bind.FindObjectIdFromPublishTx(tx, "counter", "Counter")
	require.NoError(t, err)

	// two increments
	increment := counter.Increment(counterObjectId)
	_, err = increment.Execute(ctx, bind.TxOpts{}, signer, *client)
	require.NoError(t, err)
	_, err = increment.Execute(ctx, bind.TxOpts{}, signer, *client)
	require.NoError(t, err)

	_, err = counter.GetCount(counterObjectId).Inspect((ctx), bind.TxOpts{}, signer, *client)
	require.NoError(t, err)

	value, err := counter.Inspect(ctx, counterObjectId)
	require.NoError(t, err)
	require.Equal(t, uint64(2), value)
}
