# Setup Prompt: bounty-157-beacon

Set up and validate bounty project in AU workspace.

Workspace root:
`/home/administrator/projects/bountyOS/au-workspace`

Bounty metadata:
- id: `bounty-157-beacon`
- project_key (registry): `bounty-157-beacon`
- project_dir (filesystem): `bounty-157-beacon-skill-star-share`
- issue_url: https://github.com/Scottcjn/rustchain-bounties/issues/157
- title: Star & Share beacon-skill on GitHub (25 RTC Bounty)
- target_status: `in_progress`
- target_priority: `high`
- target_owner: `unassigned`
- target_progress: `30`
- target_next_action: Star repo, post share, submit issue comment

Execution tasks:
1. Ensure registry entry exists in `config/projects.tsv` using key `bounty-157-beacon`.
2. Ensure project structure exists under:
   - `projects/bounty-157-beacon-skill-star-share/README.md`
   - `projects/bounty-157-beacon-skill-star-share/resources/raw/github/`
   - `projects/bounty-157-beacon-skill-star-share/resources/raw/local/`
   - `projects/bounty-157-beacon-skill-star-share/workflow/`
   - `projects/bounty-157-beacon-skill-star-share/scripts/`
   - `projects/bounty-157-beacon-skill-star-share/output/.gitkeep`
3. Ensure tracker entry exists in `config/bounties.tsv` for id `bounty-157-beacon`; if missing, add it.
4. Apply tracker values:
   - status `in_progress`
   - priority `high`
   - owner `unassigned`
   - progress `30`
   - next action as listed above
5. Pull GitHub issue JSON to:
   - `projects/bounty-157-beacon-skill-star-share/resources/raw/github/issue.json`
6. Pull matching local feed entry to:
   - `projects/bounty-157-beacon-skill-star-share/resources/raw/local/feed_item.json`
7. Ensure docs exist:
   - `projects/bounty-157-beacon-skill-star-share/workflow/next_step.md`
   - `projects/bounty-157-beacon-skill-star-share/workflow/execution_spec.md`
8. Append setup note:
   - `./bin/au bounty note bounty-157-beacon "setup prompt executed"`
9. Verify with:
   - `./bin/au bounty show bounty-157-beacon`
   - `./bin/au bounty board`
10. Regenerate state docs:
   - `./scripts/generate_current_state.sh`

Return:
- files created/updated
- command outputs summary
- any blockers
