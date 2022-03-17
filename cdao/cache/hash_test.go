package cache

import (
	"context"
	"fmt"
	"testing"
	"time"

	log "github.com/chuan-fu/Common/zlog"

	"github.com/chuan-fu/Common/cdao"
)

type BB struct {
	ID int64   `json:"id"`
	A  string  `json:"a"`
	B  int64   `json:"b"`
	C  float64 `json:"c"`
}

func TestGetBaseHashCache(t *testing.T) {
	op := cdao.NewBaseRedisOpWithKT("hash:A:2", time.Second*10)

	for {
		b := &BB{}
		err := GetBaseHashCache(context.TODO(), op, b, getDB)
		if err != nil {
			log.Error(err)
			return
		}
		fmt.Printf("%+v\n", b)
		time.Sleep(time.Second)
	}
}

func getDB(ctx context.Context, model interface{}) error {
	fmt.Println("getDB")
	bb, _ := model.(*BB)
	bb.ID = 1
	bb.A = "AAC"
	bb.B = 123
	bb.C = 123.1
	return nil
}
