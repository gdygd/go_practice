package worker

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

type Job interface {
	Run(ctx context.Context) error
}

var ErrPoolClosed = errors.New("worker pool closed")

type WorkerPool struct {
	count     int
	jobs      chan Job
	wg        sync.WaitGroup
	ctx       context.Context
	cancel    context.CancelFunc
	startOnce sync.Once
	stopOnce  sync.Once
	mu        sync.Mutex
	closed    bool
}

func NewPool(workerCount int, jobsBuffer int) *WorkerPool {
	if workerCount <= 0 {
		workerCount = 1
	}
	if jobsBuffer <= 0 {
		jobsBuffer = workerCount * 2
	}

	ctx, cancel := context.WithCancel(context.Background())
	return &WorkerPool{
		count:  workerCount,
		jobs:   make(chan Job, jobsBuffer),
		ctx:    ctx,
		cancel: cancel,
	}
}

func (w *WorkerPool) Start() {
	w.startOnce.Do(func() {
		for i := 0; i < w.count; i++ {
			w.wg.Add(1)
			go w.worker(i)
		}
	})
}

func (w *WorkerPool) Submit(job Job) error {
	w.mu.Lock()
	closed := w.closed
	w.mu.Unlock()

	if closed {
		return ErrPoolClosed
	}
	select {
	case <-w.ctx.Done():
		return ErrPoolClosed
	case w.jobs <- job:
		return nil
	}
}

func (w *WorkerPool) Stop(waitCtx context.Context) error {
	var retErr error
	w.stopOnce.Do(func() {
		w.mu.Lock()
		w.closed = true
		w.mu.Unlock()

		w.cancel()

		close(w.jobs)
		ch := make(chan struct{})
		go func() {
			w.wg.Wait()
			close(ch)
		}()

		select {
		case <-ch:
			// normal shutdown
		case <-waitCtx.Done():
			retErr = fmt.Errorf("stop timeout: %w", waitCtx.Err())
		}
	})
	return retErr
}

func (w *WorkerPool) worker(id int) {
	defer w.wg.Done()
	log.Printf("start woker .. %d", id)
	for {
		select {
		case <-w.ctx.Done():
			for job := range w.jobs {
				w.safeRun(job)
			}
			return
		case job, ok := <-w.jobs:
			if !ok {
				return
			}
			w.safeRun(job)
		}
	}
}

func (w *WorkerPool) safeRun(job Job) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("safeFun recover.. %v", r)
		}
	}()

	// ctx, _ := context.WithCancel(w.ctx)
	// _ = job(ctx)
	_ = job.Run(w.ctx)
}

func (w *WorkerPool) Active() bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	return !w.closed
}

func (w *WorkerPool) SubmitWithTimeout(job Job, timeout time.Duration) error {
	if timeout <= 0 {
		return w.Submit(job)
	}

	w.mu.Lock()
	closed := w.closed
	w.mu.Unlock()
	if closed {
		return ErrPoolClosed
	}
	select {
	case <-w.ctx.Done():
		return ErrPoolClosed
	case w.jobs <- job:
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("submit timeout after %s", timeout)
	}
}
