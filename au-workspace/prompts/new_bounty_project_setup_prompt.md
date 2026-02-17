# Prompt Template: New Bounty Project Setup

Use this prompt to set up one bounty project end-to-end.

---

Set up this bounty project in AU workspace.

Workspace root:
`/home/administrator/projects/bountyOS/au-workspace`

Bounty metadata:
- `id`: <bounty-id>
- `project_key`: <project-name-in-config-projects.tsv>
- `project_dir`: <folder-name-under-projects/>
- `issue_url`: <github-issue-url>
- `title`: <bounty-title>
- `priority`: <low|medium|high|critical>
- `status`: <backlog|ready|in_progress|blocked|review|submitted|paid|done|dropped>
- `progress`: <0-100>
- `owner`: <owner-name>
- `next_action`: <next-concrete-step>

Required tasks:
1. Ensure project key is registered (`./bin/au add` or `./bin/au new` using `project_key`).
2. Ensure directory structure exists under `projects/<project_dir>/`:
   - `README.md`
   - `resources/raw/github/`
   - `resources/raw/local/`
   - `workflow/`
   - `scripts/`
   - `output/.gitkeep`
3. Add/update tracker entry with `./bin/au bounty ...` commands.
4. Pull GitHub issue data into `resources/raw/github/issue.json`.
5. Pull local feed item into `resources/raw/local/feed_item.json`.
6. Create/update:
   - `workflow/next_step.md`
   - `workflow/execution_spec.md`
7. Append setup note using:
   - `./bin/au bounty note <id> "project setup completed"`
8. Verify:
   - `./bin/au bounty show <id>`
   - `./bin/au bounty board`
9. Regenerate docs:
   - `./scripts/generate_current_state.sh`

Return a completion report with exact file paths touched and command outcomes.
