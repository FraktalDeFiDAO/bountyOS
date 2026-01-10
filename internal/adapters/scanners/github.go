package scanners

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"bountyos-v8/internal/core"
	"bountyos-v8/internal/security"
)

type GitHubScanner struct {
	client      *http.Client
	token       string
	endpoints   []string
	baseURL     string
	rateLimiter *security.GitHubRateLimiter
	perPage     int
	maxPages    int
}

type GitHubScannerConfig struct {
	Labels   []string
	BaseURL  string
	PerPage  int
	MaxPages int
}

func NewGitHubScanner(token string, cfg GitHubScannerConfig) *GitHubScanner {
	labels := cfg.Labels
	if len(labels) == 0 {
		labels = []string{
			"algora-bounty",
			"polar",
			"opire",
			"gitpay",
			"issuehunt",
			"bounty",
			"funded",
		}
	}
	baseURL := strings.TrimRight(cfg.BaseURL, "/")
	if baseURL == "" {
		baseURL = "https://api.github.com"
	}
	perPage := cfg.PerPage
	if perPage <= 0 || perPage > 100 {
		perPage = 100
	}
	maxPages := cfg.MaxPages
	if maxPages <= 0 {
		maxPages = 10
	}

	return &GitHubScanner{
		client:      security.SecureHTTPClient(),
		token:       token,
		endpoints:   labels,
		baseURL:     baseURL,
		rateLimiter: security.NewGitHubRateLimiter(token),
		perPage:     perPage,
		maxPages:    maxPages,
	}
}

func (s *GitHubScanner) Name() string {
	return "GitHub Aggregator"
}

func (s *GitHubScanner) Scan(ctx context.Context) (<-chan core.Bounty, error) {
	ch := make(chan core.Bounty)

	go func() {
		defer close(ch)

		for _, label := range s.endpoints {
			for page := 1; page <= s.maxPages; page++ {
				if ctx.Err() != nil {
					return
				}

				query := fmt.Sprintf("is:issue is:open label:%s sort:created-desc", label)
				url := fmt.Sprintf("%s/search/issues?q=%s&per_page=%d&page=%d", s.baseURL, strings.ReplaceAll(query, " ", "+"), s.perPage, page)

				req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
				if err != nil {
					security.GetLogger().Error("Error creating request for %s: %v", label, err)
					break
				}

				security.SecureRequest(req, s.token)

				// Check rate limits before making request
				s.rateLimiter.CheckAndWait()

				// Execute request with retries
				resp, err := doRequestWithRetry(ctx, s.client, req)
				if err != nil {
					security.GetLogger().Error("Error fetching %s (page %d): %v", label, page, err)
					break
				}

				// Update rate limiter with response headers
				s.rateLimiter.UpdateFromHeaders(resp)

				// Validate and parse the response
				validatedResponse, err := security.ValidateGitHubResponseFromReader(resp.Body)
				resp.Body.Close()
				if err != nil {
					security.GetLogger().Error("Error validating response for %s (page %d): %v", label, page, err)
					break
				}

				if len(validatedResponse.Items) == 0 {
					break
				}

				for _, item := range validatedResponse.Items {
					createdAt, err := time.Parse(time.RFC3339, item.CreatedAt)
					if err != nil {
						continue
					}

					// Determine reward and currency from labels
					reward := "Funded"
					currency := "USD" // Default
					paymentType := "fiat"
					isFunded := false

					for _, l := range item.Labels {
						name := strings.ToLower(l.Name)
						if strings.Contains(name, "funded") {
							isFunded = true
						}
						if strings.Contains(name, "$") {
							reward = l.Name
							currency = "" // Already has $
						}
						if strings.Contains(name, "usdc") || strings.Contains(name, "eth") || strings.Contains(name, "sol") || strings.Contains(name, "usdt") {
							reward = l.Name
							currency = "" // Label likely has the currency name
							paymentType = "crypto"
						}
					}

					// Check body for payment keywords if not found in labels
					if paymentType == "fiat" {
						bodyLower := strings.ToLower(item.Body)
						if strings.Contains(bodyLower, "usdc") || strings.Contains(bodyLower, "eth") || strings.Contains(bodyLower, "sol") || strings.Contains(bodyLower, "usdt") {
							currency = "USDC/ETH/SOL"
							paymentType = "crypto"
						} else if strings.Contains(bodyLower, "paypal") {
							currency = "PAYPAL"
							paymentType = "fiat"
						} else if strings.Contains(bodyLower, "cash app") || strings.Contains(bodyLower, "cashapp") {
							currency = "CASHAPP"
							paymentType = "p2p"
						}
					}

					// Determine tags
					tags := []string{"active"}
					titleLower := strings.ToLower(item.Title)
					if strings.Contains(titleLower, "urgent") {
						tags = append(tags, "urgent")
					}
					if strings.Contains(titleLower, "fix") || strings.Contains(titleLower, "bug") {
						tags = append(tags, "dev")
					}
					if strings.Contains(titleLower, "script") || strings.Contains(titleLower, "bot") {
						tags = append(tags, "automation")
					}
					if isFunded {
						tags = append(tags, "funded")
					}

					bounty := core.Bounty{
						ID:          item.HTMLURL,
						Title:       item.Title,
						Platform:    "GITHUB/" + strings.ToUpper(label),
						Reward:      reward,
						Currency:    currency,
						URL:         item.HTMLURL,
						CreatedAt:   createdAt,
						Description: item.Body,
						Tags:        tags,
						PaymentType: paymentType,
					}

					select {
					case ch <- bounty:
					case <-ctx.Done():
						return
					}
				}

				if len(validatedResponse.Items) < s.perPage {
					break
				}

				// Rate limiting
				time.Sleep(2 * time.Second)
			}
		}
	}()

	return ch, nil
}
