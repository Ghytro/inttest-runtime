package worker

import (
	"context"

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
	if err := w.acquire(ctx); err != nil {
		return err
	}
	go func() {
		fn()
		w.release()
	}()
	return nil
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
