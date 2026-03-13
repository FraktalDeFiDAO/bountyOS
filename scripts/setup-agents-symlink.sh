#!/bin/bash
# setup-agents-symlink.sh
# Creates a convenient ~/agents symlink to the unified .qwen/agents/ directory
#
# Usage: ./scripts/setup-agents-symlink.sh

set -e

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
QWEN_AGENTS="$PROJECT_ROOT/.qwen/agents"
HOME_AGENTS_LINK="$HOME/agents"

echo "=== bountyOS Agents Symlink Setup ==="
echo ""
echo "This script creates a convenient symlink from:"
echo "  ~/agents  →  $QWEN_AGENTS"
echo ""
echo "This gives you easy access to all agents (project + global) from anywhere."
echo ""

# Check if target exists
if [ -e "$QWEN_AGENTS" ]; then
    echo "✓ Source directory exists: $QWEN_AGENTS"
else
    echo "✗ Source directory not found: $QWEN_AGENTS"
    echo "  Make sure you're running this from the bountyOS project."
    exit 1
fi

# Check if link already exists
if [ -L "$HOME_AGENTS_LINK" ]; then
    existing_target=$(readlink "$HOME_AGENTS_LINK")
    echo "! Symlink already exists: $HOME_AGENTS_LINK → $existing_target"
    read -p "  Overwrite? (y/N): " confirm
    if [[ "$confirm" != "y" && "$confirm" != "Y" ]]; then
        echo "  Skipping."
        exit 0
    fi
    rm "$HOME_AGENTS_LINK"
    echo "  Removed existing symlink."
elif [ -d "$HOME_AGENTS_LINK" ]; then
    echo "! Directory already exists: $HOME_AGENTS_LINK"
    echo "  Please remove or rename it first."
    exit 1
fi

# Create the symlink
ln -sf "$QWEN_AGENTS" "$HOME_AGENTS_LINK"
echo ""
echo "✓ Created symlink: ~/agents → $QWEN_AGENTS"
echo ""
echo "=== Available Agents ==="
ls -1 "$QWEN_AGENTS" | grep -v '\.yaml$' | grep -v '\.md$'
echo ""
echo "=== Quick Access ==="
echo "  cd ~/agents                    # Access all agents"
echo "  cat ~/agents/INDEX.md          # Onboarding guide"
echo "  cat ~/agents/workflow.yaml     # Lifecycle states"
echo "  cat ~/agents/agent-system-config.yaml  # Routing config"
echo ""
echo "=== Done ==="
