// // go:build publictestnet

package client

import (
	"context"
	"testing"
	"time"

	"github.com/block-vision/sui-go-sdk/utils"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

const ccipTestnetTimelockObjectID = "0xa514be3fe446f654389c1bd2dc4ce9dcbd85753fe537c0c64a34298607ee33b6"
const accountAddress = "0x3aacfda9282684c7013df61243f684afb9c602d14965232ee87cdcc76a2f4c45"
const coinObjectID = "0x05507095c89af9b38bbed34d5805a24ef1cf4f7dff6a529951c4edca2a371591"

func TestPublicTestnetGrpcReadObject(t *testing.T) {
	t.Parallel()

	log := logger.Test(t)
	cfg := PTBClientConfig{
		GrpcTarget:            "sui-testnet.g.alchemy.com:443",
		GrpcToken:             "lBOT0DAq5sJlkW0BFbNqG",
		TransactionTimeout:    30 * time.Second,
		MaxConcurrentRequests: 10,
		MaxGrpcConnections:    1,
	}

	ptbClient, err := NewPTBClient(log, cfg)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// obj, err := ptbClient.ReadObjectId(ctx, ccipTestnetTimelockObjectID)
	// require.NoError(t, err)
	// require.NotNil(t, obj)
	// require.Equal(t, ccipTestnetTimelockObjectID, obj.GetObjectId())

	// check balance
	coins, err := ptbClient.QueryCoinsByAddress(ctx, accountAddress, "0x2::coin::Coin<0x2::sui::SUI>")
	require.NoError(t, err)
	require.NotNil(t, coins)

	utils.PrettyPrint(coins)

	// get object details for coin
	details, err := ptbClient.ReadObjectId(ctx, coinObjectID)
	require.NoError(t, err)
	utils.PrettyPrint(details)
}
