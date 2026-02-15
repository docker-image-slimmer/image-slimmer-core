package slimmer

import (
	"context"
	"fmt"

	analyser "github.com/pnkcaht/image-slimmer-core/internal/analyser"
	digest "github.com/pnkcaht/image-slimmer-core/internal/digest"
	planner "github.com/pnkcaht/image-slimmer-core/internal/planner"
)

type Engine struct{}

func New() *Engine {
	return &Engine{}
}

type Result struct {
	Image         *analyser.Image
	Metrics       analyser.Metrics
	Deterministic *planner.DeterministicImage
	Plan          *digest.ImagePlan
}

func (e *Engine) Slim(ctx context.Context, ref string) (*Result, error) {

	// Load & analyze image
	img, metrics, err := analyser.Load(ctx, ref)
	if err != nil {
		return nil, fmt.Errorf("load failed: %w", err)
	}

	// Normalize deterministically
	det, err := planner.NewDeterministicImage(img)
	if err != nil {
		return nil, fmt.Errorf("deterministic normalization failed: %w", err)
	}

	// Create slimming plan
	plan, err := digest.NewImagePlan(img)
	if err != nil {
		return nil, fmt.Errorf("plan creation failed: %w", err)
	}

	return &Result{
		Image:         img,
		Metrics:       metrics,
		Deterministic: det,
		Plan:          plan,
	}, nil
}
