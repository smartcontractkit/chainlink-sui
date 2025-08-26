module mock_eth_token::mock_eth_token;

use sui::coin::{Self, TreasuryCap};

public struct MOCK_ETH_TOKEN has drop {}

fun init(witness: MOCK_ETH_TOKEN, ctx: &mut TxContext) {
    let (treasury, metadata) = coin::create_currency(
        witness,
        9, // decimals
        b"MOCK_ETH_TOKEN", // symbol
        b"Mock Ethereum Token", // name
        b"Mock Ethereum Token", // description
        option::none(), // icon_url
        ctx,
    );
    transfer::public_freeze_object(metadata);
    transfer::public_transfer(treasury, ctx.sender())
}

public fun mint_and_transfer(
    treasury_cap: &mut TreasuryCap<MOCK_ETH_TOKEN>,
    amount: u64,
    recipient: address,
    ctx: &mut TxContext,
) {
    let coin = coin::mint(treasury_cap, amount, ctx);
    transfer::public_transfer(coin, recipient);
}

#[allow(lint(self_transfer))]
public fun mint(treasury_cap: &mut TreasuryCap<MOCK_ETH_TOKEN>, amount: u64, ctx: &mut TxContext) {
    let coin = coin::mint(treasury_cap, amount, ctx);
    transfer::public_transfer(coin, ctx.sender());
}
