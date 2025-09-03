import { Command } from 'commander'
import { SuiClient, getFullnodeUrl } from '@mysten/sui/client'
import { Ed25519Keypair } from '@mysten/sui/keypairs/ed25519'
import { buildCcipSendPTB, type BuildArgs } from './onramp'

const privateKey = process.env.SUI_PRIVATE_KEY
const onrampPackageId = process.env.SUI_ONRAMP_PACKAGE_ID
const poolPackageId = process.env.SUI_TOKEN_POOL_ID

// Get a keypair from private key for signing
const keypair = Ed25519Keypair.fromSecretKey(privateKey!)

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
  .requiredOption('--ccip-object-ref <id>', 'CCIPObjectRef object id')
  .requiredOption('--onramp-state <id>', 'OnRampState object id')
  .requiredOption('--fee-token-metadata <id>', 'CoinMetadata<T> object id for fee coin')
  .requiredOption('--fee-token-coin <id>', 'Owned Coin<T> object id for fees')
  .requiredOption('--dest-chain-selector <u64>', 'Destination chain selector (decimal or 0x-hex)')
  .requiredOption('--receiver <bytes>', 'Receiver bytes (0x-hex or base64)')
  .option('--data <bytes>', 'Arbitrary data (0x-hex or base64)', '')
  .requiredOption('--pool-kind <kind>', 'Pool kind', (v) => {
    if (v !== 'burn_mint' && v !== 'lock_release') throw new Error('pool-kind must be burn_mint or lock_release')
    return v
  })
  .option('--onramp-pkg <id>', 'Onramp package id', onrampPackageId)
  .option('--pool-pkg <id>', 'Token pool package id', poolPackageId)
  .option('--coin-type <type>', 'Coin type', '0x2::sui::SUI')
  .option('--managed-token-state <id>', 'Managed token state object id (pool specific)', undefined)
  .option('--extra-args <bytes>', 'Extra args (0x-hex or base64)', '')
  .option('--network <net>', 'Sui network: mainnet|testnet|devnet|localnet or fullnode URL', 'testnet')
  .action(async (opts) => {
    try {
      const privateKeyB64 = process.env.SUI_PRIVATE_KEY
      if (!privateKeyB64) {
        console.error('SUI_PRIVATE_KEY is required (base64-encoded secret key bytes).')
        process.exitCode = 1
        return
      }

      // Resolve fullnode URL
      const url = ['mainnet', 'testnet', 'devnet', 'localnet'].includes(opts.network)
        ? getFullnodeUrl(opts.network as 'mainnet' | 'testnet' | 'devnet' | 'localnet')
        : opts.network

      if (!opts.onrampPkg) throw new Error('onramp-pkg is required (or set SUI_ONRAMP_PACKAGE_ID)')
      if (!opts.poolPkg) throw new Error('pool-pkg is required (or set SUI_TOKEN_POOL_ID)')

      const client = new SuiClient({ url })

      const buildArgs: BuildArgs = {
        onrampPkg: opts.onrampPkg,
        poolPkg: opts.poolPkg,
        coinType: opts.coinType,
        ccipObjectRef: opts.ccipObjectRef,
        onrampState: opts.onrampState,
        feeTokenMetadata: opts.feeTokenMetadata,
        feeTokenCoin: opts.feeTokenCoin,
        managedTokenState: opts.managedTokenState,
        destChainSelector: parseU64BigInt(opts.destChainSelector),
        receiver: parseVecU8(opts.receiver),
        data: parseVecU8(opts.data),
        extraArgs: parseVecU8(opts.extraArgs),
        poolKind: opts.poolKind,
      }

      const tx = buildCcipSendPTB(buildArgs)
      const result = await client.signAndExecuteTransaction({ signer: keypair, transaction: tx })

      console.log('Transaction submitted:')
      console.log(JSON.stringify(result, null, 2))
    } catch (err: any) {
      console.error('Error:', err?.message ?? err)
      process.exitCode = 1
    }
  })

program.parseAsync(process.argv)
