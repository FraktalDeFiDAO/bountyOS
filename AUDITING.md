# Universal Security & Operational Auditing Standard

This document outlines a **language-agnostic**, **container-based** auditing framework designed to be dropped into any software project. It covers auditing standards, event schemas, and isolated security scanning procedures.

## 1. Philosophy: The "Containerized Audit"

To ensure this process works on **any** project (Go, Python, Node, Rust, etc.) without polluting the host environment or requiring complex toolchain setups, all security assessments are performed within **ephemeral, isolated containers** (Podman or Docker).

**Principles:**
1.  **Isolation:** Tools run in containers; the codebase is mounted read-only.
2.  **Granularity:** Each language/framework uses its specific scanner container.
3.  **Standardization:** All tools output to a unified `audit-reports/` directory.

## 2. Audit Event Schema (Universal)

All applications, regardless of language, must emit audit logs adhering to this JSON structure. This ensures downstream log aggregators (Splunk, ELK, Datadog) can ingest events from polyglot microservices uniformly.

### Standard JSON Structure
```json
{
  "timestamp": "2023-10-27T10:00:00Z",
  "event_type": "AUDIT",
  "actor": {
    "id": "user_123", // or service_account_id
    "ip": "192.168.1.50",
    "role": "admin"
  },
  "action": "resource.verb", // e.g., "auth.login", "file.delete", "model.inference"
  "resource": {
    "type": "prompt", // or "database", "file", "api"
    "id": "prompt_555",
    "owner_org": "org_abc"
  },
  "status": "success", // "success", "failure", "denied"
  "metadata": {
    "language": "python", // Optional: Tracking source context
    "previous_version": "v1.0",
    "new_version": "v1.1"
  }
}
```

### Implementation Guide
*   **Go:** Use `encoding/json` with a structured logger (e.g., Zap, Logrus) to serialise the struct.
*   **Node/TS:** Use `pino` or `winston` with a custom serializer.
*   **Python:** Use `structlog` or standard `logging` with a JSON formatter.
*   **Rust:** Use `tracing` with a JSON subscriber.

## 3. Operational vs. Security Logging

| Feature | Operational Logs | Security Audit Logs |
| :--- | :--- | :--- |
| **Purpose** | Debugging, health checks, performance | Compliance, forensics, non-repudiation |
| **Retention** | Short-term (e.g., 30 days) | Long-term (e.g., 1-7 years) |
| **Content** | Stack traces, timings, debug states | Who, What, Where, When, Result |
| **Sensitivity** | Medium (Sanitized) | High (Strictly Immutable) |

## 4. Automated Security Scanning (The "Any Project" Workflow)

We use a central script (`scripts/security-audit.sh`) that auto-detects the project's languages and runs the appropriate containerized scanners.

### Supported Scanners (Extensible)

| Language/Tech | Tool | Container Image | Description |
| :--- | :--- | :--- | :--- |
| **Universal** | **Trivy** | `docker.io/aquasec/trivy` | Scans filesystem, secrets, and misconfigurations. |
| **Go** | **Gosec** | `docker.io/securego/gosec` | Static analysis for Go security flaws. |
| **Node.js** | **NPM Audit** | `docker.io/library/node` | Checks `package-lock.json` for vulnerabilities. |
| **Python** | **Bandit** | `docker.io/secfigo/bandit` | Security linter for Python code. |
| **Infrastructure** | **Checkov** | `docker.io/bridgecrew/checkov` | Scans Terraform, Kubernetes, Dockerfile. |

### Usage
Run the single entry point. It will auto-detect your project structure and apply relevant checks.

```bash
./scripts/security-audit.sh
```

## 5. AI & Model Auditing
For projects involving LLMs/Agents:
- **Prompt Versioning:** All prompts must be versioned in a database.
- **Inference Logging:** Log the `model_id`, `input_token_count`, `output_token_count`, and `cost` in the `metadata` field of the Audit Event.
- **Safety Checks:** Use the `AUDIT` event to log when a safety guardrail (e.g., PII filter) is triggered.