package nbjobqueue

import (
	"context"
	"sync"

	"github.com/rrgmc/nbchanlist"
)

type handler struct {
	ctx         context.Context
	ctxCancel   context.CancelFunc
	queue       *nbchanlist.Queue[Job]
	concurrency int
	wg          sync.WaitGroup
}

func newHandler(queue *nbchanlist.Queue[Job], concurrency int) *handler {
	ctx, cancel := context.WithCancel(context.Background())
	ret := &handler{
		ctx:         ctx,
		ctxCancel:   cancel,
		queue:       queue,
		concurrency: concurrency,
	}
	ret.start()
	return ret
}

func (h *handler) start() {
	for i := 0; i < h.concurrency; i++ {
		h.wg.Add(1)
		go func() {
			defer h.wg.Done()
			for {
				select {
				case <-h.ctx.Done():
					return
				case v, ok := <-h.queue.Get():
					if !ok {
						return
					}
					v.Run()
				}
			}
		}()
	}
}

func (h *handler) cancel() {
	h.ctxCancel()
}

func (h *handler) stop() {
	h.wg.Wait()
}
