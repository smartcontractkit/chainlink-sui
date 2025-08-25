package ccipops

import (
	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
)

type DeployCCIPSeqObjects struct {
	CCIPObjectRefObjectID           string
	OwnerCapObjectID                string
	FeeQuoterCapObjectID            string
	FeeQuoterStateObjectID          string
	NonceManagerStateObjectID       string
	NonceManagerCapObjectID         string
	ReceiverRegistryStateObjectID   string
	RMNRemoteStateObjectID          string
	TokenAdminRegistryStateObjectID string
	SourceTransferCapObjectID       string
	DestTransferCapObjectID         string
}

type DeployCCIPSeqOutput struct {
	CCIPPackageID string
	Objects       DeployCCIPSeqObjects
}

type DeployAndInitCCIPSeqInput struct {
	LinkTokenCoinMetadataObjectID string
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
				CCIPPackageID:                 deployReport.Output.PackageID,
				StateObjectID:                 deployReport.Output.Objects.CCIPObjectRefObjectID,
				OwnerCapObjectID:              deployReport.Output.Objects.OwnerCapObjectID,
				MaxFeeJuelsPerMsg:             input.MaxFeeJuelsPerMsg,
				LinkTokenCoinMetadataObjectID: input.LinkTokenCoinMetadataObjectID,
				TokenPriceStalenessThreshold:  input.TokenPriceStalenessThreshold,
				FeeTokens:                     []string{input.LinkTokenCoinMetadataObjectID},
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
				CCIPPackageID:    deployReport.Output.PackageID,
				OwnerCapObjectID: deployReport.Output.Objects.OwnerCapObjectID,
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
				CCIPPackageID:    deployReport.Output.PackageID,
				StateObjectID:    deployReport.Output.Objects.CCIPObjectRefObjectID,
				OwnerCapObjectID: deployReport.Output.Objects.OwnerCapObjectID,
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
				CCIPPackageID:    deployReport.Output.PackageID,
				StateObjectID:    deployReport.Output.Objects.CCIPObjectRefObjectID,
				OwnerCapObjectID: deployReport.Output.Objects.OwnerCapObjectID,
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
				CCIPPackageID:      deployReport.Output.PackageID,
				StateObjectID:      deployReport.Output.Objects.CCIPObjectRefObjectID,
				OwnerCapObjectID:   deployReport.Output.Objects.OwnerCapObjectID,
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
				CCIPPackageID:      deployReport.Output.PackageID,
				StateObjectID:      deployReport.Output.Objects.CCIPObjectRefObjectID,
				OwnerCapObjectID:   deployReport.Output.Objects.OwnerCapObjectID,
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
				CCIPPackageID:        deployReport.Output.PackageID,
				StateObjectID:        deployReport.Output.Objects.CCIPObjectRefObjectID,
				OwnerCapObjectID:     deployReport.Output.Objects.OwnerCapObjectID,
				AddTokens:            []string{input.LinkTokenCoinMetadataObjectID},
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
				CCIPPackageID:                     deployReport.Output.PackageID,
				StateObjectID:                     deployReport.Output.Objects.CCIPObjectRefObjectID,
				OwnerCapObjectID:                  deployReport.Output.Objects.OwnerCapObjectID,
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
				CCIPPackageID:              deployReport.Output.PackageID,
				StateObjectID:              deployReport.Output.Objects.CCIPObjectRefObjectID,
				OwnerCapObjectID:           deployReport.Output.Objects.OwnerCapObjectID,
				Tokens:                     []string{input.LinkTokenCoinMetadataObjectID},
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
				CCIPPackageID:               deployReport.Output.PackageID,
				StateObjectID:               deployReport.Output.Objects.CCIPObjectRefObjectID,
				OwnerCapObjectID:            deployReport.Output.Objects.OwnerCapObjectID,
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
			CCIPPackageID: deployReport.Output.PackageID,
			Objects: DeployCCIPSeqObjects{
				CCIPObjectRefObjectID:           deployReport.Output.Objects.CCIPObjectRefObjectID,
				OwnerCapObjectID:                deployReport.Output.Objects.OwnerCapObjectID,
				FeeQuoterCapObjectID:            issueFQCapReport.Output.Objects.FeeQuoterCapObjectID,
				FeeQuoterStateObjectID:          initFQReport.Output.Objects.FeeQuoterStateObjectID,
				NonceManagerStateObjectID:       initNMReport.Output.Objects.NonceManagerStateObjectID,
				NonceManagerCapObjectID:         initNMReport.Output.Objects.NonceManagerCapObjectID,
				ReceiverRegistryStateObjectID:   initRecRegReport.Output.Objects.ReceiverRegistryStateObjectID,
				RMNRemoteStateObjectID:          initRMNRemoteReport.Output.Objects.RMNRemoteStateObjectID,
				TokenAdminRegistryStateObjectID: initTARReport.Output.Objects.TARStateObjectID,
				SourceTransferCapObjectID:       deployReport.Output.Objects.SourceTransferCapObjectID,
				DestTransferCapObjectID:         deployReport.Output.Objects.DestTransferCapObjectID,
			},
		}, nil
	},
)
