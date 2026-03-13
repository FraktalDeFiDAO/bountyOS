# BountyOS Platform Expansion Plan

## Current Platforms (Implemented)

### ✅ GitHub Aggregator
**Status:** Implemented  
**Scanner:** `internal/adapters/scanners/github.go`  
**Bounty Labels Tracked:**
- `algora-bounty`
- `polar`
- `opire`
- `gitpay`
- `issuehunt`
- `bounty`
- `funded`

**Payment Detection:**
- Crypto: USDC, ETH, SOL, USDT
- P2P: Cash App, Venmo
- Fiat: PayPal, Stripe, Wise

---

### ✅ Superteam Earn
**Status:** Implemented (with mock fallback)  
**Scanner:** `internal/adapters/scanners/superteam.go`  
**API:** `https://earn.superteam.fun/api/listings`  
**Focus:** Solana ecosystem bounties  
**Payment:** USDC on Solana, SOL

---

### ✅ Bountycaster
**Status:** Implemented (with mock fallback)  
**Scanner:** `internal/adapters/scanners/bountycaster.go`  
**API:** `https://www.bountycaster.xyz/api/v1/bounties`  
**Focus:** Farcaster ecosystem tasks  
**Payment:** USDC, ETH, DEGEN, OP

---

## High Priority Platforms to Add

### ✅ Immunefi
**Status:** Implemented  
**Scanner:** `internal/adapters/scanners/immunefi.go`  
**API:** `https://immunefi.com/api` (with mock fallback)  
**Focus:** Bug bounties  
**Payment:** Crypto (USDC/ETH)

---

### ✅ Code4rena
**Status:** Implemented  
**Scanner:** `internal/adapters/scanners/code4rena.go`  
**API:** `https://code4rena.com/api`  
**Focus:** Audit contests  
**Payment:** Crypto (USDC)

---

### ✅ Gitcoin
**Status:** Implemented  
**Scanner:** `internal/adapters/scanners/gitcoin.go`  
**API:** `https://api.gitcoin.co/api/v1/bounties`  
**Focus:** Open source bounties  
**Payment:** Crypto (ETH/DAI)

---

### ✅ Polar.sh
**Status:** Implemented  
**Scanner:** `internal/adapters/scanners/polar.go`  
**API:** `https://api.polar.sh/v1`  
**Focus:** OSS funding  
**Payment:** Crypto/Fiat

---

### ✅ Algora
**Status:** Implemented  
**Scanner:** `internal/adapters/scanners/algora.go`  
**API:** `https://algora.io/api/bounties`  
**Focus:** GitHub bounties  
**Payment:** Crypto/Fiat

---

## Medium Priority Platforms

### 📌 IssueHunt
**Priority:** MEDIUM  
**Why:** Established OSS bounty platform  
**Type:** Issue-based bounties  
**API:** GitHub App integration  
**Payment:** Crypto, fiat  
**Implementation:**
- Already tracked via GitHub `issuehunt` label
- Consider direct API if available

**Target Endpoints:**
- `https://issuehunt.io/`
- `https://github.com/apps/issuehunt-oss`

---

### 📌 GitPay
**Priority:** MEDIUM  
**Why:** Collaborative bounty platform  
**Type:** Git issue bounties  
**API:** Unknown  
**Payment:** Crypto, fiat  
**Implementation:**
- Web scraping required
- Already tracked via GitHub `gitpay` label

**Target Endpoints:**
- `https://gitpay.me/`
- `https://gitpay.me/api/` (to verify)

---

### 📌 Opire
**Priority:** MEDIUM  
**Why:** Emerging OSS funding platform  
**Type:** GitHub issue bounties  
**API:** Unknown  
**Payment:** Crypto  
**Implementation:**
- Already tracked via GitHub `opire` label
- Monitor for API availability

**Target Endpoints:**
- `https://opire.dev/`

---

### ✅ Optimism
**Status:** Implemented  
**Scanner:** `internal/adapters/scanners/optimism.go`  
**API:** `https://api.optimism.io`  
**Focus:** Retroactive grants  
**Payment:** Crypto (OP)
---

### 📌 Uniswap Foundation
**Priority:** MEDIUM  
**Why:** Major DeFi protocol grants ($5K-$500K)  
**Type:** DeFi development, tooling  
**API:** Unknown  
**Payment:** UNI, USDC  
**Implementation:**
- Web scraping for grant rounds
- Track application cycles

**Target Endpoints:**
- `https://uniswap.foundation/`

---

### ✅ Solana Colosseum
**Status:** Implemented  
**Scanner:** `internal/adapters/scanners/colosseum.go`  
**API:** `https://api.colosseum.org`  
**Focus:** Hackathons  
**Payment:** Crypto (USDC/SOL)

---

### ✅ Base Ecosystem
**Status:** Implemented  
**Scanner:** `internal/adapters/scanners/base.go`  
**API:** `https://api.base.org`  
**Focus:** Builder grants  
**Payment:** Crypto (USDC)
---

## Lower Priority / Future Consideration

### 📋 HackerOne
**Priority:** LOW (for now)  
**Why:** Traditional bug bounties (not crypto-focused)  
**Type:** Security vulnerabilities  
**API:** Requires partnership  
**Payment:** USD, PayPal, Bitcoin  
**Note:** Consider if expanding beyond crypto focus

---

### 📋 Bugcrowd
**Priority:** LOW  
**Why:** Traditional bug bounties  
**Type:** Security vulnerabilities  
**API:** Requires partnership  
**Payment:** USD  
**Note:** Same as HackerOne

---

### 📋 Intigriti
**Priority:** LOW  
**Why:** European bug bounty platform  
**Type:** Security vulnerabilities  
**API:** Unknown  
**Payment:** EUR, USD  
**Note:** Regional focus

---

### 📋 LaborX
**Priority:** LOW  
**Why:** Crypto freelance platform  
**Type:** Freelance gigs  
**API:** Unknown  
**Payment:** Crypto  
**Note:** More freelance than bounty

---

### 📋 Hyve
**Priority:** LOW  
**Why:** DeFi freelance platform  
**Type:** Freelance gigs  
**API:** Unknown  
**Payment:** Crypto  
**Note:** Niche platform

---

### 📋 CryptoTask
**Priority:** LOW  
**Why:** Freelance platform  
**Type:** Freelance gigs  
**API:** Unknown  
**Payment:** Crypto  
**Note:** Limited traction

---

### 📋 Bondex
**Priority:** LOW  
**Why:** Professional network  
**Type:** Career opportunities  
**API:** Unknown  
**Payment:** Crypto, fiat  
**Note:** More recruiting than bounties

---

### 📋 Molt Ecosystem Platforms
**Priority:** MEDIUM (AI Agent friendly)  
**Why:** AI agent-compatible bounties  
**Platforms:**
- **ugig.net** - $500-$5,000+ SOL gigs
- **Clawlancer** - USDC on Base
- **TheAgentTimes** - BTC sats for content
- **moltlaunch.com** - ETH on Base
- **claw-jobs.com** - Lightning BTC

**Implementation:**
- Web scraping per platform
- Consider CLI/API integrations where available

---

## Implementation Priority Matrix

| Platform | Impact | Effort | ROI | Priority |
|----------|--------|--------|-----|----------|
| Immunefi | 🔥🔥🔥 | 🔴 High | 🔥🔥🔥 | CRITICAL |
| Code4rena | 🔥🔥🔥 | 🟡 Medium | 🔥🔥🔥 | CRITICAL |
| Gitcoin | 🔥🔥 | 🟡 Medium | 🔥🔥 | HIGH |
| Polar.sh | 🔥🔥 | 🟢 Low | 🔥🔥 | HIGH |
| Algora API | 🔥🔥 | 🟢 Low | 🔥🔥 | HIGH |
| IssueHunt | 🟡 | 🟢 Low | 🟡 | MEDIUM |
| Optimism RetroPGF | 🔥🔥 | 🟡 Medium | 🔥🔥 | MEDIUM |
| Solana Grants | 🔥🔥 | 🟡 Medium | 🔥🔥 | MEDIUM |
| Molt Ecosystem | 🟡 | 🟡 Medium | 🟡 | MEDIUM |

---

## Recommended Implementation Order

### Phase 1: Quick Wins (Week 1-2)
1. **Algora API** - Direct API integration (enhance existing GitHub label tracking)
2. **Polar.sh** - API investigation + integration
3. **IssueHunt** - API investigation + integration

### Phase 2: High Value (Week 3-4)
4. **Immunefi** - Web scraper implementation
5. **Code4rena** - Web scraper implementation
6. **Gitcoin** - API/scraper hybrid

### Phase 3: Ecosystem Grants (Week 5-6)
7. **Optimism RetroPGF** - Round tracker
8. **Solana Foundation** - Grant tracker
9. **Base Ecosystem** - Grant tracker

### Phase 4: AI Agent Platforms (Week 7-8)
10. **Molt ecosystem platforms** - Multiple scrapers

---

## Code Structure Recommendations

### New Scanner Template
```go
// internal/adapters/scanners/immunefi.go
type ImmunefiScanner struct {
    client      *http.Client
    baseURL     string
    minBounty   int // Filter low-value bounties
    categories  []string // DeFi, Infrastructure, etc.
}

func (s *ImmunefiScanner) Scan(ctx context.Context) (<-chan core.Bounty, error) {
    ch := make(chan core.Bounty)
    
    go func() {
        defer close(ch)
        // Scrape active programs
        // Extract bounty details
        // Emit bounties with proper scoring
    }()
    
    return ch, nil
}
```

### Scoring Enhancements
Consider adding:
- **TVL at Risk** multiplier for bug bounties
- **Deadline Urgency** for audit contests
- **Platform Reputation** score
- **Historical Payout** data

---

## API Documentation Resources

### Public APIs to Investigate
1. **Gitcoin:** `https://gitcoin.co/api/v1/docs/`
2. **Polar.sh:** Check docs at `https://docs.polar.sh/`
3. **Algora:** `https://algora.io/api/`
4. **Code4rena:** Check for undocumented API via network inspection
5. **Immunefi:** Likely no public API, scraping required

---

## Testing Strategy

### Unit Tests
- Mock HTTP responses for each platform
- Test bounty parsing logic
- Test payment type detection
- Test scoring algorithm

### Integration Tests
- Rate-limited live API calls
- Validate bounty data structure
- Test circuit breaker behavior

### E2E Tests
- Full scanner pipeline
- Storage persistence
- Notification delivery

---

## Success Metrics

### Platform Coverage
- [ ] 10+ platforms integrated
- [ ] 3+ payment types supported
- [ ] 95% uptime on scanners

### Bounty Quality
- [ ] Average bounty value > $500
- [ ] 80%+ crypto payment bounties
- [ ] Real-time detection (< 5 min delay)

### User Value
- [ ] 50+ bounties discovered daily
- [ ] 10+ high-value ($1K+) bounties weekly
- [ ] < 1% false positive rate

---

## Notes

- All new scanners must implement the `core.Scanner` interface
- Use the existing retry/circuit breaker infrastructure
- Follow security best practices (input validation, sanitization)
- Add comprehensive logging with token masking
- Include mock fallbacks for unreliable APIs
- Document API rate limits and respect them

---

**Last Updated:** February 23, 2026  
**Author:** BountyOS Team
