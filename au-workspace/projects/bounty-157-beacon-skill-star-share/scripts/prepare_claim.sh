#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
POST_URL=""
WALLET=""
PROOF_NOTE="Screenshot attached in this comment"
OUT_FILE="$ROOT_DIR/output/claim_comment_filled.md"

usage() {
  cat <<'USAGE'
Usage:
  ./scripts/prepare_claim.sh --post-url <url> --wallet <wallet> [--proof-note <text>] [--out <file>]

Example:
  ./scripts/prepare_claim.sh \
    --post-url "https://x.com/you/status/123" \
    --wallet "RTC123..." \
    --proof-note "Screenshot attached below"
USAGE
}

while (( $# > 0 )); do
  case "$1" in
    --post-url)
      POST_URL="${2:-}"
      shift 2
      ;;
    --wallet)
      WALLET="${2:-}"
      shift 2
      ;;
    --proof-note)
      PROOF_NOTE="${2:-}"
      shift 2
      ;;
    --out)
      OUT_FILE="${2:-}"
      shift 2
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "error: unknown argument: $1" >&2
      usage
      exit 1
      ;;
  esac
done

if [[ -z "$POST_URL" || -z "$WALLET" ]]; then
  echo "error: --post-url and --wallet are required" >&2
  usage
  exit 1
fi

mkdir -p "$(dirname "$OUT_FILE")"

cat > "$OUT_FILE" <<EOF2
Claiming completion for this bounty.

- GitHub star proof: ${PROOF_NOTE}
- Public share link: ${POST_URL}
- Wallet address (RTC): ${WALLET}

I starred beacon-skill and shared a public post that includes the repo link and a short description.
EOF2

echo "wrote: $OUT_FILE"
