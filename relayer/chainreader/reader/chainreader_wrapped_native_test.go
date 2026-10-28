package reader

import (
	"context"
	"database/sql"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	pkgtypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
	"github.com/smartcontractkit/chainlink-common/pkg/types/sui"

	codecLoop "github.com/smartcontractkit/chainlink-sui/codec/loop"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
)

const wrappedNativeTestKey = "0x9258181f5ceac8dbffb7030890243caed69a9599d2886d957a9cb7656af3bdb3"

var (
	wrappedNativeRouterAddress = "0x" + strings.Repeat("ab", 32)
	wrappedNativeReadID        = wrappedNativeRouterAddress + "-Router-GetWrappedNative"
)

// testDataSource satisfies sqlutil.DataSource; NewChainReader's only database work on this path is schema DDL.
type testDataSource struct{ sqlutil.DataSource }

func (testDataSource) ExecContext(context.Context, string, ...any) (sql.Result, error) {
	return nil, nil
}

func newWrappedNativeTestReader(t *testing.T, isLoop bool, wrappedNativeToken string) pkgtypes.ContractReader {
	t.Helper()

	ptbClient, err := client.NewPTBClient(logger.Test(t), client.PTBClientConfig{})
	require.NoError(t, err)

	cr, err := NewChainReader(
		t.Context(),
		logger.Test(t),
		ptbClient,
		sui.ChainReaderConfig{IsLoopPlugin: isLoop},
		testDataSource{},
		nil,
		nil,
		wrappedNativeToken,
	)
	require.NoError(t, err)

	return cr
}

func bindWrappedNativeRouter(t *testing.T, cr pkgtypes.ContractReader) {
	t.Helper()

	require.NoError(t, cr.Bind(t.Context(), []pkgtypes.BoundContract{{Name: "Router", Address: wrappedNativeRouterAddress}}))
}

func TestGetLatestValue_RouterWrappedNative(t *testing.T) {
	t.Parallel()

	wantKey, err := hex.DecodeString(strings.TrimPrefix(wrappedNativeTestKey, "0x"))
	require.NoError(t, err)
	require.Len(t, wantKey, 32)

	t.Run("returns the SUI token key bytes", func(t *testing.T) {
		cr := newWrappedNativeTestReader(t, false, wrappedNativeTestKey)
		bindWrappedNativeRouter(t, cr)

		var got []byte
		require.NoError(t, cr.GetLatestValue(t.Context(), wrappedNativeReadID, primitives.Finalized, nil, &got))
		require.Equal(t, wantKey, got)
	})

	t.Run("returns the key as a JSON string in loop plugin mode", func(t *testing.T) {
		cr := newWrappedNativeTestReader(t, true, wrappedNativeTestKey)
		bindWrappedNativeRouter(t, cr)

		var got []byte
		require.NoError(t, cr.GetLatestValue(t.Context(), wrappedNativeReadID, primitives.Finalized, nil, &got))
		require.Equal(t, `"`+wrappedNativeTestKey+`"`, string(got))
	})

	t.Run("decodes through the core loop wrapper", func(t *testing.T) {
		inner := newWrappedNativeTestReader(t, true, wrappedNativeTestKey)
		wrapped := codecLoop.NewLoopChainReader(logger.Test(t), inner)
		bindWrappedNativeRouter(t, wrapped)

		var got []byte
		require.NoError(t, wrapped.GetLatestValue(t.Context(), wrappedNativeReadID, primitives.Finalized, nil, &got))
		require.Equal(t, wantKey, got)
	})

	t.Run("fails when the Router is not bound", func(t *testing.T) {
		cr := newWrappedNativeTestReader(t, false, wrappedNativeTestKey)

		var got []byte
		err := cr.GetLatestValue(t.Context(), wrappedNativeReadID, primitives.Finalized, nil, &got)
		require.ErrorContains(t, err, "invalid binding for contract: Router")
	})

	t.Run("keeps the missing config error without a SUI key", func(t *testing.T) {
		cr := newWrappedNativeTestReader(t, false, "")
		bindWrappedNativeRouter(t, cr)

		var got []byte
		err := cr.GetLatestValue(t.Context(), wrappedNativeReadID, primitives.Finalized, nil, &got)
		require.ErrorContains(t, err, "no configuration for contract: Router")
	})
}
