package worker

import "context"

type Job func(ctx context.Context) (err error)
