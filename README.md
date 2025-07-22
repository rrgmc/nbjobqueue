# nbjobqueue - Non-blocking unbounded lock-free job queue for Golang
[![GoDoc](https://godoc.org/github.com/rrgmc/nbjobqueue?status.png)](https://godoc.org/github.com/rrgmc/nbjobqueue)

`nbjobqueue` is a non-blocking unbounded lock-free job queue for Golang.

As it is unbounded, care must be taken to avoid memory exhaustion if adding faster than reading. 

```go
import (
    "context"
    "fmt"
    "slices"
    "sync"

    "github.com/rrgmc/nbjobqueue"
)

func ExampleQueue() {
    jq := nbjobqueue.New(4) // start a job queue with 4 goroutines.

    var items []int
    var lock sync.Mutex

    // Add 10 jobs. At most, 4 jobs will run at the same time.
    // AddJob NEVER blocks.
    for i := 0; i < 10; i++ {
        jq.AddJob(func() {
            lock.Lock()
            defer lock.Unlock()
            items = append(items, i)
        })
    }

    jq.Shutdown() // wait for all jobs to be done, then release all resources.

    slices.Sort(items)
    fmt.Println(items)

    // Output: [0 1 2 3 4 5 6 7 8 9]
}
```

```go
import (
    "context"
    "fmt"
    "slices"
    "sync"

    "github.com/rrgmc/nbjobqueue"
)

func ExampleQueueCtx() {
    jq := nbjobqueue.NewWithContext(context.Background(), 4) // start a job queue with 4 goroutines.

    var items []int
    var lock sync.Mutex

    // Add 10 jobs. At most, 4 jobs will run at the same time.
    // AddJob NEVER blocks.
    for i := 0; i < 10; i++ {
        jq.AddJob(func(ctx context.Context) {
            lock.Lock()
            defer lock.Unlock()
            items = append(items, i)
        })
    }

    jq.ShutdownOpt(false, true) // cancel the context sent to jobs, but still calls all pending jobs.

    slices.Sort(items)
    fmt.Println(items)

    // Output: [0 1 2 3 4 5 6 7 8 9]
}
```

## Install

```shell
go get github.com/rrgmc/nbjobqueue
```

# License

MIT

### Author

Rangel Reale (rangelreale@gmail.com)
