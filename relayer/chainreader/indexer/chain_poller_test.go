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
	"github.com/smartcontractkit/chainlink-common/pkg/types/sui"
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

	cp := NewChainPoller(mockClient, logger.Test(t), sui.ChainPollerConfig{
		SyncTimeout: time.Minute,
	}, func() []*sui.EventFilterByMoveEventModule {
		return nil
	})

	cp.catchUp(context.Background(), 400, 502)

	require.Equal(t, []uint64{500, 501, 502}, mockClient.processed)
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

	require.Empty(t, mockClient.processed)
}

func TestIsCheckpointNotFound(t *testing.T) {
	t.Parallel()

	require.True(t, isCheckpointNotFound(status.Error(codes.NotFound, "missing")))
	require.True(t, isCheckpointNotFound(errors.New("wrapped: not found")))
	require.False(t, isCheckpointNotFound(errors.New("timeout")))
}

func strPtr(s string) *string { return &s }

func TestEventMatchesSelector(t *testing.T) {
	t.Parallel()

	const (
		originalPkg = "0x30e087460af8a8aacccbc218aa358cdcde8d43faf61ec0638d71108e276e2f1d"
		latestPkg   = "0xfa4dc9ef5e099b6dc61c90b00e2b28a90b788fda510790bae84c96d2f0b0303c"
	)

	selector := &sui.EventFilterByMoveEventModule{
		Package: originalPkg,
		Module:  "onramp",
		Event:   "CCIPMessageSent",
	}

	tests := []struct {
		name  string
		event *suirpcv2.Event
		sel   *sui.EventFilterByMoveEventModule
		want  bool
	}{
		{
			name: "upgraded package: emitting package is latest, type string carries original",
			event: &suirpcv2.Event{
				PackageId: strPtr(latestPkg),
				Module:    strPtr("onramp"),
				EventType: strPtr(originalPkg + "::onramp::CCIPMessageSent"),
			},
			sel:  selector,
			want: true,
		},
		{
			name: "exact match, no upgrade",
			event: &suirpcv2.Event{
				PackageId: strPtr(originalPkg),
				Module:    strPtr("onramp"),
				EventType: strPtr(originalPkg + "::onramp::CCIPMessageSent"),
			},
			sel:  selector,
			want: true,
		},
		{
			name: "package mismatch in type string",
			event: &suirpcv2.Event{
				PackageId: strPtr(latestPkg),
				Module:    strPtr("onramp"),
				EventType: strPtr(latestPkg + "::onramp::CCIPMessageSent"),
			},
			sel:  selector,
			want: false,
		},
		{
			name: "module mismatch",
			event: &suirpcv2.Event{
				PackageId: strPtr(originalPkg),
				Module:    strPtr("offramp"),
				EventType: strPtr(originalPkg + "::offramp::CCIPMessageSent"),
			},
			sel:  selector,
			want: false,
		},
		{
			name: "event name mismatch",
			event: &suirpcv2.Event{
				PackageId: strPtr(originalPkg),
				Module:    strPtr("onramp"),
				EventType: strPtr(originalPkg + "::onramp::ExecutionStateChanged"),
			},
			sel:  selector,
			want: false,
		},
		{
			name: "malformed event type with fewer than three segments",
			event: &suirpcv2.Event{
				PackageId: strPtr(originalPkg),
				Module:    strPtr("onramp"),
				EventType: strPtr(originalPkg + "::onramp"),
			},
			sel:  selector,
			want: false,
		},
		{
			name:  "nil event",
			event: nil,
			sel:   selector,
			want:  false,
		},
		{
			name: "nil selector",
			event: &suirpcv2.Event{
				EventType: strPtr(originalPkg + "::onramp::CCIPMessageSent"),
			},
			sel:  nil,
			want: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tc.want, eventMatchesSelector(tc.event, tc.sel))
		})
	}
}
