package worker

import (
	"context"

	"github.com/samber/lo"
	"golang.org/x/sync/semaphore"
)

type BoundedWorker struct {
	sema *semaphore.Weighted
}

func NewBoundedWorker(maxWorkers int) *BoundedWorker {
	return &BoundedWorker{
		sema: semaphore.NewWeighted(int64(maxWorkers)),
	}
}

func (w *BoundedWorker) RunAsync(ctx context.Context, fn func()) error {
	errs := make(chan error, 1)
	go func() {
		if err := w.acquire(ctx); err != nil {
			errs <- err
			return
		}
		close(errs)
		fn()
		w.release()
	}()
	err, ok := <-errs
	return lo.Ternary(ok, err, nil)
}

func (w *BoundedWorker) Run(ctx context.Context, fn func()) error {
	if err := w.acquire(ctx); err != nil {
		return err
	}
	defer w.release()
	fn()
	return nil
}

func (w *BoundedWorker) acquire(ctx context.Context) error {
	return w.sema.Acquire(ctx, 1)
}

func (w *BoundedWorker) release() {
	w.sema.Release(1)
}
