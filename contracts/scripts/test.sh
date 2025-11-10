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

patch_toml() {
  local file="$1"
  [[ ! -f "$file" ]] && return
  cp "$file" "$file.bak"

  # replace only inside [addresses] section
  awk '
    BEGIN { in_addr=0 }
    /^\[addresses\]/     { in_addr=1 }
    /^\[dev-addresses\]/ { in_addr=0 }
    {
      if (in_addr && $0 ~ /"0x0"/) gsub(/"0x0"/,"\"_\"")
      print
    }
  ' "$file" > "$file.tmp" && mv "$file.tmp" "$file"
}

# patch Move.toml files (targets + mcms dependency)
for pkg in "${PACKAGES[@]}"; do
  patch_toml "$pkg/Move.toml"
done

# restore originals even on error 
trap 'for pkg in "${PACKAGES[@]}"; do f="$pkg/Move.toml.bak"; [[ -f $f ]] && mv "$f" "${f%.bak}"; done' EXIT

# run tests
for pkg in "${PACKAGES[@]}"; do
  sui move test --path "$pkg"
done
