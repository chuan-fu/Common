package hyperloglog

import (
	"context"
	"time"

	"github.com/chuan-fu/Common/baseservice/batch"
	"github.com/chuan-fu/Common/db"
	"github.com/chuan-fu/Common/zlog"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

/*
	基数统计功能，配置了异步写入、批量写入等功能
*/

type HyperLogLogOp interface {
	AsynAdd(data ...string) error
	Add(ctx context.Context, els ...interface{}) error
	Count(ctx context.Context) (count int64, err error)
	Merge(ctx context.Context, keys ...string) error
}

type Option func(h *hyperLogLog)

// buffLen参数设置请参考 baseservice/batch/string.go:21
func WithBatch(duration time.Duration, buffLen int) Option {
	return func(h *hyperLogLog) {
		h.batch = batch.NewStringIncrease(h.add, duration, buffLen)
	}
}

func WithStore(store redis.Cmdable) Option {
	return func(h *hyperLogLog) {
		h.store = store
	}
}

type hyperLogLog struct {
	store redis.Cmdable
	key   string
	batch *batch.StringIncrease
}

func NewHyperLogLog(key string, opts ...Option) HyperLogLogOp {
	h := &hyperLogLog{
		key: key,
	}
	for _, opt := range opts {
		opt(h)
	}
	if h.store == nil {
		h.store = db.GetRedisCli()
	}
	return h
}

// 如果WithBatch开了，不用处理error
func (h *hyperLogLog) AsynAdd(data ...string) (err error) {
	if h != nil {
		for k := range data {
			h.batch.Add(data[k])
		}
		return
	}

	if err = h.add(data); err != nil { // 如果异步没开，就转同步
		log.Error(err)
	}
	return
}

func (h *hyperLogLog) add(data []string) error {
	_, err := h.store.PFAdd(context.TODO(), h.key, data).Result()
	if err != nil {
		return errors.Wrap(err, "HyperLogLogOp Asyn Add")
	}
	return nil
}

// 写入
func (h *hyperLogLog) Add(ctx context.Context, els ...interface{}) error {
	_, err := h.store.PFAdd(ctx, h.key, els).Result()
	if err != nil {
		return errors.Wrap(err, "HyperLogLogOp PfAdd")
	}
	return nil
}

// 基数统计 统计
func (h *hyperLogLog) Count(ctx context.Context) (count int64, err error) {
	count, err = h.store.PFCount(ctx, h.key).Result()
	if err != nil {
		err = errors.Wrap(err, "HyperLogLogOp PFCount")
	}
	return
}

// 合并 把keys合并到h.key中
func (h *hyperLogLog) Merge(ctx context.Context, keys ...string) error {
	_, err := h.store.PFMerge(ctx, h.key, keys...).Result()
	if err != nil {
		err = errors.Wrap(err, "HyperLogLogOp PFMerge")
	}
	return err
}
