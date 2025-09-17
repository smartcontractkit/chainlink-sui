package deployment

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
	SuiOnRampStateObjectIDType     deployment.ContractType = "SuiOnRampStateObjectID"
	SuiOffRampType                 deployment.ContractType = "SuiOffRamp"
	SuiOffRampOwnerCapObjectIDType deployment.ContractType = "SuiOffRampOwnerCapObjectId"
	SuiOffRampStateObjectIDType    deployment.ContractType = "SuiOffRampStateObjectId"
	SuiLockReleaseTPType           deployment.ContractType = "SuiLockReleaseToken"
	SuiLockReleaseTPStateType      deployment.ContractType = "SuiLockReleaseTokenState"

	// MCMS Related
	SuiMcmsPackageIDType               deployment.ContractType = "SuiManyChainMultisigPackageId"
	SuiMcmsObjectIDType                deployment.ContractType = "SuiManyChainMultisigObjectId"
	SuiMcmsRegistryObjectIdType        deployment.ContractType = "SuiManyChainMultisigRegistryObjectId"
	SuiMcmsAccountStateObjectIdType    deployment.ContractType = "SuiManyChainMultisigAccountStateObjectId"
	SuiMcmsAccountOwnerCapObjectIdType deployment.ContractType = "SuiManyChainMultisigAccountOwnerCapObjectId"
	SuiMcmsTimelockObjectIdType        deployment.ContractType = "SuiManyChainMultisigTimelockObjectId"

	// Link related
	SuiLinkTokenType             deployment.ContractType = "SuiLinkToken"
	SuiBnMTokenPoolType          deployment.ContractType = "SuiBnMTokenPool"
	SuiBnMTokenPoolStateType     deployment.ContractType = "SuiBnMTokenPoolState"
	SuiLinkTokenObjectMetadataID deployment.ContractType = "SuiLinkTokenObjectMetadataId"
	SuiLinkTokenTreasuryCapID    deployment.ContractType = "SuiLinkTokenTreasuryCapId"
)

var (
	Version1_0_0 = *semver.MustParse("1.0.0")
)
