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
	openResponse := fmt.Sprintf(`{
		"bounties": [{
			"uid":"bc-open",
			"title":"Open Bounty",
			"summary_text":"open",
			"created_at":"%s",
			"expiration_date":"%s",
			"tag_slugs":["tag-open"],
			"links":{"resource":"/bounty/0xopen","external":"https://farcaster.xyz/~/conversations/0xopen"},
			"reward_summary":{"unit_amount":"150","symbol":"USDC"},
			"platform":{"type":"farcaster","hash":"0xopen"}
		}]
	}`, now, now)
	inProgressResponse := fmt.Sprintf(`{
		"bounties": [{
			"uid":"bc-inprogress",
			"title":"In Progress Bounty",
			"summary_text":"in-progress",
			"created_at":"%s",
			"expiration_date":"%s",
			"tag_slugs":["tag-in-progress"],
			"links":{"resource":"/bounty/0xinprogress","external":"https://farcaster.xyz/~/conversations/0xinprogress"},
			"reward_summary":{"unit_amount":"300","symbol":"USDC"},
			"platform":{"type":"farcaster","hash":"0xinprogress"}
		}]
	}`, now, now)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/api/v1/bounties/open":
			fmt.Fprint(w, openResponse)
		case "/api/v1/bounties/in-progress":
			fmt.Fprint(w, inProgressResponse)
		default:
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	scanner := NewBountycasterScanner(BountycasterScannerConfig{})
	scanner.baseURL = ts.URL + "/api/v1/bounties"
	scanner.statuses = []string{"open", "in-progress"}

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
		case "Open Bounty":
			if !hasTag(b.Tags, "open") {
				t.Errorf("Expected open tag in %v", b.Tags)
			}
		case "In Progress Bounty":
			if !hasTag(b.Tags, "in-progress") {
				t.Errorf("Expected in-progress tag in %v", b.Tags)
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
