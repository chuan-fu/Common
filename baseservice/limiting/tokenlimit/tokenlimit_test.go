package tokenlimit

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/chuan-fu/Common/util"

	"github.com/chuan-fu/Common/db/redis"
	log "github.com/chuan-fu/Common/zlog"

	"github.com/chuan-fu/Common/cdao"
)

func init() {
	err := redis.ConnectRedis(redis.RedisConf{
		Addr: "127.0.0.1:6379",
	})
	if err != nil {
		log.Fatal(err)
	}
}

func TestNewTokenLimiter(t *testing.T) {
	lim := NewTokenLimiter(10, 20, cdao.NewBaseRedisOp(), "token:limiter", WithPer(10*time.Second))
	fmt.Printf("%+v", lim)
	f := func() {
		for i := 0; i < 2000; i++ {
			time.Sleep(500 * time.Millisecond)
			fmt.Println(time.Now().Format(util.DefaultTimeFormat), i, lim.Allow(context.TODO()))
		}
	}

	go f()
	go f()
	go f()

	time.Sleep(time.Hour)
}
