package analyzer

import (
	"time"

	v1 "github.com/google/go-containerregistry/pkg/v1"
)

// buildImage constructs a structured Image from a resolved v1.Image
//
// This function represents the internal build boundary of the analyzer
// All external errors are normalized into AnalyzerError
//
// Guarantees:
//   - No raw registry errors leak outside
//   - Strict structural validation
//   - Deterministic metadata extraction
func buildImage(ref string, img v1.Image, opts *options) (*Image, error) {
	const op = "build"

	// ---- NIL IMAGE CHECK ----
	if img == nil {
		return nil, NewError(
			CodeBuildFailed,
			op,
			ref,
			"cannot build from nil image",
			nil,
		)
	}

	// ---- DIGEST ----
	digest, err := img.Digest()
	if err != nil {
		return nil, NewError(
			CodeDigestFailed,
			op,
			ref,
			"failed to resolve image digest",
			err,
		)
	}

	// ---- MEDIA TYPE ----
	mediaType, err := img.MediaType()
	if err != nil {
		return nil, NewError(
			CodeMediaTypeFailed,
			op,
			ref,
			"failed to resolve image media type",
			err,
		)
	}

	// ---- SIZE ----
	size, err := img.Size()
	if err != nil {
		return nil, NewError(
			CodeSizeFailed,
			op,
			ref,
			"failed to resolve image size",
			err,
		)
	}

	// ---- LAYERS (optional) ----
	var structuredLayers []Layer

	if !opts.metadataOnly {
		rawLayers, err := img.Layers()
		if err != nil {
			return nil, NewError(
				CodeBuildFailed,
				op,
				ref,
				"failed to retrieve image layers",
				err,
			)
		}

		if len(rawLayers) == 0 {
			return nil, NewError(
				CodeNoLayers,
				op,
				ref,
				"image contains no layers",
				nil,
			)
		}

		// ---- EXTRACT LAYERS ----
		structuredLayers, err = ExtractLayers(rawLayers, ref)
		if err != nil {
			return nil, NewError(
				CodeLayerExtract,
				op,
				ref,
				"failed to extract structured layers",
				err,
			)
		}
	}

	// ---- FINAL STRUCTURE ----
	return &Image{
		Reference: ref,
		Digest:    digest.String(),
		MediaType: string(mediaType),
		Size:      size,
		Layers:    structuredLayers,
		LoadedAt:  time.Now(),
	}, nil
}
