package changesets

import (
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/deployment"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	ccipops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip"
	offrampops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_offramp"
	onrampops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_onramp"
	routerops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_router"
	mcmsops "github.com/smartcontractkit/chainlink-sui/deployment/ops/mcms"
	opregistry "github.com/smartcontractkit/chainlink-sui/deployment/ops/registry"
	"github.com/smartcontractkit/mcms"
	suisdk "github.com/smartcontractkit/mcms/sdk/sui"
)

type DeploySuiChainConfig struct {
	SuiChainSelector              uint64
	DestChainSelector             uint64 // dest chain selector
	DestChainOnRampAddressBytes   []byte // onRamp of the destination chain we are connecting to
	LinkTokenCoinMetadataObjectId string // this defines the initial feeToken
}

var _ cldf.ChangeSetV2[DeploySuiChainConfig] = DeploySuiChain{}

// DeploySuiChain deploys Sui chain packages and modules
type DeploySuiChain struct{}

// Apply implements deployment.ChangeSetV2.
func (d DeploySuiChain) Apply(e cldf.Environment, config DeploySuiChainConfig) (cldf.ChangesetOutput, error) {
	ab := cldf.NewMemoryAddressBook()
	seqReports := make([]operations.Report[any, any], 0)

	suiChain := e.BlockChains.SuiChains()[config.SuiChainSelector]
	signer := suiChain.Signer
	signerAddr, err := signer.GetAddress()
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	deps := sui_ops.OpTxDeps{
		Client: suiChain.Client,
		Signer: signer,
		GetCallOpts: func() *bind.CallOpts {
			b := uint64(500_000_000)
			return &bind.CallOpts{
				WaitForExecution: true,
				GasBudget:        &b,
			}
		},
	}

	// in case the registry is not loaded with all operations. Needed to build accept ownership proposals
	ops := make([]*cld_ops.Operation[any, any, any], len(opregistry.AllOperations))
	for i := range opregistry.AllOperations {
		ops[i] = &opregistry.AllOperations[i]
		cld_ops.RegisterOperation(e.OperationsBundle.OperationRegistry, &opregistry.AllOperations[i])
	}

	suiState, err := deployment.LoadOnchainStatesui(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	state := suiState[config.SuiChainSelector]

	mcmsPackageId := state.MCMSPackageID
	mcmsStateObjID := state.MCMSStateObjectID
	timelockObjID := state.MCMSTimelockObjectID
	accountObjID := state.MCMSAccountStateObjectID
	registryObjID := state.MCMSRegistryObjectID
	// If MCMS is not deployed, deploy it
	if mcmsPackageId == "" {
		mcmsReport, err := cld_ops.ExecuteSequence(e.OperationsBundle, mcmsops.DeployMCMSSequence, deps, mcmsops.DeployMCMSSeqInput{
			ChainSelector: config.SuiChainSelector,
		})
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy MCMS for Sui chain %d: %w", config.SuiChainSelector, err)
		}

		err = storeMCMSInAddressBook(ab, config.SuiChainSelector, mcmsReport.Output)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to store MCMS in address book for Sui chain %d: %w", config.SuiChainSelector, err)
		}

		mcmsPackageId = mcmsReport.Output.PackageId
		mcmsStateObjID = mcmsReport.Output.Objects.McmsMultisigStateObjectId
		timelockObjID = mcmsReport.Output.Objects.TimelockObjectId
		accountObjID = mcmsReport.Output.Objects.McmsAccountStateObjectId
		registryObjID = mcmsReport.Output.Objects.McmsRegistryObjectId
	}

	// Deploy Router
	// TODO: Maybe make this part of CCIP sequence
	routerReport, err := operations.ExecuteOperation(e.OperationsBundle, routerops.DeployCCIPRouterOp, deps, routerops.DeployCCIPRouterInput{
		McmsPackageId: mcmsPackageId,
		McmsOwner:     signerAddr,
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy CCIP Router for Sui chain %d: %w", config.SuiChainSelector, err)
	}

	// Transfer ownership of Router to MCMS
	_, err = operations.ExecuteOperation(e.OperationsBundle, routerops.TransferOwnershipOp, deps, routerops.TransferOwnershipInput{
		RouterPackageId:     routerReport.Output.PackageId,
		RouterStateObjectId: routerReport.Output.Objects.RouterStateObjectId,
		OwnerCapObjectId:    routerReport.Output.Objects.OwnerCapObjectId,
		NewOwner:            mcmsPackageId,
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to execute ownership transfer to MCMS Router for Sui chain %d: %w", config.SuiChainSelector, err)
	}

	// save Router address to the addressbook
	typeAndVersionRouter := cldf.NewTypeAndVersion(deployment.SuiCCIPRouterType, deployment.Version1_0_0)
	err = ab.Save(config.SuiChainSelector, routerReport.Output.PackageId, typeAndVersionRouter)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save Router address %s for Sui chain %d: %w", routerReport.Output.PackageId, config.SuiChainSelector, err)
	}

	typeAndVersionRouterObject := cldf.NewTypeAndVersion(deployment.SuiCCIPRouterStateObjectType, deployment.Version1_0_0)
	err = ab.Save(config.SuiChainSelector, routerReport.Output.Objects.RouterStateObjectId, typeAndVersionRouterObject)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save Router state object Id %s for Sui chain %d: %w", routerReport.Output.Objects.RouterStateObjectId, config.SuiChainSelector, err)
	}

	// --------------------------
	// CCIP SEQUENCE
	// --------------------------
	// DeployAndInitCCIPSequence
	// Inject chain-specific and runtime values
	ccipSeqInput := deployment.DefaultCCIPSeqConfig
	ccipSeqInput.LinkTokenCoinMetadataObjectId = config.LinkTokenCoinMetadataObjectId
	ccipSeqInput.LocalChainSelector = config.SuiChainSelector
	ccipSeqInput.DestChainSelector = config.DestChainSelector
	ccipSeqInput.DeployCCIPInput.McmsPackageId = mcmsPackageId
	ccipSeqInput.DeployCCIPInput.McmsOwner = signerAddr

	ccipSeqReport, err := operations.ExecuteSequence(e.OperationsBundle, ccipops.DeployAndInitCCIPSequence, deps, ccipSeqInput)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy CCIP for Sui chain %d: %w", config.SuiChainSelector, err)
	}
	seqReports = append(seqReports, ccipSeqReport.ExecutionReports...)

	// save CCIP address to the addressbook
	typeAndVersionCCIP := cldf.NewTypeAndVersion(deployment.SuiCCIPType, deployment.Version1_0_0)
	err = ab.Save(config.SuiChainSelector, ccipSeqReport.Output.CCIPPackageId, typeAndVersionCCIP)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save CCIP address %s for Sui chain %d: %w", ccipSeqReport.Output.CCIPPackageId, config.SuiChainSelector, err)
	}

	// save CCIP ObjectRef address to the addressbook
	typeAndVersionCCIPObjectRef := cldf.NewTypeAndVersion(deployment.SuiCCIPObjectRefType, deployment.Version1_0_0)
	err = ab.Save(config.SuiChainSelector, ccipSeqReport.Output.Objects.CCIPObjectRefObjectId, typeAndVersionCCIPObjectRef)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save CCIP objectRef Id %s for Sui chain %d: %w", ccipSeqReport.Output.Objects.CCIPObjectRefObjectId, config.SuiChainSelector, err)
	}

	// save CCIP FeeQuoterCapObjectId address to the addressbook
	typeAndVersionCCIPFeeQuoterCapIdRef := cldf.NewTypeAndVersion(deployment.SuiFeeQuoterCapType, deployment.Version1_0_0)
	err = ab.Save(config.SuiChainSelector, ccipSeqReport.Output.Objects.FeeQuoterCapObjectId, typeAndVersionCCIPFeeQuoterCapIdRef)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save CCIP FeeQuoter CapId Id %s for Sui chain %d: %w", ccipSeqReport.Output.Objects.FeeQuoterCapObjectId, config.SuiChainSelector, err)
	}

	// save CCIP ObjectRef address to the addressbook
	typeAndVersionCCIPOwnerCapObjectId := cldf.NewTypeAndVersion(deployment.SuiCCIPOwnerCapObjectIDType, deployment.Version1_0_0)
	err = ab.Save(config.SuiChainSelector, ccipSeqReport.Output.Objects.OwnerCapObjectId, typeAndVersionCCIPOwnerCapObjectId)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save CCIP ownerCapObjectId %s for Sui chain %d: %w", ccipSeqReport.Output.Objects.OwnerCapObjectId, config.SuiChainSelector, err)
	}

	typeAndVersionCCIPUpgradeCapObjectId := cldf.NewTypeAndVersion(deployment.SuiCCIPUpgradeCapObjectIDType, deployment.Version1_0_0)
	err = ab.Save(config.SuiChainSelector, ccipSeqReport.Output.Objects.UpgradeCapObjectId, typeAndVersionCCIPUpgradeCapObjectId)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save CCIP UpgradeCapObjectId %s for Sui chain %d: %w", ccipSeqReport.Output.Objects.UpgradeCapObjectId, config.SuiChainSelector, err)
	}

	// No need to store rn
	// save CCIP TransferCapId address to the addressbook
	// typeAndVersionTransferCapId := cldf.NewTypeAndVersion(deployment.SuiCCIPTransferCapIdType, deployment.Version1_0_0)
	// err = ab.Save(config.SuiChainSelector, ccipSeqReport.Output.Objects.SourceTransferCapObjectId, typeAndVersionTransferCapId)
	// if err != nil {
	// 	return cldf.ChangesetOutput{}, fmt.Errorf("failed to save CCIP TransferCapId Id %s for Sui chain %d: %w", ccipSeqReport.Output.Objects.SourceTransferCapObjectId, config.SuiChainSelector, err)
	// }

	// // save CCIP NonceManagerCapObjectId address to the addressbook
	// typeAndVersionNonceManagerCapObjectId := cldf.NewTypeAndVersion(deployment.SuiCCIPObjectRefType, deployment.Version1_0_0)
	// err = ab.Save(config.SuiChainSelector, ccipSeqReport.Output.Objects.NonceManagerCapObjectId, typeAndVersionNonceManagerCapObjectId)
	// if err != nil {
	// 	return cldf.ChangesetOutput{}, fmt.Errorf("failed to save CCIP objectRef Id %s for Sui chain %d: %w", ccipSeqReport.Output.Objects.CCIPObjectRefObjectId, config.SuiChainSelector, err)
	// }

	// --------------------------
	// CCIP ONRAMP SEQUENCE
	// --------------------------
	// Run DeployAndInitCCIPOnRampSequence
	ccipOnRampSeqInput := deployment.DefaultOnRampSeqConfig

	ccipOnRampSeqInput.DeployCCIPOnRampInput.CCIPPackageId = ccipSeqReport.Output.CCIPPackageId
	ccipOnRampSeqInput.DeployCCIPOnRampInput.MCMSPackageId = mcmsPackageId
	ccipOnRampSeqInput.DeployCCIPOnRampInput.MCMSOwnerPackageId = signerAddr
	ccipOnRampSeqInput.OnRampInitializeInput.NonceManagerCapId = ccipSeqReport.Output.Objects.NonceManagerCapObjectId
	ccipOnRampSeqInput.OnRampInitializeInput.SourceTransferCapId = ccipSeqReport.Output.Objects.SourceTransferCapObjectId
	ccipOnRampSeqInput.OnRampInitializeInput.ChainSelector = suiChain.Selector
	ccipOnRampSeqInput.OnRampInitializeInput.FeeAggregator = signerAddr
	ccipOnRampSeqInput.OnRampInitializeInput.AllowListAdmin = signerAddr
	ccipOnRampSeqInput.OnRampInitializeInput.DestChainSelectors = []uint64{config.DestChainSelector}
	ccipOnRampSeqInput.OnRampInitializeInput.DestChainRouters = []string{routerReport.Output.PackageId}
	ccipOnRampSeqInput.ApplyDestChainConfigureOnRampInput.DestChainSelector = []uint64{config.DestChainSelector}
	ccipOnRampSeqInput.ApplyAllowListUpdatesInput.DestChainSelector = []uint64{config.DestChainSelector}
	ccipOnRampSeqInput.ApplyDestChainConfigureOnRampInput.DestChainRouters = []string{routerReport.Output.PackageId}
	ccipOnRampSeqInput.ApplyDestChainConfigureOnRampInput.CCIPObjectRefId = ccipSeqReport.Output.Objects.CCIPObjectRefObjectId

	ccipOnRampSeqReport, err := operations.ExecuteSequence(e.OperationsBundle, onrampops.DeployAndInitCCIPOnRampSequence, deps, ccipOnRampSeqInput)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy CCIP for Sui chain %d: %w", config.SuiChainSelector, err)
	}
	seqReports = append(seqReports, ccipOnRampSeqReport.ExecutionReports...)

	// save onRamp address to the addressbook
	typeAndVersionOnRamp := cldf.NewTypeAndVersion(deployment.SuiOnRampType, deployment.Version1_0_0)
	err = ab.Save(config.SuiChainSelector, ccipOnRampSeqReport.Output.CCIPOnRampPackageId, typeAndVersionOnRamp)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save onRamp address %s for Sui chain %d: %w", ccipOnRampSeqReport.Output.CCIPOnRampPackageId, config.DestChainSelector, err)
	}

	// save onRampStateId address to the addressbook
	typeAndVersionOnRampStateId := cldf.NewTypeAndVersion(deployment.SuiOnRampStateObjectIDType, deployment.Version1_0_0)
	err = ab.Save(config.SuiChainSelector, ccipOnRampSeqReport.Output.Objects.StateObjectId, typeAndVersionOnRampStateId)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save onRamp state object Id  %s for Sui chain %d: %w", ccipOnRampSeqReport.Output.Objects.StateObjectId, config.DestChainSelector, err)
	}

	// save OnRampOwnerCapObjectID to addressbook
	typeAndVersionOnRampOwnerCapObjectId := cldf.NewTypeAndVersion(deployment.SuiOnRampOwnerCapObjectIDType, deployment.Version1_0_0)
	err = ab.Save(config.SuiChainSelector, ccipOnRampSeqReport.Output.Objects.OwnerCapObjectId, typeAndVersionOnRampOwnerCapObjectId)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save OnRampOwnerCapObjectID  %s for Sui chain %d: %w", ccipOnRampSeqReport.Output.Objects.StateObjectId, config.DestChainSelector, err)
	}

	// save OnRampUpgradeCapId to addressbook
	typeAndVersionOnRampUpgradeCapId := cldf.NewTypeAndVersion(deployment.SuiOnRampUpgradeCapObjectIDType, deployment.Version1_0_0)
	err = ab.Save(config.SuiChainSelector, ccipOnRampSeqReport.Output.Objects.UpgradeCapObjectId, typeAndVersionOnRampUpgradeCapId)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save OnRampUpgradeCapId  %s for Sui chain %d: %w", ccipOnRampSeqReport.Output.Objects.StateObjectId, config.DestChainSelector, err)
	}

	// --------------------------
	// CCIP OFFRAMP SEQUENCE
	// --------------------------
	//  Run DeployAndInitCCIPOffRampSequence
	ccipOffRampSeqInput := deployment.DefaultOffRampSeqConfig
	// note: this is a regression, can't acess other chains state very cleanly
	onRampBytes := [][]byte{config.DestChainOnRampAddressBytes}

	// Inject dynamic values for deployment
	ccipOffRampSeqInput.CCIPObjectRefId = ccipSeqReport.Output.Objects.CCIPObjectRefObjectId
	ccipOffRampSeqInput.DeployCCIPOffRampInput.CCIPPackageId = ccipSeqReport.Output.CCIPPackageId
	ccipOffRampSeqInput.DeployCCIPOffRampInput.MCMSPackageId = mcmsPackageId

	ccipOffRampSeqInput.InitializeOffRampInput.DestTransferCapId = ccipSeqReport.Output.Objects.DestTransferCapObjectId
	ccipOffRampSeqInput.InitializeOffRampInput.FeeQuoterCapId = ccipSeqReport.Output.Objects.FeeQuoterCapObjectId
	ccipOffRampSeqInput.InitializeOffRampInput.ChainSelector = suiChain.Selector
	ccipOffRampSeqInput.InitializeOffRampInput.SourceChainSelectors = []uint64{
		config.DestChainSelector, // Ethereum, etc.
	}
	ccipOffRampSeqInput.InitializeOffRampInput.SourceChainsOnRamp = onRampBytes

	ccipOffRampSeqReport, err := operations.ExecuteSequence(e.OperationsBundle, offrampops.DeployAndInitCCIPOffRampSequence, deps, ccipOffRampSeqInput)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy CCIP for Sui chain %d: %w", config.SuiChainSelector, err)
	}
	seqReports = append(seqReports, ccipOffRampSeqReport.ExecutionReports...)

	// save offRamp address to the addressbook
	typeAndVersionOffRamp := cldf.NewTypeAndVersion(deployment.SuiOffRampType, deployment.Version1_0_0)
	err = ab.Save(config.SuiChainSelector, ccipOffRampSeqReport.Output.CCIPOffRampPackageId, typeAndVersionOffRamp)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save offRamp address %s for Sui chain %d: %w", ccipOffRampSeqReport.Output.CCIPOffRampPackageId, config.SuiChainSelector, err)
	}

	// save offRamp ownerCapId to the addressbook
	typeAndVersionOffRampOwnerCapId := cldf.NewTypeAndVersion(deployment.SuiOffRampOwnerCapObjectIDType, deployment.Version1_0_0)
	err = ab.Save(config.SuiChainSelector, ccipOffRampSeqReport.Output.Objects.OwnerCapId, typeAndVersionOffRampOwnerCapId)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save offRamp ObjectCapId address %s for Sui chain %d: %w", ccipOffRampSeqReport.Output.CCIPOffRampPackageId, config.SuiChainSelector, err)
	}

	// save offRamp stateObjectId to the addressbook
	typeAndVersionOffRampObjectStateId := cldf.NewTypeAndVersion(deployment.SuiOffRampStateObjectIDType, deployment.Version1_0_0)
	err = ab.Save(config.SuiChainSelector, ccipOffRampSeqReport.Output.Objects.StateObjectId, typeAndVersionOffRampObjectStateId)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save offRamp StateObjectId %s for Sui chain %d: %w", ccipOffRampSeqReport.Output.Objects.StateObjectId, config.SuiChainSelector, err)
	}

	// save OnRampUpgradeCapId to addressbook
	typeAndVersionOffRampUpgradeCapId := cldf.NewTypeAndVersion(deployment.SuiOffRampUpgradeCapObjectIDType, deployment.Version1_0_0)
	err = ab.Save(config.SuiChainSelector, ccipOffRampSeqReport.Output.Objects.UpgradeCapObjectId, typeAndVersionOffRampUpgradeCapId)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save OnRampUpgradeCapId  %s for Sui chain %d: %w", ccipOnRampSeqReport.Output.Objects.StateObjectId, config.DestChainSelector, err)
	}

	// Generate the proposal to accept the ownership of the deployed contracts
	proposalInput := mcmsops.ProposalGenerateInput{
		Defs: []operations.Definition{
			ccipops.AcceptOwnershipStateObjectOp.Def(),
			routerops.AcceptOwnershipOp.Def(),
			onrampops.AcceptOwnershipOnRampOp.Def(),
			offrampops.AcceptOwnershipOffRampOp.Def(),
		},
		Inputs: []any{
			ccipops.AcceptOwnershipStateObjectInput{
				CCIPPackageId:         ccipSeqReport.Output.CCIPPackageId,
				CCIPObjectRefObjectId: ccipSeqReport.Output.Objects.CCIPObjectRefObjectId,
			},
			routerops.AcceptOwnershipInput{
				RouterPackageId:     routerReport.Output.PackageId,
				RouterStateObjectId: routerReport.Output.Objects.RouterStateObjectId,
			},
			onrampops.AcceptOwnershipOnRampInput{
				OnRampPackageId: ccipOnRampSeqReport.Output.CCIPOnRampPackageId,
				CCIPObjectRefId: ccipSeqReport.Output.Objects.CCIPObjectRefObjectId,
				StateObjectId:   ccipOnRampSeqReport.Output.Objects.StateObjectId,
			},
			offrampops.AcceptOwnershipOffRampInput{
				OffRampPackageId:     ccipOffRampSeqReport.Output.CCIPOffRampPackageId,
				OffRampRefObjectId:   ccipSeqReport.Output.Objects.CCIPObjectRefObjectId,
				OffRampStateObjectId: ccipOffRampSeqReport.Output.Objects.StateObjectId,
			},
		},
		// MCMS related
		MmcsPackageID:  mcmsPackageId,
		McmsStateObjID: mcmsStateObjID,
		TimelockObjID:  timelockObjID,
		AccountObjID:   accountObjID,
		RegistryObjID:  registryObjID,

		// Proposal
		Role: suisdk.TimelockRoleProposer,

		ChainSelector: config.SuiChainSelector,
	}

	acceptOwnershipProposalReport, err := cld_ops.ExecuteSequence(e.OperationsBundle, mcmsops.MCMSDynamicProposalGenerateSeq, deps, proposalInput)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	return cldf.ChangesetOutput{
		AddressBook:           ab,
		Reports:               seqReports,
		MCMSTimelockProposals: []mcms.TimelockProposal{acceptOwnershipProposalReport.Output},
	}, nil
}

// TODO
// VerifyPreconditions implements deployment.ChangeSetV2.
func (d DeploySuiChain) VerifyPreconditions(e cldf.Environment, config DeploySuiChainConfig) error {
	return nil
}
