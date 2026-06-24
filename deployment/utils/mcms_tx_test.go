package utils_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	"github.com/smartcontractkit/chainlink-sui/deployment/utils"
)

func TestTransactionCallToMCMSTransaction_LatestPackageID(t *testing.T) {
	t.Parallel()

	const (
		originalPackageID = "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
		latestPackageID   = "0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	)

	call := sui_ops.TransactionCall{
		PackageID:       originalPackageID,
		LatestPackageID: latestPackageID,
		Module:          "offramp",
		Function:        "apply_source_chain_config_updates",
		Data:            []byte{0x01},
		StateObjID:      "0xcccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc",
	}

	tx, err := utils.TransactionCallToMCMSTransaction(call)
	require.NoError(t, err)
	require.Equal(t, originalPackageID, tx.To)
	latestPackageIDFromTx, err := utils.TransactionLatestPackageID(tx)
	require.NoError(t, err)
	require.Equal(t, latestPackageID, latestPackageIDFromTx)
}
