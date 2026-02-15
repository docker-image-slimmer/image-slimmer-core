package digest

import (
	"fmt"
	"strings"

	analyzer "github.com/pnkcaht/image-slimmer-core/internal/analyser"
)

// RiskLevel represents the severity of modifying or removing a layer
type RiskLevel string

const (
	RiskLow    RiskLevel = "low"
	RiskMedium RiskLevel = "medium"
	RiskHigh   RiskLevel = "high"
)

// LayerRisk captures the risk assessment of a single image layer
type LayerRisk struct {
	Index  int
	Digest string
	Level  RiskLevel
	Reason string
}

// AssessLayerRisk evaluates the potential risk of modifying or removing a layer
// Heuristics used:
// - Config layers are always high risk
// - Large layers (>100MB) are medium risk
// - All other layers are considered low risk
func AssessLayerRisk(layer analyzer.Layer) LayerRisk {
	risk := LayerRisk{
		Index:  layer.Index,
		Digest: layer.Digest,
		Level:  RiskLow,
		Reason: "normal layer",
	}

	if strings.Contains(strings.ToLower(layer.MediaType), "config") {
		risk.Level = RiskHigh
		risk.Reason = "configuration layer"
	} else if layer.UncompressedSize > 100*1024*1024 {
		risk.Level = RiskMedium
		risk.Reason = "large layer (>100MB)"
	}

	return risk
}

// AssessImageRisk evaluates all layers of the image and returns a slice of LayerRisk
// Returns an error if the image is nil or has no layers
// This function provides a comprehensive risk profile for the entire image, which can be used to inform slimming decisions
func AssessImageRisk(img *analyzer.Image) ([]LayerRisk, error) {
	if img == nil {
		return nil, fmt.Errorf("image is nil")
	}

	if len(img.Layers) == 0 {
		return nil, fmt.Errorf("image has no layers to assess")
	}

	risks := make([]LayerRisk, len(img.Layers))
	for i, layer := range img.Layers {
		risks[i] = AssessLayerRisk(layer)
	}

	return risks, nil
}
