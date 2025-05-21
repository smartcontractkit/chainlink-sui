module eth_token::eth_token;

use sui::coin::{Self, TreasuryCap};

public struct ETH_TOKEN has drop {}

fun init(witness: ETH_TOKEN, ctx: &mut TxContext) {
    let (treasury, metadata) = coin::create_currency(
        witness,
        9,                  // decimals
        b"ETH_TOKEN",         // symbol
        b"",                // name
        b"",                // description
        option::none(),     // icon_url
        ctx,
    );
    transfer::public_freeze_object(metadata);
    transfer::public_transfer(treasury, ctx.sender())
}

public fun mint_and_transfer(
    treasury_cap: &mut TreasuryCap<ETH_TOKEN>,
    amount: u64,
    recipient: address,
    ctx: &mut TxContext,
) {
    let coin = coin::mint(treasury_cap, amount, ctx);
    transfer::public_transfer(coin, recipient);
}

public fun mint(
    treasury_cap: &mut TreasuryCap<ETH_TOKEN>,
    amount: u64,
    ctx: &mut TxContext,
) {
    let coin = coin::mint(treasury_cap, amount, ctx);
    transfer::public_transfer(coin, ctx.sender());
}