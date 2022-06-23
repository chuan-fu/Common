package cdao

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/chuan-fu/Common/baseservice/cast"

	"github.com/chuan-fu/Common/db/redis"
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

type AA struct {
	A int64                  `json:"a"`
	B string                 `json:"b"`
	C []byte                 `json:"c"`
	D bool                   `json:"d"`
	E map[string]interface{} `json:"e"`
	F time.Duration
	G time.Time
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
	err := NewBaseRedisOp("HSetMap", time.Hour).HSetMap(context.TODO(), m)
	if err != nil {
		log.Error(err)
		return
	}
	fmt.Println(m)

	fmt.Println(NewBaseRedisOp("HSetMap", time.Hour).HGetMap(context.TODO()))
}

func TestAA(t *testing.T) {
	b := "av"
	// b := []byte{'a', 'b'}
	// b := []string{"v", "a"}

	bb := make([]byte, 0)
	err := json.Unmarshal([]byte(b), &bb)
	if err != nil {
		log.Error(err)
		return
	}
	fmt.Println(bb)
	fmt.Println(cast.ToString(bb))
}

func TestTime(t *testing.T) {
	vt := time.Second
	v := fmt.Sprintf("%dns", uint64(vt))
	fmt.Println(v)
}

func TestTTL(t *testing.T) {
	data, err := NewBaseRedisOp("key:userId:1", time.Hour).TTL(context.TODO())
	if err != nil {
		log.Error(err)
		return
	}
	fmt.Println(data)
	fmt.Println(int64(data / time.Second))
}

func TestSetBits(t *testing.T) {
	err := NewBaseRedisOp("key:setbits:1", time.Hour).SetBits(context.TODO(), []int64{1, 2, 3, 5, 7, -1, 10000, 1231})
	if err != nil {
		log.Error(err)
		return
	}
}

func TestGetBits(t *testing.T) {
	fmt.Println(NewBaseRedisOp("key:setbits:1", time.Hour).GetBits(context.TODO(), []int64{-1, -2, 0, 1, 2, 3, 4, 5, 6, 7, 1231, 1111}))
}
