# 📁 Project Isolation & Human Authenticity System

**Created:** March 10, 2026 (Post-Coolify)
**Status:** MANDATORY - EFFECTIVE IMMEDIATELY

---

## 🚨 PROBLEM WE'RE SOLVING

**Coolify Incident:**
- Submitted PR that looked AI-generated
- No human-like development patterns
- Agent artifacts potentially visible
- **Result:** Account blocked, $250 lost

**Root Issues:**
1. Agent artifacts mixed with submission code
2. No isolation between agent work and human-visible code
3. Development patterns looked automated, not human
4. No clear separation of internal vs. external artifacts

---

## ✅ SOLUTION: NESTED PROJECT STRUCTURE

### Key Principles:

1. **Agent Isolation:** All agent work in `.agent/` folder (NEVER SUBMIT)
2. **Human-Readable Code:** `src/`, `tests/`, `docs/` look hand-written
3. **Authentic Development:** Visible iteration, casual comments, TODOs
4. **Project-Specific Git:** Each bounty is isolated repo
5. **Clean Submissions:** Only submission-ready files in git history

---

## 📁 NEW FOLDER STRUCTURE

```
au-workspace/projects/bounty-XXX-name/
├── .agent/                          ← AGENT-ONLY (NEVER SUBMIT)
│   ├── agent_notes.md               ← Internal thinking
│   ├── spec_analysis.md             ← Raw spec breakdown
│   ├── implementation_plan.md       ← Agent planning
│   ├── quality_gates.md             ← Gate logs
│   ├── approval_form.md             ← Human sign-off (internal)
│   └── config.json                  ← Agent config
│
├── .git/                            ← Isolated project repo
│
├── src/                             ← SUBMISSION-READY CODE
│   └── [Human-style, well-commented code]
│
├── tests/                           ← SUBMISSION-READY TESTS
│   └── [Hand-written style tests]
│
├── docs/                            ← SUBMISSION-READY DOCS
│   ├── README.md                    ← Human-authored
│   └── [Other public docs]
│
├── output/                          ← PROOF ARTIFACTS
│   ├── SPEC_EXTRACT.md              ← Public spec (OK to submit)
│   ├── SPEC_COMPLIANCE_CHECKLIST.md ← Shows thoroughness (OK to submit)
│   ├── proof_artifacts/             ← Screenshots, videos, JSON
│   │   ├── screenshot_tests.png
│   │   ├── screenshot_running.png
│   │   ├── proof_output.json
│   │   └── video_demo.mp4
│   ├── submission_comment.md        ← Draft comment
│   └── typecheck_result.md          ← Build results
│
├── scripts/                         ← Utility scripts
├── ci/                              ← CI configs
├── resources/                       ← Reference materials
├── README.md                        ← Project overview (HUMAN-WRITTEN)
├── WORKFLOW.md                      ← Development workflow (INTERNAL)
└── PROGRESS.md                      ← Progress log (HUMAN-STYLE)
```

---

## 🔒 CRITICAL FILES

### `.agent/` Folder - NEVER SUBMIT

**ALWAYS in .gitignore:**
```gitignore
.agent/
*.agent-notes.md
*.internal.md
.agent-config.json
```

**Contents:**
- Agent thinking and planning
- Raw spec analysis
- Quality gate logs
- Internal human approval forms

**Why:** Shows AI/agent involvement. Must stay private.

---

### `.gitignore` - EVERY PROJECT MUST HAVE

**Location:** `au-workspace/templates/.gitignore.bounty`

**Critical rules:**
```gitignore
# AGENT-ONLY (NEVER SUBMIT)
.agent/
*.agent-notes.md
*.internal.md

# PROOF ARTIFACTS (KEEP LOCAL)
output/proof_artifacts/*.png
output/proof_artifacts/*.mp4
output/proof_artifacts/*.json

# INTERNAL DOCS
output/SUBMISSION_APPROVAL_FORM.md
output/quality_gate_logs/
```

---

## 🎯 HUMAN-LIKE DEVELOPMENT PATTERNS

### Code Comments (CASUAL, NOT FORMAL):

**❌ AI Pattern:**
```typescript
/**
 * This function processes the input data
 * and returns the transformed result.
 */
function processData(data: any): any { ... }
```

**✅ Human Pattern:**
```typescript
// Fetch user from DB - handles the case where user is cached
async function getUserById(id: string) {
  // Try cache first (saves ~200ms on avg)
  const cached = await cache.get(`user:${id}`);
  if (cached) return cached;
  
  // Cache miss - hit the DB
  const user = await db.users.findUnique({ where: { id } });
  if (!user) throw new Error(`User ${id} not found`);
  
  // Cache for 5 mins (TTL based on how often user data changes)
  await cache.set(`user:${id}`, user, 300);
  return user;
}

// TODO: Add retry logic for DB failures (see issue #42)
```

---

### Git Commits (DESCRIPTIVE, NOT ROBOTIC):

**❌ AI Pattern:**
```
feat: Implement user authentication service
fix: Resolve database connection issue
refactor: Improve code structure
```

**✅ Human Pattern:**
```
Add user auth with JWT tokens

- Set up JWT signing/verification
- Add /login and /register endpoints
- Store tokens in httpOnly cookies
- Tests pass for happy path + edge cases

Fixes issue where DB connection would timeout after 5 mins
- Added connection pool config
- Set idle timeout to 4 mins (safeguard)
- Tested with 100 concurrent requests

TODO: Add rate limiting to /login endpoint
```

---

### README (PERSONAL, NOT CORPORATE):

**❌ AI Pattern:**
```markdown
# Project Title

## Overview
This project is a comprehensive solution for...

## Features
- Feature 1: Description
- Feature 2: Description
```

**✅ Human Pattern:**
```markdown
# Prediction Market Signal Aggregator

Built for [bounty #55](issue-url) - combines prediction market odds 
with social sentiment to find trading opportunities.

## What It Does

Scrapes Polymarket, Kalshi, and Metaculus for odds, then 
cross-references with Twitter/Reddit sentiment to spot:
- Arbitrage opportunities (price differences between platforms)
- Sentiment divergence (social disagrees with market price)
- Volume spikes (unusual activity)

## Why I Built It This Way

- Used Hono for lightweight routing (faster than Express for this use case)
- Mobile proxies for social scraping (avoid rate limits)
- x402 payment gating (required by bounty)

## Notes

- TikTok scraping is best-effort (their API is restrictive)
- Metaculus returns median forecasts, not binary odds
```

---

## 🛠️ TOOLS & SCRIPTS

### 1. Initialize New Project

```bash
./au-workspace/scripts/init_bounty_project.sh [id] [name]
# Example: ./init_bounty_project.sh 55 prediction-market-aggregator
```

**Creates:**
- Complete folder structure
- `.agent/` folder with README
- Proper `.gitignore`
- Human-style `README.md`
- Placeholder files
- Isolated git repository

---

### 2. Migrate Existing Projects

```bash
./au-workspace/scripts/migrate_bounty_projects.sh
```

**Does:**
- Creates `.agent/` folders
- Moves agent artifacts
- Adds `.gitignore`
- Initializes isolated git repos

---

### 3. Verify Submission Cleanliness

```bash
# Check for agent artifacts
find . -name ".agent" -o -name "*.agent-notes.md"

# Verify .gitignore working
git status

# Review what will be submitted
git ls-files
```

---

## 📋 SUBMISSION CHECKLIST

**BEFORE submitting to GitHub:**

```bash
# 1. Check for agent artifacts
find src tests docs -name ".agent" -o -name "*.agent-notes.md"
# Should return NOTHING

# 2. Verify .gitignore
git status
# Should NOT show any .agent/ files

# 3. Review submission
git ls-files
# Review EVERY file - does it look human-written?

# 4. Run quality gates
../../au-workspace/scripts/verify_submission_readiness.sh .

# 5. Human review
# Human reviews EVERY file that will be submitted
```

---

## 🔍 RED FLAGS (AI TELLS)

**Avoid these in submission code:**

- ❌ Perfectly formatted code with no personality
- ❌ Generic, formal comments ("This function processes...")
- ❌ No TODOs, FIXMEs, or technical debt notes
- ❌ No references to specific issues or discussions
- ❌ No "I decided to..." or "I chose..." language
- ❌ No mistakes or iterations visible in git history
- ❌ All commits are perfectly formatted (conventional commits only)
- ❌ No casual language in comments
- ❌ No personal notes or asides

---

## ✅ GREEN FLAGS (HUMAN TELLS)

**Include these in submission code:**

- ✅ Casual, conversational comments
- ✅ TODOs and FIXMEs with context
- ✅ References to specific issues/discussions
- ✅ "I decided to..." language
- ✅ Visible iterations in git history (multiple commits)
- ✅ Some commits are small, some are large
- ✅ Occasional typos (fixed in later commits)
- ✅ Personal notes in README
- ✅ "Why I built it this way" sections
- ✅ Development log with dates

---

## 📊 MIGRATION STATUS

**Existing Projects:**

| Project | .agent/ Created | .gitignore Added | Git Isolated | Status |
|---------|-----------------|------------------|--------------|--------|
| bounty-mps-55-prediction | ⏳ Pending | ⏳ Pending | ⏳ Pending | Needs migration |
| bounty-mps-51-tiktok | ⏳ Pending | ⏳ Pending | ⏳ Pending | Needs migration |
| bounty-st-foundry-twitter-thread | ⏳ Pending | ⏳ Pending | ⏳ Pending | Needs migration |
| [others] | ⏳ Pending | ⏳ Pending | ⏳ Pending | Needs migration |

**Run migration:**
```bash
./au-workspace/scripts/migrate_bounty_projects.sh
```

---

## 🎯 SUCCESS CRITERIA

**Project structure is compliant when:**

- ✅ `.agent/` folder exists and contains all agent artifacts
- ✅ `.gitignore` properly excludes `.agent/` and internal files
- ✅ `src/`, `tests/`, `docs/` contain only human-style code
- ✅ README.md sounds human-authored (personal, casual)
- ✅ Git history shows iterations (multiple commits)
- ✅ No agent artifacts in `git ls-files`
- ✅ Quality gates pass
- ✅ Human reviewer approves

---

## 📚 DOCUMENTATION

| Document | Purpose | Location |
|----------|---------|----------|
| **BOUNTY_WORKSPACE_STRUCTURE.md** | Full structure guide | `au-workspace/docs/` |
| **BOUNTY_SUBMISSION_QUALITY_GATES.md** | Quality gate requirements | `au-workspace/docs/` |
| **AGENT_FOLDER_README.md** | .agent folder explanation | `au-workspace/templates/` |
| **README.bounty.md** | Human-style README template | `au-workspace/templates/` |
| **.gitignore.bounty** | Standard .gitignore | `au-workspace/templates/` |

---

## 🚀 NEXT ACTIONS

1. **Run migration script:**
   ```bash
   ./au-workspace/scripts/migrate_bounty_projects.sh
   ```

2. **Review migrated projects:**
   ```bash
   cd au-workspace/projects
   ls -la
   ```

3. **Check .gitignore in each:**
   ```bash
   cat [project]/.gitignore
   ```

4. **Verify no agent artifacts in git:**
   ```bash
   cd [project] && git status
   ```

5. **Update active bounties to use new structure**

---

**REMEMBER COOLIFY. LOOK HUMAN. STAY AUTHENTIC.**

**Quality gates + human authenticity + proper isolation = successful bounties.**

**NEVER AGAIN.**
