# Security Remediation Plan
**Date:** 2026-01-09
**Based on:** Audit Reports `gosec_report_20260109_195602.json`, `trivy_report_20260109_195602.json`

## Summary
The automated audit identified several high and medium severity issues in both the Go codebase and Node.js dependencies. This plan outlines the steps to remediate these vulnerabilities.

## 1. Go Codebase Hardening (Gosec)

### 1.1. Weak TLS Cipher Suite (High)
*   **Issue:** `internal/security/secure_http.go` enables `TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA`, which is considered weak.
*   **Fix:** Remove this cipher suite from the configuration. Retain only GCM-based suites (e.g., `TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384`).

### 1.2. Slowloris Vulnerability (Medium)
*   **Issue:** `internal/adapters/ui/web.go` does not set `ReadHeaderTimeout` on the `http.Server`.
*   **Fix:** Set `ReadHeaderTimeout` to 5 seconds to prevent resource exhaustion attacks.

### 1.3. Path Traversal Risks (Medium)
*   **Issue:** `internal/config/config.go` and `cmd/obsidian/main.go` read/write files using variable paths without explicit sanitization.
*   **Fix:** Implement `filepath.Clean()` and ensure paths reside within expected directories where applicable.

### 1.4. Insecure File Permissions (Medium)
*   **Issue:** `cmd/obsidian/main.go` creates files with `0644` and directories with `0755`.
*   **Fix:** Restrict to `0600` (Read/Write Owner only) and `0700` (or `0750`) respectively to prevent unauthorized local access to logs/configs.

## 2. Dependency Vulnerabilities (Trivy)

### 2.1. NPM Packages (High)
*   **Issue:**
    *   `@modelcontextprotocol/sdk`: ReDoS vulnerability (CWE-1333).
    *   `preact`: Arbitrary script execution (CWE-843).
*   **Fix:** Update `archestra/platform/backend/package.json` and other `package.json` files to the latest safe versions.

## 3. Configuration & Secrets (Trivy)

### 3.1. Dockerfile Environment Variables (Critical/False Positive)
*   **Issue:** Trivy flagged `ARCHESTRA_AUTH_DISABLE_BASIC_AUTH` as a potential secret leak.
*   **Analysis:** These appear to be boolean feature flags, not actual secrets.
*   **Fix:** Verify contents. If non-sensitive, add to `.trivyignore`. If sensitive, switch to Docker secrets mounting.

## Execution Order
1.  Apply Go code fixes (TLS, Timeouts, Permissions).
2.  Update NPM dependencies.
3.  Re-run `scripts/security-audit.sh` to verify fixes.
