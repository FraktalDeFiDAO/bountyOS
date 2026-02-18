package security

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	stderrors "errors"

	"bountyos-v8/internal/errors"
	"bountyos-v8/internal/resilience"
)

// SecureHTTPClient creates a secure HTTP client with proper TLS configuration
// and timeout settings for secure API communications
func SecureHTTPClient() *http.Client {
	// Create TLS configuration. Avoid restricting cipher suites to prevent
	// handshake timeouts with providers that prefer newer defaults.
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	// Add system CA certificates
	caCertPool, err := x509.SystemCertPool()
	if err != nil {
		log.Printf("Warning: Could not load system cert pool: %v, using default", err)
		caCertPool = x509.NewCertPool()
	}
	tlsConfig.RootCAs = caCertPool

	// Create custom transport with security settings
	preferIPv4 := true
	if strings.EqualFold(os.Getenv("BOUNTYOS_PREFER_IPV4"), "false") {
		preferIPv4 = false
	}

	dialer := &net.Dialer{
		Timeout:   15 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: func(ctx context.Context, network, address string) (net.Conn, error) {
			host, port, err := net.SplitHostPort(address)
			if err != nil {
				return dialer.DialContext(ctx, network, address)
			}

			ipv4s := []net.IP{}
			ipv6s := []net.IP{}
			if preferIPv4 {
				if ips, err := net.DefaultResolver.LookupIP(ctx, "ip4", host); err == nil {
					ipv4s = append(ipv4s, ips...)
				}
				if len(ipv4s) == 0 {
					if ips, err := net.DefaultResolver.LookupIP(ctx, "ip6", host); err == nil {
						ipv6s = append(ipv6s, ips...)
					}
				}
			} else {
				if ips, err := net.DefaultResolver.LookupIP(ctx, "ip6", host); err == nil {
					ipv6s = append(ipv6s, ips...)
				}
				if len(ipv6s) == 0 {
					if ips, err := net.DefaultResolver.LookupIP(ctx, "ip4", host); err == nil {
						ipv4s = append(ipv4s, ips...)
					}
				}
			}

			if len(ipv4s) == 0 && len(ipv6s) == 0 {
				return dialer.DialContext(ctx, network, address)
			}

			candidates := make([]net.IP, 0, len(ipv4s)+len(ipv6s))
			if preferIPv4 {
				candidates = append(candidates, ipv4s...)
				candidates = append(candidates, ipv6s...)
			} else {
				candidates = append(candidates, ipv6s...)
				candidates = append(candidates, ipv4s...)
			}

			var lastErr error
			for _, ip := range candidates {
				conn, err := dialer.DialContext(ctx, network, net.JoinHostPort(ip.String(), port))
				if err == nil {
					return conn, nil
				}
				lastErr = err
			}
			if lastErr != nil {
				return nil, lastErr
			}

			return dialer.DialContext(ctx, network, address)
		},
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   15 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig:       tlsConfig,
	}

	// Create HTTP client with secure defaults
	client := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}

	return client
}

// MaskToken masks sensitive tokens in strings to prevent accidental logging
func MaskToken(token string) string {
	if token == "" {
		return ""
	}

	// If token is short, mask completely
	if len(token) <= 4 {
		return "****"
	}

	// Show first 2 and last 2 characters, mask the rest
	maskedLength := len(token) - 4
	if maskedLength <= 0 {
		return "****"
	}

	return token[:2] + strings.Repeat("*", maskedLength) + token[len(token)-2:]
}

// SecureRequest adds security headers and handles sensitive data in HTTP requests
func SecureRequest(req *http.Request, token string) {
	if req == nil {
		return
	}

	// Add security headers
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "BountyOS-Secure/1.0")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	// Add authorization if token is provided
	if token != "" {
		req.Header.Set("Authorization", "token "+token)
		log.Printf("Added authorization header with masked token: %s", MaskToken(token))
	}

	// Add content security headers
	req.Header.Set("X-Content-Type-Options", "nosniff")
	req.Header.Set("X-Frame-Options", "DENY")
}

// GetEnvWithFallback gets environment variable with fallback to default
func GetEnvWithFallback(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// IsDebugMode checks if application is running in debug mode
func IsDebugMode() bool {
	return strings.ToLower(GetEnvWithFallback("DEBUG", "false")) == "true"
}

type RobustHTTPClient struct {
	client         *http.Client
	circuitBreaker *resilience.CircuitBreaker
	retryConfig    resilience.RetryConfig
	logger         func(format string, args ...any)
	debug          bool
	defaultTimeout time.Duration
}

type RobustClientOption func(*RobustHTTPClient)

func WithCircuitBreaker(failureThreshold, successThreshold int, timeout time.Duration) RobustClientOption {
	return func(c *RobustHTTPClient) {
		c.circuitBreaker = resilience.NewCircuitBreaker(failureThreshold, successThreshold, timeout)
	}
}

func WithRetry(maxRetries int, initialDelay, maxDelay time.Duration) RobustClientOption {
	return func(c *RobustHTTPClient) {
		c.retryConfig = resilience.RetryConfig{
			MaxRetries:   maxRetries,
			InitialDelay: initialDelay,
			MaxDelay:     maxDelay,
			RetryableError: func(err error) bool {
				return errors.IsRetryable(err) || isTransientHTTPError(err)
			},
		}
	}
}

func WithLogger(logger func(format string, args ...any)) RobustClientOption {
	return func(c *RobustHTTPClient) {
		c.logger = logger
	}
}

func WithDebug(enabled bool) RobustClientOption {
	return func(c *RobustHTTPClient) {
		c.debug = enabled
	}
}

func WithDefaultTimeout(timeout time.Duration) RobustClientOption {
	return func(c *RobustHTTPClient) {
		c.defaultTimeout = timeout
	}
}

func NewRobustHTTPClient(opts ...RobustClientOption) *RobustHTTPClient {
	c := &RobustHTTPClient{
		client:         SecureHTTPClient(),
		retryConfig:    resilience.DefaultRetryConfig(),
		logger:         log.Printf,
		debug:          IsDebugMode(),
		defaultTimeout: 30 * time.Second,
	}

	for _, opt := range opts {
		opt(c)
	}

	if c.circuitBreaker == nil {
		c.circuitBreaker = resilience.NewCircuitBreaker(5, 2, 30*time.Second)
	}

	return c
}

func (c *RobustHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return c.DoWithContext(context.Background(), req)
}

func (c *RobustHTTPClient) DoWithContext(ctx context.Context, req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var lastErr error

	operation := req.Method + " " + req.URL.Path

	err := c.circuitBreaker.Execute(ctx, func() error {
		err := resilience.Retry(ctx, c.retryConfig, func() error {
			if c.debug {
				c.logRequest(req)
			}

			resp, lastErr = c.client.Do(req.WithContext(ctx))
			if lastErr != nil {
				return classifyHTTPError(operation, lastErr)
			}

			if resp.StatusCode >= 500 {
				body, _ := io.ReadAll(resp.Body)
				resp.Body = io.NopCloser(bytes.NewReader(body))
				return errors.NewTransientError(operation, "server error: "+resp.Status, nil)
			}

			if resp.StatusCode == 429 {
				retryAfter := parseRetryAfter(resp.Header.Get("Retry-After"))
				return errors.NewRateLimitError(operation, "rate limited", retryAfter)
			}

			return nil
		})
		return err
	})

	if err != nil {
		return nil, err
	}

	if c.debug {
		c.logResponse(resp)
	}

	return resp, nil
}

func (c *RobustHTTPClient) Get(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.NewPermanentError("GET", "failed to create request", err)
	}
	return c.DoWithContext(ctx, req)
}

func (c *RobustHTTPClient) Post(ctx context.Context, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, errors.NewPermanentError("POST", "failed to create request", err)
	}
	return c.DoWithContext(ctx, req)
}

func (c *RobustHTTPClient) PostJSON(ctx context.Context, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, errors.NewPermanentError("POST", "failed to create request", err)
	}
	req.Header.Set("Content-Type", "application/json")
	return c.DoWithContext(ctx, req)
}

func (c *RobustHTTPClient) CircuitBreakerState() resilience.State {
	return c.circuitBreaker.State()
}

func (c *RobustHTTPClient) ResetCircuitBreaker() {
	c.circuitBreaker.Reset()
}

func (c *RobustHTTPClient) logRequest(req *http.Request) {
	c.logger("[HTTP] %s %s", req.Method, req.URL.String())
	for key, values := range req.Header {
		if strings.EqualFold(key, "authorization") {
			for _, v := range values {
				c.logger("[HTTP]   %s: %s", key, MaskToken(v))
			}
		} else {
			c.logger("[HTTP]   %s: %v", key, values)
		}
	}
}

func (c *RobustHTTPClient) logResponse(resp *http.Response) {
	c.logger("[HTTP] Response: %s", resp.Status)
	for key, values := range resp.Header {
		c.logger("[HTTP]   %s: %v", key, values)
	}
}

func classifyHTTPError(operation string, err error) error {
	if err == nil {
		return nil
	}

	if urlErr, ok := err.(*url.Error); ok {
		if urlErr.Timeout() {
			return errors.NewTimeoutError(operation, "request timeout", err)
		}
		if strings.Contains(urlErr.Error(), "connection refused") ||
			strings.Contains(urlErr.Error(), "connection reset") ||
			strings.Contains(urlErr.Error(), "broken pipe") {
			return errors.NewNetworkError(operation, "connection error", err)
		}
	}

	return errors.NewTransientError(operation, "request failed", err)
}

func isTransientHTTPError(err error) bool {
	var bountyErr *errors.BountyError
	if stderrors.As(err, &bountyErr) {
		return bountyErr.Type == errors.ErrorTypeTransient ||
			bountyErr.Type == errors.ErrorTypeNetwork ||
			bountyErr.Type == errors.ErrorTypeTimeout
	}
	return false
}

func parseRetryAfter(value string) time.Duration {
	if value == "" {
		return 60 * time.Second
	}

	if strings.Contains(value, ":") {
		return 60 * time.Second
	}

	var seconds int
	if _, err := fmt.Sscanf(value, "%d", &seconds); err == nil && seconds > 0 {
		return time.Duration(seconds) * time.Second
	}

	return 60 * time.Second
}
