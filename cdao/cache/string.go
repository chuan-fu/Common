package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/chuan-fu/Common/cdao"
	"github.com/chuan-fu/Common/util"
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
	GetByCacheFunc func(ctx context.Context, b cdao.BaseRedisOp) (string, error)
	SetCacheFunc   func(ctx context.Context, b cdao.BaseRedisOp, v string) error
	DelCacheFunc   func(ctx context.Context, b cdao.BaseRedisOp) error
	GetByDBFunc    func(ctx context.Context, model interface{}) (string, error)

	BaseCacheOption func(*BaseStringCacheOptions)
)

func GetBaseStringCache(ctx context.Context, op cdao.BaseRedisOp, getByDb GetByDBFunc, opts ...BaseCacheOption) (data string, err error) {
	b := &BaseStringCacheOptions{
		GetByCache: defaultGetByCache,
		SetCache:   defaultSetCache,
		DelCache:   defaultDelCache,
	}
	for _, opt := range opts {
		opt(b)
	}
	if b.Model != nil {
		if reflect.TypeOf(b.Model).Kind() != reflect.Ptr {
			err = errors.New("Model is not ptr")
			return
		}
	}

	data, err = b.GetByCache(ctx, op)
	if err != nil {
		log.Error(errors.Wrap(err, "GetByCache"))
	}
	if data != "" {
		if b.Model == nil { // 无需解析
			return
		}
		if err = json.Unmarshal(util.StringToBytes(data), b.Model); err == nil { // 解析model并返回
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
		if b.Model == nil {
			return "", errors.New("Model、data均为空，缓存获取有误")
		}
		dataByte, _ := json.Marshal(b.Model)
		data = util.BytesToString(dataByte)
	}

	err = b.SetCache(ctx, op, data) // 写入cache
	if err != nil {
		log.Error(errors.Wrap(err, "SetCache"))
		return
	}
	return
}

func defaultGetByCache(ctx context.Context, b cdao.BaseRedisOp) (string, error) {
	return b.Get(ctx)
}

func defaultSetCache(ctx context.Context, b cdao.BaseRedisOp, v string) error {
	return b.Set(ctx, v)
}

func defaultDelCache(ctx context.Context, b cdao.BaseRedisOp) error {
	return b.Del(ctx)
}

type BaseStringCacheOptions struct {
	Model interface{}

	GetByCache GetByCacheFunc
	SetCache   SetCacheFunc
	DelCache   DelCacheFunc
}

func WithGetByCache(fn GetByCacheFunc) BaseCacheOption {
	return func(opts *BaseStringCacheOptions) {
		opts.GetByCache = fn
	}
}

func WithSetCache(fn SetCacheFunc) BaseCacheOption {
	return func(opts *BaseStringCacheOptions) {
		opts.SetCache = fn
	}
}

func WithDelCache(fn DelCacheFunc) BaseCacheOption {
	return func(opts *BaseStringCacheOptions) {
		opts.DelCache = fn
	}
}

// model需要为指针，或为空
func WithSetModel(m interface{}) BaseCacheOption {
	return func(opts *BaseStringCacheOptions) {
		opts.Model = m
	}
}
