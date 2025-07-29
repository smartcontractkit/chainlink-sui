# Onramp Integration Tests

## Notes

**Import a devnet account into the local keystore and use it for the tests**

First, install Slush browser wallet and create a new account.

```bash
sui keytool import <YOUR_PRIVATE_KEY> ed25519 --alias ccip_devnet_tests --json
```
Switch to the devnet environment using the address of the account you just imported.

```bash
sui client switch --address <YOUR_ADDRESS> --env devnet
```


Onto PTB: 

call onramp package 0x743045b0818bd0ca305e68f168be8f3ba68c3f3f7b58684f0cf49e2426f742e3::onramp::create_token_params with no args. receive a hot potato.

send to lock release token pool 0x73e8625b58b63e83d22cb7d62d9b99e6cbd7e9b435ef5656f56bdec05d5dd1b9::lock_release_token_pool::lock_or_burn with these args: 0xd62de8860aea0760ddab667438889ca98ef3f77d8c214b15557d2c2b1e5a0bf9 0x6 0xf3b4798f3436854c50c9366327c2cef68eda40e67e6f3c997782ae851c138f85 the_coin_object_id 2 the_potato. receive a hot potato

call onramp::ccip_send 0x743045b0818bd0ca305e68f168be8f3ba68c3f3f7b58684f0cf49e2426f742e3::onramp::ccip_send with 0xd62de8860aea0760ddab667438889ca98ef3f77d8c214b15557d2c2b1e5a0bf9 0xcc1836249a66bf20ba046638963ac15c128a55a7b7878ee2aba3df0caf514ac8 0x6 2 vector[] vector[] hot_potato 0x0816ae02eba12e4ea986628da3657618ae1d22955f12720edbade7c6085fcaf6 another_link_coin_for_fee vector[]


### Lessons Learned

#### Object with reference that does not exist

When you are passing an object (in this case it was a coin object)with a reference that does not exist, we will get a pure value of type NULL. You will get the following Move error: `CommandArgumentError { arg_idx: 3, kind: InvalidUsageOfPureArg } in command 1`


#### Deployed contract references a prior version of another contract

When you are passing a coin object with a reference that does not exist, we will get a pure value of type NULL. You will get the following Move error: `CommandArgumentError { arg_idx: 3, kind: InvalidUsageOfPureArg } in command 1`
