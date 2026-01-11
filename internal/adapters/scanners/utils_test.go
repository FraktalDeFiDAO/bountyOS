package scanners

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"
)

type MockRoundTripper struct {
	RoundTripFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.RoundTripFunc(req)
}

func TestDoRequestWithRetry(t *testing.T) {
	// Reduce backoff for testing
	originalBackoff := baseBackoff
	baseBackoff = 1 * time.Millisecond
	defer func() { baseBackoff = originalBackoff }()

	// 1. Test success after retries
	t.Run("SuccessAfterRetry", func(t *testing.T) {
		attempts := 0
		client := &http.Client{
			Transport: &MockRoundTripper{
				RoundTripFunc: func(req *http.Request) (*http.Response, error) {
					attempts++
					if attempts < 3 {
						return nil, errors.New("network error")
					}
					return &http.Response{
						StatusCode: 200,
						Body:       http.NoBody,
					}, nil
				},
			},
		}

		ctx := context.Background()
		req, _ := http.NewRequest("GET", "http://example.com", nil)

		// Force short backoff for test
		// Since we can't easily change the constant in the package, we just accept the delay
		// or ideally we would have made the backoff configurable.
		// For now, we'll just run it. It will wait 1s + 2s = 3s.
		// To speed it up, we could use a shorter timeout context, but that might cancel it.
		// A better way is to make baseBackoff a variable.

		_, err := doRequestWithRetry(ctx, client, req)
		if err != nil {
			t.Errorf("Expected success, got error: %v", err)
		}
		if attempts != 3 {
			t.Errorf("Expected 3 attempts, got %d", attempts)
		}
	})

	// 2. Test failure after max retries
	t.Run("FailAfterMaxRetries", func(t *testing.T) {
		client := &http.Client{
			Transport: &MockRoundTripper{
				RoundTripFunc: func(req *http.Request) (*http.Response, error) {
					return nil, errors.New("persistent error")
				},
			},
		}

		ctx := context.Background()
		req, _ := http.NewRequest("GET", "http://example.com", nil)

		_, err := doRequestWithRetry(ctx, client, req)
		if err == nil {
			t.Error("Expected error, got success")
		}
	})

	// 3. Test 5xx error retry
	t.Run("RetryOn500", func(t *testing.T) {
		attempts := 0
		client := &http.Client{
			Transport: &MockRoundTripper{
				RoundTripFunc: func(req *http.Request) (*http.Response, error) {
					attempts++
					if attempts == 1 {
						return &http.Response{StatusCode: 500, Body: http.NoBody}, nil
					}
					return &http.Response{StatusCode: 200, Body: http.NoBody}, nil
				},
			},
		}

		ctx := context.Background()
		req, _ := http.NewRequest("GET", "http://example.com", nil)

		resp, err := doRequestWithRetry(ctx, client, req)
		if err != nil {
			t.Errorf("Expected success, got error: %v", err)
		}
		if resp.StatusCode != 200 {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}
		if attempts != 2 {
			t.Errorf("Expected 2 attempts, got %d", attempts)
		}
	})
}
