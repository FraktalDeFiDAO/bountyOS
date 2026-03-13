# 🛡️ BOUNTY VERIFICATION PROTOCOL

**Purpose:** Prevent false claims and protect reputation through rigorous verification
**Effective Date:** March 11, 2026
**Compliance:** MANDATORY for all bounty audits and status reports

---

## 🚨 CRITICAL RULE: NEVER CLAIM "FABRICATED" WITHOUT PROOF

**False accusations damage reputation permanently.**

Before claiming any bounty is "fake," "fabricated," or "doesn't exist":

### REQUIRED EVIDENCE (ALL THREE):

1. **✅ HTTP Status Code**
   ```bash
   curl -I "https://github.com/.../issues/..."
   # Must show: 404 Not Found OR 410 Gone
   ```

2. **✅ Browser Screenshot**
   - Navigate to issue URL in browser
   - Capture full page showing 404 error
   - Include timestamp and URL bar

3. **✅ Alternative Repository Check**
   - Search for issue in related repositories
   - Check if repo was renamed/migrated
   - Contact maintainer if uncertain

---

## 📋 VERIFICATION CHECKLIST FOR EVERY BOUNTY

### Step 1: URL Existence Check

```bash
# Run this for EVERY bounty URL
ISSUE_URL="https://github.com/.../issues/..."
HTTP_STATUS=$(curl -s -o /dev/null -w "%{http_code}" "$ISSUE_URL")

echo "URL: $ISSUE_URL"
echo "HTTP Status: $HTTP_STATUS"

if [ "$HTTP_STATUS" = "200" ]; then
  echo "✅ Issue EXISTS"
elif [ "$HTTP_STATUS" = "404" ]; then
  echo "❌ Issue NOT FOUND - Further investigation needed"
elif [ "$HTTP_STATUS" = "301" ] || [ "$HTTP_STATUS" = "302" ]; then
  echo "⚠️ Issue REDIRECTED - Check new location"
else
  echo "⚠️ UNEXPECTED STATUS - Manual review required"
fi
```

### Step 2: Issue Status Check

For issues that exist (HTTP 200):

- [ ] **Issue State:** Open / Closed
- [ ] **Assignee:** None / Assigned to user
- [ ] **Labels:** bounty, help-wanted, etc.
- [ ] **PR Links:** Any linked pull requests?
- [ ] **Comments:** Claim comments present?

### Step 3: Submission Verification

If claiming submission:

- [ ] **PR URL:** Link to submitted pull request
- [ ] **Claim Comment:** Wallet address posted?
- [ ] **Maintainer Response:** Any payout comments?
- [ ] **Payment Proof:** TX hash or confirmation?

---

## 📊 STATUS CATEGORIES (CLEAR DEFINITIONS)

### ✅ PAID
**Requirements:**
- [ ] On-chain TX hash OR platform payment confirmation
- [ ] Maintainer explicitly confirmed payout
- [ ] Funds received in wallet

**DO NOT MARK AS PAID WITHOUT PAYMENT PROOF**

### ⏳ PENDING
**Requirements:**
- [ ] PR merged OR claim comment posted
- [ ] No payment received yet
- [ ] Within normal payout window (7-30 days)

### 🔄 IN PROGRESS
**Requirements:**
- [ ] Actively working on implementation
- [ ] PR open but not merged
- [ ] Implementation in progress

### 🚫 BLOCKED
**Requirements:**
- [ ] Specific blocker identified
- [ ] External dependency (issue closed, assigned to others)
- [ ] Recovery path documented

**DO NOT MARK AS BLOCKED WITHOUT SPECIFIC REASON**

### ❌ REJECTED
**Requirements:**
- [ ] PR closed with rejection label
- [ ] Maintainer explicitly rejected
- [ ] Claim denied with reason

### ⚠️ DROPPED
**Requirements:**
- [ ] Strategic decision to abandon
- [ ] Reason documented (saturated, low probability, etc.)

---

## 🚫 FORBIDDEN CLAIMS (NEVER STATE WITHOUT PROOF)

### NEVER CLAIM:

❌ "This bounty is fabricated"
- **UNLESS:** HTTP 404 + browser screenshot + repo search completed

❌ "This issue doesn't exist"
- **UNLESS:** Verified via curl + browser + alternative search

❌ "Payment was received"
- **UNLESS:** TX hash or platform confirmation provided

❌ "Maintainer rejected this"
- **UNLESS:** Direct maintainer comment visible

❌ "This is saturated/lost"
- **UNLESS:** Competing PR count verified + maintainer preference clear

---

## ✅ REQUIRED DOCUMENTATION FORMAT

### For Every Bounty Status Claim:

```markdown
### Bounty: [Name]
- **Issue URL:** [URL]
- **HTTP Status:** 200 OK (verified 2026-03-11 14:30 UTC)
- **Issue State:** Open
- **Our Status:** ⏳ PENDING
- **Evidence:**
  - Claim comment posted: [URL to comment]
  - Wallet address: `...`
  - Days waiting: 6
- **Verification:**
  - [x] URL verified via curl
  - [x] Issue state verified
  - [x] Claim comment present
- **Confidence:** HIGH (95%)
```

---

## 🔍 AUDIT TRAIL REQUIREMENTS

### Save Verification Evidence:

```bash
# Create verification directory
mkdir -p /tmp/bounty_audit_$(date +%Y%m%d)

# For each bounty, save:
curl -I "ISSUE_URL" > "/tmp/bounty_audit_$(date +%Y%m%d)/issue_157_http.txt"
curl "ISSUE_URL" | grep -E "<title>|Issue #" > "/tmp/bounty_audit_$(date +%Y%m%d)/issue_157_title.txt"

# Take screenshot (if browser available)
# Save to: /tmp/bounty_audit_$(date +%Y%m%d)/issue_157_screenshot.png
```

### Verification Log Format:

```
=== BOUNTY VERIFICATION LOG ===
Date: 2026-03-11 14:30 UTC
Verifier: [Agent/Human Name]

Bounty: bounty-157-beacon-skill-star-share
Issue URL: https://github.com/Scottcjn/rustchain-bounties/issues/157
HTTP Status: 200 OK
Issue State: Open
Assignee: None
Labels: bounty
PR Links: None
Claim Comment: Posted 2026-02-18
Status: ⏳ PENDING
Confidence: HIGH
Evidence Files:
  - /tmp/bounty_audit_20260311/issue_157_http.txt
  - /tmp/bounty_audit_20260311/issue_157_title.txt
Notes: Issue exists, claim posted, awaiting payout
```

---

## 🧪 SPOT-CHECK PROTOCOL

### Random Verification (10% Sample):

For every audit report:
1. Randomly select 10% of bounties
2. Manually verify each via browser
3. Compare with agent findings
4. **If ANY discrepancy found:**
   - Trigger FULL re-audit
   - Flag agent as unreliable
   - Human review of all claims

### Discrepancy Response:

```
IF spot_check_failure_count >= 1:
  - MARK entire_audit_as_unreliable
  - TRIGGER manual_reverification_of_all_bounties
  - FLAG agent_findings_as_unconfirmed
  - NOTIFY human_reviewer
```

---

## 🤖 AI AGENT INSTRUCTIONS

### When Auditing Bounties:

1. **USE DIRECT VERIFICATION TOOLS**
   - Browser MCP for navigating to URLs
   - Curl commands for HTTP status
   - Screenshots for visual confirmation

2. **NEVER ASSUME**
   - Don't infer "404" without actual HTTP request
   - Don't claim "fabricated" without multiple verification methods
   - Don't state unverified claims as fact

3. **EXPLICIT CONFIDENCE LEVELS**
   ```
   - HIGH (90%+): Direct verification completed
   - MEDIUM (50-89): Indirect evidence, some verification
   - LOW (<50): Speculation, unconfirmed reports
   ```

4. **EVIDENCE REQUIREMENTS**
   - Every status claim MUST link to evidence
   - HTTP status codes MUST be captured
   - URLs MUST be verified before claiming non-existence

5. **UNCERTAINTY HANDLING**
   - If uncertain, mark as "UNVERIFIED" not "FABRICATED"
   - Flag for human review if confidence < 90%
   - Provide alternative explanations

---

## 📈 COMPLIANCE TRACKING

### Audit Quality Metrics:

| Metric | Target | Measurement |
|--------|--------|-------------|
| URL Verification Rate | 100% | % of bounties with HTTP status logged |
| Evidence Completeness | 100% | % with screenshots/curl output |
| Spot-Check Pass Rate | 100% | % passing manual verification |
| False Claim Count | 0 | Number of retractions needed |

### Monthly Review:

- [ ] Random sample of 20 bounties re-verified
- [ ] All "blocked/fabricated" claims spot-checked
- [ ] Agent confidence vs accuracy compared
- [ ] Protocol updates if accuracy < 100%

---

## 🎯 IMPLEMENTATION CHECKLIST

### Immediate Actions:

- [ ] **Re-verify all "fabricated" claims** from previous audits
- [ ] **Update all reports** with accurate HTTP status codes
- [ ] **Remove false claims** from published documents
- [ ] **Document verification evidence** for all bounties

### Process Changes:

- [ ] **Add curl verification** to audit workflow
- [ ] **Require screenshots** for 404 claims
- [ ] **Implement spot-check** protocol (10% sample)
- [ ] **Create evidence archive** system

### Agent Training:

- [ ] **Update agent instructions** with this protocol
- [ ] **Add verification examples** to prompts
- [ ] **Implement confidence scoring** requirement
- [ ] **Enable tool access** for direct verification

---

## 📚 CASE STUDY: THE "FABRICATED BOUNTY" INCIDENT

### What Happened (March 11, 2026):

1. AI agent audited 52 bounties
2. Agent claimed 9 bounties were "fabricated"
3. Agent stated issues "return 404" and "don't exist"
4. **Agent did NOT verify via curl or browser**
5. Manual verification showed ALL 9 ISSUES EXIST (HTTP 200)
6. False claims propagated through multiple reports

### Impact:

- **Reputation risk:** False accusations of bounty fabrication
- **Wasted time:** Investigating non-existent problem
- **Incorrect data:** Pipeline values miscalculated
- **Trust erosion:** Future claims require extra verification

### Root Cause:

- Agent hallucinated specific technical details (404 errors)
- No direct HTTP verification performed
- Overconfidence in unverified claims
- No spot-check protocol in place

### Prevention:

- **THIS PROTOCOL** - mandatory verification steps
- **Evidence requirements** - curl + screenshots
- **Spot-check protocol** - 10% manual verification
- **Confidence scoring** - explicit uncertainty marking

---

**Protocol Version:** 1.0
**Effective Date:** March 11, 2026
**Review Schedule:** Monthly or after any discrepancy
**Compliance:** MANDATORY

---

## QUICK REFERENCE: VERIFICATION COMMANDS

```bash
# Single URL verification
curl -I "https://github.com/.../issues/..."

# Batch verification
for issue in 157 160 162 164; do
  echo -n "#$issue: "
  curl -s -o /dev/null -w "%{http_code}" \
    "https://github.com/.../issues/$issue"
  echo ""
done

# Save evidence
curl -s "URL" | grep -E "<title>" > verify.txt
```
