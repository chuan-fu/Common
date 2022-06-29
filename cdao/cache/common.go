package cache

import (
	"context"

	"github.com/chuan-fu/Common/baseservice/syncx"
	"github.com/pkg/errors"
)

type (
	GetJsonByDBFunc   func(ctx context.Context, model interface{}) error
	GetStringByDBFunc func(ctx context.Context) (string, error)
	GetListByDBFunc   func(ctx context.Context) ([]string, error)
	GetSetByDBFunc    func(ctx context.Context) ([]string, error)
)

var BaseGetByDBDataNilError = errors.New("getByDB获取的数据为空")

type CacheHandle struct {
	sf syncx.SingleFight
}

func NewCacheHandle(sf syncx.SingleFight) *CacheHandle {
	return &CacheHandle{sf: sf}
}

// 默认为不带超时的单飞模式
// 可替换
var C = &CacheHandle{sf: syncx.NewSingleFlight()}

var entryList = []string{""}

func checkEntryList(v []string) bool {
	return len(v) == 1 && v[0] == ""
}
