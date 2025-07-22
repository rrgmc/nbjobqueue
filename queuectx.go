package nbjobqueue

import "context"

type QueueCtx struct {
	ctx       context.Context
	ctxCancel context.CancelFunc
	queue     *Queue
}

func NewWithContext(ctx context.Context, concurrency int) *QueueCtx {
	ctx, cancel := context.WithCancel(context.Background())
	return &QueueCtx{
		ctx:       ctx,
		ctxCancel: cancel,
		queue:     New(concurrency),
	}
}

func (q *QueueCtx) Add(job JobCtx) {
	q.queue.AddJob(func() {
		job.Run(q.ctx)
	})
}

func (q *QueueCtx) AddCheck(job JobCtx) error {
	return q.queue.AddJobCheck(func() {
		job.Run(q.ctx)
	})
}

func (q *QueueCtx) AddJob(f func(context.Context)) {
	q.Add(JobCtxFunc(f))
}

func (q *QueueCtx) AddJobCheck(f func(context.Context)) error {
	return q.AddCheck(JobCtxFunc(f))
}

func (q *QueueCtx) Closed() bool {
	return q.queue.Closed()
}

func (q *QueueCtx) Stop() {
	q.queue.Stop()
}

func (q *QueueCtx) CancelAndClose() {
	if !q.queue.Closed() {
		q.ctxCancel()
	}
	q.queue.Close() // we are doing our own cancelling
}

func (q *QueueCtx) Close() {
	q.queue.Close()
	q.ctxCancel()
}
