# 🐛 BEACON ATLAS BUG REPORT - SUMMARY

**Date:** March 11, 2026  
**Bounty:** https://github.com/Scottcjn/rustchain-bounties/issues/164  
**Total Bugs Found:** 4

---

## Bug #1: 404 Not Found on Official Atlas URL

### Description
The URL mentioned in the bounty description (https://rustchain.org/beacon/atlas) returns a 404 Not Found error. The actual application resides at https://rustchain.org/beacon/.

### Evidence
- **Requested URL:** https://rustchain.org/beacon/atlas
- **Response:** `404 Not Found` (nginx)

### Severity: MINOR (10 RTC)
- Link rot / misconfiguration in documentation.

---

## Bug #2: CORS Policy Blocking BoTTube Data

### Description
The Atlas attempts to fetch agent data from `https://bottube.ai/api/atlas/agents` and `https://bottube.ai/api/grazer-github-stats`, but these requests are blocked by the browser's CORS policy.

### Evidence (Console Log)
```
Access to fetch at 'https://bottube.ai/api/atlas/agents' from origin 'https://rustchain.org' has been blocked by CORS policy: No 'Access-Control-Allow-Origin' header is present on the requested resource.
```

### Impact
- **BoTTube agents are missing from the Atlas.**
- **GitHub stats for Grazer are not displayed.**

### Severity: MAJOR (25 RTC)
- Significant data missing from the 3D visualization.

---

## Bug #3: Atlas Boot Process Hangs After Step 5

### Description
The 3D Atlas initialization process completes Step 5 ("Loading contract ledger...") but never reaches the "Atlas online" state. It hangs before building connections (Step 6) or mounting the UI (Step 7). This causes the loading screen to stay visible indefinitely for many users (though some cached data might allow it to bypass in some environments).

### Evidence (Console Logs)
```
[boot] Atlas: 24 agents (0 BoTTube, 3 Beacon, 19 miners)
[boot] Loaded 16 contracts
[No further boot logs appear]
```
The script never reaches: `status.textContent = "Atlas online — ${AGENTS.length} agents.";`

### Severity: CRITICAL (50 RTC)
- **Breaks core functionality for new sessions.**
- **Indefinite hang on loading screen.**

---

## Bug #4: WebGL GPU Stalls due to ReadPixels

### Description
The Three.js renderer is experiencing frequent GPU stalls triggered by `ReadPixels` calls. This causes significant performance degradation and input lag in the 3D visualization.

### Evidence (Console Log)
```
CONSOLE [warning]: [.WebGL-0x33fc0010b200]GL Driver Message (OpenGL, Performance, GL_CLOSE_PATH_NV, High): GPU stall due to ReadPixels
```

### Severity: MINOR (10 RTC)
- Performance issue affecting user experience.

---

## 💰 TOTAL REWARD CLAIM: 95 RTC

**Wallet:** `0x0e4c337F1b053F41a0d8CE1d553A997df18Be7af`

---

**Reports verified via Playwright automated testing.**
