#!/usr/bin/env bash
set -euo pipefail

echo "Starting deployment script"
echo "NOTE: this script is for testing and documentation purposes only!"
echo "Please use Changesets and CLD for any deployments that need to be in production."

# --------------------------------
# Config (override via environment)
# --------------------------------
SCRIPT_DIR="$(cd -- "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"
ROOT_DIR="${ROOT_DIR:-$SCRIPT_DIR/../contracts}"
GAS="${GAS:-600000000}"

# Local chain selector for THIS Sui network (used by rmn_remote & onramp init)
LOCAL_CHAIN_SELECTOR="${LOCAL_CHAIN_SELECTOR:-1}"

# OnRamp destination-chain arrays
DEST_CHAIN_SELECTORS_JSON="${DEST_CHAIN_SELECTORS_JSON:-[2]}"
DEST_CHAIN_ENABLED_JSON="${DEST_CHAIN_ENABLED_JSON:-[true]}"
DEST_CHAIN_ALLOWLIST_ENABLED_JSON="${DEST_CHAIN_ALLOWLIST_ENABLED_JSON:-[false]}"

# Admin-ish addresses (can reuse the active Sui address)
FEE_AGGREGATOR_ADDR="${FEE_AGGREGATOR_ADDR:-$(sui client active-address 2>/dev/null || echo 0x0)}"
ALLOWLIST_ADMIN_ADDR="${ALLOWLIST_ADMIN_ADDR:-$(sui client active-address 2>/dev/null || echo 0x0)}"
REBALANCER_ADDR="${REBALANCER_ADDR:-$(sui client active-address 2>/dev/null || echo 0x0)}"

# FeeQuoter parameters (examples; tune as you need)
FEE_QUOTER_LINK_RATE_WEI="${FEE_QUOTER_LINK_RATE_WEI:-1000000000000000000}"  # 1e18
FEE_QUOTER_BASE_FEE="${FEE_QUOTER_BASE_FEE:-90000000000}"

# Token type(s) used by the pools
# Example: LINK as the coin for lock_release; an ETH-like mock for burn_mint
#LR_COIN_TYPE="${LR_COIN_TYPE:-link_token::LINK_TOKEN}"      # module path within its package
OWNER="${OWNER:-$(sui client active-address)}"

CLOCK_ID="0x6"   # well-known
DENY_LIST_ID="0x403"  # per PTB doc
need() { command -v "$1" >/dev/null 2>&1 || { echo "Missing $1"; exit 1; }; }
need sui; need jq; need sed

# Optional: if the RPC is running locally, faucet the active-account
#sui client faucet

patch_move_toml() {
  local file="$1" key="$2" val="$3"
  # remove newlines just in case
  val="${val//$'\n'/}"
  if grep -Eq "^[[:space:]]*$key[[:space:]]*=" "$file"; then
    sed -Ei.bak "s|^([[:space:]]*$key[[:space:]]*=[[:space:]]*\").*(\"[[:space:]]*)$|\1$val\2|" "$file"
  else
    echo "WARN: key '$key' not found in $file (skipping)"
  fi
}

publish_and_pin() {
  local dir="$1" key="$2" tag="$3"
  local toml="$dir/Move.toml"
  local out_file="artifacts.${tag}.publish.json"

  # log to stderr
  >&2 echo "==> Publishing $tag: $dir"
  patch_move_toml "$toml" "$key" "0x0"

  pushd "$dir" >/dev/null
  sui client publish --with-unpublished-dependencies --gas-budget "$GAS" --json --silence-warnings \
    | tee >(jq -C . >&2) \
    > "$OLDPWD/$out_file"
  popd >/dev/null

  local pkg
  pkg="$(jq -r '.objectChanges[] | select(.type=="published") | .packageId' "$out_file" | head -n1)"
  # strip any stray newlines (paranoia)
  pkg="${pkg//$'\n'/}"
  if [[ -z "$pkg" || "$pkg" == "null" ]]; then
    >&2 echo "ERROR: packageId not found for $tag; see $out_file"
    exit 1
  fi

  >&2 echo "    ${tag} packageId: $pkg"
  patch_move_toml "$toml" "$key" "$pkg"

  # ONLY the id on stdout (so callers can capture cleanly)
  printf '%s' "$pkg"
}

extract_created() { jq -r --arg k "$2" '.objectChanges[] | select(.type=="created" and (.objectType|contains($k))) | .objectId' "$1"; }

echo "--- Deploying Mock LINK (for fee quoter, LR pool) ---"
LINK_DIR="$ROOT_DIR/ccip/mock_link_token"
LINK_PKG_KEY="mock_link_token"
LINK_PKG_ID="$(publish_and_pin "$LINK_DIR" "$LINK_PKG_KEY" "link")"
LINK_METADATA_ID="$(extract_created artifacts.link.publish.json '::coin::CoinMetadata<' | head -n1)"
LINK_TREASURY_CAP_ID="$(extract_created artifacts.link.publish.json '::coin::TreasuryCap<' | head -n1)"
echo "    LINK metadata: $LINK_METADATA_ID"
echo "    LINK treasury cap: $LINK_TREASURY_CAP_ID"

echo "--- Deploying MCMS ---"
MCMS_DIR="$ROOT_DIR/mcms/mcms"
patch_move_toml "$MCMS_DIR/Move.toml" "mcms_owner" "$OWNER"
MCMS_PKG_ID="$(publish_and_pin "$MCMS_DIR" "mcms" "mcms")"

echo "--- Deploying CCIP core ---"
CCIP_DIR="$ROOT_DIR/ccip/ccip"
CCIP_PKG_KEY="ccip"
patch_move_toml "$CCIP_DIR/Move.toml" "mcms" "$MCMS_PKG_ID"
patch_move_toml "$CCIP_DIR/Move.toml" "mcms_owner" "$OWNER"
CCIP_PKG_ID="$(publish_and_pin "$CCIP_DIR" "$CCIP_PKG_KEY" "ccip")"

CCIP_STATE_REF_ID="$(
  jq -r '
    .objectChanges[]
    | select(.type=="created" and (.objectType | test("::state_object::CCIPObjectRef$")))
    | .objectId
  ' artifacts.ccip.publish.json | head -n1
)"
CCIP_OWNER_CAP_ID="$(jq -r '.objectChanges[] | select(.type=="created" and (.objectType|test("OwnerCap"))) | .objectId' artifacts.ccip.publish.json | head -n1)"
CCIP_SOURCE_TRANSFER_CAP_ID="$(jq -r '.objectChanges[] | select(.type=="created" and (.objectType|test("Source.*Transfer.*Cap|source.*transfer.*cap"; "i"))) | .objectId' artifacts.ccip.publish.json | head -n1)"
[[ -n "$CCIP_STATE_REF_ID" && -n "$CCIP_OWNER_CAP_ID" ]] || { echo "Missing CCIP state/owner cap"; exit 1; }

sui client call --package "$CCIP_PKG_ID" --module upgrade_registry --function initialize \
  --args "$CCIP_STATE_REF_ID" "$CCIP_OWNER_CAP_ID" \
  --gas-budget "$GAS" --json | tee artifacts.ccip.upgrade_registry.init.json >/dev/null
  
# fee_quoter::initialize (uses LINK and SUI as fee tokens)
sui client call \
  --package "$CCIP_PKG_ID" --module fee_quoter --function initialize \
  --args "$CCIP_STATE_REF_ID" "$CCIP_OWNER_CAP_ID" \
        "$FEE_QUOTER_LINK_RATE_WEI" "$LINK_METADATA_ID" "$FEE_QUOTER_BASE_FEE" \
        "[\"$LINK_METADATA_ID\"]" \
  --gas-budget "$GAS" --json | tee artifacts.ccip.fee_quoter.init.json >/dev/null

# nonce_manager, receiver_registry, rmn_remote, token_admin_registry
sui client call --package "$CCIP_PKG_ID" --module nonce_manager --function initialize \
  --args "$CCIP_STATE_REF_ID" "$CCIP_OWNER_CAP_ID" \
  --gas-budget "$GAS" --json | tee artifacts.ccip.nonce_manager.init.json >/dev/null
NONCE_MANAGER_CAP_ID="$(jq -r '.objectChanges[] | select(.type=="created" and (.objectType|test("NonceManagerCap"))) | .objectId' artifacts.ccip.nonce_manager.init.json | head -n1)"

sui client call --package "$CCIP_PKG_ID" --module receiver_registry --function initialize \
  --args "$CCIP_STATE_REF_ID" "$CCIP_OWNER_CAP_ID" \
  --gas-budget "$GAS" --json | tee artifacts.ccip.receiver_registry.init.json >/dev/null

sui client call --package "$CCIP_PKG_ID" --module rmn_remote --function initialize \
  --args "$CCIP_STATE_REF_ID" "$CCIP_OWNER_CAP_ID" "$LOCAL_CHAIN_SELECTOR" \
  --gas-budget "$GAS" --json | tee artifacts.ccip.rmn_remote.init.json >/dev/null

sui client call --package "$CCIP_PKG_ID" --module token_admin_registry --function initialize \
  --args "$CCIP_STATE_REF_ID" "$CCIP_OWNER_CAP_ID" \
  --gas-budget "$GAS" --json | tee artifacts.ccip.token_admin_registry.init.json >/dev/null

echo "--- Deploying OnRamp ---"
ONRAMP_DIR="$ROOT_DIR/ccip/ccip_onramp"
patch_move_toml "$ONRAMP_DIR/Move.toml" "ccip" "$CCIP_PKG_ID"
ONRAMP_KEY="ccip_onramp"
ONRAMP_PKG_ID="$(publish_and_pin "$ONRAMP_DIR" "$ONRAMP_KEY" "onramp")"
ONRAMP_STATE_ID="$(
  jq -r '
    .objectChanges[]
    | select(.type=="created" and (.objectType | test("::onramp::OnRampState$")))
    | .objectId
  ' artifacts.onramp.publish.json | head -n1
)"
ONRAMP_OWNER_CAP_ID="$(jq -r '.objectChanges[] | select(.type=="created" and (.objectType|test("OwnerCap"))) | .objectId' artifacts.onramp.publish.json | head -n1)"
[[ -n "$ONRAMP_STATE_ID" && -n "$ONRAMP_OWNER_CAP_ID" ]] || { echo "Missing OnRamp state/owner cap"; exit 1; }

# NOTE: using a random address for the router
sui client call --package "$ONRAMP_PKG_ID" --module onramp --function initialize \
  --args "$ONRAMP_STATE_ID" "$ONRAMP_OWNER_CAP_ID" "$NONCE_MANAGER_CAP_ID" \
        "$CCIP_SOURCE_TRANSFER_CAP_ID" "$LOCAL_CHAIN_SELECTOR" \
        "$FEE_AGGREGATOR_ADDR" "$ALLOWLIST_ADMIN_ADDR" \
        "$DEST_CHAIN_SELECTORS_JSON" "$DEST_CHAIN_ALLOWLIST_ENABLED_JSON" '["0x4488418e4980acbb2c83b8ce98ba1b3f557dd37b392e95bbd98215233cdd5ed3"]' \
  --gas-budget "$GAS" --json | tee artifacts.onramp.init.json >/dev/null

echo "--- Deploying & initializing lock_release_token_pool ---"
LR_DIR="$ROOT_DIR/ccip/ccip_token_pools/lock_release_token_pool"
patch_move_toml "$LR_DIR/Move.toml" "ccip" "$CCIP_PKG_ID"
LR_PKG_ID="$(publish_and_pin "$LR_DIR" "lock_release_token_pool" "lock_release_tp")"

LINK_COIN_T="$(
  jq -r '
    .objectChanges[]
    | select(.type=="created" and (.objectType|test("::coin::CoinMetadata<")))
    | .objectType
  ' artifacts.link.publish.json \
    | sed -E 's/^.*CoinMetadata<([^>]+)>.*/\1/' \
    | head -n1
)"

echo "Detected LINK coin type: $LINK_COIN_T"

LR_OWNER_CAP_ID="$(jq -r '.objectChanges[] | select(.type=="created" and (.objectType|test("OwnerCap"))) | .objectId' artifacts.lock_release_tp.publish.json | head -n1)"

# lock_release_token_pool::initialize(ccip_ref, LINK metadata, LINK treasury, package_id, rebalancer)
sui client call --package "$LR_PKG_ID" --module lock_release_token_pool --function initialize \
  --type-args "$LINK_COIN_T" \
  --args "$LR_OWNER_CAP_ID" "$CCIP_STATE_REF_ID" "$LINK_METADATA_ID" "$LINK_TREASURY_CAP_ID" "$LR_PKG_ID" "$REBALANCER_ADDR" \
  --gas-budget "$GAS" --json | tee artifacts.lr_tp.init.json >/dev/null
LR_STATE_ID="$(jq -r '.objectChanges[] | select(.type=="created" and (.objectType|test("LockReleaseTokenPoolState"))) | .objectId' artifacts.lr_tp.init.json | head -n1)"

# Optional: apply_chain_updates + rate limiter (example: add chain 2)
sui client call --package "$LR_PKG_ID" --module lock_release_token_pool --function apply_chain_updates \
 --type-args "$LINK_COIN_T" \
 --args "$LR_STATE_ID" "$LR_OWNER_CAP_ID" "[]" "[2]" "[[[24, 42, 24, 42]]]" "[[0,0,0,0,0,0,0,0,0,0,0,0,24,42,24,42,24,42,24,42,24,42,24,42,24,42,24,42,24,42,24,42]]" \
 --gas-budget "$GAS" --json | tee artifacts.lr_tp.apply_chains.json >/dev/null

 sui client call --package "$LR_PKG_ID" --module lock_release_token_pool --function set_chain_rate_limiter_config \
  --type-args "$LINK_COIN_T" \
  --args "$LR_STATE_ID" "$LR_OWNER_CAP_ID" "$CLOCK_ID" "2" "false" "200000000000" "20000000000" "false" "200000000000" "20000000000" \
  --gas-budget "$GAS" --json | tee artifacts.lr_tp.rate_limiters.json >/dev/null

echo "--- Deploying mock ETH token for burn/mint pool (if not present) ---"
ETH_DIR="$ROOT_DIR/ccip/mock_eth_token"
ETH_PKG_KEY="mock_eth_token"
ETH_PKG_ID="$(publish_and_pin "$ETH_DIR" "$ETH_PKG_KEY" "eth")"
ETH_METADATA_ID="$(extract_created artifacts.eth.publish.json '::coin::CoinMetadata<' | head -n1)"
ETH_TREASURY_CAP_ID="$(extract_created artifacts.eth.publish.json '::coin::TreasuryCap<' | head -n1)"

echo "--- Minting ETH tokens ---"

# Mint some ETH so you can later test burns
sui client call --package "$ETH_PKG_ID" --module mock_eth_token --function mint \
 --args "$ETH_TREASURY_CAP_ID" "1000000000000000" \
 --gas-budget "$GAS" --json | tee artifacts.eth.mint.json >/dev/null

ETH_COIN_ID="$(jq -r '.objectChanges[] | select(.type=="created" and (.objectType|test("::coin::Coin<"))) | .objectId' artifacts.eth.mint.json | head -n1)"

echo "--- Minting LINK tokens ---"

 # Mint some LINK tokens (adjust amount as needed, respecting decimals)
sui client call \
  --package "$LINK_PKG_ID" \
  --module mock_link_token \
  --function mint \
  --args "$LINK_TREASURY_CAP_ID" "1000000000000000" \
  --gas-budget "$GAS" \
  --json | tee artifacts.link.mint.json >/dev/null

LINK_COIN_ID="$(jq -r '.objectChanges[] | select(.type=="created" and (.objectType|test("::coin::Coin<"))) | .objectId' artifacts.link.mint.json | head -n1)"

echo "--- Deploying & initializing burn_mint_token_pool ---"

BM_DIR="$ROOT_DIR/ccip/ccip_token_pools/burn_mint_token_pool"
patch_move_toml "$BM_DIR/Move.toml" "ccip" "$CCIP_PKG_ID"
BM_PKG_ID="$(publish_and_pin "$BM_DIR" "burn_mint_token_pool" "burn_mint_tp")"

ETH_COIN_T="$(
  jq -r '
    .objectChanges[]
    | select(.type=="created" and (.objectType|test("::coin::CoinMetadata<")))
    | .objectType
  ' artifacts.eth.publish.json \
    | sed -E 's/^.*CoinMetadata<([^>]+)>.*/\1/' \
    | head -n1
)"

echo "Detected ETH coin type: $ETH_COIN_T"

BM_OWNER_CAP_ID="$(jq -r '.objectChanges[] | select(.type=="created" and (.objectType|test("OwnerCap"))) | .objectId' artifacts.burn_mint_tp.publish.json | head -n1)"

# burn_mint_token_pool::initialize(ccip_ref, ETH metadata, ETH treasury cap, package_id) + <T>
sui client call --package "$BM_PKG_ID" --module burn_mint_token_pool --function initialize \
  --type-args "$ETH_COIN_T" \
  --args "$BM_OWNER_CAP_ID" "$CCIP_STATE_REF_ID" "$ETH_METADATA_ID" "$ETH_TREASURY_CAP_ID" "$BM_PKG_ID" \
  --gas-budget "$GAS" --json | tee artifacts.bm_tp.init.json >/dev/null
BM_STATE_ID="$(jq -r '.objectChanges[] | select(.type=="created" and (.objectType|test("BurnMintTokenPoolState"))) | .objectId' artifacts.bm_tp.init.json | head -n1)"

# Add chain 2; set basic rate limiters (example values)
sui client call --package "$BM_PKG_ID" --module burn_mint_token_pool --function apply_chain_updates \
 --type-args "$ETH_COIN_T" \
 --args "$BM_STATE_ID" "$BM_OWNER_CAP_ID" "[]" "[2]" "[[[24, 42, 24, 42]]]" "[[0,0,0,0,0,0,0,0,0,0,0,0,24,42,24,42,24,42,24,42,24,42,24,42,24,42,24,42,24,42,24,42]]" \
 --gas-budget "$GAS" --json | tee artifacts.bm_tp.apply_chains.json >/dev/null

sui client call --package "$BM_PKG_ID" --module burn_mint_token_pool --function set_chain_rate_limiter_config \
  --type-args "$ETH_COIN_T" \
  --args "$BM_STATE_ID" "$BM_OWNER_CAP_ID" "$CLOCK_ID" "2" "false" "200000000000" "20000000000" "false" "200000000000" "20000000000" \
  --gas-budget "$GAS" --json | tee artifacts.bm_tp.rate_limiters.json >/dev/null

echo "--- Deploying mock USDC token for managed token pool ---"
USDC_DIR="$ROOT_DIR/ccip/mock_eth_token"
USDC_PKG_KEY="mock_eth_token"
# Reusing mock_eth_token structure but treating it as USDC for this example
USDC_PKG_ID="$(publish_and_pin "$USDC_DIR" "$USDC_PKG_KEY" "usdc")"
USDC_METADATA_ID="$(extract_created artifacts.usdc.publish.json '::coin::CoinMetadata<' | head -n1)"
USDC_TREASURY_CAP_ID="$(extract_created artifacts.usdc.publish.json '::coin::TreasuryCap<' | head -n1)"

USDC_COIN_T="$(
  jq -r '
    .objectChanges[]
    | select(.type=="created" and (.objectType|test("::coin::CoinMetadata<")))
    | .objectType
  ' artifacts.usdc.publish.json \
    | sed -E 's/^.*CoinMetadata<([^>]+)>.*/\1/' \
    | head -n1
)"

echo "Detected USDC coin type: $USDC_COIN_T"

echo "--- Deploying & initializing managed_token ---"
MANAGED_TOKEN_DIR="$ROOT_DIR/ccip/managed_token"
patch_move_toml "$MANAGED_TOKEN_DIR/Move.toml" "mcms" "$MCMS_PKG_ID"
patch_move_toml "$MANAGED_TOKEN_DIR/Move.toml" "mcms_owner" "$OWNER"
MANAGED_TOKEN_PKG_ID="$(publish_and_pin "$MANAGED_TOKEN_DIR" "managed_token" "managed_token")"

# Initialize managed_token with the USDC treasury cap
MANAGED_TOKEN_PUBLISHER_ID="$(jq -r '.objectChanges[] | select(.type=="created" and (.objectType|test("Publisher"))) | .objectId' artifacts.managed_token.publish.json | head -n1)"
sui client call --package "$MANAGED_TOKEN_PKG_ID" --module managed_token --function initialize \
  --type-args "$USDC_COIN_T" \
  --args "$USDC_TREASURY_CAP_ID" "$MANAGED_TOKEN_PUBLISHER_ID" \
  --gas-budget "$GAS" --json | tee artifacts.managed_token.init.json >/dev/null

MANAGED_TOKEN_STATE_ID="$(jq -r '.objectChanges[] | select(.type=="created" and (.objectType|test("TokenState"))) | .objectId' artifacts.managed_token.init.json | head -n1)"
MANAGED_TOKEN_OWNER_CAP_ID="$(jq -r '.objectChanges[] | select(.type=="created" and (.objectType|test("OwnerCap"))) | .objectId' artifacts.managed_token.init.json | head -n1)"

echo "  Managed Token State: $MANAGED_TOKEN_STATE_ID"
echo "  Managed Token Owner Cap: $MANAGED_TOKEN_OWNER_CAP_ID"

# Configure a new minter and issue a MintCap (unlimited allowance for token pool)
ACTIVE_ADDR="$(sui client active-address 2>/dev/null)"
sui client call --package "$MANAGED_TOKEN_PKG_ID" --module managed_token --function configure_new_minter \
  --type-args "$USDC_COIN_T" \
  --args "$MANAGED_TOKEN_STATE_ID" "$MANAGED_TOKEN_OWNER_CAP_ID" "$ACTIVE_ADDR" "0" "true" \
  --gas-budget "$GAS" --json | tee artifacts.managed_token.mint_cap.json >/dev/null

MINT_CAP_ID="$(jq -r '.objectChanges[] | select(.type=="created" and (.objectType|test("MintCap"))) | .objectId' artifacts.managed_token.mint_cap.json | head -n1)"
echo "  MintCap ID: $MINT_CAP_ID"

echo "--- Deploying & initializing managed_token_pool ---"
MANAGED_TP_DIR="$ROOT_DIR/ccip/ccip_token_pools/managed_token_pool"
patch_move_toml "$MANAGED_TP_DIR/Move.toml" "ccip" "$CCIP_PKG_ID"
patch_move_toml "$MANAGED_TP_DIR/Move.toml" "mcms" "$MCMS_PKG_ID"
patch_move_toml "$MANAGED_TP_DIR/Move.toml" "mcms_owner" "$OWNER"
patch_move_toml "$MANAGED_TP_DIR/Move.toml" "managed_token" "$MANAGED_TOKEN_PKG_ID"
MANAGED_TP_PKG_ID="$(publish_and_pin "$MANAGED_TP_DIR" "managed_token_pool" "managed_tp")"

# Get the token pool administrator address (reusing active address)
TOKEN_POOL_ADMIN="$ACTIVE_ADDR"

MANAGED_TP_OWNER_CAP_ID="$(jq -r '.objectChanges[] | select(.type=="created" and (.objectType|test("OwnerCap"))) | .objectId' artifacts.managed_tp.publish.json | head -n1)"

# Initialize managed token pool with the managed token
sui client call --package "$MANAGED_TP_PKG_ID" --module managed_token_pool --function initialize_with_managed_token \
  --type-args "$USDC_COIN_T" \
  --args "$MANAGED_TP_OWNER_CAP_ID" "$CCIP_STATE_REF_ID" "$MANAGED_TOKEN_STATE_ID" "$MANAGED_TOKEN_OWNER_CAP_ID" "$USDC_METADATA_ID" "$MINT_CAP_ID" "$TOKEN_POOL_ADMIN" \
  --gas-budget "$GAS" --json | tee artifacts.managed_tp.init.json >/dev/null

MANAGED_TP_STATE_ID="$(jq -r '.objectChanges[] | select(.type=="created" and (.objectType|test("ManagedTokenPoolState"))) | .objectId' artifacts.managed_tp.init.json | head -n1)"

echo "  Managed Token Pool State: $MANAGED_TP_STATE_ID"
echo "  Managed Token Pool Owner Cap: $MANAGED_TP_OWNER_CAP_ID"

# Apply chain updates for chain 2
sui client call --package "$MANAGED_TP_PKG_ID" --module managed_token_pool --function apply_chain_updates \
  --type-args "$USDC_COIN_T" \
  --args "$MANAGED_TP_STATE_ID" "$MANAGED_TP_OWNER_CAP_ID" "[]" "[2]" "[[[24, 42, 24, 42]]]" "[[0,0,0,0,0,0,0,0,0,0,0,0,24,42,24,42,24,42,24,42,24,42,24,42,24,42,24,42,24,42,24,42]]" \
  --gas-budget "$GAS" --json | tee artifacts.managed_tp.apply_chains.json >/dev/null

# Set chain rate limiter config
sui client call --package "$MANAGED_TP_PKG_ID" --module managed_token_pool --function set_chain_rate_limiter_config \
  --type-args "$USDC_COIN_T" \
  --args "$MANAGED_TP_STATE_ID" "$MANAGED_TP_OWNER_CAP_ID" "$CLOCK_ID" "2" "false" "200000000000" "20000000000" "false" "200000000000" "20000000000" \
  --gas-budget "$GAS" --json | tee artifacts.managed_tp.rate_limiters.json >/dev/null

# echo "--- Minting USDC tokens for testing ---"
# sui client call --package "$MANAGED_TOKEN_PKG_ID" --module managed_token --function mint_and_transfer \
#   --type-args "$USDC_COIN_T" \
#   --args "$MANAGED_TOKEN_STATE_ID" "$MINT_CAP_ID" "$DENY_LIST_ID" "1000000000000000" "$ACTIVE_ADDR" \
#   --gas-budget "$GAS" --json | tee artifacts.usdc.mint.json >/dev/null

# USDC_COIN_ID="$(jq -r '.balanceChanges[] | select(.coinType | contains("mock_eth_token")) | .coinObjectId' artifacts.usdc.mint.json | head -n1)"
# echo "  USDC Coin ID: $USDC_COIN_ID"

echo "--- Debugging ---"
echo "CCIP_PKG_ID: $CCIP_PKG_ID"
echo "CCIP_STATE_REF_ID: $CCIP_STATE_REF_ID"
echo "CCIP_OWNER_CAP_ID: $CCIP_OWNER_CAP_ID"

# fee_quoter::apply_token_transfer_fee_config_updates for LINK, ETH, and USDC
echo "Applying token transfer fee config updates for LINK, ETH, and USDC..."
sui client call --package "$CCIP_PKG_ID" --module fee_quoter --function apply_token_transfer_fee_config_updates \
  --args "$CCIP_STATE_REF_ID" "$CCIP_OWNER_CAP_ID" 2 \
    "[\"$LINK_METADATA_ID\",\"$ETH_METADATA_ID\",\"$USDC_METADATA_ID\"]" "[50,50,50]" "[5000,5000,5000]" "[0,0,0]" "[180000,180000,180000]" "[640,640,640]" "[true,true,true]" "[]" \
  --gas-budget "$GAS" --json | tee artifacts.ccip.fee_quoter.token_transfer_fee_config.json >/dev/null

# fee_quoter::apply_premium_multiplier_wei_per_eth_updates for LINK, ETH, and USDC
echo "Applying premium multiplier updates for LINK, ETH, and USDC..."
sui client call --package "$CCIP_PKG_ID" --module fee_quoter --function apply_premium_multiplier_wei_per_eth_updates \
  --args "$CCIP_STATE_REF_ID" "$CCIP_OWNER_CAP_ID" \
    "[\"$LINK_METADATA_ID\",\"$ETH_METADATA_ID\",\"$USDC_METADATA_ID\"]" "[1,1,1]" \
  --gas-budget "$GAS" --json | tee artifacts.ccip.fee_quoter.premium_multiplier.json >/dev/null

# fee_quoter::apply_dest_chain_config_updates (no change needed)
echo "Applying destination chain config updates..."
sui client call --package "$CCIP_PKG_ID" --module fee_quoter --function apply_dest_chain_config_updates \
  --args "$CCIP_STATE_REF_ID" "$CCIP_OWNER_CAP_ID" 2 true 10 30000 300000 300000 16 40 3000 100 16 1 "[0x28,0x12,0xd5,0x2c]" false 25 90000 200000 1000000000000000 90000 10 \
  --gas-budget "$GAS" --json | tee artifacts.ccip.fee_quoter.dest_chain_config.json >/dev/null

# fee_quoter::update_prices_with_owner_cap to set the USD price for LINK
echo "Updating prices with owner cap for LINK..."
sui client call --package "$CCIP_PKG_ID" --module fee_quoter --function update_prices_with_owner_cap \
  --args "$CCIP_STATE_REF_ID" "$CCIP_OWNER_CAP_ID" "$CLOCK_ID" "[\"$LINK_METADATA_ID\"]" "[1000000000000000000]" "[2]" "[20]" \
  --gas-budget "$GAS" --json | tee artifacts.ccip.fee_quoter.update_prices_with_owner_cap.json >/dev/null


echo "Minting LINK coin (again) for fee token..."
sui client call \
  --package "$LINK_PKG_ID" \
  --module mock_link_token \
  --function mint \
  --args "$LINK_TREASURY_CAP_ID" "1000000000000000" \
  --gas-budget "$GAS" \
  --json | tee artifacts.link.fee_token.json >/dev/null

FEE_COIN_ID="$(jq -r '.objectChanges[] | select(.type=="created" and (.objectType|test("::coin::Coin<"))) | .objectId' artifacts.link.fee_token.json | head -n1)"


git checkout $ROOT_DIR

echo
echo "✅ Deployment complete. Artifacts written: artifacts.*.json"
echo "Packages:"
echo "  CCIP:                $CCIP_PKG_ID"
echo "  MCMS:                $MCMS_PKG_ID"
echo "  OnRamp:              $ONRAMP_PKG_ID"
echo "  LR Pool:             $LR_PKG_ID"
echo "  BM Pool:             $BM_PKG_ID"
echo "  Managed Token:       $MANAGED_TOKEN_PKG_ID"
echo "  Managed Token Pool:  $MANAGED_TP_PKG_ID"
echo ""
echo "Coin Types:"
echo "  LINK:                $LINK_COIN_T"
echo "  ETH:                 $ETH_COIN_T"
echo "  USDC:                $USDC_COIN_T"
echo ""
echo "Important State Objects:"
echo "  CCIP state:              $CCIP_STATE_REF_ID"
echo "  CCIP Owner Cap:          $CCIP_OWNER_CAP_ID"
echo "  OnRamp state:            $ONRAMP_STATE_ID"
echo "  LR Pool state:           $LR_STATE_ID"
echo "  BM Pool state:           $BM_STATE_ID"
echo "  Managed Token state:     $MANAGED_TOKEN_STATE_ID"
echo "  Managed Token Pool state: $MANAGED_TP_STATE_ID"
echo ""
echo "Treasury & Metadata:"
echo "  LINK Treasury Cap:       $LINK_TREASURY_CAP_ID"
echo "  LINK Metadata:           $LINK_METADATA_ID"
echo "  ETH Treasury Cap:        $ETH_TREASURY_CAP_ID"
echo "  ETH Metadata:            $ETH_METADATA_ID"
echo "  USDC Treasury (wrapped): $USDC_TREASURY_CAP_ID"
echo "  USDC Metadata:           $USDC_METADATA_ID"
echo "  MintCap (Managed Pool):  $MINT_CAP_ID"
echo ""
echo "Test Coins:"
echo "  LINK Coin:               $LINK_COIN_ID"
echo "  ETH Coin:                $ETH_COIN_ID"
# echo "  USDC Coin:               $USDC_COIN_ID"
echo "  FEE Coin (LINK):         $FEE_COIN_ID"
echo ""
echo "Token Pool Mapping:"
echo "  LINK   -> Lock/Release Token Pool"
echo "  ETH    -> Burn/Mint Token Pool"
echo "  USDC   -> Managed Token Pool"

# Generate .env.localnet file for ts-sdk-examples
echo ""
echo "--- Generating .env.localnet file ---"
ENV_FILE="$SCRIPT_DIR/../ts-sdk-examples/.env.localnet"

# Create ts-sdk-examples directory if it doesn't exist
mkdir -p "$(dirname "$ENV_FILE")"

# Check if .env.localnet already exists and has a private key
EXISTING_PRIVATE_KEY=""
if [[ -f "$ENV_FILE" ]]; then
  EXISTING_PRIVATE_KEY=$(grep "^SUI_PRIVATE_KEY=" "$ENV_FILE" 2>/dev/null | cut -d'=' -f2- | tr -d '"' || echo "")
fi

cat > "$ENV_FILE" << EOF
SUI_PRIVATE_KEY=${EXISTING_PRIVATE_KEY:-PLEASE_SET_YOUR_PRIVATE_KEY_HERE}

CCIP_PACKAGE_ID=$CCIP_PKG_ID
ONRAMP_PACKAGE_ID=$ONRAMP_PKG_ID
LR_POOL_PACKAGE_ID=$LR_PKG_ID
BM_POOL_PACKAGE_ID=$BM_PKG_ID

LINK_COIN_TYPE=$LINK_COIN_T
ETH_COIN_TYPE=$ETH_COIN_T

CCIP_STATE_ID=$CCIP_STATE_REF_ID
CCIP_OWNER_CAP_ID=$CCIP_OWNER_CAP_ID
ONRAMP_STATE_ID=$ONRAMP_STATE_ID

LR_POOL_STATE_ID=$LR_STATE_ID
BM_POOL_STATE_ID=$BM_STATE_ID

ETH_TREASURY_CAP_ID=$ETH_TREASURY_CAP_ID
LINK_TREASURY_CAP_ID=$LINK_TREASURY_CAP_ID
ETH_METADATA=$ETH_METADATA_ID
LINK_METADATA=$LINK_METADATA_ID
ETH_COIN_OBJECT=$ETH_COIN_ID
LINK_COIN_OBJECT=$LINK_COIN_ID

FEE_TOKEN_OBJECT=$FEE_COIN_ID
EOF

echo "✅ Generated $ENV_FILE"
if [[ -z "$EXISTING_PRIVATE_KEY" ]]; then
  echo ""
  echo "⚠️  Please set your SUI_PRIVATE_KEY in $ENV_FILE"
  echo "    You can export it with: sui keytool export --key-identity \$(sui client active-address)"
else
  echo "✅ Using existing private key from .env.localnet"
fi
echo ""
echo "You can now use the ts-sdk-examples with:"
echo "  cd $SCRIPT_DIR/../ts-sdk-examples"
echo "  bun ccip_send --dest-chain-selector 2 --receiver 0x... --pool-kind lock_release --network localnet"