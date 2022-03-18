package cdao

import (
	"context"
	"encoding/json"
	"reflect"
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
	SetTag(tag string) BaseRedisOp
	GetTag() string

	// string
	Set(ctx context.Context, value string) error
	SetNX(ctx context.Context, value string) (bool, error)
	Get(ctx context.Context) (string, error)
	GetResult(ctx context.Context, v interface{}) error

	// common
	Exists(ctx context.Context) (bool, error)
	Expire(ctx context.Context) (bool, error)
	TTL(ctx context.Context) (time.Duration, error)
	Del(ctx context.Context) error

	// 分布式锁
	SetLock(ctx context.Context) (value string, ok bool, err error)
	ExtendLock(ctx context.Context, value string) (int64, error)
	DelLock(ctx context.Context, value string) (int64, error)

	// 基数统计
	PFAdd(ctx context.Context, els ...interface{}) error
	PFCount(ctx context.Context) (count int64, err error)
	PFMerge(ctx context.Context, keys ...string) error

	// hash
	HSetModel(ctx context.Context, model interface{}) error
	HGetModel(ctx context.Context, model interface{}) error
	HSetMap(ctx context.Context, m map[string]interface{}) error
	HGetMap(ctx context.Context, m map[string]string) error
	HGet(ctx context.Context, key string) (val string, err error)
	HGetAll(ctx context.Context) (map[string]string, error)

	// zset
	ZAddString(ctx context.Context, list []string) error
	ZRangeString(ctx context.Context, start, stop int64) (data []string, err error)                 // 根据下标
	ZRangeStringWithPage(ctx context.Context, pageIndex, pageSize int64) (data []string, err error) // 根据分页
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

/* zaddScript 覆盖式写入
if redis.call("exists", KEYS[1]) == 1 then
	redis.call("del", KEYS[1])
end
redis.call("zadd", KEYS[1], unpack(ARGV))
if KEYS[2] then
	redis.call("expire", KEYS[1], KEYS[2])
end
return 1
*/
const (
	extendLockScript     = `local v2 = redis.call("get", KEYS[1]) if v2 then if v2 == ARGV[1] then redis.call("pexpire", KEYS[1], ARGV[2]) return 1 end return -1 else redis.call("set", KEYS[1], ARGV[1], "px", ARGV[2]) return 1 end`
	delLockScript        = `local v2 = redis.call("get", KEYS[1]) if v2 then if v2 == ARGV[1] then redis.call("del", KEYS[1]) return 1 end return 0 end return -1`
	DelLockStatusNotOwn  = -1 // 非本人的锁【锁已过期，且已被抢占】
	DelLockStatusExpired = 0  // 锁已过期
	DelLockStatusSuccess = 1  // 删除成功

	zaddScript = `if redis.call("exists", KEYS[1]) == 1 then redis.call("del", KEYS[1]) end redis.call("zadd", KEYS[1], unpack(ARGV)) if KEYS[2] then redis.call("expire", KEYS[1], KEYS[2]) end return 1`
)

const (
	TagJson = "json"
	TagXml  = "xml"
	TagGorm = "gorm"
)

var (
	DelLockStatusNotOwnErr  = errors.New("BaseRedisOp DelLock: 锁已过期，且已被抢占")
	DelLockStatusExpiredErr = errors.New("BaseRedisOp DelLock: 锁已过期")
)

type baseRedisOp struct {
	redisCli redis.Cmdable
	key      string
	ttl      time.Duration
	tag      string
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

func (b *baseRedisOp) GetTag() string {
	return b.tag
}

func (b *baseRedisOp) SetTag(tag string) BaseRedisOp {
	b.tag = tag
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

func (b *baseRedisOp) Exists(ctx context.Context) (bool, error) {
	i, err := b.redisCli.Exists(ctx, b.key).Result()
	if err != nil && !IsRedisNil(err) {
		return false, errors.Wrap(err, "BaseRedisOp Expire")
	}
	return i == 1, nil
}

func (b *baseRedisOp) Expire(ctx context.Context) (bool, error) {
	ok, err := b.redisCli.Expire(ctx, b.key, b.ttl).Result()
	if err != nil && !IsRedisNil(err) {
		return false, errors.Wrap(err, "BaseRedisOp Expire")
	}
	return ok, nil
}

func (b *baseRedisOp) TTL(ctx context.Context) (time.Duration, error) {
	t, err := b.redisCli.TTL(ctx, b.key).Result()
	if err != nil && !IsRedisNil(err) {
		return 0, errors.Wrap(err, "BaseRedisOp Expire")
	}
	return t, nil
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
		return errors.Wrap(err, "BaseRedisOp PfAdd")
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

// hash写入
func (b *baseRedisOp) HSetModel(ctx context.Context, model interface{}) error {
	rt := reflect.TypeOf(model)
	rv := reflect.ValueOf(model)
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
		rv = rv.Elem()
	}
	if rt.Kind() != reflect.Struct {
		return errors.New("传入的model要为结构体或结构体指针")
	}

	args := make([]interface{}, 0, 2*rt.NumField())
	for i := 0; i < rt.NumField(); i++ {
		args = append(args, getTag(rt, i, b.tag), util.ToString(rv.Field(i).Interface()))
	}

	if _, err := b.redisCli.HMSet(ctx, b.key, args...).Result(); err != nil {
		return errors.Wrap(err, "BaseRedisOp HSetModel HMSet")
	}
	if b.ttl > 0 {
		if _, err := b.redisCli.Expire(ctx, b.key, b.ttl).Result(); err != nil {
			return errors.Wrap(err, "BaseRedisOp HSetModel Expire")
		}
	}
	return nil
}

// hash读取
func (b *baseRedisOp) HGetModel(ctx context.Context, model interface{}) error {
	if !util.IsPtrStruct(model) {
		return errors.New("传入的model要为结构体指针")
	}

	rt := reflect.TypeOf(model).Elem()
	args := make([]string, rt.NumField())
	for i := 0; i < rt.NumField(); i++ {
		args[i] = getTag(rt, i, b.tag)
	}
	values, err := b.redisCli.HMGet(ctx, b.key, args...).Result()
	if err != nil {
		return errors.Wrap(err, "BaseRedisOp HGetModel HMSet")
	}

	rv := reflect.ValueOf(model).Elem()
	for i := 0; i < rv.NumField(); i++ {
		if values[i] == nil {
			continue
		}
		vStr, _ := values[i].(string)
		if err = setReflectValueByStr(rv.Field(i), vStr); err != nil {
			return err
		}
	}
	return nil
}

// hash map写入
func (b *baseRedisOp) HSetMap(ctx context.Context, m map[string]interface{}) error {
	args := make([]interface{}, 0, 2*len(m))
	for k, v := range m {
		args = append(args, k, util.ToString(v))
	}
	if _, err := b.redisCli.HMSet(ctx, b.key, args...).Result(); err != nil {
		return errors.Wrap(err, "BaseRedisOp HSetMap HMSet")
	}
	if b.ttl > 0 {
		if _, err := b.redisCli.Expire(ctx, b.key, b.ttl).Result(); err != nil {
			return errors.Wrap(err, "BaseRedisOp HSetMap Expire")
		}
	}
	return nil
}

// hash map读取
func (b *baseRedisOp) HGetMap(ctx context.Context, m map[string]string) error {
	args := make([]string, 0, len(m))
	for k := range m {
		args = append(args, k)
	}
	values, err := b.redisCli.HMGet(ctx, b.key, args...).Result()
	if err != nil {
		return errors.Wrap(err, "BaseRedisOp HGetMap HMSet")
	}
	for k := range args {
		if values[k] != nil {
			m[args[k]], _ = values[k].(string)
		}
	}
	return nil
}

// hash map读取
func (b *baseRedisOp) HGet(ctx context.Context, key string) (string, error) {
	val, err := b.redisCli.HGet(ctx, b.key, key).Result()
	if err != nil && !IsRedisNil(err) {
		return "", errors.Wrap(err, "BaseRedisOp HGet")
	}
	return val, nil
}

// hash map读取
func (b *baseRedisOp) HGetAll(ctx context.Context) (map[string]string, error) {
	m, err := b.redisCli.HGetAll(ctx, b.key).Result()
	if err != nil {
		return nil, errors.Wrap(err, "BaseRedisOp HGetAll")
	}
	return m, nil
}

// 覆盖式写入，并设置ttl
func (b *baseRedisOp) ZAddString(ctx context.Context, list []string) error {
	if len(list) == 0 {
		return nil
	}

	keys := []string{b.key}
	if t := formatSec(b.ttl); t > 0 {
		keys = append(keys, util.ToString(t))
	}
	args := make([]interface{}, 0, 2*len(list))
	for k := range list {
		args = append(args, k, list[k])
	}
	_, err := b.redisCli.Eval(ctx, zaddScript, keys, args...).Result()
	if err != nil {
		return errors.Wrap(err, "BaseRedisOp ZAdd")
	}
	return nil
}

// 获取列表数据
func (b *baseRedisOp) ZRangeString(ctx context.Context, start, stop int64) (data []string, err error) {
	data, err = b.redisCli.ZRange(ctx, b.key, start, stop).Result()
	if err != nil {
		err = errors.Wrap(err, "BaseRedisOp ZRange")
	}
	return
}

// 获取列表数据
func (b *baseRedisOp) ZRangeStringWithPage(ctx context.Context, pageIndex, pageSize int64) (data []string, err error) {
	start := (pageIndex - 1) * pageSize
	end := start + pageSize - 1

	data, err = b.redisCli.ZRange(ctx, b.key, start, end).Result()
	if err != nil {
		err = errors.Wrap(err, "BaseRedisOp ZRange")
	}
	return
}

func IsRedisNil(err error) bool {
	return errors.Is(err, redis.Nil)
}
