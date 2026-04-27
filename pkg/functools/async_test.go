package functools

import (
	"context"
	"errors"
	"testing"
	"time"
)

type Request struct{}
type Response struct{}

func TestCallAsyncOrTimeout(t *testing.T) {
	timeout := 1 * time.Second
	callHappy := func(ctx context.Context, rq *Request) (*Response, error) {
		return &Response{}, nil
	}
	t.Run("happy path", func(t *testing.T) {
		data := <-CallAsyncOrTimeout(context.Background(), timeout, &Request{}, callHappy)

		if err := data.Err; err != nil {
			t.Errorf("got unexpected err: %v", err)
		}
	})

	callTimeout := func(ctx context.Context, rq *Request) (*Response, error) {
		time.Sleep(timeout + time.Millisecond)
		return &Response{}, nil
	}
	t.Run("timeout", func(t *testing.T) {

		data := <-CallAsyncOrTimeout(context.Background(), timeout, &Request{}, callTimeout)

		if err := data.Err; !errors.Is(err, context.DeadlineExceeded) {
			t.Errorf("got unexpected err: %v", err)
		}
	})
}
