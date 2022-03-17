package cdao

import (
	"context"
	"encoding/json"
	"reflect"
	"strconv"
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
	Expire(ctx context.Context) (bool, error)
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
	SetModel(ctx context.Context, model interface{}) error
	GetModel(ctx context.Context, model interface{}) error
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

const (
	TagJson = "json"
	TagXml  = "xml"
	TagGorm = "gorm"
)

var (
	DelLockStatusNotOwnErr  = errors.Wrap(errors.New("锁已过期，且已被抢占"), "BaseRedisOp DelLock")
	DelLockStatusExpiredErr = errors.Wrap(errors.New("锁已过期"), "BaseRedisOp DelLock")
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

func (b *baseRedisOp) Expire(ctx context.Context) (bool, error) {
	ok, err := b.redisCli.Expire(ctx, b.key, b.ttl).Result()
	if err != nil && !IsRedisNil(err) {
		return false, errors.Wrap(err, "BaseRedisOp Expire")
	}
	return ok, nil
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
func (b *baseRedisOp) SetModel(ctx context.Context, model interface{}) error {
	rt := reflect.TypeOf(model)
	rv := reflect.ValueOf(model)
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
		rv = rv.Elem()
	}

	args := make([]interface{}, 0, 2*rt.NumField())
	for i := 0; i < rt.NumField(); i++ {
		args = append(args, getTag(rt, i, b.tag), rv.Field(i).Interface())
	}
	_, err := b.redisCli.HMSet(ctx, b.key, args...).Result()
	if err != nil {
		return errors.Wrap(err, "BaseRedisOp SetModel HMSet")
	}
	_, err = b.redisCli.Expire(ctx, b.key, b.ttl).Result()
	if err != nil {
		return errors.Wrap(err, "BaseRedisOp SetModel Expire")
	}
	return nil
}

// hash读取
func (b *baseRedisOp) GetModel(ctx context.Context, model interface{}) error {
	rt := reflect.TypeOf(model)
	if !(rt.Kind() == reflect.Ptr && rt.Elem().Kind() == reflect.Struct) {
		return errors.New("传入的model要为结构体指针")
	}

	rt = rt.Elem()
	args := make([]string, rt.NumField())
	for i := 0; i < rt.NumField(); i++ {
		args[i] = getTag(rt, i, b.tag)
	}
	values, err := b.redisCli.HMGet(ctx, b.key, args...).Result()
	if err != nil {
		return errors.Wrap(err, "BaseRedisOp SetModel HMSet")
	}
	if len(values) < rt.NumField() {
		return nil
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

func IsRedisNil(err error) bool {
	return errors.Is(err, redis.Nil)
}

func formatMs(dur time.Duration) int64 {
	if dur > 0 && dur < time.Millisecond {
		return 1
	}
	return int64(dur / time.Millisecond)
}

const (
	BitSize0  = 0
	BitSize8  = 8
	BitSize10 = 10
	BitSize16 = 16
	BitSize32 = 32
	BitSize64 = 64
)

func setReflectValueByStr(value reflect.Value, val string) error {
	switch value.Kind() {
	case reflect.Int:
		return setIntField(val, BitSize0, value)
	case reflect.Int8:
		return setIntField(val, BitSize8, value)
	case reflect.Int16:
		return setIntField(val, BitSize16, value)
	case reflect.Int32:
		return setIntField(val, BitSize32, value)
	case reflect.Int64:
		if _, ok := value.Interface().(time.Duration); ok {
			return setTimeDuration(val, value)
		}
		return setIntField(val, BitSize64, value)
	case reflect.Uint:
		return setUintField(val, BitSize0, value)
	case reflect.Uint8:
		return setUintField(val, BitSize8, value)
	case reflect.Uint16:
		return setUintField(val, BitSize16, value)
	case reflect.Uint32:
		return setUintField(val, BitSize32, value)
	case reflect.Uint64:
		return setUintField(val, BitSize64, value)
	case reflect.Bool:
		return setBoolField(val, value)
	case reflect.Float32:
		return setFloatField(val, BitSize32, value)
	case reflect.Float64:
		return setFloatField(val, BitSize64, value)
	case reflect.String:
		value.SetString(val)
	case reflect.Struct:
		if _, ok := value.Interface().(time.Time); ok {
			return setTimeField(val, value)
		}
		return json.Unmarshal(util.StringToBytes(val), value.Addr().Interface())
	case reflect.Map, reflect.Array:
		return json.Unmarshal(util.StringToBytes(val), value.Addr().Interface())
	case reflect.Slice:
		if _, ok := value.Interface().([]byte); ok {
			return setByteSlice(val, value)
		}
		return json.Unmarshal(util.StringToBytes(val), value.Addr().Interface())
	}
	return nil
}

func setIntField(val string, bitSize int, field reflect.Value) (err error) {
	var intVal int64
	if intVal, err = strconv.ParseInt(val, 10, bitSize); err != nil {
		return
	}
	field.SetInt(intVal)
	return
}

func setUintField(val string, bitSize int, field reflect.Value) (err error) {
	var intVal uint64
	if intVal, err = strconv.ParseUint(val, BitSize10, bitSize); err != nil {
		return
	}
	field.SetUint(intVal)
	return
}

func setBoolField(val string, field reflect.Value) (err error) {
	var boolVal bool
	if boolVal, err = strconv.ParseBool(val); err != nil {
		return
	}
	field.SetBool(boolVal)
	return
}

func setFloatField(val string, bitSize int, field reflect.Value) (err error) {
	var floatVal float64
	if floatVal, err = strconv.ParseFloat(val, bitSize); err != nil {
		return
	}
	field.SetFloat(floatVal)
	return
}

func setTimeDuration(val string, value reflect.Value) error {
	t, err := time.ParseDuration(val)
	if err != nil {
		return err
	}
	value.Set(reflect.ValueOf(t))
	return nil
}

func setTimeField(val string, value reflect.Value) error {
	t, err := time.Parse(time.RFC3339Nano, val)
	if err != nil {
		return err
	}
	value.Set(reflect.ValueOf(t))
	return nil
}

func setByteSlice(val string, value reflect.Value) error {
	value.Set(reflect.ValueOf(util.StringToBytes(val)))
	return nil
}

func getTag(field reflect.Type, i int, tag string) string {
	if tag != "" {
		if d := field.Field(i).Tag.Get(tag); d != "" && d != "-" {
			return d
		}
	}
	return field.Field(i).Name
}
