package mcmstest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
)

const (
	PackageID     = "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	StateObjectID = "0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	OwnerCapID    = "0xcccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc"
	CoinMetadata  = "0xdddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd"
	RegistryID    = "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
	Recipient     = "0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"
	CoinTypeArg   = "0x2::sui::SUI"
	DestChainSel  = uint64(16015286601757825753)
)

func Bundle(t *testing.T) cld_ops.Bundle {
	t.Helper()
	return cld_ops.NewBundle(
		func() context.Context { return t.Context() },
		logger.Test(t),
		cld_ops.NewMemoryReporter(),
	)
}

func AssertProposalDataMatches(t *testing.T, got []byte, encodedCall *bind.EncodedCall, stateObjID string, typeArgs []string) {
	t.Helper()
	var (
		expected sui_ops.TransactionCall
		err      error
	)
	if len(typeArgs) > 0 {
		expected, err = sui_ops.ToTransactionCallWithTypeArgs(encodedCall, stateObjID, typeArgs)
	} else {
		expected, err = sui_ops.ToTransactionCall(encodedCall, stateObjID)
	}
	require.NoError(t, err)
	require.Equal(t, expected.Data, got)
}
