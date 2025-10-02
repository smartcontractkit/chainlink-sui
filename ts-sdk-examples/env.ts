import dotenv from 'dotenv'

dotenv.config()

export interface EnvConfig {
    SUI_PRIVATE_KEY: string;
    CCIP_PACKAGE_ID: string;
    CCIP_STATE_ID: string;
    CCIP_OWNER_CAP_ID: string;
    ONRAMP_PACKAGE_ID: string;
    ONRAMP_STATE_ID: string;
    LR_POOL_PACKAGE_ID: string;
    LR_POOL_STATE_ID: string;
    BM_POOL_PACKAGE_ID: string;
    BM_POOL_STATE_ID: string;
    LINK_COIN_TYPE: string;
    ETH_COIN_TYPE: string;
    ETH_TREASURY_CAP_ID: string;
    LINK_TREASURY_CAP_ID: string;
    ETH_METADATA: string;
    LINK_METADATA: string;
    LINK_COIN_ID: string;
    ETH_COIN_ID: string;
}

export const env: EnvConfig = {
    SUI_PRIVATE_KEY: process.env.SUI_PRIVATE_KEY || '',
    CCIP_PACKAGE_ID: process.env.CCIP_PACKAGE_ID || '',
    CCIP_STATE_ID: process.env.CCIP_STATE_ID || '',
    CCIP_OWNER_CAP_ID: process.env.CCIP_OWNER_CAP_ID || '',
    ONRAMP_PACKAGE_ID: process.env.ONRAMP_PACKAGE_ID || '',
    ONRAMP_STATE_ID: process.env.ONRAMP_STATE_ID || '',
    LR_POOL_PACKAGE_ID: process.env.LR_POOL_PACKAGE_ID || '',
    LR_POOL_STATE_ID: process.env.LR_POOL_STATE_ID || '',
    BM_POOL_PACKAGE_ID: process.env.BM_POOL_PACKAGE_ID || '',
    BM_POOL_STATE_ID: process.env.BM_POOL_STATE_ID || '',
    LINK_COIN_TYPE: process.env.LINK_COIN_TYPE || '',
    ETH_COIN_TYPE: process.env.ETH_COIN_TYPE || '',
    ETH_TREASURY_CAP_ID: process.env.ETH_TREASURY_CAP_ID || '',
    LINK_TREASURY_CAP_ID: process.env.LINK_TREASURY_CAP_ID || '',
    ETH_METADATA: process.env.ETH_METADATA || '',
    LINK_METADATA: process.env.LINK_METADATA || '',
    LINK_COIN_ID: process.env.LINK_COIN_OBJECT || '',
    ETH_COIN_ID: process.env.ETH_COIN_OBJECT || '',
}
