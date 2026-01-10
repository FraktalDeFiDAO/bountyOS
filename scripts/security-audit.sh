#!/bin/bash
set -e

# security-audit.sh
# Runs a suite of security scanners against the codebase using containerized tools.

# Config
REPORT_DIR="./audit-reports"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
mkdir -p "$REPORT_DIR"

# Detect Container Runtime
if command -v podman &> /dev/null; then
    RUNTIME="podman"
    echo "Using Podman for isolation."
elif command -v docker &> /dev/null; then
    RUNTIME="docker"
    echo "Using Docker for isolation."
else
    echo "Error: Neither podman nor docker found. Cannot run isolated security scans."
    exit 1
fi

echo "Starting Security Audit at $TIMESTAMP..."

# 1. Trivy (Filesystem Scan)
echo "------------------------------------------------"
echo "Running Trivy (Filesystem & Dependencies)..."
$RUNTIME run --rm \
    -v "$PWD:/project:ro" \
    -v "$PWD/audit-reports:/reports:rw" \
    docker.io/aquasec/trivy:latest fs /project \
    --skip-dirs workspace \
    --scanners vuln,secret,misconfig \
    --format json \
    --output "/reports/trivy_report_$TIMESTAMP.json"

echo "Trivy scan complete. Report: $REPORT_DIR/trivy_report_$TIMESTAMP.json"

# 2. Gosec (Go Static Analysis)
if [ -f "go.mod" ]; then
    echo "------------------------------------------------"
    echo "Running Gosec (Go Code)..."
    $RUNTIME run --rm \
        -v "$PWD:/project:ro" \
        -v "$PWD/audit-reports:/reports:rw" \
        -w /project \
        docker.io/securego/gosec:latest -exclude-dir=workspace -fmt=json -out="/reports/gosec_report_$TIMESTAMP.json" ./...
    echo "Gosec scan complete. Report: $REPORT_DIR/gosec_report_$TIMESTAMP.json"
fi

# 3. NPM Audit (Node.js Dependencies)
# We run this in a standard node container to avoid polluting the host, 
# strictly matching the node version from .nvmrc if possible, defaulting to lts.
if [ -f "package.json" ] || [ -f "archestra/platform/backend/package.json" ]; then
    echo "------------------------------------------------"
    echo "Running NPM Audit..."
    # We mount the whole project to handle nested package.json files if needed
    $RUNTIME run --rm \
        -v "$PWD:/project:rw" \
        -w /project \
        docker.io/library/node:20-alpine npm audit --json > "$REPORT_DIR/npm_audit_$TIMESTAMP.json" || true
    echo "NPM Audit complete. Report: $REPORT_DIR/npm_audit_$TIMESTAMP.json"
fi

# 4. Bandit (Python Static Analysis)
if [ -f "requirements.txt" ] || [ -f "pyproject.toml" ] || [ -f "Pipfile" ] || [ -n "$(find . -maxdepth 2 -name '*.py' -print -quit)" ]; then
    echo "------------------------------------------------"
    echo "Running Bandit (Python Code)..."
    $RUNTIME run --rm \
        -v "$PWD:/project:ro" \
        -v "$PWD/audit-reports:/reports:rw" \
        -w /project \
        docker.io/secfigo/bandit:latest -x workspace -r . -f json -o /reports/bandit_report_$TIMESTAMP.json || true
    echo "Bandit scan complete. Report: $REPORT_DIR/bandit_report_$TIMESTAMP.json"
fi

# 5. Checkov (Infrastructure as Code Scan)
echo "------------------------------------------------"
echo "Running Checkov (IaC Scan)..."
$RUNTIME run --rm \
    -v "$PWD:/project:ro" \
    -v "$PWD/audit-reports:/reports:rw" \
    -w /project \
    docker.io/bridgecrew/checkov:latest --skip-path workspace -d . --output json --soft-fail > "audit-reports/checkov_report_$TIMESTAMP.json" || true
echo "Checkov scan complete. Report: $REPORT_DIR/checkov_report_$TIMESTAMP.json"

echo "------------------------------------------------"
echo "Full Security Audit Complete."
echo "All reports saved in $REPORT_DIR"
