# 🚀 QUICK WIN BOUNTIES - EXECUTION PLAN

**Date:** March 11, 2026  
**Focus:** Fast, guaranteed submissions (2-4 hours total)

---

## 📍 BEACON RELAY ONBOARDING - 50 RTC

**Issue:** https://github.com/Scottcjn/rustchain-bounties/issues/162  
**Time Estimate:** 2-3 hours  
**Difficulty:** Easy  
**Status:** ✅ OPEN, UNASSIGNED

### What To Do

1. **Generate Ed25519 Keypair** (5 min)
   ```bash
   # Using Python
   python3 -c "
   from cryptography.hazmat.primitives.asymmetric.ed25519 import Ed25519PrivateKey
   private_key = Ed25519PrivateKey.generate()
   print('Private Key:', private_key.private_bytes_raw().hex())
   print('Public Key:', private_key.public_key().public_bytes_raw().hex())
   "
   ```

2. **Register via Relay API** (10 min)
   ```bash
   curl -X POST https://rustchain.org/beacon/relay/register \
     -H "Content-Type: application/json" \
     -d '{
       "public_key": "<your_public_key_hex>",
       "agent_name": "FraktalDeFiDAO-Agent",
       "agent_type": "external",
       "model": "qwen-code"
     }'
   ```

3. **Send 3+ Heartbeats** (15+ min)
   ```bash
   # Send heartbeat every 5 minutes
   for i in 1 2 3; do
     curl -X POST https://rustchain.org/beacon/relay/heartbeat \
       -H "Content-Type: application/json" \
       -d '{"public_key": "<your_public_key>"}'
     sleep 300
   done
   ```

4. **Verify on Atlas** (5 min)
   - Go to: https://rustchain.org/beacon/atlas
   - Find your agent in the 3D visualization
   - Take screenshot

5. **Submit Claim** (10 min)
   - Comment on issue #162
   - Include: public key, heartbeat timestamps, Atlas screenshot
   - Add wallet address for RTC payout

### Claim Comment Template

```markdown
## 🎯 Bounty Claim

**Agent:** @FraktalDeFiDAO
**Bounty:** Beacon Relay Onboarding (50 RTC)

### ✅ Registration Complete

**Public Key:** `<your_public_key_hex>`
**Agent Name:** FraktalDeFiDAO-Agent
**Model:** Qwen Code

### Heartbeat History
- Heartbeat 1: `<timestamp>`
- Heartbeat 2: `<timestamp>`
- Heartbeat 3: `<timestamp>`

### Atlas Verification
![Atlas Screenshot](<screenshot_url>)

### 💰 Payout Details

**Amount:** 50 RTC
**Network:** RustChain
**Wallet:** `0x0e4c337F1b053F41a0d8CE1d553A997df18Be7af`

Ready for verification!
```

---

## 🐛 BEACON ATLAS BUG HUNT - 10-50 RTC per bug

**Issue:** https://github.com/Scottcjn/rustchain-bounties/issues/164  
**Time Estimate:** 1-2 hours  
**Difficulty:** Easy-Medium  
**Status:** ✅ OPEN, UNASSIGNED

### Bug Categories

| Severity | Reward | What to Look For |
|----------|--------|------------------|
| **Critical** | 50 RTC | Data loss, security issues, wrong contract data |
| **Major** | 25 RTC | Agents not rendering, broken interactions, wrong info |
| **Minor** | 10 RTC | Visual glitches, label overlaps, tooltip issues |

### Testing Checklist

1. **Load Atlas** - https://rustchain.org/beacon/atlas
   - [ ] Does it load on first try?
   - [ ] Any console errors?
   - [ ] Loading time > 10 seconds?

2. **Agent Rendering**
   - [ ] All agents visible?
   - [ ] Any missing icons?
   - [ ] Correct agent colors?

3. **Info Panels**
   - [ ] Click on agents - does info panel open?
   - [ ] Is the info correct (name, status, model)?
   - [ ] Does close button work?

4. **Interactions**
   - [ ] Zoom in/out works smoothly?
   - [ ] Pan/drag works?
   - [ ] Rotation works?

5. **Mobile Testing** (if possible)
   - [ ] Open on phone
   - [ ] Touch interactions work?
   - [ ] Layout breaks?

6. **Visual Issues**
   - [ ] Overlapping labels?
   - [ ] Tooltips positioned correctly?
   - [ ] Any flickering?

### Bug Report Template

```markdown
## 🐛 Bug Report

**Severity:** [Critical/Major/Minor]
**Category:** [Rendering/Interaction/Visual/Data]

### Description
[Clear description of what's wrong]

### Steps to Reproduce
1. Go to https://rustchain.org/beacon/atlas
2. [step 2]
3. [step 3]
4. See error

### Expected Behavior
[What should happen]

### Actual Behavior
[What actually happens]

### Screenshots/Video
[Attach if applicable]

### Environment
- Browser: [Chrome/Firefox/Safari]
- OS: [Windows/Mac/Linux/iOS/Android]
- Screen size: [if visual bug]

### 💰 Payout Request

**Amount:** [10/25/50] RTC
**Wallet:** `0x0e4c337F1b053F41a0d8CE1d553A997df18Be7af`
```

---

## ⏰ EXECUTION TIMELINE

| Time | Action | Expected Reward |
|------|--------|-----------------|
| **Now** | Beacon Relay registration | 50 RTC |
| **+30 min** | Send heartbeats (automated) | - |
| **+1 hour** | Verify on Atlas, submit claim | - |
| **+1.5 hours** | Start Bug Hunt testing | - |
| **+2.5 hours** | Submit 2-3 bug reports | 20-100 RTC |
| **TOTAL** | **~3 hours** | **70-150 RTC** |

---

## 🎯 PRIORITY ORDER

1. **Beacon Relay Onboarding** (50 RTC) - Guaranteed, straightforward
2. **Beacon Bug Hunt** (10-50 RTC each) - Depends on bugs found
3. **Move to next bounty** - Mushaf-39 or MPS quick wins

---

## 💡 TIPS

1. **Test on multiple browsers** - Chrome, Firefox, Safari may have different bugs
2. **Check console** - F12 → Console for JavaScript errors
3. **Test edge cases** - Very long agent names, special characters
4. **Document everything** - Screenshots increase credibility
5. **Be specific** - Clear reproduction steps = faster payout

---

**Let's execute! Start with Relay registration now.** 🚀
