package Session

import (
	"sync"
	"sync/atomic"
)

type WaitGroup struct {
	sync.WaitGroup
	m int32
}

func (w *WaitGroup) Add(i int) {
	atomic.AddInt32(&w.m, int32(i))
	w.WaitGroup.Add(i)
}

func (w *WaitGroup) Done() {
	if atomic.LoadInt32(&w.m) < 1 {
		return
	}
	atomic.AddInt32(&w.m, -1)
	w.WaitGroup.Done()
}
func (w *WaitGroup) IsWait() bool {
	return atomic.LoadInt32(&w.m) != 0
}
