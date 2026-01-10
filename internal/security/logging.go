package security

import (
	"encoding/json"
	"fmt"

	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

// SecureLogger provides secure logging with sanitization and token masking
type SecureLogger struct {
	logger     *log.Logger
	maskTokens map[string]string // original -> masked
	mu         sync.Mutex
}

// NewSecureLogger creates a new secure logger
func NewSecureLogger() *SecureLogger {
	return &SecureLogger{
		logger:     log.New(os.Stdout, "[BOUNTYOS] ", log.Ldate|log.Ltime|log.Lshortfile),
		maskTokens: make(map[string]string),
	}
}

// SetOutput swaps the logger output destination (e.g., file, stderr).
func (sl *SecureLogger) SetOutput(w io.Writer) {
	sl.mu.Lock()
	defer sl.mu.Unlock()
	sl.logger.SetOutput(w)
}

// RegisterToken registers a token for automatic masking in logs
func (sl *SecureLogger) RegisterToken(token string) {
	if token == "" {
		return
	}

	sl.mu.Lock()
	defer sl.mu.Unlock()

	masked := MaskToken(token)
	sl.maskTokens[token] = masked
}

// Info logs an informational message with sanitization
func (sl *SecureLogger) Info(format string, v ...interface{}) {
	sl.log("INFO", format, v...)
}

// Warn logs a warning message with sanitization
func (sl *SecureLogger) Warn(format string, v ...interface{}) {
	sl.log("WARN", format, v...)
}

// Error logs an error message with sanitization
func (sl *SecureLogger) Error(format string, v ...interface{}) {
	sl.log("ERROR", format, v...)
}

// Debug logs a debug message (only in debug mode)
func (sl *SecureLogger) Debug(format string, v ...interface{}) {
	if !IsDebugMode() {
		return
	}
	sl.log("DEBUG", format, v...)
}

// Audit logs a security audit event in structured JSON format
func (sl *SecureLogger) Audit(actorID, action, resourceType, resourceID string, metadata map[string]interface{}) {
	event := map[string]interface{}{
		"timestamp":  time.Now().UTC().Format(time.RFC3339),
		"event_type": "AUDIT",
		"actor": map[string]string{
			"id": actorID,
		},
		"action": action,
		"resource": map[string]string{
			"type": resourceType,
			"id":   resourceID,
		},
		"status":   "success", // Default to success, metadata can override
		"metadata": metadata,
	}

	// Sanitize metadata
	sanitizedMeta := make(map[string]interface{})
	for k, v := range metadata {
		sanitizedMeta[k] = sl.sanitizeArgument(v)
	}
	event["metadata"] = sanitizedMeta

	// Serialize to JSON
	// Note: We use a separate encoder or simple string formatting to ensure JSON validity
	// For simplicity in this logger wrapper, we'll format it as a JSON string
	// In a real production env, use encoding/json
	jsonBytes, err := json.Marshal(event)
	if err != nil {
		sl.Error("Failed to marshal audit event: %v", err)
		return
	}

	// Direct output to logger, ensuring it's treated as a single line
	sl.logger.Println(string(jsonBytes))
}

// log handles the actual logging with sanitization
func (sl *SecureLogger) log(level, format string, v ...interface{}) {
	// Sanitize the format string and arguments
	sanitizedFormat := sl.sanitizeString(format)
	sanitizedArgs := make([]interface{}, len(v))

	for i, arg := range v {
		sanitizedArgs[i] = sl.sanitizeArgument(arg)
	}

	// Add level and timestamp
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf("[%s] [%s] %s", timestamp, level, sanitizedFormat)

	// Use the logger to output the message
	sl.logger.Printf(message, sanitizedArgs...)
}

// sanitizeString sanitizes a string by masking tokens and removing dangerous content
func (sl *SecureLogger) sanitizeString(s string) string {
	sl.mu.Lock()
	defer sl.mu.Unlock()

	// Mask all registered tokens
	for original, masked := range sl.maskTokens {
		s = strings.ReplaceAll(s, original, masked)
	}

	// Additional sanitization
	return SanitizeString(s)
}

// sanitizeArgument sanitizes any type of argument
func (sl *SecureLogger) sanitizeArgument(arg interface{}) interface{} {
	switch v := arg.(type) {
	case string:
		return sl.sanitizeString(v)
	case error:
		return sl.sanitizeString(v.Error())
	case fmt.Stringer:
		return sl.sanitizeString(v.String())
	default:
		return v
	}
}

// LogRateLimitStatus logs the current rate limit status
func (sl *SecureLogger) LogRateLimitStatus(limiter *RateLimiter) {
	if limiter == nil {
		return
	}
	sl.Info("Rate limit status: %s", limiter.GetStatus())
}

// Global secure logger instance
var globalLogger *SecureLogger
var loggerInit sync.Once

// GetLogger returns the global secure logger instance
func GetLogger() *SecureLogger {
	loggerInit.Do(func() {
		globalLogger = NewSecureLogger()
	})
	return globalLogger
}
