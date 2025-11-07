import { Command } from 'commander'
import { SuiClient, getFullnodeUrl } from '@mysten/sui/client'
import { Ed25519Keypair } from '@mysten/sui/keypairs/ed25519'
import { buildCcipSendPTB, type BuildArgs } from './onramp'
import { loadEnvForNetwork } from './env'
import { Transaction } from '@mysten/sui/transactions'

// Simple helpers
function parseVecU8(input: string | undefined): Uint8Array {
  if (!input || input.length === 0) return new Uint8Array()
  const s = input.trim()
  if (s.startsWith('0x')) return Uint8Array.from(Buffer.from(s.slice(2), 'hex'))
  return Uint8Array.from(Buffer.from(s, 'base64'))
}

function parseU64BigInt(input: string): bigint {
  // accept decimal or hex (0x...)
  return input.startsWith('0x') ? BigInt(input) : BigInt(input)
}

const program = new Command()
program
  .name('ccip-cli')
  .description('Minimal CLI to build and submit a CCIP PTB on Sui')
  .version('0.1.0')

program.command('send')
  .alias('s')
  .alias('onramp')
  .description('Build and submit a CCIP send PTB')
  .option('--ccip-pkg <id>', 'CCIP package id (from env if not provided)')
  .option('--ccip-object-ref <id>', 'CCIPObjectRef object id (from env if not provided)')
  .option('--onramp-state <id>', 'OnRampState object id (from env if not provided)')
  .option('--token-metadata <id>', 'CoinMetadata<T> for token being transferred (from env if not provided)')
  .option('--token-coin <id>', 'Owned Coin<T> for token being transferred (from env if not provided)')
  .option('--coin-type <type>', 'Token type being transferred (from env if not provided)')
  .requiredOption('--dest-chain-selector <u64>', 'Destination chain selector (decimal or 0x-hex)')
  .requiredOption('--receiver <address>', 'Receiver address')
  .option('--data <bytes>', 'Arbitrary data (0x-hex or base64)', '')
  .requiredOption('--pool-kind <kind>', 'Pool kind', (v) => {
    if (v !== 'burn_mint' && v !== 'lock_release') throw new Error('pool-kind must be burn_mint or lock_release')
    return v
  })
  .option('--fee-token <id>', 'Fee token coin object ID (from env if not provided)')
  .option('--fee-token-type <type>', 'Fee token type (from env if not provided)')
  .option('--fee-token-metadata <id>', 'Fee token metadata object ID (from env if not provided)')
  .option('--onramp-pkg <id>', 'Onramp package id (from env if not provided)')
  .option('--pool-pkg <id>', 'Token pool package id (from env if not provided)')
  .option('--token-pool-state <id>', 'Token pool state object id (from env if not provided)')
  .option('--extra-args <bytes>', 'Extra args (0x-hex or base64)', '')
  .option('--network <net>', 'Sui network: mainnet|testnet|devnet|localnet or fullnode URL', 'testnet')
  .action(async (opts) => {
    try {
      // Load network-specific env file if network is specified
      const env = loadEnvForNetwork(opts.network)

      const privateKeyB64 = env.SUI_PRIVATE_KEY
      if (!privateKeyB64) {
        console.error('SUI_PRIVATE_KEY is required (base64-encoded secret key bytes).')
        process.exitCode = 1
        return
      }

      // Get a keypair from private key for signing
      const keypair = Ed25519Keypair.fromSecretKey(privateKeyB64)

      // Resolve fullnode URL
      const url = ['mainnet', 'testnet', 'devnet', 'localnet'].includes(opts.network)
        ? getFullnodeUrl(opts.network as 'mainnet' | 'testnet' | 'devnet' | 'localnet')
        : opts.network

      // Use env values as fallbacks
      const ccipPkg = opts.ccipPkg || env.CCIP_PACKAGE_ID
      const onrampPkg = opts.onrampPkg || env.ONRAMP_PACKAGE_ID
      const poolPkg = opts.poolPkg || (opts.poolKind === 'burn_mint' ? env.BM_POOL_PACKAGE_ID : env.LR_POOL_PACKAGE_ID)
      
      // For localnet: use ETH for burn_mint, LINK for lock_release
      // For testnet/others: always use LINK
      const isLocalnet = opts.network === 'localnet'
      const shouldUseEth = isLocalnet && opts.poolKind === 'burn_mint'
      
      const coinType = opts.coinType || (shouldUseEth ? env.ETH_COIN_TYPE : env.LINK_COIN_TYPE)
      const ccipObjectRef = opts.ccipObjectRef || env.CCIP_STATE_ID
      const onrampState = opts.onrampState || env.ONRAMP_STATE_ID
      const tokenMetadata = opts.tokenMetadata || (shouldUseEth ? env.ETH_METADATA : env.LINK_METADATA)
      const tokenCoin = opts.tokenCoin || (shouldUseEth ? env.ETH_COIN_ID : env.LINK_COIN_ID)
      const tokenPoolState = opts.tokenPoolState || (opts.poolKind === 'burn_mint' ? env.BM_POOL_STATE_ID : env.LR_POOL_STATE_ID)
      const feeToken = opts.feeToken || env.FEE_TOKEN_OBJECT
      const feeTokenType = opts.feeTokenType || env.LINK_COIN_TYPE
      const feeTokenMetadata = opts.feeTokenMetadata || env.LINK_METADATA

      if (!onrampPkg) throw new Error('onramp-pkg is required (or set ONRAMP_PACKAGE_ID in env)')
      if (!poolPkg) throw new Error('pool-pkg is required (or set LR_POOL_PACKAGE_ID/BM_POOL_PACKAGE_ID in env)')

      const client = new SuiClient({ url })
      const tx = new Transaction()

      const buildArgs: BuildArgs = {
        ccipPkg,
        onrampPkg,
        poolPkg,
        coinType,
        ccipObjectRef,
        onrampState,
        tokenMetadata,
        tokenCoin,
        tokenPoolState,
        feeToken,
        feeTokenType,
        feeTokenMetadata,
        destChainSelector: parseU64BigInt(opts.destChainSelector),
        receiver: opts.receiver,
        data: parseVecU8(opts.data),
        extraArgs: parseVecU8(opts.extraArgs),
        poolKind: opts.poolKind,
      }

      await buildCcipSendPTB(tx, client, buildArgs)

      tx.setGasBudget(3_000_000_000) // 3 billion MIST (3 SUI)

      const result = await client.signAndExecuteTransaction({ signer: keypair, transaction: tx })
      const txResult = await client.waitForTransaction({ digest: result.digest, options: { showEffects: true } })
      console.log('Transaction result:')
      console.log(JSON.stringify(txResult, null, 2))
    } catch (err: any) {
      console.error('Error:', err)
      process.exitCode = 1
    }
  })

program.parseAsync(process.argv)
