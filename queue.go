package nbjobqueue

import (
	"sync/atomic"

	"github.com/rrgmc/nbchanlist"
)

var (
	ErrClosed = nbchanlist.ErrStopped
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

func (q *Queue) AddJob(f func()) {
	q.Add(JobFunc(f))
}

func (q *Queue) Stop() {
	q.queue.Stop()
}

func (q *Queue) CancelAndClose() {
	q.close(true)
}

func (q *Queue) Close() {
	q.close(false)
}

func (q *Queue) close(cancel bool) {
	if q.closed.CompareAndSwap(false, true) {
		q.queue.Stop()
		if cancel {
			q.handler.cancel()
		}
		q.handler.stop()
		q.queue.Close()
	}
}
