package ccipops

import (
	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
)

type DeployCCIPSeqObjects struct {
	CCIPObjectRefObjectId           string
	OwnerCapObjectId                string
	FeeQuoterCapObjectId            string
	FeeQuoterStateObjectId          string
	NonceManagerStateObjectId       string
	NonceManagerCapObjectId         string
	ReceiverRegistryStateObjectId   string
	RMNRemoteStateObjectId          string
	TokenAdminRegistryStateObjectId string
	SourceTransferCapObjectId       string
	DestTransferCapObjectId         string
}

type DeployCCIPSeqOutput struct {
	CCIPPackageId string
	Objects       DeployCCIPSeqObjects
}

type DeployAndInitCCIPSeqInput struct {
	LinkTokenCoinMetadataObjectId string
	LocalChainSelector            uint64
	DestChainSelector             uint64
	DeployCCIPInput
	// Fee Quoter
	MaxFeeJuelsPerMsg            string
	TokenPriceStalenessThreshold uint64

	// Fee Quoter configuration
	AddMinFeeUsdCents    []uint32
	AddMaxFeeUsdCents    []uint32
	AddDeciBps           []uint16
	AddDestGasOverhead   []uint32
	AddDestBytesOverhead []uint32
	AddIsEnabled         []bool
	RemoveTokens         []string
	// Fee Quoter destination chain configuration
	IsEnabled                         bool
	MaxNumberOfTokensPerMsg           uint16
	MaxDataBytes                      uint32
	MaxPerMsgGasLimit                 uint32
	DestGasOverhead                   uint32
	DestGasPerPayloadByteBase         byte
	DestGasPerPayloadByteHigh         byte
	DestGasPerPayloadByteThreshold    uint16
	DestDataAvailabilityOverheadGas   uint32
	DestGasPerDataAvailabilityByte    uint16
	DestDataAvailabilityMultiplierBps uint16
	ChainFamilySelector               []byte
	EnforceOutOfOrder                 bool
	DefaultTokenFeeUsdCents           uint16
	DefaultTokenDestGasOverhead       uint32
	DefaultTxGasLimit                 uint32
	GasMultiplierWeiPerEth            uint64
	GasPriceStalenessThreshold        uint32
	NetworkFeeUsdCents                uint32
	// Premium multiplier updates
	PremiumMultiplierWeiPerEth []uint64

	// RMN Remote config
	RmnHomeContractConfigDigest []byte
	SignerOnchainPublicKeys     [][]byte
	NodeIndexes                 []uint64
	FSign                       uint64
}

var DeployAndInitCCIPSequence = cld_ops.NewSequence(
	"sui-deploy-ccip-seq",
	semver.MustParse("0.1.0"),
	"Deploys and sets initial CCIP configuration",
	func(env cld_ops.Bundle, deps sui_ops.OpTxDeps, input DeployAndInitCCIPSeqInput) (DeployCCIPSeqOutput, error) {
		deployReport, err := cld_ops.ExecuteOperation(env, DeployCCIPOp, deps, input.DeployCCIPInput)
		if err != nil {
			return DeployCCIPSeqOutput{}, err
		}

		initFQReport, err := cld_ops.ExecuteOperation(
			env,
			FeeQuoterInitializeOp,
			deps,
			InitFeeQuoterInput{
				CCIPPackageId:                 deployReport.Output.PackageId,
				StateObjectId:                 deployReport.Output.Objects.CCIPObjectRefObjectId,
				OwnerCapObjectId:              deployReport.Output.Objects.OwnerCapObjectId,
				MaxFeeJuelsPerMsg:             input.MaxFeeJuelsPerMsg,
				LinkTokenCoinMetadataObjectId: input.LinkTokenCoinMetadataObjectId,
				TokenPriceStalenessThreshold:  input.TokenPriceStalenessThreshold,
				FeeTokens:                     []string{input.LinkTokenCoinMetadataObjectId},
			},
		)
		if err != nil {
			return DeployCCIPSeqOutput{}, err
		}

		issueFQCapReport, err := cld_ops.ExecuteOperation(
			env,
			FeeQuoterIssueFeeQuoterCapOp,
			deps,
			IssueFeeQuoterCapInput{
				CCIPPackageId:    deployReport.Output.PackageId,
				OwnerCapObjectId: deployReport.Output.Objects.OwnerCapObjectId,
			},
		)
		if err != nil {
			return DeployCCIPSeqOutput{}, err
		}

		initNMReport, err := cld_ops.ExecuteOperation(
			env,
			NonceManagerInitializeOp,
			deps,
			InitNMInput{
				CCIPPackageId:    deployReport.Output.PackageId,
				StateObjectId:    deployReport.Output.Objects.CCIPObjectRefObjectId,
				OwnerCapObjectId: deployReport.Output.Objects.OwnerCapObjectId,
			},
		)
		if err != nil {
			return DeployCCIPSeqOutput{}, err
		}

		initRecRegReport, err := cld_ops.ExecuteOperation(
			env,
			ReceiverRegistryInitializeOp,
			deps,
			InitRecRegInput{
				CCIPPackageId:    deployReport.Output.PackageId,
				StateObjectId:    deployReport.Output.Objects.CCIPObjectRefObjectId,
				OwnerCapObjectId: deployReport.Output.Objects.OwnerCapObjectId,
			},
		)
		if err != nil {
			return DeployCCIPSeqOutput{}, err
		}

		initRMNRemoteReport, err := cld_ops.ExecuteOperation(
			env,
			RMNRemoteInitializeOp,
			deps,
			InitRMNRemoteInput{
				CCIPPackageId:      deployReport.Output.PackageId,
				StateObjectId:      deployReport.Output.Objects.CCIPObjectRefObjectId,
				OwnerCapObjectId:   deployReport.Output.Objects.OwnerCapObjectId,
				LocalChainSelector: input.LocalChainSelector,
			},
		)
		if err != nil {
			return DeployCCIPSeqOutput{}, err
		}

		initTARReport, err := cld_ops.ExecuteOperation(
			env,
			TokenAdminRegistryInitializeOp,
			deps,
			InitTARInput{
				CCIPPackageId:      deployReport.Output.PackageId,
				StateObjectId:      deployReport.Output.Objects.CCIPObjectRefObjectId,
				OwnerCapObjectId:   deployReport.Output.Objects.OwnerCapObjectId,
				LocalChainSelector: input.LocalChainSelector,
			},
		)
		if err != nil {
			return DeployCCIPSeqOutput{}, err
		}

		// apply_token_transfer_fee_config_updates
		_, err = cld_ops.ExecuteOperation(
			env,
			FeeQuoterApplyTokenTransferFeeConfigUpdatesOp,
			deps,
			FeeQuoterApplyTokenTransferFeeConfigUpdatesInput{
				CCIPPackageId:        deployReport.Output.PackageId,
				StateObjectId:        deployReport.Output.Objects.CCIPObjectRefObjectId,
				OwnerCapObjectId:     deployReport.Output.Objects.OwnerCapObjectId,
				AddTokens:            []string{input.LinkTokenCoinMetadataObjectId},
				AddMinFeeUsdCents:    input.AddMinFeeUsdCents,
				AddMaxFeeUsdCents:    input.AddMaxFeeUsdCents,
				AddDeciBps:           input.AddDeciBps,
				AddDestGasOverhead:   input.AddDestGasOverhead,
				AddDestBytesOverhead: input.AddDestBytesOverhead,
				AddIsEnabled:         input.AddIsEnabled,
				RemoveTokens:         input.RemoveTokens,
			},
		)
		if err != nil {
			return DeployCCIPSeqOutput{}, err
		}

		// apply_dest_chain_config_updates
		_, err = cld_ops.ExecuteOperation(
			env,
			FeeQuoterApplyDestChainConfigUpdatesOp,
			deps,
			FeeQuoterApplyDestChainConfigUpdatesInput{
				CCIPPackageId:                     deployReport.Output.PackageId,
				StateObjectId:                     deployReport.Output.Objects.CCIPObjectRefObjectId,
				OwnerCapObjectId:                  deployReport.Output.Objects.OwnerCapObjectId,
				DestChainSelector:                 input.DestChainSelector,
				IsEnabled:                         input.IsEnabled,
				MaxNumberOfTokensPerMsg:           input.MaxNumberOfTokensPerMsg,
				MaxDataBytes:                      input.MaxDataBytes,
				MaxPerMsgGasLimit:                 input.MaxPerMsgGasLimit,
				DestGasOverhead:                   input.DestGasOverhead,
				DestGasPerPayloadByteBase:         input.DestGasPerPayloadByteBase,
				DestGasPerPayloadByteHigh:         input.DestGasPerPayloadByteHigh,
				DestGasPerPayloadByteThreshold:    input.DestGasPerPayloadByteThreshold,
				DestDataAvailabilityOverheadGas:   input.DestDataAvailabilityOverheadGas,
				DestGasPerDataAvailabilityByte:    input.DestGasPerDataAvailabilityByte,
				DestDataAvailabilityMultiplierBps: input.DestDataAvailabilityMultiplierBps,
				ChainFamilySelector:               input.ChainFamilySelector,
				EnforceOutOfOrder:                 input.EnforceOutOfOrder,
				DefaultTokenFeeUsdCents:           input.DefaultTokenFeeUsdCents,
				DefaultTokenDestGasOverhead:       input.DefaultTokenDestGasOverhead,
				DefaultTxGasLimit:                 input.DefaultTxGasLimit,
				GasMultiplierWeiPerEth:            input.GasMultiplierWeiPerEth,
				GasPriceStalenessThreshold:        input.GasPriceStalenessThreshold,
				NetworkFeeUsdCents:                input.NetworkFeeUsdCents,
			},
		)
		if err != nil {
			return DeployCCIPSeqOutput{}, err
		}

		// apply_premium_multiplier_wei_per_eth_updates
		_, err = cld_ops.ExecuteOperation(
			env,
			FeeQuoterApplyPremiumMultiplierWeiPerEthUpdatesOp,
			deps,
			FeeQuoterApplyPremiumMultiplierWeiPerEthUpdatesInput{
				CCIPPackageId:              deployReport.Output.PackageId,
				StateObjectId:              deployReport.Output.Objects.CCIPObjectRefObjectId,
				OwnerCapObjectId:           deployReport.Output.Objects.OwnerCapObjectId,
				Tokens:                     []string{input.LinkTokenCoinMetadataObjectId},
				PremiumMultiplierWeiPerEth: input.PremiumMultiplierWeiPerEth,
			},
		)
		if err != nil {
			return DeployCCIPSeqOutput{}, err
		}

		// TODO: RMN Set config. For e2e could be disabled
		_, err = cld_ops.ExecuteOperation(
			env,
			RMNRemoteSetConfigOp,
			deps,
			RMNRemoteSetConfigInput{
				CCIPPackageId:               deployReport.Output.PackageId,
				StateObjectId:               deployReport.Output.Objects.CCIPObjectRefObjectId,
				OwnerCapObjectId:            deployReport.Output.Objects.OwnerCapObjectId,
				RmnHomeContractConfigDigest: input.RmnHomeContractConfigDigest,
				SignerOnchainPublicKeys:     input.SignerOnchainPublicKeys,
				NodeIndexes:                 input.NodeIndexes,
				FSign:                       input.FSign,
			},
		)
		if err != nil {
			return DeployCCIPSeqOutput{}, err
		}

		return DeployCCIPSeqOutput{
			CCIPPackageId: deployReport.Output.PackageId,
			Objects: DeployCCIPSeqObjects{
				CCIPObjectRefObjectId:           deployReport.Output.Objects.CCIPObjectRefObjectId,
				OwnerCapObjectId:                deployReport.Output.Objects.OwnerCapObjectId,
				FeeQuoterCapObjectId:            issueFQCapReport.Output.Objects.FeeQuoterCapObjectId,
				FeeQuoterStateObjectId:          initFQReport.Output.Objects.FeeQuoterStateObjectId,
				NonceManagerStateObjectId:       initNMReport.Output.Objects.NonceManagerStateObjectId,
				NonceManagerCapObjectId:         initNMReport.Output.Objects.NonceManagerCapObjectId,
				ReceiverRegistryStateObjectId:   initRecRegReport.Output.Objects.ReceiverRegistryStateObjectId,
				RMNRemoteStateObjectId:          initRMNRemoteReport.Output.Objects.RMNRemoteStateObjectId,
				TokenAdminRegistryStateObjectId: initTARReport.Output.Objects.TARStateObjectId,
				SourceTransferCapObjectId:       deployReport.Output.Objects.SourceTransferCapObjectId,
				DestTransferCapObjectId:         deployReport.Output.Objects.DestTransferCapObjectId,
			},
		}, nil
	},
)
