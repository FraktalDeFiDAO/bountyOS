package security

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestMaskToken(t *testing.T) {
	tests := []struct {
		token    string
		expected string
	}{
		{"ghp_1234567890", "gh**********90"},
		{"abcd", "****"},
		{"123", "****"},
		{"", ""},
		{"secret-token-long", "se*************ng"},
	}

	for _, tt := range tests {
		t.Run(tt.token, func(t *testing.T) {
			got := MaskToken(tt.token)
			if got != tt.expected {
				t.Errorf("MaskToken(%s) = %s, want %s", tt.token, got, tt.expected)
			}
		})
	}
}

func TestSanitizeString(t *testing.T) {
	input := "Hello\nWorld\tTest\r"
	expected := "Hello World Test "
	got := SanitizeString(input)
	if got != expected {
		t.Errorf("SanitizeString() = %q, want %q", got, expected)
	}
}

func TestValidateCurrency(t *testing.T) {
	tests := []struct {
		currency string
		valid    bool
	}{
		{"USDC", true},
		{"sol", true},
		{"BTC/ETH", true},
		{"INVALID", false},
		{"", false},
	}

	for _, tt := range tests {
		if ValidateCurrency(tt.currency) != tt.valid {
			t.Errorf("ValidateCurrency(%s) = %v, want %v", tt.currency, !tt.valid, tt.valid)
		}
	}
}

func TestValidateURL(t *testing.T) {
	tests := []struct {
		urlStr string
		valid  bool
	}{
		{"https://github.com", true},
		{"http://example.com", true},
		{"ftp://files.com", false},
		{"javascript:alert(1)", false},
		{"http://localhost", false},
	}

	for _, tt := range tests {
		if ValidateURL(tt.urlStr) != tt.valid {
			t.Errorf("ValidateURL(%s) = %v, want %v", tt.urlStr, !tt.valid, tt.valid)
		}
	}
}

func TestNormalizeURL(t *testing.T) {
	input := " https://example.com/path).\n"
	expected := "https://example.com/path"
	got := NormalizeURL(input)
	if got != expected {
		t.Errorf("NormalizeURL() = %q, want %q", got, expected)
	}
}

func TestValidateURLReachable(t *testing.T) {
	t.Setenv("BOUNTYOS_ALLOW_LOCAL_URLS", "true")

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/ok":
			w.WriteHeader(http.StatusOK)
		case "/forbidden":
			w.WriteHeader(http.StatusForbidden)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	ctx := context.Background()
	if !ValidateURLReachable(ctx, ts.URL+"/ok", 2*time.Second) {
		t.Errorf("expected /ok to be reachable")
	}
	if !ValidateURLReachable(ctx, ts.URL+"/forbidden", 2*time.Second) {
		t.Errorf("expected /forbidden to be treated as reachable")
	}
	if ValidateURLReachable(ctx, ts.URL+"/missing", 2*time.Second) {
		t.Errorf("expected /missing to be unreachable")
	}
}

func TestSecureRequest(t *testing.T) {
	req, _ := http.NewRequest("GET", "https://api.github.com", nil)
	SecureRequest(req, "test-token")

	if req.Header.Get("Authorization") != "token test-token" {
		t.Errorf("Authorization header not set correctly")
	}
	if req.Header.Get("User-Agent") == "" {
		t.Errorf("User-Agent header not set")
	}
}

func TestRateLimiter(t *testing.T) {
	rl := NewRateLimiter()
	rl.remaining = 10

	resp := &http.Response{
		Header: make(http.Header),
	}
	resp.Header.Set("X-RateLimit-Remaining", "4")
	resp.Header.Set("X-RateLimit-Reset", "1609459200") // 2021-01-01

	rl.UpdateFromHeaders(resp)
	if rl.remaining != 4 {
		t.Errorf("Expected remaining 4, got %d", rl.remaining)
	}

	status := rl.GetStatus()
	if !strings.Contains(status, "Remaining: 4") {
		t.Errorf("Status should contain remaining count: %s", status)
	}
}

func TestGitHubResponseValidation(t *testing.T) {
	jsonData := `{
		"items": [
			{
				"title": "Test Issue",
				"html_url": "https://github.com/test/test/issues/1",
				"created_at": "2023-01-01T00:00:00Z",
				"body": "Test body"
			}
		]
	}`

	resp, err := ValidateGitHubResponse([]byte(jsonData))
	if err != nil {
		t.Fatalf("Validation failed: %v", err)
	}
	if len(resp.Items) != 1 {
		t.Errorf("Expected 1 item, got %d", len(resp.Items))
	}

	// Test XSS detection
	xssData := `{
		"items": [
			{
				"title": "<script>alert(1)</script>",
				"html_url": "https://github.com/test/test/issues/2",
				"created_at": "2023-01-01T00:00:00Z",
				"body": "Test body"
			}
		]
	}`
	resp, err = ValidateGitHubResponse([]byte(xssData))
	if err != nil {
		t.Fatalf("Validation should not fail on XSS content: %v", err)
	}
	if len(resp.Items) != 0 {
		t.Errorf("Expected invalid XSS items to be dropped, got %d items", len(resp.Items))
	}
}
