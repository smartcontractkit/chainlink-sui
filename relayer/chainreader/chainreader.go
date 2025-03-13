package chainreader

import (
	"context"
	"errors"

	// sui "github.com/block-vision/sui-go-sdk"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
	"github.com/smartcontractkit/chainlink-common/pkg/utils"
	// TODO: enable after codec and txm is implemented
	// "github.com/smartcontractkit/chainlink-internal-integrations/sui/relayer/codec"
	// "github.com/smartcontractkit/chainlink-internal-integrations/sui/relayer/txm"
)

type suiChainReader struct {
	types.UnimplementedContractReader

	logger  logger.Logger
	config  ChainReaderConfig
	starter utils.StartStopOnce
	// moduleAddresses map[string]sui.AccountAddress
	// client *sui.NodeClient
}

func NewChainReader(lgr logger.Logger /*client *sui.SuiClient,*/, config ChainReaderConfig) types.ContractReader {
	return &suiChainReader{
		logger: logger.Named(lgr, "SuiChainReader"),
		// client: client,
		config: config,
		// moduleAddresses: map[string]sui.AccountAddress{},
	}
}

func (a *suiChainReader) Name() string {
	return a.logger.Name()
}

func (a *suiChainReader) Ready() error {
	return a.starter.Ready()
}

func (a *suiChainReader) HealthReport() map[string]error {
	return map[string]error{a.Name(): a.starter.Healthy()}
}

func (a *suiChainReader) Start(ctx context.Context) error {
	return a.starter.StartOnce(a.Name(), func() error {
		return nil
	})
}

func (a *suiChainReader) Close() error {
	return a.starter.StopOnce(a.Name(), func() error {
		return nil
	})
}

func (a *suiChainReader) GetLatestValue(ctx context.Context, readIdentifier string, confidenceLevel primitives.ConfidenceLevel, params, returnVal any) error {
	return errors.New("not implemented")
}

func (a *suiChainReader) BatchGetLatestValues(ctx context.Context, request types.BatchGetLatestValuesRequest) (types.BatchGetLatestValuesResult, error) {
	return nil, errors.New("not implemented")
}

func (a *suiChainReader) Bind(ctx context.Context, bindings []types.BoundContract) error {
	return errors.New("not implemented")
}

func (a *suiChainReader) Unbind(ctx context.Context, bindings []types.BoundContract) error {
	return errors.New("not implemented")
}

func (a *suiChainReader) QueryKey(ctx context.Context, contract types.BoundContract, filter query.KeyFilter, limitAndSort query.LimitAndSort, sequenceDataType any) ([]types.Sequence, error) {
	return nil, errors.New("not implemented")
}
