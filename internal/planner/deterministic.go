package planner

import (
	"fmt"
	"sort"

	analyser "github.com/pnkcaht/image-slimmer-core/internal/analyser"
)

// DeterministicImage represents a normalized image for deterministic operations.
type DeterministicImage struct {
	Reference string
	Digest    string
	Layers    []analyser.Layer
}

// NewDeterministicImage creates a deterministic image from a resolved image.
// Layers are sorted by index to guarantee reproducibility.
func NewDeterministicImage(img *analyser.Image) (*DeterministicImage, error) {
	if img == nil {
		return nil, fmt.Errorf("image is nil")
	}

	if len(img.Layers) == 0 {
		return nil, fmt.Errorf("image has no layers")
	}

	// Copy layers to avoid mutating the original
	layers := make([]analyser.Layer, len(img.Layers))
	copy(layers, img.Layers)

	// Sort layers by index for deterministic ordering
	sort.Slice(layers, func(i, j int) bool {
		return layers[i].Index < layers[j].Index
	})

	return &DeterministicImage{
		Reference: img.Reference,
		Digest:    img.Digest,
		Layers:    layers,
	}, nil
}

// LayerHashes returns the ordered digests of all layers.
func (d *DeterministicImage) LayerHashes() []string {
	hashes := make([]string, len(d.Layers))
	for i, l := range d.Layers {
		hashes[i] = l.Digest
	}
	return hashes
}

// LayerMediaTypes returns the ordered media types of all layers.
func (d *DeterministicImage) LayerMediaTypes() []string {
	media := make([]string, len(d.Layers))
	for i, l := range d.Layers {
		media[i] = l.MediaType
	}
	return media
}

// Summary provides a concise deterministic view of the image.
func (d *DeterministicImage) Summary() string {
	s := fmt.Sprintf("Deterministic Image: %s (digest=%s)\n", d.Reference, d.Digest)
	for _, l := range d.Layers {
		s += fmt.Sprintf("- Layer %d: %s | mediaType=%s | uncompressed=%d bytes\n",
			l.Index, l.Digest, l.MediaType, l.UncompressedSize)
	}
	return s
}
