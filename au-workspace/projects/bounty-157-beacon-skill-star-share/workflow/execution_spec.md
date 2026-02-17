# Execution Spec: bounty-157-beacon

Issue: https://github.com/Scottcjn/rustchain-bounties/issues/157
Title: Star & Share beacon-skill on GitHub (25 RTC Bounty)
Last updated: 2026-02-14

## 1) Goal

Complete the bounty with verifiable proof in minimum time while preserving account trust signals.

## 2) Deliverables Required by Bounty

- GitHub star on `https://github.com/Scottcjn/beacon-skill`
- One public social/community post mentioning beacon-skill
- Post includes:
  - repo or PyPI/npm link
  - brief description (heartbeat protocol for AI agents)
- Submission comment on issue #157 with:
  - screenshot proving star
  - public post URL
  - RTC wallet

## 3) Execution Plan

### Phase A: Claim + prep (5 min)

- Comment on issue #157 expressing interest.
- Verify you are logged into correct GitHub account.
- Confirm RTC wallet is available for payout.

### Phase B: Complete actions (10-15 min)

- Star beacon-skill repo.
- Publish one post on chosen platform (X/Dev.to/Reddit/etc.).
- Post content pattern:
  - what Beacon does
  - why useful for AI agents
  - link to repo

### Phase C: Proof capture + submission (10 min)

- Capture screenshot showing starred state.
- Copy public post URL.
- Generate final claim text:
  - `./scripts/prepare_claim.sh --post-url "<URL>" --wallet "<RTC_WALLET>"`
- Paste generated content into issue comment.

## 4) Evidence Checklist

- [ ] Star screenshot captured
- [ ] Post URL is public in incognito/private window
- [ ] Issue comment includes wallet
- [ ] Issue comment includes all proof elements

## 5) Risks and Mitigations

- Risk: post not considered sufficient.
  - Mitigation: explicitly describe Beacon as heartbeat/coordination protocol and include repo link.
- Risk: weak trust signal from fresh account.
  - Mitigation: use established account with normal activity history.
- Risk: missing proof artifact.
  - Mitigation: store screenshot and generated claim text in `output/`.

## 6) Definition of Done

- Submission comment posted on issue #157 with all required proof.
- Tracker updated:
  - status `submitted`
  - progress `100`
  - note includes comment timestamp/link.

## 7) Tracker Commands

```bash
cd /home/administrator/projects/bountyOS/au-workspace
./bin/au bounty progress bounty-157-beacon 70
./bin/au bounty move bounty-157-beacon submitted
./bin/au bounty note bounty-157-beacon "Submitted proof comment on issue #157"
```
