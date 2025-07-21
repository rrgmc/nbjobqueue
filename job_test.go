package nbjobqueue

import (
	"fmt"
	"testing"
)

func TestQueue(t *testing.T) {
	jq := New(3)
	jq.AddJob(func() {
		fmt.Println("job 1")
	})
}
