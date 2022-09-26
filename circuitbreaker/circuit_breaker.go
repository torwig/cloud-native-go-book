package circuitbreaker

import (
	"context"
	"errors"
	"sync"
	"time"
)

type Circuit func(ctx context.Context) (string, error)

func Breaker(circuit Circuit, failureThreshold int) Circuit {
	consecutiveFailures := 0
	lastAttempt := time.Now()
	var mtx sync.RWMutex

	return func(ctx context.Context) (string, error) {
		mtx.RLock()

		d := consecutiveFailures - failureThreshold

		if d >= 0 {
			retryAt := lastAttempt.Add(time.Second * 2 << d)
			if time.Now().After(retryAt) {
				mtx.RUnlock()
				return "", errors.New("service unreachable")
			}
		}

		mtx.RUnlock()

		response, err := circuit(ctx)

		mtx.Lock()
		defer mtx.Unlock()

		lastAttempt = time.Now()

		if err != nil {
			consecutiveFailures++
			return response, err
		}

		consecutiveFailures = 0

		return response, nil
	}
}
