package cache

import (
	"context"
	"fmt"
	"reflect"

	"github.com/chuan-fu/Common/baseservice/jsonx"
	"github.com/chuan-fu/Common/cdao"
	"github.com/chuan-fu/Common/zlog"
	"github.com/pkg/errors"
)

/*
func defaultGetByDB(db *gorm.DB, id int64) GetByDBFunc {
	return func(model interface{}) (data string, err error) {
		err = cdao.FindById(db, id, model)
		if err != nil {
			log.Error(err)
			return
		}
		return "", nil
	}
}
*/

type (
	GetJsonByCacheFunc func(ctx context.Context, b cdao.BaseRedisOp) (string, error)
	SetJsonCacheFunc   func(ctx context.Context, b cdao.BaseRedisOp, v string) error
	GetJsonByDBFunc    func(ctx context.Context, model interface{}) (string, error)

	BaseJsonCacheOption func(*BaseJsonCacheOptions)
)

func GetBaseJsonCache(ctx context.Context, op cdao.BaseRedisOp, model interface{}, getByDb GetJsonByDBFunc, opts ...BaseJsonCacheOption) (data string, err error) {
	if reflect.TypeOf(model).Kind() != reflect.Ptr {
		err = errors.New("Model需要为指针")
		return
	}

	b := &BaseJsonCacheOptions{
		Model: model, // 取指针

		GetByCache: defaultGetJsonByCache,
		SetCache:   defaultSetJsonCache,
		DelCache:   defaultDelCache,
	}
	for _, opt := range opts {
		opt(b)
	}

	data, err = b.GetByCache(ctx, op)
	if err != nil {
		log.Error(errors.Wrap(err, "GetByCache"))
	}
	if data != "" {
		if err = jsonx.UnmarshalObj(data, b.Model); err == nil { // 解析model并返回
			return data, nil
		}

		log.Error(errors.Wrap(err, fmt.Sprintf("Cache【%s】Unmarshal", data)))
		err = b.DelCache(ctx, op) // 解析失败，缓存有误，删除
		if err != nil {
			log.Error(errors.Wrap(err, "DelCache"))
		}
	}

	data, err = getByDb(ctx, b.Model) // 从db获取
	if err != nil {
		log.Error(errors.Wrap(err, "GetByDB"))
		return
	}

	if data == "" {
		data = jsonx.MarshalObj(b.Model)
	}

	err = b.SetCache(ctx, op, data) // 写入cache
	if err != nil {
		log.Error(errors.Wrap(err, "SetCache"))
		return data, nil
	}
	return
}

func defaultGetJsonByCache(ctx context.Context, b cdao.BaseRedisOp) (string, error) {
	return b.Get(ctx)
}

func defaultSetJsonCache(ctx context.Context, b cdao.BaseRedisOp, v string) error {
	return b.Set(ctx, v)
}

type BaseJsonCacheOptions struct {
	Model interface{}

	GetByCache GetJsonByCacheFunc
	SetCache   SetJsonCacheFunc
	DelCache   DelCacheFunc
}

func WithGetJsonByCache(fn GetJsonByCacheFunc) BaseJsonCacheOption {
	return func(opts *BaseJsonCacheOptions) {
		opts.GetByCache = fn
	}
}

func WithSetJsonCache(fn SetJsonCacheFunc) BaseJsonCacheOption {
	return func(opts *BaseJsonCacheOptions) {
		opts.SetCache = fn
	}
}

func WithDelJsonCache(fn DelCacheFunc) BaseJsonCacheOption {
	return func(opts *BaseJsonCacheOptions) {
		opts.DelCache = fn
	}
}
