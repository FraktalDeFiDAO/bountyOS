# Execution Spec: bounty-160-beacon-blog

Issue: https://github.com/Scottcjn/rustchain-bounties/issues/160
Title: Write a tutorial or blog about Beacon (50 RTC Bounty)
Last updated: 2026-02-14

## 1) Goal

Publish a high-quality, technically accurate Beacon tutorial (>=500 words) with copy-paste runnable code examples.

## 2) Bounty Requirements (Hard Constraints)

- Published on approved platform (Dev.to, Medium, Hashnode, blog, or Moltbook)
- Explain what Beacon is and why it matters
- Include working Python or JS/TS examples
- Cover at least one feature path:
  - heartbeat
  - mayday
  - contracts
  - Atlas
- Link back to `https://github.com/Scottcjn/beacon-skill`

## 3) Content Architecture

Recommended article structure:

1. Problem framing (agent coordination gaps)
2. What Beacon is
3. Quick install + setup
4. Core walkthrough (heartbeat baseline)
5. Optional second feature (mayday or contracts)
6. Practical usage pattern and failure handling
7. Conclusion + resource links

## 4) Code Example Strategy

- Build one minimal runnable script (heartbeat ping loop).
- Build one extension snippet (mayday or contracts).
- Validate both in clean environment before publishing.
- Include expected output blocks.

## 5) Implementation Plan

### Phase A: Research + outline (30 min)

- Pull latest README/examples from beacon-skill repo.
- Draft outline and choose language (Python recommended).

### Phase B: Example development (60 min)

- Implement minimal working heartbeat example.
- Implement second feature snippet.
- Verify commands end-to-end.

### Phase C: Writing (60-90 min)

- Write article body with clear sections.
- Add architecture diagram (optional bonus).
- Add comparison paragraph to alternative approaches (bonus).

### Phase D: Publish + claim (20 min)

- Publish post.
- Gather URL and screenshot if needed.
- Submit issue comment with proof and wallet.

## 6) Acceptance Criteria

- [ ] 500+ words
- [ ] runnable code examples included
- [ ] at least one required feature covered
- [ ] repo link included
- [ ] published URL submitted in issue

## 7) Risks and Mitigations

- Example drift vs latest package:
  - pin versions and test before publish.
- Post fails verification due to shallow content:
  - include concrete code and explanation of why/when to use.

## 8) Submission Package

- Published URL
- Optional local markdown source in project `output/`
- Issue comment with article URL + wallet

## 9) Tracker Commands

```bash
cd /home/administrator/projects/bountyOS/au-workspace
./bin/au bounty move bounty-160-beacon-blog in_progress
./bin/au bounty progress bounty-160-beacon-blog 25
./bin/au bounty next bounty-160-beacon-blog "Draft heartbeat tutorial and runnable code block"
```
