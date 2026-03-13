#!/bin/bash
# Bounty Verification Script
# Purpose: Verify all bounty URLs exist and capture evidence
# Usage: ./verify_bounties.sh [input_file] [output_dir]

set -e

# Configuration
INPUT_FILE="${1:-bounty_urls.txt}"
OUTPUT_DIR="${2:-./audit_trails/$(date +%Y-%m-%d)}"
TIMESTAMP=$(date +%Y-%m-%d_%H-%M-%S)
LOG_FILE="$OUTPUT_DIR/verification_log.txt"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Create output directory
mkdir -p "$OUTPUT_DIR/http_responses"
mkdir -p "$OUTPUT_DIR/screenshots"
mkdir -p "$OUTPUT_DIR/verification_logs"

echo "=== BOUNTY VERIFICATION SCRIPT ===" | tee "$LOG_FILE"
echo "Timestamp: $TIMESTAMP" | tee -a "$LOG_FILE"
echo "Input: $INPUT_FILE" | tee -a "$LOG_FILE"
echo "Output: $OUTPUT_DIR" | tee -a "$LOG_FILE"
echo "" | tee -a "$LOG_FILE"

# Initialize counters
TOTAL=0
HTTP_200=0
HTTP_404=0
HTTP_OTHER=0

# Function to verify a single URL
verify_url() {
    local url="$1"
    local bounty_name="$2"
    local issue_num=$(echo "$url" | grep -oP 'issues/\K[0-9]+' || echo "unknown")
    
    TOTAL=$((TOTAL + 1))
    
    echo -n "Verifying #$issue_num ($bounty_name)... " | tee -a "$LOG_FILE"
    
    # Get HTTP status
    HTTP_STATUS=$(curl -s -o /dev/null -w "%{http_code}" "$url" 2>/dev/null || echo "000")
    TIMESTAMP=$(date +%Y-%m-%d_%H-%M-%S)
    
    # Save evidence
    local evidence_file="$OUTPUT_DIR/http_responses/${bounty_name}_${issue_num}_${TIMESTAMP}.txt"
    curl -I "$url" 2>/dev/null > "$evidence_file" || echo "Failed to fetch headers" > "$evidence_file"
    
    # Categorize by status
    if [ "$HTTP_STATUS" = "200" ]; then
        HTTP_200=$((HTTP_200 + 1))
        echo -e "${GREEN}✅ $HTTP_STATUS OK${NC}" | tee -a "$LOG_FILE"
    elif [ "$HTTP_STATUS" = "404" ]; then
        HTTP_404=$((HTTP_404 + 1))
        echo -e "${RED}❌ $HTTP_STATUS Not Found${NC}" | tee -a "$LOG_FILE"
    else
        HTTP_OTHER=$((HTTP_OTHER + 1))
        echo -e "${YELLOW}⚠️ $HTTP_STATUS Other${NC}" | tee -a "$LOG_FILE"
    fi
    
    echo "  Evidence: $evidence_file" | tee -a "$LOG_FILE"
    echo "  URL: $url" | tee -a "$LOG_FILE"
    echo "" | tee -a "$LOG_FILE"
}

# Read URLs from input file
if [ ! -f "$INPUT_FILE" ]; then
    echo "Error: Input file not found: $INPUT_FILE" | tee -a "$LOG_FILE"
    echo "Usage: $0 [bounty_urls.txt] [output_dir]" | tee -a "$LOG_FILE"
    exit 1
fi

echo "Starting verification..." | tee -a "$LOG_FILE"
echo "" | tee -a "$LOG_FILE"

# Process each line in input file
while IFS='|' read -r bounty_name url || [ -n "$bounty_name" ]; do
    # Skip empty lines and comments
    [[ -z "$bounty_name" || "$bounty_name" =~ ^# ]] && continue
    
    # Trim whitespace
    bounty_name=$(echo "$bounty_name" | xargs)
    url=$(echo "$url" | xargs)
    
    verify_url "$url" "$bounty_name"
    
done < "$INPUT_FILE"

# Generate summary
echo "" | tee -a "$LOG_FILE"
echo "=== VERIFICATION SUMMARY ===" | tee -a "$LOG_FILE"
echo "Total Bounties: $TOTAL" | tee -a "$LOG_FILE"
echo -e "${GREEN}HTTP 200 (Exists): $HTTP_200${NC}" | tee -a "$LOG_FILE"
echo -e "${RED}HTTP 404 (Not Found): $HTTP_404${NC}" | tee -a "$LOG_FILE"
echo -e "${YELLOW}HTTP Other: $HTTP_OTHER${NC}" | tee -a "$LOG_FILE"
echo "" | tee -a "$LOG_FILE"

# Calculate percentage
if [ $TOTAL -gt 0 ]; then
    PERCENT_200=$((HTTP_200 * 100 / TOTAL))
    echo "Verification Rate: ${PERCENT_200}%" | tee -a "$LOG_FILE"
fi

echo "" | tee -a "$LOG_FILE"
echo "Evidence saved to: $OUTPUT_DIR/http_responses/" | tee -a "$LOG_FILE"
echo "Log file: $LOG_FILE" | tee -a "$LOG_FILE"
echo "" | tee -a "$LOG_FILE"
echo "=== VERIFICATION COMPLETE ===" | tee -a "$LOG_FILE"

# Export summary for reporting
cat > "$OUTPUT_DIR/summary.json" << EOF
{
  "timestamp": "$TIMESTAMP",
  "total": $TOTAL,
  "http_200": $HTTP_200,
  "http_404": $HTTP_404,
  "http_other": $HTTP_OTHER,
  "verification_rate": $PERCENT_200
}
EOF

echo "Summary JSON: $OUTPUT_DIR/summary.json"
