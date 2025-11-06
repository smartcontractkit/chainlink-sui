module link::link;

use sui::coin::{Self, Coin, TreasuryCap};
use sui::url;

public struct LINK has drop {}

fun init(witness: LINK, ctx: &mut TxContext) {
    let (treasury, metadata) = coin::create_currency(
        witness,
        9,
        b"LINK",
        b"ChainLink Token",
        b"The native token of the Chainlink Network",
        option::some(
            url::new_unsafe_from_bytes(
                b"https://d2f70xi62kby8n.cloudfront.net/tokens/link.webp?auto=compress%2Cformat",
            ),
        ),
        ctx,
    );

    transfer::public_freeze_object(metadata);
    transfer::public_transfer(treasury, ctx.sender());
}

public fun mint_and_transfer(
    treasury_cap: &mut TreasuryCap<LINK>,
    amount: u64,
    recipient: address,
    ctx: &mut TxContext,
) {
    let coin: Coin<LINK> = coin::mint(treasury_cap, amount, ctx);
    transfer::public_transfer(coin, recipient);
}

public fun mint(
    treasury_cap: &mut TreasuryCap<LINK>,
    amount: u64,
    ctx: &mut TxContext,
): Coin<LINK> {
    coin::mint(treasury_cap, amount, ctx)
}
