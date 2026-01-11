# Ops Runbook: GitHub Bounties (BountyOS v8)

This is the short, ops-focused runbook for GitHub bounties in this repo. For the full pipeline and deep dive, see `docs/github-bounties.md`.

---

## 0) Quick Facts

- **Scanner name**: `GITHUB_AGGREGATOR` (alias: `GITHUB`)
- **Config file**: `config/config.yaml`
- **Database**: `data/bounties.db` (SQLite)
- **Logs**: `data/bountyos.log` (or container logs)
- **Web API**: `http://localhost:12496`

---

## 1) Quick Health Check (2–3 minutes)

1. **Is the service running?**
   - `podman ps` (or `docker ps`)
2. **Any errors in logs?**
   - `podman logs -f bountyos-obsidian-dev`
   - or `tail -n 200 data/bountyos.log`
3. **Is GitHub enabled?**
   - `ENABLED_SCANNERS` includes `GITHUB_AGGREGATOR` or `GITHUB`.
4. **Token set?**
   - `GITHUB_TOKEN` should be set for stable rate limits.

---

## 2) Symptom → Action Map

### A) No GitHub bounties at all
**Likely causes**
- GitHub scanner disabled
- Bad token / no token + rate limit
- Label set too narrow or invalid
- Validation failures stopping scan

**Actions**
1. Confirm `ENABLED_SCANNERS` contains `GITHUB_AGGREGATOR` or `GITHUB`.
2. Check logs for `Error validating response` or HTTP errors.
3. Ensure `GITHUB_LABELS` are valid labels used on GitHub.
4. Set `GITHUB_TOKEN` if missing.

---

### B) Errors: `invalid item ... potential XSS content detected`
**Cause**
- A single GitHub issue contains content matching XSS rules.
- **Entire response is rejected** for that label/page.

**Actions**
1. Note the label + page from logs.
2. Temporarily remove that label to restore flow.
3. If this happens often, consider relaxing validation or skipping invalid items (code change).

---

### C) Errors: 401/403 from GitHub
**Cause**
- Missing or invalid token.

**Actions**
1. Set `GITHUB_TOKEN` (env or config).
2. Verify token validity and permissions.

---

### D) Errors: 429 rate limit
**Cause**
- Unauthenticated requests or too many labels/pages.

**Actions**
1. Set `GITHUB_TOKEN`.
2. Reduce `GITHUB_MAX_PAGES` or label count.
3. Increase `POLL_INTERVAL_SECONDS`.

---

### E) Bounties appear but rewards/currency are wrong
**Cause**
- Reward inference is heuristic (label/body based).

**Actions**
1. Use labels that include currency or dollar value (e.g., `100 USDC`).
2. Encourage explicit payout text in issue bodies.
3. Adjust inference rules (code change in `internal/adapters/scanners/github.go`).

---

### F) Bounties are stale or too few
**Cause**
- Small `GITHUB_MAX_PAGES`, narrow labels, or long scan time due to rate limits.

**Actions**
1. Increase `GITHUB_MAX_PAGES` (careful with rate limits).
2. Expand labels.
3. Set `GITHUB_TOKEN` for faster requests.
4. Check scan time vs `POLL_INTERVAL_SECONDS`.

---

### G) Duplicates missing (expected repeats)
**Cause**
- URL is primary key; duplicates are dropped silently.

**Actions**
- This is expected behavior. Clear `data/bounties.db` if you need to re-ingest.

---

## 3) Decision Tree (Ops Triage)

```
START
  |
  |-- No GitHub bounties?
  |       |
  |       |-- Logs show XSS validation error?
  |       |       -> Remove label or relax validation.
  |       |
  |       |-- Logs show 401/403?
  |       |       -> Set/refresh GITHUB_TOKEN.
  |       |
  |       |-- Logs show 429?
  |       |       -> Add token, reduce pages/labels.
  |       |
  |       |-- No errors?
  |               -> Verify ENABLED_SCANNERS + GITHUB_LABELS.
  |
  |-- Bounties present but wrong payouts?
  |       -> Label/body inference issue. Adjust labels or code.
  |
  |-- Bounties stale?
          -> Increase pages/labels or add token. Verify scan time.
```

---

## 4) Safe Operational Commands

**Restart services (Podman Compose)**
```
podman compose -f docker-compose.dev.yml down
podman compose -f docker-compose.dev.yml up --build -d
```

**Check logs**
```
podman logs -f bountyos-obsidian-dev
```

**Reset stored bounties** (destructive)
```
rm -f data/bounties.db
```

---

## 5) Known Failure Modes (Watchlist)

- **XSS validation rejects a page** → label scan stops for that cycle.
- **Broad labels** (e.g., `bounty`) → huge result sets; only first N pages are processed.
- **Unauthenticated scans** → slow and rate limited.
- **Special characters in labels** → query may fail or return unexpected results.

---

## 6) References

- Full pipeline doc: `docs/github-bounties.md`
- GitHub scanner: `internal/adapters/scanners/github.go`
- Validation: `internal/security/validation.go`
- Rate limiter: `internal/security/rate_limiter.go`
- Global pipeline: `cmd/obsidian/main.go`

