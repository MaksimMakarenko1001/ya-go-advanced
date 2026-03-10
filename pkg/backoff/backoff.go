// Package backoff provides retry mechanisms for operations with configurable error classification.
package backoff

import (
	"context"
	"fmt"
	"log"
	"time"
)

// ErrorClassification defines whether an error is retriable.
type ErrorClassification int

// retried represents a function that can be retried.
type retried func(ctx context.Context) (err error)

const (
	// NonRetriable indicates the operation should not be retried.
	NonRetriable ErrorClassification = iota

	// Retriable indicates the operation can be retried.
	Retriable
)

// Backoff manages retry logic with configurable error classification.
type Backoff struct {
	maxRetries      uint16
	errClassifyFunc func(err error) ErrorClassification
}

// NewBackoff creates a new Backoff instance with the specified maximum retries and error classification function.
func NewBackoff(
	maxRetries uint16,
	errClassifyFunc func(err error) ErrorClassification,
) *Backoff {
	return &Backoff{
		maxRetries:      maxRetries,
		errClassifyFunc: errClassifyFunc,
	}
}

// WithLinear returns a decorator that retries the function with linear backoff.
// The initial delay is t0 and increases by dt on each retry.
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
