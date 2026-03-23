package changesets

import (
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
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
	seqReports := make([]cld_ops.Report[any, any], 0)

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
			b := uint64(1_000_000_000)
			return &bind.CallOpts{
				WaitForExecution: true,
				GasBudget:        &b,
			}
		},
		SuiRPC: suiChain.URL,
	}

	// in case the registry is not loaded with all operations. Needed to build accept ownership proposals
	for i := range opregistry.AllOperations {
		cld_ops.RegisterOperation(e.OperationsBundle.OperationRegistry, opregistry.AllOperations[i])
	}

	suiState, err := deployment.LoadOnchainStatesui(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	state := suiState[config.SuiChainSelector]

	mcmsPackageId := state.MCMSPackageID
	// If MCMS is not deployed, deploy it
	if mcmsPackageId == "" {
		mcmsReport, err := cld_ops.ExecuteSequence(e.OperationsBundle, mcmsops.DeployMCMSSequence, deps, mcmsops.DeployMCMSSeqInput{
			ChainSelector: config.SuiChainSelector,
		})
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy MCMS for Sui chain %d: %w", config.SuiChainSelector, err)
		}

		err = storeMCMSInAddressBook(ab, config.SuiChainSelector, mcmsReport.Output, false)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to store MCMS in address book for Sui chain %d: %w", config.SuiChainSelector, err)
		}

		mcmsPackageId = mcmsReport.Output.PackageId
	}

	// Deploy Router
	// TODO: Maybe make this part of CCIP sequence
	routerReport, err := cld_ops.ExecuteOperation(e.OperationsBundle, routerops.DeployCCIPRouterOp, deps, routerops.DeployCCIPRouterInput{
		McmsPackageId: mcmsPackageId,
		McmsOwner:     signerAddr,
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy CCIP Router for Sui chain %d: %w", config.SuiChainSelector, err)
	}

	// Transfer ownership of Router to MCMS
	_, err = cld_ops.ExecuteOperation(e.OperationsBundle, routerops.TransferOwnershipOp, deps, routerops.TransferOwnershipInput{
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

	typeAndVersionRouterOwnerCap := cldf.NewTypeAndVersion(deployment.SuiCCIPRouterOwnerCapObjectIDType, deployment.Version1_0_0)
	err = ab.Save(config.SuiChainSelector, routerReport.Output.Objects.OwnerCapObjectId, typeAndVersionRouterOwnerCap)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save Router owner cap object Id %s for Sui chain %d: %w", routerReport.Output.Objects.OwnerCapObjectId, config.SuiChainSelector, err)
	}

	typeAndVersionRouterUpgradeCapId := cldf.NewTypeAndVersion(deployment.SuiRouterUpgradeCapObjectIDType, deployment.Version1_0_0)
	err = ab.Save(config.SuiChainSelector, routerReport.Output.Objects.UpgradeCapObjectId, typeAndVersionRouterUpgradeCapId)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save RouterUpgradeCapId  %s for Sui chain %d: %w", routerReport.Output.Objects.UpgradeCapObjectId, config.SuiChainSelector, err)
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

	ccipSeqReport, err := cld_ops.ExecuteSequence(e.OperationsBundle, ccipops.DeployAndInitCCIPSequence, deps, ccipSeqInput)
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

	ccipOnRampSeqReport, err := cld_ops.ExecuteSequence(e.OperationsBundle, onrampops.DeployAndInitCCIPOnRampSequence, deps, ccipOnRampSeqInput)
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

	ccipOffRampSeqReport, err := cld_ops.ExecuteSequence(e.OperationsBundle, offrampops.DeployAndInitCCIPOffRampSequence, deps, ccipOffRampSeqInput)
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

	// TODO: This could return the accept ownership proposal instead of having a different changeset for it

	return cldf.ChangesetOutput{
		AddressBook: ab,
		Reports:     seqReports,
	}, nil
}

// TODO
// VerifyPreconditions implements deployment.ChangeSetV2.
func (d DeploySuiChain) VerifyPreconditions(e cldf.Environment, config DeploySuiChainConfig) error {
	return nil
}
