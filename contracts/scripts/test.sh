#!/usr/bin/env bash
set -euxo pipefail
cd "$(dirname -- "$0")/.."

PACKAGES=(
  ccip/ccip
  ccip/ccip_router
  ccip/ccip_token_pools/managed_token_pool
  ccip/ccip_token_pools/lock_release_token_pool
  ccip/ccip_token_pools/burn_mint_token_pool
  mcms/fast_mcms
  mcms/mcms
  mcms/mcms_test
  ccip/ccip_onramp
  ccip/ccip_offramp
  ccip/managed_token
  ccip/ccip_dummy_receiver
  ccip/managed_token_faucet
)

# Sui ≥1.66 uses on-chain-like gas metering in unit tests with a ~1M default budget.
# Some tests (e.g. MCMS multi-step flows) exceed that and fail as "Test timed out".
SUI_TEST_GAS_LIMIT="${SUI_TEST_GAS_LIMIT:-500000000}"

# run tests
for pkg in "${PACKAGES[@]}"; do
  sui move test --path "$pkg" --gas-limit "$SUI_TEST_GAS_LIMIT"
done
