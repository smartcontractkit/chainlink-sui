package deployment

import (
	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

var (
	SuiCCIPRouterType              deployment.ContractType = "SuiRouter"
	SuiCCIPType                    deployment.ContractType = "SuiCCIP"
	SuiCCIPObjectRefType           deployment.ContractType = "SuiCCIPObjectRef"
	SuiCCIPOwnerCapObjectIdType    deployment.ContractType = "SuiCCIPOwnerCapObjectId"
	SuiFeeQuoterCapType            deployment.ContractType = "SuiCCIPFeeQuoterCap"
	SuiTokenPoolType               deployment.ContractType = "SuiTokenPool"
	SuiOnRampType                  deployment.ContractType = "SuiOnRamp"
	SuiOnRampStateObjectIdType     deployment.ContractType = "SuiOnRampStateObjectId"
	SuiOnRampOwnerCapObjectIdType  deployment.ContractType = "SuiOnRampOwnerCapObjectId"
	SuiOffRampType                 deployment.ContractType = "SuiOffRamp"
	SuiOffRampOwnerCapObjectIdType deployment.ContractType = "SuiOffRampOwnerCapObjectId"
	SuiOffRampStateObjectIdType    deployment.ContractType = "SuiOffRampStateObjectId"
	SuiLockReleaseTPType           deployment.ContractType = "SuiLockReleaseToken"
	SuiLockReleaseTPStateType      deployment.ContractType = "SuiLockReleaseTokenState"

	// MCMS Related
	SuiMcmsPackageIDType               deployment.ContractType = "SuiManyChainMultisigPackageId"
	SuiMcmsObjectIDType                deployment.ContractType = "SuiManyChainMultisigObjectId"
	SuiMcmsRegistryObjectIdType        deployment.ContractType = "SuiManyChainMultisigRegistryObjectId"
	SuiMcmsAccountStateObjectIdType    deployment.ContractType = "SuiManyChainMultisigAccountStateObjectId"
	SuiMcmsAccountOwnerCapObjectIdType deployment.ContractType = "SuiManyChainMultisigAccountOwnerCapObjectId"
	SuiMcmsTimelockObjectIdType        deployment.ContractType = "SuiManyChainMultisigTimelockObjectId"

	// MCMS User Related
	SuiMcmsUserPackageIDType        deployment.ContractType = "SuiMcmsUserPackageId"
	SuiMcmsUserDataObjectIDType     deployment.ContractType = "SuiMcmsUserDataObjectId"
	SuiMcmsUserOwnerCapObjectIDType deployment.ContractType = "SuiMcmsUserOwnerCapObjectId"

	// Link related
	SuiLinkTokenObjectMetadataID deployment.ContractType = "SuiLinkTokenObjectMetadataId"
	SuiLinkTokenTreasuryCapID    deployment.ContractType = "SuiLinkTokenTreasuryCapId"
	SuiMCMSType                  deployment.ContractType = "SuiManyChainMultisig"
	SuiLinkTokenType             deployment.ContractType = "SuiLinkToken"
	SuiBnMTokenPoolType          deployment.ContractType = "SuiBnMTokenPool"
	SuiBnMTokenPoolStateType     deployment.ContractType = "SuiBnMTokenPoolState"
	SuiLinkTokenObjectMetadataId deployment.ContractType = "SuiLinkTokenObjectMetadataId"
	SuiLinkTokenTreasuryCapId    deployment.ContractType = "SuiLinkTokenTreasuryCapId"
)

var (
	Version1_0_0 = *semver.MustParse("1.0.0")
)
