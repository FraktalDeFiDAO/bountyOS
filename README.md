# BountyOS Sniper Engine

**Algorithmic bounty hunting system for developers**

BountyOS is a custom-built sniping engine designed to aggregate, filter, and alert you to funded dev tasks before the general public reacts. This tool runs in your terminal, acting as your "Command Center" for finding immediate pay opportunities.

## Features

- **Multi-platform scanning**: Monitors GitHub, Superteam, and Algora for bounties
- **Real-time filtering**: Only shows developer-relevant opportunities
- **Smart deduplication**: Prevents duplicate alerts
- **Keyword-based targeting**: Focuses on "fix", "bug", "script", "api", etc.
- **Urgency detection**: Highlights urgent/important bounties
- **Comprehensive logging**: All activity is logged for analysis
- **Terminal notifications**: Visual and audio alerts for new opportunities

## Architecture

The system targets three specific data streams that hold "Immediate Pay" opportunities:

1. **GitHub Issues**: Filters for issues with `bounty`, `paid`, `reward` labels created in the last 24 hours
2. **Superteam Earn**: Queries their API for "Sprint" bounties (quick turnaround tasks)
3. **Algora**: Scans their GraphQL API for open bounties

**Logic Flow**: `Poller (Go Routines)` -> `Deduplication (Map)` -> `Keyword Filter (Regex)` -> `Alert (Terminal Bell/Print)` -> `You (Manual Execution)`

## Prerequisites

- Go 1.16 or higher
- Git
- GitHub Personal Access Token (optional, for avoiding rate limits)

## Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/bounty-sniper.git
cd bounty-sniper

# Initialize the module
go mod tidy

# Build the executable
go build -o bounty-sniper
```

## Configuration

Edit the constants in `main.go` to customize the behavior:

```go
const (
    PollInterval = 60 * time.Second  // How often to poll for new bounties
    GithubToken = ""                 // Add your GitHub token if you hit rate limits
)
```

## Usage

Run the sniper engine:

```bash
go run main.go
```

Or use the built binary:

```bash
./bounty-sniper
```

## Standard Operating Procedure (SOP)

When the tool beeps, you have roughly 15 minutes to secure the bag:

### 1. The Assessment (2 Minutes)
- Click the link
- **Check Comments**: If there are >3 comments saying "I'm working on this," **ABORT**. It's gone.
- **Check Funding**: Does the issuer have a history?

### 2. The Claim (The "Snipe")
Do not ask "Can I work on this?" **Claim it.**

- **GitHub Comment**: *"I am starting this now. I have experience with [Language]. I will have a PR ready in [Time]."*
- **Superteam**: Join their Discord immediately. Find the specific project channel. Ping the founder: *"Saw the bounty for X. I'm on it. Check my GitHub [Link]."*

### 3. The Work (The "Sprint")
- If it's a bug fix: Fork, Fix, PR.
- If it's a script: Write it, upload to a private Gist, send a video demo. **Never send the code until payment is in escrow or confirmed.**

## Git Workflow

This project follows a structured git workflow:

```bash
# Create a feature branch for enhancements
git checkout -b feature/new-data-source

# Make your changes
# ... edit files ...

# Commit with conventional commit message
git add .
git commit -m "feat: add Algora GraphQL API integration"

# Push to remote
git push origin feature/new-data-source

# Create pull request
```

## Logging

All activity is logged to `bounty_sniper.log` with timestamps and source information for analysis and debugging.

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests if applicable
5. Commit your changes (`git commit -m 'Add some amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

## License

MIT License - See LICENSE file for details.

## Disclaimer

This tool is designed for educational purposes and ethical bounty hunting. Always follow the terms of service of the platforms you interact with, and conduct yourself professionally when claiming bounties.