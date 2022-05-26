package ants

import (
	"sync"

	"github.com/chuan-fu/Common/util"
	"github.com/chuan-fu/Common/zlog"
	"github.com/panjf2000/ants/v2"
	"github.com/pkg/errors"
)

// ants.WithNonblocking(true) 和 ants.WithMaxBlockingTasks(waitTask)重复
// waitTask为0时，可无限等待
// 小于0时，等同于WithNonblocking(true)
// 大于0时，可等待waitTask个任务，其余返回失败
// ants.WithPreAlloc(true)为NewPool时预分配内存，为true不支持缩扩容pool

var (
	globalAnts *ants.Pool
	globalOnce sync.Once
)

// 全局ants尽量设一个足够大的值，达到极限内存使用量
// 并阻止达到极限后的任务等待，防止服务崩溃，同时做好失败处理
func NewGlobalPool(size int) {
	globalOnce.Do(func() {
		globalAnts, _ = ants.NewPool(size,
			// 协程池耗尽后不等待直接报错
			ants.WithNonblocking(true),
			// panic处理
			ants.WithPanicHandler(util.DeferFuncLog),
		)
	})
}

type Task func() (err error)

// 使用
func Go(f Task) error {
	if globalAnts == nil {
		return errors.New("globalAnts is nil")
	}
	if err := globalAnts.Submit(func() {
		if e := f(); e != nil {
			log.Errorf("globalAnts Task error : %s", e.Error())
		}
	}); err != nil {
		return errors.Wrap(err, "globalAnts Submit error")
	}
	return nil
}

func GoVoid(f Task) {
	if err := Go(f); err != nil {
		log.Error(err)
	}
}

// 使用完毕记得 Release() or ReleaseTimeout()
func NewLocalAnts(size int, opts ...ants.Option) (*ants.Pool, error) {
	if len(opts) == 0 {
		return ants.NewPool(size, ants.WithPanicHandler(util.DeferFuncLog))
	}
	optsNew := make([]ants.Option, 1, len(opts)+1)
	optsNew[0] = ants.WithPanicHandler(util.DeferFuncLog) // panic处理
	optsNew = append(optsNew, opts...)
	return ants.NewPool(size, optsNew...)
}
