#!/bin/bash
# bountyOS Project Structure Validator
#
# Validates that the project follows organization standards
#
# Usage: ./scripts/validate-structure.sh

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Get the directory where the script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# Root directory is parent of scripts directory
ROOT_DIR="$(dirname "$SCRIPT_DIR")"

cd "$ROOT_DIR"

echo "Validating: $ROOT_DIR"
echo ""

ERRORS=0
WARNINGS=0

echo "========================================"
echo "  bountyOS Structure Validation"
echo "========================================"
echo ""

# ============================================
# Check 1: Root Directory Files
# ============================================
echo -e "${BLUE}📁 Checking Root Directory...${NC}"

# Allowed files in root
ALLOWED_ROOT_PATTERNS=(
    "*.md"  # We'll check specific MD files below
    "*.json"
    "*.toml"
    "*.mod"
    "*.sum"
    "*.yaml"
    "*.yml"
    "*.ts"
    "*.py"
    "*.sh"
    ".*"
)

# Files that should NOT be in root
FORBIDDEN_ROOT_PATTERNS=(
    "*_STATUS.md"
    "*_SUMMARY.md"
    "*_TRACKER.md"
    "beacon_*.py"
)

echo "Checking for forbidden file patterns in root..."

for pattern in "${FORBIDDEN_ROOT_PATTERNS[@]}"; do
    files=$(find . -maxdepth 1 -name "$pattern" -type f 2>/dev/null || true)
    if [[ -n "$files" ]]; then
        echo -e "${RED}✗${NC} Found forbidden pattern '$pattern' in root:"
        echo "   $files"
        ((ERRORS++))
    fi
done

# Check for specific files that should be moved
SHOULD_MOVE=(
    "app.py"
    "approvals.ts"
    "beacon_auto_submit.py"
    "beacon_heartbeat.py"
    "beacon_register.py"
)

for file in "${SHOULD_MOVE[@]}"; do
    if [[ -f "$file" ]]; then
        echo -e "${YELLOW}⚠${NC} File should be moved: $file"
        ((WARNINGS++))
    fi
done

echo ""

# ============================================
# Check 2: Required Directories
# ============================================
echo -e "${BLUE}📂 Checking Required Directories...${NC}"

REQUIRED_DIRS=(
    "au-workspace"
    "au-workspace/projects"
    "au-workspace/research"
    "docs"
    "docs/status"
    "docs/sessions"
    "docs/guides"
    "docs/architecture"
    "docs/standards"
    "docs/archive"
    "scripts"
    "submissions"
    "config"
)

# Optional directories (may not exist yet)
OPTIONAL_DIRS=(
    ".agents"
)

for dir in "${REQUIRED_DIRS[@]}"; do
    if [[ -d "$dir" ]]; then
        echo -e "${GREEN}✓${NC} $dir"
    else
        echo -e "${RED}✗${NC} Missing: $dir"
        ((ERRORS++))
    fi
done

echo ""
echo -e "${BLUE}ℹ️  Checking Optional Directories...${NC}"

for dir in "${OPTIONAL_DIRS[@]}"; do
    if [[ -d "$dir" ]]; then
        echo -e "${GREEN}✓${NC} $dir (exists)"
    else
        echo -e "${YELLOW}⚠${NC} $dir (optional, not found)"
    fi
done

echo ""

# ============================================
# Check 3: Agent Structure
# ============================================
echo -e "${BLUE}🤖 Checking Agent Structure...${NC}"

if [[ -d ".agents" ]]; then
    agent_count=$(find .agents -maxdepth 1 -type d | wc -l)
    echo -e "${GREEN}✓${NC} Found $((agent_count - 1)) agents"
    
    # Check for orchestrator
    if [[ -d ".agents/orchestrator" ]]; then
        echo -e "${GREEN}✓${NC} Orchestrator agent exists"
    else
        echo -e "${YELLOW}⚠${NC} Missing orchestrator agent (recommended)"
        ((WARNINGS++))
    fi
else
    echo -e "${YELLOW}⚠${NC} .agents directory not found (optional)"
    ((WARNINGS++))
fi

echo ""

# ============================================
# Check 4: Project Structure
# ============================================
echo -e "${BLUE}📦 Checking Project Structure...${NC}"

if [[ -d "au-workspace/projects" ]]; then
    project_count=$(find au-workspace/projects -maxdepth 1 -type d | wc -l)
    echo -e "${GREEN}✓${NC} Found $((project_count - 1)) projects"
    
    # Check for index.md in projects
    projects_without_index=0
    for project_dir in au-workspace/projects/*/; do
        if [[ ! -f "${project_dir}index.md" ]]; then
            ((projects_without_index++))
        fi
    done
    
    if [[ $projects_without_index -gt 0 ]]; then
        echo -e "${YELLOW}⚠${NC} $projects_without_index projects missing index.md"
        ((WARNINGS++))
    else
        echo -e "${GREEN}✓${NC} All projects have index.md"
    fi
else
    echo -e "${RED}✗${NC} au-workspace/projects directory missing"
    ((ERRORS++))
fi

echo ""

# ============================================
# Check 5: Documentation Structure
# ============================================
echo -e "${BLUE}📚 Checking Documentation Structure...${NC}"

REQUIRED_DOCS=(
    "docs/PROJECT_ORGANIZATION_STANDARDS.md"
    "AGENTS.md"
)

for doc in "${REQUIRED_DOCS[@]}"; do
    if [[ -f "$doc" ]]; then
        echo -e "${GREEN}✓${NC} $doc"
    else
        echo -e "${YELLOW}⚠${NC} Missing recommended doc: $doc"
        ((WARNINGS++))
    fi
done

echo ""

# ============================================
# Check 6: Submissions Structure
# ============================================
echo -e "${BLUE}🏆 Checking Submissions Structure...${NC}"

SUBMISSION_DIRS=(
    "submissions/code4rena"
    "submissions/immunefi"
    "submissions/superteam"
)

for dir in "${SUBMISSION_DIRS[@]}"; do
    if [[ -d "$dir" ]]; then
        echo -e "${GREEN}✓${NC} $dir"
    else
        echo -e "${YELLOW}⚠${NC} Missing submission dir: $dir"
        ((WARNINGS++))
    fi
done

echo ""

# ============================================
# Summary
# ============================================
echo "========================================"
echo "  Validation Summary"
echo "========================================"
echo -e "Errors:   ${RED}$ERRORS${NC}"
echo -e "Warnings: ${YELLOW}$WARNINGS${NC}"
echo ""

if [[ $ERRORS -gt 0 ]]; then
    echo -e "${RED}❌ Validation FAILED${NC}"
    echo "Please fix the errors above."
    exit 1
elif [[ $WARNINGS -gt 0 ]]; then
    echo -e "${YELLOW}⚠️  Validation PASSED with warnings${NC}"
    echo "Consider addressing the warnings."
    exit 0
else
    echo -e "${GREEN}✅ Validation PASSED${NC}"
    echo "Project structure looks good!"
    exit 0
fi
