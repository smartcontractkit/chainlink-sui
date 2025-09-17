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
	tokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_token_pool"
	mcmsops "github.com/smartcontractkit/chainlink-sui/deployment/ops/mcms"
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
	seqReports := make([]cld_ops.Report[any, any], 0)

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
		mcmsSeqReport, err := cld_ops.ExecuteSequence(e.OperationsBundle, mcmsops.DeployMCMSSequence, deps, mcmsops.DeployMCMSSeqInput{})
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy CCIP for Sui chain %d: %w", chainSel, err)
		}
		seqReports = append(seqReports, mcmsSeqReport.ExecutionReports...)

		// save MCMS address to the addressbook
		typeAndVersionMCMS := cldf.NewTypeAndVersion(deployment.SuiMcmsPackageIDType, deployment.Version1_0_0)
		err = ab.Save(chainSel, mcmsSeqReport.Output.PackageId, typeAndVersionMCMS)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to save MCMS address %s for Sui chain %d: %w", mcmsSeqReport.Output.PackageId, chainSel, err)
		}

		// Deploy Router
		// TODO: Maybe make this part of CCIP sequence
		routerReport, err := cld_ops.ExecuteOperation(e.OperationsBundle, routerops.DeployCCIPRouterOp, deps, routerops.DeployCCIPRouterInput{
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
		ccipSeqInput.McmsPackageId = mcmsSeqReport.Output.PackageId
		ccipSeqInput.McmsOwner = signerAddr

		ccipSeqReport, err := cld_ops.ExecuteSequence(e.OperationsBundle, ccipops.DeployAndInitCCIPSequence, deps, ccipSeqInput)
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
		typeAndVersionCCIPFeeQuoterCapIDRef := cldf.NewTypeAndVersion(deployment.SuiFeeQuoterCapType, deployment.Version1_0_0)
		err = ab.Save(chainSel, ccipSeqReport.Output.Objects.FeeQuoterCapObjectId, typeAndVersionCCIPFeeQuoterCapIDRef)
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

		ccipOnRampSeqReport, err := cld_ops.ExecuteSequence(e.OperationsBundle, onrampops.DeployAndInitCCIPOnRampSequence, deps, ccipOnRampSeqInput)
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

		// save onRampStateID address to the addressbook
		typeAndVersionOnRampStateID := cldf.NewTypeAndVersion(deployment.SuiOnRampStateObjectIDType, deployment.Version1_0_0)
		err = ab.Save(chainSel, ccipOnRampSeqReport.Output.Objects.StateObjectId, typeAndVersionOnRampStateID)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to save onRamp state object Id  %s for Sui chain %d: %w", ccipOnRampSeqReport.Output.Objects.StateObjectId, chainSel, err)
		}

		//  Run DeployAndInitCCIPOffRampSequence
		ccipOffRampSeqInput := deployment.DefaultOffRampSeqConfig
		// note: this is a regression, can't acess other chains state very cleanly
		onRampBytes := [][]byte{config.DestChainOnRampAddressBytes}

		// Inject dynamic values for deployment
		ccipOffRampSeqInput.CCIPPackageId = ccipSeqReport.Output.CCIPPackageId
		ccipOffRampSeqInput.MCMSPackageId = mcmsSeqReport.Output.PackageId

		ccipOffRampSeqInput.DestTransferCapId = ccipSeqReport.Output.Objects.DestTransferCapObjectId
		ccipOffRampSeqInput.FeeQuoterCapId = ccipSeqReport.Output.Objects.FeeQuoterCapObjectId
		ccipOffRampSeqInput.ChainSelector = suiChain.Selector
		ccipOffRampSeqInput.SourceChainSelectors = []uint64{
			config.ContractParamsPerChain[chainSel].DestChainSelector, // Ethereum, etc.
		}
		ccipOffRampSeqInput.SourceChainsOnRamp = onRampBytes

		ccipOffRampSeqReport, err := cld_ops.ExecuteSequence(e.OperationsBundle, offrampops.DeployAndInitCCIPOffRampSequence, deps, ccipOffRampSeqInput)
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

		// save offRamp ownerCapID to the addressbook
		typeAndVersionOffRampOwnerCapID := cldf.NewTypeAndVersion(deployment.SuiOffRampOwnerCapObjectIDType, deployment.Version1_0_0)
		err = ab.Save(chainSel, ccipOffRampSeqReport.Output.Objects.OwnerCapId, typeAndVersionOffRampOwnerCapID)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to save offRamp ObjectCapId address %s for Sui chain %d: %w", ccipOffRampSeqReport.Output.CCIPOffRampPackageId, chainSel, err)
		}

		// save offRamp stateObjectId to the addressbook
		typeAndVersionOffRampObjectStateID := cldf.NewTypeAndVersion(deployment.SuiOffRampStateObjectIDType, deployment.Version1_0_0)
		err = ab.Save(chainSel, ccipOffRampSeqReport.Output.Objects.StateObjectId, typeAndVersionOffRampObjectStateID)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to save offRamp StateObjectId %s for Sui chain %d: %w", ccipOffRampSeqReport.Output.Objects.StateObjectId, chainSel, err)
		}

		// Deploy CCIP TokenPool
		deployTp, err := cld_ops.ExecuteOperation(e.OperationsBundle, tokenpoolops.DeployCCIPTokenPoolOp, deps,
			tokenpoolops.TokenPoolDeployInput{
				CCIPPackageId:    ccipSeqReport.Output.CCIPPackageId,
				MCMSAddress:      mcmsSeqReport.Output.PackageId,
				MCMSOwnerAddress: signerAddr,
			})
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy TokenPool for Sui chain %d: %w", chainSel, err)
		}

		// save tokenPool address in addressbook
		typeAndVersionTokenPoolID := cldf.NewTypeAndVersion(deployment.SuiTokenPoolType, deployment.Version1_0_0)
		err = ab.Save(chainSel, deployTp.Output.PackageId, typeAndVersionTokenPoolID)
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
