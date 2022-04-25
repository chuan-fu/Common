package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/chuan-fu/Common/db/redis"

	"github.com/chuan-fu/Common/cdao"
	"github.com/chuan-fu/Common/zlog"
)

func init() {
	err := redis.ConnectRedis(redis.RedisConf{
		Addr: "127.0.0.1:6379",
	})
	if err != nil {
		log.Fatal(err)
	}
}

func TestModelCache(t *testing.T) {
	op := cdao.NewBaseRedisOpWithKT("test:A:1", time.Minute)
	m := &AA{}

	data, err := GetBaseJsonCache(context.TODO(), op, m,
		func(ctx context.Context, model interface{}) (data string, err error) {
			log.Info("getByDB")
			err = json.Unmarshal([]byte(`{"id":1}`), model)
			if err != nil {
				log.Error(err)
			}
			return
		},
	)
	fmt.Println(data, err)
	fmt.Printf("%+v", m)
}