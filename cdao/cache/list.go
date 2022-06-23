package cache

import (
	"context"

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
	GetListByCacheFunc func(ctx context.Context, b cdao.BaseRedisOp, isAll bool, pageIndex, pageSize int64) ([]string, error)
	SetListCacheFunc   func(ctx context.Context, b cdao.BaseRedisOp, v []string) error
	GetListByDBFunc    func(ctx context.Context) ([]string, error)

	BaseListCacheOption func(*BaseListCacheOptions)
)

// 使用zset存储，cache使用覆盖式写入防止重复
// GetListByDBFunc 返回值的interface{} 需要为[]string 或者 slice 或者 &slice
// slice 和 &slice 会转换成[]string
// 如getByDb为空，cache会写入一个空字符串
func GetBaseListCache(ctx context.Context, op cdao.BaseRedisOp, getByDb GetListByDBFunc, opts ...BaseListCacheOption) (resp []string, err error) {
	b := &BaseListCacheOptions{
		CheckExists: defaultCheckExists,
		GetByCache:  defaultGetListByCache,
		SetCache:    defaultSetListCache,
		IsAll:       true,
	}
	for _, opt := range opts {
		opt(b)
	}

	if !b.IsAll && (b.PageIndex < 1 || b.PageSize < 1) {
		err = errors.New("页数、行数有误")
		return
	}

	var has bool
	has, err = b.CheckExists(ctx, op) // 校验缓存是否存在
	if err != nil {
		log.Error(errors.Wrap(err, "CheckExists"))
	}
	if err == nil && has {
		resp, err = b.GetByCache(ctx, op, b.IsAll, b.PageIndex, b.PageSize) // 获取缓存，并返回
		if err == nil {
			// 第一页或者全部，且只有一条空数据，代表为空
			if (b.IsAll || b.PageIndex == 1) && len(resp) == 1 && resp[0] == "" {
				resp = []string{}
			}
			return resp, nil
		} else {
			log.Error(errors.Wrap(err, "GetByCache"))
		}
	}

	resp, err = getByDb(ctx) // 从db获取
	if err != nil {
		log.Error(errors.Wrap(err, "GetByDB"))
		return
	}

	// 写入Cache
	if len(resp) == 0 {
		err = b.SetCache(ctx, op, []string{""}) // 写入空数据
	} else {
		err = b.SetCache(ctx, op, resp)
	}
	if err != nil {
		log.Error(errors.Wrap(err, "SetCache"))
		err = nil
	}

	if b.IsAll { // 返回所有
		return
	}
	// 返回值数据分页
	start := (b.PageIndex - 1) * b.PageSize
	stop := start + b.PageSize
	total := int64(len(resp))
	if start >= total {
		return []string{}, nil
	}
	if stop >= total {
		return resp[start:], nil
	}
	return resp[start:stop], nil
}

func defaultGetListByCache(ctx context.Context, b cdao.BaseRedisOp, isAll bool, pageIndex, pageSize int64) ([]string, error) {
	if isAll {
		return b.ZGetAll(ctx)
	}
	return b.ZRangeStringListWithPage(ctx, pageIndex, pageSize)
}

func defaultSetListCache(ctx context.Context, b cdao.BaseRedisOp, vs []string) error {
	return b.ZAddCoverStringList(ctx, vs)
}

type BaseListCacheOptions struct {
	CheckExists         CheckExistsFunc
	GetByCache          GetListByCacheFunc
	SetCache            SetListCacheFunc
	IsAll               bool
	PageIndex, PageSize int64
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

func WithSetListPage(pageIndex, pageSize int64) BaseListCacheOption {
	return func(opts *BaseListCacheOptions) {
		opts.IsAll = false
		opts.PageIndex = pageIndex
		opts.PageSize = pageSize
	}
}
