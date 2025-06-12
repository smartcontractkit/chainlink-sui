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
const CCIPExecuteReportFunctionName = "CCIPExecuteReport"

type SuiChainWriter struct {
	lggr       logger.Logger
	txm        txm.TxManager
	config     ChainWriterConfig
	simulate   bool
	ptbFactory *PTBConstructor
	services.StateMachine
}

func NewSuiChainWriter(lggr logger.Logger, txManager txm.TxManager, config ChainWriterConfig, simulate bool) (*SuiChainWriter, error) {
	suiClient := txManager.GetClient()
	return &SuiChainWriter{
		lggr:       logger.Named(lggr, ServiceName),
		txm:        txManager,
		config:     config,
		simulate:   simulate,
		ptbFactory: NewPTBConstructor(config, suiClient, lggr),
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

// SubmitTransaction is the primary entry point for submitting transactions via the SuiChainWriter.
// It acts as a router, determining whether to enqueue a standard smart contract call or a
// Programmable Transaction Block (PTB) based on the provided contractName.
// If contractName matches PTBChainWriterModuleName, it assumes a PTB submission and calls enqueuePTB.
// Otherwise, it treats the request as a standard Move function call and calls enqueueSmartContractCall.
// This function implements the commonTypes.ContractWriter interface.
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
	if contractName == PTBChainWriterModuleName {
		return enqueuePTB(ctx, s, contractName, method, args, transactionID, toAddress, meta)
	}

	return enqueueSmartContractCall(ctx, s, contractName, method, args, transactionID, meta)
}

// enqueueSmartContractCall handles the process of enqueuing a standard smart contract (Move function) call.
// It retrieves module and function configurations, converts arguments, constructs the full function signature,
// and then calls the TxManager's Enqueue method to generate, sign, store, and queue the transaction.
//
// Parameters:
//   - ctx: Context for the operation.
//   - s: The SuiChainWriter instance containing configuration and TxManager.
//   - contractName: The name of the contract module as defined in the configuration.
//   - method: The name of the function to call within the contract module.
//   - args: The arguments for the function call, provided as a map or struct.
//   - transactionID: The unique identifier for the transaction.
//   - meta: Transaction metadata (e.g., gas limits).
//
// Returns:
//   - error: An error if configu\ration is not found, argument conversion fails, or the TxManager Enqueue call fails.
func enqueueSmartContractCall(ctx context.Context, s *SuiChainWriter, contractName string, method string, args any, transactionID string, meta *commonTypes.TxMeta) error {
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
	// TODO: Add support for generic type args
	typeArgs := []string{}

	var arguments Arguments
	if err := mapstructure.Decode(args, &arguments); err != nil {
		return fmt.Errorf("failed to decode args: %w", err)
	}

	paramTypes, paramValues, err := convertFunctionParams(arguments.Args, functionConfig.Params)
	if err != nil {
		s.lggr.Errorw("Error converting function params", "error", err)
		return err
	}

	suiFunction := fmt.Sprintf("%s::%s::%s", moduleConfig.ModuleID, contractName, method)

	tx, err := s.txm.Enqueue(ctx, transactionID, meta, functionConfig.PublicKey, suiFunction, typeArgs, paramTypes, paramValues, s.simulate)
	if err != nil {
		s.lggr.Errorw("Error enqueuing transaction", "error", err)
		return err
	}
	s.lggr.Infow("Transaction enqueued", "transactionID", tx.TransactionID, "functionName", method)

	return nil
}

// enqueuePTB handles the process of enqueuing a Programmable Transaction Block (PTB).
// It retrieves the PTB configuration, automatically builds a PTBArgBuilder from the raw input arguments,
// constructs the PTB commands, and then calls the TxManager's EnqueuePTB method to generate, sign, store,
// and queue the PTB transaction.
//
// The function leverages the BuildFromConfig method to automatically map the input arguments to the
// appropriate command parameters based on the configuration, using the builder pattern internally.
// This approach simplifies the client API while maintaining the flexibility and type safety of the
// PTBArgBuilder.
//
// Parameters:
//   - ctx: Context for the operation.
//   - s: The SuiChainWriter instance containing configuration, PTBConstructor, and TxManager.
//   - ptbName: The name of the PTB configuration as defined in the Modules map.
//   - method: The virtual function name within the PTB configuration that defines the command sequence.
//   - args: The arguments needed to build the PTB commands, provided as a simple map or struct. These are
//     automatically mapped to the appropriate command parameters based on the PTB configuration.
//   - transactionID: The unique identifier for the transaction.
//   - meta: Transaction metadata (e.g., gas limits).
//
// Returns:
//   - error: An error if configuration is not found, argument mapping fails, PTB command building fails,
//     or the TxManager EnqueuePTB call fails.
func enqueuePTB(ctx context.Context, s *SuiChainWriter, ptbName string, method string, args any, transactionID string, toAddress string, meta *commonTypes.TxMeta) error {
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

	var arguments Arguments
	if err := mapstructure.Decode(args, &arguments); err != nil {
		return fmt.Errorf("failed to decode args: %w", err)
	}

	// TODO: Placeholder, this will be implemented in another PR
	// Use the builder with the updated PTBConstructor
	// optional pass the config overrides
	// if method == CCIPExecuteReportFunctionName {
	// 	var execArgs SuiOffRampExecCallArgs
	// 	if err := mapstructure.Decode(args, &execArgs); err != nil {
	// 		return fmt.Errorf("failed to decode args: %w", err)
	// 	}

	// 	ptbCommands, updatedArgs, err := s.ptbExpander.GetOffRampPTB(s.lggr, execArgs, functionConfig, functionConfig.PublicKey)
	// 	if err != nil {
	// 		s.lggr.Errorw("Error expanding PTB commands", "error", err)
	// 		return err
	// 	}
	// 	args = updatedArgs
	// 	fmt.Println("updatedArgs", updatedArgs)
	// 	fmt.Println("ptbCommands", ptbCommands)
	// } else {
	// 	if err := mapstructure.Decode(args, &arguments); err != nil {
	// 		return fmt.Errorf("failed to decode args: %w", err)
	// 	}
	// }

	ptbCommands, err := s.ptbFactory.BuildPTBCommands(ctx, ptbName, method, arguments, &ConfigOverrides{
		ToAddress: toAddress,
	})

	if err != nil {
		s.lggr.Errorw("Error building PTB commands", "error", err)
		return err
	}

	ptb := ptbCommands.Finish()
	tx, err := s.txm.EnqueuePTB(ctx, transactionID, meta, functionConfig.PublicKey, &ptb, s.simulate)
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
