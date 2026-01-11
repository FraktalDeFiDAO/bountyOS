package security

import (
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
	// Create custom TLS configuration
	tlsConfig := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_RSA_WITH_AES_128_CBC_SHA,
		},
	}

	// Add system CA certificates
	caCertPool, err := x509.SystemCertPool()
	if err != nil {
		log.Printf("Warning: Could not load system cert pool: %v, using default", err)
		caCertPool = x509.NewCertPool()
	}
	tlsConfig.RootCAs = caCertPool

	// Create custom transport with security settings
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
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
