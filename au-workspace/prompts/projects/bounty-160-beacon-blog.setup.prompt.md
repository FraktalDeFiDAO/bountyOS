# Setup Prompt: bounty-160-beacon-blog

Set up and validate bounty project in AU workspace.

Workspace root:
`/home/administrator/projects/bountyOS/au-workspace`

Bounty metadata:
- id: `bounty-160-beacon-blog`
- project_key (registry): `bounty-160-beacon-blog-tutorial`
- project_dir (filesystem): `bounty-160-beacon-blog-tutorial`
- issue_url: https://github.com/Scottcjn/rustchain-bounties/issues/160
- title: Write a tutorial or blog about Beacon (50 RTC Bounty)
- target_status: `backlog`
- target_priority: `medium`
- target_owner: `unassigned`
- target_progress: `0`
- target_next_action: Choose publication platform and outline Beacon 2.6 walkthrough

Execution tasks:
1. Ensure registry entry exists in `config/projects.tsv` using key `bounty-160-beacon-blog-tutorial`.
2. Ensure project structure exists under:
   - `projects/bounty-160-beacon-blog-tutorial/README.md`
   - `projects/bounty-160-beacon-blog-tutorial/resources/raw/github/`
   - `projects/bounty-160-beacon-blog-tutorial/resources/raw/local/`
   - `projects/bounty-160-beacon-blog-tutorial/workflow/`
   - `projects/bounty-160-beacon-blog-tutorial/scripts/`
   - `projects/bounty-160-beacon-blog-tutorial/output/.gitkeep`
3. Ensure tracker entry exists in `config/bounties.tsv` for id `bounty-160-beacon-blog`; if missing, add it.
4. Apply tracker values:
   - status `backlog`
   - priority `medium`
   - owner `unassigned`
   - progress `0`
   - next action as listed above
5. Pull GitHub issue JSON to:
   - `projects/bounty-160-beacon-blog-tutorial/resources/raw/github/issue.json`
6. Pull matching local feed entry to:
   - `projects/bounty-160-beacon-blog-tutorial/resources/raw/local/feed_item.json`
7. Ensure docs exist:
   - `projects/bounty-160-beacon-blog-tutorial/workflow/next_step.md`
   - `projects/bounty-160-beacon-blog-tutorial/workflow/execution_spec.md`
8. Append setup note:
   - `./bin/au bounty note bounty-160-beacon-blog "setup prompt executed"`
9. Verify with:
   - `./bin/au bounty show bounty-160-beacon-blog`
   - `./bin/au bounty board`
10. Regenerate state docs:
   - `./scripts/generate_current_state.sh`

Return:
- files created/updated
- command outputs summary
- any blockers
