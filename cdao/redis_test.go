package cdao

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/chuan-fu/Common/db/redis"
	log "github.com/chuan-fu/Common/zlog"
)

func init() {
	err := redis.ConnectRedis(redis.RedisConf{
		Addr: "127.0.0.1:6379",
	})
	if err != nil {
		log.Fatal(err)
	}
}

func TestCache(t *testing.T) {
	b := NewBaseRedisOpWithKT("setnx:1", time.Minute)
	fmt.Println(b.SetLock(context.TODO()))
	// fmt.Println(b.SetLock(context.TODO()))
	// fmt.Println(b.SetLock(context.TODO()))
	fmt.Println(b.DelLock(context.TODO(), "34d815a8-7d60-4f83-8796-a15279be87f8"))
}
