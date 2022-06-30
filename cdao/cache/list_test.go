package cache

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/chuan-fu/Common/cdao"
	"github.com/chuan-fu/Common/zlog"
)

func A(obj interface{}) {
	if reflect.TypeOf(obj).Kind() == reflect.Slice {
		s := reflect.ValueOf(obj)
		for i := 0; i < s.Len(); i++ {
			ele := s.Index(i)
			fmt.Println(ele.Interface())
		}
	}
}

func TestGetBaseListCache(t *testing.T) {
	op := cdao.NewBaseRedisOp("cache:zset:1", time.Minute)
	data, err := C.GetBaseListCache(context.TODO(), op, func(ctx context.Context) ([]string, error) {
		log.Info("getByDb")
		return []string{"A1", "A2", "A3", "A4", "A5", "A6", "A7", "A8", "A9", "A10", "A11", "A12"}, nil
	})
	fmt.Println(data, err)
}

func TestGetBaseListCache3(t *testing.T) {
	op := cdao.NewBaseRedisOp("cache:zset:1", time.Minute)
	data, err := C.GetBaseListCache(context.TODO(), op, func(ctx context.Context) ([]string, error) {
		log.Info("getByDb")
		return []string{"AA", "BB", "CC", "D", "EE", "F"}, nil
	})
	fmt.Println(data, err)
}

func TestGetBaseListCacheEntry(t *testing.T) {
	op := cdao.NewBaseRedisOp("cache:zset:entry", time.Minute)
	data, err := C.GetBaseListCache(context.TODO(), op, func(ctx context.Context) ([]string, error) {
		log.Info("getByDb")
		return []string{}, nil
	})
	fmt.Println(data, err)
}

func TestGetBaseListCacheWithPage(t *testing.T) {
	op := cdao.NewBaseRedisOp("cache:zset:1", time.Minute)
	data, err := C.GetBaseListCacheWithPage(context.TODO(), op, func(ctx context.Context) ([]string, error) {
		log.Info("getByDb")
		return []string{"AA", "BB", "CC", "D", "EE", "F"}, nil
	}, 1, 3)
	fmt.Println(data, err)
}

func TestGetBaseListCacheWithPageEntry(t *testing.T) {
	op := cdao.NewBaseRedisOp("cache:zset:entry", time.Minute)
	data, err := C.GetBaseListCacheWithPage(context.TODO(), op, func(ctx context.Context) ([]string, error) {
		log.Info("getByDb")
		return []string{}, nil
	}, 1, 3)
	fmt.Println(data, err)
}
