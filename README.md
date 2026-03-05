# BountyOS v8: Obsidian

The most advanced bounty sniping platform for Web3 developers. Automatically scans hundreds of platforms for funded tasks and prioritizes them by payment method and urgency.

## Features

- **Multi-Platform Scanning**: Monitors 20+ platforms including GitHub, Superteam, Bountycaster, Immunefi, and more
- **Intelligent Scoring**: Obsidian algorithm prioritizes Crypto > P2P > Fiat payments with urgency detection
- **Real-time Updates**: WebSocket-powered live bounty feed with connection quality monitoring
- **Persistent Storage**: SQLite database with circuit breaker protection to avoid duplicate alerts
- **Modern Web UI**: Vue 3 frontend with glass morphism design and real-time updates
- **Issue Watcher**: Separate service for monitoring GitHub issue comments and payout notifications
- **Resilient Architecture**: Rate limiting, retry logic, and fault tolerance built-in

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                      bountyOS v8                             │
├─────────────────────────────────────────────────────────────┤
│  ┌──────────────┐     ┌──────────────┐     ┌─────────────┐ │
│  │   Scanners   │     │   Scoring    │     │   Storage   │ │
│  │  (20+ plat.) │────▶│  (Obsidian)  │────▶│   (SQLite)  │ │
│  └──────────────┘     └──────────────┘     └─────────────┘ │
│                            │                                │
│                            ▼                                │
│  ┌──────────────────────────────────────────────────────┐   │
│  │              UI Adapter (Web + WebSocket)             │   │
│  └──────────────────────────────────────────────────────┘   │
│                            │                                │
└────────────────────────────┼────────────────────────────────┘
                             │
                             ▼
                    ┌─────────────────┐
                    │  Vue 3 Frontend │
                    │  (Real-time UI) │
                    └─────────────────┘
```

### Project Structure

```
bountyos-v8/
├── cmd/
│   ├── obsidian/           # Main application (TUI + Web server)
│   └── issuewatcher/       # GitHub issue monitoring service
├── internal/
│   ├── core/               # Domain entities, interfaces, scoring
│   ├── adapters/           # External integrations
│   │   ├── scanners/       # Platform bounty scanners (20+)
│   │   ├── storage/        # SQLite persistence
│   │   └── ui/             # Web server + WebSocket
│   ├── config/             # Configuration loading
│   ├── security/           # Security utilities
│   ├── resilience/         # Fault tolerance patterns
│   ├── errors/             # Error types and helpers
│   └── notify/             # Notification adapters
├── web/                    # Vue 3 frontend
│   └── src/
│       ├── components/     # Reusable UI components
│       ├── composables/    # Vue 3 composition functions
│       ├── views/          # Page components
│       ├── stores/         # Pinia state stores
│       └── utils/          # Shared utilities
├── config/                 # Runtime configuration
├── e2e-tests/              # Playwright E2E tests
└── scripts/                # Maintenance utilities
```

## Payment Priority System

**Tier 0 - Crypto King:** USDC, USDT, SOL, ETH, BTC, MATIC, AVAX, ARB, OP  
**Tier 1 - P2P Premium:** Cash App, Venmo  
**Tier 2 - Fiat Standard:** PayPal, Stripe, Wise  
**Tier 3 - Low Priority:** Everything else

### Obsidian Scoring Algorithm

| Factor | Points |
|--------|--------|
| Crypto payment (USDC, SOL, ETH) | +50 |
| P2P payment (Cash App, Venmo) | +45 |
| Fiat payment (PayPal, Stripe) | +25 |
| Urgency keywords (URGENT, ASAP) | +30 |
| Audit tasks | +35 |
| Security tasks | +25 |
| Automation tasks | +20 |
| Dev tasks | +15 |
| Recent (<1 hour) | +40 |
| Platform bonuses | +10 to +30 |

## Installation

### Quick Start

```bash
# Clone the repository
git clone https://github.com/FraktalDeFiDAO/bountyOS.git
cd bountyOS

# Install dependencies
go mod tidy

# Build the application
go build -o obsidian ./cmd/obsidian

# Run with web UI (default)
./obsidian

# Run API-only mode (no UI)
./obsidian -no-ui
```

### Docker Deployment

```bash
# Development with hot reload
podman compose -f docker-compose.dev.yml up --build

# Production
podman compose up -d
```

### Podman Network Troubleshooting (aardvark/netavark)

If startup fails with messages like `Error starting server failed to bind udp listener on 10.89.0.1:53`, clean stale compose network state and recreate:

```bash
podman compose -f docker-compose.dev.yml down --remove-orphans
podman rm -f bountyos-obsidian-dev bountyos-web-dev 2>/dev/null || true
podman network rm bountyos_bountyos-network bountyos-network 2>/dev/null || true
podman network prune -f
podman compose -f docker-compose.dev.yml up --build -d
```

### Frontend Development

```bash
cd web
npm install
npm run dev
```

## CI/CD

This repo includes two GitHub Actions workflows:

- `/.github/workflows/ci.yml` for backend build/tests and frontend build on push/PR/manual runs.
- `/.github/workflows/local-cd.yml` for local Podman deploy via manual dispatch (designed for `act` + host Podman).
- `/.github/workflows/release.yml` for tag-based releases that publish artifacts and push image tags to GHCR.

Security gate:
- CI and Release both enforce a mandatory `gitleaks` secret scan before build/release jobs.

### Run Workflows Locally With act

```bash
# List available workflows/jobs
act -l

# Run CI workflow locally (pull_request event, non-interactive runner image)
act pull_request \
  -W .github/workflows/ci.yml \
  -P ubuntu-latest=catthehacker/ubuntu:act-latest

# Run local CD workflow on your host (required for Podman access)
# This performs a real deploy.
act workflow_dispatch \
  -W .github/workflows/local-cd.yml \
  -P self-hosted=-self-hosted \
  --input stack=prod-ssl \
  --input build=true

# Dry-run release workflow locally (manual dispatch path)
act workflow_dispatch \
  -W .github/workflows/release.yml \
  -P ubuntu-latest=catthehacker/ubuntu:act-latest \
  --input release_tag=v0.0.0-local \
  -n
```

### Publish A Real Release

```bash
git tag v1.0.0
git push origin v1.0.0
```

That tag triggers `release.yml`, which:
- builds and uploads release artifacts (`obsidian-linux-amd64.tar.gz`, `web-dist.tar.gz`, `SHA256SUMS`)
- pushes container images to `ghcr.io/<owner>/bountyos-obsidian:<tag>` and `:sha-<commit>`

## 🔒 Security

BountyOS v8 includes comprehensive security features:

- **Secure HTTP Client**: TLS 1.2+ with strong cipher suites
- **Token Masking**: Automatic protection of sensitive tokens in logs
- **Input Validation**: JSON schema validation and XSS protection
- **Rate Limiting**: GitHub API rate limit tracking and enforcement
- **Secure Logging**: Sanitized logs with token protection

**Security Documentation**: See [SECURITY.md](SECURITY.md) for detailed security information.
**GitHub Bounties Ops Guide**: See [docs/github-bounties.md](docs/github-bounties.md) for the end-to-end GitHub bounty process and troubleshooting.
**GitHub Bounties Runbook**: See [docs/runbook-github-bounties.md](docs/runbook-github-bounties.md) for the short ops runbook.

### Security Configuration

```bash
# Set GitHub token securely
export GITHUB_TOKEN="your_github_personal_access_token"

# Enable debug mode (shows more detailed logs)
export DEBUG=true

# Run with security features
./obsidian
```

**Best Practices**:
- Use tokens with minimal required permissions
- Rotate tokens every 90 days
- Monitor logs regularly for suspicious activity
- Keep dependencies updated

## Documentation

- `SECURITY.md` — Security design and operational guidance
- `docs/github-bounties.md` — Full GitHub bounty pipeline and edge cases
- `docs/runbook-github-bounties.md` — Short ops runbook and decision tree

## Donations

- **ETH (mainnet + ERC-20, sidechains, L2s):** `0xB5a6102b711ADd687b12758e1C72686c434A0e90`
- **BTC:** `bc1qzhe846kmvp8ncyq09yrme4659espealac4hpff`
- **Solana:** `E97N31qNdwUz1jsgccnJJ9Eqkhhe13rhEmtREHrPLwzz`

## Configuration

The app loads `config/config.yaml` by default. You can point to a different file with `-config` and override any key via environment variables.
By default logs are written to `./data/bountyos.log`, and when the TUI is enabled the console is kept clean.

Example `.env` overrides:

```bash
GITHUB_TOKEN=your_github_personal_access_token
DISCORD_WEBHOOK_URL=your_discord_webhook
POLL_INTERVAL_SECONDS=60
MIN_SCORE=60
BOUNTYOS_DISABLE_RATE_LIMIT_SLEEP=1
LOG_PATH=./data/bountyos.log
LOG_TO_STDOUT=false
LOG_TO_STDERR=false
QUIET_UI_LOGS=true
VALIDATE_LINKS_HTTP=true
LINK_VALIDATION_TIMEOUT_SECONDS=5
```

## Usage

The application will start a terminal UI that displays bounties in real-time, sorted by priority score. High-priority bounties trigger desktop notifications.

## Web Frontend (Vue + WS)

The Go server serves the built frontend from `WEB_STATIC_DIR` (default `./web/dist`) and streams new bounties over WebSocket at `/ws`.

Dev (Podman Compose):

```bash
podman compose -f docker-compose.dev.yml up --build
```

This brings up:
- Go API/WebSocket on `http://localhost:12496`
- Vite dev server on `http://localhost:13440`

Build for production:

```bash
cd web
npm run build
```

### Keyboard Controls
- `Ctrl+C`: Exit the application

## Supported Platforms

### Category I: Flash Layer (Hours to 48h)
- **Algora** - GitHub bounties with direct API integration
- **Polar.sh** - Open source funding platform
- **GitPay** - GitHub bounties (via GitHub labels)
- **IssueHunt** - GitHub bounties (via GitHub labels)
- **Superteam Earn** - Solana ecosystem bounties
- **Bountycaster** - Farcaster social bounties
- **Gitcoin** - Open source bounties and grants
- **Proxies.sx** - Web scraping services ($50-$200)

### Category II: Big Game Hunters (Bug Bounties & Audits)
- **Immunefi** - Web3 smart contract bug bounties ($1K-$10M)
- **Code4rena** - Competitive audit contests ($10K-$200K+)

### Category III: DAO & Freelance Platforms
- **Dework** - DAO task boards and bounties
- **CharmVerse** - DAO workspace with funded tasks
- **LaborX** - Crypto freelance platform ($100-$10K+)

### Category IV: Ecosystem Grants
- **Optimism RetroPGF** - Retroactive public goods funding (10M+ OP)
- **Solana Colosseum** - Hackathons ($100K+ prize pools)
- **Base Ecosystem** - Coinbase-backed builder grants ($3K-$30K)
- **Uniswap Foundation** - DeFi protocol grants ($5K-$500K)

### Category V: AI Agent Platforms
- **ugig.net** - AI agent gigs ($500-$5K+ SOL)
- **Clawlancer** - AI + human tasks on Base (USDC)

### Category VI: Web Scraping & Automation
- **Apify** - Web scraping Actors ($100-$30K)

## Scam Filter

The system automatically filters out potential scams based on:
- No upfront payment requests
- Verified platform sources
- Established payment methods

## Development

### Code Quality

**Formatting:**
```bash
# Format all Go files with goimports
./scripts/format-go.sh

# Run linters
go vet ./...
```

**Testing:**
```bash
# Run all tests
go test ./...

# Run tests with race detector
go test -race ./...

# Speed up tests (disable rate limit sleeps)
BOUNTYOS_DISABLE_RATE_LIMIT_SLEEP=1 go test ./...
```

### Adding New Scanners

**Quick Start:** Use the BaseScanner pattern for 60% less code:

```go
package scanners

type MyScanner struct {
    BaseScanner
    // custom fields
}

func NewMyScanner(cfg MyScannerConfig) *MyScanner {
    return &MyScanner{
        BaseScanner: *NewBaseScanner(ScannerOptions{
            Name:    "My Scanner",
            BaseURL: cfg.BaseURL,
        }),
    }
}

func (s *MyScanner) Scan(ctx context.Context) (<-chan core.Bounty, error) {
    ch := make(chan core.Bounty)
    go func() {
        defer close(ch)
        var results []MyAPIResponse
        s.FetchJSON(ctx, url, &results)  // Built-in error handling
        
        for _, item := range results {
            bounty := s.CreateBounty(
                item.ID, item.Title, "PLATFORM",
                item.Reward, item.Currency, item.Description,
                item.CreatedAt, item.ExpiresAt, item.Tags, "crypto",
            )
            ch <- bounty
        }
    }()
    return ch, nil
}
```

**Documentation:** See `REFACTORING_SUMMARY.md` for complete migration guide.

### Custom Scoring

The scoring algorithm in `internal/core/score.go` can be customized to prioritize different factors.

### Frontend Development

**Composables:** Reusable Vue 3 composition functions:

```js
import { useConnection, useRetry, useBountyFilters } from '@/composables'

// Connection management
const { connected, connectionQuality, qualityColor } = useConnection()

// Retry logic
const { retryCount, canRetry, getRetryDelay } = useRetry({ maxRetries: 5 })

// Filtering
const { topBounties, platformStats, search } = useBountyFilters({ bounties })
```

**See:** `web/src/composables/` for all available composables.

## Containers (Podman Compose)

The default container workflow is Podman Compose:

```bash
podman compose -f docker-compose.yml up --build
```

If you prefer Docker, the same file works with:

```bash
docker compose -f docker-compose.yml up --build
```

## License

MIT License - See LICENSE file for details.
