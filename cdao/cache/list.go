package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/chuan-fu/Common/util"

	"github.com/chuan-fu/Common/cdao"
	"github.com/chuan-fu/Common/zlog"
	"github.com/pkg/errors"
)

// use zet 或者 zset
// 需要删除，因为可能重复

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
	GetListByCacheFunc func(ctx context.Context, b cdao.BaseRedisOp, pageIndex, pageSize int64) ([]string, error)
	SetListCacheFunc   func(ctx context.Context, b cdao.BaseRedisOp, v []string) error
	GetListByDBFunc    func(ctx context.Context) (interface{}, error)

	BaseListCacheOption func(*BaseListCacheOptions)
)

// GetListByDBFunc 返回值的interface{} 需要为[]string 或者 slice 或者 &slice
// slice 和 &slice 会转换成[]string
// 且不可为空，若为空，则每次都会打到DB
func GetBaseListCache(ctx context.Context, op cdao.BaseRedisOp, pageIndex, pageSize int64, getByDb GetListByDBFunc, opts ...BaseListCacheOption) (resp []string, err error) {
	b := &BaseListCacheOptions{
		CheckExists: defaultCheckExists,
		GetByCache:  defaultGetListByCache,
		SetCache:    defaultSetListCache,
	}
	for _, opt := range opts {
		opt(b)
	}

	var has bool
	has, err = b.CheckExists(ctx, op) // 校验缓存是否存在
	if err != nil {
		log.Error(errors.Wrap(err, "CheckExists"))
	}
	if err == nil && has {
		resp, err = b.GetByCache(ctx, op, pageIndex, pageSize) // 获取缓存，并返回
		if err == nil {
			return resp, nil
		}
		log.Error(errors.Wrap(err, "GetByCache"))
	}

	var data interface{}
	data, err = getByDb(ctx) // 从db获取
	if err != nil {
		log.Error(errors.Wrap(err, "GetByDB"))
		return
	}

	var ok bool
	resp, ok = data.([]string)
	if !ok {
		resp, err = toStringSlice(data)
		if err != nil {
			log.Error(errors.Wrap(err, "toStringSlice"))
			return
		}
	}
	if len(resp) > 0 {
		err = b.SetCache(ctx, op, resp) // 写入cache
		if err != nil {
			log.Error(errors.Wrap(err, "SetCache"))
			return
		}
	}

	// 返回值数据分页
	start := (pageIndex - 1) * pageSize
	stop := start + pageSize
	total := int64(len(resp))
	if start >= total {
		return []string{}, nil
	}
	if stop >= total {
		return resp[start:], nil
	}
	return resp[start:stop], nil
}

func toStringSlice(data interface{}) ([]string, error) {
	if !util.IsSliceOrPtrSlice(data) {
		return nil, errors.New("getByDb返回值需要为切片 或者 切片指针")
	}
	rv := reflect.ValueOf(data)
	if reflect.TypeOf(data).Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	resp := make([]string, 0, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		val := rv.Index(i)
		switch val.Kind() {
		case reflect.String:
			resp = append(resp, val.String())
		case reflect.Struct, reflect.Map, reflect.Array, reflect.Slice:
			d, err := json.Marshal(val.Interface())
			if err != nil {
				return nil, err
			}
			resp = append(resp, util.BytesToString(d))
		default:
			resp = append(resp, fmt.Sprintf("%v", val.Interface()))
		}
	}
	return resp, nil
}

func defaultGetListByCache(ctx context.Context, b cdao.BaseRedisOp, pageIndex, pageSize int64) ([]string, error) {
	return b.ZRangeStringWithPage(ctx, pageIndex, pageSize)
}

func defaultSetListCache(ctx context.Context, b cdao.BaseRedisOp, vs []string) error {
	return b.ZAddString(ctx, vs)
}

type BaseListCacheOptions struct {
	CheckExists CheckExistsFunc
	GetByCache  GetListByCacheFunc
	SetCache    SetListCacheFunc
}

func WithCheckListExists(fn CheckExistsFunc) BaseListCacheOption {
	return func(opts *BaseListCacheOptions) {
		opts.CheckExists = fn
	}
}

func WithGetListByCache(fn GetListByCacheFunc) BaseListCacheOption {
	return func(opts *BaseListCacheOptions) {
		opts.GetByCache = fn
	}
}

func WithSetListCache(fn SetListCacheFunc) BaseListCacheOption {
	return func(opts *BaseListCacheOptions) {
		opts.SetCache = fn
	}
}
