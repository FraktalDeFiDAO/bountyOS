#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'EOF'
Usage:
  scripts/dev-check.sh [--dry-run]

Runs pre-submission checks across active bounty lanes.
No claim/PR submission should happen before this passes.

Environment overrides:
  COOLIFY_7724_DIR    (default: /home/administrator/projects/coolify)
  COOLIFY_7738_DIR    (default: /home/administrator/projects/coolify-7738)
  ZIO_9877_DIR        (default: /home/administrator/projects/zio-9877)
  ZIO_9878_DIR        (default: /home/administrator/projects/zio-9878)
EOF
}

DRY_RUN=0
if [[ "${1:-}" == "--dry-run" ]]; then
  DRY_RUN=1
elif [[ "${1:-}" == "-h" || "${1:-}" == "--help" ]]; then
  usage
  exit 0
elif [[ -n "${1:-}" ]]; then
  echo "Unknown arg: $1" >&2
  usage
  exit 2
fi

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
COOLIFY_7724_DIR="${COOLIFY_7724_DIR:-/home/administrator/projects/coolify}"
COOLIFY_7738_DIR="${COOLIFY_7738_DIR:-/home/administrator/projects/coolify-7738}"
ZIO_9877_DIR="${ZIO_9877_DIR:-/home/administrator/projects/zio-9877}"
ZIO_9878_DIR="${ZIO_9878_DIR:-/home/administrator/projects/zio-9878}"

detect_runtime() {
  if [[ -n "${RUNTIME:-}" ]]; then
    printf '%s\n' "$RUNTIME"
    return
  fi
  if command -v podman >/dev/null 2>&1; then
    printf 'podman\n'
    return
  fi
  if command -v docker >/dev/null 2>&1; then
    printf 'docker\n'
    return
  fi
  echo "Error: neither podman nor docker is available." >&2
  exit 1
}

RUNTIME_BIN="$(detect_runtime)"
if [[ -n "${COOLIFY_RUNTIME:-}" ]]; then
  COOLIFY_RUNTIME_BIN="$COOLIFY_RUNTIME"
elif command -v docker >/dev/null 2>&1; then
  COOLIFY_RUNTIME_BIN="docker"
else
  COOLIFY_RUNTIME_BIN="$RUNTIME_BIN"
fi

run() {
  echo
  echo "==> $*"
  if [[ "$DRY_RUN" == "1" ]]; then
    return 0
  fi
  "$@"
}

echo "Dev Check Runner"
echo "root: $ROOT"
echo "runtime: $RUNTIME_BIN"
echo "coolify runtime: $COOLIFY_RUNTIME_BIN"
if [[ "$DRY_RUN" == "1" ]]; then
  echo "mode: dry-run"
fi

run "$RUNTIME_BIN" compose -f "$ROOT/docker-compose.dev.yml" config
run "$RUNTIME_BIN" compose -f "$ROOT/docker-compose.yml" config
run "$RUNTIME_BIN" compose -f "$ROOT/docker-compose.prod.ssl.yml" config

run env RUNTIME="$COOLIFY_RUNTIME_BIN" "$ROOT/scripts/run-coolify-tests.sh" tests/Unit/SshKeyFileSyncTest.php
run env RUNTIME="$COOLIFY_RUNTIME_BIN" COOLIFY_DIR="$COOLIFY_7738_DIR" "$ROOT/scripts/run-coolify-tests.sh" tests/Unit/ApplicationDeploymentServerEnvVariablesTest.php

run env RUNTIME="$RUNTIME_BIN" ZIO_DIR="$ZIO_9877_DIR" "$ROOT/scripts/run-zio-tests.sh" "coreTestsJVM/testOnly zio.PromiseSpec -- -t become"
run env RUNTIME="$RUNTIME_BIN" ZIO_DIR="$ZIO_9878_DIR" "$ROOT/scripts/run-zio-tests.sh" "coreJVM/compile"

echo
echo "All dev checks completed."
