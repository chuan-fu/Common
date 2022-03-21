package cache

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	log "github.com/chuan-fu/Common/zlog"

	"github.com/chuan-fu/Common/cdao"
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
	op := cdao.NewBaseRedisOpWithKT("list:A:1", time.Minute)
	data, err := GetBaseListCache(context.TODO(), op, func(ctx context.Context) (interface{}, error) {
		log.Info("getByDb")
		data := []string{"A1", "A2", "A3", "A4", "A5", "A6", "A7", "A8", "A9", "A10", "A11", "A12"}
		return data, nil
	})
	fmt.Println(data, err)
}

func TestGetBaseListCache2(t *testing.T) {
	op := cdao.NewBaseRedisOpWithKT("list:A:2", time.Minute)
	data, err := GetBaseListCache(context.TODO(), op, func(ctx context.Context) (interface{}, error) {
		log.Info("getByDb")
		data := []AA{{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5}, {ID: 6}, {ID: 7}, {ID: 8}, {ID: 9}, {ID: 10}, {ID: 11}, {ID: 12}, {ID: 13}}
		return data, nil
	}, WithSetListPage(1, 10))
	fmt.Println(data, err)
}

func TestGetBaseListCache3(t *testing.T) {
	op := cdao.NewBaseRedisOpWithKT("list:A:3", time.Minute)
	data, err := GetBaseListCache(context.TODO(), op, func(ctx context.Context) (interface{}, error) {
		log.Info("getByDb")
		data := []AA{{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5}, {ID: 6}, {ID: 7}, {ID: 8}, {ID: 9}, {ID: 10}, {ID: 11}, {ID: 12}, {ID: 13}}
		return &data, nil
	}, WithSetListPage(1, 10))
	fmt.Println(data, err)
}

func TestGetBaseListCache4(t *testing.T) {
	op := cdao.NewBaseRedisOpWithKT("list:A:4", time.Minute)
	data, err := GetBaseListCache(context.TODO(), op, func(ctx context.Context) (interface{}, error) {
		log.Info("getByDb")
		return []func(){func() {}, func() {}}, nil
	})
	fmt.Println(data, err)
}

func TestGetBaseListCache5(t *testing.T) {
	op := cdao.NewBaseRedisOpWithKT("list:A:5", time.Minute)
	data, err := GetBaseListCache(context.TODO(), op, func(ctx context.Context) (interface{}, error) {
		log.Info("getByDb")
		return []interface{}{1}, nil
	})
	fmt.Println(data, err)
}
