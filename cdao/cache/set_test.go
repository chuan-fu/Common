package cache

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/chuan-fu/Common/cdao"
	"github.com/chuan-fu/Common/zlog"
)

func TestGetBaseSetCache(t *testing.T) {
	op := cdao.NewBaseRedisOp("cache:set:1", time.Minute)
	data, err := C.GetBaseSetCache(context.TODO(), op, func(ctx context.Context) ([]string, error) {
		log.Info("getByDb")
		return []string{"A", "b"}, nil
	})
	fmt.Println(data, err)
}

func TestGetBaseSetCacheEntry(t *testing.T) {
	op := cdao.NewBaseRedisOp("cache:set:2", time.Minute)
	data, err := C.GetBaseSetCache(context.TODO(), op, func(ctx context.Context) ([]string, error) {
		log.Info("getByDb")
		return []string{}, nil
	})
	fmt.Println(data, err)
}
