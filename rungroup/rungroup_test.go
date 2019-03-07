package rungroup_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/uw-labs/sync/rungroup"
)

func run(ctx context.Context, d time.Duration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(d):
		return nil
	}
}

func TestGroup_Empty(t *testing.T) {
	assert := require.New(t)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	g, ctx := rungroup.New(ctx)
	assert.Nil(g.Wait())
}

func TestGroup_StopOnError(t *testing.T) {
	assert := require.New(t)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	var syncErr error
	expectedErr := errors.New("component stopped")
	g, ctx := rungroup.New(ctx)

	g.Go(func() error {
		syncErr = run(ctx, time.Second)
		return syncErr
	})
	g.Go(func() error {
		time.Sleep(time.Millisecond * 50)
		return expectedErr
	})

	assert.Equal(expectedErr, g.Wait())
	assert.Equal(context.Canceled, syncErr)
}

func TestGroup_StopOnTermination(t *testing.T) {
	assert := require.New(t)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	var sync1Err, sync2Err error
	g, ctx := rungroup.New(ctx)

	g.Go(func() error {
		sync2Err = run(ctx, time.Second)
		return sync2Err
	})
	g.Go(func() error {
		sync1Err = run(ctx, time.Millisecond*50)
		return sync1Err
	})

	assert.Equal(context.Canceled, g.Wait())
	assert.Equal(nil, sync1Err)
	assert.Equal(context.Canceled, sync2Err)
}
