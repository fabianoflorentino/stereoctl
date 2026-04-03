#!/usr/bin/env bash
# Removes all git tags except the latest semver tag.
# Usage: [DRY_RUN=1] ./scripts/cleanup-tags.sh <owner/repo>
set -euo pipefail

REPO="${1:?Usage: $0 <owner/repo>}"

echo "=== Fetching tags for ${REPO} ==="

LATEST=$(git tag --sort=version:refname | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+$' | tail -1)

if [ -z "$LATEST" ]; then
  echo "No semver tags found — nothing to remove."
  exit 0
fi

echo "Latest tag: ${LATEST}"

git tag --sort=version:refname | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+$' | while IFS= read -r TAG; do
  if [ "$TAG" = "$LATEST" ]; then
    echo "  KEEP   ${TAG} (latest)"
    continue
  fi

  if [ "${DRY_RUN:-0}" = "1" ]; then
    echo "  DRY-RUN  would delete tag ${TAG}"
  else
    echo "  DELETE ${TAG}"
    gh api -X DELETE "repos/${REPO}/git/refs/tags/${TAG}"
  fi
done

echo "=== Done ==="
