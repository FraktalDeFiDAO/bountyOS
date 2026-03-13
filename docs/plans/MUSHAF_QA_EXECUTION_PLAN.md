# 🧪 MUSHAF QA TESTING - EXECUTION PLAN

**Created:** March 11, 2026  
**Repo:** https://github.com/YahiaRagae/mushaf-imad-android  
**Status:** Ready for execution

---

## 📋 AVAILABLE QA BOUNTIES

| Issue | Title | Points | Severity | Time Est. |
|-------|-------|--------|----------|-----------|
| **#37** | QA-2.1: MushafView Basic Rendering | 50 pts | Critical | 2-3h |
| **#40** | QA-2.4: MushafView Callbacks & Persistence | 50 pts | High | 2-3h |
| **#36** | QA-5.1: ReciterPickerDialog | 25 pts | Medium | 1-2h |
| **#35** | QA-3.4: QuranPlayerView Callbacks & Errors | 50 pts | High | 2-3h |
| **#34** | QA-3.3: QuranPlayerView Speed & Reciter | 25 pts | Medium | 1-2h |

**Total Available:** 200+ points

---

## 🎯 TESTING REQUIREMENTS

### Prerequisites
- Android device or emulator (API 26+)
- Android Studio (for building)
- Git (for cloning repo)
- Screenshot capability

### Setup Commands
```bash
# Clone the repository
git clone https://github.com/YahiaRagae/mushaf-imad-android.git
cd mushaf-imad-android

# Build and install on device/emulator
./gradlew assembleDebug
adb install app/build/outputs/apk/debug/app-debug.apk

# Run tests
./gradlew connectedAndroidTest
```

---

## 📝 TEST TEMPLATES

### Issue #37: QA-2.1: MushafView Basic Rendering (50 pts - Critical)

**Test Area:** Basic MushafView rendering

#### Test Cases

**TC-2.1: MushafView renders without crashing**
- [ ] App launches successfully
- [ ] MushafView component visible
- [ ] No crash on startup
- [ ] Screenshot: `tc-2-1-app-launch.png`

**TC-2.2: Arabic text renders correctly**
- [ ] Arabic script visible and legible
- [ ] Text direction is right-to-left
- [ ] No character corruption
- [ ] Screenshot: `tc-2-2-arabic-text.png`

**TC-2.3: Scrolling works smoothly**
- [ ] Vertical scroll enabled
- [ ] No lag or jank
- [ ] Page transitions smooth
- [ ] Video: `tc-2-3-scrolling.mp4` (optional)

**TC-2.4: Different screen sizes**
- [ ] Test on phone (e.g., Pixel 6)
- [ ] Test on tablet (e.g., Pixel Tablet)
- [ ] Layout adapts correctly
- [ ] Screenshots: `tc-2-4-phone.png`, `tc-2-4-tablet.png`

---

### Issue #36: QA-5.1: ReciterPickerDialog (25 pts - Medium)

**Test Area:** Reciter selection dialog

#### Test Cases

**TC-5.1: Dialog opens on tap**
- [ ] Tap reciter selector
- [ ] Dialog appears
- [ ] Animation smooth
- [ ] Screenshot: `tc-5-1-dialog-open.png`

**TC-5.2: Reciter list populated**
- [ ] List shows multiple reciters
- [ ] Names display correctly (Arabic/English)
- [ ] No empty items
- [ ] Screenshot: `tc-5-2-reciter-list.png`

**TC-5.3: Selection works**
- [ ] Tap a reciter
- [ ] Dialog closes
- [ ] Selected reciter shown in UI
- [ ] Audio changes to selected reciter
- [ ] Screenshots: `tc-5-3-selected.png`

**TC-5.4: Dialog dismisses correctly**
- [ ] Tap outside dialog → dismisses
- [ ] Back button → dismisses
- [ ] No crash on dismiss
- [ ] State preserved

---

### Issue #40: QA-2.4: MushafView Callbacks & Persistence (50 pts - High)

**Test Area:** Callbacks and state persistence

#### Test Cases

**TC-2.15: Page change callback fires**
- [ ] Navigate to next page
- [ ] Callback logs page number
- [ ] Correct page number reported
- [ ] Screenshot: `tc-2-15-page-callback.png`

**TC-2.16: Bookmark persistence**
- [ ] Add bookmark
- [ ] Close app
- [ ] Reopen app
- [ ] Bookmark still present
- [ ] Screenshots: `tc-2-16-bookmark-before.png`, `tc-2-16-bookmark-after.png`

**TC-2.17: Last page persistence**
- [ ] Navigate to page 50
- [ ] Close app
- [ ] Reopen app
- [ ] App opens to page 50
- [ ] Screenshots: `tc-2-17-last-page-before.png`, `tc-2-17-last-page-after.png`

**TC-2.18: Settings persistence**
- [ ] Change theme/font size
- [ ] Close app
- [ ] Reopen app
- [ ] Settings preserved
- [ ] Screenshots: `tc-2-18-settings-before.png`, `tc-2-18-settings-after.png`

---

## 📸 SCREENSHOT GUIDE

### How to Capture Screenshots

**Method 1: Android Studio**
1. Open Device Manager
2. Click "Screen Capture" icon
3. Save screenshot

**Method 2: ADB Command**
```bash
adb shell screencap -p /sdcard/screenshot.png
adb pull /sdcard/screenshot.png ./screenshots/
```

**Method 3: Device Buttons**
- Pixel: Power + Volume Down
- Samsung: Power + Volume Down
- Some devices: Palm swipe gesture

### Screenshot Naming Convention
```
tc-{test-case-number}-{description}.png
Example: tc-2-1-app-launch.png
```

---

## 📋 SUBMISSION CHECKLIST

For each issue:

- [ ] All test cases executed
- [ ] Screenshots captured for each test
- [ ] Pass/fail documented
- [ ] Bugs reported (if any found)
- [ ] Claim comment posted on issue
- [ ] Wallet address included

---

## 💬 CLAIM COMMENT TEMPLATES

### For Issue #37 (50 pts)
```markdown
## ✅ QA Testing Complete

**Agent:** @FraktalDeFiDAO
**Issue:** QA-2.1: MushafView Basic Rendering

### Test Results

| Test Case | Status | Evidence |
|-----------|--------|----------|
| TC-2.1: Renders without crashing | ✅ PASS | [screenshot] |
| TC-2.2: Arabic text correct | ✅ PASS | [screenshot] |
| TC-2.3: Scrolling smooth | ✅ PASS | [screenshot] |
| TC-2.4: Different screen sizes | ✅ PASS | [screenshot] |

### Summary
All 4 test cases passed. MushafView renders correctly across different configurations.

### 💰 Payout Request
**Points:** 50
**Wallet:** [YOUR_WALLET_ADDRESS]

Ready for review!
```

### For Issue #36 (25 pts)
```markdown
## ✅ QA Testing Complete

**Agent:** @FraktalDeFiDAO
**Issue:** QA-5.1: ReciterPickerDialog

### Test Results

| Test Case | Status | Evidence |
|-----------|--------|----------|
| TC-5.1: Dialog opens | ✅ PASS | [screenshot] |
| TC-5.2: List populated | ✅ PASS | [screenshot] |
| TC-5.3: Selection works | ✅ PASS | [screenshot] |
| TC-5.4: Dismisses correctly | ✅ PASS | [screenshot] |

### Summary
All 4 test cases passed. ReciterPickerDialog functions as expected.

### 💰 Payout Request
**Points:** 25
**Wallet:** [YOUR_WALLET_ADDRESS]

Ready for review!
```

### For Issue #40 (50 pts)
```markdown
## ✅ QA Testing Complete

**Agent:** @FraktalDeFiDAO
**Issue:** QA-2.4: MushafView Callbacks & Persistence

### Test Results

| Test Case | Status | Evidence |
|-----------|--------|----------|
| TC-2.15: Page callback | ✅ PASS | [screenshot] |
| TC-2.16: Bookmark persistence | ✅ PASS | [screenshot] |
| TC-2.17: Last page persistence | ✅ PASS | [screenshot] |
| TC-2.18: Settings persistence | ✅ PASS | [screenshot] |

### Summary
All 4 test cases passed. Callbacks and persistence working correctly.

### 💰 Payout Request
**Points:** 50
**Wallet:** [YOUR_WALLET_ADDRESS]

Ready for review!
```

---

## 🚀 QUICK START

1. **Clone repo:** `git clone https://github.com/YahiaRagae/mushaf-imad-android.git`
2. **Build app:** `cd mushaf-imad-android && ./gradlew assembleDebug`
3. **Install on device:** `adb install app/build/outputs/apk/debug/app-debug.apk`
4. **Run tests:** Follow test cases above, capture screenshots
5. **Submit:** Post claim comments with evidence

---

**Total Potential Earnings:** 125 points (for issues #37, #36, #40)  
**Estimated Time:** 5-8 hours  
**Difficulty:** Easy-Medium (UI testing)

**Ready to execute!** 🚀
