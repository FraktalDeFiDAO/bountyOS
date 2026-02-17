# Project: bounty-157-beacon-skill-star-share

## Goal

Complete and submit GitHub bounty issue #157:
https://github.com/Scottcjn/rustchain-bounties/issues/157

## Scope

- Star the `beacon-skill` repo
- Publish one public social/community share
- Submit proof + wallet in the issue comment

## Structure

- `resources/raw/github/`: fetched issue/repo/readme snapshots
- `resources/raw/local/`: local feed snapshots from BountyOS API
- `workflow/`: manual execution docs and templates
- `scripts/pull_resources.sh`: sync latest resources into `resources/raw/*`
- `scripts/prepare_claim.sh`: render final claim comment into `output/`
- `output/`: generated claim payloads

## Commands

```bash
# sync resources (issue metadata, repo metadata, readme, local feed)
./scripts/pull_resources.sh

# generate claim text
./scripts/prepare_claim.sh \
  --post-url "https://x.com/<you>/status/<id>" \
  --wallet "<YOUR_RTC_WALLET>"
```

## Execution Spec

- `workflow/execution_spec.md`
