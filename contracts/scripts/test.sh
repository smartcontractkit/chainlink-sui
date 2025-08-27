#!/usr/bin/env bash
set -euxo pipefail

cd "$(dirname -- "$0")/.."

sui move test --path ccip/ccip
sui move test --path ccip/ccip_router
sui move test --path ccip/ccip_token_pools/token_pool
sui move test --path ccip/ccip_token_pools/managed_token_pool
sui move test --path ccip/ccip_token_pools/lock_release_token_pool
sui move test --path ccip/ccip_token_pools/burn_mint_token_pool
sui move test --path mcms/mcms
sui move test --path ccip/ccip_onramp
sui move test --path ccip/ccip_offramp
sui move test --path ccip/managed_token
