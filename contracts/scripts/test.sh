#!/usr/bin/env bash
set -euxo pipefail
cd "$(dirname -- "$0")/.."

PACKAGES=(
  ccip/ccip
  ccip/ccip_router
  ccip/ccip_token_pools/managed_token_pool
  ccip/ccip_token_pools/lock_release_token_pool
  ccip/ccip_token_pools/burn_mint_token_pool
  mcms/mcms
  mcms/mcms_test
  mcms/mcms_test_v2
  ccip/ccip_onramp
  ccip/ccip_offramp
  ccip/managed_token
  ccip/ccip_dummy_receiver
  ccip/managed_token_faucet
)

# run tests
for pkg in "${PACKAGES[@]}"; do
  sui move test --path "$pkg"
done
