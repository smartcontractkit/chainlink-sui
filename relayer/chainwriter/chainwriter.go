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

	cwConfig "github.com/smartcontractkit/chainlink-sui/relayer/chainwriter/config"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainwriter/ptb"
	"github.com/smartcontractkit/chainlink-sui/relayer/txm"
)

const ServiceName = "SuiChainWriter"

type SuiChainWriter struct {
	lggr       logger.Logger
	txm        txm.TxManager
	config     cwConfig.ChainWriterConfig
	simulate   bool
	ptbFactory *ptb.PTBConstructor
	services.StateMachine
}

func NewSuiChainWriter(lggr logger.Logger, txManager txm.TxManager, config cwConfig.ChainWriterConfig, simulate bool) (*SuiChainWriter, error) {
	suiClient := txManager.GetClient()
	return &SuiChainWriter{
		lggr:       logger.Named(lggr, ServiceName),
		txm:        txManager,
		config:     config,
		simulate:   simulate,
		ptbFactory: ptb.NewPTBConstructor(config, suiClient, lggr),
	}, nil
}

// SubmitTransaction is the primary entry point for submitting transactions via the SuiChainWriter.
// It acts as a router, determining whether to enqueue a standard smart contract call or a
// Programmable Transaction Block (PTB) based on the provided contractName.
//
// Parameters:
//   - ctx: The context for the operation, allowing for cancellation and timeouts.
//   - contractName: The identifier for the target module or a special identifier (PTBChainWriterModuleName)
//     indicating a PTB submission defined in the configuration.
//   - method: The specific function name within the module (for standard calls) or the virtual function
//     name defined in the PTB configuration.
//   - args: The arguments required by the function or PTB commands. For PTB submissions, these are automatically
//     mapped to commands based on the configuration using the builder pattern internally.
//   - transactionID: A unique identifier for this transaction attempt.
//   - toAddress: The target address for the transaction (Note: Often implicitly handled by the module/function config in Sui).
//   - meta: Transaction metadata, primarily used for specifying gas limits (*commontypes.TxMeta).
//   - _ *big.Int: An unused parameter, present for interface compatibility.
//
// Returns:
//   - error: An error if the configuration is missing, argument processing fails, or the underlying
//     transaction enqueue operation in the TxManager fails.
func (s *SuiChainWriter) SubmitTransaction(ctx context.Context, contractName string, method string, args any, transactionID string, toAddress string, meta *commonTypes.TxMeta, _ *big.Int) error {
	s.lggr.Info("RECIEVED ARGS: ", args)

	ptbName := contractName

	moduleConfig, exists := s.config.Modules[ptbName]
	if !exists {
		s.lggr.Errorw("PBT not found", "PTB name", ptbName)
		return commonTypes.ErrNotFound
	}

	functionConfig, exists := moduleConfig.Functions[method]
	if !exists {
		s.lggr.Errorw("Function not found", "functionName", method)
		return commonTypes.ErrNotFound
	}

	var arguments cwConfig.Arguments
	if err := mapstructure.Decode(args, &arguments.Args); err != nil {
		return fmt.Errorf("failed to decode args: %w", err)
	}

	s.lggr.Info("ARGUMENTS: ", arguments.Args, "PARAMS: ", functionConfig.Params, "DECODED ARGS: ", arguments)

	// overwrite ptbName
	if moduleConfig.Name != "" {
		ptbName = moduleConfig.Name
	}

	if functionConfig.Name != "" {
		method = functionConfig.Name
	}

	ptbService, err := s.ptbFactory.BuildPTBCommands(ctx, ptbName, method, arguments, toAddress, functionConfig)

	if err != nil {
		s.lggr.Errorw("Error building PTB commands", "error", err)
		return err
	}

	s.lggr.Infow("PTB commands", "ptb", ptbService, "functionConfig", functionConfig)

	tx, err := s.txm.EnqueuePTB(ctx, transactionID, meta, functionConfig.PublicKey, ptbService, s.simulate)
	if err != nil {
		s.lggr.Errorw("Error enqueuing PTB", "error", err)
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

func (s *SuiChainWriter) GetEstimateFee(ctx context.Context, contractName string, method string, args any, transactionID string, meta *commonTypes.TxMeta, _ *big.Int) (commonTypes.EstimateFee, error) {
	return commonTypes.EstimateFee{}, errors.New("GetEstimateFee not implemented")
}

// Close implements types.ContractWriter.
func (s *SuiChainWriter) Close() error {
	return s.StopOnce(ServiceName, func() error {
		s.lggr.Infow("Stopping SuiChainWriter")
		return nil
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
		return nil
	})
}

var (
	_ commonTypes.ContractWriter = &SuiChainWriter{}
	_ services.Service           = &SuiChainWriter{}
)
