#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'EOF'
Usage:
  scripts/run-coolify-tests.sh [test_target]

Environment:
  RUNTIME            Container runtime: podman|docker (auto-detected if unset)
  COOLIFY_DIR        Coolify checkout path (default: /home/administrator/projects/coolify)
  COOLIFY_IMAGE      Container image used to run php/composer tests
                     (default: docker.io/serversideup/php:8.4-fpm-nginx-alpine)
  TEST_FILTER        Optional Laravel test filter passed as --filter=...
  DRY_RUN            Set to 1 to print command without executing

Examples:
  scripts/run-coolify-tests.sh tests/Unit/SshKeyFileSyncTest.php
  TEST_FILTER=unreadable scripts/run-coolify-tests.sh tests/Unit/SshKeyFileSyncTest.php
EOF
}

if [[ "${1:-}" == "-h" || "${1:-}" == "--help" ]]; then
  usage
  exit 0
fi

detect_runtime() {
  if [[ -n "${RUNTIME:-}" ]]; then
    printf '%s\n' "$RUNTIME"
    return
  fi
  # Prefer Docker for Coolify dev image builds because the upstream Dockerfiles
  # currently include short image names that may fail under stricter Podman
  # short-name resolution policies.
  if command -v docker >/dev/null 2>&1; then
    printf 'docker\n'
    return
  fi
  if command -v podman >/dev/null 2>&1; then
    printf 'podman\n'
    return
  fi
  echo "Error: neither podman nor docker is available." >&2
  exit 1
}

RUNTIME_BIN="$(detect_runtime)"
COOLIFY_DIR="${COOLIFY_DIR:-/home/administrator/projects/coolify}"
TEST_TARGET="${1:-}"
TEST_FILTER="${TEST_FILTER:-}"
DRY_RUN="${DRY_RUN:-0}"

if [[ ! -d "$COOLIFY_DIR" ]]; then
  echo "Error: COOLIFY_DIR does not exist: $COOLIFY_DIR" >&2
  exit 1
fi
if [[ ! -f "$COOLIFY_DIR/artisan" ]]; then
  echo "Error: artisan not found in $COOLIFY_DIR" >&2
  exit 1
fi
COOLIFY_IMAGE="${COOLIFY_IMAGE:-docker.io/serversideup/php:8.4-fpm-nginx-alpine}"
CACHE_ROOT="${HOME}/.cache/bountyos/coolify"
COMPOSER_CACHE="$CACHE_ROOT/composer"
mkdir -p "$COMPOSER_CACHE"

if [[ -n "$TEST_TARGET" ]]; then
  TEST_CMD="php artisan test --compact \"$TEST_TARGET\""
else
  TEST_CMD="php artisan test --compact"
fi
if [[ -n "$TEST_FILTER" ]]; then
  TEST_CMD="$TEST_CMD --filter=\"$TEST_FILTER\""
fi

INNER_CMD=$'if [ ! -f vendor/autoload.php ]; then\n'
INNER_CMD+=$'  composer install --no-interaction --prefer-dist --no-progress;\n'
INNER_CMD+=$'fi\n'
INNER_CMD+="$TEST_CMD"

CMD=(
  "$RUNTIME_BIN" run --rm -t
  --user root
  -v "$COOLIFY_DIR:/var/www/html"
  -v "$COMPOSER_CACHE:/tmp/composer-cache"
  -w /var/www/html
  -e COMPOSER_CACHE_DIR=/tmp/composer-cache
  "$COOLIFY_IMAGE"
  sh -lc "$INNER_CMD"
)

echo "==> runtime: $RUNTIME_BIN"
echo "==> coolify: $COOLIFY_DIR"
echo "==> image: $COOLIFY_IMAGE"
echo "==> test command: $TEST_CMD"

if [[ "$DRY_RUN" == "1" ]]; then
  echo "DRY_RUN=1; not executing."
  printf 'Command: %q\n' "${CMD[@]}"
  exit 0
fi

"${CMD[@]}"
