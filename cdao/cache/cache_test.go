package cache

import (
	"fmt"
	"testing"
	"time"

	"github.com/chuan-fu/Common/cdao"

	"github.com/chuan-fu/Common/db/redis"
	log "github.com/chuan-fu/Common/zlog"
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

func TestCache(t *testing.T) {
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
