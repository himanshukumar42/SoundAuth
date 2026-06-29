package lib

import (
	"context"
	"log"
	"time"
)

// Retry Pattern

func Retry(ctx context.Context, attempts int, delay time.Duration, fn func() error) error {
	var err error

	for i := 1; i <= attempts; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err = fn(); err == nil {
			return nil
		}

		log.Printf("[Retry] attempt=%d error=%v\n", i, err)
	}
	return err
}
