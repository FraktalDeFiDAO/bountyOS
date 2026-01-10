package security

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

// RateLimiter tracks API rate limits and enforces waiting when limits are reached
type RateLimiter struct {
	mu                 sync.Mutex
	remaining          int
	resetTime          time.Time
	lastRequestTime    time.Time
	requestCount       int
	minRequestInterval time.Duration
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		remaining:          0,
		resetTime:          time.Now(),
		lastRequestTime:    time.Now(),
		requestCount:       0,
		minRequestInterval: 2 * time.Second, // Minimum 2 seconds between requests
	}
}

// UpdateFromHeaders updates rate limit information from HTTP response headers
func (rl *RateLimiter) UpdateFromHeaders(resp *http.Response) {
	if resp == nil {
		return
	}

	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Parse rate limit headers
	if remainingStr := resp.Header.Get("X-RateLimit-Remaining"); remainingStr != "" {
		if remaining, err := strconv.Atoi(remainingStr); err == nil {
			rl.remaining = remaining
		}
	}

	if resetStr := resp.Header.Get("X-RateLimit-Reset"); resetStr != "" {
		if reset, err := strconv.ParseInt(resetStr, 10, 64); err == nil {
			rl.resetTime = time.Unix(reset, 0)
		}
	}

	// Increment request count
	rl.requestCount++
	rl.lastRequestTime = time.Now()
}

// WaitIfNeeded checks if we need to wait due to rate limits
func (rl *RateLimiter) WaitIfNeeded() {
	if os.Getenv("BOUNTYOS_DISABLE_RATE_LIMIT_SLEEP") != "" {
		return
	}

	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Check if we're approaching rate limits
	if rl.remaining <= 5 {
		waitTime := time.Until(rl.resetTime)
		if waitTime > 0 {
			GetLogger().Debug("Approaching rate limit (%d remaining), waiting %v", rl.remaining, waitTime)
			time.Sleep(waitTime)
		}
	}

	// Enforce minimum request interval
	sinceLastRequest := time.Since(rl.lastRequestTime)
	if sinceLastRequest < rl.minRequestInterval {
		waitTime := rl.minRequestInterval - sinceLastRequest
		GetLogger().Debug("Enforcing minimum request interval, waiting %v", waitTime)
		time.Sleep(waitTime)
	}
}

// GetStatus returns current rate limit status
func (rl *RateLimiter) GetStatus() string {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	return fmt.Sprintf("Remaining: %d, Reset: %s, Requests: %d",
		rl.remaining, rl.resetTime.Format(time.RFC3339), rl.requestCount)
}

// GitHubRateLimiter is a specialized rate limiter for GitHub API
type GitHubRateLimiter struct {
	*RateLimiter
	token string
}

// NewGitHubRateLimiter creates a new GitHub-specific rate limiter
func NewGitHubRateLimiter(token string) *GitHubRateLimiter {
	return &GitHubRateLimiter{
		RateLimiter: NewRateLimiter(),
		token:       token,
	}
}

// CheckAndWait checks rate limits and waits if necessary
func (g *GitHubRateLimiter) CheckAndWait() {
	if g.token == "" {
		// Unauthenticated requests have lower rate limits
		g.minRequestInterval = 10 * time.Second
	} else {
		// Authenticated requests can go faster
		g.minRequestInterval = 2 * time.Second
	}
	g.WaitIfNeeded()
}
