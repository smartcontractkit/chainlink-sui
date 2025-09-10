package changesets

import (
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/deployment"
	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
	ccipops "github.com/smartcontractkit/chainlink-sui/ops/ccip"
	offrampops "github.com/smartcontractkit/chainlink-sui/ops/ccip_offramp"
	onrampops "github.com/smartcontractkit/chainlink-sui/ops/ccip_onramp"
	routerops "github.com/smartcontractkit/chainlink-sui/ops/ccip_router"
	tokenpoolops "github.com/smartcontractkit/chainlink-sui/ops/ccip_token_pool"
	mcmsops "github.com/smartcontractkit/chainlink-sui/ops/mcms"
)

type DeploySuiChainConfig struct {
	ContractParamsPerChain      map[uint64]ChainContractParams
	DestChainOnRampAddressBytes []byte
}

type ChainContractParams struct {
	DestChainSelector uint64
	FeeQuoterParams   ccipops.InitFeeQuoterInput
}

var _ cldf.ChangeSetV2[DeploySuiChainConfig] = DeploySuiChain{}

// DeploySuiChain deploys Sui chain packages and modules
type DeploySuiChain struct{}

// Apply implements deployment.ChangeSetV2.
func (d DeploySuiChain) Apply(e cldf.Environment, config DeploySuiChainConfig) (cldf.ChangesetOutput, error) {
	ab := cldf.NewMemoryAddressBook()
	seqReports := make([]operations.Report[any, any], 0)

	for chainSel := range config.ContractParamsPerChain {
		suiChain := e.BlockChains.SuiChains()[chainSel]
		signer := suiChain.Signer
		signerAddr, err := signer.GetAddress()
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}

		deps := sui_ops.OpTxDeps{
			Client: suiChain.Client,
			Signer: signer,
			GetCallOpts: func() *bind.CallOpts {
				b := uint64(400_000_000)
				return &bind.CallOpts{
					WaitForExecution: true,
					GasBudget:        &b,
				}
			},
		}

		// Deploy MCMS
		mcmsSeqReport, err := operations.ExecuteSequence(e.OperationsBundle, mcmsops.DeployMCMSSequence, deps, cld_ops.EmptyInput{})
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy CCIP for Sui chain %d: %w", chainSel, err)
		}
		seqReports = append(seqReports, mcmsSeqReport.ExecutionReports...)

		// save MCMs address to the addressbook
		typeAndVersionMCMS := cldf.NewTypeAndVersion(deployment.SuiMCMSType, deployment.Version1_0_0)
		err = ab.Save(chainSel, mcmsSeqReport.Output.PackageId, typeAndVersionMCMS)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to save MCMS address %s for Sui chain %d: %w", mcmsSeqReport.Output.PackageId, chainSel, err)
		}

		// Deploy Router
		// TODO: Maybe make this part of CCIP sequence
		routerReport, err := operations.ExecuteOperation(e.OperationsBundle, routerops.DeployCCIPRouterOp, deps, routerops.DeployCCIPRouterInput{
			McmsPackageId: mcmsSeqReport.Output.PackageId,
			McmsOwner:     signerAddr,
		})
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy CCIP Router for Sui chain %d: %w", chainSel, err)
		}

		// save Router address to the addressbook
		typeAndVersionRouter := cldf.NewTypeAndVersion(deployment.SuiCCIPRouterType, deployment.Version1_0_0)
		err = ab.Save(chainSel, routerReport.Output.PackageId, typeAndVersionRouter)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to save Router address %s for Sui chain %d: %w", routerReport.Output.PackageId, chainSel, err)
		}

		// DeployAndInitCCIPSequence

		// Inject chain-specific and runtime values
		ccipSeqInput := deployment.DefaultCCIPSeqConfig
		ccipSeqInput.LinkTokenCoinMetadataObjectId = config.ContractParamsPerChain[chainSel].FeeQuoterParams.LinkTokenCoinMetadataObjectId
		ccipSeqInput.LocalChainSelector = chainSel
		ccipSeqInput.DestChainSelector = config.ContractParamsPerChain[chainSel].DestChainSelector
		ccipSeqInput.DeployCCIPInput.McmsPackageId = mcmsSeqReport.Output.PackageId
		ccipSeqInput.DeployCCIPInput.McmsOwner = signerAddr

		ccipSeqReport, err := operations.ExecuteSequence(e.OperationsBundle, ccipops.DeployAndInitCCIPSequence, deps, ccipSeqInput)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy CCIP for Sui chain %d: %w", chainSel, err)
		}
		seqReports = append(seqReports, ccipSeqReport.ExecutionReports...)

		// save CCIP address to the addressbook
		typeAndVersionCCIP := cldf.NewTypeAndVersion(deployment.SuiCCIPType, deployment.Version1_0_0)
		err = ab.Save(chainSel, ccipSeqReport.Output.CCIPPackageId, typeAndVersionCCIP)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to save CCIP address %s for Sui chain %d: %w", ccipSeqReport.Output.CCIPPackageId, chainSel, err)
		}

		// save CCIP ObjectRef address to the addressbook
		typeAndVersionCCIPObjectRef := cldf.NewTypeAndVersion(deployment.SuiCCIPObjectRefType, deployment.Version1_0_0)
		err = ab.Save(chainSel, ccipSeqReport.Output.Objects.CCIPObjectRefObjectId, typeAndVersionCCIPObjectRef)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to save CCIP objectRef Id %s for Sui chain %d: %w", ccipSeqReport.Output.Objects.CCIPObjectRefObjectId, chainSel, err)
		}

		// save CCIP FeeQuoterCapObjectId address to the addressbook
		typeAndVersionCCIPFeeQuoterCapIdRef := cldf.NewTypeAndVersion(deployment.SuiFeeQuoterCapType, deployment.Version1_0_0)
		err = ab.Save(chainSel, ccipSeqReport.Output.Objects.FeeQuoterCapObjectId, typeAndVersionCCIPFeeQuoterCapIdRef)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to save CCIP FeeQuoter CapId Id %s for Sui chain %d: %w", ccipSeqReport.Output.Objects.FeeQuoterCapObjectId, chainSel, err)
		}

		// No need to store rn
		// save CCIP TransferCapId address to the addressbook
		// typeAndVersionTransferCapId := cldf.NewTypeAndVersion(deployment.SuiCCIPTransferCapIdType, deployment.Version1_0_0)
		// err = ab.Save(chainSel, ccipSeqReport.Output.Objects.SourceTransferCapObjectId, typeAndVersionTransferCapId)
		// if err != nil {
		// 	return cldf.ChangesetOutput{}, fmt.Errorf("failed to save CCIP TransferCapId Id %s for Sui chain %d: %w", ccipSeqReport.Output.Objects.SourceTransferCapObjectId, chainSel, err)
		// }

		// // save CCIP NonceManagerCapObjectId address to the addressbook
		// typeAndVersionNonceManagerCapObjectId := cldf.NewTypeAndVersion(deployment.SuiCCIPObjectRefType, deployment.Version1_0_0)
		// err = ab.Save(chainSel, ccipSeqReport.Output.Objects.NonceManagerCapObjectId, typeAndVersionNonceManagerCapObjectId)
		// if err != nil {
		// 	return cldf.ChangesetOutput{}, fmt.Errorf("failed to save CCIP objectRef Id %s for Sui chain %d: %w", ccipSeqReport.Output.Objects.CCIPObjectRefObjectId, chainSel, err)
		// }

		// Run DeployAndInitCCIPOnRampSequence
		ccipOnRampSeqInput := onrampops.DeployAndInitCCIPOnRampSeqInput{
			DeployCCIPOnRampInput: onrampops.DeployCCIPOnRampInput{
				CCIPPackageId:      ccipSeqReport.Output.CCIPPackageId,
				MCMSPackageId:      mcmsSeqReport.Output.PackageId,
				MCMSOwnerPackageId: signerAddr,
			},
			OnRampInitializeInput: onrampops.OnRampInitializeInput{
				NonceManagerCapId:         ccipSeqReport.Output.Objects.NonceManagerCapObjectId,   // this is from NonceManager init Op
				SourceTransferCapId:       ccipSeqReport.Output.Objects.SourceTransferCapObjectId, // this is from CCIP package publish
				ChainSelector:             suiChain.Selector,
				FeeAggregator:             signerAddr,
				AllowListAdmin:            signerAddr,
				DestChainSelectors:        []uint64{config.ContractParamsPerChain[chainSel].DestChainSelector}, // TODOD add this in input instead of hardcoding
				DestChainEnabled:          []bool{true},
				DestChainAllowListEnabled: []bool{true},
			},
			ApplyDestChainConfigureOnRampInput: onrampops.ApplyDestChainConfigureOnRampInput{
				DestChainSelector:         []uint64{config.ContractParamsPerChain[chainSel].DestChainSelector},
				DestChainEnabled:          []bool{true},
				DestChainAllowListEnabled: []bool{false},
			},
			ApplyAllowListUpdatesInput: onrampops.ApplyAllowListUpdatesInput{
				DestChainSelector:             []uint64{config.ContractParamsPerChain[chainSel].DestChainSelector},
				DestChainAllowListEnabled:     []bool{false},
				DestChainAddAllowedSenders:    [][]string{{}},
				DestChainRemoveAllowedSenders: [][]string{{}},
			},
		}

		ccipOnRampSeqReport, err := operations.ExecuteSequence(e.OperationsBundle, onrampops.DeployAndInitCCIPOnRampSequence, deps, ccipOnRampSeqInput)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy CCIP for Sui chain %d: %w", chainSel, err)
		}
		seqReports = append(seqReports, ccipOnRampSeqReport.ExecutionReports...)

		// save onRamp address to the addressbook
		typeAndVersionOnRamp := cldf.NewTypeAndVersion(deployment.SuiOnRampType, deployment.Version1_0_0)
		err = ab.Save(chainSel, ccipOnRampSeqReport.Output.CCIPOnRampPackageId, typeAndVersionOnRamp)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to save onRamp address %s for Sui chain %d: %w", ccipOnRampSeqReport.Output.CCIPOnRampPackageId, chainSel, err)
		}

		// save onRampStateId address to the addressbook
		typeAndVersionOnRampStateId := cldf.NewTypeAndVersion(deployment.SuiOnRampStateObjectIdType, deployment.Version1_0_0)
		err = ab.Save(chainSel, ccipOnRampSeqReport.Output.Objects.StateObjectId, typeAndVersionOnRampStateId)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to save onRamp state object Id  %s for Sui chain %d: %w", ccipOnRampSeqReport.Output.Objects.StateObjectId, chainSel, err)
		}

		//  Run DeployAndInitCCIPOffRampSequence
		ccipOffRampSeqInput := deployment.DefaultOffRampSeqConfig
		// note: this is a regression, can't acess other chains state very cleanly
		onRampBytes := [][]byte{config.DestChainOnRampAddressBytes}

		// Inject dynamic values for deployment
		ccipOffRampSeqInput.DeployCCIPOffRampInput.CCIPPackageId = ccipSeqReport.Output.CCIPPackageId
		ccipOffRampSeqInput.DeployCCIPOffRampInput.MCMSPackageId = mcmsSeqReport.Output.PackageId

		ccipOffRampSeqInput.InitializeOffRampInput.DestTransferCapId = ccipSeqReport.Output.Objects.DestTransferCapObjectId
		ccipOffRampSeqInput.InitializeOffRampInput.FeeQuoterCapId = ccipSeqReport.Output.Objects.FeeQuoterCapObjectId
		ccipOffRampSeqInput.InitializeOffRampInput.ChainSelector = suiChain.Selector
		ccipOffRampSeqInput.InitializeOffRampInput.SourceChainSelectors = []uint64{
			config.ContractParamsPerChain[chainSel].DestChainSelector, // Ethereum, etc.
		}
		ccipOffRampSeqInput.InitializeOffRampInput.SourceChainsOnRamp = onRampBytes

		ccipOffRampSeqReport, err := operations.ExecuteSequence(e.OperationsBundle, offrampops.DeployAndInitCCIPOffRampSequence, deps, ccipOffRampSeqInput)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy CCIP for Sui chain %d: %w", chainSel, err)
		}
		seqReports = append(seqReports, ccipOffRampSeqReport.ExecutionReports...)

		// save offRamp address to the addressbook
		typeAndVersionOffRamp := cldf.NewTypeAndVersion(deployment.SuiOffRampType, deployment.Version1_0_0)
		err = ab.Save(chainSel, ccipOffRampSeqReport.Output.CCIPOffRampPackageId, typeAndVersionOffRamp)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to save offRamp address %s for Sui chain %d: %w", ccipOffRampSeqReport.Output.CCIPOffRampPackageId, chainSel, err)
		}

		// save offRamp ownerCapId to the addressbook
		typeAndVersionOffRampOwnerCapId := cldf.NewTypeAndVersion(deployment.SuiOffRampOwnerCapObjectIdType, deployment.Version1_0_0)
		err = ab.Save(chainSel, ccipOffRampSeqReport.Output.Objects.OwnerCapId, typeAndVersionOffRampOwnerCapId)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to save offRamp ObjectCapId address %s for Sui chain %d: %w", ccipOffRampSeqReport.Output.CCIPOffRampPackageId, chainSel, err)
		}

		// save offRamp stateObjectId to the addressbook
		typeAndVersionOffRampObjectStateId := cldf.NewTypeAndVersion(deployment.SuiOffRampStateObjectIdType, deployment.Version1_0_0)
		err = ab.Save(chainSel, ccipOffRampSeqReport.Output.Objects.StateObjectId, typeAndVersionOffRampObjectStateId)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to save offRamp StateObjectId %s for Sui chain %d: %w", ccipOffRampSeqReport.Output.Objects.StateObjectId, chainSel, err)
		}

		// Deploy CCIP TokenPool
		deployTp, err := operations.ExecuteOperation(e.OperationsBundle, tokenpoolops.DeployCCIPTokenPoolOp, deps,
			tokenpoolops.TokenPoolDeployInput{
				CCIPPackageId:    ccipSeqReport.Output.CCIPPackageId,
				MCMSAddress:      mcmsSeqReport.Output.PackageId,
				MCMSOwnerAddress: signerAddr,
			})
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy TokenPool for Sui chain %d: %w", chainSel, err)
		}

		// save tokenPool address in addressbook
		typeAndVersionTokenPoolId := cldf.NewTypeAndVersion(deployment.SuiTokenPoolType, deployment.Version1_0_0)
		err = ab.Save(chainSel, deployTp.Output.PackageId, typeAndVersionTokenPoolId)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to save offRamp StateObjectId %s for Sui chain %d: %w", deployTp.Output.PackageId, chainSel, err)
		}
	}

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
