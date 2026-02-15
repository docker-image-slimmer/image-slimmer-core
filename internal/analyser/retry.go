package analyzer

import (
	"context"
	"time"
)

// retry executes fn and retries it according to the provided configuration.
// attempts defines the maximum number of retries (not total executions).
// Example: attempts=2 means 1 initial try + 2 retries = up to 3 executions.
func retry(
	ctx context.Context,
	attempts int,
	baseDelay time.Duration,
	fn func() error,
) error {

	if attempts < 0 {
		attempts = 0
	}

	var err error

	for attempt := 0; attempt <= attempts; attempt++ {
		// Check context before executing
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		err = fn()
		if err == nil {
			return nil
		}

		// If not retryable, fail fast
		if !isRetryable(err) {
			return err
		}

		// If last attempt, break immediately
		if attempt == attempts {
			break
		}

		// Exponential backoff with cap (max 30s)
		delay := exponentialBackoff(baseDelay, attempt)

		timer := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
		}
	}

	return err
}

// isRetryable determines whether an error should trigger retry.
func isRetryable(err error) bool {
	if ae, ok := AsAnalyzerError(err); ok {
		return ae.Temporary()
	}

	// Unknown errors are treated as non-retryable
	return false
}

// exponentialBackoff calculates capped exponential delay.
func exponentialBackoff(base time.Duration, attempt int) time.Duration {
	if base <= 0 {
		base = 500 * time.Millisecond
	}

	// Prevent overflow
	maxDelay := 30 * time.Second

	delay := base << attempt
	if delay <= 0 || delay > maxDelay {
		return maxDelay
	}

	return delay
}
