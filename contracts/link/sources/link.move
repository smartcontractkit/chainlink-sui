module link::link;

use sui::coin::{Coin, Self, TreasuryCap};
use sui::url;

public struct LINK has drop {}

// TODO: determine some of these values
fun init(witness: LINK, ctx: &mut TxContext) {
    let (treasury, metadata) = coin::create_currency(
        witness,
        9,
        b"LINK",
        b"ChainLink Token",
        b"Fill in the Chainlink Token Description",
        option::some(url::new_unsafe_from_bytes(b"https://chainlink.com/")),
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
