package indexer

import (
	"errors"
	"testing"

	suirpcv2 "github.com/block-vision/sui-go-sdk/pb/sui/rpc/v2"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/smartcontractkit/chainlink-sui/relayer/client"
)

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

	selector := &client.EventSelector{
		Package: originalPkg,
		Module:  "onramp",
		Event:   "CCIPMessageSent",
	}

	tests := []struct {
		name  string
		event *suirpcv2.Event
		sel   *client.EventSelector
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
