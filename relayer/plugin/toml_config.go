package plugin

import (
	"errors"
	"fmt"
	"math/rand"
	"net/url"

	"github.com/pelletier/go-toml/v2"
	"golang.org/x/exp/slices"

	"github.com/smartcontractkit/chainlink-common/pkg/config"
)

type TOMLConfigs []*TOMLConfig

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
	if f.SolidityURL != nil {
		n.SolidityURL = f.SolidityURL
	}
}

// TOMLConfig represents the configuration for a Sui chain and its nodes.
// It holds chain-specific settings and a collection of nodes for that chain.
// The configuration is typically loaded from a TOML file and used to establish
// connections to the Sui blockchain.
//
// Example Configuration:
//
//	[[Chains]]
//	ChainID = '0x1' # The Sui network ID
//	Enabled = true
//
//	# Transaction settings
//	BroadcastChanSize = 4096     # Size of the broadcast channel buffer
//	ConfirmPollPeriod = '500ms'  # Time between transaction confirmation checks
//	MaxConcurrentRequests = 5    # Maximum number of concurrent RPC requests
//	TransactionTimeout = '10s'   # Timeout for transaction requests
//	NumberRetries = 5            # Number of retries for failed requests
//	GasLimit = 10000000          # Maximum gas limit for transactions
//	RequestType = 'WaitForLocalExecution' # Transaction execution mode (WaitForLocalExecution or WaitForEffectsCert)
//
//	# Node configurations
//	[[Chains.Nodes]]
//	Name = 'sui-node-1'
//	URL = 'https://sui-rpc.example.com'
//	SolidityURL = 'https://sui-rpc.example.com'
//
//	[[Chains.Nodes]]
//	Name = 'sui-node-2'
//	URL = 'https://sui-rpc-backup.example.com'
//	SolidityURL = 'https://sui-rpc-backup.example.com'
type TOMLConfig struct {
	// ChainID is a unique identifier for the Sui chain
	ChainID *string

	// Enabled determines if this chain configuration is active
	// If nil, defaults to true. Use IsEnabled() method to check.
	Enabled *bool

	// ChainConfig holds chain-specific configuration parameters
	ChainConfig

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
	setFromChain(&c.ChainConfig, &f.ChainConfig)
	c.Nodes.SetFrom(&f.Nodes)
}

func setFromChain(c, f *ChainConfig) {
	if f.BroadcastChanSize != nil {
		c.BroadcastChanSize = f.BroadcastChanSize
	}
	if f.ConfirmPollPeriod != nil {
		c.ConfirmPollPeriod = f.ConfirmPollPeriod
	}
	if f.MaxConcurrentRequests != nil {
		c.MaxConcurrentRequests = f.MaxConcurrentRequests
	}
	if f.TransactionTimeout != nil {
		c.TransactionTimeout = f.TransactionTimeout
	}
	if f.NumberRetries != nil {
		c.NumberRetries = f.NumberRetries
	}
	if f.GasLimit != nil {
		c.GasLimit = f.GasLimit
	}
	if f.RequestType != nil {
		c.RequestType = f.RequestType
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
