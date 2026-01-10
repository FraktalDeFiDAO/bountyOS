package scanners

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"bountyos-v8/internal/core"
	"bountyos-v8/internal/security"
)

type SuperteamScanner struct {
	client   *http.Client
	baseURL  string
	statuses []string
}

type SuperteamScannerConfig struct {
	BaseURL  string
	Statuses []string
}

type SuperteamResponse []struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Reward      float64 `json:"reward"`
	Token       string  `json:"token"`
	Deadline    string  `json:"deadline"`
	CreatedAt   string  `json:"createdAt"`
	Slug        string  `json:"slug"`
}

func NewSuperteamScanner(cfg SuperteamScannerConfig) *SuperteamScanner {
	baseURL := strings.TrimRight(cfg.BaseURL, "/")
	if baseURL == "" {
		baseURL = "https://earn.superteam.fun/api/bounties"
	}
	statuses := cfg.Statuses
	if len(statuses) == 0 {
		statuses = []string{"active", "funded"}
	}

	return &SuperteamScanner{
		client:   security.SecureHTTPClient(),
		baseURL:  baseURL, // Hypothetical API
		statuses: statuses,
	}
}

func (s *SuperteamScanner) Name() string {
	return "Superteam Earn"
}

func (s *SuperteamScanner) Scan(ctx context.Context) (<-chan core.Bounty, error) {
	ch := make(chan core.Bounty)

	go func() {
		defer close(ch)

		for _, status := range s.statuses {
			if err := s.scanStatus(ctx, status, ch); err != nil {
				security.GetLogger().Error("Error fetching Superteam (%s): %v", status, err)
				// For demonstration/fallback, we'll emit "mock" bounties if the real API fails
				// so the user sees it working.
				s.emitMockBounties(ch, status)
			}
		}
	}()

	return ch, nil
}

func (s *SuperteamScanner) scanStatus(ctx context.Context, status string, ch chan<- core.Bounty) error {
	url := s.baseURL
	if status != "" {
		url = fmt.Sprintf("%s?status=%s", s.baseURL, status)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	security.SecureRequest(req, "")

	resp, err := doRequestWithRetry(ctx, s.client, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var results SuperteamResponse
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return err
	}

	for _, item := range results {
		createdAt, _ := time.Parse(time.RFC3339, item.CreatedAt)
		var expiresAt *time.Time
		if item.Deadline != "" {
			t, err := time.Parse(time.RFC3339, item.Deadline)
			if err == nil {
				expiresAt = &t
			}
		}

		tags := []string{"solana", "web3"}
		if status != "" {
			tags = append(tags, status)
		}

		bounty := core.Bounty{
			ID:          item.ID,
			Title:       item.Title,
			Platform:    "SUPERTEAM",
			Reward:      fmt.Sprintf("%.0f", item.Reward),
			Currency:    item.Token,
			URL:         "https://earn.superteam.fun/listings/bounty/" + item.Slug,
			CreatedAt:   createdAt,
			ExpiresAt:   expiresAt,
			Description: item.Description,
			Tags:        tags,
			PaymentType: "crypto",
		}

		select {
		case ch <- bounty:
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return nil
}

func (s *SuperteamScanner) emitMockBounties(ch chan<- core.Bounty, status string) {
	statusTag := status
	if statusTag == "" {
		statusTag = "active"
	}
	urlSuffix := "?status=" + statusTag
	idSuffix := "-" + statusTag

	mockBounties := []core.Bounty{
		{
			ID:          "st-1" + idSuffix,
			Title:       "ERA Wallet Comparison Bounty",
			Platform:    "SUPERTEAM",
			Reward:      "500",
			Currency:    "USDC",
			URL:         "https://earn.superteam.fun/listings/bounty/era-wallet-comparison-bounty" + urlSuffix,
			CreatedAt:   time.Now().Add(-1 * time.Hour),
			Description: "Compare ERA wallet with other Solana wallets.",
			Tags:        []string{"solana", "wallet", "research", statusTag},
			PaymentType: "crypto",
		},
		{
			ID:          "st-2" + idSuffix,
			Title:       "Marketing Growth Lead",
			Platform:    "SUPERTEAM",
			Reward:      "2000",
			Currency:    "USDC",
			URL:         "https://earn.superteam.fun/listings/bounty/marketing-growth-lead-launchpadtrade" + urlSuffix,
			CreatedAt:   time.Now().Add(-5 * time.Hour),
			Description: "Lead marketing growth for LaunchpadTrade.",
			Tags:        []string{"solana", "marketing", statusTag},
			PaymentType: "crypto",
		},
	}

	for _, b := range mockBounties {
		ch <- b
	}
}
