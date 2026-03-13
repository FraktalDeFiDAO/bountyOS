#!/usr/bin/env python3
"""
Bounty Verification Tool
Purpose: Verify bounty URLs and generate comprehensive audit reports
Usage: python3 verify_bounties.py [input_file] [output_dir]
"""

import sys
import os
import json
import requests
from datetime import datetime
from pathlib import Path
from concurrent.futures import ThreadPoolExecutor, as_completed
import re

class BountyVerifier:
    def __init__(self, output_dir=None):
        self.output_dir = output_dir or f"./audit_trails/{datetime.now().strftime('%Y-%m-%d')}"
        self.results = []
        self.stats = {
            'total': 0,
            'http_200': 0,
            'http_404': 0,
            'http_other': 0,
            'open_issues': 0,
            'closed_issues': 0
        }
        
        # Create output directories
        Path(f"{self.output_dir}/http_responses").mkdir(parents=True, exist_ok=True)
        Path(f"{self.output_dir}/screenshots").mkdir(parents=True, exist_ok=True)
        Path(f"{self.output_dir}/verification_logs").mkdir(parents=True, exist_ok=True)
    
    def verify_url(self, url: str, bounty_name: str) -> dict:
        """Verify a single bounty URL"""
        timestamp = datetime.now().strftime('%Y-%m-%d_%H-%M-%S')
        issue_num = re.search(r'issues/(\d+)', url)
        issue_num = issue_num.group(1) if issue_num else 'unknown'
        
        result = {
            'bounty_name': bounty_name,
            'url': url,
            'issue_num': issue_num,
            'timestamp': timestamp,
            'http_status': None,
            'http_status_text': None,
            'issue_state': None,
            'title': None,
            'evidence_file': None,
            'confidence': 'HIGH',
            'notes': []
        }
        
        try:
            # Get HTTP headers
            response = requests.head(url, timeout=10, allow_redirects=True)
            result['http_status'] = response.status_code
            result['http_status_text'] = f"{response.status_code} {response.reason}"
            
            # Save evidence
            evidence_file = f"{self.output_dir}/http_responses/{bounty_name}_{issue_num}_{timestamp}.txt"
            with open(evidence_file, 'w') as f:
                f.write(f"URL: {url}\n")
                f.write(f"Timestamp: {timestamp}\n")
                f.write(f"HTTP Status: {response.status_code} {response.reason}\n")
                f.write(f"Headers:\n")
                for key, value in response.headers.items():
                    f.write(f"  {key}: {value}\n")
            
            result['evidence_file'] = evidence_file
            
            # Categorize status
            if response.status_code == 200:
                self.stats['http_200'] += 1
                
                # Try to get issue state from GitHub API or HTML
                if 'github.com' in url:
                    try:
                        # Get full page to check issue state
                        page_response = requests.get(url, timeout=10)
                        if 'Open' in page_response.text and 'issue is open' in page_response.text.lower():
                            result['issue_state'] = 'Open'
                            self.stats['open_issues'] += 1
                        elif 'Closed' in page_response.text:
                            result['issue_state'] = 'Closed'
                            self.stats['closed_issues'] += 1
                    except:
                        result['issue_state'] = 'Unknown'
                        result['notes'].append('Could not determine issue state')
                
            elif response.status_code == 404:
                self.stats['http_404'] += 1
                result['notes'].append('⚠️ CRITICAL: Issue not found - requires manual verification')
                result['confidence'] = 'MEDIUM'
                
            else:
                self.stats['http_other'] += 1
                result['notes'].append(f'Unexpected status: {response.status_code}')
                result['confidence'] = 'MEDIUM'
            
            self.stats['total'] += 1
            
        except requests.exceptions.RequestException as e:
            result['http_status'] = 'ERROR'
            result['http_status_text'] = str(e)
            result['notes'].append(f'Verification failed: {e}')
            result['confidence'] = 'LOW'
            self.stats['total'] += 1
        
        self.results.append(result)
        return result
    
    def verify_batch(self, bounties: list, max_workers: int = 5):
        """Verify multiple bounties in parallel"""
        with ThreadPoolExecutor(max_workers=max_workers) as executor:
            futures = {
                executor.submit(self.verify_url, url, name): (name, url) 
                for name, url in bounties
            }
            
            for future in as_completed(futures):
                name, url = futures[future]
                try:
                    result = future.result()
                    status = result['http_status_text'] or result['http_status']
                    print(f"✅ {name}: {status}" if result['http_status'] == 200 else f"⚠️ {name}: {status}")
                except Exception as e:
                    print(f"❌ {name}: Error - {e}")
    
    def generate_report(self) -> str:
        """Generate markdown verification report"""
        report = []
        report.append("# 🔍 BOUNTY VERIFICATION REPORT")
        report.append(f"\n**Generated:** {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
        report.append(f"\n**Output Directory:** {self.output_dir}\n")
        
        # Summary
        report.append("## 📊 SUMMARY\n")
        report.append(f"| Metric | Count |")
        report.append(f"|--------|-------|")
        report.append(f"| **Total Verified** | {self.stats['total']} |")
        report.append(f"| **HTTP 200 (Exists)** | {self.stats['http_200']} |")
        report.append(f"| **HTTP 404 (Not Found)** | {self.stats['http_404']} |")
        report.append(f"| **HTTP Other** | {self.stats['http_other']} |")
        report.append(f"| **Open Issues** | {self.stats['open_issues']} |")
        report.append(f"| **Closed Issues** | {self.stats['closed_issues']} |")
        
        if self.stats['total'] > 0:
            percent = (self.stats['http_200'] * 100) // self.stats['total']
            report.append(f"\n**Verification Rate:** {percent}%\n")
        
        # Detailed results
        report.append("\n## 📋 DETAILED RESULTS\n")
        report.append("| # | Bounty | URL | HTTP Status | Issue State | Evidence | Confidence | Notes |")
        report.append("|---|--------|-----|-------------|-------------|----------|------------|-------|")
        
        for i, result in enumerate(self.results, 1):
            status_emoji = "✅" if result['http_status'] == 200 else "❌" if result['http_status'] == 404 else "⚠️"
            evidence_link = f"[View]({result['evidence_file']})" if result['evidence_file'] else "N/A"
            notes = "; ".join(result['notes']) if result['notes'] else "-"
            
            report.append(
                f"| {i} | {result['bounty_name']} | [{result['issue_num']}]({result['url']}) | "
                f"{status_emoji} {result['http_status_text']} | {result['issue_state'] or '-'} | "
                f"{evidence_link} | {result['confidence']} | {notes} |"
            )
        
        # Evidence index
        report.append("\n## 📁 EVIDENCE FILES\n")
        report.append(f"All evidence files saved to: `{self.output_dir}/http_responses/`\n")
        
        # JSON export
        json_file = f"{self.output_dir}/verification_results.json"
        with open(json_file, 'w') as f:
            json.dump({
                'timestamp': datetime.now().isoformat(),
                'stats': self.stats,
                'results': self.results
            }, f, indent=2, default=str)
        
        report.append(f"\n**JSON Export:** `{json_file}`\n")
        
        return "\n".join(report)
    
    def load_bounties_from_file(self, filepath: str) -> list:
        """Load bounty URLs from TSV or text file"""
        bounties = []
        
        with open(filepath, 'r') as f:
            for line in f:
                line = line.strip()
                if not line or line.startswith('#'):
                    continue
                
                # Try TSV format: name\turl
                if '\t' in line:
                    parts = line.split('\t')
                    if len(parts) >= 2:
                        bounties.append((parts[0].strip(), parts[1].strip()))
                # Try pipe format: name|url
                elif '|' in line:
                    parts = line.split('|')
                    if len(parts) >= 2:
                        bounties.append((parts[0].strip(), parts[1].strip()))
                # Try markdown link format: [name](url)
                elif '](' in line:
                    match = re.search(r'\[([^\]]+)\]\(([^)]+)\)', line)
                    if match:
                        bounties.append((match.group(1), match.group(2)))
        
        return bounties


def main():
    if len(sys.argv) < 2:
        print("Usage: python3 verify_bounties.py [input_file] [output_dir]")
        print("  input_file: TSV/text file with bounty URLs (default: bounty_urls.txt)")
        print("  output_dir: Directory for evidence files (default: ./audit_trails/YYYY-MM-DD)")
        sys.exit(1)
    
    input_file = sys.argv[1] if len(sys.argv) > 1 else "bounty_urls.txt"
    output_dir = sys.argv[2] if len(sys.argv) > 2 else None
    
    if not os.path.exists(input_file):
        print(f"Error: Input file not found: {input_file}")
        sys.exit(1)
    
    print("=== BOUNTY VERIFICATION TOOL ===")
    print(f"Input: {input_file}")
    
    verifier = BountyVerifier(output_dir)
    
    # Load bounties
    bounties = verifier.load_bounties_from_file(input_file)
    print(f"Loaded {len(bounties)} bounties")
    
    # Verify
    print("\nVerifying bounties...")
    verifier.verify_batch(bounties)
    
    # Generate report
    print("\nGenerating report...")
    report = verifier.generate_report()
    
    # Save report
    report_file = f"{verifier.output_dir}/verification_report.md"
    with open(report_file, 'w') as f:
        f.write(report)
    
    print(f"\n✅ Verification complete!")
    print(f"Report: {report_file}")
    print(f"Evidence: {verifier.output_dir}/http_responses/")
    
    # Print summary
    print(f"\n=== SUMMARY ===")
    print(f"Total: {verifier.stats['total']}")
    print(f"HTTP 200: {verifier.stats['http_200']}")
    print(f"HTTP 404: {verifier.stats['http_404']}")
    print(f"Other: {verifier.stats['http_other']}")


if __name__ == '__main__':
    main()
