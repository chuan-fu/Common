package cache

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/chuan-fu/Common/cdao"
	log "github.com/chuan-fu/Common/zlog"
)

func TestGetBaseSetCache(t *testing.T) {
	op := cdao.NewBaseRedisOpWithKT("set:A:3", time.Minute)
	data, err := GetBaseSetCache(context.TODO(), op, func(ctx context.Context) ([]string, error) {
		log.Info("getByDb")
		return []string{"A", "b"}, nil
	})
	fmt.Println(data, err)
}
