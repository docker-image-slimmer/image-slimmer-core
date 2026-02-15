package analyzer

import "fmt"

// Validate ensures the image is structurally consistent.
func (i *Image) Validate() error {
	if i == nil {
		return fmt.Errorf("image is nil")
	}

	if i.Reference == "" {
		return fmt.Errorf("image reference is empty")
	}

	if i.Digest == "" {
		return fmt.Errorf("image digest is empty")
	}

	if i.Size <= 0 {
		return fmt.Errorf("image size is invalid")
	}

	if len(i.Layers) == 0 {
		return fmt.Errorf("image has no layers")
	}

	for _, l := range i.Layers {
		if l.Digest == "" {
			return fmt.Errorf("layer %d has empty digest", l.Index)
		}
		if l.CompressedSize <= 0 {
			return fmt.Errorf("layer %d has invalid compressed size", l.Index)
		}
	}

	return nil
}
