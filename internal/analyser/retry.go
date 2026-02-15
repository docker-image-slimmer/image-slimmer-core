package analyzer

import (
	"context"
	"time"
)

// retry executes fn and retries it according to the provided configuration
//
// attempts defines the maximum number of retries (not total executions)
// Example: attempts=2 means:
//
//	1 initial execution + 2 retries = up to 3 total executions
//
// The retry mechanism guarantees:
//   - Respect for context cancellation and deadlines
//   - Exponential backoff with capped delay
//   - Retry only for temporary AnalyzerError instances
//   - Protection against duration overflow
//
// It returns the number of executions performed (including the first attempt)
// and the final error (if any)
func retry(
	ctx context.Context,
	attempts int,
	baseDelay time.Duration,
	fn func() error,
) (int, error) {

	if attempts < 0 {
		attempts = 0
	}

	if baseDelay <= 0 {
		baseDelay = 500 * time.Millisecond
	}

	var err error
	executions := 0

	for attempt := 0; attempt <= attempts; attempt++ {

		// Check context before executing the function
		select {
		case <-ctx.Done():
			return executions, ctx.Err()
		default:
		}

		executions++

		err = fn()
		if err == nil {
			return executions, nil
		}

		// Fail fast if error is not retryable
		if !isRetryable(err) {
			return executions, err
		}

		// If this was the last allowed attempt, stop retrying
		if attempt == attempts {
			break
		}

		// Compute exponential backoff delay (capped)
		delay := exponentialBackoff(baseDelay, attempt)

		// Respect remaining context deadline (timeout budget enforcement)
		if deadline, ok := ctx.Deadline(); ok {
			remaining := time.Until(deadline)
			if remaining <= 0 {
				return executions, ctx.Err()
			}
			if delay > remaining {
				delay = remaining
			}
		}

		// Use timer instead of time.Sleep to remain cancelable
		timer := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			timer.Stop()
			return executions, ctx.Err()
		case <-timer.C:
		}
	}

	return executions, err
}

// isRetryable determines whether an error should trigger retry
//
// Only AnalyzerError marked as Temporary() are considered retryable
// Unknown or external errors are treated as non-retryable to preserve
// deterministic and safe behavior at the analyzer boundary
func isRetryable(err error) bool {
	if ae, ok := AsAnalyzerError(err); ok {
		return ae.Temporary()
	}
	return false
}

// exponentialBackoff calculates a capped exponential delay
//
// delay = base * 2^attempt
//
// Guarantees:
//   - No overflow
//   - Upper bound enforcement
//   - Deterministic behavior
func exponentialBackoff(base time.Duration, attempt int) time.Duration {

	const maxDelay = 30 * time.Second

	// Prevent excessive bit shifting that may overflow
	if attempt > 30 {
		return maxDelay
	}

	delay := base << attempt

	// Guard against overflow or excessive growth
	if delay <= 0 || delay > maxDelay {
		return maxDelay
	}

	return delay
}
