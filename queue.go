package nbjobqueue

import (
	"sync/atomic"

	"github.com/rrgmc/nbchanlist"
)

var (
	ErrClosed = nbchanlist.ErrClosed
)

type Queue struct {
	queue   *nbchanlist.Queue[Job]
	handler *handler
	closed  atomic.Bool
}

func New(concurrency int) *Queue {
	if concurrency < 1 {
		concurrency = 1
	}
	ret := &Queue{
		queue: nbchanlist.NewQueue[Job](),
	}
	ret.handler = newHandler(ret.queue, concurrency)
	return ret
}

func (q *Queue) Add(job Job) {
	q.queue.Put(job)
}

func (q *Queue) AddCheck(job Job) error {
	return q.queue.PutCheck(job)
}

func (q *Queue) AddJob(f func()) {
	q.Add(JobFunc(f))
}

func (q *Queue) AddJobCheck(f func()) error {
	return q.AddCheck(JobFunc(f))
}

// Closed returns whether Close was called.
func (q *Queue) Closed() bool {
	return q.closed.Load()
}

// Close stops accepting new jobs. Existing jobs will still run.
func (q *Queue) Close() {
	q.queue.Close()
}

// ShutdownOpt stops accepting new jobs, and waits until all existing jobs finish.
// If drain is true, clean the list of pending jobs before waiting.
func (q *Queue) ShutdownOpt(drain bool) {
	q.shutdown(drain, false, nil)
}

// Shutdown stops accepting new jobs and waits until all existing jobs finish.
func (q *Queue) Shutdown() {
	q.shutdown(false, false, nil)
}

func (q *Queue) shutdown(drain bool, cancel bool, stoppedCB func()) {
	if q.closed.CompareAndSwap(false, true) {
		if drain {
			q.queue.Shutdown() // stop accepting new jobs and drain remaining jobs
		} else {
			q.queue.Close() // stop accepting new jobs
		}
		if cancel {
			if stoppedCB != nil {
				stoppedCB()
			}
		}

		q.handler.stop() // wait for all job handler goroutines to finish

		if !cancel {
			if stoppedCB != nil {
				stoppedCB()
			}
		}
		if !drain {
			q.queue.Shutdown() // stop queue processing and drain remaining jobs
		}
	}
}
