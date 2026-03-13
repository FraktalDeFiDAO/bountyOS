#!/usr/bin/env python3
"""
Bounty Automation System
========================
Automated PR tracking and bounty discovery for bountyOS

Features:
- PR status monitoring (TLSX, MPS, OpenClaw, etc.)
- Auto follow-up after 7 days
- Bounty discovery scanning
- Discord/Telegram notifications

Usage:
    python3 bounty_automation.py --track-prs
    python3 bounty_automation.py --scan-bounties
    python3 bounty_automation.py --all
"""

import os
import sys
import json
import time
import requests
from datetime import datetime, timedelta
from typing import Dict, List, Optional
from dataclasses import dataclass, asdict
import logging

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    handlers=[
        logging.FileHandler('bounty_automation.log'),
        logging.StreamHandler()
    ]
)
logger = logging.getLogger('bounty_automation')


# ============================================================================
# CONFIGURATION
# ============================================================================

@dataclass
class PRConfig:
    """Configuration for PR tracking"""
    repo: str
    pr_number: int
    bounty_name: str
    value_usd: float
    platform: str
    submitted_date: str
    wallet_address: Optional[str] = None
    follow_up_days: int = 7


@dataclass
class BountyConfig:
    """Configuration for bounty scanning"""
    platform: str
    url: str
    api_endpoint: Optional[str]
    scan_interval_hours: int = 6


# PRs to Track
PR_TRACKING_LIST = [
    PRConfig(
        repo="projectdiscovery/tlsx",
        pr_number=956,
        bounty_name="TLSX #819 - Deadlock Fix",
        value_usd=1200,
        platform="GitHub",
        submitted_date="2026-03-12",
        follow_up_days=7
    ),
    PRConfig(
        repo="bolivian-peru/marketplace-service-template",
        pr_number=190,
        bounty_name="MPS #51 - TikTok API",
        value_usd=75,
        platform="Proxies.sx",
        submitted_date="2026-02-21",
        wallet_address="FH84Dg6gh7bWtyZ5a1SBNLp1JBesLoCKx9mekJpr7zHR"
    ),
    PRConfig(
        repo="bolivian-peru/marketplace-service-template",
        pr_number=189,
        bounty_name="MPS #55 - Prediction Market",
        value_usd=100,
        platform="Proxies.sx",
        submitted_date="2026-02-21",
        wallet_address="FH84Dg6gh7bWtyZ5a1SBNLp1JBesLoCKx9mekJpr7zHR"
    ),
    PRConfig(
        repo="bolivian-peru/marketplace-service-template",
        pr_number=209,
        bounty_name="MPS #70 - Trend Intelligence",
        value_usd=100,
        platform="Proxies.sx",
        submitted_date="2026-03-12",
        wallet_address="FH84Dg6gh7bWtyZ5a1SBNLp1JBesLoCKx9mekJpr7zHR"
    ),
    PRConfig(
        repo="openclaw-labs/openclaw",
        pr_number=83,
        bounty_name="OpenClaw CI + Tests",
        value_usd=20,
        platform="GitHub",
        submitted_date="2026-03-12"
    ),
]

# Bounty Platforms to Scan
BOUNTY_PLATFORMS = [
    BountyConfig(
        platform="GitHub Bounties",
        url="https://github.com/topics/bounty",
        api_endpoint="https://api.github.com/search/issues?q=label:bounty+state:open+created:>%s"
    ),
    BountyConfig(
        platform="Superteam Earn",
        url="https://superteam.fun/earn",
        api_endpoint=None  # Web scraping required
    ),
    BountyConfig(
        platform="Gitcoin",
        url="https://gitcoin.co/explorer",
        api_endpoint="https://api.gitcoin.co/api/v1/bounties/"
    ),
]

# Notification Settings
DISCORD_WEBHOOK = os.getenv('DISCORD_WEBHOOK_URL')
TELEGRAM_BOT_TOKEN = os.getenv('TELEGRAM_BOT_TOKEN')
TELEGRAM_CHAT_ID = os.getenv('TELEGRAM_CHAT_ID')


# ============================================================================
# GITHUB API INTEGRATION
# ============================================================================

class GitHubAPI:
    """GitHub API wrapper for PR tracking"""
    
    def __init__(self, token: Optional[str] = None):
        self.token = token or os.getenv('GITHUB_TOKEN')
        self.base_url = "https://api.github.com"
        self.session = requests.Session()
        
        if self.token:
            self.session.headers.update({
                'Authorization': f'token {self.token}',
                'Accept': 'application/vnd.github.v3+json'
            })
        else:
            logger.warning("No GitHub token provided. Rate limits will apply.")
    
    def get_pr_status(self, repo: str, pr_number: int) -> Optional[Dict]:
        """Get PR status and details"""
        try:
            url = f"{self.base_url}/repos/{repo}/pulls/{pr_number}"
            response = self.session.get(url, timeout=10)
            
            if response.status_code == 200:
                data = response.json()
                return {
                    'state': data.get('state'),
                    'merged': data.get('merged', False),
                    'merged_at': data.get('merged_at'),
                    'merged_by': data.get('merged_by', {}).get('login') if data.get('merged_by') else None,
                    'comments': data.get('comments'),
                    'review_comments': data.get('review_comments'),
                    'commits': data.get('commits'),
                    'additions': data.get('additions'),
                    'deletions': data.get('deletions'),
                    'updated_at': data.get('updated_at'),
                    'created_at': data.get('created_at'),
                    'user': data.get('user', {}).get('login'),
                    'title': data.get('title'),
                    'labels': [label['name'] for label in data.get('labels', [])]
                }
            else:
                logger.error(f"GitHub API error: {response.status_code} - {response.text}")
                return None
                
        except Exception as e:
            logger.error(f"Error fetching PR {repo}#{pr_number}: {e}")
            return None
    
    def get_pr_comments(self, repo: str, pr_number: int) -> List[Dict]:
        """Get PR comments"""
        try:
            url = f"{self.base_url}/repos/{repo}/issues/{pr_number}/comments"
            response = self.session.get(url, timeout=10)
            
            if response.status_code == 200:
                return response.json()
            return []
        except Exception as e:
            logger.error(f"Error fetching comments: {e}")
            return []
    
    def get_pr_reviews(self, repo: str, pr_number: int) -> List[Dict]:
        """Get PR reviews"""
        try:
            url = f"{self.base_url}/repos/{repo}/pulls/{pr_number}/reviews"
            response = self.session.get(url, timeout=10)
            
            if response.status_code == 200:
                return response.json()
            return []
        except Exception as e:
            logger.error(f"Error fetching reviews: {e}")
            return []
    
    def post_comment(self, repo: str, pr_number: int, comment: str) -> bool:
        """Post comment on PR"""
        try:
            url = f"{self.base_url}/repos/{repo}/issues/{pr_number}/comments"
            response = self.session.post(url, json={'body': comment}, timeout=10)
            
            if response.status_code == 201:
                logger.info(f"Comment posted on {repo}#{pr_number}")
                return True
            else:
                logger.error(f"Failed to post comment: {response.status_code}")
                return False
        except Exception as e:
            logger.error(f"Error posting comment: {e}")
            return False


# ============================================================================
# PR TRACKER
# ============================================================================

class PRTracker:
    """Track PR status and automate follow-ups"""
    
    def __init__(self, github: GitHubAPI):
        self.github = github
        self.state_file = 'pr_tracker_state.json'
        self.load_state()
    
    def load_state(self):
        """Load tracker state from file"""
        try:
            if os.path.exists(self.state_file):
                with open(self.state_file, 'r') as f:
                    self.state = json.load(f)
            else:
                self.state = {}
        except Exception as e:
            logger.error(f"Error loading state: {e}")
            self.state = {}
    
    def save_state(self):
        """Save tracker state to file"""
        try:
            with open(self.state_file, 'w') as f:
                json.dump(self.state, f, indent=2)
        except Exception as e:
            logger.error(f"Error saving state: {e}")
    
    def track_all_prs(self) -> Dict:
        """Track all configured PRs"""
        results = {}
        
        for pr_config in PR_TRACKING_LIST:
            logger.info(f"Tracking {pr_config.bounty_name} ({pr_config.repo}#{pr_config.pr_number})")
            
            pr_status = self.github.get_pr_status(
                pr_config.repo,
                pr_config.pr_number
            )
            
            if pr_status:
                # Check if follow-up needed
                submitted_date = datetime.strptime(pr_config.submitted_date, '%Y-%m-%d')
                days_since = (datetime.now() - submitted_date).days
                
                needs_followup = (
                    days_since >= pr_config.follow_up_days and
                    pr_status['state'] == 'open' and
                    not self._has_maintainer_response(pr_config.repo, pr_config.pr_number)
                )
                
                results[pr_config.bounty_name] = {
                    'status': pr_status,
                    'days_since_submission': days_since,
                    'needs_followup': needs_followup,
                    'value_usd': pr_config.value_usd,
                    'platform': pr_config.platform,
                    'wallet': pr_config.wallet_address
                }
                
                # Auto follow-up if needed
                if needs_followup:
                    self._auto_followup(pr_config)
        
        self.save_state()
        return results
    
    def _has_maintainer_response(self, repo: str, pr_number: int) -> bool:
        """Check if PR has maintainer response"""
        comments = self.github.get_pr_comments(repo, pr_number)
        reviews = self.github.get_pr_reviews(repo, pr_number)
        
        # Check for maintainer comments (repo owner or collaborator)
        for comment in comments:
            if comment.get('user', {}).get('login') != os.getenv('GITHUB_USERNAME'):
                return True
        
        # Check for approving reviews
        for review in reviews:
            if review.get('state') == 'APPROVED':
                return True
        
        return False
    
    def _auto_followup(self, pr_config: PRConfig):
        """Send automated follow-up comment"""
        followup_comment = f"""
👋 Friendly follow-up! This PR has been open for {pr_config.follow_up_days}+ days.

**Bounty Details:**
- **Name:** {pr_config.bounty_name}
- **Value:** ${pr_config.value_usd}
- **Wallet:** `{pr_config.wallet_address or 'N/A'}`

Is there anything else needed from my side to move this forward? Happy to make any necessary changes! 🚀
"""
        
        # Don't actually post - just log for now
        logger.info(f"Follow-up needed for {pr_config.bounty_name}")
        logger.info(f"Comment: {followup_comment.strip()}")
        
        # To actually post, uncomment:
        # self.github.post_comment(pr_config.repo, pr_config.pr_number, followup_comment.strip())
        
        # Send notification
        send_discord_notification(
            f"🔔 Follow-up Needed: {pr_config.bounty_name}",
            f"PR open for {pr_config.follow_up_days}+ days\nValue: ${pr_config.value_usd}"
        )


# ============================================================================
# BOUNTY SCANNER
# ============================================================================

class BountyScanner:
    """Scan platforms for new bounties"""
    
    def __init__(self, github: GitHubAPI):
        self.github = github
        self.seen_bounties_file = 'seen_bounties.json'
        self.load_seen_bounties()
    
    def load_seen_bounties(self):
        """Load previously seen bounties"""
        try:
            if os.path.exists(self.seen_bounties_file):
                with open(self.seen_bounties_file, 'r') as f:
                    self.seen_bounties = json.load(f)
            else:
                self.seen_bounties = []
        except Exception as e:
            logger.error(f"Error loading seen bounties: {e}")
            self.seen_bounties = []
    
    def save_seen_bounties(self):
        """Save seen bounties to file"""
        try:
            with open(self.seen_bounties_file, 'w') as f:
                json.dump(self.seen_bounties, f, indent=2)
        except Exception as e:
            logger.error(f"Error saving seen bounties: {e}")
    
    def scan_all_platforms(self) -> List[Dict]:
        """Scan all configured platforms"""
        new_bounties = []
        
        for platform in BOUNTY_PLATFORMS:
            logger.info(f"Scanning {platform.platform}")
            
            if platform.platform == "GitHub Bounties":
                bounties = self._scan_github_bounties()
            elif platform.platform == "Gitcoin":
                bounties = self._scan_gitcoin()
            else:
                logger.warning(f"Platform {platform.platform} requires web scraping - skipping")
                continue
            
            # Filter out already seen bounties
            for bounty in bounties:
                if bounty['id'] not in self.seen_bounties:
                    self.seen_bounties.append(bounty['id'])
                    new_bounties.append(bounty)
                    send_discord_notification(
                        f"🎯 New Bounty: {bounty['title']}",
                        f"Platform: {platform.platform}\nValue: ${bounty.get('value', 'TBD')}\n{bounty['url']}"
                    )
        
        self.save_seen_bounties()
        return new_bounties
    
    def _scan_github_bounties(self) -> List[Dict]:
        """Scan GitHub for bounty issues"""
        try:
            # Search for issues created in last 6 hours
            six_hours_ago = (datetime.now() - timedelta(hours=6)).strftime('%Y-%m-%d')
            query = f"label:bounty state:open created:>{six_hours_ago}"
            
            url = f"{self.github.base_url}/search/issues"
            params = {'q': query, 'sort': 'created', 'order': 'desc'}
            
            response = self.github.session.get(url, params=params, timeout=10)
            
            if response.status_code == 200:
                data = response.json()
                bounties = []
                
                for item in data.get('items', [])[:10]:  # Top 10
                    bounty = {
                        'id': f"github-{item['id']}",
                        'title': item['title'],
                        'url': item['html_url'],
                        'repo': item['repository_url'].split('/')[-2:],
                        'created_at': item['created_at'],
                        'labels': [label['name'] for label in item.get('labels', [])],
                        'value': self._extract_value_from_body(item.get('body', ''))
                    }
                    bounties.append(bounty)
                
                return bounties
        except Exception as e:
            logger.error(f"Error scanning GitHub: {e}")
        
        return []
    
    def _scan_gitcoin(self) -> List[Dict]:
        """Scan Gitcoin for bounties"""
        try:
            url = "https://api.gitcoin.co/api/v1/bounties/"
            params = {'status': 'open', 'limit': 20}
            
            response = requests.get(url, params=params, timeout=10)
            
            if response.status_code == 200:
                data = response.json()
                bounties = []
                
                for item in data.get('results', [])[:10]:
                    bounty = {
                        'id': f"gitcoin-{item.get('id')}",
                        'title': item.get('title'),
                        'url': item.get('url'),
                        'value': item.get('price', 0),
                        'token': item.get('token_symbol'),
                        'created_at': item.get('created_on')
                    }
                    bounties.append(bounty)
                
                return bounties
        except Exception as e:
            logger.error(f"Error scanning Gitcoin: {e}")
        
        return []
    
    def _extract_value_from_body(self, body: str) -> Optional[str]:
        """Extract bounty value from issue body"""
        import re
        
        # Look for $X or X USD patterns
        patterns = [
            r'\$(\d+(?:,\d{3})*(?:\.\d{2})?)',
            r'(\d+(?:,\d{3})*(?:\.\d{2})?)\s*USD',
        ]
        
        for pattern in patterns:
            match = re.search(pattern, body, re.IGNORECASE)
            if match:
                return match.group(1)
        
        return None


# ============================================================================
# NOTIFICATIONS
# ============================================================================

def send_discord_notification(title: str, message: str):
    """Send notification to Discord webhook"""
    if not DISCORD_WEBHOOK:
        logger.info(f"Discord notification (not sent - no webhook): {title}")
        return
    
    try:
        payload = {
            'embeds': [{
                'title': title,
                'description': message,
                'color': 0x00ff00,
                'timestamp': datetime.now().isoformat()
            }]
        }
        
        response = requests.post(DISCORD_WEBHOOK, json=payload, timeout=10)
        
        if response.status_code in [200, 204]:
            logger.info(f"Discord notification sent: {title}")
        else:
            logger.error(f"Discord notification failed: {response.status_code}")
    except Exception as e:
        logger.error(f"Error sending Discord notification: {e}")


def send_telegram_notification(message: str):
    """Send notification to Telegram"""
    if not TELEGRAM_BOT_TOKEN or not TELEGRAM_CHAT_ID:
        logger.info(f"Telegram notification (not sent - no config): {message[:50]}")
        return
    
    try:
        url = f"https://api.telegram.org/bot{TELEGRAM_BOT_TOKEN}/sendMessage"
        payload = {
            'chat_id': TELEGRAM_CHAT_ID,
            'text': message,
            'parse_mode': 'Markdown'
        }
        
        response = requests.post(url, json=payload, timeout=10)
        
        if response.status_code == 200:
            logger.info("Telegram notification sent")
        else:
            logger.error(f"Telegram notification failed: {response.status_code}")
    except Exception as e:
        logger.error(f"Error sending Telegram notification: {e}")


# ============================================================================
# CLI & MAIN
# ============================================================================

def generate_report(pr_results: Dict, bounty_results: List[Dict]) -> str:
    """Generate automation report"""
    report = []
    report.append("=" * 60)
    report.append("BOUNTY AUTOMATION REPORT")
    report.append(f"Generated: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
    report.append("=" * 60)
    report.append("")
    
    # PR Status
    report.append("📊 PR STATUS")
    report.append("-" * 40)
    
    total_value = 0
    for name, data in pr_results.items():
        status = data['status']
        emoji = "✅" if status['merged'] else "🟡" if status['state'] == 'open' else "❌"
        report.append(f"{emoji} {name}")
        report.append(f"   State: {status['state']} | Merged: {status['merged']}")
        report.append(f"   Days Open: {data['days_since_submission']}")
        report.append(f"   Value: ${data['value_usd']}")
        
        if data['needs_followup']:
            report.append(f"   ⚠️ NEEDS FOLLOW-UP")
        
        total_value += data['value_usd']
        report.append("")
    
    report.append(f"Total Value Tracked: ${total_value}")
    report.append("")
    
    # New Bounties
    report.append("🎯 NEW BOUNTIES DISCOVERED")
    report.append("-" * 40)
    
    if bounty_results:
        for bounty in bounty_results:
            report.append(f"• {bounty['title']}")
            report.append(f"  Value: ${bounty.get('value', 'TBD')}")
            report.append(f"  URL: {bounty['url']}")
            report.append("")
    else:
        report.append("No new bounties discovered")
    
    report.append("")
    report.append("=" * 60)
    
    return "\n".join(report)


def main():
    import argparse
    
    parser = argparse.ArgumentParser(description='Bounty Automation System')
    parser.add_argument('--track-prs', action='store_true', help='Track PR status')
    parser.add_argument('--scan-bounties', action='store_true', help='Scan for new bounties')
    parser.add_argument('--all', action='store_true', help='Run all automation tasks')
    parser.add_argument('--report', action='store_true', help='Generate report')
    
    args = parser.parse_args()
    
    # Initialize
    github = GitHubAPI()
    pr_tracker = PRTracker(github)
    bounty_scanner = BountyScanner(github)
    
    pr_results = {}
    bounty_results = []
    
    # Run tasks
    if args.track_prs or args.all:
        logger.info("Running PR tracking...")
        pr_results = pr_tracker.track_all_prs()
    
    if args.scan_bounties or args.all:
        logger.info("Running bounty scanning...")
        bounty_results = bounty_scanner.scan_all_platforms()
    
    # Generate report
    if args.report or args.all or pr_results or bounty_results:
        report = generate_report(pr_results, bounty_results)
        print(report)
        
        # Save report to file
        with open('bounty_automation_report.txt', 'w') as f:
            f.write(report)


if __name__ == '__main__':
    main()
