//go:build integration

package offramp_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-sui/relayer/chainwriter/ptb/offramp"
	"github.com/smartcontractkit/chainlink-sui/relayer/testutils"
)

// TestValidateReceiverObjectOwner_Integration exercises the transmitter-ownership guard against
// a real disposable local Sui validator. The unit test TestValidateReceiverObjectOwner covers the
// decision logic with synthetic owners; this test verifies the contract that real ReadObjectMetadata
// returns Owner.Kind == ADDRESS + Owner.Address for a coin owned by the signer, so the guard
// actually rejects transmitter-owned tail objects and allows everything else.
func TestValidateReceiverObjectOwner_Integration(t *testing.T) {
	lggr := logger.Test(t)
	gasBudget := int64(1_000_000_000)

	// Disposable local validator only.
	cmd, err := testutils.StartSuiNode(testutils.CLI)
	require.NoError(t, err)
	t.Cleanup(func() {
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
		}
	})

	ctx := context.Background()

	// The transmitter: the funded test account that signs / owns gas in CCIP execute.
	keystore, transmitter, _ := testutils.SetupTestSigner(t, ctx, lggr, gasBudget)
	lggr.Infow("Using transmitter account", "address", transmitter)

	require.Eventually(t, func() bool {
		return testutils.FundWithFaucet(lggr, testutils.SuiLocalnet, transmitter) == nil
	}, 10*time.Second, time.Second)

	ptbClient, _, _ := testutils.SetupClients(t, testutils.LocalGrpcURL, keystore, lggr, gasBudget)

	// A coin owned by the transmitter — the exploit target a malicious receiver would name as a
	// tail object with &mut.
	transmitterCoins, err := ptbClient.GetCoinsByAddress(ctx, transmitter)
	require.NoError(t, err)
	require.NotEmpty(t, transmitterCoins, "transmitter must own at least one coin after funding")
	transmitterCoin := transmitterCoins[0].GetObjectId()

	// A distinct account owns a coin — must be allowed through.
	_, _, otherAddr, err := testutils.GenerateAccountKeyPair(t)
	require.NoError(t, err)
	require.Eventually(t, func() bool {
		return testutils.FundWithFaucet(lggr, testutils.SuiLocalnet, otherAddr) == nil
	}, 10*time.Second, time.Second)
	otherCoins, err := ptbClient.GetCoinsByAddress(ctx, otherAddr)
	require.NoError(t, err)
	require.NotEmpty(t, otherCoins, "other account must own at least one coin after funding")
	otherCoin := otherCoins[0].GetObjectId()

	// 0x2 is the Sui framework package — immutable, reliably readable on any network.
	const immutableObject = "0x2"

	// Transmitter-owned address object → rejected (the exploit case, via real RPC).
	require.ErrorIs(t,
		offramp.ValidateObjectOwner(ctx, ptbClient, transmitterCoin, transmitter),
		offramp.ErrTransmitterOwnedReceiverObject,
	)

	// Address-owned by a different account → allowed.
	require.NoError(t,
		offramp.ValidateObjectOwner(ctx, ptbClient, otherCoin, transmitter),
	)

	// Immutable object → allowed.
	require.NoError(t,
		offramp.ValidateObjectOwner(ctx, ptbClient, immutableObject, transmitter),
	)
}
