package analyzer

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

// Image represents a resolved container image ready for analysis.
type Image struct {
	Reference string
	Digest    string
	MediaType string
	Size      int64
	Layers    []Layer
	Raw       v1.Image
	LoadedAt  time.Time
}

// LoadImage loads a remote image with default timeout and retry.
func LoadImage(ref string) (*Image, error) {
	return LoadImageWithOptions(ref, 30*time.Second, 2)
}

// LoadImageWithOptions allows custom timeout and retry attempts.
func LoadImageWithOptions(ref string, timeout time.Duration, retries int) (*Image, error) {
	if ref == "" {
		return nil, fmt.Errorf("image reference cannot be empty")
	}

	parsedRef, err := name.ParseReference(ref)
	if err != nil {
		return nil, fmt.Errorf("invalid image reference %q: %w", ref, err)
	}

	var img v1.Image
	var lastErr error

	for attempt := 0; attempt <= retries; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)

		img, lastErr = remote.Image(parsedRef, remote.WithContext(ctx))
		cancel()

		if lastErr == nil {
			break
		}

		time.Sleep(time.Duration(attempt+1) * 500 * time.Millisecond)
	}

	if lastErr != nil {
		return nil, fmt.Errorf("failed to fetch remote image after retries: %w", lastErr)
	}

	digest, err := img.Digest()
	if err != nil {
		return nil, fmt.Errorf("failed to resolve image digest: %w", err)
	}

	mediaType, err := img.MediaType()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve image media type: %w", err)
	}

	size, err := img.Size()
	if err != nil {
		return nil, fmt.Errorf("failed to compute image size: %w", err)
	}

	rawLayers, err := img.Layers()
	if err != nil {
		return nil, fmt.Errorf("failed to extract image layers: %w", err)
	}

	if len(rawLayers) == 0 {
		return nil, fmt.Errorf("image contains no layers")
	}

	structuredLayers, err := ExtractLayers(rawLayers)
	if err != nil {
		return nil, fmt.Errorf("failed to structure layers: %w", err)
	}

	return &Image{
		Reference: ref,
		Digest:    digest.String(),
		MediaType: string(mediaType),
		Size:      size,
		Layers:    structuredLayers,
		Raw:       img,
		LoadedAt:  time.Now(),
	}, nil
}
