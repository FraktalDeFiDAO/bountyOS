package config

import (
	"errors"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

const DefaultPath = "config/config.yaml"

type Config struct {
	GitHubToken             string   `yaml:"GITHUB_TOKEN"`
	DiscordWebhookURL       string   `yaml:"DISCORD_WEBHOOK_URL"`
	PollIntervalSeconds     int      `yaml:"POLL_INTERVAL_SECONDS"`
	MinScore                int      `yaml:"MIN_SCORE"`
	StoragePath             string   `yaml:"STORAGE_PATH"`
	LogPath                 string   `yaml:"LOG_PATH"`
	LogToStdout             bool     `yaml:"LOG_TO_STDOUT"`
	LogToStderr             bool     `yaml:"LOG_TO_STDERR"`
	QuietUILogs             bool     `yaml:"QUIET_UI_LOGS"`
	ValidateLinksHTTP       bool     `yaml:"VALIDATE_LINKS_HTTP"`
	LinkValidationTimeout   int      `yaml:"LINK_VALIDATION_TIMEOUT_SECONDS"`
	WebStaticDir            string   `yaml:"WEB_STATIC_DIR"`
	WebPort                 int      `yaml:"WEB_PORT"`
	NoUI                    bool     `yaml:"NO_UI"`
	UIRefreshSeconds        int      `yaml:"UI_REFRESH_SECONDS"`
	TUIRecentLimit          int      `yaml:"TUI_RECENT_LIMIT"`
	APIBountiesLimit        int      `yaml:"API_BOUNTIES_LIMIT"`
	APIStatsLimit           int      `yaml:"API_STATS_LIMIT"`
	WebFetchIntervalSeconds int      `yaml:"WEB_FETCH_INTERVAL_SECONDS"`
	DisableRateLimitSleep   bool     `yaml:"BOUNTYOS_DISABLE_RATE_LIMIT_SLEEP"`
	EnabledScanners         []string `yaml:"ENABLED_SCANNERS"`
	GitHubLabels            []string `yaml:"GITHUB_LABELS"`
	GitHubPerPage           int      `yaml:"GITHUB_PER_PAGE"`
	GitHubMaxPages          int      `yaml:"GITHUB_MAX_PAGES"`
	GitHubBaseURL           string   `yaml:"GITHUB_BASE_URL"`
	SuperteamBaseURL        string   `yaml:"SUPERTEAM_BASE_URL"`
	SuperteamStatuses       []string `yaml:"SUPERTEAM_STATUSES"`
	BountycasterBaseURL     string   `yaml:"BOUNTYCASTER_BASE_URL"`
	BountycasterStatuses    []string `yaml:"BOUNTYCASTER_STATUSES"`
	UrgencyKeywords         []string `yaml:"URGENCY_KEYWORDS"`
	DevTaskKeywords         []string `yaml:"DEV_TASK_KEYWORDS"`
	AutomationKeywords      []string `yaml:"AUTOMATION_KEYWORDS"`
	SecurityKeywords        []string `yaml:"SECURITY_KEYWORDS"`
	AuditKeywords           []string `yaml:"AUDIT_KEYWORDS"`
	PaymentPreferences      []string `yaml:"PAYMENT_PREFERENCES"`
	CryptoCurrencies        []string `yaml:"CRYPTO_CURRENCIES"`
	P2PMethods              []string `yaml:"P2P_METHODS"`
	FiatMethods             []string `yaml:"FIAT_METHODS"`
}

func Default() Config {
	return Config{
		PollIntervalSeconds:     60,
		MinScore:                60,
		StoragePath:             "./data/bounties.db",
		LogPath:                 "./data/bountyos.log",
		LogToStdout:             true,
		LogToStderr:             false,
		QuietUILogs:             true,
		ValidateLinksHTTP:       true,
		LinkValidationTimeout:   5,
		WebStaticDir:            "./web/dist",
		WebPort:                 12496,
		UIRefreshSeconds:        5,
		TUIRecentLimit:          15,
		APIBountiesLimit:        50,
		APIStatsLimit:           100,
		WebFetchIntervalSeconds: 5,
		EnabledScanners:         []string{"GITHUB_AGGREGATOR", "SUPERTEAM", "BOUNTYCASTER"},
		GitHubLabels:            []string{"algora-bounty", "polar", "opire", "gitpay", "issuehunt", "bounty", "funded"},
		GitHubPerPage:           100,
		GitHubMaxPages:          10,
		GitHubBaseURL:           "https://api.github.com",
		SuperteamBaseURL:        "https://earn.superteam.fun/api/bounties",
		SuperteamStatuses:       []string{"active", "funded"},
		BountycasterBaseURL:     "https://www.bountycaster.xyz/api/v1/bounties",
		BountycasterStatuses:    []string{"open"},
		UrgencyKeywords:         []string{"URGENT", "ASAP", "CRITICAL", "IMMEDIATE", "EMERGENCY"},
		DevTaskKeywords:         []string{"FIX", "BUG", "API", "INTEGRATION", "SMART CONTRACT", "BLOCKCHAIN"},
		AutomationKeywords:      []string{"SCRIPT", "BOT"},
		SecurityKeywords:        []string{"SECURITY", "VULNERABILITY", "PENTEST", "HACK", "EXPLOIT"},
		AuditKeywords:           []string{"AUDIT"},
		PaymentPreferences:      []string{"USDC", "SOL", "ETH", "CASHAPP", "VENMO", "PAYPAL", "STRIPE", "WISE"},
		CryptoCurrencies:        []string{"USDC", "USDT", "SOL", "ETH", "BTC", "MATIC", "AVAX", "ARB", "OP"},
		P2PMethods:              []string{"CASHAPP", "VENMO", "CASH APP"},
		FiatMethods:             []string{"USD", "PAYPAL", "STRIPE", "WISE"},
	}
}

func Load(path string) (*Config, error) {
	cfg := Default()
	if path == "" {
		path = DefaultPath
	}

	if _, err := os.Stat(path); err == nil {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return nil, err
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	applyEnvOverrides(&cfg)
	normalize(&cfg)
	return &cfg, nil
}

func applyEnvOverrides(cfg *Config) {
	setString(&cfg.GitHubToken, "GITHUB_TOKEN")
	setString(&cfg.DiscordWebhookURL, "DISCORD_WEBHOOK_URL")
	setInt(&cfg.PollIntervalSeconds, "POLL_INTERVAL_SECONDS")
	setInt(&cfg.MinScore, "MIN_SCORE")
	setString(&cfg.StoragePath, "STORAGE_PATH")
	setString(&cfg.LogPath, "LOG_PATH")
	setBool(&cfg.LogToStdout, "LOG_TO_STDOUT")
	setBool(&cfg.LogToStderr, "LOG_TO_STDERR")
	setBool(&cfg.QuietUILogs, "QUIET_UI_LOGS")
	setBool(&cfg.ValidateLinksHTTP, "VALIDATE_LINKS_HTTP")
	setInt(&cfg.LinkValidationTimeout, "LINK_VALIDATION_TIMEOUT_SECONDS")
	setString(&cfg.WebStaticDir, "WEB_STATIC_DIR")
	setInt(&cfg.WebPort, "WEB_PORT")
	setBool(&cfg.NoUI, "NO_UI")
	setInt(&cfg.UIRefreshSeconds, "UI_REFRESH_SECONDS")
	setInt(&cfg.TUIRecentLimit, "TUI_RECENT_LIMIT")
	setInt(&cfg.APIBountiesLimit, "API_BOUNTIES_LIMIT")
	setInt(&cfg.APIStatsLimit, "API_STATS_LIMIT")
	setInt(&cfg.WebFetchIntervalSeconds, "WEB_FETCH_INTERVAL_SECONDS")
	setBool(&cfg.DisableRateLimitSleep, "BOUNTYOS_DISABLE_RATE_LIMIT_SLEEP")
	setList(&cfg.EnabledScanners, "ENABLED_SCANNERS")
	setList(&cfg.GitHubLabels, "GITHUB_LABELS")
	setInt(&cfg.GitHubPerPage, "GITHUB_PER_PAGE")
	setInt(&cfg.GitHubMaxPages, "GITHUB_MAX_PAGES")
	setString(&cfg.GitHubBaseURL, "GITHUB_BASE_URL")
	setString(&cfg.SuperteamBaseURL, "SUPERTEAM_BASE_URL")
	setList(&cfg.SuperteamStatuses, "SUPERTEAM_STATUSES")
	setString(&cfg.BountycasterBaseURL, "BOUNTYCASTER_BASE_URL")
	setList(&cfg.BountycasterStatuses, "BOUNTYCASTER_STATUSES")
	setList(&cfg.UrgencyKeywords, "URGENCY_KEYWORDS")
	setList(&cfg.DevTaskKeywords, "DEV_TASK_KEYWORDS")
	setList(&cfg.AutomationKeywords, "AUTOMATION_KEYWORDS")
	setList(&cfg.SecurityKeywords, "SECURITY_KEYWORDS")
	setList(&cfg.AuditKeywords, "AUDIT_KEYWORDS")
	setList(&cfg.PaymentPreferences, "PAYMENT_PREFERENCES")
	setList(&cfg.CryptoCurrencies, "CRYPTO_CURRENCIES")
	setList(&cfg.P2PMethods, "P2P_METHODS")
	setList(&cfg.FiatMethods, "FIAT_METHODS")
}

func normalize(cfg *Config) {
	defaults := Default()

	if cfg.PollIntervalSeconds <= 0 {
		cfg.PollIntervalSeconds = defaults.PollIntervalSeconds
	}
	if cfg.MinScore <= 0 {
		cfg.MinScore = defaults.MinScore
	}
	if cfg.StoragePath == "" {
		cfg.StoragePath = defaults.StoragePath
	}
	if cfg.LogPath == "" {
		cfg.LogPath = defaults.LogPath
	}
	if cfg.WebStaticDir == "" {
		cfg.WebStaticDir = defaults.WebStaticDir
	}
	if cfg.LinkValidationTimeout <= 0 {
		cfg.LinkValidationTimeout = defaults.LinkValidationTimeout
	}
	if cfg.WebPort <= 0 {
		cfg.WebPort = defaults.WebPort
	}
	if cfg.UIRefreshSeconds <= 0 {
		cfg.UIRefreshSeconds = defaults.UIRefreshSeconds
	}
	if cfg.TUIRecentLimit <= 0 {
		cfg.TUIRecentLimit = defaults.TUIRecentLimit
	}
	if cfg.APIBountiesLimit <= 0 {
		cfg.APIBountiesLimit = defaults.APIBountiesLimit
	}
	if cfg.APIStatsLimit <= 0 {
		cfg.APIStatsLimit = defaults.APIStatsLimit
	}
	if cfg.WebFetchIntervalSeconds <= 0 {
		cfg.WebFetchIntervalSeconds = defaults.WebFetchIntervalSeconds
	}

	cfg.EnabledScanners = normalizeUpperList(coalesceList(cfg.EnabledScanners, defaults.EnabledScanners))
	cfg.GitHubLabels = normalizeTrimList(coalesceList(cfg.GitHubLabels, defaults.GitHubLabels))
	cfg.GitHubPerPage = clampInt(cfg.GitHubPerPage, 1, 100, defaults.GitHubPerPage)
	cfg.GitHubMaxPages = clampInt(cfg.GitHubMaxPages, 1, 100, defaults.GitHubMaxPages)
	cfg.GitHubBaseURL = strings.TrimRight(firstNonEmpty(cfg.GitHubBaseURL, defaults.GitHubBaseURL), "/")
	cfg.SuperteamBaseURL = strings.TrimRight(firstNonEmpty(cfg.SuperteamBaseURL, defaults.SuperteamBaseURL), "/")
	cfg.BountycasterBaseURL = strings.TrimRight(firstNonEmpty(cfg.BountycasterBaseURL, defaults.BountycasterBaseURL), "/")
	cfg.SuperteamStatuses = normalizeLowerList(coalesceList(cfg.SuperteamStatuses, defaults.SuperteamStatuses))
	cfg.BountycasterStatuses = normalizeLowerList(coalesceList(cfg.BountycasterStatuses, defaults.BountycasterStatuses))
	cfg.UrgencyKeywords = normalizeUpperList(coalesceList(cfg.UrgencyKeywords, defaults.UrgencyKeywords))
	cfg.DevTaskKeywords = normalizeUpperList(coalesceList(cfg.DevTaskKeywords, defaults.DevTaskKeywords))
	cfg.AutomationKeywords = normalizeUpperList(coalesceList(cfg.AutomationKeywords, defaults.AutomationKeywords))
	cfg.SecurityKeywords = normalizeUpperList(coalesceList(cfg.SecurityKeywords, defaults.SecurityKeywords))
	cfg.AuditKeywords = normalizeUpperList(coalesceList(cfg.AuditKeywords, defaults.AuditKeywords))
	cfg.PaymentPreferences = normalizeUpperList(coalesceList(cfg.PaymentPreferences, defaults.PaymentPreferences))

	if len(cfg.CryptoCurrencies) == 0 && len(cfg.P2PMethods) == 0 && len(cfg.FiatMethods) == 0 {
		cfg.CryptoCurrencies = defaults.CryptoCurrencies
		cfg.P2PMethods = defaults.P2PMethods
		cfg.FiatMethods = defaults.FiatMethods
		if len(cfg.PaymentPreferences) > 0 {
			derivePaymentTiers(cfg)
		}
	} else {
		cfg.CryptoCurrencies = normalizeUpperList(coalesceList(cfg.CryptoCurrencies, defaults.CryptoCurrencies))
		cfg.P2PMethods = normalizeUpperList(coalesceList(cfg.P2PMethods, defaults.P2PMethods))
		cfg.FiatMethods = normalizeUpperList(coalesceList(cfg.FiatMethods, defaults.FiatMethods))
	}

	cfg.DevTaskKeywords = removeOverlap(cfg.DevTaskKeywords, cfg.AutomationKeywords)
	cfg.SecurityKeywords = removeOverlap(cfg.SecurityKeywords, cfg.AuditKeywords)
}

func derivePaymentTiers(cfg *Config) {
	knownCrypto := map[string]bool{
		"USDC": true, "USDT": true, "SOL": true, "ETH": true, "BTC": true, "MATIC": true, "AVAX": true, "ARB": true, "OP": true,
	}
	knownP2P := map[string]bool{"CASHAPP": true, "VENMO": true, "CASH APP": true}

	var crypto []string
	var p2p []string
	var fiat []string

	for _, item := range cfg.PaymentPreferences {
		upper := strings.ToUpper(strings.TrimSpace(item))
		if upper == "" {
			continue
		}
		if knownCrypto[upper] {
			crypto = append(crypto, upper)
		} else if knownP2P[upper] {
			p2p = append(p2p, upper)
		} else {
			fiat = append(fiat, upper)
		}
	}

	if len(crypto) > 0 {
		cfg.CryptoCurrencies = crypto
	}
	if len(p2p) > 0 {
		cfg.P2PMethods = p2p
	}
	if len(fiat) > 0 {
		cfg.FiatMethods = fiat
	}
}

func setString(target *string, key string) {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		*target = value
	}
}

func setInt(target *int, key string) {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			*target = parsed
		}
	}
}

func setBool(target *bool, key string) {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			*target = parsed
		}
	}
}

func setList(target *[]string, key string) {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		*target = splitList(value)
	}
}

func splitList(value string) []string {
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}

func coalesceList(list []string, fallback []string) []string {
	if len(list) == 0 {
		return append([]string(nil), fallback...)
	}
	return list
}

func normalizeUpperList(list []string) []string {
	seen := make(map[string]struct{}, len(list))
	out := make([]string, 0, len(list))
	for _, item := range list {
		trimmed := strings.ToUpper(strings.TrimSpace(item))
		if trimmed == "" {
			continue
		}
		if _, ok := seen[trimmed]; ok {
			continue
		}
		seen[trimmed] = struct{}{}
		out = append(out, trimmed)
	}
	return out
}

func normalizeLowerList(list []string) []string {
	seen := make(map[string]struct{}, len(list))
	out := make([]string, 0, len(list))
	for _, item := range list {
		trimmed := strings.ToLower(strings.TrimSpace(item))
		if trimmed == "" {
			continue
		}
		if _, ok := seen[trimmed]; ok {
			continue
		}
		seen[trimmed] = struct{}{}
		out = append(out, trimmed)
	}
	return out
}

func normalizeTrimList(list []string) []string {
	seen := make(map[string]struct{}, len(list))
	out := make([]string, 0, len(list))
	for _, item := range list {
		trimmed := strings.TrimSpace(item)
		if trimmed == "" {
			continue
		}
		if _, ok := seen[trimmed]; ok {
			continue
		}
		seen[trimmed] = struct{}{}
		out = append(out, trimmed)
	}
	return out
}

func removeOverlap(base []string, blacklist []string) []string {
	if len(base) == 0 || len(blacklist) == 0 {
		return base
	}
	deny := make(map[string]struct{}, len(blacklist))
	for _, item := range blacklist {
		deny[item] = struct{}{}
	}
	out := make([]string, 0, len(base))
	for _, item := range base {
		if _, ok := deny[item]; ok {
			continue
		}
		out = append(out, item)
	}
	return out
}

func clampInt(value, min, max, fallback int) int {
	if value < min || value > max {
		return fallback
	}
	return value
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
