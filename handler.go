package nbjobqueue

import (
	"sync"

	"github.com/rrgmc/nbchanlist"
)

type handler struct {
	queue       *nbchanlist.Queue[Job]
	concurrency int
	wg          sync.WaitGroup
}

func newHandler(queue *nbchanlist.Queue[Job], concurrency int) *handler {
	ret := &handler{
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

func (h *handler) stop() {
	h.wg.Wait()
}
