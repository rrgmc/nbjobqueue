package nbjobqueue

import (
	"cmp"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp/cmpopts"
	"gotest.tools/v3/assert"
)

func TestQueue(t *testing.T) {
	jq := New(3)

	var items []int
	var lock sync.Mutex

	for i := 0; i < 10; i++ {
		jq.AddJob(func() {
			lock.Lock()
			defer lock.Unlock()
			items = append(items, i)
		})
	}

	jq.Close()

	assert.Assert(t, jq.Closed())
	assert.DeepEqual(t, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, items, cmpopts.SortSlices(cmp.Less[int]))
}

func TestQueueCloseDrain(t *testing.T) {
	for _, drain := range []bool{false, true} {
		t.Run(fmt.Sprintf("drain=%t", drain), func(t *testing.T) {
			jq := New(3)

			var items []int
			var lock sync.Mutex

			for i := 0; i < 10; i++ {
				_ = jq.AddJobCheck(func() {
					time.Sleep(50 * time.Millisecond)
					lock.Lock()
					defer lock.Unlock()
					items = append(items, i)
				})
			}

			jq.CloseOpt(drain)

			if drain {
				assert.DeepEqual(t, []int{0, 1, 2}, items, cmpopts.SortSlices(cmp.Less[int]))
			} else {
				assert.DeepEqual(t, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, items, cmpopts.SortSlices(cmp.Less[int]))
			}
		})
	}
}

func TestQueueStop(t *testing.T) {
	jq := New(3)

	var items []int
	var lock sync.Mutex

	for i := 0; i < 10; i++ {
		jq.AddJob(func() {
			lock.Lock()
			defer lock.Unlock()
			items = append(items, i)
		})
	}

	jq.Stop()

	for i := 10; i < 20; i++ {
		err := jq.AddJobCheck(func() {
			lock.Lock()
			defer lock.Unlock()
			items = append(items, i)
		})
		assert.ErrorIs(t, err, ErrClosed)
	}

	jq.Close()

	assert.Assert(t, jq.Closed())
	assert.DeepEqual(t, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, items, cmpopts.SortSlices(cmp.Less[int]))
}

func TestQueueNegativeConcurrency(t *testing.T) {
	jq := New(-1)

	var items []int
	var lock sync.Mutex

	for i := 0; i < 10; i++ {
		jq.AddJob(func() {
			lock.Lock()
			defer lock.Unlock()
			items = append(items, i)
		})
	}

	jq.Close()

	assert.Assert(t, jq.Closed())
	assert.DeepEqual(t, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, items, cmpopts.SortSlices(cmp.Less[int]))
}
