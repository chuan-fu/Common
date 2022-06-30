package cache

import (
	"context"

	"github.com/chuan-fu/Common/baseservice/syncx"

	"github.com/chuan-fu/Common/cdao"
	"github.com/chuan-fu/Common/zlog"
	"github.com/pkg/errors"
)

/*
使用zset存储，cache使用覆盖式写入防止重复
如getByDb为空，cache会写入一个空字符串
list的单飞和别的不一样，因为存在行数、列数，需要在读取db时使用，而不能是在外层使用
*/

func (c *CacheHandle) GetBaseListCache(ctx context.Context, op cdao.BaseRedisOp, getByDb GetListByDBFunc) ([]string, error) {
	return GetBaseListCache(ctx, op, getByDb, c.sf)
}

func (c *CacheHandle) GetBaseListCacheWithPage(ctx context.Context, op cdao.BaseRedisOp, getByDb GetListByDBFunc, pageIndex, pageSize int64) ([]string, error) {
	return GetBaseListCacheWithPage(ctx, op, getByDb, pageIndex, pageSize, c.sf)
}

func GetBaseListCache(ctx context.Context, op cdao.BaseRedisOp, getByDb GetListByDBFunc, sf syncx.SingleFight) (resp []string, err error) {
	var has bool
	resp, has, err = op.ZGetAll(ctx)
	if err == nil && has { // 缓存存在
		if checkEntryList(resp) { // 空缓存处理
			return []string{}, nil
		}
		return
	}
	if err != nil {
		log.Error("GetByCache:", err)
	}

	resp, err = getAndSetList(ctx, op, getByDb, sf)
	if err != nil {
		log.Error(err)
	}
	return
}

func GetBaseListCacheWithPage(ctx context.Context, op cdao.BaseRedisOp, getByDb GetListByDBFunc, pageIndex, pageSize int64, sf syncx.SingleFight) (resp []string, err error) {
	if pageIndex < 1 || pageSize < 1 {
		return resp, errors.New("pageIndex、pageSize有误")
	}

	var has bool
	resp, has, err = op.ZRangeStringListWithPage(ctx, pageIndex, pageSize)
	if err == nil && has { // 缓存存在
		if pageIndex == 1 && checkEntryList(resp) { // 第一页，且为空缓存
			return []string{}, nil
		}
		return
	}
	if err != nil {
		log.Error("GetByCache:", err)
	}

	resp, err = getAndSetList(ctx, op, getByDb, sf)
	if err != nil {
		log.Error(err)
		return
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

func getAndSetList(ctx context.Context, op cdao.BaseRedisOp, getByDb GetListByDBFunc, sf syncx.SingleFight) ([]string, error) {
	f := func() (interface{}, error) {
		resp, err := getByDb(ctx)
		if err != nil {
			log.Error("GetByDB:", err)
			return nil, err
		}

		// 写入缓存
		if len(resp) == 0 {
			if err2 := op.ZAddCoverStringList(ctx, entryList); err2 != nil { // 空缓存处理
				log.Error("SetCache:", err2)
			}
		} else {
			if err2 := op.ZAddCoverStringList(ctx, resp); err2 != nil {
				log.Error("SetCache:", err2)
			}
		}
		return resp, nil
	}

	var dataInter interface{}
	var err error
	if sf != nil {
		dataInter, err = sf.Do(op.GetKey(), f)
	} else {
		dataInter, err = f()
	}
	if err != nil {
		return nil, err
	}
	return dataInter.([]string), nil
}
