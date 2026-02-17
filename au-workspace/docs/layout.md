# Layout Rationale

## Why this structure exists

- `bin/`: executable entrypoints only; keeps automation discoverable.
- `config/`: mutable workspace state and conventions in one place.
- `config/projects.tsv`: registered project map for `au run/exec/status`.
- `config/bounties.tsv`: progress tracker for all active bounties.
- `projects/`: each bounty/project is isolated to reduce cross-task spillover.
- `projects/<slug>/resources/raw/`: immutable pulled source data (API responses, READMEs, snapshots).
- `projects/<slug>/workflow/`: human task flow (checklists/templates) tied to one bounty.
- `projects/<slug>/scripts/`: repeatable operations (sync resources, generate claim text).
- `projects/<slug>/output/`: generated artifacts safe to overwrite.
- `projects/<slug>/output/progress_notes.md`: timestamped notes written by `au bounty note`.
- `scripts/`: workspace-level maintenance scripts (example: docs state snapshot generator).
- `templates/`: reusable seed content for future projects.

## Naming rules

- Project slug: `domain-issue-short-name`
- Script names: `verb_object.sh` (example: `pull_resources.sh`)
- Template files: lowercase snake/case with clear suffixes (`issue_comment.template.md`)
- Bounty statuses: `backlog, ready, in_progress, blocked, review, submitted, paid, done, dropped`
