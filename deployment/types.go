package deployment

import (
	"github.com/Masterminds/semver/v3"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

var (
	SuiCCIPRouterType                deployment.ContractType = "SuiRouter"
	SuiCCIPType                      deployment.ContractType = "SuiCCIP"
	SuiCCIPObjectRefType             deployment.ContractType = "SuiCCIPObjectRef"
	SuiCCIPOwnerCapObjectIDType      deployment.ContractType = "SuiCCIPOwnerCapObjectID"
	SuiCCIPUpgradeCapObjectIDType    deployment.ContractType = "SuiCCIPUpgradeCapObjectID"
	SuiFeeQuoterCapType              deployment.ContractType = "SuiCCIPFeeQuoterCap"
	SuiOnRampType                    deployment.ContractType = "SuiOnRamp"
	SuiOnRampStateObjectIDType       deployment.ContractType = "SuiOnRampStateObjectID"
	SuiOnRampOwnerCapObjectIDType    deployment.ContractType = "SuiOnRampOwnerCapObjectID"
	SuiOnRampUpgradeCapObjectIDType  deployment.ContractType = "SuiOnRampUpgradeCapObjectID"
	SuiOffRampType                   deployment.ContractType = "SuiOffRamp"
	SuiOffRampOwnerCapObjectIDType   deployment.ContractType = "SuiOffRampOwnerCapObjectID"
	SuiOffRampUpgradeCapObjectIDType deployment.ContractType = "SuiOffRampUpgradeCapObjectID"
	SuiOffRampStateObjectIDType      deployment.ContractType = "SuiOffRampStateObjectID"
	SuiLockReleaseTPType             deployment.ContractType = "SuiLockReleaseToken"
	SuiLockReleaseTPStateType        deployment.ContractType = "SuiLockReleaseTokenState"

	// MCMS Related
	SuiMcmsPackageIDType               deployment.ContractType = "SuiManyChainMultisigPackageID"
	SuiMcmsObjectIDType                deployment.ContractType = "SuiManyChainMultisigObjectID"
	SuiMcmsRegistryObjectIDType        deployment.ContractType = "SuiManyChainMultisigRegistryObjectID"
	SuiMcmsAccountStateObjectIDType    deployment.ContractType = "SuiManyChainMultisigAccountStateObjectID"
	SuiMcmsAccountOwnerCapObjectIDType deployment.ContractType = "SuiManyChainMultisigAccountOwnerCapObjectID"
	SuiMcmsTimelockObjectIDType        deployment.ContractType = "SuiManyChainMultisigTimelockObjectID"

	// MCMS User Related
	SuiMcmsUserPackageIDType        deployment.ContractType = "SuiMcmsUserPackageID"
	SuiMcmsUserDataObjectIDType     deployment.ContractType = "SuiMcmsUserDataObjectID"
	SuiMcmsUserOwnerCapObjectIDType deployment.ContractType = "SuiMcmsUserOwnerCapObjectID"

	// Link related
	SuiLinkTokenObjectMetadataID deployment.ContractType = "SuiLinkTokenObjectMetadataID"
	SuiLinkTokenTreasuryCapID    deployment.ContractType = "SuiLinkTokenTreasuryCapID"
	SuiMCMSType                  deployment.ContractType = "SuiManyChainMultisig"
	SuiLinkTokenType             deployment.ContractType = "SuiLinkToken"
	SuiBnMTokenPoolType          deployment.ContractType = "SuiBnMTokenPool"
	SuiBnMTokenPoolStateType     deployment.ContractType = "SuiBnMTokenPoolState"
	SuiBnMTokenPoolOwnerIDType   deployment.ContractType = "SuiBnMTokenPoolOwnerID"

	// Upgrade Related
	SuiCCIPMockV2              deployment.ContractType = "SuiCCIPMockV2PackageID"
	SuiOnRampMockV2            deployment.ContractType = "SuiOnRampMockV2PackageID"
	SuiOffRampMockV2           deployment.ContractType = "SuiOffRampMockV2PackageID"
	SuiUpgradeRegistryObjectId deployment.ContractType = "SuiUpgradeRegistryObjectId"
)

var (
	Version1_0_0 = *semver.MustParse("1.0.0")
)
