package concurrent

import (
	"github.com/gamenotover/mini_processor/ctry"
	"sync"
	"time"
)

type AsyncWorker struct {
	jobChan chan func()
	wg      sync.WaitGroup
	once    sync.Once
}

// NewAsyncWorker workers 并行任务数 buffer 待处理队列buffer长度
func NewAsyncWorker(workers, buffer int) *AsyncWorker {
	aw := &AsyncWorker{
		jobChan: make(chan func(), buffer),
	}
	for i := 0; i < workers; i++ {
		aw.wg.Add(1)
		go func() {
			defer aw.wg.Done()
			for {
				f, ok := <-aw.jobChan
				if !ok {
					return
				}
				// 只做最后兜底
				ctry.Try(f)
			}
		}()
	}
	return aw
}

// Add 添加任务 如果添加的超过任务buffer可能导致阻塞 如果已经调用过AddDone则会导致panic
func (aw *AsyncWorker) Add(job func()) {
	aw.jobChan <- job
}

// AddDone  结束添加任务, 务必调用一次否则导致goroutine泄露
func (aw *AsyncWorker) AddDone() {
	aw.once.Do(func() {
		close(aw.jobChan)
	})
}

// Wait 等待所有异步任务结束
func (aw *AsyncWorker) Wait() {
	aw.AddDone()
	aw.wg.Wait()
}

// ConditionChecker 每interval检查一个bool条件, 为true后关闭信号chan
func ConditionChecker(interval time.Duration, f func() bool) (exitChan chan struct{}) {

	exitChan = make(chan struct{})
	go func() {
		var stop bool
		for ; !stop; <-time.Tick(interval) {
			ctry.Try(func() {
				if stop = f(); stop {
					close(exitChan)
				}
			})
		}
	}()
	return
}
