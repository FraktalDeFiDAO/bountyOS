# Setup Prompt: bounty-162-wrtc

Set up and validate bounty project in AU workspace.

Workspace root:
`/home/administrator/projects/bountyOS/au-workspace`

Bounty metadata:
- id: `bounty-162-wrtc`
- project_key (registry): `bounty-162-wrtc-price-ticker-bot`
- project_dir (filesystem): `bounty-162-wrtc-price-ticker-bot`
- issue_url: https://github.com/Scottcjn/Rustchain/issues/162
- title: [BOUNTY] wRTC Price Ticker Bot — Discord/Telegram Price Alerts (20 RTC)
- target_status: `ready`
- target_priority: `high`
- target_owner: `unassigned`
- target_progress: `10`
- target_next_action: Pick Discord or Telegram implementation and scaffold bot command

Execution tasks:
1. Ensure registry entry exists in `config/projects.tsv` using key `bounty-162-wrtc-price-ticker-bot`.
2. Ensure project structure exists under:
   - `projects/bounty-162-wrtc-price-ticker-bot/README.md`
   - `projects/bounty-162-wrtc-price-ticker-bot/resources/raw/github/`
   - `projects/bounty-162-wrtc-price-ticker-bot/resources/raw/local/`
   - `projects/bounty-162-wrtc-price-ticker-bot/workflow/`
   - `projects/bounty-162-wrtc-price-ticker-bot/scripts/`
   - `projects/bounty-162-wrtc-price-ticker-bot/output/.gitkeep`
3. Ensure tracker entry exists in `config/bounties.tsv` for id `bounty-162-wrtc`; if missing, add it.
4. Apply tracker values:
   - status `ready`
   - priority `high`
   - owner `unassigned`
   - progress `10`
   - next action as listed above
5. Pull GitHub issue JSON to:
   - `projects/bounty-162-wrtc-price-ticker-bot/resources/raw/github/issue.json`
6. Pull matching local feed entry to:
   - `projects/bounty-162-wrtc-price-ticker-bot/resources/raw/local/feed_item.json`
7. Ensure docs exist:
   - `projects/bounty-162-wrtc-price-ticker-bot/workflow/next_step.md`
   - `projects/bounty-162-wrtc-price-ticker-bot/workflow/execution_spec.md`
8. Append setup note:
   - `./bin/au bounty note bounty-162-wrtc "setup prompt executed"`
9. Verify with:
   - `./bin/au bounty show bounty-162-wrtc`
   - `./bin/au bounty board`
10. Regenerate state docs:
   - `./scripts/generate_current_state.sh`

Return:
- files created/updated
- command outputs summary
- any blockers
