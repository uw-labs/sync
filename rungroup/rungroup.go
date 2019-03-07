// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This is an adaptation of golang.org/x/sync/errgroup
// that cancels the underlying context as soon as one of
// the goroutines started with the Go method terminates.
package rungroup

import (
	"context"
	"sync"
)

// Group is a collection of goroutines serving as one a part of one process. The underlying context
// is cancelled as soon as any goroutine terminates regardless of the outcome.
type Group interface {
	// Go calls the given function in a new goroutine.
	Go(f func() error)
	// Wait blocks until all function calls from the Go method have returned, then
	// returns the first non-nil error (if any) encountered.
	Wait() error
}

type group struct {
	cancel func()

	wg sync.WaitGroup

	errOnce sync.Once
	err     error
}

// New returns a new Group and an associated Context derived from the provided one.
// The context is cancelled when the first goroutine started with Go terminates
// regardless of the outcome.
func New(ctx context.Context) (*group, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	return &group{cancel: cancel}, ctx
}

func (g *group) Wait() error {
	g.wg.Wait()
	g.cancel()

	return g.err
}

func (g *group) Go(f func() error) {
	g.wg.Add(1)

	go func() {
		defer g.wg.Done()

		if err := f(); err != nil {
			g.errOnce.Do(func() {
				g.err = err
			})
		}
		g.cancel()
	}()
}
