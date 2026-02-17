# Execution Spec: bounty-163-leaderboard

Issue: https://github.com/Scottcjn/Rustchain/issues/163
Title: [BOUNTY] Miner Leaderboard — Top Earners & Hardware Showcase (20 RTC)
Last updated: 2026-02-14

## 1) Goal

Ship a live miner leaderboard with ranking, estimated earnings, hardware highlights, and network summary cards.

## 2) Required Output

Leaderboard table fields:
- rank
- miner
- hardware
- multiplier
- epochs active
- estimated earnings

Stats cards:
- total miners
- total RTC distributed
- most exotic hardware
- newest miner
- architecture diversity

Implementation constraints:
- single HTML page
- live auto-refresh
- sort support
- hardware badges/icons

## 3) Technical Design

Data sources:
- `https://50.28.86.131/api/miners`
- `https://50.28.86.131/epoch`

Earnings model:
- infer epochs active from attest timing and current epoch where possible
- estimated earnings formula from issue guidance:
  - share approx proportional to antiquity multiplier vs network sum

Processing pipeline:
1. fetch miners + epoch
2. normalize miner records
3. compute derived metrics (rank score, earnings estimate)
4. render sortable table + stat cards

## 4) Implementation Plan

### Phase A: Data adapter and ranking (45 min)

- Parse miner payload, multipliers, hardware labels.
- Implement ranking by estimated earnings.

### Phase B: Earnings estimation and stats (60 min)

- Compute network multiplier sum.
- Estimate per-miner reward share and aggregate period values.
- Build architecture diversity and newest miner metrics.

### Phase C: UI and interaction (60 min)

- Build table with sort toggles.
- Add badge styles:
  - vintage = gold
  - modern = silver
- Add refresh cadence and last-updated indicator.

### Phase D: QA + proof (30 min)

- Validate sorting and metric consistency.
- Capture screenshot with populated leaderboard and cards.

## 5) Acceptance Criteria

- [ ] leaderboard shows live miner data
- [ ] sorting works for key columns
- [ ] estimated earnings column is populated
- [ ] stats cards render from live metrics
- [ ] hardware badges distinguish classes

## 6) Risks and Mitigations

- Ambiguous epoch/attest mapping:
  - document assumptions clearly in UI tooltip.
- Missing fields in API response:
  - provide safe fallbacks and mark unknown values.
- Overly precise earnings claims:
  - label values as estimates.

## 7) Submission Package

- source files
- screenshot of leaderboard + cards
- short formula/assumption note in README
- issue comment with evidence

## 8) Tracker Commands

```bash
cd /home/administrator/projects/bountyOS/au-workspace
./bin/au bounty move bounty-163-leaderboard in_progress
./bin/au bounty progress bounty-163-leaderboard 20
./bin/au bounty next bounty-163-leaderboard "Implement earnings estimate formula and sorting"
```
