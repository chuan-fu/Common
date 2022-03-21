package cache

import (
	"context"

	"github.com/chuan-fu/Common/cdao"
)

type (
	GetSetByCacheFunc func(ctx context.Context, b cdao.BaseRedisOp, pageIndex, pageSize int64) ([]string, error)
	SetSetCacheFunc   func(ctx context.Context, b cdao.BaseRedisOp, v []string) error
	GetSetByDBFunc    func(ctx context.Context) (interface{}, error)

	BaseSetCacheOption func(*BaseSetCacheOptions)
)

type BaseSetCacheOptions struct {
	CheckExists CheckExistsFunc
	GetByCache  GetSetByCacheFunc
	SetCache    SetSetCacheFunc
}

func WithCheckSetExists(fn CheckExistsFunc) BaseSetCacheOption {
	return func(opts *BaseSetCacheOptions) {
		opts.CheckExists = fn
	}
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
