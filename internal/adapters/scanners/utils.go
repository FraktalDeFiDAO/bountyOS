package scanners

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"sync"
	"time"

	"bountyos-v8/internal/resilience"
	"bountyos-v8/internal/security"
)

const (
	maxRetries       = 3
	defaultTimeout   = 30 * time.Second
	retryAfterHeader = "Retry-After"
)

var (
	baseBackoff      = 1 * time.Second
	registry         = resilience.NewCircuitBreakerRegistry()
	scannerConfigs   = make(map[string]*ScannerConfig)
	scannerConfigsMu sync.RWMutex
)

type RequestMetrics struct {
	URL          string
	Duration     time.Duration
	Retries      int
	StatusCode   int
	Success      bool
	RetryAfter   time.Duration
	ErrorMessage string
}

type ScannerConfig struct {
	Name           string
	Timeout        time.Duration
	MaxRetries     int
	CircuitBreaker *resilience.CircuitBreaker
}

func GetScannerConfig(name string) *ScannerConfig {
	scannerConfigsMu.RLock()
	cfg, exists := scannerConfigs[name]
	scannerConfigsMu.RUnlock()

	if exists {
		return cfg
	}

	scannerConfigsMu.Lock()
	defer scannerConfigsMu.Unlock()

	if cfg, exists = scannerConfigs[name]; exists {
		return cfg
	}

	breaker := registry.Get(name, 5, 2, 30*time.Second)
	cfg = &ScannerConfig{
		Name:           name,
		Timeout:        defaultTimeout,
		MaxRetries:     maxRetries,
		CircuitBreaker: breaker,
	}
	scannerConfigs[name] = cfg
	return cfg
}

func SetScannerTimeout(name string, timeout time.Duration) {
	cfg := GetScannerConfig(name)
	scannerConfigsMu.Lock()
	cfg.Timeout = timeout
	scannerConfigsMu.Unlock()
}

type contextKey string

const scannerNameKey contextKey = "scannerName"

func ContextWithScannerName(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, scannerNameKey, name)
}

func ScannerNameFromContext(ctx context.Context) string {
	if name, ok := ctx.Value(scannerNameKey).(string); ok {
		return name
	}
	return "unknown"
}

func doRequestWithRetry(ctx context.Context, client *http.Client, req *http.Request) (*http.Response, error) {
	scannerName := ScannerNameFromContext(ctx)
	cfg := GetScannerConfig(scannerName)

	var lastErr error
	var metrics RequestMetrics
	metrics.URL = req.URL.String()
	startTime := time.Now()

	if client.Timeout == 0 {
		client.Timeout = cfg.Timeout
	}

	var resp *http.Response
	err := cfg.CircuitBreaker.Execute(ctx, func() error {
		for i := 0; i <= cfg.MaxRetries; i++ {
			if i > 0 {
				backoff := time.Duration(math.Pow(2, float64(i-1))) * baseBackoff
				retryAfter := parseRetryAfterFromContext(ctx)
				if retryAfter > 0 && retryAfter > backoff {
					backoff = retryAfter
				}

				select {
				case <-time.After(backoff):
				case <-ctx.Done():
					lastErr = ctx.Err()
					return ctx.Err()
				}
				metrics.Retries++
				security.GetLogger().Info("Retrying request to %s (attempt %d/%d)...", req.URL.String(), i, cfg.MaxRetries)
			}

			var reqErr error
			resp, reqErr = client.Do(req)
			if reqErr != nil {
				lastErr = reqErr
				continue
			}

			metrics.StatusCode = resp.StatusCode

			if resp.StatusCode >= 500 || resp.StatusCode == 429 {
				retryAfter := parseRetryAfterHeader(resp)
				if retryAfter > 0 {
					metrics.RetryAfter = retryAfter
					security.GetLogger().Warn("Rate limited on %s, Retry-After: %v", req.URL.String(), retryAfter)
				}
				resp.Body.Close()
				lastErr = fmt.Errorf("server returned status %d", resp.StatusCode)
				continue
			}

			metrics.Duration = time.Since(startTime)
			metrics.Success = true
			logRequestMetrics(metrics)
			return nil
		}

		metrics.Duration = time.Since(startTime)
		metrics.ErrorMessage = lastErr.Error()
		logRequestMetrics(metrics)
		return fmt.Errorf("after %d retries, last error: %w", cfg.MaxRetries, lastErr)
	})

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func parseRetryAfterHeader(resp *http.Response) time.Duration {
	retryAfter := resp.Header.Get(retryAfterHeader)
	if retryAfter == "" {
		return 0
	}

	if seconds, err := strconv.Atoi(retryAfter); err == nil {
		return time.Duration(seconds) * time.Second
	}

	if t, err := time.Parse(time.RFC1123, retryAfter); err == nil {
		return time.Until(t)
	}

	return 0
}

func parseRetryAfterFromContext(ctx context.Context) time.Duration {
	if v := ctx.Value(contextKey("retryAfter")); v != nil {
		if d, ok := v.(time.Duration); ok {
			return d
		}
	}
	return 0
}

func logRequestMetrics(m RequestMetrics) {
	status := "success"
	if !m.Success {
		status = "failed"
	}
	security.GetLogger().Debug("Request %s: %s (duration=%v, retries=%d, status=%d)",
		m.URL, status, m.Duration, m.Retries, m.StatusCode)
}

func GetCircuitBreakerStats() map[string]struct {
	State     resilience.State
	Failures  int
	Successes int
} {
	return registry.AllStats()
}

func ResetAllCircuitBreakers() {
	registry.ResetAll()
}
