package analyzer

import (
	"fmt"
	"io"

	v1 "github.com/google/go-containerregistry/pkg/v1"
)

// Layer represents extracted metadata from a container image layer.
type Layer struct {
	Index            int
	Digest           string
	DiffID           string
	MediaType        string
	CompressedSize   int64
	UncompressedSize int64
}

// ExtractLayers converts raw v1 layers into structured Layer metadata.
func ExtractLayers(rawLayers []v1.Layer) ([]Layer, error) {
	if len(rawLayers) == 0 {
		return nil, fmt.Errorf("no layers to extract")
	}

	layers := make([]Layer, 0, len(rawLayers))

	for i, l := range rawLayers {

		digest, err := l.Digest()
		if err != nil {
			return nil, fmt.Errorf("failed to get layer digest (index %d): %w", i, err)
		}

		diffID, err := l.DiffID()
		if err != nil {
			return nil, fmt.Errorf("failed to get layer diffID (index %d): %w", i, err)
		}

		mediaType, err := l.MediaType()
		if err != nil {
			return nil, fmt.Errorf("failed to get layer media type (index %d): %w", i, err)
		}

		compressedSize, err := l.Size()
		if err != nil {
			return nil, fmt.Errorf("failed to get compressed size (index %d): %w", i, err)
		}

		// Calculate uncompressed size manually
		rc, err := l.Uncompressed()
		if err != nil {
			return nil, fmt.Errorf("failed to get uncompressed reader (index %d): %w", i, err)
		}

		uncompressedSize, err := io.Copy(io.Discard, rc)
		rc.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to compute uncompressed size (index %d): %w", i, err)
		}

		layers = append(layers, Layer{
			Index:            i,
			Digest:           digest.String(),
			DiffID:           diffID.String(),
			MediaType:        string(mediaType),
			CompressedSize:   compressedSize,
			UncompressedSize: uncompressedSize,
		})
	}

	return layers, nil
}
