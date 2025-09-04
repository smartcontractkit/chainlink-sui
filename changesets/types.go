package changesets

import (
	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

var (
	SuiCCIPRouterType              deployment.ContractType = "SuiRouter"
	SuiCCIPType                    deployment.ContractType = "SuiCCIP"
	SuiCCIPObjectRefType           deployment.ContractType = "SuiCCIPObjectRef"
	SuiFeeQuoterCapType            deployment.ContractType = "SuiCCIPFeeQuoterCap"
	SuiTokenPoolType               deployment.ContractType = "SuiTokenPool"
	SuiOnRampType                  deployment.ContractType = "SuiOnRamp"
	SuiOnRampStateObjectIdType     deployment.ContractType = "SuiOnRampStateObjectId"
	SuiOffRampType                 deployment.ContractType = "SuiOffRamp"
	SuiOffRampOwnerCapObjectIdType deployment.ContractType = "SuiOffRampOwnerCapObjectId"
	SuiOffRampStateObjectIdType    deployment.ContractType = "SuiOffRampStateObjectId"
	SuiLockReleaseTPType           deployment.ContractType = "SuiLockReleaseToken"
	SuiLockReleaseTPStateType      deployment.ContractType = "SuiLockReleaseTokenState"
	SuiMCMSType                    deployment.ContractType = "SuiManyChainMultisig"
	SuiLinkTokenType               deployment.ContractType = "SuiLinkToken"
	SuiBnMTokenPoolType            deployment.ContractType = "SuiBnMTokenPool"
	SuiBnMTokenPoolStateType       deployment.ContractType = "SuiBnMTokenPoolState"
	SuiLinkTokenObjectMetadataId   deployment.ContractType = "SuiLinkTokenObjectMetadataId"
	SuiLinkTokenTreasuryCapId      deployment.ContractType = "SuiLinkTokenTreasuryCapId"
)

var (
	Version1_0_0 = *semver.MustParse("1.0.0")
)
