import { Command } from 'commander'
import { SuiClient, getFullnodeUrl } from '@mysten/sui/client'
import { Ed25519Keypair } from '@mysten/sui/keypairs/ed25519'
import { buildCcipSendPTB, type BuildArgs } from './onramp'
import { env } from './env'
import { Transaction, type TransactionResult } from '@mysten/sui/transactions'

const privateKey = env.SUI_PRIVATE_KEY

if (!privateKey) {
  console.error('SUI_PRIVATE_KEY is required')
  process.exitCode = 1
  process.exit(1)
}

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
  .requiredOption('--ccip-pkg <id>', 'CCIP package id', env.CCIP_PACKAGE_ID)
  .requiredOption('--ccip-object-ref <id>', 'CCIPObjectRef object id', env.CCIP_STATE_ID)
  .requiredOption('--onramp-state <id>', 'OnRampState object id', env.ONRAMP_STATE_ID)
  .requiredOption('--token-metadata <id>', 'CoinMetadata<T> for token being transferred', env.LINK_METADATA)
  .requiredOption('--token-coin <id>', 'Owned Coin<T> for token being transferred', env.LINK_COIN_ID)
  .requiredOption('--coin-type <type>', 'Token type being transferred', env.LINK_COIN_TYPE)
  .requiredOption('--dest-chain-selector <u64>', 'Destination chain selector (decimal or 0x-hex)')
  .requiredOption('--receiver <address>', 'Receiver address')
  .option('--data <bytes>', 'Arbitrary data (0x-hex or base64)', '')
  .requiredOption('--pool-kind <kind>', 'Pool kind', (v) => {
    if (v !== 'burn_mint' && v !== 'lock_release') throw new Error('pool-kind must be burn_mint or lock_release')
    return v
  })
  .requiredOption('--fee-token <id>', 'Fee token coin object ID (defaults to gas coin)', env.FEE_TOKEN_OBJECT)
  .option('--fee-token-type <type>', 'Fee token type', env.LINK_COIN_TYPE)
  .option('--fee-token-metadata <id>', 'Fee token metadata object ID', env.LINK_METADATA)
  .option('--onramp-pkg <id>', 'Onramp package id', env.ONRAMP_PACKAGE_ID)
  .option('--pool-pkg <id>', 'Token pool package id', env.LR_POOL_PACKAGE_ID)
  .option('--token-pool-state <id>', 'Token pool state object id (pool specific)', env.LR_POOL_STATE_ID)
  .option('--extra-args <bytes>', 'Extra args (0x-hex or base64)', '')
  .option('--network <net>', 'Sui network: mainnet|testnet|devnet|localnet or fullnode URL', 'testnet')
  .action(async (opts) => {
    try {
      const privateKeyB64 = env.SUI_PRIVATE_KEY
      if (!privateKeyB64) {
        console.error('SUI_PRIVATE_KEY is required (base64-encoded secret key bytes).')
        process.exitCode = 1
        return
      }

      // Resolve fullnode URL
      const url = ['mainnet', 'testnet', 'devnet', 'localnet'].includes(opts.network)
        ? getFullnodeUrl(opts.network as 'mainnet' | 'testnet' | 'devnet' | 'localnet')
        : opts.network

      if (!opts.onrampPkg) throw new Error('onramp-pkg is required (or set ONRAMP_PACKAGE_ID)')
      if (!opts.poolPkg) throw new Error('pool-pkg is required (or set LR_POOL_PACKAGE_ID)')

      const client = new SuiClient({ url })
      const tx = new Transaction()

      const buildArgs: BuildArgs = {
        ccipPkg: opts.ccipPkg,
        onrampPkg: opts.onrampPkg,
        poolPkg: opts.poolPkg,
        coinType: opts.coinType,
        ccipObjectRef: opts.ccipObjectRef,
        onrampState: opts.onrampState,
        tokenMetadata: opts.tokenMetadata,
        tokenCoin: opts.tokenCoin,
        tokenPoolState: opts.tokenPoolState,
        feeToken: opts.feeToken,
        feeTokenType: opts.feeTokenType,
        feeTokenMetadata: opts.feeTokenMetadata,
        destChainSelector: parseU64BigInt(opts.destChainSelector),
        receiver: opts.receiver,
        data: parseVecU8(opts.data),
        extraArgs: parseVecU8(opts.extraArgs),
        poolKind: opts.poolKind,
      }

      await buildCcipSendPTB(tx, client, buildArgs)

      tx.setGasBudget(1_000_000_000_000)
      tx.setGasPrice(1_000_000)

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
