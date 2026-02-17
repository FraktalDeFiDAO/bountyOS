# Setup Prompt: bounty-161-net-status

Set up and validate bounty project in AU workspace.

Workspace root:
`/home/administrator/projects/bountyOS/au-workspace`

Bounty metadata:
- id: `bounty-161-net-status`
- project_key (registry): `bounty-161-network-status-dashboard`
- project_dir (filesystem): `bounty-161-network-status-dashboard`
- issue_url: https://github.com/Scottcjn/Rustchain/issues/161
- title: [BOUNTY] RustChain Network Status Page — Public Health Dashboard (25 RTC)
- target_status: `ready`
- target_priority: `medium`
- target_owner: `unassigned`
- target_progress: `0`
- target_next_action: Build single-page status dashboard with uptime tracking

Execution tasks:
1. Ensure registry entry exists in `config/projects.tsv` using key `bounty-161-network-status-dashboard`.
2. Ensure project structure exists under:
   - `projects/bounty-161-network-status-dashboard/README.md`
   - `projects/bounty-161-network-status-dashboard/resources/raw/github/`
   - `projects/bounty-161-network-status-dashboard/resources/raw/local/`
   - `projects/bounty-161-network-status-dashboard/workflow/`
   - `projects/bounty-161-network-status-dashboard/scripts/`
   - `projects/bounty-161-network-status-dashboard/output/.gitkeep`
3. Ensure tracker entry exists in `config/bounties.tsv` for id `bounty-161-net-status`; if missing, add it.
4. Apply tracker values:
   - status `ready`
   - priority `medium`
   - owner `unassigned`
   - progress `0`
   - next action as listed above
5. Pull GitHub issue JSON to:
   - `projects/bounty-161-network-status-dashboard/resources/raw/github/issue.json`
6. Pull matching local feed entry to:
   - `projects/bounty-161-network-status-dashboard/resources/raw/local/feed_item.json`
7. Ensure docs exist:
   - `projects/bounty-161-network-status-dashboard/workflow/next_step.md`
   - `projects/bounty-161-network-status-dashboard/workflow/execution_spec.md`
8. Append setup note:
   - `./bin/au bounty note bounty-161-net-status "setup prompt executed"`
9. Verify with:
   - `./bin/au bounty show bounty-161-net-status`
   - `./bin/au bounty board`
10. Regenerate state docs:
   - `./scripts/generate_current_state.sh`

Return:
- files created/updated
- command outputs summary
- any blockers
