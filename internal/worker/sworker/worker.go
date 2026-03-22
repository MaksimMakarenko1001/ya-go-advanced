package sworker

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/MaksimMakarenko1001/ya-go-advanced/internal/worker"
)

type SimpleWorker struct {
	config Config
	pid    string
	job    worker.Job
}

func New(config Config, pid string, job worker.Job) *SimpleWorker {
	return &SimpleWorker{
		config: config,
		pid:    pid,
		job:    job,
	}
}

func (w *SimpleWorker) Start(ctx context.Context) {
	go w.run(ctx)
}

func (w *SimpleWorker) run(ctx context.Context) {
	log.Printf("worker running, pid=%v", w.pid)

	onceCh := make(chan struct{}, 1)
	defer close(onceCh)

	for {
		select {
		case <-ctx.Done():
			return
		case onceCh <- struct{}{}:
			jobCtx, cancel := context.WithTimeout(ctx, w.config.JobTimeout)
			go func() {
				ts := time.Now()
				defer cancel()

				if err := w.doJob(jobCtx); err != nil {
					log.Printf("worker failure {pid=%v, err=%v}\n", w.pid, err.Error())
				}
				time.Sleep(w.config.JobInterval - time.Since(ts))
				<-onceCh
			}()
		}
	}
}

func (w *SimpleWorker) doJob(ctx context.Context) (err error) {
	err = w.job(ctx)
	if err != nil {
		return fmt.Errorf("job err: %w", err)
	}

	return nil
}
