package cdao

import (
	"context"
	"time"
)

type BaseRedisOp interface {
	SetKey(key string)
	GetKey() string
	SetTTL(ttl time.Duration)
	GetTTL() time.Duration

	Exists(ctx context.Context) (bool, error)
	Expire(ctx context.Context) (bool, error)
	TTL(ctx context.Context) (time.Duration, error)
	Del(ctx context.Context) error

	Set(ctx context.Context, value string) error
	SetNX(ctx context.Context, value string) (bool, error)
	Get(ctx context.Context) (string, error)
	GetResult(ctx context.Context, v interface{}) error

	HSetMap(ctx context.Context, m map[string]interface{}) error
	HGetMap(ctx context.Context) (map[string]string, error)
	HGet(ctx context.Context, key string) (string, error)

	SAddCover(ctx context.Context, list []string) error     // 覆盖式写入
	SGetAll(ctx context.Context) (data []string, err error) // 读取所有 空数组表示缓存不存在

	ZAddCoverStringList(ctx context.Context, list []string) error
	ZGetAll(ctx context.Context) (data []string, has bool, err error)
	ZRangeStringList(ctx context.Context, start, stop int64) (data []string, has bool, err error)                 // 根据下标
	ZRangeStringListWithPage(ctx context.Context, pageIndex, pageSize int64) (data []string, has bool, err error) // 根据分页

	SetBits(ctx context.Context, value []int64) error
	GetBits(ctx context.Context, value []int64) (resp map[int64]struct{}, exists bool, err error)
}
