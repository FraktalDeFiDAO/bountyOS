package security

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

// GitHubAPIResponse represents the expected structure from GitHub API
type GitHubAPIResponse struct {
	Items []GitHubIssue `json:"items"`
}

// GitHubIssue represents a single GitHub issue from the search API
type GitHubIssue struct {
	Title     string `json:"title"`
	HTMLURL   string `json:"html_url"`
	CreatedAt string `json:"created_at"`
	Body      string `json:"body"`
	Labels    []struct {
		Name string `json:"name"`
	} `json:"labels"`
}

// ValidateGitHubResponse validates the structure and content of GitHub API responses
func ValidateGitHubResponse(data []byte) (*GitHubAPIResponse, error) {
	if len(data) == 0 {
		return nil, errors.New("empty response body")
	}

	var response GitHubAPIResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("invalid JSON format: %w", err)
	}

	// Validate each item
	for i, item := range response.Items {
		if err := validateGitHubIssue(item); err != nil {
			return nil, fmt.Errorf("invalid item at index %d: %w", i, err)
		}
	}

	return &response, nil
}

// ValidateGitHubResponseFromReader validates GitHub API response from io.Reader
func ValidateGitHubResponseFromReader(reader interface {
	Read(p []byte) (n int, err error)
}) (*GitHubAPIResponse, error) {
	// Read the response body
	var bodyBytes []byte
	buf := make([]byte, 4096)
	for {
		n, err := reader.Read(buf)
		if n > 0 {
			bodyBytes = append(bodyBytes, buf[:n]...)
		}
		if err != nil {
			break
		}
	}

	return ValidateGitHubResponse(bodyBytes)
}

// validateGitHubIssue validates a single GitHub issue
func validateGitHubIssue(issue GitHubIssue) error {
	// Validate required fields
	if strings.TrimSpace(issue.Title) == "" {
		return errors.New("title cannot be empty")
	}

	if strings.TrimSpace(issue.HTMLURL) == "" {
		return errors.New("html_url cannot be empty")
	}

	// Validate URL format
	if _, err := url.ParseRequestURI(issue.HTMLURL); err != nil {
		return fmt.Errorf("invalid html_url format: %w", err)
	}

	// Validate created_at format
	if strings.TrimSpace(issue.CreatedAt) == "" {
		return errors.New("created_at cannot be empty")
	}

	if _, err := time.Parse(time.RFC3339, issue.CreatedAt); err != nil {
		return fmt.Errorf("invalid created_at format (expected RFC3339): %w", err)
	}

	// Validate title length and content
	if len(issue.Title) > 500 {
		return errors.New("title too long (max 500 characters)")
	}

	// Basic XSS protection - check for script tags
	if containsScriptTags(issue.Title) || containsScriptTags(issue.Body) {
		return errors.New("potential XSS content detected")
	}

	return nil
}

// containsScriptTags checks if a string contains potential script tags
func containsScriptTags(content string) bool {
	// Case-insensitive regex for script tags
	re := regexp.MustCompile(`(?i)<script[^>]*>|</script>|javascript:|onerror\s*=|onclick\s*=|onload\s*=`)
	return re.MatchString(content)
}

// SanitizeString sanitizes strings for safe logging and display
func SanitizeString(input string) string {
	if input == "" {
		return ""
	}

	// Remove or escape potentially dangerous content
	sanitized := strings.ReplaceAll(input, "\n", " ")
	sanitized = strings.ReplaceAll(sanitized, "\r", " ")
	sanitized = strings.ReplaceAll(sanitized, "\t", " ")

	// Truncate long strings
	if len(sanitized) > 1000 {
		sanitized = sanitized[:1000] + "..."
	}

	return sanitized
}

// NormalizeURL trims whitespace and strips trailing punctuation from URLs.
func NormalizeURL(urlStr string) string {
	trimmed := strings.TrimSpace(urlStr)
	if trimmed == "" {
		return ""
	}

	trimmed = strings.ReplaceAll(trimmed, "\n", "")
	trimmed = strings.ReplaceAll(trimmed, "\r", "")
	trimmed = strings.ReplaceAll(trimmed, "\t", "")

	fields := strings.Fields(trimmed)
	if len(fields) > 0 {
		trimmed = fields[0]
	}

	trimmed = strings.TrimRight(trimmed, ".,;!?)\"'")
	return trimmed
}

// ValidateCurrency validates currency strings
func ValidateCurrency(currency string) bool {
	if currency == "" {
		return false
	}

	// List of valid currencies for the bounty system
	validCurrencies := []string{
		"USDC", "USDT", "SOL", "ETH", "BTC", "MATIC", "AVAX", "ARB", "OP",
		"CASHAPP", "VENMO", "PAYPAL", "STRIPE", "WISE", "USD", "EUR", "GBP",
	}

	upperCurrency := strings.ToUpper(strings.TrimSpace(currency))
	for _, valid := range validCurrencies {
		if upperCurrency == valid {
			return true
		}
	}

	// Allow combined currencies like "USDC/ETH"
	if strings.Contains(upperCurrency, "/") {
		parts := strings.Split(upperCurrency, "/")
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part == "" || !ValidateCurrency(part) {
				return false
			}
		}
		return true
	}

	return false
}

// ValidateURL validates URL format and safety
func ValidateURL(urlStr string) bool {
	if urlStr == "" {
		return false
	}

	parsedURL, err := url.ParseRequestURI(urlStr)
	if err != nil {
		return false
	}

	// Check for allowed schemes
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return false
	}

	// Check for potentially dangerous domains
	if strings.EqualFold(os.Getenv("BOUNTYOS_ALLOW_LOCAL_URLS"), "true") {
		return true
	}

	dangerousDomains := []string{
		"localhost", "127.0.0.1", "0.0.0.0",
		"file://", "ftp://", "javascript:",
	}

	host := strings.ToLower(parsedURL.Host)
	for _, dangerous := range dangerousDomains {
		if strings.Contains(host, dangerous) {
			return false
		}
	}

	return true
}

// ValidateURLReachable checks if a URL responds with an acceptable HTTP status.
func ValidateURLReachable(ctx context.Context, urlStr string, timeout time.Duration) bool {
	if !ValidateURL(urlStr) {
		return false
	}
	if timeout <= 0 {
		timeout = 5 * time.Second
	}

	client := SecureHTTPClient()
	client.Timeout = timeout

	statusOK := func(code int) bool {
		if code >= 200 && code < 400 {
			return true
		}
		switch code {
		case http.StatusUnauthorized, http.StatusForbidden, http.StatusMethodNotAllowed, http.StatusTooManyRequests:
			return true
		default:
			return false
		}
	}

	req, err := http.NewRequestWithContext(ctx, "HEAD", urlStr, nil)
	if err == nil {
		resp, err := client.Do(req)
		if err == nil {
			resp.Body.Close()
			if statusOK(resp.StatusCode) {
				return true
			}
			if resp.StatusCode != http.StatusMethodNotAllowed {
				return false
			}
		}
	}

	req, err = http.NewRequestWithContext(ctx, "GET", urlStr, nil)
	if err != nil {
		return false
	}
	req.Header.Set("Range", "bytes=0-0")
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	resp.Body.Close()
	return statusOK(resp.StatusCode) || resp.StatusCode == http.StatusRequestedRangeNotSatisfiable
}
