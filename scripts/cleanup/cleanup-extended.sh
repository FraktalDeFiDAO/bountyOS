#!/bin/bash
# bountyOS Extended Cleanup Script - Phase 2
# 
# Handles remaining root directory files not covered by initial cleanup
#
# Usage: ./scripts/cleanup/cleanup-extended.sh [--dry-run]

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
DRY_RUN=false

if [[ "${1:-}" == "--dry-run" ]]; then
    DRY_RUN=true
    echo -e "${YELLOW}🔍 DRY RUN MODE${NC}"
fi

MOVED=0
SKIPPED=0

move_file() {
    local src="$1"
    local dest="$2"
    local category="$3"
    
    if [[ ! -f "$src" ]]; then
        return 0
    fi
    
    if [[ "$DRY_RUN" == true ]]; then
        echo -e "${BLUE}[DRY]${NC} $src → $dest"
        MOVED=$((MOVED + 1))
    else
        mkdir -p "$(dirname "$dest")"
        if mv "$src" "$dest" 2>/dev/null; then
            echo -e "${GREEN}✓${NC} [$category] $src → $dest"
            MOVED=$((MOVED + 1))
        else
            echo -e "${RED}✗${NC} [$category] Failed: $src"
        fi
    fi
    return 0
}

cd "$ROOT_DIR"

echo "========================================"
echo "  bountyOS Extended Cleanup (Phase 2)"
echo "========================================"
echo ""

# Create additional directories
mkdir -p docs/archive/patches
mkdir -p docs/archive/code-snippets
mkdir -p docs/archive/followups
mkdir -p docs/investigations
mkdir -p docs/plans

# ============================================
# Category 1: Audit Reports → docs/archive/
# ============================================
echo -e "${YELLOW}📋 Moving Audit Reports...${NC}"

move_file "AUDIT_MITIGATION_PLAN.md" "docs/archive/AUDIT_MITIGATION_PLAN.md" "AUDIT"
move_file "BOUNTY_AUDIT_REPORT.md" "docs/archive/BOUNTY_AUDIT_REPORT.md" "AUDIT"
move_file "BOUNTY_SUBMISSION_AUDIT_2026_03_11.md" "docs/investigations/BOUNTY_SUBMISSION_AUDIT_2026_03_11.md" "AUDIT"
move_file "BOUNTY_VERIFICATION_AUDIT_2026-03-11.md" "docs/investigations/BOUNTY_VERIFICATION_AUDIT_2026-03-11.md" "AUDIT"
move_file "FABRICATED_BOUNTY_INVESTIGATION.md" "docs/investigations/FABRICATED_BOUNTY_INVESTIGATION.md" "AUDIT"

echo ""

# ============================================
# Category 2: Architecture/Design Docs
# ============================================
echo -e "${YELLOW}🏗️  Moving Architecture Docs...${NC}"

move_file "BOUNTYOS_ARCHITECTURE_V1.md" "docs/architecture/BOUNTYOS_ARCHITECTURE_V1.md" "ARCH"
move_file "PROJECT_ISOLATION_AND_HUMAN_AUTHENTICITY.md" "docs/architecture/PROJECT_ISOLATION_AND_HUMAN_AUTHENTICITY.md" "ARCH"

echo ""

# ============================================
# Category 3: Plans & Activation Docs
# ============================================
echo -e "${YELLOW}📝 Moving Plans...${NC}"

move_file "BOUNTY_OVERDRIVE_ACTIVATION.md" "docs/plans/BOUNTY_OVERDRIVE_ACTIVATION.md" "PLAN"
move_file "BOUNTY_PIPELINE_COMPREHENSIVE.md" "docs/plans/BOUNTY_PIPELINE_COMPREHENSIVE.md" "PLAN"
move_file "BOUNTY_SPEC_VERIFICATION.md" "docs/plans/BOUNTY_SPEC_VERIFICATION.md" "PLAN"
move_file "COMPLETION_PLAN.md" "docs/archive/COMPLETION_PLAN.md" "PLAN"
move_file "CORRECTED_TIMELINE.md" "docs/archive/CORRECTED_TIMELINE.md" "PLAN"
move_file "SCANNER_MIGRATION_PLAN.md" "docs/archive/SCANNER_MIGRATION_PLAN.md" "PLAN"
move_file "VERIFICATION_PROTOCOL_IMPLEMENTATION.md" "docs/archive/VERIFICATION_PROTOCOL_IMPLEMENTATION.md" "PLAN"

echo ""

# ============================================
# Category 4: Platform & Expansion
# ============================================
echo -e "${YELLOW}🌐 Moving Platform Docs...${NC}"

move_file "PLATFORM_EXPANSION_PLAN.md" "docs/research/PLATFORM_EXPANSION_PLAN.md" "RESEARCH"
move_file "POST_COOLIFY_QUALITY_GATES.md" "docs/standards/POST_COOLIFY_QUALITY_GATES.md" "STANDARD"

echo ""

# ============================================
# Category 5: Follow-up Notes
# ============================================
echo -e "${YELLOW}📌 Moving Follow-ups...${NC}"

move_file "mps_51_followup.md" "docs/archive/followups/mps_51_followup.md" "FOLLOWUP"
move_file "mps_55_followup.md" "docs/archive/followups/mps_55_followup.md" "FOLLOWUP"
move_file "coolify_7724_comment.md" "docs/archive/coolify/coolify_7724_comment.md" "COOLIFY"
move_file "coolify_8779_comment.md" "docs/archive/coolify/coolify_8779_comment.md" "COOLIFY"

echo ""

# ============================================
# Category 6: Code/Patches
# ============================================
echo -e "${YELLOW}🔧 Moving Code & Patches...${NC}"

move_file "bounty5_87.patch" "docs/archive/patches/bounty5_87.patch" "PATCH"
move_file "bounty5_87.ts" "docs/archive/code-snippets/bounty5_87.ts" "CODE"
move_file "bounty9_74.ts" "docs/archive/code-snippets/bounty9_74.ts" "CODE"
move_file "rtc_payment_middleware.patch" "docs/archive/patches/rtc_payment_middleware.patch" "PATCH"
move_file "kadeshx_audit.rs" "docs/archive/code-snippets/kadeshx_audit.rs" "CODE" 2>/dev/null || true

echo ""

# ============================================
# Category 7: Scripts & Automation
# ============================================
echo -e "${YELLOW}⚙️  Moving Scripts...${NC}"

move_file "beacon_monitor.sh" "scripts/beacon_monitor.sh" "SCRIPT"
move_file "selenium_script.py" "scripts/selenium_script.py" "SCRIPT"
move_file "capture_atlas.js" "scripts/capture_atlas.js" "SCRIPT"

echo ""

# ============================================
# Category 8: Daily Trackers & Actions
# ============================================
echo -e "${YELLOW}📅 Moving Daily Trackers...${NC}"

move_file "CRITICAL_ACTIONS_TODAY.md" "docs/archive/CRITICAL_ACTIONS_TODAY.md" "TRACKER"
move_file "NEXT_BOUNTIES.md" "docs/status/NEXT_BOUNTIES.md" "STATUS"
move_file "OPERATION_500_PLUS.md" "docs/archive/OPERATION_500_PLUS.md" "TRACKER"

echo ""

# ============================================
# Category 9: Setup & Config Docs
# ============================================
echo -e "${YELLOW}🔧 Moving Setup Docs...${NC}"

move_file "PI_LLM_SETUP.md" "docs/guides/PI_LLM_SETUP.md" "GUIDE"
move_file "ralph-bridge-pinger.md" "docs/archive/ralph-bridge-pinger.md" "ARCHIVE"

echo ""

# ============================================
# Category 10: Investigation Notes
# ============================================
echo -e "${YELLOW}🔍 Moving Investigations...${NC}"

move_file "BLOCKED_BOUNTIES_BREAKDOWN.md" "docs/investigations/BLOCKED_BOUNTIES_BREAKDOWN.md" "INVESTIGATION"
move_file "C4_BOTH_SUBMITTED.md" "docs/submissions/C4_BOTH_SUBMITTED.md" "SUBMISSION" 2>/dev/null || mkdir -p docs/submissions && move_file "C4_BOTH_SUBMITTED.md" "docs/submissions/C4_BOTH_SUBMITTED.md" "SUBMISSION"
move_file "C4_INJECTIVE_SUBMITTED.md" "docs/submissions/C4_INJECTIVE_SUBMITTED.md" "SUBMISSION"

echo ""

# ============================================
# Category 11: Keep in Root (Essential)
# ============================================
echo -e "${YELLOW}ℹ️  Files to Keep in Root...${NC}"

KEEP_FILES=(
    "AGENTS.md"
    "GEMINI.md"
    "readme.md"
    "index.ts"
)

for file in "${KEEP_FILES[@]}"; do
    if [[ -f "$file" ]]; then
        echo -e "${BLUE}ℹ${NC} Keeping: $file"
        SKIPPED=$((SKIPPED + 1))
    fi
done

echo ""

# ============================================
# Summary
# ============================================
echo "========================================"
echo "  Extended Cleanup Summary"
echo "========================================"
echo -e "Files Moved:  ${GREEN}$MOVED${NC}"
echo -e "Files Kept:   ${YELLOW}$SKIPPED${NC}"
echo ""

if [[ "$DRY_RUN" == true ]]; then
    echo -e "${YELLOW}This was a DRY RUN. Run without --dry-run to execute.${NC}"
else
    echo -e "${GREEN}Extended cleanup complete!${NC}"
fi
