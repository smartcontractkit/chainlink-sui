package config

import (
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/config"
)

// Default constant values
const (
	DefaultBroadcastChannelSize       = uint64(4096)
	DefaultMaxConcurrentRequests      = int64(500)
	DefaultRetryCount                 = int64(5)
	DefaultMaxGasLimit                = int64(10000000)
	DefaultTxTimeoutSeconds           = 120
	DefaultConfirmPollSecs            = int64(10)
	DefaultBalancePollIntervalSeconds = int64(10)

	DefaultIndexerPollIntervalSecs = uint64(10)
	DefaultIndexerSyncTimeoutSecs  = uint64(60)

	DefaultChainPollerPollIntervalSecs        = uint64(2)
	DefaultChainPollerSyncTimeoutSecs         = uint64(60)
	DefaultChainPollerChannelBufferSize       = uint64(16)
	DefaultChainPollerBackfillCheckpointCount = uint64(100)
	DefaultChainPollerMaxConcurrentWorkers    = uint64(8)
	DefaultChainPollerCatchupChunkSize        = uint64(12)
	DefaultChainPollerReplayCheckpointCount   = uint64(100)

	DefaultReaperPollSecs           = uint64(10)
	DefaultTransactionRetentionSecs = uint64(10)
)

type NodeConfig struct {
	Name       *string
	URL        *config.URL
	GrpcTarget *string
	GrpcToken  *string
}

func (n *NodeConfig) ValidateConfig() error {
	var err error
	if n.Name == nil {
		err = errors.Join(err, config.ErrMissing{Name: "Name", Msg: "required for all nodes"})
	} else if *n.Name == "" {
		err = errors.Join(err, config.ErrEmpty{Name: "Name", Msg: "required for all nodes"})
	}
	if n.URL == nil {
		err = errors.Join(err, config.ErrMissing{Name: "URL", Msg: "required for all nodes"})
	}

	// The Sui relayer drives reads, indexing, and tx submission over gRPC, so a node without a
	// gRPC endpoint cannot serve a relayer. Reject it here instead of nil-dereferencing in NewRelayer.
	// The node name is included so that layered multi-node configs identify which entry is missing gRPC.
	nameQualifier := ""
	if n.Name != nil {
		nameQualifier = fmt.Sprintf(" (node %q)", *n.Name)
	}
	grpcMsg := "required for all Sui nodes" + nameQualifier
	if n.GrpcTarget == nil {
		err = errors.Join(err, config.ErrMissing{Name: "GrpcTarget", Msg: grpcMsg})
	} else if *n.GrpcTarget == "" {
		err = errors.Join(err, config.ErrEmpty{Name: "GrpcTarget", Msg: grpcMsg})
	}
	if n.GrpcToken == nil {
		err = errors.Join(err, config.ErrMissing{Name: "GrpcToken", Msg: grpcMsg})
	} else if *n.GrpcToken == "" {
		err = errors.Join(err, config.ErrEmpty{Name: "GrpcToken", Msg: grpcMsg})
	}

	return err
}
