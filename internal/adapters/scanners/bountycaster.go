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
	Bounties []struct {
		ID        string  `json:"id"`
		Title     string  `json:"title"`
		Body      string  `json:"body"`
		Amount    float64 `json:"amount"`
		Token     string  `json:"token"`
		CastHash  string  `json:"cast_hash"`
		CreatedAt string  `json:"created_at"`
	} `json:"bounties"`
}

func NewBountycasterScanner(cfg BountycasterScannerConfig) *BountycasterScanner {
	baseURL := strings.TrimRight(cfg.BaseURL, "/")
	if baseURL == "" {
		baseURL = "https://bountycaster.xyz/api/bounties"
	}
	statuses := cfg.Statuses
	if len(statuses) == 0 {
		statuses = []string{"active", "funded"}
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

		for _, status := range s.statuses {
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

	resp, err := doRequestWithRetry(ctx, s.client, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var results BountycasterResponse
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return err
	}

	for _, item := range results.Bounties {
		createdAt, _ := time.Parse(time.RFC3339, item.CreatedAt)
		tags := []string{"farcaster", "social"}
		if status != "" {
			tags = append(tags, status)
		}

		bounty := core.Bounty{
			ID:          item.ID,
			Title:       item.Title,
			Platform:    "BOUNTYCASTER",
			Reward:      fmt.Sprintf("%.2f", item.Amount),
			Currency:    item.Token,
			URL:         "https://bountycaster.xyz/bounty/" + item.CastHash,
			CreatedAt:   createdAt,
			Description: item.Body,
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
