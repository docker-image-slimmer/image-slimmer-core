package analyzer

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// ExtractLayerToFS writes the contents of a single Layer to the target directory.
// It uses the Layer info from Image.Layers.
func ExtractLayerToFS(layer Layer, layerReader io.ReadCloser, targetDir string) error {
	defer layerReader.Close()

	tr := tar.NewReader(layerReader)

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return NewError(CodeLayerExtract, "filesystem", "", "failed to read tar entry", err)
		}

		// Sanitize path
		fpath := filepath.Join(targetDir, filepath.Clean(hdr.Name))
		if !filepath.HasPrefix(fpath, filepath.Clean(targetDir)+string(os.PathSeparator)) {
			return NewError(CodeLayerExtract, "filesystem", "", fmt.Sprintf("illegal file path in layer: %s", hdr.Name), nil)
		}

		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(fpath, os.FileMode(hdr.Mode)); err != nil {
				return NewError(CodeLayerExtract, "filesystem", "", fmt.Sprintf("failed to create directory: %s", fpath), err)
			}
		case tar.TypeReg, tar.TypeRegA:
			if err := os.MkdirAll(filepath.Dir(fpath), 0755); err != nil {
				return NewError(CodeLayerExtract, "filesystem", "", fmt.Sprintf("failed to create parent dir for file: %s", fpath), err)
			}
			f, err := os.OpenFile(fpath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(hdr.Mode))
			if err != nil {
				return NewError(CodeLayerExtract, "filesystem", "", fmt.Sprintf("failed to create file: %s", fpath), err)
			}

			if _, err := io.Copy(f, tr); err != nil {
				f.Close()
				return NewError(CodeLayerExtract, "filesystem", "", fmt.Sprintf("failed to write file: %s", fpath), err)
			}
			f.Close()
		case tar.TypeSymlink:
			if err := os.Symlink(hdr.Linkname, fpath); err != nil {
				return NewError(CodeLayerExtract, "filesystem", "", fmt.Sprintf("failed to create symlink: %s", fpath), err)
			}
		default:
			// Ignora outros tipos (hardlinks, dispositivos, etc.)
		}
	}

	return nil
}

// ExtractAllLayersToFS writes all layers from an image to targetDir.
// Recebe os Layers da struct Image e um map de readers.
func ExtractAllLayersToFS(img *Image, layerReaders map[int]io.ReadCloser, targetDir string) error {
	if img == nil {
		return NewError(CodeBuildFailed, "filesystem", "", "image is nil", nil)
	}

	for _, layer := range img.Layers {
		rc, ok := layerReaders[layer.Index]
		if !ok {
			return NewError(CodeLayerExtract, "filesystem", img.Reference, fmt.Sprintf("no reader provided for layer %d", layer.Index), nil)
		}
		if err := ExtractLayerToFS(layer, rc, targetDir); err != nil {
			return err
		}
	}

	return nil
}
