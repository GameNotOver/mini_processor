package concurrent

import (
	"context"
	"time"
)

type Concurrent interface {
	GoWithRecover(fun func(), recoverFn func(interface{}))
	Go(ctx context.Context, fn func())
	Do(ctx context.Context, fns ...func())
	SafeDo(ctx context.Context, fun func())
	DoWithLimit(ctx context.Context, limit int, fns ...func())
	NewMap() SyncMap
	NewFunctions() AsyncFunctions
	DoWithTime(ctx context.Context, fn func(), duration time.Duration)
}

type SyncMap interface {
	MustLoad(key interface{}) (value interface{})
	Load(key interface{}) (value interface{}, ok bool)
	Store(key, value interface{})
	Range(f func(key, value interface{}))
	Delete(key interface{})
}

type AsyncFunctions interface {
	Go(ctx context.Context)
	GoWithLimit(ctx context.Context, limit int)
	GoAndClear(ctx context.Context)
	GoWithLimitAndClear(ctx context.Context, limit int)
	Append(f func())
	Clear()
}
