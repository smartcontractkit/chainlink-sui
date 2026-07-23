package indexer

import (
	"context"
	"errors"
	"slices"
	"sync"
	"testing"
	"time"

	suirpcv2 "github.com/block-vision/sui-go-sdk/pb/sui/rpc/v2"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types/sui"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/testutils"
)

type checkpointTestClient struct {
	testutils.FakeSuiPTBClient
	lowestAvailable uint64
	notFoundBelow   uint64
	mu              sync.Mutex
	processed       []uint64
	// floorSeq, when set, returns successive floor values from GetCheckpointAvailability
	// (clamped to the last value once exhausted) to simulate the provider floor advancing
	// mid-catch-up.
	floorSeq []uint64
	floorIdx int
	// pruned, when set, marks specific sequence numbers as NotFound (mid-range pruning),
	// instead of the notFoundBelow bottom-pruning model.
	pruned map[uint64]bool
}

func (c *checkpointTestClient) GetCheckpointAvailability(ctx context.Context) (*suirpcv2.GetServiceInfoResponse, error) {
	c.mu.Lock()
	lowest := c.lowestAvailable
	if len(c.floorSeq) > 0 {
		idx := c.floorIdx
		if idx >= len(c.floorSeq) {
			idx = len(c.floorSeq) - 1
		}
		lowest = c.floorSeq[idx]
		c.floorIdx++
	}
	c.mu.Unlock()
	return &suirpcv2.GetServiceInfoResponse{
		LowestAvailableCheckpoint: &lowest,
	}, nil
}

func (c *checkpointTestClient) GetTransaction(ctx context.Context, digest string) (client.TransactionDetails, error) {
	return client.TransactionDetails{}, nil
}

func (c *checkpointTestClient) GetCheckpointData(ctx context.Context, checkpointSequenceNumber uint64) (*client.CheckpointData, error) {
	notFound := false
	if c.pruned != nil {
		notFound = c.pruned[checkpointSequenceNumber]
	} else if checkpointSequenceNumber < c.notFoundBelow {
		notFound = true
	}
	if notFound {
		return nil, status.Error(codes.NotFound, "checkpoint not found")
	}

	c.mu.Lock()
	c.processed = append(c.processed, checkpointSequenceNumber)
	c.mu.Unlock()

	seq := checkpointSequenceNumber
	return &client.CheckpointData{
		Checkpoint: &suirpcv2.Checkpoint{
			SequenceNumber: &seq,
		},
	}, nil
}

// getProcessed returns a copy of the fetched sequence numbers. catchUp fetches
// concurrently, so the order is non-deterministic; callers should sort before
// comparing.
func (c *checkpointTestClient) getProcessed() []uint64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make([]uint64, len(c.processed))
	copy(out, c.processed)
	return out
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

	cp := NewChainPoller(mockClient, logger.Test(t), sui.ChainPollerConfig{
		SyncTimeout: time.Minute,
	}, func() []*sui.EventFilterByMoveEventModule {
		return nil
	})

	cp.catchUp(context.Background(), 400, 502)

	// Fetches run concurrently, so compare the set rather than the order.
	processed := mockClient.getProcessed()
	slices.Sort(processed)
	require.Equal(t, []uint64{500, 501, 502}, processed)
	require.Equal(t, uint64(502), cp.lastProcessed)
}

func TestCatchUpRetriesInRangeNotFound(t *testing.T) {
	t.Parallel()

	mockClient := &checkpointTestClient{
		lowestAvailable: 100,
		notFoundBelow:   501,
	}

	cp := NewChainPoller(mockClient, logger.Test(t), sui.ChainPollerConfig{
		SyncTimeout: time.Minute,
	}, func() []*sui.EventFilterByMoveEventModule {
		return nil
	})

	cp.catchUp(context.Background(), 500, 502)

	// 500 is not found (and not pruned, not the tip), so catchUp stops without
	// advancing lastProcessed. With concurrent fetches, checkpoints ahead of 500
	// may be fetched before the commit stage stops, but none are committed past
	// the gap; the safety property is that lastProcessed stays put and 500 is
	// retried on the next poll.
	require.NotContains(t, mockClient.getProcessed(), uint64(500))
	require.Equal(t, uint64(0), cp.lastProcessed)
}

func TestCatchUpStopsOnMidCatchupFloorAdvance(t *testing.T) {
	t.Parallel()

	// floorSeq: the top clamp sees floor 100 (no clamp), then the jump sees 506, so the
	// floor "advances" mid-catch-up. pruned: 503-505 are unavailable; 500-502 and 506+
	// are available.
	mockClient := &checkpointTestClient{
		floorSeq: []uint64{100, 506},
		pruned:   map[uint64]bool{503: true, 504: true, 505: true},
	}

	cp := NewChainPoller(mockClient, logger.Test(t), sui.ChainPollerConfig{
		SyncTimeout: time.Minute,
	}, func() []*sui.EventFilterByMoveEventModule {
		return nil
	})

	// 500-502 commit; 503 is pruned and the floor has advanced to 506, so catchUp stops
	// at 502 instead of fetching/buffering the pruned 503-505.
	cp.catchUp(context.Background(), 500, 510)
	require.Equal(t, uint64(502), cp.lastProcessed)

	// The next poll resumes from lastProcessed+1=503; clampToProviderFloor raises it to
	// the current floor (506), so 506-510 commit and no available checkpoint is missed.
	// lastProcessed reaching 510 implies 506-510 all committed in order.
	cp.catchUp(context.Background(), 503, 510)
	require.Equal(t, uint64(510), cp.lastProcessed)
}

func TestIsCheckpointNotFound(t *testing.T) {
	t.Parallel()

	require.True(t, isCheckpointNotFound(status.Error(codes.NotFound, "missing")))
	require.True(t, isCheckpointNotFound(errors.New("wrapped: not found")))
	require.False(t, isCheckpointNotFound(errors.New("timeout")))
}
