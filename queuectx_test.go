package nbjobqueue

import (
	"cmp"
	"context"
	"sync"
	"testing"

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
