package cache

import (
	"context"
	"fmt"

	"github.com/chuan-fu/Common/baseservice/jsonx"
	"github.com/chuan-fu/Common/cdao"
	"github.com/chuan-fu/Common/zlog"
)

/*
func defaultGetByDB(db *gorm.DB, id int64) GetByDBFunc {
	return func(model interface{}) error {
		err := cdao.FindById(db, id, model)
		if err != nil {
			log.Error(err)
		}
		return err
	}
}
*/

/*
model应为指针
model传参错误 为开发中错误，不校验
*/

func (c *CacheHandle) GetBaseJsonCache(ctx context.Context, op cdao.BaseRedisOp, model interface{}, getByDb GetJsonByDBFunc) (string, error) {
	if c.sf == nil {
		return GetBaseJsonCache(ctx, op, model, getByDb)
	}

	dataInter, err := c.sf.Do(op.GetKey(), func() (interface{}, error) {
		return GetBaseJsonCache(ctx, op, model, getByDb)
	})
	if err != nil {
		return "", err
	}
	return dataInter.(string), nil
}

func GetBaseJsonCache(ctx context.Context, op cdao.BaseRedisOp, model interface{}, getByDb GetJsonByDBFunc) (data string, err error) {
	data, err = op.Get(ctx)
	if err != nil {
		log.Error("GetByCache:", err)
	}
	if data != "" {
		if err = jsonx.UnmarshalObj(data, model); err == nil { // 解析model并返回
			return data, nil
		}

		log.Error(fmt.Sprintf("Cache【%s】Unmarshal:", data), err)
		err = op.Del(ctx) // 解析失败，缓存有误，删除
		if err != nil {
			log.Error("DelCache:", err)
		}
	}

	err = getByDb(ctx, model) // 从db获取
	if err != nil {
		log.Error("GetByDB:", err)
		return
	}
	data = jsonx.MarshalObj(model)

	// 写入cache
	if err2 := op.Set(ctx, data); err2 != nil {
		log.Error("SetCache:", err2)
	}
	return
}
