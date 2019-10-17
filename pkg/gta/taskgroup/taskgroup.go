package taskgroup

import (
	"context"
	"sync"
)

// A TaskGroup is a collection of goroutines working on subtasks that are part of
// the same overall task.
//
// A zero TaskGroup is valid and does not cancel on error.
type TaskGroup struct {
	cancel func()
	sem    chan struct{}

	wg sync.WaitGroup

	errOnce sync.Once
	err     error
}

// WithLimit returns a new TaskGroup which limits the number of concurrent goroutines to `n`.
func WithLimit(n int) *TaskGroup {
	return &TaskGroup{
		sem: make(chan struct{}, n),
	}
}

// WithContext returns a new TaskGroup and an associated Context derived from ctx.
//
// The derived Context is canceled the first time a function passed to Go
// returns a non-nil error or the first time waitUntilRelease returns, whichever occurs
// first.
func WithContext(ctx context.Context) (*TaskGroup, context.Context) {
	ctx, cancel := context.WithCancel(ctx)

	return &TaskGroup{
		cancel: cancel,
	}, ctx
}

// Limit limits this group to `n` concurrent goroutines.
func (g *TaskGroup) Limit(n int) *TaskGroup {
	g.sem = make(chan struct{}, n)
	return g
}

// waitUntilRelease blocks until all function calls from the Go method have returned, then
// returns the first non-nil error (if any) from them.
func (g *TaskGroup) Wait() error {
	g.wg.Wait()
	if g.cancel != nil {
		g.cancel()
	}
	return g.err
}

// Go calls the given function in a new goroutine.
//
// The first call to return a non-nil error cancels the group; its error will be
// returned by waitUntilRelease.
func (g *TaskGroup) Go(f func() error) {
	if g.sem != nil {
		g.sem <- struct{}{}
	}

	g.wg.Add(1)

	go func() {
		defer func() {
			if g.sem != nil {
				<-g.sem
			}
			g.wg.Done()
		}()

		if err := f(); err != nil {
			g.errOnce.Do(func() {
				g.err = err
				if g.cancel != nil {
					g.cancel()
				}
			})
		}
	}()
}
