package mr

import (
	"context"

	"github.com/zeromicro/go-zero/core/mr"
)

func Finish(fns ...func() error) error {
	return mr.Finish(fns...)
}

func FinishWithCtx(ctx context.Context, fns ...func() error) error {
	if len(fns) == 0 {
		return nil
	}

	return mr.MapReduceVoid(func(source chan<- interface{}) {
		for _, fn := range fns {
			source <- fn
		}
	}, func(item interface{}, writer mr.Writer, cancel func(error)) {
		fn := item.(func() error)
		if err := fn(); err != nil {
			cancel(err)
		}
	}, func(pipe <-chan interface{}, cancel func(error)) {
	}, mr.WithWorkers(len(fns)), mr.WithContext(ctx))
}
