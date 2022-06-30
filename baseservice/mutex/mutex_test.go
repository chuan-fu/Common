package mutex

import (
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

func TestDistributedOnce(t *testing.T) {
	d := NewDistributedOnce("key")
	for i := 0; i < 10; i++ {
		go func(i int) {
			d.Do(func() {
				fmt.Println(i, "run")
				time.Sleep(time.Second)
			})
		}(i)
	}
	time.Sleep(time.Second)

	for i := 0; i < 10; i++ {
		go func(i int) {
			d.Do(func() {
				fmt.Println(i, "run")
				time.Sleep(time.Second)
			})
		}(i)
	}

	time.Sleep(time.Second)
}
