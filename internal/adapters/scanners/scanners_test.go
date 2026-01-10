package scanners

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"bountyos-v8/internal/core"
)

func TestGitHubScanner_Scan(t *testing.T) {
	t.Setenv("BOUNTYOS_DISABLE_RATE_LIMIT_SLEEP", "1")

	// Mock GitHub API
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify Request
		if !strings.Contains(r.URL.Path, "/search/issues") {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// Return Mock JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{
			"items": [
				{
					"title": "Urgent Security Fix Needed",
					"html_url": "https://github.com/test/repo/issues/1",
					"created_at": "`+time.Now().Format(time.RFC3339)+`",
					"body": "We need a fix for a security vulnerability.",
					"labels": [{"name": "bug"}, {"name": "urgent"}, {"name": "100 USDC"}]
				},
				{
					"title": "Old Issue",
					"html_url": "https://github.com/test/repo/issues/2",
					"created_at": "2020-01-01T00:00:00Z",
					"body": "Old stuff",
					"labels": []
				}
			]
		}`)
	}))
	defer ts.Close()

	// Initialize Scanner with Mock BaseURL
	scanner := NewGitHubScanner("dummy-token", GitHubScannerConfig{})
	scanner.baseURL = ts.URL
	scanner.endpoints = []string{"test-label"} // Reduce to 1 endpoint to speed up test

	// Run Scan
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ch, err := scanner.Scan(ctx)
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	// Collect Results
	var bounties []core.Bounty
	for b := range ch {
		bounties = append(bounties, b)
	}

	// Verify
	if len(bounties) != 2 {
		t.Fatalf("Expected 2 bounties, got %d", len(bounties))
	}

	var target *core.Bounty
	for i := range bounties {
		if bounties[i].Title == "Urgent Security Fix Needed" {
			target = &bounties[i]
			break
		}
	}

	if target == nil {
		t.Fatalf("Expected to find 'Urgent Security Fix Needed' bounty")
	}
	if target.Title != "Urgent Security Fix Needed" {
		t.Errorf("Wrong title: %s", target.Title)
	}
	if target.Reward != "100 USDC" {
		t.Errorf("Wrong reward: %s", target.Reward)
	}
	if target.PaymentType != "crypto" {
		t.Errorf("Wrong payment type: %s", target.PaymentType)
	}
}

func TestGitHubScanner_Paginates(t *testing.T) {
	t.Setenv("BOUNTYOS_DISABLE_RATE_LIMIT_SLEEP", "1")

	type label struct {
		Name string `json:"name"`
	}
	type item struct {
		Title     string  `json:"title"`
		HTMLURL   string  `json:"html_url"`
		CreatedAt string  `json:"created_at"`
		Body      string  `json:"body"`
		Labels    []label `json:"labels"`
	}
	type response struct {
		Items []item `json:"items"`
	}

	now := time.Now().UTC().Format(time.RFC3339)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "/search/issues") {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
		if perPage == 0 {
			perPage = 30
		}

		count := 0
		switch page {
		case 1:
			count = perPage
		case 2:
			count = 1
		default:
			count = 0
		}

		items := make([]item, 0, count)
		for i := 0; i < count; i++ {
			items = append(items, item{
				Title:     fmt.Sprintf("Issue p%d-%d", page, i),
				HTMLURL:   fmt.Sprintf("https://github.com/test/repo/issues/%d", page*100+i),
				CreatedAt: now,
				Body:      "Body",
				Labels:    []label{{Name: "bug"}},
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response{Items: items})
	}))
	defer ts.Close()

	cfg := GitHubScannerConfig{PerPage: 5, MaxPages: 2}
	scanner := NewGitHubScanner("dummy-token", cfg)
	scanner.baseURL = ts.URL
	scanner.endpoints = []string{"test-label"}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	ch, err := scanner.Scan(ctx)
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	count := 0
	foundPageTwo := false
	for b := range ch {
		count++
		if strings.Contains(b.Title, "Issue p2-") {
			foundPageTwo = true
		}
	}

	expected := cfg.PerPage + 1
	if count != expected {
		t.Fatalf("Expected %d bounties, got %d", expected, count)
	}
	if !foundPageTwo {
		t.Fatalf("Expected to find a page 2 bounty")
	}
}
