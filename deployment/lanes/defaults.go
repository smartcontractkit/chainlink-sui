package lanes

import suideploy "github.com/smartcontractkit/chainlink-sui/deployment"

// Default token transfer fee fields for FeeQuoterApplyTokenTransferFeeConfigUpdatesOp on the
// Sui source leg. Values match deployment.DefaultCCIPSeqConfig / connect_sui_to_evm YAML.
var (
	DefaultLinkTokenTransferMinFeeUsdCents    = suideploy.DefaultCCIPSeqConfig.AddMinFeeUsdCents[0]
	DefaultLinkTokenTransferMaxFeeUsdCents    = suideploy.DefaultCCIPSeqConfig.AddMaxFeeUsdCents[0]
	DefaultLinkTokenTransferDeciBps           = suideploy.DefaultCCIPSeqConfig.AddDeciBps[0]
	DefaultLinkTokenTransferDestGasOverhead   = suideploy.DefaultCCIPSeqConfig.AddDestGasOverhead[0]
	DefaultLinkTokenTransferDestBytesOverhead = suideploy.DefaultCCIPSeqConfig.AddDestBytesOverhead[0]
	DefaultLinkPremiumMultiplierWeiPerEth     = suideploy.DefaultCCIPSeqConfig.PremiumMultiplierWeiPerEth[0]
)
