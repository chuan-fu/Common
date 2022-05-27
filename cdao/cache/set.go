package cache

import (
	"context"

	"github.com/chuan-fu/Common/cdao"
	"github.com/chuan-fu/Common/zlog"
	"github.com/pkg/errors"
)

type (
	GetSetByCacheFunc func(ctx context.Context, b cdao.BaseRedisOp) ([]string, error)
	SetSetCacheFunc   func(ctx context.Context, b cdao.BaseRedisOp, v []string) error
	GetSetByDBFunc    func(ctx context.Context) ([]string, error)

	BaseSetCacheOption func(*BaseSetCacheOptions)
)

type BaseSetCacheOptions struct {
	GetByCache GetSetByCacheFunc
	SetCache   SetSetCacheFunc
}

// set存储，cache使用覆盖式写入防止重复
// 如getByDb为空，cache会写入一个空字符串
func GetBaseSetCache(ctx context.Context, op cdao.BaseRedisOp, getByDb GetSetByDBFunc, opts ...BaseSetCacheOption) (resp []string, err error) {
	b := &BaseSetCacheOptions{
		GetByCache: defaultGetSetByCache,
		SetCache:   defaultSetSetCache,
	}
	for _, opt := range opts {
		opt(b)
	}

	resp, err = b.GetByCache(ctx, op) // 获取缓存，并返回
	if err != nil {
		log.Error(errors.Wrap(err, "GetByCache"))
	}
	if len(resp) > 0 {
		if len(resp) == 1 && resp[0] == "" { // 只有一条空数据，代表为空
			resp = []string{}
		}
		return resp, nil
	}
	resp, err = getByDb(ctx) // 从db获取
	if err != nil {
		log.Error(errors.Wrap(err, "GetByDB"))
		return
	}

	// 写入cache
	if len(resp) == 0 {
		err = b.SetCache(ctx, op, []string{""}) // 写入空数据
	} else {
		err = b.SetCache(ctx, op, resp)
	}
	if err != nil {
		log.Error(errors.Wrap(err, "SetCache"))
		err = nil
	}

	return
}

func defaultGetSetByCache(ctx context.Context, b cdao.BaseRedisOp) ([]string, error) {
	return b.SGetAll(ctx)
}

func defaultSetSetCache(ctx context.Context, b cdao.BaseRedisOp, vs []string) error {
	return b.SAddCover(ctx, vs)
}

func WithGetSetByCache(fn GetSetByCacheFunc) BaseSetCacheOption {
	return func(opts *BaseSetCacheOptions) {
		opts.GetByCache = fn
	}
}

func WithSetSetCache(fn SetSetCacheFunc) BaseSetCacheOption {
	return func(opts *BaseSetCacheOptions) {
		opts.SetCache = fn
	}
}
