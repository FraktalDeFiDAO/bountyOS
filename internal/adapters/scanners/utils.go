package scanners

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"time"

	"bountyos-v8/internal/security"
)

const (
	maxRetries = 3
)

var baseBackoff = 1 * time.Second

// doRequestWithRetry executes an HTTP request with exponential backoff retries.
// It returns the response or the last error encountered.
func doRequestWithRetry(ctx context.Context, client *http.Client, req *http.Request) (*http.Response, error) {
	var lastErr error

	for i := 0; i <= maxRetries; i++ {
		// If this is a retry, wait before sending
		if i > 0 {
			backoff := time.Duration(math.Pow(2, float64(i-1))) * baseBackoff
			select {
			case <-time.After(backoff):
			case <-ctx.Done():
				return nil, ctx.Err()
			}
			security.GetLogger().Info("Retrying request to %s (attempt %d/%d)...", req.URL.String(), i, maxRetries)
		}

		resp, err := client.Do(req)
		if err == nil {
			// Check for 5xx or 429 status codes to retry
			if resp.StatusCode >= 500 || resp.StatusCode == 429 {
				resp.Body.Close()
				lastErr = fmt.Errorf("server returned status %d", resp.StatusCode)
				continue
			}
			return resp, nil
		}
		
		lastErr = err
	}

	return nil, fmt.Errorf("after %d retries, last error: %w", maxRetries, lastErr)
}
