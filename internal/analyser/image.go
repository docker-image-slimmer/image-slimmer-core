package analyzer

import (
	"context"
	"time"

	v1 "github.com/google/go-containerregistry/pkg/v1"
)

// Image represents a fully resolved container image ready for analysis
// It contains normalized metadata extracted from a remote registry
type Image struct {
	Reference string
	Digest    string
	MediaType string
	Size      int64
	Layers    []Layer
	LoadedAt  time.Time
}

// Load resolves and builds a container image from a remote reference
//
// It performs:
//   - Strict reference validation
//   - Controlled retry with exponential backoff
//   - Structured error normalization
//   - Metrics collection
//
// It returns the resolved Image, execution Metrics and an error (if any)
func Load(ctx context.Context, ref string, opts ...Option) (*Image, Metrics, error) {
	options := defaultOptions()
	for _, opt := range opts {
		opt(options)
	}

	collector := newMetricsCollector()

	// ---- FETCH PHASE ----

	collector.startFetch()

	var rawImg v1.Image
	var fetchErr error
	var attempts int

	attempts, fetchErr = retry(ctx, options.retries, options.backoff, func() error {
		var err error
		rawImg, err = fetchImage(ctx, ref, options)
		return err
	})

	collector.endFetch(attempts, false) // digestPinned can be improved later

	if fetchErr != nil {
		collector.markSuccess(false)
		return nil, collector.snapshot(), fetchErr
	}

	// ---- BUILD PHASE ----

	collector.startBuild()

	image, buildErr := buildImage(ref, rawImg, options)

	collector.endBuild()

	if buildErr != nil {
		collector.markSuccess(false)
		return nil, collector.snapshot(), buildErr
	}

	collector.markSuccess(true)

	return image, collector.snapshot(), nil
}
