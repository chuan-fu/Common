package cdao

import (
	"context"
	"encoding/json"
	"time"

	dbRedis "github.com/chuan-fu/Common/db/redis"
	"github.com/chuan-fu/Common/util"
	"github.com/chuan-fu/Common/zlog"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

type BaseRedisOp interface {
	SetKey(key string) BaseRedisOp
	GetKey() string
	SetTTL(ttl time.Duration) BaseRedisOp
	GetTTL() time.Duration

	Set(ctx context.Context, value string) error
	SetNX(ctx context.Context, value string) error
	Get(ctx context.Context) (string, error)
	Del(ctx context.Context) error
}

type baseRedisOp struct {
	redisCli redis.Cmdable
	key      string
	ttl      time.Duration
}

func NewBaseRedisOp(redisCli ...redis.Cmdable) BaseRedisOp {
	if len(redisCli) > 0 {
		return &baseRedisOp{redisCli: redisCli[0]}
	}
	return &baseRedisOp{redisCli: dbRedis.GetRedisCli()}
}

func NewBaseRedisOpWithKT(key string, ttl time.Duration, redisCli ...redis.Cmdable) BaseRedisOp {
	return NewBaseRedisOp(redisCli...).SetKey(key).SetTTL(ttl)
}

func (b *baseRedisOp) SetKey(key string) BaseRedisOp {
	b.key = key
	return b
}

func (b *baseRedisOp) GetKey() string {
	return b.key
}

func (b *baseRedisOp) SetTTL(ttl time.Duration) BaseRedisOp {
	b.ttl = ttl
	return b
}

func (b *baseRedisOp) GetTTL() time.Duration {
	return b.ttl
}

func (b *baseRedisOp) Set(ctx context.Context, value string) error {
	_, err := b.redisCli.Set(ctx, b.key, value, b.ttl).Result()
	if err != nil {
		return errors.Wrap(err, "BaseRedisOp Set")
	}
	return nil
}

// 不存在则写入
func (b *baseRedisOp) SetNX(ctx context.Context, value string) error {
	_, err := b.redisCli.SetNX(ctx, b.key, value, b.ttl).Result()
	if err != nil {
		return errors.Wrap(err, "BaseRedisOp SetNX")
	}
	return nil
}

// 如果不存在，data为空字符串，不报错
func (b *baseRedisOp) Get(ctx context.Context) (string, error) {
	data, err := b.redisCli.Get(ctx, b.key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return "", errors.Wrap(err, "BaseRedisOp Get")
	}
	if data != "" {
		log.Debugf("cache【%s】存在", b.key)
	}
	return data, nil
}

// v应为指针
func (b *baseRedisOp) GetResult(ctx context.Context, v interface{}) error {
	data, err := b.redisCli.Get(ctx, b.key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return errors.Wrap(err, "BaseRedisOp GetResult")
	}
	if data == "" {
		return nil
	}
	log.Debugf("cache【%s】存在", b.key)
	return json.Unmarshal(util.StringToBytes(data), v)
}

func (b *baseRedisOp) Del(ctx context.Context) error {
	_, err := b.redisCli.Del(ctx, b.key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return errors.Wrap(err, "BaseRedisOp Del")
	}
	return nil
}
