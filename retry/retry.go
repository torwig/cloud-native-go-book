package retry

import (
	"context"
	"time"
)

type Effector func(ctx context.Context) (string, error)

func Retry(effector Effector, retries int, d time.Duration) Effector {
	return func(ctx context.Context) (string, error) {
		for r := 0; ; r++ {
			response, err := effector(ctx)
			if err == nil || r > retries {
				return response, err
			}

			select {
			case <-time.After(d):
			case <-ctx.Done():
				return "", ctx.Err()
			}
		}
	}
}
