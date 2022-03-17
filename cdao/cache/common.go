package cache

import (
	"context"
	"time"

	"github.com/chuan-fu/Common/cdao"
)

type (
	CheckExistsFunc func(ctx context.Context, b cdao.BaseRedisOp) (bool, error)
	DelCacheFunc    func(ctx context.Context, b cdao.BaseRedisOp) error
)

func defaultCheckExists(ctx context.Context, b cdao.BaseRedisOp) (bool, error) {
	t, err := b.TTL(ctx)
	if err != nil {
		return false, err
	}
	if t > time.Second {
		return true, nil
	}
	return false, nil
}

func defaultDelCache(ctx context.Context, b cdao.BaseRedisOp) error {
	return b.Del(ctx)
}
