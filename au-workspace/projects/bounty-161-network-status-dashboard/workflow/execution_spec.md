# Execution Spec: bounty-161-net-status

Issue: https://github.com/Scottcjn/Rustchain/issues/161
Title: [BOUNTY] RustChain Network Status Page — Public Health Dashboard (25 RTC)
Last updated: 2026-02-14

## 1) Goal

Build a public health dashboard that monitors 3 RustChain nodes and shows status, latency, uptime windows, and core network metrics.

## 2) Required Output

Per node:
- status color (green/yellow/red)
- uptime % for 24h/7d/30d
- response time
- current epoch
- active miners count
- last block time

Implementation requirements:
- single-page HTML (GitHub Pages hostable)
- polling every 60s
- historical uptime persistence (localStorage acceptable)
- mobile-friendly layout

## 3) Technical Design

Frontend-only architecture:
- `index.html`
- `app.js`
- `styles.css`

Node config:
- Node 1: `https://50.28.86.131/health`
- Node 2: `https://50.28.86.153/health`
- Node 3: `http://76.8.228.245:8099/health`

Aux metrics:
- epoch/miners via node 1 endpoints where available

History model:
- localStorage ring-buffer per node
- rolling success/failure samples for 24h/7d/30d windows

Status thresholds:
- green: health ok and response <= 2s
- yellow: health ok but slow (>2s)
- red: timeout/error/unhealthy

## 4) Implementation Plan

### Phase A: Fetch core health (45 min)

- Implement polling engine with timeout + retry.
- Track response duration and status state.

### Phase B: Historical calculations (60 min)

- Persist sampled results in localStorage.
- Compute rolling uptime percentages by window.

### Phase C: Metrics + UI (60 min)

- Add epoch/miner/last-block panels.
- Render card layout and mobile breakpoints.
- Add timestamp and staleness indicators.

### Phase D: QA and evidence (30 min)

- Simulate node failure to validate red state.
- Capture screenshots showing mixed states.
- Prepare submission summary.

## 5) Acceptance Criteria

- [ ] all 3 nodes polled every 60s
- [ ] response time and status indicator visible
- [ ] rolling uptime windows visible
- [ ] epoch/miner/last-block metrics shown
- [ ] mobile layout verified

## 6) Risks and Mitigations

- Mixed schema between nodes:
  - normalize adapters per node.
- CORS/network restrictions on static host:
  - optional tiny proxy backend fallback.
- intermittent external node downtime:
  - classify as status signal, not app failure.

## 7) Submission Package

- hosted page or source bundle
- screenshots with live data and uptime history
- short README notes
- issue comment with proof

## 8) Tracker Commands

```bash
cd /home/administrator/projects/bountyOS/au-workspace
./bin/au bounty move bounty-161-net-status in_progress
./bin/au bounty progress bounty-161-net-status 20
./bin/au bounty next bounty-161-net-status "Implement localStorage uptime window calculations"
```
