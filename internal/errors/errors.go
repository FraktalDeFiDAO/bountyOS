package errors

import (
	"errors"
	"fmt"
	"time"
)

type ErrorType int

const (
	ErrorTypeTransient ErrorType = iota
	ErrorTypePermanent
	ErrorTypeTimeout
	ErrorTypeRateLimit
	ErrorTypeNetwork
)

func (t ErrorType) String() string {
	switch t {
	case ErrorTypeTransient:
		return "transient"
	case ErrorTypePermanent:
		return "permanent"
	case ErrorTypeTimeout:
		return "timeout"
	case ErrorTypeRateLimit:
		return "rate_limit"
	case ErrorTypeNetwork:
		return "network"
	default:
		return "unknown"
	}
}

type BountyError struct {
	Type       ErrorType
	Operation  string
	Message    string
	Cause      error
	Retryable  bool
	RetryAfter time.Duration
}

func (e *BountyError) Error() string {
	msg := fmt.Sprintf("[%s] %s", e.Type, e.Operation)
	if e.Message != "" {
		msg += ": " + e.Message
	}
	if e.Cause != nil {
		msg += fmt.Sprintf(" (cause: %v)", e.Cause)
	}
	if e.RetryAfter > 0 {
		msg += fmt.Sprintf(" (retry-after: %v)", e.RetryAfter)
	}
	return msg
}

func (e *BountyError) Unwrap() error {
	return e.Cause
}

func (e *BountyError) Is(target error) bool {
	switch target {
	case ErrTransient:
		return e.Type == ErrorTypeTransient
	case ErrPermanent:
		return e.Type == ErrorTypePermanent
	case ErrTimeout:
		return e.Type == ErrorTypeTimeout
	case ErrRateLimit:
		return e.Type == ErrorTypeRateLimit
	case ErrNetwork:
		return e.Type == ErrorTypeNetwork
	default:
		return false
	}
}

var (
	ErrTransient = &BountyError{Type: ErrorTypeTransient}
	ErrPermanent = &BountyError{Type: ErrorTypePermanent}
	ErrTimeout   = &BountyError{Type: ErrorTypeTimeout}
	ErrRateLimit = &BountyError{Type: ErrorTypeRateLimit}
	ErrNetwork   = &BountyError{Type: ErrorTypeNetwork}
)

func NewTransientError(op, msg string, cause error) *BountyError {
	return &BountyError{
		Type:      ErrorTypeTransient,
		Operation: op,
		Message:   msg,
		Cause:     cause,
		Retryable: true,
	}
}

func NewPermanentError(op, msg string, cause error) *BountyError {
	return &BountyError{
		Type:      ErrorTypePermanent,
		Operation: op,
		Message:   msg,
		Cause:     cause,
		Retryable: false,
	}
}

func NewTimeoutError(op, msg string, cause error) *BountyError {
	return &BountyError{
		Type:      ErrorTypeTimeout,
		Operation: op,
		Message:   msg,
		Cause:     cause,
		Retryable: true,
	}
}

func NewRateLimitError(op, msg string, retryAfter time.Duration) *BountyError {
	return &BountyError{
		Type:       ErrorTypeRateLimit,
		Operation:  op,
		Message:    msg,
		Retryable:  true,
		RetryAfter: retryAfter,
	}
}

func NewNetworkError(op, msg string, cause error) *BountyError {
	return &BountyError{
		Type:      ErrorTypeNetwork,
		Operation: op,
		Message:   msg,
		Cause:     cause,
		Retryable: true,
	}
}

func IsRetryable(err error) bool {
	if err == nil {
		return false
	}

	var bountyErr *BountyError
	if errors.As(err, &bountyErr) {
		return bountyErr.Retryable
	}

	return false
}

func GetRetryAfter(err error) time.Duration {
	if err == nil {
		return 0
	}

	var bountyErr *BountyError
	if errors.As(err, &bountyErr) {
		return bountyErr.RetryAfter
	}

	return 0
}

func GetErrorType(err error) ErrorType {
	if err == nil {
		return ErrorTypePermanent
	}

	var bountyErr *BountyError
	if errors.As(err, &bountyErr) {
		return bountyErr.Type
	}

	return ErrorTypePermanent
}

func IsTransient(err error) bool {
	return GetErrorType(err) == ErrorTypeTransient
}

func IsPermanent(err error) bool {
	return GetErrorType(err) == ErrorTypePermanent
}

func IsTimeout(err error) bool {
	return GetErrorType(err) == ErrorTypeTimeout
}

func IsRateLimit(err error) bool {
	return GetErrorType(err) == ErrorTypeRateLimit
}

func IsNetwork(err error) bool {
	return GetErrorType(err) == ErrorTypeNetwork
}

func Wrap(err error, op, msg string) *BountyError {
	if err == nil {
		return nil
	}

	var bountyErr *BountyError
	if errors.As(err, &bountyErr) {
		return &BountyError{
			Type:      bountyErr.Type,
			Operation: op,
			Message:   msg,
			Cause:     err,
			Retryable: bountyErr.Retryable,
		}
	}

	return &BountyError{
		Type:      ErrorTypePermanent,
		Operation: op,
		Message:   msg,
		Cause:     err,
		Retryable: false,
	}
}

func WrapAsTransient(err error, op, msg string) *BountyError {
	return &BountyError{
		Type:      ErrorTypeTransient,
		Operation: op,
		Message:   msg,
		Cause:     err,
		Retryable: true,
	}
}

func WrapAsNetwork(err error, op, msg string) *BountyError {
	return &BountyError{
		Type:      ErrorTypeNetwork,
		Operation: op,
		Message:   msg,
		Cause:     err,
		Retryable: true,
	}
}
