package plugin

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"

	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainwriter"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/txm"
)

type SuiRelayer struct {
	services.StateMachine

	chainId    string
	chainIdNum *big.Int

	cfg  *TOMLConfig
	lggr logger.Logger

	client *client.PTBClient
	txm    *txm.SuiTxm
	// TODO: balanceMonitor services.Service
}

var _ types.Relayer = &SuiRelayer{}

func NewRelayer(cfg *TOMLConfig, lggr logger.Logger, keystore core.Keystore) (*SuiRelayer, error) {
	id := *cfg.ChainID

	var idNum *big.Int
	var ok bool
	if strings.HasPrefix(id, "0x") {
		idNum, ok = new(big.Int).SetString(id[2:], 16) //nolint:mnd
	} else {
		idNum, ok = new(big.Int).SetString(id, 10) //nolint:mnd
	}

	if !ok {
		return nil, fmt.Errorf("couldn't parse chain id %s", id)
	}

	// Apply default values if not set
	cfg.ChainConfig.Defaults()

	nodeConfig, err := cfg.ListNodes().SelectRandom()
	if err != nil {
		return nil, fmt.Errorf("failed to get node config: %w", err)
	}
	store := txm.NewTxmStoreImpl()
	conf := txm.DefaultConfigSet

	// Use config values instead of constants
	suiClient, err := client.NewPTBClient(
		lggr,
		nodeConfig.URL.String(),
		nil,
		*cfg.TransactionTimeout,
		keystore,
		*cfg.MaxConcurrentRequests,
		*cfg.RequestType,
	)
	if err != nil {
		return nil, fmt.Errorf("error in NewConfigProvider chain.Reader: %w", err)
	}

	retryManager := txm.NewDefaultRetryManager(int(*cfg.NumberRetries))
	gasLimit := big.NewInt(*cfg.GasLimit)
	gasManager := txm.NewSuiGasManager(lggr, suiClient, *gasLimit, 0)

	txManager, err := txm.NewSuiTxm(lggr, suiClient, keystore, conf, store, retryManager, gasManager)

	if err != nil {
		return nil, fmt.Errorf("error in NewConfigProvider chain.Reader: %w", err)
	}

	return &SuiRelayer{
		chainId:    id,
		chainIdNum: idNum,
		cfg:        cfg,
		lggr:       logger.Named(logger.With(lggr, "chainID", id, "chain", "sui"), "SuiRelayer"),
		client:     suiClient,
		txm:        txManager,
	}, nil
}

func (r *SuiRelayer) Name() string {
	return "SuiRelayer"
}

func (r *SuiRelayer) Start(ctx context.Context) error {
	return r.StartOnce("SuiRelayer", func() error {
		r.lggr.Debug("Starting")
		r.lggr.Debug("Starting txm")
		var ms services.MultiStart

		return ms.Start(ctx, r.txm)
	})
}

func (r *SuiRelayer) Close() error {
	return r.StopOnce("SuiRelayer", func() error {
		r.lggr.Debug("Stopping")
		r.lggr.Debug("Stopping txm")

		return r.txm.Close()
	})
}

func (r *SuiRelayer) Ready() error {
	return errors.Join(
		r.StateMachine.Ready(),
		r.txm.Ready(),
	)
}

func (r *SuiRelayer) HealthReport() map[string]error {
	report := map[string]error{r.Name(): r.Healthy()}
	services.CopyHealth(report, r.txm.HealthReport())

	return report
}

// ChainService interface
func (r *SuiRelayer) GetChainStatus(ctx context.Context) (types.ChainStatus, error) {
	toml, err := r.cfg.TOMLString()
	if err != nil {
		return types.ChainStatus{}, err
	}

	return types.ChainStatus{
		ID:      r.chainId,
		Enabled: r.cfg.IsEnabled(),
		Config:  toml,
	}, nil
}

func (r *SuiRelayer) ListNodeStatuses(ctx context.Context, pageSize int32, pageToken string) ([]types.NodeStatus, string, int, error) {
	return []types.NodeStatus{}, "", 0, errors.New("not implemented")
}

func (r *SuiRelayer) Transact(ctx context.Context, from, to string, amount *big.Int, balanceCheck bool) error {
	return errors.New("not implemented")
}

// Relayer interface
func (r *SuiRelayer) NewContractWriter(_ context.Context, configBytes []byte) (types.ContractWriter, error) {
	chainConfig := chainwriter.ChainWriterConfig{}
	err := json.Unmarshal(configBytes, &chainConfig)
	if err != nil {
		return nil, fmt.Errorf("error in NewContractWriter: %w", err)
	}

	// TODO: validate chainConfig

	chainWriter, err := chainwriter.NewSuiChainWriter(r.lggr, r.txm, chainConfig, false)
	if err != nil {
		return nil, fmt.Errorf("error in NewContractWriter: %w", err)
	}

	return chainWriter, nil
}

func (r *SuiRelayer) NewContractReader(ctx context.Context, contractReaderConfig []byte) (types.ContractReader, error) {
	chainConfig := chainreader.ChainReaderConfig{}
	err := json.Unmarshal(contractReaderConfig, &chainConfig)
	if err != nil {
		return nil, fmt.Errorf("error in NewContractReader: %w", err)
	}

	// TODO: validate chainConfig

	chainReader := chainreader.NewChainReader(r.lggr, *r.client, chainConfig)

	return chainReader, nil
}

func (r *SuiRelayer) NewConfigProvider(ctx context.Context, args types.RelayArgs) (types.ConfigProvider, error) {
	return nil, errors.New("not implemented")
}

func (r *SuiRelayer) NewPluginProvider(ctx context.Context, relayargs types.RelayArgs, pluginargs types.PluginArgs) (types.PluginProvider, error) {
	// TODO: is this necessary? should we just return an error?
	return r.NewMedianProvider(ctx, relayargs, pluginargs)
}

func (r *SuiRelayer) NewLLOProvider(ctx context.Context, relayargs types.RelayArgs, pluginargs types.PluginArgs) (types.LLOProvider, error) {
	return nil, errors.New("LLO not supported for Sui")
}

// implement MedianProvider type from github.com/smartcontractkit/chainlink-common/pkg/loop/internal/types
//
// if the loop.Relayer returned by NewRelayer supports the internal loop type MedianProvider, it's called here:
// see https://github.com/smartcontractkit/chainlink-common/blob/7c11e2c2ce3677f57239c40585b04fd1c9ce1713/pkg/loop/internal/relayer/relayer.go#L493
func (r *SuiRelayer) NewMedianProvider(ctx context.Context, relayargs types.RelayArgs, pluginargs types.PluginArgs) (types.MedianProvider, error) {
	return nil, errors.New("ocr2 is not supported for sui")
}

func (r *SuiRelayer) LatestHead(ctx context.Context) (types.Head, error) {
	return types.Head{}, errors.New("not implemented")
}

// NewAutomationProvider returns a new automation provider for the given relay and plugin arguments.
// Currently not supported for Sui.
func (r *SuiRelayer) NewAutomationProvider(ctx context.Context, rargs types.RelayArgs, pargs types.PluginArgs) (types.AutomationProvider, error) {
	return nil, errors.New("automation not supported for Sui")
}

// Replay implements the transaction replay functionality.
// Currently not supported for Sui.
func (r *SuiRelayer) Replay(ctx context.Context, chainID string, data map[string]any) error {
	return errors.New("replay not supported for Sui")
}

// NewCCIPCommitProvider returns a new CCIP commit provider for the given relay and plugin arguments.
// Currently not supported for Sui.
func (r *SuiRelayer) NewCCIPCommitProvider(ctx context.Context, rargs types.RelayArgs, pargs types.PluginArgs) (types.CCIPCommitProvider, error) {
	return nil, errors.New("cCIP not supported for Sui")
}

// NewCCIPExecProvider returns a new CCIP exec provider for the given relay and plugin arguments.
// Currently not supported for Sui.
func (r *SuiRelayer) NewCCIPExecProvider(ctx context.Context, rargs types.RelayArgs, pargs types.PluginArgs) (types.CCIPExecProvider, error) {
	return nil, errors.New("cCIP not supported for Sui")
}

// NewFunctionsProvider returns a new Functions provider for the given relay and plugin arguments.
// Currently not supported for Sui.
func (r *SuiRelayer) NewFunctionsProvider(ctx context.Context, rargs types.RelayArgs, pargs types.PluginArgs) (types.FunctionsProvider, error) {
	return nil, errors.New("functions not supported for Sui")
}

// NewMercuryProvider returns a new Mercury provider for the given relay arguments.
// Currently not supported for Sui.
func (r *SuiRelayer) NewMercuryProvider(ctx context.Context, rargs types.RelayArgs, pargs types.PluginArgs) (types.MercuryProvider, error) {
	return nil, errors.New("mercury not supported for Sui")
}

// NewOCR3CapabilityProvider returns a new OCR3 capability provider for the given relay and plugin arguments.
// Currently not supported for Sui.
func (r *SuiRelayer) NewOCR3CapabilityProvider(ctx context.Context, rargs types.RelayArgs, pargs types.PluginArgs) (types.OCR3CapabilityProvider, error) {
	return nil, errors.New("ocr3 not supported for Sui")
}
