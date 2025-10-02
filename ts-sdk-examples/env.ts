import dotenv from 'dotenv'

dotenv.config()

export interface EnvConfig {
    SUI_PRIVATE_KEY: string;
    SuiOnRampStateObjectID: string;
    SuiLinkToken: string;
    SuiManyChainMultisigPackageId: string;
    SuiOffRampStateObjectId: string;
    SuiLinkTokenTreasuryCapId: string;
    SuiOffRamp: string;
    SuiTokenPool: string;
    SuiCCIP: string;
    SuiCCIPFeeQuoterCap: string;
    SuiRouter: string;
    SuiOffRampOwnerCapObjectId: string;
    SuiLinkTokenObjectMetadataId: string;
    SuiOnRamp: string;
    SuiCCIPObjectRef: string;
}

export const env: EnvConfig = {
    SUI_PRIVATE_KEY: process.env.SUI_PRIVATE_KEY || '',
    SuiOnRampStateObjectID: process.env.SuiOnRampStateObjectID || '',
    SuiLinkToken: process.env.SuiLinkToken || '',
    SuiManyChainMultisigPackageId: process.env.SuiManyChainMultisigPackageId || '',
    SuiOffRampStateObjectId: process.env.SuiOffRampStateObjectId || '',
    SuiLinkTokenTreasuryCapId: process.env.SuiLinkTokenTreasuryCapId || '',
    SuiOffRamp: process.env.SuiOffRamp || '',
    SuiTokenPool: process.env.SuiTokenPool || '',
    SuiCCIP: process.env.SuiCCIP || '',
    SuiCCIPFeeQuoterCap: process.env.SuiCCIPFeeQuoterCap || '',
    SuiRouter: process.env.SuiRouter || '',
    SuiOffRampOwnerCapObjectId: process.env.SuiOffRampOwnerCapObjectId || '',
    SuiLinkTokenObjectMetadataId: process.env.SuiLinkTokenObjectMetadataId || '',
    SuiOnRamp: process.env.SuiOnRamp || '',
    SuiCCIPObjectRef: process.env.SuiCCIPObjectRef || '',
}
