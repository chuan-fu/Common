package redislock

import (
	"context"
	"fmt"
	"testing"
	"time"

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

func TestNewRedisLock(t *testing.T) {
	f := func(id int) {
		r := NewRedisLock("redis:lock:1", time.Second*5)
		ok, err := r.SetLock(context.TODO())
		if err != nil {
			log.Error(err)
			return
		}
		if !ok {
			fmt.Println(id, "未抢到锁")
			return
		}
		fmt.Println(id, "抢到锁")
		time.Sleep(6 * time.Second)
		if err = r.ExtendLock(context.TODO()); err != nil {
			log.Error(err)
			return
		}
		fmt.Println(id, "延长锁成功")
		time.Sleep(6*time.Second + 500*time.Millisecond)

		if err = r.DelLock(context.TODO()); err != nil {
			log.Error(err)
			return
		}
		fmt.Println(id, "删除锁成功")
	}

	for i := 0; i < 20; i++ {
		time.Sleep(time.Second)
		go f(i)
	}

	time.Sleep(time.Minute)
}
