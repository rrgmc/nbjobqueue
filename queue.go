package nbjobqueue

import "github.com/rrgmc/nbchanlist"

var (
	ErrClosed = nbchanlist.ErrStopped
)

type Queue struct {
	queue *nbchanlist.Queue[Job]
}

func New() *Queue {
	return &Queue{
		queue: nbchanlist.NewQueue[Job](),
	}
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
