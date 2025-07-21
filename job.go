package nbjobqueue

import "context"

type Job interface {
	Run()
}

type JobFunc func()

func (f JobFunc) Run() {
	f()
}

type JobCtx interface {
	Run(jobCtx context.Context)
}

type JobCtxFunc func(ctx context.Context)

func (f JobCtxFunc) Run(ctx context.Context) {
	f(ctx)
}
