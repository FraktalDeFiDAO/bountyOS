#!/usr/bin/env bash
# Obsidian ARM64 Forensic Reproduction Utility (V3 - Container Native)
# Usage: ./scripts/repro_arm64_deadlock.sh [--gui] [--no-pull]

set -e

# --- Configuration ---
REPO_URL="https://github.com/projectdiscovery/tlsx"
PROJ_DIR="temp_repro_tlsx"
ARCH="arm64"
VM_DIR="tools/vm"
POC_SOURCE="au-workspace/projects/bounty-tlsx-hangs/source/internal/pdcp/race_test.go"

# --- Args ---
GUI_MODE=false
PULL_LATEST=true

for arg in "$@"; do
  case $arg in
    --gui) GUI_MODE=true ;;
    --no-pull) PULL_LATEST=false ;;
  esac
done

# --- Container Engine Detection ---
if command -v podman &> /dev/null; then
    DOCKER="podman"
else
    DOCKER="docker"
fi
echo "🐳 Using Container Engine: $DOCKER"

echo "🔍 Phase 1: Source Acquisition"
if [ "$PULL_LATEST" = true ]; then
    echo "Cloning absolute current version from GitHub..."
    rm -rf "$PROJ_DIR"
    git clone --depth 1 "$REPO_URL" "$PROJ_DIR"
else
    echo "Using existing local source..."
fi

echo "🧪 Phase 2: Exploit Injection"
mkdir -p "$PROJ_DIR/internal/pdcp"
cp "$POC_SOURCE" "$PROJ_DIR/internal/pdcp/race_test.go"

echo "🏗️ Phase 3: Cross-Compilation (AMD64 -> ARM64)"
cd "$PROJ_DIR"
GOARCH=$ARCH go test -c ./internal/pdcp -o ../$VM_DIR/pdcp_arm64.test
cd ..

echo "📦 Phase 4: Initrd Packaging"
$DOCKER run --rm --platform linux/amd64 \
  -v "$(pwd):/work" -w /work/$VM_DIR \
  docker.io/library/alpine:latest sh -c "
    apk add --no-cache cpio gzip > /dev/null
    rm -rf initramfs_contents && mkdir -p initramfs_contents && cd initramfs_contents
    zcat ../initramfs | cpio -idmv > /dev/null
    cp /work/$VM_DIR/pdcp_arm64.test .
    printf '#!/bin/sh\necho \"\n--- BOOTING ARM64 FORENSIC ENVIRONMENT ---\"\n./pdcp_arm64.test -test.v -test.run TestUploadWriterExploit\npoweroff -f\n' > init
    chmod +x init
    find . | cpio -H newc -o | gzip > ../initramfs_repro.gz
"

echo "🚀 Phase 5: AArch64 System Emulation"
QEMU_ARGS=(
    "-machine" "virt"
    "-cpu" "cortex-a57"
    "-m" "512M"
    "-kernel" "/vm/vmlinuz"
    "-initrd" "/vm/initramfs_repro.gz"
    "-serial" "mon:stdio"
)

if [ "$GUI_MODE" = true ]; then
    echo "Launching in GRAPHICAL mode..."
    QEMU_ARGS+=("-device" "virtio-gpu-pci" "-display" "gtk")
else
    QEMU_ARGS+=("-nographic" "-append" "console=ttyAMA0 quiet")
fi

# --init handles Ctrl+C; -serial mon:stdio allows Ctrl+A, X to force quit
$DOCKER run --rm -it --init \
  -v "$(pwd)/tools/vm:/vm" \
  -w /vm \
  docker.io/tianon/qemu qemu-system-aarch64 "${QEMU_ARGS[@]}"
