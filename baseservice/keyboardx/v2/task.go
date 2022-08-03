package keyboardx

import "github.com/chuan-fu/Common/zlog"

type RunTaskFunc func(t Task, key string) (isEnd bool, err error)

type Task interface {
	PreHandle(s string) error
	Handle(s string) error
	PostHandle(s string) error
	IsEnd() bool
}

type BaseTask struct {
	IsEndField bool
}

func (b *BaseTask) PreHandle(s string) error {
	return nil
}

func (b *BaseTask) Handle(s string) error {
	return nil
}

func (b *BaseTask) PostHandle(s string) error {
	return nil
}

func (b *BaseTask) IsEnd() bool {
	return b.IsEndField
}

func runTask(t Task, key string) (isEnd bool, err error) {
	if t == nil {
		return
	}

	// 前置
	if err = t.PreHandle(key); err != nil {
		log.Error(err)
		return
	}
	if t.IsEnd() {
		return t.IsEnd(), nil
	}

	// 主体
	if err = t.Handle(key); err != nil {
		log.Error(err)
		return
	}
	if t.IsEnd() {
		return t.IsEnd(), nil
	}

	// 后置
	if err = t.PostHandle(key); err != nil {
		log.Error(err)
		return
	}
	return t.IsEnd(), nil
}
