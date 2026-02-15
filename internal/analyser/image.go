package analyzer

import (
	"fmt"

	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

type Image struct {
	Reference string
	Digest    string
	MediaType string
	Size      int64
	Layers    []v1.Layer
	Raw       v1.Image
}

func LoadImage(ref string) (*Image, error) {
	if ref == "" {
		return nil, fmt.Errorf("image reference cannot be empty")
	}

	parsedRef, err := name.ParseReference(ref)
	if err != nil {
		return nil, fmt.Errorf("invalid image reference: %w", err)
	}

	img, err := remote.Image(parsedRef)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch image: %w", err)
	}

	digest, err := img.Digest()
	if err != nil {
		return nil, fmt.Errorf("failed to resolve digest: %w", err)
	}

	mediaType, err := img.MediaType()
	if err != nil {
		return nil, fmt.Errorf("failed to get media type: %w", err)
	}

	size, err := img.Size()
	if err != nil {
		return nil, fmt.Errorf("failed to get image size: %w", err)
	}

	layers, err := img.Layers()
	if err != nil {
		return nil, fmt.Errorf("failed to get layers: %w", err)
	}

	return &Image{
		Reference: ref,
		Digest:    digest.String(),
		MediaType: string(mediaType),
		Size:      size,
		Layers:    layers,
		Raw:       img,
	}, nil
}
