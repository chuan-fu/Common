package cache

import (
	"context"

	"github.com/chuan-fu/Common/cdao"
	"github.com/chuan-fu/Common/zlog"
)

/*
func defaultGetByDB(db *gorm.DB, id int64) GetByDBFunc {
	return func() (data string, err error) {
		return "11", nil
	}
}
*/

// 注意，不允许存入空字符串
func (c *CacheHandle) GetBaseStringCache(ctx context.Context, op cdao.BaseRedisOp, getByDb GetStringByDBFunc) (string, error) {
	if c.sf == nil {
		return GetBaseStringCache(ctx, op, getByDb)
	}
	dataInter, err := c.sf.Do(op.GetKey(), func() (interface{}, error) {
		return GetBaseStringCache(ctx, op, getByDb)
	})
	if err != nil {
		return "", err
	}
	return dataInter.(string), nil
}

func GetBaseStringCache(ctx context.Context, op cdao.BaseRedisOp, getByDb GetStringByDBFunc) (data string, err error) {
	data, err = op.Get(ctx)
	if err == nil && data != "" {
		return
	}
	if err != nil {
		log.Error("GetByCache:", err)
	}

	data, err = getByDb(ctx) // 从db获取
	if err != nil {
		log.Error("GetByDB:", err)
		return
	}
	if data == "" {
		return "", BaseGetByDBDataNilError
	}

	// 写入cache
	if err2 := op.Set(ctx, data); err2 != nil {
		log.Error("SetCache:", err2)
	}
	return
}
