package security

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
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
