package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"syscall"
	"time"

	"bountyos-v8/internal/adapters/scanners"
	"bountyos-v8/internal/adapters/storage"
	"bountyos-v8/internal/adapters/ui"
	"bountyos-v8/internal/config"
	"bountyos-v8/internal/core"
	"bountyos-v8/internal/notify"
	"bountyos-v8/internal/security"
	"github.com/fatih/color"
)

var logger *security.SecureLogger

func main() {
	configPath := flag.String("config", config.DefaultPath, "Path to config file")
	noUI := flag.Bool("no-ui", false, "Disable terminal UI")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.Load(*configPath)
	if err != nil {
		// Logger not initialized yet; fall back to stderr.
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	headless := strings.EqualFold(os.Getenv("HEADLESS"), "true")
	disableUI := cfg.NoUI || *noUI || headless

	// Initialize secure logger
	logger = security.GetLogger()
	logFile := openLogFile(cfg.LogPath)
	logWriters := []io.Writer{}
	if logFile != nil {
		logWriters = append(logWriters, logFile)
	}

	logToStdout := cfg.LogToStdout
	logToStderr := cfg.LogToStderr
	if !disableUI && cfg.QuietUILogs {
		logToStdout = false
		logToStderr = false
	}

	if logToStdout {
		logWriters = append(logWriters, os.Stdout)
	}
	if logToStderr {
		logWriters = append(logWriters, os.Stderr)
	}
	if len(logWriters) == 0 {
		if disableUI {
			logWriters = append(logWriters, os.Stdout)
		} else {
			logWriters = append(logWriters, io.Discard)
		}
	}
	logger.SetOutput(io.MultiWriter(logWriters...))

	if logFile != nil {
		defer logFile.Close()
	}

	if cfg.DisableRateLimitSleep {
		os.Setenv("BOUNTYOS_DISABLE_RATE_LIMIT_SLEEP", "1")
	}

	core.SetScoringConfig(core.ScoringConfig{
		UrgencyKeywords:    cfg.UrgencyKeywords,
		DevTaskKeywords:    cfg.DevTaskKeywords,
		AutomationKeywords: cfg.AutomationKeywords,
		SecurityKeywords:   cfg.SecurityKeywords,
		AuditKeywords:      cfg.AuditKeywords,
	})
	core.SetPaymentConfig(core.PaymentConfig{
		CryptoCurrencies: cfg.CryptoCurrencies,
		P2PMethods:       cfg.P2PMethods,
		FiatMethods:      cfg.FiatMethods,
	})

	githubToken := cfg.GitHubToken
	logger.RegisterToken(githubToken)
	logger.Info("Starting BountyOS v8: Obsidian with enhanced security")

	// Initialize components
	storage, err := storage.NewSQLiteStorage(cfg.StoragePath)
	if err != nil {
		logger.Error("Failed to initialize storage: %v", err)
		os.Exit(1)
	}
	defer storage.Close()

	pruned, err := storage.PurgeInvalidURLs(ctx, cfg.ValidateLinksHTTP, time.Duration(cfg.LinkValidationTimeout)*time.Second)
	if err != nil {
		logger.Warn("Failed to purge invalid URLs: %v", err)
	} else if pruned > 0 {
		logger.Info("Pruned %d invalid bounties from storage", pruned)
	}

	notifier := notify.NewDesktopNotifier()
	discordWebhook := cfg.DiscordWebhookURL
	discordNotifier := notify.NewDiscordNotifier(discordWebhook)
	if discordWebhook != "" {
		logger.Info("Discord notifications enabled")
	}

	// Initialize and start Web UI
	webUI := ui.NewWebUI(storage, cfg.WebPort, cfg.APIBountiesLimit, cfg.APIStatsLimit, cfg.WebFetchIntervalSeconds, cfg.WebStaticDir)
	if err := webUI.Start(ctx); err != nil {
		logger.Error("Failed to start Web UI: %v", err)
	}
	defer webUI.Stop()

	// Initialize scanners
	enabled := make(map[string]bool)
	for _, name := range cfg.EnabledScanners {
		enabled[strings.ToUpper(strings.TrimSpace(name))] = true
	}
	knownScanners := map[string]bool{
		"GITHUB_AGGREGATOR": true,
		"GITHUB":            true,
		"SUPERTEAM":         true,
		"BOUNTYCASTER":      true,
	}
	for name := range enabled {
		if !knownScanners[name] {
			logger.Warn("Unknown scanner in config: %s", name)
		}
	}

	var scannersList []core.Scanner
	addScanner := func(name string, scanner core.Scanner) {
		if len(enabled) == 0 || enabled[name] {
			scannersList = append(scannersList, scanner)
			return
		}
	}

	githubScanner := scanners.NewGitHubScanner(githubToken, scanners.GitHubScannerConfig{
		Labels:   cfg.GitHubLabels,
		BaseURL:  cfg.GitHubBaseURL,
		PerPage:  cfg.GitHubPerPage,
		MaxPages: cfg.GitHubMaxPages,
	})
	superteamScanner := scanners.NewSuperteamScanner(scanners.SuperteamScannerConfig{
		BaseURL:  cfg.SuperteamBaseURL,
		Statuses: cfg.SuperteamStatuses,
	})
	bountycasterScanner := scanners.NewBountycasterScanner(scanners.BountycasterScannerConfig{
		BaseURL:  cfg.BountycasterBaseURL,
		Statuses: cfg.BountycasterStatuses,
	})

	addScanner("GITHUB_AGGREGATOR", githubScanner)
	addScanner("GITHUB", githubScanner)
	addScanner("SUPERTEAM", superteamScanner)
	addScanner("BOUNTYCASTER", bountycasterScanner)

	if len(scannersList) == 0 {
		logger.Error("No scanners enabled; check ENABLED_SCANNERS in config")
		os.Exit(1)
	}

	// Channel for bounties
	bountyChan := make(chan core.Bounty, 100)

	// Start signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start scanning loop
	go func() {
		// Initial scan
		scanAll(ctx, scannersList, bountyChan)

		ticker := time.NewTicker(time.Duration(cfg.PollIntervalSeconds) * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				scanAll(ctx, scannersList, bountyChan)
			}
		}
	}()

	// Process bounties
	minScore := cfg.MinScore

	go func() {
		for bounty := range bountyChan {
			bounty.URL = security.NormalizeURL(bounty.URL)
			if bounty.URL == "" || !security.ValidateURL(bounty.URL) {
				logger.Warn("Skipping bounty with invalid URL: %s", bounty.URL)
				continue
			}

			if cfg.ValidateLinksHTTP {
				timeout := time.Duration(cfg.LinkValidationTimeout) * time.Second
				checkCtx, cancel := context.WithTimeout(ctx, timeout)
				ok := security.ValidateURLReachable(checkCtx, bounty.URL, timeout)
				cancel()
				if !ok {
					logger.Warn("Skipping bounty with unreachable URL: %s", bounty.URL)
					continue
				}
			}

			bounty.Title = security.SanitizeString(bounty.Title)
			bounty.Platform = security.SanitizeString(bounty.Platform)
			bounty.Reward = security.SanitizeString(bounty.Reward)
			bounty.Currency = security.SanitizeString(bounty.Currency)
			bounty.Description = security.SanitizeString(bounty.Description)

			isNew, err := storage.IsNew(bounty.URL)
			if err != nil {
				logger.Error("Error checking if bounty is new: %v", err)
				continue
			}

			if !isNew {
				continue
			}

			// Calculate score
			bounty.Score = core.CalculateUrgency(&bounty)

			// Save to storage
			if err := storage.Save(bounty); err != nil {
				logger.Error("Error saving bounty: %v", err)
				continue
			}
			webUI.Broadcast(bounty)

			// Send notification if score is high enough
			if bounty.Score >= minScore {
				if err := notifier.Alert(bounty); err != nil {
					logger.Error("Error sending desktop notification: %v", err)
				}
				if discordWebhook != "" {
					if err := discordNotifier.Alert(bounty); err != nil {
						logger.Error("Error sending Discord notification: %v", err)
					}
				}
			}
		}
	}()

	// Display UI if not disabled
	var uiWG sync.WaitGroup
	if !disableUI {
		uiWG.Add(1)
		go func() {
			defer uiWG.Done()
			displayUI(ctx, storage, cfg.UIRefreshSeconds, cfg.TUIRecentLimit)
		}()
	}

	// Wait for shutdown signal
	<-sigChan
	fmt.Println("\nShutting down...")
	cancel()
	uiWG.Wait()
}

func openLogFile(path string) *os.File {
	if strings.TrimSpace(path) == "" {
		return nil
	}
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create log directory %s: %v\n", dir, err)
		return nil
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open log file %s: %v\n", path, err)
		return nil
	}
	return file
}

// scanAll executes all scanners concurrently and waits for them to complete.
// It ensures that all found bounties are sent to the bountyChan before returning.
func scanAll(ctx context.Context, scanners []core.Scanner, bountyChan chan<- core.Bounty) {
	var wg sync.WaitGroup
	for _, scanner := range scanners {
		wg.Add(1)
		go func(s core.Scanner) {
			defer wg.Done()
			ch, err := s.Scan(ctx)
			if err != nil {
				logger.Error("Error scanning %s: %v", s.Name(), err)
				return
			}
			for bounty := range ch {
				bountyChan <- bounty
			}
		}(scanner)
	}
	wg.Wait()
}

func displayUI(ctx context.Context, storage *storage.SQLiteStorage, refreshSeconds int, recentLimit int) {
	// Display header
	green := color.New(color.FgGreen).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	// Initial clear
	fmt.Print("\033[2J\033[H")

	ticker := time.NewTicker(time.Duration(refreshSeconds) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			var sb strings.Builder

			// Move cursor to top-left
			sb.WriteString("\033[H")

			// Header
			sb.WriteString(green("========================================================\n"))
			sb.WriteString(green("   🕷️  BOUNTY OS v8: OBSIDIAN // SNIPER ACTIVE\n"))
			sb.WriteString(green("========================================================\n"))
			sb.WriteString("Press Ctrl+C to exit\n\n")

			// Get recent bounties
			bounties, err := storage.GetRecent(recentLimit)
			if err != nil {
				logger.Error("Error getting bounties: %v", err)
				continue
			}

			// Sort by score
			sort.Slice(bounties, func(i, j int) bool {
				return bounties[i].Score > bounties[j].Score
			})

			// Print table header
			sb.WriteString(fmt.Sprintf("%-8s | %-12s | %-10s | %-15s | %s\n", "SCORE", "PLATFORM", "PAYOUT", "PAYMENT", "TASK"))
			sb.WriteString(strings.Repeat("-", 100) + "\n")

			// Print bounties
			for _, bounty := range bounties {
				scoreStr := fmt.Sprintf("%d", bounty.Score)
				if bounty.Score >= 80 {
					scoreStr = red("⚡ " + scoreStr) // Critical
				} else if bounty.Score >= 50 {
					scoreStr = green(scoreStr) // Good
				} else if bounty.Score >= 30 {
					scoreStr = yellow(scoreStr) // Moderate
				}

				// Payment styling
				payStr := bounty.Reward
				currencyUpper := strings.ToUpper(bounty.Currency)

				if strings.Contains(currencyUpper, "USDC") ||
					strings.Contains(currencyUpper, "USDT") ||
					strings.Contains(currencyUpper, "SOL") ||
					strings.Contains(currencyUpper, "ETH") {
					payStr = cyan(payStr) // Crypto
				} else if strings.Contains(currencyUpper, "CASHAPP") {
					payStr = green(payStr) // Cash App
				} else if strings.Contains(currencyUpper, "PAYPAL") {
					payStr = yellow(payStr) // PayPal
				}

				platform := bounty.Platform
				if len(platform) > 12 {
					platform = platform[:12]
				}

				currency := bounty.Currency
				if len(currency) > 10 {
					currency = currency[:10]
				}

				sb.WriteString(fmt.Sprintf("%-8s | %-12s | %-10s | %-15s | %s\n",
					scoreStr, platform, currency, payStr, bounty.Title[:min(len(bounty.Title), 40)]))
				sb.WriteString(fmt.Sprintf("           └─ LINK: %s\n", bounty.URL))
			}

			sb.WriteString("\n")
			sb.WriteString("Last updated: " + time.Now().Format("15:04:05") + "   ")

			// Output everything at once
			fmt.Print(sb.String())
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
