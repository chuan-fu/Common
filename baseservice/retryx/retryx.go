package retryx

import (
	"context"

	"github.com/chuan-fu/Common/zlog"
)

/*
	重试
*/

type (
	RetryTask            func(ctx context.Context) (interface{}, error)
	RetryTaskWithIsRetry func(ctx context.Context) (data interface{}, isRetry bool, err error)
)

func Retry(ctx context.Context, task RetryTask, opts ...Option) (data interface{}, err error) {
	r := buildConfig(opts)

	for i := 0; i < r.retryNum; i++ {
		data, err = task(ctx)
		if err != nil {
			log.Error(err)
			continue
		}
		if !r.isRetry(data) {
			break
		}
	}
	return
}

func RetryFunc(ctx context.Context, retryNum int, task RetryTaskWithIsRetry) (data interface{}, err error) {
	if retryNum < 1 {
		retryNum = 1
	}
	var isRetry bool
	for i := 0; i < retryNum; i++ {
		data, isRetry, err = task(ctx)
		if err != nil {
			log.Error(err)
			continue
		}
		if !isRetry {
			break
		}
	}
	return
}
