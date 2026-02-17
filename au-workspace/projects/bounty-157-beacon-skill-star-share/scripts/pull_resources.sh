#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
RAW_GITHUB_DIR="$ROOT_DIR/resources/raw/github"
RAW_LOCAL_DIR="$ROOT_DIR/resources/raw/local"

mkdir -p "$RAW_GITHUB_DIR" "$RAW_LOCAL_DIR"

STAMP_UTC="$(date -u +"%Y-%m-%dT%H:%M:%SZ")"

fetch() {
  local url="$1"
  local out="$2"
  curl --fail --silent --show-error --location "$url" -o "$out"
}

# GitHub resources
fetch "https://api.github.com/repos/Scottcjn/rustchain-bounties/issues/157" "$RAW_GITHUB_DIR/issue_157.json"
fetch "https://api.github.com/repos/Scottcjn/rustchain-bounties/issues/157/comments" "$RAW_GITHUB_DIR/issue_157_comments.json"
fetch "https://api.github.com/repos/Scottcjn/beacon-skill" "$RAW_GITHUB_DIR/beacon_skill_repo.json"
fetch "https://raw.githubusercontent.com/Scottcjn/beacon-skill/main/README.md" "$RAW_GITHUB_DIR/beacon_skill_README.md"

# Local feed resources (best effort)
if curl --silent --show-error --location --max-time 10 "http://localhost:12496/api/bounties" -o "$RAW_LOCAL_DIR/feed_bounties.json"; then
  :
else
  echo "warning: could not pull http://localhost:12496/api/bounties" >&2
fi

# Project summary for quick review
{
  echo "# Resource Snapshot"
  echo
  echo "Synced at (UTC): $STAMP_UTC"
  echo
  echo "## Bounty Issue"
  jq -r '"- title: " + .title, "- state: " + .state, "- html_url: " + .html_url, "- created_at: " + .created_at, "- updated_at: " + .updated_at' "$RAW_GITHUB_DIR/issue_157.json"
  echo
  echo "## Target Repo"
  jq -r '"- full_name: " + .full_name, "- html_url: " + .html_url, "- stargazers_count: " + (.stargazers_count|tostring), "- forks_count: " + (.forks_count|tostring), "- open_issues_count: " + (.open_issues_count|tostring)' "$RAW_GITHUB_DIR/beacon_skill_repo.json"
  echo
  if [[ -f "$RAW_LOCAL_DIR/feed_bounties.json" ]]; then
    echo "## Local Feed Hit"
    jq -r '[.[] | select(.url == "https://github.com/Scottcjn/rustchain-bounties/issues/157")] | if length == 0 then "- not present in local feed snapshot" else .[0] | "- score: " + (.score|tostring), "- platform: " + .platform, "- payment_type: " + .payment_type end' "$RAW_LOCAL_DIR/feed_bounties.json"
  fi
} > "$ROOT_DIR/resources/summary.md"

echo "synced resources to: $ROOT_DIR/resources/raw"
echo "wrote summary: $ROOT_DIR/resources/summary.md"
