package plugin

import (
	"errors"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/config"

	"github.com/smartcontractkit/chainlink-sui/relayer/client"
)

// Default constant values
const (
	DefaultBroadcastChannelSize  = uint64(4096)
	DefaultMaxConcurrentRequests = int64(5)
	DefaultRetryCount            = int64(5)
	DefaultMaxGasLimit           = int64(10000000)
	DefaultConfirmPollPeriodMs   = 500
	DefaultTxTimeoutSeconds      = 10
)

// Global Sui defaults.
var defaultBroadcastChanSize = DefaultBroadcastChannelSize
var defaultConfirmPollPeriod = time.Duration(DefaultConfirmPollPeriodMs) * time.Millisecond
var defaultMaxConcurrentRequests = DefaultMaxConcurrentRequests
var defaultTxTimeout = time.Duration(DefaultTxTimeoutSeconds) * time.Second
var defaultRetryCount = DefaultRetryCount
var defaultMaxGasLimit = DefaultMaxGasLimit
var defaultRequestType = client.WaitForLocalExecution

type ChainConfig struct {
	BroadcastChanSize     *uint64
	ConfirmPollPeriod     *time.Duration
	MaxConcurrentRequests *int64
	TransactionTimeout    *time.Duration
	NumberRetries         *int64
	GasLimit              *int64
	RequestType           *client.TransactionRequestType
}

func (c *ChainConfig) Defaults() {
	if c.BroadcastChanSize == nil {
		c.BroadcastChanSize = &defaultBroadcastChanSize
	}
	if c.ConfirmPollPeriod == nil {
		c.ConfirmPollPeriod = &defaultConfirmPollPeriod
	}
	if c.MaxConcurrentRequests == nil {
		c.MaxConcurrentRequests = &defaultMaxConcurrentRequests
	}
	if c.TransactionTimeout == nil {
		c.TransactionTimeout = &defaultTxTimeout
	}
	if c.NumberRetries == nil {
		c.NumberRetries = &defaultRetryCount
	}
	if c.GasLimit == nil {
		c.GasLimit = &defaultMaxGasLimit
	}
	if c.RequestType == nil {
		c.RequestType = &defaultRequestType
	}
}

type NodeConfig struct {
	Name        *string
	URL         *config.URL
	SolidityURL *config.URL
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
	if n.SolidityURL == nil {
		err = errors.Join(err, config.ErrMissing{Name: "SolidityURL", Msg: "required for all nodes"})
	}

	return err
}
