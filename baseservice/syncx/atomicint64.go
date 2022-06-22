package syncx

import "sync/atomic"

type AtomicInt64 int64

func NewAtomicInt64() *AtomicInt64 {
	return new(AtomicInt64)
}

func ForAtomicInt64(val int64) *AtomicInt64 {
	i := NewAtomicInt64()
	i.Set(val)
	return i
}

// 将当前值与给定的old值 进行比较，如果相等，则设置为给定的 val值
func (i *AtomicInt64) CompareAndSwap(old, val int64) bool {
	return atomic.CompareAndSwapInt64((*int64)(i), old, val)
}

// 修改
func (i *AtomicInt64) Set(val int64) {
	atomic.StoreInt64((*int64)(i), val)
}

// 修改
func (i *AtomicInt64) Add(val int64) int64 {
	return atomic.AddInt64((*int64)(i), val)
}

// 原子性修改
func (i *AtomicInt64) AddAtomic(val int64) int64 {
	for {
		old := i.Val()
		nv := old + val
		if i.CompareAndSwap(old, nv) {
			return nv
		}
	}
}

// 替换 返回旧值
func (i *AtomicInt64) Swap(val int64) int64 {
	return atomic.SwapInt64((*int64)(i), val)
}

// 获取
func (i *AtomicInt64) Val() int64 {
	return atomic.LoadInt64((*int64)(i))
}

// 判断
func (i *AtomicInt64) IsVal(val int64) bool {
	return atomic.LoadInt64((*int64)(i)) == val
}
