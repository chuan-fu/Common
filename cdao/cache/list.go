package cache

import (
	"context"

	"github.com/chuan-fu/Common/cdao"
	"github.com/chuan-fu/Common/zlog"
	"github.com/pkg/errors"
)

/*
使用zset存储，cache使用覆盖式写入防止重复
如getByDb为空，cache会写入一个空字符串
*/

func (c *CacheHandle) GetBaseListCache(ctx context.Context, op cdao.BaseRedisOp, getByDb GetListByDBFunc) ([]string, error) {
	if c.sf == nil {
		return GetBaseListCache(ctx, op, getByDb)
	}

	dataInter, err := c.sf.Do(op.GetKey(), func() (interface{}, error) {
		return GetBaseListCache(ctx, op, getByDb)
	})
	if err != nil {
		return nil, err
	}
	return dataInter.([]string), nil
}

func (c *CacheHandle) GetBaseListCacheWithPage(ctx context.Context, op cdao.BaseRedisOp, getByDb GetListByDBFunc, pageIndex, pageSize int64) ([]string, error) {
	if c.sf == nil {
		return GetBaseListCacheWithPage(ctx, op, getByDb, pageIndex, pageSize)
	}

	dataInter, err := c.sf.Do(op.GetKey(), func() (interface{}, error) {
		return GetBaseListCacheWithPage(ctx, op, getByDb, pageIndex, pageSize)
	})
	if err != nil {
		return nil, err
	}
	return dataInter.([]string), nil
}

func GetBaseListCache(ctx context.Context, op cdao.BaseRedisOp, getByDb GetListByDBFunc) (resp []string, err error) {
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

	resp, err = getAndSetList(ctx, op, getByDb)
	if err != nil {
		log.Error(err)
	}
	return
}

func GetBaseListCacheWithPage(ctx context.Context, op cdao.BaseRedisOp, getByDb GetListByDBFunc, pageIndex, pageSize int64) (resp []string, err error) {
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

	resp, err = getAndSetList(ctx, op, getByDb)
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

func getAndSetList(ctx context.Context, op cdao.BaseRedisOp, getByDb GetListByDBFunc) (resp []string, err error) {
	resp, err = getByDb(ctx)
	if err != nil {
		log.Error("GetByDB:", err)
		return
	}

	// 写入缓存
	var err2 error
	if len(resp) == 0 {
		err2 = op.ZAddCoverStringList(ctx, entryList) // 空缓存处理
	} else {
		err2 = op.ZAddCoverStringList(ctx, resp)
	}
	if err2 != nil {
		log.Error("SetCache:", err2)
	}
	return
}
