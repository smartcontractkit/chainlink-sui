package chainwriter

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/mitchellh/mapstructure"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	commonTypes "github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
	"github.com/smartcontractkit/chainlink-sui/relayer/txm"
)

const ServiceName = "SuiChainWriter"

type SuiChainWriter struct {
	lggr     logger.Logger
	txm      txm.TxManager
	config   ChainWriterConfig
	simulate bool

	services.StateMachine
}

func NewSuiChainWriter(lggr logger.Logger, txManager txm.TxManager, config ChainWriterConfig, simulate bool) (*SuiChainWriter, error) {
	return &SuiChainWriter{
		lggr:     logger.Named(lggr, ServiceName),
		txm:      txManager,
		config:   config,
		simulate: simulate,
	}, nil
}

func convertFunctionParams(argMap map[string]any, params []codec.SuiFunctionParam) ([]string, []any, error) {
	types := make([]string, len(params))
	values := make([]any, len(params))

	if len(argMap) != len(params) {
		return nil, nil, errors.New("argument count mismatch")
	}

	for i, paramConfig := range params {
		argValue, ok := argMap[paramConfig.Name]
		if !ok {
			if paramConfig.Required {
				return nil, nil, fmt.Errorf("missing argument: %s", paramConfig.Name)
			}
			argValue = paramConfig.DefaultValue
		}

		types[i] = paramConfig.Type
		values[i] = argValue
	}

	return types, values, nil
}

// SubmitTransaction implements types.ContractWriter
func (s *SuiChainWriter) SubmitTransaction(ctx context.Context, contractName string, method string, args any, transactionID string, toAddress string, meta *commonTypes.TxMeta, value *big.Int) error {
	moduleConfig, exists := s.config.Modules[contractName]
	if !exists {
		s.lggr.Errorw("Contract not found", "contractName", contractName)
		return commonTypes.ErrNotFound
	}

	functionConfig, exists := moduleConfig.Functions[method]
	if !exists {
		s.lggr.Errorw("Function not found", "functionName", method)
		return commonTypes.ErrNotFound
	}

	// For now do not assume any generic type args
	typeArgs := []string{}

	argMap := make(map[string]any)
	err := mapstructure.Decode(args, &argMap)
	if err != nil {
		s.lggr.Errorw("Error decoding args", "error", err)
		return err
	}
	paramTypes, paramValues, err := convertFunctionParams(argMap, functionConfig.Params)
	if err != nil {
		s.lggr.Errorw("Error converting function params", "error", err)
		return err
	}

	suiFunction := fmt.Sprintf("%s::%s::%s", moduleConfig.ModuleID, contractName, method)

	tx, err := s.txm.Enqueue(ctx, transactionID, meta, functionConfig.FromAddress, suiFunction, typeArgs, paramTypes, paramValues, s.simulate)
	if err != nil {
		s.lggr.Errorw("Error enqueuing transaction", "error", err)
		return err
	}
	s.lggr.Infow("Transaction enqueued", "transactionID", tx.TransactionID, "functionName", method)

	return nil
}

// GetFeeComponents implements types.ContractWriter.
func (s *SuiChainWriter) GetFeeComponents(ctx context.Context) (*commonTypes.ChainFeeComponents, error) {
	return nil, errors.New("GetFeeComponents not implemented")
}

// GetTransactionStatus implements types.ContractWriter.
func (s *SuiChainWriter) GetTransactionStatus(ctx context.Context, transactionID string) (commonTypes.TransactionStatus, error) {
	return s.txm.GetTransactionStatus(ctx, transactionID)
}

// Close implements types.ContractWriter.
func (s *SuiChainWriter) Close() error {
	return s.StopOnce(ServiceName, func() error {
		s.lggr.Infow("Stopping SuiChainWriter")
		return s.txm.Close()
	})
}

// HealthReport implements types.ContractWriter.
func (s *SuiChainWriter) HealthReport() map[string]error {
	return map[string]error{s.Name(): s.Healthy()}
}

// Name implements types.ContractWriter.
func (s *SuiChainWriter) Name() string {
	return s.lggr.Name()
}

// Ready implements types.ContractWriter.
func (s *SuiChainWriter) Ready() error {
	return s.StateMachine.Ready()
}

// Start implements types.ContractWriter.
func (s *SuiChainWriter) Start(ctx context.Context) error {
	return s.StartOnce(ServiceName, func() error {
		s.lggr.Infow("Starting SuiChainWriter")
		return s.txm.Start(ctx)
	})
}

var (
	_ commonTypes.ContractWriter = &SuiChainWriter{}
	_ services.Service           = &SuiChainWriter{}
)
