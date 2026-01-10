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
	now := time.Now().UTC().Format(time.RFC3339)
	activeResponse := fmt.Sprintf(`[{
		"id":"st-active",
		"title":"Active Listing",
		"description":"active",
		"reward":500,
		"token":"USDC",
		"deadline":"",
		"createdAt":"%s",
		"slug":"active-listing"
	}]`, now)
	fundedResponse := fmt.Sprintf(`[{
		"id":"st-funded",
		"title":"Funded Listing",
		"description":"funded",
		"reward":750,
		"token":"USDC",
		"deadline":"",
		"createdAt":"%s",
		"slug":"funded-listing"
	}]`, now)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/bounties" {
			http.NotFound(w, r)
			return
		}

		status := r.URL.Query().Get("status")
		w.Header().Set("Content-Type", "application/json")
		switch status {
		case "active":
			fmt.Fprint(w, activeResponse)
		case "funded":
			fmt.Fprint(w, fundedResponse)
		default:
			fmt.Fprint(w, `[]`)
		}
	}))
	defer ts.Close()

	scanner := NewSuperteamScanner(SuperteamScannerConfig{})
	scanner.baseURL = ts.URL + "/api/bounties"
	scanner.statuses = []string{"active", "funded"}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ch, err := scanner.Scan(ctx)
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	var bounties []core.Bounty
	for b := range ch {
		bounties = append(bounties, b)
	}

	if len(bounties) != 2 {
		t.Fatalf("Expected 2 bounties, got %d", len(bounties))
	}

	for _, b := range bounties {
		switch b.Title {
		case "Active Listing":
			if !hasTag(b.Tags, "active") {
				t.Errorf("Expected active tag in %v", b.Tags)
			}
		case "Funded Listing":
			if !hasTag(b.Tags, "funded") {
				t.Errorf("Expected funded tag in %v", b.Tags)
			}
		default:
			t.Fatalf("Unexpected bounty title: %s", b.Title)
		}
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
