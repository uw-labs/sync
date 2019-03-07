package gogroup_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/uw-labs/sync/gogroup"
)

func run(ctx context.Context, d time.Duration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(d):
		return nil
	}
}

func TestGroup_StopOnTermination(t *testing.T) {
	assert := require.New(t)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	g, ctx := gogroup.New(ctx)

	g.Go(func() error {
		return run(ctx, time.Second)
	})
	g.Go(func() error {
		time.Sleep(time.Millisecond * 50)
		return nil
	})

	assert.Nil(g.Wait())
}

func TestGroup_StopOnError(t *testing.T) {
	assert := require.New(t)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	expectedErr := errors.New("component stopped")
	g, ctx := gogroup.New(ctx)

	g.Go(func() error {
		return run(ctx, time.Second)
	})
	g.Go(func() error {
		time.Sleep(time.Millisecond * 50)
		return expectedErr
	})

	assert.Equal(expectedErr, g.Wait())
}

func TestGroup_ParentCancelled(t *testing.T) {
	assert := require.New(t)
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()

	g, ctx := gogroup.New(ctx)

	g.Go(func() error {
		return run(ctx, time.Second)
	})
	g.Go(func() error {
		return run(ctx, time.Second)
	})

	err := g.Wait()
	assert.True(err == nil || err == context.DeadlineExceeded)
}
