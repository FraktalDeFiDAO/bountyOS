# bountyOS Project Organization Standards

**Version:** 1.0  
**Effective Date:** March 13, 2026  
**Status:** Active

---

## 🎯 Core Principles

### 1. Separation of Concerns
- **Root directory**: ONLY configuration, entry points, and critical documentation
- **Projects**: All active work isolated in `au-workspace/projects/`
- **Agents**: All agent code in `.agents/` with strict modular boundaries
- **Research**: All research, analysis, and knowledge in `au-workspace/research/`
- **Submissions**: All bounty submissions in `submissions/` organized by platform

### 2. Hierarchical Structure
```
bountyOS/
├── .agents/                    # Agent system (isolated from projects)
├── au-workspace/               # Active work and research
│   ├── projects/               # Individual bounty projects
│   ├── research/               # Market research, opportunities
│   ├── departments/            # Cross-project functional areas
│   └── docs/                   # Shared documentation
├── submissions/                # Completed bounty submissions
│   ├── code4rena/              # Code4rena audit findings
│   ├── immunefi/               # Immunefi bug reports
│   ├── superteam/              # Superteam bounties
│   └── ...
├── docs/                       # Project-wide documentation
├── config/                     # Configuration files
├── scripts/                    # Utility scripts
├── tools/                      # Developer tools
└── [source projects]           # Actual codebases (superteam-academy, zio, etc.)
```

### 3. File Naming Conventions

#### Documents
- **UPPER_CASE.md**: Critical status trackers, official reports
- **kebab-case.md**: Regular documentation, guides, plans
- **YYYY-MM-DD-topic.md**: Dated reports, session notes

#### Code Files
- **Rust**: `snake_case.rs` for modules, `PascalCase` for types
- **TypeScript**: `kebab-case.ts` for files, `PascalCase` for components/classes
- **Python**: `snake_case.py` for all files

#### Projects
- **bounty-{id}-{name}**: Standard bounty projects
- **bounty-c4-{target}**: Code4rena audit projects
- **bounty-mps-{id}-{name}**: MPS campaign bounties

---

## 📁 Directory Specifications

### Root Directory (`/`)

**ALLOWED:**
- `README.md` - Project overview
- `AGENTS.md` - Agent system documentation
- `package.json`, `Cargo.toml`, `go.mod` - Root package configs
- `docker-compose*.yml` - Docker configurations
- `.env.example` - Environment template
- `.gitignore`, `.dockerignore` - Git/Docker ignore files
- `.qwen/`, `.claude/`, `.github/` - Tool configurations

**NOT ALLOWED:**
- Status trackers (move to `docs/status/`)
- Session summaries (move to `docs/sessions/`)
- Random scripts (move to `scripts/`)
- Temporary files (move to `out/` or delete)

### Agent Directory (`.agents/`)

```
.agents/
├── orchestrator/               # Central routing agent
├── shared-skills/              # Reusable skills across agents
├── [agent-name]/               # Individual agents
│   ├── agent.md                # Agent definition
│   ├── skills/                 # Agent-specific skills
│   └── tests/                  # Agent tests
└── templates/                  # Agent templates
```

**Rules:**
- Each agent MUST be self-contained
- Agents MUST NOT reference project code directly
- Skills MUST be explicitly declared
- Context isolation MUST be enforced

### Projects Directory (`au-workspace/projects/`)

```
au-workspace/projects/
└── bounty-{id}-{name}/
    ├── index.md                # Project overview
    ├── spec.md                 # Specification
    ├── plan.md                 # Implementation plan
    ├── source/                 # Source code
    ├── output/                 # Generated artifacts
    ├── resources/              # Reference materials
    ├── agent/                  # Project-specific agent (if needed)
    └── tests/                  # Tests
```

**Rules:**
- Each project MUST have an `index.md`
- Source code MUST be in `source/` subdirectory
- Output artifacts MUST be in `output/`
- No project files should leak to root

### Submissions Directory (`submissions/`)

```
submissions/
├── code4rena/
│   └── {contest-name}/
│       ├── HIGH/
│       ├── MEDIUM/
│       └── LOW/
├── immunefi/
│   └── {protocol-name}/
├── superteam/
│   └── {bounty-name}/
└── ...
```

**Rules:**
- Each submission MUST include proof of concept
- Findings MUST be categorized by severity
- Status MUST be tracked in `docs/status/submissions.md`

### Documentation Directory (`docs/`)

```
docs/
├── status/                     # Status trackers
│   ├── bounties.md
│   ├── submissions.md
│   └── active-projects.md
├── sessions/                   # Session notes
├── guides/                     # How-to guides
├── architecture/               # System architecture
├── research/                   # Research summaries
└── standards/                  # Standards and protocols
```

---

## 🔧 Cleanup Protocol

### Phase 1: Root Directory Cleanup
1. Identify all non-compliant files in root
2. Categorize by type (status, session, script, config)
3. Move to appropriate directories
4. Update references and links
5. Delete obsolete files

### Phase 2: Project Isolation
1. Ensure each project has proper structure
2. Move orphaned files to correct projects
3. Create missing `index.md` files
4. Validate no cross-project dependencies

### Phase 3: Agent Modularization
1. Separate global agents from project agents
2. Extract shared skills
3. Document agent interfaces
4. Test agent isolation

### Phase 4: Obsidian Integration
1. Create Obsidian vault structure
2. Import all research and documentation
3. Set up bidirectional linking
4. Create MOCs (Maps of Content)

---

## 📊 Status Tracking

### Active Trackers (Keep in Rotation)
- `BOUNTY_PIPELINE_STATUS.md` → `docs/status/bounties.md`
- `BOUNTY_SUBMISSION_STATUS.md` → `docs/status/submissions.md`
- `IN_PROGRESS_BOUNTIES_STATUS.md` → `docs/status/active-projects.md`

### Archive Trackers (Move to Archive)
- Session summaries → `docs/sessions/`
- Completed project status → `docs/archive/`
- Historical reports → `docs/archive/`

---

## 🎯 Quality Gates

### Before Committing
- [ ] File is in correct directory
- [ ] File follows naming convention
- [ ] No temporary files included
- [ ] Links updated if files moved
- [ ] No sensitive data exposed

### Weekly Maintenance
- [ ] Root directory audit
- [ ] Archive old session notes
- [ ] Update status trackers
- [ ] Clean up `out/` directories
- [ ] Validate project structures

---

## 📝 Migration Checklist

### Immediate Actions
- [ ] Move all `*_STATUS.md` files to `docs/status/`
- [ ] Move all `*_SUMMARY.md` files to `docs/sessions/`
- [ ] Move all `*.py` scripts to `scripts/`
- [ ] Move all `*.ts` standalone files to appropriate projects
- [ ] Create `docs/` directory structure
- [ ] Create `submissions/` directory structure

### Platform Research Integration
- [ ] Create `au-workspace/research/platforms/` directory
- [ ] Import Platform Comparison Matrix as `platforms/comparison-matrix.md`
- [ ] Create `platforms/quick-payout-channels.md` from scored matrix
- [ ] Create `platforms/bounty-platforms-complete.md` from full table
- [ ] Link to Obsidian vault

---

## 🔒 Enforcement

### Automated Checks
- `.github/workflows/structure-check.yml` - Validate directory structure
- `scripts/validate-organization.sh` - Pre-commit hook
- `scripts/cleanup-root.sh` - Automated cleanup

### Manual Reviews
- Weekly structure audit
- Monthly archive sweep
- Quarterly standards review

---

## 📚 Related Documentation

- [AGENTS.md](../AGENTS.md) - Agent system guidelines
- [README.md](../readme.md) - Project overview
- [docs/architecture/](./architecture/) - System architecture

---

**Last Updated:** March 13, 2026  
**Maintained By:** bountyOS Core Team
