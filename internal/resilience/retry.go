package resilience

import (
	"context"
	"errors"
	"math/rand"
	"time"
)

type RetryConfig struct {
	MaxRetries     int
	InitialDelay   time.Duration
	MaxDelay       time.Duration
	RetryableError func(error) bool
}

type RetryOption func(*RetryConfig)

func WithMaxRetries(n int) RetryOption {
	return func(cfg *RetryConfig) {
		cfg.MaxRetries = n
	}
}

func WithInitialDelay(d time.Duration) RetryOption {
	return func(cfg *RetryConfig) {
		cfg.InitialDelay = d
	}
}

func WithMaxDelay(d time.Duration) RetryOption {
	return func(cfg *RetryConfig) {
		cfg.MaxDelay = d
	}
}

func WithRetryableError(fn func(error) bool) RetryOption {
	return func(cfg *RetryConfig) {
		cfg.RetryableError = fn
	}
}

func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries:     3,
		InitialDelay:   100 * time.Millisecond,
		MaxDelay:       30 * time.Second,
		RetryableError: IsRetryableDefault,
	}
}

func IsRetryableDefault(err error) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}
	if errors.Is(err, context.Canceled) {
		return false
	}

	return false
}

func Retry(ctx context.Context, cfg RetryConfig, fn func() error) error {
	if cfg.MaxRetries < 0 {
		cfg.MaxRetries = 0
	}
	if cfg.InitialDelay <= 0 {
		cfg.InitialDelay = 100 * time.Millisecond
	}
	if cfg.MaxDelay <= 0 {
		cfg.MaxDelay = 30 * time.Second
	}
	if cfg.RetryableError == nil {
		cfg.RetryableError = IsRetryableDefault
	}

	var lastErr error
	for attempt := 0; attempt <= cfg.MaxRetries; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}

		err := fn()
		if err == nil {
			return nil
		}

		lastErr = err

		if !cfg.RetryableError(err) {
			return err
		}

		if attempt == cfg.MaxRetries {
			break
		}

		delay := calculateBackoff(cfg.InitialDelay, cfg.MaxDelay, attempt)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}
	}

	return lastErr
}

func RetryWithBackoff(ctx context.Context, maxRetries int, initialDelay, maxDelay time.Duration, fn func() error) error {
	cfg := RetryConfig{
		MaxRetries:     maxRetries,
		InitialDelay:   initialDelay,
		MaxDelay:       maxDelay,
		RetryableError: IsRetryableDefault,
	}
	return Retry(ctx, cfg, fn)
}

func RetryWithOpts(ctx context.Context, fn func() error, opts ...RetryOption) error {
	cfg := DefaultRetryConfig()
	for _, opt := range opts {
		opt(&cfg)
	}
	return Retry(ctx, cfg, fn)
}

func calculateBackoff(initialDelay, maxDelay time.Duration, attempt int) time.Duration {
	delay := initialDelay * time.Duration(1<<uint(attempt))

	jitter := time.Duration(rand.Int63n(int64(delay) / 2))
	delay = delay + jitter

	if delay > maxDelay {
		delay = maxDelay
	}

	if delay < initialDelay {
		delay = initialDelay
	}

	return delay
}

type RetryableError struct {
	Err     error
	Message string
}

func (e *RetryableError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return e.Err.Error()
}

func (e *RetryableError) Unwrap() error {
	return e.Err
}

func NewRetryableError(err error, message string) *RetryableError {
	return &RetryableError{
		Err:     err,
		Message: message,
	}
}

func IsRetryableError(err error) bool {
	var retryable *RetryableError
	return errors.As(err, &retryable)
}
