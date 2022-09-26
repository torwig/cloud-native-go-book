package debounce

import (
	"context"
	"sync"
	"time"
)

type Circuit func(ctx context.Context) (string, error)

func First(circuit Circuit, d time.Duration) Circuit {
	var (
		threshold time.Time
		result    string
		err       error
		mtx       sync.Mutex
	)

	return func(ctx context.Context) (string, error) {
		mtx.Lock()

		defer func() {
			threshold = time.Now().Add(d)
			mtx.Unlock()
		}()

		if time.Now().Before(threshold) {
			return result, err
		}

		result, err = circuit(ctx)

		return result, err
	}
}

func Last(circuit Circuit, d time.Duration) Circuit {
	var (
		threshold = time.Now()
		ticker    *time.Ticker
		err       error
		result    string
		once      sync.Once
		mtx       sync.Mutex
	)
	return func(ctx context.Context) (string, error) {
		mtx.Lock()
		defer mtx.Unlock()

		threshold = time.Now().Add(d)

		once.Do(func() {
			ticker = time.NewTicker(100 * time.Millisecond)

			go func() {
				defer func() {
					mtx.Lock()
					ticker.Stop()
					once = sync.Once{}
					mtx.Unlock()
				}()

				for {
					select {
					case <-ticker.C:
						mtx.Lock()

						if time.Now().After(threshold) {
							result, err = circuit(ctx)
							mtx.Unlock()
							return
						}

						mtx.Unlock()
					case <-ctx.Done():
						mtx.Lock()
						result, err = "", ctx.Err()
						mtx.Unlock()
						return
					}
				}
			}()
		})

		return result, err
	}
}
