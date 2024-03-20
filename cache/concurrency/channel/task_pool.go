package channel

import (
	"context"
	"go.uber.org/atomic"
	"sync"
)

type Task func()

type TaskPool struct {
	tasks chan Task
	close *atomic.Bool
}

func NewTaskPool(numG int, capacity int) *TaskPool {
	res := &TaskPool{
		tasks: make(chan Task, capacity),
		close: atomic.NewBool(false),
	}

	for i := 0; i < numG; i++ {
		go func() {
			for t := range res.tasks {
				if res.close.Load() {
					return
				}
				t()
			}
		}()
	}

	return res
}

func (p *TaskPool) Close() error {
	p.close.Store(true)
	return nil
}

func (p *TaskPool) Submit(ctx context.Context, t Task) error {
	select {
	case p.tasks <- t:
	case <-ctx.Done():
		return ctx.Err()
	}
	return nil
}

type TaskPoolV2 struct {
	tasks     chan Task
	close     chan struct{}
	closeOnce sync.Once
}

func NewTaskPoolV2(numG int, capacity int) *TaskPoolV2 {
	res := &TaskPoolV2{
		tasks: make(chan Task, capacity),
		//close: atomic.NewBool(false),
	}

	for i := 0; i < numG; i++ {
		go func() {
			for {
				select {
				case <-res.close:
					return
				case t := <-res.tasks:
					t()
				}
			}
		}()
	}

	return res
}

func (p *TaskPoolV2) Close() error {
	//p.close <- struct {}{}             // 只会发一个，会有问题
	close(p.close) // 重复调用close会panic
	//p.closeOnce.Do(func() {
	//	close(p.close)
	//})
	return nil
}
