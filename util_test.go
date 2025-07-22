package nbjobqueue

import (
	"context"
	"errors"
	"testing"
	"time"

	"gotest.tools/v3/assert"
)

func TestMergeContextCancel(t *testing.T) {
	ctx1, cancel1 := context.WithCancelCause(context.Background())
	defer cancel1(errors.New("ctx1 canceled"))

	ctx2, cancel2 := context.WithCancelCause(context.Background())

	mergedCtx, mergedCancel := MergeContextCancel(ctx1, ctx2)
	defer mergedCancel()

	testError := errors.New("ctx2 canceled")

	cancel2(testError)

	select {
	case <-mergedCtx.Done():
		assert.ErrorIs(t, context.Cause(mergedCtx), testError)
	case <-time.After(200 * time.Millisecond):
		assert.Assert(t, true, "context should have been canceled")
	}
}
