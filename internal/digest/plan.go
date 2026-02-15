package digest

import (
	"fmt"
	"strings"

	analyzer "github.com/pnkcaht/image-slimmer-core/internal/analyser"
)

// LayerPlan represents the planned action for a specific image layer
type LayerPlan struct {
	Index       int
	Digest      string
	Action      string // "keep", "remove", "rebuild"
	Description string
}

// ImagePlan represents the overall plan for slimming an image, including actions for each layer
type ImagePlan struct {
	Reference string
	Digest    string
	Layers    []LayerPlan
}

// NewImagePlan creates a new ImagePlan based on the analyzed image data. It initializes all layers with a default action of "keep" and includes descriptive metadata for each layer. Returns an error if the input image is nil
func NewImagePlan(img *analyzer.Image) (*ImagePlan, error) {
	if img == nil {
		return nil, fmt.Errorf("image is nil")
	}

	layers := make([]LayerPlan, len(img.Layers))
	for i, l := range img.Layers {
		layers[i] = LayerPlan{
			Index:       l.Index,
			Digest:      l.Digest,
			Action:      "keep",
			Description: fmt.Sprintf("Layer %d size=%d mediaType=%s", l.Index, l.UncompressedSize, l.MediaType),
		}
	}

	return &ImagePlan{
		Reference: img.Reference,
		Digest:    img.Digest,
		Layers:    layers,
	}, nil
}

// MarkLayerForRemoval updates the plan to indicate that a specific layer should be removed. It also appends the provided reason to the layer's description for traceability. Returns an error if the specified layer index is not found in the plan
func (p *ImagePlan) MarkLayerForRemoval(index int, reason string) error {
	lp, err := p.findLayer(index)
	if err != nil {
		return err
	}
	lp.Action = "remove"
	lp.Description = strings.TrimSpace(lp.Description + " | removal reason: " + reason)
	return nil
}

// MarkLayerForRebuild updates the plan to indicate that a specific layer should be rebuilt. It also appends the provided reason to the layer's description for traceability. Returns an error if the specified layer index is not found in the plan
func (p *ImagePlan) MarkLayerForRebuild(index int, reason string) error {
	lp, err := p.findLayer(index)
	if err != nil {
		return err
	}
	lp.Action = "rebuild"
	lp.Description = strings.TrimSpace(lp.Description + " | rebuild reason: " + reason)
	return nil
}

// findLayer is a helper method to locate a LayerPlan by its index. Returns an error if the layer is not found in the plan
func (p *ImagePlan) findLayer(index int) (*LayerPlan, error) {
	for i := range p.Layers {
		if p.Layers[i].Index == index {
			return &p.Layers[i], nil
		}
	}
	return nil, fmt.Errorf("layer %d not found in plan", index)
}

// Summary generates a human-readable summary of the image plan, including the reference, digest, and planned actions for each layer. This is useful for debugging and communicating the slimming strategy to users or other components
func (p *ImagePlan) Summary() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Image Plan for %s (digest=%s)\n", p.Reference, p.Digest))
	for _, l := range p.Layers {
		sb.WriteString(fmt.Sprintf("- Layer %d: %s | Action: %s\n", l.Index, l.Digest, l.Action))
	}
	return sb.String()
}
