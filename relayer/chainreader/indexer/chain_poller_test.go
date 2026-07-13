package indexer

import (
	"context"
	"errors"
	"testing"
	"time"

	suirpcv2 "github.com/block-vision/sui-go-sdk/pb/sui/rpc/v2"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/config"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/testutils"
)

type checkpointTestClient struct {
	testutils.FakeSuiPTBClient
	lowestAvailable uint64
	notFoundBelow   uint64
	processed       []uint64
}

func (c *checkpointTestClient) GetCheckpointAvailability(ctx context.Context) (*suirpcv2.GetServiceInfoResponse, error) {
	lowest := c.lowestAvailable
	return &suirpcv2.GetServiceInfoResponse{
		LowestAvailableCheckpoint: &lowest,
	}, nil
}

func (c *checkpointTestClient) GetTransaction(ctx context.Context, digest string) (client.TransactionDetails, error) {
	return client.TransactionDetails{}, nil
}

func (c *checkpointTestClient) GetCheckpointData(ctx context.Context, checkpointSequenceNumber uint64) (*client.CheckpointData, error) {
	if checkpointSequenceNumber < c.notFoundBelow {
		return nil, status.Error(codes.NotFound, "checkpoint not found")
	}

	c.processed = append(c.processed, checkpointSequenceNumber)
	seq := checkpointSequenceNumber
	return &client.CheckpointData{
		Checkpoint: &suirpcv2.Checkpoint{
			SequenceNumber: &seq,
		},
	}, nil
}

func TestClampToProviderFloor(t *testing.T) {
	t.Parallel()

	mockClient := &checkpointTestClient{
		lowestAvailable: 500,
	}
	cp := &ChainPoller{
		client:         mockClient,
		extendedClient: asExtendedPTBClient(mockClient),
		logger:         logger.Test(t),
	}

	require.Equal(t, uint64(500), cp.clampToProviderFloor(context.Background(), 100))
	require.Equal(t, uint64(600), cp.clampToProviderFloor(context.Background(), 600))
}

func TestCatchUpJumpsOverPrunedCheckpoints(t *testing.T) {
	t.Parallel()

	mockClient := &checkpointTestClient{
		lowestAvailable: 500,
		notFoundBelow:   500,
	}

	cp := NewChainPoller(mockClient, logger.Test(t), config.ChainPollerConfig{
		SyncTimeout: time.Minute,
	}, func() []*client.EventSelector {
		return nil
	})

	cp.catchUp(context.Background(), 400, 502)

	require.Equal(t, []uint64{500, 501, 502}, mockClient.processed)
	require.Equal(t, uint64(502), cp.lastProcessed)
}

func TestCatchUpRetriesInRangeNotFound(t *testing.T) {
	t.Parallel()

	mockClient := &checkpointTestClient{
		lowestAvailable: 100,
		notFoundBelow:   501,
	}

	cp := NewChainPoller(mockClient, logger.Test(t), config.ChainPollerConfig{
		SyncTimeout: time.Minute,
	}, func() []*client.EventSelector {
		return nil
	})

	cp.catchUp(context.Background(), 500, 502)

	require.Empty(t, mockClient.processed)
	require.Equal(t, uint64(0), cp.lastProcessed)
}

func TestIsCheckpointNotFound(t *testing.T) {
	t.Parallel()

	require.True(t, isCheckpointNotFound(status.Error(codes.NotFound, "missing")))
	require.True(t, isCheckpointNotFound(errors.New("wrapped: not found")))
	require.False(t, isCheckpointNotFound(errors.New("timeout")))
}
