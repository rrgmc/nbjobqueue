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

// CloseOpt stops accepting new jobs, and waits until all existing jobs finish.
// If drain is true, the list of pending jobs is cleared before waiting.
// If cancel is true, the context is canceled before waiting.
func (q *QueueCtx) CloseOpt(drain bool, cancel bool) {
	q.close(drain, cancel)
}

// Close stops accepting new jobs, cancels the context, and waits until all existing jobs finish.
func (q *QueueCtx) Close() {
	q.close(false, true)
}

func (q *QueueCtx) close(drain bool, cancel bool) {
	q.queue.close(drain, cancel, func() {
		q.ctxCancel()
	})
}
