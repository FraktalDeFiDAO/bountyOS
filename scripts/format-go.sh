#!/usr/bin/env bash
set -euo pipefail

RUNTIME=${RUNTIME:-podman}
COMPOSE_FILE=${COMPOSE_FILE:-docker-compose.dev.yml}

$RUNTIME compose -f "$COMPOSE_FILE" run --rm obsidian sh -c \
  "find ./cmd ./internal -type f -name '*.go' -print0 | xargs -0 goimports -w"
