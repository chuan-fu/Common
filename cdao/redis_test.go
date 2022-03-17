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

type AA struct {
	A   int64  `json:"a"`
	AAA int64  `json:"a"`
	B   string `json:"b"`
	C   []byte `json:"c"`
	D   bool   `json:"d"`
}

func TestCache(t *testing.T) {
	b := NewBaseRedisOpWithKT("setnx:1", time.Minute)
	// fmt.Println(b.SetLock(context.TODO()))
	// fmt.Println(b.SetLock(context.TODO()))
	// fmt.Println(b.SetLock(context.TODO()))
	fmt.Println(b.ExtendLock(context.TODO(), "34d815a8-7d60-4f83-8796-a15279be87f8"))
}

func TestSetModel(t *testing.T) {
	b := NewBaseRedisOpWithKT("GetModel", time.Hour).SetTag("json")
	err := b.SetModel(context.TODO(), &AA{
		A: 1,
		B: "啊",
		C: []byte{'a', 'v'},
		D: true,
	})
	fmt.Println(string([]byte{'a', 'v'}))
	fmt.Println(err)
}

func TestGetModel(t *testing.T) {
	b := NewBaseRedisOpWithKT("GetModel", time.Minute).SetTag("json")
	a := &AA{}
	err := b.GetModel(context.TODO(), a)
	fmt.Println(err)
	fmt.Println(string(a.C))
	fmt.Printf("%+v", *a)
}
