package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"sync"

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
	GetByCacheFunc func(b cdao.BaseRedisOp) (string, error)
	SetCacheFunc   func(b cdao.BaseRedisOp, v string) error
	DelCacheFunc   func(b cdao.BaseRedisOp) error

	GetByDBFunc func(model interface{}) (string, error)
)

type BaseStringCache struct {
	Op    cdao.BaseRedisOp
	Model interface{} // 应为指针

	GetByCache GetByCacheFunc
	GetByDB    GetByDBFunc
	SetCache   SetCacheFunc
	DelCache   DelCacheFunc

	once sync.Once
}

func (b *BaseStringCache) init() error {
	if b.Op == nil {
		return errors.New("BaseRedisOp is nil")
	}
	if reflect.TypeOf(b.Model).Kind() != reflect.Ptr {
		return errors.New("Model is not ptr")
	}

	if b.GetByDB == nil {
		return errors.New("GetByDB is nil")
	}
	b.once.Do(func() {
		if b.GetByCache == nil {
			b.GetByCache = defaultGetByCache
		}
		if b.SetCache == nil {
			b.SetCache = defaultSetCache
		}
		if b.DelCache == nil {
			b.DelCache = defaultDelCache
		}
	})
	return nil
}

func (b *BaseStringCache) Get() (data string, err error) {
	err = b.init()
	if err != nil {
		return
	}

	data, err = b.GetByCache(b.Op)
	if err != nil {
		log.Error(errors.Wrap(err, "GetByCache"))
	}
	if data != "" {
		err = json.Unmarshal(util.StringToBytes(data), b.Model)
		if err == nil {
			return
		}

		log.Error(errors.Wrap(err, fmt.Sprintf("Cache【%s】Unmarshal", data)))
		err = b.DelCache(b.Op) // 解析失败，缓存有误，删除
		if err != nil {
			log.Error(errors.Wrap(err, "DelCache"))
		}
	}

	data, err = b.GetByDB(b.Model)
	if err != nil {
		log.Error(errors.Wrap(err, "GetByDB"))
		return
	}

	if data == "" {
		dataByte, _ := json.Marshal(b.Model)
		data = util.BytesToString(dataByte)
	}
	err = b.SetCache(b.Op, data)
	if err != nil {
		log.Error(errors.Wrap(err, "SetCache"))
		return
	}
	return
}

func defaultGetByCache(b cdao.BaseRedisOp) (string, error) {
	return b.Get(context.TODO())
}

func defaultSetCache(b cdao.BaseRedisOp, v string) error {
	return b.Set(context.TODO(), v)
}

func defaultDelCache(b cdao.BaseRedisOp) error {
	return b.Del(context.TODO())
}
