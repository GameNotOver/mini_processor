package concurrent

import (
	"context"
	"github.com/gamenotover/mini_processor/cerr"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type syncMap struct {
	m sync.Map
}

func (s *syncMap) MustLoad(key interface{}) (value interface{}) {
	var ok bool
	if value, ok = s.Load(key); !ok {
		panic(errors.Errorf("key: %v not exists", key))
	} else {
		return
	}
}

func (s *syncMap) Load(key interface{}) (value interface{}, ok bool) {
	return s.m.Load(key)
}

func (s *syncMap) Store(key, value interface{}) {
	s.m.Store(key, value)
}

func (s *syncMap) Range(f func(key, value interface{})) {
	s.m.Range(func(key, value interface{}) bool {
		f(key, value)
		return true
	})
}

func (s *syncMap) Delete(key interface{}) {
	s.m.Delete(key)
}

type asyncFunctions struct {
	functions []func()
	async     *asyncController
}

func (a *asyncFunctions) Append(f func()) {
	a.functions = append(a.functions, f)
}

func (a *asyncFunctions) Clear() {
	a.functions = make([]func(), 0, 5)
}

func (a *asyncFunctions) Go(ctx context.Context) {
	a.async.Do(ctx, a.functions...)
}

func (a *asyncFunctions) GoWithLimit(ctx context.Context, limit int) {
	a.async.DoWithLimit(ctx, limit, a.functions...)
}

func (a *asyncFunctions) GoAndClear(ctx context.Context) {
	a.Go(ctx)
	a.Clear()
}

func (a *asyncFunctions) GoWithLimitAndClear(ctx context.Context, limit int) {
	a.GoWithLimit(ctx, limit)
	a.Clear()
}

type asyncController struct{}

func NewAsyncController() Concurrent {
	return &asyncController{}
}

func (a *asyncController) NewFunctions() AsyncFunctions {
	return &asyncFunctions{
		functions: make([]func(), 0, 5),
		async:     a,
	}
}

func (a *asyncController) NewMap() SyncMap {
	return new(syncMap)
}

func (a *asyncController) GoWithRecover(fun func(), recoverFn func(interface{})) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				recoverFn(err)
			}
		}()
		fun()
	}()
}

func (a *asyncController) Go(ctx context.Context, fun func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logrus.WithContext(ctx).Errorf("%v", err)
			}
		}()
		fun()
	}()
}

func (a *asyncController) SafeDo(ctx context.Context, fun func()) {
	defer func() {
		if err := recover(); err != nil {
			logrus.WithContext(ctx).Errorf("%v", err)
		}
	}()
	fun()
}

func (a *asyncController) Do(ctx context.Context, fns ...func()) {
	a.DoWithLimit(ctx, len(fns), fns...)
}
func (a *asyncController) DoWithTime(ctx context.Context, fn func(), duration time.Duration) {
	ch := make(chan int, 1)
	go func() {
		fn()
		ch <- 1
	}()
	select {
	case <-ch:
		return
	case <-time.After(duration):
		logrus.WithContext(ctx).Infof("execute fn timeout:%v", duration)
	}
}

func (a *asyncController) DoWithLimit(ctx context.Context, limit int, fns ...func()) {
	collect := newPanicCollector(len(fns), limit)
	for _, fn := range fns {
		currentFn := fn
		collect.Go(ctx, func() { currentFn() })
	}
	collect.Done()
}

func newPanicCollector(errLen int, limit int) *panicCollector {
	return &panicCollector{
		l:         errLen,
		limitChan: make(chan struct{}, limit),
		errChan:   make(chan error, errLen),
	}
}

type panicCollector struct {
	l         int
	limitChan chan struct{}
	errChan   chan error
}

func (c *panicCollector) Go(ctx context.Context, fn func()) {
	c.limitChan <- struct{}{}
	go func() {
		defer func() {
			if err := recover(); err != nil {
				if bizerr, ok := err.(cerr.BizError); ok {
					c.errChan <- bizerr
				} else {
					c.errChan <- errors.Errorf("[OKR][SafeGoWithCtx] Goroutine Recover: %v", err)
				}
			} else {
				c.errChan <- nil
			}
			<-c.limitChan
		}()
		fn()
	}()
}

func (c *panicCollector) Done() {
	defer close(c.errChan)
	defer close(c.limitChan)
	var finalErr error
	for i := 0; i < c.l; i++ {
		if err := <-c.errChan; err != nil {
			finalErr = err
		}
	}
	if finalErr != nil {
		panic(finalErr)
	}
}
