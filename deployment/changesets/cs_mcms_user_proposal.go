package changesets

import (
	"fmt"
	"time"

	"github.com/smartcontractkit/mcms"
	suisdk "github.com/smartcontractkit/mcms/sdk/sui"
	"github.com/smartcontractkit/mcms/types"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	mcmsuser "github.com/smartcontractkit/chainlink-sui/bindings/packages/mcms/mcms_user"
	"github.com/smartcontractkit/chainlink-sui/deployment"
)

type InvokeMCMSFunctionOneConfig struct {
	// MCMS related
	MmcsPackageID      string `json:"mcmsPackageID"`
	McmsStateObjID     string `json:"mcmsStateObjID"`
	TimelockObjID      string `json:"timelockObjID"`
	AccountObjID       string `json:"accountObjID"`
	RegistryObjID      string `json:"registryObjID"`
	DeployerStateObjID string `json:"deployerStateObjID"`

	// IsFastCurse selects the fastcurse MCMS instance for the proposal.
	// MCMS fields above are auto-populated from state when empty.
	IsFastCurse bool `json:"isFastCurse,omitempty"`

	// Proposal related
	Role  suisdk.TimelockRole `json:"role"`
	Delay time.Duration       `json:"delay"`

	// MCMS User related
	McmcsUserPackageID  string `json:"mcmsUserPackageID"`
	McmsUserObjectID    string `json:"mcmsUserObjectID"`
	McmsUserOwnerCapObj string `json:"mcmsUserOwnerCapObj"`

	// Chain related
	ChainSelector uint64 `json:"chainSelector"`
}

var _ cldf.ChangeSetV2[InvokeMCMSFunctionOneConfig] = InvokeMCMSFunctionOne{}

type InvokeMCMSFunctionOne struct{}

func (d InvokeMCMSFunctionOne) Apply(e cldf.Environment, config InvokeMCMSFunctionOneConfig) (cldf.ChangesetOutput, error) {
	// Auto-populate MCMS fields from on-chain state when not provided.
	if config.MmcsPackageID == "" || config.McmsStateObjID == "" || config.TimelockObjID == "" || config.AccountObjID == "" || config.RegistryObjID == "" {
		suiState, err := deployment.LoadOnchainStatesui(e)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to load sui onchain state: %w", err)
		}
		mcmsFields := suiState[config.ChainSelector].MCMSState(config.IsFastCurse)
		config.MmcsPackageID = mcmsFields.PackageID
		config.McmsStateObjID = mcmsFields.StateObjectID
		config.TimelockObjID = mcmsFields.TimelockObjectID
		config.AccountObjID = mcmsFields.AccountStateObjectID
		config.RegistryObjID = mcmsFields.RegistryObjectID
		config.DeployerStateObjID = mcmsFields.DeployerStateObjectID
	}
	suiChains := e.BlockChains.SuiChains()

	suiChain := suiChains[config.ChainSelector]

	mcmsUserContract, err := mcmsuser.NewMCMSUser(config.McmcsUserPackageID, suiChain.Client)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to create mcms user contract instance: %w", err)
	}

	arg1 := "Updated Field A"
	arg2 := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	encodedCall, err := mcmsUserContract.MCMSUser().Encoder().FunctionOne(
		bind.Object{Id: config.McmsUserObjectID},
		bind.Object{Id: config.McmsUserOwnerCapObj},
		arg1,
		arg2,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to encode function call: %w", err)
	}

	callBytes := extractByteArgsFromEncodedCall(*encodedCall)
	transaction, err := suisdk.NewTransactionWithStateObj(
		encodedCall.Module.ModuleName,
		encodedCall.Function,
		encodedCall.Module.PackageID,
		callBytes,
		"MCMSUser",
		[]string{},
		config.McmsUserObjectID,
		[]string{},
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to create transaction: %w", err)
	}

	op := types.BatchOperation{
		ChainSelector: types.ChainSelector(config.ChainSelector),
		Transactions:  []types.Transaction{transaction},
	}

	validUntilMs := uint32(time.Now().Add(time.Hour * 24).Unix())
	metadata, err := suisdk.NewChainMetadata(0, config.Role, config.MmcsPackageID, config.McmsStateObjID, config.AccountObjID, config.RegistryObjID, config.TimelockObjID, config.DeployerStateObjID)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to create chain metadata: %w", err)
	}

	var action types.TimelockAction
	var delay *types.Duration
	switch config.Role {
	case suisdk.TimelockRoleProposer:
		action = types.TimelockActionSchedule
		delayDuration := types.NewDuration(config.Delay)
		delay = &delayDuration
	case suisdk.TimelockRoleBypasser:
		action = types.TimelockActionBypass
	default:
		return cldf.ChangesetOutput{}, fmt.Errorf("unsupported role: %v", config.Role)
	}

	builder := mcms.NewTimelockProposalBuilder().
		SetVersion("v1").
		SetValidUntil(validUntilMs).
		SetDescription("Invokes function one from MCMS user contract").
		AddTimelockAddress(types.ChainSelector(config.ChainSelector), config.TimelockObjID).
		AddChainMetadata(types.ChainSelector(config.ChainSelector), metadata).
		AddOperation(op).
		SetAction(action)

	if delay != nil {
		builder.SetDelay(*delay)
	}

	timelockProposal, err := builder.Build()
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
	}

	return cldf.ChangesetOutput{
		MCMSTimelockProposals: []mcms.TimelockProposal{*timelockProposal},
	}, nil
}

// VerifyPreconditions implements deployment.ChangeSetV2.
func (d InvokeMCMSFunctionOne) VerifyPreconditions(e cldf.Environment, config InvokeMCMSFunctionOneConfig) error {
	return nil
}

func extractByteArgsFromEncodedCall(encodedCall bind.EncodedCall) []byte {
	var args []byte
	for _, callArg := range encodedCall.CallArgs {
		if callArg.CallArg.UnresolvedObject != nil {
			args = append(args, callArg.CallArg.UnresolvedObject.ObjectId[:]...)
		}
		if callArg.CallArg.Pure != nil {
			args = append(args, callArg.CallArg.Pure.Bytes...)
		}
	}

	return args
}
