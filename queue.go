package nbjobqueue

import "github.com/rrgmc/nbchanlist"

var (
	ErrClosed = nbchanlist.ErrStopped
)

type Queue struct {
	queue   *nbchanlist.Queue[Job]
	handler *handler
}

func New(concurrency int) *Queue {
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

func (q *Queue) Close() {
	q.queue.Close()
}
