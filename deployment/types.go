package deployment

import (
	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

var (
	// CCIP
	SuiCCIPType                      deployment.ContractType = "SuiCCIP"
	SuiLatestCCIPPackageIDType       deployment.ContractType = "SuiLatestCCIPPackageID"
	SuiCCIPObjectRefType             deployment.ContractType = "SuiCCIPObjectRef"
	SuiCCIPOwnerCapObjectIDType      deployment.ContractType = "SuiCCIPOwnerCapObjectID"
	SuiCurserCapObjectIDType         deployment.ContractType = "SuiCurserCapObjectID"
	SuiCCIPUpgradeCapObjectIDType    deployment.ContractType = "SuiCCIPUpgradeCapObjectID"
	SuiFeeQuoterCapType              deployment.ContractType = "SuiCCIPFeeQuoterCap"
	SuiOnRampType                    deployment.ContractType = "SuiOnRamp"
	SuiLatestOnRampPackageIDType     deployment.ContractType = "SuiLatestOnRampPackageID"
	SuiOnRampStateObjectIDType       deployment.ContractType = "SuiOnRampStateObjectID"
	SuiOnRampOwnerCapObjectIDType    deployment.ContractType = "SuiOnRampOwnerCapObjectID"
	SuiOnRampUpgradeCapObjectIDType  deployment.ContractType = "SuiOnRampUpgradeCapObjectID"
	SuiOffRampType                   deployment.ContractType = "SuiOffRamp"
	SuiLatestOffRampPackageIDType    deployment.ContractType = "SuiLatestOffRampPackageID"
	SuiOffRampOwnerCapObjectIDType   deployment.ContractType = "SuiOffRampOwnerCapObjectID"
	SuiOffRampUpgradeCapObjectIDType deployment.ContractType = "SuiOffRampUpgradeCapObjectID"
	SuiOffRampStateObjectIDType      deployment.ContractType = "SuiOffRampStateObjectID"
	SuiLockReleaseTPType             deployment.ContractType = "SuiLockReleaseToken"
	SuiLockReleaseTPStateType        deployment.ContractType = "SuiLockReleaseTokenState"

	// CCIP Router
	SuiCCIPRouterType                 deployment.ContractType = "SuiRouter"
	SuiCCIPRouterStateObjectType      deployment.ContractType = "SuiRouterStateObjectID"
	SuiCCIPRouterOwnerCapObjectIDType deployment.ContractType = "SuiCCIPRouterOwnerCapObjectID"
	SuiRouterUpgradeCapObjectIDType   deployment.ContractType = "SuiRouterUpgradeCapObjectID"

	// MCMS Related
	SuiMCMSType                        deployment.ContractType = "SuiManyChainMultisig"
	SuiMcmsPackageIDType               deployment.ContractType = "SuiManyChainMultisigPackageID"
	SuiMcmsObjectIDType                deployment.ContractType = "SuiManyChainMultisigObjectID"
	SuiMcmsRegistryObjectIDType        deployment.ContractType = "SuiManyChainMultisigRegistryObjectID"
	SuiMcmsAccountStateObjectIDType    deployment.ContractType = "SuiManyChainMultisigAccountStateObjectID"
	SuiMcmsAccountOwnerCapObjectIDType deployment.ContractType = "SuiManyChainMultisigAccountOwnerCapObjectID"
	SuiMcmsTimelockObjectIDType        deployment.ContractType = "SuiManyChainMultisigTimelockObjectID"
	SuiMcmsDeployerObjectIDType        deployment.ContractType = "SuiManyChainMultisigDeployerObjectID"

	// MCMS User Related
	SuiMcmsUserPackageIDType        deployment.ContractType = "SuiMcmsUserPackageID"
	SuiMcmsUserDataObjectIDType     deployment.ContractType = "SuiMcmsUserDataObjectID"
	SuiMcmsUserOwnerCapObjectIDType deployment.ContractType = "SuiMcmsUserOwnerCapObjectID"

	// Link related
	SuiLinkTokenObjectMetadataID deployment.ContractType = "SuiLinkTokenObjectMetadataID"
	SuiLinkTokenTreasuryCapID    deployment.ContractType = "SuiLinkTokenTreasuryCapID"
	SuiLinkTokenUpgradeCapID     deployment.ContractType = "SuiLinkTokenUpgradeCapID"
	SuiLinkTokenType             deployment.ContractType = "SuiLinkToken"

	// Managed Token related
	// the coins under management
	SuiManagedTokenType               deployment.ContractType = "SuiManagedToken"
	SuiManagedTokenCoinMetadataIDType deployment.ContractType = "SuiManagedTokenCoinMetadataID"
	SuiManagedTokenTreasuryCapIDType  deployment.ContractType = "SuiManagedTokenTreasuryCapID"
	SuiManagedTokenUpgradeCapIDType   deployment.ContractType = "SuiManagedTokenUpgradeCapID"
	// the managed token wrapper package for the tokens
	SuiManagedTokenPackageIDType     deployment.ContractType = "SuiManagedTokenPackageID"
	SuiManagedTokenOwnerCapObjectID  deployment.ContractType = "SuiManagedTokenOwnerCapObjectID"
	SuiManagedTokenStateObjectID     deployment.ContractType = "SuiManagedTokenStateObjectID"
	SuiManagedTokenMinterCapID       deployment.ContractType = "SuiManagedTokenMinterCapID"
	SuiManagedTokenPublisherObjectId deployment.ContractType = "SuiManagedTokenPublisherObjectId"
	// Managed token faucet package
	SuiManagedTokenFaucetPackageIDType          deployment.ContractType = "SuiManagedTokenFaucetPackageID"
	SuiManagedTokenFaucetStateObjectIDType      deployment.ContractType = "SuiManagedTokenFaucetStateObjectID"
	SuiManagedTokenFaucetUpgradeCapObjectIDType deployment.ContractType = "SuiManagedTokenFaucetUpgradeCapObjectID"

	// BnM Token Pool related
	SuiBnMTokenPoolType                   deployment.ContractType = "SuiBnMTokenPool"
	SuiBnMTokenPoolStateType              deployment.ContractType = "SuiBnMTokenPoolState"
	SuiBnMTokenPoolOwnerCapObjectIDType   deployment.ContractType = "SuiBnMTokenPoolOwnerCapObjectID"
	SuiBnMTokenPoolUpgradeCapObjectIDType deployment.ContractType = "SuiBnMTokenPoolUpgradeCapObjectID"

	// LnR Token Pool related
	SuiLnRTokenPoolType                   deployment.ContractType = "SuiLnRTokenPool"
	SuiLnRTokenPoolStateType              deployment.ContractType = "SuiLnRTokenPoolState"
	SuiLnRTokenPoolOwnerCapObjectIDType   deployment.ContractType = "SuiLnRTokenPoolOwnerCapObjectID"
	SuiLnRTokenPoolUpgradeCapObjectIDType deployment.ContractType = "SuiLnRTokenPoolUpgradeCapObjectID"
	SuiLnRTokenPoolRebalancerCapIDType    deployment.ContractType = "SuiLnRTokenPoolRebalancerCapID"
	// the lock-release coin itself, recorded under its own type rather than the managed-token bucket
	SuiLnRTokenType               deployment.ContractType = "SuiLnRToken"
	SuiLnRTokenCoinMetadataIDType deployment.ContractType = "SuiLnRTokenCoinMetadataID"
	SuiLnRTokenTreasuryCapIDType  deployment.ContractType = "SuiLnRTokenTreasuryCapID"
	SuiLnRTokenUpgradeCapIDType   deployment.ContractType = "SuiLnRTokenUpgradeCapID"

	// Managed Token Pool related
	SuiManagedTokenPoolType                   deployment.ContractType = "SuiManagedTokenPool"
	SuiManagedTokenPoolStateType              deployment.ContractType = "SuiManagedTokenPoolState"
	SuiManagedTokenPoolOwnerCapObjectIDType   deployment.ContractType = "SuiManagedTokenPoolOwnerCapObjectID"
	SuiManagedTokenPoolUpgradeCapObjectIDType deployment.ContractType = "SuiManagedTokenPoolUpgradeCapObjectID"

	// Upgrade Related
	SuiCCIPMockV2              deployment.ContractType = "SuiCCIPMockV2PackageID"
	SuiOnRampMockV2            deployment.ContractType = "SuiOnRampMockV2PackageID"
	SuiOffRampMockV2           deployment.ContractType = "SuiOffRampMockV2PackageID"
	SuiUpgradeRegistryObjectId deployment.ContractType = "SuiUpgradeRegistryObjectId"
)

var (
	Version1_0_0 = *semver.MustParse("1.0.0")
)

// Coin type suffixes (module::STRUCT) for the BnM and LINK coins. The coin deploy changesets
// attach these as a coinType= label on the coin package ref so deriveSuiCoinType can build the
// correct coin type for any token, including ones whose suffix differs, from that label.
var (
	SuiCCIPBnMCoinTypeSuffix = "ccip_burn_mint_token::CCIP_BURN_MINT_TOKEN"
	SuiCCIPLnRCoinTypeSuffix = "ccip_lock_release_token::CCIP_LOCK_RELEASE_TOKEN"
	SuiLinkCoinTypeSuffix    = "link::LINK"
)
