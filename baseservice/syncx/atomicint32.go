package syncx

import "sync/atomic"

type AtomicInt32 int32

func NewAtomicInt32() *AtomicInt32 {
	return new(AtomicInt32)
}

func ForAtomicInt32(val int32) *AtomicInt32 {
	i := NewAtomicInt32()
	i.Set(val)
	return i
}

// 将当前值与给定的old值 进行比较，如果相等，则设置为给定的 val值
func (i *AtomicInt32) CompareAndSwap(old, val int32) bool {
	return atomic.CompareAndSwapInt32((*int32)(i), old, val)
}

// 修改
func (i *AtomicInt32) Set(val int32) {
	atomic.StoreInt32((*int32)(i), val)
}

// 修改
func (i *AtomicInt32) Add(val int32) int32 {
	return atomic.AddInt32((*int32)(i), val)
}

// 原子性修改
// 不知道为啥这么写 来自go-zero
func (i *AtomicInt32) AddAtomic(val int32) int32 {
	for {
		old := i.Val()
		nv := old + val
		if i.CompareAndSwap(old, nv) {
			return nv
		}
	}
}

// 替换 返回旧值
func (i *AtomicInt32) Swap(val int32) int32 {
	return atomic.SwapInt32((*int32)(i), val)
}

// 获取
func (i *AtomicInt32) Val() int32 {
	return atomic.LoadInt32((*int32)(i))
}

// 判断
func (i *AtomicInt32) IsVal(val int32) bool {
	return atomic.LoadInt32((*int32)(i)) == val
}
