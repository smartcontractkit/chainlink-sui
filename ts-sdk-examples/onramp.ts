import { Transaction, type TransactionArgument, type TransactionResult } from '@mysten/sui/transactions'
import { SuiClient } from '@mysten/sui/client'

// Onramp call args
export type BuildArgs = {
  // Packages + types
  ccipPkg: string;                  // CCIP package id
  onrampPkg: string;                // e.g. "0x…"
  poolPkg: string;                  // burn_mint or lock_release package id
  coinType: string;                 // e.g. "0xb9bf...::mock_link_token::MOCK_LINK_TOKEN"
  // Objects (owned/shared)
  ccipObjectRef: string;            // CCIPObjectRef
  onrampState: string;              // OnRampState
  tokenMetadata: string;            // &CoinMetadata<T> for token being transferred
  tokenCoin: string;                // &mut Coin<T> token being transferred (owned)
  feeToken: TransactionResult | string;                 // &mut Coin<FeeToken> for fee payment (owned)
  feeTokenType: string;             // Fee token type (e.g., "0x2::sui::SUI")
  feeTokenMetadata: string;         // &CoinMetadata<FeeToken> for fee token
  tokenPoolState: string;           // pool-specific state obj (if required)
  // Constants / params
  destChainSelector: bigint;        // u64
  receiver: string;                 // address
  data: Uint8Array;                 // vector<u8>
  extraArgs?: Uint8Array;           // vector<u8> (optional)
  // Token pool type
  poolKind: 'burn_mint' | 'lock_release';
};

export async function buildCcipSendPTB(tx: Transaction, client: SuiClient, a: BuildArgs) {
  console.debug('BuildArgs', a)

  const state = tx.object(a.onrampState)
  const clock = tx.object('0x6')
  const receiverBytes = Array.from(Buffer.from(a.receiver.replace('0x', ''), 'hex'))
  const receiver = tx.pure.vector('u8', receiverBytes)

  // Create Token State Params
  const tokenParams = tx.moveCall({
    package: a.ccipPkg,
    module: 'onramp_state_helper',
    function: 'create_token_transfer_params',
    arguments: [
      receiver,
    ]
  })

  // Token Pool call based on poolKind
  switch (a.poolKind) {
    case 'burn_mint':
      console.debug('Calling burn_mint_token_pool')
      tx.moveCall({
        package: a.poolPkg,
        module: 'burn_mint_token_pool',
        function: 'lock_or_burn',
        typeArguments: [a.coinType],
        arguments: [
          tx.object(a.ccipObjectRef),     // ccip_object_ref
          tokenParams,                    // token_transfer_params (cmd 0 result)
          tx.object(a.tokenCoin),         // coin (if the pool needs a Coin<T> input)
          tx.pure.u64(a.destChainSelector),
          clock,
          tx.object(a.tokenPoolState!)    // managed_token_state (if required)
        ]
      })
      break
    case 'lock_release':
    default:
      console.debug('Calling lock_release_token_pool', {
        destChainSelector: a.destChainSelector
      })

      tx.moveCall({
        package: a.poolPkg,
        module: 'lock_release_token_pool',
        function: 'lock_or_burn',
        typeArguments: [a.coinType],
        arguments: [
          tx.object(a.ccipObjectRef),
          tokenParams,
          tx.object(a.tokenCoin),
          tx.pure.u64(a.destChainSelector),
          clock,
          tx.object(a.tokenPoolState!)
        ]
      })
      break
  }

  console.log('receiverBytes', receiverBytes, receiverBytes.length)

  // CCIP Send
  tx.moveCall({
    package: a.onrampPkg,
    module: 'onramp',
    function: 'ccip_send',
    typeArguments: [a.feeTokenType],  // Use fee token type, not transfer token type
    arguments: [
      tx.object(a.ccipObjectRef),           // &mut CCIPObjectRef
      state,                                // &mut OnRampState
      clock,                                // &Clock (0x6)
      tx.pure.u64(a.destChainSelector),     // u64
      tx.pure('vector<u8>', receiverBytes), // vector<u8>
      tx.pure.vector('u8', a.data),         // vector<u8>
      tokenParams,                          // osh::TokenTransferParams (cmd 0)
      tx.object(a.feeTokenMetadata),        // &CoinMetadata<FeeToken>
      tx.object(a.feeToken),                // &mut Coin<FeeToken> (separate from transfer token)
      tx.pure.vector('u8', a.extraArgs ?? new Uint8Array([])) // vector<u8>
    ]
  })

  return tx
}
