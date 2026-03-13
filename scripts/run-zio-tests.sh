#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'EOF'
Usage:
  scripts/run-zio-tests.sh [sbt_task]

Environment:
  RUNTIME         Container runtime: podman|docker (auto-detected if unset)
  ZIO_DIR         ZIO checkout path (default: /home/administrator/projects/zio)
  SBT_IMAGE       SBT container image
                  (default: docker.io/sbtscala/scala-sbt:eclipse-temurin-21.0.5_1.10.5_3.6.2)
  DRY_RUN         Set to 1 to print command without executing

Examples:
  scripts/run-zio-tests.sh "coreTestsJVM/testOnly zio.PromiseSpec"
  ZIO_DIR=/home/administrator/projects/zio-9878 scripts/run-zio-tests.sh "coreJVM/compile"
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
ZIO_DIR="${ZIO_DIR:-/home/administrator/projects/zio}"
SBT_TASK="${1:-coreTestsJVM/testOnly zio.PromiseSpec}"
SBT_IMAGE="${SBT_IMAGE:-docker.io/sbtscala/scala-sbt:eclipse-temurin-21.0.5_1.10.5_3.6.2}"
DRY_RUN="${DRY_RUN:-0}"

if [[ ! -d "$ZIO_DIR" ]]; then
  echo "Error: ZIO_DIR does not exist: $ZIO_DIR" >&2
  exit 1
fi
if [[ ! -f "$ZIO_DIR/build.sbt" ]]; then
  echo "Error: build.sbt not found in $ZIO_DIR" >&2
  exit 1
fi

CACHE_ROOT="${HOME}/.cache/bountyos/zio"
IVY_CACHE="$CACHE_ROOT/ivy2"
SBT_CACHE="$CACHE_ROOT/sbt"
COURSIER_CACHE="$CACHE_ROOT/coursier"
mkdir -p "$IVY_CACHE" "$SBT_CACHE" "$COURSIER_CACHE"

CMD=(
  "$RUNTIME_BIN" run --rm -t
  -v "$ZIO_DIR:/workspace"
  -v "$IVY_CACHE:/root/.ivy2"
  -v "$SBT_CACHE:/root/.sbt"
  -v "$COURSIER_CACHE:/root/.cache/coursier"
  -w /workspace
  -e SBT_OPTS="-Dsbt.supershell=false -Dsbt.ci=true"
  "$SBT_IMAGE"
  sbt --batch "$SBT_TASK"
)

echo "==> runtime: $RUNTIME_BIN"
echo "==> zio dir: $ZIO_DIR"
echo "==> sbt image: $SBT_IMAGE"
echo "==> sbt task: $SBT_TASK"

if [[ "$DRY_RUN" == "1" ]]; then
  echo "DRY_RUN=1; not executing."
  printf 'Command: %q\n' "${CMD[@]}"
  exit 0
fi

"${CMD[@]}"
