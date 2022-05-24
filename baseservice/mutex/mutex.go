package mutex

import (
	"context"
	"time"

	"github.com/chuan-fu/Common/db"
	"github.com/chuan-fu/Common/zlog"
	"github.com/go-redis/redis/v8"
)

// 可更改
var DefaultExpiration = defaultExpiration

const (
	defaultExpiration = 10 * time.Second
)

type distributedOnce struct {
	key        string
	expiration time.Duration
	atomic     DistributeAtomic
}

func (o *distributedOnce) Do(f func()) {
	if o.atomic.AtomicSetNotExist(o.key, o.expiration) {
		f()
	}
}

// 分布式原子操作
type DistributeAtomic interface {
	AtomicSetNotExist(key string, expiration time.Duration) bool
}

type Option func(*distributedOnce)

// expiration时长内只执行一次
func WithExpiration(expiration time.Duration) Option {
	return func(once *distributedOnce) {
		once.expiration = expiration
	}
}

// 分布式的同步器
func WithDistributeAtomic(atomic DistributeAtomic) Option {
	return func(once *distributedOnce) {
		once.atomic = atomic
	}
}

// 分布式sync.Once
func NewDistributedOnce(key string, opts ...Option) *distributedOnce {
	if key == "" {
		log.Fatal("NewDistributedOnce key is empty")
	}
	o := &distributedOnce{
		key:        key,
		expiration: DefaultExpiration,
	}
	for _, opt := range opts {
		opt(o)
	}

	if o.atomic == nil {
		o.atomic = NewRedisAtomic(db.GetRedisCli())
	}

	return o
}

func NewRedisAtomic(redisCli redis.Cmdable) DistributeAtomic {
	return &redisAtomic{redisCli}
}

type redisAtomic struct {
	redisCli redis.Cmdable
}

func (a *redisAtomic) AtomicSetNotExist(key string, expiration time.Duration) bool {
	res := a.redisCli.SetNX(context.TODO(), key, time.Now().Unix(), expiration)
	if r, err := res.Result(); err != nil || !r {
		return false
	}
	return true
}
