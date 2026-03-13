#!/bin/bash
# sync_project_agents.sh
# Validates project_agents.tsv against actual project agent directories
# and syncs symlinks to .qwen/agents/
#
# Usage: ./scripts/sync_project_agents.sh [--dry-run]

set -e

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
REGISTRY_FILE="$PROJECT_ROOT/.agents/orchestrator/config/project_agents.tsv"
QWEN_AGENTS="$PROJECT_ROOT/.qwen/agents"
PROJECTS_DIR="$PROJECT_ROOT/au-workspace/projects"

DRY_RUN=false
if [[ "$1" == "--dry-run" ]]; then
    DRY_RUN=true
fi

echo "=== Project Agents Sync & Validation ==="
echo ""

# Check registry exists
if [[ ! -f "$REGISTRY_FILE" ]]; then
    echo "ERROR: Registry file not found: $REGISTRY_FILE"
    exit 1
fi

# Get list of actual project agent directories
echo "Scanning for project agent directories..."
ACTUAL_PROJECTS=()
for proj in "$PROJECTS_DIR"/bounty-*; do
    if [[ -d "$proj/agent" ]]; then
        ACTUAL_PROJECTS+=("$(basename "$proj")")
    fi
done

echo "Found ${#ACTUAL_PROJECTS[@]} project agent directories"
echo ""

# Parse registry to get projects
echo "Validating registry..."
REGISTRY_PROJECTS=()
REGISTRY_STATUS=()
while IFS=$'\t' read -r project bounty_id status priority owner agent_file role_dev role_content role_doc role_qa role_sched ops_schedule; do
    [[ "$project" == "project" ]] && continue  # Skip header
    REGISTRY_PROJECTS+=("$project")
    REGISTRY_STATUS+=("$status")
done < "$REGISTRY_FILE"

echo "Registry lists ${#REGISTRY_PROJECTS[@]} projects"
echo ""

# Find discrepancies
MISSING_IN_REGISTRY=()
EXTRA_IN_REGISTRY=()

for proj in "${ACTUAL_PROJECTS[@]}"; do
    found=false
    for reg_proj in "${REGISTRY_PROJECTS[@]}"; do
        if [[ "$proj" == "$reg_proj" ]]; then
            found=true
            break
        fi
    done
    if ! $found; then
        MISSING_IN_REGISTRY+=("$proj")
    fi
done

for reg_proj in "${REGISTRY_PROJECTS[@]}"; do
    found=false
    for proj in "${ACTUAL_PROJECTS[@]}"; do
        if [[ "$proj" == "$reg_proj" ]]; then
            found=true
            break
        fi
    done
    if ! $found; then
        EXTRA_IN_REGISTRY+=("$reg_proj")
    fi
done

# Report discrepancies
if [[ ${#MISSING_IN_REGISTRY[@]} -gt 0 ]]; then
    echo "⚠️  Projects MISSING from registry (in filesystem but not in TSV):"
    for p in "${MISSING_IN_REGISTRY[@]}"; do
        echo "    - $p"
    done
    echo ""
fi

if [[ ${#EXTRA_IN_REGISTRY[@]} -gt 0 ]]; then
    echo "⚠️  Projects EXTRA in registry (in TSV but not in filesystem):"
    for p in "${EXTRA_IN_REGISTRY[@]}"; do
        echo "    - $p"
    done
    echo ""
fi

if [[ ${#MISSING_IN_REGISTRY[@]} -eq 0 && ${#EXTRA_IN_REGISTRY[@]} -eq 0 ]]; then
    echo "✓ Registry matches filesystem"
    echo ""
fi

# Sync symlinks to .qwen/agents/
echo "Syncing project agent symlinks to .qwen/agents/..."

if $DRY_RUN; then
    echo "[DRY RUN - No changes will be made]"
    echo ""
fi

for proj in "${ACTUAL_PROJECTS[@]}"; do
    target_dir="$PROJECTS_DIR/$proj/agent"
    symlink_path="$QWEN_AGENTS/$proj"
    
    if [[ -L "$symlink_path" ]]; then
        existing_target=$(readlink "$symlink_path")
        if [[ "$existing_target" == "$target_dir" ]]; then
            continue  # Already linked correctly
        fi
        echo "  Updating symlink: $proj"
        if ! $DRY_RUN; then
            rm "$symlink_path"
            ln -s "$target_dir" "$symlink_path"
        fi
    elif [[ -e "$symlink_path" ]]; then
        echo "  ⚠️  Conflict: $proj exists as real file/directory"
    else
        echo "  Creating symlink: $proj"
        if ! $DRY_RUN; then
            ln -s "$target_dir" "$symlink_path"
        fi
    fi
done

echo ""
echo "=== Validation Complete ==="
echo ""
echo "Project agents available via ~/agents:"
ls -1 "$QWEN_AGENTS" | grep -E "^bounty-" | wc -l | xargs echo "  Total project agents:"
