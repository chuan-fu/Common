package redislock

import (
	"context"
	"time"

	"github.com/chuan-fu/Common/db"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

type RedisLockOp interface {
	SetLock(ctx context.Context) (ok bool, err error)
	ExtendLock(ctx context.Context) error
	DelLock(ctx context.Context) error
}

var (
	DelLockStatusNotOwnErr  = errors.New("RedisLockOp DelLock: 锁已过期，且已被抢占")
	DelLockStatusExpiredErr = errors.New("RedisLockOp DelLock: 锁已过期")
)

type redisLock struct {
	store redis.Cmdable
	key   string
	ttl   time.Duration
	value string
}

type Option func(r *redisLock)

func WithStore(store redis.Cmdable) Option {
	return func(r *redisLock) {
		r.store = store
	}
}

func WithValue(value string) Option {
	return func(r *redisLock) {
		r.value = value
	}
}

func NewRedisLock(key string, ttl time.Duration, opts ...Option) RedisLockOp {
	r := &redisLock{
		key:   key,
		ttl:   ttl,
		value: uuid.NewV4().String(),
	}
	for _, opt := range opts {
		opt(r)
	}
	if r.store == nil {
		r.store = db.GetRedisCli()
	}
	return r
}

// 写入锁 ok表示是否成功
func (r *redisLock) SetLock(ctx context.Context) (ok bool, err error) {
	ok, err = r.store.SetNX(ctx, r.key, r.value, r.ttl).Result()
	if err != nil {
		err = errors.Wrap(err, "RedisLockOp SetLock")
	}
	return
}

// 延长锁
func (r *redisLock) ExtendLock(ctx context.Context) error {
	resI, err := r.store.Eval(ctx, extendLockScript, []string{r.key}, r.value, formatMs(r.ttl)).Result()
	if err != nil {
		err = errors.Wrap(err, "RedisLockOp DelLock")
		return err
	}

	if res, _ := resI.(int64); res == DelLockStatusNotOwn {
		return DelLockStatusNotOwnErr
	}
	return nil
}

// 删除锁
// value为uuid，ok为是否写入成功
func (r *redisLock) DelLock(ctx context.Context) error {
	resI, err := r.store.Eval(ctx, delLockScript, []string{r.key}, r.value).Result()
	if err != nil {
		err = errors.Wrap(err, "RedisLockOp DelLock")
		return err
	}

	if res, _ := resI.(int64); res == DelLockStatusNotOwn {
		return DelLockStatusNotOwnErr
	}
	return nil
	/*
		res, _ := resI.(int64)
		switch res {
		case DelLockStatusExpired:
			return DelLockStatusExpiredErr
		case DelLockStatusNotOwn:
			return DelLockStatusNotOwnErr
		}
		return nil
	*/
}

func (r *redisLock) Value() string {
	return r.value
}
