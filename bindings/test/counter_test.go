package counter

import (
	"context"
	"testing"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/stretchr/testify/require"
)

func TestPublishCounter(t *testing.T) {
	t.Parallel()

	// signer := signer.NewSigner()
	// client := sui.NewSuiClient("http://localhost:9000")

	objID, tx, err := PublishCounter(context.Background(), bind.TxOpts{}, *signer, client)
	require.NoError(t, err)

	require.NotNil(t, objID)
	require.NotNil(t, tx)
}
