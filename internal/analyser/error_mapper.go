package analyzer

import (
	"context"
	"errors"
	"net"
	"net/url"
	"strings"
)

// MapRegistryError converts external registry-related errors into a structured AnalyzerError
// It ensures that no external or raw transport errors leak outside the analyzer boundary
func MapRegistryError(op, ref string, err error) error {
	// Fast path: no error
	if err == nil {
		return nil
	}

	// If already mapped to AnalyzerError, return as-is
	var ae *AnalyzerError
	if errors.As(err, &ae) {
		return err
	}

	// Context deadline exceeded
	if errors.Is(err, context.DeadlineExceeded) {
		return NewError(CodeTimeout, op, ref, "operation timeout", err)
	}

	// Context explicitly canceled
	if errors.Is(err, context.Canceled) {
		return NewError(CodeTimeout, op, ref, "operation canceled", err)
	}

	// URL-level errors (very common in registry HTTP calls)
	var urlErr *url.Error
	if errors.As(err, &urlErr) {
		return classifyNetworkError(op, ref, urlErr.Err)
	}

	// Generic network errors
	var netErr net.Error
	if errors.As(err, &netErr) {
		// Timeout at transport level
		if netErr.Timeout() {
			return NewError(CodeTimeout, op, ref, "network timeout", err)
		}
		// Other network-related failures
		return NewError(CodeFetchFailed, op, ref, "network error", err)
	}

	// Fallback classification using error message inspection
	// Some registry implementations return inconsistent or wrapped errors
	msg := strings.ToLower(err.Error())

	switch {
	// Unauthorized / authentication failures
	case strings.Contains(msg, "401"),
		strings.Contains(msg, "unauthorized"),
		strings.Contains(msg, "denied"):
		return NewError(CodeUnauthorized, op, ref, "unauthorized access", err)

	// Forbidden access
	case strings.Contains(msg, "403"),
		strings.Contains(msg, "forbidden"):
		return NewError(CodeUnauthorized, op, ref, "forbidden", err)

	// Image not found
	case strings.Contains(msg, "404"),
		strings.Contains(msg, "not found"):
		return NewError(CodeImageNotFound, op, ref, "image not found", err)

	// Timeout indicated in message
	case strings.Contains(msg, "timeout"):
		return NewError(CodeTimeout, op, ref, "request timeout", err)

	// Rate limiting
	case strings.Contains(msg, "429"),
		strings.Contains(msg, "too many requests"):
		return NewError(CodeFetchFailed, op, ref, "rate limited", err)

	// Registry server errors (5xx)
	case strings.Contains(msg, "500"),
		strings.Contains(msg, "502"),
		strings.Contains(msg, "503"),
		strings.Contains(msg, "504"):
		return NewError(CodeFetchFailed, op, ref, "registry server error", err)
	}

	// Final fallback classification
	return NewError(CodeUnknown, op, ref, "unknown registry error", err)
}

// classifyNetworkError normalizes lower-level transport errors into structured analyzer errors
func classifyNetworkError(op, ref string, err error) error {
	var netErr net.Error
	if errors.As(err, &netErr) {
		// Explicit network timeout
		if netErr.Timeout() {
			return NewError(CodeTimeout, op, ref, "network timeout", err)
		}
	}

	// Any other transport-level failure
	return NewError(CodeFetchFailed, op, ref, "network error", err)
}
