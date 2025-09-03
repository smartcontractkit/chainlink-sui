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
#BM_COIN_TYPE="${BM_COIN_TYPE:-eth_token::ETH_TOKEN}"
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
  sui client publish --with-unpublished-dependencies --gas-budget "$GAS" --json \
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

# fee_quoter::initialize (uses LINK as fee token list)
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
#ONRAMP_STATE_ID="$(jq -r '.objectChanges[] | select(.type=="created" and (.objectType|test("onramp.*State";"i"))) | .objectId' artifacts.onramp.publish.json | head -n1)"
ONRAMP_STATE_ID="$(
  jq -r '
    .objectChanges[]
    | select(.type=="created" and (.objectType | test("::onramp::OnRampState$")))
    | .objectId
  ' artifacts.onramp.publish.json | head -n1
)"
ONRAMP_OWNER_CAP_ID="$(jq -r '.objectChanges[] | select(.type=="created" and (.objectType|test("OwnerCap"))) | .objectId' artifacts.onramp.publish.json | head -n1)"
[[ -n "$ONRAMP_STATE_ID" && -n "$ONRAMP_OWNER_CAP_ID" ]] || { echo "Missing OnRamp state/owner cap"; exit 1; }

sui client call --package "$ONRAMP_PKG_ID" --module onramp --function initialize \
  --args "$ONRAMP_STATE_ID" "$ONRAMP_OWNER_CAP_ID" "$NONCE_MANAGER_CAP_ID" \
        "$CCIP_SOURCE_TRANSFER_CAP_ID" "$LOCAL_CHAIN_SELECTOR" \
        "$FEE_AGGREGATOR_ADDR" "$ALLOWLIST_ADMIN_ADDR" \
        "$DEST_CHAIN_SELECTORS_JSON" "$DEST_CHAIN_ENABLED_JSON" "$DEST_CHAIN_ALLOWLIST_ENABLED_JSON" \
  --gas-budget "$GAS" --json | tee artifacts.onramp.init.json >/dev/null

echo "--- Deploying token_pool (base) ---"
TP_BASE_DIR="$ROOT_DIR/ccip/ccip_token_pools/token_pool"
# The guide suggests pointing router ref to a dummy 0x1 to simplify deployment.
patch_move_toml "$TP_BASE_DIR/Move.toml" "ccip" "$CCIP_PKG_ID"
patch_move_toml "$TP_BASE_DIR/Move.toml" "mcms" "$MCMS_PKG_ID"
patch_move_toml "$TP_BASE_DIR/Move.toml" "mcms_owner" "$OWNER"
TP_BASE_PKG_ID="$(publish_and_pin "$TP_BASE_DIR" "ccip_token_pool" "token_pool_base")"

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

# lock_release_token_pool::initialize(ccip_ref, LINK metadata, LINK treasury, package_id, rebalancer)
sui client call --package "$LR_PKG_ID" --module lock_release_token_pool --function initialize \
  --type-args "$LINK_COIN_T" \
  --args "$CCIP_STATE_REF_ID" "$LINK_METADATA_ID" "$LINK_TREASURY_CAP_ID" "$LR_PKG_ID" "$REBALANCER_ADDR" \
  --gas-budget "$GAS" --json | tee artifacts.lr_tp.init.json >/dev/null
LR_STATE_ID="$(jq -r '.objectChanges[] | select(.type=="created" and (.objectType|test("LockReleaseTokenPoolState"))) | .objectId' artifacts.lr_tp.init.json | head -n1)"
LR_OWNER_CAP_ID="$(jq -r '.objectChanges[] | select(.type=="created" and (.objectType|test("OwnerCap"))) | .objectId' artifacts.lr_tp.init.json | head -n1)"

# Optional: apply_chain_updates + rate limiter (example: add chain 2)
#sui client call --package "$LR_PKG_ID" --module lock_release_token_pool --function apply_chain_updates \
#  --args "$LR_STATE_ID" "$LR_OWNER_CAP_ID" "[]" "[2]" "[]" "[]" \
#  --gas-budget "$GAS" --json | tee artifacts.lr_tp.apply_chains.json >/dev/null

echo "--- Deploying mock ETH token for burn/mint pool (if not present) ---"
ETH_DIR="$ROOT_DIR/ccip/mock_eth_token"
ETH_PKG_KEY="mock_eth_token"
ETH_PKG_ID="$(publish_and_pin "$ETH_DIR" "$ETH_PKG_KEY" "eth")"
ETH_METADATA_ID="$(extract_created artifacts.eth.publish.json '::coin::CoinMetadata<' | head -n1)"
ETH_TREASURY_CAP_ID="$(extract_created artifacts.eth.publish.json '::coin::TreasuryCap<' | head -n1)"

# (Optional) mint some ETH so you can later test burns
#sui client call --package "$ETH_PKG_ID" --module eth_token --function mint \
#  --args "$ETH_TREASURY_CAP_ID" "1000000000000" \
#  --gas-budget "$GAS" --json | tee artifacts.eth.mint.json >/dev/null

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

# burn_mint_token_pool::initialize(ccip_ref, ETH metadata, ETH treasury cap, package_id) + <T>
sui client call --package "$BM_PKG_ID" --module burn_mint_token_pool --function initialize \
  --type-args "$ETH_COIN_T" \
  --args "$CCIP_STATE_REF_ID" "$ETH_METADATA_ID" "$ETH_TREASURY_CAP_ID" "$BM_PKG_ID" \
  --gas-budget "$GAS" --json | tee artifacts.bm_tp.init.json >/dev/null
BM_STATE_ID="$(jq -r '.objectChanges[] | select(.type=="created" and (.objectType|test("BurnMintTokenPoolState"))) | .objectId' artifacts.bm_tp.init.json | head -n1)"
BM_OWNER_CAP_ID="$(jq -r '.objectChanges[] | select(.type=="created" and (.objectType|test("OwnerCap"))) | .objectId' artifacts.bm_tp.init.json | head -n1)"

# Add chain 2; set basic rate limiters (example values)
#sui client call --package "$BM_PKG_ID" --module burn_mint_token_pool --function apply_chain_updates \
#  --type-args "$ETH_PKG_ID::$BM_COIN_TYPE" \
#  --args "$BM_STATE_ID" "$BM_OWNER_CAP_ID" "[]" "[2]" "[]" "[]" \
#  --gas-budget "$GAS" --json | tee artifacts.bm_tp.apply_chains.json >/dev/null

sui client call --package "$BM_PKG_ID" --module burn_mint_token_pool --function set_chain_rate_limiter_config \
  --type-args "$ETH_COIN_T" \
  --args "$BM_STATE_ID" "$BM_OWNER_CAP_ID" "$CLOCK_ID" "2" "false" "200000000000" "20000000000" "false" "200000000000" "20000000000" \
  --gas-budget "$GAS" --json | tee artifacts.bm_tp.rate_limiters.json >/dev/null

git checkout $ROOT_DIR

echo
echo "✅ Deployment complete. Artifacts written: artifacts.*.json"
echo "Packages:"
echo "  CCIP:           $CCIP_PKG_ID"
echo "  OnRamp:         $ONRAMP_PKG_ID"
echo "  LR Pool:        $LR_PKG_ID"
echo "  BM Pool:        $BM_PKG_ID"
echo "  LINK Coin Type: $LINK_COIN_T"
echo "  ETH Coin Type:  $ETH_COIN_T"
echo "Important objects:"
echo "  CCIP state:     $CCIP_STATE_REF_ID"
echo "  CCIP Owner Cap: $CCIP_OWNER_CAP_ID"
echo "  OnRamp state:   $ONRAMP_STATE_ID"
echo "  LR state: $LR_STATE_ID"
echo "  BM state: $BM_STATE_ID"
echo "  ETH Treasury Cap: $ETH_TREASURY_CAP_ID"
echo "  LINK Treasury Cap: $LINK_TREASURY_CAP_ID"
ECHO "  ETH Metadata: $ETH_METADATA_ID"
ECHO "  LINK Metadata: $LINK_METADATA_ID"
echo "Token Pool Support:"
echo "  ETH -> BM Token Pool"
echo "  LINK -> LR Token Pool"