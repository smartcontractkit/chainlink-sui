package plugin

import (
	"errors"

	"github.com/smartcontractkit/chainlink-common/pkg/config"
)

// Default constant values
const (
	DefaultBroadcastChannelSize       = uint64(4096)
	DefaultMaxConcurrentRequests      = int64(5)
	DefaultRetryCount                 = int64(5)
	DefaultMaxGasLimit                = int64(10000000)
	DefaultTxTimeoutSeconds           = 10
	DefaultConfirmerPoolPeriodSeconds = int64(1)
)

type NodeConfig struct {
	Name *string
	URL  *config.URL
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

	return err
}
