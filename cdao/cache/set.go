package cache

import (
	"context"

	"github.com/chuan-fu/Common/cdao"
	"github.com/chuan-fu/Common/zlog"
)

/*
set存储，cache使用覆盖式写入防止重复
如getByDb为空，cache会写入一个空字符串
*/

func (c *CacheHandle) GetBaseSetCache(ctx context.Context, op cdao.BaseRedisOp, getByDb GetSetByDBFunc) ([]string, error) {
	if c.sf == nil {
		return GetBaseSetCache(ctx, op, getByDb)
	}

	dataInter, err := c.sf.Do(op.GetKey(), func() (interface{}, error) {
		return GetBaseSetCache(ctx, op, getByDb)
	})
	if err != nil {
		return nil, err
	}
	return dataInter.([]string), nil
}

func GetBaseSetCache(ctx context.Context, op cdao.BaseRedisOp, getByDb GetSetByDBFunc) (resp []string, err error) {
	resp, err = op.SGetAll(ctx) // 获取缓存，并返回
	if err == nil && len(resp) > 0 {
		if checkEntryList(resp) {
			return []string{}, nil
		}
		return resp, nil
	}
	if err != nil {
		log.Error("GetByCache:", err)
	}

	// 从db获取
	resp, err = getByDb(ctx)
	if err != nil {
		log.Error("GetByDB:", err)
		return
	}

	// 写入cache
	if len(resp) == 0 {
		if err2 := op.SAddCover(ctx, entryList); err2 != nil { // 写入空缓存
			log.Error("SetCache:", err)
		}
	} else {
		if err2 := op.SAddCover(ctx, resp); err2 != nil {
			log.Error("SetCache:", err)
		}
	}
	return
}
