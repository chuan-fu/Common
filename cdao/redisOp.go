package cdao

import (
	"context"
	"time"

	"github.com/chuan-fu/Common/baseservice/jsonx"

	"github.com/chuan-fu/Common/baseservice/cast"
	"github.com/chuan-fu/Common/db"
	"github.com/chuan-fu/Common/zlog"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

type RedisOption func(h *baseRedisOp)

func WithStore(store redis.Cmdable) RedisOption {
	return func(h *baseRedisOp) {
		h.store = store
	}
}

type baseRedisOp struct {
	store redis.Cmdable
	key   string
	ttl   time.Duration
}

func NewBaseRedisOp(key string, ttl time.Duration, opts ...RedisOption) BaseRedisOp {
	r := &baseRedisOp{
		key: key,
		ttl: ttl,
	}
	for _, opt := range opts {
		opt(r)
	}
	if r.store == nil {
		r.store = db.GetRedisCli()
	}
	return r
}

func (b *baseRedisOp) SetKey(key string) {
	b.key = key
}

func (b *baseRedisOp) GetKey() string {
	return b.key
}

func (b *baseRedisOp) SetTTL(ttl time.Duration) {
	b.ttl = ttl
}

func (b *baseRedisOp) GetTTL() time.Duration {
	return b.ttl
}

func (b *baseRedisOp) Set(ctx context.Context, value string) error {
	_, err := b.store.Set(ctx, b.key, value, b.ttl).Result()
	if err != nil {
		return errors.Wrap(err, "BaseRedisOp Set")
	}
	return nil
}

// 不存在则写入
// ok=true写入成功，ok=false写入失败
func (b *baseRedisOp) SetNX(ctx context.Context, value string) (bool, error) {
	ok, err := b.store.SetNX(ctx, b.key, value, b.ttl).Result()
	if err != nil {
		return ok, errors.Wrap(err, "BaseRedisOp SetNX")
	}
	return ok, nil
}

// 如果不存在，data为空字符串，不报错
func (b *baseRedisOp) Get(ctx context.Context) (string, error) {
	data, err := b.store.Get(ctx, b.key).Result()
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
	data, err := b.store.Get(ctx, b.key).Result()
	if err != nil && !IsRedisNil(err) {
		return errors.Wrap(err, "BaseRedisOp GetResult")
	}
	if data == "" {
		return nil
	}
	log.Debugf("cache【%s】存在", b.key)
	return jsonx.Unmarshal(data, v)
}

func (b *baseRedisOp) Exists(ctx context.Context) (bool, error) {
	i, err := b.store.Exists(ctx, b.key).Result()
	if err != nil && !IsRedisNil(err) {
		return false, errors.Wrap(err, "BaseRedisOp Expire")
	}
	return i == 1, nil
}

func (b *baseRedisOp) Expire(ctx context.Context) (bool, error) {
	ok, err := b.store.Expire(ctx, b.key, b.ttl).Result()
	if err != nil && !IsRedisNil(err) {
		return false, errors.Wrap(err, "BaseRedisOp Expire")
	}
	return ok, nil
}

func (b *baseRedisOp) TTL(ctx context.Context) (time.Duration, error) {
	t, err := b.store.TTL(ctx, b.key).Result()
	if err != nil && !IsRedisNil(err) {
		return 0, errors.Wrap(err, "BaseRedisOp Expire")
	}
	return t, nil
}

func (b *baseRedisOp) Del(ctx context.Context) error {
	_, err := b.store.Del(ctx, b.key).Result()
	if err != nil && !IsRedisNil(err) {
		return errors.Wrap(err, "BaseRedisOp Del")
	}
	return nil
}

// 覆盖式写入zset，并设置ttl
func (b *baseRedisOp) ZAddCoverStringList(ctx context.Context, list []string) error {
	if len(list) == 0 {
		return nil
	}

	keys := []string{b.key}
	if t := formatSec(b.ttl); t > 0 {
		keys = append(keys, cast.ToString(t))
	}
	args := make([]interface{}, 0, 2*len(list))
	for k := range list {
		args = append(args, k, list[k])
	}
	_, err := b.store.Eval(ctx, zaddScript, keys, args...).Result()
	if err != nil {
		return errors.Wrap(err, "BaseRedisOp ZAdd")
	}
	return nil
}

// zset获取所有
func (b *baseRedisOp) ZGetAll(ctx context.Context) ([]string, error) {
	dataInter, err := b.store.Eval(ctx, zgetallScript, []string{b.key}).Result()
	if err != nil && !IsRedisNil(err) {
		err = errors.Wrap(err, "BaseRedisOp ZGetAll")
		return nil, err
	}
	if dataInter == nil {
		return nil, nil
	}

	dataList, _ := dataInter.([]interface{})
	data := make([]string, 0, len(dataList))
	for k := range dataList {
		data = append(data, cast.ToString(dataList[k]))
	}
	return data, nil
}

// 获取列表数据
func (b *baseRedisOp) ZRangeStringList(ctx context.Context, start, stop int64) (data []string, err error) {
	data, err = b.store.ZRange(ctx, b.key, start, stop).Result()
	if err != nil {
		err = errors.Wrap(err, "BaseRedisOp ZRange")
	}
	return
}

// 获取列表数据
func (b *baseRedisOp) ZRangeStringListWithPage(ctx context.Context, pageIndex, pageSize int64) (data []string, err error) {
	start := (pageIndex - 1) * pageSize
	end := start + pageSize - 1

	data, err = b.store.ZRange(ctx, b.key, start, end).Result()
	if err != nil {
		err = errors.Wrap(err, "BaseRedisOp ZRange")
	}
	return
}

// 覆盖式写入set
func (b *baseRedisOp) SAddCover(ctx context.Context, list []string) error {
	if len(list) == 0 {
		return nil
	}
	keys := []string{b.key}
	if t := formatSec(b.ttl); t > 0 {
		keys = append(keys, cast.ToString(t))
	}
	args := make([]interface{}, len(list))
	for k := range list {
		args[k] = list[k]
	}
	_, err := b.store.Eval(ctx, saddScript, keys, args...).Result()
	if err != nil {
		return errors.Wrap(err, "BaseRedisOp ZAdd")
	}
	return nil
}

func (b *baseRedisOp) SGetAll(ctx context.Context) (data []string, err error) {
	data, err = b.store.SMembers(ctx, b.key).Result()
	if err != nil {
		err = errors.Wrap(err, "BaseRedisOp SGetAll")
	}
	return
}

func (r *baseRedisOp) HSetMap(ctx context.Context, m map[string]interface{}) error {
	if len(m) == 0 {
		return errors.New("map为空")
	}

	keys := []string{r.key}
	if t := formatSec(r.ttl); t > 0 {
		keys = append(keys, cast.ToString(t))
	}
	args := make([]interface{}, 0, 2*len(m))
	for k, v := range m {
		args = append(args, k, cast.ToString(v))
	}
	if _, err := r.store.Eval(ctx, hmsetScript, keys, args...).Result(); err != nil && !IsRedisNil(err) {
		return errors.Wrap(err, "BaseRedisOp HSetMap")
	}
	return nil
}

func (r *baseRedisOp) HGetMap(ctx context.Context) (map[string]string, error) {
	data, err := r.store.HGetAll(ctx, r.key).Result()
	if err != nil && !IsRedisNil(err) {
		return nil, errors.Wrap(err, "BaseRedisOp HGetMap")
	}
	return data, nil
}

func (r *baseRedisOp) HGet(ctx context.Context, key string) (string, error) {
	val, err := r.store.HGet(ctx, r.key, key).Result()
	if err != nil && !IsRedisNil(err) {
		return val, errors.Wrap(err, "BaseRedisOp HGet")
	}
	return val, nil
}

func (b *baseRedisOp) SetBits(ctx context.Context, value []int64) error {
	keys := []string{b.key}
	if t := formatSec(b.ttl); t > 0 {
		keys = append(keys, cast.ToString(t))
	}
	args := make([]interface{}, 0, len(value))
	for _, v := range value {
		if v >= 0 { // 下标不能为负数，数组越界
			args = append(args, cast.ToString(v))
		}
	}
	_, err := b.store.Eval(ctx, setBitsScript, keys, args...).Result()
	if err != nil {
		return errors.Wrap(err, "BaseRedisOp SetBits")
	}
	return nil
}

func (b *baseRedisOp) GetBits(ctx context.Context, value []int64) (resp map[int64]struct{}, exists bool, err error) {
	resp = make(map[int64]struct{})

	argsInt64 := make([]int64, 0, len(value))
	args := make([]interface{}, 0, len(value))
	for _, v := range value {
		if v >= 0 { // 下标不能为负数，数组越界
			args = append(args, cast.ToString(v))
			argsInt64 = append(argsInt64, v)
		}
	}
	if len(args) == 0 {
		err = errors.New("value有误")
		return
	}

	dataInter, err := b.store.Eval(ctx, getBitsScript, []string{b.key}, args...).Result()
	if err != nil && !IsRedisNil(err) {
		log.Error(err)
		return
	}
	if IsRedisNil(err) {
		return resp, false, nil
	}

	data, _ := dataInter.([]interface{})
	for k, vInter := range data {
		v, _ := vInter.(int64)
		if v == 1 {
			resp[argsInt64[k]] = struct{}{}
		}
	}
	return resp, true, nil
}

func IsRedisNil(err error) bool {
	return errors.Is(err, redis.Nil)
}
