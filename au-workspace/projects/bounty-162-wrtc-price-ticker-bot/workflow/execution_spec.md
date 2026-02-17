# Execution Spec: bounty-162-wrtc

Issue: https://github.com/Scottcjn/Rustchain/issues/162
Title: [BOUNTY] wRTC Price Ticker Bot — Discord/Telegram Price Alerts (20 RTC)
Last updated: 2026-02-14

## 1) Goal

Ship a working bot command that returns live wRTC pricing and metrics, then add optional alerting features for full payout.

## 2) Scope and Payout Mapping

Base deliverable (12 RTC):
- Working `/price` or `!price` command
- Live price data from DexScreener/Jupiter
- Output includes:
  - price in SOL
  - price in USD
  - 24h change
  - liquidity

Extended deliverable (8 RTC):
- Auto-posting schedule (hourly)
- price movement alert (>10% in 1h)
- clean formatted output

## 3) Technical Architecture

- Runtime: Python 3.11+
- Bot client: choose one
  - Discord: `discord.py`
  - Telegram: `python-telegram-bot`
- Data providers:
  - DexScreener token endpoint (primary)
  - Jupiter price API (fallback)
- Internal modules:
  - `providers.py`: API fetch/normalize
  - `formatter.py`: text message rendering
  - `bot.py`: command handlers + scheduler
  - `config.py`: token/chat/channel/env config

## 4) Implementation Plan

### Phase A: Scaffold (30 min)

- Create project structure and virtualenv.
- Add `.env.example` with bot token + channel/chat id.
- Implement provider response schema contract.

### Phase B: Core command (60 min)

- Parse DexScreener response and map fields.
- Implement fallback to Jupiter when DexScreener fails.
- Wire `/price` command and render output.
- Add retries/timeouts and graceful error messages.

### Phase C: Extended features (60-90 min)

- Add hourly scheduler.
- Persist last observed price in memory/file.
- Implement >10%/1h alert logic.
- Add startup health log and command usage docs.

### Phase D: QA + evidence (30 min)

- Test with live API.
- Capture command output screenshot.
- Capture log snippet showing periodic posts/alerts.

## 5) Acceptance Criteria

- [ ] Bot starts without runtime errors.
- [ ] `/price` returns current metrics from live data.
- [ ] Failover path works when primary API fails.
- [ ] Auto-post and alert path demonstrated (for full payout).
- [ ] README includes setup + run instructions.

## 6) Risks and Mitigations

- API schema drift:
  - parse defensively; default missing fields.
- Bot token misconfiguration:
  - include validation at startup.
- Rate limits/network issues:
  - cache last successful response and show staleness note.

## 7) Submission Package

- Source code + README
- Screenshot of `/price` command output
- Optional screenshot/log of scheduled alerting
- Issue comment summarizing implementation and proof

## 8) Tracker Commands

```bash
cd /home/administrator/projects/bountyOS/au-workspace
./bin/au bounty move bounty-162-wrtc in_progress
./bin/au bounty progress bounty-162-wrtc 40
./bin/au bounty next bounty-162-wrtc "Wire /price command and DexScreener parser"
```
