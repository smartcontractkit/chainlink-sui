#!/usr/bin/env bash
#
# Validates that all intra-repo pseudo-version commits are reachable from HEAD.
# This prevents downstream consumers from hitting "unknown revision" errors
# after squash-merges rewrite commit history.
#
# Usage: ./scripts/check-module-versions.sh

set -euo pipefail

REPO_MODULE="github.com/smartcontractkit/chainlink-sui"
ERRORS=0

while IFS= read -r gomod; do
  # Extract pseudo-versions for intra-repo modules (format: v0.0.0-YYYYMMDDHHMMSS-<commit12>)
  # Skip the Go placeholder version used for fully-replaced modules
  while IFS= read -r line; do
    [ -z "$line" ] && continue

    module=$(echo "$line" | awk '{print $1}')
    version=$(echo "$line" | awk '{print $2}')

    # Skip placeholder versions that exist solely because of replace directives
    if [[ "$version" == "v0.0.0-00010101000000-000000000000" || "$version" == "v0.0.0" ]]; then
      continue
    fi

    # Extract 12-char commit hash from pseudo-version
    commit="${version##*-}"

    if ! git merge-base --is-ancestor "$commit" HEAD 2>/dev/null; then
      echo "ERROR: $gomod requires $module@$version"
      echo "       commit $commit is NOT reachable from HEAD (develop)"
      echo "       Fix: run 'go get $module@develop' or tag a release"
      echo ""
      ERRORS=$((ERRORS + 1))
    fi
  done < <(grep -E "^\s+${REPO_MODULE}(/\S+)?\s+v0\.0\.0-[0-9]+-[0-9a-f]{12}" "$gomod" | sed 's|//.*||')
done < <(find . -name "go.mod" -not -path "./.gomodcache/*" -not -path "./.cursor-sandbox-home/*" -not -path "./.move-home/*")

if [ "$ERRORS" -gt 0 ]; then
  echo "FAILED: $ERRORS intra-repo pseudo-version(s) reference unreachable commits."
  echo ""
  echo "This typically happens when a feature branch is squash-merged,"
  echo "creating a new commit hash. The old pseudo-version becomes stale."
  echo ""
  echo "To fix, run from the repo root:"
  echo "  go tool gomods tidy"
  echo "  # or manually: go get <module>@<commit-on-develop>"
  exit 1
fi

echo "OK: all intra-repo pseudo-versions reference reachable commits."
