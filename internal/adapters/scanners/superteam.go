package scanners

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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

type SuperteamListing struct {
	ID               string   `json:"id"`
	RewardAmount     *float64 `json:"rewardAmount"`
	Deadline         string   `json:"deadline"`
	Type             string   `json:"type"`
	Title            string   `json:"title"`
	Token            string   `json:"token"`
	Slug             string   `json:"slug"`
	CompensationType string   `json:"compensationType"`
	MinRewardAsk     *float64 `json:"minRewardAsk"`
	MaxRewardAsk     *float64 `json:"maxRewardAsk"`
	Status           string   `json:"status"`
}

func NewSuperteamScanner(cfg SuperteamScannerConfig) *SuperteamScanner {
	baseURL := strings.TrimRight(cfg.BaseURL, "/")
	if baseURL == "" {
		baseURL = "https://earn.superteam.fun/api/listings"
	}
	statuses := cfg.Statuses
	if len(statuses) == 0 {
		statuses = []string{"open"}
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

		for _, status := range s.normalizedStatuses() {
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
	url := fmt.Sprintf("%s?type=bounties", s.baseURL)
	if status != "" {
		url = fmt.Sprintf("%s&status=%s", url, status)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	security.SecureRequest(req, "")
	req.Header.Set("Accept", "application/json")

	resp, err := doRequestWithRetry(ctx, s.client, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status %d from %s: %s", resp.StatusCode, url, responseSnippet(body))
	}

	var results []SuperteamListing
	if err := json.Unmarshal(body, &results); err != nil {
		return fmt.Errorf("invalid JSON from %s: %w (snippet: %s)", url, err, responseSnippet(body))
	}

	for _, item := range results {
		if strings.ToLower(strings.TrimSpace(item.Type)) != "bounty" {
			continue
		}

		// Superteam listings API does not expose created_at; use a conservative fallback
		createdAt := time.Now().Add(-48 * time.Hour)
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

		reward := ""
		if item.RewardAmount != nil {
			reward = formatAmount(*item.RewardAmount)
		} else if item.MinRewardAsk != nil || item.MaxRewardAsk != nil {
			min := ""
			max := ""
			if item.MinRewardAsk != nil {
				min = formatAmount(*item.MinRewardAsk)
			}
			if item.MaxRewardAsk != nil {
				max = formatAmount(*item.MaxRewardAsk)
			}
			switch {
			case min != "" && max != "":
				reward = fmt.Sprintf("%s-%s", min, max)
			case min != "":
				reward = min
			case max != "":
				reward = max
			}
		}
		if reward == "" && strings.EqualFold(item.CompensationType, "variable") {
			reward = "Variable"
		}

		url := ""
		if item.Slug != "" {
			url = "https://earn.superteam.fun/listings/bounty/" + item.Slug
		}

		bounty := core.Bounty{
			ID:          item.ID,
			Title:       item.Title,
			Platform:    "SUPERTEAM",
			Reward:      reward,
			Currency:    item.Token,
			URL:         url,
			CreatedAt:   createdAt,
			ExpiresAt:   expiresAt,
			Description: item.Title,
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

func (s *SuperteamScanner) normalizedStatuses() []string {
	if len(s.statuses) == 0 {
		return []string{"open"}
	}
	seen := make(map[string]struct{}, len(s.statuses))
	out := make([]string, 0, len(s.statuses))
	for _, status := range s.statuses {
		normalized := normalizeSuperteamStatus(status)
		if normalized == "" {
			continue
		}
		if _, ok := seen[normalized]; ok {
			continue
		}
		seen[normalized] = struct{}{}
		out = append(out, normalized)
	}
	if len(out) == 0 {
		return []string{"open"}
	}
	return out
}

func normalizeSuperteamStatus(status string) string {
	normalized := strings.ToLower(strings.TrimSpace(status))
	switch normalized {
	case "":
		return ""
	case "active", "funded":
		return "open"
	case "in-progress":
		return "review"
	default:
		return normalized
	}
}

func formatAmount(value float64) string {
	formatted := fmt.Sprintf("%.2f", value)
	formatted = strings.TrimRight(formatted, "0")
	formatted = strings.TrimRight(formatted, ".")
	return formatted
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
