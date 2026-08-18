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
  ccip/ccip_malicious_receiver
  ccip/managed_token_faucet
)

# Sui ≥1.66 uses on-chain-like gas metering in unit tests with a ~1M default budget.
# Some tests (e.g. MCMS multi-step flows, BCS deserialization of large extra_args) exceed that.
SUI_TEST_GAS_LIMIT="${SUI_TEST_GAS_LIMIT:-2000000000}"
SUI_BUILD_ENV="${SUI_BUILD_ENV:-testnet}"

# run tests
for pkg in "${PACKAGES[@]}"; do
  sui move test --path "$pkg" --build-env "$SUI_BUILD_ENV" --gas-limit "$SUI_TEST_GAS_LIMIT"
done
