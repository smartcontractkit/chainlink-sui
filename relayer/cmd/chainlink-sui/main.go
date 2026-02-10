package main

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	config2 "github.com/smartcontractkit/chainlink-sui/relayer/config"

	"github.com/hashicorp/go-plugin"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"

	suiplugin "github.com/smartcontractkit/chainlink-sui/relayer/plugin"
)

const (
	loggerName = "PluginSui"
)

func main() {
	s := loop.MustNewStartedServer(loggerName)
	defer s.Stop()

	p := &pluginRelayer{Plugin: loop.Plugin{Logger: s.Logger}, db: s.DataSource, lgr: s.Logger}
	defer s.Logger.ErrorIfFn(p.Close, "Failed to close")

	s.MustRegister(p)

	stopCh := make(chan struct{})
	defer close(stopCh)

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: loop.PluginRelayerHandshakeConfig(),
		Plugins: map[string]plugin.Plugin{
			loop.PluginRelayerName: &loop.GRPCPluginRelayer{
				PluginServer: p,
				BrokerConfig: loop.BrokerConfig{
					StopCh:   stopCh,
					Logger:   s.Logger,
					GRPCOpts: s.GRPCOpts,
				},
			},
		},
		GRPCServer: s.GRPCOpts.NewServer,
	})
}

type pluginRelayer struct {
	loop.Plugin
	db  sqlutil.DataSource
	lgr logger.Logger
}

var _ loop.PluginRelayer = &pluginRelayer{}

func (c *pluginRelayer) NewRelayer(ctx context.Context, rawConfig string, keystore, csa loop.Keystore, capRegistry core.CapabilitiesRegistry) (loop.Relayer, error) {
	cfg, err := config2.NewDecodedTOMLConfig(rawConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to read configs: %w", err)
	}

	rawNodes := make([]map[string]string, 0, len(cfg.Nodes))
	for _, n := range cfg.Nodes {
		if n == nil || n.URL == nil {
			continue
		}
		rawNodes = append(rawNodes, map[string]string{"URL": n.URL.String()})
	}
	chainID := ""
	if cfg.ChainID != nil {
		chainID = *cfg.ChainID
	}
	emitter := loop.NewPluginRelayerConfigEmitter(
		c.Logger,
		beholder.GetClient().Config.AuthPublicKeyHex,
		chainID,
		rawNodes,
	)
	if err := emitter.Start(ctx); err != nil {
		return nil, fmt.Errorf("failed to start plugin relayer config emitter: %w", err)
	}
	c.SubService(emitter)

	relayer, err := suiplugin.NewRelayer(cfg, c.Logger, keystore, c.db)
	if err != nil {
		return nil, fmt.Errorf("failed to create relayer: %w", err)
	}

	c.SubService(relayer)

	return relayer, nil
}
