package cache

import (
	"context"
	"encoding/json"
	"fmt"

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
	GetStringByCacheFunc func(ctx context.Context, b cdao.BaseRedisOp) (string, error)
	SetStringCacheFunc   func(ctx context.Context, b cdao.BaseRedisOp, v string) error
	GetStringByDBFunc    func(ctx context.Context, model interface{}) (string, error)

	BaseStringCacheOption func(*BaseStringCacheOptions)
)

func GetBaseStringCache(ctx context.Context, op cdao.BaseRedisOp, getByDb GetStringByDBFunc, opts ...BaseStringCacheOption) (data string, err error) {
	b := &BaseStringCacheOptions{
		GetByCache: defaultGetStringByCache,
		SetCache:   defaultSetStringCache,
		DelCache:   defaultDelCache,
	}
	for _, opt := range opts {
		opt(b)
	}
	if b.Model != nil && !util.IsPtrStruct(b.Model) {
		err = errors.New("Model需要为空 或者 结构体指针")
		return
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
		return data, nil
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
	Model interface{}

	GetByCache GetStringByCacheFunc
	SetCache   SetStringCacheFunc
	DelCache   DelCacheFunc
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

func WithDelStringCache(fn DelCacheFunc) BaseStringCacheOption {
	return func(opts *BaseStringCacheOptions) {
		opts.DelCache = fn
	}
}

// model需要为指针，或为空
func WithSetModel(m interface{}) BaseStringCacheOption {
	return func(opts *BaseStringCacheOptions) {
		opts.Model = m
	}
}
