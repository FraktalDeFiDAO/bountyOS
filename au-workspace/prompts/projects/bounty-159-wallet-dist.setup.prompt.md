# Setup Prompt: bounty-159-wallet-dist

Set up and validate bounty project in AU workspace.

Workspace root:
`/home/administrator/projects/bountyOS/au-workspace`

Bounty metadata:
- id: `bounty-159-wallet-dist`
- project_key (registry): `bounty-159-rtc-wallet-distribution`
- project_dir (filesystem): `bounty-159-rtc-wallet-distribution`
- issue_url: https://github.com/Scottcjn/Rustchain/issues/159
- title: [BOUNTY] RTC Wallet Distribution Tracker — Live Balance Dashboard (40 RTC)
- target_status: `ready`
- target_priority: `high`
- target_owner: `unassigned`
- target_progress: `0`
- target_next_action: Define API polling model and holder distribution visualization

Execution tasks:
1. Ensure registry entry exists in `config/projects.tsv` using key `bounty-159-rtc-wallet-distribution`.
2. Ensure project structure exists under:
   - `projects/bounty-159-rtc-wallet-distribution/README.md`
   - `projects/bounty-159-rtc-wallet-distribution/resources/raw/github/`
   - `projects/bounty-159-rtc-wallet-distribution/resources/raw/local/`
   - `projects/bounty-159-rtc-wallet-distribution/workflow/`
   - `projects/bounty-159-rtc-wallet-distribution/scripts/`
   - `projects/bounty-159-rtc-wallet-distribution/output/.gitkeep`
3. Ensure tracker entry exists in `config/bounties.tsv` for id `bounty-159-wallet-dist`; if missing, add it.
4. Apply tracker values:
   - status `ready`
   - priority `high`
   - owner `unassigned`
   - progress `0`
   - next action as listed above
5. Pull GitHub issue JSON to:
   - `projects/bounty-159-rtc-wallet-distribution/resources/raw/github/issue.json`
6. Pull matching local feed entry to:
   - `projects/bounty-159-rtc-wallet-distribution/resources/raw/local/feed_item.json`
7. Ensure docs exist:
   - `projects/bounty-159-rtc-wallet-distribution/workflow/next_step.md`
   - `projects/bounty-159-rtc-wallet-distribution/workflow/execution_spec.md`
8. Append setup note:
   - `./bin/au bounty note bounty-159-wallet-dist "setup prompt executed"`
9. Verify with:
   - `./bin/au bounty show bounty-159-wallet-dist`
   - `./bin/au bounty board`
10. Regenerate state docs:
   - `./scripts/generate_current_state.sh`

Return:
- files created/updated
- command outputs summary
- any blockers
