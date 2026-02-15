package analyzer

import (
	"context"
	"time"
)

func retry(ctx context.Context, attempts int, baseDelay time.Duration, fn func() error) error {
	var err error

	for i := 0; i <= attempts; i++ {
		err = fn()
		if err == nil {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		time.Sleep(baseDelay * time.Duration(1<<i))
	}

	return err
}
