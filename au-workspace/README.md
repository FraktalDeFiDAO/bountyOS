# AU Workspace (Multi-Project)

This is a structured workspace for running multiple projects and bounty tracks from one command surface.

## What changed

- Centralized config in `config/`
- All active work now lives in `projects/`
- Bounty #157 has a dedicated project: `projects/bounty-157-beacon-skill-star-share/`
- Added repeatable resource sync script and claim generator
- Added bounty progress tracking in `config/bounties.tsv` with `au bounty ...` commands

## Directory Map

- `bin/au`: multi-project CLI (`list`, `add`, `new`, `run`, `exec`, `status`)
- `config/workspace.yaml`: workspace conventions
- `config/projects.tsv`: project registry used by CLI
- `config/bounties.tsv`: bounty progress tracker used by CLI
- `docs/layout.md`: why each folder exists
- `docs/progress.md`: bounty tracker usage and status model
- `docs/README.md`: full docs index
- `docs/complete-reference.md`: complete command/config/runbook reference
- `docs/current-state.md`: current inventory snapshot and tracked bounties
- `prompts/`: orchestration and project-setup prompt pack
- `scripts/generate_current_state.sh`: regenerates `docs/current-state.md`
- `projects/`: one folder per active project/bounty
- `templates/`: starter templates for new projects

## Quickstart

```bash
cd au-workspace
./bin/au paths
./bin/au list
./bin/au status
./bin/au bounty list
```

### Add another project

```bash
./bin/au add my-repo /absolute/path/to/repo
```

### Run one command across all projects

```bash
./bin/au run git status -sb
```

### Track multiple bounties in parallel

```bash
# add a tracked bounty
./bin/au bounty add bounty-157-beacon bounty-157-beacon https://github.com/Scottcjn/rustchain-bounties/issues/157 \
  --title "Star & Share beacon-skill on GitHub (25 RTC Bounty)" \
  --status in_progress --priority high --owner me --progress 25 \
  --next "Post social proof and submit claim"

# see grouped board
./bin/au bounty board

# update progress and workflow state
./bin/au bounty progress bounty-157-beacon 60
./bin/au bounty move bounty-157-beacon review
./bin/au bounty note bounty-157-beacon "Waiting on screenshot upload"
```

### Regenerate state documentation

```bash
./scripts/generate_current_state.sh
```
