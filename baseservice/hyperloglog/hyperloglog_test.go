package hyperloglog

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/chuan-fu/Common/baseservice/cast"

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

func TestAdd(t *testing.T) {
	h := NewHyperLogLog("pf:add:5", WithBatch(time.Second, 20))
	for i := 1; i <= 1000; i++ {
		time.Sleep(time.Millisecond * time.Duration(i%10))
		_ = h.AsynAdd(cast.ToString(i))
	}
	time.Sleep(time.Second)
	fmt.Println(h.Count(context.TODO()))
}
