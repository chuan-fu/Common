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
	"github.com/satori/go.uuid"
)

type BaseRedisOp interface {
	GetCli() redis.Cmdable
	SetKey(key string) BaseRedisOp
	GetKey() string
	SetTTL(ttl time.Duration) BaseRedisOp
	GetTTL() time.Duration

	Set(ctx context.Context, value string) error
	SetNX(ctx context.Context, value string) (bool, error)
	Get(ctx context.Context) (string, error)
	Del(ctx context.Context) error

	SetLock(ctx context.Context) (value string, ok bool, err error)
	ExtendLock(ctx context.Context, value string) (int64, error)
	DelLock(ctx context.Context, value string) (int64, error)

	PFAdd(ctx context.Context, els ...interface{}) error
	PFCount(ctx context.Context) (count int64, err error)
	PFMerge(ctx context.Context, keys ...string) error
}

/* delLockScript
local v2 = redis.call("get", KEYS[1])
if v2 then
	if v2 == ARGV[1] then
		redis.call("del", KEYS[1])
		return 1
	end
	return -1
end
return 0
*/

/* extendLockScript
local v2 = redis.call("get", KEYS[1])
if v2 then
	if v2 == ARGV[1] then
		redis.call("pexpire", KEYS[1], ARGV[2])
		return 1
	end
	return -1
else
	redis.call("set", KEYS[1], ARGV[1], "px", ARGV[2])
	return 1
end
*/
const (
	extendLockScript     = `local v2 = redis.call("get", KEYS[1]) if v2 then if v2 == ARGV[1] then redis.call("pexpire", KEYS[1], ARGV[2]) return 1 end return -1 else redis.call("set", KEYS[1], ARGV[1], "px", ARGV[2]) return 1 end`
	delLockScript        = `local v2 = redis.call("get", KEYS[1]) if v2 then if v2 == ARGV[1] then redis.call("del", KEYS[1]) return 1 end return 0 end return -1`
	DelLockStatusNotOwn  = -1 // 非本人的锁【锁已过期，且已被抢占】
	DelLockStatusExpired = 0  // 锁已过期
	DelLockStatusSuccess = 1  // 删除成功
)

var (
	DelLockStatusNotOwnErr  = errors.Wrap(errors.New("锁已过期，且已被抢占"), "BaseRedisOp DelLock")
	DelLockStatusExpiredErr = errors.Wrap(errors.New("锁已过期"), "BaseRedisOp DelLock")
)

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

func (b *baseRedisOp) GetCli() redis.Cmdable {
	return b.redisCli
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
// ok=true写入成功，ok=false写入失败
func (b *baseRedisOp) SetNX(ctx context.Context, value string) (bool, error) {
	ok, err := b.redisCli.SetNX(ctx, b.key, value, b.ttl).Result()
	if err != nil {
		return ok, errors.Wrap(err, "BaseRedisOp SetNX")
	}
	return ok, nil
}

// 如果不存在，data为空字符串，不报错
func (b *baseRedisOp) Get(ctx context.Context) (string, error) {
	data, err := b.redisCli.Get(ctx, b.key).Result()
	if err != nil && !IsRedisNil(err) {
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
	if err != nil && !IsRedisNil(err) {
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
	if err != nil && !IsRedisNil(err) {
		return errors.Wrap(err, "BaseRedisOp Del")
	}
	return nil
}

// 分布式锁 写入
// value为uuid，ok为是否写入成功
func (b *baseRedisOp) SetLock(ctx context.Context) (value string, ok bool, err error) {
	value = uuid.NewV4().String()
	ok, err = b.redisCli.SetNX(ctx, b.key, value, b.ttl).Result()
	if err != nil {
		err = errors.Wrap(err, "BaseRedisOp SetLock")
		return
	}
	return
}

// 分布式锁 延长锁
func (b *baseRedisOp) ExtendLock(ctx context.Context, value string) (int64, error) {
	resI, err := b.redisCli.Eval(ctx, extendLockScript, []string{b.key}, value, formatMs(b.ttl)).Result()
	if err != nil {
		err = errors.Wrap(err, "BaseRedisOp DelLock")
		return 0, err
	}

	res, _ := resI.(int64)
	if res == DelLockStatusNotOwn {
		return res, DelLockStatusNotOwnErr
	}
	return res, nil
}

// 分布式锁 写入
// value为uuid，ok为是否写入成功
func (b *baseRedisOp) DelLock(ctx context.Context, value string) (int64, error) {
	resI, err := b.redisCli.Eval(ctx, delLockScript, []string{b.key}, value).Result()
	if err != nil {
		err = errors.Wrap(err, "BaseRedisOp DelLock")
		return 0, err
	}

	res, _ := resI.(int64)
	switch res {
	case DelLockStatusExpired:
		return res, DelLockStatusExpiredErr
	case DelLockStatusNotOwn:
		return res, DelLockStatusNotOwnErr
	}
	return res, nil
}

// 基数统计 写入
func (b *baseRedisOp) PFAdd(ctx context.Context, els ...interface{}) error {
	_, err := b.redisCli.PFAdd(ctx, b.key, els...).Result()
	if err != nil {
		err = errors.Wrap(err, "BaseRedisOp PfAdd")
		return err
	}
	return nil
}

// 基数统计 统计
// key不存在 不会redis.Nil 无需判断IsNil
func (b *baseRedisOp) PFCount(ctx context.Context) (count int64, err error) {
	count, err = b.redisCli.PFCount(ctx, b.key).Result()
	if err != nil {
		err = errors.Wrap(err, "BaseRedisOp PFCount")
	}
	return
}

// 基数统计 统计
// 把keys的数据 统计到b.key里
func (b *baseRedisOp) PFMerge(ctx context.Context, keys ...string) error {
	_, err := b.redisCli.PFMerge(ctx, b.key, keys...).Result()
	if err != nil {
		err = errors.Wrap(err, "BaseRedisOp PFMerge")
	}
	return err
}

func IsRedisNil(err error) bool {
	return errors.Is(err, redis.Nil)
}

func formatMs(dur time.Duration) int64 {
	if dur > 0 && dur < time.Millisecond {
		return 1
	}
	return int64(dur / time.Millisecond)
}
