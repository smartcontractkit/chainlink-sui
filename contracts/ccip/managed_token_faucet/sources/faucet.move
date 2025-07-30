module managed_token_faucet::faucet;

use sui::deny_list::DenyList;
use sui::coin::{Coin, CoinMetadata};

use managed_token::managed_token::{Self, MintCap, TokenState};

public struct FaucetState<phantom T> has key, store {
    id: UID,
    mint_cap: MintCap<T>,
}

public fun initialize<T>(
    mint_cap: MintCap<T>,
    ctx: &mut TxContext
): FaucetState<T> {
    FaucetState<T> {
        id: object::new(ctx),
        mint_cap,
    }
}

public fun drip<T>(
    metadata: &CoinMetadata<T>,
    state: &FaucetState<T>,
    token_state: &mut TokenState<T>,
    deny_list: &DenyList,
    ctx: &mut TxContext
): Coin<T> {
    drip_internal(metadata, state, token_state, deny_list, ctx)
}

public fun drip_and_send<T>(
    metadata: &CoinMetadata<T>,
    state: &FaucetState<T>,
    token_state: &mut TokenState<T>,
    deny_list: &DenyList,
    recipient: address,
    ctx: &mut TxContext
) {
    let coin = drip_internal(metadata, state, token_state, deny_list, ctx);
    transfer::public_transfer(coin, recipient);
}

fun drip_internal<T>(
    metadata: &CoinMetadata<T>,
    state: &FaucetState<T>,
    token_state: &mut TokenState<T>,
    deny_list: &DenyList,
    ctx: &mut TxContext
): Coin<T> {
    let decimals = metadata.get_decimals();
    let mut i = 0;
    let mut amount = 1;

    while (i < decimals) {
        amount = amount * 10;
        i = i + 1;
    };

    managed_token::mint(token_state, &state.mint_cap, deny_list, amount, ctx.sender(), ctx)
}
