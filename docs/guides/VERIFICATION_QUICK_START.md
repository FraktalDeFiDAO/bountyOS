# 🛡️ BOUNTY VERIFICATION QUICK START

**Purpose:** Get started with bounty verification in under 2 minutes
**Effective:** March 11, 2026

---

## ⚡ ONE-LINER VERIFICATION

```bash
# Verify all bounties from a URL list
cd /home/administrator/projects/bountyOS
./scripts/verify_bounties.sh au-workspace/config/bounty_urls.txt
```

**Output:** Evidence saved to `./audit_trails/YYYY-MM-DD/`

---

## 📋 COMMON COMMANDS

### Bash Script (Simple)

```bash
# Default (uses bounty_urls.txt in current dir)
./scripts/verify_bounties.sh

# Custom input and output
./scripts/verify_bounties.sh /path/to/bounties.txt /path/to/output

# Example with timestamp
./scripts/verify_bounties.sh bounties.tsv ./audit_trails/$(date +%Y-%m-%d)
```

### Python Tool (Advanced with parallel processing)

```bash
# Basic usage
python3 scripts/verify_bounties.py bounty_urls.txt

# Custom output directory
python3 scripts/verify_bounties.py bounties.tsv ./audit_trails/2026-03-11

# Results saved to:
# - verification_report.md (human-readable)
# - verification_results.json (machine-readable)
# - http_responses/ (evidence files)
```

---

## 📁 INPUT FILE FORMATS

### TSV Format (Recommended)
```tsv
bounty-157-beacon	https://github.com/.../issues/157
bounty-160-beacon	https://github.com/.../issues/160
```

### Pipe Format
```
bounty-157-beacon|https://github.com/.../issues/157
bounty-160-beacon|https://github.com/.../issues/160
```

### Markdown Format
```markdown
[bounty-157-beacon](https://github.com/.../issues/157)
[bounty-160-beacon](https://github.com/.../issues/160)
```

---

## 🔍 HTTP STATUS INTERPRETATIONS

| Status | Meaning | Action |
|--------|---------|--------|
| **200 OK** | ✅ Issue exists | Continue with status check |
| **404 Not Found** | ❌ Issue doesn't exist | Manual verification required |
| **301/302** | ⚠️ Redirected | Check new location |
| **403 Forbidden** | ⚠️ Private repo | May need authentication |
| **000/Error** | ❌ Network error | Retry or flag for manual check |

---

## 📊 OUTPUT STRUCTURE

```
audit_trails/
└── 2026-03-11/
    ├── verification_log.txt       # Bash script log
    ├── verification_report.md     # Human-readable report
    ├── verification_results.json  # Machine-readable data
    ├── summary.json               # Quick stats
    ├── http_responses/            # HTTP header evidence
    │   ├── bounty-157_157_2026-03-11_14-30-00.txt
    │   └── ...
    ├── screenshots/               # Browser screenshots (if available)
    └── verification_logs/         # Additional logs
```

---

## 🧪 EXAMPLE WORKFLOW

### Step 1: Prepare URL List

```bash
# Create bounty_urls.txt
cat > bounty_urls.txt << 'EOF'
bounty-157-beacon	https://github.com/Scottcjn/rustchain-bounties/issues/157
bounty-160-beacon	https://github.com/Scottcjn/rustchain-bounties/issues/160
bounty-162-relay	https://github.com/Scottcjn/rustchain-bounties/issues/162
EOF
```

### Step 2: Run Verification

```bash
# Run verification
./scripts/verify_bounties.sh bounty_urls.txt

# Output:
# === BOUNTY VERIFICATION SCRIPT ===
# Timestamp: 2026-03-11_14-30-00
# Verifying #157 (bounty-157-beacon)... ✅ 200 OK
# Verifying #160 (bounty-160-beacon)... ✅ 200 OK
# Verifying #162 (bounty-162-relay)... ✅ 200 OK
# 
# === VERIFICATION SUMMARY ===
# Total Bounties: 3
# HTTP 200 (Exists): 3
# HTTP 404 (Not Found): 0
# Verification Rate: 100%
```

### Step 3: Review Evidence

```bash
# View evidence files
ls -la audit_trails/2026-03-11/http_responses/

# View specific evidence
cat audit_trails/2026-03-11/http_responses/bounty-157_157_*.txt

# View full report
cat audit_trails/2026-03-11/verification_report.md
```

---

## 🚨 TROUBLESHOOTING

### "Input file not found"
```bash
# Check file exists
ls -la bounty_urls.txt

# Use absolute path
./scripts/verify_bounties.sh /full/path/to/bounty_urls.txt
```

### "Permission denied"
```bash
# Make executable
chmod +x scripts/verify_bounties.sh
chmod +x scripts/verify_bounties.py
```

### "curl: command not found"
```bash
# Install curl
sudo apt-get install curl  # Debian/Ubuntu
sudo yum install curl      # RHEL/CentOS
```

### "ModuleNotFoundError: No module named 'requests'"
```bash
# Install Python requests
pip3 install requests
```

---

## 📝 EVIDENCE FILE NAMING

```
{bounty_name}_{issue_num}_{timestamp}.txt

Example: bounty-157-beacon_157_2026-03-11_14-30-00.txt
```

---

## 🔗 RELATED DOCUMENTS

- **Full Protocol:** `BOUNTY_VERIFICATION_PROTOCOL.md`
- **Audit Trail:** `audit_trails/YYYY-MM-DD/`
- **Agent Instructions:** `templates/audit_agent_instructions.md`
- **QA Checklist:** `templates/audit_qa_checklist.md`

---

## ⚡ QUICK REFERENCE CARD

```bash
# Verify bounties
./scripts/verify_bounties.sh [input] [output]

# Python version (parallel)
python3 scripts/verify_bounties.py [input] [output]

# View results
cat audit_trails/$(date +%Y-%m-%d)/verification_report.md

# View evidence
ls audit_trails/$(date +%Y-%m-%d)/http_responses/
```

---

**Version:** 1.0 | **Effective:** March 11, 2026 | **Compliance:** MANDATORY
