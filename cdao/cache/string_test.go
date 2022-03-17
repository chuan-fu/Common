package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/chuan-fu/Common/cdao"
	"github.com/chuan-fu/Common/db/redis"
	"github.com/chuan-fu/Common/util"
	"github.com/chuan-fu/Common/zlog"
)

type AA struct {
	ID int64   `json:"id"`
	A  string  `json:"a"`
	B  int64   `json:"b"`
	C  float64 `json:"c"`
}

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

	data, err := GetBaseStringCache(context.TODO(), op,
		func(ctx context.Context, model interface{}) (data string, err error) {
			log.Info("getByDB")
			err = json.Unmarshal([]byte(`{"id":1}`), model)
			if err != nil {
				log.Error(err)
			}
			return
		},
		WithSetModel(m),
	)
	fmt.Println(data, err)
	fmt.Printf("%+v", m)
}

func TestStrCache(t *testing.T) {
	op := cdao.NewBaseRedisOpWithKT("test:A:2", time.Minute)
	data, err := GetBaseStringCache(context.TODO(), op,
		func(ctx context.Context, model interface{}) (string, error) {
			log.Info("getByDB")
			return `234`, nil
		},
	)
	fmt.Println(data, err, util.Type(data))
}
