package backoff

import (
	"context"
	"fmt"
	"log"
	"time"
)

type ErrorClassification int
type retried func(ctx context.Context) (err error)

const (
	// NonRetriable - операцию не следует повторять
	NonRetriable ErrorClassification = iota

	// Retriable - операцию можно повторить
	Retriable
)

type Backoff struct {
	maxRetries      uint16
	errClassifyFunc func(err error) ErrorClassification
}

func NewBackoff(
	maxRetries uint16,
	errClassifyFunc func(err error) ErrorClassification,
) *Backoff {
	return &Backoff{
		maxRetries:      maxRetries,
		errClassifyFunc: errClassifyFunc,
	}
}

func (r *Backoff) WithLinear(t0 time.Duration, dt time.Duration) func(retried) retried {
	return func(fn retried) retried {
		return func(ctx context.Context) error {
			delay := t0
			for attempt := range r.maxRetries + 1 {
				select {
				case <-ctx.Done():
					return ctx.Err()
				default:
					err := fn(ctx)
					if err == nil {
						return nil
					}

					if r.errClassifyFunc(err) == NonRetriable {
						return fmt.Errorf("non retriable, %w", err)
					}

					log.Printf("attempt #%d failed: %v", attempt+1, err)
					if attempt < r.maxRetries {
						log.Printf("retrying in %vs...", delay.Seconds())
						time.Sleep(delay)
					} else {
						return fmt.Errorf("max attempts reached, %w", err)
					}
				}
				delay += dt
			}
			return nil
		}
	}
}
