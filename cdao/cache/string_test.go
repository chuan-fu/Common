package cache

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/chuan-fu/Common/baseservice/syncx"

	"github.com/chuan-fu/Common/cdao"
	"github.com/chuan-fu/Common/zlog"
)

type AA struct {
	ID int64   `json:"id"`
	A  string  `json:"a"`
	B  int64   `json:"b"`
	C  float64 `json:"c"`
}

func TestStrCache(t *testing.T) {
	op := cdao.NewBaseRedisOp("cache:string:1", time.Minute)
	sg := syncx.NewSingleFlight()
	_ = sg

	ch := NewCacheHandle(sg)

	ctx := context.Background()
	f := func(ctx context.Context) (string, error) {
		log.Info("getByDB")
		return "A", nil
	}

	for i := 0; i < 100; i++ {
		go func() {
			data, err := ch.GetBaseStringCache(ctx, op, f)
			fmt.Println(data, err)
		}()
	}

	time.Sleep(3 * time.Second)
}
