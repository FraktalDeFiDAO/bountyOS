# Complete Reference

Last updated: 2026-02-14 UTC

## 1) Purpose

`au-workspace` is a multi-project operations layer for running several bounty tracks in parallel with:
- One command interface (`bin/au`)
- One project registry (`config/projects.tsv`)
- One bounty tracker (`config/bounties.tsv`)
- One isolated folder per bounty in `projects/`

## 2) Top-Level Architecture

```
au-workspace/
  bin/
    au
  config/
    workspace.yaml
    projects.tsv
    bounties.tsv
  docs/
  projects/
    <bounty-or-project-slug>/
      README.md
      resources/
      workflow/
      scripts/
      output/
  templates/
```

## 3) Naming Conventions

- Project slug format: `domain-issue-short-name`
- Bounty tracker id format: lowercase slug, e.g. `bounty-157-beacon`
- Script naming format: `verb_object.sh`
- Template naming format: lowercase plus suffixes, e.g. `issue_comment.template.md`

## 4) Configuration Contracts

## 4.1 `config/workspace.yaml`

Defines workspace conventions:
- `name`
- `version`
- `registry`
- `tracker`
- `projects_root`
- `conventions` (status list, priority list, path conventions)

## 4.2 `config/projects.tsv`

TSV schema:
1. `project_name`
2. `absolute_path`

Rules:
- `project_name` must match `^[a-z0-9][a-z0-9._-]*$`
- Paths must exist at registration time

## 4.3 `config/bounties.tsv`

TSV schema:
1. `id`
2. `issue_url`
3. `title`
4. `project`
5. `status`
6. `priority`
7. `owner`
8. `progress` (0-100)
9. `next_action`
10. `blocker`
11. `created_at` (UTC RFC3339)
12. `updated_at` (UTC RFC3339)

Allowed status values:
- `backlog`
- `ready`
- `in_progress`
- `blocked`
- `review`
- `submitted`
- `paid`
- `done`
- `dropped`

Allowed priority values:
- `low`
- `medium`
- `high`
- `critical`

## 5) CLI Reference (`bin/au`)

## 5.1 Core workspace commands

`./bin/au list`
- Lists registered projects with `ok` or `missing` state.

`./bin/au add <name> <path>`
- Registers an existing directory as a project.

`./bin/au remove <name>`
- Removes a project registration entry (does not delete files).

`./bin/au new <name>`
- Creates a new project scaffold under `projects/<name>` and auto-registers it.

`./bin/au run <command...>`
- Runs the same shell command across every registered project.
- Exit behavior:
  - `1` when no projects exist
  - `2` if one or more project runs fail

`./bin/au exec <name> <command...>`
- Runs a command in one project only.

`./bin/au status`
- Alias workflow for `run git status -sb`.

`./bin/au paths`
- Prints resolved workspace paths (`workspace_root`, `config_dir`, `registry_file`, `tracker_file`, `projects_dir`).

## 5.2 Bounty tracker commands

`./bin/au bounty list`
- Tabular list of tracked bounties with status/progress/priority/owner/next action.

`./bin/au bounty add <id> <project> <issue_url> --title "<title>" [options]`
- Creates a bounty record.
- Optional flags:
  - `--status`
  - `--priority`
  - `--owner`
  - `--progress`
  - `--next`
  - `--blocker`

`./bin/au bounty show <id>`
- Full detail view including resolved `project_path`.

`./bin/au bounty move <id> <status>`
- Updates lifecycle status.

`./bin/au bounty progress <id> <0-100>`
- Updates numeric progress.

`./bin/au bounty next <id> <text...>`
- Updates `next_action`.

`./bin/au bounty block <id> <text...>`
- Updates `blocker`.

`./bin/au bounty owner <id> <owner>`
- Updates owner assignment.

`./bin/au bounty priority <id> <level>`
- Updates priority.

`./bin/au bounty note <id> <text...>`
- Appends a timestamped note to `<project>/output/progress_notes.md`.

`./bin/au bounty board`
- Grouped board view by status.

## 6) Project Folder Contract

Each project in `projects/<slug>/` should use:
- `README.md`: concise project summary and entry commands
- `resources/raw/`: pulled source artifacts (API snapshots, docs snapshots)
- `workflow/`: checklists, templates, manual steps, and technical execution plan
- `workflow/execution_spec.md`: deep implementation plan for that specific bounty
- `scripts/`: repeatable scripts for this project
- `output/`: generated claim text, notes, reports

## 7) Bounty Operations Runbook

## 7.1 Start a new bounty

1. Create project scaffold:
   - `./bin/au new <project-slug>`
2. Add tracker entry:
   - `./bin/au bounty add <id> <project-slug> <issue_url> --title "<title>" --status backlog --priority medium --progress 0 --next "<next step>"`
3. Populate project README/workflow/resources.
4. Move to `ready` when actionable:
   - `./bin/au bounty move <id> ready`

## 7.2 Execute a bounty

1. Move to in progress:
   - `./bin/au bounty move <id> in_progress`
2. Assign owner:
   - `./bin/au bounty owner <id> <owner>`
3. Update progress and notes frequently:
   - `./bin/au bounty progress <id> <n>`
   - `./bin/au bounty note <id> "<update>"`
4. Keep `next_action` current:
   - `./bin/au bounty next <id> "<next concrete step>"`

## 7.3 Submit and close

1. Move to review/submitted:
   - `./bin/au bounty move <id> review`
   - `./bin/au bounty move <id> submitted`
2. If payout confirmed:
   - `./bin/au bounty move <id> paid`
   - `./bin/au bounty progress <id> 100`
3. Mark done:
   - `./bin/au bounty move <id> done`

## 8) Parallel Work Policy

Recommended limits:
- `in_progress`: keep 1-3 active per operator
- `ready`: maintain queue of 3-10
- `blocked`: include blocker text and next unblock action

Recommended cadences:
- Board review at start of session: `./bin/au bounty board`
- End-of-session update:
  - `progress`
  - `next`
  - `note`

## 9) Existing Project-Specific Automation

`projects/bounty-157-beacon-skill-star-share/scripts/pull_resources.sh`
- Pulls:
  - issue JSON
  - issue comments JSON
  - target repo metadata JSON
  - target repo README
  - local bounty feed snapshot (best-effort)
- Writes:
  - `resources/raw/github/*`
  - `resources/raw/local/feed_bounties.json`
  - `resources/summary.md`

`projects/bounty-157-beacon-skill-star-share/scripts/prepare_claim.sh`
- Input:
  - `--post-url`
  - `--wallet`
  - optional `--proof-note`
  - optional `--out`
- Output:
  - claim comment text file (default `output/claim_comment_filled.md`)

## 10) Prompt Pack

Prompt files for consistent setup execution:
- `prompts/orchestration_prompt.md`: batch orchestration across multiple bounties.
- `prompts/new_bounty_project_setup_prompt.md`: reusable single-bounty setup template.
- `prompts/projects/*.setup.prompt.md`: pre-filled prompts for currently tracked bounties.

Recommended usage order:
1. Use orchestration prompt for bulk setup.
2. Use project-specific prompt for one bounty follow-through.
3. Use template prompt for new bounties not yet tracked.

## 11) Error Handling and Troubleshooting

Invalid project name:
- Cause: name does not match allowed slug pattern.
- Fix: use lowercase letters/numbers/dot/underscore/dash.

Unknown project on `bounty add`:
- Cause: `project` field not in `config/projects.tsv`.
- Fix: register project first using `./bin/au add ...` or create with `./bin/au new ...`.

Invalid status/priority/progress:
- Cause: value outside allowed enums/range.
- Fix: use values documented above.

Missing project path in `bounty show`:
- Cause: project entry exists but path no longer exists.
- Fix: repair path or remove/re-add project registration.

Resource sync failure in `pull_resources.sh`:
- Cause: network/API availability or rate limit.
- Fix: rerun; for private APIs set credentials externally if needed.

## 12) Backup and Recovery

Critical state files:
- `config/projects.tsv`
- `config/bounties.tsv`

Backup suggestion:
- Commit these files regularly or copy them to a separate backup location.

Recovery:
1. Restore both TSV files.
2. Run `./bin/au list` and `./bin/au bounty board` to validate integrity.

## 13) Documentation Maintenance

Regenerate live inventory docs after tracker/project changes:

```bash
./scripts/generate_current_state.sh
```

This rewrites:
- `docs/current-state.md`

## 14) Security and Hygiene

- Do not store secrets in tracker TSV files.
- Keep API tokens in environment variables.
- Treat `resources/raw/` as source snapshots; do not hand-edit unless intentional.
- Keep generated artifacts in `output/` to avoid mixing authored docs with generated text.
