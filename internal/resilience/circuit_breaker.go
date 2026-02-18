package resilience

import (
	"context"
	"errors"
	"sync"
	"time"
)

type State int

const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

func (s State) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateOpen:
		return "open"
	case StateHalfOpen:
		return "half-open"
	default:
		return "unknown"
	}
}

var (
	ErrCircuitOpen     = errors.New("circuit breaker is open")
	ErrTooManyRequests = errors.New("too many requests in half-open state")
	ErrNilFunction     = errors.New("function cannot be nil")
)

type CircuitBreaker struct {
	mu               sync.RWMutex
	state            State
	failures         int
	successes        int
	failureThreshold int
	successThreshold int
	timeout          time.Duration
	lastFailureTime  time.Time
	onStateChange    func(from, to State)
}

type CircuitBreakerOption func(*CircuitBreaker)

func WithOnStateChange(fn func(from, to State)) CircuitBreakerOption {
	return func(cb *CircuitBreaker) {
		cb.onStateChange = fn
	}
}

func NewCircuitBreaker(failureThreshold, successThreshold int, timeout time.Duration, opts ...CircuitBreakerOption) *CircuitBreaker {
	if failureThreshold < 1 {
		failureThreshold = 5
	}
	if successThreshold < 1 {
		successThreshold = 1
	}
	if timeout < 1 {
		timeout = 30 * time.Second
	}

	cb := &CircuitBreaker{
		state:            StateClosed,
		failureThreshold: failureThreshold,
		successThreshold: successThreshold,
		timeout:          timeout,
	}

	for _, opt := range opts {
		opt(cb)
	}

	return cb
}

func (cb *CircuitBreaker) Execute(ctx context.Context, fn func() error) error {
	if fn == nil {
		return ErrNilFunction
	}

	if err := ctx.Err(); err != nil {
		return err
	}

	if !cb.allowRequest() {
		return ErrCircuitOpen
	}

	err := fn()

	select {
	case <-ctx.Done():
		cb.recordFailure()
		return ctx.Err()
	default:
	}

	if err != nil {
		cb.recordFailure()
		return err
	}

	cb.recordSuccess()
	return nil
}

func (cb *CircuitBreaker) allowRequest() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case StateClosed:
		return true
	case StateOpen:
		if time.Since(cb.lastFailureTime) > cb.timeout {
			cb.setState(StateHalfOpen)
			return true
		}
		return false
	case StateHalfOpen:
		return true
	default:
		return false
	}
}

func (cb *CircuitBreaker) recordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.failures++
	cb.lastFailureTime = time.Now()

	switch cb.state {
	case StateClosed:
		if cb.failures >= cb.failureThreshold {
			cb.setState(StateOpen)
		}
	case StateHalfOpen:
		cb.setState(StateOpen)
	}
}

func (cb *CircuitBreaker) recordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.failures = 0

	switch cb.state {
	case StateHalfOpen:
		cb.successes++
		if cb.successes >= cb.successThreshold {
			cb.setState(StateClosed)
		}
	}
}

func (cb *CircuitBreaker) setState(newState State) {
	if cb.state == newState {
		return
	}

	oldState := cb.state
	cb.state = newState

	if newState == StateClosed {
		cb.failures = 0
		cb.successes = 0
	} else if newState == StateHalfOpen {
		cb.successes = 0
	}

	if cb.onStateChange != nil {
		go cb.onStateChange(oldState, newState)
	}
}

func (cb *CircuitBreaker) State() State {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

func (cb *CircuitBreaker) Reset() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.setState(StateClosed)
}

func (cb *CircuitBreaker) FailureCount() int {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.failures
}

func (cb *CircuitBreaker) IsOpen() bool {
	return cb.State() == StateOpen
}

func (cb *CircuitBreaker) IsClosed() bool {
	return cb.State() == StateClosed
}

func (cb *CircuitBreaker) IsHalfOpen() bool {
	return cb.State() == StateHalfOpen
}

type CircuitBreakerRegistry struct {
	breakers map[string]*CircuitBreaker
	mu       sync.RWMutex
}

func NewCircuitBreakerRegistry() *CircuitBreakerRegistry {
	return &CircuitBreakerRegistry{
		breakers: make(map[string]*CircuitBreaker),
	}
}

func (r *CircuitBreakerRegistry) Get(name string, failureThreshold, successThreshold int, timeout time.Duration, opts ...CircuitBreakerOption) *CircuitBreaker {
	r.mu.RLock()
	breaker, exists := r.breakers[name]
	r.mu.RUnlock()

	if exists {
		return breaker
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if breaker, exists = r.breakers[name]; exists {
		return breaker
	}

	breaker = NewCircuitBreaker(failureThreshold, successThreshold, timeout, opts...)
	r.breakers[name] = breaker
	return breaker
}

func (r *CircuitBreakerRegistry) AllStats() map[string]struct {
	State     State
	Failures  int
	Successes int
} {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make(map[string]struct {
		State     State
		Failures  int
		Successes int
	})

	for name, breaker := range r.breakers {
		result[name] = struct {
			State     State
			Failures  int
			Successes int
		}{breaker.State(), breaker.FailureCount(), 0}
	}
	return result
}

func (r *CircuitBreakerRegistry) ResetAll() {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, breaker := range r.breakers {
		breaker.Reset()
	}
}

func AllStatsFromRegistry() map[string]struct {
	State     State
	Failures  int
	Successes int
} {
	return registry.AllStats()
}

var registry = NewCircuitBreakerRegistry()
