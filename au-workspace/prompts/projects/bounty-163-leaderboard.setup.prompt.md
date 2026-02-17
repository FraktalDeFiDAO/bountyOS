# Setup Prompt: bounty-163-leaderboard

Set up and validate bounty project in AU workspace.

Workspace root:
`/home/administrator/projects/bountyOS/au-workspace`

Bounty metadata:
- id: `bounty-163-leaderboard`
- project_key (registry): `bounty-163-miner-leaderboard`
- project_dir (filesystem): `bounty-163-miner-leaderboard`
- issue_url: https://github.com/Scottcjn/Rustchain/issues/163
- title: [BOUNTY] Miner Leaderboard — Top Earners & Hardware Showcase (20 RTC)
- target_status: `backlog`
- target_priority: `medium`
- target_owner: `unassigned`
- target_progress: `0`
- target_next_action: Design leaderboard metrics and fetch flow from miner endpoints

Execution tasks:
1. Ensure registry entry exists in `config/projects.tsv` using key `bounty-163-miner-leaderboard`.
2. Ensure project structure exists under:
   - `projects/bounty-163-miner-leaderboard/README.md`
   - `projects/bounty-163-miner-leaderboard/resources/raw/github/`
   - `projects/bounty-163-miner-leaderboard/resources/raw/local/`
   - `projects/bounty-163-miner-leaderboard/workflow/`
   - `projects/bounty-163-miner-leaderboard/scripts/`
   - `projects/bounty-163-miner-leaderboard/output/.gitkeep`
3. Ensure tracker entry exists in `config/bounties.tsv` for id `bounty-163-leaderboard`; if missing, add it.
4. Apply tracker values:
   - status `backlog`
   - priority `medium`
   - owner `unassigned`
   - progress `0`
   - next action as listed above
5. Pull GitHub issue JSON to:
   - `projects/bounty-163-miner-leaderboard/resources/raw/github/issue.json`
6. Pull matching local feed entry to:
   - `projects/bounty-163-miner-leaderboard/resources/raw/local/feed_item.json`
7. Ensure docs exist:
   - `projects/bounty-163-miner-leaderboard/workflow/next_step.md`
   - `projects/bounty-163-miner-leaderboard/workflow/execution_spec.md`
8. Append setup note:
   - `./bin/au bounty note bounty-163-leaderboard "setup prompt executed"`
9. Verify with:
   - `./bin/au bounty show bounty-163-leaderboard`
   - `./bin/au bounty board`
10. Regenerate state docs:
   - `./scripts/generate_current_state.sh`

Return:
- files created/updated
- command outputs summary
- any blockers
