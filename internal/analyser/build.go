package analyzer

import (
	"time"

	v1 "github.com/google/go-containerregistry/pkg/v1"
)

func buildImage(ref string, img v1.Image, opts *options) (*Image, error) {
	start := time.Now()

	digest, err := img.Digest()
	if err != nil {
		return nil, err
	}

	mediaType, err := img.MediaType()
	if err != nil {
		return nil, err
	}

	size, err := img.Size()
	if err != nil {
		return nil, err
	}

	var structuredLayers []Layer

	if !opts.metadataOnly {
		rawLayers, err := img.Layers()
		if err != nil {
			return nil, err
		}

		if len(rawLayers) == 0 {
			return nil, ErrNoLayers
		}

		structuredLayers, err = ExtractLayers(rawLayers)
		if err != nil {
			return nil, err
		}
	}

	return &Image{
		Reference: ref,
		Digest:    digest.String(),
		MediaType: string(mediaType),
		Size:      size,
		Layers:    structuredLayers,
		LoadedAt:  time.Now(),
	}, nil
}
