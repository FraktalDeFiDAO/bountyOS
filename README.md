# BountyOS v8: Obsidian

The most advanced bounty sniping platform for Web3 developers. Automatically scans hundreds of platforms for funded tasks and prioritizes them by payment method and urgency.

## Features

- **Multi-Platform Scanning**: Monitors GitHub, Superteam, Bountycaster, and more
- **Intelligent Scoring**: Prioritizes Crypto > Cash App > PayPal > Other payments
- **Real-time Alerts**: Desktop notifications for high-value opportunities
- **Persistent Storage**: SQLite database to avoid duplicate alerts
- **Terminal UI**: Bloomberg-style dashboard for monitoring
- **Configurable**: YAML configuration for custom settings

## Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Scanners      │    │   Core Logic    │    │   Adapters      │
│                 │    │                 │    │                 │
│ • GitHub        │───▶│ • Scoring       │───▶│ • Storage       │
│ • Superteam     │    │ • Entity        │    │ • Notifications │
│ • Bountycaster  │    │ • Interfaces    │    │ • UI            │
│ • Bug Bounties  │    │                 │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## Payment Priority System

1. **Crypto King** (Tier 0): USDC, USDT, SOL, ETH, BTC, MATIC, AVAX
2. **P2P Premium** (Tier 1): Cash App, Venmo
3. **Fiat Standard** (Tier 2): PayPal, Stripe, Wise
4. **Low Priority**: Everything else

## Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/bountyos-v8.git
cd bountyos-v8

# Initialize Go module
go mod tidy

# Build the application
go build -o obsidian ./cmd/obsidian

# Run the application
./obsidian
```

## 🔒 Security

BountyOS v8 includes comprehensive security features:

- **Secure HTTP Client**: TLS 1.2+ with strong cipher suites
- **Token Masking**: Automatic protection of sensitive tokens in logs
- **Input Validation**: JSON schema validation and XSS protection
- **Rate Limiting**: GitHub API rate limit tracking and enforcement
- **Secure Logging**: Sanitized logs with token protection

**Security Documentation**: See [SECURITY.md](SECURITY.md) for detailed security information.

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
- Algora (GitHub bounties)
- Polar.sh (GitHub bounties)
- Opire (GitHub bounties)
- GitPay (GitHub bounties)
- Superteam Earn (Solana)
- Bountycaster (Farcaster)
- IssueHunt (GitHub bounties)

### Category II: Big Game Hunters (Bug Bounties & Audits)
- Immunefi
- HackenProof
- Code4rena
- Sherlock
- Hats Finance
- Bugcrowd
- HackerOne
- Intigriti

### Category III: Automated Layer (DePIN & Compute)
- Bittensor
- Akash Network
- Render Network
- Golem
- Mysterium
- Sentinel

### Category IV: Freelance Aggregators
- LaborX
- Hyve
- CryptoTask
- Bondex
- WorkX
- Freelancer

## Scam Filter

The system automatically filters out potential scams based on:
- No upfront payment requests
- Verified platform sources
- Established payment methods

## Development

### Adding New Scanners

To add a new data source, implement the `core.Scanner` interface:

```go
type Scanner interface {
    Name() string
    Scan(ctx context.Context) (<-chan core.Bounty, error)
}
```

### Custom Scoring

The scoring algorithm in `internal/core/score.go` can be customized to prioritize different factors.

### Test Speed

To speed up tests that hit mocked HTTP endpoints, you can disable rate limiter sleeps:

```bash
BOUNTYOS_DISABLE_RATE_LIMIT_SLEEP=1 go test ./...
```

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
