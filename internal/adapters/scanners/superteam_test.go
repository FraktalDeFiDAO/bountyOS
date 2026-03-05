package scanners

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"bountyos-v8/internal/core"
)

func TestSuperteamScanner_ScanStatuses(t *testing.T) {
	openResponse := `[{
		"id":"st-open",
		"title":"Open Listing",
		"type":"bounty",
		"rewardAmount":500,
		"token":"USDC",
		"deadline":"",
		"slug":"open-listing"
	}]`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/bounties" {
			http.NotFound(w, r)
			return
		}

		status := r.URL.Query().Get("status")
		w.Header().Set("Content-Type", "application/json")
		if status == "open" || status == "" {
			fmt.Fprint(w, openResponse)
		} else {
			fmt.Fprint(w, `[]`)
		}
	}))
	defer ts.Close()

	scanner := NewSuperteamScanner(SuperteamScannerConfig{})
	scanner.baseURL = ts.URL + "/api/bounties"
	scanner.statuses = []string{"active", "funded"}

	ctx := ContextWithScannerName(context.Background(), "SUPERTEAM")
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	ch, err := scanner.Scan(ctx)
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	var bounties []core.Bounty
	for b := range ch {
		bounties = append(bounties, b)
	}

	if len(bounties) < 1 {
		t.Fatalf("Expected at least 1 bounty, got %d", len(bounties))
	}

	for _, b := range bounties {
		if b.Platform != "SUPERTEAM" {
			t.Errorf("Unexpected platform: %s", b.Platform)
		}
		if b.PaymentType != "crypto" {
			t.Errorf("Unexpected payment type: %s", b.PaymentType)
		}
	}
}

func hasTag(tags []string, target string) bool {
	for _, tag := range tags {
		if tag == target {
			return true
		}
	}
	return false
}
