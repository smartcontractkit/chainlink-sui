//go:build publictestnet

package client

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

const ccipTestnetTimelockObjectID = "0xa514be3fe446f654389c1bd2dc4ce9dcbd85753fe537c0c64a34298607ee33b6"

func TestPublicTestnetGrpcReadObject(t *testing.T) {
	t.Parallel()

	log := logger.Test(t)
	cfg := PTBClientConfig{
		GrpcTarget:            "fullnode.testnet.sui.io:443",
		GrpcToken:             "test",
		TransactionTimeout:    30 * time.Second,
		MaxConcurrentRequests: 10,
		MaxGrpcConnections:    1,
	}

	ptbClient, err := NewPTBClient(log, cfg)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	obj, err := ptbClient.ReadObjectId(ctx, ccipTestnetTimelockObjectID)
	require.NoError(t, err)
	require.NotNil(t, obj)
	require.Equal(t, ccipTestnetTimelockObjectID, obj.GetObjectId())
}
