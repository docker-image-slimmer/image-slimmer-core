package analyzer

import (
	"io"

	v1 "github.com/google/go-containerregistry/pkg/v1"
)

// Layer represents extracted metadata from a container image layer
type Layer struct {
	Index            int
	Digest           string
	DiffID           string
	MediaType        string
	CompressedSize   int64
	UncompressedSize int64
}

// ExtractLayers converts raw v1 layers into structured Layer metadata
// All errors are normalized to AnalyzerError.
func ExtractLayers(rawLayers []v1.Layer, ref string) ([]Layer, error) {
	const op = "extract_layers"

	if len(rawLayers) == 0 {
		return nil, NewError(CodeLayerExtract, op, ref, "no layers to extract", nil)
	}

	layers := make([]Layer, 0, len(rawLayers))

	for i, l := range rawLayers {

		digest, err := l.Digest()
		if err != nil {
			return nil, NewError(CodeLayerExtract, op, ref, "failed to get layer digest", err)
		}

		diffID, err := l.DiffID()
		if err != nil {
			return nil, NewError(CodeLayerExtract, op, ref, "failed to get layer diffID", err)
		}

		mediaType, err := l.MediaType()
		if err != nil {
			return nil, NewError(CodeLayerExtract, op, ref, "failed to get layer media type", err)
		}

		compressedSize, err := l.Size()
		if err != nil {
			return nil, NewError(CodeLayerExtract, op, ref, "failed to get compressed size", err)
		}

		// Calculate uncompressed size manually
		rc, err := l.Uncompressed()
		if err != nil {
			return nil, NewError(CodeLayerExtract, op, ref, "failed to get uncompressed reader", err)
		}

		uncompressedSize, err := io.Copy(io.Discard, rc)
		rc.Close()
		if err != nil {
			return nil, NewError(CodeLayerExtract, op, ref, "failed to compute uncompressed size", err)
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
