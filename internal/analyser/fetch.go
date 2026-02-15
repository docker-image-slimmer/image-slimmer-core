package analyzer

import (
	"context"
	"time"

	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

// fetchImage resolves and downloads a container image from a remote registry
// This function acts as a strict external boundary: all errors are normalized
func fetchImage(ctx context.Context, ref string, opts *options) (v1.Image, error) {
	const op = "fetch"

	// Guarantee timeout enforcement
	if _, hasDeadline := ctx.Deadline(); !hasDeadline {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, opts.timeout)
		defer cancel()
	}

	if ref == "" {
		return nil, NewError(CodeInvalidReference, op, ref, "image reference cannot be empty", nil)
	}

	// Strict parsing prevents ambiguous references
	parsedRef, err := name.ParseReference(ref, name.StrictValidation)
	if err != nil {
		return nil, NewError(CodeInvalidReference, op, ref, "invalid image reference format", err)
	}

	// Detect if reference is digest-pinned (more secure)
	_, isDigest := parsedRef.(name.Digest)

	// Prepare remote options (auth + transport extensible)
	remoteOpts := []remote.Option{
		remote.WithContext(ctx),
	}

	if opts.keychain != nil {
		remoteOpts = append(remoteOpts, remote.WithAuthFromKeychain(opts.keychain))
	}

	if opts.transport != nil {
		remoteOpts = append(remoteOpts, remote.WithTransport(opts.transport))
	}

	start := time.Now()

	img, err := remote.Image(parsedRef, remoteOpts...)

	fetchDuration := time.Since(start)

	if err != nil {
		return nil, MapRegistryError(op, ref, err)
	}

	if img == nil {
		return nil, NewError(CodeFetchFailed, op, ref, "registry returned nil image", nil)
	}

	// Force validation to ensure image is not partially resolved
	digest, err := img.Digest()
	if err != nil {
		return nil, NewError(CodeFetchFailed, op, ref, "failed to resolve image digest", err)
	}

	// If reference was digest-pinned, enforce digest match
	if isDigest {
		if parsedRef.Identifier() != digest.String() {
			return nil, NewError(
				CodeFetchFailed,
				op,
				ref,
				"digest mismatch between reference and remote image",
				nil,
			)
		}
	}

	// Attach fetch duration to context if metrics enabled
	if opts.metricsHook != nil {
		opts.metricsHook(FetchMetrics{
			Reference:    ref,
			Duration:     fetchDuration,
			Digest:       digest.String(),
			DigestPinned: isDigest,
		})
	}

	return img, nil
}
