// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This is an async adaptation of golang.org/x/sync/errgroup
package gogroup

import (
	"context"
	"sync"
)

// Group is a collection of goroutines serving as one a part of one process. The underlying context
// is cancelled as soon as any goroutine terminates regardless of the outcome.
type Group interface {
	// Go calls the given function in a new goroutine.
	Go(f func() error)
	// Wait waits until for the first function started with the Go method returns
	// or until the parent context is cancelled.
	Wait() error
}

type group struct {
	ctx    context.Context
	cancel func()

	errOnce sync.Once
	err     error
}

// New returns a new GoGroup and an associated Context derived from the provided one.
// The context is cancelled when the first goroutine started with Go terminates
// regardless of the outcome.
func New(ctx context.Context) (*group, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	return &group{ctx: ctx, cancel: cancel}, ctx
}

func (g *group) Go(f func() error) {
	go func() {
		err := f()
		g.errOnce.Do(func() {
			g.err = err
		})
		g.cancel()
	}()
}

func (g *group) Wait() error {
	<-g.ctx.Done()
	// we need to do this to avoid data race on err
	// in case the parent context is cancelled
	g.errOnce.Do(func() {
		g.err = nil
	})

	return g.err
}
