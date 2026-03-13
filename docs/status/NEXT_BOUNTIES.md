# 📋 NEXT BOUNTIES TO ATTACK

**Generated:** March 11, 2026

---

## ⏳ CURRENTLY RUNNING

### Beacon Relay Onboarding (50 RTC)
- **Status:** Auto-submitter running in background
- **PID:** 164497
- **Progress:** Retrying registration (rate limited / API issues)
- **Monitor:** `bash beacon_monitor.sh`
- **ETA:** Unknown (depends on API)

---

## 🎯 HIGH PRIORITY (Open QA Bounties)

### Mushaf QA Bounties - Multiple Available

| Issue | Title | Points | Severity | Time Est. |
|-------|-------|--------|----------|-----------|
| **#37** | QA-2.1: MushafView Basic Rendering | 50 pts | Critical | 2-3h |
| **#40** | QA-2.4: MushafView Callbacks & Persistence | 50 pts | High | 2-3h |
| **#35** | QA-3.4: QuranPlayerView Callbacks & Errors | 50 pts | High | 2-3h |
| **#36** | QA-5.1: ReciterPickerDialog | 25 pts | Medium | 1-2h |
| **#34** | QA-3.3: QuranPlayerView Speed & Reciter | 25 pts | Medium | 1-2h |

**Total Available:** 200+ points

**Repo:** https://github.com/YahiaRagae/mushaf-imad-android  
**Label:** QA, testing, bounty

---

## 🐛 BUG HUNTS

### Beacon Atlas Bug Hunt (10-50 RTC per bug)

**Issue:** https://github.com/Scottcjn/rustchain-bounties/issues/164

**Potential Bugs Found:**
1. **Atlas stuck on "Loading..."** - Major (25 RTC)
   - 3D visualization never renders
   - Only shows "Initializing renderer..."
   
2. **No error handling** - Minor (10 RTC)
   - Silent failures, no user feedback

**Total Potential:** 35+ RTC

**Status:** Bug report ready at `BEACON_BUG_REPORT.md`

---

## 💡 RECOMMENDED ORDER

### While Beacon Auto-Submitter Runs:

1. **Mushaf QA-36** (25 pts, 1-2h) - Quickest win
   - ReciterPickerDialog testing
   - Straightforward UI testing
   
2. **Beacon Bug Report** (35 RTC, 30 min)
   - Submit the loading issue report
   - Already documented

3. **Mushaf QA-37** (50 pts, 2-3h) - Critical severity
   - Basic rendering tests
   - High impact

4. **Check Beacon Status**
   - If registered → send heartbeats
   - If still failing → move on

---

## 📊 POTENTIAL EARNINGS

| Bounty | Reward | Time |
|--------|--------|------|
| Beacon Relay (running) | 50 RTC | Auto |
| Beacon Bugs | 35 RTC | 30 min |
| Mushaf QA-36 | 25 pts | 1-2h |
| Mushaf QA-37 | 50 pts | 2-3h |
| **Total** | **160 pts + RTC** | **~6h** |

---

**Ready to execute!**
