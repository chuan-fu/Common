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
	A int64                  `json:"a"`
	B string                 `json:"b"`
	C []byte                 `json:"c"`
	D bool                   `json:"d"`
	E map[string]interface{} `json:"e"`
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
	err := b.HSetModel(context.TODO(), &AA{
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
	err := b.HGetModel(context.TODO(), a)
	fmt.Println(err)
	fmt.Println(string(a.C))
	fmt.Printf("%+v", *a)
}

func TestTTL(t *testing.T) {
	ttl, err := NewBaseRedisOpWithKT("GetModel", time.Minute).TTL(context.TODO())
	if err != nil {
		log.Error(err)
		return
	}

	fmt.Println(ttl)
}

func TestZAddString(t *testing.T) {
	err := NewBaseRedisOp().SetKey("ZAdd").SetTTL(time.Hour).ZAddString(context.TODO(), []string{"A1", "A2", "A3", "A4", "A5", "A6", "A7", "A8", "A9", "A10", "A11", "A12"})
	if err != nil {
		log.Error(err)
		return
	}
}

func TestZRangeString(t *testing.T) {
	data, err := NewBaseRedisOp().SetKey("ZAdd").ZRangeString(context.TODO(), 1, 10)
	if err != nil {
		log.Error(err)
		return
	}
	fmt.Println(data)
}

func TestHSetMap(t *testing.T) {
	m := map[string]interface{}{
		"A": []string{"1a", "2b"},
		"B": map[string]string{"3c": "3c-1", "4d": "4d-2"},
		"C": AA{A: 1, B: "2", C: []byte{'a', 'v'}, D: true},
		"D": 1,
		"E": 1.2,
		"F": "F11",
	}
	err := NewBaseRedisOpWithKT("HSetMap", time.Hour).HSetMap(context.TODO(), m)
	if err != nil {
		log.Error(err)
		return
	}
	fmt.Println(m)
}

func TestHGetAll(t *testing.T) {
	m, err := NewBaseRedisOpWithKT("HSetMap2", time.Minute).HGetAll(context.TODO())
	if err != nil {
		log.Error(err)
		return
	}
	fmt.Println(m)
}

func TestHSetModel(t *testing.T) {
	err := NewBaseRedisOpWithKT("HSetModel", time.Minute).HSetModel(context.TODO(), AA{
		A: 0,
		B: "",
		C: nil,
		D: false,
		E: map[string]interface{}{
			"E1": "1",
		},
	})
	if err != nil {
		log.Error(err)
		return
	}
}
