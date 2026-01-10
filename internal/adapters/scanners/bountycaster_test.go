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

func TestBountycasterScanner_ScanStatuses(t *testing.T) {
	now := time.Now().UTC().Format(time.RFC3339)
	activeResponse := fmt.Sprintf(`{
		"bounties": [{
			"id":"bc-active",
			"title":"Active Cast",
			"body":"active",
			"amount":150,
			"token":"USDC",
			"cast_hash":"0xactive",
			"created_at":"%s"
		}]
	}`, now)
	fundedResponse := fmt.Sprintf(`{
		"bounties": [{
			"id":"bc-funded",
			"title":"Funded Cast",
			"body":"funded",
			"amount":300,
			"token":"USDC",
			"cast_hash":"0xfunded",
			"created_at":"%s"
		}]
	}`, now)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/api/bounties/active":
			fmt.Fprint(w, activeResponse)
		case "/api/bounties/funded":
			fmt.Fprint(w, fundedResponse)
		default:
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	scanner := NewBountycasterScanner(BountycasterScannerConfig{})
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
		case "Active Cast":
			if !hasTag(b.Tags, "active") {
				t.Errorf("Expected active tag in %v", b.Tags)
			}
		case "Funded Cast":
			if !hasTag(b.Tags, "funded") {
				t.Errorf("Expected funded tag in %v", b.Tags)
			}
		default:
			t.Fatalf("Unexpected bounty title: %s", b.Title)
		}
		if b.Platform != "BOUNTYCASTER" {
			t.Errorf("Unexpected platform: %s", b.Platform)
		}
		if b.PaymentType != "crypto" {
			t.Errorf("Unexpected payment type: %s", b.PaymentType)
		}
	}
}
