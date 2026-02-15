package analyzer

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

func fetchImage(ctx context.Context, ref string, opts *options) (v1.Image, error) {
	parsedRef, err := name.ParseReference(ref)
	if err != nil {
		return nil, ErrInvalidReference
	}

	var img v1.Image

	err = retry(ctx, opts.retries, 500*time.Millisecond, func() error {
		var e error
		img, e = remote.Image(parsedRef, remote.WithContext(ctx))
		return e
	})

	if err != nil {
		return nil, fmt.Errorf("fetch failed: %w", err)
	}

	return img, nil
}
