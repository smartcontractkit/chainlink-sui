package config

import (
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"golang.org/x/exp/slices"

	"github.com/smartcontractkit/chainlink-common/pkg/config"

	"github.com/smartcontractkit/chainlink-sui/relayer/client"
)

type TOMLConfigs []*TOMLConfig

// decodeConfig decodes the rawConfig as (Aptos) TOML and sets default values
func NewDecodedTOMLConfig(rawConfig string) (*TOMLConfig, error) {
	d := toml.NewDecoder(strings.NewReader(rawConfig))
	d.DisallowUnknownFields()

	var cfg TOMLConfig
	if err := d.Decode(&cfg); err != nil {
		return &TOMLConfig{}, fmt.Errorf("failed to decode config toml: %w:\n\t%s", err, rawConfig)
	}

	if err := cfg.ValidateConfig(); err != nil {
		return &TOMLConfig{}, fmt.Errorf("invalid sui config: %w", err)
	}

	if !cfg.IsEnabled() {
		return &TOMLConfig{}, fmt.Errorf("cannot create new chain with ID %v: config is disabled", cfg.ChainID)
	}

	if cfg.TransactionManager == nil {
		cfg.TransactionManager = &TransactionManagerConfig{}
	}
	cfg.TransactionManager.setDefaults()

	if cfg.BalanceMonitor == nil {
		cfg.BalanceMonitor = &BalanceMonitorConfig{}
	}
	cfg.BalanceMonitor.setDefaults()

	if cfg.TransactionsIndexer == nil {
		cfg.TransactionsIndexer = &IndexerConfig{}
	}
	cfg.TransactionsIndexer.setDefaults()

	if cfg.EventsIndexer == nil {
		cfg.EventsIndexer = &IndexerConfig{}
	}
	cfg.EventsIndexer.setDefaults()

	return &cfg, nil
}

func (cs TOMLConfigs) ValidateConfig() error {
	return cs.validateKeys()
}

func (cs TOMLConfigs) validateKeys() error {
	var err error
	// Unique chain IDs
	chainIDs := config.UniqueStrings{}
	for i, c := range cs {
		if chainIDs.IsDupe(c.ChainID) {
			err = errors.Join(err, config.NewErrDuplicate(fmt.Sprintf("%d.ChainID", i), *c.ChainID))
		}
	}

	// Unique node names
	names := config.UniqueStrings{}
	for i, c := range cs {
		for j, n := range c.Nodes {
			if names.IsDupe(n.Name) {
				err = errors.Join(err, config.NewErrDuplicate(fmt.Sprintf("%d.Nodes.%d.Name", i, j), *n.Name))
			}
		}
	}

	// Unique URLs
	urls := config.UniqueStrings{}
	for i, c := range cs {
		for j, n := range c.Nodes {
			u := (*url.URL)(n.URL)
			if urls.IsDupeFmt(u) {
				err = errors.Join(err, config.NewErrDuplicate(fmt.Sprintf("%d.Nodes.%d.URL", i, j), u.String()))
			}
		}
	}

	return err
}

func (cs *TOMLConfigs) SetFrom(fs *TOMLConfigs) error {
	if err1 := fs.validateKeys(); err1 != nil {
		return err1
	}
	for _, f := range *fs {
		if f.ChainID == nil {
			*cs = append(*cs, f)
		} else if i := slices.IndexFunc(*cs, func(c *TOMLConfig) bool {
			return c.ChainID != nil && *c.ChainID == *f.ChainID
		}); i == -1 {
			*cs = append(*cs, f)
		} else {
			(*cs)[i].SetFrom(f)
		}
	}

	return nil
}

type NodeConfigs []*NodeConfig

func (ns *NodeConfigs) SetFrom(fs *NodeConfigs) {
	for _, f := range *fs {
		if f.Name == nil {
			*ns = append(*ns, f)
		} else if i := slices.IndexFunc(*ns, func(n *NodeConfig) bool {
			return n.Name != nil && *n.Name == *f.Name
		}); i == -1 {
			*ns = append(*ns, f)
		} else {
			setFromNode((*ns)[i], f)
		}
	}
}

func (ns NodeConfigs) SelectRandom() (*NodeConfig, error) {
	if len(ns) == 0 {
		return nil, errors.New("no nodes available")
	}

	index := rand.Perm(len(ns))
	node := ns[index[0]]

	return node, nil
}

func setFromNode(n, f *NodeConfig) {
	if f.Name != nil {
		n.Name = f.Name
	}
	if f.URL != nil {
		n.URL = f.URL
	}
}

type TransactionManagerConfig struct {
	BroadcastChanSize     *uint64
	ConfirmPollSecs       *uint64
	DefaultMaxGasAmount   *uint64
	MaxTxRetryAttempts    *uint64
	RequestType           *string
	TransactionTimeout    *string
	MaxConcurrentRequests *uint64
}

type IndexerConfig struct {
	PollingIntervalSecs *uint64
	SyncTimeoutSecs     *uint64
}

func (i *IndexerConfig) setDefaults() {
	if i.PollingIntervalSecs == nil {
		v := DefaultIndexerPollIntervalSecs
		i.PollingIntervalSecs = &v
	}
	if i.SyncTimeoutSecs == nil {
		v := DefaultIndexerSyncTimeoutSecs
		i.SyncTimeoutSecs = &v
	}
}

func (t *TransactionManagerConfig) setDefaults() {
	if t.BroadcastChanSize == nil {
		defaultVal := DefaultBroadcastChannelSize
		t.BroadcastChanSize = &defaultVal
	}
	if t.MaxConcurrentRequests == nil {
		defaultVal := uint64(DefaultMaxConcurrentRequests)
		t.MaxConcurrentRequests = &defaultVal
	}
	if t.MaxTxRetryAttempts == nil {
		defaultVal := uint64(DefaultRetryCount)
		t.MaxTxRetryAttempts = &defaultVal
	}
	if t.DefaultMaxGasAmount == nil {
		defaultVal := uint64(DefaultMaxGasLimit)
		t.DefaultMaxGasAmount = &defaultVal
	}
	if t.TransactionTimeout == nil {
		defaultVal := fmt.Sprintf("%ds", DefaultTxTimeoutSeconds)
		t.TransactionTimeout = &defaultVal
	}
	if t.RequestType == nil {
		defaultVal := string(client.WaitForEffectsCert)
		t.RequestType = &defaultVal
	}
	if t.ConfirmPollSecs == nil {
		defaultVal := uint64(DefaultConfirmPollSecs)
		t.ConfirmPollSecs = &defaultVal
	}
}

type BalanceMonitorConfig struct {
	BalancePollPeriod *string
}

func (b *BalanceMonitorConfig) setDefaults() {
	if b.BalancePollPeriod == nil {
		defaultVal := fmt.Sprintf("%ds", DefaultBalancePollIntervalSeconds)
		b.BalancePollPeriod = &defaultVal
	}
}

// TOMLConfig represents the configuration for a Sui chain, typically loaded from a TOML file.
// It contains all the necessary parameters to configure a Sui relayer including chain ID,
// network details, transaction management settings, and node configurations.
//
// Example TOML configuration:
//
//	[[Sui]]
//	Enabled = true
//	ChainID = "4"
//	NetworkName = "localnet"
//	NetworkNameFull = "sui-localnet"
//
//	[[Sui.Nodes]]
//	Name = 'primary'
//	URL = "https://fullnode.devnet.sui.io:443"
//
//	[Sui.TransactionManager]
//	BroadcastChanSize = 100
//	ConfirmPollSecs = 2
//	DefaultMaxGasAmount = 200000
//	MaxTxRetryAttempts = 5
//	RequestType = 'WaitForEffectsCert'
//	TransactionTimeout = '10s'
//	MaxConcurrentRequests = 5
//
// [Sui.BalanceMonitor]
// BalancePollPeriod = '10s'

type TOMLConfig struct {
	// ChainID is a unique identifier for the Sui chain
	ChainID *string

	// Enabled determines if this chain configuration is active
	// If nil, defaults to true. Use IsEnabled() method to check.
	Enabled *bool

	// NetworkName is the name of the Sui network
	NetworkName *string

	// NetworkNameFull is the full name of the Sui network
	NetworkNameFull *string

	// ChainConfig holds chain-specific configuration parameters
	TransactionManager *TransactionManagerConfig

	// Balance monitor config
	BalanceMonitor *BalanceMonitorConfig

	// Transactions indexer configs (without any transmitter specs, transmitters are attached later)
	TransactionsIndexer *IndexerConfig

	// Events indexer configs (without any event selectors, those are attached later)
	EventsIndexer *IndexerConfig

	// Nodes is a collection of node configurations for this chain
	Nodes NodeConfigs
}

func (c *TOMLConfig) IsEnabled() bool {
	return c.Enabled == nil || *c.Enabled
}

func (c *TOMLConfig) SetFrom(f *TOMLConfig) {
	if f.ChainID != nil {
		c.ChainID = f.ChainID
	}
	if f.Enabled != nil {
		c.Enabled = f.Enabled
	}
	if f.TransactionManager != nil {
		if c.TransactionManager == nil {
			c.TransactionManager = &TransactionManagerConfig{}
		}
		setFromTransactionManager(c.TransactionManager, f.TransactionManager)
	}
	if f.BalanceMonitor != nil {
		if c.BalanceMonitor == nil {
			c.BalanceMonitor = &BalanceMonitorConfig{}
			c.BalanceMonitor.setDefaults()
		}
		setFromBalanceMonitor(c.BalanceMonitor, f.BalanceMonitor)
	}
	c.Nodes.SetFrom(&f.Nodes)
}

func setFromTransactionManager(c, f *TransactionManagerConfig) {
	if f.BroadcastChanSize != nil {
		c.BroadcastChanSize = f.BroadcastChanSize
	}
	if f.ConfirmPollSecs != nil {
		c.ConfirmPollSecs = f.ConfirmPollSecs
	}
	if f.DefaultMaxGasAmount != nil {
		c.DefaultMaxGasAmount = f.DefaultMaxGasAmount
	}
	if f.MaxTxRetryAttempts != nil {
		c.MaxTxRetryAttempts = f.MaxTxRetryAttempts
	}
	if f.RequestType != nil {
		c.RequestType = f.RequestType
	}
	if f.TransactionTimeout != nil {
		c.TransactionTimeout = f.TransactionTimeout
	}
	if f.MaxConcurrentRequests != nil {
		c.MaxConcurrentRequests = f.MaxConcurrentRequests
	}
}

func setFromBalanceMonitor(c, f *BalanceMonitorConfig) {
	if f.BalancePollPeriod != nil {
		c.BalancePollPeriod = f.BalancePollPeriod
	}
}

func (c *TOMLConfig) ValidateConfig() error {
	var err error
	if c.ChainID == nil {
		err = errors.Join(err, config.ErrMissing{Name: "ChainID", Msg: "required for all chains"})
	} else if *c.ChainID == "" {
		err = errors.Join(err, config.ErrEmpty{Name: "ChainID", Msg: "required for all chains"})
	}

	if len(c.Nodes) == 0 {
		err = errors.Join(err, config.ErrMissing{Name: "Nodes", Msg: "must have at least one node"})
	} else {
		for _, node := range c.Nodes {
			err = errors.Join(err, node.ValidateConfig())
		}
	}

	return err
}

func (c *TOMLConfig) TOMLString() (string, error) {
	b, err := toml.Marshal(c)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (c *TOMLConfig) ListNodes() NodeConfigs {
	return c.Nodes
}
