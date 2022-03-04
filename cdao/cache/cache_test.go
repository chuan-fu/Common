package cache

import (
	"fmt"
	"testing"
	"time"

	"github.com/chuan-fu/Common/util"

	"github.com/chuan-fu/Common/cdao"
	"github.com/chuan-fu/Common/db/redis"
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
	b := BaseStringCache{
		Op:    cdao.NewBaseRedisOpWithKT("test:A:1", time.Minute),
		Model: &AA{},
		GetByDB: func(model interface{}) (string, error) {
			log.Info("getByDB")
			return `{"id":1}`, nil
		},
	}

	data, err := b.Get()
	fmt.Println(data, err)
	fmt.Printf("%+v", b.Model)
}

func TestStrCache(t *testing.T) {
	b := BaseStringCache{
		Op: cdao.NewBaseRedisOpWithKT("test:A:2", time.Minute),
		GetByDB: func(model interface{}) (string, error) {
			log.Info("getByDB")
			return `234`, nil
		},
	}

	data, err := b.Get()
	fmt.Println(data, err, util.Type(data))
	fmt.Printf("%+v", b.Model)
}
