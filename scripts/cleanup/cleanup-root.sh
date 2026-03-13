#!/bin/bash
# bountyOS Root Directory Cleanup Script
# 
# This script moves files from the root directory to their proper locations
# according to the project organization standards.
#
# Usage: ./scripts/cleanup/cleanup-root.sh [--dry-run]

set -e  # Exit on error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
DRY_RUN=false

# Parse arguments
if [[ "${1:-}" == "--dry-run" ]]; then
    DRY_RUN=true
    echo -e "${YELLOW}🔍 DRY RUN MODE - No files will be moved${NC}"
    echo ""
fi

# Counters
MOVED=0
SKIPPED=0
ERRORS=0

# Function to move files
move_file() {
    local src="$1"
    local dest="$2"
    local category="$3"
    
    if [[ ! -f "$src" ]]; then
        return 0
    fi
    
    if [[ "$DRY_RUN" == true ]]; then
        echo -e "${BLUE}[DRY]${NC} Would move: $src → $dest"
        MOVED=$((MOVED + 1))
    else
        if mv "$src" "$dest" 2>/dev/null; then
            echo -e "${GREEN}✓${NC} [$category] $src → $dest"
            MOVED=$((MOVED + 1))
        else
            echo -e "${RED}✗${NC} [$category] Failed to move: $src"
            ERRORS=$((ERRORS + 1))
        fi
    fi
    return 0
}

# Function to create directory if needed
ensure_dir() {
    local dir="$1"
    if [[ ! -d "$dir" ]]; then
        if [[ "$DRY_RUN" == false ]]; then
            mkdir -p "$dir"
        fi
    fi
}

echo "========================================"
echo "  bountyOS Root Directory Cleanup"
echo "========================================"
echo ""
echo "Source: $ROOT_DIR"
echo ""

cd "$ROOT_DIR"

# Create ALL target directories upfront
echo -e "${BLUE}Creating directory structure...${NC}"
mkdir -p "$ROOT_DIR/docs/status"
mkdir -p "$ROOT_DIR/docs/sessions"
mkdir -p "$ROOT_DIR/docs/archive"
mkdir -p "$ROOT_DIR/docs/archive/investigations"
mkdir -p "$ROOT_DIR/docs/archive/bounty-code"
mkdir -p "$ROOT_DIR/docs/archive/coolify"
mkdir -p "$ROOT_DIR/docs/guides"
mkdir -p "$ROOT_DIR/docs/architecture"
mkdir -p "$ROOT_DIR/docs/research"
mkdir -p "$ROOT_DIR/scripts"
mkdir -p "$ROOT_DIR/config"
echo -e "${GREEN}✓${NC} Directory structure ready"
echo ""

# ============================================
# Category 1: Status Trackers → docs/status/
# ============================================
echo -e "${YELLOW}📊 Moving Status Trackers...${NC}"

STATUS_FILES=(
    "ALL_BOUNTIES_SUBMISSION_STATUS.md"
    "500_TODAY_STATUS_1630.md"
    "BEACON_AUTO_SUBMITTER_STATUS.md"
    "BEACON_RELAY_STATUS.md"
    "BOUNTY_PIPELINE_STATUS.md"
    "BOUNTY_PORTFOLIO_STATUS.md"
    "BOUNTY_SUBMISSION_STATUS.md"
    "BOUNTY_SUBMISSIONS_COMPLETED.md"
    "IN_PROGRESS_BOUNTIES_STATUS.md"
    "SUBMITTED_BOUNTIES_STATUS.md"
    "SCANNER_MIGRATION_PROGRESS.md"
)

for file in "${STATUS_FILES[@]}"; do
    if [[ -f "$file" ]]; then
        move_file "$file" "docs/status/$file" "STATUS"
    fi
done

echo ""

# ============================================
# Category 2: Session Summaries → docs/sessions/
# ============================================
echo -e "${YELLOW}📝 Moving Session Summaries...${NC}"

SESSION_FILES=(
    "500_TODAY_TRACKER.md"
    "BOUNTY_AGENT_IMPLEMENTATION_SUMMARY.md"
    "BOUNTY_EXECUTION_SUMMARY.md"
    "COMPLETE_INTEGRATION_SUMMARY.md"
    "FINAL_SESSION_SUMMARY.md"
    "REFACTORING_COMPLETION_REPORT.md"
    "REFACTORING_SUMMARY.md"
    "SESSION_SUMMARY_MARCH_10.md"
    "SESSION_SUMMARY_MARCH_11.md"
    "SESSION_WRAPUP.md"
)

for file in "${SESSION_FILES[@]}"; do
    if [[ -f "$file" ]]; then
        move_file "$file" "docs/sessions/$file" "SESSION"
    fi
done

echo ""

# ============================================
# Category 3: Guides → docs/guides/
# ============================================
echo -e "${YELLOW}📚 Moving Guides...${NC}"

GUIDE_FILES=(
    "QUICK_COMMENT_GUIDE.md"
    "QUICK_WIN_BOUNTIES.md"
    "VERIFICATION_QUICK_START.md"
    "BOUNTY_VERIFICATION_PROTOCOL.md"
)

for file in "${GUIDE_FILES[@]}"; do
    if [[ -f "$file" ]]; then
        move_file "$file" "docs/guides/$file" "GUIDE"
    fi
done

echo ""

# ============================================
# Category 4: Python Scripts → scripts/
# ============================================
echo -e "${YELLOW}🐍 Moving Python Scripts...${NC}"

PYTHON_SCRIPTS=(
    "app.py"
    "beacon_auto_submit.py"
    "beacon_heartbeat.py"
    "beacon_register.py"
    "create_ridima_video.py"
)

for file in "${PYTHON_SCRIPTS[@]}"; do
    if [[ -f "$file" ]]; then
        move_file "$file" "scripts/$file" "PYTHON"
    fi
done

echo ""

# ============================================
# Category 5: TypeScript Files → appropriate locations
# ============================================
echo -e "${YELLOW}📘 Moving TypeScript Files...${NC}"

# Move standalone TS files to scripts or appropriate project
if [[ -f "approvals.ts" ]]; then
    move_file "approvals.ts" "scripts/approvals.ts" "TYPESCRIPT"
fi

if [[ -f "index.ts" ]]; then
    # Keep index.ts in root if it's a main entry point
    echo -e "${BLUE}ℹ${NC} [TYPESCRIPT] Keeping index.ts in root (entry point)"
    ((SKIPPED++))
fi

echo ""

# ============================================
# Category 6: Archive Old/Completed Items
# ============================================
echo -e "${YELLOW}🗄️  Moving to Archive...${NC}"

ARCHIVE_FILES=(
    "AUDIT_MITIGATION_PLAN.md"
    "BOUNTY_OVERDRIVE_ACTIVATION.md"
    "COMPLETION_PLAN.md"
    "CORRECTED_TIMELINE.md"
    "CRITICAL_ACTIONS_TODAY.md"
    "MUSHAF_QA_EXECUTION_PLAN.md"
    "OPERATION_500_PLUS.md"
    "PLATFORM_EXPANSION_PLAN.md"
    "POST_COOLIFY_QUALITY_GATES.md"
    "SCANNER_MIGRATION_PLAN.md"
    "VERIFICATION_PROTOCOL_IMPLEMENTATION.md"
)

for file in "${ARCHIVE_FILES[@]}"; do
    if [[ -f "$file" ]]; then
        move_file "$file" "docs/archive/$file" "ARCHIVE"
    fi
done

echo ""

# ============================================
# Category 7: Investigation Reports
# ============================================
echo -e "${YELLOW}🔍 Moving Investigation Reports...${NC}"

INVESTIGATION_FILES=(
    "FABRICATED_BOUNTY_INVESTIGATION.md"
    "BOUNTY_VERIFICATION_AUDIT_2026-03-11.md"
    "BOUNTY_SUBMISSION_AUDIT_2026_03_11.md"
)

for file in "${INVESTIGATION_FILES[@]}"; do
    if [[ -f "$file" ]]; then
        move_file "$file" "docs/archive/investigations/$file" "INVESTIGATION"
    fi
done

echo ""

# ============================================
# Category 8: Coolify Comments
# ============================================
echo -e "${YELLOW}💬 Moving Coolify Comments...${NC}"

if [[ -f "coolify_7724_comment.md" ]]; then
    move_file "coolify_7724_comment.md" "docs/archive/coolify/coolify_7724_comment.md" "COOLIFY"
fi

if [[ -f "coolify_8779_comment.md" ]]; then
    move_file "coolify_8779_comment.md" "docs/archive/coolify/coolify_8779_comment.md" "COOLIFY"
fi

echo ""

# ============================================
# Category 9: Bounty Patches/Code
# ============================================
echo -e "${YELLOW}🔧 Moving Bounty Code...${NC}"

if [[ -f "bounty5_87.patch" ]]; then
    move_file "bounty5_87.patch" "docs/archive/bounty-code/bounty5_87.patch" "PATCH"
fi

if [[ -f "bounty5_87.ts" ]]; then
    move_file "bounty5_87.ts" "docs/archive/bounty-code/bounty5_87.ts" "TYPESCRIPT"
fi

if [[ -f "bounty9_74.ts" ]]; then
    move_file "bounty9_74.ts" "docs/archive/bounty-code/bounty9_74.ts" "TYPESCRIPT"
fi

if [[ -f "kadeshx_audit.rs" ]]; then
    move_file "kadeshx_audit.rs" "docs/archive/bounty-code/kadeshx_audit.rs" "RUST"
fi

echo ""

# ============================================
# Category 10: Miscellaneous
# ============================================
echo -e "${YELLOW}📦 Moving Miscellaneous...${NC}"

if [[ -f "PROJECT_ISOLATION_AND_HUMAN_AUTHENTICITY.md" ]]; then
    move_file "PROJECT_ISOLATION_AND_HUMAN_AUTHENTICITY.md" "docs/architecture/PROJECT_ISOLATION_AND_HUMAN_AUTHENTICITY.md" "ARCHITECTURE"
fi

if [[ -f "BOUNTYOS_ARCHITECTURE_V1.md" ]]; then
    move_file "BOUNTYOS_ARCHITECTURE_V1.md" "docs/architecture/BOUNTYOS_ARCHITECTURE_V1.md" "ARCHITECTURE"
fi

if [[ -f "NEW_PLATFORMS_SUMMARY.md" ]]; then
    move_file "NEW_PLATFORMS_SUMMARY.md" "docs/research/NEW_PLATFORMS_SUMMARY.md" "RESEARCH"
fi

if [[ -f "NEXT_BOUNTIES.md" ]]; then
    move_file "NEXT_BOUNTIES.md" "docs/status/NEXT_BOUNTIES.md" "STATUS"
fi

if [[ -f "PI_LLM_SETUP.md" ]]; then
    move_file "PI_LLM_SETUP.md" "docs/guides/PI_LLM_SETUP.md" "GUIDE"
fi

if [[ -f "ralph-bridge-pinger.md" ]]; then
    move_file "ralph-bridge-pinger.md" "docs/archive/ralph-bridge-pinger.md" "ARCHIVE"
fi

echo ""

# ============================================
# Summary
# ============================================
echo "========================================"
echo "  Cleanup Summary"
echo "========================================"
echo -e "Files Moved:    ${GREEN}$MOVED${NC}"
echo -e "Files Skipped:  ${YELLOW}$SKIPPED${NC}"
echo -e "Errors:         ${RED}$ERRORS${NC}"
echo ""

if [[ "$DRY_RUN" == true ]]; then
    echo -e "${YELLOW}This was a DRY RUN. No files were actually moved.${NC}"
    echo "Run without --dry-run to execute the cleanup."
else
    echo -e "${GREEN}Cleanup complete!${NC}"
fi

echo ""
echo "Next steps:"
echo "1. Review moved files in their new locations"
echo "2. Update any broken links"
echo "3. Run 'git status' to see changes"
echo "4. Commit with message: 'chore: reorganize root directory structure'"
