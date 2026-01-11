# GitHub Bounties: End-to-End Process (BountyOS v8)

This document describes **exactly how GitHub bounties are discovered, validated, scored, stored, and surfaced** in this repo. It also includes common edge cases and an ops troubleshooting decision tree.

> Scope: GitHub bounties **inside this app** (Obsidian). Not general GitHub bounty programs.
> Ops quick triage: see `docs/runbook-github-bounties.md`.

---

## 1) High-level Flow (Data Path)

1. **Config load** (defaults → `config/config.yaml` → env overrides).
2. **Scanner selection** based on `ENABLED_SCANNERS`.
3. **GitHub scanner runs**: for each label → for each page → GitHub Search API.
4. **Response validation** (schema + XSS checks). If invalid → label scan stops.
5. **Issue → Bounty mapping** (reward/currency inference, tags).
6. **Global pipeline** (URL checks, optional reachability, sanitize, dedupe).
7. **Score** using Obsidian rules.
8. **Persist** to SQLite and **broadcast** to UI / WebSocket.
9. **Notify** (desktop + Discord) if score ≥ `MIN_SCORE`.

---

## 2) Configuration (GitHub-specific)

All of these can be set in `config/config.yaml` or as env vars.

| Key | Default | Purpose | Notes |
|---|---|---|---|
| `GITHUB_TOKEN` | `""` | Auth token for GitHub API | Strongly recommended to avoid rate limits. |
| `GITHUB_LABELS` | `algora-bounty, polar, opire, gitpay, issuehunt, bounty, funded` | Labels to search | Trimmed and de-duped but **case preserved**. |
| `GITHUB_PER_PAGE` | `100` | Results per API page | Max allowed by GitHub Search API. |
| `GITHUB_MAX_PAGES` | `10` | Page cap per label | Limits total scan size per label. |
| `GITHUB_BASE_URL` | `https://api.github.com` | GitHub API base | Change for GitHub Enterprise. |

**Enable GitHub scans** by including **either** `GITHUB_AGGREGATOR` or `GITHUB` in `ENABLED_SCANNERS`.

---

## 3) GitHub Scanner Details

### 3.1 Query Construction
For each label and page:

- **Query**: `is:issue is:open label:<LABEL> sort:created-desc`
- **Request**: `GET /search/issues?q=<query>&per_page=<perPage>&page=<page>`

Notes:
- Only **open issues** are included.
- No repo/org filter → searches **all of GitHub** for that label.
- The label is only space-normalized; special characters are not URL-escaped.

### 3.2 HTTP Hardening & Auth
Requests are created with security headers:
- `Accept: application/vnd.github.v3+json`
- `User-Agent: BountyOS-Secure/1.0`
- `X-Requested-With: XMLHttpRequest`
- `Authorization: token <GITHUB_TOKEN>` (if provided)

### 3.3 Rate Limiting (GitHub)
The app enforces a minimum delay:
- **No token** → 10s between requests.
- **With token** → 2s between requests.

Additionally, if remaining GitHub quota is ≤ 5, it sleeps until the reset time.

To disable rate-limit sleeping:
- Set `BOUNTYOS_DISABLE_RATE_LIMIT_SLEEP=1`.

### 3.4 Retry Behavior
- Retries **only** on `5xx` or `429` status codes.
- Up to 3 retries with exponential backoff.
- Other errors (e.g., 401/403/422) are **not retried**.

---

## 4) Response Validation (Strict)

Every response is validated before parsing into bounties:

**Validation rules**
- JSON must decode into the expected search schema.
- Each issue must pass:
  - `title` non-empty
  - `title` length ≤ 500 chars
  - `html_url` is a valid URL
  - `created_at` RFC3339
  - **No script tags or JS event handlers** in title/body

**Critical outlier**
- If **any single issue** fails validation, the **entire response is rejected**, and the **label scan stops** for that cycle.

Example log:
```
[ERROR] Error validating response for bounty (page 5): invalid item at index 73: potential XSS content detected
```

---

## 5) Issue → Bounty Mapping

Each validated GitHub issue maps to a `core.Bounty`:

| Field | Source | Notes |
|---|---|---|
| `ID` | `html_url` | URL is unique ID. |
| `URL` | `html_url` | Stored + used for dedupe. |
| `Title` | `title` | Sanitized later. |
| `Description` | `body` | Sanitized later. |
| `CreatedAt` | `created_at` | Parsed RFC3339. |
| `Platform` | `GITHUB/<LABEL>` | Label uppercased. |
| `Reward` | inferred | See below. |
| `Currency` | inferred | See below. |
| `PaymentType` | inferred | `crypto`, `fiat`, `p2p`. |
| `Tags` | inferred | See below. |

### 5.1 Reward/Currency/PaymentType inference
Defaults:
- `Reward = "Funded"`
- `Currency = "USD"`
- `PaymentType = "fiat"`

**Label-based overrides** (checked first):
- If label contains `$` → `Reward = label`, `Currency = ""`
- If label contains `usdc|eth|sol|usdt` → `Reward = label`, `Currency = ""`, `PaymentType = "crypto"`

**Body-based overrides** (only if still fiat):
- If body contains `usdc|eth|sol|usdt` → `PaymentType = "crypto"`, `Currency = "USDC/ETH/SOL"`
- If body contains `paypal` → `Currency = "PAYPAL"`
- If body contains `cash app|cashapp` → `Currency = "CASHAPP"`, `PaymentType = "p2p"`

### 5.2 Tags
Tags are used later for scoring and UI:
- Always includes `active`.
- Adds:
  - `urgent` if title contains “urgent”
  - `dev` if title contains “fix” or “bug”
  - `automation` if title contains “script” or “bot”
  - `funded` if any label contains “funded”

---

## 6) Global Post-Processing Pipeline

After the scanner emits a bounty, the **global pipeline** enforces consistency:

1. **URL normalization**
   - Strips whitespace and trailing punctuation.
   - Uses the first whitespace-separated token.

2. **URL validation**
   - Rejects non-`http/https`.
   - Rejects local URLs unless `BOUNTYOS_ALLOW_LOCAL_URLS=true`.

3. **Optional reachability check** (`VALIDATE_LINKS_HTTP=true`)
   - Sends HEAD, then GET if needed.
   - Accepts 2xx/3xx, plus 401/403/405/429.

4. **Sanitization**
   - Title/platform/reward/currency/description sanitized and truncated.

5. **Deduplication**
   - Uses URL as primary key.
   - If already present, the bounty is dropped silently.

6. **Scoring**
   - `core.CalculateUrgency` uses payment tiers, keyword hits, recency, platform weights, and tags.

7. **Persistence & output**
   - Stored in SQLite.
   - Broadcast to Web UI.
   - Notifications sent if score ≥ `MIN_SCORE`.

---

## 7) Performance & Scale Expectations

**Max issues per scan**
- `GITHUB_LABELS` × `GITHUB_MAX_PAGES` × `GITHUB_PER_PAGE`
- Example default: `7 labels × 10 pages × 100 = 7000` issues per scan

**Rate limiting**
- Unauthenticated: minimum **10s per request** (very slow at scale).
- Authenticated: minimum **2s per request**.

**Operational impact**
- Large label sets + high page caps can cause scan cycles to exceed `POLL_INTERVAL_SECONDS`.
- If scans take too long, they overlap less often and data gets stale.

---

## 8) Outliers & Edge Cases

### Label collisions (same issue matches multiple labels)
- **Impact**: only the first label that hits the issue is stored.
- **Consequence**: `Platform`, `Reward`, and tags may reflect the earlier label, not the “best” one.
- **Mitigation**: order `GITHUB_LABELS` to prioritize preferred label.

### Broad labels (e.g., `bounty`)
- **Impact**: enormous result set; you only get the first `GITHUB_MAX_PAGES` pages.
- **Mitigation**: use narrow labels, or raise page limit cautiously.

### Validation fail on one issue (XSS)
- **Impact**: aborts the **entire** response for that label/page.
- **Mitigation**:
  - Reduce problematic labels.
  - Relax `ValidateGitHubResponse` or skip invalid items (code change).

### Missing reward data
- **Impact**: many GitHub issues don’t list payout info, so `Reward` stays `Funded` and `Currency` `USD`.
- **Mitigation**: encourage labels with currency or embed payout text in issue body.

### URL reachability (optional)
- **Impact**: if `VALIDATE_LINKS_HTTP=true`, a 404/5xx from GitHub will drop the bounty.
- **Mitigation**: set `VALIDATE_LINKS_HTTP=false` for GitHub scans, or extend status allowlist.

### Search query parsing
- **Impact**: labels with special characters can produce query errors or unexpected results.
- **Mitigation**: keep labels simple and URL-safe.

---

## 9) Ops Troubleshooting Decision Tree

```
START
  |
  |-- Are GitHub bounties appearing in the UI/API?
  |       |
  |       |-- YES --> Done.
  |       |
  |       |-- NO --> Check logs for GitHub errors.
  |                   |
  |                   |-- Error: "invalid item ... potential XSS"?
  |                   |       |
  |                   |       |-- YES --> Validation is rejecting a single issue.
  |                   |       |          Actions:
  |                   |       |          1) Identify label causing it (log includes label/page).
  |                   |       |          2) Temporarily remove that label.
  |                   |       |          3) Consider code change to skip invalid items instead of failing the page.
  |                   |       |
  |                   |       |-- NO --> Continue.
  |                   |
  |                   |-- Error: 401/403 from GitHub?
  |                   |       |
  |                   |       |-- YES --> Token missing/invalid.
  |                   |       |          Actions: set `GITHUB_TOKEN`, verify scopes.
  |                   |       |
  |                   |       |-- NO --> Continue.
  |                   |
  |                   |-- Error: 422 / bad search query?
  |                   |       |
  |                   |       |-- YES --> Label likely malformed.
  |                   |       |          Actions: simplify/remove label.
  |                   |       |
  |                   |       |-- NO --> Continue.
  |                   |
  |                   |-- Error: 429 / rate limit?
  |                   |       |
  |                   |       |-- YES --> Use `GITHUB_TOKEN` or reduce label/page count.
  |                   |       |
  |                   |       |-- NO --> Continue.
  |                   |
  |                   |-- No errors, still empty?
  |                           |
  |                           |-- Check `ENABLED_SCANNERS` includes `GITHUB_AGGREGATOR` or `GITHUB`.
  |                           |-- Check `GITHUB_LABELS` are valid labels used on GitHub issues.
  |                           |-- Check `VALIDATE_LINKS_HTTP` isn't dropping results.
  |
  |-- Are bounties present but missing/incorrect rewards?
  |       |
  |       |-- YES --> Reward inference is label/body based.
  |       |          Actions:
  |       |          - Use labels with currency (e.g., "100 USDC").
  |       |          - Ensure body contains payment hints (usdc/eth/sol/paypal/cashapp).
  |
  |-- Are bounties stale or too few?
          |
          |-- YES --> Increase `GITHUB_MAX_PAGES`, reduce broad labels, or shorten `POLL_INTERVAL_SECONDS`.
          |          Ensure token is set to avoid slow rate limits.
```

---

## 10) Recommended Baseline (Ops)

- Set `GITHUB_TOKEN` in env.
- Use 3–6 **specific** labels.
- Start with:
  - `GITHUB_PER_PAGE=100`
  - `GITHUB_MAX_PAGES=5`
- Keep `VALIDATE_LINKS_HTTP=false` for GitHub-heavy scans unless you need strict reachability.
- Periodically review logs for validation failures and adjust labels.

---

## 11) Where to Look in Code

- GitHub scanner: `internal/adapters/scanners/github.go`
- Response validation: `internal/security/validation.go`
- Rate limiter: `internal/security/rate_limiter.go`
- Scoring rules: `internal/core/score.go`
- Global pipeline: `cmd/obsidian/main.go`
- Config defaults: `internal/config/config.go`
- Config file: `config/config.yaml`
