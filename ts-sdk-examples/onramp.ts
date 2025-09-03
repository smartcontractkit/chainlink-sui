import { Transaction } from '@mysten/sui/transactions'

// Onramp call args
export type BuildArgs = {
  // Packages + types
  onrampPkg: string;                 // e.g. "0x…"
  poolPkg: string;                   // burn_mint or lock_release package id
  coinType: string;                  // e.g. "0x2::sui::SUI"
  // Objects (owned/shared)
  ccipObjectRef: string;             // CCIPObjectRef
  onrampState: string;               // OnRampState
  feeTokenMetadata: string;          // &CoinMetadata<T> used for fee
  feeTokenCoin: string;              // &mut Coin<T> used for fee payment (owned)
  managedTokenState?: string;        // pool-specific state obj (if required)
  // Constants / params
  destChainSelector: bigint;         // u64
  receiver: Uint8Array;              // vector<u8>
  data: Uint8Array;                  // vector<u8>
  extraArgs?: Uint8Array;            // vector<u8> (optional)
  // Token pool type
  poolKind: 'burn_mint' | 'lock_release';
};

export function buildCcipSendPTB(a: BuildArgs) {
  const tx = new Transaction()

  const ccipRef = tx.object(a.ccipObjectRef)
  const state = tx.object(a.onrampState)
  const clock = tx.object.clock()

  // Create Token State Params
  const tokenParams = tx.moveCall({
    package: a.onrampPkg,
    module: 'onramp_state_helper',
    function: 'create_token_transfer_params',
    typeArguments: [a.coinType],
    arguments: [
      ccipRef,
      state,
      clock
      // ...
    ]
  })

  // Token Pool call based on poolKind
  switch (a.poolKind) {
    case 'burn_mint':
      tx.moveCall({
        package: a.poolPkg,
        module: 'burn_mint_token_pool',
        function: 'lock_or_burn',
        typeArguments: [a.coinType],
        arguments: [
          ccipRef,                           // ccip_object_ref
          tokenParams,                       // token_transfer_params (cmd 0 result)
          tx.object(a.feeTokenCoin),         // coin (if the pool needs a Coin<T> input)
          tx.pure.u64(a.destChainSelector),
          clock,
          tx.object.denyList(),              // deny_list (doc: always 0x403)
          tx.object(a.managedTokenState!)    // managed_token_state (if required)
        ]
      })
      break
    case 'lock_release':
    default:
      tx.moveCall({
        package: a.poolPkg,
        module: 'lock_release_token_pool',
        function: 'lock_or_burn',
        typeArguments: [a.coinType],
        arguments: [
          ccipRef,
          tokenParams,
          tx.object(a.feeTokenCoin),
          tx.pure.u64(a.destChainSelector),
          clock,
          // pool’s state object for lock/release variant
          tx.object(a.managedTokenState!)
        ]
      })
      break
  }

  // CCIP Send
  tx.moveCall({
    package: a.poolPkg,
    module: 'onramp',
    function: 'ccip_send',
    typeArguments: [a.coinType],
    arguments: [
      ccipRef,                                 // &mut CCIPObjectRef
      state,                                   // &mut OnRampState
      clock,                                   // &Clock (0x6)
      tx.pure.u64(a.destChainSelector),        // u64
      tx.pure(a.receiver),                     // vector<u8>
      tx.pure(a.data),                         // vector<u8>
      tokenParams,                             // osh::TokenTransferParams (cmd 0)
      tx.object(a.feeTokenMetadata),           // &CoinMetadata<T>
      tx.object(a.feeTokenCoin),               // &mut Coin<T>
      tx.pure(a.extraArgs ?? new Uint8Array())// vector<u8>
    ]
  })

  return tx
}


/**
 // Example usage
 const ptb = buildCcipSendPTB({
 onrampPkg: onrampPackageId!,
 poolPkg: '0xPOOL…',
 coinType: '0x2::sui::SUI',
 ccipObjectRef: '0x…',
 onrampState: '0x…',
 feeTokenMetadata: '0x…',
 feeTokenCoin: '0x…',
 destChainSelector: 16015286601757825753n,
 receiver: new Uint8Array([...]),
 data: new Uint8Array(),
 poolKind: 'burn_mint'
 })

 // Sign & execute with your keypair:
 const ptbResponse = await suiClient.signAndExecuteTransaction({ signer: keypair, transaction: ptb })
 */