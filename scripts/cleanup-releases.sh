#!/usr/bin/env bash
# Removes all GitHub releases except the latest one.
# Usage: [DRY_RUN=1] ./scripts/cleanup-releases.sh <owner/repo>
set -euo pipefail

REPO="${1:?Usage: $0 <owner/repo>}"

echo "=== Fetching releases for ${REPO} ==="

RELEASES=$(gh api "repos/${REPO}/releases?per_page=100" --jq '.[].id')
TOTAL=$(echo "$RELEASES" | wc -l | tr -d ' ')

if [ "$TOTAL" -le 1 ]; then
  echo "Only ${TOTAL} release(s) found — nothing to remove."
  exit 0
fi

LATEST_ID=$(gh api "repos/${REPO}/releases/latest" --jq '.id')

echo "Latest release ID: ${LATEST_ID}"
echo "Total releases: ${TOTAL}"

while IFS= read -r ID; do
  if [ "$ID" = "$LATEST_ID" ]; then
    echo "  KEEP   release ${ID} (latest)"
    continue
  fi

  TAG=$(gh api "repos/${REPO}/releases/${ID}" --jq '.tag_name')

  if [ "${DRY_RUN:-0}" = "1" ]; then
    echo "  DRY-RUN  would delete release ${ID} (${TAG})"
  else
    echo "  DELETE release ${ID} (${TAG})"
    gh api -X DELETE "repos/${REPO}/releases/${ID}"
  fi
done <<< "$RELEASES"

echo "=== Done ==="
