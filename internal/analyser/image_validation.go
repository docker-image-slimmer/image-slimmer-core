package analyzer

// Validate ensures the Image is structurally consistent.
//
// It verifies:
//   - Non-nil image
//   - Required metadata presence
//   - Positive size
//   - Layer integrity
//
// All errors are normalized into AnalyzerError.
func (i *Image) Validate() error {
	const op = "validate"

	if i == nil {
		return NewError(
			CodeValidationFailed,
			op,
			"",
			"image is nil",
			nil,
		)
	}

	if i.Reference == "" {
		return NewError(
			CodeValidationFailed,
			op,
			i.Reference,
			"image reference is empty",
			nil,
		)
	}

	if i.Digest == "" {
		return NewError(
			CodeValidationFailed,
			op,
			i.Reference,
			"image digest is empty",
			nil,
		)
	}

	if i.Size <= 0 {
		return NewError(
			CodeValidationFailed,
			op,
			i.Reference,
			"image size is invalid",
			nil,
		)
	}

	// If metadataOnly mode was used, layers may be intentionally empty.
	if len(i.Layers) == 0 {
		return NewError(
			CodeValidationFailed,
			op,
			i.Reference,
			"image has no layers",
			nil,
		)
	}

	for _, l := range i.Layers {
		if l.Digest == "" {
			return NewError(
				CodeValidationFailed,
				op,
				i.Reference,
				"layer has empty digest",
				nil,
			)
		}

		if l.CompressedSize <= 0 {
			return NewError(
				CodeValidationFailed,
				op,
				i.Reference,
				"layer has invalid compressed size",
				nil,
			)
		}
	}

	return nil
}
