# Execution Spec: bounty-159-wallet-dist

Issue: https://github.com/Scottcjn/Rustchain/issues/159
Title: [BOUNTY] RTC Wallet Distribution Tracker — Live Balance Dashboard (40 RTC)
Last updated: 2026-02-14

## 1) Goal

Build a live wallet distribution tracker with concentration analytics and whale detection, designed for easy deployment.

## 2) Scope and Payout Mapping

Base deliverable (20 RTC):
- All non-zero wallets
- balance and % supply
- sorted holder table

Second tier (10 RTC):
- distribution chart
- whale alerts (>1% supply)
- Gini coefficient

Final tier (10 RTC):
- clean UI
- founder wallet labels
- auto-refresh every 5 minutes

## 3) Recommended Technical Approach

Use Option A (single HTML + JS) for fastest shipping and easiest hosting.

Architecture:
- `index.html`
- `app.js`
- `styles.css`

Data pipeline:
- Fetch `https://50.28.86.131/api/miners`
- Normalize wallet records
- Compute:
  - non-zero holders
  - % of total supply (8,300,000 RTC)
  - whale flag (`balance > 83000`)
  - Gini coefficient

Visualization:
- Top holders table
- Concentration bar/pie chart
- Whale alerts panel
- Metadata card (last refresh time)

## 4) Implementation Plan

### Phase A: Data ingestion (45 min)

- Build fetch with timeout and retry.
- Extract wallet id + balance fields.
- Filter non-zero balances.

### Phase B: Analytics (45 min)

- % of supply calculation.
- Gini implementation and unit-style sanity checks.
- Whale threshold flagging logic.

### Phase C: UI and refresh loop (60 min)

- Render holders table and chart.
- Add founder wallet labels (known ids when available).
- Implement 5-minute refresh timer and visible timestamp.

### Phase D: QA + proof (30 min)

- Validate totals against raw API sample.
- Capture screenshot of dashboard.
- Prepare usage notes and submit.

## 5) Acceptance Criteria

- [ ] Wallet list populates from live endpoint.
- [ ] % supply and totals are internally consistent.
- [ ] Whale alerts displayed for >1% wallets.
- [ ] Gini metric displayed and numerically plausible.
- [ ] Auto-refresh works without page reload.

## 6) Risks and Mitigations

- API field differences:
  - use schema normalization with safe defaults.
- CORS for static hosting:
  - if blocked, use lightweight backend proxy variant.
- Founder labels unknown:
  - keep optional config map and show fallback label.

## 7) Submission Package

- Source files (single-page app)
- screenshot with live data
- short README/run notes
- issue comment with proof

## 8) Tracker Commands

```bash
cd /home/administrator/projects/bountyOS/au-workspace
./bin/au bounty move bounty-159-wallet-dist in_progress
./bin/au bounty progress bounty-159-wallet-dist 35
./bin/au bounty next bounty-159-wallet-dist "Implement Gini + whale detection"
```
