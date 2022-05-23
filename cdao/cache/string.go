package cache

import (
	"context"

	"github.com/chuan-fu/Common/cdao"
	"github.com/chuan-fu/Common/zlog"
	"github.com/pkg/errors"
)

/*
func defaultGetByDB(db *gorm.DB, id int64) GetByDBFunc {
	return func() (data string, err error) {
		return "11", nil
	}
}
*/

type (
	GetStringByCacheFunc func(ctx context.Context, b cdao.BaseRedisOp) (string, error)
	SetStringCacheFunc   func(ctx context.Context, b cdao.BaseRedisOp, v string) error
	GetStringByDBFunc    func(ctx context.Context) (string, error)

	BaseStringCacheOption func(*BaseStringCacheOptions)
)

func GetBaseStringCache(ctx context.Context, op cdao.BaseRedisOp, getByDb GetStringByDBFunc, opts ...BaseStringCacheOption) (data string, err error) {
	b := &BaseStringCacheOptions{
		GetByCache: defaultGetStringByCache,
		SetCache:   defaultSetStringCache,
	}
	for _, opt := range opts {
		opt(b)
	}

	data, err = b.GetByCache(ctx, op)
	if err == nil && data != "" {
		return
	}
	if err != nil {
		log.Error(errors.Wrap(err, "GetByCache"))
	}

	data, err = getByDb(ctx) // 从db获取
	if err != nil {
		log.Error(errors.Wrap(err, "GetByDB"))
		return
	}
	if data == "" {
		return "", errors.New("data为空，数据获取有误")
	}

	if err2 := b.SetCache(ctx, op, data); err2 != nil { // 写入cache
		log.Error(errors.Wrap(err2, "SetCache"))
	}
	return
}

func defaultGetStringByCache(ctx context.Context, b cdao.BaseRedisOp) (string, error) {
	return b.Get(ctx)
}

func defaultSetStringCache(ctx context.Context, b cdao.BaseRedisOp, v string) error {
	return b.Set(ctx, v)
}

type BaseStringCacheOptions struct {
	GetByCache GetStringByCacheFunc
	SetCache   SetStringCacheFunc
}

func WithGetStringByCache(fn GetStringByCacheFunc) BaseStringCacheOption {
	return func(opts *BaseStringCacheOptions) {
		opts.GetByCache = fn
	}
}

func WithSetStringCache(fn SetStringCacheFunc) BaseStringCacheOption {
	return func(opts *BaseStringCacheOptions) {
		opts.SetCache = fn
	}
}
