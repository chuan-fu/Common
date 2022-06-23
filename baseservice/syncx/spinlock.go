package syncx

import (
	"runtime"
	"sync/atomic"
)

// SpinLock 用作快速执行的锁
type SpinLock struct {
	lock uint32
}

// 加锁
func (sl *SpinLock) Lock() {
	for !sl.TryLock() {
		runtime.Gosched()
	}
}

// 尝试加锁
func (sl *SpinLock) TryLock() bool {
	return atomic.CompareAndSwapUint32(&sl.lock, 0, 1)
}

// 解锁
func (sl *SpinLock) Unlock() {
	atomic.StoreUint32(&sl.lock, 0)
}
