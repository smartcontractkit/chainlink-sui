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
	tokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_token_pool"
	mcmsops "github.com/smartcontractkit/chainlink-sui/deployment/ops/mcms"
)

var _ cldf.ChangeSetV2[DeploySuiChainConfig] = DeploySuiChain{}

// DeploySuiChain deploys Sui chain packages and modules
type DeploySuiChain struct{}

// Apply implements deployment.ChangeSetV2.
func (d DeploySuiChain) Apply(e cldf.Environment, config DeploySuiChainConfig) (cldf.ChangesetOutput, error) {
	ab := cldf.NewMemoryAddressBook()
	seqReports := make([]operations.Report[any, any], 0)

	for chainSel := range config.ContractParamsPerChain {
		suiChains := e.BlockChains.SuiChains()

		suiChain := suiChains[chainSel]
		suiSigner := suiChain.Signer

		signerAddr, err := suiSigner.GetAddress()
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}

		deps := sui_ops.OpTxDeps{
			Client: suiChain.Client,
			Signer: suiSigner,
			GetCallOpts: func() *bind.CallOpts {
				b := uint64(400_000_000)
				return &bind.CallOpts{
					WaitForExecution: true,
					GasBudget:        &b,
				}
			},
		}

		// TODO RA/Sish: Should MCMS and Router be part of the main CCIP deploy sequence?
		// Deploy MCMS
		mcmsSeqReport, err := operations.ExecuteSequence(e.OperationsBundle, mcmsops.DeployMCMSSequence, deps, cld_ops.EmptyInput{})
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy MCMS for Sui chain %d: %w", chainSel, err)
		}
		seqReports = append(seqReports, mcmsSeqReport.ExecutionReports...)

		// save MCMs address to the addressbook
		typeAndVersionMCMS := cldf.NewTypeAndVersion(deployment.SuiMCMSType, deployment.Version1_0_0)
		err = ab.Save(chainSel, mcmsSeqReport.Output.PackageId, typeAndVersionMCMS)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to save MCMS address %s for Sui chain %d: %w", mcmsSeqReport.Output.PackageId, chainSel, err)
		}

		// Deploy Router
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

		// Run DeployAndInitCCIpSequence
		ccipSeqInput := ccipops.DeployAndInitCCIPSeqInput{
			LinkTokenCoinMetadataObjectId: config.ContractParamsPerChain[chainSel].FeeQuoterParams.LinkTokenCoinMetadataObjectId,
			LocalChainSelector:            chainSel,
			DestChainSelector:             config.ContractParamsPerChain[chainSel].DestChainSelector,
			MaxFeeJuelsPerMsg:             config.ContractParamsPerChain[chainSel].FeeQuoterParams.MaxFeeJuelsPerMsg,
			TokenPriceStalenessThreshold:  config.ContractParamsPerChain[chainSel].FeeQuoterParams.TokenPriceStalenessThreshold,
			DeployCCIPInput: ccipops.DeployCCIPInput{
				McmsPackageId: mcmsSeqReport.Output.PackageId,
				McmsOwner:     signerAddr,
			},
			// Fee Quoter configuration
			AddMinFeeUsdCents:    []uint32{3000},
			AddMaxFeeUsdCents:    []uint32{30000},
			AddDeciBps:           []uint16{1000},
			AddDestGasOverhead:   []uint32{1000000},
			AddDestBytesOverhead: []uint32{1000},
			AddIsEnabled:         []bool{true},
			RemoveTokens:         []string{},

			// Fee Quoter destination chain configuration
			IsEnabled:                         true,
			MaxNumberOfTokensPerMsg:           1,
			MaxDataBytes:                      30_000,
			MaxPerMsgGasLimit:                 3_000_000,
			DestGasOverhead:                   300_000,
			DestGasPerPayloadByteBase:         byte(16),
			DestGasPerPayloadByteHigh:         byte(40),
			DestGasPerPayloadByteThreshold:    uint16(3000),
			DestDataAvailabilityOverheadGas:   100,
			DestGasPerDataAvailabilityByte:    16,
			DestDataAvailabilityMultiplierBps: 1,
			ChainFamilySelector:               []byte{40, 18, 213, 44},
			EnforceOutOfOrder:                 false,
			DefaultTokenFeeUsdCents:           25,
			DefaultTokenDestGasOverhead:       90_000,
			DefaultTxGasLimit:                 200_000,
			GasMultiplierWeiPerEth:            1_000_000_000_000_000_000,
			GasPriceStalenessThreshold:        1_000_000,
			NetworkFeeUsdCents:                10,

			// apply_premium_multiplier_wei_per_eth_updates
			PremiumMultiplierWeiPerEth: []uint64{900_000_000_000_000_000},
		}

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
				DestChainSelectors:        []uint64{config.ContractParamsPerChain[chainSel].DestChainSelector},
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
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy CCIP OnRamp for Sui chain %d: %w", chainSel, err)
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

		// TODO: This should be retrieved from address book or provided in config
		// For now using placeholder - this should be the ethereum chain onRamp bytes
		onRampBytes := [][]byte{{0x0}} // placeholder for ethereum chain onRamp bytes

		// Run DeployAndInitCCIPOffRampSequence
		ccipOffRampSeqInput := offrampops.DeployAndInitCCIPOffRampSeqInput{
			DeployCCIPOffRampInput: offrampops.DeployCCIPOffRampInput{
				CCIPPackageId: ccipSeqReport.Output.CCIPPackageId,
				MCMSPackageId: mcmsSeqReport.Output.PackageId,
			},
			InitializeOffRampInput: offrampops.InitializeOffRampInput{
				DestTransferCapId:                     ccipSeqReport.Output.Objects.DestTransferCapObjectId,
				FeeQuoterCapId:                        ccipSeqReport.Output.Objects.FeeQuoterCapObjectId,
				ChainSelector:                         suiChain.Selector,
				PremissionExecThresholdSeconds:        uint32(60 * 60 * 8),
				SourceChainSelectors:                  []uint64{config.ContractParamsPerChain[chainSel].DestChainSelector}, // this is ethereum
				SourceChainsIsEnabled:                 []bool{true},
				SourceChainsIsRMNVerificationDisabled: []bool{true},
				SourceChainsOnRamp:                    onRampBytes,
			},
		}
		ccipOffRampSeqReport, err := operations.ExecuteSequence(e.OperationsBundle, offrampops.DeployAndInitCCIPOffRampSequence, deps, ccipOffRampSeqInput)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy CCIP OffRamp for Sui chain %d: %w", chainSel, err)
		}
		seqReports = append(seqReports, ccipOffRampSeqReport.ExecutionReports...)

		fmt.Println("SUI OFFRAMP: ", ccipOffRampSeqReport.Output.CCIPOffRampPackageId)

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
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to save TokenPool address %s for Sui chain %d: %w", deployTp.Output.PackageId, chainSel, err)
		}
	}

	fmt.Println("RAN CS_DEPLOY_SUI_CHAIN")
	return cldf.ChangesetOutput{
		AddressBook: ab,
		Reports:     seqReports,
	}, nil
}

// VerifyPreconditions implements deployment.ChangeSetV2.
func (d DeploySuiChain) VerifyPreconditions(e cldf.Environment, config DeploySuiChainConfig) error {
	return nil
}
