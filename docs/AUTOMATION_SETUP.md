# Bounty Automation System Setup

**Version:** 1.0  
**Date:** March 13, 2026

---

## 🚀 Quick Start

### 1. Install Dependencies

```bash
cd /home/administrator/projects/bountyOS

# Create virtual environment (if not exists)
python3 -m venv .venv

# Activate virtual environment
source .venv/bin/activate

# Install required packages
pip install requests python-dotenv
```

### 2. Configure Environment

Create `.env` file in bountyOS root:

```bash
# GitHub Configuration
GITHUB_TOKEN=your_github_personal_access_token
GITHUB_USERNAME=your_github_username

# Notification Configuration (optional)
DISCORD_WEBHOOK_URL=https://discord.com/api/webhooks/...
TELEGRAM_BOT_TOKEN=your_telegram_bot_token
TELEGRAM_CHAT_ID=your_telegram_chat_id

# Automation Settings
AUTOMATION_ENABLED=true
SCAN_INTERVAL_HOURS=6
FOLLOW_UP_DAYS=7
```

### 3. Test the Script

```bash
# Activate venv
source .venv/bin/activate

# Test PR tracking
python3 scripts/bounty_automation.py --track-prs --report

# Test bounty scanning
python3 scripts/bounty_automation.py --scan-bounties --report

# Run all automation
python3 scripts/bounty_automation.py --all --report
```

---

## ⏰ Cron Job Setup

### Option 1: Crontab (Recommended)

```bash
# Edit crontab
crontab -e

# Add these lines:

# PR tracking - every 6 hours
0 */6 * * * cd /home/administrator/projects/bountyOS && /usr/bin/python3 scripts/bounty_automation.py --track-prs >> logs/bounty_automation.log 2>&1

# Bounty scanning - every 6 hours
30 */6 * * * cd /home/administrator/projects/bountyOS && /usr/bin/python3 scripts/bounty_automation.py --scan-bounties >> logs/bounty_automation.log 2>&1

# Full report - daily at 8 AM
0 8 * * * cd /home/administrator/projects/bountyOS && /usr/bin/python3 scripts/bounty_automation.py --all --report >> logs/bounty_automation.log 2>&1
```

### Option 2: Systemd Timer

Create `/etc/systemd/system/bounty-automation.service`:

```ini
[Unit]
Description=Bounty Automation System
After=network.target

[Service]
Type=oneshot
User=administrator
WorkingDirectory=/home/administrator/projects/bountyOS
Environment="PATH=/home/administrator/projects/bountyOS/.venv/bin"
ExecStart=/home/administrator/projects/bountyOS/.venv/bin/python3 /home/administrator/projects/bountyOS/scripts/bounty_automation.py --all
StandardOutput=append:/home/administrator/projects/bountyOS/logs/bounty_automation.log
StandardError=append:/home/administrator/projects/bountyOS/logs/bounty_automation.log
```

Create `/etc/systemd/system/bounty-automation.timer`:

```ini
[Unit]
Description=Run Bounty Automation every 6 hours
Requires=bounty-automation.service

[Timer]
OnBootSec=5min
OnUnitActiveSec=6h
Unit=bounty-automation.service

[Install]
WantedBy=timers.target
```

Enable the timer:

```bash
sudo systemctl daemon-reload
sudo systemctl enable bounty-automation.timer
sudo systemctl start bounty-automation.timer
sudo systemctl status bounty-automation.timer
```

---

## 📊 Tracked PRs

The system automatically tracks these PRs:

| PR | Bounty | Value | Platform | Status |
|----|--------|-------|----------|--------|
| projectdiscovery/tlsx#956 | TLSX #819 | $1,200 | GitHub | Awaiting Review |
| bolivian-peru/marketplace-service-template#190 | MPS #51 | $75 SX | Proxies.sx | In Review |
| bolivian-peru/marketplace-service-template#189 | MPS #55 | $100 SX | Proxies.sx | In Review |
| bolivian-peru/marketplace-service-template#209 | MPS #70 | $100 SX | Proxies.sx | Submitted |
| openclaw-labs/openclaw#83 | OpenClaw CI+Tests | $20 | GitHub | Submitted |

To add more PRs, edit `PR_TRACKING_LIST` in `scripts/bounty_automation.py`.

---

## 🎯 Scanned Platforms

The system scans these platforms for new bounties:

1. **GitHub Bounties** - Issues with "bounty" label
2. **Gitcoin** - Open bounties via API
3. **Superteam Earn** - (Requires web scraping - manual for now)

To add more platforms, edit `BOUNTY_PLATFORMS` in `scripts/bounty_automation.py`.

---

## 🔔 Notifications

### Discord Setup

1. Create Discord webhook:
   - Go to Discord channel settings
   - Integrations → Webhooks → New Webhook
   - Copy webhook URL

2. Add to `.env`:
   ```bash
   DISCORD_WEBHOOK_URL=https://discord.com/api/webhooks/YOUR_WEBHOOK
   ```

### Telegram Setup

1. Create Telegram bot:
   - Message @BotFather on Telegram
   - Send `/newbot`
   - Follow instructions to get token

2. Get chat ID:
   - Message your bot
   - Visit: `https://api.telegram.org/bot<YOUR_TOKEN>/getUpdates`
   - Find your chat ID in response

3. Add to `.env`:
   ```bash
   TELEGRAM_BOT_TOKEN=your_bot_token
   TELEGRAM_CHAT_ID=your_chat_id
   ```

---

## 📁 Files Created

```
bountyOS/
├── scripts/
│   └── bounty_automation.py          # Main automation script
├── logs/
│   └── bounty_automation.log         # Automation logs
├── .env                              # Environment configuration
├── pr_tracker_state.json             # PR tracking state
├── seen_bounties.json                # Seen bounties cache
└── bounty_automation_report.txt      # Latest report
```

---

## 📈 Reports

### Automatic Reports

Reports are generated automatically and saved to:
- `bounty_automation_report.txt` (latest)
- `logs/bounty_automation.log` (historical)

### Report Contents

- PR status for all tracked bounties
- Days since submission
- Follow-up alerts
- New bounties discovered
- Total value tracked

### Example Report

```
============================================================
BOUNTY AUTOMATION REPORT
Generated: 2026-03-13 14:30:00
============================================================

📊 PR STATUS
----------------------------------------
🟡 TLSX #819 - Deadlock Fix
   State: open | Merged: False
   Days Open: 1
   Value: $1200

🟡 MPS #51 - TikTok API
   State: open | Merged: False
   Days Open: 20
   Value: $75
   ⚠️ NEEDS FOLLOW-UP

🟡 MPS #55 - Prediction Market
   State: open | Merged: False
   Days Open: 20
   Value: $100
   ⚠️ NEEDS FOLLOW-UP

Total Value Tracked: $1395

🎯 NEW BOUNTIES DISCOVERED
----------------------------------------
• New GitHub Bounty #123
  Value: $500
  URL: https://github.com/.../issues/123

============================================================
```

---

## 🛠️ Maintenance

### Update PR List

Edit `PR_TRACKING_LIST` in `scripts/bounty_automation.py`:

```python
PR_TRACKING_LIST = [
    PRConfig(
        repo="owner/repo",
        pr_number=123,
        bounty_name="Bounty Name",
        value_usd=1000,
        platform="Platform",
        submitted_date="2026-03-13",
        wallet_address="YourWalletAddress"
    ),
    # Add more PRs...
]
```

### Clear Seen Bounties

```bash
# Remove seen bounties cache
rm seen_bounties.json

# System will rebuild on next scan
```

### Reset PR Tracker State

```bash
# Remove PR tracker state
rm pr_tracker_state.json

# System will rebuild on next run
```

---

## 🐛 Troubleshooting

### GitHub Rate Limits

If you hit rate limits:
1. Get a GitHub token: https://github.com/settings/tokens
2. Add to `.env`: `GITHUB_TOKEN=your_token`
3. Restart automation

### Notifications Not Sending

Check:
1. Webhook/token is correct in `.env`
2. Network connectivity
3. Check logs: `tail -f logs/bounty_automation.log`

### Script Not Running

Check:
1. Cron is running: `systemctl status cron`
2. Python path is correct
3. Virtual environment is activated
4. Check cron logs: `grep CRON /var/log/syslog`

---

## 📊 Monitoring Dashboard

View current status:

```bash
# Latest report
cat bounty_automation_report.txt

# Live logs
tail -f logs/bounty_automation.log

# PR tracker state
cat pr_tracker_state.json | python3 -m json.tool

# Seen bounties
cat seen_bounties.json | python3 -m json.tool
```

---

## 🎯 Next Steps

### Week 1: Monitor & Adjust
- [ ] Run automation for 7 days
- [ ] Review daily reports
- [ ] Adjust follow-up timing if needed
- [ ] Add any missing PRs

### Week 2: Expand
- [ ] Add Superteam web scraping
- [ ] Add more platforms to scan
- [ ] Integrate with existing beacon automation
- [ ] Set up dashboard UI

### Week 3: Optimize
- [ ] Analyze which platforms have best bounties
- [ ] Tune scanning frequency
- [ ] Add filtering for bounty quality
- [ ] Implement auto-scoring

---

**Created:** March 13, 2026  
**Maintained By:** bountyOS Automation Team  
**Next Review:** March 20, 2026
