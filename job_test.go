package nbjobqueue

import (
	"cmp"
	"sync"
	"testing"

	"github.com/google/go-cmp/cmp/cmpopts"
	"gotest.tools/v3/assert"
)

func TestQueue(t *testing.T) {
	jq := New(3)

	var items []int
	var lock sync.Mutex

	for i := 0; i < 3; i++ {
		jq.AddJob(func() {
			lock.Lock()
			defer lock.Unlock()
			items = append(items, i)
		})
	}

	jq.Close()

	assert.DeepEqual(t, []int{0, 1, 2}, items, cmpopts.SortSlices(cmp.Less[int]))
}
