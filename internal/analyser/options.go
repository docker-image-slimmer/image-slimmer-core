package analyzer

import (
	"net/http"
	"time"

	"github.com/google/go-containerregistry/pkg/authn"
)

// FetchMetrics represents structured telemetry data emitted
// after an image fetch operation.
type FetchMetrics struct {
	// Reference is the original image reference string
	Reference string

	// Duration represents total time spent fetching
	Duration time.Duration

	// Digest is the resolved content digest of the image
	Digest string

	// DigestPinned indicates whether the original reference
	// was already pinned to a digest.
	DigestPinned bool

	// Attempts indicates how many attempts were performed
	Attempts int

	// Success indicates whether the fetch completed successfully
	Success bool
}

// options holds internal configuration for the analyzer
// It is intentionally unexported to enforce controlled construction
type options struct {
	timeout      time.Duration
	retries      int
	backoff      time.Duration
	keychain     authn.Keychain
	transport    http.RoundTripper
	metricsHook  func(FetchMetrics)
	metadataOnly bool
}

// Option defines a functional configuration modifier
type Option func(*options)

// defaultOptions returns production-safe default configuration
// These defaults are conservative and registry-safe
func defaultOptions() *options {
	return &options{
		timeout:      30 * time.Second,
		retries:      2,
		backoff:      500 * time.Millisecond,
		keychain:     authn.DefaultKeychain,
		transport:    http.DefaultTransport,
		metadataOnly: false,
	}
}

// WithTimeout configures the maximum allowed duration
// for registry communication operations
func WithTimeout(d time.Duration) Option {
	return func(o *options) {
		if d > 0 {
			o.timeout = d
		}
	}
}

// WithRetries configures how many retry attempts
// are allowed for transient registry failures
func WithRetries(n int) Option {
	return func(o *options) {
		if n >= 0 {
			o.retries = n
		}
	}
}

// WithBackoff configures the base delay between retry attempts
// This value is typically used as the starting backoff duration
func WithBackoff(d time.Duration) Option {
	return func(o *options) {
		if d > 0 {
			o.backoff = d
		}
	}
}

// WithKeychain configures the credential resolution chain used for registry authentication
func WithKeychain(k authn.Keychain) Option {
	return func(o *options) {
		if k != nil {
			o.keychain = k
		}
	}
}

// WithTransport configures a custom HTTP transport
// This allows advanced control such as proxy, TLS, or tracing
func WithTransport(t http.RoundTripper) Option {
	return func(o *options) {
		if t != nil {
			o.transport = t
		}
	}
}

// WithMetricsHook registers a callback invoked after
// a fetch operation completes. It can be used for
// logging, telemetry, or observability integration
func WithMetricsHook(h func(FetchMetrics)) Option {
	return func(o *options) {
		o.metricsHook = h
	}
}

// WithMetadataOnly configures the analyzer to skip layer extraction
// and return only high-level image metadata (digest, size, media type)
func WithMetadataOnly(enabled bool) Option {
	return func(o *options) {
		o.metadataOnly = enabled
	}
}
