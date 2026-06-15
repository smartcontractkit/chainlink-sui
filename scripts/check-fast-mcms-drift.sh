#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SLOW="${REPO_ROOT}/contracts/mcms/mcms"
FAST="${REPO_ROOT}/contracts/mcms/fast_mcms"

echo "Checking sources/ lockstep between mcms and fast_mcms..."
if ! diff -r "${SLOW}/sources" "${FAST}/sources" > /dev/null 2>&1; then
    echo "ERROR: sources/ differ between mcms and fast_mcms:"
    diff -r "${SLOW}/sources" "${FAST}/sources" || true
    echo ""
    echo "FAILED: fast_mcms sources have drifted from mcms."
    echo "Only Move.toml and Move.lock are allowed to differ."
    exit 1
fi

echo "OK: mcms and fast_mcms sources are byte-identical."
