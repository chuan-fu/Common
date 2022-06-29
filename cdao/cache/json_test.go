package cache

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/chuan-fu/Common/baseservice/syncx"

	"github.com/chuan-fu/Common/baseservice/jsonx"
	"github.com/chuan-fu/Common/cdao"
	"github.com/chuan-fu/Common/db/redis"
	"github.com/chuan-fu/Common/zlog"
)

func init() {
	err := redis.ConnectRedis(redis.RedisConf{
		Addr: []string{"127.0.0.1:6379"},
	})
	if err != nil {
		log.Fatal(err)
	}
}

func TestModelCache(t *testing.T) {
	op := cdao.NewBaseRedisOp("cache:json:1", time.Minute)
	sg := syncx.NewSingleFlight()
	_ = sg

	ch := NewCacheHandle(sg)
	// ch := NewCacheHandle(nil)

	ctx := context.Background()
	f := func(ctx context.Context, model interface{}) (err error) {
		log.Info("getByDB")
		err = jsonx.Unmarshal(`{"id":1}`, model)
		if err != nil {
			log.Error(err)
		}
		return
	}

	for i := 0; i < 100; i++ {
		go func() {
			m := &AA{}
			data, err := ch.GetBaseJsonCache(ctx, op, m, f)
			fmt.Println(data, err, fmt.Sprintf("%+v", m))
		}()
	}

	time.Sleep(3 * time.Second)
}
