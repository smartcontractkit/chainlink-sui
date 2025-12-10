package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/smartcontractkit/mcms"
	suisdk "github.com/smartcontractkit/mcms/sdk/sui"
	"github.com/smartcontractkit/mcms/types"

	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	"github.com/smartcontractkit/chainlink-sui/relayer/signer"
)

var DefaultTimelockExpirationInHours = 72

// TimelockConfig is based on chainlink/deployment proposal utils
type TimelockConfig struct {
	MCMSAction   types.TimelockAction `json:"mcmsAction"`
	MinDelay     time.Duration        `json:"minDelay"`     // delay for timelock worker to execute the transfers.
	OverrideRoot bool                 `json:"overrideRoot"` // if true, override the previous root with the new one.
}

type GenerateProposalInput struct {
	ChainSelector      uint64
	Client             sui.ISuiAPI
	MCMSPackageID      string
	MCMSStateObjID     string
	AccountObjID       string
	RegistryObjID      string
	TimelockObjID      string
	DeployerStateObjID string
	Description        string
	BatchOp            types.BatchOperation
	TimelockConfig     TimelockConfig
}

func GenerateProposal(ctx context.Context, input GenerateProposalInput) (*mcms.TimelockProposal, error) {
	// Get action and delay from role
	var delay *types.Duration
	if input.TimelockConfig.MCMSAction == types.TimelockActionSchedule {
		delayDuration := types.NewDuration(input.TimelockConfig.MinDelay)
		delay = &delayDuration
	}
	role, err := getRoleFromAction(input.TimelockConfig.MCMSAction)
	if err != nil {
		return nil, fmt.Errorf("failed to get action from role: %w", err)
	}

	// Get OP Count from inspector
	devInspectSigner := signer.NewDevInspectSigner("0x0")
	inspector, err := suisdk.NewInspector(input.Client, devInspectSigner, input.MCMSPackageID, role)
	if err != nil {
		return nil, fmt.Errorf("failed to create inspector: %w", err)
	}
	opCount, err := inspector.GetOpCount(ctx, input.MCMSStateObjID)
	if err != nil {
		return nil, fmt.Errorf("failed to get op count: %w", err)
	}

	// Build metadata
	metadata, err := suisdk.NewChainMetadata(opCount, role, input.MCMSPackageID, input.MCMSStateObjID, input.AccountObjID, input.RegistryObjID, input.TimelockObjID, input.DeployerStateObjID)
	if err != nil {
		return nil, fmt.Errorf("failed to create chain metadata: %w", err)
	}

	// Build proposal
	validUntilMs := uint32(time.Now().Add(time.Duration(DefaultTimelockExpirationInHours) * time.Hour).Unix())
	builder := mcms.NewTimelockProposalBuilder().
		SetVersion("v1").
		SetValidUntil(validUntilMs).
		SetDescription(input.Description).
		AddTimelockAddress(types.ChainSelector(input.ChainSelector), input.TimelockObjID).
		AddChainMetadata(types.ChainSelector(input.ChainSelector), metadata).
		AddOperation(input.BatchOp).
		SetAction(input.TimelockConfig.MCMSAction)

	if delay != nil {
		builder.SetDelay(*delay)
	}

	return builder.Build()
}

func ExtractTransactionCall(output interface{}, operationID string) (sui_ops.TransactionCall, error) {
	jsonBytes, err := json.Marshal(output)
	if err != nil {
		return sui_ops.TransactionCall{}, fmt.Errorf("failed to marshal operation %s output: %w", operationID, err)
	}

	var outputMap map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &outputMap); err != nil {
		return sui_ops.TransactionCall{}, fmt.Errorf("failed to unmarshal operation %s output: %w", operationID, err)
	}

	callInterface, exists := outputMap["Call"]
	if !exists {
		return sui_ops.TransactionCall{}, fmt.Errorf("operation %s output does not have a Call field", operationID)
	}

	callBytes, err := json.Marshal(callInterface)
	if err != nil {
		return sui_ops.TransactionCall{}, fmt.Errorf("failed to marshal Call field for operation %s: %w", operationID, err)
	}

	var call sui_ops.TransactionCall
	if err := json.Unmarshal(callBytes, &call); err != nil {
		return sui_ops.TransactionCall{}, fmt.Errorf("failed to unmarshal Call field for operation %s: %w", operationID, err)
	}

	return call, nil
}

func getRoleFromAction(action types.TimelockAction) (suisdk.TimelockRole, error) {
	switch action {
	case types.TimelockActionSchedule:
		return suisdk.TimelockRoleProposer, nil
	case types.TimelockActionBypass:
		return suisdk.TimelockRoleBypasser, nil
	case types.TimelockActionCancel:
		return suisdk.TimelockRoleCanceller, nil
	default:
		// NewChainMetadata will always error on invalid action, but this is a safeguard
		return 0, fmt.Errorf("unsupported action: %v", action)
	}
}
