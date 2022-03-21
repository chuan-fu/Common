package cache

import (
	"context"

	"github.com/chuan-fu/Common/cdao"
	"github.com/chuan-fu/Common/util"
	"github.com/chuan-fu/Common/zlog"
	"github.com/pkg/errors"
)

/*
func GetHashByDB(db *gorm.DB, id int64) GetHashByDBFunc {
	return func(model interface{}) (err error) {
		err = cdao.FindById(db, id, model)
		if err != nil {
			log.Error(err)
		}
		return
	}
}
*/

type (
	GetHashByCacheFunc func(ctx context.Context, b cdao.BaseRedisOp, model interface{}) error
	SetHashCacheFunc   func(ctx context.Context, b cdao.BaseRedisOp, model interface{}) error
	GetHashByDBFunc    func(ctx context.Context, model interface{}) error

	BaseHashCacheOption func(*BaseHashCacheOptions)
)

// 默认使用json的tag，如不存在，则使用key
func GetBaseHashCache(ctx context.Context, op cdao.BaseRedisOp, model interface{}, getByDb GetHashByDBFunc, opts ...BaseHashCacheOption) (err error) {
	b := &BaseHashCacheOptions{
		CheckExists: defaultCheckExists,
		GetByCache:  defaultGetHashByCache,
		SetCache:    defaultSetHashCache,
		// DelCache:    defaultDelCache,
	}
	for _, opt := range opts {
		opt(b)
	}
	if !util.IsPtrStruct(model) {
		err = errors.New("Model需要为结构体指针")
		return
	}

	var has bool
	has, err = b.CheckExists(ctx, op) // 校验缓存是否存在
	if err != nil {
		log.Error(errors.Wrap(err, "CheckExists"))
	}
	if err == nil && has {
		err = b.GetByCache(ctx, op, model) // 获取缓存，并返回
		if err == nil {
			return nil
		}
		log.Error(errors.Wrap(err, "GetByCache"))
	}

	err = getByDb(ctx, model) // 从db获取
	if err != nil {
		log.Error(errors.Wrap(err, "GetByDB"))
		return
	}

	err = b.SetCache(ctx, op, model) // 写入cache
	if err != nil {
		log.Error(errors.Wrap(err, "SetCache"))
		return nil
	}
	return
}

func defaultGetHashByCache(ctx context.Context, b cdao.BaseRedisOp, model interface{}) error {
	return b.HGetModel(ctx, model)
}

func defaultSetHashCache(ctx context.Context, b cdao.BaseRedisOp, model interface{}) error {
	return b.HSetModel(ctx, model)
}

type BaseHashCacheOptions struct {
	CheckExists CheckExistsFunc
	GetByCache  GetHashByCacheFunc
	SetCache    SetHashCacheFunc
	// DelCache    DelCacheFunc
}

func WithCheckHashExists(fn CheckExistsFunc) BaseHashCacheOption {
	return func(opts *BaseHashCacheOptions) {
		opts.CheckExists = fn
	}
}

func WithGetHashByCache(fn GetHashByCacheFunc) BaseHashCacheOption {
	return func(opts *BaseHashCacheOptions) {
		opts.GetByCache = fn
	}
}

func WithSetHashCache(fn SetHashCacheFunc) BaseHashCacheOption {
	return func(opts *BaseHashCacheOptions) {
		opts.SetCache = fn
	}
}

/*
func WithDelHashCache(fn DelCacheFunc) BaseHashCacheOption {
	return func(opts *BaseHashCacheOptions) {
		opts.DelCache = fn
	}
}
*/
