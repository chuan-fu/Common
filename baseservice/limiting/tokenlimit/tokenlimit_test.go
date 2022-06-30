package tokenlimit

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/chuan-fu/Common/baseservice/timex"

	"github.com/chuan-fu/Common/db"
	"github.com/chuan-fu/Common/db/redis"
	"github.com/chuan-fu/Common/zlog"
	"go.uber.org/goleak"
)

func init() {
	err := redis.ConnectRedis(redis.RedisConf{
		Addr: []string{"127.0.0.1:6379"},
	})
	if err != nil {
		log.Fatal(err)
	}
}

func TestNewTokenLimiter(t *testing.T) {
	defer goleak.VerifyNone(t)
	lim := NewTokenLimiter(10, 20, db.GetRedisCli(), "token:limiter", WithPer(10*time.Second))
	fmt.Printf("%+v", lim)
	f := func() {
		for i := 0; i < 2000; i++ {
			time.Sleep(500 * time.Millisecond)
			fmt.Println(timex.NewNow().FormatTime(), i, lim.Allow(context.TODO()))
		}
	}

	go f()
	go f()
	go f()

	time.Sleep(time.Hour)
}
