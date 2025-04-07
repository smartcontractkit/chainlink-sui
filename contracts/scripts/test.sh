#!/usr/bin/env bash
set -euxo pipefail

cd "$(dirname -- "$0")/.."

sui move test --path ccip/ccip
sui move test --path mcms/mcms
