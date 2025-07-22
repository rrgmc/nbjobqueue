package nbjobqueue

import (
	"cmp"
	"context"
	"sync"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp/cmpopts"
	"gotest.tools/v3/assert"
)

func TestQueueCtx(t *testing.T) {
	jq := NewWithContext(context.Background(), 3)

	var items []int
	var lock sync.Mutex

	for i := 0; i < 10; i++ {
		jq.AddJob(func(jobCtx context.Context) {
			lock.Lock()
			defer lock.Unlock()
			items = append(items, i)
		})
	}

	jq.Close()

	assert.DeepEqual(t, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, items, cmpopts.SortSlices(cmp.Less[int]))
}

func TestQueueCtxCancel(t *testing.T) {
	jq := NewWithContext(context.Background(), 3)

	var items []int
	var lock sync.Mutex

	for i := 0; i < 10; i++ {
		jq.AddJob(func(jobCtx context.Context) {
			if i%2 == 0 {
				time.Sleep(100 * time.Millisecond)
				select {
				case <-jobCtx.Done():
					return
				default:
				}
			}
			lock.Lock()
			defer lock.Unlock()
			items = append(items, i)
		})
	}

	jq.CloseOpt(false, true)

	assert.DeepEqual(t, []int{1, 3, 5, 7, 9}, items, cmpopts.SortSlices(cmp.Less[int]))
}
