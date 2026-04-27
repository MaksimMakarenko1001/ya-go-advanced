package functools

import (
	"context"
	"sync"
	"time"
)

type CallAsyncResult[T any] struct {
	Data T
	Err  error
}

func CallAsyncOrTimeout[Rq, Rs any](
	ctx context.Context,
	timeout time.Duration,
	request Rq,
	call func(ctx context.Context, rq Rq) (Rs, error),
) <-chan CallAsyncResult[Rs] {
	var once sync.Once
	ctx, cancel := context.WithTimeout(ctx, timeout)

	res := make(chan CallAsyncResult[Rs])
	writeChan := func(x CallAsyncResult[Rs]) {
		once.Do(func() {
			defer close(res)
			res <- x
		})
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		data, err := call(ctx, request)

		writeChan(CallAsyncResult[Rs]{Data: data, Err: err})
	}()

	go func() {
		defer cancel()

		select {
		case <-done:
		case <-ctx.Done():
			writeChan(CallAsyncResult[Rs]{Err: ctx.Err()})
		}

	}()

	return res
}
