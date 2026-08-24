module ccip_lock_release_token::ccip_lock_release_token;

use sui::coin::{Self, TreasuryCap};

public struct CCIP_LOCK_RELEASE_TOKEN has drop {}

fun init(witness: CCIP_LOCK_RELEASE_TOKEN, ctx: &mut TxContext) {
    let (treasury, metadata) = coin::create_currency(
        witness,
        9, // decimals
        b"CCIP LnR", // symbol
        b"CCIP Lock Release Token", // name
        b"CCIP Lock Release Token for testing", // description
        option::none(), // icon_url
        ctx,
    );
    transfer::public_freeze_object(metadata);
    transfer::public_transfer(treasury, ctx.sender())
}

public fun mint_and_transfer(
    treasury_cap: &mut TreasuryCap<CCIP_LOCK_RELEASE_TOKEN>,
    amount: u64,
    recipient: address,
    ctx: &mut TxContext,
) {
    let coin = coin::mint(treasury_cap, amount, ctx);
    transfer::public_transfer(coin, recipient);
}

#[allow(lint(self_transfer))]
public fun mint(
    treasury_cap: &mut TreasuryCap<CCIP_LOCK_RELEASE_TOKEN>,
    amount: u64,
    ctx: &mut TxContext,
) {
    let coin = coin::mint(treasury_cap, amount, ctx);
    transfer::public_transfer(coin, ctx.sender());
}
