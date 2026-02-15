package analyzer

import (
	"net/http"
	"time"

	"github.com/google/go-containerregistry/pkg/authn"
)

// FetchMetrics contains structured information about a fetch operation.
type FetchMetrics struct {
	Reference    string
	Duration     time.Duration
	Digest       string
	DigestPinned bool
}

// options defines internal configuration for the analyzer.
type options struct {
	timeout     time.Duration
	auth        authn.Authenticator
	transport   http.RoundTripper
	metricsHook func(FetchMetrics)
}

// Option is a functional option for configuring analyzer behavior.
type Option func(*options)

// defaultOptions returns production-safe defaults.
func defaultOptions() *options {
	return &options{
		timeout: 30 * time.Second,
	}
}

// WithTimeout sets the maximum duration allowed for registry operations.
func WithTimeout(d time.Duration) Option {
	return func(o *options) {
		o.timeout = d
	}
}

// WithAuth configures registry authentication.
func WithAuth(a authn.Authenticator) Option {
	return func(o *options) {
		o.auth = a
	}
}

// WithTransport allows custom HTTP transport (proxy, TLS config, etc).
func WithTransport(t http.RoundTripper) Option {
	return func(o *options) {
		o.transport = t
	}
}

// WithMetricsHook registers a callback for fetch metrics reporting.
func WithMetricsHook(h func(FetchMetrics)) Option {
	return func(o *options) {
		o.metricsHook = h
	}
}
