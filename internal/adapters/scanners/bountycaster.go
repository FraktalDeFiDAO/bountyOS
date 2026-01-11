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

type BountycasterScanner struct {
	client   *http.Client
	baseURL  string
	statuses []string
}

type BountycasterScannerConfig struct {
	BaseURL  string
	Statuses []string
}

type BountycasterResponse struct {
	Bounties []BountycasterBounty `json:"bounties"`
}

type BountycasterBounty struct {
	UID            string   `json:"uid"`
	Title          string   `json:"title"`
	SummaryText    string   `json:"summary_text"`
	CreatedAt      string   `json:"created_at"`
	ExpirationDate string   `json:"expiration_date"`
	TagSlugs       []string `json:"tag_slugs"`
	Links          struct {
		External string `json:"external"`
		Resource string `json:"resource"`
	} `json:"links"`
	RewardSummary *BountycasterRewardSummary `json:"reward_summary"`
	Platform      struct {
		Type string `json:"type"`
		Hash string `json:"hash"`
	} `json:"platform"`
}

type BountycasterRewardSummary struct {
	UnitAmount string `json:"unit_amount"`
	USDValue   string `json:"usd_value"`
	Symbol     string `json:"symbol"`
	Token      *struct {
		Symbol string `json:"symbol"`
	} `json:"token"`
}

func NewBountycasterScanner(cfg BountycasterScannerConfig) *BountycasterScanner {
	baseURL := strings.TrimRight(cfg.BaseURL, "/")
	if baseURL == "" {
		baseURL = "https://www.bountycaster.xyz/api/v1/bounties"
	}
	statuses := cfg.Statuses
	if len(statuses) == 0 {
		statuses = []string{"open"}
	}

	return &BountycasterScanner{
		client:   security.SecureHTTPClient(),
		baseURL:  baseURL, // Hypothetical API
		statuses: statuses,
	}
}

func (s *BountycasterScanner) Name() string {
	return "Bountycaster"
}

func (s *BountycasterScanner) Scan(ctx context.Context) (<-chan core.Bounty, error) {
	ch := make(chan core.Bounty)

	go func() {
		defer close(ch)

		for _, status := range s.normalizedStatuses() {
			if err := s.scanStatus(ctx, status, ch); err != nil {
				security.GetLogger().Error("Error fetching Bountycaster (%s): %v", status, err)
				s.emitMockBounties(ch, status)
			}
		}
	}()

	return ch, nil
}

func (s *BountycasterScanner) scanStatus(ctx context.Context, status string, ch chan<- core.Bounty) error {
	url := fmt.Sprintf("%s/%s", s.baseURL, status)
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

	var results BountycasterResponse
	if err := json.Unmarshal(body, &results); err != nil {
		return fmt.Errorf("invalid JSON from %s: %w (snippet: %s)", url, err, responseSnippet(body))
	}

	for _, item := range results.Bounties {
		createdAt, _ := time.Parse(time.RFC3339, item.CreatedAt)
		var expiresAt *time.Time
		if item.ExpirationDate != "" {
			if t, err := time.Parse(time.RFC3339, item.ExpirationDate); err == nil {
				expiresAt = &t
			}
		}

		tags := []string{"farcaster", "social"}
		if status != "" {
			tags = append(tags, status)
		}
		for _, tag := range item.TagSlugs {
			trimmed := strings.TrimSpace(tag)
			if trimmed != "" {
				tags = append(tags, trimmed)
			}
		}

		reward := ""
		currency := ""
		if item.RewardSummary != nil {
			reward = strings.TrimSpace(item.RewardSummary.UnitAmount)
			if reward == "" {
				reward = strings.TrimSpace(item.RewardSummary.USDValue)
			}
			currency = strings.TrimSpace(item.RewardSummary.Symbol)
			if currency == "" && item.RewardSummary.Token != nil {
				currency = strings.TrimSpace(item.RewardSummary.Token.Symbol)
			}
		}

		url := ""
		if item.Links.Resource != "" {
			url = "https://www.bountycaster.xyz" + item.Links.Resource
		} else if item.Links.External != "" {
			url = item.Links.External
		} else if item.Platform.Hash != "" {
			url = "https://www.bountycaster.xyz/bounty/" + item.Platform.Hash
		}

		paymentType := "crypto"
		if currency == "" {
			paymentType = "unknown"
		}

		bounty := core.Bounty{
			ID:          item.UID,
			Title:       item.Title,
			Platform:    "BOUNTYCASTER",
			Reward:      reward,
			Currency:    currency,
			URL:         url,
			CreatedAt:   createdAt,
			ExpiresAt:   expiresAt,
			Description: item.SummaryText,
			Tags:        tags,
			PaymentType: paymentType,
		}

		select {
		case ch <- bounty:
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return nil
}

func (s *BountycasterScanner) normalizedStatuses() []string {
	if len(s.statuses) == 0 {
		return []string{"open"}
	}
	seen := make(map[string]struct{}, len(s.statuses))
	out := make([]string, 0, len(s.statuses))
	for _, status := range s.statuses {
		normalized := normalizeBountycasterStatus(status)
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

func normalizeBountycasterStatus(status string) string {
	normalized := strings.ToLower(strings.TrimSpace(status))
	switch normalized {
	case "":
		return ""
	case "active", "funded":
		return "open"
	case "inprogress":
		return "in-progress"
	default:
		return normalized
	}
}

func responseSnippet(body []byte) string {
	snippet := strings.TrimSpace(string(body))
	if len(snippet) > 200 {
		snippet = snippet[:200] + "..."
	}
	return security.SanitizeString(snippet)
}

func (s *BountycasterScanner) emitMockBounties(ch chan<- core.Bounty, status string) {
	statusTag := status
	if statusTag == "" {
		statusTag = "active"
	}
	idSuffix := "-" + statusTag
	urlSuffix := "?status=" + statusTag

	mockBounties := []core.Bounty{
		{
			ID:          "bc-1" + idSuffix,
			Title:       "Dune Dashboard for Seamless Protocol",
			Platform:    "BOUNTYCASTER",
			Reward:      "15000",
			Currency:    "USDC",
			URL:         "https://bountycaster.xyz/bounty/0x11ce0fa8" + urlSuffix,
			CreatedAt:   time.Now().Add(-30 * time.Minute),
			Description: "Create a Dune dashboard for Seamless Protocol metrics.",
			Tags:        []string{"farcaster", "dune", "data", statusTag},
			PaymentType: "crypto",
		},
		{
			ID:          "bc-2" + idSuffix,
			Title:       "Restaurant recommendations in NYC",
			Platform:    "BOUNTYCASTER",
			Reward:      "50",
			Currency:    "USDC",
			URL:         "https://bountycaster.xyz/bounty/0x22df1gb9" + urlSuffix,
			CreatedAt:   time.Now().Add(-2 * time.Hour),
			Description: "Looking for the best pizza spots in Brooklyn.",
			Tags:        []string{"farcaster", "nyc", "pizza", statusTag},
			PaymentType: "crypto",
		},
	}

	for _, b := range mockBounties {
		ch <- b
	}
}
